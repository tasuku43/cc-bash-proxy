---
title: "Start Here"
status: implemented
date: 2026-04-18
---

# Start Here

`cmdproxy` is a local CLI that denies shell commands when they match configured
regular-expression rules.

For v1, the core workflow is:

1. Write `~/.config/cmdproxy/cmdproxy.yml` with deny rules
2. Add `block_examples` and `allow_examples` for every rule
3. Run `cmdproxy test`
4. Integrate `cmdproxy eval` into Claude Code, shell hooks, or CI

If you are contributing to the implementation, start from `docs/dev/README.md`
instead.
