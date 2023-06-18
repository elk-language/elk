package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestInt(t *testing.T) {
	tests := testTable{
		"decimal": {
			input: "23",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.INT, "23"),
			},
		},
		"decimal with leading zeros": {
			input: "00015",
			want: []*token.Token{
				V(P(0, 5, 1, 1), token.INT, "00015"),
			},
		},
		"decimal with underscores": {
			input: "23_200_123",
			want: []*token.Token{
				V(P(0, 10, 1, 1), token.INT, "23200123"),
			},
		},
		"decimal can't begin with underscores": {
			input: "_23_200_123",
			want: []*token.Token{
				V(P(0, 1, 1, 1), token.PRIVATE_IDENTIFIER, "_"),
				V(P(1, 10, 1, 2), token.INT, "23200123"),
			},
		},
		"decimal ends on last valid character": {
			input: "23bar",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.INT, "23"),
				V(P(2, 3, 1, 3), token.PUBLIC_IDENTIFIER, "bar"),
			},
		},
		"decimal int64": {
			input: "23i64",
			want: []*token.Token{
				V(P(0, 5, 1, 1), token.INT64, "23"),
			},
		},
		"hex": {
			input: "0x354ab1",
			want: []*token.Token{
				V(P(0, 8, 1, 1), token.INT, "0x354ab1"),
			},
		},
		"hex with underscores": {
			input: "0x35_4a_b1",
			want: []*token.Token{
				V(P(0, 10, 1, 1), token.INT, "0x354ab1"),
			},
		},
		"leading zeros invalidate hex": {
			input: "00x21",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.INT, "00"),
				V(P(2, 3, 1, 3), token.PUBLIC_IDENTIFIER, "x21"),
			},
		},
		"hex without digits": {
			input: "0x",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.INT, "0x"),
			},
		},
		"hex with uppercase": {
			input: "0X354Ab1",
			want: []*token.Token{
				V(P(0, 8, 1, 1), token.INT, "0x354Ab1"),
			},
		},
		"hex ends on last valid character": {
			input: "0x354fp",
			want: []*token.Token{
				V(P(0, 6, 1, 1), token.INT, "0x354f"),
				V(P(6, 1, 1, 7), token.PUBLIC_IDENTIFIER, "p"),
			},
		},
		"octal": {
			input: "0o603",
			want: []*token.Token{
				V(P(0, 5, 1, 1), token.INT, "0o603"),
			},
		},
		"octal with underscores": {
			input: "0o3201_5200",
			want: []*token.Token{
				V(P(0, 11, 1, 1), token.INT, "0o32015200"),
			},
		},
		"leading zeros invalidate octal": {
			input: "00o21",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.INT, "00"),
				V(P(2, 3, 1, 3), token.PUBLIC_IDENTIFIER, "o21"),
			},
		},
		"octal without digits": {
			input: "0o",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.INT, "0o"),
			},
		},
		"octal with uppercase": {
			input: "0O603",
			want: []*token.Token{
				V(P(0, 5, 1, 1), token.INT, "0o603"),
			},
		},
		"octal ends on last valid character": {
			input: "0o6039abc1",
			want: []*token.Token{
				V(P(0, 5, 1, 1), token.INT, "0o603"),
				V(P(5, 1, 1, 6), token.INT, "9"),
				V(P(6, 4, 1, 7), token.PUBLIC_IDENTIFIER, "abc1"),
			},
		},
		"quaternary": {
			input: "0q30212",
			want: []*token.Token{
				V(P(0, 7, 1, 1), token.INT, "0q30212"),
			},
		},
		"quaternary with underscores": {
			input: "0q3201200_23010000",
			want: []*token.Token{
				V(P(0, 18, 1, 1), token.INT, "0q320120023010000"),
			},
		},
		"leading zeros invalidate quaternary": {
			input: "00q21",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.INT, "00"),
				V(P(2, 3, 1, 3), token.PUBLIC_IDENTIFIER, "q21"),
			},
		},
		"quaternary with uppercase": {
			input: "0Q30212",
			want: []*token.Token{
				V(P(0, 7, 1, 1), token.INT, "0q30212"),
			},
		},
		"quaternary without digits": {
			input: "0q",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.INT, "0q"),
			},
		},
		"quaternary ends on last valid character": {
			input: "0q302124a",
			want: []*token.Token{
				V(P(0, 7, 1, 1), token.INT, "0q30212"),
				V(P(7, 1, 1, 8), token.INT, "4"),
				V(P(8, 1, 1, 9), token.PUBLIC_IDENTIFIER, "a"),
			},
		},
		"binary": {
			input: "0b1010",
			want: []*token.Token{
				V(P(0, 6, 1, 1), token.INT, "0b1010"),
			},
		},
		"binary with underscores": {
			input: "0b10_10_10",
			want: []*token.Token{
				V(P(0, 10, 1, 1), token.INT, "0b101010"),
			},
		},
		"leading zeros invalidate binary": {
			input: "00b21",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.INT, "00"),
				V(P(2, 3, 1, 3), token.PUBLIC_IDENTIFIER, "b21"),
			},
		},
		"binary with uppercase": {
			input: "0B1010",
			want: []*token.Token{
				V(P(0, 6, 1, 1), token.INT, "0b1010"),
			},
		},
		"binary ends on last valid character": {
			input: "0b10102dup",
			want: []*token.Token{
				V(P(0, 6, 1, 1), token.INT, "0b1010"),
				V(P(6, 1, 1, 7), token.INT, "2"),
				V(P(7, 3, 1, 8), token.PUBLIC_IDENTIFIER, "dup"),
			},
		},
		"binary without digits": {
			input: "0b",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.INT, "0b"),
			},
		},
		"duodecimal": {
			input: "0da12b3",
			want: []*token.Token{
				V(P(0, 7, 1, 1), token.INT, "0da12b3"),
			},
		},
		"duodecimal with uppercase": {
			input: "0Da12B3",
			want: []*token.Token{
				V(P(0, 7, 1, 1), token.INT, "0da12B3"),
			},
		},
		"duodecimal ends on last valid character": {
			input: "0d23a3bca3",
			want: []*token.Token{
				V(P(0, 7, 1, 1), token.INT, "0d23a3b"),
				V(P(7, 3, 1, 8), token.PUBLIC_IDENTIFIER, "ca3"),
			},
		},
		"duodecimal without digits": {
			input: "0d",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.INT, "0d"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
