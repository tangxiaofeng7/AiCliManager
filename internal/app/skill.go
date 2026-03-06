package app

import (
	"AiCliManager/internal/db/models"
)

// CreateSkillRequest 创建 Skill 的请求参数
type CreateSkillRequest struct {
	Name      string `json:"name"`
	Category  string `json:"category"`
	Trigger   string `json:"trigger"`
	Content   string `json:"content"`
	Variables string `json:"variables"`
	SortOrder int    `json:"sort_order"`
}

// UpdateSkillRequest 更新 Skill 的请求参数
type UpdateSkillRequest struct {
	Name      string `json:"name"`
	Category  string `json:"category"`
	Trigger   string `json:"trigger"`
	Content   string `json:"content"`
	Variables string `json:"variables"`
	SortOrder int    `json:"sort_order"`
}

// GetSkills 获取所有 Skill 列表
func (a *App) GetSkills() ([]models.Skill, error) {
	return a.skillService.GetAll()
}

// CreateSkill 创建新 Skill
func (a *App) CreateSkill(req CreateSkillRequest) (*models.Skill, error) {
	if req.Name == "" {
		return nil, errInvalidParam("Skill 名称不能为空")
	}
	if req.Content == "" {
		return nil, errInvalidParam("Skill 提示词内容不能为空")
	}
	category := req.Category
	if category == "" {
		category = "general"
	}
	return a.skillService.Create(req.Name, category, req.Trigger, req.Content, req.Variables, req.SortOrder)
}

// UpdateSkill 更新 Skill（内置 Skill 不可修改）
func (a *App) UpdateSkill(id int64, req UpdateSkillRequest) error {
	if req.Name == "" {
		return errInvalidParam("Skill 名称不能为空")
	}
	if req.Content == "" {
		return errInvalidParam("Skill 提示词内容不能为空")
	}
	return a.skillService.Update(id, req.Name, req.Category, req.Trigger, req.Content, req.Variables, req.SortOrder)
}

// DeleteSkill 删除 Skill（内置 Skill 不可删除）
func (a *App) DeleteSkill(id int64) error {
	return a.skillService.Delete(id)
}
