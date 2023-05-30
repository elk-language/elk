package position

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
	left := &Position{
		StartByte:  45,
		ByteLength: 3,
		Line:       2,
		Column:     31,
	}
	right := &Position{
		StartByte:  54,
		ByteLength: 6,
		Line:       3,
		Column:     2,
	}
	want := &Position{
		StartByte:  45,
		ByteLength: 15,
		Line:       2,
		Column:     31,
	}
	got := left.Join(right)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(diff)
	}

	right = nil
	want = left
	got = left.Join(right)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(diff)
	}
}

func TestPositionValid(t *testing.T) {
	tests := []*Position{
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
	tests := []*Position{
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

func TestPositionJoinLastElement(t *testing.T) {
	tests := map[string]struct {
		left  *Position
		right []Interface
		want  *Position
	}{
		"returns left position if the collection is empty": {
			left: &Position{
				StartByte:  -2,
				ByteLength: 3,
				Line:       2,
				Column:     31,
			},
			right: nil,
			want: &Position{
				StartByte:  -2,
				ByteLength: 3,
				Line:       2,
				Column:     31,
			},
		},
		"joins the left expression with the last one in the collection": {
			left: &Position{
				StartByte:  45,
				ByteLength: 3,
				Line:       2,
				Column:     31,
			},
			right: []Interface{
				&Position{
					StartByte:  2,
					ByteLength: 26,
					Line:       61,
					Column:     8,
				},
				&Position{
					StartByte:  54,
					ByteLength: 6,
					Line:       3,
					Column:     2,
				},
			},
			want: &Position{
				StartByte:  45,
				ByteLength: 15,
				Line:       2,
				Column:     31,
			},
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
				&Position{
					StartByte:  2,
					ByteLength: 26,
					Line:       61,
					Column:     8,
				},
				&Position{
					StartByte:  54,
					ByteLength: 6,
					Line:       3,
					Column:     2,
				},
			},
			want: &Position{
				StartByte:  54,
				ByteLength: 6,
				Line:       3,
				Column:     2,
			},
		},
	}

	for _, testData := range tests {
		got := OfLastElement(testData.input)
		if diff := cmp.Diff(testData.want, got); diff != "" {
			t.Fatal(diff)
		}
	}
}
