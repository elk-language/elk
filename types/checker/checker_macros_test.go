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
					var [unquote(:bar)] = "foo"
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
