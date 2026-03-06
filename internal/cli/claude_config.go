package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ClaudeConfig 实现 Claude Code 配置文件的读写
type ClaudeConfig struct{}

// NewClaudeConfig 创建 ClaudeConfig 实例
func NewClaudeConfig() *ClaudeConfig {
	return &ClaudeConfig{}
}

// ConfigPath 返回 Claude Code 配置文件的默认路径
// Windows: %USERPROFILE%\.claude\settings.json
// macOS/Linux: ~/.claude/settings.json
func (c *ClaudeConfig) ConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".claude", "settings.json")
}

// ReadConfig 读取 Claude Code 配置文件，返回 map
func (c *ClaudeConfig) ReadConfig(path string) (map[string]interface{}, error) {
	if path == "" {
		path = c.ConfigPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在时返回空配置
			return map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("读取 Claude Code 配置文件失败 [%s]: %w", path, err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析 Claude Code 配置文件失败: %w", err)
	}
	return config, nil
}

// WriteConfig 将配置写入 Claude Code 配置文件，写入前备份原文件
func (c *ClaudeConfig) WriteConfig(path string, config map[string]interface{}) error {
	if path == "" {
		path = c.ConfigPath()
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("创建配置目录失败 [%s]: %w", dir, err)
	}

	// 写入前备份原文件
	backupPath := path + ".bak"
	if err := c.backupFile(path, backupPath); err != nil {
		return fmt.Errorf("备份配置文件失败: %w", err)
	}

	// 序列化配置（带缩进，易读）
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 先写入临时文件，成功后再替换原文件（原子性保证）
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0640); err != nil {
		return fmt.Errorf("写入临时配置文件失败: %w", err)
	}

	// 将临时文件替换原文件
	if err := os.Rename(tmpPath, path); err != nil {
		// 替换失败时，尝试从备份恢复
		_ = c.restoreBackup(backupPath, path)
		return fmt.Errorf("替换配置文件失败: %w", err)
	}

	return nil
}

// backupFile 将源文件备份到目标路径（源文件不存在时跳过）
func (c *ClaudeConfig) backupFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 源文件不存在，无需备份
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

// restoreBackup 从备份文件恢复（写入失败时调用）
func (c *ClaudeConfig) restoreBackup(backupPath, targetPath string) error {
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return nil // 备份不存在，无法恢复
	}
	return os.Rename(backupPath, targetPath)
}
