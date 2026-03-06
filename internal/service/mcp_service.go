package service

import (
	"AiCliManager/internal/db/models"
	"database/sql"
	"fmt"
	"time"
)

// McpService 提供 MCP Server 相关的业务逻辑
type McpService struct {
	db *sql.DB
}

// NewMcpService 创建 McpService 实例
func NewMcpService(db *sql.DB) *McpService {
	return &McpService{db: db}
}

// GetAll 获取所有 MCP Server 列表
func (s *McpService) GetAll() ([]models.McpServer, error) {
	rows, err := s.db.Query(
		`SELECT id, name, type, command, args, env, url, description, is_enabled, sort_order, created_at, updated_at
		 FROM mcp_servers ORDER BY sort_order ASC, id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("查询 MCP Server 列表失败: %w", err)
	}
	defer rows.Close()

	var servers []models.McpServer
	for rows.Next() {
		var m models.McpServer
		if err := rows.Scan(
			&m.Id, &m.Name, &m.Type, &m.Command, &m.Args, &m.Env,
			&m.Url, &m.Description, &m.IsEnabled, &m.SortOrder,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描 MCP Server 行数据失败: %w", err)
		}
		servers = append(servers, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历 MCP Server 行数据失败: %w", err)
	}
	return servers, nil
}

// GetById 按 ID 查询 MCP Server
func (s *McpService) GetById(id int64) (*models.McpServer, error) {
	row := s.db.QueryRow(
		`SELECT id, name, type, command, args, env, url, description, is_enabled, sort_order, created_at, updated_at
		 FROM mcp_servers WHERE id=?`, id,
	)
	var m models.McpServer
	if err := row.Scan(
		&m.Id, &m.Name, &m.Type, &m.Command, &m.Args, &m.Env,
		&m.Url, &m.Description, &m.IsEnabled, &m.SortOrder,
		&m.CreatedAt, &m.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("MCP Server (id=%d) 不存在", id)
		}
		return nil, fmt.Errorf("查询 MCP Server 失败: %w", err)
	}
	return &m, nil
}

// Create 创建新 MCP Server
func (s *McpService) Create(name, mcpType, command, args, env, url, description string, isEnabled, sortOrder int) (*models.McpServer, error) {
	if args == "" {
		args = "[]"
	}
	if env == "" {
		env = "{}"
	}
	now := time.Now()
	result, err := s.db.Exec(
		`INSERT INTO mcp_servers (name, type, command, args, env, url, description, is_enabled, sort_order, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		name, mcpType, command, args, env, url, description, isEnabled, sortOrder, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("创建 MCP Server 失败: %w", err)
	}
	id, _ := result.LastInsertId()
	return &models.McpServer{
		Id:          id,
		Name:        name,
		Type:        mcpType,
		Command:     command,
		Args:        args,
		Env:         env,
		Url:         url,
		Description: description,
		IsEnabled:   isEnabled,
		SortOrder:   sortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update 更新 MCP Server
func (s *McpService) Update(id int64, name, mcpType, command, args, env, url, description string, isEnabled, sortOrder int) error {
	if args == "" {
		args = "[]"
	}
	if env == "" {
		env = "{}"
	}
	now := time.Now()
	_, err := s.db.Exec(
		`UPDATE mcp_servers SET name=?, type=?, command=?, args=?, env=?, url=?, description=?, is_enabled=?, sort_order=?, updated_at=?
		 WHERE id=?`,
		name, mcpType, command, args, env, url, description, isEnabled, sortOrder, now, id,
	)
	if err != nil {
		return fmt.Errorf("更新 MCP Server 失败: %w", err)
	}
	return nil
}

// Delete 删除 MCP Server
func (s *McpService) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM mcp_servers WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("删除 MCP Server 失败: %w", err)
	}
	return nil
}

// GetCliToolMcpServers 获取指定 CLI 工具关联的 MCP Server 列表
func (s *McpService) GetCliToolMcpServers(cliToolKey string) ([]models.McpServer, error) {
	rows, err := s.db.Query(
		`SELECT m.id, m.name, m.type, m.command, m.args, m.env, m.url, m.description, m.is_enabled, m.sort_order, m.created_at, m.updated_at
		 FROM mcp_servers m
		 INNER JOIN cli_tool_mcp_servers ctm ON ctm.mcp_server_id = m.id
		 INNER JOIN cli_tools ct ON ct.id = ctm.cli_tool_id
		 WHERE ct.key=?
		 ORDER BY m.sort_order ASC, m.id ASC`, cliToolKey,
	)
	if err != nil {
		return nil, fmt.Errorf("查询 CLI 工具 MCP Server 列表失败: %w", err)
	}
	defer rows.Close()

	var servers []models.McpServer
	for rows.Next() {
		var m models.McpServer
		if err := rows.Scan(
			&m.Id, &m.Name, &m.Type, &m.Command, &m.Args, &m.Env,
			&m.Url, &m.Description, &m.IsEnabled, &m.SortOrder,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描 MCP Server 行数据失败: %w", err)
		}
		servers = append(servers, m)
	}
	return servers, nil
}

// SetCliToolMcpServers 设置指定 CLI 工具关联的 MCP Server 列表（全量替换）
func (s *McpService) SetCliToolMcpServers(cliToolKey string, ids []int64) error {
	// 先查询 CLI 工具 ID
	row := s.db.QueryRow(`SELECT id FROM cli_tools WHERE key=?`, cliToolKey)
	var cliToolId int64
	if err := row.Scan(&cliToolId); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("CLI 工具 [%s] 不存在", cliToolKey)
		}
		return fmt.Errorf("查询 CLI 工具失败: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback()

	// 删除旧关联
	if _, err := tx.Exec(`DELETE FROM cli_tool_mcp_servers WHERE cli_tool_id=?`, cliToolId); err != nil {
		return fmt.Errorf("删除旧 MCP Server 关联失败: %w", err)
	}

	// 插入新关联
	for _, mcpId := range ids {
		if _, err := tx.Exec(
			`INSERT INTO cli_tool_mcp_servers (cli_tool_id, mcp_server_id) VALUES (?, ?)`,
			cliToolId, mcpId,
		); err != nil {
			return fmt.Errorf("插入 MCP Server 关联失败: %w", err)
		}
	}

	return tx.Commit()
}

// GetByIds 按 ID 列表批量获取 MCP Server
func (s *McpService) GetByIds(ids []int64) ([]models.McpServer, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var servers []models.McpServer
	for _, id := range ids {
		m, err := s.GetById(id)
		if err != nil {
			return nil, err
		}
		servers = append(servers, *m)
	}
	return servers, nil
}
