package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// ::Std::Sync::ArrayList
func initSyncArrayList() {
	// Instance methods
	c := value.SyncArrayListClass
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			iterator := self.IterSyncList()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			return value.SmallInt(self.LeftCapacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			other := args[1]
			return self.Subscript(other)
		},
		DefWithParameters(1),
	)
	Alias(c, "at", "[]")

	Def(
		c,
		"try_at",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			other := args[1]
			return self.TryAtVal(other)
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"[]=",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			key := args[1]
			val := args[2]
			err := self.SubscriptSet(key, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return val, value.Undefined
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"+",
		func(thread *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			other := args[1]
			return SyncArrayListConcat(thread, self, other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			other := args[1]
			return self.RepeatVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			switch other := args[1].SafeAsReference().(type) {
			case value.SyncArrayList:
				equal, err := TupleEqual(vm, self, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.BoolVal(equal), value.Undefined
			default:
				return value.False.ToValue(), value.Undefined
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			other := args[1]
			eq, err := TupleLaxEqualValue(vm, self, other)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			return value.Bool(eq).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"grow",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			nValue := args[1]
			n, ok := value.IntToGoInt(nValue)
			if !ok && n == -1 {
				return value.Undefined, value.Ref(value.NewTooLargeCapacityError(nValue.Inspect()))
			}
			if n < 0 {
				return value.Undefined, value.Ref(value.NewNegativeCapacityError(nValue.Inspect()))
			}
			if !ok {
				return value.Undefined, value.Ref(value.NewCapacityTypeError(nValue.Inspect()))
			}
			self.Grow(n)
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"append",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			values := args[1]

			self.Lock()
			defer self.Unlock()

			list := self.ArrayList()
			for val := range Iterate(vm, values) {
				err := list.AppendVal(val)
				if err.IsNotUndefined() {
					return value.Undefined, err
				}
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"remove",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			val := args[1]

			self.Lock()
			defer self.Unlock()

			list := self.ArrayList()
			var removed bool
			for i := 0; i < self.Length(); i++ {
				elem := list.AtVal(i)
				isEqual, err := Equal(vm, elem, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(isEqual) {
					list.RemoveAt(i)
					removed = true
				}
			}

			return value.BoolVal(removed), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"remove_at",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			val := args[1]
			return value.Nil, self.RemoveAtVal(val)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"remove_at_range",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			rng := args[1]
			self.Lock()
			defer self.Unlock()

			err := value.ArrayListInterfaceRemoveAtRange(self.ArrayList(), rng)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}
			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			err := self.AppendVal(args[1])
			if err.IsNotUndefined() {
				return value.Undefined, err
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "push", "<<")

	Def(
		c,
		"map_mut",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			callable := args[1]

			self.Lock()
			defer self.Unlock()

			list := self.ArrayList()

			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for i := range list.Length() {
					element := list.AtVal(i)
					result, err := vm.CallClosure(function, element)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					err = list.SetAtVal(i, result)
					if err.IsNotUndefined() {
						return value.Undefined, err
					}
				}
				return self.ToValue(), value.Undefined
			}

			// callable is another value
			for i := range list.Length() {
				element := list.AtVal(i)
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, element)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				err = list.SetAtVal(i, result)
				if err.IsNotUndefined() {
					return value.Undefined, err
				}
			}
			return self.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayList)
			callable := args[1]
			newList := value.NewSyncArrayListOfValueWithLength(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for i := range self.Length() {
					element := self.AtVal(i)
					result, err := vm.CallClosure(function, element)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					newList.List.SetAt(i, result)
				}
				return value.Ref(newList), value.Undefined
			}

			// callable is another value
			for i := range self.Length() {
				element := self.AtVal(i)
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, element)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				newList.List.SetAt(i, result)
			}
			return value.Ref(newList), value.Undefined
		},
		DefWithParameters(1),
	)

}

// ::Std::Sync::ArrayList::Iterator
func initSyncArrayListIterator() {
	// Instance methods
	c := value.SyncArrayListIteratorClass
	Def(
		c,
		"next",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.SyncArrayListIterator)
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
			self := args[0].AsReference().(value.SyncArrayListIterator)
			self.Reset()
			return args[0], value.Undefined
		},
	)
}

func SyncArrayListConcat(thread *Thread, l value.SyncArrayList, other value.Value) (value.Value, value.Value) {
	switch l := l.(type) {
	case *value.SyncArrayListOfValue:
		return value.RefErr(SyncArrayListOfValueConcat(thread, l, other))
	case value.SyncArrayList:
		return SyncArrayListInterfaceConcat(thread, l, other)
	default:
		panic(fmt.Sprintf("invalid arraylist type for concat: %T", l))
	}
}

func SyncArrayListOfValueConcat(thread *Thread, l *value.SyncArrayListOfValue, other value.Value) (*value.SyncArrayListOfValue, value.Value) {
	result, ok := l.ConcatOk(other)
	if ok {
		return result, value.Undefined
	}

	l.Lock()
	defer l.Unlock()

	lLen := l.Length()
	newList := make(value.ArrayListOfValue, lLen, lLen)
	copy(newList, *l.List)

	i := 0
	for element, err := range Iterate(thread, other) {
		if err.IsNotUndefined() {
			return nil, err
		}

		newList[lLen+i] = element
		i++
	}

	return &value.SyncArrayListOfValue{
		List: &newList,
	}, value.Undefined
}

func SyncArrayListInterfaceConcat(thread *Thread, l value.SyncArrayList, other value.Value) (value.Value, value.Value) {
	result, ok := l.ConcatValOk(other)
	if ok {
		return result, value.Undefined
	}

	l.Lock()
	defer l.Unlock()

	length := l.Length()
	newList := make(value.ArrayListOfValue, length, length)
	for i, element := range l.ArrayList().Elements() {
		newList[i] = element
	}

	i := 0
	for element, err := range Iterate(thread, other) {
		if err.IsNotUndefined() {
			return value.Undefined, err
		}

		newList[length+i] = element
		i++
	}

	return (&value.SyncArrayListOfValue{
		List: &newList,
	}).ToValue(), value.Undefined
}
