package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func TestBinaryExpressions(t *testing.T) {
	tests := testTable{
		"is a": {
			input: "3 <: ::Std::Int",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_3),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.IS_A),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.ToSymbol("Std::Int").ToValue(),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(0, 1, 1), P(0, 1, 1)), "this \"is a\" check is always true, `3` will always be an instance of `Std::Int`"),
			},
		},
		"instance of": {
			input: "3 <<: ::Std::Int",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_3),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.INSTANCE_OF),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.ToSymbol("Std::Int").ToValue(),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(0, 1, 1), P(0, 1, 1)), "this \"instance of\" check is always true, `3` will always be an instance of `Std::Int`"),
			},
		},
		"resolve static add": {
			input: "1i8 + 5i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT8), 6,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"add int": {
			input: "a := 1; a + 5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.ADD_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{},
			),
		},
		"add float": {
			input: "a := 1.2; a + 5.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.ADD_FLOAT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.Float(1.2).ToValue(),
					value.Float(5.0).ToValue(),
				},
			),
		},
		"add builtin": {
			input: "a := 1i8; a + 5i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT8), 5,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"add value": {
			input: `
				module Foo
					def +(other: Int): Int
					  other
					end
				end
				Foo + 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INT_5),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(7, 6),
				},
				[]value.Value{
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							namespaceDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.DEF_NAMESPACE), 0,
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(85, 7, 12)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Root").ToValue(),
								value.ToSymbol("Foo").ToValue(),
							},
						),
					),
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							methodDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.GET_SINGLETON),
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.LOAD_VALUE_2),
								byte(bytecode.DEF_METHOD),
								byte(bytecode.POP),
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(85, 7, 12)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpAdd,
										[]byte{
											byte(bytecode.GET_LOCAL_1),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(64, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol("+").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpAdd, 1)),
				},
			),
		},

		"resolve static subtract": {
			input: "151i32 - 25i32 - 5i32",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT32_8), 0x79,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 1, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"subtract int": {
			input: "a := 1; a - 5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.SUBTRACT_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{},
			),
		},
		"subtract float": {
			input: "a := 1.2; a - 5.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.SUBTRACT_FLOAT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.Float(1.2).ToValue(),
					value.Float(5.0).ToValue(),
				},
			),
		},
		"subtract builtin": {
			input: "a := 151i32; a - 25i32 - 5i32",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT32_8), 25,
					byte(bytecode.SUBTRACT),
					byte(bytecode.LOAD_INT32_8), 5,
					byte(bytecode.SUBTRACT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 1, 29)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 12),
				},
				[]value.Value{
					value.Int32(151).ToValue(),
				},
			),
		},
		"subtract value": {
			input: `
				module Foo
					def -(other: Int): Int
					  other
					end
				end
				Foo - 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INT_5),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(7, 6),
				},
				[]value.Value{
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							namespaceDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.DEF_NAMESPACE), 0,
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(85, 7, 12)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Root").ToValue(),
								value.ToSymbol("Foo").ToValue(),
							},
						),
					),
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							methodDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.GET_SINGLETON),
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.LOAD_VALUE_2),
								byte(bytecode.DEF_METHOD),
								byte(bytecode.POP),
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(85, 7, 12)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpSubtract,
										[]byte{
											byte(bytecode.GET_LOCAL_1),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(64, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol("-").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpSubtract, 1)),
				},
			),
		},

		"resolve static multiply": {
			input: "45.5 * 2.5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Float(113.75).ToValue(),
				},
			),
		},
		"multiply int": {
			input: "a := 1; a * 5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.MULTIPLY_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{},
			),
		},
		"multiply float": {
			input: "a := 1.2; a * 5.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.MULTIPLY_FLOAT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.Float(1.2).ToValue(),
					value.Float(5.0).ToValue(),
				},
			),
		},
		"multiply builtin": {
			input: "a := 45i8; a * 2i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT8), 45,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT8), 2,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"multiply value": {
			input: `
				module Foo
					def *(other: Int): Int
					  other
					end
				end
				Foo * 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INT_5),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(7, 6),
				},
				[]value.Value{
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							namespaceDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.DEF_NAMESPACE), 0,
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(85, 7, 12)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Root").ToValue(),
								value.ToSymbol("Foo").ToValue(),
							},
						),
					),
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							methodDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.GET_SINGLETON),
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.LOAD_VALUE_2),
								byte(bytecode.DEF_METHOD),
								byte(bytecode.POP),
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(85, 7, 12)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpMultiply,
										[]byte{
											byte(bytecode.GET_LOCAL_1),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(64, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol("*").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpMultiply, 1)),
				},
			),
		},

		"resolve static divide": {
			input: "45.5 / .5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Float(91).ToValue(),
				},
			),
		},
		"divide int": {
			input: "a := 1; a / 5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.DIVIDE_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{},
			),
		},
		"divide float": {
			input: "a := 45.5; a / .5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.DIVIDE_FLOAT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.Float(45.5).ToValue(),
					value.Float(0.5).ToValue(),
				},
			),
		},
		"divide builtin": {
			input: "a := 1i8; a / 5i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT8), 5,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"divide value": {
			input: `
				module Foo
					def /(other: Int): Int
					  other
					end
				end
				Foo / 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INT_5),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(7, 6),
				},
				[]value.Value{
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							namespaceDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.DEF_NAMESPACE), 0,
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(85, 7, 12)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Root").ToValue(),
								value.ToSymbol("Foo").ToValue(),
							},
						),
					),
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							methodDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.GET_SINGLETON),
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.LOAD_VALUE_2),
								byte(bytecode.DEF_METHOD),
								byte(bytecode.POP),
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(85, 7, 12)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpDivide,
										[]byte{
											byte(bytecode.GET_LOCAL_1),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(64, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol("/").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpDivide, 1)),
				},
			),
		},

		"resolve static exponentiate": {
			input: "-2 ** 2",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 0xFC,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"exponentiate int": {
			input: "a := -2; a ** 2",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 0xFE,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_2),
					byte(bytecode.EXPONENTIATE_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{},
			),
		},
		"exponentiate float": {
			input: "a := 1.2; a + 5.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.ADD_FLOAT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.Float(1.2).ToValue(),
					value.Float(5.0).ToValue(),
				},
			),
		},
		"exponentiate builtin": {
			input: "a := 1i8; a ** 5i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT8), 5,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"exponentiate value": {
			input: `
				module Foo
					def **(other: Int): Int
					  other
					end
				end
				Foo ** 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INT_5),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(87, 7, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(7, 6),
				},
				[]value.Value{
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							namespaceDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.DEF_NAMESPACE), 0,
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(87, 7, 13)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Root").ToValue(),
								value.ToSymbol("Foo").ToValue(),
							},
						),
					),
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							methodDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.GET_SINGLETON),
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.LOAD_VALUE_2),
								byte(bytecode.DEF_METHOD),
								byte(bytecode.POP),
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(87, 7, 13)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpExponentiate,
										[]byte{
											byte(bytecode.GET_LOCAL_1),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(65, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol("**").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpExponentiate, 1)),
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

func TestUnaryExpressions(t *testing.T) {
	tests := testTable{
		"resolve static negate": {
			input: "-5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 0xfb,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(1, 1, 2)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"negate int": {
			input: "a := 5; -a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEGATE_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{},
			),
		},
		"negate float": {
			input: "a := 5.2; -a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEGATE_FLOAT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.Float(5.2).ToValue(),
				},
			),
		},
		"negate builtin": {
			input: "a := 5i8; -a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT8), 5,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{},
			),
		},
		"negate value": {
			input: `
				module Foo
					def -@: Foo
					  self
					end
				end
				-Foo
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(70, 7, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(7, 5),
				},
				[]value.Value{
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							namespaceDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.DEF_NAMESPACE), 0,
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(70, 7, 9)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Root").ToValue(),
								value.ToSymbol("Foo").ToValue(),
							},
						),
					),
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							methodDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.GET_SINGLETON),
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.LOAD_VALUE_2),
								byte(bytecode.DEF_METHOD),
								byte(bytecode.POP),
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(70, 7, 9)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(7, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpNegate,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(52, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										0,
										0,
										nil,
									),
								),
								symbol.OpNegate.ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpNegate, 0)),
				},
			),
		},

		"resolve static bitwise not": {
			input: "~10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 0xf5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"resolve static logical not": {
			input: "!10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
		},
		"logical not": {
			input: "a := 10; !a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{},
			),
		},
		"bitwise not": {
			input: "a := 10; ~a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.BITWISE_NOT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{},
			),
		},

		"resolve static plus": {
			input: "+5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_5),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(1, 1, 2)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
		},
		"unary plus": {
			input: "a := 10; +a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.UNARY_PLUS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestComplexAssignmentLocals(t *testing.T) {
	tests := testTable{
		"increment": {
			input: "a := 1; a++",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INCREMENT_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{},
			),
		},
		"decrement": {
			input: "a := 1; a--",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.DECREMENT_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{},
			),
		},
		"add": {
			input: "a := 1; a += 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"subtract": {
			input: "a := 1; a -= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.SUBTRACT_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"multiply": {
			input: "a := 1; a *= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.MULTIPLY_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"divide": {
			input: "a := 1; a /= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.DIVIDE_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"exponentiate": {
			input: "a := 1; a **= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.EXPONENTIATE_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"modulo": {
			input: "a := 1; a %= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.MODULO_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"bitwise AND": {
			input: "a := 1; a &= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.BITWISE_AND_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"bitwise OR": {
			input: "a := 1; a |= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.BITWISE_OR_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"bitwise XOR": {
			input: "a := 1; a ^= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.BITWISE_XOR_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"left bitshift": {
			input: "a := 1; a <<= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.LBITSHIFT_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"left logical bitshift": {
			input: "a := 1u64; a <<<= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_UINT64_8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{},
			),
		},
		"right bitshift": {
			input: "a := 1; a >>= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.RBITSHIFT_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{},
			),
		},
		"right logical bitshift": {
			input: "a := 1u64; a >>>= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_UINT64_8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{},
			),
		},
		"logic OR": {
			input: "var a: Int? = 1; a ||= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.INT_3),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 1, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{},
			),
		},
		"logic AND": {
			input: "var a: Int? = 1; a &&= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.INT_3),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 1, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{},
			),
		},
		"nil coalesce": {
			input: "var a: Int? = 1; a ??= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS_NNP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.INT_3),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 1, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestComplexAssignmentInstanceVariables(t *testing.T) {
	tests := testTable{
		"increment": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a++
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(82, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
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
						L(P(0, 1, 1), P(82, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(82, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(82, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_NAME16), 0, 0,
									byte(bytecode.INCREMENT_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(73, 6, 22)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 7),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"decrement": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a--
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(82, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
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
						L(P(0, 1, 1), P(82, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(82, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(82, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_NAME16), 0, 0,
									byte(bytecode.DECREMENT_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(73, 6, 22)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 7),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"add": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a += 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.ADD_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"subtract": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a -= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.SUBTRACT_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"multiply": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a *= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.MULTIPLY_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"divide": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a /= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.DIVIDE_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"exponentiate": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a **= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(86, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
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
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.EXPONENTIATE_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(77, 6, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"modulo": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a %= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.MODULO_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"bitwise AND": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a &= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.BITWISE_AND_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"bitwise OR": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a |= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.BITWISE_OR_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"bitwise XOR": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a ^= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.BITWISE_XOR_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"left bitshift": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a <<= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(86, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.LBITSHIFT_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(77, 6, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"left logical bitshift": {
			input: `
				class Foo
					var @a: UInt64
					init(@a); end

					def foo then @a <<<= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(90, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(90, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(90, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(90, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(40, 4, 6), P(52, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.LOGIC_LBITSHIFT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(60, 6, 6), P(81, 6, 27)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"right bitshift": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a >>= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(86, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.RBITSHIFT_INT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(77, 6, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"right logical bitshift": {
			input: `
				class Foo
					var @a: Int64
					init(@a); end

					def foo then @a >>>= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(89, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(89, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(89, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(89, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(39, 4, 6), P(51, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 7),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.INT_3),
									byte(bytecode.LOGIC_RBITSHIFT),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(59, 6, 6), P(80, 6, 27)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"logic OR": {
			input: `
				class Foo
					var @a: Int?

					def foo then @a ||= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.JUMP_IF_NP), 0, 2,
									byte(bytecode.POP),
									byte(bytecode.INT_3),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(39, 5, 6), P(59, 5, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 9),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"nil coalesce": {
			input: `
				class Foo
					var @a: Int?

					def foo then @a ??= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.JUMP_UNLESS_NNP), 0, 2,
									byte(bytecode.POP),
									byte(bytecode.INT_3),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(39, 5, 6), P(59, 5, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 9),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"logic AND": {
			input: `
				class Foo
					var @a: Int?

					def foo then @a &&= 3
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
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
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<ivarIndices>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(&value.IvarIndices{
								value.ToSymbol("a"): 0,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR_0),
									byte(bytecode.JUMP_UNLESS_NP), 0, 2,
									byte(bytecode.POP),
									byte(bytecode.INT_3),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.RETURN),
								},
								L(P(39, 5, 6), P(59, 5, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 9),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
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

func TestBitwiseAnd(t *testing.T) {
	tests := testTable{
		"resolve static AND": {
			input: "23 & 10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_2),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
		},
		"resolve static nested AND": {
			input: "23 & 15 & 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 6,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"compile runtime AND": {
			input: "a := 23; a & 15 & 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 0x17,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 0x0F,
					byte(bytecode.BITWISE_AND_INT),
					byte(bytecode.LOAD_INT_8), 0x2E,
					byte(bytecode.BITWISE_AND_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestBitwiseAndNot(t *testing.T) {
	tests := testTable{
		"resolve static AND NOT": {
			input: "23 &~ 10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 0x15,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"resolve static nested AND NOT": {
			input: "23 &~ 15 &~ 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 0x10,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"compile runtime AND NOT": {
			input: "a := 23; a &~ 15 &~ 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 0x17,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 0x0F,
					byte(bytecode.BITWISE_AND_NOT),
					byte(bytecode.LOAD_INT_8), 0x2E,
					byte(bytecode.BITWISE_AND_NOT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 1, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestBitwiseOr(t *testing.T) {
	tests := testTable{
		"resolve static OR": {
			input: "23 | 10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 0x1F,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"resolve static nested OR": {
			input: "23 | 15 | 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 0x3F,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"compile runtime OR": {
			input: "a := 23; a | 15 | 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 0x17,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 0x0F,
					byte(bytecode.BITWISE_OR_INT),
					byte(bytecode.LOAD_INT_8), 0x2E,
					byte(bytecode.BITWISE_OR_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestBitwiseXor(t *testing.T) {
	tests := testTable{
		"resolve static XOR": {
			input: "23 ^ 10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 0x1D,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"resolve static nested XOR": {
			input: "23 ^ 15 ^ 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 0x36,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"compile runtime XOR": {
			input: "a := 23; a ^ 15 ^ 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 23,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 15,
					byte(bytecode.BITWISE_XOR_INT),
					byte(bytecode.LOAD_INT_8), 46,
					byte(bytecode.BITWISE_XOR_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestModulo(t *testing.T) {
	tests := testTable{
		"resolve static modulo": {
			input: "23 % 10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_3),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
		},
		"resolve static nested modulo": {
			input: "24 % 15 % 2",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
		},
		"compile runtime modulo": {
			input: "a := 24; a % 15 % 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 24,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 15,
					byte(bytecode.MODULO_INT),
					byte(bytecode.LOAD_INT_8), 46,
					byte(bytecode.MODULO_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}
