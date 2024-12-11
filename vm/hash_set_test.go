package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestHashSetEqual(t *testing.T) {
	tests := map[string]struct {
		x     *value.HashSet
		y     *value.HashSet
		equal bool
		err   value.Value
	}{
		"two empty sets should be equal": {
			x:     &value.HashSet{},
			y:     &value.HashSet{},
			equal: true,
		},
		"two sets with different number of elements": {
			x: vm.MustNewHashSetWithElements(
				nil,
				value.SmallInt(5).ToValue(),
			),
			y:     &value.HashSet{},
			equal: false,
		},
		"two equal sets": {
			x: vm.MustNewHashSetWithElements(
				nil,
				value.SmallInt(5).ToValue(),
			),
			y: vm.MustNewHashSetWithElements(
				nil,
				value.SmallInt(5).ToValue(),
			),
			equal: true,
		},
		"two sets with different values": {
			x: vm.MustNewHashSetWithElements(
				nil,
				value.SmallInt(3).ToValue(),
				value.Float(8.5).ToValue(),
			),
			y: vm.MustNewHashSetWithElements(
				nil,
				value.SmallInt(5).ToValue(),
				value.Float(8.5).ToValue(),
			),
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.HashSetEqual(v, tc.x, tc.y)
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

func TestNewHashSetWithElements(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without VM with primitives": func(t *testing.T) {
			elements := []value.Value{
				value.Ref(value.String("foo")),
				value.Ref(value.String("bar")),
			}

			set, err := vm.NewHashSetWithElements(nil, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if set.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", set.Length())
			}
			if set.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", set.Capacity())
			}
		},
		"without VM with complex types": func(t *testing.T) {
			elements := []value.Value{
				value.SmallInt(5).ToValue(),
				value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")),
			}

			set, err := vm.NewHashSetWithElements(nil, elements...)
			if !err.IsNil() {
				t.Fatalf("error is not value.Nil: %#v", err)
			}
			if set != nil {
				t.Fatalf("result should be nil, got: %#v", set)
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			elements := []value.Value{
				value.SmallInt(5).ToValue(),
				value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")),
			}

			set, err := vm.NewHashSetWithElements(vm.New(), elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if set.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", set.Length())
			}
			if set.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", set.Capacity())
			}
		},
		"with VM with complex types that implement hash": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			elements := []value.Value{
				value.SmallInt(5).ToValue(),
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
			}

			set, err := vm.NewHashSetWithElements(vm.New(), elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if set.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", set.Length())
			}
			if set.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", set.Capacity())
			}
		},
		"with VM with complex types that implements hash improperly": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.SmallInt(5).ToValue(), value.Undefined
			})

			elements := []value.Value{
				value.SmallInt(5).ToValue(),
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
			}
			wantErr := value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			))

			set, err := vm.NewHashSetWithElements(vm.New(), elements...)
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
			if set != nil {
				t.Fatalf("result should be nil, got: %#v", set)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}

func TestNewHashSetWithCapacityAndElements(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without VM with primitives and capacity equal to length": func(t *testing.T) {
			elements := []value.Value{
				value.SmallInt(5).ToValue(),
				value.Float(25.4).ToValue(),
			}

			set, err := vm.NewHashSetWithCapacityAndElements(nil, 2, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if set.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", set.Length())
			}
			if set.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", set.Capacity())
			}
		},
		"without VM with primitives and capacity greater than length": func(t *testing.T) {
			elements := []value.Value{
				value.SmallInt(5).ToValue(),
				value.Float(25.4).ToValue(),
			}

			set, err := vm.NewHashSetWithCapacityAndElements(nil, 10, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if set.Length() != 2 {
				t.Fatalf("result should be 2, got: %d", set.Length())
			}
			if set.Capacity() != 10 {
				t.Fatalf("result should be 10, got: %d", set.Capacity())
			}
		},
		"without VM with complex types": func(t *testing.T) {
			elements := []value.Value{
				value.SmallInt(5).ToValue(),
				value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")),
			}

			set, err := vm.NewHashSetWithCapacityAndElements(nil, 2, elements...)
			if !err.IsNil() {
				t.Fatalf("error is not value.Nil: %#v", err)
			}
			if set != nil {
				t.Fatalf("result should be nil, got: %#v", set)
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			elements := []value.Value{
				value.SmallInt(5).ToValue(),
				value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")),
			}

			set, err := vm.NewHashSetWithCapacityAndElements(vm.New(), 2, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if set.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", set.Length())
			}
			if set.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", set.Length())
			}
		},
		"with VM with complex types that implement hash and capacity equal to length": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			elements := []value.Value{
				value.SmallInt(5).ToValue(),
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
			}

			set, err := vm.NewHashSetWithCapacityAndElements(vm.New(), 2, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if set.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", set.Length())
			}
			if set.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", set.Length())
			}
		},
		"with VM with complex types that implement hash and capacity greater than length": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			elements := []value.Value{
				value.SmallInt(5).ToValue(),
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
			}

			set, err := vm.NewHashSetWithCapacityAndElements(vm.New(), 6, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if set.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", set.Length())
			}
			if set.Capacity() != 6 {
				t.Fatalf("capacity should be 6, got: %d", set.Capacity())
			}
		},
		"with VM with complex types that implement hash improperly": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.SmallInt(5).ToValue(), value.Undefined
			})

			elements := []value.Value{
				value.SmallInt(5).ToValue(),
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
			}
			wantErr := value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			))

			set, err := vm.NewHashSetWithCapacityAndElements(vm.New(), 2, elements...)
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
			if set != nil {
				t.Fatalf("result should be nil, got: %#v", set)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}

func TestHashSetContains(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without vm get from empty hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithElements(nil)

			result, err := vm.HashSetContains(nil, set, value.Ref(value.String("foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"without vm get missing key from full hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(nil, set, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"without vm get missing key from hashset with deleted elements": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			vm.HashSetDelete(nil, set, value.ToSymbol("foo").ToValue())

			result, err := vm.HashSetContains(nil, set, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"without vm get missing key from hashset with left capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				10,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(nil, set, value.Ref(value.String("bar")))
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"without vm get key from full hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(nil, set, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"without vm get key from full hashset 2": func(t *testing.T) {
			set := vm.MustNewHashSetWithElements(
				nil,
				value.Ref(value.String("baz")),
				value.SmallInt(1).ToValue(),
				value.Ref(value.String("foo")),
			)

			result, err := vm.HashSetContains(nil, set, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"without vm get key from hashset with deleted elements": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			vm.HashSetDelete(nil, set, value.ToSymbol("foo").ToValue())

			result, err := vm.HashSetContains(nil, set, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"without vm get key from hashset with left capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(nil, set, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"without vm get key that is a complex type": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(nil, set, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != value.Nil {
				t.Fatalf("error should be value.Nil, got: %#v", err)
			}
		},
		"with vm get from empty hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithElements(nil)

			result, err := vm.HashSetContains(vm.New(), set, value.Ref(value.String("foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"with vm get missing key from full hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(vm.New(), set, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"with vm get missing key from hashset with left capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				10,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(vm.New(), set, value.Ref(value.String("bar")))
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"with vm get key from full hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(vm.New(), set, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"with vm get key from hashset with left capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(vm.New(), set, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"with vm get key that does not implement hash": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(vm.New(), set, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"with vm get missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(vm.New(), set, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"with vm get key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True, value.Undefined
				}
				return value.False, value.Undefined
			}, vm.DefWithParameters(1))

			v := vm.New()
			set := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetContains(v, set, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}

func TestHashSetSetCapacity(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without VM with primitives and reduce capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				10,
				value.Float(25.4).ToValue(),
				value.SmallInt(5).ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Float(25.4).ToValue(),
				value.SmallInt(5).ToValue(),
			)

			err := vm.HashSetSetCapacity(nil, set, 2)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without VM with primitives and set capacity to the same value": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				10,
				value.Float(25.4).ToValue(),
				value.SmallInt(5).ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				10,
				value.Float(25.4).ToValue(),
				value.SmallInt(5).ToValue(),
			)

			err := vm.HashSetSetCapacity(nil, set, 10)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without VM with primitives and expand capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				10,
				value.Float(25.4).ToValue(),
				value.SmallInt(5).ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				25,
				value.Float(25.4).ToValue(),
				value.SmallInt(5).ToValue(),
			)

			err := vm.HashSetSetCapacity(nil, set, 25)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without VM with complex types": func(t *testing.T) {
			set := &value.HashSet{
				Table: []value.Value{
					value.SmallInt(5).ToValue(),
					value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")),
				},
				OccupiedSlots: 2,
			}

			err := vm.HashSetSetCapacity(nil, set, 25)
			if !err.IsNil() {
				t.Fatalf("error is not nil: %s", err.Inspect())
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			key := value.NewError(value.ArgumentErrorClass, "foo bar")
			set := &value.HashSet{
				Table: []value.Value{
					value.SmallInt(5).ToValue(),
					value.Ref(key),
				},
				OccupiedSlots: 2,
			}

			v := vm.New()
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				25,
				value.SmallInt(5).ToValue(),
				value.Ref(key),
			)

			err := vm.HashSetSetCapacity(v, set, 25)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %s", err.Inspect())
			}
			if !cmp.Equal(expected, set, comparer.Options()) {
				t.Fatalf("expected: %s, set: %s\n", expected.Inspect(), set.Inspect())
			}
		},
		"with VM with complex types that implement hash": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(
				&testClass.MethodContainer,
				"hash",
				func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
					return value.UInt64(10).ToValue(), value.Undefined
				},
			)
			vm.Def(
				&testClass.MethodContainer,
				"==",
				func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
					if _, ok := args[1].MustReference().(*value.Object); ok {
						return value.True, value.Undefined
					}
					return value.False, value.Undefined
				},
				vm.DefWithParameters(1),
			)

			v := vm.New()
			set := vm.MustNewHashSetWithCapacityAndElements(
				v,
				5,
				value.SmallInt(5).ToValue(),
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				10,
				value.SmallInt(5).ToValue(),
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
			)

			err := vm.HashSetSetCapacity(v, set, 10)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %s", err.Inspect())
			}
			if !cmp.Equal(expected, set, comparer.Options()) {
				t.Fatalf("expected: %s, set: %s\n", expected.Inspect(), set.Inspect())
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}

func TestHashSetAdd(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without vm set in empty hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithElements(nil)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				5,
				value.Ref(value.String("foo")),
			)

			err := vm.HashSetAppend(nil, set, value.Ref(value.String("foo")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set existing key in full hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				4,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			err := vm.HashSetAppend(nil, set, value.Ref(value.String("foo")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set existing key in hashset with left capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				10,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				10,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			err := vm.HashSetAppend(nil, set, value.Ref(value.String("foo")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set key in full hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				4,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
				value.Ref(value.String("bar")),
			)

			err := vm.HashSetAppend(nil, set, value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set key in hashset with left capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.Ref(value.String("bar")),
				value.ToSymbol("foo").ToValue(),
			)

			err := vm.HashSetAppend(nil, set, value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set key that is a complex type": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			err := vm.HashSetAppend(nil, set, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if !err.IsNil() {
				t.Fatalf("error should be value.Nil, got: %s", err.Inspect())
			}
		},
		"with vm set in empty hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithElements(nil)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				5,
				value.Ref(value.String("foo")),
			)

			err := vm.HashSetAppend(vm.New(), set, value.Ref(value.String("foo")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set existing key in full hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				4,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			err := vm.HashSetAppend(vm.New(), set, value.Ref(value.String("foo")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set existing key in hashset with left capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				10,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				10,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			err := vm.HashSetAppend(vm.New(), set, value.Ref(value.String("foo")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set key in full hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				4,
				value.ToSymbol("foo").ToValue(),
				value.Ref(value.String("foo")),
				value.Ref(value.String("bar")),
			)

			err := vm.HashSetAppend(vm.New(), set, value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set key in hashset with left capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.Ref(value.String("bar")),
				value.ToSymbol("foo").ToValue(),
			)

			err := vm.HashSetAppend(vm.New(), set, value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set key that does not implement hash": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			key := value.NewError(value.ArgumentErrorClass, "foo")
			v := vm.New()
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.Ref(key),
				value.ToSymbol("foo").ToValue(),
				value.Ref(value.String("foo")),
			)

			err := vm.HashSetAppend(vm.New(), set, value.Ref(key))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set existing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True, value.Undefined
				}
				return value.False, value.Undefined
			}, vm.DefWithParameters(1))

			v := vm.New()
			set := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.Ref(value.String("foo")),
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
				value.ToSymbol("foo").ToValue(),
				value.Ref(value.String("foo")),
			)

			err := vm.HashSetAppend(v, set, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True, value.Undefined
				}
				return value.False, value.Undefined
			}, vm.DefWithParameters(1))

			v := vm.New()
			set := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.ToSymbol("foo").ToValue(),
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
				value.Ref(value.String("foo")),
			)

			err := vm.HashSetAppend(v, set, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}

func TestHashSetDelete(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without vm delete from empty hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithElements(nil)
			expected := vm.MustNewHashSetWithElements(nil)

			result, err := vm.HashSetDelete(nil, set, value.Ref(value.String("foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm delete key from full hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(nil, set, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm delete key from hashset with left capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				6,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				6,
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(nil, set, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm delete missing key from full hashset": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(nil, set, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm delete missing key from hashset with left capacity": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(nil, set, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm delete key that is a complex type": func(t *testing.T) {
			set := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				nil,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(nil, set, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsNil() {
				t.Fatalf("error should be value.Nil, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm deletes from empty hashset": func(t *testing.T) {
			v := vm.New()
			set := vm.MustNewHashSetWithElements(v)

			result, err := vm.HashSetDelete(v, set, value.Ref(value.String("foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
		},
		"with vm delete missing key from full hashset": func(t *testing.T) {
			v := vm.New()
			set := vm.MustNewHashSetWithCapacityAndElements(
				v,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(v, set, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete missing key from hashset with left capacity": func(t *testing.T) {
			v := vm.New()
			set := vm.MustNewHashSetWithCapacityAndElements(
				v,
				10,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				10,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(v, set, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete key from full hashset": func(t *testing.T) {
			v := vm.New()
			set := vm.MustNewHashSetWithCapacityAndElements(
				v,
				2,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				2,
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(v, set, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete key from hashset with left capacity": func(t *testing.T) {
			v := vm.New()
			set := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(v, set, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete key that does not implement hash": func(t *testing.T) {
			v := vm.New()
			set := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(v, set, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.Ref(value.UInt64(5)), value.Undefined
			})

			v := vm.New()
			set := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.Ref(value.String("foo")),
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(v, set, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.Ref(value.UInt64(5)), value.Undefined
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True, value.Undefined
				}
				return value.False, value.Undefined
			}, vm.DefWithParameters(1))

			v := vm.New()
			set := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
				value.ToSymbol("foo").ToValue(),
			)
			expected := vm.MustNewHashSetWithCapacityAndElements(
				v,
				8,
				value.ToSymbol("foo").ToValue(),
			)

			result, err := vm.HashSetDelete(v, set, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if result != true {
				t.Fail()
				t.Logf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fail()
				t.Logf("error should be undefined, got: %s", err.Inspect())
			}
			if diff := cmp.Diff(expected, set, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}
