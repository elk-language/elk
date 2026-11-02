package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestNewline(t *testing.T) {
	tests := testTable{
		"can escape newline with backslash": {
			input: "5\\\n + 3",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(0, 1, 1))), token.INT, "5"),
				T(L(S(P(4, 1, 5), P(4, 1, 5))), token.PLUS),
				V(L(S(P(6, 1, 7), P(6, 1, 7))), token.INT, "3"),
			},
		},
		"can escape \\r\\n with backslash": {
			input: "5\\\r\n + 3",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(0, 1, 1))), token.INT, "5"),
				T(L(S(P(5, 1, 6), P(5, 1, 6))), token.PLUS),
				V(L(S(P(7, 1, 8), P(7, 1, 8))), token.INT, "3"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
