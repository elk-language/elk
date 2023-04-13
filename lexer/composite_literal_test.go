package lexer

import "testing"

func TestCollectionLiteral(t *testing.T) {
	tests := testTable{
		"incorrect delimiters": {
			input: "%w<foo bar   baz>",
			want: []*Token{
				{
					TokenType:  ErrorToken,
					Value:      "invalid word collection literal delimiters `%%w`",
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  LessToken,
					StartByte:  2,
					ByteLength: 1,
					Line:       1,
					Column:     3,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "foo",
					StartByte:  3,
					ByteLength: 3,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "bar",
					StartByte:  7,
					ByteLength: 3,
					Line:       1,
					Column:     8,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "baz",
					StartByte:  13,
					ByteLength: 3,
					Line:       1,
					Column:     14,
				},
				{
					TokenType:  GreaterToken,
					StartByte:  16,
					ByteLength: 1,
					Line:       1,
					Column:     17,
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

func TestArray(t *testing.T) {
	tests := testTable{
		"regular array": {
			input: "[1, 2, 3.0, 'foo', :bar]",
			want: []*Token{
				{
					TokenType:  LBracketToken,
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IntToken,
					Value:      "1",
					StartByte:  1,
					ByteLength: 1,
					Line:       1,
					Column:     2,
				},
				{
					TokenType:  CommaToken,
					StartByte:  2,
					ByteLength: 1,
					Line:       1,
					Column:     3,
				},
				{
					TokenType:  IntToken,
					Value:      "2",
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
				},
				{
					TokenType:  CommaToken,
					StartByte:  5,
					ByteLength: 1,
					Line:       1,
					Column:     6,
				},
				{
					TokenType:  FloatToken,
					Value:      "3.0",
					StartByte:  7,
					ByteLength: 3,
					Line:       1,
					Column:     8,
				},
				{
					TokenType:  CommaToken,
					StartByte:  10,
					ByteLength: 1,
					Line:       1,
					Column:     11,
				},
				{
					TokenType:  RawStringToken,
					Value:      "foo",
					StartByte:  12,
					ByteLength: 5,
					Line:       1,
					Column:     13,
				},
				{
					TokenType:  CommaToken,
					StartByte:  17,
					ByteLength: 1,
					Line:       1,
					Column:     18,
				},
				{
					TokenType:  SymbolBegToken,
					StartByte:  19,
					ByteLength: 1,
					Line:       1,
					Column:     20,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "bar",
					StartByte:  20,
					ByteLength: 3,
					Line:       1,
					Column:     21,
				},
				{
					TokenType:  RBracketToken,
					StartByte:  23,
					ByteLength: 1,
					Line:       1,
					Column:     24,
				},
			},
		},
		"word array": {
			input: `%w[
foo bar  
baz]`,
			want: []*Token{
				{
					TokenType:  WordArrayBegToken,
					StartByte:  0,
					ByteLength: 3,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  RawStringToken,
					Value:      "foo",
					StartByte:  4,
					ByteLength: 3,
					Line:       2,
					Column:     1,
				},
				{
					TokenType:  RawStringToken,
					Value:      "bar",
					StartByte:  8,
					ByteLength: 3,
					Line:       2,
					Column:     5,
				},
				{
					TokenType:  RawStringToken,
					Value:      "baz",
					StartByte:  14,
					ByteLength: 3,
					Line:       3,
					Column:     1,
				},
				{
					TokenType:  WordArrayEndToken,
					StartByte:  17,
					ByteLength: 1,
					Line:       3,
					Column:     4,
				},
			},
		},
		"symbol array": {
			input: "%s[foo bar   baz]",
			want: []*Token{
				{
					TokenType:  SymbolArrayBegToken,
					StartByte:  0,
					ByteLength: 3,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  RawStringToken,
					Value:      "foo",
					StartByte:  3,
					ByteLength: 3,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  RawStringToken,
					Value:      "bar",
					StartByte:  7,
					ByteLength: 3,
					Line:       1,
					Column:     8,
				},
				{
					TokenType:  RawStringToken,
					Value:      "baz",
					StartByte:  13,
					ByteLength: 3,
					Line:       1,
					Column:     14,
				},
				{
					TokenType:  SymbolArrayEndToken,
					StartByte:  16,
					ByteLength: 1,
					Line:       1,
					Column:     17,
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

func TestSet(t *testing.T) {
	tests := testTable{
		"regular set": {
			input: "%{1, 2, 3.0, 'foo', :bar}",
			want: []*Token{
				{
					TokenType:  SetLiteralBegToken,
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IntToken,
					Value:      "1",
					StartByte:  2,
					ByteLength: 1,
					Line:       1,
					Column:     3,
				},
				{
					TokenType:  CommaToken,
					StartByte:  3,
					ByteLength: 1,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  IntToken,
					Value:      "2",
					StartByte:  5,
					ByteLength: 1,
					Line:       1,
					Column:     6,
				},
				{
					TokenType:  CommaToken,
					StartByte:  6,
					ByteLength: 1,
					Line:       1,
					Column:     7,
				},
				{
					TokenType:  FloatToken,
					Value:      "3.0",
					StartByte:  8,
					ByteLength: 3,
					Line:       1,
					Column:     9,
				},
				{
					TokenType:  CommaToken,
					StartByte:  11,
					ByteLength: 1,
					Line:       1,
					Column:     12,
				},
				{
					TokenType:  RawStringToken,
					Value:      "foo",
					StartByte:  13,
					ByteLength: 5,
					Line:       1,
					Column:     14,
				},
				{
					TokenType:  CommaToken,
					StartByte:  18,
					ByteLength: 1,
					Line:       1,
					Column:     19,
				},
				{
					TokenType:  SymbolBegToken,
					StartByte:  20,
					ByteLength: 1,
					Line:       1,
					Column:     21,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "bar",
					StartByte:  21,
					ByteLength: 3,
					Line:       1,
					Column:     22,
				},
				{
					TokenType:  RBraceToken,
					StartByte:  24,
					ByteLength: 1,
					Line:       1,
					Column:     25,
				},
			},
		},
		"word set": {
			input: "%w{foo bar   baz}",
			want: []*Token{
				{
					TokenType:  WordSetBegToken,
					StartByte:  0,
					ByteLength: 3,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  RawStringToken,
					Value:      "foo",
					StartByte:  3,
					ByteLength: 3,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  RawStringToken,
					Value:      "bar",
					StartByte:  7,
					ByteLength: 3,
					Line:       1,
					Column:     8,
				},
				{
					TokenType:  RawStringToken,
					Value:      "baz",
					StartByte:  13,
					ByteLength: 3,
					Line:       1,
					Column:     14,
				},
				{
					TokenType:  WordSetEndToken,
					StartByte:  16,
					ByteLength: 1,
					Line:       1,
					Column:     17,
				},
			},
		},
		"symbol set": {
			input: "%s{foo bar   baz}",
			want: []*Token{
				{
					TokenType:  SymbolSetBegToken,
					StartByte:  0,
					ByteLength: 3,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  RawStringToken,
					Value:      "foo",
					StartByte:  3,
					ByteLength: 3,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  RawStringToken,
					Value:      "bar",
					StartByte:  7,
					ByteLength: 3,
					Line:       1,
					Column:     8,
				},
				{
					TokenType:  RawStringToken,
					Value:      "baz",
					StartByte:  13,
					ByteLength: 3,
					Line:       1,
					Column:     14,
				},
				{
					TokenType:  SymbolSetEndToken,
					StartByte:  16,
					ByteLength: 1,
					Line:       1,
					Column:     17,
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

func TestTuple(t *testing.T) {
	tests := testTable{
		"regular tuple": {
			input: "%(1, 2, 3.0, 'foo', :bar)",
			want: []*Token{
				{
					TokenType:  TupleLiteralBegToken,
					StartByte:  0,
					ByteLength: 2,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IntToken,
					Value:      "1",
					StartByte:  2,
					ByteLength: 1,
					Line:       1,
					Column:     3,
				},
				{
					TokenType:  CommaToken,
					StartByte:  3,
					ByteLength: 1,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  IntToken,
					Value:      "2",
					StartByte:  5,
					ByteLength: 1,
					Line:       1,
					Column:     6,
				},
				{
					TokenType:  CommaToken,
					StartByte:  6,
					ByteLength: 1,
					Line:       1,
					Column:     7,
				},
				{
					TokenType:  FloatToken,
					Value:      "3.0",
					StartByte:  8,
					ByteLength: 3,
					Line:       1,
					Column:     9,
				},
				{
					TokenType:  CommaToken,
					StartByte:  11,
					ByteLength: 1,
					Line:       1,
					Column:     12,
				},
				{
					TokenType:  RawStringToken,
					Value:      "foo",
					StartByte:  13,
					ByteLength: 5,
					Line:       1,
					Column:     14,
				},
				{
					TokenType:  CommaToken,
					StartByte:  18,
					ByteLength: 1,
					Line:       1,
					Column:     19,
				},
				{
					TokenType:  SymbolBegToken,
					StartByte:  20,
					ByteLength: 1,
					Line:       1,
					Column:     21,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "bar",
					StartByte:  21,
					ByteLength: 3,
					Line:       1,
					Column:     22,
				},
				{
					TokenType:  RParenToken,
					StartByte:  24,
					ByteLength: 1,
					Line:       1,
					Column:     25,
				},
			},
		},
		"word tuple": {
			input: "%w(foo bar   baz)",
			want: []*Token{
				{
					TokenType:  WordTupleBegToken,
					StartByte:  0,
					ByteLength: 3,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  RawStringToken,
					Value:      "foo",
					StartByte:  3,
					ByteLength: 3,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  RawStringToken,
					Value:      "bar",
					StartByte:  7,
					ByteLength: 3,
					Line:       1,
					Column:     8,
				},
				{
					TokenType:  RawStringToken,
					Value:      "baz",
					StartByte:  13,
					ByteLength: 3,
					Line:       1,
					Column:     14,
				},
				{
					TokenType:  WordTupleEndToken,
					StartByte:  16,
					ByteLength: 1,
					Line:       1,
					Column:     17,
				},
			},
		},
		"symbol tuple": {
			input: "%s(foo bar   baz)",
			want: []*Token{
				{
					TokenType:  SymbolTupleBegToken,
					StartByte:  0,
					ByteLength: 3,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  RawStringToken,
					Value:      "foo",
					StartByte:  3,
					ByteLength: 3,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  RawStringToken,
					Value:      "bar",
					StartByte:  7,
					ByteLength: 3,
					Line:       1,
					Column:     8,
				},
				{
					TokenType:  RawStringToken,
					Value:      "baz",
					StartByte:  13,
					ByteLength: 3,
					Line:       1,
					Column:     14,
				},
				{
					TokenType:  SymbolTupleEndToken,
					StartByte:  16,
					ByteLength: 1,
					Line:       1,
					Column:     17,
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
