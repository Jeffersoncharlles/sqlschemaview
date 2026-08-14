# SQL Schema View

SQL Schema View is a standalone Go CLI that extracts relational structures from SQL files and generates navigable Markdown documentation with Mermaid ER diagrams.

> Status: architecture scaffold for `0.1.0-alpha.1`. Schema extraction and documentation generation are not implemented yet.

## Goals

- Read one SQL file or one directory recursively.
- Aggregate all supported SQL files into one logical schema.
- Extract tables, columns, defaults, nullability, primary keys, and foreign keys.
- Resolve explicit relationships without database access or name-based inference.
- Generate one `README.md` and one Markdown document per table.
- Produce deterministic, cross-platform output from a standalone Go binary.

## Planned Usage

```bash
sqlschemaview database.sql
sqlschemaview ./migrations
sqlschemaview database.sql --output ./docs/database
```

## Development

Requirements:

- Go 1.24 or newer.

Common commands:

```bash
make fmt
make test
make lint
make build
```

## Architecture

```text
Input Loader
    -> SQL Extractor
    -> Schema IR
    -> Relationship Resolver
    -> Documentation Renderer
    -> Atomic Output Writer
```

The Schema IR is the central domain. Parsing, relationship resolution, Mermaid rendering, Markdown rendering, and filesystem persistence remain separate concerns.

## Contributing

See `CONTRIBUTING.md` for the development workflow and contribution requirements.

## License

Distributed under the MIT License. See `LICENSE`.
