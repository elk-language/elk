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
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.ToSymbol("Std"),
				},
			),
		},
		"absolute nested path ::Std::Float::INF": {
			input: "::Std::Float::INF",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.ToSymbol("Std::Float::INF"),
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
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(37, 3, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 3),
				},
				[]value.Value{
					nil,
					value.ToSymbol("Std::Float::INF"),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.DEF_CONST),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.INIT_NAMESPACE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 12),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Baz"),
						},
					),
					nil,
					value.ToSymbol("Foo"),
					value.ToSymbol("BAR"),
					value.SmallInt(3),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<module: Foo>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.INIT_NAMESPACE),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(84, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(4, 5),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							value.ToSymbol("Foo::Baz"),
							vm.NewBytecodeFunctionNoParams(
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
									value.ToSymbol("Std::Kernel"),
									value.ToSymbol("Foo::BAR"),
									value.NewCallSiteInfo(value.ToSymbol("println"), 1),
								},
							),
						},
					),
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
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.DEF_CONST),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					nil,
					value.ToSymbol("Root"),
					value.ToSymbol("Foo"),
					value.SmallInt(3),
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
