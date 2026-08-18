// Package schema contains the database-independent Schema IR.
package schema

import (
	"fmt"
	"strings"
)

// CanonicalName is the case-insensitive key used to compare simple SQL identifiers.
// The original SQL spelling remains in each Name field.
type CanonicalName string

// CanonicalIdentifier returns the comparison key for a simple SQL identifier.
// Quoted and qualified identifier semantics are handled by later parser work.
func CanonicalIdentifier(name string) CanonicalName {
	return CanonicalName(strings.ToLower(name))
}

// Schema is the complete logical schema aggregated from all SQL inputs.
type Schema struct {
	Tables        []Table
	Relationships []Relationship
}

// Table is one relational table as declared in SQL.
type Table struct {
	Name        string
	Source      Source
	Columns     []Column
	PrimaryKeys []PrimaryKey
	ForeignKeys []ForeignKey
}

// CanonicalName returns the table's case-insensitive comparison key.
func (t Table) CanonicalName() CanonicalName {
	return CanonicalIdentifier(t.Name)
}

// TableByName finds a table by its canonical identifier.
func (s *Schema) TableByName(name string) (*Table, bool) {
	if s == nil {
		return nil, false
	}

	canonicalName := CanonicalIdentifier(name)
	for i := range s.Tables {
		if s.Tables[i].CanonicalName() == canonicalName {
			return &s.Tables[i], true
		}
	}

	return nil, false
}

// Column preserves the structural attributes required by documentation.
type Column struct {
	Name         string
	Source       Source
	DataType     string
	Nullable     bool
	DefaultValue *string
	PrimaryKey   bool
	ForeignKey   bool
}

// CanonicalName returns the column's case-insensitive comparison key.
func (c Column) CanonicalName() CanonicalName {
	return CanonicalIdentifier(c.Name)
}

// ColumnByName finds a column by its canonical identifier.
func (t *Table) ColumnByName(name string) (*Column, bool) {
	if t == nil {
		return nil, false
	}

	canonicalName := CanonicalIdentifier(name)
	for i := range t.Columns {
		if t.Columns[i].CanonicalName() == canonicalName {
			return &t.Columns[i], true
		}
	}

	return nil, false
}

// PrimaryKey preserves inline and table-level primary-key definitions.
type PrimaryKey struct {
	Name    *string
	Columns []string
	Source  Source
}

// ForeignKey preserves a foreign key as one constraint, including composites.
type ForeignKey struct {
	Name              *string
	Columns           []string
	ReferencedTable   string
	ReferencedColumns []string
	Source            Source
}

// Relationship is created only from an explicitly declared foreign key.
type Relationship struct {
	ParentTable   string
	ParentColumns []string
	ChildTable    string
	ChildColumns  []string
	Nullable      bool
}

// Source identifies where a schema-bearing structure was declared.
type Source struct {
	Path   string
	Line   int
	Column int
}

func (s Source) String() string {
	path := s.Path
	if path == "" {
		path = "<unknown>"
	}

	return fmt.Sprintf("%s:%d:%d", path, s.Line, s.Column)
}
