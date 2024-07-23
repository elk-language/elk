package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestDoExpressions(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				a := 5
				do
					var b: Int = a
				end
			`,
		},
		"returns the last expression": {
			input: `
				a := 2
				var b: Int = do
					"foo" + "bar"
					a + 2
				end
			`,
		},
		"returns nil when empty": {
			input: `
				var b: nil = do; end
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestIfExpressions(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				if a
					var b: Int? = a
				else
					a
				end
			`,
		},
		"returns the last then expression if condition is truthy": {
			input: `
				a := 2
				var b: Int = if true
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(32, 3, 21), P(35, 3, 24)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"returns the last else expression if condition is truthy": {
			input: `
				a := 2
				var b: Float = if false
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(34, 3, 23), P(38, 3, 27)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"returns a union of both branches if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = if a
					"foo" + "bar"
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(84, 7, 7)), "type `Std::String | 2.2` cannot be assigned to type `8`"),
			},
		},
		"returns a union of all branches if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: Float? = nil
				var c: 8 = if a
					"foo" + "bar"
				else if b
					2.5
				else
					%/foo/
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 4, 16), P(134, 10, 7)), "type `Std::String | 2.5 | Std::Regex` cannot be assigned to type `8`"),
			},
		},
		"returns a union of then and nil if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = if a
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(66, 5, 7)), "type `Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"returns nil when empty": {
			input: `
				var a: Int? = nil
				var b: nil = if a; end
			`,
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a: Bool = false
				if a
					var b: true = a
				else
					var b: false = a
				end
			`,
		},
		"narrow nilable variable type by using truthiness": {
			input: `
				var a: Int? = nil
				if a
					var b: Int = a
				else
					var b: nil = a
				end
			`,
		},
		"narrow nilable variable type by using negated truthiness": {
			input: `
				var a: Int? = nil
				if !a
					var b: nil = a
				else
					var b: Int = a
				end
			`,
		},
		"narrow variable type by using <:": {
			input: `
				var a: Int? = nil
				if a <: Int
					var b: Int = a
				else
					var b: nil = a
				end
			`,
		},
		"narrow variable type by using <<:": {
			input: `
				var a: Int? = nil
				if a <<: Int
					var b: Int = a
				else
					var b: nil = a
				end
			`,
		},
		"narrow a few variables with &&": {
			input: `
				var a: Int? = nil
				var b: Bool = false
				if a && b
					var c: Int = a
					var d: true = b
				end
			`,
		},
		"narrow a few variables with ||": {
			input: `
				var a: Int? = nil
				var b: Bool = false
				if a || b
				else
					var c: nil = a
					var d: false = b
				end
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
