package parser

import "errors"

// ErrNotImplemented marks parser work intentionally deferred to its planned phase.
var ErrNotImplemented = errors.New("SQL parser is not implemented yet")

type tokenKind uint8

const (
	tokenUnknown tokenKind = iota
	tokenWord
	tokenIdentifier
	tokenString
	tokenSymbol
	tokenEOF
)

type token struct {
	kind   tokenKind
	value  string
	line   int
	column int
}
