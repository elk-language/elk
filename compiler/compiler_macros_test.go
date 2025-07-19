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
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewBinaryExpressionNode(
							L(P(16, 3, 6), P(20, 3, 10)),
							T(L(P(18, 3, 8), P(18, 3, 8)), token.PLUS),
							ast.NewIntLiteralNode(L(P(16, 3, 6), P(16, 3, 6)), "1"),
							ast.NewIntLiteralNode(L(P(20, 3, 10), P(20, 3, 10)), "2"),
						),
					),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},

		"unquote expression": {
			input: `
				quote
					1 + unquote(5)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.INT_5),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(38, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewBinaryExpressionNode(
							L(P(16, 3, 6), P(29, 3, 19)),
							T(L(P(18, 3, 8), P(18, 3, 8)), token.PLUS),
							ast.NewIntLiteralNode(L(P(16, 3, 6), P(16, 3, 6)), "1"),
							ast.NewUnquoteNode(
								L(P(20, 3, 10), P(29, 3, 19)),
								ast.UNQUOTE_EXPRESSION_KIND,
								ast.NewIntLiteralNode(L(P(28, 3, 18), P(28, 3, 18)), "5"),
							),
						),
					),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_expr_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},
		"short unquote expression": {
			input: `
				quote
					1 + !{5}
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.INT_5),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewBinaryExpressionNode(
							L(P(16, 3, 6), P(23, 3, 13)),
							T(L(P(18, 3, 8), P(18, 3, 8)), token.PLUS),
							ast.NewIntLiteralNode(L(P(16, 3, 6), P(16, 3, 6)), "1"),
							ast.NewUnquoteNode(
								L(P(20, 3, 10), P(23, 3, 13)),
								ast.UNQUOTE_EXPRESSION_KIND,
								ast.NewIntLiteralNode(L(P(22, 3, 12), P(22, 3, 12)), "5"),
							),
						),
					),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_expr_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},

		"unquote identifier": {
			input: `
				quote
					var unquote(:foo): String
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(49, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewVariableDeclarationNode(
							L(P(16, 3, 6), P(40, 3, 30)),
							"",
							ast.NewUnquoteNode(
								L(P(20, 3, 10), P(32, 3, 22)),
								ast.UNQUOTE_IDENTIFIER_KIND,
								ast.NewSimpleSymbolLiteralNode(L(P(28, 3, 18), P(31, 3, 21)), "foo"),
							),
							ast.NewPublicConstantNode(L(P(35, 3, 25), P(40, 3, 30)), "String"),
							nil,
						),
					),
					value.ToSymbol("foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_ident_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},
		"short unquote identifier": {
			input: `
				quote
					var !{:foo}: String
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewVariableDeclarationNode(
							L(P(16, 3, 6), P(34, 3, 24)),
							"",
							ast.NewUnquoteNode(
								L(P(20, 3, 10), P(26, 3, 16)),
								ast.UNQUOTE_IDENTIFIER_KIND,
								ast.NewSimpleSymbolLiteralNode(L(P(22, 3, 12), P(25, 3, 15)), "foo"),
							),
							ast.NewPublicConstantNode(L(P(29, 3, 19), P(34, 3, 24)), "String"),
							nil,
						),
					),
					value.ToSymbol("foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_ident_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},
		"unquote_ident": {
			input: `
				quote
					var unquote_ident(:foo): String
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(55, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewVariableDeclarationNode(
							L(P(16, 3, 6), P(46, 3, 36)),
							"",
							ast.NewUnquoteNode(
								L(P(20, 3, 10), P(38, 3, 28)),
								ast.UNQUOTE_IDENTIFIER_KIND,
								ast.NewSimpleSymbolLiteralNode(L(P(34, 3, 24), P(37, 3, 27)), "foo"),
							),
							ast.NewPublicConstantNode(L(P(41, 3, 31), P(46, 3, 36)), "String"),
							nil,
						),
					),
					value.ToSymbol("foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_ident_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},

		"unquote pattern expression": {
			input: `
				quote
					var ^[unquote(1)] = "foo"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.INT_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(49, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewVariablePatternDeclarationNode(
							L(P(16, 3, 6), P(40, 3, 30)),
							ast.NewSetPatternNode(
								L(P(20, 3, 10), P(32, 3, 22)),
								[]ast.PatternNode{
									ast.NewUnquoteNode(
										L(P(22, 3, 12), P(31, 3, 21)),
										ast.UNQUOTE_PATTERN_EXPRESSION_KIND,
										ast.NewIntLiteralNode(L(P(30, 3, 20), P(30, 3, 20)), "1"),
									),
								},
							),
							ast.NewDoubleQuotedStringLiteralNode(L(P(36, 3, 26), P(40, 3, 30)), "foo"),
						),
					),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_pattern_expr_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},
		"short unquote pattern expression": {
			input: `
				quote
					var ^[!{1}] = "foo"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.INT_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewVariablePatternDeclarationNode(
							L(P(16, 3, 6), P(34, 3, 24)),
							ast.NewSetPatternNode(
								L(P(20, 3, 10), P(26, 3, 16)),
								[]ast.PatternNode{
									ast.NewUnquoteNode(
										L(P(22, 3, 12), P(25, 3, 15)),
										ast.UNQUOTE_PATTERN_EXPRESSION_KIND,
										ast.NewIntLiteralNode(L(P(24, 3, 14), P(24, 3, 14)), "1"),
									),
								},
							),
							ast.NewDoubleQuotedStringLiteralNode(L(P(30, 3, 20), P(34, 3, 24)), "foo"),
						),
					),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_pattern_expr_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},

		"unquote pattern": {
			input: `
				quote
					var [unquote(Elk::AST::ListPatternNode())] = "foo"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.INSTANTIATE8), 2,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(74, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewVariablePatternDeclarationNode(
							L(P(16, 3, 6), P(65, 3, 55)),
							ast.NewListPatternNode(
								L(P(20, 3, 10), P(57, 3, 47)),
								[]ast.PatternNode{
									ast.NewUnquoteNode(
										L(P(21, 3, 11), P(56, 3, 46)),
										ast.UNQUOTE_PATTERN_KIND,
										ast.NewConstructorCallNode(
											L(P(29, 3, 19), P(55, 3, 45)),
											ast.NewPublicConstantNode(
												L(P(29, 3, 19), P(53, 3, 43)),
												"Std::Elk::AST::ListPatternNode",
											),
											[]ast.ExpressionNode{
												ast.NewUndefinedLiteralNode(L(P(29, 3, 19), P(53, 3, 43))),
												ast.NewUndefinedLiteralNode(L(P(29, 3, 19), P(53, 3, 43))),
											},
											nil,
										),
									),
								},
							),
							ast.NewDoubleQuotedStringLiteralNode(L(P(61, 3, 51), P(65, 3, 55)), "foo"),
						),
					),
					value.ToSymbol("Std::Elk::AST::ListPatternNode").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_pattern_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},
		"unquote_pattern": {
			input: `
				quote
					var [unquote_pattern(Elk::AST::ListPatternNode())] = "foo"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.INSTANTIATE8), 2,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(82, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewVariablePatternDeclarationNode(
							L(P(16, 3, 6), P(73, 3, 63)),
							ast.NewListPatternNode(
								L(P(20, 3, 10), P(65, 3, 55)),
								[]ast.PatternNode{
									ast.NewUnquoteNode(
										L(P(21, 3, 11), P(64, 3, 54)),
										ast.UNQUOTE_PATTERN_KIND,
										ast.NewConstructorCallNode(
											L(P(37, 3, 27), P(63, 3, 53)),
											ast.NewPublicConstantNode(
												L(P(37, 3, 27), P(61, 3, 51)),
												"Std::Elk::AST::ListPatternNode",
											),
											[]ast.ExpressionNode{
												ast.NewUndefinedLiteralNode(L(P(37, 3, 27), P(61, 3, 51))),
												ast.NewUndefinedLiteralNode(L(P(37, 3, 27), P(61, 3, 51))),
											},
											nil,
										),
									),
								},
							),
							ast.NewDoubleQuotedStringLiteralNode(L(P(69, 3, 59), P(73, 3, 63)), "foo"),
						),
					),
					value.ToSymbol("Std::Elk::AST::ListPatternNode").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_pattern_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},
		"short unquote pattern": {
			input: `
				quote
					var [!{Elk::AST::ListPatternNode()}] = "foo"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.INSTANTIATE8), 2,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewVariablePatternDeclarationNode(
							L(P(16, 3, 6), P(59, 3, 49)),
							ast.NewListPatternNode(
								L(P(20, 3, 10), P(51, 3, 41)),
								[]ast.PatternNode{
									ast.NewUnquoteNode(
										L(P(21, 3, 11), P(50, 3, 40)),
										ast.UNQUOTE_PATTERN_KIND,
										ast.NewConstructorCallNode(
											L(P(23, 3, 13), P(49, 3, 39)),
											ast.NewPublicConstantNode(
												L(P(23, 3, 13), P(47, 3, 37)),
												"Std::Elk::AST::ListPatternNode",
											),
											[]ast.ExpressionNode{
												ast.NewUndefinedLiteralNode(L(P(23, 3, 13), P(47, 3, 37))),
												ast.NewUndefinedLiteralNode(L(P(23, 3, 13), P(47, 3, 37))),
											},
											nil,
										),
									),
								},
							),
							ast.NewDoubleQuotedStringLiteralNode(L(P(55, 3, 45), P(59, 3, 49)), "foo"),
						),
					),
					value.ToSymbol("Std::Elk::AST::ListPatternNode").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_pattern_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},

		"unquote constant": {
			input: `
				quote
					const unquote(:Bar) = "foo"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(51, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewConstantDeclarationNode(
							L(P(16, 3, 6), P(42, 3, 32)),
							"",
							ast.NewUnquoteNode(
								L(P(22, 3, 12), P(34, 3, 24)),
								ast.UNQUOTE_CONSTANT_KIND,
								ast.NewSimpleSymbolLiteralNode(L(P(30, 3, 20), P(33, 3, 23)), "Bar"),
							),
							nil,
							ast.NewDoubleQuotedStringLiteralNode(L(P(38, 3, 28), P(42, 3, 32)), "foo"),
						),
					),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_const_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},
		"unquote_const": {
			input: `
				quote
					const unquote_const(:Bar) = "foo"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(57, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewConstantDeclarationNode(
							L(P(16, 3, 6), P(48, 3, 38)),
							"",
							ast.NewUnquoteNode(
								L(P(22, 3, 12), P(40, 3, 30)),
								ast.UNQUOTE_CONSTANT_KIND,
								ast.NewSimpleSymbolLiteralNode(L(P(36, 3, 26), P(39, 3, 29)), "Bar"),
							),
							nil,
							ast.NewDoubleQuotedStringLiteralNode(L(P(44, 3, 34), P(48, 3, 38)), "foo"),
						),
					),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_const_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},
		"short unquote constant": {
			input: `
				quote
					const !{:Bar} = "foo"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewConstantDeclarationNode(
							L(P(16, 3, 6), P(36, 3, 26)),
							"",
							ast.NewUnquoteNode(
								L(P(22, 3, 12), P(28, 3, 18)),
								ast.UNQUOTE_CONSTANT_KIND,
								ast.NewSimpleSymbolLiteralNode(L(P(24, 3, 14), P(27, 3, 17)), "Bar"),
							),
							nil,
							ast.NewDoubleQuotedStringLiteralNode(L(P(32, 3, 22), P(36, 3, 26)), "foo"),
						),
					),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_const_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
				},
			),
		},

		"unquote_ivar": {
			input: `
				quote
					var unquote_ivar(:foo): String?
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(55, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(
						ast.NewInstanceVariableDeclarationNode(
							L(P(16, 3, 6), P(46, 3, 36)),
							"",
							ast.NewUnquoteNode(
								L(P(20, 3, 10), P(37, 3, 27)),
								ast.UNQUOTE_INSTANCE_VARIABLE_KIND,
								ast.NewSimpleSymbolLiteralNode(L(P(33, 3, 23), P(36, 3, 26)), "foo"),
							),
							ast.NewNilableTypeNode(
								L(P(40, 3, 30), P(46, 3, 36)),
								ast.NewPublicConstantNode(L(P(40, 3, 30), P(45, 3, 35)), "String"),
							),
						),
					),
					value.ToSymbol("foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("to_ast_ivar_node"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#splice"), 2)),
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

func TestMacroExpansion(t *testing.T) {
	tests := testTable{
		"compile-time fibonacci": {
			input: `
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 110,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(232, 14, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(14, 3),
				},
				[]value.Value{
					value.Undefined,
				},
			),
		},
		"call a scoped macro": {
			input: `
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_INT_8), 110,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(269, 16, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(16, 3),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(269, 16, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(16, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Math").ToValue(),
						},
					)),
					value.Undefined,
				},
			),
		},
		"recursive fibonacci macro": {
			input: `
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 110,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(218, 13, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(13, 3),
				},
				[]value.Value{
					value.Undefined,
				},
			),
		},
		"define a class": {
			input: `
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
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INSTANTIATE8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(298, 16, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(13, 1),
					bytecode.NewLineInfo(15, 7),
					bytecode.NewLineInfo(16, 4),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(298, 16, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(16, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("BoxString").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.DEF_SETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(298, 16, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(16, 2),
						},
						[]value.Value{
							value.ToSymbol("BoxString").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(230, 13, 5), P(258, 13, 33)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(13, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("value").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.ToSymbol("value").ToValue(),
						},
					)),
					value.ToSymbol("BoxString").ToValue(),
					value.Ref(value.String("foo")),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("value"),
						0,
					)),
				},
			),
		},
		"define a method": {
			input: `
				using Std::Elk::AST::*

				macro reader(name: IdentifierNode, typ: TypeExpressionNode)
					ivar := PublicInstanceVariableNode(name.value)

					quote
						var unquote_ivar(ivar): !{typ.type_node}
						def !{name}: !{typ.type_node}
							!{ivar}
						end
					end
				end

				class Foo
					reader!(bar, type Int | Float?)
				end

				b := Foo()
				b.bar
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(367, 20, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(15, 4),
					bytecode.NewLineInfo(19, 6),
					bytecode.NewLineInfo(20, 4),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(367, 20, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(20, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(367, 20, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(20, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("bar"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(302, 16, 6), P(332, 16, 36)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(16, 3),
								},
								0,
								0,
								[]value.Value{
									value.ToSymbol("bar").ToValue(),
								},
							)),
							value.ToSymbol("bar").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.NewBytecodeFunction(
						value.ToSymbol("<class: Foo>"),
						[]byte{
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(287, 15, 5), P(340, 17, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(16, 1),
							bytecode.NewLineInfo(17, 1),
						},
						0,
						0,
						nil,
					)),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("bar"),
						0,
					)),
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
