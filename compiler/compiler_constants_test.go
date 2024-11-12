package compiler

// import (
// 	"testing"

// 	"github.com/elk-language/elk/bytecode"
// 	"github.com/elk-language/elk/value"
// 	"github.com/elk-language/elk/vm"
// )

// func TestGetModuleConstant(t *testing.T) {
// 	tests := testTable{
// 		"absolute path ::Std": {
// 			input: "::Std",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.ROOT),
// 					byte(bytecode.GET_MOD_CONST8), 0,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(4, 1, 5)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 4),
// 				},
// 				[]value.Value{
// 					value.ToSymbol("Std"),
// 				},
// 			),
// 		},
// 		"absolute nested path ::Std::Float::INF": {
// 			input: "::Std::Float::INF",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.ROOT),
// 					byte(bytecode.GET_MOD_CONST8), 0,
// 					byte(bytecode.GET_MOD_CONST8), 1,
// 					byte(bytecode.GET_MOD_CONST8), 2,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(16, 1, 17)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 8),
// 				},
// 				[]value.Value{
// 					value.ToSymbol("Std"),
// 					value.ToSymbol("Float"),
// 					value.ToSymbol("INF"),
// 				},
// 			),
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			compilerTest(tc, t)
// 		})
// 	}
// }

// // TODO: fix
// // func TestDefModuleConstant(t *testing.T) {
// // 	tests := testTable{
// // 		"relative path Foo": {
// // 			input: "const Foo = 3",
// // 			want: vm.NewBytecodeFunctionNoParams(
// // 				mainSymbol,
// // 				[]byte{
// // 					byte(bytecode.LOAD_VALUE8), 0,
// // 					byte(bytecode.CONSTANT_CONTAINER),
// // 					byte(bytecode.DEF_MOD_CONST8), 1,
// // 					byte(bytecode.RETURN),
// // 				},
// // 				L(P(0, 1, 1), P(12, 1, 13)),
// // 				bytecode.LineInfoList{
// // 					bytecode.NewLineInfo(1, 6),
// // 				},
// // 				[]value.Value{
// // 					value.SmallInt(3),
// // 					value.ToSymbol("Foo"),
// // 				},
// // 			),
// // 		},
// // 	}

// // 	for name, tc := range tests {
// // 		t.Run(name, func(t *testing.T) {
// // 			compilerTest(tc, t)
// // 		})
// // 	}
// // }
