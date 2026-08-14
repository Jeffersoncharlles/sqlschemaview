package schema

import "errors"

var (
	ErrNoTables             = errors.New("no tables found in SQL input")
	ErrDuplicateTable       = errors.New("duplicate table definition")
	ErrDuplicateColumn      = errors.New("duplicate column definition")
	ErrUnresolvedForeignKey = errors.New("foreign key references an unknown table or column")
)
