package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestGoExpression(t *testing.T) {
	tests := testTable{
		"with a single expression": {
			input: "go println('foo')",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.CLOSURE), 0xff,
					byte(bytecode.GO),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<closure>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.CALL_METHOD8), 2,
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(16, 1, 17)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
							value.Ref(value.NewCallSiteInfo(
								value.ToSymbol("println"),
								1,
							)),
						},
					)),
				},
			),
		},
		"with outer variables": {
			input: `
				a := 5
				go
					println("foo")
					println(a)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.CLOSURE), 2, 1, 0xff,
					byte(bytecode.GO),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(62, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionWithUpvalues(
						value.ToSymbol("<closure>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.CALL_METHOD8), 2,
							byte(bytecode.POP),
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.UNDEFINED),
							byte(bytecode.GET_UPVALUE_0),
							byte(bytecode.NEW_ARRAY_TUPLE8), 1,
							byte(bytecode.CALL_METHOD8), 3,
							byte(bytecode.RETURN),
						},
						L(P(16, 3, 5), P(54, 5, 16)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(4, 6),
							bytecode.NewLineInfo(5, 9),
						},
						0,
						0,
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
							value.Ref(value.NewCallSiteInfo(
								value.ToSymbol("println"),
								1,
							)),
							value.Ref(value.NewCallSiteInfo(
								value.ToSymbol("println"),
								1,
							)),
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

func TestForInExpression(t *testing.T) {
	tests := testTable{
		"int literal": {
			input: `
				for i in 20
					println(i)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 20,
					byte(bytecode.JUMP_UNLESS_ILT), 0, 15,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INCREMENT_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(4, 9),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"int value": {
			input: `
				j := 20
				for i in j
					println(i)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 3,
					byte(bytecode.LOAD_INT_8), 20,
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 15,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.INCREMENT_INT),
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(51, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 10),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 12),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"range literal": {
			input: `
				for i in 5...20
					println(i)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 20,
					byte(bytecode.JUMP_UNLESS_ILE), 0, 15,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INCREMENT_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(44, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(4, 9),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"range value": {
			input: `
				r := 5...20
				for i in r
					println(i)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 4,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.SET_LOCAL_4),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.JUMP_UNLESS_ILE), 0, 15,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 3,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.INCREMENT_INT),
					byte(bytecode.SET_LOCAL_4),
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(55, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 16),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 12),
				},
				[]value.Value{
					value.Ref(value.NewClosedRange(value.SmallInt(5).ToValue(), value.SmallInt(20).ToValue())),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("end"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("start"), 0)),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"iterate": {
			input: `
				for i in [1, 2, 3]
					println(i)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 13,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 17,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(4, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"with a pattern": {
			input: `
				for %[a, b] in %[%[1, 2], %[3, 4], %[5, 6]]
					println(a + b)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 3,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 59,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS_NP), 0, 33,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.INT_2),
					byte(bytecode.EQUAL_INT),
					byte(bytecode.JUMP_UNLESS_NP), 0, 24,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.INT_0),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 13,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_IF), 0, 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.THROW),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 4,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.ADD_INT),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 63,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 50),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(3, 10),
					bytecode.NewLineInfo(4, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.Ref(&value.ArrayTuple{
							value.SmallInt(1).ToValue(),
							value.SmallInt(2).ToValue(),
						}),
						value.Ref(&value.ArrayTuple{
							value.SmallInt(3).ToValue(),
							value.SmallInt(4).ToValue(),
						}),
						value.Ref(&value.ArrayTuple{
							value.SmallInt(5).ToValue(),
							value.SmallInt(6).ToValue(),
						}),
					}),
					value.Ref(value.TupleMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.Ref(value.NewError(
						value.PatternNotMatchedErrorClass,
						"assigned value does not match the pattern defined in for in loop",
					)),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"with break": {
			input: `
				for i in [1, 2, 3, 4, 5]
					println(i)
					break if i > 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 30,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 10,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.JUMP), 0, 12,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 34,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(4, 16),
					bytecode.NewLineInfo(5, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"with break with value": {
			input: `
				for i in [1, 2, 3, 4, 5]
					println(i)
					break :foo if i > 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 30,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 10,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.JUMP), 0, 12,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 34,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(78, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(4, 16),
					bytecode.NewLineInfo(5, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
					value.ToSymbol("foo").ToValue(),
				},
			),
		},
		"with labeled break": {
			input: `
				$foo: for i in [1, 2, 3, 4, 5]
					println(i)
					break[foo] if i > 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 30,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 10,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.JUMP), 0, 12,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 34,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(84, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(4, 16),
					bytecode.NewLineInfo(5, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"with continue": {
			input: `
				for i in [1, 2, 3, 4, 5]
					continue if i > 2
					println(i)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 22,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 4,
					byte(bytecode.LOOP), 0, 13,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 26,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"with labeled continue": {
			input: `
				$foo: for i in [1, 2, 3, 4, 5]
					continue[foo] if i > 2
					println(i)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 22,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_2),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 4,
					byte(bytecode.LOOP), 0, 13,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 26,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(87, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"nested with break": {
			input: `
				for c in ['a', 'b', 'c', 'd']
					for i in [1, 2, 3, 4, 5]
						break if i > 2
						println(c, i)
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 4,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 44,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.FOR_IN_BUILTIN), 0, 27,
					byte(bytecode.SET_LOCAL_4),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.INT_2),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 8,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 2,
					byte(bytecode.JUMP), 0, 18,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.NEW_ARRAY_TUPLE8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 31,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 48,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(122, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(4, 13),
					bytecode.NewLineInfo(5, 9),
					bytecode.NewLineInfo(6, 8),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.Ref(value.String("a")),
						value.Ref(value.String("b")),
						value.Ref(value.String("c")),
						value.Ref(value.String("d")),
					}),
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"nested with labeled break": {
			input: `
				$foo: for c in ['a', 'b', 'c', 'd']
					for i in [1, 2, 3, 4, 5]
						break[foo] if i > 2
						println(c, i)
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 4,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 44,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.FOR_IN_BUILTIN), 0, 27,
					byte(bytecode.SET_LOCAL_4),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.INT_2),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 8,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 4,
					byte(bytecode.JUMP), 0, 26,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.NEW_ARRAY_TUPLE8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 31,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 48,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(133, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(4, 13),
					bytecode.NewLineInfo(5, 9),
					bytecode.NewLineInfo(6, 8),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.Ref(value.String("a")),
						value.Ref(value.String("b")),
						value.Ref(value.String("c")),
						value.Ref(value.String("d")),
					}),
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"nested with continue": {
			input: `
				for c in ['a', 'b', 'c', 'd']
					for i in [1, 2, 3, 4, 5]
						continue if i > 2
						println(c, i)
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 4,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 40,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.FOR_IN_BUILTIN), 0, 23,
					byte(bytecode.SET_LOCAL_4),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.INT_2),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 4,
					byte(bytecode.LOOP), 0, 13,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.NEW_ARRAY_TUPLE8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 44,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(125, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(4, 9),
					bytecode.NewLineInfo(5, 9),
					bytecode.NewLineInfo(6, 8),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.Ref(value.String("a")),
						value.Ref(value.String("b")),
						value.Ref(value.String("c")),
						value.Ref(value.String("d")),
					}),
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"nested with labeled continue": {
			input: `
				$foo: for c in ['a', 'b', 'c', 'd']
					for i in [1, 2, 3, 4, 5]
						continue[foo] if i > 2
						println(c, i)
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 4,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 43,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.FOR_IN_BUILTIN), 0, 26,
					byte(bytecode.SET_LOCAL_4),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.INT_2),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 7,
					byte(bytecode.LEAVE_SCOPE16), 4, 2,
					byte(bytecode.LOOP), 0, 25,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_4),
					byte(bytecode.NEW_ARRAY_TUPLE8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 47,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(136, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(4, 12),
					bytecode.NewLineInfo(5, 9),
					bytecode.NewLineInfo(6, 8),
					bytecode.NewLineInfo(7, 9),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.Ref(value.String("a")),
						value.Ref(value.String("b")),
						value.Ref(value.String("c")),
						value.Ref(value.String("d")),
					}),
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(4).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
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

func TestReturnExpression(t *testing.T) {
	tests := testTable{
		"return a value": {
			input: "return 5",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_5),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(7, 1, 8), P(7, 1, 8)), "values returned in void context will be ignored"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestAwaitExpression(t *testing.T) {
	tests := testTable{
		"await in a synchronous context": {
			input: "await timeout(2.seconds)",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.INT_2),
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.AWAIT_SYNC),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 1, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("seconds"), 0)),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("timeout"), 2)),
				},
			),
		},
		"await in an asynchronous context": {
			input: `
				async def foo
					await timeout(2.seconds)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(56, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(56, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.PROMISE),
									byte(bytecode.RETURN),
									byte(bytecode.INT_2),
									byte(bytecode.CALL_METHOD8), 0,
									byte(bytecode.UNDEFINED),
									byte(bytecode.CALL_SELF8), 1,
									byte(bytecode.AWAIT),
									byte(bytecode.AWAIT_RESULT),
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(55, 4, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 2),
									bytecode.NewLineInfo(4, 1),
									bytecode.NewLineInfo(3, 8),
									bytecode.NewLineInfo(4, 1),
								},
								1,
								1,
								[]value.Value{
									value.Ref(value.NewCallSiteInfo(value.ToSymbol("seconds"), 0)),
									value.Ref(value.NewCallSiteInfo(value.ToSymbol("timeout"), 2)),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
				},
			),
		},
		"await_sync in an asynchronous context": {
			input: `
				async def foo
					await_sync timeout(2.seconds)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(61, 4, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(4, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunction(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.GET_LOCAL_1),
									byte(bytecode.PROMISE),
									byte(bytecode.RETURN),
									byte(bytecode.INT_2),
									byte(bytecode.CALL_METHOD8), 0,
									byte(bytecode.UNDEFINED),
									byte(bytecode.CALL_SELF8), 1,
									byte(bytecode.AWAIT_SYNC),
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(60, 4, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(2, 2),
									bytecode.NewLineInfo(4, 1),
									bytecode.NewLineInfo(3, 7),
									bytecode.NewLineInfo(4, 1),
								},
								1,
								1,
								[]value.Value{
									value.Ref(value.NewCallSiteInfo(value.ToSymbol("seconds"), 0)),
									value.Ref(value.NewCallSiteInfo(value.ToSymbol("timeout"), 2)),
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

func TestModifierForIn(t *testing.T) {
	tests := testTable{
		"iterate": {
			input: `println(i) for i in [1, 2, 3]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 13,
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 17,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 1, 29)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 28),
				},
				[]value.Value{
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(2).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
		},
		"with a pattern": {
			input: `println(a + b) for %[a, b] in %[%[1, 2], %[3, 4], %[5, 6]]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 3,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.FOR_IN_BUILTIN), 0, 59,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS_NP), 0, 33,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.INT_2),
					byte(bytecode.EQUAL_INT),
					byte(bytecode.JUMP_UNLESS_NP), 0, 24,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.INT_0),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 13,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.INT_1),
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_3),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_IF), 0, 2,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.THROW),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 4,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_3),
					byte(bytecode.ADD_INT),
					byte(bytecode.NEW_ARRAY_TUPLE8), 1,
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 63,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(57, 1, 58)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 73),
				},
				[]value.Value{
					value.Ref(&value.ArrayTuple{
						value.Ref(&value.ArrayTuple{
							value.SmallInt(1).ToValue(),
							value.SmallInt(2).ToValue(),
						}),
						value.Ref(&value.ArrayTuple{
							value.SmallInt(3).ToValue(),
							value.SmallInt(4).ToValue(),
						}),
						value.Ref(&value.ArrayTuple{
							value.SmallInt(5).ToValue(),
							value.SmallInt(6).ToValue(),
						}),
					}),
					value.Ref(value.TupleMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.Ref(value.NewError(
						value.PatternNotMatchedErrorClass,
						"assigned value does not match the pattern defined in for in loop",
					)),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
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

func TestIfExpression(t *testing.T) {
	tests := testTable{
		"resolve static condition with empty then and else": {
			input: `if false; end`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(3, 1, 4), P(7, 1, 8)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"empty then and else": {
			input: "a := true; if a; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.TRUE),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 4,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 14),
				},
				[]value.Value{},
			),
		},
		"resolve static condition with then branch": {
			input: `
				if true
					10
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(8, 2, 8), P(11, 2, 11)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"resolve static condition with then branch to nil": {
			input: `
				if false
					10
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(8, 2, 8), P(12, 2, 12)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewWarning(L(P(19, 3, 6), P(21, 3, 8)), "unreachable code"),
			},
		},
		"resolve static condition with then and else branches": {
			input: `
				if false
					10
				else
					5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_5),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(8, 2, 8), P(12, 2, 12)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewWarning(L(P(19, 3, 6), P(21, 3, 8)), "unreachable code"),
			},
		},
		"with then branch": {
			input: `
				var a: Int? = 5
				if a
					a = a * 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_2),
					byte(bytecode.MULTIPLY_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(52, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{},
			),
		},
		"with then and else branches": {
			input: `
				var a: Int? = 5
				if a
					a = a * 2
				else
					a = 30
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_2),
					byte(bytecode.MULTIPLY_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.LOAD_INT_8), 30,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{},
			),
		},
		"is an expression": {
			input: `
				var a: Int? = 5
				b := if a
					"foo"
				else
					5
				end
				b
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS), 0, 4,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(75, 8, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(8, 2),
				},
				[]value.Value{
					value.Ref(value.String("foo")),
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

func TestUnlessExpression(t *testing.T) {
	tests := testTable{
		"resolve static condition with empty then and else": {
			input: "unless true; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(7, 1, 8), P(10, 1, 11)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"empty then and else": {
			input: "a := true; unless a; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.TRUE),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF), 0, 4,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 1, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 14),
				},
				[]value.Value{},
			),
		},
		"resolve static condition with then branch": {
			input: `
				unless false
					10
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(33, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(12, 2, 12), P(16, 2, 16)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"resolve static condition with then branch to nil": {
			input: `
				unless true
					10
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(12, 2, 12), P(15, 2, 15)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewWarning(L(P(22, 3, 6), P(24, 3, 8)), "unreachable code"),
			},
		},
		"resolve static condition with then and else branches": {
			input: `
				unless true
					10
				else
					5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.INT_5),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(48, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(12, 2, 12), P(15, 2, 15)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewWarning(L(P(22, 3, 6), P(24, 3, 8)), "unreachable code"),
			},
		},
		"with then branch": {
			input: `
				var a: Int? = 5
				unless a
					a = 30
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF), 0, 7,
					byte(bytecode.LOAD_INT_8), 30,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(53, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{},
			),
		},
		"with then and else branches": {
			input: `
				var a: Int? = 5
				unless a
					a = 30
				else
					a = a * 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF), 0, 7,
					byte(bytecode.LOAD_INT_8), 30,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.JUMP), 0, 5,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_2),
					byte(bytecode.MULTIPLY_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(77, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{},
			),
		},
		"is an expression": {
			input: `
				var a: Int? = 5
				b := unless a
					"foo"
				else
					5
				end
				b
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF), 0, 4,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.INT_5),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(79, 8, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(8, 2),
				},
				[]value.Value{
					value.Ref(value.String("foo")),
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

func TestBreak(t *testing.T) {
	tests := testTable{
		"in top level": {
			input: "break",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(4, 1, 5)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"in top level with a label": {
			input: "break[foo]",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(8, 1, 9)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"nonexistent label": {
			input: `
				loop
					break[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(15, 3, 6), P(23, 3, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to an expression": {
			input: `
				loop
					$foo: 1 + 2
					break[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(32, 4, 6), P(40, 4, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to a different loop": {
			input: `
				$foo: loop
					println("foo")
				end

				loop
					break[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(59, 7, 6), P(67, 7, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
				diagnostic.NewWarning(L(P(49, 6, 5), P(76, 8, 7)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestContinue(t *testing.T) {
	tests := testTable{
		"in top level": {
			input: "continue",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(7, 1, 8)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"in top level with a label": {
			input: "continue[foo]",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(11, 1, 12)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"nonexistent label": {
			input: `
				loop
					continue[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(15, 3, 6), P(26, 3, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to an expression": {
			input: `
				loop
					$foo: 1 + 2
					continue[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(32, 4, 6), P(43, 4, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to a different loop": {
			input: `
				$foo: loop
					println("foo")
				end

				loop
					continue[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(59, 7, 6), P(70, 7, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
				diagnostic.NewWarning(L(P(49, 6, 5), P(79, 8, 7)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestLoopExpression(t *testing.T) {
	tests := testTable{
		"empty body": {
			input: `
				loop
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOOP), 0, 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 3, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{},
			),
		},
		"with a body": {
			input: `
				a := 0
				loop
					a = a + 1
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 9,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{},
			),
		},
		"with continue": {
			input: `
				a := 0
				loop
					a = a + 1
					continue
					println("foo")
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 17,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(77, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 5),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(55, 6, 6), P(68, 6, 19)), "unreachable code"),
			},
		},
		"with labeled continue": {
			input: `
				a := 0
				$foo: loop
					a = a + 1
					continue[foo]
					println("foo")
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 17,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(88, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(7, 5),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(66, 6, 6), P(79, 6, 19)), "unreachable code"),
			},
		},
		"continue in a nested loop": {
			input: `
			 	j := 0
				loop
					j += 1
					i := 0
					loop
						continue if i >= 5
						i += 1
					end
					continue if j >= 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_IGE), 0, 4,
					byte(bytecode.LOOP), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 18,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_IGE), 0, 9,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 36,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 47,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(134, 11, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 9),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(10, 15),
					bytecode.NewLineInfo(11, 8),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(108, 10, 6), P(125, 10, 23)), "unreachable code"),
			},
		},

		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: loop
					j += 1
					i := 0
					loop
						continue[foo] if i >= 5
						i += 1
					end
					continue if j >= 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_IGE), 0, 7,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 17,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_IGE), 0, 9,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 39,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 50,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(145, 11, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 12),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(10, 15),
					bytecode.NewLineInfo(11, 8),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(119, 10, 6), P(136, 10, 23)), "unreachable code"),
			},
		},
		"with break": {
			input: `
				a := 0
				loop
					a = a + 1
					break if a > 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
		},
		"with a labeled break": {
			input: `
				a := 0
				$foo: loop
					a = a + 1
					break[foo] if a > 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(74, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
		},
		"break in a nested loop": {
			input: `
			 	j := 0
				loop
					j += 1
					i := 0
					loop
						break if i >= 5
						i += 1
					end
					break if j >= 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_IGE), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 10,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 19,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_IGE), 0, 10,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 11,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 49,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(128, 11, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 10),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(10, 16),
					bytecode.NewLineInfo(11, 8),
				},
				[]value.Value{},
			),
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$outer: loop
					j += 1
					i := 0
					loop
						break[outer] if i >= 5
						i += 1
					end
					break if j >= 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_IGE), 0, 8,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 34,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 22,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_IGE), 0, 10,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 11,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 52,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(143, 11, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 13),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(10, 16),
					bytecode.NewLineInfo(11, 8),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(120, 10, 6), P(134, 10, 20)), "unreachable code"),
			},
		},
		"break with value": {
			input: `
				a := 0
				loop
					a = a + 1
					break true if a > 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_IGT), 0, 7,
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
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

func TestLogicalOrOperator(t *testing.T) {
	tests := testTable{
		"simple": {
			input: `
				a := "foo"
				a || true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(value.String("foo")),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(20, 3, 5), P(20, 3, 5)), "this condition will always have the same result since type `Std::String` is truthy"),
				diagnostic.NewWarning(L(P(25, 3, 10), P(28, 3, 13)), "unreachable code"),
			},
		},
		"nested": {
			input: `
				a := "foo"
				a || true || 3
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_IF_NP), 0, 2,
					// falsy 1
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// truthy 1
					byte(bytecode.JUMP_IF_NP), 0, 2,
					// falsy 2
					byte(bytecode.POP),
					byte(bytecode.INT_3),
					// truthy 2
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 12),
				},
				[]value.Value{
					value.Ref(value.String("foo")),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(20, 3, 5), P(20, 3, 5)), "this condition will always have the same result since type `Std::String` is truthy"),
				diagnostic.NewWarning(L(P(25, 3, 10), P(28, 3, 13)), "unreachable code"),
				diagnostic.NewWarning(L(P(20, 3, 5), P(28, 3, 13)), "this condition will always have the same result since type `Std::String` is truthy"),
				diagnostic.NewWarning(L(P(33, 3, 18), P(33, 3, 18)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestLogicalAndOperator(t *testing.T) {
	tests := testTable{
		"simple": {
			input: `
				a := "foo"
				a && true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					// truthy
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// falsy
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(value.String("foo")),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(20, 3, 5), P(20, 3, 5)), "this condition will always have the same result since type `Std::String` is truthy"),
			},
		},
		"nested": {
			input: `
				a := "foo"
				a && true && 3
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					// truthy 1
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// falsy 1
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					// truthy 2
					byte(bytecode.POP),
					byte(bytecode.INT_3),
					// falsy 2
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 12),
				},
				[]value.Value{
					value.Ref(value.String("foo")),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(20, 3, 5), P(20, 3, 5)), "this condition will always have the same result since type `Std::String` is truthy"),
				diagnostic.NewWarning(L(P(20, 3, 5), P(28, 3, 13)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestNilCoalescingOperator(t *testing.T) {
	tests := testTable{
		"simple": {
			input: `
				a := "foo"
				a ?? true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS_NNP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
				},
				[]value.Value{
					value.Ref(value.String("foo")),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(20, 3, 5), P(20, 3, 5)), "this condition will always have the same result since type `Std::String` can never be nil"),
				diagnostic.NewWarning(L(P(25, 3, 10), P(28, 3, 13)), "unreachable code"),
			},
		},
		"nested": {
			input: `
				a := "foo"
				a ?? true ?? 3
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.JUMP_UNLESS_NNP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS_NNP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.INT_3),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(34, 3, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 12),
				},
				[]value.Value{
					value.Ref(value.String("foo")),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(20, 3, 5), P(20, 3, 5)), "this condition will always have the same result since type `Std::String` can never be nil"),
				diagnostic.NewWarning(L(P(25, 3, 10), P(28, 3, 13)), "unreachable code"),
				diagnostic.NewWarning(L(P(20, 3, 5), P(28, 3, 13)), "this condition will always have the same result since type `Std::String` can never be nil"),
				diagnostic.NewWarning(L(P(33, 3, 18), P(33, 3, 18)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestNumericFor(t *testing.T) {
	tests := testTable{
		"for without initialiser, condition, increment and body": {
			input: `
				fornum ;;
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOOP), 0, 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 3, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{},
			),
		},
		"for without initialiser, condition and increment": {
			input: `
				a := 0
				fornum ;;
					a += 1
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 9,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{},
			),
		},
		"for with break": {
			input: `
				a := 0
				fornum ;;
					a += 1
					break if a > 10
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 22,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(66, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 14),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
		},
		"for with labeled break": {
			input: `
				a := 0
				$foo: fornum ;;
					a += 1
					break[foo] if a > 10
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 22,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(77, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 14),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
		},
		"nested for with continue": {
			input: `
				fornum a := 0;;
					fornum ;; a += 1
						continue if a > 10
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 22,
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(84, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 14),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 7),
				},
				[]value.Value{},
			),
		},
		"nested for with a labeled continue": {
			input: `
				$foo: fornum a := 0;;
					fornum ;; a += 1
						continue[foo] if a > 10
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 13,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 22,
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(95, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 14),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 7),
				},
				[]value.Value{},
			),
		},
		"nested for with break": {
			input: `
				fornum a := 0;;
					fornum ;; a += 1
						break if a > 10
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 11,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 22,
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(81, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 14),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 7),
				},
				[]value.Value{},
			),
		},
		"nested for with a labeled break": {
			input: `
				$foo: fornum a := 0;;
					fornum ;; a += 1
						break[foo] if a > 10
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 10,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.JUMP), 0, 17,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 25,
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(92, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 17),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 7),
				},
				[]value.Value{},
			),
		},

		"for with break with value": {
			input: `
				fornum a := 0;;
					a += 1
					break 5 if a > 10
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.LOAD_INT_8), 10,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 10,
					byte(bytecode.INT_5),
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.JUMP), 0, 10,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 25,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 17),
					bytecode.NewLineInfo(5, 7),
				},
				[]value.Value{},
			),
		},
		"for with initialiser, without condition and increment": {
			input: `
				fornum a := 0;;
					a += 1
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 9,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 7),
				},
				[]value.Value{},
			),
		},
		"for with initialiser, condition, without increment": {
			input: `
				fornum a := 0; a < 5;
					a += 1
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 14,
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 7),
				},
				[]value.Value{},
			),
		},
		"for with initialiser, condition and increment": {
			input: `
				a := 0
				fornum i := 0; i < 5; i += 1
					a += i
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 13,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOOP), 0, 18,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 7),
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

func TestModifierWhile(t *testing.T) {
	tests := testTable{
		"single line": {
			input: `
			  i := 0
				i += 1 while i < 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 14,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 15),
				},
				[]value.Value{},
			),
		},
		"multiline": {
			input: `
			  i := 0
				do
					i += 1
				end while i < 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 14,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(51, 5, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 4),
				},
				[]value.Value{},
			),
		},
		"with break": {
			input: `
			  i := 0
				do
					i += 1
					break if i < 5
				end while true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(70, 6, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(66, 6, 15), P(69, 6, 18)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"with labeled break": {
			input: `
				i := 0
				$foo: do
					i += 1
					break[foo] if i < 5
				end while true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(80, 6, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(76, 6, 15), P(79, 6, 18)), "this condition will always have the same result since type `true` is truthy"),
			},
		},

		"with break with value": {
			input: `
				i := 0
				do
					i += 1
					break true if i < 5
				end while true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(74, 6, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(70, 6, 15), P(73, 6, 18)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"continue in a nested loop": {
			input: `
			 	j := 0
				do
					j += 1
					i := 0
					do
						continue if i + j > 8
						i += 1
					end while i < 5
				end while j < 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(133, 10, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 13),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 2),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 5),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(10, 4),
				},
				[]value.Value{},
			),
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						continue[foo] if i + j > 8
						i += 1
					end while i < 5
				end while j < 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 8,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 18,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 48,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(144, 10, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 16),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 2),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 5),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(10, 4),
				},
				[]value.Value{},
			),
		},
		"break in a nested loop": {
			input: `
			 	j := 0
				do
					j += 1
					i := 0
					do
						break if i + j > 8
						i += 1
					end while i < 5
				end while j < 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 15,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(130, 10, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 13),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 2),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 5),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(10, 4),
				},
				[]value.Value{},
			),
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						break[foo] if i + j > 8
						i += 1
					end while i < 5
				end while j < 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 8,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 27,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 48,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(141, 10, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 16),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 2),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 5),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(10, 4),
				},
				[]value.Value{},
			),
		},
		"static infinite": {
			input: `
				do
					println("foo")
				end while true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 9,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 4, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 5),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(42, 4, 15), P(45, 4, 18)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"static one iteration": {
			input: `
				do
					println("foo")
				end while false
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 4, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(42, 4, 15), P(46, 4, 19)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestWhile(t *testing.T) {
	tests := testTable{
		"with a body": {
			input: `
			  i := 0
				while i < 5
					i += 1
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 14,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(48, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 4),
				},
				[]value.Value{},
			),
		},
		"with break": {
			input: `
			  i := 0
				while true
					i += 1
					break if i < 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(67, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(23, 3, 11), P(26, 3, 14)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"with labeled break": {
			input: `
			  i := 0
				$foo: while true
					i += 1
					break[foo] if i < 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(78, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(29, 3, 17), P(32, 3, 20)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"with break with value": {
			input: `
			  i := 0
				while true
					i += 1
					break true if i < 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(72, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(23, 3, 11), P(26, 3, 14)), "this condition will always have the same result since type `true` is truthy"),
			},
		},

		"continue in a nested loop": {
			input: `
			 	j := 0
				while j < 5
					j += 1
					i := 0
					while i < 5
						continue if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 38,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 22,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 18,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LOOP), 0, 43,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(127, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 7),
					bytecode.NewLineInfo(7, 13),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 7),
				},
				[]value.Value{},
			),
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: while j < 5
					j += 1
					i := 0
					while i < 5
						continue[foo] if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 38,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 22,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 31,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LOOP), 0, 43,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(138, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 7),
					bytecode.NewLineInfo(7, 13),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 7),
				},
				[]value.Value{},
			),
		},
		"break in a nested loop": {
			input: `
			 	j := 0
				while j < 5
					j += 1
					i := 0
					while i < 5
						break if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 38,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 22,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LOOP), 0, 43,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(124, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 7),
					bytecode.NewLineInfo(7, 13),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 7),
				},
				[]value.Value{},
			),
		},

		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: while j < 5
					j += 1
					i := 0
					while i < 5
						break[foo] if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 41,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 25,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 8,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 15,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.LOOP), 0, 46,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(135, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 7),
					bytecode.NewLineInfo(7, 16),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 7),
				},
				[]value.Value{},
			),
		},
		"without a body": {
			input: `
				i := 0
				while i < 5; end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 10,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 3, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 12),
				},
				[]value.Value{},
			),
		},
		"static infinite": {
			input: `
				while true
					println("foo")
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 9,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 5),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(11, 2, 11), P(14, 2, 14)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"static impossible": {
			input: `
				while false
					println("foo")
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(44, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(11, 2, 11), P(15, 2, 15)), "this loop will never execute since type `false` is falsy"),
				diagnostic.NewWarning(L(P(22, 3, 6), P(36, 3, 20)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestModifierUntil(t *testing.T) {
	tests := testTable{
		"single line": {
			input: `
			  i := 0
				i += 1 until i >= 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 14,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(36, 3, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 15),
				},
				[]value.Value{},
			),
		},
		"multiline": {
			input: `
			  i := 0
				do
					i += 1
				end until i >= 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 14,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(52, 5, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 4),
				},
				[]value.Value{},
			),
		},
		"with break": {
			input: `
			  i := 0
				do
					i += 1
					break if i < 5
				end until false
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 6, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(66, 6, 15), P(70, 6, 19)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"with labeled break": {
			input: `
			  i := 0
				$foo: do
					i += 1
					break[foo] if i < 5
				end until false
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(82, 6, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(77, 6, 15), P(81, 6, 19)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"with break with value": {
			input: `
			  i := 0
				do
					i += 1
					break true if i < 5
				end until false
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 6, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(71, 6, 15), P(75, 6, 19)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"continue in a nested loop": {
			input: `
			 	j := 0
				do
					j += 1
					i := 0
					do
						continue if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(135, 10, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 13),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 2),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 5),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(10, 4),
				},
				[]value.Value{},
			),
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						continue[foo] if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 8,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 18,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 48,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(146, 10, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 16),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 2),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 5),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(10, 4),
				},
				[]value.Value{},
			),
		},
		"break in a nested loop": {
			input: `
			 	j := 0
				do
					j += 1
					i := 0
					do
						break if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 15,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(132, 10, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 13),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 2),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 5),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(10, 4),
				},
				[]value.Value{},
			),
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						break[foo] if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 8,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 27,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 48,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(143, 10, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(7, 16),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 2),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 5),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(10, 4),
				},
				[]value.Value{},
			),
		},
		"static infinite": {
			input: `
				do
					println("foo")
				end until false
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 9,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 4, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 5),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(42, 4, 15), P(46, 4, 19)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"static one iteration": {
			input: `
				do
					println("foo")
				end until true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 4, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(42, 4, 15), P(45, 4, 18)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestUntil(t *testing.T) {
	tests := testTable{
		"with a body": {
			input: `
			  i := 0
				until i >= 5
					i += 1
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.LOOP), 0, 14,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(49, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 4),
				},
				[]value.Value{},
			),
		},
		"with break": {
			input: `
			  i := 0
				until false
					i += 1
					break if i < 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(23, 3, 11), P(27, 3, 15)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"with labeled break": {
			input: `
			  i := 0
				$foo: until false
					i += 1
					break[foo] if i < 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(79, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(29, 3, 17), P(33, 3, 21)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"with break with value": {
			input: `
			  i := 0
				until false
					i += 1
					break true if i < 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 7,
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 21,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 13),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(23, 3, 11), P(27, 3, 15)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"continue in a nested loop": {
			input: `
			 	j := 0
				until j >= 5
					j += 1
					i := 0
					until i >= 5
						continue if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 41,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 22,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 18,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 46,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(129, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 7),
					bytecode.NewLineInfo(7, 13),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 7),
				},
				[]value.Value{},
			),
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: until j >= 5
					j += 1
					i := 0
					until i >= 5
						continue[foo] if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 44,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 25,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 8,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 34,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 49,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(140, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 7),
					bytecode.NewLineInfo(7, 16),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 7),
				},
				[]value.Value{},
			),
		},
		"break in a nested loop": {
			input: `
			 	j := 0
				until j >= 5
					j += 1
					i := 0
					until i >= 5
						break if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 41,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 22,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 46,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(126, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 7),
					bytecode.NewLineInfo(7, 13),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 7),
				},
				[]value.Value{},
			),
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: until j >= 5
					j += 1
					i := 0
					until i >= 5
						break[foo] if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 44,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 25,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.LOAD_INT_8), 8,
					byte(bytecode.JUMP_UNLESS_IGT), 0, 8,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 15,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.INT_1),
					byte(bytecode.ADD_INT),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.LOOP), 0, 49,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(137, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 7),
					bytecode.NewLineInfo(7, 16),
					bytecode.NewLineInfo(8, 5),
					bytecode.NewLineInfo(9, 3),
					bytecode.NewLineInfo(10, 7),
				},
				[]value.Value{},
			),
		},
		"without a body": {
			input: `
				i := 0
				until i >= 5; end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_0),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.INT_5),
					byte(bytecode.JUMP_UNLESS_ILT), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 10,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(33, 3, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 12),
				},
				[]value.Value{},
			),
		},
		"static infinite": {
			input: `
				until false
					println("foo")
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 9,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(44, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 5),
				},
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
					value.Ref(value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
					)),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(11, 2, 11), P(15, 2, 15)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"static impossible": {
			input: `
				until true
					println("foo")
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(11, 2, 11), P(14, 2, 14)), "this loop will never execute since type `true` is truthy"),
				diagnostic.NewWarning(L(P(21, 3, 6), P(35, 3, 20)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestMust(t *testing.T) {
	tests := testTable{
		"with a value": {
			input: `
				var a: Int? = nil
				must a`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.MUST),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 3),
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

func TestAs(t *testing.T) {
	tests := testTable{
		"cast": {
			input: `
				var a: Int | Float = 1
				a as ::Std::Int
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.INT_1),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.AS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 3, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					value.ToSymbol("Std::Int").ToValue(),
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

func TestThrow(t *testing.T) {
	tests := testTable{
		"with a value": {
			input: `throw unchecked :foo`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.THROW),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.ToSymbol("foo").ToValue(),
				},
			),
		},
		"without a value": {
			input: `throw`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(4, 1, 5)), "thrown value of type `Std::Error` must be caught"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}

func TestCatch(t *testing.T) {
	tests := testTable{
		"simple catch": {
			input: `
				do
					throw :foo
				catch String() as str
					str
				catch :foo
					3
				end
			`,
			want: vm.NewBytecodeFunctionWithCatchEntries(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.THROW),

					byte(bytecode.JUMP), 0, 34,
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.GET_LOCAL_1),
					byte(bytecode.POP_2_SKIP_ONE),
					byte(bytecode.JUMP), 0, 12,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.INT_3),
					byte(bytecode.POP_2_SKIP_ONE),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.RETHROW),
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(88, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(4, 14),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 6),
					bytecode.NewLineInfo(7, 5),
					bytecode.NewLineInfo(8, 5),
				},
				0,
				0,
				[]value.Value{
					value.ToSymbol("foo").ToValue(),
					value.ToSymbol("Std::String").ToValue(),
				},
				[]*vm.CatchEntry{
					vm.NewCatchEntry(2, 4, 7, false),
				},
			),
		},
		"catch with stack trace": {
			input: `
				do
					throw :foo
				catch String() as str, stack_trace
					str
				catch :foo
					3
				end
			`,
			want: vm.NewBytecodeFunctionWithCatchEntries(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.THROW),
					byte(bytecode.JUMP), 0, 36,
					byte(bytecode.DUP_SECOND),
					byte(bytecode.SET_LOCAL_1),
					byte(bytecode.DUP),
					byte(bytecode.SET_LOCAL_2),
					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.GET_LOCAL_2),
					byte(bytecode.POP_2_SKIP_ONE),
					byte(bytecode.JUMP), 0, 12,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.INT_3),
					byte(bytecode.POP_2_SKIP_ONE),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.RETHROW),
					byte(bytecode.LEAVE_SCOPE16), 2, 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(101, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 6),
					bytecode.NewLineInfo(7, 5),
					bytecode.NewLineInfo(8, 5),
				},
				0,
				0,
				[]value.Value{
					value.ToSymbol("foo").ToValue(),
					value.ToSymbol("Std::String").ToValue(),
				},
				[]*vm.CatchEntry{
					vm.NewCatchEntry(2, 4, 7, false),
				},
			),
		},
		"finally": {
			input: `
				do
					println("foo")
				finally
					println("bar")
				end
			`,
			want: vm.NewBytecodeFunctionWithCatchEntries(
				mainSymbol,
				[]byte{
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE_1),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.POP),

					byte(bytecode.JUMP), 0, 40,
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.FALSE),
					byte(bytecode.JUMP), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 0,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.CALL_METHOD8), 6,
					byte(bytecode.SWAP),
					byte(bytecode.JUMP_UNLESS_UNP), 0, 2,
					byte(bytecode.POP_2),
					byte(bytecode.JUMP_TO_FINALLY),
					byte(bytecode.JUMP_IF_NP), 0, 10,
					byte(bytecode.JUMP_IF_NIL_NP), 0, 5,
					byte(bytecode.POP_2),
					byte(bytecode.POP_2_SKIP_ONE),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP_2),
					byte(bytecode.RETURN_FINALLY),
					byte(bytecode.POP_2),
					byte(bytecode.RETHROW),

					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(67, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 0),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(5, 6),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(6, 13),
					bytecode.NewLineInfo(5, 6),
					bytecode.NewLineInfo(6, 22),
				},
				0,
				0,
				[]value.Value{
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("bar"))}),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("bar"))}),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
				},
				[]*vm.CatchEntry{
					vm.NewCatchEntry(0, 5, 14, false),
					vm.NewCatchEntry(0, 5, 22, true),
				},
			),
		},
		"catch and finally": {
			input: `
				def foo! :foo
					println "foo"
					throw :foo
				end

				do
					foo()
				catch :foo
					println "bar"
				finally
					println "baz"
				end
			`,
			want: vm.NewBytecodeFunctionWithCatchEntries(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE_0),
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE_3),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.POP),
					byte(bytecode.JUMP), 0, 56,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 9,
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.CALL_METHOD8), 7,
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.FALSE),
					byte(bytecode.JUMP), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 8,
					byte(bytecode.CALL_METHOD8), 9,
					byte(bytecode.SWAP),
					byte(bytecode.JUMP_UNLESS_UNP), 0, 2,
					byte(bytecode.POP_2),
					byte(bytecode.JUMP_TO_FINALLY),
					byte(bytecode.JUMP_IF_NP), 0, 10,
					byte(bytecode.JUMP_IF_NIL_NP), 0, 5,
					byte(bytecode.POP_2),
					byte(bytecode.POP_2_SKIP_ONE),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP_2),
					byte(bytecode.RETURN_FINALLY),
					byte(bytecode.POP_2),
					byte(bytecode.RETHROW),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(153, 13, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
					bytecode.NewLineInfo(8, 4),
					bytecode.NewLineInfo(12, 6),
					bytecode.NewLineInfo(7, 3),
					bytecode.NewLineInfo(9, 7),
					bytecode.NewLineInfo(10, 9),
					bytecode.NewLineInfo(13, 13),
					bytecode.NewLineInfo(12, 6),
					bytecode.NewLineInfo(13, 22),
				},
				0,
				0,
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						methodDefinitionsSymbol,
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.LOAD_VALUE_1),
							byte(bytecode.LOAD_VALUE_2),
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(153, 13, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 7),
							bytecode.NewLineInfo(13, 2),
						},
						[]value.Value{
							value.ToSymbol("Std::Kernel").ToValue(),
							value.Ref(vm.NewBytecodeFunctionNoParams(
								value.ToSymbol("foo"),
								[]byte{
									byte(bytecode.LOAD_VALUE_0),
									byte(bytecode.CALL_SELF8), 1,
									byte(bytecode.POP),
									byte(bytecode.LOAD_VALUE_2),
									byte(bytecode.THROW),
									byte(bytecode.RETURN),
								},
								L(P(5, 2, 5), P(60, 5, 7)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(3, 4),
									bytecode.NewLineInfo(4, 2),
									bytecode.NewLineInfo(5, 1),
								},
								[]value.Value{
									value.Ref(&value.ArrayTuple{value.Ref(value.String("foo"))}),
									value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
									value.ToSymbol("foo").ToValue(),
								},
							)),
							value.ToSymbol("foo").ToValue(),
						},
					)),
					value.ToSymbol("Std::Kernel").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("foo"), 0)),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("baz"))}),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
					value.ToSymbol("foo").ToValue(),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("bar"))}),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
					value.Ref(&value.ArrayTuple{value.Ref(value.String("baz"))}),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("println"), 1)),
				},
				[]*vm.CatchEntry{
					vm.NewCatchEntry(3, 7, 16, false),
					vm.NewCatchEntry(3, 7, 40, true),
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
