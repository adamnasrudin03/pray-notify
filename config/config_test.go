package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWebhookTargetsFromEnv(t *testing.T) {
	t.Setenv("WEBHOOK_TARGETS_FILE", "")
	t.Setenv("WEBHOOK_TARGETS", `[{"Name":"Test","URL":"https://example.com"}]`)

	targets, err := LoadWebhookTargets()
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(targets))
	}
	if targets[0].Name != "Test" {
		t.Fatalf("unexpected name: %s", targets[0].Name)
	}
	if targets[0].MentionMode != "here" {
		t.Fatalf("expected default MentionMode here, got %s", targets[0].MentionMode)
	}
}

func TestLoadWebhookTargetsFromFile(t *testing.T) {
	t.Setenv("WEBHOOK_TARGETS", "")

	dir := t.TempDir()
	path := filepath.Join(dir, "targets.json")
	if err := os.WriteFile(path, []byte(`[{"Name":"FileTarget","URL":"https://example.com"}]`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("WEBHOOK_TARGETS_FILE", path)

	targets, err := LoadWebhookTargets()
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 || targets[0].Name != "FileTarget" {
		t.Fatalf("unexpected targets: %+v", targets)
	}
}

func TestLoadWebhookTargetsEnvTakesPrecedence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "targets.json")
	if err := os.WriteFile(path, []byte(`[{"Name":"FromFile","URL":"https://example.com"}]`), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("WEBHOOK_TARGETS_FILE", path)
	t.Setenv("WEBHOOK_TARGETS", `[{"Name":"FromEnv","URL":"https://example.com"}]`)

	targets, err := LoadWebhookTargets()
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 || targets[0].Name != "FromEnv" {
		t.Fatalf("env should take precedence, got %+v", targets)
	}
}
