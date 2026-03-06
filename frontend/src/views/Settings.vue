<template>
  <div class="page-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <div>
        <h2>全局设置</h2>
        <div class="subtitle">配置应用全局参数，如默认终端、主题、会话保留策略等</div>
      </div>
      <el-button type="primary" :icon="Check" @click="saveAllSettings" :loading="saving">
        保存设置
      </el-button>
    </div>

    <!-- 设置表单 -->
    <div class="settings-container" v-loading="loading">
      <!-- 终端设置 -->
      <div class="settings-section">
        <div class="section-title">终端设置</div>
        <el-form :model="settings" label-width="160px">
          <el-form-item label="默认终端">
            <el-select v-model="settings.default_terminal" style="width: 220px">
              <el-option label="系统默认" value="default" />
              <el-option label="Windows Terminal (wt)" value="wt" />
              <el-option label="PowerShell" value="powershell" />
              <el-option label="CMD" value="cmd" />
              <el-option label="WSL" value="wsl" />
              <el-option label="iTerm2 (macOS)" value="iterm2" />
              <el-option label="Terminal.app (macOS)" value="terminal" />
              <el-option label="GNOME Terminal (Linux)" value="gnome-terminal" />
              <el-option label="Tmux (Linux)" value="tmux" />
            </el-select>
            <div class="form-hint">各 CLI 工具未单独指定终端时使用此默认终端</div>
          </el-form-item>
        </el-form>
      </div>

      <el-divider />

      <!-- 会话历史设置 -->
      <div class="settings-section">
        <div class="section-title">会话历史</div>
        <el-form :model="settings" label-width="160px">
          <el-form-item label="最大保留数量">
            <el-input-number
              v-model="settings.max_sessions"
              :min="10"
              :max="10000"
              :step="50"
              style="width: 160px"
            />
            <div class="form-hint">超出此数量时自动删除最旧记录（0 表示不限制）</div>
          </el-form-item>
          <el-form-item label="自动清理天数">
            <el-input-number
              v-model="settings.session_retention_days"
              :min="0"
              :max="365"
              :step="7"
              style="width: 160px"
            />
            <div class="form-hint">超过指定天数的会话记录自动清理（0 表示不自动清理）</div>
          </el-form-item>
        </el-form>
      </div>

      <el-divider />

      <!-- 外观设置 -->
      <div class="settings-section">
        <div class="section-title">外观</div>
        <el-form :model="settings" label-width="160px">
          <el-form-item label="主题">
            <el-radio-group v-model="settings.theme">
              <el-radio-button value="dark">暗色</el-radio-button>
              <el-radio-button value="light">亮色（开发中）</el-radio-button>
              <el-radio-button value="system">跟随系统</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="字体">
            <el-select v-model="settings.font_family" style="width: 260px">
              <el-option label="Cascadia Code" value="'Cascadia Code'" />
              <el-option label="JetBrains Mono" value="'JetBrains Mono'" />
              <el-option label="Consolas" value="Consolas" />
              <el-option label="Courier New" value="'Courier New'" />
              <el-option label="Fira Code" value="'Fira Code'" />
            </el-select>
          </el-form-item>
        </el-form>
      </div>

      <el-divider />

      <!-- 安全设置 -->
      <div class="settings-section">
        <div class="section-title">安全</div>
        <el-form :model="settings" label-width="160px">
          <el-form-item label="写入前备份配置">
            <el-switch v-model="settings.backup_before_write" />
            <div class="form-hint">同步配置到 CLI 工具前，自动备份原始配置文件</div>
          </el-form-item>
          <el-form-item label="启动路径校验">
            <el-switch v-model="settings.validate_executable_path" />
            <div class="form-hint">启动前校验 CLI 工具的可执行文件路径，防止路径注入</div>
          </el-form-item>
        </el-form>
      </div>

      <el-divider />

      <!-- 关于 -->
      <div class="settings-section">
        <div class="section-title">关于</div>
        <div class="about-info">
          <div class="about-row">
            <span class="about-label">应用版本</span>
            <span class="about-value mono">v0.1.0</span>
          </div>
          <div class="about-row">
            <span class="about-label">技术栈</span>
            <span class="about-value">Wails v2.11 + Go 1.24 + Vue 3 + SQLite</span>
          </div>
          <div class="about-row">
            <span class="about-label">配置数据库</span>
            <span class="about-value mono">{{ dbPath }}</span>
          </div>
          <div class="about-row">
            <span class="about-label">项目地址</span>
            <span class="about-value">AiCliManager</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Check } from '@element-plus/icons-vue'
import { getSettings, saveSettings } from '../api'

// ---- 状态 ----
const loading = ref(false)
const saving = ref(false)
const dbPath = ref('~/.aiclimgr/data.db')

// 设置项，key 与后端 settings 表的 key 一一对应
const settings = reactive({
  default_terminal: 'default',
  max_sessions: 500,
  session_retention_days: 30,
  theme: 'dark',
  font_family: "'Cascadia Code'",
  backup_before_write: true,
  validate_executable_path: true,
})

// ---- 数据加载 ----
const loadSettings = async () => {
  loading.value = true
  try {
    const data = await getSettings().catch(() => ({} as Record<string, string>))
    if (!data) return

    if (data.default_terminal) settings.default_terminal = data.default_terminal
    if (data.max_sessions) settings.max_sessions = parseInt(data.max_sessions, 10)
    if (data.session_retention_days) settings.session_retention_days = parseInt(data.session_retention_days, 10)
    if (data.theme) settings.theme = data.theme
    if (data.font_family) settings.font_family = data.font_family
    if (data.backup_before_write !== undefined) settings.backup_before_write = data.backup_before_write === 'true'
    if (data.validate_executable_path !== undefined) settings.validate_executable_path = data.validate_executable_path === 'true'
    if (data.db_path) dbPath.value = data.db_path
  } finally {
    loading.value = false
  }
}

// ---- 保存设置 ----
const saveAllSettings = async () => {
  saving.value = true
  try {
    await saveSettings({
      default_terminal: settings.default_terminal,
      max_sessions: String(settings.max_sessions),
      session_retention_days: String(settings.session_retention_days),
      theme: settings.theme,
      font_family: settings.font_family,
      backup_before_write: String(settings.backup_before_write),
      validate_executable_path: String(settings.validate_executable_path),
    })
    ElMessage.success('设置已保存')
  } catch (err) {
    ElMessage.error(`保存失败：${err}`)
  } finally {
    saving.value = false
  }
}

onMounted(loadSettings)
</script>

<style scoped>
.settings-container {
  max-width: 720px;
}

.settings-section {
  padding: 8px 0;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
  padding-left: 8px;
  border-left: 3px solid var(--accent-color);
}

.form-hint {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
  line-height: 1.4;
}

/* 关于信息 */
.about-info {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 16px;
  background-color: var(--bg-card);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.about-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

.about-label {
  width: 100px;
  font-size: 13px;
  color: var(--text-muted);
  flex-shrink: 0;
}

.about-value {
  font-size: 13px;
  color: var(--text-primary);
}
</style>
