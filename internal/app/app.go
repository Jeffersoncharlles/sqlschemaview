// Package app coordinates the SQL Schema View application workflow.
package app

import (
	"context"
	"fmt"

	"github.com/jeffersonc/sqlschemaview/internal/input"
	"github.com/jeffersonc/sqlschemaview/internal/renderer"
	"github.com/jeffersonc/sqlschemaview/internal/schema"
)

// Config contains one complete CLI execution configuration.
type Config struct {
	InputPath string
	OutputDir string
}

// Result summarizes a successful documentation generation.
type Result struct {
	Tables      int
	Columns     int
	PrimaryKeys int
	ForeignKeys int
	Documents   []renderer.Document
}

// Loader discovers and reads supported SQL inputs.
type Loader interface {
	Load(context.Context, string) ([]input.File, error)
}

// Extractor converts SQL sources into the domain Schema IR.
type Extractor interface {
	Extract(context.Context, []input.File) (*schema.Schema, error)
}

// Resolver validates and resolves foreign-key relationships.
type Resolver interface {
	Resolve(*schema.Schema) error
}

// Renderer creates all output documents without writing to disk.
type Renderer interface {
	Render(*schema.Schema) ([]renderer.Document, error)
}

// Writer persists an already-rendered document set.
type Writer interface {
	Write(context.Context, string, []renderer.Document) error
}

// App owns the dependencies required by one execution.
type App struct {
	Loader    Loader
	Extractor Extractor
	Resolver  Resolver
	Renderer  Renderer
	Writer    Writer
}

// Run will become the single orchestration point as the implementation advances.
func (a App) Run(context.Context, Config) (Result, error) {
	return Result{}, fmt.Errorf("application workflow is not implemented yet")
}
