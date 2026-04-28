package command

import "testing"

func FuzzParseCommandPlanDoNotPanic(f *testing.F) {
	for _, seed := range commandParserRobustnessSeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		_ = Parse(input)
		_ = ParseWithRegistry(input, DefaultParserRegistry())
	})
}

func commandParserRobustnessSeeds() []string {
	return []string{
		"git status",
		"git push --force origin main",
		"bash -c 'git status'",
		"bash -c 'git push --force origin main'",
		"env bash -c 'git status'",
		"sudo -u root bash -c 'git status'",
		"timeout 10 bash -c 'git status'",
		"cat <(rm -rf /tmp/x)",
		"echo $(cat ~/.ssh/id_rsa)",
		"git status > /tmp/out",
		"(git status)",
		"git status &",
		"git status 'unterminated",
		"git status \\",
		"git status\nrm -rf /tmp/x",
		"bash -c \"git status",
		"",
	}
}
