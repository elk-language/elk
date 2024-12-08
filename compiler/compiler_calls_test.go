package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestSubscript(t *testing.T) {
	tests := testTable{
		"static": {
			input: "[5, 3][0]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(5).ToValue(),
				},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 3, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 6),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayList{
						value.SmallInt(5).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.SmallInt(1).ToValue(),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_IF_NIL), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 3, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 17),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
				},
			),
		},
		"setter": {
			input: `
				arr := [5, 3]
				arr[1] = 15
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 8),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(15),
				},
			),
		},
		"increment": {
			input: `
				arr := [5, 3]
				arr[1]++
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DUP_N), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.INCREMENT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(31, 3, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
				},
			),
		},
		"decrement": {
			input: `
				arr := [5, 3]
				arr[1]--
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DUP_N), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DECREMENT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(31, 3, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
				},
			),
		},

		"add": {
			input: `
				arr := [5, 3]
				arr[1] += 15
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.ADD),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(15),
				},
			),
		},
		"subtract": {
			input: `
				arr := [5, 3]
				arr[1] -= 2
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.SUBTRACT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(2),
				},
			),
		},
		"multiply": {
			input: `
				arr := [5, 3]
				arr[1] *= 3
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"divide": {
			input: `
				arr := [5, 3]
				arr[1] /= 10
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.DIVIDE),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(10),
				},
			),
		},
		"exponentiate": {
			input: `
				arr := [5, 3]
				arr[1] **= 12
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(36, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(12),
				},
			),
		},

		"modulo": {
			input: `
				arr := [5, 3]
				arr[1] %= 2
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.MODULO),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(2),
				},
			),
		},
		"bitwise AND": {
			input: `
				arr := [5, 3]
				arr[1] &= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(8),
				},
			),
		},
		"bitwise OR": {
			input: `
				arr := [5, 3]
				arr[1] |= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(8),
				},
			),
		},
		"bitwise XOR": {
			input: `
				arr := [5, 3]
				arr[1] ^= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(8),
				},
			),
		},
		"left bitshift": {
			input: `
				arr := [5, 3]
				arr[1] <<= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(8),
				},
			),
		},
		"right bitshift": {
			input: `
				arr := [5, 3]
				arr[1] >>= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(8),
				},
			),
		},

		"left logical bitshift": {
			input: `
				arr := [5u64, 3u64]
				arr[1] <<<= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.UInt64(5),
						value.UInt64(3),
					},
					value.SmallInt(1),
					value.SmallInt(8),
				},
			),
		},
		"right logical bitshift": {
			input: `
				arr := [5u64, 3u64]
				arr[1] >>>= 8
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.UInt64(5),
						value.UInt64(3),
					},
					value.SmallInt(1),
					value.SmallInt(8),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.JUMP_IF), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(50, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 17),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(8),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.JUMP_UNLESS), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(50, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 17),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(8),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.JUMP_IF_NIL), 0, 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(50, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 20),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.SmallInt(8),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(31, 3, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(31, 3, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(3, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					nil,
					value.ToSymbol("Foo"),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(91, 7, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(7, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
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
							bytecode.NewLineInfo(1, 23),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Bar"),
							value.ToSymbol("Foo::Bar"),
							value.ToSymbol("Baz"),
							value.ToSymbol("Foo::Bar::Baz"),
							value.ToSymbol("Std::Object"),
						},
					),
					nil,
					value.ToSymbol("Foo::Bar::Baz"),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.INSTANTIATE8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 5, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 9),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(76, 5, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(76, 5, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("#init"),
						},
					),
					value.ToSymbol("Foo"),
					value.SmallInt(1),
					value.String("lol"),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.INSTANTIATE8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(79, 5, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 9),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(79, 5, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(79, 5, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("#init"),
						},
					),
					value.ToSymbol("Foo"),
					value.SmallInt(1),
					value.String("lol"),
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
			err: error.ErrorList{
				error.NewFailure(
					L(P(83, 5, 27), P(86, 5, 30)),
					"duplicated argument `b` in call to `#init`",
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(53, 5, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(53, 5, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(53, 5, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Foo"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.CALL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(57, 5, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(57, 5, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(57, 5, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("call"),
						},
					),
					value.ToSymbol("Foo"),
					value.NewCallSiteInfo(value.ToSymbol("call"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.CALL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(53, 5, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(53, 5, 11)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(53, 5, 11)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("call"),
						},
					),
					value.ToSymbol("Foo"),
					value.NewCallSiteInfo(value.ToSymbol("call"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.CALL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(55, 5, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(55, 5, 13)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(55, 5, 13)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("call"),
						},
					),
					value.ToSymbol("Foo"),
					value.NewCallSiteInfo(value.ToSymbol("call"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
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
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 7),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(54, 5, 13)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(54, 5, 13)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Foo"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.JUMP_IF_NIL), 0, 3,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(75, 6, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(6, 10),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(75, 6, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(75, 6, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_IF_NIL), 0, 2,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(74, 6, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(6, 8),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(74, 6, 11)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(74, 6, 11)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(70, 5, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 7),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 5, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 5, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Foo"),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(82, 5, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(5, 9),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(82, 5, 22)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(82, 5, 22)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(5, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Foo"),
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 2),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_IF_NIL), 0, 6,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(103, 6, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(6, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(103, 6, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(103, 6, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 2),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(93, 6, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 9),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(93, 6, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(93, 6, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Foo"),
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 2),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(96, 6, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 9),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(96, 6, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Foo"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(96, 6, 23)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Foo"),
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 2),
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
			err: error.ErrorList{
				error.NewFailure(
					L(P(89, 5, 29), P(92, 5, 32)),
					"duplicated argument `b` in call to `foo`",
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(92, 8, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(92, 8, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(92, 8, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 6, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(63, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(63, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("call"),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("call"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(78, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(78, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
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
							bytecode.NewLineInfo(1, 18),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("call"),
							value.ToSymbol("Std::Kernel"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("a"),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("call"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(70, 6, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 6, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(70, 6, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(94, 8, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 1,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(94, 8, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
							value.ToSymbol("Std::Object"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(94, 8, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(27, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(27, 3, 10)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(3, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel"),
							vm.NewBytecodeFunctionNoParams(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Std::Kernel"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(54, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(54, 3, 18)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(3, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Std::Kernel"),
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 2),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(57, 3, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(57, 3, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(3, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo"),
						},
					),
					value.ToSymbol("Std::Kernel"),
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 2),
				},
			),
		},
		"call a function with duplicated named args": {
			input: `
				def foo(a: String, b: Int); end
				foo(b: 1, a: 'lol', b: 2)
			`,
			err: error.ErrorList{
				error.NewFailure(
					L(P(61, 3, 25), P(64, 3, 28)),
					"duplicated argument `b` in call to `foo`",
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(81, 6, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 7),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(81, 6, 14)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(81, 6, 14)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 9),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.INCREMENT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(104, 7, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 8),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(104, 7, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(104, 7, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.DECREMENT),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(104, 7, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 8),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(104, 7, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(104, 7, 12)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.ADD),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SUBTRACT),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.MULTIPLY),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.DIVIDE),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(108, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.MODULO),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(108, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(119, 7, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(119, 7, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(119, 7, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(44, 3, 29)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.Int64(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.Int64(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(108, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(108, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(119, 7, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(119, 7, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(119, 7, 20)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(44, 3, 29)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.Int64(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.Int64(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(107, 7, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(107, 7, 15)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(39, 3, 24)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.JUMP_IF), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(110, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 15),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(40, 3, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.JUMP_UNLESS), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(110, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 15),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(40, 3, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.JUMP_IF_NIL), 0, 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(110, 7, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 18),
				},
				[]value.Value{
					vm.NewBytecodeFunctionNoParams(
						namespaceDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root"),
							value.ToSymbol("Bar"),
						},
					),
					vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.LOAD_VALUE8), 3,
							byte(bytecode.LOAD_VALUE8), 4,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(110, 7, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 14),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Bar"),
							vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(40, 3, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 3),
								},
								0,
								0,
								[]value.Value{
									value.SmallInt(3),
								},
							),
							value.ToSymbol("foo"),
							vm.NewBytecodeFunction(
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
							),
							value.ToSymbol("foo="),
						},
					),
					value.ToSymbol("Bar"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1),
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
