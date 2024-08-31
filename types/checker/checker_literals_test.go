package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestArrayTupleLiteral(t *testing.T) {
	tests := testTable{
		"infer array tuple": {
			input: `
				var foo = %[1, 2]
				var a: ArrayTuple[Int] = foo
			`,
		},
		"infer array tuple with different argument types": {
			input: `
				var a = %[1, 2.2, "foo"]
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 16), P(45, 3, 16)), "type `Std::ArrayTuple[Std::Int | Std::Float | Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer empty array tuple": {
			input: `
				var foo = %[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 16), P(36, 3, 18)), "type `Std::ArrayTuple[never]` cannot be assigned to type `9`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestHashSetLiteral(t *testing.T) {
	tests := testTable{
		"infer array list": {
			input: `
				var foo = ^[1, 2]
				var a: HashSet[Int] = foo
			`,
		},
		"infer array list with different argument types": {
			input: `
				var a = ^[1, 2.2, "foo"]
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 16), P(45, 3, 16)), "type `Std::HashSet[Std::Int | Std::Float | Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer empty array list": {
			input: `
				var foo = ^[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 16), P(36, 3, 18)), "type `Std::HashSet[never]` cannot be assigned to type `9`"),
			},
		},
		"int capacity": {
			input: `
				var foo: HashSet[Float] = ^[1.2]:9
			`,
		},
		"uint8 capacity": {
			input: `
				var foo: HashSet[Float] = ^[1.2]:9u8
			`,
		},
		"invalid capacity": {
			input: `
				var foo: HashSet[Float] = ^[1.2]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(31, 2, 31), P(40, 2, 40)), "capacity must be an integer, got `9.2`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestArrayListLiteral(t *testing.T) {
	tests := testTable{
		"infer array list": {
			input: `
				var foo = [1, 2]
				var a: ArrayList[Int] = foo
			`,
		},
		"infer array list with different argument types": {
			input: `
				var a = [1, 2.2, "foo"]
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 16), P(44, 3, 16)), "type `Std::ArrayList[Std::Int | Std::Float | Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer empty array list": {
			input: `
				var foo = []
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 3, 16), P(35, 3, 18)), "type `Std::ArrayList[never]` cannot be assigned to type `9`"),
			},
		},
		"int capacity": {
			input: `
				var foo: ArrayList[Float] = [1.2]:9
			`,
		},
		"uint8 capacity": {
			input: `
				var foo: ArrayList[Float] = [1.2]:9u8
			`,
		},
		"invalid capacity": {
			input: `
				var foo: ArrayList[Float] = [1.2]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 2, 33), P(41, 2, 41)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer word array list": {
			input: `
				var foo = \w[foo bar]
				var a: ArrayList[String] = foo
			`,
		},
		"infer empty word array list": {
			input: `
				var foo = \w[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayList[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"word array list int capacity": {
			input: `
				var foo: ArrayList[String] = \w[1.2]:9
			`,
		},
		"word array list uint8 capacity": {
			input: `
				var foo: ArrayList[String] = \w[1.2]:9u8
			`,
		},
		"word array list invalid capacity": {
			input: `
				var foo: ArrayList[String] = \w[1.2]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 2, 42), P(44, 2, 44)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer symbol array list": {
			input: `
				var foo = \s[foo bar]
				var a: ArrayList[Symbol] = foo
			`,
		},
		"infer empty symbol array list": {
			input: `
				var foo = \s[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayList[Std::Symbol]` cannot be assigned to type `9`"),
			},
		},
		"symbol array list int capacity": {
			input: `
				var foo: ArrayList[Symbol] = \s[1.2]:9
			`,
		},
		"symbol array list uint8 capacity": {
			input: `
				var foo: ArrayList[Symbol] = \s[1.2]:9u8
			`,
		},
		"symbol array list invalid capacity": {
			input: `
				var foo: ArrayList[Symbol] = \s[1.2]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 2, 42), P(44, 2, 44)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer hex array list": {
			input: `
				var foo = \x[1f4 fff]
				var a: ArrayList[Int] = foo
			`,
		},
		"infer empty hex array list": {
			input: `
				var foo = \x[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayList[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"hex array list int capacity": {
			input: `
				var foo: ArrayList[Int] = \x[1f9]:9
			`,
		},
		"hex array list uint8 capacity": {
			input: `
				var foo: ArrayList[Int] = \x[fff]:9u8
			`,
		},
		"hex array list invalid capacity": {
			input: `
				var foo: ArrayList[Int] = \x[1ef]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 2, 39), P(41, 2, 41)), "capacity must be an integer, got `9.2`"),
			},
		},

		"infer bin array list": {
			input: `
				var foo = \b[111 100]
				var a: ArrayList[Int] = foo
			`,
		},
		"infer empty bin array list": {
			input: `
				var foo = \b[]
				var a: 9 = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 3, 16), P(37, 3, 18)), "type `Std::ArrayList[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"bin array list int capacity": {
			input: `
				var foo: ArrayList[Int] = \b[100]:9
			`,
		},
		"bin array list uint8 capacity": {
			input: `
				var foo: ArrayList[Int] = \x[111]:9u8
			`,
		},
		"bin array list invalid capacity": {
			input: `
				var foo: ArrayList[Int] = \b[101]:9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 2, 39), P(41, 2, 41)), "capacity must be an integer, got `9.2`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestRegexLiteral(t *testing.T) {
	tests := testTable{
		"infer regex": {
			input: `
				var foo = %/str/
				var a: Regex = foo
			`,
		},
		"infer interpolated regex": {
			input: `
				var foo = %/${1} str/
				var b: Regex = foo
			`,
		},
		"interpolate unconvertible value": {
			input: `
				class Foo; end
				var foo = %/${Foo()} str/
				var b: Regex = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(38, 3, 19), P(42, 3, 23)), "type `Foo` does not implement interface `Std::StringConvertible`:\n\n  - missing method `Std::StringConvertible.:to_string` with signature: `def to_string(): Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestStringLiteral(t *testing.T) {
	tests := testTable{
		"infer raw string": {
			input: `
				var foo = 'str'
				var a: String = foo
			`,
		},
		"assign string literal to String": {
			input: "var foo: String = 'str'",
		},
		"assign string literal to matching literal type": {
			input: "var foo: 'str' = 'str'",
		},
		"assign string literal to non matching literal type": {
			input: "var foo: 'str' = 'foo'",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `\"foo\"` cannot be assigned to type `\"str\"`"),
			},
		},
		"infer double quoted string": {
			input: `
				var foo = "str"
				var b: String = foo
			`,
		},
		"infer interpolated string": {
			input: `
				var foo = "${1} str #{5.2}"
				var b: String = foo
			`,
		},
		"interpolate unconvertible value": {
			input: `
				class Foo; end
				var foo = "${Foo()} str"
				var b: String = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(37, 3, 18), P(41, 3, 22)), "type `Foo` does not implement interface `Std::StringConvertible`:\n\n  - missing method `Std::StringConvertible.:to_string` with signature: `def to_string(): Std::String`"),
			},
		},
		"interpolate uninspectable value": {
			input: `
				class Foo < nil; end
				var foo = "#{Foo()} str"
				var b: String = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 3, 18), P(47, 3, 22)), "type `Foo` does not implement interface `Std::Inspectable`:\n\n  - missing method `Std::Inspectable.:inspect` with signature: `def inspect(): Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestSymbolLiteral(t *testing.T) {
	tests := testTable{
		"infer simple symbol": {
			input: "var foo = :str",
		},
		"infer double quoted symbol": {
			input: `var foo = :"str"`,
		},
		"infer interpolated symbol": {
			input: `var foo = :"${1} str #{5.2}"`,
		},
		"assign symbol literal to Symbol": {
			input: "var foo: Symbol = :symb",
		},
		"assign symbol literal to matching literal type": {
			input: "var foo: :symb = :symb",
		},
		"assign symbol literal to non matching literal type": {
			input: "var foo: :symb = :foob",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `:foob` cannot be assigned to type `:symb`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestCharLiteral(t *testing.T) {
	tests := testTable{
		"infer raw char": {
			input: "var foo = r`s`",
		},
		"infer char": {
			input: "var foo = `\\n`",
		},
		"assign char literal to Char": {
			input: "var foo: Char = `f`",
		},
		"assign char literal to matching literal type": {
			input: "var foo: `f` = `f`",
		},
		"assign char literal to non matching literal type": {
			input: "var foo: `b` = `f`",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(15, 1, 16), P(17, 1, 18)), "type ``f`` cannot be assigned to type ``b``"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestIntLiteral(t *testing.T) {
	tests := testTable{
		"infer int": {
			input: "var foo = 1234",
		},
		"assign int literal to Int": {
			input: "var foo: Int = 12345678",
		},
		"assign int literal to matching literal type": {
			input: "var foo: 12345 = 12345",
		},
		"assign int literal to non matching literal type": {
			input: "var foo: 23456 = 12345",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `12345` cannot be assigned to type `23456`"),
			},
		},
		"infer int64": {
			input: "var foo = 1i64",
		},
		"infer int32": {
			input: "var foo = 1i32",
		},
		"infer int16": {
			input: "var foo = 1i16",
		},
		"infer int8": {
			input: "var foo = 12i8",
		},
		"infer uint64": {
			input: "var foo = 1u64",
		},
		"infer uint32": {
			input: "var foo = 1u32",
		},
		"infer uint16": {
			input: "var foo = 1u16",
		},
		"infer uint8": {
			input: "var foo = 12u8",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestFloatLiteral(t *testing.T) {
	tests := testTable{
		"infer float": {
			input: "var foo = 12.5",
		},
		"assign float literal to Float": {
			input: "var foo: Float = 1234.6",
		},
		"assign float literal to matching literal type": {
			input: "var foo: 12.45 = 12.45",
		},
		"assign Float literal to non matching literal type": {
			input: "var foo: 23.56 = 12.45",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `12.45` cannot be assigned to type `23.56`"),
			},
		},
		"infer float64": {
			input: "var foo = 1f64",
		},
		"infer float32": {
			input: "var foo = 1f32",
		},
		"infer big float": {
			input: "var foo = 12bf",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestBoolLiteral(t *testing.T) {
	tests := testTable{
		"infer true": {
			input: "var foo = true",
		},
		"assign true literal to True": {
			input: "var foo: True = true",
		},
		"assign true literal to Bool": {
			input: "var foo: Bool = true",
		},
		"assign true literal to matching literal type": {
			input: "var foo: true = true",
		},
		"assign true literal to non matching literal type": {
			input: "var foo: false = true",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 1, 18), P(20, 1, 21)), "type `true` cannot be assigned to type `false`"),
			},
		},
		"infer false": {
			input: "var foo = false",
		},
		"assign false literal to False": {
			input: "var foo: False = false",
		},
		"assign false literal to Bool": {
			input: "var foo: Bool = false",
		},
		"assign false literal to matching literal type": {
			input: "var foo: false = false",
		},
		"assign false literal to non matching literal type": {
			input: "var foo: true = false",
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 1, 17), P(20, 1, 21)), "type `false` cannot be assigned to type `true`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestNilLiteral(t *testing.T) {
	tests := testTable{
		"infer nil": {
			input: "var foo = nil",
		},
		"assign nil literal to Nil": {
			input: "var foo: Nil = nil",
		},
		"assign nil literal to matching literal type": {
			input: "var foo: nil = nil",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
