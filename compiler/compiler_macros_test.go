package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestQuote(t *testing.T) {
	tests := testTable{
		"without unquote": {
			input: `
				quote
					1 + 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewMacroBoundaryNode(
							L(P(5, 2, 5), P(28, 4, 7)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(P(16, 3, 6), P(21, 3, 11)),
									ast.NewBinaryExpressionNode(
										L(P(16, 3, 6), P(20, 3, 10)),
										T(L(P(18, 3, 8), P(18, 3, 8)), token.PLUS),
										ast.NewIntLiteralNode(L(P(16, 3, 6), P(16, 3, 6)), "1"),
										ast.NewIntLiteralNode(L(P(20, 3, 10), P(20, 3, 10)), "2"),
									),
								),
							},
							"",
						),
					),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 3)),
				},
			),
		},
		"with unquote": {
			input: `
				quote
					1 + unquote(2)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.INT_2),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(38, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewMacroBoundaryNode(
							L(P(5, 2, 5), P(37, 4, 7)),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(P(16, 3, 6), P(30, 3, 20)),
									ast.NewBinaryExpressionNode(
										L(P(16, 3, 6), P(29, 3, 19)),
										T(L(P(18, 3, 8), P(18, 3, 8)), token.PLUS),
										ast.NewIntLiteralNode(L(P(16, 3, 6), P(16, 3, 6)), "1"),
										ast.NewUnquoteExpressionNode(
											L(P(20, 3, 10), P(29, 3, 19)),
											ast.NewIntLiteralNode(L(P(28, 3, 18), P(28, 3, 18)), "2"),
										),
									),
								),
							},
							"",
						),
					),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_expr_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 3)),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}
