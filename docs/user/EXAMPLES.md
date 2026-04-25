# Examples

These examples use the current permission shape: `command`, `env`, and
`patterns`.

## Git Read-Only Allow

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

## Git Force Push Deny

```yaml
permission:
  deny:
    - name: git force push
      command:
        name: git
        semantic:
          verb: push
          force: true
```

## AWS Identity Allow

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

## kubectl Read-Only Allow

```yaml
permission:
  allow:
    - name: kubectl read-only
      command:
        name: kubectl
        semantic:
          verb_in:
            - get
            - describe
```

## gh Read-Only PR Inspection

```yaml
permission:
  allow:
    - name: gh pr read-only
      command:
        name: gh
        semantic:
          area: pr
          verb_in:
            - view
            - list
            - diff
```

## helmfile Diff Allow

```yaml
permission:
  allow:
    - name: helmfile diff
      command:
        name: helmfile
        semantic:
          verb: diff
```

## Read-Only Shell Basics

```yaml
permission:
  allow:
    - name: read-only shell basics
      patterns:
        - "^ls(\\s|$)"
        - "^pwd$"
```

## Unknown Command Fallback

Use `patterns` when a command has no semantic schema.

```yaml
permission:
  ask:
    - name: tool preview
      patterns:
        - "^my-tool\\s+preview(\\s|$)"
```
