package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func TestLaxEqual(t *testing.T) {
	tests := testTable{
		"resolve static 25 =~ 25.0": {
			input: "25 =~ 25.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 =~ 25": {
			input: "25 =~ 25",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 =~ '25'": {
			input: "25 =~ '25'",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"compile runtime 24 =~ 98": {
			input: "a := 24; a =~ 98",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 24,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 98,
					byte(bytecode.LAX_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
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

func TestLaxNotEqual(t *testing.T) {
	tests := testTable{
		"resolve static 25 !~ 25.0": {
			input: "25 !~ 25.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 !~ 25": {
			input: "25 !~ 25",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 !~ '25'": {
			input: "25 !~ '25'",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"compile runtime 24 !~ 98": {
			input: "a := 24; a !~ 98",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 24,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 98,
					byte(bytecode.LAX_NOT_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
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

func TestEqual(t *testing.T) {
	tests := testTable{
		"resolve static 25 == 25.0": {
			input: "25 == 25.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
			err: error.ErrorList{
				error.NewWarning(L(P(0, 1, 1), P(1, 1, 2)), "this equality check is impossible, `25` cannot ever be equal to `25.0`"),
			},
		},
		"resolve static 25 == 25": {
			input: "25 == 25",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 == '25'": {
			input: "25 == '25'",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
			err: error.ErrorList{
				error.NewWarning(L(P(0, 1, 1), P(1, 1, 2)), "this equality check is impossible, `25` cannot ever be equal to `\"25\"`"),
			},
		},
		"compile runtime int 24 == 98": {
			input: "a := 24; a == 98",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 24,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 98,
					byte(bytecode.EQUAL_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
				},
			),
		},
		"compile runtime float 24.5 == 98.0": {
			input: "a := 24.5; a == 98.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
					value.Float(24.5).ToValue(),
					value.Float(98.0).ToValue(),
				},
			),
		},
		"compile runtime builtin 24i8 == 98i8": {
			input: "var a: Int8 = 24i8; a == 98i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT8), 24,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT8), 98,
					byte(bytecode.EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 1, 29)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
				},
			),
		},
		"compile runtime value 24 == 98": {
			input: "var a: any = 5; a == 98i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT8), 98,
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 1, 25)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(value.NewCallSiteInfo(symbol.OpEqual, 1)),
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

func TestNotEqual(t *testing.T) {
	tests := testTable{
		"resolve static 25 != 25.0": {
			input: "25 != 25.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
			err: error.ErrorList{
				error.NewWarning(L(P(0, 1, 1), P(1, 1, 2)), "this equality check is impossible, `25` cannot ever be equal to `25.0`"),
			},
		},
		"resolve static 25 != 25": {
			input: "25 != 25",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 != '25'": {
			input: "25 != '25'",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
			err: error.ErrorList{
				error.NewWarning(L(P(0, 1, 1), P(1, 1, 2)), "this equality check is impossible, `25` cannot ever be equal to `\"25\"`"),
			},
		},
		"compile runtime int 24 != 98": {
			input: "a := 24; a != 98",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 24,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 98,
					byte(bytecode.NOT_EQUAL_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
				},
			),
		},
		"compile runtime float 24.5 != 98.0": {
			input: "a := 24.5; a != 98.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NOT_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
					value.Float(24.5).ToValue(),
					value.Float(98.0).ToValue(),
				},
			),
		},
		"compile runtime builtin 24i8 != 98i8": {
			input: "var a: Int8 = 24i8; a != 98i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT8), 24,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT8), 98,
					byte(bytecode.NOT_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 1, 29)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
				},
			),
		},
		"compile runtime value 24 != 98": {
			input: "var a: any = 5; a != 98i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT8), 98,
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 1, 25)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 11),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(value.NewCallSiteInfo(symbol.OpEqual, 1)),
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

func TestStrictEqual(t *testing.T) {
	tests := testTable{
		"resolve static 25 === 25": {
			input: "25 === 25",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 === 25.0": {
			input: "25 === 25.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
			err: error.ErrorList{
				error.NewWarning(L(P(0, 1, 1), P(1, 1, 2)), "this strict equality check is impossible, `25` cannot ever be equal to `25.0`"),
			},
		},
		"resolve static 25 === '25'": {
			input: "25 === '25'",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
			err: error.ErrorList{
				error.NewWarning(L(P(0, 1, 1), P(1, 1, 2)), "this strict equality check is impossible, `25` cannot ever be equal to `\"25\"`"),
			},
		},
		"compile runtime 24 === 98": {
			input: "a := 24; a === 98",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 24,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 98,
					byte(bytecode.STRICT_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
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

func TestStrictNotEqual(t *testing.T) {
	tests := testTable{
		"resolve static 25 !== 25": {
			input: "25 !== 25",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 !== 25.0": {
			input: "25 !== 25.0",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
			err: error.ErrorList{
				error.NewWarning(L(P(0, 1, 1), P(1, 1, 2)), "this strict equality check is impossible, `25` cannot ever be equal to `25.0`"),
			},
		},
		"resolve static 25 !== '25'": {
			input: "25 !== '25'",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
			err: error.ErrorList{
				error.NewWarning(L(P(0, 1, 1), P(1, 1, 2)), "this strict equality check is impossible, `25` cannot ever be equal to `\"25\"`"),
			},
		},
		"compile runtime 24 !== 98": {
			input: "a := 24; a !== 98",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 24,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 98,
					byte(bytecode.STRICT_NOT_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
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

func TestGreaterThan(t *testing.T) {
	tests := testTable{
		"resolve static 3 > 3": {
			input: "3 > 3",
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
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 > 3": {
			input: "25 > 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25.2 > 25": {
			input: "25.2 > 25",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 7 > 20": {
			input: "7 > 20",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"compile runtime int": {
			input: "a := 24; a > 98",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT_8), 24,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 98,
					byte(bytecode.GREATER_INT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
				},
			),
		},
		"compile runtime float": {
			input: "a := 24.5; a > 98.5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_FLOAT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
					value.Float(24.5).ToValue(),
					value.Float(98.5).ToValue(),
				},
			),
		},
		"compile runtime builtin": {
			input: "a := 24i8; a > 98i8",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_INT8), 24,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT8), 98,
					byte(bytecode.GREATER),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 1, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Undefined,
				},
			),
		},
		"compile runtime value": {
			input: `
				module Foo
					def >(other: Foo | Int): bool
						true
					end
				end

				Foo > 98
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
					byte(bytecode.LOAD_INT_8), 98,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(92, 8, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
					bytecode.NewLineInfo(8, 7),
				},
				[]value.Value{
					value.Ref(
						vm.NewBytecodeFunctionNoParams(
							namespaceDefinitionsSymbol,
							[]byte{
								byte(bytecode.GET_CONST8), 0,
								byte(bytecode.LOAD_VALUE8), 1,
								byte(bytecode.DEF_NAMESPACE), 0,
								byte(bytecode.NIL),
								byte(bytecode.RETURN),
							},
							L(P(0, 1, 1), P(92, 8, 13)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 6),
								bytecode.NewLineInfo(8, 2),
							},
							[]value.Value{
								value.ToSymbol("Root").ToValue(),
								value.ToSymbol("Foo").ToValue(),
							},
						),
					),
					value.Ref(
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
							L(P(0, 1, 1), P(92, 8, 13)),
							bytecode.LineInfoList{
								bytecode.NewLineInfo(1, 9),
								bytecode.NewLineInfo(8, 2),
							},
							[]value.Value{
								value.ToSymbol("Foo").ToValue(),
								value.Ref(
									vm.NewBytecodeFunction(
										symbol.OpGreaterThan,
										[]byte{
											byte(bytecode.TRUE),
											byte(bytecode.RETURN),
										},
										L(P(21, 3, 6), P(69, 5, 8)),
										bytecode.LineInfoList{
											bytecode.NewLineInfo(4, 1),
											bytecode.NewLineInfo(5, 1),
										},
										1,
										0,
										nil,
									),
								),
								value.ToSymbol(">").ToValue(),
							},
						),
					),
					value.ToSymbol("Foo").ToValue(),
					value.Ref(value.NewCallSiteInfo(symbol.OpGreaterThan, 1)),
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

func TestGreaterThanEqual(t *testing.T) {
	tests := testTable{
		"resolve static 3 >= 3": {
			input: "3 >= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 >= 3": {
			input: "25 >= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25.2 >= 25": {
			input: "25.2 >= 25",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 7 >= 20": {
			input: "7 >= 20",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"compile runtime 24 >= 98": {
			input: "a := 24; a >= 98",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(24).ToValue(),
					value.SmallInt(98).ToValue(),
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

func TestLessThan(t *testing.T) {
	tests := testTable{
		"resolve static 3 < 3": {
			input: "3 < 3",
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
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 < 3": {
			input: "25 < 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25.2 < 25": {
			input: "25.2 < 25",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 7 < 20": {
			input: "7 < 20",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"compile runtime 24 < 98": {
			input: "a := 24; a < 98",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(24).ToValue(),
					value.SmallInt(98).ToValue(),
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

func TestLessThanEqual(t *testing.T) {
	tests := testTable{
		"resolve static 3 <= 3": {
			input: "3 <= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25 <= 3": {
			input: "25 <= 3",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 25.2 <= 25": {
			input: "25.2 <= 25",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"resolve static 7 <= 20": {
			input: "7 <= 20",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{value.Undefined},
			),
		},
		"compile runtime 24 <= 98": {
			input: "a := 24; a <= 98",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 13),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(24).ToValue(),
					value.SmallInt(98).ToValue(),
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
