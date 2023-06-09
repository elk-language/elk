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
				V(P(0, 2, 1, 1), token.ERROR, "invalid word collection literal delimiters `%w`"),
				T(P(2, 1, 1, 3), token.LESS),
				V(P(3, 3, 1, 4), token.PUBLIC_IDENTIFIER, "foo"),
				V(P(7, 3, 1, 8), token.PUBLIC_IDENTIFIER, "bar"),
				V(P(13, 3, 1, 14), token.PUBLIC_IDENTIFIER, "baz"),
				T(P(16, 1, 1, 17), token.GREATER),
			},
		},
		"incorrect symbol delimiters": {
			input: "%s<foo bar   baz>",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.ERROR, "invalid symbol collection literal delimiters `%s`"),
				T(P(2, 1, 1, 3), token.LESS),
				V(P(3, 3, 1, 4), token.PUBLIC_IDENTIFIER, "foo"),
				V(P(7, 3, 1, 8), token.PUBLIC_IDENTIFIER, "bar"),
				V(P(13, 3, 1, 14), token.PUBLIC_IDENTIFIER, "baz"),
				T(P(16, 1, 1, 17), token.GREATER),
			},
		},
		"incorrect hex delimiters": {
			input: "%x<45a 101   fff>",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.ERROR, "invalid hex collection literal delimiters `%x`"),
				T(P(2, 1, 1, 3), token.LESS),
				V(P(3, 2, 1, 4), token.DEC_INT, "45"),
				V(P(5, 1, 1, 6), token.PUBLIC_IDENTIFIER, "a"),
				V(P(7, 3, 1, 8), token.DEC_INT, "101"),
				V(P(13, 3, 1, 14), token.PUBLIC_IDENTIFIER, "fff"),
				T(P(16, 1, 1, 17), token.GREATER),
			},
		},
		"incorrect binary delimiters": {
			input: "%b<110 101   111>",
			want: []*token.Token{
				V(P(0, 2, 1, 1), token.ERROR, "invalid binary collection literal delimiters `%b`"),
				T(P(2, 1, 1, 3), token.LESS),
				V(P(3, 3, 1, 4), token.DEC_INT, "110"),
				V(P(7, 3, 1, 8), token.DEC_INT, "101"),
				V(P(13, 3, 1, 14), token.DEC_INT, "111"),
				T(P(16, 1, 1, 17), token.GREATER),
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
				T(P(0, 1, 1, 1), token.LBRACKET),
				V(P(1, 1, 1, 2), token.DEC_INT, "1"),
				T(P(2, 1, 1, 3), token.COMMA),
				V(P(4, 1, 1, 5), token.DEC_INT, "2"),
				T(P(5, 1, 1, 6), token.COMMA),
				V(P(7, 3, 1, 8), token.FLOAT, "3.0"),
				T(P(10, 1, 1, 11), token.COMMA),
				V(P(12, 5, 1, 13), token.RAW_STRING, "foo"),
				T(P(17, 1, 1, 18), token.COMMA),
				T(P(19, 1, 1, 20), token.COLON),
				V(P(20, 3, 1, 21), token.PUBLIC_IDENTIFIER, "bar"),
				T(P(23, 1, 1, 24), token.RBRACKET),
			},
		},
		"word array": {
			input: `%w[
foo bar  
baz]`,
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.WORD_LIST_BEG),
				V(P(4, 3, 2, 1), token.RAW_STRING, "foo"),
				V(P(8, 3, 2, 5), token.RAW_STRING, "bar"),
				V(P(14, 3, 3, 1), token.RAW_STRING, "baz"),
				T(P(17, 1, 3, 4), token.WORD_LIST_END),
			},
		},
		"symbol array": {
			input: "%s[foo bar   baz]",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.SYMBOL_LIST_BEG),
				V(P(3, 3, 1, 4), token.RAW_STRING, "foo"),
				V(P(7, 3, 1, 8), token.RAW_STRING, "bar"),
				V(P(13, 3, 1, 14), token.RAW_STRING, "baz"),
				T(P(16, 1, 1, 17), token.SYMBOL_LIST_END),
			},
		},
		"hex array": {
			input: "%x[ff 4_e   234]",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.HEX_LIST_BEG),
				V(P(3, 2, 1, 4), token.HEX_INT, "ff"),
				V(P(6, 3, 1, 7), token.HEX_INT, "4e"),
				V(P(12, 3, 1, 13), token.HEX_INT, "234"),
				T(P(15, 1, 1, 16), token.HEX_LIST_END),
			},
		},
		"hex array with invalid content": {
			input: "%x[ff 4ghij   234]",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.HEX_LIST_BEG),
				V(P(3, 2, 1, 4), token.HEX_INT, "ff"),
				V(P(6, 5, 1, 7), token.ERROR, "invalid int literal"),
				V(P(14, 3, 1, 15), token.HEX_INT, "234"),
				T(P(17, 1, 1, 18), token.HEX_LIST_END),
			},
		},
		"binary array": {
			input: "%b[11 1_0   101]",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.BIN_LIST_BEG),
				V(P(3, 2, 1, 4), token.BIN_INT, "11"),
				V(P(6, 3, 1, 7), token.BIN_INT, "10"),
				V(P(12, 3, 1, 13), token.BIN_INT, "101"),
				T(P(15, 1, 1, 16), token.BIN_LIST_END),
			},
		},
		"binary array with invalid content": {
			input: "%b[11 1ghij   101]",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.BIN_LIST_BEG),
				V(P(3, 2, 1, 4), token.BIN_INT, "11"),
				V(P(6, 5, 1, 7), token.ERROR, "invalid int literal"),
				V(P(14, 3, 1, 15), token.BIN_INT, "101"),
				T(P(17, 1, 1, 18), token.BIN_LIST_END),
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
				T(P(0, 2, 1, 1), token.SET_LITERAL_BEG),
				V(P(2, 1, 1, 3), token.DEC_INT, "1"),
				T(P(3, 1, 1, 4), token.COMMA),
				V(P(5, 1, 1, 6), token.DEC_INT, "2"),
				T(P(6, 1, 1, 7), token.COMMA),
				V(P(8, 3, 1, 9), token.FLOAT, "3.0"),
				T(P(11, 1, 1, 12), token.COMMA),
				V(P(13, 5, 1, 14), token.RAW_STRING, "foo"),
				T(P(18, 1, 1, 19), token.COMMA),
				T(P(20, 1, 1, 21), token.COLON),
				V(P(21, 3, 1, 22), token.PUBLIC_IDENTIFIER, "bar"),
				T(P(24, 1, 1, 25), token.RBRACE),
			},
		},
		"word set": {
			input: "%w{foo bar   baz}",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.WORD_SET_BEG),
				V(P(3, 3, 1, 4), token.RAW_STRING, "foo"),
				V(P(7, 3, 1, 8), token.RAW_STRING, "bar"),
				V(P(13, 3, 1, 14), token.RAW_STRING, "baz"),
				T(P(16, 1, 1, 17), token.WORD_SET_END),
			},
		},
		"symbol set": {
			input: "%s{foo bar   baz}",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.SYMBOL_SET_BEG),
				V(P(3, 3, 1, 4), token.RAW_STRING, "foo"),
				V(P(7, 3, 1, 8), token.RAW_STRING, "bar"),
				V(P(13, 3, 1, 14), token.RAW_STRING, "baz"),
				T(P(16, 1, 1, 17), token.SYMBOL_SET_END),
			},
		},
		"hex set": {
			input: "%x{ff 4_e   234}",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.HEX_SET_BEG),
				V(P(3, 2, 1, 4), token.HEX_INT, "ff"),
				V(P(6, 3, 1, 7), token.HEX_INT, "4e"),
				V(P(12, 3, 1, 13), token.HEX_INT, "234"),
				T(P(15, 1, 1, 16), token.HEX_SET_END),
			},
		},
		"hex set with invalid content": {
			input: "%x{ff 4ghij   234}",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.HEX_SET_BEG),
				V(P(3, 2, 1, 4), token.HEX_INT, "ff"),
				V(P(6, 5, 1, 7), token.ERROR, "invalid int literal"),
				V(P(14, 3, 1, 15), token.HEX_INT, "234"),
				T(P(17, 1, 1, 18), token.HEX_SET_END),
			},
		},
		"binary array": {
			input: "%b[11 1_0   101]",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.BIN_LIST_BEG),
				V(P(3, 2, 1, 4), token.BIN_INT, "11"),
				V(P(6, 3, 1, 7), token.BIN_INT, "10"),
				V(P(12, 3, 1, 13), token.BIN_INT, "101"),
				T(P(15, 1, 1, 16), token.BIN_LIST_END),
			},
		},
		"binary set with invalid content": {
			input: "%b{11 1ghij   101}",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.BIN_SET_BEG),
				V(P(3, 2, 1, 4), token.BIN_INT, "11"),
				V(P(6, 5, 1, 7), token.ERROR, "invalid int literal"),
				V(P(14, 3, 1, 15), token.BIN_INT, "101"),
				T(P(17, 1, 1, 18), token.BIN_SET_END),
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
				T(P(0, 2, 1, 1), token.TUPLE_LITERAL_BEG),
				V(P(2, 1, 1, 3), token.DEC_INT, "1"),
				T(P(3, 1, 1, 4), token.COMMA),
				V(P(5, 1, 1, 6), token.DEC_INT, "2"),
				T(P(6, 1, 1, 7), token.COMMA),
				V(P(8, 3, 1, 9), token.FLOAT, "3.0"),
				T(P(11, 1, 1, 12), token.COMMA),
				V(P(13, 5, 1, 14), token.RAW_STRING, "foo"),
				T(P(18, 1, 1, 19), token.COMMA),
				T(P(20, 1, 1, 21), token.COLON),
				V(P(21, 3, 1, 22), token.PUBLIC_IDENTIFIER, "bar"),
				T(P(24, 1, 1, 25), token.RPAREN),
			},
		},
		"word tuple": {
			input: "%w(foo bar   baz)",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.WORD_TUPLE_BEG),
				V(P(3, 3, 1, 4), token.RAW_STRING, "foo"),
				V(P(7, 3, 1, 8), token.RAW_STRING, "bar"),
				V(P(13, 3, 1, 14), token.RAW_STRING, "baz"),
				T(P(16, 1, 1, 17), token.WORD_TUPLE_END),
			},
		},
		"symbol tuple": {
			input: "%s(foo bar   baz)",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.SYMBOL_TUPLE_BEG),
				V(P(3, 3, 1, 4), token.RAW_STRING, "foo"),
				V(P(7, 3, 1, 8), token.RAW_STRING, "bar"),
				V(P(13, 3, 1, 14), token.RAW_STRING, "baz"),
				T(P(16, 1, 1, 17), token.SYMBOL_TUPLE_END),
			},
		},
		"hex tuple": {
			input: "%x(ff 4_e   234)",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.HEX_TUPLE_BEG),
				V(P(3, 2, 1, 4), token.HEX_INT, "ff"),
				V(P(6, 3, 1, 7), token.HEX_INT, "4e"),
				V(P(12, 3, 1, 13), token.HEX_INT, "234"),
				T(P(15, 1, 1, 16), token.HEX_TUPLE_END),
			},
		},
		"hex tuple with invalid content": {
			input: "%x(ff 4ghij   234)",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.HEX_TUPLE_BEG),
				V(P(3, 2, 1, 4), token.HEX_INT, "ff"),
				V(P(6, 5, 1, 7), token.ERROR, "invalid int literal"),
				V(P(14, 3, 1, 15), token.HEX_INT, "234"),
				T(P(17, 1, 1, 18), token.HEX_TUPLE_END),
			},
		},
		"binary tuple": {
			input: "%b(11 1_0   101)",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.BIN_TUPLE_BEG),
				V(P(3, 2, 1, 4), token.BIN_INT, "11"),
				V(P(6, 3, 1, 7), token.BIN_INT, "10"),
				V(P(12, 3, 1, 13), token.BIN_INT, "101"),
				T(P(15, 1, 1, 16), token.BIN_TUPLE_END),
			},
		},
		"binary tuple with invalid content": {
			input: "%b(11 1ghij   101)",
			want: []*token.Token{
				T(P(0, 3, 1, 1), token.BIN_TUPLE_BEG),
				V(P(3, 2, 1, 4), token.BIN_INT, "11"),
				V(P(6, 5, 1, 7), token.ERROR, "invalid int literal"),
				V(P(14, 3, 1, 15), token.BIN_INT, "101"),
				T(P(17, 1, 1, 18), token.BIN_TUPLE_END),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
