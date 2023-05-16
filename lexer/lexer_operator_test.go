package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestOperator(t *testing.T) {
	tests := testTable{
		"should be recognised": {
			input: ". .. ... - -= -> + += ^ ^= * *= / /= ** **= = == === =~ => =:= =!= : := :: :> :>> ~ ~= ~> > >= >> >>= < <= << <<= <: <<: & &= && &&= | |= || ||= |> ? ?? ??= ! != !== % %= <=> &! |!",
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.DOT),
				T(P(2, 2, 1, 3), token.RANGE_OP),
				T(P(5, 3, 1, 6), token.EXCLUSIVE_RANGE_OP),
				T(P(9, 1, 1, 10), token.MINUS),
				T(P(11, 2, 1, 12), token.MINUS_EQUAL),
				T(P(14, 2, 1, 15), token.THIN_ARROW),
				T(P(17, 1, 1, 18), token.PLUS),
				T(P(19, 2, 1, 20), token.PLUS_EQUAL),
				T(P(22, 1, 1, 23), token.XOR),
				T(P(24, 2, 1, 25), token.XOR_EQUAL),
				T(P(27, 1, 1, 28), token.STAR),
				T(P(29, 2, 1, 30), token.STAR_EQUAL),
				T(P(32, 1, 1, 33), token.SLASH),
				T(P(34, 2, 1, 35), token.SLASH_EQUAL),
				T(P(37, 2, 1, 38), token.STAR_STAR),
				T(P(40, 3, 1, 41), token.STAR_STAR_EQUAL),
				T(P(44, 1, 1, 45), token.EQUAL_OP),
				T(P(46, 2, 1, 47), token.EQUAL_EQUAL),
				T(P(49, 3, 1, 50), token.STRICT_EQUAL),
				T(P(53, 2, 1, 54), token.MATCH_OP),
				T(P(56, 2, 1, 57), token.THICK_ARROW),
				T(P(59, 3, 1, 60), token.REF_EQUAL),
				T(P(63, 3, 1, 64), token.REF_NOT_EQUAL),
				T(P(67, 1, 1, 68), token.COLON),
				T(P(69, 2, 1, 70), token.COLON_EQUAL),
				T(P(72, 2, 1, 73), token.SCOPE_RES_OP),
				T(P(75, 2, 1, 76), token.REVERSE_ISA_OP),
				T(P(78, 3, 1, 79), token.REVERSE_INSTANCE_OF_OP),
				T(P(82, 1, 1, 83), token.TILDE),
				T(P(84, 2, 1, 85), token.TILDE_EQUAL),
				T(P(87, 2, 1, 88), token.WIGGLY_ARROW),
				T(P(90, 1, 1, 91), token.GREATER),
				T(P(92, 2, 1, 93), token.GREATER_EQUAL),
				T(P(95, 2, 1, 96), token.RBITSHIFT),
				T(P(98, 3, 1, 99), token.RBITSHIFT_EQUAL),
				T(P(102, 1, 1, 103), token.LESS),
				T(P(104, 2, 1, 105), token.LESS_EQUAL),
				T(P(107, 2, 1, 108), token.LBITSHIFT),
				T(P(110, 3, 1, 111), token.LBITSHIFT_EQUAL),
				T(P(114, 2, 1, 115), token.ISA_OP),
				T(P(117, 3, 1, 118), token.INSTANCE_OF_OP),
				T(P(121, 1, 1, 122), token.AND),
				T(P(123, 2, 1, 124), token.AND_EQUAL),
				T(P(126, 2, 1, 127), token.AND_AND),
				T(P(129, 3, 1, 130), token.AND_AND_EQUAL),
				T(P(133, 1, 1, 134), token.OR),
				T(P(135, 2, 1, 136), token.OR_EQUAL),
				T(P(138, 2, 1, 139), token.OR_OR),
				T(P(141, 3, 1, 142), token.OR_OR_EQUAL),
				T(P(145, 2, 1, 146), token.PIPE_OP),
				T(P(148, 1, 1, 149), token.QUESTION),
				T(P(150, 2, 1, 151), token.QUESTION_QUESTION),
				T(P(153, 3, 1, 154), token.QUESTION_QUESTION_EQUAL),
				T(P(157, 1, 1, 158), token.BANG),
				T(P(159, 2, 1, 160), token.NOT_EQUAL),
				T(P(162, 3, 1, 163), token.STRICT_NOT_EQUAL),
				T(P(166, 1, 1, 167), token.PERCENT),
				T(P(168, 2, 1, 169), token.PERCENT_EQUAL),
				T(P(171, 3, 1, 172), token.SPACESHIP_OP),
				T(P(175, 2, 1, 176), token.AND_BANG),
				T(P(178, 2, 1, 179), token.OR_BANG),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
