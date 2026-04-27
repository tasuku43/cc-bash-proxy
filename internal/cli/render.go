package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tasuku43/cc-bash-guard/internal/app"
)

type colorScheme struct {
	enabled bool
}

func (c colorScheme) wrap(code string, s string) string {
	if !c.enabled {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

func (c colorScheme) green(s string) string  { return c.wrap("32", s) }
func (c colorScheme) red(s string) string    { return c.wrap("31", s) }
func (c colorScheme) yellow(s string) string { return c.wrap("33", s) }
func (c colorScheme) cyan(s string) string   { return c.wrap("36", s) }
func (c colorScheme) bold(s string) string   { return c.wrap("1", s) }
func (c colorScheme) dim(s string) string    { return c.wrap("2", s) }

func colorFor(w io.Writer, mode string) colorScheme {
	switch mode {
	case "always":
		return colorScheme{enabled: os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"}
	case "never":
		return colorScheme{}
	default:
		return colorScheme{enabled: os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb" && isTerminal(w)}
	}
}

func isTerminal(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func writeVerifyText(w io.Writer, result app.VerifyResult, color colorScheme) {
	if result.Verified {
		fmt.Fprintf(w, "%s verify\n", color.green(color.bold("PASS")))
		artifact := "verified"
		if result.ArtifactBuilt {
			artifact = "updated"
		}
		writeVerifySummary(w, result, artifact)
		return
	}
	fmt.Fprintf(w, "%s verify\n", color.red(color.bold("FAIL")))
	fmt.Fprintf(w, "  failures: %s\n", color.red(fmtInt(result.Summary.Failures)))
	fmt.Fprintf(w, "  warnings: %s\n", warningCountText(color, result.Summary.Warnings))
	fmt.Fprintln(w)
	for i, failure := range result.Diagnostics {
		writeVerifyDiagnostic(w, color, "Failure", i+1, failure)
	}
	if len(result.Warnings) > 0 {
		fmt.Fprintf(w, "%s warnings: %d\n\n", color.yellow(color.bold("WARN")), len(result.Warnings))
		for i, warning := range result.Warnings {
			writeVerifyDiagnostic(w, color, "Warning", i+1, warning)
		}
	}
	fmt.Fprintln(w, color.bold("Next:"))
	fmt.Fprintln(w, "  Fix the failures above and run:")
	fmt.Fprintln(w, "    cc-bash-guard verify")
}

func writeVerifySummary(w io.Writer, result app.VerifyResult, artifactStatus string) {
	fmt.Fprintf(w, "  config files: %d\n", result.Summary.ConfigFiles)
	fmt.Fprintf(w, "  permission rules: %d\n", result.Summary.PermissionRules)
	fmt.Fprintf(w, "  tests: %d\n", result.Summary.Tests)
	fmt.Fprintf(w, "  artifact: %s\n", artifactStatus)
}

func warningCountText(color colorScheme, count int) string {
	if count == 0 {
		return fmtInt(count)
	}
	return color.yellow(fmtInt(count))
}

func writeVerifyDiagnostic(w io.Writer, color colorScheme, label string, index int, d app.VerifyDiagnostic) {
	title := d.Title
	if title == "" {
		title = strings.ReplaceAll(d.Kind, "_", " ")
	}
	fmt.Fprintf(w, "%s %d: %s\n", label, index, color.bold(title))
	if d.Source != nil {
		fmt.Fprintf(w, "  source: %s\n", color.cyan(formatVerifySource(*d.Source)))
	}
	if d.Input != "" {
		fmt.Fprintf(w, "  input: %s\n", d.Input)
	}
	if d.Expected != "" {
		fmt.Fprintf(w, "  expected: %s\n", decisionText(color, d.Expected))
	}
	if d.Actual != "" {
		fmt.Fprintf(w, "  actual: %s\n", decisionText(color, d.Actual))
	}
	if d.Reason != "" {
		fmt.Fprintf(w, "  reason: %s\n", d.Reason)
	}
	if d.Command != "" {
		fmt.Fprintf(w, "  command: %s\n", d.Command)
	}
	if d.Field != "" {
		fmt.Fprintf(w, "  field: %s\n", d.Field)
	}
	if d.ExpectedType != "" || d.ActualType != "" {
		fmt.Fprintf(w, "  expected: %s\n", d.ExpectedType)
		fmt.Fprintf(w, "  actual: %s\n", d.ActualType)
	}
	if d.Message != "" && d.Input == "" && d.Field == "" {
		fmt.Fprintf(w, "  message: %s\n", d.Message)
	}
	if d.Decisions != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  decisions:")
		fmt.Fprintf(w, "    cc-bash-guard: %s\n", decisionText(color, d.Decisions.Policy))
		fmt.Fprintf(w, "    Claude settings: %s\n", decisionText(color, d.Decisions.ClaudeSettings))
		fmt.Fprintf(w, "    final: %s\n", decisionText(color, d.Decisions.Final))
	}
	if d.MatchedRule != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  matched rule:")
		fmt.Fprintf(w, "    source: %s\n", color.cyan(formatVerifySource(*d.MatchedRule)))
		if d.MatchedRule.Name != "" {
			fmt.Fprintf(w, "    name: %s\n", d.MatchedRule.Name)
		}
		if d.MatchedMessage != "" {
			fmt.Fprintf(w, "    message: %s\n", d.MatchedMessage)
		}
	}
	if len(d.SupportedFields) > 0 {
		fmt.Fprintln(w)
		if d.Command != "" {
			fmt.Fprintf(w, "  Supported fields for %s:\n", d.Command)
		} else {
			fmt.Fprintln(w, "  supported fields:")
		}
		fmt.Fprintf(w, "    %s\n", strings.Join(d.SupportedFields, ", "))
	}
	if d.First != nil || d.Second != nil {
		if d.First != nil {
			fmt.Fprintf(w, "  first: %s\n", color.cyan(formatVerifySource(*d.First)))
		}
		if d.Second != nil {
			fmt.Fprintf(w, "  second: %s\n", color.cyan(formatVerifySource(*d.Second)))
		}
	}
	if d.Hint != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  Hint:")
		fmt.Fprintf(w, "    %s\n", d.Hint)
	}
	fmt.Fprintln(w)
}

func decisionText(color colorScheme, decision string) string {
	switch decision {
	case "allow":
		return color.green(decision)
	case "deny":
		return color.red(decision)
	case "ask":
		return color.yellow(decision)
	default:
		return color.dim(decision)
	}
}

func formatVerifySource(src app.VerifySource) string {
	var b strings.Builder
	if src.File != "" {
		b.WriteString(src.File)
		if src.Section != "" || src.Bucket != "" {
			b.WriteByte(' ')
		}
	}
	switch {
	case src.Section == "permission" && src.Bucket != "":
		b.WriteString("permission.")
		b.WriteString(src.Bucket)
		b.WriteByte('[')
		b.WriteString(fmtInt(src.Index))
		b.WriteByte(']')
	case src.Section == "test":
		b.WriteString("test[")
		b.WriteString(fmtInt(src.Index))
		b.WriteByte(']')
	case src.Section != "":
		b.WriteString(src.Section)
	}
	if src.Name != "" {
		b.WriteString(" \"")
		b.WriteString(src.Name)
		b.WriteByte('"')
	}
	return b.String()
}

func fmtInt(v int) string {
	return fmt.Sprintf("%d", v)
}

func writeVersionText(w io.Writer, result app.VersionResult) {
	fmt.Fprintf(w, "cc-bash-guard %s\n", result.Info.Version)
	fmt.Fprintf(w, "module: %s\n", result.Info.Module)
	if result.Info.GoVersion != "" {
		fmt.Fprintf(w, "go: %s\n", result.Info.GoVersion)
	}
	if result.Info.VCSRevision != "" {
		fmt.Fprintf(w, "vcs.revision: %s\n", result.Info.VCSRevision)
	}
	if result.Info.VCSTime != "" {
		fmt.Fprintf(w, "vcs.time: %s\n", result.Info.VCSTime)
	}
	if result.Info.VCSModified != "" {
		fmt.Fprintf(w, "vcs.modified: %s\n", result.Info.VCSModified)
	}
}

func writeLine(w io.Writer, line string) {
	fmt.Fprintln(w, line)
}
