// Package bytecode implements types
// and constants that make up Elk
// bytecode.
package value

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position"
)

// A single unit of Elk bytecode.
type BytecodeFunction struct {
	Instructions           []byte
	Values                 []Value // The value pool
	Parameters             []Symbol
	OptionalParameterCount int
	LineInfoList           bytecode.LineInfoList
	Location               *position.Location
	Name                   Symbol
	NamedRestArgument      bool
}

func (*BytecodeFunction) method() {}

// Get the number of parameters
func (b *BytecodeFunction) ParameterCount() int {
	return len(b.Parameters)
}

func (*BytecodeFunction) Class() *Class {
	return nil
}

func (*BytecodeFunction) DirectClass() *Class {
	return nil
}

func (*BytecodeFunction) SingletonClass() *Class {
	return nil
}

func (*BytecodeFunction) IsFrozen() bool {
	return true
}

func (*BytecodeFunction) SetFrozen() {}

func (b *BytecodeFunction) Inspect() string {
	return fmt.Sprintf("BytecodeFunction{name: %s, location: %s}", b.Name.Name(), b.Location.String())
}

func (*BytecodeFunction) InstanceVariables() SimpleSymbolMap {
	return nil
}

// Create a new chunk of bytecode.
func NewBytecodeFunction(name Symbol, instruct []byte, loc *position.Location) *BytecodeFunction {
	return &BytecodeFunction{
		Instructions: instruct,
		Location:     loc,
		Name:         name,
	}
}

// Add a parameter to the function.
func (f *BytecodeFunction) AddParameter(name Symbol) {
	f.Parameters = append(f.Parameters, name)
}

// Add an instruction to the bytecode chunk.
func (f *BytecodeFunction) AddInstruction(lineNumber int, op bytecode.OpCode, bytes ...byte) {
	f.LineInfoList.AddLineNumber(lineNumber)
	f.Instructions = append(f.Instructions, byte(op))
	f.Instructions = append(f.Instructions, bytes...)
}

// Add bytes to the bytecode chunk.
func (f *BytecodeFunction) AddBytes(bytes ...byte) {
	f.Instructions = append(f.Instructions, bytes...)
}

// Append two bytes to the bytecode chunk.
func (f *BytecodeFunction) AppendUint16(n uint16) {
	f.Instructions = binary.BigEndian.AppendUint16(f.Instructions, n)
}

// Append four bytes to the bytecode chunk.
func (f *BytecodeFunction) AppendUint32(n uint32) {
	f.Instructions = binary.BigEndian.AppendUint32(f.Instructions, n)
}

// Size of an integer.
type IntSize uint8

// Add a constant to the constant pool.
// Returns the index of the constant.
func (f *BytecodeFunction) AddValue(obj Value) (int, IntSize) {
	var id int
	switch obj.(type) {
	case String, SmallInt, Int64, Int32, Int16,
		Int8, UInt64, UInt32, UInt16, UInt8,
		Float, Float32, Float64:
		if i := slices.Index(f.Values, obj); i != -1 {
			id = i
			break
		}
		id = len(f.Values)
		f.Values = append(f.Values, obj)
	default:
		id = len(f.Values)
		f.Values = append(f.Values, obj)
	}

	if id <= math.MaxUint8 {
		return id, bytecode.UINT8_SIZE
	}

	if id <= math.MaxUint16 {
		return id, bytecode.UINT16_SIZE
	}

	if id <= math.MaxUint32 {
		return id, bytecode.UINT32_SIZE
	}

	return id, bytecode.UINT64_SIZE
}

// Disassemble the bytecode chunk and write the
// output to stdout.
func (f *BytecodeFunction) DisassembleStdout() {
	f.Disassemble(os.Stdout)
}

// Disassemble the bytecode chunk and return a string
// containing the result.
func (f *BytecodeFunction) DisassembleString() (string, error) {
	var buffer strings.Builder
	err := f.Disassemble(&buffer)
	if err != nil {
		return buffer.String(), err
	}

	return buffer.String(), nil
}

// Disassemble the bytecode chunk and write the
// output to a writer.
func (f *BytecodeFunction) Disassemble(output io.Writer) error {
	fmt.Fprintf(output, "== Disassembly of %s at: %s ==\n\n", f.Name.Name(), f.Location.String())

	if len(f.Instructions) == 0 {
		return nil
	}

	var offset int
	var instructionIndex int
	for {
		result, err := f.DisassembleInstruction(output, offset, instructionIndex)
		if err != nil {
			return err
		}
		offset = result
		instructionIndex++
		if offset >= len(f.Instructions) {
			break
		}
	}

	for _, constant := range f.Values {
		fn, ok := constant.(*BytecodeFunction)
		if !ok {
			continue
		}
		fmt.Fprintln(output)
		fn.Disassemble(output)
	}

	return nil
}

func (f *BytecodeFunction) DisassembleInstruction(output io.Writer, offset, instructionIndex int) (int, error) {
	fmt.Fprintf(output, "%04d  ", offset)
	opcodeByte := f.Instructions[offset]
	opcode := bytecode.OpCode(opcodeByte)
	switch opcode {
	case bytecode.RETURN, bytecode.ADD, bytecode.SUBTRACT,
		bytecode.MULTIPLY, bytecode.DIVIDE, bytecode.EXPONENTIATE,
		bytecode.NEGATE, bytecode.NOT, bytecode.BITWISE_NOT,
		bytecode.TRUE, bytecode.FALSE, bytecode.NIL, bytecode.POP,
		bytecode.RBITSHIFT, bytecode.LBITSHIFT,
		bytecode.LOGIC_RBITSHIFT, bytecode.LOGIC_LBITSHIFT,
		bytecode.BITWISE_AND, bytecode.BITWISE_OR, bytecode.BITWISE_XOR, bytecode.MODULO,
		bytecode.EQUAL, bytecode.STRICT_EQUAL, bytecode.GREATER, bytecode.GREATER_EQUAL, bytecode.LESS, bytecode.LESS_EQUAL,
		bytecode.ROOT, bytecode.NOT_EQUAL, bytecode.STRICT_NOT_EQUAL,
		bytecode.CONSTANT_CONTAINER, bytecode.DEF_CLASS, bytecode.SELF, bytecode.DEF_MODULE, bytecode.DEF_METHOD,
		bytecode.UNDEFINED, bytecode.DEF_ANON_CLASS, bytecode.DEF_ANON_MODULE,
		bytecode.DEF_MIXIN, bytecode.DEF_ANON_MIXIN, bytecode.INCLUDE, bytecode.GET_SINGLETON:
		return f.disassembleOneByteInstruction(output, opcode.String(), offset, instructionIndex), nil
	case bytecode.POP_N, bytecode.SET_LOCAL8, bytecode.GET_LOCAL8, bytecode.PREP_LOCALS8:
		return f.disassembleNumericOperands(output, 1, 1, offset, instructionIndex)
	case bytecode.PREP_LOCALS16, bytecode.SET_LOCAL16, bytecode.GET_LOCAL16, bytecode.JUMP_UNLESS, bytecode.JUMP,
		bytecode.JUMP_IF, bytecode.LOOP, bytecode.JUMP_IF_NIL, bytecode.JUMP_UNLESS_UNDEF:
		return f.disassembleNumericOperands(output, 1, 2, offset, instructionIndex)
	case bytecode.LEAVE_SCOPE16:
		return f.disassembleNumericOperands(output, 2, 1, offset, instructionIndex)
	case bytecode.LEAVE_SCOPE32:
		return f.disassembleNumericOperands(output, 2, 2, offset, instructionIndex)
	case bytecode.LOAD_VALUE8, bytecode.GET_MOD_CONST8,
		bytecode.DEF_MOD_CONST8, bytecode.CALL_METHOD8,
		bytecode.CALL_FUNCTION8:
		return f.disassembleConstant(output, 2, offset, instructionIndex)
	case bytecode.LOAD_VALUE16, bytecode.GET_MOD_CONST16,
		bytecode.DEF_MOD_CONST16, bytecode.CALL_METHOD16,
		bytecode.CALL_FUNCTION16:
		return f.disassembleConstant(output, 3, offset, instructionIndex)
	case bytecode.LOAD_VALUE32, bytecode.GET_MOD_CONST32,
		bytecode.DEF_MOD_CONST32, bytecode.CALL_METHOD32,
		bytecode.CALL_FUNCTION32:
		return f.disassembleConstant(output, 5, offset, instructionIndex)
	default:
		f.printLineNumber(output, instructionIndex)
		f.dumpBytes(output, offset, 1)
		fmt.Fprintf(output, "unknown operation %d (0x%X)\n", opcodeByte, opcodeByte)
		return offset + 1, fmt.Errorf("unknown operation %d (0x%X) at offset %d (0x%X)", opcodeByte, opcodeByte, offset, offset)
	}
}

func (f *BytecodeFunction) dumpBytes(output io.Writer, offset, count int) {
	for i := offset; i < offset+count; i++ {
		fmt.Fprintf(output, "%02X ", f.Instructions[i])
	}

	for i := count; i < bytecode.MaxInstructionByteLength; i++ {
		fmt.Fprint(output, "   ")
	}
}

func (f *BytecodeFunction) disassembleOneByteInstruction(output io.Writer, name string, offset, instructionIndex int) int {
	f.printLineNumber(output, instructionIndex)
	f.dumpBytes(output, offset, 1)
	fmt.Fprintln(output, name)
	return offset + 1
}

func (f *BytecodeFunction) disassembleNumericOperands(output io.Writer, operands, operandBytes, offset, instructionIndex int) (int, error) {
	bytes := 1 + operands*operandBytes
	if result, err := f.checkBytes(output, offset, instructionIndex, bytes); err != nil {
		return result, err
	}

	opcode := bytecode.OpCode(f.Instructions[offset])

	f.printLineNumber(output, instructionIndex)
	f.dumpBytes(output, offset, bytes)
	f.printOpCode(output, opcode)

	var readFunc intReadFunc
	switch operandBytes {
	case 8:
		readFunc = readUint64
	case 4:
		readFunc = readUint32
	case 2:
		readFunc = readUint16
	case 1:
		readFunc = readUint8
	default:
		panic(fmt.Sprintf("incorrect bytesize of operands: %d", operandBytes))
	}

	for i := 0; i < operands; i++ {
		a := readFunc(f.Instructions[offset+1+i*operandBytes : offset+1+(i+1)*operandBytes])
		f.printNumField(output, a)
	}
	fmt.Fprintln(output)

	return offset + bytes, nil
}

type intReadFunc func([]byte) uint64

func readUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func readUint32(b []byte) uint64 {
	return uint64(binary.BigEndian.Uint32(b))
}

func readUint16(b []byte) uint64 {
	return uint64(binary.BigEndian.Uint16(b))
}

func readUint8(b []byte) uint64 {
	return uint64(b[0])
}

func (f *BytecodeFunction) disassembleConstant(output io.Writer, byteLength, offset, instructionIndex int) (int, error) {
	opcode := bytecode.OpCode(f.Instructions[offset])

	if result, err := f.checkBytes(output, offset, instructionIndex, byteLength); err != nil {
		return result, err
	}

	var constantIndex int
	if byteLength == 2 {
		constantIndex = int(f.Instructions[offset+1])
	} else if byteLength == 3 {
		constantIndex = int(binary.BigEndian.Uint16(f.Instructions[offset+1 : offset+3]))
	} else if byteLength == 5 {
		constantIndex = int(binary.BigEndian.Uint32(f.Instructions[offset+1 : offset+5]))
	} else {
		panic(fmt.Sprintf("%d is not a valid byteLength for a value opcode!", byteLength))
	}

	f.printLineNumber(output, instructionIndex)
	f.dumpBytes(output, offset, byteLength)
	f.printOpCode(output, opcode)

	if constantIndex >= len(f.Values) {
		msg := fmt.Sprintf("invalid value index %d (0x%X)", constantIndex, constantIndex)
		fmt.Fprintln(output, msg)
		return offset + byteLength, fmt.Errorf(msg)
	}
	constant := f.Values[constantIndex]
	fmt.Fprintln(output, constant.Inspect())

	return offset + byteLength, nil
}

func (f *BytecodeFunction) checkBytes(output io.Writer, offset, instructionIndex, byteLength int) (int, error) {
	opcode := bytecode.OpCode(f.Instructions[offset])
	if len(f.Instructions)-offset >= byteLength {
		return 0, nil
	}
	f.printLineNumber(output, instructionIndex)
	f.dumpBytes(output, offset, len(f.Instructions)-offset)
	f.printOpCode(output, opcode)
	msg := "not enough bytes"
	fmt.Fprintln(output, msg)
	return len(f.Instructions) - 1, fmt.Errorf(msg)
}

func (f *BytecodeFunction) printLineNumber(output io.Writer, instructionIndex int) {
	fmt.Fprintf(output, "%-8s", f.getLineNumberString(instructionIndex))
}

func (f *BytecodeFunction) getLineNumberString(instructionIndex int) string {
	currentLineNumber := f.LineInfoList.GetLineNumber(instructionIndex)
	if instructionIndex == 0 {
		return fmt.Sprintf("%d", currentLineNumber)
	}

	previousLineNumber := f.LineInfoList.GetLineNumber(instructionIndex - 1)
	if previousLineNumber == currentLineNumber {
		return "|"
	}

	return fmt.Sprintf("%d", currentLineNumber)
}

func (f *BytecodeFunction) printOpCode(output io.Writer, opcode bytecode.OpCode) {
	fmt.Fprintf(output, "%-18s", opcode.String())
}

func (f *BytecodeFunction) printNumField(output io.Writer, n uint64) {
	fmt.Fprintf(output, "%-16d", n)
}
