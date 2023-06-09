package position

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPositionHumanString(t *testing.T) {
	pos := New(45, 3, 2, 31)
	want := "2:31"
	got := pos.HumanString()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(diff)
	}
}

func TestPositionJoin(t *testing.T) {
	tests := map[string]struct {
		left  *Position
		right *Position
		want  *Position
	}{
		"return left when right is nil": {
			left:  New(45, 3, 2, 31),
			right: nil,
			want:  New(45, 3, 2, 31),
		},
		"return right when left is nil": {
			left:  nil,
			right: New(45, 3, 2, 31),
			want:  New(45, 3, 2, 31),
		},
		"return joined position": {
			left:  New(45, 3, 2, 31),
			right: New(54, 6, 3, 2),
			want:  New(45, 15, 2, 31),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Join(tc.right)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestPositionValid(t *testing.T) {
	tests := []*Position{
		New(45, 3, 2, 31),
		New(0, 1, 1, 1),
	}

	for _, pos := range tests {
		if !pos.Valid() {
			t.Fatalf("position should be valid: %#v", pos)
		}
	}
}

func TestPositionInvalid(t *testing.T) {
	tests := []*Position{
		New(-2, 3, 2, 31),
		New(5, 0, 2, 8),
		New(5, 4, 0, 8),
		New(5, 4, 9, 0),
	}

	for _, pos := range tests {
		if pos.Valid() {
			t.Fatalf("position should be invalid: %#v", pos)
		}
	}
}

func TestPositionJoinLastElement(t *testing.T) {
	tests := map[string]struct {
		left  *Position
		right []Interface
		want  *Position
	}{
		"returns left position if the collection is empty": {
			left:  New(-2, 3, 2, 31),
			right: nil,
			want:  New(-2, 3, 2, 31),
		},
		"joins the left expression with the last one in the collection": {
			left: New(45, 3, 2, 31),
			right: []Interface{
				New(2, 26, 61, 8),
				New(54, 6, 3, 2),
			},
			want: New(45, 15, 2, 31),
		},
	}

	for _, testData := range tests {
		got := JoinLastElement(testData.left, testData.right)
		if diff := cmp.Diff(testData.want, got); diff != "" {
			t.Fatal(diff)
		}
	}
}

func TestPositionOfLastElement(t *testing.T) {
	tests := map[string]struct {
		input []Interface
		want  *Position
	}{
		"returns nil if input is empty": {
			input: nil,
			want:  nil,
		},
		"joins the left expression with the last one in the collection": {
			input: []Interface{
				New(2, 26, 61, 8),
				New(54, 6, 3, 2),
			},
			want: New(54, 6, 3, 2),
		},
	}

	for _, testData := range tests {
		got := OfLastElement(testData.input)
		if diff := cmp.Diff(testData.want, got); diff != "" {
			t.Fatal(diff)
		}
	}
}
