package position

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLocationString(t *testing.T) {
	tests := map[string]struct {
		in   *Location
		want string
	}{
		"return empty string when nil": {
			in:   nil,
			want: "",
		},
		"return correct string": {
			in:   NewLocation("/foo/bar.elk", 35, 10, 3, 5),
			want: "/foo/bar.elk:3:5",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.String()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
