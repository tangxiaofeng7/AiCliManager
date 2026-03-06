package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // 注册 SQLite 驱动（纯 Go，无 CGO）
)

// DB 是全局数据库连接实例
var DB *sql.DB

// Init 初始化 SQLite 数据库，返回连接实例
func Init() *sql.DB {
	// 确定数据库文件路径：用户 Home 目录下 .aiclimgr/data.db
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("获取用户主目录失败: %v", err))
	}
	dbDir := filepath.Join(homeDir, ".aiclimgr")
	if err := os.MkdirAll(dbDir, 0750); err != nil {
		panic(fmt.Sprintf("创建数据库目录失败: %v", err))
	}
	dbPath := filepath.Join(dbDir, "data.db")

	// 打开数据库连接
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(fmt.Sprintf("打开数据库失败: %v", err))
	}

	// 设置连接池参数（SQLite 单文件，建议只用单连接）
	db.SetMaxOpenConns(1)

	// 设置 PRAGMA
	pragmas := []string{
		"PRAGMA journal_mode=WAL",  // WAL 模式提升并发性能
		"PRAGMA foreign_keys=ON",   // 开启外键约束
		"PRAGMA synchronous=NORMAL", // 均衡性能与安全
		"PRAGMA busy_timeout=5000", // 等待锁超时 5 秒
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			panic(fmt.Sprintf("设置 PRAGMA 失败 [%s]: %v", p, err))
		}
	}

	// 建表
	if err := createTables(db); err != nil {
		panic(fmt.Sprintf("创建数据库表失败: %v", err))
	}

	// 插入初始数据
	if err := seedInitialData(db); err != nil {
		panic(fmt.Sprintf("插入初始数据失败: %v", err))
	}

	DB = db
	return DB
}

// createTables 创建所有数据库表（如果不存在）
func createTables(db *sql.DB) error {
	statements := []string{
		// CLI 工具注册表
		`CREATE TABLE IF NOT EXISTS cli_tools (
			id                 INTEGER PRIMARY KEY AUTOINCREMENT,
			name               TEXT NOT NULL,
			key                TEXT NOT NULL UNIQUE,
			executable         TEXT DEFAULT '',
			config_path        TEXT DEFAULT '',
			preferred_terminal TEXT DEFAULT '',
			is_installed       INTEGER DEFAULT 0,
			is_enabled         INTEGER DEFAULT 1,
			sort_order         INTEGER DEFAULT 0,
			created_at         DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at         DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// API 提供商表
		`CREATE TABLE IF NOT EXISTS providers (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT NOT NULL,
			type        TEXT NOT NULL,
			api_url     TEXT NOT NULL,
			api_key     TEXT NOT NULL DEFAULT '',
			models      TEXT DEFAULT '[]',
			sort_order  INTEGER DEFAULT 0,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// 配置 Profile 表
		`CREATE TABLE IF NOT EXISTS profiles (
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			name          TEXT NOT NULL,
			provider_id   INTEGER NOT NULL,
			model         TEXT NOT NULL DEFAULT '',
			system_prompt TEXT DEFAULT '',
			temperature   REAL DEFAULT 1.0,
			max_tokens    INTEGER DEFAULT 8192,
			extra_config  TEXT DEFAULT '{}',
			created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (provider_id) REFERENCES providers(id)
		)`,

		// 代理配置表
		`CREATE TABLE IF NOT EXISTS proxies (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT NOT NULL,
			type        TEXT NOT NULL,
			host        TEXT NOT NULL,
			port        INTEGER NOT NULL,
			username    TEXT DEFAULT '',
			password    TEXT DEFAULT '',
			no_proxy    TEXT DEFAULT '',
			is_active   INTEGER DEFAULT 0,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// MCP Server 配置表
		`CREATE TABLE IF NOT EXISTS mcp_servers (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT NOT NULL,
			type        TEXT NOT NULL,
			command     TEXT DEFAULT '',
			args        TEXT DEFAULT '[]',
			env         TEXT DEFAULT '{}',
			url         TEXT DEFAULT '',
			description TEXT DEFAULT '',
			is_enabled  INTEGER DEFAULT 1,
			sort_order  INTEGER DEFAULT 0,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// CLI 工具与 MCP Server 的关联表（多对多）
		`CREATE TABLE IF NOT EXISTS cli_tool_mcp_servers (
			cli_tool_id   INTEGER NOT NULL,
			mcp_server_id INTEGER NOT NULL,
			PRIMARY KEY (cli_tool_id, mcp_server_id),
			FOREIGN KEY (cli_tool_id)   REFERENCES cli_tools(id) ON DELETE CASCADE,
			FOREIGN KEY (mcp_server_id) REFERENCES mcp_servers(id) ON DELETE CASCADE
		)`,

		// Skills / Commands 模板表
		`CREATE TABLE IF NOT EXISTS skills (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT NOT NULL,
			category    TEXT DEFAULT 'general',
			trigger     TEXT DEFAULT '',
			content     TEXT NOT NULL,
			variables   TEXT DEFAULT '[]',
			is_builtin  INTEGER DEFAULT 0,
			sort_order  INTEGER DEFAULT 0,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// CLI 工具当前激活配置表
		`CREATE TABLE IF NOT EXISTS cli_tool_active_config (
			cli_tool_id INTEGER PRIMARY KEY,
			profile_id  INTEGER,
			proxy_id    INTEGER,
			updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (cli_tool_id) REFERENCES cli_tools(id) ON DELETE CASCADE,
			FOREIGN KEY (profile_id)  REFERENCES profiles(id),
			FOREIGN KEY (proxy_id)    REFERENCES proxies(id)
		)`,

		// 全局设置表
		`CREATE TABLE IF NOT EXISTS settings (
			key         TEXT PRIMARY KEY,
			value       TEXT NOT NULL,
			description TEXT DEFAULT '',
			updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// 会话历史表
		`CREATE TABLE IF NOT EXISTS sessions (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			cli_tool_id INTEGER NOT NULL,
			profile_id  INTEGER,
			proxy_id    INTEGER,
			terminal    TEXT DEFAULT '',
			working_dir TEXT DEFAULT '',
			extra_args  TEXT DEFAULT '[]',
			status      TEXT DEFAULT 'running',
			pid         INTEGER DEFAULT 0,
			started_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
			ended_at    DATETIME,
			FOREIGN KEY (cli_tool_id) REFERENCES cli_tools(id),
			FOREIGN KEY (profile_id)  REFERENCES profiles(id),
			FOREIGN KEY (proxy_id)    REFERENCES proxies(id)
		)`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("执行建表语句失败: %w", err)
		}
	}
	return nil
}

// seedInitialData 插入初始的 CLI 工具数据（幂等操作，重复插入时忽略）
func seedInitialData(db *sql.DB) error {
	// 插入三个支持的 CLI 工具（使用 INSERT OR IGNORE 保证幂等）
	tools := []struct {
		name      string
		key       string
		sortOrder int
	}{
		{"Claude Code", "claude", 1},
		{"Codex", "codex", 2},
		{"OpenCode", "opencode", 3},
	}

	for _, t := range tools {
		_, err := db.Exec(
			`INSERT OR IGNORE INTO cli_tools (name, key, sort_order) VALUES (?, ?, ?)`,
			t.name, t.key, t.sortOrder,
		)
		if err != nil {
			return fmt.Errorf("插入 CLI 工具 [%s] 失败: %w", t.key, err)
		}
	}

	// 插入默认全局设置
	defaultSettings := []struct {
		key         string
		value       string
		description string
	}{
		{"theme", "system", "界面主题：light | dark | system"},
		{"default_terminal", "default", "默认启动终端"},
		{"session_max_count", "100", "会话历史最大保留条数"},
		{"auto_detect_tools", "true", "启动时自动检测 CLI 工具安装状态"},
	}
	for _, s := range defaultSettings {
		_, err := db.Exec(
			`INSERT OR IGNORE INTO settings (key, value, description) VALUES (?, ?, ?)`,
			s.key, s.value, s.description,
		)
		if err != nil {
			return fmt.Errorf("插入默认设置 [%s] 失败: %w", s.key, err)
		}
	}

	return nil
}
