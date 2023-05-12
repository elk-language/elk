package lexer

import "testing"

func TestCollectionLiteral(t *testing.T) {
	tests := testTable{
		"incorrect word delimiters": {
			input: "%w<foo bar   baz>",
			want: []*Token{
				V(P(0, 2, 1, 1), ErrorToken, "invalid word collection literal delimiters `%w`"),
				T(P(2, 1, 1, 3), LessToken),
				V(P(3, 3, 1, 4), PublicIdentifierToken, "foo"),
				V(P(7, 3, 1, 8), PublicIdentifierToken, "bar"),
				V(P(13, 3, 1, 14), PublicIdentifierToken, "baz"),
				T(P(16, 1, 1, 17), GreaterToken),
			},
		},
		"incorrect symbol delimiters": {
			input: "%s<foo bar   baz>",
			want: []*Token{
				V(P(0, 2, 1, 1), ErrorToken, "invalid symbol collection literal delimiters `%s`"),
				T(P(2, 1, 1, 3), LessToken),
				V(P(3, 3, 1, 4), PublicIdentifierToken, "foo"),
				V(P(7, 3, 1, 8), PublicIdentifierToken, "bar"),
				V(P(13, 3, 1, 14), PublicIdentifierToken, "baz"),
				T(P(16, 1, 1, 17), GreaterToken),
			},
		},
		"incorrect hex delimiters": {
			input: "%x<45a 101   fff>",
			want: []*Token{
				V(P(0, 2, 1, 1), ErrorToken, "invalid hex collection literal delimiters `%x`"),
				T(P(2, 1, 1, 3), LessToken),
				V(P(3, 2, 1, 4), DecIntToken, "45"),
				V(P(5, 1, 1, 6), PublicIdentifierToken, "a"),
				V(P(7, 3, 1, 8), DecIntToken, "101"),
				V(P(13, 3, 1, 14), PublicIdentifierToken, "fff"),
				T(P(16, 1, 1, 17), GreaterToken),
			},
		},
		"incorrect binary delimiters": {
			input: "%b<110 101   111>",
			want: []*Token{
				V(P(0, 2, 1, 1), ErrorToken, "invalid binary collection literal delimiters `%b`"),
				T(P(2, 1, 1, 3), LessToken),
				V(P(3, 3, 1, 4), DecIntToken, "110"),
				V(P(7, 3, 1, 8), DecIntToken, "101"),
				V(P(13, 3, 1, 14), DecIntToken, "111"),
				T(P(16, 1, 1, 17), GreaterToken),
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
				T(P(0, 1, 1, 1), LBracketToken),
				V(P(1, 1, 1, 2), DecIntToken, "1"),
				T(P(2, 1, 1, 3), CommaToken),
				V(P(4, 1, 1, 5), DecIntToken, "2"),
				T(P(5, 1, 1, 6), CommaToken),
				V(P(7, 3, 1, 8), FloatToken, "3.0"),
				T(P(10, 1, 1, 11), CommaToken),
				V(P(12, 5, 1, 13), RawStringToken, "foo"),
				T(P(17, 1, 1, 18), CommaToken),
				T(P(19, 1, 1, 20), SymbolBegToken),
				V(P(20, 3, 1, 21), PublicIdentifierToken, "bar"),
				T(P(23, 1, 1, 24), RBracketToken),
			},
		},
		"word array": {
			input: `%w[
foo bar  
baz]`,
			want: []*Token{
				T(P(0, 3, 1, 1), WordArrayBegToken),
				V(P(4, 3, 2, 1), RawStringToken, "foo"),
				V(P(8, 3, 2, 5), RawStringToken, "bar"),
				V(P(14, 3, 3, 1), RawStringToken, "baz"),
				T(P(17, 1, 3, 4), WordArrayEndToken),
			},
		},
		"symbol array": {
			input: "%s[foo bar   baz]",
			want: []*Token{
				T(P(0, 3, 1, 1), SymbolArrayBegToken),
				V(P(3, 3, 1, 4), RawStringToken, "foo"),
				V(P(7, 3, 1, 8), RawStringToken, "bar"),
				V(P(13, 3, 1, 14), RawStringToken, "baz"),
				T(P(16, 1, 1, 17), SymbolArrayEndToken),
			},
		},
		"hex array": {
			input: "%x[ff 4_e   234]",
			want: []*Token{
				T(P(0, 3, 1, 1), HexArrayBegToken),
				V(P(3, 2, 1, 4), HexIntToken, "ff"),
				V(P(6, 3, 1, 7), HexIntToken, "4e"),
				V(P(12, 3, 1, 13), HexIntToken, "234"),
				T(P(15, 1, 1, 16), HexArrayEndToken),
			},
		},
		"hex array with invalid content": {
			input: "%x[ff 4ghij   234]",
			want: []*Token{
				T(P(0, 3, 1, 1), HexArrayBegToken),
				V(P(3, 2, 1, 4), HexIntToken, "ff"),
				V(P(6, 5, 1, 7), ErrorToken, "invalid int literal"),
				V(P(14, 3, 1, 15), HexIntToken, "234"),
				T(P(17, 1, 1, 18), HexArrayEndToken),
			},
		},
		"binary array": {
			input: "%b[11 1_0   101]",
			want: []*Token{
				T(P(0, 3, 1, 1), BinArrayBegToken),
				V(P(3, 2, 1, 4), BinIntToken, "11"),
				V(P(6, 3, 1, 7), BinIntToken, "10"),
				V(P(12, 3, 1, 13), BinIntToken, "101"),
				T(P(15, 1, 1, 16), BinArrayEndToken),
			},
		},
		"binary array with invalid content": {
			input: "%b[11 1ghij   101]",
			want: []*Token{
				T(P(0, 3, 1, 1), BinArrayBegToken),
				V(P(3, 2, 1, 4), BinIntToken, "11"),
				V(P(6, 5, 1, 7), ErrorToken, "invalid int literal"),
				V(P(14, 3, 1, 15), BinIntToken, "101"),
				T(P(17, 1, 1, 18), BinArrayEndToken),
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
				T(P(0, 2, 1, 1), SetLiteralBegToken),
				V(P(2, 1, 1, 3), DecIntToken, "1"),
				T(P(3, 1, 1, 4), CommaToken),
				V(P(5, 1, 1, 6), DecIntToken, "2"),
				T(P(6, 1, 1, 7), CommaToken),
				V(P(8, 3, 1, 9), FloatToken, "3.0"),
				T(P(11, 1, 1, 12), CommaToken),
				V(P(13, 5, 1, 14), RawStringToken, "foo"),
				T(P(18, 1, 1, 19), CommaToken),
				T(P(20, 1, 1, 21), SymbolBegToken),
				V(P(21, 3, 1, 22), PublicIdentifierToken, "bar"),
				T(P(24, 1, 1, 25), RBraceToken),
			},
		},
		"word set": {
			input: "%w{foo bar   baz}",
			want: []*Token{
				T(P(0, 3, 1, 1), WordSetBegToken),
				V(P(3, 3, 1, 4), RawStringToken, "foo"),
				V(P(7, 3, 1, 8), RawStringToken, "bar"),
				V(P(13, 3, 1, 14), RawStringToken, "baz"),
				T(P(16, 1, 1, 17), WordSetEndToken),
			},
		},
		"symbol set": {
			input: "%s{foo bar   baz}",
			want: []*Token{
				T(P(0, 3, 1, 1), SymbolSetBegToken),
				V(P(3, 3, 1, 4), RawStringToken, "foo"),
				V(P(7, 3, 1, 8), RawStringToken, "bar"),
				V(P(13, 3, 1, 14), RawStringToken, "baz"),
				T(P(16, 1, 1, 17), SymbolSetEndToken),
			},
		},
		"hex set": {
			input: "%x{ff 4_e   234}",
			want: []*Token{
				T(P(0, 3, 1, 1), HexSetBegToken),
				V(P(3, 2, 1, 4), HexIntToken, "ff"),
				V(P(6, 3, 1, 7), HexIntToken, "4e"),
				V(P(12, 3, 1, 13), HexIntToken, "234"),
				T(P(15, 1, 1, 16), HexSetEndToken),
			},
		},
		"hex set with invalid content": {
			input: "%x{ff 4ghij   234}",
			want: []*Token{
				T(P(0, 3, 1, 1), HexSetBegToken),
				V(P(3, 2, 1, 4), HexIntToken, "ff"),
				V(P(6, 5, 1, 7), ErrorToken, "invalid int literal"),
				V(P(14, 3, 1, 15), HexIntToken, "234"),
				T(P(17, 1, 1, 18), HexSetEndToken),
			},
		},
		"binary array": {
			input: "%b[11 1_0   101]",
			want: []*Token{
				T(P(0, 3, 1, 1), BinArrayBegToken),
				V(P(3, 2, 1, 4), BinIntToken, "11"),
				V(P(6, 3, 1, 7), BinIntToken, "10"),
				V(P(12, 3, 1, 13), BinIntToken, "101"),
				T(P(15, 1, 1, 16), BinArrayEndToken),
			},
		},
		"binary set with invalid content": {
			input: "%b{11 1ghij   101}",
			want: []*Token{
				T(P(0, 3, 1, 1), BinSetBegToken),
				V(P(3, 2, 1, 4), BinIntToken, "11"),
				V(P(6, 5, 1, 7), ErrorToken, "invalid int literal"),
				V(P(14, 3, 1, 15), BinIntToken, "101"),
				T(P(17, 1, 1, 18), BinSetEndToken),
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
				T(P(0, 2, 1, 1), TupleLiteralBegToken),
				V(P(2, 1, 1, 3), DecIntToken, "1"),
				T(P(3, 1, 1, 4), CommaToken),
				V(P(5, 1, 1, 6), DecIntToken, "2"),
				T(P(6, 1, 1, 7), CommaToken),
				V(P(8, 3, 1, 9), FloatToken, "3.0"),
				T(P(11, 1, 1, 12), CommaToken),
				V(P(13, 5, 1, 14), RawStringToken, "foo"),
				T(P(18, 1, 1, 19), CommaToken),
				T(P(20, 1, 1, 21), SymbolBegToken),
				V(P(21, 3, 1, 22), PublicIdentifierToken, "bar"),
				T(P(24, 1, 1, 25), RParenToken),
			},
		},
		"word tuple": {
			input: "%w(foo bar   baz)",
			want: []*Token{
				T(P(0, 3, 1, 1), WordTupleBegToken),
				V(P(3, 3, 1, 4), RawStringToken, "foo"),
				V(P(7, 3, 1, 8), RawStringToken, "bar"),
				V(P(13, 3, 1, 14), RawStringToken, "baz"),
				T(P(16, 1, 1, 17), WordTupleEndToken),
			},
		},
		"symbol tuple": {
			input: "%s(foo bar   baz)",
			want: []*Token{
				T(P(0, 3, 1, 1), SymbolTupleBegToken),
				V(P(3, 3, 1, 4), RawStringToken, "foo"),
				V(P(7, 3, 1, 8), RawStringToken, "bar"),
				V(P(13, 3, 1, 14), RawStringToken, "baz"),
				T(P(16, 1, 1, 17), SymbolTupleEndToken),
			},
		},
		"hex tuple": {
			input: "%x(ff 4_e   234)",
			want: []*Token{
				T(P(0, 3, 1, 1), HexTupleBegToken),
				V(P(3, 2, 1, 4), HexIntToken, "ff"),
				V(P(6, 3, 1, 7), HexIntToken, "4e"),
				V(P(12, 3, 1, 13), HexIntToken, "234"),
				T(P(15, 1, 1, 16), HexTupleEndToken),
			},
		},
		"hex tuple with invalid content": {
			input: "%x(ff 4ghij   234)",
			want: []*Token{
				T(P(0, 3, 1, 1), HexTupleBegToken),
				V(P(3, 2, 1, 4), HexIntToken, "ff"),
				V(P(6, 5, 1, 7), ErrorToken, "invalid int literal"),
				V(P(14, 3, 1, 15), HexIntToken, "234"),
				T(P(17, 1, 1, 18), HexTupleEndToken),
			},
		},
		"binary tuple": {
			input: "%b(11 1_0   101)",
			want: []*Token{
				T(P(0, 3, 1, 1), BinTupleBegToken),
				V(P(3, 2, 1, 4), BinIntToken, "11"),
				V(P(6, 3, 1, 7), BinIntToken, "10"),
				V(P(12, 3, 1, 13), BinIntToken, "101"),
				T(P(15, 1, 1, 16), BinTupleEndToken),
			},
		},
		"binary tuple with invalid content": {
			input: "%b(11 1ghij   101)",
			want: []*Token{
				T(P(0, 3, 1, 1), BinTupleBegToken),
				V(P(3, 2, 1, 4), BinIntToken, "11"),
				V(P(6, 5, 1, 7), ErrorToken, "invalid int literal"),
				V(P(14, 3, 1, 15), BinIntToken, "101"),
				T(P(17, 1, 1, 18), BinTupleEndToken),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
