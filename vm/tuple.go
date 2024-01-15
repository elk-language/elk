package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::ArrayTuple
func init() {
	// Instance methods
	c := &value.ArrayTupleClass.MethodContainer
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayTuple)
			iterator := value.NewArrayTupleIterator(self)
			return iterator, nil
		},
	)
	Def(
		c,
		"[]",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayTuple)
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
			self := args[0].(*value.ArrayTuple)
			key := args[1]
			val := args[2]
			err := self.SubscriptSet(key, val)
			if err == nil {
				return val, nil
			}
			return nil, err
		},
		DefWithParameters("key", "val"),
		DefWithSealed(),
	)
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayTuple)
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
			self := args[0].(*value.ArrayTuple)
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
			self := args[0].(*value.ArrayTuple)
			switch other := args[1].(type) {
			case *value.List:
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
			self := args[0].(*value.ArrayTuple)
			switch other := args[1].(type) {
			case *value.ArrayTuple:
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

}

// ::Std::ArrayTupleIterator
func init() {
	// Instance methods
	c := &value.ArrayTupleIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayTupleIterator)
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
