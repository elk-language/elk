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
		h        *vm.HashMapOfValue
		val      value.Value
		contains bool
		err      value.Value
	}{
		"empty map": {
			h:        &vm.HashMapOfValue{},
			val:      value.True.ToValue(),
			contains: false,
		},
		"contains a non-duplicated value": {
			h: vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Ref(value.String("bar")),
				),
				value.MakePairOfValue(
					value.Ref(value.String("poznan")),
					value.Ref(value.String("warszawa")),
				),
			),
			val:      value.Ref(value.String("warszawa")),
			contains: true,
		},
		"contains a duplicated value": {
			h: vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Ref(value.String("bar")),
				),
				value.MakePairOfValue(
					value.Ref(value.String("poznan")),
					value.Ref(value.String("warszawa")),
				),
				value.MakePairOfValue(
					value.Ref(value.String("lodz")),
					value.Ref(value.String("warszawa")),
				),
			),
			val:      value.Ref(value.String("warszawa")),
			contains: true,
		},
		"does not contain a key": {
			h: vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Ref(value.String("bar")),
				),
				value.MakePairOfValue(
					value.Ref(value.String("poznan")),
					value.Ref(value.String("warszawa")),
				),
			),
			val:      value.Ref(value.String("poznan")),
			contains: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			contains, err := vm.HashMapOfValueContainsValue(v, tc.h, tc.val)
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
		x     *vm.HashMapOfValue
		y     *vm.HashMapOfValue
		equal bool
		err   value.Value
	}{
		"two empty maps should be equal": {
			x:     &vm.HashMapOfValue{},
			y:     &vm.HashMapOfValue{},
			equal: true,
		},
		"two maps with different number of elements": {
			x: vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(value.Ref(value.String("foo")), value.SmallInt(5).ToValue()),
			),
			y:     &vm.HashMapOfValue{},
			equal: false,
		},
		"two equal maps": {
			x: vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(value.Ref(value.String("foo")), value.SmallInt(5).ToValue()),
			),
			y: vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(value.Ref(value.String("foo")), value.SmallInt(5).ToValue()),
			),
			equal: true,
		},
		"two maps with same keys but different values": {
			x: vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(value.Ref(value.String("foo")), value.SmallInt(3).ToValue()),
				value.MakePairOfValue(value.Ref(value.String("bar")), value.Float(8.5).ToValue()),
			),
			y: vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(value.Ref(value.String("foo")), value.SmallInt(5).ToValue()),
				value.MakePairOfValue(value.Ref(value.String("bar")), value.Float(8.5).ToValue()),
			),
			equal: false,
		},
		"two maps with different keys": {
			x: vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(value.Ref(value.String("baz")), value.SmallInt(3).ToValue()),
				value.MakePairOfValue(value.Ref(value.String("bar")), value.Float(8.5).ToValue()),
			),
			y: vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(value.Ref(value.String("foo")), value.SmallInt(5).ToValue()),
				value.MakePairOfValue(value.Ref(value.String("bar")), value.Float(8.5).ToValue()),
			),
			equal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v := vm.New()
			equal, err := vm.HashMapOfValueEqual(v, tc.x, tc.y)
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
			elements := []value.PairOfValue{
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Float(25.4).ToValue(), value.Ref(value.String("bar"))),
			}

			hmap, err := vm.NewHashMapOfValueWithElements(nil, elements...)
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
			elements := []value.PairOfValue{
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")), value.Ref(value.String("bar"))),
			}

			hmap, err := vm.NewHashMapOfValueWithElements(nil, elements...)
			if err.IsUndefined() {
				t.Fatalf("error is undefined")
			}
			if hmap != nil {
				t.Fatalf("result should be nil, got: %#v", hmap)
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			elements := []value.PairOfValue{
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")), value.Ref(value.String("bar"))),
			}

			hmap, err := vm.NewHashMapOfValueWithElements(vm.New(), elements...)
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
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			elements := []value.PairOfValue{
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Ref(value.String("bar"))),
			}

			hmap, err := vm.NewHashMapOfValueWithElements(vm.New(), elements...)
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
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.SmallInt(5).ToValue(), value.Undefined
			})

			elements := []value.PairOfValue{
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Ref(value.String("bar"))),
			}
			wantErr := value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			))

			hmap, err := vm.NewHashMapOfValueWithElements(vm.New(), elements...)
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
			elements := []value.PairOfValue{
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Float(25.4).ToValue(), value.Ref(value.String("bar"))),
			}

			hmap, err := vm.NewHashMapOfValueWithCapacityAndElements(nil, 2, elements...)
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
			elements := []value.PairOfValue{
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Float(25.4).ToValue(), value.Ref(value.String("bar"))),
			}

			hmap, err := vm.NewHashMapOfValueWithCapacityAndElements(nil, 10, elements...)
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
			elements := []value.PairOfValue{
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")), value.Ref(value.String("bar"))),
			}

			hmap, err := vm.NewHashMapOfValueWithCapacityAndElements(nil, 2, elements...)
			if !err.IsNil() {
				t.Fatalf("error is not value.Nil: %#v", err)
			}
			if hmap != nil {
				t.Fatalf("result should be nil, got: %#v", hmap)
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			elements := []value.PairOfValue{
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")), value.Ref(value.String("bar"))),
			}

			hmap, err := vm.NewHashMapOfValueWithCapacityAndElements(vm.New(), 2, elements...)
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
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			elements := []value.PairOfValue{
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Ref(value.String("bar"))),
			}

			hmap, err := vm.NewHashMapOfValueWithCapacityAndElements(vm.New(), 2, elements...)
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
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			elements := []value.PairOfValue{
				value.MakePairOfValue(value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Ref(value.String("bar"))),
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
			}

			hmap, err := vm.NewHashMapOfValueWithCapacityAndElements(vm.New(), 6, elements...)
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
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.SmallInt(5).ToValue(), value.Undefined
			})

			elements := []value.PairOfValue{
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Ref(value.String("bar"))),
			}
			wantErr := value.Ref(value.NewError(
				value.TypeErrorClass,
				"`Std::Int` cannot be coerced into `Std::UInt64`",
			))

			hmap, err := vm.NewHashMapOfValueWithCapacityAndElements(vm.New(), 2, elements...)
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
			hmap := vm.MustNewHashMapOfValueWithElements(nil)

			result, err := vm.HashMapOfValueContains(nil, hmap, value.NewPairOfValue(
				value.Ref(value.String("foo")),
				value.Ref(value.String("bar")),
			))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Ref(value.Float(2.6).ToValue().AsInlineTimeSpan()),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(nil, hmap, value.NewPairOfValue(value.Ref(value.String("bar")), value.True.ToValue()))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			vm.HashMapOfValueDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapOfValueContains(nil, hmap, value.NewPairOfValue(value.Ref(value.String("bar")), value.False.ToValue()))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(nil, hmap, value.NewPairOfValue(value.Ref(value.String("bar")), value.Ref(value.String("barina"))))
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(nil, hmap, value.NewPairOfValue(value.Ref(value.String("foo")), value.Float(2.6).ToValue()))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key with wrong value from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(nil, hmap, value.NewPairOfValue(value.Ref(value.String("foo")), value.Float(35).ToValue()))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap 2": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(value.Ref(value.String("baz")), value.SmallInt(9).ToValue()),
				value.MakePairOfValue(value.SmallInt(1).ToValue(), value.Float(2.5).ToValue()),
				value.MakePairOfValue(value.Ref(value.String("foo")), value.Int64(3).ToValue()),
			)

			result, err := vm.HashMapOfValueContains(nil, hmap, value.NewPairOfValue(value.Ref(value.String("foo")), value.Int64(3).ToValue()))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			vm.HashMapOfValueDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapOfValueContains(nil, hmap, value.NewPairOfValue(value.Ref(value.String("foo")), value.Float(2.6).ToValue()))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"without vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(nil, hmap, value.NewPairOfValue(value.Ref(value.String("foo")), value.Float(2.6).ToValue()))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(nil, hmap, value.NewPairOfValue(value.Ref(value.NewError(value.ArgumentErrorClass, "foo")), value.True.ToValue()))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if err != value.Nil {
				t.Fatalf("error should be value.Nil, got: %#v", err)
			}
		},
		"with vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithElements(nil)

			result, err := vm.HashMapOfValueContains(vm.New(), hmap, value.NewPairOfValue(value.Ref(value.String("foo")), value.Nil))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(vm.New(), hmap, value.NewPairOfValue(value.Ref(value.String("bar")), value.ToSymbol("bum").ToValue()))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(vm.New(), hmap, value.NewPairOfValue(value.Ref(value.String("bar")), value.False.ToValue()))
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(vm.New(), hmap, value.NewPairOfValue(value.Ref(value.String("foo")), value.Float(2.6).ToValue()))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(vm.New(), hmap, value.NewPairOfValue(value.Ref(value.String("foo")), value.Float(2.6).ToValue()))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(vm.New(), hmap, value.NewPairOfValue(value.Ref(value.NewError(value.ArgumentErrorClass, "foo")), value.Int64(3).ToValue()))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(vm.New(), hmap, value.NewPairOfValue(value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Nil))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True.ToValue(), value.Undefined
				}
				return value.False.ToValue(), value.Undefined
			}, vm.DefWithParameters(1))

			v := vm.New()
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(v, hmap, value.NewPairOfValue(value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Float(2.6).ToValue()))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that implements necessary methods but has wrong value": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True.ToValue(), value.Undefined
				}
				return value.False.ToValue(), value.Undefined
			}, vm.DefWithParameters(1))

			v := vm.New()
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContains(
				v,
				hmap,
				value.NewPairOfValue(
					value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					value.Float(24).ToValue(),
				),
			)
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
			hmap := vm.MustNewHashMapOfValueWithElements(nil)

			result, err := vm.HashMapOfValueContainsKey(nil, hmap, value.Ref(value.String("foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(nil, hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			vm.HashMapOfValueDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapOfValueContainsKey(nil, hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(nil, hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(nil, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap 2": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(value.Ref(value.String("baz")), value.SmallInt(9).ToValue()),
				value.MakePairOfValue(value.SmallInt(1).ToValue(), value.Float(2.5).ToValue()),
				value.MakePairOfValue(value.Ref(value.String("foo")), value.Int64(3).ToValue()),
			)

			result, err := vm.HashMapOfValueContainsKey(nil, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			vm.HashMapOfValueDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapOfValueContainsKey(nil, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(nil, hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(nil, hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsNil() {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithElements(nil)

			result, err := vm.HashMapOfValueContainsKey(vm.New(), hmap, value.Ref(value.String("foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(vm.New(), hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(vm.New(), hmap, value.Ref(value.String("bar")))
			if result != false {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(vm.New(), hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(vm.New(), hmap, value.Ref(value.String("foo")))
			if result != true {
				t.Fatalf("result should be true, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(vm.New(), hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(vm.New(), hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True.ToValue(), value.Undefined
				}
				return value.False.ToValue(), value.Undefined
			}, vm.DefWithParameters(1))

			v := vm.New()
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueContainsKey(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
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
			hmap := vm.MustNewHashMapOfValueWithElements(nil)

			result, err := vm.HashMapOfValueGet(nil, hmap, value.Ref(value.String("foo")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(nil, hmap, value.Ref(value.String("bar")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			vm.HashMapOfValueDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapOfValueGet(nil, hmap, value.Ref(value.String("bar")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(nil, hmap, value.Ref(value.String("bar")))
			if !result.IsUndefined() {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(nil, hmap, value.Ref(value.String("foo")))
			if !result.IsFloat() || result.AsFloat() != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from full hashmap 2": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithElements(
				nil,
				value.MakePairOfValue(
					value.Ref(value.String("baz")),
					value.SmallInt(9).ToValue(),
				),
				value.MakePairOfValue(
					value.SmallInt(1).ToValue(),
					value.Float(2.5).ToValue(),
				),
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Int64(3).ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(nil, hmap, value.Ref(value.String("foo")))
			if !result.IsInlineInt64() || result.AsInlineInt64() != value.Int64(3) {
				t.Fatalf("result should be 3i64, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from hashmap with deleted elements": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			vm.HashMapOfValueDelete(nil, hmap, value.ToSymbol("foo").ToValue())

			result, err := vm.HashMapOfValueGet(nil, hmap, value.Ref(value.String("foo")))
			if !result.IsFloat() || result.AsFloat() != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(nil, hmap, value.Ref(value.String("foo")))
			if !result.IsFloat() || result.AsFloat() != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"without vm get key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(nil, hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsNil() {
				t.Fatalf("error should be value.Nil, got: %#v", err)
			}
		},
		"with vm get from empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithElements(nil)

			result, err := vm.HashMapOfValueGet(vm.New(), hmap, value.Ref(value.String("foo")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(vm.New(), hmap, value.Ref(value.String("bar")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(vm.New(), hmap, value.Ref(value.String("bar")))
			if !result.IsUndefined() {
				t.Logf("result: %#v, err: %#v", result, err)
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(vm.New(), hmap, value.Ref(value.String("foo")))
			if !result.IsFloat() || result.AsFloat() != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key from hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(vm.New(), hmap, value.Ref(value.String("foo")))
			if !result.IsFloat() || result.AsFloat() != value.Float(2.6) {
				t.Fatalf("result should be 2.6, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(vm.New(), hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get missing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(vm.New(), hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
			if !result.IsUndefined() {
				t.Fatalf("result should be undefined, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm get key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True.ToValue(), value.Undefined
				}
				return value.False.ToValue(), value.Undefined
			}, vm.DefWithParameters(1))

			v := vm.New()
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueGet(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
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
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(value.Float(25.4).ToValue(), value.Ref(value.String("bar"))),
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(value.Float(25.4).ToValue(), value.Ref(value.String("bar"))),
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
			)

			err := vm.HashMapOfValueSetCapacity(nil, hmap, 2)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without VM with primitives and set capacity to the same value": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(value.Float(25.4).ToValue(), value.Ref(value.String("bar"))),
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(value.Float(25.4).ToValue(), value.Ref(value.String("bar"))),
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
			)

			err := vm.HashMapOfValueSetCapacity(nil, hmap, 10)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without VM with primitives and expand capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(value.Float(25.4).ToValue(), value.Ref(value.String("bar"))),
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				25,
				value.MakePairOfValue(value.Float(25.4).ToValue(), value.Ref(value.String("bar"))),
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
			)

			err := vm.HashMapOfValueSetCapacity(nil, hmap, 25)
			if !err.IsUndefined() {
				t.Fatalf("error is not undefined: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without VM with complex types": func(t *testing.T) {
			hmap := &vm.HashMapOfValue{
				Table: []value.PairOfValue{
					value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
					value.MakePairOfValue(value.Ref(value.NewError(value.ArgumentErrorClass, "foo bar")), value.Ref(value.String("bar"))),
				},
				OccupiedSlots: 2,
			}

			err := vm.HashMapOfValueSetCapacity(nil, hmap, 25)
			if !err.IsNil() {
				t.Fatalf("error is not nil: %#v", err)
			}
		},
		"with VM with complex types that don't implement necessary methods": func(t *testing.T) {
			key := value.NewError(value.ArgumentErrorClass, "foo bar")
			hmap := &vm.HashMapOfValue{
				Table: []value.PairOfValue{
					value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
					value.MakePairOfValue(value.Ref(key), value.Ref(value.String("bar"))),
				},
				OccupiedSlots: 2,
			}
			v := vm.New()
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				25,
				value.MakePairOfValue(value.Ref(key), value.Ref(value.String("bar"))),
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
			)

			err := vm.HashMapOfValueSetCapacity(vm.New(), hmap, 25)
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
				func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
					return value.UInt64(10).ToValue(), value.Undefined
				},
			)
			vm.Def(
				&testClass.MethodContainer,
				"==",
				func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
					if _, ok := args[1].MustReference().(*value.Object); ok {
						return value.True.ToValue(), value.Undefined
					}
					return value.False.ToValue(), value.Undefined
				},
				vm.DefWithParameters(1),
			)

			v := vm.New()
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				5,
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Ref(value.String("bar"))),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				10,
				value.MakePairOfValue(value.SmallInt(5).ToValue(), value.Ref(value.String("foo"))),
				value.MakePairOfValue(value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Ref(value.String("bar"))),
			)

			err := vm.HashMapOfValueSetCapacity(v, hmap, 10)
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
			hmap := vm.MustNewHashMapOfValueWithElements(nil)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				5,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(5.9).ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(nil, hmap, value.Ref(value.String("foo")), value.Float(5.9).ToValue())
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set existing key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				4,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Ref(value.String("bar")),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(nil, hmap, value.Ref(value.String("foo")), value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set existing key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Ref(value.String("bar")),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(nil, hmap, value.Ref(value.String("foo")), value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				4,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
				value.MakePairOfValue(
					value.Ref(value.String("bar")),
					value.Float(45.8).ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(nil, hmap, value.Ref(value.String("bar")), value.Float(45.8).ToValue())
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.Ref(value.String("bar")),
					value.False.ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(nil, hmap, value.Ref(value.String("bar")), value.False.ToValue())
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"without vm set key that is a complex type": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(nil, hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")), value.True.ToValue())
			if !err.IsNil() {
				t.Fatalf("error should be nil, got: %#v", err)
			}
		},
		"with vm set in empty hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithElements(nil)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				5,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(5.9).ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(vm.New(), hmap, value.Ref(value.String("foo")), value.Float(5.9).ToValue())
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set existing key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				4,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Ref(value.String("bar")),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(vm.New(), hmap, value.Ref(value.String("foo")), value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set existing key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Ref(value.String("bar")),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(vm.New(), hmap, value.Ref(value.String("foo")), value.Ref(value.String("bar")))
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set key in full hashmap": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				4,
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.Ref(value.String("bar")),
					value.False.ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(vm.New(), hmap, value.Ref(value.String("bar")), value.False.ToValue())
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set key in hashmap with left capacity": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.Ref(value.String("bar")),
					value.UInt16(8).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(vm.New(), hmap, value.Ref(value.String("bar")), value.UInt16(8).ToValue())
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set key that does not implement hash": func(t *testing.T) {
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			v := vm.New()
			key := value.NewError(value.ArgumentErrorClass, "foo")
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(key),
					value.True.ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(vm.New(), hmap, value.Ref(key), value.True.ToValue())
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set existing key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True.ToValue(), value.Undefined
				}
				return value.False.ToValue(), value.Undefined
			}, vm.DefWithParameters(1))

			v := vm.New()
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					value.Ref(value.String("lol")),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					value.Nil,
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Nil)
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
			if diff := cmp.Diff(expected, hmap, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		},
		"with vm set key that implements necessary methods": func(t *testing.T) {
			testClass := value.NewClassWithOptions(value.ClassWithName("PizdaClass"))
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True.ToValue(), value.Undefined
				}
				return value.False.ToValue(), value.Undefined
			}, vm.DefWithParameters(1))

			v := vm.New()
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
				value.MakePairOfValue(
					value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					value.Nil,
				),
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
			)

			err := vm.HashMapOfValueSet(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))), value.Nil)
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
			hmap := vm.MustNewHashMapOfValueWithElements(nil)
			expected := vm.MustNewHashMapOfValueWithElements(nil)

			result, err := vm.HashMapOfValueDelete(nil, hmap, value.Ref(value.String("foo")))
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
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.SmallInt(5).ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.SmallInt(5).ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(nil, hmap, value.Ref(value.String("foo")))
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
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				6,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.SmallInt(5).ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				6,
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.SmallInt(5).ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(nil, hmap, value.Ref(value.String("foo")))
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
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(nil, hmap, value.Ref(value.String("bar")))
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
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(nil, hmap, value.Ref(value.String("bar")))
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
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				nil,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(nil, hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
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
			hmap := vm.MustNewHashMapOfValueWithElements(v)

			result, err := vm.HashMapOfValueDelete(v, hmap, value.Ref(value.String("foo")))
			if result != false {
				t.Fatalf("result should be false, got: %#v", result)
			}
			if !err.IsUndefined() {
				t.Fatalf("error should be undefined, got: %#v", err)
			}
		},
		"with vm delete missing key from full hashmap": func(t *testing.T) {
			v := vm.New()
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(v, hmap, value.Ref(value.String("bar")))
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
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				10,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(v, hmap, value.Ref(value.String("bar")))
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
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				2,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.SmallInt(5).ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				2,
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.SmallInt(5).ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(v, hmap, value.Ref(value.String("foo")))
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
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.SmallInt(5).ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.SmallInt(5).ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(v, hmap, value.Ref(value.String("foo")))
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
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(v, hmap, value.Ref(value.NewError(value.ArgumentErrorClass, "foo")))
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
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})

			v := vm.New()
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.String("foo")),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.True.ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
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
			vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				return value.UInt64(5).ToValue(), value.Undefined
			})
			vm.Def(&testClass.MethodContainer, "==", func(vm *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
				other := args[1]
				if other.Class() == testClass {
					return value.True.ToValue(), value.Undefined
				}
				return value.False.ToValue(), value.Undefined
			}, vm.DefWithParameters(1))

			v := vm.New()
			hmap := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.Ref(value.NewObject(value.ObjectWithClass(testClass))),
					value.Float(2.6).ToValue(),
				),
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.SmallInt(5).ToValue(),
				),
			)
			expected := vm.MustNewHashMapOfValueWithCapacityAndElements(
				v,
				8,
				value.MakePairOfValue(
					value.ToSymbol("foo").ToValue(),
					value.SmallInt(5).ToValue(),
				),
			)

			result, err := vm.HashMapOfValueDelete(v, hmap, value.Ref(value.NewObject(value.ObjectWithClass(testClass))))
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
