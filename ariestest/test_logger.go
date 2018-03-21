package ariestest

import (
	"testing"

	"shanhu.io/aries"
)

type testLog struct {
	t *testing.T
}

func (l *testLog) Print(s string) { l.t.Log(s) }

// NewLogger creates a test logger.
func NewLogger(t *testing.T) *aries.Logger {
	return aries.NewLogger(&testLog{t})
}
