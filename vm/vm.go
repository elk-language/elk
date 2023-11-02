// Package vm contains the Elk Virtual Machine.
// It interprets Elk Bytecode produced by
// the Elk compiler.
package vm

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/config"
	"github.com/elk-language/elk/value"
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
	bytecode   *value.BytecodeFunction
	ip         int           // Instruction pointer -- points to the next bytecode instruction
	stack      []value.Value // Value stack
	sp         int           // Stack pointer -- points to the offset where the next element will be pushed to
	fp         int           // Frame pointer -- points to the offset where the current frame starts
	err        value.Value   // The current error that is being thrown, nil if there has been no error or the error has already been handled
	callFrames []CallFrame   // Call stack
	Stdin      io.Reader     // standard output used by the VM
	Stdout     io.Writer     // standard input used by the VM
	Stderr     io.Writer     // standard error used by the VM
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
		stack:      make([]value.Value, VALUE_STACK_SIZE),
		callFrames: make([]CallFrame, 0, CALL_STACK_SIZE),
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}

	for _, opt := range opts {
		opt(vm)
	}

	return vm
}

// Execute the given bytecode chunk.
func (vm *VM) InterpretTopLevel(fn *value.BytecodeFunction) (value.Value, value.Value) {
	vm.bytecode = fn
	vm.ip = 0
	vm.push(value.GlobalObject)
	vm.push(value.RootModule)
	vm.run()
	return vm.peek(), vm.err
}

// Get the stored error.
func (vm *VM) Err() value.Value {
	return vm.err
}

// Get the value on top of the value stack.
func (vm *VM) StackTop() value.Value {
	return vm.peek()
}

// The main execution loop of the VM.
// Returns true when execution has been successful
// otherwise (in case of an error) returns false.
func (vm *VM) run() {
	for {
		// fmt.Println()
		// vm.bytecode.DisassembleInstruction(os.Stdout, vm.ip)
		// fmt.Println()

		instruction := bytecode.OpCode(vm.readByte())
		// BENCHMARK: replace with a jump table
		switch instruction {
		case bytecode.RETURN:
			if len(vm.callFrames) == 0 {
				return
			}
			returnValue := vm.pop()
			vm.restoreLastFrame()
			vm.push(returnValue)
		case bytecode.CONSTANT_BASE:
			vm.getLocal(1)
		case bytecode.SELF:
			vm.getLocal(0)
		case bytecode.DEF_CLASS:
			vm.defineClass()
		case bytecode.DEF_MODULE:
			vm.defineModule()
		case bytecode.ROOT:
			vm.push(value.RootModule)
		case bytecode.UNDEFINED:
			vm.push(value.Undefined)
		case bytecode.LOAD_VALUE8:
			vm.push(vm.readValue8())
		case bytecode.LOAD_VALUE16:
			vm.push(vm.readValue16())
		case bytecode.LOAD_VALUE32:
			vm.push(vm.readValue32())
		case bytecode.ADD:
			vm.add()
		case bytecode.SUBTRACT:
			vm.subtract()
		case bytecode.MULTIPLY:
			vm.multiply()
		case bytecode.DIVIDE:
			vm.divide()
		case bytecode.EXPONENTIATE:
			vm.exponentiate()
		case bytecode.NEGATE:
			vm.negate()
		case bytecode.NOT:
			vm.replace(value.ToNotBool(vm.peek()))
		case bytecode.TRUE:
			vm.push(value.True)
		case bytecode.FALSE:
			vm.push(value.False)
		case bytecode.NIL:
			vm.push(value.Nil)
		case bytecode.POP:
			vm.pop()
		case bytecode.POP_N:
			vm.popN(int(vm.readByte()))
		case bytecode.GET_LOCAL8:
			vm.getLocal(int(vm.readByte()))
		case bytecode.GET_LOCAL16:
			vm.getLocal(int(vm.readUint16()))
		case bytecode.SET_LOCAL8:
			vm.setLocal(int(vm.readByte()))
		case bytecode.SET_LOCAL16:
			vm.setLocal(int(vm.readUint16()))
		case bytecode.LEAVE_SCOPE16:
			vm.leaveScope(int(vm.readByte()), int(vm.readByte()))
		case bytecode.LEAVE_SCOPE32:
			vm.leaveScope(int(vm.readUint16()), int(vm.readUint16()))
		case bytecode.PREP_LOCALS8:
			vm.prepLocals(int(vm.readByte()))
		case bytecode.PREP_LOCALS16:
			vm.prepLocals(int(vm.readUint16()))
		case bytecode.GET_MOD_CONST8:
			vm.getModuleConstant(int(vm.readByte()))
		case bytecode.GET_MOD_CONST16:
			vm.getModuleConstant(int(vm.readUint16()))
		case bytecode.GET_MOD_CONST32:
			vm.getModuleConstant(int(vm.readUint32()))
		case bytecode.DEF_MOD_CONST8:
			vm.defModuleConstant(int(vm.readByte()))
		case bytecode.DEF_MOD_CONST16:
			vm.defModuleConstant(int(vm.readUint16()))
		case bytecode.DEF_MOD_CONST32:
			vm.defModuleConstant(int(vm.readUint32()))
		case bytecode.JUMP_UNLESS:
			if value.Falsy(vm.peek()) {
				jump := vm.readUint16()
				vm.ip += int(jump)
				break
			}
			vm.ip += 2
		case bytecode.JUMP_IF_NIL:
			if vm.peek() == value.Nil {
				jump := vm.readUint16()
				vm.ip += int(jump)
				break
			}
			vm.ip += 2
		case bytecode.JUMP_IF:
			if value.Truthy(vm.peek()) {
				jump := vm.readUint16()
				vm.ip += int(jump)
				break
			}
			vm.ip += 2
		case bytecode.JUMP:
			jump := vm.readUint16()
			vm.ip += int(jump)
		case bytecode.LOOP:
			jump := vm.readUint16()
			vm.ip -= int(jump)
		case bytecode.LBITSHIFT:
			vm.leftBitshift()
		case bytecode.LOGIC_LBITSHIFT:
			vm.logicalLeftBitshift()
		case bytecode.RBITSHIFT:
			vm.rightBitshift()
		case bytecode.LOGIC_RBITSHIFT:
			vm.logicalRightBitshift()
		case bytecode.BITWISE_AND:
			vm.bitwiseAnd()
		case bytecode.BITWISE_OR:
			vm.bitwiseOr()
		case bytecode.BITWISE_XOR:
			vm.bitwiseXor()
		case bytecode.MODULO:
			vm.modulo()
		case bytecode.EQUAL:
			vm.equal()
		case bytecode.NOT_EQUAL:
			vm.notEqual()
		case bytecode.STRICT_EQUAL:
			vm.strictEqual()
		case bytecode.GREATER:
			vm.greaterThan()
		case bytecode.GREATER_EQUAL:
			vm.greaterThanEqual()
		case bytecode.LESS:
			vm.lessThan()
		case bytecode.LESS_EQUAL:
			vm.lessThanEqual()
		default:
			panic(fmt.Sprintf("Unknown bytecode instruction: %#v", instruction))
		}

		// pp.Println(vm.stack[0:vm.sp])
		if vm.err != nil {
			return
		}
	}

}

// Restore the state of the VM to the last call frame.
func (vm *VM) restoreLastFrame() {
	lastIndex := len(vm.callFrames) - 1
	cf := vm.callFrames[lastIndex]
	// reset the popped call frame
	vm.callFrames[lastIndex] = CallFrame{}
	vm.callFrames = vm.callFrames[:lastIndex]

	vm.ip = cf.ip
	vm.fp = cf.fp
	vm.bytecode = cf.bytecode
}

// Treat the next 8 bits of bytecode as an index
// of a value and retrieve the value.
func (vm *VM) readValue8() value.Value {
	return vm.bytecode.Values[vm.readByte()]
}

// Treat the next 16 bits of bytecode as an index
// of a value and retrieve the value.
func (vm *VM) readValue16() value.Value {
	return vm.bytecode.Values[vm.readUint16()]
}

// Treat the next 32 bits of bytecode as an index
// of a value and retrieve the value.
func (vm *VM) readValue32() value.Value {
	return vm.bytecode.Values[vm.readUint32()]
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

// Define a new module
func (vm *VM) defineModule() bool {
	constantNameVal := vm.pop()
	parentModuleVal := vm.pop()
	bodyVal := vm.pop()

	constantName := constantNameVal.(value.Symbol)
	var parentModule *value.ModulelikeObject

	switch m := parentModuleVal.(type) {
	case *value.Class:
		parentModule = &m.ModulelikeObject
	case *value.Module:
		parentModule = &m.ModulelikeObject
	default:
		vm.throw(value.NewIsNotModuleError(parentModuleVal.Inspect()))
		return false
	}

	var module *value.Module

	if moduleVal, ok := parentModule.Constants.Get(constantName); ok {
		module, ok = moduleVal.(*value.Module)
		if !ok {
			vm.throw(value.NewRedefinedConstantError(parentModuleVal.Inspect(), constantName.Inspect()))
			return false
		}
	} else {
		module = value.NewModule()
		parentModule.AddConstant(constantName, module)
	}

	switch body := bodyVal.(type) {
	case *value.BytecodeFunction:
		vm.executeModuleBody(module, body)
	case value.UndefinedType:
		vm.push(module)
	default:
		panic(fmt.Sprintf("expected undefined or a bytecode function as the module body, got: %s", bodyVal.Inspect()))
	}

	return true
}

// Define a new class
func (vm *VM) defineClass() bool {
	superclassVal := vm.pop()
	constantNameVal := vm.pop()
	parentModuleVal := vm.pop()
	bodyVal := vm.pop()

	constantName := constantNameVal.(value.Symbol)
	var parentModule *value.ModulelikeObject

	switch mod := parentModuleVal.(type) {
	case *value.Class:
		parentModule = &mod.ModulelikeObject
	case *value.Module:
		parentModule = &mod.ModulelikeObject
	default:
		vm.throw(value.NewIsNotModuleError(parentModuleVal.Inspect()))
		return false
	}

	var class *value.Class

	if classVal, ok := parentModule.Constants.Get(constantName); ok {
		class, ok = classVal.(*value.Class)
		if !ok {
			vm.throw(value.NewRedefinedConstantError(parentModuleVal.Inspect(), constantName.Inspect()))
			return false
		}
		switch superclass := superclassVal.(type) {
		case *value.Class:
			if class.Parent != superclass {
				vm.throw(
					value.Errorf(
						value.TypeErrorClass,
						"superclass mismatch in %s, expected: %s, got: %s",
						class.Name,
						class.Parent.Name,
						superclass.Name,
					),
				)
				return false
			}
		case value.UndefinedType:
		default:
			vm.throw(
				value.Errorf(
					value.TypeErrorClass,
					"`%s` can't be used as a superclass", superclass.Inspect(),
				),
			)
			return false
		}
	} else {
		class = value.NewClass()
		switch superclass := superclassVal.(type) {
		case *value.Class:
			class.Parent = superclass
		case value.UndefinedType:
		default:
			vm.throw(
				value.Errorf(
					value.TypeErrorClass,
					"`%s` can't be used as a superclass", superclass.Inspect(),
				),
			)
			return false
		}
		parentModule.AddConstant(constantName, class)
	}

	switch body := bodyVal.(type) {
	case *value.BytecodeFunction:
		vm.executeModuleBody(class, body)
	case value.UndefinedType:
		vm.push(class)
	default:
		panic(fmt.Sprintf("expected undefined or a bytecode function as the class body, got: %s", bodyVal.Inspect()))
	}

	return true
}

func (vm *VM) addCallFrame(cf CallFrame) {
	if len(vm.callFrames) == CALL_STACK_SIZE {
		panic(fmt.Sprintf("Stack overflow: %d", CALL_STACK_SIZE))
	}

	vm.callFrames = append(vm.callFrames, cf)
}

// set up the vm to execute a module body
func (vm *VM) executeModuleBody(module value.Value, body *value.BytecodeFunction) {
	// preserve the current state of the vm in a call frame
	vm.addCallFrame(
		CallFrame{
			bytecode: vm.bytecode,
			ip:       vm.ip,
			fp:       vm.fp,
		},
	)

	vm.bytecode = body
	vm.fp = vm.sp
	vm.ip = 0
	// set module as `self`
	vm.push(module)
	// set module as constant base
	vm.push(module)
}

// Set a local variable or value.
func (vm *VM) setLocal(index int) {
	vm.stack[vm.fp+index] = vm.peek()
}

// Read a local variable or value.
func (vm *VM) getLocal(index int) {
	vm.push(vm.stack[vm.fp+index])
}

// Pop a module off the stack and look for a constant with the given name.
func (vm *VM) getModuleConstant(nameIndex int) bool {
	symbol := vm.bytecode.Values[nameIndex].(value.Symbol)
	mod := vm.pop()
	var constants value.SimpleSymbolMap

	switch m := mod.(type) {
	case *value.Class:
		constants = m.Constants
	case *value.Module:
		constants = m.Constants
	default:
		vm.throw(value.Errorf(value.TypeErrorClass, "`%s` is not a module", mod.Inspect()))
		return false
	}

	val, ok := constants.Get(symbol)
	if !ok {
		vm.throw(value.Errorf(value.NoConstantErrorClass, "%s doesn't have a constant named `%s`", mod.Inspect(), symbol.Inspect()))
	}

	vm.push(val)
	return true
}

// Pop two values off the stack and define a constant with the given name.
func (vm *VM) defModuleConstant(nameIndex int) bool {
	symbol := vm.bytecode.Values[nameIndex].(value.Symbol)
	mod := vm.pop()
	var constants value.SimpleSymbolMap

	switch m := mod.(type) {
	case *value.Class:
		constants = m.Constants
	case *value.Module:
		constants = m.Constants
	default:
		vm.throw(value.NewIsNotModuleError(mod.Inspect()))
		return false
	}

	val := vm.peek()
	if constants.Has(symbol) {
		vm.throw(value.NewRedefinedConstantError(mod.Inspect(), symbol.Inspect()))
		return false
	}
	constants.Set(symbol, val)
	return true
}

// Leave a local scope and pop all local variables associated with it.
func (vm *VM) leaveScope(lastLocalIndex, varsToPop int) {
	firstLocalIndex := lastLocalIndex - varsToPop
	for i := lastLocalIndex; i > firstLocalIndex; i-- {
		vm.stack[i] = nil
	}
}

// Register slots for local variables and values.
func (vm *VM) prepLocals(count int) {
	vm.sp += count
}

// Push an element on top of the value stack.
func (vm *VM) push(val value.Value) {
	vm.stack[vm.sp] = val
	vm.sp++
}

// Pop an element off the value stack.
func (vm *VM) pop() value.Value {
	if vm.sp == 0 {
		panic("tried to pop when there are no elements on the value stack!")
	}

	vm.sp--
	val := vm.stack[vm.sp]
	vm.stack[vm.sp] = nil
	return val
}

// Pop n elements off the value stack.
func (vm *VM) popN(n int) {
	if vm.sp-n < 0 {
		panic("tried to pop more elements than are available on the value stack!")
	}

	for i := vm.sp; i > vm.sp-n; i-- {
		vm.stack[vm.sp] = nil
	}
	vm.sp -= n
}

// Replaces the value on top of the stack without popping it.
func (vm *VM) replace(val value.Value) {
	vm.stack[vm.sp-1] = val
}

// Return the element on top of the stack
// without popping it.
func (vm *VM) peek() value.Value {
	if vm.sp == 0 {
		panic("tried to peek when there are no elements on the value stack!")
	}

	return vm.stack[vm.sp-1]
}

// Negate the element on top of the stack
func (vm *VM) negate() bool {
	operand := vm.peek()
	result, builtin := value.Negate(operand)
	if !builtin {
		vm.throw(value.NewNoMethodError("-", operand))
		return false
	}

	vm.replace(result)
	return true
}

type binaryOperationWithoutErrFunc func(left value.Value, right value.Value) (value.Value, bool)

func (vm *VM) binaryOperationWithoutErr(fn binaryOperationWithoutErrFunc, methodName string) bool {
	right := vm.pop()
	left := vm.peek()

	result, builtin := fn(left, right)
	if !builtin {
		vm.throw(value.NewNoMethodError(methodName, left))
		return false
	}
	vm.replace(result)
	return true
}

type binaryOperationFunc func(left value.Value, right value.Value) (value.Value, *value.Error, bool)

func (vm *VM) binaryOperation(fn binaryOperationFunc, methodName string) bool {
	right := vm.pop()
	left := vm.peek()

	result, err, builtin := fn(left, right)
	if !builtin {
		vm.throw(value.NewNoMethodError(methodName, left))
		return false
	}
	if err != nil {
		vm.throw(err)
		return false
	}
	vm.replace(result)
	return true
}

// Perform a bitwise AND and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) bitwiseAnd() bool {
	return vm.binaryOperation(value.BitwiseAnd, "&")
}

// Perform a bitwise OR and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) bitwiseOr() bool {
	return vm.binaryOperation(value.BitwiseOr, "|")
}

// Perform a bitwise XOR and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) bitwiseXor() bool {
	return vm.binaryOperation(value.BitwiseXor, "^")
}

// Perform modulo and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) modulo() bool {
	return vm.binaryOperation(value.Modulo, "%")
}

// Check whether two top elements on the stack are equal and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) equal() bool {
	return vm.binaryOperationWithoutErr(value.Equal, "==")
}

// Check whether two top elements on the stack are not and equal push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) notEqual() bool {
	return vm.binaryOperationWithoutErr(value.NotEqual, "==")
}

// Check whether two top elements on the stack are strictly equal push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) strictEqual() bool {
	return vm.binaryOperationWithoutErr(value.StrictEqual, "===")
}

// Check whether two top elements on the stack are strictly not equal push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) strictNotEqual() bool {
	return vm.binaryOperationWithoutErr(value.StrictNotEqual, "===")
}

// Check whether the first operand is greater than the second and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) greaterThan() bool {
	return vm.binaryOperation(value.GreaterThan, ">")
}

// Check whether the first operand is greater than or equal to the second and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) greaterThanEqual() bool {
	return vm.binaryOperation(value.GreaterThanEqual, ">=")
}

// Check whether the first operand is less than the second and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) lessThan() bool {
	return vm.binaryOperation(value.LessThan, "<")
}

// Check whether the first operand is less than or equal to the second and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) lessThanEqual() bool {
	return vm.binaryOperation(value.LessThanEqual, "<=")
}

// Perform a left bitshift and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) leftBitshift() bool {
	return vm.binaryOperation(value.LeftBitshift, "<<")
}

// Perform a logical left bitshift and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) logicalLeftBitshift() bool {
	return vm.binaryOperation(value.LogicalLeftBitshift, "<<<")
}

// Perform a right bitshift and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) rightBitshift() bool {
	return vm.binaryOperation(value.RightBitshift, ">>")
}

// Perform a logical right bitshift and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) logicalRightBitshift() bool {
	return vm.binaryOperation(value.LogicalRightBitshift, ">>>")
}

// Add two operands together and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) add() bool {
	return vm.binaryOperation(value.Add, "+")
}

// Subtract two operands and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) subtract() bool {
	return vm.binaryOperation(value.Subtract, "-")
}

// Multiply two operands together and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) multiply() bool {
	return vm.binaryOperation(value.Multiply, "*")
}

// Divide two operands and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) divide() bool {
	return vm.binaryOperation(value.Divide, "/")
}

// Exponentiate two operands and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) exponentiate() bool {
	return vm.binaryOperation(value.Exponentiate, "**")
}

// Throw an error and attempt to find code
// that catches it.
func (vm *VM) throw(err value.Value) {
	vm.err = err
}
