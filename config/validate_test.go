package config

import "testing"

func TestValidateWebhookURL(t *testing.T) {
	valid := "https://discord.com/api/webhooks/123456789/abcdefghijklmnop"
	if err := ValidateWebhookURL(valid); err != nil {
		t.Fatalf("expected valid: %v", err)
	}

	if err := ValidateWebhookURL("https://discord.com/api/webhooks/REPLACE_ME"); err == nil {
		t.Fatal("expected error for REPLACE_ME")
	}

	if err := ValidateWebhookURL("https://api.github.com/repos/foo/bar/dispatches"); err == nil {
		t.Fatal("expected error for non-discord URL")
	}
}
