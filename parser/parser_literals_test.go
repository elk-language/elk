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
						ast.NewStringLiteralNode(
							P(0, 36, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 34, 1, 2), "foo\nbar\rbaz\\car\t\b\"\v\f\x12\a"),
							},
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
						ast.NewStringLiteralNode(
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
						ast.NewStringLiteralNode(
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
						ast.NewStringLiteralNode(
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
						ast.NewStringLiteralNode(
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
						ast.NewStringLiteralNode(
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
				NewError(P(6, 1, 1, 7), "unexpected PublicIdentifier, expected a statement separator `\\n`, `;`"),
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
		"can omit arguments and be single line": {
			input: `-> 'foo' + .2`,
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 13, 1, 1),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(3, 10, 1, 4),
									ast.NewBinaryExpressionNode(
										P(3, 10, 1, 4),
										T(P(9, 1, 1, 10), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(3, 5, 1, 4),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(11, 2, 1, 12),
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
		"can omit arguments and be single line with braces": {
			input: `-> { 'foo' + .2 }`,
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 17, 1, 1),
							nil,
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
		"can omit arguments and be multiline with braces": {
			input: `-> {
	'foo' + .2
	nil
}`,
			want: ast.NewProgramNode(
				P(0, 23, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 23, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 23, 1, 1),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(6, 11, 2, 2),
									ast.NewBinaryExpressionNode(
										P(6, 10, 2, 2),
										T(P(12, 1, 2, 8), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(6, 5, 2, 2),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(14, 2, 2, 10),
											"0.2",
										),
									),
								),
								ast.NewExpressionStatementNode(
									P(18, 4, 3, 2),
									ast.NewNilLiteralNode(P(18, 3, 3, 2)),
								),
							},
						),
					),
				},
			),
		},
		"can omit arguments and be multiline with end": {
			input: `->
	'foo' + .2
	nil
end`,
			want: ast.NewProgramNode(
				P(0, 23, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 23, 1, 1),
						ast.NewClosureExpressionNode(
							P(0, 23, 1, 1),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(4, 11, 2, 2),
									ast.NewBinaryExpressionNode(
										P(4, 10, 2, 2),
										T(P(10, 1, 2, 8), token.PLUS),
										ast.NewRawStringLiteralNode(
											P(4, 5, 2, 2),
											"foo",
										),
										ast.NewFloatLiteralNode(
											P(12, 2, 2, 10),
											"0.2",
										),
									),
								),
								ast.NewExpressionStatementNode(
									P(16, 4, 3, 2),
									ast.NewNilLiteralNode(P(16, 3, 3, 2)),
								),
							},
						),
					),
				},
			),
		},
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
									token.NewWithValue(P(1, 1, 1, 2), token.PUBLIC_IDENTIFIER, "a"),
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
									token.NewWithValue(P(1, 1, 1, 2), token.PUBLIC_IDENTIFIER, "a"),
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
									token.NewWithValue(P(1, 1, 1, 2), token.PUBLIC_IDENTIFIER, "a"),
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
									token.NewWithValue(P(1, 1, 1, 2), token.PUBLIC_IDENTIFIER, "a"),
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
									token.NewWithValue(P(0, 1, 1, 1), token.PUBLIC_IDENTIFIER, "a"),
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
									token.NewWithValue(P(1, 1, 1, 2), token.PUBLIC_IDENTIFIER, "a"),
									ast.NewPublicConstantNode(P(4, 3, 1, 5), "Int"),
									nil,
								),
								ast.NewFormalParameterNode(
									P(9, 9, 1, 10),
									token.NewWithValue(P(9, 1, 1, 10), token.PUBLIC_IDENTIFIER, "b"),
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
									token.NewWithValue(P(1, 1, 1, 2), token.PUBLIC_IDENTIFIER, "a"),
									nil,
									ast.NewIntLiteralNode(P(5, 2, 1, 6), V(P(5, 2, 1, 6), token.DEC_INT, "32")),
								),
								ast.NewFormalParameterNode(
									P(9, 17, 1, 10),
									token.NewWithValue(P(9, 1, 1, 10), token.PUBLIC_IDENTIFIER, "b"),
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
