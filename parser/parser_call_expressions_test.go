package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
)

func TestConstructorCall(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "foo = Foo()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewConstructorCallNode(
								S(P(6, 1, 7), P(10, 1, 11)),
								ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have an empty argument list": {
			input: "Foo()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a private constant as the class": {
			input: "_Foo()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							ast.NewPrivateConstantNode(S(P(0, 1, 1), P(3, 1, 4)), "_Foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a constant lookup as the class": {
			input: "Foo::Bar()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							ast.NewConstantLookupNode(
								S(P(0, 1, 1), P(7, 1, 8)),
								ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
								ast.NewPublicConstantNode(S(P(5, 1, 6), P(7, 1, 8)), "Bar"),
							),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have positional arguments": {
			input: "Foo(.1, 'foo', :bar)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 1, 9), P(12, 1, 13)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 1, 16), P(18, 1, 19)), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have named arguments": {
			input: "Foo(bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(4, 1, 5), P(12, 1, 13)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(9, 1, 10), P(12, 1, 13)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(15, 1, 16), P(23, 1, 24)),
									"elk",
									ast.NewTrueLiteralNode(S(P(20, 1, 21), P(23, 1, 24))),
								),
							},
						),
					),
				},
			),
		},
		"can have positional and named arguments": {
			input: "Foo(.1, 'foo', :bar, bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(41, 1, 42)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(41, 1, 42)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(41, 1, 42)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 1, 9), P(12, 1, 13)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 1, 16), P(18, 1, 19)), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(21, 1, 22), P(29, 1, 30)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(26, 1, 27), P(29, 1, 30)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(32, 1, 33), P(40, 1, 41)),
									"elk",
									ast.NewTrueLiteralNode(S(P(37, 1, 38), P(40, 1, 41))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after commas": {
			input: "Foo(.1,\n'foo',\n:bar, bar: :baz,\nelk: true)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(41, 4, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(41, 4, 10)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(41, 4, 10)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 2, 1), P(12, 2, 5)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 3, 1), P(18, 3, 4)), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(21, 3, 7), P(29, 3, 15)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(26, 3, 12), P(29, 3, 15)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(32, 4, 1), P(40, 4, 9)),
									"elk",
									ast.NewTrueLiteralNode(S(P(37, 4, 6), P(40, 4, 9))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines around parentheses": {
			input: "Foo(\n.1, 'foo', :bar, bar: :baz, elk: true\n)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(43, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(43, 3, 1)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(43, 3, 1)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(5, 2, 1), P(6, 2, 2)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(9, 2, 5), P(13, 2, 9)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(16, 2, 12), P(19, 2, 15)), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(22, 2, 18), P(30, 2, 26)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(27, 2, 23), P(30, 2, 26)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(33, 2, 29), P(41, 2, 37)),
									"elk",
									ast.NewTrueLiteralNode(S(P(38, 2, 34), P(41, 2, 37))),
								),
							},
						),
					),
				},
			),
		},
		"cannot have newlines before the opening parenthesis": {
			input: "Foo\n(.1, 'foo', :bar, bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(42, 2, 39)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
					),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 2), P(42, 2, 39)),
						ast.NewFloatLiteralNode(S(P(5, 2, 2), P(6, 2, 3)), "0.1"),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(7, 2, 4), P(7, 2, 4)), "unexpected ,, expected )"),
			},
		},
		"can have positional arguments without parentheses": {
			input: "Foo .1, 'foo', :bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 1, 9), P(12, 1, 13)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 1, 16), P(18, 1, 19)), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have named arguments without parentheses": {
			input: "Foo bar: :baz, elk: true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(23, 1, 24)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(23, 1, 24)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(4, 1, 5), P(12, 1, 13)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(9, 1, 10), P(12, 1, 13)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(15, 1, 16), P(23, 1, 24)),
									"elk",
									ast.NewTrueLiteralNode(S(P(20, 1, 21), P(23, 1, 24))),
								),
							},
						),
					),
				},
			),
		},
		"can have positional and named arguments without parentheses": {
			input: "Foo .1, 'foo', :bar, bar: :baz, elk: true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(40, 1, 41)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 1, 41)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(40, 1, 41)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 1, 9), P(12, 1, 13)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 1, 16), P(18, 1, 19)), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(21, 1, 22), P(29, 1, 30)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(26, 1, 27), P(29, 1, 30)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(32, 1, 33), P(40, 1, 41)),
									"elk",
									ast.NewTrueLiteralNode(S(P(37, 1, 38), P(40, 1, 41))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after commas without parentheses": {
			input: "Foo .1,\n'foo',\n:bar, bar: :baz,\nelk: true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(40, 4, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 4, 9)),
						ast.NewConstructorCallNode(
							S(P(0, 1, 1), P(40, 4, 9)),
							ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 2, 1), P(12, 2, 5)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 3, 1), P(18, 3, 4)), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(21, 3, 7), P(29, 3, 15)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(26, 3, 12), P(29, 3, 15)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(32, 4, 1), P(40, 4, 9)),
									"elk",
									ast.NewTrueLiteralNode(S(P(37, 4, 6), P(40, 4, 9))),
								),
							},
						),
					),
				},
			),
		},
		"cannot have newlines before the arguments without parentheses": {
			input: "Foo\n.1, 'foo', :bar, bar: :baz, elk: true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 2, 2)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewPublicConstantNode(S(P(0, 1, 1), P(2, 1, 3)), "Foo"),
					),
					ast.NewExpressionStatementNode(
						S(P(4, 2, 1), P(5, 2, 2)),
						ast.NewFloatLiteralNode(S(P(4, 2, 1), P(5, 2, 2)), "0.1"),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(6, 2, 3), P(6, 2, 3)), "unexpected ,, expected a statement separator `\\n`, `;`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestMethodCall(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "foo = bar()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewFunctionCallNode(
								S(P(6, 1, 7), P(10, 1, 11)),
								"bar",
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can omit the receiver and have an empty argument list": {
			input: "foo()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewFunctionCallNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							"foo",
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a private name with an implicit receiver": {
			input: "_foo()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewFunctionCallNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							"_foo",
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver": {
			input: "foo.bar()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							false,
							"bar",
							nil,
							nil,
						),
					),
				},
			),
		},
		"can omit parentheses": {
			input: "foo.bar 1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							false,
							"bar",
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "1"),
							},
							nil,
						),
					),
				},
			),
		},
		"can use the safe navigation operator": {
			input: "foo?.bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							true,
							"bar",
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be nested with parentheses": {
			input: "foo.bar().baz()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewMethodCallNode(
								S(P(0, 1, 1), P(8, 1, 9)),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								false,
								"bar",
								nil,
								nil,
							),
							false,
							"baz",
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have any expression as the receiver": {
			input: "(foo + 2).bar()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(1, 1, 2), P(14, 1, 15)),
						ast.NewMethodCallNode(
							S(P(1, 1, 2), P(14, 1, 15)),
							ast.NewBinaryExpressionNode(
								S(P(1, 1, 2), P(7, 1, 8)),
								T(S(P(5, 1, 6), P(5, 1, 6)), token.PLUS),
								ast.NewPublicIdentifierNode(S(P(1, 1, 2), P(3, 1, 4)), "foo"),
								ast.NewIntLiteralNode(S(P(7, 1, 8), P(7, 1, 8)), "2"),
							),
							false,
							"bar",
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot call a private method on an explicit receiver": {
			input: "foo._bar()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							false,
							"_bar",
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(7, 1, 8)), "unexpected PRIVATE_IDENTIFIER, expected a public method name (public identifier, keyword or overridable operator)"),
			},
		},
		"can have an overridable operator as the method name with an explicit receiver": {
			input: "foo.+()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							false,
							"+",
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have a non overridable operator as the method name with an explicit receiver": {
			input: "foo.&&()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							false,
							"&&",
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(5, 1, 6)), "unexpected &&, expected a public method name (public identifier, keyword or overridable operator)"),
			},
		},
		"can call a private method on self": {
			input: "self._foo()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewMethodCallNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							ast.NewSelfLiteralNode(S(P(0, 1, 1), P(3, 1, 4))),
							false,
							"_foo",
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have positional arguments": {
			input: "foo(.1, 'foo', :bar)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewFunctionCallNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 1, 9), P(12, 1, 13)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 1, 16), P(18, 1, 19)), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have named arguments": {
			input: "foo(bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewFunctionCallNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							"foo",
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(4, 1, 5), P(12, 1, 13)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(9, 1, 10), P(12, 1, 13)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(15, 1, 16), P(23, 1, 24)),
									"elk",
									ast.NewTrueLiteralNode(S(P(20, 1, 21), P(23, 1, 24))),
								),
							},
						),
					),
				},
			),
		},
		"can have positional and named arguments": {
			input: "foo(.1, 'foo', :bar, bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(41, 1, 42)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(41, 1, 42)),
						ast.NewFunctionCallNode(
							S(P(0, 1, 1), P(41, 1, 42)),
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 1, 9), P(12, 1, 13)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 1, 16), P(18, 1, 19)), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(21, 1, 22), P(29, 1, 30)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(26, 1, 27), P(29, 1, 30)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(32, 1, 33), P(40, 1, 41)),
									"elk",
									ast.NewTrueLiteralNode(S(P(37, 1, 38), P(40, 1, 41))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after commas": {
			input: "foo(.1,\n'foo',\n:bar, bar: :baz,\nelk: true)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(41, 4, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(41, 4, 10)),
						ast.NewFunctionCallNode(
							S(P(0, 1, 1), P(41, 4, 10)),
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 2, 1), P(12, 2, 5)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 3, 1), P(18, 3, 4)), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(21, 3, 7), P(29, 3, 15)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(26, 3, 12), P(29, 3, 15)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(32, 4, 1), P(40, 4, 9)),
									"elk",
									ast.NewTrueLiteralNode(S(P(37, 4, 6), P(40, 4, 9))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines around parentheses": {
			input: "foo(\n.1, 'foo', :bar, bar: :baz, elk: true\n)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(43, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(43, 3, 1)),
						ast.NewFunctionCallNode(
							S(P(0, 1, 1), P(43, 3, 1)),
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(5, 2, 1), P(6, 2, 2)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(9, 2, 5), P(13, 2, 9)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(16, 2, 12), P(19, 2, 15)), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(22, 2, 18), P(30, 2, 26)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(27, 2, 23), P(30, 2, 26)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(33, 2, 29), P(41, 2, 37)),
									"elk",
									ast.NewTrueLiteralNode(S(P(38, 2, 34), P(41, 2, 37))),
								),
							},
						),
					),
				},
			),
		},
		"cannot have newlines before the opening parenthesis": {
			input: "foo\n(.1, 'foo', :bar, bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(42, 2, 39)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
					),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 2), P(42, 2, 39)),
						ast.NewFloatLiteralNode(S(P(5, 2, 2), P(6, 2, 3)), "0.1"),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(7, 2, 4), P(7, 2, 4)), "unexpected ,, expected )"),
			},
		},
		"can have positional arguments without parentheses": {
			input: "foo .1, 'foo', :bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewFunctionCallNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 1, 9), P(12, 1, 13)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 1, 16), P(18, 1, 19)), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have named arguments without parentheses": {
			input: "foo bar: :baz, elk: true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(23, 1, 24)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewFunctionCallNode(
							S(P(0, 1, 1), P(23, 1, 24)),
							"foo",
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(4, 1, 5), P(12, 1, 13)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(9, 1, 10), P(12, 1, 13)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(15, 1, 16), P(23, 1, 24)),
									"elk",
									ast.NewTrueLiteralNode(S(P(20, 1, 21), P(23, 1, 24))),
								),
							},
						),
					),
				},
			),
		},
		"can have positional and named arguments without parentheses": {
			input: "foo .1, 'foo', :bar, bar: :baz, elk: true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(40, 1, 41)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 1, 41)),
						ast.NewFunctionCallNode(
							S(P(0, 1, 1), P(40, 1, 41)),
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 1, 9), P(12, 1, 13)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 1, 16), P(18, 1, 19)), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(21, 1, 22), P(29, 1, 30)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(26, 1, 27), P(29, 1, 30)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(32, 1, 33), P(40, 1, 41)),
									"elk",
									ast.NewTrueLiteralNode(S(P(37, 1, 38), P(40, 1, 41))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after commas without parentheses": {
			input: "foo .1,\n'foo',\n:bar, bar: :baz,\nelk: true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(40, 4, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 4, 9)),
						ast.NewFunctionCallNode(
							S(P(0, 1, 1), P(40, 4, 9)),
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 2, 1), P(12, 2, 5)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 3, 1), P(18, 3, 4)), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									S(P(21, 3, 7), P(29, 3, 15)),
									"bar",
									ast.NewSimpleSymbolLiteralNode(S(P(26, 3, 12), P(29, 3, 15)), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									S(P(32, 4, 1), P(40, 4, 9)),
									"elk",
									ast.NewTrueLiteralNode(S(P(37, 4, 6), P(40, 4, 9))),
								),
							},
						),
					),
				},
			),
		},
		"cannot have newlines before the arguments without parentheses": {
			input: "foo\n.1, 'foo', :bar, bar: :baz, elk: true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 2, 2)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
					),
					ast.NewExpressionStatementNode(
						S(P(4, 2, 1), P(5, 2, 2)),
						ast.NewFloatLiteralNode(S(P(4, 2, 1), P(5, 2, 2)), "0.1"),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(6, 2, 3), P(6, 2, 3)), "unexpected ,, expected a statement separator `\\n`, `;`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestAttributeAccess(t *testing.T) {
	tests := testTable{
		"can be used on self": {
			input: "self.bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewAttributeAccessNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							ast.NewSelfLiteralNode(S(P(0, 1, 1), P(3, 1, 4))),
							"bar",
						),
					),
				},
			),
		},
		"can be called on variables": {
			input: "foo.bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewAttributeAccessNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							"bar",
						),
					),
				},
			),
		},
		"can be nested": {
			input: "foo.bar.baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewAttributeAccessNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							ast.NewAttributeAccessNode(
								S(P(0, 1, 1), P(6, 1, 7)),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								"bar",
							),
							"baz",
						),
					),
				},
			),
		},
		"can have newlines after the dot": {
			input: "foo.\nbar.\nbaz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 3, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 3, 3)),
						ast.NewAttributeAccessNode(
							S(P(0, 1, 1), P(12, 3, 3)),
							ast.NewAttributeAccessNode(
								S(P(0, 1, 1), P(7, 2, 3)),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								"bar",
							),
							"baz",
						),
					),
				},
			),
		},
		"can have newlines before the dot": {
			input: "foo\n.bar\n.baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 3, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 3, 4)),
						ast.NewAttributeAccessNode(
							S(P(0, 1, 1), P(12, 3, 4)),
							ast.NewAttributeAccessNode(
								S(P(0, 1, 1), P(7, 2, 4)),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								"bar",
							),
							"baz",
						),
					),
				},
			),
		},
		"can be nested on function calls": {
			input: "foo().bar.baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewAttributeAccessNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							ast.NewAttributeAccessNode(
								S(P(0, 1, 1), P(8, 1, 9)),
								ast.NewFunctionCallNode(
									S(P(0, 1, 1), P(4, 1, 5)),
									"foo",
									nil,
									nil,
								),
								"bar",
							),
							"baz",
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}
