package compiler

import (
	"math"
	"math/big"
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestLiterals(t *testing.T) {
	tests := testTable{
		"put UInt8": {
			input: "1u8",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.UInt8(1),
				},
			),
		},
		"put UInt16": {
			input: "25u16",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.UInt16(25),
				},
			),
		},
		"put UInt32": {
			input: "450_200u32",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.UInt32(450200),
				},
			),
		},
		"put UInt64": {
			input: "450_200u64",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.UInt64(450200),
				},
			),
		},
		"put Int8": {
			input: "1i8",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Int8(1),
				},
			),
		},
		"put Int16": {
			input: "25i16",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Int16(25),
				},
			),
		},
		"put Int32": {
			input: "450_200i32",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Int32(450200),
				},
			),
		},
		"put Int64": {
			input: "450_200i64",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Int64(450200),
				},
			),
		},
		"put SmallInt": {
			input: "450_200",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.SmallInt(450200),
				},
			),
		},
		"put BigInt": {
			input: (&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5)).String(),
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
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
					value.ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5))),
				},
			),
		},
		"put Float64": {
			input: "45.5f64",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Float64(45.5),
				},
			),
		},
		"put Float32": {
			input: "45.5f32",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Float32(45.5),
				},
			),
		},
		"put Float": {
			input: "45.5",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Float(45.5),
				},
			),
		},
		"put Raw String": {
			input: `'foo\n'`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.String(`foo\n`),
				},
			),
		},
		"put String": {
			input: `"foo\n"`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.String("foo\n"),
				},
			),
		},
		"put raw Char": {
			input: `c'I'`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Char('I'),
				},
			),
		},
		"put Char": {
			input: `c"\n"`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Char('\n'),
				},
			),
		},
		"put nil": {
			input: `nil`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"put true": {
			input: `true`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"put false": {
			input: `false`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"put simple Symbol": {
			input: `:foo`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.ToSymbol("foo"),
				},
			),
		},
		"put self": {
			input: `self`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(3, 1, 4)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestTuples(t *testing.T) {
	tests := testTable{
		"empty tuple": {
			input: "%[]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					&value.Tuple{},
				},
			),
		},
		"with static elements": {
			input: "%[1, 'foo', 5, 5.6]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					&value.Tuple{
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					&value.Tuple{
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
		"nested static tuples": {
			input: "%[1, %['bar', %[7.2]]]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 1, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					&value.Tuple{
						value.SmallInt(1),
						&value.Tuple{
							value.String("bar"),
							&value.Tuple{
								value.Float(7.2),
							},
						},
					},
				},
			),
		},
		"with static keyed and dynamic elements": {
			input: "%[1, 'foo', 5 => 5,  3 => 5.6, foo()]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_FUNCTION8), 1,
					byte(bytecode.NEW_TUPLE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(36, 1, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					&value.Tuple{
						value.SmallInt(1),
						value.String("foo"),
						value.Nil,
						value.Float(5.6),
						value.Nil,
						value.SmallInt(5),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
				},
			),
		},
		"with static and dynamic elements": {
			input: "%[1, 'foo', 5, foo(), 5, %[:foo]]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_FUNCTION8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.NEW_TUPLE8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 1, 33)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					&value.Tuple{
						value.SmallInt(1),
						value.String("foo"),
						value.SmallInt(5),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(5),
					&value.Tuple{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with dynamic elements": {
			input: "%[foo(), 5, %[:foo]]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CALL_FUNCTION8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_TUPLE8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(5),
					&value.Tuple{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with static elements and if modifiers": {
			input: `
				%[1, 5 if foo(), %[:foo]]
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.CALL_FUNCTION8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(30, 2, 30)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 12),
				},
				[]value.Value{
					&value.Tuple{
						value.SmallInt(1),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(5),
					&value.Tuple{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with static elements and unless modifiers": {
			input: `
				%[1, 5 unless foo(), %[:foo]]
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.CALL_FUNCTION8), 1,
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
				L(P(0, 1, 1), P(34, 2, 34)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 12),
				},
				[]value.Value{
					&value.Tuple{
						value.SmallInt(1),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(5),
					&value.Tuple{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with dynamic elements and if modifiers": {
			input: `
				%[self.bar, 5 if foo(), %[:foo]]
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.NEW_TUPLE8), 1,
					byte(bytecode.CALL_FUNCTION8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.APPEND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(37, 2, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 14),
				},
				[]value.Value{
					value.NewCallSiteInfo(
						value.ToSymbol("bar"),
						0,
						nil,
					),
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(5),
					&value.Tuple{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with dynamic and keyed elements": {
			input: "%[foo(), 1, 'foo', 5 => 5,  3 => 5.6]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CALL_FUNCTION8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_TUPLE8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.APPEND_AT),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.APPEND_AT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(36, 1, 37)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 12),
				},
				[]value.Value{
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(1),
					value.String("foo"),
					value.SmallInt(5),
					value.SmallInt(3),
					value.Float(5.6),
				},
			),
		},
		"with keyed and if elements": {
			input: "%[3 => 5 if foo()]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_TUPLE8), 0,
					byte(bytecode.CALL_FUNCTION8), 0,
					byte(bytecode.JUMP_UNLESS), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND_AT),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(3),
					value.SmallInt(5),
				},
			),
		},
		"with static concat": {
			input: "%[1, 2, 3] + %[4, 5, 6] + %[10]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(30, 1, 31)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					&value.Tuple{
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 1, 30)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					&value.List{
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					&value.Tuple{
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
			input: "%[1, 2, 3] + %[4, 5, 6, %[7, 8]] + %[10]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(39, 1, 40)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					&value.Tuple{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
						value.SmallInt(6),
						&value.Tuple{
							value.SmallInt(7),
							value.SmallInt(8),
						},
						value.SmallInt(10),
					},
				},
			),
		},
		"word tuple": {
			input: `%w[foo bar baz]`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					&value.Tuple{
						value.String("foo"),
						value.String("bar"),
						value.String("baz"),
					},
				},
			),
		},
		"symbol tuple": {
			input: `%s[foo bar baz]`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					&value.Tuple{
						value.ToSymbol("foo"),
						value.ToSymbol("bar"),
						value.ToSymbol("baz"),
					},
				},
			),
		},
		"hex tuple": {
			input: `%x[ab cd 5f]`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					&value.Tuple{
						value.SmallInt(0xab),
						value.SmallInt(0xcd),
						value.SmallInt(0x5f),
					},
				},
			),
		},
		"bin tuple": {
			input: `%b[101 11 10]`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					&value.Tuple{
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

func TestLists(t *testing.T) {
	tests := testTable{
		"empty list": {
			input: "[]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(1, 1, 2)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					&value.List{},
				},
			),
		},
		"with static elements": {
			input: "[1, 'foo', 5, 5.6]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					&value.List{
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					&value.List{
						value.String("foo"),
						value.String("bar"),
						value.String("baz"),
					},
				},
			),
		},
		"symbol list": {
			input: `\s[foo bar baz]`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					&value.List{
						value.ToSymbol("foo"),
						value.ToSymbol("bar"),
						value.ToSymbol("baz"),
					},
				},
			),
		},
		"hex list": {
			input: `\x[ab cd 5f]`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					&value.List{
						value.SmallInt(0xab),
						value.SmallInt(0xcd),
						value.SmallInt(0x5f),
					},
				},
			),
		},
		"bin list": {
			input: `\b[101 11 10]`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					&value.List{
						value.SmallInt(0b101),
						value.SmallInt(0b11),
						value.SmallInt(0b10),
					},
				},
			),
		},
		"with static keyed elements": {
			input: "[1, 'foo', 5 => 5,  3 => 5.6, :lol]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 1, 35)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					&value.List{
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(27, 1, 28)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					&value.List{
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					&value.List{
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_LIST8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					&value.List{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
						value.SmallInt(6),
					},
					&value.List{
						value.SmallInt(7),
						value.SmallInt(8),
					},
					value.SmallInt(10),
				},
			),
		},
		"nested static lists": {
			input: "[1, ['bar', [7.2]]]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.NEW_LIST8), 1,
					byte(bytecode.NEW_LIST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					&value.List{
						value.SmallInt(1),
					},
					&value.List{
						value.String("bar"),
					},
					&value.List{
						value.Float(7.2),
					},
				},
			),
		},
		"with static keyed and dynamic elements": {
			input: "[1, 'foo', 5 => 5,  3 => 5.6, foo()]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_FUNCTION8), 1,
					byte(bytecode.NEW_LIST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					&value.List{
						value.SmallInt(1),
						value.String("foo"),
						value.Nil,
						value.Float(5.6),
						value.Nil,
						value.SmallInt(5),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
				},
			),
		},
		"with static and dynamic elements": {
			input: "[1, 'foo', 5, foo(), 5, %[:foo]]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_FUNCTION8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.NEW_LIST8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(31, 1, 32)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					&value.List{
						value.SmallInt(1),
						value.String("foo"),
						value.SmallInt(5),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(5),
					&value.Tuple{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with dynamic elements": {
			input: "[foo(), 5, [:foo]]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CALL_FUNCTION8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.NEW_LIST8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(5),
					&value.List{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with static elements and if modifiers": {
			input: `
				[1, 5 if foo(), [:foo]]
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.CALL_FUNCTION8), 1,
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
				L(P(0, 1, 1), P(28, 2, 28)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 13),
				},
				[]value.Value{
					&value.List{
						value.SmallInt(1),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(5),
					&value.List{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with static elements and unless modifiers": {
			input: `
				[1, 5 unless foo(), [:foo]]
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.CALL_FUNCTION8), 1,
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
				L(P(0, 1, 1), P(32, 2, 32)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 13),
				},
				[]value.Value{
					&value.List{
						value.SmallInt(1),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(5),
					&value.List{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with dynamic elements and if modifiers": {
			input: `
				[self.bar, 5 if foo(), [:foo]]
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.NEW_LIST8), 1,
					byte(bytecode.CALL_FUNCTION8), 1,
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
				L(P(0, 1, 1), P(35, 2, 35)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 15),
				},
				[]value.Value{
					value.NewCallSiteInfo(
						value.ToSymbol("bar"),
						0,
						nil,
					),
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(5),
					&value.List{
						value.ToSymbol("foo"),
					},
				},
			),
		},
		"with dynamic and keyed elements": {
			input: "[foo(), 1, 'foo', 5 => 5,  3 => 5.6]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CALL_FUNCTION8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NEW_LIST8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.APPEND_AT),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.APPEND_AT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 1, 36)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 12),
				},
				[]value.Value{
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
					value.SmallInt(1),
					value.String("foo"),
					value.SmallInt(5),
					value.SmallInt(3),
					value.Float(5.6),
				},
			),
		},
		"with keyed and if elements": {
			input: "[3 => 5 if foo()]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_LIST8), 0,
					byte(bytecode.CALL_FUNCTION8), 0,
					byte(bytecode.JUMP_UNLESS), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.APPEND_AT),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{
					value.NewCallSiteInfo(
						value.ToSymbol("foo"),
						0,
						nil,
					),
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
