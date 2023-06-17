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
func (c *Chunk) AddInstruction(op OpCode, bytes ...byte) {
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
		offset, err := c.disassembleInstruction(output, offset)
		if err != nil {
			return err
		}
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
	case CONSTANT:
		return c.disassembleConstant(output, offset), nil
	default:
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
	c.dumpBytes(output, offset, 1)
	fmt.Fprintln(output, name)
	return offset + 1
}

func (c *Chunk) disassembleConstant(output io.Writer, offset int) int {
	opcode := OpCode(c.Instructions[offset])
	constantIndex := c.Instructions[offset+1]
	c.dumpBytes(output, offset, 2)
	c.printOpCode(output, opcode)

	if int(constantIndex) >= len(c.Constants) {
		fmt.Fprintf(output, "invalid constant index %d (0x%X)", constantIndex, constantIndex)
		return offset + 2
	}
	constant := c.Constants[constantIndex]
	fmt.Fprintln(output, object.Inspect(constant))

	return offset + 2
}

func (c *Chunk) printOpCode(output io.Writer, opcode OpCode) {
	fmt.Fprintf(output, "%- 16s", opcode.String())
}
