package app

import (
	"AiCliManager/internal/db/models"
	"AiCliManager/internal/service"
)

// CreateProviderRequest 创建 Provider 的请求参数
type CreateProviderRequest struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	ApiUrl    string `json:"api_url"`
	ApiKey    string `json:"api_key"`
	Models    string `json:"models"`
	SortOrder int    `json:"sort_order"`
}

// UpdateProviderRequest 更新 Provider 的请求参数
type UpdateProviderRequest struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	ApiUrl    string `json:"api_url"`
	ApiKey    string `json:"api_key"` // 为空时保持原值
	Models    string `json:"models"`
	SortOrder int    `json:"sort_order"`
}

// GetProviders 获取所有 Provider 列表（API Key 脱敏）
func (a *App) GetProviders() ([]models.Provider, error) {
	return a.providerService.GetAll()
}

// CreateProvider 创建新 Provider
func (a *App) CreateProvider(req CreateProviderRequest) (*models.Provider, error) {
	if req.Name == "" {
		return nil, errInvalidParam("Provider 名称不能为空")
	}
	if req.ApiUrl == "" {
		return nil, errInvalidParam("API URL 不能为空")
	}
	if req.ApiKey == "" {
		return nil, errInvalidParam("API Key 不能为空")
	}
	return a.providerService.Create(req.Name, req.Type, req.ApiUrl, req.ApiKey, req.Models, req.SortOrder)
}

// UpdateProvider 更新 Provider 信息
func (a *App) UpdateProvider(id int64, req UpdateProviderRequest) error {
	if req.Name == "" {
		return errInvalidParam("Provider 名称不能为空")
	}
	return a.providerService.Update(id, req.Name, req.Type, req.ApiUrl, req.ApiKey, req.Models, req.SortOrder)
}

// DeleteProvider 删除 Provider
func (a *App) DeleteProvider(id int64) error {
	return a.providerService.Delete(id)
}

// TestProvider 测试 Provider 连通性
func (a *App) TestProvider(id int64) (*service.TestResult, error) {
	return a.providerService.TestProvider(id)
}

// FetchProviderModels 从 Provider API 拉取并更新可用模型列表
func (a *App) FetchProviderModels(id int64) ([]string, error) {
	return a.providerService.FetchModels(id)
}
