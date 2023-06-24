// Package vm contains the Elk Virtual Machine.
// It interprets Elk Bytecode produced by
// the Elk compiler.
package vm

import (
	"fmt"
	"io"
	"os"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/config"
	"github.com/elk-language/elk/object"
	"github.com/k0kubun/pp"
)

// BENCHMARK: compare with a dynamically allocated array
var VALUE_STACK_SIZE int

func init() {
	val, ok := config.IntFromEnvVar("ELK_VALUE_STACK_SIZE")
	if !ok {
		VALUE_STACK_SIZE = 1024 // 1KB by default
		return
	}

	VALUE_STACK_SIZE = val
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
			vm.push(vm.readConstant8())
		case bytecode.CONSTANT16:
			vm.push(vm.readConstant16())
		case bytecode.CONSTANT32:
			vm.push(vm.readConstant32())
		case bytecode.ADD:
			vm.Add()
		case bytecode.NEGATE:
			vm.Negate()
		default:
			return RESULT_RUNTIME_ERROR
		}

		pp.Println(vm.stack[0:vm.sp])
	}

	return RESULT_OK
}

// Treat the next 8 bits of bytecode as an index
// of a constant and retrieve the constant.
func (vm *VM) readConstant8() object.Value {
	return vm.bytecode.Constants[vm.readByte()]
}

// Treat the next 16 bits of bytecode as an index
// of a constant and retrieve the constant.
func (vm *VM) readConstant16() object.Value {
	return vm.bytecode.Constants[vm.readUint16()]
}

// Treat the next 32 bits of bytecode as an index
// of a constant and retrieve the constant.
func (vm *VM) readConstant32() object.Value {
	return vm.bytecode.Constants[vm.readUint32()]
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

// Negate the element on top of the stack
func (vm *VM) Negate() {
	operand := vm.pop()
	switch o := operand.(type) {
	case object.Float64:
		vm.push(-o)
	case object.Float32:
		vm.push(-o)
	case object.Float:
		vm.push(-o)
	case object.Int64:
		vm.push(-o)
	case object.Int32:
		vm.push(-o)
	case object.Int16:
		vm.push(-o)
	case object.Int8:
		vm.push(-o)
	case object.UInt64:
		vm.push(-o)
	case object.UInt32:
		vm.push(-o)
	case object.UInt16:
		vm.push(-o)
	case object.UInt8:
		vm.push(-o)
	case object.SmallInt:
		vm.push(-o)
	case *object.BigInt:
		vm.push(o.Neg())
	default:
		panic(fmt.Sprintf("negating %s has not been implemented yet", o.Inspect()))
	}
}

// Add two operand together and push the result to the stack.
func (vm *VM) Add() {
	right := vm.pop()
	left := vm.pop()
	switch l := left.(type) {
	case object.Float64:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			panic(err)
		}
		vm.push(result)
	case object.Float32:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			panic(err)
		}
		vm.push(result)
	case object.Int64:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			panic(err)
		}
		vm.push(result)
	case object.Int32:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			panic(err)
		}
		vm.push(result)
	case object.Int16:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			panic(err)
		}
		vm.push(result)
	case object.Int8:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			panic(err)
		}
		vm.push(result)
	case object.UInt64:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			panic(err)
		}
		vm.push(result)
	case object.UInt32:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			panic(err)
		}
		vm.push(result)
	case object.UInt16:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			panic(err)
		}
		vm.push(result)
	case object.UInt8:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			panic(err)
		}
		vm.push(result)
	case object.String:
		result, err := l.Concat(right)
		if err != nil {
			panic(err)
		}
		vm.push(result)
	default:
		panic(fmt.Sprintf("adding %s and %s has not been implemented yet", left.Inspect(), right.Inspect()))
	}
}
