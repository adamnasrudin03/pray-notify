package discord

import (
	"strings"
	"testing"

	"github.com/adamnasrudin03/pray-notify/config"
	"github.com/adamnasrudin03/pray-notify/models"
)

func TestBuildWebhookContent(t *testing.T) {
	target := models.WebhookTarget{MentionMode: "here"}
	content := BuildWebhookContent(target, models.Dzuhur, models.SebelumAdzan)

	if !strings.HasPrefix(content, "📢 @here") {
		t.Fatalf("expected @here prefix, got: %q", content)
	}
	if strings.Contains(content, "Cc:") {
		t.Fatal("Cc line should not be present")
	}
}

func TestResolveUsername(t *testing.T) {
	target := models.WebhookTarget{UsernameBeforeAdzan: config.MyUsername}
	name := ResolveUsername(target, models.Dzuhur, models.SebelumAdzan, false)
	if name != config.MyUsername {
		t.Fatalf("got %q, want %q", name, config.MyUsername)
	}

	target.Username = "Bot Khusus"
	name = ResolveUsername(target, models.Dzuhur, models.SebelumAdzan, false)
	if name != "Bot Khusus" {
		t.Fatalf("Username override should take priority, got %q", name)
	}

	name = ResolveUsername(target, models.Dzuhur, models.SebelumAdzan, true)
	if name != config.Test {
		t.Fatalf("testing mode should use Test username, got %q", name)
	}
}
