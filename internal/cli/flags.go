package cli

import (
	"errors"
	"fmt"
	"strings"
)

func parseCommonFlags(args []string) (string, []string, error) {
	format := ""
	rest := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--format":
			if i+1 >= len(args) {
				return "", nil, errors.New("missing --format value")
			}
			format = args[i+1]
			i++
		default:
			rest = append(rest, args[i])
		}
	}
	if format != "" && format != "json" {
		return "", nil, fmt.Errorf("unsupported format %q", format)
	}
	return format, rest, nil
}

func parseVerifyFlags(args []string) (format string, color string, allFailures bool, rest []string, err error) {
	color = "auto"
	rest = make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--format":
			if i+1 >= len(args) {
				return "", "", false, nil, errors.New("missing --format value")
			}
			format = args[i+1]
			i++
		case arg == "--color":
			if i+1 >= len(args) {
				return "", "", false, nil, errors.New("missing --color value")
			}
			color = args[i+1]
			i++
		case arg == "--all-failures":
			allFailures = true
		case strings.HasPrefix(arg, "--format="):
			format = strings.TrimPrefix(arg, "--format=")
		case strings.HasPrefix(arg, "--color="):
			color = strings.TrimPrefix(arg, "--color=")
		default:
			rest = append(rest, arg)
		}
	}
	if format != "" && format != "json" {
		return "", "", false, nil, fmt.Errorf("unsupported format %q", format)
	}
	switch color {
	case "auto", "always", "never":
	default:
		return "", "", false, nil, fmt.Errorf("unsupported color mode %q", color)
	}
	return format, color, allFailures, rest, nil
}
