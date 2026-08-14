package input

import "testing"

func TestFileRetainsSourcePath(t *testing.T) {
	file := File{Path: "migrations/001.sql", Content: []byte("SELECT 1;")}
	if file.Path != "migrations/001.sql" {
		t.Fatalf("File.Path = %q", file.Path)
	}
}
