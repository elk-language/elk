package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestOperator(t *testing.T) {
	tests := testTable{
		"should be recognised": {
			input: ". .. ... - -= -> + += ^ ^= * *= / /= ** **= = == === =~ !~ => : := :: :> :>> ~ ~= ~> > >= >> >>= < <= << <<= <: <<: & &= && &&= | |= || ||= |> ? ?? ??= ! != !== % %= <=> &! |! <<< <<<= >>> >>>= ?. ++ -- &~",
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
				T(S(P(53, 1, 54), P(54, 1, 55)), token.LAX_EQUAL),
				T(S(P(56, 1, 57), P(57, 1, 58)), token.LAX_NOT_EQUAL),
				T(S(P(59, 1, 60), P(60, 1, 61)), token.THICK_ARROW),
				T(S(P(62, 1, 63), P(62, 1, 63)), token.COLON),
				T(S(P(64, 1, 65), P(65, 1, 66)), token.COLON_EQUAL),
				T(S(P(67, 1, 68), P(68, 1, 69)), token.SCOPE_RES_OP),
				T(S(P(70, 1, 71), P(71, 1, 72)), token.REVERSE_ISA_OP),
				T(S(P(73, 1, 74), P(75, 1, 76)), token.REVERSE_INSTANCE_OF_OP),
				T(S(P(77, 1, 78), P(77, 1, 78)), token.TILDE),
				T(S(P(79, 1, 80), P(80, 1, 81)), token.TILDE_EQUAL),
				T(S(P(82, 1, 83), P(83, 1, 84)), token.WIGGLY_ARROW),
				T(S(P(85, 1, 86), P(85, 1, 86)), token.GREATER),
				T(S(P(87, 1, 88), P(88, 1, 89)), token.GREATER_EQUAL),
				T(S(P(90, 1, 91), P(91, 1, 92)), token.RBITSHIFT),
				T(S(P(93, 1, 94), P(95, 1, 96)), token.RBITSHIFT_EQUAL),
				T(S(P(97, 1, 98), P(97, 1, 98)), token.LESS),
				T(S(P(99, 1, 100), P(100, 1, 101)), token.LESS_EQUAL),
				T(S(P(102, 1, 103), P(103, 1, 104)), token.LBITSHIFT),
				T(S(P(105, 1, 106), P(107, 1, 108)), token.LBITSHIFT_EQUAL),
				T(S(P(109, 1, 110), P(110, 1, 111)), token.ISA_OP),
				T(S(P(112, 1, 113), P(114, 1, 115)), token.INSTANCE_OF_OP),
				T(S(P(116, 1, 117), P(116, 1, 117)), token.AND),
				T(S(P(118, 1, 119), P(119, 1, 120)), token.AND_EQUAL),
				T(S(P(121, 1, 122), P(122, 1, 123)), token.AND_AND),
				T(S(P(124, 1, 125), P(126, 1, 127)), token.AND_AND_EQUAL),
				T(S(P(128, 1, 129), P(128, 1, 129)), token.OR),
				T(S(P(130, 1, 131), P(131, 1, 132)), token.OR_EQUAL),
				T(S(P(133, 1, 134), P(134, 1, 135)), token.OR_OR),
				T(S(P(136, 1, 137), P(138, 1, 139)), token.OR_OR_EQUAL),
				T(S(P(140, 1, 141), P(141, 1, 142)), token.PIPE_OP),
				T(S(P(143, 1, 144), P(143, 1, 144)), token.QUESTION),
				T(S(P(145, 1, 146), P(146, 1, 147)), token.QUESTION_QUESTION),
				T(S(P(148, 1, 149), P(150, 1, 151)), token.QUESTION_QUESTION_EQUAL),
				T(S(P(152, 1, 153), P(152, 1, 153)), token.BANG),
				T(S(P(154, 1, 155), P(155, 1, 156)), token.NOT_EQUAL),
				T(S(P(157, 1, 158), P(159, 1, 160)), token.STRICT_NOT_EQUAL),
				T(S(P(161, 1, 162), P(161, 1, 162)), token.PERCENT),
				T(S(P(163, 1, 164), P(164, 1, 165)), token.PERCENT_EQUAL),
				T(S(P(166, 1, 167), P(168, 1, 169)), token.SPACESHIP_OP),
				T(S(P(170, 1, 171), P(171, 1, 172)), token.AND_BANG),
				T(S(P(173, 1, 174), P(174, 1, 175)), token.OR_BANG),
				T(S(P(176, 1, 177), P(178, 1, 179)), token.LTRIPLE_BITSHIFT),
				T(S(P(180, 1, 181), P(183, 1, 184)), token.LTRIPLE_BITSHIFT_EQUAL),
				T(S(P(185, 1, 186), P(187, 1, 188)), token.RTRIPLE_BITSHIFT),
				T(S(P(189, 1, 190), P(192, 1, 193)), token.RTRIPLE_BITSHIFT_EQUAL),
				T(S(P(194, 1, 195), P(195, 1, 196)), token.QUESTION_DOT),
				T(S(P(197, 1, 198), P(198, 1, 199)), token.PLUS_PLUS),
				T(S(P(200, 1, 201), P(201, 1, 202)), token.MINUS_MINUS),
				T(S(P(203, 1, 204), P(204, 1, 205)), token.AND_TILDE),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
