package models

import "time"

// Proxy 表示一条代理配置记录
type Proxy struct {
	Id        int64     `json:"id"         orm:"id"`
	Name      string    `json:"name"       orm:"name"`
	Type      string    `json:"type"       orm:"type"`     // http | https | socks5
	Host      string    `json:"host"       orm:"host"`
	Port      int       `json:"port"       orm:"port"`
	Username  string    `json:"username"   orm:"username"`
	Password  string    `json:"password"   orm:"password"` // 加密存储
	NoProxy   string    `json:"no_proxy"   orm:"no_proxy"`
	IsActive  int       `json:"is_active"  orm:"is_active"`
	CreatedAt time.Time `json:"created_at" orm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" orm:"updated_at"`
}
