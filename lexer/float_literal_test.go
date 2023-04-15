package lexer

import "testing"

func TestFloat(t *testing.T) {
	tests := testTable{
		"with underscores": {
			input: "245_000.254_129",
			want: []*Token{
				V(FloatToken, "245000.254129", 0, 15, 1, 1),
			},
		},
		"ends on last valid character": {
			input: "0.36f",
			want: []*Token{
				V(FloatToken, "0.36", 0, 4, 1, 1),
				V(IdentifierToken, "f", 4, 1, 1, 5),
			},
		},
		"can only be decimal": {
			input: "0x21.36",
			want: []*Token{
				V(HexIntToken, "21", 0, 4, 1, 1),
				V(FloatToken, "0.36", 4, 3, 1, 5),
			},
		},
		"with exponent": {
			input: "0.36e2",
			want: []*Token{
				V(FloatToken, "0.36e2", 0, 6, 1, 1),
			},
		},
		"with exponent and no dot": {
			input: "25e4",
			want: []*Token{
				V(FloatToken, "25e4", 0, 4, 1, 1),
			},
		},
		"with explicit positive exponent and no dot": {
			input: "25e+4",
			want: []*Token{
				V(FloatToken, "25e+4", 0, 5, 1, 1),
			},
		},
		"with uppercase exponent": {
			input: "0.36E2",
			want: []*Token{
				V(FloatToken, "0.36e2", 0, 6, 1, 1),
			},
		},
		"with negative exponent": {
			input: "25.8e-36",
			want: []*Token{
				V(FloatToken, "25.8e-36", 0, 8, 1, 1),
			},
		},
		"without leading zero": {
			input: ".908267374623",
			want: []*Token{
				V(FloatToken, "0.908267374623", 0, 13, 1, 1),
			},
		},
		"without leading zero and with exponent": {
			input: ".8e-36",
			want: []*Token{
				V(FloatToken, "0.8e-36", 0, 6, 1, 1),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
