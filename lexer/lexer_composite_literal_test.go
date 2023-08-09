package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestCollectionLiteral(t *testing.T) {
	tests := testTable{
		"incorrect word delimiters": {
			input: "%w<foo bar   baz>",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(1, 1, 2)), token.ERROR, "invalid word collection literal delimiters `%w`"),
				T(S(P(2, 1, 3), P(2, 1, 3)), token.LESS),
				V(S(P(3, 1, 4), P(5, 1, 6)), token.PUBLIC_IDENTIFIER, "foo"),
				V(S(P(7, 1, 8), P(9, 1, 10)), token.PUBLIC_IDENTIFIER, "bar"),
				V(S(P(13, 1, 14), P(15, 1, 16)), token.PUBLIC_IDENTIFIER, "baz"),
				T(S(P(16, 1, 17), P(16, 1, 17)), token.GREATER),
			},
		},
		"incorrect symbol delimiters": {
			input: "%s<foo bar   baz>",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(1, 1, 2)), token.ERROR, "invalid symbol collection literal delimiters `%s`"),
				T(S(P(2, 1, 3), P(2, 1, 3)), token.LESS),
				V(S(P(3, 1, 4), P(5, 1, 6)), token.PUBLIC_IDENTIFIER, "foo"),
				V(S(P(7, 1, 8), P(9, 1, 10)), token.PUBLIC_IDENTIFIER, "bar"),
				V(S(P(13, 1, 14), P(15, 1, 16)), token.PUBLIC_IDENTIFIER, "baz"),
				T(S(P(16, 1, 17), P(16, 1, 17)), token.GREATER),
			},
		},
		"incorrect hex delimiters": {
			input: "%x<45a 101   fff>",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(1, 1, 2)), token.ERROR, "invalid hex collection literal delimiters `%x`"),
				T(S(P(2, 1, 3), P(2, 1, 3)), token.LESS),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "45"),
				V(S(P(5, 1, 6), P(5, 1, 6)), token.PUBLIC_IDENTIFIER, "a"),
				V(S(P(7, 1, 8), P(9, 1, 10)), token.INT, "101"),
				V(S(P(13, 1, 14), P(15, 1, 16)), token.PUBLIC_IDENTIFIER, "fff"),
				T(S(P(16, 1, 17), P(16, 1, 17)), token.GREATER),
			},
		},
		"incorrect binary delimiters": {
			input: "%b<110 101   111>",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(1, 1, 2)), token.ERROR, "invalid binary collection literal delimiters `%b`"),
				T(S(P(2, 1, 3), P(2, 1, 3)), token.LESS),
				V(S(P(3, 1, 4), P(5, 1, 6)), token.INT, "110"),
				V(S(P(7, 1, 8), P(9, 1, 10)), token.INT, "101"),
				V(S(P(13, 1, 14), P(15, 1, 16)), token.INT, "111"),
				T(S(P(16, 1, 17), P(16, 1, 17)), token.GREATER),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestArray(t *testing.T) {
	tests := testTable{
		"regular array": {
			input: "[1, 2, 3.0, 'foo', :bar]",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(0, 1, 1)), token.LBRACKET),
				V(S(P(1, 1, 2), P(1, 1, 2)), token.INT, "1"),
				T(S(P(2, 1, 3), P(2, 1, 3)), token.COMMA),
				V(S(P(4, 1, 5), P(4, 1, 5)), token.INT, "2"),
				T(S(P(5, 1, 6), P(5, 1, 6)), token.COMMA),
				V(S(P(7, 1, 8), P(9, 1, 10)), token.FLOAT, "3.0"),
				T(S(P(10, 1, 11), P(10, 1, 11)), token.COMMA),
				V(S(P(12, 1, 13), P(16, 1, 17)), token.RAW_STRING, "foo"),
				T(S(P(17, 1, 18), P(17, 1, 18)), token.COMMA),
				T(S(P(19, 1, 20), P(19, 1, 20)), token.COLON),
				V(S(P(20, 1, 21), P(22, 1, 23)), token.PUBLIC_IDENTIFIER, "bar"),
				T(S(P(23, 1, 24), P(23, 1, 24)), token.RBRACKET),
			},
		},
		"word array": {
			input: `%w[
foo bar  
baz]`,
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.WORD_LIST_BEG),
				V(S(P(4, 2, 1), P(6, 2, 3)), token.RAW_STRING, "foo"),
				V(S(P(8, 2, 5), P(10, 2, 7)), token.RAW_STRING, "bar"),
				V(S(P(14, 3, 1), P(16, 3, 3)), token.RAW_STRING, "baz"),
				T(S(P(17, 3, 4), P(17, 3, 4)), token.WORD_LIST_END),
			},
		},
		"symbol array": {
			input: "%s[foo bar   baz]",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.SYMBOL_LIST_BEG),
				V(S(P(3, 1, 4), P(5, 1, 6)), token.RAW_STRING, "foo"),
				V(S(P(7, 1, 8), P(9, 1, 10)), token.RAW_STRING, "bar"),
				V(S(P(13, 1, 14), P(15, 1, 16)), token.RAW_STRING, "baz"),
				T(S(P(16, 1, 17), P(16, 1, 17)), token.SYMBOL_LIST_END),
			},
		},
		"hex array": {
			input: "%x[ff 4_e   234]",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.HEX_LIST_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0xff"),
				V(S(P(6, 1, 7), P(8, 1, 9)), token.INT, "0x4e"),
				V(S(P(12, 1, 13), P(14, 1, 15)), token.INT, "0x234"),
				T(S(P(15, 1, 16), P(15, 1, 16)), token.HEX_LIST_END),
			},
		},
		"hex array with invalid content": {
			input: "%x[ff 4ghij   234]",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.HEX_LIST_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0xff"),
				V(S(P(6, 1, 7), P(10, 1, 11)), token.ERROR, "invalid int literal"),
				V(S(P(14, 1, 15), P(16, 1, 17)), token.INT, "0x234"),
				T(S(P(17, 1, 18), P(17, 1, 18)), token.HEX_LIST_END),
			},
		},
		"hex array with invalid content at the end": {
			input: "%x[ff 4ghij]\n",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.HEX_LIST_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0xff"),
				V(S(P(6, 1, 7), P(10, 1, 11)), token.ERROR, "invalid int literal"),
				T(S(P(11, 1, 12), P(11, 1, 12)), token.HEX_LIST_END),
				T(S(P(12, 1, 13), P(12, 1, 13)), token.NEWLINE),
			},
		},
		"binary array": {
			input: "%b[11 1_0   101]",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.BIN_LIST_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0b11"),
				V(S(P(6, 1, 7), P(8, 1, 9)), token.INT, "0b10"),
				V(S(P(12, 1, 13), P(14, 1, 15)), token.INT, "0b101"),
				T(S(P(15, 1, 16), P(15, 1, 16)), token.BIN_LIST_END),
			},
		},
		"binary array with invalid content": {
			input: "%b[11 1ghij   101]",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.BIN_LIST_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0b11"),
				V(S(P(6, 1, 7), P(10, 1, 11)), token.ERROR, "invalid int literal"),
				V(S(P(14, 1, 15), P(16, 1, 17)), token.INT, "0b101"),
				T(S(P(17, 1, 18), P(17, 1, 18)), token.BIN_LIST_END),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestSet(t *testing.T) {
	tests := testTable{
		"regular set": {
			input: "%{1, 2, 3.0, 'foo', :bar}",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.SET_LITERAL_BEG),
				V(S(P(2, 1, 3), P(2, 1, 3)), token.INT, "1"),
				T(S(P(3, 1, 4), P(3, 1, 4)), token.COMMA),
				V(S(P(5, 1, 6), P(5, 1, 6)), token.INT, "2"),
				T(S(P(6, 1, 7), P(6, 1, 7)), token.COMMA),
				V(S(P(8, 1, 9), P(10, 1, 11)), token.FLOAT, "3.0"),
				T(S(P(11, 1, 12), P(11, 1, 12)), token.COMMA),
				V(S(P(13, 1, 14), P(17, 1, 18)), token.RAW_STRING, "foo"),
				T(S(P(18, 1, 19), P(18, 1, 19)), token.COMMA),
				T(S(P(20, 1, 21), P(20, 1, 21)), token.COLON),
				V(S(P(21, 1, 22), P(23, 1, 24)), token.PUBLIC_IDENTIFIER, "bar"),
				T(S(P(24, 1, 25), P(24, 1, 25)), token.RBRACE),
			},
		},
		"word set": {
			input: "%w{foo bar   baz}",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.WORD_SET_BEG),
				V(S(P(3, 1, 4), P(5, 1, 6)), token.RAW_STRING, "foo"),
				V(S(P(7, 1, 8), P(9, 1, 10)), token.RAW_STRING, "bar"),
				V(S(P(13, 1, 14), P(15, 1, 16)), token.RAW_STRING, "baz"),
				T(S(P(16, 1, 17), P(16, 1, 17)), token.WORD_SET_END),
			},
		},
		"symbol set": {
			input: "%s{foo bar   baz}",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.SYMBOL_SET_BEG),
				V(S(P(3, 1, 4), P(5, 1, 6)), token.RAW_STRING, "foo"),
				V(S(P(7, 1, 8), P(9, 1, 10)), token.RAW_STRING, "bar"),
				V(S(P(13, 1, 14), P(15, 1, 16)), token.RAW_STRING, "baz"),
				T(S(P(16, 1, 17), P(16, 1, 17)), token.SYMBOL_SET_END),
			},
		},
		"hex set": {
			input: "%x{ff 4_e   234}",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.HEX_SET_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0xff"),
				V(S(P(6, 1, 7), P(8, 1, 9)), token.INT, "0x4e"),
				V(S(P(12, 1, 13), P(14, 1, 15)), token.INT, "0x234"),
				T(S(P(15, 1, 16), P(15, 1, 16)), token.HEX_SET_END),
			},
		},
		"hex set with invalid content": {
			input: "%x{ff 4ghij   234}",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.HEX_SET_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0xff"),
				V(S(P(6, 1, 7), P(10, 1, 11)), token.ERROR, "invalid int literal"),
				V(S(P(14, 1, 15), P(16, 1, 17)), token.INT, "0x234"),
				T(S(P(17, 1, 18), P(17, 1, 18)), token.HEX_SET_END),
			},
		},
		"binary array": {
			input: "%b[11 1_0   101]",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.BIN_LIST_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0b11"),
				V(S(P(6, 1, 7), P(8, 1, 9)), token.INT, "0b10"),
				V(S(P(12, 1, 13), P(14, 1, 15)), token.INT, "0b101"),
				T(S(P(15, 1, 16), P(15, 1, 16)), token.BIN_LIST_END),
			},
		},
		"binary set with invalid content": {
			input: "%b{11 1ghij   101}",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.BIN_SET_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0b11"),
				V(S(P(6, 1, 7), P(10, 1, 11)), token.ERROR, "invalid int literal"),
				V(S(P(14, 1, 15), P(16, 1, 17)), token.INT, "0b101"),
				T(S(P(17, 1, 18), P(17, 1, 18)), token.BIN_SET_END),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestTuple(t *testing.T) {
	tests := testTable{
		"regular tuple": {
			input: "%(1, 2, 3.0, 'foo', :bar)",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(1, 1, 2)), token.TUPLE_LITERAL_BEG),
				V(S(P(2, 1, 3), P(2, 1, 3)), token.INT, "1"),
				T(S(P(3, 1, 4), P(3, 1, 4)), token.COMMA),
				V(S(P(5, 1, 6), P(5, 1, 6)), token.INT, "2"),
				T(S(P(6, 1, 7), P(6, 1, 7)), token.COMMA),
				V(S(P(8, 1, 9), P(10, 1, 11)), token.FLOAT, "3.0"),
				T(S(P(11, 1, 12), P(11, 1, 12)), token.COMMA),
				V(S(P(13, 1, 14), P(17, 1, 18)), token.RAW_STRING, "foo"),
				T(S(P(18, 1, 19), P(18, 1, 19)), token.COMMA),
				T(S(P(20, 1, 21), P(20, 1, 21)), token.COLON),
				V(S(P(21, 1, 22), P(23, 1, 24)), token.PUBLIC_IDENTIFIER, "bar"),
				T(S(P(24, 1, 25), P(24, 1, 25)), token.RPAREN),
			},
		},
		"word tuple": {
			input: "%w(foo bar   baz)",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.WORD_TUPLE_BEG),
				V(S(P(3, 1, 4), P(5, 1, 6)), token.RAW_STRING, "foo"),
				V(S(P(7, 1, 8), P(9, 1, 10)), token.RAW_STRING, "bar"),
				V(S(P(13, 1, 14), P(15, 1, 16)), token.RAW_STRING, "baz"),
				T(S(P(16, 1, 17), P(16, 1, 17)), token.WORD_TUPLE_END),
			},
		},
		"symbol tuple": {
			input: "%s(foo bar   baz)",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.SYMBOL_TUPLE_BEG),
				V(S(P(3, 1, 4), P(5, 1, 6)), token.RAW_STRING, "foo"),
				V(S(P(7, 1, 8), P(9, 1, 10)), token.RAW_STRING, "bar"),
				V(S(P(13, 1, 14), P(15, 1, 16)), token.RAW_STRING, "baz"),
				T(S(P(16, 1, 17), P(16, 1, 17)), token.SYMBOL_TUPLE_END),
			},
		},
		"hex tuple": {
			input: "%x(ff 4_e   234)",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.HEX_TUPLE_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0xff"),
				V(S(P(6, 1, 7), P(8, 1, 9)), token.INT, "0x4e"),
				V(S(P(12, 1, 13), P(14, 1, 15)), token.INT, "0x234"),
				T(S(P(15, 1, 16), P(15, 1, 16)), token.HEX_TUPLE_END),
			},
		},
		"hex tuple with invalid content": {
			input: "%x(ff 4ghij   234)",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.HEX_TUPLE_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0xff"),
				V(S(P(6, 1, 7), P(10, 1, 11)), token.ERROR, "invalid int literal"),
				V(S(P(14, 1, 15), P(16, 1, 17)), token.INT, "0x234"),
				T(S(P(17, 1, 18), P(17, 1, 18)), token.HEX_TUPLE_END),
			},
		},
		"binary tuple": {
			input: "%b(11 1_0   101)",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.BIN_TUPLE_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0b11"),
				V(S(P(6, 1, 7), P(8, 1, 9)), token.INT, "0b10"),
				V(S(P(12, 1, 13), P(14, 1, 15)), token.INT, "0b101"),
				T(S(P(15, 1, 16), P(15, 1, 16)), token.BIN_TUPLE_END),
			},
		},
		"binary tuple with invalid content": {
			input: "%b(11 1ghij   101)",
			want: []*token.Token{
				T(S(P(0, 1, 1), P(2, 1, 3)), token.BIN_TUPLE_BEG),
				V(S(P(3, 1, 4), P(4, 1, 5)), token.INT, "0b11"),
				V(S(P(6, 1, 7), P(10, 1, 11)), token.ERROR, "invalid int literal"),
				V(S(P(14, 1, 15), P(16, 1, 17)), token.INT, "0b101"),
				T(S(P(17, 1, 18), P(17, 1, 18)), token.BIN_TUPLE_END),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
