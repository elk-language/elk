package vm

// import (
// 	"strings"
// 	"testing"

// 	"github.com/elk-language/elk/bytecode"
// 	"github.com/elk-language/elk/object"
// 	"github.com/google/go-cmp/cmp"
// )

// func TestVMInterpretBytecode(t *testing.T) {
// 	tests := map[string]struct {
// 		chunk      *bytecode.Chunk
// 		wantResult Result
// 		wantStdout string
// 	}{
// 		"add to an empty constant pool": {
// 			chunk: &bytecode.Chunk{
// 				Instructions: []byte{
// 					byte(bytecode.CONSTANT8),
// 					0x0,
// 					byte(bytecode.RETURN),
// 				},
// 				Constants: []object.Value{
// 					0x0: object.Int8(5),
// 				},
// 			},
// 			wantResult: RESULT_OK,
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			var stdout strings.Builder
// 			vm := New(WithStdout(&stdout))
// 			gotResult := vm.InterpretBytecode(tc.chunk)
// 			gotStdout := stdout.String()
// 			if diff := cmp.Diff(tc.wantResult, gotResult); diff != "" {
// 				t.Fatalf(diff)
// 			}
// 			if diff := cmp.Diff(tc.wantStdout, gotStdout); diff != "" {
// 				t.Fatalf(diff)
// 			}
// 		})
// 	}
// }
