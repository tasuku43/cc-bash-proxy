---
title: "Product Concept: Declarative Command Policy Guard"
status: implemented
date: 2026-04-18
---

# Product Concept

## 1. Purpose

This document defines the product-level purpose of `cmdproxy`: the problem it
solves, who it is for, and what value it should deliver.

## 2. Problem

Command execution policies for AI agents and shells are often implemented as
ad-hoc shell scripts. That approach is fast to start with, but it tends to
degrade over time:

1. Rules become hard to review and reason about
2. Small edits can silently widen what is allowed
3. The same policy is duplicated across Claude Code, shell hooks, and CI
4. Users get inconsistent deny messages depending on runtime integration

As a result, policy enforcement becomes fragile precisely where predictable
behavior is most important.

## 3. Primary Persona

**Operators of AI-assisted command execution**

- Use AI agents, shells, or CI that can trigger shell commands
- Need a local policy layer before command execution
- Want policies to be reviewable, testable, and portable across runtimes
- Prefer a small standalone CLI over runtime-specific shell glue

## 4. Operating Context

- `cmdproxy` runs as a local CLI, usually from a hook
- The caller provides stdin JSON describing an attempted command execution
- v1 evaluates command strings only; it does not inspect file writes, network
  fetches, or MCP calls
- The same rules should work across Claude Code, shell hooks, and CI

## 5. Core Value Proposition

`cmdproxy` should let users define command-deny policies declaratively and
verify those policies before rollout, while providing deterministic runtime
behavior when a command is denied.

## 6. Non-goals

1. Generating policies from natural language or LLM transcripts
2. Acting as a general authorization framework for all tool actions
3. Providing dashboards, telemetry products, or hosted policy management
4. Replacing runtime-specific integrations with a full adapter ecosystem in v1
