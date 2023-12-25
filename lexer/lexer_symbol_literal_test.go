package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestSymbol(t *testing.T) {
	tests := testTable{
		"cannot have whitespace between the colon and the content": {
			input: ": some_awesome_symbol",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(0, 1, 1)), token.COLON),
				V(S(P(2, 1, 3), P(20, 1, 21)), token.PUBLIC_IDENTIFIER, "some_awesome_symbol"),
			},
		},
		"can consist of an publicIdentifier": {
			input: ":some_awesome_symbol",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(0, 1, 1)), token.COLON),
				V(S(P(1, 1, 2), P(19, 1, 20)), token.PUBLIC_IDENTIFIER, "some_awesome_symbol"),
			},
		},
		"can consist of a private identifier": {
			input: ":_some_awesome_symbol",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(0, 1, 1)), token.COLON),
				V(S(P(1, 1, 2), P(20, 1, 21)), token.PRIVATE_IDENTIFIER, "_some_awesome_symbol"),
			},
		},
		"can consist of a public constant": {
			input: ":SomeAwesomeSymbol",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(0, 1, 1)), token.COLON),
				V(S(P(1, 1, 2), P(17, 1, 18)), token.PUBLIC_CONSTANT, "SomeAwesomeSymbol"),
			},
		},
		"can consist of a private constant": {
			input: ":_SomeAwesomeSymbol",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(0, 1, 1)), token.COLON),
				V(S(P(1, 1, 2), P(18, 1, 19)), token.PRIVATE_CONSTANT, "_SomeAwesomeSymbol"),
			},
		},
		"can consist of a raw string": {
			input: ":'symbol from a raw string'",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(0, 1, 1)), token.COLON),
				V(S(P(1, 1, 2), P(26, 1, 27)), token.RAW_STRING, "symbol from a raw string"),
			},
		},
		"can consist of a string": {
			input: `:"symbol from\na string"`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(0, 1, 1)), token.COLON),
				T(S(P(1, 1, 2), P(1, 1, 2)), token.STRING_BEG),
				V(S(P(2, 1, 3), P(22, 1, 23)), token.STRING_CONTENT, "symbol from\na string"),
				T(S(P(23, 1, 24), P(23, 1, 24)), token.STRING_END),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
