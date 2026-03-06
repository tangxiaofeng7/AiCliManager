package service

import (
	"AiCliManager/internal/db/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// LaunchRequest 启动 CLI 工具的请求参数
type LaunchRequest struct {
	CliToolKey   string            `json:"cli_tool_key"`
	ProfileId    int64             `json:"profile_id"`
	ProxyId      *int64            `json:"proxy_id"`
	McpServerIds []int64           `json:"mcp_server_ids"`
	SkillIds     []int64           `json:"skill_ids"`
	SkillVars    map[string]string `json:"skill_vars"`
	ExtraArgs    []string          `json:"extra_args"`
	Terminal     string            `json:"terminal"`
	WorkingDir   string            `json:"working_dir"`
}

// DetectResult CLI 工具检测结果
type DetectResult struct {
	Key         string `json:"key"`
	IsInstalled bool   `json:"is_installed"`
	Executable  string `json:"executable"`
}

// ActiveConfig CLI 工具当前激活的配置
type ActiveConfig struct {
	ProfileId *int64 `json:"profile_id"`
	ProxyId   *int64 `json:"proxy_id"`
}

// SetActiveConfigRequest 设置激活配置的请求参数
type SetActiveConfigRequest struct {
	CliToolKey string `json:"cli_tool_key"`
	ProfileId  *int64 `json:"profile_id"`
	ProxyId    *int64 `json:"proxy_id"`
}

type launchContext struct {
	request          LaunchRequest
	tool             *models.CliTool
	executable       string
	profile          *models.Profile
	provider         *models.Provider
	proxy            *models.Proxy
	mcpServers       []models.McpServer
	skills           []models.Skill
	workingDir       string
	terminal         string
	cmd              *exec.Cmd
	pid              int
	sessionProfileId *int64
}

// LauncherService 负责 CLI 工具的检测与启动
type LauncherService struct {
	db              *sql.DB
	providerService *ProviderService
	proxyService    *ProxyService
	mcpService      *McpService
	skillService    *SkillService
	sessionService  *SessionService
	terminalService *TerminalService
	syncService     *SyncService
}

// NewLauncherService 创建 LauncherService 实例
func NewLauncherService(
	db *sql.DB,
	providerSvc *ProviderService,
	proxySvc *ProxyService,
	mcpSvc *McpService,
	skillSvc *SkillService,
	sessionSvc *SessionService,
	terminalSvc *TerminalService,
	syncSvc *SyncService,
) *LauncherService {
	return &LauncherService{
		db:              db,
		providerService: providerSvc,
		proxyService:    proxySvc,
		mcpService:      mcpSvc,
		skillService:    skillSvc,
		sessionService:  sessionSvc,
		terminalService: terminalSvc,
		syncService:     syncSvc,
	}
}

// GetCliTools 获取所有 CLI 工具列表
func (s *LauncherService) GetCliTools() ([]models.CliTool, error) {
	rows, err := s.db.Query(
		`SELECT id, name, key, executable, config_path, preferred_terminal,
		        is_installed, is_enabled, sort_order, created_at, updated_at
		 FROM cli_tools ORDER BY sort_order ASC, id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("查询 CLI 工具列表失败: %w", err)
	}
	defer rows.Close()

	var tools []models.CliTool
	for rows.Next() {
		var t models.CliTool
		if err := rows.Scan(
			&t.Id, &t.Name, &t.Key, &t.Executable, &t.ConfigPath,
			&t.PreferredTerminal, &t.IsInstalled, &t.IsEnabled,
			&t.SortOrder, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描 CLI 工具行数据失败: %w", err)
		}
		tools = append(tools, t)
	}
	return tools, nil
}

// getCliToolByKey 按 key 查询 CLI 工具记录
func (s *LauncherService) getCliToolByKey(key string) (*models.CliTool, error) {
	row := s.db.QueryRow(
		`SELECT id, name, key, executable, config_path, preferred_terminal,
		        is_installed, is_enabled, sort_order, created_at, updated_at
		 FROM cli_tools WHERE key=?`, key,
	)
	var t models.CliTool
	if err := row.Scan(
		&t.Id, &t.Name, &t.Key, &t.Executable, &t.ConfigPath,
		&t.PreferredTerminal, &t.IsInstalled, &t.IsEnabled,
		&t.SortOrder, &t.CreatedAt, &t.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("CLI 工具 [%s] 不存在", key)
		}
		return nil, fmt.Errorf("查询 CLI 工具失败: %w", err)
	}
	return &t, nil
}

// DetectCliTool 检测指定 CLI 工具是否已安装，并更新数据库记录
func (s *LauncherService) DetectCliTool(key string) (*DetectResult, error) {
	result := &DetectResult{Key: key}

	candidateMap := map[string][]string{
		"claude":   {"claude", "claude.exe"},
		"codex":    {"codex", "codex.exe"},
		"opencode": {"opencode", "opencode.exe"},
	}
	candidates, ok := candidateMap[key]
	if !ok {
		return nil, fmt.Errorf("未知 CLI 工具: %s", key)
	}

	for _, name := range candidates {
		path, err := exec.LookPath(name)
		if err == nil {
			result.IsInstalled = true
			result.Executable = path
			break
		}
	}

	if !result.IsInstalled {
		extraPaths := s.getCommonInstallPaths(key)
		for _, p := range extraPaths {
			if info, err := os.Stat(p); err == nil && !info.IsDir() {
				result.IsInstalled = true
				result.Executable = p
				break
			}
		}
	}

	installed := 0
	if result.IsInstalled {
		installed = 1
	}
	_, err := s.db.Exec(
		`UPDATE cli_tools SET is_installed=?, executable=?, updated_at=datetime('now') WHERE key=?`,
		installed, result.Executable, key,
	)
	if err != nil {
		return nil, fmt.Errorf("更新 CLI 工具检测结果失败: %w", err)
	}

	return result, nil
}

// getCommonInstallPaths 返回指定 CLI 工具在各平台下的常见安装路径列表
func (s *LauncherService) getCommonInstallPaths(key string) []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	var paths []string

	switch runtime.GOOS {
	case "windows":
		commonDirs := []string{
			filepath.Join(homeDir, ".local", "bin"),
			filepath.Join(os.Getenv("APPDATA"), "npm"),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", key),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "WinGet", "Packages"),
		}
		exeNames := []string{key + ".exe", key + ".cmd", key}
		for _, dir := range commonDirs {
			for _, name := range exeNames {
				paths = append(paths, filepath.Join(dir, name))
			}
		}
	case "darwin":
		commonDirs := []string{
			filepath.Join(homeDir, ".local", "bin"),
			"/usr/local/bin",
			"/opt/homebrew/bin",
			filepath.Join(homeDir, ".npm-global", "bin"),
		}
		for _, dir := range commonDirs {
			paths = append(paths, filepath.Join(dir, key))
		}
	case "linux":
		commonDirs := []string{
			filepath.Join(homeDir, ".local", "bin"),
			"/usr/local/bin",
			"/usr/bin",
			filepath.Join(homeDir, ".npm-global", "bin"),
		}
		for _, dir := range commonDirs {
			paths = append(paths, filepath.Join(dir, key))
		}
	}

	return paths
}

// GetCliToolActiveConfig 获取指定 CLI 工具当前激活的 Profile 和 Proxy
func (s *LauncherService) GetCliToolActiveConfig(key string) (*ActiveConfig, error) {
	tool, err := s.getCliToolByKey(key)
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow(
		`SELECT profile_id, proxy_id FROM cli_tool_active_config WHERE cli_tool_id=?`, tool.Id,
	)
	var cfg ActiveConfig
	if err := row.Scan(&cfg.ProfileId, &cfg.ProxyId); err != nil {
		if err == sql.ErrNoRows {
			return &ActiveConfig{}, nil
		}
		return nil, fmt.Errorf("查询激活配置失败: %w", err)
	}
	return &cfg, nil
}

// SetCliToolActiveConfig 设置指定 CLI 工具的激活 Profile 和 Proxy
func (s *LauncherService) SetCliToolActiveConfig(req SetActiveConfigRequest) error {
	tool, err := s.getCliToolByKey(req.CliToolKey)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		`INSERT INTO cli_tool_active_config (cli_tool_id, profile_id, proxy_id, updated_at)
		 VALUES (?, ?, ?, datetime('now'))
		 ON CONFLICT(cli_tool_id) DO UPDATE SET
		   profile_id=excluded.profile_id,
		   proxy_id=excluded.proxy_id,
		   updated_at=excluded.updated_at`,
		tool.Id, req.ProfileId, req.ProxyId,
	)
	if err != nil {
		return fmt.Errorf("设置激活配置失败: %w", err)
	}
	return nil
}

// Launch 启动指定 CLI 工具：同步配置 → 构建环境变量 → 启动进程 → 写入会话历史
func (s *LauncherService) Launch(req LaunchRequest) error {
	ctx, err := s.prepareLaunchContext(req)
	if err != nil {
		return err
	}

	if err := s.syncLaunchConfig(ctx); err != nil {
		return err
	}
	if err := s.buildLaunchCommand(ctx); err != nil {
		return err
	}
	if err := s.startProcess(ctx); err != nil {
		return err
	}
	if err := s.recordSession(ctx); err != nil {
		return err
	}
	return nil
}

func (s *LauncherService) prepareLaunchContext(req LaunchRequest) (*launchContext, error) {
	tool, err := s.getCliToolByKey(req.CliToolKey)
	if err != nil {
		return nil, err
	}

	executable, err := s.resolveExecutable(req.CliToolKey, tool.Executable)
	if err != nil {
		return nil, err
	}

	profile, provider, err := s.resolveProfileProvider(req.ProfileId)
	if err != nil {
		return nil, err
	}

	proxy, err := s.resolveProxy(req.ProxyId)
	if err != nil {
		return nil, err
	}

	mcpServers, err := s.mcpService.GetByIds(req.McpServerIds)
	if err != nil {
		return nil, fmt.Errorf("查询 MCP Server 失败: %w", err)
	}

	skills, err := s.skillService.GetByIds(req.SkillIds)
	if err != nil {
		return nil, fmt.Errorf("查询 Skill 失败: %w", err)
	}

	workingDir, err := resolveWorkingDir(req.WorkingDir)
	if err != nil {
		return nil, err
	}

	terminal := req.Terminal
	if terminal == "" {
		terminal = tool.PreferredTerminal
	}

	ctx := &launchContext{
		request:    req,
		tool:       tool,
		executable: executable,
		profile:    profile,
		provider:   provider,
		proxy:      proxy,
		mcpServers: mcpServers,
		skills:     skills,
		workingDir: workingDir,
		terminal:   terminal,
	}
	if req.ProfileId > 0 {
		ctx.sessionProfileId = &req.ProfileId
	}
	return ctx, nil
}

func (s *LauncherService) syncLaunchConfig(ctx *launchContext) error {
	err := s.syncService.SyncConfig(SyncRequest{
		CliToolKey: ctx.request.CliToolKey,
		CliTool:    ctx.tool,
		Runtime: SyncRuntimeConfig{
			Profile:    ctx.profile,
			Provider:   ctx.provider,
			Proxy:      ctx.proxy,
			McpServers: ctx.mcpServers,
			Skills:     ctx.skills,
			SkillVars:  ctx.request.SkillVars,
		},
	})
	if err != nil {
		return fmt.Errorf("同步配置失败: %w", err)
	}
	return nil
}

func (s *LauncherService) buildLaunchCommand(ctx *launchContext) error {
	cmd, err := s.terminalService.BuildCmd(TerminalLaunchSpec{
		Terminal:   ctx.terminal,
		Executable: ctx.executable,
		Args:       ctx.request.ExtraArgs,
		WorkingDir: ctx.workingDir,
		KeepOpen:   true,
	})
	if err != nil {
		return fmt.Errorf("构建启动命令失败: %w", err)
	}
	if cmd.Dir == "" {
		cmd.Dir = ctx.workingDir
	}
	cmd.Env = append(filterEnv(os.Environ()), buildProxyEnv(ctx.proxy)...)
	ctx.cmd = cmd
	return nil
}

func (s *LauncherService) startProcess(ctx *launchContext) error {
	if ctx.cmd == nil {
		return fmt.Errorf("启动命令尚未构建")
	}
	if err := ctx.cmd.Start(); err != nil {
		return fmt.Errorf("启动 CLI 工具失败: %w", err)
	}
	if ctx.cmd.Process != nil {
		ctx.pid = ctx.cmd.Process.Pid
	}
	return nil
}

func (s *LauncherService) recordSession(ctx *launchContext) error {
	sess, err := s.sessionService.Create(
		ctx.tool.Id,
		ctx.sessionProfileId,
		ctx.request.ProxyId,
		ctx.terminal,
		ctx.workingDir,
		ctx.request.ExtraArgs,
		ctx.pid,
	)
	if err != nil {
		return fmt.Errorf("CLI 已启动，但写入会话记录失败: %w", err)
	}

	go func(sessionId int64, process *exec.Cmd) {
		if process == nil || process.Process == nil {
			return
		}
		waitErr := process.Wait()
		status := "exited"
		if waitErr != nil {
			status = "error"
		}
		_ = s.sessionService.UpdateStatus(sessionId, status)
	}(sess.Id, ctx.cmd)

	return nil
}

func (s *LauncherService) resolveExecutable(cliToolKey, executable string) (string, error) {
	if executable != "" {
		return executable, nil
	}
	det, err := s.DetectCliTool(cliToolKey)
	if err != nil {
		return "", fmt.Errorf("检测 CLI 工具失败: %w", err)
	}
	if !det.IsInstalled || det.Executable == "" {
		return "", fmt.Errorf("CLI 工具 [%s] 未安装或未找到可执行文件", cliToolKey)
	}
	return det.Executable, nil
}

func (s *LauncherService) resolveProfileProvider(profileId int64) (*models.Profile, *models.Provider, error) {
	if profileId <= 0 {
		return nil, nil, nil
	}

	profile := &models.Profile{}
	err := s.db.QueryRow(
		`SELECT id, name, provider_id, model, system_prompt, temperature, max_tokens, extra_config, created_at, updated_at
		 FROM profiles WHERE id=?`, profileId,
	).Scan(
		&profile.Id, &profile.Name, &profile.ProviderId, &profile.Model,
		&profile.SystemPrompt, &profile.Temperature, &profile.MaxTokens,
		&profile.ExtraConfig, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("查询 Profile 失败: %w", err)
	}

	provider := &models.Provider{}
	err = s.db.QueryRow(
		`SELECT id, name, type, api_url, api_key, models, sort_order, created_at, updated_at
		 FROM providers WHERE id=?`, profile.ProviderId,
	).Scan(
		&provider.Id, &provider.Name, &provider.Type, &provider.ApiUrl,
		&provider.ApiKey, &provider.Models, &provider.SortOrder,
		&provider.CreatedAt, &provider.UpdatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("查询 Provider 失败: %w", err)
	}

	decryptedKey, err := s.providerService.GetDecryptedApiKey(provider.Id)
	if err != nil {
		return nil, nil, fmt.Errorf("解密 API Key 失败: %w", err)
	}
	provider.ApiKey = decryptedKey
	return profile, provider, nil
}

func (s *LauncherService) resolveProxy(proxyId *int64) (*models.Proxy, error) {
	if proxyId == nil || *proxyId <= 0 {
		return nil, nil
	}
	proxy, err := s.proxyService.GetDecryptedProxy(*proxyId)
	if err != nil {
		return nil, fmt.Errorf("查询代理失败: %w", err)
	}
	return proxy, nil
}

func resolveWorkingDir(workingDir string) (string, error) {
	if workingDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("获取默认工作目录失败: %w", err)
		}
		workingDir = homeDir
	}
	if info, err := os.Stat(workingDir); err != nil {
		return "", fmt.Errorf("工作目录不可访问: %w", err)
	} else if !info.IsDir() {
		return "", fmt.Errorf("工作目录不是目录: %s", workingDir)
	}
	return workingDir, nil
}

// filterEnv 过滤掉可能导致子进程启动失败的环境变量
func filterEnv(env []string) []string {
	blocked := map[string]bool{
		"CLAUDECODE": true,
	}

	var filtered []string
	for _, e := range env {
		key := e
		if idx := strings.IndexByte(e, '='); idx >= 0 {
			key = e[:idx]
		}
		if !blocked[strings.ToUpper(key)] {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// buildProxyEnv 根据代理配置构建注入子进程的环境变量列表
func buildProxyEnv(proxy *models.Proxy) []string {
	if proxy == nil {
		return nil
	}

	proxyURL := buildProxyURL(proxy)
	envs := []string{
		"HTTP_PROXY=" + proxyURL,
		"HTTPS_PROXY=" + proxyURL,
		"http_proxy=" + proxyURL,
		"https_proxy=" + proxyURL,
	}
	if proxy.NoProxy != "" {
		envs = append(envs, "NO_PROXY="+proxy.NoProxy)
		envs = append(envs, "no_proxy="+proxy.NoProxy)
	}
	return envs
}

func buildProxyURL(proxy *models.Proxy) string {
	if proxy == nil {
		return ""
	}
	host := proxy.Host
	if proxy.Username == "" {
		return fmt.Sprintf("%s://%s:%d", proxy.Type, host, proxy.Port)
	}
	userInfo := url.User(proxy.Username)
	if proxy.Password != "" {
		userInfo = url.UserPassword(proxy.Username, proxy.Password)
	}
	return fmt.Sprintf("%s://%s@%s:%d", proxy.Type, userInfo.String(), host, proxy.Port)
}

// RelaunchSession 用历史会话的配置重新启动 CLI 工具
func (s *LauncherService) RelaunchSession(sessionId int64) error {
	sess, err := s.sessionService.GetById(sessionId)
	if err != nil {
		return err
	}

	var cliToolKey string
	if err := s.db.QueryRow(`SELECT key FROM cli_tools WHERE id=?`, sess.CliToolId).Scan(&cliToolKey); err != nil {
		return fmt.Errorf("查询 CLI 工具 key 失败: %w", err)
	}

	var extraArgs []string
	if sess.ExtraArgs != "" && sess.ExtraArgs != "[]" {
		_ = json.Unmarshal([]byte(sess.ExtraArgs), &extraArgs)
	}

	req := LaunchRequest{
		CliToolKey: cliToolKey,
		ExtraArgs:  extraArgs,
		Terminal:   sess.Terminal,
		WorkingDir: sess.WorkingDir,
		ProxyId:    sess.ProxyId,
	}
	if sess.ProfileId != nil {
		req.ProfileId = *sess.ProfileId
	}

	return s.Launch(req)
}
