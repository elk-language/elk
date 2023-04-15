package lexer

import "testing"

func TestInt(t *testing.T) {
	tests := testTable{
		"decimal": {
			input: "23",
			want: []*Token{
				V(DecIntToken, "23", 0, 2, 1, 1),
			},
		},
		"decimal with leading zeros": {
			input: "00015",
			want: []*Token{
				V(DecIntToken, "00015", 0, 5, 1, 1),
			},
		},
		"decimal with underscores": {
			input: "23_200_123",
			want: []*Token{
				V(DecIntToken, "23200123", 0, 10, 1, 1),
			},
		},
		"decimal can't begin with underscores": {
			input: "_23_200_123",
			want: []*Token{
				V(PrivateIdentifierToken, "_", 0, 1, 1, 1),
				V(DecIntToken, "23200123", 1, 10, 1, 2),
			},
		},
		"decimal ends on last valid character": {
			input: "23foo",
			want: []*Token{
				V(DecIntToken, "23", 0, 2, 1, 1),
				V(IdentifierToken, "foo", 2, 3, 1, 3),
			},
		},
		"hex": {
			input: "0x354ab1",
			want: []*Token{
				V(HexIntToken, "354ab1", 0, 8, 1, 1),
			},
		},
		"hex with underscores": {
			input: "0x35_4a_b1",
			want: []*Token{
				V(HexIntToken, "354ab1", 0, 10, 1, 1),
			},
		},
		"leading zeros invalidate hex": {
			input: "00x21",
			want: []*Token{
				V(DecIntToken, "00", 0, 2, 1, 1),
				V(IdentifierToken, "x21", 2, 3, 1, 3),
			},
		},
		"hex without digits": {
			input: "0x",
			want: []*Token{
				V(HexIntToken, "", 0, 2, 1, 1),
			},
		},
		"hex with uppercase": {
			input: "0X354Ab1",
			want: []*Token{
				V(HexIntToken, "354Ab1", 0, 8, 1, 1),
			},
		},
		"hex ends on last valid character": {
			input: "0x354fp",
			want: []*Token{
				V(HexIntToken, "354f", 0, 6, 1, 1),
				V(IdentifierToken, "p", 6, 1, 1, 7),
			},
		},
		"octal": {
			input: "0o603",
			want: []*Token{
				V(OctIntToken, "603", 0, 5, 1, 1),
			},
		},
		"octal with underscores": {
			input: "0o3201_5200",
			want: []*Token{
				V(OctIntToken, "32015200", 0, 11, 1, 1),
			},
		},
		"leading zeros invalidate octal": {
			input: "00o21",
			want: []*Token{
				V(DecIntToken, "00", 0, 2, 1, 1),
				V(IdentifierToken, "o21", 2, 3, 1, 3),
			},
		},
		"octal without digits": {
			input: "0o",
			want: []*Token{
				V(OctIntToken, "", 0, 2, 1, 1),
			},
		},
		"octal with uppercase": {
			input: "0O603",
			want: []*Token{
				V(OctIntToken, "603", 0, 5, 1, 1),
			},
		},
		"octal ends on last valid character": {
			input: "0o6039abc1",
			want: []*Token{
				V(OctIntToken, "603", 0, 5, 1, 1),
				V(DecIntToken, "9", 5, 1, 1, 6),
				V(IdentifierToken, "abc1", 6, 4, 1, 7),
			},
		},
		"quaternary": {
			input: "0q30212",
			want: []*Token{
				V(QuatIntToken, "30212", 0, 7, 1, 1),
			},
		},
		"quaternary with underscores": {
			input: "0q3201200_23010000",
			want: []*Token{
				V(QuatIntToken, "320120023010000", 0, 18, 1, 1),
			},
		},
		"leading zeros invalidate quaternary": {
			input: "00q21",
			want: []*Token{
				V(DecIntToken, "00", 0, 2, 1, 1),
				V(IdentifierToken, "q21", 2, 3, 1, 3),
			},
		},
		"quaternary with uppercase": {
			input: "0Q30212",
			want: []*Token{
				V(QuatIntToken, "30212", 0, 7, 1, 1),
			},
		},
		"quaternary without digits": {
			input: "0q",
			want: []*Token{
				V(QuatIntToken, "", 0, 2, 1, 1),
			},
		},
		"quaternary ends on last valid character": {
			input: "0q302124a",
			want: []*Token{
				V(QuatIntToken, "30212", 0, 7, 1, 1),
				V(DecIntToken, "4", 7, 1, 1, 8),
				V(IdentifierToken, "a", 8, 1, 1, 9),
			},
		},
		"binary": {
			input: "0b1010",
			want: []*Token{
				V(BinIntToken, "1010", 0, 6, 1, 1),
			},
		},
		"binary with underscores": {
			input: "0b10_10_10",
			want: []*Token{
				V(BinIntToken, "101010", 0, 10, 1, 1),
			},
		},
		"leading zeros invalidate binary": {
			input: "00b21",
			want: []*Token{
				V(DecIntToken, "00", 0, 2, 1, 1),
				V(IdentifierToken, "b21", 2, 3, 1, 3),
			},
		},
		"binary with uppercase": {
			input: "0B1010",
			want: []*Token{
				V(BinIntToken, "1010", 0, 6, 1, 1),
			},
		},
		"binary ends on last valid character": {
			input: "0b10102dup",
			want: []*Token{
				V(BinIntToken, "1010", 0, 6, 1, 1),
				V(DecIntToken, "2", 6, 1, 1, 7),
				V(IdentifierToken, "dup", 7, 3, 1, 8),
			},
		},
		"binary without digits": {
			input: "0b",
			want: []*Token{
				V(BinIntToken, "", 0, 2, 1, 1),
			},
		},
		"duodecimal": {
			input: "0da12b3",
			want: []*Token{
				V(DuoIntToken, "a12b3", 0, 7, 1, 1),
			},
		},
		"duodecimal with uppercase": {
			input: "0Da12B3",
			want: []*Token{
				V(DuoIntToken, "a12B3", 0, 7, 1, 1),
			},
		},
		"duodecimal ends on last valid character": {
			input: "0d23a3bca3",
			want: []*Token{
				V(DuoIntToken, "23a3b", 0, 7, 1, 1),
				V(IdentifierToken, "ca3", 7, 3, 1, 8),
			},
		},
		"duodecimal without digits": {
			input: "0d",
			want: []*Token{
				V(DuoIntToken, "", 0, 2, 1, 1),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
