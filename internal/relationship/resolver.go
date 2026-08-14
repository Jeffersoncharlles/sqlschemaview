// Package relationship validates foreign keys and derives navigable relationships.
package relationship

import (
	"errors"

	"github.com/jeffersonc/sqlschemaview/internal/schema"
)

// ErrNotImplemented marks relationship work deferred to its planned phase.
var ErrNotImplemented = errors.New("relationship resolver is not implemented yet")

// Resolver resolves all foreign keys only after every input has been parsed.
type Resolver struct{}

func (Resolver) Resolve(*schema.Schema) error {
	return ErrNotImplemented
}
