package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func TestSubscript(t *testing.T) {
	tests := testTable{
		"static": {
			input: "[5, 3][0]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_5),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
		},
		"dynamic": {
			input: `
				arr := [5, 3]
				arr[1]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 3, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"dynamic nil safe": {
			input: `
				var arr: List[Int]? = [5, 3]
				arr?[1]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF_NIL), 0, 7,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 3, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 12),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Ref(&value.CallSiteInfo{
						Name:          value.ToSymbol("[]"),
						ArgumentCount: 1,
					}),
				},
			),
		},
		"overload": {
			input: `
				module Foo
					overload def [](a: String): String then a
					overload def [](a: Int): Float then a.to_float
				end

				a := Foo["lol"]
				b := Foo[1]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INT_1),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(159, 8, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 6),
					bytecode.NewLineInfo(8, 8),
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
							L(P(0, 1, 1), P(159, 8, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(8, 2),
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
								byte(bytecode.LOAD_VALUE_3),
								byte(bytecode.LOAD_VALUE8), 4,
								byte(bytecode.DEF_METHOD),
								byte(bytecode.POP),
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(159, 8, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 11),
								bytecode.NewLineInfo(8, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										value.ToSymbol("[]"),
										[]byte{
											byte(bytecode.GET_LOCAL_1),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(61, 3, 46)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(3, 2),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol("[]").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										value.ToSymbol("[]@1"),
										[]byte{
											byte(bytecode.GET_LOCAL_1),
											byte(bytecode.CALL_METHOD8), 0,
											byte(bytecode.RETURN),
										},
										L(P(68, 4, 6), P(113, 4, 51)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 4),
										},
										1,
										0,
										[]value.Value{
											value.Ref(&value.CallSiteInfo{
												Name:          value.ToSymbol("to_float"),
												ArgumentCount: 0,
											}),
										},
									),
								),
								value.ToSymbol("[]@1").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.String("lol")),
					value.Ref(&value.CallSiteInfo{
						Name:          value.ToSymbol("[]"),
						ArgumentCount: 1,
					}),
					value.Ref(&value.CallSiteInfo{
						Name:          value.ToSymbol("[]@1"),
						ArgumentCount: 1,
					}),
				},
			),
		},

		"setter builtin": {
			input: `
				arr := [5, 3]
				arr[1] = 15
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_INT_8), 15,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 6),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"setter value": {
			input: `
				var arr: List[Int] = [5, 3]
				arr[1] = 15
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_INT_8), 15,
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(48, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Ref(value.NewCallSiteInfo(symbol.OpSubscriptSet, 2)),
				},
			),
		},
		"setter overload": {
			input: `
				module Foo
					overload def []=(a: String, b: String); end
					overload def []=(a: Int, b: Int); end
				end
				Foo[1] = 15
				Foo["lol"] = "foo"
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
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_INT_8), 15,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.CALL_METHOD8), 6,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(154, 7, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(6, 8),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(154, 7, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(154, 7, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("[]="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(63, 3, 48)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("[]=").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("[]=@1"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(70, 4, 6), P(106, 4, 42)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 2),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("[]=@1").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("[]=@1"), 2)),
					value.Ref(value.String("lol")),
					value.Ref(value.String("foo")),
					value.Ref(value.NewCallSiteInfo(symbol.OpSubscriptSet, 2)),
				},
			),
		},

		"increment Int": {
			input: `
				arr := [5, 3]
				arr[1]++
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.DUP_2),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.INCREMENT_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(31, 3, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"increment builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1]++
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.DUP_2),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.INCREMENT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"increment value": {
			input: `
				module Foo
					def ++: Foo
						self
					end
				end

				arr := [Foo]
				arr[1]++
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.DUP_2),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(91, 9, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 8),
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
							L(P(0, 1, 1), P(91, 9, 13)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(91, 9, 13)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpIncrement,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(51, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										0,
										0,
										nil,
									),
								),
								value.ToSymbol("++").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpIncrement, 0)),
				},
			),
		},

		"decrement Int": {
			input: `
				arr := [5, 3]
				arr[1]--
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.DUP_2),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DECREMENT_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(31, 3, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"decrement builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1]--
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.DUP_2),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DECREMENT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"decrement value": {
			input: `
				module Foo
					def --: Foo
						self
					end
				end

				arr := [Foo]
				arr[1]--
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.DUP_2),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(91, 9, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 8),
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
							L(P(0, 1, 1), P(91, 9, 13)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(91, 9, 13)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpDecrement,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(51, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										0,
										0,
										nil,
									),
								),
								value.ToSymbol("--").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpDecrement, 0)),
				},
			),
		},

		"add Int": {
			input: `
				arr := [5, 3]
				arr[1] += 15
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 15,
					byte(bytecode.ADD_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"add Float": {
			input: `
				arr := [5.5, 3.9]
				arr[1] += 3.0
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.ADD_FLOAT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.Float(5.5).ToValue(),
						value.Float(3.9).ToValue(),
					}),
					value.Float(3.0).ToValue(),
				},
			),
		},
		"add builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] += 3u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 3,
					byte(bytecode.ADD),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"add value": {
			input: `
				module Foo
					def +(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] += 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(105, 9, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpAdd,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(62, 5, 8)),
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
		"add overload": {
			input: `
				module Foo
					overload def [](other: Int): Int
						other
					end
					overload def []=(key: Int, value: Int); end

					overload def [](other: String): String
						other
					end
					overload def []=(key: String, value: String); end
				end

				Foo[1] += 8
				Foo["foo"] += "bar"
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
					byte(bytecode.INT_1),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INT_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.ADD_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.CALL_METHOD8), 6,
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.ADD),
					byte(bytecode.CALL_METHOD8), 8,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(293, 15, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(14, 14),
					bytecode.NewLineInfo(15, 16),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(293, 15, 24)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(15, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 5,
							byte(bytecode.LOAD_VALUE8), 6,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 7,
							byte(bytecode.LOAD_VALUE8), 8,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(293, 15, 24)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 21),
							bytecode.NewLineInfo(15, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("[]"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(73, 5, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 1),
									bytecode.NewLineInfo(5, 1),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("[]").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("[]@1"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.RETURN),
								},
								L(P(130, 8, 6), P(188, 10, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(9, 1),
									bytecode.NewLineInfo(10, 1),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("[]@1").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("[]="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(80, 6, 6), P(122, 6, 48)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 2),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("[]=").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("[]=@1"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(195, 11, 6), P(243, 11, 54)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(11, 2),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("[]=@1").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("[]"), 1)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("[]="), 2)),
					value.Ref(value.String("foo")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("[]@1"), 1)),
					value.Ref(value.String("bar")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("[]=@1"), 2)),
				},
			),
		},

		"subtract Int": {
			input: `
				arr := [5, 3]
				arr[1] -= 2
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.INT_2),
					byte(bytecode.SUBTRACT_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"subtract Float": {
			input: `
				arr := [5.5, 3.9]
				arr[1] -= 3.0
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.SUBTRACT_FLOAT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.Float(5.5).ToValue(),
						value.Float(3.9).ToValue(),
					}),
					value.Float(3.0).ToValue(),
				},
			),
		},
		"subtract builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] -= 3u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 3,
					byte(bytecode.SUBTRACT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"subtract value": {
			input: `
				module Foo
					def -(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] -= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(105, 9, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpSubtract,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(62, 5, 8)),
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

		"multiply Int": {
			input: `
				arr := [5, 3]
				arr[1] *= 3
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.INT_3),
					byte(bytecode.MULTIPLY_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"multiply Float": {
			input: `
				arr := [5.5, 3.9]
				arr[1] *= 3.0
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.MULTIPLY_FLOAT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.Float(5.5).ToValue(),
						value.Float(3.9).ToValue(),
					}),
					value.Float(3.0).ToValue(),
				},
			),
		},
		"multiply builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] *= 3u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 3,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"multiply value": {
			input: `
				module Foo
					def *(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] *= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(105, 9, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpMultiply,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(62, 5, 8)),
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

		"divide Int": {
			input: `
				arr := [5, 3]
				arr[1] /= 10
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.DIVIDE_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"divide Float": {
			input: `
				arr := [5.2, 3.0]
				arr[1] /= 10.9
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.DIVIDE_FLOAT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.Float(5.2).ToValue(),
						value.Float(3.0).ToValue(),
					}),
					value.Float(10.9).ToValue(),
				},
			),
		},
		"divide builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] /= 10u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 10,
					byte(bytecode.DIVIDE),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"divide value": {
			input: `
				module Foo
					def /(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] /= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(105, 9, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpDivide,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(62, 5, 8)),
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

		"exponentiate Int": {
			input: `
				arr := [5, 3]
				arr[1] **= 12
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 12,
					byte(bytecode.EXPONENTIATE_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(36, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"exponentiate builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] **= 2u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 2,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"exponentiate value": {
			input: `
				module Foo
					def **(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] **= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 9, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(107, 9, 17)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(107, 9, 17)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpExponentiate,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(63, 5, 8)),
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

		"modulo Int": {
			input: `
				arr := [5, 3]
				arr[1] %= 2
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.INT_2),
					byte(bytecode.MODULO_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"modulo Float": {
			input: `
				arr := [5.5, 3.9]
				arr[1] %= 2.5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.MODULO_FLOAT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.Float(5.5).ToValue(),
						value.Float(3.9).ToValue(),
					}),
					value.Float(2.5).ToValue(),
				},
			),
		},
		"modulo builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] %= 2u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 2,
					byte(bytecode.MODULO),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"modulo value": {
			input: `
				module Foo
					def %(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] %= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(105, 9, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpModulo,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(62, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol("%").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpModulo, 1)),
				},
			),
		},

		"bitwise AND Int": {
			input: `
				arr := [5, 3]
				arr[1] &= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.BITWISE_AND_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"bitwise AND builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] &= 8u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 8,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"bitwise AND value": {
			input: `
				module Foo
					def &(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] &= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(105, 9, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpAnd,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(62, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol("&").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpAnd, 1)),
				},
			),
		},
		"bitwise OR Int": {
			input: `
				arr := [5, 3]
				arr[1] |= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.BITWISE_OR_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"bitwise OR builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] |= 8u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 8,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"bitwise OR value": {
			input: `
				module Foo
					def |(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] |= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(105, 9, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpOr,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(62, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol("|").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpOr, 1)),
				},
			),
		},
		"bitwise XOR Int": {
			input: `
				arr := [5, 3]
				arr[1] ^= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.BITWISE_XOR_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"bitwise XOR builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] ^= 8u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 8,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"bitwise XOR value": {
			input: `
				module Foo
					def ^(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] ^= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(105, 9, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(105, 9, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpXor,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(62, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol("^").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpXor, 1)),
				},
			),
		},
		"left bitshift Int": {
			input: `
				arr := [5, 3]
				arr[1] <<= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.LBITSHIFT_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"left bitshift builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] <<= 8u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 8,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"left bitshift value": {
			input: `
				module Foo
					def <<(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] <<= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 9, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(107, 9, 17)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(107, 9, 17)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpLeftBitshift,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(63, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol("<<").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpLeftBitshift, 1)),
				},
			),
		},
		"right bitshift Int": {
			input: `
				arr := [5, 3]
				arr[1] >>= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.RBITSHIFT_INT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"right bitshift builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] >>= 8u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 8,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"right bitshift value": {
			input: `
				module Foo
					def >>(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] >>= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 9, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(107, 9, 17)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(107, 9, 17)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpRightBitshift,
										[]byte{
											byte(bytecode.SELF),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(63, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol(">>").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpRightBitshift, 1)),
				},
			),
		},

		"left logical bitshift Int": {
			input: `
				arr := [5u64, 3u64]
				arr[1] <<<= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt64(5).ToValue(),
						value.UInt64(3).ToValue(),
					}),
				},
			),
		},
		"left logical bitshift builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] <<<= 8u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 8,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 3, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"left logical bitshift value": {
			input: `
				module Foo
					def <<<(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] <<<= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(109, 9, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(109, 9, 18)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(109, 9, 18)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpLogicalLeftBitshift,
										[]byte{
											byte(bytecode.SELF),
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
								value.ToSymbol("<<<").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpLogicalLeftBitshift, 1)),
				},
			),
		},
		"right logical bitshift Int": {
			input: `
				arr := [5u64, 3u64]
				arr[1] >>>= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt64(5).ToValue(),
						value.UInt64(3).ToValue(),
					}),
				},
			),
		},
		"right logical bitshift builtin": {
			input: `
				arr := [5u8, 3u8]
				arr[1] >>>= 8u8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_UINT8), 8,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 3, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.UInt8(5).ToValue(),
						value.UInt8(3).ToValue(),
					}),
				},
			),
		},
		"right logical bitshift value": {
			input: `
				module Foo
					def >>>(other: Int): Foo
						self
					end
				end

				arr := [Foo]
				arr[1] >>>= 8
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(109, 9, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 11),
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
							L(P(0, 1, 1), P(109, 9, 18)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 5),
								bytecode.NewLineInfo(9, 2),
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
							L(P(0, 1, 1), P(109, 9, 18)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 7),
								bytecode.NewLineInfo(9, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpLogicalRightBitshift,
										[]byte{
											byte(bytecode.SELF),
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
								value.ToSymbol(">>>").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpLogicalRightBitshift, 1)),
				},
			),
		},

		"logic OR": {
			input: `
				var arr: List[Int?] = [5, 3]
				arr[1] ||= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.JUMP_IF_NP), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(50, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 15),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Ref(value.NewCallSiteInfo(symbol.OpSubscript, 1)),
					value.Ref(value.NewCallSiteInfo(symbol.OpSubscriptSet, 2)),
				},
			),
		},
		"logic AND": {
			input: `
				var arr: List[Int?] = [5, 3]
				arr[1] &&= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.JUMP_UNLESS_NP), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(50, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 15),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Ref(value.NewCallSiteInfo(symbol.OpSubscript, 1)),
					value.Ref(value.NewCallSiteInfo(symbol.OpSubscriptSet, 2)),
				},
			),
		},
		"nil coalesce": {
			input: `
				var arr: List[Int?] = [5, 3]
				arr[1] ??= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.JUMP_UNLESS_NNP), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(50, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 15),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Ref(value.NewCallSiteInfo(symbol.OpSubscript, 1)),
					value.Ref(value.NewCallSiteInfo(symbol.OpSubscriptSet, 2)),
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

func TestInstantiate(t *testing.T) {
	tests := testTable{
		"without arguments": {
			input: `
				class Foo; end
				::Foo()
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(31, 3, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(3, 5),
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
						L(P(0, 1, 1), P(31, 3, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(3, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
				},
			),
		},
		"complex constant": {
			input: `
				module Foo
					module Bar
						class Baz; end
					end
				end
				::Foo::Bar::Baz()
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(91, 7, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(7, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 5,
							byte(bytecode.GET_CONST8), 6,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(91, 7, 22)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 21),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Foo::Bar").ToValue(),
							value.ToSymbol("Baz").ToValue(),
							value.ToSymbol("Foo::Bar::Baz").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.ToSymbol("Foo::Bar::Baz").ToValue(),
				},
			),
		},
		"with positional arguments": {
			input: `
				class Foo
					init(a: Int, b: String); end
				end
				::Foo(1, 'lol')
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
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INSTANTIATE8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 5, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 7),
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
						L(P(0, 1, 1), P(76, 5, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(5, 2),
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
						L(P(0, 1, 1), P(76, 5, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(20, 3, 6), P(47, 3, 33)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.String("lol")),
				},
			),
		},
		"with named args": {
			input: `
				class Foo
					init(a: Int, b: String); end
				end
				::Foo(1, b: 'lol')
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
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INSTANTIATE8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(79, 5, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 7),
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
						L(P(0, 1, 1), P(79, 5, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(5, 2),
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
						L(P(0, 1, 1), P(79, 5, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(20, 3, 6), P(47, 3, 33)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.String("lol")),
				},
			),
		},
		"with duplicated named args": {
			input: `
				class Foo
					init(a: String, b: Int); end
				end
				::Foo(b: 1, a: 'lol', b: 2)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(83, 5, 27), P(86, 5, 30)),
					"duplicated argument `b` in call to `Foo.:#init`",
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestCallMethod(t *testing.T) {
	tests := testTable{
		"call method overloads": {
			input: `
				module Foo
					overload def foo(a: String); end
					overload def foo(a: Int); end
				end
				a := Foo
				a.foo(1)
				a.foo("lol")
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(139, 8, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 5),
					bytecode.NewLineInfo(8, 6),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(139, 8, 17)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(139, 8, 17)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(52, 3, 37)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo@1"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(59, 4, 6), P(87, 4, 34)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo@1").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo@1"), 1)),
					value.Ref(value.String("lol")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 1)),
				},
			),
		},
		"call variable overloads": {
			input: `
				module Foo
					overload def call(a: String); end
					overload def call(a: Int); end
				end
				a := Foo
				a(1)
				a("lol")
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.CALL8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(133, 8, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 5),
					bytecode.NewLineInfo(8, 6),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(133, 8, 13)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(133, 8, 13)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("call"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(53, 3, 38)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("call").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("call@1"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(60, 4, 6), P(89, 4, 35)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("call@1").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("call@1"), 1)),
					value.Ref(value.String("lol")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("call"), 1)),
				},
			),
		},
		"call a method without arguments": {
			input: `
				module Foo
					def foo; end
				end
				Foo.foo
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
				L(P(0, 1, 1), P(53, 5, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(53, 5, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
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
						L(P(0, 1, 1), P(53, 5, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(32, 3, 17)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
				},
			),
		},
		"call": {
			input: `
				module Foo
					def call; end
				end
				Foo.call()
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
					byte(bytecode.CALL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(57, 5, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 5),
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
						L(P(0, 1, 1), P(57, 5, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(57, 5, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("call"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(33, 3, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								nil,
							)),
							value.ToSymbol("call").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("call"), 0)),
				},
			),
		},
		"call special syntax": {
			input: `
				module Foo
					def call; end
				end
				Foo.()
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
					byte(bytecode.CALL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(53, 5, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 5),
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
						L(P(0, 1, 1), P(53, 5, 11)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(53, 5, 11)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("call"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(33, 3, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								nil,
							)),
							value.ToSymbol("call").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("call"), 0)),
				},
			),
		},
		"call getter": {
			input: `
				module Foo
					def call; end
				end
				Foo.call
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
					byte(bytecode.CALL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(55, 5, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 5),
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
						L(P(0, 1, 1), P(55, 5, 13)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(55, 5, 13)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("call"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(33, 3, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								nil,
							)),
							value.ToSymbol("call").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("call"), 0)),
				},
			),
		},
		"make a cascade call without arguments": {
			input: `
				module Foo
					def foo; end
				end
				Foo..foo
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
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(54, 5, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 7),
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
						L(P(0, 1, 1), P(54, 5, 13)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(54, 5, 13)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(32, 3, 17)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
				},
			),
		},
		"make a cascade call without arguments nil safe": {
			input: `
				module Foo
					def foo; end
				end
				var a: Foo? = nil
				a?..foo
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.DUP),
					byte(bytecode.JUMP_IF_NIL_NP), 0, 3,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(75, 6, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 9),
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
						L(P(0, 1, 1), P(75, 6, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(75, 6, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(32, 3, 17)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
				},
			),
		},
		"call a method without arguments nil safe": {
			input: `
				module Foo
					def foo; end
				end
				var a: Foo? = nil
				a?.foo
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF_NIL_NP), 0, 2,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(74, 6, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 7),
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
						L(P(0, 1, 1), P(74, 6, 11)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(74, 6, 11)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(32, 3, 17)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
				},
			),
		},
		"call a setter": {
			input: `
				module Foo
					def foo=(value: Int); end
				end
				Foo.foo = 3
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
					byte(bytecode.INT_3),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(70, 5, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 6),
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
						L(P(0, 1, 1), P(70, 5, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(70, 5, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(21, 3, 6), P(45, 3, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a method with positional arguments": {
			input: `
				module Foo
					def foo(a: Int, b: String); end
				end
				Foo.foo(1, 'lol')
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
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(82, 5, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 7),
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
						L(P(0, 1, 1), P(82, 5, 22)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(82, 5, 22)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(51, 3, 36)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.String("lol")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 2)),
				},
			),
		},
		"call a method with rest arguments": {
			input: `
				module Foo
					def foo(*a: Int); end
				end
				Foo.foo(1, 2, 3)
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
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 5, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 6),
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
						L(P(0, 1, 1), P(71, 5, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(71, 5, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(41, 3, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 1)),
				},
			),
		},
		"call a method with a splat argument": {
			input: `
				module Foo
					def foo(*a: Int); end
				end
				arr := [1, 2, 3]
				Foo.foo(*arr)
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(89, 6, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 6),
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
						L(P(0, 1, 1), P(89, 6, 18)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(89, 6, 18)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(41, 3, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 1)),
				},
			),
		},
		"call a method with a non-tuple splat argument": {
			input: `
				module Foo
					def foo(*a: Int); end
				end
				Foo.foo(*3)
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
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_ARRAY_TUPLE8), 0,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 8,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.APPEND),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INCREMENT_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 13,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(66, 5, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 26),
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
						L(P(0, 1, 1), P(66, 5, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(66, 5, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(41, 3, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 1)),
				},
			),
		},
		"call a method with rest and splat arguments": {
			input: `
				module Foo
					def foo(*a: Int); end
				end
				arr := [1, 2, 3]
				Foo.foo(5, *arr, 10)
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 3,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.COPY),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.FOR_IN_BUILTIN), 0, 6,
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.APPEND),
					byte(bytecode.LOOP), 0, 10,
					byte(bytecode.LEAVE_SCOPE16), 3, 2,
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.APPEND),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(96, 6, 25)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 27),
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
						L(P(0, 1, 1), P(96, 6, 25)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(96, 6, 25)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(41, 3, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(&value.ArrayTuple{
						value.SmallInt(5).ToValue(),
					}),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 1)),
				},
			),
		},

		"call a method with named rest arguments": {
			input: `
				module Foo
					def foo(**a: Int); end
				end
				Foo.foo(foo: 1, bar: 2, baz: 3)
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
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(87, 5, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(5, 6),
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
						L(P(0, 1, 1), P(87, 5, 36)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(87, 5, 36)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(42, 3, 27)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.MustNewHashRecordWithElements(
						nil,
						value.Pair{Key: value.ToSymbol("foo").ToValue(), Value: value.SmallInt(1).ToValue()},
						value.Pair{Key: value.ToSymbol("bar").ToValue(), Value: value.SmallInt(2).ToValue()},
						value.Pair{Key: value.ToSymbol("baz").ToValue(), Value: value.SmallInt(3).ToValue()},
					)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 1)),
				},
			),
		},
		"call a method with a double splat argument": {
			input: `
				module Foo
					def foo(**a: Int); end
				end
				map := { foo: 1, bar: 2, baz: 3 }
				Foo.foo(**map)
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
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(108, 6, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 6),
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
						L(P(0, 1, 1), P(108, 6, 19)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(108, 6, 19)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(42, 3, 27)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.Ref(vm.MustNewHashMapWithElements(
						nil,
						value.Pair{Key: value.ToSymbol("foo").ToValue(), Value: value.SmallInt(1).ToValue()},
						value.Pair{Key: value.ToSymbol("bar").ToValue(), Value: value.SmallInt(2).ToValue()},
						value.Pair{Key: value.ToSymbol("baz").ToValue(), Value: value.SmallInt(3).ToValue()},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 1)),
				},
			),
		},
		"call a method with a non-record double splat argument": {
			input: `
				module Foo
					def foo(**a: Int); end
				end
				arr := [Pair(:foo, 1), Pair(:bar, 2), Pair(:baz, 3)]
				Foo.foo(**arr)
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 4,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INT_1),
					byte(bytecode.INSTANTIATE8), 2,
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.INT_2),
					byte(bytecode.INSTANTIATE8), 2,
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.INT_3),
					byte(bytecode.INSTANTIATE8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 3,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_CONST8), 6,
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_HASH_RECORD8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.FOR_IN_BUILTIN), 0, 44,
					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS_NP), 0, 24,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 7,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 13,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 8,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_4),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_IF), 0, 3,
					byte(bytecode.LOAD_VALUE8), 9,
					byte(bytecode.THROW),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.MAP_SET),
					byte(bytecode.LOOP), 0, 48,
					byte(bytecode.LEAVE_SCOPE16), 4, 3,
					byte(bytecode.CALL_METHOD8), 10,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(127, 6, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 25),
					bytecode.NewLineInfo(6, 62),
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
						L(P(0, 1, 1), P(127, 6, 19)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(127, 6, 19)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(42, 3, 27)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Std::Pair").ToValue(),
					value.ToSymbol("foo").ToValue(),
					value.ToSymbol("bar").ToValue(),
					value.ToSymbol("baz").ToValue(),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("key"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("value"), 0)),
					value.Ref(value.NewError(value.PatternNotMatchedErrorClass, "assigned value does not match the pattern defined in for in loop")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 1)),
				},
			),
		},
		"call a method with rest and double splat arguments": {
			input: `
				module Foo
					def foo(**a: Int); end
				end
				map := { foo: 1, bar: 2, baz: 3 }
				Foo.foo(elo: 5, **map, pipa: 10)
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 4,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.NEW_HASH_RECORD8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.FOR_IN_BUILTIN), 0, 44,
					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 5,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS_NP), 0, 24,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 6,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 13,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 7,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_4),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_IF), 0, 3,
					byte(bytecode.LOAD_VALUE8), 8,
					byte(bytecode.THROW),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.MAP_SET),
					byte(bytecode.LOOP), 0, 48,
					byte(bytecode.LEAVE_SCOPE16), 4, 3,
					byte(bytecode.LOAD_VALUE8), 9,
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.MAP_SET),
					byte(bytecode.CALL_METHOD8), 10,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(126, 6, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 68),
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
						L(P(0, 1, 1), P(126, 6, 37)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(126, 6, 37)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(42, 3, 27)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.Ref(vm.MustNewHashMapWithElements(
						nil,
						value.Pair{Key: value.ToSymbol("foo").ToValue(), Value: value.SmallInt(1).ToValue()},
						value.Pair{Key: value.ToSymbol("bar").ToValue(), Value: value.SmallInt(2).ToValue()},
						value.Pair{Key: value.ToSymbol("baz").ToValue(), Value: value.SmallInt(3).ToValue()},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(vm.MustNewHashRecordWithElements(
						nil,
						value.Pair{Key: value.ToSymbol("elo").ToValue(), Value: value.SmallInt(5).ToValue()},
					)),
					value.ToSymbol("Std::Pair").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("key"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("value"), 0)),
					value.Ref(value.NewError(value.PatternNotMatchedErrorClass, "assigned value does not match the pattern defined in for in loop")),
					value.ToSymbol("pipa").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 1)),
				},
			),
		},

		"call a method with positional arguments nil safe": {
			input: `
				module Foo
					def foo(a: Int, b: String); end
				end
				var a: Foo? = nil
				a?.foo(1, 'lol')
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF_NIL_NP), 0, 4,
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(103, 6, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 9),
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
						L(P(0, 1, 1), P(103, 6, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(103, 6, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(51, 3, 36)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.Ref(value.String("lol")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 2)),
				},
			),
		},
		"call a method on a local variable": {
			input: `
				module Foo
					def foo(a: Int, b: String); end
				end
				a := Foo
				a.foo(1, 'lol')
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(93, 6, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 6),
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
						L(P(0, 1, 1), P(93, 6, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(93, 6, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(51, 3, 36)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.String("lol")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 2)),
				},
			),
		},
		"call a method on a local variable with named args": {
			input: `
				module Foo
					def foo(a: Int, b: String); end
				end
				a := Foo
				a.foo(1, b: 'lol')
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(96, 6, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 6),
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
						L(P(0, 1, 1), P(96, 6, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(96, 6, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(51, 3, 36)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.String("lol")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 2)),
				},
			),
		},
		"call a method with duplicated named args": {
			input: `
				module Foo
					def foo(a: String, b: Int); end
				end
				Foo.foo(b: 1, a: 'lol', b: 2)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(89, 5, 29), P(92, 5, 32)),
					"duplicated argument `b` in call to `Foo::foo`",
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestCallFunction(t *testing.T) {
	tests := testTable{
		"call a function from using all with a class": {
			input: `
				using Bar::*
				class Bar
					singleton
						def foo; end
					end
				end
				foo()
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
				L(P(0, 1, 1), P(92, 8, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(8, 5),
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
						L(P(0, 1, 1), P(92, 8, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(92, 8, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(53, 5, 7), P(64, 5, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 2),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
				},
			),
		},
		"call a function from using all with a module": {
			input: `
				using Bar::*
				module Bar
					def foo; end
				end
				foo()
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
				L(P(0, 1, 1), P(68, 6, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(6, 5),
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
						L(P(0, 1, 1), P(68, 6, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(68, 6, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(38, 4, 6), P(49, 4, 17)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 2),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
				},
			),
		},
		"call a variable": {
			input: `
				module Bar
					def call; end
				end
				a := Bar
				a()
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 4),
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
						L(P(0, 1, 1), P(63, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(63, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("call"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(33, 3, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								nil,
							)),
							value.ToSymbol("call").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("call"), 0)),
				},
			),
		},
		"call a variable instead of a method": {
			input: `
				def a; end
				module Bar
					def call; end
				end
				a := Bar
				a()
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(78, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 4),
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
						L(P(0, 1, 1), P(78, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.GET_CONST8), 3,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.LOAD_VALUE8), 5,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(78, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 16),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("call"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(36, 4, 6), P(48, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 2),
								},
								nil,
							)),
							value.ToSymbol("call").ToValue(),
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("a"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(14, 2, 14)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 2),
								},
								nil,
							)),
							value.ToSymbol("a").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("call"), 0)),
				},
			),
		},
		"call a function from using with a module": {
			input: `
				using Bar::foo
				module Bar
					def foo; end
				end
				foo()
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
				L(P(0, 1, 1), P(70, 6, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(6, 5),
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
						L(P(0, 1, 1), P(70, 6, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(70, 6, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(40, 4, 6), P(51, 4, 17)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 2),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
				},
			),
		},
		"call a function from using with a class": {
			input: `
				using Bar::foo
				class Bar
					singleton
						def foo; end
					end
				end
				foo()
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
				L(P(0, 1, 1), P(94, 8, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(8, 5),
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
						L(P(0, 1, 1), P(94, 8, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(94, 8, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(55, 5, 7), P(66, 5, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 2),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
				},
			),
		},
		"call a function without arguments": {
			input: `
				def foo; end
				foo()
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(27, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(27, 3, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(3, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(16, 2, 16)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 2),
								},
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
				},
			),
		},
		"call a function with positional arguments": {
			input: `
				def foo(a: Int, b: String); end
				foo(1, 'lol')
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(54, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(54, 3, 18)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(3, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(35, 2, 35)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 2),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.String("lol")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 2)),
				},
			),
		},
		"call a function with named args": {
			input: `
				def foo(a: Int, b: String); end
				foo(1, b: 'lol')
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(57, 3, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(57, 3, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(3, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(35, 2, 35)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 2),
								},
								2,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.String("lol")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 2)),
				},
			),
		},
		"call a function with duplicated named args": {
			input: `
				def foo(a: String, b: Int); end
				foo(b: 1, a: 'lol', b: 2)
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(61, 3, 25), P(64, 3, 28)),
					"duplicated argument `b` in call to `Std::Kernel::foo`",
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestCallSetter(t *testing.T) {
	tests := testTable{
		"call a setter": {
			input: `
				module Bar
					def foo=(value: Int); end
				end
				a := Bar
				a.foo = 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_3),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(81, 6, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 5),
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
						L(P(0, 1, 1), P(81, 6, 14)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
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
						L(P(0, 1, 1), P(81, 6, 14)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(21, 3, 6), P(45, 3, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"increment": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo++
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(104, 7, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 8),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(104, 7, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(104, 7, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("++"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"decrement": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo--
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(104, 7, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 8),
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
						L(P(0, 1, 1), P(104, 7, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(104, 7, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("--"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with add": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo += 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INT_3),
					byte(bytecode.ADD_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with subtract": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo -= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INT_3),
					byte(bytecode.SUBTRACT_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with multiply": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo *= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INT_3),
					byte(bytecode.MULTIPLY_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with divide": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo /= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INT_3),
					byte(bytecode.DIVIDE_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with exponentiate": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo **= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INT_3),
					byte(bytecode.EXPONENTIATE_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(108, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with modulo": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo %= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INT_3),
					byte(bytecode.MODULO_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with left bitshift": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo <<= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INT_3),
					byte(bytecode.LBITSHIFT_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(108, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with logic left bitshift": {
			input: `
				module Bar
					def foo: Int64 then 3i64
					def foo=(value: Int64); end
				end
				a := Bar
				a.foo <<<= 3i64
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_INT64_8), 3,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(119, 7, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 10),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(119, 7, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(119, 7, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_INT64_8), 3,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(44, 3, 29)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(51, 4, 6), P(77, 4, 32)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with right bitshift": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo >>= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INT_3),
					byte(bytecode.RBITSHIFT_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(108, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with logic right bitshift": {
			input: `
				module Bar
					def foo: Int64 then 3i64
					def foo=(value: Int64); end
				end
				a := Bar
				a.foo >>>= 3i64
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_INT64_8), 3,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(119, 7, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 10),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(119, 7, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(119, 7, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_INT64_8), 3,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(44, 3, 29)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(51, 4, 6), P(77, 4, 32)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with bitwise and": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo &= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INT_3),
					byte(bytecode.BITWISE_AND_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with bitwise or": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo |= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INT_3),
					byte(bytecode.BITWISE_OR_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with bitwise xor": {
			input: `
				module Bar
					def foo: Int then 3
					def foo=(value: Int); end
				end
				a := Bar
				a.foo ^= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INT_3),
					byte(bytecode.BITWISE_XOR_INT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(46, 4, 6), P(70, 4, 30)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with logic or": {
			input: `
				module Bar
					def foo: Int? then 3
					def foo=(value: Int?); end
				end
				a := Bar
				a.foo ||= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.JUMP_IF_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.INT_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(110, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(40, 3, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(47, 4, 6), P(72, 4, 31)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with logic and": {
			input: `
				module Bar
					def foo: Int? then 3
					def foo=(value: Int?); end
				end
				a := Bar
				a.foo &&= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.INT_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(110, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(40, 3, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(47, 4, 6), P(72, 4, 31)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
				},
			),
		},
		"call a setter with nil coalesce": {
			input: `
				module Bar
					def foo: Int? then 3
					def foo=(value: Int?); end
				end
				a := Bar
				a.foo ??= 3
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
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.JUMP_UNLESS_NNP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.INT_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(110, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 3),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 5),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Bar").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
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
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.INT_3),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(40, 3, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 2),
								},
								0,
								0,
								nil,
							)),
							value.ToSymbol("foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo="),
								[]byte{
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_FIRST_ARG),
								},
								L(P(47, 4, 6), P(72, 4, 31)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 3),
								},
								1,
								0,
								nil,
							)),
							value.ToSymbol("foo=").ToValue(),
						},
					)),
					value.ToSymbol("Bar").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo="), 1)),
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
