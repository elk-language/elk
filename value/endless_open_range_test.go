package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestEndlessOpenRangeIterator_Inspect(t *testing.T) {
	tests := map[string]struct {
		i    *value.EndlessOpenRangeIterator
		want string
	}{
		"initial": {
			i: value.NewEndlessOpenRangeIterator(
				value.NewEndlessOpenRange(value.SmallInt(3).ToValue()),
			),
			want: "Std::EndlessOpenRange::Iterator{range: 3<.., current_element: 3}",
		},
		"current element": {
			i: value.NewEndlessOpenRangeIteratorWithCurrentElement(
				value.NewEndlessOpenRange(value.SmallInt(3).ToValue()),
				value.SmallInt(7).ToValue(),
			),
			want: "Std::EndlessOpenRange::Iterator{range: 3<.., current_element: 7}",
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
