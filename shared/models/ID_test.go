package models

import (
	"sort"
	"testing"
)

func TestNewID_IsVersion7(t *testing.T) {
	id := NewID()
	if id.Version() != 7 {
		t.Fatalf("expected version 7, got %d", id.Version())
	}
}

// TestNewID_SortsInGenerationOrder checks that IDs generated in sequence
// sort in generation order, even back-to-back within the same millisecond.
func TestNewID_SortsInGenerationOrder(t *testing.T) {
	const n = 1000

	generated := make([]string, n)
	for i := range generated {
		generated[i] = NewID().String()
	}

	sorted := make([]string, n)
	copy(sorted, generated)
	sort.Strings(sorted)

	for i := range generated {
		if generated[i] != sorted[i] {
			t.Fatalf("IDs are not sorted in generation order: index %d, generated %q, sorted %q", i, generated[i], sorted[i])
		}
	}
}
