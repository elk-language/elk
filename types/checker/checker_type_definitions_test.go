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
				error.NewFailure(L("<main>", P(45, 3, 19), P(45, 3, 19)), "type `Std::Int(1)` cannot be assigned to type `Text`"),
			},
		},
		"call a method on a defined type": {
			input: `
				sealed class Std::String
					def foo; end
				end

				typedef Text = String
				var a: Text = "foo"
				a.foo
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
