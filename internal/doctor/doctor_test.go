package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tasuku43/cmdproxy/internal/config"
)

func TestExtractClaudeHookCommandFindsAbsoluteCommand(t *testing.T) {
	raw := `{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "/tmp/cmdproxy hook claude --rtk" }
        ]
      }
    ]
  }
}`
	command, absolute := extractClaudeHookCommand(raw)
	if command != "/tmp/cmdproxy hook claude --rtk" {
		t.Fatalf("command = %q", command)
	}
	if !absolute {
		t.Fatalf("expected absolute path detection")
	}
}

func TestRunWarnsWhenClaudeHookTargetIsMissing(t *testing.T) {
	home := t.TempDir()
	settingsDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(settingsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	settings := `{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "/tmp/does-not-exist-cmdproxy hook claude" }
        ]
      }
    ]
  }
}`
	if err := os.WriteFile(filepath.Join(settingsDir, "settings.json"), []byte(settings), 0o644); err != nil {
		t.Fatal(err)
	}

	report := Run(config.Loaded{}, home)
	if !hasCheck(report, "install.claude-hook-target", StatusWarn) {
		t.Fatalf("report = %+v", report)
	}
}

func hasCheck(report Report, id string, status Status) bool {
	for _, check := range report.Checks {
		if check.ID == id && check.Status == status {
			return true
		}
	}
	return false
}
