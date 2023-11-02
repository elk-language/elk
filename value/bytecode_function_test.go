package value

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position"
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

func TestBytecodeFunctionAddInstruction(t *testing.T) {
	c := &BytecodeFunction{}
	c.AddInstruction(1, bytecode.RETURN)
	want := &BytecodeFunction{
		Instructions: []byte{byte(bytecode.RETURN)},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf(diff)
	}

	c = &BytecodeFunction{}
	c.AddInstruction(1, bytecode.LOAD_VALUE8, 0x12)
	want = &BytecodeFunction{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf(diff)
	}

	c = &BytecodeFunction{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	c.AddInstruction(1, bytecode.RETURN)
	want = &BytecodeFunction{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12, byte(bytecode.RETURN)},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 2)},
	}
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf(diff)
	}

	c = &BytecodeFunction{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	c.AddInstruction(2, bytecode.RETURN)
	want = &BytecodeFunction{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12, byte(bytecode.RETURN)},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1), bytecode.NewLineInfo(2, 1)},
	}
	if diff := cmp.Diff(want, c); diff != "" {
		t.Fatalf(diff)
	}
}

func TestBytecodeFunctionAddConstant(t *testing.T) {
	tests := map[string]struct {
		chunkBefore *BytecodeFunction
		add         Value
		wantInt     int
		wantSize    IntSize
		chunkAfter  *BytecodeFunction
	}{
		"add to an empty value pool": {
			chunkBefore: &BytecodeFunction{
				Values: []Value{},
			},
			add:      Float(2.3),
			wantInt:  0,
			wantSize: bytecode.UINT8_SIZE,
			chunkAfter: &BytecodeFunction{
				Values: []Value{Float(2.3)},
			},
		},
		"add to a value pool with 255 elements": {
			chunkBefore: &BytecodeFunction{
				Values: []Value{255: Nil},
			},
			add:      Float(2.3),
			wantInt:  256,
			wantSize: bytecode.UINT16_SIZE,
			chunkAfter: &BytecodeFunction{
				Values: []Value{
					255: Nil,
					256: Float(2.3),
				},
			},
		},
		"add to a value pool with 65535 elements": {
			chunkBefore: &BytecodeFunction{
				Values: []Value{65535: Nil},
			},
			add:      Float(2.3),
			wantInt:  65536,
			wantSize: bytecode.UINT32_SIZE,
			chunkAfter: &BytecodeFunction{
				Values: []Value{
					65535: Nil,
					65536: Float(2.3),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotInt, gotSize := tc.chunkBefore.AddValue(tc.add)
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

func TestBytecodeFunctionDisassemble(t *testing.T) {
	tests := map[string]struct {
		in   *BytecodeFunction
		want string
		err  string
	}{
		"handle invalid opcodes": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{255},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       FF             unknown operation 255 (0xFF)
`,
			err: "unknown operation 255 (0xFF) at offset 0 (0x0)",
		},
		"correctly format the RETURN instruction": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.RETURN)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       00             RETURN
`,
		},
		"correctly format the LOAD_VALUE8 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0},
				Values:       []Value{SmallInt(4)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       01 00          LOAD_VALUE8     4
`,
		},
		"handle invalid LOAD_VALUE8 index": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOAD_VALUE8), 25},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Values:       []Value{SmallInt(4)},
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       01 19          LOAD_VALUE8     invalid value index 25 (0x19)
`,
			err: "invalid value index 25 (0x19)",
		},
		"handle missing bytes in LOAD_VALUE8": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOAD_VALUE8)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       01             LOAD_VALUE8     not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the LOAD_VALUE16 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOAD_VALUE16), 0x01, 0x00},
				Values:       []Value{0x1_00: SmallInt(4)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       02 01 00       LOAD_VALUE16    4
`,
		},
		"handle invalid LOAD_VALUE16 index": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOAD_VALUE16), 0x19, 0xff},
				Values:       []Value{SmallInt(4)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       02 19 FF       LOAD_VALUE16    invalid value index 6655 (0x19FF)
`,
			err: "invalid value index 6655 (0x19FF)",
		},
		"handle missing bytes in LOAD_VALUE16": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOAD_VALUE16)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       02             LOAD_VALUE16    not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the LOAD_VALUE32 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOAD_VALUE32), 0x01, 0x00, 0x00, 0x00},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Values:       []Value{0x1_00_00_00: SmallInt(4)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       03 01 00 00 00 LOAD_VALUE32    4
`,
		},
		"handle invalid LOAD_VALUE32 index": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOAD_VALUE32), 0x01, 0x00, 0x00, 0x00},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       03 01 00 00 00 LOAD_VALUE32    invalid value index 16777216 (0x1000000)
`,
			err: "invalid value index 16777216 (0x1000000)",
		},
		"handle missing bytes in LOAD_VALUE32": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOAD_VALUE32)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       03             LOAD_VALUE32    not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the ADD opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.ADD)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       04             ADD
`,
		},
		"correctly format the SUBTRACT opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.SUBTRACT)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       05             SUBTRACT
`,
		},
		"correctly format the MULTIPLY opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.MULTIPLY)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       06             MULTIPLY
`,
		},
		"correctly format the DIVIDE opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.DIVIDE)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       07             DIVIDE
`,
		},
		"correctly format the EXPONENTIATE opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.EXPONENTIATE)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       08             EXPONENTIATE
`,
		},
		"correctly format the NEGATE opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.NEGATE)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       09             NEGATE
`,
		},
		"correctly format the NOT opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.NOT)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       0A             NOT
`,
		},
		"correctly format the BITWISE_NOT opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.BITWISE_NOT)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       0B             BITWISE_NOT
`,
		},
		"correctly format the TRUE opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.TRUE)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       0C             TRUE
`,
		},
		"correctly format the FALSE opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.FALSE)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       0D             FALSE
`,
		},
		"correctly format the NIL opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.NIL)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       0E             NIL
`,
		},
		"correctly format the POP opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.POP)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       0F             POP
`,
		},
		"correctly format the POP_N opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.POP_N), 3},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       10 03          POP_N           3               
`,
		},
		"handle missing bytes in POP_N": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.POP_N)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       10             POP_N           not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the LEAVE_SCOPE16 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LEAVE_SCOPE16), 3, 2},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       11 03 02       LEAVE_SCOPE16   3               2               
`,
		},
		"correctly format the LEAVE_SCOPE32 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LEAVE_SCOPE32), 3, 2, 0, 2},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       12 03 02 00 02 LEAVE_SCOPE32   770             2               
`,
		},
		"correctly format the PREP_LOCALS8 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.PREP_LOCALS8), 3},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       13 03          PREP_LOCALS8    3               
`,
		},
		"correctly format the PREP_LOCALS16 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.PREP_LOCALS16), 3, 5},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       14 03 05       PREP_LOCALS16   773             
`,
		},
		"correctly format the SET_LOCAL8 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.SET_LOCAL8), 3},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       15 03          SET_LOCAL8      3               
`,
		},
		"correctly format the SET_LOCAL16 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.SET_LOCAL16), 3, 2},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       16 03 02       SET_LOCAL16     770             
`,
		},
		"correctly format the GET_LOCAL8 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.GET_LOCAL8), 3},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       17 03          GET_LOCAL8      3               
`,
		},
		"correctly format the GET_LOCAL16 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.GET_LOCAL16), 3, 2},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       18 03 02       GET_LOCAL16     770             
`,
		},
		"correctly format the JUMP_UNLESS opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.JUMP_UNLESS), 3, 2},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       19 03 02       JUMP_UNLESS     770             
`,
		},
		"correctly format the JUMP opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.JUMP), 3, 2},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       1A 03 02       JUMP            770             
`,
		},
		"correctly format the JUMP_IF opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.JUMP_IF), 3, 2},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       1B 03 02       JUMP_IF         770             
`,
		},
		"correctly format the LOOP opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOOP), 3, 2},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       1C 03 02       LOOP            770             
`,
		},
		"correctly format the JUMP_IF_NIL opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.JUMP_IF_NIL), 3, 2},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       1D 03 02       JUMP_IF_NIL     770             
`,
		},
		"correctly format the RBITSHIFT opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.RBITSHIFT)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       1E             RBITSHIFT
`,
		},
		"correctly format the LOGIC_RBITSHIFT opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOGIC_RBITSHIFT)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       1F             LOGIC_RBITSHIFT
`,
		},
		"correctly format the LBITSHIFT opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LBITSHIFT)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       20             LBITSHIFT
`,
		},
		"correctly format the LOGIC_LBITSHIFT opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LOGIC_LBITSHIFT)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       21             LOGIC_LBITSHIFT
`,
		},
		"correctly format the BITWISE_AND opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.BITWISE_AND)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       22             BITWISE_AND
`,
		},
		"correctly format the BITWISE_OR opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.BITWISE_OR)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       23             BITWISE_OR
`,
		},
		"correctly format the BITWISE_XOR opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.BITWISE_XOR)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       24             BITWISE_XOR
`,
		},
		"correctly format the MODULO opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.MODULO)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       25             MODULO
`,
		},
		"correctly format the EQUAL opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.EQUAL)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       26             EQUAL
`,
		},
		"correctly format the STRICT_EQUAL opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.STRICT_EQUAL)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       27             STRICT_EQUAL
`,
		},
		"correctly format the GREATER opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.GREATER)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       28             GREATER
`,
		},
		"correctly format the GREATER_EQUAL opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.GREATER_EQUAL)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       29             GREATER_EQUAL
`,
		},
		"correctly format the LESS opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LESS)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       2A             LESS
`,
		},
		"correctly format the LESS_EQUAL opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.LESS_EQUAL)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       2B             LESS_EQUAL
`,
		},
		"correctly format the GET_MOD_CONST8 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.GET_MOD_CONST8), 0},
				Values:       []Value{0: SymbolTable.Add("Foo")},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       2C 00          GET_MOD_CONST8  :Foo
`,
		},
		"correctly format the GET_MOD_CONST16 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.GET_MOD_CONST16), 0x01, 0x00},
				Values:       []Value{0x1_00: SymbolTable.Add("Bar")},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       2D 01 00       GET_MOD_CONST16 :Bar
`,
		},
		"correctly format the GET_MOD_CONST32 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.GET_MOD_CONST32), 0x01, 0x00, 0x00, 0x00},
				Values:       []Value{0x1_00_00_00: SymbolTable.Add("Bar")},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       2E 01 00 00 00 GET_MOD_CONST32 :Bar
`,
		},
		"correctly format the ROOT opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.ROOT)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       2F             ROOT
`,
		},
		"correctly format the NOT_EQUAL opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.NOT_EQUAL)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       30             NOT_EQUAL
`,
		},
		"correctly format the STRICT_NOT_EQUAL opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.STRICT_NOT_EQUAL)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       31             STRICT_NOT_EQUAL
`,
		},
		"correctly format the DEF_MOD_CONST8 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.DEF_MOD_CONST8), 0},
				Values:       []Value{0: SymbolTable.Add("Foo")},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       32 00          DEF_MOD_CONST8  :Foo
`,
		},
		"correctly format the DEF_MOD_CONST16 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.DEF_MOD_CONST16), 0x01, 0x00},
				Values:       []Value{0x1_00: SymbolTable.Add("Bar")},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       33 01 00       DEF_MOD_CONST16 :Bar
`,
		},
		"correctly format the DEF_MOD_CONST32 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.DEF_MOD_CONST32), 0x01, 0x00, 0x00, 0x00},
				Values:       []Value{0x1_00_00_00: SymbolTable.Add("Bar")},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       34 01 00 00 00 DEF_MOD_CONST32 :Bar
`,
		},
		"correctly format the CONSTANT_BASE opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.CONSTANT_BASE)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       35             CONSTANT_BASE
`,
		},
		"correctly format the DEF_CLASS16 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.DEF_CLASS)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       36             DEF_CLASS
`,
		},
		"correctly format the SELF opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.SELF)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       37             SELF
`,
		},
		"correctly format the DEF_MODULE opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.DEF_MODULE)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       38             DEF_MODULE
`,
		},
		"correctly format the CALL_METHOD8 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.CALL_METHOD8), 0},
				Values:       []Value{0: NewCallSiteInfo(SymbolTable.Add("foo"), 0)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       39 00          CALL_METHOD8    CallSiteInfo{method_name: :foo, argument_count: 0}
`,
		},
		"correctly format the CALL_METHOD16 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.CALL_METHOD16), 0x01, 0x00},
				Values:       []Value{0x1_00: NewCallSiteInfo(SymbolTable.Add("foo"), 0)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       3A 01 00       CALL_METHOD16   CallSiteInfo{method_name: :foo, argument_count: 0}
`,
		},
		"correctly format the CALL_METHOD32 opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.CALL_METHOD32), 0x01, 0x00, 0x00, 0x00},
				Values:       []Value{0x1_00_00_00: NewCallSiteInfo(SymbolTable.Add("foo"), 0)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       3B 01 00 00 00 CALL_METHOD32   CallSiteInfo{method_name: :foo, argument_count: 0}
`,
		},
		"correctly format the DEF_METHOD opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.DEF_METHOD)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       3C             DEF_METHOD
`,
		},
		"correctly format the UNDEFINED opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.UNDEFINED)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       3D             UNDEFINED
`,
		},
		"correctly format the DEF_ANON_CLASS opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.DEF_ANON_CLASS)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       3E             DEF_ANON_CLASS
`,
		},
		"correctly format the DEF_ANON_MODULE opcode": {
			in: &BytecodeFunction{
				Name:         "main",
				Instructions: []byte{byte(bytecode.DEF_ANON_MODULE)},
				LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				Location:     L(P(12, 2, 3), P(18, 2, 9)),
			},
			want: `== Disassembly of main at: /foo/bar.elk:2:3 ==

0000  1       3F             DEF_ANON_MODULE
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
