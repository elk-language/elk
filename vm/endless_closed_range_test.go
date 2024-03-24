package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestEndlessClosedRangeContains(t *testing.T) {
	tests := map[string]struct {
		r        *value.EndlessClosedRange
		val      value.Value
		contains bool
		err      value.Value
	}{
		"include int in the middle": {
			r: value.NewEndlessClosedRange(
				value.SmallInt(3),
			),
			val:      value.SmallInt(5),
			contains: true,
		},
		"include int equal to start": {
			r: value.NewEndlessClosedRange(
				value.SmallInt(3),
			),
			val:      value.SmallInt(3),
			contains: true,
		},
		"not include int lesser than start": {
			r: value.NewEndlessClosedRange(
				value.SmallInt(3),
			),
			val:      value.SmallInt(2),
			contains: false,
		},
		"include float": {
			r: value.NewEndlessClosedRange(
				value.SmallInt(3),
			),
			val:      value.Float(5.7),
			contains: true,
		},
		"throw when incomparable value": {
			r: value.NewEndlessClosedRange(
				value.SmallInt(3),
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
			contains, err := vm.EndlessClosedRangeContains(v, tc.r, tc.val)
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

func TestEndlessClosedRangeEqual(t *testing.T) {
	tests := map[string]struct {
		r     *value.EndlessClosedRange
		other *value.EndlessClosedRange
		equal bool
		err   value.Value
	}{
		"two identical ranges": {
			r: value.NewEndlessClosedRange(
				value.SmallInt(3),
			),
			other: value.NewEndlessClosedRange(
				value.SmallInt(3),
			),
			equal: true,
		},
		"different start": {
			r: value.NewEndlessClosedRange(
				value.SmallInt(4),
			),
			other: value.NewEndlessClosedRange(
				value.SmallInt(3),
			),
			equal: false,
		},
		"Two ranges with the same values of different types": {
			r: value.NewEndlessClosedRange(
				value.SmallInt(3),
			),
			other: value.NewEndlessClosedRange(
				value.Float(3),
			),
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.EndlessClosedRangeEqual(v, tc.r, tc.other)
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

func TestEndlessClosedRangeIteratorNext(t *testing.T) {
	tests := map[string]struct {
		i     *value.EndlessClosedRangeIterator
		after *value.EndlessClosedRangeIterator
		want  value.Value
		err   value.Value
	}{
		"first iteration": {
			i: value.NewEndlessClosedRangeIterator(
				value.NewEndlessClosedRange(
					value.SmallInt(2),
				),
			),
			after: value.NewEndlessClosedRangeIteratorWithCurrentElement(
				value.NewEndlessClosedRange(
					value.SmallInt(2),
				),
				value.SmallInt(3),
			),
			want: value.SmallInt(2),
		},
		"second iteration": {
			i: value.NewEndlessClosedRangeIteratorWithCurrentElement(
				value.NewEndlessClosedRange(
					value.SmallInt(2),
				),
				value.SmallInt(3),
			),
			after: value.NewEndlessClosedRangeIteratorWithCurrentElement(
				value.NewEndlessClosedRange(
					value.SmallInt(2),
				),
				value.SmallInt(4),
			),
			want: value.SmallInt(3),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			opts := comparer.Options()
			got, err := vm.EndlessClosedRangeIteratorNext(v, tc.i)
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
