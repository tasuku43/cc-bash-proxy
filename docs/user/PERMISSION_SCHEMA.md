# Permission Schema

Permission policy is grouped into `deny`, `ask`, and `allow` buckets.

```yaml
permission:
  deny: []
  ask: []
  allow: []
```

Each rule can use `command`, `env`, and `patterns`.
Each rule may also set `message`; when that rule determines `allow`, `ask`, or
`deny`, the message is returned as Claude Code's permission decision reason.

## command

Use `command` to match a command by name. For commands with semantic parser
support, add `command.semantic`.

```yaml
permission:
  allow:
    - name: git read-only
      command:
        name: git
        semantic:
          verb_in:
            - status
            - diff
            - log
            - show
```

The semantic schema is selected by `command.name`. Inspect supported commands
with:

```sh
cc-bash-guard help semantic
cc-bash-guard semantic-schema --format json
```

## env

Use `env` to require or reject environment variables for the invocation.

```yaml
permission:
  allow:
    - name: AWS identity
      command:
        name: aws
        semantic:
          service: sts
          operation: get-caller-identity
      env:
        requires:
          - AWS_PROFILE
```

`env.requires` means the variable must be present. `env.missing` means the
variable must not be present.

## patterns

Use `patterns` for raw regular expression matching against the original command
string and parsed command elements. Shell `-c` wrappers are unwrapped for
evaluation, so a pattern such as `^aws(\s|$)` also matches
`bash -c 'aws s3 ls'`. This is the fallback for commands without semantic
support.

```yaml
permission:
  allow:
    - name: read-only shell basics
      patterns:
        - "^ls(\\s|$)"
        - "^pwd$"
```

## Valid Combinations

Rules can combine fields as follows:

- `command`
- `command` plus `env`
- `command` plus `semantic`
- `command` plus `semantic` plus `env`
- `patterns`
- `patterns` plus `env`
- `env`

Use semantic matching when a command is listed by `cc-bash-guard help semantic`.
Use `patterns` for commands without semantic support or when raw regex matching
is the intended policy.

## Evaluation Order

`cc-bash-guard` policy and Claude settings permissions are permission sources.
Each source returns `deny`, `ask`, `allow`, or `abstain`.

Decision precedence is:

```text
deny > ask > allow > abstain
```

`abstain` means no matching rule. The final fallback is `ask` only when all
permission sources abstain.

## Command Evaluation

`cc-bash-guard` evaluates commands but does not rewrite them.

Parser-backed normalization is evaluation-only:

- shell `-c` wrappers can be inspected as inner commands
- absolute command paths can match by basename
- command-specific parsers can expose semantic fields such as AWS profile,
  service, and operation

The command returned to Claude Code is not changed by permission evaluation.
