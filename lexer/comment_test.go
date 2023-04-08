package lexer

import "testing"

func TestSingleLineComment(t *testing.T) {
	tests := testTable{
		"discards characters until a new line is reached": {
			input: `3 + # 25 / 3
							5`,
			want: []*Token{
				{
					TokenType:  IntToken,
					Value:      "3",
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  2,
					ByteLength: 1,
					Line:       1,
					Column:     3,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  12,
					ByteLength: 1,
					Line:       1,
					Column:     13,
				},
				{
					TokenType:  IntToken,
					Value:      "5",
					StartByte:  20,
					ByteLength: 1,
					Line:       2,
					Column:     8,
				},
			},
		},
		"can appear at the beginning of the line": {
			input: `# something awesome
							foo := 3`,
			want: []*Token{
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  19,
					ByteLength: 1,
					Line:       1,
					Column:     20,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "foo",
					StartByte:  27,
					ByteLength: 3,
					Line:       2,
					Column:     8,
				},
				{
					TokenType:  ColonEqualToken,
					Value:      "",
					StartByte:  31,
					ByteLength: 2,
					Line:       2,
					Column:     12,
				},
				{
					TokenType:  IntToken,
					Value:      "3",
					StartByte:  34,
					ByteLength: 1,
					Line:       2,
					Column:     15,
				},
			},
		},
		"can appear on consecutive lines": {
			input: `
# peace
# and love
# from Elk
println 'Hey'
			`,
			want: []*Token{
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  8,
					ByteLength: 1,
					Line:       2,
					Column:     8,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  19,
					ByteLength: 1,
					Line:       3,
					Column:     11,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  30,
					ByteLength: 1,
					Line:       4,
					Column:     11,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "println",
					StartByte:  31,
					ByteLength: 7,
					Line:       5,
					Column:     1,
				},
				{
					TokenType:  RawStringToken,
					Value:      "Hey",
					StartByte:  39,
					ByteLength: 5,
					Line:       5,
					Column:     9,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  44,
					ByteLength: 1,
					Line:       5,
					Column:     14,
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

func TestBlockComment(t *testing.T) {
	tests := testTable{
		"discards characters in the middle of the line": {
			input: `3 + #[25 / 3]# 5`,
			want: []*Token{
				{
					TokenType:  IntToken,
					Value:      "3",
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  2,
					ByteLength: 1,
					Line:       1,
					Column:     3,
				},
				{
					TokenType:  IntToken,
					Value:      "5",
					StartByte:  15,
					ByteLength: 1,
					Line:       1,
					Column:     16,
				},
			},
		},
		"must be terminated": {
			input: `3 + #[25 / 3 5`,
			want: []*Token{
				{
					TokenType:  IntToken,
					Value:      "3",
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  2,
					ByteLength: 1,
					Line:       1,
					Column:     3,
				},
				{
					TokenType:  ErrorToken,
					Value:      "unbalanced block comments, expected 1 more block comment ending(s) `]#`",
					StartByte:  4,
					ByteLength: 10,
					Line:       1,
					Column:     5,
				},
			},
		},
		"discards multiple lines": {
			input: `
class String
	#[
		def length: Integer
			len := 0
			self.each -> len += 1
			len
		end
	]#
end
			`,
			want: []*Token{
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  ClassToken,
					Value:      "",
					StartByte:  1,
					ByteLength: 5,
					Line:       2,
					Column:     1,
				},
				{
					TokenType:  ConstantToken,
					Value:      "String",
					StartByte:  7,
					ByteLength: 6,
					Line:       2,
					Column:     7,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  13,
					ByteLength: 1,
					Line:       2,
					Column:     13,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  93,
					ByteLength: 1,
					Line:       9,
					Column:     4,
				},
				{
					TokenType:  EndToken,
					Value:      "",
					StartByte:  94,
					ByteLength: 3,
					Line:       10,
					Column:     1,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  97,
					ByteLength: 1,
					Line:       10,
					Column:     4,
				},
			},
		},
		"may be be nested": {
			input: `
class String
	#[
		def length: Integer
			len := 0
			self.each ->
				len += 1
				#[
					#[ another comment ]#
					println len
				]#
			end
			len
		end
	]#
end
			`,
			want: []*Token{
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  ClassToken,
					Value:      "",
					StartByte:  1,
					ByteLength: 5,
					Line:       2,
					Column:     1,
				},
				{
					TokenType:  ConstantToken,
					Value:      "String",
					StartByte:  7,
					ByteLength: 6,
					Line:       2,
					Column:     7,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  13,
					ByteLength: 1,
					Line:       2,
					Column:     13,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  162,
					ByteLength: 1,
					Line:       15,
					Column:     4,
				},
				{
					TokenType:  EndToken,
					Value:      "",
					StartByte:  163,
					ByteLength: 3,
					Line:       16,
					Column:     1,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  166,
					ByteLength: 1,
					Line:       16,
					Column:     4,
				},
			},
		},
		"nesting must be balanced": {
			input: `
class String
	#[
		def length: Integer
			len := 0
			self.each ->
				len += 1
				#[
					#[ another comment
					println len
			end
			len
		end
	]#
end
			`,
			want: []*Token{
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  0,
					ByteLength: 1,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  ClassToken,
					Value:      "",
					StartByte:  1,
					ByteLength: 5,
					Line:       2,
					Column:     1,
				},
				{
					TokenType:  ConstantToken,
					Value:      "String",
					StartByte:  7,
					ByteLength: 6,
					Line:       2,
					Column:     7,
				},
				{
					TokenType:  SeparatorToken,
					Value:      "",
					StartByte:  13,
					ByteLength: 1,
					Line:       2,
					Column:     13,
				},
				{
					TokenType:  ErrorToken,
					Value:      "unbalanced block comments, expected 2 more block comment ending(s) `]#`",
					StartByte:  15,
					ByteLength: 145,
					Line:       3,
					Column:     2,
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
