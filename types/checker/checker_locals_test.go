package checker_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

func TestVariableAssignment(t *testing.T) {
	tests := testTable{
		"assign uninitialised variable with a matching type": {
			input: `
				var foo: Int
				foo = 5
				println foo
			`,
		},
		"assign uninitialised variable with a non-matching type": {
			input: `
				var foo: Int
				foo = 'f'
				println foo
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(28, 3, 11), P(30, 3, 13)), "type `\"f\"` cannot be assigned to type `Std::Int`"),
			},
		},
		"assign initialised variable with a matching type": {
			input: `
				var foo: Int = 5
				foo = 3
			`,
		},
		"assign initialised variable with a non-matching type": {
			input: `
				var foo: Int = 5
				foo = 'f'
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(32, 3, 11), P(34, 3, 13)), "type `\"f\"` cannot be assigned to type `Std::Int`"),
			},
		},

		"??= uninitialised variable with a non-matching and non-nilable type": {
			input: `
				var foo: Int
				foo ??= 'f'
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 5), P(24, 3, 7)), "cannot access uninitialised local `foo`"),
				diagnostic.NewWarning(L("<main>", P(22, 3, 5), P(24, 3, 7)), "this condition will always have the same result since type `Std::Int` can never be nil"),
				diagnostic.NewWarning(L("<main>", P(30, 3, 13), P(32, 3, 15)), "unreachable code"),
			},
		},
		"??= initialised variable with a non-matching type": {
			input: `
				var foo: Int? = 5
				foo ??= 'f'
				println foo
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(35, 3, 13), P(37, 3, 15)), "type `Std::Int | \"f\"` cannot be assigned to type `Std::Int?`"),
			},
		},
		"??= initialised variable with a matching nilable type": {
			input: `
				var foo: Int? = nil
				foo ??= 5
			`,
		},
		"??= initialised variable with a matching nilable union type": {
			input: `
				var foo: Int | Float | Nil = nil
				foo ??= 5
			`,
		},

		"||= uninitialised variable with a non-matching and non-falsy type": {
			input: `
				var foo: Int
				foo ||= 'f'
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 5), P(24, 3, 7)), "cannot access uninitialised local `foo`"),
				diagnostic.NewWarning(L("<main>", P(22, 3, 5), P(24, 3, 7)), "this condition will always have the same result since type `Std::Int` is truthy"),
				diagnostic.NewWarning(L("<main>", P(30, 3, 13), P(32, 3, 15)), "unreachable code"),
			},
		},
		"||= initialised variable with a non-matching and non-falsy type": {
			input: `
				var foo: Int? = 5
				foo ||= 'f'
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(35, 3, 13), P(37, 3, 15)), "type `Std::Int | \"f\"` cannot be assigned to type `Std::Int?`"),
			},
		},
		"||= initialised variable with a matching nilable type": {
			input: `
				var foo: Int? = nil
				foo ||= 5
			`,
		},
		"||= initialised variable with a matching falsy type": {
			input: `
				var foo: Int | Float | False = false
				foo ||= 5
			`,
		},

		"&&= uninitialised variable": {
			input: `
				var foo: Nil | False
				foo &&= 'f'
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(30, 3, 5), P(32, 3, 7)), "cannot access uninitialised local `foo`"),
				diagnostic.NewWarning(L("<main>", P(30, 3, 5), P(32, 3, 7)), "this condition will always have the same result since type `Std::Nil | Std::False` is falsy"),
				diagnostic.NewWarning(L("<main>", P(38, 3, 13), P(40, 3, 15)), "unreachable code"),
			},
		},
		"&&= initialised variable with a non-matching type": {
			input: `
				var foo: Int? = nil
				foo &&= 'f'
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(37, 3, 13), P(39, 3, 15)), "type `nil | \"f\"` cannot be assigned to type `Std::Int?`"),
			},
		},
		"&&= initialised variable with a matching truthy type": {
			input: `
				var foo: Int? = nil
				foo &&= 5
			`,
		},

		"+= uninitialised variable": {
			input: `
				var a: String
				a += "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(23, 3, 5), P(23, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"+= on a type with the method": {
			input: `
				var a = "foo"
				a += "bar"
			`,
		},
		"+= on a type without the method": {
			input: `
				var a = Object()
				a += "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(31, 3, 10), P(35, 3, 14)), "method `+` is not defined on type `Std::Object`"),
			},
		},

		"-= uninitialised variable": {
			input: `
				var a: Int
				a -= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"-= on a type with the method": {
			input: `
				var a = 1
				a -= 2
			`,
		},
		"-= on a type without the method": {
			input: `
				var a = Object()
				a -= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(31, 3, 10), P(35, 3, 14)), "method `-` is not defined on type `Std::Object`"),
			},
		},

		"*= uninitialised variable": {
			input: `
				var a: Int
				a *= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"*= on a type with the method": {
			input: `
				var a = 1
				a *= 2
			`,
		},
		"*= on a type without the method": {
			input: `
				var a = Object()
				a *= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(31, 3, 10), P(35, 3, 14)), "method `*` is not defined on type `Std::Object`"),
			},
		},

		"/= uninitialised variable": {
			input: `
				var a: Int
				a /= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"/= on a type with the method": {
			input: `
				var a = 1
				a /= 2
			`,
		},
		"/= on a type without the method": {
			input: `
				var a = Object()
				a /= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(31, 3, 10), P(35, 3, 14)), "method `/` is not defined on type `Std::Object`"),
			},
		},

		"**= uninitialised variable": {
			input: `
				var a: Int
				a **= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"**= on a type with the method": {
			input: `
				var a = 1
				a **= 2
			`,
		},
		"**= on a type without the method": {
			input: `
				var a = Object()
				a **= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(32, 3, 11), P(36, 3, 15)), "method `**` is not defined on type `Std::Object`"),
			},
		},

		"%= uninitialised variable": {
			input: `
				var a: Int
				a %= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"%= on a type with the method": {
			input: `
				var a = 1
				a %= 2
			`,
		},
		"%= on a type without the method": {
			input: `
				var a = Object()
				a %= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(31, 3, 10), P(35, 3, 14)), "method `%` is not defined on type `Std::Object`"),
			},
		},

		"&= uninitialised variable": {
			input: `
				var a: Int
				a &= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"&= on a type with the method": {
			input: `
				var a = 1
				a &= 2
			`,
		},
		"&= on a type without the method": {
			input: `
				var a = Object()
				a &= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(31, 3, 10), P(35, 3, 14)), "method `&` is not defined on type `Std::Object`"),
			},
		},

		"|= uninitialised variable": {
			input: `
				var a: Int
				a |= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"|= on a type with the method": {
			input: `
				var a = 1
				a |= 2
			`,
		},
		"|= on a type without the method": {
			input: `
				var a = Object()
				a |= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(31, 3, 10), P(35, 3, 14)), "method `|` is not defined on type `Std::Object`"),
			},
		},

		"^= uninitialised variable": {
			input: `
				var a: Int
				a ^= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"^= on a type with the method": {
			input: `
				var a = 1
				a ^= 2
			`,
		},
		"^= on a type without the method": {
			input: `
				var a = Object()
				a ^= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(31, 3, 10), P(35, 3, 14)), "method `^` is not defined on type `Std::Object`"),
			},
		},

		"<<= uninitialised variable": {
			input: `
				var a: Int
				a <<= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"<<= on a type with the method": {
			input: `
				var a = 1
				a <<= 2
			`,
		},
		"<<= on a type without the method": {
			input: `
				var a = Object()
				a <<= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(32, 3, 11), P(36, 3, 15)), "method `<<` is not defined on type `Std::Object`"),
			},
		},

		">>= uninitialised variable": {
			input: `
				var a: Int
				a >>= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		">>= on a type with the method": {
			input: `
				var a = 1
				a >>= 2
			`,
		},
		">>= on a type without the method": {
			input: `
				var a = Object()
				a >>= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(32, 3, 11), P(36, 3, 15)), "method `>>` is not defined on type `Std::Object`"),
			},
		},

		">>>= uninitialised variable": {
			input: `
				var a: Int64
				a >>>= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 5), P(22, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		">>>= on a type with the method": {
			input: `
				var a = 1i64
				a >>>= 2
			`,
		},
		">>>= on a type without the method": {
			input: `
				var a = Object()
				a >>>= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(33, 3, 12), P(37, 3, 16)), "method `>>>` is not defined on type `Std::Object`"),
			},
		},

		"<<<= uninitialised variable": {
			input: `
				var a: Int64
				a <<<= 3
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 5), P(22, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"<<<= on a type with the method": {
			input: `
				var a = 1i64
				a <<<= 2
			`,
		},
		"<<<= on a type without the method": {
			input: `
				var a = Object()
				a <<<= "bar"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(33, 3, 12), P(37, 3, 16)), "method `<<<` is not defined on type `Std::Object`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestVariableDeclaration(t *testing.T) {
	tests := testTable{
		"declare a recursive closure": {
			input: `
				var calc_fib = |n: Int|: Int ->
					return 1 if n < 3

					calc_fib(n - 2) + calc_fib(n - 1)
				end
			`,
		},
		"returns void when not initialised": {
			input: "var a: 9 = (var foo: Int)",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(12, 1, 13), P(23, 1, 24)), "cannot use type `void` as a value in this context"),
			},
		},
		"returns assigned value": {
			input: "var a: 9 = (var b: String? = 'foo')",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(12, 1, 13), P(33, 1, 34)), "type `\"foo\"` cannot be assigned to type `9`"),
			},
		},
		"accept variable declaration with matching initializer and type": {
			input: "var foo: Int = 5",
		},
		"accept variable declaration with inference": {
			input: "var foo = 5",
		},
		"cannot declare variable with type void": {
			input: `
				def bar; end
				var foo = bar()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(32, 3, 15), P(36, 3, 19)), "cannot use type `void` as a value in this context"),
			},
		},
		"reject variable declaration without matching initializer and type": {
			input: "var foo: Int = 5.2",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(15, 1, 16), P(17, 1, 18)), "type `5.2` cannot be assigned to type `Std::Int`"),
			},
		},
		"accept variable declaration without initializer": {
			input: "var foo: Int",
		},
		"reject variable declaration with invalid type": {
			input: "var foo: Foo",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
			},
		},
		"reject variable declaration without initializer and type": {
			input: "var foo",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(6, 1, 7)), "cannot declare a local without a type `foo`"),
			},
		},
		"reject redeclared variable": {
			input: "var foo: Int; var foo: String",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(14, 1, 15), P(28, 1, 29)), "cannot redeclare local `foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestShortVariableDeclaration(t *testing.T) {
	tests := testTable{
		"accept variable declaration with inference": {
			input: "foo := 5",
		},
		"cannot declare variable with type void": {
			input: `
				def bar; end
				foo := bar()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(29, 3, 12), P(33, 3, 16)), "cannot use type `void` as a value in this context"),
			},
		},
		"accept redeclared short declaration with matching type": {
			input: `var foo: String?; foo := "foo"`,
		},
		"reject redeclared short declaration with different type": {
			input: `var foo: Int; foo := "foo"`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(21, 1, 22), P(25, 1, 26)), "type `\"foo\"` cannot be assigned to type `Std::Int`"),
			},
		},
		"declare a recursive closure": {
			input: `
				calc_fib := |n: Int|: Int ->
					return 1 if n < 3

					calc_fib(n - 2) + calc_fib(n - 1)
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

func TestValueDeclaration(t *testing.T) {
	tests := testTable{
		"declare a recursive closure": {
			input: `
				val calc_fib = |n: Int|: Int ->
					return 1 if n < 3

					calc_fib(n - 2) + calc_fib(n - 1)
				end
			`,
		},
		"returns void when not initialised": {
			input: "var a: 9 = (val foo: Int)",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(12, 1, 13), P(23, 1, 24)), "a value must be initialised on declaration `foo`"),
				diagnostic.NewFailure(L("<main>", P(12, 1, 13), P(23, 1, 24)), "cannot use type `void` as a value in this context"),
			},
		},
		"accept value declaration with matching initializer and type": {
			input: "val foo: Int = 5",
		},
		"accept variable declaration with inference": {
			input: "val foo = 5",
		},
		"cannot declare value with type void": {
			input: `
				def bar; end
				val foo = bar()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(32, 3, 15), P(36, 3, 19)), "cannot use type `void` as a value in this context"),
			},
		},
		"reject value declaration without matching initializer and type": {
			input: "val foo: Int = 5.2",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(15, 1, 16), P(17, 1, 18)), "type `5.2` cannot be assigned to type `Std::Int`"),
			},
		},
		"reject value declaration without initializer": {
			input: "val foo: Int",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(11, 1, 12)), "a value must be initialised on declaration `foo`"),
			},
		},
		"reject value declaration with invalid type": {
			input: "val foo: Foo",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(11, 1, 12)), "a value must be initialised on declaration `foo`"),
				diagnostic.NewFailure(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
			},
		},
		"reject value declaration without initializer and type": {
			input: "val foo",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(6, 1, 7)), "a value must be initialised on declaration `foo`"),
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(6, 1, 7)), "cannot declare a local without a type `foo`"),
			},
		},
		"reject redeclared value": {
			input: "val foo: Int; val foo: String",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(11, 1, 12)), "a value must be initialised on declaration `foo`"),
				diagnostic.NewFailure(L("<main>", P(14, 1, 15), P(28, 1, 29)), "cannot redeclare local `foo`"),
				diagnostic.NewFailure(L("<main>", P(14, 1, 15), P(28, 1, 29)), "a value must be initialised on declaration `foo`"),
			},
		},
		"declaration with type lookup": {
			input: "val foo: Std::Int",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(16, 1, 17)), "a value must be initialised on declaration `foo`"),
			},
		},
		"declaration with type lookup and error in the middle": {
			input: "val foo: Std::Foo::Bar",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(21, 1, 22)), "a value must be initialised on declaration `foo`"),
				diagnostic.NewFailure(L("<main>", P(14, 1, 15), P(16, 1, 17)), "undefined type `Std::Foo`"),
			},
		},
		"declaration with type lookup and error at the start": {
			input: "val foo: Foo::Bar::Baz",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(21, 1, 22)), "a value must be initialised on declaration `foo`"),
				diagnostic.NewFailure(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
			},
		},
		"declaration with absolute type lookup": {
			input: "val foo: ::Std::Int",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(18, 1, 19)), "a value must be initialised on declaration `foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestLocalAccess(t *testing.T) {
	tests := testTable{
		"access initialised variable": {
			input: "var foo: Int = 5; foo",
		},

		"access variable initialised in if expression": {
			input: `
				var a: String
				b := true
				if b
					a = "elo"
					a
				end
`,
		},
		"access variable initialised in if else expression": {
			input: `
				var a: String
				b := true
				if b
				else
					a = "elo"
					a
				end
`,
		},
		"access variable initialised in if else expression outside": {
			input: `
				var a: String
				b := true
				if b
				else
					a = "elo"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(78, 8, 5), P(78, 8, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in exhaustive if outside": {
			input: `
				var a: String
				b := true
				if b
					a = "foo"
				else
					a = "bar"
				end
				a
`,
		},
		"access variable initialised in complex exhaustive if outside": {
			input: `
				var a: String
				b := true
				c := true
				if b
					a = "foo"
				else if c
					a = "bar"
				else
					a = "baz"
				end
				a
`,
		},
		"access variable initialised in complex if outside": {
			input: `
				var a: String
				b := true
				c := true
				if b
					a = "foo"
				else if c
					a = "bar"
				else
					println "baz"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(140, 12, 5), P(140, 12, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in if expression outside": {
			input: `
				var a: String
				b := true
				if b
					a = "elo"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(69, 7, 5), P(69, 7, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in modifier if expression outside": {
			input: `
				var a: String
				b := true
				a = "elo" if b
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(56, 5, 5), P(56, 5, 5)), "cannot access uninitialised local `a`"),
			},
		},

		"access variable initialised in unless expression": {
			input: `
				var a: String
				b := true
				unless b
					a = "elo"
					a
				end
`,
		},
		"access variable initialised in unless else expression": {
			input: `
				var a: String
				b := true
				unless b
				else
					a = "elo"
					a
				end
`,
		},
		"access variable initialised in unless else expression outside": {
			input: `
				var a: String
				b := true
				unless b
				else
					a = "elo"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(82, 8, 5), P(82, 8, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in exhaustive unless outside": {
			input: `
				var a: String
				b := true
				unless b
					a = "foo"
				else
					a = "bar"
				end
				a
`,
		},
		"access variable initialised in unless expression outside": {
			input: `
				var a: String
				b := true
				unless b
					a = "elo"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(73, 7, 5), P(73, 7, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in modifier unless expression outside": {
			input: `
				var a: String
				b := true
				a = "elo" unless b
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(60, 5, 5), P(60, 5, 5)), "cannot access uninitialised local `a`"),
			},
		},

		"access variable initialised in while": {
			input: `
				var a: String
				b := true
				while b
					a = "elo"
					a
				end
`,
		},
		"access variable initialised in while outside": {
			input: `
				var a: String
				b := true
				while b
					a = "elo"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(72, 7, 5), P(72, 7, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in modifier while outside": {
			input: `
				var a: String
				b := true
				a = "elo" while b
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(59, 5, 5), P(59, 5, 5)), "cannot access uninitialised local `a`"),
			},
		},

		"access variable initialised in until": {
			input: `
				var a: String
				b := true
				until b
					a = "elo"
					a
				end
`,
		},
		"access variable initialised in until outside": {
			input: `
				var a: String
				b := true
				until b
					a = "elo"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(72, 7, 5), P(72, 7, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in modifier until outside": {
			input: `
				var a: String
				b := true
				a = "elo" until b
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(59, 5, 5), P(59, 5, 5)), "cannot access uninitialised local `a`"),
			},
		},

		"access variable initialised in loop": {
			input: `
				var a: String
				loop
					a = "elo"
					a
				end
`,
		},
		"access variable initialised in loop outside": {
			input: `
				var a: String
				loop
					a = "elo"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(55, 6, 5), P(55, 6, 5)), "cannot access uninitialised local `a`"),
				diagnostic.NewWarning(L("<main>", P(55, 6, 5), P(55, 6, 5)), "unreachable code"),
			},
		},

		"access variable initialised in fornum outside": {
			input: `
				var a: String
				b := true
				fornum ;b;
					a = "elo"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(75, 7, 5), P(75, 7, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in fornum": {
			input: `
				var a: String
				b := true
				fornum ;b;
					a = "elo"
					a
				end
`,
		},

		"access variable initialised in for outside": {
			input: `
				var a: String
				for i in 5
					a = "elo"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(61, 6, 5), P(61, 6, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in for": {
			input: `
				var a: String
				for i in 5
					a = "elo"
					a
				end
`,
		},

		"access variable initialised in some switch cases outside": {
			input: `
				var a: String
				var b: any = nil
				switch b
				case String()
					a = "elo"
				case Int()
					println "foo"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(132, 10, 5), P(132, 10, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in all switch cases but not else outside": {
			input: `
				var a: String
				var b: any = nil
				switch b
				case String()
					a = "elo"
				case Int()
					a = "foo"
				else
					println "bar"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(156, 12, 5), P(156, 12, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in all switch cases and else outside": {
			input: `
				var a: String
				var b: any = nil
				switch b
				case String()
					a = "elo"
					a
				case Int()
					a = "foo"
					a
				else
					a = "bar"
					a
				end
				a
`,
		},
		"access variable initialised in all switch cases outside": {
			input: `
				var a: String
				var b: any = nil
				switch b
				case String()
					a = "elo"
					a
				case Int()
					a = "foo"
					a
				end
				a
`,
		},

		"access variable initialised in do": {
			input: `
				var a: String
				do
					a = "elo"
					a
				end
`,
		},
		"access variable initialised in do outside": {
			input: `
				var a: String
				do
					a = "elo"
				end
				a
`,
		},
		"access variable initialised in do catch": {
			input: `
				var a: String
				do
					a = "elo"
					a
				catch String()
					a = "foo"
					a
				end
`,
		},
		"access variable initialised in do catch outside": {
			input: `
				var a: String
				do
					a = "elo"
				catch String()
					a = "foo"
				end
				a
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(87, 8, 5), P(87, 8, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"access variable initialised in do catch finally": {
			input: `
				var a: String
				do
					a = "elo"
					a
				catch String()
					a = "foo"
					a
				finally
					a = "bar"
					a
				end
`,
		},
		"access variable initialised in do catch finally outside": {
			input: `
				var a: String
				do
					println "elo"
				catch String()
					println "foo"
				finally
					a = "bar"
				end
				a
`,
		},

		"access uninitialised variable": {
			input: "var foo: Int; foo",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(14, 1, 15), P(16, 1, 17)), "cannot access uninitialised local `foo`"),
			},
		},
		"access initialised value": {
			input: "val foo: Int = 5; foo",
		},
		"access uninitialised value": {
			input: "val foo: Int; foo",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(11, 1, 12)), "a value must be initialised on declaration `foo`"),
				diagnostic.NewFailure(L("<main>", P(14, 1, 15), P(16, 1, 17)), "cannot access uninitialised local `foo`"),
			},
		},
		"create a box to a local": {
			input: `
				a := 5
				var b: nil = &a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(29, 3, 18), P(30, 3, 19)), "type `Std::Box[Std::Int]` cannot be assigned to type `nil`"),
			},
		},
		"create a box to an immutable local": {
			input: `
				val a = 5
				var b: nil = &a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(32, 3, 18), P(33, 3, 19)), "type `Std::ImmutableBox[5]` cannot be assigned to type `nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
