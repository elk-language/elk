package lexer

import "testing"

func TestCollectionLiteral(t *testing.T) {
	tests := testTable{
		"incorrect word delimiters": {
			input: "%w<foo bar   baz>",
			want: []*Token{
				V(ErrorToken, "invalid word collection literal delimiters `%w`", 0, 2, 1, 1),
				T(LessToken, 2, 1, 1, 3),
				V(PublicIdentifierToken, "foo", 3, 3, 1, 4),
				V(PublicIdentifierToken, "bar", 7, 3, 1, 8),
				V(PublicIdentifierToken, "baz", 13, 3, 1, 14),
				T(GreaterToken, 16, 1, 1, 17),
			},
		},
		"incorrect symbol delimiters": {
			input: "%s<foo bar   baz>",
			want: []*Token{
				V(ErrorToken, "invalid symbol collection literal delimiters `%s`", 0, 2, 1, 1),
				T(LessToken, 2, 1, 1, 3),
				V(PublicIdentifierToken, "foo", 3, 3, 1, 4),
				V(PublicIdentifierToken, "bar", 7, 3, 1, 8),
				V(PublicIdentifierToken, "baz", 13, 3, 1, 14),
				T(GreaterToken, 16, 1, 1, 17),
			},
		},
		"incorrect hex delimiters": {
			input: "%x<45a 101   fff>",
			want: []*Token{
				V(ErrorToken, "invalid hex collection literal delimiters `%x`", 0, 2, 1, 1),
				T(LessToken, 2, 1, 1, 3),
				V(DecIntToken, "45", 3, 2, 1, 4),
				V(PublicIdentifierToken, "a", 5, 1, 1, 6),
				V(DecIntToken, "101", 7, 3, 1, 8),
				V(PublicIdentifierToken, "fff", 13, 3, 1, 14),
				T(GreaterToken, 16, 1, 1, 17),
			},
		},
		"incorrect binary delimiters": {
			input: "%b<110 101   111>",
			want: []*Token{
				V(ErrorToken, "invalid binary collection literal delimiters `%b`", 0, 2, 1, 1),
				T(LessToken, 2, 1, 1, 3),
				V(DecIntToken, "110", 3, 3, 1, 4),
				V(DecIntToken, "101", 7, 3, 1, 8),
				V(DecIntToken, "111", 13, 3, 1, 14),
				T(GreaterToken, 16, 1, 1, 17),
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
			want: []*Token{
				T(LBracketToken, 0, 1, 1, 1),
				V(DecIntToken, "1", 1, 1, 1, 2),
				T(CommaToken, 2, 1, 1, 3),
				V(DecIntToken, "2", 4, 1, 1, 5),
				T(CommaToken, 5, 1, 1, 6),
				V(FloatToken, "3.0", 7, 3, 1, 8),
				T(CommaToken, 10, 1, 1, 11),
				V(RawStringToken, "foo", 12, 5, 1, 13),
				T(CommaToken, 17, 1, 1, 18),
				T(SymbolBegToken, 19, 1, 1, 20),
				V(PublicIdentifierToken, "bar", 20, 3, 1, 21),
				T(RBracketToken, 23, 1, 1, 24),
			},
		},
		"word array": {
			input: `%w[
foo bar  
baz]`,
			want: []*Token{
				T(WordArrayBegToken, 0, 3, 1, 1),
				V(RawStringToken, "foo", 4, 3, 2, 1),
				V(RawStringToken, "bar", 8, 3, 2, 5),
				V(RawStringToken, "baz", 14, 3, 3, 1),
				T(WordArrayEndToken, 17, 1, 3, 4),
			},
		},
		"symbol array": {
			input: "%s[foo bar   baz]",
			want: []*Token{
				T(SymbolArrayBegToken, 0, 3, 1, 1),
				V(RawStringToken, "foo", 3, 3, 1, 4),
				V(RawStringToken, "bar", 7, 3, 1, 8),
				V(RawStringToken, "baz", 13, 3, 1, 14),
				T(SymbolArrayEndToken, 16, 1, 1, 17),
			},
		},
		"hex array": {
			input: "%x[ff 4_e   234]",
			want: []*Token{
				T(HexArrayBegToken, 0, 3, 1, 1),
				V(HexIntToken, "ff", 3, 2, 1, 4),
				V(HexIntToken, "4e", 6, 3, 1, 7),
				V(HexIntToken, "234", 12, 3, 1, 13),
				T(HexArrayEndToken, 15, 1, 1, 16),
			},
		},
		"hex array with invalid content": {
			input: "%x[ff 4ghij   234]",
			want: []*Token{
				T(HexArrayBegToken, 0, 3, 1, 1),
				V(HexIntToken, "ff", 3, 2, 1, 4),
				V(ErrorToken, "invalid int literal", 6, 5, 1, 7),
				V(HexIntToken, "234", 14, 3, 1, 15),
				T(HexArrayEndToken, 17, 1, 1, 18),
			},
		},
		"binary array": {
			input: "%b[11 1_0   101]",
			want: []*Token{
				T(BinArrayBegToken, 0, 3, 1, 1),
				V(BinIntToken, "11", 3, 2, 1, 4),
				V(BinIntToken, "10", 6, 3, 1, 7),
				V(BinIntToken, "101", 12, 3, 1, 13),
				T(BinArrayEndToken, 15, 1, 1, 16),
			},
		},
		"binary array with invalid content": {
			input: "%b[11 1ghij   101]",
			want: []*Token{
				T(BinArrayBegToken, 0, 3, 1, 1),
				V(BinIntToken, "11", 3, 2, 1, 4),
				V(ErrorToken, "invalid int literal", 6, 5, 1, 7),
				V(BinIntToken, "101", 14, 3, 1, 15),
				T(BinArrayEndToken, 17, 1, 1, 18),
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
			want: []*Token{
				T(SetLiteralBegToken, 0, 2, 1, 1),
				V(DecIntToken, "1", 2, 1, 1, 3),
				T(CommaToken, 3, 1, 1, 4),
				V(DecIntToken, "2", 5, 1, 1, 6),
				T(CommaToken, 6, 1, 1, 7),
				V(FloatToken, "3.0", 8, 3, 1, 9),
				T(CommaToken, 11, 1, 1, 12),
				V(RawStringToken, "foo", 13, 5, 1, 14),
				T(CommaToken, 18, 1, 1, 19),
				T(SymbolBegToken, 20, 1, 1, 21),
				V(PublicIdentifierToken, "bar", 21, 3, 1, 22),
				T(RBraceToken, 24, 1, 1, 25),
			},
		},
		"word set": {
			input: "%w{foo bar   baz}",
			want: []*Token{
				T(WordSetBegToken, 0, 3, 1, 1),
				V(RawStringToken, "foo", 3, 3, 1, 4),
				V(RawStringToken, "bar", 7, 3, 1, 8),
				V(RawStringToken, "baz", 13, 3, 1, 14),
				T(WordSetEndToken, 16, 1, 1, 17),
			},
		},
		"symbol set": {
			input: "%s{foo bar   baz}",
			want: []*Token{
				T(SymbolSetBegToken, 0, 3, 1, 1),
				V(RawStringToken, "foo", 3, 3, 1, 4),
				V(RawStringToken, "bar", 7, 3, 1, 8),
				V(RawStringToken, "baz", 13, 3, 1, 14),
				T(SymbolSetEndToken, 16, 1, 1, 17),
			},
		},
		"hex set": {
			input: "%x{ff 4_e   234}",
			want: []*Token{
				T(HexSetBegToken, 0, 3, 1, 1),
				V(HexIntToken, "ff", 3, 2, 1, 4),
				V(HexIntToken, "4e", 6, 3, 1, 7),
				V(HexIntToken, "234", 12, 3, 1, 13),
				T(HexSetEndToken, 15, 1, 1, 16),
			},
		},
		"hex set with invalid content": {
			input: "%x{ff 4ghij   234}",
			want: []*Token{
				T(HexSetBegToken, 0, 3, 1, 1),
				V(HexIntToken, "ff", 3, 2, 1, 4),
				V(ErrorToken, "invalid int literal", 6, 5, 1, 7),
				V(HexIntToken, "234", 14, 3, 1, 15),
				T(HexSetEndToken, 17, 1, 1, 18),
			},
		},
		"binary array": {
			input: "%b[11 1_0   101]",
			want: []*Token{
				T(BinArrayBegToken, 0, 3, 1, 1),
				V(BinIntToken, "11", 3, 2, 1, 4),
				V(BinIntToken, "10", 6, 3, 1, 7),
				V(BinIntToken, "101", 12, 3, 1, 13),
				T(BinArrayEndToken, 15, 1, 1, 16),
			},
		},
		"binary set with invalid content": {
			input: "%b{11 1ghij   101}",
			want: []*Token{
				T(BinSetBegToken, 0, 3, 1, 1),
				V(BinIntToken, "11", 3, 2, 1, 4),
				V(ErrorToken, "invalid int literal", 6, 5, 1, 7),
				V(BinIntToken, "101", 14, 3, 1, 15),
				T(BinSetEndToken, 17, 1, 1, 18),
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
			want: []*Token{
				T(TupleLiteralBegToken, 0, 2, 1, 1),
				V(DecIntToken, "1", 2, 1, 1, 3),
				T(CommaToken, 3, 1, 1, 4),
				V(DecIntToken, "2", 5, 1, 1, 6),
				T(CommaToken, 6, 1, 1, 7),
				V(FloatToken, "3.0", 8, 3, 1, 9),
				T(CommaToken, 11, 1, 1, 12),
				V(RawStringToken, "foo", 13, 5, 1, 14),
				T(CommaToken, 18, 1, 1, 19),
				T(SymbolBegToken, 20, 1, 1, 21),
				V(PublicIdentifierToken, "bar", 21, 3, 1, 22),
				T(RParenToken, 24, 1, 1, 25),
			},
		},
		"word tuple": {
			input: "%w(foo bar   baz)",
			want: []*Token{
				T(WordTupleBegToken, 0, 3, 1, 1),
				V(RawStringToken, "foo", 3, 3, 1, 4),
				V(RawStringToken, "bar", 7, 3, 1, 8),
				V(RawStringToken, "baz", 13, 3, 1, 14),
				T(WordTupleEndToken, 16, 1, 1, 17),
			},
		},
		"symbol tuple": {
			input: "%s(foo bar   baz)",
			want: []*Token{
				T(SymbolTupleBegToken, 0, 3, 1, 1),
				V(RawStringToken, "foo", 3, 3, 1, 4),
				V(RawStringToken, "bar", 7, 3, 1, 8),
				V(RawStringToken, "baz", 13, 3, 1, 14),
				T(SymbolTupleEndToken, 16, 1, 1, 17),
			},
		},
		"hex tuple": {
			input: "%x(ff 4_e   234)",
			want: []*Token{
				T(HexTupleBegToken, 0, 3, 1, 1),
				V(HexIntToken, "ff", 3, 2, 1, 4),
				V(HexIntToken, "4e", 6, 3, 1, 7),
				V(HexIntToken, "234", 12, 3, 1, 13),
				T(HexTupleEndToken, 15, 1, 1, 16),
			},
		},
		"hex tuple with invalid content": {
			input: "%x(ff 4ghij   234)",
			want: []*Token{
				T(HexTupleBegToken, 0, 3, 1, 1),
				V(HexIntToken, "ff", 3, 2, 1, 4),
				V(ErrorToken, "invalid int literal", 6, 5, 1, 7),
				V(HexIntToken, "234", 14, 3, 1, 15),
				T(HexTupleEndToken, 17, 1, 1, 18),
			},
		},
		"binary tuple": {
			input: "%b(11 1_0   101)",
			want: []*Token{
				T(BinTupleBegToken, 0, 3, 1, 1),
				V(BinIntToken, "11", 3, 2, 1, 4),
				V(BinIntToken, "10", 6, 3, 1, 7),
				V(BinIntToken, "101", 12, 3, 1, 13),
				T(BinTupleEndToken, 15, 1, 1, 16),
			},
		},
		"binary tuple with invalid content": {
			input: "%b(11 1ghij   101)",
			want: []*Token{
				T(BinTupleBegToken, 0, 3, 1, 1),
				V(BinIntToken, "11", 3, 2, 1, 4),
				V(ErrorToken, "invalid int literal", 6, 5, 1, 7),
				V(BinIntToken, "101", 14, 3, 1, 15),
				T(BinTupleEndToken, 17, 1, 1, 18),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
