// Package vm contains the Elk Virtual Machine.
// It interprets Elk Bytecode produced by
// the Elk compiler.
package vm

import (
	"github.com/elk-language/elk/bytecode"
)

// A single instance of the Elk Virtual Machine.
type VM struct {
	BytecodeChunk      *bytecode.Chunk
	instructionPointer int
}

// The main execution loop of the VM.
func (vm *VM) Run() {
	for {
		instruction := bytecode.OpCode(vm.readByte())
		switch instruction {
		case bytecode.RETURN:
			return
		default:
			return
		}
	}
}

// Read the next byte of code
func (vm *VM) readByte() byte {
	byt := vm.BytecodeChunk.Instructions[vm.instructionPointer]
	vm.instructionPointer++
	return byt
}
