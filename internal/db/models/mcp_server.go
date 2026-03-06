package models

import "time"

// McpServer 表示一个 MCP Server 配置记录
type McpServer struct {
	Id          int64     `json:"id"          orm:"id"`
	Name        string    `json:"name"        orm:"name"`
	Type        string    `json:"type"        orm:"type"`        // stdio | sse | http
	Command     string    `json:"command"     orm:"command"`
	Args        string    `json:"args"        orm:"args"`        // JSON 数组
	Env         string    `json:"env"         orm:"env"`         // JSON 对象
	Url         string    `json:"url"         orm:"url"`
	Description string    `json:"description" orm:"description"`
	IsEnabled   int       `json:"is_enabled"  orm:"is_enabled"`
	SortOrder   int       `json:"sort_order"  orm:"sort_order"`
	CreatedAt   time.Time `json:"created_at"  orm:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"  orm:"updated_at"`
}
