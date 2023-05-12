package lexer

import "testing"

func TestIdentifier(t *testing.T) {
	tests := testTable{
		"ends on the last valid character": {
			input: "foo:+",
			want: []*Token{
				V(P(0, 3, 1, 1), PublicIdentifierToken, "foo"),
				T(P(3, 1, 1, 4), SymbolBegToken),
				T(P(4, 1, 1, 5), PlusToken),
			},
		},
		"may contain letters underscores and numbers": {
			input: "some_identifier123",
			want: []*Token{
				V(P(0, 18, 1, 1), PublicIdentifierToken, "some_identifier123"),
			},
		},
		"can't start with numbers": {
			input: "3d_secure",
			want: []*Token{
				V(P(0, 1, 1, 1), DecIntToken, "3"),
				V(P(1, 8, 1, 2), PublicIdentifierToken, "d_secure"),
			},
		},
		"may contain utf-8 characters": {
			input: "zażółć_gęślą_jaźń + 2",
			want: []*Token{
				V(P(0, 26, 1, 1), PublicIdentifierToken, "zażółć_gęślą_jaźń"),
				T(P(27, 1, 1, 19), PlusToken),
				V(P(29, 1, 1, 21), DecIntToken, "2"),
			},
		},
		"may start with a utf-8 character": {
			input: "łódź",
			want: []*Token{
				V(P(0, 7, 1, 1), PublicIdentifierToken, "łódź"),
			},
		},
		"can't start with an uppercase letter": {
			input: "Dupa",
			want: []*Token{
				V(P(0, 4, 1, 1), PublicConstantToken, "Dupa"),
			},
		},
		"can't start with an underscore": {
			input: "_foo",
			want: []*Token{
				V(P(0, 4, 1, 1), PrivateIdentifierToken, "_foo"),
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
				V(P(0, 4, 1, 1), PrivateIdentifierToken, "_foo"),
				T(P(4, 1, 1, 5), SymbolBegToken),
				T(P(5, 1, 1, 6), PlusToken),
			},
		},
		"may contain letters underscores and numbers": {
			input: "_some_identifier123",
			want: []*Token{
				V(P(0, 19, 1, 1), PrivateIdentifierToken, "_some_identifier123"),
			},
		},
		"may start with a utf-8 character": {
			input: "_łódź",
			want: []*Token{
				V(P(0, 8, 1, 1), PrivateIdentifierToken, "_łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "_zażółć_gęślą_jaźń + 2",
			want: []*Token{
				V(P(0, 27, 1, 1), PrivateIdentifierToken, "_zażółć_gęślą_jaźń"),
				T(P(28, 1, 1, 20), PlusToken),
				V(P(30, 1, 1, 22), DecIntToken, "2"),
			},
		},
		"can't start with an uppercase letter": {
			input: "_Dupa",
			want: []*Token{
				V(P(0, 5, 1, 1), PrivateConstantToken, "_Dupa"),
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
				V(P(0, 3, 1, 1), PublicConstantToken, "Foo"),
				T(P(3, 1, 1, 4), SymbolBegToken),
				T(P(4, 1, 1, 5), PlusToken),
			},
		},
		"may contain letters underscores and numbers": {
			input: "Some_constant123",
			want: []*Token{
				V(P(0, 16, 1, 1), PublicConstantToken, "Some_constant123"),
			},
		},
		"can't start with numbers": {
			input: "3DSecure",
			want: []*Token{
				V(P(0, 1, 1, 1), DecIntToken, "3"),
				V(P(1, 7, 1, 2), PublicConstantToken, "DSecure"),
			},
		},
		"may contain utf-8 characters": {
			input: "ZażółćGęśląJaźń + 2",
			want: []*Token{
				V(P(0, 24, 1, 1), PublicConstantToken, "ZażółćGęśląJaźń"),
				T(P(25, 1, 1, 17), PlusToken),
				V(P(27, 1, 1, 19), DecIntToken, "2"),
			},
		},
		"may start with a utf-8 character": {
			input: "Łódź",
			want: []*Token{
				V(P(0, 7, 1, 1), PublicConstantToken, "Łódź"),
			},
		},
		"can't end with a question mark": {
			input: "Includes?",
			want: []*Token{
				V(P(0, 8, 1, 1), PublicConstantToken, "Includes"),
				T(P(8, 1, 1, 9), QuestionMarkToken),
			},
		},
		"can't end with an exclamation point": {
			input: "Map!",
			want: []*Token{
				V(P(0, 3, 1, 1), PublicConstantToken, "Map"),
				T(P(3, 1, 1, 4), BangToken),
			},
		},
		"can't start with an underscore": {
			input: "_Foo",
			want: []*Token{
				V(P(0, 4, 1, 1), PrivateConstantToken, "_Foo"),
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
				V(P(0, 4, 1, 1), PrivateConstantToken, "_Foo"),
				T(P(4, 1, 1, 5), SymbolBegToken),
				T(P(5, 1, 1, 6), PlusToken),
			},
		},
		"may contain letters underscores and numbers": {
			input: "_Some_identifier123",
			want: []*Token{
				V(P(0, 19, 1, 1), PrivateConstantToken, "_Some_identifier123"),
			},
		},
		"may start with a utf-8 character": {
			input: "_Łódź",
			want: []*Token{
				V(P(0, 8, 1, 1), PrivateConstantToken, "_Łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "_Zażółć_gęślą_jaźń + 2",
			want: []*Token{
				V(P(0, 27, 1, 1), PrivateConstantToken, "_Zażółć_gęślą_jaźń"),
				T(P(28, 1, 1, 20), PlusToken),
				V(P(30, 1, 1, 22), DecIntToken, "2"),
			},
		},
		"can't end with a question mark": {
			input: "_Includes?",
			want: []*Token{
				V(P(0, 9, 1, 1), PrivateConstantToken, "_Includes"),
				T(P(9, 1, 1, 10), QuestionMarkToken),
			},
		},
		"can't end with an exclamation point": {
			input: "_Map!",
			want: []*Token{
				V(P(0, 4, 1, 1), PrivateConstantToken, "_Map"),
				T(P(4, 1, 1, 5), BangToken),
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
				V(P(0, 4, 1, 1), InstanceVariableToken, "foo"),
				T(P(4, 1, 1, 5), SymbolBegToken),
				T(P(5, 1, 1, 6), PlusToken),
			},
		},
		"may contain letters underscores and numbers": {
			input: "@some_ivar123",
			want: []*Token{
				V(P(0, 13, 1, 1), InstanceVariableToken, "some_ivar123"),
			},
		},
		"may start with an uppercase letter": {
			input: "@SomeIvar123",
			want: []*Token{
				V(P(0, 12, 1, 1), InstanceVariableToken, "SomeIvar123"),
			},
		},
		"may start with an underscore": {
			input: "@_bar",
			want: []*Token{
				V(P(0, 5, 1, 1), InstanceVariableToken, "_bar"),
			},
		},
		"may start with a utf-8 character": {
			input: "@łódź",
			want: []*Token{
				V(P(0, 8, 1, 1), InstanceVariableToken, "łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "@zażółć_gęślą_jaźń + 2",
			want: []*Token{
				V(P(0, 27, 1, 1), InstanceVariableToken, "zażółć_gęślą_jaźń"),
				T(P(28, 1, 1, 20), PlusToken),
				V(P(30, 1, 1, 22), DecIntToken, "2"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestKeyword(t *testing.T) {
	tests := testTable{
		"has correct position": {
			input: "false",
			want: []*Token{
				T(P(0, 5, 1, 1), FalseToken),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
