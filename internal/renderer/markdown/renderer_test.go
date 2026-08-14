package markdown

import (
	"errors"
	"testing"

	"github.com/jeffersonc/sqlschemaview/internal/schema"
)

func TestRendererIsExplicitlyPending(t *testing.T) {
	_, err := (Renderer{}).Render(&schema.Schema{})
	if !errors.Is(err, ErrNotImplemented) {
		t.Fatalf("Render() error = %v, want %v", err, ErrNotImplemented)
	}
}
