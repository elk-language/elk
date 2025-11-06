package test

import (
	"fmt"
	"strings"

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
				gotInspect, err := vm.Inspect(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf("value `%s` is not truthy", gotInspect.AsString().String())
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
				gotInspect, err := vm.Inspect(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf("value `%s` is not falsy", gotInspect.AsString().String())
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
				gotInspect, err := vm.Inspect(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf("value `%s` is not `true`", gotInspect.AsString().String())
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
				gotInspect, err := vm.Inspect(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf("value `%s` is not `false`", gotInspect.AsString().String())
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
				gotInspect, err := vm.Inspect(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf("value `%s` is not `nil`", gotInspect.AsString().String())
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
				message = "value should not be `nil`"
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
		"assert_stdout",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			expected := args[1].AsString()
			callable := args[2]

			prevStdout := v.Stdout
			var buff strings.Builder
			v.Stdout = &buff

			v.CallCallable(callable)

			v.Stdout = prevStdout

			got := value.String(buff.String())
			if expected != got {
				err = value.Ref(
					value.NewError(
						AssertionErrorClass,
						fmt.Sprintf(
							"invalid stdout, expected: `%s`, got: `%s`",
							expected.Inspect(),
							got.Inspect(),
						),
					),
				)
				return value.Undefined, err
			}

			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"assert_stderr",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			expected := args[1].AsString()
			callable := args[2]

			prevStderr := v.Stderr
			var buff strings.Builder
			v.Stderr = &buff

			v.CallCallable(callable)

			v.Stderr = prevStderr

			got := value.String(buff.String())
			if expected != got {
				err = value.Ref(
					value.NewError(
						AssertionErrorClass,
						fmt.Sprintf(
							"invalid stderr, expected: `%s`, got: `%s`",
							expected.Inspect(),
							got.Inspect(),
						),
					),
				)
				return value.Undefined, err
			}

			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"assert_matches_regex",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argRegex := (*value.Regex)(args[1].Pointer())
			argString := args[2].AsString()

			if argRegex.MatchesString(argString.String()) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[3].IsUndefined() {
				message = fmt.Sprintf(
					"string `%s` does not match regex `%s`",
					argString.Inspect(),
					argRegex.Inspect(),
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
		"assert_is_a",
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			argClass := (*value.Class)(args[1].Pointer())
			argGot := args[2]

			if value.IsA(argGot, argClass) {
				return value.Nil, value.Undefined
			}

			var message string
			if args[3].IsUndefined() {
				gotInspect, err := vm.Inspect(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not an instance of `%s` or its parents",
					gotInspect.AsString().String(),
					argClass.Name,
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
				gotInspect, err := vm.Inspect(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not an instance of `%s`",
					gotInspect.AsString().String(),
					argClass.Name,
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
				gotInspect, err := vm.Inspect(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				expectedInspect, err := vm.Inspect(v, argExpected)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not equal to `%s`",
					gotInspect.AsString().String(),
					expectedInspect.AsString().String(),
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
				gotInspect, err := vm.Inspect(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				expectedInspect, err := vm.Inspect(v, argExpected)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` should not be equal to `%s`",
					gotInspect.AsString().String(),
					expectedInspect.AsString().String(),
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
				aInspect, err := vm.Inspect(v, argA)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				bInspect, err := vm.Inspect(v, argB)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not greater than `%s`",
					aInspect.AsString().String(),
					bInspect.AsString().String(),
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
				aInspect, err := vm.Inspect(v, argA)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				bInspect, err := vm.Inspect(v, argB)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not greater than or equal to `%s`",
					aInspect.AsString().String(),
					bInspect.AsString().String(),
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
				aInspect, err := vm.Inspect(v, argA)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				bInspect, err := vm.Inspect(v, argB)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not less than `%s`",
					aInspect.AsString().String(),
					bInspect.AsString().String(),
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
				aInspect, err := vm.Inspect(v, argA)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				bInspect, err := vm.Inspect(v, argB)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not less than or equal to `%s`",
					aInspect.AsString().String(),
					bInspect.AsString().String(),
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
