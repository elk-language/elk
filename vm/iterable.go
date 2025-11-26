package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Iterable::FiniteBase
func initIterableFiniteBase() {
	// Instance methods
	c := &value.IterableFiniteBase.MethodContainer

	Def(
		c,
		"contains",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			val := args[1]

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				eq, err := Equal(vm, elem, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(eq) {
					return value.True, value.Undefined
				}
			}

			return value.False, value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"is_empty",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]

			var anyElements bool
			for _, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				anyElements = true
				break
			}

			return value.ToElkBool(!anyElements), value.Undefined
		},
	)

	Def(
		c,
		"first",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				return elem, value.Undefined
			}

			selfInspect, err := Inspect(vm, self)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			err = value.Ref(
				value.Errorf(
					value.IterableNotFoundError,
					"cannot get first element of `%s`",
					selfInspect.AsString().String(),
				),
			)
			return value.Undefined, err
		},
	)

	Def(
		c,
		"try_first",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				return elem, value.Undefined
			}

			return value.Nil, value.Undefined
		},
	)

	Def(
		c,
		"last",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]

			var last value.Value
			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				last = elem
			}

			if !last.IsUndefined() {
				return last, value.Undefined
			}

			selfInspect, err := Inspect(vm, self)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			err = value.Ref(
				value.Errorf(
					value.IterableNotFoundError,
					"cannot get last element of `%s`",
					selfInspect.AsString().String(),
				),
			)
			return value.Undefined, err
		},
	)

	Def(
		c,
		"try_last",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]

			var last value.Value
			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				last = elem
			}

			if !last.IsUndefined() {
				return last, value.Undefined
			}

			return value.Nil, value.Undefined
		},
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
		"count",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			var counter value.SmallInt

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				keep, err := vm.CallCallable(fn, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(keep) {
					counter++
				}
			}

			return counter.ToValue(), value.Undefined
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
		"any",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				ok, err := vm.CallCallable(fn, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(ok) {
					return value.True, value.Undefined
				}
			}

			return value.False, value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"every",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				ok, err := vm.CallCallable(fn, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Falsy(ok) {
					return value.False, value.Undefined
				}
			}

			return value.True, value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"find",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				ok, err := vm.CallCallable(fn, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(ok) {
					return elem, value.Undefined
				}
			}

			selfInspect, err := Inspect(vm, self)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			err = value.Ref(
				value.Errorf(
					value.IterableNotFoundError,
					"could not find element of `%s`",
					selfInspect.AsString().String(),
				),
			)
			return value.Undefined, err
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"try_find",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				ok, err := vm.CallCallable(fn, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(ok) {
					return elem, value.Undefined
				}
			}

			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"index_of",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			val := args[1]

			var index value.SmallInt

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				ok, err := Equal(vm, elem, val)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(ok) {
					return index.ToValue(), value.Undefined
				}

				index++
			}

			return value.SmallInt(-1).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"find_index",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			var index value.SmallInt

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				ok, err := vm.CallCallable(fn, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if value.Truthy(ok) {
					return index.ToValue(), value.Undefined
				}

				index++
			}

			return value.SmallInt(-1).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"drop",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			count := args[1].AsInt()

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

	Def(
		c,
		"reduce",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			fn := args[1]

			var accumulator value.Value

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				if accumulator.IsUndefined() {
					accumulator = elem
					continue
				}

				newValue, err := vm.CallCallable(fn, accumulator, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				accumulator = newValue
			}

			return accumulator, value.Undefined
		},
		DefWithParameters(1),
	)

	Def(
		c,
		"fold",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]
			initial := args[1]
			fn := args[2]

			accumulator := initial

			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				newValue, err := vm.CallCallable(fn, accumulator, elem)
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				accumulator = newValue
			}

			return accumulator, value.Undefined
		},
		DefWithParameters(2),
	)

	Def(
		c,
		"to_list",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]

			var result value.ArrayList
			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				result.Append(elem)
			}

			return value.Ref(&result), value.Undefined
		},
	)
	Alias(c, "to_collection", "to_list")

	Def(
		c,
		"to_tuple",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]

			var result value.ArrayTuple
			for elem, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				result.Append(elem)
			}

			return value.Ref(&result), value.Undefined
		},
	)
	Alias(c, "to_immutable_collection", "to_tuple")

}

// Std::Iterable::Base
func initIterableBase() {
	// Instance methods
	c := &value.IterableBase.MethodContainer

	Def(
		c,
		"length",
		func(vm *VM, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0]

			var counter value.SmallInt

			for _, err := range Iterate(vm, self) {
				if !err.IsUndefined() {
					return value.Undefined, err
				}

				counter++
			}

			return counter.ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
}
