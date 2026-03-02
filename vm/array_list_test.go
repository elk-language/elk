package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestArrayListEqual(t *testing.T) {
	tests := map[string]struct {
		list  *value.ArrayListOfValue
		other *value.ArrayListOfValue
		equal bool
		err   value.Value
	}{
		"two identical lists": {
			list:  &value.ArrayListOfValue{value.Ref(value.String("foo")), value.Float(5).ToValue()},
			other: &value.ArrayListOfValue{value.Ref(value.String("foo")), value.Float(5).ToValue()},
			equal: true,
		},
		"different length": {
			list:  &value.ArrayListOfValue{value.Ref(value.String("foo")), value.Float(5).ToValue()},
			other: &value.ArrayListOfValue{value.Ref(value.String("foo")), value.Float(5).ToValue(), value.Nil},
			equal: false,
		},
		"the same values of different types": {
			list:  &value.ArrayListOfValue{value.Ref(value.String("foo")), value.SmallInt(5).ToValue()},
			other: &value.ArrayListOfValue{value.Ref(value.String("foo")), value.Float(5).ToValue()},
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.ArrayTupleEqual(v, tc.list, tc.other)
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
