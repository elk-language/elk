package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInt8Inspect(t *testing.T) {
	tests := map[string]struct {
		i    Int8
		want string
	}{
		"positive number": {
			i:    125,
			want: "125i8",
		},
		"negative number": {
			i:    -25,
			want: "-25i8",
		},
		"zero": {
			i:    0,
			want: "0i8",
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
