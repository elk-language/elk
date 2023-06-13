package bytecode

import (
	"strings"
	"testing"

	"github.com/elk-language/elk/position"
	"github.com/google/go-cmp/cmp"
)

func TestChunkAddInstruction(t *testing.T) {
	c := NewChunk(nil, nil)
	c.AddInstruction(RETURN)
	want := NewChunk([]byte{byte(RETURN)}, nil)
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf(diff)
	}

	c = NewChunk(nil, nil)
	c.AddInstruction(CONSTANT, 0x12)
	want = NewChunk([]byte{byte(CONSTANT), 0x12}, nil)
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf(diff)
	}

	c = NewChunk([]byte{byte(CONSTANT), 0x12}, nil)
	c.AddInstruction(RETURN)
	want = NewChunk([]byte{byte(CONSTANT), 0x12, byte(RETURN)}, nil)
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf(diff)
	}
}

func TestChunkDisassemble(t *testing.T) {
	tests := map[string]struct {
		in   *Chunk
		want string
		err  string
	}{
		"handle invalid opcodes": {
			in: NewChunk(
				[]byte{255},
				position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			),
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  FF          unknown operation 255 (0xFF)
`,
			err: "unknown operation 255 (0xFF) at offset 0 (0x0)",
		},
		"correctly format one byte instructions": {
			in: NewChunk(
				[]byte{byte(RETURN)},
				position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			),
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  00          RETURN
`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var buffer strings.Builder
			err := tc.in.Disassemble(&buffer)
			got := buffer.String()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			if diff := cmp.Diff(tc.err, gotErr); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
