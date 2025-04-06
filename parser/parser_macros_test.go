package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
)

func TestMacroBoundary(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
do macro
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(28, 5, 4)),
						ast.NewMacroBoundaryNode(
							S(P(1, 2, 1), P(27, 5, 3)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(11, 3, 2), P(19, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(11, 3, 2), P(18, 3, 9)),
										T(S(P(15, 3, 6), P(16, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(11, 3, 2), P(13, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(18, 3, 9), P(18, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(21, 4, 2), P(24, 4, 5)),
									ast.NewNilLiteralNode(S(P(21, 4, 2), P(23, 4, 4))),
								),
							},
							"",
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
				do macro
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 3, 8)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(21, 3, 8)),
						ast.NewMacroBoundaryNode(
							S(P(5, 2, 5), P(20, 3, 7)),
							nil,
							"",
						),
					),
				},
			),
		},
		"can have a name": {
			input: `
				do macro 'foo'
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 3, 8)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(27, 3, 8)),
						ast.NewMacroBoundaryNode(
							S(P(5, 2, 5), P(26, 3, 7)),
							nil,
							"foo",
						),
					),
				},
			),
		},
		"can be a one-liner": {
			input: "do macro 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewMacroBoundaryNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(9, 1, 10), P(9, 1, 10)),
									ast.NewIntLiteralNode(
										S(P(9, 1, 10), P(9, 1, 10)),
										"5",
									),
								),
							},
							"",
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
				bar =
					do macro
						foo += 2
					end
				nil
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(56, 6, 8)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(48, 5, 9)),
						ast.NewAssignmentExpressionNode(
							S(P(5, 2, 5), P(47, 5, 8)),
							T(S(P(9, 2, 9), P(9, 2, 9)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(5, 2, 5), P(7, 2, 7)), "bar"),
							ast.NewMacroBoundaryNode(
								S(P(16, 3, 6), P(47, 5, 8)),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(31, 4, 7), P(39, 4, 15)),
										ast.NewAssignmentExpressionNode(
											S(P(31, 4, 7), P(38, 4, 14)),
											T(S(P(35, 4, 11), P(36, 4, 12)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(31, 4, 7), P(33, 4, 9)), "foo"),
											ast.NewIntLiteralNode(S(P(38, 4, 14), P(38, 4, 14)), "2"),
										),
									),
								},
								"",
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(53, 6, 5), P(56, 6, 8)),
						ast.NewNilLiteralNode(S(P(53, 6, 5), P(55, 6, 7))),
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

func TestQuoteExpression(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
quote
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(25, 5, 4)),
						ast.NewQuoteExpressionNode(
							S(P(1, 2, 1), P(24, 5, 3)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(8, 3, 2), P(16, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(8, 3, 2), P(15, 3, 9)),
										T(S(P(12, 3, 6), P(13, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(8, 3, 2), P(10, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(15, 3, 9), P(15, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(18, 4, 2), P(21, 4, 5)),
									ast.NewNilLiteralNode(S(P(18, 4, 2), P(20, 4, 4))),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
				quote
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 3, 8)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(18, 3, 8)),
						ast.NewQuoteExpressionNode(
							S(P(5, 2, 5), P(17, 3, 7)),
							nil,
						),
					),
				},
			),
		},
		"can be a one-liner": {
			input: "quote 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewQuoteExpressionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(6, 1, 7), P(6, 1, 7)),
									ast.NewIntLiteralNode(
										S(P(6, 1, 7), P(6, 1, 7)),
										"5",
									),
								),
							},
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
				bar =
					quote
						foo += 2
					end
				nil
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(53, 6, 8)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(45, 5, 9)),
						ast.NewAssignmentExpressionNode(
							S(P(5, 2, 5), P(44, 5, 8)),
							T(S(P(9, 2, 9), P(9, 2, 9)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(5, 2, 5), P(7, 2, 7)), "bar"),
							ast.NewQuoteExpressionNode(
								S(P(16, 3, 6), P(44, 5, 8)),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(28, 4, 7), P(36, 4, 15)),
										ast.NewAssignmentExpressionNode(
											S(P(28, 4, 7), P(35, 4, 14)),
											T(S(P(32, 4, 11), P(33, 4, 12)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(28, 4, 7), P(30, 4, 9)), "foo"),
											ast.NewIntLiteralNode(S(P(35, 4, 14), P(35, 4, 14)), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(50, 6, 5), P(53, 6, 8)),
						ast.NewNilLiteralNode(S(P(50, 6, 5), P(52, 6, 7))),
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
