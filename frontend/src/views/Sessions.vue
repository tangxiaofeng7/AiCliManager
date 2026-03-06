<template>
  <div class="page-container sessions-page">
    <div class="page-header">
      <div>
        <h2>会话历史</h2>
        <div class="subtitle">浏览 AI CLI 工具的对话会话记录</div>
      </div>
      <el-button :icon="Refresh" @click="refreshPage" :loading="loading">刷新</el-button>
    </div>

    <div class="filter-bar">
      <el-select
        v-model="filterToolKey"
        placeholder="选择 CLI 工具"
        style="width: 160px"
        @change="onToolChange"
      >
        <el-option label="Claude Code" value="claude" />
        <el-option label="Codex" value="codex" />
        <el-option label="OpenCode" value="opencode" />
      </el-select>
      <el-select
        v-model="filterProject"
        placeholder="全部项目"
        clearable
        style="width: 280px"
        :disabled="toolUnsupported"
        @change="loadSessions"
      >
        <el-option
          v-for="p in projects"
          :key="p.dir_name"
          :label="p.path"
          :value="p.dir_name"
        >
          <span>{{ truncatePath(p.path, 50) }}</span>
          <span class="option-sub">{{ p.session_count }} 个会话</span>
        </el-option>
      </el-select>
    </div>

    <el-alert
      v-if="unsupportedMessage"
      :title="unsupportedMessage"
      type="info"
      :closable="false"
      show-icon
      style="margin-bottom: 16px"
    />

    <div class="session-list" v-loading="loading">
      <div
        v-for="sess in sessions"
        :key="sess.session_id"
        class="session-card"
        @click="openSession(sess)"
      >
        <div class="session-header">
          <div class="session-title-row">
            <el-tag size="small" type="primary" effect="plain" class="session-tool-tag">
              {{ toolLabel(sess.cli_tool_key) }}
            </el-tag>
            <span class="session-slug mono" v-if="sess.slug">{{ sess.slug }}</span>
            <span class="session-id mono" v-else>{{ sess.session_id.slice(0, 8) }}</span>
          </div>
          <span class="session-time">{{ formatTime(sess.last_active_at || sess.started_at) }}</span>
        </div>

        <div class="session-preview" v-if="sess.first_message">
          {{ sess.first_message }}
        </div>

        <div class="session-meta">
          <span class="meta-item" v-if="sess.model">
            <el-icon><Cpu /></el-icon>
            {{ sess.model }}
          </span>
          <span class="meta-item">
            <el-icon><ChatDotRound /></el-icon>
            {{ sess.user_count }} 轮对话
          </span>
          <span class="meta-item">
            <el-icon><FolderOpened /></el-icon>
            {{ truncatePath(sess.project, 30) }}
          </span>
        </div>
      </div>

      <el-empty
        v-if="!loading && sessions.length === 0"
        :description="unsupportedMessage || '暂无对话会话记录'"
      />
    </div>

    <el-dialog
      v-model="detailDialogVisible"
      :title="dialogTitle"
      width="780px"
      top="5vh"
      destroy-on-close
      class="session-detail-dialog"
    >
      <div class="detail-meta" v-if="selectedSession">
        <el-descriptions :column="2" size="small" border>
          <el-descriptions-item label="会话 ID">
            <span class="mono">{{ selectedSession.session_id }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="模型">
            {{ selectedSession.model || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="项目目录">
            <span class="mono">{{ selectedSession.project }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="消息数">
            {{ selectedSession.message_count }}（用户 {{ selectedSession.user_count }} / 助手 {{ selectedSession.assistant_count }}）
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <div class="message-list" v-loading="messagesLoading">
        <div
          v-for="msg in messages"
          :key="msg.uuid"
          class="message-item"
          :class="'msg-' + msg.type"
        >
          <div class="msg-header">
            <el-tag
              :type="msgTagType(msg.type)"
              size="small"
              effect="dark"
            >
              {{ msgLabel(msg.type) }}
            </el-tag>
            <span class="msg-time">{{ formatTime(msg.timestamp) }}</span>
            <span class="msg-tokens mono" v-if="msg.tokens_in || msg.tokens_out">
              ↓{{ msg.tokens_in }} ↑{{ msg.tokens_out }}
            </span>
          </div>
          <div class="msg-content">{{ msg.content }}</div>
        </div>

        <el-empty
          v-if="!messagesLoading && messages.length === 0"
          description="暂无消息"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Refresh, Cpu, ChatDotRound, FolderOpened,
} from '@element-plus/icons-vue'
import type { CliSession, CliSessionMessage, CliSessionProject } from '../types'
import { getCliSessions, getCliSessionMessages, getCliSessionProjects } from '../api'

const loading = ref(false)
const messagesLoading = ref(false)
const sessions = ref<CliSession[]>([])
const projects = ref<CliSessionProject[]>([])
const messages = ref<CliSessionMessage[]>([])
const filterToolKey = ref('claude')
const filterProject = ref('')
const detailDialogVisible = ref(false)
const selectedSession = ref<CliSession | null>(null)
const dialogTitle = ref('')
const unsupportedMessage = ref('')

const toolUnsupported = computed(() => unsupportedMessage.value !== '')

const toolLabel = (key: string): string => {
  const map: Record<string, string> = {
    claude: 'Claude Code',
    codex: 'Codex',
    opencode: 'OpenCode',
  }
  return map[key] ?? key
}

const truncatePath = (path: string, max = 40): string => {
  if (!path) return '-'
  if (path.length <= max) return path
  return '...' + path.slice(-(max - 3))
}

const formatTime = (iso: string): string => {
  if (!iso) return '-'
  try {
    return new Date(iso).toLocaleString('zh-CN', { hour12: false })
  } catch {
    return iso
  }
}

const msgTagType = (type: string) => {
  if (type === 'user') return 'primary'
  if (type === 'assistant') return 'success'
  return 'info'
}

const msgLabel = (type: string) => {
  if (type === 'user') return '用户'
  if (type === 'assistant') return '助手'
  return '系统'
}

const loadProjects = async () => {
  try {
    projects.value = await getCliSessionProjects(filterToolKey.value)
    unsupportedMessage.value = ''
  } catch (err) {
    projects.value = []
    sessions.value = []
    unsupportedMessage.value = String(err)
  }
}

const loadSessions = async () => {
  if (toolUnsupported.value) {
    sessions.value = []
    return
  }

  loading.value = true
  try {
    sessions.value = await getCliSessions({
      cli_tool_key: filterToolKey.value,
      project: filterProject.value || undefined,
      limit: 100,
    })
  } catch (err) {
    sessions.value = []
    ElMessage.error(`加载会话失败：${String(err)}`)
  } finally {
    loading.value = false
  }
}

const refreshPage = async () => {
  loading.value = true
  try {
    await loadProjects()
    if (!toolUnsupported.value) {
      await loadSessions()
    }
  } finally {
    loading.value = false
  }
}

const onToolChange = async () => {
  filterProject.value = ''
  await refreshPage()
}

const openSession = async (sess: CliSession) => {
  selectedSession.value = sess
  dialogTitle.value = sess.slug || `会话 ${sess.session_id.slice(0, 8)}...`
  detailDialogVisible.value = true
  messagesLoading.value = true
  messages.value = []

  try {
    messages.value = await getCliSessionMessages(sess.cli_tool_key, sess.session_id)
  } catch (err) {
    ElMessage.error(`加载消息失败：${String(err)}`)
  } finally {
    messagesLoading.value = false
  }
}

onMounted(async () => {
  await refreshPage()
})
</script>

<style scoped>
.sessions-page {
  display: flex;
  flex-direction: column;
}

.filter-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
}

.option-sub {
  color: var(--text-muted);
  font-size: 12px;
  margin-left: 8px;
}

.session-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
  overflow-y: auto;
}

.session-card {
  background-color: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 14px 18px;
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.session-card:hover {
  border-color: var(--accent-color);
  box-shadow: 0 0 0 1px rgba(64, 158, 255, 0.15);
}

.session-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
}

.session-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.session-tool-tag {
  flex-shrink: 0;
}

.session-slug,
.session-id {
  font-size: 13px;
  color: var(--text-primary);
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.session-time {
  font-size: 12px;
  color: var(--text-muted);
  flex-shrink: 0;
}

.session-preview {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  margin-bottom: 8px;
  word-break: break-all;
}

.session-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--text-muted);
}

.meta-item .el-icon {
  font-size: 13px;
}

.detail-meta {
  margin-bottom: 16px;
}

.message-list {
  max-height: 60vh;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message-item {
  border-radius: 8px;
  padding: 12px 14px;
}

.msg-user {
  background-color: var(--el-color-primary-light-9);
}

.msg-assistant {
  background-color: var(--el-color-success-light-9);
}

.msg-system {
  background-color: var(--el-fill-color-light);
}

.msg-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.msg-time {
  font-size: 11px;
  color: var(--text-muted);
}

.msg-tokens {
  font-size: 11px;
  color: var(--text-muted);
  margin-left: auto;
}

.msg-content {
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-primary);
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 400px;
  overflow-y: auto;
}

:deep(.session-detail-dialog .el-dialog__body) {
  padding-top: 12px;
}
</style>
