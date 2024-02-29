package parser

import (
	"testing"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/regex/parser/ast"
	"github.com/elk-language/elk/regex/token"
	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp"
)

// Represents a single parser test case.
type testCase struct {
	input string
	want  ast.Node
	err   errors.ErrorList
}

// Type of the parser test table.
type testTable map[string]testCase

// Create a new token in tests.
var T = token.New

// Create a new token with value in tests.
var V = token.NewWithValue

// Create a new source position in tests.
var P = position.New

// Create a new span in tests.
var S = position.NewSpan

// Create a new source location in tests.
var L = position.NewLocation

// Function which powers all parser tests.
// Inspects if the produced AST matches the expected one.
func parserTest(tc testCase, t *testing.T) {
	t.Helper()
	got, err := Parse(tc.input)

	opts := []cmp.Option{
		cmp.AllowUnexported(
			ast.NodeBase{},
			token.Token{},
		),
	}
	if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
		pp.Println(got)
		t.Fatal(diff)
	}

	if diff := cmp.Diff(tc.err, err); diff != "" {
		t.Fatal(diff)
	}
}

func TestChar(t *testing.T) {
	tests := testTable{
		"ascii char": {
			input: "f",
			want: ast.NewCharNode(
				S(P(0, 1, 1), P(0, 1, 1)),
				'f',
			),
		},
		"two byte char": {
			input: "ę",
			want: ast.NewCharNode(
				S(P(0, 1, 1), P(1, 1, 1)),
				'ę',
			),
		},
		"three byte char": {
			input: "€",
			want: ast.NewCharNode(
				S(P(0, 1, 1), P(2, 1, 1)),
				'€',
			),
		},
		"four byte char": {
			input: "𐍈",
			want: ast.NewCharNode(
				S(P(0, 1, 1), P(3, 1, 1)),
				'𐍈',
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestSimpleEscape(t *testing.T) {
	tests := testTable{
		"bell": {
			input: `\a`,
			want: ast.NewBellEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"form feed": {
			input: `\f`,
			want: ast.NewFormFeedEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"tab": {
			input: `\t`,
			want: ast.NewTabEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"newline": {
			input: `\n`,
			want: ast.NewNewlineEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"carriage return": {
			input: `\r`,
			want: ast.NewCarriageReturnEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"vertical tab": {
			input: `\v`,
			want: ast.NewVerticalTabEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestHexEscape(t *testing.T) {
	tests := testTable{
		"two digit": {
			input: `\x6f`,
			want: ast.NewHexEscapeNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				"6f",
			),
		},
		"two digit with invalid char": {
			input: `\x6l`,
			want: ast.NewHexEscapeNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				"6l",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(3, 1, 4), P(3, 1, 4)), "unexpected l, expected a hex digit"),
			},
		},
		"two digit with invalid meta char": {
			input: `\x6{`,
			want: ast.NewHexEscapeNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				"6",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(3, 1, 4), P(3, 1, 4)), "unexpected {, expected a hex digit"),
			},
		},
		"missing digit": {
			input: `\xf`,
			want: ast.NewHexEscapeNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				"f",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(3, 1, 4), P(2, 1, 3)), "unexpected END_OF_FILE, expected a hex digit"),
			},
		},
		"with braces": {
			input: `\x{6f}`,
			want: ast.NewHexEscapeNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				"6f",
			),
		},
		"long with braces": {
			input: `\x{6f10}`,
			want: ast.NewHexEscapeNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				"6f10",
			),
		},
		"with braces and invalid chars": {
			input: `\x{6.f{0}`,
			want: ast.NewHexEscapeNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				"6f0",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(4, 1, 5), P(4, 1, 5)), "unexpected ., expected a hex digit"),
				errors.NewError(L("regex", P(6, 1, 7), P(6, 1, 7)), "unexpected {, expected a hex digit"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestAnchor(t *testing.T) {
	tests := testTable{
		"absolute start of string": {
			input: `\A`,
			want: ast.NewAbsoluteStartOfStringAnchorNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"absolute end of string": {
			input: `\z`,
			want: ast.NewAbsoluteEndOfStringAnchorNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"start of string": {
			input: `^`,
			want: ast.NewStartOfStringAnchorNode(
				S(P(0, 1, 1), P(0, 1, 1)),
			),
		},
		"end of string": {
			input: `$`,
			want: ast.NewEndOfStringAnchorNode(
				S(P(0, 1, 1), P(0, 1, 1)),
			),
		},
		"word boundary": {
			input: `\b`,
			want: ast.NewWordBoundaryAnchorNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"not word boundary": {
			input: `\B`,
			want: ast.NewNotWordBoundaryAnchorNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestCharClass(t *testing.T) {
	tests := testTable{
		"word": {
			input: `\w`,
			want: ast.NewWordCharClassNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"not word": {
			input: `\W`,
			want: ast.NewNotWordCharClassNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"digit": {
			input: `\d`,
			want: ast.NewDigitCharClassNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"not digit": {
			input: `\D`,
			want: ast.NewNotDigitCharClassNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"whitespace": {
			input: `\s`,
			want: ast.NewWhitespaceCharClassNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
		"not whitespace": {
			input: `\S`,
			want: ast.NewNotWhitespaceCharClassNode(
				S(P(0, 1, 1), P(1, 1, 2)),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestConcatenation(t *testing.T) {
	tests := testTable{
		"ascii chars": {
			input: "foo",
			want: ast.NewConcatenationNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.PrimaryRegexNode{
					ast.NewCharNode(
						S(P(0, 1, 1), P(0, 1, 1)),
						'f',
					),
					ast.NewCharNode(
						S(P(1, 1, 2), P(1, 1, 2)),
						'o',
					),
					ast.NewCharNode(
						S(P(2, 1, 3), P(2, 1, 3)),
						'o',
					),
				},
			),
		},
		"multi-byte chars": {
			input: "fęłó€𐍈",
			want: ast.NewConcatenationNode(
				S(P(0, 1, 1), P(13, 1, 6)),
				[]ast.PrimaryRegexNode{
					ast.NewCharNode(
						S(P(0, 1, 1), P(0, 1, 1)),
						'f',
					),
					ast.NewCharNode(
						S(P(1, 1, 2), P(2, 1, 2)),
						'ę',
					),
					ast.NewCharNode(
						S(P(3, 1, 3), P(4, 1, 3)),
						'ł',
					),
					ast.NewCharNode(
						S(P(5, 1, 4), P(6, 1, 4)),
						'ó',
					),
					ast.NewCharNode(
						S(P(7, 1, 5), P(9, 1, 5)),
						'€',
					),
					ast.NewCharNode(
						S(P(10, 1, 6), P(13, 1, 6)),
						'𐍈',
					),
				},
			),
		},
		"chars escapes and anchors": {
			input: `f\n\w$`,
			want: ast.NewConcatenationNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				[]ast.PrimaryRegexNode{
					ast.NewCharNode(
						S(P(0, 1, 1), P(0, 1, 1)),
						'f',
					),
					ast.NewNewlineEscapeNode(
						S(P(1, 1, 2), P(2, 1, 3)),
					),
					ast.NewWordCharClassNode(
						S(P(3, 1, 4), P(4, 1, 5)),
					),
					ast.NewEndOfStringAnchorNode(
						S(P(5, 1, 6), P(5, 1, 6)),
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

func TestUnion(t *testing.T) {
	tests := testTable{
		"char union": {
			input: "f|o",
			want: ast.NewUnionNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'f',
				),
				ast.NewCharNode(
					S(P(2, 1, 3), P(2, 1, 3)),
					'o',
				),
			),
		},
		"concat union": {
			input: "foo|barę",
			want: ast.NewUnionNode(
				S(P(0, 1, 1), P(8, 1, 8)),
				ast.NewConcatenationNode(
					S(P(0, 1, 1), P(2, 1, 3)),
					[]ast.PrimaryRegexNode{
						ast.NewCharNode(
							S(P(0, 1, 1), P(0, 1, 1)),
							'f',
						),
						ast.NewCharNode(
							S(P(1, 1, 2), P(1, 1, 2)),
							'o',
						),
						ast.NewCharNode(
							S(P(2, 1, 3), P(2, 1, 3)),
							'o',
						),
					},
				),
				ast.NewConcatenationNode(
					S(P(4, 1, 5), P(8, 1, 8)),
					[]ast.PrimaryRegexNode{
						ast.NewCharNode(
							S(P(4, 1, 5), P(4, 1, 5)),
							'b',
						),
						ast.NewCharNode(
							S(P(5, 1, 6), P(5, 1, 6)),
							'a',
						),
						ast.NewCharNode(
							S(P(6, 1, 7), P(6, 1, 7)),
							'r',
						),
						ast.NewCharNode(
							S(P(7, 1, 8), P(8, 1, 8)),
							'ę',
						),
					},
				),
			),
		},
		"nested unions": {
			input: "foo|b|u",
			want: ast.NewUnionNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				ast.NewUnionNode(
					S(P(0, 1, 1), P(4, 1, 5)),
					ast.NewConcatenationNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						[]ast.PrimaryRegexNode{
							ast.NewCharNode(
								S(P(0, 1, 1), P(0, 1, 1)),
								'f',
							),
							ast.NewCharNode(
								S(P(1, 1, 2), P(1, 1, 2)),
								'o',
							),
							ast.NewCharNode(
								S(P(2, 1, 3), P(2, 1, 3)),
								'o',
							),
						},
					),
					ast.NewCharNode(
						S(P(4, 1, 5), P(4, 1, 5)),
						'b',
					),
				),
				ast.NewCharNode(
					S(P(6, 1, 7), P(6, 1, 7)),
					'u',
				),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}
