package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestChar(t *testing.T) {
	tests := testTable{
		"must be terminated": {
			input: "`a",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.ERROR, "unterminated character literal, missing backtick"),
			},
		},
		"can contain ascii characters": {
			input: "`a`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(2, 1, 3))), token.CHAR_LITERAL, "a"),
			},
		},
		"can contain utf8 characters": {
			input: "`ś`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 3))), token.CHAR_LITERAL, "ś"),
			},
		},
		"escapes backticks": {
			input: "`\\``",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.CHAR_LITERAL, "`"),
			},
		},
		"escapes newlines": {
			input: "`\\n`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.CHAR_LITERAL, "\n"),
			},
		},
		"escapes backslashes": {
			input: "`\\\\`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.CHAR_LITERAL, `\`),
			},
		},
		"escapes tabs": {
			input: "`\\t`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.CHAR_LITERAL, "\t"),
			},
		},
		"escapes carriage returns": {
			input: "`\\r`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.CHAR_LITERAL, "\r"),
			},
		},
		"escapes backspaces": {
			input: "`\\b`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.CHAR_LITERAL, "\b"),
			},
		},
		"escapes vertical tabs": {
			input: "`\\v`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.CHAR_LITERAL, "\v"),
			},
		},
		"escapes form feeds": {
			input: "`\\f`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.CHAR_LITERAL, "\f"),
			},
		},
		"escapes hex": {
			input: "`\\x12`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.CHAR_LITERAL, "\x12"),
			},
		},
		"escapes alerts": {
			input: "`\\a`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.CHAR_LITERAL, "\a"),
			},
		},
		"escapes unicode": {
			input: "`\\u00e9`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(7, 1, 8))), token.CHAR_LITERAL, "\u00e9"),
			},
		},
		"escapes big unicode": {
			input: "`\\U0010FFFF`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(11, 1, 12))), token.CHAR_LITERAL, "\U0010FFFF"),
			},
		},
		"cannot contain multiple characters": {
			input: "`lalala`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(7, 1, 8))), token.ERROR, "invalid char literal with more than one character"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestRawChar(t *testing.T) {
	tests := testTable{
		"must be terminated": {
			input: "r`a",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(2, 1, 3))), token.ERROR, "unterminated character literal, missing backtick"),
			},
		},
		"can contain ascii characters": {
			input: "r`a`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.RAW_CHAR_LITERAL, "a"),
			},
		},
		"can contain backslash": {
			input: "r`\\`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.RAW_CHAR_LITERAL, `\`),
			},
		},
		"can contain utf8 characters": {
			input: "r`ś`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 4))), token.RAW_CHAR_LITERAL, "ś"),
			},
		},
		"cannot escapes backticks": {
			input: "r`\\``",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.RAW_CHAR_LITERAL, `\`),
				V(L(S(P(4, 1, 5), P(4, 1, 5))), token.ERROR, "unterminated character literal, missing backtick"),
			},
		},
		"doesn't process escapes": {
			input: "r`\\n`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.ERROR, "invalid raw char literal with more than one character"),
			},
		},
		"cannot contain multiple characters": {
			input: "r`lalala`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(8, 1, 9))), token.ERROR, "invalid raw char literal with more than one character"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestRawString(t *testing.T) {
	tests := testTable{
		"must be terminated": {
			input: "'This is a raw string",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(20, 1, 21))), token.ERROR, "unterminated raw string literal, missing `'`"),
			},
		},
		"does not process escape sequences": {
			input: `'Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a'`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(71, 1, 72))), token.RAW_STRING, `Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a`),
			},
		},
		"can be multiline": {
			input: `'multiline
strings are
awesome
and really useful'`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(48, 4, 18))), token.RAW_STRING, "multiline\nstrings are\nawesome\nand really useful"),
			},
		},
		"can contain comments": {
			input: `'some string #[with elk]# comments ##[different]## types # of them'`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(66, 1, 67))), token.RAW_STRING, "some string #[with elk]# comments ##[different]## types # of them"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestString(t *testing.T) {
	tests := testTable{
		"must be terminated": {
			input: `"This is a string`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(16, 1, 17))), token.ERROR, "unterminated string literal, missing `\"`"),
			},
		},
		"processes escape sequences": {
			input: `"Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF\$\#"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(93, 1, 94))), token.STRING_CONTENT, "Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF$#"),
				T(L(S(P(94, 1, 95), P(94, 1, 95))), token.STRING_END),
			},
		},
		"reports errors for invalid escape sequences": {
			input: `"www.foo\yes.com"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(7, 1, 8))), token.STRING_CONTENT, "www.foo"),
				V(L(S(P(8, 1, 9), P(9, 1, 10))), token.ERROR, "invalid escape sequence `\\y` in string literal"),
				V(L(S(P(10, 1, 11), P(15, 1, 16))), token.STRING_CONTENT, "es.com"),
				T(L(S(P(16, 1, 17), P(16, 1, 17))), token.STRING_END),
			},
		},
		"creates errors for invalid hex escapes": {
			input: `"some\xfj string"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(4, 1, 5))), token.STRING_CONTENT, "some"),
				V(L(S(P(5, 1, 6), P(8, 1, 9))), token.ERROR, "invalid hex escape"),
				V(L(S(P(9, 1, 10), P(15, 1, 16))), token.STRING_CONTENT, " string"),
				T(L(S(P(16, 1, 17), P(16, 1, 17))), token.STRING_END),
			},
		},
		"creates errors for invalid unicode escapes": {
			input: `"some\uiaab string"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(4, 1, 5))), token.STRING_CONTENT, "some"),
				V(L(S(P(5, 1, 6), P(10, 1, 11))), token.ERROR, "invalid unicode escape"),
				V(L(S(P(11, 1, 12), P(17, 1, 18))), token.STRING_CONTENT, " string"),
				T(L(S(P(18, 1, 19), P(18, 1, 19))), token.STRING_END),
			},
		},
		"creates errors for invalid big unicode escapes": {
			input: `"some\Uiaabuj46 string"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(4, 1, 5))), token.STRING_CONTENT, "some"),
				V(L(S(P(5, 1, 6), P(14, 1, 15))), token.ERROR, "invalid unicode escape"),
				V(L(S(P(15, 1, 16), P(21, 1, 22))), token.STRING_CONTENT, " string"),
				T(L(S(P(22, 1, 23), P(22, 1, 23))), token.STRING_END),
			},
		},
		"can be multiline": {
			input: `"multiline
strings are
awesome
and really useful"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(47, 4, 17))), token.STRING_CONTENT, "multiline\nstrings are\nawesome\nand really useful"),
				T(L(S(P(48, 4, 18), P(48, 4, 18))), token.STRING_END),
			},
		},
		"can contain comments": {
			input: `"some string #[with elk]# comments ##[different]## types # of them"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(65, 1, 66))), token.STRING_CONTENT, "some string #[with elk]# comments ##[different]## types # of them"),
				T(L(S(P(66, 1, 67), P(66, 1, 67))), token.STRING_END),
			},
		},
		"can be interpolated": {
			input: `"some ${interpolated} string ${wita.expressions + 2} and end"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(5, 1, 6))), token.STRING_CONTENT, "some "),
				T(L(S(P(6, 1, 7), P(7, 1, 8))), token.STRING_INTERP_BEG),
				V(L(S(P(8, 1, 9), P(19, 1, 20))), token.PUBLIC_IDENTIFIER, "interpolated"),
				T(L(S(P(20, 1, 21), P(20, 1, 21))), token.STRING_INTERP_END),
				V(L(S(P(21, 1, 22), P(28, 1, 29))), token.STRING_CONTENT, " string "),
				T(L(S(P(29, 1, 30), P(30, 1, 31))), token.STRING_INTERP_BEG),
				V(L(S(P(31, 1, 32), P(34, 1, 35))), token.PUBLIC_IDENTIFIER, "wita"),
				T(L(S(P(35, 1, 36), P(35, 1, 36))), token.DOT),
				V(L(S(P(36, 1, 37), P(46, 1, 47))), token.PUBLIC_IDENTIFIER, "expressions"),
				T(L(S(P(48, 1, 49), P(48, 1, 49))), token.PLUS),
				V(L(S(P(50, 1, 51), P(50, 1, 51))), token.INT, "2"),
				T(L(S(P(51, 1, 52), P(51, 1, 52))), token.STRING_INTERP_END),
				V(L(S(P(52, 1, 53), P(59, 1, 60))), token.STRING_CONTENT, " and end"),
				T(L(S(P(60, 1, 61), P(60, 1, 61))), token.STRING_END),
			},
		},
		"does not generate unnecessary tokens when interpolation is right beside delimiters": {
			input: `"${interpolated}"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				T(L(S(P(1, 1, 2), P(2, 1, 3))), token.STRING_INTERP_BEG),
				V(L(S(P(3, 1, 4), P(14, 1, 15))), token.PUBLIC_IDENTIFIER, "interpolated"),
				T(L(S(P(15, 1, 16), P(15, 1, 16))), token.STRING_INTERP_END),
				T(L(S(P(16, 1, 17), P(16, 1, 17))), token.STRING_END),
			},
		},
		"can be inspect interpolated": {
			input: `"foo #{interpolated}"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(4, 1, 5))), token.STRING_CONTENT, "foo "),
				T(L(S(P(5, 1, 6), P(6, 1, 7))), token.STRING_INSPECT_INTERP_BEG),
				V(L(S(P(7, 1, 8), P(18, 1, 19))), token.PUBLIC_IDENTIFIER, "interpolated"),
				T(L(S(P(19, 1, 20), P(19, 1, 20))), token.STRING_INTERP_END),
				T(L(S(P(20, 1, 21), P(20, 1, 21))), token.STRING_END),
			},
		},
		"can contain short local and constant inspect interpolation": {
			input: `"#foo, #_foo, #Bar and #_Bar"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(4, 1, 5))), token.STRING_INSPECT_INTERP_LOCAL, "foo"),
				V(L(S(P(5, 1, 6), P(6, 1, 7))), token.STRING_CONTENT, ", "),
				V(L(S(P(7, 1, 8), P(11, 1, 12))), token.STRING_INSPECT_INTERP_LOCAL, "_foo"),
				V(L(S(P(12, 1, 13), P(13, 1, 14))), token.STRING_CONTENT, ", "),
				V(L(S(P(14, 1, 15), P(17, 1, 18))), token.STRING_INSPECT_INTERP_CONSTANT, "Bar"),
				V(L(S(P(18, 1, 19), P(22, 1, 23))), token.STRING_CONTENT, " and "),
				V(L(S(P(23, 1, 24), P(27, 1, 28))), token.STRING_INSPECT_INTERP_CONSTANT, "_Bar"),
				T(L(S(P(28, 1, 29), P(28, 1, 29))), token.STRING_END),
			},
		},
		"can contain short local and constant interpolation": {
			input: `"$foo, $_foo, $Bar and $_Bar"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(4, 1, 5))), token.STRING_INTERP_LOCAL, "foo"),
				V(L(S(P(5, 1, 6), P(6, 1, 7))), token.STRING_CONTENT, ", "),
				V(L(S(P(7, 1, 8), P(11, 1, 12))), token.STRING_INTERP_LOCAL, "_foo"),
				V(L(S(P(12, 1, 13), P(13, 1, 14))), token.STRING_CONTENT, ", "),
				V(L(S(P(14, 1, 15), P(17, 1, 18))), token.STRING_INTERP_CONSTANT, "Bar"),
				V(L(S(P(18, 1, 19), P(22, 1, 23))), token.STRING_CONTENT, " and "),
				V(L(S(P(23, 1, 24), P(27, 1, 28))), token.STRING_INTERP_CONSTANT, "_Bar"),
				T(L(S(P(28, 1, 29), P(28, 1, 29))), token.STRING_END),
			},
		},
		"raw strings can be nested in string interpolation": {
			input: `"foo ${baz + 'bar'}"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(4, 1, 5))), token.STRING_CONTENT, "foo "),
				T(L(S(P(5, 1, 6), P(6, 1, 7))), token.STRING_INTERP_BEG),
				V(L(S(P(7, 1, 8), P(9, 1, 10))), token.PUBLIC_IDENTIFIER, "baz"),
				T(L(S(P(11, 1, 12), P(11, 1, 12))), token.PLUS),
				V(L(S(P(13, 1, 14), P(17, 1, 18))), token.RAW_STRING, "bar"),
				T(L(S(P(18, 1, 19), P(18, 1, 19))), token.STRING_INTERP_END),
				T(L(S(P(19, 1, 20), P(19, 1, 20))), token.STRING_END),
			},
		},
		"strings cannot be nested in string interpolation": {
			input: `"foo ${baz + "bar"}"`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.STRING_BEG),
				V(L(S(P(1, 1, 2), P(4, 1, 5))), token.STRING_CONTENT, "foo "),
				T(L(S(P(5, 1, 6), P(6, 1, 7))), token.STRING_INTERP_BEG),
				V(L(S(P(7, 1, 8), P(9, 1, 10))), token.PUBLIC_IDENTIFIER, "baz"),
				T(L(S(P(11, 1, 12), P(11, 1, 12))), token.PLUS),
				V(L(S(P(13, 1, 14), P(17, 1, 18))), token.ERROR, "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation"),
				T(L(S(P(18, 1, 19), P(18, 1, 19))), token.STRING_INTERP_END),
				T(L(S(P(19, 1, 20), P(19, 1, 20))), token.STRING_END),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
