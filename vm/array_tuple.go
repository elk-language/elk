package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// ::Std::ArrayTuple
func initArrayTuple() {
	// Instance methods
	c := value.ArrayTupleClass
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			iterator := self.IterTuple()
			return iterator.ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_array_tuple",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"to_array_list",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			return self.ToArrayList().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_sync_array_list",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			return self.ToSyncArrayList().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
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
			self := args[0].AsReference().(value.ArrayTuple)
			other := args[1]
			return self.TryAtVal(other)
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"+",
		func(thread *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			other := args[1]
			return ArrayTupleConcat(thread, self, other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			other := args[1]
			return self.RepeatVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"immutable_box_of",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			other := args[1]
			return self.ImmutableBoxOfVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
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
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			switch other := args[1].SafeAsReference().(type) {
			case value.ArrayTuple:
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
		"map",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTuple)
			callable := args[1]
			newTuple := value.NewArrayTupleOfValueWithLength(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(Closure); ok {
				for i := range self.Length() {
					element := self.AtVal(i)
					result, err := vm.CallClosure(function, element)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					newTuple.SetAt(i, result)
				}
				return value.Ref(newTuple), value.Undefined
			}

			// callable is another value
			for i := range self.Length() {
				element := self.AtVal(i)
				result, err := vm.CallMethodByName(value.S(symbol.L_call), callable, element)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				newTuple.SetAt(i, result)
			}
			return value.Ref(newTuple), value.Undefined
		},
		DefWithParameters(1),
	)

}

// ::Std::ArrayTupleIterator
func initArrayTupleIterator() {
	// Instance methods
	c := value.ArrayTupleIteratorClass
	Def(
		c,
		"next",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(value.ArrayTupleIterator)
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
			self := args[0].AsReference().(value.ArrayTupleIterator)
			self.Reset()
			return args[0], value.Undefined
		},
	)

}

func TupleEqual(vm *Thread, x, y value.Tuple) (bool, value.Value) {
	xLen := x.Length()
	if xLen != y.Length() {
		return false, value.Undefined
	}

	for i := 0; i < xLen; i++ {
		equal, err := Equal(vm, x.AtVal(i), y.AtVal(i))
		if !err.IsUndefined() {
			return false, err
		}
		if value.Falsy(equal) {
			return false, value.Undefined
		}
	}
	return true, value.Undefined
}

func ArrayTupleLaxEqual(vm *Thread, x, y value.Tuple) (bool, value.Value) {
	xLen := x.Length()
	if xLen != y.Length() {
		return false, value.Undefined
	}

	for i := 0; i < xLen; i++ {
		equal, err := LaxEqual(vm, x.AtVal(i), y.AtVal(i))
		if !err.IsUndefined() {
			return false, err
		}
		if value.Falsy(equal) {
			return false, value.Undefined
		}
	}
	return true, value.Undefined
}

func TupleLaxEqualValue(vm *Thread, x value.Tuple, y value.Value) (bool, value.Value) {
	switch narrowY := y.SafeAsReference().(type) {
	case value.Tuple:
		return ArrayTupleLaxEqual(vm, x, narrowY)
	default:
		if !value.IsA(y, value.TupleMixin) {
			return false, value.Undefined
		}

		yLengthVal, err := vm.CallMethodByName(value.S(symbol.L_length), y)
		if err.IsNotUndefined() {
			return false, err
		}

		yLength := yLengthVal.AsInt()
		if x.Length() != yLength {
			return false, value.Undefined
		}

		i := 0
		for yVal, err := range Iterate(vm, y) {
			if err.IsNotUndefined() {
				return false, err
			}

			xVal := x.AtVal(i)
			eq, err := LaxEqual(vm, xVal, yVal)
			if err.IsNotUndefined() {
				return false, err
			}

			if value.Falsy(eq) {
				return false, value.Undefined
			}
			i++
		}

		return true, value.Undefined
	}
}

func ArrayTupleConcat(thread *Thread, l value.ArrayTuple, other value.Value) (value.Value, value.Value) {
	switch l := l.(type) {
	case *value.ArrayTupleOfValue:
		return ArrayTupleOfValueConcat(thread, l, other)
	case value.ArrayTuple:
		return ArrayTupleInterfaceConcat(thread, l, other)
	default:
		panic(fmt.Sprintf("invalid arraytuple type for concat: %T", l))
	}
}

func ArrayTupleOfValueConcat(thread *Thread, l *value.ArrayTupleOfValue, other value.Value) (value.Value, value.Value) {
	result, ok := l.ConcatValOk(other)
	if ok {
		return result, value.Undefined
	}
	length := len(*l)

	newList := make(value.ArrayTupleOfValue, length, length)
	copy(newList, *l)

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

func ArrayTupleInterfaceConcat(thread *Thread, l value.ArrayTuple, other value.Value) (value.Value, value.Value) {
	result, ok := l.ConcatValOk(other)
	if ok {
		return result, value.Undefined
	}

	length := l.Length()
	newList := make(value.ArrayTupleOfValue, length, length)
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
