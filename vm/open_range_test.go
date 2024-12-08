package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestOpenRangeContains(t *testing.T) {
	tests := map[string]struct {
		r        *value.OpenRange
		val      value.Value
		contains bool
		err      value.Value
	}{
		"include int in the middle": {
			r: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(5).ToValue(),
			contains: true,
		},
		"not include int equal to start": {
			r: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(3).ToValue(),
			contains: false,
		},
		"not include int equal to end": {
			r: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(10).ToValue(),
			contains: false,
		},
		"not include int lesser than start": {
			r: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(2).ToValue(),
			contains: false,
		},
		"not include int greater than end": {
			r: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.SmallInt(11).ToValue(),
			contains: false,
		},
		"include float": {
			r: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			val:      value.Float(5.7).ToValue(),
			contains: true,
		},
		"throw when incomparable value": {
			r: value.NewOpenRange(
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
			contains, err := vm.OpenRangeContains(v, tc.r, tc.val)
			if diff := cmp.Diff(tc.err, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
			if !err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.contains, contains, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestOpenRangeEqual(t *testing.T) {
	tests := map[string]struct {
		r     *value.OpenRange
		other *value.OpenRange
		equal bool
		err   value.Value
	}{
		"two identical ranges": {
			r: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			other: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			equal: true,
		},
		"different end": {
			r: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(11).ToValue(),
			),
			other: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			equal: false,
		},
		"different start": {
			r: value.NewOpenRange(
				value.SmallInt(4).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			other: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			equal: false,
		},
		"different start and end": {
			r: value.NewOpenRange(
				value.SmallInt(4).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			other: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(15).ToValue(),
			),
			equal: false,
		},
		"Two ranges with the same values of different types": {
			r: value.NewOpenRange(
				value.SmallInt(3).ToValue(),
				value.SmallInt(10).ToValue(),
			),
			other: value.NewOpenRange(
				value.Float(3).ToValue(),
				value.Float(10).ToValue(),
			),
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.OpenRangeEqual(v, tc.r, tc.other)
			if diff := cmp.Diff(tc.err, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
			if !err.IsUndefined() {
				return
			}
			if diff := cmp.Diff(tc.equal, equal, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestOpenRangeIteratorNext(t *testing.T) {
	tests := map[string]struct {
		i     *value.OpenRangeIterator
		after *value.OpenRangeIterator
		want  value.Value
		err   value.Value
	}{
		"empty": {
			i: value.NewOpenRangeIterator(
				value.NewOpenRange(
					value.SmallInt(3).ToValue(),
					value.SmallInt(3).ToValue(),
				),
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
		"with two elements, first iteration": {
			i: value.NewOpenRangeIterator(
				value.NewOpenRange(
					value.SmallInt(1).ToValue(),
					value.SmallInt(4).ToValue(),
				),
			),
			after: value.NewOpenRangeIteratorWithCurrentElement(
				value.NewOpenRange(
					value.SmallInt(1).ToValue(),
					value.SmallInt(4).ToValue(),
				),
				value.SmallInt(2).ToValue(),
			),
			want: value.SmallInt(2).ToValue(),
		},
		"with two elements, second iteration": {
			i: value.NewOpenRangeIteratorWithCurrentElement(
				value.NewOpenRange(
					value.SmallInt(1).ToValue(),
					value.SmallInt(4).ToValue(),
				),
				value.SmallInt(2).ToValue(),
			),
			after: value.NewOpenRangeIteratorWithCurrentElement(
				value.NewOpenRange(
					value.SmallInt(1).ToValue(),
					value.SmallInt(4).ToValue(),
				),
				value.SmallInt(3).ToValue(),
			),
			want: value.SmallInt(3).ToValue(),
		},
		"with two elements, third iteration": {
			i: value.NewOpenRangeIteratorWithCurrentElement(
				value.NewOpenRange(
					value.SmallInt(1).ToValue(),
					value.SmallInt(4).ToValue(),
				),
				value.SmallInt(3).ToValue(),
			),
			err: value.ToSymbol("stop_iteration").ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			opts := comparer.Options()
			got, err := vm.OpenRangeIteratorNext(v, tc.i)
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatalf(diff)
			}
			if !err.IsUndefined() {
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
