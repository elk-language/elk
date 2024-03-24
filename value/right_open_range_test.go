package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestRightOpenRangeIterator_Inspect(t *testing.T) {
	tests := map[string]struct {
		i    *value.RightOpenRangeIterator
		want string
	}{
		"initial": {
			i: value.NewRightOpenRangeIterator(
				value.NewRightOpenRange(value.SmallInt(3), value.SmallInt(10)),
			),
			want: "Std::RightOpenRange::Iterator{range: 3..<10, current_element: 3}",
		},
		"current element": {
			i: value.NewRightOpenRangeIteratorWithCurrentElement(
				value.NewRightOpenRange(value.SmallInt(3), value.SmallInt(10)),
				value.SmallInt(7),
			),
			want: "Std::RightOpenRange::Iterator{range: 3..<10, current_element: 7}",
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
