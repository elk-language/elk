package compiler

import (
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
	if diff := cmp.Diff(tc.want, got); diff != "" {
		t.Fatalf(diff)
	}
	if diff := cmp.Diff(tc.err, err); diff != "" {
		t.Fatalf(diff)
	}
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
				Location: position.NewLocation("main", 0, 3, 1, 1),
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
				Location: position.NewLocation("main", 0, 5, 1, 1),
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
				Location: position.NewLocation("main", 0, 10, 1, 1),
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
				Location: position.NewLocation("main", 0, 10, 1, 1),
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
				Location: position.NewLocation("main", 0, 3, 1, 1),
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
				Location: position.NewLocation("main", 0, 5, 1, 1),
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
				Location: position.NewLocation("main", 0, 10, 1, 1),
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
				Location: position.NewLocation("main", 0, 10, 1, 1),
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
				Location: position.NewLocation("main", 0, 7, 1, 1),
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
				Location: position.NewLocation("main", 0, 7, 1, 1),
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
				Location: position.NewLocation("main", 0, 4, 1, 1),
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
				Location: position.NewLocation("main", 0, 7, 1, 1),
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
				Location: position.NewLocation("main", 0, 7, 1, 1),
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
				Location: position.NewLocation("main", 0, 9, 1, 1),
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
				Location: position.NewLocation("main", 0, 21, 1, 1),
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
				Location: position.NewLocation("main", 0, 10, 1, 1),
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
				Location: position.NewLocation("main", 0, 9, 1, 1),
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
				Location: position.NewLocation("main", 0, 7, 1, 1),
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
				Location: position.NewLocation("main", 0, 2, 1, 1),
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
				Location: position.NewLocation("main", 0, 3, 1, 1),
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
				Location: position.NewLocation("main", 0, 3, 1, 1),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			compilerTest(tc, t)
		})
	}
}
