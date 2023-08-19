package object

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestArithmeticSequenceInspect(t *testing.T) {
	tests := map[string]struct {
		a    *ArithmeticSequence
		want string
	}{
		"endless inclusive": {
			a:    NewArithmeticSequence(SmallInt(3), nil, SmallInt(2), false),
			want: `3..:2`,
		},
		"endless exclusive": {
			a:    NewArithmeticSequence(SmallInt(3), nil, SmallInt(2), true),
			want: `3...:2`,
		},
		"inclusive": {
			a:    NewArithmeticSequence(Float(10.5), Float(38.2), Float(2.5), false),
			want: `10.5..38.2:2.5`,
		},
		"exclusive": {
			a:    NewArithmeticSequence(SmallInt(3), SmallInt(10), SmallInt(2), true),
			want: `3...10:2`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.a.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
