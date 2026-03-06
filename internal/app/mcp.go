package app

import (
	"AiCliManager/internal/db/models"
)

// CreateMcpServerRequest 创建 MCP Server 的请求参数
type CreateMcpServerRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Command     string `json:"command"`
	Args        string `json:"args"`
	Env         string `json:"env"`
	Url         string `json:"url"`
	Description string `json:"description"`
	IsEnabled   int    `json:"is_enabled"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateMcpServerRequest 更新 MCP Server 的请求参数
type UpdateMcpServerRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Command     string `json:"command"`
	Args        string `json:"args"`
	Env         string `json:"env"`
	Url         string `json:"url"`
	Description string `json:"description"`
	IsEnabled   int    `json:"is_enabled"`
	SortOrder   int    `json:"sort_order"`
}

// GetMcpServers 获取所有 MCP Server 列表
func (a *App) GetMcpServers() ([]models.McpServer, error) {
	return a.mcpService.GetAll()
}

// CreateMcpServer 创建新 MCP Server
func (a *App) CreateMcpServer(req CreateMcpServerRequest) (*models.McpServer, error) {
	if req.Name == "" {
		return nil, errInvalidParam("MCP Server 名称不能为空")
	}
	if req.Type == "" {
		return nil, errInvalidParam("MCP Server 类型不能为空")
	}
	// 校验类型合法性
	switch req.Type {
	case "stdio", "sse", "http":
	default:
		return nil, errInvalidParam("MCP Server 类型无效，支持：stdio | sse | http")
	}
	// stdio 类型必须有 command
	if req.Type == "stdio" && req.Command == "" {
		return nil, errInvalidParam("stdio 类型的 MCP Server 必须填写启动命令")
	}
	// sse/http 类型必须有 url
	if (req.Type == "sse" || req.Type == "http") && req.Url == "" {
		return nil, errInvalidParam("sse/http 类型的 MCP Server 必须填写服务地址")
	}
	// 默认启用
	if req.IsEnabled == 0 {
		req.IsEnabled = 1
	}
	return a.mcpService.Create(req.Name, req.Type, req.Command, req.Args, req.Env, req.Url, req.Description, req.IsEnabled, req.SortOrder)
}

// UpdateMcpServer 更新 MCP Server
func (a *App) UpdateMcpServer(id int64, req UpdateMcpServerRequest) error {
	if req.Name == "" {
		return errInvalidParam("MCP Server 名称不能为空")
	}
	return a.mcpService.Update(id, req.Name, req.Type, req.Command, req.Args, req.Env, req.Url, req.Description, req.IsEnabled, req.SortOrder)
}

// DeleteMcpServer 删除 MCP Server
func (a *App) DeleteMcpServer(id int64) error {
	return a.mcpService.Delete(id)
}

// GetCliToolMcpServers 获取指定 CLI 工具关联的 MCP Server 列表
func (a *App) GetCliToolMcpServers(cliToolKey string) ([]models.McpServer, error) {
	if cliToolKey == "" {
		return nil, errInvalidParam("CLI 工具 key 不能为空")
	}
	return a.mcpService.GetCliToolMcpServers(cliToolKey)
}

// SetCliToolMcpServers 设置指定 CLI 工具关联的 MCP Server 列表（全量替换）
func (a *App) SetCliToolMcpServers(cliToolKey string, ids []int64) error {
	if cliToolKey == "" {
		return errInvalidParam("CLI 工具 key 不能为空")
	}
	return a.mcpService.SetCliToolMcpServers(cliToolKey, ids)
}
