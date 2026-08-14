// Package markdown renders table documentation and the global README.
package markdown

import (
	"errors"

	"github.com/jeffersonc/sqlschemaview/internal/renderer"
	"github.com/jeffersonc/sqlschemaview/internal/schema"
)

// ErrNotImplemented marks renderer work deferred to its planned phase.
var ErrNotImplemented = errors.New("Markdown renderer is not implemented yet")

// Renderer creates all Markdown documents in memory.
type Renderer struct{}

func (Renderer) Render(*schema.Schema) ([]renderer.Document, error) {
	return nil, ErrNotImplemented
}
