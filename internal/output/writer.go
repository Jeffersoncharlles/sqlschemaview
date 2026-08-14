// Package output atomically persists complete rendered document sets.
package output

import (
	"context"
	"errors"

	"github.com/jeffersonc/sqlschemaview/internal/renderer"
)

// ErrNotImplemented marks output work deferred to its planned phase.
var ErrNotImplemented = errors.New("output writer is not implemented yet")

// Writer creates missing output directories without deleting unrelated files.
type Writer struct{}

func (Writer) Write(context.Context, string, []renderer.Document) error {
	return ErrNotImplemented
}
