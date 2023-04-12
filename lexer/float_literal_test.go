package lexer

import "testing"

func TestFloat(t *testing.T) {
	tests := testTable{
		"with underscores": {
			input: "245_000.254_129",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "245_000.254_129",
					StartByte:  0,
					ByteLength: 15,
					Line:       1,
					Column:     1,
				},
			},
		},
		"ends on last valid character": {
			input: "0.36f",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "0.36",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "f",
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
				},
			},
		},
		"can only be decimal": {
			input: "0x21.36",
			want: []*Token{
				{
					TokenType:  IntToken,
					Value:      "0x21",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  FloatToken,
					Value:      ".36",
					StartByte:  4,
					ByteLength: 3,
					Line:       1,
					Column:     5,
				},
			},
		},
		"with exponent": {
			input: "0.36e2",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "0.36e2",
					StartByte:  0,
					ByteLength: 6,
					Line:       1,
					Column:     1,
				},
			},
		},
		"with exponent and no dot": {
			input: "25e4",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "25e4",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
			},
		},
		"with explicit positive exponent and no dot": {
			input: "25e+4",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "25e+4",
					StartByte:  0,
					ByteLength: 5,
					Line:       1,
					Column:     1,
				},
			},
		},
		"with uppercase exponent": {
			input: "0.36E2",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "0.36E2",
					StartByte:  0,
					ByteLength: 6,
					Line:       1,
					Column:     1,
				},
			},
		},
		"with negative exponent": {
			input: "25.8e-36",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "25.8e-36",
					StartByte:  0,
					ByteLength: 8,
					Line:       1,
					Column:     1,
				},
			},
		},
		"without leading zero": {
			input: ".908267374623",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      ".908267374623",
					StartByte:  0,
					ByteLength: 13,
					Line:       1,
					Column:     1,
				},
			},
		},
		"without leading zero and with exponent": {
			input: ".8e-36",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      ".8e-36",
					StartByte:  0,
					ByteLength: 6,
					Line:       1,
					Column:     1,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
