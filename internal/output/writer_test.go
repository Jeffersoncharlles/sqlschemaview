package output

import (
	"context"
	"errors"
	"testing"
)

func TestWriterIsExplicitlyPending(t *testing.T) {
	err := (Writer{}).Write(context.Background(), t.TempDir(), nil)
	if !errors.Is(err, ErrNotImplemented) {
		t.Fatalf("Write() error = %v, want %v", err, ErrNotImplemented)
	}
}
