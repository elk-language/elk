package regex

import (
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/regex/token"
	"github.com/google/go-cmp/cmp"
)

// Represents a single parser test case.
type testCase struct {
	input string
	want  string
	err   error.ErrorList
	flags bitfield.BitField8
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

// Function which powers all transpiler tests.
// Inspects if the produced string matches the expected one.
func transpilerTest(tc testCase, t *testing.T) {
	t.Helper()
	got, err := Transpile(tc.input, tc.flags)

	if diff := cmp.Diff(tc.err, err); diff != "" {
		t.Fatal(diff)
	}

	if diff := cmp.Diff(tc.want, got); diff != "" {
		t.Fatal(diff)
	}

}

func TestChar(t *testing.T) {
	tests := testTable{
		"ascii char": {
			input: "f",
			want:  "f",
		},
		"two byte char": {
			input: "ę",
			want:  "ę",
		},
		"three byte char": {
			input: "€",
			want:  "€",
		},
		"four byte char": {
			input: "𐍈",
			want:  "𐍈",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestSimpleEscape(t *testing.T) {
	tests := testTable{
		"bell": {
			input: `\a`,
			want:  `\a`,
		},
		"form feed": {
			input: `\f`,
			want:  `\f`,
		},
		"tab": {
			input: `\t`,
			want:  `\t`,
		},
		"newline": {
			input: `\n`,
			want:  `\n`,
		},
		"carriage return": {
			input: `\r`,
			want:  `\r`,
		},
		"dot": {
			input: `\.`,
			want:  `\.`,
		},
		"question mark": {
			input: `\?`,
			want:  `\?`,
		},
		"dash": {
			input: `\-`,
			want:  `\-`,
		},
		"plus": {
			input: `\+`,
			want:  `\+`,
		},
		"star": {
			input: `\*`,
			want:  `\*`,
		},
		"caret": {
			input: `\^`,
			want:  `\^`,
		},
		"backslash": {
			input: `\\`,
			want:  `\\`,
		},
		"pipe": {
			input: `\|`,
			want:  `\|`,
		},
		"dollar": {
			input: `\$`,
			want:  `\$`,
		},
		"left paren": {
			input: `\(`,
			want:  `\(`,
		},
		"right paren": {
			input: `\)`,
			want:  `\)`,
		},
		"left bracket": {
			input: `\[`,
			want:  `\[`,
		},
		"right bracket": {
			input: `\]`,
			want:  `\]`,
		},
		"left brace": {
			input: `\{`,
			want:  `\{`,
		},
		"right brace": {
			input: `\}`,
			want:  `\}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestUnicodeCharClass(t *testing.T) {
	tests := testTable{
		"one letter": {
			input: `\pL`,
			want:  `\p{L}`,
		},
		"multi-letter": {
			input: `\p{Latin}`,
			want:  `\p{Latin}`,
		},
		"negated": {
			input: `\p{^Latin}`,
			want:  `\P{Latin}`,
		},
		"invalid multi-letter": {
			input: `\p{Latin9}`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(8, 1, 9), P(8, 1, 9)), "unexpected 9, expected an alphabetic character"),
			},
		},
		"missing end brace": {
			input: `\p{Latin`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(8, 1, 9), P(7, 1, 8)), "unexpected END_OF_FILE, expected an alphabetic character"),
			},
		},
		"invalid single char": {
			input: `\p'`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(2, 1, 3), P(2, 1, 3)), "unexpected ', expected an alphabetic character"),
				error.NewError(L("regex", P(2, 1, 3), P(2, 1, 3)), "unexpected ', expected an alphabetic character"),
			},
		},
		"missing single char": {
			input: `\p`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(2, 1, 3), P(1, 1, 2)), "unexpected END_OF_FILE, expected an alphabetic character"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestNegatedUnicodeCharClass(t *testing.T) {
	tests := testTable{
		"one letter": {
			input: `\PL`,
			want:  `\P{L}`,
		},
		"multi-letter": {
			input: `\P{Latin}`,
			want:  `\P{Latin}`,
		},
		"negated": {
			input: `\P{^Latin}`,
			want:  `\p{Latin}`,
		},
		"invalid multi-letter": {
			input: `\P{Latin9}`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(8, 1, 9), P(8, 1, 9)), "unexpected 9, expected an alphabetic character"),
			},
		},
		"missing end brace": {
			input: `\P{Latin`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(8, 1, 9), P(7, 1, 8)), "unexpected END_OF_FILE, expected an alphabetic character"),
			},
		},
		"invalid single char": {
			input: `\P'`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(2, 1, 3), P(2, 1, 3)), "unexpected ', expected an alphabetic character"),
				error.NewError(L("regex", P(2, 1, 3), P(2, 1, 3)), "unexpected ', expected an alphabetic character"),
			},
		},
		"missing single char": {
			input: `\P`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(2, 1, 3), P(1, 1, 2)), "unexpected END_OF_FILE, expected an alphabetic character"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestQuantifier(t *testing.T) {
	tests := testTable{
		"zero or one quantifier on char": {
			input: `p?`,
			want:  `p?`,
		},
		"zero or one quantifier on char class": {
			input: `\w?`,
			want:  `[\p{L}\p{Mn}\p{Nd}\p{Pc}]?`,
		},
		"zero or one quantifier on group": {
			input: `(a\w)?`,
			want:  `(a[\p{L}\p{Mn}\p{Nd}\p{Pc}])?`,
		},
		"applies to only a single preceding item": {
			input: `ep\w?`,
			want:  `ep[\p{L}\p{Mn}\p{Nd}\p{Pc}]?`,
		},
		"zero or one alt quantifier": {
			input: `p??`,
			want:  `p??`,
		},
		"zero or more quantifier": {
			input: `p*`,
			want:  `p*`,
		},
		"zero or more alt quantifier": {
			input: `p*?`,
			want:  `p*?`,
		},
		"one or more quantifier": {
			input: `p+`,
			want:  `p+`,
		},
		"one or more alt quantifier": {
			input: `p+?`,
			want:  `p+?`,
		},
		"N quantifier one digit": {
			input: `p{5}`,
			want:  `p{5}`,
		},
		"N quantifier alt": {
			input: `p{5}?`,
			want:  `p{5}?`,
		},
		"N quantifier multiple digits": {
			input: `p{164}`,
			want:  `p{164}`,
		},
		"N quantifier invalid chars": {
			input: `p{5f+9}`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(3, 1, 4), P(3, 1, 4)), "unexpected f, expected a decimal digit"),
				error.NewError(L("regex", P(4, 1, 5), P(4, 1, 5)), "unexpected +, expected a decimal digit"),
			},
		},
		"N quantifier missing right brace": {
			input: `p{5`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(3, 1, 4), P(2, 1, 3)), "unexpected END_OF_FILE, expected }"),
			},
		},
		"N quantifier missing digit": {
			input: `p{}`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(2, 1, 3), P(2, 1, 3)), "expected decimal digits"),
			},
		},
		"N quantifier missing digit and right brace": {
			input: `p{`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(2, 1, 3), P(1, 1, 2)), "unexpected END_OF_FILE, expected }"),
				error.NewError(L("regex", P(1, 1, 2), P(1, 1, 2)), "expected decimal digits"),
			},
		},
		"NM quantifier N only": {
			input: `p{5,}`,
			want:  `p{5,}`,
		},
		"NM quantifier N only missing right brace": {
			input: `p{5,`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(4, 1, 5), P(3, 1, 4)), "unexpected END_OF_FILE, expected }"),
			},
		},
		"NM quantifier N only alt": {
			input: `p{58,}?`,
			want:  `p{58,}?`,
		},
		"NM quantifier": {
			input: `p{58,153}`,
			want:  `p{58,153}`,
		},
		"NM quantifier only M": {
			input: `p{,153}`,
			want:  `p{0,153}`,
		},
		"NM quantifier only M alt": {
			input: `p{,153}?`,
			want:  `p{0,153}?`,
		},
		"NM quantifier alt": {
			input: `p{58,153}?`,
			want:  `p{58,153}?`,
		},
		"NM quantifier missing right brace": {
			input: `p{58,153`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(8, 1, 9), P(7, 1, 8)), "unexpected END_OF_FILE, expected }"),
			},
		},
		"NM quantifier invalid chars": {
			input: `p{a8,1f3}`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(2, 1, 3), P(2, 1, 3)), "unexpected a, expected a decimal digit"),
				error.NewError(L("regex", P(6, 1, 7), P(6, 1, 7)), "unexpected f, expected a decimal digit"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestCaretEscape(t *testing.T) {
	tests := testTable{
		"simple": {
			input: `\cA`,
			want:  `\x{1}`,
		},
		"consumes only a single letter": {
			input: `\czl`,
			want:  `\x{1a}l`,
		},
		"invalid char": {
			input: `\cę`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(2, 1, 3), P(3, 1, 3)), "unexpected ę, expected an ASCII letter"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestUnicodeEscape(t *testing.T) {
	tests := testTable{
		"four digit": {
			input: `\u6f45`,
			want:  `\x{6f45}`,
		},
		"four digit with invalid char": {
			input: `\u6f7l`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(5, 1, 6), P(5, 1, 6)), "unexpected l, expected a hex digit"),
			},
		},
		"four digit with invalid meta char": {
			input: `\u67f{`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(5, 1, 6), P(5, 1, 6)), "unexpected {, expected a hex digit"),
			},
		},
		"missing digit": {
			input: `\ufff`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(5, 1, 6), P(4, 1, 5)), "unexpected END_OF_FILE, expected a hex digit"),
			},
		},
		"with braces": {
			input: `\u{6f}`,
			want:  `\x{6f}`,
		},
		"missing end brace": {
			input: `\u{6f`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(5, 1, 6), P(4, 1, 5)), "unexpected END_OF_FILE, expected a hex digit"),
			},
		},
		"long with braces": {
			input: `\u{6f10}`,
			want:  `\x{6f10}`,
		},
		"with braces and invalid chars": {
			input: `\u{6.f{0}`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(4, 1, 5), P(4, 1, 5)), "unexpected ., expected a hex digit"),
				error.NewError(L("regex", P(6, 1, 7), P(6, 1, 7)), "unexpected {, expected a hex digit"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestLongUnicodeEscape(t *testing.T) {
	tests := testTable{
		"eight digit": {
			input: `\U00006f45`,
			want:  `\x{00006f45}`,
		},
		"eight digit with invalid char": {
			input: `\U00006f7l`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(9, 1, 10), P(9, 1, 10)), "unexpected l, expected a hex digit"),
			},
		},
		"eight digit with invalid meta char": {
			input: `\U000067f{`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(9, 1, 10), P(9, 1, 10)), "unexpected {, expected a hex digit"),
			},
		},
		"missing digit": {
			input: `\U0000fff`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(9, 1, 10), P(8, 1, 9)), "unexpected END_OF_FILE, expected a hex digit"),
			},
		},
		"with braces": {
			input: `\U{6f}`,
			want:  `\x{6f}`,
		},
		"missing end brace": {
			input: `\U{6f`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(5, 1, 6), P(4, 1, 5)), "unexpected END_OF_FILE, expected a hex digit"),
			},
		},
		"long with braces": {
			input: `\U{6f10}`,
			want:  `\x{6f10}`,
		},
		"with braces and invalid chars": {
			input: `\U{6.f{0}`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(4, 1, 5), P(4, 1, 5)), "unexpected ., expected a hex digit"),
				error.NewError(L("regex", P(6, 1, 7), P(6, 1, 7)), "unexpected {, expected a hex digit"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestHexEscape(t *testing.T) {
	tests := testTable{
		"two digit": {
			input: `\x6f`,
			want:  `\x{6f}`,
		},
		"two digit with invalid char": {
			input: `\x6l`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(3, 1, 4), P(3, 1, 4)), "unexpected l, expected a hex digit"),
			},
		},
		"two digit with invalid meta char": {
			input: `\x6{`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(3, 1, 4), P(3, 1, 4)), "unexpected {, expected a hex digit"),
			},
		},
		"missing digit": {
			input: `\xf`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(3, 1, 4), P(2, 1, 3)), "unexpected END_OF_FILE, expected a hex digit"),
			},
		},
		"with braces": {
			input: `\x{6f}`,
			want:  `\x{6f}`,
		},
		"missing end brace": {
			input: `\x{6f`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(5, 1, 6), P(4, 1, 5)), "unexpected END_OF_FILE, expected a hex digit"),
			},
		},
		"long with braces": {
			input: `\x{6f10}`,
			want:  `\x{6f10}`,
		},
		"with braces and invalid chars": {
			input: `\x{6.f{0}`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(4, 1, 5), P(4, 1, 5)), "unexpected ., expected a hex digit"),
				error.NewError(L("regex", P(6, 1, 7), P(6, 1, 7)), "unexpected {, expected a hex digit"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestOctalEscape(t *testing.T) {
	tests := testTable{
		"simple single digit": {
			input: `\1`,
			want:  `\001`,
		},
		"simple stops on last digit": {
			input: `\1f`,
			want:  `\001f`,
		},
		"simple two digits": {
			input: `\12`,
			want:  `\012`,
		},
		"simple three digits": {
			input: `\123`,
			want:  `\123`,
		},
		"simple too many digits": {
			input: `\1234`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(0, 1, 1), P(4, 1, 5)), `invalid octal escape: \1234`),
			},
		},
		"simple invalid digit": {
			input: `\182`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(0, 1, 1), P(3, 1, 4)), `invalid octal escape: \182`),
			},
		},

		"three digits": {
			input: `\o612`,
			want:  `\612`,
		},
		"two digit with invalid char": {
			input: `\o691`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(3, 1, 4), P(3, 1, 4)), "unexpected 9, expected an octal digit"),
			},
		},
		"two digit with invalid meta char": {
			input: `\o6{`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(3, 1, 4), P(3, 1, 4)), "unexpected {, expected an octal digit"),
				error.NewError(L("regex", P(4, 1, 5), P(3, 1, 4)), "unexpected END_OF_FILE, expected an octal digit"),
			},
		},
		"missing digit": {
			input: `\o72`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(4, 1, 5), P(3, 1, 4)), "unexpected END_OF_FILE, expected an octal digit"),
			},
		},
		"with braces": {
			input: `\o{62}`,
			want:  `\062`,
		},
		"missing end brace": {
			input: `\o{62`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(5, 1, 6), P(4, 1, 5)), "unexpected END_OF_FILE, expected an octal digit"),
			},
		},
		"long with braces": {
			input: `\o{612}`,
			want:  `\612`,
		},
		"with braces and too long": {
			input: `\o{6123}`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(0, 1, 1), P(6, 1, 7)), "too many octal digits"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestAnchor(t *testing.T) {
	tests := testTable{
		"absolute start of string": {
			input: `\A`,
			want:  `\A`,
		},
		"absolute end of string": {
			input: `\z`,
			want:  `\z`,
		},
		"start of string": {
			input: `^`,
			want:  `^`,
		},
		"end of string": {
			input: `$`,
			want:  `$`,
		},
		"word boundary": {
			input: `\b`,
			want:  `\b`,
		},
		"not word boundary": {
			input: `\B`,
			want:  `\B`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestSimpleCharClass(t *testing.T) {
	tests := testTable{
		"word": {
			input: `\w`,
			want:  `[\p{L}\p{Mn}\p{Nd}\p{Pc}]`,
		},
		"word in char class": {
			input: `[:,\w.]`,
			want:  `[:,\p{L}\p{Mn}\p{Nd}\p{Pc}.]`,
		},
		"word in negated char class": {
			input: `[^:,\w.]`,
			want:  `[^:,\p{L}\p{Mn}\p{Nd}\p{Pc}.]`,
		},
		"word ascii": {
			input: `\w`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `\w`,
		},
		"word in char class ascii": {
			input: `[:,\w.]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[:,\w.]`,
		},
		"word in negated char class ascii": {
			input: `[^:,\w.]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[^:,\w.]`,
		},
		"not word": {
			input: `\W`,
			want:  `[^\p{L}\p{Mn}\p{Nd}\p{Pc}]`,
		},
		"not word in char class": {
			input: `[ab\Wcd]`,
			want:  `(?:[abcd]|[^\p{L}\p{Mn}\p{Nd}\p{Pc}])`,
		},
		"not word in negated char class": {
			input: `[^ab\Wcd]`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(4, 1, 5), P(5, 1, 6)), `double negation of unicode-aware \W is not supported`),
			},
		},
		"not word ascii": {
			input: `\W`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `\W`,
		},
		"not word in char class ascii": {
			input: `[:.\W,]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[:.\W,]`,
		},
		"not word in negated char class ascii": {
			input: `[^:.\W,]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[^:.\W,]`,
		},
		"digit": {
			input: `\d`,
			want:  `\p{Nd}`,
		},
		"digit in char class": {
			input: `[:,\d.]`,
			want:  `[:,\p{Nd}.]`,
		},
		"digit in negated char class": {
			input: `[^:,\d.]`,
			want:  `[^:,\p{Nd}.]`,
		},
		"digit ascii": {
			input: `\d`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `\d`,
		},
		"digit in char class ascii": {
			input: `[:,\d.]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[:,\d.]`,
		},
		"digit in negated char class ascii": {
			input: `[^:,\d.]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[^:,\d.]`,
		},
		"not digit": {
			input: `\D`,
			want:  `\P{Nd}`,
		},
		"not digit in char class": {
			input: `[:,\D.]`,
			want:  `[:,\P{Nd}.]`,
		},
		"not digit in negated char class": {
			input: `[^9\D0]`,
			want:  `[^9\P{Nd}0]`,
		},
		"not digit ascii": {
			input: `\D`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `\D`,
		},
		"not digit in char class ascii": {
			input: `[9\D0]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[9\D0]`,
		},
		"not digit in negated char class ascii": {
			input: `[^9\D0]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[^9\D0]`,
		},
		"whitespace": {
			input: `\s`,
			want:  `[\s\v\p{Z}\x85]`,
		},
		"whitespace in char class": {
			input: `[:,\s.]`,
			want:  `[:,\s\v\p{Z}\x85.]`,
		},
		"whitespace in negated char class": {
			input: `[^:,\s.]`,
			want:  `[^:,\s\v\p{Z}\x85.]`,
		},
		"whitespace ascii": {
			input: `\s`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `\s`,
		},
		"whitespace in char class ascii": {
			input: `[:,\s.]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[:,\s.]`,
		},
		"whitespace in negated char class ascii": {
			input: `[^:,\s.]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[^:,\s.]`,
		},
		"not whitespace": {
			input: `\S`,
			want:  `[^\s\v\p{Z}\x85]`,
		},
		"not whitespace in char class": {
			input: `[.,\S:]`,
			want:  `(?:[.,:]|[^\s\v\p{Z}\x85])`,
		},
		"not whitespace in negated char class": {
			input: `[^.,\S:]`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(4, 1, 5), P(5, 1, 6)), `double negation of unicode-aware \S is not supported`),
			},
		},
		"not whitespace ascii": {
			input: `\S`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `\S`,
		},
		"not whitespace in char class ascii": {
			input: `[:,\S.]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[:,\S.]`,
		},
		"not whitespace in negated char class ascii": {
			input: `[^:,\S.]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[^:,\S.]`,
		},
		"horizontal whitespace": {
			input: `\h`,
			want:  `[\t\p{Zs}]`,
		},
		"horizontal whitespace in char class": {
			input: `[:,\h.]`,
			want:  `[:,\t\p{Zs}.]`,
		},
		"horizontal whitespace in negated char class": {
			input: `[^:,\h.]`,
			want:  `[^:,\t\p{Zs}.]`,
		},
		"horizontal whitespace ascii": {
			input: `\h`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[\t ]`,
		},
		"horizontal whitespace in char class ascii": {
			input: `[:,\h.]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[:,\t .]`,
		},
		"horizontal whitespace in negated char class ascii": {
			input: `[^:,\h.]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[^:,\t .]`,
		},
		"not horizontal whitespace": {
			input: `\H`,
			want:  `[^\t\p{Zs}]`,
		},
		"not horizontal whitespace in char class": {
			input: `[ab.\Hcd]`,
			want:  `(?:[ab.cd]|[^\t\p{Zs}])`,
		},
		"not horizontal whitespace in negated char class": {
			input: `[^ab.\Hcd]`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(5, 1, 6), P(6, 1, 7)), `double negation of unicode-aware \H is not supported`),
			},
		},
		"not horizontal whitespace ascii": {
			input: `\H`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[^\t ]`,
		},
		"not horizontal whitespace in char class ascii": {
			input: `[ab.\Hcd]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `(?:[ab.cd]|[^\t ])`,
		},
		"not horizontal whitespace in negated char class ascii": {
			input: `[^ab.\Hcd]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(5, 1, 6), P(6, 1, 7)), `double negation of unicode-aware \H is not supported`),
			},
		},
		"vertical whitespace": {
			input: `\v`,
			want:  `[\n\v\f\r\x85\x{2028}\x{2029}]`,
		},
		"vertical whitespace in char class": {
			input: `[ab.\vcd]`,
			want:  `[ab.\n\v\f\r\x85\x{2028}\x{2029}cd]`,
		},
		"vertical whitespace in negated char class": {
			input: `[^ab.\vcd]`,
			want:  `[^ab.\n\v\f\r\x85\x{2028}\x{2029}cd]`,
		},
		"vertical whitespace ascii": {
			input: `\v`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[\n\v\f\r]`,
		},
		"vertical whitespace in char class ascii": {
			input: `[ab.\vcd]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[ab.\n\v\f\rcd]`,
		},
		"vertical whitespace in negated char class ascii": {
			input: `[^ab.\vcd]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[^ab.\n\v\f\rcd]`,
		},
		"not vertical whitespace": {
			input: `\V`,
			want:  `[^\n\v\f\r\x85\x{2028}\x{2029}]`,
		},
		"not vertical whitespace in char class": {
			input: `[ab.\Vcd]`,
			want:  `(?:[ab.cd]|[^\n\v\f\r\x85\x{2028}\x{2029}])`,
		},
		"not vertical whitespace in negated char class": {
			input: `[^ab.\Vcd]`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(5, 1, 6), P(6, 1, 7)), `double negation of unicode-aware \V is not supported`),
			},
		},
		"not vertical whitespace ascii": {
			input: `\V`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `[^\n\v\f\r]`,
		},
		"not vertical whitespace in char class ascii": {
			input: `[ab.\Vcd]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `(?:[ab.cd]|[^\n\v\f\r])`,
		},
		"not vertical whitespace in negated char class ascii": {
			input: `[^ab.\Vcd]`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(5, 1, 6), P(6, 1, 7)), `double negation of unicode-aware \V is not supported`),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestQuotedText(t *testing.T) {
	tests := testTable{
		"simple chars": {
			input: `\Qfoo\E`,
			want:  `\Qfoo\E`,
		},
		"meta chars": {
			input: `\Q+-*.{}()[]?$^\E`,
			want:  `\Q+-*.{}()[]?$^\E`,
		},
		"invalid escapes": {
			input: `\Q\e\E`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(0, 1, 1), P(2, 1, 3)), "expected end of quoted text"),
				error.NewError(L("regex", P(4, 1, 5), P(5, 1, 6)), "invalid escape sequence: \\E"),
			},
		},
		"in char class": {
			input: `[\Q\e\E]`,
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(1, 1, 2), P(3, 1, 4)), "expected end of quoted text"),
				error.NewError(L("regex", P(5, 1, 6), P(6, 1, 7)), "invalid escape sequence: \\E"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestConcatenation(t *testing.T) {
	tests := testTable{
		"ascii chars": {
			input: "foo",
			want:  `foo`,
		},
		"with comments": {
			input: "f(?#some awesome comment)oo",
			want:  `foo`,
		},
		"multi-byte chars": {
			input: `fęłó€𐍈`,
			want:  `fęłó€𐍈`,
		},
		"chars escapes and anchors": {
			input: `f\n\w$`,
			want:  `f\n[\p{L}\p{Mn}\p{Nd}\p{Pc}]$`,
		},
		"chars escapes, anchors and groups": {
			input: `(f\n)\w$`,
			want:  `(f\n)[\p{L}\p{Mn}\p{Nd}\p{Pc}]$`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestCharClass(t *testing.T) {
	tests := testTable{
		"ascii chars": {
			input: "[foa]",
			want:  `[foa]`,
		},
		"negated": {
			input: "[^foa]",
			want:  `[^foa]`,
		},
		"unterminated": {
			input: "[foa",
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(4, 1, 5), P(3, 1, 4)), "unterminated character class, missing ]"),
			},
		},
		"invalid chars": {
			input: "[-]",
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(1, 1, 2), P(1, 1, 2)), "unexpected -, expected a char class element"),
			},
		},
		"char ranges": {
			input: `[a-z\n-\r\x22-\x7f56]`,
			want:  `[a-z\n-\r\x{22}-\x{7f}56]`,
		},
		"meta-chars": {
			input: "[*+.{}()$^|?]",
			want:  `[*+.{}()$^|?]`,
		},
		"multi-byte chars": {
			input: "[fęłó€𐍈]",
			want:  `[fęłó€𐍈]`,
		},
		"escapes and simple char classes": {
			input: `[\n\-\*\.\p{Latin}\x7f\w\s\123\o123]`,
			want:  `[\n\-\*\.\p{Latin}\x{7f}\p{L}\p{Mn}\p{Nd}\p{Pc}\s\v\p{Z}\x85\123\123]`,
		},
		"named char class": {
			input: "[[:alpha:]]",
			want:  `[[:alpha:]]`,
		},
		"named char class with invalid chars": {
			input: "[[:alphę:]]",
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(7, 1, 8), P(8, 1, 8)), "unexpected ę, expected an ASCII letter"),
			},
		},
		"named char class with other elements": {
			input: "[[:alpha:]a-zB]",
			want:  `[[:alpha:]a-zB]`,
		},
		"negated named char class": {
			input: "[[:^alpha:]]",
			want:  `[[:^alpha:]]`,
		},
		"negated named char class in negated char class": {
			input: "[^[:^alpha:]]",
			want:  `[^[:^alpha:]]`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestUnion(t *testing.T) {
	tests := testTable{
		"char union": {
			input: "f|o",
			want:  `f|o`,
		},
		"concat union": {
			input: "foo|barę",
			want:  `foo|barę`,
		},
		"group union": {
			input: "(foo)|barę",
			want:  `(foo)|barę`,
		},
		"nested unions": {
			input: "foo|b|u",
			want:  `foo|b|u`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestGroup(t *testing.T) {
	tests := testTable{
		"non capturing group": {
			input: "(?:f)",
			want:  `(?:f)`,
		},
		"named group": {
			input: "(?<foo>f)",
			want:  `(?P<foo>f)`,
		},
		"named group with single quotes": {
			input: "(?'foo'f)",
			want:  `(?P<foo>f)`,
		},
		"named group with P": {
			input: "(?P<foo>f)",
			want:  `(?P<foo>f)`,
		},
		"flags only": {
			input: "(?imU)",
			want:  `(?imU)`,
		},
		"flags and content": {
			input: "(?mi-s:f)",
			want:  `(?im-s:f)`,
		},
		"invalid flags": {
			input: "(?mihs:f)",
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(4, 1, 5), P(4, 1, 5)), "unexpected h, expected a regex flag"),
			},
		},
		"char in group": {
			input: "(f)",
			want:  `(f)`,
		},
		"missing right paren": {
			input: "(f",
			want:  ``,
			err: error.ErrorList{
				error.NewError(L("regex", P(2, 1, 3), P(1, 1, 2)), "unexpected END_OF_FILE, expected )"),
			},
		},
		"union in group": {
			input: "(foo|barę)",
			want:  `(foo|barę)`,
		},
		"nested groups": {
			input: "((foo)|barę)",
			want:  `((foo)|barę)`,
		},
		"enable ascii flag in global group": {
			input: `(?a)\wfoobar`,
			want:  `\wfoobar`,
		},
		"enable extended and ascii flag in global group": {
			input: `(?xa)\w f  o o bar # comment`,
			want:  `\wfoobar`,
		},
		"enable ascii flag and go supported flag in global group": {
			input: `(?am)\wfoobar`,
			want:  `(?m)\wfoobar`,
		},
		"enable ascii flag in the middle of the global group": {
			input: `\wfoo(?a)\wbar`,
			want:  `[\p{L}\p{Mn}\p{Nd}\p{Pc}]foo\wbar`,
		},
		"enable extended flag in the middle of the global group": {
			input: `\w fo  o(?x)  \d b     a r   # comments are awesome`,
			want:  `[\p{L}\p{Mn}\p{Nd}\p{Pc}] fo  o\p{Nd}bar`,
		},
		"enable ascii flag in a nested group": {
			input: `((?a)\wfoobar)\w`,
			want:  `(\wfoobar)[\p{L}\p{Mn}\p{Nd}\p{Pc}]`,
		},
		"enable ascii flag in the middle of a nested group": {
			input: `(\wfoo(?a)\dbar)\w`,
			want:  `([\p{L}\p{Mn}\p{Nd}\p{Pc}]foo\dbar)[\p{L}\p{Mn}\p{Nd}\p{Pc}]`,
		},
		"enable ascii flag in a nested group directly": {
			input: `(?a:\wfoobar)\w`,
			want:  `(?:\wfoobar)[\p{L}\p{Mn}\p{Nd}\p{Pc}]`,
		},
		"disable ascii flag in a nested group": {
			input: `(?a)\w(?-a:\wfoobar)\w`,
			want:  `\w(?:[\p{L}\p{Mn}\p{Nd}\p{Pc}]foobar)\w`,
		},
		"disable ascii flag in a nested group alt": {
			input: `\w(?-a:\wfoobar)\w`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `\w(?:[\p{L}\p{Mn}\p{Nd}\p{Pc}]foobar)\w`,
		},
		"disable ascii flag in a nested group with other flags": {
			input: `\w(?i-ma:\wfoobar)\w`,
			flags: bitfield.BitField8FromBitFlag(flag.ASCIIFlag),
			want:  `\w(?i-m:[\p{L}\p{Mn}\p{Nd}\p{Pc}]foobar)\w`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}

func TestExtendedFlag(t *testing.T) {
	tests := testTable{
		"ignore whitespace unless escaped": {
			input: `foo\    b   a
			r\ 
			an 	d\ 
			baz`,
			flags: bitfield.BitField8FromBitFlag(flag.ExtendedFlag),
			want:  `foo\ bar\ and\ baz`,
		},
		"comments": {
			input: `foo  b a    # awesome comment
			r
			an 	d\ # another comment
			baz`,
			flags: bitfield.BitField8FromBitFlag(flag.ExtendedFlag),
			want:  `foobarand\ baz`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transpilerTest(tc, t)
		})
	}
}
