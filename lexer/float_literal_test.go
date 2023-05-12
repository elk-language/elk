package lexer

import "testing"

func TestFloat(t *testing.T) {
	tests := testTable{
		"with underscores": {
			input: "245_000.254_129",
			want: []*Token{
				V(P(0, 15, 1, 1), FloatToken, "245000.254129"),
			},
		},
		"ends on last valid character": {
			input: "0.36f",
			want: []*Token{
				V(P(0, 4, 1, 1), FloatToken, "0.36"),
				V(P(4, 1, 1, 5), PublicIdentifierToken, "f"),
			},
		},
		"can only be decimal": {
			input: "0x21.36",
			want: []*Token{
				V(P(0, 4, 1, 1), HexIntToken, "21"),
				V(P(4, 3, 1, 5), FloatToken, "0.36"),
			},
		},
		"with exponent": {
			input: "0.36e2",
			want: []*Token{
				V(P(0, 6, 1, 1), FloatToken, "0.36e2"),
			},
		},
		"with exponent and no dot": {
			input: "25e4",
			want: []*Token{
				V(P(0, 4, 1, 1), FloatToken, "25e4"),
			},
		},
		"with explicit positive exponent and no dot": {
			input: "25e+4",
			want: []*Token{
				V(P(0, 5, 1, 1), FloatToken, "25e+4"),
			},
		},
		"with uppercase exponent": {
			input: "0.36E2",
			want: []*Token{
				V(P(0, 6, 1, 1), FloatToken, "0.36e2"),
			},
		},
		"with negative exponent": {
			input: "25.8e-36",
			want: []*Token{
				V(P(0, 8, 1, 1), FloatToken, "25.8e-36"),
			},
		},
		"without leading zero": {
			input: ".908267374623",
			want: []*Token{
				V(P(0, 13, 1, 1), FloatToken, "0.908267374623"),
			},
		},
		"without leading zero and with exponent": {
			input: ".8e-36",
			want: []*Token{
				V(P(0, 6, 1, 1), FloatToken, "0.8e-36"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
