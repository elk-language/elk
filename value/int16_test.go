package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInt16Inspect(t *testing.T) {
	tests := map[string]struct {
		i    Int16
		want string
	}{
		"positive number": {
			i:    216,
			want: "216i16",
		},
		"negative number": {
			i:    -25,
			want: "-25i16",
		},
		"zero": {
			i:    0,
			want: "0i16",
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
