package position

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPositionString(t *testing.T) {
	pos := New(45, 2, 31)
	want := "2:31"
	got := pos.String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(diff)
	}
}

func TestPositionValid(t *testing.T) {
	tests := []*Position{
		New(45, 2, 31),
		New(0, 1, 1),
	}

	for _, pos := range tests {
		if !pos.Valid() {
			t.Fatalf("position should be valid: %#v", pos)
		}
	}
}

func TestPositionInvalid(t *testing.T) {
	tests := []*Position{
		New(-2, 2, 31),
		New(5, 0, 8),
		New(5, 9, 0),
	}

	for _, pos := range tests {
		if pos.Valid() {
			t.Fatalf("position should be invalid: %#v", pos)
		}
	}
}
