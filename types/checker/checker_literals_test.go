package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

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
				error.NewFailure(L("<main>", P(38, 3, 19), P(42, 3, 23)), "type `Foo` does not implement interface `Std::StringConvertible`:\n\n  - missing method `Std::StringConvertible.:to_string` with signature: `def to_string(): Std::String`\n"),
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
				error.NewFailure(L("<main>", P(37, 3, 18), P(41, 3, 22)), "type `Foo` does not implement interface `Std::StringConvertible`:\n\n  - missing method `Std::StringConvertible.:to_string` with signature: `def to_string(): Std::String`\n"),
			},
		},
		"interpolate uninspectable value": {
			input: `
				class Foo < nil; end
				var foo = "#{Foo()} str"
				var b: String = foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 3, 18), P(47, 3, 22)), "type `Foo` does not implement interface `Std::Inspectable`:\n\n  - missing method `Std::Inspectable.:inspect` with signature: `def inspect(): Std::String`\n"),
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
