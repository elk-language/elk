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
				V(L(S(P(0, 1, 1), P(2, 1, 3))), token.PUBLIC_IDENTIFIER, "foo"),
				T(L(S(P(3, 1, 4), P(3, 1, 4))), token.COLON),
				T(L(S(P(4, 1, 5), P(4, 1, 5))), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "some_identifier123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(17, 1, 18))), token.PUBLIC_IDENTIFIER, "some_identifier123"),
			},
		},
		"cannot start with numbers": {
			input: "3d_secure",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(0, 1, 1))), token.INT, "3"),
				V(L(S(P(1, 1, 2), P(8, 1, 9))), token.PUBLIC_IDENTIFIER, "d_secure"),
			},
		},
		"may contain utf-8 characters": {
			input: "zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(25, 1, 17))), token.PUBLIC_IDENTIFIER, "zażółć_gęślą_jaźń"),
				T(L(S(P(27, 1, 19), P(27, 1, 19))), token.PLUS),
				V(L(S(P(29, 1, 21), P(29, 1, 21))), token.INT, "2"),
			},
		},
		"may start with a utf-8 character": {
			input: "łódź",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 4))), token.PUBLIC_IDENTIFIER, "łódź"),
			},
		},
		"cannot start with an uppercase letter": {
			input: "Dupa",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.PUBLIC_CONSTANT, "Dupa"),
			},
		},
		"cannot start with an underscore": {
			input: "_foo",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.PRIVATE_IDENTIFIER, "_foo"),
			},
		},

		"dollar, ends on the last valid character": {
			input: "$foo:+",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.DOLLAR_IDENTIFIER, "foo"),
				T(L(S(P(4, 1, 5), P(4, 1, 5))), token.COLON),
				T(L(S(P(5, 1, 6), P(5, 1, 6))), token.PLUS),
			},
		},
		"dollar, may contain letters underscores and numbers": {
			input: "$some_ivar123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(12, 1, 13))), token.DOLLAR_IDENTIFIER, "some_ivar123"),
			},
		},
		"dollar, may start with an uppercase letter": {
			input: "$SomeIvar123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(11, 1, 12))), token.DOLLAR_IDENTIFIER, "SomeIvar123"),
			},
		},
		"dollar, may start with a digit": {
			input: "$1",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.DOLLAR_IDENTIFIER, "1"),
			},
		},
		"dollar, may start with an underscore": {
			input: "$_bar",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.DOLLAR_IDENTIFIER, "_bar"),
			},
		},
		"dollar, may start with a utf-8 character": {
			input: "$łódź",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(7, 1, 5))), token.DOLLAR_IDENTIFIER, "łódź"),
			},
		},
		"dollar, may contain utf-8 characters": {
			input: "$zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(26, 1, 18))), token.DOLLAR_IDENTIFIER, "zażółć_gęślą_jaźń"),
				T(L(S(P(28, 1, 20), P(28, 1, 20))), token.PLUS),
				V(L(S(P(30, 1, 22), P(30, 1, 22))), token.INT, "2"),
			},
		},

		"quoted, must be terminated": {
			input: `$"This is a string`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(17, 1, 18))), token.ERROR, "unterminated quoted identifier, missing `\"`"),
			},
		},
		"quoted, processes escape sequences": {
			input: `$"Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF\$\#"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(95, 1, 96))), token.DOLLAR_IDENTIFIER, "Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF$#"),
			},
		},
		"quoted, reports errors for invalid escape sequences": {
			input: `$"www.foo\yes.com"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(8, 1, 9))), token.PUBLIC_IDENTIFIER, "www.foo"),
				V(L(S(P(9, 1, 10), P(10, 1, 11))), token.ERROR, "invalid escape sequence `\\y` in string literal"),
				V(L(S(P(11, 1, 12), P(17, 1, 18))), token.PUBLIC_IDENTIFIER, "es.com"),
			},
		},
		"quoted, creates errors for invalid hex escapes": {
			input: `$"some\xfj string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.DOLLAR_IDENTIFIER, "some"),
				V(L(S(P(6, 1, 7), P(9, 1, 10))), token.ERROR, "invalid hex escape"),
				V(L(S(P(10, 1, 11), P(17, 1, 18))), token.DOLLAR_IDENTIFIER, " string"),
			},
		},
		"quoted, creates errors for invalid unicode escapes": {
			input: `$"some\uiaab string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.DOLLAR_IDENTIFIER, "some"),
				V(L(S(P(6, 1, 7), P(11, 1, 12))), token.ERROR, "invalid unicode escape"),
				V(L(S(P(12, 1, 13), P(19, 1, 20))), token.DOLLAR_IDENTIFIER, " string"),
			},
		},
		"quoted, creates errors for invalid big unicode escapes": {
			input: `$"some\Uiaabuj46 string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.DOLLAR_IDENTIFIER, "some"),
				V(L(S(P(6, 1, 7), P(15, 1, 16))), token.ERROR, "invalid unicode escape"),
				V(L(S(P(16, 1, 17), P(23, 1, 24))), token.DOLLAR_IDENTIFIER, " string"),
			},
		},
		"quoted, can be multiline": {
			input: `$"multiline
strings are
awesome
and really useful"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(49, 4, 18))), token.DOLLAR_IDENTIFIER, "multiline\nstrings are\nawesome\nand really useful"),
			},
		},

		"raw quoted, must be terminated": {
			input: "$'This is a raw string",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(21, 1, 22))), token.ERROR, "unterminated raw quoted identifier, missing `'`"),
			},
		},
		"raw quoted, does not process escape sequences": {
			input: `$'Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a'`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(72, 1, 73))), token.DOLLAR_IDENTIFIER, `Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a`),
			},
		},
		"raw quoted, can be multiline": {
			input: `$'multiline
strings are
awesome
and really useful'`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(49, 4, 18))), token.DOLLAR_IDENTIFIER, "multiline\nstrings are\nawesome\nand really useful"),
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
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.PRIVATE_IDENTIFIER, "_foo"),
				T(L(S(P(4, 1, 5), P(4, 1, 5))), token.COLON),
				T(L(S(P(5, 1, 6), P(5, 1, 6))), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "_some_identifier123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(18, 1, 19))), token.PRIVATE_IDENTIFIER, "_some_identifier123"),
			},
		},
		"may start with a utf-8 character": {
			input: "_łódź",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(7, 1, 5))), token.PRIVATE_IDENTIFIER, "_łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "_zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(26, 1, 18))), token.PRIVATE_IDENTIFIER, "_zażółć_gęślą_jaźń"),
				T(L(S(P(28, 1, 20), P(28, 1, 20))), token.PLUS),
				V(L(S(P(30, 1, 22), P(30, 1, 22))), token.INT, "2"),
			},
		},
		"cannot start with an uppercase letter": {
			input: "_Dupa",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.PRIVATE_CONSTANT, "_Dupa"),
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
				V(L(S(P(0, 1, 1), P(2, 1, 3))), token.PUBLIC_CONSTANT, "Foo"),
				T(L(S(P(3, 1, 4), P(3, 1, 4))), token.COLON),
				T(L(S(P(4, 1, 5), P(4, 1, 5))), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "Some_constant123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(15, 1, 16))), token.PUBLIC_CONSTANT, "Some_constant123"),
			},
		},
		"cannot start with numbers": {
			input: "3DSecure",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(0, 1, 1))), token.INT, "3"),
				V(L(S(P(1, 1, 2), P(7, 1, 8))), token.PUBLIC_CONSTANT, "DSecure"),
			},
		},
		"may contain utf-8 characters": {
			input: "ZażółćGęśląJaźń + 2",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(23, 1, 15))), token.PUBLIC_CONSTANT, "ZażółćGęśląJaźń"),
				T(L(S(P(25, 1, 17), P(25, 1, 17))), token.PLUS),
				V(L(S(P(27, 1, 19), P(27, 1, 19))), token.INT, "2"),
			},
		},
		"may start with a utf-8 character": {
			input: "Łódź",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 4))), token.PUBLIC_CONSTANT, "Łódź"),
			},
		},
		"cannot end with a question mark": {
			input: "Includes?",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(7, 1, 8))), token.PUBLIC_CONSTANT, "Includes"),
				T(L(S(P(8, 1, 9), P(8, 1, 9))), token.QUESTION),
			},
		},
		"cannot end with an exclamation point": {
			input: "Map!",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(2, 1, 3))), token.PUBLIC_CONSTANT, "Map"),
				T(L(S(P(3, 1, 4), P(3, 1, 4))), token.BANG),
			},
		},
		"cannot start with an underscore": {
			input: "_Foo",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.PRIVATE_CONSTANT, "_Foo"),
			},
		},

		"dollar, ends on the last valid character": {
			input: "$$foo:+",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.PUBLIC_CONSTANT, "foo"),
				T(L(S(P(5, 1, 6), P(5, 1, 6))), token.COLON),
				T(L(S(P(6, 1, 7), P(6, 1, 7))), token.PLUS),
			},
		},
		"dollar, may contain letters underscores and numbers": {
			input: "$$some_ivar123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(13, 1, 14))), token.PUBLIC_CONSTANT, "some_ivar123"),
			},
		},
		"dollar, may start with an uppercase letter": {
			input: "$$SomeIvar123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(12, 1, 13))), token.PUBLIC_CONSTANT, "SomeIvar123"),
			},
		},
		"dollar, may start with a digit": {
			input: "$$1",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(2, 1, 3))), token.PUBLIC_CONSTANT, "1"),
			},
		},
		"dollar, may start with an underscore": {
			input: "$$_bar",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.PUBLIC_CONSTANT, "_bar"),
			},
		},
		"dollar, may start with a utf-8 character": {
			input: "$$łódź",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(8, 1, 6))), token.PUBLIC_CONSTANT, "łódź"),
			},
		},
		"dollar, may contain utf-8 characters": {
			input: "$$zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(27, 1, 19))), token.PUBLIC_CONSTANT, "zażółć_gęślą_jaźń"),
				T(L(S(P(29, 1, 21), P(29, 1, 21))), token.PLUS),
				V(L(S(P(31, 1, 23), P(31, 1, 23))), token.INT, "2"),
			},
		},

		"quoted, must be terminated": {
			input: `$$"This is a string`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(18, 1, 19))), token.ERROR, "unterminated quoted constant, missing `\"`"),
			},
		},
		"quoted, processes escape sequences": {
			input: `$$"Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF\$\#"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(96, 1, 97))), token.PUBLIC_CONSTANT, "Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF$#"),
			},
		},
		"quoted, reports errors for invalid escape sequences": {
			input: `$$"www.foo\yes.com"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(9, 1, 10))), token.PUBLIC_CONSTANT, "www.foo"),
				V(L(S(P(10, 1, 11), P(11, 1, 12))), token.ERROR, "invalid escape sequence `\\y` in string literal"),
				V(L(S(P(12, 1, 13), P(18, 1, 19))), token.PUBLIC_CONSTANT, "es.com"),
			},
		},
		"quoted, creates errors for invalid hex escapes": {
			input: `$$"some\xfj string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 7))), token.PUBLIC_CONSTANT, "some"),
				V(L(S(P(7, 1, 8), P(10, 1, 11))), token.ERROR, "invalid hex escape"),
				V(L(S(P(11, 1, 12), P(18, 1, 19))), token.PUBLIC_CONSTANT, " string"),
			},
		},
		"quoted, creates errors for invalid unicode escapes": {
			input: `$$"some\uiaab string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 7))), token.PUBLIC_CONSTANT, "some"),
				V(L(S(P(7, 1, 8), P(12, 1, 13))), token.ERROR, "invalid unicode escape"),
				V(L(S(P(13, 1, 14), P(20, 1, 21))), token.PUBLIC_CONSTANT, " string"),
			},
		},
		"quoted, creates errors for invalid big unicode escapes": {
			input: `$$"some\Uiaabuj46 string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 7))), token.PUBLIC_CONSTANT, "some"),
				V(L(S(P(7, 1, 8), P(16, 1, 17))), token.ERROR, "invalid unicode escape"),
				V(L(S(P(17, 1, 18), P(24, 1, 25))), token.PUBLIC_CONSTANT, " string"),
			},
		},
		"quoted, can be multiline": {
			input: `$$"multiline
strings are
awesome
and really useful"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(50, 4, 18))), token.PUBLIC_CONSTANT, "multiline\nstrings are\nawesome\nand really useful"),
			},
		},

		"raw quoted, must be terminated": {
			input: "$$'This is a raw string",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(22, 1, 23))), token.ERROR, "unterminated raw quoted constant, missing `'`"),
			},
		},
		"raw quoted, does not process escape sequences": {
			input: `$$'Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a'`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(73, 1, 74))), token.PUBLIC_CONSTANT, `Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a`),
			},
		},
		"raw quoted, can be multiline": {
			input: `$$'multiline
strings are
awesome
and really useful'`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(50, 4, 18))), token.PUBLIC_CONSTANT, "multiline\nstrings are\nawesome\nand really useful"),
			},
		},

		"section sign, ends on the last valid character": {
			input: "§foo:+",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 4))), token.PUBLIC_CONSTANT, "foo"),
				T(L(S(P(5, 1, 5), P(5, 1, 5))), token.COLON),
				T(L(S(P(6, 1, 6), P(6, 1, 6))), token.PLUS),
			},
		},
		"section sign, may contain letters underscores and numbers": {
			input: "§some_ivar123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(13, 1, 13))), token.PUBLIC_CONSTANT, "some_ivar123"),
			},
		},
		"section sign, may start with an uppercase letter": {
			input: "§SomeIvar123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(12, 1, 12))), token.PUBLIC_CONSTANT, "SomeIvar123"),
			},
		},
		"section sign, may start with a digit": {
			input: "§1",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(2, 1, 2))), token.PUBLIC_CONSTANT, "1"),
			},
		},
		"section sign, may start with an underscore": {
			input: "§_bar",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 5))), token.PUBLIC_CONSTANT, "_bar"),
			},
		},
		"section sign, may start with a utf-8 character": {
			input: "§łódź",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(8, 1, 5))), token.PUBLIC_CONSTANT, "łódź"),
			},
		},
		"section sign, may contain utf-8 characters": {
			input: "§zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(27, 1, 18))), token.PUBLIC_CONSTANT, "zażółć_gęślą_jaźń"),
				T(L(S(P(29, 1, 20), P(29, 1, 20))), token.PLUS),
				V(L(S(P(31, 1, 22), P(31, 1, 22))), token.INT, "2"),
			},
		},

		"quoted section sign, must be terminated": {
			input: `§"This is a string`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(18, 1, 18))), token.ERROR, "unterminated quoted constant, missing `\"`"),
			},
		},
		"quoted section sign, processes escape sequences": {
			input: `§"Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF\$\#"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(96, 1, 96))), token.PUBLIC_CONSTANT, "Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF$#"),
			},
		},
		"quoted section sign, reports errors for invalid escape sequences": {
			input: `§"www.foo\yes.com"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(9, 1, 9))), token.PUBLIC_CONSTANT, "www.foo"),
				V(L(S(P(10, 1, 10), P(11, 1, 11))), token.ERROR, "invalid escape sequence `\\y` in string literal"),
				V(L(S(P(12, 1, 12), P(18, 1, 18))), token.PUBLIC_CONSTANT, "es.com"),
			},
		},
		"quoted section sign, creates errors for invalid hex escapes": {
			input: `§"some\xfj string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 6))), token.PUBLIC_CONSTANT, "some"),
				V(L(S(P(7, 1, 7), P(10, 1, 10))), token.ERROR, "invalid hex escape"),
				V(L(S(P(11, 1, 11), P(18, 1, 18))), token.PUBLIC_CONSTANT, " string"),
			},
		},
		"quoted section sign, creates errors for invalid unicode escapes": {
			input: `§"some\uiaab string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 6))), token.PUBLIC_CONSTANT, "some"),
				V(L(S(P(7, 1, 7), P(12, 1, 12))), token.ERROR, "invalid unicode escape"),
				V(L(S(P(13, 1, 13), P(20, 1, 20))), token.PUBLIC_CONSTANT, " string"),
			},
		},
		"quoted section sign, creates errors for invalid big unicode escapes": {
			input: `§"some\Uiaabuj46 string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(6, 1, 6))), token.PUBLIC_CONSTANT, "some"),
				V(L(S(P(7, 1, 7), P(16, 1, 16))), token.ERROR, "invalid unicode escape"),
				V(L(S(P(17, 1, 17), P(24, 1, 24))), token.PUBLIC_CONSTANT, " string"),
			},
		},
		"quoted section sign, can be multiline": {
			input: `§"multiline
strings are
awesome
and really useful"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(50, 4, 18))), token.PUBLIC_CONSTANT, "multiline\nstrings are\nawesome\nand really useful"),
			},
		},

		"raw quoted section sign, must be terminated": {
			input: "§'This is a raw string",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(22, 1, 22))), token.ERROR, "unterminated raw quoted constant, missing `'`"),
			},
		},
		"raw quoted section sign, does not process escape sequences": {
			input: `§'Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a'`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(73, 1, 73))), token.PUBLIC_CONSTANT, `Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a`),
			},
		},
		"raw quoted section sign, can be multiline": {
			input: `§'multiline
strings are
awesome
and really useful'`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(50, 4, 18))), token.PUBLIC_CONSTANT, "multiline\nstrings are\nawesome\nand really useful"),
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
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.PRIVATE_CONSTANT, "_Foo"),
				T(L(S(P(4, 1, 5), P(4, 1, 5))), token.COLON),
				T(L(S(P(5, 1, 6), P(5, 1, 6))), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "_Some_identifier123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(18, 1, 19))), token.PRIVATE_CONSTANT, "_Some_identifier123"),
			},
		},
		"may start with a utf-8 character": {
			input: "_Łódź",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(7, 1, 5))), token.PRIVATE_CONSTANT, "_Łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "_Zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(26, 1, 18))), token.PRIVATE_CONSTANT, "_Zażółć_gęślą_jaźń"),
				T(L(S(P(28, 1, 20), P(28, 1, 20))), token.PLUS),
				V(L(S(P(30, 1, 22), P(30, 1, 22))), token.INT, "2"),
			},
		},
		"cannot end with a question mark": {
			input: "_Includes?",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(8, 1, 9))), token.PRIVATE_CONSTANT, "_Includes"),
				T(L(S(P(9, 1, 10), P(9, 1, 10))), token.QUESTION),
			},
		},
		"cannot end with an exclamation point": {
			input: "_Map!",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.PRIVATE_CONSTANT, "_Map"),
				T(L(S(P(4, 1, 5), P(4, 1, 5))), token.BANG),
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
				V(L(S(P(0, 1, 1), P(3, 1, 4))), token.INSTANCE_VARIABLE, "foo"),
				T(L(S(P(4, 1, 5), P(4, 1, 5))), token.COLON),
				T(L(S(P(5, 1, 6), P(5, 1, 6))), token.PLUS),
			},
		},
		"may contain letters underscores and numbers": {
			input: "@some_ivar123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(12, 1, 13))), token.INSTANCE_VARIABLE, "some_ivar123"),
			},
		},
		"may start with an uppercase letter": {
			input: "@SomeIvar123",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(11, 1, 12))), token.INSTANCE_VARIABLE, "SomeIvar123"),
			},
		},
		"may start with a digit": {
			input: "@1",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(1, 1, 2))), token.INSTANCE_VARIABLE, "1"),
			},
		},
		"may start with an underscore": {
			input: "@_bar",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(4, 1, 5))), token.INSTANCE_VARIABLE, "_bar"),
			},
		},
		"may start with a utf-8 character": {
			input: "@łódź",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(7, 1, 5))), token.INSTANCE_VARIABLE, "łódź"),
			},
		},
		"may contain utf-8 characters": {
			input: "@zażółć_gęślą_jaźń + 2",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(26, 1, 18))), token.INSTANCE_VARIABLE, "zażółć_gęślą_jaźń"),
				T(L(S(P(28, 1, 20), P(28, 1, 20))), token.PLUS),
				V(L(S(P(30, 1, 22), P(30, 1, 22))), token.INT, "2"),
			},
		},

		"quoted, must be terminated": {
			input: `@"This is a string`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(17, 1, 18))), token.ERROR, "unterminated quoted instance variable, missing `\"`"),
			},
		},
		"quoted, processes escape sequences": {
			input: `@"Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF\$\#"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(95, 1, 96))), token.INSTANCE_VARIABLE, "Some \n a\\wesome \t str\\ing \r with \\ escape \b sequences \"\v\f\x12\a\u00e9\U0010FFFF$#"),
			},
		},
		"quoted, reports errors for invalid escape sequences": {
			input: `@"www.foo\yes.com"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(8, 1, 9))), token.INSTANCE_VARIABLE, "www.foo"),
				V(L(S(P(9, 1, 10), P(10, 1, 11))), token.ERROR, "invalid escape sequence `\\y` in string literal"),
				V(L(S(P(11, 1, 12), P(17, 1, 18))), token.INSTANCE_VARIABLE, "es.com"),
			},
		},
		"quoted, creates errors for invalid hex escapes": {
			input: `@"some\xfj string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INSTANCE_VARIABLE, "some"),
				V(L(S(P(6, 1, 7), P(9, 1, 10))), token.ERROR, "invalid hex escape"),
				V(L(S(P(10, 1, 11), P(17, 1, 18))), token.INSTANCE_VARIABLE, " string"),
			},
		},
		"quoted, creates errors for invalid unicode escapes": {
			input: `@"some\uiaab string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INSTANCE_VARIABLE, "some"),
				V(L(S(P(6, 1, 7), P(11, 1, 12))), token.ERROR, "invalid unicode escape"),
				V(L(S(P(12, 1, 13), P(19, 1, 20))), token.INSTANCE_VARIABLE, " string"),
			},
		},
		"quoted, creates errors for invalid big unicode escapes": {
			input: `@"some\Uiaabuj46 string"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(5, 1, 6))), token.INSTANCE_VARIABLE, "some"),
				V(L(S(P(6, 1, 7), P(15, 1, 16))), token.ERROR, "invalid unicode escape"),
				V(L(S(P(16, 1, 17), P(23, 1, 24))), token.INSTANCE_VARIABLE, " string"),
			},
		},
		"quoted, can be multiline": {
			input: `@"multiline
strings are
awesome
and really useful"`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(49, 4, 18))), token.INSTANCE_VARIABLE, "multiline\nstrings are\nawesome\nand really useful"),
			},
		},

		"raw quoted, must be terminated": {
			input: "@'This is a raw string",
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(21, 1, 22))), token.ERROR, "unterminated raw quoted instance variable, missing `'`"),
			},
		},
		"raw quoted, does not process escape sequences": {
			input: `@'Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a'`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(72, 1, 73))), token.INSTANCE_VARIABLE, `Some \n a\wesome \t string \r with \\ escape \b sequences \"\v\f\x12\a`),
			},
		},
		"raw quoted, can be multiline": {
			input: `@'multiline
strings are
awesome
and really useful'`,
			want: []*token.Token{
				V(L(S(P(0, 1, 1), P(49, 4, 18))), token.INSTANCE_VARIABLE, "multiline\nstrings are\nawesome\nand really useful"),
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
				T(L(S(P(0, 1, 1), P(4, 1, 5))), token.FALSE),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
