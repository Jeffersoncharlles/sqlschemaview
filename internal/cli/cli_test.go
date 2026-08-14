package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		want        Options
		wantErr     bool
		wantErrText string
	}{
		{
			name: "input with default output",
			args: []string{"database.sql"},
			want: Options{InputPath: "database.sql", OutputDir: "./database-docs"},
		},
		{
			name: "short output before input",
			args: []string{"-o", "docs", "database.sql"},
			want: Options{InputPath: "database.sql", OutputDir: "docs"},
		},
		{
			name: "short output after input",
			args: []string{"database.sql", "-o", "docs"},
			want: Options{InputPath: "database.sql", OutputDir: "docs"},
		},
		{
			name: "long output before input",
			args: []string{"--output", "docs", "database.sql"},
			want: Options{InputPath: "database.sql", OutputDir: "docs"},
		},
		{
			name: "long output after input",
			args: []string{"database.sql", "--output", "docs"},
			want: Options{InputPath: "database.sql", OutputDir: "docs"},
		},
		{
			name:        "missing input",
			wantErr:     true,
			wantErrText: "exactly one input path is required",
		},
		{
			name:        "multiple inputs",
			args:        []string{"one.sql", "two.sql"},
			wantErr:     true,
			wantErrText: "exactly one input path is required",
		},
		{
			name:        "missing short output path",
			args:        []string{"database.sql", "-o"},
			wantErr:     true,
			wantErrText: "-o requires a path",
		},
		{
			name:        "missing long output path",
			args:        []string{"database.sql", "--output"},
			wantErr:     true,
			wantErrText: "--output requires a path",
		},
		{
			name:        "empty output path",
			args:        []string{"database.sql", "--output", ""},
			wantErr:     true,
			wantErrText: "--output requires a non-empty path",
		},
		{
			name:        "unknown option",
			args:        []string{"database.sql", "--unknown"},
			wantErr:     true,
			wantErrText: `unknown option "--unknown"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if tt.wantErr {
				if !errors.Is(err, ErrUsage) {
					t.Fatalf("Parse() error = %v, want ErrUsage", err)
				}
				if !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("Parse() error = %q, want text %q", err, tt.wantErrText)
				}
				return
			}

			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--help"}, &stdout, &stderr)

	if code != ExitSuccess {
		t.Fatalf("Run() exit code = %d, want %d", code, ExitSuccess)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("Run() output = %q, want usage text", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("Run() stderr = %q, want empty", stderr.String())
	}
}

func TestRunWithoutInput(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run(nil, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run() exit code = %d, want %d", code, ExitUsage)
	}
	if !strings.Contains(stderr.String(), "exactly one input path is required") {
		t.Fatalf("Run() stderr = %q, want input error", stderr.String())
	}
}

func TestRunVersion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--version"}, &stdout, &stderr)

	if code != ExitSuccess {
		t.Fatalf("Run() exit code = %d, want %d", code, ExitSuccess)
	}
	if strings.TrimSpace(stdout.String()) != Version {
		t.Fatalf("Run() output = %q, want %q", stdout.String(), Version)
	}
	if stderr.Len() != 0 {
		t.Fatalf("Run() stderr = %q, want empty", stderr.String())
	}
}

func TestRunRejectsUnknownOption(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"database.sql", "--unknown"}, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run() exit code = %d, want %d", code, ExitUsage)
	}
	if !strings.Contains(stderr.String(), `unknown option "--unknown"`) {
		t.Fatalf("Run() stderr = %q, want unknown option error", stderr.String())
	}
}
