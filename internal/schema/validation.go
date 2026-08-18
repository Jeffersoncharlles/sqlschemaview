package schema

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	ErrNoTables              = errors.New("no tables found in SQL input")
	ErrDuplicateTable        = errors.New("duplicate table definition")
	ErrDuplicateColumn       = errors.New("duplicate column definition")
	ErrForeignKeyCardinality = errors.New("foreign key column count mismatch")
	ErrUnresolvedForeignKey  = errors.New("foreign key references an unknown table or column")
)

// DuplicateTableError reports every origin for one duplicate table name.
type DuplicateTableError struct {
	Name    string
	Sources []Source
}

func (e *DuplicateTableError) Error() string {
	return fmt.Sprintf("duplicate table definition %q: declared at %s", e.Name, formatSources(e.Sources))
}

func (e *DuplicateTableError) Unwrap() error {
	return ErrDuplicateTable
}

// DuplicateColumnError reports every origin for one duplicate column name.
type DuplicateColumnError struct {
	TableName string
	Name      string
	Sources   []Source
}

func (e *DuplicateColumnError) Error() string {
	return fmt.Sprintf(
		"duplicate column definition %q in table %q: declared at %s",
		e.Name,
		e.TableName,
		formatSources(e.Sources),
	)
}

func (e *DuplicateColumnError) Unwrap() error {
	return ErrDuplicateColumn
}

// ForeignKeyCardinalityError reports mismatched local and referenced column lists.
type ForeignKeyCardinalityError struct {
	TableName             string
	ConstraintName        *string
	Columns               []string
	ReferencedTable       string
	Source                Source
	ColumnCount           int
	ReferencedColumnCount int
}

func (e *ForeignKeyCardinalityError) Error() string {
	return fmt.Sprintf(
		"foreign key %s on table %q at %s has %d columns but references %d columns in table %q",
		foreignKeyName(e.ConstraintName, e.Columns),
		e.TableName,
		e.Source,
		e.ColumnCount,
		e.ReferencedColumnCount,
		e.ReferencedTable,
	)
}

func (e *ForeignKeyCardinalityError) Unwrap() error {
	return ErrForeignKeyCardinality
}

// Validate checks invariants owned by the Schema IR without resolving references.
func (s *Schema) Validate() error {
	if s == nil || len(s.Tables) == 0 {
		return ErrNoTables
	}

	if err := validateDuplicateTables(s.Tables); err != nil {
		return err
	}
	if err := validateDuplicateColumns(s.Tables); err != nil {
		return err
	}

	return validateForeignKeyCardinality(s.Tables)
}

type namedSource struct {
	name   string
	source Source
}

func validateDuplicateTables(tables []Table) error {
	groups := make(map[CanonicalName][]namedSource, len(tables))
	for _, table := range tables {
		key := table.CanonicalName()
		groups[key] = append(groups[key], namedSource{name: table.Name, source: table.Source})
	}

	keys := duplicateKeys(groups)
	if len(keys) == 0 {
		return nil
	}

	entries := groups[keys[0]]
	return &DuplicateTableError{
		Name:    firstName(entries),
		Sources: sortedSources(entries),
	}
}

func validateDuplicateColumns(tables []Table) error {
	type duplicate struct {
		tableName string
		columnKey CanonicalName
		entries   []namedSource
	}

	var duplicates []duplicate
	for _, table := range tables {
		groups := make(map[CanonicalName][]namedSource, len(table.Columns))
		for _, column := range table.Columns {
			key := column.CanonicalName()
			groups[key] = append(groups[key], namedSource{name: column.Name, source: column.Source})
		}

		for _, key := range duplicateKeys(groups) {
			duplicates = append(duplicates, duplicate{
				tableName: table.Name,
				columnKey: key,
				entries:   groups[key],
			})
		}
	}

	if len(duplicates) == 0 {
		return nil
	}

	sort.Slice(duplicates, func(i, j int) bool {
		leftTable := CanonicalIdentifier(duplicates[i].tableName)
		rightTable := CanonicalIdentifier(duplicates[j].tableName)
		if leftTable != rightTable {
			return leftTable < rightTable
		}
		return duplicates[i].columnKey < duplicates[j].columnKey
	})

	conflict := duplicates[0]
	return &DuplicateColumnError{
		TableName: conflict.tableName,
		Name:      firstName(conflict.entries),
		Sources:   sortedSources(conflict.entries),
	}
}

func validateForeignKeyCardinality(tables []Table) error {
	var invalid []*ForeignKeyCardinalityError
	for _, table := range tables {
		for _, foreignKey := range table.ForeignKeys {
			if len(foreignKey.Columns) > 0 && len(foreignKey.Columns) == len(foreignKey.ReferencedColumns) {
				continue
			}

			invalid = append(invalid, &ForeignKeyCardinalityError{
				TableName:             table.Name,
				ConstraintName:        foreignKey.Name,
				Columns:               append([]string(nil), foreignKey.Columns...),
				ReferencedTable:       foreignKey.ReferencedTable,
				Source:                foreignKey.Source,
				ColumnCount:           len(foreignKey.Columns),
				ReferencedColumnCount: len(foreignKey.ReferencedColumns),
			})
		}
	}

	if len(invalid) == 0 {
		return nil
	}

	sort.Slice(invalid, func(i, j int) bool {
		left := foreignKeySortKey(invalid[i])
		right := foreignKeySortKey(invalid[j])
		return left < right
	})

	return invalid[0]
}

func duplicateKeys(groups map[CanonicalName][]namedSource) []CanonicalName {
	keys := make([]CanonicalName, 0)
	for key, entries := range groups {
		if len(entries) > 1 {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func firstName(entries []namedSource) string {
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.name
	}
	sort.Strings(names)
	return names[0]
}

func sortedSources(entries []namedSource) []Source {
	sources := make([]Source, len(entries))
	for i, entry := range entries {
		sources[i] = entry.source
	}
	sort.Slice(sources, func(i, j int) bool {
		if sources[i].Path != sources[j].Path {
			return sources[i].Path < sources[j].Path
		}
		if sources[i].Line != sources[j].Line {
			return sources[i].Line < sources[j].Line
		}
		return sources[i].Column < sources[j].Column
	})
	return sources
}

func formatSources(sources []Source) string {
	formatted := make([]string, len(sources))
	for i, source := range sources {
		formatted[i] = source.String()
	}
	return strings.Join(formatted, " and ")
}

func foreignKeyName(name *string, columns []string) string {
	if name != nil {
		return fmt.Sprintf("%q", *name)
	}
	if len(columns) == 0 {
		return "(<none>)"
	}
	return fmt.Sprintf("(%s)", strings.Join(columns, ", "))
}

func foreignKeySortKey(err *ForeignKeyCardinalityError) string {
	constraintName := ""
	if err.ConstraintName != nil {
		constraintName = *err.ConstraintName
	}
	return strings.Join([]string{
		string(CanonicalIdentifier(err.TableName)),
		strings.Join(err.Columns, "\x00"),
		string(CanonicalIdentifier(err.ReferencedTable)),
		constraintName,
		err.Source.Path,
		fmt.Sprintf("%010d:%010d", err.Source.Line, err.Source.Column),
	}, "\x00")
}
