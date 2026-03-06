<template>
  <div class="page-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <div>
        <h2>Skills / Commands</h2>
        <div class="subtitle">管理可跨 CLI 工具复用的提示词技能和指令模板</div>
      </div>
      <el-button type="primary" :icon="Plus" @click="openCreateDialog">新建 Skill</el-button>
    </div>

    <!-- 分类筛选 -->
    <div class="filter-bar">
      <el-radio-group v-model="filterCategory" size="small">
        <el-radio-button value="">全部</el-radio-button>
        <el-radio-button v-for="cat in categories" :key="cat" :value="cat">
          {{ cat }}
        </el-radio-button>
      </el-radio-group>
    </div>

    <!-- Skill 卡片列表 -->
    <div class="skills-grid" v-loading="loading">
      <div
        v-for="skill in filteredSkills"
        :key="skill.id"
        class="skill-card"
        :class="{ builtin: skill.is_builtin }"
      >
        <!-- 头部：名称 + 内置标签 -->
        <div class="skill-header">
          <div class="skill-name">{{ skill.name }}</div>
          <div class="skill-badges">
            <el-tag v-if="skill.is_builtin" type="warning" size="small" effect="dark">内置</el-tag>
            <el-tag type="info" size="small">{{ skill.category }}</el-tag>
          </div>
        </div>

        <!-- 触发词 -->
        <div class="skill-trigger" v-if="skill.trigger">
          <el-tag type="primary" size="small" class="mono">{{ skill.trigger }}</el-tag>
        </div>

        <!-- 内容预览 -->
        <div class="skill-content">
          {{ truncate(skill.content, 120) }}
        </div>

        <!-- 操作 -->
        <div class="skill-actions">
          <el-button size="small" :icon="View" @click="previewSkill(skill)">预览</el-button>
          <el-button size="small" :icon="Edit" @click="openEditDialog(skill)" :disabled="!!skill.is_builtin">
            编辑
          </el-button>
          <el-button
            size="small"
            type="danger"
            :icon="Delete"
            :disabled="!!skill.is_builtin"
            @click="confirmDelete(skill)"
          >
            删除
          </el-button>
        </div>
      </div>

      <el-empty v-if="!loading && filteredSkills.length === 0" description="暂无 Skill" />
    </div>

    <!-- 新建 / 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingSkill ? '编辑 Skill' : '新建 Skill'"
      width="620px"
      destroy-on-close
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="90px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="如：代码审查" clearable />
        </el-form-item>

        <el-form-item label="分类">
          <el-select v-model="form.category" allow-create filterable placeholder="选择或新建分类" style="width: 100%">
            <el-option label="general（通用）" value="general" />
            <el-option label="code（代码）" value="code" />
            <el-option label="translate（翻译）" value="translate" />
            <el-option label="docs（文档）" value="docs" />
            <el-option label="review（审查）" value="review" />
          </el-select>
        </el-form-item>

        <el-form-item label="触发词">
          <el-input v-model="form.trigger" placeholder='如：/review（以 / 开头的斜杠命令）' clearable />
          <div class="form-hint">触发词通常以 / 开头，作为快捷调用命令</div>
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="8"
            placeholder="提示词内容，支持 {{变量名}} 占位符，如：请对以下代码进行 {{language}} 风格的审查：\n\n{{code}}"
          />
          <div class="form-hint">使用 {'{'}{'{'}变量名{'}'}{'}'}  格式定义占位符，启动时会弹窗填写变量值</div>
        </el-form-item>

        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" :max="9999" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ editingSkill ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 预览对话框 -->
    <el-dialog v-model="previewVisible" :title="`预览：${previewingSkill?.name}`" width="580px">
      <div class="preview-content">
        <div class="preview-meta">
          <el-tag type="info" size="small">{{ previewingSkill?.category }}</el-tag>
          <el-tag v-if="previewingSkill?.trigger" type="primary" size="small" class="mono">
            {{ previewingSkill.trigger }}
          </el-tag>
        </div>
        <pre class="preview-text">{{ previewingSkill?.content }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus, Edit, Delete, View } from '@element-plus/icons-vue'
import type { Skill } from '../types'
import { getSkills, createSkill, updateSkill, deleteSkill } from '../api'

// ---- 状态 ----
const loading = ref(false)
const submitting = ref(false)
const skills = ref<Skill[]>([])
const filterCategory = ref('')
const dialogVisible = ref(false)
const previewVisible = ref(false)
const editingSkill = ref<Skill | null>(null)
const previewingSkill = ref<Skill | null>(null)
const formRef = ref<FormInstance>()

const form = reactive({
  name: '',
  category: 'general',
  trigger: '',
  content: '',
  sort_order: 0,
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入 Skill 名称', trigger: 'blur' }],
  content: [{ required: true, message: '请输入提示词内容', trigger: 'blur' }],
}

// ---- 计算 ----
const categories = computed<string[]>(() => {
  const cats = new Set(skills.value.map(s => s.category).filter(Boolean))
  return Array.from(cats)
})

const filteredSkills = computed(() => {
  if (!filterCategory.value) return skills.value
  return skills.value.filter(s => s.category === filterCategory.value)
})

// ---- 辅助 ----
const truncate = (str: string, len: number) =>
  str.length > len ? str.slice(0, len) + '…' : str

// ---- 数据加载 ----
const loadSkills = async () => {
  loading.value = true
  try {
    skills.value = (await getSkills().catch(() => [])) ?? []
  } finally {
    loading.value = false
  }
}

// ---- 对话框 ----
const openCreateDialog = () => {
  editingSkill.value = null
  Object.assign(form, { name: '', category: 'general', trigger: '', content: '', sort_order: 0 })
  dialogVisible.value = true
}

const openEditDialog = (skill: Skill) => {
  if (skill.is_builtin) return
  editingSkill.value = skill
  Object.assign(form, {
    name: skill.name,
    category: skill.category,
    trigger: skill.trigger ?? '',
    content: skill.content,
    sort_order: skill.sort_order,
  })
  dialogVisible.value = true
}

const previewSkill = (skill: Skill) => {
  previewingSkill.value = skill
  previewVisible.value = true
}

const submitForm = async () => {
  await formRef.value?.validate().catch(() => { throw new Error('validation') })
  submitting.value = true
  try {
    if (editingSkill.value) {
      await updateSkill(editingSkill.value.id, { ...form })
      ElMessage.success('Skill 已更新')
    } else {
      await createSkill({ ...form })
      ElMessage.success('Skill 已创建')
    }
    dialogVisible.value = false
    await loadSkills()
  } catch (err) {
    if (String(err) !== 'Error: validation') ElMessage.error(`操作失败：${err}`)
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async (skill: Skill) => {
  if (skill.is_builtin) return
  await ElMessageBox.confirm(
    `确认删除 Skill "${skill.name}"？`,
    '删除确认',
    { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
  )
  try {
    await deleteSkill(skill.id)
    ElMessage.success('Skill 已删除')
    await loadSkills()
  } catch (err) {
    ElMessage.error(`删除失败：${err}`)
  }
}

onMounted(loadSkills)
</script>

<style scoped>
/* 分类筛选 */
.filter-bar {
  margin-bottom: 16px;
}

/* Skill 卡片网格 */
.skills-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 14px;
  overflow-y: auto;
}

.skill-card {
  background-color: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  transition: border-color 0.2s;
}

.skill-card:hover {
  border-color: var(--accent-color);
}

.skill-card.builtin {
  border-left: 3px solid var(--warning-color);
}

.skill-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
}

.skill-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.skill-badges {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.skill-trigger {
  display: flex;
  align-items: center;
}

.skill-content {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
  flex: 1;
}

.skill-actions {
  display: flex;
  gap: 6px;
  margin-top: auto;
}

/* 预览 */
.preview-meta {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.preview-text {
  background-color: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 14px;
  font-family: 'Cascadia Code', 'JetBrains Mono', monospace;
  font-size: 13px;
  color: var(--text-primary);
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 400px;
  overflow-y: auto;
}

.form-hint {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}
</style>
