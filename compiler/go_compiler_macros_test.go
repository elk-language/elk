package compiler_test

import (
	"testing"
)

func TestGoQuote(t *testing.T) {
	tests := goTestTable{
		"without unquote": {
			input: `
				quote
					1 + 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 ast.Node
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t1 = ast.DeepCopy(ast.NewBinaryExpressionNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(20, 3, 10))), token.New(position.NewLocation("<main>", position.NewSpan(position.New(18, 3, 8), position.New(18, 3, 8))), token.PLUS), ast.NewIntLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 6), position.New(16, 3, 6))), "1"), ast.NewIntLiteralNode(position.NewLocation("<main>", position.NewSpan(position.New(20, 3, 10), position.New(20, 3, 10))), "2")))
}
`,
		},

		// "unquote expression": {
		// 	input: `
		// 		quote
		// 			1 + unquote(5)
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.INT_5),
		// 			byte(bytecode.CALL_METHOD_NT8), 2,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(38, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 3),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewBinaryExpressionNode(
		// 					L(P(16, 3, 6), P(29, 3, 19)),
		// 					T(L(P(18, 3, 8), P(18, 3, 8)), token.PLUS),
		// 					ast.NewIntLiteralNode(L(P(16, 3, 6), P(16, 3, 6)), "1"),
		// 					ast.NewUnquoteNode(
		// 						L(P(20, 3, 10), P(29, 3, 19)),
		// 						ast.UNQUOTE_EXPRESSION_KIND,
		// 						ast.NewIntLiteralNode(L(P(28, 3, 18), P(28, 3, 18)), "5"),
		// 					),
		// 				),
		// 			),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.IntClass, "to_ast_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },
		// "short unquote expression": {
		// 	input: `
		// 		quote
		// 			1 + !{5}
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.INT_5),
		// 			byte(bytecode.CALL_METHOD_NT8), 2,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(32, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 3),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewBinaryExpressionNode(
		// 					L(P(16, 3, 6), P(23, 3, 13)),
		// 					T(L(P(18, 3, 8), P(18, 3, 8)), token.PLUS),
		// 					ast.NewIntLiteralNode(L(P(16, 3, 6), P(16, 3, 6)), "1"),
		// 					ast.NewUnquoteNode(
		// 						L(P(20, 3, 10), P(23, 3, 13)),
		// 						ast.UNQUOTE_EXPRESSION_KIND,
		// 						ast.NewIntLiteralNode(L(P(22, 3, 12), P(22, 3, 12)), "5"),
		// 					),
		// 				),
		// 			),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.IntClass, "to_ast_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },

		// "unquote identifier": {
		// 	input: `
		// 		quote
		// 			var unquote(:foo): String
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.LOAD_VALUE_2),
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 4,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(49, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 3),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewVariableDeclarationNode(
		// 					L(P(16, 3, 6), P(40, 3, 30)),
		// 					"",
		// 					ast.NewUnquoteNode(
		// 						L(P(20, 3, 10), P(32, 3, 22)),
		// 						ast.UNQUOTE_IDENTIFIER_KIND,
		// 						ast.NewSimpleSymbolLiteralNode(L(P(28, 3, 18), P(31, 3, 21)), "foo"),
		// 					),
		// 					ast.NewPublicConstantNode(L(P(35, 3, 25), P(40, 3, 30)), "String"),
		// 					nil,
		// 				),
		// 			),
		// 			value.ToSymbol("foo").ToValue(),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.SymbolClass, "to_ast_ident_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },
		// "short unquote identifier": {
		// 	input: `
		// 		quote
		// 			var !{:foo}: String
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.LOAD_VALUE_2),
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 4,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(43, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 3),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewVariableDeclarationNode(
		// 					L(P(16, 3, 6), P(34, 3, 24)),
		// 					"",
		// 					ast.NewUnquoteNode(
		// 						L(P(20, 3, 10), P(26, 3, 16)),
		// 						ast.UNQUOTE_IDENTIFIER_KIND,
		// 						ast.NewSimpleSymbolLiteralNode(L(P(22, 3, 12), P(25, 3, 15)), "foo"),
		// 					),
		// 					ast.NewPublicConstantNode(L(P(29, 3, 19), P(34, 3, 24)), "String"),
		// 					nil,
		// 				),
		// 			),
		// 			value.ToSymbol("foo").ToValue(),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.SymbolClass, "to_ast_ident_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },
		// "unquote_ident": {
		// 	input: `
		// 		quote
		// 			var unquote_ident(:foo): String
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.LOAD_VALUE_2),
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 4,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(55, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 3),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewVariableDeclarationNode(
		// 					L(P(16, 3, 6), P(46, 3, 36)),
		// 					"",
		// 					ast.NewUnquoteNode(
		// 						L(P(20, 3, 10), P(38, 3, 28)),
		// 						ast.UNQUOTE_IDENTIFIER_KIND,
		// 						ast.NewSimpleSymbolLiteralNode(L(P(34, 3, 24), P(37, 3, 27)), "foo"),
		// 					),
		// 					ast.NewPublicConstantNode(L(P(41, 3, 31), P(46, 3, 36)), "String"),
		// 					nil,
		// 				),
		// 			),
		// 			value.ToSymbol("foo").ToValue(),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.SymbolClass, "to_ast_ident_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },

		// "unquote pattern expression": {
		// 	input: `
		// 		quote
		// 			var ^[unquote(1)] = "foo"
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.CALL_METHOD_NT8), 2,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(49, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 3),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewVariablePatternDeclarationNode(
		// 					L(P(16, 3, 6), P(40, 3, 30)),
		// 					ast.NewSetPatternNode(
		// 						L(P(20, 3, 10), P(32, 3, 22)),
		// 						[]ast.PatternNode{
		// 							ast.NewUnquoteNode(
		// 								L(P(22, 3, 12), P(31, 3, 21)),
		// 								ast.UNQUOTE_PATTERN_EXPRESSION_KIND,
		// 								ast.NewIntLiteralNode(L(P(30, 3, 20), P(30, 3, 20)), "1"),
		// 							),
		// 						},
		// 					),
		// 					ast.NewDoubleQuotedStringLiteralNode(L(P(36, 3, 26), P(40, 3, 30)), "foo"),
		// 				),
		// 			),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.IntClass, "to_ast_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },
		// "short unquote pattern expression": {
		// 	input: `
		// 		quote
		// 			var ^[!{1}] = "foo"
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.CALL_METHOD_NT8), 2,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(43, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 3),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewVariablePatternDeclarationNode(
		// 					L(P(16, 3, 6), P(34, 3, 24)),
		// 					ast.NewSetPatternNode(
		// 						L(P(20, 3, 10), P(26, 3, 16)),
		// 						[]ast.PatternNode{
		// 							ast.NewUnquoteNode(
		// 								L(P(22, 3, 12), P(25, 3, 15)),
		// 								ast.UNQUOTE_PATTERN_EXPRESSION_KIND,
		// 								ast.NewIntLiteralNode(L(P(24, 3, 14), P(24, 3, 14)), "1"),
		// 							),
		// 						},
		// 					),
		// 					ast.NewDoubleQuotedStringLiteralNode(L(P(30, 3, 20), P(34, 3, 24)), "foo"),
		// 				),
		// 			),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.IntClass, "to_ast_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },

		// "unquote pattern": {
		// 	input: `
		// 		quote
		// 			var [unquote(Elk::AST::ListPatternNode())] = "foo"
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.GET_CONST8), 2,
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.INSTANTIATE8), 2,
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 4,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(74, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 8),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewVariablePatternDeclarationNode(
		// 					L(P(16, 3, 6), P(65, 3, 55)),
		// 					ast.NewListPatternNode(
		// 						L(P(20, 3, 10), P(57, 3, 47)),
		// 						[]ast.PatternNode{
		// 							ast.NewUnquoteNode(
		// 								L(P(21, 3, 11), P(56, 3, 46)),
		// 								ast.UNQUOTE_PATTERN_KIND,
		// 								ast.NewConstructorCallNode(
		// 									L(P(29, 3, 19), P(55, 3, 45)),
		// 									ast.NewPublicConstantNode(
		// 										L(P(29, 3, 19), P(53, 3, 43)),
		// 										"Std::Elk::AST::ListPatternNode",
		// 									),
		// 									[]ast.ExpressionNode{
		// 										ast.NewUndefinedLiteralNode(L(P(29, 3, 19), P(53, 3, 43))),
		// 										ast.NewUndefinedLiteralNode(L(P(29, 3, 19), P(53, 3, 43))),
		// 									},
		// 									nil,
		// 								),
		// 							),
		// 						},
		// 					),
		// 					ast.NewDoubleQuotedStringLiteralNode(L(P(61, 3, 51), P(65, 3, 55)), "foo"),
		// 				),
		// 			),
		// 			value.ToSymbol("Std::Elk::AST::ListPatternNode").ToValue(),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.PatternNodeMixin, "to_ast_pattern_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },
		// "unquote_pattern": {
		// 	input: `
		// 		quote
		// 			var [unquote_pattern(Elk::AST::ListPatternNode())] = "foo"
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.GET_CONST8), 2,
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.INSTANTIATE8), 2,
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 4,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(82, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 8),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewVariablePatternDeclarationNode(
		// 					L(P(16, 3, 6), P(73, 3, 63)),
		// 					ast.NewListPatternNode(
		// 						L(P(20, 3, 10), P(65, 3, 55)),
		// 						[]ast.PatternNode{
		// 							ast.NewUnquoteNode(
		// 								L(P(21, 3, 11), P(64, 3, 54)),
		// 								ast.UNQUOTE_PATTERN_KIND,
		// 								ast.NewConstructorCallNode(
		// 									L(P(37, 3, 27), P(63, 3, 53)),
		// 									ast.NewPublicConstantNode(
		// 										L(P(37, 3, 27), P(61, 3, 51)),
		// 										"Std::Elk::AST::ListPatternNode",
		// 									),
		// 									[]ast.ExpressionNode{
		// 										ast.NewUndefinedLiteralNode(L(P(37, 3, 27), P(61, 3, 51))),
		// 										ast.NewUndefinedLiteralNode(L(P(37, 3, 27), P(61, 3, 51))),
		// 									},
		// 									nil,
		// 								),
		// 							),
		// 						},
		// 					),
		// 					ast.NewDoubleQuotedStringLiteralNode(L(P(69, 3, 59), P(73, 3, 63)), "foo"),
		// 				),
		// 			),
		// 			value.ToSymbol("Std::Elk::AST::ListPatternNode").ToValue(),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.PatternNodeMixin, "to_ast_pattern_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },
		// "short unquote pattern": {
		// 	input: `
		// 		quote
		// 			var [!{Elk::AST::ListPatternNode()}] = "foo"
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.GET_CONST8), 2,
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.INSTANTIATE8), 2,
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 4,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(68, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 8),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewVariablePatternDeclarationNode(
		// 					L(P(16, 3, 6), P(59, 3, 49)),
		// 					ast.NewListPatternNode(
		// 						L(P(20, 3, 10), P(51, 3, 41)),
		// 						[]ast.PatternNode{
		// 							ast.NewUnquoteNode(
		// 								L(P(21, 3, 11), P(50, 3, 40)),
		// 								ast.UNQUOTE_PATTERN_KIND,
		// 								ast.NewConstructorCallNode(
		// 									L(P(23, 3, 13), P(49, 3, 39)),
		// 									ast.NewPublicConstantNode(
		// 										L(P(23, 3, 13), P(47, 3, 37)),
		// 										"Std::Elk::AST::ListPatternNode",
		// 									),
		// 									[]ast.ExpressionNode{
		// 										ast.NewUndefinedLiteralNode(L(P(23, 3, 13), P(47, 3, 37))),
		// 										ast.NewUndefinedLiteralNode(L(P(23, 3, 13), P(47, 3, 37))),
		// 									},
		// 									nil,
		// 								),
		// 							),
		// 						},
		// 					),
		// 					ast.NewDoubleQuotedStringLiteralNode(L(P(55, 3, 45), P(59, 3, 49)), "foo"),
		// 				),
		// 			),
		// 			value.ToSymbol("Std::Elk::AST::ListPatternNode").ToValue(),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.PatternNodeMixin, "to_ast_pattern_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },

		// "unquote constant": {
		// 	input: `
		// 		quote
		// 			const unquote(:Bar) = "foo"
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.LOAD_VALUE_2),
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 4,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(51, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 3),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewConstantDeclarationNode(
		// 					L(P(16, 3, 6), P(42, 3, 32)),
		// 					"",
		// 					ast.NewUnquoteNode(
		// 						L(P(22, 3, 12), P(34, 3, 24)),
		// 						ast.UNQUOTE_CONSTANT_KIND,
		// 						ast.NewSimpleSymbolLiteralNode(L(P(30, 3, 20), P(33, 3, 23)), "Bar"),
		// 					),
		// 					nil,
		// 					ast.NewDoubleQuotedStringLiteralNode(L(P(38, 3, 28), P(42, 3, 32)), "foo"),
		// 				),
		// 			),
		// 			value.ToSymbol("Bar").ToValue(),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.SymbolClass, "to_ast_const_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },
		// "unquote_const": {
		// 	input: `
		// 		quote
		// 			const unquote_const(:Bar) = "foo"
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.LOAD_VALUE_2),
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 4,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(57, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 3),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewConstantDeclarationNode(
		// 					L(P(16, 3, 6), P(48, 3, 38)),
		// 					"",
		// 					ast.NewUnquoteNode(
		// 						L(P(22, 3, 12), P(40, 3, 30)),
		// 						ast.UNQUOTE_CONSTANT_KIND,
		// 						ast.NewSimpleSymbolLiteralNode(L(P(36, 3, 26), P(39, 3, 29)), "Bar"),
		// 					),
		// 					nil,
		// 					ast.NewDoubleQuotedStringLiteralNode(L(P(44, 3, 34), P(48, 3, 38)), "foo"),
		// 				),
		// 			),
		// 			value.ToSymbol("Bar").ToValue(),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.SymbolClass, "to_ast_const_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },
		// "short unquote constant": {
		// 	input: `
		// 		quote
		// 			const !{:Bar} = "foo"
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.LOAD_VALUE_2),
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 4,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(45, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 3),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewConstantDeclarationNode(
		// 					L(P(16, 3, 6), P(36, 3, 26)),
		// 					"",
		// 					ast.NewUnquoteNode(
		// 						L(P(22, 3, 12), P(28, 3, 18)),
		// 						ast.UNQUOTE_CONSTANT_KIND,
		// 						ast.NewSimpleSymbolLiteralNode(L(P(24, 3, 14), P(27, 3, 17)), "Bar"),
		// 					),
		// 					nil,
		// 					ast.NewDoubleQuotedStringLiteralNode(L(P(32, 3, 22), P(36, 3, 26)), "foo"),
		// 				),
		// 			),
		// 			value.ToSymbol("Bar").ToValue(),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.SymbolClass, "to_ast_const_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },

		// "unquote_ivar": {
		// 	input: `
		// 		quote
		// 			var unquote_ivar(:foo): String?
		// 		end
		// 	`,
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		nil,
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.GET_CONST8), 0,
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.UNDEFINED),
		// 			byte(bytecode.LOAD_VALUE_2),
		// 			byte(bytecode.CALL_METHOD_NT8), 3,
		// 			byte(bytecode.NEW_ARRAY_TUPLE8), 1,
		// 			byte(bytecode.CALL_METHOD_NT8), 4,
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(55, 4, 8)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 0),
		// 			bytecode.NewLineInfo(2, 4),
		// 			bytecode.NewLineInfo(3, 3),
		// 			bytecode.NewLineInfo(4, 2),
		// 			bytecode.NewLineInfo(2, 2),
		// 			bytecode.NewLineInfo(4, 1),
		// 		},
		// 		[]value.Value{
		// 			value.ToSymbol("Std::Kernel").ToValue(),
		// 			value.Ref(
		// 				ast.NewInstanceVariableDeclarationNode(
		// 					L(P(16, 3, 6), P(46, 3, 36)),
		// 					"",
		// 					ast.NewUnquoteNode(
		// 						L(P(20, 3, 10), P(37, 3, 27)),
		// 						ast.UNQUOTE_INSTANCE_VARIABLE_KIND,
		// 						ast.NewSimpleSymbolLiteralNode(L(P(33, 3, 23), P(36, 3, 26)), "foo"),
		// 					),
		// 					ast.NewNilableTypeNode(
		// 						L(P(40, 3, 30), P(46, 3, 36)),
		// 						ast.NewPublicConstantNode(L(P(40, 3, 30), P(45, 3, 35)), "String"),
		// 					),
		// 				),
		// 			),
		// 			value.ToSymbol("foo").ToValue(),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.SymbolClass, "to_ast_ivar_node"), 0)),
		// 			value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "#splice"), 2)),
		// 		},
		// 	),
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoMacroExpansion(t *testing.T) {
	tests := goTestTable{
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

				a := fib!(10) * 2
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(110)).ToValue()
}
`,
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

				a := Math::fib!(10) * 2
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Module // Math
var sym0 = value.ToSymbol("Math")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(110)).ToValue()
}

func initGlobalEnv() {
	var parentNamespace value.Value
	_ = parentNamespace
	var namespace value.Value
	_ = namespace
	var class *value.Class
	_ = class
	var superclass *value.Class
	_ = superclass
	var mixin *value.Mixin
	_ = mixin

	parentNamespace = (value.RootModule).ToValue()
	const0 = value.NewModule()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}
`,
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

				a := fib!(10) * 2
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(110)).ToValue()
}
`,
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
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym6 = value.ToSymbol("#init")
var fn_method1 vm.NativeFunction // BoxString.:value

var const0 *value.Class // BoxString
var sym0 = value.ToSymbol("BoxString")

var sym1 = value.ToSymbol("BoxString.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: BoxString.:#init, loc: <main>:13:5
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("value")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var b: BoxString
	_ = l0
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var t2 []value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	fn_method1 = vm.MethodToFunc((const0).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
	defer thread.PopNativeCallFrame()
	t1, err = fn_method0(thread, const0.CreateInstance(), (value.String("foo")).ToValue()) // receiver: BoxString, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l0 = t1
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = l0
	callFrame.SetNativeLineNumber(16)
	_, err = fn_method1(thread, t2) // receiver: BoxString, name: value
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}

func initGlobalEnv() {
	var parentNamespace value.Value
	_ = parentNamespace
	var namespace value.Value
	_ = namespace
	var class *value.Class
	_ = class
	var superclass *value.Class
	_ = superclass
	var mixin *value.Mixin
	_ = mixin

	parentNamespace = (value.RootModule).ToValue()
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("value"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // BoxString
	vm.Def(class, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
	vm.DefineGetter(class, sym3, 0)
	vm.DefineSetter(class, sym3, 0)
}
`,
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
				r := b.bar
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym3 = value.ToSymbol("main")
var sym5 = value.ToSymbol("bar")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:bar")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:bar, loc: <main>:16:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.GetInstanceVariable(self, 0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var b: Foo
	_ = l0
	var l1 value.Value // var r: Std::Int | Std::Float | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = const0.CreateInstance()
	t1, err = fn_method0(thread, l0) // receiver: Foo, name: bar
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}

func initGlobalEnv() {
	var parentNamespace value.Value
	_ = parentNamespace
	var namespace value.Value
	_ = namespace
	var class *value.Class
	_ = class
	var superclass *value.Class
	_ = superclass
	var mixin *value.Mixin
	_ = mixin

	parentNamespace = (value.RootModule).ToValue()
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("bar"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(class, "bar", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}
