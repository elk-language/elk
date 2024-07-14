package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

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
