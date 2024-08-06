package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestTypeDefinition(t *testing.T) {
	tests := testTable{
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(19, 2, 19), P(21, 2, 21)), "undefined type `Bar`"),
			},
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
		"define a generic type with valid content": {
			input: `
				typedef Dupa[V] = V | String
			`,
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
