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

func TestVM_LoadConstant(t *testing.T) {
	tests := bytecodeTestTable{
		"load 8bit constant": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x20,
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
					byte(bytecode.CONSTANT16), 0x01, 0x00,
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
					byte(bytecode.CONSTANT32), 0x01, 0x00, 0x00, 0x00,
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

func TestVM_Negate(t *testing.T) {
	tests := bytecodeTestTable{
		"negate BigFloat": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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

func TestVM_PutValue(t *testing.T) {
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

func TestVM_BoolNot(t *testing.T) {
	tests := bytecodeTestTable{
		"bool not string": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
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
					byte(bytecode.CONSTANT8), 0x0,
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

func TestVM_Add(t *testing.T) {
	tests := bytecodeTestTable{
		"add Int8 to Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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

func TestVM_Subtract(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 - Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.String("foobar"),
					0x1: object.String("bar"),
				},
			},

			wantStackTop: object.String("foo"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_Multiply(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 * Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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

func TestVM_Divide(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 / Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
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

func TestVM_RightBitshift(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 >> Int64": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(35),
					0x1: object.Int64(1),
				},
			},

			wantStackTop: object.Int8(17),
		},
		"Int >> Int": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.SmallInt(12),
					0x1: object.SmallInt(2),
				},
			},

			wantStackTop: object.SmallInt(3),
		},
		"-Int16 >> UInt32": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int16(-6),
					0x1: object.UInt32(1),
				},
			},

			wantStackTop: object.Int16(-3),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_LeftBitshift(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 << Int64": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(35),
					0x1: object.Int64(1),
				},
			},

			wantStackTop: object.Int8(70),
		},
		"Int << Int": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.SmallInt(12),
					0x1: object.SmallInt(2),
				},
			},

			wantStackTop: object.SmallInt(48),
		},
		"-Int16 << UInt32": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int16(-6),
					0x1: object.UInt32(1),
				},
			},

			wantStackTop: object.Int16(-12),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_LogicalRightBitshift(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 >>> Int64": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(35),
					0x1: object.Int64(1),
				},
			},

			wantStackTop: object.Int8(17),
		},
		"-Int16 >>> UInt32": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int16(-6),
					0x1: object.UInt32(1),
				},
			},

			wantStackTop: object.Int16(32765),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_LogicalLeftBitshift(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 <<< Int64": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(35),
					0x1: object.Int64(1),
				},
			},

			wantStackTop: object.Int8(70),
		},
		"UInt16 <<< Int": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.UInt16(12),
					0x1: object.SmallInt(2),
				},
			},

			wantStackTop: object.UInt16(48),
		},
		"-Int16 <<< UInt32": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int16(-6),
					0x1: object.UInt32(1),
				},
			},

			wantStackTop: object.Int16(-12),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_BitwiseAnd(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 & Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(23),
					0x1: object.Int8(10),
				},
			},

			wantStackTop: object.Int8(2),
		},
		"UInt16 & UInt16": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.UInt16(235),
					0x1: object.UInt16(58),
				},
			},

			wantStackTop: object.UInt16(42),
		},
		"SmallInt & SmallInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.SmallInt(235),
					0x1: object.SmallInt(58),
				},
			},

			wantStackTop: object.SmallInt(42),
		},
		"SmallInt & BigInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.SmallInt(255),
					0x1: object.ParseBigIntPanic("9223372036857247042", 10),
				},
			},

			wantStackTop: object.SmallInt(66),
		},
		"BigInt & SmallInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.ParseBigIntPanic("9223372036857247042", 10),
					0x1: object.SmallInt(255),
				},
			},

			wantStackTop: object.SmallInt(66),
		},
		"BigInt & BigInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.ParseBigIntPanic("9223372036857247042", 10),
					0x1: object.ParseBigIntPanic("10223372099998981329", 10),
				},
			},

			wantStackTop: object.ParseBigIntPanic("9223372036855043136", 10),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_BitwiseOr(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 | Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(23),
					0x1: object.Int8(10),
				},
			},

			wantStackTop: object.Int8(31),
		},
		"UInt16 | UInt16": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.UInt16(235),
					0x1: object.UInt16(58),
				},
			},

			wantStackTop: object.UInt16(251),
		},
		"SmallInt | SmallInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.SmallInt(235),
					0x1: object.SmallInt(58),
				},
			},

			wantStackTop: object.SmallInt(251),
		},
		"SmallInt | BigInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.SmallInt(255),
					0x1: object.ParseBigIntPanic("9223372036857247042", 10),
				},
			},

			wantStackTop: object.ParseBigIntPanic("9223372036857247231", 10),
		},
		"BigInt | SmallInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.ParseBigIntPanic("9223372036857247042", 10),
					0x1: object.SmallInt(255),
				},
			},

			wantStackTop: object.ParseBigIntPanic("9223372036857247231", 10),
		},
		"BigInt | BigInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.ParseBigIntPanic("9223372036857247042", 10),
					0x1: object.ParseBigIntPanic("10223372099998981329", 10),
				},
			},

			wantStackTop: object.ParseBigIntPanic("10223372100001185235", 10),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_BitwiseXor(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 ^ Int8": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.Int8(23),
					0x1: object.Int8(10),
				},
			},

			wantStackTop: object.Int8(29),
		},
		"UInt16 ^ UInt16": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.UInt16(235),
					0x1: object.UInt16(58),
				},
			},

			wantStackTop: object.UInt16(209),
		},
		"SmallInt ^ SmallInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.SmallInt(235),
					0x1: object.SmallInt(58),
				},
			},

			wantStackTop: object.SmallInt(209),
		},
		"SmallInt ^ BigInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.SmallInt(255),
					0x1: object.ParseBigIntPanic("9223372036857247042", 10),
				},
			},

			wantStackTop: object.ParseBigIntPanic("9223372036857247165", 10),
		},
		"BigInt ^ SmallInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.ParseBigIntPanic("9223372036857247042", 10),
					0x1: object.SmallInt(255),
				},
			},

			wantStackTop: object.ParseBigIntPanic("9223372036857247165", 10),
		},
		"BigInt ^ BigInt": {
			chunk: &bytecode.Chunk{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []object.Value{
					0x0: object.ParseBigIntPanic("9223372036857247042", 10),
					0x1: object.ParseBigIntPanic("10223372099998981329", 10),
				},
			},

			wantStackTop: object.SmallInt(1000000063146142099),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}
