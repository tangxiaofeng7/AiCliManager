<template>
  <!-- 无边框窗口自定义标题栏 -->
  <div class="titlebar">
    <!-- 左侧：应用图标 + 名称 -->
    <div class="titlebar-left">
      <div class="app-icon">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="2" y="3" width="20" height="14" rx="2" stroke="#409eff" stroke-width="2"/>
          <path d="M8 21h8M12 17v4" stroke="#409eff" stroke-width="2" stroke-linecap="round"/>
          <path d="M6 8l3 3-3 3M11 14h6" stroke="#67c23a" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </div>
      <span class="app-name">AiCliManager</span>
    </div>

    <!-- 中间：可拖拽区域 -->
    <div class="titlebar-drag" style="-webkit-app-region: drag; flex: 1; height: 100%;" />

    <!-- 右侧：窗口控制按钮 -->
    <div class="titlebar-controls">
      <button class="ctrl-btn minimize" title="最小化" @click="minimize">
        <svg width="10" height="1" viewBox="0 0 10 1">
          <line x1="0" y1="0.5" x2="10" y2="0.5" stroke="currentColor" stroke-width="1.2"/>
        </svg>
      </button>
      <button class="ctrl-btn maximize" title="最大化 / 还原" @click="toggleMaximize">
        <svg width="10" height="10" viewBox="0 0 10 10">
          <rect x="0.5" y="0.5" width="9" height="9" rx="1" stroke="currentColor" stroke-width="1.2" fill="none"/>
        </svg>
      </button>
      <button class="ctrl-btn close" title="关闭" @click="quit">
        <svg width="10" height="10" viewBox="0 0 10 10">
          <line x1="0" y1="0" x2="10" y2="10" stroke="currentColor" stroke-width="1.2"/>
          <line x1="10" y1="0" x2="0" y2="10" stroke="currentColor" stroke-width="1.2"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
// 调用 Wails 运行时 API 控制窗口
const minimize = () => {
  try {
    // @ts-ignore
    window.runtime?.WindowMinimise()
  } catch {
    console.warn('[TitleBar] WindowMinimise 不可用')
  }
}

const toggleMaximize = () => {
  try {
    // @ts-ignore
    window.runtime?.WindowToggleMaximise()
  } catch {
    console.warn('[TitleBar] WindowToggleMaximise 不可用')
  }
}

const quit = () => {
  try {
    // @ts-ignore
    window.runtime?.Quit()
  } catch {
    console.warn('[TitleBar] Quit 不可用')
  }
}
</script>

<style scoped>
.titlebar {
  display: flex;
  align-items: center;
  height: var(--titlebar-height);
  background-color: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
  /* 标题栏整体可拖拽，子元素覆盖此属性 */
  -webkit-app-region: drag;
  user-select: none;
}

.titlebar-left {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 12px;
  -webkit-app-region: no-drag;
  pointer-events: none; /* 左侧仅展示，不拦截拖拽 */
}

.app-icon {
  display: flex;
  align-items: center;
}

.app-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: 0.3px;
}

/* 窗口控制按钮组 */
.titlebar-controls {
  display: flex;
  align-items: stretch;
  -webkit-app-region: no-drag;
}

.ctrl-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 46px;
  height: var(--titlebar-height);
  border: none;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  transition: background-color 0.15s, color 0.15s;
}

.ctrl-btn:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.ctrl-btn.close:hover {
  background-color: #e81123;
  color: #fff;
}

.ctrl-btn svg {
  pointer-events: none;
}
</style>
