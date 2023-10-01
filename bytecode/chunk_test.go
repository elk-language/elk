package bytecode

import (
	"testing"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

const testFileName = "/foo/bar.elk"

// Create a new position in tests
var P = position.New

// Create a new span in tests
var S = position.NewSpan

// Create a new location in tests
func L(startPos, endPos *position.Position) *position.Location {
	return position.NewLocation(testFileName, startPos, endPos)
}

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
		add         value.Value
		wantInt     int
		wantSize    IntSize
		chunkAfter  *Chunk
	}{
		"add to an empty constant pool": {
			chunkBefore: &Chunk{
				Constants: []value.Value{},
			},
			add:      value.Float(2.3),
			wantInt:  0,
			wantSize: UINT8_SIZE,
			chunkAfter: &Chunk{
				Constants: []value.Value{value.Float(2.3)},
			},
		},
		"add to a constant pool with 255 elements": {
			chunkBefore: &Chunk{
				Constants: []value.Value{255: value.Nil},
			},
			add:      value.Float(2.3),
			wantInt:  256,
			wantSize: UINT16_SIZE,
			chunkAfter: &Chunk{
				Constants: []value.Value{
					255: value.Nil,
					256: value.Float(2.3),
				},
			},
		},
		"add to a constant pool with 65535 elements": {
			chunkBefore: &Chunk{
				Constants: []value.Value{65535: value.Nil},
			},
			add:      value.Float(2.3),
			wantInt:  65536,
			wantSize: UINT32_SIZE,
			chunkAfter: &Chunk{
				Constants: []value.Value{
					65535: value.Nil,
					65536: value.Float(2.3),
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
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
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
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       00             RETURN
`,
		},
		"correctly format the CONSTANT8 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT8), 0},
				Constants:    []value.Value{value.SmallInt(4)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       01 00          CONSTANT8       4
`,
		},
		"handle invalid CONSTANT8 index": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT8), 25},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Constants:    []value.Value{value.SmallInt(4)},
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
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       01             CONSTANT8       not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the CONSTANT16 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT16), 0x01, 0x00},
				Constants:    []value.Value{0x1_00: value.SmallInt(4)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       02 01 00       CONSTANT16      4
`,
		},
		"handle invalid CONSTANT16 index": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT16), 0x19, 0xff},
				Constants:    []value.Value{value.SmallInt(4)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
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
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
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
				Constants:    []value.Value{0x1_00_00_00: value.SmallInt(4)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       03 01 00 00 00 CONSTANT32      4
`,
		},
		"handle invalid CONSTANT32 index": {
			in: &Chunk{
				Instructions: []byte{byte(CONSTANT32), 0x01, 0x00, 0x00, 0x00},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
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
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
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
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       04             ADD
`,
		},
		"correctly format the SUBTRACT opcode": {
			in: &Chunk{
				Instructions: []byte{byte(SUBTRACT)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       05             SUBTRACT
`,
		},
		"correctly format the MULTIPLY opcode": {
			in: &Chunk{
				Instructions: []byte{byte(MULTIPLY)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       06             MULTIPLY
`,
		},
		"correctly format the DIVIDE opcode": {
			in: &Chunk{
				Instructions: []byte{byte(DIVIDE)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       07             DIVIDE
`,
		},
		"correctly format the EXPONENTIATE opcode": {
			in: &Chunk{
				Instructions: []byte{byte(EXPONENTIATE)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       08             EXPONENTIATE
`,
		},
		"correctly format the NEGATE opcode": {
			in: &Chunk{
				Instructions: []byte{byte(NEGATE)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       09             NEGATE
`,
		},
		"correctly format the NOT opcode": {
			in: &Chunk{
				Instructions: []byte{byte(NOT)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       0A             NOT
`,
		},
		"correctly format the BITWISE_NOT opcode": {
			in: &Chunk{
				Instructions: []byte{byte(BITWISE_NOT)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       0B             BITWISE_NOT
`,
		},
		"correctly format the TRUE opcode": {
			in: &Chunk{
				Instructions: []byte{byte(TRUE)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       0C             TRUE
`,
		},
		"correctly format the FALSE opcode": {
			in: &Chunk{
				Instructions: []byte{byte(FALSE)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       0D             FALSE
`,
		},
		"correctly format the NIL opcode": {
			in: &Chunk{
				Instructions: []byte{byte(NIL)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       0E             NIL
`,
		},
		"correctly format the POP opcode": {
			in: &Chunk{
				Instructions: []byte{byte(POP)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       0F             POP
`,
		},
		"correctly format the POP_N opcode": {
			in: &Chunk{
				Instructions: []byte{byte(POP_N), 3},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       10 03          POP_N           3               
`,
		},
		"handle missing bytes in POP_N": {
			in: &Chunk{
				Instructions: []byte{byte(POP_N)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       10             POP_N           not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the LEAVE_SCOPE16 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(LEAVE_SCOPE16), 3, 2},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       11 03 02       LEAVE_SCOPE16   3               2               
`,
		},
		"correctly format the LEAVE_SCOPE32 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(LEAVE_SCOPE32), 3, 2, 0, 2},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       12 03 02 00 02 LEAVE_SCOPE32   770             2               
`,
		},
		"correctly format the PREP_LOCALS8 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(PREP_LOCALS8), 3},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       13 03          PREP_LOCALS8    3               
`,
		},
		"correctly format the PREP_LOCALS16 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(PREP_LOCALS16), 3, 5},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       14 03 05       PREP_LOCALS16   773             
`,
		},
		"correctly format the SET_LOCAL8 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(SET_LOCAL8), 3},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       15 03          SET_LOCAL8      3               
`,
		},
		"correctly format the SET_LOCAL16 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(SET_LOCAL16), 3, 2},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       16 03 02       SET_LOCAL16     770             
`,
		},
		"correctly format the GET_LOCAL8 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(GET_LOCAL8), 3},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       17 03          GET_LOCAL8      3               
`,
		},
		"correctly format the GET_LOCAL16 opcode": {
			in: &Chunk{
				Instructions: []byte{byte(GET_LOCAL16), 3, 2},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       18 03 02       GET_LOCAL16     770             
`,
		},
		"correctly format the JUMP_UNLESS opcode": {
			in: &Chunk{
				Instructions: []byte{byte(JUMP_UNLESS), 3, 2},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       19 03 02       JUMP_UNLESS     770             
`,
		},
		"correctly format the JUMP opcode": {
			in: &Chunk{
				Instructions: []byte{byte(JUMP), 3, 2},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       1A 03 02       JUMP            770             
`,
		},
		"correctly format the JUMP_IF opcode": {
			in: &Chunk{
				Instructions: []byte{byte(JUMP_IF), 3, 2},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       1B 03 02       JUMP_IF         770             
`,
		},
		"correctly format the LOOP opcode": {
			in: &Chunk{
				Instructions: []byte{byte(LOOP), 3, 2},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       1C 03 02       LOOP            770             
`,
		},
		"correctly format the JUMP_IF_NIL opcode": {
			in: &Chunk{
				Instructions: []byte{byte(JUMP_IF_NIL), 3, 2},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       1D 03 02       JUMP_IF_NIL     770             
`,
		},
		"correctly format the RBITSHIFT opcode": {
			in: &Chunk{
				Instructions: []byte{byte(RBITSHIFT)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       1E             RBITSHIFT
`,
		},
		"correctly format the LOGIC_RBITSHIFT opcode": {
			in: &Chunk{
				Instructions: []byte{byte(LOGIC_RBITSHIFT)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       1F             LOGIC_RBITSHIFT
`,
		},
		"correctly format the LBITSHIFT opcode": {
			in: &Chunk{
				Instructions: []byte{byte(LBITSHIFT)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       20             LBITSHIFT
`,
		},
		"correctly format the LOGIC_LBITSHIFT opcode": {
			in: &Chunk{
				Instructions: []byte{byte(LOGIC_LBITSHIFT)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       21             LOGIC_LBITSHIFT
`,
		},
		"correctly format the BITWISE_AND opcode": {
			in: &Chunk{
				Instructions: []byte{byte(BITWISE_AND)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       22             BITWISE_AND
`,
		},
		"correctly format the BITWISE_OR opcode": {
			in: &Chunk{
				Instructions: []byte{byte(BITWISE_OR)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       23             BITWISE_OR
`,
		},
		"correctly format the BITWISE_XOR opcode": {
			in: &Chunk{
				Instructions: []byte{byte(BITWISE_XOR)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       24             BITWISE_XOR
`,
		},
		"correctly format the MODULO opcode": {
			in: &Chunk{
				Instructions: []byte{byte(MODULO)},
				LineInfoList: LineInfoList{NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of bytecode chunk at: /foo/bar.elk:2:3 ==

0000  1       25             MODULO
`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.in.DisassembleString()
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
