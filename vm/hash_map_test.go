package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestNewHashMapWithElements(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without VM with primitives": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.Float(25.4), Value: value.String("bar")},
			}
			result := &value.HashMap{
				Table: []value.Pair{
					{Key: value.Float(25.4), Value: value.String("bar")},
					{Key: value.SmallInt(5), Value: value.String("foo")},
				},
				Count: 2,
			}

			hmap, err := vm.NewHashMapWithElements(nil, elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if diff := cmp.Diff(result, hmap, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without VM with complex types": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.NewError(value.ArgumentErrorClass, "foo bar"), Value: value.String("bar")},
			}

			hmap, err := vm.NewHashMapWithElements(nil, elements...)
			if err != value.Nil {
				t.Fatalf("error is not value.Nil: %#v", err)
			}
			if hmap != nil {
				t.Fatalf("result should be nil, got: %#v", hmap)
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.NewError(value.ArgumentErrorClass, "foo bar"), Value: value.String("bar")},
			}
			wantErr := value.NewError(
				value.NoMethodErrorClass,
				"method `hash` is not available to value of class `Std::ArgumentError`: Std::ArgumentError{message: \"foo bar\"}",
			)

			hmap, err := vm.NewHashMapWithElements(vm.New(), elements...)
			if diff := cmp.Diff(wantErr, err, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
			if hmap != nil {
				t.Fatalf("result should be nil, got: %#v", hmap)
			}
		},
		"with VM with complex types that implements hash": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})

			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
			}

			result := &value.HashMap{
				Table: []value.Pair{
					{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
					{Key: value.SmallInt(5), Value: value.String("foo")},
				},
				Count: 2,
			}

			hmap, err := vm.NewHashMapWithElements(vm.New(), elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if diff := cmp.Diff(result, hmap, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with VM with complex types that implements hash improperly": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.SmallInt(5), nil
			})

			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
			}
			wantErr := value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			)

			hmap, err := vm.NewHashMapWithElements(vm.New(), elements...)
			if diff := cmp.Diff(wantErr, err, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
			if hmap != nil {
				t.Fatalf("result should be nil, got: %#v", hmap)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}

func TestNewHashMapWithCapacityAndElements(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without VM with primitives and capacity equal to length": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.Float(25.4), Value: value.String("bar")},
			}
			result := &value.HashMap{
				Table: []value.Pair{
					{Key: value.Float(25.4), Value: value.String("bar")},
					{Key: value.SmallInt(5), Value: value.String("foo")},
				},
				Count: 2,
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(nil, 2, elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if diff := cmp.Diff(result, hmap, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without VM with primitives and capacity greater than length": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.Float(25.4), Value: value.String("bar")},
			}
			result := &value.HashMap{
				Table: []value.Pair{
					{},
					{},
					{},
					{},
					{},
					{Key: value.SmallInt(5), Value: value.String("foo")},
					{Key: value.Float(25.4), Value: value.String("bar")},
					{},
					{},
					{},
				},
				Count: 2,
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(nil, 10, elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if diff := cmp.Diff(result, hmap, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without VM with complex types": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.NewError(value.ArgumentErrorClass, "foo bar"), Value: value.String("bar")},
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(nil, 2, elements...)
			if err != value.Nil {
				t.Fatalf("error is not value.Nil: %#v", err)
			}
			if hmap != nil {
				t.Fatalf("result should be nil, got: %#v", hmap)
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.NewError(value.ArgumentErrorClass, "foo bar"), Value: value.String("bar")},
			}
			wantErr := value.NewError(
				value.NoMethodErrorClass,
				"method `hash` is not available to value of class `Std::ArgumentError`: Std::ArgumentError{message: \"foo bar\"}",
			)

			hmap, err := vm.NewHashMapWithCapacityAndElements(vm.New(), 2, elements...)
			if diff := cmp.Diff(wantErr, err, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
			if hmap != nil {
				t.Fatalf("result should be nil, got: %#v", hmap)
			}
		},
		"with VM with complex types that implement hash and capacity equal to length": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})

			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
			}

			result := &value.HashMap{
				Table: []value.Pair{
					{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
					{Key: value.SmallInt(5), Value: value.String("foo")},
				},
				Count: 2,
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(vm.New(), 2, elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if diff := cmp.Diff(result, hmap, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with VM with complex types that implement hash and capacity greater than length": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})

			elements := []value.Pair{
				{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
				{Key: value.SmallInt(5), Value: value.String("foo")},
			}

			result := &value.HashMap{
				Table: []value.Pair{
					{},
					{Key: value.SmallInt(5), Value: value.String("foo")},
					{},
					{},
					{},
					{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
				},
				Count: 2,
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(vm.New(), 6, elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if diff := cmp.Diff(result, hmap, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with VM with complex types that implement hash improperly": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.SmallInt(5), nil
			})

			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
			}
			wantErr := value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			)

			hmap, err := vm.NewHashMapWithCapacityAndElements(vm.New(), 2, elements...)
			if diff := cmp.Diff(wantErr, err, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
			if hmap != nil {
				t.Fatalf("result should be nil, got: %#v", hmap)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}
