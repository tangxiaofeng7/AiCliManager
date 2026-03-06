<template>
  <!-- 底部状态栏 -->
  <div class="statusbar">
    <span class="status-item">AiCliManager</span>
    <span class="status-divider">|</span>
    <span class="status-item">{{ appVersion }}</span>
    <span class="status-divider">|</span>
    <span class="status-item">Wails v2 + Go 1.24 + Vue 3</span>
    <div class="status-spacer" />
    <span class="status-item">{{ currentTime }}</span>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import packageJson from '../../../package.json'

// 当前时间显示
const currentTime = ref('')
const appVersion = `v${packageJson.version}`

const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour12: false })
}

let timer: ReturnType<typeof setInterval>

onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<style scoped>
.statusbar {
  display: flex;
  align-items: center;
  height: var(--statusbar-height);
  background-color: var(--accent-color);
  padding: 0 12px;
  flex-shrink: 0;
  gap: 6px;
}

.status-item {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.9);
  white-space: nowrap;
}

.status-divider {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
}

.status-spacer {
  flex: 1;
}
</style>
