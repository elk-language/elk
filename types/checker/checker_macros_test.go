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
		"can take an expression convertible value": {
			input: `
				quote
					a + unquote(5)
				end
			`,
		},
		"cannot appear outside of quote": {
			input: `
				unquote(5)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(5, 2, 5), P(14, 2, 14)), "unquote expressions cannot appear in this context"),
			},
		},
		"report an error when the argument is not expression convertible": {
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
		"cannot be nested": {
			input: `
				quote
					a + unquote(unqoute(5))
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(28, 3, 18), P(37, 3, 27)), "method `unqoute` is not defined on type `Std::Object`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
