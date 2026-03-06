package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// OpenCodeConfig 实现 OpenCode CLI 配置文件的读写
type OpenCodeConfig struct{}

// NewOpenCodeConfig 创建 OpenCodeConfig 实例
func NewOpenCodeConfig() *OpenCodeConfig {
	return &OpenCodeConfig{}
}

// ConfigPath 返回 OpenCode 配置文件的默认路径
// Windows: %APPDATA%\OpenCode\config.json
// macOS: ~/Library/Application Support/OpenCode/config.json
// Linux: ~/.config/opencode/config.json
func (c *OpenCodeConfig) ConfigPath() string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			homeDir, _ := os.UserHomeDir()
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		return filepath.Join(appData, "OpenCode", "config.json")
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, "Library", "Application Support", "OpenCode", "config.json")
	default:
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".config", "opencode", "config.json")
	}
}

// ReadConfig 读取 OpenCode 配置文件，返回 map
func (c *OpenCodeConfig) ReadConfig(path string) (map[string]interface{}, error) {
	if path == "" {
		path = c.ConfigPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("读取 OpenCode 配置文件失败 [%s]: %w", path, err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析 OpenCode 配置文件失败: %w", err)
	}
	return config, nil
}

// WriteConfig 将配置写入 OpenCode 配置文件，写入前备份原文件
func (c *OpenCodeConfig) WriteConfig(path string, config map[string]interface{}) error {
	if path == "" {
		path = c.ConfigPath()
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("创建 OpenCode 配置目录失败 [%s]: %w", dir, err)
	}

	backupPath := path + ".bak"
	if err := backupFile(path, backupPath); err != nil {
		return fmt.Errorf("备份 OpenCode 配置文件失败: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 OpenCode 配置失败: %w", err)
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0640); err != nil {
		return fmt.Errorf("写入 OpenCode 临时配置文件失败: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = restoreBackup(backupPath, path)
		return fmt.Errorf("替换 OpenCode 配置文件失败: %w", err)
	}

	return nil
}
