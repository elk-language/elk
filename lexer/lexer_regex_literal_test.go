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
			input: `%/foo/ bar/`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.REGEX_BEG),
				V(S(P(2, 1, 3), P(4, 1, 5)), token.REGEX_CONTENT, `foo`),
				T(S(P(5, 1, 6), P(5, 1, 6)), token.REGEX_END),
				V(S(P(7, 1, 8), P(9, 1, 10)), token.PUBLIC_IDENTIFIER, "bar"),
				T(S(P(10, 1, 11), P(10, 1, 11)), token.SLASH),
			},
		},
		"can have flags": {
			input: `%/foo\/bar/Uimaxs`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.REGEX_BEG),
				V(S(P(2, 1, 3), P(9, 1, 10)), token.REGEX_CONTENT, `foo\/bar`),
				T(S(P(10, 1, 11), P(10, 1, 11)), token.REGEX_END),
				T(S(P(11, 1, 12), P(11, 1, 12)), token.REGEX_FLAG_U),
				T(S(P(12, 1, 13), P(12, 1, 13)), token.REGEX_FLAG_i),
				T(S(P(13, 1, 14), P(13, 1, 14)), token.REGEX_FLAG_m),
				T(S(P(14, 1, 15), P(14, 1, 15)), token.REGEX_FLAG_a),
				T(S(P(15, 1, 16), P(15, 1, 16)), token.REGEX_FLAG_x),
				T(S(P(16, 1, 17), P(16, 1, 17)), token.REGEX_FLAG_s),
			},
		},
		"cannot have invalid flags": {
			input: `%/foo\/bar/UTs`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.REGEX_BEG),
				V(S(P(2, 1, 3), P(9, 1, 10)), token.REGEX_CONTENT, `foo\/bar`),
				T(S(P(10, 1, 11), P(10, 1, 11)), token.REGEX_END),
				T(S(P(11, 1, 12), P(11, 1, 12)), token.REGEX_FLAG_U),
				V(S(P(12, 1, 13), P(12, 1, 13)), token.ERROR, `invalid regex flag`),
				T(S(P(13, 1, 14), P(13, 1, 14)), token.REGEX_FLAG_s),
			},
		},
		"can have flags and operators without whitespace": {
			input: `%/foo\/bar/i+5`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.REGEX_BEG),
				V(S(P(2, 1, 3), P(9, 1, 10)), token.REGEX_CONTENT, `foo\/bar`),
				T(S(P(10, 1, 11), P(10, 1, 11)), token.REGEX_END),
				T(S(P(11, 1, 12), P(11, 1, 12)), token.REGEX_FLAG_i),
				T(S(P(12, 1, 13), P(12, 1, 13)), token.PLUS),
				V(S(P(13, 1, 14), P(13, 1, 14)), token.INT, `5`),
			},
		},
		"cannot have flags with whitespace": {
			input: `%/foo\/bar/ i`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.REGEX_BEG),
				V(S(P(2, 1, 3), P(9, 1, 10)), token.REGEX_CONTENT, `foo\/bar`),
				T(S(P(10, 1, 11), P(10, 1, 11)), token.REGEX_END),
				V(S(P(12, 1, 13), P(12, 1, 13)), token.PUBLIC_IDENTIFIER, `i`),
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
		"can be interpolated": {
			input: `%/foo${1 + 3.5}bar/i`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.REGEX_BEG),
				V(S(P(2, 1, 3), P(4, 1, 5)), token.REGEX_CONTENT, `foo`),
				T(S(P(5, 1, 6), P(6, 1, 7)), token.REGEX_INTERP_BEG),
				V(S(P(7, 1, 8), P(7, 1, 8)), token.INT, "1"),
				T(S(P(9, 1, 10), P(9, 1, 10)), token.PLUS),
				V(S(P(11, 1, 12), P(13, 1, 14)), token.FLOAT, "3.5"),
				T(S(P(14, 1, 15), P(14, 1, 15)), token.REGEX_INTERP_END),
				V(S(P(15, 1, 16), P(17, 1, 18)), token.REGEX_CONTENT, `bar`),
				T(S(P(18, 1, 19), P(18, 1, 19)), token.REGEX_END),
				T(S(P(19, 1, 20), P(19, 1, 20)), token.REGEX_FLAG_i),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
