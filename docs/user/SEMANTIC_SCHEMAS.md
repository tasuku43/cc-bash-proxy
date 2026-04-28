# Semantic Schemas

Semantic matching lets policy match parsed command meaning instead of only raw
text. Semantic fields are command-specific, and the schema is selected by
`command.name`.

Semantic fields live directly under `command.semantic`; `command.name` selects
which parser namespace validates those fields:

```yaml
permission:
  deny:
    - command:
        name: git
        semantic:
          verb: push
          force: true
```

Inspect the installed registry:

```sh
cc-bash-guard help semantic
cc-bash-guard help semantic git
cc-bash-guard semantic-schema --format json
cc-bash-guard semantic-schema git --format json
```

`cc-bash-guard help semantic` lists commands from the semantic schema registry.
`cc-bash-guard help semantic <command>` shows fields, types, descriptions,
boolean notes, and examples for one command. `semantic-schema --format json`
prints the same registry as machine-readable JSON.

## Supported Commands

The current registry supports:

- `git`
- `aws`
- `kubectl`
- `gh`
- `gws`
- `helmfile`
- `argocd`

Treat the CLI output as the source of truth for the installed binary. Commands
without a semantic schema should use `patterns`.

## Field Types

- `string`: exact value match.
- `[]string`: list values checked by the field-specific matcher.
- `bool`: `true` or `false` for parser-recognized command properties.

Unsupported semantic fields fail `verify` and include the supported fields for
that command. Unsupported value types fail `verify`; for example, `force:
"true"` is rejected because `force` must be a bool.

Semantic rules must set `command.name`; nested tool-name forms are rejected
because the command name is already the discriminator.

## Flags Fields

`flags_contains` and `flags_prefixes` match option tokens recognized by the
command-specific parser. They do not scan raw argv words, positional arguments,
or commands handled only by GenericParser fallback. Whether a token is available
for semantic flags therefore depends on the parser for the selected command.

## git

Use Git semantic fields for verbs, remotes, branches, refs, push force syntax,
reset mode, clean mode, diff staging, and parser-recognized flags.

Boolean fields:

- `force`: for `git push`, true only for `--force` or `-f`; for `git clean`,
  true for `-f` or `--force`.
- `force_with_lease`: for `git push`, true for `--force-with-lease`.
- `force_if_includes`: for `git push`, true for `--force-if-includes`.
- `hard`: true for `git reset --hard`.
- `recursive`: true for `git clean -d`.
- `include_ignored`: true for `git clean -x` or `--ignored`.
- `cached` and `staged`: true for `git diff --cached` or `--staged`.

To block every force-like push syntax, include all three force fields in deny
rules or write separate rules:

```yaml
permission:
  deny:
    - name: git destructive force push
      command:
        name: git
        semantic:
          verb: push
          force: true
    - name: git force-with-lease push
      command:
        name: git
        semantic:
          verb: push
          force_with_lease: true
    - name: git force-if-includes push
      command:
        name: git
        semantic:
          verb: push
          force_if_includes: true
```

## aws

Use AWS semantic fields for service, operation, profile, region, endpoint, dry
run, pager behavior, and parser-recognized flags.

Boolean fields:

- `dry_run`: true when `--dry-run` is present, false when `--no-dry-run` is
  present, and unset when neither form is recognized.
- `no_cli_pager`: true when `--no-cli-pager` is present.

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

Boolean fields:

- `namespace_missing`: true when no namespace was selected.
- `all_namespaces`: true for `-A` or `--all-namespaces`.
- `selector_missing`: true when no selector was selected.
- `dry_run`: true for `--dry-run` or a dry-run value other than `none`; false
  for `--dry-run=none`; unset when absent.
- `force`: true for `--force`.
- `recursive`: true for `-R` or `--recursive`.

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

Use GitHub CLI semantic fields for `api`, `pr`, `issue`, `repo`, `release`,
`secret`, `search`, `workflow`, `auth`, and `run` workflows, including
repository selectors, API method and endpoint, issue metadata, release tags,
secret names, workflow refs, search query text, PR merge options, run rerun
options, and parser-recognized flags.

Boolean fields include:

- `web`: true for `-w` or `--web`.
- `paginate`, `input`, `silent`, `include_headers`: true for the corresponding
  `gh api` options.
- `draft`, `fill`: true for the corresponding `gh pr create` options.
- `draft`, `prerelease`, `latest`: true for the corresponding `gh release`
  options.
- `force`: true for `gh pr checkout --force` or `-f`, and `gh run rerun
  --force`.
- `admin`, `auto`, `delete_branch`: true for the corresponding `gh pr merge`
  options.
- `failed`, `debug`, `exit_status`: true for the corresponding `gh run`
  options.

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

Issue examples:

```yaml
permission:
  allow:
    - name: gh issue read-only
      command:
        name: gh
        semantic:
          area: issue
          verb_in: [view, list, status]
  ask:
    - name: gh issue comment critical repo
      command:
        name: gh
        semantic:
          area: issue
          verb: comment
          repo_in: [owner/prod]
```

Additional area examples:

```yaml
permission:
  allow:
    - name: gh safe read operations
      command:
        name: gh
        semantic:
          area_in: [repo, release, secret, workflow, auth]
          verb_in: [view, list, status]
  deny:
    - name: gh token exposure
      command:
        name: gh
        semantic:
          area: auth
          verb: token
    - name: gh remove secrets
      command:
        name: gh
        semantic:
          area: secret
          verb: remove
  ask:
    - name: production release create
      command:
        name: gh
        semantic:
          area: release
          verb: create
          prerelease: false
    - name: production workflow run
      command:
        name: gh
        semantic:
          area: workflow
          verb: run
          repo_in: [owner/prod]
```

## gws

Use Google Workspace CLI semantic fields for dynamic Discovery commands and
hand-written helper commands. The `gws` command surface is built largely from
Google Discovery Service, so the parser intentionally exposes generic
`service`, `resource_path`, and `method` fields rather than a closed enum of
every API method.

Core fields:

- `service`, `service_in`: first action token after `gws`, such as `drive`,
  `gmail`, `calendar`, `sheets`, `docs`, `chat`, or `auth`.
- `resource_path`, `resource_path_contains`: resource tokens between service
  and method, such as `[files]` or `[spreadsheets, values]`.
- `method`, `method_in`: final Discovery method or helper command, such as
  `list`, `get`, `create`, `delete`, `export`, `login`, `+send`, or `+upload`.
- `helper`: true when the method starts with `+`.
- `mutating`, `destructive`, `read_only`: conservative method-name
  classifications.
- `dry_run`, `page_all`, `upload`, `sanitize`, `params`, `json_body`,
  `unmasked`: true when the corresponding parser-recognized option is present.
- `scopes`: scopes selected by `--scopes` or `-s`, split on commas and spaces.

```yaml
permission:
  allow:
    - name: gws drive list files
      command:
        name: gws
        semantic:
          service: drive
          resource_path:
            - files
          method: list
      test:
        allow:
          - "gws drive files list --params '{\"pageSize\": 5}'"
        abstain:
          - "gws drive files delete --params '{\"fileId\":\"abc\"}'"

  deny:
    - name: gws unmasked credential export
      command:
        name: gws
        semantic:
          service: auth
          method: export
          unmasked: true
      test:
        deny:
          - "gws auth export --unmasked"
        abstain:
          - "gws auth login"
```

## helmfile

Use helmfile semantic fields for verb, environment, file, namespace,
kube-context, selector, dry run, wait behavior, delete behavior, state values,
and parser-recognized flags.

Boolean fields include:

- `environment_missing`, `file_missing`, `namespace_missing`,
  `kube_context_missing`, `selector_missing`: true when the corresponding
  value was not selected.
- `interactive`: true for `--interactive`.
- `dry_run`: true for `--dry-run`.
- `wait`: true for `--wait`.
- `wait_for_jobs`: true for `--wait-for-jobs`.
- `skip_diff`: true for `--skip-diff`.
- `skip_needs`: true for `--skip-needs`.
- `include_needs`: true for `--include-needs`.
- `include_transitive_needs`: true for `--include-transitive-needs`.
- `purge`: true for `--purge`.
- `delete_wait`: true for `--delete-wait`.

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

## argocd

Use Argo CD semantic fields for app action path, app name, project, revision,
and parser-recognized flags. App verbs are represented as action paths such as
`app get`, `app list`, `app diff`, `app sync`, `app rollback`, and `app
delete`.

```yaml
permission:
  allow:
    - name: argocd app read-only
      command:
        name: argocd
        semantic:
          verb_in: [app get, app list, app diff]
  ask:
    - name: argocd app sync
      command:
        name: argocd
        semantic:
          verb: app sync
  deny:
    - name: argocd destructive app ops
      command:
        name: argocd
        semantic:
          verb_in: [app delete, app rollback]
```

## Unsupported Commands

When a command is not listed by `cc-bash-guard help semantic`, write a
`patterns` rule instead of `command.semantic`:

```yaml
permission:
  ask:
    - name: tool preview
      patterns:
        - "^my-tool\\s+preview(\\s|$)"
```
