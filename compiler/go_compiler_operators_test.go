package compiler_test

import (
	"testing"
)

// func TestBytecodeBinaryExpressions(t *testing.T) {
// 	tests := bytecodeTestTable{
// 		"is a": {
// 			input: "3 <: ::Std::Int",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.INT_3),
// 					byte(bytecode.GET_CONST8), 0,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(14, 1, 15)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 5),
// 				},
// 				[]value.Value{
// 					value.ToSymbol("Std::Int").ToValue(),
// 				},
// 			),
// 			err: diagnostic.DiagnosticList{
// 				diagnostic.NewWarning(L(P(0, 1, 1), P(0, 1, 1)), "this \"is a\" check is always true, `3` will always be an instance of `Std::Int`"),
// 			},
// 		},
// 		"instance of": {
// 			input: "3 <<: ::Std::Int",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.INT_3),
// 					byte(bytecode.GET_CONST8), 0,
// 					byte(bytecode.INSTANCE_OF),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(15, 1, 16)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 5),
// 				},
// 				[]value.Value{
// 					value.ToSymbol("Std::Int").ToValue(),
// 				},
// 			),
// 			err: diagnostic.DiagnosticList{
// 				diagnostic.NewWarning(L(P(0, 1, 1), P(0, 1, 1)), "this \"instance of\" check is always true, `3` will always be an instance of `Std::Int`"),
// 			},
// 		},
// 		"resolve static add": {
// 			input: "1i8 + 5i8",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_INT8), 6,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(8, 1, 9)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"add int": {
// 			input: "a := 1; a + 5",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.INT_1),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.INT_5),
// 					byte(bytecode.ADD_INT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(12, 1, 13)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"add float": {
// 			input: "a := 1.2; a + 5.0",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.ADD_FLOAT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(16, 1, 17)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{
// 					value.Float(1.2).ToValue(),
// 					value.Float(5.0).ToValue(),
// 				},
// 			),
// 		},
// 		"add builtin": {
// 			input: "a := 1i8; a + 5i8",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_INT8), 1,
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.LOAD_INT8), 5,
// 					byte(bytecode.ADD),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(16, 1, 17)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 10),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"add value": {
// 			input: `
// 				module Foo
// 					def +(other: Int): Int
// 					  other
// 					end
// 				end
// 				Foo + 5
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_CONST8), 2,
// 					byte(bytecode.INT_5),
// 					byte(bytecode.CALL_METHOD8), 3,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 12)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 6),
// 					bytecode.NewLineInfo(7, 6),
// 				},
// 				[]value.Value{
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							namespaceDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.DEF_NAMESPACE), 0,
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(85, 7, 12)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 5),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Root").ToValue(),
// 								value.ToSymbol("Foo").ToValue(),
// 							},
// 						),
// 					),
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							methodDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.GET_SINGLETON),
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.LOAD_VALUE_2),
// 								byte(bytecode.DEF_METHOD),
// 								byte(bytecode.POP),
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(85, 7, 12)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 7),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Foo").ToValue(),
// 								value.Ref(
// 									vm.NewBytecodeFunction(
// 										value.ToSymbol("Foo::+"),
// 										[]byte{
// 											byte(bytecode.GET_LOCAL_1),
// 											byte(bytecode.RETURN),
// 										},
// 										L(P(21, 3, 6), P(64, 5, 8)),
// 										bytecode.LineInfoList{
// 											bytecode.NewLineInfo(4, 1),
// 											bytecode.NewLineInfo(5, 1),
// 										},
// 										1,
// 										0,
// 										nil,
// 									),
// 								),
// 								value.ToSymbol("+").ToValue(),
// 							},
// 						),
// 					),
// 					value.ToSymbol("Foo").ToValue(),
// 					value.Ref(value.NewCallSiteInfo(symbol.OpAdd, 1)),
// 				},
// 			),
// 		},

// 		"resolve static subtract": {
// 			input: "151i32 - 25i32 - 5i32",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_INT32_8), 0x79,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(20, 1, 21)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"subtract int": {
// 			input: "a := 1; a - 5",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.INT_1),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.INT_5),
// 					byte(bytecode.SUBTRACT_INT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(12, 1, 13)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"subtract float": {
// 			input: "a := 1.2; a - 5.0",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.SUBTRACT_FLOAT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(16, 1, 17)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{
// 					value.Float(1.2).ToValue(),
// 					value.Float(5.0).ToValue(),
// 				},
// 			),
// 		},
// 		"subtract builtin": {
// 			input: "a := 151i32; a - 25i32 - 5i32",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.LOAD_INT32_8), 25,
// 					byte(bytecode.SUBTRACT),
// 					byte(bytecode.LOAD_INT32_8), 5,
// 					byte(bytecode.SUBTRACT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(28, 1, 29)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 12),
// 				},
// 				[]value.Value{
// 					value.Int32(151).ToValue(),
// 				},
// 			),
// 		},
// 		"subtract value": {
// 			input: `
// 				module Foo
// 					def -(other: Int): Int
// 					  other
// 					end
// 				end
// 				Foo - 5
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_CONST8), 2,
// 					byte(bytecode.INT_5),
// 					byte(bytecode.CALL_METHOD8), 3,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 12)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 6),
// 					bytecode.NewLineInfo(7, 6),
// 				},
// 				[]value.Value{
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							namespaceDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.DEF_NAMESPACE), 0,
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(85, 7, 12)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 5),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Root").ToValue(),
// 								value.ToSymbol("Foo").ToValue(),
// 							},
// 						),
// 					),
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							methodDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.GET_SINGLETON),
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.LOAD_VALUE_2),
// 								byte(bytecode.DEF_METHOD),
// 								byte(bytecode.POP),
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(85, 7, 12)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 7),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Foo").ToValue(),
// 								value.Ref(
// 									vm.NewBytecodeFunction(
// 										value.ToSymbol("Foo::-"),
// 										[]byte{
// 											byte(bytecode.GET_LOCAL_1),
// 											byte(bytecode.RETURN),
// 										},
// 										L(P(21, 3, 6), P(64, 5, 8)),
// 										bytecode.LineInfoList{
// 											bytecode.NewLineInfo(4, 1),
// 											bytecode.NewLineInfo(5, 1),
// 										},
// 										1,
// 										0,
// 										nil,
// 									),
// 								),
// 								value.ToSymbol("-").ToValue(),
// 							},
// 						),
// 					),
// 					value.ToSymbol("Foo").ToValue(),
// 					value.Ref(value.NewCallSiteInfo(symbol.OpSubtract, 1)),
// 				},
// 			),
// 		},

// 		"resolve static multiply": {
// 			input: "45.5 * 2.5",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(9, 1, 10)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 2),
// 				},
// 				[]value.Value{
// 					value.Float(113.75).ToValue(),
// 				},
// 			),
// 		},
// 		"multiply int": {
// 			input: "a := 1; a * 5",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.INT_1),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.INT_5),
// 					byte(bytecode.MULTIPLY_INT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(12, 1, 13)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"multiply float": {
// 			input: "a := 1.2; a * 5.0",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.MULTIPLY_FLOAT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(16, 1, 17)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{
// 					value.Float(1.2).ToValue(),
// 					value.Float(5.0).ToValue(),
// 				},
// 			),
// 		},
// 		"multiply builtin": {
// 			input: "a := 45i8; a * 2i8",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_INT8), 45,
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.LOAD_INT8), 2,
// 					byte(bytecode.MULTIPLY),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(17, 1, 18)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 10),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"multiply value": {
// 			input: `
// 				module Foo
// 					def *(other: Int): Int
// 					  other
// 					end
// 				end
// 				Foo * 5
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_CONST8), 2,
// 					byte(bytecode.INT_5),
// 					byte(bytecode.CALL_METHOD8), 3,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 12)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 6),
// 					bytecode.NewLineInfo(7, 6),
// 				},
// 				[]value.Value{
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							namespaceDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.DEF_NAMESPACE), 0,
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(85, 7, 12)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 5),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Root").ToValue(),
// 								value.ToSymbol("Foo").ToValue(),
// 							},
// 						),
// 					),
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							methodDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.GET_SINGLETON),
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.LOAD_VALUE_2),
// 								byte(bytecode.DEF_METHOD),
// 								byte(bytecode.POP),
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(85, 7, 12)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 7),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Foo").ToValue(),
// 								value.Ref(
// 									vm.NewBytecodeFunction(
// 										value.ToSymbol("Foo::*"),
// 										[]byte{
// 											byte(bytecode.GET_LOCAL_1),
// 											byte(bytecode.RETURN),
// 										},
// 										L(P(21, 3, 6), P(64, 5, 8)),
// 										bytecode.LineInfoList{
// 											bytecode.NewLineInfo(4, 1),
// 											bytecode.NewLineInfo(5, 1),
// 										},
// 										1,
// 										0,
// 										nil,
// 									),
// 								),
// 								value.ToSymbol("*").ToValue(),
// 							},
// 						),
// 					),
// 					value.ToSymbol("Foo").ToValue(),
// 					value.Ref(value.NewCallSiteInfo(symbol.OpMultiply, 1)),
// 				},
// 			),
// 		},

// 		"resolve static divide": {
// 			input: "45.5 / .5",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(8, 1, 9)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 2),
// 				},
// 				[]value.Value{
// 					value.Float(91).ToValue(),
// 				},
// 			),
// 		},
// 		"divide int": {
// 			input: "a := 1; a / 5",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.INT_1),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.INT_5),
// 					byte(bytecode.DIVIDE_INT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(12, 1, 13)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"divide float": {
// 			input: "a := 45.5; a / .5",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.DIVIDE_FLOAT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(16, 1, 17)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{
// 					value.Float(45.5).ToValue(),
// 					value.Float(0.5).ToValue(),
// 				},
// 			),
// 		},
// 		"divide builtin": {
// 			input: "a := 1i8; a / 5i8",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_INT8), 1,
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.LOAD_INT8), 5,
// 					byte(bytecode.DIVIDE),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(16, 1, 17)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 10),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"divide value": {
// 			input: `
// 				module Foo
// 					def /(other: Int): Int
// 					  other
// 					end
// 				end
// 				Foo / 5
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_CONST8), 2,
// 					byte(bytecode.INT_5),
// 					byte(bytecode.CALL_METHOD8), 3,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 12)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 6),
// 					bytecode.NewLineInfo(7, 6),
// 				},
// 				[]value.Value{
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							namespaceDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.DEF_NAMESPACE), 0,
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(85, 7, 12)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 5),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Root").ToValue(),
// 								value.ToSymbol("Foo").ToValue(),
// 							},
// 						),
// 					),
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							methodDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.GET_SINGLETON),
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.LOAD_VALUE_2),
// 								byte(bytecode.DEF_METHOD),
// 								byte(bytecode.POP),
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(85, 7, 12)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 7),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Foo").ToValue(),
// 								value.Ref(
// 									vm.NewBytecodeFunction(
// 										value.ToSymbol("Foo::/"),
// 										[]byte{
// 											byte(bytecode.GET_LOCAL_1),
// 											byte(bytecode.RETURN),
// 										},
// 										L(P(21, 3, 6), P(64, 5, 8)),
// 										bytecode.LineInfoList{
// 											bytecode.NewLineInfo(4, 1),
// 											bytecode.NewLineInfo(5, 1),
// 										},
// 										1,
// 										0,
// 										nil,
// 									),
// 								),
// 								value.ToSymbol("/").ToValue(),
// 							},
// 						),
// 					),
// 					value.ToSymbol("Foo").ToValue(),
// 					value.Ref(value.NewCallSiteInfo(symbol.OpDivide, 1)),
// 				},
// 			),
// 		},

// 		"resolve static exponentiate": {
// 			input: "-2 ** 2",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_INT_8), 0xFC,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(6, 1, 7)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"exponentiate int": {
// 			input: "a := -2; a ** 2",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_INT_8), 0xFE,
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.INT_2),
// 					byte(bytecode.EXPONENTIATE_INT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(14, 1, 15)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"exponentiate float": {
// 			input: "a := 1.2; a + 5.0",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.ADD_FLOAT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(16, 1, 17)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{
// 					value.Float(1.2).ToValue(),
// 					value.Float(5.0).ToValue(),
// 				},
// 			),
// 		},
// 		"exponentiate builtin": {
// 			input: "a := 1i8; a ** 5i8",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_INT8), 1,
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.LOAD_INT8), 5,
// 					byte(bytecode.EXPONENTIATE),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(17, 1, 18)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 10),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"exponentiate value": {
// 			input: `
// 				module Foo
// 					def **(other: Int): Int
// 					  other
// 					end
// 				end
// 				Foo ** 5
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_CONST8), 2,
// 					byte(bytecode.INT_5),
// 					byte(bytecode.CALL_METHOD8), 3,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(87, 7, 13)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 6),
// 					bytecode.NewLineInfo(7, 6),
// 				},
// 				[]value.Value{
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							namespaceDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.DEF_NAMESPACE), 0,
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(87, 7, 13)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 5),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Root").ToValue(),
// 								value.ToSymbol("Foo").ToValue(),
// 							},
// 						),
// 					),
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							methodDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.GET_SINGLETON),
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.LOAD_VALUE_2),
// 								byte(bytecode.DEF_METHOD),
// 								byte(bytecode.POP),
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(87, 7, 13)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 7),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Foo").ToValue(),
// 								value.Ref(
// 									vm.NewBytecodeFunction(
// 										value.ToSymbol("Foo::**"),
// 										[]byte{
// 											byte(bytecode.GET_LOCAL_1),
// 											byte(bytecode.RETURN),
// 										},
// 										L(P(21, 3, 6), P(65, 5, 8)),
// 										bytecode.LineInfoList{
// 											bytecode.NewLineInfo(4, 1),
// 											bytecode.NewLineInfo(5, 1),
// 										},
// 										1,
// 										0,
// 										nil,
// 									),
// 								),
// 								value.ToSymbol("**").ToValue(),
// 							},
// 						),
// 					),
// 					value.ToSymbol("Foo").ToValue(),
// 					value.Ref(value.NewCallSiteInfo(symbol.OpExponentiate, 1)),
// 				},
// 			),
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			bytecodeCompilerTest(tc, t)
// 		})
// 	}
// }

// func TestBytecodeUnaryExpressions(t *testing.T) {
// 	tests := bytecodeTestTable{
// 		"resolve static negate": {
// 			input: "-5",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_INT_8), 0xfb,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(1, 1, 2)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"negate int": {
// 			input: "a := 5; -a",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.INT_5),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.NEGATE_INT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(9, 1, 10)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 7),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"negate float": {
// 			input: "a := 5.2; -a",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.NEGATE_FLOAT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(11, 1, 12)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 7),
// 				},
// 				[]value.Value{
// 					value.Float(5.2).ToValue(),
// 				},
// 			),
// 		},
// 		"negate builtin": {
// 			input: "a := 5i8; -a",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_INT8), 5,
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.NEGATE),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(11, 1, 12)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"negate value": {
// 			input: `
// 				module Foo
// 					def -@: Foo
// 					  self
// 					end
// 				end
// 				-Foo
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_CONST8), 2,
// 					byte(bytecode.CALL_METHOD8), 3,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(70, 7, 9)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 6),
// 					bytecode.NewLineInfo(7, 5),
// 				},
// 				[]value.Value{
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							namespaceDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.DEF_NAMESPACE), 0,
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(70, 7, 9)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 5),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Root").ToValue(),
// 								value.ToSymbol("Foo").ToValue(),
// 							},
// 						),
// 					),
// 					value.Ref(
// 						vm.NewBytecodeFunctionNoParams(
// 							methodDefinitionsSymbol,
// 							[]byte{
// 								byte(bytecode.GET_CONST8), 0,
// 								byte(bytecode.GET_SINGLETON),
// 								byte(bytecode.LOAD_VALUE_1),
// 								byte(bytecode.LOAD_VALUE_2),
// 								byte(bytecode.DEF_METHOD),
// 								byte(bytecode.POP),
// 								byte(bytecode.NIL),
// 								byte(bytecode.RETURN),
// 							},
// 							L(P(0, 1, 1), P(70, 7, 9)),
// 							bytecode.LineInfoList{
// 								bytecode.NewLineInfo(1, 7),
// 								bytecode.NewLineInfo(7, 2),
// 							},
// 							[]value.Value{
// 								value.ToSymbol("Foo").ToValue(),
// 								value.Ref(
// 									vm.NewBytecodeFunction(
// 										value.ToSymbol("Foo::-@"),
// 										[]byte{
// 											byte(bytecode.SELF),
// 											byte(bytecode.RETURN),
// 										},
// 										L(P(21, 3, 6), P(52, 5, 8)),
// 										bytecode.LineInfoList{
// 											bytecode.NewLineInfo(4, 1),
// 											bytecode.NewLineInfo(5, 1),
// 										},
// 										0,
// 										0,
// 										nil,
// 									),
// 								),
// 								symbol.OpNegate.ToValue(),
// 							},
// 						),
// 					),
// 					value.ToSymbol("Foo").ToValue(),
// 					value.Ref(value.NewCallSiteInfo(symbol.OpNegate, 0)),
// 				},
// 			),
// 		},

// 		"resolve static bitwise not": {
// 			input: "~10",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_INT_8), 0xf5,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(2, 1, 3)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"resolve static logical not": {
// 			input: "!10",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.FALSE),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(2, 1, 3)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 2),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"logical not": {
// 			input: "a := 10; !a",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_INT_8), 10,
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.NOT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(10, 1, 11)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"bitwise not": {
// 			input: "a := 10; ~a",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_INT_8), 10,
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.BITWISE_NOT),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(10, 1, 11)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{},
// 			),
// 		},

// 		"resolve static plus": {
// 			input: "+5",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.INT_5),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(1, 1, 2)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 2),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 		"unary plus": {
// 			input: "a := 10; +a",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_INT_8), 10,
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.UNARY_PLUS),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(10, 1, 11)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{},
// 			),
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			bytecodeCompilerTest(tc, t)
// 		})
// 	}
// }

func TestGoComplexAssignmentLocals(t *testing.T) {
	tests := goTestTable{
		// TODO:
		// 		"increment": {
		// 			input: "a := 1; a++",
		// 			want: `
		// `,
		// 		},
		// "decrement": {
		// 	input: "a := 1; a--",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.DECREMENT_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(10, 1, 11)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 9),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "add": {
		// 	input: "a := 1; a += 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.ADD_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "subtract": {
		// 	input: "a := 1; a -= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.SUBTRACT_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "multiply": {
		// 	input: "a := 1; a *= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.MULTIPLY_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "divide": {
		// 	input: "a := 1; a /= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.DIVIDE_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "exponentiate": {
		// 	input: "a := 1; a **= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.EXPONENTIATE_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(14, 1, 15)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "modulo": {
		// 	input: "a := 1; a %= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.MODULO_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "bitwise AND": {
		// 	input: "a := 1; a &= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.BITWISE_AND_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "bitwise OR": {
		// 	input: "a := 1; a |= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.BITWISE_OR_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "bitwise XOR": {
		// 	input: "a := 1; a ^= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.BITWISE_XOR_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(13, 1, 14)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "left bitshift": {
		// 	input: "a := 1; a <<= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.LBITSHIFT_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(14, 1, 15)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "left logical bitshift": {
		// 	input: "a := 1u64; a <<<= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.LOAD_UINT64_8), 1,
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.LOGIC_LBITSHIFT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(18, 1, 19)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 11),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "right bitshift": {
		// 	input: "a := 1; a >>= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.RBITSHIFT_INT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(14, 1, 15)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 10),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "right logical bitshift": {
		// 	input: "a := 1u64; a >>>= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.LOAD_UINT64_8), 1,
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.LOGIC_RBITSHIFT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(18, 1, 19)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 11),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "logic OR": {
		// 	input: "var a: Int? = 1; a ||= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.JUMP_IF_NP), 0, 2,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(23, 1, 24)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 13),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "logic AND": {
		// 	input: "var a: Int? = 1; a &&= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.JUMP_UNLESS_NP), 0, 2,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(23, 1, 24)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 13),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
		// "nil coalesce": {
		// 	input: "var a: Int? = 1; a ??= 3",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.GET_LOCAL_1),
		// 			byte(bytecode.JUMP_UNLESS_NNP), 0, 2,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.INT_3),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(23, 1, 24)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 13),
		// 		},
		// 		[]value.Value{},
		// 	),
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

// func TestBytecodeComplexAssignmentInstanceVariables(t *testing.T) {
// 	tests := bytecodeTestTable{
// 		"increment": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a++
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(82, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_NAME16), 0, 0,
// 									byte(bytecode.INCREMENT_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(73, 6, 22)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 7),
// 								},
// 								[]value.Value{
// 									value.ToSymbol("a").ToValue(),
// 								},
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"decrement": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a--
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(82, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(82, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_NAME16), 0, 0,
// 									byte(bytecode.DECREMENT_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(73, 6, 22)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 7),
// 								},
// 								[]value.Value{
// 									value.ToSymbol("a").ToValue(),
// 								},
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"add": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a += 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.ADD_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"subtract": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a -= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.SUBTRACT_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"multiply": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a *= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.MULTIPLY_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"divide": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a /= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.DIVIDE_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"exponentiate": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a **= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(86, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.EXPONENTIATE_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(77, 6, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"modulo": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a %= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.MODULO_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"bitwise AND": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a &= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.BITWISE_AND_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"bitwise OR": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a |= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.BITWISE_OR_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"bitwise XOR": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a ^= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(85, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.BITWISE_XOR_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(76, 6, 25)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"left bitshift": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a <<= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(86, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.LBITSHIFT_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(77, 6, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"left logical bitshift": {
// 			input: `
// 				class Foo
// 					var @a: UInt64
// 					init(@a); end

// 					def foo then @a <<<= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(90, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(90, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(90, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(90, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(40, 4, 6), P(52, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.LOGIC_LBITSHIFT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(60, 6, 6), P(81, 6, 27)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"right bitshift": {
// 			input: `
// 				class Foo
// 					var @a: Int
// 					init(@a); end

// 					def foo then @a >>= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(86, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(37, 4, 6), P(49, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.RBITSHIFT_INT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(57, 6, 6), P(77, 6, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"right logical bitshift": {
// 			input: `
// 				class Foo
// 					var @a: Int64
// 					init(@a); end

// 					def foo then @a >>>= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(89, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(7, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(89, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(89, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.LOAD_VALUE8), 4,
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(89, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("Foo.:#init"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.POP),
// 									byte(bytecode.RETURN_SELF),
// 								},
// 								L(P(39, 4, 6), P(51, 4, 18)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 7),
// 								},
// 								1,
// 								0,
// 								nil,
// 							)),
// 							value.ToSymbol("#init").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.LOGIC_RBITSHIFT),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(59, 6, 6), P(80, 6, 27)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(6, 6),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"logic OR": {
// 			input: `
// 				class Foo
// 					var @a: Int?

// 					def foo then @a ||= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(68, 6, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(6, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 6),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.JUMP_IF_NP), 0, 2,
// 									byte(bytecode.POP),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(39, 5, 6), P(59, 5, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(5, 9),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"nil coalesce": {
// 			input: `
// 				class Foo
// 					var @a: Int?

// 					def foo then @a ??= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(68, 6, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(6, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 6),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.JUMP_UNLESS_NNP), 0, 2,
// 									byte(bytecode.POP),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(39, 5, 6), P(59, 5, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(5, 9),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"logic AND": {
// 			input: `
// 				class Foo
// 					var @a: Int?

// 					def foo then @a &&= 3
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(68, 6, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 9),
// 					bytecode.NewLineInfo(6, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<ivarIndices>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_IVARS),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 4),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(&value.IvarIndices{
// 								value.ToSymbol("a"): 0,
// 							}),
// 						},
// 					)),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						value.ToSymbol("<methodDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(68, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 6),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								value.ToSymbol("Foo.:foo"),
// 								[]byte{
// 									byte(bytecode.GET_IVAR_0),
// 									byte(bytecode.JUMP_UNLESS_NP), 0, 2,
// 									byte(bytecode.POP),
// 									byte(bytecode.INT_3),
// 									byte(bytecode.DUP),
// 									byte(bytecode.SET_IVAR_0),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(39, 5, 6), P(59, 5, 26)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(5, 9),
// 								},
// 								nil,
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			bytecodeCompilerTest(tc, t)
// 		})
// 	}
// }

func TestGoBitwiseAnd(t *testing.T) {
	tests := goTestTable{
		"resolve static AND": {
			input: "a := 23 & 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(2)).ToValue()
}
`,
		},
		"resolve static nested AND": {
			input: "a := 23 & 15 & 46",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(6)).ToValue()
}
`,
		},

		"and smallint smallint": {
			input: `
				val a = 23
				val b = 10
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.SmallInt // var b: 10
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = value.SmallInt(10)
	l2 = ((l0).BitwiseAndSmallInt(l1)).ToValue()
}
`,
		},
		"and smallint bigint": {
			input: `
				val a = 23
				c := a & 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var c: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).BitwiseAndBigInt(bi0)
}
`,
		},
		"and smallint int": {
			input: `
				val a = 23
				b := 5
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).BitwiseAndInt(l1)
}
`,
		},

		"and bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 10
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.SmallInt // var b: 10
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(10)
	l2 = (l0).BitwiseAndSmallInt(l1)
}
`,
		},
		"and bigint bigint": {
			input: `
				val a = 18446744073709551616
				c := a & 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var c: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (l0).BitwiseAndBigInt(bi0)
}
`,
		},
		"and bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).BitwiseAndInt(l1)
}
`,
		},

		"and int64": {
			input: `
				a := 23i64
				b := 10i64
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(23)
	l1 = value.Int64(10)
	l2 = (l0) & (l1)
}
`,
		},
		"and int32": {
			input: `
				a := 23i32
				b := 10i32
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(23)
	l1 = value.Int32(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and int16": {
			input: `
				a := 23i16
				b := 10i16
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(23)
	l1 = value.Int16(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and int8": {
			input: `
				a := 23i8
				b := 10i8
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(23)
	l1 = value.Int8(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and uint64": {
			input: `
				a := 23u64
				b := 10u64
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(23)
	l1 = value.UInt64(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and uint32": {
			input: `
				a := 23u32
				b := 10u32
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(23)
	l1 = value.UInt32(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and uint16": {
			input: `
				a := 23u16
				b := 10u16
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(23)
	l1 = value.UInt16(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and uint8": {
			input: `
				a := 23u8
				b := 10u8
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(23)
	l1 = value.UInt8(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and uint": {
			input: `
				a := 23u
				b := 10u
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(23)
	l1 = value.UInt(10)
	l2 = (l0) & (l1)
}
`,
		},

		"and ints": {
			input: `
				a := 23
				b := 10
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	l2 = value.BitwiseAndInts(l0, l1)
}
`,
		},
		"and custom object": {
			input: `
				module Foo
					def &(other: Int): Int
						5 & other
					end
				end

				a := Foo
				b := 10
				c := a & b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::&")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__and_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::&, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).BitwiseAndInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	callFrame.SetNativeLineNumber(10)
	t1, err = Foo_ns__and_(thread, l0, l1) // receiver: Foo, name: &
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "&", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__and_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
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

func TestGoBitwiseAndNot(t *testing.T) {
	tests := goTestTable{
		"resolve static AND NOT": {
			input: "a := 23 &~ 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(21)).ToValue()
}
`,
		},
		"resolve static nested AND NOT": {
			input: "a := 23 &~ 15 &~ 46",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(16)).ToValue()
}
`,
		},

		"and not int64": {
			input: `
				a := 23i64
				b := 10i64
				c := a &~ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(23)
	l1 = value.Int64(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not int32": {
			input: `
				a := 23i32
				b := 10i32
				c := a &~ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(23)
	l1 = value.Int32(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not int16": {
			input: `
				a := 23i16
				b := 10i16
				c := a &~ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(23)
	l1 = value.Int16(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not int8": {
			input: `
				a := 23i8
				b := 10i8
				c := a &~ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(23)
	l1 = value.Int8(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not uint64": {
			input: `
				a := 23u64
				b := 10u64
				c := a &~ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(23)
	l1 = value.UInt64(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not uint32": {
			input: `
				a := 23u32
				b := 10u32
				c := a &~ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(23)
	l1 = value.UInt32(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not uint16": {
			input: `
				a := 23u16
				b := 10u16
				c := a &~ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(23)
	l1 = value.UInt16(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not uint8": {
			input: `
				a := 23u8
				b := 10u8
				c := a &~ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(23)
	l1 = value.UInt8(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not uint": {
			input: `
				a := 23u
				b := 10u
				c := a &~ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(23)
	l1 = value.UInt(10)
	l2 = (l0) &^ (l1)
}
`,
		},

		"and not ints": {
			input: `
				a := 23
				b := 10
				c := a &~ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	l2 = value.BitwiseAndNotInts(l0, l1)
}
`,
		},
		"and not custom object": {
			input: `
				module Foo
					def &~(other: Int): Int
						5 &~ other
					end
				end

				a := Foo
				b := 10
				c := a &~ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::&~")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__andnot_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::&~, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.BitwiseAndNotInts((value.SmallInt(5)).ToValue(), l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	callFrame.SetNativeLineNumber(10)
	t1, err = Foo_ns__andnot_(thread, l0, l1) // receiver: Foo, name: &~
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "&~", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__andnot_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
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

func TestGoBitwiseOr(t *testing.T) {
	tests := goTestTable{
		"resolve static OR": {
			input: "a := 23 | 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(31)).ToValue()
}
`,
		},
		"resolve static nested OR": {
			input: "a := 23 | 15 | 46",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(63)).ToValue()
}
`,
		},

		"or int64": {
			input: `
				a := 23i64
				b := 10i64
				c := a | b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(23)
	l1 = value.Int64(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or int32": {
			input: `
				a := 23i32
				b := 10i32
				c := a | b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(23)
	l1 = value.Int32(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or int16": {
			input: `
				a := 23i16
				b := 10i16
				c := a | b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(23)
	l1 = value.Int16(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or int8": {
			input: `
				a := 23i8
				b := 10i8
				c := a | b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(23)
	l1 = value.Int8(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or uint64": {
			input: `
				a := 23u64
				b := 10u64
				c := a | b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(23)
	l1 = value.UInt64(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or uint32": {
			input: `
				a := 23u32
				b := 10u32
				c := a | b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(23)
	l1 = value.UInt32(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or uint16": {
			input: `
				a := 23u16
				b := 10u16
				c := a | b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(23)
	l1 = value.UInt16(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or uint8": {
			input: `
				a := 23u8
				b := 10u8
				c := a | b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(23)
	l1 = value.UInt8(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or uint": {
			input: `
				a := 23u
				b := 10u
				c := a | b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(23)
	l1 = value.UInt(10)
	l2 = (l0) | (l1)
}
`,
		},

		"or ints": {
			input: `
				a := 23
				b := 10
				c := a | b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	l2 = value.BitwiseOrInts(l0, l1)
}
`,
		},
		"or custom object": {
			input: `
				module Foo
					def |(other: Int): Int
						5 | other
					end
				end

				a := Foo
				b := 10
				c := a | b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::|")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__or_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::|, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.BitwiseOrInts((value.SmallInt(5)).ToValue(), l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	callFrame.SetNativeLineNumber(10)
	t1, err = Foo_ns__or_(thread, l0, l1) // receiver: Foo, name: |
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "|", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__or_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
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

func TestGoBitwiseXor(t *testing.T) {
	tests := goTestTable{
		"resolve static XOR": {
			input: "a := 23 ^ 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(29)).ToValue()
}
`,
		},
		"resolve static nested XOR": {
			input: "a := 23 ^ 15 ^ 46",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(54)).ToValue()
}
`,
		},

		"xor int64": {
			input: `
				a := 23i64
				b := 10i64
				c := a ^ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(23)
	l1 = value.Int64(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor int32": {
			input: `
				a := 23i32
				b := 10i32
				c := a ^ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(23)
	l1 = value.Int32(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor int16": {
			input: `
				a := 23i16
				b := 10i16
				c := a ^ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(23)
	l1 = value.Int16(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor int8": {
			input: `
				a := 23i8
				b := 10i8
				c := a ^ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(23)
	l1 = value.Int8(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor uint64": {
			input: `
				a := 23u64
				b := 10u64
				c := a ^ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(23)
	l1 = value.UInt64(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor uint32": {
			input: `
				a := 23u32
				b := 10u32
				c := a ^ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(23)
	l1 = value.UInt32(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor uint16": {
			input: `
				a := 23u16
				b := 10u16
				c := a ^ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(23)
	l1 = value.UInt16(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor uint8": {
			input: `
				a := 23u8
				b := 10u8
				c := a ^ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(23)
	l1 = value.UInt8(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor uint": {
			input: `
				a := 23u
				b := 10u
				c := a ^ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(23)
	l1 = value.UInt(10)
	l2 = (l0) ^ (l1)
}
`,
		},

		"xor ints": {
			input: `
				a := 23
				b := 10
				c := a ^ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	l2 = value.BitwiseXorInts(l0, l1)
}
`,
		},
		"xor custom object": {
			input: `
				module Foo
					def ^(other: Int): Int
						5 ^ other
					end
				end

				a := Foo
				b := 10
				c := a ^ b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::^")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__xor_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::^, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return value.BitwiseXorInts((value.SmallInt(5)).ToValue(), l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	callFrame.SetNativeLineNumber(10)
	t1, err = Foo_ns__xor_(thread, l0, l1) // receiver: Foo, name: ^
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "^", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__xor_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
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

func TestGoModulo(t *testing.T) {
	tests := goTestTable{
		"resolve static modulo": {
			input: "a := 23 % 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(3)).ToValue()
}
`,
		},
		"modulo ints": {
			input: `
				a := 23
				b := a % 10
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.ModuloInts(l0, (value.SmallInt(10)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"modulo int": {
			input: `
				a := 23
				b := a % 10.5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.ModuloInt(l0, (value.Float(10.5)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = (t1).AsFloat()
}
`,
		},

		"modulo smallint smallint": {
			input: `
				val a = 23
				b := a % 10
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	callFrame.SetNativeLineNumber(3)
	t1, err = (l0).ModuloSmallInt(value.SmallInt(10))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"modulo smallint bigint": {
			input: `
				val a = 23
				b := a % 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	callFrame.SetNativeLineNumber(3)
	t1, err = (l0).ModuloBigInt(bi0)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"modulo smallint float": {
			input: `
				val a = 23
				b := a % 2.5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).ModuloFloat(value.Float(2.5))
}
`,
		},
		"modulo smallint bigfloat": {
			input: `
				val a = 23
				b := a % 2.5bf
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).ModuloBigFloat(bf0)
}
`,
		},
		"modulo smallint int": {
			input: `
				val a = 23
				b := 5
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloInt(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"modulo smallint value": {
			input: `
				val a = 23
				var b: Int | Float = 5
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloVal(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 5
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.SmallInt // var b: 5
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloSmallInt(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"modulo bigint float": {
			input: `
				val a = 18446744073709551616
				val b = 5.5
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Float // var b: 5.5
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.Float(5.5)
	l2 = (l0).ModuloFloat(l1)
}
`,
		},
		"modulo bigint bigfloat": {
			input: `
				val a = 18446744073709551616
				val b = 5.5bf
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)
var bf0 = value.ParseBigFloatPanic("5.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigFloat // var b: 5.5bf
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bf0
	l2 = (l0).ModuloBigFloat(l1)
}
`,
		},
		"modulo bigint bigint": {
			input: `
				val a = 18446744073709551616
				val b = 18446744073709551616
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigInt // var b: 18446744073709551616
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bi0
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloBigInt(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"modulo bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloInt(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"modulo bigint value": {
			input: `
				var a: Int | Float = 10
				b := 5
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int | Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(10)).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.ModuloVal(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo int64": {
			input: `
				a := 6i64
				b := 5i64
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var t1 value.Int64
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(6)
	l1 = value.Int64(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloInt64(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo int32": {
			input: `
				a := 6i32
				b := 5i32
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var t1 value.Int32
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(6)
	l1 = value.Int32(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloInt32(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo int16": {
			input: `
				a := 6i16
				b := 5i16
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var t1 value.Int16
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(6)
	l1 = value.Int16(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloInt16(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo int8": {
			input: `
				a := 6i8
				b := 5i8
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var t1 value.Int8
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(6)
	l1 = value.Int8(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloInt8(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo uint": {
			input: `
				a := 6u
				b := 5u
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var t1 value.UInt
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(6)
	l1 = value.UInt(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloUInt(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo uint64": {
			input: `
				a := 6u64
				b := 5u64
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var t1 value.UInt64
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(6)
	l1 = value.UInt64(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloUInt64(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo uint32": {
			input: `
				a := 6u32
				b := 5u32
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var t1 value.UInt32
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(6)
	l1 = value.UInt32(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloUInt32(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo uint16": {
			input: `
				a := 6u16
				b := 5u16
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var t1 value.UInt16
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(6)
	l1 = value.UInt16(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloUInt16(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo uint8": {
			input: `
				a := 6u8
				b := 5u8
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var t1 value.UInt8
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(6)
	l1 = value.UInt8(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloUInt8(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo float smallint": {
			input: `
				a := 2.5
				val b = 20
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.SmallInt // var b: 20
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.SmallInt(20)
	l2 = (l0).ModuloSmallInt(l1)
}
`,
		},
		"modulo float bigint": {
			input: `
				a := 2.5
				c := a % 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var c: Std::Float
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (l0).ModuloBigInt(bi0)
}
`,
		},
		"modulo float float": {
			input: `
				a := 2.5
				b := 0.1
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Float(0.1)
	l2 = (l0).ModuloFloat(l1)
}
`,
		},
		"modulo float bigfloat": {
			input: `
				a := 2.5
				b := 0.1bf
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("0.1")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = bf0
	l2 = (l0).ModuloBigFloat(l1)
}
`,
		},
		"modulo float int": {
			input: `
				a := 2.5
				b := 1
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (value.SmallInt(1)).ToValue()
	l2 = (l0).ModuloInt(l1)
}
`,
		},
		"modulo float value": {
			input: `
				var a: Float | Int = 2.5
				b := 1
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.ModuloVal(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"modulo float64": {
			input: `
				a := 2.5f64
				b := 1f64
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float64 // var a: Std::Float64
	_ = l0
	var l1 value.Float64 // var b: Std::Float64
	_ = l1
	var l2 value.Float64 // var c: Std::Float64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float64(2.5)
	l1 = value.Float64(1)
	l2 = (l0).ModuloFloat64(l1)
}
`,
		},

		"modulo float32": {
			input: `
				a := 2.5f32
				b := 1f32
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float32 // var a: Std::Float32
	_ = l0
	var l1 value.Float32 // var b: Std::Float32
	_ = l1
	var l2 value.Float32 // var c: Std::Float32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float32(2.5)
	l1 = value.Float32(1)
	l2 = (l0).ModuloFloat32(l1)
}
`,
		},

		"modulo bigfloat smallint": {
			input: `
				a := 2.5bf
				val b = 1
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.SmallInt // var b: 1
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.SmallInt(1)
	l2 = (l0).ModuloSmallInt(l1)
}
`,
		},
		"modulo bigfloat bigint": {
			input: `
				a := 2.5bf
				c := a % 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var c: Std::BigFloat
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (l0).ModuloBigInt(bi0)
}
`,
		},
		"modulo bigfloat float": {
			input: `
				a := 2.5bf
				val b = 1.0
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Float // var b: 1.0
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Float(1)
	l2 = (l0).ModuloFloat(l1)
}
`,
		},
		"modulo bigfloat bigfloat": {
			input: `
				a := 2.5bf
				val b = 1.0bf
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")
var bf1 = value.ParseBigFloatPanic("1.0")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var b: 1.0bf
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = bf1
	l2 = (l0).ModuloBigFloat(l1)
}
`,
		},
		"modulo bigfloat int": {
			input: `
				a := 2.5bf
				b := 1
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = (l0).ModuloInt(l1)
}
`,
		},
		"modulo bigfloat value": {
			input: `
				a := 2.5bf
				var b: Int | Float = 1
				c := a % b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).ModuloVal(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = (*value.BigFloat)((t1).Pointer())
}
`,
		},

		"modulo builtin values": {
			input: `
				var a: Float | Int = 2.5
				b := a % 0.1
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::CoercibleNumeric
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.ModuloVal(l0, (value.Float(0.1)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},

		"modulo value": {
			input: `
				module Foo
					def %(other: Int): Int
						5 % other
					end
				end
				a := Foo
				b := a % 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::%")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__mod_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::%, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	callFrame.SetNativeLineNumber(4)
	t1, err = (value.SmallInt(5)).ModuloInt(l0)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	callFrame.SetNativeLineNumber(8)
	t1, err = Foo_ns__mod_(thread, l0, (value.SmallInt(5)).ToValue()) // receiver: Foo, name: %
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "%", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__mod_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"resolve static nested modulo": {
			input: "a := 24 % 15 % 2",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(1)).ToValue()
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

func TestGoDivide(t *testing.T) {
	tests := goTestTable{
		"resolve static divide": {
			input: "a := 23 / 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(2)).ToValue()
}
`,
		},
		"div ints": {
			input: `
				a := 23
				b := a / 10
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.DivideInts(l0, (value.SmallInt(10)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"div int": {
			input: `
				a := 23
				b := a / 10.5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.DivideInt(l0, (value.Float(10.5)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = (t1).AsFloat()
}
`,
		},

		"div smallint smallint": {
			input: `
				val a = 23
				b := a / 10
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	callFrame.SetNativeLineNumber(3)
	t1, err = (l0).DivideSmallInt(value.SmallInt(10))
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"div smallint bigint": {
			input: `
				val a = 23
				b := a / 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	callFrame.SetNativeLineNumber(3)
	t1, err = (l0).DivideBigInt(bi0)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},
		"div smallint float": {
			input: `
				val a = 23
				b := a / 2.5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).DivideFloat(value.Float(2.5))
}
`,
		},
		"div smallint bigfloat": {
			input: `
				val a = 23
				b := a / 2.5bf
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).DivideBigFloat(bf0)
}
`,
		},
		"div smallint int": {
			input: `
				val a = 23
				b := 5
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideInt(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"div smallint value": {
			input: `
				val a = 23
				var b: Int | Float = 5
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideVal(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 5
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.SmallInt // var b: 5
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideSmallInt(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"div bigint float": {
			input: `
				val a = 18446744073709551616
				val b = 5.5
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Float // var b: 5.5
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.Float(5.5)
	l2 = (l0).DivideFloat(l1)
}
`,
		},
		"div bigint bigfloat": {
			input: `
				val a = 18446744073709551616
				val b = 5.5bf
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)
var bf0 = value.ParseBigFloatPanic("5.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigFloat // var b: 5.5bf
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bf0
	l2 = (l0).DivideBigFloat(l1)
}
`,
		},
		"div bigint bigint": {
			input: `
				val a = 18446744073709551616
				val b = 18446744073709551616
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigInt // var b: 18446744073709551616
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bi0
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideBigInt(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"div bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideInt(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},
		"div bigint value": {
			input: `
				var a: Int | Float = 10
				b := 5
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int | Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(10)).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.DivideVal(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div int64": {
			input: `
				a := 6i64
				b := 5i64
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var t1 value.Int64
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(6)
	l1 = value.Int64(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideInt64(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div int32": {
			input: `
				a := 6i32
				b := 5i32
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var t1 value.Int32
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(6)
	l1 = value.Int32(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideInt32(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div int16": {
			input: `
				a := 6i16
				b := 5i16
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var t1 value.Int16
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(6)
	l1 = value.Int16(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideInt16(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div int8": {
			input: `
				a := 6i8
				b := 5i8
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var t1 value.Int8
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(6)
	l1 = value.Int8(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideInt8(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div uint": {
			input: `
				a := 6u
				b := 5u
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var t1 value.UInt
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(6)
	l1 = value.UInt(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideUInt(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div uint64": {
			input: `
				a := 6u64
				b := 5u64
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var t1 value.UInt64
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(6)
	l1 = value.UInt64(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideUInt64(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div uint32": {
			input: `
				a := 6u32
				b := 5u32
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var t1 value.UInt32
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(6)
	l1 = value.UInt32(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideUInt32(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div uint16": {
			input: `
				a := 6u16
				b := 5u16
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var t1 value.UInt16
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(6)
	l1 = value.UInt16(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideUInt16(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div uint8": {
			input: `
				a := 6u8
				b := 5u8
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var t1 value.UInt8
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(6)
	l1 = value.UInt8(5)
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideUInt8(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div float smallint": {
			input: `
				a := 2.5
				val b = 20
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.SmallInt // var b: 20
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.SmallInt(20)
	l2 = (l0).DivideSmallInt(l1)
}
`,
		},
		"div float bigint": {
			input: `
				a := 2.5
				c := a / 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var c: Std::Float
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (l0).DivideBigInt(bi0)
}
`,
		},
		"div float float": {
			input: `
				a := 2.5
				b := 0.1
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Float(0.1)
	l2 = (l0).DivideFloat(l1)
}
`,
		},
		"div float bigfloat": {
			input: `
				a := 2.5
				b := 0.1bf
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("0.1")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = bf0
	l2 = (l0).DivideBigFloat(l1)
}
`,
		},
		"div float int": {
			input: `
				a := 2.5
				b := 1
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (value.SmallInt(1)).ToValue()
	l2 = (l0).DivideInt(l1)
}
`,
		},
		"div float value": {
			input: `
				var a: Float | Int = 2.5
				b := 1
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.DivideVal(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"div float64": {
			input: `
				a := 2.5f64
				b := 1f64
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float64 // var a: Std::Float64
	_ = l0
	var l1 value.Float64 // var b: Std::Float64
	_ = l1
	var l2 value.Float64 // var c: Std::Float64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float64(2.5)
	l1 = value.Float64(1)
	l2 = (l0) / (l1)
}
`,
		},

		"div float32": {
			input: `
				a := 2.5f32
				b := 1f32
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float32 // var a: Std::Float32
	_ = l0
	var l1 value.Float32 // var b: Std::Float32
	_ = l1
	var l2 value.Float32 // var c: Std::Float32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float32(2.5)
	l1 = value.Float32(1)
	l2 = (l0) / (l1)
}
`,
		},

		"div bigfloat smallint": {
			input: `
				a := 2.5bf
				val b = 1
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.SmallInt // var b: 1
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.SmallInt(1)
	l2 = (l0).DivideSmallInt(l1)
}
`,
		},

		"div bigfloat bigint": {
			input: `
				a := 2.5bf
				c := a / 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var c: Std::BigFloat
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (l0).DivideBigInt(bi0)
}
`,
		},
		"div bigfloat float": {
			input: `
				a := 2.5bf
				val b = 1.0
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Float // var b: 1.0
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Float(1)
	l2 = (l0).DivideFloat(l1)
}
`,
		},
		"div bigfloat bigfloat": {
			input: `
				a := 2.5bf
				val b = 1.0bf
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")
var bf1 = value.ParseBigFloatPanic("1.0")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var b: 1.0bf
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = bf1
	l2 = (l0).DivideBigFloat(l1)
}
`,
		},
		"div bigfloat int": {
			input: `
				a := 2.5bf
				b := 1
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = (l0).DivideInt(l1)
}
`,
		},
		"div bigfloat value": {
			input: `
				a := 2.5bf
				var b: Int | Float = 1
				c := a / b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).DivideVal(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = (*value.BigFloat)((t1).Pointer())
}
`,
		},

		"div builtin values": {
			input: `
				var a: Float | Int = 2.5
				b := a / 0.1
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::CoercibleNumeric
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.DivideVal(l0, (value.Float(0.1)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},

		"div value": {
			input: `
				module Foo
					def /(other: Int): Int
						5 / other
					end
				end
				a := Foo
				b := a / 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::/")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__div_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::/, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	callFrame.SetNativeLineNumber(4)
	t1, err = (value.SmallInt(5)).DivideInt(l0)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	callFrame.SetNativeLineNumber(8)
	t1, err = Foo_ns__div_(thread, l0, (value.SmallInt(5)).ToValue()) // receiver: Foo, name: /
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "/", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__div_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"resolve static nested div": {
			input: "a := 90 / 15 / 2",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(3)).ToValue()
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

func TestGoAdd(t *testing.T) {
	tests := goTestTable{
		"resolve static add": {
			input: "a := 23 + 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(33)).ToValue()
}
`,
		},
		"add ints": {
			input: `
						a := 23
						b := a + 10
					`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = value.AddInts(l0, (value.SmallInt(10)).ToValue())
}
`,
		},
		"add int": {
			input: `
						a := 23
						b := a + 10.5
					`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.AddInt(l0, (value.Float(10.5)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = (t1).AsFloat()
}
`,
		},

		"add smallint smallint": {
			input: `
						val a = 23
						b := a + 10
					`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).AddSmallInt(value.SmallInt(10))
}
`,
		},
		"add smallint bigint": {
			input: `
				val a = 23
				b := a + 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).AddBigInt(bi0)
}
`,
		},
		"add smallint float": {
			input: `
						val a = 23
						b := a + 2.5
					`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).AddFloat(value.Float(2.5))
}
`,
		},
		"add smallint bigfloat": {
			input: `
						val a = 23
						b := a + 2.5bf
					`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).AddBigFloat(bf0)
}
`,
		},
		"add smallint int": {
			input: `
						val a = 23
						b := 5
						c := a + b
					`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).AddInt(l1)
}
`,
		},
		"add smallint value": {
			input: `
						val a = 23
						var b: Int | Float = 5
						c := a + b
					`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).AddVal(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"add bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 5
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.SmallInt // var b: 5
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(5)
	l2 = (l0).AddSmallInt(l1)
}
`,
		},
		"add bigint float": {
			input: `
				val a = 18446744073709551616
				val b = 5.5
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Float // var b: 5.5
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.Float(5.5)
	l2 = (l0).AddFloat(l1)
}
`,
		},
		"add bigint bigfloat": {
			input: `
				val a = 18446744073709551616
				val b = 5.5bf
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)
var bf0 = value.ParseBigFloatPanic("5.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigFloat // var b: 5.5bf
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bf0
	l2 = (l0).AddBigFloat(l1)
}
`,
		},
		"add bigint bigint": {
			input: `
				val a = 18446744073709551616
				val b = 18446744073709551616
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigInt // var b: 18446744073709551616
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bi0
	l2 = (l0).AddBigInt(l1)
}
`,
		},
		"add bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).AddInt(l1)
}
`,
		},
		"add bigint value": {
			input: `
				var a: Int | Float = 10
				b := 5
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int | Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(10)).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.AddVal(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"add int64": {
			input: `
				a := 6i64
				b := 5i64
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(6)
	l1 = value.Int64(5)
	l2 = (l0) + (l1)
}
`,
		},

		"add int32": {
			input: `
				a := 6i32
				b := 5i32
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(6)
	l1 = value.Int32(5)
	l2 = (l0) + (l1)
}
`,
		},

		"add int16": {
			input: `
				a := 6i16
				b := 5i16
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(6)
	l1 = value.Int16(5)
	l2 = (l0) + (l1)
}
`,
		},

		"add int8": {
			input: `
				a := 6i8
				b := 5i8
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(6)
	l1 = value.Int8(5)
	l2 = (l0) + (l1)
}
`,
		},

		"add uint": {
			input: `
				a := 6u
				b := 5u
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(6)
	l1 = value.UInt(5)
	l2 = (l0) + (l1)
}
`,
		},

		"add uint64": {
			input: `
				a := 6u64
				b := 5u64
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(6)
	l1 = value.UInt64(5)
	l2 = (l0) + (l1)
}
`,
		},

		"add uint32": {
			input: `
				a := 6u32
				b := 5u32
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(6)
	l1 = value.UInt32(5)
	l2 = (l0) + (l1)
}
`,
		},

		"add uint16": {
			input: `
				a := 6u16
				b := 5u16
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(6)
	l1 = value.UInt16(5)
	l2 = (l0) + (l1)
}
`,
		},

		"add uint8": {
			input: `
				a := 6u8
				b := 5u8
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(6)
	l1 = value.UInt8(5)
	l2 = (l0) + (l1)
}
`,
		},

		"add float smallint": {
			input: `
				a := 2.5
				val b = 20
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.SmallInt // var b: 20
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.SmallInt(20)
	l2 = (l0).AddSmallInt(l1)
}
`,
		},
		"add float bigint": {
			input: `
				a := 2.5
				c := a + 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var c: Std::Float
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (l0).AddBigInt(bi0)
}
`,
		},
		"add float float": {
			input: `
				a := 2.5
				b := 0.1
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Float(0.1)
	l2 = (l0).AddFloat(l1)
}
`,
		},
		"add float bigfloat": {
			input: `
				a := 2.5
				b := 0.1bf
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("0.1")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = bf0
	l2 = (l0).AddBigFloat(l1)
}
`,
		},
		"add float int": {
			input: `
				a := 2.5
				b := 1
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (value.SmallInt(1)).ToValue()
	l2 = (l0).AddInt(l1)
}
`,
		},
		"add float value": {
			input: `
				var a: Float | Int = 2.5
				b := 1
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.AddVal(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"add float64": {
			input: `
				a := 2.5f64
				b := 1f64
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float64 // var a: Std::Float64
	_ = l0
	var l1 value.Float64 // var b: Std::Float64
	_ = l1
	var l2 value.Float64 // var c: Std::Float64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float64(2.5)
	l1 = value.Float64(1)
	l2 = (l0) + (l1)
}
`,
		},

		"add float32": {
			input: `
				a := 2.5f32
				b := 1f32
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float32 // var a: Std::Float32
	_ = l0
	var l1 value.Float32 // var b: Std::Float32
	_ = l1
	var l2 value.Float32 // var c: Std::Float32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float32(2.5)
	l1 = value.Float32(1)
	l2 = (l0) + (l1)
}
`,
		},

		"add bigfloat smallint": {
			input: `
				a := 2.5bf
				val b = 1
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.SmallInt // var b: 1
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.SmallInt(1)
	l2 = (l0).AddSmallInt(l1)
}
`,
		},

		"add bigfloat bigint": {
			input: `
				a := 2.5bf
				c := a + 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var c: Std::BigFloat
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (l0).AddBigInt(bi0)
}
`,
		},
		"add bigfloat float": {
			input: `
				a := 2.5bf
				val b = 1.0
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Float // var b: 1.0
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Float(1)
	l2 = (l0).AddFloat(l1)
}
`,
		},
		"add bigfloat bigfloat": {
			input: `
				a := 2.5bf
				val b = 1.0bf
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")
var bf1 = value.ParseBigFloatPanic("1.0")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var b: 1.0bf
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = bf1
	l2 = (l0).AddBigFloat(l1)
}
`,
		},
		"add bigfloat int": {
			input: `
				a := 2.5bf
				b := 1
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = (l0).AddInt(l1)
}
`,
		},
		"add bigfloat value": {
			input: `
				a := 2.5bf
				var b: Int | Float = 1
				c := a + b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).AddVal(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = (*value.BigFloat)((t1).Pointer())
}
`,
		},

		"add builtin values": {
			input: `
				var a: Float | Int = 2.5
				b := a + 0.1
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::CoercibleNumeric
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.AddVal(l0, (value.Float(0.1)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},

		"add value": {
			input: `
				module Foo
					def +(other: Int): Int
						5 + other
					end
				end
				a := Foo
				b := a + 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::+")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__add_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::+, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).AddInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	callFrame.SetNativeLineNumber(8)
	t1, err = Foo_ns__add_(thread, l0, (value.SmallInt(5)).ToValue()) // receiver: Foo, name: +
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "+", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__add_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"resolve static nested add": {
			input: "a := 90 + 15 + 2",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(107)).ToValue()
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

func TestGoSubtract(t *testing.T) {
	tests := goTestTable{
		"resolve static subtract": {
			input: "a := 23 - 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(13)).ToValue()
}
`,
		},
		"subtract ints": {
			input: `
				a := 23
				b := a - 10
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = value.SubtractInts(l0, (value.SmallInt(10)).ToValue())
}
`,
		},
		"subtract int": {
			input: `
						a := 23
						b := a - 10.5
					`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.SubtractInt(l0, (value.Float(10.5)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = (t1).AsFloat()
}
`,
		},

		"subtract smallint smallint": {
			input: `
						val a = 23
						b := a - 10
					`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).SubtractSmallInt(value.SmallInt(10))
}
`,
		},
		"subtract smallint bigint": {
			input: `
				val a = 23
				b := a - 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).SubtractBigInt(bi0)
}
`,
		},
		"subtract smallint float": {
			input: `
						val a = 23
						b := a - 2.5
					`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).SubtractFloat(value.Float(2.5))
}
`,
		},
		"subtract smallint bigfloat": {
			input: `
				val a = 23
				b := a - 2.5bf
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).SubtractBigFloat(bf0)
}
`,
		},
		"subtract smallint int": {
			input: `
				val a = 23
				b := 5
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).SubtractInt(l1)
}
`,
		},
		"subtract smallint value": {
			input: `
				val a = 23
				var b: Int | Float = 5
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).SubtractVal(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"subtract bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 5
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.SmallInt // var b: 5
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(5)
	l2 = (l0).SubtractSmallInt(l1)
}
`,
		},
		"subtract bigint float": {
			input: `
				val a = 18446744073709551616
				val b = 5.5
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Float // var b: 5.5
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.Float(5.5)
	l2 = (l0).SubtractFloat(l1)
}
`,
		},
		"subtract bigint bigfloat": {
			input: `
				val a = 18446744073709551616
				val b = 5.5bf
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)
var bf0 = value.ParseBigFloatPanic("5.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigFloat // var b: 5.5bf
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bf0
	l2 = (l0).SubtractBigFloat(l1)
}
`,
		},
		"subtract bigint bigint": {
			input: `
				val a = 18446744073709551616
				val b = 18446744073709551616
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigInt // var b: 18446744073709551616
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bi0
	l2 = (l0).SubtractBigInt(l1)
}
`,
		},
		"subtract bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).SubtractInt(l1)
}
`,
		},
		"subtract bigint value": {
			input: `
				var a: Int | Float = 10
				b := 5
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int | Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(10)).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.SubtractVal(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"subtract int64": {
			input: `
				a := 6i64
				b := 5i64
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(6)
	l1 = value.Int64(5)
	l2 = (l0) - (l1)
}
`,
		},

		"subtract int32": {
			input: `
				a := 6i32
				b := 5i32
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(6)
	l1 = value.Int32(5)
	l2 = (l0) - (l1)
}
`,
		},

		"subtract int16": {
			input: `
				a := 6i16
				b := 5i16
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(6)
	l1 = value.Int16(5)
	l2 = (l0) - (l1)
}
`,
		},

		"subtract int8": {
			input: `
				a := 6i8
				b := 5i8
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(6)
	l1 = value.Int8(5)
	l2 = (l0) - (l1)
}
`,
		},

		"subtract uint": {
			input: `
				a := 6u
				b := 5u
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(6)
	l1 = value.UInt(5)
	l2 = (l0) - (l1)
}
`,
		},

		"subtract uint64": {
			input: `
				a := 6u64
				b := 5u64
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(6)
	l1 = value.UInt64(5)
	l2 = (l0) - (l1)
}
`,
		},

		"subtract uint32": {
			input: `
				a := 6u32
				b := 5u32
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(6)
	l1 = value.UInt32(5)
	l2 = (l0) - (l1)
}
`,
		},

		"subtract uint16": {
			input: `
				a := 6u16
				b := 5u16
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(6)
	l1 = value.UInt16(5)
	l2 = (l0) - (l1)
}
`,
		},

		"subtract uint8": {
			input: `
				a := 6u8
				b := 5u8
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(6)
	l1 = value.UInt8(5)
	l2 = (l0) - (l1)
}
`,
		},

		"subtract float smallint": {
			input: `
				a := 2.5
				val b = 20
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.SmallInt // var b: 20
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.SmallInt(20)
	l2 = (l0).SubtractSmallInt(l1)
}
`,
		},
		"subtract float bigint": {
			input: `
				a := 2.5
				c := a - 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var c: Std::Float
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (l0).SubtractBigInt(bi0)
}
`,
		},
		"subtract float float": {
			input: `
				a := 2.5
				b := 0.1
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Float(0.1)
	l2 = (l0).SubtractFloat(l1)
}
`,
		},
		"subtract float bigfloat": {
			input: `
				a := 2.5
				b := 0.1bf
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("0.1")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = bf0
	l2 = (l0).SubtractBigFloat(l1)
}
`,
		},
		"subtract float int": {
			input: `
				a := 2.5
				b := 1
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (value.SmallInt(1)).ToValue()
	l2 = (l0).SubtractInt(l1)
}
`,
		},
		"subtract float value": {
			input: `
				var a: Float | Int = 2.5
				b := 1
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.SubtractVal(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"subtract float64": {
			input: `
				a := 2.5f64
				b := 1f64
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float64 // var a: Std::Float64
	_ = l0
	var l1 value.Float64 // var b: Std::Float64
	_ = l1
	var l2 value.Float64 // var c: Std::Float64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float64(2.5)
	l1 = value.Float64(1)
	l2 = (l0) - (l1)
}
`,
		},

		"subtract float32": {
			input: `
				a := 2.5f32
				b := 1f32
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float32 // var a: Std::Float32
	_ = l0
	var l1 value.Float32 // var b: Std::Float32
	_ = l1
	var l2 value.Float32 // var c: Std::Float32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float32(2.5)
	l1 = value.Float32(1)
	l2 = (l0) - (l1)
}
`,
		},

		"subtract bigfloat smallint": {
			input: `
				a := 2.5bf
				val b = 1
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.SmallInt // var b: 1
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.SmallInt(1)
	l2 = (l0).SubtractSmallInt(l1)
}
`,
		},

		"subtract bigfloat bigint": {
			input: `
				a := 2.5bf
				c := a - 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var c: Std::BigFloat
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (l0).SubtractBigInt(bi0)
}
`,
		},
		"subtract bigfloat float": {
			input: `
				a := 2.5bf
				val b = 1.0
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Float // var b: 1.0
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Float(1)
	l2 = (l0).SubtractFloat(l1)
}
`,
		},
		"subtract bigfloat bigfloat": {
			input: `
				a := 2.5bf
				val b = 1.0bf
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")
var bf1 = value.ParseBigFloatPanic("1.0")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var b: 1.0bf
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = bf1
	l2 = (l0).SubtractBigFloat(l1)
}
`,
		},
		"subtract bigfloat int": {
			input: `
				a := 2.5bf
				b := 1
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = (l0).SubtractInt(l1)
}
`,
		},
		"subtract bigfloat value": {
			input: `
				a := 2.5bf
				var b: Int | Float = 1
				c := a - b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).SubtractVal(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = (*value.BigFloat)((t1).Pointer())
}
`,
		},

		"subtract builtin values": {
			input: `
				var a: Float | Int = 2.5
				b := a - 0.1
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::CoercibleNumeric
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.SubtractVal(l0, (value.Float(0.1)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},

		"subtract value": {
			input: `
				module Foo
					def -(other: Int): Int
						5 - other
					end
				end
				a := Foo
				b := a - 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::-")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__sub_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::-, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).SubtractInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	callFrame.SetNativeLineNumber(8)
	t1, err = Foo_ns__sub_(thread, l0, (value.SmallInt(5)).ToValue()) // receiver: Foo, name: -
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "-", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__sub_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"resolve static nested subtract": {
			input: "a := 90 - 15 - 2",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(73)).ToValue()
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

func TestGoMultiply(t *testing.T) {
	tests := goTestTable{
		"resolve static multiply": {
			input: "a := 23 * 10",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(230)).ToValue()
}
`,
		},
		"multiply ints": {
			input: `
				a := 23
				b := a * 10
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	l1 = value.MultiplyInts(l0, (value.SmallInt(10)).ToValue())
}
`,
		},
		"multiply int": {
			input: `
				a := 23
				b := a * 10.5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(23)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.MultiplyInt(l0, (value.Float(10.5)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = (t1).AsFloat()
}
`,
		},

		"multiply smallint smallint": {
			input: `
						val a = 23
						b := a * 10
					`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).MultiplySmallInt(value.SmallInt(10))
}
`,
		},
		"multiply smallint bigint": {
			input: `
				val a = 23
				b := a * 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).MultiplyBigInt(bi0)
}
`,
		},
		"multiply smallint float": {
			input: `
				val a = 23
				b := a * 2.5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).MultiplyFloat(value.Float(2.5))
}
`,
		},
		"multiply smallint bigfloat": {
			input: `
				val a = 23
				b := a * 2.5bf
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (l0).MultiplyBigFloat(bf0)
}
`,
		},
		"multiply smallint int": {
			input: `
				val a = 23
				b := 5
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).MultiplyInt(l1)
}
`,
		},
		"multiply smallint value": {
			input: `
				val a = 23
				var b: Int | Float = 5
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 23
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(23)
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).MultiplyVal(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"multiply bigint smallint": {
			input: `
				val a = 18446744073709551616
				val b = 5
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.SmallInt // var b: 5
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.SmallInt(5)
	l2 = (l0).MultiplySmallInt(l1)
}
`,
		},
		"multiply bigint float": {
			input: `
				val a = 18446744073709551616
				val b = 5.5
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Float // var b: 5.5
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = value.Float(5.5)
	l2 = (l0).MultiplyFloat(l1)
}
`,
		},
		"multiply bigint bigfloat": {
			input: `
				val a = 18446744073709551616
				val b = 5.5bf
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)
var bf0 = value.ParseBigFloatPanic("5.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigFloat // var b: 5.5bf
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bf0
	l2 = (l0).MultiplyBigFloat(l1)
}
`,
		},
		"multiply bigint bigint": {
			input: `
				val a = 18446744073709551616
				val b = 18446744073709551616
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 *value.BigInt // var b: 18446744073709551616
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = bi0
	l2 = (l0).MultiplyBigInt(l1)
}
`,
		},
		"multiply bigint int": {
			input: `
				val a = 18446744073709551616
				b := 5
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigInt // var a: 18446744073709551616
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bi0
	l1 = (value.SmallInt(5)).ToValue()
	l2 = (l0).MultiplyInt(l1)
}
`,
		},
		"multiply bigint value": {
			input: `
				var a: Int | Float = 10
				b := 5
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int | Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(10)).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.MultiplyVal(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"multiply int64": {
			input: `
				a := 6i64
				b := 5i64
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int64 // var a: Std::Int64
	_ = l0
	var l1 value.Int64 // var b: Std::Int64
	_ = l1
	var l2 value.Int64 // var c: Std::Int64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int64(6)
	l1 = value.Int64(5)
	l2 = (l0) * (l1)
}
`,
		},

		"multiply int32": {
			input: `
				a := 6i32
				b := 5i32
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int32 // var a: Std::Int32
	_ = l0
	var l1 value.Int32 // var b: Std::Int32
	_ = l1
	var l2 value.Int32 // var c: Std::Int32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int32(6)
	l1 = value.Int32(5)
	l2 = (l0) * (l1)
}
`,
		},

		"multiply int16": {
			input: `
				a := 6i16
				b := 5i16
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int16 // var a: Std::Int16
	_ = l0
	var l1 value.Int16 // var b: Std::Int16
	_ = l1
	var l2 value.Int16 // var c: Std::Int16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int16(6)
	l1 = value.Int16(5)
	l2 = (l0) * (l1)
}
`,
		},

		"multiply int8": {
			input: `
				a := 6i8
				b := 5i8
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Int8 // var a: Std::Int8
	_ = l0
	var l1 value.Int8 // var b: Std::Int8
	_ = l1
	var l2 value.Int8 // var c: Std::Int8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Int8(6)
	l1 = value.Int8(5)
	l2 = (l0) * (l1)
}
`,
		},

		"multiply uint": {
			input: `
				a := 6u
				b := 5u
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt // var a: Std::UInt
	_ = l0
	var l1 value.UInt // var b: Std::UInt
	_ = l1
	var l2 value.UInt // var c: Std::UInt
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt(6)
	l1 = value.UInt(5)
	l2 = (l0) * (l1)
}
`,
		},

		"multiply uint64": {
			input: `
				a := 6u64
				b := 5u64
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt64 // var a: Std::UInt64
	_ = l0
	var l1 value.UInt64 // var b: Std::UInt64
	_ = l1
	var l2 value.UInt64 // var c: Std::UInt64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt64(6)
	l1 = value.UInt64(5)
	l2 = (l0) * (l1)
}
`,
		},

		"multiply uint32": {
			input: `
				a := 6u32
				b := 5u32
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt32 // var a: Std::UInt32
	_ = l0
	var l1 value.UInt32 // var b: Std::UInt32
	_ = l1
	var l2 value.UInt32 // var c: Std::UInt32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt32(6)
	l1 = value.UInt32(5)
	l2 = (l0) * (l1)
}
`,
		},

		"multiply uint16": {
			input: `
				a := 6u16
				b := 5u16
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt16 // var a: Std::UInt16
	_ = l0
	var l1 value.UInt16 // var b: Std::UInt16
	_ = l1
	var l2 value.UInt16 // var c: Std::UInt16
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt16(6)
	l1 = value.UInt16(5)
	l2 = (l0) * (l1)
}
`,
		},

		"multiply uint8": {
			input: `
				a := 6u8
				b := 5u8
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.UInt8 // var a: Std::UInt8
	_ = l0
	var l1 value.UInt8 // var b: Std::UInt8
	_ = l1
	var l2 value.UInt8 // var c: Std::UInt8
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.UInt8(6)
	l1 = value.UInt8(5)
	l2 = (l0) * (l1)
}
`,
		},

		"multiply float smallint": {
			input: `
				a := 2.5
				val b = 20
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.SmallInt // var b: 20
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.SmallInt(20)
	l2 = (l0).MultiplySmallInt(l1)
}
`,
		},
		"multiply float bigint": {
			input: `
				a := 2.5
				c := a * 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var c: Std::Float
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (l0).MultiplyBigInt(bi0)
}
`,
		},
		"multiply float float": {
			input: `
				a := 2.5
				b := 0.1
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Float // var b: Std::Float
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = value.Float(0.1)
	l2 = (l0).MultiplyFloat(l1)
}
`,
		},
		"multiply float bigfloat": {
			input: `
				a := 2.5
				b := 0.1bf
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("0.1")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 *value.BigFloat // var b: Std::BigFloat
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = bf0
	l2 = (l0).MultiplyBigFloat(l1)
}
`,
		},
		"multiply float int": {
			input: `
				a := 2.5
				b := 1
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Float // var c: Std::Float
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(2.5)
	l1 = (value.SmallInt(1)).ToValue()
	l2 = (l0).MultiplyInt(l1)
}
`,
		},
		"multiply float value": {
			input: `
				var a: Float | Int = 2.5
				b := 1
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var c: Std::CoercibleNumeric
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.MultiplyVal(l0, l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = t1
}
`,
		},

		"multiply float64": {
			input: `
				a := 2.5f64
				b := 1f64
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float64 // var a: Std::Float64
	_ = l0
	var l1 value.Float64 // var b: Std::Float64
	_ = l1
	var l2 value.Float64 // var c: Std::Float64
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float64(2.5)
	l1 = value.Float64(1)
	l2 = (l0) * (l1)
}
`,
		},

		"multiply float32": {
			input: `
				a := 2.5f32
				b := 1f32
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float32 // var a: Std::Float32
	_ = l0
	var l1 value.Float32 // var b: Std::Float32
	_ = l1
	var l2 value.Float32 // var c: Std::Float32
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float32(2.5)
	l1 = value.Float32(1)
	l2 = (l0) * (l1)
}
`,
		},

		"multiply bigfloat smallint": {
			input: `
				a := 2.5bf
				val b = 1
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.SmallInt // var b: 1
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.SmallInt(1)
	l2 = (l0).MultiplySmallInt(l1)
}
`,
		},

		"multiply bigfloat bigint": {
			input: `
				a := 2.5bf
				c := a * 18446744073709551616
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")
var bi0 = value.ParseBigIntPanic("18446744073709551616", 0)

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var c: Std::BigFloat
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (l0).MultiplyBigInt(bi0)
}
`,
		},
		"multiply bigfloat float": {
			input: `
				a := 2.5bf
				val b = 1.0
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Float // var b: 1.0
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = value.Float(1)
	l2 = (l0).MultiplyFloat(l1)
}
`,
		},
		"multiply bigfloat bigfloat": {
			input: `
				a := 2.5bf
				val b = 1.0bf
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")
var bf1 = value.ParseBigFloatPanic("1.0")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 *value.BigFloat // var b: 1.0bf
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = bf1
	l2 = (l0).MultiplyBigFloat(l1)
}
`,
		},
		"multiply bigfloat int": {
			input: `
				a := 2.5bf
				b := 1
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	l2 = (l0).MultiplyInt(l1)
}
`,
		},
		"multiply bigfloat value": {
			input: `
				a := 2.5bf
				var b: Int | Float = 1
				c := a * b
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var bf0 = value.ParseBigFloatPanic("2.5")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 *value.BigFloat // var a: Std::BigFloat
	_ = l0
	var l1 value.Value // var b: Std::Int | Std::Float
	_ = l1
	var l2 *value.BigFloat // var c: Std::BigFloat
	_ = l2
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = bf0
	l1 = (value.SmallInt(1)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = (l0).MultiplyVal(l1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l2 = (*value.BigFloat)((t1).Pointer())
}
`,
		},

		"multiply builtin values": {
			input: `
				var a: Float | Int = 2.5
				b := a * 0.1
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Float | Std::Int
	_ = l0
	var l1 value.Value // var b: Std::CoercibleNumeric
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(2.5)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = value.MultiplyVal(l0, (value.Float(0.1)).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t1
}
`,
		},

		"multiply value": {
			input: `
				module Foo
					def *(other: Int): Int
						5 * other
					end
				end
				a := Foo
				b := a * 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var Foo *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::*")
var sym2 = value.ToSymbol("<main>")

func Foo_ns__mul_(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Foo::*, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 3)
	defer thread.PopNativeCallFrame()
	return (value.SmallInt(5)).MultiplyInt(l0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (Foo).ToValue()
	callFrame.SetNativeLineNumber(8)
	t1, err = Foo_ns__mul_(thread, l0, (value.SmallInt(5)).ToValue()) // receiver: Foo, name: *
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
	Foo = value.NewModule()
	namespace = value.Ref(Foo)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (Foo).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "*", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := Foo_ns__mul_(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},

		"resolve static nested multiply": {
			input: "a := 90 * 15 * 2",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(2700)).ToValue()
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
