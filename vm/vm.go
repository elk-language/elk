// Package vm contains the Elk Virtual Machine.
// It interprets Elk Bytecode produced by
// the Elk compiler.
package vm

import (
	"fmt"
	"io"
	"os"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/object"
)

type Result uint8 // Result of the interpreted program

const (
	RESULT_OK            Result = iota // the program has successfully finished
	RESULT_COMPILE_ERROR               // the program couldn't be compiled
	RESULT_RUNTIME_ERROR               // the program has halted because of a runtime error
)

// A single instance of the Elk Virtual Machine.
type VM struct {
	bytecodeChunk *bytecode.Chunk
	ip            int       // Instruction pointer
	Stdin         io.Reader // standard output used by the VM
	Stdout        io.Writer // standard input used by the VM
	Stderr        io.Writer // standard error used by the VM
}

type Option func(*VM) // constructor option function

// Assign the given io.Reader as the Stdin of the VM.
func WithStdin(stdin io.Reader) Option {
	return func(vm *VM) {
		vm.Stdin = stdin
	}
}

// Assign the given io.Writer as the Stdout of the VM.
func WithStdout(stdout io.Writer) Option {
	return func(vm *VM) {
		vm.Stdout = stdout
	}
}

// Assign the given io.Writer as the Stderr of the VM.
func WithStderr(stderr io.Writer) Option {
	return func(vm *VM) {
		vm.Stderr = stderr
	}
}

// Create a new VM instance.
func New(opts ...Option) *VM {
	vm := &VM{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	for _, opt := range opts {
		opt(vm)
	}

	return vm
}

// Execute the given bytecode chunk.
func (vm *VM) InterpretBytecode(chunk *bytecode.Chunk) Result {
	vm.bytecodeChunk = chunk
	vm.ip = 0
	return vm.run()
}

// The main execution loop of the VM.
func (vm *VM) run() Result {
	for {
		vm.bytecodeChunk.DisassembleInstruction(vm.Stdout, vm.ip)

		instruction := bytecode.OpCode(vm.readByte())
		switch instruction {
		case bytecode.RETURN:
			return RESULT_OK
		case bytecode.CONSTANT8:
			index := vm.readByte()
			constant := vm.bytecodeChunk.Constants[index]
			fmt.Fprintln(vm.Stdout, object.Inspect(constant))
		case bytecode.CONSTANT16:
			index := vm.readUint16()
			constant := vm.bytecodeChunk.Constants[index]
			fmt.Fprintln(vm.Stdout, object.Inspect(constant))
		case bytecode.CONSTANT32:
			index := vm.readUint32()
			constant := vm.bytecodeChunk.Constants[index]
			fmt.Fprintln(vm.Stdout, object.Inspect(constant))
		default:
			return RESULT_RUNTIME_ERROR
		}
	}
}

// Read the next byte of code
func (vm *VM) readByte() byte {
	byt := vm.bytecodeChunk.Instructions[vm.ip]
	vm.ip++
	return byt
}

// Read the next 2 bytes of code
func (vm *VM) readUint16() uint16 {
	result := uint16(vm.bytecodeChunk.Instructions[vm.ip])<<8 |
		uint16(vm.bytecodeChunk.Instructions[vm.ip+1])

	vm.ip += 2

	return result
}

// Read the next 4 bytes of code
func (vm *VM) readUint32() uint32 {
	result := uint32(vm.bytecodeChunk.Instructions[vm.ip])<<24 |
		uint32(vm.bytecodeChunk.Instructions[vm.ip+1])<<16 |
		uint32(vm.bytecodeChunk.Instructions[vm.ip+2])<<8 |
		uint32(vm.bytecodeChunk.Instructions[vm.ip+3])

	vm.ip += 4

	return result
}
