package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestRegex(t *testing.T) {
	tests := testTable{
		"must be terminated": {
			input: `%/This is a regex`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.REGEX_BEG),
				V(S(P(2, 1, 3), P(16, 1, 17)), token.ERROR, "unterminated regex literal, missing `/`"),
			},
		},
		"does not process string escape sequences": {
			input: `%/Some \w+\d \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF/`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.REGEX_BEG),
				V(S(P(2, 1, 3), P(96, 1, 97)), token.REGEX_CONTENT, `Some \w+\d \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF`),
				T(S(P(97, 1, 98), P(97, 1, 98)), token.REGEX_END),
			},
		},
		"supports slash escapes": {
			input: `%/foo\/bar/`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.REGEX_BEG),
				V(S(P(2, 1, 3), P(9, 1, 10)), token.REGEX_CONTENT, `foo\/bar`),
				T(S(P(10, 1, 11), P(10, 1, 11)), token.REGEX_END),
			},
		},
		"ends on first unescaped slash": {
			input: `%/foo/bar/`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.REGEX_BEG),
				V(S(P(2, 1, 3), P(4, 1, 5)), token.REGEX_CONTENT, `foo`),
				T(S(P(5, 1, 6), P(5, 1, 6)), token.REGEX_END),
				V(S(P(6, 1, 7), P(8, 1, 9)), token.PUBLIC_IDENTIFIER, "bar"),
				T(S(P(9, 1, 10), P(9, 1, 10)), token.SLASH),
			},
		},
		"can be multiline": {
			input: `%/multiline
regexes are
\w+
and really \w+/`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.REGEX_BEG),
				V(S(P(2, 1, 3), P(41, 4, 14)), token.REGEX_CONTENT, "multiline\nregexes are\n\\w+\nand really \\w+"),
				T(S(P(42, 4, 15), P(42, 4, 15)), token.REGEX_END),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
