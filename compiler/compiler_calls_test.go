package compiler

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestSubscript(t *testing.T) {
	tests := testTable{
		"static": {
			input: "[5, 3][0]",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.SmallInt(5),
				},
			),
		},
		"dynamic": {
			input: `
				arr := [5, 3]
				arr[1]
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(29, 3, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{
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
				arr[1] = :foo
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT_SET),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(36, 3, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					&value.ArrayList{
						value.SmallInt(5),
						value.SmallInt(3),
					},
					value.SmallInt(1),
					value.ToSymbol("foo"),
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
			input: "::Foo()",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.INSTANTIATE8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
					value.NewCallSiteInfo(value.ToSymbol("#init"), 0, nil),
				},
			),
		},
		"complex constant": {
			input: "::Foo::Bar::Baz()",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.GET_MOD_CONST8), 2,
					byte(bytecode.INSTANTIATE8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
					value.ToSymbol("Bar"),
					value.ToSymbol("Baz"),
					value.NewCallSiteInfo(value.ToSymbol("#init"), 0, nil),
				},
			),
		},
		"with positional arguments": {
			input: "::Foo(1, 'lol')",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.INSTANTIATE8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("#init"), 2, nil),
				},
			),
		},
		"with named args": {
			input: `::Foo(1, b: 'lol')`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.INSTANTIATE8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.ToSymbol("Foo"),
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("#init"), 2, []value.Symbol{value.ToSymbol("b")}),
				},
			),
		},
		"with duplicated named args": {
			input: "::Foo(b: 1, a: 'lol', b: 2)",
			err: errors.ErrorList{
				errors.NewError(
					L(P(22, 1, 23), P(25, 1, 26)),
					"duplicated named argument in call: :b",
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
			input: "self.foo",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
				},
			),
		},
		"call a setter": {
			input: "self.foo = 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a method with positional arguments": {
			input: "self.foo(1, 'lol')",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 2, nil),
				},
			),
		},
		"call a method on a local variable": {
			input: `
				a := 25
				a.foo(1, 'lol')
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
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 3, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					value.SmallInt(25),
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 2, nil),
				},
			),
		},
		"call a method on a local variable with named args": {
			input: `
				a := 25
				a.foo(1, b: 'lol')
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
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(35, 3, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 5),
				},
				[]value.Value{
					value.SmallInt(25),
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 2, []value.Symbol{value.ToSymbol("b")}),
				},
			),
		},
		"call a method with duplicated named args": {
			input: "self.foo(b: 1, a: 'lol', b: 2)",
			err: errors.ErrorList{
				errors.NewError(
					L(P(25, 1, 26), P(28, 1, 29)),
					"duplicated named argument in call: :b",
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
		"call a function without arguments": {
			input: "foo()",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.CALL_FUNCTION8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
				},
			),
		},
		"call a function with positional arguments": {
			input: "foo(1, 'lol')",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.CALL_FUNCTION8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 2, nil),
				},
			),
		},
		"call a function with named args": {
			input: "foo(1, b: 'lol')",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.CALL_FUNCTION8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.SmallInt(1),
					value.String("lol"),
					value.NewCallSiteInfo(value.ToSymbol("foo"), 2, []value.Symbol{value.ToSymbol("b")}),
				},
			),
		},
		"call a function with duplicated named args": {
			input: "foo(b: 1, a: 'lol', b: 2)",
			err: errors.ErrorList{
				errors.NewError(
					L(P(20, 1, 21), P(23, 1, 24)),
					"duplicated named argument in call: :b",
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
			input: "self.foo = 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CALL_METHOD8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with add": {
			input: "self.foo += 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with subtract": {
			input: "self.foo -= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with multiply": {
			input: "self.foo *= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with divide": {
			input: "self.foo /= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DIVIDE),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with exponentiate": {
			input: "self.foo **= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with modulo": {
			input: "self.foo %= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.MODULO),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with left bitshift": {
			input: "self.foo <<= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with logic left bitshift": {
			input: "self.foo <<<= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with right bitshift": {
			input: "self.foo >>= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with logic right bitshift": {
			input: "self.foo >>>= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with bitwise and": {
			input: "self.foo &= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with bitwise or": {
			input: "self.foo |= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with bitwise xor": {
			input: "self.foo ^= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(12, 1, 13)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with logic or": {
			input: "self.foo ||= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.JUMP_IF), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with logic and": {
			input: "self.foo &&= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
				},
			),
		},
		"call a setter with nil coalesce": {
			input: "self.foo ??= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.SELF),
					byte(bytecode.CALL_METHOD8), 0,
					byte(bytecode.JUMP_IF_NIL), 0, 3,
					byte(bytecode.JUMP), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.CALL_METHOD8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil),
					value.SmallInt(3),
					value.NewCallSiteInfo(value.ToSymbol("foo="), 1, nil),
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
