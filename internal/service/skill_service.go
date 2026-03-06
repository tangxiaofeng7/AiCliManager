package service

import (
	"AiCliManager/internal/db/models"
	"database/sql"
	"fmt"
	"time"
)

// SkillService 提供 Skill/Command 模板相关的业务逻辑
type SkillService struct {
	db *sql.DB
}

// NewSkillService 创建 SkillService 实例
func NewSkillService(db *sql.DB) *SkillService {
	return &SkillService{db: db}
}

// GetAll 获取所有 Skill 列表
func (s *SkillService) GetAll() ([]models.Skill, error) {
	rows, err := s.db.Query(
		`SELECT id, name, category, trigger, content, variables, is_builtin, sort_order, created_at, updated_at
		 FROM skills ORDER BY sort_order ASC, id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("查询 Skill 列表失败: %w", err)
	}
	defer rows.Close()

	var skills []models.Skill
	for rows.Next() {
		var sk models.Skill
		if err := rows.Scan(
			&sk.Id, &sk.Name, &sk.Category, &sk.Trigger, &sk.Content,
			&sk.Variables, &sk.IsBuiltin, &sk.SortOrder, &sk.CreatedAt, &sk.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描 Skill 行数据失败: %w", err)
		}
		skills = append(skills, sk)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历 Skill 行数据失败: %w", err)
	}
	return skills, nil
}

// GetById 按 ID 查询 Skill
func (s *SkillService) GetById(id int64) (*models.Skill, error) {
	row := s.db.QueryRow(
		`SELECT id, name, category, trigger, content, variables, is_builtin, sort_order, created_at, updated_at
		 FROM skills WHERE id=?`, id,
	)
	var sk models.Skill
	if err := row.Scan(
		&sk.Id, &sk.Name, &sk.Category, &sk.Trigger, &sk.Content,
		&sk.Variables, &sk.IsBuiltin, &sk.SortOrder, &sk.CreatedAt, &sk.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Skill (id=%d) 不存在", id)
		}
		return nil, fmt.Errorf("查询 Skill 失败: %w", err)
	}
	return &sk, nil
}

// Create 创建新 Skill
func (s *SkillService) Create(name, category, trigger, content, variables string, sortOrder int) (*models.Skill, error) {
	if variables == "" {
		variables = "[]"
	}
	now := time.Now()
	result, err := s.db.Exec(
		`INSERT INTO skills (name, category, trigger, content, variables, is_builtin, sort_order, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, 0, ?, ?, ?)`,
		name, category, trigger, content, variables, sortOrder, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("创建 Skill 失败: %w", err)
	}
	id, _ := result.LastInsertId()
	return &models.Skill{
		Id:        id,
		Name:      name,
		Category:  category,
		Trigger:   trigger,
		Content:   content,
		Variables: variables,
		IsBuiltin: 0,
		SortOrder: sortOrder,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update 更新 Skill（内置 Skill 不可修改）
func (s *SkillService) Update(id int64, name, category, trigger, content, variables string, sortOrder int) error {
	// 检查是否为内置 Skill
	sk, err := s.GetById(id)
	if err != nil {
		return err
	}
	if sk.IsBuiltin == 1 {
		return fmt.Errorf("内置 Skill 不可修改")
	}
	if variables == "" {
		variables = "[]"
	}
	now := time.Now()
	_, err = s.db.Exec(
		`UPDATE skills SET name=?, category=?, trigger=?, content=?, variables=?, sort_order=?, updated_at=?
		 WHERE id=?`,
		name, category, trigger, content, variables, sortOrder, now, id,
	)
	if err != nil {
		return fmt.Errorf("更新 Skill 失败: %w", err)
	}
	return nil
}

// Delete 删除 Skill（内置 Skill 不可删除）
func (s *SkillService) Delete(id int64) error {
	sk, err := s.GetById(id)
	if err != nil {
		return err
	}
	if sk.IsBuiltin == 1 {
		return fmt.Errorf("内置 Skill 不可删除")
	}
	_, err = s.db.Exec(`DELETE FROM skills WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("删除 Skill 失败: %w", err)
	}
	return nil
}

// GetByIds 按 ID 列表批量获取 Skill
func (s *SkillService) GetByIds(ids []int64) ([]models.Skill, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var skills []models.Skill
	for _, id := range ids {
		sk, err := s.GetById(id)
		if err != nil {
			return nil, err
		}
		skills = append(skills, *sk)
	}
	return skills, nil
}
