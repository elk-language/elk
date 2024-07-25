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
		"narrow named nilable variable type by using truthiness": {
			input: `
				typedef Foo = Int?
				var a: Foo = nil
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
		"narrow named nilable variable type by using negated truthiness": {
			input: `
				typedef Foo = Int?
				var a: Foo = nil
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

func TestLogicalAnd(t *testing.T) {
	tests := testTable{
		"returns the right type when the left type is truthy": {
			input: `
				var a = "foo"
				var b = 2
				var c: Int = a && b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(50, 4, 18), P(50, 4, 18)), "this condition will always have the same result since type `Std::String` is truthy"),
			},
		},
		"returns the left type when the left type is falsy": {
			input: `
				var a = nil
				var b = 2
				var c: nil = a && b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(48, 4, 18), P(48, 4, 18)), "this condition will always have the same result since type `Std::Nil` is falsy"),
			},
		},
		"returns a union of both types with only nil when the left can be both truthy and falsy": {
			input: `
				var a: String? = "foo"
				var b = 2
				var c: 9 = a && b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 16), P(62, 4, 21)), "type `nil | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types with only false when the left can be both truthy and falsy": {
			input: `
				var a = false
				var b = 2
				var c: 9 = a && b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 4, 16), P(53, 4, 21)), "type `false | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without duplication": {
			input: `
				var a: false | nil | Int = nil
				var b: Float | Int | nil = 2.2
				var c: 9 = a && b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(86, 4, 16), P(91, 4, 21)), "type `false | nil | Std::Float | Std::Int` cannot be assigned to type `9`"),
			},
		},

		"narrow left variable to non falsy": {
			input: `
				var a: false | nil | Int = nil
				a && a + 2
			`,
		},
		"narrow a few variables to non falsy": {
			input: `
				var a: Int? = nil
				var b: Int? = nil
				a && b && a + b
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestLogicalOr(t *testing.T) {
	tests := testTable{
		"returns the left type when it is truthy": {
			input: `
				var a = "foo"
				var b = 2
				var c: String = a || b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(53, 4, 21), P(53, 4, 21)), "this condition will always have the same result since type `Std::String` is truthy"),
			},
		},
		"returns the right type when the left type is falsy": {
			input: `
				var a = nil
				var b = 2
				var c: Int = a || b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(48, 4, 18), P(48, 4, 18)), "this condition will always have the same result since type `Std::Nil` is falsy"),
			},
		},
		"returns a union of both types without nil when the left can be both truthy and falsy": {
			input: `
				var a: String? = "foo"
				var b = 2
				var c: 9 = a || b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 16), P(62, 4, 21)), "type `Std::String | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without false when the left can be both truthy and falsy": {
			input: `
				var a = false
				var b = 2
				var c: 9 = a || b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 4, 16), P(53, 4, 21)), "type `true | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without duplication": {
			input: `
				var a: String | Int | nil = nil
				var b: Float | Int | nil = 2.2
				var c: 9 = a || b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(87, 4, 16), P(92, 4, 21)), "type `Std::String | Std::Int | Std::Float | nil` cannot be assigned to type `9`"),
			},
		},

		"narrow left variable to falsy": {
			input: `
				var a: false | nil | Int = nil
				a || var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(56, 3, 21), P(56, 3, 21)), "type `false | nil` cannot be assigned to type `9`"),
			},
		},
		"narrow a few variables to non falsy": {
			input: `
				var a: Int? = nil
				var b: Int? = nil
				a || b || var c: 9 = a && b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(70, 4, 26), P(70, 4, 26)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewFailure(L("<main>", P(70, 4, 26), P(75, 4, 31)), "type `nil` cannot be assigned to type `9`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}