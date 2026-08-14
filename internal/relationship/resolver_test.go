package relationship

import (
	"errors"
	"testing"

	"github.com/jeffersonc/sqlschemaview/internal/schema"
)

func TestResolverIsExplicitlyPending(t *testing.T) {
	err := (Resolver{}).Resolve(&schema.Schema{})
	if !errors.Is(err, ErrNotImplemented) {
		t.Fatalf("Resolve() error = %v, want %v", err, ErrNotImplemented)
	}
}
