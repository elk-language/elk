package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestRightOpenRangeContains(t *testing.T) {
	tests := map[string]struct {
		r        *value.RightOpenRange
		val      value.Value
		contains bool
		err      value.Value
	}{
		"include int in the middle": {
			r: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(10),
			),
			val:      value.SmallInt(5),
			contains: true,
		},
		"include int equal to start": {
			r: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(10),
			),
			val:      value.SmallInt(3),
			contains: true,
		},
		"not include int equal to end": {
			r: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(10),
			),
			val:      value.SmallInt(10),
			contains: false,
		},
		"not include int lesser than start": {
			r: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(10),
			),
			val:      value.SmallInt(2),
			contains: false,
		},
		"not include int greater than end": {
			r: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(10),
			),
			val:      value.SmallInt(11),
			contains: false,
		},
		"include float": {
			r: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(10),
			),
			val:      value.Float(5.7),
			contains: true,
		},
		"throw when incomparable value": {
			r: value.NewRightOpenRange(
				value.SmallInt(3),
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
			contains, err := vm.RightOpenRangeContains(v, tc.r, tc.val)
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

func TestRightOpenRangeEqual(t *testing.T) {
	tests := map[string]struct {
		r     *value.RightOpenRange
		other *value.RightOpenRange
		equal bool
		err   value.Value
	}{
		"two identical ranges": {
			r: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(10),
			),
			other: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(10),
			),
			equal: true,
		},
		"different end": {
			r: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(11),
			),
			other: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(10),
			),
			equal: false,
		},
		"different start": {
			r: value.NewRightOpenRange(
				value.SmallInt(4),
				value.SmallInt(10),
			),
			other: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(10),
			),
			equal: false,
		},
		"different start and end": {
			r: value.NewRightOpenRange(
				value.SmallInt(4),
				value.SmallInt(10),
			),
			other: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(15),
			),
			equal: false,
		},
		"Two ranges with the same values of different types": {
			r: value.NewRightOpenRange(
				value.SmallInt(3),
				value.SmallInt(10),
			),
			other: value.NewRightOpenRange(
				value.Float(3),
				value.Float(10),
			),
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.RightOpenRangeEqual(v, tc.r, tc.other)
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

func TestRightOpenRangeIteratorNext(t *testing.T) {
	tests := map[string]struct {
		i     *value.RightOpenRangeIterator
		after *value.RightOpenRangeIterator
		want  value.Value
		err   value.Value
	}{
		"empty": {
			i: value.NewRightOpenRangeIterator(
				value.NewRightOpenRange(
					value.SmallInt(3),
					value.SmallInt(3),
				),
			),
			err: value.ToSymbol("stop_iteration"),
		},
		"with two elements, first iteration": {
			i: value.NewRightOpenRangeIterator(
				value.NewRightOpenRange(
					value.SmallInt(1),
					value.SmallInt(3),
				),
			),
			after: value.NewRightOpenRangeIteratorWithCurrentElement(
				value.NewRightOpenRange(
					value.SmallInt(1),
					value.SmallInt(3),
				),
				value.SmallInt(2),
			),
			want: value.SmallInt(1),
		},
		"with two elements, second iteration": {
			i: value.NewRightOpenRangeIteratorWithCurrentElement(
				value.NewRightOpenRange(
					value.SmallInt(1),
					value.SmallInt(3),
				),
				value.SmallInt(2),
			),
			after: value.NewRightOpenRangeIteratorWithCurrentElement(
				value.NewRightOpenRange(
					value.SmallInt(1),
					value.SmallInt(3),
				),
				value.SmallInt(3),
			),
			want: value.SmallInt(2),
		},
		"with two elements, third iteration": {
			i: value.NewRightOpenRangeIteratorWithCurrentElement(
				value.NewRightOpenRange(
					value.SmallInt(1),
					value.SmallInt(3),
				),
				value.SmallInt(3),
			),
			err: value.ToSymbol("stop_iteration"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			opts := comparer.Options()
			got, err := vm.RightOpenRangeIteratorNext(v, tc.i)
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if tc.err != nil {
				return
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.after, tc.i, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
