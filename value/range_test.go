package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestRangeInspect(t *testing.T) {
	tests := map[string]struct {
		r    *value.Range
		want string
	}{
		"beginless inclusive": {
			r:    value.NewRange(nil, value.SmallInt(6), false),
			want: `..6`,
		},
		"beginless exclusive": {
			r:    value.NewRange(nil, value.SmallInt(5), true),
			want: `...5`,
		},
		"endless inclusive": {
			r:    value.NewRange(value.SmallInt(3), nil, false),
			want: `3..`,
		},
		"endless exclusive": {
			r:    value.NewRange(value.SmallInt(3), nil, true),
			want: `3...`,
		},
		"inclusive": {
			r:    value.NewRange(value.Float(10.5), value.Float(38.2), false),
			want: `10.5..38.2`,
		},
		"exclusive": {
			r:    value.NewRange(value.SmallInt(3), value.SmallInt(10), true),
			want: `3...10`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.r.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
