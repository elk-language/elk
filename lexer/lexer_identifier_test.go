package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestIdentifier(t *testing.T) {
	tests := testTable{
		"ends on the last valid character": {
			input: "foo:+",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(2, 1, 3)), token.PUBLIC_IDENTIFIER, "foo"),
				T(S(P(3, 1, 4), P(3, 1, 4)), token.COLON),
				T(S(P(4, 1, 5), P(4, 1, 5)), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "some_identifier123",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(17, 1, 18)), token.PUBLIC_IDENTIFIER, "some_identifier123"),
			},
		},
		"can't start with numbers": {
			input: "3d_secure",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(0, 1, 1)), token.INT, "3"),
				V(S(P(1, 1, 2), P(8, 1, 9)), token.PUBLIC_IDENTIFIER, "d_secure"),
			},
		},
		"may contain utf-8 characters": {
			input: "zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(24, 1, 17)), token.PUBLIC_IDENTIFIER, "zażółć_gęślą_jaźń"),
				T(S(P(27, 1, 19), P(27, 1, 19)), token.PLUS),
				V(S(P(29, 1, 21), P(29, 1, 21)), token.INT, "2"),
			},
		},
		"may start with a utf-8 character": {
			input: "łódź",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(5, 1, 4)), token.PUBLIC_IDENTIFIER, "łódź"),
			},
		},
		"can't start with an uppercase letter": {
			input: "Dupa",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(3, 1, 4)), token.PUBLIC_CONSTANT, "Dupa"),
			},
		},
		"can't start with an underscore": {
			input: "_foo",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(3, 1, 4)), token.PRIVATE_IDENTIFIER, "_foo"),
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
			want: []*token.Token{
				V(S(P(0, 1, 1), P(3, 1, 4)), token.PRIVATE_IDENTIFIER, "_foo"),
				T(S(P(4, 1, 5), P(4, 1, 5)), token.COLON),
				T(S(P(5, 1, 6), P(5, 1, 6)), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "_some_identifier123",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(18, 1, 19)), token.PRIVATE_IDENTIFIER, "_some_identifier123"),
			},
		},
		"may start with a utf-8 character": {
			input: "_łódź",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(6, 1, 5)), token.PRIVATE_IDENTIFIER, "_łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "_zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(25, 1, 18)), token.PRIVATE_IDENTIFIER, "_zażółć_gęślą_jaźń"),
				T(S(P(28, 1, 20), P(28, 1, 20)), token.PLUS),
				V(S(P(30, 1, 22), P(30, 1, 22)), token.INT, "2"),
			},
		},
		"can't start with an uppercase letter": {
			input: "_Dupa",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(4, 1, 5)), token.PRIVATE_CONSTANT, "_Dupa"),
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
			want: []*token.Token{
				V(S(P(0, 1, 1), P(2, 1, 3)), token.PUBLIC_CONSTANT, "Foo"),
				T(S(P(3, 1, 4), P(3, 1, 4)), token.COLON),
				T(S(P(4, 1, 5), P(4, 1, 5)), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "Some_constant123",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(15, 1, 16)), token.PUBLIC_CONSTANT, "Some_constant123"),
			},
		},
		"can't start with numbers": {
			input: "3DSecure",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(0, 1, 1)), token.INT, "3"),
				V(S(P(1, 1, 2), P(7, 1, 8)), token.PUBLIC_CONSTANT, "DSecure"),
			},
		},
		"may contain utf-8 characters": {
			input: "ZażółćGęśląJaźń + 2",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(22, 1, 15)), token.PUBLIC_CONSTANT, "ZażółćGęśląJaźń"),
				T(S(P(25, 1, 17), P(25, 1, 17)), token.PLUS),
				V(S(P(27, 1, 19), P(27, 1, 19)), token.INT, "2"),
			},
		},
		"may start with a utf-8 character": {
			input: "Łódź",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(5, 1, 4)), token.PUBLIC_CONSTANT, "Łódź"),
			},
		},
		"can't end with a question mark": {
			input: "Includes?",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(7, 1, 8)), token.PUBLIC_CONSTANT, "Includes"),
				T(S(P(8, 1, 9), P(8, 1, 9)), token.QUESTION),
			},
		},
		"can't end with an exclamation point": {
			input: "Map!",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(2, 1, 3)), token.PUBLIC_CONSTANT, "Map"),
				T(S(P(3, 1, 4), P(3, 1, 4)), token.BANG),
			},
		},
		"can't start with an underscore": {
			input: "_Foo",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(3, 1, 4)), token.PRIVATE_CONSTANT, "_Foo"),
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
			want: []*token.Token{
				V(S(P(0, 1, 1), P(3, 1, 4)), token.PRIVATE_CONSTANT, "_Foo"),
				T(S(P(4, 1, 5), P(4, 1, 5)), token.COLON),
				T(S(P(5, 1, 6), P(5, 1, 6)), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "_Some_identifier123",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(18, 1, 19)), token.PRIVATE_CONSTANT, "_Some_identifier123"),
			},
		},
		"may start with a utf-8 character": {
			input: "_Łódź",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(6, 1, 5)), token.PRIVATE_CONSTANT, "_Łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "_Zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(25, 1, 18)), token.PRIVATE_CONSTANT, "_Zażółć_gęślą_jaźń"),
				T(S(P(28, 1, 20), P(28, 1, 20)), token.PLUS),
				V(S(P(30, 1, 22), P(30, 1, 22)), token.INT, "2"),
			},
		},
		"can't end with a question mark": {
			input: "_Includes?",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(8, 1, 9)), token.PRIVATE_CONSTANT, "_Includes"),
				T(S(P(9, 1, 10), P(9, 1, 10)), token.QUESTION),
			},
		},
		"can't end with an exclamation point": {
			input: "_Map!",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(3, 1, 4)), token.PRIVATE_CONSTANT, "_Map"),
				T(S(P(4, 1, 5), P(4, 1, 5)), token.BANG),
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
			want: []*token.Token{
				V(S(P(0, 1, 1), P(3, 1, 4)), token.INSTANCE_VARIABLE, "foo"),
				T(S(P(4, 1, 5), P(4, 1, 5)), token.COLON),
				T(S(P(5, 1, 6), P(5, 1, 6)), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "@some_ivar123",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(12, 1, 13)), token.INSTANCE_VARIABLE, "some_ivar123"),
			},
		},
		"may start with an uppercase letter": {
			input: "@SomeIvar123",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(11, 1, 12)), token.INSTANCE_VARIABLE, "SomeIvar123"),
			},
		},
		"may start with an underscore": {
			input: "@_bar",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(4, 1, 5)), token.INSTANCE_VARIABLE, "_bar"),
			},
		},
		"may start with a utf-8 character": {
			input: "@łódź",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(6, 1, 5)), token.INSTANCE_VARIABLE, "łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "@zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(S(P(0, 1, 1), P(25, 1, 18)), token.INSTANCE_VARIABLE, "zażółć_gęślą_jaźń"),
				T(S(P(28, 1, 20), P(28, 1, 20)), token.PLUS),
				V(S(P(30, 1, 22), P(30, 1, 22)), token.INT, "2"),
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
			want: []*token.Token{
				T(S(P(0, 1, 1), P(4, 1, 5)), token.FALSE),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
