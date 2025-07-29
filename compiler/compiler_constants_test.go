package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestGetConstant(t *testing.T) {
	tests := testTable{
		"absolute path ::Std": {
			input: "::Std",
			want: vm.NewBytecodeFunctionNoParams(
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
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.INT_3),
					byte(bytecode.DEF_CONST),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(1, 1),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
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
					value.Undefined,
					value.ToSymbol("Foo").ToValue(),
					value.ToSymbol("BAR").ToValue(),
					value.Ref(vm.NewBytecodeFunctionNoParams(
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
								value.ToSymbol("<module: Foo::Baz>"),
								[]byte{
									byte(bytecode.GET_CONST8), 0,
									byte(bytecode.UNDEFINED),
									byte(bytecode.GET_CONST8), 1,
									byte(bytecode.NEW_ARRAY_TUPLE8), 1,
									byte(bytecode.CALL_METHOD8), 2,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.RETURN),
								},
								L(P(40, 4, 6), P(76, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 9),
									bytecode.NewLineInfo(6, 3),
								},
								[]value.Value{
									value.ToSymbol("Std::Kernel").ToValue(),
									value.ToSymbol("Foo::BAR").ToValue(),
									value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
								},
							)),
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

func TestDefConstant(t *testing.T) {
	tests := testTable{
		"relative path Foo": {
			input: "const Foo = 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_2),
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
					value.Undefined,
					value.ToSymbol("Root").ToValue(),
					value.ToSymbol("Foo").ToValue(),
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
