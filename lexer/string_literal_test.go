package lexer

import "testing"

func TestRawString(t *testing.T) {
	tests := testTable{
		"must be terminated": {
			input: "'This is a raw string",
			want: []*Token{
				V(ErrorToken, "unterminated raw string literal, missing `'`", 0, 21, 1, 1),
			},
		},
		"does not process escape sequences": {
			input: `'Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a'`,
			want: []*Token{
				V(RawStringToken, `Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a`, 0, 72, 1, 1),
			},
		},
		"can be multiline": {
			input: `'multiline
strings are
awesome
and really useful'`,
			want: []*Token{
				V(RawStringToken, "multiline\nstrings are\nawesome\nand really useful", 0, 49, 1, 1),
			},
		},
		"can contain comments": {
			input: `'some string #[with elk]# comments ##[different]## types # of them'`,
			want: []*Token{
				V(RawStringToken, "some string #[with elk]# comments ##[different]## types # of them", 0, 67, 1, 1),
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
				T(StringBegToken, 0, 1, 1, 1),
				V(ErrorToken, "unterminated string literal, missing `\"`", 1, 16, 1, 2),
			},
		},
		"processes escape sequences": {
			input: `"Some \n a\wesome \t str\ing \r with \\ escape \b sequences \"\v\f\x12\a"`,
			want: []*Token{
				T(StringBegToken, 0, 1, 1, 1),
				V(StringContentToken, "Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a", 1, 71, 1, 2),
				T(StringEndToken, 72, 1, 1, 73),
			},
		},
		"creates errors for invalid hex escapes": {
			input: `"some\xfj string"`,
			want: []*Token{
				T(StringBegToken, 0, 1, 1, 1),
				V(ErrorToken, "invalid hex escape", 1, 6, 1, 2),
				V(StringContentToken, "fj string", 7, 9, 1, 8),
				T(StringEndToken, 16, 1, 1, 17),
			},
		},
		"can be multiline": {
			input: `"multiline
strings are
awesome
and really useful"`,
			want: []*Token{
				T(StringBegToken, 0, 1, 1, 1),
				V(StringContentToken, "multiline\nstrings are\nawesome\nand really useful", 1, 47, 1, 2),
				T(StringEndToken, 48, 1, 4, 18),
			},
		},
		"can contain comments": {
			input: `"some string #[with elk]# comments ##[different]## types # of them"`,
			want: []*Token{
				T(StringBegToken, 0, 1, 1, 1),
				V(StringContentToken, "some string #[with elk]# comments ##[different]## types # of them", 1, 65, 1, 2),
				T(StringEndToken, 66, 1, 1, 67),
			},
		},
		"can be interpolated": {
			input: `"some ${interpolated} string ${with.expressions + 2} and end"`,
			want: []*Token{
				T(StringBegToken, 0, 1, 1, 1),
				V(StringContentToken, "some ", 1, 5, 1, 2),
				T(StringInterpBegToken, 6, 2, 1, 7),
				V(IdentifierToken, "interpolated", 8, 12, 1, 9),
				T(StringInterpEndToken, 20, 1, 1, 21),
				V(StringContentToken, " string ", 21, 8, 1, 22),
				T(StringInterpBegToken, 29, 2, 1, 30),
				V(IdentifierToken, "with", 31, 4, 1, 32),
				T(DotToken, 35, 1, 1, 36),
				V(IdentifierToken, "expressions", 36, 11, 1, 37),
				T(PlusToken, 48, 1, 1, 49),
				V(DecIntToken, "2", 50, 1, 1, 51),
				T(StringInterpEndToken, 51, 1, 1, 52),
				V(StringContentToken, " and end", 52, 8, 1, 53),
				T(StringEndToken, 60, 1, 1, 61),
			},
		},
		"does not generate unnecessary tokens when interpolation is right beside delimiters": {
			input: `"${interpolated}"`,
			want: []*Token{
				T(StringBegToken, 0, 1, 1, 1),
				T(StringInterpBegToken, 1, 2, 1, 2),
				V(IdentifierToken, "interpolated", 3, 12, 1, 4),
				T(StringInterpEndToken, 15, 1, 1, 16),
				T(StringEndToken, 16, 1, 1, 17),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
