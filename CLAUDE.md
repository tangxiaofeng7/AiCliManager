# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概览

AiCliManager 是一个基于 **Wails v2 + Go + Vue 3 + SQLite** 的桌面应用，用来统一管理并启动多个 AI CLI 工具（目前内置 Claude Code、Codex、OpenCode）。

应用的核心职责不是“自己实现聊天”，而是：
- 管理 Provider / Profile / Proxy / MCP / Skill 配置
- 在启动 CLI 工具前把配置同步写入对应工具的本地配置文件
- 通过不同终端启动 CLI 子进程
- 同时维护两套“会话”能力：
  1. `sessions` 表：AiCliManager 自己记录的启动历史
  2. `CliSessionService`：直接读取 CLI 工具本地会话文件（当前主要支持 Claude 的 `~/.claude/projects`）

## 常用命令

### 安装依赖
```bash
go mod tidy
cd frontend && npm install
```

### 开发
```bash
wails dev
```

### 前端单独开发
```bash
cd frontend && npm run dev
```

### 构建
```bash
cd frontend && npm run build
wails build
```

### 测试
```bash
go test ./...
```

### 运行单个测试
```bash
go test ./internal/service -run TestName
```

### 当前仓库现状
- 没有配置独立的 lint 命令；不要在 CLAUDE.md 中假设存在 `npm run lint` 或 golangci-lint。
- `frontend/package.json` 里只有 `dev` / `build` / `preview`。
- 当前 `go test ./...` 可以执行，但仓库里基本还没有测试文件。

## 高层架构

### 1. Wails 入口与依赖装配
- `main.go` 初始化 SQLite，构造 `app.App`，并通过 Wails 绑定给前端。
- `internal/app/app.go` 是后端依赖注入中心：集中创建各个 service，并在 `App` 上暴露给前端调用的方法。
- 前端实际调用的是 Wails 生成的绑定，但项目又包了一层 `frontend/src/api/index.ts`，避免页面直接依赖生成代码。

### 2. 后端分层

#### `internal/app`
Wails 绑定层。这里的方法主要做：
- 参数校验
- 将前端请求结构转换为 service 层结构
- 把 service 返回值直接暴露给前端

不要把核心业务逻辑堆在这里；逻辑通常应落在 `internal/service`。

#### `internal/service`
主要业务层，关键点如下：
- `LauncherService`：启动主流程编排者，负责查询配置、解密敏感字段、触发配置同步、构建终端命令、注入环境变量、记录会话。
- `SyncService`：根据工具类型把 Provider / Profile / MCP / Skill 写入不同 CLI 的配置格式。
- `TerminalService`：按平台构造终端启动命令，Windows/macOS/Linux 分开处理。
- `ProviderService` / `ProxyService`：除了 CRUD，还负责解密敏感字段供启动阶段使用。
- `SessionService`：操作 SQLite 里的启动历史。
- `CliSessionService`：直接解析 AI CLI 本地 JSONL 会话文件，用于“真实对话历史”展示。

#### `internal/cli`
不同 CLI 工具配置文件的读写适配层。`SyncService` 会按 `cliToolKey` 选用对应 writer。

#### `internal/db`
- `db.go` 直接执行建表 SQL 和初始数据 seed，不使用 migration 框架。
- 数据库文件位于用户目录：`~/.aiclimgr/data.db`
- 初始化时会写入默认 `cli_tools` 记录和部分全局设置。

### 3. 启动链路

CLI 启动的关键调用链：
1. 前端 Dashboard 调用 `launchCliTool(...)`
2. `internal/app/launcher.go` 做参数透传
3. `internal/service/launcher_service.go` 执行主流程：
   - 查询 `cli_tools`
   - 必要时自动检测可执行文件
   - 查询 Profile / Provider / Proxy / MCP / Skill
   - 解密 API Key / 代理密码
   - 调用 `SyncService.SyncConfig(...)`
   - 通过 `TerminalService.BuildCmd(...)` 构建终端命令
   - 注入代理环境变量
   - 启动子进程并写入 `sessions` 表
   - goroutine 等待进程退出并更新状态

理解启动功能时，优先读：
- `internal/service/launcher_service.go`
- `internal/service/sync_service.go`
- `internal/service/terminal_service.go`
- `internal/cli/*.go`

### 4. 配置同步模型

`SyncService` 是项目最核心的业务之一。

它接收聚合后的 `SyncRequest`，再按不同工具写入不同字段：
- Claude: `apiKey` / `baseUrl` / `model` / `systemPrompt` / `mcpServers` / `customCommands`
- Codex / OpenCode: 字段名与 Claude 类似但不完全一致

这意味着：
- 新增一个“可启动的 CLI 工具”时，通常不只是加一条数据库记录
- 还要同时补齐：检测逻辑、默认配置路径、配置读写 adapter、同步字段映射、终端启动兼容性

### 5. 两套会话系统不要混淆

#### 启动历史（SQLite）
- 表：`sessions`
- 服务：`SessionService`
- 用途：记录“某次启动了哪个工具，用了哪个 Profile/Proxy，在哪个终端和目录启动”

#### CLI 真实对话历史（本地 JSONL）
- 服务：`CliSessionService`
- 当前主要支持 Claude 本地目录：`~/.claude/projects`
- 用途：在 UI 中浏览 AI CLI 实际对话内容

`frontend/src/views/Sessions.vue` 使用的是第二套能力，不是数据库里的 `sessions` 表。

## 前端结构

### API 封装层
`frontend/src/api/index.ts` 很关键：
- 统一封装 Wails Go 方法调用
- 在非 Wails 环境下自动降级，返回空数组或 `null`
- 所以前端页面在浏览器单独跑时不会因为缺少 Wails runtime 直接崩掉

如果你改了后端绑定方法，通常也要同步更新这里的封装。

### 路由与页面
- `frontend/src/router/index.ts` 使用 `createWebHashHistory()`，这是为了适配 Wails 打包后的本地静态资源场景。
- 页面基本按功能模块分：Dashboard / Providers / Profiles / Proxy / MCP / Skills / Sessions / Settings。

### Dashboard 页
`frontend/src/views/Dashboard.vue` 是启动主入口页面，负责：
- 加载工具列表
- 逐个触发安装检测
- 加载每个工具当前激活的 Profile / Proxy
- 弹出启动对话框并调用后端启动

这个页面体现了产品主流程，改启动相关功能时优先看它。

## 重要实现细节

### 数据库与迁移方式
- 表结构定义直接写在 `internal/db/db.go`。
- 如果需要改 schema，直接修改建表 SQL，并补上必要的手动兼容逻辑；不要引入 AutoMigrate 思路。

### 敏感信息处理
- Provider 的 API Key、Proxy 密码是加密存储的。
- 启动阶段由 service 层解密后只在内存中使用，不应回传前端。

### 环境变量过滤
- `LauncherService` 会过滤 `CLAUDECODE`，避免 Claude Code 因嵌套会话检测导致启动失败。
- 这是启动链路里的特定兼容逻辑，修改环境变量传递时不要误删。

### 终端能力是平台相关的
`TerminalService` 针对不同 OS 构建不同命令：
- Windows: `wt` / `cmd` / `powershell` / `wsl`
- macOS: `Terminal.app` / `iTerm2`
- Linux: `gnome-terminal` / `tmux` 等

调试启动失败时，先区分是：
1. CLI 本身不可执行
2. 终端包装命令有问题
3. 配置同步写坏了目标 CLI 配置文件

## 修改代码时的仓库约定

- 代码注释使用中文，标识符保持英文。
- 不要手改 `frontend/wailsjs/` 下的生成文件。
- 新增前端调用前，先确认后端 `App` 上已经存在对应绑定方法。
- 涉及启动、同步、加密、会话解析时，优先保持现有分层，不要把逻辑挪到 view 或 app binding 层。

## 未来维护时最值得先读的文件

- `main.go`
- `internal/app/app.go`
- `internal/app/launcher.go`
- `internal/service/launcher_service.go`
- `internal/service/sync_service.go`
- `internal/service/terminal_service.go`
- `internal/service/cli_session_service.go`
- `internal/db/db.go`
- `frontend/src/api/index.ts`
- `frontend/src/views/Dashboard.vue`
- `frontend/src/views/Sessions.vue`
