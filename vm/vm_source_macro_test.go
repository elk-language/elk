package vm_test

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

func TestVMSource_Quote(t *testing.T) {
	tests := sourceTestTable{
		"without unquote": {
			source: `
				quote
					1 + 2
				end
			`,
			wantStackTop: value.Ref(
				ast.NewBinaryExpressionNode(
					L(P(16, 3, 6), P(20, 3, 10)),
					T(L(P(18, 3, 8), P(18, 3, 8)), token.PLUS),
					ast.NewIntLiteralNode(L(P(16, 3, 6), P(16, 3, 6)), "1"),
					ast.NewIntLiteralNode(L(P(20, 3, 10), P(20, 3, 10)), "2"),
				),
			),
		},
		"unquote expression": {
			source: `
				a := 5
				quote
					1 + unquote(a)
				end
			`,
			wantStackTop: value.Ref(
				ast.NewBinaryExpressionNode(
					L(P(27, 4, 6), P(40, 4, 19)),
					T(L(P(29, 4, 8), P(29, 4, 8)), token.PLUS),
					ast.NewIntLiteralNode(L(P(27, 4, 6), P(27, 4, 6)), "1"),
					ast.NewIntLiteralNode(position.ZeroLocation, "5"),
				),
			),
		},
		"unquote several items": {
			source: `
				a := :foo
				b := 5
				c := 2.5

				quote
					unquote_ident(a) := 1 + !{b} - !{c}
				end
			`,
			wantStackTop: value.Ref(
				ast.NewAssignmentExpressionNode(
					L(P(55, 7, 6), P(89, 7, 40)),
					T(L(P(72, 7, 23), P(73, 7, 24)), token.COLON_EQUAL),
					ast.NewPublicIdentifierNode(position.ZeroLocation, "foo"),
					ast.NewBinaryExpressionNode(
						L(P(75, 7, 26), P(89, 7, 40)),
						T(L(P(84, 7, 35), P(84, 7, 35)), token.MINUS),
						ast.NewBinaryExpressionNode(
							L(P(75, 7, 26), P(82, 7, 33)),
							T(L(P(77, 7, 28), P(77, 7, 28)), token.PLUS),
							ast.NewIntLiteralNode(L(P(75, 7, 26), P(75, 7, 26)), "1"),
							ast.NewIntLiteralNode(position.ZeroLocation, "5"),
						),
						ast.NewFloatLiteralNode(position.ZeroLocation, "2.5"),
					),
				),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_ExpandMacro(t *testing.T) {
	tests := sourceTestTable{
		"compile-time fibonacci": {
			source: `
				using Std::Elk::AST::*

				macro fib(i: IntLiteralNode)
					calc_fib := |n: Int|: Int ->
						return 1 if n < 3

						calc_fib(n - 2) + calc_fib(n - 1)
					end

					calc_fib(i.to_int).to_ast_node
				end

				fib!(10) * 2
			`,
			wantStackTop: value.SmallInt(110).ToValue(),
		},
		"recursive fibonacci macro": {
			source: `
				using Std::Elk::AST::*

				macro fib(i: IntLiteralNode)
					int := i.to_int
					return try IntLiteralNode('1') if int < 3

					quote
						fib!(!{int - 1}) + fib!(!{int - 2})
					end
				end

				fib!(10) * 2
			`,
			wantStackTop: value.SmallInt(110).ToValue(),
		},
		"call a scoped macro": {
			source: `
				using Std::Elk::AST::*

				module Math
					macro fib(i: IntLiteralNode)
						calc_fib := |n: Int|: Int ->
							return 1 if n < 3

							calc_fib(n - 2) + calc_fib(n - 1)
						end

						calc_fib(i.to_int).to_ast_node
					end
				end

				Math::fib!(10) * 2
			`,
			wantStackTop: value.SmallInt(110).ToValue(),
		},
		"define a class and call a getter": {
			source: `
				using Std::Elk::AST::*

				macro box(name: ConstantNode, typ: TypeExpressionNode)
					quote
						class !{name}
							attr value: !{typ.type_node}
							init(@value: !{typ.type_node}); end
						end
					end
				end

				box!(BoxString, type String?)

				b := BoxString("foo")
				b.value
			`,
			wantStackTop: value.Ref(value.String("foo")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
