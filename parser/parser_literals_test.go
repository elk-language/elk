package parser

import (
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/token"
)

func TestFloatLiteral(t *testing.T) {
	tests := testTable{
		"can have underscores": {
			input: `245_000.254_129`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewFloatLiteralNode(S(P(0, 1, 1), P(14, 1, 15)), `245000.254129`),
					),
				},
			),
		},
		"ends on the last valid character": {
			input: `0.36p`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewFloatLiteralNode(S(P(0, 1, 1), P(3, 1, 4)), `0.36`),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(4, 1, 5), P(4, 1, 5)), "unexpected PUBLIC_IDENTIFIER, expected a statement separator `\\n`, `;`"),
			},
		},
		"can only be decimal": {
			input: `0x21.36`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(3, 1, 4)), `0x21`),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(4, 1, 5), P(6, 1, 7)), "unexpected FLOAT, expected a statement separator `\\n`, `;`"),
			},
		},
		"can have an exponent": {
			input: `0.36e2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewFloatLiteralNode(S(P(0, 1, 1), P(5, 1, 6)), `0.36e2`),
					),
				},
			),
		},
		"with exponent and no dot": {
			input: `25e4`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewFloatLiteralNode(S(P(0, 1, 1), P(3, 1, 4)), `25e4`),
					),
				},
			),
		},
		"with an uppercase exponent": {
			input: `25E4`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewFloatLiteralNode(S(P(0, 1, 1), P(3, 1, 4)), `25e4`),
					),
				},
			),
		},
		"with an explicit positive exponent": {
			input: `25E+4`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewFloatLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `25e+4`),
					),
				},
			),
		},
		"with a negative exponent": {
			input: `25E-4`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewFloatLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `25e-4`),
					),
				},
			),
		},
		"without a leading zero": {
			input: `.908267374623`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewFloatLiteralNode(S(P(0, 1, 1), P(12, 1, 13)), `0.908267374623`),
					),
				},
			),
		},
		"BigFloat without a dot": {
			input: `24bf`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewBigFloatLiteralNode(S(P(0, 1, 1), P(3, 1, 4)), `24`),
					),
				},
			),
		},
		"BigFloat with a dot": {
			input: `24.5bf`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewBigFloatLiteralNode(S(P(0, 1, 1), P(5, 1, 6)), `24.5`),
					),
				},
			),
		},
		"BigFloat with an exponent": {
			input: `24e5_bf`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewBigFloatLiteralNode(S(P(0, 1, 1), P(6, 1, 7)), `24e5`),
					),
				},
			),
		},
		"BigFloat with an exponent and dot": {
			input: `24.5e5_bf`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewBigFloatLiteralNode(S(P(0, 1, 1), P(8, 1, 9)), `24.5e5`),
					),
				},
			),
		},
		"float64 without a dot": {
			input: `24f64`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewFloat64LiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `24`),
					),
				},
			),
		},
		"float64 with a dot": {
			input: `24.5f64`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewFloat64LiteralNode(S(P(0, 1, 1), P(6, 1, 7)), `24.5`),
					),
				},
			),
		},
		"float64 with an exponent": {
			input: `24e5f64`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewFloat64LiteralNode(S(P(0, 1, 1), P(6, 1, 7)), `24e5`),
					),
				},
			),
		},
		"float64 with an exponent and dot": {
			input: `24.5e5f64`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewFloat64LiteralNode(S(P(0, 1, 1), P(8, 1, 9)), `24.5e5`),
					),
				},
			),
		},
		"float32 without a dot": {
			input: `24f32`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewFloat32LiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `24`),
					),
				},
			),
		},
		"float32 with a dot": {
			input: `24.5f32`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewFloat32LiteralNode(S(P(0, 1, 1), P(6, 1, 7)), `24.5`),
					),
				},
			),
		},
		"float32 with an exponent": {
			input: `24e5f32`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewFloat32LiteralNode(S(P(0, 1, 1), P(6, 1, 7)), `24e5`),
					),
				},
			),
		},
		"float32 with an exponent and dot": {
			input: `24.5e5f32`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewFloat32LiteralNode(S(P(0, 1, 1), P(8, 1, 9)), `24.5e5`),
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

func TestIntLiteral(t *testing.T) {
	tests := testTable{
		"decimal": {
			input: `23`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(1, 1, 2)), `23`),
					),
				},
			),
		},
		"decimal int64": {
			input: `23i64`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewInt64LiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `23`),
					),
				},
			),
		},
		"decimal uint64": {
			input: `23u64`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewUInt64LiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `23`),
					),
				},
			),
		},
		"decimal int32": {
			input: `23i32`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewInt32LiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `23`),
					),
				},
			),
		},
		"decimal uint32": {
			input: `23u32`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewUInt32LiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `23`),
					),
				},
			),
		},
		"decimal int16": {
			input: `23i16`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewInt16LiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `23`),
					),
				},
			),
		},
		"decimal uint16": {
			input: `23u16`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewUInt16LiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `23`),
					),
				},
			),
		},
		"decimal int8": {
			input: `23i8`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewInt8LiteralNode(S(P(0, 1, 1), P(3, 1, 4)), `23`),
					),
				},
			),
		},
		"decimal uint8": {
			input: `23u8`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewUInt8LiteralNode(S(P(0, 1, 1), P(3, 1, 4)), `23`),
					),
				},
			),
		},
		"decimal with leading zeros": {
			input: `00015`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `00015`),
					),
				},
			),
		},
		"decimal with underscores": {
			input: `23_200_123`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(9, 1, 10)), `23200123`),
					),
				},
			),
		},
		"hex": {
			input: `0xff24`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(5, 1, 6)), `0xff24`),
					),
				},
			),
		},
		"duodecimal": {
			input: `0d2a4`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `0d2a4`),
					),
				},
			),
		},
		"octal": {
			input: `0o723`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `0o723`),
					),
				},
			),
		},
		"quaternary": {
			input: `0q323`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), `0q323`),
					),
				},
			),
		},
		"binary": {
			input: `0b1101`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewIntLiteralNode(S(P(0, 1, 1), P(5, 1, 6)), `0b1101`),
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

func TestStringLiteral(t *testing.T) {
	tests := testTable{
		"processes escape sequences": {
			input: `"foo\nbar\rbaz\\car\t\b\"\v\f\x12\a\u00e9\U0010FFFF"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(51, 1, 52)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(51, 1, 52)),
						ast.NewDoubleQuotedStringLiteralNode(
							S(P(0, 1, 1), P(51, 1, 52)),
							"foo\nbar\rbaz\\car\t\b\"\v\f\x12\a\u00e9\U0010FFFF",
						),
					),
				},
			),
		},
		"can be empty": {
			input: `""`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewDoubleQuotedStringLiteralNode(
							S(P(0, 1, 1), P(1, 1, 2)),
							"",
						),
					),
				},
			),
		},
		"reports errors for invalid hex escapes": {
			input: `"foo \xgh bar"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewInterpolatedStringLiteralNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(S(P(1, 1, 2), P(4, 1, 5)), "foo "),
								ast.NewInvalidNode(S(P(5, 1, 6), P(8, 1, 9)), V(S(P(5, 1, 6), P(8, 1, 9)), token.ERROR, "invalid hex escape")),
								ast.NewStringLiteralContentSectionNode(S(P(9, 1, 10), P(12, 1, 13)), " bar"),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(5, 1, 6), P(8, 1, 9)), "invalid hex escape"),
			},
		},
		"reports errors for invalid unicode escapes": {
			input: `"foo \u7fgf bar"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewInterpolatedStringLiteralNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(S(P(1, 1, 2), P(4, 1, 5)), "foo "),
								ast.NewInvalidNode(S(P(5, 1, 6), P(10, 1, 11)), V(S(P(5, 1, 6), P(10, 1, 11)), token.ERROR, "invalid unicode escape")),
								ast.NewStringLiteralContentSectionNode(S(P(11, 1, 12), P(14, 1, 15)), " bar"),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(5, 1, 6), P(10, 1, 11)), "invalid unicode escape"),
			},
		},
		"reports errors for invalid big unicode escapes": {
			input: `"foo \U7fgf0234 bar"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewInterpolatedStringLiteralNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(S(P(1, 1, 2), P(4, 1, 5)), "foo "),
								ast.NewInvalidNode(S(P(5, 1, 6), P(14, 1, 15)), V(S(P(5, 1, 6), P(14, 1, 15)), token.ERROR, "invalid unicode escape")),
								ast.NewStringLiteralContentSectionNode(S(P(15, 1, 16), P(18, 1, 19)), " bar"),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(5, 1, 6), P(14, 1, 15)), "invalid unicode escape"),
			},
		},
		"reports errors for nonexistent escape sequences": {
			input: `"foo \q bar"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewInterpolatedStringLiteralNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(S(P(1, 1, 2), P(4, 1, 5)), "foo "),
								ast.NewInvalidNode(S(P(5, 1, 6), P(6, 1, 7)), V(S(P(5, 1, 6), P(6, 1, 7)), token.ERROR, "invalid escape sequence `\\q` in string literal")),
								ast.NewStringLiteralContentSectionNode(S(P(7, 1, 8), P(10, 1, 11)), " bar"),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(5, 1, 6), P(6, 1, 7)), "invalid escape sequence `\\q` in string literal"),
			},
		},
		"can contain interpolated expressions": {
			input: `"foo ${bar + 2} baz ${fudge}"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 1, 29)),
						ast.NewInterpolatedStringLiteralNode(
							S(P(0, 1, 1), P(28, 1, 29)),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(S(P(1, 1, 2), P(4, 1, 5)), "foo "),
								ast.NewStringInterpolationNode(
									S(P(5, 1, 6), P(14, 1, 15)),
									ast.NewBinaryExpressionNode(
										S(P(7, 1, 8), P(13, 1, 14)),
										T(S(P(11, 1, 12), P(11, 1, 12)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
										ast.NewIntLiteralNode(S(P(13, 1, 14), P(13, 1, 14)), "2"),
									),
								),
								ast.NewStringLiteralContentSectionNode(S(P(15, 1, 16), P(19, 1, 20)), " baz "),
								ast.NewStringInterpolationNode(
									S(P(20, 1, 21), P(27, 1, 28)),
									ast.NewPublicIdentifierNode(S(P(22, 1, 23), P(26, 1, 27)), "fudge"),
								),
							},
						),
					),
				},
			),
		},
		"cannot contain string literals inside interpolation": {
			input: `"foo ${"bar" + 2} baza"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewInterpolatedStringLiteralNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(S(P(1, 1, 2), P(4, 1, 5)), "foo "),
								ast.NewStringInterpolationNode(
									S(P(5, 1, 6), P(16, 1, 17)),
									ast.NewBinaryExpressionNode(
										S(P(7, 1, 8), P(15, 1, 16)),
										T(S(P(13, 1, 14), P(13, 1, 14)), token.PLUS),
										ast.NewInvalidNode(S(P(7, 1, 8), P(11, 1, 12)), V(S(P(7, 1, 8), P(11, 1, 12)), token.ERROR, "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation")),
										ast.NewIntLiteralNode(S(P(15, 1, 16), P(15, 1, 16)), "2"),
									),
								),
								ast.NewStringLiteralContentSectionNode(S(P(17, 1, 18), P(21, 1, 22)), " baza"),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(7, 1, 8), P(11, 1, 12)), "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation"),
			},
		},
		"can contain raw string literals inside interpolation": {
			input: `"foo ${'bar' + 2} baza"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewInterpolatedStringLiteralNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(S(P(1, 1, 2), P(4, 1, 5)), "foo "),
								ast.NewStringInterpolationNode(
									S(P(5, 1, 6), P(16, 1, 17)),
									ast.NewBinaryExpressionNode(
										S(P(7, 1, 8), P(15, 1, 16)),
										T(S(P(13, 1, 14), P(13, 1, 14)), token.PLUS),
										ast.NewRawStringLiteralNode(S(P(7, 1, 8), P(11, 1, 12)), "bar"),
										ast.NewIntLiteralNode(S(P(15, 1, 16), P(15, 1, 16)), "2"),
									),
								),
								ast.NewStringLiteralContentSectionNode(S(P(17, 1, 18), P(21, 1, 22)), " baza"),
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
				S(P(0, 1, 1), P(35, 1, 36)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(35, 1, 36)),
						ast.NewRawStringLiteralNode(S(P(0, 1, 1), P(35, 1, 36)), `foo\nbar\rbaz\\car\t\b\"\v\f\x12\a`),
					),
				},
			),
		},
		"cannot contain interpolated expressions": {
			input: `'foo ${bar + 2} baz ${fudge}'`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 1, 29)),
						ast.NewRawStringLiteralNode(S(P(0, 1, 1), P(28, 1, 29)), `foo ${bar + 2} baz ${fudge}`),
					),
				},
			),
		},
		"can contain double quotes": {
			input: `'foo ${"bar" + 2} baza'`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewRawStringLiteralNode(S(P(0, 1, 1), P(22, 1, 23)), `foo ${"bar" + 2} baza`),
					),
				},
			),
		},
		"doesn't allow escaping single quotes": {
			input: `'foo\'s house'`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewRawStringLiteralNode(S(P(0, 1, 1), P(5, 1, 6)), "foo\\"),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(6, 1, 7), P(6, 1, 7)), "unexpected PUBLIC_IDENTIFIER, expected a statement separator `\\n`, `;`"),
				errors.NewError(L("<main>", P(13, 1, 14), P(13, 1, 14)), "unterminated raw string literal, missing `'`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestCharLiteral(t *testing.T) {
	tests := testTable{
		"must be terminated": {
			input: "`a",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewInvalidNode(S(P(0, 1, 1), P(1, 1, 2)), V(S(P(0, 1, 1), P(1, 1, 2)), token.ERROR, "unterminated character literal, missing backtick")),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(0, 1, 1), P(1, 1, 2)), "unterminated character literal, missing backtick"),
			},
		},
		"can contain ascii characters": {
			input: "`a`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						ast.NewCharLiteralNode(S(P(0, 1, 1), P(2, 1, 3)), 'a'),
					),
				},
			),
		},
		"can contain utf8 characters": {
			input: "`ś`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 3)),
						ast.NewCharLiteralNode(S(P(0, 1, 1), P(3, 1, 3)), 'ś'),
					),
				},
			),
		},
		"escapes backticks": {
			input: "`\\``",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewCharLiteralNode(S(P(0, 1, 1), P(3, 1, 4)), '`'),
					),
				},
			),
		},
		"cannot contain multiple characters": {
			input: "`lalala`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewInvalidNode(S(P(0, 1, 1), P(7, 1, 8)), V(S(P(0, 1, 1), P(7, 1, 8)), token.ERROR, "invalid char literal with more than one character")),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(0, 1, 1), P(7, 1, 8)), "invalid char literal with more than one character"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestRawCharLiteral(t *testing.T) {
	tests := testTable{
		"must be terminated": {
			input: "r`a",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						ast.NewInvalidNode(S(P(0, 1, 1), P(2, 1, 3)), V(S(P(0, 1, 1), P(2, 1, 3)), token.ERROR, "unterminated character literal, missing backtick")),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(0, 1, 1), P(2, 1, 3)), "unterminated character literal, missing backtick"),
			},
		},
		"can contain ascii characters": {
			input: "r`a`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewRawCharLiteralNode(S(P(0, 1, 1), P(3, 1, 4)), 'a'),
					),
				},
			),
		},
		"can contain utf8 characters": {
			input: "r`ś`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 4)),
						ast.NewRawCharLiteralNode(S(P(0, 1, 1), P(4, 1, 4)), 'ś'),
					),
				},
			),
		},
		"cannot escape backticks": {
			input: "r`\\``",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewRawCharLiteralNode(S(P(0, 1, 1), P(3, 1, 4)), '\\'),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(4, 1, 5), P(4, 1, 5)), "unterminated character literal, missing backtick"),
			},
		},
		"cannot contain multiple characters": {
			input: "r`lalala`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(8, 1, 9)),
						ast.NewInvalidNode(S(P(0, 1, 1), P(8, 1, 9)), V(S(P(0, 1, 1), P(8, 1, 9)), token.ERROR, "invalid raw char literal with more than one character")),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(0, 1, 1), P(8, 1, 9)), "invalid raw char literal with more than one character"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestFunction(t *testing.T) {
	tests := testTable{
		"can have arguments and be single line": {
			input: `|a| -> 'foo' + .2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(7, 1, 8), P(16, 1, 17)),
									ast.NewBinaryExpressionNode(
										S(P(7, 1, 8), P(16, 1, 17)),
										T(S(P(13, 1, 14), P(13, 1, 14)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(7, 1, 8), P(11, 1, 12)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(15, 1, 16), P(16, 1, 17)),
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
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(9, 1, 10), P(18, 1, 19)),
									ast.NewBinaryExpressionNode(
										S(P(9, 1, 10), P(18, 1, 19)),
										T(S(P(15, 1, 16), P(15, 1, 16)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(9, 1, 10), P(13, 1, 14)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(17, 1, 18), P(18, 1, 19)),
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
				S(P(0, 1, 1), P(26, 4, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 4, 1)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(26, 4, 1)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(10, 2, 2), P(20, 2, 12)),
									ast.NewBinaryExpressionNode(
										S(P(10, 2, 2), P(19, 2, 11)),
										T(S(P(16, 2, 8), P(16, 2, 8)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(10, 2, 2), P(14, 2, 6)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(18, 2, 10), P(19, 2, 11)),
											"0.2",
										),
									),
								),
								ast.NewExpressionStatementNode(
									S(P(22, 3, 2), P(25, 3, 5)),
									ast.NewNilLiteralNode(S(P(22, 3, 2), P(24, 3, 4))),
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
				S(P(0, 1, 1), P(26, 4, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 4, 3)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(26, 4, 3)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(8, 2, 2), P(18, 2, 12)),
									ast.NewBinaryExpressionNode(
										S(P(8, 2, 2), P(17, 2, 11)),
										T(S(P(14, 2, 8), P(14, 2, 8)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(8, 2, 2), P(12, 2, 6)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(16, 2, 10), P(17, 2, 11)),
											"0.2",
										),
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
		"can omit pipes when there's a single argument": {
			input: `a -> 'foo' + .2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(0, 1, 1), P(0, 1, 1)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(5, 1, 6), P(14, 1, 15)),
									ast.NewBinaryExpressionNode(
										S(P(5, 1, 6), P(14, 1, 15)),
										T(S(P(11, 1, 12), P(11, 1, 12)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(5, 1, 6), P(9, 1, 10)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(13, 1, 14), P(14, 1, 15)),
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
		"can have a positional rest argument": {
			input: "|a, b, *c| -> nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"b",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(7, 1, 8), P(8, 1, 9)),
									"c",
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(14, 1, 15), P(16, 1, 17)),
									ast.NewNilLiteralNode(S(P(14, 1, 15), P(16, 1, 17))),
								),
							},
						),
					),
				},
			),
		},
		"can have a positional rest argument in the middle": {
			input: "|a, b, *c, d| -> nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"b",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(7, 1, 8), P(8, 1, 9)),
									"c",
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(11, 1, 12), P(11, 1, 12)),
									"d",
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(17, 1, 18), P(19, 1, 20)),
									ast.NewNilLiteralNode(S(P(17, 1, 18), P(19, 1, 20))),
								),
							},
						),
					),
				},
			),
		},
		"cannot have multiple positional rest arguments": {
			input: "|a, b, *c, *d| -> nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"b",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(7, 1, 8), P(8, 1, 9)),
									"c",
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(11, 1, 12), P(12, 1, 13)),
									"d",
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(18, 1, 19), P(20, 1, 21)),
									ast.NewNilLiteralNode(S(P(18, 1, 19), P(20, 1, 21))),
								),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(11, 1, 12), P(12, 1, 13)), "there should be only a single positional rest parameter"),
			},
		},
		"can have a positional rest argument with a type": {
			input: "|a, b, *c: String| -> nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"b",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(7, 1, 8), P(16, 1, 17)),
									"c",
									ast.NewPublicConstantNode(S(P(11, 1, 12), P(16, 1, 17)), "String"),
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(22, 1, 23), P(24, 1, 25)),
									ast.NewNilLiteralNode(S(P(22, 1, 23), P(24, 1, 25))),
								),
							},
						),
					),
				},
			),
		},
		"can have a named rest argument": {
			input: "|a, b, **c| -> nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"b",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(7, 1, 8), P(9, 1, 10)),
									"c",
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(15, 1, 16), P(17, 1, 18)),
									ast.NewNilLiteralNode(S(P(15, 1, 16), P(17, 1, 18))),
								),
							},
						),
					),
				},
			),
		},
		"can have a named rest argument with a type": {
			input: "|a, b, **c: String| -> nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"b",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(7, 1, 8), P(17, 1, 18)),
									"c",
									ast.NewPublicConstantNode(S(P(12, 1, 13), P(17, 1, 18)), "String"),
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(23, 1, 24), P(25, 1, 26)),
									ast.NewNilLiteralNode(S(P(23, 1, 24), P(25, 1, 26))),
								),
							},
						),
					),
				},
			),
		},
		"cannot have parameters after a named rest argument": {
			input: "|a, b, **c, d| -> nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"b",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(7, 1, 8), P(9, 1, 10)),
									"c",
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(12, 1, 13), P(12, 1, 13)),
									"d",
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(18, 1, 19), P(20, 1, 21)),
									ast.NewNilLiteralNode(S(P(18, 1, 19), P(20, 1, 21))),
								),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(12, 1, 13), P(12, 1, 13)), "named rest parameters should appear last"),
			},
		},
		"can have a positional and named rest parameter": {
			input: "|a, b, *c, **d| -> nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									"a",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									"b",
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(7, 1, 8), P(8, 1, 9)),
									"c",
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(11, 1, 12), P(13, 1, 14)),
									"d",
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(19, 1, 20), P(21, 1, 22)),
									ast.NewNilLiteralNode(S(P(19, 1, 20), P(21, 1, 22))),
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
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(6, 1, 7)),
									"a",
									ast.NewPublicConstantNode(S(P(4, 1, 5), P(6, 1, 7)), "Int"),
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(9, 1, 10), P(17, 1, 18)),
									"b",
									ast.NewPublicConstantNode(S(P(12, 1, 13), P(17, 1, 18)), "String"),
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(23, 1, 24), P(32, 1, 33)),
									ast.NewBinaryExpressionNode(
										S(P(23, 1, 24), P(32, 1, 33)),
										T(S(P(29, 1, 30), P(29, 1, 30)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(23, 1, 24), P(27, 1, 28)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(31, 1, 32), P(32, 1, 33)),
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
				S(P(0, 1, 1), P(40, 1, 41)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(40, 1, 41)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(40, 1, 41)),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									S(P(1, 1, 2), P(6, 1, 7)),
									"a",
									nil,
									ast.NewIntLiteralNode(S(P(5, 1, 6), P(6, 1, 7)), "32"),
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									S(P(9, 1, 10), P(25, 1, 26)),
									"b",
									ast.NewPublicConstantNode(S(P(12, 1, 13), P(17, 1, 18)), "String"),
									ast.NewRawStringLiteralNode(S(P(21, 1, 22), P(25, 1, 26)), "foo"),
									ast.NormalParameterKind,
								),
							},
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(31, 1, 32), P(40, 1, 41)),
									ast.NewBinaryExpressionNode(
										S(P(31, 1, 32), P(40, 1, 41)),
										T(S(P(37, 1, 38), P(37, 1, 38)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(31, 1, 32), P(35, 1, 36)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(39, 1, 40), P(40, 1, 41)),
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
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(6, 1, 7), P(15, 1, 16)),
									ast.NewBinaryExpressionNode(
										S(P(6, 1, 7), P(15, 1, 16)),
										T(S(P(12, 1, 13), P(12, 1, 13)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(6, 1, 7), P(10, 1, 11)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(14, 1, 15), P(15, 1, 16)),
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
		"can have omit the argument list": {
			input: `-> 'foo' + .2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(3, 1, 4), P(12, 1, 13)),
									ast.NewBinaryExpressionNode(
										S(P(3, 1, 4), P(12, 1, 13)),
										T(S(P(9, 1, 10), P(9, 1, 10)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(3, 1, 4), P(7, 1, 8)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(11, 1, 12), P(12, 1, 13)),
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
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							nil,
							ast.NewNilableTypeNode(S(P(4, 1, 5), P(10, 1, 11)), ast.NewPublicConstantNode(S(P(4, 1, 5), P(9, 1, 10)), "String")),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(15, 1, 16), P(24, 1, 25)),
									ast.NewBinaryExpressionNode(
										S(P(15, 1, 16), P(24, 1, 25)),
										T(S(P(21, 1, 22), P(21, 1, 22)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(15, 1, 16), P(19, 1, 20)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(23, 1, 24), P(24, 1, 25)),
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
		"can have a throw type": {
			input: `||! RuntimeError -> 'foo' + .2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(29, 1, 30)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(29, 1, 30)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(29, 1, 30)),
							nil,
							nil,
							ast.NewPublicConstantNode(S(P(4, 1, 5), P(15, 1, 16)), "RuntimeError"),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(20, 1, 21), P(29, 1, 30)),
									ast.NewBinaryExpressionNode(
										S(P(20, 1, 21), P(29, 1, 30)),
										T(S(P(26, 1, 27), P(26, 1, 27)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(20, 1, 21), P(24, 1, 25)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(28, 1, 29), P(29, 1, 30)),
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
		"can have a return and throw type": {
			input: `||: String? ! RuntimeError -> 'foo' + .2`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(39, 1, 40)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(39, 1, 40)),
						ast.NewFunctionLiteralNode(
							S(P(0, 1, 1), P(39, 1, 40)),
							nil,
							ast.NewNilableTypeNode(S(P(4, 1, 5), P(10, 1, 11)), ast.NewPublicConstantNode(S(P(4, 1, 5), P(9, 1, 10)), "String")),
							ast.NewPublicConstantNode(S(P(14, 1, 15), P(25, 1, 26)), "RuntimeError"),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(30, 1, 31), P(39, 1, 40)),
									ast.NewBinaryExpressionNode(
										S(P(30, 1, 31), P(39, 1, 40)),
										T(S(P(36, 1, 37), P(36, 1, 37)), token.PLUS),
										ast.NewRawStringLiteralNode(
											S(P(30, 1, 31), P(34, 1, 35)),
											"foo",
										),
										ast.NewFloatLiteralNode(
											S(P(38, 1, 39), P(39, 1, 40)),
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
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewSimpleSymbolLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), "foo"),
					),
				},
			),
		},
		"can have a public identifier as the content": {
			input: ":foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewSimpleSymbolLiteralNode(S(P(0, 1, 1), P(3, 1, 4)), "foo"),
					),
				},
			),
		},
		"can have a private identifier as the content": {
			input: ":_foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewSimpleSymbolLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), "_foo"),
					),
				},
			),
		},
		"can have a public constant as the content": {
			input: ":Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewSimpleSymbolLiteralNode(S(P(0, 1, 1), P(3, 1, 4)), "Foo"),
					),
				},
			),
		},
		"can have a private constant as the content": {
			input: ":_Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewSimpleSymbolLiteralNode(S(P(0, 1, 1), P(4, 1, 5)), "_Foo"),
					),
				},
			),
		},
		"can have a keyword as the content": {
			input: ":var",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewSimpleSymbolLiteralNode(S(P(0, 1, 1), P(3, 1, 4)), "var"),
					),
				},
			),
		},
		"can have a raw string as the content": {
			input: ":'foo bar'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewSimpleSymbolLiteralNode(S(P(0, 1, 1), P(9, 1, 10)), "foo bar"),
					),
				},
			),
		},
		"can have a double quoted string as the content": {
			input: `:"foo bar"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewSimpleSymbolLiteralNode(S(P(0, 1, 1), P(9, 1, 10)), "foo bar"),
					),
				},
			),
		},
		"can have an overridable operator as the content": {
			input: ":+",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewSimpleSymbolLiteralNode(S(P(0, 1, 1), P(1, 1, 2)), "+"),
					),
				},
			),
		},
		"cannot have a not overridable operator as the content": {
			input: ":&&",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						ast.NewInvalidNode(S(P(0, 1, 1), P(2, 1, 3)), T(S(P(1, 1, 2), P(2, 1, 3)), token.AND_AND)),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(1, 1, 2), P(2, 1, 3)), "unexpected &&, expected an identifier, overridable operator or string literal"),
			},
		},
		"can have a string as the content": {
			input: `:"foo ${bar}"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewInterpolatedSymbolLiteralNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							ast.NewInterpolatedStringLiteralNode(
								S(P(1, 1, 2), P(12, 1, 13)),
								[]ast.StringLiteralContentNode{
									ast.NewStringLiteralContentSectionNode(
										S(P(2, 1, 3), P(5, 1, 6)),
										"foo ",
									),
									ast.NewStringInterpolationNode(
										S(P(6, 1, 7), P(11, 1, 12)),
										ast.NewPublicIdentifierNode(S(P(8, 1, 9), P(10, 1, 11)), "bar"),
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

func TestArrayListLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(1, 1, 2)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be empty with capacity": {
			input: "[]:20",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							nil,
							ast.NewIntLiteralNode(S(P(3, 1, 4), P(4, 1, 5)), "20"),
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 3, 1)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(3, 3, 1)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can contain if modifiers": {
			input: "[.1, 'foo', :bar, baz + 5 if baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(1, 1, 2), P(2, 1, 3)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(5, 1, 6), P(9, 1, 10)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(12, 1, 13), P(15, 1, 16)), "bar"),
								ast.NewModifierNode(
									S(P(18, 1, 19), P(31, 1, 32)),
									T(S(P(26, 1, 27), P(27, 1, 28)), token.IF),
									ast.NewBinaryExpressionNode(
										S(P(18, 1, 19), P(24, 1, 25)),
										T(S(P(22, 1, 23), P(22, 1, 23)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(18, 1, 19), P(20, 1, 21)), "baz"),
										ast.NewIntLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "5"),
									),
									ast.NewPublicIdentifierNode(S(P(29, 1, 30), P(31, 1, 32)), "baz"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can contain unless modifiers": {
			input: "[.1, 'foo', :bar, baz + 5 unless baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(36, 1, 37)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(36, 1, 37)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(36, 1, 37)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(1, 1, 2), P(2, 1, 3)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(5, 1, 6), P(9, 1, 10)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(12, 1, 13), P(15, 1, 16)), "bar"),
								ast.NewModifierNode(
									S(P(18, 1, 19), P(35, 1, 36)),
									T(S(P(26, 1, 27), P(31, 1, 32)), token.UNLESS),
									ast.NewBinaryExpressionNode(
										S(P(18, 1, 19), P(24, 1, 25)),
										T(S(P(22, 1, 23), P(22, 1, 23)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(18, 1, 19), P(20, 1, 21)), "baz"),
										ast.NewIntLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "5"),
									),
									ast.NewPublicIdentifierNode(S(P(33, 1, 34), P(35, 1, 36)), "baz"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can contain for modifiers": {
			input: "[.1, 'foo', :bar, baz + 5 for baz in bazz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(41, 1, 42)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(41, 1, 42)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(41, 1, 42)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(1, 1, 2), P(2, 1, 3)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(5, 1, 6), P(9, 1, 10)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(12, 1, 13), P(15, 1, 16)), "bar"),
								ast.NewModifierForInNode(
									S(P(18, 1, 19), P(40, 1, 41)),
									ast.NewBinaryExpressionNode(
										S(P(18, 1, 19), P(24, 1, 25)),
										T(S(P(22, 1, 23), P(22, 1, 23)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(18, 1, 19), P(20, 1, 21)), "baz"),
										ast.NewIntLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "5"),
									),
									ast.NewPublicIdentifierNode(S(P(30, 1, 31), P(32, 1, 33)), "baz"),
									ast.NewPublicIdentifierNode(S(P(37, 1, 38), P(40, 1, 41)), "bazz"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have elements": {
			input: "[.1, 'foo', :bar, baz + 5]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(25, 1, 26)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(25, 1, 26)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(25, 1, 26)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(1, 1, 2), P(2, 1, 3)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(5, 1, 6), P(9, 1, 10)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(12, 1, 13), P(15, 1, 16)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(18, 1, 19), P(24, 1, 25)),
									T(S(P(22, 1, 23), P(22, 1, 23)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(18, 1, 19), P(20, 1, 21)), "baz"),
									ast.NewIntLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "5"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have elements and capacity": {
			input: "[.1, 'foo', :bar, baz + 5]:n",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(1, 1, 2), P(2, 1, 3)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(5, 1, 6), P(9, 1, 10)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(12, 1, 13), P(15, 1, 16)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(18, 1, 19), P(24, 1, 25)),
									T(S(P(22, 1, 23), P(22, 1, 23)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(18, 1, 19), P(20, 1, 21)), "baz"),
									ast.NewIntLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "5"),
								),
							},
							ast.NewPublicIdentifierNode(S(P(27, 1, 28), P(27, 1, 28)), "n"),
						),
					),
				},
			),
		},
		"can have a trailing comma": {
			input: "[.1, 'foo', :bar, baz + 5,]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(1, 1, 2), P(2, 1, 3)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(5, 1, 6), P(9, 1, 10)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(12, 1, 13), P(15, 1, 16)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(18, 1, 19), P(24, 1, 25)),
									T(S(P(22, 1, 23), P(22, 1, 23)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(18, 1, 19), P(20, 1, 21)), "baz"),
									ast.NewIntLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "5"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have explicit indices": {
			input: "[.1, 'foo', 10 => :bar, baz => baz + 5]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 1, 39)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(38, 1, 39)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(38, 1, 39)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(1, 1, 2), P(2, 1, 3)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(5, 1, 6), P(9, 1, 10)), "foo"),
								ast.NewKeyValueExpressionNode(
									S(P(12, 1, 13), P(21, 1, 22)),
									ast.NewIntLiteralNode(S(P(12, 1, 13), P(13, 1, 14)), "10"),
									ast.NewSimpleSymbolLiteralNode(S(P(18, 1, 19), P(21, 1, 22)), "bar"),
								),
								ast.NewKeyValueExpressionNode(
									S(P(24, 1, 25), P(37, 1, 38)),
									ast.NewPublicIdentifierNode(S(P(24, 1, 25), P(26, 1, 27)), "baz"),
									ast.NewBinaryExpressionNode(
										S(P(31, 1, 32), P(37, 1, 38)),
										T(S(P(35, 1, 36), P(35, 1, 36)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(31, 1, 32), P(33, 1, 34)), "baz"),
										ast.NewIntLiteralNode(S(P(37, 1, 38), P(37, 1, 38)), "5"),
									),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have explicit indices with modifiers": {
			input: "[.1, 'foo', 10 => :bar if bar, baz => baz + 5 for baz in bazz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(61, 1, 62)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(61, 1, 62)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(61, 1, 62)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(1, 1, 2), P(2, 1, 3)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(5, 1, 6), P(9, 1, 10)), "foo"),
								ast.NewModifierNode(
									S(P(12, 1, 13), P(28, 1, 29)),
									T(S(P(23, 1, 24), P(24, 1, 25)), token.IF),
									ast.NewKeyValueExpressionNode(
										S(P(12, 1, 13), P(21, 1, 22)),
										ast.NewIntLiteralNode(S(P(12, 1, 13), P(13, 1, 14)), "10"),
										ast.NewSimpleSymbolLiteralNode(S(P(18, 1, 19), P(21, 1, 22)), "bar"),
									),
									ast.NewPublicIdentifierNode(S(P(26, 1, 27), P(28, 1, 29)), "bar"),
								),
								ast.NewModifierForInNode(
									S(P(31, 1, 32), P(60, 1, 61)),
									ast.NewKeyValueExpressionNode(
										S(P(31, 1, 32), P(44, 1, 45)),
										ast.NewPublicIdentifierNode(S(P(31, 1, 32), P(33, 1, 34)), "baz"),
										ast.NewBinaryExpressionNode(
											S(P(38, 1, 39), P(44, 1, 45)),
											T(S(P(42, 1, 43), P(42, 1, 43)), token.PLUS),
											ast.NewPublicIdentifierNode(S(P(38, 1, 39), P(40, 1, 41)), "baz"),
											ast.NewIntLiteralNode(S(P(44, 1, 45), P(44, 1, 45)), "5"),
										),
									),
									ast.NewPublicIdentifierNode(S(P(50, 1, 51), P(52, 1, 53)), "baz"),
									ast.NewPublicIdentifierNode(S(P(57, 1, 58), P(60, 1, 61)), "bazz"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: "[\n.1\n,\n'foo'\n,\n:bar\n,\nbaz + 5\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 9, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(30, 9, 1)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(30, 9, 1)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 2, 1), P(3, 2, 2)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(7, 4, 1), P(11, 4, 5)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 6, 1), P(18, 6, 4)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(22, 8, 1), P(28, 8, 7)),
									T(S(P(26, 8, 5), P(26, 8, 5)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(22, 8, 1), P(24, 8, 3)), "baz"),
									ast.NewIntLiteralNode(S(P(28, 8, 7), P(28, 8, 7)), "5"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can span multiple lines with a trailing comma": {
			input: "[\n.1\n,\n'foo'\n,\n:bar\n,\nbaz + 5,\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 9, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 9, 1)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(31, 9, 1)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 2, 1), P(3, 2, 2)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(7, 4, 1), P(11, 4, 5)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(15, 6, 1), P(18, 6, 4)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(22, 8, 1), P(28, 8, 7)),
									T(S(P(26, 8, 5), P(26, 8, 5)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(22, 8, 1), P(24, 8, 3)), "baz"),
									ast.NewIntLiteralNode(S(P(28, 8, 7), P(28, 8, 7)), "5"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can be nested": {
			input: "[[.1, :+], .2]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewArrayListLiteralNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							[]ast.ExpressionNode{
								ast.NewArrayListLiteralNode(
									S(P(1, 1, 2), P(8, 1, 9)),
									[]ast.ExpressionNode{
										ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
										ast.NewSimpleSymbolLiteralNode(S(P(6, 1, 7), P(7, 1, 8)), "+"),
									},
									nil,
								),
								ast.NewFloatLiteralNode(S(P(11, 1, 12), P(12, 1, 13)), "0.2"),
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

func TestWordArrayListLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "\\w[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewWordArrayListLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be empty with capacity": {
			input: "\\w[]:20",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewWordArrayListLiteralNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							nil,
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(6, 1, 7)), "20"),
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "\\w[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewWordArrayListLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "\\w[foo bar]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewWordArrayListLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.WordCollectionContentNode{
								ast.NewRawStringLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "foo"),
								ast.NewRawStringLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have content with capacity": {
			input: "\\w[foo bar]:n",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewWordArrayListLiteralNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							[]ast.WordCollectionContentNode{
								ast.NewRawStringLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "foo"),
								ast.NewRawStringLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
							},
							ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(12, 1, 13)), "n"),
						),
					),
				},
			),
		},
		"content is interpreted as strings separated by spaces": {
			input: "\\w[.1, 'foo', :bar, baz + 5 if baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(34, 1, 35)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(34, 1, 35)),
						ast.NewWordArrayListLiteralNode(
							S(P(0, 1, 1), P(34, 1, 35)),
							[]ast.WordCollectionContentNode{
								ast.NewRawStringLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), ".1,"),
								ast.NewRawStringLiteralNode(S(P(7, 1, 8), P(12, 1, 13)), "'foo',"),
								ast.NewRawStringLiteralNode(S(P(14, 1, 15), P(18, 1, 19)), ":bar,"),
								ast.NewRawStringLiteralNode(S(P(20, 1, 21), P(22, 1, 23)), "baz"),
								ast.NewRawStringLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "+"),
								ast.NewRawStringLiteralNode(S(P(26, 1, 27), P(26, 1, 27)), "5"),
								ast.NewRawStringLiteralNode(S(P(28, 1, 29), P(29, 1, 30)), "if"),
								ast.NewRawStringLiteralNode(S(P(31, 1, 32), P(33, 1, 34)), "baz"),
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

func TestSymbolArrayListLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "\\s[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewSymbolArrayListLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be empty with capacity": {
			input: "\\s[]:20",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewSymbolArrayListLiteralNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							nil,
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(6, 1, 7)), "20"),
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "\\s[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewSymbolArrayListLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "\\s[foo bar]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewSymbolArrayListLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.SymbolCollectionContentNode{
								ast.NewSimpleSymbolLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have content and capacity": {
			input: "\\s[foo bar]:n",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewSymbolArrayListLiteralNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							[]ast.SymbolCollectionContentNode{
								ast.NewSimpleSymbolLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
							},
							ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(12, 1, 13)), "n"),
						),
					),
				},
			),
		},
		"content is interpreted as strings separated by spaces": {
			input: "\\s[.1, 'foo', :bar, baz + 5 if baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(34, 1, 35)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(34, 1, 35)),
						ast.NewSymbolArrayListLiteralNode(
							S(P(0, 1, 1), P(34, 1, 35)),
							[]ast.SymbolCollectionContentNode{
								ast.NewSimpleSymbolLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), ".1,"),
								ast.NewSimpleSymbolLiteralNode(S(P(7, 1, 8), P(12, 1, 13)), "'foo',"),
								ast.NewSimpleSymbolLiteralNode(S(P(14, 1, 15), P(18, 1, 19)), ":bar,"),
								ast.NewSimpleSymbolLiteralNode(S(P(20, 1, 21), P(22, 1, 23)), "baz"),
								ast.NewSimpleSymbolLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "+"),
								ast.NewSimpleSymbolLiteralNode(S(P(26, 1, 27), P(26, 1, 27)), "5"),
								ast.NewSimpleSymbolLiteralNode(S(P(28, 1, 29), P(29, 1, 30)), "if"),
								ast.NewSimpleSymbolLiteralNode(S(P(31, 1, 32), P(33, 1, 34)), "baz"),
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

func TestHexArrayListLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "\\x[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewHexArrayListLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be empty with capacity": {
			input: "\\x[]:20",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewHexArrayListLiteralNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							nil,
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(6, 1, 7)), "20"),
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "\\x[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewHexArrayListLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "\\x[fff e12]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewHexArrayListLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0xfff"),
								ast.NewIntLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "0xe12"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have content and capacity": {
			input: "\\x[fff e12]:n",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewHexArrayListLiteralNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0xfff"),
								ast.NewIntLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "0xe12"),
							},
							ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(12, 1, 13)), "n"),
						),
					),
				},
			),
		},
		"reports errors about incorrect hex values": {
			input: "\\x[fff fufu 12]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewHexArrayListLiteralNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0xfff"),
								ast.NewInvalidNode(S(P(7, 1, 8), P(10, 1, 11)), V(S(P(7, 1, 8), P(10, 1, 11)), token.ERROR, "invalid int literal")),
								ast.NewIntLiteralNode(S(P(12, 1, 13), P(13, 1, 14)), "0x12"),
							},
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(7, 1, 8), P(10, 1, 11)), "invalid int literal"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestBinArrayListLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "\\b[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewBinArrayListLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be empty with capacity": {
			input: "\\b[]:20",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewBinArrayListLiteralNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							nil,
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(6, 1, 7)), "20"),
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "\\b[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewBinArrayListLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "\\b[111 100]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewBinArrayListLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0b111"),
								ast.NewIntLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "0b100"),
							},
							nil,
						),
					),
				},
			),
		},
		"reports errors about incorrect hex values": {
			input: "\\b[101 fufu 10]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewBinArrayListLiteralNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0b101"),
								ast.NewInvalidNode(S(P(7, 1, 8), P(10, 1, 11)), V(S(P(7, 1, 8), P(10, 1, 11)), token.ERROR, "invalid int literal")),
								ast.NewIntLiteralNode(S(P(12, 1, 13), P(13, 1, 14)), "0b10"),
							},
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(7, 1, 8), P(10, 1, 11)), "invalid int literal"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestArrayTupleLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "%[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(2, 1, 3)),
							nil,
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "%[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 3, 1)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(4, 3, 1)),
							nil,
						),
					),
				},
			),
		},
		"can contain if modifiers": {
			input: "%[.1, 'foo', :bar, baz + 5 if baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewModifierNode(
									S(P(19, 1, 20), P(32, 1, 33)),
									T(S(P(27, 1, 28), P(28, 1, 29)), token.IF),
									ast.NewBinaryExpressionNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
										ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
									),
									ast.NewPublicIdentifierNode(S(P(30, 1, 31), P(32, 1, 33)), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can contain unless modifiers": {
			input: "%[.1, 'foo', :bar, baz + 5 unless baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(37, 1, 38)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(37, 1, 38)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(37, 1, 38)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewModifierNode(
									S(P(19, 1, 20), P(36, 1, 37)),
									T(S(P(27, 1, 28), P(32, 1, 33)), token.UNLESS),
									ast.NewBinaryExpressionNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
										ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
									),
									ast.NewPublicIdentifierNode(S(P(34, 1, 35), P(36, 1, 37)), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can contain for modifiers": {
			input: "%[.1, 'foo', :bar, baz + 5 for baz in bazz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(42, 1, 43)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(42, 1, 43)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(42, 1, 43)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewModifierForInNode(
									S(P(19, 1, 20), P(41, 1, 42)),
									ast.NewBinaryExpressionNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
										ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
									),
									ast.NewPublicIdentifierNode(S(P(31, 1, 32), P(33, 1, 34)), "baz"),
									ast.NewPublicIdentifierNode(S(P(38, 1, 39), P(41, 1, 42)), "bazz"),
								),
							},
						),
					),
				},
			),
		},
		"can have elements": {
			input: "%[.1, 'foo', :bar, baz + 5]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(19, 1, 20), P(25, 1, 26)),
									T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
									ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
								),
							},
						),
					),
				},
			),
		},
		"can have a trailing comma": {
			input: "%[.1, 'foo', :bar, baz + 5,]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(19, 1, 20), P(25, 1, 26)),
									T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
									ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
								),
							},
						),
					),
				},
			),
		},
		"can have explicit indices": {
			input: "%[.1, 'foo', 10 => :bar, baz => baz + 5]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(39, 1, 40)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(39, 1, 40)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(39, 1, 40)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewKeyValueExpressionNode(
									S(P(13, 1, 14), P(22, 1, 23)),
									ast.NewIntLiteralNode(S(P(13, 1, 14), P(14, 1, 15)), "10"),
									ast.NewSimpleSymbolLiteralNode(S(P(19, 1, 20), P(22, 1, 23)), "bar"),
								),
								ast.NewKeyValueExpressionNode(
									S(P(25, 1, 26), P(38, 1, 39)),
									ast.NewPublicIdentifierNode(S(P(25, 1, 26), P(27, 1, 28)), "baz"),
									ast.NewBinaryExpressionNode(
										S(P(32, 1, 33), P(38, 1, 39)),
										T(S(P(36, 1, 37), P(36, 1, 37)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(32, 1, 33), P(34, 1, 35)), "baz"),
										ast.NewIntLiteralNode(S(P(38, 1, 39), P(38, 1, 39)), "5"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can have explicit indices with modifiers": {
			input: "%[.1, 'foo', 10 => :bar if bar, baz => baz + 5 for baz in bazz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(62, 1, 63)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(62, 1, 63)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(62, 1, 63)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewModifierNode(
									S(P(13, 1, 14), P(29, 1, 30)),
									T(S(P(24, 1, 25), P(25, 1, 26)), token.IF),
									ast.NewKeyValueExpressionNode(
										S(P(13, 1, 14), P(22, 1, 23)),
										ast.NewIntLiteralNode(S(P(13, 1, 14), P(14, 1, 15)), "10"),
										ast.NewSimpleSymbolLiteralNode(S(P(19, 1, 20), P(22, 1, 23)), "bar"),
									),
									ast.NewPublicIdentifierNode(S(P(27, 1, 28), P(29, 1, 30)), "bar"),
								),
								ast.NewModifierForInNode(
									S(P(32, 1, 33), P(61, 1, 62)),
									ast.NewKeyValueExpressionNode(
										S(P(32, 1, 33), P(45, 1, 46)),
										ast.NewPublicIdentifierNode(S(P(32, 1, 33), P(34, 1, 35)), "baz"),
										ast.NewBinaryExpressionNode(
											S(P(39, 1, 40), P(45, 1, 46)),
											T(S(P(43, 1, 44), P(43, 1, 44)), token.PLUS),
											ast.NewPublicIdentifierNode(S(P(39, 1, 40), P(41, 1, 42)), "baz"),
											ast.NewIntLiteralNode(S(P(45, 1, 46), P(45, 1, 46)), "5"),
										),
									),
									ast.NewPublicIdentifierNode(S(P(51, 1, 52), P(53, 1, 54)), "baz"),
									ast.NewPublicIdentifierNode(S(P(58, 1, 59), P(61, 1, 62)), "bazz"),
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: "%[\n.1\n,\n'foo'\n,\n:bar\n,\nbaz + 5\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 9, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 9, 1)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(31, 9, 1)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(3, 2, 1), P(4, 2, 2)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 4, 1), P(12, 4, 5)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(16, 6, 1), P(19, 6, 4)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(23, 8, 1), P(29, 8, 7)),
									T(S(P(27, 8, 5), P(27, 8, 5)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(23, 8, 1), P(25, 8, 3)), "baz"),
									ast.NewIntLiteralNode(S(P(29, 8, 7), P(29, 8, 7)), "5"),
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines with a trailing comma": {
			input: "%[\n.1\n,\n'foo'\n,\n:bar\n,\nbaz + 5,\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 9, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 9, 1)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(32, 9, 1)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(3, 2, 1), P(4, 2, 2)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 4, 1), P(12, 4, 5)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(16, 6, 1), P(19, 6, 4)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(23, 8, 1), P(29, 8, 7)),
									T(S(P(27, 8, 5), P(27, 8, 5)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(23, 8, 1), P(25, 8, 3)), "baz"),
									ast.NewIntLiteralNode(S(P(29, 8, 7), P(29, 8, 7)), "5"),
								),
							},
						),
					),
				},
			),
		},
		"can be nested": {
			input: "%[%[.1, :+], .2]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewArrayTupleLiteralNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							[]ast.ExpressionNode{
								ast.NewArrayTupleLiteralNode(
									S(P(2, 1, 3), P(10, 1, 11)),
									[]ast.ExpressionNode{
										ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
										ast.NewSimpleSymbolLiteralNode(S(P(8, 1, 9), P(9, 1, 10)), "+"),
									},
								),
								ast.NewFloatLiteralNode(S(P(13, 1, 14), P(14, 1, 15)), "0.2"),
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

func TestWordArrayTupleLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "%w[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewWordArrayTupleLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "%w[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewWordArrayTupleLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "%w[foo bar]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewWordArrayTupleLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.WordCollectionContentNode{
								ast.NewRawStringLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "foo"),
								ast.NewRawStringLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
							},
						),
					),
				},
			),
		},
		"content is interpreted as strings separated by spaces": {
			input: "%w[.1, 'foo', :bar, baz + 5 if baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(34, 1, 35)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(34, 1, 35)),
						ast.NewWordArrayTupleLiteralNode(
							S(P(0, 1, 1), P(34, 1, 35)),
							[]ast.WordCollectionContentNode{
								ast.NewRawStringLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), ".1,"),
								ast.NewRawStringLiteralNode(S(P(7, 1, 8), P(12, 1, 13)), "'foo',"),
								ast.NewRawStringLiteralNode(S(P(14, 1, 15), P(18, 1, 19)), ":bar,"),
								ast.NewRawStringLiteralNode(S(P(20, 1, 21), P(22, 1, 23)), "baz"),
								ast.NewRawStringLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "+"),
								ast.NewRawStringLiteralNode(S(P(26, 1, 27), P(26, 1, 27)), "5"),
								ast.NewRawStringLiteralNode(S(P(28, 1, 29), P(29, 1, 30)), "if"),
								ast.NewRawStringLiteralNode(S(P(31, 1, 32), P(33, 1, 34)), "baz"),
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

func TestSymbolArrayTupleLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "%s[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewSymbolArrayTupleLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "%s[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewSymbolArrayTupleLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "%s[foo bar]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewSymbolArrayTupleLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.SymbolCollectionContentNode{
								ast.NewSimpleSymbolLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
							},
						),
					),
				},
			),
		},
		"content is interpreted as strings separated by spaces": {
			input: "%s[.1, 'foo', :bar, baz + 5 if baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(34, 1, 35)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(34, 1, 35)),
						ast.NewSymbolArrayTupleLiteralNode(
							S(P(0, 1, 1), P(34, 1, 35)),
							[]ast.SymbolCollectionContentNode{
								ast.NewSimpleSymbolLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), ".1,"),
								ast.NewSimpleSymbolLiteralNode(S(P(7, 1, 8), P(12, 1, 13)), "'foo',"),
								ast.NewSimpleSymbolLiteralNode(S(P(14, 1, 15), P(18, 1, 19)), ":bar,"),
								ast.NewSimpleSymbolLiteralNode(S(P(20, 1, 21), P(22, 1, 23)), "baz"),
								ast.NewSimpleSymbolLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "+"),
								ast.NewSimpleSymbolLiteralNode(S(P(26, 1, 27), P(26, 1, 27)), "5"),
								ast.NewSimpleSymbolLiteralNode(S(P(28, 1, 29), P(29, 1, 30)), "if"),
								ast.NewSimpleSymbolLiteralNode(S(P(31, 1, 32), P(33, 1, 34)), "baz"),
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

func TestHexArrayTupleLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "%x[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewHexArrayTupleLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "%x[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewHexArrayTupleLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "%x[fff e12]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewHexArrayTupleLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0xfff"),
								ast.NewIntLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "0xe12"),
							},
						),
					),
				},
			),
		},
		"reports errors about incorrect hex values": {
			input: "%x[fff fufu 12]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewHexArrayTupleLiteralNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0xfff"),
								ast.NewInvalidNode(S(P(7, 1, 8), P(10, 1, 11)), V(S(P(7, 1, 8), P(10, 1, 11)), token.ERROR, "invalid int literal")),
								ast.NewIntLiteralNode(S(P(12, 1, 13), P(13, 1, 14)), "0x12"),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(7, 1, 8), P(10, 1, 11)), "invalid int literal"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestBinArrayTupleLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "%b[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewBinArrayTupleLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "%b[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewBinArrayTupleLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "%b[111 100]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewBinArrayTupleLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0b111"),
								ast.NewIntLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "0b100"),
							},
						),
					),
				},
			),
		},
		"reports errors about incorrect hex values": {
			input: "%b[101 fufu 10]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewBinArrayTupleLiteralNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0b101"),
								ast.NewInvalidNode(S(P(7, 1, 8), P(10, 1, 11)), V(S(P(7, 1, 8), P(10, 1, 11)), token.ERROR, "invalid int literal")),
								ast.NewIntLiteralNode(S(P(12, 1, 13), P(13, 1, 14)), "0b10"),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(7, 1, 8), P(10, 1, 11)), "invalid int literal"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestHashSetLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "^[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(2, 1, 3)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be empty with capacity": {
			input: "^[]:20",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 1, 6)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(5, 1, 6)),
							nil,
							ast.NewIntLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "20"),
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "^[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 3, 1)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(4, 3, 1)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can contain if modifiers": {
			input: "^[.1, 'foo', :bar, baz + 5 if baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewModifierNode(
									S(P(19, 1, 20), P(32, 1, 33)),
									T(S(P(27, 1, 28), P(28, 1, 29)), token.IF),
									ast.NewBinaryExpressionNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
										ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
									),
									ast.NewPublicIdentifierNode(S(P(30, 1, 31), P(32, 1, 33)), "baz"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can contain unless modifiers": {
			input: "^[.1, 'foo', :bar, baz + 5 unless baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(37, 1, 38)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(37, 1, 38)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(37, 1, 38)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewModifierNode(
									S(P(19, 1, 20), P(36, 1, 37)),
									T(S(P(27, 1, 28), P(32, 1, 33)), token.UNLESS),
									ast.NewBinaryExpressionNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
										ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
									),
									ast.NewPublicIdentifierNode(S(P(34, 1, 35), P(36, 1, 37)), "baz"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can contain for modifiers": {
			input: "^[.1, 'foo', :bar, baz + 5 for baz in bazz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(42, 1, 43)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(42, 1, 43)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(42, 1, 43)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewModifierForInNode(
									S(P(19, 1, 20), P(41, 1, 42)),
									ast.NewBinaryExpressionNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
										ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
									),
									ast.NewPublicIdentifierNode(S(P(31, 1, 32), P(33, 1, 34)), "baz"),
									ast.NewPublicIdentifierNode(S(P(38, 1, 39), P(41, 1, 42)), "bazz"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have elements": {
			input: "^[.1, 'foo', :bar, baz + 5]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(19, 1, 20), P(25, 1, 26)),
									T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
									ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have elements and capacity": {
			input: "^[.1, 'foo', :bar, baz + 5]:n",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 1, 29)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(28, 1, 29)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(19, 1, 20), P(25, 1, 26)),
									T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
									ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
								),
							},
							ast.NewPublicIdentifierNode(S(P(28, 1, 29), P(28, 1, 29)), "n"),
						),
					),
				},
			),
		},
		"can have a trailing comma": {
			input: "^[.1, 'foo', :bar, baz + 5,]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(19, 1, 20), P(25, 1, 26)),
									T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
									ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot have explicit indices": {
			input: "^[.1, 'foo', 10 => :bar, baz => baz + 5]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(39, 1, 40)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(16, 1, 17), P(39, 1, 40)),
						ast.NewInvalidNode(
							S(P(16, 1, 17), P(17, 1, 18)),
							T(S(P(16, 1, 17), P(17, 1, 18)), token.THICK_ARROW),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(16, 1, 17), P(17, 1, 18)), "unexpected =>, expected ]"),
			},
		},
		"can span multiple lines": {
			input: "^[\n.1\n,\n'foo'\n,\n:bar\n,\nbaz + 5\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(31, 9, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(31, 9, 1)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(31, 9, 1)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(3, 2, 1), P(4, 2, 2)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 4, 1), P(12, 4, 5)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(16, 6, 1), P(19, 6, 4)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(23, 8, 1), P(29, 8, 7)),
									T(S(P(27, 8, 5), P(27, 8, 5)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(23, 8, 1), P(25, 8, 3)), "baz"),
									ast.NewIntLiteralNode(S(P(29, 8, 7), P(29, 8, 7)), "5"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can span multiple lines with a trailing comma": {
			input: "^[\n.1\n,\n'foo'\n,\n:bar\n,\nbaz + 5,\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 9, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 9, 1)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(32, 9, 1)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(3, 2, 1), P(4, 2, 2)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(8, 4, 1), P(12, 4, 5)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(16, 6, 1), P(19, 6, 4)), "bar"),
								ast.NewBinaryExpressionNode(
									S(P(23, 8, 1), P(29, 8, 7)),
									T(S(P(27, 8, 5), P(27, 8, 5)), token.PLUS),
									ast.NewPublicIdentifierNode(S(P(23, 8, 1), P(25, 8, 3)), "baz"),
									ast.NewIntLiteralNode(S(P(29, 8, 7), P(29, 8, 7)), "5"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can be nested": {
			input: "^[^[.1, :+], .2]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewHashSetLiteralNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							[]ast.ExpressionNode{
								ast.NewHashSetLiteralNode(
									S(P(2, 1, 3), P(10, 1, 11)),
									[]ast.ExpressionNode{
										ast.NewFloatLiteralNode(S(P(4, 1, 5), P(5, 1, 6)), "0.1"),
										ast.NewSimpleSymbolLiteralNode(S(P(8, 1, 9), P(9, 1, 10)), "+"),
									},
									nil,
								),
								ast.NewFloatLiteralNode(S(P(13, 1, 14), P(14, 1, 15)), "0.2"),
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

func TestWordHashSetLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "^w[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewWordHashSetLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be empty with capacity": {
			input: "^w[]:20",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewWordHashSetLiteralNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							nil,
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(6, 1, 7)), "20"),
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "^w[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewWordHashSetLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "^w[foo bar]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewWordHashSetLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.WordCollectionContentNode{
								ast.NewRawStringLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "foo"),
								ast.NewRawStringLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have content and capacity": {
			input: "^w[foo bar]:n",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewWordHashSetLiteralNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							[]ast.WordCollectionContentNode{
								ast.NewRawStringLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "foo"),
								ast.NewRawStringLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
							},
							ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(12, 1, 13)), "n"),
						),
					),
				},
			),
		},
		"content is interpreted as strings separated by spaces": {
			input: "^w[.1, 'foo', :bar, baz + 5 if baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(34, 1, 35)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(34, 1, 35)),
						ast.NewWordHashSetLiteralNode(
							S(P(0, 1, 1), P(34, 1, 35)),
							[]ast.WordCollectionContentNode{
								ast.NewRawStringLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), ".1,"),
								ast.NewRawStringLiteralNode(S(P(7, 1, 8), P(12, 1, 13)), "'foo',"),
								ast.NewRawStringLiteralNode(S(P(14, 1, 15), P(18, 1, 19)), ":bar,"),
								ast.NewRawStringLiteralNode(S(P(20, 1, 21), P(22, 1, 23)), "baz"),
								ast.NewRawStringLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "+"),
								ast.NewRawStringLiteralNode(S(P(26, 1, 27), P(26, 1, 27)), "5"),
								ast.NewRawStringLiteralNode(S(P(28, 1, 29), P(29, 1, 30)), "if"),
								ast.NewRawStringLiteralNode(S(P(31, 1, 32), P(33, 1, 34)), "baz"),
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

func TestSymbolHashSetLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "^s[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewSymbolHashSetLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be empty with capacity": {
			input: "^s[]:20",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewSymbolHashSetLiteralNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							nil,
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(6, 1, 7)), "20"),
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "^s[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewSymbolHashSetLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "^s[foo bar]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewSymbolHashSetLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.SymbolCollectionContentNode{
								ast.NewSimpleSymbolLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have content and capacity": {
			input: "^s[foo bar]:n",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewSymbolHashSetLiteralNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							[]ast.SymbolCollectionContentNode{
								ast.NewSimpleSymbolLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
							},
							ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(12, 1, 13)), "n"),
						),
					),
				},
			),
		},
		"content is interpreted as strings separated by spaces": {
			input: "^s[.1, 'foo', :bar, baz + 5 if baz]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(34, 1, 35)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(34, 1, 35)),
						ast.NewSymbolHashSetLiteralNode(
							S(P(0, 1, 1), P(34, 1, 35)),
							[]ast.SymbolCollectionContentNode{
								ast.NewSimpleSymbolLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), ".1,"),
								ast.NewSimpleSymbolLiteralNode(S(P(7, 1, 8), P(12, 1, 13)), "'foo',"),
								ast.NewSimpleSymbolLiteralNode(S(P(14, 1, 15), P(18, 1, 19)), ":bar,"),
								ast.NewSimpleSymbolLiteralNode(S(P(20, 1, 21), P(22, 1, 23)), "baz"),
								ast.NewSimpleSymbolLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "+"),
								ast.NewSimpleSymbolLiteralNode(S(P(26, 1, 27), P(26, 1, 27)), "5"),
								ast.NewSimpleSymbolLiteralNode(S(P(28, 1, 29), P(29, 1, 30)), "if"),
								ast.NewSimpleSymbolLiteralNode(S(P(31, 1, 32), P(33, 1, 34)), "baz"),
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

func TestHexHashSetLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "^x[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewHexHashSetLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be empty with capacity": {
			input: "^x[]:20",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewHexHashSetLiteralNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							nil,
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(6, 1, 7)), "20"),
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "^x[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewHexHashSetLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "^x[fff e12]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewHexHashSetLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0xfff"),
								ast.NewIntLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "0xe12"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have content and capacity": {
			input: "^x[fff e12]:n",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewHexHashSetLiteralNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0xfff"),
								ast.NewIntLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "0xe12"),
							},
							ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(12, 1, 13)), "n"),
						),
					),
				},
			),
		},
		"reports errors about incorrect hex values": {
			input: "^x[fff fufu 12]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewHexHashSetLiteralNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0xfff"),
								ast.NewInvalidNode(S(P(7, 1, 8), P(10, 1, 11)), V(S(P(7, 1, 8), P(10, 1, 11)), token.ERROR, "invalid int literal")),
								ast.NewIntLiteralNode(S(P(12, 1, 13), P(13, 1, 14)), "0x12"),
							},
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(7, 1, 8), P(10, 1, 11)), "invalid int literal"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestBinHashSetLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "^b[]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewBinHashSetLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be empty with capacity": {
			input: "^b[]:20",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewBinHashSetLiteralNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							nil,
							ast.NewIntLiteralNode(S(P(5, 1, 6), P(6, 1, 7)), "20"),
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "^b[\n\n]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(5, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(5, 3, 1)),
						ast.NewBinHashSetLiteralNode(
							S(P(0, 1, 1), P(5, 3, 1)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have content": {
			input: "^b[111 100]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewBinHashSetLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0b111"),
								ast.NewIntLiteralNode(S(P(7, 1, 8), P(9, 1, 10)), "0b100"),
							},
							nil,
						),
					),
				},
			),
		},
		"reports errors about incorrect hex values": {
			input: "^b[101 fufu 10]",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewBinHashSetLiteralNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							[]ast.IntCollectionContentNode{
								ast.NewIntLiteralNode(S(P(3, 1, 4), P(5, 1, 6)), "0b101"),
								ast.NewInvalidNode(S(P(7, 1, 8), P(10, 1, 11)), V(S(P(7, 1, 8), P(10, 1, 11)), token.ERROR, "invalid int literal")),
								ast.NewIntLiteralNode(S(P(12, 1, 13), P(13, 1, 14)), "0b10"),
							},
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(7, 1, 8), P(10, 1, 11)), "invalid int literal"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestHashMapLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "{}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(1, 1, 2)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(1, 1, 2)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be empty with capacity": {
			input: "{}:20",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							nil,
							ast.NewIntLiteralNode(S(P(3, 1, 4), P(4, 1, 5)), "20"),
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "{\n\n}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 3, 1)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(3, 3, 1)),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot contain elements other than key value pairs and identifiers": {
			input: "{.1, 'foo', :bar, baz + 5 if baz}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(32, 1, 33)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(32, 1, 33)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(32, 1, 33)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(1, 1, 2), P(2, 1, 3)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(5, 1, 6), P(9, 1, 10)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(12, 1, 13), P(15, 1, 16)), "bar"),
								ast.NewModifierNode(
									S(P(18, 1, 19), P(31, 1, 32)),
									T(S(P(26, 1, 27), P(27, 1, 28)), token.IF),
									ast.NewBinaryExpressionNode(
										S(P(18, 1, 19), P(24, 1, 25)),
										T(S(P(22, 1, 23), P(22, 1, 23)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(18, 1, 19), P(20, 1, 21)), "baz"),
										ast.NewIntLiteralNode(S(P(24, 1, 25), P(24, 1, 25)), "5"),
									),
									ast.NewPublicIdentifierNode(S(P(29, 1, 30), P(31, 1, 32)), "baz"),
								),
							},
							nil,
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(1, 1, 2), P(2, 1, 3)), "expected a key-value pair, map literals should consist of key-value pairs"),
				errors.NewError(L("<main>", P(5, 1, 6), P(9, 1, 10)), "expected a key-value pair, map literals should consist of key-value pairs"),
				errors.NewError(L("<main>", P(12, 1, 13), P(15, 1, 16)), "expected a key-value pair, map literals should consist of key-value pairs"),
				errors.NewError(L("<main>", P(18, 1, 19), P(24, 1, 25)), "expected a key-value pair, map literals should consist of key-value pairs"),
			},
		},
		"can contain any expression as key with thick arrows": {
			input: "{Math::PI => 3, foo => foo && bar, 5 => 'bar', 'baz' => :bar, a + 5 => 1, n.to_string() => n}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(92, 1, 93)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(92, 1, 93)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(92, 1, 93)),
							[]ast.ExpressionNode{
								ast.NewKeyValueExpressionNode(
									S(P(1, 1, 2), P(13, 1, 14)),
									ast.NewConstantLookupNode(
										S(P(1, 1, 2), P(8, 1, 9)),
										ast.NewPublicConstantNode(S(P(1, 1, 2), P(4, 1, 5)), "Math"),
										ast.NewPublicConstantNode(S(P(7, 1, 8), P(8, 1, 9)), "PI"),
									),
									ast.NewIntLiteralNode(S(P(13, 1, 14), P(13, 1, 14)), "3"),
								),
								ast.NewKeyValueExpressionNode(
									S(P(16, 1, 17), P(32, 1, 33)),
									ast.NewPublicIdentifierNode(S(P(16, 1, 17), P(18, 1, 19)), "foo"),
									ast.NewLogicalExpressionNode(
										S(P(23, 1, 24), P(32, 1, 33)),
										T(S(P(27, 1, 28), P(28, 1, 29)), token.AND_AND),
										ast.NewPublicIdentifierNode(S(P(23, 1, 24), P(25, 1, 26)), "foo"),
										ast.NewPublicIdentifierNode(S(P(30, 1, 31), P(32, 1, 33)), "bar"),
									),
								),
								ast.NewKeyValueExpressionNode(
									S(P(35, 1, 36), P(44, 1, 45)),
									ast.NewIntLiteralNode(S(P(35, 1, 36), P(35, 1, 36)), "5"),
									ast.NewRawStringLiteralNode(S(P(40, 1, 41), P(44, 1, 45)), "bar"),
								),
								ast.NewKeyValueExpressionNode(
									S(P(47, 1, 48), P(59, 1, 60)),
									ast.NewRawStringLiteralNode(S(P(47, 1, 48), P(51, 1, 52)), "baz"),
									ast.NewSimpleSymbolLiteralNode(S(P(56, 1, 57), P(59, 1, 60)), "bar"),
								),
								ast.NewKeyValueExpressionNode(
									S(P(62, 1, 63), P(71, 1, 72)),
									ast.NewBinaryExpressionNode(
										S(P(62, 1, 63), P(66, 1, 67)),
										T(S(P(64, 1, 65), P(64, 1, 65)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(62, 1, 63), P(62, 1, 63)), "a"),
										ast.NewIntLiteralNode(S(P(66, 1, 67), P(66, 1, 67)), "5"),
									),
									ast.NewIntLiteralNode(S(P(71, 1, 72), P(71, 1, 72)), "1"),
								),
								ast.NewKeyValueExpressionNode(
									S(P(74, 1, 75), P(91, 1, 92)),
									ast.NewMethodCallNode(
										S(P(74, 1, 75), P(86, 1, 87)),
										ast.NewPublicIdentifierNode(S(P(74, 1, 75), P(74, 1, 75)), "n"),
										false,
										"to_string",
										nil,
										nil,
									),
									ast.NewPublicIdentifierNode(S(P(91, 1, 92), P(91, 1, 92)), "n"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have shorthand symbol keys": {
			input: "{foo: :bar}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(1, 1, 2), P(9, 1, 10)),
									"foo",
									ast.NewSimpleSymbolLiteralNode(S(P(6, 1, 7), P(9, 1, 10)), "bar"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have elements and capacity": {
			input: "{foo: :bar}:n",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(1, 1, 2), P(9, 1, 10)),
									"foo",
									ast.NewSimpleSymbolLiteralNode(S(P(6, 1, 7), P(9, 1, 10)), "bar"),
								),
							},
							ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(12, 1, 13)), "n"),
						),
					),
				},
			),
		},
		"can contain for modifiers": {
			input: "{foo: bar, baz => baz.to_int for baz in bazz}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(44, 1, 45)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(44, 1, 45)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(44, 1, 45)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(1, 1, 2), P(8, 1, 9)),
									"foo",
									ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
								),
								ast.NewModifierForInNode(
									S(P(11, 1, 12), P(43, 1, 44)),
									ast.NewKeyValueExpressionNode(
										S(P(11, 1, 12), P(27, 1, 28)),
										ast.NewPublicIdentifierNode(S(P(11, 1, 12), P(13, 1, 14)), "baz"),
										ast.NewAttributeAccessNode(
											S(P(18, 1, 19), P(27, 1, 28)),
											ast.NewPublicIdentifierNode(S(P(18, 1, 19), P(20, 1, 21)), "baz"),
											"to_int",
										),
									),
									ast.NewPublicIdentifierNode(S(P(33, 1, 34), P(35, 1, 36)), "baz"),
									ast.NewPublicIdentifierNode(S(P(40, 1, 41), P(43, 1, 44)), "bazz"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can contain if modifiers": {
			input: "{foo: bar, baz => baz.to_int if baz}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(35, 1, 36)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(35, 1, 36)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(35, 1, 36)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(1, 1, 2), P(8, 1, 9)),
									"foo",
									ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
								),
								ast.NewModifierNode(
									S(P(11, 1, 12), P(34, 1, 35)),
									T(S(P(29, 1, 30), P(30, 1, 31)), token.IF),
									ast.NewKeyValueExpressionNode(
										S(P(11, 1, 12), P(27, 1, 28)),
										ast.NewPublicIdentifierNode(S(P(11, 1, 12), P(13, 1, 14)), "baz"),
										ast.NewAttributeAccessNode(
											S(P(18, 1, 19), P(27, 1, 28)),
											ast.NewPublicIdentifierNode(S(P(18, 1, 19), P(20, 1, 21)), "baz"),
											"to_int",
										),
									),
									ast.NewPublicIdentifierNode(S(P(32, 1, 33), P(34, 1, 35)), "baz"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have a trailing comma": {
			input: "{foo: bar, baz => baz.to_int if baz,}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(36, 1, 37)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(36, 1, 37)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(36, 1, 37)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(1, 1, 2), P(8, 1, 9)),
									"foo",
									ast.NewPublicIdentifierNode(S(P(6, 1, 7), P(8, 1, 9)), "bar"),
								),
								ast.NewModifierNode(
									S(P(11, 1, 12), P(34, 1, 35)),
									T(S(P(29, 1, 30), P(30, 1, 31)), token.IF),
									ast.NewKeyValueExpressionNode(
										S(P(11, 1, 12), P(27, 1, 28)),
										ast.NewPublicIdentifierNode(S(P(11, 1, 12), P(13, 1, 14)), "baz"),
										ast.NewAttributeAccessNode(
											S(P(18, 1, 19), P(27, 1, 28)),
											ast.NewPublicIdentifierNode(S(P(18, 1, 19), P(20, 1, 21)), "baz"),
											"to_int",
										),
									),
									ast.NewPublicIdentifierNode(S(P(32, 1, 33), P(34, 1, 35)), "baz"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: "{\nfoo:\nbar,\nbaz =>\nbaz.to_int if\nbaz\n}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(37, 7, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(37, 7, 1)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(37, 7, 1)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(2, 2, 1), P(9, 3, 3)),
									"foo",
									ast.NewPublicIdentifierNode(S(P(7, 3, 1), P(9, 3, 3)), "bar"),
								),
								ast.NewModifierNode(
									S(P(12, 4, 1), P(35, 6, 3)),
									T(S(P(30, 5, 12), P(31, 5, 13)), token.IF),
									ast.NewKeyValueExpressionNode(
										S(P(12, 4, 1), P(28, 5, 10)),
										ast.NewPublicIdentifierNode(S(P(12, 4, 1), P(14, 4, 3)), "baz"),
										ast.NewAttributeAccessNode(
											S(P(19, 5, 1), P(28, 5, 10)),
											ast.NewPublicIdentifierNode(S(P(19, 5, 1), P(21, 5, 3)), "baz"),
											"to_int",
										),
									),
									ast.NewPublicIdentifierNode(S(P(33, 6, 1), P(35, 6, 3)), "baz"),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can span multiple lines with a trailing comma": {
			input: "{\nfoo:\nbar,\nbaz =>\nbaz.to_int if\nbaz,\n}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 7, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(38, 7, 1)),
						ast.NewHashMapLiteralNode(
							S(P(0, 1, 1), P(38, 7, 1)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(2, 2, 1), P(9, 3, 3)),
									"foo",
									ast.NewPublicIdentifierNode(S(P(7, 3, 1), P(9, 3, 3)), "bar"),
								),
								ast.NewModifierNode(
									S(P(12, 4, 1), P(35, 6, 3)),
									T(S(P(30, 5, 12), P(31, 5, 13)), token.IF),
									ast.NewKeyValueExpressionNode(
										S(P(12, 4, 1), P(28, 5, 10)),
										ast.NewPublicIdentifierNode(S(P(12, 4, 1), P(14, 4, 3)), "baz"),
										ast.NewAttributeAccessNode(
											S(P(19, 5, 1), P(28, 5, 10)),
											ast.NewPublicIdentifierNode(S(P(19, 5, 1), P(21, 5, 3)), "baz"),
											"to_int",
										),
									),
									ast.NewPublicIdentifierNode(S(P(33, 6, 1), P(35, 6, 3)), "baz"),
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

func TestRecordLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "%{}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						ast.NewHashRecordLiteralNode(
							S(P(0, 1, 1), P(2, 1, 3)),
							nil,
						),
					),
				},
			),
		},
		"can be empty with newlines": {
			input: "%{\n\n}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 3, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 3, 1)),
						ast.NewHashRecordLiteralNode(
							S(P(0, 1, 1), P(4, 3, 1)),
							nil,
						),
					),
				},
			),
		},
		"cannot contain elements other than key value pairs and identifiers": {
			input: "%{.1, 'foo', :bar, baz + 5 if baz}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(33, 1, 34)),
						ast.NewHashRecordLiteralNode(
							S(P(0, 1, 1), P(33, 1, 34)),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(S(P(2, 1, 3), P(3, 1, 4)), "0.1"),
								ast.NewRawStringLiteralNode(S(P(6, 1, 7), P(10, 1, 11)), "foo"),
								ast.NewSimpleSymbolLiteralNode(S(P(13, 1, 14), P(16, 1, 17)), "bar"),
								ast.NewModifierNode(
									S(P(19, 1, 20), P(32, 1, 33)),
									T(S(P(27, 1, 28), P(28, 1, 29)), token.IF),
									ast.NewBinaryExpressionNode(
										S(P(19, 1, 20), P(25, 1, 26)),
										T(S(P(23, 1, 24), P(23, 1, 24)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
										ast.NewIntLiteralNode(S(P(25, 1, 26), P(25, 1, 26)), "5"),
									),
									ast.NewPublicIdentifierNode(S(P(30, 1, 31), P(32, 1, 33)), "baz"),
								),
							},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(2, 1, 3), P(3, 1, 4)), "expected a key-value pair, map literals should consist of key-value pairs"),
				errors.NewError(L("<main>", P(6, 1, 7), P(10, 1, 11)), "expected a key-value pair, map literals should consist of key-value pairs"),
				errors.NewError(L("<main>", P(13, 1, 14), P(16, 1, 17)), "expected a key-value pair, map literals should consist of key-value pairs"),
				errors.NewError(L("<main>", P(19, 1, 20), P(25, 1, 26)), "expected a key-value pair, map literals should consist of key-value pairs"),
			},
		},
		"can contain any expression as key with thick arrows": {
			input: "%{Math::PI => 3, foo => foo && bar, 5 => 'bar', 'baz' => :bar, a + 5 => 1, n.to_string() => n}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(93, 1, 94)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(93, 1, 94)),
						ast.NewHashRecordLiteralNode(
							S(P(0, 1, 1), P(93, 1, 94)),
							[]ast.ExpressionNode{
								ast.NewKeyValueExpressionNode(
									S(P(2, 1, 3), P(14, 1, 15)),
									ast.NewConstantLookupNode(
										S(P(2, 1, 3), P(9, 1, 10)),
										ast.NewPublicConstantNode(S(P(2, 1, 3), P(5, 1, 6)), "Math"),
										ast.NewPublicConstantNode(S(P(8, 1, 9), P(9, 1, 10)), "PI"),
									),
									ast.NewIntLiteralNode(S(P(14, 1, 15), P(14, 1, 15)), "3"),
								),
								ast.NewKeyValueExpressionNode(
									S(P(17, 1, 18), P(33, 1, 34)),
									ast.NewPublicIdentifierNode(S(P(17, 1, 18), P(19, 1, 20)), "foo"),
									ast.NewLogicalExpressionNode(
										S(P(24, 1, 25), P(33, 1, 34)),
										T(S(P(28, 1, 29), P(29, 1, 30)), token.AND_AND),
										ast.NewPublicIdentifierNode(S(P(24, 1, 25), P(26, 1, 27)), "foo"),
										ast.NewPublicIdentifierNode(S(P(31, 1, 32), P(33, 1, 34)), "bar"),
									),
								),
								ast.NewKeyValueExpressionNode(
									S(P(36, 1, 37), P(45, 1, 46)),
									ast.NewIntLiteralNode(S(P(36, 1, 37), P(36, 1, 37)), "5"),
									ast.NewRawStringLiteralNode(S(P(41, 1, 42), P(45, 1, 46)), "bar"),
								),
								ast.NewKeyValueExpressionNode(
									S(P(48, 1, 49), P(60, 1, 61)),
									ast.NewRawStringLiteralNode(S(P(48, 1, 49), P(52, 1, 53)), "baz"),
									ast.NewSimpleSymbolLiteralNode(S(P(57, 1, 58), P(60, 1, 61)), "bar"),
								),
								ast.NewKeyValueExpressionNode(
									S(P(63, 1, 64), P(72, 1, 73)),
									ast.NewBinaryExpressionNode(
										S(P(63, 1, 64), P(67, 1, 68)),
										T(S(P(65, 1, 66), P(65, 1, 66)), token.PLUS),
										ast.NewPublicIdentifierNode(S(P(63, 1, 64), P(63, 1, 64)), "a"),
										ast.NewIntLiteralNode(S(P(67, 1, 68), P(67, 1, 68)), "5"),
									),
									ast.NewIntLiteralNode(S(P(72, 1, 73), P(72, 1, 73)), "1"),
								),
								ast.NewKeyValueExpressionNode(
									S(P(75, 1, 76), P(92, 1, 93)),
									ast.NewMethodCallNode(
										S(P(75, 1, 76), P(87, 1, 88)),
										ast.NewPublicIdentifierNode(S(P(75, 1, 76), P(75, 1, 76)), "n"),
										false,
										"to_string",
										nil,
										nil,
									),
									ast.NewPublicIdentifierNode(S(P(92, 1, 93), P(92, 1, 93)), "n"),
								),
							},
						),
					),
				},
			),
		},
		"can have shorthand symbol keys": {
			input: "%{foo: :bar}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewHashRecordLiteralNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(2, 1, 3), P(10, 1, 11)),
									"foo",
									ast.NewSimpleSymbolLiteralNode(S(P(7, 1, 8), P(10, 1, 11)), "bar"),
								),
							},
						),
					),
				},
			),
		},
		"can contain for modifiers": {
			input: "%{foo: bar, baz => baz.to_int for baz in bazz}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(45, 1, 46)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(45, 1, 46)),
						ast.NewHashRecordLiteralNode(
							S(P(0, 1, 1), P(45, 1, 46)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(2, 1, 3), P(9, 1, 10)),
									"foo",
									ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
								),
								ast.NewModifierForInNode(
									S(P(12, 1, 13), P(44, 1, 45)),
									ast.NewKeyValueExpressionNode(
										S(P(12, 1, 13), P(28, 1, 29)),
										ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(14, 1, 15)), "baz"),
										ast.NewAttributeAccessNode(
											S(P(19, 1, 20), P(28, 1, 29)),
											ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
											"to_int",
										),
									),
									ast.NewPublicIdentifierNode(S(P(34, 1, 35), P(36, 1, 37)), "baz"),
									ast.NewPublicIdentifierNode(S(P(41, 1, 42), P(44, 1, 45)), "bazz"),
								),
							},
						),
					),
				},
			),
		},
		"can contain if modifiers": {
			input: "%{foo: bar, baz => baz.to_int if baz}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(36, 1, 37)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(36, 1, 37)),
						ast.NewHashRecordLiteralNode(
							S(P(0, 1, 1), P(36, 1, 37)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(2, 1, 3), P(9, 1, 10)),
									"foo",
									ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
								),
								ast.NewModifierNode(
									S(P(12, 1, 13), P(35, 1, 36)),
									T(S(P(30, 1, 31), P(31, 1, 32)), token.IF),
									ast.NewKeyValueExpressionNode(
										S(P(12, 1, 13), P(28, 1, 29)),
										ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(14, 1, 15)), "baz"),
										ast.NewAttributeAccessNode(
											S(P(19, 1, 20), P(28, 1, 29)),
											ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
											"to_int",
										),
									),
									ast.NewPublicIdentifierNode(S(P(33, 1, 34), P(35, 1, 36)), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can have a trailing comma": {
			input: "%{foo: bar, baz => baz.to_int if baz,}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(37, 1, 38)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(37, 1, 38)),
						ast.NewHashRecordLiteralNode(
							S(P(0, 1, 1), P(37, 1, 38)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(2, 1, 3), P(9, 1, 10)),
									"foo",
									ast.NewPublicIdentifierNode(S(P(7, 1, 8), P(9, 1, 10)), "bar"),
								),
								ast.NewModifierNode(
									S(P(12, 1, 13), P(35, 1, 36)),
									T(S(P(30, 1, 31), P(31, 1, 32)), token.IF),
									ast.NewKeyValueExpressionNode(
										S(P(12, 1, 13), P(28, 1, 29)),
										ast.NewPublicIdentifierNode(S(P(12, 1, 13), P(14, 1, 15)), "baz"),
										ast.NewAttributeAccessNode(
											S(P(19, 1, 20), P(28, 1, 29)),
											ast.NewPublicIdentifierNode(S(P(19, 1, 20), P(21, 1, 22)), "baz"),
											"to_int",
										),
									),
									ast.NewPublicIdentifierNode(S(P(33, 1, 34), P(35, 1, 36)), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines": {
			input: "%{\nfoo:\nbar,\nbaz =>\nbaz.to_int if\nbaz\n}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 7, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(38, 7, 1)),
						ast.NewHashRecordLiteralNode(
							S(P(0, 1, 1), P(38, 7, 1)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(3, 2, 1), P(10, 3, 3)),
									"foo",
									ast.NewPublicIdentifierNode(S(P(8, 3, 1), P(10, 3, 3)), "bar"),
								),
								ast.NewModifierNode(
									S(P(13, 4, 1), P(36, 6, 3)),
									T(S(P(31, 5, 12), P(32, 5, 13)), token.IF),
									ast.NewKeyValueExpressionNode(
										S(P(13, 4, 1), P(29, 5, 10)),
										ast.NewPublicIdentifierNode(S(P(13, 4, 1), P(15, 4, 3)), "baz"),
										ast.NewAttributeAccessNode(
											S(P(20, 5, 1), P(29, 5, 10)),
											ast.NewPublicIdentifierNode(S(P(20, 5, 1), P(22, 5, 3)), "baz"),
											"to_int",
										),
									),
									ast.NewPublicIdentifierNode(S(P(34, 6, 1), P(36, 6, 3)), "baz"),
								),
							},
						),
					),
				},
			),
		},
		"can span multiple lines with a trailing comma": {
			input: "%{\nfoo:\nbar,\nbaz =>\nbaz.to_int if\nbaz,\n}",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(39, 7, 1)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(39, 7, 1)),
						ast.NewHashRecordLiteralNode(
							S(P(0, 1, 1), P(39, 7, 1)),
							[]ast.ExpressionNode{
								ast.NewSymbolKeyValueExpressionNode(
									S(P(3, 2, 1), P(10, 3, 3)),
									"foo",
									ast.NewPublicIdentifierNode(S(P(8, 3, 1), P(10, 3, 3)), "bar"),
								),
								ast.NewModifierNode(
									S(P(13, 4, 1), P(36, 6, 3)),
									T(S(P(31, 5, 12), P(32, 5, 13)), token.IF),
									ast.NewKeyValueExpressionNode(
										S(P(13, 4, 1), P(29, 5, 10)),
										ast.NewPublicIdentifierNode(S(P(13, 4, 1), P(15, 4, 3)), "baz"),
										ast.NewAttributeAccessNode(
											S(P(20, 5, 1), P(29, 5, 10)),
											ast.NewPublicIdentifierNode(S(P(20, 5, 1), P(22, 5, 3)), "baz"),
											"to_int",
										),
									),
									ast.NewPublicIdentifierNode(S(P(34, 6, 1), P(36, 6, 3)), "baz"),
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

func TestRegexLiteral(t *testing.T) {
	tests := testTable{
		"can be empty": {
			input: "%//",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						ast.NewUninterpolatedRegexLiteralNode(
							S(P(0, 1, 1), P(2, 1, 3)),
							"",
							bitfield.BitField8{},
						),
					),
				},
			),
		},
		"can be nested in string interpolation": {
			input: `"foo: ${%/bar\w+/i}"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewInterpolatedStringLiteralNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(S(P(1, 1, 2), P(5, 1, 6)), "foo: "),
								ast.NewStringInterpolationNode(
									S(P(6, 1, 7), P(18, 1, 19)),
									ast.NewUninterpolatedRegexLiteralNode(
										S(P(8, 1, 9), P(17, 1, 18)),
										`bar\w+`,
										bitfield.BitField8FromBitFlag(flag.CaseInsensitiveFlag),
									),
								),
							},
						),
					),
				},
			),
		},
		"can be empty with flags": {
			input: "%//im",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewUninterpolatedRegexLiteralNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							"",
							bitfield.BitField8FromBitFlag(flag.CaseInsensitiveFlag|flag.MultilineFlag),
						),
					),
				},
			),
		},
		"cannot have invalid flags": {
			input: "%//ipm",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewUninterpolatedRegexLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							"",
							bitfield.BitField8FromBitFlag(flag.CaseInsensitiveFlag),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(4, 1, 5), P(4, 1, 5)), "invalid regex flag"),
			},
		},
		"can have content": {
			input: `%/foo\/\w+bar/`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewUninterpolatedRegexLiteralNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							`foo\/\w+bar`,
							bitfield.BitField8{},
						),
					),
				},
			),
		},
		"can be interpolated": {
			input: `%/foo${oompa + loompa}\w+bar/`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(28, 1, 29)),
						ast.NewInterpolatedRegexLiteralNode(
							S(P(0, 1, 1), P(28, 1, 29)),
							[]ast.RegexLiteralContentNode{
								ast.NewRegexLiteralContentSectionNode(
									S(P(2, 1, 3), P(4, 1, 5)),
									"foo",
								),
								ast.NewRegexInterpolationNode(
									S(P(5, 1, 6), P(21, 1, 22)),
									ast.NewBinaryExpressionNode(
										S(P(7, 1, 8), P(20, 1, 21)),
										T(S(P(13, 1, 14), P(13, 1, 14)), token.PLUS),
										ast.NewPublicIdentifierNode(
											S(P(7, 1, 8), P(11, 1, 12)),
											"oompa",
										),
										ast.NewPublicIdentifierNode(
											S(P(15, 1, 16), P(20, 1, 21)),
											"loompa",
										),
									),
								),
								ast.NewRegexLiteralContentSectionNode(
									S(P(22, 1, 23), P(27, 1, 28)),
									`\w+bar`,
								),
							},
							bitfield.BitField8{},
						),
					),
				},
			),
		},
		"can have content and flags": {
			input: `%/foo\/bar/xUs`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewUninterpolatedRegexLiteralNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							`foo\/bar`,
							bitfield.BitField8FromBitFlag(flag.ExtendedFlag|flag.UngreedyFlag|flag.DotAllFlag),
						),
					),
				},
			),
		},
		"can repeat flags": {
			input: `%/foo\/bar/xUsxxxxss`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewUninterpolatedRegexLiteralNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							`foo\/bar`,
							bitfield.BitField8FromBitFlag(flag.ExtendedFlag|flag.UngreedyFlag|flag.DotAllFlag),
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
		"can be beginless and closed": {
			input: "...5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							T(S(P(0, 1, 1), P(2, 1, 3)), token.CLOSED_RANGE_OP),
							nil,
							ast.NewIntLiteralNode(S(P(3, 1, 4), P(3, 1, 4)), "5"),
						),
					),
				},
			),
		},
		"can be beginless and right open": {
			input: "..<5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							T(S(P(0, 1, 1), P(2, 1, 3)), token.RIGHT_OPEN_RANGE_OP),
							nil,
							ast.NewIntLiteralNode(S(P(3, 1, 4), P(3, 1, 4)), "5"),
						),
					),
				},
			),
		},
		"can be beginless and left open": {
			input: "<..5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							T(S(P(0, 1, 1), P(2, 1, 3)), token.LEFT_OPEN_RANGE_OP),
							nil,
							ast.NewIntLiteralNode(S(P(3, 1, 4), P(3, 1, 4)), "5"),
						),
					),
				},
			),
		},
		"can be beginless and open": {
			input: "<.<5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							T(S(P(0, 1, 1), P(2, 1, 3)), token.OPEN_RANGE_OP),
							nil,
							ast.NewIntLiteralNode(S(P(3, 1, 4), P(3, 1, 4)), "5"),
						),
					),
				},
			),
		},
		"can be endless and closed": {
			input: "5...",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							T(S(P(1, 1, 2), P(3, 1, 4)), token.CLOSED_RANGE_OP),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "5"),
							nil,
						),
					),
				},
			),
		},
		"can be endless and left open": {
			input: "5<..",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							T(S(P(1, 1, 2), P(3, 1, 4)), token.LEFT_OPEN_RANGE_OP),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "5"),
							nil,
						),
					),
				},
			),
		},
		"can be endless and right open": {
			input: "5..<",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							T(S(P(1, 1, 2), P(3, 1, 4)), token.RIGHT_OPEN_RANGE_OP),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "5"),
							nil,
						),
					),
				},
			),
		},
		"can be endless and open": {
			input: "5<.<",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(3, 1, 4)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(3, 1, 4)),
							T(S(P(1, 1, 2), P(3, 1, 4)), token.OPEN_RANGE_OP),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "5"),
							nil,
						),
					),
				},
			),
		},
		"can have a beginning and be closed": {
			input: "2...5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							T(S(P(1, 1, 2), P(3, 1, 4)), token.CLOSED_RANGE_OP),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "2"),
							ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "5"),
						),
					),
				},
			),
		},
		"can have a beginning and be right open": {
			input: "2..<5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							T(S(P(1, 1, 2), P(3, 1, 4)), token.RIGHT_OPEN_RANGE_OP),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "2"),
							ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "5"),
						),
					),
				},
			),
		},
		"can have a beginning and be left open": {
			input: "2<..5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							T(S(P(1, 1, 2), P(3, 1, 4)), token.LEFT_OPEN_RANGE_OP),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "2"),
							ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "5"),
						),
					),
				},
			),
		},
		"can have a beginning and be open": {
			input: "2<.<5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewRangeLiteralNode(
							S(P(0, 1, 1), P(4, 1, 5)),
							T(S(P(1, 1, 2), P(3, 1, 4)), token.OPEN_RANGE_OP),
							ast.NewIntLiteralNode(S(P(0, 1, 1), P(0, 1, 1)), "2"),
							ast.NewIntLiteralNode(S(P(4, 1, 5), P(4, 1, 5)), "5"),
						),
					),
				},
			),
		},
		"can have any expressions as operands": {
			input: "(2 * 5)...'foo'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(1, 1, 2), P(14, 1, 15)),
						ast.NewRangeLiteralNode(
							S(P(1, 1, 2), P(14, 1, 15)),
							T(S(P(7, 1, 8), P(9, 1, 10)), token.CLOSED_RANGE_OP),
							ast.NewBinaryExpressionNode(
								S(P(1, 1, 2), P(5, 1, 6)),
								T(S(P(3, 1, 4), P(3, 1, 4)), token.STAR),
								ast.NewIntLiteralNode(S(P(1, 1, 2), P(1, 1, 2)), "2"),
								ast.NewIntLiteralNode(S(P(5, 1, 6), P(5, 1, 6)), "5"),
							),
							ast.NewRawStringLiteralNode(
								S(P(10, 1, 11), P(14, 1, 15)),
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
