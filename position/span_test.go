package position

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSpanJoin(t *testing.T) {
	tests := map[string]struct {
		left  *Span
		right *Span
		want  *Span
	}{
		"return left when right is nil": {
			left:  NewSpan(New(45, 2, 31), New(48, 2, 34)),
			right: nil,
			want:  NewSpan(New(45, 2, 31), New(48, 2, 34)),
		},
		"return right when left is nil": {
			left:  nil,
			right: NewSpan(New(45, 2, 31), New(48, 2, 34)),
			want:  NewSpan(New(45, 2, 31), New(48, 2, 34)),
		},
		"return joined": {
			left:  NewSpan(New(45, 2, 31), New(48, 2, 34)),
			right: NewSpan(New(49, 3, 1), New(53, 3, 5)),
			want:  NewSpan(New(45, 2, 31), New(53, 3, 5)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Join(tc.right)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestSpanJoinLastElement(t *testing.T) {
	tests := map[string]struct {
		left  *Span
		right []SpanInterface
		want  *Span
	}{
		"returns left if the collection is empty": {
			left:  NewSpan(New(45, 2, 31), New(48, 2, 34)),
			right: nil,
			want:  NewSpan(New(45, 2, 31), New(48, 2, 34)),
		},
		"joins the left expression with the last one in the collection": {
			left: NewSpan(New(45, 2, 31), New(48, 2, 34)),
			right: []SpanInterface{
				NewSpan(New(51, 3, 1), New(53, 3, 3)),
				NewSpan(New(54, 3, 4), New(60, 3, 10)),
			},
			want: NewSpan(New(45, 2, 31), New(60, 3, 10)),
		},
	}

	for _, testData := range tests {
		got := JoinSpanOfLastElement(testData.left, testData.right)
		if diff := cmp.Diff(testData.want, got); diff != "" {
			t.Fatal(diff)
		}
	}
}

func TestSpanOfLastElement(t *testing.T) {
	tests := map[string]struct {
		input []SpanInterface
		want  *Span
	}{
		"returns nil if input is empty": {
			input: nil,
			want:  nil,
		},
		"joins the left expression with the last one in the collection": {
			input: []SpanInterface{
				NewSpan(New(51, 3, 1), New(53, 3, 3)),
				NewSpan(New(54, 3, 4), New(60, 3, 10)),
			},
			want: NewSpan(New(54, 3, 4), New(60, 3, 10)),
		},
	}

	for _, testData := range tests {
		got := SpanOfLastElement(testData.input)
		if diff := cmp.Diff(testData.want, got); diff != "" {
			t.Fatal(diff)
		}
	}
}
