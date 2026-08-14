// Package parser extracts the supported structural SQL subset into Schema IR.
package parser

import (
	"context"

	"github.com/jeffersonc/sqlschemaview/internal/input"
	"github.com/jeffersonc/sqlschemaview/internal/schema"
)

// Parser is the SQL extraction entry point.
type Parser struct{}

// Extract will tokenize all files before extracting supported schema statements.
func (Parser) Extract(context.Context, []input.File) (*schema.Schema, error) {
	return nil, ErrNotImplemented
}
