// Package cli owns argument parsing, usage text, and process exit semantics.
package cli

import (
	"fmt"
	"io"
)

const (
	// Version is the current application version.
	Version = "0.1.0-alpha.1"

	ExitSuccess = 0
	ExitError   = 1
	ExitUsage   = 2
)

const help = `SQL Schema View

Generate Markdown database documentation from SQL.

Usage:
  sqlschemaview <input-path> [flags]

Flags:
  -o, --output <path>   Output directory (default: ./database-docs)
  -h, --help            Show help
  -v, --version         Show version
`

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

	if len(args) == 0 {
		fmt.Fprintln(stderr, "error: exactly one input path is required")
		fmt.Fprint(stderr, help)
		return ExitUsage
	}

	fmt.Fprintln(stderr, "error: application workflow is not implemented yet")
	return ExitError
}
