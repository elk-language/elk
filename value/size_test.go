package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestSize_Bytes(t *testing.T) {
	tests := map[string]struct {
		s    value.Size
		want value.Float
	}{
		"returns self as float": {
			s:    value.Megabyte,
			want: value.Float(1024 * 1024),
		},
		"handles more complex values": {
			s:    2*value.Megabyte + 3*value.Kilobyte,
			want: value.Float(28475845.248754),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.Bytes()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
