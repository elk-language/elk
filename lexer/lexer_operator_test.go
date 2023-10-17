package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestOperator(t *testing.T) {
	tests := testTable{
		"should be recognised": {
			input: ". .. ... - -= -> + += ^ ^= * *= / /= ** **= = == === =~ => =:= =!= : := :: :> :>> ~ ~= ~> > >= >> >>= < <= << <<= <: <<: & &= && &&= | |= || ||= |> ? ?? ??= ! != !== % %= <=> &! |! <<< <<<= >>> >>>= ?. ++ --",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(0, 1, 1)), token.DOT),
				T(S(P(2, 1, 3), P(3, 1, 4)), token.RANGE_OP),
				T(S(P(5, 1, 6), P(7, 1, 8)), token.EXCLUSIVE_RANGE_OP),
				T(S(P(9, 1, 10), P(9, 1, 10)), token.MINUS),
				T(S(P(11, 1, 12), P(12, 1, 13)), token.MINUS_EQUAL),
				T(S(P(14, 1, 15), P(15, 1, 16)), token.THIN_ARROW),
				T(S(P(17, 1, 18), P(17, 1, 18)), token.PLUS),
				T(S(P(19, 1, 20), P(20, 1, 21)), token.PLUS_EQUAL),
				T(S(P(22, 1, 23), P(22, 1, 23)), token.XOR),
				T(S(P(24, 1, 25), P(25, 1, 26)), token.XOR_EQUAL),
				T(S(P(27, 1, 28), P(27, 1, 28)), token.STAR),
				T(S(P(29, 1, 30), P(30, 1, 31)), token.STAR_EQUAL),
				T(S(P(32, 1, 33), P(32, 1, 33)), token.SLASH),
				T(S(P(34, 1, 35), P(35, 1, 36)), token.SLASH_EQUAL),
				T(S(P(37, 1, 38), P(38, 1, 39)), token.STAR_STAR),
				T(S(P(40, 1, 41), P(42, 1, 43)), token.STAR_STAR_EQUAL),
				T(S(P(44, 1, 45), P(44, 1, 45)), token.EQUAL_OP),
				T(S(P(46, 1, 47), P(47, 1, 48)), token.EQUAL_EQUAL),
				T(S(P(49, 1, 50), P(51, 1, 52)), token.STRICT_EQUAL),
				T(S(P(53, 1, 54), P(54, 1, 55)), token.MATCH_OP),
				T(S(P(56, 1, 57), P(57, 1, 58)), token.THICK_ARROW),
				T(S(P(59, 1, 60), P(61, 1, 62)), token.REF_EQUAL),
				T(S(P(63, 1, 64), P(65, 1, 66)), token.REF_NOT_EQUAL),
				T(S(P(67, 1, 68), P(67, 1, 68)), token.COLON),
				T(S(P(69, 1, 70), P(70, 1, 71)), token.COLON_EQUAL),
				T(S(P(72, 1, 73), P(73, 1, 74)), token.SCOPE_RES_OP),
				T(S(P(75, 1, 76), P(76, 1, 77)), token.REVERSE_ISA_OP),
				T(S(P(78, 1, 79), P(80, 1, 81)), token.REVERSE_INSTANCE_OF_OP),
				T(S(P(82, 1, 83), P(82, 1, 83)), token.TILDE),
				T(S(P(84, 1, 85), P(85, 1, 86)), token.TILDE_EQUAL),
				T(S(P(87, 1, 88), P(88, 1, 89)), token.WIGGLY_ARROW),
				T(S(P(90, 1, 91), P(90, 1, 91)), token.GREATER),
				T(S(P(92, 1, 93), P(93, 1, 94)), token.GREATER_EQUAL),
				T(S(P(95, 1, 96), P(96, 1, 97)), token.RBITSHIFT),
				T(S(P(98, 1, 99), P(100, 1, 101)), token.RBITSHIFT_EQUAL),
				T(S(P(102, 1, 103), P(102, 1, 103)), token.LESS),
				T(S(P(104, 1, 105), P(105, 1, 106)), token.LESS_EQUAL),
				T(S(P(107, 1, 108), P(108, 1, 109)), token.LBITSHIFT),
				T(S(P(110, 1, 111), P(112, 1, 113)), token.LBITSHIFT_EQUAL),
				T(S(P(114, 1, 115), P(115, 1, 116)), token.ISA_OP),
				T(S(P(117, 1, 118), P(119, 1, 120)), token.INSTANCE_OF_OP),
				T(S(P(121, 1, 122), P(121, 1, 122)), token.AND),
				T(S(P(123, 1, 124), P(124, 1, 125)), token.AND_EQUAL),
				T(S(P(126, 1, 127), P(127, 1, 128)), token.AND_AND),
				T(S(P(129, 1, 130), P(131, 1, 132)), token.AND_AND_EQUAL),
				T(S(P(133, 1, 134), P(133, 1, 134)), token.OR),
				T(S(P(135, 1, 136), P(136, 1, 137)), token.OR_EQUAL),
				T(S(P(138, 1, 139), P(139, 1, 140)), token.OR_OR),
				T(S(P(141, 1, 142), P(143, 1, 144)), token.OR_OR_EQUAL),
				T(S(P(145, 1, 146), P(146, 1, 147)), token.PIPE_OP),
				T(S(P(148, 1, 149), P(148, 1, 149)), token.QUESTION),
				T(S(P(150, 1, 151), P(151, 1, 152)), token.QUESTION_QUESTION),
				T(S(P(153, 1, 154), P(155, 1, 156)), token.QUESTION_QUESTION_EQUAL),
				T(S(P(157, 1, 158), P(157, 1, 158)), token.BANG),
				T(S(P(159, 1, 160), P(160, 1, 161)), token.NOT_EQUAL),
				T(S(P(162, 1, 163), P(164, 1, 165)), token.STRICT_NOT_EQUAL),
				T(S(P(166, 1, 167), P(166, 1, 167)), token.PERCENT),
				T(S(P(168, 1, 169), P(169, 1, 170)), token.PERCENT_EQUAL),
				T(S(P(171, 1, 172), P(173, 1, 174)), token.SPACESHIP_OP),
				T(S(P(175, 1, 176), P(176, 1, 177)), token.AND_BANG),
				T(S(P(178, 1, 179), P(179, 1, 180)), token.OR_BANG),
				T(S(P(181, 1, 182), P(183, 1, 184)), token.LTRIPLE_BITSHIFT),
				T(S(P(185, 1, 186), P(188, 1, 189)), token.LTRIPLE_BITSHIFT_EQUAL),
				T(S(P(190, 1, 191), P(192, 1, 193)), token.RTRIPLE_BITSHIFT),
				T(S(P(194, 1, 195), P(197, 1, 198)), token.RTRIPLE_BITSHIFT_EQUAL),
				T(S(P(199, 1, 200), P(200, 1, 201)), token.QUESTION_DOT),
				T(S(P(202, 1, 203), P(203, 1, 204)), token.PLUS_PLUS),
				T(S(P(205, 1, 206), P(206, 1, 207)), token.MINUS_MINUS),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
