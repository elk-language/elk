package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestConstantAccess(t *testing.T) {
	tests := testTable{
		"access class constant": {
			input: "Int",
		},
		"access module constant": {
			input: "Std",
		},
		"access undefined constant": {
			input: "Foo",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(2, 1, 3)), "undefined constant `Foo`"),
			},
		},
		"constant lookup": {
			input: "Std::Int",
		},
		"constant lookup with error in the middle": {
			input: "Std::Foo::Bar",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(5, 1, 6), P(7, 1, 8)), "undefined constant `Std::Foo`"),
			},
		},
		"constant lookup with error at the start": {
			input: "Foo::Bar::Baz",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(2, 1, 3)), "undefined constant `Foo`"),
			},
		},
		"absolute constant lookup": {
			input: "::Std::Int",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestConstantDeclarations(t *testing.T) {
	tests := testTable{
		"declare with explicit type": {
			input: "const Foo: Int = 5",
		},
		"declare with implicit type": {
			input: "const Foo = 5",
		},
		"declare with implicit type and assign to literal type": {
			input: `
				const Foo = 5
				var foo: 5 = Foo
			`,
		},
		"declare with incorrect explicit type": {
			input: "const Foo: String = 5",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 1, 21), P(20, 1, 21)), "type `Std::Int(5)` cannot be assigned to type `Std::String`"),
			},
		},
		"declare without initialising": {
			input: "const Foo: String",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(16, 1, 17)), "constant `Foo` has to be initialised"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
