package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestDoExpression(t *testing.T) {
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

func TestNumericForExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				fornum ;;
					var b: Int? = a
				end
			`,
		},
		"use variables defined in the header": {
			input: `
				fornum i := 0; i < 8; i++
					var b: Int = i
				end
			`,
		},
		"typecheck the header and body": {
			input: `
				fornum a; b; c
					d
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 2, 12), P(12, 2, 12)), "undefined local `a`"),
				error.NewFailure(L("<main>", P(15, 2, 15), P(15, 2, 15)), "undefined local `b`"),
				error.NewFailure(L("<main>", P(18, 2, 18), P(18, 2, 18)), "undefined local `c`"),
				error.NewFailure(L("<main>", P(25, 3, 6), P(25, 3, 6)), "undefined local `d`"),
			},
		},
		"returns never if condition is truthy": {
			input: `
				a := 2
				b := fornum ;true;
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(29, 3, 18), P(32, 3, 21)), "this condition will always have the same result since type `true` is truthy"),
				error.NewFailure(L("<main>", P(81, 7, 9), P(81, 7, 9)), "type `3` cannot be assigned to type `never`"),
			},
		},
		"returns never when there is no condition": {
			input: `
				a := 2
				b := fornum ;;
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(77, 7, 9), P(77, 7, 9)), "type `3` cannot be assigned to type `never`"),
			},
		},
		"returns nil if condition is falsy": {
			input: `
				a := 2
				var b: nil = fornum ;false;
					"foo" + "bar"
					a + 2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(37, 3, 26), P(41, 3, 30)), "this loop will never execute since type `false` is falsy"),
			},
		},
		"returns a nilable body type if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = fornum ;a;
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(72, 5, 7)), "type `Std::String?` cannot be assigned to type `8`"),
			},
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				fornum ;a;
					var b: true = a
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

func TestWhileExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				while a
					var b: Int? = a
				end
			`,
		},
		"returns never if condition is truthy": {
			input: `
				a := 2
				b := while true
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 3, 16), P(30, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
				error.NewFailure(L("<main>", P(78, 7, 9), P(78, 7, 9)), "type `3` cannot be assigned to type `never`"),
			},
		},
		"returns nil expression if condition is falsy": {
			input: `
				a := 2
				var b: nil = while false
					"foo" + "bar"
					a + 2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(35, 3, 24), P(39, 3, 28)), "this loop will never execute since type `false` is falsy"),
			},
		},
		"returns a nilable body type if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = while a
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(69, 5, 7)), "type `Std::String?` cannot be assigned to type `8`"),
			},
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				while a
					var b: true = a
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

func TestUnlessExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				unless a
					var b: Int? = a
				else
					a
				end
			`,
		},
		"returns the last else expression if condition is truthy": {
			input: `
				a := 2
				var b: Float = unless true
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(38, 3, 27), P(41, 3, 30)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"returns the last then expression if condition is falsy": {
			input: `
				a := 2
				var b: Int = unless false
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(36, 3, 25), P(40, 3, 29)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"returns a union of both branches if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = unless a
					"foo" + "bar"
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(88, 7, 7)), "type `Std::String | 2.2` cannot be assigned to type `8`"),
			},
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				unless a
					var b: false = a
				else
					var b: true = a
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

func TestIfExpression(t *testing.T) {
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
				var a = false
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
		"narrow union type by using <:": {
			input: `
				var a: Int | String = "foo"
				if a <: Int
					var b: Int = a
				else
					var b: String = a
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
				var b = false
				if a && b
					var c: Int = a
					var d: true = b
				end
			`,
		},
		"narrow with an impossible && branch": {
			input: `
				var a: Int? = 3
				if a && nil
					a = :foo
				else
					a = :bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(46, 4, 10), P(49, 4, 13)), "type `:foo` cannot be assigned to type `never`"),
				error.NewFailure(L("<main>", P(69, 6, 10), P(72, 6, 13)), "type `:bar` cannot be assigned to type `Std::Int?`"),
				error.NewWarning(L("<main>", P(28, 3, 8), P(35, 3, 15)), "this condition will always have the same result since type `nil` is falsy"),
			},
		},
		"narrow a few variables with ||": {
			input: `
				var a: Int? = nil
				var b = false
				if a || b
				else
					var c: nil = a
					var d: false = b
				end
			`,
		},
		"narrow with an impossible || branch": {
			input: `
				var a: Int? = nil
				if a || 3
					a = :foo
				else
					a = :bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(46, 4, 10), P(49, 4, 13)), "type `:foo` cannot be assigned to type `Std::Int?`"),
				error.NewFailure(L("<main>", P(69, 6, 10), P(72, 6, 13)), "type `:bar` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(30, 3, 8), P(35, 3, 13)), "this condition will always have the same result since type `Std::Int` is truthy"),
			},
		},

		"narrow with ===": {
			input: `
				var a: Int | Float = 1
				var b: Float? = .2
				if a === b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(75, 5, 10), P(78, 5, 13)), "type `:foo` cannot be assigned to type `Std::Float`"),
				error.NewFailure(L("<main>", P(89, 6, 10), P(92, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float`"),
				error.NewFailure(L("<main>", P(112, 8, 10), P(115, 8, 13)), "type `:baz` cannot be assigned to type `Std::Int | Std::Float`"),
				error.NewFailure(L("<main>", P(126, 9, 10), P(130, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float?`"),
			},
		},
		"narrow with an impossible ===": {
			input: `
				var a: Int = 1
				var b: Float = .2
				if a === b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(49, 4, 8), P(49, 4, 8)), "this strict equality check is impossible, `Std::Int` cannot ever be equal to `Std::Float`"),
				error.NewFailure(L("<main>", P(66, 5, 10), P(69, 5, 13)), "type `:foo` cannot be assigned to type `never`"),
				error.NewFailure(L("<main>", P(80, 6, 10), P(83, 6, 13)), "type `:bar` cannot be assigned to type `never`"),
				error.NewFailure(L("<main>", P(103, 8, 10), P(106, 8, 13)), "type `:baz` cannot be assigned to type `Std::Int`"),
				error.NewFailure(L("<main>", P(117, 9, 10), P(121, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float`"),
			},
		},
		"narrow with !==": {
			input: `
				var a: Int | Float = 1
				var b: Float? = .2
				if a !== b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(75, 5, 10), P(78, 5, 13)), "type `:foo` cannot be assigned to type `Std::Int | Std::Float`"),
				error.NewFailure(L("<main>", P(89, 6, 10), P(92, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float?`"),
				error.NewFailure(L("<main>", P(112, 8, 10), P(115, 8, 13)), "type `:baz` cannot be assigned to type `Std::Float`"),
				error.NewFailure(L("<main>", P(126, 9, 10), P(130, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float`"),
			},
		},
		"narrow with an impossible !==": {
			input: `
				var a: Int = 1
				var b: Float = .2
				if a !== b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(49, 4, 8), P(49, 4, 8)), "this strict equality check is impossible, `Std::Int` cannot ever be equal to `Std::Float`"),
				error.NewFailure(L("<main>", P(66, 5, 10), P(69, 5, 13)), "type `:foo` cannot be assigned to type `Std::Int`"),
				error.NewFailure(L("<main>", P(80, 6, 10), P(83, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float`"),
				error.NewFailure(L("<main>", P(103, 8, 10), P(106, 8, 13)), "type `:baz` cannot be assigned to type `never`"),
				error.NewFailure(L("<main>", P(117, 9, 10), P(121, 9, 14)), "type `:fizz` cannot be assigned to type `never`"),
			},
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
				error.NewWarning(L("<main>", P(48, 4, 18), P(48, 4, 18)), "this condition will always have the same result since type `nil` is falsy"),
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
				error.NewWarning(L("<main>", P(48, 4, 18), P(48, 4, 18)), "this condition will always have the same result since type `nil` is falsy"),
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

func TestNilCoalescing(t *testing.T) {
	tests := testTable{
		"returns the left type when it is not nilable": {
			input: `
				var a = "foo"
				var b = 2
				var c: String = a ?? b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(53, 4, 21), P(53, 4, 21)), "this condition will always have the same result since type `Std::String` can never be nil"),
			},
		},
		"returns the right type when the left type is nil": {
			input: `
				var a = nil
				var b = 2
				var c: Int = a ?? b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(48, 4, 18), P(48, 4, 18)), "this condition will always have the same result"),
			},
		},
		"returns a union of both types without nil when the left can be both nil and not nil": {
			input: `
				var a: String? = "foo"
				var b = 2
				var c: 9 = a ?? b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 16), P(62, 4, 21)), "type `Std::String | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without duplication": {
			input: `
				var a: String | Int | nil = nil
				var b: Float | Int | nil = 2.2
				var c: 9 = a ?? b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(87, 4, 16), P(92, 4, 21)), "type `Std::String | Std::Int | Std::Float | nil` cannot be assigned to type `9`"),
			},
		},

		"narrow left variable to non nilable": {
			input: `
				var a: false | nil | Int = nil
				a ?? var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(56, 3, 21), P(56, 3, 21)), "type `nil` cannot be assigned to type `9`"),
			},
		},
		"narrow a few variables to non nilable": {
			input: `
				var a: Int? = nil
				var b: Int? = nil
				a ?? b ?? var c: 9 = a && b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(70, 4, 26), P(70, 4, 26)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewFailure(L("<main>", P(70, 4, 26), P(75, 4, 31)), "type `nil` cannot be assigned to type `9`"),
			},
		},
		"narrow nested ||": {
			input: `
				var a: bool? = false
				var b: bool? = false
				(a || b) ?? do
					a = :foo
					b = :bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(79, 5, 10), P(82, 5, 13)), "type `:foo` cannot be assigned to type `nil`"),
				error.NewFailure(L("<main>", P(93, 6, 10), P(96, 6, 13)), "type `:bar` cannot be assigned to type `nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
