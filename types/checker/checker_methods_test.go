package checker

import (
	"testing"

	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/ast"
)

func TestModuleMethod(t *testing.T) {
	// globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
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
								nil,
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
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(0, 1, 1), P(13, 1, 14)), "argument `c` is missing in call to `baz`"),
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
								nil,
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
								nil,
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
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(8, 1, 9), P(12, 1, 13)), "expected type `Std::String` for parameter `bar`, got type `Std::Float(123.4)`"),
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
								nil,
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
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(0, 1, 1), P(26, 1, 27)), "expected 2 arguments, got 5"),
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
								nil,
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
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(0, 1, 1), P(18, 1, 19)), "argument `c` is missing in call to `baz`"),
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
								nil,
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
								nil,
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
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(14, 1, 15), P(23, 1, 24)), "expected type `Std::String` for parameter `bar`, got type `Std::Float(123.4)`"),
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
								nil,
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
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(18, 1, 19), P(23, 1, 24)), "duplicated argument `bar` in call to `baz`"),
				errors.NewError(L("<main>", P(18, 1, 19), P(23, 1, 24)), "expected type `Std::String` for parameter `bar`, got type `Std::Int(9)`"),
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
								nil,
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
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(18, 1, 19), P(21, 1, 22)), "duplicated argument `c` in call to `baz`"),
				errors.NewError(L("<main>", P(24, 1, 25), P(27, 1, 28)), "duplicated argument `c` in call to `baz`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, true)
		})
	}
}
