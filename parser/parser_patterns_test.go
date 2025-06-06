package parser

import (
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/token"
)

func TestSwitch(t *testing.T) {
	tests := testTable{
		"cannot be empty": {
			input: `
switch foo
end
`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 3, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(15, 3, 4))),
						ast.NewSwitchExpressionNode(
							L(S(P(1, 2, 1), P(14, 3, 3))),
							ast.NewPublicIdentifierNode(L(S(P(8, 2, 8), P(10, 2, 10))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(1, 2, 1), P(14, 3, 3))), "switch cannot be empty"),
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
				L(S(P(0, 1, 1), P(43, 7, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(39, 6, 5))),
						ast.NewAssignmentExpressionNode(
							L(S(P(1, 2, 1), P(38, 6, 4))),
							T(L(S(P(5, 2, 5), P(5, 2, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(1, 2, 1), P(3, 2, 3))), "bar"),
							ast.NewSwitchExpressionNode(
								L(S(P(8, 3, 2), P(38, 6, 4))),
								ast.NewPublicIdentifierNode(L(S(P(15, 3, 9), P(17, 3, 11))), "foo"),
								[]*ast.CaseNode{
									ast.NewCaseNode(
										L(S(P(20, 4, 2), P(34, 5, 8))),
										ast.NewPublicIdentifierNode(L(S(P(25, 4, 7), P(25, 4, 7))), "n"),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												L(S(P(29, 5, 3), P(34, 5, 8))),
												ast.NewBinaryExpressionNode(
													L(S(P(29, 5, 3), P(33, 5, 7))),
													T(L(S(P(31, 5, 5), P(31, 5, 5))), token.PLUS),
													ast.NewPublicIdentifierNode(L(S(P(29, 5, 3), P(29, 5, 3))), "n"),
													ast.NewIntLiteralNode(L(S(P(33, 5, 7), P(33, 5, 7))), "2"),
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
						L(S(P(40, 7, 1), P(43, 7, 4))),
						ast.NewNilLiteralNode(L(S(P(40, 7, 1), P(42, 7, 3)))),
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
				L(S(P(0, 1, 1), P(28, 5, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(28, 5, 4))),
						ast.NewSwitchExpressionNode(
							L(S(P(1, 2, 1), P(27, 5, 3))),
							ast.NewPublicIdentifierNode(L(S(P(8, 2, 8), P(10, 2, 10))), "foo"),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(19, 4, 3), P(24, 4, 8))),
									ast.NewBinaryExpressionNode(
										L(S(P(19, 4, 3), P(23, 4, 7))),
										T(L(S(P(21, 4, 5), P(21, 4, 5))), token.PLUS),
										ast.NewPublicIdentifierNode(L(S(P(19, 4, 3), P(19, 4, 3))), "n"),
										ast.NewIntLiteralNode(L(S(P(23, 4, 7), P(23, 4, 7))), "2"),
									),
								),
							},
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(1, 2, 1), P(27, 5, 3))), "switch cannot only consist of else"),
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
				L(S(P(0, 1, 1), P(50, 9, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(50, 9, 4))),
						ast.NewSwitchExpressionNode(
							L(S(P(1, 2, 1), P(49, 9, 3))),
							ast.NewPublicIdentifierNode(L(S(P(8, 2, 8), P(10, 2, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(12, 3, 1), P(22, 4, 4))),
									ast.NewPublicIdentifierNode(L(S(P(17, 3, 6), P(17, 3, 6))), "n"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(21, 4, 3), P(22, 4, 4))),
											ast.NewPublicIdentifierNode(L(S(P(21, 4, 3), P(21, 4, 3))), "n"),
										),
									},
								),
								ast.NewCaseNode(
									L(S(P(23, 5, 1), P(33, 6, 4))),
									ast.NewPublicIdentifierNode(L(S(P(28, 5, 6), P(28, 5, 6))), "m"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(32, 6, 3), P(33, 6, 4))),
											ast.NewPublicIdentifierNode(L(S(P(32, 6, 3), P(32, 6, 3))), "m"),
										),
									},
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(41, 8, 3), P(46, 8, 8))),
									ast.NewBinaryExpressionNode(
										L(S(P(41, 8, 3), P(45, 8, 7))),
										T(L(S(P(43, 8, 5), P(43, 8, 5))), token.PLUS),
										ast.NewPublicIdentifierNode(L(S(P(41, 8, 3), P(41, 8, 3))), "n"),
										ast.NewIntLiteralNode(L(S(P(45, 8, 7), P(45, 8, 7))), "2"),
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
				L(S(P(0, 1, 1), P(54, 6, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(54, 6, 4))),
						ast.NewSwitchExpressionNode(
							L(S(P(1, 2, 1), P(53, 6, 3))),
							ast.NewPublicIdentifierNode(L(S(P(8, 2, 8), P(10, 2, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(12, 3, 1), P(24, 3, 13))),
									ast.NewPublicIdentifierNode(L(S(P(17, 3, 6), P(17, 3, 6))), "n"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(24, 3, 13), P(24, 3, 13))),
											ast.NewPublicIdentifierNode(L(S(P(24, 3, 13), P(24, 3, 13))), "n"),
										),
									},
								),
								ast.NewCaseNode(
									L(S(P(26, 4, 1), P(38, 4, 13))),
									ast.NewPublicIdentifierNode(L(S(P(31, 4, 6), P(31, 4, 6))), "m"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(38, 4, 13), P(38, 4, 13))),
											ast.NewPublicIdentifierNode(L(S(P(38, 4, 13), P(38, 4, 13))), "m"),
										),
									},
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(45, 5, 6), P(49, 5, 10))),
									ast.NewBinaryExpressionNode(
										L(S(P(45, 5, 6), P(49, 5, 10))),
										T(L(S(P(47, 5, 8), P(47, 5, 8))), token.PLUS),
										ast.NewPublicIdentifierNode(L(S(P(45, 5, 6), P(45, 5, 6))), "n"),
										ast.NewIntLiteralNode(L(S(P(49, 5, 10), P(49, 5, 10))), "2"),
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
				L(S(P(0, 1, 1), P(52, 1, 53))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(52, 1, 53))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(52, 1, 53))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(23, 1, 24))),
									ast.NewPublicIdentifierNode(L(S(P(16, 1, 17), P(16, 1, 17))), "n"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(23, 1, 24), P(23, 1, 24))),
											ast.NewPublicIdentifierNode(L(S(P(23, 1, 24), P(23, 1, 24))), "n"),
										),
									},
								),
								ast.NewCaseNode(
									L(S(P(25, 1, 26), P(37, 1, 38))),
									ast.NewPublicIdentifierNode(L(S(P(30, 1, 31), P(30, 1, 31))), "m"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(37, 1, 38), P(37, 1, 38))),
											ast.NewPublicIdentifierNode(L(S(P(37, 1, 38), P(37, 1, 38))), "m"),
										),
									},
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(44, 1, 45), P(48, 1, 49))),
									ast.NewBinaryExpressionNode(
										L(S(P(44, 1, 45), P(48, 1, 49))),
										T(L(S(P(46, 1, 47), P(46, 1, 47))), token.PLUS),
										ast.NewPublicIdentifierNode(L(S(P(44, 1, 45), P(44, 1, 45))), "n"),
										ast.NewIntLiteralNode(L(S(P(48, 1, 49), P(48, 1, 49))), "2"),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewTrueLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20)))),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(29, 1, 30))),
									ast.NewFalseLiteralNode(L(S(P(16, 1, 17), P(20, 1, 21)))),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(27, 1, 28), P(29, 1, 30))),
											ast.NewNilLiteralNode(L(S(P(27, 1, 28), P(29, 1, 30)))),
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
				L(S(P(0, 1, 1), P(31, 1, 32))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(31, 1, 32))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(31, 1, 32))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(27, 1, 28))),
									ast.NewNilLiteralNode(L(S(P(16, 1, 17), P(18, 1, 19)))),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(25, 1, 26), P(27, 1, 28))),
											ast.NewNilLiteralNode(L(S(P(25, 1, 26), P(27, 1, 28)))),
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
				L(S(P(0, 1, 1), P(31, 1, 32))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(31, 1, 32))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(31, 1, 32))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(27, 1, 28))),
									ast.NewCharLiteralNode(L(S(P(16, 1, 17), P(18, 1, 19))), 'f'),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(25, 1, 26), P(27, 1, 28))),
											ast.NewNilLiteralNode(L(S(P(25, 1, 26), P(27, 1, 28)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewRawCharLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), 'f'),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewRawStringLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "fo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewDoubleQuotedStringLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "fo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewUninterpolatedRegexLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "f", bitfield.BitField8{}),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewSimpleSymbolLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "foo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewSimpleSymbolLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "&"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewPublicIdentifierNode(L(S(P(16, 1, 17), P(19, 1, 20))), "foof"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewPrivateIdentifierNode(L(S(P(16, 1, 17), P(19, 1, 20))), "_foo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "1234"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewInt64LiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewUInt64LiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewInt32LiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewUInt32LiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewInt16LiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewUInt16LiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewInt8LiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "12"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewUInt8LiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "12"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewFloatLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "12.5"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewBigFloatLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "12"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewFloat64LiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewFloat32LiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "1"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(29, 1, 30))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(20, 1, 21))),
										T(L(S(P(16, 1, 17), P(16, 1, 17))), token.MINUS),
										ast.NewFloat32LiteralNode(L(S(P(17, 1, 18), P(20, 1, 21))), "1"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(27, 1, 28), P(29, 1, 30))),
											ast.NewNilLiteralNode(L(S(P(27, 1, 28), P(29, 1, 30)))),
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
				L(S(P(0, 1, 1), P(31, 1, 32))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(31, 1, 32))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(31, 1, 32))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(27, 1, 28))),
									ast.NewPublicConstantNode(L(S(P(16, 1, 17), P(18, 1, 19))), "Foo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(25, 1, 26), P(27, 1, 28))),
											ast.NewNilLiteralNode(L(S(P(25, 1, 26), P(27, 1, 28)))),
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
				L(S(P(0, 1, 1), P(31, 1, 32))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(31, 1, 32))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(31, 1, 32))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(27, 1, 28))),
									ast.NewPrivateConstantNode(L(S(P(16, 1, 17), P(18, 1, 19))), "_Fo"),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(25, 1, 26), P(27, 1, 28))),
											ast.NewNilLiteralNode(L(S(P(25, 1, 26), P(27, 1, 28)))),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(29, 1, 30))),
									ast.NewConstantLookupNode(
										L(S(P(16, 1, 17), P(20, 1, 21))),
										nil,
										ast.NewPublicConstantNode(L(S(P(18, 1, 19), P(20, 1, 21))), "Foo"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(27, 1, 28), P(29, 1, 30))),
											ast.NewNilLiteralNode(L(S(P(27, 1, 28), P(29, 1, 30)))),
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
				L(S(P(0, 1, 1), P(36, 1, 37))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(36, 1, 37))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(36, 1, 37))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(32, 1, 33))),
									ast.NewConstantLookupNode(
										L(S(P(16, 1, 17), P(23, 1, 24))),
										ast.NewPublicConstantNode(L(S(P(16, 1, 17), P(18, 1, 19))), "Foo"),
										ast.NewPublicConstantNode(L(S(P(21, 1, 22), P(23, 1, 24))), "Bar"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(30, 1, 31), P(32, 1, 33))),
											ast.NewNilLiteralNode(L(S(P(30, 1, 31), P(32, 1, 33)))),
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
				L(S(P(0, 1, 1), P(38, 1, 39))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(38, 1, 39))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(38, 1, 39))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(34, 1, 35))),
									ast.NewConstantLookupNode(
										L(S(P(16, 1, 17), P(25, 1, 26))),
										ast.NewConstantLookupNode(
											L(S(P(16, 1, 17), P(20, 1, 21))),
											nil,
											ast.NewPublicConstantNode(L(S(P(18, 1, 19), P(20, 1, 21))), "Foo"),
										),
										ast.NewPublicConstantNode(L(S(P(23, 1, 24), P(25, 1, 26))), "Bar"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(32, 1, 33), P(34, 1, 35))),
											ast.NewNilLiteralNode(L(S(P(32, 1, 33), P(34, 1, 35)))),
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
				L(S(P(0, 1, 1), P(31, 1, 32))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(31, 1, 32))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(31, 1, 32))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(27, 1, 28))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(18, 1, 19))),
										T(L(S(P(16, 1, 17), P(16, 1, 17))), token.LESS),
										ast.NewIntLiteralNode(L(S(P(18, 1, 19), P(18, 1, 19))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(25, 1, 26), P(27, 1, 28))),
											ast.NewNilLiteralNode(L(S(P(25, 1, 26), P(27, 1, 28)))),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(29, 1, 30))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(20, 1, 21))),
										T(L(S(P(16, 1, 17), P(16, 1, 17))), token.LESS),
										ast.NewPublicConstantNode(
											L(S(P(18, 1, 19), P(20, 1, 21))),
											"Foo",
										),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(27, 1, 28), P(29, 1, 30))),
											ast.NewNilLiteralNode(L(S(P(27, 1, 28), P(29, 1, 30)))),
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
				L(S(P(0, 1, 1), P(31, 1, 32))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(31, 1, 32))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(31, 1, 32))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(27, 1, 28))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(18, 1, 19))),
										T(L(S(P(16, 1, 17), P(16, 1, 17))), token.GREATER),
										ast.NewIntLiteralNode(L(S(P(18, 1, 19), P(18, 1, 19))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(25, 1, 26), P(27, 1, 28))),
											ast.NewNilLiteralNode(L(S(P(25, 1, 26), P(27, 1, 28)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(19, 1, 20))),
										T(L(S(P(16, 1, 17), P(17, 1, 18))), token.LESS_EQUAL),
										ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(19, 1, 20))),
										T(L(S(P(16, 1, 17), P(17, 1, 18))), token.GREATER_EQUAL),
										ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(19, 1, 20))),
										T(L(S(P(16, 1, 17), P(17, 1, 18))), token.EQUAL_EQUAL),
										ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(19, 1, 20))),
										T(L(S(P(16, 1, 17), P(17, 1, 18))), token.NOT_EQUAL),
										ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(19, 1, 20))),
										T(L(S(P(16, 1, 17), P(17, 1, 18))), token.LAX_EQUAL),
										ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(19, 1, 20))),
										T(L(S(P(16, 1, 17), P(17, 1, 18))), token.LAX_NOT_EQUAL),
										ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(29, 1, 30))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(20, 1, 21))),
										T(L(S(P(16, 1, 17), P(18, 1, 19))), token.STRICT_EQUAL),
										ast.NewIntLiteralNode(L(S(P(20, 1, 21), P(20, 1, 21))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(27, 1, 28), P(29, 1, 30))),
											ast.NewNilLiteralNode(L(S(P(27, 1, 28), P(29, 1, 30)))),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(29, 1, 30))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(20, 1, 21))),
										T(L(S(P(16, 1, 17), P(18, 1, 19))), token.STRICT_NOT_EQUAL),
										ast.NewIntLiteralNode(L(S(P(20, 1, 21), P(20, 1, 21))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(27, 1, 28), P(29, 1, 30))),
											ast.NewNilLiteralNode(L(S(P(27, 1, 28), P(29, 1, 30)))),
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
				L(S(P(0, 1, 1), P(34, 1, 35))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(34, 1, 35))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(34, 1, 35))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(30, 1, 31))),
									ast.NewUnaryExpressionNode(
										L(S(P(16, 1, 17), P(21, 1, 22))),
										T(L(S(P(16, 1, 17), P(18, 1, 19))), token.STRICT_NOT_EQUAL),
										ast.NewUnaryExpressionNode(
											L(S(P(20, 1, 21), P(21, 1, 22))),
											T(L(S(P(20, 1, 21), P(20, 1, 21))), token.MINUS),
											ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "5"),
										),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(28, 1, 29), P(30, 1, 31))),
											ast.NewNilLiteralNode(L(S(P(28, 1, 29), P(30, 1, 31)))),
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
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewBinaryPatternNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										T(L(S(P(20, 1, 21), P(21, 1, 22))), token.AND_AND),
										ast.NewUnaryExpressionNode(
											L(S(P(16, 1, 17), P(18, 1, 19))),
											T(L(S(P(16, 1, 17), P(16, 1, 17))), token.GREATER),
											ast.NewIntLiteralNode(L(S(P(18, 1, 19), P(18, 1, 19))), "5"),
										),
										ast.NewUnaryExpressionNode(
											L(S(P(23, 1, 24), P(26, 1, 27))),
											T(L(S(P(23, 1, 24), P(23, 1, 24))), token.LESS),
											ast.NewIntLiteralNode(L(S(P(25, 1, 26), P(26, 1, 27))), "10"),
										),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
				L(S(P(0, 1, 1), P(36, 1, 37))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(36, 1, 37))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(36, 1, 37))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(32, 1, 33))),
									ast.NewBinaryPatternNode(
										L(S(P(16, 1, 17), P(23, 1, 24))),
										T(L(S(P(20, 1, 21), P(21, 1, 22))), token.OR_OR),
										ast.NewUnaryExpressionNode(
											L(S(P(16, 1, 17), P(18, 1, 19))),
											T(L(S(P(16, 1, 17), P(16, 1, 17))), token.GREATER),
											ast.NewIntLiteralNode(L(S(P(18, 1, 19), P(18, 1, 19))), "5"),
										),
										ast.NewIntLiteralNode(L(S(P(23, 1, 24), P(23, 1, 24))), "2"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(30, 1, 31), P(32, 1, 33))),
											ast.NewNilLiteralNode(L(S(P(30, 1, 31), P(32, 1, 33)))),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(29, 1, 30))),
									ast.NewRangeLiteralNode(
										L(S(P(16, 1, 17), P(20, 1, 21))),
										T(L(S(P(17, 1, 18), P(19, 1, 20))), token.CLOSED_RANGE_OP),
										ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "2"),
										ast.NewIntLiteralNode(L(S(P(20, 1, 21), P(20, 1, 21))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(27, 1, 28), P(29, 1, 30))),
											ast.NewNilLiteralNode(L(S(P(27, 1, 28), P(29, 1, 30)))),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(29, 1, 30))),
									ast.NewRangeLiteralNode(
										L(S(P(16, 1, 17), P(20, 1, 21))),
										T(L(S(P(17, 1, 18), P(19, 1, 20))), token.CLOSED_RANGE_OP),
										ast.NewPublicConstantNode(L(S(P(16, 1, 17), P(16, 1, 17))), "A"),
										ast.NewPublicConstantNode(L(S(P(20, 1, 21), P(20, 1, 21))), "B"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(27, 1, 28), P(29, 1, 30))),
											ast.NewNilLiteralNode(L(S(P(27, 1, 28), P(29, 1, 30)))),
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
				L(S(P(0, 1, 1), P(35, 1, 36))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(35, 1, 36))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(35, 1, 36))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(31, 1, 32))),
									ast.NewRangeLiteralNode(
										L(S(P(16, 1, 17), P(22, 1, 23))),
										T(L(S(P(18, 1, 19), P(20, 1, 21))), token.CLOSED_RANGE_OP),
										ast.NewUnaryExpressionNode(
											L(S(P(16, 1, 17), P(17, 1, 18))),
											T(L(S(P(16, 1, 17), P(16, 1, 17))), token.MINUS),
											ast.NewIntLiteralNode(L(S(P(17, 1, 18), P(17, 1, 18))), "2"),
										),
										ast.NewUnaryExpressionNode(
											L(S(P(21, 1, 22), P(22, 1, 23))),
											T(L(S(P(21, 1, 22), P(21, 1, 22))), token.PLUS),
											ast.NewIntLiteralNode(L(S(P(22, 1, 23), P(22, 1, 23))), "5"),
										),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(29, 1, 30), P(31, 1, 32))),
											ast.NewNilLiteralNode(L(S(P(29, 1, 30), P(31, 1, 32)))),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(29, 1, 30))),
									ast.NewRangeLiteralNode(
										L(S(P(16, 1, 17), P(20, 1, 21))),
										T(L(S(P(17, 1, 18), P(19, 1, 20))), token.OPEN_RANGE_OP),
										ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "2"),
										ast.NewIntLiteralNode(L(S(P(20, 1, 21), P(20, 1, 21))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(27, 1, 28), P(29, 1, 30))),
											ast.NewNilLiteralNode(L(S(P(27, 1, 28), P(29, 1, 30)))),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(29, 1, 30))),
									ast.NewRangeLiteralNode(
										L(S(P(16, 1, 17), P(20, 1, 21))),
										T(L(S(P(17, 1, 18), P(19, 1, 20))), token.LEFT_OPEN_RANGE_OP),
										ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "2"),
										ast.NewIntLiteralNode(L(S(P(20, 1, 21), P(20, 1, 21))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(27, 1, 28), P(29, 1, 30))),
											ast.NewNilLiteralNode(L(S(P(27, 1, 28), P(29, 1, 30)))),
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
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(29, 1, 30))),
									ast.NewRangeLiteralNode(
										L(S(P(16, 1, 17), P(20, 1, 21))),
										T(L(S(P(17, 1, 18), P(19, 1, 20))), token.RIGHT_OPEN_RANGE_OP),
										ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "2"),
										ast.NewIntLiteralNode(L(S(P(20, 1, 21), P(20, 1, 21))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(27, 1, 28), P(29, 1, 30))),
											ast.NewNilLiteralNode(L(S(P(27, 1, 28), P(29, 1, 30)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewRangeLiteralNode(
										L(S(P(16, 1, 17), P(19, 1, 20))),
										T(L(S(P(16, 1, 17), P(18, 1, 19))), token.CLOSED_RANGE_OP),
										nil,
										ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewRangeLiteralNode(
										L(S(P(16, 1, 17), P(19, 1, 20))),
										T(L(S(P(16, 1, 17), P(18, 1, 19))), token.RIGHT_OPEN_RANGE_OP),
										nil,
										ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "5"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewRangeLiteralNode(
										L(S(P(16, 1, 17), P(19, 1, 20))),
										T(L(S(P(17, 1, 18), P(19, 1, 20))), token.CLOSED_RANGE_OP),
										ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "2"),
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(28, 1, 29))),
									ast.NewRangeLiteralNode(
										L(S(P(16, 1, 17), P(19, 1, 20))),
										T(L(S(P(17, 1, 18), P(19, 1, 20))), token.LEFT_OPEN_RANGE_OP),
										ast.NewIntLiteralNode(L(S(P(16, 1, 17), P(16, 1, 17))), "2"),
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(26, 1, 27), P(28, 1, 29))),
											ast.NewNilLiteralNode(L(S(P(26, 1, 27), P(28, 1, 29)))),
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

		"empty set pattern": {
			input: `switch foo case ^[] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(31, 1, 32))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(31, 1, 32))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(31, 1, 32))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(27, 1, 28))),
									ast.NewSetPatternNode(
										L(S(P(16, 1, 17), P(18, 1, 19))),
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(25, 1, 26), P(27, 1, 28))),
											ast.NewNilLiteralNode(L(S(P(25, 1, 26), P(27, 1, 28)))),
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
		"word set pattern": {
			input: `switch foo case ^w[foo bar] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewWordHashSetLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.WordCollectionContentNode{
											ast.NewRawStringLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "foo"),
											ast.NewRawStringLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "bar"),
										},
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
		"symbol set pattern": {
			input: `switch foo case ^s[foo bar] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewSymbolHashSetLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.SymbolCollectionContentNode{
											ast.NewSimpleSymbolLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "foo"),
											ast.NewSimpleSymbolLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "bar"),
										},
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
		"hex set pattern": {
			input: `switch foo case ^x[f5f 9e2] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewHexHashSetLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.IntCollectionContentNode{
											ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "0xf5f"),
											ast.NewIntLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "0x9e2"),
										},
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
		"bin set pattern": {
			input: `switch foo case ^b[101 111] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewBinHashSetLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.IntCollectionContentNode{
											ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "0b101"),
											ast.NewIntLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "0b111"),
										},
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
		"set with subpatterns": {
			input: `switch foo case ^[-1, "string", *, _] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(49, 1, 50))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(49, 1, 50))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(49, 1, 50))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(45, 1, 46))),
									ast.NewSetPatternNode(
										L(S(P(16, 1, 17), P(36, 1, 37))),
										[]ast.PatternNode{
											ast.NewUnaryExpressionNode(
												L(S(P(18, 1, 19), P(19, 1, 20))),
												T(L(S(P(18, 1, 19), P(18, 1, 19))), token.MINUS),
												ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "1"),
											),
											ast.NewDoubleQuotedStringLiteralNode(L(S(P(22, 1, 23), P(29, 1, 30))), "string"),
											ast.NewRestPatternNode(L(S(P(32, 1, 33), P(32, 1, 33))), nil),
											ast.NewPrivateIdentifierNode(L(S(P(35, 1, 36), P(35, 1, 36))), "_"),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(43, 1, 44), P(45, 1, 46))),
											ast.NewNilLiteralNode(L(S(P(43, 1, 44), P(45, 1, 46)))),
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
		"set with constants": {
			input: `switch foo case ^[-1, Bea::Fin, *, _Foo] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(52, 1, 53))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(52, 1, 53))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(52, 1, 53))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(48, 1, 49))),
									ast.NewSetPatternNode(
										L(S(P(16, 1, 17), P(39, 1, 40))),
										[]ast.PatternNode{
											ast.NewUnaryExpressionNode(
												L(S(P(18, 1, 19), P(19, 1, 20))),
												T(L(S(P(18, 1, 19), P(18, 1, 19))), token.MINUS),
												ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "1"),
											),
											ast.NewConstantLookupNode(
												L(S(P(22, 1, 23), P(29, 1, 30))),
												ast.NewPublicConstantNode(L(S(P(22, 1, 23), P(24, 1, 25))), "Bea"),
												ast.NewPublicConstantNode(L(S(P(27, 1, 28), P(29, 1, 30))), "Fin"),
											),
											ast.NewRestPatternNode(L(S(P(32, 1, 33), P(32, 1, 33))), nil),
											ast.NewPrivateConstantNode(L(S(P(35, 1, 36), P(38, 1, 39))), "_Foo"),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(46, 1, 47), P(48, 1, 49))),
											ast.NewNilLiteralNode(L(S(P(46, 1, 47), P(48, 1, 49)))),
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
		"set with invalid private identifier": {
			input: `switch foo case ^[-1, "string", *, _foo] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(52, 1, 53))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(52, 1, 53))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(52, 1, 53))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(48, 1, 49))),
									ast.NewSetPatternNode(
										L(S(P(16, 1, 17), P(39, 1, 40))),
										[]ast.PatternNode{
											ast.NewUnaryExpressionNode(
												L(S(P(18, 1, 19), P(19, 1, 20))),
												T(L(S(P(18, 1, 19), P(18, 1, 19))), token.MINUS),
												ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "1"),
											),
											ast.NewDoubleQuotedStringLiteralNode(L(S(P(22, 1, 23), P(29, 1, 30))), "string"),
											ast.NewRestPatternNode(L(S(P(32, 1, 33), P(32, 1, 33))), nil),
											ast.NewPrivateIdentifierNode(L(S(P(35, 1, 36), P(38, 1, 39))), "_foo"),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(46, 1, 47), P(48, 1, 49))),
											ast.NewNilLiteralNode(L(S(P(46, 1, 47), P(48, 1, 49)))),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(35, 1, 36), P(38, 1, 39))), "set patterns cannot contain identifiers other than _"),
			},
		},
		"empty list pattern": {
			input: `switch foo case [] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(30, 1, 31))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(30, 1, 31))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(30, 1, 31))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(26, 1, 27))),
									ast.NewListPatternNode(
										L(S(P(16, 1, 17), P(17, 1, 18))),
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(24, 1, 25), P(26, 1, 27))),
											ast.NewNilLiteralNode(L(S(P(24, 1, 25), P(26, 1, 27)))),
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
		"word list pattern": {
			input: `switch foo case \w[foo bar] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewWordArrayListLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.WordCollectionContentNode{
											ast.NewRawStringLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "foo"),
											ast.NewRawStringLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "bar"),
										},
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
		"symbol list pattern": {
			input: `switch foo case \s[foo bar] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewSymbolArrayListLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.SymbolCollectionContentNode{
											ast.NewSimpleSymbolLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "foo"),
											ast.NewSimpleSymbolLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "bar"),
										},
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
		"hex list pattern": {
			input: `switch foo case \x[f5f 9e2] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewHexArrayListLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.IntCollectionContentNode{
											ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "0xf5f"),
											ast.NewIntLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "0x9e2"),
										},
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
		"bin list pattern": {
			input: `switch foo case \b[101 111] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewBinArrayListLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.IntCollectionContentNode{
											ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "0b101"),
											ast.NewIntLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "0b111"),
										},
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
			input: `switch foo case [a, > 6 && < 20, [*b, :foo]] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(56, 1, 57))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(56, 1, 57))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(56, 1, 57))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(52, 1, 53))),
									ast.NewListPatternNode(
										L(S(P(16, 1, 17), P(43, 1, 44))),
										[]ast.PatternNode{
											ast.NewPublicIdentifierNode(
												L(S(P(17, 1, 18), P(17, 1, 18))),
												"a",
											),
											ast.NewBinaryPatternNode(
												L(S(P(20, 1, 21), P(30, 1, 31))),
												T(L(S(P(24, 1, 25), P(25, 1, 26))), token.AND_AND),
												ast.NewUnaryExpressionNode(
													L(S(P(20, 1, 21), P(22, 1, 23))),
													T(L(S(P(20, 1, 21), P(20, 1, 21))), token.GREATER),
													ast.NewIntLiteralNode(L(S(P(22, 1, 23), P(22, 1, 23))), "6"),
												),
												ast.NewUnaryExpressionNode(
													L(S(P(27, 1, 28), P(30, 1, 31))),
													T(L(S(P(27, 1, 28), P(27, 1, 28))), token.LESS),
													ast.NewIntLiteralNode(L(S(P(29, 1, 30), P(30, 1, 31))), "20"),
												),
											),
											ast.NewListPatternNode(
												L(S(P(33, 1, 34), P(42, 1, 43))),
												[]ast.PatternNode{
													ast.NewRestPatternNode(
														L(S(P(34, 1, 35), P(35, 1, 36))),
														ast.NewPublicIdentifierNode(L(S(P(35, 1, 36), P(35, 1, 36))), "b"),
													),
													ast.NewSimpleSymbolLiteralNode(L(S(P(38, 1, 39), P(41, 1, 42))), "foo"),
												},
											),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(50, 1, 51), P(52, 1, 53))),
											ast.NewNilLiteralNode(L(S(P(50, 1, 51), P(52, 1, 53)))),
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
		"list pattern with unnamed rest element": {
			input: `switch foo case [*, 2] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(34, 1, 35))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(34, 1, 35))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(34, 1, 35))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(30, 1, 31))),
									ast.NewListPatternNode(
										L(S(P(16, 1, 17), P(21, 1, 22))),
										[]ast.PatternNode{
											ast.NewRestPatternNode(
												L(S(P(17, 1, 18), P(17, 1, 18))),
												nil,
											),
											ast.NewIntLiteralNode(
												L(S(P(20, 1, 21), P(20, 1, 21))),
												"2",
											),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(28, 1, 29), P(30, 1, 31))),
											ast.NewNilLiteralNode(L(S(P(28, 1, 29), P(30, 1, 31)))),
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
		"empty tuple pattern": {
			input: `switch foo case %[] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(31, 1, 32))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(31, 1, 32))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(31, 1, 32))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(27, 1, 28))),
									ast.NewTuplePatternNode(
										L(S(P(16, 1, 17), P(18, 1, 19))),
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(25, 1, 26), P(27, 1, 28))),
											ast.NewNilLiteralNode(L(S(P(25, 1, 26), P(27, 1, 28)))),
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
		"word tuple pattern": {
			input: `switch foo case %w[foo bar] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewWordArrayTupleLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.WordCollectionContentNode{
											ast.NewRawStringLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "foo"),
											ast.NewRawStringLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "bar"),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
		"symbol tuple pattern": {
			input: `switch foo case %s[foo bar] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewSymbolArrayTupleLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.SymbolCollectionContentNode{
											ast.NewSimpleSymbolLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "foo"),
											ast.NewSimpleSymbolLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "bar"),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
		"hex tuple pattern": {
			input: `switch foo case %x[f5f 9e2] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewHexArrayTupleLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.IntCollectionContentNode{
											ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "0xf5f"),
											ast.NewIntLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "0x9e2"),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
		"bin tuple pattern": {
			input: `switch foo case %b[101 111] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(39, 1, 40))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(39, 1, 40))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(39, 1, 40))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(35, 1, 36))),
									ast.NewBinArrayTupleLiteralNode(
										L(S(P(16, 1, 17), P(26, 1, 27))),
										[]ast.IntCollectionContentNode{
											ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(21, 1, 22))), "0b101"),
											ast.NewIntLiteralNode(L(S(P(23, 1, 24), P(25, 1, 26))), "0b111"),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(33, 1, 34), P(35, 1, 36))),
											ast.NewNilLiteralNode(L(S(P(33, 1, 34), P(35, 1, 36)))),
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
		"tuple with subpatterns": {
			input: `switch foo case %[a, > 6 && < 20, %[*b, :foo]] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(58, 1, 59))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(58, 1, 59))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(58, 1, 59))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(54, 1, 55))),
									ast.NewTuplePatternNode(
										L(S(P(16, 1, 17), P(45, 1, 46))),
										[]ast.PatternNode{
											ast.NewPublicIdentifierNode(
												L(S(P(18, 1, 19), P(18, 1, 19))),
												"a",
											),
											ast.NewBinaryPatternNode(
												L(S(P(21, 1, 22), P(31, 1, 32))),
												T(L(S(P(25, 1, 26), P(26, 1, 27))), token.AND_AND),
												ast.NewUnaryExpressionNode(
													L(S(P(21, 1, 22), P(23, 1, 24))),
													T(L(S(P(21, 1, 22), P(21, 1, 22))), token.GREATER),
													ast.NewIntLiteralNode(L(S(P(23, 1, 24), P(23, 1, 24))), "6"),
												),
												ast.NewUnaryExpressionNode(
													L(S(P(28, 1, 29), P(31, 1, 32))),
													T(L(S(P(28, 1, 29), P(28, 1, 29))), token.LESS),
													ast.NewIntLiteralNode(L(S(P(30, 1, 31), P(31, 1, 32))), "20"),
												),
											),
											ast.NewTuplePatternNode(
												L(S(P(34, 1, 35), P(44, 1, 45))),
												[]ast.PatternNode{
													ast.NewRestPatternNode(
														L(S(P(36, 1, 37), P(37, 1, 38))),
														ast.NewPublicIdentifierNode(L(S(P(37, 1, 38), P(37, 1, 38))), "b"),
													),
													ast.NewSimpleSymbolLiteralNode(L(S(P(40, 1, 41), P(43, 1, 44))), "foo"),
												},
											),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(52, 1, 53), P(54, 1, 55))),
											ast.NewNilLiteralNode(L(S(P(52, 1, 53), P(54, 1, 55)))),
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
		"tuple pattern with unnamed rest element": {
			input: `switch foo case %[*, 2] then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(35, 1, 36))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(35, 1, 36))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(35, 1, 36))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(31, 1, 32))),
									ast.NewTuplePatternNode(
										L(S(P(16, 1, 17), P(22, 1, 23))),
										[]ast.PatternNode{
											ast.NewRestPatternNode(
												L(S(P(18, 1, 19), P(18, 1, 19))),
												nil,
											),
											ast.NewIntLiteralNode(
												L(S(P(21, 1, 22), P(21, 1, 22))),
												"2",
											),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(29, 1, 30), P(31, 1, 32))),
											ast.NewNilLiteralNode(L(S(P(29, 1, 30), P(31, 1, 32)))),
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

		"empty map pattern": {
			input: `switch foo case {} then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(30, 1, 31))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(30, 1, 31))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(30, 1, 31))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(26, 1, 27))),
									ast.NewMapPatternNode(
										L(S(P(16, 1, 17), P(17, 1, 18))),
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(24, 1, 25), P(26, 1, 27))),
											ast.NewNilLiteralNode(L(S(P(24, 1, 25), P(26, 1, 27)))),
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
		"map with subpatterns": {
			input: `switch foo case {a, 1 => > 6 && < 20, foo: { "foo" => ["baz", *] } } then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(80, 1, 81))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(80, 1, 81))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(80, 1, 81))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(76, 1, 77))),
									ast.NewMapPatternNode(
										L(S(P(16, 1, 17), P(67, 1, 68))),
										[]ast.PatternNode{
											ast.NewPublicIdentifierNode(
												L(S(P(17, 1, 18), P(17, 1, 18))),
												"a",
											),
											ast.NewKeyValuePatternNode(
												L(S(P(20, 1, 21), P(35, 1, 36))),
												ast.NewIntLiteralNode(L(S(P(20, 1, 21), P(20, 1, 21))), "1"),
												ast.NewBinaryPatternNode(
													L(S(P(25, 1, 26), P(35, 1, 36))),
													T(L(S(P(29, 1, 30), P(30, 1, 31))), token.AND_AND),
													ast.NewUnaryExpressionNode(
														L(S(P(25, 1, 26), P(27, 1, 28))),
														T(L(S(P(25, 1, 26), P(25, 1, 26))), token.GREATER),
														ast.NewIntLiteralNode(L(S(P(27, 1, 28), P(27, 1, 28))), "6"),
													),
													ast.NewUnaryExpressionNode(
														L(S(P(32, 1, 33), P(35, 1, 36))),
														T(L(S(P(32, 1, 33), P(32, 1, 33))), token.LESS),
														ast.NewIntLiteralNode(L(S(P(34, 1, 35), P(35, 1, 36))), "20"),
													),
												),
											),
											ast.NewSymbolKeyValuePatternNode(
												L(S(P(38, 1, 39), P(65, 1, 66))),
												"foo",
												ast.NewMapPatternNode(
													L(S(P(43, 1, 44), P(65, 1, 66))),
													[]ast.PatternNode{
														ast.NewKeyValuePatternNode(
															L(S(P(45, 1, 46), P(63, 1, 64))),
															ast.NewDoubleQuotedStringLiteralNode(L(S(P(45, 1, 46), P(49, 1, 50))), "foo"),
															ast.NewListPatternNode(
																L(S(P(54, 1, 55), P(63, 1, 64))),
																[]ast.PatternNode{
																	ast.NewDoubleQuotedStringLiteralNode(L(S(P(55, 1, 56), P(59, 1, 60))), "baz"),
																	ast.NewRestPatternNode(L(S(P(62, 1, 63), P(62, 1, 63))), nil),
																},
															),
														),
													},
												),
											),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(74, 1, 75), P(76, 1, 77))),
											ast.NewNilLiteralNode(L(S(P(74, 1, 75), P(76, 1, 77)))),
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

		"empty record pattern": {
			input: `switch foo case %{} then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(31, 1, 32))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(31, 1, 32))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(31, 1, 32))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(27, 1, 28))),
									ast.NewRecordPatternNode(
										L(S(P(16, 1, 17), P(18, 1, 19))),
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(25, 1, 26), P(27, 1, 28))),
											ast.NewNilLiteralNode(L(S(P(25, 1, 26), P(27, 1, 28)))),
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
		"record with subpatterns": {
			input: `switch foo case %{a, 1 => > 6 && < 20, foo: %{ "foo" => ["baz", *] } } then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(82, 1, 83))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(82, 1, 83))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(82, 1, 83))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(78, 1, 79))),
									ast.NewRecordPatternNode(
										L(S(P(16, 1, 17), P(69, 1, 70))),
										[]ast.PatternNode{
											ast.NewPublicIdentifierNode(
												L(S(P(18, 1, 19), P(18, 1, 19))),
												"a",
											),
											ast.NewKeyValuePatternNode(
												L(S(P(21, 1, 22), P(36, 1, 37))),
												ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "1"),
												ast.NewBinaryPatternNode(
													L(S(P(26, 1, 27), P(36, 1, 37))),
													T(L(S(P(30, 1, 31), P(31, 1, 32))), token.AND_AND),
													ast.NewUnaryExpressionNode(
														L(S(P(26, 1, 27), P(28, 1, 29))),
														T(L(S(P(26, 1, 27), P(26, 1, 27))), token.GREATER),
														ast.NewIntLiteralNode(L(S(P(28, 1, 29), P(28, 1, 29))), "6"),
													),
													ast.NewUnaryExpressionNode(
														L(S(P(33, 1, 34), P(36, 1, 37))),
														T(L(S(P(33, 1, 34), P(33, 1, 34))), token.LESS),
														ast.NewIntLiteralNode(L(S(P(35, 1, 36), P(36, 1, 37))), "20"),
													),
												),
											),
											ast.NewSymbolKeyValuePatternNode(
												L(S(P(39, 1, 40), P(67, 1, 68))),
												"foo",
												ast.NewRecordPatternNode(
													L(S(P(44, 1, 45), P(67, 1, 68))),
													[]ast.PatternNode{
														ast.NewKeyValuePatternNode(
															L(S(P(47, 1, 48), P(65, 1, 66))),
															ast.NewDoubleQuotedStringLiteralNode(L(S(P(47, 1, 48), P(51, 1, 52))), "foo"),
															ast.NewListPatternNode(
																L(S(P(56, 1, 57), P(65, 1, 66))),
																[]ast.PatternNode{
																	ast.NewDoubleQuotedStringLiteralNode(L(S(P(57, 1, 58), P(61, 1, 62))), "baz"),
																	ast.NewRestPatternNode(L(S(P(64, 1, 65), P(64, 1, 65))), nil),
																},
															),
														),
													},
												),
											),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(76, 1, 77), P(78, 1, 79))),
											ast.NewNilLiteralNode(L(S(P(76, 1, 77), P(78, 1, 79)))),
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

		"empty object pattern": {
			input: `switch foo case ::Foo() then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(35, 1, 36))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(35, 1, 36))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(35, 1, 36))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(31, 1, 32))),
									ast.NewObjectPatternNode(
										L(S(P(16, 1, 17), P(22, 1, 23))),
										ast.NewConstantLookupNode(
											L(S(P(16, 1, 17), P(20, 1, 21))),
											nil,
											ast.NewPublicConstantNode(
												L(S(P(18, 1, 19), P(20, 1, 21))),
												"Foo",
											),
										),
										nil,
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(29, 1, 30), P(31, 1, 32))),
											ast.NewNilLiteralNode(L(S(P(29, 1, 30), P(31, 1, 32)))),
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
		"object with subpatterns": {
			input: `switch foo case Foo(a, bar: Bar(x: > 6 && < 20), foo: %{ "foo" => ["baz", *] }) then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(91, 1, 92))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(91, 1, 92))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(91, 1, 92))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(87, 1, 88))),
									ast.NewObjectPatternNode(
										L(S(P(16, 1, 17), P(78, 1, 79))),
										ast.NewPublicConstantNode(
											L(S(P(16, 1, 17), P(18, 1, 19))),
											"Foo",
										),
										[]ast.PatternNode{
											ast.NewPublicIdentifierNode(
												L(S(P(20, 1, 21), P(20, 1, 21))),
												"a",
											),
											ast.NewSymbolKeyValuePatternNode(
												L(S(P(23, 1, 24), P(46, 1, 47))),
												"bar",
												ast.NewObjectPatternNode(
													L(S(P(28, 1, 29), P(46, 1, 47))),
													ast.NewPublicConstantNode(
														L(S(P(28, 1, 29), P(30, 1, 31))),
														"Bar",
													),
													[]ast.PatternNode{
														ast.NewSymbolKeyValuePatternNode(
															L(S(P(32, 1, 33), P(45, 1, 46))),
															"x",
															ast.NewBinaryPatternNode(
																L(S(P(35, 1, 36), P(45, 1, 46))),
																T(L(S(P(39, 1, 40), P(40, 1, 41))), token.AND_AND),
																ast.NewUnaryExpressionNode(
																	L(S(P(35, 1, 36), P(37, 1, 38))),
																	T(L(S(P(35, 1, 36), P(35, 1, 36))), token.GREATER),
																	ast.NewIntLiteralNode(L(S(P(37, 1, 38), P(37, 1, 38))), "6"),
																),
																ast.NewUnaryExpressionNode(
																	L(S(P(42, 1, 43), P(45, 1, 46))),
																	T(L(S(P(42, 1, 43), P(42, 1, 43))), token.LESS),
																	ast.NewIntLiteralNode(L(S(P(44, 1, 45), P(45, 1, 46))), "20"),
																),
															),
														),
													},
												),
											),
											ast.NewSymbolKeyValuePatternNode(
												L(S(P(49, 1, 50), P(77, 1, 78))),
												"foo",
												ast.NewRecordPatternNode(
													L(S(P(54, 1, 55), P(77, 1, 78))),
													[]ast.PatternNode{
														ast.NewKeyValuePatternNode(
															L(S(P(57, 1, 58), P(75, 1, 76))),
															ast.NewDoubleQuotedStringLiteralNode(L(S(P(57, 1, 58), P(61, 1, 62))), "foo"),
															ast.NewListPatternNode(
																L(S(P(66, 1, 67), P(75, 1, 76))),
																[]ast.PatternNode{
																	ast.NewDoubleQuotedStringLiteralNode(L(S(P(67, 1, 68), P(71, 1, 72))), "baz"),
																	ast.NewRestPatternNode(L(S(P(74, 1, 75), P(74, 1, 75))), nil),
																},
															),
														),
													},
												),
											),
										},
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(85, 1, 86), P(87, 1, 88))),
											ast.NewNilLiteralNode(L(S(P(85, 1, 86), P(87, 1, 88)))),
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
		"as pattern": {
			input: `switch foo case > 5 || 2 as bar then nil end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(43, 1, 44))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(43, 1, 44))),
						ast.NewSwitchExpressionNode(
							L(S(P(0, 1, 1), P(43, 1, 44))),
							ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
							[]*ast.CaseNode{
								ast.NewCaseNode(
									L(S(P(11, 1, 12), P(39, 1, 40))),
									ast.NewAsPatternNode(
										L(S(P(16, 1, 17), P(30, 1, 31))),
										ast.NewBinaryPatternNode(
											L(S(P(16, 1, 17), P(23, 1, 24))),
											T(L(S(P(20, 1, 21), P(21, 1, 22))), token.OR_OR),
											ast.NewUnaryExpressionNode(
												L(S(P(16, 1, 17), P(18, 1, 19))),
												T(L(S(P(16, 1, 17), P(16, 1, 17))), token.GREATER),
												ast.NewIntLiteralNode(L(S(P(18, 1, 19), P(18, 1, 19))), "5"),
											),
											ast.NewIntLiteralNode(L(S(P(23, 1, 24), P(23, 1, 24))), "2"),
										),
										ast.NewPublicIdentifierNode(L(S(P(28, 1, 29), P(30, 1, 31))), "bar"),
									),
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(37, 1, 38), P(39, 1, 40))),
											ast.NewNilLiteralNode(L(S(P(37, 1, 38), P(39, 1, 40)))),
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
