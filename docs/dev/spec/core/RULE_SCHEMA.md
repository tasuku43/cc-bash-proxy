---
title: "Rule Schema"
status: proposed
date: 2026-04-19
---

# Rule Schema

## 1. Scope

This document defines the target directive-based YAML schema for `cmdproxy`.

## 2. Transitional Implemented Shape

The currently implemented on-disk configuration shape is:

```yaml
version: 1
rules:
  - id: aws-profile-to-env
    match:
      command: aws
      args_contains:
        - "--profile"
    rewrite:
      move_flag_to_env:
        flag: "--profile"
        env: "AWS_PROFILE"
    block_examples:
      - "aws --profile prod s3 ls"
    allow_examples:
      - "AWS_PROFILE=prod aws s3 ls"

  - id: no-shell-dash-c
    match:
      command_in: ["bash", "sh", "zsh", "dash", "ksh"]
      args_contains: ["-c"]
    reject:
      message: "shell -c must not pass through unchanged."
    block_examples:
      - "bash -c 'git status && git diff'"
    allow_examples:
      - "git status"
```

This transitional shape keeps `version: 1` and still uses
`block_examples` / `allow_examples`.

## 3. Target Shape

The target configuration shape is:

```yaml
version: 2
rules:
  - id: aws-profile-to-env
    match:
      command: aws
      args_contains:
        - "--profile"
    rewrite:
      move_flag_to_env:
        flag: "--profile"
        env: "AWS_PROFILE"
    examples:
      - in: "aws --profile prod s3 ls"
        out: "AWS_PROFILE=prod aws s3 ls"

  - id: no-shell-dash-c
    match:
      command_in: ["bash", "sh", "zsh", "dash", "ksh"]
      args_contains: ["-c"]
    reject:
      message: "shell -c must not pass through unchanged."
    examples:
      - in: "bash -c 'git status && git diff'"
        reject: true
```

## 4. Top-Level Fields

- `version`: required integer, currently `1`, target value `2`
- `rules`: required non-empty array of rule objects

Unknown top-level keys are invalid.

## 5. Rule Fields

Each rule object must contain:

- `id`: required string
- exactly one of `match` or `pattern`
- exactly one of `rewrite` or `reject`
- examples, currently as `block_examples` and `allow_examples`

Unknown rule-level keys are invalid.

## 6. Matcher Fields

### `match`

`match` is the preferred matcher model.

Target fields:

- `command`
- `command_in`
- `subcommand`
- `args_contains`
- `args_prefixes`
- `env_requires`
- `env_missing`

The matcher operates on `cmdproxy`'s internal normalized invocation model.

### `pattern`

`pattern` remains available as an escape hatch for invocation shapes that are
not yet well represented by structured matchers.

- Must compile as Go RE2
- Matches against the raw command string
- Should be used sparingly where structured matching is insufficient

## 7. Directive Fields

### `rewrite`

`rewrite` contains a typed rewrite primitive.

Initial target primitives are intentionally narrow, for example:

- `move_flag_to_env`
- `unwrap_shell_dash_c`
- `strip_wrapper`

The currently implemented primitives are:

- `move_flag_to_env`
- `unwrap_shell_dash_c`

Each primitive should have a dedicated structured payload. Free-form string
templates are out of scope.

### `reject`

`reject` contains:

- `message`: required string

Future metadata is possible, but the minimal contract is a stable human-readable
explanation.

## 8. Examples

Examples validate rule intent at the directive level.

Today, examples are stored as:

- `block_examples`
- `allow_examples`

Each example should describe one of:

- `in` + `out` for rewrite behavior
- `in` + `reject: true` for reject behavior
- `in` + `pass: true` when a non-match example is useful

The exact final example schema may still evolve, but examples remain mandatory.

## 9. Validation Model

Validation is strict and aggregate.

- parsing should report all discovered schema issues in one run
- invalid matcher combinations are validation errors
- invalid directive payloads are validation errors
- missing examples are validation errors
- empty or ambiguous rules are validation errors

## 10. Out Of Scope

The following remain out of scope for the target initial model:

- arbitrary shell templating
- user-defined rewrite plugins
- implicit multi-step rewrite pipelines within one rule
- remote includes or hosted policy packs
