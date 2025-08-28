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
					Key:   value.Ref(value.String("foo")),
					Value: value.Ref(value.String("bar")),
				},
				value.Pair{
					Key:   value.Ref(value.String("poznan")),
					Value: value.Ref(value.String("warszawa")),
				},
			),
			val:      value.Ref(value.String("warszawa")),
			contains: true,
		},
		"contains a duplicated value": {
			h: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Ref(value.String("bar")),
				},
				value.Pair{
					Key:   value.Ref(value.String("poznan")),
					Value: value.Ref(value.String("warszawa")),
				},
				value.Pair{
					Key:   value.Ref(value.String("lodz")),
					Value: value.Ref(value.String("warszawa")),
				},
			),
			val:      value.Ref(value.String("warszawa")),
			contains: true,
		},
		"does not contain a key": {
			h: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Ref(value.String("bar")),
				},
				value.Pair{
					Key:   value.Ref(value.String("poznan")),
					Value: value.Ref(value.String("warszawa")),
				},
			),
			val:      value.Ref(value.String("poznan")),
			contains: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			contains, err := vm.HashMapContainsValue(v, tc.h, tc.val)
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
				value.Pair{Key: value.Ref(value.String("foo")), Value: value.SmallInt(5).ToValue()},
			),
			y:     &value.HashMap{},
			equal: false,
		},
		"two equal maps": {
			x: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.Ref(value.String("foo")), Value: value.SmallInt(5).ToValue()},
			),
			y: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.Ref(value.String("foo")), Value: value.SmallInt(5).ToValue()},
			),
			equal: true,
		},
		"two maps with same keys but different values": {
			x: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.Ref(value.String("foo")), Value: value.SmallInt(3).ToValue()},
				value.Pair{Key: value.Ref(value.String("bar")), Value: value.Float(8.5).ToValue()},
			),
			y: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.Ref(value.String("foo")), Value: value.SmallInt(5).ToValue()},
				value.Pair{Key: value.Ref(value.String("bar")), Value: value.Float(8.5).ToValue()},
			),
			equal: false,
		},
		"two maps with different keys": {
			x: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.Ref(value.String("baz")), Value: value.SmallInt(3).ToValue()},
				value.Pair{Key: value.Ref(value.String("bar")), Value: value.Float(8.5).ToValue()},
			),
			y: vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.Ref(value.String("foo")), Value: value.SmallInt(5).ToValue()},
				value.Pair{Key: value.Ref(value.String("bar")), Value: value.Float(8.5).ToValue()},
			),
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.HashMapEqual(v, tc.x, tc.y)
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

func TestNewHashMapWithElements(t *testing.T) {
	tests := map[string]func(*testing.T){
		"without VM with primitives": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				{Key: value.Float(25.4).ToValue(), Value: value.Ref(value.String("bar"))},
			}

			hmap, err := vm.NewHashMapWithElements(nil, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
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
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				{Key: value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")), Value: value.Ref(value.String("bar"))},
			}

			hmap, err := vm.NewHashMapWithElements(nil, elements...)
			if err.IsUndefined() {
				t.Fatalf("error is undefined")
			}
			if hmap != nil {
				t.Fatalf("result should be nil, got: %#v", hmap)
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				{Key: value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")), Value: value.Ref(value.String("bar"))},
			}

			hmap, err := vm.NewHashMapWithElements(vm.New(), elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if hmap.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", hmap.Length())
			}
			if hmap.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", hmap.Capacity())
			}
		},
		"with VM with complex types that implements hash": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			elements := []value.Pair{
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				{Key: value.Ref(value.NewObject(value.ObjectWithClass(testClass))), Value: value.Ref(value.String("bar"))},
			}

			hmap, err := vm.NewHashMapWithElements(vm.New(), elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
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
				return value.SmallInt(5).ToValue(), value.Undefined
			})

			elements := []value.Pair{
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				{Key: value.Ref(value.NewObject(value.ObjectWithClass(testClass))), Value: value.Ref(value.String("bar"))},
			}
			wantErr := value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			))

			hmap, err := vm.NewHashMapWithElements(vm.New(), elements...)
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
				t.Fatal(diff)
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
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				{Key: value.Float(25.4).ToValue(), Value: value.Ref(value.String("bar"))},
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(nil, 2, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
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
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				{Key: value.Float(25.4).ToValue(), Value: value.Ref(value.String("bar"))},
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(nil, 10, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if hmap.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", hmap.Length())
			}
			if hmap.Capacity() != 10 {
				t.Fatalf("capacity should be 10, got: %d", hmap.Capacity())
			}
		},
		"without VM with complex types": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				{Key: value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")), Value: value.Ref(value.String("bar"))},
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(nil, 2, elements...)
			if !err.IsNil() {
				t.Fatalf("error is not value.Nil: %#v", err)
			}
			if hmap != nil {
				t.Fatalf("result should be nil, got: %#v", hmap)
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			elements := []value.Pair{
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				{Key: value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")), Value: value.Ref(value.String("bar"))},
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(vm.New(), 2, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if hmap.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", hmap.Length())
			}
			if hmap.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", hmap.Length())
			}
		},
		"with VM with complex types that implement hash and capacity equal to length": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			elements := []value.Pair{
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				{Key: value.Ref(value.NewObject(value.ObjectWithClass(testClass))), Value: value.Ref(value.String("bar"))},
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(vm.New(), 2, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if hmap.Length() != 2 {
				t.Fatalf("length should be 2, got: %d", hmap.Length())
			}
			if hmap.Capacity() != 2 {
				t.Fatalf("capacity should be 2, got: %d", hmap.Capacity())
			}
		},
		"with VM with complex types that implement hash and capacity greater than length": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			elements := []value.Pair{
				{Key: value.Ref(value.NewObject(value.ObjectWithClass(testClass))), Value: value.Ref(value.String("bar"))},
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
			}

			hmap, err := vm.NewHashMapWithCapacityAndElements(vm.New(), 6, elements...)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
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
				return value.SmallInt(5).ToValue(), value.Undefined
			})

			elements := []value.Pair{
				{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				{Key: value.Ref(value.NewObject(value.ObjectWithClass(testClass))), Value: value.Ref(value.String("bar"))},
			}
			wantErr := value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			))

			hmap, err := vm.NewHashMapWithCapacityAndElements(vm.New(), 2, elements...)
			if diff := cmp.Diff(wantErr, err, comparer.Options()); diff != "" {
				t.Fatal(diff)
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

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{
				Key:   value.Ref(value.String("foo")),
				Value: value.Ref(value.String("bar")),
			})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Ref(value.Float(2.6).ToValue().AsInlineTimeSpan()),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.Ref(value.String("bar")), Value: value.True})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.Ref(value.String("bar")), Value: value.False})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.Ref(value.String("bar")), Value: value.Ref(value.String("barina"))})
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.Ref(value.String("foo")), Value: value.Float(2.6).ToValue()})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key with wrong value from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.Ref(value.String("foo")), Value: value.Float(35).ToValue()})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap 2": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.Ref(value.String("baz")), Value: value.SmallInt(9).ToValue()},
				value.Pair{Key: value.SmallInt(1).ToValue(), Value: value.Float(2.5).ToValue()},
				value.Pair{Key: value.Ref(value.String("foo")), Value: value.Int64(3).ToValue()},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.Ref(value.String("foo")), Value: value.Int64(3).ToValue()})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.Ref(value.String("foo")), Value: value.Float(2.6).ToValue()})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.Ref(value.String("foo")), Value: value.Float(2.6).ToValue()})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(nil, hmap, &value.Pair{Key: value.Ref(value.NewError(value.ArgumentErrorClass, "foo")), Value: value.True})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != value.Nil {
				t.Fatalf("error should be value.Nil, got: %#v", err)
			}
		},
		"with vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.Ref(value.String("foo")), Value: value.Nil})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.Ref(value.String("bar")), Value: value.ToSymbol("bum").ToValue()})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.Ref(value.String("bar")), Value: value.False})
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.Ref(value.String("foo")), Value: value.Float(2.6).ToValue()})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.Ref(value.String("foo")), Value: value.Float(2.6).ToValue()})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.Ref(value.NewError(value.ArgumentErrorClass, "foo")), Value: value.Int64(3).ToValue()})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(vm.New(), hmap, &value.Pair{Key: value.Ref(value.NewObject(value.ObjectWithClass(testClass))), Value: value.Nil})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
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
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(v, hmap, &value.Pair{Key: value.Ref(value.NewObject(value.ObjectWithClass(testClass))), Value: value.Float(2.6).ToValue()})
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that implements necessary methods but has wrong value": func(t *testing.T) {
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
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContains(v, hmap, &value.Pair{Key: value.Ref(value.NewObject(value.ObjectWithClass(testClass))), Value: value.Float(24).ToValue()})
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
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

			result, err := vm.HashMapContainsKey(nil, hmap, value.Ref(value.String("foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapContainsKey(nil, hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap 2": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.Ref(value.String("baz")), Value: value.SmallInt(9).ToValue()},
				value.Pair{Key: value.SmallInt(1).ToValue(), Value: value.Float(2.5).ToValue()},
				value.Pair{Key: value.Ref(value.String("foo")), Value: value.Int64(3).ToValue()},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapContainsKey(nil, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(nil, hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsNil() {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.Ref(value.String("foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(vm.New(), hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
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
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapContainsKey(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
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

			result, err := vm.HashMapGet(nil, hmap, value.Ref(value.String("foo")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(nil, hmap, value.Ref(value.String("bar")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapGet(nil, hmap, value.Ref(value.String("bar")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(nil, hmap, value.Ref(value.String("bar")))
			if !result.IsUndefined() {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(nil, hmap, value.Ref(value.String("foo")))
			if !result.IsFloat() || result.AsFloat() != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap 2": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(
				nil,
				value.Pair{Key: value.Ref(value.String("baz")), Value: value.SmallInt(9).ToValue()},
				value.Pair{Key: value.SmallInt(1).ToValue(), Value: value.Float(2.5).ToValue()},
				value.Pair{Key: value.Ref(value.String("foo")), Value: value.Int64(3).ToValue()},
			)

			result, err := vm.HashMapGet(nil, hmap, value.Ref(value.String("foo")))
			if !result.IsInlineInt64() || result.AsInlineInt64() != value.Int64(3) {
				t.Fatalf("result should be 3i64, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			vm.HashMapDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapGet(nil, hmap, value.Ref(value.String("foo")))
			if !result.IsFloat() || result.AsFloat() != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(nil, hmap, value.Ref(value.String("foo")))
			if !result.IsFloat() || result.AsFloat() != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(nil, hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsNil() {
				t.Fatalf("error should be value.Nil, got: %#v", err)
			}
		},
		"with vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)

			result, err := vm.HashMapGet(vm.New(), hmap, value.Ref(value.String("foo")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.Ref(value.String("bar")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.Ref(value.String("bar")))
			if !result.IsUndefined() {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.Ref(value.String("foo")))
			if !result.IsFloat() || result.AsFloat() != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.Ref(value.String("foo")))
			if !result.IsFloat() || result.AsFloat() != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(vm.New(), hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that implements necessary methods": func(t *testing.T) {
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
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapGet(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if !result.IsFloat() || result.AsFloat() != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
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
				value.Pair{Key: value.Float(25.4).ToValue(), Value: value.Ref(value.String("bar"))},
				value.Pair{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{Key: value.Float(25.4).ToValue(), Value: value.Ref(value.String("bar"))},
				value.Pair{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
			)

			err := vm.HashMapSetCapacity(nil, hmap, 2)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without VM with primitives and set capacity to the same value": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{Key: value.Float(25.4).ToValue(), Value: value.Ref(value.String("bar"))},
				value.Pair{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{Key: value.Float(25.4).ToValue(), Value: value.Ref(value.String("bar"))},
				value.Pair{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
			)

			err := vm.HashMapSetCapacity(nil, hmap, 10)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without VM with primitives and expand capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{Key: value.Float(25.4).ToValue(), Value: value.Ref(value.String("bar"))},
				value.Pair{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				25,
				value.Pair{Key: value.Float(25.4).ToValue(), Value: value.Ref(value.String("bar"))},
				value.Pair{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
			)

			err := vm.HashMapSetCapacity(nil, hmap, 25)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without VM with complex types": func(t *testing.T) {
			hmap := &value.HashMap{
				Table: []value.Pair{
					{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
					{Key: value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")), Value: value.Ref(value.String("bar"))},
				},
				OccupiedSlots: 2,
			}

			err := vm.HashMapSetCapacity(nil, hmap, 25)
			if !err.IsNil() {
				t.Fatalf("error is not nil: %#v", err)
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			key := value.NewError(value.ArgumentErrorClass, "foo bar")
			hmap := &value.HashMap{
				Table: []value.Pair{
					{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
					{Key: value.Ref(key), Value: value.Ref(value.String("bar"))},
				},
				OccupiedSlots: 2,
			}
			v := vm.New()
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				25,
				value.Pair{Key: value.Ref(key), Value: value.Ref(value.String("bar"))},
				value.Pair{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
			)

			err := vm.HashMapSetCapacity(vm.New(), hmap, 25)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if !cmp.Equal(expected, hmap, comparer.Options()) {
				t.Fatalf("expected: %s, hmap: %s\n", expected.Inspect(), hmap.Inspect())
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
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				5,
				value.Pair{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				value.Pair{Key: value.Ref(value.NewObject(value.ObjectWithClass(testClass))), Value: value.Ref(value.String("bar"))},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				10,
				value.Pair{Key: value.SmallInt(5).ToValue(), Value: value.Ref(value.String("foo"))},
				value.Pair{Key: value.Ref(value.NewObject(value.ObjectWithClass(testClass))), Value: value.Ref(value.String("bar"))},
			)

			err := vm.HashMapSetCapacity(v, hmap, 10)
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
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
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(5.9).ToValue(),
				},
			)

			err := vm.HashMapSet(nil, hmap, value.Ref(value.String("foo")), value.Float(5.9).ToValue())
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set existing key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Ref(value.String("bar")),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(nil, hmap, value.Ref(value.String("foo")), value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set existing key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Ref(value.String("bar")),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(nil, hmap, value.Ref(value.String("foo")), value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
				value.Pair{
					Key:   value.Ref(value.String("bar")),
					Value: value.Float(45.8).ToValue(),
				},
			)

			err := vm.HashMapSet(nil, hmap, value.Ref(value.String("bar")), value.Float(45.8).ToValue())
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.Ref(value.String("bar")),
					Value: value.False,
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(nil, hmap, value.Ref(value.String("bar")), value.False)
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(nil, hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")), value.True)
			if !err.IsNil() {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm set in empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithElements(nil)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				5,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(5.9).ToValue(),
				},
			)

			err := vm.HashMapSet(vm.New(), hmap, value.Ref(value.String("foo")), value.Float(5.9).ToValue())
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set existing key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Ref(value.String("bar")),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(vm.New(), hmap, value.Ref(value.String("foo")), value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set existing key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Ref(value.String("bar")),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(vm.New(), hmap, value.Ref(value.String("foo")), value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				4,
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.Ref(value.String("bar")),
					Value: value.False,
				},
			)

			err := vm.HashMapSet(vm.New(), hmap, value.Ref(value.String("bar")), value.False)
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.Ref(value.String("bar")),
					Value: value.UInt16(8).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			err := vm.HashMapSet(vm.New(), hmap, value.Ref(value.String("bar")), value.UInt16(8).ToValue())
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			v := vm.New()
			key := value.NewError(value.ArgumentErrorClass, "foo")
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(key),
					Value: value.True,
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
			)

			err := vm.HashMapSet(vm.New(), hmap, value.Ref(key), value.True)
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
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
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					Value: value.Ref(value.String("lol")),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					Value: value.Nil,
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
			)

			err := vm.HashMapSet(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Nil)
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
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
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
				value.Pair{
					Key:   value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					Value: value.Nil,
				},
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
			)

			err := vm.HashMapSet(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Nil)
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
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

			result, err := vm.HashMapDelete(nil, hmap, value.Ref(value.String("foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm delete key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.SmallInt(5).ToValue(),
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.SmallInt(5).ToValue(),
				},
			)

			result, err := vm.HashMapDelete(nil, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm delete key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				6,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.SmallInt(5).ToValue(),
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				6,
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.SmallInt(5).ToValue(),
				},
			)

			result, err := vm.HashMapDelete(nil, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm delete missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(nil, hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm delete missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(nil, hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm delete key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				nil,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(nil, hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != value.Nil {
				t.Fatalf("error should be value.Nil, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm deletes from empty hashmap": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithElements(v)

			result, err := vm.HashMapDelete(v, hmap, value.Ref(value.String("foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm delete missing key from full hashmap": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete missing key from hashmap with left capacity": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				10,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete key from full hashmap": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				2,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.SmallInt(5).ToValue(),
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				2,
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.SmallInt(5).ToValue(),
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete key from hashmap with left capacity": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.SmallInt(5).ToValue(),
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.SmallInt(5).ToValue(),
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete key that does not implement hash": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			v := vm.New()
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.String("foo")),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.True,
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm delete key that implements necessary methods": func(t *testing.T) {
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
			hmap := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					Value: value.Float(2.6).ToValue(),
				},
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.SmallInt(5).ToValue(),
				},
			)
			expected := vm.MustNewHashMapWithCapacityAndElements(
				v,
				8,
				value.Pair{
					Key:   value.ToSymbol("foo").ToValue(),
					Value: value.SmallInt(5).ToValue(),
				},
			)

			result, err := vm.HashMapDelete(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
	}

	for name, tc := range tests {
		t.Run(name, tc)
	}
}
