package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestInt64Inspect(t *testing.T) {
	tests := map[string]struct {
		i    value.Int64
		want string
	}{
		"positive number": {
			i:    216,
			want: "216i64",
		},
		"negative number": {
			i:    -25,
			want: "-25i64",
		},
		"zero": {
			i:    0,
			want: "0i64",
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
