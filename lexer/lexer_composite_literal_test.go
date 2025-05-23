package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestArrayList(t *testing.T) {
	tests := testTable{
		"regular list": {
			input: "[1, 2, 3.0, 'foo', :bar]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(0, 1, 1))), token.LBRACKET),
				V(L(S(P(1, 1, 2), P(1, 1, 2))), token.INT, "1"),
				T(L(S(P(2, 1, 3), P(2, 1, 3))), token.COMMA),
				V(L(S(P(4, 1, 5), P(4, 1, 5))), token.INT, "2"),
				T(L(S(P(5, 1, 6), P(5, 1, 6))), token.COMMA),
				V(L(S(P(7, 1, 8), P(9, 1, 10))), token.FLOAT, "3.0"),
				T(L(S(P(10, 1, 11), P(10, 1, 11))), token.COMMA),
				V(L(S(P(12, 1, 13), P(16, 1, 17))), token.RAW_STRING, "foo"),
				T(L(S(P(17, 1, 18), P(17, 1, 18))), token.COMMA),
				T(L(S(P(19, 1, 20), P(19, 1, 20))), token.COLON),
				V(L(S(P(20, 1, 21), P(22, 1, 23))), token.PUBLIC_IDENTIFIER, "bar"),
				T(L(S(P(23, 1, 24), P(23, 1, 24))), token.RBRACKET),
			},
		},
		"word list": {
			input: `\w[
foo bar  
baz]`,
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.WORD_ARRAY_LIST_BEG),
				V(L(S(P(4, 2, 1), P(6, 2, 3))), token.RAW_STRING, "foo"),
				V(L(S(P(8, 2, 5), P(10, 2, 7))), token.RAW_STRING, "bar"),
				V(L(S(P(14, 3, 1), P(16, 3, 3))), token.RAW_STRING, "baz"),
				T(L(S(P(17, 3, 4), P(17, 3, 4))), token.WORD_ARRAY_LIST_END),
			},
		},
		"symbol list": {
			input: "\\s[foo bar   baz]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.SYMBOL_ARRAY_LIST_BEG),
				V(L(S(P(3, 1, 4), P(5, 1, 6))), token.RAW_STRING, "foo"),
				V(L(S(P(7, 1, 8), P(9, 1, 10))), token.RAW_STRING, "bar"),
				V(L(S(P(13, 1, 14), P(15, 1, 16))), token.RAW_STRING, "baz"),
				T(L(S(P(16, 1, 17), P(16, 1, 17))), token.SYMBOL_ARRAY_LIST_END),
			},
		},
		"hex list": {
			input: "\\x[ff 4_e   234]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.HEX_ARRAY_LIST_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0xff"),
				V(L(S(P(6, 1, 7), P(8, 1, 9))), token.INT, "0x4e"),
				V(L(S(P(12, 1, 13), P(14, 1, 15))), token.INT, "0x234"),
				T(L(S(P(15, 1, 16), P(15, 1, 16))), token.HEX_ARRAY_LIST_END),
			},
		},
		"hex list with invalid content": {
			input: "\\x[ff 4ghij   234]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.HEX_ARRAY_LIST_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0xff"),
				V(L(S(P(6, 1, 7), P(10, 1, 11))), token.ERROR, "invalid int literal"),
				V(L(S(P(14, 1, 15), P(16, 1, 17))), token.INT, "0x234"),
				T(L(S(P(17, 1, 18), P(17, 1, 18))), token.HEX_ARRAY_LIST_END),
			},
		},
		"hex list with invalid content at the end": {
			input: "\\x[ff 4ghij]\n",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.HEX_ARRAY_LIST_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0xff"),
				V(L(S(P(6, 1, 7), P(10, 1, 11))), token.ERROR, "invalid int literal"),
				T(L(S(P(11, 1, 12), P(11, 1, 12))), token.HEX_ARRAY_LIST_END),
				T(L(S(P(12, 1, 13), P(12, 1, 13))), token.NEWLINE),
			},
		},
		"binary list": {
			input: "\\b[11 1_0   101]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.BIN_ARRAY_LIST_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0b11"),
				V(L(S(P(6, 1, 7), P(8, 1, 9))), token.INT, "0b10"),
				V(L(S(P(12, 1, 13), P(14, 1, 15))), token.INT, "0b101"),
				T(L(S(P(15, 1, 16), P(15, 1, 16))), token.BIN_ARRAY_LIST_END),
			},
		},
		"binary list with invalid content": {
			input: "\\b[11 1ghij   101]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.BIN_ARRAY_LIST_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0b11"),
				V(L(S(P(6, 1, 7), P(10, 1, 11))), token.ERROR, "invalid int literal"),
				V(L(S(P(14, 1, 15), P(16, 1, 17))), token.INT, "0b101"),
				T(L(S(P(17, 1, 18), P(17, 1, 18))), token.BIN_ARRAY_LIST_END),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestHashSet(t *testing.T) {
	tests := testTable{
		"regular set": {
			input: "^[1, 2, 3.0, 'foo', :bar]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(1, 1, 2))), token.HASH_SET_LITERAL_BEG),
				V(L(S(P(2, 1, 3), P(2, 1, 3))), token.INT, "1"),
				T(L(S(P(3, 1, 4), P(3, 1, 4))), token.COMMA),
				V(L(S(P(5, 1, 6), P(5, 1, 6))), token.INT, "2"),
				T(L(S(P(6, 1, 7), P(6, 1, 7))), token.COMMA),
				V(L(S(P(8, 1, 9), P(10, 1, 11))), token.FLOAT, "3.0"),
				T(L(S(P(11, 1, 12), P(11, 1, 12))), token.COMMA),
				V(L(S(P(13, 1, 14), P(17, 1, 18))), token.RAW_STRING, "foo"),
				T(L(S(P(18, 1, 19), P(18, 1, 19))), token.COMMA),
				T(L(S(P(20, 1, 21), P(20, 1, 21))), token.COLON),
				V(L(S(P(21, 1, 22), P(23, 1, 24))), token.PUBLIC_IDENTIFIER, "bar"),
				T(L(S(P(24, 1, 25), P(24, 1, 25))), token.RBRACKET),
			},
		},
		"word set": {
			input: "^w[foo bar   baz]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.WORD_HASH_SET_BEG),
				V(L(S(P(3, 1, 4), P(5, 1, 6))), token.RAW_STRING, "foo"),
				V(L(S(P(7, 1, 8), P(9, 1, 10))), token.RAW_STRING, "bar"),
				V(L(S(P(13, 1, 14), P(15, 1, 16))), token.RAW_STRING, "baz"),
				T(L(S(P(16, 1, 17), P(16, 1, 17))), token.WORD_HASH_SET_END),
			},
		},
		"symbol set": {
			input: "^s[foo bar   baz]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.SYMBOL_HASH_SET_BEG),
				V(L(S(P(3, 1, 4), P(5, 1, 6))), token.RAW_STRING, "foo"),
				V(L(S(P(7, 1, 8), P(9, 1, 10))), token.RAW_STRING, "bar"),
				V(L(S(P(13, 1, 14), P(15, 1, 16))), token.RAW_STRING, "baz"),
				T(L(S(P(16, 1, 17), P(16, 1, 17))), token.SYMBOL_HASH_SET_END),
			},
		},
		"hex set": {
			input: "^x[ff 4_e   234]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.HEX_HASH_SET_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0xff"),
				V(L(S(P(6, 1, 7), P(8, 1, 9))), token.INT, "0x4e"),
				V(L(S(P(12, 1, 13), P(14, 1, 15))), token.INT, "0x234"),
				T(L(S(P(15, 1, 16), P(15, 1, 16))), token.HEX_HASH_SET_END),
			},
		},
		"hex set with invalid content": {
			input: "^x[ff 4ghij   234]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.HEX_HASH_SET_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0xff"),
				V(L(S(P(6, 1, 7), P(10, 1, 11))), token.ERROR, "invalid int literal"),
				V(L(S(P(14, 1, 15), P(16, 1, 17))), token.INT, "0x234"),
				T(L(S(P(17, 1, 18), P(17, 1, 18))), token.HEX_HASH_SET_END),
			},
		},
		"binary set": {
			input: "^b[11 1_0   101]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.BIN_HASH_SET_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0b11"),
				V(L(S(P(6, 1, 7), P(8, 1, 9))), token.INT, "0b10"),
				V(L(S(P(12, 1, 13), P(14, 1, 15))), token.INT, "0b101"),
				T(L(S(P(15, 1, 16), P(15, 1, 16))), token.BIN_HASH_SET_END),
			},
		},
		"binary set with invalid content": {
			input: "^b[11 1ghij   101]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.BIN_HASH_SET_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0b11"),
				V(L(S(P(6, 1, 7), P(10, 1, 11))), token.ERROR, "invalid int literal"),
				V(L(S(P(14, 1, 15), P(16, 1, 17))), token.INT, "0b101"),
				T(L(S(P(17, 1, 18), P(17, 1, 18))), token.BIN_HASH_SET_END),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestArrayTuple(t *testing.T) {
	tests := testTable{
		"regular arrayTuple": {
			input: "%[1, 2, 3.0, 'foo', :bar]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(1, 1, 2))), token.TUPLE_LITERAL_BEG),
				V(L(S(P(2, 1, 3), P(2, 1, 3))), token.INT, "1"),
				T(L(S(P(3, 1, 4), P(3, 1, 4))), token.COMMA),
				V(L(S(P(5, 1, 6), P(5, 1, 6))), token.INT, "2"),
				T(L(S(P(6, 1, 7), P(6, 1, 7))), token.COMMA),
				V(L(S(P(8, 1, 9), P(10, 1, 11))), token.FLOAT, "3.0"),
				T(L(S(P(11, 1, 12), P(11, 1, 12))), token.COMMA),
				V(L(S(P(13, 1, 14), P(17, 1, 18))), token.RAW_STRING, "foo"),
				T(L(S(P(18, 1, 19), P(18, 1, 19))), token.COMMA),
				T(L(S(P(20, 1, 21), P(20, 1, 21))), token.COLON),
				V(L(S(P(21, 1, 22), P(23, 1, 24))), token.PUBLIC_IDENTIFIER, "bar"),
				T(L(S(P(24, 1, 25), P(24, 1, 25))), token.RBRACKET),
			},
		},
		"word arrayTuple": {
			input: "%w[foo bar   baz]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.WORD_ARRAY_TUPLE_BEG),
				V(L(S(P(3, 1, 4), P(5, 1, 6))), token.RAW_STRING, "foo"),
				V(L(S(P(7, 1, 8), P(9, 1, 10))), token.RAW_STRING, "bar"),
				V(L(S(P(13, 1, 14), P(15, 1, 16))), token.RAW_STRING, "baz"),
				T(L(S(P(16, 1, 17), P(16, 1, 17))), token.WORD_ARRAY_TUPLE_END),
			},
		},
		"symbol arrayTuple": {
			input: "%s[foo bar   baz]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.SYMBOL_ARRAY_TUPLE_BEG),
				V(L(S(P(3, 1, 4), P(5, 1, 6))), token.RAW_STRING, "foo"),
				V(L(S(P(7, 1, 8), P(9, 1, 10))), token.RAW_STRING, "bar"),
				V(L(S(P(13, 1, 14), P(15, 1, 16))), token.RAW_STRING, "baz"),
				T(L(S(P(16, 1, 17), P(16, 1, 17))), token.SYMBOL_ARRAY_TUPLE_END),
			},
		},
		"hex arrayTuple": {
			input: "%x[ff 4_e   234]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.HEX_ARRAY_TUPLE_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0xff"),
				V(L(S(P(6, 1, 7), P(8, 1, 9))), token.INT, "0x4e"),
				V(L(S(P(12, 1, 13), P(14, 1, 15))), token.INT, "0x234"),
				T(L(S(P(15, 1, 16), P(15, 1, 16))), token.HEX_ARRAY_TUPLE_END),
			},
		},
		"hex arrayTuple with invalid content": {
			input: "%x[ff 4ghij   234]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.HEX_ARRAY_TUPLE_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0xff"),
				V(L(S(P(6, 1, 7), P(10, 1, 11))), token.ERROR, "invalid int literal"),
				V(L(S(P(14, 1, 15), P(16, 1, 17))), token.INT, "0x234"),
				T(L(S(P(17, 1, 18), P(17, 1, 18))), token.HEX_ARRAY_TUPLE_END),
			},
		},
		"binary arrayTuple": {
			input: "%b[11 1_0   101]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.BIN_ARRAY_TUPLE_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0b11"),
				V(L(S(P(6, 1, 7), P(8, 1, 9))), token.INT, "0b10"),
				V(L(S(P(12, 1, 13), P(14, 1, 15))), token.INT, "0b101"),
				T(L(S(P(15, 1, 16), P(15, 1, 16))), token.BIN_ARRAY_TUPLE_END),
			},
		},
		"binary arrayTuple with invalid content": {
			input: "%b[11 1ghij   101]",
			want: []*token.Token{
				T(L(S(P(0, 1, 1), P(2, 1, 3))), token.BIN_ARRAY_TUPLE_BEG),
				V(L(S(P(3, 1, 4), P(4, 1, 5))), token.INT, "0b11"),
				V(L(S(P(6, 1, 7), P(10, 1, 11))), token.ERROR, "invalid int literal"),
				V(L(S(P(14, 1, 15), P(16, 1, 17))), token.INT, "0b101"),
				T(L(S(P(17, 1, 18), P(17, 1, 18))), token.BIN_ARRAY_TUPLE_END),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
