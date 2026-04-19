# cmdproxy

Declarative, testable command-string policy engine for AI agents and shells.

> **Status:** v1 core implementation in progress. See
> [`docs/dev/spec/README.md`](docs/dev/spec/README.md) for the v1.0 implementation contracts and
> [`docs/concepts/product-concept.md`](docs/concepts/product-concept.md) for the
> product concept.

## What it does

`cmdproxy` is a tiny hook that decides whether a shell command is allowed to
run. It is called from Claude Code `PreToolUse`, `zsh` `preexec`,
`pre-commit`, CI, or anywhere else a command-string policy is useful.

Rules are declared in YAML. Every rule ships with block/allow examples, and
`cmdproxy test` runs them as unit tests — so a rule change that would let
through a command it used to block fails CI, not production.

```yaml
# ~/.config/cmdproxy/cmdproxy.yml
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
      - "# git -C in comment"
```

Rules may use either:

- `match`: structured predicate matching, recommended for new rules
- `pattern`: a raw RE2 regular expression, kept as an escape hatch

## Non-goals

- LLM-assisted rule authoring and transcript mining live in a separate
  `cmdproxy-claude-plugin` repository, so the core CLI has no LLM
  dependency.
- Non-`exec` action types (`write`, `fetch`, `mcp_call`) are post-v1.

See [`docs/README.md`](docs/README.md) for the current documentation map.

## Install

Not yet released. Once v1 ships:

```sh
brew install tasuku43/tap/cmdproxy
# or
go install github.com/tasuku43/cmdproxy/cmd/cmdproxy@latest
```

## Setup

### 1. Initialize the user config

```sh
cmdproxy init
```

This creates the default rule file at:

```text
~/.config/cmdproxy/cmdproxy.yml
```

### 2. Edit rules

Update `~/.config/cmdproxy/cmdproxy.yml` directly.
For every rule, keep both:

- `block_examples`
- `allow_examples`

### 3. Validate changes

Run the main authoring command after every rule edit:

```sh
cmdproxy test
```

Use `cmdproxy check` for spot checks against concrete commands:

```sh
cmdproxy check --format json 'git -C repo status'
cmdproxy check --format json 'AWS_PROFILE=read-only-profile aws s3 ls'
```

## Claude Code Hook Setup

Register `cmdproxy eval` as a `PreToolUse` hook for `Bash`.

Example:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "cmdproxy eval" }
        ]
      }
    ]
  }
}
```

### Hook ordering with other tools

If you also use another `PreToolUse` Bash hook such as `rtk hook claude`,
register `cmdproxy eval` first.

Recommended order:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "cmdproxy eval" }
        ]
      },
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "rtk hook claude" }
        ]
      }
    ]
  }
}
```

This keeps `cmdproxy` as the first deny gate before other hook-side behavior
runs.

## License

MIT.
