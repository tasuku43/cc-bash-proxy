package invocation

import (
	"reflect"
	"testing"
)

func TestParseUnwrapsCommonWrappers(t *testing.T) {
	parsed := Parse("sudo -u root /usr/bin/env bash -c 'echo hi'")
	if parsed.Command != "bash" {
		t.Fatalf("Command = %q", parsed.Command)
	}
	if len(parsed.Args) < 2 || parsed.Args[0] != "-c" || parsed.Args[1] != "echo hi" {
		t.Fatalf("Args = %#v", parsed.Args)
	}
}

func TestTokensPreserveQuotedPayload(t *testing.T) {
	got := Tokens("bash -c 'git status'")
	if len(got) != 3 || got[2] != "git status" {
		t.Fatalf("Tokens() = %#v", got)
	}
}

func TestJoinRoundTripPreservesQuotedArgs(t *testing.T) {
	command := `aws s3 cp "hello world" s3://bucket/key`

	want := Tokens(command)
	got := Tokens(Join(want))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokens(Join(Tokens(command))) = %#v, want %#v", got, want)
	}
}

func TestJoinRoundTripPreservesEnvAssignmentWithSpaces(t *testing.T) {
	command := `FOO="hello world" env`

	want := Tokens(command)
	got := Tokens(Join(want))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokens(Join(Tokens(command))) = %#v, want %#v", got, want)
	}
}

func TestIsEnvAssignmentAcceptsEmptyValue(t *testing.T) {
	if !IsEnvAssignment("FOO=") {
		t.Fatal("expected empty env assignment to be treated as env assignment")
	}
}

func TestIsSafeSingleCommandRejectsCompoundPayload(t *testing.T) {
	if IsSafeSingleCommand("git status && git diff") {
		t.Fatal("expected safe single command check to fail")
	}
}
