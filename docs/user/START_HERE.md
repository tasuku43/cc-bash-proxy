# Start Here

`cmdproxy` is a local CLI that sits in front of command execution and enforces
policy-approved invocation shape.

Its main job is to normalize command shape so downstream permission systems
evaluate the invocation you intended, not a drifted wrapper-heavy form.

## Quick Start

1. Create the user config

```sh
cmdproxy init
```

2. Edit `~/.config/cmdproxy/cmdproxy.yml`

3. Validate the config after each change

```sh
cmdproxy test
cmdproxy doctor --format json
```

4. Spot-check individual commands

```sh
cmdproxy check aws --profile read-only-profile s3 ls
cmdproxy check bash -c 'git status'
```

5. Register `cmdproxy eval` in your hook runner

## Claude Code

For Claude Code, add `cmdproxy eval` as a `PreToolUse` Bash hook.

```json
{
  "matcher": "Bash",
  "hooks": [
    { "type": "command", "command": "cmdproxy eval" }
  ]
}
```

If you also use another Bash hook such as `rtk hook claude`, place
`cmdproxy eval` first.

That ordering matters because `cmdproxy` should canonicalize or reject the
invocation before later hook-side processing and before Claude Code permissions
evaluate the final command shape.

## Current Rule Model

The current config format is still `version: 1`.

- rules use `match` or `pattern`
- rules use one directive: `rewrite` or `reject`
- examples are still written as `block_examples` and `allow_examples`

If you are contributing to the implementation, start from
`docs/dev/README.md` instead.
