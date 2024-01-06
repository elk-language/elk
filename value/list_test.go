package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestList_Inspect(t *testing.T) {
	tests := map[string]struct {
		l    *value.List
		want string
	}{
		"empty": {
			l:    &value.List{},
			want: "[]",
		},
		"with one element": {
			l:    &value.List{value.SmallInt(3)},
			want: `[3]`,
		},
		"with elements": {
			l:    &value.List{value.SmallInt(3), value.String("foo")},
			want: `[3, "foo"]`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.l.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestList_Expand(t *testing.T) {
	tests := map[string]struct {
		l    *value.List
		new  int
		want *value.List
	}{
		"add 3 to an empty list": {
			l:    &value.List{},
			new:  3,
			want: &value.List{value.Nil, value.Nil, value.Nil},
		},
		"add 0 to an empty list": {
			l:    &value.List{},
			new:  0,
			want: &value.List{},
		},
		"add 2 to a filled list": {
			l:    &value.List{value.SmallInt(-3), value.Float(10.5)},
			new:  2,
			want: &value.List{value.SmallInt(-3), value.Float(10.5), value.Nil, value.Nil},
		},
		"add 0 to a filled list": {
			l:    &value.List{value.SmallInt(-3), value.Float(10.5)},
			new:  0,
			want: &value.List{value.SmallInt(-3), value.Float(10.5)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.l.Expand(tc.new)
			if diff := cmp.Diff(tc.want, tc.l); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestList_Append(t *testing.T) {
	tests := map[string]struct {
		l    *value.List
		val  value.Value
		want *value.List
	}{
		"append to an empty list": {
			l:    &value.List{},
			val:  value.SmallInt(3),
			want: &value.List{value.SmallInt(3)},
		},
		"append to a filled list": {
			l:    &value.List{value.SmallInt(3)},
			val:  value.String("foo"),
			want: &value.List{value.SmallInt(3), value.String("foo")},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.l.Append(tc.val)
			if diff := cmp.Diff(tc.want, tc.l); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
