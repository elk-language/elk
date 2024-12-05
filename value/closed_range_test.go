package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestClosedRangeIterator_Inspect(t *testing.T) {
	tests := map[string]struct {
		i    *value.ClosedRangeIterator
		want string
	}{
		"initial": {
			i: value.NewClosedRangeIterator(
				value.NewClosedRange(value.SmallInt(3).ToValue(), value.SmallInt(10).ToValue()),
			),
			want: "Std::ClosedRange::Iterator{range: 3...10, current_element: 3}",
		},
		"current element": {
			i: value.NewClosedRangeIteratorWithCurrentElement(
				value.NewClosedRange(value.SmallInt(3).ToValue(), value.SmallInt(10).ToValue()),
				value.SmallInt(7).ToValue(),
			),
			want: "Std::ClosedRange::Iterator{range: 3...10, current_element: 7}",
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
