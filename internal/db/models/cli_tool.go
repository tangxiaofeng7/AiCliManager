package models

import "time"

// CliTool 表示一个支持的 CLI 工具记录
type CliTool struct {
	Id                int64     `json:"id"                 orm:"id"`
	Name              string    `json:"name"               orm:"name"`
	Key               string    `json:"key"                orm:"key"`
	Executable        string    `json:"executable"         orm:"executable"`
	ConfigPath        string    `json:"config_path"        orm:"config_path"`
	PreferredTerminal string    `json:"preferred_terminal" orm:"preferred_terminal"`
	IsInstalled       int       `json:"is_installed"       orm:"is_installed"`
	IsEnabled         int       `json:"is_enabled"         orm:"is_enabled"`
	SortOrder         int       `json:"sort_order"         orm:"sort_order"`
	CreatedAt         time.Time `json:"created_at"         orm:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"         orm:"updated_at"`
}
