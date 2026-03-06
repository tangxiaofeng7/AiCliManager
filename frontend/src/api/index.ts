/**
 * API 封装层：统一处理 Wails Go 方法调用
 * 在浏览器开发模式下做有限降级；Wails 环境中的绑定缺失或调用异常必须显式抛错。
 */

import type {
  CliTool,
  Provider,
  Profile,
  Proxy,
  McpServer,
  Skill,
  Session,
  CliSession,
  CliSessionMessage,
  CliSessionProject,
  ActiveConfig,
  TerminalInfo,
  TestResult,
  LaunchRequest,
  GetSessionsRequest,
  GetCliSessionsRequest,
  CreateProviderRequest,
  UpdateProviderRequest,
  CreateProfileRequest,
  UpdateProfileRequest,
  CreateProxyRequest,
  UpdateProxyRequest,
  CreateMcpServerRequest,
  UpdateMcpServerRequest,
  CreateSkillRequest,
  UpdateSkillRequest,
  SetActiveConfigRequest,
  DetectResult,
} from '../types'

const isWailsEnv = (): boolean => {
  return typeof window !== 'undefined' && !!(window as unknown as Record<string, unknown>)['go']
}

let wailsModule: Record<string, (...args: unknown[]) => Promise<unknown>> | null = null

const getWailsModule = async () => {
  if (wailsModule) return wailsModule
  try {
    const mod = await import('../../wailsjs/go/app/App')
    wailsModule = mod as unknown as Record<string, (...args: unknown[]) => Promise<unknown>>
    return wailsModule
  } catch (error) {
    if (isWailsEnv()) {
      throw new Error(`Wails 绑定加载失败: ${String(error)}`)
    }
    return null
  }
}

const devFallbackFactories: Record<string, () => unknown> = {
  GetCliTools: () => [],
  DetectCliTool: () => null,
  GetCliToolActiveConfig: () => ({ profile_id: null, proxy_id: null }),
  ListAvailableTerminals: () => [],
  GetProviders: () => [],
  CreateProvider: () => null,
  GetProfiles: () => [],
  CreateProfile: () => null,
  GetProxies: () => [],
  CreateProxy: () => null,
  GetMcpServers: () => [],
  CreateMcpServer: () => null,
  GetCliToolMcpServers: () => [],
  GetSkills: () => [],
  CreateSkill: () => null,
  GetSessions: () => [],
  GetCliSessions: () => [],
  GetCliSessionMessages: () => [],
  GetCliSessionProjects: () => [],
  GetSettings: () => ({}),
  TestProvider: () => null,
  FetchProviderModels: () => [],
}

const callGo = async <T>(method: string, ...args: unknown[]): Promise<T> => {
  const mod = await getWailsModule()
  if (!mod || typeof mod[method] !== 'function') {
    if (!isWailsEnv()) {
      const fallback = devFallbackFactories[method]
      if (fallback) return fallback() as T
      return undefined as T
    }
    throw new Error(`Wails 方法不可用: ${method}`)
  }

  try {
    return await mod[method](...args) as T
  } catch (error) {
    throw error instanceof Error ? error : new Error(String(error))
  }
}

export const getCliTools = (): Promise<CliTool[]> =>
  callGo<CliTool[]>('GetCliTools')

export const detectCliTool = (key: string): Promise<DetectResult | null> =>
  callGo<DetectResult | null>('DetectCliTool', key)

export const launchCliTool = (req: LaunchRequest): Promise<void> =>
  callGo<void>('LaunchCliTool', req)

export const getCliToolActiveConfig = (key: string): Promise<ActiveConfig | null> =>
  callGo<ActiveConfig | null>('GetCliToolActiveConfig', key)

export const setCliToolActiveConfig = (req: SetActiveConfigRequest): Promise<void> =>
  callGo<void>('SetCliToolActiveConfig', req)

export const listAvailableTerminals = (): Promise<TerminalInfo[]> =>
  callGo<TerminalInfo[]>('ListAvailableTerminals')

export const getProviders = (): Promise<Provider[]> =>
  callGo<Provider[]>('GetProviders')

export const createProvider = (req: CreateProviderRequest): Promise<Provider | null> =>
  callGo<Provider | null>('CreateProvider', req)

export const updateProvider = (id: number, req: UpdateProviderRequest): Promise<void> =>
  callGo<void>('UpdateProvider', id, req)

export const deleteProvider = (id: number): Promise<void> =>
  callGo<void>('DeleteProvider', id)

export const testProvider = (id: number): Promise<TestResult | null> =>
  callGo<TestResult | null>('TestProvider', id)

export const fetchProviderModels = (id: number): Promise<string[]> =>
  callGo<string[]>('FetchProviderModels', id)

export const getProfiles = (): Promise<Profile[]> =>
  callGo<Profile[]>('GetProfiles')

export const createProfile = (req: CreateProfileRequest): Promise<Profile | null> =>
  callGo<Profile | null>('CreateProfile', req)

export const updateProfile = (id: number, req: UpdateProfileRequest): Promise<void> =>
  callGo<void>('UpdateProfile', id, req)

export const deleteProfile = (id: number): Promise<void> =>
  callGo<void>('DeleteProfile', id)

export const getProxies = (): Promise<Proxy[]> =>
  callGo<Proxy[]>('GetProxies')

export const createProxy = (req: CreateProxyRequest): Promise<Proxy | null> =>
  callGo<Proxy | null>('CreateProxy', req)

export const updateProxy = (id: number, req: UpdateProxyRequest): Promise<void> =>
  callGo<void>('UpdateProxy', id, req)

export const deleteProxy = (id: number): Promise<void> =>
  callGo<void>('DeleteProxy', id)

export const setGlobalProxy = (id: number): Promise<void> =>
  callGo<void>('SetGlobalProxy', id)

export const clearGlobalProxy = (): Promise<void> =>
  callGo<void>('ClearGlobalProxy')

export const getMcpServers = (): Promise<McpServer[]> =>
  callGo<McpServer[]>('GetMcpServers')

export const createMcpServer = (req: CreateMcpServerRequest): Promise<McpServer | null> =>
  callGo<McpServer | null>('CreateMcpServer', req)

export const updateMcpServer = (id: number, req: UpdateMcpServerRequest): Promise<void> =>
  callGo<void>('UpdateMcpServer', id, req)

export const deleteMcpServer = (id: number): Promise<void> =>
  callGo<void>('DeleteMcpServer', id)

export const getCliToolMcpServers = (cliToolKey: string): Promise<McpServer[]> =>
  callGo<McpServer[]>('GetCliToolMcpServers', cliToolKey)

export const setCliToolMcpServers = (cliToolKey: string, ids: number[]): Promise<void> =>
  callGo<void>('SetCliToolMcpServers', cliToolKey, ids)

export const getSkills = (): Promise<Skill[]> =>
  callGo<Skill[]>('GetSkills')

export const createSkill = (req: CreateSkillRequest): Promise<Skill | null> =>
  callGo<Skill | null>('CreateSkill', req)

export const updateSkill = (id: number, req: UpdateSkillRequest): Promise<void> =>
  callGo<void>('UpdateSkill', id, req)

export const deleteSkill = (id: number): Promise<void> =>
  callGo<void>('DeleteSkill', id)

export const getSessions = (req: GetSessionsRequest): Promise<Session[]> =>
  callGo<Session[]>('GetSessions', req)

export const deleteSession = (id: number): Promise<void> =>
  callGo<void>('DeleteSession', id)

export const clearSessions = (cliToolKey: string): Promise<void> =>
  callGo<void>('ClearSessions', cliToolKey)

export const relaunchSession = (sessionId: number): Promise<void> =>
  callGo<void>('RelaunchSession', sessionId)

export const getCliSessions = (req: GetCliSessionsRequest): Promise<CliSession[]> =>
  callGo<CliSession[]>('GetCliSessions', req)

export const getCliSessionMessages = (cliToolKey: string, sessionId: string): Promise<CliSessionMessage[]> =>
  callGo<CliSessionMessage[]>('GetCliSessionMessages', cliToolKey, sessionId)

export const getCliSessionProjects = (cliToolKey: string): Promise<CliSessionProject[]> =>
  callGo<CliSessionProject[]>('GetCliSessionProjects', cliToolKey)

export const getSettings = (): Promise<Record<string, string>> =>
  callGo<Record<string, string>>('GetSettings')

export const setSetting = (key: string, value: string): Promise<void> =>
  callGo<void>('SetSetting', key, value)

export const saveSettings = (settings: Record<string, string>): Promise<void> =>
  callGo<void>('SaveSettings', settings)
