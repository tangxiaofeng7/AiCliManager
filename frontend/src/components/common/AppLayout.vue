<template>
  <div class="app-layout">
    <!-- 自定义标题栏 -->
    <TitleBar />

    <!-- 主体区域：侧边栏 + 内容 -->
    <div class="app-body">
      <Sidebar v-model:collapsed="sidebarCollapsed" />
      <main class="app-main" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
        <RouterView />
      </main>
    </div>

    <!-- 底部状态栏 -->
    <StatusBar />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import TitleBar from './TitleBar.vue'
import Sidebar from './Sidebar.vue'
import StatusBar from './StatusBar.vue'

// 侧边栏折叠状态，默认展开
const sidebarCollapsed = ref(false)
</script>

<style scoped>
.app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  background-color: var(--bg-primary);
}

.app-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.app-main {
  flex: 1;
  overflow: hidden;
  transition: margin-left 0.2s ease;
}
</style>
