package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestStringBuilder(t *testing.T) {
	tests := map[string]struct {
		s    *value.StringBuilder
		want value.String
	}{
		"build empty string": {
			s:    &value.StringBuilder{},
			want: value.String(""),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.s.ToString()

			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
