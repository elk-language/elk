package compiler

import (
	"math"
	"math/big"
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

var classSymbol value.Symbol = value.SymbolTable.Add("class")
var moduleSymbol value.Symbol = value.SymbolTable.Add("module")
var mixinSymbol value.Symbol = value.SymbolTable.Add("mixin")
var mainSymbol value.Symbol = value.SymbolTable.Add("main")

// Represents a single compiler test case.
type testCase struct {
	input string
	want  *vm.BytecodeMethod
	err   errors.ErrorList
}

// Type of the compiler test table.
type testTable map[string]testCase

func compilerTest(tc testCase, t *testing.T) {
	t.Helper()

	got, err := CompileSource("main", tc.input)
	opts := comparer.Comparer
	if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
		t.Fatal(diff)
	}
	if tc.err != nil {
		return
	}
	if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
		// wantDisasm, _ := tc.want.DisassembleString()
		// gotDisasm, _ := got.DisassembleString()
		// t.Log(cmp.Diff(wantDisasm, gotDisasm, opts...))
		t.Log(got.DisassembleString())
		t.Fatal(diff)
	}
}

// Create a new position in tests
var P = position.New

// Create a new span in tests
var S = position.NewSpan

const testFileName = "main"

// Create a new source location in tests.
// Create a new location in tests
func L(startPos, endPos *position.Position) *position.Location {
	return position.NewLocation(testFileName, startPos, endPos)
}

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
					value.SymbolTable.Add("foo"),
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

func TestBinaryExpressions(t *testing.T) {
	tests := testTable{
		"resolve static add": {
			input: "1i8 + 5i8",
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
					value.Int8(6),
				},
			),
		},
		"add": {
			input: "a := 1i8; a + 5i8",
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
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.Int8(1),
					value.Int8(5),
				},
			),
		},
		"resolve static subtract": {
			input: "151i32 - 25i32 - 5i32",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 1, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.Int32(121),
				},
			),
		},
		"subtract": {
			input: "a := 151i32; a - 25i32 - 5i32",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SUBTRACT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(28, 1, 29)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.Int32(151),
					value.Int32(25),
					value.Int32(5),
				},
			),
		},
		"resolve static multiply": {
			input: "45.5 * 2.5",
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
					value.Float(113.75),
				},
			),
		},
		"multiply": {
			input: "a := 45.5; a * 2.5",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.Float(45.5),
					value.Float(2.5),
				},
			),
		},
		"resolve static divide": {
			input: "45.5 / .5",
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
					value.Float(91),
				},
			),
		},
		"divide": {
			input: "a := 45.5; a / .5",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.Float(45.5),
					value.Float(0.5),
				},
			),
		},
		"resolve static exponentiate": {
			input: "-2 ** 2",
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
					value.SmallInt(-4),
				},
			),
		},
		"exponentiate": {
			input: "a := -2; a ** 2",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SmallInt(-2),
					value.SmallInt(2),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(1, 1, 2)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.SmallInt(-5),
				},
			),
		},
		"negate": {
			input: "a := 5; -a",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.SmallInt(5),
				},
			),
		},
		"bitwise not": {
			input: "~10",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.BITWISE_NOT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.SmallInt(10),
				},
			),
		},
		"resolve static logical not": {
			input: "!10",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(2, 1, 3)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"logical not": {
			input: "a := 10; !a",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.SmallInt(10),
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

func TestLocalVariables(t *testing.T) {
	tests := testTable{
		"declare": {
			input: "var a",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				nil,
			),
		},
		"declare with a type": {
			input: "var a: Int",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				nil,
			),
		},
		"declare and initialise": {
			input: "var a = 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.SmallInt(3),
				},
			),
		},
		"read undeclared": {
			input: "a",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(0, 1, 1)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				nil,
			),
			err: errors.ErrorList{
				errors.NewError(L(P(0, 1, 1), P(0, 1, 1)), "undeclared variable: a"),
			},
		},
		"assign undeclared": {
			input: "a = 3",
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
					value.SmallInt(3),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L(P(0, 1, 1), P(4, 1, 5)), "undeclared variable: a"),
			},
		},
		"assign uninitialised": {
			input: `
				var a
				a = 'foo'
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 3),
				},
				[]value.Value{
					value.String("foo"),
				},
			),
		},
		"assign initialised": {
			input: `
				var a = 'foo'
				a = 'bar'
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
				},
				[]value.Value{
					value.String("foo"),
					value.String("bar"),
				},
			),
		},
		"read uninitialised": {
			input: `
				var a
				a + 2
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 3),
				},
				[]value.Value{
					value.SmallInt(2),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L(P(15, 3, 5), P(15, 3, 5)), "can't access an uninitialised local: a"),
			},
		},
		"read initialised": {
			input: `
				var a = 5
				a + 2
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
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
				},
			),
		},
		"read initialised in child scope": {
			input: `
				var a = 5
				do
					a + 2
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
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
				},
			),
		},
		"shadow in child scope": {
			input: `
				var a = 5
				2 + do
					var a = 10
					a + 12
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
					value.SmallInt(10),
					value.SmallInt(12),
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

func TestLocalValues(t *testing.T) {
	tests := testTable{
		"declare": {
			input: "val a",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				nil,
			),
		},
		"declare with a type": {
			input: "val a: Int",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				nil,
			),
		},
		"declare and initialise": {
			input: "val a = 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.SmallInt(3),
				},
			),
		},
		"declare and initialise 255 variables": {
			input: `
				do
					a0:=1;a1:=1;a2:=1;a3:=1;a4:=1;a5:=1;a6:=1;a7:=1;a8:=1;a9:=1;a10:=1;a11:=1;a12:=1;a13:=1;a14:=1;a15:=1;a16:=1;a17:=1;a18:=1;a19:=1;a20:=1;a21:=1;a22:=1;a23:=1;a24:=1;a25:=1;a26:=1;a27:=1;a28:=1;a29:=1;a30:=1;a31:=1;a32:=1;a33:=1;a34:=1;a35:=1;a36:=1;a37:=1;a38:=1;a39:=1;a40:=1;a41:=1;a42:=1;a43:=1;a44:=1;a45:=1;a46:=1;a47:=1;a48:=1;a49:=1;a50:=1;a51:=1;a52:=1;a53:=1;a54:=1;a55:=1;a56:=1;a57:=1;a58:=1;a59:=1;a60:=1;a61:=1;a62:=1;a63:=1;a64:=1;a65:=1;a66:=1;a67:=1;a68:=1;a69:=1;a70:=1;a71:=1;a72:=1;a73:=1;a74:=1;a75:=1;a76:=1;a77:=1;a78:=1;a79:=1;a80:=1;a81:=1;a82:=1;a83:=1;a84:=1;a85:=1;a86:=1;a87:=1;a88:=1;a89:=1;a90:=1;a91:=1;a92:=1;a93:=1;a94:=1;a95:=1;a96:=1;a97:=1;a98:=1;a99:=1;a100:=1;a101:=1;a102:=1;a103:=1;a104:=1;a105:=1;a106:=1;a107:=1;a108:=1;a109:=1;a110:=1;a111:=1;a112:=1;a113:=1;a114:=1;a115:=1;a116:=1;a117:=1;a118:=1;a119:=1;a120:=1;a121:=1;a122:=1;a123:=1;a124:=1;a125:=1;a126:=1;a127:=1;a128:=1;a129:=1;a130:=1;a131:=1;a132:=1;a133:=1;a134:=1;a135:=1;a136:=1;a137:=1;a138:=1;a139:=1;a140:=1;a141:=1;a142:=1;a143:=1;a144:=1;a145:=1;a146:=1;a147:=1;a148:=1;a149:=1;a150:=1;a151:=1;a152:=1;a153:=1;a154:=1;a155:=1;a156:=1;a157:=1;a158:=1;a159:=1;a160:=1;a161:=1;a162:=1;a163:=1;a164:=1;a165:=1;a166:=1;a167:=1;a168:=1;a169:=1;a170:=1;a171:=1;a172:=1;a173:=1;a174:=1;a175:=1;a176:=1;a177:=1;a178:=1;a179:=1;a180:=1;a181:=1;a182:=1;a183:=1;a184:=1;a185:=1;a186:=1;a187:=1;a188:=1;a189:=1;a190:=1;a191:=1;a192:=1;a193:=1;a194:=1;a195:=1;a196:=1;a197:=1;a198:=1;a199:=1;a200:=1;a201:=1;a202:=1;a203:=1;a204:=1;a205:=1;a206:=1;a207:=1;a208:=1;a209:=1;a210:=1;a211:=1;a212:=1;a213:=1;a214:=1;a215:=1;a216:=1;a217:=1;a218:=1;a219:=1;a220:=1;a221:=1;a222:=1;a223:=1;a224:=1;a225:=1;a226:=1;a227:=1;a228:=1;a229:=1;a230:=1;a231:=1;a232:=1;a233:=1;a234:=1;a235:=1;a236:=1;a237:=1;a238:=1;a239:=1;a240:=1;a241:=1;a242:=1;a243:=1;a244:=1;a245:=1;a246:=1;a247:=1;a248:=1;a249:=1;a250:=1;a251:=1;a252:=1;a253:=1;a254:=1;a255:=1
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				append(
					append(
						[]byte{
							byte(bytecode.PREP_LOCALS16),
							1,
							0,
						},
						declareNVariables(253)...,
					),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL16), 1, 0,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL16), 1, 1,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL16), 1, 2,
					byte(bytecode.LEAVE_SCOPE32), 1, 2, 1, 0,
					byte(bytecode.RETURN),
				),
				L(P(0, 1, 1), P(1958, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 768),
					bytecode.NewLineInfo(4, 2),
				},
				[]value.Value{
					value.SmallInt(1),
				},
			),
		},
		"assign uninitialised": {
			input: `
				val a
				a = 'foo'
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 3),
				},
				[]value.Value{
					value.String("foo"),
				},
			),
		},
		"assign initialised": {
			input: `
				val a = 'foo'
				a = 'bar'
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(32, 3, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
				},
				[]value.Value{
					value.String("foo"),
					value.String("bar"),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L(P(23, 3, 5), P(31, 3, 13)), "can't reassign a val: a"),
			},
		},
		"read uninitialised": {
			input: `
				val a
				a + 2
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(20, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 3),
				},
				[]value.Value{
					value.SmallInt(2),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L(P(15, 3, 5), P(15, 3, 5)), "can't access an uninitialised local: a"),
			},
		},
		"read initialised": {
			input: `
				val a = 5
				a + 2
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
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 3, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 4),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
				},
			),
		},
		"read initialised in child scope": {
			input: `
				val a = 5
				do
					a + 2
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
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
				},
			),
		},
		"shadow in child scope": {
			input: `
				val a = 5
				2 + do
					val a = 10
					a + 12
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 6, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(6, 1),
				},
				[]value.Value{
					value.SmallInt(5),
					value.SmallInt(2),
					value.SmallInt(10),
					value.SmallInt(12),
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

func declareNVariables(n int) []byte {
	b := make([]byte, 0, n*4)
	for i := 0; i < n; i++ {
		b = append(
			b,
			byte(bytecode.LOAD_VALUE8), 0,
			byte(bytecode.SET_LOCAL8), byte(i+3),
			byte(bytecode.POP),
		)
	}

	return b
}

func TestComplexAssignment(t *testing.T) {
	tests := testTable{
		"add": {
			input: "a := 1; a += 3",
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
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"subtract": {
			input: "a := 1; a -= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"multiply": {
			input: "a := 1; a *= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"divide": {
			input: "a := 1; a /= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DIVIDE),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"exponentiate": {
			input: "a := 1; a **= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"modulo": {
			input: "a := 1; a %= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.MODULO),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"bitwise AND": {
			input: "a := 1; a &= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"bitwise OR": {
			input: "a := 1; a |= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"bitwise XOR": {
			input: "a := 1; a |= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"left bitshift": {
			input: "a := 1; a <<= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"left logical bitshift": {
			input: "a := 1; a <<<= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"right bitshift": {
			input: "a := 1; a >>= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
					value.SmallInt(3),
				},
			),
		},
		"right logical bitshift": {
			input: "a := 1; a >>>= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 9),
				},
				[]value.Value{
					value.SmallInt(1),
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

func TestBitwiseAnd(t *testing.T) {
	tests := testTable{
		"resolve static AND": {
			input: "23 & 10",
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
					value.SmallInt(2),
				},
			),
		},
		"resolve static nested AND": {
			input: "23 & 15 & 46",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.SmallInt(6),
				},
			),
		},
		"compile runtime AND": {
			input: "a := 23; a & 15 & 46",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.SmallInt(23),
					value.SmallInt(15),
					value.SmallInt(46),
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
					value.SmallInt(31),
				},
			),
		},
		"resolve static nested OR": {
			input: "23 | 15 | 46",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.SmallInt(63),
				},
			),
		},
		"compile runtime OR": {
			input: "a := 23; a | 15 | 46",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.SmallInt(23),
					value.SmallInt(15),
					value.SmallInt(46),
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
					value.SmallInt(29),
				},
			),
		},
		"resolve static nested XOR": {
			input: "23 ^ 15 ^ 46",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(11, 1, 12)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.SmallInt(54),
				},
			),
		},
		"compile runtime XOR": {
			input: "a := 23; a ^ 15 ^ 46",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.SmallInt(23),
					value.SmallInt(15),
					value.SmallInt(46),
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
					value.SmallInt(3),
				},
			),
		},
		"resolve static nested modulo": {
			input: "24 % 15 % 2",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				[]value.Value{
					value.SmallInt(1),
				},
			),
		},
		"compile runtime modulo": {
			input: "a := 24; a % 15 % 46",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.MODULO),
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.MODULO),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 1, 20)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				[]value.Value{
					value.SmallInt(24),
					value.SmallInt(15),
					value.SmallInt(46),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25 == '25'": {
			input: "25 == '25'",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"compile runtime 24 == 98": {
			input: "a := 24; a == 98",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SmallInt(24),
					value.SmallInt(98),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25 != '25'": {
			input: "25 != '25'",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"compile runtime 24 != 98": {
			input: "a := 24; a != 98",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.NOT_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SmallInt(24),
					value.SmallInt(98),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25 === 25.0": {
			input: "25 === 25.0",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25 === '25'": {
			input: "25 === '25'",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"compile runtime 24 === 98": {
			input: "a := 24; a === 98",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.STRICT_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SmallInt(24),
					value.SmallInt(98),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25 !== 25.0": {
			input: "25 !== 25.0",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25 !== '25'": {
			input: "25 !== '25'",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"compile runtime 24 !== 98": {
			input: "a := 24; a !== 98",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.STRICT_NOT_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SmallInt(24),
					value.SmallInt(98),
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
		"resolve static 25 > 3": {
			input: "25 > 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25.2 > 25": {
			input: "25.2 > 25",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 7 > 20": {
			input: "7 > 20",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"compile runtime 24 > 98": {
			input: "a := 24; a > 98",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SmallInt(24),
					value.SmallInt(98),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25 >= 3": {
			input: "25 >= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25.2 >= 25": {
			input: "25.2 >= 25",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 7 >= 20": {
			input: "7 >= 20",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"compile runtime 24 >= 98": {
			input: "a := 24; a >= 98",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.GREATER_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SmallInt(24),
					value.SmallInt(98),
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
		"resolve static 25 < 3": {
			input: "25 < 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25.2 < 25": {
			input: "25.2 < 25",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(8, 1, 9)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 7 < 20": {
			input: "7 < 20",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"compile runtime 24 < 98": {
			input: "a := 24; a < 98",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SmallInt(24),
					value.SmallInt(98),
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
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(5, 1, 6)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25 <= 3": {
			input: "25 <= 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 25.2 <= 25": {
			input: "25.2 <= 25",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"resolve static 7 <= 20": {
			input: "7 <= 20",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(6, 1, 7)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				nil,
			),
		},
		"compile runtime 24 <= 98": {
			input: "a := 24; a <= 98",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS_EQUAL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SmallInt(24),
					value.SmallInt(98),
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
					byte(bytecode.LOOP), 0, 3,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(19, 3, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 3),
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
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 11,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 4),
				},
				[]value.Value{
					value.SmallInt(0),
					value.SmallInt(1),
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
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 11,
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(37, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 5),
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
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 12,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 20,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 3, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(43, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(4, 6),
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
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.LESS),
					byte(bytecode.JUMP_UNLESS), 0, 19,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 3,
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 3,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 4,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 4,
					byte(bytecode.LOOP), 0, 27,
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.LEAVE_SCOPE16), 4, 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(61, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(3, 4),
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 1),
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

func TestGetModuleConstant(t *testing.T) {
	tests := testTable{
		"absolute path ::Std": {
			input: "::Std",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(4, 1, 5)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				[]value.Value{
					value.SymbolTable.Add("Std"),
				},
			),
		},
		"absolute nested path ::Std::Float::INF": {
			input: "::Std::Float::INF",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.GET_MOD_CONST8), 2,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.SymbolTable.Add("Std"),
					value.SymbolTable.Add("Float"),
					value.SymbolTable.Add("INF"),
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

func TestDefModuleConstant(t *testing.T) {
	tests := testTable{
		"relative path Foo": {
			input: "Foo := 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.DEF_MOD_CONST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(7, 1, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.SmallInt(3),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"absolute path ::Foo": {
			input: "::Foo := 3",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.DEF_MOD_CONST8), 1,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				[]value.Value{
					value.SmallInt(3),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"absolute nested path ::Std::Float::Foo": {
			input: "::Std::Float::Foo := 'bar'",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.GET_MOD_CONST8), 2,
					byte(bytecode.DEF_MOD_CONST8), 3,
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 1, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.String("bar"),
					value.SymbolTable.Add("Std"),
					value.SymbolTable.Add("Float"),
					value.SymbolTable.Add("Foo"),
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

func TestDefClass(t *testing.T) {
	tests := testTable{
		"anonymous class without a body": {
			input: "class; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				nil,
			),
		},
		"class with a relative name without a body": {
			input: "class Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"anonymous class with an absolute parent": {
			input: "class < ::Bar; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.DEF_ANON_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(17, 1, 18)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.SymbolTable.Add("Bar"),
				},
			),
		},
		"class with an absolute parent": {
			input: "class Foo < ::Bar; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(21, 1, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
					value.SymbolTable.Add("Bar"),
				},
			),
		},
		"anonymous class with a nested parent": {
			input: "class < ::Baz::Bar; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.DEF_ANON_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 1, 23)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.SymbolTable.Add("Baz"),
					value.SymbolTable.Add("Bar"),
				},
			),
		},
		"class with an absolute nested parent": {
			input: "class Foo < ::Baz::Bar; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.GET_MOD_CONST8), 2,
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(26, 1, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
					value.SymbolTable.Add("Baz"),
					value.SymbolTable.Add("Bar"),
				},
			),
		},
		"class with an absolute name without a body": {
			input: "class ::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"class with an absolute nested name without a body": {
			input: "class ::Std::Int::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 1, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				[]value.Value{
					value.SymbolTable.Add("Std"),
					value.SymbolTable.Add("Int"),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"class with a body": {
			input: `
				class Foo
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(44, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"class with a an error": {
			input: "class A then def a then a",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(24, 1, 25)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(0, 1, 1), P(24, 1, 25)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 6),
						},
						[]value.Value{
							vm.NewBytecodeMethodNoParams(
								value.SymbolTable.Add("a"),
								[]byte{
									byte(bytecode.RETURN),
								},
								L(P(13, 1, 14), P(24, 1, 25)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(1, 1),
								},
								nil,
							),
							value.SymbolTable.Add("a"),
						},
					),
					value.SymbolTable.Add("A"),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L(P(24, 1, 25), P(24, 1, 25)), "undeclared variable: a"),
			},
		},
		"anonymous class with a body": {
			input: `
				class
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(40, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
				},
			),
		},
		"nested classes": {
			input: `
				class Foo
					class Bar
						a := 1
						a + 2
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.UNDEFINED),
							byte(bytecode.DEF_CLASS),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(70, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 5),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							vm.NewBytecodeMethodNoParams(
								classSymbol,
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.SET_LOCAL8), 3,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL8), 3,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.ADD),
									byte(bytecode.POP),
									byte(bytecode.SELF),
									byte(bytecode.RETURN),
								},
								L(P(20, 3, 6), P(62, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 3),
								},
								[]value.Value{
									value.SmallInt(1),
									value.SmallInt(2),
								},
							),
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
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

func TestDefModule(t *testing.T) {
	tests := testTable{
		"anonymous module without a body": {
			input: "module; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(10, 1, 11)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				nil,
			),
		},
		"module with a relative name without a body": {
			input: "module Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(14, 1, 15)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"class with an absolute name without a body": {
			input: "module ::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(16, 1, 17)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"class with an absolute nested name without a body": {
			input: "module ::Std::Int::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(26, 1, 27)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.SymbolTable.Add("Std"),
					value.SymbolTable.Add("Int"),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"anonymous module with a body": {
			input: `
				module
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_ANON_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(42, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(41, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
				},
			),
		},
		"module with a body": {
			input: `
				module Foo
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(46, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(45, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"nested modules": {
			input: `
				module Foo
					module Bar
						a := 1
						a + 2
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(73, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_MODULE),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(72, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							vm.NewBytecodeMethodNoParams(
								moduleSymbol,
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.SET_LOCAL8), 3,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL8), 3,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.ADD),
									byte(bytecode.POP),
									byte(bytecode.SELF),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(64, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 3),
								},
								[]value.Value{
									value.SmallInt(1),
									value.SmallInt(2),
								},
							),
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
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

func TestDefMethod(t *testing.T) {
	tests := testTable{
		"define method in top level": {
			input: `
				def foo then :bar
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(22, 2, 22)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						value.SymbolTable.Add("foo"),
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(21, 2, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(2, 2),
						},
						[]value.Value{
							value.SymbolTable.Add("bar"),
						},
					),
					value.SymbolTable.Add("foo"),
				},
			),
		},
		"define method with required parameters in top level": {
			input: `
				def foo(a, b)
					c := 5
					a + b + c
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(53, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethod(
						value.SymbolTable.Add("foo"),
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 1,
							byte(bytecode.GET_LOCAL8), 2,
							byte(bytecode.ADD),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.ADD),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(52, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 5),
							bytecode.NewLineInfo(5, 1),
						},
						[]value.Symbol{
							value.SymbolTable.Add("a"),
							value.SymbolTable.Add("b"),
						},
						0,
						-1,
						false,
						false,
						[]value.Value{
							value.SmallInt(5),
						},
					),
					value.SymbolTable.Add("foo"),
				},
			),
		},
		"define method with optional parameters in top level": {
			input: `
				def foo(a, b = 5.2, c = 10)
					d := 5
					a + b + c + d
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_METHOD),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethod(
						value.SymbolTable.Add("foo"),
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,

							byte(bytecode.GET_LOCAL8), 2,
							byte(bytecode.JUMP_UNLESS_UNDEF), 0, 5,
							byte(bytecode.POP),
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 2,
							byte(bytecode.POP),

							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.JUMP_UNLESS_UNDEF), 0, 5,
							byte(bytecode.POP),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),

							byte(bytecode.LOAD_VALUE8), 2,
							byte(bytecode.SET_LOCAL8), 4,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 1,
							byte(bytecode.GET_LOCAL8), 2,
							byte(bytecode.ADD),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.ADD),
							byte(bytecode.GET_LOCAL8), 4,
							byte(bytecode.ADD),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(70, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(2, 13),
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(4, 7),
							bytecode.NewLineInfo(5, 1),
						},
						[]value.Symbol{
							value.SymbolTable.Add("a"),
							value.SymbolTable.Add("b"),
							value.SymbolTable.Add("c"),
						},
						2,
						-1,
						false,
						false,
						[]value.Value{
							value.Float(5.2),
							value.SmallInt(10),
							value.SmallInt(5),
						},
					),
					value.SymbolTable.Add("foo"),
				},
			),
		},
		"define method with required parameters in a class": {
			input: `
				class Bar
					def foo(a, b)
						c := 5
						a + b + c
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(79, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(78, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							vm.NewBytecodeMethod(
								value.SymbolTable.Add("foo"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.SET_LOCAL8), 3,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.GET_LOCAL8), 2,
									byte(bytecode.ADD),
									byte(bytecode.GET_LOCAL8), 3,
									byte(bytecode.ADD),
									byte(bytecode.RETURN),
								},
								L(P(20, 3, 6), P(70, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 5),
									bytecode.NewLineInfo(6, 1),
								},
								[]value.Symbol{
									value.SymbolTable.Add("a"),
									value.SymbolTable.Add("b"),
								},
								0,
								-1,
								false,
								false,
								[]value.Value{
									value.SmallInt(5),
								},
							),
							value.SymbolTable.Add("foo"),
						},
					),
					value.SymbolTable.Add("Bar"),
				},
			),
		},
		"define method with required parameters in a module": {
			input: `
				module Bar
					def foo(a, b)
						c := 5
						a + b + c
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(80, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_METHOD),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(79, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 3),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							vm.NewBytecodeMethod(
								value.SymbolTable.Add("foo"),
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.SET_LOCAL8), 3,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL8), 1,
									byte(bytecode.GET_LOCAL8), 2,
									byte(bytecode.ADD),
									byte(bytecode.GET_LOCAL8), 3,
									byte(bytecode.ADD),
									byte(bytecode.RETURN),
								},
								L(P(21, 3, 6), P(71, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 5),
									bytecode.NewLineInfo(6, 1),
								},
								[]value.Symbol{
									value.SymbolTable.Add("a"),
									value.SymbolTable.Add("b"),
								},
								0,
								-1,
								false,
								false,
								[]value.Value{
									value.SmallInt(5),
								},
							),
							value.SymbolTable.Add("foo"),
						},
					),
					value.SymbolTable.Add("Bar"),
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
					value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 0, nil),
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
					value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 2, nil),
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
					value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 2, nil),
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
					value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 2, []value.Symbol{value.ToSymbol("b")}),
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
					value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 0, nil),
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
					value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 2, nil),
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
					value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 2, []value.Symbol{value.ToSymbol("b")}),
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

func TestDefMixin(t *testing.T) {
	tests := testTable{
		"anonymous mixin without a body": {
			input: "mixin; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(9, 1, 10)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				nil,
			),
		},
		"mixin with a relative name without a body": {
			input: "mixin Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(13, 1, 14)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"mixin with an absolute name without a body": {
			input: "mixin ::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(15, 1, 16)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				[]value.Value{
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"mixin with an absolute nested name without a body": {
			input: "mixin ::Std::Int::Foo; end",
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.LOAD_VALUE8), 2,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(25, 1, 26)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				[]value.Value{
					value.SymbolTable.Add("Std"),
					value.SymbolTable.Add("Int"),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"anonymous mixin with a body": {
			input: `
				mixin
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_ANON_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 2),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(40, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
				},
			),
		},
		"mixin with a body": {
			input: `
				mixin Foo
					a := 1
					a + 2
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(45, 5, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(5, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.SET_LOCAL8), 3,
							byte(bytecode.POP),
							byte(bytecode.GET_LOCAL8), 3,
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.ADD),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(44, 5, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(4, 3),
							bytecode.NewLineInfo(5, 3),
						},
						[]value.Value{
							value.SmallInt(1),
							value.SmallInt(2),
						},
					),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"nested mixins": {
			input: `
				mixin Foo
					mixin Bar
						a := 1
						a + 2
					end
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(71, 7, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(7, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.LOAD_VALUE8), 1,
							byte(bytecode.DEF_MIXIN),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(70, 7, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 4),
							bytecode.NewLineInfo(7, 3),
						},
						[]value.Value{
							vm.NewBytecodeMethodNoParams(
								mixinSymbol,
								[]byte{
									byte(bytecode.PREP_LOCALS8), 1,
									byte(bytecode.LOAD_VALUE8), 0,
									byte(bytecode.SET_LOCAL8), 3,
									byte(bytecode.POP),
									byte(bytecode.GET_LOCAL8), 3,
									byte(bytecode.LOAD_VALUE8), 1,
									byte(bytecode.ADD),
									byte(bytecode.POP),
									byte(bytecode.SELF),
									byte(bytecode.RETURN),
								},
								L(P(20, 3, 6), P(62, 6, 8)),
								bytecode.LineInfoList{
									bytecode.NewLineInfo(4, 4),
									bytecode.NewLineInfo(5, 3),
									bytecode.NewLineInfo(6, 3),
								},
								[]value.Value{
									value.SmallInt(1),
									value.SmallInt(2),
								},
							),
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
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

func TestInclude(t *testing.T) {
	tests := testTable{
		"include a global constant in a class": {
			input: `
				class Foo
					include ::Bar
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(40, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 5),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"include two constants in a class": {
			input: `
				class Foo
					include ::Bar, ::Baz
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(48, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.INCLUDE),

							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 1,
							byte(bytecode.SELF),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),

							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(47, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 9),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
							value.SymbolTable.Add("Baz"),
						},
					),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"include in top level": {
			input: `include ::Bar`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(0, 1, 1), P(12, 1, 13)),
					"can't include mixins in the top level",
				),
			},
		},
		"include in a module": {
			input: `
				module Foo
					include ::Bar
				end
			`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(21, 3, 6), P(33, 3, 18)),
					"can't include mixins in a module",
				),
			},
		},
		"include in a method": {
			input: `
				def foo
					include ::Bar
				end
			`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(18, 3, 6), P(30, 3, 18)),
					"can't include mixins in a method",
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

func TestExtend(t *testing.T) {
	tests := testTable{
		"extend a class with a global constant": {
			input: `
				class Foo
					extend ::Bar
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(39, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 6),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"extend a module with a global constant": {
			input: `
				module Foo
					extend ::Bar
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(41, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						moduleSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(40, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 6),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"extend a mixin with a global constant": {
			input: `
				mixin Foo
					extend ::Bar
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(40, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						mixinSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(39, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 6),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
						},
					),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"extend a class with two global constants": {
			input: `
				class Foo
					extend ::Bar, ::Baz
				end
			`,
			want: vm.NewBytecodeMethodNoParams(
				mainSymbol,
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				L(P(0, 1, 1), P(47, 4, 8)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
					bytecode.NewLineInfo(4, 1),
				},
				[]value.Value{
					vm.NewBytecodeMethodNoParams(
						classSymbol,
						[]byte{
							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 0,
							byte(bytecode.SELF),
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.INCLUDE),

							byte(bytecode.ROOT),
							byte(bytecode.GET_MOD_CONST8), 1,
							byte(bytecode.SELF),
							byte(bytecode.GET_SINGLETON),
							byte(bytecode.INCLUDE),
							byte(bytecode.NIL),
							byte(bytecode.POP),

							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(46, 4, 7)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(3, 11),
							bytecode.NewLineInfo(4, 3),
						},
						[]value.Value{
							value.SymbolTable.Add("Bar"),
							value.SymbolTable.Add("Baz"),
						},
					),
					value.SymbolTable.Add("Foo"),
				},
			),
		},
		"extend in top level": {
			input: `extend ::Bar`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(0, 1, 1), P(11, 1, 12)),
					"can't extend mixins in the top level",
				),
			},
		},
		"extend in a method": {
			input: `
				def foo
					extend ::Bar
				end
			`,
			err: errors.ErrorList{
				errors.NewError(
					L(P(18, 3, 6), P(29, 3, 17)),
					"can't extend mixins in a method",
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
