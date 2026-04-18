---
title: "Configuration Model"
status: implemented
date: 2026-04-18
---

# Configuration Model

## 1. Scope

This document defines where `cmdguard` looks for configuration and how multiple
configuration layers interact in v1.

## 2. Supported Locations

v1 supports up to two configuration files:

1. Project-local: `.cmdguard.yml`
2. User-wide: `$XDG_CONFIG_HOME/cmdguard/cmdguard.yml`, or
   `~/.config/cmdguard/cmdguard.yml` by default

Each layer is optional. If neither file exists, evaluation proceeds with an
empty rule set.

## 3. Project-Local Resolution

Project-local configuration is resolved from the current working context of the
`cmdguard` process.

- If `.cmdguard.yml` exists in the current working directory, use it
- v1 does not search parent directories by default unless that behavior is
  specified separately in implementation docs

This keeps path resolution explicit and easy to reason about in hooks and CI.

## 4. Layer Precedence

Layer precedence is defined in `EVALUATION.md`.

- Project-local rules are evaluated before user-wide rules
- Within a file, source order is preserved

Precedence affects evaluation order, not schema differences. Both layers use the
same file format.

## 5. ID Collision Policy

Rule IDs must be unique across the effective configuration set.

- Duplicate IDs within one file are errors
- Duplicate IDs across layers are errors

v1 does not provide an override mechanism based on matching IDs.

## 6. Empty and Invalid States

- Missing file: allowed, treated as absent layer
- Empty file: invalid configuration
- Invalid YAML: invalid configuration
- Valid YAML with schema errors: invalid configuration

Invalid configuration causes `cmdguard eval` to exit with error rather than
silently falling back to partial policy enforcement.

## 7. Future Extensions

These are post-v1 concerns:

- `include:` directives
- rule packs
- explicit override semantics
- additional config layers such as repo-global or team-managed paths
