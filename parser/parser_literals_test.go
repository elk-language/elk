package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
)

func TestStringLiteral(t *testing.T) {
	tests := testTable{
		"processes escape sequences": {
			input: `"foo\nbar\rbaz\\car\t\b\"\v\f\x12\a"`,
			want: ast.NewProgramNode(
				P(0, 36, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 36, 1, 1),
						ast.NewDoubleQuotedStringLiteralNode(
							P(0, 36, 1, 1),
							"foo\nbar\rbaz\\car\t\b\"\v\f\x12\a",
						),
					),
				},
			),
		},
		"reports errors for invalid hex escapes": {
			input: `"foo \xgh bar"`,
			want: ast.NewProgramNode(
				P(0, 14, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 14, 1, 1),
						ast.NewInterpolatedStringLiteralNode(
							P(0, 14, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 4, 1, 2), "foo "),
								ast.NewInvalidNode(P(5, 4, 1, 6), V(P(5, 4, 1, 6), token.ERROR, "invalid hex escape in string literal")),
								ast.NewStringLiteralContentSectionNode(P(9, 4, 1, 10), " bar"),
							},
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(5, 4, 1, 6), "invalid hex escape in string literal"),
			},
		},
		"reports errors for nonexistent escape sequences": {
			input: `"foo \q bar"`,
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 12, 1, 1),
						ast.NewInterpolatedStringLiteralNode(
							P(0, 12, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 4, 1, 2), "foo "),
								ast.NewInvalidNode(P(5, 2, 1, 6), V(P(5, 2, 1, 6), token.ERROR, "invalid escape sequence `\\q` in string literal")),
								ast.NewStringLiteralContentSectionNode(P(7, 4, 1, 8), " bar"),
							},
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(5, 2, 1, 6), "invalid escape sequence `\\q` in string literal"),
			},
		},
		"can contain interpolated expressions": {
			input: `"foo ${bar + 2} baz ${fudge}"`,
			want: ast.NewProgramNode(
				P(0, 29, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 29, 1, 1),
						ast.NewInterpolatedStringLiteralNode(
							P(0, 29, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 4, 1, 2), "foo "),
								ast.NewStringInterpolationNode(
									P(5, 10, 1, 6),
									ast.NewBinaryExpressionNode(
										P(7, 7, 1, 8),
										T(P(11, 1, 1, 12), token.PLUS),
										ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
										ast.NewIntLiteralNode(P(13, 1, 1, 14), V(P(13, 1, 1, 14), token.DEC_INT, "2")),
									),
								),
								ast.NewStringLiteralContentSectionNode(P(15, 5, 1, 16), " baz "),
								ast.NewStringInterpolationNode(
									P(20, 8, 1, 21),
									ast.NewPublicIdentifierNode(P(22, 5, 1, 23), "fudge"),
								),
							},
						),
					),
				},
			),
		},
		"can't contain string literals inside interpolation": {
			input: `"foo ${"bar" + 2} baza"`,
			want: ast.NewProgramNode(
				P(0, 23, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 23, 1, 1),
						ast.NewInterpolatedStringLiteralNode(
							P(0, 23, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 4, 1, 2), "foo "),
								ast.NewStringInterpolationNode(
									P(5, 12, 1, 6),
									ast.NewBinaryExpressionNode(
										P(7, 9, 1, 8),
										T(P(13, 1, 1, 14), token.PLUS),
										ast.NewInvalidNode(P(7, 5, 1, 8), V(P(7, 5, 1, 8), token.ERROR, "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation")),
										ast.NewIntLiteralNode(P(15, 1, 1, 16), V(P(15, 1, 1, 16), token.DEC_INT, "2")),
									),
								),
								ast.NewStringLiteralContentSectionNode(P(17, 5, 1, 18), " baza"),
							},
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(7, 5, 1, 8), "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation"),
			},
		},
		"can contain raw string literals inside interpolation": {
			input: `"foo ${'bar' + 2} baza"`,
			want: ast.NewProgramNode(
				P(0, 23, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 23, 1, 1),
						ast.NewInterpolatedStringLiteralNode(
							P(0, 23, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 4, 1, 2), "foo "),
								ast.NewStringInterpolationNode(
									P(5, 12, 1, 6),
									ast.NewBinaryExpressionNode(
										P(7, 9, 1, 8),
										T(P(13, 1, 1, 14), token.PLUS),
										ast.NewRawStringLiteralNode(P(7, 5, 1, 8), "bar"),
										ast.NewIntLiteralNode(P(15, 1, 1, 16), V(P(15, 1, 1, 16), token.DEC_INT, "2")),
									),
								),
								ast.NewStringLiteralContentSectionNode(P(17, 5, 1, 18), " baza"),
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

func TestRawStringLiteral(t *testing.T) {
	tests := testTable{
		"doesn't process escape sequences": {
			input: `'foo\nbar\rbaz\\car\t\b\"\v\f\x12\a'`,
			want: ast.NewProgramNode(
				P(0, 36, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 36, 1, 1),
						ast.NewRawStringLiteralNode(P(0, 36, 1, 1), `foo\nbar\rbaz\\car\t\b\"\v\f\x12\a`),
					),
				},
			),
		},
		"can't contain interpolated expressions": {
			input: `'foo ${bar + 2} baz ${fudge}'`,
			want: ast.NewProgramNode(
				P(0, 29, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 29, 1, 1),
						ast.NewRawStringLiteralNode(P(0, 29, 1, 1), `foo ${bar + 2} baz ${fudge}`),
					),
				},
			),
		},
		"can contain double quotes": {
			input: `'foo ${"bar" + 2} baza'`,
			want: ast.NewProgramNode(
				P(0, 23, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 23, 1, 1),
						ast.NewRawStringLiteralNode(P(0, 23, 1, 1), `foo ${"bar" + 2} baza`),
					),
				},
			),
		},
		"doesn't allow escaping single quotes": {
			input: `'foo\'s house'`,
			want: ast.NewProgramNode(
				P(0, 14, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewRawStringLiteralNode(P(0, 6, 1, 1), "foo\\"),
					),
				},
			),
			err: ErrorList{
				NewError(P(6, 1, 1, 7), "unexpected PUBLIC_IDENTIFIER, expected a statement separator `\\n`, `;`"),
				NewError(P(13, 1, 1, 14), "unterminated raw string literal, missing `'`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestClosure(t *testing.T) {
	tests := testTable{
		"can have arguments and be single line": {
			input: `|a| -> 'foo' + .2`,
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 17, 1, 1),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									P(1, 1, 1, 2),
									"a",
									nil,
									nil,
								),
							},
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(7, 10, 1, 8),
									ast.NewBinaryExpressionNode(
										P(7, 10, 1, 8),
										T(P(13, 1, 1, 14), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(7, 5, 1, 8),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(15, 2, 1, 16),
											"0.2",
										),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have arguments and be single line with braces": {
			input: `|a| -> { 'foo' + .2 }`,
			want: ast.NewProgramNode(
				P(0, 21, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 21, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 21, 1, 1),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									P(1, 1, 1, 2),
									"a",
									nil,
									nil,
								),
							},
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(9, 10, 1, 10),
									ast.NewBinaryExpressionNode(
										P(9, 10, 1, 10),
										T(P(15, 1, 1, 16), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(9, 5, 1, 10),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(17, 2, 1, 18),
											"0.2",
										),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have arguments and be multiline with braces": {
			input: `|a| -> {
	'foo' + .2
	nil
}`,
			want: ast.NewProgramNode(
				P(0, 27, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 27, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 27, 1, 1),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									P(1, 1, 1, 2),
									"a",
									nil,
									nil,
								),
							},
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(10, 11, 2, 2),
									ast.NewBinaryExpressionNode(
										P(10, 10, 2, 2),
										T(P(16, 1, 2, 8), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(10, 5, 2, 2),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(18, 2, 2, 10),
											"0.2",
										),
									),
								),
								ast.NewExpressionStatementNode(
									P(22, 4, 3, 2),
									ast.NewNilLiteralNode(P(22, 3, 3, 2)),
								),
							},
						),
					),
				},
			),
		},
		"can have arguments and be multiline with end": {
			input: `|a| ->
	'foo' + .2
	nil
end`,
			want: ast.NewProgramNode(
				P(0, 27, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 27, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 27, 1, 1),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									P(1, 1, 1, 2),
									"a",
									nil,
									nil,
								),
							},
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(8, 11, 2, 2),
									ast.NewBinaryExpressionNode(
										P(8, 10, 2, 2),
										T(P(14, 1, 2, 8), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(8, 5, 2, 2),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(16, 2, 2, 10),
											"0.2",
										),
									),
								),
								ast.NewExpressionStatementNode(
									P(20, 4, 3, 2),
									ast.NewNilLiteralNode(P(20, 3, 3, 2)),
								),
							},
						),
					),
				},
			),
		},
		"can omit pipes when there's a single argument": {
			input: `a -> 'foo' + .2`,
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 15, 1, 1),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									P(0, 1, 1, 1),
									"a",
									nil,
									nil,
								),
							},
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(5, 10, 1, 6),
									ast.NewBinaryExpressionNode(
										P(5, 10, 1, 6),
										T(P(11, 1, 1, 12), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(5, 5, 1, 6),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(13, 2, 1, 14),
											"0.2",
										),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have arguments with types": {
			input: `|a: Int, b: String| -> 'foo' + .2`,
			want: ast.NewProgramNode(
				P(0, 33, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 33, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 33, 1, 1),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									P(1, 6, 1, 2),
									"a",
									ast.NewPublicConstantNode(P(4, 3, 1, 5), "Int"),
									nil,
								),
								ast.NewFormalParameterNode(
									P(9, 9, 1, 10),
									"b",
									ast.NewPublicConstantNode(P(12, 6, 1, 13), "String"),
									nil,
								),
							},
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(23, 10, 1, 24),
									ast.NewBinaryExpressionNode(
										P(23, 10, 1, 24),
										T(P(29, 1, 1, 30), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(23, 5, 1, 24),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(31, 2, 1, 32),
											"0.2",
										),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have arguments with initialisers": {
			input: `|a = 32, b: String = 'foo'| -> 'foo' + .2`,
			want: ast.NewProgramNode(
				P(0, 41, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 41, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 41, 1, 1),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									P(1, 6, 1, 2),
									"a",
									nil,
									ast.NewIntLiteralNode(P(5, 2, 1, 6), V(P(5, 2, 1, 6), token.DEC_INT, "32")),
								),
								ast.NewFormalParameterNode(
									P(9, 17, 1, 10),
									"b",
									ast.NewPublicConstantNode(P(12, 6, 1, 13), "String"),
									ast.NewRawStringLiteralNode(P(21, 5, 1, 22), "foo"),
								),
							},
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(31, 10, 1, 32),
									ast.NewBinaryExpressionNode(
										P(31, 10, 1, 32),
										T(P(37, 1, 1, 38), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(31, 5, 1, 32),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(39, 2, 1, 40),
											"0.2",
										),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty argument list": {
			input: `|| -> 'foo' + .2`,
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 16, 1, 1),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(6, 10, 1, 7),
									ast.NewBinaryExpressionNode(
										P(6, 10, 1, 7),
										T(P(12, 1, 1, 13), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(6, 5, 1, 7),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(14, 2, 1, 15),
											"0.2",
										),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a return type": {
			input: `||: String? -> 'foo' + .2`,
			want: ast.NewProgramNode(
				P(0, 25, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 25, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 25, 1, 1),
							nil,
							ast.NewNilableTypeNode(P(4, 7, 1, 5), ast.NewPublicConstantNode(P(4, 6, 1, 5), "String")),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(15, 10, 1, 16),
									ast.NewBinaryExpressionNode(
										P(15, 10, 1, 16),
										T(P(21, 1, 1, 22), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(15, 5, 1, 16),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(23, 2, 1, 24),
											"0.2",
										),
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

func TestSymbolLiteral(t *testing.T) {
	tests := testTable{
		"can have spaces between the colon and the content": {
			input: ": foo",
			want: ast.NewProgramNode(
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewSimpleSymbolLiteralNode(P(0, 5, 1, 1), "foo"),
					),
				},
			),
		},
		"can have a public identifier as the content": {
			input: ":foo",
			want: ast.NewProgramNode(
				P(0, 4, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewSimpleSymbolLiteralNode(P(0, 4, 1, 1), "foo"),
					),
				},
			),
		},
		"can have a private identifier as the content": {
			input: ":_foo",
			want: ast.NewProgramNode(
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewSimpleSymbolLiteralNode(P(0, 5, 1, 1), "_foo"),
					),
				},
			),
		},
		"can have a public constant as the content": {
			input: ":Foo",
			want: ast.NewProgramNode(
				P(0, 4, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewSimpleSymbolLiteralNode(P(0, 4, 1, 1), "Foo"),
					),
				},
			),
		},
		"can have a private constant as the content": {
			input: ":_Foo",
			want: ast.NewProgramNode(
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewSimpleSymbolLiteralNode(P(0, 5, 1, 1), "_Foo"),
					),
				},
			),
		},
		"can have a keyword as the content": {
			input: ":var",
			want: ast.NewProgramNode(
				P(0, 4, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewSimpleSymbolLiteralNode(P(0, 4, 1, 1), "var"),
					),
				},
			),
		},
		"can have a raw string as the content": {
			input: ":'foo bar'",
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewSimpleSymbolLiteralNode(P(0, 10, 1, 1), "foo bar"),
					),
				},
			),
		},
		"can have a double quoted string as the content": {
			input: `:"foo bar"`,
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewSimpleSymbolLiteralNode(P(0, 10, 1, 1), "foo bar"),
					),
				},
			),
		},
		"can have an overridable operator as the content": {
			input: ":+",
			want: ast.NewProgramNode(
				P(0, 2, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewSimpleSymbolLiteralNode(P(0, 2, 1, 1), "+"),
					),
				},
			),
		},
		"can't have a not overridable operator as the content": {
			input: ":&&",
			want: ast.NewProgramNode(
				P(0, 3, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 3, 1, 1),
						ast.NewInvalidNode(P(0, 3, 1, 1), T(P(1, 2, 1, 2), token.AND_AND)),
					),
				},
			),
			err: ErrorList{
				NewError(P(1, 2, 1, 2), "unexpected &&, expected an identifier, overridable operator or string literal"),
			},
		},
		"can have a string as the content": {
			input: `:"foo ${bar}"`,
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewInterpolatedSymbolLiteral(
							P(0, 13, 1, 1),
							ast.NewInterpolatedStringLiteralNode(
								P(1, 12, 1, 2),
								[]ast.StringLiteralContentNode{
									ast.NewStringLiteralContentSectionNode(
										P(2, 4, 1, 3),
										"foo ",
									),
									ast.NewStringInterpolationNode(
										P(6, 6, 1, 7),
										ast.NewPublicIdentifierNode(P(8, 3, 1, 9), "bar"),
									),
								},
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

func TestNamedValueLiteral(t *testing.T) {
	tests := testTable{
		"can have spaces between the colon and the name": {
			input: ": foo{.5}",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewNamedValueLiteralNode(
							P(0, 9, 1, 1),
							"foo",
							ast.NewFloatLiteralNode(P(6, 2, 1, 7), "0.5"),
						),
					),
				},
			),
		},
		"can have a public identifier as the name": {
			input: ":foo{.5}",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewNamedValueLiteralNode(
							P(0, 8, 1, 1),
							"foo",
							ast.NewFloatLiteralNode(P(5, 2, 1, 6), "0.5"),
						),
					),
				},
			),
		},
		"can have an expression as the value": {
			input: ":foo{.5 + 'hej'}",
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewNamedValueLiteralNode(
							P(0, 16, 1, 1),
							"foo",
							ast.NewBinaryExpressionNode(
								P(5, 10, 1, 6),
								T(P(8, 1, 1, 9), token.PLUS),
								ast.NewFloatLiteralNode(P(5, 2, 1, 6), "0.5"),
								ast.NewRawStringLiteralNode(P(10, 5, 1, 11), "hej"),
							),
						),
					),
				},
			),
		},
		"can omit the value": {
			input: ":foo{}",
			want: ast.NewProgramNode(
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewNamedValueLiteralNode(
							P(0, 6, 1, 1),
							"foo",
							nil,
						),
					),
				},
			),
		},
		"can have a private identifier as the name": {
			input: ":_foo{.5}",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewNamedValueLiteralNode(
							P(0, 9, 1, 1),
							"_foo",
							ast.NewFloatLiteralNode(P(6, 2, 1, 7), "0.5"),
						),
					),
				},
			),
		},
		"can have a public constant as the name": {
			input: ":Foo{.5}",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewNamedValueLiteralNode(
							P(0, 8, 1, 1),
							"Foo",
							ast.NewFloatLiteralNode(P(5, 2, 1, 6), "0.5"),
						),
					),
				},
			),
		},
		"can have a private constant as the name": {
			input: ":_Foo{}",
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewNamedValueLiteralNode(
							P(0, 7, 1, 1),
							"_Foo",
							nil,
						),
					),
				},
			),
		},
		"can have a raw string as the name": {
			input: ":'foo bar'{}",
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 12, 1, 1),
						ast.NewNamedValueLiteralNode(
							P(0, 12, 1, 1),
							"foo bar",
							nil,
						),
					),
				},
			),
		},
		"can have an overridable operator as the name": {
			input: ":+{}",
			want: ast.NewProgramNode(
				P(0, 4, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewNamedValueLiteralNode(
							P(0, 4, 1, 1),
							"+",
							nil,
						),
					),
				},
			),
		},
		"can't have a not overridable operator as the name": {
			input: ":&&{}",
			want: ast.NewProgramNode(
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewInvalidNode(P(0, 3, 1, 1), T(P(1, 2, 1, 2), token.AND_AND)),
					),
				},
			),
			err: ErrorList{
				NewError(P(1, 2, 1, 2), "unexpected &&, expected an identifier, overridable operator or string literal"),
			},
		},
		"can't have a string as the name": {
			input: `:"foo ${bar}"{}`,
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewInterpolatedSymbolLiteral(
							P(0, 13, 1, 1),
							ast.NewInterpolatedStringLiteralNode(
								P(1, 12, 1, 2),
								[]ast.StringLiteralContentNode{
									ast.NewStringLiteralContentSectionNode(
										P(2, 4, 1, 3),
										"foo ",
									),
									ast.NewStringInterpolationNode(
										P(6, 6, 1, 7),
										ast.NewPublicIdentifierNode(P(8, 3, 1, 9), "bar"),
									),
								},
							),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(13, 1, 1, 14), "unexpected {, expected a statement separator `\\n`, `;`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestListLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "[]",
			want: ast.NewProgramNode(
				P(0, 2, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewListLiteralNode(
							P(0, 2, 1, 1),
							nil,
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "[\n\n]",
			want: ast.NewProgramNode(
				P(0, 4, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewListLiteralNode(
							P(0, 4, 1, 1),
							nil,
						),
					),
				},
			),
		},
		"can contain if modifiers": {
			input: "[.1, 'foo', :bar, baz + 5 if baz]",
			want: ast.NewProgramNode(
				P(0, 33, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 33, 1, 1),
						ast.NewListLiteralNode(
							P(0, 33, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(1, 2, 1, 2), "0.1"),
								ast.NewRawStringLiteralNode(P(5, 5, 1, 6), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(12, 4, 1, 13), "bar"),
								ast.NewModifierNode(
									P(18, 14, 1, 19),
									T(P(26, 2, 1, 27), token.IF),
									ast.NewBinaryExpressionNode(
										P(18, 7, 1, 19),
										T(P(22, 1, 1, 23), token.PLUS),
										ast.NewPublicIdentifierNode(P(18, 3, 1, 19), "baz"),
										ast.NewIntLiteralNode(P(24, 1, 1, 25), V(P(24, 1, 1, 25), token.DEC_INT, "5")),
									),
									ast.NewPublicIdentifierNode(P(29, 3, 1, 30), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can contain unless modifiers": {
			input: "[.1, 'foo', :bar, baz + 5 unless baz]",
			want: ast.NewProgramNode(
				P(0, 37, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 37, 1, 1),
						ast.NewListLiteralNode(
							P(0, 37, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(1, 2, 1, 2), "0.1"),
								ast.NewRawStringLiteralNode(P(5, 5, 1, 6), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(12, 4, 1, 13), "bar"),
								ast.NewModifierNode(
									P(18, 18, 1, 19),
									T(P(26, 6, 1, 27), token.UNLESS),
									ast.NewBinaryExpressionNode(
										P(18, 7, 1, 19),
										T(P(22, 1, 1, 23), token.PLUS),
										ast.NewPublicIdentifierNode(P(18, 3, 1, 19), "baz"),
										ast.NewIntLiteralNode(P(24, 1, 1, 25), V(P(24, 1, 1, 25), token.DEC_INT, "5")),
									),
									ast.NewPublicIdentifierNode(P(33, 3, 1, 34), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can contain for modifiers": {
			input: "[.1, 'foo', :bar, baz + 5 for baz in bazz]",
			want: ast.NewProgramNode(
				P(0, 42, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 42, 1, 1),
						ast.NewListLiteralNode(
							P(0, 42, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(1, 2, 1, 2), "0.1"),
								ast.NewRawStringLiteralNode(P(5, 5, 1, 6), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(12, 4, 1, 13), "bar"),
								ast.NewModifierForInNode(
									P(18, 23, 1, 19),
									ast.NewBinaryExpressionNode(
										P(18, 7, 1, 19),
										T(P(22, 1, 1, 23), token.PLUS),
										ast.NewPublicIdentifierNode(P(18, 3, 1, 19), "baz"),
										ast.NewIntLiteralNode(P(24, 1, 1, 25), V(P(24, 1, 1, 25), token.DEC_INT, "5")),
									),
									[]ast.ParameterNode{
										ast.NewLoopParameterNode(P(30, 3, 1, 31), "baz", nil),
									},
									ast.NewPublicIdentifierNode(P(37, 4, 1, 38), "bazz"),
								),
							},
						),
					),
				},
			),
		},
		"can have elements": {
			input: "[.1, 'foo', :bar, baz + 5]",
			want: ast.NewProgramNode(
				P(0, 26, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 26, 1, 1),
						ast.NewListLiteralNode(
							P(0, 26, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(1, 2, 1, 2), "0.1"),
								ast.NewRawStringLiteralNode(P(5, 5, 1, 6), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(12, 4, 1, 13), "bar"),
								ast.NewBinaryExpressionNode(
									P(18, 7, 1, 19),
									T(P(22, 1, 1, 23), token.PLUS),
									ast.NewPublicIdentifierNode(P(18, 3, 1, 19), "baz"),
									ast.NewIntLiteralNode(P(24, 1, 1, 25), V(P(24, 1, 1, 25), token.DEC_INT, "5")),
								),
							},
						),
					),
				},
			),
		},
		"can have explicit indices": {
			input: "[.1, 'foo', 10 => :bar, baz => baz + 5]",
			want: ast.NewProgramNode(
				P(0, 39, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 39, 1, 1),
						ast.NewListLiteralNode(
							P(0, 39, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(1, 2, 1, 2), "0.1"),
								ast.NewRawStringLiteralNode(P(5, 5, 1, 6), "foo"),
								ast.NewKeyValueExpressionNode(
									P(12, 10, 1, 13),
									ast.NewIntLiteralNode(P(12, 2, 1, 13), V(P(12, 2, 1, 13), token.DEC_INT, "10")),
									ast.NewSimpleSymbolLiteralNode(P(18, 4, 1, 19), "bar"),
								),
								ast.NewKeyValueExpressionNode(
									P(24, 14, 1, 25),
									ast.NewPublicIdentifierNode(P(24, 3, 1, 25), "baz"),
									ast.NewBinaryExpressionNode(
										P(31, 7, 1, 32),
										T(P(35, 1, 1, 36), token.PLUS),
										ast.NewPublicIdentifierNode(P(31, 3, 1, 32), "baz"),
										ast.NewIntLiteralNode(P(37, 1, 1, 38), V(P(37, 1, 1, 38), token.DEC_INT, "5")),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have explicit indices with modifiers": {
			input: "[.1, 'foo', 10 => :bar if bar, baz => baz + 5 for baz in bazz]",
			want: ast.NewProgramNode(
				P(0, 62, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 62, 1, 1),
						ast.NewListLiteralNode(
							P(0, 62, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(1, 2, 1, 2), "0.1"),
								ast.NewRawStringLiteralNode(P(5, 5, 1, 6), "foo"),
								ast.NewModifierNode(
									P(12, 17, 1, 13),
									T(P(23, 2, 1, 24), token.IF),
									ast.NewKeyValueExpressionNode(
										P(12, 10, 1, 13),
										ast.NewIntLiteralNode(P(12, 2, 1, 13), V(P(12, 2, 1, 13), token.DEC_INT, "10")),
										ast.NewSimpleSymbolLiteralNode(P(18, 4, 1, 19), "bar"),
									),
									ast.NewPublicIdentifierNode(P(26, 3, 1, 27), "bar"),
								),
								ast.NewModifierForInNode(
									P(31, 30, 1, 32),
									ast.NewKeyValueExpressionNode(
										P(31, 14, 1, 32),
										ast.NewPublicIdentifierNode(P(31, 3, 1, 32), "baz"),
										ast.NewBinaryExpressionNode(
											P(38, 7, 1, 39),
											T(P(42, 1, 1, 43), token.PLUS),
											ast.NewPublicIdentifierNode(P(38, 3, 1, 39), "baz"),
											ast.NewIntLiteralNode(P(44, 1, 1, 45), V(P(44, 1, 1, 45), token.DEC_INT, "5")),
										),
									),
									[]ast.ParameterNode{
										ast.NewLoopParameterNode(P(50, 3, 1, 51), "baz", nil),
									},
									ast.NewPublicIdentifierNode(P(57, 4, 1, 58), "bazz"),
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: "[\n.1\n,\n'foo'\n,\n:bar\n,\nbaz + 5\n]",
			want: ast.NewProgramNode(
				P(0, 31, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 31, 1, 1),
						ast.NewListLiteralNode(
							P(0, 31, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(2, 2, 2, 1), "0.1"),
								ast.NewRawStringLiteralNode(P(7, 5, 4, 1), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(15, 4, 6, 1), "bar"),
								ast.NewBinaryExpressionNode(
									P(22, 7, 8, 1),
									T(P(26, 1, 8, 5), token.PLUS),
									ast.NewPublicIdentifierNode(P(22, 3, 8, 1), "baz"),
									ast.NewIntLiteralNode(P(28, 1, 8, 7), V(P(28, 1, 8, 7), token.DEC_INT, "5")),
								),
							},
						),
					),
				},
			),
		},
		"can be nested": {
			input: "[[.1, :+], .2]",
			want: ast.NewProgramNode(
				P(0, 14, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 14, 1, 1),
						ast.NewListLiteralNode(
							P(0, 14, 1, 1),
							[]ast.ExpressionNode{
								ast.NewListLiteralNode(
									P(1, 8, 1, 2),
									[]ast.ExpressionNode{
										ast.NewFloatLiteralNode(P(2, 2, 1, 3), "0.1"),
										ast.NewSimpleSymbolLiteralNode(P(6, 2, 1, 7), "+"),
									},
								),
								ast.NewFloatLiteralNode(P(11, 2, 1, 12), "0.2"),
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

func TestTupleLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "%()",
			want: ast.NewProgramNode(
				P(0, 3, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 3, 1, 1),
						ast.NewTupleLiteralNode(
							P(0, 3, 1, 1),
							nil,
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "%(\n\n)",
			want: ast.NewProgramNode(
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewTupleLiteralNode(
							P(0, 5, 1, 1),
							nil,
						),
					),
				},
			),
		},
		"can contain if modifiers": {
			input: "%(.1, 'foo', :bar, baz + 5 if baz)",
			want: ast.NewProgramNode(
				P(0, 34, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 34, 1, 1),
						ast.NewTupleLiteralNode(
							P(0, 34, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(2, 2, 1, 3), "0.1"),
								ast.NewRawStringLiteralNode(P(6, 5, 1, 7), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(13, 4, 1, 14), "bar"),
								ast.NewModifierNode(
									P(19, 14, 1, 20),
									T(P(27, 2, 1, 28), token.IF),
									ast.NewBinaryExpressionNode(
										P(19, 7, 1, 20),
										T(P(23, 1, 1, 24), token.PLUS),
										ast.NewPublicIdentifierNode(P(19, 3, 1, 20), "baz"),
										ast.NewIntLiteralNode(P(25, 1, 1, 26), V(P(25, 1, 1, 26), token.DEC_INT, "5")),
									),
									ast.NewPublicIdentifierNode(P(30, 3, 1, 31), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can contain unless modifiers": {
			input: "%(.1, 'foo', :bar, baz + 5 unless baz)",
			want: ast.NewProgramNode(
				P(0, 38, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 38, 1, 1),
						ast.NewTupleLiteralNode(
							P(0, 38, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(2, 2, 1, 3), "0.1"),
								ast.NewRawStringLiteralNode(P(6, 5, 1, 7), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(13, 4, 1, 14), "bar"),
								ast.NewModifierNode(
									P(19, 18, 1, 20),
									T(P(27, 6, 1, 28), token.UNLESS),
									ast.NewBinaryExpressionNode(
										P(19, 7, 1, 20),
										T(P(23, 1, 1, 24), token.PLUS),
										ast.NewPublicIdentifierNode(P(19, 3, 1, 20), "baz"),
										ast.NewIntLiteralNode(P(25, 1, 1, 26), V(P(25, 1, 1, 26), token.DEC_INT, "5")),
									),
									ast.NewPublicIdentifierNode(P(34, 3, 1, 35), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can contain for modifiers": {
			input: "%(.1, 'foo', :bar, baz + 5 for baz in bazz)",
			want: ast.NewProgramNode(
				P(0, 43, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 43, 1, 1),
						ast.NewTupleLiteralNode(
							P(0, 43, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(2, 2, 1, 3), "0.1"),
								ast.NewRawStringLiteralNode(P(6, 5, 1, 7), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(13, 4, 1, 14), "bar"),
								ast.NewModifierForInNode(
									P(19, 23, 1, 20),
									ast.NewBinaryExpressionNode(
										P(19, 7, 1, 20),
										T(P(23, 1, 1, 24), token.PLUS),
										ast.NewPublicIdentifierNode(P(19, 3, 1, 20), "baz"),
										ast.NewIntLiteralNode(P(25, 1, 1, 26), V(P(25, 1, 1, 26), token.DEC_INT, "5")),
									),
									[]ast.ParameterNode{
										ast.NewLoopParameterNode(P(31, 3, 1, 32), "baz", nil),
									},
									ast.NewPublicIdentifierNode(P(38, 4, 1, 39), "bazz"),
								),
							},
						),
					),
				},
			),
		},
		"can have elements": {
			input: "%(.1, 'foo', :bar, baz + 5)",
			want: ast.NewProgramNode(
				P(0, 27, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 27, 1, 1),
						ast.NewTupleLiteralNode(
							P(0, 27, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(2, 2, 1, 3), "0.1"),
								ast.NewRawStringLiteralNode(P(6, 5, 1, 7), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(13, 4, 1, 14), "bar"),
								ast.NewBinaryExpressionNode(
									P(19, 7, 1, 20),
									T(P(23, 1, 1, 24), token.PLUS),
									ast.NewPublicIdentifierNode(P(19, 3, 1, 20), "baz"),
									ast.NewIntLiteralNode(P(25, 1, 1, 26), V(P(25, 1, 1, 26), token.DEC_INT, "5")),
								),
							},
						),
					),
				},
			),
		},
		"can have explicit indices": {
			input: "%(.1, 'foo', 10 => :bar, baz => baz + 5)",
			want: ast.NewProgramNode(
				P(0, 40, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 40, 1, 1),
						ast.NewTupleLiteralNode(
							P(0, 40, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(2, 2, 1, 3), "0.1"),
								ast.NewRawStringLiteralNode(P(6, 5, 1, 7), "foo"),
								ast.NewKeyValueExpressionNode(
									P(13, 10, 1, 14),
									ast.NewIntLiteralNode(P(13, 2, 1, 14), V(P(13, 2, 1, 14), token.DEC_INT, "10")),
									ast.NewSimpleSymbolLiteralNode(P(19, 4, 1, 20), "bar"),
								),
								ast.NewKeyValueExpressionNode(
									P(25, 14, 1, 26),
									ast.NewPublicIdentifierNode(P(25, 3, 1, 26), "baz"),
									ast.NewBinaryExpressionNode(
										P(32, 7, 1, 33),
										T(P(36, 1, 1, 37), token.PLUS),
										ast.NewPublicIdentifierNode(P(32, 3, 1, 33), "baz"),
										ast.NewIntLiteralNode(P(38, 1, 1, 39), V(P(38, 1, 1, 39), token.DEC_INT, "5")),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have explicit indices with modifiers": {
			input: "%(.1, 'foo', 10 => :bar if bar, baz => baz + 5 for baz in bazz)",
			want: ast.NewProgramNode(
				P(0, 63, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 63, 1, 1),
						ast.NewTupleLiteralNode(
							P(0, 63, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(2, 2, 1, 3), "0.1"),
								ast.NewRawStringLiteralNode(P(6, 5, 1, 7), "foo"),
								ast.NewModifierNode(
									P(13, 17, 1, 14),
									T(P(24, 2, 1, 25), token.IF),
									ast.NewKeyValueExpressionNode(
										P(13, 10, 1, 14),
										ast.NewIntLiteralNode(P(13, 2, 1, 14), V(P(13, 2, 1, 14), token.DEC_INT, "10")),
										ast.NewSimpleSymbolLiteralNode(P(19, 4, 1, 20), "bar"),
									),
									ast.NewPublicIdentifierNode(P(27, 3, 1, 28), "bar"),
								),
								ast.NewModifierForInNode(
									P(32, 30, 1, 33),
									ast.NewKeyValueExpressionNode(
										P(32, 14, 1, 33),
										ast.NewPublicIdentifierNode(P(32, 3, 1, 33), "baz"),
										ast.NewBinaryExpressionNode(
											P(39, 7, 1, 40),
											T(P(43, 1, 1, 44), token.PLUS),
											ast.NewPublicIdentifierNode(P(39, 3, 1, 40), "baz"),
											ast.NewIntLiteralNode(P(45, 1, 1, 46), V(P(45, 1, 1, 46), token.DEC_INT, "5")),
										),
									),
									[]ast.ParameterNode{
										ast.NewLoopParameterNode(P(51, 3, 1, 52), "baz", nil),
									},
									ast.NewPublicIdentifierNode(P(58, 4, 1, 59), "bazz"),
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: "%(\n.1\n,\n'foo'\n,\n:bar\n,\nbaz + 5\n)",
			want: ast.NewProgramNode(
				P(0, 32, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 32, 1, 1),
						ast.NewTupleLiteralNode(
							P(0, 32, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(3, 2, 2, 1), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 4, 1), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(16, 4, 6, 1), "bar"),
								ast.NewBinaryExpressionNode(
									P(23, 7, 8, 1),
									T(P(27, 1, 8, 5), token.PLUS),
									ast.NewPublicIdentifierNode(P(23, 3, 8, 1), "baz"),
									ast.NewIntLiteralNode(P(29, 1, 8, 7), V(P(29, 1, 8, 7), token.DEC_INT, "5")),
								),
							},
						),
					),
				},
			),
		},
		"can be nested": {
			input: "%(%(.1, :+), .2)",
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewTupleLiteralNode(
							P(0, 16, 1, 1),
							[]ast.ExpressionNode{
								ast.NewTupleLiteralNode(
									P(2, 9, 1, 3),
									[]ast.ExpressionNode{
										ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
										ast.NewSimpleSymbolLiteralNode(P(8, 2, 1, 9), "+"),
									},
								),
								ast.NewFloatLiteralNode(P(13, 2, 1, 14), "0.2"),
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

func TestSetLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "%{}",
			want: ast.NewProgramNode(
				P(0, 3, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 3, 1, 1),
						ast.NewSetLiteralNode(
							P(0, 3, 1, 1),
							nil,
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "%{\n\n}",
			want: ast.NewProgramNode(
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewSetLiteralNode(
							P(0, 5, 1, 1),
							nil,
						),
					),
				},
			),
		},
		"can contain if modifiers": {
			input: "%{.1, 'foo', :bar, baz + 5 if baz}",
			want: ast.NewProgramNode(
				P(0, 34, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 34, 1, 1),
						ast.NewSetLiteralNode(
							P(0, 34, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(2, 2, 1, 3), "0.1"),
								ast.NewRawStringLiteralNode(P(6, 5, 1, 7), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(13, 4, 1, 14), "bar"),
								ast.NewModifierNode(
									P(19, 14, 1, 20),
									T(P(27, 2, 1, 28), token.IF),
									ast.NewBinaryExpressionNode(
										P(19, 7, 1, 20),
										T(P(23, 1, 1, 24), token.PLUS),
										ast.NewPublicIdentifierNode(P(19, 3, 1, 20), "baz"),
										ast.NewIntLiteralNode(P(25, 1, 1, 26), V(P(25, 1, 1, 26), token.DEC_INT, "5")),
									),
									ast.NewPublicIdentifierNode(P(30, 3, 1, 31), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can contain unless modifiers": {
			input: "%{.1, 'foo', :bar, baz + 5 unless baz}",
			want: ast.NewProgramNode(
				P(0, 38, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 38, 1, 1),
						ast.NewSetLiteralNode(
							P(0, 38, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(2, 2, 1, 3), "0.1"),
								ast.NewRawStringLiteralNode(P(6, 5, 1, 7), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(13, 4, 1, 14), "bar"),
								ast.NewModifierNode(
									P(19, 18, 1, 20),
									T(P(27, 6, 1, 28), token.UNLESS),
									ast.NewBinaryExpressionNode(
										P(19, 7, 1, 20),
										T(P(23, 1, 1, 24), token.PLUS),
										ast.NewPublicIdentifierNode(P(19, 3, 1, 20), "baz"),
										ast.NewIntLiteralNode(P(25, 1, 1, 26), V(P(25, 1, 1, 26), token.DEC_INT, "5")),
									),
									ast.NewPublicIdentifierNode(P(34, 3, 1, 35), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can contain for modifiers": {
			input: "%{.1, 'foo', :bar, baz + 5 for baz in bazz}",
			want: ast.NewProgramNode(
				P(0, 43, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 43, 1, 1),
						ast.NewSetLiteralNode(
							P(0, 43, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(2, 2, 1, 3), "0.1"),
								ast.NewRawStringLiteralNode(P(6, 5, 1, 7), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(13, 4, 1, 14), "bar"),
								ast.NewModifierForInNode(
									P(19, 23, 1, 20),
									ast.NewBinaryExpressionNode(
										P(19, 7, 1, 20),
										T(P(23, 1, 1, 24), token.PLUS),
										ast.NewPublicIdentifierNode(P(19, 3, 1, 20), "baz"),
										ast.NewIntLiteralNode(P(25, 1, 1, 26), V(P(25, 1, 1, 26), token.DEC_INT, "5")),
									),
									[]ast.ParameterNode{
										ast.NewLoopParameterNode(P(31, 3, 1, 32), "baz", nil),
									},
									ast.NewPublicIdentifierNode(P(38, 4, 1, 39), "bazz"),
								),
							},
						),
					),
				},
			),
		},
		"can have elements": {
			input: "%{.1, 'foo', :bar, baz + 5}",
			want: ast.NewProgramNode(
				P(0, 27, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 27, 1, 1),
						ast.NewSetLiteralNode(
							P(0, 27, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(2, 2, 1, 3), "0.1"),
								ast.NewRawStringLiteralNode(P(6, 5, 1, 7), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(13, 4, 1, 14), "bar"),
								ast.NewBinaryExpressionNode(
									P(19, 7, 1, 20),
									T(P(23, 1, 1, 24), token.PLUS),
									ast.NewPublicIdentifierNode(P(19, 3, 1, 20), "baz"),
									ast.NewIntLiteralNode(P(25, 1, 1, 26), V(P(25, 1, 1, 26), token.DEC_INT, "5")),
								),
							},
						),
					),
				},
			),
		},
		"can't have explicit indices": {
			input: "%{.1, 'foo', 10 => :bar, baz => baz + 5}",
			want: ast.NewProgramNode(
				P(0, 40, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(16, 24, 1, 17),
						ast.NewInvalidNode(
							P(16, 2, 1, 17),
							T(P(16, 2, 1, 17), token.THICK_ARROW),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(16, 2, 1, 17), "unexpected =>, expected }"),
			},
		},
		"can span multiple lines": {
			input: "%{\n.1\n,\n'foo'\n,\n:bar\n,\nbaz + 5\n}",
			want: ast.NewProgramNode(
				P(0, 32, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 32, 1, 1),
						ast.NewSetLiteralNode(
							P(0, 32, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(3, 2, 2, 1), "0.1"),
								ast.NewRawStringLiteralNode(P(8, 5, 4, 1), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(16, 4, 6, 1), "bar"),
								ast.NewBinaryExpressionNode(
									P(23, 7, 8, 1),
									T(P(27, 1, 8, 5), token.PLUS),
									ast.NewPublicIdentifierNode(P(23, 3, 8, 1), "baz"),
									ast.NewIntLiteralNode(P(29, 1, 8, 7), V(P(29, 1, 8, 7), token.DEC_INT, "5")),
								),
							},
						),
					),
				},
			),
		},
		"can be nested": {
			input: "%{%{.1, :+}, .2}",
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewSetLiteralNode(
							P(0, 16, 1, 1),
							[]ast.ExpressionNode{
								ast.NewSetLiteralNode(
									P(2, 9, 1, 3),
									[]ast.ExpressionNode{
										ast.NewFloatLiteralNode(P(4, 2, 1, 5), "0.1"),
										ast.NewSimpleSymbolLiteralNode(P(8, 2, 1, 9), "+"),
									},
								),
								ast.NewFloatLiteralNode(P(13, 2, 1, 14), "0.2"),
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

func TestMapLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "{}",
			want: ast.NewProgramNode(
				P(0, 2, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewMapLiteralNode(
							P(0, 2, 1, 1),
							nil,
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "{\n\n}",
			want: ast.NewProgramNode(
				P(0, 4, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewMapLiteralNode(
							P(0, 4, 1, 1),
							nil,
						),
					),
				},
			),
		},
		"can't contain elements other than key value pairs and identifiers": {
			input: "{.1, 'foo', :bar, baz + 5 if baz}",
			want: ast.NewProgramNode(
				P(0, 33, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 33, 1, 1),
						ast.NewMapLiteralNode(
							P(0, 33, 1, 1),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(P(1, 2, 1, 2), "0.1"),
								ast.NewRawStringLiteralNode(P(5, 5, 1, 6), "foo"),
								ast.NewSimpleSymbolLiteralNode(P(12, 4, 1, 13), "bar"),
								ast.NewModifierNode(
									P(18, 14, 1, 19),
									T(P(26, 2, 1, 27), token.IF),
									ast.NewBinaryExpressionNode(
										P(18, 7, 1, 19),
										T(P(22, 1, 1, 23), token.PLUS),
										ast.NewPublicIdentifierNode(P(18, 3, 1, 19), "baz"),
										ast.NewIntLiteralNode(P(24, 1, 1, 25), V(P(24, 1, 1, 25), token.DEC_INT, "5")),
									),
									ast.NewPublicIdentifierNode(P(29, 3, 1, 30), "baz"),
								),
							},
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(1, 2, 1, 2), "expected a key-value pair, map literals should consist of key-value pairs"),
				NewError(P(5, 5, 1, 6), "expected a key-value pair, map literals should consist of key-value pairs"),
				NewError(P(12, 4, 1, 13), "expected a key-value pair, map literals should consist of key-value pairs"),
				NewError(P(18, 7, 1, 19), "expected a key-value pair, map literals should consist of key-value pairs"),
			},
		},
		"can contain any expression as key with thick arrows": {
			input: "{Math::PI => 3, foo => foo && bar, 5 => 'bar', 'baz' => :bar, a + 5 => 1, n.to_string() => n}",
			want: ast.NewProgramNode(
				P(0, 93, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 93, 1, 1),
						ast.NewMapLiteralNode(
							P(0, 93, 1, 1),
							[]ast.ExpressionNode{
								ast.NewKeyValueExpressionNode(
									P(1, 13, 1, 2),
									ast.NewConstantLookupNode(
										P(1, 8, 1, 2),
										ast.NewPublicConstantNode(P(1, 4, 1, 2), "Math"),
										ast.NewPublicConstantNode(P(7, 2, 1, 8), "PI"),
									),
									ast.NewIntLiteralNode(P(13, 1, 1, 14), V(P(13, 1, 1, 14), token.DEC_INT, "3")),
								),
								ast.NewKeyValueExpressionNode(
									P(16, 17, 1, 17),
									ast.NewPublicIdentifierNode(P(16, 3, 1, 17), "foo"),
									ast.NewLogicalExpressionNode(
										P(23, 10, 1, 24),
										T(P(27, 2, 1, 28), token.AND_AND),
										ast.NewPublicIdentifierNode(P(23, 3, 1, 24), "foo"),
										ast.NewPublicIdentifierNode(P(30, 3, 1, 31), "bar"),
									),
								),
								ast.NewKeyValueExpressionNode(
									P(35, 10, 1, 36),
									ast.NewIntLiteralNode(P(35, 1, 1, 36), V(P(35, 1, 1, 36), token.DEC_INT, "5")),
									ast.NewRawStringLiteralNode(P(40, 5, 1, 41), "bar"),
								),
								ast.NewKeyValueExpressionNode(
									P(47, 13, 1, 48),
									ast.NewRawStringLiteralNode(P(47, 5, 1, 48), "baz"),
									ast.NewSimpleSymbolLiteralNode(P(56, 4, 1, 57), "bar"),
								),
								ast.NewKeyValueExpressionNode(
									P(62, 10, 1, 63),
									ast.NewBinaryExpressionNode(
										P(62, 5, 1, 63),
										T(P(64, 1, 1, 65), token.PLUS),
										ast.NewPublicIdentifierNode(P(62, 1, 1, 63), "a"),
										ast.NewIntLiteralNode(P(66, 1, 1, 67), V(P(66, 1, 1, 67), token.DEC_INT, "5")),
									),
									ast.NewIntLiteralNode(P(71, 1, 1, 72), V(P(71, 1, 1, 72), token.DEC_INT, "1")),
								),
								ast.NewKeyValueExpressionNode(
									P(74, 18, 1, 75),
									ast.NewMethodCallNode(
										P(74, 13, 1, 75),
										ast.NewPublicIdentifierNode(P(74, 1, 1, 75), "n"),
										false,
										"to_string",
										nil,
										nil,
									),
									ast.NewPublicIdentifierNode(P(91, 1, 1, 92), "n"),
								),
							},
						),
					),
				},
			),
		},
		"can have shorthand symbol keys": {
			input: "{foo: :bar}",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewMapLiteralNode(
							P(0, 11, 1, 1),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									P(1, 9, 1, 2),
									"foo",
									ast.NewSimpleSymbolLiteralNode(P(6, 4, 1, 7), "bar"),
								),
							},
						),
					),
				},
			),
		},
		"can contain for modifiers": {
			input: "{foo: bar, baz => baz.to_int for baz in bazz}",
			want: ast.NewProgramNode(
				P(0, 45, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 45, 1, 1),
						ast.NewMapLiteralNode(
							P(0, 45, 1, 1),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									P(1, 8, 1, 2),
									"foo",
									ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
								),
								ast.NewModifierForInNode(
									P(11, 33, 1, 12),
									ast.NewKeyValueExpressionNode(
										P(11, 17, 1, 12),
										ast.NewPublicIdentifierNode(P(11, 3, 1, 12), "baz"),
										ast.NewMethodCallNode(
											P(18, 10, 1, 19),
											ast.NewPublicIdentifierNode(P(18, 3, 1, 19), "baz"),
											false,
											"to_int",
											nil,
											nil,
										),
									),
									[]ast.ParameterNode{
										ast.NewLoopParameterNode(P(33, 3, 1, 34), "baz", nil),
									},
									ast.NewPublicIdentifierNode(P(40, 4, 1, 41), "bazz"),
								),
							},
						),
					),
				},
			),
		},
		"can contain if modifiers": {
			input: "{foo: bar, baz => baz.to_int if baz}",
			want: ast.NewProgramNode(
				P(0, 36, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 36, 1, 1),
						ast.NewMapLiteralNode(
							P(0, 36, 1, 1),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									P(1, 8, 1, 2),
									"foo",
									ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
								),
								ast.NewModifierNode(
									P(11, 24, 1, 12),
									T(P(29, 2, 1, 30), token.IF),
									ast.NewKeyValueExpressionNode(
										P(11, 17, 1, 12),
										ast.NewPublicIdentifierNode(P(11, 3, 1, 12), "baz"),
										ast.NewMethodCallNode(
											P(18, 10, 1, 19),
											ast.NewPublicIdentifierNode(P(18, 3, 1, 19), "baz"),
											false,
											"to_int",
											nil,
											nil,
										),
									),
									ast.NewPublicIdentifierNode(P(32, 3, 1, 33), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: "{\nfoo:\nbar,\nbaz =>\nbaz.to_int if\nbaz\n}",
			want: ast.NewProgramNode(
				P(0, 38, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 38, 1, 1),
						ast.NewMapLiteralNode(
							P(0, 38, 1, 1),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									P(2, 8, 2, 1),
									"foo",
									ast.NewPublicIdentifierNode(P(7, 3, 3, 1), "bar"),
								),
								ast.NewModifierNode(
									P(12, 24, 4, 1),
									T(P(30, 2, 5, 12), token.IF),
									ast.NewKeyValueExpressionNode(
										P(12, 17, 4, 1),
										ast.NewPublicIdentifierNode(P(12, 3, 4, 1), "baz"),
										ast.NewMethodCallNode(
											P(19, 10, 5, 1),
											ast.NewPublicIdentifierNode(P(19, 3, 5, 1), "baz"),
											false,
											"to_int",
											nil,
											nil,
										),
									),
									ast.NewPublicIdentifierNode(P(33, 3, 6, 1), "baz"),
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

func TestRangeLiteral(t *testing.T) {
	tests := testTable{
		"can be beginless and inclusive": {
			input: "..5",
			want: ast.NewProgramNode(
				P(0, 3, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 3, 1, 1),
						ast.NewRangeLiteralNode(
							P(0, 3, 1, 1),
							false,
							nil,
							ast.NewIntLiteralNode(P(2, 1, 1, 3), V(P(2, 1, 1, 3), token.DEC_INT, "5")),
						),
					),
				},
			),
		},
		"can be beginless and exclusive": {
			input: "...5",
			want: ast.NewProgramNode(
				P(0, 4, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewRangeLiteralNode(
							P(0, 4, 1, 1),
							true,
							nil,
							ast.NewIntLiteralNode(P(3, 1, 1, 4), V(P(3, 1, 1, 4), token.DEC_INT, "5")),
						),
					),
				},
			),
		},
		"can be endless and inclusive": {
			input: "5..",
			want: ast.NewProgramNode(
				P(0, 3, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 3, 1, 1),
						ast.NewRangeLiteralNode(
							P(0, 3, 1, 1),
							false,
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.DEC_INT, "5")),
							nil,
						),
					),
				},
			),
		},
		"can be endless and exclusive": {
			input: "5...",
			want: ast.NewProgramNode(
				P(0, 4, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewRangeLiteralNode(
							P(0, 4, 1, 1),
							true,
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.DEC_INT, "5")),
							nil,
						),
					),
				},
			),
		},
		"can have a beginning and be inclusive": {
			input: "2..5",
			want: ast.NewProgramNode(
				P(0, 4, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewRangeLiteralNode(
							P(0, 4, 1, 1),
							false,
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.DEC_INT, "2")),
							ast.NewIntLiteralNode(P(3, 1, 1, 4), V(P(3, 1, 1, 4), token.DEC_INT, "5")),
						),
					),
				},
			),
		},
		"can have a beginning and be exclusive": {
			input: "2...5",
			want: ast.NewProgramNode(
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewRangeLiteralNode(
							P(0, 5, 1, 1),
							true,
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.DEC_INT, "2")),
							ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.DEC_INT, "5")),
						),
					),
				},
			),
		},
		"has higher precedence than method calls": {
			input: "2...5.to_string",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewMethodCallNode(
							P(0, 15, 1, 1),
							ast.NewRangeLiteralNode(
								P(0, 5, 1, 1),
								true,
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), token.DEC_INT, "2")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), token.DEC_INT, "5")),
							),
							false,
							"to_string",
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have any expressions as operands": {
			input: "(2 * 5)...'foo'",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(1, 14, 1, 2),
						ast.NewRangeLiteralNode(
							P(1, 14, 1, 2),
							true,
							ast.NewBinaryExpressionNode(
								P(1, 5, 1, 2),
								T(P(3, 1, 1, 4), token.STAR),
								ast.NewIntLiteralNode(P(1, 1, 1, 2), V(P(1, 1, 1, 2), token.DEC_INT, "2")),
								ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), token.DEC_INT, "5")),
							),
							ast.NewRawStringLiteralNode(
								P(10, 5, 1, 11),
								"foo",
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

func TestTypeLiteral(t *testing.T) {
	tests := testTable{
		"can have a constant as a type": {
			input: "type String",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewTypeLiteralNode(
							P(0, 11, 1, 1),
							ast.NewPublicConstantNode(P(5, 6, 1, 6), "String"),
						),
					),
				},
			),
		},
		"can have a nilable type": {
			input: "type String?",
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 12, 1, 1),
						ast.NewTypeLiteralNode(
							P(0, 12, 1, 1),
							ast.NewNilableTypeNode(
								P(5, 7, 1, 6),
								ast.NewPublicConstantNode(P(5, 6, 1, 6), "String"),
							),
						),
					),
				},
			),
		},
		"can have a union type": {
			input: "type Int | String",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewTypeLiteralNode(
							P(0, 17, 1, 1),
							ast.NewBinaryTypeExpressionNode(
								P(5, 12, 1, 6),
								T(P(9, 1, 1, 10), token.OR),
								ast.NewPublicConstantNode(P(5, 3, 1, 6), "Int"),
								ast.NewPublicConstantNode(P(11, 6, 1, 12), "String"),
							),
						),
					),
				},
			),
		},
		"can have an intersection type": {
			input: "type Int & String",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewTypeLiteralNode(
							P(0, 17, 1, 1),
							ast.NewBinaryTypeExpressionNode(
								P(5, 12, 1, 6),
								T(P(9, 1, 1, 10), token.AND),
								ast.NewPublicConstantNode(P(5, 3, 1, 6), "Int"),
								ast.NewPublicConstantNode(P(11, 6, 1, 12), "String"),
							),
						),
					),
				},
			),
		},
		"can have a generic type": {
			input: "type Std::Map[Std::Symbol, List[String]]",
			want: ast.NewProgramNode(
				P(0, 40, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 40, 1, 1),
						ast.NewTypeLiteralNode(
							P(0, 40, 1, 1),
							ast.NewGenericConstantNode(
								P(5, 35, 1, 6),
								ast.NewConstantLookupNode(
									P(5, 8, 1, 6),
									ast.NewPublicConstantNode(P(5, 3, 1, 6), "Std"),
									ast.NewPublicConstantNode(P(10, 3, 1, 11), "Map"),
								),
								[]ast.ComplexConstantNode{
									ast.NewConstantLookupNode(
										P(14, 11, 1, 15),
										ast.NewPublicConstantNode(P(14, 3, 1, 15), "Std"),
										ast.NewPublicConstantNode(P(19, 6, 1, 20), "Symbol"),
									),
									ast.NewGenericConstantNode(
										P(27, 12, 1, 28),
										ast.NewPublicConstantNode(P(27, 4, 1, 28), "List"),
										[]ast.ComplexConstantNode{
											ast.NewPublicConstantNode(P(32, 6, 1, 33), "String"),
										},
									),
								},
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
