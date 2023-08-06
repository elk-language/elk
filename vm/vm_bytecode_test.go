package vm

import (
	"strings"
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/object"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Represents a single VM test case.
type bytecodeTestCase struct {
	chunk        *bytecode.Chunk
	wantStackTop object.Value
	wantStdout   string
	wantErr      object.Value
}

// Type of the compiler test table.
type bytecodeTestTable map[string]bytecodeTestCase

func vmBytecodeTest(tc bytecodeTestCase, t *testing.T) {
	t.Helper()

	var stdout strings.Builder
	vm := New(WithStdout(&stdout))
	gotStackTop, gotErr := vm.InterpretBytecode(tc.chunk)
	gotStdout := stdout.String()
	opts := []cmp.Option{
		cmp.AllowUnexported(object.Error{}, object.BigFloat{}, object.BigInt{}),
		cmpopts.IgnoreUnexported(object.Class{}),
		cmpopts.IgnoreFields(object.Class{}, "ConstructorFunc"),
	}
	if diff := cmp.Diff(tc.wantErr, gotErr, opts...); diff != "" {
		t.Fatalf(diff)
	}
	if diff := cmp.Diff(tc.wantStdout, gotStdout, opts...); diff != "" {
		t.Fatalf(diff)
	}
	if diff := cmp.Diff(tc.wantStackTop, gotStackTop, opts...); diff != "" {
		t.Fatalf(diff)
	}
}

func TestVMLoadConstant(t *testing.T) {
	tests := bytecodeTestTable{
		"load 8bit constant": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x20,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x20: object.Int8(5),
				},
			},

			wantStackTop: object.Int8(5),
		},
		"load 16bit constant": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT16),
					0x01,
					0x00,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x100: object.Int8(5),
				},
			},

			wantStackTop: object.Int8(5),
		},
		"load 32bit constant": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT32),
					0x01,
					0x00,
					0x00,
					0x00,
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x1000000: object.Int8(5),
				},
			},

			wantStackTop: object.Int8(5),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVMNegate(t *testing.T) {
	tests := bytecodeTestTable{
		"negate BigFloat": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.NewBigFloat(25.3),
				},
			},
			wantStackTop: object.NewBigFloat(-25.3),
		},
		"negate Float": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Float(25.3),
				},
			},
			wantStackTop: object.Float(-25.3),
		},
		"negate Float64": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Float64(25.3),
				},
			},
			wantStackTop: object.Float64(-25.3),
		},
		"negate Float32": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Float32(25.3),
				},
			},
			wantStackTop: object.Float32(-25.3),
		},
		"negate BigInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.NewBigInt(5),
				},
			},
			wantStackTop: object.NewBigInt(-5),
		},
		"negate SmallInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.SmallInt(5),
				},
			},
			wantStackTop: object.SmallInt(-5),
		},
		"negate Int64": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int64(5),
				},
			},
			wantStackTop: object.Int64(-5),
		},
		"negate UInt64": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.UInt64(5),
				},
			},
			wantStackTop: object.UInt64(18446744073709551611),
		},
		"negate Int32": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int32(5),
				},
			},
			wantStackTop: object.Int32(-5),
		},
		"negate UInt32": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.UInt32(5),
				},
			},
			wantStackTop: object.UInt32(4294967291),
		},
		"negate Int16": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int16(5),
				},
			},

			wantStackTop: object.Int16(-5),
		},
		"negate UInt16": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.UInt16(5),
				},
			},

			wantStackTop: object.UInt16(65531),
		},
		"negate Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(5),
				},
			},

			wantStackTop: object.Int8(-5),
		},
		"negate UInt8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.UInt8(5),
				},
			},

			wantStackTop: object.UInt8(251),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVMPutValue(t *testing.T) {
	tests := bytecodeTestTable{
		"put false": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
			},

			wantStackTop: object.False,
		},
		"put true": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
			},

			wantStackTop: object.True,
		},
		"put nil": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
			},

			wantStackTop: object.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVMBoolNot(t *testing.T) {
	tests := bytecodeTestTable{
		"bool not string": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.String("foo"),
				},
			},

			wantStackTop: object.False,
		},
		"bool not int": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(0),
				},
			},

			wantStackTop: object.False,
		},
		"bool not true": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.TRUE),
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
			},

			wantStackTop: object.False,
		},
		"bool not false": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.FALSE),
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(0),
				},
			},

			wantStackTop: object.True,
		},
		"bool not nil": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(0),
				},
			},

			wantStackTop: object.True,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVMAdd(t *testing.T) {
	tests := bytecodeTestTable{
		"add Int8 to Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(5),
					0x1: object.Int8(25),
				},
			},

			wantStackTop: object.Int8(30),
		},
		"add String to String": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.String("foo"),
					0x1: object.String("bar"),
				},
			},

			wantStackTop: object.String("foobar"),
		},
		"add String to Char": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Char('f'),
					0x1: object.String("oo"),
				},
			},

			wantStackTop: object.String("foo"),
		},
		"add Int8 to String": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.String("foo"),
					0x1: object.Int8(5),
				},
			},
			wantStackTop: object.String("foo"),
			wantErr:      object.Errorf(object.TypeErrorClass, `can't concat 5i8 to string "foo"`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVMSubtract(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 - Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(5),
					0x1: object.Int8(25),
				},
			},

			wantStackTop: object.Int8(-20),
		},
		"String - String": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.String("foo"),
					0x1: object.String("bar"),
				},
			},

			wantStackTop: object.String("foo"),
			wantErr:      object.NewNoMethodError("-", object.String("foo")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVMMultiply(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 * Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(5),
					0x1: object.Int8(25),
				},
			},

			wantStackTop: object.Int8(125),
		},
		"String * SmallInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.String("foo"),
					0x1: object.SmallInt(3),
				},
			},

			wantStackTop: object.String("foofoofoo"),
		},
		"Char * SmallInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Char('a'),
					0x1: object.SmallInt(3),
				},
			},

			wantStackTop: object.String("aaa"),
		},
		"BigFloat * Float": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.NewBigFloat(5.5),
					0x1: object.Float(10),
				},
			},

			wantStackTop: object.NewBigFloat(55),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVMDivide(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 / Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(35),
					0x1: object.Int8(5),
				},
			},

			wantStackTop: object.Int8(7),
		},
		"String / SmallInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.String("foo"),
					0x1: object.SmallInt(3),
				},
			},

			wantStackTop: object.String("foo"),
			wantErr:      object.NewNoMethodError("/", object.String("foo")),
		},
		"BigFloat / Float": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8),
					0x0,
					byte(bytecode.CONSTANT8),
					0x1,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.NewBigFloat(6.8),
					0x1: object.Float(2),
				},
			},

			wantStackTop: object.NewBigFloat(3.4),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}
