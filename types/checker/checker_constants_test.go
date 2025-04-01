package checker

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

func TestConstantAccess(t *testing.T) {
	tests := testTable{
		"access class constant": {
			input: "Int",
		},
		"access module constant": {
			input: "Std",
		},
		"access undefined constant": {
			input: "Foo",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(2, 1, 3)), "undefined constant `Foo`"),
			},
		},
		"constant lookup": {
			input: "Std::Int",
		},
		"constant lookup with error in the middle": {
			input: "Std::Foo::Bar",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(5, 1, 6), P(7, 1, 8)), "undefined constant `Std::Foo`"),
			},
		},
		"constant lookup with error at the start": {
			input: "Foo::Bar::Baz",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(2, 1, 3)), "undefined constant `Foo`"),
			},
		},
		"absolute constant lookup": {
			input: "::Std::Int",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestConstantDeclarations(t *testing.T) {
	tests := testTable{
		"declare in extend where": {
			input: `
				class E[T]
					extend where T < String
						const D = 3
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(51, 4, 7), P(61, 4, 17)), "constants cannot be declared in this context"),
			},
		},
		"declare in a module": {
			input: `
				module F
					const D = 3
				end
			`,
		},
		"declare in a mixin": {
			input: `
				mixin F
					const D = 3
				end
			`,
		},
		"declare in a class": {
			input: `
				class F
					const D = 3
				end
			`,
		},
		"declare in a singleton": {
			input: `
				class F
					singleton
						const D = 3
					end
				end
			`,
		},
		"declare with explicit type": {
			input: "const Foo: Int = 5",
		},
		"declare with implicit type": {
			input: "const Foo = 5",
		},
		"declare with implicit type and assign to literal type": {
			input: `
				const Foo = 5
				var foo: 5 = Foo
			`,
		},
		"declare with incorrect explicit type": {
			input: "const Foo: String = 5",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 1, 21), P(20, 1, 21)), "type `5` cannot be assigned to type `Std::String`"),
			},
		},
		"declare without initialising": {
			input: "const Foo: String",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(0, 1, 1), P(16, 1, 17)), "constants must be initialised"),
			},
		},
		"declare with a static initialiser without an explicit type": {
			input: "const Foo = 3",
		},
		"declare with a circular reference in method": {
			input: `
				const Bar: String = Foo.foo
				module Foo
					def foo: String
						Bar
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L("<main>", P(53, 4, 6), P(86, 6, 8)),
					"method `Foo::foo` circularly refers to constant `Bar` because it gets called in its initializer",
				),
			},
		},
		"declare with a circular reference in method to related constant": {
			input: `
				const Baz: String = Bar + "lol"
				const Bar: String = Foo.foo
				module Foo
					def foo: String
						Baz
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L("<main>", P(89, 5, 6), P(122, 7, 8)),
					"method `Foo::foo` circularly refers to constant `Baz` because it gets called in its initializer",
				),
			},
		},
		"declare with a circular reference in method to related constant that is already defined": {
			input: `
				const Bar: String = Foo.foo
				const Baz: String = Bar + "lol"
				module Foo
					def foo: String
						Baz
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L("<main>", P(89, 5, 6), P(122, 7, 8)),
					"method `Foo::foo` circularly refers to constant `Baz` because it gets called in its initializer",
				),
			},
		},
		"declare with a circular reference in related method to related constant": {
			input: `
				const Baz: String = Bar + "lol"
				const Bar: String = Foo.foo
				module Foo
					def foo: String
						bar()
					end

					def bar: String
						Baz
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L("<main>", P(132, 9, 6), P(165, 11, 8)),
					"method `Foo::bar` circularly refers to constant `Baz` because it gets called in its initializer",
				),
			},
		},
		"declare with a dynamic initialiser": {
			input: `
				const Bar: String = Foo.foo
				module Foo
					def foo: String
						"foo"
					end
				end
			`,
		},
		"declare with a dynamic initialiser without an explicit type": {
			input: `
				module Foo
					def foo: String
						"foo"
					end
				end
				const Bar = Foo.foo
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(70, 7, 5), P(88, 7, 23)), "non-static constants must have an explicit type"),
			},
		},
		"declare with a simple circular reference": {
			input: `
				const Foo: String = Foo
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(25, 2, 25), P(27, 2, 27)), "constant `Foo` circularly references itself"),
			},
		},
		"declare with a complex circular reference": {
			input: `
				const Foo: String = Bar
				const Bar: String = Baz
				const Baz: String = Foo
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(81, 4, 25), P(83, 4, 27)), "constant `Foo` circularly references itself"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
