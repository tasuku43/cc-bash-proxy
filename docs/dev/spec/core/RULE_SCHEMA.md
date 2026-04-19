---
title: "Rule Schema"
status: implemented
date: 2026-04-18
---

# Rule Schema

## 1. Scope

This document defines the v1 YAML schema for `cmdproxy` rule files.

## 2. Top-Level Shape

The v1 configuration file is a single YAML document with this shape:

```yaml
version: 1
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
```

## 3. Top-Level Fields

- `version`: required, integer, must be `1`
- `rules`: required, non-empty array of rule objects

Unknown top-level keys are invalid in v1.

## 4. Rule Fields

Each rule object must contain:

- `id`: required string
- exactly one of `match` or `pattern`
- `message`: required string
- `block_examples`: required non-empty array of strings
- `allow_examples`: required non-empty array of strings

Unknown rule-level keys are invalid in v1.

## 5. Field Constraints

### `id`

- Must match: `[a-z0-9][a-z0-9-]*`
- Must be unique across all loaded layers
- Should remain stable over time so tests and runtime output can refer to it

### `match`

`match` is the recommended matcher for new rules.

Supported fields in v1:

- `command`: exact executable basename match
- `command_in`: executable basename must be one of these values
- `subcommand`: first argument after the executable must match exactly
- `args_contains`: each listed argument token must be present exactly
- `args_prefixes`: each listed prefix must match at least one argument token
- `env_requires`: each listed environment variable must be present in command
  prefixes such as `AWS_PROFILE=dev aws s3 ls`
- `env_missing`: each listed environment variable must be absent from command
  prefixes

The matcher is evaluated against `cmdproxy`'s internal command parsing, not
against the raw input string directly.

### `pattern`

- Must compile as a Go RE2 regular expression
- Is evaluated against the full command string as provided by the caller
- Is kept as an escape hatch when `match` is not expressive enough
- v1 does not define capture-group semantics or replacement behavior

### `message`

- Must be a non-empty human-readable string
- Should explain both the reason for the deny and the preferred alternative
- v1 does not attempt to lint natural-language quality heuristically

### `block_examples`

- Must contain at least one example
- Every example must match the rule's matcher

### `allow_examples`

- Must contain at least one example
- Every example must not match the rule's matcher

## 6. Validation Model

Validation is strict and aggregate.

- Parsing should report all discovered schema issues in one run
- Regex compilation failures are validation errors
- Invalid or empty `match` objects are validation errors
- Missing required fields are validation errors
- Empty example arrays are validation errors

The goal is to make rule authoring reviewable and safe before runtime.

## 7. What v1 Does Not Support

The following are intentionally out of scope for v1:

- `allow` rules
- explicit exception clauses such as `except` or `unless`
- rule priority fields
- file includes or remote rule packs
- structured message metadata
- custom parser plugins or external matcher hooks

If these capabilities are needed later, they should be added as new schema
features rather than inferred from v1 fields.
