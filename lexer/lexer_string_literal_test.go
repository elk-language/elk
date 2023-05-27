package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestRawString(t *testing.T) {
	tests := testTable{
		"must be terminated": {
			input: "'This is a raw string",
			want: []*token.Token{
				V(P(0, 21, 1, 1), token.ERROR, "unterminated raw string literal, missing `'`"),
			},
		},
		"does not process escape sequences": {
			input: `'Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a'`,
			want: []*token.Token{
				V(P(0, 72, 1, 1), token.RAW_STRING, `Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a`),
			},
		},
		"can be multiline": {
			input: `'multiline
strings are
awesome
and really useful'`,
			want: []*token.Token{
				V(P(0, 49, 1, 1), token.RAW_STRING, "multiline\nstrings are\nawesome\nand really useful"),
			},
		},
		"can contain comments": {
			input: `'some string #[with elk]# comments ##[different]## types # of them'`,
			want: []*token.Token{
				V(P(0, 67, 1, 1), token.RAW_STRING, "some string #[with elk]# comments ##[different]## types # of them"),
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
				T(P(0, 1, 1, 1), token.STRING_BEG),
				V(P(1, 16, 1, 2), token.ERROR, "unterminated string literal, missing `\"`"),
			},
		},
		"processes escape sequences": {
			input: `"Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a"`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.STRING_BEG),
				V(P(1, 73, 1, 2), token.STRING_CONTENT, "Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a"),
				T(P(74, 1, 1, 75), token.STRING_END),
			},
		},
		"reports errors for invalid escape sequences": {
			input: `"www.foo\yes.com"`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.STRING_BEG),
				V(P(1, 7, 1, 2), token.STRING_CONTENT, "www.foo"),
				V(P(8, 2, 1, 9), token.ERROR, "invalid escape sequence `\\y` in string literal"),
				V(P(10, 6, 1, 11), token.STRING_CONTENT, "es.com"),
				T(P(16, 1, 1, 17), token.STRING_END),
			},
		},
		"creates errors for invalid hex escapes": {
			input: `"some\xfj string"`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.STRING_BEG),
				V(P(1, 4, 1, 2), token.STRING_CONTENT, "some"),
				V(P(5, 4, 1, 6), token.ERROR, "invalid hex escape in string literal"),
				V(P(9, 7, 1, 10), token.STRING_CONTENT, " string"),
				T(P(16, 1, 1, 17), token.STRING_END),
			},
		},
		"can be multiline": {
			input: `"multiline
strings are
awesome
and really useful"`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.STRING_BEG),
				V(P(1, 47, 1, 2), token.STRING_CONTENT, "multiline\nstrings are\nawesome\nand really useful"),
				T(P(48, 1, 4, 18), token.STRING_END),
			},
		},
		"can contain comments": {
			input: `"some string #[with elk]# comments ##[different]## types # of them"`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.STRING_BEG),
				V(P(1, 65, 1, 2), token.STRING_CONTENT, "some string #[with elk]# comments ##[different]## types # of them"),
				T(P(66, 1, 1, 67), token.STRING_END),
			},
		},
		"can be interpolated": {
			input: `"some ${interpolated} string ${with.expressions + 2} and end"`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.STRING_BEG),
				V(P(1, 5, 1, 2), token.STRING_CONTENT, "some "),
				T(P(6, 2, 1, 7), token.STRING_INTERP_BEG),
				V(P(8, 12, 1, 9), token.PUBLIC_IDENTIFIER, "interpolated"),
				T(P(20, 1, 1, 21), token.STRING_INTERP_END),
				V(P(21, 8, 1, 22), token.STRING_CONTENT, " string "),
				T(P(29, 2, 1, 30), token.STRING_INTERP_BEG),
				V(P(31, 4, 1, 32), token.PUBLIC_IDENTIFIER, "with"),
				T(P(35, 1, 1, 36), token.DOT),
				V(P(36, 11, 1, 37), token.PUBLIC_IDENTIFIER, "expressions"),
				T(P(48, 1, 1, 49), token.PLUS),
				V(P(50, 1, 1, 51), token.DEC_INT, "2"),
				T(P(51, 1, 1, 52), token.STRING_INTERP_END),
				V(P(52, 8, 1, 53), token.STRING_CONTENT, " and end"),
				T(P(60, 1, 1, 61), token.STRING_END),
			},
		},
		"does not generate unnecessary tokens when interpolation is right beside delimiters": {
			input: `"${interpolated}"`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.STRING_BEG),
				T(P(1, 2, 1, 2), token.STRING_INTERP_BEG),
				V(P(3, 12, 1, 4), token.PUBLIC_IDENTIFIER, "interpolated"),
				T(P(15, 1, 1, 16), token.STRING_INTERP_END),
				T(P(16, 1, 1, 17), token.STRING_END),
			},
		},
		"raw strings can be nested in string interpolation": {
			input: `"foo ${baz + 'bar'}"`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.STRING_BEG),
				V(P(1, 4, 1, 2), token.STRING_CONTENT, "foo "),
				T(P(5, 2, 1, 6), token.STRING_INTERP_BEG),
				V(P(7, 3, 1, 8), token.PUBLIC_IDENTIFIER, "baz"),
				T(P(11, 1, 1, 12), token.PLUS),
				V(P(13, 5, 1, 14), token.RAW_STRING, "bar"),
				T(P(18, 1, 1, 19), token.STRING_INTERP_END),
				T(P(19, 1, 1, 20), token.STRING_END),
			},
		},
		"strings can't be nested in string interpolation": {
			input: `"foo ${baz + "bar"}"`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.STRING_BEG),
				V(P(1, 4, 1, 2), token.STRING_CONTENT, "foo "),
				T(P(5, 2, 1, 6), token.STRING_INTERP_BEG),
				V(P(7, 3, 1, 8), token.PUBLIC_IDENTIFIER, "baz"),
				T(P(11, 1, 1, 12), token.PLUS),
				V(P(13, 5, 1, 14), token.ERROR, "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation"),
				T(P(18, 1, 1, 19), token.STRING_INTERP_END),
				T(P(19, 1, 1, 20), token.STRING_END),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}