package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
)

func TestAddition(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 + 2 + 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), token.PLUS),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.PLUS),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 +\n2 +\n3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 2, 3), token.PLUS),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.PLUS),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 2, 1), V(P(4, 1, 2, 1), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 3, 1), V(P(8, 1, 3, 1), token.INT, "3")),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n+ 2\n+ 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
					),
					ast.NewExpressionStatementNode(
						P(2, 4, 2, 1),
						ast.NewUnaryExpressionNode(
							P(2, 3, 2, 1),
							T(P(2, 1, 2, 1), token.PLUS),
							ast.NewIntLiteralNode(P(4, 1, 2, 3), V(P(4, 1, 2, 3), token.INT, "2")),
						),
					),
					ast.NewExpressionStatementNode(
						P(6, 3, 3, 1),
						ast.NewUnaryExpressionNode(
							P(6, 3, 3, 1),
							T(P(6, 1, 3, 1), token.PLUS),
							ast.NewIntLiteralNode(P(8, 1, 3, 3), V(P(8, 1, 3, 3), token.INT, "3")),
						),
					),
				},
			),
		},
		"has higher precedence than bitshifts": {
			input: "foo >> bar + baz",
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 16, 1, 1),
							T(P(4, 2, 1, 5), token.RBITSHIFT),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(7, 9, 1, 8),
								T(P(11, 1, 1, 12), token.PLUS),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(13, 3, 1, 14), "baz"),
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
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), token.MINUS),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.MINUS),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 -\n2 -\n3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 2, 3), token.MINUS),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.MINUS),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 2, 1), V(P(4, 1, 2, 1), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 3, 1), V(P(8, 1, 3, 1), token.INT, "3")),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n- 2\n- 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
					),
					ast.NewExpressionStatementNode(
						P(2, 4, 2, 1),
						ast.NewUnaryExpressionNode(
							P(2, 3, 2, 1),
							T(P(2, 1, 2, 1), token.MINUS),
							ast.NewIntLiteralNode(P(4, 1, 2, 3), V(P(4, 1, 2, 3), token.INT, "2")),
						),
					),
					ast.NewExpressionStatementNode(
						P(6, 3, 3, 1),
						ast.NewUnaryExpressionNode(
							P(6, 3, 3, 1),
							T(P(6, 1, 3, 1), token.MINUS),
							ast.NewIntLiteralNode(P(8, 1, 3, 3), V(P(8, 1, 3, 3), token.INT, "3")),
						),
					),
				},
			),
		},
		"has the same precedence as addition": {
			input: "1 + 2 - 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), token.MINUS),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.PLUS),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "3")),
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
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), token.STAR),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.STAR),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 *\n2 *\n3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 2, 3), token.STAR),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.STAR),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 2, 1), V(P(4, 1, 2, 1), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 3, 1), V(P(8, 1, 3, 1), token.INT, "3")),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n* 2\n* 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
					),
					ast.NewExpressionStatementNode(
						P(2, 4, 2, 1),
						ast.NewInvalidNode(P(2, 1, 2, 1), T(P(2, 1, 2, 1), token.STAR)),
					),
					ast.NewExpressionStatementNode(
						P(6, 3, 3, 1),
						ast.NewInvalidNode(P(6, 1, 3, 1), T(P(6, 1, 3, 1), token.STAR)),
					),
				},
			),
			err: ErrorList{
				NewError(P(2, 1, 2, 1), "unexpected *, expected an expression"),
				NewError(P(6, 1, 3, 1), "unexpected *, expected an expression"),
			},
		},
		"has higher precedence than addition": {
			input: "1 + 2 * 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(2, 1, 1, 3), token.PLUS),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
							ast.NewBinaryExpressionNode(
								P(4, 5, 1, 5),
								T(P(6, 1, 1, 7), token.STAR),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.INT, "2")),
								ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "3")),
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
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), token.SLASH),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.SLASH),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 /\n2 /\n3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 2, 3), token.SLASH),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.SLASH),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 2, 1), V(P(4, 1, 2, 1), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 3, 1), V(P(8, 1, 3, 1), token.INT, "3")),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n/ 2\n/ 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
					),
					ast.NewExpressionStatementNode(
						P(2, 4, 2, 1),
						ast.NewInvalidNode(P(2, 1, 2, 1), T(P(2, 1, 2, 1), token.SLASH)),
					),
					ast.NewExpressionStatementNode(
						P(6, 3, 3, 1),
						ast.NewInvalidNode(P(6, 1, 3, 1), T(P(6, 1, 3, 1), token.SLASH)),
					),
				},
			),
			err: ErrorList{
				NewError(P(2, 1, 2, 1), "unexpected /, expected an expression"),
				NewError(P(6, 1, 3, 1), "unexpected /, expected an expression"),
			},
		},
		"has the same precedence as multiplication": {
			input: "1 * 2 / 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), token.SLASH),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.STAR),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "3")),
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
		"plus can be nested": {
			input: "+++1.5",
			want: ast.NewProgramNode(
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(0, 1, 1, 1), token.PLUS),
							ast.NewUnaryExpressionNode(
								P(1, 5, 1, 2),
								T(P(1, 1, 1, 2), token.PLUS),
								ast.NewUnaryExpressionNode(
									P(2, 4, 1, 3),
									T(P(2, 1, 1, 3), token.PLUS),
									ast.NewFloatLiteralNode(P(3, 3, 1, 4), "1.5"),
								),
							),
						),
					),
				},
			),
		},
		"minus can be nested": {
			input: "---1.5",
			want: ast.NewProgramNode(
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(0, 1, 1, 1), token.MINUS),
							ast.NewUnaryExpressionNode(
								P(1, 5, 1, 2),
								T(P(1, 1, 1, 2), token.MINUS),
								ast.NewUnaryExpressionNode(
									P(2, 4, 1, 3),
									T(P(2, 1, 1, 3), token.MINUS),
									ast.NewFloatLiteralNode(P(3, 3, 1, 4), "1.5"),
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
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(0, 1, 1, 1), token.BANG),
							ast.NewUnaryExpressionNode(
								P(1, 5, 1, 2),
								T(P(1, 1, 1, 2), token.BANG),
								ast.NewUnaryExpressionNode(
									P(2, 4, 1, 3),
									T(P(2, 1, 1, 3), token.BANG),
									ast.NewFloatLiteralNode(P(3, 3, 1, 4), "1.5"),
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
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(0, 1, 1, 1), token.TILDE),
							ast.NewUnaryExpressionNode(
								P(1, 5, 1, 2),
								T(P(1, 1, 1, 2), token.TILDE),
								ast.NewUnaryExpressionNode(
									P(2, 4, 1, 3),
									T(P(2, 1, 1, 3), token.TILDE),
									ast.NewFloatLiteralNode(P(3, 3, 1, 4), "1.5"),
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
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(0, 1, 1, 1), token.BANG),
							ast.NewUnaryExpressionNode(
								P(1, 5, 1, 2),
								T(P(1, 1, 1, 2), token.PLUS),
								ast.NewUnaryExpressionNode(
									P(2, 4, 1, 3),
									T(P(2, 1, 1, 3), token.TILDE),
									ast.NewFloatLiteralNode(P(3, 3, 1, 4), "1.5"),
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
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(10, 1, 1, 11), token.PLUS),
							ast.NewBinaryExpressionNode(
								P(0, 9, 1, 1),
								T(P(6, 1, 1, 7), token.STAR),
								ast.NewUnaryExpressionNode(
									P(0, 5, 1, 1),
									T(P(0, 1, 1, 1), token.BANG),
									ast.NewUnaryExpressionNode(
										P(1, 4, 1, 2),
										T(P(1, 1, 1, 2), token.BANG),
										ast.NewFloatLiteralNode(P(2, 3, 1, 3), "1.5"),
									),
								),
								ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "2")),
							),
							ast.NewUnaryExpressionNode(
								P(12, 3, 1, 13),
								T(P(12, 1, 1, 13), token.TILDE),
								ast.NewFloatLiteralNode(P(13, 2, 1, 14), "0.5"),
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
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 11, 1, 1),
							T(P(2, 2, 1, 3), token.STAR_STAR),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
							ast.NewBinaryExpressionNode(
								P(5, 6, 1, 6),
								T(P(7, 2, 1, 8), token.STAR_STAR),
								ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), token.INT, "2")),
								ast.NewIntLiteralNode(P(10, 1, 1, 11), V(P(10, 1, 1, 11), token.INT, "3")),
							),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 **\n2 **\n3",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 11, 1, 1),
							T(P(2, 2, 1, 3), token.STAR_STAR),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
							ast.NewBinaryExpressionNode(
								P(5, 6, 2, 1),
								T(P(7, 2, 2, 3), token.STAR_STAR),
								ast.NewIntLiteralNode(P(5, 1, 2, 1), V(P(5, 1, 2, 1), token.INT, "2")),
								ast.NewIntLiteralNode(P(10, 1, 3, 1), V(P(10, 1, 3, 1), token.INT, "3")),
							),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n** 2\n** 3",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
					),
					ast.NewExpressionStatementNode(
						P(2, 5, 2, 1),
						ast.NewInvalidNode(P(2, 2, 2, 1), T(P(2, 2, 2, 1), token.STAR_STAR)),
					),
					ast.NewExpressionStatementNode(
						P(7, 4, 3, 1),
						ast.NewInvalidNode(P(7, 2, 3, 1), T(P(7, 2, 3, 1), token.STAR_STAR)),
					),
				},
			),
			err: ErrorList{
				NewError(P(2, 2, 2, 1), "unexpected **, expected an expression"),
				NewError(P(7, 2, 3, 1), "unexpected **, expected an expression"),
			},
		},
		"has higher precedence than unary expressions": {
			input: "-2 ** 3",
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 7, 1, 1),
							T(P(0, 1, 1, 1), token.MINUS),
							ast.NewBinaryExpressionNode(
								P(1, 6, 1, 2),
								T(P(3, 2, 1, 4), token.STAR_STAR),
								ast.NewIntLiteralNode(P(1, 1, 1, 2), V(P(1, 1, 1, 2), token.INT, "2")),
								ast.NewIntLiteralNode(P(6, 1, 1, 7), V(P(6, 1, 1, 7), token.INT, "3")),
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
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), token.OR),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.OR),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 |\n2 |\n3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 2, 3), token.OR),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.OR),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 2, 1), V(P(4, 1, 2, 1), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 3, 1), V(P(8, 1, 3, 1), token.INT, "3")),
						),
					),
				},
			),
		},
		"has higher precedence than logical and": {
			input: "foo && bar | baz",
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 16, 1, 1),
							T(P(4, 2, 1, 5), token.AND_AND),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(7, 9, 1, 8),
								T(P(11, 1, 1, 12), token.OR),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(13, 3, 1, 14), "baz"),
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
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), token.XOR),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.XOR),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 ^\n2 ^\n3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 2, 3), token.XOR),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.XOR),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 2, 1), V(P(4, 1, 2, 1), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 3, 1), V(P(8, 1, 3, 1), token.INT, "3")),
						),
					),
				},
			),
		},
		"has higher precedence than bitwise or": {
			input: "foo | bar ^ baz",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(4, 1, 1, 5), token.OR),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(6, 9, 1, 7),
								T(P(10, 1, 1, 11), token.XOR),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
								ast.NewPublicIdentifierNode(P(12, 3, 1, 13), "baz"),
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
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), token.AND),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.AND),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), token.INT, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 &\n2 &\n3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 2, 3), token.AND),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), token.AND),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(4, 1, 2, 1), V(P(4, 1, 2, 1), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 3, 1), V(P(8, 1, 3, 1), token.INT, "3")),
						),
					),
				},
			),
		},
		"has higher precedence than bitwise xor": {
			input: "foo ^ bar & baz",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(4, 1, 1, 5), token.XOR),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(6, 9, 1, 7),
								T(P(10, 1, 1, 11), token.AND),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
								ast.NewPublicIdentifierNode(P(12, 3, 1, 13), "baz"),
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
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 11, 1, 1),
							T(P(7, 2, 1, 8), token.RBITSHIFT),
							ast.NewBinaryExpressionNode(
								P(0, 6, 1, 1),
								T(P(2, 2, 1, 3), token.LBITSHIFT),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(10, 1, 1, 11), V(P(10, 1, 1, 11), token.INT, "3")),
						),
					),
				},
			),
		},
		"can be triple": {
			input: "1 <<< 2 >>> 3",
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 13, 1, 1),
							T(P(8, 3, 1, 9), token.RTRIPLE_BITSHIFT),
							ast.NewBinaryExpressionNode(
								P(0, 7, 1, 1),
								T(P(2, 3, 1, 3), token.LTRIPLE_BITSHIFT),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(6, 1, 1, 7), V(P(6, 1, 1, 7), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(12, 1, 1, 13), V(P(12, 1, 1, 13), token.INT, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 <<\n2 >>\n3",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 11, 1, 1),
							T(P(7, 2, 2, 3), token.RBITSHIFT),
							ast.NewBinaryExpressionNode(
								P(0, 6, 1, 1),
								T(P(2, 2, 1, 3), token.LBITSHIFT),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.INT, "1")),
								ast.NewIntLiteralNode(P(5, 1, 2, 1), V(P(5, 1, 2, 1), token.INT, "2")),
							),
							ast.NewIntLiteralNode(P(10, 1, 3, 1), V(P(10, 1, 3, 1), token.INT, "3")),
						),
					),
				},
			),
		},
		"has higher precedence than comparisons": {
			input: "foo > bar << baz",
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 16, 1, 1),
							T(P(4, 1, 1, 5), token.GREATER),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(6, 10, 1, 7),
								T(P(10, 2, 1, 11), token.LBITSHIFT),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
								ast.NewPublicIdentifierNode(P(13, 3, 1, 14), "baz"),
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
