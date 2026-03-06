package models

import "time"

// Session 表示一次 CLI 工具启动会话记录
type Session struct {
	Id         int64      `json:"id"          orm:"id"`
	CliToolId  int64      `json:"cli_tool_id" orm:"cli_tool_id"`
	ProfileId  *int64     `json:"profile_id"  orm:"profile_id"`
	ProxyId    *int64     `json:"proxy_id"    orm:"proxy_id"`
	Terminal   string     `json:"terminal"    orm:"terminal"`
	WorkingDir string     `json:"working_dir" orm:"working_dir"`
	ExtraArgs  string     `json:"extra_args"  orm:"extra_args"` // JSON 数组
	Status     string     `json:"status"      orm:"status"`     // running | exited | error
	Pid        int        `json:"pid"         orm:"pid"`
	StartedAt  time.Time  `json:"started_at"  orm:"started_at"`
	EndedAt    *time.Time `json:"ended_at"    orm:"ended_at"`
}
