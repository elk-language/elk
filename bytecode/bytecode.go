// Package bytecode implements types
// and constants that make up Elk
// bytecode.
package bytecode

import (
	"fmt"
	"io"
	"os"

	"github.com/elk-language/elk/position"
)

// A single unit of Elk bytecode.
type Chunk struct {
	Instructions []byte
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
	instruction := OpCode(c.Instructions[offset])
	switch instruction {
	case RETURN:
		return c.disassembleOneByteInstruction(output, "RETURN", offset), nil
	default:
		c.dumpBytes(output, offset, 1)
		fmt.Fprintf(output, "unknown operation %d (0x%X)\n", instruction, instruction)
		return offset + 1, fmt.Errorf("unknown operation %d (0x%X) at offset %d (0x%X)", instruction, instruction, offset, offset)
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
