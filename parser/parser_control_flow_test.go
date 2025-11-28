package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/token"
)

func TestGoExpression(t *testing.T) {
	tests := testTable{
		"can be single line": {
			input: "go println('foo')",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewGoExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(3, 1, 4), P(16, 1, 17))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(3, 1, 4), P(16, 1, 17))),
										ast.NewPublicIdentifierNode(L(S(P(3, 1, 4), P(9, 1, 10))), "println"),
										[]ast.ExpressionNode{
											ast.NewRawStringLiteralNode(
												L(S(P(11, 1, 12), P(15, 1, 16))),
												"foo",
											),
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
		"can be multi line": {
			input: `
				go
					a := bar()
					println('foo')
				end
			`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(51, 5, 8))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(
						L(S(P(0, 1, 1), P(0, 1, 1))),
					),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 5), P(51, 5, 8))),
						ast.NewGoExpressionNode(
							L(S(P(5, 2, 5), P(43, 4, 20))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(13, 3, 6), P(23, 3, 16))),
									ast.NewAssignmentExpressionNode(
										L(S(P(13, 3, 6), P(22, 3, 15))),
										T(L(S(P(15, 3, 8), P(16, 3, 9))), token.COLON_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(13, 3, 6), P(13, 3, 6))), "a"),
										ast.NewReceiverlessMethodCallNode(
											L(S(P(18, 3, 11), P(22, 3, 15))),
											ast.NewPublicIdentifierNode(L(S(P(18, 3, 11), P(20, 3, 13))), "bar"),
											nil,
											nil,
										),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(29, 4, 6), P(43, 4, 20))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(29, 4, 6), P(42, 4, 19))),
										ast.NewPublicIdentifierNode(L(S(P(29, 4, 6), P(35, 4, 12))), "println"),
										[]ast.ExpressionNode{
											ast.NewRawStringLiteralNode(
												L(S(P(37, 4, 14), P(41, 4, 18))),
												"foo",
											),
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

func TestModifierExpression(t *testing.T) {
	tests := testTable{
		"has lower precedence than assignment": {
			input: "foo = bar if baz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							T(L(S(P(10, 1, 11), P(11, 1, 12))), token.IF),
							ast.NewAssignmentExpressionNode(
								L(S(P(0, 1, 1), P(8, 1, 9))),
								T(L(S(P(4, 1, 5), P(4, 1, 5))), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(15, 1, 16))), "baz"),
						),
					),
				},
			),
		},
		"can have newlines after the modifier keyword": {
			input: "foo = bar if\nbaz",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 2, 3))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 2, 3))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(15, 2, 3))),
							T(L(S(P(10, 1, 11), P(11, 1, 12))), token.IF),
							ast.NewAssignmentExpressionNode(
								L(S(P(0, 1, 1), P(8, 1, 9))),
								T(L(S(P(4, 1, 5), P(4, 1, 5))), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(13, 2, 1), P(15, 2, 3))), "baz"),
						),
					),
				},
			),
		},
		"if can contain else": {
			input: "foo = bar if baz else car = red",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(30, 1, 31))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(30, 1, 31))),
						ast.NewModifierIfElseNode(
							L(S(P(0, 1, 1), P(30, 1, 31))),
							ast.NewAssignmentExpressionNode(
								L(S(P(0, 1, 1), P(8, 1, 9))),
								T(L(S(P(4, 1, 5), P(4, 1, 5))), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(15, 1, 16))), "baz"),
							ast.NewAssignmentExpressionNode(
								L(S(P(22, 1, 23), P(30, 1, 31))),
								T(L(S(P(26, 1, 27), P(26, 1, 27))), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(24, 1, 25))), "car"),
								ast.NewPublicIdentifierNode(L(S(P(28, 1, 29), P(30, 1, 31))), "red"),
							),
						),
					),
				},
			),
		},
		"if else can span multiple lines": {
			input: "foo = bar if\nbaz else\ncar = red",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(30, 3, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(30, 3, 9))),
						ast.NewModifierIfElseNode(
							L(S(P(0, 1, 1), P(30, 3, 9))),
							ast.NewAssignmentExpressionNode(
								L(S(P(0, 1, 1), P(8, 1, 9))),
								T(L(S(P(4, 1, 5), P(4, 1, 5))), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(13, 2, 1), P(15, 2, 3))), "baz"),
							ast.NewAssignmentExpressionNode(
								L(S(P(22, 3, 1), P(30, 3, 9))),
								T(L(S(P(26, 3, 5), P(26, 3, 5))), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(L(S(P(22, 3, 1), P(24, 3, 3))), "car"),
								ast.NewPublicIdentifierNode(L(S(P(28, 3, 7), P(30, 3, 9))), "red"),
							),
						),
					),
				},
			),
		},
		"can have for loops": {
			input: "println(i) for i in [1, 2, 3]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(28, 1, 29))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(28, 1, 29))),
						ast.NewModifierForInNode(
							L(S(P(0, 1, 1), P(28, 1, 29))),
							ast.NewReceiverlessMethodCallNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(6, 1, 7))), "println"),
								[]ast.ExpressionNode{
									ast.NewPublicIdentifierNode(L(S(P(8, 1, 9), P(8, 1, 9))), "i"),
								},
								nil,
							),
							ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(15, 1, 16))), "i"),
							ast.NewArrayListLiteralNode(
								L(S(P(20, 1, 21), P(28, 1, 29))),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "1"),
									ast.NewIntLiteralNode(L(S(P(24, 1, 25), P(24, 1, 25))), "2"),
									ast.NewIntLiteralNode(L(S(P(27, 1, 28), P(27, 1, 28))), "3"),
								},
								nil,
							),
						),
					),
				},
			),
		},
		"can have for loops with patterns": {
			input: "println(a, b) for [a, b] in [[1, 2], [3, 4]]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(43, 1, 44))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(43, 1, 44))),
						ast.NewModifierForInNode(
							L(S(P(0, 1, 1), P(43, 1, 44))),
							ast.NewReceiverlessMethodCallNode(
								L(S(P(0, 1, 1), P(12, 1, 13))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(6, 1, 7))), "println"),
								[]ast.ExpressionNode{
									ast.NewPublicIdentifierNode(L(S(P(8, 1, 9), P(8, 1, 9))), "a"),
									ast.NewPublicIdentifierNode(L(S(P(11, 1, 12), P(11, 1, 12))), "b"),
								},
								nil,
							),
							ast.NewListPatternNode(
								L(S(P(18, 1, 19), P(23, 1, 24))),
								[]ast.PatternNode{
									ast.NewPublicIdentifierNode(L(S(P(19, 1, 20), P(19, 1, 20))), "a"),
									ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(22, 1, 23))), "b"),
								},
							),
							ast.NewArrayListLiteralNode(
								L(S(P(28, 1, 29), P(43, 1, 44))),
								[]ast.ExpressionNode{
									ast.NewArrayListLiteralNode(
										L(S(P(29, 1, 30), P(34, 1, 35))),
										[]ast.ExpressionNode{
											ast.NewIntLiteralNode(L(S(P(30, 1, 31), P(30, 1, 31))), "1"),
											ast.NewIntLiteralNode(L(S(P(33, 1, 34), P(33, 1, 34))), "2"),
										},
										nil,
									),
									ast.NewArrayListLiteralNode(
										L(S(P(37, 1, 38), P(42, 1, 43))),
										[]ast.ExpressionNode{
											ast.NewIntLiteralNode(L(S(P(38, 1, 39), P(38, 1, 39))), "3"),
											ast.NewIntLiteralNode(L(S(P(41, 1, 42), P(41, 1, 42))), "4"),
										},
										nil,
									),
								},
								nil,
							),
						),
					),
				},
			),
		},
		"can have for loops with pattern without variables": {
			input: "println(a, b) for [1, 2] in [[1, 2], [3, 4]]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(43, 1, 44))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(43, 1, 44))),
						ast.NewModifierForInNode(
							L(S(P(0, 1, 1), P(43, 1, 44))),
							ast.NewReceiverlessMethodCallNode(
								L(S(P(0, 1, 1), P(12, 1, 13))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(6, 1, 7))), "println"),
								[]ast.ExpressionNode{
									ast.NewPublicIdentifierNode(L(S(P(8, 1, 9), P(8, 1, 9))), "a"),
									ast.NewPublicIdentifierNode(L(S(P(11, 1, 12), P(11, 1, 12))), "b"),
								},
								nil,
							),
							ast.NewListPatternNode(
								L(S(P(18, 1, 19), P(23, 1, 24))),
								[]ast.PatternNode{
									ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "1"),
									ast.NewIntLiteralNode(L(S(P(22, 1, 23), P(22, 1, 23))), "2"),
								},
							),
							ast.NewArrayListLiteralNode(
								L(S(P(28, 1, 29), P(43, 1, 44))),
								[]ast.ExpressionNode{
									ast.NewArrayListLiteralNode(
										L(S(P(29, 1, 30), P(34, 1, 35))),
										[]ast.ExpressionNode{
											ast.NewIntLiteralNode(L(S(P(30, 1, 31), P(30, 1, 31))), "1"),
											ast.NewIntLiteralNode(L(S(P(33, 1, 34), P(33, 1, 34))), "2"),
										},
										nil,
									),
									ast.NewArrayListLiteralNode(
										L(S(P(37, 1, 38), P(42, 1, 43))),
										[]ast.ExpressionNode{
											ast.NewIntLiteralNode(L(S(P(38, 1, 39), P(38, 1, 39))), "3"),
											ast.NewIntLiteralNode(L(S(P(41, 1, 42), P(41, 1, 42))), "4"),
										},
										nil,
									),
								},
								nil,
							),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(18, 1, 19), P(23, 1, 24))), "patterns in for in loops should define at least one variable"),
			},
		},
		"for loops can span multiple lines": {
			input: "println(i) for\ni\nin\n[1,\n2,\n3]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(28, 6, 2))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(28, 6, 2))),
						ast.NewModifierForInNode(
							L(S(P(0, 1, 1), P(28, 6, 2))),
							ast.NewReceiverlessMethodCallNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(6, 1, 7))), "println"),
								[]ast.ExpressionNode{
									ast.NewPublicIdentifierNode(L(S(P(8, 1, 9), P(8, 1, 9))), "i"),
								},
								nil,
							),
							ast.NewPublicIdentifierNode(L(S(P(15, 2, 1), P(15, 2, 1))), "i"),
							ast.NewArrayListLiteralNode(
								L(S(P(20, 4, 1), P(28, 6, 2))),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(L(S(P(21, 4, 2), P(21, 4, 2))), "1"),
									ast.NewIntLiteralNode(L(S(P(24, 5, 1), P(24, 5, 1))), "2"),
									ast.NewIntLiteralNode(L(S(P(27, 6, 1), P(27, 6, 1))), "3"),
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
				L(S(P(0, 1, 1), P(52, 4, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.IF),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(11, 2, 1), P(25, 2, 15))),
						ast.NewModifierNode(
							L(S(P(11, 2, 1), P(24, 2, 14))),
							T(L(S(P(15, 2, 5), P(20, 2, 10))), token.UNLESS),
							ast.NewPublicIdentifierNode(L(S(P(11, 2, 1), P(13, 2, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(22, 2, 12), P(24, 2, 14))), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(26, 3, 1), P(39, 3, 14))),
						ast.NewModifierNode(
							L(S(P(26, 3, 1), P(38, 3, 13))),
							T(L(S(P(30, 3, 5), P(34, 3, 9))), token.WHILE),
							ast.NewPublicIdentifierNode(L(S(P(26, 3, 1), P(28, 3, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(36, 3, 11), P(38, 3, 13))), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(40, 4, 1), P(52, 4, 13))),
						ast.NewModifierNode(
							L(S(P(40, 4, 1), P(52, 4, 13))),
							T(L(S(P(44, 4, 5), P(48, 4, 9))), token.UNTIL),
							ast.NewPublicIdentifierNode(L(S(P(40, 4, 1), P(42, 4, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(50, 4, 11), P(52, 4, 13))), "bar"),
						),
					),
				},
			),
		},
		"cannot be nested": {
			input: "foo = bar if baz if false\n3",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(26, 2, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							T(L(S(P(10, 1, 11), P(11, 1, 12))), token.IF),
							ast.NewAssignmentExpressionNode(
								L(S(P(0, 1, 1), P(8, 1, 9))),
								T(L(S(P(4, 1, 5), P(4, 1, 5))), token.EQUAL_OP),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
							),
							ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(15, 1, 16))), "baz"),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(26, 2, 1), P(26, 2, 1))),
						ast.NewIntLiteralNode(L(S(P(26, 2, 1), P(26, 2, 1))), "3"),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(17, 1, 18), P(18, 1, 19))), "unexpected if, expected a statement separator `\\n`, `;`"),
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
				L(S(P(0, 1, 1), P(30, 5, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(30, 5, 4))),
						ast.NewIfExpressionNode(
							L(S(P(1, 2, 1), P(29, 5, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(4, 2, 4), P(10, 2, 10))),
								T(L(S(P(8, 2, 8), P(8, 2, 8))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(4, 2, 4), P(6, 2, 6))), "foo"),
								ast.NewIntLiteralNode(L(S(P(10, 2, 10), P(10, 2, 10))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(13, 3, 2), P(21, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(13, 3, 2), P(20, 3, 9))),
										T(L(S(P(17, 3, 6), P(18, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(13, 3, 2), P(15, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(20, 3, 9), P(20, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(23, 4, 2), P(26, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(23, 4, 2), P(25, 4, 4)))),
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
				L(S(P(0, 1, 1), P(15, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(15, 3, 4))),
						ast.NewIfExpressionNode(
							L(S(P(1, 2, 1), P(14, 3, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(4, 2, 4), P(10, 2, 10))),
								T(L(S(P(8, 2, 8), P(8, 2, 8))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(4, 2, 4), P(6, 2, 6))), "foo"),
								ast.NewIntLiteralNode(L(S(P(10, 2, 10), P(10, 2, 10))), "0"),
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
				L(S(P(0, 1, 1), P(38, 6, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(34, 5, 5))),
						ast.NewAssignmentExpressionNode(
							L(S(P(1, 2, 1), P(33, 5, 4))),
							T(L(S(P(5, 2, 5), P(5, 2, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(1, 2, 1), P(3, 2, 3))), "bar"),
							ast.NewIfExpressionNode(
								L(S(P(8, 3, 2), P(33, 5, 4))),
								ast.NewBinaryExpressionNode(
									L(S(P(11, 3, 5), P(17, 3, 11))),
									T(L(S(P(15, 3, 9), P(15, 3, 9))), token.GREATER),
									ast.NewPublicIdentifierNode(L(S(P(11, 3, 5), P(13, 3, 7))), "foo"),
									ast.NewIntLiteralNode(L(S(P(17, 3, 11), P(17, 3, 11))), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										L(S(P(21, 4, 3), P(29, 4, 11))),
										ast.NewAssignmentExpressionNode(
											L(S(P(21, 4, 3), P(28, 4, 10))),
											T(L(S(P(25, 4, 7), P(26, 4, 8))), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(L(S(P(21, 4, 3), P(23, 4, 5))), "foo"),
											ast.NewIntLiteralNode(L(S(P(28, 4, 10), P(28, 4, 10))), "2"),
										),
									),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(35, 6, 1), P(38, 6, 4))),
						ast.NewNilLiteralNode(L(S(P(35, 6, 1), P(37, 6, 3)))),
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
				L(S(P(0, 1, 1), P(29, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(25, 2, 25))),
						ast.NewIfExpressionNode(
							L(S(P(1, 2, 1), P(24, 2, 24))),
							ast.NewBinaryExpressionNode(
								L(S(P(4, 2, 4), P(10, 2, 10))),
								T(L(S(P(8, 2, 8), P(8, 2, 8))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(4, 2, 4), P(6, 2, 6))), "foo"),
								ast.NewIntLiteralNode(L(S(P(10, 2, 10), P(10, 2, 10))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(17, 2, 17), P(24, 2, 24))),
									ast.NewAssignmentExpressionNode(
										L(S(P(17, 2, 17), P(24, 2, 24))),
										T(L(S(P(21, 2, 21), P(22, 2, 22))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(17, 2, 17), P(19, 2, 19))), "foo"),
										ast.NewIntLiteralNode(L(S(P(24, 2, 24), P(24, 2, 24))), "2"),
									),
								),
							},
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(26, 3, 1), P(29, 3, 4))),
						ast.NewNilLiteralNode(L(S(P(26, 3, 1), P(28, 3, 3)))),
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
				L(S(P(0, 1, 1), P(55, 9, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(51, 8, 4))),
						ast.NewIfExpressionNode(
							L(S(P(1, 2, 1), P(50, 8, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(4, 2, 4), P(10, 2, 10))),
								T(L(S(P(8, 2, 8), P(8, 2, 8))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(4, 2, 4), P(6, 2, 6))), "foo"),
								ast.NewIntLiteralNode(L(S(P(10, 2, 10), P(10, 2, 10))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(13, 3, 2), P(21, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(13, 3, 2), P(20, 3, 9))),
										T(L(S(P(17, 3, 6), P(18, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(13, 3, 2), P(15, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(20, 3, 9), P(20, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(23, 4, 2), P(26, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(23, 4, 2), P(25, 4, 4)))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(34, 6, 3), P(42, 6, 11))),
									ast.NewAssignmentExpressionNode(
										L(S(P(34, 6, 3), P(41, 6, 10))),
										T(L(S(P(38, 6, 7), P(39, 6, 8))), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(34, 6, 3), P(36, 6, 5))), "foo"),
										ast.NewIntLiteralNode(L(S(P(41, 6, 10), P(41, 6, 10))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(44, 7, 2), P(47, 7, 5))),
									ast.NewNilLiteralNode(L(S(P(44, 7, 2), P(46, 7, 4)))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(52, 9, 1), P(55, 9, 4))),
						ast.NewNilLiteralNode(L(S(P(52, 9, 1), P(54, 9, 3)))),
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
				L(S(P(0, 1, 1), P(43, 4, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(39, 3, 14))),
						ast.NewIfExpressionNode(
							L(S(P(1, 2, 1), P(38, 3, 13))),
							ast.NewBinaryExpressionNode(
								L(S(P(4, 2, 4), P(10, 2, 10))),
								T(L(S(P(8, 2, 8), P(8, 2, 8))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(4, 2, 4), P(6, 2, 6))), "foo"),
								ast.NewIntLiteralNode(L(S(P(10, 2, 10), P(10, 2, 10))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(17, 2, 17), P(24, 2, 24))),
									ast.NewAssignmentExpressionNode(
										L(S(P(17, 2, 17), P(24, 2, 24))),
										T(L(S(P(21, 2, 21), P(22, 2, 22))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(17, 2, 17), P(19, 2, 19))), "foo"),
										ast.NewIntLiteralNode(L(S(P(24, 2, 24), P(24, 2, 24))), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(31, 3, 6), P(38, 3, 13))),
									ast.NewAssignmentExpressionNode(
										L(S(P(31, 3, 6), P(38, 3, 13))),
										T(L(S(P(35, 3, 10), P(36, 3, 11))), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(31, 3, 6), P(33, 3, 8))), "foo"),
										ast.NewIntLiteralNode(L(S(P(38, 3, 13), P(38, 3, 13))), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(40, 4, 1), P(43, 4, 4))),
						ast.NewNilLiteralNode(L(S(P(40, 4, 1), P(42, 4, 3)))),
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
				L(S(P(0, 1, 1), P(52, 5, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(39, 3, 14))),
						ast.NewIfExpressionNode(
							L(S(P(1, 2, 1), P(38, 3, 13))),
							ast.NewBinaryExpressionNode(
								L(S(P(4, 2, 4), P(10, 2, 10))),
								T(L(S(P(8, 2, 8), P(8, 2, 8))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(4, 2, 4), P(6, 2, 6))), "foo"),
								ast.NewIntLiteralNode(L(S(P(10, 2, 10), P(10, 2, 10))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(17, 2, 17), P(24, 2, 24))),
									ast.NewAssignmentExpressionNode(
										L(S(P(17, 2, 17), P(24, 2, 24))),
										T(L(S(P(21, 2, 21), P(22, 2, 22))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(17, 2, 17), P(19, 2, 19))), "foo"),
										ast.NewIntLiteralNode(L(S(P(24, 2, 24), P(24, 2, 24))), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(31, 3, 6), P(38, 3, 13))),
									ast.NewAssignmentExpressionNode(
										L(S(P(31, 3, 6), P(38, 3, 13))),
										T(L(S(P(35, 3, 10), P(36, 3, 11))), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(31, 3, 6), P(33, 3, 8))), "foo"),
										ast.NewIntLiteralNode(L(S(P(38, 3, 13), P(38, 3, 13))), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(40, 4, 1), P(48, 4, 9))),
						ast.NewInvalidNode(L(S(P(40, 4, 1), P(43, 4, 4))), T(L(S(P(40, 4, 1), P(43, 4, 4))), token.ELSE)),
					),
					ast.NewExpressionStatementNode(
						L(S(P(49, 5, 1), P(52, 5, 4))),
						ast.NewNilLiteralNode(L(S(P(49, 5, 1), P(51, 5, 3)))),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(40, 4, 1), P(43, 4, 4))), "unexpected else, expected an expression"),
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
				L(S(P(0, 1, 1), P(103, 13, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(99, 12, 4))),
						ast.NewIfExpressionNode(
							L(S(P(1, 2, 1), P(98, 12, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(4, 2, 4), P(10, 2, 10))),
								T(L(S(P(8, 2, 8), P(8, 2, 8))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(4, 2, 4), P(6, 2, 6))), "foo"),
								ast.NewIntLiteralNode(L(S(P(10, 2, 10), P(10, 2, 10))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(13, 3, 2), P(21, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(13, 3, 2), P(20, 3, 9))),
										T(L(S(P(17, 3, 6), P(18, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(13, 3, 2), P(15, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(20, 3, 9), P(20, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(23, 4, 2), P(26, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(23, 4, 2), P(25, 4, 4)))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(27, 5, 1), P(51, 6, 11))),
									ast.NewIfExpressionNode(
										L(S(P(27, 5, 1), P(51, 6, 11))),
										ast.NewBinaryExpressionNode(
											L(S(P(33, 5, 7), P(39, 5, 13))),
											T(L(S(P(37, 5, 11), P(37, 5, 11))), token.LESS),
											ast.NewPublicIdentifierNode(L(S(P(33, 5, 7), P(35, 5, 9))), "foo"),
											ast.NewIntLiteralNode(L(S(P(39, 5, 13), P(39, 5, 13))), "5"),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												L(S(P(42, 6, 2), P(51, 6, 11))),
												ast.NewAssignmentExpressionNode(
													L(S(P(42, 6, 2), P(50, 6, 10))),
													T(L(S(P(46, 6, 6), P(47, 6, 7))), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(L(S(P(42, 6, 2), P(44, 6, 4))), "foo"),
													ast.NewIntLiteralNode(L(S(P(49, 6, 9), P(50, 6, 10))), "10"),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												L(S(P(52, 7, 1), P(75, 8, 10))),
												ast.NewIfExpressionNode(
													L(S(P(52, 7, 1), P(98, 12, 3))),
													ast.NewBinaryExpressionNode(
														L(S(P(58, 7, 7), P(64, 7, 13))),
														T(L(S(P(62, 7, 11), P(62, 7, 11))), token.LESS),
														ast.NewPublicIdentifierNode(L(S(P(58, 7, 7), P(60, 7, 9))), "foo"),
														ast.NewIntLiteralNode(L(S(P(64, 7, 13), P(64, 7, 13))), "0"),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															L(S(P(67, 8, 2), P(75, 8, 10))),
															ast.NewAssignmentExpressionNode(
																L(S(P(67, 8, 2), P(74, 8, 9))),
																T(L(S(P(71, 8, 6), P(72, 8, 7))), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(L(S(P(67, 8, 2), P(69, 8, 4))), "foo"),
																ast.NewIntLiteralNode(L(S(P(74, 8, 9), P(74, 8, 9))), "3"),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															L(S(P(82, 10, 2), P(90, 10, 10))),
															ast.NewAssignmentExpressionNode(
																L(S(P(82, 10, 2), P(89, 10, 9))),
																T(L(S(P(86, 10, 6), P(87, 10, 7))), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(L(S(P(82, 10, 2), P(84, 10, 4))), "foo"),
																ast.NewIntLiteralNode(L(S(P(89, 10, 9), P(89, 10, 9))), "2"),
															),
														),
														ast.NewExpressionStatementNode(
															L(S(P(92, 11, 2), P(95, 11, 5))),
															ast.NewNilLiteralNode(L(S(P(92, 11, 2), P(94, 11, 4)))),
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
						L(S(P(100, 13, 1), P(103, 13, 4))),
						ast.NewNilLiteralNode(L(S(P(100, 13, 1), P(102, 13, 3)))),
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
				L(S(P(0, 1, 1), P(100, 6, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(96, 5, 14))),
						ast.NewIfExpressionNode(
							L(S(P(1, 2, 1), P(95, 5, 13))),
							ast.NewBinaryExpressionNode(
								L(S(P(4, 2, 4), P(10, 2, 10))),
								T(L(S(P(8, 2, 8), P(8, 2, 8))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(4, 2, 4), P(6, 2, 6))), "foo"),
								ast.NewIntLiteralNode(L(S(P(10, 2, 10), P(10, 2, 10))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(17, 2, 17), P(24, 2, 24))),
									ast.NewAssignmentExpressionNode(
										L(S(P(17, 2, 17), P(24, 2, 24))),
										T(L(S(P(21, 2, 21), P(22, 2, 22))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(17, 2, 17), P(19, 2, 19))), "foo"),
										ast.NewIntLiteralNode(L(S(P(24, 2, 24), P(24, 2, 24))), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(26, 3, 1), P(53, 3, 28))),
									ast.NewIfExpressionNode(
										L(S(P(26, 3, 1), P(53, 3, 28))),
										ast.NewBinaryExpressionNode(
											L(S(P(32, 3, 7), P(38, 3, 13))),
											T(L(S(P(36, 3, 11), P(36, 3, 11))), token.LESS),
											ast.NewPublicIdentifierNode(L(S(P(32, 3, 7), P(34, 3, 9))), "foo"),
											ast.NewIntLiteralNode(L(S(P(38, 3, 13), P(38, 3, 13))), "5"),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												L(S(P(45, 3, 20), P(53, 3, 28))),
												ast.NewAssignmentExpressionNode(
													L(S(P(45, 3, 20), P(53, 3, 28))),
													T(L(S(P(49, 3, 24), P(50, 3, 25))), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(L(S(P(45, 3, 20), P(47, 3, 22))), "foo"),
													ast.NewIntLiteralNode(L(S(P(52, 3, 27), P(53, 3, 28))), "10"),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												L(S(P(55, 4, 1), P(81, 4, 27))),
												ast.NewIfExpressionNode(
													L(S(P(55, 4, 1), P(95, 5, 13))),
													ast.NewBinaryExpressionNode(
														L(S(P(61, 4, 7), P(67, 4, 13))),
														T(L(S(P(65, 4, 11), P(65, 4, 11))), token.LESS),
														ast.NewPublicIdentifierNode(L(S(P(61, 4, 7), P(63, 4, 9))), "foo"),
														ast.NewIntLiteralNode(L(S(P(67, 4, 13), P(67, 4, 13))), "0"),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															L(S(P(74, 4, 20), P(81, 4, 27))),
															ast.NewAssignmentExpressionNode(
																L(S(P(74, 4, 20), P(81, 4, 27))),
																T(L(S(P(78, 4, 24), P(79, 4, 25))), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(L(S(P(74, 4, 20), P(76, 4, 22))), "foo"),
																ast.NewIntLiteralNode(L(S(P(81, 4, 27), P(81, 4, 27))), "3"),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															L(S(P(88, 5, 6), P(95, 5, 13))),
															ast.NewAssignmentExpressionNode(
																L(S(P(88, 5, 6), P(95, 5, 13))),
																T(L(S(P(92, 5, 10), P(93, 5, 11))), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(L(S(P(88, 5, 6), P(90, 5, 8))), "foo"),
																ast.NewIntLiteralNode(L(S(P(95, 5, 13), P(95, 5, 13))), "2"),
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
						L(S(P(97, 6, 1), P(100, 6, 4))),
						ast.NewNilLiteralNode(L(S(P(97, 6, 1), P(99, 6, 3)))),
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
				L(S(P(0, 1, 1), P(107, 13, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(103, 12, 4))),
						ast.NewIfExpressionNode(
							L(S(P(1, 2, 1), P(102, 12, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(4, 2, 4), P(10, 2, 10))),
								T(L(S(P(8, 2, 8), P(8, 2, 8))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(4, 2, 4), P(6, 2, 6))), "foo"),
								ast.NewIntLiteralNode(L(S(P(10, 2, 10), P(10, 2, 10))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(13, 3, 2), P(21, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(13, 3, 2), P(20, 3, 9))),
										T(L(S(P(17, 3, 6), P(18, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(13, 3, 2), P(15, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(20, 3, 9), P(20, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(23, 4, 2), P(26, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(23, 4, 2), P(25, 4, 4)))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(32, 5, 6), P(102, 12, 3))),
									ast.NewIfExpressionNode(
										L(S(P(32, 5, 6), P(102, 12, 3))),
										ast.NewBinaryExpressionNode(
											L(S(P(35, 5, 9), P(41, 5, 15))),
											T(L(S(P(39, 5, 13), P(39, 5, 13))), token.LESS),
											ast.NewPublicIdentifierNode(L(S(P(35, 5, 9), P(37, 5, 11))), "foo"),
											ast.NewIntLiteralNode(L(S(P(41, 5, 15), P(41, 5, 15))), "5"),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												L(S(P(44, 6, 2), P(53, 6, 11))),
												ast.NewAssignmentExpressionNode(
													L(S(P(44, 6, 2), P(52, 6, 10))),
													T(L(S(P(48, 6, 6), P(49, 6, 7))), token.STAR_EQUAL),
													ast.NewPublicIdentifierNode(L(S(P(44, 6, 2), P(46, 6, 4))), "foo"),
													ast.NewIntLiteralNode(L(S(P(51, 6, 9), P(52, 6, 10))), "10"),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												L(S(P(59, 7, 6), P(102, 12, 3))),
												ast.NewIfExpressionNode(
													L(S(P(59, 7, 6), P(102, 12, 3))),
													ast.NewBinaryExpressionNode(
														L(S(P(62, 7, 9), P(68, 7, 15))),
														T(L(S(P(66, 7, 13), P(66, 7, 13))), token.LESS),
														ast.NewPublicIdentifierNode(L(S(P(62, 7, 9), P(64, 7, 11))), "foo"),
														ast.NewIntLiteralNode(L(S(P(68, 7, 15), P(68, 7, 15))), "0"),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															L(S(P(71, 8, 2), P(79, 8, 10))),
															ast.NewAssignmentExpressionNode(
																L(S(P(71, 8, 2), P(78, 8, 9))),
																T(L(S(P(75, 8, 6), P(76, 8, 7))), token.PERCENT_EQUAL),
																ast.NewPublicIdentifierNode(L(S(P(71, 8, 2), P(73, 8, 4))), "foo"),
																ast.NewIntLiteralNode(L(S(P(78, 8, 9), P(78, 8, 9))), "3"),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															L(S(P(86, 10, 2), P(94, 10, 10))),
															ast.NewAssignmentExpressionNode(
																L(S(P(86, 10, 2), P(93, 10, 9))),
																T(L(S(P(90, 10, 6), P(91, 10, 7))), token.MINUS_EQUAL),
																ast.NewPublicIdentifierNode(L(S(P(86, 10, 2), P(88, 10, 4))), "foo"),
																ast.NewIntLiteralNode(L(S(P(93, 10, 9), P(93, 10, 9))), "2"),
															),
														),
														ast.NewExpressionStatementNode(
															L(S(P(96, 11, 2), P(99, 11, 5))),
															ast.NewNilLiteralNode(L(S(P(96, 11, 2), P(98, 11, 4)))),
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
						L(S(P(104, 13, 1), P(107, 13, 4))),
						ast.NewNilLiteralNode(L(S(P(104, 13, 1), P(106, 13, 3)))),
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
				L(S(P(0, 1, 1), P(34, 5, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(34, 5, 4))),
						ast.NewUnlessExpressionNode(
							L(S(P(1, 2, 1), P(33, 5, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(8, 2, 8), P(14, 2, 14))),
								T(L(S(P(12, 2, 12), P(12, 2, 12))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(8, 2, 8), P(10, 2, 10))), "foo"),
								ast.NewIntLiteralNode(L(S(P(14, 2, 14), P(14, 2, 14))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(17, 3, 2), P(25, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(17, 3, 2), P(24, 3, 9))),
										T(L(S(P(21, 3, 6), P(22, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(17, 3, 2), P(19, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(24, 3, 9), P(24, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(27, 4, 2), P(30, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(27, 4, 2), P(29, 4, 4)))),
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
				L(S(P(0, 1, 1), P(19, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(19, 3, 4))),
						ast.NewUnlessExpressionNode(
							L(S(P(1, 2, 1), P(18, 3, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(8, 2, 8), P(14, 2, 14))),
								T(L(S(P(12, 2, 12), P(12, 2, 12))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(8, 2, 8), P(10, 2, 10))), "foo"),
								ast.NewIntLiteralNode(L(S(P(14, 2, 14), P(14, 2, 14))), "0"),
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
				L(S(P(0, 1, 1), P(42, 6, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(38, 5, 5))),
						ast.NewAssignmentExpressionNode(
							L(S(P(1, 2, 1), P(37, 5, 4))),
							T(L(S(P(5, 2, 5), P(5, 2, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(1, 2, 1), P(3, 2, 3))), "bar"),
							ast.NewUnlessExpressionNode(
								L(S(P(8, 3, 2), P(37, 5, 4))),
								ast.NewBinaryExpressionNode(
									L(S(P(15, 3, 9), P(21, 3, 15))),
									T(L(S(P(19, 3, 13), P(19, 3, 13))), token.GREATER),
									ast.NewPublicIdentifierNode(L(S(P(15, 3, 9), P(17, 3, 11))), "foo"),
									ast.NewIntLiteralNode(L(S(P(21, 3, 15), P(21, 3, 15))), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										L(S(P(25, 4, 3), P(33, 4, 11))),
										ast.NewAssignmentExpressionNode(
											L(S(P(25, 4, 3), P(32, 4, 10))),
											T(L(S(P(29, 4, 7), P(30, 4, 8))), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(L(S(P(25, 4, 3), P(27, 4, 5))), "foo"),
											ast.NewIntLiteralNode(L(S(P(32, 4, 10), P(32, 4, 10))), "2"),
										),
									),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(39, 6, 1), P(42, 6, 4))),
						ast.NewNilLiteralNode(L(S(P(39, 6, 1), P(41, 6, 3)))),
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
				L(S(P(0, 1, 1), P(33, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(29, 2, 29))),
						ast.NewUnlessExpressionNode(
							L(S(P(1, 2, 1), P(28, 2, 28))),
							ast.NewBinaryExpressionNode(
								L(S(P(8, 2, 8), P(14, 2, 14))),
								T(L(S(P(12, 2, 12), P(12, 2, 12))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(8, 2, 8), P(10, 2, 10))), "foo"),
								ast.NewIntLiteralNode(L(S(P(14, 2, 14), P(14, 2, 14))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(21, 2, 21), P(28, 2, 28))),
									ast.NewAssignmentExpressionNode(
										L(S(P(21, 2, 21), P(28, 2, 28))),
										T(L(S(P(25, 2, 25), P(26, 2, 26))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(21, 2, 21), P(23, 2, 23))), "foo"),
										ast.NewIntLiteralNode(L(S(P(28, 2, 28), P(28, 2, 28))), "2"),
									),
								),
							},
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(30, 3, 1), P(33, 3, 4))),
						ast.NewNilLiteralNode(L(S(P(30, 3, 1), P(32, 3, 3)))),
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
				L(S(P(0, 1, 1), P(58, 9, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(54, 8, 4))),
						ast.NewUnlessExpressionNode(
							L(S(P(1, 2, 1), P(53, 8, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(8, 2, 8), P(14, 2, 14))),
								T(L(S(P(12, 2, 12), P(12, 2, 12))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(8, 2, 8), P(10, 2, 10))), "foo"),
								ast.NewIntLiteralNode(L(S(P(14, 2, 14), P(14, 2, 14))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(17, 3, 2), P(25, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(17, 3, 2), P(24, 3, 9))),
										T(L(S(P(21, 3, 6), P(22, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(17, 3, 2), P(19, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(24, 3, 9), P(24, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(27, 4, 2), P(30, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(27, 4, 2), P(29, 4, 4)))),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(37, 6, 2), P(45, 6, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(37, 6, 2), P(44, 6, 9))),
										T(L(S(P(41, 6, 6), P(42, 6, 7))), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(37, 6, 2), P(39, 6, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(44, 6, 9), P(44, 6, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(47, 7, 2), P(50, 7, 5))),
									ast.NewNilLiteralNode(L(S(P(47, 7, 2), P(49, 7, 4)))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(55, 9, 1), P(58, 9, 4))),
						ast.NewNilLiteralNode(L(S(P(55, 9, 1), P(57, 9, 3)))),
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
				L(S(P(0, 1, 1), P(47, 4, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(43, 3, 14))),
						ast.NewUnlessExpressionNode(
							L(S(P(1, 2, 1), P(42, 3, 13))),
							ast.NewBinaryExpressionNode(
								L(S(P(8, 2, 8), P(14, 2, 14))),
								T(L(S(P(12, 2, 12), P(12, 2, 12))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(8, 2, 8), P(10, 2, 10))), "foo"),
								ast.NewIntLiteralNode(L(S(P(14, 2, 14), P(14, 2, 14))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(21, 2, 21), P(28, 2, 28))),
									ast.NewAssignmentExpressionNode(
										L(S(P(21, 2, 21), P(28, 2, 28))),
										T(L(S(P(25, 2, 25), P(26, 2, 26))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(21, 2, 21), P(23, 2, 23))), "foo"),
										ast.NewIntLiteralNode(L(S(P(28, 2, 28), P(28, 2, 28))), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(35, 3, 6), P(42, 3, 13))),
									ast.NewAssignmentExpressionNode(
										L(S(P(35, 3, 6), P(42, 3, 13))),
										T(L(S(P(39, 3, 10), P(40, 3, 11))), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(35, 3, 6), P(37, 3, 8))), "foo"),
										ast.NewIntLiteralNode(L(S(P(42, 3, 13), P(42, 3, 13))), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(44, 4, 1), P(47, 4, 4))),
						ast.NewNilLiteralNode(L(S(P(44, 4, 1), P(46, 4, 3)))),
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
				L(S(P(0, 1, 1), P(56, 5, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(43, 3, 14))),
						ast.NewUnlessExpressionNode(
							L(S(P(1, 2, 1), P(42, 3, 13))),
							ast.NewBinaryExpressionNode(
								L(S(P(8, 2, 8), P(14, 2, 14))),
								T(L(S(P(12, 2, 12), P(12, 2, 12))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(8, 2, 8), P(10, 2, 10))), "foo"),
								ast.NewIntLiteralNode(L(S(P(14, 2, 14), P(14, 2, 14))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(21, 2, 21), P(28, 2, 28))),
									ast.NewAssignmentExpressionNode(
										L(S(P(21, 2, 21), P(28, 2, 28))),
										T(L(S(P(25, 2, 25), P(26, 2, 26))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(21, 2, 21), P(23, 2, 23))), "foo"),
										ast.NewIntLiteralNode(L(S(P(28, 2, 28), P(28, 2, 28))), "2"),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(35, 3, 6), P(42, 3, 13))),
									ast.NewAssignmentExpressionNode(
										L(S(P(35, 3, 6), P(42, 3, 13))),
										T(L(S(P(39, 3, 10), P(40, 3, 11))), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(35, 3, 6), P(37, 3, 8))), "foo"),
										ast.NewIntLiteralNode(L(S(P(42, 3, 13), P(42, 3, 13))), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(44, 4, 1), P(52, 4, 9))),
						ast.NewInvalidNode(L(S(P(44, 4, 1), P(47, 4, 4))), T(L(S(P(44, 4, 1), P(47, 4, 4))), token.ELSE)),
					),
					ast.NewExpressionStatementNode(
						L(S(P(53, 5, 1), P(56, 5, 4))),
						ast.NewNilLiteralNode(L(S(P(53, 5, 1), P(55, 5, 3)))),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(44, 4, 1), P(47, 4, 4))), "unexpected else, expected an expression"),
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
				L(S(P(0, 1, 1), P(33, 5, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(33, 5, 4))),
						ast.NewWhileExpressionNode(
							L(S(P(1, 2, 1), P(32, 5, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(7, 2, 7), P(13, 2, 13))),
								T(L(S(P(11, 2, 11), P(11, 2, 11))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(7, 2, 7), P(9, 2, 9))), "foo"),
								ast.NewIntLiteralNode(L(S(P(13, 2, 13), P(13, 2, 13))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(16, 3, 2), P(24, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(16, 3, 2), P(23, 3, 9))),
										T(L(S(P(20, 3, 6), P(21, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(16, 3, 2), P(18, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(23, 3, 9), P(23, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(26, 4, 2), P(29, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(26, 4, 2), P(28, 4, 4)))),
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
				L(S(P(0, 1, 1), P(18, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(18, 3, 4))),
						ast.NewWhileExpressionNode(
							L(S(P(1, 2, 1), P(17, 3, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(7, 2, 7), P(13, 2, 13))),
								T(L(S(P(11, 2, 11), P(11, 2, 11))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(7, 2, 7), P(9, 2, 9))), "foo"),
								ast.NewIntLiteralNode(L(S(P(13, 2, 13), P(13, 2, 13))), "0"),
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
				L(S(P(0, 1, 1), P(41, 6, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(37, 5, 5))),
						ast.NewAssignmentExpressionNode(
							L(S(P(1, 2, 1), P(36, 5, 4))),
							T(L(S(P(5, 2, 5), P(5, 2, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(1, 2, 1), P(3, 2, 3))), "bar"),
							ast.NewWhileExpressionNode(
								L(S(P(8, 3, 2), P(36, 5, 4))),
								ast.NewBinaryExpressionNode(
									L(S(P(14, 3, 8), P(20, 3, 14))),
									T(L(S(P(18, 3, 12), P(18, 3, 12))), token.GREATER),
									ast.NewPublicIdentifierNode(L(S(P(14, 3, 8), P(16, 3, 10))), "foo"),
									ast.NewIntLiteralNode(L(S(P(20, 3, 14), P(20, 3, 14))), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										L(S(P(24, 4, 3), P(32, 4, 11))),
										ast.NewAssignmentExpressionNode(
											L(S(P(24, 4, 3), P(31, 4, 10))),
											T(L(S(P(28, 4, 7), P(29, 4, 8))), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(L(S(P(24, 4, 3), P(26, 4, 5))), "foo"),
											ast.NewIntLiteralNode(L(S(P(31, 4, 10), P(31, 4, 10))), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(38, 6, 1), P(41, 6, 4))),
						ast.NewNilLiteralNode(L(S(P(38, 6, 1), P(40, 6, 3)))),
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
				L(S(P(0, 1, 1), P(32, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(28, 2, 28))),
						ast.NewWhileExpressionNode(
							L(S(P(1, 2, 1), P(27, 2, 27))),
							ast.NewBinaryExpressionNode(
								L(S(P(7, 2, 7), P(13, 2, 13))),
								T(L(S(P(11, 2, 11), P(11, 2, 11))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(7, 2, 7), P(9, 2, 9))), "foo"),
								ast.NewIntLiteralNode(L(S(P(13, 2, 13), P(13, 2, 13))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(20, 2, 20), P(27, 2, 27))),
									ast.NewAssignmentExpressionNode(
										L(S(P(20, 2, 20), P(27, 2, 27))),
										T(L(S(P(24, 2, 24), P(25, 2, 25))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(20, 2, 20), P(22, 2, 22))), "foo"),
										ast.NewIntLiteralNode(L(S(P(27, 2, 27), P(27, 2, 27))), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(29, 3, 1), P(32, 3, 4))),
						ast.NewNilLiteralNode(L(S(P(29, 3, 1), P(31, 3, 3)))),
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
				L(S(P(0, 1, 1), P(57, 9, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(53, 8, 4))),
						ast.NewWhileExpressionNode(
							L(S(P(1, 2, 1), P(52, 8, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(7, 2, 7), P(13, 2, 13))),
								T(L(S(P(11, 2, 11), P(11, 2, 11))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(7, 2, 7), P(9, 2, 9))), "foo"),
								ast.NewIntLiteralNode(L(S(P(13, 2, 13), P(13, 2, 13))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(16, 3, 2), P(24, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(16, 3, 2), P(23, 3, 9))),
										T(L(S(P(20, 3, 6), P(21, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(16, 3, 2), P(18, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(23, 3, 9), P(23, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(26, 4, 2), P(29, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(26, 4, 2), P(28, 4, 4)))),
								),
								ast.NewExpressionStatementNode(
									L(S(P(30, 5, 1), P(34, 5, 5))),
									ast.NewInvalidNode(L(S(P(30, 5, 1), P(33, 5, 4))), T(L(S(P(30, 5, 1), P(33, 5, 4))), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									L(S(P(36, 6, 2), P(44, 6, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(36, 6, 2), P(43, 6, 9))),
										T(L(S(P(40, 6, 6), P(41, 6, 7))), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(36, 6, 2), P(38, 6, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(43, 6, 9), P(43, 6, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(46, 7, 2), P(49, 7, 5))),
									ast.NewNilLiteralNode(L(S(P(46, 7, 2), P(48, 7, 4)))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(54, 9, 1), P(57, 9, 4))),
						ast.NewNilLiteralNode(L(S(P(54, 9, 1), P(56, 9, 3)))),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(30, 5, 1), P(33, 5, 4))), "unexpected else, expected an expression"),
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
				L(S(P(0, 1, 1), P(33, 5, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(33, 5, 4))),
						ast.NewUntilExpressionNode(
							L(S(P(1, 2, 1), P(32, 5, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(7, 2, 7), P(13, 2, 13))),
								T(L(S(P(11, 2, 11), P(11, 2, 11))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(7, 2, 7), P(9, 2, 9))), "foo"),
								ast.NewIntLiteralNode(L(S(P(13, 2, 13), P(13, 2, 13))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(16, 3, 2), P(24, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(16, 3, 2), P(23, 3, 9))),
										T(L(S(P(20, 3, 6), P(21, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(16, 3, 2), P(18, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(23, 3, 9), P(23, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(26, 4, 2), P(29, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(26, 4, 2), P(28, 4, 4)))),
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
				L(S(P(0, 1, 1), P(18, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(18, 3, 4))),
						ast.NewUntilExpressionNode(
							L(S(P(1, 2, 1), P(17, 3, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(7, 2, 7), P(13, 2, 13))),
								T(L(S(P(11, 2, 11), P(11, 2, 11))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(7, 2, 7), P(9, 2, 9))), "foo"),
								ast.NewIntLiteralNode(L(S(P(13, 2, 13), P(13, 2, 13))), "0"),
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
				L(S(P(0, 1, 1), P(41, 6, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(37, 5, 5))),
						ast.NewAssignmentExpressionNode(
							L(S(P(1, 2, 1), P(36, 5, 4))),
							T(L(S(P(5, 2, 5), P(5, 2, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(1, 2, 1), P(3, 2, 3))), "bar"),
							ast.NewUntilExpressionNode(
								L(S(P(8, 3, 2), P(36, 5, 4))),
								ast.NewBinaryExpressionNode(
									L(S(P(14, 3, 8), P(20, 3, 14))),
									T(L(S(P(18, 3, 12), P(18, 3, 12))), token.GREATER),
									ast.NewPublicIdentifierNode(L(S(P(14, 3, 8), P(16, 3, 10))), "foo"),
									ast.NewIntLiteralNode(L(S(P(20, 3, 14), P(20, 3, 14))), "0"),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										L(S(P(24, 4, 3), P(32, 4, 11))),
										ast.NewAssignmentExpressionNode(
											L(S(P(24, 4, 3), P(31, 4, 10))),
											T(L(S(P(28, 4, 7), P(29, 4, 8))), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(L(S(P(24, 4, 3), P(26, 4, 5))), "foo"),
											ast.NewIntLiteralNode(L(S(P(31, 4, 10), P(31, 4, 10))), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(38, 6, 1), P(41, 6, 4))),
						ast.NewNilLiteralNode(L(S(P(38, 6, 1), P(40, 6, 3)))),
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
				L(S(P(0, 1, 1), P(32, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(28, 2, 28))),
						ast.NewUntilExpressionNode(
							L(S(P(1, 2, 1), P(27, 2, 27))),
							ast.NewBinaryExpressionNode(
								L(S(P(7, 2, 7), P(13, 2, 13))),
								T(L(S(P(11, 2, 11), P(11, 2, 11))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(7, 2, 7), P(9, 2, 9))), "foo"),
								ast.NewIntLiteralNode(L(S(P(13, 2, 13), P(13, 2, 13))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(20, 2, 20), P(27, 2, 27))),
									ast.NewAssignmentExpressionNode(
										L(S(P(20, 2, 20), P(27, 2, 27))),
										T(L(S(P(24, 2, 24), P(25, 2, 25))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(20, 2, 20), P(22, 2, 22))), "foo"),
										ast.NewIntLiteralNode(L(S(P(27, 2, 27), P(27, 2, 27))), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(29, 3, 1), P(32, 3, 4))),
						ast.NewNilLiteralNode(L(S(P(29, 3, 1), P(31, 3, 3)))),
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
				L(S(P(0, 1, 1), P(57, 9, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(53, 8, 4))),
						ast.NewUntilExpressionNode(
							L(S(P(1, 2, 1), P(52, 8, 3))),
							ast.NewBinaryExpressionNode(
								L(S(P(7, 2, 7), P(13, 2, 13))),
								T(L(S(P(11, 2, 11), P(11, 2, 11))), token.GREATER),
								ast.NewPublicIdentifierNode(L(S(P(7, 2, 7), P(9, 2, 9))), "foo"),
								ast.NewIntLiteralNode(L(S(P(13, 2, 13), P(13, 2, 13))), "0"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(16, 3, 2), P(24, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(16, 3, 2), P(23, 3, 9))),
										T(L(S(P(20, 3, 6), P(21, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(16, 3, 2), P(18, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(23, 3, 9), P(23, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(26, 4, 2), P(29, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(26, 4, 2), P(28, 4, 4)))),
								),
								ast.NewExpressionStatementNode(
									L(S(P(30, 5, 1), P(34, 5, 5))),
									ast.NewInvalidNode(L(S(P(30, 5, 1), P(33, 5, 4))), T(L(S(P(30, 5, 1), P(33, 5, 4))), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									L(S(P(36, 6, 2), P(44, 6, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(36, 6, 2), P(43, 6, 9))),
										T(L(S(P(40, 6, 6), P(41, 6, 7))), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(36, 6, 2), P(38, 6, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(43, 6, 9), P(43, 6, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(46, 7, 2), P(49, 7, 5))),
									ast.NewNilLiteralNode(L(S(P(46, 7, 2), P(48, 7, 4)))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(54, 9, 1), P(57, 9, 4))),
						ast.NewNilLiteralNode(L(S(P(54, 9, 1), P(56, 9, 3)))),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(30, 5, 1), P(33, 5, 4))), "unexpected else, expected an expression"),
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
				L(S(P(0, 1, 1), P(24, 5, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(24, 5, 4))),
						ast.NewLoopExpressionNode(
							L(S(P(1, 2, 1), P(23, 5, 3))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(7, 3, 2), P(15, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(7, 3, 2), P(14, 3, 9))),
										T(L(S(P(11, 3, 6), P(12, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(7, 3, 2), P(9, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(14, 3, 9), P(14, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(17, 4, 2), P(20, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(17, 4, 2), P(19, 4, 4)))),
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
				L(S(P(0, 1, 1), P(9, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(9, 3, 4))),
						ast.NewLoopExpressionNode(
							L(S(P(1, 2, 1), P(8, 3, 3))),
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
				L(S(P(0, 1, 1), P(32, 6, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(28, 5, 5))),
						ast.NewAssignmentExpressionNode(
							L(S(P(1, 2, 1), P(27, 5, 4))),
							T(L(S(P(5, 2, 5), P(5, 2, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(1, 2, 1), P(3, 2, 3))), "bar"),
							ast.NewLoopExpressionNode(
								L(S(P(8, 3, 2), P(27, 5, 4))),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										L(S(P(15, 4, 3), P(23, 4, 11))),
										ast.NewAssignmentExpressionNode(
											L(S(P(15, 4, 3), P(22, 4, 10))),
											T(L(S(P(19, 4, 7), P(20, 4, 8))), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(L(S(P(15, 4, 3), P(17, 4, 5))), "foo"),
											ast.NewIntLiteralNode(L(S(P(22, 4, 10), P(22, 4, 10))), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(29, 6, 1), P(32, 6, 4))),
						ast.NewNilLiteralNode(L(S(P(29, 6, 1), P(31, 6, 3)))),
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
				L(S(P(0, 1, 1), P(18, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(14, 2, 14))),
						ast.NewLoopExpressionNode(
							L(S(P(1, 2, 1), P(13, 2, 13))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(6, 2, 6), P(13, 2, 13))),
									ast.NewAssignmentExpressionNode(
										L(S(P(6, 2, 6), P(13, 2, 13))),
										T(L(S(P(10, 2, 10), P(11, 2, 11))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(6, 2, 6), P(8, 2, 8))), "foo"),
										ast.NewIntLiteralNode(L(S(P(13, 2, 13), P(13, 2, 13))), "2"),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(15, 3, 1), P(18, 3, 4))),
						ast.NewNilLiteralNode(L(S(P(15, 3, 1), P(17, 3, 3)))),
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
				L(S(P(0, 1, 1), P(48, 9, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(44, 8, 4))),
						ast.NewLoopExpressionNode(
							L(S(P(1, 2, 1), P(43, 8, 3))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(7, 3, 2), P(15, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(7, 3, 2), P(14, 3, 9))),
										T(L(S(P(11, 3, 6), P(12, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(7, 3, 2), P(9, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(14, 3, 9), P(14, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(17, 4, 2), P(20, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(17, 4, 2), P(19, 4, 4)))),
								),
								ast.NewExpressionStatementNode(
									L(S(P(21, 5, 1), P(25, 5, 5))),
									ast.NewInvalidNode(L(S(P(21, 5, 1), P(24, 5, 4))), T(L(S(P(21, 5, 1), P(24, 5, 4))), token.ELSE)),
								),
								ast.NewExpressionStatementNode(
									L(S(P(27, 6, 2), P(35, 6, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(27, 6, 2), P(34, 6, 9))),
										T(L(S(P(31, 6, 6), P(32, 6, 7))), token.MINUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(27, 6, 2), P(29, 6, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(34, 6, 9), P(34, 6, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(37, 7, 2), P(40, 7, 5))),
									ast.NewNilLiteralNode(L(S(P(37, 7, 2), P(39, 7, 4)))),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(45, 9, 1), P(48, 9, 4))),
						ast.NewNilLiteralNode(L(S(P(45, 9, 1), P(47, 9, 3)))),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(21, 5, 1), P(24, 5, 4))), "unexpected else, expected an expression"),
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
				L(S(P(0, 1, 1), P(22, 5, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(22, 5, 4))),
						ast.NewDoExpressionNode(
							L(S(P(1, 2, 1), P(21, 5, 3))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(5, 3, 2), P(13, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(5, 3, 2), P(12, 3, 9))),
										T(L(S(P(9, 3, 6), P(10, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(5, 3, 2), P(7, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(12, 3, 9), P(12, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(15, 4, 2), P(18, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(15, 4, 2), P(17, 4, 4)))),
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have finally": {
			input: `
do
	foo += 2
	nil
finally
  println("foo")
end
`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(47, 7, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(47, 7, 4))),
						ast.NewDoExpressionNode(
							L(S(P(1, 2, 1), P(46, 7, 3))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(5, 3, 2), P(13, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(5, 3, 2), P(12, 3, 9))),
										T(L(S(P(9, 3, 6), P(10, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(5, 3, 2), P(7, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(12, 3, 9), P(12, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(15, 4, 2), P(18, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(15, 4, 2), P(17, 4, 4)))),
								),
							},
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(29, 6, 3), P(43, 6, 17))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(29, 6, 3), P(42, 6, 16))),
										ast.NewPublicIdentifierNode(L(S(P(29, 6, 3), P(35, 6, 9))), "println"),
										[]ast.ExpressionNode{
											ast.NewDoubleQuotedStringLiteralNode(L(S(P(37, 6, 11), P(41, 6, 15))), "foo"),
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
		"can have catch": {
			input: `
do
	foo += 2
	nil
catch Error() as e
	println(e)
catch Symbol() as s, stack_trace
	println(s)
end
`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(98, 9, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(98, 9, 4))),
						ast.NewDoExpressionNode(
							L(S(P(1, 2, 1), P(97, 9, 3))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(5, 3, 2), P(13, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(5, 3, 2), P(12, 3, 9))),
										T(L(S(P(9, 3, 6), P(10, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(5, 3, 2), P(7, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(12, 3, 9), P(12, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(15, 4, 2), P(18, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(15, 4, 2), P(17, 4, 4)))),
								),
							},
							[]*ast.CatchNode{
								ast.NewCatchNode(
									L(S(P(19, 5, 1), P(49, 6, 12))),
									ast.NewAsPatternNode(
										L(S(P(25, 5, 7), P(36, 5, 18))),
										ast.NewObjectPatternNode(
											L(S(P(25, 5, 7), P(31, 5, 13))),
											ast.NewPublicConstantNode(L(S(P(25, 5, 7), P(29, 5, 11))), "Error"),
											nil,
										),
										ast.NewPublicIdentifierNode(L(S(P(36, 5, 18), P(36, 5, 18))), "e"),
									),
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(39, 6, 2), P(49, 6, 12))),
											ast.NewReceiverlessMethodCallNode(
												L(S(P(39, 6, 2), P(48, 6, 11))),
												ast.NewPublicIdentifierNode(L(S(P(39, 6, 2), P(45, 6, 8))), "println"),
												[]ast.ExpressionNode{
													ast.NewPublicIdentifierNode(L(S(P(47, 6, 10), P(47, 6, 10))), "e"),
												},
												nil,
											),
										),
									},
								),
								ast.NewCatchNode(
									L(S(P(50, 7, 1), P(94, 8, 12))),
									ast.NewAsPatternNode(
										L(S(P(56, 7, 7), P(68, 7, 19))),
										ast.NewObjectPatternNode(
											L(S(P(56, 7, 7), P(63, 7, 14))),
											ast.NewPublicConstantNode(L(S(P(56, 7, 7), P(61, 7, 12))), "Symbol"),
											nil,
										),
										ast.NewPublicIdentifierNode(L(S(P(68, 7, 19), P(68, 7, 19))), "s"),
									),
									ast.NewPublicIdentifierNode(L(S(P(71, 7, 22), P(81, 7, 32))), "stack_trace"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(84, 8, 2), P(94, 8, 12))),
											ast.NewReceiverlessMethodCallNode(
												L(S(P(84, 8, 2), P(93, 8, 11))),
												ast.NewPublicIdentifierNode(L(S(P(84, 8, 2), P(90, 8, 8))), "println"),
												[]ast.ExpressionNode{
													ast.NewPublicIdentifierNode(L(S(P(92, 8, 10), P(92, 8, 10))), "s"),
												},
												nil,
											),
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
		"can have an empty body": {
			input: `
do
end
`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(7, 3, 4))),
						ast.NewDoExpressionNode(
							L(S(P(1, 2, 1), P(6, 3, 3))),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be a one-liner": {
			input: "do 5",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(3, 1, 4))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(3, 1, 4))),
						ast.NewDoExpressionNode(
							L(S(P(0, 1, 1), P(3, 1, 4))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(3, 1, 4), P(3, 1, 4))),
									ast.NewIntLiteralNode(
										L(S(P(3, 1, 4), P(3, 1, 4))),
										"5",
									),
								),
							},
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
	do
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(30, 6, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(26, 5, 5))),
						ast.NewAssignmentExpressionNode(
							L(S(P(1, 2, 1), P(25, 5, 4))),
							T(L(S(P(5, 2, 5), P(5, 2, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(1, 2, 1), P(3, 2, 3))), "bar"),
							ast.NewDoExpressionNode(
								L(S(P(8, 3, 2), P(25, 5, 4))),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										L(S(P(13, 4, 3), P(21, 4, 11))),
										ast.NewAssignmentExpressionNode(
											L(S(P(13, 4, 3), P(20, 4, 10))),
											T(L(S(P(17, 4, 7), P(18, 4, 8))), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(L(S(P(13, 4, 3), P(15, 4, 5))), "foo"),
											ast.NewIntLiteralNode(L(S(P(20, 4, 10), P(20, 4, 10))), "2"),
										),
									),
								},
								nil,
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(27, 6, 1), P(30, 6, 4))),
						ast.NewNilLiteralNode(L(S(P(27, 6, 1), P(29, 6, 3)))),
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
				L(S(P(0, 1, 1), P(4, 1, 5))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(4, 1, 5))),
						ast.NewBreakExpressionNode(
							L(S(P(0, 1, 1), P(4, 1, 5))),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a label": {
			input: `break[foo]`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewBreakExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							nil,
						),
					),
				},
			),
		},
		"can have a modifier if without an argument": {
			input: `break if true`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 1, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(12, 1, 13))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(12, 1, 13))),
							T(L(S(P(6, 1, 7), P(7, 1, 8))), token.IF),
							ast.NewBreakExpressionNode(
								L(S(P(0, 1, 1), P(4, 1, 5))),
								nil,
								nil,
							),
							ast.NewTrueLiteralNode(L(S(P(9, 1, 10), P(12, 1, 13)))),
						),
					),
				},
			),
		},
		"can have a modifier if with an argument": {
			input: `break :foo if true`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(17, 1, 18))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(17, 1, 18))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(17, 1, 18))),
							T(L(S(P(11, 1, 12), P(12, 1, 13))), token.IF),
							ast.NewBreakExpressionNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								nil,
								ast.NewSimpleSymbolLiteralNode(
									L(S(P(6, 1, 7), P(9, 1, 10))),
									"foo",
								),
							),
							ast.NewTrueLiteralNode(L(S(P(14, 1, 15), P(17, 1, 18)))),
						),
					),
				},
			),
		},
		"can have an argument": {
			input: `break 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewBreakExpressionNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							nil,
							ast.NewIntLiteralNode(
								L(S(P(6, 1, 7), P(6, 1, 7))),
								"2",
							),
						),
					),
				},
			),
		},
		"can have a label and argument": {
			input: `break[foo] 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewBreakExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							ast.NewIntLiteralNode(
								L(S(P(11, 1, 12), P(11, 1, 12))),
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
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewBreakExpressionNode(L(S(P(7, 1, 8), P(11, 1, 12))), nil, nil),
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

func TestAwaitSync(t *testing.T) {
	tests := testTable{
		"cannot stand alone": {
			input: `await_sync`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewInvalidNode(
								L(S(P(10, 1, 11), P(9, 1, 10))),
								T(L(S(P(10, 1, 11), P(9, 1, 10))), token.END_OF_FILE),
							),
							true,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(10, 1, 11), P(9, 1, 10))), "unexpected END_OF_FILE, expected an expression"),
			},
		},
		"can have an argument": {
			input: `await_sync 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewIntLiteralNode(L(S(P(11, 1, 12), P(11, 1, 12))), "2"),
							true,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && await_sync foo`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(20, 1, 21))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(20, 1, 21))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(20, 1, 21))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewAwaitExpressionNode(
								L(S(P(7, 1, 8), P(20, 1, 21))),
								ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(20, 1, 21))), "foo"),
								true,
							),
						),
					),
				},
			),
		},
		"can resemble a method call": {
			input: `foo.await_sync`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(13, 1, 14))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(13, 1, 14))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(13, 1, 14))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							true,
						),
					),
				},
			),
		},
		"cannot have arguments": {
			input: `foo.await(2)`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							false,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(9, 1, 10), P(9, 1, 10))), "unexpected (, expected a statement separator `\\n`, `;`"),
			},
		},
		"can be chained": {
			input: `foo.await_sync.elo().await_sync`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(30, 1, 31))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(30, 1, 31))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(30, 1, 31))),
							ast.NewMethodCallNode(
								L(S(P(0, 1, 1), P(19, 1, 20))),
								ast.NewAwaitExpressionNode(
									L(S(P(0, 1, 1), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
									true,
								),
								T(L(S(P(14, 1, 15), P(14, 1, 15))), token.DOT),
								ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(17, 1, 18))), "elo"),
								nil,
								nil,
							),
							true,
						),
					),
				},
			),
		},
		"invalid operator": {
			input: `foo..await_sync`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							true,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(3, 1, 4), P(4, 1, 5))), "invalid await_sync operator"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestAwait(t *testing.T) {
	tests := testTable{
		"cannot stand alone": {
			input: `await`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(4, 1, 5))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(4, 1, 5))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(4, 1, 5))),
							ast.NewInvalidNode(
								L(S(P(5, 1, 6), P(4, 1, 5))),
								T(L(S(P(5, 1, 6), P(4, 1, 5))), token.END_OF_FILE),
							),
							false,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(5, 1, 6), P(4, 1, 5))), "unexpected END_OF_FILE, expected an expression"),
			},
		},
		"can have an argument": {
			input: `await 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							ast.NewIntLiteralNode(L(S(P(6, 1, 7), P(6, 1, 7))), "2"),
							false,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && await foo`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewAwaitExpressionNode(
								L(S(P(7, 1, 8), P(15, 1, 16))),
								ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(15, 1, 16))), "foo"),
								false,
							),
						),
					),
				},
			),
		},
		"can resemble a method call": {
			input: `foo.await`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							false,
						),
					),
				},
			),
		},
		"cannot have arguments": {
			input: `foo.await(2)`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							false,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(9, 1, 10), P(9, 1, 10))), "unexpected (, expected a statement separator `\\n`, `;`"),
			},
		},
		"can be chained": {
			input: `foo.await.elo().await`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(20, 1, 21))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(20, 1, 21))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(20, 1, 21))),
							ast.NewMethodCallNode(
								L(S(P(0, 1, 1), P(14, 1, 15))),
								ast.NewAwaitExpressionNode(
									L(S(P(0, 1, 1), P(8, 1, 9))),
									ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
									false,
								),
								T(L(S(P(9, 1, 10), P(9, 1, 10))), token.DOT),
								ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(12, 1, 13))), "elo"),
								nil,
								nil,
							),
							false,
						),
					),
				},
			),
		},
		"invalid operator": {
			input: `foo..await`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewAwaitExpressionNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							false,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(3, 1, 4), P(4, 1, 5))), "invalid await operator"),
			},
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
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewReturnExpressionNode(L(S(P(0, 1, 1), P(5, 1, 6))), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "return\n1",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 2, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewReturnExpressionNode(L(S(P(0, 1, 1), P(5, 1, 6))), nil),
					),
					ast.NewExpressionStatementNode(
						L(S(P(7, 2, 1), P(7, 2, 1))),
						ast.NewIntLiteralNode(L(S(P(7, 2, 1), P(7, 2, 1))), "1"),
					),
				},
			),
		},
		"can have a modifier if without an argument": {
			input: `return if true`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(13, 1, 14))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(13, 1, 14))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(13, 1, 14))),
							T(L(S(P(7, 1, 8), P(8, 1, 9))), token.IF),
							ast.NewReturnExpressionNode(L(S(P(0, 1, 1), P(5, 1, 6))), nil),
							ast.NewTrueLiteralNode(L(S(P(10, 1, 11), P(13, 1, 14)))),
						),
					),
				},
			),
		},
		"can have a modifier if with an argument": {
			input: `return :foo if true`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(18, 1, 19))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(18, 1, 19))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(18, 1, 19))),
							T(L(S(P(12, 1, 13), P(13, 1, 14))), token.IF),
							ast.NewReturnExpressionNode(
								L(S(P(0, 1, 1), P(10, 1, 11))),
								ast.NewSimpleSymbolLiteralNode(
									L(S(P(7, 1, 8), P(10, 1, 11))),
									"foo",
								),
							),
							ast.NewTrueLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19)))),
						),
					),
				},
			),
		},
		"can have an argument": {
			input: `return 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewReturnExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewIntLiteralNode(L(S(P(7, 1, 8), P(7, 1, 8))), "2"),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && return`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 1, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(12, 1, 13))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(12, 1, 13))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewReturnExpressionNode(L(S(P(7, 1, 8), P(12, 1, 13))), nil),
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

func TestYield(t *testing.T) {
	tests := testTable{
		"can stand alone at the end": {
			input: `yield`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(4, 1, 5))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(4, 1, 5))),
						ast.NewYieldExpressionNode(L(S(P(0, 1, 1), P(4, 1, 5))), false, nil),
					),
				},
			),
		},
		"cannot stand alone with forwarding": {
			input: `yield*`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewYieldExpressionNode(
							L(S(P(0, 1, 1), P(5, 1, 6))),
							true,
							ast.NewInvalidNode(
								L(S(P(6, 1, 7), P(5, 1, 6))),
								T(L(S(P(6, 1, 7), P(5, 1, 6))), token.END_OF_FILE),
							),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 1, 7), P(5, 1, 6))), "unexpected END_OF_FILE, expected an expression"),
			},
		},
		"can stand alone in the middle": {
			input: "yield\n1",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 2, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewYieldExpressionNode(L(S(P(0, 1, 1), P(4, 1, 5))), false, nil),
					),
					ast.NewExpressionStatementNode(
						L(S(P(6, 2, 1), P(6, 2, 1))),
						ast.NewIntLiteralNode(L(S(P(6, 2, 1), P(6, 2, 1))), "1"),
					),
				},
			),
		},
		"can have a modifier if without an argument": {
			input: `yield if true`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 1, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(12, 1, 13))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(12, 1, 13))),
							T(L(S(P(6, 1, 7), P(7, 1, 8))), token.IF),
							ast.NewYieldExpressionNode(L(S(P(0, 1, 1), P(4, 1, 5))), false, nil),
							ast.NewTrueLiteralNode(L(S(P(9, 1, 10), P(12, 1, 13)))),
						),
					),
				},
			),
		},
		"can have a modifier if with an argument": {
			input: `yield :foo if true`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(17, 1, 18))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(17, 1, 18))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(17, 1, 18))),
							T(L(S(P(11, 1, 12), P(12, 1, 13))), token.IF),
							ast.NewYieldExpressionNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								false,
								ast.NewSimpleSymbolLiteralNode(
									L(S(P(6, 1, 7), P(9, 1, 10))),
									"foo",
								),
							),
							ast.NewTrueLiteralNode(L(S(P(14, 1, 15), P(17, 1, 18)))),
						),
					),
				},
			),
		},
		"can have an argument with forwarding": {
			input: `yield* 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewYieldExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							true,
							ast.NewIntLiteralNode(L(S(P(7, 1, 8), P(7, 1, 8))), "2"),
						),
					),
				},
			),
		},
		"can have an argument": {
			input: `yield 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewYieldExpressionNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							false,
							ast.NewIntLiteralNode(L(S(P(6, 1, 7), P(6, 1, 7))), "2"),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && yield`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewYieldExpressionNode(L(S(P(7, 1, 8), P(11, 1, 12))), false, nil),
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
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewContinueExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "continue\n1",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 2, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewContinueExpressionNode(L(S(P(0, 1, 1), P(7, 1, 8))), nil, nil),
					),
					ast.NewExpressionStatementNode(
						L(S(P(9, 2, 1), P(9, 2, 1))),
						ast.NewIntLiteralNode(L(S(P(9, 2, 1), P(9, 2, 1))), "1"),
					),
				},
			),
		},
		"can have a label": {
			input: `continue[foo]`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 1, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(12, 1, 13))),
						ast.NewContinueExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewPublicIdentifierNode(L(S(P(9, 1, 10), P(11, 1, 12))), "foo"),
							nil,
						),
					),
				},
			),
		},
		"can have an argument": {
			input: `continue 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewContinueExpressionNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							nil,
							ast.NewIntLiteralNode(L(S(P(9, 1, 10), P(9, 1, 10))), "2"),
						),
					),
				},
			),
		},
		"can have a modifier if without an argument": {
			input: `continue if true`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							T(L(S(P(9, 1, 10), P(10, 1, 11))), token.IF),
							ast.NewContinueExpressionNode(L(S(P(0, 1, 1), P(7, 1, 8))), nil, nil),
							ast.NewTrueLiteralNode(L(S(P(12, 1, 13), P(15, 1, 16)))),
						),
					),
				},
			),
		},
		"can have a modifier if with an argument": {
			input: `continue :foo if true`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(20, 1, 21))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(20, 1, 21))),
						ast.NewModifierNode(
							L(S(P(0, 1, 1), P(20, 1, 21))),
							T(L(S(P(14, 1, 15), P(15, 1, 16))), token.IF),
							ast.NewContinueExpressionNode(
								L(S(P(0, 1, 1), P(12, 1, 13))),
								nil,
								ast.NewSimpleSymbolLiteralNode(
									L(S(P(9, 1, 10), P(12, 1, 13))),
									"foo",
								),
							),
							ast.NewTrueLiteralNode(L(S(P(17, 1, 18), P(20, 1, 21)))),
						),
					),
				},
			),
		},
		"can have a label and argument": {
			input: `continue[foo] 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewContinueExpressionNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							ast.NewPublicIdentifierNode(L(S(P(9, 1, 10), P(11, 1, 12))), "foo"),
							ast.NewIntLiteralNode(
								L(S(P(14, 1, 15), P(14, 1, 15))),
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
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewContinueExpressionNode(L(S(P(7, 1, 8), P(14, 1, 15))), nil, nil),
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

func TestTry(t *testing.T) {
	tests := testTable{
		"can have an argument": {
			input: `try 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(4, 1, 5))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(4, 1, 5))),
						ast.NewTryExpressionNode(
							L(S(P(0, 1, 1), P(4, 1, 5))),
							ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "2"),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && try bar()`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewTryExpressionNode(
								L(S(P(7, 1, 8), P(15, 1, 16))),
								ast.NewReceiverlessMethodCallNode(
									L(S(P(11, 1, 12), P(15, 1, 16))),
									ast.NewPublicIdentifierNode(L(S(P(11, 1, 12), P(13, 1, 14))), "bar"),
									nil,
									nil,
								),
							),
						),
					),
				},
			),
		},
		"chained": {
			input: `foo && bar().try`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewTryExpressionNode(
								L(S(P(7, 1, 8), P(15, 1, 16))),
								ast.NewReceiverlessMethodCallNode(
									L(S(P(7, 1, 8), P(11, 1, 12))),
									ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "bar"),
									nil,
									nil,
								),
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

func TestMust(t *testing.T) {
	tests := testTable{
		"can have an argument": {
			input: `must 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewMustExpressionNode(
							L(S(P(0, 1, 1), P(5, 1, 6))),
							ast.NewIntLiteralNode(L(S(P(5, 1, 6), P(5, 1, 6))), "2"),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && must bar()`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewMustExpressionNode(
								L(S(P(7, 1, 8), P(16, 1, 17))),
								ast.NewReceiverlessMethodCallNode(
									L(S(P(12, 1, 13), P(16, 1, 17))),
									ast.NewPublicIdentifierNode(L(S(P(12, 1, 13), P(14, 1, 15))), "bar"),
									nil,
									nil,
								),
							),
						),
					),
				},
			),
		},
		"can resemble a method call": {
			input: `foo.must`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewMustExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
						),
					),
				},
			),
		},
		"cannot have arguments": {
			input: `foo.must(2)`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewMustExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(8, 1, 9), P(8, 1, 9))), "unexpected (, expected a statement separator `\\n`, `;`"),
			},
		},
		"can be chained": {
			input: `foo.must.elo().must`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(18, 1, 19))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(18, 1, 19))),
						ast.NewMustExpressionNode(
							L(S(P(0, 1, 1), P(18, 1, 19))),
							ast.NewMethodCallNode(
								L(S(P(0, 1, 1), P(13, 1, 14))),
								ast.NewMustExpressionNode(
									L(S(P(0, 1, 1), P(7, 1, 8))),
									ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								),
								T(L(S(P(8, 1, 9), P(8, 1, 9))), token.DOT),
								ast.NewPublicIdentifierNode(L(S(P(9, 1, 10), P(11, 1, 12))), "elo"),
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"invalid operator": {
			input: `foo..must`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewMustExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(3, 1, 4), P(4, 1, 5))), "invalid must operator"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestAs(t *testing.T) {
	tests := testTable{
		"can have a public constant as a type": {
			input: `a as String`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewAsExpressionNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(0, 1, 1))), "a"),
							ast.NewPublicConstantNode(L(S(P(5, 1, 6), P(10, 1, 11))), "String"),
						),
					),
				},
			),
		},
		"can have a private constant as a type": {
			input: `a as _String`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewAsExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(0, 1, 1))), "a"),
							ast.NewPrivateConstantNode(L(S(P(5, 1, 6), P(11, 1, 12))), "_String"),
						),
					),
				},
			),
		},
		"can have a constant lookup as a type": {
			input: `a as ::Std::String`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(17, 1, 18))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(17, 1, 18))),
						ast.NewAsExpressionNode(
							L(S(P(0, 1, 1), P(17, 1, 18))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(0, 1, 1))), "a"),
							ast.NewConstantLookupNode(
								L(S(P(5, 1, 6), P(17, 1, 18))),
								ast.NewConstantLookupNode(
									L(S(P(5, 1, 6), P(9, 1, 10))),
									nil,
									ast.NewPublicConstantNode(L(S(P(7, 1, 8), P(9, 1, 10))), "Std"),
								),
								ast.NewPublicConstantNode(L(S(P(12, 1, 13), P(17, 1, 18))), "String"),
							),
						),
					),
				},
			),
		},
		"cannot have a public identifier as a type": {
			input: `a as string`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(0, 1, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(0, 1, 1))),
						ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(0, 1, 1))), "a"),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(2, 1, 3), P(3, 1, 4))), "unexpected as, expected a statement separator `\\n`, `;`"),
			},
		},
		"is an expression": {
			input: `foo := a as String`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(17, 1, 18))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(17, 1, 18))),
						ast.NewAssignmentExpressionNode(
							L(S(P(0, 1, 1), P(17, 1, 18))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.COLON_EQUAL),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewAsExpressionNode(
								L(S(P(7, 1, 8), P(17, 1, 18))),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(7, 1, 8))), "a"),
								ast.NewPublicConstantNode(L(S(P(12, 1, 13), P(17, 1, 18))), "String"),
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

func TestThrow(t *testing.T) {
	tests := testTable{
		"can stand alone at the end": {
			input: `throw`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(4, 1, 5))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(4, 1, 5))),
						ast.NewThrowExpressionNode(L(S(P(0, 1, 1), P(4, 1, 5))), false, nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "throw\n1",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 2, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewThrowExpressionNode(L(S(P(0, 1, 1), P(4, 1, 5))), false, nil),
					),
					ast.NewExpressionStatementNode(
						L(S(P(6, 2, 1), P(6, 2, 1))),
						ast.NewIntLiteralNode(L(S(P(6, 2, 1), P(6, 2, 1))), "1"),
					),
				},
			),
		},
		"can have an argument": {
			input: `throw 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewThrowExpressionNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							false,
							ast.NewIntLiteralNode(L(S(P(6, 1, 7), P(6, 1, 7))), "2"),
						),
					),
				},
			),
		},
		"can be unchecked": {
			input: `throw unchecked 2`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewThrowExpressionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							true,
							ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "2"),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && throw`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewLogicalExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewThrowExpressionNode(L(S(P(7, 1, 8), P(11, 1, 12))), false, nil),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewForInExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(4, 1, 5))), "i"),
							ast.NewArrayListLiteralNode(
								L(S(P(9, 1, 10), P(17, 1, 18))),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(L(S(P(10, 1, 11), P(10, 1, 11))), "1"),
									ast.NewIntLiteralNode(L(S(P(13, 1, 14), P(13, 1, 14))), "2"),
									ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "3"),
								},
								nil,
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(24, 1, 25), P(33, 1, 34))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(24, 1, 25), P(33, 1, 34))),
										ast.NewPublicIdentifierNode(L(S(P(24, 1, 25), P(30, 1, 31))), "println"),
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(L(S(P(32, 1, 33), P(32, 1, 33))), "i"),
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
		"can have patterns": {
			input: `for [a, b] in [[1, 2], [3, 4]] then println(a, b)`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(48, 1, 49))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(48, 1, 49))),
						ast.NewForInExpressionNode(
							L(S(P(0, 1, 1), P(48, 1, 49))),
							ast.NewListPatternNode(
								L(S(P(4, 1, 5), P(9, 1, 10))),
								[]ast.PatternNode{
									ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(5, 1, 6))), "a"),
									ast.NewPublicIdentifierNode(L(S(P(8, 1, 9), P(8, 1, 9))), "b"),
								},
							),
							ast.NewArrayListLiteralNode(
								L(S(P(14, 1, 15), P(29, 1, 30))),
								[]ast.ExpressionNode{
									ast.NewArrayListLiteralNode(
										L(S(P(15, 1, 16), P(20, 1, 21))),
										[]ast.ExpressionNode{
											ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "1"),
											ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "2"),
										},
										nil,
									),
									ast.NewArrayListLiteralNode(
										L(S(P(23, 1, 24), P(28, 1, 29))),
										[]ast.ExpressionNode{
											ast.NewIntLiteralNode(L(S(P(24, 1, 25), P(24, 1, 25))), "3"),
											ast.NewIntLiteralNode(L(S(P(27, 1, 28), P(27, 1, 28))), "4"),
										},
										nil,
									),
								},
								nil,
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(36, 1, 37), P(48, 1, 49))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(36, 1, 37), P(48, 1, 49))),
										ast.NewPublicIdentifierNode(L(S(P(36, 1, 37), P(42, 1, 43))), "println"),
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(L(S(P(44, 1, 45), P(44, 1, 45))), "a"),
											ast.NewPublicIdentifierNode(L(S(P(47, 1, 48), P(47, 1, 48))), "b"),
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
		"cannot have patterns without variables": {
			input: `for [1, 2] in [[1, 2], [3, 4]] then println(a, b)`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(48, 1, 49))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(48, 1, 49))),
						ast.NewForInExpressionNode(
							L(S(P(0, 1, 1), P(48, 1, 49))),
							ast.NewListPatternNode(
								L(S(P(4, 1, 5), P(9, 1, 10))),
								[]ast.PatternNode{
									ast.NewIntLiteralNode(L(S(P(5, 1, 6), P(5, 1, 6))), "1"),
									ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "2"),
								},
							),
							ast.NewArrayListLiteralNode(
								L(S(P(14, 1, 15), P(29, 1, 30))),
								[]ast.ExpressionNode{
									ast.NewArrayListLiteralNode(
										L(S(P(15, 1, 16), P(20, 1, 21))),
										[]ast.ExpressionNode{
											ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "1"),
											ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "2"),
										},
										nil,
									),
									ast.NewArrayListLiteralNode(
										L(S(P(23, 1, 24), P(28, 1, 29))),
										[]ast.ExpressionNode{
											ast.NewIntLiteralNode(L(S(P(24, 1, 25), P(24, 1, 25))), "3"),
											ast.NewIntLiteralNode(L(S(P(27, 1, 28), P(27, 1, 28))), "4"),
										},
										nil,
									),
								},
								nil,
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(36, 1, 37), P(48, 1, 49))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(36, 1, 37), P(48, 1, 49))),
										ast.NewPublicIdentifierNode(L(S(P(36, 1, 37), P(42, 1, 43))), "println"),
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(L(S(P(44, 1, 45), P(44, 1, 45))), "a"),
											ast.NewPublicIdentifierNode(L(S(P(47, 1, 48), P(47, 1, 48))), "b"),
										},
										nil,
									),
								),
							},
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(4, 1, 5), P(9, 1, 10))), "patterns in for in loops should define at least one variable"),
			},
		},
		"can be multiline": {
			input: `for i in [1, 2, 3]
  println(i)
  nil
end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(40, 4, 3))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(40, 4, 3))),
						ast.NewForInExpressionNode(
							L(S(P(0, 1, 1), P(40, 4, 3))),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(4, 1, 5))), "i"),
							ast.NewArrayListLiteralNode(
								L(S(P(9, 1, 10), P(17, 1, 18))),
								[]ast.ExpressionNode{
									ast.NewIntLiteralNode(L(S(P(10, 1, 11), P(10, 1, 11))), "1"),
									ast.NewIntLiteralNode(L(S(P(13, 1, 14), P(13, 1, 14))), "2"),
									ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "3"),
								},
								nil,
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(21, 2, 3), P(31, 2, 13))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(21, 2, 3), P(30, 2, 12))),
										ast.NewPublicIdentifierNode(L(S(P(21, 2, 3), P(27, 2, 9))), "println"),
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(L(S(P(29, 2, 11), P(29, 2, 11))), "i"),
										},
										nil,
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(34, 3, 3), P(37, 3, 6))),
									ast.NewNilLiteralNode(L(S(P(34, 3, 3), P(36, 3, 5)))),
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
				L(S(P(0, 1, 1), P(43, 1, 44))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(43, 1, 44))),
						ast.NewNumericForExpressionNode(
							L(S(P(0, 1, 1), P(43, 1, 44))),
							ast.NewAssignmentExpressionNode(
								L(S(P(7, 1, 8), P(12, 1, 13))),
								T(L(S(P(9, 1, 10), P(10, 1, 11))), token.COLON_EQUAL),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(7, 1, 8))), "i"),
								ast.NewIntLiteralNode(L(S(P(12, 1, 13), P(12, 1, 13))), "0"),
							),
							ast.NewBinaryExpressionNode(
								L(S(P(15, 1, 16), P(19, 1, 20))),
								T(L(S(P(17, 1, 18), P(17, 1, 18))), token.LESS),
								ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(15, 1, 16))), "i"),
								ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "5"),
							),
							ast.NewAssignmentExpressionNode(
								L(S(P(22, 1, 23), P(27, 1, 28))),
								T(L(S(P(24, 1, 25), P(25, 1, 26))), token.PLUS_EQUAL),
								ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(22, 1, 23))), "i"),
								ast.NewIntLiteralNode(L(S(P(27, 1, 28), P(27, 1, 28))), "1"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(34, 1, 35), P(43, 1, 44))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(34, 1, 35), P(43, 1, 44))),
										ast.NewPublicIdentifierNode(L(S(P(34, 1, 35), P(40, 1, 41))), "println"),
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(L(S(P(42, 1, 43), P(42, 1, 43))), "i"),
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
				L(S(P(0, 1, 1), P(24, 1, 25))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(24, 1, 25))),
						ast.NewNumericForExpressionNode(
							L(S(P(0, 1, 1), P(24, 1, 25))),
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(15, 1, 16), P(24, 1, 25))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(15, 1, 16), P(24, 1, 25))),
										ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(21, 1, 22))), "println"),
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(L(S(P(23, 1, 24), P(23, 1, 24))), "i"),
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
				L(S(P(0, 1, 1), P(50, 4, 3))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(50, 4, 3))),
						ast.NewNumericForExpressionNode(
							L(S(P(0, 1, 1), P(50, 4, 3))),
							ast.NewAssignmentExpressionNode(
								L(S(P(7, 1, 8), P(12, 1, 13))),
								T(L(S(P(9, 1, 10), P(10, 1, 11))), token.COLON_EQUAL),
								ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(7, 1, 8))), "i"),
								ast.NewIntLiteralNode(L(S(P(12, 1, 13), P(12, 1, 13))), "0"),
							),
							ast.NewBinaryExpressionNode(
								L(S(P(15, 1, 16), P(19, 1, 20))),
								T(L(S(P(17, 1, 18), P(17, 1, 18))), token.LESS),
								ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(15, 1, 16))), "i"),
								ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "5"),
							),
							ast.NewAssignmentExpressionNode(
								L(S(P(22, 1, 23), P(27, 1, 28))),
								T(L(S(P(24, 1, 25), P(25, 1, 26))), token.PLUS_EQUAL),
								ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(22, 1, 23))), "i"),
								ast.NewIntLiteralNode(L(S(P(27, 1, 28), P(27, 1, 28))), "1"),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(31, 2, 3), P(41, 2, 13))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(31, 2, 3), P(40, 2, 12))),
										ast.NewPublicIdentifierNode(L(S(P(31, 2, 3), P(37, 2, 9))), "println"),
										[]ast.ExpressionNode{
											ast.NewPublicIdentifierNode(L(S(P(39, 2, 11), P(39, 2, 11))), "i"),
										},
										nil,
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(44, 3, 3), P(47, 3, 6))),
									ast.NewNilLiteralNode(L(S(P(44, 3, 3), P(46, 3, 5)))),
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
