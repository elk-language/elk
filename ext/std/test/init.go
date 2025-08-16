package test

import (
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
	testModule := initTest()
	initAssertions(testModule)
}

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
