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
				NewError(P(6, 1, 1, 7), "unexpected PublicIdentifier, expected a statement separator `\\n`, `;` or end of file"),
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
