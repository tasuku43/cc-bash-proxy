---
title: "cmdguard eval"
status: implemented
date: 2026-04-18
---

# cmdguard eval

## Purpose

`cmdguard eval` is the hook entrypoint for v1. It reads stdin JSON describing
an attempted command execution, evaluates loaded rules, and exits with the
result.

## Input Sources

`cmdguard eval` supports:

- Generic execution input defined in `../core/EVALUATION.md`
- Claude Code `PreToolUse` Bash payloads adapted into the generic execution
  model

Unsupported or malformed input is an error.

## Evaluation Behavior

`cmdguard eval` must:

1. Read stdin fully
2. Parse input JSON
3. Normalize supported input into an execution command string
4. Load configuration layers
5. Evaluate rules using first-match deny semantics
6. Emit output according to `../core/OUTPUT_CONTRACT.md`
7. Exit with `0`, `2`, or `1`

## Flags

- `--format json`: emit machine-readable output

Additional debug or trace flags are post-v1 unless they are specified
separately.
