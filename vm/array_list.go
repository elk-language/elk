package vm

import (
	"github.com/elk-language/elk/value"
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
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			switch other := args[1].(type) {
			case *value.ArrayList:
				selfLen := self.Length()
				if selfLen != other.Length() {
					return value.False, nil
				}

				for i := 0; i < selfLen; i++ {
					equal, err := vm.CallMethod(equalSymbol, (*self)[i], (*other)[i])
					if err != nil {
						return nil, err
					}
					switch equal.(type) {
					case value.FalseType, value.NilType:
						return value.False, nil
					}
				}
				return value.True, nil
			case *value.ArrayTuple:
				selfLen := self.Length()
				if selfLen != other.Length() {
					return value.False, nil
				}

				for i := 0; i < selfLen; i++ {
					equal, err := vm.CallMethod(equalSymbol, (*self)[i], (*other)[i])
					if err != nil {
						return nil, err
					}
					switch equal.(type) {
					case value.FalseType, value.NilType:
						return value.False, nil
					}
				}
				return value.True, nil
			default:
				return value.False, nil
			}
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"===",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayList)
			switch other := args[1].(type) {
			case *value.ArrayList:
				selfLen := self.Length()
				if selfLen != other.Length() {
					return value.False, nil
				}

				for i := 0; i < selfLen; i++ {
					equal, err := vm.CallMethod(strictEqualSymbol, (*self)[i], (*other)[i])
					if err != nil {
						return nil, err
					}
					switch equal.(type) {
					case value.FalseType, value.NilType:
						return value.False, nil
					}
				}
				return value.True, nil
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
