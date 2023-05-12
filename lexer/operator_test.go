package lexer

import "testing"

func TestOperator(t *testing.T) {
	tests := testTable{
		"should be recognised": {
			input: ". .. ... - -= -> + += ^ ^= * *= / /= ** **= = == === =~ => =:= =!= : := :: :> :>> ~ ~= ~> > >= >> >>= < <= << <<= <: <<: & &= && &&= | |= || ||= |> ? ?? ??= ! != !== % %= <=>",
			want: []*Token{
				T(P(0, 1, 1, 1), DotToken),
				T(P(2, 2, 1, 3), RangeOpToken),
				T(P(5, 3, 1, 6), ExclusiveRangeOpToken),
				T(P(9, 1, 1, 10), MinusToken),
				T(P(11, 2, 1, 12), MinusEqualToken),
				T(P(14, 2, 1, 15), ThinArrowToken),
				T(P(17, 1, 1, 18), PlusToken),
				T(P(19, 2, 1, 20), PlusEqualToken),
				T(P(22, 1, 1, 23), XorToken),
				T(P(24, 2, 1, 25), XorEqualToken),
				T(P(27, 1, 1, 28), StarToken),
				T(P(29, 2, 1, 30), StarEqualToken),
				T(P(32, 1, 1, 33), SlashToken),
				T(P(34, 2, 1, 35), SlashEqualToken),
				T(P(37, 2, 1, 38), StarStarToken),
				T(P(40, 3, 1, 41), StarStarEqualToken),
				T(P(44, 1, 1, 45), EqualToken),
				T(P(46, 2, 1, 47), EqualEqualToken),
				T(P(49, 3, 1, 50), StrictEqualToken),
				T(P(53, 2, 1, 54), MatchOpToken),
				T(P(56, 2, 1, 57), ThickArrowToken),
				T(P(59, 3, 1, 60), RefEqualToken),
				T(P(63, 3, 1, 64), RefNotEqualToken),
				T(P(67, 1, 1, 68), ColonToken),
				T(P(69, 2, 1, 70), ColonEqualToken),
				T(P(72, 2, 1, 73), ScopeResOpToken),
				T(P(75, 2, 1, 76), ReverseSubtypeToken),
				T(P(78, 3, 1, 79), ReverseInstanceOfToken),
				T(P(82, 1, 1, 83), TildeToken),
				T(P(84, 2, 1, 85), TildeEqualToken),
				T(P(87, 2, 1, 88), WigglyArrowToken),
				T(P(90, 1, 1, 91), GreaterToken),
				T(P(92, 2, 1, 93), GreaterEqualToken),
				T(P(95, 2, 1, 96), RBitShiftToken),
				T(P(98, 3, 1, 99), RBitShiftEqualToken),
				T(P(102, 1, 1, 103), LessToken),
				T(P(104, 2, 1, 105), LessEqualToken),
				T(P(107, 2, 1, 108), LBitShiftToken),
				T(P(110, 3, 1, 111), LBitShiftEqualToken),
				T(P(114, 2, 1, 115), SubtypeToken),
				T(P(117, 3, 1, 118), InstanceOfToken),
				T(P(121, 1, 1, 122), AndToken),
				T(P(123, 2, 1, 124), AndEqualToken),
				T(P(126, 2, 1, 127), AndAndToken),
				T(P(129, 3, 1, 130), AndAndEqualToken),
				T(P(133, 1, 1, 134), OrToken),
				T(P(135, 2, 1, 136), OrEqualToken),
				T(P(138, 2, 1, 139), OrOrToken),
				T(P(141, 3, 1, 142), OrOrEqualToken),
				T(P(145, 2, 1, 146), PipeOpToken),
				T(P(148, 1, 1, 149), QuestionMarkToken),
				T(P(150, 2, 1, 151), QuestionQuestionToken),
				T(P(153, 3, 1, 154), QuestionQuestionEqualToken),
				T(P(157, 1, 1, 158), BangToken),
				T(P(159, 2, 1, 160), NotEqualToken),
				T(P(162, 3, 1, 163), StrictNotEqualToken),
				T(P(166, 1, 1, 167), PercentToken),
				T(P(168, 2, 1, 169), PercentEqualToken),
				T(P(171, 3, 1, 172), SpaceshipOpToken),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
