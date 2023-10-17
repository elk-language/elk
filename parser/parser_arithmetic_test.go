package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
)

func TestAddition(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 + 2 + 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(6, 1, 7), P(6, 1, 7)), token.PLUS),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 1, 5)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.PLUS),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 +\n2 +\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 3, 1)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 3, 1)),
							T(S(P(6, 2, 3), P(6, 2, 3)), token.PLUS),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 2, 1)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.PLUS),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 2, 1), P(4, 2, 1)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 3, 1), P(8, 3, 1)), "3"),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n+ 2\n+ 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
					),
					ast.NewExpressionStatementNode(
						S(P(2, 2, 1), P(5, 2, 4)),
						ast.NewUnaryExpressionNode(
							S(P(2, 2, 1), P(4, 2, 3)),
							T(S(P(2, 2, 1), P(2, 2, 1)), token.PLUS),
							ast.NewIntLiteralNode(S(P(4, 2, 3), P(4, 2, 3)), "2"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(6, 3, 1), P(8, 3, 3)),
						ast.NewUnaryExpressionNode(
							S(P(6, 3, 1), P(8, 3, 3)),
							T(S(P(6, 3, 1), P(6, 3, 1)), token.PLUS),
							ast.NewIntLiteralNode(S(P(8, 3, 3), P(8, 3, 3)), "3"),
						),
					),
				},
			),
		},
		"has higher precedence than bitshifts": {
			input: "foo >> bar + baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							T(S(P(4, 1, 5), P(5, 1, 6)), token.RBITSHIFT),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewBinaryExpressionNode(
								S(P(7, 1, 8), P(15, 1, 16)),
								T(S(P(11, 1, 12), P(11, 1, 12)), token.PLUS),
								ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
								ast.NewPublicIdentifierNode(S(P(13, 1, 14), P(15, 1, 16)), "baz"),
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
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(6, 1, 7), P(6, 1, 7)), token.MINUS),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 1, 5)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.MINUS),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 -\n2 -\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 3, 1)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 3, 1)),
							T(S(P(6, 2, 3), P(6, 2, 3)), token.MINUS),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 2, 1)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.MINUS),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 2, 1), P(4, 2, 1)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 3, 1), P(8, 3, 1)), "3"),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n- 2\n- 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
					),
					ast.NewExpressionStatementNode(
						S(P(2, 2, 1), P(5, 2, 4)),
						ast.NewUnaryExpressionNode(
							S(P(2, 2, 1), P(4, 2, 3)),
							T(S(P(2, 2, 1), P(2, 2, 1)), token.MINUS),
							ast.NewIntLiteralNode(S(P(4, 2, 3), P(4, 2, 3)), "2"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(6, 3, 1), P(8, 3, 3)),
						ast.NewUnaryExpressionNode(
							S(P(6, 3, 1), P(8, 3, 3)),
							T(S(P(6, 3, 1), P(6, 3, 1)), token.MINUS),
							ast.NewIntLiteralNode(S(P(8, 3, 3), P(8, 3, 3)), "3"),
						),
					),
				},
			),
		},
		"has the same precedence as addition": {
			input: "1 + 2 - 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(6, 1, 7), P(6, 1, 7)), token.MINUS),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 1, 5)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.PLUS),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
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
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(6, 1, 7), P(6, 1, 7)), token.STAR),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 1, 5)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.STAR),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 *\n2 *\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 3, 1)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 3, 1)),
							T(S(P(6, 2, 3), P(6, 2, 3)), token.STAR),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 2, 1)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.STAR),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 2, 1), P(4, 2, 1)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 3, 1), P(8, 3, 1)), "3"),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n* 2\n* 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
					),
					ast.NewExpressionStatementNode(
						S(P(2, 2, 1), P(5, 2, 4)),
						ast.NewInvalidNode(S(P(2, 2, 1), P(2, 2, 1)), T(S(P(2, 2, 1), P(2, 2, 1)), token.STAR)),
					),
					ast.NewExpressionStatementNode(
						S(P(6, 3, 1), P(8, 3, 3)),
						ast.NewInvalidNode(S(P(6, 3, 1), P(6, 3, 1)), T(S(P(6, 3, 1), P(6, 3, 1)), token.STAR)),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(2, 2, 1), P(2, 2, 1)), "unexpected *, expected an expression"),
				errors.NewError(L("main", P(6, 3, 1), P(6, 3, 1)), "unexpected *, expected an expression"),
			},
		},
		"has higher precedence than addition": {
			input: "1 + 2 * 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.PLUS),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
							ast.NewBinaryExpressionNode(
								S(P(4, 1, 5), P(8, 1, 9)),
								T(S(P(6, 1, 7), P(6, 1, 7)), token.STAR),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
								ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
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
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(6, 1, 7), P(6, 1, 7)), token.SLASH),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 1, 5)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.SLASH),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 /\n2 /\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 3, 1)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 3, 1)),
							T(S(P(6, 2, 3), P(6, 2, 3)), token.SLASH),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 2, 1)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.SLASH),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 2, 1), P(4, 2, 1)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 3, 1), P(8, 3, 1)), "3"),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n/ 2\n/ 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
					),
					ast.NewExpressionStatementNode(
						S(P(2, 2, 1), P(5, 2, 4)),
						ast.NewInvalidNode(S(P(2, 2, 1), P(2, 2, 1)), T(S(P(2, 2, 1), P(2, 2, 1)), token.SLASH)),
					),
					ast.NewExpressionStatementNode(
						S(P(6, 3, 1), P(8, 3, 3)),
						ast.NewInvalidNode(S(P(6, 3, 1), P(6, 3, 1)), T(S(P(6, 3, 1), P(6, 3, 1)), token.SLASH)),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(2, 2, 1), P(2, 2, 1)), "unexpected /, expected an expression"),
				errors.NewError(L("main", P(6, 3, 1), P(6, 3, 1)), "unexpected /, expected an expression"),
			},
		},
		"has the same precedence as multiplication": {
			input: "1 * 2 / 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(6, 1, 7), P(6, 1, 7)), token.SLASH),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 1, 5)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.STAR),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
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
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(6, 1, 7), P(6, 1, 7)), token.PERCENT),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 1, 5)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.PERCENT),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 %\n2 %\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 3, 1)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 3, 1)),
							T(S(P(6, 2, 3), P(6, 2, 3)), token.PERCENT),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 2, 1)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.PERCENT),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 2, 1), P(4, 2, 1)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 3, 1), P(8, 3, 1)), "3"),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n% 2\n% 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
					),
					ast.NewExpressionStatementNode(
						S(P(2, 2, 1), P(5, 2, 4)),
						ast.NewInvalidNode(S(P(2, 2, 1), P(2, 2, 1)), T(S(P(2, 2, 1), P(2, 2, 1)), token.PERCENT)),
					),
					ast.NewExpressionStatementNode(
						S(P(6, 3, 1), P(8, 3, 3)),
						ast.NewInvalidNode(S(P(6, 3, 1), P(6, 3, 1)), T(S(P(6, 3, 1), P(6, 3, 1)), token.PERCENT)),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(2, 2, 1), P(2, 2, 1)), "unexpected %, expected an expression"),
				errors.NewError(L("main", P(6, 3, 1), P(6, 3, 1)), "unexpected %, expected an expression"),
			},
		},
		"has the same precedence as multiplication": {
			input: "1 * 2 % 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(6, 1, 7), P(6, 1, 7)), token.PERCENT),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 1, 5)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.STAR),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
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
		"plus can't be nested without spaces": {
			input: "+++1.5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.PLUS),
							ast.NewInvalidNode(S(P(0, 1, 1), P(1, 1, 2)), T(S(P(0, 1, 1), P(1, 1, 2)), token.PLUS_PLUS)),
							ast.NewFloatLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "1.5"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(1, 1, 2)), "unexpected ++, expected an expression"),
			},
		},
		"minus can't be nested without spaces": {
			input: "---1.5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.MINUS),
							ast.NewInvalidNode(S(P(0, 1, 1), P(1, 1, 2)), T(S(P(0, 1, 1), P(1, 1, 2)), token.MINUS_MINUS)),
							ast.NewFloatLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "1.5"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(1, 1, 2)), "unexpected --, expected an expression"),
			},
		},
		"plus can be nested": {
			input: "+ + +1.5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewUnaryExpressionNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							T(S(P(0, 1, 1), P(0, 1, 1)), token.PLUS),
							ast.NewUnaryExpressionNode(
								S(P(2, 1, 3), P(7, 1, 8)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.PLUS),
								ast.NewUnaryExpressionNode(
									S(P(4, 1, 5), P(7, 1, 8)),
									T(S(P(4, 1, 5), P(4, 1, 5)), token.PLUS),
									ast.NewFloatLiteralNode(S(P(5, 1, 6), P(7, 1, 8)), "1.5"),
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
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewUnaryExpressionNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							T(S(P(0, 1, 1), P(0, 1, 1)), token.MINUS),
							ast.NewUnaryExpressionNode(
								S(P(2, 1, 3), P(7, 1, 8)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.MINUS),
								ast.NewUnaryExpressionNode(
									S(P(4, 1, 5), P(7, 1, 8)),
									T(S(P(4, 1, 5), P(4, 1, 5)), token.MINUS),
									ast.NewFloatLiteralNode(S(P(5, 1, 6), P(7, 1, 8)), "1.5"),
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
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewUnaryExpressionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							T(S(P(0, 1, 1), P(0, 1, 1)), token.BANG),
							ast.NewUnaryExpressionNode(
								S(P(1, 1, 2), P(5, 1, 6)),
								T(S(P(1, 1, 2), P(1, 1, 2)), token.BANG),
								ast.NewUnaryExpressionNode(
									S(P(2, 1, 3), P(5, 1, 6)),
									T(S(P(2, 1, 3), P(2, 1, 3)), token.BANG),
									ast.NewFloatLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "1.5"),
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
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewUnaryExpressionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							T(S(P(0, 1, 1), P(0, 1, 1)), token.TILDE),
							ast.NewUnaryExpressionNode(
								S(P(1, 1, 2), P(5, 1, 6)),
								T(S(P(1, 1, 2), P(1, 1, 2)), token.TILDE),
								ast.NewUnaryExpressionNode(
									S(P(2, 1, 3), P(5, 1, 6)),
									T(S(P(2, 1, 3), P(2, 1, 3)), token.TILDE),
									ast.NewFloatLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "1.5"),
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
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewUnaryExpressionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							T(S(P(0, 1, 1), P(0, 1, 1)), token.BANG),
							ast.NewUnaryExpressionNode(
								S(P(1, 1, 2), P(5, 1, 6)),
								T(S(P(1, 1, 2), P(1, 1, 2)), token.PLUS),
								ast.NewUnaryExpressionNode(
									S(P(2, 1, 3), P(5, 1, 6)),
									T(S(P(2, 1, 3), P(2, 1, 3)), token.TILDE),
									ast.NewFloatLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "1.5"),
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
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							T(S(P(10, 1, 11), P(10, 1, 11)), token.PLUS),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(8, 1, 9)),
								T(S(P(6, 1, 7), P(6, 1, 7)), token.STAR),
								ast.NewUnaryExpressionNode(
									S(P(0, 1, 1), P(4, 1, 5)),
									T(S(P(0, 1, 1), P(0, 1, 1)), token.BANG),
									ast.NewUnaryExpressionNode(
										S(P(1, 1, 2), P(4, 1, 5)),
										T(S(P(1, 1, 2), P(1, 1, 2)), token.BANG),
										ast.NewFloatLiteralNode(S(P(2, 1, 3), P(4, 1, 5)), "1.5"),
									),
								),
								ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "2"),
							),
							ast.NewUnaryExpressionNode(
								S(P(12, 1, 13), P(14, 1, 15)),
								T(S(P(12, 1, 13), P(12, 1, 13)), token.TILDE),
								ast.NewFloatLiteralNode(S(P(13, 1, 14), P(14, 1, 15)), "0.5"),
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
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							T(S(P(2, 1, 3), P(3, 1, 4)), token.STAR_STAR),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
							ast.NewBinaryExpressionNode(
								S(P(5, 1, 6), P(10, 1, 11)),
								T(S(P(7, 1, 8), P(8, 1, 9)), token.STAR_STAR),
								ast.NewIntLiteralNode(S(P(5, 1, 6), P(5, 1, 6)), "2"),
								ast.NewIntLiteralNode(S(P(10, 1, 11), P(10, 1, 11)), "3"),
							),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 **\n2 **\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 3, 1)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(10, 3, 1)),
							T(S(P(2, 1, 3), P(3, 1, 4)), token.STAR_STAR),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
							ast.NewBinaryExpressionNode(
								S(P(5, 2, 1), P(10, 3, 1)),
								T(S(P(7, 2, 3), P(8, 2, 4)), token.STAR_STAR),
								ast.NewIntLiteralNode(S(P(5, 2, 1), P(5, 2, 1)), "2"),
								ast.NewIntLiteralNode(S(P(10, 3, 1), P(10, 3, 1)), "3"),
							),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n** 2\n** 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 3, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
					),
					ast.NewExpressionStatementNode(
						S(P(2, 2, 1), P(6, 2, 5)),
						ast.NewInvalidNode(S(P(2, 2, 1), P(3, 2, 2)), T(S(P(2, 2, 1), P(3, 2, 2)), token.STAR_STAR)),
					),
					ast.NewExpressionStatementNode(
						S(P(7, 3, 1), P(10, 3, 4)),
						ast.NewInvalidNode(S(P(7, 3, 1), P(8, 3, 2)), T(S(P(7, 3, 1), P(8, 3, 2)), token.STAR_STAR)),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(2, 2, 1), P(3, 2, 2)), "unexpected **, expected an expression"),
				errors.NewError(L("main", P(7, 3, 1), P(8, 3, 2)), "unexpected **, expected an expression"),
			},
		},
		"has higher precedence than unary expressions": {
			input: "-2 ** 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewUnaryExpressionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							T(S(P(0, 1, 1), P(0, 1, 1)), token.MINUS),
							ast.NewBinaryExpressionNode(
								S(P(1, 1, 2), P(6, 1, 7)),
								T(S(P(3, 1, 4), P(4, 1, 5)), token.STAR_STAR),
								ast.NewIntLiteralNode(S(P(1, 1, 2), P(1, 1, 2)), "2"),
								ast.NewIntLiteralNode(S(P(6, 1, 7), P(6, 1, 7)), "3"),
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
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(6, 1, 7), P(6, 1, 7)), token.OR),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 1, 5)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.OR),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 |\n2 |\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 3, 1)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 3, 1)),
							T(S(P(6, 2, 3), P(6, 2, 3)), token.OR),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 2, 1)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.OR),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 2, 1), P(4, 2, 1)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 3, 1), P(8, 3, 1)), "3"),
						),
					),
				},
			),
		},
		"has higher precedence than logical and": {
			input: "foo && bar | baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewLogicalExpressionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							T(S(P(4, 1, 5), P(5, 1, 6)), token.AND_AND),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewBinaryExpressionNode(
								S(P(7, 1, 8), P(15, 1, 16)),
								T(S(P(11, 1, 12), P(11, 1, 12)), token.OR),
								ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
								ast.NewPublicIdentifierNode(S(P(13, 1, 14), P(15, 1, 16)), "baz"),
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
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(6, 1, 7), P(6, 1, 7)), token.XOR),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 1, 5)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.XOR),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 ^\n2 ^\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 3, 1)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 3, 1)),
							T(S(P(6, 2, 3), P(6, 2, 3)), token.XOR),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 2, 1)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.XOR),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 2, 1), P(4, 2, 1)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 3, 1), P(8, 3, 1)), "3"),
						),
					),
				},
			),
		},
		"has higher precedence than bitwise or": {
			input: "foo | bar ^ baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.OR),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewBinaryExpressionNode(
								S(P(6, 1, 7), P(14, 1, 15)),
								T(S(P(10, 1, 11), P(10, 1, 11)), token.XOR),
								ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
								ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(14, 1, 15)), "baz"),
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
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							T(S(P(6, 1, 7), P(6, 1, 7)), token.AND),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 1, 5)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.AND),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 &\n2 &\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 3, 1)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(8, 3, 1)),
							T(S(P(6, 2, 3), P(6, 2, 3)), token.AND),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(4, 2, 1)),
								T(S(P(2, 1, 3), P(2, 1, 3)), token.AND),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(4, 2, 1), P(4, 2, 1)), "2"),
							),
							ast.NewIntLiteralNode(S(P(8, 3, 1), P(8, 3, 1)), "3"),
						),
					),
				},
			),
		},
		"has higher precedence than bitwise xor": {
			input: "foo ^ bar & baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.XOR),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewBinaryExpressionNode(
								S(P(6, 1, 7), P(14, 1, 15)),
								T(S(P(10, 1, 11), P(10, 1, 11)), token.AND),
								ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
								ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(14, 1, 15)), "baz"),
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
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							T(S(P(7, 1, 8), P(8, 1, 9)), token.RBITSHIFT),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(5, 1, 6)),
								T(S(P(2, 1, 3), P(3, 1, 4)), token.LBITSHIFT),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(5, 1, 6), P(5, 1, 6)), "2"),
							),
							ast.NewIntLiteralNode(S(P(10, 1, 11), P(10, 1, 11)), "3"),
						),
					),
				},
			),
		},
		"can be triple": {
			input: "1 <<< 2 >>> 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							T(S(P(8, 1, 9), P(10, 1, 11)), token.RTRIPLE_BITSHIFT),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(6, 1, 7)),
								T(S(P(2, 1, 3), P(4, 1, 5)), token.LTRIPLE_BITSHIFT),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(6, 1, 7), P(6, 1, 7)), "2"),
							),
							ast.NewIntLiteralNode(S(P(12, 1, 13), P(12, 1, 13)), "3"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 <<\n2 >>\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 3, 1)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(10, 3, 1)),
							T(S(P(7, 2, 3), P(8, 2, 4)), token.RBITSHIFT),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(5, 2, 1)),
								T(S(P(2, 1, 3), P(3, 1, 4)), token.LBITSHIFT),
								ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "1"),
								ast.NewIntLiteralNode(S(P(5, 2, 1), P(5, 2, 1)), "2"),
							),
							ast.NewIntLiteralNode(S(P(10, 3, 1), P(10, 3, 1)), "3"),
						),
					),
				},
			),
		},
		"has higher precedence than comparisons": {
			input: "foo > bar << baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.GREATER),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewBinaryExpressionNode(
								S(P(6, 1, 7), P(15, 1, 16)),
								T(S(P(10, 1, 11), P(11, 1, 12)), token.LBITSHIFT),
								ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
								ast.NewPublicIdentifierNode(S(P(13, 1, 14), P(15, 1, 16)), "baz"),
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
