package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func TestBytecodeGetConstant(t *testing.T) {
	tests := bytecodeTestTable{
		"absolute path ::Std": {
			input: "::Std",
			want: vm.NewBytecodeFunctionNoParams(
				nil,
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.ToSymbol("Std").ToValue(),
				},
			),
		},
		"absolute nested path ::Std::Float::INF": {
			input: "::Std::Float::INF",
			want: vm.NewBytecodeFunctionNoParams(
				nil,
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.ToSymbol("Std::Float::INF").ToValue(),
				},
			),
		},
		"relative path from using": {
			input: `
				using Std::Float::INF as I
				I
			`,
			want: vm.NewBytecodeFunctionNoParams(
				nil,
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(37, 3, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 3),
				},
				[]value.Value{
					value.ToSymbol("Std::Float::INF").ToValue(),
				},
			),
		},
		"relative path": {
			input: `
				module Foo
					const BAR = 3
					module Baz
						println BAR
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				nil,
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.INT_3),
					byte(bytecode.DEF_CONST),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						nil,
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_NAMESPACE), 0,
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
							value.ToSymbol("Baz").ToValue(),
						},
					)),
					value.ToSymbol("Foo").ToValue(),
					value.ToSymbol("BAR").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						nil,
						value.ToSymbol("<module: Foo>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.INIT_NAMESPACE),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(84, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(4, 4),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							value.ToSymbol("Foo::Baz").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								nil,
								value.ToSymbol("<module: Foo::Baz>"),
								[]byte{
									byte(bytecode.GET_CONST8), 0,
									byte(bytecode.GET_CONST8), 1,
									byte(bytecode.CALL_METHOD_NT8), 2,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(40, 4, 6), P(76, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 6),
									bytecode.NewLineInfo(6, 3),
								},
								[]value.Value{
									value.ToSymbol("Std::Kernel").ToValue(),
									value.ToSymbol("Foo::BAR").ToValue(),
									value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "println@1"), 1)),
								},
							)),
						},
					)),
				},
			),
		},
		"infer namespace in method argument": {
			input: `
				class Color
					const RED: Color = new(255u8, 0u8, 0u8)
					const GREEN: Color = new(0u8, 255u8, 0u8)
					const BLUE: Color = new(0u8, 0u8, 255u8)

					getter red: UInt8
					getter green: UInt8
					getter blue: UInt8
					init(@red, @green, @blue); end
				end

				def baz(col: Color)
					println "col: #col"
				end
				baz(.::RED)
			`,
			wantFn: func(btc bytecodeTestCase) *vm.BytecodeFunction {
				var baz *vm.BytecodeFunction
				return vm.NewBytecodeFunctionNoParams(
					nil,
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
						byte(bytecode.GET_CONST8), 3,
						byte(bytecode.LOAD_VALUE8), 4,
						byte(bytecode.GET_CONST8), 3,
						byte(bytecode.LOAD_UINT8), 255,
						byte(bytecode.LOAD_UINT8), 0,
						byte(bytecode.LOAD_UINT8), 0,
						byte(bytecode.INSTANTIATE8), 3,
						byte(bytecode.DEF_CONST),
						byte(bytecode.GET_CONST8), 3,
						byte(bytecode.LOAD_VALUE8), 5,
						byte(bytecode.GET_CONST8), 3,
						byte(bytecode.LOAD_UINT8), 0,
						byte(bytecode.LOAD_UINT8), 255,
						byte(bytecode.LOAD_UINT8), 0,
						byte(bytecode.INSTANTIATE8), 3,
						byte(bytecode.DEF_CONST),
						byte(bytecode.GET_CONST8), 3,
						byte(bytecode.LOAD_VALUE8), 6,
						byte(bytecode.GET_CONST8), 3,
						byte(bytecode.LOAD_UINT8), 0,
						byte(bytecode.LOAD_UINT8), 0,
						byte(bytecode.LOAD_UINT8), 255,
						byte(bytecode.INSTANTIATE8), 3,
						byte(bytecode.DEF_CONST),
						byte(bytecode.GET_CONST8), 7,
						byte(bytecode.GET_CONST8), 8,
						byte(bytecode.CALL_METHOD_BC8), 9,
						byte(bytecode.RETURN),
					},
					L(P(0, 1, 1), P(345, 16, 16)),
					bytecode.LineInfoList{
						bytecode.NewLineInfo(1, 9),
						bytecode.NewLineInfo(3, 4),
						bytecode.NewLineInfo(0, 2),
						bytecode.NewLineInfo(3, 9),
						bytecode.NewLineInfo(4, 4),
						bytecode.NewLineInfo(0, 2),
						bytecode.NewLineInfo(4, 9),
						bytecode.NewLineInfo(5, 4),
						bytecode.NewLineInfo(0, 2),
						bytecode.NewLineInfo(5, 9),
						bytecode.NewLineInfo(16, 7),
					},
					[]value.Value{
						value.Ref(vm.NewBytecodeFunctionNoParams(
							nil,
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
							L(P(0, 1, 1), P(345, 16, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 10),
								bytecode.NewLineInfo(16, 2),
							},
							[]value.Value{
								value.ToSymbol("Root").ToValue(),
								value.ToSymbol("Color").ToValue(),
								value.ToSymbol("Std::Object").ToValue(),
							},
						)),
						value.Ref(vm.NewBytecodeFunctionNoParams(
							nil,
							ivarIndicesSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.DEF_IVARS),
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(345, 16, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 4),
								bytecode.NewLineInfo(16, 2),
							},
							[]value.Value{
								value.ToSymbol("Color").ToValue(),
								value.Ref(&value.IvarIndices{
									symbol.ToSymbol("blue"):  0,
									symbol.ToSymbol("green"): 1,
									symbol.ToSymbol("red"):   2,
								}),
							},
						)),
						value.Ref(vm.NewBytecodeFunctionNoParams(
							nil,
							methodDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.LOAD_VALUE_1),
								byte(bytecode.LOAD_VALUE_2),
								byte(bytecode.DEF_METHOD),
								byte(bytecode.LOAD_VALUE_3),
								byte(bytecode.INT_0),
								byte(bytecode.DEF_GETTER),
								byte(bytecode.LOAD_VALUE8), 4,
								byte(bytecode.INT_1),
								byte(bytecode.DEF_GETTER),
								byte(bytecode.LOAD_VALUE8), 5,
								byte(bytecode.INT_2),
								byte(bytecode.DEF_GETTER),
								byte(bytecode.POP),
								byte(bytecode.GET_CONST8), 6,
								byte(bytecode.GET_SINGLETON),
								byte(bytecode.LOAD_VALUE8), 7,
								byte(bytecode.LOAD_VALUE8), 8,
								byte(bytecode.DEF_METHOD),
								byte(bytecode.POP),
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(345, 16, 16)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 26),
								bytecode.NewLineInfo(16, 2),
							},
							[]value.Value{
								value.ToSymbol("Color").ToValue(),
								value.Ref(vm.NewBytecodeFunction(
									nil,
									value.ToSymbol("Color.:#init"),
									[]byte{
										byte(bytecode.GET_LOCAL_1),
										byte(bytecode.DUP),
										byte(bytecode.SET_IVAR_2),
										byte(bytecode.POP),
										byte(bytecode.GET_LOCAL_2),
										byte(bytecode.DUP),
										byte(bytecode.SET_IVAR_1),
										byte(bytecode.POP),
										byte(bytecode.GET_LOCAL_3),
										byte(bytecode.DUP),
										byte(bytecode.SET_IVAR_0),
										byte(bytecode.POP),
										byte(bytecode.NIL),
										byte(bytecode.POP),
										byte(bytecode.RETURN_SELF),
									},
									L(P(233, 10, 6), P(262, 10, 35)),
									bytecode.LineInfoList{
										bytecode.NewLineInfo(10, 15),
									},
									3,
									0,
									nil,
								)),
								value.ToSymbol("#init").ToValue(),
								value.ToSymbol("blue").ToValue(),
								value.ToSymbol("green").ToValue(),
								value.ToSymbol("red").ToValue(),
								value.ToSymbol("Std::Kernel").ToValue(),
								value.Ref(set(&baz, vm.NewBytecodeFunction(
									nil,
									value.ToSymbol("Std::Kernel::baz"),
									[]byte{
										byte(bytecode.SELF),
										byte(bytecode.LOAD_VALUE_0),
										byte(bytecode.GET_LOCAL_1),
										byte(bytecode.CALL_METHOD8), 1,
										byte(bytecode.NEW_STRING8), 2,
										byte(bytecode.CALL_METHOD_NT8), 2,
										byte(bytecode.RETURN),
									},
									L(P(277, 13, 5), P(328, 15, 7)),
									bytecode.LineInfoList{
										bytecode.NewLineInfo(14, 9),
										bytecode.NewLineInfo(15, 1),
									},
									1,
									0,
									[]value.Value{
										value.Ref(value.String("col: ")),
										value.Ref(vm.NewCallSiteInfo(value.ToSymbol("inspect"), 0)),
										value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "println@1"), 1)),
									},
								))),
								value.ToSymbol("baz").ToValue(),
							},
						)),
						value.ToSymbol("Color").ToValue(),
						value.ToSymbol("RED").ToValue(),
						value.ToSymbol("GREEN").ToValue(),
						value.ToSymbol("BLUE").ToValue(),
						value.ToSymbol("Std::Kernel").ToValue(),
						value.ToSymbol("Color::RED").ToValue(),
						value.Ref(vm.NewBytecodeCallSiteInfo(baz, 1, false)),
					},
				)
			},
		},
		"infer namespace in local declaration": {
			input: `
				class Color
					const RED: Color = new(255u8, 0u8, 0u8)
					const GREEN: Color = new(0u8, 255u8, 0u8)
					const BLUE: Color = new(0u8, 0u8, 255u8)

					getter red: UInt8
					getter green: UInt8
					getter blue: UInt8
					init(@red, @green, @blue); end
				end

				var baz: Color = .::RED
			`,
			want: vm.NewBytecodeFunctionNoParams(
				nil,
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
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_UINT8), 255,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.INSTANTIATE8), 3,
					byte(bytecode.DEF_CONST),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.LOAD_UINT8), 255,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.INSTANTIATE8), 3,
					byte(bytecode.DEF_CONST),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.LOAD_UINT8), 255,
					byte(bytecode.INSTANTIATE8), 3,
					byte(bytecode.DEF_CONST),
					byte(bytecode.GET_CONST8), 7,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(300, 13, 28)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(0, 2),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(0, 2),
					bytecode.NewLineInfo(4, 9),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(0, 2),
					bytecode.NewLineInfo(5, 9),
					bytecode.NewLineInfo(13, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						nil,
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
						L(P(0, 1, 1), P(300, 13, 28)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(13, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Color").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						nil,
						ivarIndicesSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(300, 13, 28)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(13, 2),
						},
						[]value.Value{
							value.ToSymbol("Color").ToValue(),
							value.Ref(&value.IvarIndices{
								symbol.ToSymbol("blue"):  0,
								symbol.ToSymbol("green"): 1,
								symbol.ToSymbol("red"):   2,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						nil,
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.INT_0),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.INT_1),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE8), 5,
							byte(bytecode.INT_2),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(300, 13, 28)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 17),
							bytecode.NewLineInfo(13, 2),
						},
						[]value.Value{
							value.ToSymbol("Color").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								nil,
								value.ToSymbol("Color.:#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_2),
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_1),
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL_3),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(233, 10, 6), P(262, 10, 35)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(10, 15),
								},
								3,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.ToSymbol("blue").ToValue(),
							value.ToSymbol("green").ToValue(),
							value.ToSymbol("red").ToValue(),
						},
					)),
					value.ToSymbol("Color").ToValue(),
					value.ToSymbol("RED").ToValue(),
					value.ToSymbol("GREEN").ToValue(),
					value.ToSymbol("BLUE").ToValue(),
					value.ToSymbol("Color::RED").ToValue(),
				},
			),
		},
		"infer namespace in switch pattern": {
			input: `
				class Color
					const RED: Color = new(255u8, 0u8, 0u8)
					const GREEN: Color = new(0u8, 255u8, 0u8)
					const BLUE: Color = new(0u8, 0u8, 255u8)

					getter red: UInt8
					getter green: UInt8
					getter blue: UInt8
					init(@red, @green, @blue); end
				end

				baz := Color::RED
				switch baz
				case .::GREEN
					println("green")
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				nil,
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
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_UINT8), 255,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.INSTANTIATE8), 3,
					byte(bytecode.DEF_CONST),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.LOAD_UINT8), 255,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.INSTANTIATE8), 3,
					byte(bytecode.DEF_CONST),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.LOAD_UINT8), 0,
					byte(bytecode.LOAD_UINT8), 255,
					byte(bytecode.INSTANTIATE8), 3,
					byte(bytecode.DEF_CONST),
					byte(bytecode.GET_CONST8), 7,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 8,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 9,
					byte(bytecode.LOAD_VALUE8), 10,
					byte(bytecode.CALL_METHOD_NT8), 11,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(357, 17, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(0, 2),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(0, 2),
					bytecode.NewLineInfo(4, 9),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(0, 2),
					bytecode.NewLineInfo(5, 9),
					bytecode.NewLineInfo(13, 3),
					bytecode.NewLineInfo(14, 1),
					bytecode.NewLineInfo(15, 8),
					bytecode.NewLineInfo(16, 9),
					bytecode.NewLineInfo(14, 1),
					bytecode.NewLineInfo(17, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						nil,
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
						L(P(0, 1, 1), P(357, 17, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
							bytecode.NewLineInfo(17, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Color").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						nil,
						ivarIndicesSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.DEF_IVARS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(357, 17, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 4),
							bytecode.NewLineInfo(17, 2),
						},
						[]value.Value{
							value.ToSymbol("Color").ToValue(),
							value.Ref(&value.IvarIndices{
								symbol.ToSymbol("blue"):  0,
								symbol.ToSymbol("green"): 1,
								symbol.ToSymbol("red"):   2,
							}),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						nil,
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE_3),
							byte(bytecode.INT_0),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.INT_1),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.LOAD_VALUE8), 5,
							byte(bytecode.INT_2),
							byte(bytecode.DEF_GETTER),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(357, 17, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 17),
							bytecode.NewLineInfo(17, 2),
						},
						[]value.Value{
							value.ToSymbol("Color").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								nil,
								value.ToSymbol("Color.:#init"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_2),
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL_2),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_1),
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL_3),
									byte(bytecode.DUP),
									byte(bytecode.SET_IVAR_0),
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(233, 10, 6), P(262, 10, 35)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(10, 15),
								},
								3,
								0,
								nil,
							)),
							value.ToSymbol("#init").ToValue(),
							value.ToSymbol("blue").ToValue(),
							value.ToSymbol("green").ToValue(),
							value.ToSymbol("red").ToValue(),
						},
					)),
					value.ToSymbol("Color").ToValue(),
					value.ToSymbol("RED").ToValue(),
					value.ToSymbol("GREEN").ToValue(),
					value.ToSymbol("BLUE").ToValue(),
					value.ToSymbol("Color::RED").ToValue(),
					value.ToSymbol("Color::GREEN").ToValue(),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.String("green")),
					value.Ref(vm.NewNativeCallSiteInfo(nativeMethodStr(value.KernelModule.SingletonClass(), "println@1"), 1)),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bytecodeCompilerTest(tc, t)
		})
	}
}

func TestBytecodeDefConstant(t *testing.T) {
	tests := bytecodeTestTable{
		"relative path Foo": {
			input: "const Foo = 3",
			want: vm.NewBytecodeFunctionNoParams(
				nil,
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.INT_3),
					byte(bytecode.DEF_CONST),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.ToSymbol("Root").ToValue(),
					value.ToSymbol("Foo").ToValue(),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bytecodeCompilerTest(tc, t)
		})
	}
}
