package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
)

func TestEquality(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "bar == baz == 1",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(11, 2, 1, 12), token.EQUAL_EQUAL),
							ast.NewBinaryExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), token.EQUAL_EQUAL),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "bar"),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "baz"),
							),
							ast.NewIntLiteralNode(P(14, 1, 1, 15), "1"),
						),
					),
				},
			),
		},
		"can have endlines after the operator": {
			input: "bar ==\nbaz ==\n1",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(11, 2, 2, 5), token.EQUAL_EQUAL),
							ast.NewBinaryExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), token.EQUAL_EQUAL),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "bar"),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 1), "baz"),
							),
							ast.NewIntLiteralNode(P(14, 1, 3, 1), "1"),
						),
					),
				},
			),
		},
		"can't have endlines before the operator": {
			input: "bar\n== baz\n== 1",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "bar"),
					),
					ast.NewExpressionStatementNode(
						P(4, 7, 2, 1),
						ast.NewInvalidNode(P(4, 2, 2, 1), T(P(4, 2, 2, 1), token.EQUAL_EQUAL)),
					),
					ast.NewExpressionStatementNode(
						P(11, 4, 3, 1),
						ast.NewInvalidNode(P(11, 2, 3, 1), T(P(11, 2, 3, 1), token.EQUAL_EQUAL)),
					),
				},
			),
			err: position.ErrorList{
				position.NewError(L("main", 4, 2, 2, 1), "unexpected ==, expected an expression"),
				position.NewError(L("main", 11, 2, 3, 1), "unexpected ==, expected an expression"),
			},
		},
		"has many versions": {
			input: "a == b != c === d !== e =:= f =!= g",
			want: ast.NewProgramNode(
				P(0, 35, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 35, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 35, 1, 1),
							T(P(30, 3, 1, 31), token.REF_NOT_EQUAL),
							ast.NewBinaryExpressionNode(
								P(0, 29, 1, 1),
								T(P(24, 3, 1, 25), token.REF_EQUAL),
								ast.NewBinaryExpressionNode(
									P(0, 23, 1, 1),
									T(P(18, 3, 1, 19), token.STRICT_NOT_EQUAL),
									ast.NewBinaryExpressionNode(
										P(0, 17, 1, 1),
										T(P(12, 3, 1, 13), token.STRICT_EQUAL),
										ast.NewBinaryExpressionNode(
											P(0, 11, 1, 1),
											T(P(7, 2, 1, 8), token.NOT_EQUAL),
											ast.NewBinaryExpressionNode(
												P(0, 6, 1, 1),
												T(P(2, 2, 1, 3), token.EQUAL_EQUAL),
												ast.NewPublicIdentifierNode(P(0, 1, 1, 1), "a"),
												ast.NewPublicIdentifierNode(P(5, 1, 1, 6), "b"),
											),
											ast.NewPublicIdentifierNode(P(10, 1, 1, 11), "c"),
										),
										ast.NewPublicIdentifierNode(P(16, 1, 1, 17), "d"),
									),
									ast.NewPublicIdentifierNode(P(22, 1, 1, 23), "e"),
								),
								ast.NewPublicIdentifierNode(P(28, 1, 1, 29), "f"),
							),
							ast.NewPublicIdentifierNode(P(34, 1, 1, 35), "g"),
						),
					),
				},
			),
		},
		"has higher precedence than bitwise and": {
			input: "foo & bar == baz",
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 16, 1, 1),
							T(P(4, 1, 1, 5), token.AND),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(6, 10, 1, 7),
								T(P(10, 2, 1, 11), token.EQUAL_EQUAL),
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

func TestComparison(t *testing.T) {
	tests := testTable{
		"is processed from left to right": {
			input: "foo > bar > baz",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(10, 1, 1, 11), token.GREATER),
							ast.NewBinaryExpressionNode(
								P(0, 9, 1, 1),
								T(P(4, 1, 1, 5), token.GREATER),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
							),
							ast.NewPublicIdentifierNode(P(12, 3, 1, 13), "baz"),
						),
					),
				},
			),
		},
		"can have endlines after the operator": {
			input: "foo >\nbar >\nbaz",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(10, 1, 2, 5), token.GREATER),
							ast.NewBinaryExpressionNode(
								P(0, 9, 1, 1),
								T(P(4, 1, 1, 5), token.GREATER),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(6, 3, 2, 1), "bar"),
							),
							ast.NewPublicIdentifierNode(P(12, 3, 3, 1), "baz"),
						),
					),
				},
			),
		},
		"can't have endlines before the operator": {
			input: "bar\n> baz\n> baz",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "bar"),
					),
					ast.NewExpressionStatementNode(
						P(4, 6, 2, 1),
						ast.NewInvalidNode(P(4, 1, 2, 1), T(P(4, 1, 2, 1), token.GREATER)),
					),
					ast.NewExpressionStatementNode(
						P(10, 5, 3, 1),
						ast.NewInvalidNode(P(10, 1, 3, 1), T(P(10, 1, 3, 1), token.GREATER)),
					),
				},
			),
			err: position.ErrorList{
				position.NewError(L("main", 4, 1, 2, 1), "unexpected >, expected an expression"),
				position.NewError(L("main", 10, 1, 3, 1), "unexpected >, expected an expression"),
			},
		},
		"has many versions": {
			input: "a < b <= c > d >= e <: f :> g <<: h :>> i <=> j",
			want: ast.NewProgramNode(
				P(0, 47, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 47, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 47, 1, 1),
							T(P(42, 3, 1, 43), token.SPACESHIP_OP),
							ast.NewBinaryExpressionNode(
								P(0, 41, 1, 1),
								T(P(36, 3, 1, 37), token.REVERSE_INSTANCE_OF_OP),
								ast.NewBinaryExpressionNode(
									P(0, 35, 1, 1),
									T(P(30, 3, 1, 31), token.INSTANCE_OF_OP),
									ast.NewBinaryExpressionNode(
										P(0, 29, 1, 1),
										T(P(25, 2, 1, 26), token.REVERSE_ISA_OP),
										ast.NewBinaryExpressionNode(
											P(0, 24, 1, 1),
											T(P(20, 2, 1, 21), token.ISA_OP),
											ast.NewBinaryExpressionNode(
												P(0, 19, 1, 1),
												T(P(15, 2, 1, 16), token.GREATER_EQUAL),
												ast.NewBinaryExpressionNode(
													P(0, 14, 1, 1),
													T(P(11, 1, 1, 12), token.GREATER),
													ast.NewBinaryExpressionNode(
														P(0, 10, 1, 1),
														T(P(6, 2, 1, 7), token.LESS_EQUAL),
														ast.NewBinaryExpressionNode(
															P(0, 5, 1, 1),
															T(P(2, 1, 1, 3), token.LESS),
															ast.NewPublicIdentifierNode(P(0, 1, 1, 1), "a"),
															ast.NewPublicIdentifierNode(P(4, 1, 1, 5), "b"),
														),
														ast.NewPublicIdentifierNode(P(9, 1, 1, 10), "c"),
													),
													ast.NewPublicIdentifierNode(P(13, 1, 1, 14), "d"),
												),
												ast.NewPublicIdentifierNode(P(18, 1, 1, 19), "e"),
											),
											ast.NewPublicIdentifierNode(P(23, 1, 1, 24), "f"),
										),
										ast.NewPublicIdentifierNode(P(28, 1, 1, 29), "g"),
									),
									ast.NewPublicIdentifierNode(P(34, 1, 1, 35), "h"),
								),
								ast.NewPublicIdentifierNode(P(40, 1, 1, 41), "i"),
							),
							ast.NewPublicIdentifierNode(P(46, 1, 1, 47), "j"),
						),
					),
				},
			),
		},
		"has higher precedence than equality operators": {
			input: "foo == bar >= baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 17, 1, 1),
							T(P(4, 2, 1, 5), token.EQUAL_EQUAL),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(7, 10, 1, 8),
								T(P(11, 2, 1, 12), token.GREATER_EQUAL),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
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
