package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestArrayTupleContains(t *testing.T) {
	tests := map[string]struct {
		tuple    *value.ArrayTuple
		val      value.Value
		contains bool
		err      value.Value
	}{
		"empty tuple": {
			tuple:    &value.ArrayTuple{},
			val:      value.SmallInt(5),
			contains: false,
		},
		"coercible elements": {
			tuple:    &value.ArrayTuple{value.String("foo"), value.Float(5)},
			val:      value.SmallInt(5),
			contains: false,
		},
		"has the value": {
			tuple:    &value.ArrayTuple{value.String("foo"), value.SmallInt(5), value.Float(9.3)},
			val:      value.SmallInt(5),
			contains: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			contains, err := vm.ArrayTupleContains(v, tc.tuple, tc.val)
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

func TestArrayTupleEqual(t *testing.T) {
	tests := map[string]struct {
		tuple *value.ArrayTuple
		other *value.ArrayTuple
		equal bool
		err   value.Value
	}{
		"two identical tuples": {
			tuple: &value.ArrayTuple{value.String("foo"), value.Float(5)},
			other: &value.ArrayTuple{value.String("foo"), value.Float(5)},
			equal: true,
		},
		"different length": {
			tuple: &value.ArrayTuple{value.String("foo"), value.Float(5)},
			other: &value.ArrayTuple{value.String("foo"), value.Float(5), value.Nil},
			equal: false,
		},
		"the same values of different types": {
			tuple: &value.ArrayTuple{value.String("foo"), value.SmallInt(5)},
			other: &value.ArrayTuple{value.String("foo"), value.Float(5)},
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.ArrayTupleEqual(v, tc.tuple, tc.other)
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
