package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInt32Inspect(t *testing.T) {
	tests := map[string]struct {
		i    Int32
		want string
	}{
		"positive number": {
			i:    216,
			want: "216i32",
		},
		"negative number": {
			i:    -25,
			want: "-25i32",
		},
		"zero": {
			i:    0,
			want: "0i32",
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
