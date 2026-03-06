# AiCliManager

AiCliManager 是一个基于 **Wails v2 + Go + Vue 3 + SQLite** 的桌面应用，用来统一管理并启动多个 AI CLI 工具。

[Release 下载](#release-产物) · [贡献指南](./CONTRIBUTING.md) · [安全策略](./SECURITY.md) · [行为准则](./CODE_OF_CONDUCT.md) · [许可证](./LICENSE)

当前内置支持：
- Claude Code
- Codex
- OpenCode

它的核心目标不是替代这些 CLI，而是作为它们的 **启动器 + 配置管理中心**：
- 统一管理 Provider、Profile、Proxy、MCP Server、Skill / Command
- 在启动前把配置同步写入对应 CLI 的本地配置文件
- 通过不同终端启动 CLI 子进程
- 记录启动历史，并读取部分 CLI 的本地对话历史

---

## 技术栈

### 后端
- Go 1.24
- Wails v2.11
- SQLite（`modernc.org/sqlite`）

### 前端
- Vue 3
- TypeScript
- Vite
- Element Plus
- Pinia
- Vue Router

---

## 当前功能

- CLI 工具检测与启动
- Provider 管理
  - API Key 加密存储
  - 连通性测试
  - 模型列表拉取
- Profile 管理
- Proxy 管理
- MCP Server 管理
- Skill / Command 管理
- 启动历史记录
- Claude 本地会话历史读取与展示

---

## 开发环境要求

- Go 1.24+
- Node.js 20+
- Wails CLI

安装 Wails CLI：

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

---

## 安装依赖

```bash
go mod tidy
cd frontend && npm install
```

---

## 本地开发

启动 Wails 开发模式：

```bash
wails dev
```

如果只想单独运行前端：

```bash
cd frontend && npm run dev
```

---

## 构建

前端构建：

```bash
cd frontend && npm run build
```

桌面应用构建：

```bash
wails build
```

---

## 测试

运行全部 Go 测试：

```bash
go test ./...
```

运行单个包内的单测：

```bash
go test ./internal/service -run TestName
```

> 当前仓库里基本还没有测试文件，因此 `go test ./...` 目前主要用于确认编译通过。

---

## 项目结构

```text
main.go                 Wails 入口
internal/app/           Wails 绑定层
internal/service/       核心业务逻辑
internal/cli/           各 CLI 配置文件读写适配
internal/db/            SQLite 初始化与数据模型
internal/crypto/        敏感信息加解密
frontend/src/views/     页面级 Vue 组件
frontend/src/api/       前端对 Wails 绑定的封装层
```

---

## 架构说明

### 1. 后端分层

#### `internal/app`
Wails 绑定层，负责：
- 对前端暴露可调用方法
- 做基础参数校验
- 把请求转给 service 层

#### `internal/service`
业务核心层，主要包含：
- `LauncherService`：启动流程编排
- `SyncService`：把配置写入不同 CLI 的配置文件
- `TerminalService`：按平台构建终端启动命令
- `SessionService`：管理 SQLite 中的启动历史
- `CliSessionService`：读取 CLI 本地 JSONL 会话文件

#### `internal/cli`
不同 CLI 工具配置格式的适配层。新增新的 CLI 工具时，通常需要在这里补配置读写实现。

#### `internal/db`
项目不使用 migration 框架，表结构直接在 `internal/db/db.go` 中创建和初始化。

---

### 2. CLI 启动流程

用户从 Dashboard 发起启动后，后端大致按以下顺序执行：

1. 查询 CLI 工具记录
2. 检测可执行文件路径
3. 查询 Profile / Provider / Proxy / MCP / Skill
4. 解密 API Key 和代理密码
5. 调用 `SyncService` 写入目标 CLI 配置文件
6. 调用 `TerminalService` 构造终端命令
7. 注入代理环境变量后启动子进程
8. 将本次启动写入 `sessions` 表

关键文件：
- `internal/service/launcher_service.go`
- `internal/service/sync_service.go`
- `internal/service/terminal_service.go`

---

### 3. 两套“会话”能力

项目里有两类会话，不要混淆：

#### 启动历史
保存在 SQLite `sessions` 表里，表示“某次启动了哪个 CLI、用了什么配置、在哪个目录和终端启动”。

#### CLI 实际对话历史
由 `CliSessionService` 直接读取 CLI 工具本地会话文件。

目前代码中：
- Claude 会话目录来自 `~/.claude/projects`
- Codex / OpenCode 的本地会话读取暂未实现

---

## 前端说明

前端不是直接在页面里调用 Wails 生成代码，而是统一通过：

- `frontend/src/api/index.ts`

这一层做封装。它还支持在没有 Wails runtime 的情况下返回空数据，方便单独调试前端页面。

路由使用 `createWebHashHistory()`，这是为了适配 Wails 打包后的本地资源加载场景。

---

## 数据与本地文件

### SQLite 数据库
数据库文件默认位于：

```text
~/.aiclimgr/data.db
```

### Claude 本地会话目录

```text
~/.claude/projects
```

---

## 安全相关

- Provider API Key 会在本地加密存储
- Proxy 密码也会加密存储
- 前端拿到的是脱敏后的 API Key，不是明文
- 启动 CLI 时会过滤部分环境变量（例如 `CLAUDECODE`），避免嵌套启动问题

---

## 当前仓库现状

- 已有完整的 Wails + Vue + SQLite 基础骨架
- 主要功能模块已经分层完成
- 当前没有单独配置 lint 命令
- 当前测试覆盖较少

---

## 开源协作

- 贡献流程：[`CONTRIBUTING.md`](./CONTRIBUTING.md)
- 安全漏洞反馈：[`SECURITY.md`](./SECURITY.md)
- 社区行为准则：[`CODE_OF_CONDUCT.md`](./CODE_OF_CONDUCT.md)
- 开源许可证：[`LICENSE`](./LICENSE)

提交 Issue 或 PR 前，建议先阅读以上文档。

---

## Release 产物

从 `v0.1.2` 开始，仓库包含面向 GitHub Actions 的多平台 Release 流程：

- `push` / `pull_request` 会自动执行 `go test ./...` 与前端构建
- `v*` tag 或手动触发可执行 Release 构建工作流
- Release 工作流会先构建前端，再执行 Wails 多平台打包并统一上传 GitHub Release 附件

当前支持的发布矩阵：

- Windows amd64
- Windows arm64
- Linux amd64
- Linux arm64
- macOS amd64
- macOS arm64

当前产物格式：

- Windows：`.zip`
- Linux：`.tar.gz`
- macOS：`.zip`

Release 附件命名格式：

- `AiCliManager-${RELEASE_VERSION}-windows-amd64.zip`
- `AiCliManager-${RELEASE_VERSION}-windows-arm64.zip`
- `AiCliManager-${RELEASE_VERSION}-linux-amd64.tar.gz`
- `AiCliManager-${RELEASE_VERSION}-linux-arm64.tar.gz`
- `AiCliManager-${RELEASE_VERSION}-macos-amd64.zip`
- `AiCliManager-${RELEASE_VERSION}-macos-arm64.zip`

当前明确不提供：

- `.dmg`
- `.deb`
- `.rpm`
- `.msi`

注意事项：

- macOS 当前为未签名构建，首次运行时可能需要用户在系统安全设置中手动放行。
- Linux 当前仅提供 `.tar.gz` 主产物，不包含 `.deb`、`.rpm` 或 AppImage。

发布方式建议：

1. 确认 `go test ./...` 与 `cd frontend && npm run build` 通过
2. 确认 `wails.json` 与 `frontend/package.json` 版本一致
3. 创建并推送版本 tag，例如 `v0.1.2`
4. 等待 GitHub Actions 完成 Release 工作流
5. 到 GitHub Release 页面下载对应平台与架构的附件

你可以在 GitHub Release 页面下载已发布版本的构建附件。

---

## 相关文件

- `CLAUDE.md`：提供给 Claude Code 的仓库协作说明
- `wails.json`：Wails 构建与前端命令配置
- `frontend/package.json`：前端脚本与依赖
