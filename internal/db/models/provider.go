package models

import "time"

// Provider 表示一个 API 提供商记录
type Provider struct {
	Id        int64     `json:"id"         orm:"id"`
	Name      string    `json:"name"       orm:"name"`
	Type      string    `json:"type"       orm:"type"`       // anthropic | openai | custom
	ApiUrl    string    `json:"api_url"    orm:"api_url"`
	ApiKey    string    `json:"api_key"    orm:"api_key"`    // 加密存储
	Models    string    `json:"models"     orm:"models"`     // JSON 数组字符串
	SortOrder int       `json:"sort_order" orm:"sort_order"`
	CreatedAt time.Time `json:"created_at" orm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" orm:"updated_at"`
}
