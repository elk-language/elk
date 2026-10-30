package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Set interface {
	value.Reference
	value.NativeIterable
	All() iter.Seq[value.Value]
	Length() int
	Contains(*Thread, value.Value) (result bool, err value.Value)
}

type HashSet interface {
	Set
	UnionVal(*Thread, value.Value) (result value.Value, err value.Value)
	IntersectionVal(*Thread, value.Value) (result value.Value, err value.Value)
	Equal(*Thread, value.Value) (result bool, err value.Value)
	LaxEqual(*Thread, value.Value) (result bool, err value.Value)
	AppendVal(*Thread, value.Value) (added bool, err value.Value)
	RemoveVal(thread *Thread, other value.Value) (removed bool, err value.Value)
	IterSet() value.NativeResettableIterator
	CloneHashSet(thread *Thread, capacity int) (result HashSet, err value.Value)
	NewHashSet(capacity int) HashSet
	ToSyncHashSet() SyncHashSet
}

func hashSetConstructor(class *value.Class) value.Value {
	return NewHashSetOfValue(5).ToValue()
}

// ::Std::HashSet
func initHashSet() {
	value.HashSetClass.ConstructorFunc = hashSetConstructor

	// Instance methods
	c := value.HashSetClass
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			iterator := self.IterSet()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_hash_set",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"to_sync_hash_set",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			return self.ToSyncHashSet().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"remove",
		func(thread *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			other := args[1]
			removed, err := self.RemoveVal(thread, other)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			return value.BoolVal(removed), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"+",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			otherVal := args[1]
			return self.UnionVal(vm, otherVal)
		},
		DefWithParameters(1),
	)
	Alias(c, "union", "+")
	Alias(c, "|", "+")

	Def(
		c,
		"&",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			otherVal := args[1]
			return self.IntersectionVal(vm, otherVal)
		},
		DefWithParameters(1),
	)
	Alias(c, "intersection", "&")

	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			otherVal := args[1]
			result, err := self.Equal(vm, otherVal)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			return value.BoolVal(result), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"=~",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			otherVal := args[1]
			result, err := self.LaxEqual(vm, otherVal)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			return value.BoolVal(result), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			val := args[1]
			contains, err := self.Contains(vm, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"append",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			for element, err := range Iterate(vm, args[1]) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				_, err := self.AppendVal(vm, element)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			val := args[1]
			_, err := self.AppendVal(vm, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"push",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			val := args[1]
			added, err := self.AppendVal(vm, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.BoolVal(added), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(HashSet)
			callable := args[1]
			newSet := NewHashSetOfValue(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for val := range self.All() {
					result, err := vm.CallClosure(function, val)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					_, err = newSet.AppendVal(vm, result)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
				}
				return value.Ref(newSet), value.Undefined
			}

			// callable is another value
			for val := range self.All() {
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				_, err = newSet.AppendVal(vm, result)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
			}
			return value.Ref(newSet), value.Undefined
		},
		DefWithParameters(1),
	)
}

// ::Std::HashSet::Iterator
func initHashSetIterator() {
	// Instance methods
	c := value.HashSetIteratorClass
	Def(
		c,
		"next",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.NativeResettableIterator)
			return self.NextValue()
		},
	)
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"reset",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.NativeResettableIterator)
			self.Reset()
			return args[0], value.Undefined
		},
	)

}
