package policy

import (
	"testing"

	commandpkg "github.com/tasuku43/cc-bash-guard/internal/domain/command"
)

type securityTraceWant struct {
	name   string
	effect string
}

func TestSecurityRegressionMatrixEvaluationBoundaries(t *testing.T) {
	gitRule := func(subcommand string) PermissionRuleSpec {
		return PermissionRuleSpec{Command: PermissionCommandSpec{Name: "git", Semantic: &SemanticMatchSpec{Verb: subcommand}}}
	}
	allowGitReadOnly := []PermissionRuleSpec{gitRule("status"), gitRule("diff"), gitRule("log")}

	tests := []struct {
		name       string
		category   string
		command    string
		permission PermissionSpec
		want       string
		shape      commandpkg.ShellShapeKind
		flags      []string
		trace      []securityTraceWant
	}{
		{
			name:       "and list composes per command",
			category:   "compound",
			command:    "git status && git diff",
			permission: PermissionSpec{Allow: allowGitReadOnly},
			want:       "allow",
			shape:      commandpkg.ShellShapeCompound,
			flags:      []string{"conditional"},
			trace:      []securityTraceWant{{name: "composition", effect: "allow"}},
		},
		{
			name:       "or list composes per command",
			category:   "compound",
			command:    "git status || git diff",
			permission: PermissionSpec{Allow: allowGitReadOnly},
			want:       "allow",
			shape:      commandpkg.ShellShapeCompound,
			flags:      []string{"conditional"},
			trace:      []securityTraceWant{{name: "composition", effect: "allow"}},
		},
		{
			name:       "pipeline requires right side allow",
			category:   "compound",
			command:    "git status | sh",
			permission: PermissionSpec{Allow: []PermissionRuleSpec{gitRule("status")}},
			want:       "ask",
			shape:      commandpkg.ShellShapeCompound,
			flags:      []string{"pipeline"},
			trace:      []securityTraceWant{{name: "composition", effect: "ask"}},
		},
		{
			name:       "sequence composes per command",
			category:   "compound",
			command:    "git status; git diff",
			permission: PermissionSpec{Allow: allowGitReadOnly},
			want:       "allow",
			shape:      commandpkg.ShellShapeCompound,
			flags:      []string{"sequence"},
			trace:      []securityTraceWant{{name: "composition", effect: "allow"}},
		},
		{
			name:       "nested compound fails closed",
			category:   "compound",
			command:    "git status && (git diff)",
			permission: PermissionSpec{Allow: allowGitReadOnly},
			want:       "ask",
			shape:      commandpkg.ShellShapeCompound,
			flags:      []string{"conditional", "subshell"},
			trace:      []securityTraceWant{{name: "fail_closed", effect: "ask"}, {name: "composition", effect: "ask"}},
		},
		{
			name:       "subshell is never auto allowed",
			category:   "shell_features",
			command:    "(git status)",
			permission: PermissionSpec{Allow: allowGitReadOnly},
			want:       "ask",
			shape:      commandpkg.ShellShapeCompound,
			flags:      []string{"subshell"},
			trace:      []securityTraceWant{{name: "fail_closed", effect: "ask"}, {name: "composition", effect: "ask"}},
		},
		{
			name:       "command substitution is never auto allowed",
			category:   "shell_features",
			command:    "echo $(git status)",
			permission: PermissionSpec{Allow: append([]PermissionRuleSpec{{Command: PermissionCommandSpec{Name: "echo"}}}, allowGitReadOnly...)},
			want:       "ask",
			shape:      commandpkg.ShellShapeCompound,
			flags:      []string{"command_substitution"},
			trace:      []securityTraceWant{{name: "fail_closed", effect: "ask"}, {name: "composition", effect: "ask"}},
		},
		{
			name:       "process substitution extracts deny but cannot allow",
			category:   "shell_features",
			command:    "cat <(rm -rf /tmp/x)",
			permission: PermissionSpec{Deny: []PermissionRuleSpec{{Command: PermissionCommandSpec{Name: "rm"}}}},
			want:       "deny",
			shape:      commandpkg.ShellShapeCompound,
			flags:      []string{"process_substitution"},
			trace:      []securityTraceWant{{name: "fail_closed", effect: "ask"}, {name: "composition", effect: "deny"}},
		},
		{
			name:       "redirection is never auto allowed",
			category:   "shell_features",
			command:    "git status > /tmp/out",
			permission: PermissionSpec{Allow: allowGitReadOnly},
			want:       "ask",
			shape:      commandpkg.ShellShapeCompound,
			flags:      []string{"redirection"},
			trace:      []securityTraceWant{{name: "fail_closed", effect: "ask"}, {name: "composition", effect: "ask"}},
		},
		{
			name:       "patterns deny cannot be upgraded by allow",
			category:   "permission",
			command:    "git status",
			permission: PermissionSpec{Deny: []PermissionRuleSpec{{Patterns: []string{`^\s*git\s+status\s*$`}}}, Allow: []PermissionRuleSpec{gitRule("status")}},
			want:       "deny",
			shape:      commandpkg.ShellShapeSimple,
			trace:      []securityTraceWant{{effect: "deny"}},
		},
		{
			name:       "patterns ask cannot be upgraded by allow",
			category:   "permission",
			command:    "git status",
			permission: PermissionSpec{Ask: []PermissionRuleSpec{{Patterns: []string{`^\s*git\s+status\s*$`}}}, Allow: []PermissionRuleSpec{gitRule("status")}},
			want:       "ask",
			shape:      commandpkg.ShellShapeSimple,
			trace:      []securityTraceWant{{effect: "ask"}},
		},
		{
			name:       "patterns allow does not broaden across compound commands",
			category:   "permission",
			command:    "git status && rm -rf /tmp/x",
			permission: PermissionSpec{Allow: []PermissionRuleSpec{{Patterns: []string{`^\s*git\s+status\s*&&\s*rm\s+-rf\s+/tmp/x\s*$`}}}},
			want:       "ask",
			shape:      commandpkg.ShellShapeCompound,
			flags:      []string{"conditional"},
			trace:      []securityTraceWant{{name: "composition", effect: "ask"}},
		},
		{
			name:       "command deny beats broad patterns allow",
			category:   "permission",
			command:    "git status && rm -rf /tmp/x",
			permission: PermissionSpec{Deny: []PermissionRuleSpec{{Command: PermissionCommandSpec{Name: "rm"}}}, Allow: []PermissionRuleSpec{{Patterns: []string{`.*`}, Message: "broad patterns allow"}}},
			want:       "deny",
			shape:      commandpkg.ShellShapeCompound,
			flags:      []string{"conditional"},
			trace:      []securityTraceWant{{name: "composition", effect: "deny"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.category+"/"+tt.name, func(t *testing.T) {
			p := NewPipeline(PipelineSpec{Permission: tt.permission}, Source{})
			got, err := Evaluate(p, tt.command)
			if err != nil {
				t.Fatalf("Evaluate() error = %v", err)
			}
			assertSecurityDecision(t, got, tt.want, tt.shape, tt.flags, tt.trace)
		})
	}
}

func TestUnsafeShellShapesDoNotBecomeStructuredAllow(t *testing.T) {
	allowVisibleCommands := PermissionSpec{Allow: []PermissionRuleSpec{
		{Command: PermissionCommandSpec{Name: "git", Semantic: &SemanticMatchSpec{Verb: "status"}}},
		{Command: PermissionCommandSpec{Name: "git", Semantic: &SemanticMatchSpec{Verb: "diff"}}},
		{Command: PermissionCommandSpec{Name: "echo"}},
		{Command: PermissionCommandSpec{Name: "cat"}},
	}}

	tests := []string{
		"git status > /tmp/out",
		"echo $(git status)",
		"cat <(git status)",
		"(git status)",
		"git status &",
		"git status && (git diff)",
	}

	p := NewPipeline(PipelineSpec{Permission: allowVisibleCommands}, Source{})
	for _, command := range tests {
		t.Run(command, func(t *testing.T) {
			got, err := Evaluate(p, command)
			if err != nil {
				t.Fatalf("Evaluate() error = %v", err)
			}
			if got.Outcome == "allow" {
				t.Fatalf("unsafe shell shape was auto-allowed; decision=%+v", got)
			}
		})
	}
}

func TestSupportedWrappersDoNotHideDangerousInnerCommands(t *testing.T) {
	force := true
	p := NewPipeline(PipelineSpec{Permission: PermissionSpec{
		Deny: []PermissionRuleSpec{{
			Command: PermissionCommandSpec{Name: "git", Semantic: &SemanticMatchSpec{Verb: "push", Force: &force}},
		}},
		Allow: []PermissionRuleSpec{{Command: PermissionCommandSpec{Name: "git"}}},
	}}, Source{})

	tests := []string{
		"bash -c 'git push --force origin main'",
		"sh -c 'git push --force origin main'",
		"/bin/bash -c 'git push --force origin main'",
		"env bash -c 'git push --force origin main'",
		"command bash -c 'git push --force origin main'",
		"sudo -u root bash -c 'git push --force origin main'",
		"timeout 10 bash -c 'git push --force origin main'",
		"busybox sh -c 'git push --force origin main'",
	}

	for _, command := range tests {
		t.Run(command, func(t *testing.T) {
			got, err := Evaluate(p, command)
			if err != nil {
				t.Fatalf("Evaluate() error = %v", err)
			}
			if got.Outcome != "deny" {
				t.Fatalf("Outcome = %q, want deny; decision=%+v", got.Outcome, got)
			}
		})
	}
}

func TestSemanticParserFallbackDoesNotWidenSemanticRuleToCommandAllow(t *testing.T) {
	force := true
	p := NewPipeline(PipelineSpec{Permission: PermissionSpec{
		Allow: []PermissionRuleSpec{{
			Command: PermissionCommandSpec{Name: "git", Semantic: &SemanticMatchSpec{Verb: "push", Force: &force}},
		}},
	}}, Source{})

	cmd := commandpkg.NewInvocation("git push --force origin main")
	parsed, ok := commandpkg.NewCommandParserRegistry().Parse(cmd)
	if !ok {
		t.Fatal("generic parser did not parse git command")
	}
	if parsed.SemanticParser != "" || parsed.Git != nil {
		t.Fatalf("test setup expected generic command without git semantic fields; cmd=%+v", parsed)
	}

	decision := evaluatePreparedCommand(p.prepared.Deny, p.prepared.Ask, p.prepared.Allow, parsed)
	if decision.Outcome == "allow" {
		t.Fatalf("semantic parser fallback widened semantic allow; decision=%+v cmd=%+v", decision, parsed)
	}
	if decision.Outcome != "ask" {
		t.Fatalf("Outcome = %q, want ask; decision=%+v", decision.Outcome, decision)
	}
}

func TestPermissionPrecedenceProperties(t *testing.T) {
	gitStatus := PermissionRuleSpec{Command: PermissionCommandSpec{Name: "git", Semantic: &SemanticMatchSpec{Verb: "status"}}}
	tests := []struct {
		name       string
		permission PermissionSpec
		command    string
		want       string
	}{
		{name: "deny beats ask", permission: PermissionSpec{Deny: []PermissionRuleSpec{gitStatus}, Ask: []PermissionRuleSpec{gitStatus}}, command: "git status", want: "deny"},
		{name: "deny beats allow", permission: PermissionSpec{Deny: []PermissionRuleSpec{gitStatus}, Allow: []PermissionRuleSpec{gitStatus}}, command: "git status", want: "deny"},
		{name: "ask beats allow", permission: PermissionSpec{Ask: []PermissionRuleSpec{gitStatus}, Allow: []PermissionRuleSpec{gitStatus}}, command: "git status", want: "ask"},
		{name: "all abstain stays abstain", permission: PermissionSpec{Allow: []PermissionRuleSpec{{Command: PermissionCommandSpec{Name: "git", Semantic: &SemanticMatchSpec{Verb: "diff"}}}}}, command: "git status", want: "abstain"},
		{name: "unsafe fallback asks", permission: PermissionSpec{Allow: []PermissionRuleSpec{gitStatus}}, command: "git status > /tmp/out", want: "ask"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPipeline(PipelineSpec{Permission: tt.permission}, Source{})
			got, err := Evaluate(p, tt.command)
			if err != nil {
				t.Fatalf("Evaluate() error = %v", err)
			}
			if got.Outcome != tt.want {
				t.Fatalf("Outcome = %q, want %q; decision=%+v", got.Outcome, tt.want, got)
			}
		})
	}
}

func TestSecurityRegressionMatrixParserBoundaries(t *testing.T) {
	tests := []struct {
		name           string
		command        string
		registry       *commandpkg.CommandParserRegistry
		permission     PermissionSpec
		want           string
		wantParser     string
		wantSemantic   string
		wantActionPath []string
	}{
		{
			name:           "git semantic parser preserves deny",
			command:        "git status",
			registry:       commandpkg.DefaultParserRegistry(),
			permission:     PermissionSpec{Deny: []PermissionRuleSpec{{Command: PermissionCommandSpec{Name: "git", Semantic: &SemanticMatchSpec{Verb: "status"}}}}, Allow: []PermissionRuleSpec{{Command: PermissionCommandSpec{Name: "git"}}}},
			want:           "deny",
			wantParser:     "git",
			wantSemantic:   "git",
			wantActionPath: []string{"status"},
		},
		{
			name:         "generic fallback asks instead of widening to command allow",
			command:      "git -C repo status",
			registry:     commandpkg.NewCommandParserRegistry(),
			permission:   PermissionSpec{Deny: []PermissionRuleSpec{{Command: PermissionCommandSpec{Name: "git", Semantic: &SemanticMatchSpec{Verb: "status"}}}}, Allow: []PermissionRuleSpec{{Command: PermissionCommandSpec{Name: "git"}}}},
			want:         "ask",
			wantParser:   "generic",
			wantSemantic: "",
		},
		{
			name:         "unknown command asks by default",
			command:      "unknown-tool status",
			registry:     commandpkg.DefaultParserRegistry(),
			permission:   PermissionSpec{Allow: []PermissionRuleSpec{{Command: PermissionCommandSpec{Name: "git"}}}},
			want:         "ask",
			wantParser:   "generic",
			wantSemantic: "",
		},
	}

	for _, tt := range tests {
		t.Run("parser/"+tt.name, func(t *testing.T) {
			plan := commandpkg.ParseWithRegistry(tt.command, tt.registry)
			if plan.Shape.Kind != commandpkg.ShellShapeSimple {
				t.Fatalf("Shape.Kind = %q, want simple; plan=%+v", plan.Shape.Kind, plan)
			}
			if len(plan.Commands) != 1 {
				t.Fatalf("len(Commands) = %d, want 1; plan=%+v", len(plan.Commands), plan)
			}
			cmd := plan.Commands[0]
			if cmd.Parser != tt.wantParser || cmd.SemanticParser != tt.wantSemantic {
				t.Fatalf("parser=(%q,%q), want (%q,%q); cmd=%+v", cmd.Parser, cmd.SemanticParser, tt.wantParser, tt.wantSemantic, cmd)
			}
			if len(tt.wantActionPath) > 0 && !sameStrings(cmd.ActionPath, tt.wantActionPath) {
				t.Fatalf("ActionPath=%#v, want %#v", cmd.ActionPath, tt.wantActionPath)
			}

			p := NewPipeline(PipelineSpec{Permission: tt.permission}, Source{})
			decision := evaluatePreparedCommand(p.prepared.Deny, p.prepared.Ask, p.prepared.Allow, cmd)
			if decision.Outcome != tt.want {
				t.Fatalf("Outcome = %q, want %q; decision=%+v cmd=%+v", decision.Outcome, tt.want, decision, cmd)
			}
			if tt.want == "ask" && decision.Outcome == "allow" {
				t.Fatalf("unsafe parser fallback widened to allow; decision=%+v", decision)
			}
		})
	}
}

func TestSecurityRegressionMatrixEvaluationNormalizationBoundaries(t *testing.T) {
	tests := []struct {
		name       string
		command    string
		permission PermissionSpec
		want       string
		wantCmd    string
		shape      commandpkg.ShellShapeKind
		trace      []securityTraceWant
	}{
		{
			name:    "shell dash c evaluates inner command without rewriting",
			command: "bash -c 'git status'",
			permission: PermissionSpec{Allow: []PermissionRuleSpec{{
				Command: PermissionCommandSpec{Name: "git", Semantic: &SemanticMatchSpec{Verb: "status"}},
			}}},
			want:    "allow",
			wantCmd: "bash -c 'git status'",
			shape:   commandpkg.ShellShapeSimple,
			trace:   []securityTraceWant{{name: "normalized_command"}, {effect: "allow"}},
		},
		{
			name:    "aws profile flag is semantic and raw deny still sees original command",
			command: "aws --profile dev sts get-caller-identity",
			permission: PermissionSpec{
				Deny: []PermissionRuleSpec{{Patterns: []string{`^\s*aws\s+--profile\s+dev\s+`}}},
				Allow: []PermissionRuleSpec{{
					Command: PermissionCommandSpec{Name: "aws", Semantic: &SemanticMatchSpec{Service: "sts", Profile: "dev"}},
				}},
			},
			want:    "deny",
			wantCmd: "aws --profile dev sts get-caller-identity",
			shape:   commandpkg.ShellShapeSimple,
			trace:   []securityTraceWant{{effect: "deny"}},
		},
		{
			name:    "wrapper plus aws profile flag is semantic but not env rewrite",
			command: "env aws --profile dev sts get-caller-identity",
			permission: PermissionSpec{Allow: []PermissionRuleSpec{{
				Command: PermissionCommandSpec{Name: "aws", Semantic: &SemanticMatchSpec{Service: "sts", Profile: "dev"}},
			}}},
			want:    "allow",
			wantCmd: "env aws --profile dev sts get-caller-identity",
			shape:   commandpkg.ShellShapeSimple,
			trace:   []securityTraceWant{{effect: "allow"}},
		},
		{
			name:    "unsafe shell dash c payload remains ask",
			command: "bash -c 'git status && rm -rf /tmp/x'",
			permission: PermissionSpec{Allow: []PermissionRuleSpec{{
				Command: PermissionCommandSpec{Name: "git", Semantic: &SemanticMatchSpec{Verb: "status"}},
			}}},
			want:    "ask",
			wantCmd: "bash -c 'git status && rm -rf /tmp/x'",
			shape:   commandpkg.ShellShapeCompound,
			trace:   []securityTraceWant{{effect: "ask"}},
		},
	}

	for _, tt := range tests {
		t.Run("evaluation/"+tt.name, func(t *testing.T) {
			p := NewPipeline(PipelineSpec{Permission: tt.permission}, Source{})
			got, err := Evaluate(p, tt.command)
			if err != nil {
				t.Fatalf("Evaluate() error = %v", err)
			}
			if got.Command != tt.wantCmd {
				t.Fatalf("Command = %q, want %q; decision=%+v", got.Command, tt.wantCmd, got)
			}
			assertSecurityDecision(t, got, tt.want, tt.shape, nil, tt.trace)
		})
	}
}

func assertSecurityDecision(t *testing.T, got Decision, wantOutcome string, wantShape commandpkg.ShellShapeKind, wantFlags []string, wantTrace []securityTraceWant) {
	t.Helper()
	if got.Outcome != wantOutcome {
		t.Fatalf("Outcome = %q, want %q; decision=%+v", got.Outcome, wantOutcome, got)
	}
	if len(got.Trace) == 0 {
		t.Fatalf("Trace is empty; decision=%+v", got)
	}
	plan := commandpkg.Parse(got.Command)
	if plan.Shape.Kind != wantShape {
		t.Fatalf("Shape.Kind = %q, want %q; command=%q decision=%+v plan=%+v", plan.Shape.Kind, wantShape, got.Command, got, plan)
	}
	for _, flag := range wantFlags {
		if !containsString(plan.Shape.Flags(), flag) {
			t.Fatalf("Shape.Flags() = %#v, want %q; command=%q decision=%+v", plan.Shape.Flags(), flag, got.Command, got)
		}
	}
	for _, want := range wantTrace {
		if !traceContains(got.Trace, want) {
			t.Fatalf("trace does not contain %+v; trace=%+v", want, got.Trace)
		}
	}
}

func traceContains(trace []TraceStep, want securityTraceWant) bool {
	for _, step := range trace {
		if want.name != "" && step.Name != want.name {
			continue
		}
		if want.effect != "" && step.Effect != want.effect {
			continue
		}
		return true
	}
	return false
}

func sameStrings(got []string, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
