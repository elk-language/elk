package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/ast"
	"github.com/elk-language/elk/value/symbol"
)

func TestMethodDefinitionOverride(t *testing.T) {
	tests := simplifiedTestTable{
		"add additional optional params": {
			input: `
				class Foo
					def baz(a: Int): Int then a
				end

				class Bar < Foo
					def baz(a: Int, b: Int? = nil): Int then a
				end
			`,
		},
		"invalid override": {
			input: `
				class Foo
					def baz(a: Int): Int then a
				end

				class Bar < Foo
					def baz(); end
				end
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(82, 7, 6), P(95, 7, 19)), "cannot redeclare method `baz` with less parameters\n  previous definition found in `Foo`, with signature: sig baz(a: Std::Int): Std::Int"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			simplifiedCheckerTest(tc, t)
		})
	}
}

func TestMethodDefinition(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"override the method with additional optional params": {
			before: `def baz(a: Int): Int then a`,
			input:  `def baz(a: Int, b: Int = 2): Int then a`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 1, 39)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(38, 1, 39)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(38, 1, 39)),
							"baz",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(13, 1, 14)),
									"a",
									false,
									ast.NewPublicConstantNode(
										S(P(11, 1, 12), P(13, 1, 14)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(16, 1, 17), P(25, 1, 26)),
									"b",
									false,
									ast.NewPublicConstantNode(
										S(P(19, 1, 20), P(21, 1, 22)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
									ast.NewIntLiteralNode(
										S(P(25, 1, 26), P(25, 1, 26)),
										"2",
										types.NewIntLiteral("2"),
									),
									ast.NormalParameterKind,
								),
							},
							ast.NewPublicConstantNode(
								S(P(29, 1, 30), P(31, 1, 32)),
								"Int",
								globalEnv.StdSubtypeString("Int"),
							),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(38, 1, 39), P(38, 1, 39)),
									ast.NewPublicIdentifierNode(
										S(P(38, 1, 39), P(38, 1, 39)),
										"a",
										globalEnv.StdSubtypeString("Int"),
									),
								),
							},
						),
					),
				},
			),
		},
		"override the method with different param name": {
			before: `def baz(a: Int): Int then a`,
			input:  `def baz(b: Int): Int then b`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							"baz",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(13, 1, 14)),
									"b",
									false,
									ast.NewPublicConstantNode(
										S(P(11, 1, 12), P(13, 1, 14)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
									nil,
									ast.NormalParameterKind,
								),
							},
							ast.NewPublicConstantNode(
								S(P(17, 1, 18), P(19, 1, 20)),
								"Int",
								globalEnv.StdSubtypeString("Int"),
							),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(26, 1, 27), P(26, 1, 27)),
									ast.NewPublicIdentifierNode(
										S(P(26, 1, 27), P(26, 1, 27)),
										"b",
										globalEnv.StdSubtypeString("Int"),
									),
								),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(8, 1, 9), P(13, 1, 14)), "cannot redeclare method `baz` with invalid parameter name, is `b`, should be `a`\n  previous definition found in `Std::Object`, with signature: sig baz(a: Std::Int): Std::Int"),
			},
		},
		"override the method with different param type": {
			before: `def baz(a: Int): Int then a`,
			input:  `def baz(a: Char): Int then 1`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							"baz",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(14, 1, 15)),
									"a",
									false,
									ast.NewPublicConstantNode(
										S(P(11, 1, 12), P(14, 1, 15)),
										"Char",
										globalEnv.StdSubtypeString("Char"),
									),
									nil,
									ast.NormalParameterKind,
								),
							},
							ast.NewPublicConstantNode(
								S(P(18, 1, 19), P(20, 1, 21)),
								"Int",
								globalEnv.StdSubtypeString("Int"),
							),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(27, 1, 28), P(27, 1, 28)),
									ast.NewIntLiteralNode(
										S(P(27, 1, 28), P(27, 1, 28)),
										"1",
										types.NewIntLiteral("1"),
									),
								),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(8, 1, 9), P(14, 1, 15)), "cannot redeclare method `baz` with invalid parameter type, is `Std::Char`, should be `Std::Int`\n  previous definition found in `Std::Object`, with signature: sig baz(a: Std::Int): Std::Int"),
			},
		},
		"override the method with different return type": {
			before: `def baz(a: Int): Int then a`,
			input:  "def baz(a: Int): Char then `a`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							"baz",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(13, 1, 14)),
									"a",
									false,
									ast.NewPublicConstantNode(
										S(P(11, 1, 12), P(13, 1, 14)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
									nil,
									ast.NormalParameterKind,
								),
							},
							ast.NewPublicConstantNode(
								S(P(17, 1, 18), P(20, 1, 21)),
								"Char",
								globalEnv.StdSubtypeString("Char"),
							),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(27, 1, 28), P(29, 1, 30)),
									ast.NewCharLiteralNode(
										S(P(27, 1, 28), P(29, 1, 30)),
										'a',
										types.NewCharLiteral('a'),
									),
								),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(17, 1, 18), P(20, 1, 21)), "cannot redeclare method `baz` with a different return type, is `Std::Char`, should be `Std::Int`\n  previous definition found in `Std::Object`, with signature: sig baz(a: Std::Int): Std::Int"),
			},
		},
		"methods get hoisted to the top": {
			input: `
			  foo()
				def foo; end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 3, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(16, 3, 5), P(28, 3, 17)),
						ast.NewMethodDefinitionNode(
							S(P(16, 3, 5), P(27, 3, 16)),
							"foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						S(P(6, 2, 6), P(11, 2, 11)),
						ast.NewReceiverlessMethodCallNode(
							S(P(6, 2, 6), P(10, 2, 10)),
							"foo",
							nil,
							types.Void{},
						),
					),
				},
			),
		},
		"methods can reference each other": {
			input: `
				def foo then bar()
				def bar then foo()
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(46, 3, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(23, 2, 23)),
						ast.NewMethodDefinitionNode(
							S(P(5, 2, 5), P(22, 2, 22)),
							"foo",
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(18, 2, 18), P(22, 2, 22)),
									ast.NewReceiverlessMethodCallNode(
										S(P(18, 2, 18), P(22, 2, 22)),
										"bar",
										nil,
										types.Void{},
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(28, 3, 5), P(46, 3, 23)),
						ast.NewMethodDefinitionNode(
							S(P(28, 3, 5), P(45, 3, 22)),
							"bar",
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(41, 3, 18), P(45, 3, 22)),
									ast.NewReceiverlessMethodCallNode(
										S(P(41, 3, 18), P(45, 3, 22)),
										"foo",
										nil,
										types.Void{},
									),
								),
							},
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, true)
		})
	}
}

func TestMethodCalls(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"call has the same return type as the method": {
			before: `
				module Foo
					def baz(a: Int): Int then a
				end
			`,
			input: `Foo.baz(5)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "5", types.NewIntLiteral("5")),
							},
							globalEnv.StdSubtypeString("Int"),
						),
					),
				},
			),
		},
		"cannot make nil-safe call on a non nilable receiver": {
			before: `
				module Foo
					def baz(a: Int): Int then a
				end
			`,
			input: `Foo?.baz(5)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							true,
							"baz",
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(S(P(9, 1, 10), P(9, 1, 10)), "5", types.NewIntLiteral("5")),
							},
							globalEnv.StdSubtypeString("Int"),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(0, 1, 1), P(10, 1, 11)), "cannot make a nil-safe call on type `Foo` which is not nilable"),
			},
		},
		"can make nil-safe call on a nilable receiver": {
			before: `
				module Foo
					def baz(a: Int): Int then a
				end
				const NilableFoo: Foo? = Foo
			`,
			input: `NilableFoo?.baz(5)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(9, 1, 10)),
								"NilableFoo",
								types.NewNilable(types.NewModule("Foo")),
							),
							true,
							"baz",
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "5", types.NewIntLiteral("5")),
							},
							types.NewNilable(globalEnv.StdSubtypeString("Int")),
						),
					),
				},
			),
		},
		"missing required argument": {
			before: `
				module Foo
					def baz(bar: String, c: Int); end
				end
			`,
			input: `Foo.baz("foo")`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(0, 1, 1), P(13, 1, 14)), "argument `c` is missing in call to `baz`"),
			},
		},
		"all required positional arguments": {
			before: `
				module Foo
					def baz(bar: String, c: Int); end
				end
			`,
			input: `Foo.baz("foo", 5)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewIntLiteralNode(S(P(15, 1, 16), P(15, 1, 16)), "5", types.NewIntLiteral("5")),
							},
							types.Void{},
						),
					),
				},
			),
		},
		"all required positional arguments with wrong type": {
			before: `
				module Foo
					def baz(bar: String, c: Int); end
				end
			`,
			input: `Foo.baz(123.4, 5)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"123.4",
									types.NewFloatLiteral("123.4"),
								),
								ast.NewIntLiteralNode(S(P(15, 1, 16), P(15, 1, 16)), "5", types.NewIntLiteral("5")),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String` for parameter `bar` in call to `baz`, got type `Std::Float(123.4)`"),
			},
		},
		"too many positional arguments": {
			before: `
				module Foo
					def baz(bar: String, c: Int); end
				end
			`,
			input: `Foo.baz("foo", 5, 28, 9, 0)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewIntLiteralNode(S(P(15, 1, 16), P(15, 1, 16)), "5", types.NewIntLiteral("5")),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(0, 1, 1), P(26, 1, 27)), "expected 2 arguments in call to `baz`, got 5"),
			},
		},
		"missing required argument with named argument": {
			before: `
				module Foo
					def baz(bar: String, c: Int); end
				end
			`,
			input: `Foo.baz(bar: "foo")`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(13, 1, 14), P(17, 1, 18)),
									"foo",
									types.NewStringLiteral("foo"),
								),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(0, 1, 1), P(18, 1, 19)), "argument `c` is missing in call to `baz`"),
			},
		},
		"all required named arguments": {
			before: `
				module Foo
					def baz(bar: String, c: Int); end
				end
			`,
			input: `Foo.baz(c: 5, bar: "foo")`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(19, 1, 20), P(23, 1, 24)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewIntLiteralNode(S(P(11, 1, 12), P(11, 1, 12)), "5", types.NewIntLiteral("5")),
							},
							types.Void{},
						),
					),
				},
			),
		},
		"all required named arguments with wrong type": {
			before: `
				module Foo
					def baz(bar: String, c: Int); end
				end
			`,
			input: `Foo.baz(c: 5, bar: 123.4)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(
									S(P(19, 1, 20), P(23, 1, 24)),
									"123.4",
									types.NewFloatLiteral("123.4"),
								),
								ast.NewIntLiteralNode(S(P(11, 1, 12), P(11, 1, 12)), "5", types.NewIntLiteral("5")),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(14, 1, 15), P(23, 1, 24)), "expected type `Std::String` for parameter `bar` in call to `baz`, got type `Std::Float(123.4)`"),
			},
		},
		"duplicated positional argument as named argument": {
			before: `
				module Foo
					def baz(bar: String, c: Int); end
				end
			`,
			input: `Foo.baz("foo", 5, bar: 9)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewIntLiteralNode(S(P(15, 1, 16), P(15, 1, 16)), "5", types.NewIntLiteral("5")),
								ast.NewIntLiteralNode(S(P(23, 1, 24), P(23, 1, 24)), "9", types.NewIntLiteral("9")),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(18, 1, 19), P(23, 1, 24)), "duplicated argument `bar` in call to `baz`"),
				error.NewError(L("<main>", P(18, 1, 19), P(23, 1, 24)), "expected type `Std::String` for parameter `bar` in call to `baz`, got type `Std::Int(9)`"),
			},
		},
		"duplicated named argument": {
			before: `
				module Foo
					def baz(bar: String, c: Int); end
				end
			`,
			input: `Foo.baz("foo", 2, c: 3, c: 9)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 1, 29)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(28, 1, 29)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewIntLiteralNode(S(P(15, 1, 16), P(15, 1, 16)), "2", types.NewIntLiteral("2")),
								ast.NewIntLiteralNode(S(P(21, 1, 22), P(21, 1, 22)), "3", types.NewIntLiteral("3")),
								ast.NewIntLiteralNode(S(P(27, 1, 28), P(27, 1, 28)), "9", types.NewIntLiteral("9")),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(18, 1, 19), P(21, 1, 22)), "duplicated argument `c` in call to `baz`"),
				error.NewError(L("<main>", P(24, 1, 25), P(27, 1, 28)), "duplicated argument `c` in call to `baz`"),
			},
		},
		"call with missing optional argument": {
			before: `
				module Foo
					def baz(bar: String, c: Int = 3); end
				end
			`,
			input: `Foo.baz("foo")`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewUndefinedLiteralNode(S(P(0, 1, 1), P(13, 1, 14))),
							},
							types.Void{},
						),
					),
				},
			),
		},
		"call with optional argument": {
			before: `
				module Foo
					def baz(bar: String, c: Int = 3); end
				end
			`,
			input: `Foo.baz("foo", 9)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewIntLiteralNode(S(P(15, 1, 16), P(15, 1, 16)), "9", types.NewIntLiteral("9")),
							},
							types.Void{},
						),
					),
				},
			),
		},
		"call with missing rest arguments": {
			before: `
				module Foo
					def baz(*b: Float); end
				end
			`,
			input: `Foo.baz`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewArrayTupleLiteralNode(
									S(P(0, 1, 1), P(6, 1, 7)),
									nil,
								),
							},
							types.Void{},
						),
					),
				},
			),
		},
		"call with rest arguments": {
			before: `
				module Foo
					def baz(*b: Float); end
				end
			`,
			input: `Foo.baz 1.2, 56.9, .5`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewArrayTupleLiteralNode(
									S(P(0, 1, 1), P(20, 1, 21)),
									[]ast.ExpressionNode{
										ast.NewFloatLiteralNode(
											S(P(8, 1, 9), P(10, 1, 11)),
											"1.2",
											types.NewFloatLiteral("1.2"),
										),
										ast.NewFloatLiteralNode(
											S(P(13, 1, 14), P(16, 1, 17)),
											"56.9",
											types.NewFloatLiteral("56.9"),
										),
										ast.NewFloatLiteralNode(
											S(P(19, 1, 20), P(20, 1, 21)),
											"0.5",
											types.NewFloatLiteral("0.5"),
										),
									},
								),
							},
							types.Void{},
						),
					),
				},
			),
		},
		"call with rest arguments with wrong type": {
			before: `
				module Foo
					def baz(*b: Float); end
				end
			`,
			input: `Foo.baz 1.2, 5, "foo", .5`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewArrayTupleLiteralNode(
									S(P(0, 1, 1), P(24, 1, 25)),
									[]ast.ExpressionNode{
										ast.NewFloatLiteralNode(
											S(P(8, 1, 9), P(10, 1, 11)),
											"1.2",
											types.NewFloatLiteral("1.2"),
										),
										ast.NewIntLiteralNode(
											S(P(13, 1, 14), P(13, 1, 14)),
											"5",
											types.NewIntLiteral("5"),
										),
										ast.NewDoubleQuotedStringLiteralNode(
											S(P(16, 1, 17), P(20, 1, 21)),
											"foo",
											types.NewStringLiteral("foo"),
										),
										ast.NewFloatLiteralNode(
											S(P(23, 1, 24), P(24, 1, 25)),
											"0.5",
											types.NewFloatLiteral("0.5"),
										),
									},
								),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(13, 1, 14), P(13, 1, 14)), "expected type `Std::Float` for rest parameter `*b` in call to `baz`, got type `Std::Int(5)`"),
				error.NewError(L("<main>", P(16, 1, 17), P(20, 1, 21)), "expected type `Std::Float` for rest parameter `*b` in call to `baz`, got type `Std::String(\"foo\")`"),
			},
		},
		"call with rest argument given by name": {
			before: `
				module Foo
					def baz(*b: Float); end
				end
			`,
			input: `Foo.baz b: []`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewArrayTupleLiteralNode(
									S(P(0, 1, 1), P(12, 1, 13)),
									nil,
								),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(8, 1, 9), P(12, 1, 13)), "nonexistent parameter `b` given in call to `baz`"),
			},
		},
		"call with required post arguments": {
			before: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
			`,
			input: `Foo.baz("foo", 3)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewArrayTupleLiteralNode(
									S(P(0, 1, 1), P(16, 1, 17)),
									nil,
								),
								ast.NewIntLiteralNode(
									S(P(15, 1, 16), P(15, 1, 16)),
									"3",
									types.NewIntLiteral("3"),
								),
							},
							types.Void{},
						),
					),
				},
			),
		},
		"call with missing post argument": {
			before: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
			`,
			input: `Foo.baz("foo")`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewArrayTupleLiteralNode(
									S(P(0, 1, 1), P(13, 1, 14)),
									nil,
								),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(0, 1, 1), P(13, 1, 14)), "argument `c` is missing in call to `baz`"),
			},
		},
		"call with rest and post arguments": {
			before: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
			`,
			input: `Foo.baz("foo", 2.5, .9, 128.1, 3)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewArrayTupleLiteralNode(
									S(P(0, 1, 1), P(32, 1, 33)),
									[]ast.ExpressionNode{
										ast.NewFloatLiteralNode(
											S(P(15, 1, 16), P(17, 1, 18)),
											"2.5",
											types.NewFloatLiteral("2.5"),
										),
										ast.NewFloatLiteralNode(
											S(P(20, 1, 21), P(21, 1, 22)),
											"0.9",
											types.NewFloatLiteral("0.9"),
										),
										ast.NewFloatLiteralNode(
											S(P(24, 1, 25), P(28, 1, 29)),
											"128.1",
											types.NewFloatLiteral("128.1"),
										),
									},
								),
								ast.NewIntLiteralNode(
									S(P(31, 1, 32), P(31, 1, 32)),
									"3",
									types.NewIntLiteral("3"),
								),
							},
							types.Void{},
						),
					),
				},
			),
		},
		"call with rest and post arguments and wrong type in post": {
			before: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
			`,
			input: `Foo.baz("foo", 2.5, .9, 128.1, 3.2)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(34, 1, 35)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(34, 1, 35)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(34, 1, 35)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewArrayTupleLiteralNode(
									S(P(0, 1, 1), P(34, 1, 35)),
									[]ast.ExpressionNode{
										ast.NewFloatLiteralNode(
											S(P(15, 1, 16), P(17, 1, 18)),
											"2.5",
											types.NewFloatLiteral("2.5"),
										),
										ast.NewFloatLiteralNode(
											S(P(20, 1, 21), P(21, 1, 22)),
											"0.9",
											types.NewFloatLiteral("0.9"),
										),
										ast.NewFloatLiteralNode(
											S(P(24, 1, 25), P(28, 1, 29)),
											"128.1",
											types.NewFloatLiteral("128.1"),
										),
									},
								),
								ast.NewFloatLiteralNode(
									S(P(31, 1, 32), P(33, 1, 34)),
									"3.2",
									types.NewFloatLiteral("3.2"),
								),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(31, 1, 32), P(33, 1, 34)), "expected type `Std::Int` for parameter `c` in call to `baz`, got type `Std::Float(3.2)`"),
			},
		},
		"call with rest and post arguments and wrong type in rest": {
			before: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
			`,
			input: `Foo.baz("foo", 212, .9, '282', 3)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewArrayTupleLiteralNode(
									S(P(0, 1, 1), P(32, 1, 33)),
									[]ast.ExpressionNode{
										ast.NewIntLiteralNode(
											S(P(15, 1, 16), P(17, 1, 18)),
											"212",
											types.NewIntLiteral("212"),
										),
										ast.NewFloatLiteralNode(
											S(P(20, 1, 21), P(21, 1, 22)),
											"0.9",
											types.NewFloatLiteral("0.9"),
										),
										ast.NewRawStringLiteralNode(
											S(P(24, 1, 25), P(28, 1, 29)),
											"282",
											types.NewStringLiteral("282"),
										),
									},
								),
								ast.NewIntLiteralNode(
									S(P(31, 1, 32), P(31, 1, 32)),
									"3",
									types.NewIntLiteral("3"),
								),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(15, 1, 16), P(17, 1, 18)), "expected type `Std::Float` for rest parameter `*b` in call to `baz`, got type `Std::Int(212)`"),
				error.NewError(L("<main>", P(24, 1, 25), P(28, 1, 29)), "expected type `Std::Float` for rest parameter `*b` in call to `baz`, got type `Std::String(\"282\")`"),
			},
		},
		"call with rest arguments and missing post argument": {
			before: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
			`,
			input: `Foo.baz("foo", 2.5, .9, 128.1)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewArrayTupleLiteralNode(
									S(P(0, 1, 1), P(29, 1, 30)),
									[]ast.ExpressionNode{
										ast.NewFloatLiteralNode(
											S(P(15, 1, 16), P(17, 1, 18)),
											"2.5",
											types.NewFloatLiteral("2.5"),
										),
										ast.NewFloatLiteralNode(
											S(P(20, 1, 21), P(21, 1, 22)),
											"0.9",
											types.NewFloatLiteral("0.9"),
										),
									},
								),
								ast.NewFloatLiteralNode(
									S(P(24, 1, 25), P(28, 1, 29)),
									"128.1",
									types.NewFloatLiteral("128.1"),
								),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(24, 1, 25), P(28, 1, 29)), "expected type `Std::Int` for parameter `c` in call to `baz`, got type `Std::Float(128.1)`"),
			},
		},
		"call with named post argument": {
			before: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
			`,
			input: `Foo.baz("foo", c: 3)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewArrayTupleLiteralNode(
									S(P(0, 1, 1), P(19, 1, 20)),
									nil,
								),
								ast.NewIntLiteralNode(
									S(P(18, 1, 19), P(18, 1, 19)),
									"3",
									types.NewIntLiteral("3"),
								),
							},
							types.Void{},
						),
					),
				},
			),
		},
		"call with named pre rest argument": {
			before: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
			`,
			input: `Foo.baz(bar: "foo", c: 3)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							nil,
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(0, 1, 1), P(24, 1, 25)), "expected 1... positional arguments in call to `baz`, got 0"),
			},
		},
		"call without named rest arguments": {
			before: `
				module Foo
					def baz(bar: String, c: Int, **rest: Int); end
				end
			`,
			input: `Foo.baz("foo", 5)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewIntLiteralNode(
									S(P(15, 1, 16), P(15, 1, 16)),
									"5",
									types.NewIntLiteral("5"),
								),
								ast.NewHashRecordLiteralNode(
									S(P(0, 1, 1), P(16, 1, 17)),
									nil,
								),
							},
							types.Void{},
						),
					),
				},
			),
		},
		"call with named rest arguments": {
			before: `
				module Foo
					def baz(bar: String, c: Int, **rest: Int); end
				end
			`,
			input: `Foo.baz("foo", d: 25, c: 5, e: 11)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewIntLiteralNode(
									S(P(25, 1, 26), P(25, 1, 26)),
									"5",
									types.NewIntLiteral("5"),
								),
								ast.NewHashRecordLiteralNode(
									S(P(0, 1, 1), P(33, 1, 34)),
									[]ast.ExpressionNode{
										ast.NewSymbolKeyValueExpressionNode(
											S(P(15, 1, 16), P(19, 1, 20)),
											"d",
											ast.NewIntLiteralNode(
												S(P(18, 1, 19), P(19, 1, 20)),
												"25",
												types.NewIntLiteral("25"),
											),
										),
										ast.NewSymbolKeyValueExpressionNode(
											S(P(28, 1, 29), P(32, 1, 33)),
											"e",
											ast.NewIntLiteralNode(
												S(P(31, 1, 32), P(32, 1, 33)),
												"11",
												types.NewIntLiteral("11"),
											),
										),
									},
								),
							},
							types.Void{},
						),
					),
				},
			),
		},
		"call with named rest arguments with wrong type": {
			before: `
				module Foo
					def baz(bar: String, c: Int, **rest: Int); end
				end
			`,
			input: `Foo.baz("foo", d: .2, c: 5, e: .1)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewModule("Foo"),
							),
							false,
							"baz",
							[]ast.ExpressionNode{
								ast.NewDoubleQuotedStringLiteralNode(
									S(P(8, 1, 9), P(12, 1, 13)),
									"foo",
									types.NewStringLiteral("foo"),
								),
								ast.NewIntLiteralNode(
									S(P(25, 1, 26), P(25, 1, 26)),
									"5",
									types.NewIntLiteral("5"),
								),
								ast.NewHashRecordLiteralNode(
									S(P(0, 1, 1), P(33, 1, 34)),
									[]ast.ExpressionNode{
										ast.NewSymbolKeyValueExpressionNode(
											S(P(15, 1, 16), P(19, 1, 20)),
											"d",
											ast.NewFloatLiteralNode(
												S(P(18, 1, 19), P(19, 1, 20)),
												"0.2",
												types.NewFloatLiteral("0.2"),
											),
										),
										ast.NewSymbolKeyValueExpressionNode(
											S(P(28, 1, 29), P(32, 1, 33)),
											"e",
											ast.NewFloatLiteralNode(
												S(P(31, 1, 32), P(32, 1, 33)),
												"0.1",
												types.NewFloatLiteral("0.1"),
											),
										),
									},
								),
							},
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(15, 1, 16), P(19, 1, 20)), "expected type `Std::Int` for named rest parameter `**rest` in call to `baz`, got type `Std::Float(0.2)`"),
				error.NewError(L("<main>", P(28, 1, 29), P(32, 1, 33)), "expected type `Std::Int` for named rest parameter `**rest` in call to `baz`, got type `Std::Float(0.1)`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, true)
		})
	}
}

func TestInitDefinition(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"define in outer context": {
			input: `init; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewInvalidNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(0, 1, 1), P(8, 1, 9)), "init definitions cannot appear outside of classes"),
			},
		},
		"define in module": {
			input: `
				module Foo
					init; end
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 4, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(38, 4, 8)),
						ast.NewModuleDeclarationNode(
							S(P(5, 2, 5), P(37, 4, 7)),
							ast.NewPublicConstantNode(
								S(P(12, 2, 12), P(14, 2, 14)),
								"Foo",
								types.NewModule("Foo"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 3, 6), P(30, 3, 15)),
									ast.NewInvalidNode(
										S(P(21, 3, 6), P(29, 3, 14)),
										nil,
									),
								),
							},
							types.NewModule("Foo"),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(21, 3, 6), P(29, 3, 14)), "init definitions cannot appear outside of classes"),
			},
		},
		"define in class": {
			input: `
				class Foo
					init; end
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(37, 4, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(37, 4, 8)),
						ast.NewClassDeclarationNode(
							S(P(5, 2, 5), P(36, 4, 7)),
							false,
							false,
							ast.NewPublicConstantNode(
								S(P(11, 2, 11), P(13, 2, 13)),
								"Foo",
								types.NewClass(
									"Foo",
									globalEnv.StdSubtypeClass(symbol.Object),
								),
							),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(20, 3, 6), P(29, 3, 15)),
									ast.NewInitDefinitionNode(
										S(P(20, 3, 6), P(28, 3, 14)),
										nil,
										nil,
										nil,
									),
								),
							},
							types.NewClass(
								"Foo",
								globalEnv.StdSubtypeClass(symbol.Object),
							),
						),
					),
				},
			),
		},
		"with parameters": {
			input: `
				class Foo
					init(a: Int); end
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(45, 4, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(45, 4, 8)),
						ast.NewClassDeclarationNode(
							S(P(5, 2, 5), P(44, 4, 7)),
							false,
							false,
							ast.NewPublicConstantNode(
								S(P(11, 2, 11), P(13, 2, 13)),
								"Foo",
								types.NewClass(
									"Foo",
									globalEnv.StdSubtypeClass(symbol.Object),
								),
							),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(20, 3, 6), P(37, 3, 23)),
									ast.NewInitDefinitionNode(
										S(P(20, 3, 6), P(36, 3, 22)),
										[]ast.ParameterNode{
											ast.NewMethodParameterNode(
												S(P(25, 3, 11), P(30, 3, 16)),
												"a",
												false,
												ast.NewPublicConstantNode(
													S(P(28, 3, 14), P(30, 3, 16)),
													"Int",
													globalEnv.StdSubtypeString("Int"),
												),
												nil,
												ast.NormalParameterKind,
											),
										},
										nil,
										nil,
									),
								),
							},
							types.NewClass(
								"Foo",
								globalEnv.StdSubtypeClass(symbol.Object),
							),
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, true)
		})
	}
}

func TestConstructorCall(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"instantiate a class without a constructor": {
			before: `
				class Foo; end
			`,
			input: `Foo()`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewClass(
									"Foo",
									globalEnv.StdSubtypeClass(symbol.Object),
								),
							),
							nil,
							types.NewClass(
								"Foo",
								globalEnv.StdSubtypeClass(symbol.Object),
							),
						),
					),
				},
			),
		},
		"instantiate a class with a constructor": {
			before: `
				class Foo
					init(a: Int); end
				end
			`,
			input: `Foo(1)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewClass(
									"Foo",
									globalEnv.StdSubtypeClass(symbol.Object),
								),
							),
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"1",
									types.NewIntLiteral("1"),
								),
							},
							types.NewClass(
								"Foo",
								globalEnv.StdSubtypeClass(symbol.Object),
							),
						),
					),
				},
			),
		},
		"instantiate a class with a constructor with a wrong type": {
			before: `
				class Foo
					init(a: String); end
				end
			`,
			input: `Foo(1)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewClass(
									"Foo",
									globalEnv.StdSubtypeClass(symbol.Object),
								),
							),
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"1",
									types.NewIntLiteral("1"),
								),
							},
							types.NewClass(
								"Foo",
								globalEnv.StdSubtypeClass(symbol.Object),
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(4, 1, 5), P(4, 1, 5)), "expected type `Std::String` for parameter `a` in call to `#init`, got type `Std::Int(1)`"),
			},
		},
		"instantiate a class with an inherited constructor": {
			before: `
				class Bar
					init(a: Int); end
				end

				class Foo < Bar; end
			`,
			input: `Foo(1)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewClass(
									"Foo",
									types.NewClass(
										"Bar",
										globalEnv.StdSubtypeClass(symbol.Object),
									),
								),
							),
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"1",
									types.NewIntLiteral("1"),
								),
							},
							types.NewClass(
								"Foo",
								types.NewClass(
									"Bar",
									globalEnv.StdSubtypeClass(symbol.Object),
								),
							),
						),
					),
				},
			),
		},
		"instantiate a class with an inherited constructor with a wrong type": {
			before: `
				class Bar
					init(a: String); end
				end

				class Foo < Bar; end
			`,
			input: `Foo(1)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							ast.NewPublicConstantNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"Foo",
								types.NewClass(
									"Foo",
									types.NewClass(
										"Bar",
										globalEnv.StdSubtypeClass(symbol.Object),
									),
								),
							),
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"1",
									types.NewIntLiteral("1"),
								),
							},
							types.NewClass(
								"Foo",
								types.NewClass(
									"Bar",
									globalEnv.StdSubtypeClass(symbol.Object),
								),
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(4, 1, 5), P(4, 1, 5)), "expected type `Std::String` for parameter `a` in call to `#init`, got type `Std::Int(1)`"),
			},
		},
		"call a method on an instantiated instance": {
			before: `
				class Foo
					init(a: String); end

					def bar; end
				end

				var foo = Foo("foo")
			`,
			input: `foo.bar`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							ast.NewPublicIdentifierNode(
								S(P(0, 1, 1), P(2, 1, 3)),
								"foo",
								types.NewClass(
									"Foo",
									globalEnv.StdSubtypeClass(symbol.Object),
								),
							),
							false,
							"bar",
							nil,
							types.Void{},
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, true)
		})
	}
}

func TestMethodInheritance(t *testing.T) {
	tests := simplifiedTestTable{
		"call a method inherited from superclass": {
			input: `
				class Foo
					def baz(a: Int): Int then a
				end

				class Bar < Foo; end
				var bar = Bar()
				bar.baz(5)
			`,
		},
		"call a method inherited from mixin": {
			input: `
				mixin Bar
					def baz(a: Int): Int then a
				end

				class Foo
					include Bar
				end

				var foo = Foo()
				foo.baz(5)
			`,
		},
		"call a method on a mixin type": {
			input: `
				mixin Bar
					def baz(a: Int): Int then a
				end

				class Foo
					include Bar
				end

				var bar: Bar = Foo()
				bar.baz(5)
			`,
		},
		"call an inherited method on a mixin type": {
			input: `
				mixin Baz
					def baz(a: Int): Int then a
				end

				mixin Bar
				  include Baz
				end

				class Foo
					include Bar
				end

				var bar: Bar = Foo()
				bar.baz(5)
			`,
		},
		"call a class instance method on a mixin type": {
			input: `
				mixin Bar
					def baz(a: Int): Int then a
				end

				class Foo
					include Bar

					def foo; end
				end

				var bar: Bar = Foo()
				bar.foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(145, 13, 5), P(151, 13, 11)), "method `foo` is not defined on type `Bar`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			simplifiedCheckerTest(tc, t)
		})
	}
}
