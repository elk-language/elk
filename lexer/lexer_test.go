package lexer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testCase struct {
	input string
	want  []*Token
}
type testTable map[string]testCase

func tokenTest(tc testCase, t *testing.T) {
	lex := New([]byte(tc.input))
	var got []*Token
	for {
		tok := lex.Next()
		if tok.IsEOF() {
			break
		}
		got = append(got, tok)
	}
	diff := cmp.Diff(tc.want, got)
	if diff != "" {
		t.Fatalf(diff)
	}
}

func TestInt(t *testing.T) {
	tests := testTable{
		"decimal": {
			input: "23",
			want: []*Token{
				{
					TokenType:  IntToken,
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
					TokenType:  IntToken,
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
					TokenType:  IntToken,
					Value:      "23_200_123",
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
					TokenType:  IntToken,
					Value:      "23_200_123",
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
					TokenType:  IntToken,
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
					TokenType:  IntToken,
					Value:      "0x354ab1",
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
					TokenType:  IntToken,
					Value:      "0x35_4a_b1",
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
					TokenType:  IntToken,
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
					TokenType:  IntToken,
					Value:      "0x",
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
					TokenType:  IntToken,
					Value:      "0X354Ab1",
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
					TokenType:  IntToken,
					Value:      "0x354f",
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
					TokenType:  IntToken,
					Value:      "0o603",
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
					TokenType:  IntToken,
					Value:      "0o3201_5200",
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
					TokenType:  IntToken,
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
					TokenType:  IntToken,
					Value:      "0o",
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
					TokenType:  IntToken,
					Value:      "0O603",
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
					TokenType:  IntToken,
					Value:      "0o603",
					StartByte:  0,
					ByteLength: 5,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IntToken,
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
					TokenType:  IntToken,
					Value:      "0q30212",
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
					TokenType:  IntToken,
					Value:      "0q3201200_23010000",
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
					TokenType:  IntToken,
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
					TokenType:  IntToken,
					Value:      "0Q30212",
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
					TokenType:  IntToken,
					Value:      "0q",
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
					TokenType:  IntToken,
					Value:      "0q30212",
					StartByte:  0,
					ByteLength: 7,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IntToken,
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
					TokenType:  IntToken,
					Value:      "0b1010",
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
					TokenType:  IntToken,
					Value:      "0b10_10_10",
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
					TokenType:  IntToken,
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
					TokenType:  IntToken,
					Value:      "0B1010",
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
					TokenType:  IntToken,
					Value:      "0b1010",
					StartByte:  0,
					ByteLength: 6,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IntToken,
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
					TokenType:  IntToken,
					Value:      "0b",
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
					TokenType:  IntToken,
					Value:      "0da12b3",
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
					TokenType:  IntToken,
					Value:      "0Da12B3",
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
					TokenType:  IntToken,
					Value:      "0d23a3b",
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
					TokenType:  IntToken,
					Value:      "0d",
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

func TestFloat(t *testing.T) {
	tests := testTable{
		"with underscores": {
			input: "245_000.254_129",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "245_000.254_129",
					StartByte:  0,
					ByteLength: 15,
					Line:       1,
					Column:     1,
				},
			},
		},
		"ends on last valid character": {
			input: "0.36f",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "0.36",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "f",
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
				},
			},
		},
		"can only be decimal": {
			input: "0x21.36",
			want: []*Token{
				{
					TokenType:  IntToken,
					Value:      "0x21",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  FloatToken,
					Value:      ".36",
					StartByte:  4,
					ByteLength: 3,
					Line:       1,
					Column:     5,
				},
			},
		},
		"with exponent": {
			input: "0.36e2",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "0.36e2",
					StartByte:  0,
					ByteLength: 6,
					Line:       1,
					Column:     1,
				},
			},
		},
		"with exponent and no dot": {
			input: "25e4",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "25e4",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
			},
		},
		"with uppercase exponent": {
			input: "0.36E2",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "0.36E2",
					StartByte:  0,
					ByteLength: 6,
					Line:       1,
					Column:     1,
				},
			},
		},
		"with negative exponent": {
			input: "25.8e-36",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      "25.8e-36",
					StartByte:  0,
					ByteLength: 8,
					Line:       1,
					Column:     1,
				},
			},
		},
		"without leading zero": {
			input: ".908267374623",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      ".908267374623",
					StartByte:  0,
					ByteLength: 13,
					Line:       1,
					Column:     1,
				},
			},
		},
		"without leading zero and with exponent": {
			input: ".8e-36",
			want: []*Token{
				{
					TokenType:  FloatToken,
					Value:      ".8e-36",
					StartByte:  0,
					ByteLength: 6,
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

func TestIdentifier(t *testing.T) {
	tests := testTable{
		"ends on the last valid character": {
			input: "foo:+",
			want: []*Token{
				{
					TokenType:  IdentifierToken,
					Value:      "foo",
					StartByte:  0,
					ByteLength: 3,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  SymbolBegToken,
					Value:      "",
					StartByte:  3,
					ByteLength: 1,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
				},
			},
		},
		"may contain letters underscores and numbers": {
			input: "some_identifier123",
			want: []*Token{
				{
					TokenType:  IdentifierToken,
					Value:      "some_identifier123",
					StartByte:  0,
					ByteLength: 18,
					Line:       1,
					Column:     1,
				},
			},
		},
		"can't start with numbers": {
			input: "3d_secure",
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
					TokenType:  IdentifierToken,
					Value:      "d_secure",
					StartByte:  1,
					ByteLength: 8,
					Line:       1,
					Column:     2,
				},
			},
		},
		"may contain utf-8 characters": {
			input: "zażółć_gęślą_jaźń + 2",
			want: []*Token{
				{
					TokenType:  IdentifierToken,
					Value:      "zażółć_gęślą_jaźń",
					StartByte:  0,
					ByteLength: 26,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  27,
					ByteLength: 1,
					Line:       1,
					Column:     19,
				},
				{
					TokenType:  IntToken,
					Value:      "2",
					StartByte:  29,
					ByteLength: 1,
					Line:       1,
					Column:     21,
				},
			},
		},
		"may start with a utf-8 character": {
			input: "łódź",
			want: []*Token{
				{
					TokenType:  IdentifierToken,
					Value:      "łódź",
					StartByte:  0,
					ByteLength: 7,
					Line:       1,
					Column:     1,
				},
			},
		},
		"may end with a question mark": {
			input: "includes?",
			want: []*Token{
				{
					TokenType:  IdentifierToken,
					Value:      "includes?",
					StartByte:  0,
					ByteLength: 9,
					Line:       1,
					Column:     1,
				},
			},
		},
		"may end with an exclamation point": {
			input: "map!",
			want: []*Token{
				{
					TokenType:  IdentifierToken,
					Value:      "map!",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
			},
		},
		"can't start with an uppercase letter": {
			input: "Dupa",
			want: []*Token{
				{
					TokenType:  ConstantToken,
					Value:      "Dupa",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
			},
		},
		"can't start with an underscore": {
			input: "_foo",
			want: []*Token{
				{
					TokenType:  PrivateIdentifierToken,
					Value:      "_foo",
					StartByte:  0,
					ByteLength: 4,
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

func TestPrivateIdentifier(t *testing.T) {
	tests := testTable{
		"ends on the last valid character": {
			input: "_foo:+",
			want: []*Token{
				{
					TokenType:  PrivateIdentifierToken,
					Value:      "_foo",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  SymbolBegToken,
					Value:      "",
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  5,
					ByteLength: 1,
					Line:       1,
					Column:     6,
				},
			},
		},
		"may contain letters underscores and numbers": {
			input: "_some_identifier123",
			want: []*Token{
				{
					TokenType:  PrivateIdentifierToken,
					Value:      "_some_identifier123",
					StartByte:  0,
					ByteLength: 19,
					Line:       1,
					Column:     1,
				},
			},
		},
		"may start with a utf-8 character": {
			input: "_łódź",
			want: []*Token{
				{
					TokenType:  PrivateIdentifierToken,
					Value:      "_łódź",
					StartByte:  0,
					ByteLength: 8,
					Line:       1,
					Column:     1,
				},
			},
		},
		"may contain utf-8 characters": {
			input: "_zażółć_gęślą_jaźń + 2",
			want: []*Token{
				{
					TokenType:  PrivateIdentifierToken,
					Value:      "_zażółć_gęślą_jaźń",
					StartByte:  0,
					ByteLength: 27,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  28,
					ByteLength: 1,
					Line:       1,
					Column:     20,
				},
				{
					TokenType:  IntToken,
					Value:      "2",
					StartByte:  30,
					ByteLength: 1,
					Line:       1,
					Column:     22,
				},
			},
		},
		"may end with a question mark": {
			input: "_includes?",
			want: []*Token{
				{
					TokenType:  PrivateIdentifierToken,
					Value:      "_includes?",
					StartByte:  0,
					ByteLength: 10,
					Line:       1,
					Column:     1,
				},
			},
		},
		"may end with an exclamation point": {
			input: "_map!",
			want: []*Token{
				{
					TokenType:  PrivateIdentifierToken,
					Value:      "_map!",
					StartByte:  0,
					ByteLength: 5,
					Line:       1,
					Column:     1,
				},
			},
		},
		"can't start with an uppercase letter": {
			input: "_Dupa",
			want: []*Token{
				{
					TokenType:  PrivateConstantToken,
					Value:      "_Dupa",
					StartByte:  0,
					ByteLength: 5,
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

func TestConstant(t *testing.T) {
	tests := testTable{
		"ends on the last valid character": {
			input: "Foo:+",
			want: []*Token{
				{
					TokenType:  ConstantToken,
					Value:      "Foo",
					StartByte:  0,
					ByteLength: 3,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  SymbolBegToken,
					Value:      "",
					StartByte:  3,
					ByteLength: 1,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
				},
			},
		},
		"may contain letters underscores and numbers": {
			input: "Some_constant123",
			want: []*Token{
				{
					TokenType:  ConstantToken,
					Value:      "Some_constant123",
					StartByte:  0,
					ByteLength: 16,
					Line:       1,
					Column:     1,
				},
			},
		},
		"can't start with numbers": {
			input: "3DSecure",
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
					TokenType:  ConstantToken,
					Value:      "DSecure",
					StartByte:  1,
					ByteLength: 7,
					Line:       1,
					Column:     2,
				},
			},
		},
		"may contain utf-8 characters": {
			input: "ZażółćGęśląJaźń + 2",
			want: []*Token{
				{
					TokenType:  ConstantToken,
					Value:      "ZażółćGęśląJaźń",
					StartByte:  0,
					ByteLength: 24,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  25,
					ByteLength: 1,
					Line:       1,
					Column:     17,
				},
				{
					TokenType:  IntToken,
					Value:      "2",
					StartByte:  27,
					ByteLength: 1,
					Line:       1,
					Column:     19,
				},
			},
		},
		"may start with a utf-8 character": {
			input: "Łódź",
			want: []*Token{
				{
					TokenType:  ConstantToken,
					Value:      "Łódź",
					StartByte:  0,
					ByteLength: 7,
					Line:       1,
					Column:     1,
				},
			},
		},
		"can't end with a question mark": {
			input: "Includes?",
			want: []*Token{
				{
					TokenType:  ConstantToken,
					Value:      "Includes",
					StartByte:  0,
					ByteLength: 8,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  QuestionMarkToken,
					Value:      "",
					StartByte:  8,
					ByteLength: 1,
					Line:       1,
					Column:     9,
				},
			},
		},
		"can't end with an exclamation point": {
			input: "Map!",
			want: []*Token{
				{
					TokenType:  ConstantToken,
					Value:      "Map",
					StartByte:  0,
					ByteLength: 3,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  BangToken,
					Value:      "",
					StartByte:  3,
					ByteLength: 1,
					Line:       1,
					Column:     4,
				},
			},
		},
		"can't start with an underscore": {
			input: "_Foo",
			want: []*Token{
				{
					TokenType:  PrivateConstantToken,
					Value:      "_Foo",
					StartByte:  0,
					ByteLength: 4,
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

func TestPrivateConstant(t *testing.T) {
	tests := testTable{
		"ends on the last valid character": {
			input: "_Foo:+",
			want: []*Token{
				{
					TokenType:  PrivateConstantToken,
					Value:      "_Foo",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  SymbolBegToken,
					Value:      "",
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  5,
					ByteLength: 1,
					Line:       1,
					Column:     6,
				},
			},
		},
		"may contain letters underscores and numbers": {
			input: "_Some_identifier123",
			want: []*Token{
				{
					TokenType:  PrivateConstantToken,
					Value:      "_Some_identifier123",
					StartByte:  0,
					ByteLength: 19,
					Line:       1,
					Column:     1,
				},
			},
		},
		"may start with a utf-8 character": {
			input: "_Łódź",
			want: []*Token{
				{
					TokenType:  PrivateConstantToken,
					Value:      "_Łódź",
					StartByte:  0,
					ByteLength: 8,
					Line:       1,
					Column:     1,
				},
			},
		},
		"may contain utf-8 characters": {
			input: "_Zażółć_gęślą_jaźń + 2",
			want: []*Token{
				{
					TokenType:  PrivateConstantToken,
					Value:      "_Zażółć_gęślą_jaźń",
					StartByte:  0,
					ByteLength: 27,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  28,
					ByteLength: 1,
					Line:       1,
					Column:     20,
				},
				{
					TokenType:  IntToken,
					Value:      "2",
					StartByte:  30,
					ByteLength: 1,
					Line:       1,
					Column:     22,
				},
			},
		},
		"can't end with a question mark": {
			input: "_Includes?",
			want: []*Token{
				{
					TokenType:  PrivateConstantToken,
					Value:      "_Includes",
					StartByte:  0,
					ByteLength: 9,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  QuestionMarkToken,
					Value:      "",
					StartByte:  9,
					ByteLength: 1,
					Line:       1,
					Column:     10,
				},
			},
		},
		"can't end with an exclamation point": {
			input: "_Map!",
			want: []*Token{
				{
					TokenType:  PrivateConstantToken,
					Value:      "_Map",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  BangToken,
					Value:      "",
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
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

func TestInstanceVariable(t *testing.T) {
	tests := testTable{
		"ends on the last valid character": {
			input: "@foo:+",
			want: []*Token{
				{
					TokenType:  InstanceVariableToken,
					Value:      "foo",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  SymbolBegToken,
					Value:      "",
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  5,
					ByteLength: 1,
					Line:       1,
					Column:     6,
				},
			},
		},
		"may contain letters underscores and numbers": {
			input: "@some_ivar123",
			want: []*Token{
				{
					TokenType:  InstanceVariableToken,
					Value:      "some_ivar123",
					StartByte:  0,
					ByteLength: 13,
					Line:       1,
					Column:     1,
				},
			},
		},
		"may start with an uppercase letter": {
			input: "@SomeIvar123",
			want: []*Token{
				{
					TokenType:  InstanceVariableToken,
					Value:      "SomeIvar123",
					StartByte:  0,
					ByteLength: 12,
					Line:       1,
					Column:     1,
				},
			},
		},
		"may start with an underscore": {
			input: "@_bar",
			want: []*Token{
				{
					TokenType:  InstanceVariableToken,
					Value:      "_bar",
					StartByte:  0,
					ByteLength: 5,
					Line:       1,
					Column:     1,
				},
			},
		},
		"may start with a utf-8 character": {
			input: "@łódź",
			want: []*Token{
				{
					TokenType:  InstanceVariableToken,
					Value:      "łódź",
					StartByte:  0,
					ByteLength: 8,
					Line:       1,
					Column:     1,
				},
			},
		},
		"may contain utf-8 characters": {
			input: "@zażółć_gęślą_jaźń + 2",
			want: []*Token{
				{
					TokenType:  InstanceVariableToken,
					Value:      "zażółć_gęślą_jaźń",
					StartByte:  0,
					ByteLength: 27,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  PlusToken,
					Value:      "",
					StartByte:  28,
					ByteLength: 1,
					Line:       1,
					Column:     20,
				},
				{
					TokenType:  IntToken,
					Value:      "2",
					StartByte:  30,
					ByteLength: 1,
					Line:       1,
					Column:     22,
				},
			},
		},
		"can't end with a question mark": {
			input: "@includes?",
			want: []*Token{
				{
					TokenType:  InstanceVariableToken,
					Value:      "includes",
					StartByte:  0,
					ByteLength: 9,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  QuestionMarkToken,
					Value:      "",
					StartByte:  9,
					ByteLength: 1,
					Line:       1,
					Column:     10,
				},
			},
		},
		"can't end with an exclamation point": {
			input: "@map!",
			want: []*Token{
				{
					TokenType:  InstanceVariableToken,
					Value:      "map",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  BangToken,
					Value:      "",
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
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
