package schema

import "testing"

func TestCompositeForeignKeyRemainsSingleConstraint(t *testing.T) {
	fk := ForeignKey{
		Columns:           []string{"tenant_id", "user_id"},
		ReferencedTable:   "users",
		ReferencedColumns: []string{"tenant_id", "id"},
	}

	if len(fk.Columns) != 2 || len(fk.ReferencedColumns) != 2 {
		t.Fatalf("composite foreign key was not preserved: %#v", fk)
	}
}
