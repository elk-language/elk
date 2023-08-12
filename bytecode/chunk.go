// Package bytecode implements types
// and constants that make up Elk
// bytecode.
package bytecode

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"slices"

	"github.com/elk-language/elk/object"
	"github.com/elk-language/elk/position"
)

// A single unit of Elk bytecode.
type Chunk struct {
	Instructions []byte
	Constants    []object.Value // The constant pool
	LineInfoList LineInfoList
	Location     *position.Location
}

// Create a new chunk of bytecode.
func NewChunk(instruct []byte, loc *position.Location) *Chunk {
	return &Chunk{
		Instructions: instruct,
		Location:     loc,
	}
}

// Add an instruction to the bytecode chunk.
func (c *Chunk) AddInstruction(lineNumber int, op OpCode, bytes ...byte) {
	c.LineInfoList.AddLineNumber(lineNumber)
	c.Instructions = append(c.Instructions, byte(op))
	c.Instructions = append(c.Instructions, bytes...)
}

// Add bytes to the bytecode chunk.
func (c *Chunk) AddBytes(bytes ...byte) {
	c.Instructions = append(c.Instructions, bytes...)
}

// Append two bytes to the bytecode chunk.
func (c *Chunk) AppendUint16(n uint16) {
	c.Instructions = binary.BigEndian.AppendUint16(c.Instructions, n)
}

// Append four bytes to the bytecode chunk.
func (c *Chunk) AppendUint32(n uint32) {
	c.Instructions = binary.BigEndian.AppendUint32(c.Instructions, n)
}

// Size of an integer.
type IntSize uint8

const (
	UINT8_SIZE  = iota // The integer fits in a uint8
	UINT16_SIZE        // The integer fits in a uint16
	UINT32_SIZE        // The integer fits in a uint32
	UINT64_SIZE        // The integer fits in a uint64
)

// Add a constant to the constant pool.
// Returns the index of the constant.
func (c *Chunk) AddConstant(obj object.Value) (int, IntSize) {
	var id int
	switch obj.(type) {
	case object.String, object.SmallInt, object.Int64, object.Int32, object.Int16,
		object.Int8, object.UInt64, object.UInt32, object.UInt16, object.UInt8,
		object.Float, object.Float32, object.Float64:
		if i := slices.Index(c.Constants, obj); i != -1 {
			id = i
			break
		}
		id = len(c.Constants)
		c.Constants = append(c.Constants, obj)
	default:
		id = len(c.Constants)
		c.Constants = append(c.Constants, obj)
	}

	if id <= math.MaxUint8 {
		return id, UINT8_SIZE
	}

	if id <= math.MaxUint16 {
		return id, UINT16_SIZE
	}

	if id <= math.MaxUint32 {
		return id, UINT32_SIZE
	}

	return id, UINT64_SIZE
}

// Disassemble the bytecode chunk and write the
// output to stdout.
func (c *Chunk) DisassembleStdout() {
	c.Disassemble(os.Stdout)
}

// Disassemble the bytecode chunk and write the
// output to a writer.
func (c *Chunk) Disassemble(output io.Writer) error {
	fmt.Fprintf(output, "== Disassembly of bytecode chunk at: %s ==\n\n", c.Location.String())

	if len(c.Instructions) == 0 {
		return nil
	}

	var offset int
	var instructionIndex int
	for {
		result, err := c.DisassembleInstruction(output, offset, instructionIndex)
		if err != nil {
			return err
		}
		offset = result
		instructionIndex++
		if offset >= len(c.Instructions) {
			break
		}
	}

	return nil
}

func (c *Chunk) DisassembleInstruction(output io.Writer, offset, instructionIndex int) (int, error) {
	fmt.Fprintf(output, "%04d  ", offset)
	opcodeByte := c.Instructions[offset]
	opcode := OpCode(opcodeByte)
	switch opcode {
	case RETURN, ADD, SUBTRACT,
		MULTIPLY, DIVIDE, EXPONENTIATE,
		NEGATE, NOT, BITWISE_NOT,
		TRUE, FALSE, NIL, POP:
		return c.disassembleOneByteInstruction(output, opcode.String(), offset, instructionIndex), nil
	case POP_N, SET_LOCAL8, GET_LOCAL8, PREP_LOCALS8:
		return c.disassembleNumericOperands(output, 1, 1, offset, instructionIndex)
	case PREP_LOCALS16, SET_LOCAL16, GET_LOCAL16:
		return c.disassembleNumericOperands(output, 1, 2, offset, instructionIndex)
	case LEAVE_SCOPE16:
		return c.disassembleNumericOperands(output, 2, 1, offset, instructionIndex)
	case LEAVE_SCOPE32:
		return c.disassembleNumericOperands(output, 2, 2, offset, instructionIndex)
	case CONSTANT8:
		return c.disassembleConstant(output, 2, offset, instructionIndex)
	case CONSTANT16:
		return c.disassembleConstant(output, 3, offset, instructionIndex)
	case CONSTANT32:
		return c.disassembleConstant(output, 5, offset, instructionIndex)
	default:
		c.printLineNumber(output, instructionIndex)
		c.dumpBytes(output, offset, 1)
		fmt.Fprintf(output, "unknown operation %d (0x%X)\n", opcodeByte, opcodeByte)
		return offset + 1, fmt.Errorf("unknown operation %d (0x%X) at offset %d (0x%X)", opcodeByte, opcodeByte, offset, offset)
	}
}

func (c *Chunk) dumpBytes(output io.Writer, offset, count int) {
	for i := offset; i < offset+count; i++ {
		fmt.Fprintf(output, "%02X ", c.Instructions[i])
	}

	for i := count; i < maxInstructionByteLength; i++ {
		fmt.Fprint(output, "   ")
	}
}

func (c *Chunk) disassembleOneByteInstruction(output io.Writer, name string, offset, instructionIndex int) int {
	c.printLineNumber(output, instructionIndex)
	c.dumpBytes(output, offset, 1)
	fmt.Fprintln(output, name)
	return offset + 1
}

func (c *Chunk) disassembleNumericOperands(output io.Writer, operands, operandBytes, offset, instructionIndex int) (int, error) {
	bytes := 1 + operands*operandBytes
	if result, err := c.checkBytes(output, offset, instructionIndex, bytes); err != nil {
		return result, err
	}

	opcode := OpCode(c.Instructions[offset])

	c.printLineNumber(output, instructionIndex)
	c.dumpBytes(output, offset, bytes)
	c.printOpCode(output, opcode)

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
		a := readFunc(c.Instructions[offset+1+i*operandBytes : offset+1+(i+1)*operandBytes])
		c.printNumField(output, a)
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

func (c *Chunk) disassembleConstant(output io.Writer, byteLength, offset, instructionIndex int) (int, error) {
	opcode := OpCode(c.Instructions[offset])

	if result, err := c.checkBytes(output, offset, instructionIndex, byteLength); err != nil {
		return result, err
	}

	var constantIndex int
	if byteLength == 2 {
		constantIndex = int(c.Instructions[offset+1])
	} else if byteLength == 3 {
		constantIndex = int(binary.BigEndian.Uint16(c.Instructions[offset+1 : offset+3]))
	} else if byteLength == 5 {
		constantIndex = int(binary.BigEndian.Uint32(c.Instructions[offset+1 : offset+5]))
	} else {
		panic(fmt.Sprintf("%d is not a valid byteLength for a constant opcode!", byteLength))
	}

	c.printLineNumber(output, instructionIndex)
	c.dumpBytes(output, offset, byteLength)
	c.printOpCode(output, opcode)

	if constantIndex >= len(c.Constants) {
		msg := fmt.Sprintf("invalid constant index %d (0x%X)", constantIndex, constantIndex)
		fmt.Fprintln(output, msg)
		return offset + byteLength, fmt.Errorf(msg)
	}
	constant := c.Constants[constantIndex]
	fmt.Fprintln(output, constant.Inspect())

	return offset + byteLength, nil
}

func (c *Chunk) checkBytes(output io.Writer, offset, instructionIndex, byteLength int) (int, error) {
	opcode := OpCode(c.Instructions[offset])
	if len(c.Instructions)-offset >= byteLength {
		return 0, nil
	}
	c.printLineNumber(output, instructionIndex)
	c.dumpBytes(output, offset, len(c.Instructions)-offset)
	c.printOpCode(output, opcode)
	msg := "not enough bytes"
	fmt.Fprintln(output, msg)
	return len(c.Instructions) - 1, fmt.Errorf(msg)
}

func (c *Chunk) printLineNumber(output io.Writer, instructionIndex int) {
	fmt.Fprintf(output, "%-8s", c.getLineNumberString(instructionIndex))
}

func (c *Chunk) getLineNumberString(instructionIndex int) string {
	currentLineNumber := c.LineInfoList.GetLineNumber(instructionIndex)
	if instructionIndex == 0 {
		return fmt.Sprintf("%d", currentLineNumber)
	}

	previousLineNumber := c.LineInfoList.GetLineNumber(instructionIndex - 1)
	if previousLineNumber == currentLineNumber {
		return "|"
	}

	return fmt.Sprintf("%d", currentLineNumber)
}

func (c *Chunk) printOpCode(output io.Writer, opcode OpCode) {
	fmt.Fprintf(output, "%-16s", opcode.String())
}

func (c *Chunk) printNumField(output io.Writer, n uint64) {
	fmt.Fprintf(output, "%-16d", n)
}
