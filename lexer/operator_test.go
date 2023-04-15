package lexer

import "testing"

func TestOperator(t *testing.T) {
	tests := testTable{
		"should be recognised": {
			input: ". .. ... - -= -> + += ^ ^= * *= ** **= = == === =~ => =:= =!= : := :: :> :>> ~ ~= ~> > >= >> >>= < <= << <<= <: <<: & &= && &&= | |= || ||= |> ? ?? ??= ! != !== % %=",
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
				T(PowerToken, 32, 2, 1, 33),
				T(PowerEqualToken, 35, 3, 1, 36),
				T(AssignToken, 39, 1, 1, 40),
				T(EqualToken, 41, 2, 1, 42),
				T(StrictEqualToken, 44, 3, 1, 45),
				T(MatchOpToken, 48, 2, 1, 49),
				T(ThickArrowToken, 51, 2, 1, 52),
				T(RefEqualToken, 54, 3, 1, 55),
				T(RefNotEqualToken, 58, 3, 1, 59),
				T(ColonToken, 62, 1, 1, 63),
				T(ColonEqualToken, 64, 2, 1, 65),
				T(ScopeResOpToken, 67, 2, 1, 68),
				T(ReverseSubtypeToken, 70, 2, 1, 71),
				T(ReverseInstanceOfToken, 73, 3, 1, 74),
				T(TildeToken, 77, 1, 1, 78),
				T(TildeEqualToken, 79, 2, 1, 80),
				T(WigglyArrowToken, 82, 2, 1, 83),
				T(GreaterToken, 85, 1, 1, 86),
				T(GreaterEqualToken, 87, 2, 1, 88),
				T(RBitShiftToken, 90, 2, 1, 91),
				T(RBitShiftEqualToken, 93, 3, 1, 94),
				T(LessToken, 97, 1, 1, 98),
				T(LessEqualToken, 99, 2, 1, 100),
				T(LBitShiftToken, 102, 2, 1, 103),
				T(LBitShiftEqualToken, 105, 3, 1, 106),
				T(SubtypeToken, 109, 2, 1, 110),
				T(InstanceOfToken, 112, 3, 1, 113),
				T(AndToken, 116, 1, 1, 117),
				T(AndEqualToken, 118, 2, 1, 119),
				T(AndAndToken, 121, 2, 1, 122),
				T(AndAndEqualToken, 124, 3, 1, 125),
				T(OrToken, 128, 1, 1, 129),
				T(OrEqualToken, 130, 2, 1, 131),
				T(OrOrToken, 133, 2, 1, 134),
				T(OrOrEqualToken, 136, 3, 1, 137),
				T(PipeOpToken, 140, 2, 1, 141),
				T(QuestionMarkToken, 143, 1, 1, 144),
				T(NilCoalesceToken, 145, 2, 1, 146),
				T(NilCoalesceEqualToken, 148, 3, 1, 149),
				T(BangToken, 152, 1, 1, 153),
				T(NotEqualToken, 154, 2, 1, 155),
				T(StrictNotEqualToken, 157, 3, 1, 158),
				T(PercentToken, 161, 1, 1, 162),
				T(PercentEqualToken, 163, 2, 1, 164),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
