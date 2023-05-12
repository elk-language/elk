package lexer

import "testing"

func TestSymbol(t *testing.T) {
	tests := testTable{
		"can't have whitespace between the colon and the content": {
			input: ": some_awesome_symbol",
			want: []*Token{
				T(P(0, 1, 1, 1), ColonToken),
				V(P(2, 19, 1, 3), PublicIdentifierToken, "some_awesome_symbol"),
			},
		},
		"can consist of an publicIdentifier": {
			input: ":some_awesome_symbol",
			want: []*Token{
				T(P(0, 1, 1, 1), SymbolBegToken),
				V(P(1, 19, 1, 2), PublicIdentifierToken, "some_awesome_symbol"),
			},
		},
		"can consist of a private identifier": {
			input: ":_some_awesome_symbol",
			want: []*Token{
				T(P(0, 1, 1, 1), SymbolBegToken),
				V(P(1, 20, 1, 2), PrivateIdentifierToken, "_some_awesome_symbol"),
			},
		},
		"can consist of a public constant": {
			input: ":SomeAwesomeSymbol",
			want: []*Token{
				T(P(0, 1, 1, 1), SymbolBegToken),
				V(P(1, 17, 1, 2), PublicConstantToken, "SomeAwesomeSymbol"),
			},
		},
		"can consist of a private constant": {
			input: ":_SomeAwesomeSymbol",
			want: []*Token{
				T(P(0, 1, 1, 1), SymbolBegToken),
				V(P(1, 18, 1, 2), PrivateConstantToken, "_SomeAwesomeSymbol"),
			},
		},
		"can consist of a raw string": {
			input: ":'symbol from a raw string'",
			want: []*Token{
				T(P(0, 1, 1, 1), SymbolBegToken),
				V(P(1, 26, 1, 2), RawStringToken, "symbol from a raw string"),
			},
		},
		"can consist of a string": {
			input: `:"symbol from\na string"`,
			want: []*Token{
				T(P(0, 1, 1, 1), SymbolBegToken),
				T(P(1, 1, 1, 2), StringBegToken),
				V(P(2, 21, 1, 3), StringContentToken, "symbol from\na string"),
				T(P(23, 1, 1, 24), StringEndToken),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
