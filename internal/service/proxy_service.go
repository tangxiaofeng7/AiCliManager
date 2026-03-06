package service

import (
	"AiCliManager/internal/crypto"
	"AiCliManager/internal/db/models"
	"database/sql"
	"fmt"
	"time"
)

// ProxyService 提供代理配置相关的业务逻辑
type ProxyService struct {
	db *sql.DB
}

// NewProxyService 创建 ProxyService 实例
func NewProxyService(db *sql.DB) *ProxyService {
	return &ProxyService{db: db}
}

// GetAll 获取所有代理配置（密码脱敏）
func (s *ProxyService) GetAll() ([]models.Proxy, error) {
	rows, err := s.db.Query(
		`SELECT id, name, type, host, port, username, password, no_proxy, is_active, created_at, updated_at
		 FROM proxies ORDER BY id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("查询代理列表失败: %w", err)
	}
	defer rows.Close()

	var proxies []models.Proxy
	for rows.Next() {
		var p models.Proxy
		if err := rows.Scan(
			&p.Id, &p.Name, &p.Type, &p.Host, &p.Port,
			&p.Username, &p.Password, &p.NoProxy, &p.IsActive,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描代理行数据失败: %w", err)
		}
		// 密码脱敏
		if p.Password != "" {
			p.Password = "****"
		}
		proxies = append(proxies, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历代理行数据失败: %w", err)
	}
	return proxies, nil
}

// GetById 按 ID 查询代理（内部使用，返回加密密码）
func (s *ProxyService) GetById(id int64) (*models.Proxy, error) {
	row := s.db.QueryRow(
		`SELECT id, name, type, host, port, username, password, no_proxy, is_active, created_at, updated_at
		 FROM proxies WHERE id=?`, id,
	)
	var p models.Proxy
	if err := row.Scan(
		&p.Id, &p.Name, &p.Type, &p.Host, &p.Port,
		&p.Username, &p.Password, &p.NoProxy, &p.IsActive,
		&p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("代理 (id=%d) 不存在", id)
		}
		return nil, fmt.Errorf("查询代理失败: %w", err)
	}
	return &p, nil
}

// Create 创建新代理配置，密码加密存储
func (s *ProxyService) Create(name, proxyType, host string, port int, username, password, noProxy string) (*models.Proxy, error) {
	encPassword, err := crypto.Encrypt(password)
	if err != nil {
		return nil, fmt.Errorf("加密代理密码失败: %w", err)
	}
	now := time.Now()
	result, err := s.db.Exec(
		`INSERT INTO proxies (name, type, host, port, username, password, no_proxy, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`,
		name, proxyType, host, port, username, encPassword, noProxy, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("创建代理失败: %w", err)
	}
	id, _ := result.LastInsertId()
	masked := ""
	if password != "" {
		masked = "****"
	}
	return &models.Proxy{
		Id:        id,
		Name:      name,
		Type:      proxyType,
		Host:      host,
		Port:      port,
		Username:  username,
		Password:  masked,
		NoProxy:   noProxy,
		IsActive:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update 更新代理配置，若 password 不为空则重新加密存储
func (s *ProxyService) Update(id int64, name, proxyType, host string, port int, username, password, noProxy string) error {
	old, err := s.GetById(id)
	if err != nil {
		return err
	}
	encPassword := old.Password // 默认保持原加密值
	if password != "" {
		encPassword, err = crypto.Encrypt(password)
		if err != nil {
			return fmt.Errorf("加密代理密码失败: %w", err)
		}
	}
	now := time.Now()
	_, err = s.db.Exec(
		`UPDATE proxies SET name=?, type=?, host=?, port=?, username=?, password=?, no_proxy=?, updated_at=?
		 WHERE id=?`,
		name, proxyType, host, port, username, encPassword, noProxy, now, id,
	)
	if err != nil {
		return fmt.Errorf("更新代理失败: %w", err)
	}
	return nil
}

// Delete 删除代理配置
func (s *ProxyService) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM proxies WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("删除代理失败: %w", err)
	}
	return nil
}

// SetGlobalProxy 设置全局激活代理（互斥：先清空所有，再激活指定条目）
func (s *ProxyService) SetGlobalProxy(id int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback()

	// 清空所有激活状态
	if _, err := tx.Exec(`UPDATE proxies SET is_active=0, updated_at=?`, time.Now()); err != nil {
		return fmt.Errorf("清空代理激活状态失败: %w", err)
	}
	// 激活指定代理
	if _, err := tx.Exec(`UPDATE proxies SET is_active=1, updated_at=? WHERE id=?`, time.Now(), id); err != nil {
		return fmt.Errorf("设置全局代理失败: %w", err)
	}
	return tx.Commit()
}

// ClearGlobalProxy 清除全局激活代理
func (s *ProxyService) ClearGlobalProxy() error {
	_, err := s.db.Exec(`UPDATE proxies SET is_active=0, updated_at=?`, time.Now())
	if err != nil {
		return fmt.Errorf("清除全局代理失败: %w", err)
	}
	return nil
}

// GetActiveProxy 获取当前激活的代理（含解密密码，内部使用）
func (s *ProxyService) GetActiveProxy() (*models.Proxy, error) {
	row := s.db.QueryRow(
		`SELECT id, name, type, host, port, username, password, no_proxy, is_active, created_at, updated_at
		 FROM proxies WHERE is_active=1 LIMIT 1`,
	)
	var p models.Proxy
	if err := row.Scan(
		&p.Id, &p.Name, &p.Type, &p.Host, &p.Port,
		&p.Username, &p.Password, &p.NoProxy, &p.IsActive,
		&p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 没有激活的代理
		}
		return nil, fmt.Errorf("查询激活代理失败: %w", err)
	}
	// 解密密码
	if p.Password != "" {
		decrypted, err := crypto.Decrypt(p.Password)
		if err == nil {
			p.Password = decrypted
		}
	}
	return &p, nil
}

// GetDecryptedProxy 获取解密后的代理信息（按 ID，内部使用）
func (s *ProxyService) GetDecryptedProxy(id int64) (*models.Proxy, error) {
	p, err := s.GetById(id)
	if err != nil {
		return nil, err
	}
	if p.Password != "" {
		decrypted, err := crypto.Decrypt(p.Password)
		if err == nil {
			p.Password = decrypted
		}
	}
	return p, nil
}
