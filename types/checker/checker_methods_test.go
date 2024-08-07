package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestAttrDefinition(t *testing.T) {
	tests := testTable{
		"declare within a method": {
			input: `
				def foo
					attr foo: String?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(34, 3, 22)), "method definitions cannot appear in this context"),
			},
		},
		"declare an attr and call a getter": {
			input: `
				class Foo
					attr foo: String?
				end
				Foo().foo
			`,
		},
		"assign the return value of a getter to an incompatible type": {
			input: `
				class Foo
					attr foo: String?
				end
				var a: Int = Foo().foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(63, 5, 18), P(71, 5, 26)), "type `Std::String?` cannot be assigned to type `Std::Int`"),
			},
		},
		"use an instance variable declared by an attr": {
			input: `
				class Foo
					attr foo: String?

					def bar
						@foo
					end
				end
			`,
		},
		"assign an instance variable declared by an attr to an incompatible type": {
			input: `
				class Foo
					attr foo: String?

					def bar
						var a: Int = @foo
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(71, 6, 20), P(74, 6, 23)), "type `Std::String?` cannot be assigned to type `Std::Int`"),
			},
		},
		"redeclare an attr with the same type": {
			input: `
				class Foo
					attr foo: String?
					attr foo: String?
				end
			`,
		},
		"redeclare an attr with a different type": {
			input: `
				class Foo
					attr foo: String?
					attr foo: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(53, 4, 16), P(55, 4, 18)), "type `Std::Int` cannot be assigned to instance variable `@foo` of type `Std::String?`"),
				error.NewFailure(L("<main>", P(48, 4, 11), P(55, 4, 18)), "cannot redeclare instance variable `@foo` with a different type, is `Std::Int`, should be `Std::String?`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(48, 4, 11), P(55, 4, 18)), "type `Std::String?` cannot be assigned to type `Std::Int`"),
			},
		},
		"redeclare an instance variable using an attr with the same type": {
			input: `
				class Foo
					var @foo: String?
					attr foo: String?
				end
			`,
		},
		"redeclare an instance variable using an attr with a different type": {
			input: `
				class Foo
					var @foo: String?
					attr foo: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(53, 4, 16), P(55, 4, 18)), "type `Std::Int` cannot be assigned to instance variable `@foo` of type `Std::String?`"),
				error.NewFailure(L("<main>", P(48, 4, 11), P(55, 4, 18)), "cannot redeclare instance variable `@foo` with a different type, is `Std::Int`, should be `Std::String?`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(48, 4, 11), P(55, 4, 18)), "type `Std::String?` cannot be assigned to type `Std::Int`"),
			},
		},
		"override an attr with the same type in a child class": {
			input: `
				class Foo
					attr foo: Int
				end
				class Bar < Foo
					attr foo: Int
				end
			`,
		},
		"override an attr with a different type in a child class": {
			input: `
				class Foo
					attr foo: Int
				end
				class Bar < Foo
					attr foo: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(77, 6, 16), P(82, 6, 21)), "type `Std::String` cannot be assigned to instance variable `@foo` of type `Std::Int`"),
				error.NewFailure(L("<main>", P(72, 6, 11), P(82, 6, 21)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String`, should be `Std::Int`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(77, 6, 16), P(82, 6, 21)), "cannot override method `foo=` with invalid parameter type, is `Std::String`, should be `Std::Int`\n  previous definition found in `Foo`, with signature: `def foo=(foo: Std::Int): void`"),
				error.NewFailure(L("<main>", P(77, 6, 16), P(82, 6, 21)), "cannot override method `foo` with a different return type, is `Std::String`, should be `Std::Int`\n  previous definition found in `Foo`, with signature: `def foo(): Std::Int`"),
				error.NewFailure(L("<main>", P(72, 6, 11), P(82, 6, 21)), "type `Std::Int` cannot be assigned to type `Std::String`"),
			},
		},
		"override an instance variable using an attr with the same type in a child class": {
			input: `
				class Foo
					var @foo: Int
				end
				class Bar < Foo
					attr foo: Int
				end
			`,
		},
		"override an instance variable using an attr with a different type in a child class": {
			input: `
				class Foo
					var @foo: Int
				end
				class Bar < Foo
					attr foo: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(77, 6, 16), P(82, 6, 21)), "type `Std::String` cannot be assigned to instance variable `@foo` of type `Std::Int`"),
				error.NewFailure(L("<main>", P(72, 6, 11), P(82, 6, 21)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String`, should be `Std::Int`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(72, 6, 11), P(82, 6, 21)), "type `Std::Int` cannot be assigned to type `Std::String`"),
			},
		},
		"override a method using an attr with the same parameter type in a child class": {
			input: `
				class Foo
					def foo=(foo: String); end
				end
				class Bar < Foo
					attr foo: String
				end
			`,
		},
		"override a method using an attr with a different parameter type in a child class": {
			input: `
				class Foo
					def foo=(foo: String); end
				end
				class Bar < Foo
					attr foo: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(90, 6, 16), P(92, 6, 18)), "cannot override method `foo=` with invalid parameter type, is `Std::Int`, should be `Std::String`\n  previous definition found in `Foo`, with signature: `def foo=(foo: Std::String): void`"),
			},
		},
		"override a method using an attr with a different return type in a child class": {
			input: `
				class Foo
					def foo: String then "foo"
				end
				class Bar < Foo
					attr foo: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(90, 6, 16), P(92, 6, 18)), "cannot override method `foo` with a different return type, is `Std::Int`, should be `Std::String`\n  previous definition found in `Foo`, with signature: `def foo(): Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestGetterDefinition(t *testing.T) {
	tests := testTable{
		"declare within a method": {
			input: `
				def foo
					getter foo: String?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(36, 3, 24)), "method definitions cannot appear in this context"),
			},
		},
		"declare a getter and call it": {
			input: `
				class Foo
					getter foo: String?
				end
				Foo().foo
			`,
		},
		"assign the return value of a getter to an incompatible type": {
			input: `
				class Foo
					getter foo: String?
				end
				var a: Int = Foo().foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(65, 5, 18), P(73, 5, 26)), "type `Std::String?` cannot be assigned to type `Std::Int`"),
			},
		},
		"use an instance variable declared by a getter": {
			input: `
				class Foo
					getter foo: String?

					def bar
						@foo
					end
				end
			`,
		},
		"assign an instance variable declared by a getter to an incompatible type": {
			input: `
				class Foo
					getter foo: String?

					def bar
						var a: Int = @foo
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(73, 6, 20), P(76, 6, 23)), "type `Std::String?` cannot be assigned to type `Std::Int`"),
			},
		},
		"redeclare a getter with the same type": {
			input: `
				class Foo
					getter foo: String?
					getter foo: String?
				end
			`,
		},
		"redeclare a getter with a different type": {
			input: `
				class Foo
					getter foo: String?
					getter foo: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(52, 4, 13), P(59, 4, 20)), "cannot redeclare instance variable `@foo` with a different type, is `Std::Int`, should be `Std::String?`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(52, 4, 13), P(59, 4, 20)), "type `Std::String?` cannot be assigned to type `Std::Int`"),
			},
		},
		"redeclare an instance variable using a getter with the same type": {
			input: `
				class Foo
					var @foo: String?
					getter foo: String?
				end
			`,
		},
		"redeclare an instance variable using a getter with a different type": {
			input: `
				class Foo
					var @foo: String?
					getter foo: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(50, 4, 13), P(57, 4, 20)), "cannot redeclare instance variable `@foo` with a different type, is `Std::Int`, should be `Std::String?`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(50, 4, 13), P(57, 4, 20)), "type `Std::String?` cannot be assigned to type `Std::Int`"),
			},
		},
		"override a getter with the same type in a child class": {
			input: `
				class Foo
					getter foo: Int
				end
				class Bar < Foo
					getter foo: Int
				end
			`,
		},
		"override a getter with a different type in a child class": {
			input: `
				class Foo
					getter foo: Int
				end
				class Bar < Foo
					getter foo: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(76, 6, 13), P(86, 6, 23)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String`, should be `Std::Int`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(81, 6, 18), P(86, 6, 23)), "cannot override method `foo` with a different return type, is `Std::String`, should be `Std::Int`\n  previous definition found in `Foo`, with signature: `def foo(): Std::Int`"),
				error.NewFailure(L("<main>", P(76, 6, 13), P(86, 6, 23)), "type `Std::Int` cannot be assigned to type `Std::String`"),
			},
		},
		"override an instance variable using getter with the same type in a child class": {
			input: `
				class Foo
					var @foo: Int
				end
				class Bar < Foo
					getter foo: Int
				end
			`,
		},
		"override an instance variable using getter with a different type in a child class": {
			input: `
				class Foo
					var @foo: Int
				end
				class Bar < Foo
					getter foo: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(74, 6, 13), P(84, 6, 23)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String`, should be `Std::Int`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(74, 6, 13), P(84, 6, 23)), "type `Std::Int` cannot be assigned to type `Std::String`"),
			},
		},
		"override a method using a getter with the same return type in a child class": {
			input: `
				class Foo
					def foo: String then "foo"
				end
				class Bar < Foo
					getter foo: String
				end
			`,
		},
		"override a method using a getter with a different return type in a child class": {
			input: `
				class Foo
					def foo: String then "foo"
				end
				class Bar < Foo
					getter foo: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(92, 6, 18), P(94, 6, 20)), "cannot override method `foo` with a different return type, is `Std::Int`, should be `Std::String`\n  previous definition found in `Foo`, with signature: `def foo(): Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestSetterDefinition(t *testing.T) {
	tests := testTable{
		"declare within a method": {
			input: `
				def foo
					setter foo: String?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(36, 3, 24)), "method definitions cannot appear in this context"),
			},
		},
		"declare a setter and call a getter": {
			input: `
				class Foo
					setter foo: String?
				end
				Foo().foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(52, 5, 5), P(60, 5, 13)), "method `foo` is not defined on type `Foo`"),
			},
		},
		"use an instance variable declared by a setter": {
			input: `
				class Foo
					setter foo: String?

					def bar
						@foo
					end
				end
			`,
		},
		"assign an instance variable declared by a setter to an incompatible type": {
			input: `
				class Foo
					setter foo: String?

					def bar
						var a: Int = @foo
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(73, 6, 20), P(76, 6, 23)), "type `Std::String?` cannot be assigned to type `Std::Int`"),
			},
		},
		"redeclare a setter with the same type": {
			input: `
				class Foo
					setter foo: String?
					setter foo: String?
				end
			`,
		},
		"redeclare a setter with a different type": {
			input: `
				class Foo
					setter foo: String?
					setter foo: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 18), P(59, 4, 20)), "type `Std::Int` cannot be assigned to instance variable `@foo` of type `Std::String?`"),
				error.NewFailure(L("<main>", P(52, 4, 13), P(59, 4, 20)), "cannot redeclare instance variable `@foo` with a different type, is `Std::Int`, should be `Std::String?`, previous definition found in `Foo`"),
			},
		},
		"redeclare an instance variable using a setter with the same type": {
			input: `
				class Foo
					var @foo: String?
					setter foo: String?
				end
			`,
		},
		"redeclare an instance variable using a setter with a different type": {
			input: `
				class Foo
					var @foo: String?
					setter foo: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(55, 4, 18), P(57, 4, 20)), "type `Std::Int` cannot be assigned to instance variable `@foo` of type `Std::String?`"),
				error.NewFailure(L("<main>", P(50, 4, 13), P(57, 4, 20)), "cannot redeclare instance variable `@foo` with a different type, is `Std::Int`, should be `Std::String?`, previous definition found in `Foo`"),
			},
		},
		"override a setter with the same type in a child class": {
			input: `
				class Foo
					setter foo: Int
				end
				class Bar < Foo
					setter foo: Int
				end
			`,
		},
		"override a setter with a different type in a child class": {
			input: `
				class Foo
					setter foo: Int
				end
				class Bar < Foo
					setter foo: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(81, 6, 18), P(86, 6, 23)), "type `Std::String` cannot be assigned to instance variable `@foo` of type `Std::Int`"),
				error.NewFailure(L("<main>", P(76, 6, 13), P(86, 6, 23)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String`, should be `Std::Int`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(81, 6, 18), P(86, 6, 23)), "cannot override method `foo=` with invalid parameter type, is `Std::String`, should be `Std::Int`\n  previous definition found in `Foo`, with signature: `def foo=(foo: Std::Int): void`"),
			},
		},
		"override an instance variable using setter with the same type in a child class": {
			input: `
				class Foo
					var @foo: Int
				end
				class Bar < Foo
					setter foo: Int
				end
			`,
		},
		"override an instance variable using setter with a different type in a child class": {
			input: `
				class Foo
					var @foo: Int
				end
				class Bar < Foo
					setter foo: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(79, 6, 18), P(84, 6, 23)), "type `Std::String` cannot be assigned to instance variable `@foo` of type `Std::Int`"),
				error.NewFailure(L("<main>", P(74, 6, 13), P(84, 6, 23)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String`, should be `Std::Int`, previous definition found in `Foo`"),
			},
		},
		"override a method using a setter with the same parameter type in a child class": {
			input: `
				class Foo
					def foo=(foo: String); end
				end
				class Bar < Foo
					setter foo: String
				end
			`,
		},
		"override a method using a setter with a different parameter type in a child class": {
			input: `
				class Foo
					def foo=(foo: String); end
				end
				class Bar < Foo
					setter foo: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(92, 6, 18), P(94, 6, 20)), "cannot override method `foo=` with invalid parameter type, is `Std::Int`, should be `Std::String`\n  previous definition found in `Foo`, with signature: `def foo=(foo: Std::String): void`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestAliasDeclaration(t *testing.T) {
	tests := testTable{
		"declare within a method": {
			input: `
				def foo
					alias bar foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(30, 3, 18)), "method definitions cannot appear in this context"),
			},
		},
		"declare an alias": {
			input: `
				class Foo
					def foo: String then "foo"
					alias bar foo
				end
				var a: String = Foo().bar
			`,
		},
		"declare an alias of a nonexistent method": {
			input: `
				class Foo
					alias bar foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(26, 3, 12), P(32, 3, 18)), "method `foo` is not defined on type `Foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestMethodDefinitionOverride(t *testing.T) {
	tests := testTable{
		"invalid override": {
			input: `
				class Foo
					def baz(a: Int): Int then a
				end

				class Bar < Foo
					def baz(); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(82, 7, 6), P(95, 7, 19)), "cannot override method `baz` with a different return type, is `void`, should be `Std::Int`\n  previous definition found in `Foo`, with signature: `def baz(a: Std::Int): Std::Int`"),
				error.NewFailure(L("<main>", P(82, 7, 6), P(95, 7, 19)), "cannot override method `baz` with less parameters\n  previous definition found in `Foo`, with signature: `def baz(a: Std::Int): Std::Int`"),
			},
		},
		"invalid override in included mixin": {
			input: `
				mixin Foo
					def baz(a: Int): Int then a
				end

				class Bar
					include Foo
					def baz(); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(93, 8, 6), P(106, 8, 19)), "cannot override method `baz` with a different return type, is `void`, should be `Std::Int`\n  previous definition found in `Foo`, with signature: `def baz(a: Std::Int): Std::Int`"),
				error.NewFailure(L("<main>", P(93, 8, 6), P(106, 8, 19)), "cannot override method `baz` with less parameters\n  previous definition found in `Foo`, with signature: `def baz(a: Std::Int): Std::Int`"),
			},
		},
		"invalid override in mixin included in mixin": {
			input: `
				mixin Foo
					def baz(a: Int): Int then a
				end

				mixin Bar
					include Foo
					def baz(); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(93, 8, 6), P(106, 8, 19)), "cannot override method `baz` with a different return type, is `void`, should be `Std::Int`\n  previous definition found in `Foo`, with signature: `def baz(a: Std::Int): Std::Int`"),
				error.NewFailure(L("<main>", P(93, 8, 6), P(106, 8, 19)), "cannot override method `baz` with less parameters\n  previous definition found in `Foo`, with signature: `def baz(a: Std::Int): Std::Int`"),
			},
		},
		"override sealed method": {
			input: `
				class Bar
					sealed def baz(a: Int): Int then a
				end
				class Foo < Bar
					def baz(a: Int): Int then a
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(88, 6, 6), P(114, 6, 32)), "cannot override sealed method `baz`\n  previous definition found in `Bar`, with signature: `sealed def baz(a: Std::Int): Std::Int`"),
			},
		},
		"redeclare sealed method in the same class": {
			input: `
				class Bar
					sealed def baz(a: Int): Int then a
					def baz(a: Int): Int then a
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 4, 6), P(86, 4, 32)), "cannot override sealed method `baz`\n  previous definition found in `Bar`, with signature: `sealed def baz(a: Std::Int): Std::Int`"),
			},
		},
		"redeclare method with a new sealed modifier": {
			input: `
				class Bar
					def baz(a: Int): Int then a
					sealed def baz(a: Int): Int then a
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(53, 4, 6), P(86, 4, 39)), "cannot redeclare method `baz` with a different modifier, is `sealed`, should be `default`"),
			},
		},
		"redeclare method with a new abstract modifier": {
			input: `
				abstract class Bar
					def baz(a: Int): Int then a
					abstract def baz(a: Int): Int; end
				end
			`,
		},
		"override method with a new abstract modifier": {
			input: `
				class Foo
					def baz(a: Int): Int then a
				end

				abstract class Bar < Foo
					abstract def baz(a: Int): Int; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(91, 7, 6), P(124, 7, 39)), "cannot override method `baz` with a different modifier, is `abstract`, should be `default`\n  previous definition found in `Foo`, with signature: `def baz(a: Std::Int): Std::Int`"),
			},
		},
		"redeclare abstract method with a new sealed modifier ": {
			input: `
				abstract class Bar
					abstract def baz(a: Int): Int; end
					sealed def baz(a: Int): Int then a
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(69, 4, 6), P(102, 4, 39)), "cannot redeclare method `baz` with a different modifier, is `sealed`, should be `abstract`"),
			},
		},
		"override the method with additional optional params": {
			input: `
				class Bar
					def baz(a: Int): Int then a
				end
				class Foo < Bar
					def baz(a: Int, b: Int = 2): Int then a
				end
			`,
		},
		"override the method with different param name": {
			input: `
				class Bar
					def baz(a: Int): Int then a
				end
				class Foo < Bar
					def baz(b: Int): Int then b
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(89, 6, 14), P(94, 6, 19)), "cannot override method `baz` with invalid parameter name, is `b`, should be `a`\n  previous definition found in `Bar`, with signature: `def baz(a: Std::Int): Std::Int`"),
			},
		},
		"override the method with incompatible param type": {
			input: `
				class Bar
					def baz(a: Int): Int then a
				end
				class Foo < Bar
					def baz(a: Char): Int then 1
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(89, 6, 14), P(95, 6, 20)), "cannot override method `baz` with invalid parameter type, is `Std::Char`, should be `Std::Int`\n  previous definition found in `Bar`, with signature: `def baz(a: Std::Int): Std::Int`"),
			},
		},
		"override the method with narrower param type": {
			input: `
				class Bar
					def baz(a: Value): Value then a
				end
				class Foo < Bar
					def baz(a: Int): Value then 1
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(93, 6, 14), P(98, 6, 19)), "cannot override method `baz` with invalid parameter type, is `Std::Int`, should be `Std::Value`\n  previous definition found in `Bar`, with signature: `def baz(a: Std::Value): Std::Value`"),
			},
		},
		"override the method with wider param type": {
			input: `
				class Bar
					def baz(a: Int): Int then a
				end
				class Foo < Bar
					def baz(a: Value): Int then 1
				end
			`,
		},
		"override the method with incompatible return type": {
			input: `
				class Bar
					def baz(a: Int): Int then a
				end
				class Foo < Bar
					def baz(a: Int): String then "a"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(98, 6, 23), P(103, 6, 28)), "cannot override method `baz` with a different return type, is `Std::String`, should be `Std::Int`\n  previous definition found in `Bar`, with signature: `def baz(a: Std::Int): Std::Int`"),
			},
		},
		"override the method with narrower return type": {
			input: `
				class Bar
					def baz(a: Object): Value then a
				end
				class Foo < Bar
					def baz(a: Object): String then "a"
				end
			`,
		},
		"override the method with wider return type": {
			input: `
				class Bar
					def baz(a: String): String then a
				end
				class Foo < Bar
					def baz(a: String): Value then "a"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(107, 6, 26), P(111, 6, 30)), "cannot override method `baz` with a different return type, is `Std::Value`, should be `Std::String`\n  previous definition found in `Bar`, with signature: `def baz(a: Std::String): Std::String`"),
			},
		},
		"override the method with no return type": {
			input: `
				class Bar
					def baz(a: String): String then a
				end
				class Foo < Bar
					def baz(a: String); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(87, 6, 6), P(109, 6, 28)), "cannot override method `baz` with a different return type, is `void`, should be `Std::String`\n  previous definition found in `Bar`, with signature: `def baz(a: Std::String): Std::String`"),
			},
		},
		"override void method with a new return type": {
			input: `
				class Bar
					def baz(a: String); end
				end
				class Foo < Bar
					def baz(a: String): String then a
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

func TestMethodDefinition(t *testing.T) {
	tests := testTable{
		"declare within a method": {
			input: `
				def foo
					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(29, 3, 17)), "method definitions cannot appear in this context"),
			},
		},
		"declare with an invalid implicit return value": {
			input: `
				def foo: String
					5
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(26, 3, 6), P(26, 3, 6)), "type `5` cannot be assigned to type `Std::String`"),
			},
		},
		"declare with a valid implicit return value": {
			input: `
				def foo: String
					"lol"
				end
			`,
		},
		"redeclare the method in the same class with incompatible signature": {
			input: `
				class Foo
					def baz(a: Int): String then "a"
					def baz(): void; end
				end
			`,
		},
		"declare an abstract method with a body": {
			input: `
				abstract class Foo
					abstract def baz(a: Int)
						3
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(29, 3, 6), P(69, 5, 8)), "method `baz` cannot have a body because it is abstract"),
			},
		},
		"declare an interface method with a body": {
			input: `
				interface Foo
					def baz(a: Int)
						3
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 3, 6), P(55, 5, 8)), "method `baz` cannot have a body because it is abstract"),
			},
		},
		"declare an abstract method in an abstract class": {
			input: `
				abstract class Foo
					abstract def baz(a: Int); end
				end
			`,
		},
		"declare an abstract sig in an abstract class": {
			input: `
				abstract class Foo
					sig baz(a: Int)
				end
			`,
		},
		"declare an abstract method in a non-abstract class": {
			input: `
				class Foo
					abstract def baz(a: Int); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(48, 3, 34)), "cannot declare abstract method `baz` in non-abstract class `Foo`"),
			},
		},
		"declare an abstract sig in a non-abstract class": {
			input: `
				class Foo
					sig baz(a: Int)
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(34, 3, 20)), "cannot declare abstract method `baz` in non-abstract class `Foo`"),
			},
		},
		"declare an abstract method in an abstract mixin": {
			input: `
				abstract mixin Foo
					abstract def baz(a: Int); end
				end
			`,
		},
		"declare an abstract method in a non-abstract mixin": {
			input: `
				mixin Foo
					abstract def baz(a: Int); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(48, 3, 34)), "cannot declare abstract method `baz` in non-abstract mixin `Foo`"),
			},
		},
		"declare an abstract method in a module": {
			input: `
				module Foo
					abstract def baz(a: Int); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(21, 3, 6), P(49, 3, 34)), "cannot declare abstract method `baz` in this context"),
			},
		},
		"methods get hoisted to the top": {
			input: `
			  foo()
				def foo; end
			`,
		},
		"methods can reference each other": {
			input: `
				def foo then bar()
				def bar then foo()
			`,
		},
		"instance variable parameter without a type": {
			input: `
				class Foo
					def baz(@a); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 3, 14), P(29, 3, 15)), "cannot infer the type of instance variable `a`"),
			},
		},
		"inferred instance variable parameter type": {
			input: `
				class Foo
				  var @a: String
					def baz(@a); end
				end
			`,
		},
		"explicit instance variable parameter type": {
			input: `
				class Foo
				  var @a: String
					def baz(@a: String); end
				end
			`,
		},
		"explicit instance variable parameter subtype": {
			input: `
				class Foo
				  var @a: String?
					def baz(@a: String); end
				end
			`,
		},
		"explicit instance variable parameter supertype": {
			input: `
				class Foo
				  var @a: String
					def baz(@a: String?); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(53, 4, 18), P(59, 4, 24)), "type `Std::String?` cannot be assigned to instance variable `@a` of type `Std::String`"),
			},
		},
		"instance variable parameter takes the explicit type": {
			input: `
				class Foo
				  var @a: String?
					def baz(@a: String)
						var b: String = a
					end
				end
			`,
		},
		"instance variable parameter cannot be assigned to incompatible variable": {
			input: `
				class Foo
				  var @a: String?
					def baz(@a: String)
						var b: Int = a
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(81, 5, 20), P(81, 5, 20)), "type `Std::String` cannot be assigned to type `Std::Int`"),
			},
		},
		"instance variable parameter declares an instance variable": {
			input: `
				class Foo
					def baz(@a: String)
						var b: String = @a
					end
				end
			`,
		},
		"instance variable parameter declares an instance variable tha can be redeclared with the same type": {
			input: `
				class Foo
					def baz(@a: String)
						var b: String = @a
					end

					var @a: String
				end
			`,
		},
		"instance variable parameter declares an instance variable tha cannot be redeclared with a different type": {
			input: `
				class Foo
					def baz(@a: String)
						var b: String = @a
					end

					var @a: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(80, 7, 6), P(90, 7, 16)), "cannot redeclare instance variable `@a` with a different type, is `Std::Int`, should be `Std::String`, previous definition found in `Foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestMethodCalls(t *testing.T) {
	tests := testTable{
		"call has the same return type as the method": {
			input: `
				module Foo
					def baz(a: Int): Int then a
				end
				var a: Int = Foo.baz(5)
			`,
		},
		"cannot make nil-safe call on a non nilable receiver": {
			input: `
				module Foo
					def baz(a: Int): Int then a
				end
				Foo?.baz(5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(61, 5, 5), P(71, 5, 15)), "cannot make a nil-safe call on type `Foo` which is not nilable"),
			},
		},
		"can make nil-safe call on a nilable receiver": {
			input: `
				module Foo
					def baz(a: Int): Int then a
				end
				var nilableFoo: Foo? = Foo
				nilableFoo?.baz(5)
			`,
		},

		"cascade call returns the receiver": {
			input: `
				module Foo
					def baz(a: Int): Int then a
				end
				var a: Foo = Foo..baz(5)
			`,
		},
		"cannot make nil-safe cascade call on a non nilable receiver": {
			input: `
				module Foo
					def baz(a: Int): Int then a
				end
				Foo?..baz(5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(61, 5, 5), P(72, 5, 16)), "cannot make a nil-safe call on type `Foo` which is not nilable"),
			},
		},
		"can make nil-safe cascade call on a nilable receiver": {
			input: `
				module Foo
					def baz(a: Int): Int then a
				end
				var nilableFoo: Foo? = Foo
				nilableFoo?..baz(5)
			`,
		},
		"nil-safe cascade call returns a nilable receiver": {
			input: `
				module Foo
					def baz(a: Int): Int then a
				end
				var nilableFoo: Foo? = Foo
				var b: 8 = nilableFoo?..baz(5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(103, 6, 16), P(121, 6, 34)), "type `Foo?` cannot be assigned to type `8`"),
			},
		},

		"missing required argument": {
			input: `
				module Foo
					def baz(bar: String, c: Int); end
				end
				Foo.baz("foo")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(67, 5, 5), P(80, 5, 18)), "argument `c` is missing in call to `baz`"),
			},
		},
		"all required positional arguments": {
			input: `
				module Foo
					def baz(bar: String, c: Int); end
				end
				Foo.baz("foo", 5)
			`,
		},
		"all required positional arguments with wrong type": {
			input: `
				module Foo
					def baz(bar: String, c: Int); end
				end
				Foo.baz(123.4, 5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(75, 5, 13), P(79, 5, 17)), "expected type `Std::String` for parameter `bar` in call to `baz`, got type `123.4`"),
			},
		},
		"too many positional arguments": {
			input: `
				module Foo
					def baz(bar: String, c: Int); end
				end
				Foo.baz("foo", 5, 28, 9, 0)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(67, 5, 5), P(93, 5, 31)), "expected 2 arguments in call to `baz`, got 5"),
			},
		},
		"missing required argument with named argument": {
			input: `
				module Foo
					def baz(bar: String, c: Int); end
				end
				Foo.baz(bar: "foo")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(67, 5, 5), P(85, 5, 23)), "argument `c` is missing in call to `baz`"),
			},
		},
		"all required named arguments": {
			input: `
				module Foo
					def baz(bar: String, c: Int); end
				end
				Foo.baz(c: 5, bar: "foo")
			`,
		},
		"all required named arguments with wrong type": {
			input: `
				module Foo
					def baz(bar: String, c: Int); end
				end
				Foo.baz(c: 5, bar: 123.4)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(81, 5, 19), P(90, 5, 28)), "expected type `Std::String` for parameter `bar` in call to `baz`, got type `123.4`"),
			},
		},
		"duplicated positional argument as named argument": {
			input: `
				module Foo
					def baz(bar: String, c: Int); end
				end
				Foo.baz("foo", 5, bar: 9)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 5, 23), P(90, 5, 28)), "duplicated argument `bar` in call to `baz`"),
				error.NewFailure(L("<main>", P(85, 5, 23), P(90, 5, 28)), "expected type `Std::String` for parameter `bar` in call to `baz`, got type `9`"),
			},
		},
		"duplicated named argument": {
			input: `
				module Foo
					def baz(bar: String, c: Int); end
				end
				Foo.baz("foo", 2, c: 3, c: 9)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 5, 23), P(88, 5, 26)), "duplicated argument `c` in call to `baz`"),
				error.NewFailure(L("<main>", P(91, 5, 29), P(94, 5, 32)), "duplicated argument `c` in call to `baz`"),
			},
		},
		"call with missing optional argument": {
			input: `
				module Foo
					def baz(bar: String, c: Int = 3); end
				end
				Foo.baz("foo")
			`,
		},
		"call with optional argument": {
			input: `
				module Foo
					def baz(bar: String, c: Int = 3); end
				end
				Foo.baz("foo", 9)
			`,
		},
		"call with missing rest arguments": {
			input: `
				module Foo
					def baz(*b: Float); end
				end
				Foo.baz
			`,
		},
		"call with rest arguments": {
			input: `
				module Foo
					def baz(*b: Float); end
				end
				Foo.baz 1.2, 56.9, .5
			`,
		},
		"call with rest arguments with wrong type": {
			input: `
				module Foo
					def baz(*b: Float); end
				end
				Foo.baz 1.2, 5, "foo", .5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(70, 5, 18), P(70, 5, 18)), "expected type `Std::Float` for rest parameter `*b` in call to `baz`, got type `5`"),
				error.NewFailure(L("<main>", P(73, 5, 21), P(77, 5, 25)), "expected type `Std::Float` for rest parameter `*b` in call to `baz`, got type `\"foo\"`"),
			},
		},
		"call with rest argument given by name": {
			input: `
				module Foo
					def baz(*b: Float); end
				end
				Foo.baz b: []
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(65, 5, 13), P(69, 5, 17)), "nonexistent parameter `b` given in call to `baz`"),
			},
		},
		"call with required post arguments": {
			input: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
				Foo.baz("foo", 3)
			`,
		},
		"call with missing post argument": {
			input: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
				Foo.baz("foo")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(78, 5, 5), P(91, 5, 18)), "argument `c` is missing in call to `baz`"),
			},
		},
		"call with rest and post arguments": {
			input: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
				Foo.baz("foo", 2.5, .9, 128.1, 3)
			`,
		},
		"call with rest and post arguments and wrong type in post": {
			input: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
				Foo.baz("foo", 2.5, .9, 128.1, 3.2)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(109, 5, 36), P(111, 5, 38)), "expected type `Std::Int` for parameter `c` in call to `baz`, got type `3.2`"),
			},
		},
		"call with rest and post arguments and wrong type in rest": {
			input: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
				Foo.baz("foo", 212, .9, '282', 3)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(93, 5, 20), P(95, 5, 22)), "expected type `Std::Float` for rest parameter `*b` in call to `baz`, got type `212`"),
				error.NewFailure(L("<main>", P(102, 5, 29), P(106, 5, 33)), "expected type `Std::Float` for rest parameter `*b` in call to `baz`, got type `\"282\"`"),
			},
		},
		"call with rest arguments and missing post argument": {
			input: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
				Foo.baz("foo", 2.5, .9, 128.1)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(102, 5, 29), P(106, 5, 33)), "expected type `Std::Int` for parameter `c` in call to `baz`, got type `128.1`"),
			},
		},
		"call with named post argument": {
			input: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
				Foo.baz("foo", c: 3)
			`,
		},
		"call with named pre rest argument": {
			input: `
				module Foo
					def baz(bar: String, *b: Float, c: Int); end
				end
				Foo.baz(bar: "foo", c: 3)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(78, 5, 5), P(102, 5, 29)), "expected 1... positional arguments in call to `baz`, got 0"),
			},
		},
		"call without named rest arguments": {
			input: `
				module Foo
					def baz(bar: String, c: Int, **rest: Int); end
				end
				Foo.baz("foo", 5)
			`,
		},
		"call with named rest arguments": {
			input: `
				module Foo
					def baz(bar: String, c: Int, **rest: Int); end
				end
				Foo.baz("foo", d: 25, c: 5, e: 11)
			`,
		},
		"call with named rest arguments with wrong type": {
			input: `
				module Foo
					def baz(bar: String, c: Int, **rest: Int); end
				end
				Foo.baz("foo", d: .2, c: 5, e: .1)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(95, 5, 20), P(99, 5, 24)), "expected type `Std::Int` for named rest parameter `**rest` in call to `baz`, got type `0.2`"),
				error.NewFailure(L("<main>", P(108, 5, 33), P(112, 5, 37)), "expected type `Std::Int` for named rest parameter `**rest` in call to `baz`, got type `0.1`"),
			},
		},

		"call setter with matching argument": {
			input: `
				class Foo
					def baz=(bar: String); end
				end
				Foo().baz = "bar"
			`,
		},
		"the return type of a setter is the same as the argument": {
			input: `
				class Foo
					def baz=(bar: String); end
				end
				var a: String = Foo().baz = "bar"
			`,
		},
		"call setter with non-matching argument": {
			input: `
				class Foo
					def baz=(bar: String); end
				end
				Foo().baz = 1
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(71, 5, 17), P(71, 5, 17)), "expected type `Std::String` for parameter `bar` in call to `baz=`, got type `1`"),
			},
		},
		"call nonexistent setter": {
			input: `
				class Foo
					def baz=(bar: String); end
				end
				Foo().foo = 1
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(59, 5, 5), P(71, 5, 17)), "method `foo=` is not defined on type `Foo`"),
			},
		},

		"call subscript setter with matching argument": {
			input: `
				class Foo
					def []=(key: String, value: Int); end
				end
				Foo()["foo"] = 1
			`,
		},
		"the return type of a subscript setter is the same as the value": {
			input: `
				class Foo
					def []=(key: String, value: Int); end
				end
				var a: Int = Foo()["foo"] = 1
			`,
		},
		"call subscript setter with non-matching argument": {
			input: `
				class Foo
					def []=(key: String, value: Int); end
				end
				Foo()[1] = "foo"
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(76, 5, 11), P(76, 5, 11)), "expected type `Std::String` for parameter `key` in call to `[]=`, got type `1`"),
				error.NewFailure(L("<main>", P(81, 5, 16), P(85, 5, 20)), "expected type `Std::Int` for parameter `value` in call to `[]=`, got type `\"foo\"`"),
			},
		},
		"call nonexistent subscript setter": {
			input: `
				class Foo; end
				Foo()["foo"] = 1
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 3, 5), P(39, 3, 20)), "method `[]=` is not defined on type `Foo`"),
			},
		},

		"call subscript with matching argument": {
			input: `
				class Foo
					def [](key: String): String then key
				end
				var f: String = Foo()["foo"]
			`,
		},
		"call subscript with non-matching argument": {
			input: `
				class Foo
					def [](key: String): String then key
				end
				Foo()[1]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(75, 5, 11), P(75, 5, 11)), "expected type `Std::String` for parameter `key` in call to `[]`, got type `1`"),
			},
		},
		"call nonexistent subscript": {
			input: `
				class Foo; end
				Foo()["foo"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 3, 5), P(35, 3, 16)), "method `[]` is not defined on type `Foo`"),
			},
		},

		"call nil-safe subscript on non nilable type": {
			input: `
				class Foo
					def [](key: String): String then key
				end
				var f: String = Foo()?["foo"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 5, 21), P(97, 5, 33)), "cannot make a nil-safe call on type `Foo` which is not nilable"),
			},
		},
		"call nil-safe subscript with matching argument": {
			input: `
				class Foo
					def [](key: String): String then key
				end
				var a: Foo? = Foo()
				var f: String? = a?["foo"]
			`,
		},
		"call nil-safe subscript and make the return type nilable": {
			input: `
				class Foo
					def [](key: String): String then key
				end
				var a: Foo? = Foo()
				var f: String = a?["foo"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(109, 6, 21), P(117, 6, 29)), "type `Std::String?` cannot be assigned to type `Std::String`"),
			},
		},
		"call nil-safe subscript with non-matching argument": {
			input: `
				class Foo
					def [](key: String): String then key
				end
				var a: Foo? = Foo()
				a?[1]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(96, 6, 8), P(96, 6, 8)), "expected type `Std::String` for parameter `key` in call to `[]`, got type `1`"),
			},
		},
		"call nonexistent nil-safe subscript": {
			input: `
				class Foo; end
				var a: Foo? = Foo()
				a?["foo"]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 4, 5), P(56, 4, 13)), "method `[]` is not defined on type `Foo`"),
			},
		},

		"call nonexistent increment": {
			input: `
				class Foo; end
				f := Foo()
				f++
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 4, 5), P(41, 4, 7)), "method `++` is not defined on type `Foo`"),
			},
		},
		"call increment": {
			input: `
				f := 5
				f++
			`,
		},
		"the return type of increment is as expected": {
			input: `
				f := 5
				var g: 2 = f++
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 16), P(29, 3, 18)), "type `Std::Int` cannot be assigned to type `2`"),
			},
		},
		"increment with incompatible return type": {
			input: `
				class Foo
					def ++: String
						"foo"
					end
				end

				f := Foo()
				f++
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(84, 9, 5), P(86, 9, 7)), "type `Std::String` cannot be assigned to type `Foo`"),
			},
		},
		"call increment on nonexistent variable": {
			input: `
				f++
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(5, 2, 5), P(5, 2, 5)), "undefined local `f`"),
			},
		},

		"call nonexistent decrement": {
			input: `
				class Foo; end
				f := Foo()
				f--
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 4, 5), P(41, 4, 7)), "method `--` is not defined on type `Foo`"),
			},
		},
		"call decrement": {
			input: `
				f := 5
				f--
			`,
		},
		"the return type of decrement is as expected": {
			input: `
				f := 5
				var g: 2 = f--
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 16), P(29, 3, 18)), "type `Std::Int` cannot be assigned to type `2`"),
			},
		},
		"decrement with incompatible return type": {
			input: `
				class Foo
					def --: String
						"foo"
					end
				end

				f := Foo()
				f--
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(84, 9, 5), P(86, 9, 7)), "type `Std::String` cannot be assigned to type `Foo`"),
			},
		},
		"call decrement on nonexistent variable": {
			input: `
				f--
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(5, 2, 5), P(5, 2, 5)), "undefined local `f`"),
			},
		},

		"pipe operator add an inexistent argument to a constructor call": {
			input: `
				class Foo; end
				1 |> Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(29, 3, 10), P(33, 3, 14)), "expected 0 arguments in call to `#init`, got 1"),
			},
		},
		"pipe operator add an argument with an incompatible type to a constructor call": {
			input: `
				class Foo
					init(a: Float); end
				end
				1 |> Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(52, 5, 5), P(52, 5, 5)), "expected type `Std::Float` for parameter `a` in call to `#init`, got type `1`"),
			},
		},
		"pipe operator add an additional required argument with a compatible type to a constructor call": {
			input: `
				class Foo
					init(a: Int, b: Float); end
				end
				1 |> Foo(2.2)
			`,
		},
		"pipe operator add an argument with a compatible type to a constructor call": {
			input: `
				class Foo
					init(a: Int); end
				end
				1 |> Foo()
			`,
		},

		"pipe operator add an inexistent argument to a method call": {
			input: `
				class Foo
					def foo; end
				end
				f := Foo()
				1 |> f.foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(65, 6, 10), P(71, 6, 16)), "expected 0 arguments in call to `foo`, got 1"),
			},
		},
		"pipe operator add an argument with an incompatible type to a method call": {
			input: `
				class Foo
					def foo(a: Float); end
				end
				f := Foo()
				1 |> f.foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(70, 6, 5), P(70, 6, 5)), "expected type `Std::Float` for parameter `a` in call to `foo`, got type `1`"),
			},
		},
		"pipe operator add an additional required argument with a compatible type to a method call": {
			input: `
				class Foo
					def foo(a: Int, b: Float); end
				end
				f := Foo()
				1 |> f.foo(2.2)
			`,
		},
		"pipe operator add an argument with a compatible type to a method call": {
			input: `
				class Foo
					def foo(a: Int); end
				end
				f := Foo()
				1 |> f.foo()
			`,
		},

		"pipe has the same return type as the method": {
			input: `
				class Foo
					def foo(a: Int): String then "f"
				end
				f := Foo()
				var b: 9 = 1 |> f.foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(96, 6, 21), P(102, 6, 27)), "type `Std::String` cannot be assigned to type `9`"),
			},
		},

		"pipe operator add an inexistent argument to a receiverless method call": {
			input: `
				module Foo
					def foo; end
					1 |> foo()
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 4, 11), P(48, 4, 15)), "expected 0 arguments in call to `foo`, got 1"),
			},
		},
		"pipe operator add an argument with an incompatible type to a receiverless method call": {
			input: `
				module Foo
					def foo(a: Float); end
					1 |> foo()
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(49, 4, 6), P(49, 4, 6)), "expected type `Std::Float` for parameter `a` in call to `foo`, got type `1`"),
			},
		},
		"pipe operator add an additional required argument with a compatible type to a receiverless method call": {
			input: `
				module Foo
					def foo(a: Int, b: Float); end
					1 |> foo(2.2)
				end
			`,
		},
		"pipe operator add an argument with a compatible type to a receiverless method call": {
			input: `
				module Foo
					def foo(a: Int); end
					1 |> foo()
				end
			`,
		},

		"pipe operator add an inexistent argument to attribute access": {
			input: `
				module Foo
					def foo; end
					1 |> self.foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 4, 6), P(51, 4, 18)), "expected 0 arguments in call to `foo`, got 1"),
			},
		},
		"pipe operator add an argument with an incompatible type to attribute access": {
			input: `
				module Foo
					def foo(a: Float); end
					1 |> self.foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(49, 4, 6), P(49, 4, 6)), "expected type `Std::Float` for parameter `a` in call to `foo`, got type `1`"),
			},
		},
		"pipe operator add an argument with a compatible type to attribute access": {
			input: `
				module Foo
					def foo(a: Int); end
					1 |> self.foo
				end
			`,
		},

		"pipe operator add an inexistent argument to a call": {
			input: `
				class Foo
					def call; end
				end
				f := Foo()
				1 |> f.()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(66, 6, 10), P(69, 6, 13)), "expected 0 arguments in call to `call`, got 1"),
			},
		},
		"pipe operator add an argument with an incompatible type to a call": {
			input: `
				class Foo
					def call(a: Float); end
				end
				f := Foo()
				1 |> f.()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(71, 6, 5), P(71, 6, 5)), "expected type `Std::Float` for parameter `a` in call to `call`, got type `1`"),
			},
		},
		"pipe operator add an additional required argument with a compatible type to a call": {
			input: `
				class Foo
					def call(a: Int, b: Float); end
				end
				f := Foo()
				1 |> f.(2.2)
			`,
		},
		"pipe operator add an argument with a compatible type to a call": {
			input: `
				class Foo
					def call(a: Int); end
				end
				f := Foo()
				1 |> f.()
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestInitDefinition(t *testing.T) {
	tests := testTable{
		"define within a method": {
			input: `def foo; init; end; end`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 1, 10), P(17, 1, 18)), "method definitions cannot appear in this context"),
			},
		},
		"define in outer context": {
			input: `init; end`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(8, 1, 9)), "init definitions cannot appear outside of classes"),
			},
		},
		"define in module": {
			input: `
				module Foo
					init; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(21, 3, 6), P(29, 3, 14)), "init definitions cannot appear outside of classes"),
			},
		},
		"define in class": {
			input: `
				class Foo
					init; end
				end
			`,
		},
		"with parameters": {
			input: `
				class Foo
					init(a: Int); end
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

func TestConstructorCall(t *testing.T) {
	tests := testTable{
		"instantiate a class without a constructor": {
			input: `
				class Foo; end
				Foo()
			`,
		},
		"instantiate an abstract class": {
			input: `
				abstract class Foo; end
				Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 3, 5), P(37, 3, 9)), "cannot instantiate abstract class `Foo`"),
			},
		},
		"instantiate a class with a constructor": {
			input: `
				class Foo
					init(a: Int); end
				end
				Foo(1)
			`,
		},
		"instantiate a class with a constructor with a wrong type": {
			input: `
				class Foo
					init(a: String); end
				end
				Foo(1)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 5, 9), P(57, 5, 9)), "expected type `Std::String` for parameter `a` in call to `#init`, got type `1`"),
			},
		},
		"instantiate a class with an inherited constructor": {
			input: `
				class Bar
					init(a: Int); end
				end

				class Foo < Bar; end
				Foo(1)
			`,
		},
		"instantiate a class with an inherited constructor with a wrong type": {
			input: `
				class Bar
					init(a: String); end
				end

				class Foo < Bar; end
				Foo(1)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(83, 7, 9), P(83, 7, 9)), "expected type `Std::String` for parameter `a` in call to `#init`, got type `1`"),
			},
		},
		"call a method on an instantiated instance": {
			input: `
				class Foo
					init(a: String); end

					def bar; end
				end

				var foo = Foo("foo")
				foo.bar
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestMethodInheritance(t *testing.T) {
	tests := testTable{
		"call a method inherited from superclass": {
			input: `
				class Foo
					def baz(a: Int): Int then a
				end

				class Bar < Foo; end
				var bar = Bar()
				bar.baz(5)
			`,
		},
		"call a method inherited from mixin": {
			input: `
				mixin Bar
					def baz(a: Int): Int then a
				end

				class Foo
					include Bar
				end

				var foo = Foo()
				foo.baz(5)
			`,
		},
		"call a method on a mixin type": {
			input: `
				mixin Bar
					def baz(a: Int): Int then a
				end

				class Foo
					include Bar
				end

				var bar: Bar = Foo()
				bar.baz(5)
			`,
		},
		"call an inherited method on a mixin type": {
			input: `
				mixin Baz
					def baz(a: Int): Int then a
				end

				mixin Bar
				  include Baz
				end

				class Foo
					include Bar
				end

				var bar: Bar = Foo()
				bar.baz(5)
			`,
		},
		"call a class instance method on a mixin type": {
			input: `
				mixin Bar
					def baz(a: Int): Int then a
				end

				class Foo
					include Bar

					def foo; end
				end

				var bar: Bar = Foo()
				bar.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(145, 13, 5), P(151, 13, 11)), "method `foo` is not defined on type `Bar`"),
			},
		},
		"call a method on singleton type": {
			input: `
				class Foo
					singleton
						def foo; end
					end
				end

				var foo = Foo
				foo.foo
			`,
		},
		"call an inherited method on singleton type": {
			input: `
				class Foo
					singleton
						def foo; end
					end
				end

				class Bar < Foo; end

				var foo = Bar
				foo.foo
			`,
		},
		"call inherited singleton method from mixin": {
			input: `
				mixin Foo
					def foo; end
				end
				class Bar
					singleton
						include Foo
					end
				end

				var foo = Bar
				foo.foo
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestClosureLiteral(t *testing.T) {
	tests := testTable{
		"assign a valid closure to a variable": {
			input: `
				a := |a: Int|: Int -> a
			`,
		},
		"invalid parameter default value and return value": {
			input: `
				a := |a: Int = 2.3|: String -> a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 2, 20), P(22, 2, 22)), "type `2.3` cannot be assigned to type `Std::Int`"),
				error.NewFailure(L("<main>", P(36, 2, 36), P(36, 2, 36)), "type `Std::Int` cannot be assigned to type `Std::String`"),
			},
		},
		"without param type": {
			input: `
				a := |a| -> 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(11, 2, 11), P(11, 2, 11)), "cannot declare parameter `a` without a type"),
			},
		},
		"assign an invalid value to a closure type": {
			input: `
				a := |a: Int|: Int -> a
				a = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(37, 3, 9), P(37, 3, 9)), "type `Std::Int` does not implement interface `|a: Std::Int|: Std::Int`:\n\n  - missing method `call` with signature: `def call(a: Std::Int): Std::Int`\n"),
				error.NewFailure(L("<main>", P(37, 3, 9), P(37, 3, 9)), "type `3` cannot be assigned to type `|a: Std::Int|: Std::Int`"),
			},
		},
		"assign a compatible value to a closure type": {
			input: `
				var a: |a: Int|: String
				a = |a: Int, b: String = "foo"|: String -> b
			`,
		},
		"assign an incompatible closure to a closure type": {
			input: `
				a := |a: Int|: Int -> a
				a = |a: Float, b: String = "foo"|: String -> b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(37, 3, 9), P(78, 3, 50)), "type `|a: Std::Float, b?: Std::String|: Std::String` does not implement interface `|a: Std::Int|: Std::Int`:\n\n  - incorrect implementation of `call`\n      is:        `def call(a: Std::Float, b?: Std::String): Std::String`\n      should be: `def call(a: Std::Int): Std::Int`\n"),
				error.NewFailure(L("<main>", P(37, 3, 9), P(78, 3, 50)), "type `|a: Std::Float, b?: Std::String|: Std::String` cannot be assigned to type `|a: Std::Int|: Std::Int`"),
			},
		},
		"take param and return types from closure defined as a formal parameter": {
			input: `
				def foo(fn: |a: String|: Int); end
				foo() |a| -> 5
			`,
		},
		"invalid closure argument": {
			input: `
				def foo(fn: |a: String|: Int); end
				foo() |i| -> 2.5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 3, 12), P(51, 3, 12)), "cannot override method `call` with invalid parameter name, is `i`, should be `a`\n  previous definition found in `|a: Std::String|: Std::Int`, with signature: `abstract def call(a: Std::String): Std::Int`"),
				error.NewFailure(L("<main>", P(57, 3, 18), P(59, 3, 20)), "type `2.5` cannot be assigned to type `Std::Int`"),
			},
		},
		"call a closure": {
			input: `
				a := |a: Int|: Int -> a
				a.(9)
				a.call(3 + 8)
			`,
		},
		"call a closure with invalid arguments": {
			input: `
				a := |a: Int|: Int -> a
				a.(2.5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 8), P(38, 3, 10)), "expected type `Std::Int` for parameter `a` in call to `call`, got type `2.5`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
