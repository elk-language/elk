package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::ArrayList
func init() {
	// Instance methods
	c := &value.ListClass.MethodContainer
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
