package compiler_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestStringLiteral(t *testing.T) {
	tests := testTable{
		"static string": {
			input: `"foo bar"`,
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
					nil,
					value.String("foo bar"),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_STRING8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(62, 4, 33)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 14),
				},
				[]value.Value{
					nil,
					value.Float(15.2),
					value.SmallInt(1),
					value.String("foo: "),
					value.SmallInt(2),
					value.String(", bar: "),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.ADD),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CALL_METHOD8), 7,
					byte(bytecode.NEW_STRING8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(62, 4, 33)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 18),
				},
				[]value.Value{
					nil,
					value.Float(15.2),
					value.SmallInt(1),
					value.String("foo: "),
					value.SmallInt(2),
					value.NewCallSiteInfo(value.ToSymbol("inspect"), 0),
					value.String(", bar: "),
					value.NewCallSiteInfo(value.ToSymbol("inspect"), 0),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.NewClosedRange(value.SmallInt(2), value.SmallInt(5)),
				},
			),
		},
		"static open range": {
			input: `2<.<5`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.NewOpenRange(value.SmallInt(2), value.SmallInt(5)),
				},
			),
		},
		"static left open range": {
			input: `2<..5`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.NewLeftOpenRange(value.SmallInt(2), value.SmallInt(5)),
				},
			),
		},
		"static right open range": {
			input: `2..<5`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.NewRightOpenRange(value.SmallInt(2), value.SmallInt(5)),
				},
			),
		},
		"static beginless closed range": {
			input: `...5`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.NewBeginlessClosedRange(value.SmallInt(5)),
				},
			),
		},
		"static beginless open range": {
			input: `..<5`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.NewBeginlessOpenRange(value.SmallInt(5)),
				},
			),
		},
		"static endless closed range": {
			input: `2...`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.NewEndlessClosedRange(value.SmallInt(2)),
				},
			),
		},
		"static endless open range": {
			input: `2<..`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.NewEndlessOpenRange(value.SmallInt(2)),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_RANGE), bytecode.CLOSED_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_RANGE), bytecode.OPEN_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_RANGE), bytecode.LEFT_OPEN_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_RANGE), bytecode.RIGHT_OPEN_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_RANGE), bytecode.BEGINLESS_CLOSED_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 3, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
				},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_RANGE), bytecode.BEGINLESS_OPEN_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 3, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
				},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_RANGE), bytecode.ENDLESS_CLOSED_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 3, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
				},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_RANGE), bytecode.ENDLESS_OPEN_RANGE_FLAG,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 3, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
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

func TestLiterals(t *testing.T) {
	tests := testTable{
		"put UInt8": {
			input: "1u8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.UInt8(1),
				},
			),
		},
		"put UInt16": {
			input: "25u16",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.UInt16(25),
				},
			),
		},
		"put UInt32": {
			input: "450_200u32",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.UInt32(450200),
				},
			),
		},
		"put UInt64": {
			input: "450_200u64",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.UInt64(450200),
				},
			),
		},
		"put Int8": {
			input: "1i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.Int8(1),
				},
			),
		},
		"put Int16": {
			input: "25i16",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.Int16(25),
				},
			),
		},
		"put Int32": {
			input: "450_200i32",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.Int32(450200),
				},
			),
		},
		"put Int64": {
			input: "450_200i64",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.Int64(450200),
				},
			),
		},
		"put SmallInt": {
			input: "450_200",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.SmallInt(450200),
				},
			),
		},
		"put BigInt": {
			input: (&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5)).String(),
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
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
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5))),
				},
			),
		},
		"put Float64": {
			input: "45.5f64",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.Float64(45.5),
				},
			),
		},
		"put Float32": {
			input: "45.5f32",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.Float32(45.5),
				},
			),
		},
		"put Float": {
			input: "45.5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.Float(45.5),
				},
			),
		},
		"put Raw String": {
			input: `'foo\n'`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.String(`foo\n`),
				},
			),
		},
		"put String": {
			input: `"foo\n"`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.String("foo\n"),
				},
			),
		},
		"put raw Char": {
			input: "`I`",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.Char('I'),
				},
			),
		},
		"put Char": {
			input: "`\\n`",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.Char('\n'),
				},
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
				[]value.Value{nil},
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
				[]value.Value{nil},
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
				[]value.Value{nil},
			),
		},
		"put simple Symbol": {
			input: `:foo`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.ToSymbol("foo"),
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
				[]value.Value{nil},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{},
				},
			),
		},
		"with static elements": {
			input: "%[1, 'foo', 5, 5.6]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(1),
						value.String("foo"),
						value.SmallInt(5),
						value.Float(5.6),
					},
				},
			),
		},
		"with static keyed elements": {
			input: "%[1, 'foo', 5 => 5,  3 => 5.6, :lol]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(1),
						value.String("foo"),
						value.Nil,
						value.Float(5.6),
						value.Nil,
						value.SmallInt(5),
						value.ToSymbol("lol"),
					},
				},
			),
		},
		"nested static arrayTuples": {
			input: "%[1, %['bar', %[7.2]]]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 1, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(1),
						&value.ArrayTuple{
							value.String("bar"),
							&value.ArrayTuple{
								value.Float(7.2),
							},
						},
					},
				},
			),
		},
		"nested static with mutable elements": {
			input: "%[1, %['bar', [7.2]]]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 1, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 12),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(1),
					},
					&value.ArrayTuple{
						value.String("bar"),
					},
					&value.ArrayList{
						value.Float(7.2),
					},
				},
			),
		},
		"with static keyed and dynamic elements": {
			input: "%[1, 'foo', 5 => 5,  3 => 5.6, String.name]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 1, 43)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(1),
						value.String("foo"),
						value.Nil,
						value.Float(5.6),
						value.Nil,
						value.SmallInt(5),
					},
					value.ToSymbol("Std::String"),
					value.NewCallSiteInfo(
						value.ToSymbol("name"),
						0,
					),
				},
			),
		},
		"with static and dynamic elements": {
			input: "%[1, 'foo', 5, Object(), 5, %[:foo]]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.NEW_ARRAY_TUPLE8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(1),
						value.String("foo"),
						value.SmallInt(5),
					},
					value.ToSymbol("Std::Object"),
					value.SmallInt(5),
					&value.ArrayTuple{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with dynamic elements": {
			input: "%[Object(), 5, %[:foo]]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.NEW_ARRAY_TUPLE8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 1, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 12),
				},
				[]value.Value{
					nil,
					value.ToSymbol("Std::Object"),
					value.SmallInt(5),
					&value.ArrayTuple{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(56, 3, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 20),
				},
				[]value.Value{
					nil,
					value.ToSymbol("Std::Object"),
					&value.ArrayTuple{
						value.SmallInt(1),
					},
					value.SmallInt(5),
					&value.ArrayTuple{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_IF), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(55, 3, 30)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 20),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(1),
					},
					value.SmallInt(5),
					&value.ArrayTuple{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.FOR_IN), 0, 15,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.MULTIPLY),
					byte(bytecode.APPEND),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(44, 2, 44)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 37),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(1),
					},
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
					},
					value.SmallInt(2),
					&value.ArrayTuple{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 3, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 24),
				},
				[]value.Value{
					nil,
					value.ToSymbol("Std::String"),
					value.NewCallSiteInfo(
						value.ToSymbol("name"),
						0,
					),
					value.SmallInt(5),
					&value.ArrayTuple{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with dynamic and keyed elements": {
			input: "%[Object(), 1, 'foo', 5 => 5,  3 => 5.6]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.NEW_ARRAY_TUPLE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.APPEND_AT),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.APPEND_AT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(39, 1, 40)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 22),
				},
				[]value.Value{
					nil,
					value.ToSymbol("Std::Object"),
					value.SmallInt(1),
					value.String("foo"),
					value.SmallInt(5),
					value.SmallInt(3),
					value.Float(5.6),
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_ARRAY_TUPLE8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND_AT),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(44, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 19),
				},
				[]value.Value{
					nil,
					value.SmallInt(3),
					value.SmallInt(5),
				},
			),
		},

		"with static concat": {
			input: "%[1, 2, 3] + %[4, 5, 6] + %[10]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(30, 1, 31)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
						value.SmallInt(6),
						value.SmallInt(10),
					},
				},
			),
		},
		"with static concat with list": {
			input: "%[1, 2, 3] + [4, 5, 6] + %[10]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 1, 30)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
						value.SmallInt(6),
						value.SmallInt(10),
					},
				},
			),
		},
		"with static repeat": {
			input: "%[1, 2, 3] * 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
					},
				},
			),
		},
		"with static concat and nested tuples": {
			input: "%[1, 2, 3] + %[4, 5, 6, %[7, 8]] + %[10]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(39, 1, 40)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
						value.SmallInt(6),
						&value.ArrayTuple{
							value.SmallInt(7),
							value.SmallInt(8),
						},
						value.SmallInt(10),
					},
				},
			),
		},
		"word arrayTuple": {
			input: `%w[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.String("foo"),
						value.String("bar"),
						value.String("baz"),
					},
				},
			),
		},
		"symbol arrayTuple": {
			input: `%s[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.ToSymbol("foo"),
						value.ToSymbol("bar"),
						value.ToSymbol("baz"),
					},
				},
			),
		},
		"hex arrayTuple": {
			input: `%x[ab cd 5f]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(0xab),
						value.SmallInt(0xcd),
						value.SmallInt(0x5f),
					},
				},
			),
		},
		"bin arrayTuple": {
			input: `%b[101 11 10]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					&value.ArrayTuple{
						value.SmallInt(0b101),
						value.SmallInt(0b11),
						value.SmallInt(0b10),
					},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(1, 1, 2)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					&value.ArrayList{},
				},
			),
		},
		"with static elements": {
			input: "[1, 'foo', 5, 5.6]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(1),
						value.String("foo"),
						value.SmallInt(5),
						value.Float(5.6),
					},
				},
			),
		},
		"with static elements and static capacity": {
			input: "[1, 'foo', 5, 5.6]:10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 1, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(10),
					&value.ArrayList{
						value.SmallInt(1),
						value.String("foo"),
						value.SmallInt(5),
						value.Float(5.6),
					},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
					&value.ArrayList{
						value.SmallInt(1),
						value.String("foo"),
						value.SmallInt(5),
						value.Float(5.6),
					},
				},
			),
		},
		"word list": {
			input: `\w[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.String("foo"),
						value.String("bar"),
						value.String("baz"),
					},
				},
			),
		},
		"word list with capacity": {
			input: `\w[foo bar baz]:15`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(15),
					&value.ArrayList{
						value.String("foo"),
						value.String("bar"),
						value.String("baz"),
					},
				},
			),
		},
		"symbol list": {
			input: `\s[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.ToSymbol("foo"),
						value.ToSymbol("bar"),
						value.ToSymbol("baz"),
					},
				},
			),
		},

		"symbol list with capacity": {
			input: `\s[foo bar baz]:15`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(15),
					&value.ArrayList{
						value.ToSymbol("foo"),
						value.ToSymbol("bar"),
						value.ToSymbol("baz"),
					},
				},
			),
		},
		"hex list": {
			input: `\x[ab cd 5f]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(0xab),
						value.SmallInt(0xcd),
						value.SmallInt(0x5f),
					},
				},
			),
		},
		"hex list with capacity": {
			input: `\x[ab cd 5f]:2`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
					&value.ArrayList{
						value.SmallInt(0xab),
						value.SmallInt(0xcd),
						value.SmallInt(0x5f),
					},
				},
			),
		},
		"bin list": {
			input: `\b[101 11 10]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(0b101),
						value.SmallInt(0b11),
						value.SmallInt(0b10),
					},
				},
			),
		},
		"bin list with capacity": {
			input: `\b[101 11 10]:3`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(3),
					&value.ArrayList{
						value.SmallInt(0b101),
						value.SmallInt(0b11),
						value.SmallInt(0b10),
					},
				},
			),
		},

		"with static keyed elements": {
			input: "[1, 'foo', 5 => 5,  3 => 5.6, :lol]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 1, 35)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(1),
						value.String("foo"),
						value.Nil,
						value.Float(5.6),
						value.Nil,
						value.SmallInt(5),
						value.ToSymbol("lol"),
					},
				},
			),
		},
		"with static keyed elements and static capacity": {
			input: "[1, 'foo', 5 => 5,  3 => 5.6, :lol]:6",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(36, 1, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(6),
					&value.ArrayList{
						value.SmallInt(1),
						value.String("foo"),
						value.Nil,
						value.Float(5.6),
						value.Nil,
						value.SmallInt(5),
						value.ToSymbol("lol"),
					},
				},
			),
		},
		"with static concat": {
			input: "[1, 2, 3] + [4, 5, 6] + [10]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(27, 1, 28)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
						value.SmallInt(6),
						value.SmallInt(10),
					},
				},
			),
		},
		"with static repeat": {
			input: "[1, 2, 3] * 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
					},
				},
			),
		},
		"with static concat and nested lists": {
			input: "[1, 2, 3] + [4, 5, 6, [7, 8]] + [10]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.NEW_ARRAY_LIST8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
						value.SmallInt(6),
					},
					&value.ArrayList{
						value.SmallInt(7),
						value.SmallInt(8),
					},
					value.SmallInt(10),
				},
			),
		},

		"nested static lists": {
			input: "[1, ['bar', [7.2]]]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 14),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(1),
					},
					&value.ArrayList{
						value.String("bar"),
					},
					&value.ArrayList{
						value.Float(7.2),
					},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(48, 3, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 8),
				},
				[]value.Value{
					nil,
					value.SmallInt(5),
					&value.ArrayList{
						value.SmallInt(1),
						value.String("foo"),
						value.Nil,
						value.Float(5.6),
						value.Nil,
						value.SmallInt(5),
					},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(51, 3, 40)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 9),
				},
				[]value.Value{
					nil,
					value.SmallInt(5),
					value.SmallInt(15),
					&value.ArrayList{
						value.SmallInt(1),
						value.String("foo"),
						value.Nil,
						value.Float(5.6),
						value.Nil,
						value.SmallInt(5),
					},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.NEW_ARRAY_LIST8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(53, 3, 33)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 12),
				},
				[]value.Value{
					nil,
					value.SmallInt(3),
					&value.ArrayList{
						value.SmallInt(1),
						value.String("foo"),
						value.SmallInt(5),
					},
					value.SmallInt(5),
					&value.ArrayTuple{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.NEW_ARRAY_LIST8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(30, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 12),
				},
				[]value.Value{
					nil,
					value.SmallInt(3),
					value.SmallInt(5),
					&value.ArrayList{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 3, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 23),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(1),
					},
					value.SmallInt(5),
					&value.ArrayList{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with static elements, if modifiers and capacity": {
			input: `
				var a: Int? = nil
				[1, 5 if a, [:foo]]:45
			`,
			err: error.ErrorList{
				error.NewFailure(
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_IF), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(50, 3, 28)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 23),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(1),
					},
					value.SmallInt(5),
					&value.ArrayList{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.FOR_IN), 0, 15,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.MULTIPLY),
					byte(bytecode.APPEND),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 2, 43)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 39),
				},
				[]value.Value{
					nil,
					&value.ArrayList{
						value.SmallInt(1),
					},
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
					},
					value.SmallInt(2),
					&value.ArrayTuple{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(53, 3, 31)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 26),
				},
				[]value.Value{
					nil,
					value.ToSymbol("Std::Object"),
					value.SmallInt(5),
					&value.ArrayList{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.NEW_ARRAY_LIST8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.APPEND_AT),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.APPEND_AT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(59, 3, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 22),
				},
				[]value.Value{
					nil,
					&value.ArrayList{},
					value.SmallInt(1),
					value.String("foo"),
					value.SmallInt(5),
					value.SmallInt(3),
					value.Float(5.6),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.NEW_ARRAY_LIST8), 3,
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.APPEND_AT),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.APPEND_AT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(50, 3, 39)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 23),
				},
				[]value.Value{
					nil,
					value.SmallInt(3),
					value.SmallInt(7),
					&value.ArrayList{},
					value.SmallInt(1),
					value.String("foo"),
					value.SmallInt(5),
					value.Float(5.6),
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND_AT),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 20),
				},
				[]value.Value{
					nil,
					value.SmallInt(3),
					value.SmallInt(5),
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

func TestHashSet(t *testing.T) {
	tests := testTable{
		"empty list": {
			input: "^[]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					&value.HashSet{},
				},
			),
		},
		"with static elements": {
			input: "^[1, 'foo', 5, 5.6]",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(1),
						value.String("foo"),
						value.SmallInt(5),
						value.Float(5.6),
					),
				},
			),
		},
		"with static elements and static capacity": {
			input: "^[1, 'foo', 5, 5.6]:10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 1, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(10),
					vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
						nil,
						4,
						1,
						value.SmallInt(1),
						value.String("foo"),
						value.SmallInt(5),
						value.Float(5.6),
					),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 3, 28)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
					vm.MustNewHashSetWithCapacityAndElementsMaxLoad(
						nil,
						4,
						1,
						value.SmallInt(1),
						value.String("foo"),
						value.SmallInt(5),
						value.Float(5.6),
					),
				},
			),
		},

		"word set": {
			input: `^w[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					vm.MustNewHashSetWithElements(
						nil,
						value.String("foo"),
						value.String("bar"),
						value.String("baz"),
					),
				},
			),
		},
		"word set with capacity": {
			input: `^w[foo bar baz]:15`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(15),
					vm.MustNewHashSetWithElements(
						nil,
						value.String("foo"),
						value.String("bar"),
						value.String("baz"),
					),
				},
			),
		},
		"symbol set": {
			input: `^s[foo bar baz]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					vm.MustNewHashSetWithElements(
						nil,
						value.ToSymbol("foo"),
						value.ToSymbol("bar"),
						value.ToSymbol("baz"),
					),
				},
			),
		},
		"symbol set with capacity": {
			input: `^s[foo bar baz]:15`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(15),
					vm.MustNewHashSetWithElements(
						nil,
						value.ToSymbol("foo"),
						value.ToSymbol("bar"),
						value.ToSymbol("baz"),
					),
				},
			),
		},
		"hex set": {
			input: `^x[ab cd 5f]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(0xab),
						value.SmallInt(0xcd),
						value.SmallInt(0x5f),
					),
				},
			),
		},
		"hex set with capacity": {
			input: `^x[ab cd 5f]:2`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
					vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(0xab),
						value.SmallInt(0xcd),
						value.SmallInt(0x5f),
					),
				},
			),
		},

		"bin set": {
			input: `^b[101 11 10]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(0b101),
						value.SmallInt(0b11),
						value.SmallInt(0b10),
					),
				},
			),
		},
		"bin set with capacity": {
			input: `^b[101 11 10]:3`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(3),
					vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(0b101),
						value.SmallInt(0b11),
						value.SmallInt(0b10),
					),
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_HASH_SET8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 3, 25)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					nil,
					vm.MustNewHashSetWithCapacityAndElements(
						nil,
						5,
						value.String("foo"),
						value.SmallInt(5),
						value.SmallInt(1),
					),
					value.SmallInt(5),
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_HASH_SET8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					nil,
					vm.MustNewHashSetWithCapacityAndElements(
						nil,
						2,
					),
					value.SmallInt(5),
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(39, 3, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 19),
				},
				[]value.Value{
					nil,
					vm.MustNewHashSetWithCapacityAndElements(
						nil,
						2,
						value.SmallInt(1),
					),
					value.SmallInt(5),
				},
			),
		},
		"with static elements, if modifiers and capacity": {
			input: `
				var a: Int? = nil
				^[1, 5 if a]:45
			`,
			err: error.ErrorList{
				error.NewFailure(
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_IF), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 3, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 19),
				},
				[]value.Value{
					nil,
					vm.MustNewHashSetWithCapacityAndElements(
						nil,
						2,
						value.SmallInt(1),
					),
					value.SmallInt(5),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_HASH_SET8), 0,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.FOR_IN), 0, 15,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.MULTIPLY),
					byte(bytecode.APPEND),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(38, 2, 38)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 39),
				},
				[]value.Value{
					nil,
					vm.MustNewHashSetWithCapacityAndElements(
						nil,
						3,
						value.SmallInt(1),
					),
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
					},
					value.SmallInt(2),
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.NEW_HASH_SET8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 3, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 23),
				},
				[]value.Value{
					nil,
					vm.MustNewHashSetWithCapacityAndElements(
						nil,
						2,
					),
					value.ToSymbol("Std::Object"),
					value.SmallInt(5),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(1, 1, 2)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					value.NewHashMap(0),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_HASH_MAP8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 3, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					nil,
					value.SmallInt(3),
					value.NewHashMap(0),
					value.ToSymbol("foo"),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_HASH_MAP8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(27, 3, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					nil,
					value.SmallInt(3),
					value.NewHashMap(0),
					value.ToSymbol("_foo"),
				},
			),
		},
		"with static elements": {
			input: `{ 1 => 'foo', foo: 5, "bar" => 5.6 }`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					nil,
					vm.MustNewHashMapWithElements(
						nil,
						value.Pair{
							Key:   value.SmallInt(1),
							Value: value.String("foo"),
						},
						value.Pair{
							Key:   value.ToSymbol("foo"),
							Value: value.SmallInt(5),
						},
						value.Pair{
							Key:   value.String("bar"),
							Value: value.Float(5.6),
						},
					),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.FOR_IN), 0, 17,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.MAP_SET),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 22,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(55, 1, 56)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 45),
				},
				[]value.Value{
					nil,
					vm.MustNewHashMapWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(1),
							Value: value.String("foo"),
						},
					),
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
					},
					value.SmallInt(2),
					value.Float(5.6),
				},
			),
		},
		"with static elements and static capacity": {
			input: `{ 1 => 'foo', foo: 5, "bar" => 5.6 }:10`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(38, 1, 39)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(10),
					vm.MustNewHashMapWithElements(
						nil,
						value.Pair{
							Key:   value.SmallInt(1),
							Value: value.String("foo"),
						},
						value.Pair{
							Key:   value.ToSymbol("foo"),
							Value: value.SmallInt(5),
						},
						value.Pair{
							Key:   value.String("bar"),
							Value: value.Float(5.6),
						},
					),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(58, 3, 45)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					nil,
					value.SmallInt(2),
					vm.MustNewHashMapWithElements(
						nil,
						value.Pair{
							Key:   value.SmallInt(1),
							Value: value.String("foo"),
						},
						value.Pair{
							Key:   value.ToSymbol("foo"),
							Value: value.SmallInt(5),
						},
						value.Pair{
							Key:   value.String("bar"),
							Value: value.Float(5.6),
						},
					),
				},
			),
		},
		"nested static": {
			input: "{ 1 => { 'bar' => [7.2] } }",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.COPY),
					byte(bytecode.NEW_HASH_MAP8), 1,
					byte(bytecode.NEW_HASH_MAP8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(26, 1, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 18),
				},
				[]value.Value{
					nil,
					value.NewHashMap(1),
					value.SmallInt(1),
					value.NewHashMap(1),
					value.String("bar"),
					&value.ArrayList{
						value.Float(7.2),
					},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.NEW_HASH_MAP8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(52, 3, 41)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 14),
				},
				[]value.Value{
					nil,
					value.SmallInt(5),
					vm.MustNewHashMapWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(1),
							Value: value.String("foo"),
						},
					),
					&value.ArrayTuple{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 3, 39)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 27),
				},
				[]value.Value{
					nil,
					vm.MustNewHashMapWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(2),
							Value: value.SmallInt(5),
						},
					),
					value.SmallInt(1),
					value.SmallInt(5),
					value.ToSymbol("a"),
					&value.ArrayList{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with static elements, if modifiers and capacity": {
			input: `
				var a: Int? = nil
				{ 1 => 5 if a, 6 => [:foo] }:45
			`,
			err: error.ErrorList{
				error.NewFailure(
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_IF), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(59, 3, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 27),
				},
				[]value.Value{
					nil,
					value.NewHashMap(2),
					value.SmallInt(1),
					value.SmallInt(5),
					value.SmallInt(9),
					&value.ArrayList{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_HASH_MAP8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(62, 3, 40)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 29),
				},
				[]value.Value{
					nil,
					value.NewHashMap(2),
					value.ToSymbol("Std::Object"),
					value.SmallInt(5),
					value.SmallInt(0),
					&value.ArrayList{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.NewHashRecord(0),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_HASH_RECORD8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(26, 3, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 8),
				},
				[]value.Value{
					nil,
					value.SmallInt(3),
					value.ToSymbol("foo"),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_HASH_RECORD8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 8),
				},
				[]value.Value{
					nil,
					value.SmallInt(3),
					value.ToSymbol("_foo"),
				},
			),
		},
		"with static elements": {
			input: `%{ 1 => 'foo', foo: 5, "bar" => 5.6 }`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(36, 1, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					vm.MustNewHashRecordWithElements(
						nil,
						value.Pair{
							Key:   value.SmallInt(1),
							Value: value.String("foo"),
						},
						value.Pair{
							Key:   value.ToSymbol("foo"),
							Value: value.SmallInt(5),
						},
						value.Pair{
							Key:   value.String("bar"),
							Value: value.Float(5.6),
						},
					),
				},
			),
		},
		"with static elements and for loops": {
			input: `%{ 1 => 'foo', i => i ** 2 for i in [1, 2, 3], 2 => 5.6 }`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_HASH_RECORD8), 0,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.FOR_IN), 0, 17,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.MAP_SET),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 22,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(56, 1, 57)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 44),
				},
				[]value.Value{
					nil,
					vm.MustNewHashRecordWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(1),
							Value: value.String("foo"),
						},
					),
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
					},
					value.SmallInt(2),
					value.Float(5.6),
				},
			),
		},
		"nested static": {
			input: "%{ 'foo' => 9, 1 => %{ 'bar' => [7.2] } }",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.COPY),
					byte(bytecode.NEW_HASH_RECORD8), 1,
					byte(bytecode.NEW_HASH_RECORD8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 1, 41)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 16),
				},
				[]value.Value{
					nil,
					vm.MustNewHashRecordWithCapacityAndElements(
						nil,
						2,
						value.Pair{Key: value.String("foo"), Value: value.SmallInt(9)},
					),
					value.SmallInt(1),
					value.NewHashRecord(1),
					value.String("bar"),
					&value.ArrayList{
						value.Float(7.2),
					},
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.NEW_HASH_RECORD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 3, 42)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 13),
				},
				[]value.Value{
					nil,
					vm.MustNewHashRecordWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(1),
							Value: value.String("foo"),
						},
					),
					value.SmallInt(5),
					&value.ArrayTuple{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NEW_HASH_RECORD8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(62, 3, 40)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 26),
				},
				[]value.Value{
					nil,
					vm.MustNewHashRecordWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(2),
							Value: value.SmallInt(5),
						},
					),
					value.SmallInt(1),
					value.SmallInt(5),
					value.ToSymbol("a"),
					&value.ArrayList{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_HASH_RECORD8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_IF), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(60, 3, 38)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 25),
				},
				[]value.Value{
					nil,
					value.SmallInt(1),
					value.SmallInt(5),
					value.SmallInt(9),
					&value.ArrayList{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_HASH_RECORD8), 0,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.INSTANTIATE8), 0,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.MAP_SET),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.COPY),
					byte(bytecode.MAP_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 3, 41)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 27),
				},
				[]value.Value{
					nil,
					value.ToSymbol("Std::Object"),
					value.SmallInt(5),
					value.SmallInt(0),
					&value.ArrayList{
						value.ToSymbol("foo"),
					},
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.MustCompileRegex("", bitfield.BitField8FromBitFlag(0)),
				},
			),
		},
		"empty with flags": {
			input: "%//imx",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.MustCompileRegex("", bitfield.BitField8FromBitFlag(flag.CaseInsensitiveFlag|flag.MultilineFlag|flag.ExtendedFlag)),
				},
			),
		},
		"with content": {
			input: `%/foo \w+ bar/i`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					nil,
					value.MustCompileRegex(`foo \w+ bar`, bitfield.BitField8FromBitFlag(flag.CaseInsensitiveFlag)),
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
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.NEW_REGEX8), byte(flag.CaseInsensitiveFlag), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 3, 25)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 10),
				},
				[]value.Value{
					nil,
					value.String("baz"),
					value.String("foo \\w+ "),
					value.String(" bar"),
				},
			),
		},
		"with compile error": {
			input: `%/foo\y/i`,
			err: error.ErrorList{
				error.NewFailure(
					L(P(5, 1, 6), P(6, 1, 7)),
					`invalid escape sequence: \y`,
				),
			},
		},
		"with compile error from Go": {
			input: ` %/foo{1000000}/i`,
			err: error.ErrorList{
				error.NewFailure(
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
