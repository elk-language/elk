package test

import (
	"fmt"

	"github.com/elk-language/elk/ext"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func Init() {
	ext.Register("std/test", runtimeInit, typecheckerInit)
}

func runtimeInit() {
	testModule := value.NewModule()
	value.StdModule.AddConstantString("Test", value.Ref(testModule))

	assertionsMixin := value.NewMixin()
	testModule.AddConstantString("Assertions", value.Ref(assertionsMixin))
	assertionsMixin.SingletonClass().IncludeMixin(assertionsMixin)

	assertionErrorClass := value.NewClassWithOptions(
		value.ClassWithSuperclass(value.ErrorClass),
	)
	testModule.AddConstantString("AssertionError", value.Ref(assertionErrorClass))

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
				message = fmt.Sprintf("value `%s` is not truthy", gotInspect)
			} else {
				message = args[2].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
				message = fmt.Sprintf("value `%s` is not falsy", gotInspect)
			} else {
				message = args[2].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
				message = fmt.Sprintf("value `%s` is not `true`", gotInspect)
			} else {
				message = args[2].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
				message = fmt.Sprintf("value `%s` is not `false`", gotInspect)
			} else {
				message = args[2].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
				message = fmt.Sprintf("value `%s` is not `nil`", gotInspect)
			} else {
				message = args[2].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
				message = args[2].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
				gotInspect, err := vm.Inspect(v, argGot)
				if !err.IsUndefined() {
					return value.Undefined, err
				}
				message = fmt.Sprintf(
					"value `%s` is not an instance of `%s` or its parents",
					gotInspect,
					argClass.Name,
				)
			} else {
				message = args[3].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
					gotInspect,
					argClass.Name,
				)
			} else {
				message = args[3].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
					gotInspect,
					expectedInspect,
				)
			} else {
				message = args[3].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
					gotInspect,
					expectedInspect,
				)
			} else {
				message = args[3].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
					aInspect,
					bInspect,
				)
			} else {
				message = args[3].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
					aInspect,
					bInspect,
				)
			} else {
				message = args[3].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
					aInspect,
					bInspect,
				)
			} else {
				message = args[3].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
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
					aInspect,
					bInspect,
				)
			} else {
				message = args[3].AsReference().(value.String).String()
			}

			err = value.Ref(
				value.NewError(
					assertionErrorClass,
					message,
				),
			)
			return value.Undefined, err
		},
		vm.DefWithParameters(3),
	)
}

// TODO: test macros
func typecheckerInit(checker types.Checker) {
	env := checker.Env()
	testModule := env.Std().DefineModule("", symbol.Test, env)
	assertionsMixin := testModule.DefineMixin("", false, value.ToSymbol("Assertions"), env)
	expressionNodeMixin := env.ExpressionNode()

	n := assertionsMixin.Singleton()
	types.DefMacro(
		n,
		"",
		"assert!",
		[]*types.Parameter{
			types.NewParameter(
				value.ToSymbol("expression"),
				expressionNodeMixin,
				types.NormalParameterKind,
				false,
			),
		},
		expressionNodeMixin,
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			expr := args[1].AsReference().(ast.ExpressionNode)

			assertionsConstant := ast.NewConstantLookupNode(
				position.ZeroLocation,
				ast.NewConstantLookupNode(
					position.ZeroLocation,
					ast.NewConstantLookupNode(
						position.ZeroLocation,
						nil,
						ast.NewPublicConstantNode(
							position.ZeroLocation,
							"Std",
						),
					),
					ast.NewPublicConstantNode(
						position.ZeroLocation,
						"Test",
					),
				),
				ast.NewPublicConstantNode(
					position.ZeroLocation,
					"Assertions",
				),
			)

			binExpr, ok := expr.(*ast.BinaryExpressionNode)
			if !ok {
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_truthy",
					),
					[]ast.ExpressionNode{ast.Unhygienic(expr)},
					nil,
				)
				return value.Ref(result), value.Undefined
			}

			switch binExpr.Op.Type {
			case token.EQUAL_EQUAL:
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_equal",
					),
					[]ast.ExpressionNode{
						ast.Unhygienic(binExpr.Left),
						ast.Unhygienic(binExpr.Right),
					},
					nil,
				)
				return value.Ref(result), value.Undefined
			case token.NOT_EQUAL:
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_not_equal",
					),
					[]ast.ExpressionNode{
						ast.Unhygienic(binExpr.Left),
						ast.Unhygienic(binExpr.Right),
					},
					nil,
				)
				return value.Ref(result), value.Undefined
			case token.GREATER:
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_greater",
					),
					[]ast.ExpressionNode{
						ast.Unhygienic(binExpr.Left),
						ast.Unhygienic(binExpr.Right),
					},
					nil,
				)
				return value.Ref(result), value.Undefined
			case token.GREATER_EQUAL:
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_greater_equal",
					),
					[]ast.ExpressionNode{
						ast.Unhygienic(binExpr.Left),
						ast.Unhygienic(binExpr.Right),
					},
					nil,
				)
				return value.Ref(result), value.Undefined
			case token.LESS:
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_less",
					),
					[]ast.ExpressionNode{
						ast.Unhygienic(binExpr.Left),
						ast.Unhygienic(binExpr.Right),
					},
					nil,
				)
				return value.Ref(result), value.Undefined
			case token.LESS_EQUAL:
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_less_equal",
					),
					[]ast.ExpressionNode{
						ast.Unhygienic(binExpr.Left),
						ast.Unhygienic(binExpr.Right),
					},
					nil,
				)
				return value.Ref(result), value.Undefined
			case token.ISA_OP:
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_is_a",
					),
					[]ast.ExpressionNode{
						ast.Unhygienic(binExpr.Left),
						ast.Unhygienic(binExpr.Right),
					},
					nil,
				)
				return value.Ref(result), value.Undefined
			case token.REVERSE_ISA_OP:
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_is_a",
					),
					[]ast.ExpressionNode{
						ast.Unhygienic(binExpr.Right),
						ast.Unhygienic(binExpr.Left),
					},
					nil,
				)
				return value.Ref(result), value.Undefined
			case token.INSTANCE_OF_OP:
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_instance_of",
					),
					[]ast.ExpressionNode{
						ast.Unhygienic(binExpr.Left),
						ast.Unhygienic(binExpr.Right),
					},
					nil,
				)
				return value.Ref(result), value.Undefined
			case token.REVERSE_INSTANCE_OF_OP:
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_instance_of",
					),
					[]ast.ExpressionNode{
						ast.Unhygienic(binExpr.Right),
						ast.Unhygienic(binExpr.Left),
					},
					nil,
				)
				return value.Ref(result), value.Undefined
			default:
				result := ast.NewMethodCallNode(
					position.ZeroLocation,
					assertionsConstant,
					token.New(position.ZeroLocation, token.DOT),
					ast.NewPublicIdentifierNode(
						position.ZeroLocation,
						"assert_truthy",
					),
					[]ast.ExpressionNode{ast.Unhygienic(expr)},
					nil,
				)
				return value.Ref(result), value.Undefined
			}
		},
	)
}
