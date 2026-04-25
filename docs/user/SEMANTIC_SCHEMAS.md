# Semantic Schemas

Semantic matching lets policy match parsed command meaning instead of only raw
text. Semantic fields are command-specific, and the schema is selected by
`command.name`.

Inspect the current registry:

```sh
cc-bash-guard help semantic
cc-bash-guard help semantic git
cc-bash-guard semantic-schema --format json
cc-bash-guard semantic-schema git --format json
```

## Supported Commands

- `git`
- `aws`
- `kubectl`
- `gh`
- `helmfile`

The CLI list is generated from the semantic schema registry. Treat
`cc-bash-guard help semantic` and `semantic-schema` as the source of truth for
the installed binary.

## Field Types

Common field types are:

- `string`: exact value match
- `[]string`: all listed values are checked by the field-specific matcher
- `bool`: `true` or `false` for parser-recognized command properties

Boolean fields are parser-defined. For example, Git `force` is true for
`git push` when `--force`, `-f`, `--force-with-lease`, or
`--force-if-includes` is present. It is also true for `git clean` when `-f` or
`--force` is present.

`flags_contains` and `flags_prefixes` match parser-recognized option tokens.
They are not raw argv matchers.

## git

Use Git semantic fields for verbs, remotes, branches, refs, force-like options,
reset mode, clean mode, and diff staging.

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

Inspect all fields:

```sh
cc-bash-guard help semantic git
```

## aws

Use AWS semantic fields for service, operation, profile, region, endpoint, dry
run, and parser-recognized flags.

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

See also `docs/user/AWS_GUIDELINES.md`.

## kubectl

Use kubectl semantic fields for verb, resource, namespace, context, filename,
selector, container, dry run, force, recursion, and parser-recognized flags.

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

## gh

Use GitHub CLI semantic fields for `api`, `pr`, and `run` workflows.

```yaml
permission:
  deny:
    - name: mutating GitHub API
      command:
        name: gh
        semantic:
          area: api
          method_in:
            - POST
            - PATCH
            - PUT
            - DELETE
```

## helmfile

Use helmfile semantic fields for verb, environment, file, namespace,
kube-context, selector, dry run, wait behavior, and values options.

```yaml
permission:
  ask:
    - name: production helmfile destroy
      command:
        name: helmfile
        semantic:
          verb: destroy
          environment: prod
```
