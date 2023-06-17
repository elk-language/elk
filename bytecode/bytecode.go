// Package bytecode implements types
// and constants that make up Elk
// bytecode.
package bytecode

import (
	"fmt"
	"io"
	"os"

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

// Add a constant to the constant pool.
// Returns the index of the constant.
func (c *Chunk) AddConstant(obj object.Value) int {
	c.Constants = append(c.Constants, obj)
	return len(c.Constants) - 1
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

	offset := 0
	for {
		result, err := c.disassembleInstruction(output, offset)
		if err != nil {
			return err
		}
		offset = result
		if offset >= len(c.Instructions) {
			break
		}
	}

	return nil
}

func (c *Chunk) disassembleInstruction(output io.Writer, offset int) (int, error) {
	fmt.Fprintf(output, "%04d  ", offset)
	opcodeByte := c.Instructions[offset]
	opcode := OpCode(opcodeByte)
	switch opcode {
	case RETURN:
		return c.disassembleOneByteInstruction(output, opcode.String(), offset), nil
	case CONSTANT8:
		return c.disassembleConstant(output, 2, offset)
	case CONSTANT16:
		return c.disassembleConstant(output, 3, offset)
	case CONSTANT32:
		return c.disassembleConstant(output, 5, offset)
	default:
		c.printLineNumber(output, offset)
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

func (c *Chunk) disassembleOneByteInstruction(output io.Writer, name string, offset int) int {
	c.printLineNumber(output, offset)
	c.dumpBytes(output, offset, 1)
	fmt.Fprintln(output, name)
	return offset + 1
}

func (c *Chunk) disassembleConstant(output io.Writer, byteLength, offset int) (int, error) {
	opcode := OpCode(c.Instructions[offset])

	if result, err := c.checkBytes(output, offset, byteLength); err != nil {
		return result, err
	}

	var constantIndex int
	if byteLength == 2 {
		constantIndex = int(c.Instructions[offset+1])
	} else if byteLength == 3 {
		constantIndex = (int(c.Instructions[offset+1])<<8 | int(c.Instructions[offset+2]))
	} else if byteLength == 5 {
		constantIndex = (int(c.Instructions[offset+1])<<24 | int(c.Instructions[offset+2])<<16 | int(c.Instructions[offset+3])<<8 | int(c.Instructions[offset+4]))
	} else {
		panic(fmt.Sprintf("%d is not a valid byteLength for a constant opcode!", byteLength))
	}

	c.printLineNumber(output, offset)
	c.dumpBytes(output, offset, byteLength)
	c.printOpCode(output, opcode)

	if int(constantIndex) >= len(c.Constants) {
		msg := fmt.Sprintf("invalid constant index %d (0x%X)", constantIndex, constantIndex)
		fmt.Fprintln(output, msg)
		return offset + byteLength, fmt.Errorf(msg)
	}
	constant := c.Constants[constantIndex]
	fmt.Fprintln(output, object.Inspect(constant))

	return offset + byteLength, nil
}

func (c *Chunk) checkBytes(output io.Writer, offset, byteLength int) (int, error) {
	opcode := OpCode(c.Instructions[offset])
	if len(c.Instructions)-offset >= byteLength {
		return 0, nil
	}
	c.printLineNumber(output, offset)
	c.dumpBytes(output, offset, len(c.Instructions)-offset)
	c.printOpCode(output, opcode)
	msg := "not enough bytes"
	fmt.Fprintln(output, msg)
	return len(c.Instructions) - 1, fmt.Errorf(msg)
}

func (c *Chunk) printLineNumber(output io.Writer, offset int) {
	fmt.Fprintf(output, "%- 8s", c.getLineNumberString(offset))
}

func (c *Chunk) getLineNumberString(offset int) string {
	currentLineNumber := c.LineInfoList.GetLineNumber(offset)
	if offset == 0 {
		return fmt.Sprintf("%d", currentLineNumber)
	}

	previousLineNumber := c.LineInfoList.GetLineNumber(offset - 1)
	if previousLineNumber == currentLineNumber {
		return "|"
	}

	return fmt.Sprintf("%d", currentLineNumber)
}

func (c *Chunk) printOpCode(output io.Writer, opcode OpCode) {
	fmt.Fprintf(output, "%- 16s", opcode.String())
}
