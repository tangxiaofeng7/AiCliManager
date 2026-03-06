package service

import (
	"AiCliManager/internal/db/models"
	"testing"
)

func TestClaudeSyncAdapterAppliesRuntimeFields(t *testing.T) {
	config := map[string]interface{}{}
	adapter := claudeSyncAdapter{}
	adapter.Apply(config, SyncRuntimeConfig{
		Provider: &models.Provider{ApiKey: "secret", ApiUrl: "https://api.example.com"},
		Profile:  &models.Profile{Model: "claude-opus", SystemPrompt: "hello"},
		McpServers: []models.McpServer{{
			Name:      "local",
			Type:      "stdio",
			Command:   "node",
			Args:      `["server.js"]`,
			Env:       `{"TOKEN":"1"}`,
			IsEnabled: 1,
		}},
		Skills:    []models.Skill{{Id: 1, Name: "review", Trigger: "/review", Content: "Hello {{name}}"}},
		SkillVars: map[string]string{"1.name": "Claude"},
	})

	if config["apiKey"] != "secret" {
		t.Fatalf("expected apiKey to be written")
	}
	if config["baseUrl"] != "https://api.example.com" {
		t.Fatalf("expected baseUrl to be written")
	}
	if config["model"] != "claude-opus" {
		t.Fatalf("expected model to be written")
	}
	if config["systemPrompt"] != "hello" {
		t.Fatalf("expected systemPrompt to be written")
	}
	commands, ok := config["customCommands"].([]map[string]interface{})
	if !ok || len(commands) != 1 {
		t.Fatalf("expected one custom command, got %#v", config["customCommands"])
	}
	if commands[0]["prompt"] != "Hello Claude" {
		t.Fatalf("expected skill vars to be applied, got %#v", commands[0]["prompt"])
	}
	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok || len(mcpServers) != 1 {
		t.Fatalf("expected one mcp server, got %#v", config["mcpServers"])
	}
}

func TestOpenCodeSyncAdapterDoesNotWriteSystemPromptWhenProfileMissing(t *testing.T) {
	config := map[string]interface{}{}
	adapter := openCodeSyncAdapter{}
	adapter.Apply(config, SyncRuntimeConfig{})
	if len(config) != 0 {
		t.Fatalf("expected empty config, got %#v", config)
	}
}

func TestApplySkillVars(t *testing.T) {
	result := applySkillVars("Hello {{name}} from {{city}}", 2, map[string]string{
		"2.name": "Alice",
		"2.city": "Shanghai",
		"3.name": "Bob",
	})
	if result != "Hello Alice from Shanghai" {
		t.Fatalf("unexpected result: %q", result)
	}
}
