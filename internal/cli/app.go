package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tasuku43/cmdproxy/internal/config"
	"github.com/tasuku43/cmdproxy/internal/doctor"
	"github.com/tasuku43/cmdproxy/internal/domain/policy"
	"github.com/tasuku43/cmdproxy/internal/input"
)

const (
	exitAllow  = 0
	exitError  = 1
	exitReject = 2
)

type Streams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Env struct {
	Cwd           string
	Home          string
	XDGConfigHome string
	XDGCacheHome  string
}

func Run(args []string, streams Streams, env Env) int {
	if len(args) == 0 {
		writeUsage(streams.Stdout)
		return exitError
	}

	switch args[0] {
	case "eval":
		return runEval(args[1:], streams, env)
	case "check":
		return runCheck(args[1:], streams, env)
	case "test":
		return runTest(args[1:], streams, env)
	case "doctor":
		return runDoctor(args[1:], streams, env)
	case "init":
		return runInit(args[1:], streams, env)
	case "-h", "--help", "help":
		if len(args) > 1 {
			writeCommandHelp(streams.Stdout, args[1])
		} else {
			writeUsage(streams.Stdout)
		}
		return exitAllow
	default:
		writeErr(streams.Stderr, "unknown command: "+args[0])
		return exitError
	}
}

func runEval(args []string, streams Streams, env Env) int {
	if wantsHelp(args) {
		writeCommandHelp(streams.Stdout, "eval")
		return exitAllow
	}
	format, rest, err := parseCommonFlags(args)
	if err != nil || len(rest) != 0 {
		writeCommandHelp(streams.Stderr, "eval")
		return exitError
	}
	raw, err := io.ReadAll(streams.Stdin)
	if err != nil {
		return emitError(streams, format, "runtime_error", err.Error())
	}

	req, err := input.Normalize(raw)
	if err != nil {
		return emitError(streams, format, "invalid_input", err.Error())
	}
	return evaluateRequest(req, format, streams, env)
}

func runCheck(args []string, streams Streams, env Env) int {
	if wantsHelp(args) {
		writeCommandHelp(streams.Stdout, "check")
		return exitAllow
	}
	format, rest, err := parseCommonFlags(args)
	if err != nil || len(rest) == 0 {
		writeCommandHelp(streams.Stderr, "check")
		return exitError
	}
	req := input.ExecRequest{Action: "exec", Command: strings.Join(rest, " ")}
	return evaluateRequest(req, format, streams, env)
}

func runTest(args []string, streams Streams, env Env) int {
	if wantsHelp(args) {
		writeCommandHelp(streams.Stdout, "test")
		return exitAllow
	}
	if len(args) != 0 {
		writeCommandHelp(streams.Stderr, "test")
		return exitError
	}
	loaded := config.LoadEffective(env.Home, env.XDGConfigHome)
	if len(loaded.Errors) > 0 {
		for _, msg := range policy.ErrorStrings(loaded.Errors) {
			writeErr(streams.Stderr, msg)
		}
		return exitError
	}

	report := doctor.Run(loaded, env.Home)
	for _, check := range report.Checks {
		if check.ID == "rules.examples-pass" && check.Status == doctor.StatusFail {
			writeErr(streams.Stderr, check.Message)
			return exitError
		}
	}

	ruleCount := len(loaded.Rules)
	exampleCount := 0
	for _, r := range loaded.Rules {
		exampleCount += len(r.BlockExamples) + len(r.AllowExamples)
	}
	fmt.Fprintf(streams.Stdout, "ok: %d rules, %d examples checked\n", ruleCount, exampleCount)
	return exitAllow
}

func runDoctor(args []string, streams Streams, env Env) int {
	if wantsHelp(args) {
		writeCommandHelp(streams.Stdout, "doctor")
		return exitAllow
	}
	format, rest, err := parseCommonFlags(args)
	if err != nil || len(rest) != 0 {
		writeCommandHelp(streams.Stderr, "doctor")
		return exitError
	}
	loaded := config.LoadEffective(env.Home, env.XDGConfigHome)
	report := doctor.Run(loaded, env.Home)

	if format == "json" {
		enc := json.NewEncoder(streams.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(report); err != nil {
			writeErr(streams.Stderr, err.Error())
			return exitError
		}
	} else {
		for _, check := range report.Checks {
			fmt.Fprintf(streams.Stdout, "[%s] %s: %s\n", strings.ToUpper(string(check.Status)), check.ID, check.Message)
		}
	}

	if doctor.HasFailures(report) {
		return exitError
	}
	return exitAllow
}

func runInit(args []string, streams Streams, env Env) int {
	if wantsHelp(args) {
		writeCommandHelp(streams.Stdout, "init")
		return exitAllow
	}
	if len(args) != 0 {
		writeCommandHelp(streams.Stderr, "init")
		return exitError
	}
	configDir := filepath.Join(userConfigBase(env.Home, env.XDGConfigHome), "cmdproxy")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		writeErr(streams.Stderr, err.Error())
		return exitError
	}
	configPath := filepath.Join(configDir, "cmdproxy.yml")
	created := false
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(configPath, []byte(starterConfig), 0o644); err != nil {
			writeErr(streams.Stderr, err.Error())
			return exitError
		}
		created = true
	}

	if created {
		fmt.Fprintf(streams.Stdout, "created %s\n", configPath)
	} else {
		fmt.Fprintf(streams.Stdout, "exists %s\n", configPath)
	}
	fmt.Fprintf(streams.Stdout, "user config: %s\n", configPath)

	claudeSettings := filepath.Join(env.Home, ".claude", "settings.json")
	if _, err := os.Stat(claudeSettings); err == nil {
		fmt.Fprintf(streams.Stdout, "detected Claude Code settings: %s\n", claudeSettings)
	} else {
		fmt.Fprintf(streams.Stdout, "Claude Code settings not found: %s\n", claudeSettings)
	}

	fmt.Fprintln(streams.Stdout, "hook snippet:")
	fmt.Fprintln(streams.Stdout, `{"matcher":"Bash","hooks":[{"type":"command","command":"cmdproxy eval"}]}`)
	return exitAllow
}

func evaluateRequest(req input.ExecRequest, format string, streams Streams, env Env) int {
	loaded := config.LoadEffectiveForEval(env.Home, env.XDGConfigHome, env.XDGCacheHome)
	if len(loaded.Errors) > 0 {
		return emitError(streams, format, "invalid_config", strings.Join(policy.ErrorStrings(loaded.Errors), "; "))
	}

	decision, err := policy.Evaluate(loaded.Rules, req.Command)
	if err != nil {
		return emitError(streams, format, "runtime_error", err.Error())
	}
	if decision.Outcome == "pass" {
		if format == "json" {
			_ = json.NewEncoder(streams.Stdout).Encode(map[string]any{
				"decision": "pass",
				"command":  decision.Command,
			})
		}
		return exitAllow
	}

	if decision.Outcome == "rewrite" {
		if format == "json" {
			payload := map[string]any{
				"decision":         "rewrite",
				"rule_id":          decision.Rule.ID,
				"command":          decision.Command,
				"original_command": decision.OriginalCommand,
				"source":           decision.Rule.Source,
			}
			_ = json.NewEncoder(streams.Stdout).Encode(payload)
		} else {
			fmt.Fprintln(streams.Stdout, decision.Command)
		}
		return exitAllow
	}

	if format == "json" {
		payload := map[string]any{
			"decision": "reject",
			"rule_id":  decision.Rule.ID,
			"message":  decision.Rule.RejectMessage(),
			"command":  decision.Command,
			"source":   decision.Rule.Source,
		}
		_ = json.NewEncoder(streams.Stdout).Encode(payload)
	} else {
		fmt.Fprintf(streams.Stderr, "[%s] %s\n", decision.Rule.ID, decision.Rule.RejectMessage())
	}
	return exitReject
}

func emitError(streams Streams, format string, code string, message string) int {
	if format == "json" {
		_ = json.NewEncoder(streams.Stdout).Encode(map[string]any{
			"decision": "error",
			"error": map[string]string{
				"code":    code,
				"message": message,
			},
		})
	} else {
		writeErr(streams.Stderr, message)
	}
	return exitError
}

func parseCommonFlags(args []string) (string, []string, error) {
	format := ""
	rest := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--format":
			if i+1 >= len(args) {
				return "", nil, errors.New("missing --format value")
			}
			format = args[i+1]
			i++
		default:
			rest = append(rest, args[i])
		}
	}
	if format != "" && format != "json" {
		return "", nil, fmt.Errorf("unsupported format %q", format)
	}
	return format, rest, nil
}

func writeUsage(w io.Writer) {
	fmt.Fprint(w, `cmdproxy

Declarative, testable command policy for AI-agent shell commands.

Typical workflow:
  1. Edit ~/.config/cmdproxy/cmdproxy.yml
  2. Add block_examples and allow_examples for every rule
  3. Run cmdproxy test
  4. Use cmdproxy check for spot checks
  5. Let Claude Code call cmdproxy eval from PreToolUse

Usage:
  cmdproxy <command> [flags]

Commands:
  init     create the user config and print the Claude Code hook snippet
  test     validate every rule example; this is the main authoring command
  check    evaluate one command string interactively
  doctor   inspect config quality and installation state
  eval     hook entrypoint used by Claude Code and other callers

Help:
  cmdproxy help <command>
  cmdproxy <command> --help

Examples:
  cmdproxy init
  cmdproxy test
  cmdproxy check --format json 'git -C repo status'
  cmdproxy doctor --format json
`)
}

func writeCommandHelp(w io.Writer, command string) {
	switch command {
	case "init":
		fmt.Fprint(w, `cmdproxy init

Create ~/.config/cmdproxy/cmdproxy.yml when it does not exist and print the
Claude Code PreToolUse hook snippet.

Usage:
  cmdproxy init

Typical use:
  cmdproxy init
`)
	case "test":
		fmt.Fprint(w, `cmdproxy test

Validate every rule in ~/.config/cmdproxy/cmdproxy.yml.
This is the main command to run after editing rules.

Usage:
  cmdproxy test

What it checks:
  - every block_examples entry matches its rule matcher
  - every allow_examples entry does not match its rule matcher

Typical use:
  $EDITOR ~/.config/cmdproxy/cmdproxy.yml
  cmdproxy test
`)
	case "check":
		fmt.Fprint(w, `cmdproxy check

Evaluate one command string against the current rule set.
Use this while authoring rules before relying on Claude Code hooks.

Usage:
  cmdproxy check [--format json] <command>

Examples:
  cmdproxy check 'git -C repo status'
  cmdproxy check --format json 'AWS_PROFILE=read-only-profile aws s3 ls'
`)
	case "doctor":
		fmt.Fprint(w, `cmdproxy doctor

Inspect config validity, rule quality, and Claude Code hook registration.

Usage:
  cmdproxy doctor [--format json]

Examples:
  cmdproxy doctor
  cmdproxy doctor --format json
`)
	case "eval":
		fmt.Fprint(w, `cmdproxy eval

Hook entrypoint for Claude Code and other callers.
Reads stdin JSON and returns allow, deny, or error.

Usage:
  cmdproxy eval [--format json]

Note:
  You usually do not run this manually. Edit rules and use cmdproxy test or
  cmdproxy check instead.
`)
	default:
		writeUsage(w)
	}
}

func wantsHelp(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func writeErr(w io.Writer, msg string) {
	fmt.Fprintln(w, msg)
}

func userConfigBase(home string, xdgConfigHome string) string {
	if xdgConfigHome != "" {
		return xdgConfigHome
	}
	return filepath.Join(home, ".config")
}

const starterConfig = `version: 1
rules:
  - id: no-git-dash-c
    match:
      command: git
      args_contains:
        - "-C"
    message: "git -C is blocked. Change into the target directory and rerun the command."
    block_examples:
      - "git -C repos/foo status"
    allow_examples:
      - "git status"
      - "# git -C in comment"
`
