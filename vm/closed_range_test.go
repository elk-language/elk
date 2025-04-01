package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestClosedRangeContains(t *testing.T) {
	tests := map[string]struct {
		r        *value.ClosedRange
		val      value.Value
		contains bool
		err      value.Value
	}{
		"include int in the middle": {
			r: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(5).ToValue(),
			contains: true,
		},
		"include int equal to start": {
			r: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(3).ToValue(),
			contains: true,
		},
		"include int equal to end": {
			r: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(10).ToValue(),
			contains: true,
		},
		"not include int lesser than start": {
			r: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(2).ToValue(),
			contains: false,
		},
		"not include int greater than end": {
			r: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(11).ToValue(),
			contains: false,
		},
		"include float": {
			r: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.Float(5.7).ToValue(),
			contains: true,
		},
		"throw when incomparable value": {
			r: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
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
			contains, err := vm.ClosedRangeContains(v, tc.r, tc.val)
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

func TestClosedRangeEqual(t *testing.T) {
	tests := map[string]struct {
		r     *value.ClosedRange
		other *value.ClosedRange
		equal bool
		err   value.Value
	}{
		"two identical ranges": {
			r: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			other: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			equal: true,
		},
		"different end": {
			r: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(11).ToValue(),
			),
			other: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			equal: false,
		},
		"different start": {
			r: value.NewClosedRange(
				value.SmallInt(4).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			other: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			equal: false,
		},
		"different start and end": {
			r: value.NewClosedRange(
				value.SmallInt(4).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			other: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(15).ToValue(),
			),
			equal: false,
		},
		"Two ranges with the same values of different types": {
			r: value.NewClosedRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			other: value.NewClosedRange(
				value.Float(3).ToValue(),
				value.Float(10).ToValue(),
			),
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.ClosedRangeEqual(v, tc.r, tc.other)
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

func TestClosedRangeIteratorNext(t *testing.T) {
	tests := map[string]struct {
		i     *value.ClosedRangeIterator
		after *value.ClosedRangeIterator
		want  value.Value
		err   value.Value
	}{
		"empty": {
			i: value.NewClosedRangeIterator(
				value.NewClosedRange(
					value.SmallInt(3).ToValue(),
					value.SmallInt(2).ToValue(),
				),
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
		"with two elements, first iteration": {
			i: value.NewClosedRangeIterator(
				value.NewClosedRange(
					value.SmallInt(2).ToValue(),
					value.SmallInt(3).ToValue(),
				),
			),
			after: value.NewClosedRangeIteratorWithCurrentElement(
				value.NewClosedRange(
					value.SmallInt(2).ToValue(),
					value.SmallInt(3).ToValue(),
				),
				value.SmallInt(3).ToValue(),
			),
			want: value.SmallInt(2).ToValue(),
		},
		"with two elements, second iteration": {
			i: value.NewClosedRangeIteratorWithCurrentElement(
				value.NewClosedRange(
					value.SmallInt(2).ToValue(),
					value.SmallInt(3).ToValue(),
				),
				value.SmallInt(3).ToValue(),
			),
			after: value.NewClosedRangeIteratorWithCurrentElement(
				value.NewClosedRange(
					value.SmallInt(2).ToValue(),
					value.SmallInt(3).ToValue(),
				),
				value.SmallInt(4).ToValue(),
			),
			want: value.SmallInt(3).ToValue(),
		},
		"with two elements, third iteration": {
			i: value.NewClosedRangeIteratorWithCurrentElement(
				value.NewClosedRange(
					value.SmallInt(2).ToValue(),
					value.SmallInt(3).ToValue(),
				),
				value.SmallInt(4).ToValue(),
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			opts := comparer.Options()
			got, err := vm.ClosedRangeIteratorNext(v, tc.i)
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
			if !err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.after, tc.i, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
