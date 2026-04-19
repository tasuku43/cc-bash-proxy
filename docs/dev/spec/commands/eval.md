---
title: "cmdproxy eval"
status: proposed
date: 2026-04-19
---

# cmdproxy eval

## Purpose

`cmdproxy eval` is the current hook entrypoint name. In the target architecture
it acts as the invocation-policy entrypoint that parses input, applies the first
matching directive, and returns pass / rewrite / reject / error.

The command name may change later if a better entrypoint name emerges, but this
document keeps `eval` as the current shell-facing contract.

## Input Sources

`cmdproxy eval` supports:

- generic execution input defined in `../core/INPUT_MODEL.md`
- Claude Code `PreToolUse` Bash payloads adapted into that execution model

Unsupported or malformed input is an error.

## Runtime Behavior

The target flow is:

1. Read stdin fully
2. Parse input JSON
3. Normalize supported input into a requested invocation
4. Load the effective config
5. Parse the invocation internally
6. Evaluate rules using first-match directive semantics
7. Emit output according to `../core/OUTPUT_CONTRACT.md`

## Implemented Rewrite Support

The current implementation already supports rewrite outcomes for:

- `rewrite.unwrap_shell_dash_c`
- `rewrite.move_flag_to_env`

If a rewrite primitive matches but cannot safely rewrite the invocation,
evaluation continues and the original command may still pass unless a later
`reject` rule matches.

## Notes

- The earlier deny-only implementation is transitional
- The target model is directive-driven, not allow / deny only
- Downstream permission systems remain the final execution authority
