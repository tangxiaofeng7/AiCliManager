<template>
  <div class="page-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <div>
        <h2>代理配置</h2>
        <div class="subtitle">统一管理代理，启动 CLI 工具时自动注入环境变量</div>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreateDialog">新建代理</el-button>
    </div>

    <!-- 代理列表 -->
    <el-table :data="proxies" v-loading="loading" border stripe style="width: 100%">
      <el-table-column label="状态" width="70">
        <template #default="{ row }">
          <el-tooltip :content="row.is_active ? '当前全局激活' : '未激活'" placement="top">
            <el-icon :style="{ color: row.is_active ? 'var(--success-color)' : 'var(--text-muted)' }">
              <CircleCheckFilled v-if="row.is_active" />
              <CircleCheck v-else />
            </el-icon>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column label="类型" width="90">
        <template #default="{ row }">
          <el-tag size="small" :type="row.type === 'socks5' ? 'warning' : 'info'">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="地址" min-width="200">
        <template #default="{ row }">
          <span class="mono">{{ row.host }}:{{ row.port }}</span>
        </template>
      </el-table-column>
      <el-table-column label="认证" width="100">
        <template #default="{ row }">
          <span class="text-muted">{{ row.username ? row.username : '无' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="不走代理域名" min-width="160">
        <template #default="{ row }">
          <span class="text-muted mono" :title="row.no_proxy">
            {{ row.no_proxy ? truncate(row.no_proxy, 30) : '-' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="{ row }">
          <el-button
            size="small"
            :type="row.is_active ? 'success' : 'default'"
            :icon="row.is_active ? CircleCheckFilled : Check"
            @click="toggleGlobalProxy(row)"
          >
            {{ row.is_active ? '已激活' : '设为全局' }}
          </el-button>
          <el-button size="small" :icon="Edit" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" :icon="Delete" @click="confirmDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 清除全局代理 -->
    <div style="margin-top: 12px" v-if="hasActiveProxy">
      <el-button :icon="Close" @click="clearProxy">清除全局代理</el-button>
    </div>

    <!-- 新建 / 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingProxy ? '编辑代理' : '新建代理'"
      width="500px"
      destroy-on-close
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="110px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="如：本地 Clash" clearable />
        </el-form-item>

        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" style="width: 100%">
            <el-option label="HTTP" value="http" />
            <el-option label="HTTPS" value="https" />
            <el-option label="SOCKS5" value="socks5" />
          </el-select>
        </el-form-item>

        <el-form-item label="主机" prop="host">
          <el-input v-model="form.host" placeholder="如：127.0.0.1" clearable />
        </el-form-item>

        <el-form-item label="端口" prop="port">
          <el-input-number v-model="form.port" :min="1" :max="65535" style="width: 160px" />
        </el-form-item>

        <el-form-item label="用户名">
          <el-input v-model="form.username" placeholder="可选，代理认证用户名" clearable />
        </el-form-item>

        <el-form-item label="密码">
          <el-input
            v-model="form.password"
            type="password"
            :placeholder="editingProxy ? '留空则保持不变' : '可选，代理认证密码'"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item label="不走代理">
          <el-input
            v-model="form.no_proxy"
            placeholder="如：localhost,127.0.0.1,*.local（逗号分隔）"
            clearable
          />
          <div class="form-hint">不经过代理的域名/IP，逗号分隔</div>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ editingProxy ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus, Edit, Delete, CircleCheckFilled, CircleCheck, Check, Close } from '@element-plus/icons-vue'
import type { Proxy } from '../types'
import {
  getProxies, createProxy, updateProxy, deleteProxy,
  setGlobalProxy, clearGlobalProxy,
} from '../api'

// ---- 状态 ----
const loading = ref(false)
const submitting = ref(false)
const proxies = ref<Proxy[]>([])
const dialogVisible = ref(false)
const editingProxy = ref<Proxy | null>(null)
const formRef = ref<FormInstance>()

const hasActiveProxy = computed(() => proxies.value.some(p => p.is_active))

const form = reactive({
  name: '',
  type: 'http' as 'http' | 'https' | 'socks5',
  host: '127.0.0.1',
  port: 7890,
  username: '',
  password: '',
  no_proxy: 'localhost,127.0.0.1',
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入代理名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  host: [{ required: true, message: '请输入主机地址', trigger: 'blur' }],
  port: [{ required: true, message: '请输入端口', trigger: 'blur' }],
}

const truncate = (str: string, len: number) =>
  str.length > len ? str.slice(0, len) + '…' : str

// ---- 数据加载 ----
const loadProxies = async () => {
  loading.value = true
  try {
    proxies.value = (await getProxies().catch(() => [])) ?? []
  } finally {
    loading.value = false
  }
}

// ---- 对话框 ----
const openCreateDialog = () => {
  editingProxy.value = null
  Object.assign(form, {
    name: '', type: 'http', host: '127.0.0.1', port: 7890,
    username: '', password: '', no_proxy: 'localhost,127.0.0.1',
  })
  dialogVisible.value = true
}

const openEditDialog = (proxy: Proxy) => {
  editingProxy.value = proxy
  Object.assign(form, {
    name: proxy.name, type: proxy.type, host: proxy.host,
    port: proxy.port, username: proxy.username ?? '',
    password: '', // 密码不回填
    no_proxy: proxy.no_proxy ?? '',
  })
  dialogVisible.value = true
}

const submitForm = async () => {
  await formRef.value?.validate().catch(() => { throw new Error('validation') })
  submitting.value = true
  try {
    if (editingProxy.value) {
      const req: Record<string, unknown> = {
        name: form.name, type: form.type, host: form.host,
        port: form.port, username: form.username, no_proxy: form.no_proxy,
      }
      if (form.password) req.password = form.password
      await updateProxy(editingProxy.value.id, req as Parameters<typeof updateProxy>[1])
      ElMessage.success('代理配置已更新')
    } else {
      await createProxy({ ...form })
      ElMessage.success('代理配置已创建')
    }
    dialogVisible.value = false
    await loadProxies()
  } catch (err) {
    if (String(err) !== 'Error: validation') ElMessage.error(`操作失败：${err}`)
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async (proxy: Proxy) => {
  await ElMessageBox.confirm(
    `确认删除代理 "${proxy.name}"？`,
    '删除确认',
    { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
  )
  try {
    await deleteProxy(proxy.id)
    ElMessage.success('代理已删除')
    await loadProxies()
  } catch (err) {
    ElMessage.error(`删除失败：${err}`)
  }
}

// ---- 全局代理 ----
const toggleGlobalProxy = async (proxy: Proxy) => {
  if (proxy.is_active) return // 已激活不重复操作
  try {
    await setGlobalProxy(proxy.id)
    ElMessage.success(`已将 "${proxy.name}" 设为全局代理`)
    await loadProxies()
  } catch (err) {
    ElMessage.error(`设置失败：${err}`)
  }
}

const clearProxy = async () => {
  try {
    await clearGlobalProxy()
    ElMessage.success('全局代理已清除')
    await loadProxies()
  } catch (err) {
    ElMessage.error(`清除失败：${err}`)
  }
}

onMounted(loadProxies)
</script>

<style scoped>
.form-hint {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}
</style>
