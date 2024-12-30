package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestBinaryExpressions(t *testing.T) {
	tests := testTable{
		"is a": {
			input: "3 <: ::Std::Int",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(3).ToValue(),
					value.ToSymbol("Std::Int").ToValue(),
				},
			),
			err: error.ErrorList{
				error.NewWarning(L(P(0, 1, 1), P(0, 1, 1)), "this \"is a\" check is always true, `3` will always be an instance of `Std::Int`"),
			},
		},
		"instance of": {
			input: "3 <<: ::Std::Int",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.INSTANCE_OF),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(3).ToValue(),
					value.ToSymbol("Std::Int").ToValue(),
				},
			),
			err: error.ErrorList{
				error.NewWarning(L(P(0, 1, 1), P(0, 1, 1)), "this \"instance of\" check is always true, `3` will always be an instance of `Std::Int`"),
			},
		},
		"resolve static add": {
			input: "1i8 + 5i8",
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
					value.Int8(6).ToValue(),
				},
			),
		},
		"add": {
			input: "a := 1i8; a + 5i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{
					value.Undefined,
					value.Int8(1).ToValue(),
					value.Int8(5).ToValue(),
				},
			),
		},
		"resolve static subtract": {
			input: "151i32 - 25i32 - 5i32",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 1, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Undefined,
					value.Int32(121).ToValue(),
				},
			),
		},
		"subtract": {
			input: "a := 151i32; a - 25i32 - 5i32",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBTRACT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.SUBTRACT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 1, 29)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 16),
				},
				[]value.Value{
					value.Undefined,
					value.Int32(151).ToValue(),
					value.Int32(25).ToValue(),
					value.Int32(5).ToValue(),
				},
			),
		},
		"resolve static multiply": {
			input: "45.5 * 2.5",
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
					value.Undefined,
					value.Float(113.75).ToValue(),
				},
			),
		},
		"multiply": {
			input: "a := 45.5; a * 2.5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{
					value.Undefined,
					value.Float(45.5).ToValue(),
					value.Float(2.5).ToValue(),
				},
			),
		},
		"resolve static divide": {
			input: "45.5 / .5",
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
					value.Float(91).ToValue(),
				},
			),
		},
		"divide": {
			input: "a := 45.5; a / .5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{
					value.Undefined,
					value.Float(45.5).ToValue(),
					value.Float(0.5).ToValue(),
				},
			),
		},
		"resolve static exponentiate": {
			input: "-2 ** 2",
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
					value.Undefined,
					value.SmallInt(-4).ToValue(),
				},
			),
		},
		"exponentiate": {
			input: "a := -2; a ** 2",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(-2).ToValue(),
					value.SmallInt(2).ToValue(),
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

func TestUnaryExpressions(t *testing.T) {
	tests := testTable{
		"resolve static negate": {
			input: "-5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(1, 1, 2)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(-5).ToValue(),
				},
			),
		},
		"resolve static plus": {
			input: "+5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(1, 1, 2)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(5).ToValue(),
				},
			),
		},
		"negate": {
			input: "a := 5; -a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(5).ToValue(),
				},
			),
		},
		"resolve static bitwise not": {
			input: "~10",
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
					value.Undefined,
					value.SmallInt(-11).ToValue(),
				},
			),
		},
		"resolve static logical not": {
			input: "!10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Undefined,
				},
			),
		},
		"logical not": {
			input: "a := 10; !a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(10).ToValue(),
				},
			),
		},
		"bitwise not": {
			input: "a := 10; ~a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.BITWISE_NOT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(10).ToValue(),
				},
			),
		},
		"unary plus": {
			input: "a := 10; +a",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.UNARY_PLUS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(10).ToValue(),
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

func TestComplexAssignmentLocals(t *testing.T) {
	tests := testTable{
		"increment": {
			input: "a := 1; a++",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.INCREMENT),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
				},
			),
		},
		"decrement": {
			input: "a := 1; a--",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DECREMENT),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
				},
			),
		},
		"add": {
			input: "a := 1; a += 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"subtract": {
			input: "a := 1; a -= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBTRACT),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"multiply": {
			input: "a := 1; a *= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"divide": {
			input: "a := 1; a /= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DIVIDE),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"exponentiate": {
			input: "a := 1; a **= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"modulo": {
			input: "a := 1; a %= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.MODULO),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"bitwise AND": {
			input: "a := 1; a &= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"bitwise OR": {
			input: "a := 1; a |= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"bitwise XOR": {
			input: "a := 1; a ^= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"left bitshift": {
			input: "a := 1; a <<= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"left logical bitshift": {
			input: "a := 1u64; a <<<= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.UInt64(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"right bitshift": {
			input: "a := 1; a >>= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"right logical bitshift": {
			input: "a := 1u64; a >>>= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 15),
				},
				[]value.Value{
					value.Undefined,
					value.UInt64(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"logic OR": {
			input: "var a: Int? = 1; a ||= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_IF), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 1, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 18),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"logic AND": {
			input: "var a: Int? = 1; a &&= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_UNLESS), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 1, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 18),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"nil coalesce": {
			input: "var a: Int? = 1; a ??= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.JUMP_IF_NIL), 0, 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 1, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 21),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
					value.SmallInt(3).ToValue(),
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

func TestComplexAssignmentInstanceVariables(t *testing.T) {
	tests := testTable{
		"increment": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a++
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(82, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
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
						L(P(0, 1, 1), P(82, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(82, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.INCREMENT),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(73, 6, 22)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"decrement": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a--
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(82, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
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
						L(P(0, 1, 1), P(82, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(82, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.DECREMENT),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(73, 6, 22)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 6),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"add": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a += 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.ADD),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"subtract": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a -= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.SUBTRACT),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"multiply": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a *= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.MULTIPLY),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"divide": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a /= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.DIVIDE),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"exponentiate": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a **= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(86, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
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
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.EXPONENTIATE),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(77, 6, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"modulo": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a %= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.MODULO),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"bitwise AND": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a &= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.BITWISE_AND),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"bitwise OR": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a |= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.BITWISE_OR),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"bitwise XOR": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a ^= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(85, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.BITWISE_XOR),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(76, 6, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"left bitshift": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a <<= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(86, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.LBITSHIFT),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(77, 6, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"left logical bitshift": {
			input: `
				class Foo
					var @a: UInt64
					init(@a); end

					def foo then @a <<<= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(90, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
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
						L(P(0, 1, 1), P(90, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(90, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(40, 4, 6), P(52, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.LOGIC_LBITSHIFT),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(60, 6, 6), P(81, 6, 27)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"right bitshift": {
			input: `
				class Foo
					var @a: Int
					init(@a); end

					def foo then @a >>= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(86, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(86, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(37, 4, 6), P(49, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.RBITSHIFT),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(57, 6, 6), P(77, 6, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"right logical bitshift": {
			input: `
				class Foo
					var @a: Int64
					init(@a); end

					def foo then @a >>>= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(89, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(89, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
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
						L(P(0, 1, 1), P(89, 7, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 13),
							bytecode.NewLineInfo(7, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("#init"),
								[]byte{
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.POP),
									byte(bytecode.NIL),
									byte(bytecode.POP),
									byte(bytecode.RETURN_SELF),
								},
								L(P(39, 4, 6), P(51, 4, 18)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 8),
								},
								1,
								0,
								[]value.Value{
									value.ToSymbol("a").ToValue(),
								},
							)),
							value.ToSymbol("#init").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.LOGIC_RBITSHIFT),
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(59, 6, 6), P(80, 6, 27)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(6, 8),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"logic OR": {
			input: `
				class Foo
					var @a: Int?

					def foo then @a ||= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.JUMP_IF), 0, 3,
									byte(bytecode.POP),
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(39, 5, 6), P(59, 5, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 11),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"nil coalesce": {
			input: `
				class Foo
					var @a: Int?

					def foo then @a ??= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), bytecode.DEF_CLASS_FLAG,
							byte(bytecode.GET_CONST8), 1,
							byte(bytecode.GET_CONST8), 2,
							byte(bytecode.SET_SUPERCLASS),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.JUMP_IF_NIL), 0, 3,
									byte(bytecode.JUMP), 0, 3,
									byte(bytecode.POP),
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(39, 5, 6), P(59, 5, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 14),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"logic AND": {
			input: `
				class Foo
					var @a: Int?

					def foo then @a &&= 3
				end
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
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
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
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 11),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
							value.ToSymbol("Std::Object").ToValue(),
						},
					)),
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<methodDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(68, 6, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 8),
							bytecode.NewLineInfo(6, 2),
						},
						[]value.Value{
							value.ToSymbol("Foo").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_IVAR8), 0,
									byte(bytecode.JUMP_UNLESS), 0, 3,
									byte(bytecode.POP),
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.SET_IVAR8), 0,
									byte(bytecode.RETURN),
								},
								L(P(39, 5, 6), P(59, 5, 26)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(5, 11),
								},
								[]value.Value{
									value.ToSymbol("a").ToValue(),
									value.SmallInt(3).ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
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

func TestBitwiseAnd(t *testing.T) {
	tests := testTable{
		"resolve static AND": {
			input: "23 & 10",
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
					value.Undefined,
					value.SmallInt(2).ToValue(),
				},
			),
		},
		"resolve static nested AND": {
			input: "23 & 15 & 46",
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
					value.Undefined,
					value.SmallInt(6).ToValue(),
				},
			),
		},
		"compile runtime AND": {
			input: "a := 23; a & 15 & 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 16),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(23).ToValue(),
					value.SmallInt(15).ToValue(),
					value.SmallInt(46).ToValue(),
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

func TestBitwiseAndNot(t *testing.T) {
	tests := testTable{
		"resolve static AND NOT": {
			input: "23 &~ 10",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(21).ToValue(),
				},
			),
		},
		"resolve static nested AND NOT": {
			input: "23 &~ 15 &~ 46",
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
					value.Undefined,
					value.SmallInt(16).ToValue(),
				},
			),
		},
		"compile runtime AND NOT": {
			input: "a := 23; a &~ 15 &~ 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.BITWISE_AND_NOT),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.BITWISE_AND_NOT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 1, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 16),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(23).ToValue(),
					value.SmallInt(15).ToValue(),
					value.SmallInt(46).ToValue(),
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

func TestBitwiseOr(t *testing.T) {
	tests := testTable{
		"resolve static OR": {
			input: "23 | 10",
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
					value.Undefined,
					value.SmallInt(31).ToValue(),
				},
			),
		},
		"resolve static nested OR": {
			input: "23 | 15 | 46",
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
					value.Undefined,
					value.SmallInt(63).ToValue(),
				},
			),
		},
		"compile runtime OR": {
			input: "a := 23; a | 15 | 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 16),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(23).ToValue(),
					value.SmallInt(15).ToValue(),
					value.SmallInt(46).ToValue(),
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

func TestBitwiseXor(t *testing.T) {
	tests := testTable{
		"resolve static XOR": {
			input: "23 ^ 10",
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
					value.Undefined,
					value.SmallInt(29).ToValue(),
				},
			),
		},
		"resolve static nested XOR": {
			input: "23 ^ 15 ^ 46",
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
					value.Undefined,
					value.SmallInt(54).ToValue(),
				},
			),
		},
		"compile runtime XOR": {
			input: "a := 23; a ^ 15 ^ 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 16),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(23).ToValue(),
					value.SmallInt(15).ToValue(),
					value.SmallInt(46).ToValue(),
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

func TestModulo(t *testing.T) {
	tests := testTable{
		"resolve static modulo": {
			input: "23 % 10",
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
					value.Undefined,
					value.SmallInt(3).ToValue(),
				},
			),
		},
		"resolve static nested modulo": {
			input: "24 % 15 % 2",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(1).ToValue(),
				},
			),
		},
		"compile runtime modulo": {
			input: "a := 24; a % 15 % 46",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.MODULO),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.MODULO),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 16),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(24).ToValue(),
					value.SmallInt(15).ToValue(),
					value.SmallInt(46).ToValue(),
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
