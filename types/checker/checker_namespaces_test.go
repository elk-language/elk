package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestModule(t *testing.T) {
	tests := testTable{
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
				error.NewError(L("<main>", P(7, 1, 8), P(9, 1, 10)), "undefined namespace `Int`"),
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
				error.NewError(L("<main>", P(7, 1, 8), P(9, 1, 10)), "undefined namespace `Foo`"),
			},
		},
		"module with undefined constant in the middle": {
			input: `module Std::Foo::Bar; end`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "undefined namespace `Std::Foo`"),
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
				error.NewError(L("<main>", P(49, 5, 5), P(51, 5, 7)), "undefined constant `Bar`"),
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
				error.NewError(L("<main>", P(21, 3, 6), P(38, 4, 8)), "cannot declare a singleton class in this context"),
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
		"class with public constant": {
			input: `class Foo; end`,
		},
		"class with nonexistent superclass": {
			input: `class Foo < Bar; end`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "undefined type `Bar`"),
				error.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "`void` is not a class"),
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
				error.NewError(L("<main>", P(43, 3, 17), P(45, 3, 19)), "cannot inherit from sealed class `Bar`"),
			},
		},
		"class with module superclass": {
			input: `
				module Bar; end
				class Foo < Bar; end
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(37, 3, 17), P(39, 3, 19)), "`Bar` is not a class"),
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
				error.NewError(L("<main>", P(89, 6, 11), P(91, 6, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(177, 10, 11), P(179, 10, 13)), "missing abstract method implementation `Bar.:bar` with signature: sig bar(): void"),
				error.NewError(L("<main>", P(177, 10, 11), P(179, 10, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(156, 11, 11), P(158, 11, 13)), "missing abstract method implementation `Bar.:bar` with signature: sig bar(): void"),
				error.NewError(L("<main>", P(156, 11, 11), P(158, 11, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(92, 6, 11), P(94, 6, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(189, 12, 11), P(191, 12, 13)), "missing abstract method implementation `Bar.:bar` with signature: sig bar(): void"),
				error.NewError(L("<main>", P(189, 12, 11), P(191, 12, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(156, 11, 11), P(158, 11, 13)), "missing abstract method implementation `Bar.:bar` with signature: sig bar(): void"),
				error.NewError(L("<main>", P(156, 11, 11), P(158, 11, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(57, 5, 11), P(59, 5, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(123, 10, 11), P(125, 10, 13)), "missing abstract method implementation `Bar.:bar` with signature: sig bar(): void"),
				error.NewError(L("<main>", P(123, 10, 11), P(125, 10, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(60, 6, 19), P(60, 6, 19)), "type `Foo` cannot be assigned to type `&Foo`"),
			},
		},
		"assign class to a class singleton type": {
			input: `
				class Foo; end

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
				error.NewError(L("<main>", P(51, 6, 5), P(100, 8, 7)), "cannot redeclare class `Bar` with a different modifier, is `abstract`, should be `default`"),
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
				error.NewError(L("<main>", P(51, 6, 5), P(98, 8, 7)), "cannot redeclare class `Bar` with a different modifier, is `sealed`, should be `default`"),
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
				error.NewError(L("<main>", P(60, 6, 5), P(107, 8, 7)), "cannot redeclare class `Bar` with a different modifier, is `sealed`, should be `abstract`"),
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
				error.NewError(L("<main>", P(60, 6, 5), P(100, 8, 7)), "cannot redeclare class `Bar` with a different modifier, is `default`, should be `abstract`"),
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
				error.NewError(L("<main>", P(51, 6, 5), P(85, 8, 7)), "superclass mismatch in `Bar`, got `Std::Object`, expected `Foo`"),
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
				error.NewError(L("<main>", P(8, 1, 9), P(10, 1, 11)), "undefined type `Foo`"),
				error.NewError(L("<main>", P(8, 1, 9), P(10, 1, 11)), "only mixins can be included"),
			},
		},
		"include in top level": {
			input: `
				mixin Foo; end
				include Foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(32, 3, 13), P(34, 3, 15)), "cannot include mixins in this context"),
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
				error.NewError(L("<main>", P(49, 4, 14), P(51, 4, 16)), "cannot include mixins in this context"),
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
				error.NewError(L("<main>", P(52, 4, 14), P(54, 4, 16)), "cannot include mixins in this context"),
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
				error.NewError(L("<main>", P(49, 4, 14), P(51, 4, 16)), "only mixins can be included"),
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
				error.NewError(L("<main>", P(48, 4, 14), P(50, 4, 16)), "only mixins can be included"),
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
				error.NewError(L("<main>", P(52, 4, 14), P(54, 4, 16)), "only mixins can be included"),
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
			input: `implement Foo`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(10, 1, 11), P(12, 1, 13)), "undefined type `Foo`"),
				error.NewError(L("<main>", P(10, 1, 11), P(12, 1, 13)), "only interfaces can be implemented"),
			},
		},
		"implement in top level": {
			input: `
				interface Foo; end
				implement Foo
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(38, 3, 15), P(40, 3, 17)), "cannot implement interfaces in this context"),
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
				error.NewError(L("<main>", P(55, 4, 16), P(57, 4, 18)), "cannot implement interfaces in this context"),
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
				error.NewError(L("<main>", P(51, 4, 16), P(53, 4, 18)), "only interfaces can be implemented"),
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
				error.NewError(L("<main>", P(50, 4, 16), P(52, 4, 18)), "only interfaces can be implemented"),
			},
		},
		"include mixin": {
			input: `
				mixin Foo; end
				class  Bar
					implement Foo
				end
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(50, 4, 16), P(52, 4, 18)), "only interfaces can be implemented"),
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
				error.NewError(L("<main>", P(57, 5, 18), P(61, 5, 22)), "type `Foo` cannot be assigned to type `Bar`"),
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
				error.NewError(L("<main>", P(102, 9, 19), P(102, 9, 19)), "type `Foo` cannot be assigned to type `&Foo`"),
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
				error.NewError(L("<main>", P(24, 3, 5), P(46, 3, 27)), "cannot redeclare mixin `Bar` with a different modifier, is `abstract`, should be `default`"),
			},
		},
		"modifier was abstract, is default": {
			input: `
				abstract mixin Bar; end
				mixin Bar; end
			`,
			err: error.ErrorList{
				error.NewError(L("<main>", P(33, 3, 5), P(46, 3, 18)), "cannot redeclare mixin `Bar` with a different modifier, is `default`, should be `abstract`"),
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
				error.NewError(L("<main>", P(89, 6, 11), P(91, 6, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(57, 5, 11), P(59, 5, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(123, 10, 11), P(125, 10, 13)), "missing abstract method implementation `Bar.:bar` with signature: sig bar(): void"),
				error.NewError(L("<main>", P(123, 10, 11), P(125, 10, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(137, 10, 11), P(139, 10, 13)), "missing abstract method implementation `Bar.:bar` with signature: sig bar(): void"),
				error.NewError(L("<main>", P(137, 10, 11), P(139, 10, 13)), "missing abstract method implementation `Foo.:foo` with signature: sig foo(): void"),
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
				error.NewError(L("<main>", P(111, 11, 5), P(117, 11, 11)), "method `foo` is not defined on type `&Bar`"),
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
				error.NewError(L("<main>", P(121, 11, 5), P(127, 11, 11)), "method `foo` is not defined on type `&Bar`"),
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
				error.NewError(L("<main>", P(82, 7, 18), P(86, 7, 22)), "type `Bar` does not implement interface `Foo`:\n\n  - missing method `Foo.:foo` with signature: sig foo(): void\n"),
				error.NewError(L("<main>", P(82, 7, 18), P(86, 7, 22)), "type `Bar` cannot be assigned to type `Foo`"),
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
				error.NewError(L("<main>", P(189, 14, 18), P(189, 14, 18)), "type `Bar` does not implement interface `Foo`:\n\n  - missing method `Foo.:foo` with signature: sig foo(): void\n"),
				error.NewError(L("<main>", P(189, 14, 18), P(189, 14, 18)), "type `Bar` cannot be assigned to type `Foo`"),
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
				error.NewError(L("<main>", P(85, 6, 19), P(85, 6, 19)), "type `Foo` cannot be assigned to type `&Foo`"),
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
