package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::ImmutableCollection::Base
func initImmutableCollection() {
	// Instance methods
	c := &value.ImmutableCollectionBaseMixin.MethodContainer

	Def(
		c,
		"map",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			var result value.ArrayTuple

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				newElem, err := vm.CallCallable(fn, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				result.Append(newElem)
			}

			return value.Ref(&result), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"filter",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			var result value.ArrayTuple

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				keep, err := vm.CallCallable(fn, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(keep) {
					result.Append(elem)
				}
			}

			return value.Ref(&result), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"reject",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			var result value.ArrayTuple

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				keep, err := vm.CallCallable(fn, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Falsy(keep) {
					result.Append(elem)
				}
			}

			return value.Ref(&result), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"drop",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			count := args[1].AsInt()

			if count < 0 {
				return value.Undefined, value.Ref(
					value.Errorf(
						value.OutOfRangeErrorClass,
						"tried to drop a negative amount of values `%d` from an iterable",
						count,
					),
				)
			}
			var result value.ArrayTuple

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if count > 0 {
					count--
					continue
				}

				result.Append(elem)
			}

			return value.Ref(&result), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"drop_while",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			var result value.ArrayTuple
			var collect bool

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if collect {
					result.Append(elem)
					continue
				}

				ok, err := vm.CallCallable(fn, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(ok) {
					continue
				}

				collect = true
				result.Append(elem)
			}

			return value.Ref(&result), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"take",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			count := args[1].AsInt()

			if count < 0 {
				return value.Undefined, value.Ref(
					value.Errorf(
						value.OutOfRangeErrorClass,
						"tried to take a negative amount of values `%d` from an iterable",
						count,
					),
				)
			}
			var result value.ArrayTuple

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if count <= 0 {
					break
				}

				count--
				result.Append(elem)
				continue
			}

			return value.Ref(&result), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"take_while",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			var result value.ArrayTuple

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				ok, err := vm.CallCallable(fn, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Falsy(ok) {
					break
				}

				result.Append(elem)
			}

			return value.Ref(&result), value.Undefined
		},
		DefWithParameters(1),
	)

}
