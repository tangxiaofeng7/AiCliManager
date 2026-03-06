<template>
  <div class="page-container dashboard">
    <div class="page-header">
      <div>
        <h2>启动面板</h2>
        <div class="subtitle">选择 AI CLI 工具，配置并一键启动</div>
      </div>
      <el-button type="primary" :icon="Refresh" @click="reloadDashboard" :loading="loading">
        刷新检测
      </el-button>
    </div>

    <div class="tools-grid" v-loading="loading">
      <div
        v-for="tool in cliTools"
        :key="tool.id"
        class="tool-card"
        :class="{ disabled: !tool.is_enabled }"
      >
        <div class="tool-header">
          <div class="tool-info">
            <div class="tool-name">{{ tool.name }}</div>
            <div class="tool-key mono">{{ tool.key }}</div>
          </div>
          <el-tag
            :type="tool.is_installed ? 'success' : 'danger'"
            size="small"
            effect="dark"
          >
            {{ tool.is_installed ? '已安装' : '未安装' }}
          </el-tag>
        </div>

        <div class="tool-config">
          <div class="config-row" v-if="activeConfigs[tool.key]?.profile_id">
            <el-icon class="config-icon"><User /></el-icon>
            <span class="config-value">
              {{ getProfileName(activeConfigs[tool.key]?.profile_id) }}
            </span>
          </div>
          <div class="config-row" v-else>
            <el-icon class="config-icon"><InfoFilled /></el-icon>
            <span class="config-value text-muted">未选择 Profile</span>
          </div>

          <div class="config-row" v-if="activeConfigs[tool.key]?.proxy_id">
            <el-icon class="config-icon"><Share /></el-icon>
            <span class="config-value">
              {{ getProxyName(activeConfigs[tool.key]?.proxy_id) }}
            </span>
          </div>
          <div class="config-row" v-else>
            <el-icon class="config-icon"><Share /></el-icon>
            <span class="config-value text-muted">不使用代理</span>
          </div>

          <div class="config-row" v-if="tool.executable">
            <el-icon class="config-icon"><FolderOpened /></el-icon>
            <span class="config-value mono" :title="tool.executable">
              {{ truncatePath(tool.executable) }}
            </span>
          </div>
        </div>

        <div class="tool-actions">
          <el-button
            type="primary"
            :icon="VideoPlay"
            :disabled="!tool.is_installed"
            @click="openLaunchDialog(tool)"
            class="launch-btn"
          >
            启动
          </el-button>
          <el-button
            :icon="Setting"
            circle
            @click="openConfigDialog(tool)"
            title="配置默认参数"
          />
        </div>
      </div>

      <el-empty v-if="!loading && cliTools.length === 0" description="未检测到任何 CLI 工具" />
    </div>

    <el-dialog
      v-model="launchDialogVisible"
      :title="`启动 ${selectedTool?.name ?? ''}`"
      width="560px"
      destroy-on-close
    >
      <el-form :model="launchForm" label-width="90px" class="launch-form">
        <el-form-item label="Profile">
          <el-select
            v-model="launchForm.profile_id"
            placeholder="请选择 Profile"
            style="width: 100%"
            clearable
          >
            <el-option
              v-for="p in profiles"
              :key="p.id"
              :label="p.name"
              :value="p.id"
            >
              <span>{{ p.name }}</span>
              <span class="option-sub mono">{{ p.model }}</span>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="代理">
          <el-select
            v-model="launchForm.proxy_id"
            placeholder="不使用代理"
            style="width: 100%"
            clearable
          >
            <el-option
              v-for="p in proxies"
              :key="p.id"
              :label="p.name"
              :value="p.id"
            >
              <span>{{ p.name }}</span>
              <span class="option-sub mono">{{ p.type }}://{{ p.host }}:{{ p.port }}</span>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="MCP Servers">
          <el-select
            v-model="launchForm.mcp_server_ids"
            multiple
            placeholder="选择 MCP Server（可多选）"
            style="width: 100%"
            collapse-tags
            collapse-tags-tooltip
          >
            <el-option
              v-for="m in mcpServers"
              :key="m.id"
              :label="m.name"
              :value="m.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="终端">
          <el-select
            v-model="launchForm.terminal"
            placeholder="使用工具默认终端"
            style="width: 100%"
            clearable
          >
            <el-option
              v-for="t in availableTerminals"
              :key="t.id"
              :label="t.name"
              :value="t.id"
              :disabled="!t.is_available"
            >
              <span>{{ t.name }}</span>
              <el-tag v-if="!t.is_available" type="info" size="small" style="margin-left: 8px">不可用</el-tag>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="工作目录">
          <el-input
            v-model="launchForm.working_dir"
            placeholder="留空则使用用户 Home 目录"
            clearable
          />
        </el-form-item>

        <el-form-item label="额外参数">
          <el-input
            v-model="extraArgsInput"
            placeholder="如：--verbose --no-auto-update（空格分隔）"
            clearable
          />
          <div class="form-hint">以空格分隔的额外 CLI 启动参数</div>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="launchDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :icon="VideoPlay"
          :loading="launching"
          @click="confirmLaunch"
        >
          启动
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="configDialogVisible"
      :title="`${selectedTool?.name ?? ''} 默认配置`"
      width="480px"
      destroy-on-close
    >
      <el-form :model="configForm" label-width="90px">
        <el-form-item label="默认 Profile">
          <el-select v-model="configForm.profile_id" placeholder="选择默认 Profile" style="width: 100%" clearable>
            <el-option v-for="p in profiles" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="默认代理">
          <el-select v-model="configForm.proxy_id" placeholder="不使用代理" style="width: 100%" clearable>
            <el-option v-for="p in proxies" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveActiveConfig">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Refresh, VideoPlay, Setting, User, Share,
  InfoFilled, FolderOpened,
} from '@element-plus/icons-vue'
import type { CliTool, Profile, Proxy, McpServer, TerminalInfo, ActiveConfig } from '../types'
import {
  getCliTools, detectCliTool, getProfiles, getProxies, getMcpServers,
  listAvailableTerminals, launchCliTool, getCliToolActiveConfig,
  setCliToolActiveConfig,
} from '../api'

const loading = ref(false)
const launching = ref(false)
const cliTools = ref<CliTool[]>([])
const profiles = ref<Profile[]>([])
const proxies = ref<Proxy[]>([])
const mcpServers = ref<McpServer[]>([])
const availableTerminals = ref<TerminalInfo[]>([])
const activeConfigs = ref<Record<string, ActiveConfig | null>>({})

const launchDialogVisible = ref(false)
const selectedTool = ref<CliTool | null>(null)
const extraArgsInput = ref('')
const launchForm = ref({
  profile_id: 0,
  proxy_id: null as number | null,
  mcp_server_ids: [] as number[],
  terminal: '',
  working_dir: '',
})

const configDialogVisible = ref(false)
const configForm = ref({
  profile_id: null as number | null,
  proxy_id: null as number | null,
})

const getProfileName = (id: number | null | undefined): string => {
  if (!id) return '未选择'
  return profiles.value.find(p => p.id === id)?.name ?? `Profile #${id}`
}

const getProxyName = (id: number | null | undefined): string => {
  if (!id) return ''
  const p = proxies.value.find(p => p.id === id)
  return p ? `${p.name} (${p.type}://${p.host}:${p.port})` : `Proxy #${id}`
}

const truncatePath = (path: string): string => {
  if (path.length <= 40) return path
  return '...' + path.slice(-37)
}

const parseExtraArgs = (input: string): string[] => {
  return input.trim().split(/\s+/).filter(Boolean)
}

const loadMetaData = async () => {
  const [pf, px, mcp, terms] = await Promise.all([
    getProfiles(),
    getProxies(),
    getMcpServers(),
    listAvailableTerminals(),
  ])
  profiles.value = pf
  proxies.value = px
  mcpServers.value = mcp
  availableTerminals.value = terms
}

const loadActiveConfigs = async (tools: CliTool[]) => {
  const entries = await Promise.all(
    tools.map(async (tool) => [tool.key, await getCliToolActiveConfig(tool.key)] as const),
  )
  activeConfigs.value = Object.fromEntries(entries)
}

const loadCliToolsWithDetection = async () => {
  const tools = await getCliTools()
  await Promise.all(tools.map(tool => detectCliTool(tool.key)))
  const updatedTools = await getCliTools()
  cliTools.value = updatedTools
  await loadActiveConfigs(updatedTools)
}

const reloadDashboard = async () => {
  loading.value = true
  try {
    await Promise.all([loadMetaData(), loadCliToolsWithDetection()])
  } catch (err) {
    ElMessage.error(`加载启动面板失败：${String(err)}`)
  } finally {
    loading.value = false
  }
}

const openLaunchDialog = (tool: CliTool) => {
  selectedTool.value = tool
  const cfg = activeConfigs.value[tool.key]
  launchForm.value = {
    profile_id: cfg?.profile_id ?? 0,
    proxy_id: cfg?.proxy_id ?? null,
    mcp_server_ids: [],
    terminal: tool.preferred_terminal ?? '',
    working_dir: '',
  }
  extraArgsInput.value = ''
  launchDialogVisible.value = true
}

const confirmLaunch = async () => {
  if (!selectedTool.value) return
  if (!launchForm.value.profile_id) {
    ElMessage.warning('请先选择一个 Profile')
    return
  }

  launching.value = true
  try {
    await launchCliTool({
      cli_tool_key: selectedTool.value.key,
      profile_id: launchForm.value.profile_id,
      proxy_id: launchForm.value.proxy_id,
      mcp_server_ids: launchForm.value.mcp_server_ids,
      skill_ids: [],
      skill_vars: {},
      extra_args: parseExtraArgs(extraArgsInput.value),
      terminal: launchForm.value.terminal,
      working_dir: launchForm.value.working_dir,
    })

    ElMessage.success(`${selectedTool.value.name} 已成功启动`)
    launchDialogVisible.value = false
  } catch (err) {
    ElMessage.error(`启动失败：${String(err)}`)
  } finally {
    launching.value = false
  }
}

const openConfigDialog = (tool: CliTool) => {
  selectedTool.value = tool
  const cfg = activeConfigs.value[tool.key]
  configForm.value = {
    profile_id: cfg?.profile_id ?? null,
    proxy_id: cfg?.proxy_id ?? null,
  }
  configDialogVisible.value = true
}

const saveActiveConfig = async () => {
  if (!selectedTool.value) return
  try {
    await setCliToolActiveConfig({
      cli_tool_key: selectedTool.value.key,
      profile_id: configForm.value.profile_id,
      proxy_id: configForm.value.proxy_id,
    })
    activeConfigs.value[selectedTool.value.key] = await getCliToolActiveConfig(selectedTool.value.key)
    ElMessage.success('默认配置已保存')
    configDialogVisible.value = false
  } catch (err) {
    ElMessage.error(`保存失败：${String(err)}`)
  }
}

onMounted(async () => {
  await reloadDashboard()
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
  flex: 1;
}

.tool-card {
  background-color: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.tool-card:hover {
  border-color: var(--accent-color);
  box-shadow: 0 0 0 1px rgba(64, 158, 255, 0.2);
}

.tool-card.disabled {
  opacity: 0.5;
  pointer-events: none;
}

.tool-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.tool-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.tool-key {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 2px;
}

.tool-config {
  display: flex;
  flex-direction: column;
  gap: 6px;
  background-color: var(--bg-secondary);
  border-radius: 6px;
  padding: 10px 12px;
}

.config-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.config-icon {
  color: var(--text-muted);
  flex-shrink: 0;
  font-size: 14px;
}

.config-value {
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tool-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.launch-btn {
  flex: 1;
}

.option-sub {
  color: var(--text-muted);
  font-size: 12px;
  margin-left: 8px;
}

.form-hint {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}
</style>
