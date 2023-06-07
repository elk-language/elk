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
				V(P(0, 3, 1, 1), token.PUBLIC_IDENTIFIER, "foo"),
				T(P(3, 1, 1, 4), token.COLON),
				T(P(4, 1, 1, 5), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "some_identifier123",
			want: []*token.Token{
				V(P(0, 18, 1, 1), token.PUBLIC_IDENTIFIER, "some_identifier123"),
			},
		},
		"can't start with numbers": {
			input: "3d_secure",
			want: []*token.Token{
				V(P(0, 1, 1, 1), token.DEC_INT, "3"),
				V(P(1, 8, 1, 2), token.PUBLIC_IDENTIFIER, "d_secure"),
			},
		},
		"may contain utf-8 characters": {
			input: "zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(P(0, 26, 1, 1), token.PUBLIC_IDENTIFIER, "zażółć_gęślą_jaźń"),
				T(P(27, 1, 1, 19), token.PLUS),
				V(P(29, 1, 1, 21), token.DEC_INT, "2"),
			},
		},
		"may start with a utf-8 character": {
			input: "łódź",
			want: []*token.Token{
				V(P(0, 7, 1, 1), token.PUBLIC_IDENTIFIER, "łódź"),
			},
		},
		"can't start with an uppercase letter": {
			input: "Dupa",
			want: []*token.Token{
				V(P(0, 4, 1, 1), token.PUBLIC_CONSTANT, "Dupa"),
			},
		},
		"can't start with an underscore": {
			input: "_foo",
			want: []*token.Token{
				V(P(0, 4, 1, 1), token.PRIVATE_IDENTIFIER, "_foo"),
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
				V(P(0, 4, 1, 1), token.PRIVATE_IDENTIFIER, "_foo"),
				T(P(4, 1, 1, 5), token.COLON),
				T(P(5, 1, 1, 6), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "_some_identifier123",
			want: []*token.Token{
				V(P(0, 19, 1, 1), token.PRIVATE_IDENTIFIER, "_some_identifier123"),
			},
		},
		"may start with a utf-8 character": {
			input: "_łódź",
			want: []*token.Token{
				V(P(0, 8, 1, 1), token.PRIVATE_IDENTIFIER, "_łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "_zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(P(0, 27, 1, 1), token.PRIVATE_IDENTIFIER, "_zażółć_gęślą_jaźń"),
				T(P(28, 1, 1, 20), token.PLUS),
				V(P(30, 1, 1, 22), token.DEC_INT, "2"),
			},
		},
		"can't start with an uppercase letter": {
			input: "_Dupa",
			want: []*token.Token{
				V(P(0, 5, 1, 1), token.PRIVATE_CONSTANT, "_Dupa"),
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
				V(P(0, 3, 1, 1), token.PUBLIC_CONSTANT, "Foo"),
				T(P(3, 1, 1, 4), token.COLON),
				T(P(4, 1, 1, 5), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "Some_constant123",
			want: []*token.Token{
				V(P(0, 16, 1, 1), token.PUBLIC_CONSTANT, "Some_constant123"),
			},
		},
		"can't start with numbers": {
			input: "3DSecure",
			want: []*token.Token{
				V(P(0, 1, 1, 1), token.DEC_INT, "3"),
				V(P(1, 7, 1, 2), token.PUBLIC_CONSTANT, "DSecure"),
			},
		},
		"may contain utf-8 characters": {
			input: "ZażółćGęśląJaźń + 2",
			want: []*token.Token{
				V(P(0, 24, 1, 1), token.PUBLIC_CONSTANT, "ZażółćGęśląJaźń"),
				T(P(25, 1, 1, 17), token.PLUS),
				V(P(27, 1, 1, 19), token.DEC_INT, "2"),
			},
		},
		"may start with a utf-8 character": {
			input: "Łódź",
			want: []*token.Token{
				V(P(0, 7, 1, 1), token.PUBLIC_CONSTANT, "Łódź"),
			},
		},
		"can't end with a question mark": {
			input: "Includes?",
			want: []*token.Token{
				V(P(0, 8, 1, 1), token.PUBLIC_CONSTANT, "Includes"),
				T(P(8, 1, 1, 9), token.QUESTION),
			},
		},
		"can't end with an exclamation point": {
			input: "Map!",
			want: []*token.Token{
				V(P(0, 3, 1, 1), token.PUBLIC_CONSTANT, "Map"),
				T(P(3, 1, 1, 4), token.BANG),
			},
		},
		"can't start with an underscore": {
			input: "_Foo",
			want: []*token.Token{
				V(P(0, 4, 1, 1), token.PRIVATE_CONSTANT, "_Foo"),
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
				V(P(0, 4, 1, 1), token.PRIVATE_CONSTANT, "_Foo"),
				T(P(4, 1, 1, 5), token.COLON),
				T(P(5, 1, 1, 6), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "_Some_identifier123",
			want: []*token.Token{
				V(P(0, 19, 1, 1), token.PRIVATE_CONSTANT, "_Some_identifier123"),
			},
		},
		"may start with a utf-8 character": {
			input: "_Łódź",
			want: []*token.Token{
				V(P(0, 8, 1, 1), token.PRIVATE_CONSTANT, "_Łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "_Zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(P(0, 27, 1, 1), token.PRIVATE_CONSTANT, "_Zażółć_gęślą_jaźń"),
				T(P(28, 1, 1, 20), token.PLUS),
				V(P(30, 1, 1, 22), token.DEC_INT, "2"),
			},
		},
		"can't end with a question mark": {
			input: "_Includes?",
			want: []*token.Token{
				V(P(0, 9, 1, 1), token.PRIVATE_CONSTANT, "_Includes"),
				T(P(9, 1, 1, 10), token.QUESTION),
			},
		},
		"can't end with an exclamation point": {
			input: "_Map!",
			want: []*token.Token{
				V(P(0, 4, 1, 1), token.PRIVATE_CONSTANT, "_Map"),
				T(P(4, 1, 1, 5), token.BANG),
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
				V(P(0, 4, 1, 1), token.INSTANCE_VARIABLE, "foo"),
				T(P(4, 1, 1, 5), token.COLON),
				T(P(5, 1, 1, 6), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "@some_ivar123",
			want: []*token.Token{
				V(P(0, 13, 1, 1), token.INSTANCE_VARIABLE, "some_ivar123"),
			},
		},
		"may start with an uppercase letter": {
			input: "@SomeIvar123",
			want: []*token.Token{
				V(P(0, 12, 1, 1), token.INSTANCE_VARIABLE, "SomeIvar123"),
			},
		},
		"may start with an underscore": {
			input: "@_bar",
			want: []*token.Token{
				V(P(0, 5, 1, 1), token.INSTANCE_VARIABLE, "_bar"),
			},
		},
		"may start with a utf-8 character": {
			input: "@łódź",
			want: []*token.Token{
				V(P(0, 8, 1, 1), token.INSTANCE_VARIABLE, "łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "@zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(P(0, 27, 1, 1), token.INSTANCE_VARIABLE, "zażółć_gęślą_jaźń"),
				T(P(28, 1, 1, 20), token.PLUS),
				V(P(30, 1, 1, 22), token.DEC_INT, "2"),
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
				T(P(0, 5, 1, 1), token.FALSE),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
