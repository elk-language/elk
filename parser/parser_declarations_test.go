package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/token"
)

func TestImport(t *testing.T) {
	tests := testTable{
		"import without a path": {
			input: "import",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewInvalidNode(
						S(P(6, 1, 7), P(5, 1, 6)),
						T(S(P(6, 1, 7), P(5, 1, 6)), token.END_OF_FILE),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(5, 1, 6)), "unexpected END_OF_FILE, expected a string literal"),
			},
		},
		"import with a non string path": {
			input: "import 3",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewInvalidNode(
						S(P(7, 1, 8), P(7, 1, 8)),
						V(S(P(7, 1, 8), P(7, 1, 8)), token.INT, "3"),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(7, 1, 8), P(7, 1, 8)), "unexpected INT, expected a string literal"),
			},
		},
		"import with a raw string": {
			input: "import 'foo'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewImportStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewRawStringLiteralNode(
							S(P(7, 1, 8), P(11, 1, 12)),
							"foo",
						),
					),
				},
			),
		},
		"import with a double quoted string": {
			input: `import "foo"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewImportStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewDoubleQuotedStringLiteralNode(
							S(P(7, 1, 8), P(11, 1, 12)),
							"foo",
						),
					),
				},
			),
		},
		"import with an interpolated string": {
			input: `import "foo${1}"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewImportStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewInterpolatedStringLiteralNode(
							S(P(7, 1, 8), P(15, 1, 16)),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(
									S(P(8, 1, 9), P(10, 1, 11)),
									"foo",
								),
								ast.NewStringInterpolationNode(
									S(P(11, 1, 12), P(14, 1, 15)),
									ast.NewIntLiteralNode(
										S(P(13, 1, 14), P(13, 1, 14)),
										"1",
									),
								),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(7, 1, 8), P(15, 1, 16)), "cannot interpolate strings in this context"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

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
		"cannot appear in expressions": {
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(8, 3, 2), P(28, 4, 11)), "singleton definitions cannot appear in expressions"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}
func TestExtendWhereBlock(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
extend where T < String
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(43, 5, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(43, 5, 4)),
						ast.NewExtendWhereBlockExpressionNode(
							S(P(1, 2, 1), P(42, 5, 3)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(26, 3, 2), P(34, 3, 10)),
									ast.NewAssignmentExpressionNode(
										S(P(26, 3, 2), P(33, 3, 9)),
										T(S(P(30, 3, 6), P(31, 3, 7)), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(S(P(26, 3, 2), P(28, 3, 4)), "foo"),
										ast.NewIntLiteralNode(S(P(33, 3, 9), P(33, 3, 9)), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(36, 4, 2), P(39, 4, 5)),
									ast.NewNilLiteralNode(S(P(36, 4, 2), P(38, 4, 4))),
								),
							},
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(14, 2, 14), P(23, 2, 23)),
									ast.INVARIANT,
									"T",
									nil,
									ast.NewPublicConstantNode(S(P(18, 2, 18), P(23, 2, 23)), "String"),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
extend where T > Float
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(27, 3, 4)),
						ast.NewExtendWhereBlockExpressionNode(
							S(P(1, 2, 1), P(26, 3, 3)),
							nil,
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(14, 2, 14), P(22, 2, 22)),
									ast.INVARIANT,
									"T",
									ast.NewPublicConstantNode(S(P(18, 2, 18), P(22, 2, 22)), "Float"),
									nil,
								),
							},
						),
					),
				},
			),
		},
		"can have multiple type parameters": {
			input: `
extend where T > Float, Y > String, Z = Int
end
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(48, 3, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(48, 3, 4)),
						ast.NewExtendWhereBlockExpressionNode(
							S(P(1, 2, 1), P(47, 3, 3)),
							nil,
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(14, 2, 14), P(22, 2, 22)),
									ast.INVARIANT,
									"T",
									ast.NewPublicConstantNode(S(P(18, 2, 18), P(22, 2, 22)), "Float"),
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(25, 2, 25), P(34, 2, 34)),
									ast.INVARIANT,
									"Y",
									ast.NewPublicConstantNode(S(P(29, 2, 29), P(34, 2, 34)), "String"),
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(37, 2, 37), P(43, 2, 43)),
									ast.INVARIANT,
									"Z",
									ast.NewPublicConstantNode(S(P(41, 2, 41), P(43, 2, 43)), "Int"),
									ast.NewPublicConstantNode(S(P(41, 2, 41), P(43, 2, 43)), "Int"),
								),
							},
						),
					),
				},
			),
		},
		"cannot appear in expressions": {
			input: `
bar =
	extend where T > Float
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(50, 6, 4)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(1, 2, 1), P(46, 5, 5)),
						ast.NewAssignmentExpressionNode(
							S(P(1, 2, 1), P(45, 5, 4)),
							T(S(P(5, 2, 5), P(5, 2, 5)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(1, 2, 1), P(3, 2, 3)), "bar"),
							ast.NewExtendWhereBlockExpressionNode(
								S(P(8, 3, 2), P(45, 5, 4)),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										S(P(33, 4, 3), P(41, 4, 11)),
										ast.NewAssignmentExpressionNode(
											S(P(33, 4, 3), P(40, 4, 10)),
											T(S(P(37, 4, 7), P(38, 4, 8)), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(S(P(33, 4, 3), P(35, 4, 5)), "foo"),
											ast.NewIntLiteralNode(S(P(40, 4, 10), P(40, 4, 10)), "2"),
										),
									),
								},
								[]ast.TypeParameterNode{
									ast.NewVariantTypeParameterNode(
										S(P(21, 3, 15), P(29, 3, 23)),
										ast.INVARIANT,
										"T",
										ast.NewPublicConstantNode(S(P(25, 3, 19), P(29, 3, 23)), "Float"),
										nil,
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(47, 6, 1), P(50, 6, 4)),
						ast.NewNilLiteralNode(S(P(47, 6, 1), P(49, 6, 3))),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(8, 3, 2), P(41, 4, 11)), "extend where definitions cannot appear in expressions"),
			},
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
						S(P(9, 1, 10), P(8, 1, 9)),
						ast.NewInvalidNode(
							S(P(9, 1, 10), P(8, 1, 9)),
							T(S(P(9, 1, 10), P(8, 1, 9)), token.END_OF_FILE),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 1, 10), P(8, 1, 9)), "unexpected END_OF_FILE, expected an expression"),
				error.NewFailure(L("<main>", P(9, 1, 10), P(8, 1, 9)), "doc comments cannot be attached to this expression"),
			},
		},
		"cannot be nested": {
			input: "##[foo]## ##[bar]## 1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(20, 1, 21), P(20, 1, 21)),
						ast.NewIntLiteralNode(S(P(20, 1, 21), P(20, 1, 21)), "1"),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(10, 1, 11), P(18, 1, 19)), "doc comments cannot document one another"),
				error.NewFailure(L("<main>", P(20, 1, 21), P(20, 1, 21)), "doc comments cannot be attached to this expression"),
			},
		},
		"can be empty": {
			input: "##[]## def foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(7, 1, 8), P(18, 1, 19)),
						ast.NewMethodDefinitionNode(
							S(P(7, 1, 8), P(18, 1, 19)),
							"",
							false,
							false,
							"foo",
							nil,
							nil,
							nil,
							nil,
							nil,
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
						S(P(39, 6, 5), P(51, 6, 17)),
						ast.NewMethodDefinitionNode(
							S(P(39, 6, 5), P(50, 6, 16)),
							"foo\nbar",
							false,
							false,
							"foo",
							nil,
							nil,
							nil,
							nil,
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

func TestUsingExpression(t *testing.T) {
	tests := testTable{
		"cannot omit the argument": {
			input: "using",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							[]ast.UsingEntryNode{
								ast.NewInvalidNode(
									S(P(5, 1, 6), P(4, 1, 5)),
									T(S(P(5, 1, 6), P(4, 1, 5)), token.END_OF_FILE),
								),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(5, 1, 6), P(4, 1, 5)), "unexpected END_OF_FILE, expected a constant"),
			},
		},
		"can have a public constant as the argument": {
			input: "using Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							[]ast.UsingEntryNode{
								ast.NewPublicConstantNode(
									S(P(6, 1, 7), P(15, 1, 16)),
									"Enumerable",
								),
							},
						),
					),
				},
			),
		},
		"can specify all members of a namespace": {
			input: "using Enumerable::*",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							[]ast.UsingEntryNode{
								ast.NewUsingAllEntryNode(
									S(P(6, 1, 7), P(18, 1, 19)),
									ast.NewPublicConstantNode(
										S(P(6, 1, 7), P(15, 1, 16)),
										"Enumerable",
									),
								),
							},
						),
					),
				},
			),
		},
		"can specify members of a namespace": {
			input: "using Enumerable::{Foo, bar}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							[]ast.UsingEntryNode{
								ast.NewUsingEntryWithSubentriesNode(
									S(P(6, 1, 7), P(27, 1, 28)),
									ast.NewPublicConstantNode(
										S(P(6, 1, 7), P(15, 1, 16)),
										"Enumerable",
									),
									[]ast.UsingSubentryNode{
										ast.NewPublicConstantNode(
											S(P(19, 1, 20), P(21, 1, 22)),
											"Foo",
										),
										ast.NewPublicIdentifierNode(
											S(P(24, 1, 25), P(26, 1, 27)),
											"bar",
										),
									},
								),
							},
						),
					),
				},
			),
		},
		"can specify members of a namespace with changed names": {
			input: "using Enumerable::{Foo as F, bar as b}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(37, 1, 38)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(37, 1, 38)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(37, 1, 38)),
							[]ast.UsingEntryNode{
								ast.NewUsingEntryWithSubentriesNode(
									S(P(6, 1, 7), P(37, 1, 38)),
									ast.NewPublicConstantNode(
										S(P(6, 1, 7), P(15, 1, 16)),
										"Enumerable",
									),
									[]ast.UsingSubentryNode{
										ast.NewPublicConstantAsNode(
											S(P(19, 1, 20), P(26, 1, 27)),
											ast.NewPublicConstantNode(
												S(P(19, 1, 20), P(21, 1, 22)),
												"Foo",
											),
											"F",
										),
										ast.NewPublicIdentifierAsNode(
											S(P(29, 1, 30), P(36, 1, 37)),
											ast.NewPublicIdentifierNode(
												S(P(29, 1, 30), P(31, 1, 32)),
												"bar",
											),
											"b",
										),
									},
								),
							},
						),
					),
				},
			),
		},
		"cannot appear in expressions": {
			input: "var a = using Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(23, 1, 24)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(23, 1, 24)),
							"",
							"a",
							nil,
							ast.NewUsingExpressionNode(
								S(P(8, 1, 9), P(23, 1, 24)),
								[]ast.UsingEntryNode{
									ast.NewPublicConstantNode(
										S(P(14, 1, 15), P(23, 1, 24)),
										"Enumerable",
									),
								},
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(8, 1, 9), P(12, 1, 13)), "using declarations cannot appear in expressions"),
			},
		},
		"can have many arguments": {
			input: "using Enumerable, Memoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							[]ast.UsingEntryNode{
								ast.NewPublicConstantNode(
									S(P(6, 1, 7), P(15, 1, 16)),
									"Enumerable",
								),
								ast.NewPublicConstantNode(
									S(P(18, 1, 19), P(27, 1, 28)),
									"Memoizable",
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after the comma": {
			input: "using Enumerable,\nMemoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 2, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 2, 10)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(27, 2, 10)),
							[]ast.UsingEntryNode{
								ast.NewPublicConstantNode(
									S(P(6, 1, 7), P(15, 1, 16)),
									"Enumerable",
								),
								ast.NewPublicConstantNode(
									S(P(18, 2, 1), P(27, 2, 10)),
									"Memoizable",
								),
							},
						),
					),
				},
			),
		},
		"can have a private constant as the argument": {
			input: "using _Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							[]ast.UsingEntryNode{
								ast.NewPrivateConstantNode(
									S(P(6, 1, 7), P(16, 1, 17)),
									"_Enumerable",
								),
							},
						),
					),
				},
			),
		},
		"can have a constant lookup as the argument": {
			input: "using Std::Memoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							[]ast.UsingEntryNode{
								ast.NewConstantLookupNode(
									S(P(6, 1, 7), P(20, 1, 21)),
									ast.NewPublicConstantNode(
										S(P(6, 1, 7), P(8, 1, 9)),
										"Std",
									),
									ast.NewPublicConstantNode(
										S(P(11, 1, 12), P(20, 1, 21)),
										"Memoizable",
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a constant lookup with as": {
			input: "using Std::Memoizable as M",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							[]ast.UsingEntryNode{
								ast.NewConstantAsNode(
									S(P(6, 1, 7), P(20, 1, 21)),
									ast.NewConstantLookupNode(
										S(P(6, 1, 7), P(20, 1, 21)),
										ast.NewPublicConstantNode(
											S(P(6, 1, 7), P(8, 1, 9)),
											"Std",
										),
										ast.NewPublicConstantNode(
											S(P(11, 1, 12), P(20, 1, 21)),
											"Memoizable",
										),
									),
									"M",
								),
							},
						),
					),
				},
			),
		},
		"can have a constant lookup with as and identifier": {
			input: "using Std::Memoizable as a",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							[]ast.UsingEntryNode{
								ast.NewInvalidNode(
									S(P(25, 1, 26), P(25, 1, 26)),
									V(S(P(25, 1, 26), P(25, 1, 26)), token.PUBLIC_IDENTIFIER, "a"),
								),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(25, 1, 26), P(25, 1, 26)), "unexpected PUBLIC_IDENTIFIER, expected PUBLIC_CONSTANT"),
			},
		},
		"can have a method lookup as the argument": {
			input: "using Std::Memoizable::memo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							[]ast.UsingEntryNode{
								ast.NewMethodLookupNode(
									S(P(6, 1, 7), P(26, 1, 27)),
									ast.NewConstantLookupNode(
										S(P(6, 1, 7), P(20, 1, 21)),
										ast.NewPublicConstantNode(
											S(P(6, 1, 7), P(8, 1, 9)),
											"Std",
										),
										ast.NewPublicConstantNode(
											S(P(11, 1, 12), P(20, 1, 21)),
											"Memoizable",
										),
									),
									"memo",
								),
							},
						),
					),
				},
			),
		},
		"can have a method lookup with as and public identifier": {
			input: "using Std::Memoizable::memo as m",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 1, 32)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 1, 32)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							[]ast.UsingEntryNode{
								ast.NewMethodLookupAsNode(
									S(P(6, 1, 7), P(20, 1, 21)),
									ast.NewMethodLookupNode(
										S(P(6, 1, 7), P(26, 1, 27)),
										ast.NewConstantLookupNode(
											S(P(6, 1, 7), P(20, 1, 21)),
											ast.NewPublicConstantNode(
												S(P(6, 1, 7), P(8, 1, 9)),
												"Std",
											),
											ast.NewPublicConstantNode(
												S(P(11, 1, 12), P(20, 1, 21)),
												"Memoizable",
											),
										),
										"memo",
									),
									"m",
								),
							},
						),
					),
				},
			),
		},
		"cannot have a method lookup with as and constant": {
			input: "using Std::Memoizable::memo as M",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 1, 32)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 1, 32)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(31, 1, 32)),
							[]ast.UsingEntryNode{
								ast.NewInvalidNode(
									S(P(31, 1, 32), P(31, 1, 32)),
									V(S(P(31, 1, 32), P(31, 1, 32)), token.PUBLIC_CONSTANT, "M"),
								),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(31, 1, 32), P(31, 1, 32)), "unexpected PUBLIC_CONSTANT, expected PUBLIC_IDENTIFIER"),
			},
		},
		"can have a generic constant as the argument": {
			input: "using Enumerable[String]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewUsingExpressionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							[]ast.UsingEntryNode{
								ast.NewPublicConstantNode(S(P(6, 1, 7), P(15, 1, 16)), "Enumerable"),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 1, 17), P(16, 1, 17)), "unexpected [, expected a statement separator `\\n`, `;`"),
			},
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(7, 1, 8), P(6, 1, 7)), "unexpected END_OF_FILE, expected a constant"),
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
		"cannot appear in expressions": {
			input: "var a = include Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							"",
							"a",
							nil,
							ast.NewIncludeExpressionNode(
								S(P(8, 1, 9), P(25, 1, 26)),
								[]ast.ComplexConstantNode{
									ast.NewPublicConstantNode(
										S(P(16, 1, 17), P(25, 1, 26)),
										"Enumerable",
									),
								},
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(8, 1, 9), P(25, 1, 26)), "this definition cannot appear in expressions"),
			},
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
									[]ast.TypeNode{
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

func TestImplementExpression(t *testing.T) {
	tests := testTable{
		"cannot omit the argument": {
			input: "implement",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewImplementExpressionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							[]ast.ComplexConstantNode{
								ast.NewInvalidNode(
									S(P(9, 1, 10), P(8, 1, 9)),
									T(S(P(9, 1, 10), P(8, 1, 9)), token.END_OF_FILE),
								),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 1, 10), P(8, 1, 9)), "unexpected END_OF_FILE, expected a constant"),
			},
		},
		"can have a public constant as the argument": {
			input: "implement Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewImplementExpressionNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(10, 1, 11), P(19, 1, 20)),
									"Enumerable",
								),
							},
						),
					),
				},
			),
		},
		"cannot appear in expressions": {
			input: "var a = implement Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							"",
							"a",
							nil,
							ast.NewImplementExpressionNode(
								S(P(8, 1, 9), P(27, 1, 28)),
								[]ast.ComplexConstantNode{
									ast.NewPublicConstantNode(
										S(P(18, 1, 19), P(27, 1, 28)),
										"Enumerable",
									),
								},
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(8, 1, 9), P(27, 1, 28)), "this definition cannot appear in expressions"),
			},
		},
		"can have many arguments": {
			input: "implement Enumerable, Memoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 1, 32)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 1, 32)),
						ast.NewImplementExpressionNode(
							S(P(0, 1, 1), P(31, 1, 32)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(10, 1, 11), P(19, 1, 20)),
									"Enumerable",
								),
								ast.NewPublicConstantNode(
									S(P(22, 1, 23), P(31, 1, 32)),
									"Memoizable",
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after the comma": {
			input: "implement Enumerable,\nMemoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 2, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 2, 10)),
						ast.NewImplementExpressionNode(
							S(P(0, 1, 1), P(31, 2, 10)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(
									S(P(10, 1, 11), P(19, 1, 20)),
									"Enumerable",
								),
								ast.NewPublicConstantNode(
									S(P(22, 2, 1), P(31, 2, 10)),
									"Memoizable",
								),
							},
						),
					),
				},
			),
		},
		"can have a private constant as the argument": {
			input: "implement _Enumerable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewImplementExpressionNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							[]ast.ComplexConstantNode{
								ast.NewPrivateConstantNode(
									S(P(10, 1, 11), P(20, 1, 21)),
									"_Enumerable",
								),
							},
						),
					),
				},
			),
		},
		"can have a constant lookup as the argument": {
			input: "implement Std::Memoizable",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewImplementExpressionNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							[]ast.ComplexConstantNode{
								ast.NewConstantLookupNode(
									S(P(10, 1, 11), P(24, 1, 25)),
									ast.NewPublicConstantNode(
										S(P(10, 1, 11), P(12, 1, 13)),
										"Std",
									),
									ast.NewPublicConstantNode(
										S(P(15, 1, 16), P(24, 1, 25)),
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
			input: "implement Enumerable[String]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewImplementExpressionNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							[]ast.ComplexConstantNode{
								ast.NewGenericConstantNode(
									S(P(10, 1, 11), P(27, 1, 28)),
									ast.NewPublicConstantNode(S(P(10, 1, 11), P(19, 1, 20)), "Enumerable"),
									[]ast.TypeNode{
										ast.NewPublicConstantNode(S(P(21, 1, 22), P(26, 1, 27)), "String"),
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
				implement Foo
				implement Bar
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(36, 3, 18)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(18, 2, 18)),
						ast.NewImplementExpressionNode(
							S(P(5, 2, 5), P(17, 2, 17)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(S(P(15, 2, 15), P(17, 2, 17)), "Foo"),
							},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(23, 3, 5), P(36, 3, 18)),
						ast.NewImplementExpressionNode(
							S(P(23, 3, 5), P(35, 3, 17)),
							[]ast.ComplexConstantNode{
								ast.NewPublicConstantNode(S(P(33, 3, 15), P(35, 3, 17)), "Bar"),
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
							"foo",
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
								"foo",
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
							"_foo",
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
							"foo",
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(7, 1, 8)), "instance variables cannot be declared using `val`"),
			},
		},
		"cannot have a constant as the value name": {
			input: "val Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							"Foo",
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(6, 1, 7)), "variable names cannot resemble constants"),
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
							"foo",
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
							"foo",
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
							"foo",
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
							"foo",
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
							"foo",
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
							"foo",
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
							"foo",
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
							"foo",
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
							"foo",
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
							"foo",
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
							"foo",
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
							"foo",
							ast.NewGenericConstantNode(
								S(P(9, 1, 10), P(43, 1, 44)),
								ast.NewConstantLookupNode(
									S(P(9, 1, 10), P(16, 1, 17)),
									ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Std"),
									ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Map"),
								),
								[]ast.TypeNode{
									ast.NewConstantLookupNode(
										S(P(18, 1, 19), P(28, 1, 29)),
										ast.NewPublicConstantNode(S(P(18, 1, 19), P(20, 1, 21)), "Std"),
										ast.NewPublicConstantNode(S(P(23, 1, 24), P(28, 1, 29)), "Symbol"),
									),
									ast.NewGenericConstantNode(
										S(P(31, 1, 32), P(42, 1, 43)),
										ast.NewPublicConstantNode(S(P(31, 1, 32), P(34, 1, 35)), "List"),
										[]ast.TypeNode{
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
		"can have a pattern": {
			input: "val [a, { b, c: 2 }] = bar()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewValuePatternDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewListPatternNode(
								S(P(4, 1, 5), P(19, 1, 20)),
								[]ast.PatternNode{
									ast.NewPublicIdentifierNode(S(P(5, 1, 6), P(5, 1, 6)), "a"),
									ast.NewMapPatternNode(
										S(P(8, 1, 9), P(18, 1, 19)),
										[]ast.PatternNode{
											ast.NewPublicIdentifierNode(
												S(P(10, 1, 11), P(10, 1, 11)),
												"b",
											),
											ast.NewSymbolKeyValuePatternNode(
												S(P(13, 1, 14), P(16, 1, 17)),
												"c",
												ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "2"),
											),
										},
									),
								},
							),
							ast.NewReceiverlessMethodCallNode(
								S(P(23, 1, 24), P(27, 1, 28)),
								"bar",
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"cannot have a pattern without variables": {
			input: "val [1, 2] = bar()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewValuePatternDeclarationNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							ast.NewListPatternNode(
								S(P(4, 1, 5), P(9, 1, 10)),
								[]ast.PatternNode{
									ast.NewIntLiteralNode(S(P(5, 1, 6), P(5, 1, 6)), "1"),
									ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "2"),
								},
							),
							ast.NewReceiverlessMethodCallNode(
								S(P(13, 1, 14), P(17, 1, 18)),
								"bar",
								nil,
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(9, 1, 10)), "patterns in value declarations should define at least one value"),
			},
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
							"",
							"foo",
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
								"",
								"foo",
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
							"",
							"_foo",
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have an instance variable as the variable name": {
			input: "var @foo: Float",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewInstanceVariableDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							"",
							"foo",
							ast.NewPublicConstantNode(
								S(P(10, 1, 11), P(14, 1, 15)),
								"Float",
							),
						),
					),
				},
			),
		},
		"instance variable declarations cannot appear in expressions": {
			input: "f = var @foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(
								S(P(0, 1, 1), P(0, 1, 1)),
								"f",
							),
							ast.NewInstanceVariableDeclarationNode(
								S(P(4, 1, 5), P(11, 1, 12)),
								"",
								"foo",
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(11, 1, 12)), "instance variable declarations cannot appear in expressions"),
				error.NewFailure(L("<main>", P(4, 1, 5), P(11, 1, 12)), "instance variable declarations must have an explicit type"),
			},
		},
		"instance variables cannot be initialised": {
			input: "var @foo = 2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewInstanceVariableDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							"",
							"foo",
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 1, 12), P(11, 1, 12)), "instance variables cannot be initialised when declared"),
				error.NewFailure(L("<main>", P(0, 1, 1), P(11, 1, 12)), "instance variable declarations must have an explicit type"),
			},
		},
		"cannot have a constant as the variable name": {
			input: "var Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							"",
							"Foo",
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(6, 1, 7)), "variable names cannot resemble constants"),
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
							"",
							"foo",
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
							"",
							"foo",
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
							"",
							"foo",
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
							"",
							"foo",
							ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Int"),
							nil,
						),
					),
				},
			),
		},
		"can have never": {
			input: "var foo: never",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							"",
							"foo",
							ast.NewNeverTypeNode(S(P(9, 1, 10), P(13, 1, 14))),
							nil,
						),
					),
				},
			),
		},
		"cannot have void": {
			input: "var foo: void",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							"",
							"foo",
							ast.NewVoidTypeNode(S(P(9, 1, 10), P(12, 1, 13))),
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 1, 10), P(12, 1, 13)), "type `void` cannot be used in this context"),
			},
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
							"",
							"foo",
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
							"",
							"foo",
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
							"",
							"foo",
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
							"",
							"foo",
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
							"",
							"foo",
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
							"",
							"foo",
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
							"",
							"foo",
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
							"",
							"foo",
							ast.NewGenericConstantNode(
								S(P(9, 1, 10), P(43, 1, 44)),
								ast.NewConstantLookupNode(
									S(P(9, 1, 10), P(16, 1, 17)),
									ast.NewPublicConstantNode(S(P(9, 1, 10), P(11, 1, 12)), "Std"),
									ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Map"),
								),
								[]ast.TypeNode{
									ast.NewConstantLookupNode(
										S(P(18, 1, 19), P(28, 1, 29)),
										ast.NewPublicConstantNode(S(P(18, 1, 19), P(20, 1, 21)), "Std"),
										ast.NewPublicConstantNode(S(P(23, 1, 24), P(28, 1, 29)), "Symbol"),
									),
									ast.NewGenericConstantNode(
										S(P(31, 1, 32), P(42, 1, 43)),
										ast.NewPublicConstantNode(S(P(31, 1, 32), P(34, 1, 35)), "List"),
										[]ast.TypeNode{
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
							"",
							"foo",
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
							"",
							"foo",
							ast.NewSingletonTypeNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								ast.NewNilableTypeNode(
									S(P(10, 1, 11), P(13, 1, 14)),
									ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Int"),
								),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a pattern": {
			input: "var [a, { b, c: 2 }] = bar()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewVariablePatternDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewListPatternNode(
								S(P(4, 1, 5), P(19, 1, 20)),
								[]ast.PatternNode{
									ast.NewPublicIdentifierNode(S(P(5, 1, 6), P(5, 1, 6)), "a"),
									ast.NewMapPatternNode(
										S(P(8, 1, 9), P(18, 1, 19)),
										[]ast.PatternNode{
											ast.NewPublicIdentifierNode(
												S(P(10, 1, 11), P(10, 1, 11)),
												"b",
											),
											ast.NewSymbolKeyValuePatternNode(
												S(P(13, 1, 14), P(16, 1, 17)),
												"c",
												ast.NewIntLiteralNode(S(P(16, 1, 17), P(16, 1, 17)), "2"),
											),
										},
									),
								},
							),
							ast.NewReceiverlessMethodCallNode(
								S(P(23, 1, 24), P(27, 1, 28)),
								"bar",
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"cannot have a pattern without variables": {
			input: "var [1, 2] = bar()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewVariablePatternDeclarationNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							ast.NewListPatternNode(
								S(P(4, 1, 5), P(9, 1, 10)),
								[]ast.PatternNode{
									ast.NewIntLiteralNode(S(P(5, 1, 6), P(5, 1, 6)), "1"),
									ast.NewIntLiteralNode(S(P(8, 1, 9), P(8, 1, 9)), "2"),
								},
							),
							ast.NewReceiverlessMethodCallNode(
								S(P(13, 1, 14), P(17, 1, 18)),
								"bar",
								nil,
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(9, 1, 10)), "patterns in variable declarations should define at least one variable"),
			},
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
		"is invalid without an initialiser and with a type": {
			input: "const Foo: String",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							"",
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
							),
							ast.NewPublicConstantNode(S(P(11, 1, 12), P(16, 1, 17)), "String"),
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(16, 1, 17)), "constants must be initialised"),
			},
		},
		"is not valid without an initialiser and without a type": {
			input: "const Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							"",
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
							),
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(8, 1, 9)), "constants must be initialised"),
			},
		},
		"cannot be a part of an expression": {
			input: "a = const _Foo: String = 'bar'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
							ast.NewConstantDeclarationNode(
								S(P(4, 1, 5), P(29, 1, 30)),
								"",
								ast.NewPrivateConstantNode(S(P(10, 1, 11), P(13, 1, 14)), "_Foo"),
								ast.NewPublicConstantNode(S(P(16, 1, 17), P(21, 1, 22)), "String"),
								ast.NewRawStringLiteralNode(
									S(P(25, 1, 26), P(29, 1, 30)),
									"bar",
								),
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(29, 1, 30)), "constant declarations cannot appear in expressions"),
			},
		},
		"can have a private constant as the name": {
			input: "const _Foo: String = 'bar'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							"",
							ast.NewPrivateConstantNode(
								S(P(6, 1, 7), P(9, 1, 10)),
								"_Foo",
							),
							ast.NewPublicConstantNode(S(P(12, 1, 13), P(17, 1, 18)), "String"),
							ast.NewRawStringLiteralNode(
								S(P(21, 1, 22), P(25, 1, 26)),
								"bar",
							),
						),
					),
				},
			),
		},
		"cannot have an instance variable as the name": {
			input: "const @foo: String = 'bar'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							"",
							ast.NewInstanceVariableNode(
								S(P(6, 1, 7), P(9, 1, 10)),
								"foo",
							),
							ast.NewPublicConstantNode(S(P(12, 1, 13), P(17, 1, 18)), "String"),
							ast.NewRawStringLiteralNode(
								S(P(21, 1, 22), P(25, 1, 26)),
								"bar",
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(9, 1, 10)), "invalid constant name"),
			},
		},
		"cannot have a lowercase identifier as the name": {
			input: "const foo: String = 'bar'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							"",
							ast.NewPublicIdentifierNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"foo",
							),
							ast.NewPublicConstantNode(S(P(11, 1, 12), P(16, 1, 17)), "String"),
							ast.NewRawStringLiteralNode(
								S(P(20, 1, 21), P(24, 1, 25)),
								"bar",
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(8, 1, 9)), "invalid constant name"),
			},
		},
		"can have a static initialiser without a type": {
			input: "const Foo = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							"",
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
							),
							nil,
							ast.NewIntLiteralNode(S(P(12, 1, 13), P(12, 1, 13)), "5"),
						),
					),
				},
			),
		},
		"can have an initialiser without a type": {
			input: "const Foo = f",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							"",
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
							),
							nil,
							ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(12, 1, 13)), "f"),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "const Foo: String =\n5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 2, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 2, 1)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(20, 2, 1)),
							"",
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
							),
							ast.NewPublicConstantNode(S(P(11, 1, 12), P(16, 1, 17)), "String"),
							ast.NewIntLiteralNode(S(P(20, 2, 1), P(20, 2, 1)), "5"),
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
							"",
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
							),
							ast.NewPublicConstantNode(S(P(11, 1, 12), P(13, 1, 14)), "Int"),
							ast.NewIntLiteralNode(S(P(17, 1, 18), P(17, 1, 18)), "5"),
						),
					),
				},
			),
		},
		"can have a complex constant": {
			input: "const Foo::Bar: Int = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							"",
							ast.NewConstantLookupNode(
								S(P(6, 1, 7), P(13, 1, 14)),
								ast.NewPublicConstantNode(
									S(P(6, 1, 7), P(8, 1, 9)),
									"Foo",
								),
								ast.NewPublicConstantNode(
									S(P(11, 1, 12), P(13, 1, 14)),
									"Bar",
								),
							),
							ast.NewPublicConstantNode(S(P(16, 1, 17), P(18, 1, 19)), "Int"),
							ast.NewIntLiteralNode(S(P(22, 1, 23), P(22, 1, 23)), "5"),
						),
					),
				},
			),
		},
		"can have never": {
			input: "const Foo: never = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							"",
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
							),
							ast.NewNeverTypeNode(S(P(11, 1, 12), P(15, 1, 16))),
							ast.NewIntLiteralNode(S(P(19, 1, 20), P(19, 1, 20)), "5"),
						),
					),
				},
			),
		},
		"cannot have void": {
			input: "const Foo: void = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							"",
							ast.NewPublicConstantNode(
								S(P(6, 1, 7), P(8, 1, 9)),
								"Foo",
							),
							ast.NewVoidTypeNode(S(P(11, 1, 12), P(14, 1, 15))),
							ast.NewIntLiteralNode(S(P(18, 1, 19), P(18, 1, 19)), "5"),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 1, 12), P(14, 1, 15)), "type `void` cannot be used in this context"),
			},
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 1, 12), P(10, 1, 11)), "unexpected END_OF_FILE, expected ="),
			},
		},
		"can be generic": {
			input: "typedef Foo[+V > Bar, -T < Baz] = V | T",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 1, 39)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(38, 1, 39)),
						ast.NewGenericTypeDefinitionNode(
							S(P(0, 1, 1), P(38, 1, 39)),
							"",
							ast.NewPublicConstantNode(S(P(8, 1, 9), P(10, 1, 11)), "Foo"),
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(12, 1, 13), P(19, 1, 20)),
									ast.COVARIANT,
									"V",
									ast.NewPublicConstantNode(
										S(P(17, 1, 18), P(19, 1, 20)),
										"Bar",
									),
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(22, 1, 23), P(29, 1, 30)),
									ast.CONTRAVARIANT,
									"T",
									nil,
									ast.NewPublicConstantNode(
										S(P(27, 1, 28), P(29, 1, 30)),
										"Baz",
									),
								),
							},
							ast.NewBinaryTypeExpressionNode(
								S(P(34, 1, 35), P(38, 1, 39)),
								T(S(P(36, 1, 37), P(36, 1, 37)), token.OR),
								ast.NewPublicConstantNode(
									S(P(34, 1, 35), P(34, 1, 35)),
									"V",
								),
								ast.NewPublicConstantNode(
									S(P(38, 1, 39), P(38, 1, 39)),
									"T",
								),
							),
						),
					),
				},
			),
		},
		"cannot be a part of an expression": {
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
								"",
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(24, 1, 25)), "type definitions cannot appear in expressions"),
			},
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
							"",
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
							"",
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
							"",
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
							"",
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(8, 1, 9), P(11, 1, 12)), "unexpected INSTANCE_VARIABLE, expected a constant"),
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
							"",
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(8, 1, 9), P(10, 1, 11)), "unexpected PUBLIC_IDENTIFIER, expected a constant"),
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
		"cannot be a part of an expression": {
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
								"",
								[]ast.ParameterNode{
									ast.NewAttributeParameterNode(
										S(P(11, 1, 12), P(13, 1, 14)),
										"foo",
										nil,
										nil,
									),
								},
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(9, 1, 10)), "getter declarations cannot appear in expressions"),
			},
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
							"",
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(7, 1, 8), P(15, 1, 16)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(12, 1, 13), P(15, 1, 16)),
										ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
									),
									nil,
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
							"",
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(7, 1, 8), P(15, 1, 16)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(12, 1, 13), P(15, 1, 16)),
										ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
									),
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(18, 1, 19), P(20, 1, 21)),
									"bar",
									nil,
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
									nil,
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
							"",
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(12, 2, 12), P(20, 2, 20)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(17, 2, 17), P(20, 2, 20)),
										ast.NewPublicConstantNode(S(P(17, 2, 17), P(19, 2, 19)), "Bar"),
									),
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(31, 3, 9), P(33, 3, 11)),
									"bar",
									nil,
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
									nil,
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
		"cannot be a part of an expression": {
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
								"",
								[]ast.ParameterNode{
									ast.NewAttributeParameterNode(
										S(P(11, 1, 12), P(13, 1, 14)),
										"foo",
										nil,
										nil,
									),
								},
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(9, 1, 10)), "setter declarations cannot appear in expressions"),
			},
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
							"",
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(7, 1, 8), P(15, 1, 16)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(12, 1, 13), P(15, 1, 16)),
										ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
									),
									nil,
								),
							},
						),
					),
				},
			),
		},
		"cannot have an initialiser": {
			input: "setter foo: Bar? = 1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewSetterDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							"",
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(7, 1, 8), P(19, 1, 20)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(12, 1, 13), P(15, 1, 16)),
										ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
									),
									ast.NewIntLiteralNode(
										S(P(19, 1, 20), P(19, 1, 20)), "1",
									),
								),
							},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(19, 1, 20), P(19, 1, 20)), "setter declarations cannot have initialisers"),
			},
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
							"",
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(7, 1, 8), P(15, 1, 16)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(12, 1, 13), P(15, 1, 16)),
										ast.NewPublicConstantNode(S(P(12, 1, 13), P(14, 1, 15)), "Bar"),
									),
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(18, 1, 19), P(20, 1, 21)),
									"bar",
									nil,
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
									nil,
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
							"",
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(12, 2, 12), P(20, 2, 20)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(17, 2, 17), P(20, 2, 20)),
										ast.NewPublicConstantNode(S(P(17, 2, 17), P(19, 2, 19)), "Bar"),
									),
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(31, 3, 9), P(33, 3, 11)),
									"bar",
									nil,
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
									nil,
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
		"cannot be a part of an expression": {
			input: "a = attr     foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewAssignmentExpressionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							T(S(P(2, 1, 3), P(2, 1, 3)), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(S(P(0, 1, 1), P(0, 1, 1)), "a"),
							ast.NewAttrDeclarationNode(
								S(P(4, 1, 5), P(15, 1, 16)),
								"",
								[]ast.ParameterNode{
									ast.NewAttributeParameterNode(
										S(P(13, 1, 14), P(15, 1, 16)),
										"foo",
										nil,
										nil,
									),
								},
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(7, 1, 8)), "attr declarations cannot appear in expressions"),
			},
		},
		"can have a type": {
			input: "attr     foo: Bar?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewAttrDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							"",
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(9, 1, 10), P(17, 1, 18)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(14, 1, 15), P(17, 1, 18)),
										ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Bar"),
									),
									nil,
								),
							},
						),
					),
				},
			),
		},
		"can have an initialiser": {
			input: "attr     foo: Bar? = 1",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewAttrDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							"",
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(9, 1, 10), P(21, 1, 22)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(14, 1, 15), P(17, 1, 18)),
										ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Bar"),
									),
									ast.NewIntLiteralNode(
										S(P(21, 1, 22), P(21, 1, 22)), "1",
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a few attributes": {
			input: "attr     foo: Bar?, bar, baz: Int | Float",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(40, 1, 41)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 1, 41)),
						ast.NewAttrDeclarationNode(
							S(P(0, 1, 1), P(40, 1, 41)),
							"",
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(9, 1, 10), P(17, 1, 18)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(14, 1, 15), P(17, 1, 18)),
										ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Bar"),
									),
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(20, 1, 21), P(22, 1, 23)),
									"bar",
									nil,
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
									nil,
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: `
				attr     foo: Bar?,
							 bar,
							 baz: Int | Float
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(62, 4, 25)),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(S(P(0, 1, 1), P(0, 1, 1))),
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(62, 4, 25)),
						ast.NewAttrDeclarationNode(
							S(P(5, 2, 5), P(61, 4, 24)),
							"",
							[]ast.ParameterNode{
								ast.NewAttributeParameterNode(
									S(P(14, 2, 14), P(22, 2, 22)),
									"foo",
									ast.NewNilableTypeNode(
										S(P(19, 2, 19), P(22, 2, 22)),
										ast.NewPublicConstantNode(S(P(19, 2, 19), P(21, 2, 21)), "Bar"),
									),
									nil,
								),
								ast.NewAttributeParameterNode(
									S(P(33, 3, 9), P(35, 3, 11)),
									"bar",
									nil,
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
									nil,
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
		"cannot be a part of an expression": {
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(16, 1, 17)), "alias definitions cannot appear in expressions"),
			},
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(9, 1, 10)), "unexpected INSTANCE_VARIABLE, expected a method name (identifier, overridable operator)"),
				error.NewFailure(L("<main>", P(11, 1, 12), P(14, 1, 15)), "unexpected INSTANCE_VARIABLE, expected a method name (identifier, overridable operator)"),
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
		"cannot be anonymous": {
			input: `class; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							"",
							false,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(9, 1, 10)), "anonymous classes are not supported"),
			},
		},
		"cannot be a part of an expression": {
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
								"",
								false,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(15, 1, 16)), "class declarations cannot appear in expressions"),
				error.NewFailure(L("<main>", P(6, 1, 7), P(15, 1, 16)), "anonymous classes are not supported"),
			},
		},
		"cannot be anonymous with a superclass": {
			input: `class < Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							"",
							false,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(15, 1, 16)), "anonymous classes are not supported"),
			},
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
							"",
							false,
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(10, 1, 11), P(10, 1, 11)),
									ast.INVARIANT,
									"V",
									nil,
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(13, 1, 14), P(14, 1, 15)),
									ast.COVARIANT,
									"T",
									nil,
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(17, 1, 18), P(18, 1, 19)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
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
							"",
							false,
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(10, 1, 11), P(24, 1, 25)),
									ast.INVARIANT,
									"V",
									nil,
									ast.NewConstantLookupNode(
										S(P(14, 1, 15), P(24, 1, 25)),
										ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Std"),
										ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									),
								),
								ast.NewVariantTypeParameterNode(
									S(P(27, 1, 28), P(34, 1, 35)),
									ast.COVARIANT,
									"T",
									nil,
									ast.NewPublicConstantNode(S(P(32, 1, 33), P(34, 1, 35)), "Foo"),
								),
								ast.NewVariantTypeParameterNode(
									S(P(37, 1, 38), P(45, 1, 46)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
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
							"",
							false,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(10, 1, 11), P(10, 1, 11)), "unexpected ], expected a list of type variables"),
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
							"",
							true,
							false,
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
		"cannot appear in expressions with modifiers": {
			input: `var foo = abstract class Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							"",
							"foo",
							nil,
							ast.NewClassDeclarationNode(
								S(P(10, 1, 11), P(32, 1, 33)),
								"",
								true,
								false,
								false,
								ast.NewPublicConstantNode(S(P(25, 1, 26), P(27, 1, 28)), "Foo"),
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(19, 1, 20), P(32, 1, 33)), "class declarations cannot appear in expressions"),
			},
		},
		"cannot appear in expressions with doc comments": {
			input: `var foo = ##[ab]## class Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							"",
							"foo",
							nil,
							ast.NewClassDeclarationNode(
								S(P(19, 1, 20), P(32, 1, 33)),
								"ab",
								false,
								false,
								false,
								ast.NewPublicConstantNode(S(P(25, 1, 26), P(27, 1, 28)), "Foo"),
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(19, 1, 20), P(32, 1, 33)), "class declarations cannot appear in expressions"),
			},
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
							"",
							true,
							false,
							false,
							ast.NewPublicConstantNode(S(P(24, 1, 25), P(26, 1, 27)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(7, 1, 8)), "the abstract modifier can only be attached once"),
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
							"",
							true,
							true,
							false,
							ast.NewPublicConstantNode(S(P(22, 1, 23), P(24, 1, 25)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(7, 1, 8)), "the abstract modifier cannot be attached to sealed classes"),
			},
		},
		"can be primitive": {
			input: `primitive class Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(23, 1, 24)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(23, 1, 24)),
							"",
							false,
							false,
							true,
							ast.NewPublicConstantNode(S(P(16, 1, 17), P(18, 1, 19)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
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
							"",
							false,
							true,
							false,
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
							"",
							true,
							true,
							false,
							ast.NewPublicConstantNode(S(P(22, 1, 23), P(24, 1, 25)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(5, 1, 6)), "the sealed modifier cannot be attached to abstract classes"),
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
							"",
							false,
							true,
							false,
							ast.NewPublicConstantNode(S(P(20, 1, 21), P(22, 1, 23)), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(5, 1, 6)), "the sealed modifier can only be attached once"),
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
							"",
							false,
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
							"",
							false,
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
							"",
							false,
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
							"",
							false,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(8, 1, 9)), "invalid class name, expected a constant"),
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
							"",
							false,
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
		"can have nil as a superclass": {
			input: `class Foo < nil; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewClassDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							"",
							false,
							false,
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							ast.NewNilLiteralNode(S(P(12, 1, 13), P(14, 1, 15))),
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
							"",
							false,
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
							"",
							false,
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
							"",
							false,
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
								[]ast.TypeNode{
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
							"",
							false,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 1, 13), P(14, 1, 15)), "unexpected PUBLIC_IDENTIFIER, expected a constant"),
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
							"",
							false,
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
							"",
							false,
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
		"cannot be anonymous": {
			input: `module; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							"",
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(10, 1, 11)), "anonymous modules are not supported"),
			},
		},
		"cannot be a part of an expression": {
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
								"",
								nil,
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(16, 1, 17)), "module declarations cannot appear in expressions"),
				error.NewFailure(L("<main>", P(6, 1, 7), P(16, 1, 17)), "anonymous modules are not supported"),
			},
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
							"",
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(10, 1, 11), P(20, 1, 21)), "modules cannot be generic"),
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
							"",
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
							"",
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
							"",
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
							"",
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(7, 1, 8), P(9, 1, 10)), "invalid module name, expected a constant"),
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
							"",
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
							"",
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
		"cannot be anonymous": {
			input: `mixin; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							"",
							false,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(9, 1, 10)), "anonymous mixins are not supported"),
			},
		},
		"cannot be a part of an expression": {
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
								"",
								false,
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(15, 1, 16)), "mixin declarations cannot appear in expressions"),
				error.NewFailure(L("<main>", P(6, 1, 7), P(15, 1, 16)), "anonymous mixins are not supported"),
			},
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
							"",
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(10, 1, 11), P(10, 1, 11)),
									ast.INVARIANT,
									"V",
									nil,
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(13, 1, 14), P(14, 1, 15)),
									ast.COVARIANT,
									"T",
									nil,
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(17, 1, 18), P(18, 1, 19)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
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
							"",
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(10, 1, 11), P(24, 1, 25)),
									ast.INVARIANT,
									"V",
									nil,
									ast.NewConstantLookupNode(
										S(P(14, 1, 15), P(24, 1, 25)),
										ast.NewPublicConstantNode(S(P(14, 1, 15), P(16, 1, 17)), "Std"),
										ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									),
								),
								ast.NewVariantTypeParameterNode(
									S(P(27, 1, 28), P(34, 1, 35)),
									ast.COVARIANT,
									"T",
									nil,
									ast.NewPublicConstantNode(S(P(32, 1, 33), P(34, 1, 35)), "Foo"),
								),
								ast.NewVariantTypeParameterNode(
									S(P(37, 1, 38), P(45, 1, 46)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
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
							"",
							false,
							ast.NewPublicConstantNode(S(P(6, 1, 7), P(8, 1, 9)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(10, 1, 11), P(10, 1, 11)), "unexpected ], expected a list of type variables"),
			},
		},
		"cannot be sealed": {
			input: `sealed mixin Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(7, 1, 8), P(20, 1, 21)),
						ast.NewMixinDeclarationNode(
							S(P(7, 1, 8), P(20, 1, 21)),
							"",
							false,
							ast.NewPublicConstantNode(S(P(13, 1, 14), P(15, 1, 16)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(7, 1, 8), P(20, 1, 21)), "the sealed modifier can only be attached to classes and methods"),
			},
		},
		"can be abstract": {
			input: `abstract mixin Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewMixinDeclarationNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							"",
							true,
							ast.NewPublicConstantNode(S(P(15, 1, 16), P(17, 1, 18)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
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
							"",
							false,
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
							"",
							false,
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
							"",
							false,
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
							"",
							false,
							ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(8, 1, 9)), "invalid mixin name, expected a constant"),
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
							"",
							false,
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
							"",
							false,
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
		"cannot be anonymous": {
			input: `interface; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewInterfaceDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							"",
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(13, 1, 14)), "anonymous interfaces are not supported"),
			},
		},
		"cannot be a part of an expression": {
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
								"",
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(19, 1, 20)), "interface declarations cannot appear in expressions"),
				error.NewFailure(L("<main>", P(6, 1, 7), P(19, 1, 20)), "anonymous interfaces are not supported"),
			},
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
							"",
							ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Foo"),
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(14, 1, 15), P(14, 1, 15)),
									ast.INVARIANT,
									"V",
									nil,
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(17, 1, 18), P(18, 1, 19)),
									ast.COVARIANT,
									"T",
									nil,
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(21, 1, 22), P(22, 1, 23)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
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
							"",
							ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Foo"),
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(14, 1, 15), P(28, 1, 29)),
									ast.INVARIANT,
									"V",
									nil,
									ast.NewConstantLookupNode(
										S(P(18, 1, 19), P(28, 1, 29)),
										ast.NewPublicConstantNode(S(P(18, 1, 19), P(20, 1, 21)), "Std"),
										ast.NewPublicConstantNode(S(P(23, 1, 24), P(28, 1, 29)), "String"),
									),
								),
								ast.NewVariantTypeParameterNode(
									S(P(31, 1, 32), P(38, 1, 39)),
									ast.COVARIANT,
									"T",
									nil,
									ast.NewPublicConstantNode(S(P(36, 1, 37), P(38, 1, 39)), "Foo"),
								),
								ast.NewVariantTypeParameterNode(
									S(P(41, 1, 42), P(49, 1, 50)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
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
							"",
							ast.NewPublicConstantNode(S(P(10, 1, 11), P(12, 1, 13)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(14, 1, 15), P(14, 1, 15)), "unexpected ], expected a list of type variables"),
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
							"",
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
							"",
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
							"",
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
							"",
							ast.NewPublicIdentifierNode(S(P(10, 1, 11), P(12, 1, 13)), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(10, 1, 11), P(12, 1, 13)), "invalid interface name, expected a constant"),
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
							"",
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
							"",
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
		"cannot be anonymous": {
			input: `struct; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							"",
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(10, 1, 11)), "anonymous structs are not supported"),
			},
		},
		"cannot be a part of an expression": {
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
								"",
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(16, 1, 17)), "struct declarations cannot appear in expressions"),
				error.NewFailure(L("<main>", P(6, 1, 7), P(16, 1, 17)), "anonymous structs are not supported"),
			},
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
							"",
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									ast.INVARIANT,
									"V",
									nil,
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(14, 1, 15), P(15, 1, 16)),
									ast.COVARIANT,
									"T",
									nil,
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(18, 1, 19), P(19, 1, 20)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
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
							"",
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(11, 1, 12), P(25, 1, 26)),
									ast.INVARIANT,
									"V",
									nil,
									ast.NewConstantLookupNode(
										S(P(15, 1, 16), P(25, 1, 26)),
										ast.NewPublicConstantNode(S(P(15, 1, 16), P(17, 1, 18)), "Std"),
										ast.NewPublicConstantNode(S(P(20, 1, 21), P(25, 1, 26)), "String"),
									),
								),
								ast.NewVariantTypeParameterNode(
									S(P(28, 1, 29), P(35, 1, 36)),
									ast.COVARIANT,
									"T",
									nil,
									ast.NewPublicConstantNode(S(P(33, 1, 34), P(35, 1, 36)), "Foo"),
								),
								ast.NewVariantTypeParameterNode(
									S(P(38, 1, 39), P(46, 1, 47)),
									ast.CONTRAVARIANT,
									"Z",
									nil,
									ast.NewPrivateConstantNode(S(P(43, 1, 44), P(46, 1, 47)), "_Bar"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have type variables with lower bounds": {
			input: `struct Foo[V > Std::String, +T > Foo < _Bar]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(48, 1, 49)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(48, 1, 49)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(48, 1, 49)),
							"",
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(11, 1, 12), P(25, 1, 26)),
									ast.INVARIANT,
									"V",
									ast.NewConstantLookupNode(
										S(P(15, 1, 16), P(25, 1, 26)),
										ast.NewPublicConstantNode(S(P(15, 1, 16), P(17, 1, 18)), "Std"),
										ast.NewPublicConstantNode(S(P(20, 1, 21), P(25, 1, 26)), "String"),
									),
									nil,
								),
								ast.NewVariantTypeParameterNode(
									S(P(28, 1, 29), P(42, 1, 43)),
									ast.COVARIANT,
									"T",
									ast.NewPublicConstantNode(S(P(33, 1, 34), P(35, 1, 36)), "Foo"),
									ast.NewPrivateConstantNode(S(P(39, 1, 40), P(42, 1, 43)), "_Bar"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have type variables with fixed bounds": {
			input: `struct Foo[V = Std::String, +T > Foo < _Bar]; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(48, 1, 49)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(48, 1, 49)),
						ast.NewStructDeclarationNode(
							S(P(0, 1, 1), P(48, 1, 49)),
							"",
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(11, 1, 12), P(25, 1, 26)),
									ast.INVARIANT,
									"V",
									ast.NewConstantLookupNode(
										S(P(15, 1, 16), P(25, 1, 26)),
										ast.NewPublicConstantNode(S(P(15, 1, 16), P(17, 1, 18)), "Std"),
										ast.NewPublicConstantNode(S(P(20, 1, 21), P(25, 1, 26)), "String"),
									),
									ast.NewConstantLookupNode(
										S(P(15, 1, 16), P(25, 1, 26)),
										ast.NewPublicConstantNode(S(P(15, 1, 16), P(17, 1, 18)), "Std"),
										ast.NewPublicConstantNode(S(P(20, 1, 21), P(25, 1, 26)), "String"),
									),
								),
								ast.NewVariantTypeParameterNode(
									S(P(28, 1, 29), P(42, 1, 43)),
									ast.COVARIANT,
									"T",
									ast.NewPublicConstantNode(S(P(33, 1, 34), P(35, 1, 36)), "Foo"),
									ast.NewPrivateConstantNode(S(P(39, 1, 40), P(42, 1, 43)), "_Bar"),
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
							"",
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 1, 12), P(11, 1, 12)), "unexpected ], expected a list of type variables"),
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
							"",
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
							"",
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
							"",
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
							"",
							ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(7, 1, 8), P(9, 1, 10)), "invalid struct name, expected a constant"),
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
							"",
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							nil,
							[]ast.StructBodyStatementNode{
								ast.NewParameterStatementNode(
									S(P(13, 2, 3), P(16, 2, 6)),
									ast.NewAttributeParameterNode(
										S(P(13, 2, 3), P(15, 2, 5)),
										"foo",
										nil,
										nil,
									),
								),
								ast.NewParameterStatementNode(
									S(P(19, 3, 3), P(31, 3, 15)),
									ast.NewAttributeParameterNode(
										S(P(19, 3, 3), P(30, 3, 14)),
										"bar",
										ast.NewNilableTypeNode(
											S(P(24, 3, 8), P(30, 3, 14)),
											ast.NewPublicConstantNode(S(P(24, 3, 8), P(29, 3, 13)), "String"),
										),
										nil,
									),
								),
								ast.NewParameterStatementNode(
									S(P(34, 4, 3), P(47, 4, 16)),
									ast.NewAttributeParameterNode(
										S(P(34, 4, 3), P(46, 4, 15)),
										"baz",
										ast.NewPublicConstantNode(S(P(39, 4, 8), P(41, 4, 10)), "Int"),
										ast.NewFloatLiteralNode(S(P(45, 4, 14), P(46, 4, 15)), "0.3"),
									),
								),
								ast.NewParameterStatementNode(
									S(P(50, 5, 3), P(61, 5, 14)),
									ast.NewAttributeParameterNode(
										S(P(50, 5, 3), P(60, 5, 13)),
										"ban",
										nil,
										ast.NewRawStringLiteralNode(S(P(56, 5, 9), P(60, 5, 13)), "hey"),
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
							"",
							ast.NewPublicConstantNode(S(P(7, 1, 8), P(9, 1, 10)), "Foo"),
							nil,
							[]ast.StructBodyStatementNode{
								ast.NewParameterStatementNode(
									S(P(16, 1, 17), P(23, 1, 24)),
									ast.NewAttributeParameterNode(
										S(P(16, 1, 17), P(23, 1, 24)),
										"foo",
										ast.NewPublicConstantNode(S(P(21, 1, 22), P(23, 1, 24)), "Int"),
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

func TestMethodDefinition(t *testing.T) {
	tests := testTable{
		"cannot be a part of an expression": {
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
								"",
								false,
								false,
								"foo",
								nil,
								nil,
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(17, 1, 18)), "method definitions cannot appear in expressions"),
			},
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
							"",
							false,
							false,
							"foo",
							nil,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have type variables": {
			input: "def foo[V]; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							"",
							false,
							false,
							"foo",
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									ast.INVARIANT,
									"V",
									nil,
									nil,
								),
							},
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be sealed": {
			input: "sealed def foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							"",
							false,
							true,
							"foo",
							nil,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot repeat sealed": {
			input: "sealed sealed def foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							"",
							false,
							true,
							"foo",
							nil,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(5, 1, 6)), "the sealed modifier can only be attached once"),
			},
		},
		"cannot attach sealed to an abstract method": {
			input: "sealed abstract def foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							"",
							true,
							true,
							"foo",
							nil,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(5, 1, 6)), "the sealed modifier cannot be attached to abstract methods"),
			},
		},
		"can be abstract": {
			input: "abstract def foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							"",
							true,
							false,
							"foo",
							nil,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot repeat abstract": {
			input: "abstract abstract def foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							"",
							true,
							false,
							"foo",
							nil,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(7, 1, 8)), "the abstract modifier can only be attached once"),
			},
		},
		"cannot attach abstract to a sealed method": {
			input: "abstract sealed def foo; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							"",
							true,
							true,
							"foo",
							nil,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(7, 1, 8)), "the abstract modifier cannot be attached to sealed methods"),
			},
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
							"",
							false,
							false,
							"foo=",
							nil,
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
							"",
							false,
							false,
							"foo=",
							nil,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(13, 1, 14), P(18, 1, 19)), "setter methods cannot be defined with custom return types"),
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
							"",
							false,
							false,
							"foo=",
							nil,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 1, 13), P(15, 1, 16)), "setter methods must have a single parameter, got: 3"),
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
							"",
							false,
							false,
							"fo=",
							nil,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(6, 1, 7)), "setter methods must have a single parameter, got: 0"),
			},
		},
		"can have subscript setter with two arguments": {
			input: "def []=(k, v); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							"",
							false,
							false,
							"[]=",
							nil,
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"k",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
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
		"cannot have subscript setter with custom return type": {
			input: "def []=(k, v): Int; end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							"",
							false,
							false,
							"[]=",
							nil,
							[]ast.ParameterNode{
								ast.NewMethodParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"k",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewMethodParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"v",
									false,
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							ast.NewPublicConstantNode(
								S(P(15, 1, 16), P(17, 1, 18)),
								"Int",
							),
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(15, 1, 16), P(17, 1, 18)), "setter methods cannot be defined with custom return types"),
			},
		},
		"cannot define subscript setter with less parameters": {
			input: "def []=(v); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							"",
							false,
							false,
							"[]=",
							nil,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(8, 1, 9), P(8, 1, 9)), "subscript setter methods must have two parameters, got: 1"),
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
							"",
							false,
							false,
							"_foo",
							nil,
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
							"",
							false,
							false,
							"class",
							nil,
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
							"",
							false,
							false,
							"+",
							nil,
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
							"",
							false,
							false,
							"[]",
							nil,
							nil,
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
							"",
							false,
							false,
							"Foo",
							nil,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(6, 1, 7)), "unexpected PUBLIC_CONSTANT, expected a method name (identifier, overridable operator)"),
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
							"",
							false,
							false,
							"&&",
							nil,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(5, 1, 6)), "unexpected &&, expected a method name (identifier, overridable operator)"),
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
							"",
							false,
							false,
							"_Foo",
							nil,
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(7, 1, 8)), "unexpected PRIVATE_CONSTANT, expected a method name (identifier, overridable operator)"),
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(14, 1, 15), P(19, 1, 20)), "rest parameters cannot have default values"),
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 1, 19), P(22, 1, 23)), "optional parameters cannot appear after rest parameters"),
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
							"",
							false,
							false,
							"foo",
							nil,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 1, 19), P(19, 1, 20)), "there should be only a single positional rest parameter"),
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(14, 1, 15), P(20, 1, 21)), "rest parameters cannot have default values"),
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(19, 1, 20), P(19, 1, 20)), "named rest parameters should appear last"),
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(21, 1, 22), P(23, 1, 24)), "named rest parameters cannot appear after a post parameter"),
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
							"",
							false,
							false,
							"foo",
							nil,
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
		"can have never as a param type": {
			input: "def foo(a: Int, b: never); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							"",
							false,
							false,
							"foo",
							nil,
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
									S(P(16, 1, 17), P(23, 1, 24)),
									"b",
									false,
									ast.NewNeverTypeNode(
										S(P(19, 1, 20), P(23, 1, 24)),
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
		"cannot have void as a param type": {
			input: "def foo(a: Int, b: void); end",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 1, 29)),
						ast.NewMethodDefinitionNode(
							S(P(0, 1, 1), P(28, 1, 29)),
							"",
							false,
							false,
							"foo",
							nil,
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
									S(P(16, 1, 17), P(22, 1, 23)),
									"b",
									false,
									ast.NewVoidTypeNode(
										S(P(19, 1, 20), P(22, 1, 23)),
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(19, 1, 20), P(22, 1, 23)), "type `void` cannot be used in this context"),
			},
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 1, 17), P(24, 1, 25)), "required parameters cannot appear after optional parameters"),
				error.NewFailure(L("<main>", P(37, 1, 38), P(37, 1, 38)), "required parameters cannot appear after optional parameters"),
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
							"",
							false,
							false,
							"foo",
							nil,
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
							"",
							false,
							false,
							"foo",
							nil,
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
		"cannot be a part of an expression": {
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(14, 1, 15)), "method definitions cannot appear in expressions"),
			},
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(4, 1, 5)), "unexpected :, expected a statement separator `\\n`, `;`"),
				error.NewFailure(L("<main>", P(12, 1, 13), P(12, 1, 13)), "unexpected ?, expected a statement separator `\\n`, `;`"),
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(15, 1, 16), P(16, 1, 17)), "there should be only a single positional rest parameter"),
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 1, 17), P(16, 1, 17)), "named rest parameters should appear last"),
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(13, 1, 14), P(21, 1, 22)), "required parameters cannot appear after optional parameters"),
				error.NewFailure(L("<main>", P(34, 1, 35), P(34, 1, 35)), "required parameters cannot appear after optional parameters"),
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
		"cannot be a part of an expression": {
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
								"",
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(6, 1, 7), P(12, 1, 13)), "signature definitions cannot appear in expressions"),
			},
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
							"",
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
		"can have type parameters": {
			input: "sig foo[V]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							"",
							"foo",
							[]ast.TypeParameterNode{
								ast.NewVariantTypeParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									ast.INVARIANT,
									"V",
									nil,
									nil,
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
		"can have a private identifier as a name": {
			input: "sig _foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							"",
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
			input: "sig class",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							"",
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
			input: "sig +",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							"",
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
		"cannot have a public constant as a name": {
			input: "sig Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							"",
							"Foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(6, 1, 7)), "unexpected PUBLIC_CONSTANT, expected a method name (identifier, overridable operator)"),
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
							"",
							"&&",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(5, 1, 6)), "unexpected &&, expected a method name (identifier, overridable operator)"),
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
							"",
							"_Foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(4, 1, 5), P(7, 1, 8)), "unexpected PRIVATE_CONSTANT, expected a method name (identifier, overridable operator)"),
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
							"",
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
			input: "sig foo: String?",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							"",
							"foo",
							nil,
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
							"",
							"foo",
							nil,
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
							"",
							"foo",
							nil,
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
							"",
							"foo",
							nil,
							[]ast.ParameterNode{
								ast.NewSignatureParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									nil,
									false,
									ast.NormalParameterKind,
								),
								ast.NewSignatureParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"b",
									nil,
									false,
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
		"can have arguments with types": {
			input: "sig foo(a: Int, b: String?)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							"",
							"foo",
							nil,
							[]ast.ParameterNode{
								ast.NewSignatureParameterNode(
									S(P(8, 1, 9), P(13, 1, 14)),
									"a",
									ast.NewPublicConstantNode(S(P(11, 1, 12), P(13, 1, 14)), "Int"),
									false,
									ast.NormalParameterKind,
								),
								ast.NewSignatureParameterNode(
									S(P(16, 1, 17), P(25, 1, 26)),
									"b",
									ast.NewNilableTypeNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									),
									false,
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
		"can have optional arguments": {
			input: "sig foo(a, b?, c?: String?)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							"",
							"foo",
							nil,
							[]ast.ParameterNode{
								ast.NewSignatureParameterNode(
									S(P(8, 1, 9), P(8, 1, 9)),
									"a",
									nil,
									false,
									ast.NormalParameterKind,
								),
								ast.NewSignatureParameterNode(
									S(P(11, 1, 12), P(12, 1, 13)),
									"b",
									nil,
									true,
									ast.NormalParameterKind,
								),
								ast.NewSignatureParameterNode(
									S(P(15, 1, 16), P(25, 1, 26)),
									"c",
									ast.NewNilableTypeNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										ast.NewPublicConstantNode(S(P(19, 1, 20), P(24, 1, 25)), "String"),
									),
									true,
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
		"cannot have required parameters after optional ones": {
			input: "sig foo(a?, b, c?, d)",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewMethodSignatureDefinitionNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							"",
							"foo",
							nil,
							[]ast.ParameterNode{
								ast.NewSignatureParameterNode(
									S(P(8, 1, 9), P(9, 1, 10)),
									"a",
									nil,
									true,
									ast.NormalParameterKind,
								),
								ast.NewSignatureParameterNode(
									S(P(12, 1, 13), P(12, 1, 13)),
									"b",
									nil,
									false,
									ast.NormalParameterKind,
								),
								ast.NewSignatureParameterNode(
									S(P(15, 1, 16), P(16, 1, 17)),
									"c",
									nil,
									true,
									ast.NormalParameterKind,
								),
								ast.NewSignatureParameterNode(
									S(P(19, 1, 20), P(19, 1, 20)),
									"d",
									nil,
									false,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 1, 13), P(12, 1, 13)), "required parameters cannot appear after optional parameters"),
				error.NewFailure(L("<main>", P(19, 1, 20), P(19, 1, 20)), "required parameters cannot appear after optional parameters"),
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(10, 1, 11), P(10, 1, 11)), "unexpected =, expected )"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}
