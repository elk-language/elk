package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestSingleton(t *testing.T) {
	tests := testTable{
		"has its own local scope": {
			input: `
				class Foo
					a := 5
					singleton
						a
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 5, 7), P(48, 5, 7)), "undefined local `a`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestModule(t *testing.T) {
	tests := testTable{
		"has its own local scope": {
			input: `
				a := 5
				module Foo
					a
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(32, 4, 6), P(32, 4, 6)), "undefined local `a`"),
			},
		},
		"module with public constant": {
			input: `module Foo; end`,
		},
		"module with conflicting constant with Std": {
			input: `module Int; end`,
		},
		"module with private constant": {
			input: `module _Fo; end`,
		},
		"module with simple constant lookup": {
			input: `module Std::Foo; end`,
		},
		"module with non obvious constant lookup": {
			input: `module Int::Foo; end`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(7, 1, 8), P(9, 1, 10)), "undefined namespace `Int`"),
			},
		},
		"resolve module with non obvious constant lookup": {
			input: `
				module Int
				  module Foo; end
				end
			  ::Int::Foo
			`,
		},
		"module with undefined root constant": {
			input: `module Foo::Bar; end`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(7, 1, 8), P(9, 1, 10)), "undefined namespace `Foo`"),
			},
		},
		"module with undefined constant in the middle": {
			input: `module Std::Foo::Bar; end`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 1, 13), P(14, 1, 15)), "undefined namespace `Std::Foo`"),
			},
		},
		"nested modules": {
			input: `
				module Foo
					module Bar; end
				end
			`,
		},
		"resolve constant inside of new module": {
			input: `
				module Foo
					module Bar; end
					Bar
				end
			`,
		},
		"resolve constant outside of new module": {
			input: `
				module Foo
					module Bar; end
				end
				Bar
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(49, 5, 5), P(51, 5, 7)), "undefined constant `Bar`"),
			},
		},
		"define singleton class": {
			input: `
				module Foo
					singleton
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(21, 3, 6), P(38, 4, 8)), "singleton definitions cannot appear in this context"),
			},
		},
		"within method": {
			input: `
				def foo
					module Foo; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(32, 3, 20)), "module definitions cannot appear in this context"),
			},
		},
		"within singleton": {
			input: `
				class Foo
					singleton
						module Bar; end
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 4, 7), P(50, 4, 21)), "module definitions cannot appear in this context"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestStruct(t *testing.T) {
	tests := testTable{
		"struct with public constant": {
			input: `struct Foo; end`,
		},
		"instantiate a struct with all attributes": {
			input: `
				struct Foo
					a: String
					b: Int = 5
				end

				var f = Foo("a", 2)
			`,
		},
		"instantiate a struct without optional attributes": {
			input: `
				struct Foo
					a: String
					b: Int = 5
				end

				var f = Foo("a")
			`,
		},
		"instantiate a struct with invalid attributes": {
			input: `
				struct Foo
					a: String
					b: Int = 5
				end

				var f = Foo(5.2, 'b', :bar)
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(72, 7, 17), P(74, 7, 19)), "expected type `Std::String` for parameter `a` in call to `#init`, got type `5.2`"),
				error.NewFailure(L("<main>", P(77, 7, 22), P(79, 7, 24)), "expected type `Std::Int` for parameter `b` in call to `#init`, got type `\"b\"`"),
				error.NewFailure(L("<main>", P(68, 7, 13), P(86, 7, 31)), "expected 1...2 arguments in call to `#init`, got 3"),
			},
		},
		"call a getter on a struct": {
			input: `
				struct Foo
					a: String
					b: Int = 5
				end

				var f = Foo("a")
				var a: String = f.a
			`,
		},
		"call a getter on a struct and assign to a wrong type": {
			input: `
				struct Foo
					a: String
					b: Int = 5
				end

				var f = Foo("a")
				var a: String = f.b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(97, 8, 21), P(99, 8, 23)), "type `Std::Int` cannot be assigned to type `Std::String`"),
			},
		},
		"assign struct type to a singleton type": {
			input: `
				struct Foo
				end

				var a = Foo()
				var b: &Foo = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(61, 6, 19), P(61, 6, 19)), "type `Foo` cannot be assigned to type `&Foo`"),
			},
		},
		"assign struct to a singleton type": {
			input: `
				struct Foo; end

				var a: &Foo = Foo
			`,
		},
		"within method": {
			input: `
				def foo
					struct Foo; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(25, 3, 13), P(27, 3, 15)), "struct definitions cannot appear in this context"),
			},
		},
		"within singleton": {
			input: `
				class Foo
					singleton
						struct Bar; end
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 4, 14), P(45, 4, 16)), "struct definitions cannot appear in this context"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestClass(t *testing.T) {
	tests := testTable{
		"has its own local scope": {
			input: `
				a := 5
				class Foo
					a
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(31, 4, 6), P(31, 4, 6)), "undefined local `a`"),
			},
		},
		"class with public constant": {
			input: `class Foo; end`,
		},
		"class without a superclass": {
			input: `class Foo < nil; end`,
		},
		"call a builtin method on a class without a superclass": {
			input: `
				class Foo < nil; end
				Foo() == 2
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 3, 5), P(39, 3, 14)), "method `==` is not defined on type `Foo`"),
			},
		},
		"class with nonexistent superclass": {
			input: `class Foo < Bar; end`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 1, 13), P(14, 1, 15)), "undefined type `Bar`"),
				error.NewFailure(L("<main>", P(12, 1, 13), P(14, 1, 15)), "`void` is not a class"),
			},
		},
		"class with superclass": {
			input: `
				class Bar; end
				class Foo < Bar; end
			`,
		},
		"class with sealed superclass": {
			input: `
				sealed class Bar; end
				class Foo < Bar; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 3, 17), P(45, 3, 19)), "cannot inherit from sealed class `Bar`"),
			},
		},
		"primitive class with primitive superclass": {
			input: `
				primitive class Bar; end
				primitive class Foo < Bar; end
			`,
		},
		"class with primitive superclass": {
			input: `
				primitive class Bar; end
				class Foo < Bar; end
			`,
		},
		"class with module superclass": {
			input: `
				module Bar; end
				class Foo < Bar; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(37, 3, 17), P(39, 3, 19)), "`Bar` is not a class"),
			},
		},
		"include mixin with instance variables": {
			input: `
				mixin Foo
					var @foo: String?
				end
				class Bar
					include Foo

					def bar
						var f: String? = @foo
					end
				end
			`,
		},
		"include mixin with instance variables in a primitive": {
			input: `
				mixin Foo
					var @foo: String?
				end
				primitive class Bar
					include Foo

					def bar
						var f: String? = @foo
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(83, 6, 14), P(85, 6, 16)), "cannot include mixin with instance variables `Foo` in primitive `Bar`"),
				error.NewFailure(L("<main>", P(124, 9, 24), P(127, 9, 27)), "cannot use instance variables in this context"),
			},
		},
		"report errors for missing abstract methods from parent": {
			input: `
				abstract class Foo
					abstract def foo(); end
					def bar; end
				end
				class Bar < Foo; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(89, 6, 11), P(91, 6, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"report errors for missing abstract methods from parents": {
			input: `
				abstract class Foo
					abstract def foo(); end
					def fooo(); end
				end
				abstract class Bar < Foo
					abstract def bar(); end
					def barr; end
				end
				class Baz < Bar; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(177, 10, 11), P(179, 10, 13)), "missing abstract method implementation `Bar.:bar` with signature: `def bar(): void`"),
				error.NewFailure(L("<main>", P(177, 10, 11), P(179, 10, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"report errors for missing abstract methods from interfaces in parents": {
			input: `
				interface Foo
					def foo(); end
				end
				abstract class Bar
					implement Foo

					abstract def bar(); end
					def barr; end
				end
				class Baz < Bar; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(156, 11, 11), P(158, 11, 13)), "missing abstract method implementation `Bar.:bar` with signature: `def bar(): void`"),
				error.NewFailure(L("<main>", P(156, 11, 11), P(158, 11, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"report errors for missing abstract methods from mixin": {
			input: `
				abstract mixin Foo
					abstract def foo(); end
					def fooo(); end
				end
				class Bar
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(92, 6, 11), P(94, 6, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"report errors for missing abstract methods from mixins": {
			input: `
				abstract mixin Foo
					abstract def foo(); end
					def fooo(); end
				end
				abstract mixin Bar
					include Foo

					abstract def bar(); end
					def barr; end
				end
				class Baz
					include Bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(189, 12, 11), P(191, 12, 13)), "missing abstract method implementation `Bar.:bar` with signature: `def bar(): void`"),
				error.NewFailure(L("<main>", P(189, 12, 11), P(191, 12, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"report errors for missing abstract methods from interfaces in mixins": {
			input: `
				interface Foo
					def foo(); end
				end
				abstract mixin Bar
					implement Foo

					abstract def bar(); end
					def barr; end
				end
				class Baz
					include Bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(156, 11, 11), P(158, 11, 13)), "missing abstract method implementation `Bar.:bar` with signature: `def bar(): void`"),
				error.NewFailure(L("<main>", P(156, 11, 11), P(158, 11, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"report errors for missing abstract methods from interface": {
			input: `
				interface Foo
					def foo(); end
				end
				class Bar
					implement Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 5, 11), P(59, 5, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"report errors for missing abstract methods from interfaces": {
			input: `
				interface Foo
					def foo(); end
				end
				interface Bar
					implement Foo

					def bar(); end
				end
				class Baz
					implement Bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(123, 10, 11), P(125, 10, 13)), "missing abstract method implementation `Bar.:bar` with signature: `def bar(): void`"),
				error.NewFailure(L("<main>", P(123, 10, 11), P(125, 10, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"define and call a singleton method": {
			input: `
				class Foo
					singleton
						def foo; end
					end
					foo()
				end

				Foo.foo
			`,
		},
		"define and call a singleton method from parent": {
			input: `
				class Foo
					singleton
						def foo; end
					end
				end
				class Bar < Foo
				end

				Bar.foo
			`,
		},
		"assign class type to a class singleton type": {
			input: `
				class Foo
				end

				var a = Foo()
				var b: &Foo = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 6, 19), P(60, 6, 19)), "type `Foo` cannot be assigned to type `&Foo`"),
			},
		},
		"assign class to a class singleton type": {
			input: `
				class Foo; end

				var a: &Foo = Foo
			`,
		},
		"declare within method": {
			input: `
				def foo
					class Foo; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(31, 3, 19)), "class definitions cannot appear in this context"),
			},
		},
		"declare within singleton": {
			input: `
				class Foo
					singleton
						class Bar; end
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 4, 7), P(49, 4, 20)), "class definitions cannot appear in this context"),
			},
		},
		"declare a class inheriting from itself": {
			input: `
				class Foo < Foo; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(17, 2, 17), P(19, 2, 19)), "Type `Foo` circularly references itself"),
			},
		},
		"declare a class inheriting from its child": {
			input: `
				class Foo < Bar; end
				class Bar < Foo; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 3, 17), P(44, 3, 19)), "Type `Foo` circularly references itself"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestInstanceVariables(t *testing.T) {
	tests := testTable{
		"declare an instance variable in a primitive class": {
			input: `
				primitive class Foo
					var @foo: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 3, 6), P(45, 3, 21)), "cannot declare instance variable `foo` in a primitive `Foo`"),
			},
		},
		"declare an instance variable in a class": {
			input: `
				class Foo
					var @foo: String
				end
			`,
		},
		"redeclare an instance variable in a class": {
			input: `
				class Foo
					var @foo: String
					var @foo: Int
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 4, 6), P(54, 4, 18)), "cannot redeclare instance variable `@foo` with a different type, is `Std::Int`, should be `Std::String`, previous definition found in `Foo`"),
			},
		},
		"redeclare an instance variable in a class with a supertype": {
			input: `
				class Foo
					var @foo: String
					var @foo: String?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 4, 6), P(58, 4, 22)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String?`, should be `Std::String`, previous definition found in `Foo`"),
			},
		},
		"redeclare an instance variable in a class with a subtype": {
			input: `
				class Foo
					var @foo: String?
					var @foo: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 4, 6), P(58, 4, 21)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String`, should be `Std::String?`, previous definition found in `Foo`"),
			},
		},
		"redeclare an instance variable in a class with the same type": {
			input: `
				class Foo
					var @foo: String
					var @foo: String
				end
			`,
		},
		"declare an instance variable in a singleton class": {
			input: `
				class Foo
					singleton
						var @foo: String
					end
				end
			`,
		},
		"declare an instance variable in a mixin": {
			input: `
				mixin Foo
					var @foo: String
				end
			`,
		},
		"declare an instance variable in a module": {
			input: `
				module Foo
					var @foo: String
				end
			`,
		},
		"declare an instance variable in an interface": {
			input: `
				interface Foo
					var @foo: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 3, 6), P(39, 3, 21)), "cannot declare instance variable `@foo` in this context"),
			},
		},
		"use instance variable in a class": {
			input: `
				class Foo
					var @foo: String
					@foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 4, 6), P(45, 4, 9)), "undefined instance variable `@foo` in type `&Foo`"),
			},
		},
		"use instance variable in an instance method of a class": {
			input: `
				class Foo
					var @foo: String
					def bar then @foo
				end
			`,
		},
		"use singleton instance variable in a class": {
			input: `
				class Foo
				  singleton
						var @foo: String
					end
					@foo
				end
			`,
		},
		"use instance variable in a module": {
			input: `
				module Foo
					var @foo: String
					@foo
					def bar then @foo
				end
			`,
		},
		"use instance variable in a mixin": {
			input: `
				mixin Foo
					var @foo: String
					@foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 4, 6), P(45, 4, 9)), "undefined instance variable `@foo` in type `&Foo`"),
			},
		},
		"use instance variable in an instance method of a mixin": {
			input: `
				mixin Foo
					var @foo: String
					def bar then @foo
				end
			`,
		},
		"use singleton instance variable in a mixin": {
			input: `
				mixin Foo
				  singleton
						var @foo: String
					end
					@foo
				end
			`,
		},
		"use singleton instance variable in an interface": {
			input: `
				interface Foo
				  singleton
						var @foo: String
					end
					@foo
				end
			`,
		},

		"assign an instance variable with a matching type": {
			input: `
				module Foo
					var @foo: String
					@foo = "foo"
				end
			`,
		},
		"assign an instance variable with a non-matching type": {
			input: `
				module Foo
					var @foo: String
					@foo = 2
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(50, 4, 13), P(50, 4, 13)), "type `2` cannot be assigned to type `Std::String`"),
			},
		},
		"assign an inexistent instance variable": {
			input: `
				module Foo
					@foo = 2
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(21, 3, 6), P(24, 3, 9)), "undefined instance variable `@foo` in type `Foo`"),
			},
		},
		"declare within a method": {
			input: `
				def foo
					var @a: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(31, 3, 19)), "instance variable definitions cannot appear in this context"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestClassOverride(t *testing.T) {
	tests := testTable{
		"superclass matches": {
			input: `
				class Foo; end

				class Bar < Foo; end

				class Bar < Foo
					def bar; end
				end
			`,
		},
		"sealed modifier matches": {
			input: `
				class Foo; end

				sealed class Bar < Foo; end

				sealed class Bar < Foo
					def bar; end
				end
			`,
		},
		"abstract modifier matches": {
			input: `
				class Foo; end

				abstract class Bar < Foo; end

				abstract class Bar < Foo
					def bar; end
				end
			`,
		},
		"modifier was default, is abstract": {
			input: `
				class Foo; end

				class Bar < Foo; end

				abstract class Bar < Foo
					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 6, 5), P(100, 8, 7)), "cannot redeclare class `Bar` with a different modifier, is `abstract`, should be `default`"),
			},
		},
		"modifier was default, is primitive": {
			input: `
				class Foo; end

				class Bar < Foo; end

				primitive class Bar < Foo
					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 6, 5), P(101, 8, 7)), "cannot redeclare class `Bar` with a different modifier, is `primitive`, should be `default`"),
			},
		},
		"modifier was default, is sealed": {
			input: `
				class Foo; end

				class Bar < Foo; end

				sealed class Bar < Foo
					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 6, 5), P(98, 8, 7)), "cannot redeclare class `Bar` with a different modifier, is `sealed`, should be `default`"),
			},
		},
		"modifier was abstract, is sealed": {
			input: `
				class Foo; end

				abstract class Bar < Foo; end

				sealed class Bar < Foo
					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 6, 5), P(107, 8, 7)), "cannot redeclare class `Bar` with a different modifier, is `sealed`, should be `abstract`"),
			},
		},
		"modifier was abstract, is default": {
			input: `
				class Foo; end

				abstract class Bar < Foo; end

				class Bar < Foo
					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 6, 5), P(100, 8, 7)), "cannot redeclare class `Bar` with a different modifier, is `default`, should be `abstract`"),
			},
		},
		"modifier was primitive, is default": {
			input: `
				primitive class Foo; end

				primitive class Bar < Foo; end

				class Bar < Foo
					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(71, 6, 5), P(111, 8, 7)), "cannot redeclare class `Bar` with a different modifier, is `default`, should be `primitive`"),
			},
		},
		"superclass does not match": {
			input: `
				class Foo; end

				class Bar < Foo; end

				class Bar
					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 6, 5), P(85, 8, 7)), "superclass mismatch in `Bar`, got `Std::Object`, expected `Foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestInclude(t *testing.T) {
	tests := testTable{
		"include inexistent mixin": {
			input: `include Foo`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(10, 1, 11)), "cannot include mixins in this context"),
			},
		},
		"include in top level": {
			input: `
				mixin Foo; end
				include Foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 3, 5), P(34, 3, 15)), "cannot include mixins in this context"),
			},
		},
		"include in module": {
			input: `
				mixin Foo; end
			  module Bar
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(41, 4, 6), P(51, 4, 16)), "cannot include mixins in this context"),
			},
		},
		"include in interface": {
			input: `
				mixin Foo; end
			  interface Bar
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 4, 6), P(54, 4, 16)), "cannot include mixins in this context"),
			},
		},
		"include in class": {
			input: `
				mixin Foo; end
			  class  Bar
					include Foo
				end
			`,
		},
		"include in singleton class": {
			input: `
				mixin Foo; end
			  class  Bar
					singleton
						include Foo
					end
				end
			`,
		},
		"include in mixin": {
			input: `
				mixin Foo; end
				mixin Bar
					include Foo
				end
			`,
		},
		"include module": {
			input: `
				module Foo; end
				class  Bar
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(49, 4, 14), P(51, 4, 16)), "only mixins can be included"),
			},
		},
		"include class": {
			input: `
				class Foo; end
				class  Bar
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 4, 14), P(50, 4, 16)), "only mixins can be included"),
			},
		},
		"include interface": {
			input: `
				interface Foo; end
				class  Bar
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(52, 4, 14), P(54, 4, 16)), "only mixins can be included"),
			},
		},
		"include mixin with compatible methods": {
			input: `
				class Foo
					def foo(f: String?): String?; end
				end
				mixin Bar
					def foo(f: Value): String then "bar"
				end
				class Baz < Foo
					include Bar
				end
			`,
		},
		"include mixin with incompatible methods": {
			input: `
				class Foo
					def foo(f: Object): String then "bar"
				end
				mixin Bar
					def foo(f: String?): String?; end
				end
				class Baz < Foo
					include Bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(160, 9, 14), P(162, 9, 16)), "cannot include `Bar` in `Baz`:\n\n  - incompatible definitions of method `foo`\n      `Bar` has: `def foo(f: Std::String?): Std::String?`\n      `Foo` has: `def foo(f: Std::Object): Std::String`\n"),
			},
		},
		"include mixin with incompatible methods in parent": {
			input: `
				class Foo
					def foo(f: Object): String then "bar"
				end
				class Fooo < Foo; end
				mixin Bar
					def foo(f: String?): String?; end
				end
				mixin Barr
					include Bar
				end
				class Baz < Fooo
					include Barr
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(227, 13, 14), P(230, 13, 17)), "cannot include `Barr` in `Baz`:\n\n  - incompatible definitions of method `foo`\n      `Bar` has: `def foo(f: Std::String?): Std::String?`\n      `Foo` has: `def foo(f: Std::Object): Std::String`\n"),
			},
		},
		"include mixin with incompatible instance variables": {
			input: `
				class Foo
					var @foo: Object?
				end
				mixin Bar
					var @foo: String?
				end
				class Baz < Foo
					include Bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(124, 9, 14), P(126, 9, 16)), "cannot include `Bar` in `Baz`:\n\n  - incompatible definitions of instance variable `@foo`\n      `Bar` has: `var @foo: Std::String?`\n      `Foo` has: `var @foo: Std::Object?`\n"),
			},
		},
		"include mixin with incompatible instance variables in parent": {
			input: `
				class Foo
					var @foo: Object?
				end
				class Fooo < Foo; end
				mixin Bar
					var @foo: String?
				end
				mixin Barr
					include Bar
				end
				class Baz < Fooo
					include Barr
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(191, 13, 14), P(194, 13, 17)), "cannot include `Barr` in `Baz`:\n\n  - incompatible definitions of instance variable `@foo`\n      `Bar` has: `var @foo: Std::String?`\n      `Foo` has: `var @foo: Std::Object?`\n"),
			},
		},
		"include within a method": {
			input: `
				def foo
					include Bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(28, 3, 16)), "cannot include mixins in this context"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestImplement(t *testing.T) {
	tests := testTable{
		"implement inexistent interface": {
			input: `
				class Foo
					implement Bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 3, 16), P(32, 3, 18)), "undefined type `Bar`"),
				error.NewFailure(L("<main>", P(30, 3, 16), P(32, 3, 18)), "only interfaces can be implemented"),
			},
		},
		"implement in top level": {
			input: `
				interface Foo; end
				implement Foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 3, 05), P(40, 3, 17)), "cannot implement interfaces in this context"),
			},
		},
		"implement in module": {
			input: `
				interface Foo; end
			  module Bar
					implement Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 4, 6), P(57, 4, 18)), "cannot implement interfaces in this context"),
			},
		},
		"implement in interface": {
			input: `
				interface Foo; end
			  interface Bar
					implement Foo
				end
			`,
		},
		"implement in class": {
			input: `
				interface Foo; end
			  class  Bar
					implement Foo
				end
			`,
		},
		"implement in mixin": {
			input: `
				interface Foo; end
				mixin Bar
					implement Foo
				end
			`,
		},
		"implement module": {
			input: `
				module Foo; end
				class  Bar
					implement Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 4, 16), P(53, 4, 18)), "only interfaces can be implemented"),
			},
		},
		"implement class": {
			input: `
				class Foo; end
				class  Bar
					implement Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(50, 4, 16), P(52, 4, 18)), "only interfaces can be implemented"),
			},
		},
		"implement mixin": {
			input: `
				mixin Foo; end
				class  Bar
					implement Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(50, 4, 16), P(52, 4, 18)), "only interfaces can be implemented"),
			},
		},
		"implement within a method": {
			input: `
				def foo
					implement Bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(30, 3, 18)), "cannot implement interfaces in this context"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestMixinType(t *testing.T) {
	tests := testTable{
		"assign instance of related class to mixin": {
			input: `
				mixin Bar; end
				class Foo
					include Bar
				end

				var a: Bar = Foo()
			`,
		},
		"assign instance of unrelated class to mixin": {
			input: `
				mixin Bar; end
				class Foo; end

				var a: Bar = Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 5, 18), P(61, 5, 22)), "type `Foo` cannot be assigned to type `Bar`"),
			},
		},
		"assign mixin type to the same mixin type": {
			input: `
				mixin Bar; end
				class Foo
					include Bar
				end

				var a: Bar = Foo()
				var b: Bar = a
			`,
		},
		"assign related mixin type to a mixin type": {
			input: `
				mixin Baz; end

				mixin Bar
					include Baz
				end

				class Foo
					include Bar
				end

				var a: Bar = Foo()
				var b: Baz = a
			`,
		},
		"assign mixin type to a mixin singleton type": {
			input: `
				mixin Foo; end

				class Bar
					include Foo
				end

				var a: Foo = Bar()
				var b: &Foo = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(102, 9, 19), P(102, 9, 19)), "type `Foo` cannot be assigned to type `&Foo`"),
			},
		},
		"assign mixin to a mixin singleton type": {
			input: `
				mixin Foo; end

				var a: &Foo = Foo
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestMixinOverride(t *testing.T) {
	tests := testTable{
		"default modifier matches": {
			input: `
				mixin Bar; end
				mixin Bar; end
			`,
		},
		"abstract modifier matches": {
			input: `
				abstract mixin Bar; end
				abstract mixin Bar; end
			`,
		},
		"modifier was default, is abstract": {
			input: `
				mixin Bar; end
				abstract mixin Bar; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 3, 5), P(46, 3, 27)), "cannot redeclare mixin `Bar` with a different modifier, is `abstract`, should be `default`"),
			},
		},
		"modifier was abstract, is default": {
			input: `
				abstract mixin Bar; end
				mixin Bar; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 3, 5), P(46, 3, 18)), "cannot redeclare mixin `Bar` with a different modifier, is `default`, should be `abstract`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestMixin(t *testing.T) {
	tests := testTable{
		"has its own local scope": {
			input: `
				a := 5
				mixin Foo
					a
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(31, 4, 6), P(31, 4, 6)), "undefined local `a`"),
			},
		},
		"report errors for missing abstract methods from mixin parent": {
			input: `
				abstract mixin Foo
					abstract def foo(); end
					def bar; end
				end
				mixin Bar
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(89, 6, 11), P(91, 6, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"report errors for missing abstract methods from interface": {
			input: `
				interface Foo
					def foo(); end
				end
				mixin Bar
					implement Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 5, 11), P(59, 5, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"report errors for missing abstract methods from interfaces": {
			input: `
				interface Foo
					def foo(); end
				end
				interface Bar
					implement Foo

					def bar(); end
				end
				mixin Baz
					implement Bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(123, 10, 11), P(125, 10, 13)), "missing abstract method implementation `Bar.:bar` with signature: `def bar(): void`"),
				error.NewFailure(L("<main>", P(123, 10, 11), P(125, 10, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"report errors for missing abstract methods from interfaces in mixins": {
			input: `
				interface Foo
					def foo(); end
				end
				abstract mixin Bar
					implement Foo

					abstract def bar(); end
				end
				mixin Baz
					include Bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(137, 10, 11), P(139, 10, 13)), "missing abstract method implementation `Bar.:bar` with signature: `def bar(): void`"),
				error.NewFailure(L("<main>", P(137, 10, 11), P(139, 10, 13)), "missing abstract method implementation `Foo.:foo` with signature: `def foo(): void`"),
			},
		},
		"define and call a singleton method": {
			input: `
				mixin Foo
					singleton
						def foo; end
					end
					foo()
				end

				Foo.foo
			`,
		},
		"define and call a singleton method from parent": {
			input: `
				mixin Foo
					singleton
						def foo; end
					end
				end
				mixin Bar
				  include Foo
				end

				Bar.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(111, 11, 5), P(117, 11, 11)), "method `foo` is not defined on type `&Bar`"),
			},
		},
		"within a method": {
			input: `
				def foo
					mixin Bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(31, 3, 19)), "mixin definitions cannot appear in this context"),
			},
		},
		"include itself": {
			input: `
				mixin Foo
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 3, 14), P(30, 3, 16)), "Type `Foo` circularly references itself"),
			},
		},
		"circular include": {
			input: `
				mixin Foo
					include Bar
				end

				mixin Bar
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(68, 7, 14), P(70, 7, 16)), "Type `Foo` circularly references itself"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestInterface(t *testing.T) {
	tests := testTable{
		"has its own local scope": {
			input: `
				a := 5
				interface Foo
					a
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 4, 6), P(35, 4, 6)), "undefined local `a`"),
			},
		},
		"define and call a singleton method": {
			input: `
				interface Foo
					singleton
						def foo; end
					end
					foo()
				end

				Foo.foo
			`,
		},
		"define and call a singleton method from parent": {
			input: `
				interface Foo
					singleton
						def foo; end
					end
				end
				interface Bar
				  implement Foo
				end

				Bar.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(121, 11, 5), P(127, 11, 11)), "method `foo` is not defined on type `&Bar`"),
			},
		},
		"within a method": {
			input: `
				def foo
					interface Bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(18, 3, 6), P(35, 3, 23)), "interface definitions cannot appear in this context"),
			},
		},
		"implement itself": {
			input: `
				interface Foo
					implement Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 16), P(36, 3, 18)), "Type `Foo` circularly references itself"),
			},
		},
		"circular implement": {
			input: `
				interface Foo
					implement Bar
				end

				interface Bar
					implement Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(80, 7, 16), P(82, 7, 18)), "Type `Foo` circularly references itself"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestInterfaceType(t *testing.T) {
	tests := testTable{
		"assign instance of class that implements the interface explicitly": {
			input: `
				interface Foo
					def foo; end
				end
				class Bar
					implement Foo

					def foo; end
				end

				var a: Foo = Bar()
			`,
		},
		"assign instance of class that implements the interface implicitly": {
			input: `
				interface Foo
					def foo; end
				end
				class Bar
					def foo; end
				end

				var a: Foo = Bar()
			`,
		},
		"assign instance of class that does not implement the interface": {
			input: `
				interface Foo
					def foo; end
				end
				class Bar; end

				var a: Foo = Bar()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(82, 7, 18), P(86, 7, 22)), "type `Bar` does not implement interface `Foo`:\n\n  - missing method `Foo.:foo` with signature: `def foo(): void`\n"),
				error.NewFailure(L("<main>", P(82, 7, 18), P(86, 7, 22)), "type `Bar` cannot be assigned to type `Foo`"),
			},
		},
		"assign interface type to the same interface type": {
			input: `
				interface Foo
					def foo; end
				end
				class Bar
					def foo; end
				end

				var a: Foo = Bar()
				var b: Foo = a
			`,
		},
		"assign interface that implements another interface explicitly": {
			input: `
				interface Foo
					def foo; end
				end
				interface Bar
					implement Foo

					def bar; end
				end
				class Baz
				  def foo; end
					def bar; end
				end

				var a: Bar = Baz()
				var b: Foo = a
			`,
		},
		"assign unrelated interface type to an interface type": {
			input: `
				interface Foo
					def foo; end
				end
				interface Bar
				  def foo; end
					def bar; end
				end
				class Baz
				  def foo; end
					def bar; end
				end

				var a: Bar = Baz()
				var b: Foo = a
			`,
		},
		"assign interface that implements another interface implicitly": {
			input: `
				interface Foo
					def foo; end
				end
				interface Bar
				  def foo; end
					def bar; end
				end
				class Baz
				  def foo; end
					def bar; end
				end

				var a: Bar = Baz()
				var b: Foo = a
			`,
		},
		"assign interface that does not implement another interface": {
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

				var a: Bar = Baz()
				var b: Foo = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(189, 14, 18), P(189, 14, 18)), "type `Bar` does not implement interface `Foo`:\n\n  - missing method `Foo.:foo` with signature: `def foo(): void`\n"),
				error.NewFailure(L("<main>", P(189, 14, 18), P(189, 14, 18)), "type `Bar` cannot be assigned to type `Foo`"),
			},
		},
		"assign interface type to an interface singleton type": {
			input: `
				interface Foo; end
				class Bar; end

				var a: Foo = Bar()
				var b: &Foo = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 6, 19), P(85, 6, 19)), "type `Foo` cannot be assigned to type `&Foo`"),
			},
		},
		"assign interface to an interface singleton type": {
			input: `
				interface Foo; end

				var a: &Foo = Foo
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
