package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestIsA(t *testing.T) {
	tests := testTable{
		"impossible check": {
			input: `
				1.2 <: Int
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(5, 2, 5), P(7, 2, 7)), "impossible \"is a\" check, `Std::Float(1.2)` cannot ever be an instance of a descendant of `Std::Int`"),
			},
		},
		"always true check": {
			input: `
				1 <: Int
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(5, 2, 5), P(5, 2, 5)), "this \"is a\" check is always true, `Std::Int(1)` will always be an instance of `Std::Int`"),
			},
		},
		"valid check with class": {
			input: `
				var a: String | Int = 1
				a <: Int
			`,
		},
		"valid check with subclass": {
			input: `
				class Foo; end
				class Bar < Foo; end
				var a: Bar? = nil
				a <: Foo
			`,
		},
		"valid nested check with subclass": {
			input: `
				class Foo; end
				class Bar < Foo; end
				var a: String | Bar? = nil
				a <: Foo
			`,
		},
		"valid check with mixin": {
			input: `
				mixin Foo; end
				class Bar
					include Foo
				end
				var a: Bar? = nil
				a <: Foo
			`,
		},
		"valid nested check with mixin": {
			input: `
				mixin Foo; end
				class Bar
					include Foo
				end
				var a: String | Bar? = nil
				a <: Foo
			`,
		},
		"invalid right operand": {
			input: `
				1 <: 5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(10, 2, 10), P(10, 2, 10)), "only classes and mixins are allowed as the right operand of the is a operator `<:`"),
			},
		},
		"invalid right operand - module": {
			input: `
				module Foo; end
				1 <: Foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 3, 10), P(32, 3, 12)), "only classes and mixins are allowed as the right operand of the is a operator `<:`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestInstanceOf(t *testing.T) {
	tests := testTable{
		"impossible check": {
			input: `
				1.2 <<: Int
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(5, 2, 5), P(7, 2, 7)), "impossible \"instance of\" check, `Std::Float(1.2)` cannot ever be an instance of `Std::Int`"),
			},
		},
		"always true check": {
			input: `
				1 <<: Int
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(5, 2, 5), P(5, 2, 5)), "this \"instance of\" check is always true, `Std::Int(1)` will always be an instance of `Std::Int`"),
			},
		},
		"valid check with class": {
			input: `
				var a: String | Int = 1
				a <<: Int
			`,
		},
		"impossible check with subclass": {
			input: `
				class Foo; end
				class Bar < Foo; end
				var a: Bar? = nil
				a <<: Foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(71, 5, 5), P(71, 5, 5)), "impossible \"instance of\" check, `Bar?` cannot ever be an instance of `Foo`"),
			},
		},
		"invalid right operand": {
			input: `
				1 <<: 5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 2, 11), P(11, 2, 11)), "only classes are allowed as the right operand of the instance of operator `<<:`"),
			},
		},
		"invalid right operand - module": {
			input: `
				module Foo; end
				1 <<: Foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(31, 3, 11), P(33, 3, 13)), "only classes are allowed as the right operand of the instance of operator `<<:`"),
			},
		},
		"invalid right operand - mixin": {
			input: `
				mixin Foo; end
				1 <<: Foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 3, 11), P(32, 3, 13)), "only classes are allowed as the right operand of the instance of operator `<<:`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestBinaryOpMethod(t *testing.T) {
	tests := testTable{
		"Call custom add method": {
			input: `
				class Foo
					def +(other: String): String
						other
					end
				end

				var a: String = Foo() + "lol"
			`,
		},
		"Call add method on a type without it": {
			input: `
				class Foo; end

				var a: String = Foo() + "lol"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(41, 4, 21), P(53, 4, 33)), "method `+` is not defined on type `Foo`"),
				error.NewFailure(L("<main>", P(41, 4, 21), P(53, 4, 33)), "type `void` cannot be assigned to type `Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestAdd(t *testing.T) {
	tests := testTable{
		"Int - String => Error": {
			input: `
				var a: Int = 1 + "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 2, 22), P(26, 2, 26)), "expected type `Std::Int` for parameter `other` in call to `+`, got type `Std::String(\"foo\")`"),
			},
		},
		"Int + Int => Int": {
			input: `
				var a: Int = 1 + 2
			`,
		},
		"Int + Float => Float": {
			input: `
				var a: Float = 1 + 2.2
			`,
		},
		"Int + BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1 + 2.2bf
			`,
		},
		"Float + Int => Float": {
			input: `
				var a: Float = 1.2 + 2
			`,
		},
		"Float + Float => Float": {
			input: `
				var a: Float = 1.2 + 2.9
			`,
		},
		"Float + BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1.2 + 2.9bf
			`,
		},
		"BigFloat + Int => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf + 2
			`,
		},
		"BigFloat + Float => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf + .2
			`,
		},
		"BigFloat + BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf + .2bf
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := testTable{
		"Int - String => Error": {
			input: `
				var a: Int = 1 - "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 2, 22), P(26, 2, 26)), "expected type `Std::Int` for parameter `other` in call to `-`, got type `Std::String(\"foo\")`"),
			},
		},
		"Int - Int => Int": {
			input: `
				var a: Int = 1 - 2
			`,
		},
		"Int - Float => Float": {
			input: `
				var a: Float = 1 - 2.2
			`,
		},
		"Int - BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1 - 2.2bf
			`,
		},
		"Float - Int => Float": {
			input: `
				var a: Float = 1.2 - 2
			`,
		},
		"Float - Float => Float": {
			input: `
				var a: Float = 1.2 - 2.9
			`,
		},
		"Float - BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1.2 - 2.9bf
			`,
		},
		"BigFloat - Int => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf - 2
			`,
		},
		"BigFloat - Float => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf - .2
			`,
		},
		"BigFloat - BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf - .2bf
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := testTable{
		"Int * String => Error": {
			input: `
				var a: Int = 1 * "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 2, 22), P(26, 2, 26)), "expected type `Std::Int` for parameter `other` in call to `*`, got type `Std::String(\"foo\")`"),
			},
		},
		"Int * Int => Int": {
			input: `
				var a: Int = 1 * 2
			`,
		},
		"Int * Float => Float": {
			input: `
				var a: Float = 1 * 2.2
			`,
		},
		"Int * BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1 * 2.2bf
			`,
		},
		"Float * Int => Float": {
			input: `
				var a: Float = 1.2 * 2
			`,
		},
		"Float * Float => Float": {
			input: `
				var a: Float = 1.2 * 2.9
			`,
		},
		"Float * BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1.2 * 2.9bf
			`,
		},
		"BigFloat * Int => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf * 2
			`,
		},
		"BigFloat * Float => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf * .2
			`,
		},
		"BigFloat * BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf * .2bf
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestDivide(t *testing.T) {
	tests := testTable{
		"Int / String => Error": {
			input: `
				var a: Int = 1 / "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 2, 22), P(26, 2, 26)), "expected type `Std::Int` for parameter `other` in call to `/`, got type `Std::String(\"foo\")`"),
			},
		},
		"Int / Int => Int": {
			input: `
				var a: Int = 1 / 2
			`,
		},
		"Int / Float => Float": {
			input: `
				var a: Float = 1 / 2.2
			`,
		},
		"Int / BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1 / 2.2bf
			`,
		},
		"Float / Int => Float": {
			input: `
				var a: Float = 1.2 / 2
			`,
		},
		"Float / Float => Float": {
			input: `
				var a: Float = 1.2 / 2.9
			`,
		},
		"Float / BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1.2 / 2.9bf
			`,
		},
		"BigFloat / Int => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf / 2
			`,
		},
		"BigFloat / Float => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf / .2
			`,
		},
		"BigFloat / BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf / .2bf
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestExponentiate(t *testing.T) {
	tests := testTable{
		"Int ** String => Error": {
			input: `
				var a: Int = 1 ** "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(23, 2, 23), P(27, 2, 27)), "expected type `Std::Int` for parameter `other` in call to `**`, got type `Std::String(\"foo\")`"),
			},
		},
		"Int ** Int => Int": {
			input: `
				var a: Int = 1 ** 2
			`,
		},
		"Int ** Float => Float": {
			input: `
				var a: Float = 1 ** 2.2
			`,
		},
		"Int ** BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1 ** 2.2bf
			`,
		},
		"Float ** Int => Float": {
			input: `
				var a: Float = 1.2 ** 2
			`,
		},
		"Float ** Float => Float": {
			input: `
				var a: Float = 1.2 ** 2.9
			`,
		},
		"Float ** BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1.2 ** 2.9bf
			`,
		},
		"BigFloat ** Int => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf ** 2
			`,
		},
		"BigFloat ** Float => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf ** .2
			`,
		},
		"BigFloat ** BigFloat => BigFloat": {
			input: `
				var a: BigFloat = 1.2bf ** .2bf
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
