package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestEmbellishedText(t *testing.T) {
	tests := testTable{
		"with only text": {
			input: `foo bar baz
1 + 2
3
			`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(22, 4, 3))), token.TEXT, "foo bar baz\n1 + 2\n3\n\t\t\t"),
			},
		},
		"with single backtick at the end": {
			input: "foo `bar + 1`",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.TEXT, "foo "),
				V(L(S(P(5, 1, 6), P(7, 1, 8))), token.PUBLIC_IDENTIFIER, "bar"),
				T(L(S(P(9, 1, 10), P(9, 1, 10))), token.PLUS),
				V(L(S(P(11, 1, 12), P(11, 1, 12))), token.INT, "1"),
				T(L(S(P(12, 1, 13), P(12, 1, 13))), token.TEXT),
			},
		},
		"with unterminated single backtick": {
			input: "foo `bar + 1",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.TEXT, "foo "),
				V(L(S(P(5, 1, 6), P(7, 1, 8))), token.PUBLIC_IDENTIFIER, "bar"),
				T(L(S(P(9, 1, 10), P(9, 1, 10))), token.PLUS),
				V(L(S(P(11, 1, 12), P(11, 1, 12))), token.INT, "1"),
			},
		},
		"with single backtick": {
			input: "foo `bar + 1` baz",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.TEXT, "foo "),
				V(L(S(P(5, 1, 6), P(7, 1, 8))), token.PUBLIC_IDENTIFIER, "bar"),
				T(L(S(P(9, 1, 10), P(9, 1, 10))), token.PLUS),
				V(L(S(P(11, 1, 12), P(11, 1, 12))), token.INT, "1"),
				V(L(S(P(12, 1, 13), P(16, 1, 17))), token.TEXT, " baz"),
			},
		},
		"with double backtick": {
			input: "foo ``bar + `1` `` baz",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.TEXT, "foo "),
				V(L(S(P(6, 1, 7), P(8, 1, 9))), token.PUBLIC_IDENTIFIER, "bar"),
				T(L(S(P(10, 1, 11), P(10, 1, 11))), token.PLUS),
				V(L(S(P(12, 1, 13), P(14, 1, 15))), token.CHAR_LITERAL, "1"),
				V(L(S(P(16, 1, 17), P(21, 1, 22))), token.TEXT, " baz"),
			},
		},
		"with tripple backtick": {
			input: "foo ```bar + `1` ``` baz",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 7))), token.TEXT, "foo "),
				V(L(S(P(7, 1, 8), P(9, 1, 10))), token.PUBLIC_IDENTIFIER, "bar"),
				T(L(S(P(11, 1, 12), P(11, 1, 12))), token.PLUS),
				V(L(S(P(13, 1, 14), P(15, 1, 16))), token.CHAR_LITERAL, "1"),
				V(L(S(P(17, 1, 18), P(23, 1, 24))), token.TEXT, " baz"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTestWithMode(tc, embellishmentMode, t)
		})
	}
}
