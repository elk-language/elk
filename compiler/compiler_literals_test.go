package compiler_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestClosureLiteral(t *testing.T) {
	tests := testTable{
		"recursive closure": {
			input: `
				var calc_fib: |n: Int|: Int = |n| ->
					return 1 if n < 3

					calc_fib(n - 2) + calc_fib(n - 1)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.CLOSURE), 2, 1, 0xff,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(112, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionWithUpvalues(
						functionSymbol,
						[]byte{
							byte(bytecode.GET_LOCAL_1),
							byte(bytecode.INT_3),
							byte(bytecode.JUMP_UNLESS_ILT), 0, 2,
							byte(bytecode.INT_1),
							byte(bytecode.RETURN),
							byte(bytecode.GET_UPVALUE_0),
							byte(bytecode.GET_LOCAL_1),
							byte(bytecode.INT_2),
							byte(bytecode.SUBTRACT_INT),
							byte(bytecode.CALL8), 0,
							byte(bytecode.GET_UPVALUE_0),
							byte(bytecode.GET_LOCAL_1),
							byte(bytecode.INT_1),
							byte(bytecode.SUBTRACT_INT),
							byte(bytecode.CALL8), 1,
							byte(bytecode.ADD_INT),
							byte(bytecode.RETURN),
						},
						L(P(35, 2, 35), P(111, 6, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 7),
							bytecode.NewLineInfo(5, 13),
							bytecode.NewLineInfo(6, 1),
						},
						1,
						0,
						[]value.Value{
							value.Ref(value.NewCallSiteInfo(value.ToSymbol("call"), 1)),
							value.Ref(value.NewCallSiteInfo(value.ToSymbol("call"), 1)),
						},
						1,
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

func TestStringLiteral(t *testing.T) {
	tests := testTable{
		"static string": {
			input: `"foo bar"`,
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
					value.Ref(value.String("foo bar")),
				},
			),
		},
		"interpolated string": {
			input: `
				bar := 15.2
				foo := 1
				"foo: ${foo + 2}, bar: $bar"
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_STRING8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(62, 4, 33)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 9),
				},
				[]value.Value{
					value.Float(15.2).ToValue(),
					value.Ref(value.String("foo: ")),
					value.Ref(value.String(", bar: ")),
				},
			),
		},
		"inspect interpolated string": {
			input: `
				bar := 15.2
				foo := 1
				"foo: #{foo + 2}, bar: #bar"
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.ADD_INT),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.NEW_STRING8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(62, 4, 33)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 14),
				},
				[]value.Value{
					value.Float(15.2).ToValue(),
					value.Ref(value.String("foo: ")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("inspect"), 0)),
					value.Ref(value.String(", bar: ")),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("inspect"), 0)),
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

func TestRangeLiteral(t *testing.T) {
	tests := testTable{
		"static closed range": {
			input: `2...5`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.NewClosedRange(value.SmallInt(2).ToValue(), value.SmallInt(5).ToValue())),
				},
			),
		},
		"static open range": {
			input: `2<.<5`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.NewOpenRange(value.SmallInt(2).ToValue(), value.SmallInt(5).ToValue())),
				},
			),
		},
		"static left open range": {
			input: `2<..5`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.NewLeftOpenRange(value.SmallInt(2).ToValue(), value.SmallInt(5).ToValue())),
				},
			),
		},
		"static right open range": {
			input: `2..<5`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.NewRightOpenRange(value.SmallInt(2).ToValue(), value.SmallInt(5).ToValue())),
				},
			),
		},
		"static beginless closed range": {
			input: `...5`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.NewBeginlessClosedRange(value.SmallInt(5).ToValue())),
				},
			),
		},
		"static beginless open range": {
			input: `..<5`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.NewBeginlessOpenRange(value.SmallInt(5).ToValue())),
				},
			),
		},
		"static endless closed range": {
			input: `2...`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.NewEndlessClosedRange(value.SmallInt(2).ToValue())),
				},
			),
		},
		"static endless open range": {
			input: `2<..`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.NewEndlessOpenRange(value.SmallInt(2).ToValue())),
				},
			),
		},
		"closed range": {
			input: `
			  a := 2
				a...5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.NEW_RANGE), bytecode.CLOSED_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{},
			),
		},
		"open range": {
			input: `
			  a := 2
				a<.<5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.NEW_RANGE), bytecode.OPEN_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{},
			),
		},
		"left open range": {
			input: `
			  a := 2
				a<..5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.NEW_RANGE), bytecode.LEFT_OPEN_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{},
			),
		},
		"right open range": {
			input: `
			  a := 2
				a..<5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.NEW_RANGE), bytecode.RIGHT_OPEN_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{},
			),
		},
		"beginless closed range": {
			input: `
			  a := 2
				...a
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_RANGE), bytecode.BEGINLESS_CLOSED_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 3, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{},
			),
		},
		"beginless open range": {
			input: `
			  a := 2
				..<a
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_RANGE), bytecode.BEGINLESS_OPEN_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 3, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{},
			),
		},
		"endless closed range": {
			input: `
			  a := 2
				a...
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_RANGE), bytecode.ENDLESS_CLOSED_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 3, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{},
			),
		},
		"endless open range": {
			input: `
			  a := 2
				a<..
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_RANGE), bytecode.ENDLESS_OPEN_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 3, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
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

func TestLiterals(t *testing.T) {
	tests := testTable{
		"put UInt8": {
			input: "1u8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_UINT8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"put UInt16": {
			input: "25u16",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_UINT16_8), 25,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"put UInt32": {
			input: "450_200u32",
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
					value.UInt32(450200).ToValue(),
				},
			),
		},
		"put UInt64": {
			input: "450_200u64",
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
					value.UInt64(450200).ToValue(),
				},
			),
		},
		"put Int8": {
			input: "1i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"put Int16": {
			input: "25i16",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT16_8), 25,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"put Int32": {
			input: "450_200i32",
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
					value.Int32(450200).ToValue(),
				},
			),
		},
		"put Int64": {
			input: "450_200i64",
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
					value.Int64(450200).ToValue(),
				},
			),
		},
		"put SmallInt": {
			input: "450_200",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.SmallInt(450200).ToValue(),
				},
			),
		},
		"put BigInt": {
			input: (&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5)).String(),
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(
					P(0, 1, 1),
					P(
						len((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5)).String())-1,
						1,
						len((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5)).String()),
					),
				),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5)))),
				},
			),
		},
		"put Float64": {
			input: "45.5f64",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Float64(45.5).ToValue(),
				},
			),
		},
		"put Float32": {
			input: "45.5f32",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Float32(45.5).ToValue(),
				},
			),
		},
		"put Float": {
			input: "45.5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Float(45.5).ToValue(),
				},
			),
		},
		"put Raw String": {
			input: `'foo\n'`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.String(`foo\n`)),
				},
			),
		},
		"put String": {
			input: `"foo\n"`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.String("foo\n")),
				},
			),
		},
		"put raw Char": {
			input: "`I`",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_CHAR_8), 'I',
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"put Char": {
			input: "`\\n`",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_CHAR_8), '\n',
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{},
			),
		},
		"put nil": {
			input: `nil`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
		},
		"put true": {
			input: `true`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
		},
		"put false": {
			input: `false`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
		},
		"put simple Symbol": {
			input: `:foo`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.ToSymbol("foo").ToValue(),
				},
			),
		},
		"put self": {
			input: `self`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
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

func TestArrayTuples(t *testing.T) {
	tests := testTable{
		"empty arrayTuple": {
			input: "%[]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{}),
				},
			),
		},
		"with static elements": {
			input: "%[1, 'foo', 5, 5.6]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.SmallInt(5).ToValue(),
						value.Float(5.6).ToValue(),
					}),
				},
			),
		},
		"with static keyed elements": {
			input: "%[1, 'foo', 5 => 5,  3 => 5.6, :lol]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.Nil,
						value.Float(5.6).ToValue(),
						value.Nil,
						value.SmallInt(5).ToValue(),
						value.ToSymbol("lol").ToValue(),
					}),
				},
			),
		},
		"nested static arrayTuples": {
			input: "%[1, %['bar', %[7.2]]]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 1, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.Ref(&value.ArrayTuple{
							value.Ref(value.String("bar")),
							value.Ref(&value.ArrayTuple{
								value.Float(7.2).ToValue(),
							}),
						}),
					}),
				},
			),
		},
		"nested static with mutable elements": {
			input: "%[1, %['bar', [7.2]]]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.COPY),
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 1, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
					}),
					value.Ref(&value.ArrayTuple{
						value.Ref(value.String("bar")),
					}),
					value.Ref(&value.ArrayList{
						value.Float(7.2).ToValue(),
					}),
				},
			),
		},
		"with static keyed and dynamic elements": {
			input: "%[1, 'foo', 5 => 5,  3 => 5.6, String.name]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 1, 43)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.Nil,
						value.Float(5.6).ToValue(),
						value.Nil,
						value.SmallInt(5).ToValue(),
					}),
					value.ToSymbol("Std::String").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("name"),
						0,
					)),
				},
			),
		},
		"with static and dynamic elements": {
			input: "%[1, 'foo', 5, Object(), 5, %[:foo]]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.INT_5),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.NEW_ARRAY_TUPLE8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.SmallInt(5).ToValue(),
					}),
					value.ToSymbol("Std::Object").ToValue(),
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with dynamic elements": {
			input: "%[Object(), 5, %[:foo]]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.INT_5),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.NEW_ARRAY_TUPLE8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 1, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.ToSymbol("Std::Object").ToValue(),
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with static elements and if modifiers": {
			input: `
				var a: Object? = Object()
				%[1, 5 if a, %[:foo]]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.INT_5),
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(56, 3, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					value.ToSymbol("Std::Object").ToValue(),
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
					}),
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},

		"with static elements and unless modifiers": {
			input: `
				var a: Object? = nil
				%[1, 5 unless a, %[:foo]]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.INT_5),
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(55, 3, 30)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
					}),
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with static elements and for in loops": {
			input: `
				%[1, i * 2 for i in [1, 2, 3], %[:foo]]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 8,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.MULTIPLY_INT),
					byte(bytecode.APPEND),
					byte(bytecode.LOOP), 0, 12,
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(44, 2, 44)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 24),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
					}),
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with dynamic elements and if modifiers": {
			input: `
				var a: Object? = nil
				%[String.name, 5 if a, %[:foo]]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.INT_5),
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 3, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 19),
				},
				[]value.Value{
					value.ToSymbol("Std::String").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("name"),
						0,
					)),
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with dynamic and keyed elements": {
			input: "%[Object(), 1, 'foo', 5 => 5,  3 => 5.6]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.NEW_ARRAY_TUPLE8), 3,
					byte(bytecode.INT_5),
					byte(bytecode.INT_5),
					byte(bytecode.APPEND_AT),
					byte(bytecode.INT_3),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.APPEND_AT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(39, 1, 40)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 16),
				},
				[]value.Value{
					value.ToSymbol("Std::Object").ToValue(),
					value.Ref(value.String("foo")),
					value.Float(5.6).ToValue(),
				},
			),
		},
		"with keyed and if elements": {
			input: `
				var a: String? = nil
				%[3 => 5 if a]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_ARRAY_TUPLE8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.INT_3),
					byte(bytecode.INT_5),
					byte(bytecode.APPEND_AT),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(44, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{},
			),
		},

		"with static concat": {
			input: "%[1, 2, 3] + %[4, 5, 6] + %[10]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(30, 1, 31)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
						value.SmallInt(6).ToValue(),
						value.SmallInt(10).ToValue(),
					}),
				},
			),
		},
		"with static concat with list": {
			input: "%[1, 2, 3] + [4, 5, 6] + %[10]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 1, 30)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
						value.SmallInt(6).ToValue(),
						value.SmallInt(10).ToValue(),
					}),
				},
			),
		},
		"with static repeat": {
			input: "%[1, 2, 3] * 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"with static concat and nested tuples": {
			input: "%[1, 2, 3] + %[4, 5, 6, %[7, 8]] + %[10]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(39, 1, 40)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
						value.SmallInt(6).ToValue(),
						value.Ref(&value.ArrayTuple{
							value.SmallInt(7).ToValue(),
							value.SmallInt(8).ToValue(),
						}),
						value.SmallInt(10).ToValue(),
					}),
				},
			),
		},
		"word arrayTuple": {
			input: `%w[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.Ref(value.String("foo")),
						value.Ref(value.String("bar")),
						value.Ref(value.String("baz")),
					}),
				},
			),
		},
		"symbol arrayTuple": {
			input: `%s[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
						value.ToSymbol("bar").ToValue(),
						value.ToSymbol("baz").ToValue(),
					}),
				},
			),
		},
		"hex arrayTuple": {
			input: `%x[ab cd 5f]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(0xab).ToValue(),
						value.SmallInt(0xcd).ToValue(),
						value.SmallInt(0x5f).ToValue(),
					}),
				},
			),
		},
		"bin arrayTuple": {
			input: `%b[101 11 10]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.SmallInt(0b101).ToValue(),
						value.SmallInt(0b11).ToValue(),
						value.SmallInt(0b10).ToValue(),
					}),
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

func TestArrayLists(t *testing.T) {
	tests := testTable{
		"empty list": {
			input: "[]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(1, 1, 2)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{}),
				},
			),
		},
		"with static elements": {
			input: "[1, 'foo', 5, 5.6]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.SmallInt(5).ToValue(),
						value.Float(5.6).ToValue(),
					}),
				},
			),
		},
		"with static elements and static capacity": {
			input: "[1, 'foo', 5, 5.6]:10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 1, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.SmallInt(5).ToValue(),
						value.Float(5.6).ToValue(),
					}),
				},
			),
		},

		"with static elements and dynamic capacity": {
			input: `
				cap := 2
				[1, 'foo', 5, 5.6]:cap
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.SmallInt(5).ToValue(),
						value.Float(5.6).ToValue(),
					}),
				},
			),
		},
		"word list": {
			input: `\w[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.Ref(value.String("foo")),
						value.Ref(value.String("bar")),
						value.Ref(value.String("baz")),
					}),
				},
			),
		},
		"word list with capacity": {
			input: `\w[foo bar baz]:15`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 15,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.Ref(value.String("foo")),
						value.Ref(value.String("bar")),
						value.Ref(value.String("baz")),
					}),
				},
			),
		},
		"symbol list": {
			input: `\s[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
						value.ToSymbol("bar").ToValue(),
						value.ToSymbol("baz").ToValue(),
					}),
				},
			),
		},

		"symbol list with capacity": {
			input: `\s[foo bar baz]:15`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 15,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
						value.ToSymbol("bar").ToValue(),
						value.ToSymbol("baz").ToValue(),
					}),
				},
			),
		},
		"hex list": {
			input: `\x[ab cd 5f]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(0xab).ToValue(),
						value.SmallInt(0xcd).ToValue(),
						value.SmallInt(0x5f).ToValue(),
					}),
				},
			),
		},
		"hex list with capacity": {
			input: `\x[ab cd 5f]:2`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_2),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(0xab).ToValue(),
						value.SmallInt(0xcd).ToValue(),
						value.SmallInt(0x5f).ToValue(),
					}),
				},
			),
		},
		"bin list": {
			input: `\b[101 11 10]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(0b101).ToValue(),
						value.SmallInt(0b11).ToValue(),
						value.SmallInt(0b10).ToValue(),
					}),
				},
			),
		},
		"bin list with capacity": {
			input: `\b[101 11 10]:3`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_3),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(0b101).ToValue(),
						value.SmallInt(0b11).ToValue(),
						value.SmallInt(0b10).ToValue(),
					}),
				},
			),
		},

		"with static keyed elements": {
			input: "[1, 'foo', 5 => 5,  3 => 5.6, :lol]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 1, 35)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.Nil,
						value.Float(5.6).ToValue(),
						value.Nil,
						value.SmallInt(5).ToValue(),
						value.ToSymbol("lol").ToValue(),
					}),
				},
			),
		},
		"with static keyed elements and static capacity": {
			input: "[1, 'foo', 5 => 5,  3 => 5.6, :lol]:6",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 6,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(36, 1, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.Nil,
						value.Float(5.6).ToValue(),
						value.Nil,
						value.SmallInt(5).ToValue(),
						value.ToSymbol("lol").ToValue(),
					}),
				},
			),
		},
		"with static concat": {
			input: "[1, 2, 3] + [4, 5, 6] + [10]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(27, 1, 28)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
						value.SmallInt(6).ToValue(),
						value.SmallInt(10).ToValue(),
					}),
				},
			),
		},
		"with static repeat": {
			input: "[1, 2, 3] * 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},
		"with static concat and nested lists": {
			input: "[1, 2, 3] + [4, 5, 6, [7, 8]] + [10]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.NEW_ARRAY_LIST8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
						value.SmallInt(6).ToValue(),
					}),
					value.Ref(&value.ArrayList{
						value.SmallInt(7).ToValue(),
						value.SmallInt(8).ToValue(),
					}),
				},
			),
		},

		"nested static lists": {
			input: "[1, ['bar', [7.2]]]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.COPY),
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
					}),
					value.Ref(&value.ArrayList{
						value.Ref(value.String("bar")),
					}),
					value.Ref(&value.ArrayList{
						value.Float(7.2).ToValue(),
					}),
				},
			),
		},
		"with static keyed and dynamic elements": {
			input: `
				a := 5
				[1, 'foo', 5 => 5,  3 => 5.6, a]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(48, 3, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 6),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.Nil,
						value.Float(5.6).ToValue(),
						value.Nil,
						value.SmallInt(5).ToValue(),
					}),
				},
			),
		},
		"with static keyed, dynamic elements and capacity": {
			input: `
				a := 5
				[1, 'foo', 5 => 5,  3 => 5.6, a]:15
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 15,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(51, 3, 40)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.Nil,
						value.Float(5.6).ToValue(),
						value.Nil,
						value.SmallInt(5).ToValue(),
					}),
				},
			),
		},
		"with static and dynamic elements": {
			input: `
				var a: Int? = 3
				[1, 'foo', 5, a, 5, %[:foo]]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_3),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.NEW_ARRAY_LIST8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(53, 3, 33)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 8),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.SmallInt(5).ToValue(),
					}),
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with dynamic elements": {
			input: `
				a := 3
				[a, 5, [:foo]]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_3),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.NEW_ARRAY_LIST8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(30, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},

		"with static elements and if modifiers": {
			input: `
				var a: Int? = nil
				[1, 5 if a, [:foo]]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.INT_5),
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 3, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 17),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
					}),
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with static elements, if modifiers and capacity": {
			input: `
				var a: Int? = nil
				[1, 5 if a, [:foo]]:45
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(47, 3, 25), P(48, 3, 26)),
					"capacity cannot be specified in collection literals with conditional elements or loops",
				),
			},
		},
		"with static elements and unless modifiers": {
			input: `
				var a: Int? = nil
				[1, 5 unless a, [:foo]]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.INT_5),
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(50, 3, 28)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 17),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
					}),
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with static elements and for in loops": {
			input: `
				[1, i * 2 for i in [1, 2, 3], %[:foo]]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 8,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.MULTIPLY_INT),
					byte(bytecode.APPEND),
					byte(bytecode.LOOP), 0, 12,
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 2, 43)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 26),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
					}),
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with dynamic elements and if modifiers": {
			input: `
				var a: Int? = nil
				[Object(), 5 if a, [:foo]]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.INT_5),
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(53, 3, 31)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 21),
				},
				[]value.Value{
					value.ToSymbol("Std::Object").ToValue(),
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},

		"with dynamic and keyed elements": {
			input: `
				var a: Int? = nil
				[a, 1, 'foo', 5 => 5,  3 => 5.6]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.NEW_ARRAY_LIST8), 3,
					byte(bytecode.INT_5),
					byte(bytecode.INT_5),
					byte(bytecode.APPEND_AT),
					byte(bytecode.INT_3),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.APPEND_AT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(59, 3, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{}),
					value.Ref(value.String("foo")),
					value.Float(5.6).ToValue(),
				},
			),
		},
		"with dynamic, keyed elements and capacity": {
			input: `
				a := 3
				[a, 1, 'foo', 5 => 5,  3 => 5.6]:7
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_3),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 7,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.NEW_ARRAY_LIST8), 3,
					byte(bytecode.INT_5),
					byte(bytecode.INT_5),
					byte(bytecode.APPEND_AT),
					byte(bytecode.INT_3),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.APPEND_AT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(50, 3, 39)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 15),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{}),
					value.Ref(value.String("foo")),
					value.Float(5.6).ToValue(),
				},
			),
		},
		"with keyed and if elements": {
			input: `
				var a: Int? = nil
				[3 => 5 if a]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.INT_3),
					byte(bytecode.INT_5),
					byte(bytecode.APPEND_AT),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 15),
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

func TestHashSet(t *testing.T) {
	tests := testTable{
		"empty list": {
			input: "^[]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(&value.HashSet{}),
				},
			),
		},
		"with static elements": {
			input: "^[1, 'foo', 5, 5.6]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.SmallInt(5).ToValue(),
						value.Float(5.6).ToValue(),
					)),
				},
			),
		},
		"with static elements and static capacity": {
			input: "^[1, 'foo', 5, 5.6]:10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 1, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
						nil,
						4,
						1,
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.SmallInt(5).ToValue(),
						value.Float(5.6).ToValue(),
					)),
				},
			),
		},
		"with static elements and dynamic capacity": {
			input: `
				cap := 2
				^[1, 'foo', 5, 5.6]:cap
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 3, 28)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
						nil,
						4,
						1,
						value.SmallInt(1).ToValue(),
						value.Ref(value.String("foo")),
						value.SmallInt(5).ToValue(),
						value.Float(5.6).ToValue(),
					)),
				},
			),
		},

		"word set": {
			input: `^w[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.Ref(value.String("foo")),
						value.Ref(value.String("bar")),
						value.Ref(value.String("baz")),
					)),
				},
			),
		},
		"word set with capacity": {
			input: `^w[foo bar baz]:15`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 15,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.Ref(value.String("foo")),
						value.Ref(value.String("bar")),
						value.Ref(value.String("baz")),
					)),
				},
			),
		},
		"symbol set": {
			input: `^s[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.ToSymbol("foo").ToValue(),
						value.ToSymbol("bar").ToValue(),
						value.ToSymbol("baz").ToValue(),
					)),
				},
			),
		},
		"symbol set with capacity": {
			input: `^s[foo bar baz]:15`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 15,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.ToSymbol("foo").ToValue(),
						value.ToSymbol("bar").ToValue(),
						value.ToSymbol("baz").ToValue(),
					)),
				},
			),
		},
		"hex set": {
			input: `^x[ab cd 5f]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(0xab).ToValue(),
						value.SmallInt(0xcd).ToValue(),
						value.SmallInt(0x5f).ToValue(),
					)),
				},
			),
		},
		"hex set with capacity": {
			input: `^x[ab cd 5f]:2`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_2),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(0xab).ToValue(),
						value.SmallInt(0xcd).ToValue(),
						value.SmallInt(0x5f).ToValue(),
					)),
				},
			),
		},

		"bin set": {
			input: `^b[101 11 10]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(0b101).ToValue(),
						value.SmallInt(0b11).ToValue(),
						value.SmallInt(0b10).ToValue(),
					)),
				},
			),
		},
		"bin set with capacity": {
			input: `^b[101 11 10]:3`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_3),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(0b101).ToValue(),
						value.SmallInt(0b11).ToValue(),
						value.SmallInt(0b10).ToValue(),
					)),
				},
			),
		},
		"with static and dynamic elements": {
			input: `
				var a: Int? = nil
				^[1, 'foo', 5, a, 5]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.NEW_HASH_SET8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 3, 25)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithCapacityAndElements(
						nil,
						5,
						value.Ref(value.String("foo")),
						value.SmallInt(5).ToValue(),
						value.SmallInt(1).ToValue(),
					)),
				},
			),
		},
		"with dynamic elements": {
			input: `
				var a: Int? = nil
				^[a, 5]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.NEW_HASH_SET8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithCapacityAndElements(
						nil,
						2,
					)),
				},
			),
		},
		"with static elements and if modifiers": {
			input: `
				var a: Int? = nil
				^[1, 5 if a]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.INT_5),
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(39, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithCapacityAndElements(
						nil,
						2,
						value.SmallInt(1).ToValue(),
					)),
				},
			),
		},
		"with static elements, if modifiers and capacity": {
			input: `
				var a: Int? = nil
				^[1, 5 if a]:45
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(40, 3, 18), P(41, 3, 19)),
					"capacity cannot be specified in collection literals with conditional elements or loops",
				),
			},
		},
		"with static elements and unless modifiers": {
			input: `
				var a: Int? = nil
				^[1, 5 unless a]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.INT_5),
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 3, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithCapacityAndElements(
						nil,
						2,
						value.SmallInt(1).ToValue(),
					)),
				},
			),
		},
		"with static elements and for in loops": {
			input: `
				^[1, i * 2 for i in [1, 2, 3], 2]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 8,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.MULTIPLY_INT),
					byte(bytecode.APPEND),
					byte(bytecode.LOOP), 0, 12,
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.INT_2),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(38, 2, 38)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 26),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithCapacityAndElements(
						nil,
						3,
						value.SmallInt(1).ToValue(),
					)),
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
				},
			),
		},

		"with dynamic elements and if modifiers": {
			input: `
				var a: Int? = nil
				^[Object(), 5 if a]
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.NEW_HASH_SET8), 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.INT_5),
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 3, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 18),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashSetWithCapacityAndElements(
						nil,
						2,
					)),
					value.ToSymbol("Std::Object").ToValue(),
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

func TestHashMap(t *testing.T) {
	tests := testTable{
		"empty": {
			input: "{}",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(1, 1, 2)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(value.NewHashMap(0)),
				},
			),
		},
		"shorthand local": {
			input: `
				foo := 3
				{ foo }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_3),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_HASH_MAP8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 3, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(value.NewHashMap(0)),
					value.ToSymbol("foo").ToValue(),
				},
			),
		},
		"shorthand private local": {
			input: `
				_foo := 3
				{ _foo }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_3),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_HASH_MAP8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(27, 3, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(value.NewHashMap(0)),
					value.ToSymbol("_foo").ToValue(),
				},
			),
		},
		"with static elements": {
			input: `{ 1 => 'foo', foo: 5, "bar" => 5.6 }`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashMapWithElements(
						nil,
						value.Pair{
							Key:   value.SmallInt(1).ToValue(),
							Value: value.Ref(value.String("foo")),
						},
						value.Pair{
							Key:   value.ToSymbol("foo").ToValue(),
							Value: value.SmallInt(5).ToValue(),
						},
						value.Pair{
							Key:   value.Ref(value.String("bar")),
							Value: value.Float(5.6).ToValue(),
						},
					)),
				},
			),
		},
		"with static elements and for loops": {
			input: `{ 1 => 'foo', i => i ** 2 for i in [1, 2, 3], 2 => 5.6 }`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 9,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.EXPONENTIATE_INT),
					byte(bytecode.MAP_SET),
					byte(bytecode.LOOP), 0, 13,
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.INT_2),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(55, 1, 56)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 30),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashMapWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(1).ToValue(),
							Value: value.Ref(value.String("foo")),
						},
					)),
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Float(5.6).ToValue(),
				},
			),
		},
		"with static elements and static capacity": {
			input: `{ 1 => 'foo', foo: 5, "bar" => 5.6 }:10`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(38, 1, 39)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashMapWithElements(
						nil,
						value.Pair{
							Key:   value.SmallInt(1).ToValue(),
							Value: value.Ref(value.String("foo")),
						},
						value.Pair{
							Key:   value.ToSymbol("foo").ToValue(),
							Value: value.SmallInt(5).ToValue(),
						},
						value.Pair{
							Key:   value.Ref(value.String("bar")),
							Value: value.Float(5.6).ToValue(),
						},
					)),
				},
			),
		},
		"with static elements and dynamic capacity": {
			input: `
				cap := 2
				{ 1 => 'foo', foo: 5, "bar" => 5.6 }:cap
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_2),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(58, 3, 45)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashMapWithElements(
						nil,
						value.Pair{
							Key:   value.SmallInt(1).ToValue(),
							Value: value.Ref(value.String("foo")),
						},
						value.Pair{
							Key:   value.ToSymbol("foo").ToValue(),
							Value: value.SmallInt(5).ToValue(),
						},
						value.Pair{
							Key:   value.Ref(value.String("bar")),
							Value: value.Float(5.6).ToValue(),
						},
					)),
				},
			),
		},
		"nested static": {
			input: "{ 1 => { 'bar' => [7.2] } }",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.INT_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.COPY),
					byte(bytecode.NEW_HASH_MAP8), 1,
					byte(bytecode.NEW_HASH_MAP8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(26, 1, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 14),
				},
				[]value.Value{
					value.Ref(value.NewHashMap(1)),
					value.Ref(value.NewHashMap(1)),
					value.Ref(value.String("bar")),
					value.Ref(&value.ArrayList{
						value.Float(7.2).ToValue(),
					}),
				},
			),
		},
		"with static and dynamic elements": {
			input: `
				a := 5
				{ 1 => 'foo', 5 => a, 5 => %[:foo] }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.INT_5),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.NEW_HASH_MAP8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(52, 3, 41)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashMapWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(1).ToValue(),
							Value: value.Ref(value.String("foo")),
						},
					)),
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with static elements and if modifiers": {
			input: `
				var a: Int? = nil
				{ 2 => 5, 1 => 5 if a, a: [:foo] }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.INT_1),
					byte(bytecode.INT_5),
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 3, 39)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 19),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashMapWithCapacityAndElements(
						nil,
						2,
						value.Pair{
							Key:   value.SmallInt(2).ToValue(),
							Value: value.SmallInt(5).ToValue(),
						},
					)),
					value.ToSymbol("a").ToValue(),
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with static elements, if modifiers and capacity": {
			input: `
				var a: Int? = nil
				{ 1 => 5 if a, 6 => [:foo] }:45
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(56, 3, 34), P(57, 3, 35)),
					"capacity cannot be specified in collection literals with conditional elements or loops",
				),
			},
		},
		"with static elements and unless modifiers": {
			input: `
				var a: Int? = nil
				{ 1 => 5 unless a, 9 => [:foo] }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF), 0, 6,
					byte(bytecode.INT_1),
					byte(bytecode.INT_5),
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.LOAD_INT_8), 9,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(59, 3, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 20),
				},
				[]value.Value{
					value.Ref(value.NewHashMap(2)),
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with dynamic elements and if modifiers": {
			input: `
				var a: Int? = nil
				{ Object() => 5 if a, 0 => [:foo] }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 9,
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.INT_5),
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.INT_0),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(62, 3, 40)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 22),
				},
				[]value.Value{
					value.Ref(value.NewHashMap(2)),
					value.ToSymbol("Std::Object").ToValue(),
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
					}),
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

func TestHashRecord(t *testing.T) {
	tests := testTable{
		"empty": {
			input: "%{}",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.NewHashRecord(0)),
				},
			),
		},
		"shorthand local": {
			input: `
				foo := 3
				%{ foo }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_3),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_HASH_RECORD8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(26, 3, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 6),
				},
				[]value.Value{
					value.ToSymbol("foo").ToValue(),
				},
			),
		},
		"shorthand private local": {
			input: `
				_foo := 3
				%{ _foo }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_3),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_HASH_RECORD8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 6),
				},
				[]value.Value{
					value.ToSymbol("_foo").ToValue(),
				},
			),
		},
		"with static elements": {
			input: `%{ 1 => 'foo', foo: 5, "bar" => 5.6 }`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(36, 1, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashRecordWithElements(
						nil,
						value.Pair{
							Key:   value.SmallInt(1).ToValue(),
							Value: value.Ref(value.String("foo")),
						},
						value.Pair{
							Key:   value.ToSymbol("foo").ToValue(),
							Value: value.SmallInt(5).ToValue(),
						},
						value.Pair{
							Key:   value.Ref(value.String("bar")),
							Value: value.Float(5.6).ToValue(),
						},
					)),
				},
			),
		},
		"with static elements and for loops": {
			input: `%{ 1 => 'foo', i => i ** 2 for i in [1, 2, 3], 2 => 5.6 }`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_RECORD8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 9,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.EXPONENTIATE_INT),
					byte(bytecode.MAP_SET),
					byte(bytecode.LOOP), 0, 13,
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.INT_2),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(56, 1, 57)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 29),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(1).ToValue(),
							Value: value.Ref(value.String("foo")),
						},
					)),
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Float(5.6).ToValue(),
				},
			),
		},
		"nested static": {
			input: "%{ 'foo' => 9, 1 => %{ 'bar' => [7.2] } }",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.INT_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.COPY),
					byte(bytecode.NEW_HASH_RECORD8), 1,
					byte(bytecode.NEW_HASH_RECORD8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 1, 41)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 12),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
						nil,
						2,
						value.Pair{Key: value.Ref(value.String("foo")), Value: value.SmallInt(9).ToValue()},
					)),
					value.Ref(value.NewHashRecord(1)),
					value.Ref(value.String("bar")),
					value.Ref(&value.ArrayList{
						value.Float(7.2).ToValue(),
					}),
				},
			),
		},
		"with static and dynamic elements": {
			input: `
				var a: Int? = nil
				%{ 1 => 'foo', 5 => a, 5 => %[:foo] }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.INT_5),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.NEW_HASH_RECORD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 3, 42)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 8),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(1).ToValue(),
							Value: value.Ref(value.String("foo")),
						},
					)),
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with static elements and if modifiers": {
			input: `
				var a: Int? = nil
				%{ 2 => 5, 1 => 5 if a, a: [:foo] }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.NEW_HASH_RECORD8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.INT_1),
					byte(bytecode.INT_5),
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(62, 3, 40)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 18),
				},
				[]value.Value{
					value.Ref(vm.MustNewHashRecordWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(2).ToValue(),
							Value: value.SmallInt(5).ToValue(),
						},
					)),
					value.ToSymbol("a").ToValue(),
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with static elements and unless modifiers": {
			input: `
				var a: Int? = nil
				%{ 1 => 5 unless a, 9 => [:foo] }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_HASH_RECORD8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF), 0, 6,
					byte(bytecode.INT_1),
					byte(bytecode.INT_5),
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.LOAD_INT_8), 9,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(60, 3, 38)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 19),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
					}),
				},
			),
		},
		"with dynamic elements and if modifiers": {
			input: `
				var a: Int? = nil
				%{ Object() => 5 if a, 0 => [:foo] }
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_HASH_RECORD8), 0,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 9,
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.INT_5),
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 0,
					byte(bytecode.INT_0),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 3, 41)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 21),
				},
				[]value.Value{
					value.ToSymbol("Std::Object").ToValue(),
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
					}),
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

func TestRegex(t *testing.T) {
	tests := testTable{
		"empty": {
			input: "%//",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.MustCompileRegex("", bitfield.BitField8FromBitFlag(0))),
				},
			),
		},
		"empty with flags": {
			input: "%//imx",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.MustCompileRegex("", bitfield.BitField8FromBitFlag(flag.CaseInsensitiveFlag|flag.MultilineFlag|flag.ExtendedFlag))),
				},
			),
		},
		"with content": {
			input: `%/foo \w+ bar/i`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Ref(value.MustCompileRegex(`foo \w+ bar`, bitfield.BitField8FromBitFlag(flag.CaseInsensitiveFlag))),
				},
			),
		},
		"with interpolation": {
			input: `
				a := "baz"
				%/foo \w+ ${a} bar/i
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_2),
					byte(bytecode.NEW_REGEX8), byte(flag.CaseInsensitiveFlag), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 25)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(value.String("baz")),
					value.Ref(value.String("foo \\w+ ")),
					value.Ref(value.String(" bar")),
				},
			),
		},
		"with compile error": {
			input: `%/foo\y/i`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(5, 1, 6), P(6, 1, 7)),
					`invalid escape sequence: \y`,
				),
			},
		},
		"with compile error from Go": {
			input: ` %/foo{1000000}/i`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(1, 1, 2), P(16, 1, 17)),
					"error parsing regexp: invalid repeat count: `{1000000}`",
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
