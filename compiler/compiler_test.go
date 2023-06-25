package compiler

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/object"
	"github.com/elk-language/elk/position"
	"github.com/google/go-cmp/cmp"
)

func TestCompilerCompileSource(t *testing.T) {
	tests := map[string]struct {
		input string
		want  *bytecode.Chunk
		err   position.ErrorList
	}{
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := CompileSource("main", []byte(tc.input))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.err, err); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
