package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/token"
)

func TestPipeExpression(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "foo = 2 |> Foo()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewAssignmentExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							T(L(S(P(4, 1, 5), P(4, 1, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewBinaryExpressionNode(
								L(S(P(6, 1, 7), P(15, 1, 16))),
								T(L(S(P(8, 1, 9), P(9, 1, 10))), token.PIPE_OP),
								ast.NewIntLiteralNode(L(S(P(6, 1, 7), P(6, 1, 7))), "2"),
								ast.NewConstructorCallNode(
									L(S(P(11, 1, 12), P(15, 1, 16))),
									ast.NewPublicConstantNode(L(S(P(11, 1, 12), P(13, 1, 14))), "Foo"),
									nil,
									nil,
								),
							),
						),
					),
				},
			),
		},
		"can have a constructor call on the right hand side": {
			input: "2 |> Foo()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							T(L(S(P(2, 1, 3), P(3, 1, 4))), token.PIPE_OP),
							ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "2"),
							ast.NewConstructorCallNode(
								L(S(P(5, 1, 6), P(9, 1, 10))),
								ast.NewPublicConstantNode(L(S(P(5, 1, 6), P(7, 1, 8))), "Foo"),
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have a receiverless method call on the right hand side": {
			input: "2 |> foo()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							T(L(S(P(2, 1, 3), P(3, 1, 4))), token.PIPE_OP),
							ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "2"),
							ast.NewReceiverlessMethodCallNode(
								L(S(P(5, 1, 6), P(9, 1, 10))),
								ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "foo"),
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have a method call on the right hand side": {
			input: "2 |> a.foo()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							T(L(S(P(2, 1, 3), P(3, 1, 4))), token.PIPE_OP),
							ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "2"),
							ast.NewMethodCallNode(
								L(S(P(5, 1, 6), P(11, 1, 12))),
								ast.NewPublicIdentifierNode(
									L(S(P(5, 1, 6), P(5, 1, 6))),
									"a",
								),
								T(L(S(P(6, 1, 7), P(6, 1, 7))), token.DOT),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have a function call on the right hand side": {
			input: "2 |> a.()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(2, 1, 3), P(3, 1, 4))), token.PIPE_OP),
							ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "2"),
							ast.NewCallNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								ast.NewPublicIdentifierNode(
									L(S(P(5, 1, 6), P(5, 1, 6))),
									"a",
								),
								false,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"cannot have non method calls on the right": {
			input: "2 |> foo() + 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(13, 1, 14))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(13, 1, 14))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(13, 1, 14))),
							T(L(S(P(2, 1, 3), P(3, 1, 4))), token.PIPE_OP),
							ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "2"),
							ast.NewBinaryExpressionNode(
								L(S(P(5, 1, 6), P(13, 1, 14))),
								T(L(S(P(11, 1, 12), P(11, 1, 12))), token.PLUS),
								ast.NewReceiverlessMethodCallNode(
									L(S(P(5, 1, 6), P(9, 1, 10))),
									ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "foo"),
									nil,
									nil,
								),
								ast.NewIntLiteralNode(L(S(P(13, 1, 14), P(13, 1, 14))), "2"),
							),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(5, 1, 6), P(13, 1, 14))), "invalid right hand side of a pipe expression, only method and function calls are allowed"),
			},
		},
		"can be chained": {
			input: "2 |> foo() |> bar(9.2)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 1, 22))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 1, 22))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(21, 1, 22))),
							T(L(S(P(11, 1, 12), P(12, 1, 13))), token.PIPE_OP),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								T(L(S(P(2, 1, 3), P(3, 1, 4))), token.PIPE_OP),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "2"),
								ast.NewReceiverlessMethodCallNode(
									L(S(P(5, 1, 6), P(9, 1, 10))),
									ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "foo"),
									nil,
									nil,
								),
							),
							ast.NewReceiverlessMethodCallNode(
								L(S(P(14, 1, 15), P(21, 1, 22))),
								ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "bar"),
								[]ast.ExpressionNode{
									ast.NewFloatLiteralNode(L(S(P(18, 1, 19), P(20, 1, 21))), "9.2"),
								},
								nil,
							),
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

func TestConstructorCall(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "foo = Foo()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewAssignmentExpressionNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							T(L(S(P(4, 1, 5), P(4, 1, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewConstructorCallNode(
								L(S(P(6, 1, 7), P(10, 1, 11))),
								ast.NewPublicConstantNode(L(S(P(6, 1, 7), P(8, 1, 9))), "Foo"),
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
				L(S(P(0, 1, 1), P(4, 1, 5))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(4, 1, 5))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(4, 1, 5))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have type arguments": {
			input: "Foo::[1 | 2, String]()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 1, 22))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 1, 22))),
						ast.NewGenericConstructorCallNode(
							L(S(P(0, 1, 1), P(21, 1, 22))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							[]ast.TypeNode{
								ast.NewBinaryTypeNode(
									L(S(P(6, 1, 7), P(10, 1, 11))),
									T(L(S(P(8, 1, 9), P(8, 1, 9))), token.OR),
									ast.NewIntLiteralNode(L(S(P(6, 1, 7), P(6, 1, 7))), "1"),
									ast.NewIntLiteralNode(L(S(P(10, 1, 11), P(10, 1, 11))), "2"),
								),
								ast.NewPublicConstantNode(L(S(P(13, 1, 14), P(18, 1, 19))), "String"),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have type arguments and regular arguments": {
			input: "Foo::[1 | 2, String](8)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(22, 1, 23))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(22, 1, 23))),
						ast.NewGenericConstructorCallNode(
							L(S(P(0, 1, 1), P(22, 1, 23))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							[]ast.TypeNode{
								ast.NewBinaryTypeNode(
									L(S(P(6, 1, 7), P(10, 1, 11))),
									T(L(S(P(8, 1, 9), P(8, 1, 9))), token.OR),
									ast.NewIntLiteralNode(L(S(P(6, 1, 7), P(6, 1, 7))), "1"),
									ast.NewIntLiteralNode(L(S(P(10, 1, 11), P(10, 1, 11))), "2"),
								),
								ast.NewPublicConstantNode(L(S(P(13, 1, 14), P(18, 1, 19))), "String"),
							},
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "8"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have type arguments and regular arguments without parens": {
			input: "Foo::[1 | 2, String] 8",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 1, 22))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 1, 22))),
						ast.NewGenericConstructorCallNode(
							L(S(P(0, 1, 1), P(21, 1, 22))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							[]ast.TypeNode{
								ast.NewBinaryTypeNode(
									L(S(P(6, 1, 7), P(10, 1, 11))),
									T(L(S(P(8, 1, 9), P(8, 1, 9))), token.OR),
									ast.NewIntLiteralNode(L(S(P(6, 1, 7), P(6, 1, 7))), "1"),
									ast.NewIntLiteralNode(L(S(P(10, 1, 11), P(10, 1, 11))), "2"),
								),
								ast.NewPublicConstantNode(L(S(P(13, 1, 14), P(18, 1, 19))), "String"),
							},
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "8"),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot have type arguments without regular arguments": {
			input: "Foo::[1 | 2, String]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(19, 1, 20))),
						ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(0, 1, 1), P(19, 1, 20))), "invalid generic constant"),
			},
		},
		"can have multiline type arguments": {
			input: `
				Foo::[
					1 | 2,
					String,
				]()
			`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(44, 5, 8))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(
						L(S(P(0, 1, 1), P(0, 1, 1))),
					),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 5), P(44, 5, 8))),
						ast.NewGenericConstructorCallNode(
							L(S(P(5, 2, 5), P(43, 5, 7))),
							ast.NewPublicConstantNode(L(S(P(5, 2, 5), P(7, 2, 7))), "Foo"),
							[]ast.TypeNode{
								ast.NewBinaryTypeNode(
									L(S(P(17, 3, 6), P(21, 3, 10))),
									T(L(S(P(19, 3, 8), P(19, 3, 8))), token.OR),
									ast.NewIntLiteralNode(L(S(P(17, 3, 6), P(17, 3, 6))), "1"),
									ast.NewIntLiteralNode(L(S(P(21, 3, 10), P(21, 3, 10))), "2"),
								),
								ast.NewPublicConstantNode(L(S(P(29, 4, 6), P(34, 4, 11))), "String"),
							},
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
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(5, 1, 6))),
							ast.NewPrivateConstantNode(L(S(P(0, 1, 1), P(3, 1, 4))), "_Foo"),
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
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewConstantLookupNode(
								L(S(P(0, 1, 1), P(7, 1, 8))),
								ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
								ast.NewPublicConstantNode(L(S(P(5, 1, 6), P(7, 1, 8))), "Bar"),
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
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(19, 1, 20))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(19, 1, 20))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have a trailing comma": {
			input: "Foo(.1, 'foo', :bar,)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(20, 1, 21))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(20, 1, 21))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(20, 1, 21))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
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
				L(S(P(0, 1, 1), P(24, 1, 25))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(24, 1, 25))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(24, 1, 25))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(4, 1, 5), P(12, 1, 13))),
									ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(9, 1, 10), P(12, 1, 13))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(15, 1, 16), P(23, 1, 24))),
									ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(17, 1, 18))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(20, 1, 21), P(23, 1, 24)))),
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
				L(S(P(0, 1, 1), P(41, 1, 42))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(41, 1, 42))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(41, 1, 42))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 1, 22), P(29, 1, 30))),
									ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(23, 1, 24))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 1, 27), P(29, 1, 30))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 1, 33), P(40, 1, 41))),
									ast.NewPublicIdentifierNode(L(S(P(32, 1, 33), P(34, 1, 35))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 1, 38), P(40, 1, 41)))),
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
				L(S(P(0, 1, 1), P(41, 4, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(41, 4, 10))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(41, 4, 10))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 2, 1), P(12, 2, 5))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 3, 1), P(18, 3, 4))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 3, 7), P(29, 3, 15))),
									ast.NewPublicIdentifierNode(L(S(P(21, 3, 7), P(23, 3, 9))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 3, 12), P(29, 3, 15))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 4, 1), P(40, 4, 9))),
									ast.NewPublicIdentifierNode(L(S(P(32, 4, 1), P(34, 4, 3))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 4, 6), P(40, 4, 9)))),
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
				L(S(P(0, 1, 1), P(43, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(43, 3, 1))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(43, 3, 1))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(5, 2, 1), P(6, 2, 2))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(9, 2, 5), P(13, 2, 9))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(16, 2, 12), P(19, 2, 15))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(22, 2, 18), P(30, 2, 26))),
									ast.NewPublicIdentifierNode(L(S(P(22, 2, 18), P(24, 2, 20))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(27, 2, 23), P(30, 2, 26))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(33, 2, 29), P(41, 2, 37))),
									ast.NewPublicIdentifierNode(L(S(P(33, 2, 29), P(35, 2, 31))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(38, 2, 34), P(41, 2, 37)))),
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
				L(S(P(0, 1, 1), P(42, 2, 39))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(3, 1, 4))),
						ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
					),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 2), P(42, 2, 39))),
						ast.NewFloatLiteralNode(L(S(P(5, 2, 2), P(6, 2, 3))), "0.1"),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(7, 2, 4), P(7, 2, 4))), "unexpected ,, expected )"),
			},
		},
		"can have positional arguments without parentheses": {
			input: "Foo .1, 'foo', :bar",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(18, 1, 19))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(18, 1, 19))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(18, 1, 19))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
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
				L(S(P(0, 1, 1), P(23, 1, 24))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(23, 1, 24))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(23, 1, 24))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(4, 1, 5), P(12, 1, 13))),
									ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(9, 1, 10), P(12, 1, 13))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(15, 1, 16), P(23, 1, 24))),
									ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(17, 1, 18))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(20, 1, 21), P(23, 1, 24)))),
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
				L(S(P(0, 1, 1), P(40, 1, 41))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(40, 1, 41))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(40, 1, 41))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 1, 22), P(29, 1, 30))),
									ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(23, 1, 24))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 1, 27), P(29, 1, 30))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 1, 33), P(40, 1, 41))),
									ast.NewPublicIdentifierNode(L(S(P(32, 1, 33), P(34, 1, 35))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 1, 38), P(40, 1, 41)))),
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
				L(S(P(0, 1, 1), P(40, 4, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(40, 4, 9))),
						ast.NewConstructorCallNode(
							L(S(P(0, 1, 1), P(40, 4, 9))),
							ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 2, 1), P(12, 2, 5))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 3, 1), P(18, 3, 4))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 3, 7), P(29, 3, 15))),
									ast.NewPublicIdentifierNode(L(S(P(21, 3, 7), P(23, 3, 9))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 3, 12), P(29, 3, 15))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 4, 1), P(40, 4, 9))),
									ast.NewPublicIdentifierNode(L(S(P(32, 4, 1), P(34, 4, 3))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 4, 6), P(40, 4, 9)))),
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
				L(S(P(0, 1, 1), P(5, 2, 2))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(3, 1, 4))),
						ast.NewPublicConstantNode(L(S(P(0, 1, 1), P(2, 1, 3))), "Foo"),
					),
					ast.NewExpressionStatementNode(
						L(S(P(4, 2, 1), P(5, 2, 2))),
						ast.NewFloatLiteralNode(L(S(P(4, 2, 1), P(5, 2, 2))), "0.1"),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 2, 3), P(6, 2, 3))), "unexpected ,, expected a statement separator `\\n`, `;`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestNewExpression(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "foo = new",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewAssignmentExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(4, 1, 5), P(4, 1, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewNewExpressionNode(
								L(S(P(6, 1, 7), P(8, 1, 9))),
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have an empty argument list": {
			input: "new()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(4, 1, 5))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(4, 1, 5))),
						ast.NewNewExpressionNode(
							L(S(P(0, 1, 1), P(4, 1, 5))),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have positional arguments": {
			input: "new(.1, 'foo', :bar)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(19, 1, 20))),
						ast.NewNewExpressionNode(
							L(S(P(0, 1, 1), P(19, 1, 20))),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have a trailing comma": {
			input: "new(.1, 'foo', :bar,)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(20, 1, 21))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(20, 1, 21))),
						ast.NewNewExpressionNode(
							L(S(P(0, 1, 1), P(20, 1, 21))),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have named arguments": {
			input: "new(bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(24, 1, 25))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(24, 1, 25))),
						ast.NewNewExpressionNode(
							L(S(P(0, 1, 1), P(24, 1, 25))),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(4, 1, 5), P(12, 1, 13))),
									ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(9, 1, 10), P(12, 1, 13))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(15, 1, 16), P(23, 1, 24))),
									ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(17, 1, 18))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(20, 1, 21), P(23, 1, 24)))),
								),
							},
						),
					),
				},
			),
		},
		"can have positional and named arguments": {
			input: "new(.1, 'foo', :bar, bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(41, 1, 42))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(41, 1, 42))),
						ast.NewNewExpressionNode(
							L(S(P(0, 1, 1), P(41, 1, 42))),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 1, 22), P(29, 1, 30))),
									ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(23, 1, 24))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 1, 27), P(29, 1, 30))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 1, 33), P(40, 1, 41))),
									ast.NewPublicIdentifierNode(L(S(P(32, 1, 33), P(34, 1, 35))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 1, 38), P(40, 1, 41)))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after commas": {
			input: "new(.1,\n'foo',\n:bar, bar: :baz,\nelk: true)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(41, 4, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(41, 4, 10))),
						ast.NewNewExpressionNode(
							L(S(P(0, 1, 1), P(41, 4, 10))),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 2, 1), P(12, 2, 5))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 3, 1), P(18, 3, 4))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 3, 7), P(29, 3, 15))),
									ast.NewPublicIdentifierNode(L(S(P(21, 3, 7), P(23, 3, 9))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 3, 12), P(29, 3, 15))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 4, 1), P(40, 4, 9))),
									ast.NewPublicIdentifierNode(L(S(P(32, 4, 1), P(34, 4, 3))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 4, 6), P(40, 4, 9)))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines around parentheses": {
			input: "new(\n.1, 'foo', :bar, bar: :baz, elk: true\n)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(43, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(43, 3, 1))),
						ast.NewNewExpressionNode(
							L(S(P(0, 1, 1), P(43, 3, 1))),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(5, 2, 1), P(6, 2, 2))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(9, 2, 5), P(13, 2, 9))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(16, 2, 12), P(19, 2, 15))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(22, 2, 18), P(30, 2, 26))),
									ast.NewPublicIdentifierNode(L(S(P(22, 2, 18), P(24, 2, 20))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(27, 2, 23), P(30, 2, 26))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(33, 2, 29), P(41, 2, 37))),
									ast.NewPublicIdentifierNode(L(S(P(33, 2, 29), P(35, 2, 31))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(38, 2, 34), P(41, 2, 37)))),
								),
							},
						),
					),
				},
			),
		},
		"can have positional arguments without parentheses": {
			input: "new .1, 'foo', :bar",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(18, 1, 19))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(18, 1, 19))),
						ast.NewNewExpressionNode(
							L(S(P(0, 1, 1), P(18, 1, 19))),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have named arguments without parentheses": {
			input: "new bar: :baz, elk: true",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(23, 1, 24))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(23, 1, 24))),
						ast.NewNewExpressionNode(
							L(S(P(0, 1, 1), P(23, 1, 24))),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(4, 1, 5), P(12, 1, 13))),
									ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(9, 1, 10), P(12, 1, 13))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(15, 1, 16), P(23, 1, 24))),
									ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(17, 1, 18))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(20, 1, 21), P(23, 1, 24)))),
								),
							},
						),
					),
				},
			),
		},
		"can have positional and named arguments without parentheses": {
			input: "new .1, 'foo', :bar, bar: :baz, elk: true",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(40, 1, 41))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(40, 1, 41))),
						ast.NewNewExpressionNode(
							L(S(P(0, 1, 1), P(40, 1, 41))),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 1, 22), P(29, 1, 30))),
									ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(23, 1, 24))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 1, 27), P(29, 1, 30))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 1, 33), P(40, 1, 41))),
									ast.NewPublicIdentifierNode(L(S(P(32, 1, 33), P(34, 1, 35))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 1, 38), P(40, 1, 41)))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after commas without parentheses": {
			input: "new .1,\n'foo',\n:bar, bar: :baz,\nelk: true",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(40, 4, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(40, 4, 9))),
						ast.NewNewExpressionNode(
							L(S(P(0, 1, 1), P(40, 4, 9))),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 2, 1), P(12, 2, 5))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 3, 1), P(18, 3, 4))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 3, 7), P(29, 3, 15))),
									ast.NewPublicIdentifierNode(L(S(P(21, 3, 7), P(23, 3, 9))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 3, 12), P(29, 3, 15))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 4, 1), P(40, 4, 9))),
									ast.NewPublicIdentifierNode(L(S(P(32, 4, 1), P(34, 4, 3))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 4, 6), P(40, 4, 9)))),
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
			parserTest(tc, t)
		})
	}
}

func TestMethodCall(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "foo = bar()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewAssignmentExpressionNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							T(L(S(P(4, 1, 5), P(4, 1, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewReceiverlessMethodCallNode(
								L(S(P(6, 1, 7), P(10, 1, 11))),
								ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
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
				L(S(P(0, 1, 1), P(4, 1, 5))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(4, 1, 5))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(4, 1, 5))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can omit the receiver and have type arguments": {
			input: "foo::[String]()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewGenericReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.TypeNode{
								ast.NewPublicConstantNode(L(S(P(6, 1, 7), P(11, 1, 12))), "String"),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be multiline, omit the receiver and have type arguments": {
			input: `
				foo::[
					String,
				]()
			`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(32, 4, 8))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(
						L(S(P(0, 1, 1), P(0, 1, 1))),
					),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 5), P(32, 4, 8))),
						ast.NewGenericReceiverlessMethodCallNode(
							L(S(P(5, 2, 5), P(31, 4, 7))),
							ast.NewPublicIdentifierNode(L(S(P(5, 2, 5), P(7, 2, 7))), "foo"),
							[]ast.TypeNode{
								ast.NewPublicConstantNode(L(S(P(17, 3, 6), P(22, 3, 11))), "String"),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can omit the receiver, arguments and have a trailing closure": {
			input: "foo() |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(17, 1, 18))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(17, 1, 18))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(17, 1, 18))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(6, 1, 7), P(17, 1, 18))),
									[]ast.ParameterNode{
										ast.NewFormalParameterNode(
											L(S(P(7, 1, 8), P(7, 1, 8))),
											ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(7, 1, 8))), "i"),
											nil,
											nil,
											ast.NormalParameterKind,
										),
									},
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(13, 1, 14), P(17, 1, 18))),
											ast.NewBinaryExpressionNode(
												L(S(P(13, 1, 14), P(17, 1, 18))),
												T(L(S(P(15, 1, 16), P(15, 1, 16))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "i"),
												ast.NewIntLiteralNode(L(S(P(17, 1, 18), P(17, 1, 18))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can omit the receiver, arguments and have a trailing closure without pipes": {
			input: "foo() -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(13, 1, 14))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(13, 1, 14))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(13, 1, 14))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(6, 1, 7), P(13, 1, 14))),
									nil,
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(9, 1, 10), P(13, 1, 14))),
											ast.NewBinaryExpressionNode(
												L(S(P(9, 1, 10), P(13, 1, 14))),
												T(L(S(P(11, 1, 12), P(11, 1, 12))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(9, 1, 10), P(9, 1, 10))), "i"),
												ast.NewIntLiteralNode(L(S(P(13, 1, 14), P(13, 1, 14))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can omit the receiver, arguments and have a trailing closure without arguments": {
			input: "foo() || -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(6, 1, 7), P(16, 1, 17))),
									nil,
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(12, 1, 13), P(16, 1, 17))),
											ast.NewBinaryExpressionNode(
												L(S(P(12, 1, 13), P(16, 1, 17))),
												T(L(S(P(14, 1, 15), P(14, 1, 15))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(12, 1, 13), P(12, 1, 13))), "i"),
												ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can omit the receiver and have a trailing closure": {
			input: "foo(1, 5) |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 1, 22))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 1, 22))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(21, 1, 22))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "1"),
								ast.NewIntLiteralNode(L(S(P(7, 1, 8), P(7, 1, 8))), "5"),
								ast.NewClosureLiteralNode(
									L(S(P(10, 1, 11), P(21, 1, 22))),
									[]ast.ParameterNode{
										ast.NewFormalParameterNode(
											L(S(P(11, 1, 12), P(11, 1, 12))),
											ast.NewPublicIdentifierNode(L(S(P(11, 1, 12), P(11, 1, 12))), "i"),
											nil,
											nil,
											ast.NormalParameterKind,
										),
									},
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(17, 1, 18), P(21, 1, 22))),
											ast.NewBinaryExpressionNode(
												L(S(P(17, 1, 18), P(21, 1, 22))),
												T(L(S(P(19, 1, 20), P(19, 1, 20))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "i"),
												ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can omit the receiver have named arguments and a trailing closure": {
			input: "foo(f: 5) |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 1, 22))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 1, 22))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(21, 1, 22))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(4, 1, 5), P(7, 1, 8))),
									ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(4, 1, 5))), "f"),
									ast.NewIntLiteralNode(L(S(P(7, 1, 8), P(7, 1, 8))), "5"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(10, 1, 11), P(21, 1, 22))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(21, 1, 22))), "func"),
									ast.NewClosureLiteralNode(
										L(S(P(10, 1, 11), P(21, 1, 22))),
										[]ast.ParameterNode{
											ast.NewFormalParameterNode(
												L(S(P(11, 1, 12), P(11, 1, 12))),
												ast.NewPublicIdentifierNode(L(S(P(11, 1, 12), P(11, 1, 12))), "i"),
												nil,
												nil,
												ast.NormalParameterKind,
											),
										},
										nil,
										nil,
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												L(S(P(17, 1, 18), P(21, 1, 22))),
												ast.NewBinaryExpressionNode(
													L(S(P(17, 1, 18), P(21, 1, 22))),
													T(L(S(P(19, 1, 20), P(19, 1, 20))), token.STAR),
													ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "i"),
													ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "2"),
												),
											),
										},
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a private name with an implicit receiver": {
			input: "_foo()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(5, 1, 6))),
							ast.NewPrivateIdentifierNode(L(S(P(0, 1, 1), P(3, 1, 4))), "_foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have unquote as the name": {
			input: "foo.unquote(bar * 2)()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 1, 22))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 1, 22))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(21, 1, 22))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewUnquoteNode(
								L(S(P(4, 1, 5), P(19, 1, 20))),
								ast.UNQUOTE_IDENTIFIER_KIND,
								ast.NewBinaryExpressionNode(
									L(S(P(12, 1, 13), P(18, 1, 19))),
									T(L(S(P(16, 1, 17), P(16, 1, 17))), token.STAR),
									ast.NewPublicIdentifierNode(L(S(P(12, 1, 13), P(14, 1, 15))), "bar"),
									ast.NewIntLiteralNode(L(S(P(18, 1, 19), P(18, 1, 19))), "2"),
								),
							),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have short unquote as the name": {
			input: "foo.!{bar * 2}()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewUnquoteNode(
								L(S(P(4, 1, 5), P(13, 1, 14))),
								ast.UNQUOTE_IDENTIFIER_KIND,
								ast.NewBinaryExpressionNode(
									L(S(P(6, 1, 7), P(12, 1, 13))),
									T(L(S(P(10, 1, 11), P(10, 1, 11))), token.STAR),
									ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
									ast.NewIntLiteralNode(L(S(P(12, 1, 13), P(12, 1, 13))), "2"),
								),
							),
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
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver with type arguments": {
			input: "foo.bar::[String]()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(18, 1, 19))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(18, 1, 19))),
						ast.NewGenericMethodCallNode(
							L(S(P(0, 1, 1), P(18, 1, 19))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							[]ast.TypeNode{
								ast.NewPublicConstantNode(L(S(P(10, 1, 11), P(15, 1, 16))), "String"),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be multiline, have an explicit receiver with type arguments": {
			input: `
				foo.bar::[
					String,
				]()
			`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(36, 4, 8))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(
						L(S(P(0, 1, 1), P(0, 1, 1))),
					),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 5), P(36, 4, 8))),
						ast.NewGenericMethodCallNode(
							L(S(P(5, 2, 5), P(35, 4, 7))),
							ast.NewPublicIdentifierNode(L(S(P(5, 2, 5), P(7, 2, 7))), "foo"),
							T(L(S(P(8, 2, 8), P(8, 2, 8))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(9, 2, 9), P(11, 2, 11))), "bar"),
							[]ast.TypeNode{
								ast.NewPublicConstantNode(L(S(P(21, 3, 6), P(26, 3, 11))), "String"),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver and a trailing closure without pipes": {
			input: "foo.bar() -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(17, 1, 18))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(17, 1, 18))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(17, 1, 18))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(10, 1, 11), P(17, 1, 18))),
									nil,
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(13, 1, 14), P(17, 1, 18))),
											ast.NewBinaryExpressionNode(
												L(S(P(13, 1, 14), P(17, 1, 18))),
												T(L(S(P(15, 1, 16), P(15, 1, 16))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "i"),
												ast.NewIntLiteralNode(L(S(P(17, 1, 18), P(17, 1, 18))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver and a trailing closure without arguments": {
			input: "foo.bar() || -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(20, 1, 21))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(20, 1, 21))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(20, 1, 21))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(10, 1, 11), P(20, 1, 21))),
									nil,
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(16, 1, 17), P(20, 1, 21))),
											ast.NewBinaryExpressionNode(
												L(S(P(16, 1, 17), P(20, 1, 21))),
												T(L(S(P(18, 1, 19), P(18, 1, 19))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(16, 1, 17), P(16, 1, 17))), "i"),
												ast.NewIntLiteralNode(L(S(P(20, 1, 21), P(20, 1, 21))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver and a trailing closure": {
			input: "foo.bar() |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 1, 22))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 1, 22))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(21, 1, 22))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(10, 1, 11), P(21, 1, 22))),
									[]ast.ParameterNode{
										ast.NewFormalParameterNode(
											L(S(P(11, 1, 12), P(11, 1, 12))),
											ast.NewPublicIdentifierNode(L(S(P(11, 1, 12), P(11, 1, 12))), "i"),
											nil,
											nil,
											ast.NormalParameterKind,
										),
									},
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(17, 1, 18), P(21, 1, 22))),
											ast.NewBinaryExpressionNode(
												L(S(P(17, 1, 18), P(21, 1, 22))),
												T(L(S(P(19, 1, 20), P(19, 1, 20))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "i"),
												ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver, arguments and a trailing closure": {
			input: "foo.bar(1, 5) |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(25, 1, 26))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(25, 1, 26))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(25, 1, 26))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "1"),
								ast.NewIntLiteralNode(L(S(P(11, 1, 12), P(11, 1, 12))), "5"),
								ast.NewClosureLiteralNode(
									L(S(P(14, 1, 15), P(25, 1, 26))),
									[]ast.ParameterNode{
										ast.NewFormalParameterNode(
											L(S(P(15, 1, 16), P(15, 1, 16))),
											ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(15, 1, 16))), "i"),
											nil,
											nil,
											ast.NormalParameterKind,
										),
									},
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(21, 1, 22), P(25, 1, 26))),
											ast.NewBinaryExpressionNode(
												L(S(P(21, 1, 22), P(25, 1, 26))),
												T(L(S(P(23, 1, 24), P(23, 1, 24))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(21, 1, 22))), "i"),
												ast.NewIntLiteralNode(L(S(P(25, 1, 26), P(25, 1, 26))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver, named arguments and a trailing closure": {
			input: "foo.bar(f: 5) |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(25, 1, 26))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(25, 1, 26))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(25, 1, 26))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(8, 1, 9), P(11, 1, 12))),
									ast.NewPublicIdentifierNode(L(S(P(8, 1, 9), P(8, 1, 9))), "f"),
									ast.NewIntLiteralNode(L(S(P(11, 1, 12), P(11, 1, 12))), "5"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(14, 1, 15), P(25, 1, 26))),
									ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(25, 1, 26))), "func"),
									ast.NewClosureLiteralNode(
										L(S(P(14, 1, 15), P(25, 1, 26))),
										[]ast.ParameterNode{
											ast.NewFormalParameterNode(
												L(S(P(15, 1, 16), P(15, 1, 16))),
												ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(15, 1, 16))), "i"),
												nil,
												nil,
												ast.NormalParameterKind,
											),
										},
										nil,
										nil,
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												L(S(P(21, 1, 22), P(25, 1, 26))),
												ast.NewBinaryExpressionNode(
													L(S(P(21, 1, 22), P(25, 1, 26))),
													T(L(S(P(23, 1, 24), P(23, 1, 24))), token.STAR),
													ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(21, 1, 22))), "i"),
													ast.NewIntLiteralNode(L(S(P(25, 1, 26), P(25, 1, 26))), "2"),
												),
											),
										},
									),
								),
							},
						),
					),
				},
			),
		},
		"can omit parentheses": {
			input: "foo.bar 1",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "1"),
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
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(4, 1, 5))), token.QUESTION_DOT),
							ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "bar"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can use the cascade call operator": {
			input: "foo..bar",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(4, 1, 5))), token.DOT_DOT),
							ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "bar"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can use the safe cascade call operator": {
			input: "foo?..bar",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(5, 1, 6))), token.QUESTION_DOT_DOT),
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
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
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							ast.NewMethodCallNode(
								L(S(P(0, 1, 1), P(8, 1, 9))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
								ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
								nil,
								nil,
							),
							T(L(S(P(9, 1, 10), P(9, 1, 10))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(12, 1, 13))), "baz"),
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
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(1, 1, 2), P(14, 1, 15))),
						ast.NewMethodCallNode(
							L(S(P(1, 1, 2), P(14, 1, 15))),
							ast.NewBinaryExpressionNode(
								L(S(P(1, 1, 2), P(7, 1, 8))),
								T(L(S(P(5, 1, 6), P(5, 1, 6))), token.PLUS),
								ast.NewPublicIdentifierNode(L(S(P(1, 1, 2), P(3, 1, 4))), "foo"),
								ast.NewIntLiteralNode(L(S(P(7, 1, 8), P(7, 1, 8))), "2"),
							),
							T(L(S(P(9, 1, 10), P(9, 1, 10))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(12, 1, 13))), "bar"),
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
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewPrivateIdentifierNode(L(S(P(4, 1, 5), P(7, 1, 8))), "_bar"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(4, 1, 5), P(7, 1, 8))), "unexpected PRIVATE_IDENTIFIER, expected a public method name (public identifier, keyword or overridable operator)"),
			},
		},
		"can have an overridable operator as the method name with an explicit receiver": {
			input: "foo.+()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(4, 1, 5))), "+"),
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
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
							ast.NewInvalidNode(
								L(S(P(4, 1, 5), P(5, 1, 6))),
								T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(4, 1, 5), P(5, 1, 6))), "unexpected &&, expected a public method name (public identifier, keyword or overridable operator)"),
			},
		},
		"can call a private method on self": {
			input: "self._foo()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewMethodCallNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
							T(L(S(P(4, 1, 5), P(4, 1, 5))), token.DOT),
							ast.NewPrivateIdentifierNode(L(S(P(5, 1, 6), P(8, 1, 9))), "_foo"),
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
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(19, 1, 20))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(19, 1, 20))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have splat arguments": {
			input: "foo(*baz, 'foo', *bar)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 1, 22))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 1, 22))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(21, 1, 22))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewSplatExpressionNode(
									L(S(P(4, 1, 5), P(7, 1, 8))),
									ast.NewPublicIdentifierNode(
										L(S(P(5, 1, 6), P(7, 1, 8))),
										"baz",
									),
								),
								ast.NewRawStringLiteralNode(L(S(P(10, 1, 11), P(14, 1, 15))), "foo"),
								ast.NewSplatExpressionNode(
									L(S(P(17, 1, 18), P(20, 1, 21))),
									ast.NewPublicIdentifierNode(
										L(S(P(18, 1, 19), P(20, 1, 21))),
										"bar",
									),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have a trailing comma": {
			input: "foo(.1, 'foo', :bar,)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(20, 1, 21))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(20, 1, 21))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(20, 1, 21))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
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
				L(S(P(0, 1, 1), P(24, 1, 25))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(24, 1, 25))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(24, 1, 25))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(4, 1, 5), P(12, 1, 13))),
									ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(9, 1, 10), P(12, 1, 13))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(15, 1, 16), P(23, 1, 24))),
									ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(17, 1, 18))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(20, 1, 21), P(23, 1, 24)))),
								),
							},
						),
					),
				},
			),
		},
		"can have double splat arguments": {
			input: "foo(**f, bar: :baz, **dupa(), elk: true)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewDoubleSplatExpressionNode(
									L(S(P(4, 1, 5), P(6, 1, 7))),
									ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(6, 1, 7))), "f"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(9, 1, 10), P(17, 1, 18))),
									ast.NewPublicIdentifierNode(L(S(P(9, 1, 10), P(11, 1, 12))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(14, 1, 15), P(17, 1, 18))), "baz"),
								),
								ast.NewDoubleSplatExpressionNode(
									L(S(P(20, 1, 21), P(27, 1, 28))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(22, 1, 23), P(27, 1, 28))),
										ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(25, 1, 26))), "dupa"),
										nil,
										nil,
									),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(30, 1, 31), P(38, 1, 39))),
									ast.NewPublicIdentifierNode(L(S(P(30, 1, 31), P(32, 1, 33))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(35, 1, 36), P(38, 1, 39)))),
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
				L(S(P(0, 1, 1), P(41, 1, 42))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(41, 1, 42))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(41, 1, 42))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 1, 22), P(29, 1, 30))),
									ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(23, 1, 24))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 1, 27), P(29, 1, 30))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 1, 33), P(40, 1, 41))),
									ast.NewPublicIdentifierNode(L(S(P(32, 1, 33), P(34, 1, 35))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 1, 38), P(40, 1, 41)))),
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
				L(S(P(0, 1, 1), P(41, 4, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(41, 4, 10))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(41, 4, 10))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 2, 1), P(12, 2, 5))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 3, 1), P(18, 3, 4))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 3, 7), P(29, 3, 15))),
									ast.NewPublicIdentifierNode(L(S(P(21, 3, 7), P(23, 3, 9))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 3, 12), P(29, 3, 15))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 4, 1), P(40, 4, 9))),
									ast.NewPublicIdentifierNode(L(S(P(32, 4, 1), P(34, 4, 3))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 4, 6), P(40, 4, 9)))),
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
				L(S(P(0, 1, 1), P(43, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(43, 3, 1))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(43, 3, 1))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(5, 2, 1), P(6, 2, 2))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(9, 2, 5), P(13, 2, 9))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(16, 2, 12), P(19, 2, 15))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(22, 2, 18), P(30, 2, 26))),
									ast.NewPublicIdentifierNode(L(S(P(22, 2, 18), P(24, 2, 20))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(27, 2, 23), P(30, 2, 26))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(33, 2, 29), P(41, 2, 37))),
									ast.NewPublicIdentifierNode(L(S(P(33, 2, 29), P(35, 2, 31))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(38, 2, 34), P(41, 2, 37)))),
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
				L(S(P(0, 1, 1), P(42, 2, 39))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(3, 1, 4))),
						ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
					),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 2), P(42, 2, 39))),
						ast.NewFloatLiteralNode(L(S(P(5, 2, 2), P(6, 2, 3))), "0.1"),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(7, 2, 4), P(7, 2, 4))), "unexpected ,, expected )"),
			},
		},
		"can have positional arguments without parentheses": {
			input: "foo .1, 'foo', :bar",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(18, 1, 19))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(18, 1, 19))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(18, 1, 19))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
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
				L(S(P(0, 1, 1), P(23, 1, 24))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(23, 1, 24))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(23, 1, 24))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(4, 1, 5), P(12, 1, 13))),
									ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(9, 1, 10), P(12, 1, 13))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(15, 1, 16), P(23, 1, 24))),
									ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(17, 1, 18))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(20, 1, 21), P(23, 1, 24)))),
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
				L(S(P(0, 1, 1), P(40, 1, 41))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(40, 1, 41))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(40, 1, 41))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 1, 22), P(29, 1, 30))),
									ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(23, 1, 24))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 1, 27), P(29, 1, 30))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 1, 33), P(40, 1, 41))),
									ast.NewPublicIdentifierNode(L(S(P(32, 1, 33), P(34, 1, 35))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 1, 38), P(40, 1, 41)))),
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
				L(S(P(0, 1, 1), P(40, 4, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(40, 4, 9))),
						ast.NewReceiverlessMethodCallNode(
							L(S(P(0, 1, 1), P(40, 4, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(8, 2, 1), P(12, 2, 5))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(15, 3, 1), P(18, 3, 4))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(21, 3, 7), P(29, 3, 15))),
									ast.NewPublicIdentifierNode(L(S(P(21, 3, 7), P(23, 3, 9))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(26, 3, 12), P(29, 3, 15))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(32, 4, 1), P(40, 4, 9))),
									ast.NewPublicIdentifierNode(L(S(P(32, 4, 1), P(34, 4, 3))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(37, 4, 6), P(40, 4, 9)))),
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
				L(S(P(0, 1, 1), P(5, 2, 2))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(3, 1, 4))),
						ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
					),
					ast.NewExpressionStatementNode(
						L(S(P(4, 2, 1), P(5, 2, 2))),
						ast.NewFloatLiteralNode(L(S(P(4, 2, 1), P(5, 2, 2))), "0.1"),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 2, 3), P(6, 2, 3))), "unexpected ,, expected a statement separator `\\n`, `;`"),
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
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewAttributeAccessNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
							ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "bar"),
						),
					),
				},
			),
		},
		"can have unquote as the name": {
			input: "foo.unquote(foo + 1)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(19, 1, 20))),
						ast.NewAttributeAccessNode(
							L(S(P(0, 1, 1), P(19, 1, 20))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewUnquoteNode(
								L(S(P(4, 1, 5), P(19, 1, 20))),
								ast.UNQUOTE_IDENTIFIER_KIND,
								ast.NewBinaryExpressionNode(
									L(S(P(12, 1, 13), P(18, 1, 19))),
									T(L(S(P(16, 1, 17), P(16, 1, 17))), token.PLUS),
									ast.NewPublicIdentifierNode(L(S(P(12, 1, 13), P(14, 1, 15))), "foo"),
									ast.NewIntLiteralNode(L(S(P(18, 1, 19), P(18, 1, 19))), "1"),
								),
							),
						),
					),
				},
			),
		},
		"can have short unquote as the name": {
			input: "foo.!{foo + 1}",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(13, 1, 14))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(13, 1, 14))),
						ast.NewAttributeAccessNode(
							L(S(P(0, 1, 1), P(13, 1, 14))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewUnquoteNode(
								L(S(P(4, 1, 5), P(13, 1, 14))),
								ast.UNQUOTE_IDENTIFIER_KIND,
								ast.NewBinaryExpressionNode(
									L(S(P(6, 1, 7), P(12, 1, 13))),
									T(L(S(P(10, 1, 11), P(10, 1, 11))), token.PLUS),
									ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
									ast.NewIntLiteralNode(L(S(P(12, 1, 13), P(12, 1, 13))), "1"),
								),
							),
						),
					),
				},
			),
		},
		"can be called on variables": {
			input: "foo.bar",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewAttributeAccessNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
						),
					),
				},
			),
		},
		"can be nested": {
			input: "foo.bar.baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewAttributeAccessNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							ast.NewAttributeAccessNode(
								L(S(P(0, 1, 1), P(6, 1, 7))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(8, 1, 9), P(10, 1, 11))), "baz"),
						),
					),
				},
			),
		},
		"can have newlines after the dot": {
			input: "foo.\nbar.\nbaz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 3, 3))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(12, 3, 3))),
						ast.NewAttributeAccessNode(
							L(S(P(0, 1, 1), P(12, 3, 3))),
							ast.NewAttributeAccessNode(
								L(S(P(0, 1, 1), P(7, 2, 3))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(5, 2, 1), P(7, 2, 3))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(10, 3, 1), P(12, 3, 3))), "baz"),
						),
					),
				},
			),
		},
		"can have newlines before the dot": {
			input: "foo\n.bar\n.baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 3, 4))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(12, 3, 4))),
						ast.NewAttributeAccessNode(
							L(S(P(0, 1, 1), P(12, 3, 4))),
							ast.NewAttributeAccessNode(
								L(S(P(0, 1, 1), P(7, 2, 4))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(5, 2, 2), P(7, 2, 4))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(10, 3, 2), P(12, 3, 4))), "baz"),
						),
					),
				},
			),
		},
		"can be nested on function calls": {
			input: "foo().bar.baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 1, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(12, 1, 13))),
						ast.NewAttributeAccessNode(
							L(S(P(0, 1, 1), P(12, 1, 13))),
							ast.NewAttributeAccessNode(
								L(S(P(0, 1, 1), P(8, 1, 9))),
								ast.NewReceiverlessMethodCallNode(
									L(S(P(0, 1, 1), P(4, 1, 5))),
									ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
									nil,
									nil,
								),
								ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(12, 1, 13))), "baz"),
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

func TestCall(t *testing.T) {
	tests := testTable{
		"can be used on self": {
			input: "self.()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewCallNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
							false,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have arguments": {
			input: "self.(.1, 'foo', :bar, bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(43, 1, 44))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(43, 1, 44))),
						ast.NewCallNode(
							L(S(P(0, 1, 1), P(43, 1, 44))),
							ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
							false,
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(6, 1, 7), P(7, 1, 8))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(10, 1, 11), P(14, 1, 15))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(17, 1, 18), P(20, 1, 21))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(23, 1, 24), P(31, 1, 32))),
									ast.NewPublicIdentifierNode(L(S(P(23, 1, 24), P(25, 1, 26))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(28, 1, 29), P(31, 1, 32))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(34, 1, 35), P(42, 1, 43))),
									ast.NewPublicIdentifierNode(L(S(P(34, 1, 35), P(36, 1, 37))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(39, 1, 40), P(42, 1, 43)))),
								),
							},
						),
					),
				},
			),
		},
		"can be nil-safe": {
			input: "self?.()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewCallNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
							true,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be called on variables": {
			input: "foo.()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewCallNode(
							L(S(P(0, 1, 1), P(5, 1, 6))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							false,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be nested": {
			input: "foo.().()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewCallNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewCallNode(
								L(S(P(0, 1, 1), P(5, 1, 6))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								false,
								nil,
								nil,
							),
							false,
							nil,
							nil,
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

func TestSubscript(t *testing.T) {
	tests := testTable{
		"can be used on self": {
			input: "self[5]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
							ast.NewIntLiteralNode(
								L(S(P(5, 1, 6), P(5, 1, 6))),
								"5",
							),
						),
					),
				},
			),
		},
		"can be multiline": {
			input: `
				self[
					5
				]
			`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(23, 4, 6))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 5), P(23, 4, 6))),
						ast.NewSubscriptExpressionNode(
							L(S(P(5, 2, 5), P(22, 4, 5))),
							ast.NewSelfLiteralNode(L(S(P(5, 2, 5), P(8, 2, 8)))),
							ast.NewIntLiteralNode(
								L(S(P(16, 3, 6), P(16, 3, 6))),
								"5",
							),
						),
					),
				},
			),
		},
		"can have any expression inside brackets": {
			input: "self[5 + foo]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 1, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(12, 1, 13))),
						ast.NewSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(12, 1, 13))),
							ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
							ast.NewBinaryExpressionNode(
								L(S(P(5, 1, 6), P(11, 1, 12))),
								T(L(S(P(7, 1, 8), P(7, 1, 8))), token.PLUS),
								ast.NewIntLiteralNode(
									L(S(P(5, 1, 6), P(5, 1, 6))),
									"5",
								),
								ast.NewPublicIdentifierNode(L(S(P(9, 1, 10), P(11, 1, 12))), "foo"),
							),
						),
					),
				},
			),
		},
		"can be used on attribute access": {
			input: "self.foo[5]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							ast.NewAttributeAccessNode(
								L(S(P(0, 1, 1), P(7, 1, 8))),
								ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
								ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "foo"),
							),
							ast.NewIntLiteralNode(
								L(S(P(9, 1, 10), P(9, 1, 10))),
								"5",
							),
						),
					),
				},
			),
		},
		"can be mixed with method calls": {
			input: "self.foo[5][2].bar[0]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(20, 1, 21))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(20, 1, 21))),
						ast.NewSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(20, 1, 21))),
							ast.NewAttributeAccessNode(
								L(S(P(0, 1, 1), P(17, 1, 18))),
								ast.NewSubscriptExpressionNode(
									L(S(P(0, 1, 1), P(13, 1, 14))),
									ast.NewSubscriptExpressionNode(
										L(S(P(0, 1, 1), P(10, 1, 11))),
										ast.NewAttributeAccessNode(
											L(S(P(0, 1, 1), P(7, 1, 8))),
											ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
											ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "foo"),
										),
										ast.NewIntLiteralNode(
											L(S(P(9, 1, 10), P(9, 1, 10))),
											"5",
										),
									),
									ast.NewIntLiteralNode(
										L(S(P(12, 1, 13), P(12, 1, 13))),
										"2",
									),
								),
								ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(17, 1, 18))), "bar"),
							),
							ast.NewIntLiteralNode(
								L(S(P(19, 1, 20), P(19, 1, 20))),
								"0",
							),
						),
					),
				},
			),
		},
		"can be used on method calls": {
			input: "self.foo()[5]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 1, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(12, 1, 13))),
						ast.NewSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(12, 1, 13))),
							ast.NewMethodCallNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
								T(L(S(P(4, 1, 5), P(4, 1, 5))), token.DOT),
								ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "foo"),
								nil,
								nil,
							),
							ast.NewIntLiteralNode(
								L(S(P(11, 1, 12), P(11, 1, 12))),
								"5",
							),
						),
					),
				},
			),
		},
		"can be used on function calls": {
			input: "foo()[5]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewReceiverlessMethodCallNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								nil,
								nil,
							),
							ast.NewIntLiteralNode(
								L(S(P(6, 1, 7), P(6, 1, 7))),
								"5",
							),
						),
					),
				},
			),
		},
		"can be nested": {
			input: "foo[5][20][3]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 1, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(12, 1, 13))),
						ast.NewSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(12, 1, 13))),
							ast.NewSubscriptExpressionNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								ast.NewSubscriptExpressionNode(
									L(S(P(0, 1, 1), P(5, 1, 6))),
									ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
									ast.NewIntLiteralNode(
										L(S(P(4, 1, 5), P(4, 1, 5))),
										"5",
									),
								),
								ast.NewIntLiteralNode(
									L(S(P(7, 1, 8), P(8, 1, 9))),
									"20",
								),
							),
							ast.NewIntLiteralNode(
								L(S(P(11, 1, 12), P(11, 1, 12))),
								"3",
							),
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

func TestNilSafeSubscript(t *testing.T) {
	tests := testTable{
		"can be used on self": {
			input: "self?[5]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewNilSafeSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
							ast.NewIntLiteralNode(
								L(S(P(6, 1, 7), P(6, 1, 7))),
								"5",
							),
						),
					),
				},
			),
		},
		"can be multiline": {
			input: `
				self?[
					5
				]
			`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(24, 4, 6))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 5), P(24, 4, 6))),
						ast.NewNilSafeSubscriptExpressionNode(
							L(S(P(5, 2, 5), P(23, 4, 5))),
							ast.NewSelfLiteralNode(L(S(P(5, 2, 5), P(8, 2, 8)))),
							ast.NewIntLiteralNode(
								L(S(P(17, 3, 6), P(17, 3, 6))),
								"5",
							),
						),
					),
				},
			),
		},
		"can have any expression inside brackets": {
			input: "self?[5 + foo]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(13, 1, 14))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(13, 1, 14))),
						ast.NewNilSafeSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(13, 1, 14))),
							ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
							ast.NewBinaryExpressionNode(
								L(S(P(6, 1, 7), P(12, 1, 13))),
								T(L(S(P(8, 1, 9), P(8, 1, 9))), token.PLUS),
								ast.NewIntLiteralNode(
									L(S(P(6, 1, 7), P(6, 1, 7))),
									"5",
								),
								ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(12, 1, 13))), "foo"),
							),
						),
					),
				},
			),
		},
		"can be used on attribute access": {
			input: "self.foo?[5]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewNilSafeSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewAttributeAccessNode(
								L(S(P(0, 1, 1), P(7, 1, 8))),
								ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
								ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "foo"),
							),
							ast.NewIntLiteralNode(
								L(S(P(10, 1, 11), P(10, 1, 11))),
								"5",
							),
						),
					),
				},
			),
		},
		"can be used on method calls": {
			input: "self.foo()?[5]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(13, 1, 14))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(13, 1, 14))),
						ast.NewNilSafeSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(13, 1, 14))),
							ast.NewMethodCallNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
								T(L(S(P(4, 1, 5), P(4, 1, 5))), token.DOT),
								ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "foo"),
								nil,
								nil,
							),
							ast.NewIntLiteralNode(
								L(S(P(12, 1, 13), P(12, 1, 13))),
								"5",
							),
						),
					),
				},
			),
		},
		"can be used on function calls": {
			input: "foo()?[5]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewNilSafeSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewReceiverlessMethodCallNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								nil,
								nil,
							),
							ast.NewIntLiteralNode(
								L(S(P(7, 1, 8), P(7, 1, 8))),
								"5",
							),
						),
					),
				},
			),
		},
		"can be nested": {
			input: "foo?[5]?[20]?[3]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewNilSafeSubscriptExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							ast.NewNilSafeSubscriptExpressionNode(
								L(S(P(0, 1, 1), P(11, 1, 12))),
								ast.NewNilSafeSubscriptExpressionNode(
									L(S(P(0, 1, 1), P(6, 1, 7))),
									ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
									ast.NewIntLiteralNode(
										L(S(P(5, 1, 6), P(5, 1, 6))),
										"5",
									),
								),
								ast.NewIntLiteralNode(
									L(S(P(9, 1, 10), P(10, 1, 11))),
									"20",
								),
							),
							ast.NewIntLiteralNode(
								L(S(P(14, 1, 15), P(14, 1, 15))),
								"3",
							),
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
