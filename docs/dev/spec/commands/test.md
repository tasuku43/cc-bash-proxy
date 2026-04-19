---
title: "cmdproxy test"
status: proposed
date: 2026-04-19
---

# cmdproxy test

## Purpose

`cmdproxy test` validates that configured rules behave as claimed by their
embedded examples.

## Target Behavior

For every loaded rule, `cmdproxy test` should verify the rule's declared
directive behavior.

That means examples should be able to express at least:

- rewrite: input becomes a specific canonical output
- reject: input is blocked
- pass: input is not matched by the rule

## Scope

`cmdproxy test` is for rule-local verification.

It should:

- parse and validate the config
- validate matcher and directive payloads
- run each example against the rule's own matcher and directive

It should not:

- simulate full multi-rule downstream permission behavior
- replace dedicated integration tests for Claude Code hook semantics

## Exit Behavior

- `0`: all rule examples pass
- `1`: any example fails or configuration is invalid
