package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
)

func TestConstructorCall(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "foo = Foo()",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 11, 1, 1),
							T(P(4, 1, 1, 5), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewConstructorCallNode(
								P(6, 5, 1, 7),
								ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
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
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 5, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
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
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 6, 1, 1),
							ast.NewPrivateConstantNode(P(0, 4, 1, 1), "_Foo"),
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
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 10, 1, 1),
							ast.NewConstantLookupNode(
								P(0, 8, 1, 1),
								ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
								ast.NewPublicConstantNode(P(5, 3, 1, 6), "Bar"),
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
				P(0, 20, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 20, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 20, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 1, 9), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 1, 16), "bar"),
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
				P(0, 25, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 25, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 25, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(4, 9, 1, 5),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(9, 4, 1, 10), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(15, 9, 1, 16),
									"elk",
									ast.NewTrueLiteralNode(P(20, 4, 1, 21)),
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
				P(0, 42, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 42, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 42, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 1, 9), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 1, 16), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(21, 9, 1, 22),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(26, 4, 1, 27), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(32, 9, 1, 33),
									"elk",
									ast.NewTrueLiteralNode(P(37, 4, 1, 38)),
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
				P(0, 42, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 42, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 42, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 2, 1), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 3, 1), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(21, 9, 3, 7),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(26, 4, 3, 12), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(32, 9, 4, 1),
									"elk",
									ast.NewTrueLiteralNode(P(37, 4, 4, 6)),
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
				P(0, 44, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 44, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 44, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(5, 2, 2, 1), "0.1"),
								ast.NewRawStringLiteralNode(P(9, 5, 2, 5), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(16, 4, 2, 12), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(22, 9, 2, 18),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(27, 4, 2, 23), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(33, 9, 2, 29),
									"elk",
									ast.NewTrueLiteralNode(P(38, 4, 2, 34)),
								),
							},
						),
					),
				},
			),
		},
		"can't have newlines before the opening parenthesis": {
			input: "Foo\n(.1, 'foo', :bar, bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				P(0, 43, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
					),
					ast.NewExpressionStatementNode(
						P(5, 38, 2, 2),
						ast.NewFloatLiteralNode(P(5, 2, 2, 2), "0.1"),
					),
				},
			),
			err: ErrorList{
				NewError(P(7, 1, 2, 4), "unexpected ,, expected )"),
			},
		},
		"can have positional arguments without parentheses": {
			input: "Foo .1, 'foo', :bar",
			want: ast.NewProgramNode(
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 19, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 19, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 1, 9), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 1, 16), "bar"),
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
				P(0, 24, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 24, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 24, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(4, 9, 1, 5),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(9, 4, 1, 10), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(15, 9, 1, 16),
									"elk",
									ast.NewTrueLiteralNode(P(20, 4, 1, 21)),
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
				P(0, 41, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 41, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 41, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 1, 9), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 1, 16), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(21, 9, 1, 22),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(26, 4, 1, 27), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(32, 9, 1, 33),
									"elk",
									ast.NewTrueLiteralNode(P(37, 4, 1, 38)),
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
				P(0, 41, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 41, 1, 1),
						ast.NewConstructorCallNode(
							P(0, 41, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 2, 1), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 3, 1), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(21, 9, 3, 7),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(26, 4, 3, 12), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(32, 9, 4, 1),
									"elk",
									ast.NewTrueLiteralNode(P(37, 4, 4, 6)),
								),
							},
						),
					),
				},
			),
		},
		"can't have newlines before the arguments without parentheses": {
			input: "Foo\n.1, 'foo', :bar, bar: :baz, elk: true",
			want: ast.NewProgramNode(
				P(0, 41, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
					),
					ast.NewExpressionStatementNode(
						P(4, 2, 2, 1),
						ast.NewFloatLiteralNode(P(4, 2, 2, 1), "0.1"),
					),
				},
			),
			err: ErrorList{
				NewError(P(6, 1, 2, 3), "unexpected ,, expected a statement separator `\\n`, `;`"),
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
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 11, 1, 1),
							T(P(4, 1, 1, 5), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewMethodCallNode(
								P(6, 5, 1, 7),
								nil,
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
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewMethodCallNode(
							P(0, 5, 1, 1),
							nil,
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
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewMethodCallNode(
							P(0, 6, 1, 1),
							nil,
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
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewMethodCallNode(
							P(0, 9, 1, 1),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							"bar",
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
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(1, 14, 1, 2),
						ast.NewMethodCallNode(
							P(1, 14, 1, 2),
							ast.NewBinaryExpressionNode(
								P(1, 7, 1, 2),
								T(P(5, 1, 1, 6), token.PLUS),
								ast.NewPublicIdentifierNode(P(1, 3, 1, 2), "foo"),
								ast.NewIntLiteralNode(P(7, 1, 1, 8), V(P(7, 1, 1, 8), token.DEC_INT, "2")),
							),
							"bar",
							nil,
							nil,
						),
					),
				},
			),
		},
		"can't call a private method on an explicit receiver": {
			input: "foo._bar()",
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewMethodCallNode(
							P(0, 10, 1, 1),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							"_bar",
							nil,
							nil,
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(4, 4, 1, 5), "unexpected PRIVATE_IDENTIFIER, expected a public method name (public identifier, keyword or overridable operator)"),
			},
		},
		"can have an overridable operator as the method name with an explicit receiver": {
			input: "foo.+()",
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewMethodCallNode(
							P(0, 7, 1, 1),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							"+",
							nil,
							nil,
						),
					),
				},
			),
		},
		"can't have a non overridable operator as the method name with an explicit receiver": {
			input: "foo.&&()",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewMethodCallNode(
							P(0, 8, 1, 1),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							"&&",
							nil,
							nil,
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(4, 2, 1, 5), "unexpected &&, expected a public method name (public identifier, keyword or overridable operator)"),
			},
		},
		"can call a private method on self": {
			input: "self._foo()",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewMethodCallNode(
							P(0, 11, 1, 1),
							ast.NewSelfLiteralNode(P(0, 4, 1, 1)),
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
				P(0, 20, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 20, 1, 1),
						ast.NewMethodCallNode(
							P(0, 20, 1, 1),
							nil,
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 1, 9), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 1, 16), "bar"),
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
				P(0, 25, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 25, 1, 1),
						ast.NewMethodCallNode(
							P(0, 25, 1, 1),
							nil,
							"foo",
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(4, 9, 1, 5),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(9, 4, 1, 10), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(15, 9, 1, 16),
									"elk",
									ast.NewTrueLiteralNode(P(20, 4, 1, 21)),
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
				P(0, 42, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 42, 1, 1),
						ast.NewMethodCallNode(
							P(0, 42, 1, 1),
							nil,
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 1, 9), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 1, 16), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(21, 9, 1, 22),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(26, 4, 1, 27), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(32, 9, 1, 33),
									"elk",
									ast.NewTrueLiteralNode(P(37, 4, 1, 38)),
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
				P(0, 42, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 42, 1, 1),
						ast.NewMethodCallNode(
							P(0, 42, 1, 1),
							nil,
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 2, 1), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 3, 1), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(21, 9, 3, 7),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(26, 4, 3, 12), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(32, 9, 4, 1),
									"elk",
									ast.NewTrueLiteralNode(P(37, 4, 4, 6)),
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
				P(0, 44, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 44, 1, 1),
						ast.NewMethodCallNode(
							P(0, 44, 1, 1),
							nil,
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(5, 2, 2, 1), "0.1"),
								ast.NewRawStringLiteralNode(P(9, 5, 2, 5), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(16, 4, 2, 12), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(22, 9, 2, 18),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(27, 4, 2, 23), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(33, 9, 2, 29),
									"elk",
									ast.NewTrueLiteralNode(P(38, 4, 2, 34)),
								),
							},
						),
					),
				},
			),
		},
		"can't have newlines before the opening parenthesis": {
			input: "foo\n(.1, 'foo', :bar, bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				P(0, 43, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
					),
					ast.NewExpressionStatementNode(
						P(5, 38, 2, 2),
						ast.NewFloatLiteralNode(P(5, 2, 2, 2), "0.1"),
					),
				},
			),
			err: ErrorList{
				NewError(P(7, 1, 2, 4), "unexpected ,, expected )"),
			},
		},
		"can have positional arguments without parentheses": {
			input: "foo .1, 'foo', :bar",
			want: ast.NewProgramNode(
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 19, 1, 1),
						ast.NewMethodCallNode(
							P(0, 19, 1, 1),
							nil,
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 1, 9), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 1, 16), "bar"),
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
				P(0, 24, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 24, 1, 1),
						ast.NewMethodCallNode(
							P(0, 24, 1, 1),
							nil,
							"foo",
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(4, 9, 1, 5),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(9, 4, 1, 10), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(15, 9, 1, 16),
									"elk",
									ast.NewTrueLiteralNode(P(20, 4, 1, 21)),
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
				P(0, 41, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 41, 1, 1),
						ast.NewMethodCallNode(
							P(0, 41, 1, 1),
							nil,
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 1, 9), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 1, 16), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(21, 9, 1, 22),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(26, 4, 1, 27), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(32, 9, 1, 33),
									"elk",
									ast.NewTrueLiteralNode(P(37, 4, 1, 38)),
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
				P(0, 41, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 41, 1, 1),
						ast.NewMethodCallNode(
							P(0, 41, 1, 1),
							nil,
							"foo",
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 2, 1), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 3, 1), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									P(21, 9, 3, 7),
									"bar",
									ast.NewSimpleSymbolLiteralNode(P(26, 4, 3, 12), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									P(32, 9, 4, 1),
									"elk",
									ast.NewTrueLiteralNode(P(37, 4, 4, 6)),
								),
							},
						),
					),
				},
			),
		},
		"can't have newlines before the arguments without parentheses": {
			input: "foo\n.1, 'foo', :bar, bar: :baz, elk: true",
			want: ast.NewProgramNode(
				P(0, 41, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
					),
					ast.NewExpressionStatementNode(
						P(4, 2, 2, 1),
						ast.NewFloatLiteralNode(P(4, 2, 2, 1), "0.1"),
					),
				},
			),
			err: ErrorList{
				NewError(P(6, 1, 2, 3), "unexpected ,, expected a statement separator `\\n`, `;`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}
