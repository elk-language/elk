package lexer

import "testing"

func TestRawString(t *testing.T) {
	tests := testTable{
		"must be terminated": {
			input: "'This is a raw string",
			want: []*Token{
				V(P(0, 21, 1, 1), ErrorToken, "unterminated raw string literal, missing `'`"),
			},
		},
		"does not process escape sequences": {
			input: `'Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a'`,
			want: []*Token{
				V(P(0, 72, 1, 1), RawStringToken, `Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a`),
			},
		},
		"can be multiline": {
			input: `'multiline
strings are
awesome
and really useful'`,
			want: []*Token{
				V(P(0, 49, 1, 1), RawStringToken, "multiline\nstrings are\nawesome\nand really useful"),
			},
		},
		"can contain comments": {
			input: `'some string #[with elk]# comments ##[different]## types # of them'`,
			want: []*Token{
				V(P(0, 67, 1, 1), RawStringToken, "some string #[with elk]# comments ##[different]## types # of them"),
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
			want: []*Token{
				T(P(0, 1, 1, 1), StringBegToken),
				V(P(1, 16, 1, 2), ErrorToken, "unterminated string literal, missing `\"`"),
			},
		},
		"processes escape sequences": {
			input: `"Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a"`,
			want: []*Token{
				T(P(0, 1, 1, 1), StringBegToken),
				V(P(1, 73, 1, 2), StringContentToken, "Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a"),
				T(P(74, 1, 1, 75), StringEndToken),
			},
		},
		"reports errors for invalid escape sequences": {
			input: `"www.foo\yes.com"`,
			want: []*Token{
				T(P(0, 1, 1, 1), StringBegToken),
				V(P(1, 7, 1, 2), StringContentToken, "www.foo"),
				V(P(8, 2, 1, 9), ErrorToken, "invalid escape sequence `\\y` in string literal"),
				V(P(10, 6, 1, 11), StringContentToken, "es.com"),
				T(P(16, 1, 1, 17), StringEndToken),
			},
		},
		"creates errors for invalid hex escapes": {
			input: `"some\xfj string"`,
			want: []*Token{
				T(P(0, 1, 1, 1), StringBegToken),
				V(P(1, 4, 1, 2), StringContentToken, "some"),
				V(P(5, 4, 1, 6), ErrorToken, "invalid hex escape in string literal"),
				V(P(9, 7, 1, 10), StringContentToken, " string"),
				T(P(16, 1, 1, 17), StringEndToken),
			},
		},
		"can be multiline": {
			input: `"multiline
strings are
awesome
and really useful"`,
			want: []*Token{
				T(P(0, 1, 1, 1), StringBegToken),
				V(P(1, 47, 1, 2), StringContentToken, "multiline\nstrings are\nawesome\nand really useful"),
				T(P(48, 1, 4, 18), StringEndToken),
			},
		},
		"can contain comments": {
			input: `"some string #[with elk]# comments ##[different]## types # of them"`,
			want: []*Token{
				T(P(0, 1, 1, 1), StringBegToken),
				V(P(1, 65, 1, 2), StringContentToken, "some string #[with elk]# comments ##[different]## types # of them"),
				T(P(66, 1, 1, 67), StringEndToken),
			},
		},
		"can be interpolated": {
			input: `"some ${interpolated} string ${with.expressions + 2} and end"`,
			want: []*Token{
				T(P(0, 1, 1, 1), StringBegToken),
				V(P(1, 5, 1, 2), StringContentToken, "some "),
				T(P(6, 2, 1, 7), StringInterpBegToken),
				V(P(8, 12, 1, 9), PublicIdentifierToken, "interpolated"),
				T(P(20, 1, 1, 21), StringInterpEndToken),
				V(P(21, 8, 1, 22), StringContentToken, " string "),
				T(P(29, 2, 1, 30), StringInterpBegToken),
				V(P(31, 4, 1, 32), PublicIdentifierToken, "with"),
				T(P(35, 1, 1, 36), DotToken),
				V(P(36, 11, 1, 37), PublicIdentifierToken, "expressions"),
				T(P(48, 1, 1, 49), PlusToken),
				V(P(50, 1, 1, 51), DecIntToken, "2"),
				T(P(51, 1, 1, 52), StringInterpEndToken),
				V(P(52, 8, 1, 53), StringContentToken, " and end"),
				T(P(60, 1, 1, 61), StringEndToken),
			},
		},
		"does not generate unnecessary tokens when interpolation is right beside delimiters": {
			input: `"${interpolated}"`,
			want: []*Token{
				T(P(0, 1, 1, 1), StringBegToken),
				T(P(1, 2, 1, 2), StringInterpBegToken),
				V(P(3, 12, 1, 4), PublicIdentifierToken, "interpolated"),
				T(P(15, 1, 1, 16), StringInterpEndToken),
				T(P(16, 1, 1, 17), StringEndToken),
			},
		},
		"raw strings can be nested in string interpolation": {
			input: `"foo ${baz + 'bar'}"`,
			want: []*Token{
				T(P(0, 1, 1, 1), StringBegToken),
				V(P(1, 4, 1, 2), StringContentToken, "foo "),
				T(P(5, 2, 1, 6), StringInterpBegToken),
				V(P(7, 3, 1, 8), PublicIdentifierToken, "baz"),
				T(P(11, 1, 1, 12), PlusToken),
				V(P(13, 5, 1, 14), RawStringToken, "bar"),
				T(P(18, 1, 1, 19), StringInterpEndToken),
				T(P(19, 1, 1, 20), StringEndToken),
			},
		},
		"strings can't be nested in string interpolation": {
			input: `"foo ${baz + "bar"}"`,
			want: []*Token{
				T(P(0, 1, 1, 1), StringBegToken),
				V(P(1, 4, 1, 2), StringContentToken, "foo "),
				T(P(5, 2, 1, 6), StringInterpBegToken),
				V(P(7, 3, 1, 8), PublicIdentifierToken, "baz"),
				T(P(11, 1, 1, 12), PlusToken),
				V(P(13, 5, 1, 14), ErrorToken, "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation"),
				T(P(18, 1, 1, 19), StringInterpEndToken),
				T(P(19, 1, 1, 20), StringEndToken),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
