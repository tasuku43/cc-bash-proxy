# AWS Guidelines

Prefer environment-prefixed AWS profiles in project guidance:

```sh
AWS_PROFILE=myprof aws eks list-clusters
```

Discourage profile flags as the normal project style:

```sh
aws --profile myprof eks list-clusters
```

`cc-bash-guard` does not convert one form to another. It evaluates the command
that was requested.

The AWS parser can still evaluate AWS `profile`, `service`, and `operation`
semantically. For example:

```yaml
permission:
  allow:
    - name: AWS identity for configured profile
      command:
        name: aws
        semantic:
          service: sts
          operation: get-caller-identity
      env:
        requires:
          - AWS_PROFILE
```

Teams should document their preferred AWS command style in project guidance and
write permission policy that matches that style.
