package vm

import (
	"iter"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type SyncOrderedHashSet interface {
	Set
	Lock()
	Unlock()
	OrderedHashSet() OrderedHashSet
	All() iter.Seq[value.Value]
	Length() int
	// Normalizes the index so that `-1` becomes `length-1`.
	// Does bounds checking and returns (undefined, undefined) if the index is out of range.
	TryAtVal(value.Value) (value.Value, value.Value)
	// Get the element with the given index
	// Does not do bounds checking or index normalisation.
	// Panics if the index is invalid.
	AtVal(int) value.Value
	// Get the element with the given index
	// Does not do bounds checking or index normalisation.
	// Returns an error (undefined, err) if the index is out of range.
	GetAtVal(int) (value.Value, value.Value)

	UnionVal(*Thread, value.Value) (result value.Value, err value.Value)
	IntersectionVal(*Thread, value.Value) (result value.Value, err value.Value)
	Equal(*Thread, value.Value) (result bool, err value.Value)
	LaxEqual(*Thread, value.Value) (result bool, err value.Value)
	Contains(*Thread, value.Value) (result bool, err value.Value)
	AppendVal(*Thread, value.Value) (added bool, err value.Value)
	RemoveVal(thread *Thread, other value.Value) (removed bool, err value.Value)
	IterSyncOrderedSet() value.NativeResettableIterator
	CloneSyncOrderedHashSet(thread *Thread, capacity int) (result SyncOrderedHashSet, err value.Value)
	NewSyncOrderedHashSet(capacity int) SyncOrderedHashSet
}

func syncOrderedHashSetConstructor(class *value.Class) value.Value {
	return NewSyncOrderedHashSetOfValue(5).ToValue()
}

// ::Std::Sync::OrderedHashSet
func initSyncOrderedHashSet() {
	value.SyncOrderedHashSetClass.ConstructorFunc = syncOrderedHashSetConstructor

	// Instance methods
	c := value.SyncOrderedHashSetClass
	Def(
		c,
		"try_first",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashSet)
			length := self.Length()
			if length == 0 {
				return value.Nil, value.Undefined
			}

			return self.AtVal(0), value.Undefined
		},
	)
	Def(
		c,
		"first",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashSet)
			length := self.Length()
			if length == 0 {
				return value.Undefined, value.Ref(
					value.Errorf(
						value.IterableNotFoundErrorClass,
						"cannot get first element of `%s`",
						self.Inspect(),
					),
				)
			}

			return self.AtVal(0), value.Undefined
		},
	)
	Def(
		c,
		"try_last",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashSet)
			length := self.Length()
			if length == 0 {
				return value.Nil, value.Undefined
			}

			return self.AtVal(length - 1), value.Undefined
		},
	)
	Def(
		c,
		"last",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashSet)
			length := self.Length()
			if length == 0 {
				return value.Undefined, value.Ref(
					value.Errorf(
						value.IterableNotFoundErrorClass,
						"cannot get last element of `%s`",
						self.Inspect(),
					),
				)
			}

			return self.AtVal(length - 1), value.Undefined
		},
	)
	Def(
		c,
		"get_at",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashSet)
			index := args[1].AsAnyInt()
			return self.GetAtVal(index)
		},
		DefWithParameters(1),
	)
	Alias(c, "at", "get_at")

	Def(
		c,
		"try_at",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashSet)
			other := args[1]
			return self.TryAtVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashSet)
			iterator := self.IterSyncOrderedSet()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashSet)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"remove",
		func(thread *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(SyncOrderedHashSet)
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
			self := args[0].AsReference().(SyncOrderedHashSet)
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
			self := args[0].AsReference().(SyncOrderedHashSet)
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
			self := args[0].AsReference().(SyncOrderedHashSet)
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
			self := args[0].AsReference().(SyncOrderedHashSet)
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
			self := args[0].AsReference().(SyncOrderedHashSet)
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
			self := args[0].AsReference().(SyncOrderedHashSet)
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
			self := args[0].AsReference().(SyncOrderedHashSet)
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
			self := args[0].AsReference().(SyncOrderedHashSet)
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
			self := args[0].AsReference().(SyncOrderedHashSet)
			callable := args[1]
			newSet := NewSyncOrderedHashSetOfValue(self.Length())

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

// ::Std::Sync::OrderedHashSet::Iterator
func initSyncOrderedHashSetIterator() {
	// Instance methods
	c := value.SyncOrderedHashSetIteratorClass
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
