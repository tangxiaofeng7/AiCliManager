package service

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// TerminalInfo 描述一个可用的终端
type TerminalInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"is_available"`
}

// TerminalLaunchSpec 描述一次 CLI 启动请求
type TerminalLaunchSpec struct {
	Terminal   string
	Executable string
	Args       []string
	WorkingDir string
	KeepOpen   bool
}

// TerminalService 提供跨平台终端适配逻辑
type TerminalService struct{}

// NewTerminalService 创建 TerminalService 实例
func NewTerminalService() *TerminalService {
	return &TerminalService{}
}

// ListAvailableTerminals 返回当前平台所有支持的终端及其可用状态
func (s *TerminalService) ListAvailableTerminals() []TerminalInfo {
	var terminals []TerminalInfo
	switch runtime.GOOS {
	case "windows":
		candidates := []TerminalInfo{
			{Id: "wt", Name: "Windows Terminal"},
			{Id: "powershell", Name: "PowerShell"},
			{Id: "cmd", Name: "命令提示符 (cmd)"},
			{Id: "wsl", Name: "WSL"},
			{Id: "default", Name: "系统默认终端"},
		}
		for _, t := range candidates {
			t.IsAvailable = s.isTerminalAvailable(t.Id)
			terminals = append(terminals, t)
		}
	case "darwin":
		candidates := []TerminalInfo{
			{Id: "iterm2", Name: "iTerm2"},
			{Id: "terminal", Name: "Terminal.app"},
			{Id: "default", Name: "系统默认终端"},
		}
		for _, t := range candidates {
			t.IsAvailable = s.isTerminalAvailable(t.Id)
			terminals = append(terminals, t)
		}
	default:
		candidates := []TerminalInfo{
			{Id: "gnome-terminal", Name: "GNOME Terminal"},
			{Id: "tmux", Name: "Tmux"},
			{Id: "default", Name: "系统默认终端"},
		}
		for _, t := range candidates {
			t.IsAvailable = s.isTerminalAvailable(t.Id)
			terminals = append(terminals, t)
		}
	}
	return terminals
}

// isTerminalAvailable 检测指定终端是否已安装
func (s *TerminalService) isTerminalAvailable(terminalId string) bool {
	switch terminalId {
	case "default":
		return true
	case "wt":
		_, err := exec.LookPath("wt.exe")
		return err == nil
	case "powershell":
		_, err := exec.LookPath("powershell.exe")
		return err == nil
	case "cmd":
		_, err := exec.LookPath("cmd.exe")
		return err == nil
	case "wsl":
		_, err := exec.LookPath("wsl.exe")
		return err == nil
	case "iterm2":
		cmd := exec.Command("osascript", "-e", `application "iTerm" exists`)
		out, err := cmd.Output()
		return err == nil && strings.TrimSpace(string(out)) == "true"
	case "terminal":
		return runtime.GOOS == "darwin"
	case "gnome-terminal":
		_, err := exec.LookPath("gnome-terminal")
		return err == nil
	case "tmux":
		_, err := exec.LookPath("tmux")
		return err == nil
	default:
		_, err := exec.LookPath(terminalId)
		return err == nil
	}
}

// BuildCmd 将 CLI 命令包装为在指定终端中执行的 *exec.Cmd
func (s *TerminalService) BuildCmd(spec TerminalLaunchSpec) (*exec.Cmd, error) {
	if spec.Executable == "" {
		return nil, fmt.Errorf("CLI 可执行文件不能为空")
	}

	switch runtime.GOOS {
	case "windows":
		return s.buildWindowsCmd(spec)
	case "darwin":
		return s.buildDarwinCmd(spec)
	default:
		return s.buildLinuxCmd(spec)
	}
}

func (s *TerminalService) buildWindowsCmd(spec TerminalLaunchSpec) (*exec.Cmd, error) {
	terminal := spec.Terminal
	if terminal == "" || terminal == "default" {
		if s.isTerminalAvailable("wt") {
			terminal = "wt"
		} else {
			terminal = "cmd"
		}
	}

	switch terminal {
	case "wt":
		cmdLine := buildWindowsCommandLine(spec.Executable, spec.Args)
		wrapper := buildWindowsCmdShellCommand(spec.WorkingDir, cmdLine, spec.KeepOpen)
		return exec.Command("wt.exe", "new-tab", "cmd.exe", "/d", "/k", wrapper), nil
	case "cmd":
		cmdLine := buildWindowsCommandLine(spec.Executable, spec.Args)
		wrapper := buildWindowsCmdShellCommand(spec.WorkingDir, cmdLine, spec.KeepOpen)
		return exec.Command("cmd.exe", "/c", "start", "", "cmd.exe", "/d", "/k", wrapper), nil
	case "powershell":
		cmdLine := buildWindowsCommandLine(spec.Executable, spec.Args)
		wrapper := buildWindowsPowerShellCommand(spec.WorkingDir, cmdLine, spec.KeepOpen)
		return exec.Command("cmd.exe", "/c", "start", "", "powershell.exe", "-NoLogo", "-NoExit", "-Command", wrapper), nil
	case "wsl":
		args := []string{"--cd", windowsPathToWSL(spec.WorkingDir), "--", spec.Executable}
		args = append(args, spec.Args...)
		return exec.Command("wsl.exe", args...), nil
	default:
		return nil, fmt.Errorf("不支持的终端类型: %s", spec.Terminal)
	}
}

func (s *TerminalService) buildDarwinCmd(spec TerminalLaunchSpec) (*exec.Cmd, error) {
	cliCmd := buildPosixCommandLine(spec.Executable, spec.Args)
	wrapped := wrapWithPosixWorkingDir(spec.WorkingDir, cliCmd)
	if spec.KeepOpen {
		wrapped += "; exec $SHELL"
	}

	switch spec.Terminal {
	case "iterm2":
		script := fmt.Sprintf(`tell application "iTerm"
			create window with default profile command %s
			activate
		end tell`, quoteAppleScriptString(wrapped))
		return exec.Command("osascript", "-e", script), nil
	case "", "default", "terminal":
		script := fmt.Sprintf(`tell application "Terminal"
			do script %s
			activate
		end tell`, quoteAppleScriptString(wrapped))
		return exec.Command("osascript", "-e", script), nil
	default:
		return nil, fmt.Errorf("不支持的终端类型: %s", spec.Terminal)
	}
}

func (s *TerminalService) buildLinuxCmd(spec TerminalLaunchSpec) (*exec.Cmd, error) {
	cliCmd := buildPosixCommandLine(spec.Executable, spec.Args)
	wrapped := wrapWithPosixWorkingDir(spec.WorkingDir, cliCmd)
	if spec.KeepOpen {
		wrapped += "; exec bash"
	}

	terminal := spec.Terminal
	if terminal == "" || terminal == "default" {
		for _, candidate := range []string{"gnome-terminal", "xterm", "konsole"} {
			if _, err := exec.LookPath(candidate); err == nil {
				terminal = candidate
				break
			}
		}
		if terminal == "" || terminal == "default" {
			cmd := exec.Command(spec.Executable, spec.Args...)
			cmd.Dir = spec.WorkingDir
			return cmd, nil
		}
	}

	switch terminal {
	case "gnome-terminal":
		return exec.Command("gnome-terminal", "--", "bash", "-lc", wrapped), nil
	case "xterm":
		return exec.Command("xterm", "-e", "bash", "-lc", wrapped), nil
	case "konsole":
		return exec.Command("konsole", "-e", "bash", "-lc", wrapped), nil
	case "tmux":
		return exec.Command("tmux", "new-window", wrapped), nil
	default:
		return nil, fmt.Errorf("不支持的终端类型: %s", spec.Terminal)
	}
}

func buildWindowsCmdShellCommand(workingDir, cliCmd string, keepOpen bool) string {
	parts := make([]string, 0, 3)
	if workingDir != "" {
		parts = append(parts, fmt.Sprintf("cd /d %s", quoteWindowsArg(workingDir)))
	}
	parts = append(parts, cliCmd)
	if keepOpen {
		parts = append(parts, "echo.")
	}
	return strings.Join(parts, " && ")
}

func buildWindowsPowerShellCommand(workingDir, cliCmd string, keepOpen bool) string {
	parts := make([]string, 0, 3)
	if workingDir != "" {
		parts = append(parts, fmt.Sprintf("Set-Location -LiteralPath %s", quotePowerShellString(workingDir)))
	}
	parts = append(parts, cliCmd)
	if keepOpen {
		parts = append(parts, "Write-Host ''")
	}
	return strings.Join(parts, "; ")
}

func wrapWithPosixWorkingDir(workingDir, cliCmd string) string {
	if workingDir == "" {
		return cliCmd
	}
	return "cd " + quotePosixArg(workingDir) + " && " + cliCmd
}

func buildWindowsCommandLine(executable string, args []string) string {
	parts := []string{quoteWindowsArg(executable)}
	for _, arg := range args {
		parts = append(parts, quoteWindowsArg(arg))
	}
	return strings.Join(parts, " ")
}

func buildPosixCommandLine(executable string, args []string) string {
	parts := []string{quotePosixArg(executable)}
	for _, arg := range args {
		parts = append(parts, quotePosixArg(arg))
	}
	return strings.Join(parts, " ")
}

func quoteWindowsArg(value string) string {
	if value == "" {
		return `""`
	}

	needsQuotes := strings.ContainsAny(value, " \t\n\v\"")
	if !needsQuotes {
		return value
	}

	var b strings.Builder
	b.WriteByte('"')
	backslashes := 0
	for _, r := range value {
		switch r {
		case '\\':
			backslashes++
		case '"':
			b.WriteString(strings.Repeat("\\", backslashes*2+1))
			b.WriteRune('"')
			backslashes = 0
		default:
			if backslashes > 0 {
				b.WriteString(strings.Repeat("\\", backslashes))
				backslashes = 0
			}
			b.WriteRune(r)
		}
	}
	if backslashes > 0 {
		b.WriteString(strings.Repeat("\\", backslashes*2))
	}
	b.WriteByte('"')
	return b.String()
}

func quotePosixArg(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", `'"'"'`) + "'"
}

func quotePowerShellString(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func quoteAppleScriptString(value string) string {
	replacer := strings.NewReplacer(`\\`, `\\\\`, `"`, `\\"`, "\n", `\\n`)
	return `"` + replacer.Replace(value) + `"`
}

func windowsPathToWSL(path string) string {
	if path == "" {
		return "~"
	}
	volume := filepath.VolumeName(path)
	cleanPath := filepath.Clean(path)
	if len(volume) == 2 && volume[1] == ':' {
		drive := strings.ToLower(string(volume[0]))
		rest := strings.TrimPrefix(cleanPath, volume)
		rest = strings.ReplaceAll(rest, `\`, "/")
		if rest == "" {
			rest = "/"
		}
		return "/mnt/" + drive + rest
	}
	return strings.ReplaceAll(cleanPath, `\`, "/")
}
