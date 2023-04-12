package lexer

import "testing"

func TestRawString(t *testing.T) {
	tests := testTable{
		"must be terminated": {
			input: "'This is a raw string",
			want: []*Token{
				{
					TokenType:  ErrorToken,
					Value:      "unterminated raw string literal, missing `'`",
					StartByte:  0,
					ByteLength: 21,
					Line:       1,
					Column:     1,
				},
			},
		},
		"does not process escape sequences": {
			input: `'Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a'`,
			want: []*Token{
				{
					TokenType:  RawStringToken,
					Value:      `Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a`,
					StartByte:  0,
					ByteLength: 72,
					Line:       1,
					Column:     1,
				},
			},
		},
		"can be multiline": {
			input: `'multiline
strings are
awesome
and really useful'`,
			want: []*Token{
				{
					TokenType:  RawStringToken,
					Value:      "multiline\nstrings are\nawesome\nand really useful",
					StartByte:  0,
					ByteLength: 49,
					Line:       1,
					Column:     1,
				},
			},
		},
		"can contain comments": {
			input: `'some string #[with elk]# comments ##[different]## types # of them'`,
			want: []*Token{
				{
					TokenType:  RawStringToken,
					Value:      "some string #[with elk]# comments ##[different]## types # of them",
					StartByte:  0,
					ByteLength: 67,
					Line:       1,
					Column:     1,
				},
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
				{
					TokenType:  StringBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  ErrorToken,
					Value:      "unterminated string literal, missing `\"`",
					StartByte:  1,
					ByteLength: 16,
					Line:       1,
					Column:     2,
				},
			},
		},
		"processes escape sequences": {
			input: `"Some \n a\wesome \t str\ing \r with \\ escape \b sequences \"\v\f\x12\a"`,
			want: []*Token{
				{
					TokenType:  StringBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  StringContentToken,
					Value:      "Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a",
					StartByte:  1,
					ByteLength: 71,
					Line:       1,
					Column:     2,
				},
				{
					TokenType:  StringEndToken,
					StartByte:  72,
					ByteLength: 1,
					Line:       1,
					Column:     73,
				},
			},
		},
		"creates errors for invalid hex escapes": {
			input: `"some\xfj string"`,
			want: []*Token{
				{
					TokenType:  StringBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  ErrorToken,
					Value:      "invalid hex escape",
					StartByte:  1,
					ByteLength: 6,
					Line:       1,
					Column:     2,
				},
				{
					TokenType:  StringContentToken,
					Value:      "fj string",
					StartByte:  7,
					ByteLength: 9,
					Line:       1,
					Column:     8,
				},
				{
					TokenType:  StringEndToken,
					StartByte:  16,
					ByteLength: 1,
					Line:       1,
					Column:     17,
				},
			},
		},
		"can be multiline": {
			input: `"multiline
strings are
awesome
and really useful"`,
			want: []*Token{
				{
					TokenType:  StringBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  StringContentToken,
					Value:      "multiline\nstrings are\nawesome\nand really useful",
					StartByte:  1,
					ByteLength: 47,
					Line:       1,
					Column:     2,
				},
				{
					TokenType:  StringEndToken,
					StartByte:  48,
					ByteLength: 1,
					Line:       4,
					Column:     18,
				},
			},
		},
		"can contain comments": {
			input: `"some string #[with elk]# comments ##[different]## types # of them"`,
			want: []*Token{
				{
					TokenType:  StringBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  StringContentToken,
					Value:      "some string #[with elk]# comments ##[different]## types # of them",
					StartByte:  1,
					ByteLength: 65,
					Line:       1,
					Column:     2,
				},
				{
					TokenType:  StringEndToken,
					StartByte:  66,
					ByteLength: 1,
					Line:       1,
					Column:     67,
				},
			},
		},
		"can be interpolated": {
			input: `"some ${interpolated} string ${with.expressions + 2} and end"`,
			want: []*Token{
				{
					TokenType:  StringBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  StringContentToken,
					Value:      "some ",
					StartByte:  1,
					ByteLength: 5,
					Line:       1,
					Column:     2,
				},
				{
					TokenType:  StringInterpBegToken,
					StartByte:  6,
					ByteLength: 2,
					Line:       1,
					Column:     7,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "interpolated",
					StartByte:  8,
					ByteLength: 12,
					Line:       1,
					Column:     9,
				},
				{
					TokenType:  StringInterpEndToken,
					StartByte:  20,
					ByteLength: 1,
					Line:       1,
					Column:     21,
				},
				{
					TokenType:  StringContentToken,
					Value:      " string ",
					StartByte:  21,
					ByteLength: 8,
					Line:       1,
					Column:     22,
				},
				{
					TokenType:  StringInterpBegToken,
					StartByte:  29,
					ByteLength: 2,
					Line:       1,
					Column:     30,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "with",
					StartByte:  31,
					ByteLength: 4,
					Line:       1,
					Column:     32,
				},
				{
					TokenType:  DotToken,
					StartByte:  35,
					ByteLength: 1,
					Line:       1,
					Column:     36,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "expressions",
					StartByte:  36,
					ByteLength: 11,
					Line:       1,
					Column:     37,
				},
				{
					TokenType:  PlusToken,
					StartByte:  48,
					ByteLength: 1,
					Line:       1,
					Column:     49,
				},
				{
					TokenType:  IntToken,
					Value:      "2",
					StartByte:  50,
					ByteLength: 1,
					Line:       1,
					Column:     51,
				},
				{
					TokenType:  StringInterpEndToken,
					StartByte:  51,
					ByteLength: 1,
					Line:       1,
					Column:     52,
				},
				{
					TokenType:  StringContentToken,
					Value:      " and end",
					StartByte:  52,
					ByteLength: 8,
					Line:       1,
					Column:     53,
				},
				{
					TokenType:  StringEndToken,
					StartByte:  60,
					ByteLength: 1,
					Line:       1,
					Column:     61,
				},
			},
		},
		"does not generate unnecessary tokens when interpolation is right beside delimiters": {
			input: `"${interpolated}"`,
			want: []*Token{
				{
					TokenType:  StringBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  StringInterpBegToken,
					StartByte:  1,
					ByteLength: 2,
					Line:       1,
					Column:     2,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "interpolated",
					StartByte:  3,
					ByteLength: 12,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  StringInterpEndToken,
					StartByte:  15,
					ByteLength: 1,
					Line:       1,
					Column:     16,
				},
				{
					TokenType:  StringEndToken,
					StartByte:  16,
					ByteLength: 1,
					Line:       1,
					Column:     17,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
