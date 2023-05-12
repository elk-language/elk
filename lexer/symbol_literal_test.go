package lexer

import "testing"

func TestSymbol(t *testing.T) {
	tests := testTable{
		"can't have whitespace between the colon and the content": {
			input: ": some_awesome_symbol",
			want: []*Token{
				T(ColonToken, 0, 1, 1, 1),
				V(PublicIdentifierToken, "some_awesome_symbol", 2, 19, 1, 3),
			},
		},
		"can consist of an publicIdentifier": {
			input: ":some_awesome_symbol",
			want: []*Token{
				T(SymbolBegToken, 0, 1, 1, 1),
				V(PublicIdentifierToken, "some_awesome_symbol", 1, 19, 1, 2),
			},
		},
		"can consist of a private identifier": {
			input: ":_some_awesome_symbol",
			want: []*Token{
				T(SymbolBegToken, 0, 1, 1, 1),
				V(PrivateIdentifierToken, "_some_awesome_symbol", 1, 20, 1, 2),
			},
		},
		"can consist of a public constant": {
			input: ":SomeAwesomeSymbol",
			want: []*Token{
				T(SymbolBegToken, 0, 1, 1, 1),
				V(PublicConstantToken, "SomeAwesomeSymbol", 1, 17, 1, 2),
			},
		},
		"can consist of a private constant": {
			input: ":_SomeAwesomeSymbol",
			want: []*Token{
				T(SymbolBegToken, 0, 1, 1, 1),
				V(PrivateConstantToken, "_SomeAwesomeSymbol", 1, 18, 1, 2),
			},
		},
		"can consist of a raw string": {
			input: ":'symbol from a raw string'",
			want: []*Token{
				T(SymbolBegToken, 0, 1, 1, 1),
				V(RawStringToken, "symbol from a raw string", 1, 26, 1, 2),
			},
		},
		"can consist of a string": {
			input: `:"symbol from\na string"`,
			want: []*Token{
				T(SymbolBegToken, 0, 1, 1, 1),
				T(StringBegToken, 1, 1, 1, 2),
				V(StringContentToken, "symbol from\na string", 2, 21, 1, 3),
				T(StringEndToken, 23, 1, 1, 24),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
