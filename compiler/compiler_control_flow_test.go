package compiler

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestForInExpression(t *testing.T) {
	tests := testTable{
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 14,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 19,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 12),
				},
				[]value.Value{
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 65,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 36,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 26,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 14,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 5,
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_IF), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.THROW),
					byte(bytecode.POP),

					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 5,
					byte(bytecode.ADD),
					byte(bytecode.CALL_SELF8), 7,

					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 5, 2,
					byte(bytecode.LOOP), 0, 70,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(2, 51),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 12),
				},
				[]value.Value{
					&value.ArrayTuple{
						&value.ArrayTuple{
							value.SmallInt(1),
							value.SmallInt(2),
						},
						&value.ArrayTuple{
							value.SmallInt(3),
							value.SmallInt(4),
						},
						&value.ArrayTuple{
							value.SmallInt(5),
							value.SmallInt(6),
						},
					},
					value.TupleMixin,
					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
					value.SmallInt(2),
					value.SmallInt(0),
					value.SmallInt(1),
					value.NewError(
						value.PatternNotMatchedErrorClass,
						"assigned value does not match the pattern defined in for in loop",
					),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 36,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 2,
					byte(bytecode.JUMP), 0, 16,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 41,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 21),
					bytecode.NewLineInfo(5, 12),
				},
				[]value.Value{
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
					value.SmallInt(2),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 37,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 12,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LEAVE_SCOPE16), 4, 2,
					byte(bytecode.JUMP), 0, 16,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 42,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(78, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 22),
					bytecode.NewLineInfo(5, 12),
				},
				[]value.Value{
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
					value.SmallInt(2),
					value.ToSymbol("foo"),
				},
			),
		},
		"with labeled break": {
			input: `
				$foo: for i in [1, 2, 3, 4, 5]
					println(i)
					break$foo if i > 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 36,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 2,
					byte(bytecode.JUMP), 0, 16,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 41,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(83, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 21),
					bytecode.NewLineInfo(5, 12),
				},
				[]value.Value{
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
					value.SmallInt(2),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 35,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 23,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.CALL_SELF8), 2,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 40,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(3, 21),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 12),
				},
				[]value.Value{
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
					},
					value.SmallInt(2),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
				},
			),
		},
		"with labeled continue": {
			input: `
				$foo: for i in [1, 2, 3, 4, 5]
					continue$foo if i > 2
					println(i)
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 35,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 23,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.CALL_SELF8), 2,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 40,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(86, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(3, 21),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 12),
				},
				[]value.Value{
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
					},
					value.SmallInt(2),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 64,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),

					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 5,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 5,
					byte(bytecode.FOR_IN), 0, 38,
					byte(bytecode.SET_LOCAL8), 6,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 6,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 6, 2,
					byte(bytecode.JUMP), 0, 23,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 6,
					byte(bytecode.CALL_SELF8), 3,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 6, 1,
					byte(bytecode.LOOP), 0, 43,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 5, 1,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 69,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(122, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(7, 1),
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(4, 22),
					bytecode.NewLineInfo(5, 6),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					&value.ArrayList{
						value.String("a"),
						value.String("b"),
						value.String("c"),
						value.String("d"),
					},
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
					},
					value.SmallInt(2),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						2,
						nil,
					),
				},
			),
		},
		"nested with labeled break": {
			input: `
				$foo: for c in ['a', 'b', 'c', 'd']
					for i in [1, 2, 3, 4, 5]
						break$foo if i > 2
						println(c, i)
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 4,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 64,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),

					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 5,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 5,
					byte(bytecode.FOR_IN), 0, 38,
					byte(bytecode.SET_LOCAL8), 6,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 6,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 6, 4,
					byte(bytecode.JUMP), 0, 34,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 6,
					byte(bytecode.CALL_SELF8), 3,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 6, 1,
					byte(bytecode.LOOP), 0, 43,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 5, 1,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 69,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(132, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(7, 1),
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(4, 22),
					bytecode.NewLineInfo(5, 6),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					&value.ArrayList{
						value.String("a"),
						value.String("b"),
						value.String("c"),
						value.String("d"),
					},
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
					},
					value.SmallInt(2),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						2,
						nil,
					),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 63,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),

					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 5,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 5,
					byte(bytecode.FOR_IN), 0, 37,
					byte(bytecode.SET_LOCAL8), 6,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 6,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 6, 1,
					byte(bytecode.LOOP), 0, 23,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 6,
					byte(bytecode.CALL_SELF8), 3,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 6, 1,
					byte(bytecode.LOOP), 0, 42,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 5, 1,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 68,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(125, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(7, 1),
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(4, 21),
					bytecode.NewLineInfo(5, 6),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					&value.ArrayList{
						value.String("a"),
						value.String("b"),
						value.String("c"),
						value.String("d"),
					},
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
					},
					value.SmallInt(2),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						2,
						nil,
					),
				},
			),
		},
		"nested with labeled continue": {
			input: `
				$foo: for c in ['a', 'b', 'c', 'd']
					for i in [1, 2, 3, 4, 5]
						continue$foo if i > 2
						println(c, i)
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 4,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 63,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),

					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 5,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 5,
					byte(bytecode.FOR_IN), 0, 37,
					byte(bytecode.SET_LOCAL8), 6,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 6,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 6, 3,
					byte(bytecode.LOOP), 0, 38,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 6,
					byte(bytecode.CALL_SELF8), 3,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 6, 1,
					byte(bytecode.LOOP), 0, 42,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 5, 1,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 68,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(135, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(7, 1),
					bytecode.NewLineInfo(2, 8),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(4, 21),
					bytecode.NewLineInfo(5, 6),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 12),
				},
				[]value.Value{
					&value.ArrayList{
						value.String("a"),
						value.String("b"),
						value.String("c"),
						value.String("d"),
					},
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
						value.SmallInt(4),
						value.SmallInt(5),
					},
					value.SmallInt(2),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						2,
						nil,
					),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
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

func TestModifierForIn(t *testing.T) {
	tests := testTable{
		"iterate": {
			input: `println(i) for i in [1, 2, 3]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 14,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 19,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 1, 29)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 33),
				},
				[]value.Value{
					&value.ArrayList{
						value.SmallInt(1),
						value.SmallInt(2),
						value.SmallInt(3),
					},
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
				},
			),
		},
		"with a pattern": {
			input: `println(a + b) for %[a, b] in %[%[1, 2], %[3, 4], %[5, 6]]`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 3,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.GET_ITERATOR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.FOR_IN), 0, 65,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 36,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 26,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 14,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 5,
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_IF), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.THROW),
					byte(bytecode.POP),

					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 5,
					byte(bytecode.ADD),
					byte(bytecode.CALL_SELF8), 7,

					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 5, 2,
					byte(bytecode.LOOP), 0, 70,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(57, 1, 58)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 83),
				},
				[]value.Value{
					&value.ArrayTuple{
						&value.ArrayTuple{
							value.SmallInt(1),
							value.SmallInt(2),
						},
						&value.ArrayTuple{
							value.SmallInt(3),
							value.SmallInt(4),
						},
						&value.ArrayTuple{
							value.SmallInt(5),
							value.SmallInt(6),
						},
					},
					value.TupleMixin,
					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
					value.SmallInt(2),
					value.SmallInt(0),
					value.SmallInt(1),
					value.NewError(
						value.PatternNotMatchedErrorClass,
						"assigned value does not match the pattern defined in for in loop",
					),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
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
			input: "if true; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"empty then and else": {
			input: "a := true; if a; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.TRUE),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.JUMP_UNLESS), 0, 5,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 2,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 19),
				},
				nil,
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.SmallInt(10),
				},
			),
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
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(4, 1),
				},
				nil,
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{
					value.SmallInt(5),
				},
			),
		},
		"with then branch": {
			input: `
				a := 5
				if a
					a = a * 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.JUMP_UNLESS), 0, 11,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.JUMP), 0, 2,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(4, 7),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
				},
			),
		},
		"with then and else branches": {
			input: `
				a := 5
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.JUMP_UNLESS), 0, 11,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.JUMP), 0, 5,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(4, 7),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
					value.SmallInt(30),
				},
			),
		},
		"is an expression": {
			input: `
			a := 5
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.JUMP_UNLESS), 0, 6,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.JUMP), 0, 3,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					// end
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(59, 8, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(6, 2),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(7, 1),
					bytecode.NewLineInfo(8, 3),
				},
				[]value.Value{
					value.SmallInt(5),
					value.String("foo"),
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
				nil,
			),
		},
		"empty then and else": {
			input: "a := true; unless a; end",
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.TRUE),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.JUMP_IF), 0, 5,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 2,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 1, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 19),
				},
				nil,
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(33, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.SmallInt(10),
				},
			),
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
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(4, 1),
				},
				nil,
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(48, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{
					value.SmallInt(5),
				},
			),
		},
		"with then branch": {
			input: `
				a := 5
				unless a
					a = a * 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.JUMP_IF), 0, 11,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.JUMP), 0, 2,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(4, 7),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
				},
			),
		},
		"with then and else branches": {
			input: `
				a := 5
				unless a
					a = a * 2
				else
					a = 30
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.JUMP_IF), 0, 11,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.JUMP), 0, 5,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(4, 7),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
					value.SmallInt(30),
				},
			),
		},
		"is an expression": {
			input: `
			a := 5
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.JUMP_IF), 0, 6,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.JUMP), 0, 3,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					// end
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 8, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(6, 2),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(7, 1),
					bytecode.NewLineInfo(8, 3),
				},
				[]value.Value{
					value.SmallInt(5),
					value.String("foo"),
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
			err: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(4, 1, 5)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"in top level with a label": {
			input: "break$foo",
			err: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(8, 1, 9)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"nonexistent label": {
			input: `
				loop
					break$foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L(P(15, 3, 6), P(23, 3, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to an expression": {
			input: `
				loop
					$foo: 1 + 2
					break$foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L(P(32, 4, 6), P(40, 4, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to a different loop": {
			input: `
				$foo: loop
					println("foo")
				end

				loop
					break$foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L(P(59, 7, 6), P(67, 7, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
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
			err: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(7, 1, 8)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"in top level with a label": {
			input: "continue$foo",
			err: error.ErrorList{
				error.NewFailure(L(P(0, 1, 1), P(11, 1, 12)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"nonexistent label": {
			input: `
				loop
					continue$foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L(P(15, 3, 6), P(26, 3, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to an expression": {
			input: `
				loop
					$foo: 1 + 2
					continue$foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L(P(32, 4, 6), P(43, 4, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to a different loop": {
			input: `
				$foo: loop
					println("foo")
				end

				loop
					continue$foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L(P(59, 7, 6), P(70, 7, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
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
					bytecode.NewLineInfo(3, 4),
				},
				nil,
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					// loop body
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 11,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 7),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					// loop body
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 11, // continue
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.CALL_SELF8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(77, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(7, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.String("foo"),
					value.NewCallSiteInfo(value.ToSymbol("println"), 1, nil),
				},
			),
		},
		"with labeled continue": {
			input: `
				a := 0
				$foo: loop
					a = a + 1
					continue$foo
					println("foo")
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					// loop body
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 11, // continue
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.CALL_SELF8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(87, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(6, 4),
					bytecode.NewLineInfo(7, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.String("foo"),
					value.NewCallSiteInfo(value.ToSymbol("println"), 1, nil),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 12,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 29,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 58,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 70,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(134, 11, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 18),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(10, 20),
					bytecode.NewLineInfo(11, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: loop
					j += 1
					i := 0
					loop
						continue$foo if i >= 5
						i += 1
					end
					continue if j >= 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 28,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 32,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 61,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 73,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(144, 11, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 21),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(10, 20),
					bytecode.NewLineInfo(11, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					// loop body
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
		},
		"with a labeled break": {
			input: `
				a := 0
				$foo: loop
					a = a + 1
					break$foo if a > 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					// loop body
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 17,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.JUMP), 0, 12,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 72,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(128, 11, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 19),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(10, 21),
					bytecode.NewLineInfo(11, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$outer: loop
					j += 1
					i := 0
					loop
						break$outer if i >= 5
						i += 1
					end
					break if j >= 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.JUMP), 0, 46,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 33,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.JUMP), 0, 12,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 75,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(142, 11, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 22),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(10, 21),
					bytecode.NewLineInfo(11, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					// loop body
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
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

func TestLogicalOrOperator(t *testing.T) {
	tests := testTable{
		"simple": {
			input: `
				"foo" || true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.JUMP_IF), 0, 2,
					// falsy
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// truthy
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 2, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
				},
				[]value.Value{
					value.String("foo"),
				},
			),
		},
		"nested": {
			input: `
				"foo" || true || 3
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.JUMP_IF), 0, 2,
					// falsy 1
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// truthy 1
					byte(bytecode.JUMP_IF), 0, 3,
					// falsy 2
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					// truthy 2
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 2, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 14),
				},
				[]value.Value{
					value.String("foo"),
					value.SmallInt(3),
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

func TestLogicalAndOperator(t *testing.T) {
	tests := testTable{
		"simple": {
			input: `
				"foo" && true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.JUMP_UNLESS), 0, 2,
					// truthy
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// falsy
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 2, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
				},
				[]value.Value{
					value.String("foo"),
				},
			),
		},
		"nested": {
			input: `
				"foo" && true && 3
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.JUMP_UNLESS), 0, 2,
					// truthy 1
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// falsy 1
					byte(bytecode.JUMP_UNLESS), 0, 3,
					// truthy 2
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					// falsy 2
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 2, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 14),
				},
				[]value.Value{
					value.String("foo"),
					value.SmallInt(3),
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

func TestNilCoalescingOperator(t *testing.T) {
	tests := testTable{
		"simple": {
			input: `
				"foo" ?? true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.JUMP_IF_NIL), 0, 3,
					byte(bytecode.JUMP), 0, 2,
					// nil
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// not nil
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(18, 2, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 11),
				},
				[]value.Value{
					value.String("foo"),
				},
			),
		},
		"nested": {
			input: `
				"foo" ?? true ?? 3
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.JUMP_IF_NIL), 0, 3,
					byte(bytecode.JUMP), 0, 2,
					// nil 1
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// not nil 1
					byte(bytecode.JUMP_IF_NIL), 0, 3,
					byte(bytecode.JUMP), 0, 3,
					// nil 2
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					// not nil 2
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(23, 2, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 20),
				},
				[]value.Value{
					value.String("foo"),
					value.SmallInt(3),
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
					bytecode.NewLineInfo(3, 4),
				},
				nil,
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 11,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 7),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(66, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(10),
				},
			),
		},
		"for with labeled break": {
			input: `
				a := 0
				$foo: fornum ;;
					a += 1
					break$foo if a > 10
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(10),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 5,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.LOOP), 0, 35,
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(84, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 18),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(6, 7),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(10),
					value.SmallInt(1),
				},
			),
		},
		"nested for with a labeled continue": {
			input: `
				$foo: fornum a := 0;;
					fornum ;; a += 1
						continue$foo if a > 10
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 16,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.LOOP), 0, 35,
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(94, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 18),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(6, 7),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(10),
					value.SmallInt(1),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 16,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.LOOP), 0, 35,
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(81, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 18),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(6, 7),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(10),
					value.SmallInt(1),
				},
			),
		},
		"nested for with a labeled break": {
			input: `
				$foo: fornum a := 0;;
					fornum ;; a += 1
						break$foo if a > 10
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.JUMP), 0, 22,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 33,
					byte(bytecode.LOOP), 0, 38,
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(91, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 21),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(5, 4),
					bytecode.NewLineInfo(6, 7),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(10),
					value.SmallInt(1),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 12,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.JUMP), 0, 11,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 34,
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(3, 8),
					bytecode.NewLineInfo(4, 22),
					bytecode.NewLineInfo(5, 7),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(10),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.LOOP), 0, 11,
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 7),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 12,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(2, 9),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 20,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 28,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(3, 9),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(4, 7),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(5, 9),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 22),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(51, 5, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 7),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(70, 6, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
		},
		"with labeled break": {
			input: `
			  i := 0
				$foo: do
					i += 1
					break$foo if i < 5
				end while true
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(80, 6, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(75, 6, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 13,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 42,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 72,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(133, 10, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 22),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(10, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(8),
					value.SmallInt(5),
				},
			),
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						continue$foo if i + j > 8
						i += 1
					end while i < 5
				end while j < 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.JUMP), 0, 30,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 75,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(143, 10, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 25),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(10, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(8),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 27,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 42,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 72,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(130, 10, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 22),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(10, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(8),
					value.SmallInt(5),
				},
			),
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						break$foo if i + j > 8
						i += 1
					end while i < 5
				end while j < 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.JUMP), 0, 44,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 75,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(140, 10, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 25),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(10, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(8),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 8,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 4, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 5),
				},
				[]value.Value{
					value.String("foo"),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 4, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.String("foo"),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 12,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(48, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 11),
					bytecode.NewLineInfo(4, 7),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(67, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
		},
		"with labeled break": {
			input: `
			  i := 0
				$foo: while true
					i += 1
					break$foo if i < 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(77, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(72, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 65,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 34,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 26,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.LOOP), 0, 42,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 73,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(127, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 11),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 22),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
					value.SmallInt(8),
				},
			),
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: while j < 5
					j += 1
					i := 0
					while i < 5
						continue$foo if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 68,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 37,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 53,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 76,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(137, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 11),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 25),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
					value.SmallInt(8),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 65,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 34,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 17,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.LOOP), 0, 42,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 73,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(124, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 11),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 22),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
					value.SmallInt(8),
				},
			),
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: while j < 5
					j += 1
					i := 0
					while i < 5
						break$foo if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 68,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 37,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.JUMP), 0, 24,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 76,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(134, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 11),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 25),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
					value.SmallInt(8),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.POP_N), 2,
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 14,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 3, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 17),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 8,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 5),
				},
				[]value.Value{
					value.String("foo"),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
				},
			),
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
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(4, 1),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(36, 3, 24)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 22),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(52, 5, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 7),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 6, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
		},
		"with labeled break": {
			input: `
			  i := 0
				$foo: do
					i += 1
					break$foo if i < 5
				end until false
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(81, 6, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 6, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 13,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 42,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 72,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(135, 10, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 22),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(10, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(8),
					value.SmallInt(5),
				},
			),
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						continue$foo if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.JUMP), 0, 30,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 75,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(145, 10, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 25),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(10, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(8),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 27,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 42,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 72,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(132, 10, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 22),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(10, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(8),
					value.SmallInt(5),
				},
			),
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						break$foo if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.JUMP), 0, 44,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.POP_N), 2,
					byte(bytecode.LOOP), 0, 75,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(142, 10, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(7, 25),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 5),
					bytecode.NewLineInfo(6, 5),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(10, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(8),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 8,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 4, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 5),
				},
				[]value.Value{
					value.String("foo"),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 4, 19)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					value.String("foo"),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 12,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(49, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 11),
					bytecode.NewLineInfo(4, 7),
					bytecode.NewLineInfo(5, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
		},
		"with labeled break": {
			input: `
			  i := 0
				$foo: until false
					i += 1
					break$foo if i < 5
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(78, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 9,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 18),
					bytecode.NewLineInfo(6, 5),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 65,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 34,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 26,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.LOOP), 0, 42,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 73,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(129, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 11),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 22),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
					value.SmallInt(8),
				},
			),
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: until j >= 5
					j += 1
					i := 0
					until i >= 5
						continue$foo if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 68,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 37,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 53,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 76,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(139, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 11),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 25),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
					value.SmallInt(8),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 65,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 34,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 17,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.LOOP), 0, 42,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 73,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(126, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 11),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 22),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
					value.SmallInt(8),
				},
			),
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: until j >= 5
					j += 1
					i := 0
					until i >= 5
						break$foo if i + j > 8
						i += 1
					end
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 68,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 37,
					byte(bytecode.POP_N), 2,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.GREATER),
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.JUMP), 0, 24,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.LOOP), 0, 45,
					byte(bytecode.POP),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.LOOP), 0, 76,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(136, 10, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 11),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 11),
					bytecode.NewLineInfo(7, 25),
					bytecode.NewLineInfo(8, 7),
					bytecode.NewLineInfo(9, 4),
					bytecode.NewLineInfo(10, 8),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
					value.SmallInt(1),
					value.SmallInt(8),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_IF), 0, 6,
					byte(bytecode.POP_N), 2,
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 14,
					byte(bytecode.POP),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(33, 3, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 17),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(5),
				},
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
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 8,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(44, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 5),
				},
				[]value.Value{
					value.String("foo"),
					value.NewCallSiteInfo(
						value.ToSymbol("println"),
						1,
						nil,
					),
				},
			),
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
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(4, 1),
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

func TestThrow(t *testing.T) {
	tests := testTable{
		"with a value": {
			input: `throw :foo`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.THROW),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.ToSymbol("foo"),
				},
			),
		},
		"without a value": {
			input: `throw`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.THROW),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.NewError(value.ErrorClass, "error"),
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

func TestCatch(t *testing.T) {
	tests := testTable{
		"simple catch": {
			input: `
				do
					throw :foo
				catch ::String() as str
					str
				catch :foo
					3
				end
			`,
			want: vm.NewBytecodeFunctionWithCatchEntries(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.THROW),

					byte(bytecode.JUMP), 0, 41,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.DUP),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.POP_N_SKIP_ONE), 2,
					byte(bytecode.JUMP), 0, 18,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.POP_N_SKIP_ONE), 2,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.RETHROW),

					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(90, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(4, 16),
					bytecode.NewLineInfo(5, 8),
					bytecode.NewLineInfo(6, 8),
					bytecode.NewLineInfo(7, 8),
					bytecode.NewLineInfo(8, 2),
				},
				nil,
				0,
				-1,
				false,
				false,
				[]value.Value{
					value.ToSymbol("foo"),
					value.ToSymbol("String"),
					value.SmallInt(3),
				},
				[]*vm.CatchEntry{
					vm.NewCatchEntry(2, 5, 8, false),
				},
			),
		},
		"finally": {
			input: `
				do
					foo()
				finally
					bar()
				end
			`,
			want: vm.NewBytecodeFunctionWithCatchEntries(
				mainSymbol,
				[]byte{
					byte(bytecode.CALL_SELF8), 0,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.POP),

					byte(bytecode.JUMP), 0, 41,
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.FALSE),
					byte(bytecode.JUMP), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.CALL_SELF8), 2,
					byte(bytecode.SWAP),
					byte(bytecode.JUMP_UNLESS_UNDEF), 0, 3,
					byte(bytecode.POP_N), 2,
					byte(bytecode.JUMP_TO_FINALLY),
					byte(bytecode.JUMP_IF), 0, 13,
					byte(bytecode.JUMP_IF_NIL), 0, 7,
					byte(bytecode.POP_N), 2,
					byte(bytecode.POP_N_SKIP_ONE), 2,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.POP_N), 2,
					byte(bytecode.RETURN_FINALLY),
					byte(bytecode.POP_N), 2,
					byte(bytecode.RETHROW),

					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(49, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(6, 13),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 27),
				},
				nil,
				0,
				-1,
				false,
				false,
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.NewCallSiteInfo(value.ToSymbol("bar"), 0, nil),
					value.NewCallSiteInfo(value.ToSymbol("bar"), 0, nil),
				},
				[]*vm.CatchEntry{
					vm.NewCatchEntry(0, 2, 8, false),
					vm.NewCatchEntry(0, 2, 16, true),
				},
			),
		},
		"catch and finally": {
			input: `
				do
					foo()
				catch :foo
					bar()
				finally
					baz()
				end
			`,
			want: vm.NewBytecodeFunctionWithCatchEntries(
				mainSymbol,
				[]byte{
					byte(bytecode.CALL_SELF8), 0,
					byte(bytecode.CALL_SELF8), 1,
					byte(bytecode.POP),

					byte(bytecode.JUMP), 0, 55,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.POP),
					byte(bytecode.CALL_SELF8), 3,
					byte(bytecode.JUMP), 0, 5,

					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.FALSE),
					byte(bytecode.JUMP), 0, 5,
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.CALL_SELF8), 4,
					byte(bytecode.SWAP),
					byte(bytecode.JUMP_UNLESS_UNDEF), 0, 3,
					byte(bytecode.POP_N), 2,
					byte(bytecode.JUMP_TO_FINALLY),
					byte(bytecode.JUMP_IF), 0, 13,
					byte(bytecode.JUMP_IF_NIL), 0, 7,
					byte(bytecode.POP_N), 2,
					byte(bytecode.POP_N_SKIP_ONE), 2,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.POP_N), 2,
					byte(bytecode.RETURN_FINALLY),
					byte(bytecode.POP_N), 2,
					byte(bytecode.RETHROW),

					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(75, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(7, 2),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 8),
					bytecode.NewLineInfo(5, 6),
					bytecode.NewLineInfo(8, 13),
					bytecode.NewLineInfo(7, 2),
					bytecode.NewLineInfo(8, 27),
				},
				nil,
				0,
				-1,
				false,
				false,
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.NewCallSiteInfo(value.ToSymbol("baz"), 0, nil),
					value.ToSymbol("foo"),
					value.NewCallSiteInfo(value.ToSymbol("bar"), 0, nil),
					value.NewCallSiteInfo(value.ToSymbol("baz"), 0, nil),
				},
				[]*vm.CatchEntry{
					vm.NewCatchEntry(0, 2, 8, false),
					vm.NewCatchEntry(0, 2, 30, true),
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
