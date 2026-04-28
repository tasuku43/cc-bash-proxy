package hookinput

import "testing"

func TestNormalizeGenericExec(t *testing.T) {
	req, err := Normalize([]byte(`{"action":"exec","command":"git status"}`))
	if err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
	if req.Command != "git status" {
		t.Fatalf("got %q", req.Command)
	}
}

func TestNormalizeClaudeBash(t *testing.T) {
	req, err := Normalize([]byte(`{"tool_name":"Bash","tool_input":{"command":"git status"}}`))
	if err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
	if req.Action != "exec" || req.Command != "git status" {
		t.Fatalf("got %+v", req)
	}
}

func TestNormalizeClaudeBashPreservesToolInput(t *testing.T) {
	req, err := Normalize([]byte(`{"tool_name":"Bash","tool_input":{"command":"git status","description":"check status","extra":{"nested":true}}}`))
	if err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
	if req.Command != "git status" {
		t.Fatalf("got command %q", req.Command)
	}
	if req.OriginalToolInput["description"] != "check status" {
		t.Fatalf("tool input = %+v", req.OriginalToolInput)
	}
	extra, ok := req.OriginalToolInput["extra"].(map[string]any)
	if !ok || extra["nested"] != true {
		t.Fatalf("tool input = %+v", req.OriginalToolInput)
	}
}

func TestNormalizeRejectsUnknownAction(t *testing.T) {
	if _, err := Normalize([]byte(`{"action":"write","command":"x"}`)); err == nil {
		t.Fatal("expected error")
	}
}
