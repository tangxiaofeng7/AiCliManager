package service

import (
	"AiCliManager/internal/db/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// SessionService 提供会话历史相关的业务逻辑
type SessionService struct {
	db *sql.DB
}

// NewSessionService 创建 SessionService 实例
func NewSessionService(db *sql.DB) *SessionService {
	return &SessionService{db: db}
}

// GetSessionsRequest 查询会话列表的请求参数
type GetSessionsRequest struct {
	CliToolKey string `json:"cli_tool_key"` // 按工具筛选，空=全部
	Page       int    `json:"page"`         // 页码，从 1 开始
	PageSize   int    `json:"page_size"`    // 每页条数，默认 20
}

// GetAll 获取会话历史列表，支持分页和工具筛选
func (s *SessionService) GetAll(req GetSessionsRequest) ([]models.Session, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	offset := (req.Page - 1) * req.PageSize

	var rows *sql.Rows
	var err error
	if req.CliToolKey != "" {
		rows, err = s.db.Query(
			`SELECT s.id, s.cli_tool_id, s.profile_id, s.proxy_id, s.terminal, s.working_dir,
			        s.extra_args, s.status, s.pid, s.started_at, s.ended_at
			 FROM sessions s
			 INNER JOIN cli_tools ct ON ct.id = s.cli_tool_id
			 WHERE ct.key=?
			 ORDER BY s.started_at DESC
			 LIMIT ? OFFSET ?`,
			req.CliToolKey, req.PageSize, offset,
		)
	} else {
		rows, err = s.db.Query(
			`SELECT id, cli_tool_id, profile_id, proxy_id, terminal, working_dir,
			        extra_args, status, pid, started_at, ended_at
			 FROM sessions
			 ORDER BY started_at DESC
			 LIMIT ? OFFSET ?`,
			req.PageSize, offset,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("查询会话列表失败: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var sess models.Session
		if err := rows.Scan(
			&sess.Id, &sess.CliToolId, &sess.ProfileId, &sess.ProxyId,
			&sess.Terminal, &sess.WorkingDir, &sess.ExtraArgs,
			&sess.Status, &sess.Pid, &sess.StartedAt, &sess.EndedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描会话行数据失败: %w", err)
		}
		sessions = append(sessions, sess)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历会话行数据失败: %w", err)
	}
	return sessions, nil
}

// GetById 按 ID 查询会话
func (s *SessionService) GetById(id int64) (*models.Session, error) {
	row := s.db.QueryRow(
		`SELECT id, cli_tool_id, profile_id, proxy_id, terminal, working_dir,
		        extra_args, status, pid, started_at, ended_at
		 FROM sessions WHERE id=?`, id,
	)
	var sess models.Session
	if err := row.Scan(
		&sess.Id, &sess.CliToolId, &sess.ProfileId, &sess.ProxyId,
		&sess.Terminal, &sess.WorkingDir, &sess.ExtraArgs,
		&sess.Status, &sess.Pid, &sess.StartedAt, &sess.EndedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("会话 (id=%d) 不存在", id)
		}
		return nil, fmt.Errorf("查询会话失败: %w", err)
	}
	return &sess, nil
}

// Create 创建新会话记录
func (s *SessionService) Create(cliToolId int64, profileId *int64, proxyId *int64, terminal, workingDir string, extraArgs []string, pid int) (*models.Session, error) {
	argsJSON, err := json.Marshal(extraArgs)
	if err != nil {
		argsJSON = []byte("[]")
	}
	now := time.Now()
	result, err := s.db.Exec(
		`INSERT INTO sessions (cli_tool_id, profile_id, proxy_id, terminal, working_dir, extra_args, status, pid, started_at)
		 VALUES (?, ?, ?, ?, ?, ?, 'running', ?, ?)`,
		cliToolId, profileId, proxyId, terminal, workingDir, string(argsJSON), pid, now,
	)
	if err != nil {
		return nil, fmt.Errorf("创建会话记录失败: %w", err)
	}
	id, _ := result.LastInsertId()
	return &models.Session{
		Id:         id,
		CliToolId:  cliToolId,
		ProfileId:  profileId,
		ProxyId:    proxyId,
		Terminal:   terminal,
		WorkingDir: workingDir,
		ExtraArgs:  string(argsJSON),
		Status:     "running",
		Pid:        pid,
		StartedAt:  now,
	}, nil
}

// UpdateStatus 更新会话状态（进程退出后调用）
func (s *SessionService) UpdateStatus(id int64, status string) error {
	now := time.Now()
	_, err := s.db.Exec(
		`UPDATE sessions SET status=?, ended_at=? WHERE id=?`,
		status, now, id,
	)
	if err != nil {
		return fmt.Errorf("更新会话状态失败: %w", err)
	}
	return nil
}

// Delete 删除单条会话记录
func (s *SessionService) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM sessions WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("删除会话失败: %w", err)
	}
	return nil
}

// ClearByCliTool 清空指定 CLI 工具的历史（cliToolKey 为空则清空全部）
func (s *SessionService) ClearByCliTool(cliToolKey string) error {
	var err error
	if cliToolKey == "" {
		_, err = s.db.Exec(`DELETE FROM sessions`)
	} else {
		_, err = s.db.Exec(
			`DELETE FROM sessions WHERE cli_tool_id IN (SELECT id FROM cli_tools WHERE key=?)`,
			cliToolKey,
		)
	}
	if err != nil {
		return fmt.Errorf("清空会话历史失败: %w", err)
	}
	return nil
}

// AutoCleanup 按最大条数自动清理旧记录
func (s *SessionService) AutoCleanup(maxCount int) error {
	if maxCount <= 0 {
		return nil
	}
	_, err := s.db.Exec(
		`DELETE FROM sessions WHERE id NOT IN (
			SELECT id FROM sessions ORDER BY started_at DESC LIMIT ?
		)`, maxCount,
	)
	if err != nil {
		return fmt.Errorf("自动清理会话记录失败: %w", err)
	}
	return nil
}
