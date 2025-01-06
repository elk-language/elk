package vm

import (
	"github.com/elk-language/elk/value"
)

// ::Std::ArrayTuple
func initArrayTuple() {
	// Instance methods
	c := &value.ArrayTupleClass.MethodContainer
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayTuple)
			iterator := value.NewArrayTupleIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayTuple)
			return value.SmallInt(self.Length()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"[]",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayTuple)
			other := args[1]
			return self.Subscript(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"[]=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayTuple)
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
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayTuple)
			other := args[1]
			return self.Concat(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayTuple)
			other := args[1]
			return value.RefErr(self.Repeat(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayTuple)
			contains, err := ArrayTupleContains(vm, self, args[1])
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.ToElkBool(contains), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayTuple)
			switch other := args[1].SafeAsReference().(type) {
			case *value.ArrayList:
				equal, err := ArrayTupleEqual(vm, self, (*value.ArrayTuple)(other))
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				return value.ToElkBool(equal), value.Undefined
			case *value.ArrayTuple:
				equal, err := ArrayTupleEqual(vm, self, other)
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
		"==",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayTuple)
			switch other := args[1].SafeAsReference().(type) {
			case *value.ArrayTuple:
				equal, err := ArrayTupleEqual(vm, self, other)
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
		"map",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayTuple)
			callable := args[1]
			newTuple := value.NewArrayTupleWithLength(self.Length())

			// callable is a closure
			if function, ok := callable.SafeAsReference().(*Closure); ok {
				for i := range self.Length() {
					element := self.At(i)
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
				element := self.At(i)
				result, err := vm.CallMethodByName(callSymbol, callable, element)
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
	c := &value.ArrayTupleIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*value.ArrayTupleIterator)
			return self.Next()
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)

}

func ArrayTupleContains(vm *VM, tuple *value.ArrayTuple, val value.Value) (bool, value.Value) {
	return ArrayListContains(vm, (*value.ArrayList)(tuple), val)
}

func ArrayTupleEqual(vm *VM, x, y *value.ArrayTuple) (bool, value.Value) {
	return ArrayListEqual(vm, (*value.ArrayList)(x), (*value.ArrayList)(y))
}
