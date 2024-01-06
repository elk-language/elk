package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestTuple_Inspect(t *testing.T) {
	tests := map[string]struct {
		t    *value.Tuple
		want string
	}{
		"empty": {
			t:    &value.Tuple{},
			want: "%[]",
		},
		"with one element": {
			t:    &value.Tuple{value.SmallInt(3)},
			want: `%[3]`,
		},
		"with elements": {
			t:    &value.Tuple{value.SmallInt(3), value.String("foo")},
			want: `%[3, "foo"]`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.t.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestTuple_Expand(t *testing.T) {
	tests := map[string]struct {
		t    *value.Tuple
		new  int
		want *value.Tuple
	}{
		"add 3 to an empty list": {
			t:    &value.Tuple{},
			new:  3,
			want: &value.Tuple{value.Nil, value.Nil, value.Nil},
		},
		"add 0 to an empty list": {
			t:    &value.Tuple{},
			new:  0,
			want: &value.Tuple{},
		},
		"add 2 to a filled list": {
			t:    &value.Tuple{value.SmallInt(-3), value.Float(10.5)},
			new:  2,
			want: &value.Tuple{value.SmallInt(-3), value.Float(10.5), value.Nil, value.Nil},
		},
		"add 0 to a filled list": {
			t:    &value.Tuple{value.SmallInt(-3), value.Float(10.5)},
			new:  0,
			want: &value.Tuple{value.SmallInt(-3), value.Float(10.5)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.t.Expand(tc.new)
			if diff := cmp.Diff(tc.want, tc.t); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestTuple_Append(t *testing.T) {
	tests := map[string]struct {
		t    *value.Tuple
		val  value.Value
		want *value.Tuple
	}{
		"append to an empty list": {
			t:    &value.Tuple{},
			val:  value.SmallInt(3),
			want: &value.Tuple{value.SmallInt(3)},
		},
		"append to a filled list": {
			t:    &value.Tuple{value.SmallInt(3)},
			val:  value.String("foo"),
			want: &value.Tuple{value.SmallInt(3), value.String("foo")},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.t.Append(tc.val)
			if diff := cmp.Diff(tc.want, tc.t); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
