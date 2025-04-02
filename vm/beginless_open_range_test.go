package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestBeginlessOpenRangeContains(t *testing.T) {
	tests := map[string]struct {
		r        *value.BeginlessOpenRange
		val      value.Value
		contains bool
		err      value.Value
	}{
		"include int in the middle": {
			r: value.NewBeginlessOpenRange(
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(5).ToValue(),
			contains: true,
		},
		"not include int equal to end": {
			r: value.NewBeginlessOpenRange(
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(10).ToValue(),
			contains: false,
		},
		"not include int greater than end": {
			r: value.NewBeginlessOpenRange(
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(11).ToValue(),
			contains: false,
		},
		"include float": {
			r: value.NewBeginlessOpenRange(
				value.SmallInt(10).ToValue(),
			),
			val:      value.Float(5.7).ToValue(),
			contains: true,
		},
		"throw when incomparable value": {
			r: value.NewBeginlessOpenRange(
				value.SmallInt(10).ToValue(),
			),
			val: value.Ref(value.String("foo")),
			err: value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::String`",
			)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			contains, err := vm.BeginlessOpenRangeContains(v, tc.r, tc.val)
			if diff := cmp.Diff(tc.err, err, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
			if !err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.contains, contains, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBeginlessOpenRangeEqual(t *testing.T) {
	tests := map[string]struct {
		r     *value.BeginlessOpenRange
		other *value.BeginlessOpenRange
		equal bool
		err   value.Value
	}{
		"two identical ranges": {
			r: value.NewBeginlessOpenRange(
				value.SmallInt(10).ToValue(),
			),
			other: value.NewBeginlessOpenRange(
				value.SmallInt(10).ToValue(),
			),
			equal: true,
		},
		"different end": {
			r: value.NewBeginlessOpenRange(
				value.SmallInt(11).ToValue(),
			),
			other: value.NewBeginlessOpenRange(
				value.SmallInt(10).ToValue(),
			),
			equal: false,
		},
		"Two ranges with the same values of different types": {
			r: value.NewBeginlessOpenRange(
				value.SmallInt(10).ToValue(),
			),
			other: value.NewBeginlessOpenRange(
				value.Float(10).ToValue(),
			),
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.BeginlessOpenRangeEqual(v, tc.r, tc.other)
			if diff := cmp.Diff(tc.err, err, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
			if !err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.equal, equal, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
