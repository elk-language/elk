package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestUInt16Inspect(t *testing.T) {
	tests := map[string]struct {
		i    value.UInt16
		want string
	}{
		"positive number": {
			i:    216,
			want: "216u16",
		},
		"zero": {
			i:    0,
			want: "0u16",
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
