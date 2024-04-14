package vm_test

import (
	"testing"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

var mainSymbol = value.ToSymbol("main")

func TestBytecodeMethod_AddInstruction(t *testing.T) {
	c := &vm.BytecodeMethod{}
	c.AddInstruction(1, bytecode.RETURN)
	want := &vm.BytecodeMethod{
		Instructions: []byte{byte(bytecode.RETURN)},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	if diff := cmp.Diff(want, c, comparer.Options()...); diff != "" {
		t.Fatalf(diff)
	}

	c = &vm.BytecodeMethod{}
	c.AddInstruction(1, bytecode.LOAD_VALUE8, 0x12)
	want = &vm.BytecodeMethod{
		Instructions: []byte{byte(bytecode.LOAD_VALUE8), 0x12},
		LineInfoList: bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
	}
	if diff := cmp.Diff(want, c, comparer.Options()); diff != "" {
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
	if diff := cmp.Diff(want, c, comparer.Options()); diff != "" {
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
	if diff := cmp.Diff(want, c, comparer.Options()); diff != "" {
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
			if diff := cmp.Diff(tc.chunkAfter, tc.chunkBefore, comparer.Options()); diff != "" {
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

0000  1       FF                unknown operation 255 (0xFF)
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

0000  1       01                RETURN
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

0000  1       02 00             LOAD_VALUE8       4
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

0000  1       02 19             LOAD_VALUE8       invalid value index 25 (0x19)
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

0000  1       02                LOAD_VALUE8       not enough bytes
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

0000  1       03 01 00          LOAD_VALUE16      4
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

0000  1       03 19 FF          LOAD_VALUE16      invalid value index 6655 (0x19FF)
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

0000  1       03                LOAD_VALUE16      not enough bytes
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

0000  1       04 01 00 00 00    LOAD_VALUE32      4
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

0000  1       04 01 00 00 00    LOAD_VALUE32      invalid value index 16777216 (0x1000000)
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

0000  1       04                LOAD_VALUE32      not enough bytes
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

0000  1       05                ADD
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

0000  1       06                SUBTRACT
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

0000  1       07                MULTIPLY
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

0000  1       08                DIVIDE
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

0000  1       09                EXPONENTIATE
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

0000  1       0A                NEGATE
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

0000  1       0B                NOT
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

0000  1       0C                BITWISE_NOT
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

0000  1       0D                TRUE
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

0000  1       0E                FALSE
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

0000  1       0F                NIL
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

0000  1       10                POP
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

0000  1       11 03             POP_N             3               
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

0000  1       11                POP_N             not enough bytes
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

0000  1       12 03 02          LEAVE_SCOPE16     3               2               
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

0000  1       13 03 02 00 02    LEAVE_SCOPE32     770             2               
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

0000  1       14 03             PREP_LOCALS8      3               
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

0000  1       15 03 05          PREP_LOCALS16     773             
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

0000  1       16 03             SET_LOCAL8        3               
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

0000  1       17 03 02          SET_LOCAL16       770             
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

0000  1       18 03             GET_LOCAL8        3               
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

0000  1       19 03 02          GET_LOCAL16       770             
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

0000  1       1A 03 02          JUMP_UNLESS       770             
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

0000  1       1B 03 02          JUMP              770             
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

0000  1       1C 03 02          JUMP_IF           770             
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

0000  1       1D 03 02          LOOP              770             
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

0000  1       1E 03 02          JUMP_IF_NIL       770             
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

0000  1       1F                RBITSHIFT
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

0000  1       20                LOGIC_RBITSHIFT
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

0000  1       21                LBITSHIFT
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

0000  1       22                LOGIC_LBITSHIFT
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

0000  1       23                BITWISE_AND
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

0000  1       24                BITWISE_OR
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

0000  1       25                BITWISE_XOR
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

0000  1       26                MODULO
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

0000  1       27                EQUAL
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

0000  1       28                STRICT_EQUAL
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

0000  1       29                GREATER
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

0000  1       2A                GREATER_EQUAL
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

0000  1       2B                LESS
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

0000  1       2C                LESS_EQUAL
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
				[]value.Value{0: value.ToSymbol("Foo")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       2D 00             GET_MOD_CONST8    :Foo
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
				[]value.Value{0x1_00: value.ToSymbol("Bar")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       2E 01 00          GET_MOD_CONST16   :Bar
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
				[]value.Value{0x1_00_00_00: value.ToSymbol("Bar")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       2F 01 00 00 00    GET_MOD_CONST32   :Bar
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

0000  1       30                ROOT
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

0000  1       31                NOT_EQUAL
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

0000  1       32                STRICT_NOT_EQUAL
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
				[]value.Value{0: value.ToSymbol("Foo")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       33 00             DEF_MOD_CONST8    :Foo
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
				[]value.Value{0x1_00: value.ToSymbol("Bar")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       34 01 00          DEF_MOD_CONST16   :Bar
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
				[]value.Value{0x1_00_00_00: value.ToSymbol("Bar")},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       35 01 00 00 00    DEF_MOD_CONST32   :Bar
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

0000  1       36                CONSTANT_CONTAINER
`,
		},
		"correctly format the DEF_CLASS opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_CLASS), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       37 00             DEF_CLASS         0               
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

0000  1       38                SELF
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

0000  1       39                DEF_MODULE
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
				[]value.Value{0: value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       3A 00             CALL_METHOD8      CallSiteInfo{name: :foo, argument_count: 0}
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
				[]value.Value{0x1_00: value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       3B 01 00          CALL_METHOD16     CallSiteInfo{name: :foo, argument_count: 0}
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
				[]value.Value{0x1_00_00_00: value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       3C 01 00 00 00    CALL_METHOD32     CallSiteInfo{name: :foo, argument_count: 0}
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

0000  1       3D                DEF_METHOD
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

0000  1       3E                UNDEFINED
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

0000  1       3F                DEF_ANON_CLASS
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

0000  1       40                DEF_ANON_MODULE
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
				[]value.Value{0: value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       41 00             CALL_FUNCTION8    CallSiteInfo{name: :foo, argument_count: 0}
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
				[]value.Value{0x1_00: value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       42 01 00          CALL_FUNCTION16   CallSiteInfo{name: :foo, argument_count: 0}
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
				[]value.Value{0x1_00_00_00: value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       43 01 00 00 00    CALL_FUNCTION32   CallSiteInfo{name: :foo, argument_count: 0}
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

0000  1       44                DEF_MIXIN
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

0000  1       45                DEF_ANON_MIXIN
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

0000  1       46                INCLUDE
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

0000  1       47                GET_SINGLETON
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

0000  1       48 03 02          JUMP_UNLESS_UNDEF 770             
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

0000  1       49                DEF_ALIAS
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

0000  1       4A                METHOD_CONTAINER
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

0000  1       4B                COMPARE
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

0000  1       4C                DOC_COMMENT
`,
		},
		"correctly format the DEF_GETTER opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_GETTER)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       4D                DEF_GETTER
`,
		},
		"correctly format the DEF_SETTER opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_SETTER)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       4E                DEF_SETTER
`,
		},
		"correctly format the DEF_SINGLETON opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DEF_SINGLETON)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       4F                DEF_SINGLETON
`,
		},
		"correctly format the RETURN_FIRST_ARG opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.RETURN_FIRST_ARG)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       50                RETURN_FIRST_ARG
`,
		},
		"correctly format the INSTANTIATE8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.INSTANTIATE8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0: value.NewCallSiteInfo(value.ToSymbol("#init"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       51 00             INSTANTIATE8      CallSiteInfo{name: :"#init", argument_count: 0}
`,
		},
		"correctly format the INSTANTIATE16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.INSTANTIATE16), 0x01, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00: value.NewCallSiteInfo(value.ToSymbol("#init"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       52 01 00          INSTANTIATE16     CallSiteInfo{name: :"#init", argument_count: 0}
`,
		},
		"correctly format the INSTANTIATE32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.INSTANTIATE32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00_00_00: value.NewCallSiteInfo(value.ToSymbol("#init"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       53 01 00 00 00    INSTANTIATE32     CallSiteInfo{name: :"#init", argument_count: 0}
`,
		},
		"correctly format the RETURN_SELF opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.RETURN_SELF)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       54                RETURN_SELF
`,
		},
		"correctly format the GET_IVAR8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GET_IVAR8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{value.SmallInt(4)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       55 00             GET_IVAR8         4
`,
		},
		"correctly format the GET_IVAR16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GET_IVAR16), 0x01, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00: value.SmallInt(4)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       56 01 00          GET_IVAR16        4
`,
		},
		"correctly format the GET_IVAR32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GET_IVAR32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00_00_00: value.SmallInt(4)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       57 01 00 00 00    GET_IVAR32        4
`,
		},
		"correctly format the SET_IVAR8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.SET_IVAR8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{value.SmallInt(4)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       58 00             SET_IVAR8         4
`,
		},
		"correctly format the SET_IVAR16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.SET_IVAR16), 0x01, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00: value.SmallInt(4)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       59 01 00          SET_IVAR16        4
`,
		},
		"correctly format the SET_IVAR32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.SET_IVAR32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00_00_00: value.SmallInt(4)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       5A 01 00 00 00    SET_IVAR32        4
`,
		},
		"correctly format the NEW_ARRAY_TUPLE8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_ARRAY_TUPLE8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       5B 00             NEW_ARRAY_TUPLE8  0               
`,
		},
		"correctly format the NEW_ARRAY_TUPLE32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_ARRAY_TUPLE32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       5C 01 00 00 00    NEW_ARRAY_TUPLE32 16777216        
`,
		},
		"correctly format the APPEND opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.APPEND)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       5D                APPEND
`,
		},
		"correctly format the COPY opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.COPY)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       5E                COPY
`,
		},
		"correctly format the SUBSCRIPT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.SUBSCRIPT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       5F                SUBSCRIPT
`,
		},
		"correctly format the SUBSCRIPT_SET opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.SUBSCRIPT_SET)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       60                SUBSCRIPT_SET
`,
		},
		"correctly format the APPEND_AT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.APPEND_AT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       61                APPEND_AT
`,
		},
		"correctly format the NEW_ARRAY_LIST8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_ARRAY_LIST8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       62 00             NEW_ARRAY_LIST8   0               
`,
		},
		"correctly format the NEW_ARRAY_LIST32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_ARRAY_LIST32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       63 01 00 00 00    NEW_ARRAY_LIST32  16777216        
`,
		},
		"correctly format the GET_ITERATOR opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.GET_ITERATOR)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       64                GET_ITERATOR
`,
		},
		"correctly format the FOR_IN opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.FOR_IN), 3, 2},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       65 03 02          FOR_IN            770             
`,
		},
		"correctly format the NEW_STRING8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_STRING8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       66 00             NEW_STRING8       0               
`,
		},
		"correctly format the NEW_STRING32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_STRING32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       67 01 00 00 00    NEW_STRING32      16777216        
`,
		},
		"correctly format the NEW_HASH_MAP8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_HASH_MAP8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       68 00             NEW_HASH_MAP8     0               
`,
		},
		"correctly format the NEW_HASH_MAP32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_HASH_MAP32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       69 01 00 00 00    NEW_HASH_MAP32    16777216        
`,
		},
		"correctly format the MAP_SET opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.MAP_SET)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       6A                MAP_SET
`,
		},
		"correctly format the NEW_HASH_RECORD8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_HASH_RECORD8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       6B 00             NEW_HASH_RECORD8  0               
`,
		},
		"correctly format the NEW_HASH_RECORD32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_HASH_RECORD32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       6C 01 00 00 00    NEW_HASH_RECORD32 16777216        
`,
		},
		"correctly format the LAX_EQUAL opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LAX_EQUAL)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       6D                LAX_EQUAL
`,
		},
		"correctly format the LAX_NOT_EQUAL opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.LAX_NOT_EQUAL)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       6E                LAX_NOT_EQUAL
`,
		},
		"correctly format the NEW_REGEX8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_REGEX8), 3, 4},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       6F 03 04          NEW_REGEX8        im-sUxa         4               
`,
		},
		"correctly format the NEW_REGEX32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_REGEX32), 5, 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       70 05 01 00 00 00 NEW_REGEX32       is-mUxa         16777216        
`,
		},
		"correctly format the BITWISE_AND_NOT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.BITWISE_AND_NOT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       71                BITWISE_AND_NOT
`,
		},
		"correctly format the UNARY_PLUS opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.UNARY_PLUS)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       72                UNARY_PLUS
`,
		},
		"correctly format the INCREMENT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.INCREMENT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       73                INCREMENT
`,
		},
		"correctly format the DECREMENT opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DECREMENT)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       74                DECREMENT
`,
		},
		"correctly format the DUP opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DUP)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       75                DUP
`,
		},
		"correctly format the DUP_N opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.DUP_N), 8},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       76 08             DUP_N             8               
`,
		},
		"correctly format the POP_N_SKIP_ONE opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.POP_N_SKIP_ONE), 8},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       77 08             POP_N_SKIP_ONE    8               
`,
		},
		"correctly format the NEW_SYMBOL8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_SYMBOL8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       78 00             NEW_SYMBOL8       0               
`,
		},
		"correctly format the NEW_SYMBOL32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_SYMBOL32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       79 01 00 00 00    NEW_SYMBOL32      16777216        
`,
		},
		"correctly format the SWAP opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.SWAP)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7A                SWAP
`,
		},
		"correctly format the NEW_RANGE opcode with closed arg": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_RANGE), bytecode.CLOSED_RANGE_FLAG},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7B 00             NEW_RANGE         x...y           
`,
		},
		"correctly format the NEW_RANGE opcode with open arg": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_RANGE), bytecode.OPEN_RANGE_FLAG},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7B 01             NEW_RANGE         x<.<y           
`,
		},
		"correctly format the NEW_RANGE opcode with left open arg": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_RANGE), bytecode.LEFT_OPEN_RANGE_FLAG},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7B 02             NEW_RANGE         x<..y           
`,
		},
		"correctly format the NEW_RANGE opcode with right open arg": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_RANGE), bytecode.RIGHT_OPEN_RANGE_FLAG},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7B 03             NEW_RANGE         x..<y           
`,
		},
		"correctly format the NEW_RANGE opcode with beginless closed arg": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_RANGE), bytecode.BEGINLESS_CLOSED_RANGE_FLAG},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7B 04             NEW_RANGE         ...x            
`,
		},
		"correctly format the NEW_RANGE opcode with beginless open arg": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_RANGE), bytecode.BEGINLESS_OPEN_RANGE_FLAG},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7B 05             NEW_RANGE         ..<x            
`,
		},
		"correctly format the NEW_RANGE opcode with endless closed arg": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_RANGE), bytecode.ENDLESS_CLOSED_RANGE_FLAG},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7B 06             NEW_RANGE         x...            
`,
		},
		"correctly format the NEW_RANGE opcode with endless open arg": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_RANGE), bytecode.ENDLESS_OPEN_RANGE_FLAG},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7B 07             NEW_RANGE         x<..            
`,
		},
		"correctly format the CALL_PATTERN8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.CALL_PATTERN8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0: value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7C 00             CALL_PATTERN8     CallSiteInfo{name: :foo, argument_count: 0}
`,
		},
		"correctly format the CALL_PATTERN16 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.CALL_PATTERN16), 0x01, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00: value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7D 01 00          CALL_PATTERN16    CallSiteInfo{name: :foo, argument_count: 0}
`,
		},
		"correctly format the CALL_PATTERN32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.CALL_PATTERN32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				[]value.Value{0x1_00_00_00: value.NewCallSiteInfo(value.ToSymbol("foo"), 0, nil)},
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7E 01 00 00 00    CALL_PATTERN32    CallSiteInfo{name: :foo, argument_count: 0}
`,
		},
		"correctly format the INSTANCE_OF opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.INSTANCE_OF)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       7F                INSTANCE_OF
`,
		},
		"correctly format the IS_A opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.IS_A)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       80                IS_A
`,
		},
		"correctly format the POP_SKIP_ONE opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.POP_SKIP_ONE)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       81                POP_SKIP_ONE
`,
		},
		"correctly format the INSPECT_STACK opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.INSPECT_STACK)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       82                INSPECT_STACK
`,
		},
		"correctly format the NEW_HASH_SET8 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_HASH_SET8), 0},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       83 00             NEW_HASH_SET8     0               
`,
		},
		"correctly format the NEW_HASH_SET32 opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.NEW_HASH_SET32), 0x01, 0x00, 0x00, 0x00},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       84 01 00 00 00    NEW_HASH_SET32    16777216        
`,
		},
		"correctly format the THROW opcode": {
			in: vm.NewBytecodeMethod(
				mainSymbol,
				[]byte{byte(bytecode.THROW)},
				L(P(12, 2, 3), P(18, 2, 9)),
				bytecode.LineInfoList{bytecode.NewLineInfo(1, 1)},
				nil,
				0,
				-1,
				false, false,
				nil,
			),
			want: `== Disassembly of main at: sourceName:2:3 ==

0000  1       85                THROW
`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.in.DisassembleString()
			if diff := cmp.Diff(tc.want, got, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			if diff := cmp.Diff(tc.err, gotErr, comparer.Options()); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
