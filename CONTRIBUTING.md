# Contributing to AiCliManager

感谢你关注 AiCliManager。

## 开始之前

提交 Issue 或 PR 前，请先阅读以下文档：

- [README.md](./README.md)
- [SECURITY.md](./SECURITY.md)
- [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md)

## 开发环境

请确保本地具备以下环境：

- Go 1.24+
- Node.js 20+
- Wails CLI

安装依赖：

```bash
go mod tidy
cd frontend && npm install
```

本地开发：

```bash
wails dev
```

仅运行前端：

```bash
cd frontend && npm run dev
```

## 提交前自检

请在提交前至少执行以下检查：

```bash
go test ./...
cd frontend && npm run build
```

当前仓库没有独立 lint 命令，请不要在 PR 中引用不存在的 lint 流程。

## 提交与 PR 说明

为了便于维护，请尽量做到：

- 一个 PR 只解决一个明确问题
- 变更说明聚焦“为什么改”以及“如何验证”
- 涉及 UI 的变更附截图或录屏
- 涉及平台差异时说明操作系统与终端环境
- 不要手改 `frontend/wailsjs/` 生成文件

## Bug 反馈建议包含的信息

提交 bug 时，建议提供：

- 应用版本
- 操作系统版本
- 使用的 CLI 工具类型与终端类型
- 复现步骤
- 实际行为与期望行为
- 相关日志、报错、截图
- 是否与 Provider / Proxy / MCP / Skill 配置有关

## 代码约定

- 代码注释使用中文，标识符保持英文
- 尽量保持现有分层：`internal/app` 负责绑定，核心逻辑放在 `internal/service`
- 启动、配置同步、加密、会话解析相关逻辑不要直接下沉到前端页面
- Schema 变更请直接修改 `internal/db/db.go`，不要引入 migration 框架

## 安全问题

如果你发现的是安全漏洞，请不要公开提交 Issue。

请改为阅读并遵循 [SECURITY.md](./SECURITY.md) 中的私下反馈流程。
