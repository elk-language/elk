package vm_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestHashMapContainsValue(t *testing.T) {
	tests := map[string]struct {
		h        *value.HashMap
		val      value.Value
		contains bool
		err      value.Value
	}{
		"empty map": {
			h:        &value.HashMap{},
			val:      value.True,
			contains: false,
		},
		"contains a non-duplicated value": {
			h: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.String("bar"),
				},
				value.Pair{
					Key:   value.String("poznan"),
					Value: value.String("warszawa"),
				},
			),
			val:      value.String("warszawa"),
			contains: true,
		},
		"contains a duplicated value": {
			h: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.String("bar"),
				},
				value.Pair{
					Key:   value.String("poznan"),
					Value: value.String("warszawa"),
				},
				value.Pair{
					Key:   value.String("lodz"),
					Value: value.String("warszawa"),
				},
			),
			val:      value.String("warszawa"),
			contains: true,
		},
		"does not contain a key": {
			h: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.String("bar"),
				},
				value.Pair{
					Key:   value.String("poznan"),
					Value: value.String("warszawa"),
				},
			),
			val:      value.String("poznan"),
			contains: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			contains, err := vm.HashMapContainsValue(v, tc.h, tc.val)
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

func TestHashMapEqual(t *testing.T) {
	tests := map[string]struct {
		x     *value.HashMap
		y     *value.HashMap
		equal bool
		err   value.Value
	}{
		"two empty maps should be equal": {
			x:     &value.HashMap{},
			y:     &value.HashMap{},
			equal: true,
		},
		"two maps with different number of elements": {
			x: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.String("foo"), Value: value.SmallInt(5)},
			),
			y:     &value.HashMap{},
			equal: false,
		},
		"two equal maps": {
			x: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.String("foo"), Value: value.SmallInt(5)},
			),
			y: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.String("foo"), Value: value.SmallInt(5)},
			),
			equal: true,
		},
		"two maps with same keys but different values": {
			x: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.String("foo"), Value: value.SmallInt(3)},
				value.Pair{Key: value.String("bar"), Value: value.Float(8.5)},
			),
			y: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.String("foo"), Value: value.SmallInt(5)},
				value.Pair{Key: value.String("bar"), Value: value.Float(8.5)},
			),
			equal: false,
		},
		"two maps with different keys": {
			x: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.String("baz"), Value: value.SmallInt(3)},
				value.Pair{Key: value.String("bar"), Value: value.Float(8.5)},
			),
			y: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.String("foo"), Value: value.SmallInt(5)},
				value.Pair{Key: value.String("bar"), Value: value.Float(8.5)},
			),
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.HashMapEqual(v, tc.x, tc.y)
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

func TestNewHashMapWithElements(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without VM with primitives": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.Float(25.4), Value: value.String("bar")},
			}

			hmap, err := vm.NewHashMapWithElements(nil, elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if hmap.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", hmap.Length())
			}
			if hmap.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", hmap.Capacity())
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
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
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

			hmap, err := vm.NewHashMapWithElements(vm.New(), elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if hmap.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", hmap.Length())
			}
			if hmap.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", hmap.Capacity())
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
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
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

			hmap, err := vm.NewHashMapWithCapacityAndElements(nil, 2, elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if hmap.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", hmap.Length())
			}
			if hmap.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", hmap.Capacity())
			}
		},
		"without VM with primitives and capacity greater than length": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5), Value: value.String("foo")},
				{Key: value.Float(25.4), Value: value.String("bar")},
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(nil, 10, elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if hmap.Length() != 2 {
				t.Fatalf("result should be 2, got: %d", hmap.Length())
			}
			if hmap.Capacity() != 10 {
				t.Fatalf("result should be 10, got: %d", hmap.Capacity())
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
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
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

			hmap, err := vm.NewHashMapWithCapacityAndElements(vm.New(), 2, elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if hmap.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", hmap.Length())
			}
			if hmap.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", hmap.Length())
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

			hmap, err := vm.NewHashMapWithCapacityAndElements(vm.New(), 6, elements...)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if hmap.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", hmap.Length())
			}
			if hmap.Capacity() != 6 {
				t.Fatalf("capacity should be 6, got: %d", hmap.Capacity())
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
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
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

func TestHashMapContains(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.String("foo"), Value: value.String("bar")})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.String("bar"), Value: value.True})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo"))

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.String("bar"), Value: value.False})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.String("bar"), Value: value.String("barina")})
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.String("foo"), Value: value.Float(2.6)})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key with wrong value from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.String("foo"), Value: value.Float(35)})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from full hashmap 2": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.String("baz"), Value: value.SmallInt(9)},
				value.Pair{Key: value.SmallInt(1), Value: value.Float(2.5)},
				value.Pair{Key: value.String("foo"), Value: value.Int64(3)},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.String("foo"), Value: value.Int64(3)})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo"))

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.String("foo"), Value: value.Float(2.6)})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.String("foo"), Value: value.Float(2.6)})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.NewError(value.ArgumentErrorClass, "foo"), Value: value.True})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != value.Nil {
				t.Fatalf("error should be value.Nil, got: %#v", err)
			}
		},
		"with vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.String("foo"), Value: value.Nil})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.String("bar"), Value: value.ToSymbol("bum")})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.String("bar"), Value: value.False})
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.String("foo"), Value: value.Float(2.6)})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.String("foo"), Value: value.Float(2.6)})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			wantErr := value.NewError(
				value.NoMethodErrorClass,
				"method `hash` is not available to value of class `Std::ArgumentError`: Std::ArgumentError{message: \"foo\"}",
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.NewError(value.ArgumentErrorClass, "foo"), Value: value.Int64(3)})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm get missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})

			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.Nil})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True, nil
				}
				return value.False, nil
			}, vm.DefWithParameters("other"))

			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.NewObject(value.ObjectWithClass(testClass)),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(v, hmap, &value.Pair{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.Float(2.6)})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key that implements necessary methods but has wrong value": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True, nil
				}
				return value.False, nil
			}, vm.DefWithParameters("other"))

			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.NewObject(value.ObjectWithClass(testClass)),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(v, hmap, &value.Pair{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.Float(24)})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}
func TestHashMapContainsKey(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)

			result, err := vm.HashMapContainsKey(nil, hmap, value.String("foo"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.String("bar"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo"))

			result, err := vm.HashMapContainsKey(nil, hmap, value.String("bar"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.String("bar"))
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.String("foo"))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from full hashmap 2": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.String("baz"), Value: value.SmallInt(9)},
				value.Pair{Key: value.SmallInt(1), Value: value.Float(2.5)},
				value.Pair{Key: value.String("foo"), Value: value.Int64(3)},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.String("foo"))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo"))

			result, err := vm.HashMapContainsKey(nil, hmap, value.String("foo"))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.String("foo"))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.NewError(value.ArgumentErrorClass, "foo"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != value.Nil {
				t.Fatalf("error should be value.Nil, got: %#v", err)
			}
		},
		"with vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.String("foo"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.String("bar"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.String("bar"))
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.String("foo"))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.String("foo"))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			wantErr := value.NewError(
				value.NoMethodErrorClass,
				"method `hash` is not available to value of class `Std::ArgumentError`: Std::ArgumentError{message: \"foo\"}",
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.NewError(value.ArgumentErrorClass, "foo"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm get missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})

			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.NewObject(value.ObjectWithClass(testClass)))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True, nil
				}
				return value.False, nil
			}, vm.DefWithParameters("other"))

			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.NewObject(value.ObjectWithClass(testClass)),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(v, hmap, value.NewObject(value.ObjectWithClass(testClass)))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}

func TestHashMapGet(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)

			result, err := vm.HashMapGet(nil, hmap, value.String("foo"))
			if result != nil {
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(nil, hmap, value.String("bar"))
			if result != nil {
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo"))

			result, err := vm.HashMapGet(nil, hmap, value.String("bar"))
			if result != nil {
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(nil, hmap, value.String("bar"))
			if result != nil {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(nil, hmap, value.String("foo"))
			if result != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from full hashmap 2": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.String("baz"), Value: value.SmallInt(9)},
				value.Pair{Key: value.SmallInt(1), Value: value.Float(2.5)},
				value.Pair{Key: value.String("foo"), Value: value.Int64(3)},
			)

			result, err := vm.HashMapGet(nil, hmap, value.String("foo"))
			if result != value.Int64(3) {
				t.Fatalf("result should be 3i64, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo"))

			result, err := vm.HashMapGet(nil, hmap, value.String("foo"))
			if result != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(nil, hmap, value.String("foo"))
			if result != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(nil, hmap, value.NewError(value.ArgumentErrorClass, "foo"))
			if result != nil {
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if err != value.Nil {
				t.Fatalf("error should be value.Nil, got: %#v", err)
			}
		},
		"with vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)

			result, err := vm.HashMapGet(vm.New(), hmap, value.String("foo"))
			if result != nil {
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.String("bar"))
			if result != nil {
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.String("bar"))
			if result != nil {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.String("foo"))
			if result != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.String("foo"))
			if result != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			wantErr := value.NewError(
				value.NoMethodErrorClass,
				"method `hash` is not available to value of class `Std::ArgumentError`: Std::ArgumentError{message: \"foo\"}",
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.NewError(value.ArgumentErrorClass, "foo"))
			if result != nil {
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm get missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})

			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.NewObject(value.ObjectWithClass(testClass)))
			if result != nil {
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True, nil
				}
				return value.False, nil
			}, vm.DefWithParameters("other"))

			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.NewObject(value.ObjectWithClass(testClass)),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(v, hmap, value.NewObject(value.ObjectWithClass(testClass)))
			if result != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}

func TestHashMapSetCapacity(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without VM with primitives and reduce capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{Key: value.Float(25.4), Value: value.String("bar")},
				value.Pair{Key: value.SmallInt(5), Value: value.String("foo")},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{Key: value.Float(25.4), Value: value.String("bar")},
				value.Pair{Key: value.SmallInt(5), Value: value.String("foo")},
			)

			err := vm.HashMapSetCapacity(nil, hmap, 2)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without VM with primitives and set capacity to the same value": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{Key: value.Float(25.4), Value: value.String("bar")},
				value.Pair{Key: value.SmallInt(5), Value: value.String("foo")},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{Key: value.Float(25.4), Value: value.String("bar")},
				value.Pair{Key: value.SmallInt(5), Value: value.String("foo")},
			)

			err := vm.HashMapSetCapacity(nil, hmap, 10)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without VM with primitives and expand capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{Key: value.Float(25.4), Value: value.String("bar")},
				value.Pair{Key: value.SmallInt(5), Value: value.String("foo")},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				25,
				value.Pair{Key: value.Float(25.4), Value: value.String("bar")},
				value.Pair{Key: value.SmallInt(5), Value: value.String("foo")},
			)

			err := vm.HashMapSetCapacity(nil, hmap, 25)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without VM with complex types": func(t *testing.T) {
			hmap := &value.HashMap{
				Table: []value.Pair{
					{Key: value.SmallInt(5), Value: value.String("foo")},
					{Key: value.NewError(value.ArgumentErrorClass, "foo bar"), Value: value.String("bar")},
				},
				OccupiedSlots: 2,
			}

			err := vm.HashMapSetCapacity(nil, hmap, 25)
			if err != value.Nil {
				t.Fatalf("error is not value.Nil: %#v", err)
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			hmap := &value.HashMap{
				Table: []value.Pair{
					{Key: value.SmallInt(5), Value: value.String("foo")},
					{Key: value.NewError(value.ArgumentErrorClass, "foo bar"), Value: value.String("bar")},
				},
				OccupiedSlots: 2,
			}

			wantErr := value.NewError(
				value.NoMethodErrorClass,
				"method `hash` is not available to value of class `Std::ArgumentError`: Std::ArgumentError{message: \"foo bar\"}",
			)

			err := vm.HashMapSetCapacity(vm.New(), hmap, 25)
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with VM with complex types that implement hash": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(
				&testClass.MethodContainer,
				"hash",
				func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
					return value.UInt64(10), nil
				},
			)
			vm.Def(
				&testClass.MethodContainer,
				"==",
				func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
					if _, ok := args[1].(*value.Object); ok {
						return value.True, nil
					}
					return value.False, nil
				},
				vm.DefWithParameters("other"),
			)

			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				5,
				value.Pair{Key: value.SmallInt(5), Value: value.String("foo")},
				value.Pair{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				10,
				value.Pair{Key: value.SmallInt(5), Value: value.String("foo")},
				value.Pair{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
			)

			err := vm.HashMapSetCapacity(v, hmap, 10)
			if err != nil {
				t.Fatalf("error is not nil: %#v", err)
			}
			if !cmp.Equal(expected, hmap, comparer.Options()) {
				t.Fatalf("expected: %s, hmap: %s\n", expected.Inspect(), hmap.Inspect())
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}

func TestHashMapSet(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without vm set in empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(5.9),
				},
			)

			err := vm.HashMapSet(nil, hmap, value.String("foo"), value.Float(5.9))
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without vm set existing key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.String("bar"),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(nil, hmap, value.String("foo"), value.String("bar"))
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without vm set existing key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.String("bar"),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(nil, hmap, value.String("foo"), value.String("bar"))
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without vm set key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
				value.Pair{
					Key:   value.String("bar"),
					Value: value.Float(45.8),
				},
			)

			err := vm.HashMapSet(nil, hmap, value.String("bar"), value.Float(45.8))
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without vm set key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.String("bar"),
					Value: value.False,
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(nil, hmap, value.String("bar"), value.False)
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without vm set key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(nil, hmap, value.NewError(value.ArgumentErrorClass, "foo"), value.True)
			if err != value.Nil {
				t.Fatalf("error should be value.Nil, got: %#v", err)
			}
		},
		"with vm set in empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(5.9),
				},
			)

			err := vm.HashMapSet(vm.New(), hmap, value.String("foo"), value.Float(5.9))
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm set existing key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.String("bar"),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(vm.New(), hmap, value.String("foo"), value.String("bar"))
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm set existing key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.String("bar"),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(vm.New(), hmap, value.String("foo"), value.String("bar"))
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm set key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.String("bar"),
					Value: value.False,
				},
			)

			err := vm.HashMapSet(vm.New(), hmap, value.String("bar"), value.False)
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm set key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.String("bar"),
					Value: value.UInt16(8),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(vm.New(), hmap, value.String("bar"), value.UInt16(8))
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm set key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			wantErr := value.NewError(
				value.NoMethodErrorClass,
				"method `hash` is not available to value of class `Std::ArgumentError`: Std::ArgumentError{message: \"foo\"}",
			)

			err := vm.HashMapSet(vm.New(), hmap, value.NewError(value.ArgumentErrorClass, "foo"), value.True)
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm set existing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True, nil
				}
				return value.False, nil
			}, vm.DefWithParameters("other"))

			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.NewObject(value.ObjectWithClass(testClass)),
					Value: value.String("lol"),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.NewObject(value.ObjectWithClass(testClass)),
					Value: value.Nil,
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
			)

			err := vm.HashMapSet(v, hmap, value.NewObject(value.ObjectWithClass(testClass)), value.Nil)
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm set key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True, nil
				}
				return value.False, nil
			}, vm.DefWithParameters("other"))

			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
				value.Pair{
					Key:   value.NewObject(value.ObjectWithClass(testClass)),
					Value: value.Nil,
				},
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
			)

			err := vm.HashMapSet(v, hmap, value.NewObject(value.ObjectWithClass(testClass)), value.Nil)
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}

func TestHashMapDelete(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without vm delete from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)
			expected := vm.MustNewHashMapWithElements(nil)

			result, err := vm.HashMapDelete(nil, hmap, value.String("foo"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without vm delete key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.SmallInt(5),
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.SmallInt(5),
				},
			)

			result, err := vm.HashMapDelete(nil, hmap, value.String("foo"))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without vm delete key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				6,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.SmallInt(5),
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				6,
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.SmallInt(5),
				},
			)

			result, err := vm.HashMapDelete(nil, hmap, value.String("foo"))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without vm delete missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(nil, hmap, value.String("bar"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without vm delete missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(nil, hmap, value.String("bar"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"without vm delete key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(nil, hmap, value.NewError(value.ArgumentErrorClass, "foo"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != value.Nil {
				t.Fatalf("error should be value.Nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm deletes from empty hashmap": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithElements(v)

			result, err := vm.HashMapDelete(v, hmap, value.String("foo"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm delete missing key from full hashmap": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.String("bar"))
			if result != false {
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm delete missing key from hashmap with left capacity": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				10,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.String("bar"))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm delete key from full hashmap": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				2,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.SmallInt(5),
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				2,
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.SmallInt(5),
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.String("foo"))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm delete key from hashmap with left capacity": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.SmallInt(5),
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.SmallInt(5),
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.String("foo"))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm delete key that does not implement hash": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			wantErr := value.NewError(
				value.NoMethodErrorClass,
				"method `hash` is not available to value of class `Std::ArgumentError`: Std::ArgumentError{message: \"foo\"}",
			)

			result, err := vm.HashMapDelete(v, hmap, value.NewError(value.ArgumentErrorClass, "foo"))
			if result != false {
				t.Fatalf("result should be nil, got: %#v", result)
			}
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm delete missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})

			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.String("foo"),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.NewObject(value.ObjectWithClass(testClass)))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
		"with vm delete key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5), nil
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True, nil
				}
				return value.False, nil
			}, vm.DefWithParameters("other"))

			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.NewObject(value.ObjectWithClass(testClass)),
					Value: value.Float(2.6),
				},
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.SmallInt(5),
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.ToSymbol("foo"),
					Value: value.SmallInt(5),
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.NewObject(value.ObjectWithClass(testClass)))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if err != nil {
				t.Fatalf("error should be nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}
