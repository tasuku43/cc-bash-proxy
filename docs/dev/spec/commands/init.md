---
title: "cmdproxy init"
status: implemented
date: 2026-04-18
---

# cmdproxy init

## Purpose

`cmdproxy init` bootstraps a local `cmdproxy` setup without silently modifying
existing user configuration in unsafe ways.

## v1 Responsibilities

`cmdproxy init` may:

- create a starter user-wide config when one does not exist
- explain where the user-wide config lives
- detect compatible Claude Code settings files
- print the hook snippet needed to register `cmdproxy eval`

## Safety Principle

v1 `init` should optimize for idempotence and non-destructive setup.

- If the user-wide config already exists, do not overwrite it
- If Claude Code hook registration already exists or settings are non-trivial,
  prefer reporting status and proposed changes over blind mutation
- If automatic mutation is supported, it should be conservative and explicit

## Recommended Starter Config

The starter config should:

- use schema version `1`
- include at least one sample deny rule
- demonstrate both `block_examples` and `allow_examples`
- be valid under `cmdproxy test`

## Output

Default output should clearly separate:

- what was created
- what was detected
- what still requires manual action

This is especially important because `init` often runs once, long before a user
debugs hook behavior later.
