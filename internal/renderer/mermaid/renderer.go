// Package mermaid renders ER diagrams from resolved Schema IR relationships.
package mermaid

import (
	"errors"

	"github.com/jeffersonc/sqlschemaview/internal/schema"
)

// ErrNotImplemented marks renderer work deferred to its planned phase.
var ErrNotImplemented = errors.New("Mermaid renderer is not implemented yet")

// Renderer sanitizes only display identifiers and fails on collisions.
type Renderer struct{}

func (Renderer) Render(*schema.Schema) (string, error) {
	return "", ErrNotImplemented
}
