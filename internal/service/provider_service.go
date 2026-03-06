package service

import (
	"AiCliManager/internal/crypto"
	"AiCliManager/internal/db/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ProviderService 提供 Provider 相关的业务逻辑
type ProviderService struct {
	db *sql.DB
}

// NewProviderService 创建 ProviderService 实例
func NewProviderService(db *sql.DB) *ProviderService {
	return &ProviderService{db: db}
}

// GetAll 获取所有 Provider 列表（API Key 脱敏返回）
func (s *ProviderService) GetAll() ([]models.Provider, error) {
	rows, err := s.db.Query(
		`SELECT id, name, type, api_url, api_key, models, sort_order, created_at, updated_at
		 FROM providers ORDER BY sort_order ASC, id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("查询 Provider 列表失败: %w", err)
	}
	defer rows.Close()

	var providers []models.Provider
	for rows.Next() {
		var p models.Provider
		if err := rows.Scan(
			&p.Id, &p.Name, &p.Type, &p.ApiUrl, &p.ApiKey,
			&p.Models, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描 Provider 行数据失败: %w", err)
		}
		// API Key 脱敏处理：解密后只返回掩码，不暴露原文
		if p.ApiKey != "" {
			decrypted, decErr := crypto.Decrypt(p.ApiKey)
			if decErr == nil {
				p.ApiKey = crypto.MaskApiKey(decrypted)
			} else {
				p.ApiKey = "****"
			}
		}
		providers = append(providers, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历 Provider 行数据失败: %w", err)
	}
	return providers, nil
}

// GetById 按 ID 查询 Provider（内部使用，返回加密的 ApiKey）
func (s *ProviderService) GetById(id int64) (*models.Provider, error) {
	row := s.db.QueryRow(
		`SELECT id, name, type, api_url, api_key, models, sort_order, created_at, updated_at
		 FROM providers WHERE id = ?`, id,
	)
	var p models.Provider
	if err := row.Scan(
		&p.Id, &p.Name, &p.Type, &p.ApiUrl, &p.ApiKey,
		&p.Models, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Provider (id=%d) 不存在", id)
		}
		return nil, fmt.Errorf("查询 Provider 失败: %w", err)
	}
	return &p, nil
}

// Create 创建新 Provider，API Key 加密存储
func (s *ProviderService) Create(name, providerType, apiUrl, apiKey, modelsJSON string, sortOrder int) (*models.Provider, error) {
	// 加密 API Key
	encryptedKey, err := crypto.Encrypt(apiKey)
	if err != nil {
		return nil, fmt.Errorf("加密 API Key 失败: %w", err)
	}
	if modelsJSON == "" {
		modelsJSON = "[]"
	}

	now := time.Now()
	result, err := s.db.Exec(
		`INSERT INTO providers (name, type, api_url, api_key, models, sort_order, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		name, providerType, apiUrl, encryptedKey, modelsJSON, sortOrder, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("创建 Provider 失败: %w", err)
	}
	id, _ := result.LastInsertId()

	// 返回时 API Key 脱敏
	p := &models.Provider{
		Id:        id,
		Name:      name,
		Type:      providerType,
		ApiUrl:    apiUrl,
		ApiKey:    crypto.MaskApiKey(apiKey),
		Models:    modelsJSON,
		SortOrder: sortOrder,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return p, nil
}

// Update 更新 Provider 信息，若 apiKey 不为空则重新加密存储
func (s *ProviderService) Update(id int64, name, providerType, apiUrl, apiKey, modelsJSON string, sortOrder int) error {
	// 先查询旧记录
	old, err := s.GetById(id)
	if err != nil {
		return err
	}

	// 决定使用的加密 Key
	encryptedKey := old.ApiKey // 默认保持原加密值
	if apiKey != "" {
		encryptedKey, err = crypto.Encrypt(apiKey)
		if err != nil {
			return fmt.Errorf("加密 API Key 失败: %w", err)
		}
	}
	if modelsJSON == "" {
		modelsJSON = old.Models
	}

	now := time.Now()
	_, err = s.db.Exec(
		`UPDATE providers SET name=?, type=?, api_url=?, api_key=?, models=?, sort_order=?, updated_at=?
		 WHERE id=?`,
		name, providerType, apiUrl, encryptedKey, modelsJSON, sortOrder, now, id,
	)
	if err != nil {
		return fmt.Errorf("更新 Provider 失败: %w", err)
	}
	return nil
}

// Delete 删除 Provider
func (s *ProviderService) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM providers WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("删除 Provider 失败: %w", err)
	}
	return nil
}

// GetDecryptedApiKey 获取解密后的 API Key（内部使用，禁止返回给前端）
func (s *ProviderService) GetDecryptedApiKey(id int64) (string, error) {
	p, err := s.GetById(id)
	if err != nil {
		return "", err
	}
	decrypted, err := crypto.Decrypt(p.ApiKey)
	if err != nil {
		return "", fmt.Errorf("解密 API Key 失败: %w", err)
	}
	return decrypted, nil
}

// TestResult 连通性测试结果
type TestResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Latency int64  `json:"latency_ms"` // 响应延迟（毫秒）
}

// TestProvider 测试 Provider 连通性
func (s *ProviderService) TestProvider(id int64) (*TestResult, error) {
	p, err := s.GetById(id)
	if err != nil {
		return nil, err
	}

	// 解密 API Key
	apiKey, err := crypto.Decrypt(p.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("解密 API Key 失败: %w", err)
	}

	// 构建测试请求 URL
	baseURL := strings.TrimRight(p.ApiUrl, "/")
	testURL := baseURL + "/v1/models"

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建测试请求失败: %w", err)
	}

	// 根据 Provider 类型设置认证头
	switch p.Type {
	case "anthropic":
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	default:
		// openai 和 custom 使用 Bearer token
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并计算延迟
	client := &http.Client{Timeout: 15 * time.Second}
	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return &TestResult{
			Success: false,
			Message: fmt.Sprintf("连接失败: %v", err),
			Latency: latency,
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return &TestResult{
			Success: true,
			Message: fmt.Sprintf("连接成功（HTTP %d）", resp.StatusCode),
			Latency: latency,
		}, nil
	}

	// 读取错误响应体（最多 512 字节）
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
	return &TestResult{
		Success: false,
		Message: fmt.Sprintf("服务器返回 HTTP %d: %s", resp.StatusCode, string(body)),
		Latency: latency,
	}, nil
}

// FetchModels 从 Provider API 获取可用模型列表，并更新数据库
func (s *ProviderService) FetchModels(id int64) ([]string, error) {
	p, err := s.GetById(id)
	if err != nil {
		return nil, err
	}

	// 解密 API Key
	apiKey, err := crypto.Decrypt(p.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("解密 API Key 失败: %w", err)
	}

	baseURL := strings.TrimRight(p.ApiUrl, "/")
	reqURL := baseURL + "/v1/models"

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	switch p.Type {
	case "anthropic":
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	default:
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求模型列表失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("获取模型列表失败（HTTP %d）: %s", resp.StatusCode, string(body))
	}

	// 解析响应，兼容 OpenAI 和 Anthropic 格式
	var models []string
	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("解析模型列表响应失败: %w", err)
	}

	// OpenAI 格式：{"data": [{"id": "..."}, ...]}
	if data, ok := raw["data"].([]interface{}); ok {
		for _, item := range data {
			if m, ok := item.(map[string]interface{}); ok {
				if id, ok := m["id"].(string); ok {
					models = append(models, id)
				}
			}
		}
	}

	// Anthropic 格式：{"models": [{"id": "..."}, ...]} 或 {"data": [...]}
	if len(models) == 0 {
		if data, ok := raw["models"].([]interface{}); ok {
			for _, item := range data {
				if m, ok := item.(map[string]interface{}); ok {
					if id, ok := m["id"].(string); ok {
						models = append(models, id)
					}
				}
			}
		}
	}

	// 将模型列表序列化后更新数据库
	modelsJSON, err := json.Marshal(models)
	if err != nil {
		return nil, fmt.Errorf("序列化模型列表失败: %w", err)
	}
	_, err = s.db.Exec(
		`UPDATE providers SET models=?, updated_at=? WHERE id=?`,
		string(modelsJSON), time.Now(), id,
	)
	if err != nil {
		return nil, fmt.Errorf("更新模型列表失败: %w", err)
	}

	return models, nil
}
