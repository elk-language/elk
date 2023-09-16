package compiler

import (
	"math"
	"math/big"
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/object"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/google/go-cmp/cmp"
)

// Represents a single compiler test case.
type testCase struct {
	input string
	want  *bytecode.Chunk
	err   errors.ErrorList
}

// Type of the compiler test table.
type testTable map[string]testCase

func compilerTest(tc testCase, t *testing.T) {
	t.Helper()

	got, err := CompileSource("main", tc.input)
	opts := []cmp.Option{
		cmp.AllowUnexported(object.BigInt{}),
	}
	if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
		t.Log(got.DisassembleString())
		t.Fatal(diff)
	}
	if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.UInt8(1),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(2, 1, 3)),
			},
		},
		"put UInt16": {
			input: "25u16",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.UInt16(25),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
		},
		"put UInt32": {
			input: "450_200u32",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.UInt32(450200),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"put UInt64": {
			input: "450_200u64",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.UInt64(450200),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"put Int8": {
			input: "1i8",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Int8(1),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(2, 1, 3)),
			},
		},
		"put Int16": {
			input: "25i16",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Int16(25),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
		},
		"put Int32": {
			input: "450_200i32",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Int32(450200),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"put Int64": {
			input: "450_200i64",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Int64(450200),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"put SmallInt": {
			input: "450_200",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(450200),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
			},
		},
		"put BigInt": {
			input: (&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5)).String(),
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5))),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(
					P(0, 1, 1),
					P(
						len((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5)).String())-1,
						1,
						len((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5)).String()),
					),
				),
			},
		},
		"put Float64": {
			input: "45.5f64",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Float64(45.5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
			},
		},
		"put Float32": {
			input: "45.5f32",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Float32(45.5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
			},
		},
		"put Float": {
			input: "45.5",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Float(45.5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(3, 1, 4)),
			},
		},
		"put Raw String": {
			input: `'foo\n'`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String(`foo\n`),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
			},
		},
		"put String": {
			input: `"foo\n"`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String("foo\n"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
			},
		},
		"put raw Char": {
			input: `c'I'`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Char('I'),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(3, 1, 4)),
			},
		},
		"put Char": {
			input: `c"\n"`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Char('\n'),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
		},
		"put nil": {
			input: `nil`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(2, 1, 3)),
			},
		},
		"put true": {
			input: `true`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(3, 1, 4)),
			},
		},
		"put false": {
			input: `false`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
		},
		"put simple Symbol": {
			input: `:foo`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SymbolTable.Add("foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(3, 1, 4)),
			},
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Int8(6),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(8, 1, 9)),
			},
		},
		"add": {
			input: "a := 1i8; a + 5i8",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Int8(1),
					object.Int8(5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				Location: L(P(0, 1, 1), P(16, 1, 17)),
			},
		},
		"resolve static subtract": {
			input: "151i32 - 25i32 - 5i32",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Int32(121),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(20, 1, 21)),
			},
		},
		"subtract": {
			input: "a := 151i32; a - 25i32 - 5i32",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.CONSTANT8), 2,
					byte(bytecode.SUBTRACT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Int32(151),
					object.Int32(25),
					object.Int32(5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 10),
				},
				Location: L(P(0, 1, 1), P(28, 1, 29)),
			},
		},
		"resolve static multiply": {
			input: "45.5 * 2.5",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Float(113.75),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"multiply": {
			input: "a := 45.5; a * 2.5",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Float(45.5),
					object.Float(2.5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				Location: L(P(0, 1, 1), P(17, 1, 18)),
			},
		},
		"resolve static divide": {
			input: "45.5 / .5",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Float(91),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(8, 1, 9)),
			},
		},
		"divide": {
			input: "a := 45.5; a / .5",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.Float(45.5),
					object.Float(0.5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				Location: L(P(0, 1, 1), P(16, 1, 17)),
			},
		},
		"resolve static exponentiate": {
			input: "-2 ** 2",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(-4),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
			},
		},
		"exponentiate": {
			input: "a := -2; a ** 2",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(-2),
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 8),
				},
				Location: L(P(0, 1, 1), P(14, 1, 15)),
			},
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(-5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(1, 1, 2)),
			},
		},
		"negate": {
			input: "a := 5; -a",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"bitwise not": {
			input: "~10",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.BITWISE_NOT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(10),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(2, 1, 3)),
			},
		},
		"resolve static logical not": {
			input: "!10",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(2, 1, 3)),
			},
		},
		"logical not": {
			input: "a := 10; !a",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(10),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 7),
				},
				Location: L(P(0, 1, 1), P(10, 1, 11)),
			},
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
		},
		"declare with a type": {
			input: "var a: Int",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"declare and initialise": {
			input: "var a = 3",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(3),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				Location: L(P(0, 1, 1), P(8, 1, 9)),
			},
		},
		"read undeclared": {
			input: "a",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(0, 1, 1)),
			},
			err: errors.ErrorList{
				errors.NewError(L(P(0, 1, 1), P(0, 1, 1)), "undeclared variable: a"),
			},
		},
		"assign undeclared": {
			input: "a = 3",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(3),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
			err: errors.ErrorList{
				errors.NewError(L(P(0, 1, 1), P(4, 1, 5)), "undeclared variable: a"),
			},
		},
		"assign uninitialised": {
			input: `
				var a
				a = 'foo'
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String("foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 3),
				},
				Location: L(P(0, 1, 1), P(24, 3, 14)),
			},
		},
		"assign initialised": {
			input: `
				var a = 'foo'
				a = 'bar'
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String("foo"),
					object.String("bar"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
				},
				Location: L(P(0, 1, 1), P(32, 3, 14)),
			},
		},
		"read uninitialised": {
			input: `
				var a
				a + 2
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 3),
				},
				Location: L(P(0, 1, 1), P(20, 3, 10)),
			},
			err: errors.ErrorList{
				errors.NewError(L(P(15, 3, 5), P(15, 3, 5)), "can't access an uninitialised local: a"),
			},
		},
		"read initialised": {
			input: `
				var a = 5
				a + 2
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 4),
				},
				Location: L(P(0, 1, 1), P(24, 3, 10)),
			},
		},
		"read initialised in child scope": {
			input: `
				var a = 5
				do
					a + 2
				end
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 1),
				},
				Location: L(P(0, 1, 1), P(40, 5, 8)),
			},
		},
		"shadow in child scope": {
			input: `
				var a = 5
				2 + do
					var a = 10
					a + 12
				end
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.CONSTANT8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CONSTANT8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
					object.SmallInt(10),
					object.SmallInt(12),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(6, 1),
				},
				Location: L(P(0, 1, 1), P(61, 6, 8)),
			},
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
		},
		"declare with a type": {
			input: "val a: Int",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"declare and initialise": {
			input: "val a = 3",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(3),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				Location: L(P(0, 1, 1), P(8, 1, 9)),
			},
		},
		"declare and initialise 257 variables": {
			input: `
				do
					a0:=1;a1:=1;a2:=1;a3:=1;a4:=1;a5:=1;a6:=1;a7:=1;a8:=1;a9:=1;a10:=1;a11:=1;a12:=1;a13:=1;a14:=1;a15:=1;a16:=1;a17:=1;a18:=1;a19:=1;a20:=1;a21:=1;a22:=1;a23:=1;a24:=1;a25:=1;a26:=1;a27:=1;a28:=1;a29:=1;a30:=1;a31:=1;a32:=1;a33:=1;a34:=1;a35:=1;a36:=1;a37:=1;a38:=1;a39:=1;a40:=1;a41:=1;a42:=1;a43:=1;a44:=1;a45:=1;a46:=1;a47:=1;a48:=1;a49:=1;a50:=1;a51:=1;a52:=1;a53:=1;a54:=1;a55:=1;a56:=1;a57:=1;a58:=1;a59:=1;a60:=1;a61:=1;a62:=1;a63:=1;a64:=1;a65:=1;a66:=1;a67:=1;a68:=1;a69:=1;a70:=1;a71:=1;a72:=1;a73:=1;a74:=1;a75:=1;a76:=1;a77:=1;a78:=1;a79:=1;a80:=1;a81:=1;a82:=1;a83:=1;a84:=1;a85:=1;a86:=1;a87:=1;a88:=1;a89:=1;a90:=1;a91:=1;a92:=1;a93:=1;a94:=1;a95:=1;a96:=1;a97:=1;a98:=1;a99:=1;a100:=1;a101:=1;a102:=1;a103:=1;a104:=1;a105:=1;a106:=1;a107:=1;a108:=1;a109:=1;a110:=1;a111:=1;a112:=1;a113:=1;a114:=1;a115:=1;a116:=1;a117:=1;a118:=1;a119:=1;a120:=1;a121:=1;a122:=1;a123:=1;a124:=1;a125:=1;a126:=1;a127:=1;a128:=1;a129:=1;a130:=1;a131:=1;a132:=1;a133:=1;a134:=1;a135:=1;a136:=1;a137:=1;a138:=1;a139:=1;a140:=1;a141:=1;a142:=1;a143:=1;a144:=1;a145:=1;a146:=1;a147:=1;a148:=1;a149:=1;a150:=1;a151:=1;a152:=1;a153:=1;a154:=1;a155:=1;a156:=1;a157:=1;a158:=1;a159:=1;a160:=1;a161:=1;a162:=1;a163:=1;a164:=1;a165:=1;a166:=1;a167:=1;a168:=1;a169:=1;a170:=1;a171:=1;a172:=1;a173:=1;a174:=1;a175:=1;a176:=1;a177:=1;a178:=1;a179:=1;a180:=1;a181:=1;a182:=1;a183:=1;a184:=1;a185:=1;a186:=1;a187:=1;a188:=1;a189:=1;a190:=1;a191:=1;a192:=1;a193:=1;a194:=1;a195:=1;a196:=1;a197:=1;a198:=1;a199:=1;a200:=1;a201:=1;a202:=1;a203:=1;a204:=1;a205:=1;a206:=1;a207:=1;a208:=1;a209:=1;a210:=1;a211:=1;a212:=1;a213:=1;a214:=1;a215:=1;a216:=1;a217:=1;a218:=1;a219:=1;a220:=1;a221:=1;a222:=1;a223:=1;a224:=1;a225:=1;a226:=1;a227:=1;a228:=1;a229:=1;a230:=1;a231:=1;a232:=1;a233:=1;a234:=1;a235:=1;a236:=1;a237:=1;a238:=1;a239:=1;a240:=1;a241:=1;a242:=1;a243:=1;a244:=1;a245:=1;a246:=1;a247:=1;a248:=1;a249:=1;a250:=1;a251:=1;a252:=1;a253:=1;a254:=1;a255:=1;a256:=1
				end
			`,
			want: &bytecode.Chunk{
				Instructions: append(
					append(
						[]byte{
							byte(bytecode.PREP_LOCALS16),
							1,
							1,
						},
						declareNVariables(256)...,
					),
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL16), 1, 0,
					byte(bytecode.LEAVE_SCOPE32), 1, 0, 1, 1,
					byte(bytecode.RETURN),
				),
				Constants: []object.Value{
					object.SmallInt(1),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 771),
					bytecode.NewLineInfo(4, 2),
				},
				Location: L(P(0, 1, 1), P(1966, 4, 8)),
			},
		},
		"assign uninitialised": {
			input: `
				val a
				a = 'foo'
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String("foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 3),
				},
				Location: L(P(0, 1, 1), P(24, 3, 14)),
			},
		},
		"assign initialised": {
			input: `
				val a = 'foo'
				a = 'bar'
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String("foo"),
					object.String("bar"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
				},
				Location: L(P(0, 1, 1), P(32, 3, 14)),
			},
			err: errors.ErrorList{
				errors.NewError(L(P(23, 3, 5), P(31, 3, 13)), "can't reassign a val: a"),
			},
		},
		"read uninitialised": {
			input: `
				val a
				a + 2
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 3),
				},
				Location: L(P(0, 1, 1), P(20, 3, 10)),
			},
			err: errors.ErrorList{
				errors.NewError(L(P(15, 3, 5), P(15, 3, 5)), "can't access an uninitialised local: a"),
			},
		},
		"read initialised": {
			input: `
				val a = 5
				a + 2
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 4),
				},
				Location: L(P(0, 1, 1), P(24, 3, 10)),
			},
		},
		"read initialised in child scope": {
			input: `
				val a = 5
				do
					a + 2
				end
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 1),
				},
				Location: L(P(0, 1, 1), P(40, 5, 8)),
			},
		},
		"shadow in child scope": {
			input: `
				val a = 5
				2 + do
					val a = 10
					a + 12
				end
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.CONSTANT8), 2,
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.CONSTANT8), 3,
					byte(bytecode.ADD),
					byte(bytecode.LEAVE_SCOPE16), 1, 1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
					object.SmallInt(10),
					object.SmallInt(12),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(4, 3),
					bytecode.NewLineInfo(5, 3),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(6, 1),
				},
				Location: L(P(0, 1, 1), P(61, 6, 8)),
			},
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
			byte(bytecode.CONSTANT8), 0,
			byte(bytecode.SET_LOCAL8), byte(i),
			byte(bytecode.POP),
		)
	}

	return b
}

func TestIfExpression(t *testing.T) {
	tests := testTable{
		"resolve static condition with empty then and else": {
			input: "if true; end",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(11, 1, 12)),
			},
		},
		"empty then and else": {
			input: "a := true; if a; end",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.TRUE),
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
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
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 12),
				},
				Location: L(P(0, 1, 1), P(19, 1, 20)),
			},
		},
		"resolve static condition with then branch": {
			input: `
				if true
					10
				end
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(10),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(4, 1),
				},
				Location: L(P(0, 1, 1), P(28, 4, 8)),
			},
		},
		"resolve static condition with then branch to nil": {
			input: `
				if false
					10
				end
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(4, 2),
				},
				Location: L(P(0, 1, 1), P(29, 4, 8)),
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(6, 1),
				},
				Location: L(P(0, 1, 1), P(45, 6, 8)),
			},
		},
		"with then branch": {
			input: `
				a := 5
				if a
					a = a * 2
				end
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.JUMP_UNLESS), 0, 11,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.JUMP), 0, 2,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(5, 1),
				},
				Location: L(P(0, 1, 1), P(43, 5, 8)),
			},
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.JUMP_UNLESS), 0, 11,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.JUMP), 0, 5,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 2,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
					object.SmallInt(30),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(6, 2),
					bytecode.NewLineInfo(7, 1),
				},
				Location: L(P(0, 1, 1), P(64, 7, 8)),
			},
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.JUMP_UNLESS), 0, 6,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.JUMP), 0, 3,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 0,
					// end
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.String("foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(7, 1),
					bytecode.NewLineInfo(8, 2),
				},
				Location: L(P(0, 1, 1), P(59, 8, 5)),
			},
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(15, 1, 16)),
			},
		},
		"empty then and else": {
			input: "a := true; unless a; end",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.TRUE),
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
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
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 12),
				},
				Location: L(P(0, 1, 1), P(23, 1, 24)),
			},
		},
		"resolve static condition with then branch": {
			input: `
				unless false
					10
				end
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(10),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(4, 1),
				},
				Location: L(P(0, 1, 1), P(33, 4, 8)),
			},
		},
		"resolve static condition with then branch to nil": {
			input: `
				unless true
					10
				end
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(4, 2),
				},
				Location: L(P(0, 1, 1), P(32, 4, 8)),
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(5, 1),
					bytecode.NewLineInfo(6, 1),
				},
				Location: L(P(0, 1, 1), P(48, 6, 8)),
			},
		},
		"with then branch": {
			input: `
				a := 5
				unless a
					a = a * 2
				end
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.JUMP_IF), 0, 11,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.JUMP), 0, 2,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(5, 1),
				},
				Location: L(P(0, 1, 1), P(47, 5, 8)),
			},
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.JUMP_IF), 0, 11,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.JUMP), 0, 5,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 2,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
					object.SmallInt(30),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(6, 2),
					bytecode.NewLineInfo(7, 1),
				},
				Location: L(P(0, 1, 1), P(68, 7, 8)),
			},
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 2,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.JUMP_IF), 0, 6,
					// then branch
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.JUMP), 0, 3,
					// else branch
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 0,
					// end
					byte(bytecode.SET_LOCAL8), 1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL8), 1,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.String("foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
					bytecode.NewLineInfo(4, 1),
					bytecode.NewLineInfo(3, 2),
					bytecode.NewLineInfo(6, 1),
					bytecode.NewLineInfo(3, 1),
					bytecode.NewLineInfo(7, 1),
					bytecode.NewLineInfo(8, 2),
				},
				Location: L(P(0, 1, 1), P(63, 8, 5)),
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 5,
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(3, 4),
				},
				Location: L(P(0, 1, 1), P(17, 3, 8)),
			},
		},
		"with a body": {
			input: `
				a := 0
				loop
					a = a + 1
				end
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.PREP_LOCALS8), 1,
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					// loop body
					byte(bytecode.GET_LOCAL8), 0,
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.ADD),
					byte(bytecode.SET_LOCAL8), 0,
					byte(bytecode.POP),
					byte(bytecode.LOOP), 0, 11,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.SmallInt(0),
					object.SmallInt(1),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 4),
					bytecode.NewLineInfo(5, 3),
				},
				Location: L(P(0, 1, 1), P(43, 5, 8)),
			},
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
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.JUMP_IF), 0, 2,
					// falsy
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// truthy
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String("foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
				},
				Location: L(P(0, 1, 1), P(18, 2, 18)),
			},
		},
		"nested": {
			input: `
				"foo" || true || 3
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.JUMP_IF), 0, 2,
					// falsy 1
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// truthy 1
					byte(bytecode.JUMP_IF), 0, 3,
					// falsy 2
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 1,
					// truthy 2
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String("foo"),
					object.SmallInt(3),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
				},
				Location: L(P(0, 1, 1), P(23, 2, 23)),
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
				"foo" && true
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.JUMP_UNLESS), 0, 2,
					// truthy
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// falsy
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String("foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 5),
				},
				Location: L(P(0, 1, 1), P(18, 2, 18)),
			},
		},
		"nested": {
			input: `
				"foo" && true && 3
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.JUMP_UNLESS), 0, 2,
					// truthy 1
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// falsy 1
					byte(bytecode.JUMP_UNLESS), 0, 3,
					// truthy 2
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8), 1,
					// falsy 2
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String("foo"),
					object.SmallInt(3),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 8),
				},
				Location: L(P(0, 1, 1), P(23, 2, 23)),
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
				"foo" ?? true
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.JUMP_IF_NIL), 0, 3,
					byte(bytecode.JUMP), 0, 2,
					// nil
					byte(bytecode.POP),
					byte(bytecode.TRUE),
					// not nil
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String("foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 6),
				},
				Location: L(P(0, 1, 1), P(18, 2, 18)),
			},
		},
		"nested": {
			input: `
				"foo" ?? true ?? 3
			`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
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
					byte(bytecode.CONSTANT8), 1,
					// not nil 2
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					object.String("foo"),
					object.SmallInt(3),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 10),
				},
				Location: L(P(0, 1, 1), P(23, 2, 23)),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}
