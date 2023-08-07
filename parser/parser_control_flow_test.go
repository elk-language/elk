package parser

import (
	"testing"

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
							ast.NewPublicIdentifierNode(S(P(13, 2, 1), P(15, 2, 16)), "baz"),
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
							ast.NewPublicIdentifierNode(S(P(13, 2, 1), P(15, 2, 16)), "baz"),
							ast.NewAssignmentExpressionNode(
								S(P(22, 3, 1), P(30, 3, 31)),
								T(S(P(26, 3, 5), P(26, 3, 27)), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(S(P(22, 3, 1), P(24, 3, 25)), "car"),
								ast.NewPublicIdentifierNode(S(P(28, 3, 7), P(30, 3, 31)), "red"),
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
							[]ast.ParameterNode{
								ast.NewLoopParameterNode(S(P(15, 1, 16), P(15, 1, 16)), "i", nil),
							},
							ast.NewListLiteralNode(
								S(P(20, 1, 21), P(28, 1, 29)),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(S(P(21, 1, 22), P(21, 1, 22)), "1"),
									ast.NewIntLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "2"),
									ast.NewIntLiteralNode(S(P(27, 1, 28), P(27, 1, 28)), "3"),
								},
							),
						),
					),
				},
			),
		},
		"can have multiple parameters in for loops": {
			input: "println(i) for i, j: Int in [1, 2, 3]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(36, 1, 37)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(36, 1, 37)),
						ast.NewModifierForInNode(
							S(P(0, 1, 1), P(36, 1, 37)),
							ast.NewFunctionCallNode(
								S(P(0, 1, 1), P(9, 1, 10)),
								"println",
								[]ast.ExpressionNode{
									ast.NewPublicIdentifierNode(S(P(8, 1, 9), P(8, 1, 9)), "i"),
								},
								nil,
							),
							[]ast.ParameterNode{
								ast.NewLoopParameterNode(S(P(15, 1, 16), P(15, 1, 16)), "i", nil),
								ast.NewLoopParameterNode(S(P(18, 1, 19), P(23, 1, 24)), "j", ast.NewPublicConstantNode(S(P(21, 1, 22), P(23, 1, 24)), "Int")),
							},
							ast.NewListLiteralNode(
								S(P(28, 1, 29), P(36, 1, 37)),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(S(P(29, 1, 30), P(29, 1, 30)), "1"),
									ast.NewIntLiteralNode(S(P(32, 1, 33), P(32, 1, 33)), "2"),
									ast.NewIntLiteralNode(S(P(35, 1, 36), P(35, 1, 36)), "3"),
								},
							),
						),
					),
				},
			),
		},
		"for loops can span multiple lines": {
			input: "println(i) for\ni,\nj: Int\nin\n[1,\n2,\n3]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(36, 1, 37)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(36, 1, 37)),
						ast.NewModifierForInNode(
							S(P(0, 1, 1), P(36, 1, 37)),
							ast.NewFunctionCallNode(
								S(P(0, 1, 1), P(9, 1, 10)),
								"println",
								[]ast.ExpressionNode{
									ast.NewPublicIdentifierNode(S(P(8, 1, 9), P(8, 1, 9)), "i"),
								},
								nil,
							),
							[]ast.ParameterNode{
								ast.NewLoopParameterNode(S(P(15, 2, 1), P(15, 2, 16)), "i", nil),
								ast.NewLoopParameterNode(S(P(18, 3, 1), P(23, 3, 24)), "j", ast.NewPublicConstantNode(S(P(21, 3, 4), P(23, 3, 24)), "Int")),
							},
							ast.NewListLiteralNode(
								S(P(28, 5, 1), P(36, 5, 37)),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(S(P(29, 5, 2), P(29, 5, 30)), "1"),
									ast.NewIntLiteralNode(S(P(32, 6, 1), P(32, 6, 33)), "2"),
									ast.NewIntLiteralNode(S(P(35, 7, 1), P(35, 7, 36)), "3"),
								},
							),
						),
					),
				},
			),
		},
		"has many versions": {
			input: "foo if bar\nfoo unless bar\nfoo while bar\nfoo until bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(52, 1, 53)),
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
						S(P(11, 2, 1), P(25, 2, 26)),
						ast.NewModifierNode(
							S(P(11, 2, 1), P(24, 2, 25)),
							T(S(P(15, 2, 5), P(20, 2, 21)), token.UNLESS),
							ast.NewPublicIdentifierNode(S(P(11, 2, 1), P(13, 2, 14)), "foo"),
							ast.NewPublicIdentifierNode(S(P(22, 2, 12), P(24, 2, 25)), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(26, 3, 1), P(39, 3, 40)),
						ast.NewModifierNode(
							S(P(26, 3, 1), P(38, 3, 39)),
							T(S(P(30, 3, 5), P(34, 3, 35)), token.WHILE),
							ast.NewPublicIdentifierNode(S(P(26, 3, 1), P(28, 3, 29)), "foo"),
							ast.NewPublicIdentifierNode(S(P(36, 3, 11), P(38, 3, 39)), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(40, 4, 1), P(52, 4, 53)),
						ast.NewModifierNode(
							S(P(40, 4, 1), P(52, 4, 53)),
							T(S(P(44, 4, 5), P(48, 4, 49)), token.UNTIL),
							ast.NewPublicIdentifierNode(S(P(40, 4, 1), P(42, 4, 43)), "foo"),
							ast.NewPublicIdentifierNode(S(P(50, 4, 11), P(52, 4, 53)), "bar"),
						),
					),
				},
			),
		},
		"can't be nested": {
			input: "foo = bar if baz if false\n3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
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
						S(P(26, 2, 1), P(26, 2, 27)),
						ast.NewIntLiteralNode(S(P(26, 2, 1), P(26, 2, 27)), "3"),
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
				S(P(0, 1, 1), P(30, 1, 31)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(30, 2, 31)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(29, 2, 30)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 11)),
								T(S(P(8, 2, 8), P(8, 2, 9)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 7)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 11)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(13, 3, 2), P(21, 3, 22)),
									ast.NewAssignmentExpressionNode(
										S(P(13, 3, 2), P(20, 3, 21)),
										T(S(P(17, 3, 6), P(18, 3, 19)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(13, 3, 2), P(15, 3, 16)), "foo"),
										ast.NewIntLiteralNode(S(P(20, 3, 9), P(20, 3, 21)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(23, 4, 2), P(26, 4, 27)),
									ast.NewNilLiteralNode(S(P(23, 4, 2), P(25, 4, 26))),
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
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(15, 2, 16)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(14, 2, 15)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 11)),
								T(S(P(8, 2, 8), P(8, 2, 9)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 7)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 11)), "0"),
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
				S(P(0, 1, 1), P(38, 1, 39)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(34, 2, 35)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(33, 2, 34)),
							T(S(P(5, 2, 5), P(5, 2, 6)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 4)), "bar"),
							ast.NewIfExpressionNode(
								S(P(8, 3, 2), P(33, 3, 34)),
								ast.NewBinaryExpressionNode(
									S(P(11, 3, 5), P(17, 3, 18)),
									T(S(P(15, 3, 9), P(15, 3, 16)), token.GREATER),
									ast.NewPublicIdentifierNode(S(P(11, 3, 5), P(13, 3, 14)), "foo"),
									ast.NewIntLiteralNode(S(P(17, 3, 11), P(17, 3, 18)), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(21, 4, 3), P(29, 4, 30)),
										ast.NewAssignmentExpressionNode(
											S(P(21, 4, 3), P(28, 4, 29)),
											T(S(P(25, 4, 7), P(26, 4, 27)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(21, 4, 3), P(23, 4, 24)), "foo"),
											ast.NewIntLiteralNode(S(P(28, 4, 10), P(28, 4, 29)), "2"),
										),
									),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(35, 6, 1), P(38, 6, 39)),
						ast.NewNilLiteralNode(S(P(35, 6, 1), P(37, 6, 38))),
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
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(25, 2, 26)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(24, 2, 25)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 11)),
								T(S(P(8, 2, 8), P(8, 2, 9)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 7)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 11)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 2, 17), P(24, 2, 25)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 2, 17), P(24, 2, 25)),
										T(S(P(21, 2, 21), P(22, 2, 23)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 2, 17), P(19, 2, 20)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 2, 24), P(24, 2, 25)), "2"),
									),
								),
							},
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						S(P(26, 3, 1), P(29, 3, 30)),
						ast.NewNilLiteralNode(S(P(26, 3, 1), P(28, 3, 29))),
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
				S(P(0, 1, 1), P(55, 1, 56)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(51, 2, 52)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(50, 2, 51)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 11)),
								T(S(P(8, 2, 8), P(8, 2, 9)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 7)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 11)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(13, 3, 2), P(21, 3, 22)),
									ast.NewAssignmentExpressionNode(
										S(P(13, 3, 2), P(20, 3, 21)),
										T(S(P(17, 3, 6), P(18, 3, 19)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(13, 3, 2), P(15, 3, 16)), "foo"),
										ast.NewIntLiteralNode(S(P(20, 3, 9), P(20, 3, 21)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(23, 4, 2), P(26, 4, 27)),
									ast.NewNilLiteralNode(S(P(23, 4, 2), P(25, 4, 26))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(34, 6, 3), P(42, 6, 43)),
									ast.NewAssignmentExpressionNode(
										S(P(34, 6, 3), P(41, 6, 42)),
										T(S(P(38, 6, 7), P(39, 6, 40)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(34, 6, 3), P(36, 6, 37)), "foo"),
										ast.NewIntLiteralNode(S(P(41, 6, 10), P(41, 6, 42)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(44, 7, 2), P(47, 7, 48)),
									ast.NewNilLiteralNode(S(P(44, 7, 2), P(46, 7, 47))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(52, 9, 1), P(55, 9, 56)),
						ast.NewNilLiteralNode(S(P(52, 9, 1), P(54, 9, 55))),
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
				S(P(0, 1, 1), P(43, 1, 44)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(39, 2, 40)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(38, 2, 39)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 11)),
								T(S(P(8, 2, 8), P(8, 2, 9)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 7)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 11)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 2, 17), P(24, 2, 25)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 2, 17), P(24, 2, 25)),
										T(S(P(21, 2, 21), P(22, 2, 23)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 2, 17), P(19, 2, 20)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 2, 24), P(24, 2, 25)), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(31, 3, 6), P(38, 3, 39)),
									ast.NewAssignmentExpressionNode(
										S(P(31, 3, 6), P(38, 3, 39)),
										T(S(P(35, 3, 10), P(36, 3, 37)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(31, 3, 6), P(33, 3, 34)), "foo"),
										ast.NewIntLiteralNode(S(P(38, 3, 13), P(38, 3, 39)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(40, 4, 1), P(43, 4, 44)),
						ast.NewNilLiteralNode(S(P(40, 4, 1), P(42, 4, 43))),
					),
				},
			),
		},
		"can't have two elses": {
			input: `
if foo > 0 then foo += 2
else foo -= 2
else bar
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(52, 1, 53)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(39, 2, 40)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(38, 2, 39)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 11)),
								T(S(P(8, 2, 8), P(8, 2, 9)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 7)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 11)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 2, 17), P(24, 2, 25)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 2, 17), P(24, 2, 25)),
										T(S(P(21, 2, 21), P(22, 2, 23)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 2, 17), P(19, 2, 20)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 2, 24), P(24, 2, 25)), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(31, 3, 6), P(38, 3, 39)),
									ast.NewAssignmentExpressionNode(
										S(P(31, 3, 6), P(38, 3, 39)),
										T(S(P(35, 3, 10), P(36, 3, 37)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(31, 3, 6), P(33, 3, 34)), "foo"),
										ast.NewIntLiteralNode(S(P(38, 3, 13), P(38, 3, 39)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(40, 4, 1), P(48, 4, 49)),
						ast.NewInvalidNode(S(P(40, 4, 1), P(43, 4, 44)), T(S(P(40, 4, 1), P(43, 4, 44)), token.ELSE)),
					),
					ast.NewExpressionStatementNode(
						S(P(49, 5, 1), P(52, 5, 53)),
						ast.NewNilLiteralNode(S(P(49, 5, 1), P(51, 5, 52))),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(40, 4, 1), P(43, 4, 44)), "unexpected else, expected an expression"),
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
				S(P(0, 1, 1), P(103, 1, 104)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(99, 2, 100)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(98, 2, 99)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 11)),
								T(S(P(8, 2, 8), P(8, 2, 9)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 7)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 11)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(13, 3, 2), P(21, 3, 22)),
									ast.NewAssignmentExpressionNode(
										S(P(13, 3, 2), P(20, 3, 21)),
										T(S(P(17, 3, 6), P(18, 3, 19)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(13, 3, 2), P(15, 3, 16)), "foo"),
										ast.NewIntLiteralNode(S(P(20, 3, 9), P(20, 3, 21)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(23, 4, 2), P(26, 4, 27)),
									ast.NewNilLiteralNode(S(P(23, 4, 2), P(25, 4, 26))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(27, 5, 1), P(51, 5, 52)),
									ast.NewIfExpressionNode(
										S(P(27, 5, 1), P(51, 5, 52)),
										ast.NewBinaryExpressionNode(
											S(P(33, 5, 7), P(39, 5, 40)),
											T(S(P(37, 5, 11), P(37, 5, 38)), token.LESS),
											ast.NewPublicIdentifierNode(S(P(33, 5, 7), P(35, 5, 36)), "foo"),
											ast.NewIntLiteralNode(S(P(39, 5, 13), P(39, 5, 40)), "5"),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(42, 6, 2), P(51, 6, 52)),
												ast.NewAssignmentExpressionNode(
													S(P(42, 6, 2), P(50, 6, 51)),
													T(S(P(46, 6, 6), P(47, 6, 48)), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(S(P(42, 6, 2), P(44, 6, 45)), "foo"),
													ast.NewIntLiteralNode(S(P(49, 6, 9), P(50, 6, 51)), "10"),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(52, 7, 1), P(98, 7, 99)),
												ast.NewIfExpressionNode(
													S(P(52, 7, 1), P(98, 7, 99)),
													ast.NewBinaryExpressionNode(
														S(P(58, 7, 7), P(64, 7, 65)),
														T(S(P(62, 7, 11), P(62, 7, 63)), token.LESS),
														ast.NewPublicIdentifierNode(S(P(58, 7, 7), P(60, 7, 61)), "foo"),
														ast.NewIntLiteralNode(S(P(64, 7, 13), P(64, 7, 65)), "0"),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(67, 8, 2), P(75, 8, 76)),
															ast.NewAssignmentExpressionNode(
																S(P(67, 8, 2), P(74, 8, 75)),
																T(S(P(71, 8, 6), P(72, 8, 73)), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(S(P(67, 8, 2), P(69, 8, 70)), "foo"),
																ast.NewIntLiteralNode(S(P(74, 8, 9), P(74, 8, 75)), "3"),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(82, 10, 2), P(90, 10, 91)),
															ast.NewAssignmentExpressionNode(
																S(P(82, 10, 2), P(89, 10, 90)),
																T(S(P(86, 10, 6), P(87, 10, 88)), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(S(P(82, 10, 2), P(84, 10, 85)), "foo"),
																ast.NewIntLiteralNode(S(P(89, 10, 9), P(89, 10, 90)), "2"),
															),
														),
														ast.NewExpressionStatementNode(
															S(P(92, 11, 2), P(95, 11, 96)),
															ast.NewNilLiteralNode(S(P(92, 11, 2), P(94, 11, 95))),
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
						S(P(100, 13, 1), P(103, 13, 104)),
						ast.NewNilLiteralNode(S(P(100, 13, 1), P(102, 13, 103))),
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
				S(P(0, 1, 1), P(100, 1, 101)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(96, 2, 97)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(95, 2, 96)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 11)),
								T(S(P(8, 2, 8), P(8, 2, 9)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 7)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 11)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 2, 17), P(24, 2, 25)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 2, 17), P(24, 2, 25)),
										T(S(P(21, 2, 21), P(22, 2, 23)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 2, 17), P(19, 2, 20)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 2, 24), P(24, 2, 25)), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(26, 3, 1), P(53, 3, 54)),
									ast.NewIfExpressionNode(
										S(P(26, 3, 1), P(53, 3, 54)),
										ast.NewBinaryExpressionNode(
											S(P(32, 3, 7), P(38, 3, 39)),
											T(S(P(36, 3, 11), P(36, 3, 37)), token.LESS),
											ast.NewPublicIdentifierNode(S(P(32, 3, 7), P(34, 3, 35)), "foo"),
											ast.NewIntLiteralNode(S(P(38, 3, 13), P(38, 3, 39)), "5"),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(45, 3, 20), P(53, 3, 54)),
												ast.NewAssignmentExpressionNode(
													S(P(45, 3, 20), P(53, 3, 54)),
													T(S(P(49, 3, 24), P(50, 3, 51)), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(S(P(45, 3, 20), P(47, 3, 48)), "foo"),
													ast.NewIntLiteralNode(S(P(52, 3, 27), P(53, 3, 54)), "10"),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(55, 4, 1), P(95, 4, 96)),
												ast.NewIfExpressionNode(
													S(P(55, 4, 1), P(95, 4, 96)),
													ast.NewBinaryExpressionNode(
														S(P(61, 4, 7), P(67, 4, 68)),
														T(S(P(65, 4, 11), P(65, 4, 66)), token.LESS),
														ast.NewPublicIdentifierNode(S(P(61, 4, 7), P(63, 4, 64)), "foo"),
														ast.NewIntLiteralNode(S(P(67, 4, 13), P(67, 4, 68)), "0"),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(74, 4, 20), P(81, 4, 82)),
															ast.NewAssignmentExpressionNode(
																S(P(74, 4, 20), P(81, 4, 82)),
																T(S(P(78, 4, 24), P(79, 4, 80)), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(S(P(74, 4, 20), P(76, 4, 77)), "foo"),
																ast.NewIntLiteralNode(S(P(81, 4, 27), P(81, 4, 82)), "3"),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(88, 5, 6), P(95, 5, 96)),
															ast.NewAssignmentExpressionNode(
																S(P(88, 5, 6), P(95, 5, 96)),
																T(S(P(92, 5, 10), P(93, 5, 94)), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(S(P(88, 5, 6), P(90, 5, 91)), "foo"),
																ast.NewIntLiteralNode(S(P(95, 5, 13), P(95, 5, 96)), "2"),
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
						S(P(97, 6, 1), P(100, 6, 101)),
						ast.NewNilLiteralNode(S(P(97, 6, 1), P(99, 6, 100))),
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
				S(P(0, 1, 1), P(107, 1, 108)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(103, 2, 104)),
						ast.NewIfExpressionNode(
							S(P(1, 2, 1), P(102, 2, 103)),
							ast.NewBinaryExpressionNode(
								S(P(4, 2, 4), P(10, 2, 11)),
								T(S(P(8, 2, 8), P(8, 2, 9)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(4, 2, 4), P(6, 2, 7)), "foo"),
								ast.NewIntLiteralNode(S(P(10, 2, 10), P(10, 2, 11)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(13, 3, 2), P(21, 3, 22)),
									ast.NewAssignmentExpressionNode(
										S(P(13, 3, 2), P(20, 3, 21)),
										T(S(P(17, 3, 6), P(18, 3, 19)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(13, 3, 2), P(15, 3, 16)), "foo"),
										ast.NewIntLiteralNode(S(P(20, 3, 9), P(20, 3, 21)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(23, 4, 2), P(26, 4, 27)),
									ast.NewNilLiteralNode(S(P(23, 4, 2), P(25, 4, 26))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(32, 5, 6), P(102, 5, 103)),
									ast.NewIfExpressionNode(
										S(P(32, 5, 6), P(102, 5, 103)),
										ast.NewBinaryExpressionNode(
											S(P(35, 5, 9), P(41, 5, 42)),
											T(S(P(39, 5, 13), P(39, 5, 40)), token.LESS),
											ast.NewPublicIdentifierNode(S(P(35, 5, 9), P(37, 5, 38)), "foo"),
											ast.NewIntLiteralNode(S(P(41, 5, 15), P(41, 5, 42)), "5"),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(44, 6, 2), P(53, 6, 54)),
												ast.NewAssignmentExpressionNode(
													S(P(44, 6, 2), P(52, 6, 53)),
													T(S(P(48, 6, 6), P(49, 6, 50)), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(S(P(44, 6, 2), P(46, 6, 47)), "foo"),
													ast.NewIntLiteralNode(S(P(51, 6, 9), P(52, 6, 53)), "10"),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												S(P(59, 7, 6), P(102, 7, 103)),
												ast.NewIfExpressionNode(
													S(P(59, 7, 6), P(102, 7, 103)),
													ast.NewBinaryExpressionNode(
														S(P(62, 7, 9), P(68, 7, 69)),
														T(S(P(66, 7, 13), P(66, 7, 67)), token.LESS),
														ast.NewPublicIdentifierNode(S(P(62, 7, 9), P(64, 7, 65)), "foo"),
														ast.NewIntLiteralNode(S(P(68, 7, 15), P(68, 7, 69)), "0"),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(71, 8, 2), P(79, 8, 80)),
															ast.NewAssignmentExpressionNode(
																S(P(71, 8, 2), P(78, 8, 79)),
																T(S(P(75, 8, 6), P(76, 8, 77)), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(S(P(71, 8, 2), P(73, 8, 74)), "foo"),
																ast.NewIntLiteralNode(S(P(78, 8, 9), P(78, 8, 79)), "3"),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															S(P(86, 10, 2), P(94, 10, 95)),
															ast.NewAssignmentExpressionNode(
																S(P(86, 10, 2), P(93, 10, 94)),
																T(S(P(90, 10, 6), P(91, 10, 92)), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(S(P(86, 10, 2), P(88, 10, 89)), "foo"),
																ast.NewIntLiteralNode(S(P(93, 10, 9), P(93, 10, 94)), "2"),
															),
														),
														ast.NewExpressionStatementNode(
															S(P(96, 11, 2), P(99, 11, 100)),
															ast.NewNilLiteralNode(S(P(96, 11, 2), P(98, 11, 99))),
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
						S(P(104, 13, 1), P(107, 13, 108)),
						ast.NewNilLiteralNode(S(P(104, 13, 1), P(106, 13, 107))),
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
				S(P(0, 1, 1), P(34, 1, 35)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(34, 2, 35)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(33, 2, 34)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 15)),
								T(S(P(12, 2, 12), P(12, 2, 13)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 11)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 15)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 3, 2), P(25, 3, 26)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 3, 2), P(24, 3, 25)),
										T(S(P(21, 3, 6), P(22, 3, 23)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 3, 2), P(19, 3, 20)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 3, 9), P(24, 3, 25)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(27, 4, 2), P(30, 4, 31)),
									ast.NewNilLiteralNode(S(P(27, 4, 2), P(29, 4, 30))),
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
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(19, 2, 20)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(18, 2, 19)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 15)),
								T(S(P(12, 2, 12), P(12, 2, 13)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 11)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 15)), "0"),
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
				S(P(0, 1, 1), P(42, 1, 43)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(38, 2, 39)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(37, 2, 38)),
							T(S(P(5, 2, 5), P(5, 2, 6)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 4)), "bar"),
							ast.NewUnlessExpressionNode(
								S(P(8, 3, 2), P(37, 3, 38)),
								ast.NewBinaryExpressionNode(
									S(P(15, 3, 9), P(21, 3, 22)),
									T(S(P(19, 3, 13), P(19, 3, 20)), token.GREATER),
									ast.NewPublicIdentifierNode(S(P(15, 3, 9), P(17, 3, 18)), "foo"),
									ast.NewIntLiteralNode(S(P(21, 3, 15), P(21, 3, 22)), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(25, 4, 3), P(33, 4, 34)),
										ast.NewAssignmentExpressionNode(
											S(P(25, 4, 3), P(32, 4, 33)),
											T(S(P(29, 4, 7), P(30, 4, 31)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(25, 4, 3), P(27, 4, 28)), "foo"),
											ast.NewIntLiteralNode(S(P(32, 4, 10), P(32, 4, 33)), "2"),
										),
									),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(39, 6, 1), P(42, 6, 43)),
						ast.NewNilLiteralNode(S(P(39, 6, 1), P(41, 6, 42))),
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
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(29, 2, 30)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(28, 2, 29)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 15)),
								T(S(P(12, 2, 12), P(12, 2, 13)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 11)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 15)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 2, 21), P(28, 2, 29)),
									ast.NewAssignmentExpressionNode(
										S(P(21, 2, 21), P(28, 2, 29)),
										T(S(P(25, 2, 25), P(26, 2, 27)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(21, 2, 21), P(23, 2, 24)), "foo"),
										ast.NewIntLiteralNode(S(P(28, 2, 28), P(28, 2, 29)), "2"),
									),
								),
							},
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						S(P(30, 3, 1), P(33, 3, 34)),
						ast.NewNilLiteralNode(S(P(30, 3, 1), P(32, 3, 33))),
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
				S(P(0, 1, 1), P(58, 1, 59)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(54, 2, 55)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(53, 2, 54)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 15)),
								T(S(P(12, 2, 12), P(12, 2, 13)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 11)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 15)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 3, 2), P(25, 3, 26)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 3, 2), P(24, 3, 25)),
										T(S(P(21, 3, 6), P(22, 3, 23)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 3, 2), P(19, 3, 20)), "foo"),
										ast.NewIntLiteralNode(S(P(24, 3, 9), P(24, 3, 25)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(27, 4, 2), P(30, 4, 31)),
									ast.NewNilLiteralNode(S(P(27, 4, 2), P(29, 4, 30))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(37, 6, 2), P(45, 6, 46)),
									ast.NewAssignmentExpressionNode(
										S(P(37, 6, 2), P(44, 6, 45)),
										T(S(P(41, 6, 6), P(42, 6, 43)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(37, 6, 2), P(39, 6, 40)), "foo"),
										ast.NewIntLiteralNode(S(P(44, 6, 9), P(44, 6, 45)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(47, 7, 2), P(50, 7, 51)),
									ast.NewNilLiteralNode(S(P(47, 7, 2), P(49, 7, 50))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(55, 9, 1), P(58, 9, 59)),
						ast.NewNilLiteralNode(S(P(55, 9, 1), P(57, 9, 58))),
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
				S(P(0, 1, 1), P(47, 1, 48)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(43, 2, 44)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(42, 2, 43)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 15)),
								T(S(P(12, 2, 12), P(12, 2, 13)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 11)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 15)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 2, 21), P(28, 2, 29)),
									ast.NewAssignmentExpressionNode(
										S(P(21, 2, 21), P(28, 2, 29)),
										T(S(P(25, 2, 25), P(26, 2, 27)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(21, 2, 21), P(23, 2, 24)), "foo"),
										ast.NewIntLiteralNode(S(P(28, 2, 28), P(28, 2, 29)), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(35, 3, 6), P(42, 3, 43)),
									ast.NewAssignmentExpressionNode(
										S(P(35, 3, 6), P(42, 3, 43)),
										T(S(P(39, 3, 10), P(40, 3, 41)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(35, 3, 6), P(37, 3, 38)), "foo"),
										ast.NewIntLiteralNode(S(P(42, 3, 13), P(42, 3, 43)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(44, 4, 1), P(47, 4, 48)),
						ast.NewNilLiteralNode(S(P(44, 4, 1), P(46, 4, 47))),
					),
				},
			),
		},
		"can't have two elses": {
			input: `
unless foo > 0 then foo += 2
else foo -= 2
else bar
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(56, 1, 57)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(43, 2, 44)),
						ast.NewUnlessExpressionNode(
							S(P(1, 2, 1), P(42, 2, 43)),
							ast.NewBinaryExpressionNode(
								S(P(8, 2, 8), P(14, 2, 15)),
								T(S(P(12, 2, 12), P(12, 2, 13)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(8, 2, 8), P(10, 2, 11)), "foo"),
								ast.NewIntLiteralNode(S(P(14, 2, 14), P(14, 2, 15)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 2, 21), P(28, 2, 29)),
									ast.NewAssignmentExpressionNode(
										S(P(21, 2, 21), P(28, 2, 29)),
										T(S(P(25, 2, 25), P(26, 2, 27)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(21, 2, 21), P(23, 2, 24)), "foo"),
										ast.NewIntLiteralNode(S(P(28, 2, 28), P(28, 2, 29)), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(35, 3, 6), P(42, 3, 43)),
									ast.NewAssignmentExpressionNode(
										S(P(35, 3, 6), P(42, 3, 43)),
										T(S(P(39, 3, 10), P(40, 3, 41)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(35, 3, 6), P(37, 3, 38)), "foo"),
										ast.NewIntLiteralNode(S(P(42, 3, 13), P(42, 3, 43)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(44, 4, 1), P(52, 4, 53)),
						ast.NewInvalidNode(S(P(44, 4, 1), P(47, 4, 48)), T(S(P(44, 4, 1), P(47, 4, 48)), token.ELSE)),
					),
					ast.NewExpressionStatementNode(
						S(P(53, 5, 1), P(56, 5, 57)),
						ast.NewNilLiteralNode(S(P(53, 5, 1), P(55, 5, 56))),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(44, 4, 1), P(47, 4, 48)), "unexpected else, expected an expression"),
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
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(33, 2, 34)),
						ast.NewWhileExpressionNode(
							S(P(1, 2, 1), P(32, 2, 33)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 14)),
								T(S(P(11, 2, 11), P(11, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 14)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(16, 3, 2), P(24, 3, 25)),
									ast.NewAssignmentExpressionNode(
										S(P(16, 3, 2), P(23, 3, 24)),
										T(S(P(20, 3, 6), P(21, 3, 22)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(16, 3, 2), P(18, 3, 19)), "foo"),
										ast.NewIntLiteralNode(S(P(23, 3, 9), P(23, 3, 24)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(26, 4, 2), P(29, 4, 30)),
									ast.NewNilLiteralNode(S(P(26, 4, 2), P(28, 4, 29))),
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
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(18, 2, 19)),
						ast.NewWhileExpressionNode(
							S(P(1, 2, 1), P(17, 2, 18)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 14)),
								T(S(P(11, 2, 11), P(11, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 14)), "0"),
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
				S(P(0, 1, 1), P(41, 1, 42)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(37, 2, 38)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(36, 2, 37)),
							T(S(P(5, 2, 5), P(5, 2, 6)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 4)), "bar"),
							ast.NewWhileExpressionNode(
								S(P(8, 3, 2), P(36, 3, 37)),
								ast.NewBinaryExpressionNode(
									S(P(14, 3, 8), P(20, 3, 21)),
									T(S(P(18, 3, 12), P(18, 3, 19)), token.GREATER),
									ast.NewPublicIdentifierNode(S(P(14, 3, 8), P(16, 3, 17)), "foo"),
									ast.NewIntLiteralNode(S(P(20, 3, 14), P(20, 3, 21)), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(24, 4, 3), P(32, 4, 33)),
										ast.NewAssignmentExpressionNode(
											S(P(24, 4, 3), P(31, 4, 32)),
											T(S(P(28, 4, 7), P(29, 4, 30)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(24, 4, 3), P(26, 4, 27)), "foo"),
											ast.NewIntLiteralNode(S(P(31, 4, 10), P(31, 4, 32)), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(38, 6, 1), P(41, 6, 42)),
						ast.NewNilLiteralNode(S(P(38, 6, 1), P(40, 6, 41))),
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
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(28, 2, 29)),
						ast.NewWhileExpressionNode(
							S(P(1, 2, 1), P(27, 2, 28)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 14)),
								T(S(P(11, 2, 11), P(11, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 14)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(20, 2, 20), P(27, 2, 28)),
									ast.NewAssignmentExpressionNode(
										S(P(20, 2, 20), P(27, 2, 28)),
										T(S(P(24, 2, 24), P(25, 2, 26)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(20, 2, 20), P(22, 2, 23)), "foo"),
										ast.NewIntLiteralNode(S(P(27, 2, 27), P(27, 2, 28)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(29, 3, 1), P(32, 3, 33)),
						ast.NewNilLiteralNode(S(P(29, 3, 1), P(31, 3, 32))),
					),
				},
			),
		},
		"can't have else": {
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
				S(P(0, 1, 1), P(57, 1, 58)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(53, 2, 54)),
						ast.NewWhileExpressionNode(
							S(P(1, 2, 1), P(52, 2, 53)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 14)),
								T(S(P(11, 2, 11), P(11, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 14)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(16, 3, 2), P(24, 3, 25)),
									ast.NewAssignmentExpressionNode(
										S(P(16, 3, 2), P(23, 3, 24)),
										T(S(P(20, 3, 6), P(21, 3, 22)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(16, 3, 2), P(18, 3, 19)), "foo"),
										ast.NewIntLiteralNode(S(P(23, 3, 9), P(23, 3, 24)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(26, 4, 2), P(29, 4, 30)),
									ast.NewNilLiteralNode(S(P(26, 4, 2), P(28, 4, 29))),
								),
								ast.NewExpressionStatementNode(
									S(P(30, 5, 1), P(34, 5, 35)),
									ast.NewInvalidNode(S(P(30, 5, 1), P(33, 5, 34)), T(S(P(30, 5, 1), P(33, 5, 34)), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									S(P(36, 6, 2), P(44, 6, 45)),
									ast.NewAssignmentExpressionNode(
										S(P(36, 6, 2), P(43, 6, 44)),
										T(S(P(40, 6, 6), P(41, 6, 42)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(36, 6, 2), P(38, 6, 39)), "foo"),
										ast.NewIntLiteralNode(S(P(43, 6, 9), P(43, 6, 44)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(46, 7, 2), P(49, 7, 50)),
									ast.NewNilLiteralNode(S(P(46, 7, 2), P(48, 7, 49))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(54, 9, 1), P(57, 9, 58)),
						ast.NewNilLiteralNode(S(P(54, 9, 1), P(56, 9, 57))),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(30, 5, 1), P(33, 5, 34)), "unexpected else, expected an expression"),
			},
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
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(33, 2, 34)),
						ast.NewUntilExpressionNode(
							S(P(1, 2, 1), P(32, 2, 33)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 14)),
								T(S(P(11, 2, 11), P(11, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 14)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(16, 3, 2), P(24, 3, 25)),
									ast.NewAssignmentExpressionNode(
										S(P(16, 3, 2), P(23, 3, 24)),
										T(S(P(20, 3, 6), P(21, 3, 22)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(16, 3, 2), P(18, 3, 19)), "foo"),
										ast.NewIntLiteralNode(S(P(23, 3, 9), P(23, 3, 24)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(26, 4, 2), P(29, 4, 30)),
									ast.NewNilLiteralNode(S(P(26, 4, 2), P(28, 4, 29))),
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
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(18, 2, 19)),
						ast.NewUntilExpressionNode(
							S(P(1, 2, 1), P(17, 2, 18)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 14)),
								T(S(P(11, 2, 11), P(11, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 14)), "0"),
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
				S(P(0, 1, 1), P(41, 1, 42)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(37, 2, 38)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(36, 2, 37)),
							T(S(P(5, 2, 5), P(5, 2, 6)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 4)), "bar"),
							ast.NewUntilExpressionNode(
								S(P(8, 3, 2), P(36, 3, 37)),
								ast.NewBinaryExpressionNode(
									S(P(14, 3, 8), P(20, 3, 21)),
									T(S(P(18, 3, 12), P(18, 3, 19)), token.GREATER),
									ast.NewPublicIdentifierNode(S(P(14, 3, 8), P(16, 3, 17)), "foo"),
									ast.NewIntLiteralNode(S(P(20, 3, 14), P(20, 3, 21)), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(24, 4, 3), P(32, 4, 33)),
										ast.NewAssignmentExpressionNode(
											S(P(24, 4, 3), P(31, 4, 32)),
											T(S(P(28, 4, 7), P(29, 4, 30)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(24, 4, 3), P(26, 4, 27)), "foo"),
											ast.NewIntLiteralNode(S(P(31, 4, 10), P(31, 4, 32)), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(38, 6, 1), P(41, 6, 42)),
						ast.NewNilLiteralNode(S(P(38, 6, 1), P(40, 6, 41))),
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
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(28, 2, 29)),
						ast.NewUntilExpressionNode(
							S(P(1, 2, 1), P(27, 2, 28)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 14)),
								T(S(P(11, 2, 11), P(11, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 14)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(20, 2, 20), P(27, 2, 28)),
									ast.NewAssignmentExpressionNode(
										S(P(20, 2, 20), P(27, 2, 28)),
										T(S(P(24, 2, 24), P(25, 2, 26)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(20, 2, 20), P(22, 2, 23)), "foo"),
										ast.NewIntLiteralNode(S(P(27, 2, 27), P(27, 2, 28)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(29, 3, 1), P(32, 3, 33)),
						ast.NewNilLiteralNode(S(P(29, 3, 1), P(31, 3, 32))),
					),
				},
			),
		},
		"can't have else": {
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
				S(P(0, 1, 1), P(57, 1, 58)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(53, 2, 54)),
						ast.NewUntilExpressionNode(
							S(P(1, 2, 1), P(52, 2, 53)),
							ast.NewBinaryExpressionNode(
								S(P(7, 2, 7), P(13, 2, 14)),
								T(S(P(11, 2, 11), P(11, 2, 12)), token.GREATER),
								ast.NewPublicIdentifierNode(S(P(7, 2, 7), P(9, 2, 10)), "foo"),
								ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 14)), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(16, 3, 2), P(24, 3, 25)),
									ast.NewAssignmentExpressionNode(
										S(P(16, 3, 2), P(23, 3, 24)),
										T(S(P(20, 3, 6), P(21, 3, 22)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(16, 3, 2), P(18, 3, 19)), "foo"),
										ast.NewIntLiteralNode(S(P(23, 3, 9), P(23, 3, 24)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(26, 4, 2), P(29, 4, 30)),
									ast.NewNilLiteralNode(S(P(26, 4, 2), P(28, 4, 29))),
								),
								ast.NewExpressionStatementNode(
									S(P(30, 5, 1), P(34, 5, 35)),
									ast.NewInvalidNode(S(P(30, 5, 1), P(33, 5, 34)), T(S(P(30, 5, 1), P(33, 5, 34)), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									S(P(36, 6, 2), P(44, 6, 45)),
									ast.NewAssignmentExpressionNode(
										S(P(36, 6, 2), P(43, 6, 44)),
										T(S(P(40, 6, 6), P(41, 6, 42)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(36, 6, 2), P(38, 6, 39)), "foo"),
										ast.NewIntLiteralNode(S(P(43, 6, 9), P(43, 6, 44)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(46, 7, 2), P(49, 7, 50)),
									ast.NewNilLiteralNode(S(P(46, 7, 2), P(48, 7, 49))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(54, 9, 1), P(57, 9, 58)),
						ast.NewNilLiteralNode(S(P(54, 9, 1), P(56, 9, 57))),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(30, 5, 1), P(33, 5, 34)), "unexpected else, expected an expression"),
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
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(24, 2, 25)),
						ast.NewLoopExpressionNode(
							S(P(1, 2, 1), P(23, 2, 24)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(7, 3, 2), P(15, 3, 16)),
									ast.NewAssignmentExpressionNode(
										S(P(7, 3, 2), P(14, 3, 15)),
										T(S(P(11, 3, 6), P(12, 3, 13)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(7, 3, 2), P(9, 3, 10)), "foo"),
										ast.NewIntLiteralNode(S(P(14, 3, 9), P(14, 3, 15)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(17, 4, 2), P(20, 4, 21)),
									ast.NewNilLiteralNode(S(P(17, 4, 2), P(19, 4, 20))),
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
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(9, 2, 10)),
						ast.NewLoopExpressionNode(
							S(P(1, 2, 1), P(8, 2, 9)),
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
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(28, 2, 29)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(27, 2, 28)),
							T(S(P(5, 2, 5), P(5, 2, 6)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 4)), "bar"),
							ast.NewLoopExpressionNode(
								S(P(8, 3, 2), P(27, 3, 28)),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(15, 4, 3), P(23, 4, 24)),
										ast.NewAssignmentExpressionNode(
											S(P(15, 4, 3), P(22, 4, 23)),
											T(S(P(19, 4, 7), P(20, 4, 21)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(15, 4, 3), P(17, 4, 18)), "foo"),
											ast.NewIntLiteralNode(S(P(22, 4, 10), P(22, 4, 23)), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(29, 6, 1), P(32, 6, 33)),
						ast.NewNilLiteralNode(S(P(29, 6, 1), P(31, 6, 32))),
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
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(14, 2, 15)),
						ast.NewLoopExpressionNode(
							S(P(1, 2, 1), P(13, 2, 14)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(6, 2, 6), P(13, 2, 14)),
									ast.NewAssignmentExpressionNode(
										S(P(6, 2, 6), P(13, 2, 14)),
										T(S(P(10, 2, 10), P(11, 2, 12)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(6, 2, 6), P(8, 2, 9)), "foo"),
										ast.NewIntLiteralNode(S(P(13, 2, 13), P(13, 2, 14)), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(15, 3, 1), P(18, 3, 19)),
						ast.NewNilLiteralNode(S(P(15, 3, 1), P(17, 3, 18))),
					),
				},
			),
		},
		"can't have else": {
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
				S(P(0, 1, 1), P(48, 1, 49)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(44, 2, 45)),
						ast.NewLoopExpressionNode(
							S(P(1, 2, 1), P(43, 2, 44)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(7, 3, 2), P(15, 3, 16)),
									ast.NewAssignmentExpressionNode(
										S(P(7, 3, 2), P(14, 3, 15)),
										T(S(P(11, 3, 6), P(12, 3, 13)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(7, 3, 2), P(9, 3, 10)), "foo"),
										ast.NewIntLiteralNode(S(P(14, 3, 9), P(14, 3, 15)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(17, 4, 2), P(20, 4, 21)),
									ast.NewNilLiteralNode(S(P(17, 4, 2), P(19, 4, 20))),
								),
								ast.NewExpressionStatementNode(
									S(P(21, 5, 1), P(25, 5, 26)),
									ast.NewInvalidNode(S(P(21, 5, 1), P(24, 5, 25)), T(S(P(21, 5, 1), P(24, 5, 25)), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									S(P(27, 6, 2), P(35, 6, 36)),
									ast.NewAssignmentExpressionNode(
										S(P(27, 6, 2), P(34, 6, 35)),
										T(S(P(31, 6, 6), P(32, 6, 33)), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(27, 6, 2), P(29, 6, 30)), "foo"),
										ast.NewIntLiteralNode(S(P(34, 6, 9), P(34, 6, 35)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(37, 7, 2), P(40, 7, 41)),
									ast.NewNilLiteralNode(S(P(37, 7, 2), P(39, 7, 40))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(45, 9, 1), P(48, 9, 49)),
						ast.NewNilLiteralNode(S(P(45, 9, 1), P(47, 9, 48))),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(21, 5, 1), P(24, 5, 25)), "unexpected else, expected an expression"),
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
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(22, 2, 23)),
						ast.NewDoExpressionNode(
							S(P(1, 2, 1), P(21, 2, 22)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(5, 3, 2), P(13, 3, 14)),
									ast.NewAssignmentExpressionNode(
										S(P(5, 3, 2), P(12, 3, 13)),
										T(S(P(9, 3, 6), P(10, 3, 11)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(5, 3, 2), P(7, 3, 8)), "foo"),
										ast.NewIntLiteralNode(S(P(12, 3, 9), P(12, 3, 13)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(15, 4, 2), P(18, 4, 19)),
									ast.NewNilLiteralNode(S(P(15, 4, 2), P(17, 4, 18))),
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
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(7, 2, 8)),
						ast.NewDoExpressionNode(
							S(P(1, 2, 1), P(6, 2, 7)),
							nil,
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
				S(P(0, 1, 1), P(30, 1, 31)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(26, 2, 27)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(25, 2, 26)),
							T(S(P(5, 2, 5), P(5, 2, 6)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 4)), "bar"),
							ast.NewDoExpressionNode(
								S(P(8, 3, 2), P(25, 3, 26)),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(13, 4, 3), P(21, 4, 22)),
										ast.NewAssignmentExpressionNode(
											S(P(13, 4, 3), P(20, 4, 21)),
											T(S(P(17, 4, 7), P(18, 4, 19)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(13, 4, 3), P(15, 4, 16)), "foo"),
											ast.NewIntLiteralNode(S(P(20, 4, 10), P(20, 4, 21)), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(27, 6, 1), P(30, 6, 31)),
						ast.NewNilLiteralNode(S(P(27, 6, 1), P(29, 6, 30))),
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
						ast.NewBreakExpressionNode(S(P(0, 1, 1), P(4, 1, 5))),
					),
				},
			),
		},
		"can't have an argument": {
			input: `break 2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewBreakExpressionNode(S(P(0, 1, 1), P(4, 1, 5))),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(
					L("main", P(6, 1, 7), P(6, 1, 7)),
					"unexpected INT, expected a statement separator `\\n`, `;`",
				),
			},
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
							ast.NewBreakExpressionNode(S(P(7, 1, 8), P(11, 1, 12))),
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
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewReturnExpressionNode(S(P(0, 1, 1), P(5, 1, 6)), nil),
					),
					ast.NewExpressionStatementNode(
						S(P(7, 2, 1), P(7, 2, 8)),
						ast.NewIntLiteralNode(S(P(7, 2, 1), P(7, 2, 8)), "1"),
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
						ast.NewContinueExpressionNode(S(P(0, 1, 1), P(7, 1, 8)), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "continue\n1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewContinueExpressionNode(S(P(0, 1, 1), P(7, 1, 8)), nil),
					),
					ast.NewExpressionStatementNode(
						S(P(9, 2, 1), P(9, 2, 10)),
						ast.NewIntLiteralNode(S(P(9, 2, 1), P(9, 2, 10)), "1"),
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
							ast.NewIntLiteralNode(S(P(9, 1, 10), P(9, 1, 10)), "2"),
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
							ast.NewContinueExpressionNode(S(P(7, 1, 8), P(14, 1, 15)), nil),
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
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewThrowExpressionNode(S(P(0, 1, 1), P(4, 1, 5)), nil),
					),
					ast.NewExpressionStatementNode(
						S(P(6, 2, 1), P(6, 2, 7)),
						ast.NewIntLiteralNode(S(P(6, 2, 1), P(6, 2, 7)), "1"),
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

func TestFor(t *testing.T) {
	tests := testTable{
		"can be single-line with then": {
			input: `for i in [1, 2, 3] then println(i)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewForExpressionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							[]ast.ParameterNode{
								ast.NewLoopParameterNode(S(P(4, 1, 5), P(4, 1, 5)), "i", nil),
							},
							ast.NewListLiteralNode(
								S(P(9, 1, 10), P(17, 1, 18)),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(S(P(10, 1, 11), P(10, 1, 11)), "1"),
									ast.NewIntLiteralNode(S(P(13, 1, 14), P(13, 1, 14)), "2"),
									ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "3"),
								},
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
				S(P(0, 1, 1), P(40, 1, 41)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 1, 41)),
						ast.NewForExpressionNode(
							S(P(0, 1, 1), P(40, 1, 41)),
							[]ast.ParameterNode{
								ast.NewLoopParameterNode(S(P(4, 1, 5), P(4, 1, 5)), "i", nil),
							},
							ast.NewListLiteralNode(
								S(P(9, 1, 10), P(17, 1, 18)),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(S(P(10, 1, 11), P(10, 1, 11)), "1"),
									ast.NewIntLiteralNode(S(P(13, 1, 14), P(13, 1, 14)), "2"),
									ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "3"),
								},
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 2, 3), P(31, 2, 32)),
									ast.NewFunctionCallNode(
										S(P(21, 2, 3), P(30, 2, 31)),
										"println",
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(S(P(29, 2, 11), P(29, 2, 30)), "i"),
										},
										nil,
									),
								),
								ast.NewExpressionStatementNode(
									S(P(34, 3, 3), P(37, 3, 38)),
									ast.NewNilLiteralNode(S(P(34, 3, 3), P(36, 3, 37))),
								),
							},
						),
					),
				},
			),
		},
		"can have multiple parameters": {
			input: `for i, j: Int in [1, 2, 3] then println(i)`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(41, 1, 42)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(41, 1, 42)),
						ast.NewForExpressionNode(
							S(P(0, 1, 1), P(41, 1, 42)),
							[]ast.ParameterNode{
								ast.NewLoopParameterNode(S(P(4, 1, 5), P(4, 1, 5)), "i", nil),
								ast.NewLoopParameterNode(S(P(7, 1, 8), P(12, 1, 13)), "j", ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Int")),
							},
							ast.NewListLiteralNode(
								S(P(17, 1, 18), P(25, 1, 26)),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(S(P(18, 1, 19), P(18, 1, 19)), "1"),
									ast.NewIntLiteralNode(S(P(21, 1, 22), P(21, 1, 22)), "2"),
									ast.NewIntLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "3"),
								},
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(32, 1, 33), P(41, 1, 42)),
									ast.NewFunctionCallNode(
										S(P(32, 1, 33), P(41, 1, 42)),
										"println",
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(S(P(40, 1, 41), P(40, 1, 41)), "i"),
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}
