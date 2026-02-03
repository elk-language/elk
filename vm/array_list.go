package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// ::Std::ArrayList
func initArrayList() {
	// Instance methods
	c := &value.ArrayListClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			iterator := value.NewArrayListOfValueIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Def(
		c,
		"capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			return value.SmallInt(self.Capacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			return value.SmallInt(self.LeftCapacity()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"box_of",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			other := args[1]
			b, err := self.BoxOfVal(other)
			return value.Ref(b), err
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"immutable_box_of",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			other := args[1]
			b, err := self.ImmutableBoxOfVal(other)
			return value.Ref(b), err
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"[]",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			other := args[1]
			return self.Subscript(other)
		},
		DefWithParameters(1),
	)
	Alias(c, "at", "[]")
	Def(
		c,
		"[]=",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
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
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			other := args[1]
			return value.RefErr(self.Concat(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			other := args[1]
			return value.RefErr(self.Repeat(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			contains, err := ArrayListOfValueContains(vm, self, args[1])
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			switch other := args[1].SafeAsReference().(type) {
			case *value.ArrayListOfValue:
				equal, err := ArrayListOfValueEqual(vm, self, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.ToElkBool(equal), value.Undefined
			default:
				return value.False, value.Undefined
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			switch other := args[1].SafeAsReference().(type) {
			case *value.ArrayListOfValue:
				equal, err := ArrayListOfValueEqual(vm, self, other)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.ToElkBool(equal), value.Undefined
			case *value.ArrayTupleOfValue:
				equal, err := ArrayListOfValueEqual(vm, self, (*value.ArrayListOfValue)(other))
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.ToElkBool(equal), value.Undefined
			default:
				return value.False, value.Undefined
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"grow",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
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
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"append",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			values := args[1]
			for val := range Iterate(vm, values) {
				self.Append(val)
			}
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"remove",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			val := args[1]

			var removed bool
			for i := 0; i < self.Length(); i++ {
				elem := self.At(i)
				isEqual, err := Equal(vm, elem, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(isEqual) {
					self.RemoveAt(i)
					removed = true
				}
			}

			return value.ToElkBool(removed), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"remove_at",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			val := args[1].AsInt()
			return value.Nil, self.RemoveAtErr(val)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			self.Append(args[1])
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)
	Alias(c, "push", "<<")

	Def(
		c,
		"map_mut",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			callable := args[1]
			// callable is a closure
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for i := range self.Length() {
					element := self.At(i)
					result, err := vm.CallClosure(function, element)
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					self.SetAt(i, result)
				}
				return value.Ref(self), value.Undefined
			}

			// callable is another value
			for i := range self.Length() {
				element := self.At(i)
				result, err := vm.CallMethodByName(callSymbol, callable, element)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				self.SetAt(i, result)
			}
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValue)(args[0].Pointer())
			callable := args[1]
			newList := value.NewArrayListOfValueWithLength(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for i := range self.Length() {
					element := self.At(i)
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
				element := self.At(i)
				result, err := vm.CallMethodByName(callSymbol, callable, element)
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
	c := &value.ArrayListOfValueIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.ArrayListOfValueIterator)(args[0].Pointer())
			return self.Next()
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
			self := (*value.ArrayListOfValueIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
		},
	)
}

func ArrayListOfValueContains(vm *Thread, list *value.ArrayListOfValue, val value.Value) (bool, value.Value) {
	for _, element := range *list {
		equal, err := vm.CallMethodByName(symbol.OpEqual, element, val)
		if !err.IsUndefined() {
			return false, err
		}
		if value.Truthy(equal) {
			return true, value.Undefined
		}
	}
	return false, value.Undefined
}

func ArrayListOfValueEqual(vm *Thread, x, y *value.ArrayListOfValue) (bool, value.Value) {
	xLen := x.Length()
	if xLen != y.Length() {
		return false, value.Undefined
	}

	for i := 0; i < xLen; i++ {
		equal, err := vm.CallMethodByName(symbol.OpEqual, (*x)[i], (*y)[i])
		if !err.IsUndefined() {
			return false, err
		}
		if value.Falsy(equal) {
			return false, value.Undefined
		}
	}
	return true, value.Undefined
}
