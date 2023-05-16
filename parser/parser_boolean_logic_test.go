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
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(4, 2, 1, 5), token.AND_AND),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(7, 10, 1, 8),
								T(P(11, 2, 1, 12), token.EQUAL_EQUAL),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
							),
						),
					),
				},
			),
		},
		"or has lower precedence than and": {
			input: "foo || bar && baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(4, 2, 1, 5), token.OR_OR),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewLogicalExpressionNode(
								P(7, 10, 1, 8),
								T(P(11, 2, 1, 12), token.AND_AND),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
							),
						),
					),
				},
			),
		},
		"nil coalescing operator has lower precedence than and": {
			input: "foo ?? bar && baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(4, 2, 1, 5), token.QUESTION_QUESTION),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewLogicalExpressionNode(
								P(7, 10, 1, 8),
								T(P(11, 2, 1, 12), token.AND_AND),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
							),
						),
					),
				},
			),
		},
		"or expression sequencing operator has lower precedence than and": {
			input: "foo |! bar && baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(4, 2, 1, 5), token.OR_BANG),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewLogicalExpressionNode(
								P(7, 10, 1, 8),
								T(P(11, 2, 1, 12), token.AND_AND),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
							),
						),
					),
				},
			),
		},
		"and expression sequencing operator has the same precedence as and": {
			input: "foo &! bar && baz &! boo",
			want: ast.NewProgramNode(
				P(0, 24, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 24, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 24, 1, 1),
							T(P(18, 2, 1, 19), token.AND_BANG),
							ast.NewLogicalExpressionNode(
								P(0, 17, 1, 1),
								T(P(11, 2, 1, 12), token.AND_AND),
								ast.NewLogicalExpressionNode(
									P(0, 10, 1, 1),
									T(P(4, 2, 1, 5), token.AND_BANG),
									ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
									ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
							),
							ast.NewPublicIdentifierNode(P(21, 3, 1, 22), "boo"),
						),
					),
				},
			),
		},
		"nil coalescing operator has the same precedence as or": {
			input: "foo ?? bar || baz ?? boo",
			want: ast.NewProgramNode(
				P(0, 24, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 24, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 24, 1, 1),
							T(P(18, 2, 1, 19), token.QUESTION_QUESTION),
							ast.NewLogicalExpressionNode(
								P(0, 17, 1, 1),
								T(P(11, 2, 1, 12), token.OR_OR),
								ast.NewLogicalExpressionNode(
									P(0, 10, 1, 1),
									T(P(4, 2, 1, 5), token.QUESTION_QUESTION),
									ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
									ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
							),
							ast.NewPublicIdentifierNode(P(21, 3, 1, 22), "boo"),
						),
					),
				},
			),
		},
		"or expression sequencing operator has the same precedence as or": {
			input: "foo |! bar || baz |! boo",
			want: ast.NewProgramNode(
				P(0, 24, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 24, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 24, 1, 1),
							T(P(18, 2, 1, 19), token.OR_BANG),
							ast.NewLogicalExpressionNode(
								P(0, 17, 1, 1),
								T(P(11, 2, 1, 12), token.OR_OR),
								ast.NewLogicalExpressionNode(
									P(0, 10, 1, 1),
									T(P(4, 2, 1, 5), token.OR_BANG),
									ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
									ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
							),
							ast.NewPublicIdentifierNode(P(21, 3, 1, 22), "boo"),
						),
					),
				},
			),
		},
		"or is evaluated from left to right": {
			input: "foo || bar || baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(11, 2, 1, 12), token.OR_OR),
							ast.NewLogicalExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), token.OR_OR),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
							),
							ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
						),
					),
				},
			),
		},
		"and is evaluated from left to right": {
			input: "foo && bar && baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(11, 2, 1, 12), token.AND_AND),
							ast.NewLogicalExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), token.AND_AND),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
							),
							ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
						),
					),
				},
			),
		},
		"nil coalescing operator is evaluated from left to right": {
			input: "foo ?? bar ?? baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(11, 2, 1, 12), token.QUESTION_QUESTION),
							ast.NewLogicalExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), token.QUESTION_QUESTION),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
							),
							ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
						),
					),
				},
			),
		},
		"or expression sequencing operator is evaluated from left to right": {
			input: "foo |! bar |! baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(11, 2, 1, 12), token.OR_BANG),
							ast.NewLogicalExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), token.OR_BANG),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
							),
							ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
						),
					),
				},
			),
		},
		"and expression sequencing operator is evaluated from left to right": {
			input: "foo &! bar &! baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(11, 2, 1, 12), token.AND_BANG),
							ast.NewLogicalExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), token.AND_BANG),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
							),
							ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
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
