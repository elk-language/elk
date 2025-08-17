package checker_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

func TestUsing(t *testing.T) {
	tests := testTable{
		"not a namespace with star": {
			input: `
				typedef Lol = 3
				using Lol::*
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(31, 3, 11), P(36, 3, 16)), "type `Lol` is not a namespace"),
			},
		},
		"undefined type with star": {
			input: `
				using Lol::*
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(11, 2, 11), P(16, 2, 16)), "undefined namespace `Lol`"),
			},
		},
		"star in top level": {
			input: `
				using Foo::*

				var a: Bar = 3

				class Foo
				 	class Bar; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(36, 4, 18), P(36, 4, 18)), "type `3` cannot be assigned to type `Foo::Bar`"),
			},
		},
		"using with star in module, resolve in methods": {
			input: `
				module Baz
					using Foo::*

					def baz: Bar then Bar()
				end

				class Foo
				 	class Bar; end
				end

				var a: 9 = Baz.baz
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(132, 12, 16), P(138, 12, 22)), "type `Foo::Bar` cannot be assigned to type `9`"),
			},
		},
		"using with multiple namespaces with star": {
			input: `
				class Foo
					class Bar; end
				end
				class Lol
					class Grub; end
				end

				module Baz
					using Foo::*, Lol::*

					var a: 9 = Bar()
					var b: 12 = Grub()
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(145, 12, 17), P(149, 12, 21)), "type `Foo::Bar` cannot be assigned to type `9`"),
				diagnostic.NewFailure(L("<main>", P(168, 13, 18), P(173, 13, 23)), "type `Lol::Grub` cannot be assigned to type `12`"),
			},
		},
		"using goes out of scope": {
			input: `
				module Baz
					using Foo::*
				end

				class Foo
				 	class Bar; end
				end

				Bar
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(91, 10, 5), P(93, 10, 7)), "undefined constant `Bar`"),
			},
		},
		"using only accepts absolute constants": {
			input: `
				module Baz
					class Foo; end
					using Foo::*
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(47, 4, 12), P(52, 4, 17)), "undefined namespace `Foo`"),
			},
		},

		"using with a few namespaces": {
			input: `
				using Foo::{Bar, Baz}

				class Foo
					class Bar; end
					class Baz; end
				end

				var a: Bar = Bar()
				var b: Bar = 9
				var c: &Bar = Bar

				var d: Baz = Baz()
				var e: Baz = 9
				var f: &Baz = Baz
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(131, 10, 18), P(131, 10, 18)), "type `9` cannot be assigned to type `Foo::Bar`"),
				diagnostic.NewFailure(L("<main>", P(196, 14, 18), P(196, 14, 18)), "type `9` cannot be assigned to type `Foo::Baz`"),
			},
		},
		"using with a few namespaces and as": {
			input: `
				using Foo::{Bar as R, Baz as Z}

				class Foo
					class Bar; end
					class Baz; end
				end

				var a: R = R()
				var b: R = 9
				var c: &R = R

				var d: Z = Z()
				var e: Z = 9
				var f: &Z = Z
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(135, 10, 16), P(135, 10, 16)), "type `9` cannot be assigned to type `Foo::Bar`"),
				diagnostic.NewFailure(L("<main>", P(190, 14, 16), P(190, 14, 16)), "type `9` cannot be assigned to type `Foo::Baz`"),
			},
		},

		"using with a single namespace": {
			input: `
				using Foo::Bar

				class Foo
					class Bar; end
				end

				var a: Bar = Bar()
				var b: Bar = 9
				var c: &Bar = Bar
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(104, 9, 18), P(104, 9, 18)), "type `9` cannot be assigned to type `Foo::Bar`"),
			},
		},
		"using with a single namespace after its declaration": {
			input: `
				class Foo
					class Bar; end
				end

				using Foo::Bar

				var a: Bar = Bar()
				var b: Bar = 9
				var c: &Bar = Bar
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(104, 9, 18), P(104, 9, 18)), "type `9` cannot be assigned to type `Foo::Bar`"),
			},
		},
		"using with a single namespace and as": {
			input: `
				using Foo::Bar as B

				class Foo
					class Bar; end
				end

				var a: B = B()
				var b: B = 9
				var c: &B = B
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(103, 9, 16), P(103, 9, 16)), "type `9` cannot be assigned to type `Foo::Bar`"),
			},
		},
		"using with a single constant": {
			input: `
				using Foo::Bar

				class Foo
					const Bar = 3
				end

				var a: Bar = 9
				var c: 3 = Bar
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(74, 8, 12), P(76, 8, 14)), "undefined type `Foo::Bar`"),
			},
		},
		"using with a single macro in braces": {
			input: `
				using Foo::{bar!}

				class Foo
					macro bar
						3.to_ast_expr_node
					end
				end

				var a: nil = bar!
				var c: 3 = bar!
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(113, 10, 18), P(115, 10, 20)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"using with a single macro in braces with as": {
			input: `
				using Foo::{bar! as b!}

				class Foo
					macro bar
						3.to_ast_expr_node
					end
				end

				var a: nil = b!
				var c: 3 = b!
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(119, 10, 18), P(119, 10, 18)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"using with a single macro": {
			input: `
				using Foo::bar!

				class Foo
					macro bar
						3.to_ast_expr_node
					end
				end

				var a: nil = bar!
				var c: 3 = bar!
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(111, 10, 18), P(113, 10, 20)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"using with a single macro and as": {
			input: `
				using Foo::bar! as b!

				class Foo
					macro bar
						3.to_ast_expr_node
					end
				end

				var a: nil = b!
				var c: 3 = b!
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(117, 10, 18), P(117, 10, 18)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"using with a single constant and as": {
			input: `
				using Foo::Bar as B

				class Foo
					const Bar = 3
				end

				var a: B = 9
				var c: 3 = B
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(79, 8, 12), P(79, 8, 12)), "undefined type `Foo::Bar`"),
			},
		},
		"using with a single constant and as after declaration": {
			input: `
				class Foo
					const Bar = 3
				end
				using Foo::Bar as B

				var a: B = 9
				var c: 3 = B
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(78, 7, 12), P(78, 7, 12)), "undefined type `Foo::Bar`"),
			},
		},
		"using with a single type": {
			input: `
				using Foo::Bar

				class Foo
					typedef Bar = 3
				end

				var a: Bar = 9
				var c: 3 = Bar
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(82, 8, 18), P(82, 8, 18)), "type `9` cannot be assigned to type `Foo::Bar`"),
				diagnostic.NewFailure(L("<main>", P(99, 9, 16), P(101, 9, 18)), "`Foo::Bar` cannot be used as a value in expressions"),
			},
		},
		"using with a single type and as": {
			input: `
				using Foo::Bar as B

				class Foo
					typedef Bar = 3
				end

				var a: B = 9
				var c: 3 = B
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(85, 8, 16), P(85, 8, 16)), "type `9` cannot be assigned to type `Foo::Bar`"),
				diagnostic.NewFailure(L("<main>", P(102, 9, 16), P(102, 9, 16)), "`Foo::Bar` cannot be used as a value in expressions"),
			},
		},
		"using with a single nonexistent constant": {
			input: `
				using Foo::Bar

				module Foo; end

				var a: Bar = 9
				var c: 3 = Bar
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(11, 2, 11), P(18, 2, 18)), "undefined type or constant `Foo::Bar`"),
				diagnostic.NewFailure(L("<main>", P(53, 6, 12), P(55, 6, 14)), "undefined type `Foo::Bar`"),
				diagnostic.NewFailure(L("<main>", P(76, 7, 16), P(78, 7, 18)), "`Foo::Bar` cannot be used as a value in expressions"),
			},
		},
		"using with a single nonexistent constant and as": {
			input: `
				using Foo::Bar as B

				module Foo; end

				var a: B = 9
				var c: 3 = B
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(11, 2, 11), P(18, 2, 18)), "undefined type or constant `Foo::Bar`"),
				diagnostic.NewFailure(L("<main>", P(58, 6, 12), P(58, 6, 12)), "undefined type `Foo::Bar`"),
				diagnostic.NewFailure(L("<main>", P(79, 7, 16), P(79, 7, 16)), "`Foo::Bar` cannot be used as a value in expressions"),
			},
		},
		"using with a single class goes out of scope": {
			input: `
				module Baz
					using Foo::Bar

					var a: Bar = Bar()
					var b: Bar = 9
				end

				class Foo
				 	class Bar; end
				end

				Bar
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(79, 6, 19), P(79, 6, 19)), "type `9` cannot be assigned to type `Foo::Bar`"),
				diagnostic.NewFailure(L("<main>", P(138, 13, 5), P(140, 13, 7)), "undefined constant `Bar`"),
			},
		},

		"using with a single method": {
			input: `
				using Foo::bar

				module Foo
					def bar: Int then 3
				end

				var a: Int = bar()
				var b: 9 = bar()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(108, 9, 16), P(112, 9, 20)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"using with a single class method": {
			input: `
				using Foo::bar

				class Foo
					singleton
						def bar: Int then 3
					end
				end

				var a: Int = bar()
				var b: 9 = bar()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(132, 11, 16), P(136, 11, 20)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"using with a single class method after its declaration": {
			input: `
				class Foo
					singleton
						def bar: Int then 3
					end
				end

				using Foo::bar

				var a: Int = bar()
				var b: 9 = bar()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(132, 11, 16), P(136, 11, 20)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"using with a single method and as": {
			input: `
				using Foo::bar as b

				module Foo
					def bar: Int then 3
				end

				var a: Int = b()
				var b: 9 = b()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(111, 9, 16), P(113, 9, 18)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"using with a single method under a nonexistent namespace": {
			input: `
				using Foo::bar

				var a: Int = bar()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(11, 2, 11), P(13, 2, 13)), "undefined namespace `Foo`"),
				diagnostic.NewFailure(L("<main>", P(38, 4, 18), P(40, 4, 20)), "method `bar` is not defined on type `Std::Object`"),
			},
		},
		"using with a single nonexistent method": {
			input: `
				using Foo::bar
				module Foo; end

				var a: Int = bar()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(11, 2, 11), P(18, 2, 18)), "undefined method `Foo::bar`"),
			},
		},
		"using with a single nonexistent macro": {
			input: `
				using Foo::bar!
				module Foo; end

				var a: Int = bar!()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(11, 2, 11), P(19, 2, 19)), "undefined method `Foo::bar!`"),
			},
		},
		"using with a few methods": {
			input: `
				using Foo::{bar, baz}

				module Foo
					def bar: Int then 3
					def baz: Float then .3
				end

				var a: Int = bar()
				var b: 9 = bar()

				var c: Float = baz()
				var d: 9.2 = baz()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(143, 10, 16), P(147, 10, 20)), "type `Std::Int` cannot be assigned to type `9`"),
				diagnostic.NewFailure(L("<main>", P(192, 13, 18), P(196, 13, 22)), "type `Std::Float` cannot be assigned to type `9.2`"),
			},
		},
		"using with a few methods and as": {
			input: `
				using Foo::{bar as r, baz as z}

				module Foo
					def bar: Int then 3
					def baz: Float then .3
				end

				var a: Int = r()
				var b: 9 = r()

				var c: Float = z()
				var d: 9.2 = z()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(151, 10, 16), P(153, 10, 18)), "type `Std::Int` cannot be assigned to type `9`"),
				diagnostic.NewFailure(L("<main>", P(196, 13, 18), P(198, 13, 20)), "type `Std::Float` cannot be assigned to type `9.2`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
