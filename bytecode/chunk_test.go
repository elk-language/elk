package bytecode

import (
	"strings"
	"testing"

	"github.com/elk-language/elk/object"
	"github.com/elk-language/elk/position"
	"github.com/google/go-cmp/cmp"
)

func TestChunkAddInstruction(t *testing.T) {
	c := &Chunk{}
	c.AddInstruction(1, RETURN)
	want := &Chunk{
		Instructions: []byte{byte(RETURN)},
		LineInfoList: LineInfoList{NewLineInfo(1, 1)},
	}
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf(diff)
	}

	c = &Chunk{}
	c.AddInstruction(1, CONSTANT8, 0x12)
	want = &Chunk{
		Instructions: []byte{byte(CONSTANT8), 0x12},
		LineInfoList: LineInfoList{NewLineInfo(1, 1)},
	}
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf(diff)
	}

	c = &Chunk{
		Instructions: []byte{byte(CONSTANT8), 0x12},
		LineInfoList: LineInfoList{NewLineInfo(1, 1)},
	}
	c.AddInstruction(1, RETURN)
	want = &Chunk{
		Instructions: []byte{byte(CONSTANT8), 0x12, byte(RETURN)},
		LineInfoList: LineInfoList{NewLineInfo(1, 2)},
	}
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf(diff)
	}

	c = &Chunk{
		Instructions: []byte{byte(CONSTANT8), 0x12},
		LineInfoList: LineInfoList{NewLineInfo(1, 1)},
	}
	c.AddInstruction(2, RETURN)
	want = &Chunk{
		Instructions: []byte{byte(CONSTANT8), 0x12, byte(RETURN)},
		LineInfoList: LineInfoList{NewLineInfo(1, 1), NewLineInfo(2, 1)},
	}
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf(diff)
	}
}

func TestChunkAddConstant(t *testing.T) {
	tests := map[string]struct {
		chunkBefore *Chunk
		add         object.Value
		wantInt     int
		wantSize    IntSize
		chunkAfter  *Chunk
	}{
		"add to an empty constant pool": {
			chunkBefore: &Chunk{
				Constants: []object.Value{},
			},
			add:      object.Float(2.3),
			wantInt:  0,
			wantSize: UINT8_SIZE,
			chunkAfter: &Chunk{
				Constants: []object.Value{object.Float(2.3)},
			},
		},
		"add to a constant pool with 255 elements": {
			chunkBefore: &Chunk{
				Constants: []object.Value{255: object.Nil},
			},
			add:      object.Float(2.3),
			wantInt:  256,
			wantSize: UINT16_SIZE,
			chunkAfter: &Chunk{
				Constants: []object.Value{
					255: object.Nil,
					256: object.Float(2.3),
				},
			},
		},
		"add to a constant pool with 65535 elements": {
			chunkBefore: &Chunk{
				Constants: []object.Value{65535: object.Nil},
			},
			add:      object.Float(2.3),
			wantInt:  65536,
			wantSize: UINT32_SIZE,
			chunkAfter: &Chunk{
				Constants: []object.Value{
					65535: object.Nil,
					65536: object.Float(2.3),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotInt, gotSize := tc.chunkBefore.AddConstant(tc.add)
			if diff := cmp.Diff(tc.wantInt, gotInt); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.wantSize, gotSize); diff != "" {
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.chunkAfter, tc.chunkBefore); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestChunkDisassemble(t *testing.T) {
	tests := map[string]struct {
		in   *Chunk
		want string
		err  string
	}{
		"handle invalid opcodes": {
			in: &Chunk{
				Instructions: []byte{255},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       FF             unknown operation 255 (0xFF)
`,
			err: "unknown operation 255 (0xFF) at offset 0 (0x0)",
		},
		"correctly format the RETURN instruction": {
			in: &Chunk{
				Instructions: []byte{byte(RETURN)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       00             RETURN
`,
		},
		"correctly format the CONSTANT8 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT8), 0},
				Constants:    []object.Value{object.SmallInt(4)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       01 00          CONSTANT8       4
`,
		},
		"handle invalid CONSTANT8 index": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT8), 25},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Constants:    []object.Value{object.SmallInt(4)},
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       01 19          CONSTANT8       invalid constant index 25 (0x19)
`,
			err: "invalid constant index 25 (0x19)",
		},
		"handle missing bytes in CONSTANT8": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT8)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       01             CONSTANT8       not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the CONSTANT16 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT16), 0x01, 0x00},
				Constants:    []object.Value{0x1_00: object.SmallInt(4)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       02 01 00       CONSTANT16      4
`,
		},
		"handle invalid CONSTANT16 index": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT16), 0x19, 0xff},
				Constants:    []object.Value{object.SmallInt(4)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       02 19 FF       CONSTANT16      invalid constant index 6655 (0x19FF)
`,
			err: "invalid constant index 6655 (0x19FF)",
		},
		"handle missing bytes in CONSTANT16": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT16)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       02             CONSTANT16      not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the CONSTANT32 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT32), 0x01, 0x00, 0x00, 0x00},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Constants:    []object.Value{0x1_00_00_00: object.SmallInt(4)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       03 01 00 00 00 CONSTANT32      4
`,
		},
		"handle invalid CONSTANT32 index": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT32), 0x01, 0x00, 0x00, 0x00},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       03 01 00 00 00 CONSTANT32      invalid constant index 16777216 (0x1000000)
`,
			err: "invalid constant index 16777216 (0x1000000)",
		},
		"handle missing bytes in CONSTANT32": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT32)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       03             CONSTANT32      not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the ADD opcode": {
			in: &Chunk{
				Instructions: []byte{byte(ADD)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       04             ADD
`,
		},
		"correctly format the SUBTRACT opcode": {
			in: &Chunk{
				Instructions: []byte{byte(SUBTRACT)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       05             SUBTRACT
`,
		},
		"correctly format the MULTIPLY opcode": {
			in: &Chunk{
				Instructions: []byte{byte(MULTIPLY)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       06             MULTIPLY
`,
		},
		"correctly format the DIVIDE opcode": {
			in: &Chunk{
				Instructions: []byte{byte(DIVIDE)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       07             DIVIDE
`,
		},
		"correctly format the EXPONENTIATE opcode": {
			in: &Chunk{
				Instructions: []byte{byte(EXPONENTIATE)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       08             EXPONENTIATE
`,
		},
		"correctly format the NEGATE opcode": {
			in: &Chunk{
				Instructions: []byte{byte(NEGATE)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       09             NEGATE
`,
		},
		"correctly format the NOT opcode": {
			in: &Chunk{
				Instructions: []byte{byte(NOT)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       0A             NOT
`,
		},
		"correctly format the BITWISE_NOT opcode": {
			in: &Chunk{
				Instructions: []byte{byte(BITWISE_NOT)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       0B             BITWISE_NOT
`,
		},
		"correctly format the TRUE opcode": {
			in: &Chunk{
				Instructions: []byte{byte(TRUE)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       0C             TRUE
`,
		},
		"correctly format the FALSE opcode": {
			in: &Chunk{
				Instructions: []byte{byte(FALSE)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       0D             FALSE
`,
		},
		"correctly format the NIL opcode": {
			in: &Chunk{
				Instructions: []byte{byte(NIL)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     position.NewLocation("/foo/bar.elk", 12, 6, 2, 3),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       0E             NIL
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
