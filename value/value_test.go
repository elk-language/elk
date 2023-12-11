package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestValueToBool(t *testing.T) {
	tests := map[string]struct {
		val  value.Value
		want value.Bool
	}{
		"positive number to true": {
			val:  value.Float(5),
			want: value.True,
		},
		"negative number to true": {
			val:  value.Float(-5),
			want: value.True,
		},
		"zero to true": {
			val:  value.SmallInt(0),
			want: value.True,
		},
		"string to true": {
			val:  value.String("foo"),
			want: value.True,
		},
		"empty string to true": {
			val:  value.String(""),
			want: value.True,
		},
		"true to true": {
			val:  value.True,
			want: value.True,
		},
		"nil to false": {
			val:  value.Nil,
			want: value.False,
		},
		"false to false": {
			val:  value.False,
			want: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.ToBool(tc.val)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestValue_InspectSlice(t *testing.T) {
	tests := map[string]struct {
		val  []value.Value
		want string
	}{
		"nil slice": {
			val:  nil,
			want: "[]",
		},
		"empty slice": {
			val:  make([]value.Value, 0),
			want: "[]",
		},
		"with values": {
			val: []value.Value{
				value.SmallInt(5),
				value.Float(10.5),
				value.String("foo"),
				value.Char('a'),
				value.ToSymbol("bar"),
			},
			want: `[5, 10.5, "foo", c"a", :bar]`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.InspectSlice(tc.val)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestValueToNotBool(t *testing.T) {
	tests := map[string]struct {
		val  value.Value
		want value.Bool
	}{
		"positive number to false": {
			val:  value.Float(5),
			want: value.False,
		},
		"negative number to false": {
			val:  value.Float(-5),
			want: value.False,
		},
		"zero to false": {
			val:  value.SmallInt(0),
			want: value.False,
		},
		"string to false": {
			val:  value.String("foo"),
			want: value.False,
		},
		"empty string to false": {
			val:  value.String(""),
			want: value.False,
		},
		"true to false": {
			val:  value.True,
			want: value.False,
		},
		"nil to true": {
			val:  value.Nil,
			want: value.True,
		},
		"false to true": {
			val:  value.False,
			want: value.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.ToNotBool(tc.val)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestValueTruthy(t *testing.T) {
	tests := map[string]struct {
		val  value.Value
		want bool
	}{
		"positive number to true": {
			val:  value.Float(5),
			want: true,
		},
		"negative number to true": {
			val:  value.Float(-5),
			want: true,
		},
		"zero to true": {
			val:  value.SmallInt(0),
			want: true,
		},
		"string to true": {
			val:  value.String("foo"),
			want: true,
		},
		"empty string to true": {
			val:  value.String(""),
			want: true,
		},
		"true to true": {
			val:  value.True,
			want: true,
		},
		"nil to false": {
			val:  value.Nil,
			want: false,
		},
		"false to false": {
			val:  value.False,
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.Truthy(tc.val)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestValueFalsy(t *testing.T) {
	tests := map[string]struct {
		val  value.Value
		want bool
	}{
		"positive number to false": {
			val:  value.Float(5),
			want: false,
		},
		"negative number to false": {
			val:  value.Float(-5),
			want: false,
		},
		"zero to false": {
			val:  value.SmallInt(0),
			want: false,
		},
		"string to false": {
			val:  value.String("foo"),
			want: false,
		},
		"empty string to false": {
			val:  value.String(""),
			want: false,
		},
		"true to false": {
			val:  value.True,
			want: false,
		},
		"nil to true": {
			val:  value.Nil,
			want: true,
		},
		"false to true": {
			val:  value.False,
			want: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.Falsy(tc.val)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
