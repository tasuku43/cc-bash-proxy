package semantic

type Schema struct {
	Command     string    `json:"command"`
	Description string    `json:"description"`
	Parser      string    `json:"parser"`
	Fields      []Field   `json:"fields"`
	Examples    []Example `json:"examples"`
	Notes       []string  `json:"notes,omitempty"`
}

type Field struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Values      []string `json:"values,omitempty"`
	Since       string   `json:"since,omitempty"`
}

type Example struct {
	Title string `json:"title"`
	YAML  string `json:"yaml"`
}

var schemas = []Schema{
	{
		Command:     "git",
		Description: "Git operations such as push, clean, reset, diff, checkout, switch, and status.",
		Parser:      "git",
		Fields: []Field{
			stringField("verb", "Git verb parsed after global git options."),
			stringListField("verb_in", "Allowed Git verbs."),
			stringField("remote", "Remote positional for commands such as git push."),
			stringListField("remote_in", "Allowed remotes."),
			stringField("branch", "Branch positional for push, checkout, or switch."),
			stringListField("branch_in", "Allowed branches."),
			stringField("ref", "Ref positional for push, reset, checkout, or switch."),
			stringListField("ref_in", "Allowed refs."),
			boolField("force", "For git push, true when --force, -f, --force-with-lease, or --force-if-includes is present. For git clean, true when -f or --force is present."),
			boolField("hard", "True for git reset --hard."),
			boolField("recursive", "True for git clean -d."),
			boolField("include_ignored", "True for git clean -x or --ignored."),
			boolField("cached", "True for git diff --cached or --staged."),
			boolField("staged", "True for git diff --cached or --staged."),
			stringListField("flags_contains", "Parser-recognized git option tokens that must be present."),
			stringListField("flags_prefixes", "Parser-recognized git option tokens that must start with these prefixes."),
		},
		Examples: []Example{
			{Title: "Deny force pushes", YAML: `permission:
  deny:
    - command:
        name: git
        semantic:
          verb: push
          force: true`},
		},
		Notes: []string{
			"`force` intentionally treats --force-with-lease and --force-if-includes as force-like for git push.",
			"`flags_contains` and `flags_prefixes` inspect parser-recognized option tokens, not raw argv words.",
		},
	},
	{
		Command:     "aws",
		Description: "AWS CLI service, operation, profile, region, endpoint, and dry-run matching.",
		Parser:      "aws",
		Fields: []Field{
			stringField("service", "AWS service name such as s3 or iam."),
			stringListField("service_in", "Allowed AWS services."),
			stringField("operation", "AWS operation name."),
			stringListField("operation_in", "Allowed AWS operations."),
			stringField("profile", "AWS profile selected by --profile or AWS_PROFILE."),
			stringListField("profile_in", "Allowed AWS profiles."),
			stringField("region", "AWS region selected by --region or environment."),
			stringListField("region_in", "Allowed AWS regions."),
			stringField("endpoint_url", "Exact --endpoint-url value."),
			stringField("endpoint_url_prefix", "--endpoint-url prefix."),
			boolField("dry_run", "True when the parser recognizes an AWS dry-run option or operation shape."),
			boolField("no_cli_pager", "True when --no-cli-pager is present."),
			stringListField("flags_contains", "Parser-recognized AWS option tokens that must be present."),
			stringListField("flags_prefixes", "Parser-recognized AWS option tokens that must start with these prefixes."),
		},
		Examples: []Example{
			{Title: "Ask for IAM writes", YAML: `permission:
  ask:
    - command:
        name: aws
        semantic:
          service: iam`},
		},
	},
	{
		Command:     "kubectl",
		Description: "Kubernetes verb, resource, namespace, context, filename, selector, and container matching.",
		Parser:      "kubectl",
		Fields: []Field{
			stringField("verb", "kubectl verb such as get, apply, delete, or exec."),
			stringListField("verb_in", "Allowed kubectl verbs."),
			stringField("subverb", "Secondary action for compound kubectl commands."),
			stringListField("subverb_in", "Allowed kubectl subverbs."),
			stringField("resource_type", "Kubernetes resource type."),
			stringListField("resource_type_in", "Allowed resource types."),
			stringField("resource_name", "Kubernetes resource name."),
			stringListField("resource_name_in", "Allowed resource names."),
			stringField("namespace", "Namespace selected by -n or --namespace."),
			stringListField("namespace_in", "Allowed namespaces."),
			boolField("namespace_missing", "True when no namespace was selected."),
			stringField("context", "Context selected by --context."),
			stringListField("context_in", "Allowed contexts."),
			stringField("kubeconfig", "Kubeconfig path selected by --kubeconfig."),
			boolField("all_namespaces", "True when -A or --all-namespaces is present."),
			stringField("filename", "Filename selected by -f or --filename."),
			stringListField("filename_in", "Allowed filenames."),
			stringField("filename_prefix", "Filename prefix selected by -f or --filename."),
			stringField("selector", "Selector selected by -l or --selector."),
			stringListField("selector_in", "Allowed selectors."),
			stringListField("selector_contains", "Selectors that must be present."),
			boolField("selector_missing", "True when no selector was selected."),
			stringField("container", "Container selected by -c or --container."),
			boolField("dry_run", "True when --dry-run or a dry-run value is present."),
			boolField("force", "True when --force is present."),
			boolField("recursive", "True when -R or --recursive is present."),
			stringListField("flags_contains", "Parser-recognized kubectl option tokens that must be present."),
			stringListField("flags_prefixes", "Parser-recognized kubectl option tokens that must start with these prefixes."),
		},
		Examples: []Example{
			{Title: "Deny production deletes", YAML: `permission:
  deny:
    - command:
        name: kubectl
        semantic:
          verb: delete
          namespace: prod`},
		},
	},
	{
		Command:     "gh",
		Description: "GitHub CLI api, pr, and run operations.",
		Parser:      "gh",
		Fields: []Field{
			stringField("area", "Top-level gh area such as api, pr, or run."),
			stringListField("area_in", "Allowed gh areas."),
			stringField("verb", "gh subcommand verb inside the selected area."),
			stringListField("verb_in", "Allowed gh verbs."),
			stringField("repo", "Repository selected by -R or --repo."),
			stringListField("repo_in", "Allowed repositories."),
			stringField("hostname", "Hostname selected by --hostname."),
			stringListField("hostname_in", "Allowed hostnames."),
			boolField("web", "True when -w or --web is present."),
			stringField("method", "HTTP method for gh api."),
			stringListField("method_in", "Allowed gh api HTTP methods."),
			stringField("endpoint", "Normalized gh api endpoint."),
			stringField("endpoint_prefix", "Normalized gh api endpoint prefix."),
			stringListField("endpoint_contains", "Endpoint substrings that must be present."),
			boolField("paginate", "True when gh api --paginate is present."),
			boolField("input", "True when gh api --input is present."),
			boolField("silent", "True when gh api --silent is present."),
			boolField("include_headers", "True when gh api -i or --include is present."),
			stringListField("field_keys_contains", "gh api -F/--field keys that must be present."),
			stringListField("raw_field_keys_contains", "gh api -f/--raw-field keys that must be present."),
			stringListField("header_keys_contains", "gh api -H/--header keys that must be present."),
			stringField("pr_number", "Pull request number positional for gh pr commands."),
			stringField("base", "Base branch selected by --base."),
			stringField("head", "Head branch selected by --head."),
			boolField("draft", "True when gh pr create --draft is present."),
			boolField("fill", "True when gh pr create --fill is present."),
			boolField("force", "True for gh pr checkout --force or -f, and gh run rerun --force."),
			boolField("admin", "True when gh pr merge --admin is present."),
			boolField("auto", "True when gh pr merge --auto is present."),
			boolField("delete_branch", "True when gh pr merge --delete-branch is present."),
			stringField("merge_strategy", "Merge strategy selected by gh pr merge."),
			stringListField("merge_strategy_in", "Allowed merge strategies."),
			stringField("run_id", "Run ID positional for gh run commands."),
			boolField("failed", "True when gh run rerun --failed is present."),
			stringField("job", "Job selected by gh run rerun --job."),
			boolField("debug", "True when gh run rerun --debug is present."),
			boolField("exit_status", "True when gh run view --exit-status is present."),
			stringListField("flags_contains", "Parser-recognized gh option tokens that must be present."),
			stringListField("flags_prefixes", "Parser-recognized gh option tokens that must start with these prefixes."),
		},
		Examples: []Example{
			{Title: "Deny mutating GitHub API calls", YAML: `permission:
  deny:
    - command:
        name: gh
        semantic:
          area: api
          method_in: [POST, PATCH, PUT, DELETE]`},
		},
	},
	{
		Command:     "helmfile",
		Description: "Helmfile apply, sync, destroy, diff, environment, file, selector, namespace, and values matching.",
		Parser:      "helmfile",
		Fields: []Field{
			stringField("verb", "helmfile verb such as apply, sync, destroy, or diff."),
			stringListField("verb_in", "Allowed helmfile verbs."),
			stringField("environment", "Environment selected by -e or --environment."),
			stringListField("environment_in", "Allowed environments."),
			boolField("environment_missing", "True when no environment was selected."),
			stringField("file", "State file selected by -f or --file."),
			stringListField("file_in", "Allowed state files."),
			stringField("file_prefix", "State file prefix."),
			boolField("file_missing", "True when no state file was selected."),
			stringField("namespace", "Namespace selected by --namespace."),
			stringListField("namespace_in", "Allowed namespaces."),
			boolField("namespace_missing", "True when no namespace was selected."),
			stringField("kube_context", "Kube context selected by --kube-context."),
			stringListField("kube_context_in", "Allowed kube contexts."),
			boolField("kube_context_missing", "True when no kube context was selected."),
			stringField("selector", "Selector selected by -l or --selector."),
			stringListField("selector_in", "Allowed selectors."),
			stringListField("selector_contains", "Selectors that must be present."),
			boolField("selector_missing", "True when no selector was selected."),
			boolField("interactive", "True when --interactive is present."),
			boolField("dry_run", "True when --dry-run is present."),
			boolField("wait", "True when --wait is present."),
			boolField("wait_for_jobs", "True when --wait-for-jobs is present."),
			boolField("skip_diff", "True when --skip-diff is present."),
			boolField("skip_needs", "True when --skip-needs is present."),
			boolField("include_needs", "True when --include-needs is present."),
			boolField("include_transitive_needs", "True when --include-transitive-needs is present."),
			boolField("purge", "True when --purge is present."),
			stringField("cascade", "Cascade value selected by --cascade."),
			stringListField("cascade_in", "Allowed cascade values."),
			boolField("delete_wait", "True when --delete-wait is present."),
			stringField("state_values_file", "State values file selected by --state-values-file."),
			stringListField("state_values_file_in", "Allowed state values files."),
			stringListField("state_values_set_keys_contains", "Keys selected by --state-values-set that must be present."),
			stringListField("state_values_set_string_keys_contains", "Keys selected by --state-values-set-string that must be present."),
			stringListField("flags_contains", "Parser-recognized helmfile option tokens that must be present."),
			stringListField("flags_prefixes", "Parser-recognized helmfile option tokens that must start with these prefixes."),
		},
		Examples: []Example{
			{Title: "Ask before production destroy", YAML: `permission:
  ask:
    - command:
        name: helmfile
        semantic:
          verb: destroy
          environment: prod`},
		},
	},
}

func AllSchemas() []Schema {
	out := append([]Schema(nil), schemas...)
	return out
}

func Lookup(command string) (Schema, bool) {
	for _, schema := range schemas {
		if schema.Command == command {
			return schema, true
		}
	}
	return Schema{}, false
}

func SupportedCommands() []string {
	commands := make([]string, 0, len(schemas))
	for _, schema := range schemas {
		commands = append(commands, schema.Command)
	}
	return commands
}

func FieldNames(command string) []string {
	schema, ok := Lookup(command)
	if !ok {
		return nil
	}
	names := make([]string, 0, len(schema.Fields))
	for _, field := range schema.Fields {
		names = append(names, field.Name)
	}
	return names
}

func IsFieldSupported(command, field string) bool {
	for _, supported := range FieldNames(command) {
		if supported == field {
			return true
		}
	}
	return false
}

func stringField(name, description string) Field {
	return Field{Name: name, Type: "string", Description: description}
}

func stringListField(name, description string) Field {
	return Field{Name: name, Type: "[]string", Description: description}
}

func boolField(name, description string) Field {
	return Field{Name: name, Type: "bool", Description: description}
}
