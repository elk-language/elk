package vm_test

import (
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
