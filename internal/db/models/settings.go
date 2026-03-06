package models

import "time"

// Settings 表示一条全局设置的 KV 记录
type Settings struct {
	Key         string    `json:"key"         orm:"key"`
	Value       string    `json:"value"       orm:"value"`
	Description string    `json:"description" orm:"description"`
	UpdatedAt   time.Time `json:"updated_at"  orm:"updated_at"`
}
