<template>
  <div class="page-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <div>
        <h2>MCP Servers</h2>
        <div class="subtitle">管理 Model Context Protocol 服务器，按需分配给各 CLI 工具</div>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreateDialog">新建 MCP Server</el-button>
    </div>

    <!-- MCP Server 列表 -->
    <el-table :data="mcpServers" v-loading="loading" border stripe style="width: 100%">
      <el-table-column label="状态" width="70">
        <template #default="{ row }">
          <el-switch :model-value="!!row.is_enabled" @change="toggleEnabled(row)" />
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column label="类型" width="90">
        <template #default="{ row }">
          <el-tag
            size="small"
            :type="row.type === 'stdio' ? 'primary' : row.type === 'sse' ? 'success' : 'warning'"
          >
            {{ row.type }}
          </el-tag>
        </template>
      </el-table-column>
      <!-- stdio 类型显示 command -->
      <el-table-column label="命令 / URL" min-width="220">
        <template #default="{ row }">
          <template v-if="row.type === 'stdio'">
            <span class="mono">{{ row.command }}</span>
            <span class="text-muted mono" v-if="row.args">
              {{ formatArgs(row.args) }}
            </span>
          </template>
          <template v-else>
            <span class="mono">{{ row.url }}</span>
          </template>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="160">
        <template #default="{ row }">
          <span class="text-muted">{{ row.description || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" :icon="Edit" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" :icon="Delete" @click="confirmDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新建 / 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingServer ? '编辑 MCP Server' : '新建 MCP Server'"
      width="580px"
      destroy-on-close
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="如：Filesystem MCP" clearable />
        </el-form-item>

        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio-button value="stdio">stdio（本地进程）</el-radio-button>
            <el-radio-button value="sse">SSE（HTTP 流）</el-radio-button>
            <el-radio-button value="http">HTTP</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <!-- stdio 专属字段 -->
        <template v-if="form.type === 'stdio'">
          <el-form-item label="启动命令" prop="command">
            <el-input v-model="form.command" placeholder="如：npx 或 python3" clearable />
          </el-form-item>
          <el-form-item label="命令参数">
            <el-input
              v-model="form.args"
              type="textarea"
              :rows="3"
              placeholder='JSON 数组，如：["-y", "@modelcontextprotocol/server-filesystem", "/workspace"]'
            />
          </el-form-item>
          <el-form-item label="环境变量">
            <el-input
              v-model="form.env"
              type="textarea"
              :rows="2"
              placeholder='JSON 对象，如：{"API_KEY": "sk-..."}'
            />
          </el-form-item>
        </template>

        <!-- sse / http 专属字段 -->
        <template v-else>
          <el-form-item label="服务地址" prop="url">
            <el-input v-model="form.url" placeholder="如：http://localhost:3000/sse" clearable />
          </el-form-item>
        </template>

        <el-form-item label="描述">
          <el-input v-model="form.description" placeholder="用途描述（可选）" clearable />
        </el-form-item>

        <el-form-item label="启用">
          <el-switch v-model="form.is_enabled" />
        </el-form-item>

        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" :max="9999" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ editingServer ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import type { McpServer } from '../types'
import { getMcpServers, createMcpServer, updateMcpServer, deleteMcpServer } from '../api'

// ---- 状态 ----
const loading = ref(false)
const submitting = ref(false)
const mcpServers = ref<McpServer[]>([])
const dialogVisible = ref(false)
const editingServer = ref<McpServer | null>(null)
const formRef = ref<FormInstance>()

const form = reactive({
  name: '',
  type: 'stdio' as 'stdio' | 'sse' | 'http',
  command: '',
  args: '',
  env: '',
  url: '',
  description: '',
  is_enabled: true,
  sort_order: 0,
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  command: [
    {
      validator: (_rule, value, callback) => {
        if (form.type === 'stdio' && !value) callback(new Error('stdio 类型必须填写启动命令'))
        else callback()
      },
      trigger: 'blur',
    },
  ],
  url: [
    {
      validator: (_rule, value, callback) => {
        if ((form.type === 'sse' || form.type === 'http') && !value) {
          callback(new Error('请填写服务地址'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

// ---- 辅助 ----
const formatArgs = (args: string): string => {
  try {
    const arr = JSON.parse(args) as string[]
    return ' ' + arr.join(' ')
  } catch {
    return ' ' + args
  }
}

// ---- 数据加载 ----
const loadData = async () => {
  loading.value = true
  try {
    mcpServers.value = (await getMcpServers().catch(() => [])) ?? []
  } finally {
    loading.value = false
  }
}

// ---- 对话框 ----
const openCreateDialog = () => {
  editingServer.value = null
  Object.assign(form, {
    name: '', type: 'stdio', command: '', args: '',
    env: '', url: '', description: '', is_enabled: true, sort_order: 0,
  })
  dialogVisible.value = true
}

const openEditDialog = (server: McpServer) => {
  editingServer.value = server
  Object.assign(form, {
    name: server.name,
    type: server.type,
    command: server.command ?? '',
    args: server.args ?? '',
    env: server.env ?? '',
    url: server.url ?? '',
    description: server.description ?? '',
    is_enabled: !!server.is_enabled,
    sort_order: server.sort_order,
  })
  dialogVisible.value = true
}

const submitForm = async () => {
  await formRef.value?.validate().catch(() => { throw new Error('validation') })
  submitting.value = true
  try {
    const payload = {
      ...form,
      is_enabled: form.is_enabled ? 1 : 0,
    }
    if (editingServer.value) {
      await updateMcpServer(editingServer.value.id, payload)
      ElMessage.success('MCP Server 已更新')
    } else {
      await createMcpServer(payload)
      ElMessage.success('MCP Server 已创建')
    }
    dialogVisible.value = false
    await loadData()
  } catch (err) {
    if (String(err) !== 'Error: validation') ElMessage.error(`操作失败：${err}`)
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async (server: McpServer) => {
  await ElMessageBox.confirm(
    `确认删除 MCP Server "${server.name}"？`,
    '删除确认',
    { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
  )
  try {
    await deleteMcpServer(server.id)
    ElMessage.success('MCP Server 已删除')
    await loadData()
  } catch (err) {
    ElMessage.error(`删除失败：${err}`)
  }
}

// ---- 启用 / 禁用切换 ----
const toggleEnabled = async (server: McpServer) => {
  try {
    await updateMcpServer(server.id, { is_enabled: server.is_enabled ? 0 : 1 })
    await loadData()
  } catch (err) {
    ElMessage.error(`操作失败：${err}`)
  }
}

onMounted(loadData)
</script>
