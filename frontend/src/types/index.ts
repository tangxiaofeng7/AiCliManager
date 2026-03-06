// 类型定义文件：所有前端使用的 TypeScript 接口与类型

export interface CliTool {
  id: number
  name: string
  key: string
  executable: string
  config_path: string
  preferred_terminal: string
  is_installed: number
  is_enabled: number
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Provider {
  id: number
  name: string
  type: 'anthropic' | 'openai' | 'custom'
  api_url: string
  api_key: string
  models: string
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Profile {
  id: number
  name: string
  provider_id: number
  model: string
  system_prompt: string
  temperature: number
  max_tokens: number
  extra_config: string
  created_at: string
  updated_at: string
}

export interface Proxy {
  id: number
  name: string
  type: 'http' | 'https' | 'socks5'
  host: string
  port: number
  username: string
  password: string
  no_proxy: string
  is_active: number
  created_at: string
  updated_at: string
}

export interface McpServer {
  id: number
  name: string
  type: 'stdio' | 'sse' | 'http'
  command: string
  args: string
  env: string
  url: string
  description: string
  is_enabled: number
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Skill {
  id: number
  name: string
  category: string
  trigger: string
  content: string
  variables: string
  is_builtin: number
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Session {
  id: number
  cli_tool_id: number
  profile_id: number | null
  proxy_id: number | null
  terminal: string
  working_dir: string
  extra_args: string
  status: 'running' | 'exited' | 'error'
  pid: number
  started_at: string
  ended_at: string | null
}

export interface CliSession {
  session_id: string
  cli_tool_key: string
  project: string
  project_dir: string
  slug: string
  first_message: string
  message_count: number
  user_count: number
  assistant_count: number
  model: string
  started_at: string
  last_active_at: string
}

export interface CliSessionMessage {
  type: 'user' | 'assistant' | 'system'
  content: string
  timestamp: string
  model?: string
  tokens_in?: number
  tokens_out?: number
  uuid: string
}

export interface CliSessionProject {
  dir_name: string
  path: string
  session_count: number
}

export interface GetCliSessionsRequest {
  cli_tool_key?: string
  project?: string
  limit?: number
}

export interface ActiveConfig {
  profile_id: number | null
  proxy_id: number | null
}

export interface TerminalInfo {
  id: string
  name: string
  is_available: boolean
}

export interface TestResult {
  success: boolean
  message: string
  latency_ms: number
}

export interface LaunchRequest {
  cli_tool_key: string
  profile_id: number
  proxy_id: number | null
  mcp_server_ids: number[]
  skill_ids: number[]
  skill_vars: Record<string, string>
  extra_args: string[]
  terminal: string
  working_dir: string
}

export interface GetSessionsRequest {
  cli_tool_key?: string
  page?: number
  page_size?: number
}

export interface CreateProviderRequest {
  name: string
  type: 'anthropic' | 'openai' | 'custom'
  api_url: string
  api_key: string
  models?: string
  sort_order?: number
}

export interface UpdateProviderRequest {
  name?: string
  type?: 'anthropic' | 'openai' | 'custom'
  api_url?: string
  api_key?: string
  models?: string
  sort_order?: number
}

export interface CreateProfileRequest {
  name: string
  provider_id: number
  model: string
  system_prompt?: string
  temperature?: number
  max_tokens?: number
  extra_config?: string
}

export interface UpdateProfileRequest {
  name?: string
  provider_id?: number
  model?: string
  system_prompt?: string
  temperature?: number
  max_tokens?: number
  extra_config?: string
}

export interface CreateProxyRequest {
  name: string
  type: 'http' | 'https' | 'socks5'
  host: string
  port: number
  username?: string
  password?: string
  no_proxy?: string
}

export interface UpdateProxyRequest {
  name?: string
  type?: 'http' | 'https' | 'socks5'
  host?: string
  port?: number
  username?: string
  password?: string
  no_proxy?: string
}

export interface CreateMcpServerRequest {
  name: string
  type: 'stdio' | 'sse' | 'http'
  command?: string
  args?: string
  env?: string
  url?: string
  description?: string
  is_enabled?: number
  sort_order?: number
}

export interface UpdateMcpServerRequest {
  name?: string
  type?: 'stdio' | 'sse' | 'http'
  command?: string
  args?: string
  env?: string
  url?: string
  description?: string
  is_enabled?: number
  sort_order?: number
}

export interface CreateSkillRequest {
  name: string
  category?: string
  trigger?: string
  content: string
  variables?: string
  sort_order?: number
}

export interface UpdateSkillRequest {
  name?: string
  category?: string
  trigger?: string
  content?: string
  variables?: string
  sort_order?: number
}

export interface SetActiveConfigRequest {
  cli_tool_key: string
  profile_id: number | null
  proxy_id: number | null
}

export interface DetectResult {
  key: string
  is_installed: boolean
  executable: string
}
