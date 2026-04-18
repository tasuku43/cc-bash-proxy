---
title: "Output Contract"
status: implemented
date: 2026-04-18
---

# Output Contract

## 1. Scope

This document defines exit codes and output behavior for `cmdguard eval` v1.

## 2. Exit Codes

`cmdguard eval` uses the following process exit codes:

- `0`: allow
- `2`: deny
- `1`: error

`deny` and `error` are intentionally distinct. A policy decision is not the
same as a malformed input or configuration failure.

## 3. Default Output Mode

The default output mode is human-readable text.

- Allow: no output
- Deny: write a concise explanation to `stderr`
- Error: write an error explanation to `stderr`

v1 should avoid writing human-readable decision text to `stdout` so callers can
use `stdout` safely in scripts.

## 4. Deny Output Requirements

In default human-readable mode, a deny response must include:

- The selected `rule_id`
- The configured deny `message`

The deny output may also include the evaluated command and source file path when
that adds useful context, but `rule_id` and `message` are the minimum contract.

## 5. JSON Output Mode

`cmdguard eval --format json` emits a single JSON object describing the result.

### Allow payload

```json
{
  "decision": "allow"
}
```

### Deny payload

```json
{
  "decision": "deny",
  "rule_id": "no-git-dash-c",
  "message": "git -C は禁止。cd で移動してから実行してください。",
  "command": "git -C repos/foo status",
  "source": {
    "layer": "user",
    "path": "/home/alice/.config/cmdguard/cmdguard.yml"
  }
}
```

### Error payload

```json
{
  "decision": "error",
  "error": {
    "code": "invalid_input",
    "message": "action must be exec"
  }
}
```

## 6. Error Classes

v1 error payloads should distinguish at least these classes:

- `invalid_input`: stdin JSON shape is unsupported or incomplete
- `invalid_config`: YAML parsing or schema validation failed
- `runtime_error`: unexpected internal failure

Exact wording may change, but the top-level error class should remain stable.
