import { createRouter, createWebHashHistory } from 'vue-router'

// 使用 Hash 模式路由，适配 Wails 打包后的静态资源加载
const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      component: () => import('../views/Dashboard.vue'),
      meta: { title: '启动面板' },
    },
    {
      path: '/providers',
      component: () => import('../views/Providers.vue'),
      meta: { title: 'Provider 管理' },
    },
    {
      path: '/profiles',
      component: () => import('../views/Profiles.vue'),
      meta: { title: 'Profile 管理' },
    },
    {
      path: '/proxy',
      component: () => import('../views/Proxy.vue'),
      meta: { title: '代理配置' },
    },
    {
      path: '/mcp',
      component: () => import('../views/MCP.vue'),
      meta: { title: 'MCP Servers' },
    },
    {
      path: '/skills',
      component: () => import('../views/Skills.vue'),
      meta: { title: 'Skills / Commands' },
    },
    {
      path: '/sessions',
      component: () => import('../views/Sessions.vue'),
      meta: { title: '会话历史' },
    },
    {
      path: '/settings',
      component: () => import('../views/Settings.vue'),
      meta: { title: '全局设置' },
    },
  ],
})

export default router
