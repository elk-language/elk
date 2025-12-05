package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Std::Collection::Base
func initCollection() {
	// Instance methods
	c := &value.CollectionBaseMixin.MethodContainer

	Def(
		c,
		"remove_all",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			values := args[1]

			var deleted bool

			for val, err := range Iterate(vm, values) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				deletedElement, err := vm.CallMethodByName(symbol.L_remove, self, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(deletedElement) {
					deleted = true
				}
			}

			return value.ToElkBool(deleted), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"append",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			values := args[1]

			for val, err := range Iterate(vm, values) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				_, err := vm.CallMethodByName(symbol.L_push, self, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
			}

			return self, value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"<<",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			val := args[1]

			_, err = vm.CallMethodByName(symbol.L_push, self, val)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			return self, value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"map",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			var result value.ArrayList

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

			var result value.ArrayList

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

			var result value.ArrayList

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
			var result value.ArrayList

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

			var result value.ArrayList
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
			var result value.ArrayList

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

			var result value.ArrayList

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
