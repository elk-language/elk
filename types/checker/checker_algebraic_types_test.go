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
				error.NewFailure(L("<main>", P(72, 3, 16), P(72, 3, 16)), "type `Std::Float | Std::String | Std::Int | Std::Char | nil` cannot be assigned to type `9`"),
			},
		},
		"normalise Int | ~Int": {
			input: `
				var a: Int | ~Int = 3
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 3, 16), P(42, 3, 16)), "type `any` cannot be assigned to type `9`"),
			},
		},
		"normalise ~Int | Int": {
			input: `
				var a: ~Int | Int = 3
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 3, 16), P(42, 3, 16)), "type `any` cannot be assigned to type `9`"),
			},
		},
		"normalise Bool | False": {
			input: `
				var a: Bool | False = :foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 2, 27), P(30, 2, 30)), "type `:foo` cannot be assigned to type `Std::Bool`"),
			},
		},
		"normalise False | Bool": {
			input: `
				var a: False | Bool = :foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 2, 27), P(30, 2, 30)), "type `:foo` cannot be assigned to type `Std::Bool`"),
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
				error.NewFailure(L("<main>", P(257, 11, 68), P(257, 11, 68)), "type `3` cannot be assigned to type `StringConvertible & IntConvertible & SigmaConvertible`"),
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
				error.NewFailure(L("<main>", P(27, 2, 27), P(27, 2, 27)), "type `3` cannot be assigned to type `never`"),
			},
		},
		"normalise intersection type with multiple modules": {
			input: `
				module Foo; end
				var a: Std & Foo = 3
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 24), P(44, 3, 24)), "type `3` cannot be assigned to type `never`"),
			},
		},
		"normalise intersection type with the same module repeated": {
			input: `
				var a: Std & Std = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 2, 24), P(24, 2, 24)), "type `3` cannot be assigned to type `Std`"),
			},
		},
		"normalise intersection type with the same class repeated": {
			input: `
				var a: String & String = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 2, 30), P(30, 2, 30)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},
		"normalise intersection type with the same literal repeated": {
			input: `
				var a: 9 & 9 = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 2, 20), P(20, 2, 20)), "type `3` cannot be assigned to type `9`"),
			},
		},
		"normalise intersection type with different literals": {
			input: `
				var a: 9 & 3 = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 2, 20), P(20, 2, 20)), "type `3` cannot be assigned to type `never`"),
			},
		},
		"normalise Int & ~Int": {
			input: `
				var a: Int & ~Int = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(25, 2, 25), P(25, 2, 25)), "type `3` cannot be assigned to type `never`"),
			},
		},
		"normalise Float & Int & ~Int": {
			input: `
				var a: Float & Int & ~Int = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 2, 33), P(33, 2, 33)), "type `3` cannot be assigned to type `never`"),
			},
		},
		"normalise (Float | Int) & ~Int": {
			input: `
				var a: (Float | Int) & ~Int = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 2, 35), P(35, 2, 35)), "type `3` cannot be assigned to type `Std::Float`"),
			},
		},
		"normalise intersection of unions": {
			input: `
				var a: (1 | 2) & (2 | 3) = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(32, 2, 32), P(32, 2, 32)), "type `3` cannot be assigned to type `2`"),
			},
		},
		"normalise intersection of union and negation": {
			input: `
				var a: (String | Float | Int) & ~Float = 2.5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(46, 2, 46), P(48, 2, 48)), "type `2.5` cannot be assigned to type `Std::String | Std::Int`"),
			},
		},
		"normalise Float & 9.2": {
			input: `
				var a: Float & 9.2 = 2.5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(26, 2, 26), P(28, 2, 28)), "type `2.5` cannot be assigned to type `9.2`"),
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

func TestNotType(t *testing.T) {
	tests := testTable{
		"assign Int to not Int": {
			input: `
				var a = 3
				var b: ~Int = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 3, 19), P(33, 3, 19)), "type `Std::Int` cannot be assigned to type `~Std::Int`"),
			},
		},
		"assign any non Int value to not Int": {
			input: `
				var b: ~Int = "foo"
				b = 2.5
				b = 9u8
			`,
		},
		"cannot assign a superclass of the negated type to the not type": {
			input: `
				var a: ~Object = Value()
				var b: ~String = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 2, 22), P(28, 2, 28)), "type `Std::Value` cannot be assigned to type `~Std::Object`"),
			},
		},
		"assign ~String to ~Object": {
			input: `
				var a: ~String = 5
				var b: ~Object = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 22), P(45, 3, 22)), "type `~Std::String` cannot be assigned to type `~Std::Object`"),
			},
		},
		"normalise nested not types": {
			input: `
				var a: ~(~Int) = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 2, 22), P(26, 2, 26)), "type `\"foo\"` cannot be assigned to type `Std::Int`"),
			},
		},
		"normalise not any to never": {
			input: `
				var a: ~(any) = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(21, 2, 21), P(25, 2, 25)), "type `\"foo\"` cannot be assigned to type `never`"),
			},
		},
		"normalise not never to any": {
			input: `
				var a: ~(never) = "foo"
				a.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 3, 5), P(37, 3, 9)), "method `foo` is not defined on type `any`"),
			},
		},
		"normalise ~Float & String": {
			input: `
				var a: ~Float & String = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 2, 30), P(32, 2, 32)), "type `1.2` cannot be assigned to type `Std::String`"),
			},
		},
		"normalise String & ~Float": {
			input: `
				var a: String & ~Float = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 2, 30), P(32, 2, 32)), "type `1.2` cannot be assigned to type `Std::String`"),
			},
		},
		"normalise named type": {
			input: `
				interface Foo
					def foo; end
				end

				var a: AnyInt & Foo = 1
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(72, 6, 27), P(72, 6, 27)), "type `Std::Int` does not implement interface `Foo`:\n\n  - missing method `Foo.:foo` with signature: `def foo(): void`\n"),
				error.NewFailure(L("<main>", P(72, 6, 27), P(72, 6, 27)), "type `1` cannot be assigned to type `Std::AnyInt & Foo`"),
			},
		},
		"normalise Bool & ~False": {
			input: `
				var a: Bool & ~False = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 2, 28), P(30, 2, 30)), "type `1.2` cannot be assigned to type `Std::Bool & ~Std::False`"),
			},
		},
		"normalise ~Bool & False": {
			input: `
				var a: ~Bool & False = 1.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 2, 28), P(30, 2, 30)), "type `1.2` cannot be assigned to type `never`"),
			},
		},
		"normalise nil & ~false & ~nil": {
			input: `
				var a: nil & ~false & ~nil = false
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 2, 34), P(38, 2, 38)), "type `false` cannot be assigned to type `never`"),
			},
		},
		"normalise intersection of negated unions": {
			input: `
				var a: (Bool | Int | nil) & ~(false | nil) = 2.5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(50, 2, 50), P(52, 2, 52)), "type `2.5` cannot be assigned to type `(Std::Bool & ~false) | Std::Int`"),
			},
		},
		"normalise intersection with a named union": {
			input: `
				typedef Foo = Bool | Int | nil
				typedef Bar = false | nil
				var a: Foo & ~Bar = 2.5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(90, 4, 25), P(92, 4, 27)), "type `2.5` cannot be assigned to type `(Std::Bool & ~Bar) | Std::Int`"),
			},
		},
		"normalise intersection with a named union containing bool": {
			input: `
				typedef Foo = bool | Int | nil
				typedef Bar = false | nil
				var a: Foo & ~Bar = 2.5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(90, 4, 25), P(92, 4, 27)), "type `2.5` cannot be assigned to type `true | Std::Int`"),
			},
		},
		"normalise intersection with a named intersection": {
			input: `
				interface Foo
					def foo; end
				end
				typedef Bar = Int & Foo
				var a: Bar & ~Foo = 2.5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(97, 6, 25), P(99, 6, 27)), "type `2.5` cannot be assigned to type `Std::Int`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestDifferenceType(t *testing.T) {
	tests := testTable{
		"normalise Int / Int to never": {
			input: `
				var a: Int / Int = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 2, 24), P(28, 2, 28)), "type `\"foo\"` cannot be assigned to type `never`"),
			},
		},
		"normalise Int? / nil to Int": {
			input: `
				var a: Int? / nil = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(25, 2, 25), P(29, 2, 29)), "type `\"foo\"` cannot be assigned to type `Std::Int`"),
			},
		},
		"normalise with union": {
			input: `
				var a: (Int | String | Float) / String = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(46, 2, 46), P(50, 2, 50)), "type `\"foo\"` cannot be assigned to type `Std::Int | Std::Float`"),
			},
		},
		"normalise with intersection": {
			input: `
				interface Foo
					sig foo
				end
				var a: (Int & Foo) / Foo = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(71, 5, 32), P(75, 5, 36)), "type `\"foo\"` cannot be assigned to type `never`"),
			},
		},
		"normalise with named types": {
			input: `
				typedef Foo = 0 | 1 | 2 | 3 | 4 | 5
				typedef Bar = 0 | 2 | 4
				var a: Foo / Bar = 2.5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(92, 4, 24), P(94, 4, 26)), "type `2.5` cannot be assigned to type `1 | 3 | 5`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
