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
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.INT, "23"),
			},
		},
		"decimal with leading zeros": {
			input: "00015",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT, "00015"),
			},
		},
		"decimal with underscores": {
			input: "23_200_123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(9, 1, 10))), token.INT, "23200123"),
			},
		},
		"decimal cannot begin with underscores": {
			input: "_23_200_123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(0, 1, 1))), token.PRIVATE_IDENTIFIER, "_"),
				V(L(S(P(1, 1, 2), P(10, 1, 11))), token.INT, "23200123"),
			},
		},
		"decimal ends on last valid character": {
			input: "23kar",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.INT, "23"),
				V(L(S(P(2, 1, 3), P(4, 1, 5))), token.PUBLIC_IDENTIFIER, "kar"),
			},
		},
		"decimal int64": {
			input: "23i64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT64, "23"),
			},
		},
		"decimal uint64": {
			input: "23u64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.UINT64, "23"),
			},
		},
		"decimal int32": {
			input: "23i32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT32, "23"),
			},
		},
		"decimal uint32": {
			input: "23u32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.UINT32, "23"),
			},
		},
		"decimal int16": {
			input: "23i16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT16, "23"),
			},
		},
		"decimal uint16": {
			input: "23u16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.UINT16, "23"),
			},
		},
		"decimal int8": {
			input: "23i8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.INT8, "23"),
			},
		},
		"decimal uint8": {
			input: "23u8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.UINT8, "23"),
			},
		},
		"hex": {
			input: "0x354ab1",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(7, 1, 8))), token.INT, "0x354ab1"),
			},
		},
		"hex with underscores": {
			input: "0x35_4a_b1",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(9, 1, 10))), token.INT, "0x354ab1"),
			},
		},
		"leading zeros invalidate hex": {
			input: "00x21",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.INT, "00"),
				V(L(S(P(2, 1, 3), P(4, 1, 5))), token.PUBLIC_IDENTIFIER, "x21"),
			},
		},
		"hex without digits": {
			input: "0x",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.INT, "0x"),
			},
		},
		"hex with uppercase": {
			input: "0X354Ab1",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(7, 1, 8))), token.INT, "0x354Ab1"),
			},
		},
		"hex ends on last valid character": {
			input: "0x354fp",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT, "0x354f"),
				V(L(S(P(6, 1, 7), P(6, 1, 7))), token.PUBLIC_IDENTIFIER, "p"),
			},
		},
		"hex int64": {
			input: "0xfi64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT64, "0xf"),
			},
		},
		"hex uint64": {
			input: "0xfu64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT64, "0xf"),
			},
		},
		"hex int32": {
			input: "0xfi32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT32, "0xf"),
			},
		},
		"hex uint32": {
			input: "0xfu32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT32, "0xf"),
			},
		},
		"hex int16": {
			input: "0xfi16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT16, "0xf"),
			},
		},
		"hex uint16": {
			input: "0xfu16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT16, "0xf"),
			},
		},
		"hex int8": {
			input: "0xfi8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT8, "0xf"),
			},
		},
		"hex uint8": {
			input: "0xfu8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.UINT8, "0xf"),
			},
		},
		"octal": {
			input: "0o603",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT, "0o603"),
			},
		},
		"octal with underscores": {
			input: "0o3201_5200",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(10, 1, 11))), token.INT, "0o32015200"),
			},
		},
		"leading zeros invalidate octal": {
			input: "00o21",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.INT, "00"),
				V(L(S(P(2, 1, 3), P(4, 1, 5))), token.PUBLIC_IDENTIFIER, "o21"),
			},
		},
		"octal without digits": {
			input: "0o",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.INT, "0o"),
			},
		},
		"octal with uppercase": {
			input: "0O603",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT, "0o603"),
			},
		},
		"octal ends on last valid character": {
			input: "0o6039abc1",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT, "0o603"),
				V(L(S(P(5, 1, 6), P(5, 1, 6))), token.INT, "9"),
				V(L(S(P(6, 1, 7), P(9, 1, 10))), token.PUBLIC_IDENTIFIER, "abc1"),
			},
		},
		"octal int64": {
			input: "0o5i64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT64, "0o5"),
			},
		},
		"octal uint64": {
			input: "0o5u64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT64, "0o5"),
			},
		},
		"octal int32": {
			input: "0o5i32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT32, "0o5"),
			},
		},
		"octal uint32": {
			input: "0o5u32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT32, "0o5"),
			},
		},
		"octal int16": {
			input: "0o5i16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT16, "0o5"),
			},
		},
		"octal uint16": {
			input: "0o5u16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT16, "0o5"),
			},
		},
		"octal int8": {
			input: "0o5i8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT8, "0o5"),
			},
		},
		"octal uint8": {
			input: "0o5u8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.UINT8, "0o5"),
			},
		},
		"quaternary": {
			input: "0q30212",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 7))), token.INT, "0q30212"),
			},
		},
		"quaternary with underscores": {
			input: "0q3201200_23010000",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(17, 1, 18))), token.INT, "0q320120023010000"),
			},
		},
		"leading zeros invalidate quaternary": {
			input: "00q21",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.INT, "00"),
				V(L(S(P(2, 1, 3), P(4, 1, 5))), token.PUBLIC_IDENTIFIER, "q21"),
			},
		},
		"quaternary with uppercase": {
			input: "0Q30212",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 7))), token.INT, "0q30212"),
			},
		},
		"quaternary without digits": {
			input: "0q",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.INT, "0q"),
			},
		},
		"quaternary ends on last valid character": {
			input: "0q302124a",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 7))), token.INT, "0q30212"),
				V(L(S(P(7, 1, 8), P(7, 1, 8))), token.INT, "4"),
				V(L(S(P(8, 1, 9), P(8, 1, 9))), token.PUBLIC_IDENTIFIER, "a"),
			},
		},
		"quaternary int64": {
			input: "0q2i64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT64, "0q2"),
			},
		},
		"quaternary uint64": {
			input: "0q2u64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT64, "0q2"),
			},
		},
		"quaternary int32": {
			input: "0q2i32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT32, "0q2"),
			},
		},
		"quaternary uint32": {
			input: "0q2u32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT32, "0q2"),
			},
		},
		"quaternary int16": {
			input: "0q2i16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT16, "0q2"),
			},
		},
		"quaternary uint16": {
			input: "0q2u16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT16, "0q2"),
			},
		},
		"quaternary int8": {
			input: "0q2i8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT8, "0q2"),
			},
		},
		"quaternary uint8": {
			input: "0q2u8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.UINT8, "0q2"),
			},
		},
		"binary": {
			input: "0b1010",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT, "0b1010"),
			},
		},
		"binary with underscores": {
			input: "0b10_10_10",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(9, 1, 10))), token.INT, "0b101010"),
			},
		},
		"leading zeros invalidate binary": {
			input: "00b21",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(2, 1, 3))), token.ERROR, "invalid big numeric literal"),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "21"),
			},
		},
		"binary with uppercase": {
			input: "0B1010",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT, "0b1010"),
			},
		},
		"binary ends on last valid character": {
			input: "0b10102dup",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT, "0b1010"),
				V(L(S(P(6, 1, 7), P(6, 1, 7))), token.INT, "2"),
				V(L(S(P(7, 1, 8), P(9, 1, 10))), token.PUBLIC_IDENTIFIER, "dup"),
			},
		},
		"binary without digits": {
			input: "0b",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.INT, "0b"),
			},
		},
		"binary int64": {
			input: "0b1i64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT64, "0b1"),
			},
		},
		"binary uint64": {
			input: "0b1u64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT64, "0b1"),
			},
		},
		"binary int32": {
			input: "0b1i32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT32, "0b1"),
			},
		},
		"binary uint32": {
			input: "0b1u32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT32, "0b1"),
			},
		},
		"binary int16": {
			input: "0b1i16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT16, "0b1"),
			},
		},
		"binary uint16": {
			input: "0b1u16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT16, "0b1"),
			},
		},
		"binary int8": {
			input: "0b1i8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT8, "0b1"),
			},
		},
		"binary uint8": {
			input: "0b1u8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.UINT8, "0b1"),
			},
		},
		"duodecimal": {
			input: "0da12b3",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 7))), token.INT, "0da12b3"),
			},
		},
		"duodecimal with uppercase": {
			input: "0Da12B3",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 7))), token.INT, "0da12B3"),
			},
		},
		"duodecimal ends on last valid character": {
			input: "0d23a3bca3",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 7))), token.INT, "0d23a3b"),
				V(L(S(P(7, 1, 8), P(9, 1, 10))), token.PUBLIC_IDENTIFIER, "ca3"),
			},
		},
		"duodecimal without digits": {
			input: "0d",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.INT, "0d"),
			},
		},
		"duodecimal int64": {
			input: "0d1i64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT64, "0d1"),
			},
		},
		"duodecimal uint64": {
			input: "0d1u64",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT64, "0d1"),
			},
		},
		"duodecimal int32": {
			input: "0d1i32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT32, "0d1"),
			},
		},
		"duodecimal uint32": {
			input: "0d1u32",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT32, "0d1"),
			},
		},
		"duodecimal int16": {
			input: "0d1i16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INT16, "0d1"),
			},
		},
		"duodecimal uint16": {
			input: "0d1u16",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.UINT16, "0d1"),
			},
		},
		"duodecimal int8": {
			input: "0d1i8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INT8, "0d1"),
			},
		},
		"duodecimal uint8": {
			input: "0d1u8",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.UINT8, "0d1"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
