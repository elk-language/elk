package compiler

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

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