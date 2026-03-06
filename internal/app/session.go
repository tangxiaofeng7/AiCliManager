package app

import (
	"AiCliManager/internal/db/models"
	"AiCliManager/internal/service"
)

// GetSessionsRequest 查询会话列表的请求参数（透传给 service 层）
type GetSessionsRequest struct {
	CliToolKey string `json:"cli_tool_key"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

// GetSessions 获取会话历史列表（AiCliManager 启动记录），支持分页和工具筛选
func (a *App) GetSessions(req GetSessionsRequest) ([]models.Session, error) {
	return a.sessionService.GetAll(service.GetSessionsRequest{
		CliToolKey: req.CliToolKey,
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
}

// DeleteSession 删除单条会话记录
func (a *App) DeleteSession(id int64) error {
	return a.sessionService.Delete(id)
}

// ClearSessions 清空指定工具的历史（cliToolKey 为空则清空全部）
func (a *App) ClearSessions(cliToolKey string) error {
	return a.sessionService.ClearByCliTool(cliToolKey)
}

// RelaunchSession 用历史会话配置重新启动 CLI 工具
func (a *App) RelaunchSession(sessionId int64) error {
	return a.launcherService.RelaunchSession(sessionId)
}

// ---- CLI 工具实际对话会话 ----

// GetCliSessionsRequest 查询 CLI 工具对话会话的请求参数
type GetCliSessionsRequest struct {
	CliToolKey string `json:"cli_tool_key"`
	Project    string `json:"project"`
	Limit      int    `json:"limit"`
}

// GetCliSessions 获取 CLI 工具的实际对话会话列表（读取 CLI 工具本地会话文件）
func (a *App) GetCliSessions(req GetCliSessionsRequest) ([]service.CliSession, error) {
	return a.cliSessionService.GetSessions(service.GetCliSessionsRequest{
		CliToolKey: req.CliToolKey,
		Project:    req.Project,
		Limit:      req.Limit,
	})
}

// GetCliSessionMessages 获取指定对话会话的完整消息列表
func (a *App) GetCliSessionMessages(cliToolKey, sessionId string) ([]service.CliSessionMessage, error) {
	return a.cliSessionService.GetMessages(cliToolKey, sessionId)
}

// GetCliSessionProjects 获取 CLI 工具的所有项目目录列表
func (a *App) GetCliSessionProjects(cliToolKey string) ([]service.CliSessionProject, error) {
	return a.cliSessionService.GetProjects(cliToolKey)
}
