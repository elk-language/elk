package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUInt64Inspect(t *testing.T) {
	tests := map[string]struct {
		i    UInt64
		want string
	}{
		"positive number": {
			i:    216,
			want: "216u64",
		},
		"zero": {
			i:    0,
			want: "0u64",
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
