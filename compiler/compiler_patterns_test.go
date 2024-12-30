package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func TestSwitch(t *testing.T) {
	tests := testTable{
		"with a few literal cases": {
			input: `
			  var a: any = 0
				switch a
				case true then "a"
				case false then "b"
				case 0 then "c"
				case 1 then "d"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.TRUE),
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.JUMP), 0, 47,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.FALSE),
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 33,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 18,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(128, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 14),
					bytecode.NewLineInfo(5, 14),
					bytecode.NewLineInfo(6, 15),
					bytecode.NewLineInfo(7, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(8, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("a")),
					value.Ref(value.String("b")),
					value.Ref(value.String("c")),
					value.SmallInt(1).ToValue(),
					value.Ref(value.String("d")),
				},
			),
		},
		"with else": {
			input: `
			  var a: any = 0
				switch a
				case true then "a"
				case false then "b"
				else "c"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.TRUE),
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.JUMP), 0, 18,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.FALSE),
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(101, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 14),
					bytecode.NewLineInfo(5, 14),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(6, 2),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("a")),
					value.Ref(value.String("b")),
					value.Ref(value.String("c")),
				},
			),
		},
		"literal true": {
			input: `
			  var a: any = 0
				switch a
				case true then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.TRUE),
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 14),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal false": {
			input: `
			  var a: any = 0
				switch a
				case false then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.FALSE),
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(65, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 14),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},

		"literal nil": {
			input: `
			  var a: any = 0
				switch a
				case nil then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.NIL),
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 14),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal string": {
			input: `
			  var a: any = 0
				switch a
				case "foo" then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(65, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("foo")),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal raw string": {
			input: `
			  var a: any = 0
				switch a
				case 'foo' then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(65, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("foo")),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal interpolated string": {
			input: `
			  var a: any = 0
				switch a
				case "f${a}" then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_STRING8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(67, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 19),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("f")),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal symbol": {
			input: `
			  var a: any = 0
				switch a
				case :foo then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.ToSymbol("foo").ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal interpolated symbol": {
			input: `
			  var a: any = 0
				switch a
				case :"f${a}" then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.NEW_SYMBOL8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(68, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 19),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("f")),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal int": {
			input: `
			  var a: any = 0
				switch a
				case 5 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal int64": {
			input: `
			  var a: any = 0
				switch a
				case 5i64 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Int64(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal uint64": {
			input: `
			  var a: any = 0
				switch a
				case 5u64 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.UInt64(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal int32": {
			input: `
			  var a: any = 0
				switch a
				case 5i32 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Int32(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal uint32": {
			input: `
			  var a: any = 0
				switch a
				case 5u32 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.UInt32(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal int16": {
			input: `
			  var a: any = 0
				switch a
				case 5i16 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Int16(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal uint16": {
			input: `
			  var a: any = 0
				switch a
				case 5u16 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.UInt16(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal int8": {
			input: `
			  var a: any = 0
				switch a
				case 5i8 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Int8(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal uint8": {
			input: `
			  var a: any = 0
				switch a
				case 5u8 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.UInt8(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal float": {
			input: `
			  var a: any = 0
				switch a
				case 5.8 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Float(5.8).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal float64": {
			input: `
			  var a: any = 0
				switch a
				case 5.8f64 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(66, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Float64(5.8).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal float32": {
			input: `
			  var a: any = 0
				switch a
				case 5.8f32 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(66, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Float32(5.8).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal negative float32": {
			input: `
			  var a: any = 0
				switch a
				case -5.8f32 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(67, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Float32(-5.8).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"literal big float": {
			input: `
			  var a: any = 0
				switch a
				case 5.8bf then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(65, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.NewBigFloat(5.8)),
					value.Ref(value.String("a")),
				},
			),
		},
		"root constant": {
			input: `
				const Foo = 3
			  var a: any = 0
				switch a
				case ::Foo then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.DEF_CONST),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(83, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(5, 15),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Root").ToValue(),
					value.ToSymbol("Foo").ToValue(),
					value.SmallInt(3).ToValue(),
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"constant lookup": {
			input: `
				module Foo
					const Bar = 3
				end
			  var a: any = 0
				switch a
				case ::Foo::Bar then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.DEF_CONST),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 6,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(112, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 2),
					bytecode.NewLineInfo(7, 15),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(8, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(112, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Undefined,
					value.ToSymbol("Foo").ToValue(),
					value.ToSymbol("Bar").ToValue(),
					value.SmallInt(3).ToValue(),
					value.SmallInt(0).ToValue(),
					value.ToSymbol("Foo::Bar").ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"negative constant lookup": {
			input: `
				module Foo
					const Bar = 3
				end
			  var a: any = 0
				switch a
				case -::Foo::Bar then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.EXEC),
					byte(bytecode.POP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.DEF_CONST),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 6,
					byte(bytecode.NEGATE),
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(113, 8, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(5, 5),
					bytecode.NewLineInfo(6, 2),
					bytecode.NewLineInfo(7, 16),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(8, 2),
				},
				[]value.Value{
					value.Ref(vm.NewBytecodeFunctionNoParams(
						value.ToSymbol("<namespaceDefinitions>"),
						[]byte{
							byte(bytecode.GET_CONST8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_NAMESPACE), 0,
							byte(bytecode.NIL),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(113, 8, 8)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
							bytecode.NewLineInfo(8, 2),
						},
						[]value.Value{
							value.ToSymbol("Root").ToValue(),
							value.ToSymbol("Foo").ToValue(),
						},
					)),
					value.Undefined,
					value.ToSymbol("Foo").ToValue(),
					value.ToSymbol("Bar").ToValue(),
					value.SmallInt(3).ToValue(),
					value.SmallInt(0).ToValue(),
					value.ToSymbol("Foo::Bar").ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"less pattern": {
			input: `
			  var a: any = 0
				switch a
				case < 5 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 31),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"less than root constant": {
			input: `
				const Foo = 5
			  var a: any = 0
				switch a
				case < ::Foo then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.DEF_CONST),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(5, 31),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Root").ToValue(),
					value.ToSymbol("Foo").ToValue(),
					value.SmallInt(5).ToValue(),
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"less than negative root constant": {
			input: `
				const Foo = 2
			  var a: any = 0
				switch a
				case < -::Foo then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.DEF_CONST),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.NEGATE),
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(86, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(5, 32),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Root").ToValue(),
					value.ToSymbol("Foo").ToValue(),
					value.SmallInt(2).ToValue(),
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"less equal pattern": {
			input: `
			  var a: any = 0
				switch a
				case <= 5 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS_EQUAL),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 31),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"greater pattern": {
			input: `
			  var a: any = 0
				switch a
				case > 5 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.GREATER),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(63, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 31),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"greater equal pattern": {
			input: `
			  var a: any = 0
				switch a
				case >= 5 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 31),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"equal pattern": {
			input: `
			  var a: any = 0
				switch a
				case == 5 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"equal regex pattern": {
			input: `
			  var a: any = 0
				switch a
				case == %/fo+/ then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(69, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.MustCompileRegex("fo+", bitfield.BitField8{})),
					value.Ref(value.String("a")),
				},
			),
		},
		"equal local pattern": {
			input: `
			  var a: any = 0
				b := 2
				switch a
				case == b then "a"
				end
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
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(75, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(5, 15),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(2).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"not equal pattern": {
			input: `
			  var a: any = 0
				switch a
				case != 5 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.NOT_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"lax equal pattern": {
			input: `
			  var a: any = 0
				b := 2
				switch a
				case =~ b then "a"
				end
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
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.LAX_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(75, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 5),
					bytecode.NewLineInfo(4, 2),
					bytecode.NewLineInfo(5, 15),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(6, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(2).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"lax not equal pattern": {
			input: `
			  var a: any = 0
				switch a
				case !~ 5 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LAX_NOT_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(64, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"strict equal pattern": {
			input: `
			  var a: any = 0
				switch a
				case === 5 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.STRICT_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(65, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"strict not equal pattern": {
			input: `
			  var a: any = 0
				switch a
				case !== 5 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.STRICT_NOT_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(65, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"strict not equal negative pattern": {
			input: `
			  var a: any = 0
				switch a
				case !== -5 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.STRICT_NOT_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(66, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 15),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(-5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},

		"regex pattern": {
			input: `
			  a := "foo"
				switch a
				case %/fo+/ then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SWAP),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(62, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 17),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(value.String("foo")),
					value.Ref(value.MustCompileRegex("fo+", bitfield.BitField8{})),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("matches"), 1)),
					value.Ref(value.String("a")),
				},
			),
		},
		"variable pattern": {
			input: `
			  a := 0
				switch a
				case n then n + 2
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 13,
					byte(bytecode.POP_2),
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(55, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 23),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.SmallInt(2).ToValue(),
				},
			),
		},
		"range": {
			input: `
			  var a: any = 0
				switch a
				case -2...9 then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SWAP),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(66, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 17),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.SmallInt(0).ToValue(),
					value.Ref(value.NewClosedRange(value.SmallInt(-2).ToValue(), value.SmallInt(9).ToValue())),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#contains"), 1)),
					value.Ref(value.String("a")),
				},
			),
		},
		"range with constants": {
			input: `
				const Foo = 3
				const Bar = 10
			  var a: any = 0
				switch a
				case ::Foo...-::Bar then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.DEF_CONST),
					byte(bytecode.GET_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.DEF_CONST),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.GET_CONST8), 4,
					byte(bytecode.NEGATE),
					byte(bytecode.NEW_RANGE), 0,
					byte(bytecode.SWAP),
					byte(bytecode.CALL_METHOD8), 7,
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 8,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(111, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 7),
					bytecode.NewLineInfo(3, 7),
					bytecode.NewLineInfo(4, 5),
					bytecode.NewLineInfo(5, 2),
					bytecode.NewLineInfo(6, 22),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(7, 2),
				},
				[]value.Value{
					value.Undefined,
					value.ToSymbol("Root").ToValue(),
					value.ToSymbol("Foo").ToValue(),
					value.SmallInt(3).ToValue(),
					value.ToSymbol("Bar").ToValue(),
					value.SmallInt(10).ToValue(),
					value.SmallInt(0).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("#contains"), 1)),
					value.Ref(value.String("a")),
				},
			),
		},
		"set pattern": {
			input: `
			  a := ^[1, 5, -4]
				switch a
				case ^[1, _, -4] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 30,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 20,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.CALL_METHOD8), 6,
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.CALL_METHOD8), 8,
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 9,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 48),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(5).ToValue(),
						value.SmallInt(1).ToValue(),
						value.SmallInt(-4).ToValue(),
					)),
					value.Ref(value.SetMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(3).ToValue(),
					value.SmallInt(1).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("contains"), 1)),
					value.SmallInt(-4).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("contains"), 1)),
					value.Ref(value.String("a")),
				},
			),
		},
		"set pattern with rest elements": {
			input: `
			  a := ^[1, 5, -4]
				switch a
				case ^[1, *, -4] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 30,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 20,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.CALL_METHOD8), 6,
					byte(bytecode.JUMP_UNLESS), 0, 11,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.CALL_METHOD8), 8,
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 9,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 48),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(5).ToValue(),
						value.SmallInt(1).ToValue(),
						value.SmallInt(-4).ToValue(),
					)),
					value.Ref(value.SetMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(2).ToValue(),
					value.SmallInt(1).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("contains"), 1)),
					value.SmallInt(-4).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("contains"), 1)),
					value.Ref(value.String("a")),
				},
			),
		},

		"word set pattern": {
			input: `
			  a := ^['foo', 'bar']
				switch a
				case ^w[foo bar] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(77, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 24),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.Ref(value.String("foo")),
						value.Ref(value.String("bar")),
					)),
					value.Ref(value.SetMixin),
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.Ref(value.String("bar")),
						value.Ref(value.String("foo")),
					)),
					value.Ref(value.String("a")),
				},
			),
		},
		"symbol set pattern": {
			input: `
			  a := ^[:foo, :bar]
				switch a
				case ^s[foo bar] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(75, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 24),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.ToSymbol("foo").ToValue(),
						value.ToSymbol("bar").ToValue(),
					)),
					value.Ref(value.SetMixin),
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.ToSymbol("foo").ToValue(),
						value.ToSymbol("bar").ToValue(),
					)),
					value.Ref(value.String("a")),
				},
			),
		},
		"hex set pattern": {
			input: `
			  a := ^[0xff, 0x26]
				switch a
				case ^x[ff 26] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 24),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(38).ToValue(),
						value.SmallInt(255).ToValue(),
					)),
					value.Ref(value.SetMixin),
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(38).ToValue(),
						value.SmallInt(255).ToValue(),
					)),
					value.Ref(value.String("a")),
				},
			),
		},

		"bin set pattern": {
			input: `
			  a := ^[0b11, 0b10]
				switch a
				case ^b[11 10] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 24),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(3).ToValue(),
						value.SmallInt(2).ToValue(),
					)),
					value.Ref(value.SetMixin),
					value.Ref(vm.MustNewHashSetWithElements(
						nil,
						value.SmallInt(3).ToValue(),
						value.SmallInt(2).ToValue(),
					)),
					value.Ref(value.String("a")),
				},
			),
		},
		"list pattern": {
			input: `
			  a := [1, 5, [8, 3]]
				switch a
				case [1, < 8, [a, > 1 && < 5]] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 147,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 137,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.EQUAL),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 124,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 8,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 95,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 9,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 77,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 10,
					byte(bytecode.LOAD_VALUE8), 9,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 67,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 55,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.GREATER),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.JUMP_UNLESS), 0, 21,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 11,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 12,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(90, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 11),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 171),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.Ref(&value.ArrayList{
						value.SmallInt(8).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Ref(value.ListMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(3).ToValue(),
					value.SmallInt(0).ToValue(),
					value.SmallInt(1).ToValue(),
					value.SmallInt(8).ToValue(),
					value.SmallInt(2).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"word list pattern": {
			input: `
			  a := ['foo', 'bar']
				switch a
				case \w[foo bar] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(76, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 24),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayList{
						value.Ref(value.String("foo")),
						value.Ref(value.String("bar")),
					}),
					value.Ref(value.ListMixin),
					value.Ref(&value.ArrayList{
						value.Ref(value.String("foo")),
						value.Ref(value.String("bar")),
					}),
					value.Ref(value.String("a")),
				},
			),
		},

		"symbol list pattern": {
			input: `
			  a := [:foo, :bar]
				switch a
				case \s[foo bar] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(74, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 24),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
						value.ToSymbol("bar").ToValue(),
					}),
					value.Ref(value.ListMixin),
					value.Ref(&value.ArrayList{
						value.ToSymbol("foo").ToValue(),
						value.ToSymbol("bar").ToValue(),
					}),
					value.Ref(value.String("a")),
				},
			),
		},
		"hex list pattern": {
			input: `
			  a := [0xff, 0x26]
				switch a
				case \x[ff 26] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(72, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 24),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayList{
						value.SmallInt(255).ToValue(),
						value.SmallInt(38).ToValue(),
					}),
					value.Ref(value.ListMixin),
					value.Ref(&value.ArrayList{
						value.SmallInt(255).ToValue(),
						value.SmallInt(38).ToValue(),
					}),
					value.Ref(value.String("a")),
				},
			),
		},
		"bin list pattern": {
			input: `
			  a := [0b11, 0b10]
				switch a
				case \b[11 10] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 6,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.COPY),
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(72, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 24),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayList{
						value.SmallInt(3).ToValue(),
						value.SmallInt(2).ToValue(),
					}),
					value.Ref(value.ListMixin),
					value.Ref(&value.ArrayList{
						value.SmallInt(3).ToValue(),
						value.SmallInt(2).ToValue(),
					}),
					value.Ref(value.String("a")),
				},
			),
		},
		"list pattern with rest elements": {
			input: `
			  a := [1, 5, [-2, 8, 3, 6]]
				switch a
				case [*b, [< 0, *c]] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 7,
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 160,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 148,
					byte(bytecode.POP),

					// adjust the length variable
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SUBTRACT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),

					// create the iterator variable
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),

					// loop header
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 19,

					// loop body
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.SWAP),
					byte(bytecode.APPEND),
					byte(bytecode.POP),

					// i++
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.INCREMENT),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),

					byte(bytecode.LOOP), 0, 27,

					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.SUBSCRIPT),

					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.SET_LOCAL8), 5,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 76,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 7,
					byte(bytecode.SET_LOCAL8), 6,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 64,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 35,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SET_LOCAL8), 7,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 7,
					byte(bytecode.GET_LOCAL8), 6,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 19,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.GET_LOCAL8), 7,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.GET_LOCAL8), 5,
					byte(bytecode.SWAP),
					byte(bytecode.APPEND),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 7,
					byte(bytecode.INCREMENT),
					byte(bytecode.SET_LOCAL8), 7,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.INCREMENT),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 8,
					byte(bytecode.LEAVE_SCOPE16), 7, 6,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.LEAVE_SCOPE16), 7, 6,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(87, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 11),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 191),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.Ref(&value.ArrayList{
						value.SmallInt(-2).ToValue(),
						value.SmallInt(8).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(6).ToValue(),
					}),
					value.Ref(value.ListMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(1).ToValue(),
					value.SmallInt(0).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.Ref(value.String("a")),
				},
			),
		},
		"list pattern with unnamed rest elements": {
			input: `
			  a := [1, 5, [-2, 8, 3, 6]]
				switch a
				case [*, [< 0, *]] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 5,
					byte(bytecode.UNDEFINED),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.COPY),
					byte(bytecode.NEW_ARRAY_LIST8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 92,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 4,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 80,
					byte(bytecode.POP),

					// create the iterator variable
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SUBTRACT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.SUBSCRIPT),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 48,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 6,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 36,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.SET_LOCAL8), 5,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.INCREMENT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 8,
					byte(bytecode.LEAVE_SCOPE16), 5, 4,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.LEAVE_SCOPE16), 5, 4,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(85, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 11),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 116),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayList{
						value.SmallInt(1).ToValue(),
						value.SmallInt(5).ToValue(),
					}),
					value.Ref(&value.ArrayList{
						value.SmallInt(-2).ToValue(),
						value.SmallInt(8).ToValue(),
						value.SmallInt(3).ToValue(),
						value.SmallInt(6).ToValue(),
					}),
					value.Ref(value.ListMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(1).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"tuple pattern": {
			input: `
			  a := %[1, 5, %[8, 3]]
				switch a
				case %[1, < 8, %[a, > 1 && < 5]] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 147,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 137,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.EQUAL),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 124,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 95,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 8,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 77,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 9,
					byte(bytecode.LOAD_VALUE8), 8,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 67,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 55,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.GREATER),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.JUMP_UNLESS), 0, 21,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 10,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 11,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.LEAVE_SCOPE16), 2, 1,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(94, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 171),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.SmallInt(5).ToValue(),
						value.Ref(&value.ArrayTuple{
							value.SmallInt(8).ToValue(),
							value.SmallInt(3).ToValue(),
						}),
					}),
					value.Ref(value.TupleMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(3).ToValue(),
					value.SmallInt(0).ToValue(),
					value.SmallInt(1).ToValue(),
					value.SmallInt(8).ToValue(),
					value.SmallInt(2).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},

		"word tuple pattern": {
			input: `
			  a := %['foo', 'bar']
				switch a
				case %w[foo bar] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(77, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 23),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayTuple{
						value.Ref(value.String("foo")),
						value.Ref(value.String("bar")),
					}),
					value.Ref(value.TupleMixin),
					value.Ref(&value.ArrayTuple{
						value.Ref(value.String("foo")),
						value.Ref(value.String("bar")),
					}),
					value.Ref(value.String("a")),
				},
			),
		},
		"symbol tuple pattern": {
			input: `
			  a := %[:foo, :bar]
				switch a
				case %s[foo bar] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(75, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 23),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
						value.ToSymbol("bar").ToValue(),
					}),
					value.Ref(value.TupleMixin),
					value.Ref(&value.ArrayTuple{
						value.ToSymbol("foo").ToValue(),
						value.ToSymbol("bar").ToValue(),
					}),
					value.Ref(value.String("a")),
				},
			),
		},
		"hex tuple pattern": {
			input: `
			  a := %[0xff, 0x26]
				switch a
				case %x[ff 26] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 23),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayTuple{
						value.SmallInt(255).ToValue(),
						value.SmallInt(38).ToValue(),
					}),
					value.Ref(value.TupleMixin),
					value.Ref(&value.ArrayTuple{
						value.SmallInt(255).ToValue(),
						value.SmallInt(38).ToValue(),
					}),
					value.Ref(value.String("a")),
				},
			),
		},
		"bin tuple pattern": {
			input: `
			  a := %[0b11, 0b10]
				switch a
				case %b[11 10] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LAX_EQUAL),

					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.JUMP), 0, 3,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 23),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayTuple{
						value.SmallInt(3).ToValue(),
						value.SmallInt(2).ToValue(),
					}),
					value.Ref(value.TupleMixin),
					value.Ref(&value.ArrayTuple{
						value.SmallInt(3).ToValue(),
						value.SmallInt(2).ToValue(),
					}),
					value.Ref(value.String("a")),
				},
			),
		},
		"tuple pattern with rest elements": {
			input: `
			  a := %[1, 5, %[-2, 8, 3, 6]]
				switch a
				case %[*b, %[< 0, *c]] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 7,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 160,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 148,
					byte(bytecode.POP),

					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SUBTRACT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),

					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),

					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 19,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.SWAP),
					byte(bytecode.APPEND),
					byte(bytecode.POP),

					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.INCREMENT),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),

					byte(bytecode.LOOP), 0, 27,

					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.SUBSCRIPT),

					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.NEW_ARRAY_LIST8), 0,
					byte(bytecode.SET_LOCAL8), 5,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 76,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 6,
					byte(bytecode.SET_LOCAL8), 6,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 64,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 35,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SET_LOCAL8), 7,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 7,
					byte(bytecode.GET_LOCAL8), 6,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 19,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.GET_LOCAL8), 7,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.GET_LOCAL8), 5,
					byte(bytecode.SWAP),
					byte(bytecode.APPEND),
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 7,
					byte(bytecode.INCREMENT),
					byte(bytecode.SET_LOCAL8), 7,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.INCREMENT),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.LEAVE_SCOPE16), 7, 6,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.LEAVE_SCOPE16), 7, 6,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(91, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 191),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.SmallInt(5).ToValue(),
						value.Ref(&value.ArrayTuple{
							value.SmallInt(-2).ToValue(),
							value.SmallInt(8).ToValue(),
							value.SmallInt(3).ToValue(),
							value.SmallInt(6).ToValue(),
						}),
					}),
					value.Ref(value.TupleMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(1).ToValue(),
					value.SmallInt(0).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.Ref(value.String("a")),
				},
			),
		},
		"tuple pattern with unnamed rest elements": {
			input: `
			  a := %[1, 5, %[-2, 8, 3, 6]]
				switch a
				case %[*, %[< 0, *]] then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 5,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 92,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 80,
					byte(bytecode.POP),

					byte(bytecode.GET_LOCAL8), 2,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SUBTRACT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.SUBSCRIPT),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 48,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 5,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 36,
					byte(bytecode.POP),

					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 7,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.SET_LOCAL8), 5,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.INCREMENT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.LEAVE_SCOPE16), 5, 4,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.LEAVE_SCOPE16), 5, 4,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(89, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 116),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayTuple{
						value.SmallInt(1).ToValue(),
						value.SmallInt(5).ToValue(),
						value.Ref(&value.ArrayTuple{
							value.SmallInt(-2).ToValue(),
							value.SmallInt(8).ToValue(),
							value.SmallInt(3).ToValue(),
							value.SmallInt(6).ToValue(),
						}),
					}),
					value.Ref(value.TupleMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(1).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(0).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},

		"map pattern": {
			input: `
			  a := { 1 => 2, foo: :bar, "baz" => { dupa: [8, 3] } }
				switch a
				case { 1 => < 8, foo, "baz" => { dupa: [a, > 1 && < 5] } } then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 3,
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
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 149,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 8,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 120,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 9,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 95,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 10,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 77,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 11,
					byte(bytecode.LOAD_VALUE8), 12,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 67,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 13,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 55,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.GREATER),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.JUMP_UNLESS), 0, 21,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 14,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 15,
					byte(bytecode.LEAVE_SCOPE16), 3, 2,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.LEAVE_SCOPE16), 3, 2,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(152, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 20),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 173),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(vm.MustNewHashMapWithCapacityAndElements(
						nil,
						3,
						value.Pair{
							Key:   value.SmallInt(1).ToValue(),
							Value: value.SmallInt(2).ToValue(),
						},
						value.Pair{
							Key:   value.ToSymbol("foo").ToValue(),
							Value: value.ToSymbol("bar").ToValue(),
						},
					)),
					value.Ref(value.String("baz")),
					value.Ref(vm.MustNewHashMapWithCapacityAndElements(nil, 1)),
					value.ToSymbol("dupa").ToValue(),
					value.Ref(&value.ArrayList{
						value.SmallInt(8).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Ref(value.MapMixin),
					value.SmallInt(1).ToValue(),
					value.SmallInt(8).ToValue(),
					value.ToSymbol("foo").ToValue(),
					value.Ref(value.ListMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(2).ToValue(),
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"record pattern": {
			input: `
			  a := %{ 1 => 2, foo: :bar, "baz" => %{ dupa: [8, 3] } }
				switch a
				case %{ 1 => < 8, foo, "baz" => %{ dupa: [a, > 1 && < 5] } } then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.COPY),
					byte(bytecode.NEW_HASH_RECORD8), 1,
					byte(bytecode.NEW_HASH_RECORD8), 1,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 149,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 8,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 120,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 9,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 6,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 95,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 10,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 77,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 11,
					byte(bytecode.LOAD_VALUE8), 12,
					byte(bytecode.EQUAL),
					byte(bytecode.JUMP_UNLESS), 0, 67,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 13,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 55,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.SUBSCRIPT),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.GREATER),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.JUMP_UNLESS), 0, 21,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 14,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 2,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 15,
					byte(bytecode.LEAVE_SCOPE16), 3, 2,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.LEAVE_SCOPE16), 3, 2,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(156, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 18),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 173),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(vm.MustNewHashRecordWithElements(
						nil,
						value.Pair{
							Key:   value.SmallInt(1).ToValue(),
							Value: value.SmallInt(2).ToValue(),
						},
						value.Pair{
							Key:   value.ToSymbol("foo").ToValue(),
							Value: value.ToSymbol("bar").ToValue(),
						},
					)),
					value.Ref(value.String("baz")),
					value.Ref(&value.HashRecord{}),
					value.ToSymbol("dupa").ToValue(),
					value.Ref(&value.ArrayList{
						value.SmallInt(8).ToValue(),
						value.SmallInt(3).ToValue(),
					}),
					value.Ref(value.RecordMixin),
					value.SmallInt(1).ToValue(),
					value.SmallInt(8).ToValue(),
					value.ToSymbol("foo").ToValue(),
					value.Ref(value.ListMixin),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(2).ToValue(),
					value.SmallInt(0).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.String("a")),
				},
			),
		},
		"object pattern": {
			input: `
			  a := [0b11, 0b10]
				switch a
				case ::Std::ArrayList(length: > 1 && < 5 as l, first) then "a"
				end
			`,
			want: vm.NewBytecodeFunctionNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.COPY),
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.DUP),
					byte(bytecode.GET_CONST8), 2,
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 62,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 3,
					byte(bytecode.SET_LOCAL8), 2,
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 4,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.GREATER),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.JUMP_UNLESS), 0, 21,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.LOAD_VALUE8), 5,
					byte(bytecode.DUP_2),
					byte(bytecode.SWAP),
					byte(bytecode.GET_CLASS),
					byte(bytecode.IS_A),
					byte(bytecode.JUMP_UNLESS), 0, 5,
					byte(bytecode.POP),
					byte(bytecode.LESS),
					byte(bytecode.JUMP), 0, 4,
					byte(bytecode.POP),
					byte(bytecode.POP_2),
					byte(bytecode.FALSE),
					byte(bytecode.POP_SKIP_ONE),
					byte(bytecode.JUMP_UNLESS), 0, 8,
					byte(bytecode.POP),
					byte(bytecode.DUP),
					byte(bytecode.CALL_METHOD8), 6,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					byte(bytecode.JUMP_UNLESS), 0, 10,
					byte(bytecode.POP_2),
					byte(bytecode.LOAD_VALUE8), 7,
					byte(bytecode.LEAVE_SCOPE16), 3, 2,
					byte(bytecode.JUMP), 0, 6,
					byte(bytecode.LEAVE_SCOPE16), 3, 2,
					byte(bytecode.POP),

					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(111, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
					bytecode.NewLineInfo(2, 6),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(4, 86),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(5, 2),
				},
				[]value.Value{
					value.Undefined,
					value.Ref(&value.ArrayList{
						value.SmallInt(3).ToValue(),
						value.SmallInt(2).ToValue(),
					}),
					value.ToSymbol("Std::ArrayList").ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
					value.SmallInt(1).ToValue(),
					value.SmallInt(5).ToValue(),
					value.Ref(value.NewCallSiteInfo(value.ToSymbol("first"), 0)),
					value.Ref(value.String("a")),
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
