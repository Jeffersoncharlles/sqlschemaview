// Package schema contains the database-independent Schema IR.
package schema

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

// Column preserves the structural attributes required by documentation.
type Column struct {
	Name         string
	DataType     string
	Nullable     bool
	DefaultValue *string
	PrimaryKey   bool
	ForeignKey   bool
}

// PrimaryKey preserves inline and table-level primary-key definitions.
type PrimaryKey struct {
	Name    *string
	Columns []string
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
