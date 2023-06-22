// Package vm contains the Elk Virtual Machine.
// It interprets Elk Bytecode produced by
// the Elk compiler.
package vm

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/object"
	"github.com/k0kubun/pp"
)

// BENCHMARK: compare with a dynamically allocated array
var VALUE_STACK_SIZE int

func init() {
	valueStackSize, ok := os.LookupEnv("ELK_VALUE_STACK_SIZE")
	if !ok {
		VALUE_STACK_SIZE = 1024 // 1KB by default
		return
	}

	valInt, err := strconv.Atoi(valueStackSize)
	if err != nil {
		panic(fmt.Sprintf("invalid value for ELK_VALUE_STACK_SIZE, expected int, got %v", valueStackSize))
	}
	VALUE_STACK_SIZE = valInt
}

type Result uint8 // Result of the interpreted program

const (
	RESULT_OK            Result = iota // the program has successfully finished
	RESULT_COMPILE_ERROR               // the program couldn't be compiled
	RESULT_RUNTIME_ERROR               // the program has halted because of a runtime error
)

// A single instance of the Elk Virtual Machine.
type VM struct {
	bytecode *bytecode.Chunk
	ip       int            // Instruction pointer -- points to the next bytecode instruction
	stack    []object.Value // value stack
	sp       int            // Stack pointer -- points to the offset where the next element will be pushed to
	Stdin    io.Reader      // standard output used by the VM
	Stdout   io.Writer      // standard input used by the VM
	Stderr   io.Writer      // standard error used by the VM
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
		stack:  make([]object.Value, VALUE_STACK_SIZE),
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
	vm.bytecode = chunk
	vm.ip = 0
	return vm.run()
}

// The main execution loop of the VM.
func (vm *VM) run() Result {
	for vm.ip < len(vm.bytecode.Instructions) {
		fmt.Println()
		vm.bytecode.DisassembleInstruction(os.Stdout, vm.ip)
		fmt.Println()

		instruction := bytecode.OpCode(vm.readByte())
		switch instruction {
		case bytecode.RETURN:
			vm.pop()
			return RESULT_OK
		case bytecode.CONSTANT8:
			index := vm.readByte()
			vm.push(vm.bytecode.Constants[index])
		case bytecode.CONSTANT16:
			index := vm.readUint16()
			vm.push(vm.bytecode.Constants[index])
		case bytecode.CONSTANT32:
			index := vm.readUint32()
			vm.push(vm.bytecode.Constants[index])
		case bytecode.ADD:
			left := vm.pop()
			right := vm.pop()
			switch l := left.(type) {
			case object.Int64:
				r, ok := right.(object.Int64)
				if !ok {
					panic(fmt.Sprintf("can't add Int64 to %s", r.Inspect()))
				}
				vm.push(l + r)
			case object.Int32:
				r, ok := right.(object.Int32)
				if !ok {
					panic(fmt.Sprintf("can't add Int32 to %s", r.Inspect()))
				}
				vm.push(l + r)
			case object.Int16:
				r, ok := right.(object.Int16)
				if !ok {
					panic(fmt.Sprintf("can't add Int16 to %s", r.Inspect()))
				}
				vm.push(l + r)
			case object.Int8:
				r, ok := right.(object.Int8)
				if !ok {
					panic(fmt.Sprintf("can't add Int8 to %s", r.Inspect()))
				}
				vm.push(l + r)
			default:
				panic(fmt.Sprintf("adding %s and %s has not been implemented yet", left.Inspect(), right.Inspect()))
			}
		default:
			return RESULT_RUNTIME_ERROR
		}

		pp.Println(vm.stack[0:vm.sp])
	}

	return RESULT_OK
}

// Read the next byte of code
func (vm *VM) readByte() byte {
	// BENCHMARK: compare pointer arithmetic to offsets
	byt := vm.bytecode.Instructions[vm.ip]
	vm.ip++
	return byt
}

// Read the next 2 bytes of code
func (vm *VM) readUint16() uint16 {
	// BENCHMARK: compare binary.BigEndian.Uint16
	result := uint16(vm.bytecode.Instructions[vm.ip])<<8 |
		uint16(vm.bytecode.Instructions[vm.ip+1])

	vm.ip += 2

	return result
}

// Read the next 4 bytes of code
func (vm *VM) readUint32() uint32 {
	// BENCHMARK: compare binary.BigEndian.Uint32
	result := uint32(vm.bytecode.Instructions[vm.ip])<<24 |
		uint32(vm.bytecode.Instructions[vm.ip+1])<<16 |
		uint32(vm.bytecode.Instructions[vm.ip+2])<<8 |
		uint32(vm.bytecode.Instructions[vm.ip+3])

	vm.ip += 4

	return result
}

// Push an element on top of the value stack.
func (vm *VM) push(val object.Value) {
	vm.stack[vm.sp] = val
	vm.sp++
}

// Pop an element off the value stack.
func (vm *VM) pop() object.Value {
	vm.sp--
	return vm.stack[vm.sp]
}
