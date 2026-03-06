package service

import (
	"database/sql"
	"fmt"
	"time"
)

// SettingsService 提供全局设置 KV 读写业务逻辑
type SettingsService struct {
	db *sql.DB
}

// NewSettingsService 创建 SettingsService 实例
func NewSettingsService(db *sql.DB) *SettingsService {
	return &SettingsService{db: db}
}

// Get 按 key 读取设置值，不存在时返回 defaultValue
func (s *SettingsService) Get(key, defaultValue string) (string, error) {
	row := s.db.QueryRow(`SELECT value FROM settings WHERE key=?`, key)
	var value string
	if err := row.Scan(&value); err != nil {
		if err == sql.ErrNoRows {
			return defaultValue, nil
		}
		return defaultValue, fmt.Errorf("读取设置 [%s] 失败: %w", key, err)
	}
	return value, nil
}

// Set 写入或更新一条设置
func (s *SettingsService) Set(key, value string) error {
	now := time.Now()
	_, err := s.db.Exec(
		`INSERT INTO settings (key, value, updated_at) VALUES (?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET value=excluded.value, updated_at=excluded.updated_at`,
		key, value, now,
	)
	if err != nil {
		return fmt.Errorf("写入设置 [%s] 失败: %w", key, err)
	}
	return nil
}

// GetAll 获取所有设置项
func (s *SettingsService) GetAll() (map[string]string, error) {
	rows, err := s.db.Query(`SELECT key, value FROM settings ORDER BY key ASC`)
	if err != nil {
		return nil, fmt.Errorf("查询所有设置失败: %w", err)
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, fmt.Errorf("扫描设置行数据失败: %w", err)
		}
		result[k] = v
	}
	return result, nil
}

// Delete 删除一条设置
func (s *SettingsService) Delete(key string) error {
	_, err := s.db.Exec(`DELETE FROM settings WHERE key=?`, key)
	if err != nil {
		return fmt.Errorf("删除设置 [%s] 失败: %w", key, err)
	}
	return nil
}
