package schema

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestCanonicalLookupsPreserveOriginalNames(t *testing.T) {
	schema := &Schema{Tables: []Table{{
		Name: "Users",
		Columns: []Column{{
			Name: "Email",
		}},
	}}}

	if got, want := CanonicalIdentifier("USERS"), CanonicalName("users"); got != want {
		t.Fatalf("CanonicalIdentifier() = %q, want %q", got, want)
	}

	table, ok := schema.TableByName("users")
	if !ok {
		t.Fatal("TableByName() did not find a case-insensitive match")
	}
	if table.Name != "Users" {
		t.Fatalf("TableByName().Name = %q, want original name %q", table.Name, "Users")
	}

	column, ok := table.ColumnByName("EMAIL")
	if !ok {
		t.Fatal("ColumnByName() did not find a case-insensitive match")
	}
	if column.Name != "Email" {
		t.Fatalf("ColumnByName().Name = %q, want original name %q", column.Name, "Email")
	}

	if _, ok := schema.TableByName("posts"); ok {
		t.Fatal("TableByName() found an unknown table")
	}
	if _, ok := (*Schema)(nil).TableByName("users"); ok {
		t.Fatal("TableByName() found a table in a nil schema")
	}
	if _, ok := (*Table)(nil).ColumnByName("email"); ok {
		t.Fatal("ColumnByName() found a column in a nil table")
	}
}

func TestValidateAcceptsTablesWithoutPrimaryOrForeignKeys(t *testing.T) {
	schema := &Schema{Tables: []Table{{
		Name:    "events",
		Source:  Source{Path: "schema.sql", Line: 1, Column: 1},
		Columns: []Column{{Name: "payload", Source: Source{Path: "schema.sql", Line: 2, Column: 5}}},
	}}}

	if err := schema.Validate(); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func TestValidatePreservesCompositeConstraints(t *testing.T) {
	primaryKeyName := "pk_memberships"
	foreignKeyName := "fk_memberships_users"
	schema := &Schema{Tables: []Table{{
		Name: "memberships",
		PrimaryKeys: []PrimaryKey{{
			Name:    &primaryKeyName,
			Columns: []string{"tenant_id", "user_id"},
			Source:  Source{Path: "schema.sql", Line: 5, Column: 5},
		}},
		ForeignKeys: []ForeignKey{{
			Name:              &foreignKeyName,
			Columns:           []string{"tenant_id", "user_id"},
			ReferencedTable:   "users",
			ReferencedColumns: []string{"tenant_id", "id"},
			Source:            Source{Path: "schema.sql", Line: 6, Column: 5},
		}},
	}}}
	want := &Schema{Tables: []Table{{
		Name: "memberships",
		PrimaryKeys: []PrimaryKey{{
			Name:    &primaryKeyName,
			Columns: []string{"tenant_id", "user_id"},
			Source:  Source{Path: "schema.sql", Line: 5, Column: 5},
		}},
		ForeignKeys: []ForeignKey{{
			Name:              &foreignKeyName,
			Columns:           []string{"tenant_id", "user_id"},
			ReferencedTable:   "users",
			ReferencedColumns: []string{"tenant_id", "id"},
			Source:            Source{Path: "schema.sql", Line: 6, Column: 5},
		}},
	}}}

	if err := schema.Validate(); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(schema, want) {
		t.Fatalf("Validate() changed composite constraints:\n got: %#v\nwant: %#v", schema, want)
	}
}

func TestValidateRejectsSchemaWithoutTables(t *testing.T) {
	for _, schema := range []*Schema{nil, {}} {
		if err := schema.Validate(); !errors.Is(err, ErrNoTables) {
			t.Fatalf("Validate() error = %v, want %v", err, ErrNoTables)
		}
	}
}

func TestValidateRejectsDuplicateTablesWithAllOrigins(t *testing.T) {
	schema := &Schema{Tables: []Table{
		{Name: "users", Source: Source{Path: "002.sql", Line: 8, Column: 1}},
		{Name: "Users", Source: Source{Path: "001.sql", Line: 3, Column: 1}},
		{Name: "USERS", Source: Source{Path: "003.sql", Line: 2, Column: 4}},
	}}

	err := schema.Validate()
	if !errors.Is(err, ErrDuplicateTable) {
		t.Fatalf("Validate() error = %v, want %v", err, ErrDuplicateTable)
	}

	var duplicate *DuplicateTableError
	if !errors.As(err, &duplicate) {
		t.Fatalf("Validate() error type = %T, want *DuplicateTableError", err)
	}
	if duplicate.Name != "USERS" {
		t.Fatalf("DuplicateTableError.Name = %q, want deterministic name %q", duplicate.Name, "USERS")
	}
	wantSources := []Source{
		{Path: "001.sql", Line: 3, Column: 1},
		{Path: "002.sql", Line: 8, Column: 1},
		{Path: "003.sql", Line: 2, Column: 4},
	}
	if !reflect.DeepEqual(duplicate.Sources, wantSources) {
		t.Fatalf("DuplicateTableError.Sources = %#v, want %#v", duplicate.Sources, wantSources)
	}
	for _, part := range []string{"USERS", "001.sql:3:1", "002.sql:8:1", "003.sql:2:4"} {
		if !strings.Contains(err.Error(), part) {
			t.Errorf("Validate() error %q does not contain %q", err, part)
		}
	}
}

func TestValidateRejectsDuplicateColumnsWithinTable(t *testing.T) {
	schema := &Schema{Tables: []Table{
		{
			Name: "users",
			Columns: []Column{
				{Name: "email", Source: Source{Path: "schema.sql", Line: 4, Column: 5}},
				{Name: "Email", Source: Source{Path: "schema.sql", Line: 7, Column: 5}},
			},
		},
		{Name: "audit", Columns: []Column{{Name: "email"}}},
	}}

	err := schema.Validate()
	if !errors.Is(err, ErrDuplicateColumn) {
		t.Fatalf("Validate() error = %v, want %v", err, ErrDuplicateColumn)
	}

	var duplicate *DuplicateColumnError
	if !errors.As(err, &duplicate) {
		t.Fatalf("Validate() error type = %T, want *DuplicateColumnError", err)
	}
	if duplicate.TableName != "users" || duplicate.Name != "Email" {
		t.Fatalf("duplicate column = %q.%q, want users.Email", duplicate.TableName, duplicate.Name)
	}
	for _, part := range []string{"users", "Email", "schema.sql:4:5", "schema.sql:7:5"} {
		if !strings.Contains(err.Error(), part) {
			t.Errorf("Validate() error %q does not contain %q", err, part)
		}
	}
}

func TestValidateAllowsSameColumnNameInDifferentTables(t *testing.T) {
	schema := &Schema{Tables: []Table{
		{Name: "users", Columns: []Column{{Name: "id"}}},
		{Name: "posts", Columns: []Column{{Name: "ID"}}},
	}}

	if err := schema.Validate(); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func TestValidateForeignKeyCardinality(t *testing.T) {
	tests := []struct {
		name              string
		columns           []string
		referencedColumns []string
		wantError         bool
	}{
		{name: "composite", columns: []string{"tenant_id", "user_id"}, referencedColumns: []string{"tenant_id", "id"}},
		{name: "more local columns", columns: []string{"tenant_id", "user_id"}, referencedColumns: []string{"id"}, wantError: true},
		{name: "more referenced columns", columns: []string{"user_id"}, referencedColumns: []string{"tenant_id", "id"}, wantError: true},
		{name: "empty lists", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraintName := "fk_posts_users"
			schema := &Schema{Tables: []Table{{
				Name: "posts",
				ForeignKeys: []ForeignKey{{
					Name:              &constraintName,
					Columns:           tt.columns,
					ReferencedTable:   "users",
					ReferencedColumns: tt.referencedColumns,
					Source:            Source{Path: "schema.sql", Line: 12, Column: 5},
				}},
			}}}

			err := schema.Validate()
			if !tt.wantError {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}
			if !errors.Is(err, ErrForeignKeyCardinality) {
				t.Fatalf("Validate() error = %v, want %v", err, ErrForeignKeyCardinality)
			}

			var cardinality *ForeignKeyCardinalityError
			if !errors.As(err, &cardinality) {
				t.Fatalf("Validate() error type = %T, want *ForeignKeyCardinalityError", err)
			}
			if cardinality.ColumnCount != len(tt.columns) || cardinality.ReferencedColumnCount != len(tt.referencedColumns) {
				t.Fatalf(
					"cardinality counts = %d -> %d, want %d -> %d",
					cardinality.ColumnCount,
					cardinality.ReferencedColumnCount,
					len(tt.columns),
					len(tt.referencedColumns),
				)
			}
			for _, part := range []string{"fk_posts_users", "posts", "users", "schema.sql:12:5"} {
				if !strings.Contains(err.Error(), part) {
					t.Errorf("Validate() error %q does not contain %q", err, part)
				}
			}
		})
	}
}

func TestValidateUsesDeterministicErrorPrecedence(t *testing.T) {
	schema := &Schema{Tables: []Table{
		{
			Name:   "users",
			Source: Source{Path: "002.sql", Line: 1, Column: 1},
			Columns: []Column{
				{Name: "id"},
				{Name: "ID"},
			},
			ForeignKeys: []ForeignKey{{ReferencedTable: "tenants"}},
		},
		{Name: "USERS", Source: Source{Path: "001.sql", Line: 1, Column: 1}},
	}}

	if err := schema.Validate(); !errors.Is(err, ErrDuplicateTable) {
		t.Fatalf("Validate() error = %v, want duplicate table to take precedence", err)
	}
}

func TestSourceStringIncludesUnknownOrigin(t *testing.T) {
	if got, want := (Source{}).String(), "<unknown>:0:0"; got != want {
		t.Fatalf("Source.String() = %q, want %q", got, want)
	}
}
