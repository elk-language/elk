package lexer

import "testing"

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
					StartByte:  3,
					ByteLength: 1,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  PlusToken,
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
					TokenType:  DecIntToken,
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
					StartByte:  27,
					ByteLength: 1,
					Line:       1,
					Column:     19,
				},
				{
					TokenType:  DecIntToken,
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
		"can't have a question mark in the middle": {
			input: "foo?bar",
			want: []*Token{
				{
					TokenType:  IdentifierToken,
					Value:      "foo?",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "bar",
					StartByte:  4,
					ByteLength: 3,
					Line:       1,
					Column:     5,
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
		"can't have an exclamation point in the middle": {
			input: "foo!bar",
			want: []*Token{
				{
					TokenType:  IdentifierToken,
					Value:      "foo!",
					StartByte:  0,
					ByteLength: 4,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "bar",
					StartByte:  4,
					ByteLength: 3,
					Line:       1,
					Column:     5,
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
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
				},
				{
					TokenType:  PlusToken,
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
					StartByte:  28,
					ByteLength: 1,
					Line:       1,
					Column:     20,
				},
				{
					TokenType:  DecIntToken,
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
		"can't have a question mark in the middle": {
			input: "_foo?bar",
			want: []*Token{
				{
					TokenType:  PrivateIdentifierToken,
					Value:      "_foo?",
					StartByte:  0,
					ByteLength: 5,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "bar",
					StartByte:  5,
					ByteLength: 3,
					Line:       1,
					Column:     6,
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
		"can't have an exclamation point in the middle": {
			input: "_foo!bar",
			want: []*Token{
				{
					TokenType:  PrivateIdentifierToken,
					Value:      "_foo!",
					StartByte:  0,
					ByteLength: 5,
					Line:       1,
					Column:     1,
				},
				{
					TokenType:  IdentifierToken,
					Value:      "bar",
					StartByte:  5,
					ByteLength: 3,
					Line:       1,
					Column:     6,
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
					StartByte:  3,
					ByteLength: 1,
					Line:       1,
					Column:     4,
				},
				{
					TokenType:  PlusToken,
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
					TokenType:  DecIntToken,
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
					StartByte:  25,
					ByteLength: 1,
					Line:       1,
					Column:     17,
				},
				{
					TokenType:  DecIntToken,
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
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
				},
				{
					TokenType:  PlusToken,
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
					StartByte:  28,
					ByteLength: 1,
					Line:       1,
					Column:     20,
				},
				{
					TokenType:  DecIntToken,
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
					StartByte:  4,
					ByteLength: 1,
					Line:       1,
					Column:     5,
				},
				{
					TokenType:  PlusToken,
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
					StartByte:  28,
					ByteLength: 1,
					Line:       1,
					Column:     20,
				},
				{
					TokenType:  DecIntToken,
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
