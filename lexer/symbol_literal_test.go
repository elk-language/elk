package lexer

import "testing"

func TestSymbol(t *testing.T) {
	tests := testTable{
		"can't have whitespace between the colon and the content": {
			input: ": some_awesome_symbol",
			want: []*Token{
				{
					TokenType:  ColonToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "some_awesome_symbol",
					StartByte:  2,
					ByteLength: 19,
					Line:       1,
					Column:     3,
				},
			},
		},
		"can consist of an identifier": {
			input: ":some_awesome_symbol",
			want: []*Token{
				{
					TokenType:  SymbolBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "some_awesome_symbol",
					StartByte:  1,
					ByteLength: 19,
					Line:       1,
					Column:     2,
				},
			},
		},
		"can consist of a private identifier": {
			input: ":_some_awesome_symbol",
			want: []*Token{
				{
					TokenType:  SymbolBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  PrivateIdentifierToken,
					Value:      "_some_awesome_symbol",
					StartByte:  1,
					ByteLength: 20,
					Line:       1,
					Column:     2,
				},
			},
		},
		"can consist of a constant": {
			input: ":SomeAwesomeSymbol",
			want: []*Token{
				{
					TokenType:  SymbolBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  ConstantToken,
					Value:      "SomeAwesomeSymbol",
					StartByte:  1,
					ByteLength: 17,
					Line:       1,
					Column:     2,
				},
			},
		},
		"can consist of a private constant": {
			input: ":_SomeAwesomeSymbol",
			want: []*Token{
				{
					TokenType:  SymbolBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  PrivateConstantToken,
					Value:      "_SomeAwesomeSymbol",
					StartByte:  1,
					ByteLength: 18,
					Line:       1,
					Column:     2,
				},
			},
		},
		"can consist of a raw string": {
			input: ":'symbol from a raw string'",
			want: []*Token{
				{
					TokenType:  SymbolBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  RawStringToken,
					Value:      "symbol from a raw string",
					StartByte:  1,
					ByteLength: 26,
					Line:       1,
					Column:     2,
				},
			},
		},
		"can consist of a string": {
			input: `:"symbol from\na string"`,
			want: []*Token{
				{
					TokenType:  SymbolBegToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  StringBegToken,
					StartByte:  1,
					ByteLength: 1,
					Line:       1,
					Column:     2,
				},
				{
					TokenType:  StringContentToken,
					Value:      "symbol from\na string",
					StartByte:  2,
					ByteLength: 21,
					Line:       1,
					Column:     3,
				},
				{
					TokenType:  StringEndToken,
					StartByte:  23,
					ByteLength: 1,
					Line:       1,
					Column:     24,
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
