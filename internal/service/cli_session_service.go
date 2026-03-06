package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// CliSession 表示 AI CLI 工具的一个对话会话（从 CLI 工具本地会话文件读取）
type CliSession struct {
	SessionId      string `json:"session_id"`
	CliToolKey     string `json:"cli_tool_key"`
	Project        string `json:"project"`
	ProjectDir     string `json:"project_dir"`
	Slug           string `json:"slug"`
	FirstMessage   string `json:"first_message"`
	MessageCount   int    `json:"message_count"`
	UserCount      int    `json:"user_count"`
	AssistantCount int    `json:"assistant_count"`
	Model          string `json:"model"`
	StartedAt      string `json:"started_at"`
	LastActiveAt   string `json:"last_active_at"`
}

// CliSessionMessage 表示一条对话消息
type CliSessionMessage struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	Model     string `json:"model,omitempty"`
	TokensIn  int    `json:"tokens_in,omitempty"`
	TokensOut int    `json:"tokens_out,omitempty"`
	UUID      string `json:"uuid"`
}

// CliSessionProject 表示一个项目目录及其会话数量
type CliSessionProject struct {
	DirName      string `json:"dir_name"`
	Path         string `json:"path"`
	SessionCount int    `json:"session_count"`
}

// GetCliSessionsRequest 查询 CLI 会话的请求参数
type GetCliSessionsRequest struct {
	CliToolKey string `json:"cli_tool_key"`
	Project    string `json:"project"`
	Limit      int    `json:"limit"`
}

// CliSessionService 负责读取 AI CLI 工具的实际对话会话文件
type CliSessionService struct{}

// NewCliSessionService 创建实例
func NewCliSessionService() *CliSessionService {
	return &CliSessionService{}
}

type sessionFileInfo struct {
	path    string
	modTime time.Time
	projDir string
}

// GetSessions 获取 CLI 工具的对话会话列表
func (s *CliSessionService) GetSessions(req GetCliSessionsRequest) ([]CliSession, error) {
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.CliToolKey == "" {
		req.CliToolKey = "claude"
	}

	baseDir, err := s.getSessionsBaseDir(req.CliToolKey)
	if err != nil {
		return nil, err
	}

	files, err := s.collectSessionFiles(baseDir, req.Project)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.After(files[j].modTime)
	})
	if len(files) > req.Limit {
		files = files[:req.Limit]
	}

	sessions := make([]CliSession, 0, len(files))
	for _, f := range files {
		sess, err := s.parseSessionMetadata(f.path, f.projDir, req.CliToolKey)
		if err != nil || sess == nil {
			continue
		}
		if sess.UserCount == 0 {
			continue
		}
		sessions = append(sessions, *sess)
	}

	return sessions, nil
}

// GetMessages 获取指定会话的完整消息列表
func (s *CliSessionService) GetMessages(cliToolKey, sessionId string) ([]CliSessionMessage, error) {
	if cliToolKey == "" {
		cliToolKey = "claude"
	}

	filePath, err := s.findSessionFile(cliToolKey, sessionId)
	if err != nil {
		return nil, err
	}

	return s.parseSessionMessages(filePath)
}

// GetProjects 获取指定 CLI 工具的所有项目目录列表
func (s *CliSessionService) GetProjects(cliToolKey string) ([]CliSessionProject, error) {
	if cliToolKey == "" {
		cliToolKey = "claude"
	}

	baseDir, err := s.getSessionsBaseDir(cliToolKey)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("读取项目目录失败: %w", err)
	}

	projects := make([]CliSessionProject, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirName := entry.Name()
		projDir := filepath.Join(baseDir, dirName)
		jsonlFiles, _ := filepath.Glob(filepath.Join(projDir, "*.jsonl"))
		if len(jsonlFiles) == 0 {
			continue
		}

		projPath := s.extractProjectPath(jsonlFiles[0])
		if projPath == "" {
			projPath = dirName
		}

		projects = append(projects, CliSessionProject{
			DirName:      dirName,
			Path:         projPath,
			SessionCount: len(jsonlFiles),
		})
	}

	sort.Slice(projects, func(i, j int) bool {
		return strings.ToLower(projects[i].Path) < strings.ToLower(projects[j].Path)
	})

	return projects, nil
}

func (s *CliSessionService) getSessionsBaseDir(cliToolKey string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户目录失败: %w", err)
	}

	switch cliToolKey {
	case "claude":
		return filepath.Join(homeDir, ".claude", "projects"), nil
	case "codex":
		return "", fmt.Errorf("Codex 暂不支持会话历史读取")
	case "opencode":
		return "", fmt.Errorf("OpenCode 暂不支持会话历史读取")
	default:
		return "", fmt.Errorf("未知 CLI 工具: %s", cliToolKey)
	}
}

func (s *CliSessionService) collectSessionFiles(baseDir, projectFilter string) ([]sessionFileInfo, error) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("读取项目目录失败: %w", err)
	}

	var files []sessionFileInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirName := entry.Name()
		if projectFilter != "" && dirName != projectFilter {
			continue
		}

		projDir := filepath.Join(baseDir, dirName)
		jsonlFiles, err := os.ReadDir(projDir)
		if err != nil {
			continue
		}

		for _, f := range jsonlFiles {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".jsonl") {
				continue
			}
			info, err := f.Info()
			if err != nil {
				continue
			}
			files = append(files, sessionFileInfo{
				path:    filepath.Join(projDir, f.Name()),
				modTime: info.ModTime(),
				projDir: dirName,
			})
		}
	}

	return files, nil
}

func (s *CliSessionService) findSessionFile(cliToolKey, sessionId string) (string, error) {
	baseDir, err := s.getSessionsBaseDir(cliToolKey)
	if err != nil {
		return "", err
	}
	pattern := filepath.Join(baseDir, "*", sessionId+".jsonl")
	matches, _ := filepath.Glob(pattern)
	if len(matches) == 0 {
		return "", fmt.Errorf("会话文件未找到: %s", sessionId)
	}
	return matches[0], nil
}

type rawSessionLine struct {
	Type      string          `json:"type"`
	SessionId string          `json:"sessionId"`
	Slug      string          `json:"slug"`
	Cwd       string          `json:"cwd"`
	Timestamp json.RawMessage `json:"timestamp"`
	Message   *rawMessage     `json:"message"`
	Content   string          `json:"content"`
	UUID      string          `json:"uuid"`
	IsMeta    bool            `json:"isMeta"`
	UserType  string          `json:"userType"`
}

type rawMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
	Model   string          `json:"model"`
	Id      string          `json:"id"`
	Usage   *rawUsage       `json:"usage"`
}

type rawUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func (s *CliSessionService) parseSessionMetadata(filePath, projDir, cliToolKey string) (*CliSession, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sess := &CliSession{
		CliToolKey: cliToolKey,
		ProjectDir: projDir,
		SessionId:  strings.TrimSuffix(filepath.Base(filePath), ".jsonl"),
	}

	scanner := newJSONLScanner(file)
	firstUserFound := false
	firstAssistantFound := false

	for scanner.Scan() {
		var raw rawSessionLine
		if err := json.Unmarshal(scanner.Bytes(), &raw); err != nil {
			continue
		}
		if shouldSkipSessionLine(raw.Type) {
			continue
		}

		if sess.SessionId == "" && raw.SessionId != "" {
			sess.SessionId = raw.SessionId
		}
		if sess.Slug == "" && raw.Slug != "" {
			sess.Slug = raw.Slug
		}
		if sess.Project == "" && raw.Cwd != "" {
			sess.Project = raw.Cwd
		}

		ts := s.parseTimestamp(raw.Timestamp)
		switch raw.Type {
		case "user":
			if raw.IsMeta {
				continue
			}
			sess.UserCount++
			sess.MessageCount++
			if ts != "" {
				sess.LastActiveAt = ts
			}
			if !firstUserFound {
				firstUserFound = true
				if ts != "" {
					sess.StartedAt = ts
				}
				if raw.Message != nil {
					content := extractDisplayContent(raw.Message.Content)
					if len(content) > 200 {
						content = content[:200] + "..."
					}
					sess.FirstMessage = content
				}
			}
		case "assistant":
			sess.AssistantCount++
			sess.MessageCount++
			if ts != "" {
				sess.LastActiveAt = ts
			}
			if !firstAssistantFound && raw.Message != nil && raw.Message.Model != "" {
				firstAssistantFound = true
				sess.Model = raw.Message.Model
			}
		case "system":
			sess.MessageCount++
			if ts != "" {
				sess.LastActiveAt = ts
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取会话文件失败: %w", err)
	}
	return sess, nil
}

func (s *CliSessionService) parseSessionMessages(filePath string) ([]CliSessionMessage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开会话文件失败: %w", err)
	}
	defer file.Close()

	var messages []CliSessionMessage
	scanner := newJSONLScanner(file)
	for scanner.Scan() {
		var raw rawSessionLine
		if err := json.Unmarshal(scanner.Bytes(), &raw); err != nil {
			continue
		}
		if shouldSkipSessionLine(raw.Type) {
			continue
		}

		ts := s.parseTimestamp(raw.Timestamp)
		switch raw.Type {
		case "user":
			if raw.IsMeta {
				continue
			}
			content := ""
			if raw.Message != nil {
				content = extractDisplayContent(raw.Message.Content)
			}
			if content == "" {
				continue
			}
			messages = append(messages, CliSessionMessage{Type: "user", Content: content, Timestamp: ts, UUID: raw.UUID})
		case "assistant":
			if raw.Message == nil {
				continue
			}
			content := extractAssistantContent(raw.Message.Content)
			if content == "" {
				continue
			}
			msg := CliSessionMessage{
				Type:      "assistant",
				Content:   content,
				Timestamp: ts,
				Model:     raw.Message.Model,
				UUID:      raw.UUID,
			}
			if raw.Message.Usage != nil {
				msg.TokensIn = raw.Message.Usage.InputTokens
				msg.TokensOut = raw.Message.Usage.OutputTokens
			}
			messages = append(messages, msg)
		case "system":
			if raw.Content == "" {
				continue
			}
			messages = append(messages, CliSessionMessage{Type: "system", Content: raw.Content, Timestamp: ts, UUID: raw.UUID})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取会话消息失败: %w", err)
	}
	return messages, nil
}

func (s *CliSessionService) extractProjectPath(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := newJSONLScanner(file)
	for scanner.Scan() {
		var raw rawSessionLine
		if json.Unmarshal(scanner.Bytes(), &raw) == nil && raw.Cwd != "" {
			return raw.Cwd
		}
	}
	return ""
}

func (s *CliSessionService) parseTimestamp(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var str string
	if json.Unmarshal(raw, &str) == nil && str != "" {
		return str
	}
	var ms float64
	if json.Unmarshal(raw, &ms) == nil && ms > 0 {
		return time.UnixMilli(int64(ms)).UTC().Format(time.RFC3339Nano)
	}
	return ""
}

func newJSONLScanner(file *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)
	return scanner
}

func shouldSkipSessionLine(lineType string) bool {
	switch lineType {
	case "queue-operation", "file-history-snapshot", "progress":
		return true
	default:
		return false
	}
}

// extractDisplayContent 从消息的 content 字段提取用于显示的纯文本
func extractDisplayContent(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var str string
	if json.Unmarshal(raw, &str) == nil {
		return cleanDisplayText(str)
	}
	var blocks []json.RawMessage
	if json.Unmarshal(raw, &blocks) != nil {
		return ""
	}

	var parts []string
	for _, blockRaw := range blocks {
		var block map[string]interface{}
		if json.Unmarshal(blockRaw, &block) != nil {
			continue
		}
		if blockType, _ := block["type"].(string); blockType == "text" {
			if text, ok := block["text"].(string); ok {
				cleaned := cleanDisplayText(text)
				if cleaned != "" {
					parts = append(parts, cleaned)
				}
			}
		}
	}
	return strings.Join(parts, "\n")
}

// extractAssistantContent 从 assistant 消息的 content 字段提取显示内容
func extractAssistantContent(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var str string
	if json.Unmarshal(raw, &str) == nil {
		return str
	}
	var blocks []json.RawMessage
	if json.Unmarshal(raw, &blocks) != nil {
		return ""
	}

	var parts []string
	for _, blockRaw := range blocks {
		var block map[string]interface{}
		if json.Unmarshal(blockRaw, &block) != nil {
			continue
		}
		blockType, _ := block["type"].(string)
		switch blockType {
		case "text":
			if text, ok := block["text"].(string); ok && text != "" {
				parts = append(parts, text)
			}
		case "tool_use":
			name, _ := block["name"].(string)
			if name != "" {
				parts = append(parts, fmt.Sprintf("[调用工具: %s]", name))
			}
		}
	}
	return strings.Join(parts, "\n")
}

var blockRemoveIdeFile = regexp.MustCompile(`(?s)<ide_opened_file>.*?</ide_opened_file>`)
var blockRemoveSysReminder = regexp.MustCompile(`(?s)<system-reminder>.*?</system-reminder>`)
var blockRemoveFileHistory = regexp.MustCompile(`(?s)<file-history-snapshot>.*?</file-history-snapshot>`)
var inlineStripTags = regexp.MustCompile(`</?(?:command-name|command-message|command-args|bash-input|local-command-stdout|local-command-stderr|antml:[a-z_]+)>`)

// cleanDisplayText 清理显示文本中的特殊 XML 标签
func cleanDisplayText(text string) string {
	text = blockRemoveIdeFile.ReplaceAllString(text, "")
	text = blockRemoveSysReminder.ReplaceAllString(text, "")
	text = blockRemoveFileHistory.ReplaceAllString(text, "")
	text = inlineStripTags.ReplaceAllString(text, "")
	return strings.TrimSpace(text)
}
