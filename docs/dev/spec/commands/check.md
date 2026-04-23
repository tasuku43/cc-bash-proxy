---
title: "cc-bash-proxy check"
status: proposed
date: 2026-04-19
---

# cc-bash-proxy check

## Purpose

`cc-bash-proxy check` evaluates a single invocation interactively without requiring
stdin JSON from an external hook.

## Relationship To `hook`

`cc-bash-proxy check` is the interactive convenience wrapper over the same directive
application logic used by `cc-bash-proxy hook`.

- it accepts shell argv and reconstructs one command string internally
- it constructs the canonical execution request internally
- it applies the same parse, match, and directive flow
- it emits the same pass / rewrite / reject / error outcomes

For shell-sensitive examples, callers should prefer a single quoted command
string, for example:

```sh
cmdproxy check 'git -C repo status'
cmdproxy check 'bash -c '"'"'echo hello world'"'"''
```

This keeps quoting and nested arguments intact before `cmdproxy` rebuilds the
command string from CLI argv.

## Use Cases

- ad-hoc debugging while authoring rules
- checking whether a command would be rewritten
- confirming whether a command would be rejected
- observing the canonicalized form before relying on Claude Code hooks
