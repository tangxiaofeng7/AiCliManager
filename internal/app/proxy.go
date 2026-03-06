package app

import (
	"AiCliManager/internal/db/models"
)

// CreateProxyRequest 创建代理配置的请求参数
type CreateProxyRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	NoProxy  string `json:"no_proxy"`
}

// UpdateProxyRequest 更新代理配置的请求参数
type UpdateProxyRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"` // 为空时保持原值
	NoProxy  string `json:"no_proxy"`
}

// GetProxies 获取所有代理配置列表（密码脱敏）
func (a *App) GetProxies() ([]models.Proxy, error) {
	return a.proxyService.GetAll()
}

// CreateProxy 创建新代理配置
func (a *App) CreateProxy(req CreateProxyRequest) (*models.Proxy, error) {
	if req.Name == "" {
		return nil, errInvalidParam("代理名称不能为空")
	}
	if req.Host == "" {
		return nil, errInvalidParam("代理主机不能为空")
	}
	if req.Port <= 0 || req.Port > 65535 {
		return nil, errInvalidParam("代理端口无效，范围为 1-65535")
	}
	proxyType := req.Type
	if proxyType == "" {
		proxyType = "http"
	}
	return a.proxyService.Create(req.Name, proxyType, req.Host, req.Port, req.Username, req.Password, req.NoProxy)
}

// UpdateProxy 更新代理配置
func (a *App) UpdateProxy(id int64, req UpdateProxyRequest) error {
	if req.Name == "" {
		return errInvalidParam("代理名称不能为空")
	}
	if req.Host == "" {
		return errInvalidParam("代理主机不能为空")
	}
	return a.proxyService.Update(id, req.Name, req.Type, req.Host, req.Port, req.Username, req.Password, req.NoProxy)
}

// DeleteProxy 删除代理配置
func (a *App) DeleteProxy(id int64) error {
	return a.proxyService.Delete(id)
}

// SetGlobalProxy 设置全局激活代理（互斥：同时只有一条 is_active=1）
func (a *App) SetGlobalProxy(id int64) error {
	return a.proxyService.SetGlobalProxy(id)
}

// ClearGlobalProxy 清除全局激活代理
func (a *App) ClearGlobalProxy() error {
	return a.proxyService.ClearGlobalProxy()
}
