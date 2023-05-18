package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
)

func TestModifierExpression(t *testing.T) {
	tests := testTable{
		"has lower precedence than assignment": {
			input: "foo = bar if baz",
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewModifierNode(
							P(0, 16, 1, 1),
							T(P(10, 2, 1, 11), token.IF),
							ast.NewAssignmentExpressionNode(
								P(0, 9, 1, 1),
								T(P(4, 1, 1, 5), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
							),
							ast.NewPublicIdentifierNode(P(13, 3, 1, 14), "baz"),
						),
					),
				},
			),
		},
		"if can contain else": {
			input: "foo = bar if baz else car = red",
			want: ast.NewProgramNode(
				P(0, 31, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 31, 1, 1),
						ast.NewModifierIfElseNode(
							P(0, 31, 1, 1),
							ast.NewAssignmentExpressionNode(
								P(0, 9, 1, 1),
								T(P(4, 1, 1, 5), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
							),
							ast.NewPublicIdentifierNode(P(13, 3, 1, 14), "baz"),
							ast.NewAssignmentExpressionNode(
								P(22, 9, 1, 23),
								T(P(26, 1, 1, 27), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(P(22, 3, 1, 23), "car"),
								ast.NewPublicIdentifierNode(P(28, 3, 1, 29), "red"),
							),
						),
					),
				},
			),
		},
		"has many versions": {
			input: "foo if bar\nfoo unless bar\nfoo while bar\nfoo until bar",
			want: ast.NewProgramNode(
				P(0, 53, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewModifierNode(
							P(0, 10, 1, 1),
							T(P(4, 2, 1, 5), token.IF),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						P(11, 15, 2, 1),
						ast.NewModifierNode(
							P(11, 14, 2, 1),
							T(P(15, 6, 2, 5), token.UNLESS),
							ast.NewPublicIdentifierNode(P(11, 3, 2, 1), "foo"),
							ast.NewPublicIdentifierNode(P(22, 3, 2, 12), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						P(26, 14, 3, 1),
						ast.NewModifierNode(
							P(26, 13, 3, 1),
							T(P(30, 5, 3, 5), token.WHILE),
							ast.NewPublicIdentifierNode(P(26, 3, 3, 1), "foo"),
							ast.NewPublicIdentifierNode(P(36, 3, 3, 11), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						P(40, 13, 4, 1),
						ast.NewModifierNode(
							P(40, 13, 4, 1),
							T(P(44, 5, 4, 5), token.UNTIL),
							ast.NewPublicIdentifierNode(P(40, 3, 4, 1), "foo"),
							ast.NewPublicIdentifierNode(P(50, 3, 4, 11), "bar"),
						),
					),
				},
			),
		},
		"can't be nested": {
			input: "foo = bar if baz if false\n3",
			want: ast.NewProgramNode(
				P(0, 27, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewModifierNode(
							P(0, 16, 1, 1),
							T(P(10, 2, 1, 11), token.IF),
							ast.NewAssignmentExpressionNode(
								P(0, 9, 1, 1),
								T(P(4, 1, 1, 5), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
							),
							ast.NewPublicIdentifierNode(P(13, 3, 1, 14), "baz"),
						),
					),
					ast.NewExpressionStatementNode(
						P(26, 1, 2, 1),
						ast.NewIntLiteralNode(P(26, 1, 2, 1), V(P(26, 1, 2, 1), token.DEC_INT, "3")),
					),
				},
			),
			err: ErrorList{
				NewError(P(17, 2, 1, 18), "unexpected if, expected a statement separator `\\n`, `;`"),
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
				P(0, 31, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 30, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 29, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), token.GREATER),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(13, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(13, 8, 3, 2),
										T(P(17, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(13, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(20, 1, 3, 9), V(P(20, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(23, 4, 4, 2),
									ast.NewNilLiteralNode(P(23, 3, 4, 2)),
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
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 15, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 14, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), token.GREATER),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), token.DEC_INT, "0")),
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
				P(0, 39, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 34, 2, 1),
						ast.NewAssignmentExpressionNode(
							P(1, 33, 2, 1),
							T(P(5, 1, 2, 5), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(1, 3, 2, 1), "bar"),
							ast.NewIfExpressionNode(
								P(8, 26, 3, 2),
								ast.NewBinaryExpressionNode(
									P(11, 7, 3, 5),
									T(P(15, 1, 3, 9), token.GREATER),
									ast.NewPublicIdentifierNode(P(11, 3, 3, 5), "foo"),
									ast.NewIntLiteralNode(P(17, 1, 3, 11), V(P(17, 1, 3, 11), token.DEC_INT, "0")),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										P(21, 9, 4, 3),
										ast.NewAssignmentExpressionNode(
											P(21, 8, 4, 3),
											T(P(25, 2, 4, 7), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(P(21, 3, 4, 3), "foo"),
											ast.NewIntLiteralNode(P(28, 1, 4, 10), V(P(28, 1, 4, 10), token.DEC_INT, "2")),
										),
									),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						P(35, 4, 6, 1),
						ast.NewNilLiteralNode(P(35, 3, 6, 1)),
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
				P(0, 30, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 25, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 24, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), token.GREATER),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 8, 2, 17),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 2, 17),
										T(P(21, 2, 2, 21), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(17, 3, 2, 17), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 2, 24), V(P(24, 1, 2, 24), token.DEC_INT, "2")),
									),
								),
							},
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						P(26, 4, 3, 1),
						ast.NewNilLiteralNode(P(26, 3, 3, 1)),
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
				P(0, 56, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 51, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 50, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), token.GREATER),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(13, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(13, 8, 3, 2),
										T(P(17, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(13, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(20, 1, 3, 9), V(P(20, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(23, 4, 4, 2),
									ast.NewNilLiteralNode(P(23, 3, 4, 2)),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(34, 9, 6, 3),
									ast.NewAssignmentExpressionNode(
										P(34, 8, 6, 3),
										T(P(38, 2, 6, 7), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(P(34, 3, 6, 3), "foo"),
										ast.NewIntLiteralNode(P(41, 1, 6, 10), V(P(41, 1, 6, 10), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(44, 4, 7, 2),
									ast.NewNilLiteralNode(P(44, 3, 7, 2)),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(52, 4, 9, 1),
						ast.NewNilLiteralNode(P(52, 3, 9, 1)),
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
				P(0, 44, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 39, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 38, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), token.GREATER),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 8, 2, 17),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 2, 17),
										T(P(21, 2, 2, 21), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(17, 3, 2, 17), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 2, 24), V(P(24, 1, 2, 24), token.DEC_INT, "2")),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(31, 8, 3, 6),
									ast.NewAssignmentExpressionNode(
										P(31, 8, 3, 6),
										T(P(35, 2, 3, 10), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(P(31, 3, 3, 6), "foo"),
										ast.NewIntLiteralNode(P(38, 1, 3, 13), V(P(38, 1, 3, 13), token.DEC_INT, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(40, 4, 4, 1),
						ast.NewNilLiteralNode(P(40, 3, 4, 1)),
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
				P(0, 53, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 39, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 38, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), token.GREATER),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 8, 2, 17),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 2, 17),
										T(P(21, 2, 2, 21), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(17, 3, 2, 17), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 2, 24), V(P(24, 1, 2, 24), token.DEC_INT, "2")),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(31, 8, 3, 6),
									ast.NewAssignmentExpressionNode(
										P(31, 8, 3, 6),
										T(P(35, 2, 3, 10), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(P(31, 3, 3, 6), "foo"),
										ast.NewIntLiteralNode(P(38, 1, 3, 13), V(P(38, 1, 3, 13), token.DEC_INT, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(40, 9, 4, 1),
						ast.NewInvalidNode(P(40, 4, 4, 1), T(P(40, 4, 4, 1), token.ELSE)),
					),
					ast.NewExpressionStatementNode(
						P(49, 4, 5, 1),
						ast.NewNilLiteralNode(P(49, 3, 5, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(P(40, 4, 4, 1), "unexpected else, expected an expression"),
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
				P(0, 104, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 99, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 98, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), token.GREATER),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(13, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(13, 8, 3, 2),
										T(P(17, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(13, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(20, 1, 3, 9), V(P(20, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(23, 4, 4, 2),
									ast.NewNilLiteralNode(P(23, 3, 4, 2)),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(27, 25, 5, 1),
									ast.NewIfExpressionNode(
										P(27, 25, 5, 1),
										ast.NewBinaryExpressionNode(
											P(33, 7, 5, 7),
											T(P(37, 1, 5, 11), token.LESS),
											ast.NewPublicIdentifierNode(P(33, 3, 5, 7), "foo"),
											ast.NewIntLiteralNode(P(39, 1, 5, 13), V(P(39, 1, 5, 13), token.DEC_INT, "5")),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(42, 10, 6, 2),
												ast.NewAssignmentExpressionNode(
													P(42, 9, 6, 2),
													T(P(46, 2, 6, 6), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(P(42, 3, 6, 2), "foo"),
													ast.NewIntLiteralNode(P(49, 2, 6, 9), V(P(49, 2, 6, 9), token.DEC_INT, "10")),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(52, 47, 7, 1),
												ast.NewIfExpressionNode(
													P(52, 47, 7, 1),
													ast.NewBinaryExpressionNode(
														P(58, 7, 7, 7),
														T(P(62, 1, 7, 11), token.LESS),
														ast.NewPublicIdentifierNode(P(58, 3, 7, 7), "foo"),
														ast.NewIntLiteralNode(P(64, 1, 7, 13), V(P(64, 1, 7, 13), token.DEC_INT, "0")),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(67, 9, 8, 2),
															ast.NewAssignmentExpressionNode(
																P(67, 8, 8, 2),
																T(P(71, 2, 8, 6), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(P(67, 3, 8, 2), "foo"),
																ast.NewIntLiteralNode(P(74, 1, 8, 9), V(P(74, 1, 8, 9), token.DEC_INT, "3")),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(82, 9, 10, 2),
															ast.NewAssignmentExpressionNode(
																P(82, 8, 10, 2),
																T(P(86, 2, 10, 6), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(P(82, 3, 10, 2), "foo"),
																ast.NewIntLiteralNode(P(89, 1, 10, 9), V(P(89, 1, 10, 9), token.DEC_INT, "2")),
															),
														),
														ast.NewExpressionStatementNode(
															P(92, 4, 11, 2),
															ast.NewNilLiteralNode(P(92, 3, 11, 2)),
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
						P(100, 4, 13, 1),
						ast.NewNilLiteralNode(P(100, 3, 13, 1)),
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
				P(0, 101, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 96, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 95, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), token.GREATER),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 8, 2, 17),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 2, 17),
										T(P(21, 2, 2, 21), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(17, 3, 2, 17), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 2, 24), V(P(24, 1, 2, 24), token.DEC_INT, "2")),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(26, 28, 3, 1),
									ast.NewIfExpressionNode(
										P(26, 28, 3, 1),
										ast.NewBinaryExpressionNode(
											P(32, 7, 3, 7),
											T(P(36, 1, 3, 11), token.LESS),
											ast.NewPublicIdentifierNode(P(32, 3, 3, 7), "foo"),
											ast.NewIntLiteralNode(P(38, 1, 3, 13), V(P(38, 1, 3, 13), token.DEC_INT, "5")),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(45, 9, 3, 20),
												ast.NewAssignmentExpressionNode(
													P(45, 9, 3, 20),
													T(P(49, 2, 3, 24), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(P(45, 3, 3, 20), "foo"),
													ast.NewIntLiteralNode(P(52, 2, 3, 27), V(P(52, 2, 3, 27), token.DEC_INT, "10")),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(55, 41, 4, 1),
												ast.NewIfExpressionNode(
													P(55, 41, 4, 1),
													ast.NewBinaryExpressionNode(
														P(61, 7, 4, 7),
														T(P(65, 1, 4, 11), token.LESS),
														ast.NewPublicIdentifierNode(P(61, 3, 4, 7), "foo"),
														ast.NewIntLiteralNode(P(67, 1, 4, 13), V(P(67, 1, 4, 13), token.DEC_INT, "0")),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(74, 8, 4, 20),
															ast.NewAssignmentExpressionNode(
																P(74, 8, 4, 20),
																T(P(78, 2, 4, 24), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(P(74, 3, 4, 20), "foo"),
																ast.NewIntLiteralNode(P(81, 1, 4, 27), V(P(81, 1, 4, 27), token.DEC_INT, "3")),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(88, 8, 5, 6),
															ast.NewAssignmentExpressionNode(
																P(88, 8, 5, 6),
																T(P(92, 2, 5, 10), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(P(88, 3, 5, 6), "foo"),
																ast.NewIntLiteralNode(P(95, 1, 5, 13), V(P(95, 1, 5, 13), token.DEC_INT, "2")),
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
						P(97, 4, 6, 1),
						ast.NewNilLiteralNode(P(97, 3, 6, 1)),
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
				P(0, 108, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 103, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 102, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), token.GREATER),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(13, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(13, 8, 3, 2),
										T(P(17, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(13, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(20, 1, 3, 9), V(P(20, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(23, 4, 4, 2),
									ast.NewNilLiteralNode(P(23, 3, 4, 2)),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(32, 71, 5, 6),
									ast.NewIfExpressionNode(
										P(32, 71, 5, 6),
										ast.NewBinaryExpressionNode(
											P(35, 7, 5, 9),
											T(P(39, 1, 5, 13), token.LESS),
											ast.NewPublicIdentifierNode(P(35, 3, 5, 9), "foo"),
											ast.NewIntLiteralNode(P(41, 1, 5, 15), V(P(41, 1, 5, 15), token.DEC_INT, "5")),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(44, 10, 6, 2),
												ast.NewAssignmentExpressionNode(
													P(44, 9, 6, 2),
													T(P(48, 2, 6, 6), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(P(44, 3, 6, 2), "foo"),
													ast.NewIntLiteralNode(P(51, 2, 6, 9), V(P(51, 2, 6, 9), token.DEC_INT, "10")),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(59, 44, 7, 6),
												ast.NewIfExpressionNode(
													P(59, 44, 7, 6),
													ast.NewBinaryExpressionNode(
														P(62, 7, 7, 9),
														T(P(66, 1, 7, 13), token.LESS),
														ast.NewPublicIdentifierNode(P(62, 3, 7, 9), "foo"),
														ast.NewIntLiteralNode(P(68, 1, 7, 15), V(P(68, 1, 7, 15), token.DEC_INT, "0")),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(71, 9, 8, 2),
															ast.NewAssignmentExpressionNode(
																P(71, 8, 8, 2),
																T(P(75, 2, 8, 6), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(P(71, 3, 8, 2), "foo"),
																ast.NewIntLiteralNode(P(78, 1, 8, 9), V(P(78, 1, 8, 9), token.DEC_INT, "3")),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(86, 9, 10, 2),
															ast.NewAssignmentExpressionNode(
																P(86, 8, 10, 2),
																T(P(90, 2, 10, 6), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(P(86, 3, 10, 2), "foo"),
																ast.NewIntLiteralNode(P(93, 1, 10, 9), V(P(93, 1, 10, 9), token.DEC_INT, "2")),
															),
														),
														ast.NewExpressionStatementNode(
															P(96, 4, 11, 2),
															ast.NewNilLiteralNode(P(96, 3, 11, 2)),
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
						P(104, 4, 13, 1),
						ast.NewNilLiteralNode(P(104, 3, 13, 1)),
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
				P(0, 35, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 34, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 33, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), token.GREATER),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 3, 2),
										T(P(21, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(17, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 3, 9), V(P(24, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(27, 4, 4, 2),
									ast.NewNilLiteralNode(P(27, 3, 4, 2)),
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
				P(0, 20, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 19, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 18, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), token.GREATER),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), token.DEC_INT, "0")),
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
				P(0, 43, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 38, 2, 1),
						ast.NewAssignmentExpressionNode(
							P(1, 37, 2, 1),
							T(P(5, 1, 2, 5), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(1, 3, 2, 1), "bar"),
							ast.NewUnlessExpressionNode(
								P(8, 30, 3, 2),
								ast.NewBinaryExpressionNode(
									P(15, 7, 3, 9),
									T(P(19, 1, 3, 13), token.GREATER),
									ast.NewPublicIdentifierNode(P(15, 3, 3, 9), "foo"),
									ast.NewIntLiteralNode(P(21, 1, 3, 15), V(P(21, 1, 3, 15), token.DEC_INT, "0")),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										P(25, 9, 4, 3),
										ast.NewAssignmentExpressionNode(
											P(25, 8, 4, 3),
											T(P(29, 2, 4, 7), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(P(25, 3, 4, 3), "foo"),
											ast.NewIntLiteralNode(P(32, 1, 4, 10), V(P(32, 1, 4, 10), token.DEC_INT, "2")),
										),
									),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						P(39, 4, 6, 1),
						ast.NewNilLiteralNode(P(39, 3, 6, 1)),
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
				P(0, 34, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 29, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 28, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), token.GREATER),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(21, 8, 2, 21),
									ast.NewAssignmentExpressionNode(
										P(21, 8, 2, 21),
										T(P(25, 2, 2, 25), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(21, 3, 2, 21), "foo"),
										ast.NewIntLiteralNode(P(28, 1, 2, 28), V(P(28, 1, 2, 28), token.DEC_INT, "2")),
									),
								),
							},
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						P(30, 4, 3, 1),
						ast.NewNilLiteralNode(P(30, 3, 3, 1)),
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
				P(0, 59, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 54, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 53, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), token.GREATER),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 3, 2),
										T(P(21, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(17, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 3, 9), V(P(24, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(27, 4, 4, 2),
									ast.NewNilLiteralNode(P(27, 3, 4, 2)),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(37, 9, 6, 2),
									ast.NewAssignmentExpressionNode(
										P(37, 8, 6, 2),
										T(P(41, 2, 6, 6), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(P(37, 3, 6, 2), "foo"),
										ast.NewIntLiteralNode(P(44, 1, 6, 9), V(P(44, 1, 6, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(47, 4, 7, 2),
									ast.NewNilLiteralNode(P(47, 3, 7, 2)),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(55, 4, 9, 1),
						ast.NewNilLiteralNode(P(55, 3, 9, 1)),
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
				P(0, 48, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 43, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 42, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), token.GREATER),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(21, 8, 2, 21),
									ast.NewAssignmentExpressionNode(
										P(21, 8, 2, 21),
										T(P(25, 2, 2, 25), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(21, 3, 2, 21), "foo"),
										ast.NewIntLiteralNode(P(28, 1, 2, 28), V(P(28, 1, 2, 28), token.DEC_INT, "2")),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(35, 8, 3, 6),
									ast.NewAssignmentExpressionNode(
										P(35, 8, 3, 6),
										T(P(39, 2, 3, 10), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(P(35, 3, 3, 6), "foo"),
										ast.NewIntLiteralNode(P(42, 1, 3, 13), V(P(42, 1, 3, 13), token.DEC_INT, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(44, 4, 4, 1),
						ast.NewNilLiteralNode(P(44, 3, 4, 1)),
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
				P(0, 57, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 43, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 42, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), token.GREATER),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(21, 8, 2, 21),
									ast.NewAssignmentExpressionNode(
										P(21, 8, 2, 21),
										T(P(25, 2, 2, 25), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(21, 3, 2, 21), "foo"),
										ast.NewIntLiteralNode(P(28, 1, 2, 28), V(P(28, 1, 2, 28), token.DEC_INT, "2")),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(35, 8, 3, 6),
									ast.NewAssignmentExpressionNode(
										P(35, 8, 3, 6),
										T(P(39, 2, 3, 10), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(P(35, 3, 3, 6), "foo"),
										ast.NewIntLiteralNode(P(42, 1, 3, 13), V(P(42, 1, 3, 13), token.DEC_INT, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(44, 9, 4, 1),
						ast.NewInvalidNode(P(44, 4, 4, 1), T(P(44, 4, 4, 1), token.ELSE)),
					),
					ast.NewExpressionStatementNode(
						P(53, 4, 5, 1),
						ast.NewNilLiteralNode(P(53, 3, 5, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(P(44, 4, 4, 1), "unexpected else, expected an expression"),
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
				P(0, 34, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 33, 2, 1),
						ast.NewWhileExpressionNode(
							P(1, 32, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), token.GREATER),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(16, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(16, 8, 3, 2),
										T(P(20, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(16, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(23, 1, 3, 9), V(P(23, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(26, 4, 4, 2),
									ast.NewNilLiteralNode(P(26, 3, 4, 2)),
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
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 18, 2, 1),
						ast.NewWhileExpressionNode(
							P(1, 17, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), token.GREATER),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), token.DEC_INT, "0")),
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
				P(0, 42, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 37, 2, 1),
						ast.NewAssignmentExpressionNode(
							P(1, 36, 2, 1),
							T(P(5, 1, 2, 5), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(1, 3, 2, 1), "bar"),
							ast.NewWhileExpressionNode(
								P(8, 29, 3, 2),
								ast.NewBinaryExpressionNode(
									P(14, 7, 3, 8),
									T(P(18, 1, 3, 12), token.GREATER),
									ast.NewPublicIdentifierNode(P(14, 3, 3, 8), "foo"),
									ast.NewIntLiteralNode(P(20, 1, 3, 14), V(P(20, 1, 3, 14), token.DEC_INT, "0")),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										P(24, 9, 4, 3),
										ast.NewAssignmentExpressionNode(
											P(24, 8, 4, 3),
											T(P(28, 2, 4, 7), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(P(24, 3, 4, 3), "foo"),
											ast.NewIntLiteralNode(P(31, 1, 4, 10), V(P(31, 1, 4, 10), token.DEC_INT, "2")),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						P(38, 4, 6, 1),
						ast.NewNilLiteralNode(P(38, 3, 6, 1)),
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
				P(0, 33, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 28, 2, 1),
						ast.NewWhileExpressionNode(
							P(1, 27, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), token.GREATER),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(20, 8, 2, 20),
									ast.NewAssignmentExpressionNode(
										P(20, 8, 2, 20),
										T(P(24, 2, 2, 24), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(20, 3, 2, 20), "foo"),
										ast.NewIntLiteralNode(P(27, 1, 2, 27), V(P(27, 1, 2, 27), token.DEC_INT, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(29, 4, 3, 1),
						ast.NewNilLiteralNode(P(29, 3, 3, 1)),
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
				P(0, 58, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 53, 2, 1),
						ast.NewWhileExpressionNode(
							P(1, 52, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), token.GREATER),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(16, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(16, 8, 3, 2),
										T(P(20, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(16, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(23, 1, 3, 9), V(P(23, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(26, 4, 4, 2),
									ast.NewNilLiteralNode(P(26, 3, 4, 2)),
								),
								ast.NewExpressionStatementNode(
									P(30, 5, 5, 1),
									ast.NewInvalidNode(P(30, 4, 5, 1), T(P(30, 4, 5, 1), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									P(36, 9, 6, 2),
									ast.NewAssignmentExpressionNode(
										P(36, 8, 6, 2),
										T(P(40, 2, 6, 6), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(P(36, 3, 6, 2), "foo"),
										ast.NewIntLiteralNode(P(43, 1, 6, 9), V(P(43, 1, 6, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(46, 4, 7, 2),
									ast.NewNilLiteralNode(P(46, 3, 7, 2)),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(54, 4, 9, 1),
						ast.NewNilLiteralNode(P(54, 3, 9, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(P(30, 4, 5, 1), "unexpected else, expected an expression"),
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
				P(0, 34, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 33, 2, 1),
						ast.NewUntilExpressionNode(
							P(1, 32, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), token.GREATER),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(16, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(16, 8, 3, 2),
										T(P(20, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(16, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(23, 1, 3, 9), V(P(23, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(26, 4, 4, 2),
									ast.NewNilLiteralNode(P(26, 3, 4, 2)),
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
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 18, 2, 1),
						ast.NewUntilExpressionNode(
							P(1, 17, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), token.GREATER),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), token.DEC_INT, "0")),
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
				P(0, 42, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 37, 2, 1),
						ast.NewAssignmentExpressionNode(
							P(1, 36, 2, 1),
							T(P(5, 1, 2, 5), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(1, 3, 2, 1), "bar"),
							ast.NewUntilExpressionNode(
								P(8, 29, 3, 2),
								ast.NewBinaryExpressionNode(
									P(14, 7, 3, 8),
									T(P(18, 1, 3, 12), token.GREATER),
									ast.NewPublicIdentifierNode(P(14, 3, 3, 8), "foo"),
									ast.NewIntLiteralNode(P(20, 1, 3, 14), V(P(20, 1, 3, 14), token.DEC_INT, "0")),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										P(24, 9, 4, 3),
										ast.NewAssignmentExpressionNode(
											P(24, 8, 4, 3),
											T(P(28, 2, 4, 7), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(P(24, 3, 4, 3), "foo"),
											ast.NewIntLiteralNode(P(31, 1, 4, 10), V(P(31, 1, 4, 10), token.DEC_INT, "2")),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						P(38, 4, 6, 1),
						ast.NewNilLiteralNode(P(38, 3, 6, 1)),
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
				P(0, 33, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 28, 2, 1),
						ast.NewUntilExpressionNode(
							P(1, 27, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), token.GREATER),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(20, 8, 2, 20),
									ast.NewAssignmentExpressionNode(
										P(20, 8, 2, 20),
										T(P(24, 2, 2, 24), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(20, 3, 2, 20), "foo"),
										ast.NewIntLiteralNode(P(27, 1, 2, 27), V(P(27, 1, 2, 27), token.DEC_INT, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(29, 4, 3, 1),
						ast.NewNilLiteralNode(P(29, 3, 3, 1)),
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
				P(0, 58, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 53, 2, 1),
						ast.NewUntilExpressionNode(
							P(1, 52, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), token.GREATER),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), token.DEC_INT, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(16, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(16, 8, 3, 2),
										T(P(20, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(16, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(23, 1, 3, 9), V(P(23, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(26, 4, 4, 2),
									ast.NewNilLiteralNode(P(26, 3, 4, 2)),
								),
								ast.NewExpressionStatementNode(
									P(30, 5, 5, 1),
									ast.NewInvalidNode(P(30, 4, 5, 1), T(P(30, 4, 5, 1), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									P(36, 9, 6, 2),
									ast.NewAssignmentExpressionNode(
										P(36, 8, 6, 2),
										T(P(40, 2, 6, 6), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(P(36, 3, 6, 2), "foo"),
										ast.NewIntLiteralNode(P(43, 1, 6, 9), V(P(43, 1, 6, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(46, 4, 7, 2),
									ast.NewNilLiteralNode(P(46, 3, 7, 2)),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(54, 4, 9, 1),
						ast.NewNilLiteralNode(P(54, 3, 9, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(P(30, 4, 5, 1), "unexpected else, expected an expression"),
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
				P(0, 25, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 24, 2, 1),
						ast.NewLoopExpressionNode(
							P(1, 23, 2, 1),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(7, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(7, 8, 3, 2),
										T(P(11, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(7, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(14, 1, 3, 9), V(P(14, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(17, 4, 4, 2),
									ast.NewNilLiteralNode(P(17, 3, 4, 2)),
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
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 9, 2, 1),
						ast.NewLoopExpressionNode(
							P(1, 8, 2, 1),
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
				P(0, 33, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 28, 2, 1),
						ast.NewAssignmentExpressionNode(
							P(1, 27, 2, 1),
							T(P(5, 1, 2, 5), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(1, 3, 2, 1), "bar"),
							ast.NewLoopExpressionNode(
								P(8, 20, 3, 2),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										P(15, 9, 4, 3),
										ast.NewAssignmentExpressionNode(
											P(15, 8, 4, 3),
											T(P(19, 2, 4, 7), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(P(15, 3, 4, 3), "foo"),
											ast.NewIntLiteralNode(P(22, 1, 4, 10), V(P(22, 1, 4, 10), token.DEC_INT, "2")),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						P(29, 4, 6, 1),
						ast.NewNilLiteralNode(P(29, 3, 6, 1)),
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
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 14, 2, 1),
						ast.NewLoopExpressionNode(
							P(1, 13, 2, 1),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(6, 8, 2, 6),
									ast.NewAssignmentExpressionNode(
										P(6, 8, 2, 6),
										T(P(10, 2, 2, 10), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(6, 3, 2, 6), "foo"),
										ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), token.DEC_INT, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(15, 4, 3, 1),
						ast.NewNilLiteralNode(P(15, 3, 3, 1)),
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
				P(0, 49, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 44, 2, 1),
						ast.NewLoopExpressionNode(
							P(1, 43, 2, 1),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(7, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(7, 8, 3, 2),
										T(P(11, 2, 3, 6), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(7, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(14, 1, 3, 9), V(P(14, 1, 3, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(17, 4, 4, 2),
									ast.NewNilLiteralNode(P(17, 3, 4, 2)),
								),
								ast.NewExpressionStatementNode(
									P(21, 5, 5, 1),
									ast.NewInvalidNode(P(21, 4, 5, 1), T(P(21, 4, 5, 1), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									P(27, 9, 6, 2),
									ast.NewAssignmentExpressionNode(
										P(27, 8, 6, 2),
										T(P(31, 2, 6, 6), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(P(27, 3, 6, 2), "foo"),
										ast.NewIntLiteralNode(P(34, 1, 6, 9), V(P(34, 1, 6, 9), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(37, 4, 7, 2),
									ast.NewNilLiteralNode(P(37, 3, 7, 2)),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(45, 4, 9, 1),
						ast.NewNilLiteralNode(P(45, 3, 9, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(P(21, 4, 5, 1), "unexpected else, expected an expression"),
			},
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
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewBreakExpressionNode(P(0, 5, 1, 1)),
					),
				},
			),
		},
		"can't have an argument": {
			input: `break 2`,
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewBreakExpressionNode(P(0, 5, 1, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(
					P(6, 1, 1, 7),
					"unexpected DecInt, expected a statement separator `\\n`, `;`",
				),
			},
		},
		"is an expression": {
			input: `foo && break`,
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 12, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 12, 1, 1),
							T(P(4, 2, 1, 5), token.AND_AND),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBreakExpressionNode(P(7, 5, 1, 8)),
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
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewReturnExpressionNode(P(0, 6, 1, 1), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "return\n1",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewReturnExpressionNode(P(0, 6, 1, 1), nil),
					),
					ast.NewExpressionStatementNode(
						P(7, 1, 2, 1),
						ast.NewIntLiteralNode(P(7, 1, 2, 1), V(P(7, 1, 2, 1), token.DEC_INT, "1")),
					),
				},
			),
		},
		"can have an argument": {
			input: `return 2`,
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewReturnExpressionNode(
							P(0, 8, 1, 1),
							ast.NewIntLiteralNode(P(7, 1, 1, 8), V(P(7, 1, 1, 8), token.DEC_INT, "2")),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && return`,
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 13, 1, 1),
							T(P(4, 2, 1, 5), token.AND_AND),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewReturnExpressionNode(P(7, 6, 1, 8), nil),
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
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewContinueExpressionNode(P(0, 8, 1, 1), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "continue\n1",
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewContinueExpressionNode(P(0, 8, 1, 1), nil),
					),
					ast.NewExpressionStatementNode(
						P(9, 1, 2, 1),
						ast.NewIntLiteralNode(P(9, 1, 2, 1), V(P(9, 1, 2, 1), token.DEC_INT, "1")),
					),
				},
			),
		},
		"can have an argument": {
			input: `continue 2`,
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewContinueExpressionNode(
							P(0, 10, 1, 1),
							ast.NewIntLiteralNode(P(9, 1, 1, 10), V(P(9, 1, 1, 10), token.DEC_INT, "2")),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && continue`,
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 15, 1, 1),
							T(P(4, 2, 1, 5), token.AND_AND),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewContinueExpressionNode(P(7, 8, 1, 8), nil),
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
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewThrowExpressionNode(P(0, 5, 1, 1), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "throw\n1",
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewThrowExpressionNode(P(0, 5, 1, 1), nil),
					),
					ast.NewExpressionStatementNode(
						P(6, 1, 2, 1),
						ast.NewIntLiteralNode(P(6, 1, 2, 1), V(P(6, 1, 2, 1), token.DEC_INT, "1")),
					),
				},
			),
		},
		"can have an argument": {
			input: `throw 2`,
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewThrowExpressionNode(
							P(0, 7, 1, 1),
							ast.NewIntLiteralNode(P(6, 1, 1, 7), V(P(6, 1, 1, 7), token.DEC_INT, "2")),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && throw`,
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 12, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 12, 1, 1),
							T(P(4, 2, 1, 5), token.AND_AND),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewThrowExpressionNode(P(7, 5, 1, 8), nil),
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
