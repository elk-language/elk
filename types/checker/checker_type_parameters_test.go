package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestTypeParameters(t *testing.T) {
	tests := testTable{
		"assign type variable to its upper bound": {
			input: `
				class Foo; end

				class Baz[V < Foo]
					def foo(a: V)
						var b: Foo = a
					end
				end
			`,
		},
		"assign type variable to the parent of its upper bound": {
			input: `
				class Foo; end
				class Bar < Foo; end

				class Baz[V < Bar]
					def foo(a: V)
						var b: Foo = a
					end
				end
			`,
		},
		"assign type variable to the child of its upper bound": {
			input: `
				class Foo; end
				class Bar < Foo; end

				class Baz[V < Foo]
					def foo(a: V)
						var b: Bar = a
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(107, 7, 20), P(107, 7, 20)), "type `V` cannot be assigned to type `Bar`"),
			},
		},

		"assign upper bound to type variable": {
			input: `
				class Foo; end

				class Baz[V < Foo]
					def foo
						var a: V = Foo()
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(74, 6, 18), P(78, 6, 22)), "type `Foo` cannot be assigned to type `V`"),
			},
		},
		"assign parent of upper bound to type variable": {
			input: `
				class Foo; end
				class Bar < Foo; end

				class Baz[V < Bar]
					def foo
						var a: V = Foo()
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(99, 7, 18), P(103, 7, 22)), "type `Foo` cannot be assigned to type `V`"),
			},
		},
		"assign child of upper bound to type variable": {
			input: `
				class Foo; end
				class Bar < Foo; end

				class Baz[V < Foo]
					def foo
						var a: V = Bar()
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(99, 7, 18), P(103, 7, 22)), "type `Bar` cannot be assigned to type `V`"),
			},
		},

		"assign type variable to its lower bound": {
			input: `
				class Foo; end

				class Baz[V > Foo]
					def foo(a: V)
						var b: Foo = a
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(82, 6, 20), P(82, 6, 20)), "type `V` cannot be assigned to type `Foo`"),
			},
		},
		"assign type variable to the parent of its lower bound": {
			input: `
				class Foo; end
				class Bar < Foo; end

				class Baz[V > Bar]
					def foo(a: V)
						var b: Foo = a
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(107, 7, 20), P(107, 7, 20)), "type `V` cannot be assigned to type `Foo`"),
			},
		},
		"assign type variable to the child of its lower bound": {
			input: `
				class Foo; end
				class Bar < Foo; end

				class Baz[V > Foo]
					def foo(a: V)
						var b: Bar = a
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(107, 7, 20), P(107, 7, 20)), "type `V` cannot be assigned to type `Bar`"),
			},
		},

		"assign lower bound to type variable": {
			input: `
				class Foo; end

				class Baz[V > Foo]
					def foo
						var a: V = Foo()
					end
				end
			`,
		},
		"assign parent of lower bound to type variable": {
			input: `
				class Foo; end
				class Bar < Foo; end

				class Baz[V > Bar]
					def foo
						var a: V = Foo()
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(99, 7, 18), P(103, 7, 22)), "type `Foo` cannot be assigned to type `V`"),
			},
		},
		"assign child of lower bound to type variable": {
			input: `
				class Foo; end
				class Bar < Foo; end

				class Baz[V > Foo]
					def foo
						var a: V = Bar()
					end
				end
			`,
		},
		"call a method on a type parameter and return its singleton type": {
			input: `
				class Lol[V < Value]
					def foo(a: V): &V
						a.class
					end
				end

				var a: 9 = Lol::[String]().foo("bar")
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(96, 8, 16), P(121, 8, 41)), "type `&Std::String` cannot be assigned to type `9`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
