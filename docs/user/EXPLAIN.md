# Explain Output

`cc-bash-guard explain` inspects a command without executing it. It uses the
same verified artifact as `cc-bash-guard hook`, so run `cc-bash-guard verify`
after editing policy or included files.

```sh
cc-bash-guard explain "bash -c 'git status'"
cc-bash-guard explain --format json "git push --force origin main"
cc-bash-guard explain --why-not allow "git status > /tmp/out"
```

## What To Look For

- parsed shell shape: how the shell command was classified
- shape flags: parser-derived flags such as redirects or unsafe shell shapes
- normalized command names: for example `/usr/bin/git` matching `git`
- evaluated inner command: for supported shell `-c` wrappers
- semantic fields: command-specific fields such as `git` `verb` and `force`
- policy outcome: cc-bash-guard policy result, including `abstain`
- Claude settings outcome: native Claude Code permission contribution
- final outcome: merged hook decision
- trace: why each decision step happened

`abstain` means "no matching rule" or "no opinion" for one permission source.
The final hook output never remains `abstain`: if cc-bash-guard policy and
Claude settings both abstain, the final decision falls back to `ask`.

## Targeted Why-Not Mode

Use `--why-not allow|ask|deny` when you expected a specific outcome and want a
direct explanation:

```sh
cc-bash-guard explain --why-not allow "git status > /tmp/out"
cc-bash-guard explain --why-not deny "git push origin main"
cc-bash-guard explain --format json --why-not allow "bash -c 'git status'"
```

Normal `explain` output is unchanged unless `--why-not` is passed. Why-not mode
reports:

- requested outcome
- actual cc-bash-guard policy outcome
- actual Claude settings outcome
- actual final outcome
- matched rule, if any
- command shape, shape flags, parser, and semantic fields
- concise reasons why the requested outcome did not happen
- safe suggestions such as running `verify`, adding a rule with `suggest`, or
  reviewing a higher-priority `deny` or `ask` rule

The JSON output is shaped for agents:

```json
{
  "command": "git diff",
  "requested_outcome": "allow",
  "actual": {
    "policy": "abstain",
    "claude_settings": "abstain",
    "final": "ask"
  },
  "reasons": [
    {
      "kind": "no_policy_match",
      "message": "cc-bash-guard policy abstained"
    },
    {
      "kind": "fallback_ask",
      "message": "all permission sources abstained; final fallback is ask"
    }
  ],
  "suggestions": [
    {
      "kind": "add_policy_rule",
      "message": "Use cc-bash-guard suggest to generate a starter rule"
    }
  ]
}
```

## Common Reads

When policy allows and Claude settings have no matching rule:

```text
policy: allow
claude_settings: abstain
final: allow
```

When neither policy nor Claude settings has a matching rule:

```text
policy: abstain
claude_settings: abstain
final: ask
```

This is the expected conservative fallback. Add a semantic `allow`, `ask`, or
`deny` rule plus tests when you want the command to have a more specific
decision.

## Shell Wrappers

For supported shell `-c` wrappers, `explain` shows that the inner command was
evaluated for policy. For example, `bash -c 'git status'` can match the same
semantic `git` rule as `git status`.

This is not command rewriting. The default hook still passes the original
command through unchanged.
