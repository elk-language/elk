// Package vm contains the Elk Virtual Machine.
// It interprets Elk Bytecode produced by
// the Elk compiler.
package vm

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/config"
	"github.com/elk-language/elk/object"
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

// A single instance of the Elk Virtual Machine.
type VM struct {
	bytecode *bytecode.Chunk
	ip       int            // Instruction pointer -- points to the next bytecode instruction
	stack    []object.Value // Value stack
	sp       int            // Stack pointer -- points to the offset where the next element will be pushed to
	err      *object.Error  // The current error that is being thrown, nil if there has been no error or the error has already been handled
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
func (vm *VM) InterpretBytecode(chunk *bytecode.Chunk) *object.Error {
	vm.bytecode = chunk
	vm.ip = 0
	vm.run()
	return vm.err
}

// The main execution loop of the VM.
// Returns true when execution has been successful
// otherwise (in case of an error) returns false.
func (vm *VM) run() bool {
	for vm.ip < len(vm.bytecode.Instructions) {
		// fmt.Println()
		// vm.bytecode.DisassembleInstruction(os.Stdout, vm.ip)
		// fmt.Println()

		instruction := bytecode.OpCode(vm.readByte())
		// BENCHMARK: replace with a jump table
		switch instruction {
		case bytecode.RETURN:
			return true
		case bytecode.CONSTANT8:
			vm.push(vm.readConstant8())
		case bytecode.CONSTANT16:
			vm.push(vm.readConstant16())
		case bytecode.CONSTANT32:
			vm.push(vm.readConstant32())
		case bytecode.ADD:
			vm.add()
		case bytecode.NEGATE:
			vm.negate()
		case bytecode.NOT:
			vm.replace(object.ToNotBool(vm.peek()))
		case bytecode.TRUE:
			vm.push(object.True)
		case bytecode.FALSE:
			vm.push(object.False)
		case bytecode.NIL:
			vm.push(object.Nil)
		default:
			return false
		}

		// pp.Println(vm.stack[0:vm.sp])
	}

	return true
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
	// BENCHMARK: compare manual bit shifts
	result := binary.BigEndian.Uint16(vm.bytecode.Instructions[vm.ip : vm.ip+2])
	vm.ip += 2

	return result
}

// Read the next 4 bytes of code
func (vm *VM) readUint32() uint32 {
	// BENCHMARK: compare manual bit shifts
	result := binary.BigEndian.Uint32(vm.bytecode.Instructions[vm.ip : vm.ip+4])

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
	if vm.sp == 0 {
		panic("tried to pop when there are no elements on the value stack!")
	}

	vm.sp--
	val := vm.stack[vm.sp]
	vm.stack[vm.sp] = nil
	return val
}

// Replaces the value on top of the stack without popping it.
func (vm *VM) replace(val object.Value) {
	vm.stack[vm.sp-1] = val
}

// Return the element on top of the stack
// without popping it.
func (vm *VM) peek() object.Value {
	if vm.sp == 0 {
		panic("tried to peek when there are no elements on the value stack!")
	}

	return vm.stack[vm.sp-1]
}

// Negate the element on top of the stack
func (vm *VM) negate() bool {
	operand := vm.peek()
	switch o := operand.(type) {
	case object.Float64:
		vm.replace(-o)
	case object.Float32:
		vm.replace(-o)
	case object.Float:
		vm.replace(-o)
	case *object.BigFloat:
		vm.replace(o.Negate())
	case object.Int64:
		vm.replace(-o)
	case object.Int32:
		vm.replace(-o)
	case object.Int16:
		vm.replace(-o)
	case object.Int8:
		vm.replace(-o)
	case object.UInt64:
		vm.replace(-o)
	case object.UInt32:
		vm.replace(-o)
	case object.UInt16:
		vm.replace(-o)
	case object.UInt8:
		vm.replace(-o)
	case object.SmallInt:
		vm.replace(-o)
	case *object.BigInt:
		vm.replace(o.Negate())
	default:
		vm.throw(object.NewNoMethodError("-", o))
		return false
	}

	return true
}

// Add two operands together and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) add() bool {
	right := vm.pop()
	left := vm.peek()
	// TODO: Implement SmallInt, BigInt and other type addition
	switch l := left.(type) {
	case object.Float:
		result, err := l.Add(right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case *object.BigFloat:
		result, err := l.Add(right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.Float64:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.Float32:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.Int64:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.Int32:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.Int16:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.Int8:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.UInt64:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.UInt32:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.UInt16:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.UInt8:
		result, err := object.StrictNumericAdd(l, right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.String:
		result, err := l.Concat(right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	case object.Char:
		result, err := l.Concat(right)
		if err != nil {
			vm.throw(err)
			return false
		}
		vm.replace(result)
	default:
		vm.throw(object.NewNoMethodError("+", left))
		return false
	}

	return true
}

// Throw an error and attempt to find code
// that catches it.
func (vm *VM) throw(err *object.Error) {
	vm.err = err
}
