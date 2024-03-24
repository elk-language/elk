package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestEndlessClosedRangeIterator_Inspect(t *testing.T) {
	tests := map[string]struct {
		i    *value.EndlessClosedRangeIterator
		want string
	}{
		"initial": {
			i: value.NewEndlessClosedRangeIterator(
				value.NewEndlessClosedRange(value.SmallInt(3)),
			),
			want: "Std::EndlessClosedRange::Iterator{range: 3..., current_element: 3}",
		},
		"current element": {
			i: value.NewEndlessClosedRangeIteratorWithCurrentElement(
				value.NewEndlessClosedRange(value.SmallInt(3)),
				value.SmallInt(7),
			),
			want: "Std::EndlessClosedRange::Iterator{range: 3..., current_element: 7}",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.i.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
