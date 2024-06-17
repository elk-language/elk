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
		"include in class": {
			input: `
				mixin Foo; end
			  class  Bar
					include Foo
				end
			`,
		},
		"include in mixin": {
			input: `
				mixin Foo; end
			  mixin  Bar
					include Foo
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
