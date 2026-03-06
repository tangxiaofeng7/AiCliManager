package app

import (
	"AiCliManager/internal/db/models"
	"AiCliManager/internal/service"
)

// LaunchRequest 启动 CLI 工具的请求参数（透传给 service 层）
type LaunchRequest struct {
	CliToolKey   string            `json:"cli_tool_key"`
	ProfileId    int64             `json:"profile_id"`
	ProxyId      *int64            `json:"proxy_id"`
	McpServerIds []int64           `json:"mcp_server_ids"`
	SkillIds     []int64           `json:"skill_ids"`
	SkillVars    map[string]string `json:"skill_vars"`
	ExtraArgs    []string          `json:"extra_args"`
	Terminal     string            `json:"terminal"`
	WorkingDir   string            `json:"working_dir"`
}

// GetCliTools 获取所有 CLI 工具列表及安装状态
func (a *App) GetCliTools() ([]models.CliTool, error) {
	return a.launcherService.GetCliTools()
}

// DetectCliTool 检测指定 CLI 工具是否已安装
func (a *App) DetectCliTool(key string) (*service.DetectResult, error) {
	if key == "" {
		return nil, errInvalidParam("CLI 工具 key 不能为空")
	}
	return a.launcherService.DetectCliTool(key)
}

// LaunchCliTool 启动指定 CLI 工具
func (a *App) LaunchCliTool(req LaunchRequest) error {
	if req.CliToolKey == "" {
		return errInvalidParam("CLI 工具 key 不能为空")
	}
	return a.launcherService.Launch(service.LaunchRequest{
		CliToolKey:   req.CliToolKey,
		ProfileId:    req.ProfileId,
		ProxyId:      req.ProxyId,
		McpServerIds: req.McpServerIds,
		SkillIds:     req.SkillIds,
		SkillVars:    req.SkillVars,
		ExtraArgs:    req.ExtraArgs,
		Terminal:     req.Terminal,
		WorkingDir:   req.WorkingDir,
	})
}

// GetCliToolActiveConfig 获取指定 CLI 工具当前激活的 Profile 和 Proxy
func (a *App) GetCliToolActiveConfig(key string) (*service.ActiveConfig, error) {
	if key == "" {
		return nil, errInvalidParam("CLI 工具 key 不能为空")
	}
	return a.launcherService.GetCliToolActiveConfig(key)
}

// SetActiveConfigRequest 设置激活配置的请求参数
type SetActiveConfigRequest struct {
	CliToolKey string `json:"cli_tool_key"`
	ProfileId  *int64 `json:"profile_id"`
	ProxyId    *int64 `json:"proxy_id"`
}

// SetCliToolActiveConfig 设置指定 CLI 工具的激活 Profile 和 Proxy
func (a *App) SetCliToolActiveConfig(req SetActiveConfigRequest) error {
	if req.CliToolKey == "" {
		return errInvalidParam("CLI 工具 key 不能为空")
	}
	return a.launcherService.SetCliToolActiveConfig(service.SetActiveConfigRequest{
		CliToolKey: req.CliToolKey,
		ProfileId:  req.ProfileId,
		ProxyId:    req.ProxyId,
	})
}

// ListAvailableTerminals 返回当前平台可用的终端列表
func (a *App) ListAvailableTerminals() []service.TerminalInfo {
	return a.terminalService.ListAvailableTerminals()
}
