package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValueToBool(t *testing.T) {
	tests := map[string]struct {
		val  Value
		want Bool
	}{
		"positive number to true": {
			val:  Float(5),
			want: True,
		},
		"negative number to true": {
			val:  Float(-5),
			want: True,
		},
		"zero to true": {
			val:  SmallInt(0),
			want: True,
		},
		"string to true": {
			val:  String("foo"),
			want: True,
		},
		"empty string to true": {
			val:  String(""),
			want: True,
		},
		"true to true": {
			val:  True,
			want: True,
		},
		"nil to false": {
			val:  Nil,
			want: False,
		},
		"false to false": {
			val:  False,
			want: False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := ToBool(tc.val)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestValueToNotBool(t *testing.T) {
	tests := map[string]struct {
		val  Value
		want Bool
	}{
		"positive number to false": {
			val:  Float(5),
			want: False,
		},
		"negative number to false": {
			val:  Float(-5),
			want: False,
		},
		"zero to false": {
			val:  SmallInt(0),
			want: False,
		},
		"string to false": {
			val:  String("foo"),
			want: False,
		},
		"empty string to false": {
			val:  String(""),
			want: False,
		},
		"true to false": {
			val:  True,
			want: False,
		},
		"nil to true": {
			val:  Nil,
			want: True,
		},
		"false to true": {
			val:  False,
			want: True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := ToNotBool(tc.val)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestValueTruthy(t *testing.T) {
	tests := map[string]struct {
		val  Value
		want bool
	}{
		"positive number to true": {
			val:  Float(5),
			want: true,
		},
		"negative number to true": {
			val:  Float(-5),
			want: true,
		},
		"zero to true": {
			val:  SmallInt(0),
			want: true,
		},
		"string to true": {
			val:  String("foo"),
			want: true,
		},
		"empty string to true": {
			val:  String(""),
			want: true,
		},
		"true to true": {
			val:  True,
			want: true,
		},
		"nil to false": {
			val:  Nil,
			want: false,
		},
		"false to false": {
			val:  False,
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := Truthy(tc.val)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestValueFalsy(t *testing.T) {
	tests := map[string]struct {
		val  Value
		want bool
	}{
		"positive number to false": {
			val:  Float(5),
			want: false,
		},
		"negative number to false": {
			val:  Float(-5),
			want: false,
		},
		"zero to false": {
			val:  SmallInt(0),
			want: false,
		},
		"string to false": {
			val:  String("foo"),
			want: false,
		},
		"empty string to false": {
			val:  String(""),
			want: false,
		},
		"true to false": {
			val:  True,
			want: false,
		},
		"nil to true": {
			val:  Nil,
			want: true,
		},
		"false to true": {
			val:  False,
			want: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := Falsy(tc.val)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
