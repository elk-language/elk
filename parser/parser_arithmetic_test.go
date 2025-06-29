package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/token"
)

func TestAddition(t *testing.T) {
	tests := testTable{
		"can contain unquote": {
			input: "1 + 2 + unquote(foo)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(19, 1, 20))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(19, 1, 20))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.PLUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewUnquoteNode(
								L(S(P(8, 1, 9), P(19, 1, 20))),
								ast.UNQUOTE_EXPRESSION_KIND,
								ast.NewPublicIdentifierNode(L(S(P(16, 1, 17), P(18, 1, 19))), "foo"),
							),
						),
					),
				},
			),
		},
		"can contain unquote_expr": {
			input: "1 + 2 + unquote_expr(foo)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(24, 1, 25))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(24, 1, 25))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(24, 1, 25))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.PLUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewUnquoteNode(
								L(S(P(8, 1, 9), P(24, 1, 25))),
								ast.UNQUOTE_EXPRESSION_KIND,
								ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(23, 1, 24))), "foo"),
							),
						),
					),
				},
			),
		},
		"can contain unquote_ident": {
			input: "1 + 2 + unquote_ident(foo)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(25, 1, 26))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(25, 1, 26))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(25, 1, 26))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.PLUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewUnquoteNode(
								L(S(P(8, 1, 9), P(25, 1, 26))),
								ast.UNQUOTE_IDENTIFIER_KIND,
								ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(24, 1, 25))), "foo"),
							),
						),
					),
				},
			),
		},
		"can contain unquote_const": {
			input: "1 + 2 + unquote_const(foo)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(25, 1, 26))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(25, 1, 26))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(25, 1, 26))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.PLUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewUnquoteNode(
								L(S(P(8, 1, 9), P(25, 1, 26))),
								ast.UNQUOTE_CONSTANT_KIND,
								ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(24, 1, 25))), "foo"),
							),
						),
					),
				},
			),
		},
		"can contain unquote_ivar": {
			input: "1 + 2 + unquote_ivar(foo)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(24, 1, 25))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(24, 1, 25))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(24, 1, 25))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.PLUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewUnquoteNode(
								L(S(P(8, 1, 9), P(24, 1, 25))),
								ast.UNQUOTE_INSTANCE_VARIABLE_KIND,
								ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(23, 1, 24))), "foo"),
							),
						),
					),
				},
			),
		},
		"can contain short unquote": {
			input: "1 + 2 + !{foo}",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(13, 1, 14))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(13, 1, 14))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(13, 1, 14))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.PLUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewUnquoteNode(
								L(S(P(8, 1, 9), P(13, 1, 14))),
								ast.UNQUOTE_EXPRESSION_KIND,
								ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(12, 1, 13))), "foo"),
							),
						),
					),
				},
			),
		},
		"is evaluated from left to right": {
			input: "1 + 2 + 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.PLUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 +\n2 +\n3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 3, 1))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 3, 1))),
							T(L(S(P(6, 2, 3), P(6, 2, 3))), token.PLUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 2, 1))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 2, 1), P(4, 2, 1))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 3, 1), P(8, 3, 1))), "3"),
						),
					),
				},
			),
		},
		"cannot have newlines before the operator": {
			input: "1\n+ 2\n+ 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 3))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(1, 1, 2))),
						ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
					),
					ast.NewExpressionStatementNode(
						L(S(P(2, 2, 1), P(5, 2, 4))),
						ast.NewUnaryExpressionNode(
							L(S(P(2, 2, 1), P(4, 2, 3))),
							T(L(S(P(2, 2, 1), P(2, 2, 1))), token.PLUS),
							ast.NewIntLiteralNode(L(S(P(4, 2, 3), P(4, 2, 3))), "2"),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(6, 3, 1), P(8, 3, 3))),
						ast.NewUnaryExpressionNode(
							L(S(P(6, 3, 1), P(8, 3, 3))),
							T(L(S(P(6, 3, 1), P(6, 3, 1))), token.PLUS),
							ast.NewIntLiteralNode(L(S(P(8, 3, 3), P(8, 3, 3))), "3"),
						),
					),
				},
			),
		},
		"has higher precedence than bitshifts": {
			input: "foo >> bar + baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.RBITSHIFT),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewBinaryExpressionNode(
								L(S(P(7, 1, 8), P(15, 1, 16))),
								T(L(S(P(11, 1, 12), P(11, 1, 12))), token.PLUS),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
								ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(15, 1, 16))), "baz"),
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

func TestSubtraction(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 - 2 - 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.MINUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.MINUS),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 -\n2 -\n3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 3, 1))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 3, 1))),
							T(L(S(P(6, 2, 3), P(6, 2, 3))), token.MINUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 2, 1))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.MINUS),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 2, 1), P(4, 2, 1))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 3, 1), P(8, 3, 1))), "3"),
						),
					),
				},
			),
		},
		"cannot have newlines before the operator": {
			input: "1\n- 2\n- 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 3))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(1, 1, 2))),
						ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
					),
					ast.NewExpressionStatementNode(
						L(S(P(2, 2, 1), P(5, 2, 4))),
						ast.NewUnaryExpressionNode(
							L(S(P(2, 2, 1), P(4, 2, 3))),
							T(L(S(P(2, 2, 1), P(2, 2, 1))), token.MINUS),
							ast.NewIntLiteralNode(L(S(P(4, 2, 3), P(4, 2, 3))), "2"),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(6, 3, 1), P(8, 3, 3))),
						ast.NewUnaryExpressionNode(
							L(S(P(6, 3, 1), P(8, 3, 3))),
							T(L(S(P(6, 3, 1), P(6, 3, 1))), token.MINUS),
							ast.NewIntLiteralNode(L(S(P(8, 3, 3), P(8, 3, 3))), "3"),
						),
					),
				},
			),
		},
		"has the same precedence as addition": {
			input: "1 + 2 - 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.MINUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
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

func TestMultiplication(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 * 2 * 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.STAR),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.STAR),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 *\n2 *\n3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 3, 1))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 3, 1))),
							T(L(S(P(6, 2, 3), P(6, 2, 3))), token.STAR),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 2, 1))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.STAR),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 2, 1), P(4, 2, 1))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 3, 1), P(8, 3, 1))), "3"),
						),
					),
				},
			),
		},
		"cannot have newlines before the operator": {
			input: "1\n* 2\n* 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 3))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(1, 1, 2))),
						ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
					),
					ast.NewExpressionStatementNode(
						L(S(P(2, 2, 1), P(5, 2, 4))),
						ast.NewInvalidNode(L(S(P(2, 2, 1), P(2, 2, 1))), T(L(S(P(2, 2, 1), P(2, 2, 1))), token.STAR)),
					),
					ast.NewExpressionStatementNode(
						L(S(P(6, 3, 1), P(8, 3, 3))),
						ast.NewInvalidNode(L(S(P(6, 3, 1), P(6, 3, 1))), T(L(S(P(6, 3, 1), P(6, 3, 1))), token.STAR)),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(2, 2, 1), P(2, 2, 1))), "unexpected *, expected an expression"),
				diagnostic.NewFailure(L(S(P(6, 3, 1), P(6, 3, 1))), "unexpected *, expected an expression"),
			},
		},
		"has higher precedence than addition": {
			input: "1 + 2 * 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
							ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
							ast.NewBinaryExpressionNode(
								L(S(P(4, 1, 5), P(8, 1, 9))),
								T(L(S(P(6, 1, 7), P(6, 1, 7))), token.STAR),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
								ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
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

func TestDivision(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 / 2 / 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.SLASH),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.SLASH),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 /\n2 /\n3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 3, 1))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 3, 1))),
							T(L(S(P(6, 2, 3), P(6, 2, 3))), token.SLASH),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 2, 1))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.SLASH),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 2, 1), P(4, 2, 1))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 3, 1), P(8, 3, 1))), "3"),
						),
					),
				},
			),
		},
		"cannot have newlines before the operator": {
			input: "1\n/ 2\n/ 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 3))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(1, 1, 2))),
						ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
					),
					ast.NewExpressionStatementNode(
						L(S(P(2, 2, 1), P(5, 2, 4))),
						ast.NewInvalidNode(L(S(P(2, 2, 1), P(2, 2, 1))), T(L(S(P(2, 2, 1), P(2, 2, 1))), token.SLASH)),
					),
					ast.NewExpressionStatementNode(
						L(S(P(6, 3, 1), P(8, 3, 3))),
						ast.NewInvalidNode(L(S(P(6, 3, 1), P(6, 3, 1))), T(L(S(P(6, 3, 1), P(6, 3, 1))), token.SLASH)),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(2, 2, 1), P(2, 2, 1))), "unexpected /, expected an expression"),
				diagnostic.NewFailure(L(S(P(6, 3, 1), P(6, 3, 1))), "unexpected /, expected an expression"),
			},
		},
		"has the same precedence as multiplication": {
			input: "1 * 2 / 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.SLASH),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.STAR),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
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

func TestModulo(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 % 2 % 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.PERCENT),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PERCENT),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 %\n2 %\n3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 3, 1))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 3, 1))),
							T(L(S(P(6, 2, 3), P(6, 2, 3))), token.PERCENT),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 2, 1))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PERCENT),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 2, 1), P(4, 2, 1))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 3, 1), P(8, 3, 1))), "3"),
						),
					),
				},
			),
		},
		"cannot have newlines before the operator": {
			input: "1\n% 2\n% 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 3))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(1, 1, 2))),
						ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
					),
					ast.NewExpressionStatementNode(
						L(S(P(2, 2, 1), P(5, 2, 4))),
						ast.NewInvalidNode(L(S(P(2, 2, 1), P(2, 2, 1))), T(L(S(P(2, 2, 1), P(2, 2, 1))), token.PERCENT)),
					),
					ast.NewExpressionStatementNode(
						L(S(P(6, 3, 1), P(8, 3, 3))),
						ast.NewInvalidNode(L(S(P(6, 3, 1), P(6, 3, 1))), T(L(S(P(6, 3, 1), P(6, 3, 1))), token.PERCENT)),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(2, 2, 1), P(2, 2, 1))), "unexpected %, expected an expression"),
				diagnostic.NewFailure(L(S(P(6, 3, 1), P(6, 3, 1))), "unexpected %, expected an expression"),
			},
		},
		"has the same precedence as multiplication": {
			input: "1 * 2 % 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.PERCENT),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.STAR),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
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

func TestUnaryExpressions(t *testing.T) {
	tests := testTable{
		"plus cannot be nested without spaces": {
			input: "+++1.5",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(5, 1, 6))),
							T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
							ast.NewInvalidNode(L(S(P(0, 1, 1), P(1, 1, 2))), T(L(S(P(0, 1, 1), P(1, 1, 2))), token.PLUS_PLUS)),
							ast.NewFloatLiteralNode(L(S(P(3, 1, 4), P(5, 1, 6))), "1.5"),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(0, 1, 1), P(1, 1, 2))), "unexpected ++, expected an expression"),
			},
		},
		"minus cannot be nested without spaces": {
			input: "---1.5",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(5, 1, 6))),
							T(L(S(P(2, 1, 3), P(2, 1, 3))), token.MINUS),
							ast.NewInvalidNode(L(S(P(0, 1, 1), P(1, 1, 2))), T(L(S(P(0, 1, 1), P(1, 1, 2))), token.MINUS_MINUS)),
							ast.NewFloatLiteralNode(L(S(P(3, 1, 4), P(5, 1, 6))), "1.5"),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(0, 1, 1), P(1, 1, 2))), "unexpected --, expected an expression"),
			},
		},
		"plus can be nested": {
			input: "+ + +1.5",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewUnaryExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							T(L(S(P(0, 1, 1), P(0, 1, 1))), token.PLUS),
							ast.NewUnaryExpressionNode(
								L(S(P(2, 1, 3), P(7, 1, 8))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.PLUS),
								ast.NewUnaryExpressionNode(
									L(S(P(4, 1, 5), P(7, 1, 8))),
									T(L(S(P(4, 1, 5), P(4, 1, 5))), token.PLUS),
									ast.NewFloatLiteralNode(L(S(P(5, 1, 6), P(7, 1, 8))), "1.5"),
								),
							),
						),
					),
				},
			),
		},
		"minus can be nested": {
			input: "- - -1.5",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewUnaryExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							T(L(S(P(0, 1, 1), P(0, 1, 1))), token.MINUS),
							ast.NewUnaryExpressionNode(
								L(S(P(2, 1, 3), P(7, 1, 8))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.MINUS),
								ast.NewUnaryExpressionNode(
									L(S(P(4, 1, 5), P(7, 1, 8))),
									T(L(S(P(4, 1, 5), P(4, 1, 5))), token.MINUS),
									ast.NewFloatLiteralNode(L(S(P(5, 1, 6), P(7, 1, 8))), "1.5"),
								),
							),
						),
					),
				},
			),
		},
		"logical not can be nested": {
			input: "!!!1.5",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewUnaryExpressionNode(
							L(S(P(0, 1, 1), P(5, 1, 6))),
							T(L(S(P(0, 1, 1), P(0, 1, 1))), token.BANG),
							ast.NewUnaryExpressionNode(
								L(S(P(1, 1, 2), P(5, 1, 6))),
								T(L(S(P(1, 1, 2), P(1, 1, 2))), token.BANG),
								ast.NewUnaryExpressionNode(
									L(S(P(2, 1, 3), P(5, 1, 6))),
									T(L(S(P(2, 1, 3), P(2, 1, 3))), token.BANG),
									ast.NewFloatLiteralNode(L(S(P(3, 1, 4), P(5, 1, 6))), "1.5"),
								),
							),
						),
					),
				},
			),
		},
		"bitwise not can be nested": {
			input: "~~~1.5",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewUnaryExpressionNode(
							L(S(P(0, 1, 1), P(5, 1, 6))),
							T(L(S(P(0, 1, 1), P(0, 1, 1))), token.TILDE),
							ast.NewUnaryExpressionNode(
								L(S(P(1, 1, 2), P(5, 1, 6))),
								T(L(S(P(1, 1, 2), P(1, 1, 2))), token.TILDE),
								ast.NewUnaryExpressionNode(
									L(S(P(2, 1, 3), P(5, 1, 6))),
									T(L(S(P(2, 1, 3), P(2, 1, 3))), token.TILDE),
									ast.NewFloatLiteralNode(L(S(P(3, 1, 4), P(5, 1, 6))), "1.5"),
								),
							),
						),
					),
				},
			),
		},
		"all have the same precedence": {
			input: "!+~1.5",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewUnaryExpressionNode(
							L(S(P(0, 1, 1), P(5, 1, 6))),
							T(L(S(P(0, 1, 1), P(0, 1, 1))), token.BANG),
							ast.NewUnaryExpressionNode(
								L(S(P(1, 1, 2), P(5, 1, 6))),
								T(L(S(P(1, 1, 2), P(1, 1, 2))), token.PLUS),
								ast.NewUnaryExpressionNode(
									L(S(P(2, 1, 3), P(5, 1, 6))),
									T(L(S(P(2, 1, 3), P(2, 1, 3))), token.TILDE),
									ast.NewFloatLiteralNode(L(S(P(3, 1, 4), P(5, 1, 6))), "1.5"),
								),
							),
						),
					),
				},
			),
		},
		"have higher precedence than multiplicative expressions": {
			input: "!!1.5 * 2 + ~.5",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							T(L(S(P(10, 1, 11), P(10, 1, 11))), token.PLUS),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(8, 1, 9))),
								T(L(S(P(6, 1, 7), P(6, 1, 7))), token.STAR),
								ast.NewUnaryExpressionNode(
									L(S(P(0, 1, 1), P(4, 1, 5))),
									T(L(S(P(0, 1, 1), P(0, 1, 1))), token.BANG),
									ast.NewUnaryExpressionNode(
										L(S(P(1, 1, 2), P(4, 1, 5))),
										T(L(S(P(1, 1, 2), P(1, 1, 2))), token.BANG),
										ast.NewFloatLiteralNode(L(S(P(2, 1, 3), P(4, 1, 5))), "1.5"),
									),
								),
								ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "2"),
							),
							ast.NewUnaryExpressionNode(
								L(S(P(12, 1, 13), P(14, 1, 15))),
								T(L(S(P(12, 1, 13), P(12, 1, 13))), token.TILDE),
								ast.NewFloatLiteralNode(L(S(P(13, 1, 14), P(14, 1, 15))), "0.5"),
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

func TestExponentiation(t *testing.T) {
	tests := testTable{
		"is evaluated from right to left": {
			input: "1 ** 2 ** 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							T(L(S(P(2, 1, 3), P(3, 1, 4))), token.STAR_STAR),
							ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
							ast.NewBinaryExpressionNode(
								L(S(P(5, 1, 6), P(10, 1, 11))),
								T(L(S(P(7, 1, 8), P(8, 1, 9))), token.STAR_STAR),
								ast.NewIntLiteralNode(L(S(P(5, 1, 6), P(5, 1, 6))), "2"),
								ast.NewIntLiteralNode(L(S(P(10, 1, 11), P(10, 1, 11))), "3"),
							),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 **\n2 **\n3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 3, 1))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(10, 3, 1))),
							T(L(S(P(2, 1, 3), P(3, 1, 4))), token.STAR_STAR),
							ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
							ast.NewBinaryExpressionNode(
								L(S(P(5, 2, 1), P(10, 3, 1))),
								T(L(S(P(7, 2, 3), P(8, 2, 4))), token.STAR_STAR),
								ast.NewIntLiteralNode(L(S(P(5, 2, 1), P(5, 2, 1))), "2"),
								ast.NewIntLiteralNode(L(S(P(10, 3, 1), P(10, 3, 1))), "3"),
							),
						),
					),
				},
			),
		},
		"cannot have newlines before the operator": {
			input: "1\n** 2\n** 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 3, 4))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(1, 1, 2))),
						ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
					),
					ast.NewExpressionStatementNode(
						L(S(P(2, 2, 1), P(6, 2, 5))),
						ast.NewInvalidNode(L(S(P(2, 2, 1), P(3, 2, 2))), T(L(S(P(2, 2, 1), P(3, 2, 2))), token.STAR_STAR)),
					),
					ast.NewExpressionStatementNode(
						L(S(P(7, 3, 1), P(10, 3, 4))),
						ast.NewInvalidNode(L(S(P(7, 3, 1), P(8, 3, 2))), T(L(S(P(7, 3, 1), P(8, 3, 2))), token.STAR_STAR)),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(2, 2, 1), P(3, 2, 2))), "unexpected **, expected an expression"),
				diagnostic.NewFailure(L(S(P(7, 3, 1), P(8, 3, 2))), "unexpected **, expected an expression"),
			},
		},
		"has higher precedence than unary expressions": {
			input: "-2 ** 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewUnaryExpressionNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							T(L(S(P(0, 1, 1), P(0, 1, 1))), token.MINUS),
							ast.NewBinaryExpressionNode(
								L(S(P(1, 1, 2), P(6, 1, 7))),
								T(L(S(P(3, 1, 4), P(4, 1, 5))), token.STAR_STAR),
								ast.NewIntLiteralNode(L(S(P(1, 1, 2), P(1, 1, 2))), "2"),
								ast.NewIntLiteralNode(L(S(P(6, 1, 7), P(6, 1, 7))), "3"),
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

func TestBitwiseOr(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 | 2 | 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.OR),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.OR),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 |\n2 |\n3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 3, 1))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 3, 1))),
							T(L(S(P(6, 2, 3), P(6, 2, 3))), token.OR),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 2, 1))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.OR),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 2, 1), P(4, 2, 1))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 3, 1), P(8, 3, 1))), "3"),
						),
					),
				},
			),
		},
		"has higher precedence than logical and": {
			input: "foo && bar | baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewBinaryExpressionNode(
								L(S(P(7, 1, 8), P(15, 1, 16))),
								T(L(S(P(11, 1, 12), P(11, 1, 12))), token.OR),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
								ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(15, 1, 16))), "baz"),
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

func TestBitwiseXor(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 ^ 2 ^ 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.XOR),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.XOR),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 ^\n2 ^\n3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 3, 1))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 3, 1))),
							T(L(S(P(6, 2, 3), P(6, 2, 3))), token.XOR),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 2, 1))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.XOR),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 2, 1), P(4, 2, 1))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 3, 1), P(8, 3, 1))), "3"),
						),
					),
				},
			),
		},
		"has higher precedence than bitwise or": {
			input: "foo | bar ^ baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							T(L(S(P(4, 1, 5), P(4, 1, 5))), token.OR),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewBinaryExpressionNode(
								L(S(P(6, 1, 7), P(14, 1, 15))),
								T(L(S(P(10, 1, 11), P(10, 1, 11))), token.XOR),
								ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
								ast.NewPublicIdentifierNode(L(S(P(12, 1, 13), P(14, 1, 15))), "baz"),
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

func TestBitwiseAnd(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 & 2 & 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.AND),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.AND),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 &\n2 &\n3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 3, 1))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(8, 3, 1))),
							T(L(S(P(6, 2, 3), P(6, 2, 3))), token.AND),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(4, 2, 1))),
								T(L(S(P(2, 1, 3), P(2, 1, 3))), token.AND),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(4, 2, 1), P(4, 2, 1))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(8, 3, 1), P(8, 3, 1))), "3"),
						),
					),
				},
			),
		},
		"has higher precedence than bitwise xor": {
			input: "foo ^ bar & baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							T(L(S(P(4, 1, 5), P(4, 1, 5))), token.XOR),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewBinaryExpressionNode(
								L(S(P(6, 1, 7), P(14, 1, 15))),
								T(L(S(P(10, 1, 11), P(10, 1, 11))), token.AND),
								ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
								ast.NewPublicIdentifierNode(L(S(P(12, 1, 13), P(14, 1, 15))), "baz"),
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

func TestBitwiseShift(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 << 2 >> 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							T(L(S(P(7, 1, 8), P(8, 1, 9))), token.RBITSHIFT),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(5, 1, 6))),
								T(L(S(P(2, 1, 3), P(3, 1, 4))), token.LBITSHIFT),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(5, 1, 6), P(5, 1, 6))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(10, 1, 11), P(10, 1, 11))), "3"),
						),
					),
				},
			),
		},
		"can be triple": {
			input: "1 <<< 2 >>> 3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 1, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(12, 1, 13))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(12, 1, 13))),
							T(L(S(P(8, 1, 9), P(10, 1, 11))), token.RTRIPLE_BITSHIFT),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(6, 1, 7))),
								T(L(S(P(2, 1, 3), P(4, 1, 5))), token.LTRIPLE_BITSHIFT),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(6, 1, 7), P(6, 1, 7))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(12, 1, 13), P(12, 1, 13))), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 <<\n2 >>\n3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 3, 1))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(10, 3, 1))),
							T(L(S(P(7, 2, 3), P(8, 2, 4))), token.RBITSHIFT),
							ast.NewBinaryExpressionNode(
								L(S(P(0, 1, 1), P(5, 2, 1))),
								T(L(S(P(2, 1, 3), P(3, 1, 4))), token.LBITSHIFT),
								ast.NewIntLiteralNode(L(S(P(0, 1, 1), P(0, 1, 1))), "1"),
								ast.NewIntLiteralNode(L(S(P(5, 2, 1), P(5, 2, 1))), "2"),
							),
							ast.NewIntLiteralNode(L(S(P(10, 3, 1), P(10, 3, 1))), "3"),
						),
					),
				},
			),
		},
		"has higher precedence than comparisons": {
			input: "foo > bar << baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewBinaryExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							T(L(S(P(4, 1, 5), P(4, 1, 5))), token.GREATER),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewBinaryExpressionNode(
								L(S(P(6, 1, 7), P(15, 1, 16))),
								T(L(S(P(10, 1, 11), P(11, 1, 12))), token.LBITSHIFT),
								ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
								ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(15, 1, 16))), "baz"),
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
