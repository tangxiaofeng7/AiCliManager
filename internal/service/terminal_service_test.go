package service

import (
	"AiCliManager/internal/db/models"
	"strings"
	"testing"
)

func TestBuildProxyURLEscapesCredentials(t *testing.T) {
	proxyURL := buildProxyURL(nil)
	if proxyURL != "" {
		t.Fatalf("expected empty proxy url, got %q", proxyURL)
	}

	url := buildProxyURL(&models.Proxy{Type: "http", Host: "proxy.local", Port: 8080, Username: "user name", Password: "p@ss word"})
	if !strings.Contains(url, "user%20name") {
		t.Fatalf("expected encoded username, got %q", url)
	}
	if !strings.Contains(url, "p%40ss%20word") {
		t.Fatalf("expected encoded password, got %q", url)
	}
}

func TestFilterEnvRemovesClaudeCode(t *testing.T) {
	env := []string{"PATH=/bin", "CLAUDECODE=1", "claudecode=2", "HOME=/tmp"}
	filtered := filterEnv(env)
	for _, item := range filtered {
		if strings.HasPrefix(strings.ToUpper(item), "CLAUDECODE=") {
			t.Fatalf("unexpected CLAUDECODE entry: %q", item)
		}
	}
	if len(filtered) != 2 {
		t.Fatalf("expected 2 env entries, got %d", len(filtered))
	}
}

func TestWindowsCommandLineQuotesArgs(t *testing.T) {
	cmd := buildWindowsCommandLine(`C:\Program Files\Claude\claude.exe`, []string{"--model", "claude sonnet", `say \"hi\"`})
	if !strings.Contains(cmd, `"C:\Program Files\Claude\claude.exe"`) {
		t.Fatalf("expected quoted executable, got %q", cmd)
	}
	if !strings.Contains(cmd, `"claude sonnet"`) {
		t.Fatalf("expected quoted arg, got %q", cmd)
	}
}

func TestWindowsPathToWSL(t *testing.T) {
	result := windowsPathToWSL(`C:\Users\demo\workspace`)
	if result != "/mnt/c/Users/demo/workspace" {
		t.Fatalf("unexpected wsl path: %q", result)
	}
}
