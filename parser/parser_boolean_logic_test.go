package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
)

func TestBooleanLogic(t *testing.T) {
	tests := testTable{
		"has lower precedence than equality": {
			input: "foo && bar == baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewBinaryExpressionNode(
								L(S(P(7, 1, 8), P(16, 1, 17))),
								T(L(S(P(11, 1, 12), P(12, 1, 13))), token.EQUAL_EQUAL),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
								ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
							),
						),
					),
				},
			),
		},
		"or has lower precedence than and": {
			input: "foo || bar && baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.OR_OR),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewLogicalExpressionNode(
								L(S(P(7, 1, 8), P(16, 1, 17))),
								T(L(S(P(11, 1, 12), P(12, 1, 13))), token.AND_AND),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
								ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
							),
						),
					),
				},
			),
		},
		"nil coalescing operator has lower precedence than and": {
			input: "foo ?? bar && baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.QUESTION_QUESTION),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewLogicalExpressionNode(
								L(S(P(7, 1, 8), P(16, 1, 17))),
								T(L(S(P(11, 1, 12), P(12, 1, 13))), token.AND_AND),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
								ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
							),
						),
					),
				},
			),
		},
		"or expression sequencing operator has lower precedence than and": {
			input: "foo |! bar && baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.OR_BANG),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewLogicalExpressionNode(
								L(S(P(7, 1, 8), P(16, 1, 17))),
								T(L(S(P(11, 1, 12), P(12, 1, 13))), token.AND_AND),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
								ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
							),
						),
					),
				},
			),
		},
		"and expression sequencing operator has the same precedence as and": {
			input: "foo &! bar && baz &! boo",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(23, 1, 24))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(23, 1, 24))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(23, 1, 24))),
							T(L(S(P(18, 1, 19), P(19, 1, 20))), token.AND_BANG),
							ast.NewLogicalExpressionNode(
								L(S(P(0, 1, 1), P(16, 1, 17))),
								T(L(S(P(11, 1, 12), P(12, 1, 13))), token.AND_AND),
								ast.NewLogicalExpressionNode(
									L(S(P(0, 1, 1), P(9, 1, 10))),
									T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_BANG),
									ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
									ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
								),
								ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
							),
							ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(23, 1, 24))), "boo"),
						),
					),
				},
			),
		},
		"nil coalescing operator has the same precedence as or": {
			input: "foo ?? bar || baz ?? boo",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(23, 1, 24))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(23, 1, 24))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(23, 1, 24))),
							T(L(S(P(18, 1, 19), P(19, 1, 20))), token.QUESTION_QUESTION),
							ast.NewLogicalExpressionNode(
								L(S(P(0, 1, 1), P(16, 1, 17))),
								T(L(S(P(11, 1, 12), P(12, 1, 13))), token.OR_OR),
								ast.NewLogicalExpressionNode(
									L(S(P(0, 1, 1), P(9, 1, 10))),
									T(L(S(P(4, 1, 5), P(5, 1, 6))), token.QUESTION_QUESTION),
									ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
									ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
								),
								ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
							),
							ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(23, 1, 24))), "boo"),
						),
					),
				},
			),
		},
		"or expression sequencing operator has the same precedence as or": {
			input: "foo |! bar || baz |! boo",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(23, 1, 24))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(23, 1, 24))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(23, 1, 24))),
							T(L(S(P(18, 1, 19), P(19, 1, 20))), token.OR_BANG),
							ast.NewLogicalExpressionNode(
								L(S(P(0, 1, 1), P(16, 1, 17))),
								T(L(S(P(11, 1, 12), P(12, 1, 13))), token.OR_OR),
								ast.NewLogicalExpressionNode(
									L(S(P(0, 1, 1), P(9, 1, 10))),
									T(L(S(P(4, 1, 5), P(5, 1, 6))), token.OR_BANG),
									ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
									ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
								),
								ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
							),
							ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(23, 1, 24))), "boo"),
						),
					),
				},
			),
		},
		"or is evaluated from left to right": {
			input: "foo || bar || baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							T(L(S(P(11, 1, 12), P(12, 1, 13))), token.OR_OR),
							ast.NewLogicalExpressionNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								T(L(S(P(4, 1, 5), P(5, 1, 6))), token.OR_OR),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
						),
					),
				},
			),
		},
		"and is evaluated from left to right": {
			input: "foo && bar && baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							T(L(S(P(11, 1, 12), P(12, 1, 13))), token.AND_AND),
							ast.NewLogicalExpressionNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
						),
					),
				},
			),
		},
		"nil coalescing operator is evaluated from left to right": {
			input: "foo ?? bar ?? baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							T(L(S(P(11, 1, 12), P(12, 1, 13))), token.QUESTION_QUESTION),
							ast.NewLogicalExpressionNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								T(L(S(P(4, 1, 5), P(5, 1, 6))), token.QUESTION_QUESTION),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
						),
					),
				},
			),
		},
		"or expression sequencing operator is evaluated from left to right": {
			input: "foo |! bar |! baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							T(L(S(P(11, 1, 12), P(12, 1, 13))), token.OR_BANG),
							ast.NewLogicalExpressionNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								T(L(S(P(4, 1, 5), P(5, 1, 6))), token.OR_BANG),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
						),
					),
				},
			),
		},
		"and expression sequencing operator is evaluated from left to right": {
			input: "foo &! bar &! baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							T(L(S(P(11, 1, 12), P(12, 1, 13))), token.AND_BANG),
							ast.NewLogicalExpressionNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_BANG),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(16, 1, 17))), "baz"),
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
