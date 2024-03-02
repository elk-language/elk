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
		"dot": {
			input: `\.`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'.',
			),
		},
		"question mark": {
			input: `\?`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'?',
			),
		},
		"dash": {
			input: `\-`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'-',
			),
		},
		"plus": {
			input: `\+`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'+',
			),
		},
		"star": {
			input: `\*`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'*',
			),
		},
		"caret": {
			input: `\^`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'^',
			),
		},
		"backslash": {
			input: `\\`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'\\',
			),
		},
		"pipe": {
			input: `\|`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'|',
			),
		},
		"dollar": {
			input: `\$`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'$',
			),
		},
		"left paren": {
			input: `\(`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'(',
			),
		},
		"right paren": {
			input: `\)`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				')',
			),
		},
		"left bracket": {
			input: `\[`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'[',
			),
		},
		"right bracket": {
			input: `\]`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				']',
			),
		},
		"left brace": {
			input: `\{`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'{',
			),
		},
		"right brace": {
			input: `\}`,
			want: ast.NewMetaCharEscapeNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				'}',
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestUnicodeCharClass(t *testing.T) {
	tests := testTable{
		"one letter": {
			input: `\pL`,
			want: ast.NewUnicodeCharClassNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				"L",
			),
		},
		"multi-letter": {
			input: `\p{Latin}`,
			want: ast.NewUnicodeCharClassNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				"Latin",
			),
		},
		"negated": {
			input: `\p{^Latin}`,
			want: ast.NewNegatedUnicodeCharClassNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				"Latin",
			),
		},
		"invalid multi-letter": {
			input: `\p{Latin9}`,
			want: ast.NewUnicodeCharClassNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				"Latin9",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(8, 1, 9), P(8, 1, 9)), "unexpected 9, expected an alphabetic character"),
			},
		},
		"missing end brace": {
			input: `\p{Latin`,
			want: ast.NewUnicodeCharClassNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				"Latin",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(8, 1, 9), P(7, 1, 8)), "unexpected END_OF_FILE, expected an alphabetic character"),
			},
		},
		"invalid single char": {
			input: `\p'`,
			want: ast.NewUnicodeCharClassNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				"'",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(2, 1, 3), P(2, 1, 3)), "unexpected ', expected an alphabetic character"),
			},
		},
		"missing single char": {
			input: `\p`,
			want: ast.NewUnicodeCharClassNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				"E",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(2, 1, 3), P(1, 1, 2)), "unexpected END_OF_FILE, expected CHAR"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestNegatedUnicodeCharClass(t *testing.T) {
	tests := testTable{
		"one letter": {
			input: `\PL`,
			want: ast.NewNegatedUnicodeCharClassNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				"L",
			),
		},
		"multi-letter": {
			input: `\P{Latin}`,
			want: ast.NewNegatedUnicodeCharClassNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				"Latin",
			),
		},
		"negated": {
			input: `\P{^Latin}`,
			want: ast.NewUnicodeCharClassNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				"Latin",
			),
		},
		"invalid multi-letter": {
			input: `\P{Latin9}`,
			want: ast.NewNegatedUnicodeCharClassNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				"Latin9",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(8, 1, 9), P(8, 1, 9)), "unexpected 9, expected an alphabetic character"),
			},
		},
		"missing end brace": {
			input: `\P{Latin`,
			want: ast.NewNegatedUnicodeCharClassNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				"Latin",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(8, 1, 9), P(7, 1, 8)), "unexpected END_OF_FILE, expected an alphabetic character"),
			},
		},
		"invalid single char": {
			input: `\P'`,
			want: ast.NewNegatedUnicodeCharClassNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				"'",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(2, 1, 3), P(2, 1, 3)), "unexpected ', expected an alphabetic character"),
			},
		},
		"missing single char": {
			input: `\P`,
			want: ast.NewNegatedUnicodeCharClassNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				"E",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(2, 1, 3), P(1, 1, 2)), "unexpected END_OF_FILE, expected CHAR"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestQuantifier(t *testing.T) {
	tests := testTable{
		"zero or one quantifier on char": {
			input: `p?`,
			want: ast.NewZeroOrOneQuantifierNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				false,
			),
		},
		"zero or one quantifier on char class": {
			input: `\w?`,
			want: ast.NewZeroOrOneQuantifierNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				ast.NewWordCharClassNode(
					S(P(0, 1, 1), P(1, 1, 2)),
				),
				false,
			),
		},
		"zero or one quantifier on group": {
			input: `(a\w)?`,
			want: ast.NewZeroOrOneQuantifierNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				ast.NewGroupNode(
					S(P(0, 1, 1), P(4, 1, 5)),
					ast.NewConcatenationNode(
						S(P(1, 1, 2), P(3, 1, 4)),
						[]ast.ConcatenationElementNode{
							ast.NewCharNode(
								S(P(1, 1, 2), P(1, 1, 2)),
								'a',
							),
							ast.NewWordCharClassNode(
								S(P(2, 1, 3), P(3, 1, 4)),
							),
						},
					),
					"",
					"",
					false,
				),
				false,
			),
		},
		"applies to only a single preceding item": {
			input: `ep\w?`,
			want: ast.NewConcatenationNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				[]ast.ConcatenationElementNode{
					ast.NewCharNode(
						S(P(0, 1, 1), P(0, 1, 1)),
						'e',
					),
					ast.NewCharNode(
						S(P(1, 1, 2), P(1, 1, 2)),
						'p',
					),
					ast.NewZeroOrOneQuantifierNode(
						S(P(2, 1, 3), P(4, 1, 5)),
						ast.NewWordCharClassNode(
							S(P(2, 1, 3), P(3, 1, 4)),
						),
						false,
					),
				},
			),
		},
		"zero or one alt quantifier": {
			input: `p??`,
			want: ast.NewZeroOrOneQuantifierNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				true,
			),
		},
		"zero or more quantifier": {
			input: `p*`,
			want: ast.NewZeroOrMoreQuantifierNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				false,
			),
		},
		"zero or more alt quantifier": {
			input: `p*?`,
			want: ast.NewZeroOrMoreQuantifierNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				true,
			),
		},
		"one or more quantifier": {
			input: `p+`,
			want: ast.NewOneOrMoreQuantifierNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				false,
			),
		},
		"one or more alt quantifier": {
			input: `p+?`,
			want: ast.NewOneOrMoreQuantifierNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				true,
			),
		},
		"N quantifier one digit": {
			input: `p{5}`,
			want: ast.NewNQuantifierNode(
				S(P(0, 1, 1), P(3, 1, 4)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"5",
			),
		},
		"N quantifier alt": {
			input: `p{5}?`,
			want: ast.NewNQuantifierNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"5",
			),
		},
		"N quantifier multiple digits": {
			input: `p{164}`,
			want: ast.NewNQuantifierNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"164",
			),
		},
		"N quantifier invalid chars": {
			input: `p{5f+9}`,
			want: ast.NewNQuantifierNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"5f9",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(3, 1, 4), P(3, 1, 4)), "unexpected f, expected a decimal digit"),
				errors.NewError(L("regex", P(4, 1, 5), P(4, 1, 5)), "unexpected +, expected a decimal digit"),
			},
		},
		"N quantifier missing right brace": {
			input: `p{5`,
			want: ast.NewNQuantifierNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"5",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(3, 1, 4), P(2, 1, 3)), "unexpected END_OF_FILE, expected }"),
			},
		},
		"N quantifier missing digit": {
			input: `p{}`,
			want: ast.NewNQuantifierNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(2, 1, 3), P(2, 1, 3)), "expected decimal digits"),
			},
		},
		"N quantifier missing digit and right brace": {
			input: `p{`,
			want: ast.NewNQuantifierNode(
				S(P(0, 1, 1), P(1, 1, 2)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(2, 1, 3), P(1, 1, 2)), "unexpected END_OF_FILE, expected }"),
				errors.NewError(L("regex", P(1, 1, 2), P(1, 1, 2)), "expected decimal digits"),
			},
		},
		"NM quantifier N only": {
			input: `p{5,}`,
			want: ast.NewNMQuantifierNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"5",
				"",
				false,
			),
		},
		"NM quantifier N only missing right brace": {
			input: `p{5,`,
			want: ast.NewNMQuantifierNode(
				S(P(0, 1, 1), P(0, 1, 1)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"5",
				"",
				false,
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(4, 1, 5), P(3, 1, 4)), "unexpected END_OF_FILE, expected }"),
			},
		},
		"NM quantifier N only alt": {
			input: `p{58,}?`,
			want: ast.NewNMQuantifierNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"58",
				"",
				true,
			),
		},
		"NM quantifier": {
			input: `p{58,153}`,
			want: ast.NewNMQuantifierNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"58",
				"153",
				false,
			),
		},
		"NM quantifier only M": {
			input: `p{,153}`,
			want: ast.NewNMQuantifierNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"",
				"153",
				false,
			),
		},
		"NM quantifier only M alt": {
			input: `p{,153}?`,
			want: ast.NewNMQuantifierNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"",
				"153",
				true,
			),
		},
		"NM quantifier alt": {
			input: `p{58,153}?`,
			want: ast.NewNMQuantifierNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"58",
				"153",
				true,
			),
		},
		"NM quantifier missing right brace": {
			input: `p{58,153`,
			want: ast.NewNMQuantifierNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"58",
				"153",
				false,
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(8, 1, 9), P(7, 1, 8)), "unexpected END_OF_FILE, expected }"),
			},
		},
		"NM quantifier invalid chars": {
			input: `p{a8,1f3}`,
			want: ast.NewNMQuantifierNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				ast.NewCharNode(
					S(P(0, 1, 1), P(0, 1, 1)),
					'p',
				),
				"a8",
				"1f3",
				false,
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(2, 1, 3), P(2, 1, 3)), "unexpected a, expected a decimal digit"),
				errors.NewError(L("regex", P(6, 1, 7), P(6, 1, 7)), "unexpected f, expected a decimal digit"),
			},
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
		"missing end brace": {
			input: `\x{6f`,
			want: ast.NewHexEscapeNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				"6f",
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(5, 1, 6), P(4, 1, 5)), "unexpected END_OF_FILE, expected a hex digit"),
			},
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
				[]ast.ConcatenationElementNode{
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
				[]ast.ConcatenationElementNode{
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
				[]ast.ConcatenationElementNode{
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
		"chars escapes, anchors and groups": {
			input: `(f\n)\w$`,
			want: ast.NewConcatenationNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.ConcatenationElementNode{
					ast.NewGroupNode(
						S(P(0, 1, 1), P(4, 1, 5)),
						ast.NewConcatenationNode(
							S(P(1, 1, 2), P(3, 1, 4)),
							[]ast.ConcatenationElementNode{
								ast.NewCharNode(
									S(P(1, 1, 2), P(1, 1, 2)),
									'f',
								),
								ast.NewNewlineEscapeNode(
									S(P(2, 1, 3), P(3, 1, 4)),
								),
							},
						),
						"",
						"",
						false,
					),
					ast.NewWordCharClassNode(
						S(P(5, 1, 6), P(6, 1, 7)),
					),
					ast.NewEndOfStringAnchorNode(
						S(P(7, 1, 8), P(7, 1, 8)),
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
					[]ast.ConcatenationElementNode{
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
					[]ast.ConcatenationElementNode{
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
		"group union": {
			input: "(foo)|barę",
			want: ast.NewUnionNode(
				S(P(0, 1, 1), P(10, 1, 10)),
				ast.NewGroupNode(
					S(P(0, 1, 1), P(4, 1, 5)),
					ast.NewConcatenationNode(
						S(P(1, 1, 2), P(3, 1, 4)),
						[]ast.ConcatenationElementNode{
							ast.NewCharNode(
								S(P(1, 1, 2), P(1, 1, 2)),
								'f',
							),
							ast.NewCharNode(
								S(P(2, 1, 3), P(2, 1, 3)),
								'o',
							),
							ast.NewCharNode(
								S(P(3, 1, 4), P(3, 1, 4)),
								'o',
							),
						},
					),
					"",
					"",
					false,
				),
				ast.NewConcatenationNode(
					S(P(6, 1, 7), P(10, 1, 10)),
					[]ast.ConcatenationElementNode{
						ast.NewCharNode(
							S(P(6, 1, 7), P(6, 1, 7)),
							'b',
						),
						ast.NewCharNode(
							S(P(7, 1, 8), P(7, 1, 8)),
							'a',
						),
						ast.NewCharNode(
							S(P(8, 1, 9), P(8, 1, 9)),
							'r',
						),
						ast.NewCharNode(
							S(P(9, 1, 10), P(10, 1, 10)),
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
						[]ast.ConcatenationElementNode{
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

func TestGroup(t *testing.T) {
	tests := testTable{
		"non capturing group": {
			input: "(?:f)",
			want: ast.NewGroupNode(
				S(P(0, 1, 1), P(4, 1, 5)),
				ast.NewCharNode(
					S(P(3, 1, 4), P(3, 1, 4)),
					'f',
				),
				"",
				"",
				true,
			),
		},
		"named group": {
			input: "(?<foo>f)",
			want: ast.NewGroupNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				ast.NewCharNode(
					S(P(7, 1, 8), P(7, 1, 8)),
					'f',
				),
				"foo",
				"",
				false,
			),
		},
		"named group with P": {
			input: "(?P<foo>f)",
			want: ast.NewGroupNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				ast.NewCharNode(
					S(P(8, 1, 9), P(8, 1, 9)),
					'f',
				),
				"foo",
				"",
				false,
			),
		},
		"flags only": {
			input: "(?imU)",
			want: ast.NewGroupNode(
				S(P(0, 1, 1), P(5, 1, 6)),
				nil,
				"",
				"imU",
				false,
			),
		},
		"flags and content": {
			input: "(?mi-s:f)",
			want: ast.NewGroupNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				ast.NewCharNode(
					S(P(7, 1, 8), P(7, 1, 8)),
					'f',
				),
				"",
				"mi-s",
				false,
			),
		},
		"invalid flags": {
			input: "(?mihs:f)",
			want: ast.NewGroupNode(
				S(P(0, 1, 1), P(8, 1, 9)),
				ast.NewCharNode(
					S(P(7, 1, 8), P(7, 1, 8)),
					'f',
				),
				"",
				"mihs",
				false,
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(4, 1, 5), P(4, 1, 5)), "unexpected h, expected a regex flag"),
			},
		},
		"char in group": {
			input: "(f)",
			want: ast.NewGroupNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				ast.NewCharNode(
					S(P(1, 1, 2), P(1, 1, 2)),
					'f',
				),
				"",
				"",
				false,
			),
		},
		"missing right paren": {
			input: "(f",
			want: ast.NewInvalidNode(
				S(P(2, 1, 3), P(1, 1, 2)),
				T(S(P(2, 1, 3), P(1, 1, 2)), token.END_OF_FILE),
			),
			err: errors.ErrorList{
				errors.NewError(L("regex", P(2, 1, 3), P(1, 1, 2)), "unexpected END_OF_FILE, expected )"),
			},
		},
		"union in group": {
			input: "(foo|barę)",
			want: ast.NewGroupNode(
				S(P(0, 1, 1), P(10, 1, 10)),
				ast.NewUnionNode(
					S(P(1, 1, 2), P(9, 1, 9)),
					ast.NewConcatenationNode(
						S(P(1, 1, 2), P(3, 1, 4)),
						[]ast.ConcatenationElementNode{
							ast.NewCharNode(
								S(P(1, 1, 2), P(1, 1, 2)),
								'f',
							),
							ast.NewCharNode(
								S(P(2, 1, 3), P(2, 1, 3)),
								'o',
							),
							ast.NewCharNode(
								S(P(3, 1, 4), P(3, 1, 4)),
								'o',
							),
						},
					),
					ast.NewConcatenationNode(
						S(P(5, 1, 6), P(9, 1, 9)),
						[]ast.ConcatenationElementNode{
							ast.NewCharNode(
								S(P(5, 1, 6), P(5, 1, 6)),
								'b',
							),
							ast.NewCharNode(
								S(P(6, 1, 7), P(6, 1, 7)),
								'a',
							),
							ast.NewCharNode(
								S(P(7, 1, 8), P(7, 1, 8)),
								'r',
							),
							ast.NewCharNode(
								S(P(8, 1, 9), P(9, 1, 9)),
								'ę',
							),
						},
					),
				),
				"",
				"",
				false,
			),
		},
		"nested groups": {
			input: "((foo)|barę)",
			want: ast.NewGroupNode(
				S(P(0, 1, 1), P(12, 1, 12)),
				ast.NewUnionNode(
					S(P(1, 1, 2), P(11, 1, 11)),
					ast.NewGroupNode(
						S(P(1, 1, 2), P(5, 1, 6)),
						ast.NewConcatenationNode(
							S(P(2, 1, 3), P(4, 1, 5)),
							[]ast.ConcatenationElementNode{
								ast.NewCharNode(
									S(P(2, 1, 3), P(2, 1, 3)),
									'f',
								),
								ast.NewCharNode(
									S(P(3, 1, 4), P(3, 1, 4)),
									'o',
								),
								ast.NewCharNode(
									S(P(4, 1, 5), P(4, 1, 5)),
									'o',
								),
							},
						),
						"",
						"",
						false,
					),
					ast.NewConcatenationNode(
						S(P(7, 1, 8), P(11, 1, 11)),
						[]ast.ConcatenationElementNode{
							ast.NewCharNode(
								S(P(7, 1, 8), P(7, 1, 8)),
								'b',
							),
							ast.NewCharNode(
								S(P(8, 1, 9), P(8, 1, 9)),
								'a',
							),
							ast.NewCharNode(
								S(P(9, 1, 10), P(9, 1, 10)),
								'r',
							),
							ast.NewCharNode(
								S(P(10, 1, 11), P(11, 1, 11)),
								'ę',
							),
						},
					),
				),
				"",
				"",
				false,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}
