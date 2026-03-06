package service

import (
	"AiCliManager/internal/cli"
	"AiCliManager/internal/db/models"
	"encoding/json"
	"fmt"
	"strings"
)

// SyncRequest 配置同步请求参数
type SyncRequest struct {
	CliToolKey string
	CliTool    *models.CliTool
	Runtime    SyncRuntimeConfig
}

// SyncRuntimeConfig 表示启动时聚合后的运行态配置
type SyncRuntimeConfig struct {
	Profile    *models.Profile
	Provider   *models.Provider
	Proxy      *models.Proxy
	McpServers []models.McpServer
	Skills     []models.Skill
	SkillVars  map[string]string
}

// CliConfigWriter 定义 CLI 配置文件读写接口
type CliConfigWriter interface {
	ReadConfig(path string) (map[string]interface{}, error)
	WriteConfig(path string, config map[string]interface{}) error
	ConfigPath() string
}

type syncAdapter interface {
	Apply(config map[string]interface{}, runtime SyncRuntimeConfig)
}

// SyncService 负责将 AiCliManager 中的配置同步写入各 CLI 工具的配置文件
type SyncService struct{}

// NewSyncService 创建 SyncService 实例
func NewSyncService() *SyncService {
	return &SyncService{}
}

// SyncConfig 在启动 CLI 工具前调用，将配置写入对应工具的配置文件
func (s *SyncService) SyncConfig(req SyncRequest) error {
	writer, adapter, err := s.getSyncTarget(req.CliToolKey)
	if err != nil {
		return err
	}

	configPath := ""
	if req.CliTool != nil && req.CliTool.ConfigPath != "" {
		configPath = req.CliTool.ConfigPath
	}

	config, err := writer.ReadConfig(configPath)
	if err != nil {
		return fmt.Errorf("读取 CLI 工具配置文件失败: %w", err)
	}

	adapter.Apply(config, req.Runtime)

	if err := writer.WriteConfig(configPath, config); err != nil {
		return fmt.Errorf("写入 CLI 工具配置文件失败: %w", err)
	}

	return nil
}

func (s *SyncService) getSyncTarget(cliToolKey string) (CliConfigWriter, syncAdapter, error) {
	switch cliToolKey {
	case "claude":
		return cli.NewClaudeConfig(), claudeSyncAdapter{}, nil
	case "codex":
		return cli.NewCodexConfig(), codexSyncAdapter{}, nil
	case "opencode":
		return cli.NewOpenCodeConfig(), openCodeSyncAdapter{}, nil
	default:
		return nil, nil, fmt.Errorf("不支持的 CLI 工具: %s", cliToolKey)
	}
}

type claudeSyncAdapter struct{}

type codexSyncAdapter struct{}

type openCodeSyncAdapter struct{}

func (claudeSyncAdapter) Apply(config map[string]interface{}, runtime SyncRuntimeConfig) {
	applyProviderFields(config, runtime.Provider, "apiKey", "baseUrl")
	applyProfileFields(config, runtime.Profile, true)
	applyMcpServers(config, runtime.McpServers, true)
	applyCustomCommands(config, runtime.Skills, runtime.SkillVars)
}

func (codexSyncAdapter) Apply(config map[string]interface{}, runtime SyncRuntimeConfig) {
	applyProviderFields(config, runtime.Provider, "apiKey", "baseURL")
	applyProfileFields(config, runtime.Profile, false)
	applyMcpServers(config, runtime.McpServers, false)
}

func (openCodeSyncAdapter) Apply(config map[string]interface{}, runtime SyncRuntimeConfig) {
	applyProviderFields(config, runtime.Provider, "apiKey", "apiBaseURL")
	applyProfileFields(config, runtime.Profile, true)
	applyMcpServers(config, runtime.McpServers, false)
}

func applyProviderFields(config map[string]interface{}, provider *models.Provider, apiKeyField, baseURLField string) {
	if provider == nil {
		return
	}
	config[apiKeyField] = provider.ApiKey
	if provider.ApiUrl != "" {
		config[baseURLField] = provider.ApiUrl
	}
}

func applyProfileFields(config map[string]interface{}, profile *models.Profile, includeSystemPrompt bool) {
	if profile == nil {
		return
	}
	if profile.Model != "" {
		config["model"] = profile.Model
	}
	if includeSystemPrompt && profile.SystemPrompt != "" {
		config["systemPrompt"] = profile.SystemPrompt
	}
}

func applyMcpServers(config map[string]interface{}, servers []models.McpServer, includeEnv bool) {
	if len(servers) == 0 {
		return
	}
	mcpServers := map[string]interface{}{}
	for _, server := range servers {
		if server.IsEnabled == 0 {
			continue
		}
		entry := map[string]interface{}{
			"type": server.Type,
		}
		switch server.Type {
		case "stdio":
			entry["command"] = server.Command
			if argsArr := parseStringArray(server.Args); len(argsArr) > 0 {
				entry["args"] = argsArr
			}
			if includeEnv {
				if envMap := parseStringMap(server.Env); len(envMap) > 0 {
					entry["env"] = envMap
				}
			}
		case "sse", "http":
			entry["url"] = server.Url
		}
		mcpServers[server.Name] = entry
	}
	config["mcpServers"] = mcpServers
}

func applyCustomCommands(config map[string]interface{}, skills []models.Skill, skillVars map[string]string) {
	if len(skills) == 0 {
		return
	}
	commands := make([]map[string]interface{}, 0, len(skills))
	for _, sk := range skills {
		content := applySkillVars(sk.Content, sk.Id, skillVars)
		command := map[string]interface{}{
			"name":        sk.Name,
			"description": sk.Name,
			"prompt":      content,
		}
		if sk.Trigger != "" {
			command["trigger"] = sk.Trigger
		}
		commands = append(commands, command)
	}
	config["customCommands"] = commands
}

func parseStringArray(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var result []string
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil
	}
	return result
}

func parseStringMap(raw string) map[string]string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil
	}
	return result
}

// applySkillVars 将 Skill 内容中的 {{变量}} 占位符替换为实际值
// skillVars 的 key 格式为 "{skillId}.{varName}"
func applySkillVars(content string, skillId int64, skillVars map[string]string) string {
	if skillVars == nil {
		return content
	}
	prefix := fmt.Sprintf("%d.", skillId)
	for k, v := range skillVars {
		if strings.HasPrefix(k, prefix) {
			varName := strings.TrimPrefix(k, prefix)
			content = strings.ReplaceAll(content, "{{"+varName+"}}", v)
		}
	}
	return content
}
