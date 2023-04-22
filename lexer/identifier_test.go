package lexer

import "testing"

func TestIdentifier(t *testing.T) {
	tests := testTable{
		"ends on the last valid character": {
			input: "foo:+",
			want: []*Token{
				V(IdentifierToken, "foo", 0, 3, 1, 1),
				T(SymbolBegToken, 3, 1, 1, 4),
				T(PlusToken, 4, 1, 1, 5),
			},
		},
		"may contain letters underscores and numbers": {
			input: "some_identifier123",
			want: []*Token{
				V(IdentifierToken, "some_identifier123", 0, 18, 1, 1),
			},
		},
		"can't start with numbers": {
			input: "3d_secure",
			want: []*Token{
				V(DecIntToken, "3", 0, 1, 1, 1),
				V(IdentifierToken, "d_secure", 1, 8, 1, 2),
			},
		},
		"may contain utf-8 characters": {
			input: "zażółć_gęślą_jaźń + 2",
			want: []*Token{
				V(IdentifierToken, "zażółć_gęślą_jaźń", 0, 26, 1, 1),
				T(PlusToken, 27, 1, 1, 19),
				V(DecIntToken, "2", 29, 1, 1, 21),
			},
		},
		"may start with a utf-8 character": {
			input: "łódź",
			want: []*Token{
				V(IdentifierToken, "łódź", 0, 7, 1, 1),
			},
		},
		"can't have a question mark in the middle": {
			input: "foo?bar",
			want: []*Token{
				V(IdentifierToken, "foo?", 0, 4, 1, 1),
				V(IdentifierToken, "bar", 4, 3, 1, 5),
			},
		},
		"may end with a question mark": {
			input: "includes?",
			want: []*Token{
				V(IdentifierToken, "includes?", 0, 9, 1, 1),
			},
		},
		"can't have an exclamation point in the middle": {
			input: "foo!bar",
			want: []*Token{
				V(IdentifierToken, "foo!", 0, 4, 1, 1),
				V(IdentifierToken, "bar", 4, 3, 1, 5),
			},
		},
		"may end with an exclamation point": {
			input: "map!",
			want: []*Token{
				V(IdentifierToken, "map!", 0, 4, 1, 1),
			},
		},
		"can't start with an uppercase letter": {
			input: "Dupa",
			want: []*Token{
				V(ConstantToken, "Dupa", 0, 4, 1, 1),
			},
		},
		"can't start with an underscore": {
			input: "_foo",
			want: []*Token{
				V(PrivateIdentifierToken, "_foo", 0, 4, 1, 1),
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
				V(PrivateIdentifierToken, "_foo", 0, 4, 1, 1),
				T(SymbolBegToken, 4, 1, 1, 5),
				T(PlusToken, 5, 1, 1, 6),
			},
		},
		"may contain letters underscores and numbers": {
			input: "_some_identifier123",
			want: []*Token{
				V(PrivateIdentifierToken, "_some_identifier123", 0, 19, 1, 1),
			},
		},
		"may start with a utf-8 character": {
			input: "_łódź",
			want: []*Token{
				V(PrivateIdentifierToken, "_łódź", 0, 8, 1, 1),
			},
		},
		"may contain utf-8 characters": {
			input: "_zażółć_gęślą_jaźń + 2",
			want: []*Token{
				V(PrivateIdentifierToken, "_zażółć_gęślą_jaźń", 0, 27, 1, 1),
				T(PlusToken, 28, 1, 1, 20),
				V(DecIntToken, "2", 30, 1, 1, 22),
			},
		},
		"may end with a question mark": {
			input: "_includes?",
			want: []*Token{
				V(PrivateIdentifierToken, "_includes?", 0, 10, 1, 1),
			},
		},
		"can't have a question mark in the middle": {
			input: "_foo?bar",
			want: []*Token{
				V(PrivateIdentifierToken, "_foo?", 0, 5, 1, 1),
				V(IdentifierToken, "bar", 5, 3, 1, 6),
			},
		},
		"may end with an exclamation point": {
			input: "_map!",
			want: []*Token{
				V(PrivateIdentifierToken, "_map!", 0, 5, 1, 1),
			},
		},
		"can't have an exclamation point in the middle": {
			input: "_foo!bar",
			want: []*Token{
				V(PrivateIdentifierToken, "_foo!", 0, 5, 1, 1),
				V(IdentifierToken, "bar", 5, 3, 1, 6),
			},
		},
		"can't start with an uppercase letter": {
			input: "_Dupa",
			want: []*Token{
				V(PrivateConstantToken, "_Dupa", 0, 5, 1, 1),
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
				V(ConstantToken, "Foo", 0, 3, 1, 1),
				T(SymbolBegToken, 3, 1, 1, 4),
				T(PlusToken, 4, 1, 1, 5),
			},
		},
		"may contain letters underscores and numbers": {
			input: "Some_constant123",
			want: []*Token{
				V(ConstantToken, "Some_constant123", 0, 16, 1, 1),
			},
		},
		"can't start with numbers": {
			input: "3DSecure",
			want: []*Token{
				V(DecIntToken, "3", 0, 1, 1, 1),
				V(ConstantToken, "DSecure", 1, 7, 1, 2),
			},
		},
		"may contain utf-8 characters": {
			input: "ZażółćGęśląJaźń + 2",
			want: []*Token{
				V(ConstantToken, "ZażółćGęśląJaźń", 0, 24, 1, 1),
				T(PlusToken, 25, 1, 1, 17),
				V(DecIntToken, "2", 27, 1, 1, 19),
			},
		},
		"may start with a utf-8 character": {
			input: "Łódź",
			want: []*Token{
				V(ConstantToken, "Łódź", 0, 7, 1, 1),
			},
		},
		"can't end with a question mark": {
			input: "Includes?",
			want: []*Token{
				V(ConstantToken, "Includes", 0, 8, 1, 1),
				T(QuestionMarkToken, 8, 1, 1, 9),
			},
		},
		"can't end with an exclamation point": {
			input: "Map!",
			want: []*Token{
				V(ConstantToken, "Map", 0, 3, 1, 1),
				T(BangToken, 3, 1, 1, 4),
			},
		},
		"can't start with an underscore": {
			input: "_Foo",
			want: []*Token{
				V(PrivateConstantToken, "_Foo", 0, 4, 1, 1),
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
				V(PrivateConstantToken, "_Foo", 0, 4, 1, 1),
				T(SymbolBegToken, 4, 1, 1, 5),
				T(PlusToken, 5, 1, 1, 6),
			},
		},
		"may contain letters underscores and numbers": {
			input: "_Some_identifier123",
			want: []*Token{
				V(PrivateConstantToken, "_Some_identifier123", 0, 19, 1, 1),
			},
		},
		"may start with a utf-8 character": {
			input: "_Łódź",
			want: []*Token{
				V(PrivateConstantToken, "_Łódź", 0, 8, 1, 1),
			},
		},
		"may contain utf-8 characters": {
			input: "_Zażółć_gęślą_jaźń + 2",
			want: []*Token{
				V(PrivateConstantToken, "_Zażółć_gęślą_jaźń", 0, 27, 1, 1),
				T(PlusToken, 28, 1, 1, 20),
				V(DecIntToken, "2", 30, 1, 1, 22),
			},
		},
		"can't end with a question mark": {
			input: "_Includes?",
			want: []*Token{
				V(PrivateConstantToken, "_Includes", 0, 9, 1, 1),
				T(QuestionMarkToken, 9, 1, 1, 10),
			},
		},
		"can't end with an exclamation point": {
			input: "_Map!",
			want: []*Token{
				V(PrivateConstantToken, "_Map", 0, 4, 1, 1),
				T(BangToken, 4, 1, 1, 5),
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
				V(InstanceVariableToken, "foo", 0, 4, 1, 1),
				T(SymbolBegToken, 4, 1, 1, 5),
				T(PlusToken, 5, 1, 1, 6),
			},
		},
		"may contain letters underscores and numbers": {
			input: "@some_ivar123",
			want: []*Token{
				V(InstanceVariableToken, "some_ivar123", 0, 13, 1, 1),
			},
		},
		"may start with an uppercase letter": {
			input: "@SomeIvar123",
			want: []*Token{
				V(InstanceVariableToken, "SomeIvar123", 0, 12, 1, 1),
			},
		},
		"may start with an underscore": {
			input: "@_bar",
			want: []*Token{
				V(InstanceVariableToken, "_bar", 0, 5, 1, 1),
			},
		},
		"may start with a utf-8 character": {
			input: "@łódź",
			want: []*Token{
				V(InstanceVariableToken, "łódź", 0, 8, 1, 1),
			},
		},
		"may contain utf-8 characters": {
			input: "@zażółć_gęślą_jaźń + 2",
			want: []*Token{
				V(InstanceVariableToken, "zażółć_gęślą_jaźń", 0, 27, 1, 1),
				T(PlusToken, 28, 1, 1, 20),
				V(DecIntToken, "2", 30, 1, 1, 22),
			},
		},
		"can't end with a question mark": {
			input: "@includes?",
			want: []*Token{
				V(InstanceVariableToken, "includes", 0, 9, 1, 1),
				T(QuestionMarkToken, 9, 1, 1, 10),
			},
		},
		"can't end with an exclamation point": {
			input: "@map!",
			want: []*Token{
				V(InstanceVariableToken, "map", 0, 4, 1, 1),
				T(BangToken, 4, 1, 1, 5),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}