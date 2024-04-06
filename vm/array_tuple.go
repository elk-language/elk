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
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayTuple)
			return value.SmallInt(self.Length()), nil
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
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayTuple)
			contains, err := ArrayTupleContains(vm, self, args[1])
			if err != nil {
				return nil, err
			}
			return value.ToElkBool(contains), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"=~",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayTuple)
			switch other := args[1].(type) {
			case *value.ArrayList:
				equal, err := ArrayTupleEqual(vm, self, (*value.ArrayTuple)(other))
				if err != nil {
					return nil, err
				}
				return value.ToElkBool(equal), nil
			case *value.ArrayTuple:
				equal, err := ArrayTupleEqual(vm, self, other)
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
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.ArrayTuple)
			switch other := args[1].(type) {
			case *value.ArrayTuple:
				equal, err := ArrayTupleEqual(vm, self, other)
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

func ArrayTupleContains(vm *VM, tuple *value.ArrayTuple, val value.Value) (bool, value.Value) {
	return ArrayListContains(vm, (*value.ArrayList)(tuple), val)
}

func ArrayTupleEqual(vm *VM, x, y *value.ArrayTuple) (bool, value.Value) {
	return ArrayListEqual(vm, (*value.ArrayList)(x), (*value.ArrayList)(y))
}
