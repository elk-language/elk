package lexer

import "testing"

func TestInt(t *testing.T) {
	tests := testTable{
		"decimal": {
			input: "23",
			want: []*Token{
				{
					TokenType:  DecIntToken,
					Value:      "23",
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
			},
		},
		"decimal with leading zeros": {
			input: "00015",
			want: []*Token{
				{
					TokenType:  DecIntToken,
					Value:      "00015",
					StartByte:  0,
					ByteLength: 5,
					Line:       1,
					Column:     1,
				},
			},
		},
		"decimal with underscores": {
			input: "23_200_123",
			want: []*Token{
				{
					TokenType:  DecIntToken,
					Value:      "23200123",
					StartByte:  0,
					ByteLength: 10,
					Line:       1,
					Column:     1,
				},
			},
		},
		"decimal can't begin with underscores": {
			input: "_23_200_123",
			want: []*Token{
				{
					TokenType:  PrivateIdentifierToken,
					Value:      "_",
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  DecIntToken,
					Value:      "23200123",
					StartByte:  1,
					ByteLength: 10,
					Line:       1,
					Column:     2,
				},
			},
		},
		"decimal ends on last valid character": {
			input: "23foo",
			want: []*Token{
				{
					TokenType:  DecIntToken,
					Value:      "23",
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "foo",
					StartByte:  2,
					ByteLength: 3,
					Line:       1,
					Column:     3,
				},
			},
		},
		"hex": {
			input: "0x354ab1",
			want: []*Token{
				{
					TokenType:  HexIntToken,
					Value:      "354ab1",
					StartByte:  0,
					ByteLength: 8,
					Line:       1,
					Column:     1,
				},
			},
		},
		"hex with underscores": {
			input: "0x35_4a_b1",
			want: []*Token{
				{
					TokenType:  HexIntToken,
					Value:      "354ab1",
					StartByte:  0,
					ByteLength: 10,
					Line:       1,
					Column:     1,
				},
			},
		},
		"leading zeros invalidate hex": {
			input: "00x21",
			want: []*Token{
				{
					TokenType:  DecIntToken,
					Value:      "00",
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "x21",
					StartByte:  2,
					ByteLength: 3,
					Line:       1,
					Column:     3,
				},
			},
		},
		"hex without digits": {
			input: "0x",
			want: []*Token{
				{
					TokenType:  HexIntToken,
					Value:      "",
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
			},
		},
		"hex with uppercase": {
			input: "0X354Ab1",
			want: []*Token{
				{
					TokenType:  HexIntToken,
					Value:      "354Ab1",
					StartByte:  0,
					ByteLength: 8,
					Line:       1,
					Column:     1,
				},
			},
		},
		"hex ends on last valid character": {
			input: "0x354fp",
			want: []*Token{
				{
					TokenType:  HexIntToken,
					Value:      "354f",
					StartByte:  0,
					ByteLength: 6,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "p",
					StartByte:  6,
					ByteLength: 1,
					Line:       1,
					Column:     7,
				},
			},
		},
		"octal": {
			input: "0o603",
			want: []*Token{
				{
					TokenType:  OctIntToken,
					Value:      "603",
					StartByte:  0,
					ByteLength: 5,
					Line:       1,
					Column:     1,
				},
			},
		},
		"octal with underscores": {
			input: "0o3201_5200",
			want: []*Token{
				{
					TokenType:  OctIntToken,
					Value:      "32015200",
					StartByte:  0,
					ByteLength: 11,
					Line:       1,
					Column:     1,
				},
			},
		},
		"leading zeros invalidate octal": {
			input: "00o21",
			want: []*Token{
				{
					TokenType:  DecIntToken,
					Value:      "00",
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "o21",
					StartByte:  2,
					ByteLength: 3,
					Line:       1,
					Column:     3,
				},
			},
		},
		"octal without digits": {
			input: "0o",
			want: []*Token{
				{
					TokenType:  OctIntToken,
					Value:      "",
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
			},
		},
		"octal with uppercase": {
			input: "0O603",
			want: []*Token{
				{
					TokenType:  OctIntToken,
					Value:      "603",
					StartByte:  0,
					ByteLength: 5,
					Line:       1,
					Column:     1,
				},
			},
		},
		"octal ends on last valid character": {
			input: "0o6039abc1",
			want: []*Token{
				{
					TokenType:  OctIntToken,
					Value:      "603",
					StartByte:  0,
					ByteLength: 5,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  DecIntToken,
					Value:      "9",
					StartByte:  5,
					ByteLength: 1,
					Line:       1,
					Column:     6,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "abc1",
					StartByte:  6,
					ByteLength: 4,
					Line:       1,
					Column:     7,
				},
			},
		},
		"quaternary": {
			input: "0q30212",
			want: []*Token{
				{
					TokenType:  QuatIntToken,
					Value:      "30212",
					StartByte:  0,
					ByteLength: 7,
					Line:       1,
					Column:     1,
				},
			},
		},
		"quaternary with underscores": {
			input: "0q3201200_23010000",
			want: []*Token{
				{
					TokenType:  QuatIntToken,
					Value:      "320120023010000",
					StartByte:  0,
					ByteLength: 18,
					Line:       1,
					Column:     1,
				},
			},
		},
		"leading zeros invalidate quaternary": {
			input: "00q21",
			want: []*Token{
				{
					TokenType:  DecIntToken,
					Value:      "00",
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "q21",
					StartByte:  2,
					ByteLength: 3,
					Line:       1,
					Column:     3,
				},
			},
		},
		"quaternary with uppercase": {
			input: "0Q30212",
			want: []*Token{
				{
					TokenType:  QuatIntToken,
					Value:      "30212",
					StartByte:  0,
					ByteLength: 7,
					Line:       1,
					Column:     1,
				},
			},
		},
		"quaternary without digits": {
			input: "0q",
			want: []*Token{
				{
					TokenType:  QuatIntToken,
					Value:      "",
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
			},
		},
		"quaternary ends on last valid character": {
			input: "0q302124a",
			want: []*Token{
				{
					TokenType:  QuatIntToken,
					Value:      "30212",
					StartByte:  0,
					ByteLength: 7,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  DecIntToken,
					Value:      "4",
					StartByte:  7,
					ByteLength: 1,
					Line:       1,
					Column:     8,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "a",
					StartByte:  8,
					ByteLength: 1,
					Line:       1,
					Column:     9,
				},
			},
		},
		"binary": {
			input: "0b1010",
			want: []*Token{
				{
					TokenType:  BinIntToken,
					Value:      "1010",
					StartByte:  0,
					ByteLength: 6,
					Line:       1,
					Column:     1,
				},
			},
		},
		"binary with underscores": {
			input: "0b10_10_10",
			want: []*Token{
				{
					TokenType:  BinIntToken,
					Value:      "101010",
					StartByte:  0,
					ByteLength: 10,
					Line:       1,
					Column:     1,
				},
			},
		},
		"leading zeros invalidate binary": {
			input: "00b21",
			want: []*Token{
				{
					TokenType:  DecIntToken,
					Value:      "00",
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "b21",
					StartByte:  2,
					ByteLength: 3,
					Line:       1,
					Column:     3,
				},
			},
		},
		"binary with uppercase": {
			input: "0B1010",
			want: []*Token{
				{
					TokenType:  BinIntToken,
					Value:      "1010",
					StartByte:  0,
					ByteLength: 6,
					Line:       1,
					Column:     1,
				},
			},
		},
		"binary ends on last valid character": {
			input: "0b10102dup",
			want: []*Token{
				{
					TokenType:  BinIntToken,
					Value:      "1010",
					StartByte:  0,
					ByteLength: 6,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  DecIntToken,
					Value:      "2",
					StartByte:  6,
					ByteLength: 1,
					Line:       1,
					Column:     7,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "dup",
					StartByte:  7,
					ByteLength: 3,
					Line:       1,
					Column:     8,
				},
			},
		},
		"binary without digits": {
			input: "0b",
			want: []*Token{
				{
					TokenType:  BinIntToken,
					Value:      "",
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
			},
		},
		"duodecimal": {
			input: "0da12b3",
			want: []*Token{
				{
					TokenType:  DuoIntToken,
					Value:      "a12b3",
					StartByte:  0,
					ByteLength: 7,
					Line:       1,
					Column:     1,
				},
			},
		},
		"duodecimal with uppercase": {
			input: "0Da12B3",
			want: []*Token{
				{
					TokenType:  DuoIntToken,
					Value:      "a12B3",
					StartByte:  0,
					ByteLength: 7,
					Line:       1,
					Column:     1,
				},
			},
		},
		"duodecimal ends on last valid character": {
			input: "0d23a3bca3",
			want: []*Token{
				{
					TokenType:  DuoIntToken,
					Value:      "23a3b",
					StartByte:  0,
					ByteLength: 7,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "ca3",
					StartByte:  7,
					ByteLength: 3,
					Line:       1,
					Column:     8,
				},
			},
		},
		"duodecimal without digits": {
			input: "0d",
			want: []*Token{
				{
					TokenType:  DuoIntToken,
					Value:      "",
					StartByte:  0,
					ByteLength: 2,
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
