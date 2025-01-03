package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestValueToBool(t *testing.T) {
	tests := map[string]struct {
		val  value.Value
		want value.Value
	}{
		"positive number to true": {
			val:  value.Float(5).ToValue(),
			want: value.True,
		},
		"negative number to true": {
			val:  value.Float(-5).ToValue(),
			want: value.True,
		},
		"zero to true": {
			val:  value.SmallInt(0).ToValue(),
			want: value.True,
		},
		"string to true": {
			val:  value.Ref(value.String("foo")),
			want: value.True,
		},
		"empty string to true": {
			val:  value.Ref(value.String("")),
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
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
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
				value.SmallInt(5).ToValue(),
				value.Float(10.5).ToValue(),
				value.Ref(value.String("foo")),
				value.Char('a').ToValue(),
				value.ToSymbol("bar").ToValue(),
			},
			want: "[5, 10.5, \"foo\", `a`, :bar]",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.InspectSlice(tc.val)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestValueToNotBool(t *testing.T) {
	tests := map[string]struct {
		val  value.Value
		want value.Value
	}{
		"positive number to false": {
			val:  value.Float(5).ToValue(),
			want: value.False,
		},
		"negative number to false": {
			val:  value.Float(-5).ToValue(),
			want: value.False,
		},
		"zero to false": {
			val:  value.SmallInt(0).ToValue(),
			want: value.False,
		},
		"string to false": {
			val:  value.Ref(value.String("foo")),
			want: value.False,
		},
		"empty string to false": {
			val:  value.Ref(value.String("")),
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
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
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
			val:  value.Float(5).ToValue(),
			want: true,
		},
		"negative number to true": {
			val:  value.Float(-5).ToValue(),
			want: true,
		},
		"zero to true": {
			val:  value.SmallInt(0).ToValue(),
			want: true,
		},
		"string to true": {
			val:  value.Ref(value.String("foo")),
			want: true,
		},
		"empty string to true": {
			val:  value.Ref(value.String("")),
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
				t.Fatal(diff)
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
			val:  value.Float(5).ToValue(),
			want: false,
		},
		"negative number to false": {
			val:  value.Float(-5).ToValue(),
			want: false,
		},
		"zero to false": {
			val:  value.SmallInt(0).ToValue(),
			want: false,
		},
		"string to false": {
			val:  value.Ref(value.String("foo")),
			want: false,
		},
		"empty string to false": {
			val:  value.Ref(value.String("")),
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
				t.Fatal(diff)
			}
		})
	}
}

func TestValue_InstanceOf(t *testing.T) {
	tests := map[string]struct {
		val   value.Value
		class *value.Class
		want  bool
	}{
		"true for direct instance": {
			val:   value.Float(5).ToValue(),
			class: value.FloatClass,
			want:  true,
		},
		"false for another class's instance": {
			val:   value.Float(5).ToValue(),
			class: value.IntClass,
			want:  false,
		},
		"false for superclass": {
			val:   value.Float(5).ToValue(),
			class: value.ObjectClass,
			want:  false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.InstanceOf(tc.val, tc.class)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestValue_ClassIsA(t *testing.T) {
	tests := map[string]struct {
		val   value.Value
		class *value.Class
		want  bool
	}{
		"true for direct instance": {
			val:   value.Float(5).ToValue(),
			class: value.FloatClass,
			want:  true,
		},
		"false for another class's instance": {
			val:   value.Float(5).ToValue(),
			class: value.IntClass,
			want:  false,
		},
		"true for superclass": {
			val:   value.Float(5).ToValue(),
			class: value.ObjectClass,
			want:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.IsA(tc.val, tc.class)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestValue_MixinIsA(t *testing.T) {
	tests := map[string]struct {
		val   value.Value
		mixin *value.Mixin
		want  bool
	}{
		"true for direct mixin": {
			val:   value.Ref(&value.HashMap{}),
			mixin: value.MapMixin,
			want:  true,
		},
		"true for indirect mixin": {
			val:   value.Ref(&value.HashMap{}),
			mixin: value.RecordMixin,
			want:  true,
		},
		"false for invalid mixin": {
			val:   value.Ref(&value.HashMap{}),
			mixin: value.ListMixin,
			want:  false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := value.IsA(tc.val, tc.mixin)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
