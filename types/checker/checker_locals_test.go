package checker

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
		"reject redeclared variable": {
			input: `var foo: Int; foo := "foo"`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(14, 1, 15), P(25, 1, 26)), "cannot redeclare local `foo`"),
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
		"accept value declaration without initializer": {
			input: "val foo: Int",
		},
		"reject value declaration with invalid type": {
			input: "val foo: Foo",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
			},
		},
		"reject value declaration without initializer and type": {
			input: "val foo",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(6, 1, 7)), "cannot declare a local without a type `foo`"),
			},
		},
		"reject redeclared value": {
			input: "val foo: Int; val foo: String",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(14, 1, 15), P(28, 1, 29)), "cannot redeclare local `foo`"),
			},
		},
		"declaration with type lookup": {
			input: "val foo: Std::Int",
		},
		"declaration with type lookup and error in the middle": {
			input: "val foo: Std::Foo::Bar",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(14, 1, 15), P(16, 1, 17)), "undefined type `Std::Foo`"),
			},
		},
		"declaration with type lookup and error at the start": {
			input: "val foo: Foo::Bar::Baz",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
			},
		},
		"declaration with absolute type lookup": {
			input: "val foo: ::Std::Int",
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
				diagnostic.NewFailure(L("<main>", P(14, 1, 15), P(16, 1, 17)), "cannot access uninitialised local `foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
