package vm

import (
	"strings"
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/object"
	"github.com/google/go-cmp/cmp"
)

func TestVMInterpretBytecode(t *testing.T) {
	tests := map[string]struct {
		chunk        *bytecode.Chunk
		wantStackTop object.Value
		wantResult   bool
		wantStdout   string
	}{
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
			wantResult:   true,
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
			wantResult:   true,
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
			wantResult:   true,
			wantStackTop: object.Int8(5),
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
			wantResult:   true,
			wantStackTop: object.Int8(-5),
		},
		"add Int8": {
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
			wantResult:   true,
			wantStackTop: object.Int8(30),
		},
		"add String": {
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
			wantResult:   true,
			wantStackTop: object.String("foobar"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var stdout strings.Builder
			vm := New(WithStdout(&stdout))
			gotResult := vm.InterpretBytecode(tc.chunk)
			gotStdout := stdout.String()
			gotStackTop := vm.peek()
			if diff := cmp.Diff(tc.wantResult, gotResult); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.wantStdout, gotStdout); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.wantStackTop, gotStackTop); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
