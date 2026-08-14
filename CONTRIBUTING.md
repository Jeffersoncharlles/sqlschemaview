# Contributing

Contributions to SQL Schema View are welcome. The alpha prioritizes correctness, deterministic output, and a deliberately small SQL subset.

## Development Setup

1. Install Go 1.24 or newer.
2. Fork and clone the repository.
3. Create a focused branch from the default branch.
4. Run `make test` before submitting changes.

## Engineering Rules

- Keep the Schema IR independent from SQL dialects and output formats.
- Do not parse complete SQL structures primarily with regular expressions.
- Create relationships only from explicit foreign keys.
- Preserve composite constraints as single domain objects.
- Keep output deterministic and use LF line endings.
- Add unit tests for parser and domain behavior.
- Add or update golden tests for renderer changes.
- Do not expand the alpha scope without prior discussion.

## Pull Requests

- Keep changes small and focused.
- Explain behavioral changes and tradeoffs.
- Include fixtures for supported and rejected SQL.
- Ensure `go test ./...` and `go vet ./...` pass.
- Update `CHANGELOG.md` for user-visible changes.

## Reporting Bugs

Include the smallest SQL input that reproduces the problem, the command used, the expected result, the actual result, the operating system, and the SQL Schema View version.
