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
					byte(bytecode.APPEND_COLLECTION),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.APPEND_COLLECTION),
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
					byte(bytecode.APPEND_COLLECTION),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.APPEND_COLLECTION),
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}
