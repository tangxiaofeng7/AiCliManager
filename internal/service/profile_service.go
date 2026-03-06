package service

import (
	"AiCliManager/internal/db/models"
	"database/sql"
	"fmt"
	"time"
)

// ProfileService 提供 Profile 相关的业务逻辑
type ProfileService struct {
	db *sql.DB
}

// NewProfileService 创建 ProfileService 实例
func NewProfileService(db *sql.DB) *ProfileService {
	return &ProfileService{db: db}
}

// GetAll 获取所有 Profile 列表
func (s *ProfileService) GetAll() ([]models.Profile, error) {
	rows, err := s.db.Query(
		`SELECT id, name, provider_id, model, system_prompt, temperature, max_tokens, extra_config, created_at, updated_at
		 FROM profiles ORDER BY id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("查询 Profile 列表失败: %w", err)
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var p models.Profile
		if err := rows.Scan(
			&p.Id, &p.Name, &p.ProviderId, &p.Model, &p.SystemPrompt,
			&p.Temperature, &p.MaxTokens, &p.ExtraConfig, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描 Profile 行数据失败: %w", err)
		}
		profiles = append(profiles, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历 Profile 行数据失败: %w", err)
	}
	return profiles, nil
}

// GetById 按 ID 查询 Profile
func (s *ProfileService) GetById(id int64) (*models.Profile, error) {
	row := s.db.QueryRow(
		`SELECT id, name, provider_id, model, system_prompt, temperature, max_tokens, extra_config, created_at, updated_at
		 FROM profiles WHERE id=?`, id,
	)
	var p models.Profile
	if err := row.Scan(
		&p.Id, &p.Name, &p.ProviderId, &p.Model, &p.SystemPrompt,
		&p.Temperature, &p.MaxTokens, &p.ExtraConfig, &p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Profile (id=%d) 不存在", id)
		}
		return nil, fmt.Errorf("查询 Profile 失败: %w", err)
	}
	return &p, nil
}

// Create 创建新 Profile
func (s *ProfileService) Create(name string, providerId int64, model, systemPrompt string, temperature float64, maxTokens int, extraConfig string) (*models.Profile, error) {
	if extraConfig == "" {
		extraConfig = "{}"
	}
	now := time.Now()
	result, err := s.db.Exec(
		`INSERT INTO profiles (name, provider_id, model, system_prompt, temperature, max_tokens, extra_config, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		name, providerId, model, systemPrompt, temperature, maxTokens, extraConfig, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("创建 Profile 失败: %w", err)
	}
	id, _ := result.LastInsertId()
	return &models.Profile{
		Id:           id,
		Name:         name,
		ProviderId:   providerId,
		Model:        model,
		SystemPrompt: systemPrompt,
		Temperature:  temperature,
		MaxTokens:    maxTokens,
		ExtraConfig:  extraConfig,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Update 更新 Profile
func (s *ProfileService) Update(id int64, name string, providerId int64, model, systemPrompt string, temperature float64, maxTokens int, extraConfig string) error {
	if extraConfig == "" {
		extraConfig = "{}"
	}
	now := time.Now()
	_, err := s.db.Exec(
		`UPDATE profiles SET name=?, provider_id=?, model=?, system_prompt=?, temperature=?, max_tokens=?, extra_config=?, updated_at=?
		 WHERE id=?`,
		name, providerId, model, systemPrompt, temperature, maxTokens, extraConfig, now, id,
	)
	if err != nil {
		return fmt.Errorf("更新 Profile 失败: %w", err)
	}
	return nil
}

// Delete 删除 Profile
func (s *ProfileService) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM profiles WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("删除 Profile 失败: %w", err)
	}
	return nil
}
