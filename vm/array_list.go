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
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			iterator := value.NewArrayListIterator(self)
			return value.Ref(iterator), value.Nil
		},
	)
	Def(
		c,
		"capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			return value.SmallInt(self.Capacity()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			return value.SmallInt(self.Length()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			return value.SmallInt(self.LeftCapacity()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"[]",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			other := args[1]
			return self.Subscript(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"[]=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			key := args[1]
			val := args[2]
			err := self.SubscriptSet(key, val)
			if !err.IsNil() {
				return value.Nil, err
			}
			return val, value.Nil
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			other := args[1]
			return value.RefErr(self.Concat(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			other := args[1]
			return value.RefErr(self.Repeat(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			contains, err := ArrayListContains(vm, self, args[1])
			if !err.IsNil() {
				return value.Nil, err
			}
			return value.ToElkBool(contains), value.Nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			switch other := args[1].SafeAsReference().(type) {
			case *value.ArrayList:
				equal, err := ArrayListEqual(vm, self, other)
				if !err.IsNil() {
					return value.Nil, err
				}
				return value.ToElkBool(equal), value.Nil
			default:
				return value.False, value.Nil
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			switch other := args[1].SafeAsReference().(type) {
			case *value.ArrayList:
				equal, err := ArrayListEqual(vm, self, other)
				if !err.IsNil() {
					return value.Nil, err
				}
				return value.ToElkBool(equal), value.Nil
			case *value.ArrayTuple:
				equal, err := ArrayListEqual(vm, self, (*value.ArrayList)(other))
				if !err.IsNil() {
					return value.Nil, err
				}
				return value.ToElkBool(equal), value.Nil
			default:
				return value.False, value.Nil
			}
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"grow",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			nValue := args[1]
			n, ok := value.IntToGoInt(nValue)
			if !ok && n == -1 {
				return value.Nil, value.Ref(value.NewTooLargeCapacityError(nValue.Inspect()))
			}
			if n < 0 {
				return value.Nil, value.Ref(value.NewNegativeCapacityError(nValue.Inspect()))
			}
			if !ok {
				return value.Nil, value.Ref(value.NewCapacityTypeError(nValue.Inspect()))
			}
			self.Grow(n)
			return value.Ref(self), value.Nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"append",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			values := args[1].MustReference().(*value.ArrayList)
			self.Append(*values...)
			return value.Ref(self), value.Nil
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<<",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			self.Append(args[1])
			return value.Ref(self), value.Nil
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map_mut",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			callable := args[1]
			// callable is a closure
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for i := range self.Length() {
					element := self.At(i)
					result, err := vm.CallClosure(function, element)
					if !err.IsNil() {
						return value.Nil, err
					}
					self.SetAt(i, result)
				}
				return value.Ref(self), value.Nil
			}

			// callable is another value
			for i := range self.Length() {
				element := self.At(i)
				result, err := vm.CallMethodByName(callSymbol, callable, element)
				if !err.IsNil() {
					return value.Nil, err
				}
				self.SetAt(i, result)
			}
			return value.Ref(self), value.Nil
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayList)
			callable := args[1]
			newList := value.NewArrayListWithLength(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for i := range self.Length() {
					element := self.At(i)
					result, err := vm.CallClosure(function, element)
					if !err.IsNil() {
						return value.Nil, err
					}
					newList.SetAt(i, result)
				}
				return value.Ref(newList), value.Nil
			}

			// callable is another value
			for i := range self.Length() {
				element := self.At(i)
				result, err := vm.CallMethodByName(callSymbol, callable, element)
				if !err.IsNil() {
					return value.Nil, err
				}
				newList.SetAt(i, result)
			}
			return value.Ref(newList), value.Nil
		},
		DefWithParameters(1),
	)

}

// ::Std::ArrayList::Iterator
func initArrayListIterator() {
	// Instance methods
	c := &value.ArrayListIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayListIterator)
			return self.Next()
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Nil
		},
	)
}

func ArrayListContains(vm *VM, list *value.ArrayList, val value.Value) (bool, value.Value) {
	for _, element := range *list {
		equal, err := vm.CallMethodByName(symbol.OpEqual, element, val)
		if !err.IsNil() {
			return false, err
		}
		if value.Truthy(equal) {
			return true, value.Nil
		}
	}
	return false, value.Nil
}

func ArrayListEqual(vm *VM, x, y *value.ArrayList) (bool, value.Value) {
	xLen := x.Length()
	if xLen != y.Length() {
		return false, value.Nil
	}

	for i := 0; i < xLen; i++ {
		equal, err := vm.CallMethodByName(symbol.OpEqual, (*x)[i], (*y)[i])
		if !err.IsNil() {
			return false, err
		}
		if value.Falsy(equal) {
			return false, value.Nil
		}
	}
	return true, value.Nil
}
