---
title: "cmdproxy test"
status: implemented
date: 2026-04-18
---

# cmdproxy test

## Purpose

`cmdproxy test` validates that configured rules behave as claimed by their
embedded examples.

## Core Behavior

For every loaded rule, `cmdproxy test` must verify:

- every `block_examples` entry matches the rule's matcher
- every `allow_examples` entry does not match the rule's matcher

If all checks pass, the command exits successfully.

## Exit Behavior

- `0`: all rule examples pass
- `1`: any example fails, configuration is invalid, or runtime execution fails

`cmdproxy test` does not use the `deny` exit code because it is not evaluating a
live command decision.

## Scope of Verification

v1 `cmdproxy test` verifies rule-local claims, not full merged runtime behavior.

Specifically, it does:

- parse and validate loaded config
- validate rule matchers
- check example truth against the rule's own matcher

It does not:

- simulate first-match selection across multiple rules
- fail on shadowing between rules
- evaluate external adapter payloads

Those concerns belong to `doctor` or dedicated integration tests.

## Output

Default output should be concise and human-readable.

- Success: summary of rules and examples checked
- Failure: each failing rule/example pair with the expected and actual result

Machine-readable output is post-v1 unless specified separately.
