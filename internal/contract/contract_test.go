package contract

import (
	"strings"
	"testing"

	"github.com/tasuku43/cmdproxy/internal/domain/policy"
)

func TestValidateRewritesAcceptsStrictAWSMappings(t *testing.T) {
	issues := ValidateRewrites([]policy.RewriteStepSpec{{
		Match: policy.MatchSpec{Command: "aws"},
		MoveFlagToEnv: policy.MoveFlagToEnvSpec{
			Flag: "--profile",
			Env:  "AWS_PROFILE",
		},
		Test: policy.RewriteTestSpec{
			Expect: []policy.RewriteExpectCase{{In: "aws --profile dev sts get-caller-identity", Out: "AWS_PROFILE=dev aws sts get-caller-identity"}},
			Pass:   []string{"AWS_PROFILE=dev aws sts get-caller-identity"},
		},
	}})
	if len(issues) != 0 {
		t.Fatalf("issues = %v", issues)
	}
}

func TestValidateRewritesRejectsUnknownEnvMapping(t *testing.T) {
	issues := ValidateRewrites([]policy.RewriteStepSpec{{
		Match: policy.MatchSpec{Command: "aws"},
		MoveFlagToEnv: policy.MoveFlagToEnvSpec{
			Flag: "--profile",
			Env:  "HOGE",
		},
		Test: policy.RewriteTestSpec{
			Expect: []policy.RewriteExpectCase{{In: "aws --profile dev sts get-caller-identity", Out: "HOGE=dev aws sts get-caller-identity"}},
			Pass:   []string{"aws sts get-caller-identity"},
		},
	}})
	if len(issues) == 0 || !strings.Contains(issues[0], "AWS_PROFILE") {
		t.Fatalf("issues = %v", issues)
	}
}

func TestValidateRewritesAllowsRelaxedKubectlMapping(t *testing.T) {
	strict := false
	issues := ValidateRewrites([]policy.RewriteStepSpec{{
		Match: policy.MatchSpec{Command: "kubectl"},
		MoveFlagToEnv: policy.MoveFlagToEnvSpec{
			Flag: "--kubeconfig",
			Env:  "KUBECONFIG",
		},
		Strict: &strict,
		Test: policy.RewriteTestSpec{
			Expect: []policy.RewriteExpectCase{{In: "kubectl --kubeconfig /tmp/dev get pods", Out: "KUBECONFIG=/tmp/dev kubectl get pods"}},
			Pass:   []string{"KUBECONFIG=/tmp/dev kubectl get pods"},
		},
	}})
	if len(issues) != 0 {
		t.Fatalf("issues = %v", issues)
	}
}

func TestValidateRewritesSkipsStripCommandPathContract(t *testing.T) {
	issues := ValidateRewrites([]policy.RewriteStepSpec{{
		Match:            policy.MatchSpec{CommandIsAbsolutePath: true},
		StripCommandPath: true,
		Test: policy.RewriteTestSpec{
			Expect: []policy.RewriteExpectCase{{In: "/bin/ls", Out: "ls"}},
			Pass:   []string{"ls"},
		},
	}})
	if len(issues) != 0 {
		t.Fatalf("issues = %v", issues)
	}
}
