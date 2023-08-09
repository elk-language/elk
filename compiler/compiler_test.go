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
		t.Fatalf(diff)
	}
	if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
		t.Fatalf(diff)
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
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.UInt8(1),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(2, 1, 3)),
			},
		},
		"put UInt16": {
			input: "25u16",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.UInt16(25),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
		},
		"put UInt32": {
			input: "450_200u32",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.UInt32(450200),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"put UInt64": {
			input: "450_200u64",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.UInt64(450200),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"put Int8": {
			input: "1i8",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.Int8(1),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(2, 1, 3)),
			},
		},
		"put Int16": {
			input: "25i16",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.Int16(25),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
		},
		"put Int32": {
			input: "450_200i32",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.Int32(450200),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"put Int64": {
			input: "450_200i64",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.Int64(450200),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"put SmallInt": {
			input: "450_200",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.SmallInt(450200),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
			},
		},
		"put BigInt": {
			input: (&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5)).String(),
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.ToElkBigInt((&big.Int{}).Add(big.NewInt(math.MaxInt64), big.NewInt(5))),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
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
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.Float64(45.5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
			},
		},
		"put Float32": {
			input: "45.5f32",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.Float32(45.5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
			},
		},
		"put Float": {
			input: "45.5",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.Float(45.5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(3, 1, 4)),
			},
		},
		"put Raw String": {
			input: `'foo\n'`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.String(`foo\n`),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
			},
		},
		"put String": {
			input: `"foo\n"`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.String("foo\n"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
			},
		},
		"put raw Char": {
			input: `c'I'`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.Char('I'),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(3, 1, 4)),
			},
		},
		"put Char": {
			input: `c"\n"`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.Char('\n'),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
		},
		"put nil": {
			input: `nil`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.NIL),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(2, 1, 3)),
			},
		},
		"put true": {
			input: `true`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.TRUE),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(3, 1, 4)),
			},
		},
		"put false": {
			input: `false`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.FALSE),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
		},
		"put simple Symbol": {
			input: `:foo`,
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.SymbolTable.Add("foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
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
		"add": {
			input: "1i8 + 5i8",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.CONSTANT8),
					1,
					byte(bytecode.ADD),
				},
				Constants: []object.Value{
					object.Int8(1),
					object.Int8(5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(8, 1, 9)),
			},
		},
		"subtract": {
			input: "151i32 - 25i32 - 5i32",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.CONSTANT8),
					1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.CONSTANT8),
					2,
					byte(bytecode.SUBTRACT),
				},
				Constants: []object.Value{
					object.Int32(151),
					object.Int32(25),
					object.Int32(5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 5),
				},
				Location: L(P(0, 1, 1), P(20, 1, 21)),
			},
		},
		"multiply": {
			input: "45.5 * 2.5",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.CONSTANT8),
					1,
					byte(bytecode.MULTIPLY),
				},
				Constants: []object.Value{
					object.Float(45.5),
					object.Float(2.5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"divide": {
			input: "45.5 / .5",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.CONSTANT8),
					1,
					byte(bytecode.DIVIDE),
				},
				Constants: []object.Value{
					object.Float(45.5),
					object.Float(0.5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(8, 1, 9)),
			},
		},
		"exponentiate": {
			input: "-2 ** 2",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.CONSTANT8),
					1,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.NEGATE),
				},
				Constants: []object.Value{
					object.SmallInt(2),
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				Location: L(P(0, 1, 1), P(6, 1, 7)),
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
		"negate": {
			input: "-5",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.NEGATE),
				},
				Constants: []object.Value{
					object.SmallInt(5),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(1, 1, 2)),
			},
		},
		"bitwise not": {
			input: "~10",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.BITWISE_NOT),
				},
				Constants: []object.Value{
					object.SmallInt(10),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(2, 1, 3)),
			},
		},
		"logical not": {
			input: "!10",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.NOT),
				},
				Constants: []object.Value{
					object.SmallInt(10),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(2, 1, 3)),
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
					byte(bytecode.REGISTER_LOCALS),
					1,
					byte(bytecode.NIL),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(4, 1, 5)),
			},
		},
		"declare with a type": {
			input: "var a: Int",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.REGISTER_LOCALS),
					1,
					byte(bytecode.NIL),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 2),
				},
				Location: L(P(0, 1, 1), P(9, 1, 10)),
			},
		},
		"declare and initialise": {
			input: "var a = 3",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.REGISTER_LOCALS),
					1,
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.SET_LOCAL),
					0,
				},
				Constants: []object.Value{
					object.SmallInt(3),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(8, 1, 9)),
			},
		},
		"read undeclared": {
			input: "a",
			want: &bytecode.Chunk{
				Instructions: []byte{},
				Location:     L(P(0, 1, 1), P(0, 1, 1)),
			},
			err: errors.ErrorList{
				errors.NewError(L(P(0, 1, 1), P(0, 1, 1)), "undeclared variable: a"),
			},
		},
		"assign undeclared": {
			input: "a = 3",
			want: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0,
				},
				Constants: []object.Value{
					object.SmallInt(3),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 1),
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
					byte(bytecode.REGISTER_LOCALS),
					1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.SET_LOCAL),
					0,
				},
				Constants: []object.Value{
					object.String("foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 2),
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
					byte(bytecode.REGISTER_LOCALS),
					1,
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.SET_LOCAL),
					0,
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8),
					1,
					byte(bytecode.SET_LOCAL),
					0,
				},
				Constants: []object.Value{
					object.String("foo"),
					object.String("bar"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 2),
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
					byte(bytecode.REGISTER_LOCALS),
					1,
					byte(bytecode.NIL),
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.ADD),
				},
				Constants: []object.Value{
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 3),
					bytecode.NewLineInfo(3, 2),
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
					byte(bytecode.REGISTER_LOCALS),
					1,
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.SET_LOCAL),
					0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL),
					0,
					byte(bytecode.CONSTANT8),
					1,
					byte(bytecode.ADD),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(3, 3),
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
					byte(bytecode.REGISTER_LOCALS),
					1,
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.SET_LOCAL),
					0,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL),
					0,
					byte(bytecode.CONSTANT8),
					1,
					byte(bytecode.ADD),
				},
				Constants: []object.Value{
					object.SmallInt(5),
					object.SmallInt(2),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 4),
					bytecode.NewLineInfo(4, 3),
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
					byte(bytecode.REGISTER_LOCALS),
					2,
					byte(bytecode.CONSTANT8),
					0,
					byte(bytecode.SET_LOCAL),
					0,
					byte(bytecode.POP),
					byte(bytecode.CONSTANT8),
					1,
					byte(bytecode.CONSTANT8),
					2,
					byte(bytecode.SET_LOCAL),
					1,
					byte(bytecode.POP),
					byte(bytecode.GET_LOCAL),
					1,
					byte(bytecode.CONSTANT8),
					3,
					byte(bytecode.ADD),
					byte(bytecode.LEAVE_SCOPE),
					1,
					1,
					byte(bytecode.ADD),
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
