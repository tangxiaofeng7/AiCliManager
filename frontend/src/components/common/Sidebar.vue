<template>
  <aside class="sidebar" :class="{ collapsed: collapsed }">
    <!-- 折叠切换按钮 -->
    <div class="sidebar-toggle" @click="toggleCollapse" :title="collapsed ? '展开侧边栏' : '折叠侧边栏'">
      <el-icon><ArrowLeft v-if="!collapsed" /><ArrowRight v-else /></el-icon>
    </div>

    <!-- 导航菜单 -->
    <el-menu
      :default-active="currentPath"
      :collapse="collapsed"
      :collapse-transition="false"
      router
      class="sidebar-menu"
    >
      <el-menu-item index="/">
        <el-icon><Monitor /></el-icon>
        <template #title>启动面板</template>
      </el-menu-item>

      <el-divider class="menu-divider" />

      <div class="menu-group-label" v-if="!collapsed">配置管理</div>

      <el-menu-item index="/providers">
        <el-icon><Connection /></el-icon>
        <template #title>Providers</template>
      </el-menu-item>

      <el-menu-item index="/profiles">
        <el-icon><User /></el-icon>
        <template #title>Profiles</template>
      </el-menu-item>

      <el-menu-item index="/proxy">
        <el-icon><Share /></el-icon>
        <template #title>代理配置</template>
      </el-menu-item>

      <el-menu-item index="/mcp">
        <el-icon><Grid /></el-icon>
        <template #title>MCP Servers</template>
      </el-menu-item>

      <el-menu-item index="/skills">
        <el-icon><MagicStick /></el-icon>
        <template #title>Skills</template>
      </el-menu-item>

      <el-divider class="menu-divider" />

      <div class="menu-group-label" v-if="!collapsed">历史 & 设置</div>

      <el-menu-item index="/sessions">
        <el-icon><Clock /></el-icon>
        <template #title>会话历史</template>
      </el-menu-item>

      <el-menu-item index="/settings">
        <el-icon><Setting /></el-icon>
        <template #title>全局设置</template>
      </el-menu-item>
    </el-menu>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  Monitor,
  Connection,
  User,
  Share,
  Grid,
  MagicStick,
  Clock,
  Setting,
  ArrowLeft,
  ArrowRight,
} from '@element-plus/icons-vue'

// Props
interface Props {
  collapsed: boolean
}
const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  'update:collapsed': [value: boolean]
}>()

const route = useRoute()

// 当前激活的路由路径
const currentPath = computed(() => route.path)

// 切换折叠状态
const toggleCollapse = () => {
  emit('update:collapsed', !props.collapsed)
}
</script>

<style scoped>
.sidebar {
  width: var(--sidebar-width);
  min-width: var(--sidebar-width);
  background-color: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: width 0.2s ease, min-width 0.2s ease;
  position: relative;
  flex-shrink: 0;
}

.sidebar.collapsed {
  width: var(--sidebar-collapsed-width);
  min-width: var(--sidebar-collapsed-width);
}

/* 折叠切换按钮 */
.sidebar-toggle {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 6px 10px;
  cursor: pointer;
  color: var(--text-muted);
  transition: color 0.15s;
  border-bottom: 1px solid var(--border-color);
  min-height: 32px;
}

.sidebar-toggle:hover {
  color: var(--text-primary);
}

.collapsed .sidebar-toggle {
  justify-content: center;
}

/* 菜单 */
.sidebar-menu {
  flex: 1;
  border-right: none !important;
  overflow-y: auto;
  overflow-x: hidden;
}

/* 折叠时菜单宽度 */
.sidebar.collapsed .sidebar-menu {
  width: var(--sidebar-collapsed-width) !important;
}

/* 分组标签 */
.menu-group-label {
  padding: 8px 16px 4px;
  font-size: 11px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.8px;
  white-space: nowrap;
}

/* 分隔线 */
.menu-divider {
  margin: 6px 0 !important;
  border-color: var(--border-color) !important;
}

/* 菜单项图标对齐 */
:deep(.el-menu-item) {
  height: 40px;
  line-height: 40px;
  font-size: 13px;
  border-radius: 6px;
  margin: 2px 8px;
  width: calc(100% - 16px);
}

:deep(.el-menu--collapse .el-menu-item) {
  width: calc(100% - 8px);
  margin: 2px 4px;
  justify-content: center;
}
</style>
