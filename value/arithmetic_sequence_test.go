package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestArithmeticSequenceInspect(t *testing.T) {
	tests := map[string]struct {
		a    *value.ArithmeticSequence
		want string
	}{
		"endless inclusive": {
			a:    value.NewArithmeticSequence(value.SmallInt(3), nil, value.SmallInt(2), false),
			want: `3..:2`,
		},
		"endless exclusive": {
			a:    value.NewArithmeticSequence(value.SmallInt(3), nil, value.SmallInt(2), true),
			want: `3...:2`,
		},
		"inclusive": {
			a:    value.NewArithmeticSequence(value.Float(10.5), value.Float(38.2), value.Float(2.5), false),
			want: `10.5..38.2:2.5`,
		},
		"exclusive": {
			a:    value.NewArithmeticSequence(value.SmallInt(3), value.SmallInt(10), value.SmallInt(2), true),
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
