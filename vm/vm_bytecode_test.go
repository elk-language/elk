package vm

import (
	"strings"
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Represents a single VM test case.
type bytecodeTestCase struct {
	chunk        *value.BytecodeFunction
	wantStackTop value.Value
	wantStdout   string
	wantErr      value.Value
	teardown     func()
}

// Type of the compiler test table.
type bytecodeTestTable map[string]bytecodeTestCase

func vmBytecodeTest(tc bytecodeTestCase, t *testing.T) {
	t.Helper()

	var stdout strings.Builder
	vm := New(WithStdout(&stdout))
	gotStackTop, gotErr := vm.InterpretTopLevel(tc.chunk)
	gotStdout := stdout.String()
	if tc.teardown != nil {
		tc.teardown()
	}
	opts := []cmp.Option{
		cmp.AllowUnexported(value.Error{}, value.BigFloat{}, value.BigInt{}),
		cmpopts.IgnoreUnexported(value.Class{}, value.Module{}),
		cmpopts.IgnoreFields(value.Class{}, "ConstructorFunc"),
		value.FloatComparer,
		value.Float32Comparer,
		value.Float64Comparer,
		value.BigFloatComparer,
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x20,
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x20: value.Int8(5),
				},
			},

			wantStackTop: value.Int8(5),
		},
		"load 16bit constant": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT16), 0x01, 0x00,
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x100: value.Int8(5),
				},
			},

			wantStackTop: value.Int8(5),
		},
		"load 32bit constant": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT32), 0x01, 0x00, 0x00, 0x00,
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x1000000: value.Int8(5),
				},
			},

			wantStackTop: value.Int8(5),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.NewBigFloat(25.3),
				},
			},
			wantStackTop: value.NewBigFloat(-25.3),
		},
		"negate Float": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Float(25.3),
				},
			},
			wantStackTop: value.Float(-25.3),
		},
		"negate Float64": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Float64(25.3),
				},
			},
			wantStackTop: value.Float64(-25.3),
		},
		"negate Float32": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Float32(25.3),
				},
			},
			wantStackTop: value.Float32(-25.3),
		},
		"negate BigInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.NewBigInt(5),
				},
			},
			wantStackTop: value.NewBigInt(-5),
		},
		"negate SmallInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.SmallInt(5),
				},
			},
			wantStackTop: value.SmallInt(-5),
		},
		"negate Int64": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int64(5),
				},
			},
			wantStackTop: value.Int64(-5),
		},
		"negate UInt64": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.UInt64(5),
				},
			},
			wantStackTop: value.UInt64(18446744073709551611),
		},
		"negate Int32": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int32(5),
				},
			},
			wantStackTop: value.Int32(-5),
		},
		"negate UInt32": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.UInt32(5),
				},
			},
			wantStackTop: value.UInt32(4294967291),
		},
		"negate Int16": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int16(5),
				},
			},

			wantStackTop: value.Int16(-5),
		},
		"negate UInt16": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.UInt16(5),
				},
			},

			wantStackTop: value.UInt16(65531),
		},
		"negate Int8": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(5),
				},
			},

			wantStackTop: value.Int8(-5),
		},
		"negate UInt8": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.UInt8(5),
				},
			},

			wantStackTop: value.UInt8(251),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
			},

			wantStackTop: value.False,
		},
		"put true": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
			},

			wantStackTop: value.True,
		},
		"put nil": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
			},

			wantStackTop: value.Nil,
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.String("foo"),
				},
			},

			wantStackTop: value.False,
		},
		"bool not int": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(0),
				},
			},

			wantStackTop: value.False,
		},
		"bool not true": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.TRUE),
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
			},

			wantStackTop: value.False,
		},
		"bool not false": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.FALSE),
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(0),
				},
			},

			wantStackTop: value.True,
		},
		"bool not nil": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(0),
				},
			},

			wantStackTop: value.True,
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(5),
					0x1: value.Int8(25),
				},
			},

			wantStackTop: value.Int8(30),
		},
		"add String to String": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.String("foo"),
					0x1: value.String("bar"),
				},
			},

			wantStackTop: value.String("foobar"),
		},
		"add String to Char": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Char('f'),
					0x1: value.String("oo"),
				},
			},

			wantStackTop: value.String("foo"),
		},
		"add Int8 to String": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.String("foo"),
					0x1: value.Int8(5),
				},
			},
			wantStackTop: value.String("foo"),
			wantErr:      value.Errorf(value.TypeErrorClass, `can't concat 5i8 to string "foo"`),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(5),
					0x1: value.Int8(25),
				},
			},

			wantStackTop: value.Int8(-20),
		},
		"String - String": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.String("foobar"),
					0x1: value.String("bar"),
				},
			},

			wantStackTop: value.String("foo"),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(5),
					0x1: value.Int8(25),
				},
			},

			wantStackTop: value.Int8(125),
		},
		"String * SmallInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.String("foo"),
					0x1: value.SmallInt(3),
				},
			},

			wantStackTop: value.String("foofoofoo"),
		},
		"Char * SmallInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Char('a'),
					0x1: value.SmallInt(3),
				},
			},

			wantStackTop: value.String("aaa"),
		},
		"BigFloat * Float": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.NewBigFloat(5.5),
					0x1: value.Float(10),
				},
			},

			wantStackTop: value.NewBigFloat(55),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(35),
					0x1: value.Int8(5),
				},
			},

			wantStackTop: value.Int8(7),
		},
		"String / SmallInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.String("foo"),
					0x1: value.SmallInt(3),
				},
			},

			wantStackTop: value.String("foo"),
			wantErr:      value.NewNoMethodError("/", value.String("foo")),
		},
		"BigFloat / Float": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.NewBigFloat(6.8),
					0x1: value.Float(2),
				},
			},

			wantStackTop: value.NewBigFloat(3.4),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_Exponentiate(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 ** Int8": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(2),
					0x1: value.Int8(5),
				},
			},

			wantStackTop: value.Int8(32),
		},
		"String ** SmallInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.String("foo"),
					0x1: value.SmallInt(3),
				},
			},

			wantStackTop: value.String("foo"),
			wantErr:      value.NewNoMethodError("**", value.String("foo")),
		},
		"BigFloat ** Float": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.NewBigFloat(6.5),
					0x1: value.Float(2),
				},
			},

			wantStackTop: value.NewBigFloat(42.25),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_Modulo(t *testing.T) {
	tests := bytecodeTestTable{
		"Int8 % Int8": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.MODULO),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(25),
					0x1: value.Int8(4),
				},
			},

			wantStackTop: value.Int8(1),
		},
		"BigFloat % Float": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.MODULO),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.NewBigFloat(68.5),
					0x1: value.Float(20.5),
				},
			},

			wantStackTop: value.NewBigFloat(7),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(35),
					0x1: value.Int64(1),
				},
			},

			wantStackTop: value.Int8(17),
		},
		"Int >> Int": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.SmallInt(12),
					0x1: value.SmallInt(2),
				},
			},

			wantStackTop: value.SmallInt(3),
		},
		"-Int16 >> UInt32": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int16(-6),
					0x1: value.UInt32(1),
				},
			},

			wantStackTop: value.Int16(-3),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(35),
					0x1: value.Int64(1),
				},
			},

			wantStackTop: value.Int8(70),
		},
		"Int << Int": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.SmallInt(12),
					0x1: value.SmallInt(2),
				},
			},

			wantStackTop: value.SmallInt(48),
		},
		"-Int16 << UInt32": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int16(-6),
					0x1: value.UInt32(1),
				},
			},

			wantStackTop: value.Int16(-12),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(35),
					0x1: value.Int64(1),
				},
			},

			wantStackTop: value.Int8(17),
		},
		"-Int16 >>> UInt32": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int16(-6),
					0x1: value.UInt32(1),
				},
			},

			wantStackTop: value.Int16(32765),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(35),
					0x1: value.Int64(1),
				},
			},

			wantStackTop: value.Int8(70),
		},
		"UInt16 <<< Int": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.UInt16(12),
					0x1: value.SmallInt(2),
				},
			},

			wantStackTop: value.UInt16(48),
		},
		"-Int16 <<< UInt32": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int16(-6),
					0x1: value.UInt32(1),
				},
			},

			wantStackTop: value.Int16(-12),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(23),
					0x1: value.Int8(10),
				},
			},

			wantStackTop: value.Int8(2),
		},
		"UInt16 & UInt16": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.UInt16(235),
					0x1: value.UInt16(58),
				},
			},

			wantStackTop: value.UInt16(42),
		},
		"SmallInt & SmallInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.SmallInt(235),
					0x1: value.SmallInt(58),
				},
			},

			wantStackTop: value.SmallInt(42),
		},
		"SmallInt & BigInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.SmallInt(255),
					0x1: value.ParseBigIntPanic("9223372036857247042", 10),
				},
			},

			wantStackTop: value.SmallInt(66),
		},
		"BigInt & SmallInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.ParseBigIntPanic("9223372036857247042", 10),
					0x1: value.SmallInt(255),
				},
			},

			wantStackTop: value.SmallInt(66),
		},
		"BigInt & BigInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.ParseBigIntPanic("9223372036857247042", 10),
					0x1: value.ParseBigIntPanic("10223372099998981329", 10),
				},
			},

			wantStackTop: value.ParseBigIntPanic("9223372036855043136", 10),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(23),
					0x1: value.Int8(10),
				},
			},

			wantStackTop: value.Int8(31),
		},
		"UInt16 | UInt16": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.UInt16(235),
					0x1: value.UInt16(58),
				},
			},

			wantStackTop: value.UInt16(251),
		},
		"SmallInt | SmallInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.SmallInt(235),
					0x1: value.SmallInt(58),
				},
			},

			wantStackTop: value.SmallInt(251),
		},
		"SmallInt | BigInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.SmallInt(255),
					0x1: value.ParseBigIntPanic("9223372036857247042", 10),
				},
			},

			wantStackTop: value.ParseBigIntPanic("9223372036857247231", 10),
		},
		"BigInt | SmallInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.ParseBigIntPanic("9223372036857247042", 10),
					0x1: value.SmallInt(255),
				},
			},

			wantStackTop: value.ParseBigIntPanic("9223372036857247231", 10),
		},
		"BigInt | BigInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.ParseBigIntPanic("9223372036857247042", 10),
					0x1: value.ParseBigIntPanic("10223372099998981329", 10),
				},
			},

			wantStackTop: value.ParseBigIntPanic("10223372100001185235", 10),
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
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.Int8(23),
					0x1: value.Int8(10),
				},
			},

			wantStackTop: value.Int8(29),
		},
		"UInt16 ^ UInt16": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.UInt16(235),
					0x1: value.UInt16(58),
				},
			},

			wantStackTop: value.UInt16(209),
		},
		"SmallInt ^ SmallInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.SmallInt(235),
					0x1: value.SmallInt(58),
				},
			},

			wantStackTop: value.SmallInt(209),
		},
		"SmallInt ^ BigInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.SmallInt(255),
					0x1: value.ParseBigIntPanic("9223372036857247042", 10),
				},
			},

			wantStackTop: value.ParseBigIntPanic("9223372036857247165", 10),
		},
		"BigInt ^ SmallInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.ParseBigIntPanic("9223372036857247042", 10),
					0x1: value.SmallInt(255),
				},
			},

			wantStackTop: value.ParseBigIntPanic("9223372036857247165", 10),
		},
		"BigInt ^ BigInt": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0x0,
					byte(bytecode.CONSTANT8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					0x0: value.ParseBigIntPanic("9223372036857247042", 10),
					0x1: value.ParseBigIntPanic("10223372099998981329", 10),
				},
			},

			wantStackTop: value.SmallInt(1000000063146142099),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_GetModConst(t *testing.T) {
	tests := bytecodeTestTable{
		"get constant under Root": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					value.SymbolTable.Add("Std"),
				},
			},

			wantStackTop: value.StdModule,
		},
		"get nested constants": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.GET_MOD_CONST8), 2,
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					value.SymbolTable.Add("Std"),
					value.SymbolTable.Add("Float"),
					value.SymbolTable.Add("INF"),
				},
			},

			wantStackTop: value.FloatInf(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_DefModConst(t *testing.T) {
	tests := bytecodeTestTable{
		"define constant under Root": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.DEF_MOD_CONST8), 1,
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					value.String("constant!"),
					value.SymbolTable.Add("Foo"),
				},
			},
			teardown:     func() { value.RootModule.Constants.DeleteString("Foo") },
			wantStackTop: value.String("constant!"),
		},
		"define constant under Root and read it": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.DEF_MOD_CONST8), 1,
					byte(bytecode.POP),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					value.String("constant!"),
					value.SymbolTable.Add("Foo"),
				},
			},
			teardown:     func() { value.RootModule.Constants.DeleteString("Foo") },
			wantStackTop: value.String("constant!"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_DefClass(t *testing.T) {
	tests := bytecodeTestTable{
		"define class without a body or superclass": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.CONSTANT_BASE),
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.NIL),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					value.SymbolTable.Add("Foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				Location: L(P(0, 1, 1), P(13, 1, 14)),
			},
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
			),
		},
		"define class with a superclass without a body": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.CONSTANT_BASE),
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.GET_MOD_CONST8), 2,
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					value.SymbolTable.Add("Foo"),
					value.SymbolTable.Add("Std"),
					value.SymbolTable.Add("Error"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				Location: L(P(0, 1, 1), P(13, 1, 14)),
			},
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithParent(value.ErrorClass),
			),
		},
		"define class with a body": {
			chunk: &value.BytecodeFunction{
				Instructions: []byte{
					byte(bytecode.CONSTANT8), 0,
					byte(bytecode.CONSTANT_BASE),
					byte(bytecode.CONSTANT8), 1,
					byte(bytecode.NIL),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				Constants: []value.Value{
					&value.BytecodeFunction{
						Instructions: []byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.CONSTANT8), 0,
							byte(bytecode.CONSTANT_BASE),
							byte(bytecode.DEF_MOD_CONST8), 1,
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						Constants: []value.Value{
							value.SmallInt(1),
							value.SymbolTable.Add("Bar"),
						},
						LineInfoList: bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
						},
						Location: L(P(5, 2, 5), P(44, 5, 7)),
					},
					value.SymbolTable.Add("Foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				Location: L(P(0, 1, 1), P(45, 5, 8)),
			},
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithName("Foo"),
				value.ClassWithConstants(
					value.SimpleSymbolMap{
						value.SymbolTable.Add("Bar"): value.SmallInt(1),
					},
				),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}
