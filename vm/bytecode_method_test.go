package vm_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

var mainSymbol = value.SymbolTable.Add("main")

func TestBytecodeMethod_AddInstruction(t *testing.T) {
	c := &vm.BytecodeMethod{}
	c.AddInstruction(1, bytecode.RETURN)
	want := &vm.BytecodeMethod{
		Instructions: []byte{byte(bytecode.RETURN)},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	if diff := cmp.Diff(want, c, comparer.Comparer...); diff != "" {
		t.Fatalf(diff)
	}

	c = &vm.BytecodeMethod{}
	c.AddInstruction(1, bytecode.LOAD_VALUE8, 0x12)
	want = &vm.BytecodeMethod{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	if diff := cmp.Diff(want, c, comparer.Comparer); diff != "" {
		t.Fatalf(diff)
	}

	c = &vm.BytecodeMethod{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	c.AddInstruction(1, bytecode.RETURN)
	want = &vm.BytecodeMethod{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12, byte(bytecode.RETURN)},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 2)},
	}
	if diff := cmp.Diff(want, c, comparer.Comparer); diff != "" {
		t.Fatalf(diff)
	}

	c = &vm.BytecodeMethod{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	c.AddInstruction(2, bytecode.RETURN)
	want = &vm.BytecodeMethod{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12, byte(bytecode.RETURN)},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1), bytecode.NewLineInfo(2, 1)},
	}
	if diff := cmp.Diff(want, c, comparer.Comparer); diff != "" {
		t.Fatalf(diff)
	}
}

func TestBytecodeMethod_AddConstant(t *testing.T) {
	tests := map[string]struct {
		chunkBefore *vm.BytecodeMethod
		add         value.Value
		wantInt     int
		wantSize    vm.IntSize
		chunkAfter  *vm.BytecodeMethod
	}{
		"add to an empty value pool": {
			chunkBefore: &vm.BytecodeMethod{
				Values: []value.Value{},
			},
			add:      value.Float(2.3),
			wantInt:  0,
			wantSize: bytecode.UINT8_SIZE,
			chunkAfter: &vm.BytecodeMethod{
				Values: []value.Value{value.Float(2.3)},
			},
		},
		"add to a value pool with 255 elements": {
			chunkBefore: &vm.BytecodeMethod{
				Values: []value.Value{255: value.Nil},
			},
			add:      value.Float(2.3),
			wantInt:  256,
			wantSize: bytecode.UINT16_SIZE,
			chunkAfter: &vm.BytecodeMethod{
				Values: []value.Value{
					255: value.Nil,
					256: value.Float(2.3),
				},
			},
		},
		"add to a value pool with 65535 elements": {
			chunkBefore: &vm.BytecodeMethod{
				Values: []value.Value{65535: value.Nil},
			},
			add:      value.Float(2.3),
			wantInt:  65536,
			wantSize: bytecode.UINT32_SIZE,
			chunkAfter: &vm.BytecodeMethod{
				Values: []value.Value{
					65535: value.Nil,
					65536: value.Float(2.3),
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
			if diff := cmp.Diff(tc.chunkAfter, tc.chunkBefore, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestBytecodeMethod_Disassemble(t *testing.T) {
	tests := map[string]struct {
		in   *vm.BytecodeMethod
		want string
		err  string
	}{
		"handle invalid opcodes": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{255},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       FF             unknown operation 255 (0xFF)
`,
			err: "unknown operation 255 (0xFF) at offset 0 (0x0)",
		},
		"correctly format the RETURN instruction": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.RETURN)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       00             RETURN
`,
		},
		"correctly format the LOAD_VALUE8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOAD_VALUE8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{value.SmallInt(4)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       01 00          LOAD_VALUE8       4
`,
		},
		"handle invalid LOAD_VALUE8 index": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOAD_VALUE8), 25},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{value.SmallInt(4)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       01 19          LOAD_VALUE8       invalid value index 25 (0x19)
`,
			err: "invalid value index 25 (0x19)",
		},
		"handle missing bytes in LOAD_VALUE8": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOAD_VALUE8)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       01             LOAD_VALUE8       not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the LOAD_VALUE16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOAD_VALUE16), 0x01, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00: value.SmallInt(4)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       02 01 00       LOAD_VALUE16      4
`,
		},
		"handle invalid LOAD_VALUE16 index": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOAD_VALUE16), 0x19, 0xff},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{value.SmallInt(4)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       02 19 FF       LOAD_VALUE16      invalid value index 6655 (0x19FF)
`,
			err: "invalid value index 6655 (0x19FF)",
		},
		"handle missing bytes in LOAD_VALUE16": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOAD_VALUE16)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       02             LOAD_VALUE16      not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the LOAD_VALUE32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOAD_VALUE32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00_00_00: value.SmallInt(4)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       03 01 00 00 00 LOAD_VALUE32      4
`,
		},
		"handle invalid LOAD_VALUE32 index": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOAD_VALUE32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       03 01 00 00 00 LOAD_VALUE32      invalid value index 16777216 (0x1000000)
`,
			err: "invalid value index 16777216 (0x1000000)",
		},
		"handle missing bytes in LOAD_VALUE32": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOAD_VALUE32)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       03             LOAD_VALUE32      not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the ADD opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.ADD)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       04             ADD
`,
		},
		"correctly format the SUBTRACT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.SUBTRACT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       05             SUBTRACT
`,
		},
		"correctly format the MULTIPLY opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.MULTIPLY)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       06             MULTIPLY
`,
		},
		"correctly format the DIVIDE opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DIVIDE)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       07             DIVIDE
`,
		},
		"correctly format the EXPONENTIATE opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.EXPONENTIATE)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       08             EXPONENTIATE
`,
		},
		"correctly format the NEGATE opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEGATE)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       09             NEGATE
`,
		},
		"correctly format the NOT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NOT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       0A             NOT
`,
		},
		"correctly format the BITWISE_NOT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.BITWISE_NOT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       0B             BITWISE_NOT
`,
		},
		"correctly format the TRUE opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.TRUE)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       0C             TRUE
`,
		},
		"correctly format the FALSE opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.FALSE)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       0D             FALSE
`,
		},
		"correctly format the NIL opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NIL)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       0E             NIL
`,
		},
		"correctly format the POP opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.POP)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       0F             POP
`,
		},
		"correctly format the POP_N opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.POP_N), 3},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       10 03          POP_N             3               
`,
		},
		"handle missing bytes in POP_N": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.POP_N)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       10             POP_N             not enough bytes
`,
			err: "not enough bytes",
		},
		"correctly format the LEAVE_SCOPE16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LEAVE_SCOPE16), 3, 2},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       11 03 02       LEAVE_SCOPE16     3               2               
`,
		},
		"correctly format the LEAVE_SCOPE32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LEAVE_SCOPE32), 3, 2, 0, 2},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       12 03 02 00 02 LEAVE_SCOPE32     770             2               
`,
		},
		"correctly format the PREP_LOCALS8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.PREP_LOCALS8), 3},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       13 03          PREP_LOCALS8      3               
`,
		},
		"correctly format the PREP_LOCALS16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.PREP_LOCALS16), 3, 5},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       14 03 05       PREP_LOCALS16     773             
`,
		},
		"correctly format the SET_LOCAL8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.SET_LOCAL8), 3},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       15 03          SET_LOCAL8        3               
`,
		},
		"correctly format the SET_LOCAL16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.SET_LOCAL16), 3, 2},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       16 03 02       SET_LOCAL16       770             
`,
		},
		"correctly format the GET_LOCAL8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GET_LOCAL8), 3},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       17 03          GET_LOCAL8        3               
`,
		},
		"correctly format the GET_LOCAL16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GET_LOCAL16), 3, 2},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       18 03 02       GET_LOCAL16       770             
`,
		},
		"correctly format the JUMP_UNLESS opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.JUMP_UNLESS), 3, 2},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       19 03 02       JUMP_UNLESS       770             
`,
		},
		"correctly format the JUMP opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.JUMP), 3, 2},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       1A 03 02       JUMP              770             
`,
		},
		"correctly format the JUMP_IF opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.JUMP_IF), 3, 2},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       1B 03 02       JUMP_IF           770             
`,
		},
		"correctly format the LOOP opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOOP), 3, 2},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       1C 03 02       LOOP              770             
`,
		},
		"correctly format the JUMP_IF_NIL opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.JUMP_IF_NIL), 3, 2},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       1D 03 02       JUMP_IF_NIL       770             
`,
		},
		"correctly format the RBITSHIFT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.RBITSHIFT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       1E             RBITSHIFT
`,
		},
		"correctly format the LOGIC_RBITSHIFT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOGIC_RBITSHIFT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       1F             LOGIC_RBITSHIFT
`,
		},
		"correctly format the LBITSHIFT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LBITSHIFT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       20             LBITSHIFT
`,
		},
		"correctly format the LOGIC_LBITSHIFT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LOGIC_LBITSHIFT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       21             LOGIC_LBITSHIFT
`,
		},
		"correctly format the BITWISE_AND opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.BITWISE_AND)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       22             BITWISE_AND
`,
		},
		"correctly format the BITWISE_OR opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.BITWISE_OR)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       23             BITWISE_OR
`,
		},
		"correctly format the BITWISE_XOR opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.BITWISE_XOR)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       24             BITWISE_XOR
`,
		},
		"correctly format the MODULO opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.MODULO)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       25             MODULO
`,
		},
		"correctly format the EQUAL opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.EQUAL)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       26             EQUAL
`,
		},
		"correctly format the STRICT_EQUAL opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.STRICT_EQUAL)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       27             STRICT_EQUAL
`,
		},
		"correctly format the GREATER opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GREATER)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       28             GREATER
`,
		},
		"correctly format the GREATER_EQUAL opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GREATER_EQUAL)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       29             GREATER_EQUAL
`,
		},
		"correctly format the LESS opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LESS)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       2A             LESS
`,
		},
		"correctly format the LESS_EQUAL opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LESS_EQUAL)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       2B             LESS_EQUAL
`,
		},
		"correctly format the GET_MOD_CONST8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GET_MOD_CONST8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0: value.SymbolTable.Add("Foo")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       2C 00          GET_MOD_CONST8    :Foo
`,
		},
		"correctly format the GET_MOD_CONST16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GET_MOD_CONST16), 0x01, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00: value.SymbolTable.Add("Bar")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       2D 01 00       GET_MOD_CONST16   :Bar
`,
		},
		"correctly format the GET_MOD_CONST32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GET_MOD_CONST32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00_00_00: value.SymbolTable.Add("Bar")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       2E 01 00 00 00 GET_MOD_CONST32   :Bar
`,
		},
		"correctly format the ROOT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.ROOT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       2F             ROOT
`,
		},
		"correctly format the NOT_EQUAL opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NOT_EQUAL)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       30             NOT_EQUAL
`,
		},
		"correctly format the STRICT_NOT_EQUAL opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.STRICT_NOT_EQUAL)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       31             STRICT_NOT_EQUAL
`,
		},
		"correctly format the DEF_MOD_CONST8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_MOD_CONST8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0: value.SymbolTable.Add("Foo")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       32 00          DEF_MOD_CONST8    :Foo
`,
		},
		"correctly format the DEF_MOD_CONST16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_MOD_CONST16), 0x01, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00: value.SymbolTable.Add("Bar")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       33 01 00       DEF_MOD_CONST16   :Bar
`,
		},
		"correctly format the DEF_MOD_CONST32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_MOD_CONST32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00_00_00: value.SymbolTable.Add("Bar")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       34 01 00 00 00 DEF_MOD_CONST32   :Bar
`,
		},
		"correctly format the CONSTANT_CONTAINER opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.CONSTANT_CONTAINER)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       35             CONSTANT_CONTAINER
`,
		},
		"correctly format the DEF_CLASS16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_CLASS)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       36             DEF_CLASS
`,
		},
		"correctly format the SELF opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.SELF)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       37             SELF
`,
		},
		"correctly format the DEF_MODULE opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_MODULE)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       38             DEF_MODULE
`,
		},
		"correctly format the CALL_METHOD8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.CALL_METHOD8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0: value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       39 00          CALL_METHOD8      CallSiteInfo{name: :foo, argument_count: 0}
`,
		},
		"correctly format the CALL_METHOD16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.CALL_METHOD16), 0x01, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00: value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       3A 01 00       CALL_METHOD16     CallSiteInfo{name: :foo, argument_count: 0}
`,
		},
		"correctly format the CALL_METHOD32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.CALL_METHOD32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00_00_00: value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       3B 01 00 00 00 CALL_METHOD32     CallSiteInfo{name: :foo, argument_count: 0}
`,
		},
		"correctly format the DEF_METHOD opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_METHOD)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       3C             DEF_METHOD
`,
		},
		"correctly format the UNDEFINED opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.UNDEFINED)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       3D             UNDEFINED
`,
		},
		"correctly format the DEF_ANON_CLASS opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_ANON_CLASS)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       3E             DEF_ANON_CLASS
`,
		},
		"correctly format the DEF_ANON_MODULE opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_ANON_MODULE)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       3F             DEF_ANON_MODULE
`,
		},
		"correctly format the CALL_FUNCTION8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.CALL_FUNCTION8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0: value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       40 00          CALL_FUNCTION8    CallSiteInfo{name: :foo, argument_count: 0}
`,
		},
		"correctly format the CALL_FUNCTION16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.CALL_FUNCTION16), 0x01, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00: value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       41 01 00       CALL_FUNCTION16   CallSiteInfo{name: :foo, argument_count: 0}
`,
		},
		"correctly format the CALL_FUNCTION32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.CALL_FUNCTION32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00_00_00: value.NewCallSiteInfo(value.SymbolTable.Add("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       42 01 00 00 00 CALL_FUNCTION32   CallSiteInfo{name: :foo, argument_count: 0}
`,
		},
		"correctly format the DEF_MIXIN opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_MIXIN)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       43             DEF_MIXIN
`,
		},
		"correctly format the DEF_ANON_MIXIN opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_ANON_MIXIN)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       44             DEF_ANON_MIXIN
`,
		},
		"correctly format the INCLUDE opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.INCLUDE)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       45             INCLUDE
`,
		},
		"correctly format the GET_SINGLETON opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GET_SINGLETON)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       46             GET_SINGLETON
`,
		},
		"correctly format the JUMP_IF_UNDEF opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.JUMP_UNLESS_UNDEF), 3, 2},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       47 03 02       JUMP_UNLESS_UNDEF 770             
`,
		},
		"correctly format the DEF_ALIAS opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_ALIAS)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       48             DEF_ALIAS
`,
		},
		"correctly format the METHOD_CONTAINER opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.METHOD_CONTAINER)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       49             METHOD_CONTAINER
`,
		},
		"correctly format the COMPARE opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.COMPARE)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       4A             COMPARE
`,
		},
		"correctly format the DOC_COMMENT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DOC_COMMENT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       4B             DOC_COMMENT
`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.in.DisassembleString()
			if diff := cmp.Diff(tc.want, got, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			if diff := cmp.Diff(tc.err, gotErr, comparer.Comparer); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
