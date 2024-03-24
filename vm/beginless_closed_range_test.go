package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestBeginlessClosedRangeContains(t *testing.T) {
	tests := map[string]struct {
		r        *value.BeginlessClosedRange
		val      value.Value
		contains bool
		err      value.Value
	}{
		"include int in the middle": {
			r: value.NewBeginlessClosedRange(
				value.SmallInt(5),
			),
			val:      value.SmallInt(3),
			contains: true,
		},
		"include int equal to end": {
			r: value.NewBeginlessClosedRange(
				value.SmallInt(10),
			),
			val:      value.SmallInt(10),
			contains: true,
		},
		"not include int greater than end": {
			r: value.NewBeginlessClosedRange(
				value.SmallInt(10),
			),
			val:      value.SmallInt(11),
			contains: false,
		},
		"include float": {
			r: value.NewBeginlessClosedRange(
				value.SmallInt(10),
			),
			val:      value.Float(5.7),
			contains: true,
		},
		"throw when incomparable value": {
			r: value.NewBeginlessClosedRange(
				value.SmallInt(10),
			),
			val: value.String("foo"),
			err: value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::String`",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			contains, err := vm.BeginlessClosedRangeContains(v, tc.r, tc.val)
			if diff := cmp.Diff(tc.err, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tc.contains, contains, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBeginlessClosedRangeEqual(t *testing.T) {
	tests := map[string]struct {
		r     *value.BeginlessClosedRange
		other *value.BeginlessClosedRange
		equal bool
		err   value.Value
	}{
		"two identical ranges": {
			r: value.NewBeginlessClosedRange(
				value.SmallInt(10),
			),
			other: value.NewBeginlessClosedRange(
				value.SmallInt(10),
			),
			equal: true,
		},
		"different end": {
			r: value.NewBeginlessClosedRange(
				value.SmallInt(11),
			),
			other: value.NewBeginlessClosedRange(
				value.SmallInt(10),
			),
			equal: false,
		},
		"Two ranges with the same values of different types": {
			r: value.NewBeginlessClosedRange(
				value.SmallInt(10),
			),
			other: value.NewBeginlessClosedRange(
				value.Float(10),
			),
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.BeginlessClosedRangeEqual(v, tc.r, tc.other)
			if diff := cmp.Diff(tc.err, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tc.equal, equal, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
