// Package renderer defines output documents independently from persistence.
package renderer

// Document is one complete UTF-8 output file with LF line endings.
type Document struct {
	Path    string
	Content []byte
}
