---
title: "cmdproxy hook"
status: proposed
date: 2026-04-22
---

# cmdproxy hook

## Purpose

`cmdproxy hook claude` is the Claude Code hook entrypoint. It reads the
Claude Code `PreToolUse` Bash payload from stdin, applies the configured
rewrite and permission pipeline, and emits Claude Code hook JSON on stdout.

## Input Sources

`cmdproxy hook claude` supports:

- Claude Code `PreToolUse` Bash payloads

Unsupported or malformed input is converted into a deny response for the hook
caller.

## Runtime Behavior

The current flow is:

1. Read stdin fully
2. Parse Claude Code hook JSON
3. Normalize the Bash command into an invocation request
4. Load the verified artifact for the effective config
5. Evaluate the rewrite pipeline
6. Evaluate permissions on the rewritten command
7. Emit Claude Code hook JSON:
   - `allow`: `permissionDecision: "allow"`
   - `ask`: no `permissionDecision`, so Claude prompts
   - `deny`: `permissionDecision: "deny"`
   - `error`: deny response

## Implemented Rewrite Support

The current implementation already supports rewrite outcomes for:

- `move_flag_to_env`
- `move_env_to_flag`
- `unwrap_shell_dash_c`
- `unwrap_wrapper`
- `strip_command_path`

If a rewrite primitive matches but cannot safely rewrite the invocation,
evaluation continues with the current command.

## Permission Source Of Truth

`cmdproxy` is the source of truth for command permission evaluation.

Claude Code `settings.json` is no longer the primary place to express shell
command permission policy. The Claude hook only receives the final `allow`,
`ask`, or `deny` result produced by `cmdproxy`.

## RTK Integration

When `cmdproxy hook claude --rtk` is used, the runtime order is:

1. evaluate `cmdproxy` rewrite pipeline
2. evaluate `cmdproxy` permission pipeline
3. if not denied, apply the final `rtk` rewrite
4. emit the final `updatedInput.command`

This keeps permission decisions stable even when external Bash hooks are not
executed serially.
