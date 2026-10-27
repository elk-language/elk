package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// ::Std::ArrayList
func initArrayList() {
	// Instance methods
	c := value.ArrayListClass
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayList)
			iterator := self.IterList()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayList)
			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayList)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayList)
			return value.SmallInt(self.LeftCapacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"box_of",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayList)
			other := args[1]
			b, err := self.BoxOfVal(other)
			return b, err
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"immutable_box_of",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayList)
			other := args[1]
			b, err := self.ImmutableBoxOfVal(other)
			return b, err
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"[]",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayList)
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
			self := args[0].AsReference().(value.ArrayList)
			other := args[1]
			return self.TryAtVal(other)
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"[]=",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayList)
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
			self := args[0].AsReference().(value.ArrayList)
			other := args[1]
			return ArrayListConcat(thread, self, other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayList)
			other := args[1]
			return self.RepeatVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayList)
			switch other := args[1].SafeAsReference().(type) {
			case value.ArrayList:
				equal, err := ArrayTupleEqual(vm, self, other)
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
			self := args[0].AsReference().(value.ArrayList)
			other := args[1]
			eq, err := ArrayTupleLaxEqualValue(vm, self, other)
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
			self := args[0].AsReference().(value.ArrayList)
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
			self := args[0].AsReference().(value.ArrayList)
			values := args[1]
			for val := range Iterate(vm, values) {
				err := self.AppendVal(val)
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
			self := args[0].AsReference().(value.ArrayList)
			val := args[1]

			var removed bool
			for i := 0; i < self.Length(); i++ {
				elem := self.AtVal(i)
				isEqual, err := Equal(vm, elem, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(isEqual) {
					self.RemoveAt(i)
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
			self := args[0].AsReference().(value.ArrayList)
			val := args[1]
			return value.Nil, self.RemoveAtVal(val)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"remove_at_range",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayList)
			rng := args[1]
			err := value.ArrayListInterfaceRemoveAtRange(self, rng)
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
			self := args[0].AsReference().(value.ArrayList)
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
			self := args[0].AsReference().(value.ArrayList)
			callable := args[1]
			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for i := range self.Length() {
					element := self.AtVal(i)
					result, err := vm.CallClosure(function, element)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					err = self.SetAtVal(i, result)
					if err.IsNotUndefined() {
						return value.Undefined, err
					}
				}
				return self.ToValue(), value.Undefined
			}

			// callable is another value
			for i := range self.Length() {
				element := self.AtVal(i)
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, element)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				err = self.SetAtVal(i, result)
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
			self := args[0].AsReference().(value.ArrayList)
			callable := args[1]
			newList := value.NewArrayListOfValueWithLength(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for i := range self.Length() {
					element := self.AtVal(i)
					result, err := vm.CallClosure(function, element)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					newList.SetAt(i, result)
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
				newList.SetAt(i, result)
			}
			return value.Ref(newList), value.Undefined
		},
		DefWithParameters(1),
	)

}

// ::Std::ArrayList::Iterator
func initArrayListIterator() {
	// Instance methods
	c := value.ArrayListIteratorClass
	Def(
		c,
		"next",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayListIterator)
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
			self := args[0].AsReference().(value.ArrayListIterator)
			self.Reset()
			return args[0], value.Undefined
		},
	)
}

func ArrayListConcat(thread *Thread, l value.ArrayList, other value.Value) (value.Value, value.Value) {
	switch l := l.(type) {
	case *value.ArrayListOfValue:
		return value.RefErr(ArrayListOfValueConcat(thread, l, other))
	case value.ArrayList:
		return ArrayListInterfaceConcat(thread, l, other)
	default:
		panic(fmt.Sprintf("invalid arraylist type for concat: %T", l))
	}
}

func ArrayListOfValueConcat(thread *Thread, l *value.ArrayListOfValue, other value.Value) (*value.ArrayListOfValue, value.Value) {
	result, ok := l.ConcatOk(other)
	if ok {
		return result, value.Undefined
	}

	newList := make(value.ArrayListOfValue, len(*l), len(*l))
	copy(newList, *l)

	i := 0
	for element, err := range Iterate(thread, other) {
		if err.IsNotUndefined() {
			return nil, err
		}

		newList[len(*l)+i] = element
		i++
	}

	return &newList, value.Undefined
}

func ArrayListInterfaceConcat(thread *Thread, l value.ArrayList, other value.Value) (value.Value, value.Value) {
	result, ok := l.ConcatValOk(other)
	if ok {
		return result, value.Undefined
	}

	length := l.Length()
	newList := make(value.ArrayListOfValue, length, length)
	for i, element := range l.Elements() {
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

	return (&newList).ToValue(), value.Undefined
}
