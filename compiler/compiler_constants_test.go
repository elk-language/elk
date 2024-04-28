package compiler

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestGetModuleConstant(t *testing.T) {
	tests := testTable{
		"absolute path ::Std": {
			input: "::Std",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.ToSymbol("Std"),
				},
			),
		},
		"absolute nested path ::Std::Float::INF": {
			input: "::Std::Float::INF",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.GET_MOD_CONST8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.ToSymbol("Std"),
					value.ToSymbol("Float"),
					value.ToSymbol("INF"),
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

func TestDefModuleConstant(t *testing.T) {
	tests := testTable{
		"relative path Foo": {
			input: "Foo := 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.DEF_MOD_CONST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.SmallInt(3),
					value.ToSymbol("Foo"),
				},
			),
		},
		"absolute path ::Foo": {
			input: "::Foo := 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.DEF_MOD_CONST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.SmallInt(3),
					value.ToSymbol("Foo"),
				},
			),
		},
		"absolute nested path ::Std::Float::Foo": {
			input: "::Std::Float::Foo := 'bar'",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.GET_MOD_CONST8), 2,
					byte(bytecode.DEF_MOD_CONST8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 1, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.String("bar"),
					value.ToSymbol("Std"),
					value.ToSymbol("Float"),
					value.ToSymbol("Foo"),
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
