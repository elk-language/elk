package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/token"
)

func TestEquality(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "bar == baz == 1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							T(S(P(11, 1, 12), P(12, 1, 13)), token.EQUAL_EQUAL),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(9, 1, 10)),
								T(S(P(4, 1, 5), P(5, 1, 6)), token.EQUAL_EQUAL),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "bar"),
								ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "baz"),
							),
							ast.NewIntLiteralNode(S(P(14, 1, 15), P(14, 1, 15)), "1"),
						),
					),
				},
			),
		},
		"can have endlines after the operator": {
			input: "bar ==\nbaz ==\n1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 3, 1)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(14, 3, 1)),
							T(S(P(11, 2, 5), P(12, 2, 6)), token.EQUAL_EQUAL),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(9, 2, 3)),
								T(S(P(4, 1, 5), P(5, 1, 6)), token.EQUAL_EQUAL),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "bar"),
								ast.NewPublicIdentifierNode(S(P(7, 2, 1), P(9, 2, 3)), "baz"),
							),
							ast.NewIntLiteralNode(S(P(14, 3, 1), P(14, 3, 1)), "1"),
						),
					),
				},
			),
		},
		"cannot have endlines before the operator": {
			input: "bar\n== baz\n== 1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 3, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "bar"),
					),
					ast.NewExpressionStatementNode(
						S(P(4, 2, 1), P(10, 2, 7)),
						ast.NewInvalidNode(S(P(4, 2, 1), P(5, 2, 2)), T(S(P(4, 2, 1), P(5, 2, 2)), token.EQUAL_EQUAL)),
					),
					ast.NewExpressionStatementNode(
						S(P(11, 3, 1), P(14, 3, 4)),
						ast.NewInvalidNode(S(P(11, 3, 1), P(12, 3, 2)), T(S(P(11, 3, 1), P(12, 3, 2)), token.EQUAL_EQUAL)),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(4, 2, 1), P(5, 2, 2)), "unexpected ==, expected an expression"),
				error.NewError(L("<main>", P(11, 3, 1), P(12, 3, 2)), "unexpected ==, expected an expression"),
			},
		},
		"has many versions": {
			input: "a == b != c === d !== e =~ f !~ g",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							T(S(P(29, 1, 30), P(30, 1, 31)), token.LAX_NOT_EQUAL),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(27, 1, 28)),
								T(S(P(24, 1, 25), P(25, 1, 26)), token.LAX_EQUAL),
								ast.NewBinaryExpressionNode(
									S(P(0, 1, 1), P(22, 1, 23)),
									T(S(P(18, 1, 19), P(20, 1, 21)), token.STRICT_NOT_EQUAL),
									ast.NewBinaryExpressionNode(
										S(P(0, 1, 1), P(16, 1, 17)),
										T(S(P(12, 1, 13), P(14, 1, 15)), token.STRICT_EQUAL),
										ast.NewBinaryExpressionNode(
											S(P(0, 1, 1), P(10, 1, 11)),
											T(S(P(7, 1, 8), P(8, 1, 9)), token.NOT_EQUAL),
											ast.NewBinaryExpressionNode(
												S(P(0, 1, 1), P(5, 1, 6)),
												T(S(P(2, 1, 3), P(3, 1, 4)), token.EQUAL_EQUAL),
												ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
												ast.NewPublicIdentifierNode(S(P(5, 1, 6), P(5, 1, 6)), "b"),
											),
											ast.NewPublicIdentifierNode(S(P(10, 1, 11), P(10, 1, 11)), "c"),
										),
										ast.NewPublicIdentifierNode(S(P(16, 1, 17), P(16, 1, 17)), "d"),
									),
									ast.NewPublicIdentifierNode(S(P(22, 1, 23), P(22, 1, 23)), "e"),
								),
								ast.NewPublicIdentifierNode(S(P(27, 1, 28), P(27, 1, 28)), "f"),
							),
							ast.NewPublicIdentifierNode(S(P(32, 1, 33), P(32, 1, 33)), "g"),
						),
					),
				},
			),
		},
		"has higher precedence than bitwise and": {
			input: "foo & bar == baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.AND),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewBinaryExpressionNode(
								S(P(6, 1, 7), P(15, 1, 16)),
								T(S(P(10, 1, 11), P(11, 1, 12)), token.EQUAL_EQUAL),
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

func TestComparison(t *testing.T) {
	tests := testTable{
		"is processed from left to right": {
			input: "foo > bar > baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							T(S(P(10, 1, 11), P(10, 1, 11)), token.GREATER),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(8, 1, 9)),
								T(S(P(4, 1, 5), P(4, 1, 5)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
							),
							ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(14, 1, 15)), "baz"),
						),
					),
				},
			),
		},
		"can have endlines after the operator": {
			input: "foo >\nbar >\nbaz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 3, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 3, 3)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(14, 3, 3)),
							T(S(P(10, 2, 5), P(10, 2, 5)), token.GREATER),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(8, 2, 3)),
								T(S(P(4, 1, 5), P(4, 1, 5)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								ast.NewPublicIdentifierNode(S(P(6, 2, 1), P(8, 2, 3)), "bar"),
							),
							ast.NewPublicIdentifierNode(S(P(12, 3, 1), P(14, 3, 3)), "baz"),
						),
					),
				},
			),
		},
		"cannot have endlines before the operator": {
			input: "bar\n> baz\n> baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 3, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "bar"),
					),
					ast.NewExpressionStatementNode(
						S(P(4, 2, 1), P(9, 2, 6)),
						ast.NewInvalidNode(S(P(4, 2, 1), P(4, 2, 1)), T(S(P(4, 2, 1), P(4, 2, 1)), token.GREATER)),
					),
					ast.NewExpressionStatementNode(
						S(P(10, 3, 1), P(14, 3, 5)),
						ast.NewInvalidNode(S(P(10, 3, 1), P(10, 3, 1)), T(S(P(10, 3, 1), P(10, 3, 1)), token.GREATER)),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(4, 2, 1), P(4, 2, 1)), "unexpected >, expected an expression"),
				error.NewError(L("<main>", P(10, 3, 1), P(10, 3, 1)), "unexpected >, expected an expression"),
			},
		},
		"has many versions": {
			input: "a < b <= c > d >= e <: f :> g <<: h :>> i <=> j",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(46, 1, 47)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(46, 1, 47)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(46, 1, 47)),
							T(S(P(42, 1, 43), P(44, 1, 45)), token.SPACESHIP_OP),
							ast.NewBinaryExpressionNode(
								S(P(0, 1, 1), P(40, 1, 41)),
								T(S(P(36, 1, 37), P(38, 1, 39)), token.REVERSE_INSTANCE_OF_OP),
								ast.NewBinaryExpressionNode(
									S(P(0, 1, 1), P(34, 1, 35)),
									T(S(P(30, 1, 31), P(32, 1, 33)), token.INSTANCE_OF_OP),
									ast.NewBinaryExpressionNode(
										S(P(0, 1, 1), P(28, 1, 29)),
										T(S(P(25, 1, 26), P(26, 1, 27)), token.REVERSE_ISA_OP),
										ast.NewBinaryExpressionNode(
											S(P(0, 1, 1), P(23, 1, 24)),
											T(S(P(20, 1, 21), P(21, 1, 22)), token.ISA_OP),
											ast.NewBinaryExpressionNode(
												S(P(0, 1, 1), P(18, 1, 19)),
												T(S(P(15, 1, 16), P(16, 1, 17)), token.GREATER_EQUAL),
												ast.NewBinaryExpressionNode(
													S(P(0, 1, 1), P(13, 1, 14)),
													T(S(P(11, 1, 12), P(11, 1, 12)), token.GREATER),
													ast.NewBinaryExpressionNode(
														S(P(0, 1, 1), P(9, 1, 10)),
														T(S(P(6, 1, 7), P(7, 1, 8)), token.LESS_EQUAL),
														ast.NewBinaryExpressionNode(
															S(P(0, 1, 1), P(4, 1, 5)),
															T(S(P(2, 1, 3), P(2, 1, 3)), token.LESS),
															ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
															ast.NewPublicIdentifierNode(S(P(4, 1, 5), P(4, 1, 5)), "b"),
														),
														ast.NewPublicIdentifierNode(S(P(9, 1, 10), P(9, 1, 10)), "c"),
													),
													ast.NewPublicIdentifierNode(S(P(13, 1, 14), P(13, 1, 14)), "d"),
												),
												ast.NewPublicIdentifierNode(S(P(18, 1, 19), P(18, 1, 19)), "e"),
											),
											ast.NewPublicIdentifierNode(S(P(23, 1, 24), P(23, 1, 24)), "f"),
										),
										ast.NewPublicIdentifierNode(S(P(28, 1, 29), P(28, 1, 29)), "g"),
									),
									ast.NewPublicIdentifierNode(S(P(34, 1, 35), P(34, 1, 35)), "h"),
								),
								ast.NewPublicIdentifierNode(S(P(40, 1, 41), P(40, 1, 41)), "i"),
							),
							ast.NewPublicIdentifierNode(S(P(46, 1, 47), P(46, 1, 47)), "j"),
						),
					),
				},
			),
		},
		"has higher precedence than equality operators": {
			input: "foo == bar >= baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewBinaryExpressionNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							T(S(P(4, 1, 5), P(5, 1, 6)), token.EQUAL_EQUAL),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewBinaryExpressionNode(
								S(P(7, 1, 8), P(16, 1, 17)),
								T(S(P(11, 1, 12), P(12, 1, 13)), token.GREATER_EQUAL),
								ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
								ast.NewPublicIdentifierNode(S(P(14, 1, 15), P(16, 1, 17)), "baz"),
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
