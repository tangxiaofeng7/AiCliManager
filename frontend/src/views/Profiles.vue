<template>
  <div class="page-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <div>
        <h2>Profile 管理</h2>
        <div class="subtitle">基于 Provider 创建多套参数配置，各 CLI 工具按需选用</div>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreateDialog">新建 Profile</el-button>
    </div>

    <!-- Profile 表格 -->
    <el-table :data="profiles" v-loading="loading" border stripe style="width: 100%">
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column label="Provider" min-width="140">
        <template #default="{ row }">
          <span>{{ getProviderName(row.provider_id) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="model" label="模型" min-width="200">
        <template #default="{ row }">
          <span class="mono">{{ row.model }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Temperature" width="110">
        <template #default="{ row }">
          <span>{{ row.temperature }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Max Tokens" width="110">
        <template #default="{ row }">
          <span>{{ row.max_tokens }}</span>
        </template>
      </el-table-column>
      <el-table-column label="System Prompt" min-width="180">
        <template #default="{ row }">
          <el-tooltip :content="row.system_prompt" placement="top" v-if="row.system_prompt">
            <span class="text-muted">{{ truncate(row.system_prompt, 40) }}</span>
          </el-tooltip>
          <span class="text-muted" v-else>-</span>
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
      :title="editingProfile ? '编辑 Profile' : '新建 Profile'"
      width="580px"
      destroy-on-close
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="110px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="如：工作-Claude3.5" clearable />
        </el-form-item>

        <el-form-item label="Provider" prop="provider_id">
          <el-select
            v-model="form.provider_id"
            placeholder="选择 Provider"
            style="width: 100%"
            @change="onProviderChange"
          >
            <el-option
              v-for="p in providers"
              :key="p.id"
              :label="p.name"
              :value="p.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="模型" prop="model">
          <el-select
            v-model="form.model"
            placeholder="选择或手动输入模型名称"
            style="width: 100%"
            filterable
            allow-create
          >
            <el-option
              v-for="m in currentModels"
              :key="m"
              :label="m"
              :value="m"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="Temperature">
          <div style="display: flex; align-items: center; gap: 12px; width: 100%">
            <el-slider v-model="form.temperature" :min="0" :max="2" :step="0.1" style="flex: 1" />
            <el-input-number v-model="form.temperature" :min="0" :max="2" :step="0.1" :precision="1" style="width: 100px" />
          </div>
        </el-form-item>

        <el-form-item label="Max Tokens">
          <el-input-number v-model="form.max_tokens" :min="1" :max="200000" :step="1024" style="width: 180px" />
        </el-form-item>

        <el-form-item label="System Prompt">
          <el-input
            v-model="form.system_prompt"
            type="textarea"
            :rows="4"
            placeholder="可选：设置默认 System Prompt"
          />
        </el-form-item>

        <el-form-item label="扩展配置">
          <el-input
            v-model="form.extra_config"
            type="textarea"
            :rows="2"
            placeholder='可选：JSON 格式的扩展参数，如 {"stream": true}'
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ editingProfile ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import type { Profile, Provider } from '../types'
import {
  getProfiles, createProfile, updateProfile, deleteProfile,
  getProviders,
} from '../api'

// ---- 状态 ----
const loading = ref(false)
const submitting = ref(false)
const profiles = ref<Profile[]>([])
const providers = ref<Provider[]>([])
const dialogVisible = ref(false)
const editingProfile = ref<Profile | null>(null)
const formRef = ref<FormInstance>()

const form = reactive({
  name: '',
  provider_id: 0,
  model: '',
  system_prompt: '',
  temperature: 1.0,
  max_tokens: 8192,
  extra_config: '',
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入 Profile 名称', trigger: 'blur' }],
  provider_id: [{ required: true, message: '请选择 Provider', trigger: 'change' }],
  model: [{ required: true, message: '请选择或输入模型名称', trigger: 'blur' }],
}

// 当前选中 Provider 支持的模型列表
const currentModels = computed<string[]>(() => {
  const provider = providers.value.find(p => p.id === form.provider_id)
  if (!provider?.models) return []
  try {
    return JSON.parse(provider.models) as string[]
  } catch {
    return []
  }
})

// ---- 辅助 ----
const getProviderName = (id: number): string =>
  providers.value.find(p => p.id === id)?.name ?? `Provider #${id}`

const truncate = (str: string, len: number): string =>
  str.length > len ? str.slice(0, len) + '…' : str

// Provider 切换时清空模型选择
const onProviderChange = () => {
  form.model = ''
}

// ---- 数据加载 ----
const loadData = async () => {
  loading.value = true
  try {
    const [pf, pv] = await Promise.all([
      getProfiles().catch(() => []),
      getProviders().catch(() => []),
    ])
    profiles.value = pf ?? []
    providers.value = pv ?? []
  } finally {
    loading.value = false
  }
}

// ---- 对话框 ----
const openCreateDialog = () => {
  editingProfile.value = null
  Object.assign(form, { name: '', provider_id: 0, model: '', system_prompt: '', temperature: 1.0, max_tokens: 8192, extra_config: '' })
  dialogVisible.value = true
}

const openEditDialog = (profile: Profile) => {
  editingProfile.value = profile
  Object.assign(form, {
    name: profile.name,
    provider_id: profile.provider_id,
    model: profile.model,
    system_prompt: profile.system_prompt ?? '',
    temperature: profile.temperature,
    max_tokens: profile.max_tokens,
    extra_config: profile.extra_config ?? '',
  })
  dialogVisible.value = true
}

const submitForm = async () => {
  await formRef.value?.validate().catch(() => { throw new Error('validation') })
  submitting.value = true
  try {
    if (editingProfile.value) {
      await updateProfile(editingProfile.value.id, { ...form })
      ElMessage.success('Profile 已更新')
    } else {
      await createProfile({ ...form })
      ElMessage.success('Profile 已创建')
    }
    dialogVisible.value = false
    await loadData()
  } catch (err) {
    if (String(err) !== 'Error: validation') ElMessage.error(`操作失败：${err}`)
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async (profile: Profile) => {
  await ElMessageBox.confirm(
    `确认删除 Profile "${profile.name}"？`,
    '删除确认',
    { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
  )
  try {
    await deleteProfile(profile.id)
    ElMessage.success('Profile 已删除')
    await loadData()
  } catch (err) {
    ElMessage.error(`删除失败：${err}`)
  }
}

onMounted(loadData)
</script>
