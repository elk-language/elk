package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Std::Tuple
func initTuple() {
	// Instance methods
	c := &value.TupleMixin.MethodContainer

	Def(
		c,
		"[]",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			index := args[1]

			return vm.CallMethodByName(symbol.L_at, self, index)
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"[]@1",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			rangeVal := args[1]

			return vm.CallMethodByName(symbol.L_view, self, rangeVal)
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"slice",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			rangeVal := args[1].AsReference()

			lengthVal, err := vm.CallMethodByName(symbol.L_length, self)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			length := lengthVal.AsInt()

			var start int
			end := lengthVal.AsInt()

			switch r := rangeVal.(type) {
			case *value.ClosedRange:
				start = r.Start.AsInt()
				end = r.End.AsInt()
			case *value.LeftOpenRange:
				start = r.Start.AsInt() + 1
				end = r.End.AsInt()
			case *value.RightOpenRange:
				start = r.Start.AsInt()
				end = r.End.AsInt() - 1
			case *value.OpenRange:
				start = r.Start.AsInt() + 1
				end = r.End.AsInt() - 1
			case *value.BeginlessOpenRange:
				end = r.End.AsInt() - 1
			case *value.BeginlessClosedRange:
				end = r.End.AsInt()
			case *value.EndlessOpenRange:
				start = r.Start.AsInt() + 1
			case *value.EndlessClosedRange:
				start = r.Start.AsInt()
			}

			start, err = value.NormalizeArrayIndex(start, length)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			end, err = value.NormalizeArrayIndex(end, length)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			var result value.ArrayTuple
			for i := start; i <= end; i++ {
				element, err := vm.CallMethodByName(symbol.L_at, self, value.SmallInt(i).ToValue())
				if err.IsNotUndefined() {
					return value.Undefined, err
				}

				result = append(result, element)
			}

			return value.Ref(&result), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"try_at",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			index := args[1]

			val, err := vm.CallMethodByName(symbol.L_at, self, index)
			if !err.IsUndefined() {
				if err.Class() == value.IndexErrorClass {
					return value.Nil, value.Undefined
				}
				return value.Undefined, err
			}

			return val, value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"+",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			other := args[1]

			var result value.ArrayTuple
			for elem, err := range Iterate(vm, self) {
				if err.IsNotUndefined() {
					return value.Undefined, err
				}
				result = append(result, elem)
			}
			for elem, err := range Iterate(vm, other) {
				if err.IsNotUndefined() {
					return value.Undefined, err
				}
				result = append(result, elem)
			}

			return value.Ref(&result), value.Undefined
		},
		DefWithParameters(1),
	)

}
