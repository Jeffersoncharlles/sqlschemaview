// Package input discovers SQL files and loads their contents.
package input

// File is one UTF-8 SQL source loaded from the input path.
type File struct {
	Path    string
	Content []byte
}

// Loader will implement regular-file validation, recursive directory walking,
// symlink exclusion, and deterministic path ordering.
type Loader struct{}
