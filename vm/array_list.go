package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// ::Std::ArrayList
func init() {
	// Instance methods
	c := &value.ArrayListClass.MethodContainer
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			iterator := value.NewArrayListIterator(self)
			return iterator, nil
		},
	)
	Def(
		c,
		"capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			return value.SmallInt(self.Capacity()), nil
		},
	)
	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			return value.SmallInt(self.Length()), nil
		},
	)
	Def(
		c,
		"left_capacity",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			return value.SmallInt(self.LeftCapacity()), nil
		},
	)
	Def(
		c,
		"[]",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			other := args[1]
			return value.ToValueErr(self.Subscript(other))
		},
		DefWithParameters("key"),
		DefWithSealed(),
	)
	Def(
		c,
		"[]=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			key := args[1]
			val := args[2]
			err := self.SubscriptSet(key, val)
			if err != nil {
				return nil, err
			}
			return val, nil
		},
		DefWithParameters("key", "value"),
		DefWithSealed(),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			other := args[1]
			return value.ToValueErr(self.Concat(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			other := args[1]
			return value.ToValueErr(self.Repeat(other))
		},
		DefWithParameters("n"),
		DefWithSealed(),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			contains, err := ArrayListContains(vm, self, args[1])
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters("value"),
		DefWithSealed(),
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			switch other := args[1].(type) {
			case *value.ArrayList:
				equal, err := ArrayListEqual(vm, self, other)
				if err != nil {
					return nil, err
				}
				return value.ToElkBool(equal), nil
			default:
				return value.False, nil
			}
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"=~",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			switch other := args[1].(type) {
			case *value.ArrayList:
				equal, err := ArrayListEqual(vm, self, other)
				if err != nil {
					return nil, err
				}
				return value.ToElkBool(equal), nil
			case *value.ArrayTuple:
				equal, err := ArrayListEqual(vm, self, (*value.ArrayList)(other))
				if err != nil {
					return nil, err
				}
				return value.ToElkBool(equal), nil
			default:
				return value.False, nil
			}
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"grow",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			nValue := args[1]
			n, ok := value.IntToGoInt(nValue)
			if !ok && n == -1 {
				return nil, value.NewTooLargeCapacityError(nValue.Inspect())
			}
			if n < 0 {
				return nil, value.NewNegativeCapacityError(nValue.Inspect())
			}
			if !ok {
				return nil, value.NewCapacityTypeError(nValue.Inspect())
			}
			self.Grow(n)
			return self, nil
		},
		DefWithParameters("new_slots"),
		DefWithSealed(),
	)
	Def(
		c,
		"append",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			values := args[1].(*value.ArrayList)
			self.Append(*values...)
			return self, nil
		},
		DefWithParameters("values"),
		DefWithPositionalRestParameter(),
		DefWithSealed(),
	)
	Def(
		c,
		"<<",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			self.Append(args[1])
			return self, nil
		},
		DefWithParameters("value"),
		DefWithSealed(),
	)

	Def(
		c,
		"map_mut",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			callable := args[1]
			// callable is a closure
			if function, ok := callable.(*Closure); ok {
				for i := range self.Length() {
					element := self.At(i)
					result, err := vm.CallClosure(function, element)
					if err != nil {
						return nil, err
					}
					self.SetAt(i, result)
				}
				return self, nil
			}

			// callable is another value
			for i := range self.Length() {
				element := self.At(i)
				result, err := vm.CallMethod(callSymbol, callable, element)
				if err != nil {
					return nil, err
				}
				self.SetAt(i, result)
			}
			return self, nil
		},
		DefWithParameters("func"),
	)

	Def(
		c,
		"map",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			callable := args[1]
			newList := value.NewArrayListWithLength(self.Length())

			// callable is a closure
			if function, ok := callable.(*Closure); ok {
				for i := range self.Length() {
					element := self.At(i)
					result, err := vm.CallClosure(function, element)
					if err != nil {
						return nil, err
					}
					newList.SetAt(i, result)
				}
				return newList, nil
			}

			// callable is another value
			for i := range self.Length() {
				element := self.At(i)
				result, err := vm.CallMethod(callSymbol, callable, element)
				if err != nil {
					return nil, err
				}
				newList.SetAt(i, result)
			}
			return newList, nil
		},
		DefWithParameters("func"),
	)

}

// ::Std::ArrayList::Iterator
func init() {
	// Instance methods
	c := &value.ArrayListIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayListIterator)
			return self.Next()
		},
	)
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
	)
}

func ArrayListContains(vm *VM, list *value.ArrayList, val value.Value) (bool, value.Value) {
	for _, element := range *list {
		equal, err := vm.CallMethod(symbol.OpEqual, element, val)
		if err != nil {
			return false, err
		}
		if value.Truthy(equal) {
			return true, nil
		}
	}
	return false, nil
}

func ArrayListEqual(vm *VM, x, y *value.ArrayList) (bool, value.Value) {
	xLen := x.Length()
	if xLen != y.Length() {
		return false, nil
	}

	for i := 0; i < xLen; i++ {
		equal, err := vm.CallMethod(symbol.OpEqual, (*x)[i], (*y)[i])
		if err != nil {
			return false, err
		}
		if value.Falsy(equal) {
			return false, nil
		}
	}
	return true, nil
}
