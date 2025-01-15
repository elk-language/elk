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
		"module with constant in the middle defined later": {
			input: `
				class Foo::Bar::Baz; end
				class Foo
					class Bar; end
				end
			`,
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
				error.NewWarning(L("<main>", P(30, 3, 5), P(34, 3, 9)), "this equality check is impossible, `Foo` cannot ever be equal to `2`"),
			},
		},
		"class with nonexistent superclass": {
			input: `class Foo < Bar; end`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 1, 13), P(14, 1, 15)), "undefined type `Bar`"),
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
				error.NewFailure(L("<main>", P(89, 6, 11), P(91, 6, 13)), "missing abstract method implementations in `Bar`:\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
			},
		},
		"report errors for missing abstract methods from generic parent": {
			input: `
				abstract class Foo[V]
					abstract def foo(): V; end
					def bar; end
				end
				class Bar < Foo[String]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(95, 6, 11), P(97, 6, 13)), "missing abstract method implementations in `Bar`:\n\n  - method `Foo.:foo`:\n      `def foo(): Std::String`\n"),
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
				error.NewFailure(L("<main>", P(177, 10, 11), P(179, 10, 13)), "missing abstract method implementations in `Baz`:\n\n  - method `Bar.:bar`:\n      `def bar(): void`\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
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
				error.NewFailure(L("<main>", P(156, 11, 11), P(158, 11, 13)), "missing abstract method implementations in `Baz`:\n\n  - method `Bar.:bar`:\n      `def bar(): void`\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
			},
		},
		"report errors for missing abstract methods from generic interfaces in parents": {
			input: `
				interface Foo[V]
					sig foo(): V
				end
				abstract class Bar
					implement Foo[Int]

					abstract def bar(); end
					def barr; end
				end
				class Baz < Bar; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(162, 11, 11), P(164, 11, 13)), "missing abstract method implementations in `Baz`:\n\n  - method `Bar.:bar`:\n      `def bar(): void`\n\n  - method `Foo.:foo`:\n      `def foo(): Std::Int`\n"),
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
				error.NewFailure(L("<main>", P(92, 6, 11), P(94, 6, 13)), "missing abstract method implementations in `Bar`:\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
			},
		},
		"report errors for missing abstract methods from generic mixin": {
			input: `
				abstract mixin Foo[V]
					abstract def foo(): V; end
					def fooo(); end
				end
				class Bar
					include Foo[Float]
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(98, 6, 11), P(100, 6, 13)), "missing abstract method implementations in `Bar`:\n\n  - method `Foo.:foo`:\n      `def foo(): Std::Float`\n"),
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
				error.NewFailure(L("<main>", P(189, 12, 11), P(191, 12, 13)), "missing abstract method implementations in `Baz`:\n\n  - method `Bar.:bar`:\n      `def bar(): void`\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
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
				error.NewFailure(L("<main>", P(156, 11, 11), P(158, 11, 13)), "missing abstract method implementations in `Baz`:\n\n  - method `Bar.:bar`:\n      `def bar(): void`\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
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
				error.NewFailure(L("<main>", P(57, 5, 11), P(59, 5, 13)), "missing abstract method implementations in `Bar`:\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
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
				error.NewFailure(L("<main>", P(123, 10, 11), P(125, 10, 13)), "missing abstract method implementations in `Baz`:\n\n  - method `Bar.:bar`:\n      `def bar(): void`\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
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
				error.NewFailure(L("<main>", P(17, 2, 17), P(19, 2, 19)), "type `Foo` circularly references itself"),
			},
		},
		"declare a class inheriting from its child": {
			input: `
				class Foo < Bar; end
				class Bar < Foo; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 3, 17), P(44, 3, 19)), "type `Foo` circularly references itself"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestGenericClass(t *testing.T) {
	tests := testTable{
		"redeclare non generic class as generic": {
			input: `
				class Foo; end
				class Foo[V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 3, 5), P(40, 3, 21)), "type parameter count mismatch in `Foo`, got: 1, expected: 0"),
			},
		},
		"redeclare generic class as non generic": {
			input: `
				class Foo[V]; end
				class Foo; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 5), P(40, 3, 18)), "type parameter count mismatch in `Foo`, got: 0, expected: 1"),
			},
		},
		"redeclare generic class with missing type parameter": {
			input: `
				class Foo[V, T]; end
				class Foo[V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 3, 5), P(46, 3, 21)), "type parameter count mismatch in `Foo`, got: 1, expected: 2"),
			},
		},
		"redeclare generic class with additional type parameter": {
			input: `
				class Foo[V]; end
				class Foo[V, T]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 5), P(46, 3, 24)), "type parameter count mismatch in `Foo`, got: 2, expected: 1"),
			},
		},
		"redeclare generic class with matching type parameters": {
			input: `
				class Foo[+V < String]; end
				class Foo[+V < String]; end
			`,
		},
		"redeclare generic class with wrong type param name": {
			input: `
				class Foo[V]; end
				class Foo[T]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 5), P(43, 3, 21)), "type parameter mismatch in `Foo`, is `T`, should be `V`"),
			},
		},
		"redeclare generic class with wrong type param variance": {
			input: `
				class Foo[-V]; end
				class Foo[+V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 3, 5), P(45, 3, 22)), "type parameter mismatch in `Foo`, is `+V`, should be `-V`"),
			},
		},
		"redeclare generic class with wrong type param upper bound": {
			input: `
				class Foo[V < String]; end
				class Foo[V < Int]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 5), P(58, 3, 27)), "type parameter mismatch in `Foo`, is `V < Std::Int`, should be `V < Std::String`"),
			},
		},
		"redeclare generic class with wrong type param lower bound": {
			input: `
				class Foo[V > String]; end
				class Foo[V > Int]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 5), P(58, 3, 27)), "type parameter mismatch in `Foo`, is `V > Std::Int`, should be `V > Std::String`"),
			},
		},

		"assign to related generic mixin with correct type args": {
			input: `
				mixin Bar[K, V]; end

				class Foo[V]
					include Bar[String, V]
				end

				var a: Bar[String, Int] = Foo::[Int]()
			`,
		},
		"assign to related generic mixin with incorrect type args": {
			input: `
				mixin Bar[K, V]; end

				class Foo[V]
					include Bar[String, V]
				end

				var a: Bar[String, Float] = Foo::[Int]()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(113, 8, 33), P(124, 8, 44)), "type `Foo[Std::Int]` cannot be assigned to type `Bar[Std::String, Std::Float]`"),
			},
		},
		"assign to distantly related generic mixin with correct type args": {
			input: `
				mixin Qux[E]; end
				mixin Baz[F]
					include Qux[F]
				end
				mixin Bar[K, V]
					include Baz[V]
				end

				class Foo[V]
					include Bar[String, V]
				end

				var a: Qux[Int] = Foo::[Int]()
			`,
		},

		"assign to related generic class with correct type args": {
			input: `
				class Bar[K, V]; end
				class Foo[V] < Bar[String, V]; end

				var a: Bar[String, Int] = Foo::[Int]()
			`,
		},
		"assign to related generic class with incorrect type args": {
			input: `
				class Bar[K, V]; end
				class Foo[V] < Bar[String, V]; end

				var a: Bar[String, Float] = Foo::[Int]()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(98, 5, 33), P(109, 5, 44)), "type `Foo[Std::Int]` cannot be assigned to type `Bar[Std::String, Std::Float]`"),
			},
		},
		"assign to distantly related generic class with correct type args": {
			input: `
				class Qux[E]; end
				class Baz[F] < Qux[F]; end
				class Bar[K, V] < Baz[V]; end
				class Foo[V] < Bar[String, V]; end

				var a: Qux[Int] = Foo::[Int]()
			`,
		},

		"inherit from generic class specifying type arguments": {
			input: `
				class Foo[V]; end
				class Bar < Foo[String]; end
			`,
		},
		"inherit from generic class forwarding type arguments": {
			input: `
				class Foo[V]; end
				class Bar[V] < Foo[V]; end
			`,
		},
		"inherit from generic class forwarding covariant type argument as covariant": {
			input: `
				class Foo[+V]; end
				class Bar[+V] < Foo[V]; end
			`,
		},
		"inherit from generic class forwarding contravariant type argument as contravariant": {
			input: `
				class Foo[-V]; end
				class Bar[-V] < Foo[V]; end
			`,
		},
		"inherit from generic class forwarding invariant type argument as covariant": {
			input: `
				class Foo[+V]; end
				class Bar[V] < Foo[V]; end
			`,
		},
		"inherit from generic class forwarding invariant type argument as contravariant": {
			input: `
				class Foo[-V]; end
				class Bar[V] < Foo[V]; end
			`,
		},
		"inherit from generic class forwarding covariant type argument as invariant": {
			input: `
				class Foo[V]; end
				class Bar[+V] < Foo[V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(47, 3, 25), P(47, 3, 25)), "covariant type `V` cannot appear in invariant position"),
			},
		},
		"inherit from generic class forwarding contravariant type argument as invariant": {
			input: `
				class Foo[V]; end
				class Bar[-V] < Foo[V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(47, 3, 25), P(47, 3, 25)), "contravariant type `V` cannot appear in invariant position"),
			},
		},
		"inherit from generic class forwarding covariant type argument as contravariant": {
			input: `
				class Foo[-V]; end
				class Bar[+V] < Foo[V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 3, 25), P(48, 3, 25)), "covariant type `V` cannot appear in contravariant position"),
			},
		},
		"inherit from generic class forwarding contravariant type argument as covariant": {
			input: `
				class Foo[+V]; end
				class Bar[-V] < Foo[V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 3, 25), P(48, 3, 25)), "contravariant type `V` cannot appear in covariant position"),
			},
		},
		"inherit from generic class without type arguments": {
			input: `
				class Foo[V]; end
				class Bar < Foo; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 3, 17), P(41, 3, 19)), "`Foo` requires 1 type argument(s), got: 0"),
			},
		},
		"inherit from generic class specifying too many type arguments": {
			input: `
				class Foo[V]; end
				class Bar < Foo[String, Float]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 3, 17), P(41, 3, 19)), "`Foo` requires 1 type argument(s), got: 2"),
			},
		},
		"call a method on a class that inherits from a generic class with specified type arguments": {
			input: `
				class Foo[V]
					def foo(a: V): V then a
				end
				class Bar < Foo[String]; end

				var a: String = Bar().foo("elo")
			`,
		},

		"return self type from instance method": {
			input: `
				class Foo
					def foo: self then self
				end

				var a: 9 = Foo().foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(68, 6, 16), P(76, 6, 24)), "type `Foo` cannot be assigned to type `9`"),
			},
		},
		"return self type from instance method of child": {
			input: `
				class Foo
					def foo: self then self
				end
				class Bar < Foo; end

				var a: 9 = Bar().foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(93, 7, 16), P(101, 7, 24)), "type `Bar` cannot be assigned to type `9`"),
			},
		},
		"return self type from singleton method": {
			input: `
				class Foo
					singleton
						def foo: self then self
					end
				end

				var a: 9 = Foo.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(93, 8, 16), P(99, 8, 22)), "type `&Foo` cannot be assigned to type `9`"),
			},
		},
		"return self type from singleton method of child": {
			input: `
				class Foo
					singleton
						def foo: self then self
					end
				end
				class Bar < Foo; end

				var a: 9 = Bar.foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(118, 9, 16), P(124, 9, 22)), "type `&Bar` cannot be assigned to type `9`"),
			},
		},

		"declare a generic class": {
			input: `
				class Foo[V]; end
			`,
		},
		"use a type parameter as a value": {
			input: `
				class Foo[V]
					V
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(23, 3, 6), P(23, 3, 6)), "`Foo::V` cannot be used as a value in expressions"),
			},
		},
		"use a type parameter from the outside": {
			input: `
				class Foo[V]; end
				var a: Foo::V
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 3, 17), P(39, 3, 17)), "type parameter `V` cannot be used in this context"),
			},
		},
		"use a type parameter from the outside within a method": {
			input: `
				class Foo[V]; end
				def foo
					var a: Foo::V
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(52, 4, 18), P(52, 4, 18)), "undefined type `V`"),
			},
		},
		"declare a generic class with an upper bound referencing itself": {
			input: `
				interface Lol[T]
					sig lol: T
				end
				class Foo[V < Lol[V]]; end

				class Bar
					def lol: Bar then self
				end
				Foo::[Bar]()
			`,
		},
		"declare a generic class with bounds": {
			input: `
				class Foo[V > Baz < Object]; end
				class Bar; end
				class Baz < Bar; end
			`,
		},
		"declare a generic class with invalid bounds": {
			input: `
				class Foo[V > Baz < Bar]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(19, 2, 19), P(21, 2, 21)), "undefined type `Baz`"),
				error.NewFailure(L("<main>", P(25, 2, 25), P(27, 2, 27)), "undefined type `Bar`"),
			},
		},
		"declare a generic class with invalid default": {
			input: `
				class Foo[V < String = Int]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 2, 28), P(30, 2, 30)), "type parameter `V` has an invalid default `Std::Int`, should be a subtype of `Std::String` and supertype of `never`"),
			},
		},
		"declare a generic class with required type args after optionals": {
			input: `
				class Foo[V = Int, Y]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 2, 24), P(24, 2, 24)), "required type parameter `Y` cannot appear after optional type parameters"),
			},
		},
		"declare a generic class with type params as defaults": {
			input: `
				class Foo[V, Y = V]; end
			`,
		},
		"use a generic class with type params as defaults": {
			input: `
				class Foo[V, Y = V]; end
				var a: Foo[String] = 5
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(55, 3, 26), P(55, 3, 26)), "type `5` cannot be assigned to type `Foo[Std::String, Std::String]`"),
			},
		},
		"use a class without optional type parameters": {
			input: `
				class Foo[V, Y = Int]; end
				var a: Foo[String] = nil
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 3, 26), P(59, 3, 28)), "type `nil` cannot be assigned to type `Foo[Std::String, Std::Int]`"),
			},
		},
		"use a class overriding optional type parameters": {
			input: `
				class Foo[V, Y = Int]; end
				var a: Foo[String, Char] = nil
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(63, 3, 32), P(65, 3, 34)), "type `nil` cannot be assigned to type `Foo[Std::String, Std::Char]`"),
			},
		},
		"use a class without required type parameters": {
			input: `
				class Foo[V, Y = Int]; end
				var a: Foo = nil
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 3, 12), P(45, 3, 14)), "`Foo` requires 1...2 type argument(s), got: 0"),
			},
		},
		"use a class with too many": {
			input: `
				class Foo[V, Y = Int]; end
				var a: Foo[Int, Char, String] = nil
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 3, 12), P(45, 3, 14)), "`Foo` requires 1...2 type argument(s), got: 3"),
			},
		},
		"assign related generic class to its parent with the same type argument": {
			input: `
				class Foo[V]; end
				class Bar[V] < Foo[V]; end

				var a = Foo::[Int]()
				a = Bar::[Int]()
			`,
		},
		"assign related generic class to its child with the same type argument": {
			input: `
				class Foo[V]; end
				class Bar[V] < Foo[V]; end

				var a = Bar::[Int]()
				a = Foo::[Int]()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(88, 6, 9), P(99, 6, 20)), "type `Foo[Std::Int]` cannot be assigned to type `Bar[Std::Int]`"),
			},
		},

		"invariant type param - assign to self": {
			input: `
				class Foo[V]; end

				var a = Foo::[Int]()
				a = Foo::[Int]()
			`,
		},
		"invariant type param - assign to parent": {
			input: `
				class Foo[V]; end

				var a = Foo::[Value]()
				a = Foo::[Int]()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(59, 5, 9), P(70, 5, 20)), "type `Foo[Std::Int]` cannot be assigned to type `Foo[Std::Value]`"),
			},
		},
		"invariant type param - assign to child": {
			input: `
				class Foo[V]; end

				var a = Foo::[Int]()
				a = Foo::[Value]()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 5, 9), P(70, 5, 22)), "type `Foo[Std::Value]` cannot be assigned to type `Foo[Std::Int]`"),
			},
		},

		"covariant type param - assign to self": {
			input: `
				class Foo[+V]; end

				var a = Foo::[Int]()
				a = Foo::[Int]()
			`,
		},
		"covariant type param - assign to parent": {
			input: `
				class Foo[+V]; end

				var a = Foo::[Value]()
				a = Foo::[Int]()
			`,
		},
		"covariant type param - assign to child": {
			input: `
				class Foo[+V]; end

				var a = Foo::[Int]()
				a = Foo::[Value]()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(58, 5, 9), P(71, 5, 22)), "type `Foo[Std::Value]` cannot be assigned to type `Foo[Std::Int]`"),
			},
		},

		"contravariant type param - assign to self": {
			input: `
				class Foo[-V]; end

				var a = Foo::[Int]()
				a = Foo::[Int]()
			`,
		},
		"contravariant type param - assign to parent": {
			input: `
				class Foo[-V]; end

				var a = Foo::[Value]()
				a = Foo::[Int]()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 5, 9), P(71, 5, 20)), "type `Foo[Std::Int]` cannot be assigned to type `Foo[Std::Value]`"),
			},
		},
		"contravariant type param - assign to child": {
			input: `
				class Foo[-V]; end

				var a = Foo::[Int]()
				a = Foo::[Value]()
			`,
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
		"declare an instance variable in an extend where block": {
			input: `
				class Foo[T]
					extend where T < String
						var @foo: String
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(53, 4, 7), P(68, 4, 22)), "cannot declare instance variable `@foo` in this context"),
			},
		},
		"declare an instance variable in a primitive class": {
			input: `
				primitive class Foo
					var @foo: String?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 3, 6), P(46, 3, 22)), "cannot declare instance variable `foo` in a primitive `Foo`"),
			},
		},
		"include non-nilable instance variables in a class": {
			input: `
				mixin Foo
					var @foo: String
				end

				class Bar
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(50, 6, 5), P(83, 8, 7)), "instance variable `var @foo: Std::String` must be initialised in the constructor, since it is not nilable"),
			},
		},
		"include non-nilable instance variables in a class and initialise them in init": {
			input: `
				mixin Foo
					var @foo: String
				end

				class Bar
					include Foo

					init(@foo); end
				end
			`,
		},
		"declare a non-nilable instance variable in a class": {
			input: `
				class Foo
					var @foo: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(5, 2, 5), P(43, 4, 7)), "instance variable `var @foo: Std::String` must be initialised in the constructor, since it is not nilable"),
			},
		},
		"declare a non-nilable instance variable in a class and assign it in init": {
			input: `
				class Foo
					var @foo: String

					init
						@foo = "lol"
					end
				end
			`,
		},
		"declare a non-nilable instance variable in a class and assign it in init param": {
			input: `
				class Foo
					var @foo: String

					init(@foo); end
				end
			`,
		},
		"declare a non-nilable instance variable in a class and assign it in init after method calls": {
			input: `
				class Foo
					var @foo: String

					init
						bar()
						@foo = "lol"
					end

					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(54, 6, 7), P(58, 6, 11)), "instance variable `var @foo: Std::String` must be initialised before `self` can be used, since it is non-nilable"),
			},
		},
		"declare a non-nilable instance variable in a class and assign it in init after reading self": {
			input: `
				class Foo
					var @foo: String

					init
						self
						@foo = "lol"
					end

					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(54, 6, 7), P(57, 6, 10)), "instance variable `var @foo: Std::String` must be initialised before `self` can be used, since it is non-nilable"),
			},
		},
		"declare a non-nilable instance variable in a class and assign it in init after reading instance variables": {
			input: `
				class Foo
					var @foo: String

					init
						@foo
						@foo = "lol"
					end

					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(54, 6, 7), P(57, 6, 10)), "instance variable `var @foo: Std::String` must be initialised before `self` can be used, since it is non-nilable"),
			},
		},
		"declare an instance variable in a class": {
			input: `
				class Foo
					var @foo: String?
				end
			`,
		},
		"redeclare an instance variable in a class": {
			input: `
				class Foo
					var @foo: String?
					var @foo: Int?
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 4, 6), P(56, 4, 19)), "cannot redeclare instance variable `@foo` with a different type, is `Std::Int?`, should be `Std::String?`, previous definition found in `Foo`"),
			},
		},
		"redeclare an instance variable in a class with a supertype": {
			input: `
				class Foo
					var @foo: String?
					var @foo: String | Float | nil
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 4, 6), P(72, 4, 35)), "cannot redeclare instance variable `@foo` with a different type, is `Std::String | Std::Float | nil`, should be `Std::String?`, previous definition found in `Foo`"),
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
					var @foo: String?
					var @foo: String?
				end
			`,
		},
		"declare an instance variable in a singleton class": {
			input: `
				class Foo
					singleton
						var @foo: String?
					end
				end
			`,
		},
		"declare a non-nilable instance variable in a mixin": {
			input: `
				mixin Foo
					var @foo: String
				end
			`,
		},
		"declare an instance variable in a mixin": {
			input: `
				mixin Foo
					var @foo: String?
				end
			`,
		},
		"declare a non-nilable instance variable in a module": {
			input: `
				module Foo
					var @foo: String
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(21, 3, 6), P(36, 3, 21)), "instance variable `@foo` must be declared as nilable"),
			},
		},
		"declare an instance variable in a module": {
			input: `
				module Foo
					var @foo: String?
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
					var @foo: String?
					@foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(43, 4, 6), P(46, 4, 9)), "undefined instance variable `@foo` in type `&Foo`"),
			},
		},
		"use instance variable in an instance method of a class": {
			input: `
				class Foo
					var @foo: String?
					def bar then @foo
				end
			`,
		},
		"declare non-nilable instance variable in a singleton class": {
			input: `
				class Foo
				  singleton
						var @foo: String
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(37, 4, 7), P(52, 4, 22)), "instance variable `@foo` must be declared as nilable"),
			},
		},
		"use singleton instance variable in a class": {
			input: `
				class Foo
				  singleton
						var @foo: String?
					end
					@foo
				end
			`,
		},
		"use instance variable in a module": {
			input: `
				module Foo
					var @foo: String?
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
					var @foo: String?
					def bar then @foo
				end
			`,
		},
		"use singleton instance variable in a mixin": {
			input: `
				mixin Foo
				  singleton
						var @foo: String?
					end
					@foo
				end
			`,
		},
		"use singleton instance variable in an interface": {
			input: `
				interface Foo
				  singleton
						var @foo: String?
					end
					@foo
				end
			`,
		},

		"assign an instance variable with a matching type": {
			input: `
				module Foo
					var @foo: String?
					@foo = "foo"
				end
			`,
		},
		"assign an instance variable with a non-matching type": {
			input: `
				module Foo
					var @foo: String?
					@foo = 2
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 4, 13), P(51, 4, 13)), "type `2` cannot be assigned to type `Std::String?`"),
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
		"resolve an instance variable inherited from a generic parent": {
			input: `
				class Foo[V]
					var @foo: V?
				end
				class Bar < Foo[String]
					def bar: String? then @foo
				end
			`,
		},
		"resolve an instance variable inherited from a distant generic parent": {
			input: `
				class Qux[E]
					var @foo: E?
				end
				class Foo[V] < Qux[V]; end
				class Bar < Foo[String]
					def bar: String? then @foo
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

func TestClassOverride(t *testing.T) {
	tests := testTable{
		"superclass was nil, is Foo": {
			input: `
				class Foo; end

				class Bar < nil; end

				class Bar < Foo
					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(63, 6, 17), P(65, 6, 19)), "superclass mismatch in `Bar`, got `Foo`, expected `nil`"),
			},
		},
		"superclass matches": {
			input: `
				class Foo; end

				class Bar < Foo; end

				class Bar < Foo
					def bar; end
				end
			`,
		},
		"generic superclass matches": {
			input: `
				class Foo[V]; end

				class Bar < Foo[String]; end

				class Bar < Foo[String]
					def bar; end
				end
			`,
		},
		"generic superclass has a different type argument": {
			input: `
				class Foo[V]; end

				class Bar < Foo[String]; end

				class Bar < Foo[Int]
					def bar; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(74, 6, 17), P(81, 6, 24)), "superclass mismatch in `Bar`, got `Foo[Std::Int]`, expected `Foo[Std::String]`"),
			},
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
		"include mixin multiple times": {
			input: `
				mixin Foo; end
			  class Bar
					include Foo, Foo
					include Foo
				end
			`,
		},
		"include generic mixin multiple times with different type args": {
			input: `
				mixin Foo[V]; end
			  class Bar
					include Foo[String], Foo[Int]
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(64, 4, 27), P(71, 4, 34)), "cannot include mixin `Foo[Std::Int]` since `Foo[Std::String]` has already been included"),
			},
		},
		"include generic mixin multiple times with the same type args": {
			input: `
				mixin Foo[V]; end
			  class Bar
					include Foo[String], Foo[String]
					include Foo[String]
				end
			`,
		},
		"include generic mixin with type arguments in class": {
			input: `
				mixin Foo[V]; end
			  class Bar
					include Foo[String]
				end
			`,
		},
		"include generic mixin without type arguments in class": {
			input: `
				mixin Foo[V]; end
			  class Bar
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 4, 14), P(53, 4, 16)), "`Foo` requires 1 type argument(s), got: 0"),
			},
		},
		"include generic mixin with too many type arguments in class": {
			input: `
				mixin Foo[V]; end
			  class Bar
					include Foo[String, Int]
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 4, 14), P(53, 4, 16)), "`Foo` requires 1 type argument(s), got: 2"),
			},
		},
		"include generic mixin forwarding type arguments in class": {
			input: `
				mixin Foo[V]; end
			  class Bar[V]
					include Foo[V]
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
		"include generic mixin with type arguments in mixin": {
			input: `
				mixin Foo[V]; end
			  mixin Bar
					include Foo[String]
				end
			`,
		},
		"include generic mixin without type arguments in mixin": {
			input: `
				mixin Foo[V]; end
			  mixin Bar
					include Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 4, 14), P(53, 4, 16)), "`Foo` requires 1 type argument(s), got: 0"),
			},
		},
		"include generic mixin with too many type arguments in mixin": {
			input: `
				mixin Foo[V]; end
			  mixin Bar
					include Foo[String, Int]
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(51, 4, 14), P(53, 4, 16)), "`Foo` requires 1 type argument(s), got: 2"),
			},
		},
		"include generic mixin forwarding type arguments in mixin": {
			input: `
				mixin Foo[V]; end
			  mixin Bar[V]
					include Foo[V]
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
		"implement multiple times": {
			input: `
				interface Foo; end
			  class  Bar
					implement Foo, Foo
				end
			`,
		},
		"implement generic interface without type arguments in class": {
			input: `
				interface Foo[V]; end
			  class Bar
					implement Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 16), P(59, 4, 18)), "`Foo` requires 1 type argument(s), got: 0"),
			},
		},
		"implement generic interface with type arguments in class": {
			input: `
				interface Foo[V]; end
			  class Bar
					implement Foo[String]
				end
			`,
		},
		"implement generic interface with too many type arguments in class": {
			input: `
				interface Foo[V]; end
			  class Bar
					implement Foo[String, Int]
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 16), P(59, 4, 18)), "`Foo` requires 1 type argument(s), got: 2"),
			},
		},
		"implement generic interface forwarding type arguments in class": {
			input: `
				interface Foo[V]; end
			  class Bar[V]
					implement Foo[V]
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
		"implement generic interface without type arguments in mixin": {
			input: `
				interface Foo[V]; end
			  mixin Bar
					implement Foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 16), P(59, 4, 18)), "`Foo` requires 1 type argument(s), got: 0"),
			},
		},
		"implement generic interface with too many type arguments in mixin": {
			input: `
				interface Foo[V]; end
			  mixin Bar
					implement Foo[String, Int]
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 16), P(59, 4, 18)), "`Foo` requires 1 type argument(s), got: 2"),
			},
		},
		"implement generic interface with type arguments in mixin": {
			input: `
				interface Foo[V]; end
			  mixin Bar
					implement Foo[String]
				end
			`,
		},
		"implement generic interface forwarding type arguments in mixin": {
			input: `
				interface Foo[V]; end
			  mixin Bar[V]
					implement Foo[V]
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
		"assign to related generic mixin with correct type args": {
			input: `
				mixin Baz[K, V]; end

				mixin Bar[V]
					include Baz[String, V]
				end

				class Foo
					include Bar[Int]
				end

				var a: Bar[Int] = Foo()
				var b: Baz[String, Int] = a
			`,
		},
		"assign to related generic mixin with incorrect type args": {
			input: `
				mixin Baz[K, V]; end

				mixin Bar[V]
					include Baz[String, V]
				end

				class Foo
					include Bar[Int]
				end

				var a: Bar[Int] = Foo()
				var b: Baz[String, Float] = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(186, 13, 33), P(186, 13, 33)), "type `Bar[Std::Int]` cannot be assigned to type `Baz[Std::String, Std::Float]`"),
			},
		},
		"assign to distantly related generic mixin with correct type args": {
			input: `
				mixin Qux[E]; end
				mixin Baz[F]
					include Qux[F]
				end
				mixin Bar[K, V]
					include Baz[V]
				end
				mixin Foob[V]
					include Bar[String, V]
				end
				class Foo
					include Foob[Int]
				end

				var a: Foob[Int] = Foo()
				var b: Qux[Int] = a
			`,
		},
		"assign class to generic mixin": {
			input: `
				mixin Bar[V]; end
				class Foo
					include Bar[String]
				end

				var a: Bar[String] = Foo()
			`,
		},
		"assign class to generic mixin with wrong parameters": {
			input: `
				mixin Bar[V]; end
				class Foo
					include Bar[String]
				end

				var a: Bar[Int] = Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(93, 7, 23), P(97, 7, 27)), "type `Foo` cannot be assigned to type `Bar[Std::Int]`"),
			},
		},
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
		"redeclare non generic mixin as generic": {
			input: `
				mixin Foo; end
				mixin Foo[V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 3, 5), P(40, 3, 21)), "type parameter count mismatch in `Foo`, got: 1, expected: 0"),
			},
		},
		"redeclare generic mixin as non generic": {
			input: `
				mixin Foo[V]; end
				mixin Foo; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 5), P(40, 3, 18)), "type parameter count mismatch in `Foo`, got: 0, expected: 1"),
			},
		},
		"redeclare generic mixin with missing type parameter": {
			input: `
				mixin Foo[V, T]; end
				mixin Foo[V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(30, 3, 5), P(46, 3, 21)), "type parameter count mismatch in `Foo`, got: 1, expected: 2"),
			},
		},
		"redeclare generic mixin with additional type parameter": {
			input: `
				mixin Foo[V]; end
				mixin Foo[V, T]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 5), P(46, 3, 24)), "type parameter count mismatch in `Foo`, got: 2, expected: 1"),
			},
		},
		"redeclare generic mixin with matching type parameters": {
			input: `
				mixin Foo[+V < String]; end
				mixin Foo[+V < String]; end
			`,
		},
		"redeclare generic mixin with wrong type param name": {
			input: `
				mixin Foo[V]; end
				mixin Foo[T]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(27, 3, 5), P(43, 3, 21)), "type parameter mismatch in `Foo`, is `T`, should be `V`"),
			},
		},
		"redeclare generic mixin with wrong type param variance": {
			input: `
				mixin Foo[-V]; end
				mixin Foo[+V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 3, 5), P(45, 3, 22)), "type parameter mismatch in `Foo`, is `+V`, should be `-V`"),
			},
		},
		"redeclare generic mixin with wrong type param upper bound": {
			input: `
				mixin Foo[V < String]; end
				mixin Foo[V < Int]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 5), P(58, 3, 27)), "type parameter mismatch in `Foo`, is `V < Std::Int`, should be `V < Std::String`"),
			},
		},
		"redeclare generic mixin with wrong type param lower bound": {
			input: `
				mixin Foo[V > String]; end
				mixin Foo[V > Int]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 5), P(58, 3, 27)), "type parameter mismatch in `Foo`, is `V > Std::Int`, should be `V > Std::String`"),
			},
		},

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
				error.NewFailure(L("<main>", P(89, 6, 11), P(91, 6, 13)), "missing abstract method implementations in `Bar`:\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
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
				error.NewFailure(L("<main>", P(57, 5, 11), P(59, 5, 13)), "missing abstract method implementations in `Bar`:\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
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
				error.NewFailure(L("<main>", P(123, 10, 11), P(125, 10, 13)), "missing abstract method implementations in `Baz`:\n\n  - method `Bar.:bar`:\n      `def bar(): void`\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
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
				error.NewFailure(L("<main>", P(137, 10, 11), P(139, 10, 13)), "missing abstract method implementations in `Baz`:\n\n  - method `Bar.:bar`:\n      `def bar(): void`\n\n  - method `Foo.:foo`:\n      `def foo(): void`\n"),
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
				error.NewFailure(L("<main>", P(28, 3, 14), P(30, 3, 16)), "type `Foo` circularly references itself"),
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
				error.NewFailure(L("<main>", P(68, 7, 14), P(70, 7, 16)), "type `Foo` circularly references itself"),
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
		"redeclare non generic interface as generic": {
			input: `
				interface Foo; end
				interface Foo[V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 3, 5), P(48, 3, 25)), "type parameter count mismatch in `Foo`, got: 1, expected: 0"),
			},
		},
		"redeclare generic interface as non generic": {
			input: `
				interface Foo[V]; end
				interface Foo; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(31, 3, 5), P(48, 3, 22)), "type parameter count mismatch in `Foo`, got: 0, expected: 1"),
			},
		},
		"redeclare generic interface with missing type parameter": {
			input: `
				interface Foo[V, T]; end
				interface Foo[V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 5), P(54, 3, 25)), "type parameter count mismatch in `Foo`, got: 1, expected: 2"),
			},
		},
		"redeclare generic interface with additional type parameter": {
			input: `
				interface Foo[V]; end
				interface Foo[V, T]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(31, 3, 5), P(54, 3, 28)), "type parameter count mismatch in `Foo`, got: 2, expected: 1"),
			},
		},
		"redeclare generic interface with matching type parameters": {
			input: `
				interface Foo[+V < String]; end
				interface Foo[+V < String]; end
			`,
		},
		"redeclare generic interface with wrong type param name": {
			input: `
				interface Foo[V]; end
				interface Foo[T]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(31, 3, 5), P(51, 3, 25)), "type parameter mismatch in `Foo`, is `T`, should be `V`"),
			},
		},
		"redeclare generic interface with wrong type param variance": {
			input: `
				interface Foo[-V]; end
				interface Foo[+V]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(32, 3, 5), P(53, 3, 26)), "type parameter mismatch in `Foo`, is `+V`, should be `-V`"),
			},
		},
		"redeclare generic interface with wrong type param upper bound": {
			input: `
				interface Foo[V < String]; end
				interface Foo[V < Int]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(40, 3, 5), P(66, 3, 31)), "type parameter mismatch in `Foo`, is `V < Std::Int`, should be `V < Std::String`"),
			},
		},
		"redeclare generic interface with wrong type param lower bound": {
			input: `
				interface Foo[V > String]; end
				interface Foo[V > Int]; end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(40, 3, 5), P(66, 3, 31)), "type parameter mismatch in `Foo`, is `V > Std::Int`, should be `V > Std::String`"),
			},
		},

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
				error.NewFailure(L("<main>", P(34, 3, 16), P(36, 3, 18)), "type `Foo` circularly references itself"),
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
				error.NewFailure(L("<main>", P(80, 7, 16), P(82, 7, 18)), "type `Foo` circularly references itself"),
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
		"handle circular definitions in return type in implicit interface implementations": {
			input: `
				class Foo
					def foo: Foo then loop; end
				end
				interface Bar
					def foo: Bar; end
				end
				var a: Bar = Foo()
			`,
		},
		"handle circular definitions in return type in implicit generic interface implementations": {
			input: `
				interface C[E]
					def c: C[E]; end
				end

				class D
					def c: D then loop; end
				end

				var a: C[String] = D()
			`,
		},
		"handle circular definitions in return type in implicit generic interface implementations with a generic class": {
			input: `
				interface C[E]
					def c: C[E]; end
				end

				class D[E]
					def c: D[E] then loop; end
				end

				var a: C[String] = D::[String]()
			`,
		},
		"handle circular definitions in param type in implicit interface implementations": {
			input: `
				class Foo
					def foo(a: Bar); end
				end
				interface Bar
					def foo(a: Foo); end
				end
				var a: Bar = Foo()
			`,
		},
		"assign to related generic interface with correct type args": {
			input: `
				interface Baz[K, V]
					def baz: V; end
				end

				interface Bar[V]
					implement Baz[String, V]
				end

				class Foo
					implement Bar[Int]

					def baz: Int then 3
				end

				var a: Bar[Int] = Foo()
				var b: Baz[String, Int] = a
			`,
		},
		"assign to implicit generic interface with correct type args": {
			input: `
				interface Baz[K, V]
					def baz: V; end
				end

				interface Bar[V]
					implement Baz[String, V]
				end

				class Foo
					def baz: Int then 3
				end

				var a: Bar[Int] = Foo()
				var b: Baz[String, Int] = a
			`,
		},
		"assign to related generic interface with incorrect type args": {
			input: `
				interface Baz[K, V]
					def baz: V; end
				end

				interface Bar[V]
					implement Baz[String, V]
				end

				class Foo
					implement Bar[Int]

					def baz: Int then 3
				end

				var a: Bar[Int] = Foo()
				var b: Baz[String, Float] = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(248, 17, 33), P(248, 17, 33)), "type `Bar[Std::Int]` cannot be assigned to type `Baz[Std::String, Std::Float]`"),
			},
		},
		"assign to implicit generic interface with incorrect type args": {
			input: `
				interface Baz[K, V]
					def baz: V; end
				end

				interface Bar[V]
					implement Baz[String, V]
				end

				class Foo
					def baz: Int then 3
				end

				var a: Bar[Int] = Foo()
				var b: Baz[String, Float] = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(223, 15, 33), P(223, 15, 33)), "type `Bar[Std::Int]` cannot be assigned to type `Baz[Std::String, Std::Float]`"),
			},
		},
		"assign to distantly related generic interface with correct type args": {
			input: `
				interface Qux[E]; end
				interface Baz[F]
					implement Qux[F]
				end
				interface Bar[K, V]
					implement Baz[V]
				end
				interface Foob[V]
					implement Bar[String, V]
				end
				class Foo
					implement Foob[Int]
				end

				var a: Foob[Int] = Foo()
				var b: Qux[Int] = a
			`,
		},
		"assign to distantly implicitly related generic interface with correct type args": {
			input: `
				interface Qux[E]; end
				interface Baz[F]
					implement Qux[F]
				end
				interface Bar[K, V]
					implement Baz[V]
				end
				interface Foob[V]
					implement Bar[String, V]
				end
				class Foo; end

				var a: Foob[Int] = Foo()
				var b: Qux[Int] = a
			`,
		},

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
				error.NewFailure(L("<main>", P(82, 7, 18), P(86, 7, 22)), "type `Bar` does not implement interface `Foo`:\n\n  - missing method `Foo.:foo` with signature: `def foo(): void`"),
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
				error.NewFailure(L("<main>", P(189, 14, 18), P(189, 14, 18)), "type `Bar` does not implement interface `Foo`:\n\n  - missing method `Foo.:foo` with signature: `def foo(): void`"),
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

func TestExtendWhere(t *testing.T) {
	tests := testTable{
		"declare extend in the top level": {
			input: `extend where T < String; end`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(27, 1, 28)), "cannot declare extend blocks in this context"),
			},
		},
		"declare extend in a module": {
			input: `
				module Foo
					extend where T < String; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(21, 3, 6), P(48, 3, 33)), "cannot declare extend blocks in this context"),
			},
		},
		"declare extend in an interface": {
			input: `
				interface Foo
					extend where T < String; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(24, 3, 6), P(51, 3, 33)), "cannot declare extend blocks in this context"),
			},
		},
		"declare extend in a non-generic class": {
			input: `
				class Foo
					extend where T < String; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(47, 3, 33)), "cannot use `extend where` since namespace `Foo` is not generic"),
			},
		},

		"declare extend with a nonexistent type parameter": {
			input: `
				class Foo[T]
					extend where W < String; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 19), P(45, 3, 28)), "cannot add where constraints to nonexistent type parameter `W`"),
			},
		},
		"declare extend with an invalid type parameter upper bound": {
			input: `
				class Foo[T < String]
					extend where T < Int; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 19), P(51, 3, 25)), "type parameter `T` in where clause should have a narrower upper bound, has `Std::Int`, should have `Std::String` or its subtype"),
			},
		},
		"declare extend with a narrower type parameter upper bound": {
			input: `
				class Foo[T < Value]
					extend where T < Int; end
				end
			`,
		},
		"declare extend with a wider type parameter upper bound": {
			input: `
				class Foo[T < Int]
					extend where T < Value; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 3, 19), P(50, 3, 27)), "type parameter `T` in where clause should have a narrower upper bound, has `Std::Value`, should have `Std::Int` or its subtype"),
			},
		},
		"declare extend with a narrower type parameter lower bound": {
			input: `
				class Foo[T > Value]
					extend where T > Int; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 3, 19), P(50, 3, 25)), "type parameter `T` in where clause should have a wider lower bound, has `Std::Int`, should have `Std::Value` or its supertype"),
			},
		},
		"declare extend with a wider type parameter lower bound": {
			input: `
				class Foo[T > Int]
					extend where T > Value; end
				end
			`,
		},
		"assign class with extend where to an interface when the conditions are satisfied": {
			input: `
				interface Foo
					sig foo: Int
				end

				class Bar[T]
					extend where T < Int
						def foo: Int then 3
					end
				end

				var a: Foo = Bar::[Int]()
			`,
		},
		"assign class with extend where to a generic interface when the conditions are satisfied": {
			input: `
				interface Foo[T]
					sig foo: T
				end

				class Bar[T]
					extend where T < Int
						def foo: Int then 3
					end
				end

				var a: Foo[Int] = Bar::[Int]()
			`,
		},
		"assign class with extend where to an interface when the conditions are not satisfied": {
			input: `
				interface Foo
					sig foo: Int
				end

				class Bar[T]
					extend where T < Int
						def foo: Int then 3
					end
				end

				var a: Foo = Bar::[String]()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(150, 12, 18), P(164, 12, 32)), "type `Bar[Std::String]` does not implement interface `Foo`:\n\n  - missing method `Foo.:foo` with signature: `def foo(): Std::Int`"),
				error.NewFailure(L("<main>", P(150, 12, 18), P(164, 12, 32)), "type `Bar[Std::String]` cannot be assigned to type `Foo`"),
			},
		},
		"assign class with extend where to a generic interface when the conditions are not satisfied": {
			input: `
				interface Foo[T]
					sig foo: T
				end

				class Bar[T]
					extend where T < Int
						def foo: Int then 3
					end
				end

				var a: Foo[Int] = Bar::[String]()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(156, 12, 23), P(170, 12, 37)), "type `Bar[Std::String]` does not implement interface `Foo[Std::Int]`:\n\n  - missing method `Foo.:foo` with signature: `def foo(): Std::Int`"),
				error.NewFailure(L("<main>", P(156, 12, 23), P(170, 12, 37)), "type `Bar[Std::String]` cannot be assigned to type `Foo[Std::Int]`"),
			},
		},

		"declare extend in a non-generic mixin": {
			input: `
				mixin Foo
					extend where T < String; end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(20, 3, 6), P(47, 3, 33)), "cannot use `extend where` since namespace `Foo` is not generic"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
