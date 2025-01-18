package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestAttrDefinition(t *testing.T) {
	tests := testTable{
		"declare within extend where": {
			input: `
				class E[T]
					extend where T < String
						attr foo: String?
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(56, 4, 12), P(67, 4, 23)), "cannot declare instance variable `@foo` in this context"),
			},
		},
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
				error.NewFailure(L("<main>", P(48, 4, 11), P(55, 4, 18)), "method `Foo.:foo=` is not a valid override of `Foo.:foo=`\n  is:        `def foo=(foo: Std::Int): void`\n  should be: `def foo=(foo: Std::String?): void`\n\n  - has an incompatible parameter, is `foo: Std::Int`, should be `foo: Std::String?`"),
				error.NewFailure(L("<main>", P(48, 4, 11), P(55, 4, 18)), "method `Foo.:foo` is not a valid override of `Foo.:foo`\n  is:        `def foo(): Std::Int`\n  should be: `def foo(): Std::String?`\n\n  - has a different return type, is `Std::Int`, should be `Std::String?`"),
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
					attr foo: Int?
				end
				class Bar < Foo
					attr foo: Int?
				end
			`,
		},
		"override an attr with a different type in a child class": {
			input: `
				class Foo
					attr foo: Int?
				end
				class Bar < Foo
					attr foo: String?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(78, 6, 16), P(84, 6, 22)), "type `Std::String?` cannot be assigned to instance variable `@foo` of type `Std::Int?`"),
				error.NewFailure(L("<main>", P(73, 6, 11), P(84, 6, 22)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String?`, should be `Std::Int?`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(73, 6, 11), P(84, 6, 22)), "method `Bar.:foo=` is not a valid override of `Foo.:foo=`\n  is:        `def foo=(foo: Std::String?): void`\n  should be: `def foo=(foo: Std::Int?): void`\n\n  - has an incompatible parameter, is `foo: Std::String?`, should be `foo: Std::Int?`"),
				error.NewFailure(L("<main>", P(73, 6, 11), P(84, 6, 22)), "method `Bar.:foo` is not a valid override of `Foo.:foo`\n  is:        `def foo(): Std::String?`\n  should be: `def foo(): Std::Int?`\n\n  - has a different return type, is `Std::String?`, should be `Std::Int?`"),
				error.NewFailure(L("<main>", P(73, 6, 11), P(84, 6, 22)), "type `Std::Int?` cannot be assigned to type `Std::String?`"),
			},
		},
		"override an instance variable using an attr with the same type in a child class": {
			input: `
				class Foo
					var @foo: Int?
				end
				class Bar < Foo
					attr foo: Int?
				end
			`,
		},
		"override an instance variable using an attr with a different type in a child class": {
			input: `
				class Foo
					var @foo: Int?
				end
				class Bar < Foo
					attr foo: String?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(78, 6, 16), P(84, 6, 22)), "type `Std::String?` cannot be assigned to instance variable `@foo` of type `Std::Int?`"),
				error.NewFailure(L("<main>", P(73, 6, 11), P(84, 6, 22)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String?`, should be `Std::Int?`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(73, 6, 11), P(84, 6, 22)), "type `Std::Int?` cannot be assigned to type `Std::String?`"),
			},
		},
		"override a method using an attr with the same parameter type in a child class": {
			input: `
				class Foo
					def foo=(foo: String?); end
				end
				class Bar < Foo
					attr foo: String?
				end
			`,
		},
		"override a method using an attr with a different parameter type in a child class": {
			input: `
				class Foo
					def foo=(foo: String); end
				end
				class Bar < Foo
					attr foo: Int?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 6, 11), P(93, 6, 19)), "method `Bar.:foo=` is not a valid override of `Foo.:foo=`\n  is:        `def foo=(foo: Std::Int?): void`\n  should be: `def foo=(foo: Std::String): void`\n\n  - has an incompatible parameter, is `foo: Std::Int?`, should be `foo: Std::String`"),
			},
		},
		"override a method using an attr with a different return type in a child class": {
			input: `
				class Foo
					def foo: String then "foo"
				end
				class Bar < Foo
					attr foo: Int?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 6, 11), P(93, 6, 19)), "method `Bar.:foo` is not a valid override of `Foo.:foo`\n  is:        `def foo(): Std::Int?`\n  should be: `def foo(): Std::String`\n\n  - has a different return type, is `Std::Int?`, should be `Std::String`"),
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
		"declare within extend where": {
			input: `
				class E[T]
					extend where T < String
						getter foo: String
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(58, 4, 14), P(68, 4, 24)), "cannot declare instance variable `@foo` in this context"),
				error.NewFailure(L("<main>", P(58, 4, 14), P(68, 4, 24)), "undefined instance variable `@foo` in type `E`"),
			},
		},
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
				error.NewFailure(L("<main>", P(52, 4, 13), P(59, 4, 20)), "method `Foo.:foo` is not a valid override of `Foo.:foo`\n  is:        `def foo(): Std::Int`\n  should be: `def foo(): Std::String?`\n\n  - has a different return type, is `Std::Int`, should be `Std::String?`"),
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
					getter foo: Int?
				end
				class Bar < Foo
					getter foo: Int?
				end
			`,
		},
		"override a getter with a different type in a child class": {
			input: `
				class Foo
					getter foo: Int?
				end
				class Bar < Foo
					getter foo: String?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(77, 6, 13), P(88, 6, 24)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String?`, should be `Std::Int?`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(77, 6, 13), P(88, 6, 24)), "method `Bar.:foo` is not a valid override of `Foo.:foo`\n  is:        `def foo(): Std::String?`\n  should be: `def foo(): Std::Int?`\n\n  - has a different return type, is `Std::String?`, should be `Std::Int?`"),
				error.NewFailure(L("<main>", P(77, 6, 13), P(88, 6, 24)), "type `Std::Int?` cannot be assigned to type `Std::String?`"),
			},
		},
		"override an instance variable using getter with the same type in a child class": {
			input: `
				class Foo
					var @foo: Int?
				end
				class Bar < Foo
					getter foo: Int?
				end
			`,
		},
		"override an instance variable using getter with a different type in a child class": {
			input: `
				class Foo
					var @foo: Int?
				end
				class Bar < Foo
					getter foo: String?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(75, 6, 13), P(86, 6, 24)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String?`, should be `Std::Int?`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(75, 6, 13), P(86, 6, 24)), "type `Std::Int?` cannot be assigned to type `Std::String?`"),
			},
		},
		"override a method using a getter with the same return type in a child class": {
			input: `
				class Foo
					def foo: String? then "foo"
				end
				class Bar < Foo
					getter foo: String?
				end
			`,
		},
		"override a method using a getter with a different return type in a child class": {
			input: `
				class Foo
					def foo: String then "foo"
				end
				class Bar < Foo
					getter foo: Int?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(87, 6, 13), P(95, 6, 21)), "method `Bar.:foo` is not a valid override of `Foo.:foo`\n  is:        `def foo(): Std::Int?`\n  should be: `def foo(): Std::String`\n\n  - has a different return type, is `Std::Int?`, should be `Std::String`"),
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
		"declare within extend where": {
			input: `
				class E[T]
					extend where T < String
						setter foo: String?
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(58, 4, 14), P(69, 4, 25)), "cannot declare instance variable `@foo` in this context"),
			},
		},
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
				error.NewFailure(L("<main>", P(52, 4, 13), P(59, 4, 20)), "method `Foo.:foo=` is not a valid override of `Foo.:foo=`\n  is:        `def foo=(foo: Std::Int): void`\n  should be: `def foo=(foo: Std::String?): void`\n\n  - has an incompatible parameter, is `foo: Std::Int`, should be `foo: Std::String?`"),
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
					setter foo: Int?
				end
				class Bar < Foo
					setter foo: Int?
				end
			`,
		},
		"override a setter with a different type in a child class": {
			input: `
				class Foo
					setter foo: Int?
				end
				class Bar < Foo
					setter foo: String?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(82, 6, 18), P(88, 6, 24)), "type `Std::String?` cannot be assigned to instance variable `@foo` of type `Std::Int?`"),
				error.NewFailure(L("<main>", P(77, 6, 13), P(88, 6, 24)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String?`, should be `Std::Int?`, previous definition found in `Foo`"),
				error.NewFailure(L("<main>", P(77, 6, 13), P(88, 6, 24)), "method `Bar.:foo=` is not a valid override of `Foo.:foo=`\n  is:        `def foo=(foo: Std::String?): void`\n  should be: `def foo=(foo: Std::Int?): void`\n\n  - has an incompatible parameter, is `foo: Std::String?`, should be `foo: Std::Int?`"),
			},
		},
		"override an instance variable using setter with the same type in a child class": {
			input: `
				class Foo
					var @foo: Int?
				end
				class Bar < Foo
					setter foo: Int?
				end
			`,
		},
		"override an instance variable using setter with a different type in a child class": {
			input: `
				class Foo
					var @foo: Int?
				end
				class Bar < Foo
					setter foo: String?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(80, 6, 18), P(86, 6, 24)), "type `Std::String?` cannot be assigned to instance variable `@foo` of type `Std::Int?`"),
				error.NewFailure(L("<main>", P(75, 6, 13), P(86, 6, 24)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String?`, should be `Std::Int?`, previous definition found in `Foo`"),
			},
		},
		"override a method using a setter with the same parameter type in a child class": {
			input: `
				class Foo
					def foo=(foo: String?); end
				end
				class Bar < Foo
					setter foo: String?
				end
			`,
		},
		"override a method using a setter with a different parameter type in a child class": {
			input: `
				class Foo
					def foo=(foo: String); end
				end
				class Bar < Foo
					setter foo: Int?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(87, 6, 13), P(95, 6, 21)), "method `Bar.:foo=` is not a valid override of `Foo.:foo=`\n  is:        `def foo=(foo: Std::Int?): void`\n  should be: `def foo=(foo: Std::String): void`\n\n  - has an incompatible parameter, is `foo: Std::Int?`, should be `foo: Std::String`"),
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
				error.NewFailure(L("<main>", P(82, 7, 6), P(95, 7, 19)), "method `Bar.:baz` is not a valid override of `Foo.:baz`\n  is:        `def baz(): void`\n  should be: `def baz(a: Std::Int): Std::Int`\n\n  - has a different return type, is `void`, should be `Std::Int`\n  - has less parameters"),
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
				error.NewFailure(L("<main>", P(93, 8, 6), P(106, 8, 19)), "method `Bar.:baz` is not a valid override of `Foo.:baz`\n  is:        `def baz(): void`\n  should be: `def baz(a: Std::Int): Std::Int`\n\n  - has a different return type, is `void`, should be `Std::Int`\n  - has less parameters"),
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
				error.NewFailure(L("<main>", P(93, 8, 6), P(106, 8, 19)), "method `Bar.:baz` is not a valid override of `Foo.:baz`\n  is:        `def baz(): void`\n  should be: `def baz(a: Std::Int): Std::Int`\n\n  - has a different return type, is `void`, should be `Std::Int`\n  - has less parameters"),
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
				error.NewFailure(L("<main>", P(88, 6, 6), P(114, 6, 32)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz(a: Std::Int): Std::Int`\n  should be: `sealed def baz(a: Std::Int): Std::Int`\n\n  - method `Bar.:baz` is sealed and cannot be overridden"),
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
				error.NewFailure(L("<main>", P(60, 4, 6), P(86, 4, 32)), "method `Bar.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz(a: Std::Int): Std::Int`\n  should be: `sealed def baz(a: Std::Int): Std::Int`\n\n  - method `Bar.:baz` is sealed and cannot be overridden"),
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
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(62, 4, 6), P(95, 4, 39)), "method `Bar.:baz` is not a valid override of `Bar.:baz`\n  is:        `abstract def baz(a: Std::Int): Std::Int`\n  should be: `def baz(a: Std::Int): Std::Int`\n\n  - has a different modifier, is `abstract`, should be `default`"),
			},
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
				error.NewFailure(L("<main>", P(91, 7, 6), P(124, 7, 39)), "method `Bar.:baz` is not a valid override of `Foo.:baz`\n  is:        `abstract def baz(a: Std::Int): Std::Int`\n  should be: `def baz(a: Std::Int): Std::Int`\n\n  - has a different modifier, is `abstract`, should be `default`"),
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
				error.NewFailure(L("<main>", P(81, 6, 6), P(107, 6, 32)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz(b: Std::Int): Std::Int`\n  should be: `def baz(a: Std::Int): Std::Int`\n\n  - has an incompatible parameter, is `b: Std::Int`, should be `a: Std::Int`"),
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
				error.NewFailure(L("<main>", P(81, 6, 6), P(108, 6, 33)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz(a: Std::Char): Std::Int`\n  should be: `def baz(a: Std::Int): Std::Int`\n\n  - has an incompatible parameter, is `a: Std::Char`, should be `a: Std::Int`"),
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
				error.NewFailure(L("<main>", P(85, 6, 6), P(113, 6, 34)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz(a: Std::Int): Std::Value`\n  should be: `def baz(a: Std::Value): Std::Value`\n\n  - has an incompatible parameter, is `a: Std::Int`, should be `a: Std::Value`"),
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
				error.NewFailure(L("<main>", P(81, 6, 6), P(112, 6, 37)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz(a: Std::Int): Std::String`\n  should be: `def baz(a: Std::Int): Std::Int`\n\n  - has a different return type, is `Std::String`, should be `Std::Int`"),
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
				error.NewFailure(L("<main>", P(87, 6, 6), P(120, 6, 39)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz(a: Std::String): Std::Value`\n  should be: `def baz(a: Std::String): Std::String`\n\n  - has a different return type, is `Std::Value`, should be `Std::String`"),
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
				error.NewFailure(L("<main>", P(87, 6, 6), P(109, 6, 28)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz(a: Std::String): void`\n  should be: `def baz(a: Std::String): Std::String`\n\n  - has a different return type, is `void`, should be `Std::String`"),
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

		"override with more type parameters": {
			input: `
				class Bar
					def baz(a: String); end
				end
				class Foo < Bar
					def baz[V](a: String): String then a
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(77, 6, 6), P(112, 6, 41)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz[V](a: Std::String): Std::String`\n  should be: `def baz(a: Std::String): void`\n\n  - has a different number of type parameters, has `1`, should have `0`"),
			},
		},
		"override with less type parameters": {
			input: `
				class Bar
					def baz[V, T](a: String); end
				end
				class Foo < Bar
					def baz[V](a: String): String then a
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(83, 6, 6), P(118, 6, 41)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz[V](a: Std::String): Std::String`\n  should be: `def baz[V, T](a: Std::String): void`\n\n  - has a different number of type parameters, has `1`, should have `2`"),
			},
		},
		"override with different names of type parameters": {
			input: `
				class Bar
					def baz[V, T]; end
				end
				class Foo < Bar
					def baz[E, L]; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(72, 6, 6), P(89, 6, 23)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz[E, L](): void`\n  should be: `def baz[V, T](): void`\n\n  - has an incompatible type parameter, is `E`, should be `V`\n  - has an incompatible type parameter, is `L`, should be `T`"),
			},
		},

		"override with the same invariant type parameter": {
			input: `
				class Bar
					def baz[V](a: V); end
				end
				class Foo < Bar
					def baz[V](a: V); end
				end
			`,
		},
		"override with invariant type parameter and the same upper bound": {
			input: `
				class Bar
					def baz[V < Object](a: V); end
				end
				class Foo < Bar
					def baz[V < Object](a: V); end
				end
			`,
		},
		"override with invariant type parameter and wider upper bound": {
			input: `
				class Bar
					def baz[V < String](a: V); end
				end
				class Foo < Bar
					def baz[V < Object](a: V); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(84, 6, 6), P(113, 6, 35)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz[V < Std::Object](a: V): void`\n  should be: `def baz[V < Std::String](a: V): void`\n\n  - has an incompatible type parameter, is `V < Std::Object`, should be `V < Std::String`"),
			},
		},
		"override with invariant type parameter and narrower upper bound": {
			input: `
				class Bar
					def baz[V < Object](a: V); end
				end
				class Foo < Bar
					def baz[V < String](a: V); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(84, 6, 6), P(113, 6, 35)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz[V < Std::String](a: V): void`\n  should be: `def baz[V < Std::Object](a: V): void`\n\n  - has an incompatible type parameter, is `V < Std::String`, should be `V < Std::Object`"),
			},
		},

		"override with the same covariant type parameter": {
			input: `
				class Bar
					def baz[+V]: V then loop; end
				end
				class Foo < Bar
					def baz[+V]: V then loop; end
				end
			`,
		},
		"override with covariant type parameter and the same upper bound": {
			input: `
				class Bar
					def baz[+V < Object]: V then loop; end
				end
				class Foo < Bar
					def baz[+V < Object]: V then loop; end
				end
			`,
		},
		"override with covariant type parameter and wider upper bound": {
			input: `
				class Bar
					def baz[+V < String]: V then loop; end
				end
				class Foo < Bar
					def baz[+V < Value]: V then loop; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(92, 6, 6), P(128, 6, 42)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz[+V < Std::Value](): V`\n  should be: `def baz[+V < Std::String](): V`\n\n  - has an incompatible type parameter, is `+V < Std::Value`, should be `+V < Std::String`"),
			},
		},
		"override with covariant type parameter and narrower upper bound": {
			input: `
				class Bar
					def baz[+V < Value]: V then loop; end
				end
				class Foo < Bar
					def baz[+V < String]: V then loop; end
				end
			`,
		},

		"override with covariant type parameter and the same lower bound": {
			input: `
				class Bar
					def baz[+V > Object]: V then loop; end
				end
				class Foo < Bar
					def baz[+V > Object]: V then loop; end
				end
			`,
		},
		"override with covariant type parameter and wider lower bound": {
			input: `
				class Bar
					def baz[+V > String]: V then loop; end
				end
				class Foo < Bar
					def baz[+V > Value]: V then loop; end
				end
			`,
		},
		"override with covariant type parameter and narrower lower bound": {
			input: `
				class Bar
					def baz[+V > Value]: V then loop; end
				end
				class Foo < Bar
					def baz[+V > String]: V then loop; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(91, 6, 6), P(128, 6, 43)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz[+V > Std::String](): V`\n  should be: `def baz[+V > Std::Value](): V`\n\n  - has an incompatible type parameter, is `+V > Std::String`, should be `+V > Std::Value`"),
			},
		},

		"override with the same contravariant type parameter": {
			input: `
				class Bar
					def baz[-V](a: V); end
				end
				class Foo < Bar
					def baz[-V](a: V); end
				end
			`,
		},
		"override with contravariant type parameter and the same upper bound": {
			input: `
				class Bar
					def baz[-V < Object](a: V); end
				end
				class Foo < Bar
					def baz[-V < Object](a: V); end
				end
			`,
		},
		"override with contravariant type parameter and wider upper bound": {
			input: `
				class Bar
					def baz[-V < String](a: V); end
				end
				class Foo < Bar
					def baz[-V < Value](a: V); end
				end
			`,
		},
		"override with contravariant type parameter and narrower upper bound": {
			input: `
				class Bar
					def baz[-V < Value](a: V); end
				end
				class Foo < Bar
					def baz[-V < String](a: V); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(84, 6, 6), P(114, 6, 36)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz[-V < Std::String](a: V): void`\n  should be: `def baz[-V < Std::Value](a: V): void`\n\n  - has an incompatible type parameter, is `-V < Std::String`, should be `-V < Std::Value`"),
			},
		},

		"override with contravariant type parameter and the same lower bound": {
			input: `
				class Bar
					def baz[-V > Object](a: V); end
				end
				class Foo < Bar
					def baz[-V > Object](a: V); end
				end
			`,
		},
		"override with contravariant type parameter and wider lower bound": {
			input: `
				class Bar
					def baz[-V > String](a: V); end
				end
				class Foo < Bar
					def baz[V > Value](a: V); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 6, 6), P(113, 6, 34)), "method `Foo.:baz` is not a valid override of `Bar.:baz`\n  is:        `def baz[V > Std::Value](a: V): void`\n  should be: `def baz[-V > Std::String](a: V): void`\n\n  - has an incompatible type parameter, is `V > Std::Value`, should be `-V > Std::String`"),
			},
		},
		"override with contravariant type parameter and narrower lower bound": {
			input: `
				class Bar
					def baz[-V > Value](a: V); end
				end
				class Foo < Bar
					def baz[-V > String](a: V); end
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

func TestSpecialMethodDefinition(t *testing.T) {
	tests := testTable{
		"declare an equal method without params": {
			input: `
				class Foo
					def ==; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(30, 3, 16)), "equality operator `==` must return `bool`"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(30, 3, 16)), "equality operator `==` must accept a single parameter, got 0"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(30, 3, 16)), "method `Foo.:==` is not a valid override of `Std::Value.:==`\n  is:        `def ==(): void`\n  should be: `native def ==(other: any): bool`\n\n  - has a different return type, is `void`, should be `bool`\n  - has less parameters"),
			},
		},
		"declare an equal method without params using an alias": {
			input: `
				class Foo
					def lol; end
					alias == lol
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 4, 12), P(49, 4, 17)), "method `Foo.:lol` is not a valid override of `Std::Value.:==`\n  is:        `def lol(): void`\n  should be: `native def ==(other: any): bool`\n\n  - has a different return type, is `void`, should be `bool`\n  - has less parameters"),
				error.NewFailure(L("<main>", P(44, 4, 12), P(49, 4, 17)), "equality operator `==` must return `bool`"),
				error.NewFailure(L("<main>", P(44, 4, 12), P(49, 4, 17)), "equality operator `==` must accept a single parameter, got 0"),
			},
		},
		"declare an equal method with too many params": {
			input: `
				class Foo
					def ==(a: String, b: Int); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(49, 3, 35)), "equality operator `==` must return `bool`"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(49, 3, 35)), "equality operator `==` must accept a single parameter, got 2"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(49, 3, 35)), "method `Foo.:==` is not a valid override of `Std::Value.:==`\n  is:        `def ==(a: Std::String, b: Std::Int): void`\n  should be: `native def ==(other: any): bool`\n\n  - has a different return type, is `void`, should be `bool`\n  - has an incompatible parameter, is `a: Std::String`, should be `other: any`\n  - has an additional required parameter `b: Std::Int`"),
			},
		},
		"declare an equal method with invalid parameter type": {
			input: `
				class Foo
					def ==(a: String): bool then false
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 13), P(35, 3, 21)), "parameter `a` of equality operator `==` must be of type `any`"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(53, 3, 39)), "method `Foo.:==` is not a valid override of `Std::Value.:==`\n  is:        `def ==(a: Std::String): bool`\n  should be: `native def ==(other: any): bool`\n\n  - has an incompatible parameter, is `a: Std::String`, should be `other: any`"),
			},
		},
		"declare a valid equal method": {
			input: `
				class Foo
					def ==(other: any): bool then false
				end
			`,
		},

		"declare a lax equal method without params": {
			input: `
				class Foo
					def =~; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(30, 3, 16)), "equality operator `=~` must return `bool`"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(30, 3, 16)), "equality operator `=~` must accept a single parameter, got 0"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(30, 3, 16)), "method `Foo.:=~` is not a valid override of `Std::Value.:=~`\n  is:        `def =~(): void`\n  should be: `native def =~(other: any): bool`\n\n  - has a different return type, is `void`, should be `bool`\n  - has less parameters"),
			},
		},
		"declare a lax equal method without params using an alias": {
			input: `
				class Foo
					def lol; end
					alias =~ lol
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 4, 12), P(49, 4, 17)), "method `Foo.:lol` is not a valid override of `Std::Value.:=~`\n  is:        `def lol(): void`\n  should be: `native def =~(other: any): bool`\n\n  - has a different return type, is `void`, should be `bool`\n  - has less parameters"),
				error.NewFailure(L("<main>", P(44, 4, 12), P(49, 4, 17)), "equality operator `=~` must return `bool`"),
				error.NewFailure(L("<main>", P(44, 4, 12), P(49, 4, 17)), "equality operator `=~` must accept a single parameter, got 0"),
			},
		},
		"declare a lax equal method with too many params": {
			input: `
				class Foo
					def =~(a: String, b: Int); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(49, 3, 35)), "equality operator `=~` must return `bool`"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(49, 3, 35)), "equality operator `=~` must accept a single parameter, got 2"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(49, 3, 35)), "method `Foo.:=~` is not a valid override of `Std::Value.:=~`\n  is:        `def =~(a: Std::String, b: Std::Int): void`\n  should be: `native def =~(other: any): bool`\n\n  - has a different return type, is `void`, should be `bool`\n  - has an incompatible parameter, is `a: Std::String`, should be `other: any`\n  - has an additional required parameter `b: Std::Int`"),
			},
		},
		"declare a lax equal method with invalid parameter type": {
			input: `
				class Foo
					def =~(a: String): bool then false
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 13), P(35, 3, 21)), "parameter `a` of equality operator `=~` must be of type `any`"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(53, 3, 39)), "method `Foo.:=~` is not a valid override of `Std::Value.:=~`\n  is:        `def =~(a: Std::String): bool`\n  should be: `native def =~(other: any): bool`\n\n  - has an incompatible parameter, is `a: Std::String`, should be `other: any`"),
			},
		},
		"declare a valid lax equal method": {
			input: `
				class Foo
					def =~(other: any): bool then false
				end
			`,
		},

		"declare a relational operator method without params": {
			input: `
				class Foo
					def <; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(29, 3, 15)), "relational operator `<` must return `bool`"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(29, 3, 15)), "relational operator `<` must accept a single parameter, got 0"),
			},
		},
		"declare a relational operator method without params using an alias": {
			input: `
				class Foo
					def lol; end
					alias < lol
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 4, 12), P(48, 4, 16)), "relational operator `<` must return `bool`"),
				error.NewFailure(L("<main>", P(44, 4, 12), P(48, 4, 16)), "relational operator `<` must accept a single parameter, got 0"),
			},
		},
		"declare a relational operator method with too many params": {
			input: `
				class Foo
					def >(a: String, b: Int); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(48, 3, 34)), "relational operator `>` must return `bool`"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(48, 3, 34)), "relational operator `>` must accept a single parameter, got 2"),
			},
		},
		"declare a relational operator method with invalid parameter type": {
			input: `
				class Foo
					def <=(a: String): bool then false
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 13), P(35, 3, 21)), "parameter `a` of relational operator `<=` must accept `Foo`"),
			},
		},
		"declare a valid relational operator method": {
			input: `
				class Foo
					def >=(other: Foo): bool then false
				end
			`,
		},
		"declare a valid relational operator method with wider param type": {
			input: `
				class Foo
					def >=(other: any): bool then false
				end
			`,
		},
		"declare a relational operator method in an interface": {
			input: `
				interface Foo
					sig <=(a: String): bool
				end
			`,
		},

		"declare a binary operator without params": {
			input: `
				class Foo
					def +; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(29, 3, 15)), "method `+` cannot be void"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(29, 3, 15)), "method `+` must define exactly 1 parameters, got 0"),
			},
		},
		"declare a binary operator without params using an alias": {
			input: `
				class Foo
					def lol; end
					alias + lol
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 4, 12), P(48, 4, 16)), "method `+` cannot be void"),
				error.NewFailure(L("<main>", P(44, 4, 12), P(48, 4, 16)), "method `+` must define exactly 1 parameters, got 0"),
			},
		},
		"declare a binary operator with too many params": {
			input: `
				class Foo
					def -(a: String, b: Int); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(48, 3, 34)), "method `-` cannot be void"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(48, 3, 34)), "method `-` must define exactly 1 parameters, got 2"),
			},
		},
		"declare a valid binary operator": {
			input: `
				class Foo
					def *(other: String): String then other
				end
			`,
		},

		"declare an increment method without params and return type": {
			input: `
				class Foo
					def ++; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(30, 3, 16)), "method `++` cannot be void"),
			},
		},
		"declare an increment method without params and return type using an alias": {
			input: `
				class Foo
					def lol; end
					alias ++ lol
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 4, 12), P(49, 4, 17)), "method `++` cannot be void"),
			},
		},
		"declare a decrement method with too many params": {
			input: `
				class Foo
					def --(a: String, b: Int); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(49, 3, 35)), "method `--` cannot be void"),
				error.NewFailure(L("<main>", P(20, 3, 6), P(49, 3, 35)), "method `--` must define exactly 0 parameters, got 2"),
			},
		},
		"declare a negate method with too many params": {
			input: `
				class Foo
					def -@(a: String): bool then false
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(53, 3, 39)), "method `-@` must define exactly 0 parameters, got 1"),
			},
		},
		"declare a unary plus method with too many params": {
			input: `
				class Foo
					def +@(a: String): bool then false
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(53, 3, 39)), "method `+@` must define exactly 0 parameters, got 1"),
			},
		},
		"declare a bitwise not method with too many params": {
			input: `
				class Foo
					def ~(a: String): bool then false
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(52, 3, 38)), "method `~` must define exactly 0 parameters, got 1"),
			},
		},
		"declare a valid unary method": {
			input: `
				class Foo
					def ++: bool then false
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
		"positional rest params have tuple types": {
			input: `
				module Foo
					def baz(*b: Float)
						var c: nil = b
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(59, 4, 20), P(59, 4, 20)), "type `Std::Tuple[Std::Float]` cannot be assigned to type `nil`"),
			},
		},
		"named rest params have record types": {
			input: `
				module Foo
					def baz(**b: Float)
						var c: nil = b
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 4, 20), P(60, 4, 20)), "type `Std::Record[Std::Symbol, Std::Float]` cannot be assigned to type `nil`"),
			},
		},
		"declare with type parameters": {
			input: `
				def foo[V](a: V): V
					a
				end
			`,
		},
		"typecheck bounds of type params": {
			input: `
				def foo[V < Foo](a: V): V
					a
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 2, 17), P(19, 2, 19)), "undefined type `Foo`"),
			},
		},
		"call methods on type params with upper bounds": {
			input: `
				def foo[V < Int](a: V): String
					a.to_string
				end
			`,
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
		"declare a generator": {
			input: `
				def *foo: String
					"lol"
				end
			`,
		},
		"declare a generator with a throw type": {
			input: `
				def *foo: String ! Int
					"lol"
				end
			`,
		},
		"parameters are values in generators": {
			input: `
				def *foo(a: Int): String ! Int
					a = 5
					"lol"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(41, 3, 6), P(41, 3, 6)), "local value `a` cannot be reassigned"),
			},
		},
		"redeclare the method in the same class with incompatible signature": {
			input: `
				class Foo
					def baz(a: Int): String then "a"
					def baz(): void; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(58, 4, 6), P(77, 4, 25)), "method `Foo.:baz` is not a valid override of `Foo.:baz`\n  is:        `def baz(): void`\n  should be: `def baz(a: Std::Int): Std::String`\n\n  - has a different return type, is `void`, should be `Std::String`\n  - has less parameters"),
			},
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
		"declare sig with type parameters": {
			input: `
				interface Bar
					sig foo[V](a: V): V
				end
			`,
		},
		"typecheck bounds of type params in sig": {
			input: `
				interface Bar
					sig foo[V < Foo](a: V): V
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 18), P(38, 3, 20)), "undefined type `Foo`"),
			},
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
				  var @a: String?
					def baz(@a); end
				end
			`,
		},
		"explicit instance variable parameter type": {
			input: `
				class Foo
				  var @a: String?
					def baz(@a: String?); end
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
				  var @a: String?
					def baz(@a: String | Float | nil); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(54, 4, 18), P(73, 4, 37)), "type `Std::String | Std::Float | nil` cannot be assigned to instance variable `@a` of type `Std::String?`"),
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
					def baz(@a: String?)
						var b: String? = @a
					end
				end
			`,
		},
		"instance variable parameter declares an instance variable that can be redeclared with the same type": {
			input: `
				class Foo
					def baz(@a: String?)
						var b: String? = @a
					end

					var @a: String?
				end
			`,
		},
		"instance variable parameter declares an instance variable tha cannot be redeclared with a different type": {
			input: `
				class Foo
					def baz(@a: String?)
						var b: String? = @a
					end

					var @a: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(82, 7, 6), P(92, 7, 16)), "cannot redeclare instance variable `@a` with a different type, is `Std::Int`, should be `Std::String?`, previous definition found in `Foo`"),
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
		"call a variable": {
			input: `
				module Foo
					def call(a: Int): Int then a
				end
				a := Foo
				a(1)
			`,
		},
		"call a variable instead of a method": {
			input: `
				def a: Int then 20

				module Foo
					def call(a: Int): Int then a
				end
				a := Foo
				a(1)
			`,
		},
		"call a variable without call method": {
			input: `
				a := 5
				a(1)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 3, 5), P(19, 3, 8)), "method `call` is not defined on type `Std::Int`"),
			},
		},
		"call to a generator returns a generator": {
			input: `
				module Foo
					def* baz(a: Int): Int
						yield 5

						10
					end
				end
				var a: nil = Foo.baz(5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(101, 9, 18), P(110, 9, 27)), "type `Std::Generator[Std::Int, never]` cannot be assigned to type `nil`"),
			},
		},
		"call to a generator returns a generator that throws": {
			input: `
				module Foo
					def* baz(a: Int): Int ! Error
						yield 5

						10
					end
				end
				var a: nil = Foo.baz(5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(109, 9, 18), P(118, 9, 27)), "type `Std::Generator[Std::Int, Std::Error]` cannot be assigned to type `nil`"),
			},
		},
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

		"call a method on a type parameter": {
			input: `
				class Bar[V]
					init(@a: V); end

					def bar
						@a.foo()
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 6, 7), P(67, 6, 14)), "method `foo` is not defined on type `any`"),
			},
		},
		"call a method on a type parameter with an upper bound": {
			input: `
				class Foo
					def foo; end
				end

				class Bar[V < Foo]
					init(@a: V); end

					def bar
						@a.foo()
					end
				end
			`,
		},
		"call a method on a type parameter with a lower bound": {
			input: `
				class Foo
					def foo; end
				end

				class Bar[V > Foo]
					init(@a: V); end

					def bar
						@a.foo()
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(107, 10, 7), P(114, 10, 14)), "method `foo` is not defined on type `any`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestGenericMethodCalls(t *testing.T) {
	tests := testTable{
		"infer type parameter from closure's return type": {
			input: `
				var a: 9 = HashMap::[String, Int]().map_pairs() |pair| -> Pair(1u8, 1.2)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 2, 16), P(76, 2, 76)), "type `Std::HashMap[Std::UInt8, Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"infer type parameter from closure's return type with different param names": {
			input: `
				var a: 9 = HashMap::[String, Int]().map_pairs() |p| -> Pair(1u8, 1.2)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 2, 16), P(73, 2, 73)), "type `Std::HashMap[Std::UInt8, Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"call a generic method with explicit type arguments": {
			input: `
				module Foo
					def baz[V](a: V): V then a
				end
				var a: Int = Foo.baz::[Int](5)
			`,
		},
		"call a generic method with explicit type argument that satisfies the upper bound": {
			input: `
				class Bar; end
				module Foo
					def baz[V < Bar](a: V): V then a
				end
				var a: Bar = Foo.baz::[Bar](Bar())
			`,
		},
		"call a generic method with explicit type argument that does not satisfy the upper bound": {
			input: `
				class Bar; end
				module Foo
					def baz[V < Bar](a: V): V then a
				end
				var a: Int = Foo.baz::[Int](5)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(108, 6, 28), P(110, 6, 30)), "type `Std::Int` does not satisfy the upper bound `Bar`"),
			},
		},
		"call a generic method with explicit type argument that satisfies the lower bound": {
			input: `
				class Bar; end
				module Foo
					def baz[V > Bar](a: V): V then a
				end
				var a: Bar = Foo.baz::[Object](Object())
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(98, 6, 18), P(124, 6, 44)), "type `Std::Object` cannot be assigned to type `Bar`"),
			},
		},
		"call a generic method with explicit type argument that does not satisfy the lower bound": {
			input: `
				class Bar; end
				module Foo
					def baz[V > Bar](a: V): V then a
				end
				var a: Int = Foo.baz::[3](3)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(108, 6, 28), P(108, 6, 28)), "type `3` does not satisfy the lower bound `Bar`"),
			},
		},
		"call a non-generic method with explicit type argument": {
			input: `
				class Bar; end
				module Foo
					def baz(a: String): String then a
				end
				var a: Int = Foo.baz::[3](3)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(99, 6, 18), P(113, 6, 32)), "`Foo::baz` requires 0 type argument(s), got: 1"),
			},
		},

		"infer type argument incompatible with upper bound": {
			input: `
				module Foo
					def foo[V < Int](a: V): V then a
				end
				var a: 9 = Foo.foo("foo")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 5, 24), P(89, 5, 28)), "type `\"foo\"` does not satisfy the upper bound `Std::Int`"),
				error.NewFailure(L("<main>", P(77, 5, 16), P(90, 5, 29)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"infer type argument from upper bound": {
			input: `
				module Foo
					def foo[V < Int]: V then loop; end
				end
				var a: 9 = Foo.foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(79, 5, 16), P(87, 5, 24)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"infer type argument from upper bound with type parameters": {
			input: `
				module Foo
					def foo[V < Comparable[V]]: V then loop; end
				end
				var a: 9 = Foo.foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 29), P(44, 3, 29)), "undefined type `V`"),
				error.NewFailure(L("<main>", P(89, 5, 16), P(97, 5, 24)), "type `Std::Comparable[untyped]` cannot be assigned to type `9`"),
			},
		},
		"infer type argument from lower bound": {
			input: `
				module Foo
					def foo[V > Int]: V then loop; end
				end
				var a: 9 = Foo.foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(79, 5, 16), P(87, 5, 24)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"infer type argument from lower bound with type parameters": {
			input: `
				module Foo
					def foo[V > Comparable[V]]: V then loop; end
				end
				var a: 9 = Foo.foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 29), P(44, 3, 29)), "undefined type `V`"),
				error.NewFailure(L("<main>", P(89, 5, 16), P(97, 5, 24)), "type `Std::Comparable[untyped]` cannot be assigned to type `9`"),
			},
		},
		"infer type argument incompatible with lower bound": {
			input: `
				module Foo
					def foo[V > Int](a: V): V then a
				end
				var a: 9 = Foo.foo("foo")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 5, 24), P(89, 5, 28)), "type `\"foo\"` does not satisfy the lower bound `Std::Int`"),
				error.NewFailure(L("<main>", P(77, 5, 16), P(90, 5, 29)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"infer simple type argument": {
			input: `
				module Foo
					def foo[V](a: V): V then a
				end
				var a: 9 = Foo.foo("foo")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(71, 5, 16), P(84, 5, 29)), "type `Std::String` cannot be assigned to type `9`"),
			},
		},
		"infer based on first argument": {
			input: `
				module Foo
					def foo[V](a: V, b: V): V then b
				end
				var b: 9 = Foo.foo("foo", 2)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(92, 5, 31), P(92, 5, 31)), "expected type `Std::String` for parameter `b` in call to `foo`, got type `2`"),
				error.NewFailure(L("<main>", P(77, 5, 16), P(93, 5, 32)), "type `Std::String` cannot be assigned to type `9`"),
			},
		},
		"infer two type arguments": {
			input: `
				module Foo
					def foo[V, T](a: V | T): V | T then a
				end
				var a: Int | Float = 2
				var b: 9 = Foo.foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(109, 6, 16), P(118, 6, 25)), "type `Std::Int | Std::Float` cannot be assigned to type `9`"),
			},
		},

		"call a generic receiverless method with explicit type arguments": {
			input: `
				module Foo
					def baz[V](a: V): V then a
					var a: Int = baz::[Int](5)
				end
			`,
		},
		"call a generic receiverless method with explicit type argument that satisfies the upper bound": {
			input: `
				class Bar; end
				module Foo
					def baz[V < Bar](a: V): V then a
					var a: Bar = baz::[Bar](Bar())
				end
			`,
		},
		"call a generic receiverless method with explicit type argument that does not satisfy the upper bound": {
			input: `
				class Bar; end
				module Foo
					def baz[V < Bar](a: V): V then a
					var a: Int = baz::[Int](5)
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(97, 5, 25), P(99, 5, 27)), "type `Std::Int` does not satisfy the upper bound `Bar`"),
			},
		},
		"call a generic receiverless method with explicit type argument that satisfies the lower bound": {
			input: `
				class Bar; end
				module Foo
					def baz[V > Bar](a: V): V then a
					var a: Bar = baz::[Object](Object())
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(91, 5, 19), P(113, 5, 41)), "type `Std::Object` cannot be assigned to type `Bar`"),
			},
		},
		"call a generic receiverless method with explicit type argument that does not satisfy the lower bound": {
			input: `
				class Bar; end
				module Foo
					def baz[V > Bar](a: V): V then a
					var a: Int = baz::[3](3)
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(97, 5, 25), P(97, 5, 25)), "type `3` does not satisfy the lower bound `Bar`"),
			},
		},
		"call a non-generic receiverless method with explicit type argument": {
			input: `
				class Bar; end
				module Foo
					def baz(a: String): String then a
					var a: Int = baz::[3](3)
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(92, 5, 19), P(102, 5, 29)), "`Foo::baz` requires 0 type argument(s), got: 1"),
			},
		},

		"infer simple type argument in receiverless call": {
			input: `
				module Foo
					def foo[V](a: V): V then a
					var a: 9 = foo("foo")
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(64, 4, 17), P(73, 4, 26)), "type `Std::String` cannot be assigned to type `9`"),
			},
		},
		"infer based on first argument in receiverless call": {
			input: `
				module Foo
					def foo[V](a: V, b: V): V then b
					var b: 9 = foo("foo", 2)
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(81, 4, 28), P(81, 4, 28)), "expected type `Std::String` for parameter `b` in call to `foo`, got type `2`"),
				error.NewFailure(L("<main>", P(70, 4, 17), P(82, 4, 29)), "type `Std::String` cannot be assigned to type `9`"),
			},
		},
		"infer two type arguments in receiverless call": {
			input: `
				module Foo
					def foo[V, T](a: V | T): V | T then a
					var a: Int | Float = 2
					var b: 9 = foo(a)
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(103, 5, 17), P(108, 5, 22)), "type `Std::Int | Std::Float` cannot be assigned to type `9`"),
			},
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
		"instantiate a noinit class": {
			input: `
				noinit class Foo; end
				Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(31, 3, 5), P(35, 3, 9)), "cannot instantiate class `Foo` marked as `noinit`"),
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
		"instantiate a generic class": {
			input: `
				class Foo[V]
					init(a: V); end
				end
				Foo::[String]("foo")
			`,
		},
		"instantiate a generic class with valid type arguments": {
			input: `
				class Foo[V]
					init(a: V); end
				end
				Foo::[String]("foo")
			`,
		},
		"instantiate a generic class with valid type arguments and invalid arguments": {
			input: `
				class Foo[V]
					init(a: V); end
				end
				Foo::[String](1)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(65, 5, 19), P(65, 5, 19)), "expected type `Std::String` for parameter `a` in call to `#init`, got type `1`"),
			},
		},
		"instantiate a generic class with satisfied upper bound": {
			input: `
				class Bar; end
				class Baz < Bar; end

				class Foo[V < Bar]
					init(a: V); end
				end
				Foo::[Bar](Bar())
				Foo::[Baz](Baz())
			`,
		},
		"instantiate a generic class with unsatisfied upper bound": {
			input: `
				class Bar; end
				class Baz < Bar; end

				class Foo[V < Bar]
					init(a: V); end
				end
				Foo::[String]("foo")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(108, 8, 11), P(113, 8, 16)), "type `Std::String` does not satisfy the upper bound `Bar`"),
			},
		},
		"instantiate a generic class with satisfied lower bound": {
			input: `
				class Bar; end
				class Baz < Bar; end

				class Foo[V > Baz]
					init(a: V); end
				end
				Foo::[Bar](Bar())
				Foo::[Object](Object())
			`,
		},
		"instantiate a generic class with unsatisfied lower bound": {
			input: `
				class Bar; end
				class Baz < Bar; end

				class Foo[V > Baz]
					init(a: V); end
				end
				Foo::[Int](1)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(108, 8, 11), P(110, 8, 13)), "type `Std::Int` does not satisfy the lower bound `Baz`"),
			},
		},

		"new - instantiate a class without a constructor": {
			input: `
				class Foo
					def foo: self
						new
					end
				end
			`,
		},
		"new - instantiate an abstract class": {
			input: `
				abstract class Foo
					def foo: self
						new
					end
				end
			`,
		},
		"new - instantiate a class with a constructor": {
			input: `
				class Foo
					init(a: Int); end

					def foo: self
						new(1)
					end
				end
			`,
		},
		"new - instantiate a class with a constructor with a wrong type": {
			input: `
				class Foo
					init(a: String); end
					def foo: self
						new(1)
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(70, 5, 11), P(70, 5, 11)), "expected type `Std::String` for parameter `a` in call to `#init`, got type `1`"),
			},
		},
		"new - instantiate a class with an inherited constructor": {
			input: `
				class Bar
					init(a: Int); end
				end

				class Foo < Bar
					def foo: self
						new(1)
					end
				end
			`,
		},
		"new - instantiate a class with an inherited constructor with a wrong type": {
			input: `
				class Bar
					init(a: String); end
				end

				class Foo < Bar
					def foo: self
						new(1)
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(99, 8, 11), P(99, 8, 11)), "expected type `Std::String` for parameter `a` in call to `#init`, got type `1`"),
			},
		},
		"new - call a method on an instantiated instance": {
			input: `
				class Foo
					init(a: String); end

					def foo
						a := new("foo")
						a.bar
					end

					def bar; end
				end
			`,
		},
		"new - instantiate a generic class": {
			input: `
				class Foo[V]
					init(@a: V); end

					def foo: self
						new(@a)
					end
				end
			`,
		},
		"new - instantiate a generic class with invalid arguments": {
			input: `
				class Foo[V]
					init(a: V); end

					def foo: self
						new("foo")
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(69, 6, 11), P(73, 6, 15)), "expected type `V` for parameter `a` in call to `#init`, got type `\"foo\"`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestConstructorCallInference(t *testing.T) {
	tests := testTable{
		"infer type argument with empty constructor": {
			input: `
				class Foo[V]; end
				var a: 9 = Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(38, 3, 16), P(42, 3, 20)), "type `Foo[any]` cannot be assigned to type `9`"),
			},
		},
		"infer type argument based on upper bound": {
			input: `
				class Foo[V < String]; end
				var a: 9 = Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(47, 3, 16), P(51, 3, 20)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer type argument based on upper bound with type parameters": {
			input: `
				class Foo[V < Comparable[V]]; end
				var a: 9 = Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(54, 3, 16), P(58, 3, 20)), "cannot infer type argument for `V` in call to `#init`"),
				error.NewFailure(L("<main>", P(54, 3, 16), P(58, 3, 20)), "type `Foo[untyped]` cannot be assigned to type `9`"),
			},
		},
		"infer type argument based on lower bound": {
			input: `
				class Foo[V > String]; end
				var a: 9 = Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(47, 3, 16), P(51, 3, 20)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer type argument based on lower bound with type parameters": {
			input: `
				class Foo[V > Comparable[V]]; end
				var a: 9 = Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(54, 3, 16), P(58, 3, 20)), "type `Foo[any]` cannot be assigned to type `9`"),
			},
		},
		"infer type argument based on lower bound and upper bound": {
			input: `
				class Bar; end
				class Foo[V > Bar < Object]; end
				var a: 9 = Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(72, 4, 16), P(76, 4, 20)), "type `Foo[Bar]` cannot be assigned to type `9`"),
			},
		},
		"infer type argument incompatible with upper bound": {
			input: `
				class Foo[V < String]
					init(a: V); end
				end
				var a: 9 = Foo(9)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(75, 5, 20), P(75, 5, 20)), "type `9` does not satisfy the upper bound `Std::String`"),
				error.NewFailure(L("<main>", P(71, 5, 16), P(76, 5, 21)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer type argument incompatible with lower bound": {
			input: `
				class Foo[V > String]
					init(a: V); end
				end
				var a: 9 = Foo(9)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(75, 5, 20), P(75, 5, 20)), "type `9` does not satisfy the lower bound `Std::String`"),
				error.NewFailure(L("<main>", P(71, 5, 16), P(76, 5, 21)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer simple type argument": {
			input: `
				class Foo[V]
					init(a: V); end
				end
				var a: 9 = Foo("foo")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(62, 5, 16), P(71, 5, 25)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer based on first argument": {
			input: `
				class Foo[V]
					init(a: V, b: V); end
				end
				var b: 9 = Foo("foo", 2)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(79, 5, 27), P(79, 5, 27)), "expected type `Std::String` for parameter `b` in call to `#init`, got type `2`"),
				error.NewFailure(L("<main>", P(68, 5, 16), P(80, 5, 28)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"infer two type arguments": {
			input: `
				class Foo[V, T]
					init(a: V | T); end
				end
				var a: Int | Float = 2
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(96, 6, 16), P(101, 6, 21)), "type `Foo[Std::Int | Std::Float, Std::Int | Std::Float]` cannot be assigned to type `9`"),
			},
		},

		"param is a union, argument is an exact match": {
			input: `
				class Foo[V]
					init(a: V | Int); end
				end
				var a: Int | String = "lol"
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(100, 6, 16), P(105, 6, 21)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"param is a union, argument is a union with additional types": {
			input: `
				class Foo[V]
					init(a: V | Int); end
				end
				var a: Int | String | Float = "lol"
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(108, 6, 16), P(113, 6, 21)), "type `Foo[Std::String | Std::Float]` cannot be assigned to type `9`"),
			},
		},
		"param is a union, argument is not": {
			input: `
				class Foo[V]
					init(a: V | Int); end
				end
				var b: 9 = Foo("foo")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(68, 5, 16), P(77, 5, 25)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},

		"param is nilable, argument is not": {
			input: `
				class Foo[V]
					init(a: V?); end
				end
				var b: 9 = Foo("foo")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(63, 5, 16), P(72, 5, 25)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"param and argument are nilable": {
			input: `
				class Foo[V]
					init(a: V?); end
				end
				var a: String? = "foo"
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(90, 6, 16), P(95, 6, 21)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"param is nilable, argument is a union with nil": {
			input: `
				class Foo[V]
					init(a: V?); end
				end
				var a: String | nil = "foo"
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(95, 6, 16), P(100, 6, 21)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"param is nilable, argument is a union with Nil": {
			input: `
				class Foo[V]
					init(a: V?); end
				end
				var a: String | Nil = "foo"
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(95, 6, 16), P(100, 6, 21)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"param is nilable, argument is a broad union with Nil": {
			input: `
				class Foo[V]
					init(a: V?); end
				end
				var a: String | Int | Nil = "foo"
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(101, 6, 16), P(106, 6, 21)), "type `Foo[Std::String | Std::Int]` cannot be assigned to type `9`"),
			},
		},

		"param is a generic, argument is also": {
			input: `
				class Bar[V]
					init(a: V); end
				end

				class Foo[V]
					init(a: Bar[V]); end
				end
				a := Bar::[String]("foo")
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(144, 10, 16), P(149, 10, 21)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"param is a generic, argument is a subtype": {
			input: `
				class Bar[V]
					init(a: V); end
				end
				class Baz[V] < Bar[V]; end

				class Foo[V]
					init(a: Bar[V]); end
				end
				a := Baz::[String]("foo")
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(175, 11, 16), P(180, 11, 21)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"param is a generic, argument is a subtype with more type parameters": {
			input: `
				class Bar[V]
					init(a: V); end
				end
				class Baz[V, T] < Bar[V]; end

				class Foo[V]
					init(a: Bar[V]); end
				end
				a := Baz::[Int, Float](1)
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(178, 11, 16), P(183, 11, 21)), "type `Foo[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"param is a generic, argument is not": {
			input: `
				class Bar[V]
					init(a: V); end
				end

				class Foo[V]
					init(a: Bar[V]); end
				end
				var b: 9 = Foo(1)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(118, 9, 20), P(118, 9, 20)), "expected type `Bar[V]` for parameter `a` in call to `#init`, got type `1`"),
				error.NewFailure(L("<main>", P(114, 9, 16), P(119, 9, 21)), "type `Foo[any]` cannot be assigned to type `9`"),
			},
		},

		"param is a singleton, argument is also": {
			input: `
				class Foo[V < Value]
					init(a: &V); end
				end
				var b: 9 = Foo(String)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(71, 5, 16), P(81, 5, 26)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"param is a singleton, argument is not": {
			input: `
				class Foo[V < Value]
					init(a: &V); end
				end
				class Bar; end
				var b: 9 = Foo("")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(94, 6, 20), P(95, 6, 21)), "expected type `&V` for parameter `a` in call to `#init`, got type `\"\"`"),
				error.NewFailure(L("<main>", P(90, 6, 16), P(96, 6, 22)), "type `Foo[Std::Value]` cannot be assigned to type `9`"),
			},
		},

		"param is an instance, argument is a class instance": {
			input: `
				class Foo[V < Class]
					init(a: ^V); end
				end
				var b: 9 = Foo("")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(71, 5, 16), P(77, 5, 22)), "type `Foo[&Std::String]` cannot be assigned to type `9`"),
			},
		},
		"param is an instance, argument is not": {
			input: `
				class Foo[V < Class]
					init(a: ^V); end
				end
				var b: 9 = Foo(String)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(75, 5, 20), P(80, 5, 25)), "expected type `^V` for parameter `a` in call to `#init`, got type `&Std::String`"),
				error.NewFailure(L("<main>", P(71, 5, 16), P(81, 5, 26)), "type `Foo[Std::Class]` cannot be assigned to type `9`"),
			},
		},

		"param is a not type, argument is also": {
			input: `
				class Foo[V]
					init(a: ~V); end
				end
				var a: ~Int = 2.9
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 6, 16), P(90, 6, 21)), "type `Foo[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"param is a not type, argument is not": {
			input: `
				class Foo[V]
					init(a: ~V); end
				end
				var b: 9 = Foo(2.9)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(67, 5, 20), P(69, 5, 22)), "expected type `~V` for parameter `a` in call to `#init`, got type `2.9`"),
				error.NewFailure(L("<main>", P(63, 5, 16), P(70, 5, 23)), "type `Foo[any]` cannot be assigned to type `9`"),
			},
		},

		"param is an intersection, argument is also": {
			input: `
				class Foo[V]
					init(a: V & StringConvertible); end
				end
				var a: Int & StringConvertible = 2
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(121, 6, 16), P(126, 6, 21)), "type `Foo[Std::Int]` cannot be assigned to type `9`"),
			},
		},

		"param is a closure with a type param in one param, argument is also": {
			input: `
				class Foo[V]
					init(a: |a: V|: void); end
				end
				a := |a: Int| -> a
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(96, 6, 16), P(101, 6, 21)), "type `Foo[Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"param is a closure with a type param in the return type, argument is also": {
			input: `
				class Foo[V]
					init(a: |a: Int|: V); end
				end
				a := |a: Int|: String -> a.to_string
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(113, 6, 16), P(118, 6, 21)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"param is a closure with a type param in the throw type, argument is also": {
			input: `
				class Foo[V]
					init(a: |a: Int| ! V); end
				end
				a := |a: Int| ! String -> a.to_string
				var b: 9 = Foo(a)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(115, 6, 16), P(120, 6, 21)), "type `Foo[Std::String]` cannot be assigned to type `9`"),
			},
		},
		"param is a closure with a type param in one param, argument is not": {
			input: `
				class Foo[V]
					init(a: |a: V|: void); end
				end
				var b: 9 = Foo("foo")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(77, 5, 20), P(81, 5, 24)), "type `Std::String` does not implement closure `|a: V|: void`:\n\n  - missing method `call` with signature: `def call(a: V): void`\n"),
				error.NewFailure(L("<main>", P(77, 5, 20), P(81, 5, 24)), "expected type `|a: V|: void` for parameter `a` in call to `#init`, got type `\"foo\"`"),
				error.NewFailure(L("<main>", P(73, 5, 16), P(82, 5, 25)), "type `Foo[any]` cannot be assigned to type `9`"),
			},
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
		"call a method inherited from generic superclass": {
			input: `
				class Foo[V]
					def baz(a: V): V then a
				end

				class Bar < Foo[Int]; end
				var bar = Bar()
				bar.baz(5)
			`,
		},
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
		"call a method inherited from generic mixin": {
			input: `
				mixin Bar[V]
					def baz(a: V): V then a
				end

				class Foo
					include Bar[Int]
				end

				var foo = Foo()
				foo.baz(5)
			`,
		},
		"call a method inherited from generic mixin included in mixin": {
			input: `
				mixin Bar[V]
					def bar(a: V): V then a
				end

				mixin Baz
					include Bar[Int]
				end

				class Foo
					include Baz
				end

				var foo = Foo()
				foo.bar(5)
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
		"call method inherited from mixin that was included to the parent mixin afterwards": {
			input: `
				mixin Foo
					def foo; end
				end

				class Bar
					include Foo
				end

				mixin Fuzz
					def fuzz; end
				end

				mixin Hejo
					include Fuzz
					def hejo; end
				end

				mixin Foo
					include Hejo
				end

				b := Bar()
				b.foo
				b.hejo
				b.fuzz
			`,
		},
		"call method inherited from interface that was implemented to the parent interface afterwards": {
			input: `
				interface Foo
					def foo; end
				end

				interface Bar
					implement Foo
				end

				interface Fuzz
					def fuzz; end
				end

				interface Hejo
					implement Fuzz
					def hejo; end
				end

				interface Foo
					implement Hejo
				end

				class Elo
					def foo; end
					def hejo; end
					def fuzz; end
				end

				var b: Bar = Elo()
				b.foo
				b.hejo
				b.fuzz
			`,
		},
		"call method from extend where when the conditions are satisfied": {
			input: `
				class Foo[T]
					extend where T < String
						def bar(a: T); end
					end
				end

				f := Foo::[String]()
				f.bar("lol")
			`,
		},
		"call method from extend where in parent mixin when the conditions are satisfied": {
			input: `
				mixin Foo[T]
					extend where T < String
						def foo(a: T); end
					end
				end

				class Bar
					include Foo[String]
				end

				b := Bar()
				b.foo("lol")
			`,
		},
		"call method from extend where when the conditions are not satisfied": {
			input: `
				class Foo[T]
					extend where T < String
						def bar(a: T); end
					end
				end

				f := Foo::[Int]()
				f.bar("lol")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(116, 9, 5), P(127, 9, 16)), "method `bar` is not defined on type `Foo[Std::Int]`"),
			},
		},
		"call method from extend where in parent mixin when the conditions are not satisfied": {
			input: `
				mixin Foo[T]
					extend where T < String
						def foo(a: T); end
					end
				end

				class Bar
					include Foo[Int]
				end

				b := Bar()
				b.foo("lol")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(154, 13, 5), P(165, 13, 16)), "method `foo` is not defined on type `Bar`"),
			},
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
		"use variables from outer scope": {
			input: `
				a := 6
				b := |c: Int|: Int -> c + a
			`,
		},
		"infer return type": {
			input: `
				var a: 8 = |a: Int| -> 9.2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 2, 16), P(30, 2, 30)), "type `|a: Std::Int|: 9.2` cannot be assigned to type `8`"),
			},
		},
		"infer throw type": {
			input: `
				var a: 8 = |a: Int| ->
					throw "dupa" if a == 5

					a * 2
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 2, 16), P(74, 6, 7)), "type `|a: Std::Int|: Std::Int ! \"dupa\"` cannot be assigned to type `8`"),
			},
		},
		"infer return type in multiline closure": {
			input: `
				var a: 8 = |a: Int| ->
					return 9.2 if a == 9

					nil
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(16, 2, 16), P(70, 6, 7)), "type `|a: Std::Int|: 9.2 | nil` cannot be assigned to type `8`"),
			},
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
		"invalid return value": {
			input: `
				a := |a: Int = 2|: String -> 5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 2, 34), P(34, 2, 34)), "type `5` cannot be assigned to type `Std::String`"),
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
				error.NewFailure(L("<main>", P(37, 3, 9), P(37, 3, 9)), "type `Std::Int` does not implement closure `|a: Std::Int|: Std::Int`:\n\n  - missing method `call` with signature: `def call(a: Std::Int): Std::Int`\n"),
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
				error.NewFailure(L("<main>", P(37, 3, 9), P(78, 3, 50)), "type `|a: Std::Float, b?: Std::String|: Std::String` does not implement closure `|a: Std::Int|: Std::Int`:\n\n  - incorrect implementation of `call`\n      is:        `def call(a: Std::Float, b?: Std::String): Std::String`\n      should be: `def call(a: Std::Int): Std::Int`\n"),
				error.NewFailure(L("<main>", P(37, 3, 9), P(78, 3, 50)), "type `|a: Std::Float, b?: Std::String|: Std::String` cannot be assigned to type `|a: Std::Int|: Std::Int`"),
			},
		},
		"take param and return types from closure defined as a method parameter": {
			input: `
				def foo(fn: |a: String|: Int); end
				foo() |a| ->
					var b: nil = a
					5
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(75, 4, 19), P(75, 4, 19)), "type `Std::String` cannot be assigned to type `nil`"),
			},
		},
		"invalid closure argument": {
			input: `
				def foo(fn: |a: String|: Int); end
				foo() |i| -> 2.5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(50, 3, 11), P(59, 3, 20)), "type `|i: Std::String|: 2.5` does not implement closure `|a: Std::String|: Std::Int`:\n\n  - incorrect implementation of `call`\n      is:        `def call(i: Std::String): 2.5`\n      should be: `def call(a: Std::String): Std::Int`\n"),
				error.NewFailure(L("<main>", P(50, 3, 11), P(59, 3, 20)), "expected type `|a: Std::String|: Std::Int` for parameter `fn` in call to `foo`, got type `|i: Std::String|: 2.5`"),
			},
		},
		"accept closure argument with different param names": {
			input: `
				def foo(fn: |a: String|: Float); end
				foo() |i| ->
					var b: nil = i
					2.5
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(77, 4, 19), P(77, 4, 19)), "type `Std::String` cannot be assigned to type `nil`"),
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
		"closure call returns the specified type": {
			input: `
				a := |a: Int|: Int -> a
				var b: nil = a.(2)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(46, 3, 18), P(50, 3, 22)), "type `Std::Int` cannot be assigned to type `nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
