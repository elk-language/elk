package vm_test

import (
	"strings"
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

// Represents a single VM test case.
type bytecodeTestCase struct {
	chunk        *vm.BytecodeMethod
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
	vm := vm.New(vm.WithStdout(&stdout))
	gotStackTop, gotErr := vm.InterpretTopLevel(tc.chunk)
	gotStdout := stdout.String()
	if tc.teardown != nil {
		tc.teardown()
	}
	opts := comparer.Comparer
	if diff := cmp.Diff(tc.wantErr, gotErr, opts...); diff != "" {
		t.Fatalf(diff)
	}
	if diff := cmp.Diff(tc.wantStdout, gotStdout, opts...); diff != "" {
		t.Fatalf(diff)
	}
	if tc.wantErr != nil {
		return
	}

	if diff := cmp.Diff(tc.wantStackTop, gotStackTop, opts...); diff != "" {
		t.Fatalf(diff)
	}
}

func TestVM_LoadConstant(t *testing.T) {
	tests := bytecodeTestTable{
		"load 8bit constant": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x20,
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x20: value.Int8(5),
				},
			},

			wantStackTop: value.Int8(5),
		},
		"load 16bit constant": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE16), 0x01, 0x00,
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x100: value.Int8(5),
				},
			},

			wantStackTop: value.Int8(5),
		},
		"load 32bit constant": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE32), 0x01, 0x00, 0x00, 0x00,
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.NewBigFloat(25.3),
				},
			},
			wantStackTop: value.NewBigFloat(-25.3),
		},
		"negate Float": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Float(25.3),
				},
			},
			wantStackTop: value.Float(-25.3),
		},
		"negate Float64": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Float64(25.3),
				},
			},
			wantStackTop: value.Float64(-25.3),
		},
		"negate Float32": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Float32(25.3),
				},
			},
			wantStackTop: value.Float32(-25.3),
		},
		"negate BigInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.NewBigInt(5),
				},
			},
			wantStackTop: value.NewBigInt(-5),
		},
		"negate SmallInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.SmallInt(5),
				},
			},
			wantStackTop: value.SmallInt(-5),
		},
		"negate Int64": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int64(5),
				},
			},
			wantStackTop: value.Int64(-5),
		},
		"negate UInt64": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.UInt64(5),
				},
			},
			wantStackTop: value.UInt64(18446744073709551611),
		},
		"negate Int32": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int32(5),
				},
			},
			wantStackTop: value.Int32(-5),
		},
		"negate UInt32": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.UInt32(5),
				},
			},
			wantStackTop: value.UInt32(4294967291),
		},
		"negate Int16": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int16(5),
				},
			},

			wantStackTop: value.Int16(-5),
		},
		"negate UInt16": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.UInt16(5),
				},
			},

			wantStackTop: value.UInt16(65531),
		},
		"negate Int8": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(5),
				},
			},

			wantStackTop: value.Int8(-5),
		},
		"negate UInt8": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NEGATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.FALSE),
					byte(bytecode.RETURN),
				},
			},

			wantStackTop: value.False,
		},
		"put true": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.TRUE),
					byte(bytecode.RETURN),
				},
			},

			wantStackTop: value.True,
		},
		"put nil": {
			chunk: &vm.BytecodeMethod{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.String("foo"),
				},
			},

			wantStackTop: value.False,
		},
		"bool not int": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(0),
				},
			},

			wantStackTop: value.False,
		},
		"bool not true": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.TRUE),
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
			},

			wantStackTop: value.False,
		},
		"bool not false": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.FALSE),
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(0),
				},
			},

			wantStackTop: value.True,
		},
		"bool not nil": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.NIL),
					byte(bytecode.NOT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(5),
					0x1: value.Int8(25),
				},
			},

			wantStackTop: value.Int8(30),
		},
		"add String to String": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.String("foo"),
					0x1: value.String("bar"),
				},
			},

			wantStackTop: value.String("foobar"),
		},
		"add String to Char": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Char('f'),
					0x1: value.String("oo"),
				},
			},

			wantStackTop: value.String("foo"),
		},
		"add Int8 to String": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.ADD),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(5),
					0x1: value.Int8(25),
				},
			},

			wantStackTop: value.Int8(-20),
		},
		"String - String": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.SUBTRACT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(5),
					0x1: value.Int8(25),
				},
			},

			wantStackTop: value.Int8(125),
		},
		"String * SmallInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.String("foo"),
					0x1: value.SmallInt(3),
				},
			},

			wantStackTop: value.String("foofoofoo"),
		},
		"Char * SmallInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Char('a'),
					0x1: value.SmallInt(3),
				},
			},

			wantStackTop: value.String("aaa"),
		},
		"BigFloat * Float": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.MULTIPLY),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(35),
					0x1: value.Int8(5),
				},
			},

			wantStackTop: value.Int8(7),
		},
		"String / SmallInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.String("foo"),
					0x1: value.SmallInt(3),
				},
			},

			wantStackTop: value.String("foo"),
			wantErr:      value.NewNoMethodError("/", value.String("foo")),
		},
		"BigFloat / Float": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.DIVIDE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(2),
					0x1: value.Int8(5),
				},
			},

			wantStackTop: value.Int8(32),
		},
		"String ** SmallInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.String("foo"),
					0x1: value.SmallInt(3),
				},
			},

			wantStackTop: value.String("foo"),
			wantErr:      value.NewNoMethodError("**", value.String("foo")),
		},
		"BigFloat ** Float": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.EXPONENTIATE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.MODULO),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(25),
					0x1: value.Int8(4),
				},
			},

			wantStackTop: value.Int8(1),
		},
		"BigFloat % Float": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.MODULO),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(35),
					0x1: value.Int64(1),
				},
			},

			wantStackTop: value.Int8(17),
		},
		"Int >> Int": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.SmallInt(12),
					0x1: value.SmallInt(2),
				},
			},

			wantStackTop: value.SmallInt(3),
		},
		"-Int16 >> UInt32": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(35),
					0x1: value.Int64(1),
				},
			},

			wantStackTop: value.Int8(70),
		},
		"Int << Int": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.SmallInt(12),
					0x1: value.SmallInt(2),
				},
			},

			wantStackTop: value.SmallInt(48),
		},
		"-Int16 << UInt32": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(35),
					0x1: value.Int64(1),
				},
			},

			wantStackTop: value.Int8(17),
		},
		"-Int16 >>> UInt32": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.LOGIC_RBITSHIFT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(35),
					0x1: value.Int64(1),
				},
			},

			wantStackTop: value.Int8(70),
		},
		"UInt16 <<< Int": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.UInt16(12),
					0x1: value.SmallInt(2),
				},
			},

			wantStackTop: value.UInt16(48),
		},
		"-Int16 <<< UInt32": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.LOGIC_LBITSHIFT),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(23),
					0x1: value.Int8(10),
				},
			},

			wantStackTop: value.Int8(2),
		},
		"UInt16 & UInt16": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.UInt16(235),
					0x1: value.UInt16(58),
				},
			},

			wantStackTop: value.UInt16(42),
		},
		"SmallInt & SmallInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.SmallInt(235),
					0x1: value.SmallInt(58),
				},
			},

			wantStackTop: value.SmallInt(42),
		},
		"SmallInt & BigInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.SmallInt(255),
					0x1: value.ParseBigIntPanic("9223372036857247042", 10),
				},
			},

			wantStackTop: value.SmallInt(66),
		},
		"BigInt & SmallInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.ParseBigIntPanic("9223372036857247042", 10),
					0x1: value.SmallInt(255),
				},
			},

			wantStackTop: value.SmallInt(66),
		},
		"BigInt & BigInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_AND),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(23),
					0x1: value.Int8(10),
				},
			},

			wantStackTop: value.Int8(31),
		},
		"UInt16 | UInt16": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.UInt16(235),
					0x1: value.UInt16(58),
				},
			},

			wantStackTop: value.UInt16(251),
		},
		"SmallInt | SmallInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.SmallInt(235),
					0x1: value.SmallInt(58),
				},
			},

			wantStackTop: value.SmallInt(251),
		},
		"SmallInt | BigInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.SmallInt(255),
					0x1: value.ParseBigIntPanic("9223372036857247042", 10),
				},
			},

			wantStackTop: value.ParseBigIntPanic("9223372036857247231", 10),
		},
		"BigInt | SmallInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.ParseBigIntPanic("9223372036857247042", 10),
					0x1: value.SmallInt(255),
				},
			},

			wantStackTop: value.ParseBigIntPanic("9223372036857247231", 10),
		},
		"BigInt | BigInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_OR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.Int8(23),
					0x1: value.Int8(10),
				},
			},

			wantStackTop: value.Int8(29),
		},
		"UInt16 ^ UInt16": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.UInt16(235),
					0x1: value.UInt16(58),
				},
			},

			wantStackTop: value.UInt16(209),
		},
		"SmallInt ^ SmallInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.SmallInt(235),
					0x1: value.SmallInt(58),
				},
			},

			wantStackTop: value.SmallInt(209),
		},
		"SmallInt ^ BigInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.SmallInt(255),
					0x1: value.ParseBigIntPanic("9223372036857247042", 10),
				},
			},

			wantStackTop: value.ParseBigIntPanic("9223372036857247165", 10),
		},
		"BigInt ^ SmallInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					0x0: value.ParseBigIntPanic("9223372036857247042", 10),
					0x1: value.SmallInt(255),
				},
			},

			wantStackTop: value.ParseBigIntPanic("9223372036857247165", 10),
		},
		"BigInt ^ BigInt": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0x0,
					byte(bytecode.LOAD_VALUE8), 0x1,
					byte(bytecode.BITWISE_XOR),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.ToSymbol("Std"),
				},
			},

			wantStackTop: value.StdModule,
		},
		"get nested constants": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.GET_MOD_CONST8), 2,
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.ToSymbol("Std"),
					value.ToSymbol("Float"),
					value.ToSymbol("INF"),
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.DEF_MOD_CONST8), 1,
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.String("constant!"),
					value.ToSymbol("Foo"),
				},
			},
			teardown:     func() { value.RootModule.Constants.DeleteString("Foo") },
			wantStackTop: value.String("constant!"),
		},
		"define constant under Root and read it": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.DEF_MOD_CONST8), 1,
					byte(bytecode.POP),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.String("constant!"),
					value.ToSymbol("Foo"),
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.ToSymbol("Foo"),
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.GET_MOD_CONST8), 2,
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.ToSymbol("Foo"),
					value.ToSymbol("Std"),
					value.ToSymbol("Error"),
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
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_CLASS),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					&vm.BytecodeMethod{
						Instructions: []byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.DEF_MOD_CONST8), 1,
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						Values: []value.Value{
							value.SmallInt(1),
							value.ToSymbol("Bar"),
						},
						LineInfoList: bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
						},
						Location: L(P(5, 2, 5), P(44, 5, 7)),
					},
					value.ToSymbol("Foo"),
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
					value.SymbolMap{
						value.ToSymbol("Bar"): value.SmallInt(1),
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

func TestVM_DefAnonClass(t *testing.T) {
	tests := bytecodeTestTable{
		"define class without a body or superclass": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_CLASS),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 4),
				},
				Location: L(P(0, 1, 1), P(13, 1, 14)),
			},
			wantStackTop: value.NewClass(),
		},
		"define class with a superclass without a body": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.ROOT),
					byte(bytecode.GET_MOD_CONST8), 0,
					byte(bytecode.GET_MOD_CONST8), 1,
					byte(bytecode.DEF_ANON_CLASS),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.ToSymbol("Std"),
					value.ToSymbol("Error"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				Location: L(P(0, 1, 1), P(13, 1, 14)),
			},
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithParent(value.ErrorClass),
			),
		},
		"define class with a body": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_CLASS),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					&vm.BytecodeMethod{
						Instructions: []byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.DEF_MOD_CONST8), 1,
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						Values: []value.Value{
							value.SmallInt(1),
							value.ToSymbol("Bar"),
						},
						LineInfoList: bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
						},
						Location: L(P(5, 2, 5), P(44, 5, 7)),
					},
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				Location: L(P(0, 1, 1), P(45, 5, 8)),
			},
			wantStackTop: value.NewClassWithOptions(
				value.ClassWithConstants(
					value.SymbolMap{
						value.ToSymbol("Bar"): value.SmallInt(1),
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

func TestVM_DefModule(t *testing.T) {
	tests := bytecodeTestTable{
		"define module without a body": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.ToSymbol("Foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				Location: L(P(0, 1, 1), P(13, 1, 14)),
			},
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithName("Foo"),
			),
		},
		"define module with a body": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MODULE),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					&vm.BytecodeMethod{
						Instructions: []byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.DEF_MOD_CONST8), 1,
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						Values: []value.Value{
							value.SmallInt(1),
							value.ToSymbol("Bar"),
						},
						LineInfoList: bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
						},
						Location: L(P(5, 2, 5), P(44, 5, 7)),
					},
					value.ToSymbol("Foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				Location: L(P(0, 1, 1), P(45, 5, 8)),
			},
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithName("Foo"),
				value.ModuleWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithParent(value.ModuleClass),
					),
				),
				value.ModuleWithConstants(
					value.SymbolMap{
						value.ToSymbol("Bar"): value.SmallInt(1),
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

func TestVM_DefAnonModule(t *testing.T) {
	tests := bytecodeTestTable{
		"define module without a body": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_MODULE),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(13, 1, 14)),
			},
			wantStackTop: value.NewModule(),
		},
		"define module with a body": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_ANON_MODULE),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(45, 5, 8)),
				Values: []value.Value{
					&vm.BytecodeMethod{
						Instructions: []byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.DEF_MOD_CONST8), 1,
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						Values: []value.Value{
							value.SmallInt(1),
							value.ToSymbol("Bar"),
						},
						LineInfoList: bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
						},
						Location: L(P(5, 2, 5), P(44, 5, 7)),
					},
				},
			},
			wantStackTop: value.NewModuleWithOptions(
				value.ModuleWithClass(
					value.NewClassWithOptions(
						value.ClassWithSingleton(),
						value.ClassWithParent(value.ModuleClass),
					),
				),
				value.ModuleWithConstants(
					value.SymbolMap{
						value.ToSymbol("Bar"): value.SmallInt(1),
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

func TestVM_DefMethod(t *testing.T) {
	tests := bytecodeTestTable{
		"define a method": {
			chunk: vm.NewBytecodeMethod(
				value.ToSymbol("main"),
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
				nil,
				0,
				-1,
				false,
				false,
				[]value.Value{
					vm.NewBytecodeMethod(
						value.ToSymbol("foo"),
						[]byte{
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.RETURN),
						},
						L(P(5, 2, 5), P(21, 2, 21)),
						bytecode.LineInfoList{
							bytecode.NewLineInfo(2, 2),
						},
						nil,
						0,
						-1,
						false,
						false,
						[]value.Value{
							value.ToSymbol("bar"),
						},
					),
					value.ToSymbol("foo"),
				},
			),
			wantStackTop: vm.NewBytecodeMethod(
				value.ToSymbol("foo"),
				[]byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.RETURN),
				},
				L(P(5, 2, 5), P(21, 2, 21)),
				bytecode.LineInfoList{
					bytecode.NewLineInfo(2, 2),
				},
				nil,
				0,
				-1,
				false,
				false,
				[]value.Value{
					value.ToSymbol("bar"),
				},
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}

func TestVM_DefMixin(t *testing.T) {
	tests := bytecodeTestTable{
		"define mixin without a body": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.ToSymbol("Foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				Location: L(P(0, 1, 1), P(13, 1, 14)),
			},
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
			),
		},
		"define module with a body": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.CONSTANT_CONTAINER),
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.DEF_MIXIN),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					&vm.BytecodeMethod{
						Instructions: []byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.DEF_MOD_CONST8), 1,
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						Values: []value.Value{
							value.SmallInt(1),
							value.ToSymbol("Bar"),
						},
						LineInfoList: bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
						},
						Location: L(P(5, 2, 5), P(44, 5, 7)),
					},
					value.ToSymbol("Foo"),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 6),
				},
				Location: L(P(0, 1, 1), P(45, 5, 8)),
			},
			teardown: func() { value.RootModule.Constants.DeleteString("Foo") },
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithName("Foo"),
				value.MixinWithConstants(
					value.SymbolMap{
						value.ToSymbol("Bar"): value.SmallInt(1),
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

func TestVM_DefAnonMixin(t *testing.T) {
	tests := bytecodeTestTable{
		"define mixin without a body": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.UNDEFINED),
					byte(bytecode.DEF_ANON_MIXIN),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(13, 1, 14)),
			},
			wantStackTop: value.NewMixin(),
		},
		"define mixin with a body": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.DEF_ANON_MIXIN),
					byte(bytecode.RETURN),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(45, 5, 8)),
				Values: []value.Value{
					&vm.BytecodeMethod{
						Instructions: []byte{
							byte(bytecode.PREP_LOCALS8), 1,
							byte(bytecode.LOAD_VALUE8), 0,
							byte(bytecode.CONSTANT_CONTAINER),
							byte(bytecode.DEF_MOD_CONST8), 1,
							byte(bytecode.POP),
							byte(bytecode.SELF),
							byte(bytecode.RETURN),
						},
						Values: []value.Value{
							value.SmallInt(1),
							value.ToSymbol("Bar"),
						},
						LineInfoList: bytecode.LineInfoList{
							bytecode.NewLineInfo(1, 10),
						},
						Location: L(P(5, 2, 5), P(44, 5, 7)),
					},
				},
			},
			wantStackTop: value.NewMixinWithOptions(
				value.MixinWithConstants(
					value.SymbolMap{
						value.ToSymbol("Bar"): value.SmallInt(1),
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

func TestVM_Include(t *testing.T) {
	tests := bytecodeTestTable{
		"include a mixin in a class": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.INCLUDE),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.NewMixin(),
					value.ObjectClass,
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(13, 1, 14)),
			},
			teardown:     func() { value.ObjectClass.Parent = value.PrimitiveObjectClass },
			wantStackTop: value.Nil,
		},
		"include a mixin in an Int": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.INCLUDE),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.NewMixin(),
					value.SmallInt(4),
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(13, 1, 14)),
			},
			wantErr: value.NewError(
				value.TypeErrorClass,
				"can't include into an instance of Std::Int: `4`",
			),
		},
		"include a non mixin value": {
			chunk: &vm.BytecodeMethod{
				Instructions: []byte{
					byte(bytecode.LOAD_VALUE8), 0,
					byte(bytecode.LOAD_VALUE8), 1,
					byte(bytecode.INCLUDE),
					byte(bytecode.NIL),
					byte(bytecode.RETURN),
				},
				Values: []value.Value{
					value.String("foo"),
					value.ObjectClass,
				},
				LineInfoList: bytecode.LineInfoList{
					bytecode.NewLineInfo(1, 3),
				},
				Location: L(P(0, 1, 1), P(13, 1, 14)),
			},
			wantErr: value.NewError(
				value.TypeErrorClass,
				"`\"foo\"` is not a mixin",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmBytecodeTest(tc, t)
		})
	}
}
