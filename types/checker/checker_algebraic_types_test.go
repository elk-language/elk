package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestNilableSubtype(t *testing.T) {
	tests := testTable{
		"assign String to nilable String": {
			input: `
				var a = "foo"
				var b: String? = a
			`,
		},
		"assign nil to nilable String": {
			input: `
				var a = nil
				var b: String? = a
			`,
		},
		"assign Int to nilable String": {
			input: `
				var a = 3
				var b: String? = a
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(36, 3, 22), P(36, 3, 22)), "type `Std::Int` cannot be assigned to type `Std::String?`"),
			},
		},
		"assign nilable String to union type with String and nil": {
			input: `
				var a: String? = "foo"
				var b: String | Float | nil = a
			`,
		},
		"assign nilable String to union type without nil": {
			input: `
				var a: String? = "foo"
				var b: String | Float = a
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(56, 3, 29), P(56, 3, 29)), "type `Std::String?` cannot be assigned to type `Std::String | Std::Float`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestNilableTypeMethodCall(t *testing.T) {
	tests := testTable{
		"missing method on nil": {
			input: `
			  class Foo
				  def foo; end
				end
				var a: Foo? = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(69, 6, 5), P(73, 6, 9)), "method `foo` is not defined on type `Std::Nil`"),
			},
		},
		"missing method on nilable type": {
			input: `
			  class Foo; end
				class Std::Nil
				  def foo; end
				end
				var a: Foo? = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(93, 7, 5), P(97, 7, 9)), "method `foo` is not defined on type `Foo`"),
			},
		},
		"missing method on both types": {
			input: `
			  class Foo; end
				var a: Foo? = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(47, 4, 5), P(51, 4, 9)), "method `foo` is not defined on type `Std::Nil`"),
				error.NewError(L("<main>", P(47, 4, 5), P(51, 4, 9)), "method `foo` is not defined on type `Foo`"),
			},
		},
		"method with different number of arguments": {
			input: `
				class Foo
					def foo(a: Int, b: String); end
				end
				class Std::Nil
					def foo(a: Int); end
				end
				var a: Foo? = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(139, 9, 5), P(143, 9, 9)), "method `Foo.:foo` has a required parameter missing in `Std::Nil.:foo`, got `b`"),
			},
		},
		"method with different return types": {
			input: `
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Int): Nil then nil
				end
				var a: Foo? = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(144, 9, 5), P(148, 9, 9)), "method `Std::Nil.:foo` has a different return type than `Foo.:foo`, has `Std::Nil`, should have `Std::Int`"),
			},
		},
		"method with different param types": {
			input: `
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Float): Int then 5
				end
				var a: Foo? = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(144, 9, 5), P(148, 9, 9)), "method `Std::Nil.:foo` has a different type for parameter `a` than `Foo.:foo`, has `Std::Float`, should have `Std::Int`"),
			},
		},
		"method with additional optional params": {
			input: `
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Int, b: Float = .2): Int then a
				end
				var a: Foo? = nil
				a.foo(5)
			`,
		},
		"method with additional optional params in call": {
			input: `
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Int, b: Float = .2): Int then a
				end
				var a: Foo? = nil
				a.foo(5, 2.5)
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(157, 9, 5), P(169, 9, 17)), "expected 1 arguments in call to `foo`, got 2"),
			},
		},
		"method with additional rest param": {
			input: `
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Int, *b: Float): Int then a
				end
				var a: Foo? = nil
				a.foo(5, 2.5)
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(153, 9, 5), P(165, 9, 17)), "method `Std::Nil.:foo` has a required parameter missing in `Foo.:foo`, got `b`"),
			},
		},
		"method with additional named rest param": {
			input: `
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Int, **b: Float): Int then a
				end
				var a: Foo? = nil
				a.foo(5, a: 2.5)
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(154, 9, 5), P(169, 9, 20)), "method `Std::Nil.:foo` has a required parameter missing in `Foo.:foo`, got `b`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestUnionTypeSubtype(t *testing.T) {
	tests := testTable{
		"assign Int to union type with Int": {
			input: `
				var a = 3
				var b: String | Int = a
			`,
		},
		"assign Int to union type without Int": {
			input: `
				var a = 3
				var b: String | Float = a
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(43, 3, 29), P(43, 3, 29)), "type `Std::Int` cannot be assigned to type `Std::String | Std::Float`"),
			},
		},
		"assign union type to non union type": {
			input: `
				var a = 3
				var b: String | Float = a
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(43, 3, 29), P(43, 3, 29)), "type `Std::Int` cannot be assigned to type `Std::String | Std::Float`"),
			},
		},
		"assign union type to more general union type": {
			input: `
				var a: String | Int = 3
				var b: Int | Float | String = a
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestUnionTypeMethodCall(t *testing.T) {
	tests := testTable{
		"missing method": {
			input: `
			  class Foo
				  def foo; end
				end
				class Bar
				end
				var a: Foo | Bar | Nil = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(102, 8, 5), P(106, 8, 9)), "method `foo` is not defined on type `Bar`"),
				error.NewError(L("<main>", P(102, 8, 5), P(106, 8, 9)), "method `foo` is not defined on type `Std::Nil`"),
			},
		},
		"method with different number of arguments": {
			input: `
				class Bar
					def foo(a: Int, b: String, c: String); end
				end
				class Foo
					def foo(a: Int, b: String); end
				end
				class Std::Nil
					def foo(a: Int); end
				end
				var a: Foo | Bar | Nil = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(220, 12, 5), P(224, 12, 9)), "method `Foo.:foo` has a required parameter missing in `Std::Nil.:foo`, got `b`"),
				error.NewError(L("<main>", P(220, 12, 5), P(224, 12, 9)), "method `Bar.:foo` has a required parameter missing in `Std::Nil.:foo`, got `b`"),
				error.NewError(L("<main>", P(220, 12, 5), P(224, 12, 9)), "method `Bar.:foo` has a required parameter missing in `Std::Nil.:foo`, got `c`"),
			},
		},
		"method with different return types": {
			input: `
				class Bar
					def foo(a: Int): String then a
				end
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Int): Nil then nil
				end
				var a: Foo | Bar | Nil = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(213, 12, 5), P(217, 12, 9)), "method `Bar.:foo` has a different return type than `Foo.:foo`, has `Std::String`, should have `Std::Int`"),
				error.NewError(L("<main>", P(213, 12, 5), P(217, 12, 9)), "method `Std::Nil.:foo` has a different return type than `Foo.:foo`, has `Std::Nil`, should have `Std::Int`"),
			},
		},
		"method with different param types": {
			input: `
				class Bar
					def foo(a: String): Int then a
				end
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Float): Int then 5
				end
				var a: Foo | Bar | Nil = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(213, 12, 5), P(217, 12, 9)), "method `Bar.:foo` has a different type for parameter `a` than `Foo.:foo`, has `Std::String`, should have `Std::Int`"),
				error.NewError(L("<main>", P(213, 12, 5), P(217, 12, 9)), "method `Std::Nil.:foo` has a different type for parameter `a` than `Foo.:foo`, has `Std::Float`, should have `Std::Int`"),
			},
		},
		"method with wider param type": {
			input: `
				class Foo
					def foo(a: String): Int then a
				end
				class Bar
					def foo(a: Object): Int then 5
				end
				var a: Foo | Bar = Foo()
				a.foo("b")
			`,
		},
		"method with narrower param type": {
			input: `
				class Foo
					def foo(a: String): Int then a
				end
				class Bar
					def foo(a: Object): Int then 5
				end
				var a: Bar | Foo = Foo()
				a.foo("b")
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(150, 9, 5), P(159, 9, 14)), "method `Foo.:foo` has a different type for parameter `a` than `Bar.:foo`, has `Std::String`, should have `Std::Object`"),
			},
		},
		"method with wider return type": {
			input: `
				class Foo
					def foo(a: String): Int then a
				end
				class Bar
					def foo(a: String): Object then 5
				end
				var a: Foo | Bar = Foo()
				a.foo("b")
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(153, 9, 5), P(162, 9, 14)), "method `Bar.:foo` has a different return type than `Foo.:foo`, has `Std::Object`, should have `Std::Int`"),
			},
		},
		"method with narrower return type": {
			input: `
				class Foo
					def foo(a: String): Int then a
				end
				class Bar
					def foo(a: String): Object then 5
				end
				var a: Bar | Foo = Foo()
				a.foo("b")
			`,
		},
		"method with additional optional params": {
			input: `
				class Bar
					def foo(a: Int, c: String = "c"): Int then a
				end
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Int, b: Float = .2): Int then a
				end
				var a: Foo | Bar | Nil = nil
				a.foo(5)
			`,
		},
		"method with additional optional params in call": {
			input: `
				class Bar
					def foo(a: Int, c: String = "c"): Int then a
				end
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Int, b: Float = .2): Int then a
				end
				var a: Foo | Bar | Nil = nil
				a.foo(5, 2.5)
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(240, 12, 5), P(252, 12, 17)), "expected 1 arguments in call to `foo`, got 2"),
			},
		},
		"method with additional rest param": {
			input: `
				class Bar
					def foo(a: Int): Int then a
				end
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Int, *b: Float): Int then a
				end
				var a: Foo | Bar | Nil = nil
				a.foo(5, 2.5)
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(219, 12, 5), P(231, 12, 17)), "method `Std::Nil.:foo` has a required parameter missing in `Foo.:foo`, got `b`"),
			},
		},
		"method with additional named rest param": {
			input: `
				class Bar
					def foo(a: Int): Int then a
				end
				class Foo
					def foo(a: Int): Int then a
				end
				class Std::Nil
					def foo(a: Int, **b: Float): Int then a
				end
				var a: Foo | Bar | Nil = nil
				a.foo(5, a: 2.5)
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(220, 12, 5), P(235, 12, 20)), "method `Std::Nil.:foo` has a required parameter missing in `Foo.:foo`, got `b`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
