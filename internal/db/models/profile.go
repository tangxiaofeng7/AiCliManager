package models

import "time"

// Profile 表示一个配置 Profile，绑定 Provider 与参数
type Profile struct {
	Id           int64     `json:"id"            orm:"id"`
	Name         string    `json:"name"          orm:"name"`
	ProviderId   int64     `json:"provider_id"   orm:"provider_id"`
	Model        string    `json:"model"         orm:"model"`
	SystemPrompt string    `json:"system_prompt" orm:"system_prompt"`
	Temperature  float64   `json:"temperature"   orm:"temperature"`
	MaxTokens    int       `json:"max_tokens"    orm:"max_tokens"`
	ExtraConfig  string    `json:"extra_config"  orm:"extra_config"` // JSON 扩展配置
	CreatedAt    time.Time `json:"created_at"    orm:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"    orm:"updated_at"`
}
