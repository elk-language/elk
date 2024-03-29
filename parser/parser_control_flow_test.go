package parser

import (
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
)

func TestModifierExpression(t *testing.T) {
	tests := testTable{
		"has lower precedence than assignment": {
			input: "foo = bar if baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewModifierNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							T(S(P(10, 1, 11), P(11, 1, 12)), token.IF),
							ast.NewAssignmentExpressionNode(
								S(P(0, 1, 1), P(8, 1, 9)),
								T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
							),
							ast.NewPublicIdentifierNode(S(P(13, 1, 14), P(15, 1, 16)), "baz"),
						),
					),
				},
			),
		},
		"can have newlines after the modifier keyword": {
			input: "foo = bar if\nbaz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 2, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 2, 3)),
						ast.NewModifierNode(
							S(P(0, 1, 1), P(15, 2, 3)),
							T(S(P(10, 1, 11), P(11, 1, 12)), token.IF),
							ast.NewAssignmentExpressionNode(
								S(P(0, 1, 1), P(8, 1, 9)),
								T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
							),
							ast.NewPublicIdentifierNode(S(P(13, 2, 1), P(15, 2, 3)), "baz"),
						),
					),
				},
			),
		},
		"if can contain else": {
			input: "foo = bar if baz else car = red",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 1, 31)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(30, 1, 31)),
						ast.NewModifierIfElseNode(
							S(P(0, 1, 1), P(30, 1, 31)),
							ast.NewAssignmentExpressionNode(
								S(P(0, 1, 1), P(8, 1, 9)),
								T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
							),
							ast.NewPublicIdentifierNode(S(P(13, 1, 14), P(15, 1, 16)), "baz"),
							ast.NewAssignmentExpressionNode(
								S(P(22, 1, 23), P(30, 1, 31)),
								T(S(P(26, 1, 27), P(26, 1, 27)), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(S(P(22, 1, 23), P(24, 1, 25)), "car"),
								ast.NewPublicIdentifierNode(S(P(28, 1, 29), P(30, 1, 31)), "red"),
							),
						),
					),
				},
			),
		},
		"if else can span multiple lines": {
			input: "foo = bar if\nbaz else\ncar = red",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 3, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(30, 3, 9)),
						ast.NewModifierIfElseNode(
							S(P(0, 1, 1), P(30, 3, 9)),
							ast.NewAssignmentExpressionNode(
								S(P(0, 1, 1), P(8, 1, 9)),
								T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
							),
							ast.NewPublicIdentifierNode(S(P(13, 2, 1), P(15, 2, 3)), "baz"),
							ast.NewAssignmentExpressionNode(
								S(P(22, 3, 1), P(30, 3, 9)),
								T(S(P(26, 3, 5), P(26, 3, 5)), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(S(P(22, 3, 1), P(24, 3, 3)), "car"),
								ast.NewPublicIdentifierNode(S(P(28, 3, 7), P(30, 3, 9)), "red"),
							),
						),
					),
				},
			),
		},
		"can have for loops": {
			input: "println(i) for i in [1, 2, 3]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 1, 29)),
						ast.NewModifierForInNode(
							S(P(0, 1, 1), P(28, 1, 29)),
							ast.NewFunctionCallNode(
								S(P(0, 1, 1), P(9, 1, 10)),
								"println",
								[]ast.ExpressionNode{
									ast.NewPublicIdentifierNode(S(P(8, 1, 9), P(8, 1, 9)), "i"),
								},
								nil,
							),
							ast.NewPublicIdentifierNode(S(P(15, 1, 16), P(15, 1, 16)), "i"),
							ast.NewArrayListLiteralNode(
								S(P(20, 1, 21), P(28, 1, 29)),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(S(P(21, 1, 22), P(21, 1, 22)), "1"),
									ast.NewIntLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "2"),
									ast.NewIntLiteralNode(S(P(27, 1, 28), P(27, 1, 28)), "3"),
								},
								nil,
							),
						),
					),
				},
			),
		},
		"for loops can span multiple lines": {
			input: "println(i) for\ni\nin\n[1,\n2,\n3]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 6, 2)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 6, 2)),
						ast.NewModifierForInNode(
							S(P(0, 1, 1), P(28, 6, 2)),
							ast.NewFunctionCallNode(
								S(P(0, 1, 1), P(9, 1, 10)),
								"println",
								[]ast.ExpressionNode{
									ast.NewPublicIdentifierNode(S(P(8, 1, 9), P(8, 1, 9)), "i"),
								},
								nil,
							),
							ast.NewPublicIdentifierNode(S(P(15, 2, 1), P(15, 2, 1)), "i"),
							ast.NewArrayListLiteralNode(
								S(P(20, 4, 1), P(28, 6, 2)),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(S(P(21, 4, 2), P(21, 4, 2)), "1"),
									ast.NewIntLiteralNode(S(P(24, 5, 1), P(24, 5, 1)), "2"),
									ast.NewIntLiteralNode(S(P(27, 6, 1), P(27, 6, 1)), "3"),
								},
								nil,
							),
						),
					),
				},
			),
		},
		"has many versions": {
			input: "foo if bar\nfoo unless bar\nfoo while bar\nfoo until bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(52, 4, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewModifierNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							T(S(P(4, 1, 5), P(5, 1, 6)), token.IF),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(11, 2, 1), P(25, 2, 15)),
						ast.NewModifierNode(
							S(P(11, 2, 1), P(24, 2, 14)),
							T(S(P(15, 2, 5), P(20, 2, 10)), token.UNLESS),
							ast.NewPublicIdentifierNode(S(P(11, 2, 1), P(13, 2, 3)), "foo"),
							ast.NewPublicIdentifierNode(S(P(22, 2, 12), P(24, 2, 14)), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(26, 3, 1), P(39, 3, 14)),
						ast.NewModifierNode(
							S(P(26, 3, 1), P(38, 3, 13)),
							T(S(P(30, 3, 5), P(34, 3, 9)), token.WHILE),
							ast.NewPublicIdentifierNode(S(P(26, 3, 1), P(28, 3, 3)), "foo"),
							ast.NewPublicIdentifierNode(S(P(36, 3, 11), P(38, 3, 13)), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(40, 4, 1), P(52, 4, 13)),
						ast.NewModifierNode(
							S(P(40, 4, 1), P(52, 4, 13)),
							T(S(P(44, 4, 5), P(48, 4, 9)), token.UNTIL),
							ast.NewPublicIdentifierNode(S(P(40, 4, 1), P(42, 4, 3)), "foo"),
							ast.NewPublicIdentifierNode(S(P(50, 4, 11), P(52, 4, 13)), "bar"),
						),
					),
				},
			),
		},
		"cannot be nested": {
			input: "foo = bar if baz if false\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 2, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewModifierNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							T(S(P(10, 1, 11), P(11, 1, 12)), token.IF),
							ast.NewAssignmentExpressionNode(
								S(P(0, 1, 1), P(8, 1, 9)),
								T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
								ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
							),
							ast.NewPublicIdentifierNode(S(P(13, 1, 14), P(15, 1, 16)), "baz"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(26, 2, 1), P(26, 2, 1)),
						ast.NewIntLiteralNode(S(P(26, 2, 1), P(26, 2, 1)), "3"),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(17, 1, 18), P(18, 1, 19)), "unexpected if, expected a statement separator `\\n`, `;`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestIf(t *testing.T) {
	tests := testTable{
		"can have one branch": {
			input: `
if foo > 0
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(30, 5, 4)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(29, 5, 3)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 10)),
								T(S(P(8, 2, 8), P(8, 2, 8)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 6)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 10)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(13, 3, 2), P(21, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(13, 3, 2), P(20, 3, 9)),
										T(S(P(17, 3, 6), P(18, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(13, 3, 2), P(15, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(20, 3, 9), P(20, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(23, 4, 2), P(26, 4, 5)),
									ast.NewNilLiteralNode(S(P(23, 4, 2), P(25, 4, 4))),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
if foo > 0
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(15, 3, 4)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(14, 3, 3)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 10)),
								T(S(P(8, 2, 8), P(8, 2, 8)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 6)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 10)), "0"),
							),
							nil,
							nil,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
bar =
	if foo > 0
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 6, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(34, 5, 5)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(33, 5, 4)),
							T(S(P(5, 2, 5), P(5, 2, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 3)), "bar"),
							ast.NewIfExpressionNode(
								S(P(8, 3, 2), P(33, 5, 4)),
								ast.NewBinaryExpressionNode(
									S(P(11, 3, 5), P(17, 3, 11)),
									T(S(P(15, 3, 9), P(15, 3, 9)), token.GREATER),
									ast.NewPublicIdentifierNode(S(P(11, 3, 5), P(13, 3, 7)), "foo"),
									ast.NewIntLiteralNode(S(P(17, 3, 11), P(17, 3, 11)), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(21, 4, 3), P(29, 4, 11)),
										ast.NewAssignmentExpressionNode(
											S(P(21, 4, 3), P(28, 4, 10)),
											T(S(P(25, 4, 7), P(26, 4, 8)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(21, 4, 3), P(23, 4, 5)), "foo"),
											ast.NewIntLiteralNode(S(P(28, 4, 10), P(28, 4, 10)), "2"),
										),
									),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(35, 6, 1), P(38, 6, 4)),
						ast.NewNilLiteralNode(S(P(35, 6, 1), P(37, 6, 3))),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
if foo > 0 then foo += 2
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(25, 2, 25)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(24, 2, 24)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 10)),
								T(S(P(8, 2, 8), P(8, 2, 8)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 6)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 10)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 2, 17), P(24, 2, 24)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 2, 17), P(24, 2, 24)),
										T(S(P(21, 2, 21), P(22, 2, 22)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 2, 17), P(19, 2, 19)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 2, 24), P(24, 2, 24)), "2"),
									),
								),
							},
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						S(P(26, 3, 1), P(29, 3, 4)),
						ast.NewNilLiteralNode(S(P(26, 3, 1), P(28, 3, 3))),
					),
				},
			),
		},
		"can have else": {
			input: `
if foo > 0
	foo += 2
	nil
else
  foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(55, 9, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(51, 8, 4)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(50, 8, 3)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 10)),
								T(S(P(8, 2, 8), P(8, 2, 8)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 6)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 10)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(13, 3, 2), P(21, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(13, 3, 2), P(20, 3, 9)),
										T(S(P(17, 3, 6), P(18, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(13, 3, 2), P(15, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(20, 3, 9), P(20, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(23, 4, 2), P(26, 4, 5)),
									ast.NewNilLiteralNode(S(P(23, 4, 2), P(25, 4, 4))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(34, 6, 3), P(42, 6, 11)),
									ast.NewAssignmentExpressionNode(
										S(P(34, 6, 3), P(41, 6, 10)),
										T(S(P(38, 6, 7), P(39, 6, 8)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(34, 6, 3), P(36, 6, 5)), "foo"),
										ast.NewIntLiteralNode(S(P(41, 6, 10), P(41, 6, 10)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(44, 7, 2), P(47, 7, 5)),
									ast.NewNilLiteralNode(S(P(44, 7, 2), P(46, 7, 4))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(52, 9, 1), P(55, 9, 4)),
						ast.NewNilLiteralNode(S(P(52, 9, 1), P(54, 9, 3))),
					),
				},
			),
		},
		"can have else in short form": {
			input: `
if foo > 0 then foo += 2
else foo -= 2
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(43, 4, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(39, 3, 14)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(38, 3, 13)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 10)),
								T(S(P(8, 2, 8), P(8, 2, 8)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 6)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 10)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 2, 17), P(24, 2, 24)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 2, 17), P(24, 2, 24)),
										T(S(P(21, 2, 21), P(22, 2, 22)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 2, 17), P(19, 2, 19)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 2, 24), P(24, 2, 24)), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(31, 3, 6), P(38, 3, 13)),
									ast.NewAssignmentExpressionNode(
										S(P(31, 3, 6), P(38, 3, 13)),
										T(S(P(35, 3, 10), P(36, 3, 11)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(31, 3, 6), P(33, 3, 8)), "foo"),
										ast.NewIntLiteralNode(S(P(38, 3, 13), P(38, 3, 13)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(40, 4, 1), P(43, 4, 4)),
						ast.NewNilLiteralNode(S(P(40, 4, 1), P(42, 4, 3))),
					),
				},
			),
		},
		"cannot have two elses": {
			input: `
if foo > 0 then foo += 2
else foo -= 2
else bar
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(52, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(39, 3, 14)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(38, 3, 13)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 10)),
								T(S(P(8, 2, 8), P(8, 2, 8)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 6)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 10)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 2, 17), P(24, 2, 24)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 2, 17), P(24, 2, 24)),
										T(S(P(21, 2, 21), P(22, 2, 22)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 2, 17), P(19, 2, 19)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 2, 24), P(24, 2, 24)), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(31, 3, 6), P(38, 3, 13)),
									ast.NewAssignmentExpressionNode(
										S(P(31, 3, 6), P(38, 3, 13)),
										T(S(P(35, 3, 10), P(36, 3, 11)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(31, 3, 6), P(33, 3, 8)), "foo"),
										ast.NewIntLiteralNode(S(P(38, 3, 13), P(38, 3, 13)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(40, 4, 1), P(48, 4, 9)),
						ast.NewInvalidNode(S(P(40, 4, 1), P(43, 4, 4)), T(S(P(40, 4, 1), P(43, 4, 4)), token.ELSE)),
					),
					ast.NewExpressionStatementNode(
						S(P(49, 5, 1), P(52, 5, 4)),
						ast.NewNilLiteralNode(S(P(49, 5, 1), P(51, 5, 3))),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(40, 4, 1), P(43, 4, 4)), "unexpected else, expected an expression"),
			},
		},
		"can have many elsif blocks": {
			input: `
if foo > 0
	foo += 2
	nil
elsif foo < 5
	foo *= 10
elsif foo < 0
	foo %= 3
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(103, 13, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(99, 12, 4)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(98, 12, 3)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 10)),
								T(S(P(8, 2, 8), P(8, 2, 8)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 6)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 10)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(13, 3, 2), P(21, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(13, 3, 2), P(20, 3, 9)),
										T(S(P(17, 3, 6), P(18, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(13, 3, 2), P(15, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(20, 3, 9), P(20, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(23, 4, 2), P(26, 4, 5)),
									ast.NewNilLiteralNode(S(P(23, 4, 2), P(25, 4, 4))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(27, 5, 1), P(51, 6, 11)),
									ast.NewIfExpressionNode(
										S(P(27, 5, 1), P(51, 6, 11)),
										ast.NewBinaryExpressionNode(
											S(P(33, 5, 7), P(39, 5, 13)),
											T(S(P(37, 5, 11), P(37, 5, 11)), token.LESS),
											ast.NewPublicIdentifierNode(S(P(33, 5, 7), P(35, 5, 9)), "foo"),
											ast.NewIntLiteralNode(S(P(39, 5, 13), P(39, 5, 13)), "5"),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(42, 6, 2), P(51, 6, 11)),
												ast.NewAssignmentExpressionNode(
													S(P(42, 6, 2), P(50, 6, 10)),
													T(S(P(46, 6, 6), P(47, 6, 7)), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(S(P(42, 6, 2), P(44, 6, 4)), "foo"),
													ast.NewIntLiteralNode(S(P(49, 6, 9), P(50, 6, 10)), "10"),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(52, 7, 1), P(75, 8, 10)),
												ast.NewIfExpressionNode(
													S(P(52, 7, 1), P(98, 12, 3)),
													ast.NewBinaryExpressionNode(
														S(P(58, 7, 7), P(64, 7, 13)),
														T(S(P(62, 7, 11), P(62, 7, 11)), token.LESS),
														ast.NewPublicIdentifierNode(S(P(58, 7, 7), P(60, 7, 9)), "foo"),
														ast.NewIntLiteralNode(S(P(64, 7, 13), P(64, 7, 13)), "0"),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(67, 8, 2), P(75, 8, 10)),
															ast.NewAssignmentExpressionNode(
																S(P(67, 8, 2), P(74, 8, 9)),
																T(S(P(71, 8, 6), P(72, 8, 7)), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(S(P(67, 8, 2), P(69, 8, 4)), "foo"),
																ast.NewIntLiteralNode(S(P(74, 8, 9), P(74, 8, 9)), "3"),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(82, 10, 2), P(90, 10, 10)),
															ast.NewAssignmentExpressionNode(
																S(P(82, 10, 2), P(89, 10, 9)),
																T(S(P(86, 10, 6), P(87, 10, 7)), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(S(P(82, 10, 2), P(84, 10, 4)), "foo"),
																ast.NewIntLiteralNode(S(P(89, 10, 9), P(89, 10, 9)), "2"),
															),
														),
														ast.NewExpressionStatementNode(
															S(P(92, 11, 2), P(95, 11, 5)),
															ast.NewNilLiteralNode(S(P(92, 11, 2), P(94, 11, 4))),
														),
													},
												),
											),
										},
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(100, 13, 1), P(103, 13, 4)),
						ast.NewNilLiteralNode(S(P(100, 13, 1), P(102, 13, 3))),
					),
				},
			),
		},
		"can have elsifs in short form": {
			input: `
if foo > 0 then foo += 2
elsif foo < 5 then foo *= 10
elsif foo < 0 then foo %= 3
else foo -= 2
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(100, 6, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(96, 5, 14)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(95, 5, 13)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 10)),
								T(S(P(8, 2, 8), P(8, 2, 8)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 6)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 10)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 2, 17), P(24, 2, 24)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 2, 17), P(24, 2, 24)),
										T(S(P(21, 2, 21), P(22, 2, 22)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 2, 17), P(19, 2, 19)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 2, 24), P(24, 2, 24)), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(26, 3, 1), P(53, 3, 28)),
									ast.NewIfExpressionNode(
										S(P(26, 3, 1), P(53, 3, 28)),
										ast.NewBinaryExpressionNode(
											S(P(32, 3, 7), P(38, 3, 13)),
											T(S(P(36, 3, 11), P(36, 3, 11)), token.LESS),
											ast.NewPublicIdentifierNode(S(P(32, 3, 7), P(34, 3, 9)), "foo"),
											ast.NewIntLiteralNode(S(P(38, 3, 13), P(38, 3, 13)), "5"),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(45, 3, 20), P(53, 3, 28)),
												ast.NewAssignmentExpressionNode(
													S(P(45, 3, 20), P(53, 3, 28)),
													T(S(P(49, 3, 24), P(50, 3, 25)), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(S(P(45, 3, 20), P(47, 3, 22)), "foo"),
													ast.NewIntLiteralNode(S(P(52, 3, 27), P(53, 3, 28)), "10"),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(55, 4, 1), P(81, 4, 27)),
												ast.NewIfExpressionNode(
													S(P(55, 4, 1), P(95, 5, 13)),
													ast.NewBinaryExpressionNode(
														S(P(61, 4, 7), P(67, 4, 13)),
														T(S(P(65, 4, 11), P(65, 4, 11)), token.LESS),
														ast.NewPublicIdentifierNode(S(P(61, 4, 7), P(63, 4, 9)), "foo"),
														ast.NewIntLiteralNode(S(P(67, 4, 13), P(67, 4, 13)), "0"),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(74, 4, 20), P(81, 4, 27)),
															ast.NewAssignmentExpressionNode(
																S(P(74, 4, 20), P(81, 4, 27)),
																T(S(P(78, 4, 24), P(79, 4, 25)), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(S(P(74, 4, 20), P(76, 4, 22)), "foo"),
																ast.NewIntLiteralNode(S(P(81, 4, 27), P(81, 4, 27)), "3"),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(88, 5, 6), P(95, 5, 13)),
															ast.NewAssignmentExpressionNode(
																S(P(88, 5, 6), P(95, 5, 13)),
																T(S(P(92, 5, 10), P(93, 5, 11)), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(S(P(88, 5, 6), P(90, 5, 8)), "foo"),
																ast.NewIntLiteralNode(S(P(95, 5, 13), P(95, 5, 13)), "2"),
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
					),
					ast.NewExpressionStatementNode(
						S(P(97, 6, 1), P(100, 6, 4)),
						ast.NewNilLiteralNode(S(P(97, 6, 1), P(99, 6, 3))),
					),
				},
			),
		},
		"else if is also possible": {
			input: `
if foo > 0
	foo += 2
	nil
else if foo < 5
	foo *= 10
else if foo < 0
	foo %= 3
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(107, 13, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(103, 12, 4)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(102, 12, 3)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 10)),
								T(S(P(8, 2, 8), P(8, 2, 8)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 6)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 10)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(13, 3, 2), P(21, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(13, 3, 2), P(20, 3, 9)),
										T(S(P(17, 3, 6), P(18, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(13, 3, 2), P(15, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(20, 3, 9), P(20, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(23, 4, 2), P(26, 4, 5)),
									ast.NewNilLiteralNode(S(P(23, 4, 2), P(25, 4, 4))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(32, 5, 6), P(102, 12, 3)),
									ast.NewIfExpressionNode(
										S(P(32, 5, 6), P(102, 12, 3)),
										ast.NewBinaryExpressionNode(
											S(P(35, 5, 9), P(41, 5, 15)),
											T(S(P(39, 5, 13), P(39, 5, 13)), token.LESS),
											ast.NewPublicIdentifierNode(S(P(35, 5, 9), P(37, 5, 11)), "foo"),
											ast.NewIntLiteralNode(S(P(41, 5, 15), P(41, 5, 15)), "5"),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(44, 6, 2), P(53, 6, 11)),
												ast.NewAssignmentExpressionNode(
													S(P(44, 6, 2), P(52, 6, 10)),
													T(S(P(48, 6, 6), P(49, 6, 7)), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(S(P(44, 6, 2), P(46, 6, 4)), "foo"),
													ast.NewIntLiteralNode(S(P(51, 6, 9), P(52, 6, 10)), "10"),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(59, 7, 6), P(102, 12, 3)),
												ast.NewIfExpressionNode(
													S(P(59, 7, 6), P(102, 12, 3)),
													ast.NewBinaryExpressionNode(
														S(P(62, 7, 9), P(68, 7, 15)),
														T(S(P(66, 7, 13), P(66, 7, 13)), token.LESS),
														ast.NewPublicIdentifierNode(S(P(62, 7, 9), P(64, 7, 11)), "foo"),
														ast.NewIntLiteralNode(S(P(68, 7, 15), P(68, 7, 15)), "0"),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(71, 8, 2), P(79, 8, 10)),
															ast.NewAssignmentExpressionNode(
																S(P(71, 8, 2), P(78, 8, 9)),
																T(S(P(75, 8, 6), P(76, 8, 7)), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(S(P(71, 8, 2), P(73, 8, 4)), "foo"),
																ast.NewIntLiteralNode(S(P(78, 8, 9), P(78, 8, 9)), "3"),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(86, 10, 2), P(94, 10, 10)),
															ast.NewAssignmentExpressionNode(
																S(P(86, 10, 2), P(93, 10, 9)),
																T(S(P(90, 10, 6), P(91, 10, 7)), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(S(P(86, 10, 2), P(88, 10, 4)), "foo"),
																ast.NewIntLiteralNode(S(P(93, 10, 9), P(93, 10, 9)), "2"),
															),
														),
														ast.NewExpressionStatementNode(
															S(P(96, 11, 2), P(99, 11, 5)),
															ast.NewNilLiteralNode(S(P(96, 11, 2), P(98, 11, 4))),
														),
													},
												),
											),
										},
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(104, 13, 1), P(107, 13, 4)),
						ast.NewNilLiteralNode(S(P(104, 13, 1), P(106, 13, 3))),
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

func TestUnless(t *testing.T) {
	tests := testTable{
		"can have one branch": {
			input: `
unless foo > 0
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(34, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(34, 5, 4)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(33, 5, 3)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 14)),
								T(S(P(12, 2, 12), P(12, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 14)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 3, 2), P(25, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 3, 2), P(24, 3, 9)),
										T(S(P(21, 3, 6), P(22, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 3, 2), P(19, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 3, 9), P(24, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(27, 4, 2), P(30, 4, 5)),
									ast.NewNilLiteralNode(S(P(27, 4, 2), P(29, 4, 4))),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
unless foo > 0
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(19, 3, 4)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(18, 3, 3)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 14)),
								T(S(P(12, 2, 12), P(12, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 14)), "0"),
							),
							nil,
							nil,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
bar =
	unless foo > 0
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(42, 6, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(38, 5, 5)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(37, 5, 4)),
							T(S(P(5, 2, 5), P(5, 2, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 3)), "bar"),
							ast.NewUnlessExpressionNode(
								S(P(8, 3, 2), P(37, 5, 4)),
								ast.NewBinaryExpressionNode(
									S(P(15, 3, 9), P(21, 3, 15)),
									T(S(P(19, 3, 13), P(19, 3, 13)), token.GREATER),
									ast.NewPublicIdentifierNode(S(P(15, 3, 9), P(17, 3, 11)), "foo"),
									ast.NewIntLiteralNode(S(P(21, 3, 15), P(21, 3, 15)), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(25, 4, 3), P(33, 4, 11)),
										ast.NewAssignmentExpressionNode(
											S(P(25, 4, 3), P(32, 4, 10)),
											T(S(P(29, 4, 7), P(30, 4, 8)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(25, 4, 3), P(27, 4, 5)), "foo"),
											ast.NewIntLiteralNode(S(P(32, 4, 10), P(32, 4, 10)), "2"),
										),
									),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(39, 6, 1), P(42, 6, 4)),
						ast.NewNilLiteralNode(S(P(39, 6, 1), P(41, 6, 3))),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
unless foo > 0 then foo += 2
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(29, 2, 29)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(28, 2, 28)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 14)),
								T(S(P(12, 2, 12), P(12, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 14)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 2, 21), P(28, 2, 28)),
									ast.NewAssignmentExpressionNode(
										S(P(21, 2, 21), P(28, 2, 28)),
										T(S(P(25, 2, 25), P(26, 2, 26)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(21, 2, 21), P(23, 2, 23)), "foo"),
										ast.NewIntLiteralNode(S(P(28, 2, 28), P(28, 2, 28)), "2"),
									),
								),
							},
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						S(P(30, 3, 1), P(33, 3, 4)),
						ast.NewNilLiteralNode(S(P(30, 3, 1), P(32, 3, 3))),
					),
				},
			),
		},
		"can have else": {
			input: `
unless foo > 0
	foo += 2
	nil
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(58, 9, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(54, 8, 4)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(53, 8, 3)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 14)),
								T(S(P(12, 2, 12), P(12, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 14)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 3, 2), P(25, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 3, 2), P(24, 3, 9)),
										T(S(P(21, 3, 6), P(22, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 3, 2), P(19, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 3, 9), P(24, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(27, 4, 2), P(30, 4, 5)),
									ast.NewNilLiteralNode(S(P(27, 4, 2), P(29, 4, 4))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(37, 6, 2), P(45, 6, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(37, 6, 2), P(44, 6, 9)),
										T(S(P(41, 6, 6), P(42, 6, 7)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(37, 6, 2), P(39, 6, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(44, 6, 9), P(44, 6, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(47, 7, 2), P(50, 7, 5)),
									ast.NewNilLiteralNode(S(P(47, 7, 2), P(49, 7, 4))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(55, 9, 1), P(58, 9, 4)),
						ast.NewNilLiteralNode(S(P(55, 9, 1), P(57, 9, 3))),
					),
				},
			),
		},
		"can have else in short form": {
			input: `
unless foo > 0 then foo += 2
else foo -= 2
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(47, 4, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(43, 3, 14)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(42, 3, 13)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 14)),
								T(S(P(12, 2, 12), P(12, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 14)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 2, 21), P(28, 2, 28)),
									ast.NewAssignmentExpressionNode(
										S(P(21, 2, 21), P(28, 2, 28)),
										T(S(P(25, 2, 25), P(26, 2, 26)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(21, 2, 21), P(23, 2, 23)), "foo"),
										ast.NewIntLiteralNode(S(P(28, 2, 28), P(28, 2, 28)), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(35, 3, 6), P(42, 3, 13)),
									ast.NewAssignmentExpressionNode(
										S(P(35, 3, 6), P(42, 3, 13)),
										T(S(P(39, 3, 10), P(40, 3, 11)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(35, 3, 6), P(37, 3, 8)), "foo"),
										ast.NewIntLiteralNode(S(P(42, 3, 13), P(42, 3, 13)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(44, 4, 1), P(47, 4, 4)),
						ast.NewNilLiteralNode(S(P(44, 4, 1), P(46, 4, 3))),
					),
				},
			),
		},
		"cannot have two elses": {
			input: `
unless foo > 0 then foo += 2
else foo -= 2
else bar
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(56, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(43, 3, 14)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(42, 3, 13)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 14)),
								T(S(P(12, 2, 12), P(12, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 14)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 2, 21), P(28, 2, 28)),
									ast.NewAssignmentExpressionNode(
										S(P(21, 2, 21), P(28, 2, 28)),
										T(S(P(25, 2, 25), P(26, 2, 26)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(21, 2, 21), P(23, 2, 23)), "foo"),
										ast.NewIntLiteralNode(S(P(28, 2, 28), P(28, 2, 28)), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(35, 3, 6), P(42, 3, 13)),
									ast.NewAssignmentExpressionNode(
										S(P(35, 3, 6), P(42, 3, 13)),
										T(S(P(39, 3, 10), P(40, 3, 11)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(35, 3, 6), P(37, 3, 8)), "foo"),
										ast.NewIntLiteralNode(S(P(42, 3, 13), P(42, 3, 13)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(44, 4, 1), P(52, 4, 9)),
						ast.NewInvalidNode(S(P(44, 4, 1), P(47, 4, 4)), T(S(P(44, 4, 1), P(47, 4, 4)), token.ELSE)),
					),
					ast.NewExpressionStatementNode(
						S(P(53, 5, 1), P(56, 5, 4)),
						ast.NewNilLiteralNode(S(P(53, 5, 1), P(55, 5, 3))),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(44, 4, 1), P(47, 4, 4)), "unexpected else, expected an expression"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestWhile(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
while foo > 0
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(33, 5, 4)),
						ast.NewWhileExpressionNode(
							S(P(1, 2, 1), P(32, 5, 3)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 13)),
								T(S(P(11, 2, 11), P(11, 2, 11)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 9)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 13)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(16, 3, 2), P(24, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(16, 3, 2), P(23, 3, 9)),
										T(S(P(20, 3, 6), P(21, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(16, 3, 2), P(18, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(23, 3, 9), P(23, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(26, 4, 2), P(29, 4, 5)),
									ast.NewNilLiteralNode(S(P(26, 4, 2), P(28, 4, 4))),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
while foo > 0
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(18, 3, 4)),
						ast.NewWhileExpressionNode(
							S(P(1, 2, 1), P(17, 3, 3)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 13)),
								T(S(P(11, 2, 11), P(11, 2, 11)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 9)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 13)), "0"),
							),
							nil,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
bar =
	while foo > 0
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(41, 6, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(37, 5, 5)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(36, 5, 4)),
							T(S(P(5, 2, 5), P(5, 2, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 3)), "bar"),
							ast.NewWhileExpressionNode(
								S(P(8, 3, 2), P(36, 5, 4)),
								ast.NewBinaryExpressionNode(
									S(P(14, 3, 8), P(20, 3, 14)),
									T(S(P(18, 3, 12), P(18, 3, 12)), token.GREATER),
									ast.NewPublicIdentifierNode(S(P(14, 3, 8), P(16, 3, 10)), "foo"),
									ast.NewIntLiteralNode(S(P(20, 3, 14), P(20, 3, 14)), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(24, 4, 3), P(32, 4, 11)),
										ast.NewAssignmentExpressionNode(
											S(P(24, 4, 3), P(31, 4, 10)),
											T(S(P(28, 4, 7), P(29, 4, 8)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(24, 4, 3), P(26, 4, 5)), "foo"),
											ast.NewIntLiteralNode(S(P(31, 4, 10), P(31, 4, 10)), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(38, 6, 1), P(41, 6, 4)),
						ast.NewNilLiteralNode(S(P(38, 6, 1), P(40, 6, 3))),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
while foo > 0 then foo += 2
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(28, 2, 28)),
						ast.NewWhileExpressionNode(
							S(P(1, 2, 1), P(27, 2, 27)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 13)),
								T(S(P(11, 2, 11), P(11, 2, 11)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 9)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 13)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(20, 2, 20), P(27, 2, 27)),
									ast.NewAssignmentExpressionNode(
										S(P(20, 2, 20), P(27, 2, 27)),
										T(S(P(24, 2, 24), P(25, 2, 25)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(20, 2, 20), P(22, 2, 22)), "foo"),
										ast.NewIntLiteralNode(S(P(27, 2, 27), P(27, 2, 27)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(29, 3, 1), P(32, 3, 4)),
						ast.NewNilLiteralNode(S(P(29, 3, 1), P(31, 3, 3))),
					),
				},
			),
		},
		"cannot have else": {
			input: `
while foo > 0
	foo += 2
	nil
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(57, 9, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(53, 8, 4)),
						ast.NewWhileExpressionNode(
							S(P(1, 2, 1), P(52, 8, 3)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 13)),
								T(S(P(11, 2, 11), P(11, 2, 11)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 9)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 13)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(16, 3, 2), P(24, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(16, 3, 2), P(23, 3, 9)),
										T(S(P(20, 3, 6), P(21, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(16, 3, 2), P(18, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(23, 3, 9), P(23, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(26, 4, 2), P(29, 4, 5)),
									ast.NewNilLiteralNode(S(P(26, 4, 2), P(28, 4, 4))),
								),
								ast.NewExpressionStatementNode(
									S(P(30, 5, 1), P(34, 5, 5)),
									ast.NewInvalidNode(S(P(30, 5, 1), P(33, 5, 4)), T(S(P(30, 5, 1), P(33, 5, 4)), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									S(P(36, 6, 2), P(44, 6, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(36, 6, 2), P(43, 6, 9)),
										T(S(P(40, 6, 6), P(41, 6, 7)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(36, 6, 2), P(38, 6, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(43, 6, 9), P(43, 6, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(46, 7, 2), P(49, 7, 5)),
									ast.NewNilLiteralNode(S(P(46, 7, 2), P(48, 7, 4))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(54, 9, 1), P(57, 9, 4)),
						ast.NewNilLiteralNode(S(P(54, 9, 1), P(56, 9, 3))),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(30, 5, 1), P(33, 5, 4)), "unexpected else, expected an expression"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestSwitch(t *testing.T) {
	tests := testTable{
		"cannot be empty": {
			input: `
switch foo
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(15, 3, 4)),
						ast.NewSwitchExpressionNode(
							S(P(1, 2, 1), P(14, 3, 3)),
							ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 10)), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(1, 2, 1), P(14, 3, 3)), "switch cannot be empty"),
			},
		},
		"is an expression": {
			input: `
bar =
	switch foo
	case n
		n + 2
	end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(43, 7, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(39, 6, 5)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(38, 6, 4)),
							T(S(P(5, 2, 5), P(5, 2, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 3)), "bar"),
							ast.NewSwitchExpressionNode(
								S(P(8, 3, 2), P(38, 6, 4)),
								ast.NewPublicIdentifierNode(S(P(15, 3, 9), P(17, 3, 11)), "foo"),
								[]*ast.CaseNode{
									ast.NewCaseNode(
										S(P(20, 4, 2), P(34, 5, 8)),
										ast.NewPublicIdentifierNode(S(P(25, 4, 7), P(25, 4, 7)), "n"),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(29, 5, 3), P(34, 5, 8)),
												ast.NewBinaryExpressionNode(
													S(P(29, 5, 3), P(33, 5, 7)),
													T(S(P(31, 5, 5), P(31, 5, 5)), token.PLUS),
													ast.NewPublicIdentifierNode(S(P(29, 5, 3), P(29, 5, 3)), "n"),
													ast.NewIntLiteralNode(S(P(33, 5, 7), P(33, 5, 7)), "2"),
												),
											),
										},
									),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(40, 7, 1), P(43, 7, 4)),
						ast.NewNilLiteralNode(S(P(40, 7, 1), P(42, 7, 3))),
					),
				},
			),
		},
		"cannot have only have else": {
			input: `
switch foo
else
  n + 2
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(28, 5, 4)),
						ast.NewSwitchExpressionNode(
							S(P(1, 2, 1), P(27, 5, 3)),
							ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 10)), "foo"),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(19, 4, 3), P(24, 4, 8)),
									ast.NewBinaryExpressionNode(
										S(P(19, 4, 3), P(23, 4, 7)),
										T(S(P(21, 4, 5), P(21, 4, 5)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(19, 4, 3), P(19, 4, 3)), "n"),
										ast.NewIntLiteralNode(S(P(23, 4, 7), P(23, 4, 7)), "2"),
									),
								),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(1, 2, 1), P(27, 5, 3)), "switch cannot only consist of else"),
			},
		},
		"can have multiple branches": {
			input: `
switch foo
case n
  n
case m
  m
else
  n + 2
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(50, 9, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(50, 9, 4)),
						ast.NewSwitchExpressionNode(
							S(P(1, 2, 1), P(49, 9, 3)),
							ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(12, 3, 1), P(22, 4, 4)),
									ast.NewPublicIdentifierNode(S(P(17, 3, 6), P(17, 3, 6)), "n"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(21, 4, 3), P(22, 4, 4)),
											ast.NewPublicIdentifierNode(S(P(21, 4, 3), P(21, 4, 3)), "n"),
										),
									},
								),
								ast.NewCaseNode(
									S(P(23, 5, 1), P(33, 6, 4)),
									ast.NewPublicIdentifierNode(S(P(28, 5, 6), P(28, 5, 6)), "m"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(32, 6, 3), P(33, 6, 4)),
											ast.NewPublicIdentifierNode(S(P(32, 6, 3), P(32, 6, 3)), "m"),
										),
									},
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(41, 8, 3), P(46, 8, 8)),
									ast.NewBinaryExpressionNode(
										S(P(41, 8, 3), P(45, 8, 7)),
										T(S(P(43, 8, 5), P(43, 8, 5)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(41, 8, 3), P(41, 8, 3)), "n"),
										ast.NewIntLiteralNode(S(P(45, 8, 7), P(45, 8, 7)), "2"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have short branches with then": {
			input: `
switch foo
case n then n
case m then m
else n + 2
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(54, 6, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(54, 6, 4)),
						ast.NewSwitchExpressionNode(
							S(P(1, 2, 1), P(53, 6, 3)),
							ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(12, 3, 1), P(24, 3, 13)),
									ast.NewPublicIdentifierNode(S(P(17, 3, 6), P(17, 3, 6)), "n"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(24, 3, 13), P(24, 3, 13)),
											ast.NewPublicIdentifierNode(S(P(24, 3, 13), P(24, 3, 13)), "n"),
										),
									},
								),
								ast.NewCaseNode(
									S(P(26, 4, 1), P(38, 4, 13)),
									ast.NewPublicIdentifierNode(S(P(31, 4, 6), P(31, 4, 6)), "m"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(38, 4, 13), P(38, 4, 13)),
											ast.NewPublicIdentifierNode(S(P(38, 4, 13), P(38, 4, 13)), "m"),
										),
									},
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(45, 5, 6), P(49, 5, 10)),
									ast.NewBinaryExpressionNode(
										S(P(45, 5, 6), P(49, 5, 10)),
										T(S(P(47, 5, 8), P(47, 5, 8)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(45, 5, 6), P(45, 5, 6)), "n"),
										ast.NewIntLiteralNode(S(P(49, 5, 10), P(49, 5, 10)), "2"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can be single-line with then": {
			input: `switch foo case n then n case m then m else n + 2 end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(52, 1, 53)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(52, 1, 53)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(52, 1, 53)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(23, 1, 24)),
									ast.NewPublicIdentifierNode(S(P(16, 1, 17), P(16, 1, 17)), "n"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(23, 1, 24), P(23, 1, 24)),
											ast.NewPublicIdentifierNode(S(P(23, 1, 24), P(23, 1, 24)), "n"),
										),
									},
								),
								ast.NewCaseNode(
									S(P(25, 1, 26), P(37, 1, 38)),
									ast.NewPublicIdentifierNode(S(P(30, 1, 31), P(30, 1, 31)), "m"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(37, 1, 38), P(37, 1, 38)),
											ast.NewPublicIdentifierNode(S(P(37, 1, 38), P(37, 1, 38)), "m"),
										),
									},
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(44, 1, 45), P(48, 1, 49)),
									ast.NewBinaryExpressionNode(
										S(P(44, 1, 45), P(48, 1, 49)),
										T(S(P(46, 1, 47), P(46, 1, 47)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(44, 1, 45), P(44, 1, 45)), "n"),
										ast.NewIntLiteralNode(S(P(48, 1, 49), P(48, 1, 49)), "2"),
									),
								),
							},
						),
					),
				},
			),
		},
		"pattern can be true": {
			input: `switch foo case true then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewTrueLiteralNode(S(P(16, 1, 17), P(19, 1, 20))),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be false": {
			input: `switch foo case false then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(29, 1, 30)),
									ast.NewFalseLiteralNode(S(P(16, 1, 17), P(20, 1, 21))),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(27, 1, 28), P(29, 1, 30)),
											ast.NewNilLiteralNode(S(P(27, 1, 28), P(29, 1, 30))),
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
		"pattern can be nil": {
			input: `switch foo case nil then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 1, 32)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 1, 32)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(31, 1, 32)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(27, 1, 28)),
									ast.NewNilLiteralNode(S(P(16, 1, 17), P(18, 1, 19))),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(25, 1, 26), P(27, 1, 28)),
											ast.NewNilLiteralNode(S(P(25, 1, 26), P(27, 1, 28))),
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
		"pattern can be a char": {
			input: "switch foo case `f` then nil end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 1, 32)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 1, 32)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(31, 1, 32)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(27, 1, 28)),
									ast.NewCharLiteralNode(S(P(16, 1, 17), P(18, 1, 19)), 'f'),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(25, 1, 26), P(27, 1, 28)),
											ast.NewNilLiteralNode(S(P(25, 1, 26), P(27, 1, 28))),
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
		"pattern can be a raw char": {
			input: "switch foo case r`f` then nil end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewRawCharLiteralNode(S(P(16, 1, 17), P(19, 1, 20)), 'f'),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a raw string": {
			input: "switch foo case 'fo' then nil end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewRawStringLiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "fo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a string": {
			input: `switch foo case "fo" then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewDoubleQuotedStringLiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "fo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a regex": {
			input: `switch foo case %/f/ then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewUninterpolatedRegexLiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "f", bitfield.BitField8{}),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a simple symbol": {
			input: `switch foo case :foo then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewSimpleSymbolLiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "foo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a symbol with quotes": {
			input: `switch foo case :'&' then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewSimpleSymbolLiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "&"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a public identifier": {
			input: `switch foo case foof then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewPublicIdentifierNode(S(P(16, 1, 17), P(19, 1, 20)), "foof"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a private identifier": {
			input: `switch foo case _foo then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewPrivateIdentifierNode(S(P(16, 1, 17), P(19, 1, 20)), "_foo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be an int": {
			input: `switch foo case 1234 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewIntLiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "1234"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be an int64": {
			input: `switch foo case 1i64 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewInt64LiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a uint64": {
			input: `switch foo case 1u64 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewUInt64LiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be an int32": {
			input: `switch foo case 1i32 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewInt32LiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a uint32": {
			input: `switch foo case 1u32 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewUInt32LiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be an int16": {
			input: `switch foo case 1i16 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewInt16LiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a uint16": {
			input: `switch foo case 1u16 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewUInt16LiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be an int8": {
			input: `switch foo case 12i8 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewInt8LiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "12"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a uint8": {
			input: `switch foo case 12u8 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewUInt8LiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "12"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a float": {
			input: `switch foo case 12.5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewFloatLiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "12.5"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a big float": {
			input: `switch foo case 12bf then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewBigFloatLiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "12"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a float64": {
			input: `switch foo case 1f64 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewFloat64LiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a float32": {
			input: `switch foo case 1f32 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewFloat32LiteralNode(S(P(16, 1, 17), P(19, 1, 20)), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"pattern can be a negative float32": {
			input: `switch foo case -1f32 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(29, 1, 30)),
									ast.NewUnaryExpressionNode(
										S(P(16, 1, 17), P(20, 1, 21)),
										T(S(P(16, 1, 17), P(16, 1, 17)), token.MINUS),
										ast.NewFloat32LiteralNode(S(P(17, 1, 18), P(20, 1, 21)), "1"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(27, 1, 28), P(29, 1, 30)),
											ast.NewNilLiteralNode(S(P(27, 1, 28), P(29, 1, 30))),
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
		"pattern public constant": {
			input: `switch foo case Foo then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 1, 32)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 1, 32)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(31, 1, 32)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(27, 1, 28)),
									ast.NewPublicConstantNode(S(P(16, 1, 17), P(18, 1, 19)), "Foo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(25, 1, 26), P(27, 1, 28)),
											ast.NewNilLiteralNode(S(P(25, 1, 26), P(27, 1, 28))),
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
		"pattern private constant": {
			input: `switch foo case _Fo then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 1, 32)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 1, 32)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(31, 1, 32)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(27, 1, 28)),
									ast.NewPrivateConstantNode(S(P(16, 1, 17), P(18, 1, 19)), "_Fo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(25, 1, 26), P(27, 1, 28)),
											ast.NewNilLiteralNode(S(P(25, 1, 26), P(27, 1, 28))),
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
		"pattern root constant lookup": {
			input: `switch foo case ::Foo then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(29, 1, 30)),
									ast.NewConstantLookupNode(
										S(P(16, 1, 17), P(20, 1, 21)),
										nil,
										ast.NewPublicConstantNode(S(P(18, 1, 19), P(20, 1, 21)), "Foo"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(27, 1, 28), P(29, 1, 30)),
											ast.NewNilLiteralNode(S(P(27, 1, 28), P(29, 1, 30))),
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
		"pattern constant lookup": {
			input: `switch foo case Foo::Bar then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(36, 1, 37)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(36, 1, 37)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(36, 1, 37)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(32, 1, 33)),
									ast.NewConstantLookupNode(
										S(P(16, 1, 17), P(23, 1, 24)),
										ast.NewPublicConstantNode(S(P(16, 1, 17), P(18, 1, 19)), "Foo"),
										ast.NewPublicConstantNode(S(P(21, 1, 22), P(23, 1, 24)), "Bar"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(30, 1, 31), P(32, 1, 33)),
											ast.NewNilLiteralNode(S(P(30, 1, 31), P(32, 1, 33))),
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
		"pattern nested constant lookup": {
			input: `switch foo case ::Foo::Bar then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 1, 39)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(38, 1, 39)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(38, 1, 39)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(34, 1, 35)),
									ast.NewConstantLookupNode(
										S(P(16, 1, 17), P(25, 1, 26)),
										ast.NewConstantLookupNode(
											S(P(16, 1, 17), P(20, 1, 21)),
											nil,
											ast.NewPublicConstantNode(S(P(18, 1, 19), P(20, 1, 21)), "Foo"),
										),
										ast.NewPublicConstantNode(S(P(23, 1, 24), P(25, 1, 26)), "Bar"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(32, 1, 33), P(34, 1, 35)),
											ast.NewNilLiteralNode(S(P(32, 1, 33), P(34, 1, 35))),
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
		"unary pattern less": {
			input: `switch foo case < 5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 1, 32)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 1, 32)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(31, 1, 32)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(27, 1, 28)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(18, 1, 19)),
										T(S(P(16, 1, 17), P(16, 1, 17)), token.LESS),
										ast.NewIntLiteralNode(S(P(18, 1, 19), P(18, 1, 19)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(25, 1, 26), P(27, 1, 28)),
											ast.NewNilLiteralNode(S(P(25, 1, 26), P(27, 1, 28))),
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
		"unary pattern less public constant": {
			input: `switch foo case < Foo then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(29, 1, 30)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(20, 1, 21)),
										T(S(P(16, 1, 17), P(16, 1, 17)), token.LESS),
										ast.NewPublicConstantNode(
											S(P(18, 1, 19), P(20, 1, 21)),
											"Foo",
										),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(27, 1, 28), P(29, 1, 30)),
											ast.NewNilLiteralNode(S(P(27, 1, 28), P(29, 1, 30))),
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
		"unary pattern greater": {
			input: `switch foo case > 5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 1, 32)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 1, 32)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(31, 1, 32)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(27, 1, 28)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(18, 1, 19)),
										T(S(P(16, 1, 17), P(16, 1, 17)), token.GREATER),
										ast.NewIntLiteralNode(S(P(18, 1, 19), P(18, 1, 19)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(25, 1, 26), P(27, 1, 28)),
											ast.NewNilLiteralNode(S(P(25, 1, 26), P(27, 1, 28))),
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
		"unary pattern less equal": {
			input: `switch foo case <= 5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(19, 1, 20)),
										T(S(P(16, 1, 17), P(17, 1, 18)), token.LESS_EQUAL),
										ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"unary pattern greater equal": {
			input: `switch foo case >= 5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(19, 1, 20)),
										T(S(P(16, 1, 17), P(17, 1, 18)), token.GREATER_EQUAL),
										ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"unary pattern equal": {
			input: `switch foo case == 5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(19, 1, 20)),
										T(S(P(16, 1, 17), P(17, 1, 18)), token.EQUAL_EQUAL),
										ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"unary pattern not equal": {
			input: `switch foo case != 5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(19, 1, 20)),
										T(S(P(16, 1, 17), P(17, 1, 18)), token.NOT_EQUAL),
										ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"unary pattern lax equal": {
			input: `switch foo case =~ 5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(19, 1, 20)),
										T(S(P(16, 1, 17), P(17, 1, 18)), token.LAX_EQUAL),
										ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"unary pattern lax not equal": {
			input: `switch foo case !~ 5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(19, 1, 20)),
										T(S(P(16, 1, 17), P(17, 1, 18)), token.LAX_NOT_EQUAL),
										ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"unary pattern strict equal": {
			input: `switch foo case === 5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(29, 1, 30)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(20, 1, 21)),
										T(S(P(16, 1, 17), P(18, 1, 19)), token.STRICT_EQUAL),
										ast.NewIntLiteralNode(S(P(20, 1, 21), P(20, 1, 21)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(27, 1, 28), P(29, 1, 30)),
											ast.NewNilLiteralNode(S(P(27, 1, 28), P(29, 1, 30))),
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
		"unary pattern strict not equal": {
			input: `switch foo case !== 5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(29, 1, 30)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(20, 1, 21)),
										T(S(P(16, 1, 17), P(18, 1, 19)), token.STRICT_NOT_EQUAL),
										ast.NewIntLiteralNode(S(P(20, 1, 21), P(20, 1, 21)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(27, 1, 28), P(29, 1, 30)),
											ast.NewNilLiteralNode(S(P(27, 1, 28), P(29, 1, 30))),
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
		"unary pattern with unary minus": {
			input: `switch foo case !== -5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(34, 1, 35)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(34, 1, 35)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(34, 1, 35)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(30, 1, 31)),
									ast.NewUnaryPatternNode(
										S(P(16, 1, 17), P(21, 1, 22)),
										T(S(P(16, 1, 17), P(18, 1, 19)), token.STRICT_NOT_EQUAL),
										ast.NewUnaryExpressionNode(
											S(P(20, 1, 21), P(21, 1, 22)),
											T(S(P(20, 1, 21), P(20, 1, 21)), token.MINUS),
											ast.NewIntLiteralNode(S(P(21, 1, 22), P(21, 1, 22)), "5"),
										),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(28, 1, 29), P(30, 1, 31)),
											ast.NewNilLiteralNode(S(P(28, 1, 29), P(30, 1, 31))),
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
		"binary and pattern": {
			input: `switch foo case > 5 && < 10 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(39, 1, 40)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(39, 1, 40)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(39, 1, 40)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(35, 1, 36)),
									ast.NewBinaryPatternNode(
										S(P(16, 1, 17), P(26, 1, 27)),
										T(S(P(20, 1, 21), P(21, 1, 22)), token.AND_AND),
										ast.NewUnaryPatternNode(
											S(P(16, 1, 17), P(18, 1, 19)),
											T(S(P(16, 1, 17), P(16, 1, 17)), token.GREATER),
											ast.NewIntLiteralNode(S(P(18, 1, 19), P(18, 1, 19)), "5"),
										),
										ast.NewUnaryPatternNode(
											S(P(23, 1, 24), P(26, 1, 27)),
											T(S(P(23, 1, 24), P(23, 1, 24)), token.LESS),
											ast.NewIntLiteralNode(S(P(25, 1, 26), P(26, 1, 27)), "10"),
										),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(33, 1, 34), P(35, 1, 36)),
											ast.NewNilLiteralNode(S(P(33, 1, 34), P(35, 1, 36))),
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
		"binary or pattern": {
			input: `switch foo case > 5 || 2 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(36, 1, 37)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(36, 1, 37)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(36, 1, 37)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(32, 1, 33)),
									ast.NewBinaryPatternNode(
										S(P(16, 1, 17), P(23, 1, 24)),
										T(S(P(20, 1, 21), P(21, 1, 22)), token.OR_OR),
										ast.NewUnaryPatternNode(
											S(P(16, 1, 17), P(18, 1, 19)),
											T(S(P(16, 1, 17), P(16, 1, 17)), token.GREATER),
											ast.NewIntLiteralNode(S(P(18, 1, 19), P(18, 1, 19)), "5"),
										),
										ast.NewIntLiteralNode(S(P(23, 1, 24), P(23, 1, 24)), "2"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(30, 1, 31), P(32, 1, 33)),
											ast.NewNilLiteralNode(S(P(30, 1, 31), P(32, 1, 33))),
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
		"closed range": {
			input: `switch foo case 2...5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(29, 1, 30)),
									ast.NewRangeLiteralNode(
										S(P(16, 1, 17), P(20, 1, 21)),
										T(S(P(17, 1, 18), P(19, 1, 20)), token.CLOSED_RANGE_OP),
										ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "2"),
										ast.NewIntLiteralNode(S(P(20, 1, 21), P(20, 1, 21)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(27, 1, 28), P(29, 1, 30)),
											ast.NewNilLiteralNode(S(P(27, 1, 28), P(29, 1, 30))),
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
		"closed range constant": {
			input: `switch foo case A...B then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(29, 1, 30)),
									ast.NewRangeLiteralNode(
										S(P(16, 1, 17), P(20, 1, 21)),
										T(S(P(17, 1, 18), P(19, 1, 20)), token.CLOSED_RANGE_OP),
										ast.NewPublicConstantNode(S(P(16, 1, 17), P(16, 1, 17)), "A"),
										ast.NewPublicConstantNode(S(P(20, 1, 21), P(20, 1, 21)), "B"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(27, 1, 28), P(29, 1, 30)),
											ast.NewNilLiteralNode(S(P(27, 1, 28), P(29, 1, 30))),
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
		"closed range with unary plus and minus": {
			input: `switch foo case -2...+5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(35, 1, 36)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(35, 1, 36)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(35, 1, 36)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(31, 1, 32)),
									ast.NewRangeLiteralNode(
										S(P(16, 1, 17), P(22, 1, 23)),
										T(S(P(18, 1, 19), P(20, 1, 21)), token.CLOSED_RANGE_OP),
										ast.NewUnaryExpressionNode(
											S(P(16, 1, 17), P(17, 1, 18)),
											T(S(P(16, 1, 17), P(16, 1, 17)), token.MINUS),
											ast.NewIntLiteralNode(S(P(17, 1, 18), P(17, 1, 18)), "2"),
										),
										ast.NewUnaryExpressionNode(
											S(P(21, 1, 22), P(22, 1, 23)),
											T(S(P(21, 1, 22), P(21, 1, 22)), token.PLUS),
											ast.NewIntLiteralNode(S(P(22, 1, 23), P(22, 1, 23)), "5"),
										),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(29, 1, 30), P(31, 1, 32)),
											ast.NewNilLiteralNode(S(P(29, 1, 30), P(31, 1, 32))),
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
		"open range": {
			input: `switch foo case 2<.<5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(29, 1, 30)),
									ast.NewRangeLiteralNode(
										S(P(16, 1, 17), P(20, 1, 21)),
										T(S(P(17, 1, 18), P(19, 1, 20)), token.OPEN_RANGE_OP),
										ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "2"),
										ast.NewIntLiteralNode(S(P(20, 1, 21), P(20, 1, 21)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(27, 1, 28), P(29, 1, 30)),
											ast.NewNilLiteralNode(S(P(27, 1, 28), P(29, 1, 30))),
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
		"left open range": {
			input: `switch foo case 2<..5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(29, 1, 30)),
									ast.NewRangeLiteralNode(
										S(P(16, 1, 17), P(20, 1, 21)),
										T(S(P(17, 1, 18), P(19, 1, 20)), token.LEFT_OPEN_RANGE_OP),
										ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "2"),
										ast.NewIntLiteralNode(S(P(20, 1, 21), P(20, 1, 21)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(27, 1, 28), P(29, 1, 30)),
											ast.NewNilLiteralNode(S(P(27, 1, 28), P(29, 1, 30))),
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
		"right open range": {
			input: `switch foo case 2..<5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(29, 1, 30)),
									ast.NewRangeLiteralNode(
										S(P(16, 1, 17), P(20, 1, 21)),
										T(S(P(17, 1, 18), P(19, 1, 20)), token.RIGHT_OPEN_RANGE_OP),
										ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "2"),
										ast.NewIntLiteralNode(S(P(20, 1, 21), P(20, 1, 21)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(27, 1, 28), P(29, 1, 30)),
											ast.NewNilLiteralNode(S(P(27, 1, 28), P(29, 1, 30))),
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
		"beginless closed range": {
			input: `switch foo case ...5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewRangeLiteralNode(
										S(P(16, 1, 17), P(19, 1, 20)),
										T(S(P(16, 1, 17), P(18, 1, 19)), token.CLOSED_RANGE_OP),
										nil,
										ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"beginless open range": {
			input: `switch foo case ..<5 then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewRangeLiteralNode(
										S(P(16, 1, 17), P(19, 1, 20)),
										T(S(P(16, 1, 17), P(18, 1, 19)), token.RIGHT_OPEN_RANGE_OP),
										nil,
										ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"endless closed range": {
			input: `switch foo case 2... then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewRangeLiteralNode(
										S(P(16, 1, 17), P(19, 1, 20)),
										T(S(P(17, 1, 18), P(19, 1, 20)), token.CLOSED_RANGE_OP),
										ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "2"),
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"endless open range": {
			input: `switch foo case 2<.. then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(28, 1, 29)),
									ast.NewRangeLiteralNode(
										S(P(16, 1, 17), P(19, 1, 20)),
										T(S(P(17, 1, 18), P(19, 1, 20)), token.LEFT_OPEN_RANGE_OP),
										ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "2"),
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(26, 1, 27), P(28, 1, 29)),
											ast.NewNilLiteralNode(S(P(26, 1, 27), P(28, 1, 29))),
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
		"empty list pattern": {
			input: `switch foo case [] then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 1, 31)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(30, 1, 31)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(30, 1, 31)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(26, 1, 27)),
									ast.NewListPatternNode(
										S(P(16, 1, 17), P(17, 1, 18)),
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(24, 1, 25), P(26, 1, 27)),
											ast.NewNilLiteralNode(S(P(24, 1, 25), P(26, 1, 27))),
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
		"list with subpatterns": {
			input: `switch foo case [a, > 6 && < 20, [b, :foo]] then nil end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(55, 1, 56)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(55, 1, 56)),
						ast.NewSwitchExpressionNode(
							S(P(0, 1, 1), P(55, 1, 56)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									S(P(11, 1, 12), P(51, 1, 52)),
									ast.NewListPatternNode(
										S(P(16, 1, 17), P(42, 1, 43)),
										[]ast.PatternNode{
											ast.NewPublicIdentifierNode(
												S(P(17, 1, 18), P(17, 1, 18)),
												"a",
											),
											ast.NewBinaryPatternNode(
												S(P(20, 1, 21), P(30, 1, 31)),
												T(S(P(24, 1, 25), P(25, 1, 26)), token.AND_AND),
												ast.NewUnaryPatternNode(
													S(P(20, 1, 21), P(22, 1, 23)),
													T(S(P(20, 1, 21), P(20, 1, 21)), token.GREATER),
													ast.NewIntLiteralNode(S(P(22, 1, 23), P(22, 1, 23)), "6"),
												),
												ast.NewUnaryPatternNode(
													S(P(27, 1, 28), P(30, 1, 31)),
													T(S(P(27, 1, 28), P(27, 1, 28)), token.LESS),
													ast.NewIntLiteralNode(S(P(29, 1, 30), P(30, 1, 31)), "20"),
												),
											),
											ast.NewListPatternNode(
												S(P(33, 1, 34), P(41, 1, 42)),
												[]ast.PatternNode{
													ast.NewPublicIdentifierNode(S(P(34, 1, 35), P(34, 1, 35)), "b"),
													ast.NewSimpleSymbolLiteralNode(S(P(37, 1, 38), P(40, 1, 41)), "foo"),
												},
											),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											S(P(49, 1, 50), P(51, 1, 52)),
											ast.NewNilLiteralNode(S(P(49, 1, 50), P(51, 1, 52))),
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestUntil(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
until foo > 0
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(33, 5, 4)),
						ast.NewUntilExpressionNode(
							S(P(1, 2, 1), P(32, 5, 3)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 13)),
								T(S(P(11, 2, 11), P(11, 2, 11)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 9)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 13)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(16, 3, 2), P(24, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(16, 3, 2), P(23, 3, 9)),
										T(S(P(20, 3, 6), P(21, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(16, 3, 2), P(18, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(23, 3, 9), P(23, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(26, 4, 2), P(29, 4, 5)),
									ast.NewNilLiteralNode(S(P(26, 4, 2), P(28, 4, 4))),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
until foo > 0
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(18, 3, 4)),
						ast.NewUntilExpressionNode(
							S(P(1, 2, 1), P(17, 3, 3)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 13)),
								T(S(P(11, 2, 11), P(11, 2, 11)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 9)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 13)), "0"),
							),
							nil,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
bar =
	until foo > 0
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(41, 6, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(37, 5, 5)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(36, 5, 4)),
							T(S(P(5, 2, 5), P(5, 2, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 3)), "bar"),
							ast.NewUntilExpressionNode(
								S(P(8, 3, 2), P(36, 5, 4)),
								ast.NewBinaryExpressionNode(
									S(P(14, 3, 8), P(20, 3, 14)),
									T(S(P(18, 3, 12), P(18, 3, 12)), token.GREATER),
									ast.NewPublicIdentifierNode(S(P(14, 3, 8), P(16, 3, 10)), "foo"),
									ast.NewIntLiteralNode(S(P(20, 3, 14), P(20, 3, 14)), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(24, 4, 3), P(32, 4, 11)),
										ast.NewAssignmentExpressionNode(
											S(P(24, 4, 3), P(31, 4, 10)),
											T(S(P(28, 4, 7), P(29, 4, 8)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(24, 4, 3), P(26, 4, 5)), "foo"),
											ast.NewIntLiteralNode(S(P(31, 4, 10), P(31, 4, 10)), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(38, 6, 1), P(41, 6, 4)),
						ast.NewNilLiteralNode(S(P(38, 6, 1), P(40, 6, 3))),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
until foo > 0 then foo += 2
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(28, 2, 28)),
						ast.NewUntilExpressionNode(
							S(P(1, 2, 1), P(27, 2, 27)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 13)),
								T(S(P(11, 2, 11), P(11, 2, 11)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 9)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 13)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(20, 2, 20), P(27, 2, 27)),
									ast.NewAssignmentExpressionNode(
										S(P(20, 2, 20), P(27, 2, 27)),
										T(S(P(24, 2, 24), P(25, 2, 25)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(20, 2, 20), P(22, 2, 22)), "foo"),
										ast.NewIntLiteralNode(S(P(27, 2, 27), P(27, 2, 27)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(29, 3, 1), P(32, 3, 4)),
						ast.NewNilLiteralNode(S(P(29, 3, 1), P(31, 3, 3))),
					),
				},
			),
		},
		"cannot have else": {
			input: `
until foo > 0
	foo += 2
	nil
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(57, 9, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(53, 8, 4)),
						ast.NewUntilExpressionNode(
							S(P(1, 2, 1), P(52, 8, 3)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 13)),
								T(S(P(11, 2, 11), P(11, 2, 11)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 9)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 13)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(16, 3, 2), P(24, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(16, 3, 2), P(23, 3, 9)),
										T(S(P(20, 3, 6), P(21, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(16, 3, 2), P(18, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(23, 3, 9), P(23, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(26, 4, 2), P(29, 4, 5)),
									ast.NewNilLiteralNode(S(P(26, 4, 2), P(28, 4, 4))),
								),
								ast.NewExpressionStatementNode(
									S(P(30, 5, 1), P(34, 5, 5)),
									ast.NewInvalidNode(S(P(30, 5, 1), P(33, 5, 4)), T(S(P(30, 5, 1), P(33, 5, 4)), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									S(P(36, 6, 2), P(44, 6, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(36, 6, 2), P(43, 6, 9)),
										T(S(P(40, 6, 6), P(41, 6, 7)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(36, 6, 2), P(38, 6, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(43, 6, 9), P(43, 6, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(46, 7, 2), P(49, 7, 5)),
									ast.NewNilLiteralNode(S(P(46, 7, 2), P(48, 7, 4))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(54, 9, 1), P(57, 9, 4)),
						ast.NewNilLiteralNode(S(P(54, 9, 1), P(56, 9, 3))),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(30, 5, 1), P(33, 5, 4)), "unexpected else, expected an expression"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestLoop(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
loop
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(24, 5, 4)),
						ast.NewLoopExpressionNode(
							S(P(1, 2, 1), P(23, 5, 3)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(7, 3, 2), P(15, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(7, 3, 2), P(14, 3, 9)),
										T(S(P(11, 3, 6), P(12, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(7, 3, 2), P(9, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(14, 3, 9), P(14, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(17, 4, 2), P(20, 4, 5)),
									ast.NewNilLiteralNode(S(P(17, 4, 2), P(19, 4, 4))),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
loop
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(9, 3, 4)),
						ast.NewLoopExpressionNode(
							S(P(1, 2, 1), P(8, 3, 3)),
							nil,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
bar =
	loop
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 6, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(28, 5, 5)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(27, 5, 4)),
							T(S(P(5, 2, 5), P(5, 2, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 3)), "bar"),
							ast.NewLoopExpressionNode(
								S(P(8, 3, 2), P(27, 5, 4)),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(15, 4, 3), P(23, 4, 11)),
										ast.NewAssignmentExpressionNode(
											S(P(15, 4, 3), P(22, 4, 10)),
											T(S(P(19, 4, 7), P(20, 4, 8)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(15, 4, 3), P(17, 4, 5)), "foo"),
											ast.NewIntLiteralNode(S(P(22, 4, 10), P(22, 4, 10)), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(29, 6, 1), P(32, 6, 4)),
						ast.NewNilLiteralNode(S(P(29, 6, 1), P(31, 6, 3))),
					),
				},
			),
		},
		"can be single line without end": {
			input: `
loop foo += 2
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(14, 2, 14)),
						ast.NewLoopExpressionNode(
							S(P(1, 2, 1), P(13, 2, 13)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(6, 2, 6), P(13, 2, 13)),
									ast.NewAssignmentExpressionNode(
										S(P(6, 2, 6), P(13, 2, 13)),
										T(S(P(10, 2, 10), P(11, 2, 11)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(6, 2, 6), P(8, 2, 8)), "foo"),
										ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 13)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(15, 3, 1), P(18, 3, 4)),
						ast.NewNilLiteralNode(S(P(15, 3, 1), P(17, 3, 3))),
					),
				},
			),
		},
		"cannot have else": {
			input: `
loop
	foo += 2
	nil
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(48, 9, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(44, 8, 4)),
						ast.NewLoopExpressionNode(
							S(P(1, 2, 1), P(43, 8, 3)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(7, 3, 2), P(15, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(7, 3, 2), P(14, 3, 9)),
										T(S(P(11, 3, 6), P(12, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(7, 3, 2), P(9, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(14, 3, 9), P(14, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(17, 4, 2), P(20, 4, 5)),
									ast.NewNilLiteralNode(S(P(17, 4, 2), P(19, 4, 4))),
								),
								ast.NewExpressionStatementNode(
									S(P(21, 5, 1), P(25, 5, 5)),
									ast.NewInvalidNode(S(P(21, 5, 1), P(24, 5, 4)), T(S(P(21, 5, 1), P(24, 5, 4)), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									S(P(27, 6, 2), P(35, 6, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(27, 6, 2), P(34, 6, 9)),
										T(S(P(31, 6, 6), P(32, 6, 7)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(27, 6, 2), P(29, 6, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(34, 6, 9), P(34, 6, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(37, 7, 2), P(40, 7, 5)),
									ast.NewNilLiteralNode(S(P(37, 7, 2), P(39, 7, 4))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(45, 9, 1), P(48, 9, 4)),
						ast.NewNilLiteralNode(S(P(45, 9, 1), P(47, 9, 3))),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(21, 5, 1), P(24, 5, 4)), "unexpected else, expected an expression"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestDo(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
do
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(22, 5, 4)),
						ast.NewDoExpressionNode(
							S(P(1, 2, 1), P(21, 5, 3)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(5, 3, 2), P(13, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(5, 3, 2), P(12, 3, 9)),
										T(S(P(9, 3, 6), P(10, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(5, 3, 2), P(7, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(12, 3, 9), P(12, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(15, 4, 2), P(18, 4, 5)),
									ast.NewNilLiteralNode(S(P(15, 4, 2), P(17, 4, 4))),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
do
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(7, 3, 4)),
						ast.NewDoExpressionNode(
							S(P(1, 2, 1), P(6, 3, 3)),
							nil,
						),
					),
				},
			),
		},
		"can be a one-liner": {
			input: "do 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewDoExpressionNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(3, 1, 4), P(3, 1, 4)),
									ast.NewIntLiteralNode(
										S(P(3, 1, 4), P(3, 1, 4)),
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
	do
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 6, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(26, 5, 5)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(25, 5, 4)),
							T(S(P(5, 2, 5), P(5, 2, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 3)), "bar"),
							ast.NewDoExpressionNode(
								S(P(8, 3, 2), P(25, 5, 4)),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(13, 4, 3), P(21, 4, 11)),
										ast.NewAssignmentExpressionNode(
											S(P(13, 4, 3), P(20, 4, 10)),
											T(S(P(17, 4, 7), P(18, 4, 8)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(13, 4, 3), P(15, 4, 5)), "foo"),
											ast.NewIntLiteralNode(S(P(20, 4, 10), P(20, 4, 10)), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(27, 6, 1), P(30, 6, 4)),
						ast.NewNilLiteralNode(S(P(27, 6, 1), P(29, 6, 3))),
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

func TestBreak(t *testing.T) {
	tests := testTable{
		"can stand alone": {
			input: `break`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewBreakExpressionNode(S(P(0, 1, 1), P(4, 1, 5)), "", nil),
					),
				},
			),
		},
		"can have a label": {
			input: `break$foo`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBreakExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							"foo",
							nil,
						),
					),
				},
			),
		},
		"can have a modifier if without an argument": {
			input: `break if true`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewModifierNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							T(S(P(6, 1, 7), P(7, 1, 8)), token.IF),
							ast.NewBreakExpressionNode(S(P(0, 1, 1), P(4, 1, 5)), "", nil),
							ast.NewTrueLiteralNode(S(P(9, 1, 10), P(12, 1, 13))),
						),
					),
				},
			),
		},
		"can have a modifier if with an argument": {
			input: `break :foo if true`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewModifierNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							T(S(P(11, 1, 12), P(12, 1, 13)), token.IF),
							ast.NewBreakExpressionNode(
								S(P(0, 1, 1), P(9, 1, 10)),
								"",
								ast.NewSimpleSymbolLiteralNode(
									S(P(6, 1, 7), P(9, 1, 10)),
									"foo",
								),
							),
							ast.NewTrueLiteralNode(S(P(14, 1, 15), P(17, 1, 18))),
						),
					),
				},
			),
		},
		"can have an argument": {
			input: `break 2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewBreakExpressionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							"",
							ast.NewIntLiteralNode(
								S(P(6, 1, 7), P(6, 1, 7)),
								"2",
							),
						),
					),
				},
			),
		},
		"can have a label and argument": {
			input: `break$foo 2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewBreakExpressionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							"foo",
							ast.NewIntLiteralNode(
								S(P(10, 1, 11), P(10, 1, 11)),
								"2",
							),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && break`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewLogicalExpressionNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							T(S(P(4, 1, 5), P(5, 1, 6)), token.AND_AND),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewBreakExpressionNode(S(P(7, 1, 8), P(11, 1, 12)), "", nil),
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

func TestReturn(t *testing.T) {
	tests := testTable{
		"can stand alone at the end": {
			input: `return`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewReturnExpressionNode(S(P(0, 1, 1), P(5, 1, 6)), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "return\n1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 2, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewReturnExpressionNode(S(P(0, 1, 1), P(5, 1, 6)), nil),
					),
					ast.NewExpressionStatementNode(
						S(P(7, 2, 1), P(7, 2, 1)),
						ast.NewIntLiteralNode(S(P(7, 2, 1), P(7, 2, 1)), "1"),
					),
				},
			),
		},
		"can have a modifier if without an argument": {
			input: `return if true`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewModifierNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							T(S(P(7, 1, 8), P(8, 1, 9)), token.IF),
							ast.NewReturnExpressionNode(S(P(0, 1, 1), P(5, 1, 6)), nil),
							ast.NewTrueLiteralNode(S(P(10, 1, 11), P(13, 1, 14))),
						),
					),
				},
			),
		},
		"can have a modifier if with an argument": {
			input: `return :foo if true`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewModifierNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							T(S(P(12, 1, 13), P(13, 1, 14)), token.IF),
							ast.NewReturnExpressionNode(
								S(P(0, 1, 1), P(10, 1, 11)),
								ast.NewSimpleSymbolLiteralNode(
									S(P(7, 1, 8), P(10, 1, 11)),
									"foo",
								),
							),
							ast.NewTrueLiteralNode(S(P(15, 1, 16), P(18, 1, 19))),
						),
					),
				},
			),
		},
		"can have an argument": {
			input: `return 2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewReturnExpressionNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							ast.NewIntLiteralNode(S(P(7, 1, 8), P(7, 1, 8)), "2"),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && return`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewLogicalExpressionNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							T(S(P(4, 1, 5), P(5, 1, 6)), token.AND_AND),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewReturnExpressionNode(S(P(7, 1, 8), P(12, 1, 13)), nil),
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

func TestContinue(t *testing.T) {
	tests := testTable{
		"can stand alone at the end": {
			input: `continue`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewContinueExpressionNode(S(P(0, 1, 1), P(7, 1, 8)), "", nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "continue\n1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 2, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewContinueExpressionNode(S(P(0, 1, 1), P(7, 1, 8)), "", nil),
					),
					ast.NewExpressionStatementNode(
						S(P(9, 2, 1), P(9, 2, 1)),
						ast.NewIntLiteralNode(S(P(9, 2, 1), P(9, 2, 1)), "1"),
					),
				},
			),
		},
		"can have a label": {
			input: `continue$foo`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewContinueExpressionNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							"foo",
							nil,
						),
					),
				},
			),
		},
		"can have an argument": {
			input: `continue 2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewContinueExpressionNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							"",
							ast.NewIntLiteralNode(S(P(9, 1, 10), P(9, 1, 10)), "2"),
						),
					),
				},
			),
		},
		"can have a modifier if without an argument": {
			input: `continue if true`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewModifierNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							T(S(P(9, 1, 10), P(10, 1, 11)), token.IF),
							ast.NewContinueExpressionNode(S(P(0, 1, 1), P(7, 1, 8)), "", nil),
							ast.NewTrueLiteralNode(S(P(12, 1, 13), P(15, 1, 16))),
						),
					),
				},
			),
		},
		"can have a modifier if with an argument": {
			input: `continue :foo if true`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewModifierNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							T(S(P(14, 1, 15), P(15, 1, 16)), token.IF),
							ast.NewContinueExpressionNode(
								S(P(0, 1, 1), P(12, 1, 13)),
								"",
								ast.NewSimpleSymbolLiteralNode(
									S(P(9, 1, 10), P(12, 1, 13)),
									"foo",
								),
							),
							ast.NewTrueLiteralNode(S(P(17, 1, 18), P(20, 1, 21))),
						),
					),
				},
			),
		},
		"can have a label and argument": {
			input: `continue$foo 2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewContinueExpressionNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							"foo",
							ast.NewIntLiteralNode(
								S(P(13, 1, 14), P(13, 1, 14)),
								"2",
							),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && continue`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewLogicalExpressionNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							T(S(P(4, 1, 5), P(5, 1, 6)), token.AND_AND),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewContinueExpressionNode(S(P(7, 1, 8), P(14, 1, 15)), "", nil),
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

func TestThrow(t *testing.T) {
	tests := testTable{
		"can stand alone at the end": {
			input: `throw`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewThrowExpressionNode(S(P(0, 1, 1), P(4, 1, 5)), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "throw\n1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 2, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewThrowExpressionNode(S(P(0, 1, 1), P(4, 1, 5)), nil),
					),
					ast.NewExpressionStatementNode(
						S(P(6, 2, 1), P(6, 2, 1)),
						ast.NewIntLiteralNode(S(P(6, 2, 1), P(6, 2, 1)), "1"),
					),
				},
			),
		},
		"can have an argument": {
			input: `throw 2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewThrowExpressionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							ast.NewIntLiteralNode(S(P(6, 1, 7), P(6, 1, 7)), "2"),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && throw`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewLogicalExpressionNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							T(S(P(4, 1, 5), P(5, 1, 6)), token.AND_AND),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewThrowExpressionNode(S(P(7, 1, 8), P(11, 1, 12)), nil),
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

func TestForIn(t *testing.T) {
	tests := testTable{
		"can be single-line with then": {
			input: `for i in [1, 2, 3] then println(i)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewForInExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							ast.NewPublicIdentifierNode(S(P(4, 1, 5), P(4, 1, 5)), "i"),
							ast.NewArrayListLiteralNode(
								S(P(9, 1, 10), P(17, 1, 18)),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(S(P(10, 1, 11), P(10, 1, 11)), "1"),
									ast.NewIntLiteralNode(S(P(13, 1, 14), P(13, 1, 14)), "2"),
									ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "3"),
								},
								nil,
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(24, 1, 25), P(33, 1, 34)),
									ast.NewFunctionCallNode(
										S(P(24, 1, 25), P(33, 1, 34)),
										"println",
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(S(P(32, 1, 33), P(32, 1, 33)), "i"),
										},
										nil,
									),
								),
							},
						),
					),
				},
			),
		},
		"can be multiline": {
			input: `for i in [1, 2, 3]
  println(i)
  nil
end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(40, 4, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 4, 3)),
						ast.NewForInExpressionNode(
							S(P(0, 1, 1), P(40, 4, 3)),
							ast.NewPublicIdentifierNode(S(P(4, 1, 5), P(4, 1, 5)), "i"),
							ast.NewArrayListLiteralNode(
								S(P(9, 1, 10), P(17, 1, 18)),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(S(P(10, 1, 11), P(10, 1, 11)), "1"),
									ast.NewIntLiteralNode(S(P(13, 1, 14), P(13, 1, 14)), "2"),
									ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "3"),
								},
								nil,
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 2, 3), P(31, 2, 13)),
									ast.NewFunctionCallNode(
										S(P(21, 2, 3), P(30, 2, 12)),
										"println",
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(S(P(29, 2, 11), P(29, 2, 11)), "i"),
										},
										nil,
									),
								),
								ast.NewExpressionStatementNode(
									S(P(34, 3, 3), P(37, 3, 6)),
									ast.NewNilLiteralNode(S(P(34, 3, 3), P(36, 3, 5))),
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

func TestNumericFor(t *testing.T) {
	tests := testTable{
		"can be single-line with then": {
			input: `fornum i := 0; i < 5; i += 1 then println(i)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(43, 1, 44)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(43, 1, 44)),
						ast.NewNumericForExpressionNode(
							S(P(0, 1, 1), P(43, 1, 44)),
							ast.NewAssignmentExpressionNode(
								S(P(7, 1, 8), P(12, 1, 13)),
								T(S(P(9, 1, 10), P(10, 1, 11)), token.COLON_EQUAL),
								ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(7, 1, 8)), "i"),
								ast.NewIntLiteralNode(S(P(12, 1, 13), P(12, 1, 13)), "0"),
							),
							ast.NewBinaryExpressionNode(
								S(P(15, 1, 16), P(19, 1, 20)),
								T(S(P(17, 1, 18), P(17, 1, 18)), token.LESS),
								ast.NewPublicIdentifierNode(S(P(15, 1, 16), P(15, 1, 16)), "i"),
								ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "5"),
							),
							ast.NewAssignmentExpressionNode(
								S(P(22, 1, 23), P(27, 1, 28)),
								T(S(P(24, 1, 25), P(25, 1, 26)), token.PLUS_EQUAL),
								ast.NewPublicIdentifierNode(S(P(22, 1, 23), P(22, 1, 23)), "i"),
								ast.NewIntLiteralNode(S(P(27, 1, 28), P(27, 1, 28)), "1"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(34, 1, 35), P(43, 1, 44)),
									ast.NewFunctionCallNode(
										S(P(34, 1, 35), P(43, 1, 44)),
										"println",
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(S(P(42, 1, 43), P(42, 1, 43)), "i"),
										},
										nil,
									),
								),
							},
						),
					),
				},
			),
		},
		"can have empty fields": {
			input: `fornum ;; then println(i)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewNumericForExpressionNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(15, 1, 16), P(24, 1, 25)),
									ast.NewFunctionCallNode(
										S(P(15, 1, 16), P(24, 1, 25)),
										"println",
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(S(P(23, 1, 24), P(23, 1, 24)), "i"),
										},
										nil,
									),
								),
							},
						),
					),
				},
			),
		},
		"can be multiline": {
			input: `fornum i := 0; i < 5; i += 1
  println(i)
  nil
end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(50, 4, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(50, 4, 3)),
						ast.NewNumericForExpressionNode(
							S(P(0, 1, 1), P(50, 4, 3)),
							ast.NewAssignmentExpressionNode(
								S(P(7, 1, 8), P(12, 1, 13)),
								T(S(P(9, 1, 10), P(10, 1, 11)), token.COLON_EQUAL),
								ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(7, 1, 8)), "i"),
								ast.NewIntLiteralNode(S(P(12, 1, 13), P(12, 1, 13)), "0"),
							),
							ast.NewBinaryExpressionNode(
								S(P(15, 1, 16), P(19, 1, 20)),
								T(S(P(17, 1, 18), P(17, 1, 18)), token.LESS),
								ast.NewPublicIdentifierNode(S(P(15, 1, 16), P(15, 1, 16)), "i"),
								ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "5"),
							),
							ast.NewAssignmentExpressionNode(
								S(P(22, 1, 23), P(27, 1, 28)),
								T(S(P(24, 1, 25), P(25, 1, 26)), token.PLUS_EQUAL),
								ast.NewPublicIdentifierNode(S(P(22, 1, 23), P(22, 1, 23)), "i"),
								ast.NewIntLiteralNode(S(P(27, 1, 28), P(27, 1, 28)), "1"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(31, 2, 3), P(41, 2, 13)),
									ast.NewFunctionCallNode(
										S(P(31, 2, 3), P(40, 2, 12)),
										"println",
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(S(P(39, 2, 11), P(39, 2, 11)), "i"),
										},
										nil,
									),
								),
								ast.NewExpressionStatementNode(
									S(P(44, 3, 3), P(47, 3, 6)),
									ast.NewNilLiteralNode(S(P(44, 3, 3), P(46, 3, 5))),
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
