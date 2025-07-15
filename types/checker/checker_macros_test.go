package checker

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

func TestMacroBoundary(t *testing.T) {
	tests := testTable{
		"has it's own scope": {
			input: `
				do macro
					a := 5
				end
				a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(38, 5, 5), P(38, 5, 5)), "undefined local `a`"),
			},
		},
		"cannot access variables from outer scopes": {
			input: `
				a := 5
				do macro
					a
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(30, 4, 6), P(30, 4, 6)), "undefined local `a`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestQuote(t *testing.T) {
	tests := testTable{
		"should return an expression node": {
			input: `
				n := quote
					a + 2
				end
				var b: 1 = n
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(50, 5, 16), P(50, 5, 16)), "type `Std::Elk::AST::ExpressionNode` cannot be assigned to type `1`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestUnquote(t *testing.T) {
	tests := testTable{
		"cannot appear outside of quote": {
			input: `
				unquote(5)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(5, 2, 5), P(14, 2, 14)), "unquote expressions cannot appear in this context"),
			},
		},
		"cannot be nested": {
			input: `
				quote
					a + unquote(unquote(5))
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(28, 3, 18), P(37, 3, 27)), "unquote expressions cannot appear in this context"),
			},
		},

		"can take an expression convertible value": {
			input: `
				quote
					a + unquote(5)
				end
			`,
		},
		"report an error when the argument to unquote is not expression node convertible": {
			input: `
				quote
					a + unquote(Time.now)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(36, 3, 26)), "type `Std::Time` does not implement interface `Std::Elk::AST::ExpressionNode::Convertible`:\n\n  - missing method `Std::Elk::AST::ExpressionNode::Convertible.:to_ast_expr_node` with signature: `def to_ast_expr_node(): Std::Elk::AST::ExpressionNode`"),
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(36, 3, 26)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::ExpressionNode::Convertible`"),
			},
		},
		"report an error when the argument to unquote_expr is not expression node convertible": {
			input: `
				quote
					a + unquote_expr(Time.now)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(41, 3, 31)), "type `Std::Time` does not implement interface `Std::Elk::AST::ExpressionNode::Convertible`:\n\n  - missing method `Std::Elk::AST::ExpressionNode::Convertible.:to_ast_expr_node` with signature: `def to_ast_expr_node(): Std::Elk::AST::ExpressionNode`"),
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(41, 3, 31)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::ExpressionNode::Convertible`"),
			},
		},
		"report an error when the argument to short unquote is not expression node convertible": {
			input: `
				quote
					a + !{Time.now}
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(30, 3, 20)), "type `Std::Time` does not implement interface `Std::Elk::AST::ExpressionNode::Convertible`:\n\n  - missing method `Std::Elk::AST::ExpressionNode::Convertible.:to_ast_expr_node` with signature: `def to_ast_expr_node(): Std::Elk::AST::ExpressionNode`"),
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(30, 3, 20)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::ExpressionNode::Convertible`"),
			},
		},

		"can take an identifier convertible value": {
			input: `
				quote
					var unquote(:foo): String
				end
			`,
		},
		"report an error when the argument to unquote is not identifier node convertible": {
			input: `
				quote
					var unquote(Time.now): String
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(36, 3, 26)), "type `Std::Time` does not implement interface `Std::Elk::AST::IdentifierNode::Convertible`:\n\n  - missing method `Std::Elk::AST::IdentifierNode::Convertible.:to_ast_ident_node` with signature: `def to_ast_ident_node(): Std::Elk::AST::IdentifierNode`"),
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(36, 3, 26)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::IdentifierNode::Convertible`"),
			},
		},
		"report an error when the argument to unquote_ident is not identifier node convertible": {
			input: `
				quote
					var unquote_ident(Time.now): String
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(42, 3, 32)), "type `Std::Time` does not implement interface `Std::Elk::AST::IdentifierNode::Convertible`:\n\n  - missing method `Std::Elk::AST::IdentifierNode::Convertible.:to_ast_ident_node` with signature: `def to_ast_ident_node(): Std::Elk::AST::IdentifierNode`"),
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(42, 3, 32)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::IdentifierNode::Convertible`"),
			},
		},
		"report an error when the argument to short unquote is not identifier node convertible": {
			input: `
				quote
					var !{Time.now}: String
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(30, 3, 20)), "type `Std::Time` does not implement interface `Std::Elk::AST::IdentifierNode::Convertible`:\n\n  - missing method `Std::Elk::AST::IdentifierNode::Convertible.:to_ast_ident_node` with signature: `def to_ast_ident_node(): Std::Elk::AST::IdentifierNode`"),
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(30, 3, 20)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::IdentifierNode::Convertible`"),
			},
		},

		"can take a pattern expression convertible value": {
			input: `
				quote
					var ^[unquote(1)] = "foo"
				end
			`,
		},
		"report an error when the argument to unquote is not pattern expression node convertible": {
			input: `
				quote
					var ^[unquote(Time.now)] = "foo"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 12), P(38, 3, 28)), "type `Std::Time` does not implement interface `Std::Elk::AST::PatternExpressionNode::Convertible`:\n\n  - missing method `Std::Elk::AST::PatternExpressionNode::Convertible.:to_ast_pattern_expr_node` with signature: `def to_ast_pattern_expr_node(): Std::Elk::AST::PatternExpressionNode`"),
				diagnostic.NewFailure(L("<main>", P(22, 3, 12), P(38, 3, 28)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::PatternExpressionNode::Convertible`"),
			},
		},
		"report an error when the argument to short unquote is not pattern expression node convertible": {
			input: `
				quote
					var ^[!{Time.now}] = "foo"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 12), P(32, 3, 22)), "type `Std::Time` does not implement interface `Std::Elk::AST::PatternExpressionNode::Convertible`:\n\n  - missing method `Std::Elk::AST::PatternExpressionNode::Convertible.:to_ast_pattern_expr_node` with signature: `def to_ast_pattern_expr_node(): Std::Elk::AST::PatternExpressionNode`"),
				diagnostic.NewFailure(L("<main>", P(22, 3, 12), P(32, 3, 22)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::PatternExpressionNode::Convertible`"),
			},
		},

		"can take a constant convertible value": {
			input: `
				quote
					const unquote(:Bar) = "foo"
				end
			`,
		},
		"report an error when the argument to unquote is not constant node convertible": {
			input: `
				quote
					const unquote(Time.now) = "foo"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 12), P(38, 3, 28)), "type `Std::Time` does not implement interface `Std::Elk::AST::ConstantNode::Convertible`:\n\n  - missing method `Std::Elk::AST::ConstantNode::Convertible.:to_ast_const_node` with signature: `def to_ast_const_node(): Std::Elk::AST::ConstantNode`"),
				diagnostic.NewFailure(L("<main>", P(22, 3, 12), P(38, 3, 28)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::ConstantNode::Convertible`"),
			},
		},
		"report an error when the argument to unquote_const is not constant node convertible": {
			input: `
				quote
					const unquote_const(Time.now) = "foo"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 12), P(44, 3, 34)), "type `Std::Time` does not implement interface `Std::Elk::AST::ConstantNode::Convertible`:\n\n  - missing method `Std::Elk::AST::ConstantNode::Convertible.:to_ast_const_node` with signature: `def to_ast_const_node(): Std::Elk::AST::ConstantNode`"),
				diagnostic.NewFailure(L("<main>", P(22, 3, 12), P(44, 3, 34)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::ConstantNode::Convertible`"),
			},
		},
		"report an error when the argument to short unquote is not constant node convertible": {
			input: `
				quote
					const !{Time.now} = "foo"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(22, 3, 12), P(32, 3, 22)), "type `Std::Time` does not implement interface `Std::Elk::AST::ConstantNode::Convertible`:\n\n  - missing method `Std::Elk::AST::ConstantNode::Convertible.:to_ast_const_node` with signature: `def to_ast_const_node(): Std::Elk::AST::ConstantNode`"),
				diagnostic.NewFailure(L("<main>", P(22, 3, 12), P(32, 3, 22)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::ConstantNode::Convertible`"),
			},
		},

		"can take a pattern convertible value": {
			input: `
				quote
					var [unquote(Elk::AST::ListPatternNode())] = "foo"
				end
			`,
		},
		"report an error when the argument to unquote is not pattern node convertible": {
			input: `
				quote
					var [unquote(Time.now)] = "foo"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(21, 3, 11), P(37, 3, 27)), "type `Std::Time` does not implement interface `Std::Elk::AST::PatternNode::Convertible`:\n\n  - missing method `Std::Elk::AST::PatternNode::Convertible.:to_ast_pattern_node` with signature: `def to_ast_pattern_node(): Std::Elk::AST::PatternNode`"),
				diagnostic.NewFailure(L("<main>", P(21, 3, 11), P(37, 3, 27)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::PatternNode::Convertible`"),
			},
		},
		"report an error when the argument to unquote_pattern is not pattern node convertible": {
			input: `
				quote
					var [unquote_pattern(Time.now)] = "foo"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(21, 3, 11), P(45, 3, 35)), "type `Std::Time` does not implement interface `Std::Elk::AST::PatternNode::Convertible`:\n\n  - missing method `Std::Elk::AST::PatternNode::Convertible.:to_ast_pattern_node` with signature: `def to_ast_pattern_node(): Std::Elk::AST::PatternNode`"),
				diagnostic.NewFailure(L("<main>", P(21, 3, 11), P(45, 3, 35)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::PatternNode::Convertible`"),
			},
		},
		"report an error when the argument to short unquote is not pattern node convertible": {
			input: `
				quote
					var [!{Time.now}] = "foo"
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(21, 3, 11), P(31, 3, 21)), "type `Std::Time` does not implement interface `Std::Elk::AST::PatternNode::Convertible`:\n\n  - missing method `Std::Elk::AST::PatternNode::Convertible.:to_ast_pattern_node` with signature: `def to_ast_pattern_node(): Std::Elk::AST::PatternNode`"),
				diagnostic.NewFailure(L("<main>", P(21, 3, 11), P(31, 3, 21)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::PatternNode::Convertible`"),
			},
		},

		"can take a instance variable node convertible value": {
			input: `
				quote
					var unquote_ivar(:foo): String?
				end
			`,
		},
		"report an error when the argument to unquote_ivar is not pattern node convertible": {
			input: `
				quote
					var unquote_ivar(Time.now): String?
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(41, 3, 31)), "type `Std::Time` does not implement interface `Std::Elk::AST::InstanceVariableNode::Convertible`:\n\n  - missing method `Std::Elk::AST::InstanceVariableNode::Convertible.:to_ast_ivar_node` with signature: `def to_ast_ivar_node(): Std::Elk::AST::ExpressionNode`"),
				diagnostic.NewFailure(L("<main>", P(20, 3, 10), P(41, 3, 31)), "type `Std::Time` cannot be assigned to type `Std::Elk::AST::InstanceVariableNode::Convertible`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestExpandMacro(t *testing.T) {
	tests := testTable{
		"define a class in a top level macro": {
			input: `
				using Std::Elk::AST::*

				macro klass(name: ConstantNode)
					quote
						class !{name}
							def foo: String
								"lol"
							end
						end
					end
				end
				klass!(Bar)

				b := Bar()
				var a: String = b.foo
			`,
		},
		"inherit from a generated class": {
			input: `
				using Std::Elk::AST::*

				macro klass(name: ConstantNode)
					quote
						class !{name}
							def foo: String
								"lol"
							end
						end
					end
				end
				klass!(Bar)

				class Baz < Bar; end
				b := Baz()
				var a: String = b.foo
			`,
		},
		"define a class in an expression level macro": {
			input: `
				using Std::Elk::AST::*

				macro klass(name: ConstantNode)
					quote
						class !{name}
							def foo: String
								"lol"
							end
						end
					end
				end
				klass!(Bar) + 5

				b := Bar()
				var a: String = b.foo
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(175, 13, 5), P(185, 13, 15)), "class definitions cannot appear in this context"),
				diagnostic.NewFailure(L("<main>", P(175, 13, 5), P(185, 13, 15)), "method definitions cannot appear in this context"),
				diagnostic.NewFailure(L("<main>", P(175, 13, 5), P(189, 13, 19)), "method `+` is not defined on type `Std::Nil`"),
				diagnostic.NewFailure(L("<main>", P(201, 15, 10), P(203, 15, 12)), "undefined type `Bar`"),
			},
		},
		"generate a method": {
			input: `
				using Std::Elk::AST::*

				macro fn(name: IdentifierNode)
					quote
						def !{name}: String
							!{name.value}
						end
					end
				end

				module Foo
					fn!(lol)
				end

				var a: String = Foo.lol
			`,
		},
		"use a generated method in a constant": {
			input: `
				using Std::Elk::AST::*

				macro fn(name: IdentifierNode)
					quote
						def !{name}: String
							!{name.value}
						end
					end
				end

				module Foo
					fn!(lol)
					const BAR: String = "#{Foo.lol()} hey"
				end

				var a: String = Foo::BAR
			`,
		},
		"generate a method in an expression level macro": {
			input: `
				using Std::Elk::AST::*

				macro fn(name: IdentifierNode)
					quote
						def !{name}: String
							!{name.value}
						end
					end
				end

				module Foo
					fn!(lol) * 5
				end

				var a: String = Foo.lol
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(170, 13, 6), P(177, 13, 13)), "method definitions cannot appear in this context"),
				diagnostic.NewFailure(L("<main>", P(216, 16, 25), P(218, 16, 27)), "method `lol` is not defined on type `Foo`"),
			},
		},
		"generate an expression": {
			input: `
				using Std::Elk::AST::*

				macro fib(i: IntLiteralNode)
					calc_fib := |n: Int|: Int ->
						return 1 if n < 3

						calc_fib(n - 2) + calc_fib(n - 1)
					end

					calc_fib(i.to_int).to_ast_node
				end

				timeout := fib!(15) + 2
			`,
		},
		"use a scoped macro": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro fib(i: IntLiteralNode)
						calc_fib := |n: Int|: Int ->
							return 1 if n < 3

							calc_fib(n - 2) + calc_fib(n - 1)
						end

						calc_fib(i.to_int).to_ast_node
					end
				end

				timeout := Foo::fib!(15) + 2
			`,
		},
		"use an invalid namespace for a macro": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro fib(i: IntLiteralNode)
						calc_fib := |n: Int|: Int ->
							return 1 if n < 3

							calc_fib(n - 2) + calc_fib(n - 1)
						end

						calc_fib(i.to_int).to_ast_node
					end
				end

				timeout := Foob::fib!(15) + 2
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(261, 16, 16), P(264, 16, 19)), "undefined constant `Foob`"),
			},
		},
		"call a nonexistent macro": {
			input: `
				module Foo; end
				timeout := Foo::fib!(15) + 2
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(41, 3, 21), P(43, 3, 23)), "undefined macro `fib!`"),
			},
		},
		"inherit a macro": {
			input: `
				using Std::Elk::AST::*

				class Foo
					macro fib(i: IntLiteralNode)
						calc_fib := |n: Int|: Int ->
							return 1 if n < 3

							calc_fib(n - 2) + calc_fib(n - 1)
						end

						calc_fib(i.to_int).to_ast_node
					end
				end

				class Bar < Foo; end

				timeout := Bar::fib!(15) + 2
			`,
		},
		"throw an error in a macro": {
			input: `
				using Std::Elk::AST::*

				macro foo then throw unchecked 5
				foo!
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L("<main>", P(70, 5, 5), P(72, 5, 7)),
					"error while executing macro `Std::Kernel::foo!`: 5\nStack trace (the most recent call is last)\n 0: <main>:4, in `foo!`\n",
				),
			},
		},
		"cannot make nil-safe macro calls": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(a: IntLiteralNode) then a
				end
				Foo?.baz!(5)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(100, 7, 8), P(101, 7, 9)), "invalid macro call operator"),
			},
		},
		"cannot make cascade macro calls": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(a: IntLiteralNode) then a
				end
				Foo..baz!(5)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(100, 7, 8), P(101, 7, 9)), "invalid macro call operator"),
			},
		},
		"cannot make nil-safe cascade macro calls": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(a: IntLiteralNode) then a
				end
				Foo?..baz!(5)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(100, 7, 8), P(102, 7, 10)), "invalid macro call operator"),
			},
		},

		"missing required argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode) then bar
					baz!("foo")
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(116, 6, 6), P(126, 6, 16)), "argument `c` is missing in call to `Foo::baz!`"),
			},
		},
		"all required positional arguments": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode) then c
					baz!("foo", 5)
				end
			`,
		},
		"all required positional arguments with wrong type": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode) then c
					baz!(123.4, 5)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(119, 6, 11), P(123, 6, 15)), "expected type `Std::Elk::AST::StringLiteralNode` for parameter `bar` in call to `Foo::baz!`, got type `Std::Elk::AST::FloatLiteralNode`"),
			},
		},
		"too many positional arguments": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode) then c
					baz!("foo", 5, 28, 9, 0)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(114, 6, 6), P(137, 6, 29)), "expected 2 arguments in call to `Foo::baz!`, got 5"),
			},
		},
		"missing required argument with named argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode) then c
					baz!(bar: "foo")
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(114, 6, 6), P(129, 6, 21)), "argument `c` is missing in call to `Foo::baz!`"),
			},
		},
		"all required named arguments": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode) then c
					baz!(c: 5, bar: "foo")
				end
			`,
		},
		"all required named arguments with wrong type": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode) then c
					baz!(c: 5, bar: 123.4)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(125, 6, 17), P(134, 6, 26)), "expected type `Std::Elk::AST::StringLiteralNode` for parameter `bar` in call to `Foo::baz!`, got type `Std::Elk::AST::FloatLiteralNode`"),
			},
		},
		"duplicated positional argument as named argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode) then c
					baz!("foo", 5, bar: 9)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(129, 6, 21), P(134, 6, 26)), "duplicated argument `bar` in call to `Foo::baz!`"),
				diagnostic.NewFailure(L("<main>", P(129, 6, 21), P(134, 6, 26)), "expected type `Std::Elk::AST::StringLiteralNode` for parameter `bar` in call to `Foo::baz!`, got type `Std::Elk::AST::IntLiteralNode`"),
			},
		},
		"duplicated named argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode) then c
					baz!("foo", 2, c: 3, c: 9)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(129, 6, 21), P(132, 6, 24)), "duplicated argument `c` in call to `Foo::baz!`"),
				diagnostic.NewFailure(L("<main>", P(135, 6, 27), P(138, 6, 30)), "duplicated argument `c` in call to `Foo::baz!`"),
			},
		},
		"call with missing optional argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode = 3.to_ast_node) then c
					var a: 3 = baz!("foo")
				end
			`,
		},
		"call with optional argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode = 3.to_ast_node) then c
					var a: 9 = baz!("foo", 9)
				end
			`,
		},
		"call with missing rest arguments": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(*b: FloatLiteralNode) then 3.to_ast_node
					baz!
				end
			`,
		},
		"call with rest arguments": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(*b: FloatLiteralNode) then 3.to_ast_node
					baz! 1.2, 56.9, .5
				end
			`,
		},
		"call with splat argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(*b: FloatLiteralNode) then 1.to_ast_node
					arr := [1.2, 2.6, 3.1]
					baz!(*arr)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(138, 7, 11), P(141, 7, 14)), "expected type `Std::Elk::AST::FloatLiteralNode` for rest parameter `*b` in call to `Foo::baz!`, got type `Std::Elk::AST::SplatExpressionNode`"),
			},
		},
		"call with rest argument given by name": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(*b: FloatLiteralNode) then 3.to_ast_node
					baz! b: []
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(110, 6, 11), P(114, 6, 15)), "nonexistent parameter `b` given in call to `Foo::baz!`"),
			},
		},
		"call with required post arguments": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, *b: FloatLiteralNode, c: IntLiteralNode) then c
					baz!("foo", 3)
				end
			`,
		},
		"call with missing post argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, *b: FloatLiteralNode, c: IntLiteralNode) then 3.to_ast_node
					baz!("foo")
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(148, 6, 6), P(158, 6, 16)), "argument `c` is missing in call to `Foo::baz!`"),
			},
		},
		"call with rest and post arguments": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, *b: FloatLiteralNode, c: IntLiteralNode) then 3.to_ast_node
					baz!("foo", 2.5, .9, 128.1, 3)
				end
			`,
		},
		"call with rest and post arguments and wrong type in post": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, *b: FloatLiteralNode, c: IntLiteralNode) then c
					baz!("foo", 2.5, .9, 128.1, 3.2)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(164, 6, 34), P(166, 6, 36)), "expected type `Std::Elk::AST::IntLiteralNode` for parameter `c` in call to `Foo::baz!`, got type `Std::Elk::AST::FloatLiteralNode`"),
			},
		},
		"call with rest and post arguments and wrong type in rest": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, *b: FloatLiteralNode, c: IntLiteralNode) then c
					baz!("foo", 212, .9, '282', 3)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(148, 6, 18), P(150, 6, 20)), "expected type `Std::Elk::AST::FloatLiteralNode` for rest parameter `*b` in call to `Foo::baz!`, got type `Std::Elk::AST::IntLiteralNode`"),
				diagnostic.NewFailure(L("<main>", P(157, 6, 27), P(161, 6, 31)), "expected type `Std::Elk::AST::FloatLiteralNode` for rest parameter `*b` in call to `Foo::baz!`, got type `Std::Elk::AST::RawStringLiteralNode`"),
			},
		},
		"call with rest arguments and missing post argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, *b: FloatLiteralNode, c: IntLiteralNode) then c
					baz!("foo", 2.5, .9, 128.1)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(157, 6, 27), P(161, 6, 31)), "expected type `Std::Elk::AST::IntLiteralNode` for parameter `c` in call to `Foo::baz!`, got type `Std::Elk::AST::FloatLiteralNode`"),
			},
		},
		"call with named post argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, *b: FloatLiteralNode, c: IntLiteralNode) then c
					baz!("foo", c: 3)
				end
			`,
		},
		"call with named pre rest argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, *b: FloatLiteralNode, c: IntLiteralNode) then c
					baz!(bar: "foo", c: 3)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(136, 6, 6), P(157, 6, 27)), "expected 1... positional arguments in call to `Foo::baz!`, got 0"),
			},
		},
		"call without named rest arguments": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode, **rest: IntLiteralNode) then c
					baz!("foo", 5)
				end
			`,
		},
		"call with named rest arguments": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode, **rest: IntLiteralNode) then c
					baz!("foo", d: 25, c: 5, e: 11)
				end
			`,
		},
		"call with named rest arguments with wrong type": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(bar: StringLiteralNode, c: IntLiteralNode, **rest: IntLiteralNode) then c
					baz!("foo", d: .2, c: 5, e: .1)
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(150, 6, 18), P(154, 6, 22)), "expected type `Std::Elk::AST::IntLiteralNode` for named rest parameter `**rest` in call to `baz!`, got type `Std::Elk::AST::FloatLiteralNode`"),
				diagnostic.NewFailure(L("<main>", P(163, 6, 31), P(167, 6, 35)), "expected type `Std::Elk::AST::IntLiteralNode` for named rest parameter `**rest` in call to `baz!`, got type `Std::Elk::AST::FloatLiteralNode`"),
			},
		},
		"call with double splat argument": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(**b: FloatLiteralNode) then 3.to_ast_node
					map := { foo: 1.2, bar: 29.9 }
					baz! a: 1.2, **map, b: .5
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(155, 7, 19), P(159, 7, 23)), "double splat arguments cannot be used in macro call `Foo::baz!`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestMacroDefinition(t *testing.T) {
	tests := testTable{
		"declare within a macro": {
			input: `
				macro foo
					macro bar; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 3, 6), P(33, 3, 19)), "macro definitions cannot appear in this context"),
			},
		},
		"param types must inherit from ExpressionNode": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(b: Float)
						loop; end
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(59, 5, 16), P(66, 5, 23)), "type `Std::Float` does not inherit from `Std::Elk::AST::ExpressionNode`, macro parameters must be expression nodes"),
			},
		},
		"returned value must inherit from ExpressionNode": {
			input: `
				module Foo
					macro baz
						5
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(37, 4, 7), P(37, 4, 7)), "type `5` cannot be assigned to type `Std::Elk::AST::ExpressionNode`"),
			},
		},
		"positional rest params have tuple types": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(*b: FloatLiteralNode)
						var c: nil = b
						loop; end
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(100, 6, 20), P(100, 6, 20)), "type `Std::Tuple[Std::Elk::AST::FloatLiteralNode]` cannot be assigned to type `nil`"),
			},
		},
		"named rest params have record types": {
			input: `
				using Std::Elk::AST::*

				module Foo
					macro baz(**b: FloatLiteralNode)
						var c: nil = b
						loop; end
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(101, 6, 20), P(101, 6, 20)), "type `Std::Record[Std::Symbol, Std::Elk::AST::FloatLiteralNode]` cannot be assigned to type `nil`"),
			},
		},
		"cannot declare with type parameters": {
			input: `
				macro foo[V](a: V)
					a
					loop; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(14, 2, 14), P(14, 2, 14)), "unexpected [, expected a statement separator `\\n`, `;`"),
				diagnostic.NewFailure(L("<main>", P(17, 2, 17), P(17, 2, 17)), "unexpected (, expected a statement separator `\\n`, `;`"),
			},
		},
		"cannot declare with a return type": {
			input: `
				macro foo: String
					5
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(14, 2, 14), P(14, 2, 14)), "unexpected :, expected a statement separator `\\n`, `;`"),
			},
		},
		"redeclare the macro in the same class": {
			input: `
				using Std::Elk::AST::*

				class Foo
					macro baz(a: IntLiteralNode) then a
					macro baz then try IntLiteralNode("123")

					baz!
				end
			`,
		},
		"macros get hoisted to the top": {
			input: `
				using Std::Elk::AST::*

			  baz!
				macro baz then try IntLiteralNode("123")
			`,
		},
		"macros cannot reference other macros": {
			input: `
				using Std::Elk::AST::*

				macro bar then try IntLiteralNode("1")
				macro foo then bar!()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(91, 5, 20), P(93, 5, 22)), "macros cannot be used in macro definitions"),
			},
		},
		"macros cannot use instance variable parameters": {
			input: `
				using Std::Elk::AST::*

				macro bar(@a: IntLiteralNode) then a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(43, 4, 15), P(44, 4, 16)), "unexpected INSTANCE_VARIABLE, expected an identifier as the name of the declared signature parameter"),
				diagnostic.NewFailure(L("<main>", P(45, 4, 17), P(45, 4, 17)), "unexpected :, expected )"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
