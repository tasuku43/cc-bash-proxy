---
title: "cmdproxy test"
status: proposed
date: 2026-04-22
---

# cmdproxy test

## Purpose

`cmdproxy test` validates that the configured pipeline behaves as claimed by
its embedded examples.

## Target Behavior

`cmdproxy test` should validate three test layers:

- rewrite-step-local tests
- permission-rule-local tests
- top-level end-to-end tests

That means examples should be able to express at least:

- rewrite: input becomes a specific canonical output
- permission: rewritten input matches or does not match a specific effect bucket
- E2E: input becomes an optional rewritten command and a final `allow`, `ask`,
  or `deny` decision

## Scope

`cmdproxy test` is for local pipeline verification.

It should:

- parse and validate the config
- validate matcher, rewrite, and permission payloads
- run each rewrite-step-local test
- run each permission-rule-local test
- run each top-level E2E test

It should not:

- depend on Claude Code `settings.json`
- replace dedicated integration tests for hook transport behavior

## Exit Behavior

- `0`: all configured tests pass
- `1`: any test fails or configuration is invalid
