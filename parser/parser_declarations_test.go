package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
)

func TestSingletonBlock(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
singleton
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(29, 5, 4)),
						ast.NewSingletonBlockExpressionNode(
							S(P(1, 2, 1), P(28, 5, 3)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(12, 3, 2), P(20, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(12, 3, 2), P(19, 3, 9)),
										T(S(P(16, 3, 6), P(17, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(12, 3, 2), P(14, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(19, 3, 9), P(19, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(22, 4, 2), P(25, 4, 5)),
									ast.NewNilLiteralNode(S(P(22, 4, 2), P(24, 4, 4))),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
singleton
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(14, 3, 4)),
						ast.NewSingletonBlockExpressionNode(
							S(P(1, 2, 1), P(13, 3, 3)),
							nil,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
bar =
	singleton
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(37, 6, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(33, 5, 5)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(32, 5, 4)),
							T(S(P(5, 2, 5), P(5, 2, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 3)), "bar"),
							ast.NewSingletonBlockExpressionNode(
								S(P(8, 3, 2), P(32, 5, 4)),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(20, 4, 3), P(28, 4, 11)),
										ast.NewAssignmentExpressionNode(
											S(P(20, 4, 3), P(27, 4, 10)),
											T(S(P(24, 4, 7), P(25, 4, 8)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(20, 4, 3), P(22, 4, 5)), "foo"),
											ast.NewIntLiteralNode(S(P(27, 4, 10), P(27, 4, 10)), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(34, 6, 1), P(37, 6, 4)),
						ast.NewNilLiteralNode(S(P(34, 6, 1), P(36, 6, 3))),
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

func TestDocComment(t *testing.T) {
	tests := testTable{
		"cannot omit the argument": {
			input: "##[foo]##",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewDocCommentNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							"foo",
							ast.NewInvalidNode(
								S(P(9, 1, 10), P(8, 1, 9)),
								T(S(P(9, 1, 10), P(8, 1, 9)), token.END_OF_FILE),
							),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(9, 1, 10), P(8, 1, 9)), "unexpected END_OF_FILE, expected an expression"),
			},
		},
		"cannot be nested": {
			input: "##[foo]## ##[bar]## 1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewDocCommentNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							"foo",
							ast.NewDocCommentNode(
								S(P(10, 1, 11), P(20, 1, 21)),
								"bar",
								ast.NewIntLiteralNode(S(P(20, 1, 21), P(20, 1, 21)), "1"),
							),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(10, 1, 11), P(18, 1, 19)), "doc comments cannot document one another"),
			},
		},
		"can be empty": {
			input: "##[]## def foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewDocCommentNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							"",
							ast.NewMethodDefinitionNode(
								S(P(7, 1, 8), P(18, 1, 19)),
								"foo",
								nil,
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: `
				##[
					foo
					bar
				]##
				def foo; end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(51, 6, 17)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(
						S(P(0, 1, 1), P(0, 1, 1)),
					),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(51, 6, 17)),
						ast.NewDocCommentNode(
							S(P(5, 2, 5), P(50, 6, 16)),
							"foo\nbar",
							ast.NewMethodDefinitionNode(
								S(P(39, 6, 5), P(50, 6, 16)),
								"foo",
								nil,
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can document expressions": {
			input: "##[foo]## 1 + class Foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewDocCommentNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							"foo",
							ast.NewBinaryExpressionNode(
								S(P(10, 1, 11), P(27, 1, 28)),
								T(S(P(12, 1, 13), P(12, 1, 13)), token.PLUS),
								ast.NewIntLiteralNode(S(P(10, 1, 11), P(10, 1, 11)), "1"),
								ast.NewClassDeclarationNode(
									S(P(14, 1, 15), P(27, 1, 28)),
									false,
									false,
									ast.NewPublicConstantNode(S(P(20, 1, 21), P(22, 1, 23)), "Foo"),
									nil,
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

func TestIncludeExpression(t *testing.T) {
	tests := testTable{
		"cannot omit the argument": {
			input: "include",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewIncludeExpressionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							[]ast.ComplexConstantNode{
								ast.NewInvalidNode(
									S(P(7, 1, 8), P(6, 1, 7)),
									T(S(P(7, 1, 8), P(6, 1, 7)), token.END_OF_FILE),
								),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(7, 1, 8), P(6, 1, 7)), "unexpected END_OF_FILE, expected a constant"),
			},
		},
		"can have a public constant as the argument": {
			input: "include Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewIncludeExpressionNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(8, 1, 9), P(17, 1, 18)),
									"Enumerable",
								),
							},
						),
					),
				},
			),
		},
		"can have many arguments": {
			input: "include Enumerable, Memoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewIncludeExpressionNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(8, 1, 9), P(17, 1, 18)),
									"Enumerable",
								),
								ast.NewPublicConstantNode(
									S(P(20, 1, 21), P(29, 1, 30)),
									"Memoizable",
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after the comma": {
			input: "include Enumerable,\nMemoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 2, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 2, 10)),
						ast.NewIncludeExpressionNode(
							S(P(0, 1, 1), P(29, 2, 10)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(8, 1, 9), P(17, 1, 18)),
									"Enumerable",
								),
								ast.NewPublicConstantNode(
									S(P(20, 2, 1), P(29, 2, 10)),
									"Memoizable",
								),
							},
						),
					),
				},
			),
		},
		"can have a private constant as the argument": {
			input: "include _Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewIncludeExpressionNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							[]ast.ComplexConstantNode{
								ast.NewPrivateConstantNode(
									S(P(8, 1, 9), P(18, 1, 19)),
									"_Enumerable",
								),
							},
						),
					),
				},
			),
		},
		"can have a constant lookup as the argument": {
			input: "include Std::Memoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewIncludeExpressionNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							[]ast.ComplexConstantNode{
								ast.NewConstantLookupNode(
									S(P(8, 1, 9), P(22, 1, 23)),
									ast.NewPublicConstantNode(
										S(P(8, 1, 9), P(10, 1, 11)),
										"Std",
									),
									ast.NewPublicConstantNode(
										S(P(13, 1, 14), P(22, 1, 23)),
										"Memoizable",
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a generic constant as the argument": {
			input: "include Enumerable[String]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewIncludeExpressionNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							[]ast.ComplexConstantNode{
								ast.NewGenericConstantNode(
									S(P(8, 1, 9), P(25, 1, 26)),
									ast.NewPublicConstantNode(S(P(8, 1, 9), P(17, 1, 18)), "Enumerable"),
									[]ast.ComplexConstantNode{
										ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									},
								),
							},
						),
					),
				},
			),
		},
		"can be repeated": {
			input: `
				include Foo
				include Bar
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 3, 16)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(16, 2, 16)),
						ast.NewIncludeExpressionNode(
							S(P(5, 2, 5), P(15, 2, 15)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(S(P(13, 2, 13), P(15, 2, 15)), "Foo"),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(21, 3, 5), P(32, 3, 16)),
						ast.NewIncludeExpressionNode(
							S(P(21, 3, 5), P(31, 3, 15)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(S(P(29, 3, 13), P(31, 3, 15)), "Bar"),
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

func TestExtendExpression(t *testing.T) {
	tests := testTable{
		"cannot omit the argument": {
			input: "extend",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewExtendExpressionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							[]ast.ComplexConstantNode{
								ast.NewInvalidNode(
									S(P(6, 1, 7), P(5, 1, 6)),
									T(S(P(6, 1, 7), P(5, 1, 6)), token.END_OF_FILE),
								),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(6, 1, 7), P(5, 1, 6)), "unexpected END_OF_FILE, expected a constant"),
			},
		},
		"can have a public constant as the argument": {
			input: "extend Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewExtendExpressionNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(7, 1, 8), P(16, 1, 17)),
									"Enumerable",
								),
							},
						),
					),
				},
			),
		},
		"can have many arguments": {
			input: "extend Enumerable, Memoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 1, 29)),
						ast.NewExtendExpressionNode(
							S(P(0, 1, 1), P(28, 1, 29)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(7, 1, 8), P(16, 1, 17)),
									"Enumerable",
								),
								ast.NewPublicConstantNode(
									S(P(19, 1, 20), P(28, 1, 29)),
									"Memoizable",
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after the comma": {
			input: "extend Enumerable,\nMemoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 2, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 2, 10)),
						ast.NewExtendExpressionNode(
							S(P(0, 1, 1), P(28, 2, 10)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(7, 1, 8), P(16, 1, 17)),
									"Enumerable",
								),
								ast.NewPublicConstantNode(
									S(P(19, 2, 1), P(28, 2, 10)),
									"Memoizable",
								),
							},
						),
					),
				},
			),
		},
		"can have a private constant as the argument": {
			input: "extend _Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewExtendExpressionNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							[]ast.ComplexConstantNode{
								ast.NewPrivateConstantNode(
									S(P(7, 1, 8), P(17, 1, 18)),
									"_Enumerable",
								),
							},
						),
					),
				},
			),
		},
		"can have a constant lookup as the argument": {
			input: "extend Std::Memoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewExtendExpressionNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							[]ast.ComplexConstantNode{
								ast.NewConstantLookupNode(
									S(P(7, 1, 8), P(21, 1, 22)),
									ast.NewPublicConstantNode(
										S(P(7, 1, 8), P(9, 1, 10)),
										"Std",
									),
									ast.NewPublicConstantNode(
										S(P(12, 1, 13), P(21, 1, 22)),
										"Memoizable",
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a generic constant as the argument": {
			input: "extend Enumerable[String]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewExtendExpressionNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							[]ast.ComplexConstantNode{
								ast.NewGenericConstantNode(
									S(P(7, 1, 8), P(24, 1, 25)),
									ast.NewPublicConstantNode(S(P(7, 1, 8), P(16, 1, 17)), "Enumerable"),
									[]ast.ComplexConstantNode{
										ast.NewPublicConstantNode(S(P(18, 1, 19), P(23, 1, 24)), "String"),
									},
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

func TestEnhanceExpression(t *testing.T) {
	tests := testTable{
		"cannot omit the argument": {
			input: "enhance",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewEnhanceExpressionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							[]ast.ComplexConstantNode{
								ast.NewInvalidNode(
									S(P(7, 1, 8), P(6, 1, 7)),
									T(S(P(7, 1, 8), P(6, 1, 7)), token.END_OF_FILE),
								),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(7, 1, 8), P(6, 1, 7)), "unexpected END_OF_FILE, expected a constant"),
			},
		},
		"can have a public constant as the argument": {
			input: "enhance Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewEnhanceExpressionNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(8, 1, 9), P(17, 1, 18)),
									"Enumerable",
								),
							},
						),
					),
				},
			),
		},
		"can have many arguments": {
			input: "enhance Enumerable, Memoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewEnhanceExpressionNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(8, 1, 9), P(17, 1, 18)),
									"Enumerable",
								),
								ast.NewPublicConstantNode(
									S(P(20, 1, 21), P(29, 1, 30)),
									"Memoizable",
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after the comma": {
			input: "enhance Enumerable,\nMemoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 2, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 2, 10)),
						ast.NewEnhanceExpressionNode(
							S(P(0, 1, 1), P(29, 2, 10)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(8, 1, 9), P(17, 1, 18)),
									"Enumerable",
								),
								ast.NewPublicConstantNode(
									S(P(20, 2, 1), P(29, 2, 10)),
									"Memoizable",
								),
							},
						),
					),
				},
			),
		},
		"can have a private constant as the argument": {
			input: "enhance _Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewEnhanceExpressionNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							[]ast.ComplexConstantNode{
								ast.NewPrivateConstantNode(
									S(P(8, 1, 9), P(18, 1, 19)),
									"_Enumerable",
								),
							},
						),
					),
				},
			),
		},
		"can have a constant lookup as the argument": {
			input: "enhance Std::Memoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewEnhanceExpressionNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							[]ast.ComplexConstantNode{
								ast.NewConstantLookupNode(
									S(P(8, 1, 9), P(22, 1, 23)),
									ast.NewPublicConstantNode(
										S(P(8, 1, 9), P(10, 1, 11)),
										"Std",
									),
									ast.NewPublicConstantNode(
										S(P(13, 1, 14), P(22, 1, 23)),
										"Memoizable",
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a generic constant as the argument": {
			input: "enhance Enumerable[String]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewEnhanceExpressionNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							[]ast.ComplexConstantNode{
								ast.NewGenericConstantNode(
									S(P(8, 1, 9), P(25, 1, 26)),
									ast.NewPublicConstantNode(S(P(8, 1, 9), P(17, 1, 18)), "Enumerable"),
									[]ast.ComplexConstantNode{
										ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									},
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

func TestValueDeclaration(t *testing.T) {
	tests := testTable{
		"is valid without a type or initialiser": {
			input: "val foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be a part of an expression": {
			input: "a = val foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
							ast.NewValueDeclarationNode(
								S(P(4, 1, 5), P(10, 1, 11)),
								V(S(P(8, 1, 9), P(10, 1, 11)), token.PUBLIC_IDENTIFIER, "foo"),
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have a private identifier as the value name": {
			input: "val _foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							V(S(P(4, 1, 5), P(7, 1, 8)), token.PRIVATE_IDENTIFIER, "_foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have an instance variable as the value name": {
			input: "val @foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							V(S(P(4, 1, 5), P(7, 1, 8)), token.INSTANCE_VARIABLE, "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(7, 1, 8)), "instance variables cannot be declared using `val`"),
			},
		},
		"cannot have a constant as the value name": {
			input: "val Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(4, 1, 5), P(6, 1, 7)),
						ast.NewInvalidNode(
							S(P(4, 1, 5), P(6, 1, 7)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_CONSTANT, "Foo"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(6, 1, 7)), "unexpected PUBLIC_CONSTANT, expected an identifier as the name of the declared value"),
			},
		},
		"can have an initialiser without a type": {
			input: "val foo = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewIntLiteralNode(S(P(10, 1, 11), P(10, 1, 11)), "5"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "val foo =\n5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 2, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 2, 1)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(10, 2, 1)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewIntLiteralNode(S(P(10, 2, 1), P(10, 2, 1)), "5"),
						),
					),
				},
			),
		},
		"can have an initialiser with a type": {
			input: "val foo: Int = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
							ast.NewIntLiteralNode(S(P(15, 1, 16), P(15, 1, 16)), "5"),
						),
					),
				},
			),
		},
		"can have a type": {
			input: "val foo: Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable type": {
			input: "val foo: Int?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewNilableTypeNode(
								S(P(9, 1, 10), P(12, 1, 13)),
								ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a union type": {
			input: "val foo: Int | String",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								S(P(9, 1, 10), P(20, 1, 21)),
								T(S(P(13, 1, 14), P(13, 1, 14)), token.OR),
								ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
								ast.NewPublicConstantNode(S(P(15, 1, 16), P(20, 1, 21)), "String"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nested union type": {
			input: "val foo: Int | String | Symbol",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								S(P(9, 1, 10), P(29, 1, 30)),
								T(S(P(22, 1, 23), P(22, 1, 23)), token.OR),
								ast.NewBinaryTypeExpressionNode(
									S(P(9, 1, 10), P(20, 1, 21)),
									T(S(P(13, 1, 14), P(13, 1, 14)), token.OR),
									ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
									ast.NewPublicConstantNode(S(P(15, 1, 16), P(20, 1, 21)), "String"),
								),
								ast.NewPublicConstantNode(S(P(24, 1, 25), P(29, 1, 30)), "Symbol"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable union type": {
			input: "val foo: (Int | String)?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(23, 1, 24)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(23, 1, 24)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewNilableTypeNode(
								S(P(10, 1, 11), P(23, 1, 24)),
								ast.NewBinaryTypeExpressionNode(
									S(P(10, 1, 11), P(21, 1, 22)),
									T(S(P(14, 1, 15), P(14, 1, 15)), token.OR),
									ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Int"),
									ast.NewPublicConstantNode(S(P(16, 1, 17), P(21, 1, 22)), "String"),
								),
							),
							nil,
						),
					),
				},
			),
		},
		"can have an intersection type": {
			input: "val foo: Int & String",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								S(P(9, 1, 10), P(20, 1, 21)),
								T(S(P(13, 1, 14), P(13, 1, 14)), token.AND),
								ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
								ast.NewPublicConstantNode(S(P(15, 1, 16), P(20, 1, 21)), "String"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nested intersection type": {
			input: "val foo: Int & String & Symbol",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								S(P(9, 1, 10), P(29, 1, 30)),
								T(S(P(22, 1, 23), P(22, 1, 23)), token.AND),
								ast.NewBinaryTypeExpressionNode(
									S(P(9, 1, 10), P(20, 1, 21)),
									T(S(P(13, 1, 14), P(13, 1, 14)), token.AND),
									ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
									ast.NewPublicConstantNode(S(P(15, 1, 16), P(20, 1, 21)), "String"),
								),
								ast.NewPublicConstantNode(S(P(24, 1, 25), P(29, 1, 30)), "Symbol"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable intersection type": {
			input: "val foo: (Int & String)?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(23, 1, 24)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(23, 1, 24)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewNilableTypeNode(
								S(P(10, 1, 11), P(23, 1, 24)),
								ast.NewBinaryTypeExpressionNode(
									S(P(10, 1, 11), P(21, 1, 22)),
									T(S(P(14, 1, 15), P(14, 1, 15)), token.AND),
									ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Int"),
									ast.NewPublicConstantNode(S(P(16, 1, 17), P(21, 1, 22)), "String"),
								),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a generic type": {
			input: "val foo: Std::Map[Std::Symbol, List[String]]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(43, 1, 44)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(43, 1, 44)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(43, 1, 44)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewGenericConstantNode(
								S(P(9, 1, 10), P(43, 1, 44)),
								ast.NewConstantLookupNode(
									S(P(9, 1, 10), P(16, 1, 17)),
									ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Std"),
									ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Map"),
								),
								[]ast.ComplexConstantNode{
									ast.NewConstantLookupNode(
										S(P(18, 1, 19), P(28, 1, 29)),
										ast.NewPublicConstantNode(S(P(18, 1, 19), P(20, 1, 21)), "Std"),
										ast.NewPublicConstantNode(S(P(23, 1, 24), P(28, 1, 29)), "Symbol"),
									),
									ast.NewGenericConstantNode(
										S(P(31, 1, 32), P(42, 1, 43)),
										ast.NewPublicConstantNode(S(P(31, 1, 32), P(34, 1, 35)), "List"),
										[]ast.ComplexConstantNode{
											ast.NewPublicConstantNode(S(P(36, 1, 37), P(41, 1, 42)), "String"),
										},
									),
								},
							),
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

func TestVariableDeclaration(t *testing.T) {
	tests := testTable{
		"is valid without a type or initialiser": {
			input: "var foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be a part of an expression": {
			input: "a = var foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
							ast.NewVariableDeclarationNode(
								S(P(4, 1, 5), P(10, 1, 11)),
								V(S(P(8, 1, 9), P(10, 1, 11)), token.PUBLIC_IDENTIFIER, "foo"),
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have a private identifier as the variable name": {
			input: "var _foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							V(S(P(4, 1, 5), P(7, 1, 8)), token.PRIVATE_IDENTIFIER, "_foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have an instance variable as the variable name": {
			input: "var @foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							V(S(P(4, 1, 5), P(7, 1, 8)), token.INSTANCE_VARIABLE, "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"instance variables cannot be initialised": {
			input: "var @foo = 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							V(S(P(4, 1, 5), P(7, 1, 8)), token.INSTANCE_VARIABLE, "foo"),
							nil,
							ast.NewIntLiteralNode(
								S(P(11, 1, 12), P(11, 1, 12)),
								"2",
							),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(11, 1, 12), P(11, 1, 12)), "instance variables cannot be initialised when declared"),
			},
		},
		"cannot have a constant as the variable name": {
			input: "var Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(4, 1, 5), P(6, 1, 7)),
						ast.NewInvalidNode(
							S(P(4, 1, 5), P(6, 1, 7)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_CONSTANT, "Foo"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(6, 1, 7)), "unexpected PUBLIC_CONSTANT, expected an identifier as the name of the declared variable"),
			},
		},
		"can have an initialiser without a type": {
			input: "var foo = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewIntLiteralNode(S(P(10, 1, 11), P(10, 1, 11)), "5"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "var foo =\n5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 2, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 2, 1)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(10, 2, 1)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewIntLiteralNode(S(P(10, 2, 1), P(10, 2, 1)), "5"),
						),
					),
				},
			),
		},
		"can have an initialiser with a type": {
			input: "var foo: Int = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
							ast.NewIntLiteralNode(S(P(15, 1, 16), P(15, 1, 16)), "5"),
						),
					),
				},
			),
		},
		"can have a type": {
			input: "var foo: Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable type": {
			input: "var foo: Int?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewNilableTypeNode(
								S(P(9, 1, 10), P(12, 1, 13)),
								ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a union type": {
			input: "var foo: Int | String",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								S(P(9, 1, 10), P(20, 1, 21)),
								T(S(P(13, 1, 14), P(13, 1, 14)), token.OR),
								ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
								ast.NewPublicConstantNode(S(P(15, 1, 16), P(20, 1, 21)), "String"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nested union type": {
			input: "var foo: Int | String | Symbol",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								S(P(9, 1, 10), P(29, 1, 30)),
								T(S(P(22, 1, 23), P(22, 1, 23)), token.OR),
								ast.NewBinaryTypeExpressionNode(
									S(P(9, 1, 10), P(20, 1, 21)),
									T(S(P(13, 1, 14), P(13, 1, 14)), token.OR),
									ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
									ast.NewPublicConstantNode(S(P(15, 1, 16), P(20, 1, 21)), "String"),
								),
								ast.NewPublicConstantNode(S(P(24, 1, 25), P(29, 1, 30)), "Symbol"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable union type": {
			input: "var foo: (Int | String)?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(23, 1, 24)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(23, 1, 24)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewNilableTypeNode(
								S(P(10, 1, 11), P(23, 1, 24)),
								ast.NewBinaryTypeExpressionNode(
									S(P(10, 1, 11), P(21, 1, 22)),
									T(S(P(14, 1, 15), P(14, 1, 15)), token.OR),
									ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Int"),
									ast.NewPublicConstantNode(S(P(16, 1, 17), P(21, 1, 22)), "String"),
								),
							),
							nil,
						),
					),
				},
			),
		},
		"can have an intersection type": {
			input: "var foo: Int & String",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								S(P(9, 1, 10), P(20, 1, 21)),
								T(S(P(13, 1, 14), P(13, 1, 14)), token.AND),
								ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
								ast.NewPublicConstantNode(S(P(15, 1, 16), P(20, 1, 21)), "String"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nested intersection type": {
			input: "var foo: Int & String & Symbol",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								S(P(9, 1, 10), P(29, 1, 30)),
								T(S(P(22, 1, 23), P(22, 1, 23)), token.AND),
								ast.NewBinaryTypeExpressionNode(
									S(P(9, 1, 10), P(20, 1, 21)),
									T(S(P(13, 1, 14), P(13, 1, 14)), token.AND),
									ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
									ast.NewPublicConstantNode(S(P(15, 1, 16), P(20, 1, 21)), "String"),
								),
								ast.NewPublicConstantNode(S(P(24, 1, 25), P(29, 1, 30)), "Symbol"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable intersection type": {
			input: "var foo: (Int & String)?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(23, 1, 24)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(23, 1, 24)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewNilableTypeNode(
								S(P(10, 1, 11), P(23, 1, 24)),
								ast.NewBinaryTypeExpressionNode(
									S(P(10, 1, 11), P(21, 1, 22)),
									T(S(P(14, 1, 15), P(14, 1, 15)), token.AND),
									ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Int"),
									ast.NewPublicConstantNode(S(P(16, 1, 17), P(21, 1, 22)), "String"),
								),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a generic type": {
			input: "var foo: Std::Map[Std::Symbol, List[String]]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(43, 1, 44)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(43, 1, 44)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(43, 1, 44)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewGenericConstantNode(
								S(P(9, 1, 10), P(43, 1, 44)),
								ast.NewConstantLookupNode(
									S(P(9, 1, 10), P(16, 1, 17)),
									ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Std"),
									ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Map"),
								),
								[]ast.ComplexConstantNode{
									ast.NewConstantLookupNode(
										S(P(18, 1, 19), P(28, 1, 29)),
										ast.NewPublicConstantNode(S(P(18, 1, 19), P(20, 1, 21)), "Std"),
										ast.NewPublicConstantNode(S(P(23, 1, 24), P(28, 1, 29)), "Symbol"),
									),
									ast.NewGenericConstantNode(
										S(P(31, 1, 32), P(42, 1, 43)),
										ast.NewPublicConstantNode(S(P(31, 1, 32), P(34, 1, 35)), "List"),
										[]ast.ComplexConstantNode{
											ast.NewPublicConstantNode(S(P(36, 1, 37), P(41, 1, 42)), "String"),
										},
									),
								},
							),
							nil,
						),
					),
				},
			),
		},
		"can have a singleton type": {
			input: "var foo: &Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewSingletonTypeNode(
								S(P(9, 1, 10), P(12, 1, 13)),
								ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Int"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable singleton type": {
			input: "var foo: &Int?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewNilableTypeNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								ast.NewSingletonTypeNode(
									S(P(9, 1, 10), P(12, 1, 13)),
									ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Int"),
								),
							),
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

func TestConstantDeclaration(t *testing.T) {
	tests := testTable{
		"is not valid without an initialiser": {
			input: "const Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							V(S(P(6, 1, 7), P(8, 1, 9)), token.PUBLIC_CONSTANT, "Foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(8, 1, 9)), "constants must be initialised"),
			},
		},
		"can be a part of an expression": {
			input: "a = const _Foo = 'bar'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
							ast.NewConstantDeclarationNode(
								S(P(4, 1, 5), P(21, 1, 22)),
								V(S(P(10, 1, 11), P(13, 1, 14)), token.PRIVATE_CONSTANT, "_Foo"),
								nil,
								ast.NewRawStringLiteralNode(
									S(P(17, 1, 18), P(21, 1, 22)),
									"bar",
								),
							),
						),
					),
				},
			),
		},
		"can have a private constant as the name": {
			input: "const _Foo = 'bar'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							V(S(P(6, 1, 7), P(9, 1, 10)), token.PRIVATE_CONSTANT, "_Foo"),
							nil,
							ast.NewRawStringLiteralNode(
								S(P(13, 1, 14), P(17, 1, 18)),
								"bar",
							),
						),
					),
				},
			),
		},
		"cannot have an instance variable as the name": {
			input: "const @foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(6, 1, 7), P(9, 1, 10)),
						ast.NewInvalidNode(
							S(P(6, 1, 7), P(9, 1, 10)),
							V(S(P(6, 1, 7), P(9, 1, 10)), token.INSTANCE_VARIABLE, "foo"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(6, 1, 7), P(9, 1, 10)), "unexpected INSTANCE_VARIABLE, expected an uppercase identifier as the name of the declared constant"),
			},
		},
		"cannot have a lowercase identifier as the name": {
			input: "const foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(6, 1, 7), P(8, 1, 9)),
						ast.NewInvalidNode(
							S(P(6, 1, 7), P(8, 1, 9)),
							V(S(P(6, 1, 7), P(8, 1, 9)), token.PUBLIC_IDENTIFIER, "foo"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(6, 1, 7), P(8, 1, 9)), "unexpected PUBLIC_IDENTIFIER, expected an uppercase identifier as the name of the declared constant"),
			},
		},
		"can have an initialiser without a type": {
			input: "const Foo = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							V(S(P(6, 1, 7), P(8, 1, 9)), token.PUBLIC_CONSTANT, "Foo"),
							nil,
							ast.NewIntLiteralNode(S(P(12, 1, 13), P(12, 1, 13)), "5"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "const Foo =\n5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 2, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 2, 1)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(12, 2, 1)),
							V(S(P(6, 1, 7), P(8, 1, 9)), token.PUBLIC_CONSTANT, "Foo"),
							nil,
							ast.NewIntLiteralNode(S(P(12, 2, 1), P(12, 2, 1)), "5"),
						),
					),
				},
			),
		},
		"can have an initialiser with a type": {
			input: "const Foo: Int = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							V(S(P(6, 1, 7), P(8, 1, 9)), token.PUBLIC_CONSTANT, "Foo"),
							ast.NewPublicConstantNode(S(P(11, 1, 12), P(13, 1, 14)), "Int"),
							ast.NewIntLiteralNode(S(P(17, 1, 18), P(17, 1, 18)), "5"),
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

func TestTypeDefinition(t *testing.T) {
	tests := testTable{
		"is not valid without an initialiser": {
			input: "typedef Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(11, 1, 12), P(10, 1, 11)),
						ast.NewInvalidNode(
							S(P(11, 1, 12), P(10, 1, 11)),
							T(S(P(11, 1, 12), P(10, 1, 11)), token.END_OF_FILE),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(11, 1, 12), P(10, 1, 11)), "unexpected END_OF_FILE, expected ="),
			},
		},
		"can be a part of an expression": {
			input: "a = typedef Foo = String?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
							ast.NewTypeDefinitionNode(
								S(P(4, 1, 5), P(24, 1, 25)),
								ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Foo"),
								ast.NewNilableTypeNode(
									S(P(18, 1, 19), P(24, 1, 25)),
									ast.NewPublicConstantNode(
										S(P(18, 1, 19), P(23, 1, 24)),
										"String",
									),
								),
							),
						),
					),
				},
			),
		},
		"can have a public constant as the name": {
			input: "typedef Foo = String?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewTypeDefinitionNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							ast.NewPublicConstantNode(S(P(8, 1, 9), P(10, 1, 11)), "Foo"),
							ast.NewNilableTypeNode(
								S(P(14, 1, 15), P(20, 1, 21)),
								ast.NewPublicConstantNode(
									S(P(14, 1, 15), P(19, 1, 20)),
									"String",
								),
							),
						),
					),
				},
			),
		},
		"can have newlines after the assignment operator": {
			input: "typedef Foo =\nString?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 2, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 2, 7)),
						ast.NewTypeDefinitionNode(
							S(P(0, 1, 1), P(20, 2, 7)),
							ast.NewPublicConstantNode(S(P(8, 1, 9), P(10, 1, 11)), "Foo"),
							ast.NewNilableTypeNode(
								S(P(14, 2, 1), P(20, 2, 7)),
								ast.NewPublicConstantNode(
									S(P(14, 2, 1), P(19, 2, 6)),
									"String",
								),
							),
						),
					),
				},
			),
		},
		"can have a private constant as the name": {
			input: "typedef _Foo = String?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewTypeDefinitionNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							ast.NewPrivateConstantNode(S(P(8, 1, 9), P(11, 1, 12)), "_Foo"),
							ast.NewNilableTypeNode(
								S(P(15, 1, 16), P(21, 1, 22)),
								ast.NewPublicConstantNode(
									S(P(15, 1, 16), P(20, 1, 21)),
									"String",
								),
							),
						),
					),
				},
			),
		},
		"cannot have an instance variable as the name": {
			input: "typedef @foo = Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewTypeDefinitionNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							ast.NewInvalidNode(
								S(P(8, 1, 9), P(11, 1, 12)),
								V(S(P(8, 1, 9), P(11, 1, 12)), token.INSTANCE_VARIABLE, "foo"),
							),
							ast.NewPublicConstantNode(
								S(P(15, 1, 16), P(17, 1, 18)),
								"Int",
							),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(8, 1, 9), P(11, 1, 12)), "unexpected INSTANCE_VARIABLE, expected a constant"),
			},
		},
		"cannot have a lowercase identifier as the name": {
			input: "typedef foo = Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewTypeDefinitionNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							ast.NewInvalidNode(
								S(P(8, 1, 9), P(10, 1, 11)),
								V(S(P(8, 1, 9), P(10, 1, 11)), token.PUBLIC_IDENTIFIER, "foo"),
							),
							ast.NewPublicConstantNode(
								S(P(14, 1, 15), P(16, 1, 17)),
								"Int",
							),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(8, 1, 9), P(10, 1, 11)), "unexpected PUBLIC_IDENTIFIER, expected a constant"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestGetterDeclaration(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "a = getter foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
							ast.NewGetterDeclarationNode(
								S(P(4, 1, 5), P(13, 1, 14)),
								[]ast.ParameterNode{
									ast.NewAttributeParameterNode(
										S(P(11, 1, 12), P(13, 1, 14)),
										"foo",
										nil,
									),
								},
							),
						),
					),
				},
			),
		},
		"can have a type": {
			input: "getter foo: Bar?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewGetterDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(7, 1, 8), P(15, 1, 16)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(12, 1, 13), P(15, 1, 16)),
										ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a few attributes": {
			input: "getter foo: Bar?, bar, baz: Int | Float",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 1, 39)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(38, 1, 39)),
						ast.NewGetterDeclarationNode(
							S(P(0, 1, 1), P(38, 1, 39)),
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(7, 1, 8), P(15, 1, 16)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(12, 1, 13), P(15, 1, 16)),
										ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
									),
								),
								ast.NewAttributeParameterNode(
									S(P(18, 1, 19), P(20, 1, 21)),
									"bar",
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(23, 1, 24), P(38, 1, 39)),
									"baz",
									ast.NewBinaryTypeExpressionNode(
										S(P(28, 1, 29), P(38, 1, 39)),
										T(S(P(32, 1, 33), P(32, 1, 33)), token.OR),
										ast.NewPublicConstantNode(S(P(28, 1, 29), P(30, 1, 31)), "Int"),
										ast.NewPublicConstantNode(S(P(34, 1, 35), P(38, 1, 39)), "Float"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: `
				getter foo: Bar?,
							 bar,
							 baz: Int | Float
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(60, 4, 25)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(60, 4, 25)),
						ast.NewGetterDeclarationNode(
							S(P(5, 2, 5), P(59, 4, 24)),
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(12, 2, 12), P(20, 2, 20)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(17, 2, 17), P(20, 2, 20)),
										ast.NewPublicConstantNode(S(P(17, 2, 17), P(19, 2, 19)), "Bar"),
									),
								),
								ast.NewAttributeParameterNode(
									S(P(31, 3, 9), P(33, 3, 11)),
									"bar",
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(44, 4, 9), P(59, 4, 24)),
									"baz",
									ast.NewBinaryTypeExpressionNode(
										S(P(49, 4, 14), P(59, 4, 24)),
										T(S(P(53, 4, 18), P(53, 4, 18)), token.OR),
										ast.NewPublicConstantNode(S(P(49, 4, 14), P(51, 4, 16)), "Int"),
										ast.NewPublicConstantNode(S(P(55, 4, 20), P(59, 4, 24)), "Float"),
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

func TestSetterDeclaration(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "a = setter foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
							ast.NewSetterDeclarationNode(
								S(P(4, 1, 5), P(13, 1, 14)),
								[]ast.ParameterNode{
									ast.NewAttributeParameterNode(
										S(P(11, 1, 12), P(13, 1, 14)),
										"foo",
										nil,
									),
								},
							),
						),
					),
				},
			),
		},
		"can have a type": {
			input: "setter foo: Bar?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewSetterDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(7, 1, 8), P(15, 1, 16)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(12, 1, 13), P(15, 1, 16)),
										ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a few attributes": {
			input: "setter foo: Bar?, bar, baz: Int | Float",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 1, 39)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(38, 1, 39)),
						ast.NewSetterDeclarationNode(
							S(P(0, 1, 1), P(38, 1, 39)),
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(7, 1, 8), P(15, 1, 16)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(12, 1, 13), P(15, 1, 16)),
										ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
									),
								),
								ast.NewAttributeParameterNode(
									S(P(18, 1, 19), P(20, 1, 21)),
									"bar",
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(23, 1, 24), P(38, 1, 39)),
									"baz",
									ast.NewBinaryTypeExpressionNode(
										S(P(28, 1, 29), P(38, 1, 39)),
										T(S(P(32, 1, 33), P(32, 1, 33)), token.OR),
										ast.NewPublicConstantNode(S(P(28, 1, 29), P(30, 1, 31)), "Int"),
										ast.NewPublicConstantNode(S(P(34, 1, 35), P(38, 1, 39)), "Float"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: `
				setter foo: Bar?,
							 bar,
							 baz: Int | Float
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(60, 4, 25)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(60, 4, 25)),
						ast.NewSetterDeclarationNode(
							S(P(5, 2, 5), P(59, 4, 24)),
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(12, 2, 12), P(20, 2, 20)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(17, 2, 17), P(20, 2, 20)),
										ast.NewPublicConstantNode(S(P(17, 2, 17), P(19, 2, 19)), "Bar"),
									),
								),
								ast.NewAttributeParameterNode(
									S(P(31, 3, 9), P(33, 3, 11)),
									"bar",
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(44, 4, 9), P(59, 4, 24)),
									"baz",
									ast.NewBinaryTypeExpressionNode(
										S(P(49, 4, 14), P(59, 4, 24)),
										T(S(P(53, 4, 18), P(53, 4, 18)), token.OR),
										ast.NewPublicConstantNode(S(P(49, 4, 14), P(51, 4, 16)), "Int"),
										ast.NewPublicConstantNode(S(P(55, 4, 20), P(59, 4, 24)), "Float"),
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

func TestAccessorDeclaration(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "a = accessor foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
							ast.NewAccessorDeclarationNode(
								S(P(4, 1, 5), P(15, 1, 16)),
								[]ast.ParameterNode{
									ast.NewAttributeParameterNode(
										S(P(13, 1, 14), P(15, 1, 16)),
										"foo",
										nil,
									),
								},
							),
						),
					),
				},
			),
		},
		"can have a type": {
			input: "accessor foo: Bar?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewAccessorDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(9, 1, 10), P(17, 1, 18)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(14, 1, 15), P(17, 1, 18)),
										ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Bar"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a few attributes": {
			input: "accessor foo: Bar?, bar, baz: Int | Float",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(40, 1, 41)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 1, 41)),
						ast.NewAccessorDeclarationNode(
							S(P(0, 1, 1), P(40, 1, 41)),
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(9, 1, 10), P(17, 1, 18)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(14, 1, 15), P(17, 1, 18)),
										ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Bar"),
									),
								),
								ast.NewAttributeParameterNode(
									S(P(20, 1, 21), P(22, 1, 23)),
									"bar",
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(25, 1, 26), P(40, 1, 41)),
									"baz",
									ast.NewBinaryTypeExpressionNode(
										S(P(30, 1, 31), P(40, 1, 41)),
										T(S(P(34, 1, 35), P(34, 1, 35)), token.OR),
										ast.NewPublicConstantNode(S(P(30, 1, 31), P(32, 1, 33)), "Int"),
										ast.NewPublicConstantNode(S(P(36, 1, 37), P(40, 1, 41)), "Float"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: `
				accessor foo: Bar?,
							 bar,
							 baz: Int | Float
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(62, 4, 25)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(62, 4, 25)),
						ast.NewAccessorDeclarationNode(
							S(P(5, 2, 5), P(61, 4, 24)),
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(14, 2, 14), P(22, 2, 22)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(19, 2, 19), P(22, 2, 22)),
										ast.NewPublicConstantNode(S(P(19, 2, 19), P(21, 2, 21)), "Bar"),
									),
								),
								ast.NewAttributeParameterNode(
									S(P(33, 3, 9), P(35, 3, 11)),
									"bar",
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(46, 4, 9), P(61, 4, 24)),
									"baz",
									ast.NewBinaryTypeExpressionNode(
										S(P(51, 4, 14), P(61, 4, 24)),
										T(S(P(55, 4, 18), P(55, 4, 18)), token.OR),
										ast.NewPublicConstantNode(S(P(51, 4, 14), P(53, 4, 16)), "Int"),
										ast.NewPublicConstantNode(S(P(57, 4, 20), P(61, 4, 24)), "Float"),
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

func TestAliasDeclaration(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "a = alias foo bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
							ast.NewAliasDeclarationNode(
								S(P(4, 1, 5), P(16, 1, 17)),
								[]*ast.AliasDeclarationEntry{
									ast.NewAliasDeclarationEntry(
										S(P(10, 1, 11), P(16, 1, 17)),
										"foo",
										"bar",
									),
								},
							),
						),
					),
				},
			),
		},
		"can have a few entries": {
			input: "alias foo bar, add plus, remove delete",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(37, 1, 38)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(37, 1, 38)),
						ast.NewAliasDeclarationNode(
							S(P(0, 1, 1), P(37, 1, 38)),
							[]*ast.AliasDeclarationEntry{
								ast.NewAliasDeclarationEntry(
									S(P(6, 1, 7), P(12, 1, 13)),
									"foo",
									"bar",
								),
								ast.NewAliasDeclarationEntry(
									S(P(15, 1, 16), P(22, 1, 23)),
									"add",
									"plus",
								),
								ast.NewAliasDeclarationEntry(
									S(P(25, 1, 26), P(37, 1, 38)),
									"remove",
									"delete",
								),
							},
						),
					),
				},
			),
		},
		"can have public identifiers as names": {
			input: "alias foo bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewAliasDeclarationNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							[]*ast.AliasDeclarationEntry{
								ast.NewAliasDeclarationEntry(
									S(P(6, 1, 7), P(12, 1, 13)),
									"foo",
									"bar",
								),
							},
						),
					),
				},
			),
		},
		"can have overridable operators as names": {
			input: "alias + -",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewAliasDeclarationNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							[]*ast.AliasDeclarationEntry{
								ast.NewAliasDeclarationEntry(
									S(P(6, 1, 7), P(8, 1, 9)),
									"+",
									"-",
								),
							},
						),
					),
				},
			),
		},
		"can have setters as names": {
			input: "alias foo= bar=",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewAliasDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							[]*ast.AliasDeclarationEntry{
								ast.NewAliasDeclarationEntry(
									S(P(6, 1, 7), P(14, 1, 15)),
									"foo=",
									"bar=",
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: "alias\nfoo\nbar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 3, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 3, 3)),
						ast.NewAliasDeclarationNode(
							S(P(0, 1, 1), P(12, 3, 3)),
							[]*ast.AliasDeclarationEntry{
								ast.NewAliasDeclarationEntry(
									S(P(6, 2, 1), P(12, 3, 3)),
									"foo",
									"bar",
								),
							},
						),
					),
				},
			),
		},
		"can have private identifiers as names": {
			input: "alias _foo _bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewAliasDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							[]*ast.AliasDeclarationEntry{
								ast.NewAliasDeclarationEntry(
									S(P(6, 1, 7), P(14, 1, 15)),
									"_foo",
									"_bar",
								),
							},
						),
					),
				},
			),
		},
		"cannot have instance variables as names": {
			input: "alias @foo @bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewAliasDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							[]*ast.AliasDeclarationEntry{
								ast.NewAliasDeclarationEntry(
									S(P(6, 1, 7), P(14, 1, 15)),
									"foo",
									"bar",
								),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(6, 1, 7), P(9, 1, 10)), "unexpected INSTANCE_VARIABLE, expected a method name (identifier, overridable operator)"),
				errors.NewError(L("main", P(11, 1, 12), P(14, 1, 15)), "unexpected INSTANCE_VARIABLE, expected a method name (identifier, overridable operator)"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestClassDeclaration(t *testing.T) {
	tests := testTable{
		"can be anonymous": {
			input: `class; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							false,
							false,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be a part of an expression": {
			input: `foo = class; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewClassDeclarationNode(
								S(P(6, 1, 7), P(15, 1, 16)),
								false,
								false,
								nil,
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can be anonymous with a superclass": {
			input: `class < Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							false,
							false,
							nil,
							nil,
							ast.NewPublicConstantNode(S(P(8, 1, 9), P(10, 1, 11)), "Foo"),
							nil,
						),
					),
				},
			),
		},
		"can have type variables": {
			input: `class Foo[V, +T, -Z]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							[]ast.TypeVariableNode{
								ast.NewVariantTypeVariableNode(
									S(P(10, 1, 11), P(10, 1, 11)),
									ast.INVARIANT,
									"V",
									nil,
								),
								ast.NewVariantTypeVariableNode(
									S(P(13, 1, 14), P(14, 1, 15)),
									ast.COVARIANT,
									"T",
									nil,
								),
								ast.NewVariantTypeVariableNode(
									S(P(17, 1, 18), P(18, 1, 19)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have type variables with upper bounds": {
			input: `class Foo[V < Std::String, +T < Foo, -Z < _Bar]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(51, 1, 52)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(51, 1, 52)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(51, 1, 52)),
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							[]ast.TypeVariableNode{
								ast.NewVariantTypeVariableNode(
									S(P(10, 1, 11), P(24, 1, 25)),
									ast.INVARIANT,
									"V",
									ast.NewConstantLookupNode(
										S(P(14, 1, 15), P(24, 1, 25)),
										ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Std"),
										ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									),
								),
								ast.NewVariantTypeVariableNode(
									S(P(27, 1, 28), P(34, 1, 35)),
									ast.COVARIANT,
									"T",
									ast.NewPublicConstantNode(S(P(32, 1, 33), P(34, 1, 35)), "Foo"),
								),
								ast.NewVariantTypeVariableNode(
									S(P(37, 1, 38), P(45, 1, 46)),
									ast.CONTRAVARIANT,
									"Z",
									ast.NewPrivateConstantNode(S(P(42, 1, 43), P(45, 1, 46)), "_Bar"),
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have an empty type variable list": {
			input: `class Foo[]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(10, 1, 11), P(10, 1, 11)), "unexpected ], expected a list of type variables"),
			},
		},
		"can be abstract": {
			input: `abstract class Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							true,
							false,
							ast.NewPublicConstantNode(S(P(15, 1, 16), P(17, 1, 18)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot repeat abstract": {
			input: `abstract abstract class Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 1, 32)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 1, 32)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(31, 1, 32)),
							true,
							false,
							ast.NewPublicConstantNode(S(P(24, 1, 25), P(26, 1, 27)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(7, 1, 8)), "the abstract modifier can only be attached once"),
			},
		},
		"cannot attach abstract to a sealed class": {
			input: `abstract sealed class Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							true,
							true,
							ast.NewPublicConstantNode(S(P(22, 1, 23), P(24, 1, 25)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(7, 1, 8)), "the abstract modifier cannot be attached to sealed classes"),
			},
		},
		"can be sealed": {
			input: `sealed class Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							false,
							true,
							ast.NewPublicConstantNode(S(P(13, 1, 14), P(15, 1, 16)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot attach sealed to abstract classes": {
			input: `sealed abstract class Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							true,
							true,
							ast.NewPublicConstantNode(S(P(22, 1, 23), P(24, 1, 25)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(5, 1, 6)), "the sealed modifier cannot be attached to abstract classes"),
			},
		},
		"cannot repeat sealed": {
			input: `sealed sealed class Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							false,
							true,
							ast.NewPublicConstantNode(S(P(20, 1, 21), P(22, 1, 23)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(0, 1, 1), P(5, 1, 6)), "the sealed modifier can only be attached once"),
			},
		},
		"can have a public constant as a name": {
			input: `class Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a private constant as a name": {
			input: `class _Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							false,
							false,
							ast.NewPrivateConstantNode(S(P(6, 1, 7), P(9, 1, 10)), "_Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a constant lookup as a name": {
			input: `class Foo::Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							false,
							false,
							ast.NewConstantLookupNode(
								S(P(6, 1, 7), P(13, 1, 14)),
								ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
								ast.NewPublicConstantNode(S(P(11, 1, 12), P(13, 1, 14)), "Bar"),
							),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have an identifier as a name": {
			input: `class foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							false,
							false,
							ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(6, 1, 7), P(8, 1, 9)), "invalid class name, expected a constant"),
			},
		},
		"can have a public constant as a superclass": {
			input: `class Foo < Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
							nil,
						),
					),
				},
			),
		},
		"can have a private constant as a superclass": {
			input: `class Foo < _Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							ast.NewPrivateConstantNode(S(P(12, 1, 13), P(15, 1, 16)), "_Bar"),
							nil,
						),
					),
				},
			),
		},
		"can have a constant lookup as a superclass": {
			input: `class Foo < Bar::Baz; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							ast.NewConstantLookupNode(
								S(P(12, 1, 13), P(19, 1, 20)),
								ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
								ast.NewPublicConstantNode(S(P(17, 1, 18), P(19, 1, 20)), "Baz"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a generic constant as a superclass": {
			input: `class Foo < Std::Map[Symbol, String]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(40, 1, 41)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 1, 41)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(40, 1, 41)),
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							ast.NewGenericConstantNode(
								S(P(12, 1, 13), P(35, 1, 36)),
								ast.NewConstantLookupNode(
									S(P(12, 1, 13), P(19, 1, 20)),
									ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Std"),
									ast.NewPublicConstantNode(S(P(17, 1, 18), P(19, 1, 20)), "Map"),
								),
								[]ast.ComplexConstantNode{
									ast.NewPublicConstantNode(S(P(21, 1, 22), P(26, 1, 27)), "Symbol"),
									ast.NewPublicConstantNode(S(P(29, 1, 30), P(34, 1, 35)), "String"),
								},
							),
							nil,
						),
					),
				},
			),
		},
		"cannot have an identifier as a superclass": {
			input: `class Foo < bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							ast.NewInvalidNode(S(P(12, 1, 13), P(14, 1, 15)), V(S(P(12, 1, 13), P(14, 1, 15)), token.PUBLIC_IDENTIFIER, "bar")),
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(12, 1, 13), P(14, 1, 15)), "unexpected PUBLIC_IDENTIFIER, expected a constant"),
			},
		},
		"can have a multiline body": {
			input: `class Foo
	foo = 2
	nil
end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 4, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 4, 3)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(26, 4, 3)),
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(11, 2, 2), P(18, 2, 9)),
									ast.NewAssignmentExpressionNode(
										S(P(11, 2, 2), P(17, 2, 8)),
										T(S(P(15, 2, 6), P(15, 2, 6)), token.EQUAL_OP),
										ast.NewPublicIdentifierNode(S(P(11, 2, 2), P(13, 2, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(17, 2, 8), P(17, 2, 8)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(20, 3, 2), P(23, 3, 5)),
									ast.NewNilLiteralNode(S(P(20, 3, 2), P(22, 3, 4))),
								),
							},
						),
					),
				},
			),
		},
		"can have a single line body with then": {
			input: `class Foo then .1 * .2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(15, 1, 16), P(21, 1, 22)),
									ast.NewBinaryExpressionNode(
										S(P(15, 1, 16), P(21, 1, 22)),
										T(S(P(18, 1, 19), P(18, 1, 19)), token.STAR),
										ast.NewFloatLiteralNode(S(P(15, 1, 16), P(16, 1, 17)), "0.1"),
										ast.NewFloatLiteralNode(S(P(20, 1, 21), P(21, 1, 22)), "0.2"),
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

func TestModuleDeclaration(t *testing.T) {
	tests := testTable{
		"can be anonymous": {
			input: `module; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be a part of an expression": {
			input: `foo = module; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewModuleDeclarationNode(
								S(P(6, 1, 7), P(16, 1, 17)),
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"cannot be generic": {
			input: `module Foo[V, +T, -Z]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(10, 1, 11), P(20, 1, 21)), "modules cannot be generic"),
			},
		},
		"can have a public constant as a name": {
			input: `module Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							nil,
						),
					),
				},
			),
		},
		"can have a private constant as a name": {
			input: `module _Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							ast.NewPrivateConstantNode(S(P(7, 1, 8), P(10, 1, 11)), "_Foo"),
							nil,
						),
					),
				},
			),
		},
		"can have a constant lookup as a name": {
			input: `module Foo::Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewConstantLookupNode(
								S(P(7, 1, 8), P(14, 1, 15)),
								ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
								ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
							),
							nil,
						),
					),
				},
			),
		},
		"cannot have an identifier as a name": {
			input: `module foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(7, 1, 8), P(9, 1, 10)), "invalid module name, expected a constant"),
			},
		},
		"can have a multiline body": {
			input: `module Foo
	foo = 2
	nil
end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 4, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 4, 3)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(27, 4, 3)),
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(12, 2, 2), P(19, 2, 9)),
									ast.NewAssignmentExpressionNode(
										S(P(12, 2, 2), P(18, 2, 8)),
										T(S(P(16, 2, 6), P(16, 2, 6)), token.EQUAL_OP),
										ast.NewPublicIdentifierNode(S(P(12, 2, 2), P(14, 2, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(18, 2, 8), P(18, 2, 8)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(21, 3, 2), P(24, 3, 5)),
									ast.NewNilLiteralNode(S(P(21, 3, 2), P(23, 3, 4))),
								),
							},
						),
					),
				},
			),
		},
		"can have a single line body with then": {
			input: `module Foo then .1 * .2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(16, 1, 17), P(22, 1, 23)),
									ast.NewBinaryExpressionNode(
										S(P(16, 1, 17), P(22, 1, 23)),
										T(S(P(19, 1, 20), P(19, 1, 20)), token.STAR),
										ast.NewFloatLiteralNode(S(P(16, 1, 17), P(17, 1, 18)), "0.1"),
										ast.NewFloatLiteralNode(S(P(21, 1, 22), P(22, 1, 23)), "0.2"),
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

func TestMixinDeclaration(t *testing.T) {
	tests := testTable{
		"can be anonymous": {
			input: `mixin; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be a part of an expression": {
			input: `foo = mixin; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewMixinDeclarationNode(
								S(P(6, 1, 7), P(15, 1, 16)),
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have type variables": {
			input: `mixin Foo[V, +T, -Z]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							[]ast.TypeVariableNode{
								ast.NewVariantTypeVariableNode(
									S(P(10, 1, 11), P(10, 1, 11)),
									ast.INVARIANT,
									"V",
									nil,
								),
								ast.NewVariantTypeVariableNode(
									S(P(13, 1, 14), P(14, 1, 15)),
									ast.COVARIANT,
									"T",
									nil,
								),
								ast.NewVariantTypeVariableNode(
									S(P(17, 1, 18), P(18, 1, 19)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have type variables with upper bounds": {
			input: `mixin Foo[V < Std::String, +T < Foo, -Z < _Bar]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(51, 1, 52)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(51, 1, 52)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(51, 1, 52)),
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							[]ast.TypeVariableNode{
								ast.NewVariantTypeVariableNode(
									S(P(10, 1, 11), P(24, 1, 25)),
									ast.INVARIANT,
									"V",
									ast.NewConstantLookupNode(
										S(P(14, 1, 15), P(24, 1, 25)),
										ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Std"),
										ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									),
								),
								ast.NewVariantTypeVariableNode(
									S(P(27, 1, 28), P(34, 1, 35)),
									ast.COVARIANT,
									"T",
									ast.NewPublicConstantNode(S(P(32, 1, 33), P(34, 1, 35)), "Foo"),
								),
								ast.NewVariantTypeVariableNode(
									S(P(37, 1, 38), P(45, 1, 46)),
									ast.CONTRAVARIANT,
									"Z",
									ast.NewPrivateConstantNode(S(P(42, 1, 43), P(45, 1, 46)), "_Bar"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot have an empty type variable list": {
			input: `mixin Foo[]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(10, 1, 11), P(10, 1, 11)), "unexpected ], expected a list of type variables"),
			},
		},
		"cannot be abstract": {
			input: `abstract mixin Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(9, 1, 10), P(22, 1, 23)),
						ast.NewMixinDeclarationNode(
							S(P(9, 1, 10), P(22, 1, 23)),
							ast.NewPublicConstantNode(S(P(15, 1, 16), P(17, 1, 18)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(9, 1, 10), P(22, 1, 23)), "the abstract modifier can only be attached to classes"),
			},
		},
		"can have a public constant as a name": {
			input: `mixin Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a private constant as a name": {
			input: `mixin _Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewPrivateConstantNode(S(P(6, 1, 7), P(9, 1, 10)), "_Foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a constant lookup as a name": {
			input: `mixin Foo::Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							ast.NewConstantLookupNode(
								S(P(6, 1, 7), P(13, 1, 14)),
								ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
								ast.NewPublicConstantNode(S(P(11, 1, 12), P(13, 1, 14)), "Bar"),
							),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have an identifier as a name": {
			input: `mixin foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(6, 1, 7), P(8, 1, 9)), "invalid mixin name, expected a constant"),
			},
		},
		"can have a multiline body": {
			input: `mixin Foo
	foo = 2
	nil
end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 4, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 4, 3)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(26, 4, 3)),
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(11, 2, 2), P(18, 2, 9)),
									ast.NewAssignmentExpressionNode(
										S(P(11, 2, 2), P(17, 2, 8)),
										T(S(P(15, 2, 6), P(15, 2, 6)), token.EQUAL_OP),
										ast.NewPublicIdentifierNode(S(P(11, 2, 2), P(13, 2, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(17, 2, 8), P(17, 2, 8)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(20, 3, 2), P(23, 3, 5)),
									ast.NewNilLiteralNode(S(P(20, 3, 2), P(22, 3, 4))),
								),
							},
						),
					),
				},
			),
		},
		"can have a single line body with then": {
			input: `mixin Foo then .1 * .2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(15, 1, 16), P(21, 1, 22)),
									ast.NewBinaryExpressionNode(
										S(P(15, 1, 16), P(21, 1, 22)),
										T(S(P(18, 1, 19), P(18, 1, 19)), token.STAR),
										ast.NewFloatLiteralNode(S(P(15, 1, 16), P(16, 1, 17)), "0.1"),
										ast.NewFloatLiteralNode(S(P(20, 1, 21), P(21, 1, 22)), "0.2"),
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

func TestInterfaceDeclaration(t *testing.T) {
	tests := testTable{
		"can be anonymous": {
			input: `interface; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewInterfaceDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be a part of an expression": {
			input: `foo = interface; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewInterfaceDeclarationNode(
								S(P(6, 1, 7), P(19, 1, 20)),
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have type variables": {
			input: `interface Foo[V, +T, -Z]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 1, 29)),
						ast.NewInterfaceDeclarationNode(
							S(P(0, 1, 1), P(28, 1, 29)),
							ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Foo"),
							[]ast.TypeVariableNode{
								ast.NewVariantTypeVariableNode(
									S(P(14, 1, 15), P(14, 1, 15)),
									ast.INVARIANT,
									"V",
									nil,
								),
								ast.NewVariantTypeVariableNode(
									S(P(17, 1, 18), P(18, 1, 19)),
									ast.COVARIANT,
									"T",
									nil,
								),
								ast.NewVariantTypeVariableNode(
									S(P(21, 1, 22), P(22, 1, 23)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have type variables with upper bounds": {
			input: `interface Foo[V < Std::String, +T < Foo, -Z < _Bar]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(55, 1, 56)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(55, 1, 56)),
						ast.NewInterfaceDeclarationNode(
							S(P(0, 1, 1), P(55, 1, 56)),
							ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Foo"),
							[]ast.TypeVariableNode{
								ast.NewVariantTypeVariableNode(
									S(P(14, 1, 15), P(28, 1, 29)),
									ast.INVARIANT,
									"V",
									ast.NewConstantLookupNode(
										S(P(18, 1, 19), P(28, 1, 29)),
										ast.NewPublicConstantNode(S(P(18, 1, 19), P(20, 1, 21)), "Std"),
										ast.NewPublicConstantNode(S(P(23, 1, 24), P(28, 1, 29)), "String"),
									),
								),
								ast.NewVariantTypeVariableNode(
									S(P(31, 1, 32), P(38, 1, 39)),
									ast.COVARIANT,
									"T",
									ast.NewPublicConstantNode(S(P(36, 1, 37), P(38, 1, 39)), "Foo"),
								),
								ast.NewVariantTypeVariableNode(
									S(P(41, 1, 42), P(49, 1, 50)),
									ast.CONTRAVARIANT,
									"Z",
									ast.NewPrivateConstantNode(S(P(46, 1, 47), P(49, 1, 50)), "_Bar"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot have an empty type variable list": {
			input: `interface Foo[]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewInterfaceDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(14, 1, 15), P(14, 1, 15)), "unexpected ], expected a list of type variables"),
			},
		},
		"can have a public constant as a name": {
			input: `interface Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewInterfaceDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a private constant as a name": {
			input: `interface _Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewInterfaceDeclarationNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							ast.NewPrivateConstantNode(S(P(10, 1, 11), P(13, 1, 14)), "_Foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a constant lookup as a name": {
			input: `interface Foo::Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewInterfaceDeclarationNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							ast.NewConstantLookupNode(
								S(P(10, 1, 11), P(17, 1, 18)),
								ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Foo"),
								ast.NewPublicConstantNode(S(P(15, 1, 16), P(17, 1, 18)), "Bar"),
							),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have an identifier as a name": {
			input: `interface foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewInterfaceDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							ast.NewPublicIdentifierNode(S(P(10, 1, 11), P(12, 1, 13)), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(10, 1, 11), P(12, 1, 13)), "invalid interface name, expected a constant"),
			},
		},
		"can have a multiline body": {
			input: `interface Foo
	foo = 2
	nil
end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 4, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(30, 4, 3)),
						ast.NewInterfaceDeclarationNode(
							S(P(0, 1, 1), P(30, 4, 3)),
							ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Foo"),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(15, 2, 2), P(22, 2, 9)),
									ast.NewAssignmentExpressionNode(
										S(P(15, 2, 2), P(21, 2, 8)),
										T(S(P(19, 2, 6), P(19, 2, 6)), token.EQUAL_OP),
										ast.NewPublicIdentifierNode(S(P(15, 2, 2), P(17, 2, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(21, 2, 8), P(21, 2, 8)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(24, 3, 2), P(27, 3, 5)),
									ast.NewNilLiteralNode(S(P(24, 3, 2), P(26, 3, 4))),
								),
							},
						),
					),
				},
			),
		},
		"can have a single line body with then": {
			input: `interface Foo then .1 * .2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewInterfaceDeclarationNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Foo"),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(19, 1, 20), P(25, 1, 26)),
									ast.NewBinaryExpressionNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										T(S(P(22, 1, 23), P(22, 1, 23)), token.STAR),
										ast.NewFloatLiteralNode(S(P(19, 1, 20), P(20, 1, 21)), "0.1"),
										ast.NewFloatLiteralNode(S(P(24, 1, 25), P(25, 1, 26)), "0.2"),
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

func TestStructDeclaration(t *testing.T) {
	tests := testTable{
		"can be anonymous": {
			input: `struct; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be a part of an expression": {
			input: `foo = struct; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "foo"),
							ast.NewStructDeclarationNode(
								S(P(6, 1, 7), P(16, 1, 17)),
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have type variables": {
			input: `struct Foo[V, +T, -Z]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							[]ast.TypeVariableNode{
								ast.NewVariantTypeVariableNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									ast.INVARIANT,
									"V",
									nil,
								),
								ast.NewVariantTypeVariableNode(
									S(P(14, 1, 15), P(15, 1, 16)),
									ast.COVARIANT,
									"T",
									nil,
								),
								ast.NewVariantTypeVariableNode(
									S(P(18, 1, 19), P(19, 1, 20)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have type variables with upper bounds": {
			input: `struct Foo[V < Std::String, +T < Foo, -Z < _Bar]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(52, 1, 53)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(52, 1, 53)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(52, 1, 53)),
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							[]ast.TypeVariableNode{
								ast.NewVariantTypeVariableNode(
									S(P(11, 1, 12), P(25, 1, 26)),
									ast.INVARIANT,
									"V",
									ast.NewConstantLookupNode(
										S(P(15, 1, 16), P(25, 1, 26)),
										ast.NewPublicConstantNode(S(P(15, 1, 16), P(17, 1, 18)), "Std"),
										ast.NewPublicConstantNode(S(P(20, 1, 21), P(25, 1, 26)), "String"),
									),
								),
								ast.NewVariantTypeVariableNode(
									S(P(28, 1, 29), P(35, 1, 36)),
									ast.COVARIANT,
									"T",
									ast.NewPublicConstantNode(S(P(33, 1, 34), P(35, 1, 36)), "Foo"),
								),
								ast.NewVariantTypeVariableNode(
									S(P(38, 1, 39), P(46, 1, 47)),
									ast.CONTRAVARIANT,
									"Z",
									ast.NewPrivateConstantNode(S(P(43, 1, 44), P(46, 1, 47)), "_Bar"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot have an empty type variable list": {
			input: `struct Foo[]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(11, 1, 12), P(11, 1, 12)), "unexpected ], expected a list of type variables"),
			},
		},
		"can have a public constant as a name": {
			input: `struct Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a private constant as a name": {
			input: `struct _Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							ast.NewPrivateConstantNode(S(P(7, 1, 8), P(10, 1, 11)), "_Foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a constant lookup as a name": {
			input: `struct Foo::Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewConstantLookupNode(
								S(P(7, 1, 8), P(14, 1, 15)),
								ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
								ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
							),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have an identifier as a name": {
			input: `struct foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(7, 1, 8), P(9, 1, 10)), "invalid struct name, expected a constant"),
			},
		},
		"can have a multiline body": {
			input: `struct Foo
  foo
  bar: String?
  baz: Int = .3
  ban = 'hey'
end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(64, 6, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(64, 6, 3)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(64, 6, 3)),
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							nil,
							[]ast.StructBodyStatementNode{
								ast.NewParameterStatementNode(
									S(P(13, 2, 3), P(16, 2, 6)),
									ast.NewFormalParameterNode(
										S(P(13, 2, 3), P(15, 2, 5)),
										"foo",
										nil,
										nil,
										ast.NormalParameterKind,
									),
								),
								ast.NewParameterStatementNode(
									S(P(19, 3, 3), P(31, 3, 15)),
									ast.NewFormalParameterNode(
										S(P(19, 3, 3), P(30, 3, 14)),
										"bar",
										ast.NewNilableTypeNode(
											S(P(24, 3, 8), P(30, 3, 14)),
											ast.NewPublicConstantNode(S(P(24, 3, 8), P(29, 3, 13)), "String"),
										),
										nil,
										ast.NormalParameterKind,
									),
								),
								ast.NewParameterStatementNode(
									S(P(34, 4, 3), P(47, 4, 16)),
									ast.NewFormalParameterNode(
										S(P(34, 4, 3), P(46, 4, 15)),
										"baz",
										ast.NewPublicConstantNode(S(P(39, 4, 8), P(41, 4, 10)), "Int"),
										ast.NewFloatLiteralNode(S(P(45, 4, 14), P(46, 4, 15)), "0.3"),
										ast.NormalParameterKind,
									),
								),
								ast.NewParameterStatementNode(
									S(P(50, 5, 3), P(61, 5, 14)),
									ast.NewFormalParameterNode(
										S(P(50, 5, 3), P(60, 5, 13)),
										"ban",
										nil,
										ast.NewRawStringLiteralNode(S(P(56, 5, 9), P(60, 5, 13)), "hey"),
										ast.NormalParameterKind,
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a single line body with then": {
			input: `struct Foo then foo: Int`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(23, 1, 24)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(23, 1, 24)),
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							nil,
							[]ast.StructBodyStatementNode{
								ast.NewParameterStatementNode(
									S(P(16, 1, 17), P(23, 1, 24)),
									ast.NewFormalParameterNode(
										S(P(16, 1, 17), P(23, 1, 24)),
										"foo",
										ast.NewPublicConstantNode(S(P(21, 1, 22), P(23, 1, 24)), "Int"),
										nil,
										ast.NormalParameterKind,
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

func TestMethodDefinition(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "bar = def foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "bar"),
							ast.NewMethodDefinitionNode(
								S(P(6, 1, 7), P(17, 1, 18)),
								"foo",
								nil,
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have a public identifier as a name": {
			input: "def foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							"foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a setter as a name": {
			input: "def foo=(v); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							"foo=",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(9, 1, 10), P(9, 1, 10)),
									"v",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"setters cannot have custom return types": {
			input: "def foo=(v): String; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(23, 1, 24)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(23, 1, 24)),
							"foo=",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(9, 1, 10), P(9, 1, 10)),
									"v",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							ast.NewPublicConstantNode(
								S(P(13, 1, 14), P(18, 1, 19)),
								"String",
							),
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(13, 1, 14), P(18, 1, 19)), "setter methods cannot be defined with custom return types"),
			},
		},
		"setters cannot have multiple parameters": {
			input: "def foo=(a, b, c); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							"foo=",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(9, 1, 10), P(9, 1, 10)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(12, 1, 13), P(12, 1, 13)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(15, 1, 16), P(15, 1, 16)),
									"c",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(12, 1, 13), P(15, 1, 16)), "setter methods must have a single parameter, got: 3"),
			},
		},
		"setters must have a parameter": {
			input: "def fo=; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							"fo=",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(6, 1, 7)), "setter methods must have a single parameter, got: 0"),
			},
		},
		"can have a private identifier as a name": {
			input: "def _foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							"_foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a keyword as a name": {
			input: "def class; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							"class",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have an overridable operator as a name": {
			input: "def +; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							"+",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have brackets as a name": {
			input: "def []; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							"[]",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have brackets setter as a name": {
			input: "def []=(v); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							"[]=",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"v",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have a public constant as a name": {
			input: "def Foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							"Foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(6, 1, 7)), "unexpected PUBLIC_CONSTANT, expected a method name (identifier, overridable operator)"),
			},
		},
		"cannot have a non overridable operator as a name": {
			input: "def &&; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							"&&",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(5, 1, 6)), "unexpected &&, expected a method name (identifier, overridable operator)"),
			},
		},
		"cannot have a private constant as a name": {
			input: "def _Foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							"_Foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(7, 1, 8)), "unexpected PRIVATE_CONSTANT, expected a method name (identifier, overridable operator)"),
			},
		},
		"can have an empty argument list": {
			input: "def foo(); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							"foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a return type and omit arguments": {
			input: "def foo: String?; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							"foo",
							nil,
							ast.NewNilableTypeNode(
								S(P(9, 1, 10), P(15, 1, 16)),
								ast.NewPublicConstantNode(S(P(9, 1, 10), P(14, 1, 15)), "String"),
							),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a throw type and omit arguments": {
			input: "def foo! NoMethodError | TypeError; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 1, 39)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(38, 1, 39)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(38, 1, 39)),
							"foo",
							nil,
							nil,
							ast.NewBinaryTypeExpressionNode(
								S(P(9, 1, 10), P(33, 1, 34)),
								T(S(P(23, 1, 24), P(23, 1, 24)), token.OR),
								ast.NewPublicConstantNode(S(P(9, 1, 10), P(21, 1, 22)), "NoMethodError"),
								ast.NewPublicConstantNode(S(P(25, 1, 26), P(33, 1, 34)), "TypeError"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a return and throw type and omit arguments": {
			input: "def foo : String? ! NoMethodError | TypeError; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(49, 1, 50)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(49, 1, 50)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(49, 1, 50)),
							"foo",
							nil,
							ast.NewNilableTypeNode(
								S(P(10, 1, 11), P(16, 1, 17)),
								ast.NewPublicConstantNode(S(P(10, 1, 11), P(15, 1, 16)), "String"),
							),
							ast.NewBinaryTypeExpressionNode(
								S(P(20, 1, 21), P(44, 1, 45)),
								T(S(P(34, 1, 35), P(34, 1, 35)), token.OR),
								ast.NewPublicConstantNode(S(P(20, 1, 21), P(32, 1, 33)), "NoMethodError"),
								ast.NewPublicConstantNode(S(P(36, 1, 37), P(44, 1, 45)), "TypeError"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have arguments": {
			input: "def foo(a, b); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a trailing comma in parameters": {
			input: "def foo(a, b,); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have multiline parameters": {
			input: "def foo(\na,\nb\n); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 4, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 4, 6)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(19, 4, 6)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(9, 2, 1), P(9, 2, 1)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(12, 3, 1), P(12, 3, 1)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a trailing comma in multiline parameters": {
			input: "def foo(\na,\nb,\n); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 4, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 4, 6)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(20, 4, 6)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(9, 2, 1), P(9, 2, 1)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(12, 3, 1), P(12, 3, 1)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a positional rest parameter": {
			input: "def foo(a, b, *c); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(15, 1, 16)),
									"c",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have a positional rest parameter with a default value": {
			input: "def foo(a, b, *c = 3); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(19, 1, 20)),
									"c",
									false,
									nil,
									ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "3"),
									ast.PositionalRestParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(14, 1, 15), P(19, 1, 20)), "rest parameters cannot have default values"),
			},
		},
		"can have a positional rest parameter in the middle": {
			input: "def foo(a, b, *c, d); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(15, 1, 16)),
									"c",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(18, 1, 19), P(18, 1, 19)),
									"d",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have an optional parameter after a positional rest parameter": {
			input: "def foo(a, b, *c, d = 3); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 1, 29)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(28, 1, 29)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(15, 1, 16)),
									"c",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(18, 1, 19), P(22, 1, 23)),
									"d",
									false,
									nil,
									ast.NewIntLiteralNode(S(P(22, 1, 23), P(22, 1, 23)), "3"),
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(18, 1, 19), P(22, 1, 23)), "optional parameters cannot appear after rest parameters"),
			},
		},
		"cannot have multiple positional rest parameters": {
			input: "def foo(a, b, *c, *d); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(15, 1, 16)),
									"c",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(18, 1, 19), P(19, 1, 20)),
									"d",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(18, 1, 19), P(19, 1, 20)), "there should be only a single positional rest parameter"),
			},
		},
		"can have a positional rest parameter with a type": {
			input: "def foo(a, b, *c: String); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(23, 1, 24)),
									"c",
									false,
									ast.NewPublicConstantNode(S(P(18, 1, 19), P(23, 1, 24)), "String"),
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a named rest parameter": {
			input: "def foo(a, b, **c); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(16, 1, 17)),
									"c",
									false,
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have a named rest parameter with a default value": {
			input: "def foo(a, b, **c = 3); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(20, 1, 21)),
									"c",
									false,
									nil,
									ast.NewIntLiteralNode(S(P(20, 1, 21), P(20, 1, 21)), "3"),
									ast.NamedRestParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(14, 1, 15), P(20, 1, 21)), "rest parameters cannot have default values"),
			},
		},
		"can have a named rest parameter with a type": {
			input: "def foo(a, b, **c: String); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 1, 31)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(30, 1, 31)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(30, 1, 31)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(24, 1, 25)),
									"c",
									false,
									ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have parameters after a named rest parameter": {
			input: "def foo(a, b, **c, d); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(16, 1, 17)),
									"c",
									false,
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(19, 1, 20), P(19, 1, 20)),
									"d",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(19, 1, 20), P(19, 1, 20)), "named rest parameters should appear last"),
			},
		},
		"can have a positional and named rest parameter": {
			input: "def foo(a, b, *c, **d); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(15, 1, 16)),
									"c",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(18, 1, 19), P(20, 1, 21)),
									"d",
									false,
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have a post parameter and a named rest parameter": {
			input: "def foo(a, b, *c, d, **e); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(15, 1, 16)),
									"c",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(18, 1, 19), P(18, 1, 19)),
									"d",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(21, 1, 22), P(23, 1, 24)),
									"e",
									false,
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(21, 1, 22), P(23, 1, 24)), "named rest parameters cannot appear after a post parameter"),
			},
		},
		"can have arguments with types": {
			input: "def foo(a: Int, b: String?); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 1, 32)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 1, 32)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(31, 1, 32)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(13, 1, 14)),
									"a",
									false,
									ast.NewPublicConstantNode(S(P(11, 1, 12), P(13, 1, 14)), "Int"),
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(16, 1, 17), P(25, 1, 26)),
									"b",
									false,
									ast.NewNilableTypeNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									),
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have arguments with initialisers": {
			input: "def foo(a = 32, b: String = 'foo'); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 1, 39)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(38, 1, 39)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(38, 1, 39)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(13, 1, 14)),
									"a",
									false,
									nil,
									ast.NewIntLiteralNode(S(P(12, 1, 13), P(13, 1, 14)), "32"),
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(16, 1, 17), P(32, 1, 33)),
									"b",
									false,
									ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									ast.NewRawStringLiteralNode(S(P(28, 1, 29), P(32, 1, 33)), "foo"),
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have arguments that set instance variables": {
			input: "def foo(@a = 32, @b: String = 'foo'); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(40, 1, 41)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 1, 41)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(40, 1, 41)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(14, 1, 15)),
									"a",
									true,
									nil,
									ast.NewIntLiteralNode(S(P(13, 1, 14), P(14, 1, 15)), "32"),
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(17, 1, 18), P(34, 1, 35)),
									"b",
									true,
									ast.NewPublicConstantNode(S(P(21, 1, 22), P(26, 1, 27)), "String"),
									ast.NewRawStringLiteralNode(S(P(30, 1, 31), P(34, 1, 35)), "foo"),
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have required arguments after optional ones": {
			input: "def foo(a = 32, b: String, c = true, d); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(43, 1, 44)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(43, 1, 44)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(43, 1, 44)),
							"foo",
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(13, 1, 14)),
									"a",
									false,
									nil,
									ast.NewIntLiteralNode(S(P(12, 1, 13), P(13, 1, 14)), "32"),
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(16, 1, 17), P(24, 1, 25)),
									"b",
									false,
									ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(27, 1, 28), P(34, 1, 35)),
									"c",
									false,
									nil,
									ast.NewTrueLiteralNode(S(P(31, 1, 32), P(34, 1, 35))),
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(37, 1, 38), P(37, 1, 38)),
									"d",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(16, 1, 17), P(24, 1, 25)), "required parameters cannot appear after optional parameters"),
				errors.NewError(L("main", P(37, 1, 38), P(37, 1, 38)), "required parameters cannot appear after optional parameters"),
			},
		},
		"can have a multiline body": {
			input: `def foo
  a := .5
  a += .7
end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 4, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(30, 4, 3)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(30, 4, 3)),
							"foo",
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(10, 2, 3), P(17, 2, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(10, 2, 3), P(16, 2, 9)),
										T(S(P(12, 2, 5), P(13, 2, 6)), token.COLON_EQUAL),
										ast.NewPublicIdentifierNode(S(P(10, 2, 3), P(10, 2, 3)), "a"),
										ast.NewFloatLiteralNode(S(P(15, 2, 8), P(16, 2, 9)), "0.5"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(20, 3, 3), P(27, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(20, 3, 3), P(26, 3, 9)),
										T(S(P(22, 3, 5), P(23, 3, 6)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(20, 3, 3), P(20, 3, 3)), "a"),
										ast.NewFloatLiteralNode(S(P(25, 3, 8), P(26, 3, 9)), "0.7"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can be single line with then": {
			input: `def foo then .3 + .4`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							"foo",
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(13, 1, 14), P(19, 1, 20)),
									ast.NewBinaryExpressionNode(
										S(P(13, 1, 14), P(19, 1, 20)),
										T(S(P(16, 1, 17), P(16, 1, 17)), token.PLUS),
										ast.NewFloatLiteralNode(S(P(13, 1, 14), P(14, 1, 15)), "0.3"),
										ast.NewFloatLiteralNode(S(P(18, 1, 19), P(19, 1, 20)), "0.4"),
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

func TestInitDefinition(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "bar = init; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "bar"),
							ast.NewInitDefinitionNode(
								S(P(6, 1, 7), P(14, 1, 15)),
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have an empty argument list": {
			input: "init(); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have a return type": {
			input: "init: String?; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(4, 1, 5), P(11, 1, 12)),
									ast.NewSimpleSymbolLiteralNode(
										S(P(4, 1, 5), P(11, 1, 12)),
										"String",
									),
								),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(4, 1, 5)), "unexpected :, expected a statement separator `\\n`, `;`"),
				errors.NewError(L("main", P(12, 1, 13), P(12, 1, 13)), "unexpected ?, expected a statement separator `\\n`, `;`"),
			},
		},
		"can have a throw type and omit arguments": {
			input: "init! NoMethodError | TypeError; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(35, 1, 36)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(35, 1, 36)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(35, 1, 36)),
							nil,
							ast.NewBinaryTypeExpressionNode(
								S(P(6, 1, 7), P(30, 1, 31)),
								T(S(P(20, 1, 21), P(20, 1, 21)), token.OR),
								ast.NewPublicConstantNode(S(P(6, 1, 7), P(18, 1, 19)), "NoMethodError"),
								ast.NewPublicConstantNode(S(P(22, 1, 23), P(30, 1, 31)), "TypeError"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have arguments": {
			input: "init(a, b); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(5, 1, 6)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a positional rest parameter": {
			input: "init(a, b, *c); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(5, 1, 6)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(12, 1, 13)),
									"c",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a positional rest parameter in the middle": {
			input: "init(a, b, *c, d); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(5, 1, 6)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(12, 1, 13)),
									"c",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(15, 1, 16), P(15, 1, 16)),
									"d",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have multiple positional rest parameters": {
			input: "init(a, b, *c, *d); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(5, 1, 6)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(12, 1, 13)),
									"c",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(15, 1, 16), P(16, 1, 17)),
									"d",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(15, 1, 16), P(16, 1, 17)), "there should be only a single positional rest parameter"),
			},
		},
		"can have a positional rest parameter with a type": {
			input: "init(a, b, *c: String); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(5, 1, 6)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(20, 1, 21)),
									"c",
									false,
									ast.NewPublicConstantNode(S(P(15, 1, 16), P(20, 1, 21)), "String"),
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a named rest parameter": {
			input: "init(a, b, **c); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(5, 1, 6)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(13, 1, 14)),
									"c",
									false,
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a named rest parameter with a type": {
			input: "init(a, b, **c: String); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(5, 1, 6)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(21, 1, 22)),
									"c",
									false,
									ast.NewPublicConstantNode(S(P(16, 1, 17), P(21, 1, 22)), "String"),
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have parameters after a named rest parameter": {
			input: "init(a, b, **c, d); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(5, 1, 6)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(13, 1, 14)),
									"c",
									false,
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(16, 1, 17), P(16, 1, 17)),
									"d",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(16, 1, 17), P(16, 1, 17)), "named rest parameters should appear last"),
			},
		},
		"can have a positional and named rest parameter": {
			input: "init(a, b, *c, **d); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(23, 1, 24)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(23, 1, 24)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(5, 1, 6)),
									"a",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"b",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(12, 1, 13)),
									"c",
									false,
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(15, 1, 16), P(17, 1, 18)),
									"d",
									false,
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have arguments with types": {
			input: "init(a: Int, b: String?); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 1, 29)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(28, 1, 29)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(10, 1, 11)),
									"a",
									false,
									ast.NewPublicConstantNode(S(P(8, 1, 9), P(10, 1, 11)), "Int"),
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(13, 1, 14), P(22, 1, 23)),
									"b",
									false,
									ast.NewNilableTypeNode(
										S(P(16, 1, 17), P(22, 1, 23)),
										ast.NewPublicConstantNode(S(P(16, 1, 17), P(21, 1, 22)), "String"),
									),
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have arguments with initialisers": {
			input: "init(a = 32, b: String = 'foo'); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(35, 1, 36)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(35, 1, 36)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(35, 1, 36)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(10, 1, 11)),
									"a",
									false,
									nil,
									ast.NewIntLiteralNode(S(P(9, 1, 10), P(10, 1, 11)), "32"),
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(13, 1, 14), P(29, 1, 30)),
									"b",
									false,
									ast.NewPublicConstantNode(S(P(16, 1, 17), P(21, 1, 22)), "String"),
									ast.NewRawStringLiteralNode(S(P(25, 1, 26), P(29, 1, 30)), "foo"),
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have arguments that set instance variables": {
			input: "init(@a = 32, @b: String = 'foo'); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(37, 1, 38)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(37, 1, 38)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(37, 1, 38)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(11, 1, 12)),
									"a",
									true,
									nil,
									ast.NewIntLiteralNode(S(P(10, 1, 11), P(11, 1, 12)), "32"),
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(14, 1, 15), P(31, 1, 32)),
									"b",
									true,
									ast.NewPublicConstantNode(S(P(18, 1, 19), P(23, 1, 24)), "String"),
									ast.NewRawStringLiteralNode(S(P(27, 1, 28), P(31, 1, 32)), "foo"),
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have required arguments after optional ones": {
			input: "init(a = 32, b: String, c = true, d); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(40, 1, 41)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 1, 41)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(40, 1, 41)),
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(5, 1, 6), P(10, 1, 11)),
									"a",
									false,
									nil,
									ast.NewIntLiteralNode(S(P(9, 1, 10), P(10, 1, 11)), "32"),
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(13, 1, 14), P(21, 1, 22)),
									"b",
									false,
									ast.NewPublicConstantNode(S(P(16, 1, 17), P(21, 1, 22)), "String"),
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(24, 1, 25), P(31, 1, 32)),
									"c",
									false,
									nil,
									ast.NewTrueLiteralNode(S(P(28, 1, 29), P(31, 1, 32))),
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(34, 1, 35), P(34, 1, 35)),
									"d",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(13, 1, 14), P(21, 1, 22)), "required parameters cannot appear after optional parameters"),
				errors.NewError(L("main", P(34, 1, 35), P(34, 1, 35)), "required parameters cannot appear after optional parameters"),
			},
		},
		"can have a multiline body": {
			input: `init
  a := .5
  a += .7
end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 4, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 4, 3)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(27, 4, 3)),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(7, 2, 3), P(14, 2, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(7, 2, 3), P(13, 2, 9)),
										T(S(P(9, 2, 5), P(10, 2, 6)), token.COLON_EQUAL),
										ast.NewPublicIdentifierNode(S(P(7, 2, 3), P(7, 2, 3)), "a"),
										ast.NewFloatLiteralNode(S(P(12, 2, 8), P(13, 2, 9)), "0.5"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(17, 3, 3), P(24, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(17, 3, 3), P(23, 3, 9)),
										T(S(P(19, 3, 5), P(20, 3, 6)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(17, 3, 3), P(17, 3, 3)), "a"),
										ast.NewFloatLiteralNode(S(P(22, 3, 8), P(23, 3, 9)), "0.7"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can be single line with then": {
			input: `init then .3 + .4`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewInitDefinitionNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(10, 1, 11), P(16, 1, 17)),
									ast.NewBinaryExpressionNode(
										S(P(10, 1, 11), P(16, 1, 17)),
										T(S(P(13, 1, 14), P(13, 1, 14)), token.PLUS),
										ast.NewFloatLiteralNode(S(P(10, 1, 11), P(11, 1, 12)), "0.3"),
										ast.NewFloatLiteralNode(S(P(15, 1, 16), P(16, 1, 17)), "0.4"),
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

func TestMethodSignatureDefinition(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "bar = sig foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							T(S(P(4, 1, 5), P(4, 1, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(2, 1, 3)), "bar"),
							ast.NewMethodSignatureDefinitionNode(
								S(P(6, 1, 7), P(12, 1, 13)),
								"foo",
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have a public identifier as a name": {
			input: "sig foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							"foo",
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a private identifier as a name": {
			input: "sig _foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							"_foo",
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a keyword as a name": {
			input: "sig class",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							"class",
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have an overridable operator as a name": {
			input: "sig +",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							"+",
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have a public constant as a name": {
			input: "sig Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							"Foo",
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(6, 1, 7)), "unexpected PUBLIC_CONSTANT, expected a method name (identifier, overridable operator)"),
			},
		},
		"cannot have a non overridable operator as a name": {
			input: "sig &&",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							"&&",
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(5, 1, 6)), "unexpected &&, expected a method name (identifier, overridable operator)"),
			},
		},
		"cannot have a private constant as a name": {
			input: "sig _Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							"_Foo",
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(4, 1, 5), P(7, 1, 8)), "unexpected PRIVATE_CONSTANT, expected a method name (identifier, overridable operator)"),
			},
		},
		"can have an empty argument list": {
			input: "sig foo()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							"foo",
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a return type and omit arguments": {
			input: "sig foo: String?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							"foo",
							nil,
							ast.NewNilableTypeNode(
								S(P(9, 1, 10), P(15, 1, 16)),
								ast.NewPublicConstantNode(S(P(9, 1, 10), P(14, 1, 15)), "String"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a throw type and omit arguments": {
			input: "sig foo! NoMethodError | TypeError",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							"foo",
							nil,
							nil,
							ast.NewBinaryTypeExpressionNode(
								S(P(9, 1, 10), P(33, 1, 34)),
								T(S(P(23, 1, 24), P(23, 1, 24)), token.OR),
								ast.NewPublicConstantNode(S(P(9, 1, 10), P(21, 1, 22)), "NoMethodError"),
								ast.NewPublicConstantNode(S(P(25, 1, 26), P(33, 1, 34)), "TypeError"),
							),
						),
					),
				},
			),
		},
		"can have a return and throw type and omit arguments": {
			input: "sig foo : String? ! NoMethodError | TypeError",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(44, 1, 45)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(44, 1, 45)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(44, 1, 45)),
							"foo",
							nil,
							ast.NewNilableTypeNode(
								S(P(10, 1, 11), P(16, 1, 17)),
								ast.NewPublicConstantNode(S(P(10, 1, 11), P(15, 1, 16)), "String"),
							),
							ast.NewBinaryTypeExpressionNode(
								S(P(20, 1, 21), P(44, 1, 45)),
								T(S(P(34, 1, 35), P(34, 1, 35)), token.OR),
								ast.NewPublicConstantNode(S(P(20, 1, 21), P(32, 1, 33)), "NoMethodError"),
								ast.NewPublicConstantNode(S(P(36, 1, 37), P(44, 1, 45)), "TypeError"),
							),
						),
					),
				},
			),
		},
		"can have arguments": {
			input: "sig foo(a, b)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							"foo",
							[]ast.ParameterNode{
								ast.NewSignatureParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									nil,
									false,
								),
								ast.NewSignatureParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									nil,
									false,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have arguments with types": {
			input: "sig foo(a: Int, b: String?)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							"foo",
							[]ast.ParameterNode{
								ast.NewSignatureParameterNode(
									S(P(8, 1, 9), P(13, 1, 14)),
									"a",
									ast.NewPublicConstantNode(S(P(11, 1, 12), P(13, 1, 14)), "Int"),
									false,
								),
								ast.NewSignatureParameterNode(
									S(P(16, 1, 17), P(25, 1, 26)),
									"b",
									ast.NewNilableTypeNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									),
									false,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have optional arguments": {
			input: "sig foo(a, b?, c?: String?)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							"foo",
							[]ast.ParameterNode{
								ast.NewSignatureParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									nil,
									false,
								),
								ast.NewSignatureParameterNode(
									S(P(11, 1, 12), P(12, 1, 13)),
									"b",
									nil,
									true,
								),
								ast.NewSignatureParameterNode(
									S(P(15, 1, 16), P(25, 1, 26)),
									"c",
									ast.NewNilableTypeNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									),
									true,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have required parameters after optional ones": {
			input: "sig foo(a?, b, c?, d)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							"foo",
							[]ast.ParameterNode{
								ast.NewSignatureParameterNode(
									S(P(8, 1, 9), P(9, 1, 10)),
									"a",
									nil,
									true,
								),
								ast.NewSignatureParameterNode(
									S(P(12, 1, 13), P(12, 1, 13)),
									"b",
									nil,
									false,
								),
								ast.NewSignatureParameterNode(
									S(P(15, 1, 16), P(16, 1, 17)),
									"c",
									nil,
									true,
								),
								ast.NewSignatureParameterNode(
									S(P(19, 1, 20), P(19, 1, 20)),
									"d",
									nil,
									false,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(12, 1, 13), P(12, 1, 13)), "required parameters cannot appear after optional parameters"),
				errors.NewError(L("main", P(19, 1, 20), P(19, 1, 20)), "required parameters cannot appear after optional parameters"),
			},
		},
		"cannot have arguments with initialisers": {
			input: "sig foo(a = 32, b: String = 'foo')",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(10, 1, 11), P(33, 1, 34)),
						ast.NewInvalidNode(
							S(P(10, 1, 11), P(10, 1, 11)),
							T(S(P(10, 1, 11), P(10, 1, 11)), token.EQUAL_OP),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("main", P(10, 1, 11), P(10, 1, 11)), "unexpected =, expected )"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}
