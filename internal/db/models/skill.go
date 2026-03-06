package models

import "time"

// Skill 表示一个技能/指令模板记录
type Skill struct {
	Id        int64     `json:"id"         orm:"id"`
	Name      string    `json:"name"       orm:"name"`
	Category  string    `json:"category"   orm:"category"`
	Trigger   string    `json:"trigger"    orm:"trigger"`
	Content   string    `json:"content"    orm:"content"`
	Variables string    `json:"variables"  orm:"variables"` // JSON 变量定义
	IsBuiltin int       `json:"is_builtin" orm:"is_builtin"`
	SortOrder int       `json:"sort_order" orm:"sort_order"`
	CreatedAt time.Time `json:"created_at" orm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" orm:"updated_at"`
}
