package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTimestampSupportsMillis(t *testing.T) {
	svc := NewCliSessionService()
	result := svc.parseTimestamp([]byte(`1700000000000`))
	if result == "" {
		t.Fatal("expected parsed timestamp")
	}
}

func TestCleanDisplayTextRemovesSystemTags(t *testing.T) {
	input := "hello<system-reminder>secret</system-reminder><command-name>/test</command-name>world"
	result := cleanDisplayText(input)
	if result != "hello/testworld" {
		t.Fatalf("unexpected cleaned text: %q", result)
	}
}

func TestParseSessionMessages(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "session.jsonl")
	content := "{\"type\":\"user\",\"timestamp\":\"2026-01-01T00:00:00Z\",\"uuid\":\"u1\",\"message\":{\"content\":\"hello\"}}\n" +
		"{\"type\":\"assistant\",\"timestamp\":1700000000000,\"uuid\":\"a1\",\"message\":{\"model\":\"claude\",\"content\":[{\"type\":\"text\",\"text\":\"hi\"}],\"usage\":{\"input_tokens\":1,\"output_tokens\":2}}}\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	svc := NewCliSessionService()
	messages, err := svc.parseSessionMessages(filePath)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}
	if messages[1].TokensIn != 1 || messages[1].TokensOut != 2 {
		t.Fatalf("unexpected token usage: %#v", messages[1])
	}
}
