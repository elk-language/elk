package compiler

// import (
// 	"testing"

// 	"github.com/elk-language/elk/bitfield"
// 	"github.com/elk-language/elk/bytecode"
// 	"github.com/elk-language/elk/value"
// 	"github.com/elk-language/elk/vm"
// )

// func TestSwitch(t *testing.T) {
// 	tests := testTable{
// 		"with a few literal cases": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case true then "a"
// 				case false then "b"
// 				case 0 then "c"
// 				case 1 then "d"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.JUMP), 0, 47,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.FALSE),
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 33,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 18,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(120, 8, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 14),
// 					bytecode.NewLineInfo(5, 14),
// 					bytecode.NewLineInfo(6, 15),
// 					bytecode.NewLineInfo(7, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(8, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.String("a"),
// 					value.String("b"),
// 					value.String("c"),
// 					value.SmallInt(1),
// 					value.String("d"),
// 				},
// 			),
// 		},
// 		"with else": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case true then "a"
// 				case false then "b"
// 				else "c"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.JUMP), 0, 18,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.FALSE),
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(93, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 14),
// 					bytecode.NewLineInfo(5, 14),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(6, 2),
// 					bytecode.NewLineInfo(7, 1),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.String("a"),
// 					value.String("b"),
// 					value.String("c"),
// 				},
// 			),
// 		},
// 		"literal true": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case true then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 14),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal false": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case false then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.FALSE),
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(57, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 14),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal nil": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case nil then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(55, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 14),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal string": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case "foo" then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(57, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.String("foo"),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal raw string": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 'foo' then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(57, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.String("foo"),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal interpolated string": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case "f${a}" then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.NEW_STRING8), 2,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(59, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 19),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.String("f"),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal symbol": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case :foo then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.ToSymbol("foo"),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal interpolated symbol": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case :"f${a}" then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.NEW_SYMBOL8), 2,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(60, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 19),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.String("f"),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal int": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(53, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal int64": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5i64 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.Int64(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal uint64": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5u64 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.UInt64(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal int32": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5i32 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.Int32(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal uint32": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5u32 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.UInt32(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal int16": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5i16 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.Int16(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal uint16": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5u16 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.UInt16(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal int8": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5i8 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(55, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.Int8(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal uint8": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5u8 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(55, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.UInt8(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal float": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5.8 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(55, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.Float(5.8),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal float64": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5.8f64 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(58, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.Float64(5.8),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal float32": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5.8f32 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(58, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.Float32(5.8),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal negative float32": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case -5.8f32 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(59, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.Float32(-5.8),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"literal big float": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case 5.8bf then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(57, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.NewBigFloat(5.8),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"root constant": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case ::Foo then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.ROOT),
// 					byte(bytecode.GET_MOD_CONST8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(57, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 16),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.ToSymbol("Foo"),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"constant lookup": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case ::Foo::Bar then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.ROOT),
// 					byte(bytecode.GET_MOD_CONST8), 1,
// 					byte(bytecode.GET_MOD_CONST8), 2,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(62, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 18),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.ToSymbol("Foo"),
// 					value.ToSymbol("Bar"),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"negative constant lookup": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case -::Foo::Bar then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.ROOT),
// 					byte(bytecode.GET_MOD_CONST8), 1,
// 					byte(bytecode.GET_MOD_CONST8), 2,
// 					byte(bytecode.NEGATE),
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(63, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 19),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.ToSymbol("Foo"),
// 					value.ToSymbol("Bar"),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"less pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case < 5 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(55, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 31),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"less than root constant": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case < ::Foo then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.ROOT),
// 					byte(bytecode.GET_MOD_CONST8), 1,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(59, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 32),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.ToSymbol("Foo"),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"less than negative root constant": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case < -::Foo then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.ROOT),
// 					byte(bytecode.GET_MOD_CONST8), 1,
// 					byte(bytecode.NEGATE),
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(60, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 33),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.ToSymbol("Foo"),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"less equal pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case <= 5 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS_EQUAL),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 31),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"greater pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case > 5 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.GREATER),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(55, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 31),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"greater equal pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case >= 5 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.GREATER_EQUAL),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 31),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"equal pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case == 5 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"equal regex pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case == %/fo+/ then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(61, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.MustCompileRegex("fo+", bitfield.BitField8{}),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"equal local pattern": {
// 			input: `
// 			  a := 0
// 				b := 2
// 				switch a
// 				case == b then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 2,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.GET_LOCAL8), 4,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(67, 6, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 5),
// 					bytecode.NewLineInfo(4, 2),
// 					bytecode.NewLineInfo(5, 15),
// 					bytecode.NewLineInfo(4, 1),
// 					bytecode.NewLineInfo(6, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(2),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"not equal pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case != 5 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.NOT_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"lax equal pattern": {
// 			input: `
// 			  a := 0
// 				b := 2
// 				switch a
// 				case =~ b then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 2,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.GET_LOCAL8), 4,
// 					byte(bytecode.LAX_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(67, 6, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 5),
// 					bytecode.NewLineInfo(4, 2),
// 					bytecode.NewLineInfo(5, 15),
// 					bytecode.NewLineInfo(4, 1),
// 					bytecode.NewLineInfo(6, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(2),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"lax not equal pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case !~ 5 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.LAX_NOT_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"strict equal pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case === 5 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.STRICT_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(57, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"strict not equal pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case !== 5 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.STRICT_NOT_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(57, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"strict not equal negative pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case !== -5 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.STRICT_NOT_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(58, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 15),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(-5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"regex pattern": {
// 			input: `
// 			  a := "foo"
// 				switch a
// 				case %/fo+/ then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.CALL_METHOD8), 2,
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(62, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 17),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.String("foo"),
// 					value.MustCompileRegex("fo+", bitfield.BitField8{}),
// 					value.NewCallSiteInfo(value.ToSymbol("matches"), 1, nil),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"variable pattern": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case n then n + 2
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 2,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 13,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.GET_LOCAL8), 4,
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.ADD),
// 					byte(bytecode.LEAVE_SCOPE16), 4, 1,
// 					byte(bytecode.JUMP), 0, 6,
// 					byte(bytecode.LEAVE_SCOPE16), 4, 1,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(55, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 23),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.SmallInt(2),
// 				},
// 			),
// 		},
// 		"range": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case -2...9 then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.CALL_METHOD8), 2,
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(58, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 17),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.NewClosedRange(value.SmallInt(-2), value.SmallInt(9)),
// 					value.NewCallSiteInfo(value.ToSymbol("#contains"), 1, nil),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"range with constants": {
// 			input: `
// 			  a := 0
// 				switch a
// 				case ::Foo...-::Bar then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.ROOT),
// 					byte(bytecode.GET_MOD_CONST8), 1,
// 					byte(bytecode.ROOT),
// 					byte(bytecode.GET_MOD_CONST8), 2,
// 					byte(bytecode.NEGATE),
// 					byte(bytecode.NEW_RANGE), bytecode.CLOSED_RANGE_FLAG,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.CALL_METHOD8), 3,
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(66, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 24),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.SmallInt(0),
// 					value.ToSymbol("Foo"),
// 					value.ToSymbol("Bar"),
// 					value.NewCallSiteInfo(value.ToSymbol("#contains"), 1, nil),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"set pattern": {
// 			input: `
// 			  a := ^[1, 5, -4]
// 				switch a
// 				case ^[1, _, -4] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.COPY),
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 30,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 20,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.CALL_METHOD8), 5,
// 					byte(bytecode.JUMP_UNLESS), 0, 11,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.CALL_METHOD8), 7,
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 8,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(73, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 48),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					vm.MustNewHashSetWithElements(
// 						nil,
// 						value.SmallInt(5),
// 						value.SmallInt(1),
// 						value.SmallInt(-4),
// 					),
// 					value.SetMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(3),
// 					value.SmallInt(1),
// 					value.NewCallSiteInfo(value.ToSymbol("contains"), 1, nil),
// 					value.SmallInt(-4),
// 					value.NewCallSiteInfo(value.ToSymbol("contains"), 1, nil),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"set pattern with rest elements": {
// 			input: `
// 			  a := ^[1, 5, -4]
// 				switch a
// 				case ^[1, *, -4] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.COPY),
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 30,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.GREATER_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 20,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.CALL_METHOD8), 5,
// 					byte(bytecode.JUMP_UNLESS), 0, 11,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.CALL_METHOD8), 7,
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 8,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(73, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 48),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					vm.MustNewHashSetWithElements(
// 						nil,
// 						value.SmallInt(5),
// 						value.SmallInt(1),
// 						value.SmallInt(-4),
// 					),
// 					value.SetMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(2),
// 					value.SmallInt(1),
// 					value.NewCallSiteInfo(value.ToSymbol("contains"), 1, nil),
// 					value.SmallInt(-4),
// 					value.NewCallSiteInfo(value.ToSymbol("contains"), 1, nil),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"word set pattern": {
// 			input: `
// 			  a := ^['foo', 'bar']
// 				switch a
// 				case ^w[foo bar] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.COPY),
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 6,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.COPY),
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(77, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 24),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					vm.MustNewHashSetWithElements(
// 						nil,
// 						value.String("foo"),
// 						value.String("bar"),
// 					),
// 					value.SetMixin,
// 					vm.MustNewHashSetWithElements(
// 						nil,
// 						value.String("foo"),
// 						value.String("bar"),
// 					),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"symbol set pattern": {
// 			input: `
// 			  a := ^[:foo, :bar]
// 				switch a
// 				case ^s[foo bar] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.COPY),
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 6,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.COPY),
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(75, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 24),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					vm.MustNewHashSetWithElements(
// 						nil,
// 						value.ToSymbol("foo"),
// 						value.ToSymbol("bar"),
// 					),
// 					value.SetMixin,
// 					vm.MustNewHashSetWithElements(
// 						nil,
// 						value.ToSymbol("foo"),
// 						value.ToSymbol("bar"),
// 					),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"hex set pattern": {
// 			input: `
// 			  a := ^[0xff, 0x26]
// 				switch a
// 				case ^x[ff 26] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.COPY),
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 6,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.COPY),
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(73, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 24),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					vm.MustNewHashSetWithElements(
// 						nil,
// 						value.SmallInt(255),
// 						value.SmallInt(38),
// 					),
// 					value.SetMixin,
// 					vm.MustNewHashSetWithElements(
// 						nil,
// 						value.SmallInt(255),
// 						value.SmallInt(38),
// 					),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"bin set pattern": {
// 			input: `
// 			  a := ^[0b11, 0b10]
// 				switch a
// 				case ^b[11 10] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.COPY),
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 6,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.COPY),
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(73, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 24),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					vm.MustNewHashSetWithElements(
// 						nil,
// 						value.SmallInt(3),
// 						value.SmallInt(2),
// 					),
// 					value.SetMixin,
// 					vm.MustNewHashSetWithElements(
// 						nil,
// 						value.SmallInt(3),
// 						value.SmallInt(2),
// 					),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"list pattern": {
// 			input: `
// 			  a := [1, 5, [8, 3]]
// 				switch a
// 				case [1, < 8, [a, > 1 && < 5]] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 2,
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.COPY),
// 					byte(bytecode.NEW_ARRAY_LIST8), 1,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 147,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 3,
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 137,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 124,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 7,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 95,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 8,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 9,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 77,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 10,
// 					byte(bytecode.LOAD_VALUE8), 8,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 67,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 55,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.GREATER),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.JUMP_UNLESS), 0, 21,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 11,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 10,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 12,
// 					byte(bytecode.LEAVE_SCOPE16), 4, 1,
// 					byte(bytecode.JUMP), 0, 6,
// 					byte(bytecode.LEAVE_SCOPE16), 4, 1,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(90, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 13),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 171),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayList{
// 						value.SmallInt(1),
// 						value.SmallInt(5),
// 					},
// 					&value.ArrayList{
// 						value.SmallInt(8),
// 						value.SmallInt(3),
// 					},
// 					value.ListMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(3),
// 					value.SmallInt(0),
// 					value.SmallInt(1),
// 					value.SmallInt(8),
// 					value.SmallInt(2),
// 					value.ListMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"word list pattern": {
// 			input: `
// 			  a := ['foo', 'bar']
// 				switch a
// 				case \w[foo bar] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.COPY),
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 6,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.COPY),
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(76, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 24),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayList{
// 						value.String("foo"),
// 						value.String("bar"),
// 					},
// 					value.ListMixin,
// 					&value.ArrayList{
// 						value.String("foo"),
// 						value.String("bar"),
// 					},
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"symbol list pattern": {
// 			input: `
// 			  a := [:foo, :bar]
// 				switch a
// 				case \s[foo bar] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.COPY),
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 6,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.COPY),
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(74, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 24),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayList{
// 						value.ToSymbol("foo"),
// 						value.ToSymbol("bar"),
// 					},
// 					value.ListMixin,
// 					&value.ArrayList{
// 						value.ToSymbol("foo"),
// 						value.ToSymbol("bar"),
// 					},
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"hex list pattern": {
// 			input: `
// 			  a := [0xff, 0x26]
// 				switch a
// 				case \x[ff 26] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.COPY),
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 6,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.COPY),
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(72, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 24),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayList{
// 						value.SmallInt(255),
// 						value.SmallInt(38),
// 					},
// 					value.ListMixin,
// 					&value.ArrayList{
// 						value.SmallInt(255),
// 						value.SmallInt(38),
// 					},
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"bin list pattern": {
// 			input: `
// 			  a := [0b11, 0b10]
// 				switch a
// 				case \b[11 10] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.COPY),
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 6,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.COPY),
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(72, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 24),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayList{
// 						value.SmallInt(3),
// 						value.SmallInt(2),
// 					},
// 					value.ListMixin,
// 					&value.ArrayList{
// 						value.SmallInt(3),
// 						value.SmallInt(2),
// 					},
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"list pattern with rest elements": {
// 			input: `
// 			  a := [1, 5, [-2, 8, 3, 6]]
// 				switch a
// 				case [*b, [< 0, *c]] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 7,
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.COPY),
// 					byte(bytecode.NEW_ARRAY_LIST8), 1,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.NEW_ARRAY_LIST8), 0,
// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 160,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 3,
// 					byte(bytecode.SET_LOCAL8), 5,
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.GREATER_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 148,
// 					byte(bytecode.POP),

// 					// adjust the length variable
// 					byte(bytecode.GET_LOCAL8), 5,
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.SUBTRACT),
// 					byte(bytecode.SET_LOCAL8), 5,
// 					byte(bytecode.POP),

// 					// create the iterator variable
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.SET_LOCAL8), 6,
// 					byte(bytecode.POP),

// 					// loop header
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.GET_LOCAL8), 5,
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP_UNLESS), 0, 19,

// 					// loop body
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.GET_LOCAL8), 4,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.APPEND),
// 					byte(bytecode.POP),

// 					// i++
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.INCREMENT),
// 					byte(bytecode.SET_LOCAL8), 6,
// 					byte(bytecode.POP),

// 					byte(bytecode.LOOP), 0, 27,

// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.SUBSCRIPT),

// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.NEW_ARRAY_LIST8), 0,
// 					byte(bytecode.SET_LOCAL8), 7,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 76,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 7,
// 					byte(bytecode.SET_LOCAL8), 8,
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.GREATER_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 64,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 35,
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.SET_LOCAL8), 9,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 9,
// 					byte(bytecode.GET_LOCAL8), 8,
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP_UNLESS), 0, 19,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.GET_LOCAL8), 9,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.GET_LOCAL8), 7,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.APPEND),
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 9,
// 					byte(bytecode.INCREMENT),
// 					byte(bytecode.SET_LOCAL8), 9,
// 					byte(bytecode.POP),
// 					byte(bytecode.LOOP), 0, 27,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 8,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.INCREMENT),
// 					byte(bytecode.SET_LOCAL8), 6,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 10,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 8,
// 					byte(bytecode.LEAVE_SCOPE16), 9, 6,
// 					byte(bytecode.JUMP), 0, 6,
// 					byte(bytecode.LEAVE_SCOPE16), 9, 6,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(87, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 13),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 191),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayList{
// 						value.SmallInt(1),
// 						value.SmallInt(5),
// 					},
// 					&value.ArrayList{
// 						value.SmallInt(-2),
// 						value.SmallInt(8),
// 						value.SmallInt(3),
// 						value.SmallInt(6),
// 					},
// 					value.ListMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(1),
// 					value.SmallInt(0),
// 					value.ListMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"list pattern with unnamed rest elements": {
// 			input: `
// 			  a := [1, 5, [-2, 8, 3, 6]]
// 				switch a
// 				case [*, [< 0, *]] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 5,
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.COPY),
// 					byte(bytecode.NEW_ARRAY_LIST8), 1,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 92,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 3,
// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.GREATER_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 80,
// 					byte(bytecode.POP),

// 					// create the iterator variable
// 					byte(bytecode.GET_LOCAL8), 4,
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.SUBTRACT),
// 					byte(bytecode.SET_LOCAL8), 5,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.GET_LOCAL8), 5,
// 					byte(bytecode.SUBSCRIPT),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 48,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 6,
// 					byte(bytecode.SET_LOCAL8), 6,
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.GREATER_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 36,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 7,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 7,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.SET_LOCAL8), 7,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 8,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 5,
// 					byte(bytecode.INCREMENT),
// 					byte(bytecode.SET_LOCAL8), 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 10,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 8,
// 					byte(bytecode.LEAVE_SCOPE16), 7, 4,
// 					byte(bytecode.JUMP), 0, 6,
// 					byte(bytecode.LEAVE_SCOPE16), 7, 4,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(85, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 13),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 116),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayList{
// 						value.SmallInt(1),
// 						value.SmallInt(5),
// 					},
// 					&value.ArrayList{
// 						value.SmallInt(-2),
// 						value.SmallInt(8),
// 						value.SmallInt(3),
// 						value.SmallInt(6),
// 					},
// 					value.ListMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(1),
// 					value.ListMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(0),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"tuple pattern": {
// 			input: `
// 			  a := %[1, 5, %[8, 3]]
// 				switch a
// 				case %[1, < 8, %[a, > 1 && < 5]] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 2,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 147,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 137,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 124,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 95,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 7,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 8,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 77,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 9,
// 					byte(bytecode.LOAD_VALUE8), 7,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 67,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 55,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.GREATER),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.JUMP_UNLESS), 0, 21,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 10,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 10,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 11,
// 					byte(bytecode.LEAVE_SCOPE16), 4, 1,
// 					byte(bytecode.JUMP), 0, 6,
// 					byte(bytecode.LEAVE_SCOPE16), 4, 1,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(94, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 171),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayTuple{
// 						value.SmallInt(1),
// 						value.SmallInt(5),
// 						&value.ArrayTuple{
// 							value.SmallInt(8),
// 							value.SmallInt(3),
// 						},
// 					},
// 					value.TupleMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(3),
// 					value.SmallInt(0),
// 					value.SmallInt(1),
// 					value.SmallInt(8),
// 					value.SmallInt(2),
// 					value.TupleMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"word tuple pattern": {
// 			input: `
// 			  a := %['foo', 'bar']
// 				switch a
// 				case %w[foo bar] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(77, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 23),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayTuple{
// 						value.String("foo"),
// 						value.String("bar"),
// 					},
// 					value.TupleMixin,
// 					&value.ArrayTuple{
// 						value.String("foo"),
// 						value.String("bar"),
// 					},
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"symbol tuple pattern": {
// 			input: `
// 			  a := %[:foo, :bar]
// 				switch a
// 				case %s[foo bar] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(75, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 23),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayTuple{
// 						value.ToSymbol("foo"),
// 						value.ToSymbol("bar"),
// 					},
// 					value.TupleMixin,
// 					&value.ArrayTuple{
// 						value.ToSymbol("foo"),
// 						value.ToSymbol("bar"),
// 					},
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"hex tuple pattern": {
// 			input: `
// 			  a := %[0xff, 0x26]
// 				switch a
// 				case %x[ff 26] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(73, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 23),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayTuple{
// 						value.SmallInt(255),
// 						value.SmallInt(38),
// 					},
// 					value.TupleMixin,
// 					&value.ArrayTuple{
// 						value.SmallInt(255),
// 						value.SmallInt(38),
// 					},
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"bin tuple pattern": {
// 			input: `
// 			  a := %[0b11, 0b10]
// 				switch a
// 				case %b[11 10] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 1,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.LAX_EQUAL),

// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.JUMP), 0, 3,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(73, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 23),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayTuple{
// 						value.SmallInt(3),
// 						value.SmallInt(2),
// 					},
// 					value.TupleMixin,
// 					&value.ArrayTuple{
// 						value.SmallInt(3),
// 						value.SmallInt(2),
// 					},
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"tuple pattern with rest elements": {
// 			input: `
// 			  a := %[1, 5, %[-2, 8, 3, 6]]
// 				switch a
// 				case %[*b, %[< 0, *c]] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 7,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.NEW_ARRAY_LIST8), 0,
// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 160,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 2,
// 					byte(bytecode.SET_LOCAL8), 5,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.GREATER_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 148,
// 					byte(bytecode.POP),

// 					// adjust the length variable
// 					byte(bytecode.GET_LOCAL8), 5,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.SUBTRACT),
// 					byte(bytecode.SET_LOCAL8), 5,
// 					byte(bytecode.POP),

// 					// create the iterator variable
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.SET_LOCAL8), 6,
// 					byte(bytecode.POP),

// 					// loop header
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.GET_LOCAL8), 5,
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP_UNLESS), 0, 19,

// 					// loop body
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.GET_LOCAL8), 4,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.APPEND),
// 					byte(bytecode.POP),

// 					// i++
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.INCREMENT),
// 					byte(bytecode.SET_LOCAL8), 6,
// 					byte(bytecode.POP),

// 					byte(bytecode.LOOP), 0, 27,

// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.SUBSCRIPT),

// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.NEW_ARRAY_LIST8), 0,
// 					byte(bytecode.SET_LOCAL8), 7,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 76,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 6,
// 					byte(bytecode.SET_LOCAL8), 8,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.GREATER_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 64,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 35,
// 					byte(bytecode.POP),
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.SET_LOCAL8), 9,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 9,
// 					byte(bytecode.GET_LOCAL8), 8,
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP_UNLESS), 0, 19,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.GET_LOCAL8), 9,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.GET_LOCAL8), 7,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.APPEND),
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 9,
// 					byte(bytecode.INCREMENT),
// 					byte(bytecode.SET_LOCAL8), 9,
// 					byte(bytecode.POP),
// 					byte(bytecode.LOOP), 0, 27,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 8,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.INCREMENT),
// 					byte(bytecode.SET_LOCAL8), 6,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 10,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 7,
// 					byte(bytecode.LEAVE_SCOPE16), 9, 6,
// 					byte(bytecode.JUMP), 0, 6,
// 					byte(bytecode.LEAVE_SCOPE16), 9, 6,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(91, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 191),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayTuple{
// 						value.SmallInt(1),
// 						value.SmallInt(5),
// 						&value.ArrayTuple{
// 							value.SmallInt(-2),
// 							value.SmallInt(8),
// 							value.SmallInt(3),
// 							value.SmallInt(6),
// 						},
// 					},
// 					value.TupleMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(1),
// 					value.SmallInt(0),
// 					value.TupleMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"tuple pattern with unnamed rest elements": {
// 			input: `
// 			  a := %[1, 5, %[-2, 8, 3, 6]]
// 				switch a
// 				case %[*, %[< 0, *]] then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 5,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 92,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 2,
// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.GREATER_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 80,
// 					byte(bytecode.POP),

// 					// create the iterator variable
// 					byte(bytecode.GET_LOCAL8), 4,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.SUBTRACT),
// 					byte(bytecode.SET_LOCAL8), 5,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.GET_LOCAL8), 5,
// 					byte(bytecode.SUBSCRIPT),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 48,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 5,
// 					byte(bytecode.SET_LOCAL8), 6,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.GREATER_EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 36,
// 					byte(bytecode.POP),

// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 7,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 6,
// 					byte(bytecode.SET_LOCAL8), 7,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 8,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 5,
// 					byte(bytecode.INCREMENT),
// 					byte(bytecode.SET_LOCAL8), 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 10,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 7,
// 					byte(bytecode.LEAVE_SCOPE16), 7, 4,
// 					byte(bytecode.JUMP), 0, 6,
// 					byte(bytecode.LEAVE_SCOPE16), 7, 4,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(89, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 7),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 116),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayTuple{
// 						value.SmallInt(1),
// 						value.SmallInt(5),
// 						&value.ArrayTuple{
// 							value.SmallInt(-2),
// 							value.SmallInt(8),
// 							value.SmallInt(3),
// 							value.SmallInt(6),
// 						},
// 					},
// 					value.TupleMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(1),
// 					value.TupleMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(0),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"map pattern": {
// 			input: `
// 			  a := { 1 => 2, foo: :bar, "baz" => { dupa: [8, 3] } }
// 				switch a
// 				case { 1 => < 8, foo, "baz" => { dupa: [a, > 1 && < 5] } } then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 3,
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.COPY),
// 					byte(bytecode.NEW_HASH_MAP8), 1,
// 					byte(bytecode.NEW_HASH_MAP8), 1,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 149,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 7,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 120,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 8,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 9,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 95,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 10,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 77,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 11,
// 					byte(bytecode.LOAD_VALUE8), 12,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 67,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 13,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.SET_LOCAL8), 5,
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 55,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.GREATER),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.JUMP_UNLESS), 0, 21,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 14,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 10,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 15,
// 					byte(bytecode.LEAVE_SCOPE16), 5, 2,
// 					byte(bytecode.JUMP), 0, 6,
// 					byte(bytecode.LEAVE_SCOPE16), 5, 2,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(152, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 22),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 173),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					vm.MustNewHashMapWithCapacityAndElements(
// 						nil,
// 						3,
// 						value.Pair{
// 							Key:   value.SmallInt(1),
// 							Value: value.SmallInt(2),
// 						},
// 						value.Pair{
// 							Key:   value.ToSymbol("foo"),
// 							Value: value.ToSymbol("bar"),
// 						},
// 					),
// 					value.String("baz"),
// 					value.NewHashMap(1),
// 					value.ToSymbol("dupa"),
// 					&value.ArrayList{
// 						value.SmallInt(8),
// 						value.SmallInt(3),
// 					},
// 					value.MapMixin,
// 					value.SmallInt(1),
// 					value.SmallInt(8),
// 					value.ToSymbol("foo"),
// 					value.MapMixin,
// 					value.ListMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(2),
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"record pattern": {
// 			input: `
// 			  a := %{ 1 => 2, foo: :bar, "baz" => %{ dupa: [8, 3] } }
// 				switch a
// 				case %{ 1 => < 8, foo, "baz" => %{ dupa: [a, > 1 && < 5] } } then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 3,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.LOAD_VALUE8), 2,
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.COPY),
// 					byte(bytecode.NEW_HASH_RECORD8), 1,
// 					byte(bytecode.NEW_HASH_RECORD8), 1,
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 149,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 7,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 120,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 8,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 1,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 9,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 95,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 3,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 10,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 77,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 11,
// 					byte(bytecode.LOAD_VALUE8), 12,
// 					byte(bytecode.EQUAL),
// 					byte(bytecode.JUMP_UNLESS), 0, 67,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 13,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.SET_LOCAL8), 5,
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 55,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 6,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.GREATER),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.JUMP_UNLESS), 0, 21,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 14,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 10,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 15,
// 					byte(bytecode.LEAVE_SCOPE16), 5, 2,
// 					byte(bytecode.JUMP), 0, 6,
// 					byte(bytecode.LEAVE_SCOPE16), 5, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(156, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 20),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 173),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					vm.MustNewHashRecordWithCapacityAndElements(
// 						nil,
// 						3,
// 						value.Pair{
// 							Key:   value.SmallInt(1),
// 							Value: value.SmallInt(2),
// 						},
// 						value.Pair{
// 							Key:   value.ToSymbol("foo"),
// 							Value: value.ToSymbol("bar"),
// 						},
// 					),
// 					value.String("baz"),
// 					value.NewHashRecord(1),
// 					value.ToSymbol("dupa"),
// 					&value.ArrayList{
// 						value.SmallInt(8),
// 						value.SmallInt(3),
// 					},
// 					value.RecordMixin,
// 					value.SmallInt(1),
// 					value.SmallInt(8),
// 					value.ToSymbol("foo"),
// 					value.RecordMixin,
// 					value.ListMixin,
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(2),
// 					value.SmallInt(0),
// 					value.SmallInt(5),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 		"object pattern": {
// 			input: `
// 			  a := [0b11, 0b10]
// 				switch a
// 				case ::Std::ArrayList(length: > 1 && < 5 as l, uppercase) then "a"
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 3,
// 					byte(bytecode.LOAD_VALUE8), 0,
// 					byte(bytecode.COPY),
// 					byte(bytecode.SET_LOCAL8), 3,
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_LOCAL8), 3,
// 					byte(bytecode.DUP),
// 					byte(bytecode.ROOT),
// 					byte(bytecode.GET_MOD_CONST8), 1,
// 					byte(bytecode.GET_MOD_CONST8), 2,
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 62,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 3,
// 					byte(bytecode.SET_LOCAL8), 4,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 4,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.GREATER),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.JUMP_UNLESS), 0, 21,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE8), 5,
// 					byte(bytecode.DUP_N), 2,
// 					byte(bytecode.SWAP),
// 					byte(bytecode.GET_CLASS),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS), 0, 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LESS),
// 					byte(bytecode.JUMP), 0, 4,
// 					byte(bytecode.POP),
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.FALSE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS), 0, 8,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 6,
// 					byte(bytecode.SET_LOCAL8), 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_UNLESS), 0, 10,
// 					byte(bytecode.POP_N), 2,
// 					byte(bytecode.LOAD_VALUE8), 7,
// 					byte(bytecode.LEAVE_SCOPE16), 5, 2,
// 					byte(bytecode.JUMP), 0, 6,
// 					byte(bytecode.LEAVE_SCOPE16), 5, 2,
// 					byte(bytecode.POP),

// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(115, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(2, 8),
// 					bytecode.NewLineInfo(3, 2),
// 					bytecode.NewLineInfo(4, 89),
// 					bytecode.NewLineInfo(3, 1),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					&value.ArrayList{
// 						value.SmallInt(3),
// 						value.SmallInt(2),
// 					},
// 					value.ToSymbol("Std"),
// 					value.ToSymbol("ArrayList"),
// 					value.NewCallSiteInfo(value.ToSymbol("length"), 0, nil),
// 					value.SmallInt(1),
// 					value.SmallInt(5),
// 					value.NewCallSiteInfo(value.ToSymbol("uppercase"), 0, nil),
// 					value.String("a"),
// 				},
// 			),
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			compilerTest(tc, t)
// 		})
// 	}
// }
