<template>
  <div class="page-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <div>
        <h2>Provider 管理</h2>
        <div class="subtitle">管理 API 提供商，供所有 CLI 工具的 Profile 复用</div>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreateDialog">新建 Provider</el-button>
    </div>

    <!-- Provider 表格 -->
    <el-table :data="providers" v-loading="loading" border stripe style="width: 100%">
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column prop="type" label="类型" width="110">
        <template #default="{ row }">
          <el-tag :type="typeTagType(row.type)" size="small">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="api_url" label="API Endpoint" min-width="200">
        <template #default="{ row }">
          <span class="mono">{{ row.api_url }}</span>
        </template>
      </el-table-column>
      <el-table-column label="API Key" width="130">
        <template #default="{ row }">
          <span class="mono text-muted">{{ maskApiKey(row.api_key) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="支持模型" min-width="120">
        <template #default="{ row }">
          <span class="text-muted">{{ parseModelsCount(row.models) }} 个模型</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="{ row }">
          <el-button size="small" :icon="Connection" @click="testProvider(row)" :loading="testingId === row.id">
            测试
          </el-button>
          <el-button size="small" :icon="Download" @click="fetchModels(row)" :loading="fetchingId === row.id">
            拉取模型
          </el-button>
          <el-button size="small" :icon="Edit" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" :icon="Delete" @click="confirmDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 测试结果展示 -->
    <el-alert
      v-if="testResult"
      :title="testResult.success ? `连通性测试成功（延迟 ${testResult.latency_ms}ms）` : `连通性测试失败：${testResult.message}`"
      :type="testResult.success ? 'success' : 'error'"
      show-icon
      closable
      @close="testResult = null"
      style="margin-top: 16px"
    />

    <!-- 新建 / 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingProvider ? '编辑 Provider' : '新建 Provider'"
      width="520px"
      destroy-on-close
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="如：Anthropic 官方" clearable />
        </el-form-item>

        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" style="width: 100%">
            <el-option label="Anthropic" value="anthropic" />
            <el-option label="OpenAI 兼容" value="openai" />
            <el-option label="自定义 Endpoint" value="custom" />
          </el-select>
        </el-form-item>

        <el-form-item label="API Endpoint" prop="api_url">
          <el-input v-model="form.api_url" placeholder="如：https://api.anthropic.com" clearable />
        </el-form-item>

        <el-form-item label="API Key" prop="api_key">
          <el-input
            v-model="form.api_key"
            type="password"
            :placeholder="editingProvider ? '留空则保持不变' : '请输入 API Key'"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item label="模型列表">
          <el-input
            v-model="form.models"
            type="textarea"
            :rows="3"
            placeholder='JSON 数组，如：["claude-3-5-sonnet-20241022","claude-3-haiku-20240307"]'
          />
          <div class="form-hint">可手动填写或点击"拉取模型"自动获取</div>
        </el-form-item>

        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" :max="9999" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ editingProvider ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus, Edit, Delete, Connection, Download } from '@element-plus/icons-vue'
import type { Provider, TestResult } from '../types'
import {
  getProviders, createProvider, updateProvider,
  deleteProvider, testProvider as apiTestProvider,
  fetchProviderModels,
} from '../api'

// ---- 状态 ----
const loading = ref(false)
const submitting = ref(false)
const testingId = ref<number | null>(null)
const fetchingId = ref<number | null>(null)
const providers = ref<Provider[]>([])
const testResult = ref<TestResult | null>(null)

// 对话框
const dialogVisible = ref(false)
const editingProvider = ref<Provider | null>(null)
const formRef = ref<FormInstance>()

// 表单数据
const form = reactive({
  name: '',
  type: 'anthropic' as 'anthropic' | 'openai' | 'custom',
  api_url: '',
  api_key: '',
  models: '',
  sort_order: 0,
})

// 表单验证规则
const rules: FormRules = {
  name: [{ required: true, message: '请输入 Provider 名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  api_url: [
    { required: true, message: '请输入 API Endpoint', trigger: 'blur' },
    { type: 'url', message: '请输入有效的 URL', trigger: 'blur' },
  ],
  api_key: [
    {
      validator: (_rule, value, callback) => {
        if (!editingProvider.value && !value) {
          callback(new Error('请输入 API Key'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

// ---- 辅助函数 ----
const typeTagType = (type: string) => {
  if (type === 'anthropic') return 'warning'
  if (type === 'openai') return 'success'
  return 'info'
}

const maskApiKey = (key: string): string => {
  if (!key || key.length <= 8) return '••••••••'
  return key.slice(0, 4) + '••••' + key.slice(-4)
}

const parseModelsCount = (models: string): number => {
  try {
    const arr = JSON.parse(models)
    return Array.isArray(arr) ? arr.length : 0
  } catch {
    return 0
  }
}

// ---- 数据加载 ----
const loadProviders = async () => {
  loading.value = true
  try {
    providers.value = (await getProviders().catch(() => [])) ?? []
  } finally {
    loading.value = false
  }
}

// ---- 对话框操作 ----
const openCreateDialog = () => {
  editingProvider.value = null
  Object.assign(form, { name: '', type: 'anthropic', api_url: '', api_key: '', models: '', sort_order: 0 })
  dialogVisible.value = true
}

const openEditDialog = (provider: Provider) => {
  editingProvider.value = provider
  Object.assign(form, {
    name: provider.name,
    type: provider.type,
    api_url: provider.api_url,
    api_key: '', // 编辑时不回填 API Key
    models: provider.models ?? '',
    sort_order: provider.sort_order,
  })
  dialogVisible.value = true
}

const submitForm = async () => {
  await formRef.value?.validate().catch(() => { throw new Error('validation') })
  submitting.value = true
  try {
    if (editingProvider.value) {
      const req: Record<string, unknown> = {
        name: form.name,
        type: form.type,
        api_url: form.api_url,
        models: form.models,
        sort_order: form.sort_order,
      }
      if (form.api_key) req.api_key = form.api_key
      await updateProvider(editingProvider.value.id, req as Parameters<typeof updateProvider>[1])
      ElMessage.success('Provider 已更新')
    } else {
      await createProvider({
        name: form.name,
        type: form.type,
        api_url: form.api_url,
        api_key: form.api_key,
        models: form.models,
        sort_order: form.sort_order,
      })
      ElMessage.success('Provider 已创建')
    }
    dialogVisible.value = false
    await loadProviders()
  } catch (err) {
    if (String(err) !== 'Error: validation') {
      ElMessage.error(`操作失败：${err}`)
    }
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async (provider: Provider) => {
  await ElMessageBox.confirm(
    `确认删除 Provider "${provider.name}"？此操作不可撤销。`,
    '删除确认',
    { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
  )
  try {
    await deleteProvider(provider.id)
    ElMessage.success('Provider 已删除')
    await loadProviders()
  } catch (err) {
    ElMessage.error(`删除失败：${err}`)
  }
}

// ---- 测试 & 拉取模型 ----
const testProvider = async (provider: Provider) => {
  testingId.value = provider.id
  testResult.value = null
  try {
    const result = await apiTestProvider(provider.id)
    testResult.value = result
    if (result?.success) {
      ElMessage.success(`连通性测试成功（${result.latency_ms}ms）`)
    } else {
      ElMessage.error(`连通性测试失败：${result?.message}`)
    }
  } catch (err) {
    ElMessage.error(`测试请求失败：${err}`)
  } finally {
    testingId.value = null
  }
}

const fetchModels = async (provider: Provider) => {
  fetchingId.value = provider.id
  try {
    const models = await fetchProviderModels(provider.id)
    if (models && models.length > 0) {
      // 更新 models 字段
      await updateProvider(provider.id, { models: JSON.stringify(models) })
      ElMessage.success(`已拉取 ${models.length} 个模型`)
      await loadProviders()
    } else {
      ElMessage.warning('未获取到任何模型')
    }
  } catch (err) {
    ElMessage.error(`拉取模型失败：${err}`)
  } finally {
    fetchingId.value = null
  }
}

onMounted(loadProviders)
</script>

<style scoped>
.form-hint {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}
</style>
