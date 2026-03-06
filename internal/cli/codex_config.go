package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// CodexConfig 实现 Codex CLI 配置文件的读写
type CodexConfig struct{}

// NewCodexConfig 创建 CodexConfig 实例
func NewCodexConfig() *CodexConfig {
	return &CodexConfig{}
}

// ConfigPath 返回 Codex 配置文件的默认路径
// Windows: %APPDATA%\Codex\config.json
// macOS: ~/Library/Application Support/Codex/config.json
// Linux: ~/.config/codex/config.json
func (c *CodexConfig) ConfigPath() string {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			homeDir, _ := os.UserHomeDir()
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		return filepath.Join(appData, "Codex", "config.json")
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, "Library", "Application Support", "Codex", "config.json")
	default:
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".config", "codex", "config.json")
	}
}

// ReadConfig 读取 Codex 配置文件，返回 map
func (c *CodexConfig) ReadConfig(path string) (map[string]interface{}, error) {
	if path == "" {
		path = c.ConfigPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("读取 Codex 配置文件失败 [%s]: %w", path, err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析 Codex 配置文件失败: %w", err)
	}
	return config, nil
}

// WriteConfig 将配置写入 Codex 配置文件，写入前备份原文件
func (c *CodexConfig) WriteConfig(path string, config map[string]interface{}) error {
	if path == "" {
		path = c.ConfigPath()
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("创建 Codex 配置目录失败 [%s]: %w", dir, err)
	}

	backupPath := path + ".bak"
	if err := backupFile(path, backupPath); err != nil {
		return fmt.Errorf("备份 Codex 配置文件失败: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化 Codex 配置失败: %w", err)
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0640); err != nil {
		return fmt.Errorf("写入 Codex 临时配置文件失败: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = restoreBackup(backupPath, path)
		return fmt.Errorf("替换 Codex 配置文件失败: %w", err)
	}

	return nil
}

// backupFile 将源文件备份到目标路径
func backupFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// restoreBackup 从备份文件恢复
func restoreBackup(backupPath, targetPath string) error {
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return nil
	}
	return os.Rename(backupPath, targetPath)
}
