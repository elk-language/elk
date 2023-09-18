package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRangeInspect(t *testing.T) {
	tests := map[string]struct {
		r    *Range
		want string
	}{
		"beginless inclusive": {
			r:    NewRange(nil, SmallInt(6), false),
			want: `..6`,
		},
		"beginless exclusive": {
			r:    NewRange(nil, SmallInt(5), true),
			want: `...5`,
		},
		"endless inclusive": {
			r:    NewRange(SmallInt(3), nil, false),
			want: `3..`,
		},
		"endless exclusive": {
			r:    NewRange(SmallInt(3), nil, true),
			want: `3...`,
		},
		"inclusive": {
			r:    NewRange(Float(10.5), Float(38.2), false),
			want: `10.5..38.2`,
		},
		"exclusive": {
			r:    NewRange(SmallInt(3), SmallInt(10), true),
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
