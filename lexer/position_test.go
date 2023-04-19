package lexer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPositionHumanString(t *testing.T) {
	pos := Position{
		StartByte:  45,
		ByteLength: 3,
		Line:       2,
		Column:     31,
	}
	want := "2:31"
	got := pos.HumanString()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(diff)
	}
}

func TestPositionJoin(t *testing.T) {
	left := Position{
		StartByte:  45,
		ByteLength: 3,
		Line:       2,
		Column:     31,
	}
	right := Position{
		StartByte:  54,
		ByteLength: 6,
		Line:       3,
		Column:     2,
	}
	want := Position{
		StartByte:  45,
		ByteLength: 15,
		Line:       2,
		Column:     31,
	}
	got := left.Join(right)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(diff)
	}
}

func TestPositionValid(t *testing.T) {
	tests := []Position{
		{
			StartByte:  45,
			ByteLength: 3,
			Line:       2,
			Column:     31,
		},
		{
			StartByte:  0,
			ByteLength: 1,
			Line:       1,
			Column:     1,
		},
	}

	for _, pos := range tests {
		if !pos.Valid() {
			t.Fatalf("position should be valid: %#v", pos)
		}
	}
}

func TestPositionInvalid(t *testing.T) {
	tests := []Position{
		{
			StartByte:  -2,
			ByteLength: 3,
			Line:       2,
			Column:     31,
		},
		{
			StartByte:  5,
			ByteLength: 0,
			Line:       2,
			Column:     8,
		},
		{
			StartByte:  5,
			ByteLength: 4,
			Line:       0,
			Column:     8,
		},
		{
			StartByte:  5,
			ByteLength: 4,
			Line:       9,
			Column:     0,
		},
	}

	for _, pos := range tests {
		if pos.Valid() {
			t.Fatalf("position should be invalid: %#v", pos)
		}
	}
}
