package lexer

import "testing"

func TestOperator(t *testing.T) {
	tests := testTable{
		"should be recognised": {
			input: ". .. ... - -= -> + += ^ ^= * *= / /= ** **= = == === =~ => =:= =!= : := :: :> :>> ~ ~= ~> > >= >> >>= < <= << <<= <: <<: & &= && &&= | |= || ||= |> ? ?? ??= ! != !== % %=",
			want: []*Token{
				T(DotToken, 0, 1, 1, 1),
				T(RangeOpToken, 2, 2, 1, 3),
				T(ExclusiveRangeOpToken, 5, 3, 1, 6),
				T(MinusToken, 9, 1, 1, 10),
				T(MinusEqualToken, 11, 2, 1, 12),
				T(ThinArrowToken, 14, 2, 1, 15),
				T(PlusToken, 17, 1, 1, 18),
				T(PlusEqualToken, 19, 2, 1, 20),
				T(XorToken, 22, 1, 1, 23),
				T(XorEqualToken, 24, 2, 1, 25),
				T(StarToken, 27, 1, 1, 28),
				T(StarEqualToken, 29, 2, 1, 30),
				T(SlashToken, 32, 1, 1, 33),
				T(SlashEqualToken, 34, 2, 1, 35),
				T(PowerToken, 37, 2, 1, 38),
				T(PowerEqualToken, 40, 3, 1, 41),
				T(AssignToken, 44, 1, 1, 45),
				T(EqualToken, 46, 2, 1, 47),
				T(StrictEqualToken, 49, 3, 1, 50),
				T(MatchOpToken, 53, 2, 1, 54),
				T(ThickArrowToken, 56, 2, 1, 57),
				T(RefEqualToken, 59, 3, 1, 60),
				T(RefNotEqualToken, 63, 3, 1, 64),
				T(ColonToken, 67, 1, 1, 68),
				T(ColonEqualToken, 69, 2, 1, 70),
				T(ScopeResOpToken, 72, 2, 1, 73),
				T(ReverseSubtypeToken, 75, 2, 1, 76),
				T(ReverseInstanceOfToken, 78, 3, 1, 79),
				T(TildeToken, 82, 1, 1, 83),
				T(TildeEqualToken, 84, 2, 1, 85),
				T(WigglyArrowToken, 87, 2, 1, 88),
				T(GreaterToken, 90, 1, 1, 91),
				T(GreaterEqualToken, 92, 2, 1, 93),
				T(RBitShiftToken, 95, 2, 1, 96),
				T(RBitShiftEqualToken, 98, 3, 1, 99),
				T(LessToken, 102, 1, 1, 103),
				T(LessEqualToken, 104, 2, 1, 105),
				T(LBitShiftToken, 107, 2, 1, 108),
				T(LBitShiftEqualToken, 110, 3, 1, 111),
				T(SubtypeToken, 114, 2, 1, 115),
				T(InstanceOfToken, 117, 3, 1, 118),
				T(AndToken, 121, 1, 1, 122),
				T(AndEqualToken, 123, 2, 1, 124),
				T(AndAndToken, 126, 2, 1, 127),
				T(AndAndEqualToken, 129, 3, 1, 130),
				T(OrToken, 133, 1, 1, 134),
				T(OrEqualToken, 135, 2, 1, 136),
				T(OrOrToken, 138, 2, 1, 139),
				T(OrOrEqualToken, 141, 3, 1, 142),
				T(PipeOpToken, 145, 2, 1, 146),
				T(QuestionMarkToken, 148, 1, 1, 149),
				T(NilCoalesceToken, 150, 2, 1, 151),
				T(NilCoalesceEqualToken, 153, 3, 1, 154),
				T(BangToken, 157, 1, 1, 158),
				T(NotEqualToken, 159, 2, 1, 160),
				T(StrictNotEqualToken, 162, 3, 1, 163),
				T(PercentToken, 166, 1, 1, 167),
				T(PercentEqualToken, 168, 2, 1, 169),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
