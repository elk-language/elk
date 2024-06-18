package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestVariableDeclaration(t *testing.T) {
	tests := testTable{
		"accept variable declaration with matching initializer and type": {
			input: "var foo: Int = 5",
		},
		"accept variable declaration with inference": {
			input: "var foo = 5",
		},
		"cannot declare variable with type void": {
			input: `
				def bar; end
				var foo = bar()
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(32, 3, 15), P(36, 3, 19)), "cannot declare variable `foo` with type `void`"),
			},
		},
		"reject variable declaration without matching initializer and type": {
			input: "var foo: Int = 5.2",
			err: error.ErrorList{
				error.NewError(L("<main>", P(15, 1, 16), P(17, 1, 18)), "type `Std::Float(5.2)` cannot be assigned to type `Std::Int`"),
			},
		},
		"accept variable declaration without initializer": {
			input: "var foo: Int",
		},
		"reject variable declaration with invalid type": {
			input: "var foo: Foo",
			err: error.ErrorList{
				error.NewError(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
			},
		},
		"reject variable declaration without initializer and type": {
			input: "var foo",
			err: error.ErrorList{
				error.NewError(L("<main>", P(0, 1, 1), P(6, 1, 7)), "cannot declare a variable without a type `foo`"),
			},
		},
		"reject redeclared variable": {
			input: "var foo: Int; var foo: String",
			err: error.ErrorList{
				error.NewError(L("<main>", P(14, 1, 15), P(28, 1, 29)), "cannot redeclare local `foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestValueDeclaration(t *testing.T) {
	tests := testTable{
		"accept value declaration with matching initializer and type": {
			input: "val foo: Int = 5",
		},
		"accept variable declaration with inference": {
			input: "val foo = 5",
		},
		"cannot declare value with type void": {
			input: `
				def bar; end
				val foo = bar()
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(32, 3, 15), P(36, 3, 19)), "cannot declare value `foo` with type `void`"),
			},
		},
		"reject value declaration without matching initializer and type": {
			input: "val foo: Int = 5.2",
			err: error.ErrorList{
				error.NewError(L("<main>", P(15, 1, 16), P(17, 1, 18)), "type `Std::Float(5.2)` cannot be assigned to type `Std::Int`"),
			},
		},
		"accept value declaration without initializer": {
			input: "val foo: Int",
		},
		"reject value declaration with invalid type": {
			input: "val foo: Foo",
			err: error.ErrorList{
				error.NewError(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
			},
		},
		"reject value declaration without initializer and type": {
			input: "val foo",
			err: error.ErrorList{
				error.NewError(L("<main>", P(0, 1, 1), P(6, 1, 7)), "cannot declare a value without a type `foo`"),
			},
		},
		"reject redeclared value": {
			input: "val foo: Int; val foo: String",
			err: error.ErrorList{
				error.NewError(L("<main>", P(14, 1, 15), P(28, 1, 29)), "cannot redeclare local `foo`"),
			},
		},
		"declaration with type lookup": {
			input: "val foo: Std::Int",
		},
		"declaration with type lookup and error in the middle": {
			input: "val foo: Std::Foo::Bar",
			err: error.ErrorList{
				error.NewError(L("<main>", P(14, 1, 15), P(16, 1, 17)), "undefined type `Std::Foo`"),
			},
		},
		"declaration with type lookup and error at the start": {
			input: "val foo: Foo::Bar::Baz",
			err: error.ErrorList{
				error.NewError(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
			},
		},
		"declaration with absolute type lookup": {
			input: "val foo: ::Std::Int",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestLocalAccess(t *testing.T) {
	tests := testTable{
		"access initialised variable": {
			input: "var foo: Int = 5; foo",
		},
		"access uninitialised variable": {
			input: "var foo: Int; foo",
			err: error.ErrorList{
				error.NewError(L("<main>", P(14, 1, 15), P(16, 1, 17)), "cannot access uninitialised local `foo`"),
			},
		},
		"access initialised value": {
			input: "val foo: Int = 5; foo",
		},
		"access uninitialised value": {
			input: "val foo: Int; foo",
			err: error.ErrorList{
				error.NewError(L("<main>", P(14, 1, 15), P(16, 1, 17)), "cannot access uninitialised local `foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
