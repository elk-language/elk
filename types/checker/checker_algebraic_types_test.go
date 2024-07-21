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
				error.NewFailure(L("<main>", P(36, 3, 22), P(36, 3, 22)), "type `Std::Int` cannot be assigned to type `Std::String?`"),
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
				error.NewFailure(L("<main>", P(56, 3, 29), P(56, 3, 29)), "type `Std::String?` cannot be assigned to type `Std::String | Std::Float`"),
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
				error.NewFailure(L("<main>", P(69, 6, 5), P(73, 6, 9)), "method `foo` is not defined on type `Std::Nil`"),
			},
		},
		"missing method on nilable type": {
			input: `
			  class Foo; end
				sealed primitive class Std::Nil
				  def foo; end
				end
				var a: Foo? = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(110, 7, 5), P(114, 7, 9)), "method `foo` is not defined on type `Foo`"),
			},
		},
		"missing method on both types": {
			input: `
			  class Foo; end
				var a: Foo? = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(47, 4, 5), P(51, 4, 9)), "method `foo` is not defined on type `Std::Nil`"),
				error.NewFailure(L("<main>", P(47, 4, 5), P(51, 4, 9)), "method `foo` is not defined on type `Foo`"),
			},
		},
		"method with different number of arguments": {
			input: `
				class Foo
					def foo(a: Int, b: String); end
				end
				sealed primitive class Std::Nil
					def foo(a: Int); end
				end
				var a: Foo? = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(156, 9, 5), P(160, 9, 9)), "method `Foo.:foo` has a required parameter missing in `Std::Nil.:foo`, got `b`"),
			},
		},
		"method with different return types": {
			input: `
				class Foo
					def foo(a: Int): Int then a
				end
				sealed primitive class Std::Nil
					def foo(a: Int): Nil then nil
				end
				var a: Foo? = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(161, 9, 5), P(165, 9, 9)), "method `Std::Nil.:foo` has a different return type than `Foo.:foo`, has `Std::Nil`, should have `Std::Int`"),
			},
		},
		"method with different param types": {
			input: `
				class Foo
					def foo(a: Int): Int then a
				end
				sealed primitive class Std::Nil
					def foo(a: Float): Int then 5
				end
				var a: Foo? = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(161, 9, 5), P(165, 9, 9)), "method `Std::Nil.:foo` has a different type for parameter `a` than `Foo.:foo`, has `Std::Float`, should have `Std::Int`"),
			},
		},
		"method with additional optional params": {
			input: `
				class Foo
					def foo(a: Int): Int then a
				end
				sealed primitive class Std::Nil
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
				sealed primitive class Std::Nil
					def foo(a: Int, b: Float = .2): Int then a
				end
				var a: Foo? = nil
				a.foo(5, 2.5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(174, 9, 5), P(186, 9, 17)), "expected 1 arguments in call to `foo`, got 2"),
			},
		},
		"method with additional rest param": {
			input: `
				class Foo
					def foo(a: Int): Int then a
				end
				sealed primitive class Std::Nil
					def foo(a: Int, *b: Float): Int then a
				end
				var a: Foo? = nil
				a.foo(5, 2.5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(170, 9, 5), P(182, 9, 17)), "method `Std::Nil.:foo` has a required parameter missing in `Foo.:foo`, got `b`"),
			},
		},
		"method with additional named rest param": {
			input: `
				class Foo
					def foo(a: Int): Int then a
				end
				sealed primitive class Std::Nil
					def foo(a: Int, **b: Float): Int then a
				end
				var a: Foo? = nil
				a.foo(5, a: 2.5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(171, 9, 5), P(186, 9, 20)), "method `Std::Nil.:foo` has a required parameter missing in `Foo.:foo`, got `b`"),
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
				error.NewFailure(L("<main>", P(43, 3, 29), P(43, 3, 29)), "type `Std::Int` cannot be assigned to type `Std::String | Std::Float`"),
			},
		},
		"assign union type to more wide union type": {
			input: `
				var a: String | Int = 3
				var b: Int | Float | String = a
			`,
		},
		"assign union type to more narrow union type": {
			input: `
				var a: Int | Float | String = 3
				var b: Int | String = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(63, 3, 27), P(63, 3, 27)), "type `Std::Int | Std::Float | Std::String` cannot be assigned to type `Std::Int | Std::String`"),
			},
		},
		"normalise union type": {
			input: `
				var a: (String | (Int | Char | Float))? | Float = 3
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(72, 3, 16), P(72, 3, 16)), "type `Std::String | Std::Int | Std::Char | Std::Float | nil` cannot be assigned to type `Std::Int(9)`"),
			},
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
				error.NewFailure(L("<main>", P(102, 8, 5), P(106, 8, 9)), "method `foo` is not defined on type `Bar`"),
				error.NewFailure(L("<main>", P(102, 8, 5), P(106, 8, 9)), "method `foo` is not defined on type `Std::Nil`"),
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
				sealed primitive class Std::Nil
					def foo(a: Int); end
				end
				var a: Foo | Bar | Nil = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(237, 12, 5), P(241, 12, 9)), "method `Foo.:foo` has a required parameter missing in `Std::Nil.:foo`, got `b`"),
				error.NewFailure(L("<main>", P(237, 12, 5), P(241, 12, 9)), "method `Bar.:foo` has a required parameter missing in `Std::Nil.:foo`, got `b`"),
				error.NewFailure(L("<main>", P(237, 12, 5), P(241, 12, 9)), "method `Bar.:foo` has a required parameter missing in `Std::Nil.:foo`, got `c`"),
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
				sealed primitive class Std::Nil
					def foo(a: Int): Nil then nil
				end
				var a: Foo | Bar | Nil = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(230, 12, 5), P(234, 12, 9)), "method `Bar.:foo` has a different return type than `Foo.:foo`, has `Std::String`, should have `Std::Int`"),
				error.NewFailure(L("<main>", P(230, 12, 5), P(234, 12, 9)), "method `Std::Nil.:foo` has a different return type than `Foo.:foo`, has `Std::Nil`, should have `Std::Int`"),
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
				sealed primitive class Std::Nil
					def foo(a: Float): Int then 5
				end
				var a: Foo | Bar | Nil = nil
				a.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(230, 12, 5), P(234, 12, 9)), "method `Bar.:foo` has a different type for parameter `a` than `Foo.:foo`, has `Std::String`, should have `Std::Int`"),
				error.NewFailure(L("<main>", P(230, 12, 5), P(234, 12, 9)), "method `Std::Nil.:foo` has a different type for parameter `a` than `Foo.:foo`, has `Std::Float`, should have `Std::Int`"),
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
				error.NewFailure(L("<main>", P(150, 9, 5), P(159, 9, 14)), "method `Foo.:foo` has a different type for parameter `a` than `Bar.:foo`, has `Std::String`, should have `Std::Object`"),
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
				error.NewFailure(L("<main>", P(153, 9, 5), P(162, 9, 14)), "method `Bar.:foo` has a different return type than `Foo.:foo`, has `Std::Object`, should have `Std::Int`"),
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
				sealed primitive class Std::Nil
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
				sealed primitive class Std::Nil
					def foo(a: Int, b: Float = .2): Int then a
				end
				var a: Foo | Bar | Nil = nil
				a.foo(5, 2.5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(257, 12, 5), P(269, 12, 17)), "expected 1 arguments in call to `foo`, got 2"),
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
				sealed primitive class Std::Nil
					def foo(a: Int, *b: Float): Int then a
				end
				var a: Foo | Bar | Nil = nil
				a.foo(5, 2.5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(236, 12, 5), P(248, 12, 17)), "method `Std::Nil.:foo` has a required parameter missing in `Foo.:foo`, got `b`"),
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
				sealed primitive class Std::Nil
					def foo(a: Int, **b: Float): Int then a
				end
				var a: Foo | Bar | Nil = nil
				a.foo(5, a: 2.5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(237, 12, 5), P(252, 12, 20)), "method `Std::Nil.:foo` has a required parameter missing in `Foo.:foo`, got `b`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestIntersectionTypeSubtype(t *testing.T) {
	tests := testTable{
		"assign Int to intersection type with Int": {
			input: `
				var a = 3
				var b: Int & StringConvertible = a
			`,
		},
		"assign Int to intersection type with two compatible types": {
			input: `
				var a = 3
				var b: Object & StringConvertible = a
			`,
		},
		"assign intersection type to more wide intersection type": {
			input: `
				interface StringConvertible
					sig to_string: String
				end
				interface IntConvertible
					sig to_int: Int
				end
				interface FloatConvertible
					sig to_float: Float
				end
				var a: StringConvertible & IntConvertible = 3
				var b: StringConvertible & IntConvertible & FloatConvertible = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(307, 12, 68), P(307, 12, 68)), "type `StringConvertible & IntConvertible` cannot be assigned to type `StringConvertible & IntConvertible & FloatConvertible`"),
			},
		},
		"assign intersection type to more narrow intersection type": {
			input: `
				interface StringConvertible
					sig to_string: String
				end
				interface IntConvertible
					sig to_int: Int
				end
				interface FloatConvertible
					sig to_float: Float
				end
				var a: StringConvertible & IntConvertible & FloatConvertible = 3
				var b: StringConvertible & IntConvertible = a
			`,
		},
		"assign a value that does not implement one interface in the intersection": {
			input: `
				interface StringConvertible
					sig to_string: String
				end
				interface IntConvertible
					sig to_int: Int
				end
				interface SigmaConvertible
					sig to_sigma: Float
				end
				var a: StringConvertible & IntConvertible & SigmaConvertible = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(257, 11, 68), P(257, 11, 68)), "type `Std::Int` does not implement interface `SigmaConvertible`:\n\n  - missing method `SigmaConvertible.:to_sigma` with signature: `def to_sigma(): Std::Float`\n"),
				error.NewFailure(L("<main>", P(257, 11, 68), P(257, 11, 68)), "type `Std::Int(3)` cannot be assigned to type `StringConvertible & IntConvertible & SigmaConvertible`"),
			},
		},
		"assign a value that does not implement a few interfaces in the intersection": {
			input: `
				interface BarConvertible
					sig to_bar: String
				end
				interface FooConvertible
					sig to_foo: Int
				end
				interface SigmaConvertible
					sig to_sigma: Float
				end
				var a: FooConvertible & BarConvertible & SigmaConvertible = Object()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(248, 11, 65), P(255, 11, 72)), "type `Std::Object` does not implement interface `FooConvertible`:\n\n  - missing method `FooConvertible.:to_foo` with signature: `def to_foo(): Std::Int`\n"),
				error.NewFailure(L("<main>", P(248, 11, 65), P(255, 11, 72)), "type `Std::Object` does not implement interface `BarConvertible`:\n\n  - missing method `BarConvertible.:to_bar` with signature: `def to_bar(): Std::String`\n"),
				error.NewFailure(L("<main>", P(248, 11, 65), P(255, 11, 72)), "type `Std::Object` does not implement interface `SigmaConvertible`:\n\n  - missing method `SigmaConvertible.:to_sigma` with signature: `def to_sigma(): Std::Float`\n"),
				error.NewFailure(L("<main>", P(248, 11, 65), P(255, 11, 72)), "type `Std::Object` cannot be assigned to type `FooConvertible & BarConvertible & SigmaConvertible`"),
			},
		},
		"normalise intersection type with multiple classes": {
			input: `
				var a: String & Int = 3
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 2, 27), P(27, 2, 27)), "type `Std::Int(3)` cannot be assigned to type `never`"),
			},
		},
		"normalise intersection type with multiple modules": {
			input: `
				module Foo; end
				var a: Std & Foo = 3
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 24), P(44, 3, 24)), "type `Std::Int(3)` cannot be assigned to type `never`"),
			},
		},
		"normalise intersection type with the same module repeated": {
			input: `
				var a: Std & Std = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 2, 24), P(24, 2, 24)), "type `Std::Int(3)` cannot be assigned to type `Std`"),
			},
		},
		"normalise intersection type with the same class repeated": {
			input: `
				var a: String & String = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 2, 30), P(30, 2, 30)), "type `Std::Int(3)` cannot be assigned to type `Std::String`"),
			},
		},
		"normalise intersection type with the same literal repeated": {
			input: `
				var a: 9 & 9 = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 2, 20), P(20, 2, 20)), "type `Std::Int(3)` cannot be assigned to type `Std::Int(9)`"),
			},
		},
		"normalise intersection type with different literals": {
			input: `
				var a: 9 & 3 = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 2, 20), P(20, 2, 20)), "type `Std::Int(3)` cannot be assigned to type `never`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestIntersectionTypeMethodCall(t *testing.T) {
	tests := testTable{
		"cal method only present in one type": {
			input: `
			  interface Foo
				  def foo; end
				end
				interface Bar
					def bar; end
				end
				class Baz
					def foo; end
					def bar; end
				end

				var a: Foo & Bar = Baz()
				a.foo
				a.bar
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
