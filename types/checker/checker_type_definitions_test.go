package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestTypeDefinition(t *testing.T) {
	tests := testTable{
		"define types with circular dependencies": {
			input: `
				typedef Foo = Bar
				typedef Bar = Foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(41, 3, 19), P(43, 3, 21)), "Type `Foo` circularly references itself"),
			},
		},
		"define a type and assign a compatible value": {
			input: `
				typedef Text = String
				var a: Text = "foo"
			`,
		},
		"define a type and assign an incompatible value": {
			input: `
				typedef Text = String
				var a: Text = 1
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 19), P(45, 3, 19)), "type `1` cannot be assigned to type `Text`"),
			},
		},
		"call a method on a defined type": {
			input: `
				sealed primitive class Std::String < Value
					def foo; end
				end

				typedef Text = String
				var a: Text = "foo"
				a.foo
			`,
		},
		"define a type using a class before its declaration": {
			input: `
				typedef Bar = Baz | nil
				class Baz; end

				var b: Bar / 3 = 9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(70, 5, 22), P(72, 5, 24)), "type `9.2` cannot be assigned to type `Baz | nil`"),
			},
		},
		"define a type using another type before its declaration": {
			input: `
				typedef Foo = Bar | nil
				typedef Bar = 1 | 2
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestGenericTypeDefinition(t *testing.T) {
	tests := testTable{
		"define generic types with circular dependencies": {
			input: `
				typedef Foo[V] = V | Bar
				typedef Bar = Foo[Int]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 3, 19), P(50, 3, 21)), "Type `Foo` circularly references itself"),
			},
		},
		"define generic types with circular dependencies in the bounds": {
			input: `
				typedef Foo[V < Bar] = V | Float
				typedef Bar = Foo[Int]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(56, 3, 19), P(58, 3, 21)), "Type `Foo` circularly references itself"),
			},
		},
		"define a generic type with valid content": {
			input: `
				typedef Dupa[V] = V | String
			`,
		},
		"use a generic type under a namespace": {
			input: `
				module Foo
					typedef Bar[V] = V | String
				end
				var a: Foo::Bar[Int] = nil
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(84, 5, 28), P(86, 5, 30)), "type `nil` cannot be assigned to type `Std::Int | Std::String`"),
			},
		},
		"define a generic type with invalid content": {
			input: `
				typedef Dupa[V] = T | String
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(23, 2, 23), P(23, 2, 23)), "undefined type `T`"),
			},
		},
		"define a generic type with valid upper bound": {
			input: `
				typedef Dupa[V < Object] = V | String
			`,
		},
		"define a generic type with invalid upper bound": {
			input: `
				typedef Dupa[V < Foo] = V | String
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 2, 22), P(24, 2, 24)), "undefined type `Foo`"),
			},
		},
		"define a generic type with valid lower bound": {
			input: `
				typedef Dupa[V > Int] = V | String
			`,
		},
		"define a generic type with invalid lower bound": {
			input: `
				typedef Dupa[V > Foo] = V | String
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 2, 22), P(24, 2, 24)), "undefined type `Foo`"),
			},
		},
		"define a generic type with invalid upper and lower bound": {
			input: `
				typedef Dupa[V > Foo < Bar] = V | String
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 2, 22), P(24, 2, 24)), "undefined type `Foo`"),
				error.NewFailure(L("<main>", P(28, 2, 28), P(30, 2, 30)), "undefined type `Bar`"),
			},
		},

		"use a generic type with a valid type argument": {
			input: `
				typedef Foo[V] = V | String
				var a: Foo[Int] = 2.4
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(55, 3, 23), P(57, 3, 25)), "type `2.4` cannot be assigned to type `Std::Int | Std::String`"),
			},
		},
		"use a generic type with an invalid number of type arguments": {
			input: `
				typedef Foo[V] = V | String
				var a: Foo[Int, Float] = 2.4
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 12), P(46, 3, 14)), "generic type `Foo[V]` requires 1 type argument(s)"),
			},
		},
		"use a generic type with a satisfied upper bound": {
			input: `
				typedef Foo[V < Float] = V | String
				var a: Foo[Float] = 2.4
			`,
		},
		"use a generic type with an unsatisfied upper bound": {
			input: `
				typedef Foo[V < Float] = V | String
				var a: Foo[Int] = 2.4
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(56, 3, 16), P(58, 3, 18)), "type `Std::Int` does not satisfy the upper bound `Std::Float`"),
			},
		},
		"use a generic type with a satisfied lower bound": {
			input: `
				typedef Foo[V > Float] = V | String
				var a: Foo[Value] = 2.4
			`,
		},
		"use a generic type with an unsatisfied lower bound": {
			input: `
				typedef Foo[V > Float] = V | String
				var a: Foo[Int] = 2.4
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(56, 3, 16), P(58, 3, 18)), "type `Std::Int` does not satisfy the lower bound `Std::Float`"),
			},
		},
		"use a generic type with an unsatisfied upper and lower bound": {
			input: `
				typedef Foo[V > Float < Object] = V | String
				var a: Foo[Int] = 2.4
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(65, 3, 16), P(67, 3, 18)), "type `Std::Int` does not satisfy the upper bound `Std::Object`"),
				error.NewFailure(L("<main>", P(65, 3, 16), P(67, 3, 18)), "type `Std::Int` does not satisfy the lower bound `Std::Float`"),
			},
		},
		"use a generic type without type arguments": {
			input: `
				typedef Dupa[V] = V | String
				var a: Dupa = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 12), P(48, 3, 15)), "generic type `Dupa[V]` requires 1 type argument(s)"),
			},
		},
		"use a generic type under a namespace without type arguments": {
			input: `
				module Foo
					typedef Bar[V] = V | String
				end
				var a: Foo::Bar = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(68, 5, 12), P(75, 5, 19)), "generic type `Foo::Bar[V]` requires 1 type argument(s)"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
