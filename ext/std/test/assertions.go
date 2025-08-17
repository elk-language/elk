package test

import (
	"fmt"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

var AssertionErrorClass *value.Class // Std::Test::AssertionError

func initAssertions(testModule *value.Module) {
	assertionsMixin := value.NewMixin()
	testModule.AddConstantString("Assertions", value.Ref(assertionsMixin))
	assertionsMixin.SingletonClass().IncludeMixin(assertionsMixin)

	AssertionErrorClass = value.NewClassWithOptions(
		value.ClassWithSuperclass(value.ErrorClass),
	)
	testModule.AddConstantString("AssertionError", value.Ref(AssertionErrorClass))

	c := &assertionsMixin.MethodContainer
	vm.Def(
		c,
		"assert_truthy",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argGot := args[1]

			if value.Truthy(argGot) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[2].IsUndefined() {
				gotInspect, err := vm.InspectWithColor(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf("value `%s` is not truthy", gotInspect)
			} else {
				message = args[2].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(2),
	)
	vm.Def(
		c,
		"assert_falsy",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argGot := args[1]

			if value.Falsy(argGot) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[2].IsUndefined() {
				gotInspect, err := vm.InspectWithColor(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf("value `%s` is not falsy", gotInspect)
			} else {
				message = args[2].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"assert_true",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argGot := args[1]

			if argGot.IsTrue() {
				return value.Nil, value.Undefined
			}

			var message string
			if args[2].IsUndefined() {
				gotInspect, err := vm.InspectWithColor(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf("value `%s` is not `%s`", gotInspect, lexer.Colorize("true"))
			} else {
				message = args[2].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(2),
	)
	vm.Def(
		c,
		"assert_false",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argGot := args[1]

			if argGot.IsFalse() {
				return value.Nil, value.Undefined
			}

			var message string
			if args[2].IsUndefined() {
				gotInspect, err := vm.InspectWithColor(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf("value `%s` is not `%s`", gotInspect, lexer.Colorize("false"))
			} else {
				message = args[2].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(2),
	)
	vm.Def(
		c,
		"assert_nil",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argGot := args[1]

			if argGot.IsNil() {
				return value.Nil, value.Undefined
			}

			var message string
			if args[2].IsUndefined() {
				gotInspect, err := vm.InspectWithColor(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf("value `%s` is not `%s`", gotInspect, lexer.Colorize("nil"))
			} else {
				message = args[2].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(2),
	)
	vm.Def(
		c,
		"assert_not_nil",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argGot := args[1]

			if !argGot.IsNil() {
				return value.Nil, value.Undefined
			}

			var message string
			if args[2].IsUndefined() {
				message = fmt.Sprintf("value should not be `%s`", lexer.Colorize("nil"))
			} else {
				message = args[2].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(2),
	)
	vm.Def(
		c,
		"assert_is_a",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argClass := (*value.Class)(args[1].Pointer())
			argGot := args[2]

			if value.IsA(argGot, argClass) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[3].IsUndefined() {
				gotInspect, err := vm.InspectWithColor(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not an instance of `%s` or its parents",
					gotInspect,
					lexer.Colorize(argClass.Name),
				)
			} else {
				message = args[3].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(3),
	)
	vm.Def(
		c,
		"assert_instance_of",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argClass := (*value.Class)(args[1].Pointer())
			argGot := args[2]

			if value.InstanceOf(argGot, argClass) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[3].IsUndefined() {
				gotInspect, err := vm.InspectWithColor(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not an instance of `%s`",
					gotInspect,
					lexer.Colorize(argClass.Name),
				)
			} else {
				message = args[3].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(3),
	)
	vm.Def(
		c,
		"assert_equal",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argExpected := args[1]
			argGot := args[2]

			areEqual, err := vm.Equal(v, argExpected, argGot)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			if value.Truthy(areEqual) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[3].IsUndefined() {
				gotInspect, err := vm.InspectWithColor(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				expectedInspect, err := vm.InspectWithColor(v, argExpected)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not equal to `%s`",
					gotInspect,
					expectedInspect,
				)
			} else {
				message = args[3].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(3),
	)
	vm.Def(
		c,
		"assert_not_equal",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argExpected := args[1]
			argGot := args[2]

			areEqual, err := vm.Equal(v, argExpected, argGot)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			if value.Falsy(areEqual) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[3].IsUndefined() {
				gotInspect, err := vm.InspectWithColor(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				expectedInspect, err := vm.InspectWithColor(v, argExpected)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` should not be equal to `%s`",
					gotInspect,
					expectedInspect,
				)
			} else {
				message = args[3].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"assert_greater",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argA := args[1]
			argB := args[2]

			isGreater, err := vm.GreaterThan(v, argA, argB)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			if value.Truthy(isGreater) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[3].IsUndefined() {
				aInspect, err := vm.InspectWithColor(v, argA)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				bInspect, err := vm.InspectWithColor(v, argB)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not greater than `%s`",
					aInspect,
					bInspect,
				)
			} else {
				message = args[3].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(3),
	)
	vm.Def(
		c,
		"assert_greater_equal",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argA := args[1]
			argB := args[2]

			isGreaterEqual, err := vm.GreaterThanEqual(v, argA, argB)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			if value.Truthy(isGreaterEqual) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[3].IsUndefined() {
				aInspect, err := vm.InspectWithColor(v, argA)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				bInspect, err := vm.InspectWithColor(v, argB)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not greater than or equal to `%s`",
					aInspect,
					bInspect,
				)
			} else {
				message = args[3].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(3),
	)
	vm.Def(
		c,
		"assert_less",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argA := args[1]
			argB := args[2]

			isLess, err := vm.LessThan(v, argA, argB)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			if value.Truthy(isLess) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[3].IsUndefined() {
				aInspect, err := vm.InspectWithColor(v, argA)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				bInspect, err := vm.InspectWithColor(v, argB)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not less than `%s`",
					aInspect,
					bInspect,
				)
			} else {
				message = args[3].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(3),
	)
	vm.Def(
		c,
		"assert_less_equal",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argA := args[1]
			argB := args[2]

			isLessEqual, err := vm.LessThanEqual(v, argA, argB)
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			if value.Truthy(isLessEqual) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[3].IsUndefined() {
				aInspect, err := vm.InspectWithColor(v, argA)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				bInspect, err := vm.InspectWithColor(v, argB)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not less than or equal to `%s`",
					aInspect,
					bInspect,
				)
			} else {
				message = args[3].AsString().String()
			}

			err = value.Ref(
				value.NewError(
					AssertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(3),
	)
}
