package checker

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

func TestTypeParameters(t *testing.T) {
	tests := testTable{
		"use invariant type parameter in instance variable": {
			input: `
				class Foo[V]
					var @a: V?
				end
			`,
		},
		"use covariant type parameter in instance variable": {
			input: `
				class Foo[+V]
					var @a: V?
				end
			`,
		},
		"use contravariant type parameter in instance variable": {
			input: `
				class Foo[-V]
					var @a: V?
				end
			`,
		},
		"use type parameter of enclosing namespace in an instance variable": {
			input: `
				class Foo[V]
					class Bar[W]
						var @a: V?
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(50, 4, 15), P(50, 4, 15)), "undefined type `V`"),
			},
		},

		"use type parameter outside of methods": {
			input: `
				class Foo[V]
					var a: V
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(30, 3, 13), P(30, 3, 13)), "type parameter `V` cannot be used in this context"),
			},
		},
		"use type parameter of enclosing namespace in a method parameter": {
			input: `
				class Foo[V]
					class Bar[W]
						def baz(a: V); end
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(53, 4, 18), P(53, 4, 18)), "undefined type `V`"),
			},
		},
		"use type parameter of enclosing namespace in a method return type": {
			input: `
				class Foo[V]
					class Bar[W]
						def baz: V; end
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(51, 4, 16), P(51, 4, 16)), "undefined type `V`"),
			},
		},
		"use type parameter of enclosing namespace in a method throw type": {
			input: `
				class Foo[V]
					class Bar[W]
						def baz! V; end
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(51, 4, 16), P(51, 4, 16)), "undefined type `V`"),
			},
		},
		"use type parameter of enclosing namespace in a method body": {
			input: `
				class Foo[V]
					class Bar[W]
						def baz
							var a: V
						end
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(64, 5, 15), P(64, 5, 15)), "undefined type `V`"),
			},
		},
		"use type parameter of namespace in a method body": {
			input: `
				class Bar[V]
					def baz
						var a: V
					end
				end
			`,
		},
		"use type parameter of namespace and method in a method body": {
			input: `
				class Bar[V]
					def baz[W]
						var a: W
						var b: V
					end
				end
			`,
		},

		"invariant type parameter in method parameters": {
			input: `
				class Foo[V]
					def foo(a: V); end
				end
			`,
		},
		"invariant type parameter in method return type": {
			input: `
				class Foo[V]
					def foo(a: V): V then a
				end
			`,
		},
		"invariant type parameter in method throw type": {
			input: `
				class Foo[V]
					def foo(a: V)! V; end
				end
			`,
		},

		"covariant type parameter in method parameters": {
			input: `
				class Foo[+V]
					def foo(a: V); end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(35, 3, 17), P(35, 3, 17)), "covariant type parameter `V` cannot appear in input positions"),
			},
		},
		"covariant type parameter in method return type": {
			input: `
				class Foo[+V]
					def foo: V then loop; end
				end
			`,
		},
		"covariant type parameter in method throw type": {
			input: `
				class Foo[+V]
					def foo! V; end
				end
			`,
		},
		"covariant type parameter in closure param's param type": {
			input: `
				class Foo[+V]
					def foo(a: |a: V|); end
				end
			`,
		},
		"covariant type parameter in closure param's return type": {
			input: `
				class Foo[+V]
					def foo(a: ||: V); end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(39, 3, 21), P(39, 3, 21)), "covariant type parameter `V` cannot appear in input positions"),
			},
		},
		"covariant type parameter in closure return type's param type": {
			input: `
				class Foo[+V]
					def foo(): |a: V| then loop; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(39, 3, 21), P(39, 3, 21)), "covariant type parameter `V` cannot appear in input positions"),
			},
		},
		"covariant type parameter in closure return type's return type": {
			input: `
				class Foo[+V]
					def foo(): ||: V then loop; end
				end
			`,
		},

		"contravariant type parameter in method parameters": {
			input: `
				class Foo[-V]
					def foo(a: V); end
				end
			`,
		},
		"contravariant type parameter in method return type": {
			input: `
				class Foo[-V]
					def foo: V; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(33, 3, 15), P(33, 3, 15)), "contravariant type parameter `V` cannot appear in output positions"),
			},
		},
		"contravariant type parameter in method throw type": {
			input: `
				class Foo[-V]
					def foo! V; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(33, 3, 15), P(33, 3, 15)), "contravariant type parameter `V` cannot appear in output positions"),
			},
		},
		"contravariant type parameter in closure param's param type": {
			input: `
				class Foo[-V]
					def foo(a: |a: V|); end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(39, 3, 21), P(39, 3, 21)), "contravariant type parameter `V` cannot appear in output positions"),
			},
		},
		"contravariant type parameter in closure param's return type": {
			input: `
				class Foo[-V]
					def foo(a: ||: V); end
				end
			`,
		},
		"contravariant type parameter in closure return type's param type": {
			input: `
				class Foo[-V]
					def foo(): |a: V| then loop; end
				end
			`,
		},
		"contravariant type parameter in closure return type's return type": {
			input: `
				class Foo[-V]
					def foo(): ||: V then loop; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(39, 3, 21), P(39, 3, 21)), "contravariant type parameter `V` cannot appear in output positions"),
			},
		},

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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(107, 7, 20), P(107, 7, 20)), "type `V` cannot be assigned to type `Bar`"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(74, 6, 18), P(78, 6, 22)), "type `Foo` cannot be assigned to type `V`"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(99, 7, 18), P(103, 7, 22)), "type `Foo` cannot be assigned to type `V`"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(99, 7, 18), P(103, 7, 22)), "type `Bar` cannot be assigned to type `V`"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(82, 6, 20), P(82, 6, 20)), "type `V` cannot be assigned to type `Foo`"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(107, 7, 20), P(107, 7, 20)), "type `V` cannot be assigned to type `Foo`"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(107, 7, 20), P(107, 7, 20)), "type `V` cannot be assigned to type `Bar`"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(99, 7, 18), P(103, 7, 22)), "type `Foo` cannot be assigned to type `V`"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(96, 8, 16), P(121, 8, 41)), "type `&Std::String` cannot be assigned to type `9`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
