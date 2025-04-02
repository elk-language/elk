package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestEndlessOpenRangeContains(t *testing.T) {
	tests := map[string]struct {
		r        *value.EndlessOpenRange
		val      value.Value
		contains bool
		err      value.Value
	}{
		"include int in the middle": {
			r: value.NewEndlessOpenRange(
				value.SmallInt(3).ToValue(),
			),
			val:      value.SmallInt(5).ToValue(),
			contains: true,
		},
		"not include int equal to start": {
			r: value.NewEndlessOpenRange(
				value.SmallInt(3).ToValue(),
			),
			val:      value.SmallInt(3).ToValue(),
			contains: false,
		},
		"not include int lesser than start": {
			r: value.NewEndlessOpenRange(
				value.SmallInt(3).ToValue(),
			),
			val:      value.SmallInt(2).ToValue(),
			contains: false,
		},
		"include float": {
			r: value.NewEndlessOpenRange(
				value.SmallInt(3).ToValue(),
			),
			val:      value.Float(5.7).ToValue(),
			contains: true,
		},
		"throw when incomparable value": {
			r: value.NewEndlessOpenRange(
				value.SmallInt(3).ToValue(),
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
			contains, err := vm.EndlessOpenRangeContains(v, tc.r, tc.val)
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

func TestEndlessOpenRangeEqual(t *testing.T) {
	tests := map[string]struct {
		r     *value.EndlessOpenRange
		other *value.EndlessOpenRange
		equal bool
		err   value.Value
	}{
		"two identical ranges": {
			r: value.NewEndlessOpenRange(
				value.SmallInt(3).ToValue(),
			),
			other: value.NewEndlessOpenRange(
				value.SmallInt(3).ToValue(),
			),
			equal: true,
		},
		"different start": {
			r: value.NewEndlessOpenRange(
				value.SmallInt(4).ToValue(),
			),
			other: value.NewEndlessOpenRange(
				value.SmallInt(3).ToValue(),
			),
			equal: false,
		},
		"Two ranges with the same values of different types": {
			r: value.NewEndlessOpenRange(
				value.SmallInt(3).ToValue(),
			),
			other: value.NewEndlessOpenRange(
				value.Float(3).ToValue(),
			),
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.EndlessOpenRangeEqual(v, tc.r, tc.other)
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

func TestEndlessOpenRangeIteratorNext(t *testing.T) {
	tests := map[string]struct {
		i     *value.EndlessOpenRangeIterator
		after *value.EndlessOpenRangeIterator
		want  value.Value
		err   value.Value
	}{
		"with two elements, first iteration": {
			i: value.NewEndlessOpenRangeIterator(
				value.NewEndlessOpenRange(
					value.SmallInt(1).ToValue(),
				),
			),
			after: value.NewEndlessOpenRangeIteratorWithCurrentElement(
				value.NewEndlessOpenRange(
					value.SmallInt(1).ToValue(),
				),
				value.SmallInt(2).ToValue(),
			),
			want: value.SmallInt(2).ToValue(),
		},
		"with two elements, second iteration": {
			i: value.NewEndlessOpenRangeIteratorWithCurrentElement(
				value.NewEndlessOpenRange(
					value.SmallInt(1).ToValue(),
				),
				value.SmallInt(2).ToValue(),
			),
			after: value.NewEndlessOpenRangeIteratorWithCurrentElement(
				value.NewEndlessOpenRange(
					value.SmallInt(1).ToValue(),
				),
				value.SmallInt(3).ToValue(),
			),
			want: value.SmallInt(3).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			opts := comparer.Options()
			got, err := vm.EndlessOpenRangeIteratorNext(v, tc.i)
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
