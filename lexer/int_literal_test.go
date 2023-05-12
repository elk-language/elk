package lexer

import "testing"

func TestInt(t *testing.T) {
	tests := testTable{
		"decimal": {
			input: "23",
			want: []*Token{
				V(P(0, 2, 1, 1), DecIntToken, "23"),
			},
		},
		"decimal with leading zeros": {
			input: "00015",
			want: []*Token{
				V(P(0, 5, 1, 1), DecIntToken, "00015"),
			},
		},
		"decimal with underscores": {
			input: "23_200_123",
			want: []*Token{
				V(P(0, 10, 1, 1), DecIntToken, "23200123"),
			},
		},
		"decimal can't begin with underscores": {
			input: "_23_200_123",
			want: []*Token{
				V(P(0, 1, 1, 1), PrivateIdentifierToken, "_"),
				V(P(1, 10, 1, 2), DecIntToken, "23200123"),
			},
		},
		"decimal ends on last valid character": {
			input: "23foo",
			want: []*Token{
				V(P(0, 2, 1, 1), DecIntToken, "23"),
				V(P(2, 3, 1, 3), PublicIdentifierToken, "foo"),
			},
		},
		"hex": {
			input: "0x354ab1",
			want: []*Token{
				V(P(0, 8, 1, 1), HexIntToken, "354ab1"),
			},
		},
		"hex with underscores": {
			input: "0x35_4a_b1",
			want: []*Token{
				V(P(0, 10, 1, 1), HexIntToken, "354ab1"),
			},
		},
		"leading zeros invalidate hex": {
			input: "00x21",
			want: []*Token{
				V(P(0, 2, 1, 1), DecIntToken, "00"),
				V(P(2, 3, 1, 3), PublicIdentifierToken, "x21"),
			},
		},
		"hex without digits": {
			input: "0x",
			want: []*Token{
				V(P(0, 2, 1, 1), HexIntToken, ""),
			},
		},
		"hex with uppercase": {
			input: "0X354Ab1",
			want: []*Token{
				V(P(0, 8, 1, 1), HexIntToken, "354Ab1"),
			},
		},
		"hex ends on last valid character": {
			input: "0x354fp",
			want: []*Token{
				V(P(0, 6, 1, 1), HexIntToken, "354f"),
				V(P(6, 1, 1, 7), PublicIdentifierToken, "p"),
			},
		},
		"octal": {
			input: "0o603",
			want: []*Token{
				V(P(0, 5, 1, 1), OctIntToken, "603"),
			},
		},
		"octal with underscores": {
			input: "0o3201_5200",
			want: []*Token{
				V(P(0, 11, 1, 1), OctIntToken, "32015200"),
			},
		},
		"leading zeros invalidate octal": {
			input: "00o21",
			want: []*Token{
				V(P(0, 2, 1, 1), DecIntToken, "00"),
				V(P(2, 3, 1, 3), PublicIdentifierToken, "o21"),
			},
		},
		"octal without digits": {
			input: "0o",
			want: []*Token{
				V(P(0, 2, 1, 1), OctIntToken, ""),
			},
		},
		"octal with uppercase": {
			input: "0O603",
			want: []*Token{
				V(P(0, 5, 1, 1), OctIntToken, "603"),
			},
		},
		"octal ends on last valid character": {
			input: "0o6039abc1",
			want: []*Token{
				V(P(0, 5, 1, 1), OctIntToken, "603"),
				V(P(5, 1, 1, 6), DecIntToken, "9"),
				V(P(6, 4, 1, 7), PublicIdentifierToken, "abc1"),
			},
		},
		"quaternary": {
			input: "0q30212",
			want: []*Token{
				V(P(0, 7, 1, 1), QuatIntToken, "30212"),
			},
		},
		"quaternary with underscores": {
			input: "0q3201200_23010000",
			want: []*Token{
				V(P(0, 18, 1, 1), QuatIntToken, "320120023010000"),
			},
		},
		"leading zeros invalidate quaternary": {
			input: "00q21",
			want: []*Token{
				V(P(0, 2, 1, 1), DecIntToken, "00"),
				V(P(2, 3, 1, 3), PublicIdentifierToken, "q21"),
			},
		},
		"quaternary with uppercase": {
			input: "0Q30212",
			want: []*Token{
				V(P(0, 7, 1, 1), QuatIntToken, "30212"),
			},
		},
		"quaternary without digits": {
			input: "0q",
			want: []*Token{
				V(P(0, 2, 1, 1), QuatIntToken, ""),
			},
		},
		"quaternary ends on last valid character": {
			input: "0q302124a",
			want: []*Token{
				V(P(0, 7, 1, 1), QuatIntToken, "30212"),
				V(P(7, 1, 1, 8), DecIntToken, "4"),
				V(P(8, 1, 1, 9), PublicIdentifierToken, "a"),
			},
		},
		"binary": {
			input: "0b1010",
			want: []*Token{
				V(P(0, 6, 1, 1), BinIntToken, "1010"),
			},
		},
		"binary with underscores": {
			input: "0b10_10_10",
			want: []*Token{
				V(P(0, 10, 1, 1), BinIntToken, "101010"),
			},
		},
		"leading zeros invalidate binary": {
			input: "00b21",
			want: []*Token{
				V(P(0, 2, 1, 1), DecIntToken, "00"),
				V(P(2, 3, 1, 3), PublicIdentifierToken, "b21"),
			},
		},
		"binary with uppercase": {
			input: "0B1010",
			want: []*Token{
				V(P(0, 6, 1, 1), BinIntToken, "1010"),
			},
		},
		"binary ends on last valid character": {
			input: "0b10102dup",
			want: []*Token{
				V(P(0, 6, 1, 1), BinIntToken, "1010"),
				V(P(6, 1, 1, 7), DecIntToken, "2"),
				V(P(7, 3, 1, 8), PublicIdentifierToken, "dup"),
			},
		},
		"binary without digits": {
			input: "0b",
			want: []*Token{
				V(P(0, 2, 1, 1), BinIntToken, ""),
			},
		},
		"duodecimal": {
			input: "0da12b3",
			want: []*Token{
				V(P(0, 7, 1, 1), DuoIntToken, "a12b3"),
			},
		},
		"duodecimal with uppercase": {
			input: "0Da12B3",
			want: []*Token{
				V(P(0, 7, 1, 1), DuoIntToken, "a12B3"),
			},
		},
		"duodecimal ends on last valid character": {
			input: "0d23a3bca3",
			want: []*Token{
				V(P(0, 7, 1, 1), DuoIntToken, "23a3b"),
				V(P(7, 3, 1, 8), PublicIdentifierToken, "ca3"),
			},
		},
		"duodecimal without digits": {
			input: "0d",
			want: []*Token{
				V(P(0, 2, 1, 1), DuoIntToken, ""),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
