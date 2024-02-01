package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestNewHashMapWithElements(t *testing.T) {
	tests := map[string]struct {
		vm       *vm.VM
		elements []value.Pair
		result   *value.HashMap
		err      value.Value
	}{
		"without VM with primitives": {
			elements: []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.Float(25.4), Value: value.String("bar")},
			},
			result: &value.HashMap{
				Table: []value.Pair{
					{Key: value.Float(25.4), Value: value.String("bar")},
					{Key: value.SmallInt(5), Value: value.String("foo")},
				},
				Count: 2,
			},
		},
		"without VM with complex types": {
			elements: []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.NewError(value.ArgumentErrorClass, "foo bar"), Value: value.String("bar")},
			},
			result: nil,
			err:    value.Nil,
		},
		"with VM with complex types that don't implement necessary methods": {
			vm: vm.New(),
			elements: []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.NewError(value.ArgumentErrorClass, "foo bar"), Value: value.String("bar")},
			},
			result: nil,
			err: value.NewError(
				value.NoMethodErrorClass,
				"method `hash` is not available to value of class `Std::ArgumentError`: Std::ArgumentError{message: \"foo bar\"}",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hmap, err := vm.NewHashMapWithElements(tc.vm, tc.elements...)
			if diff := cmp.Diff(tc.err, err, comparer.Comparer); diff != "" {
				t.Logf("result: %#v, err: %#v\n", hmap, err)
				t.Fatalf(diff)
			}
			if err != nil {
				return
			}

			if diff := cmp.Diff(tc.result, hmap, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
