// Package cli owns argument parsing, usage text, and process exit semantics.
package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	// Version is the current application version.
	Version = "0.1.0-alpha.1"

	ExitSuccess = 0
	ExitError   = 1
	ExitUsage   = 2
)

// ErrUsage identifies invalid command-line arguments.
var ErrUsage = errors.New("invalid command usage")

// Options contains the arguments required to run one documentation generation.
type Options struct {
	InputPath string
	OutputDir string
}

const help = `SQL Schema View

Generate Markdown database documentation from SQL.

Usage:
  sqlschemaview <input-path> [flags]

Flags:
  -o, --output <path>   Output directory (default: ./database-docs)
  -h, --help            Show help
  -v, --version         Show version
`

// Parse validates one processing invocation independently from application execution.
func Parse(args []string) (Options, error) {
	options := Options{OutputDir: "./database-docs"}
	inputPaths := make([]string, 0, 1)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-o", "--output":
			if i+1 >= len(args) {
				return Options{}, fmt.Errorf("%w: %s requires a path", ErrUsage, arg)
			}

			i++
			if args[i] == "" {
				return Options{}, fmt.Errorf("%w: %s requires a non-empty path", ErrUsage, arg)
			}
			options.OutputDir = args[i]
		default:
			if strings.HasPrefix(arg, "-") {
				return Options{}, fmt.Errorf("%w: unknown option %q", ErrUsage, arg)
			}
			inputPaths = append(inputPaths, arg)
		}
	}

	if len(inputPaths) != 1 {
		return Options{}, fmt.Errorf("%w: exactly one input path is required", ErrUsage)
	}

	options.InputPath = inputPaths[0]
	return options, nil
}

// Run is the process boundary used by main and CLI integration tests.
func Run(args []string, stdout, stderr io.Writer) int {
	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			fmt.Fprint(stdout, help)
			return ExitSuccess
		case "-v", "--version":
			fmt.Fprintln(stdout, Version)
			return ExitSuccess
		}
	}

	if _, err := Parse(args); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		fmt.Fprint(stderr, help)
		return ExitUsage
	}

	fmt.Fprintln(stderr, "error: application workflow is not implemented yet")
	return ExitError
}
