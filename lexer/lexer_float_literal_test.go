package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestFloat(t *testing.T) {
	tests := testTable{
		"with underscores": {
			input: "245_000.254_129",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(14, 1, 15)), token.FLOAT, "245000.254129"),
			},
		},
		"ends on last valid character": {
			input: "0.36p",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(3, 1, 4)), token.FLOAT, "0.36"),
				V(S(P(4, 1, 5), P(4, 1, 5)), token.PUBLIC_IDENTIFIER, "p"),
			},
		},
		"can only be decimal": {
			input: "0x21.36",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(3, 1, 4)), token.INT, "0x21"),
				V(S(P(4, 1, 5), P(6, 1, 7)), token.FLOAT, "0.36"),
			},
		},
		"with exponent": {
			input: "0.36e2",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(5, 1, 6)), token.FLOAT, "0.36e2"),
			},
		},
		"with exponent and no dot": {
			input: "25e4",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(3, 1, 4)), token.FLOAT, "25e4"),
			},
		},
		"with explicit positive exponent and no dot": {
			input: "25e+4",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(4, 1, 5)), token.FLOAT, "25e+4"),
			},
		},
		"with uppercase exponent": {
			input: "0.36E2",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(5, 1, 6)), token.FLOAT, "0.36e2"),
			},
		},
		"with negative exponent": {
			input: "25.8e-36",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(7, 1, 8)), token.FLOAT, "25.8e-36"),
			},
		},
		"without leading zero": {
			input: ".908267374623",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(12, 1, 13)), token.FLOAT, "0.908267374623"),
			},
		},
		"without leading zero and with exponent": {
			input: ".8e-36",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(5, 1, 6)), token.FLOAT, "0.8e-36"),
			},
		},
		"with trailing dot": {
			input: "1.",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(0, 1, 1)), token.INT, "1"),
				T(S(P(1, 1, 2), P(1, 1, 2)), token.DOT),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestBigFloat(t *testing.T) {
	tests := testTable{
		"with underscores": {
			input: "245_000.254_129bf",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(16, 1, 17)), token.BIG_FLOAT, "245000.254129"),
			},
		},
		"can only be decimal": {
			input: "0x21.36bf",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(3, 1, 4)), token.INT, "0x21"),
				V(S(P(4, 1, 5), P(8, 1, 9)), token.BIG_FLOAT, "0.36"),
			},
		},
		"with exponent": {
			input: "0.36e2_bf",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(8, 1, 9)), token.BIG_FLOAT, "0.36e2"),
			},
		},
		"with exponent and no dot": {
			input: "25e4_bf",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(6, 1, 7)), token.BIG_FLOAT, "25e4"),
			},
		},
		"with explicit positive exponent and no dot": {
			input: "25e+4_bf",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(7, 1, 8)), token.BIG_FLOAT, "25e+4"),
			},
		},
		"with uppercase exponent": {
			input: "0.36E2_bf",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(8, 1, 9)), token.BIG_FLOAT, "0.36e2"),
			},
		},
		"with negative exponent": {
			input: "25.8e-36_bf",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(10, 1, 11)), token.BIG_FLOAT, "25.8e-36"),
			},
		},
		"without leading zero": {
			input: ".908267374623bf",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(14, 1, 15)), token.BIG_FLOAT, "0.908267374623"),
			},
		},
		"without leading zero and with exponent": {
			input: ".8e-36_bf",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(8, 1, 9)), token.BIG_FLOAT, "0.8e-36"),
			},
		},
		"with trailing dot": {
			input: "1.",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(0, 1, 1)), token.INT, "1"),
				T(S(P(1, 1, 2), P(1, 1, 2)), token.DOT),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
