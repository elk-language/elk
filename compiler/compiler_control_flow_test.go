package compiler

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestReturnExpression(t *testing.T) {
	tests := testTable{
		"return a value": {
			input: "return 5",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
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

func TestIfExpression(t *testing.T) {
	tests := testTable{
		"resolve static condition with empty then and else": {
			input: "if true; end",
			want: vm.NewBytecodeMethodNoParams(
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(1, 12),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 1),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(4, 2),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(5, 1),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(3, 3),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(6, 2),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(7, 1),
					bytecode.NewLineInfo(8, 2),
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
			want: vm.NewBytecodeMethodNoParams(
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(1, 12),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(33, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 1),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(4, 2),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(48, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(5, 1),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(3, 3),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(6, 2),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(7, 1),
					bytecode.NewLineInfo(8, 2),
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

func TestLoopExpression(t *testing.T) {
	tests := testTable{
		"empty body": {
			input: `
				loop
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOOP), 0, 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 3, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 2),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 3),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 10),
					bytecode.NewLineInfo(6, 3),
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
				a := 0
				loop
					a = a + 1
					break true if a > 5
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 10),
					bytecode.NewLineInfo(6, 3),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 5),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 8),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 5),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 8),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 6),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 10),
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
				for ;;
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 5,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 3, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(3, 4),
				},
				nil,
			),
		},
		"for without initialiser, condition and increment": {
			input: `
				a := 0
				for ;;
					a += 1
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
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
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 2),
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
				for ;;
					a += 1
					break if a > 10
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
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
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 30,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 10),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(10),
				},
			),
		},
		"for with break with value": {
			input: `
				a := 0
				for ;;
					a += 1
					break 5 if a > 10
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
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
					byte(bytecode.JUMP_UNLESS), 0, 9,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 8,
					byte(bytecode.JUMP), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LOOP), 0, 31,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(65, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 10),
					bytecode.NewLineInfo(6, 2),
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
				for a := 0;;
					a += 1
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
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
				L(P(0, 1, 1), P(37, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(2, 1),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 3),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
				},
			),
		},
		"for with initialiser, condition, without increment": {
			input: `
				for a := 0; a < 5;
					a += 1
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
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
				L(P(0, 1, 1), P(43, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 4),
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
				for i := 0; i < 5; i += 1
					a += i
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
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
				L(P(0, 1, 1), P(61, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 5),
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

func TestWhile(t *testing.T) {
	tests := testTable{
		"with a body": {
			input: `
			  i := 0
				while i < 5
					i += 1
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 6),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 3),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 10),
					bytecode.NewLineInfo(6, 3),
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
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 10),
					bytecode.NewLineInfo(6, 3),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
					value.SmallInt(5),
				},
			),
		},
		"without a body": {
			input: `
				i := 0
				while i < 5; end
			`,
			want: vm.NewBytecodeMethodNoParams(
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
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 10),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_FUNCTION8), 1,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 8,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 3),
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
			want: vm.NewBytecodeMethodNoParams(
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
