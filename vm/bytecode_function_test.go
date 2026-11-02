package vm_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

func TestBytecodeFunction_AddInstruction(t *testing.T) {
	c := &vm.BytecodeFunction{}
	c.AddInstruction(1, bytecode.RETURN)
	want := &vm.BytecodeFunction{
		Instructions: []byte{byte(bytecode.RETURN)},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	if diff := cmp.Diff(want, c, comparer.Options()...); diff != "" {
		t.Fatal(diff)
	}

	c = &vm.BytecodeFunction{}
	c.AddInstruction(1, bytecode.LOAD_VALUE8, 0x12)
	want = &vm.BytecodeFunction{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 2)},
	}
	if diff := cmp.Diff(want, c, comparer.Options()); diff != "" {
		t.Fatal(diff)
	}

	c = &vm.BytecodeFunction{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	c.AddInstruction(1, bytecode.RETURN)
	want = &vm.BytecodeFunction{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12, byte(bytecode.RETURN)},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 2)},
	}
	if diff := cmp.Diff(want, c, comparer.Options()); diff != "" {
		t.Fatal(diff)
	}

	c = &vm.BytecodeFunction{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	c.AddInstruction(2, bytecode.RETURN)
	want = &vm.BytecodeFunction{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12, byte(bytecode.RETURN)},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1), bytecode.NewLineInfo(2, 1)},
	}
	if diff := cmp.Diff(want, c, comparer.Options()); diff != "" {
		t.Fatal(diff)
	}
}

func TestBytecodeFunction_AddConstant(t *testing.T) {
	tests := map[string]struct {
		chunkBefore *vm.BytecodeFunction
		add         value.Value
		wantInt     int
		wantSize    vm.IntSize
		chunkAfter  *vm.BytecodeFunction
	}{
		"add to an empty value pool": {
			chunkBefore: &vm.BytecodeFunction{
				Values: []value.Value{},
			},
			add:      value.Float(2.3).ToValue(),
			wantInt:  0,
			wantSize: bytecode.UINT8_SIZE,
			chunkAfter: &vm.BytecodeFunction{
				Values: []value.Value{value.Float(2.3).ToValue()},
			},
		},
		"add to a value pool with 255 elements": {
			chunkBefore: &vm.BytecodeFunction{
				Values: []value.Value{255: value.Nil},
			},
			add:      value.Float(2.3).ToValue(),
			wantInt:  256,
			wantSize: bytecode.UINT16_SIZE,
			chunkAfter: &vm.BytecodeFunction{
				Values: []value.Value{
					255: value.Nil,
					256: value.Float(2.3).ToValue(),
				},
			},
		},
		"add to a value pool with 65535 elements": {
			chunkBefore: &vm.BytecodeFunction{
				Values: []value.Value{65535: value.Nil},
			},
			add:      value.Float(2.3).ToValue(),
			wantInt:  65536,
			wantSize: bytecode.UINT32_SIZE,
			chunkAfter: &vm.BytecodeFunction{
				Values: []value.Value{
					65535: value.Nil,
					65536: value.Float(2.3).ToValue(),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotInt, gotSize := tc.chunkBefore.AddValue(tc.add)
			if diff := cmp.Diff(tc.wantInt, gotInt); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.wantSize, gotSize); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.chunkAfter, tc.chunkBefore, comparer.Options()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBytecodeFunction_AllInstructions(t *testing.T) {
	t.Run("matches disassembly offsets for every fixed opcode", func(t *testing.T) {
		for op := bytecode.NOOP; op <= bytecode.EXEC_DEFER; op++ {
			if op == bytecode.CLOSURE || op == bytecode.CLOSED_CLOSURE {
				continue
			}
			fn := &vm.BytecodeFunction{
				Instructions: []byte{byte(op), 0, 0, 0},
				Values:       []value.Value{value.Nil, value.Nil, value.Nil, value.Nil},
			}
			var gotOp bytecode.OpCode
			var gotOperands []byte
			var count int
			fn.AllInstructions()(func(opcode bytecode.OpCode, operands []byte) bool {
				gotOp = opcode
				gotOperands = operands
				count++
				return false
			})

			var buf strings.Builder
			next, err := fn.DisassembleInstruction(&buf, 0)
			if err != nil {
				t.Fatalf("%s: %v", op, err)
			}
			if count != 1 || gotOp != op || !bytes.Equal(gotOperands, fn.Instructions[1:next]) {
				t.Fatalf("%s: count=%d opcode=%s operands=%v want %v", op, count, gotOp, gotOperands, fn.Instructions[1:next])
			}
		}
	})

	t.Run("walks a closure and the following instruction", func(t *testing.T) {
		fn := &vm.BytecodeFunction{
			Instructions: []byte{
				byte(bytecode.CLOSURE),
				byte(vm.UpvalueLocalFlag),
				0x05,
				byte(vm.UpvalueLongIndexFlag),
				0x01, 0x02,
				vm.ClosureTerminatorFlag,
				byte(bytecode.RETURN),
			},
		}

		var got []bytecode.OpCode
		var operandLens []int
		for opcode, operands := range fn.AllInstructions() {
			got = append(got, opcode)
			operandLens = append(operandLens, len(operands))
			if opcode == bytecode.CLOSURE && !bytes.Equal(operands, fn.Instructions[1:7]) {
				t.Fatalf("closure operands: %v", operands)
			}
		}
		if diff := cmp.Diff([]bytecode.OpCode{bytecode.CLOSURE, bytecode.RETURN}, got); diff != "" {
			t.Fatal(diff)
		}
		if diff := cmp.Diff([]int{6, 0}, operandLens); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("stops when the consumer stops", func(t *testing.T) {
		fn := &vm.BytecodeFunction{
			Instructions: []byte{byte(bytecode.NOOP), byte(bytecode.RETURN)},
		}
		var count int
		for range fn.AllInstructions() {
			count++
			break
		}
		if count != 1 {
			t.Fatalf("yielded %d instructions", count)
		}
	})

	t.Run("panics on a truncated instruction", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic")
			}
		}()
		fn := &vm.BytecodeFunction{Instructions: []byte{byte(bytecode.LOAD_INT_16)}}
		for range fn.AllInstructions() {
		}
	})

	t.Run("panics on an unknown opcode", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic")
			}
		}()
		fn := &vm.BytecodeFunction{Instructions: []byte{0xff}}
		for range fn.AllInstructions() {
		}
	})
}
