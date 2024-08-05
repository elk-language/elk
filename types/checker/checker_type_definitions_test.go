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
				sealed primitive class Std::String
					def foo; end
				end

				typedef Text = String
				var a: Text = "foo"
				a.foo
			`,
		},
		"define a type using a type before its declaration with difference": {
			input: `
				typedef Foo = Bar / nil
				typedef Bar = Baz | nil
				class Baz; end

				var b: Foo / 3 = 9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(98, 6, 22), P(100, 6, 24)), "type `9.2` cannot be assigned to type `Baz`"),
			},
		},
		"define a type using a type before its declaration with union": {
			input: `
				typedef Foo = Bar | 3
				typedef Bar = Baz | nil
				class Baz; end

				var b: Foo / 9 = 9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(96, 6, 22), P(98, 6, 24)), "type `9.2` cannot be assigned to type `3 | Baz | nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
