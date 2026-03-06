package app

import (
	"AiCliManager/internal/db/models"
)

// CreateProfileRequest 创建 Profile 的请求参数
type CreateProfileRequest struct {
	Name         string  `json:"name"`
	ProviderId   int64   `json:"provider_id"`
	Model        string  `json:"model"`
	SystemPrompt string  `json:"system_prompt"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int     `json:"max_tokens"`
	ExtraConfig  string  `json:"extra_config"`
}

// UpdateProfileRequest 更新 Profile 的请求参数
type UpdateProfileRequest struct {
	Name         string  `json:"name"`
	ProviderId   int64   `json:"provider_id"`
	Model        string  `json:"model"`
	SystemPrompt string  `json:"system_prompt"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int     `json:"max_tokens"`
	ExtraConfig  string  `json:"extra_config"`
}

// GetProfiles 获取所有 Profile 列表
func (a *App) GetProfiles() ([]models.Profile, error) {
	return a.profileService.GetAll()
}

// CreateProfile 创建新 Profile
func (a *App) CreateProfile(req CreateProfileRequest) (*models.Profile, error) {
	if req.Name == "" {
		return nil, errInvalidParam("Profile 名称不能为空")
	}
	if req.ProviderId <= 0 {
		return nil, errInvalidParam("必须关联一个有效的 Provider")
	}
	if req.Model == "" {
		return nil, errInvalidParam("模型名称不能为空")
	}
	// 默认参数
	if req.Temperature == 0 {
		req.Temperature = 1.0
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 8192
	}
	return a.profileService.Create(
		req.Name, req.ProviderId, req.Model,
		req.SystemPrompt, req.Temperature, req.MaxTokens, req.ExtraConfig,
	)
}

// UpdateProfile 更新 Profile
func (a *App) UpdateProfile(id int64, req UpdateProfileRequest) error {
	if req.Name == "" {
		return errInvalidParam("Profile 名称不能为空")
	}
	return a.profileService.Update(
		id, req.Name, req.ProviderId, req.Model,
		req.SystemPrompt, req.Temperature, req.MaxTokens, req.ExtraConfig,
	)
}

// DeleteProfile 删除 Profile
func (a *App) DeleteProfile(id int64) error {
	return a.profileService.Delete(id)
}
