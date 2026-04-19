---
title: "cmdproxy check"
status: implemented
date: 2026-04-18
---

# cmdproxy check

## Purpose

`cmdproxy check` evaluates a single command string interactively without
requiring stdin JSON from an external hook.

## Relationship to `eval`

`cmdproxy check` is a convenience wrapper over the same evaluation logic used by
`cmdproxy eval`.

- It accepts a command string as a CLI argument or flag
- It constructs the canonical execution input internally
- It applies the same config loading and first-match deny logic
- It uses the same output contract and exit codes as `eval`

## Use Cases

- ad-hoc debugging while authoring rules
- reproducing a deny decision outside the hook runtime
- confirming which rule ID would fire for a candidate command

## Output

By default, `cmdproxy check` should present the same decision shape as
`cmdproxy eval`.

- allow: no output, exit `0`
- deny: human-readable stderr, exit `2`
- error: human-readable stderr, exit `1`

JSON output support should mirror `eval` when `--format json` is used.
