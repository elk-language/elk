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
	if ok {
		VALUE_STACK_SIZE = val
	} else {
		VALUE_STACK_SIZE = 1024 // 1KB by default
	}
}

// VM mode
type mode uint8

const (
	normalMode           mode = iota
	singleMethodCallMode      // the VM should halt after executing a single method
)

// A single instance of the Elk Virtual Machine.
type VM struct {
	bytecode   *BytecodeMethod
	ip         int           // Instruction pointer -- points to the next bytecode instruction
	sp         int           // Stack pointer -- points to the offset where the next element will be pushed to
	fp         int           // Frame pointer -- points to the offset where the current frame starts
	stack      []value.Value // Value stack
	callFrames []CallFrame   // Call stack
	err        value.Value   // The current error that is being thrown, nil if there has been no error or the error has already been handled
	Stdin      io.Reader     // standard output used by the VM
	Stdout     io.Writer     // standard input used by the VM
	Stderr     io.Writer     // standard error used by the VM
	mode       mode
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
func (vm *VM) InterpretTopLevel(fn *BytecodeMethod) (value.Value, value.Value) {
	vm.bytecode = fn
	vm.ip = 0
	vm.push(value.GlobalObject)
	vm.push(value.RootModule)
	vm.push(value.GlobalObjectSingletonClass)
	vm.run()
	return vm.peek(), vm.err
}

// Execute the given bytecode chunk.
func (vm *VM) InterpretREPL(fn *BytecodeMethod) (value.Value, value.Value) {
	vm.bytecode = fn
	vm.ip = 0
	if vm.sp == 0 {
		// populate the predeclared local variables
		vm.push(value.GlobalObject)               // populate self
		vm.push(value.RootModule)                 // populate constant container
		vm.push(value.GlobalObjectSingletonClass) // populate method container
	} else {
		// pop the return value of the last run
		vm.pop()
	}
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

func (vm *VM) Stack() []value.Value {
	return vm.stack[:vm.sp]
}

func (vm *VM) throwIfErr(err value.Value) {
	if err != nil {
		vm.throw(err)
	}
}

// Call an Elk method from Go code, preserving the state of the VM.
func (vm *VM) CallMethod(name value.Symbol, args ...value.Value) (value.Value, value.Value) {
	self := args[0]
	class := self.DirectClass()
	method := class.LookupMethod(name)
	if method == nil {
		return nil, value.NewNoMethodError(name.ToString(), self)
	}
	if method.ParameterCount() != len(args) {
		return nil, value.NewWrongArgumentCountError(len(args), method.ParameterCount())
	}

	switch m := method.(type) {
	case *BytecodeMethod:
		vm.createCurrentCallFrame()
		vm.bytecode = m
		vm.fp = vm.sp
		vm.ip = 0
		vm.mode = singleMethodCallMode
		for _, arg := range args {
			vm.push(arg)
		}
		vm.run()
		vm.mode = normalMode
		if vm.err != nil {
			err := vm.err
			vm.err = nil
			vm.restoreLastFrame()

			return nil, err
		}
		return vm.pop(), nil
	case *NativeMethod:
		return m.Function(vm, args)
	case *GetterMethod:
		return m.Call(self)
	case *SetterMethod:
		return m.Call(self, args[1])
	default:
		panic(fmt.Sprintf("tried to call an invalid method: %#v", method))
	}
}

// Call a method without preprocessing its arguments, directly
// on the stack as it is.
func (vm *VM) callMethodOnStack(name value.Symbol, args int) value.Value {
	self := vm.stack[vm.sp-args-1]
	class := self.DirectClass()
	method := class.LookupMethod(name)
	if method == nil {
		return value.NewNoMethodError(name.ToString(), self)
	}

	switch m := method.(type) {
	case *BytecodeMethod:
		vm.createCurrentCallFrame()
		vm.bytecode = m
		vm.fp = vm.sp - args - 1
		vm.ip = 0
	case *NativeMethod:
		result, err := m.Function(vm, vm.stack[vm.sp-args-1:vm.sp])
		if err != nil {
			return err
		}
		vm.popN(args + 1)
		vm.push(result)
	default:
		panic(fmt.Sprintf("tried to call a method that is neither bytecode nor native: %#v", method))
	}

	return nil
}

// The main execution loop of the VM.
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
			vm.returnFromFunction()
			if vm.mode == singleMethodCallMode {
				return
			}
		case bytecode.CONSTANT_CONTAINER:
			vm.constantContainer()
		case bytecode.METHOD_CONTAINER:
			vm.methodContainer()
		case bytecode.SELF:
			vm.self()
		case bytecode.GET_SINGLETON:
			vm.throwIfErr(vm.getSingletonClass())
		case bytecode.DEF_ALIAS:
			vm.throwIfErr(vm.defineAlias())
		case bytecode.DEF_GETTER:
			vm.throwIfErr(vm.defineGetter())
		case bytecode.DEF_CLASS:
			vm.throwIfErr(vm.defineClass())
		case bytecode.DEF_ANON_CLASS:
			vm.throwIfErr(vm.defineAnonymousClass())
		case bytecode.DEF_MODULE:
			vm.throwIfErr(vm.defineModule())
		case bytecode.DEF_ANON_MODULE:
			vm.defineAnonymousModule()
		case bytecode.DEF_MIXIN:
			vm.throwIfErr(vm.defineMixin())
		case bytecode.DEF_ANON_MIXIN:
			vm.defineAnonymousMixin()
		case bytecode.DEF_METHOD:
			vm.throwIfErr(vm.defineMethod())
		case bytecode.INCLUDE:
			vm.throwIfErr(vm.includeMixin())
		case bytecode.DOC_COMMENT:
			vm.throwIfErr(vm.docComment())
		case bytecode.CALL_METHOD8:
			vm.throwIfErr(
				vm.callMethod(int(vm.readByte())),
			)
		case bytecode.CALL_METHOD16:
			vm.throwIfErr(
				vm.callMethod(int(vm.readUint16())),
			)
		case bytecode.CALL_METHOD32:
			vm.throwIfErr(
				vm.callMethod(int(vm.readUint32())),
			)
		case bytecode.CALL_FUNCTION8:
			vm.throwIfErr(
				vm.callFunction(int(vm.readByte())),
			)
		case bytecode.CALL_FUNCTION16:
			vm.throwIfErr(
				vm.callFunction(int(vm.readUint16())),
			)
		case bytecode.CALL_FUNCTION32:
			vm.throwIfErr(
				vm.callFunction(int(vm.readUint32())),
			)
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
			vm.throwIfErr(vm.add())
		case bytecode.SUBTRACT:
			vm.throwIfErr(vm.subtract())
		case bytecode.MULTIPLY:
			vm.throwIfErr(vm.multiply())
		case bytecode.DIVIDE:
			vm.throwIfErr(vm.divide())
		case bytecode.EXPONENTIATE:
			vm.throwIfErr(vm.exponentiate())
		case bytecode.NEGATE:
			vm.throwIfErr(vm.negate())
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
			vm.throwIfErr(vm.getModuleConstant(int(vm.readByte())))
		case bytecode.GET_MOD_CONST16:
			vm.throwIfErr(
				vm.getModuleConstant(int(vm.readUint16())),
			)
		case bytecode.GET_MOD_CONST32:
			vm.throwIfErr(
				vm.getModuleConstant(int(vm.readUint32())),
			)
		case bytecode.DEF_MOD_CONST8:
			vm.throwIfErr(
				vm.defModuleConstant(int(vm.readByte())),
			)
		case bytecode.DEF_MOD_CONST16:
			vm.throwIfErr(
				vm.defModuleConstant(int(vm.readUint16())),
			)
		case bytecode.DEF_MOD_CONST32:
			vm.throwIfErr(
				vm.defModuleConstant(int(vm.readUint32())),
			)
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
		case bytecode.JUMP_UNLESS_UNDEF:
			if vm.peek() != value.Undefined {
				jump := vm.readUint16()
				vm.ip += int(jump)
				break
			}
			vm.ip += 2
		case bytecode.LOOP:
			jump := vm.readUint16()
			vm.ip -= int(jump)
		case bytecode.LBITSHIFT:
			vm.throwIfErr(vm.leftBitshift())
		case bytecode.LOGIC_LBITSHIFT:
			vm.throwIfErr(vm.logicalLeftBitshift())
		case bytecode.RBITSHIFT:
			vm.throwIfErr(vm.rightBitshift())
		case bytecode.LOGIC_RBITSHIFT:
			vm.throwIfErr(vm.logicalRightBitshift())
		case bytecode.BITWISE_AND:
			vm.throwIfErr(vm.bitwiseAnd())
		case bytecode.BITWISE_OR:
			vm.throwIfErr(vm.bitwiseOr())
		case bytecode.BITWISE_XOR:
			vm.throwIfErr(vm.bitwiseXor())
		case bytecode.MODULO:
			vm.throwIfErr(vm.modulo())
		case bytecode.COMPARE:
			vm.throwIfErr(vm.compare())
		case bytecode.EQUAL:
			vm.throwIfErr(vm.equal())
		case bytecode.NOT_EQUAL:
			vm.throwIfErr(vm.notEqual())
		case bytecode.STRICT_EQUAL:
			vm.throwIfErr(vm.strictEqual())
		case bytecode.STRICT_NOT_EQUAL:
			vm.throwIfErr(vm.strictNotEqual())
		case bytecode.GREATER:
			vm.throwIfErr(vm.greaterThan())
		case bytecode.GREATER_EQUAL:
			vm.throwIfErr(vm.greaterThanEqual())
		case bytecode.LESS:
			vm.throwIfErr(vm.lessThan())
		case bytecode.LESS_EQUAL:
			vm.throwIfErr(vm.lessThanEqual())
		default:
			panic(fmt.Sprintf("Unknown bytecode instruction: %#v", instruction))
		}

		// pp.Println(vm.stack[0:vm.sp])
		if vm.err != nil {
			return
		}
	}

}

func (vm *VM) returnFromFunction() {
	returnValue := vm.pop()
	vm.restoreLastFrame()
	vm.push(returnValue)
}

// Restore the state of the VM to the last call frame.
func (vm *VM) restoreLastFrame() {
	lastIndex := len(vm.callFrames) - 1
	cf := vm.callFrames[lastIndex]
	// reset the popped call frame
	vm.callFrames[lastIndex] = CallFrame{}
	vm.callFrames = vm.callFrames[:lastIndex]

	vm.ip = cf.ip
	vm.popN(vm.sp - vm.fp)
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

func (vm *VM) constantContainer() {
	vm.getLocal(1)
}

func (vm *VM) methodContainer() {
	vm.getLocal(2)
}

func (vm *VM) methodContainerValue() value.Value {
	return vm.getLocalValue(2)
}

func (vm *VM) constantContainerValue() value.Value {
	return vm.getLocalValue(1)
}

func (vm *VM) self() {
	vm.getLocal(0)
}

func (vm *VM) getSingletonClass() (err value.Value) {
	val := vm.pop()
	singleton := val.SingletonClass()
	if singleton == nil {
		return value.Errorf(
			value.TypeErrorClass,
			"value `%s` can't have a singleton class",
			val.Inspect(),
		)
	}

	vm.push(singleton)
	return nil
}

func (vm *VM) selfValue() value.Value {
	return vm.getLocalValue(0)
}

// Call a method with an implicit receiver
func (vm *VM) callFunction(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].(*value.CallSiteInfo)

	self := vm.selfValue()
	class := self.DirectClass()

	method := class.LookupMethod(callInfo.Name)
	if method == nil {
		return value.NewNoMethodError(callInfo.Name.ToString(), self)
	}

	// shift all arguments one slot forward to make room for self
	for i := 0; i < callInfo.ArgumentCount; i++ {
		vm.stack[vm.sp-i] = vm.stack[vm.sp-i-1]
	}
	vm.stack[vm.sp-callInfo.ArgumentCount] = self
	vm.sp++

	switch m := method.(type) {
	case *BytecodeMethod:
		return vm.callBytecodeMethod(m, callInfo)
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	}

	panic(fmt.Sprintf("tried to call a method that is neither bytecode nor native: %#v", method))
}

// Call a method with an explicit receiver
func (vm *VM) callMethod(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].(*value.CallSiteInfo)

	self := vm.stack[vm.sp-callInfo.ArgumentCount-1]
	class := self.DirectClass()

	method := class.LookupMethod(callInfo.Name)
	if method == nil {
		return value.NewNoMethodError(callInfo.Name.ToString(), self)
	}
	switch m := method.(type) {
	case *BytecodeMethod:
		return vm.callBytecodeMethod(m, callInfo)
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	case *GetterMethod:
		if callInfo.ArgumentCount != 0 {
			return value.NewWrongArgumentCountError(callInfo.ArgumentCount, 0)
		}
		result, err := m.Call(self)
		if err != nil {
			return err
		}
		vm.push(result)
		return nil
	case *SetterMethod:
		if callInfo.ArgumentCount != 1 {
			return value.NewWrongArgumentCountError(callInfo.ArgumentCount, 1)
		}
		other := vm.stack[vm.sp-1]
		result, err := m.Call(self, other)
		if err != nil {
			return err
		}
		vm.push(result)
		return nil
	default:
		panic(fmt.Sprintf("tried to call an invalid method: %#v", method))
	}
}

// set up the vm to execute a bytecode method
func (vm *VM) callNativeMethod(method *NativeMethod, callInfo *value.CallSiteInfo) (err value.Value) {
	if err := vm.prepareArguments(method, callInfo); err != nil {
		return err
	}

	returnVal, err := method.Function(vm, vm.stack[vm.sp-method.ParameterCount()-1:vm.sp])
	if err != nil {
		return err
	}
	vm.popN(method.ParameterCount() + 1)
	vm.push(returnVal)
	return nil
}

// set up the vm to execute a bytecode method
func (vm *VM) callBytecodeMethod(method *BytecodeMethod, callInfo *value.CallSiteInfo) (err value.Value) {
	if err := vm.prepareArguments(method, callInfo); err != nil {
		return err
	}

	vm.createCurrentCallFrame()

	vm.bytecode = method
	vm.fp = vm.sp - method.ParameterCount() - 1
	vm.ip = 0

	return nil
}

func (vm *VM) prepareArguments(method value.Method, callInfo *value.CallSiteInfo) (err value.Value) {
	namedArgCount := callInfo.NamedArgumentCount()

	if namedArgCount == 0 {
		if err := vm.preparePositionalArguments(method, callInfo); err != nil {
			return err
		}
	} else if err := vm.prepareNamedArguments(method, callInfo); err != nil {
		return err
	}

	return nil
}

func (vm *VM) prepareNamedArguments(method value.Method, callInfo *value.CallSiteInfo) (err value.Value) {
	paramCount := method.ParameterCount()
	namedArgCount := callInfo.NamedArgumentCount()
	reqParamCount := paramCount - method.OptionalParameterCount()
	posArgCount := callInfo.PositionalArgumentCount()

	// create a slice containing the given arguments
	// in original order
	namedArgs := make([]value.Value, namedArgCount)
	copy(namedArgs, vm.stack[vm.sp-namedArgCount:vm.sp])

	var posParamNames []value.Symbol
	var namedParamNames []value.Symbol
	var spIncrease int
	paramNames := method.Parameters()

	if method.PostRestParameterCount() >= 0 {
		requiredPosParamCount := paramCount - method.OptionalParameterCount() - method.PostRestParameterCount() - 1
		if posArgCount < requiredPosParamCount {
			return value.NewWrongPositionalArgumentCountError(
				posArgCount,
				requiredPosParamCount,
			)
		}

		firstPosRestArg := paramCount - method.PostRestParameterCount() - 1
		lastPosRestArg := callInfo.ArgumentCount - method.PostRestParameterCount() - 1
		posRestArgCount := lastPosRestArg - firstPosRestArg + 1
		postArgCount := callInfo.ArgumentCount - lastPosRestArg - 1
		var postArgs []value.Value
		var restList value.List
		if postArgCount > 0 {
			postArgs = make([]value.Value, postArgCount)
			copy(postArgs, vm.stack[vm.sp-postArgCount:vm.sp])
			vm.popN(postArgCount)
		}

		if posRestArgCount > 0 {
			restList = make(value.List, posRestArgCount)
		}
		for i := 1; i <= posRestArgCount; i++ {
			restList[posRestArgCount-i] = vm.pop()
		}
		vm.push(restList)
		for _, postArg := range postArgs {
			vm.push(postArg)
		}

		posParamNames = paramNames[:firstPosRestArg+1]
		namedParamNames = paramNames[paramCount-(callInfo.ArgumentCount-posArgCount):]
		spIncrease = paramCount - (callInfo.ArgumentCount - posRestArgCount + 1)
	} else {
		posParamNames = paramNames[:posArgCount]
		namedParamNames = paramNames[posArgCount:]
		spIncrease = paramCount - callInfo.ArgumentCount
	}

	var foundNamedArgCount int

	for _, paramName := range posParamNames {
		for j := 0; j < namedArgCount; j++ {
			if paramName == callInfo.NamedArguments[j] {
				return value.NewDuplicatedArgumentError(
					method.Name().ToString(),
					paramName.ToString(),
				)
			}
		}
	}

methodParamLoop:
	for i, paramName := range namedParamNames {
		found := false
		targetIndex := vm.sp - namedArgCount + i
	namedArgLoop:
		for j := 0; j < namedArgCount; j++ {
			if paramName != callInfo.NamedArguments[j] {
				continue namedArgLoop
			}

			found = true
			foundNamedArgCount++
			if i == j {
				break namedArgLoop
			}

			vm.stack[targetIndex] = namedArgs[j]
			// mark the found value as undefined
			namedArgs[j] = value.Undefined
		}

		if found {
			continue methodParamLoop
		}

		// the parameter is required
		// but is not present in the call
		if posArgCount+i < reqParamCount {
			return value.NewRequiredArgumentMissingError(
				method.Name().ToString(),
				paramName.ToString(),
			)
		}

		vm.stack[targetIndex] = value.Undefined
	}

	unknownNamedArgCount := namedArgCount - foundNamedArgCount
	if unknownNamedArgCount != 0 {
		// construct a slice that contains
		// the names of unknown named arguments
		// that have been given
		unknownNamedArgNames := make([]value.Symbol, unknownNamedArgCount)
		for i, namedArg := range namedArgs {
			if namedArg == value.Undefined {
				continue
			}

			unknownNamedArgNames[i] = callInfo.NamedArguments[i]
		}

		return value.NewUnknownArgumentsError(unknownNamedArgNames)
	}

	vm.sp += spIncrease
	return nil
}

func (vm *VM) preparePositionalArguments(method value.Method, callInfo *value.CallSiteInfo) (err value.Value) {
	optParamCount := method.OptionalParameterCount()
	postParamCount := method.PostRestParameterCount()
	paramCount := method.ParameterCount()
	preRestParamCount := paramCount - postParamCount - 1
	reqParamCount := paramCount - optParamCount
	if postParamCount >= 0 {
		reqParamCount -= 1
	}

	if callInfo.ArgumentCount < reqParamCount {
		if postParamCount == -1 {
			return value.NewWrongArgumentCountRangeError(
				callInfo.ArgumentCount,
				reqParamCount,
				paramCount,
			)
		} else {
			return value.NewWrongArgumentCountRestError(
				callInfo.ArgumentCount,
				reqParamCount,
			)
		}
	}

	if optParamCount > 0 {
		// populate missing optional arguments with undefined
		missingArgCount := preRestParamCount - callInfo.ArgumentCount
		for i := 0; i < missingArgCount; i++ {
			vm.push(value.Undefined)
		}
	} else if postParamCount == -1 && paramCount != callInfo.ArgumentCount {
		return value.NewWrongArgumentCountError(
			callInfo.ArgumentCount,
			paramCount,
		)
	}

	if postParamCount >= 0 {
		var postArgs []value.Value
		if postParamCount > 0 {
			postArgs = make([]value.Value, postParamCount)
			copy(postArgs, vm.stack[vm.sp-postParamCount:vm.sp])
			vm.popN(postParamCount)
		}
		var restList value.List
		restArgCount := callInfo.ArgumentCount - preRestParamCount - postParamCount
		if restArgCount > 0 {
			// rest arguments
			restList = make(value.List, restArgCount)
			for i := 1; i <= restArgCount; i++ {
				restList[restArgCount-i] = vm.pop()
			}
		}
		vm.push(restList)
		for _, postArg := range postArgs {
			vm.push(postArg)
		}
	}

	return nil
}

// Include a mixin in a class/mixin.
func (vm *VM) includeMixin() (err value.Value) {
	targetValue := vm.pop()
	mixinVal := vm.pop()

	mixin, ok := mixinVal.(*value.Mixin)
	if !ok {
		return value.NewIsNotMixinError(mixinVal.Inspect())
	}

	switch target := targetValue.(type) {
	case *value.Class:
		target.IncludeMixin(mixin)
	case *value.Mixin:
		target.IncludeMixin(mixin)
	default:
		return value.Errorf(
			value.TypeErrorClass,
			"can't include into an instance of %s: `%s`",
			targetValue.Class().PrintableName(),
			target.Inspect(),
		)
	}

	return nil
}

// Attach a doc comment to an object.
func (vm *VM) docComment() (err value.Value) {
	// targetValue := vm.pop()
	// docStringVal := vm.pop()

	// docString := docStringVal.(*value.String)

	// switch target := targetValue.(type) {
	// case *value.Class:
	// 	target.IncludeMixin(mixin)
	// case *value.Mixin:
	// 	target.IncludeMixin(mixin)
	// default:
	// 	return value.Errorf(
	// 		value.TypeErrorClass,
	// 		"can't include into an instance of %s: `%s`",
	// 		targetValue.Class().PrintableName(),
	// 		target.Inspect(),
	// 	)
	// }

	return nil
}

// Define a new method
func (vm *VM) defineMethod() value.Value {
	nameVal := vm.pop()
	bodyVal := vm.pop()

	body := bodyVal.(*BytecodeMethod)
	name := nameVal.(value.Symbol)

	methodContainer := vm.methodContainerValue()

	switch m := methodContainer.(type) {
	case *value.Class:
		if !m.CanOverride(name) {
			return value.NewCantOverrideAFrozenMethod(name.ToString())
		}
		m.Methods[name] = body
	case *value.Mixin:
		if !m.CanOverride(name) {
			return value.NewCantOverrideAFrozenMethod(name.ToString())
		}
		m.Methods[name] = body
	default:
		panic(fmt.Sprintf("invalid method container: %s", methodContainer.Inspect()))
	}

	vm.push(body)
	return nil
}

// Define a new anonymous mixin
func (vm *VM) defineAnonymousMixin() {
	bodyVal := vm.pop()

	mixin := value.NewMixin()

	switch body := bodyVal.(type) {
	case *BytecodeMethod:
		vm.executeMixinBody(mixin, body)
	case value.UndefinedType:
		vm.push(mixin)
	default:
		panic(fmt.Sprintf("expected undefined or a bytecode function as the mixin body, got: %s", bodyVal.Inspect()))
	}
}

// Define a new mixin
func (vm *VM) defineMixin() (err value.Value) {
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
	case *value.Mixin:
		parentModule = &m.ModulelikeObject
	default:
		return value.NewIsNotModuleError(parentModuleVal.Inspect())
	}

	var mixin *value.Mixin

	if mixinVal, ok := parentModule.Constants.Get(constantName); ok {
		mixin, ok = mixinVal.(*value.Mixin)
		if !ok {
			return value.NewRedefinedConstantError(parentModuleVal.Inspect(), constantName.Inspect())
		}
	} else {
		mixin = value.NewMixin()
		parentModule.AddConstant(constantName, mixin)
	}

	switch body := bodyVal.(type) {
	case *BytecodeMethod:
		vm.executeMixinBody(mixin, body)
	case value.UndefinedType:
		vm.push(mixin)
	default:
		panic(fmt.Sprintf("expected undefined or a bytecode function as the mixin body, got: %s", bodyVal.Inspect()))
	}

	return nil
}

// Define a new anonymous module
func (vm *VM) defineAnonymousModule() {
	bodyVal := vm.pop()

	module := value.NewModule()

	switch body := bodyVal.(type) {
	case *BytecodeMethod:
		vm.executeModuleBody(module, body)
	case value.UndefinedType:
		vm.push(module)
	default:
		panic(fmt.Sprintf("expected undefined or a bytecode function as the module body, got: %s", bodyVal.Inspect()))
	}
}

// Define a new module
func (vm *VM) defineModule() (err value.Value) {
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
	case *value.Mixin:
		parentModule = &m.ModulelikeObject
	default:
		return value.NewIsNotModuleError(parentModuleVal.Inspect())
	}

	var module *value.Module

	if moduleVal, ok := parentModule.Constants.Get(constantName); ok {
		module, ok = moduleVal.(*value.Module)
		if !ok {
			return value.NewRedefinedConstantError(parentModuleVal.Inspect(), constantName.Inspect())
		}
	} else {
		module = value.NewModule()
		parentModule.AddConstant(constantName, module)
	}

	switch body := bodyVal.(type) {
	case *BytecodeMethod:
		vm.executeModuleBody(module, body)
	case value.UndefinedType:
		vm.push(module)
	default:
		panic(fmt.Sprintf("expected undefined or a bytecode function as the module body, got: %s", bodyVal.Inspect()))
	}

	return nil
}

// Define a new anonymous class
func (vm *VM) defineAnonymousClass() (err value.Value) {
	superclassVal := vm.pop()
	bodyVal := vm.pop()

	class := value.NewClass()
	switch superclass := superclassVal.(type) {
	case *value.Class:
		class.Parent = superclass
	case value.UndefinedType:
	default:
		return value.Errorf(
			value.TypeErrorClass,
			"`%s` can't be used as a superclass", superclass.Inspect(),
		)
	}

	switch body := bodyVal.(type) {
	case *BytecodeMethod:
		vm.executeClassBody(class, body)
	case value.UndefinedType:
		vm.push(class)
	default:
		panic(fmt.Sprintf("expected undefined or a bytecode function as the class body, got: %s", bodyVal.Inspect()))
	}

	return nil
}

// Define a getter method
func (vm *VM) defineGetter() value.Value {
	name := vm.pop().(value.Symbol)

	var container *value.MethodContainer
	methodContainerValue := vm.methodContainerValue()

	switch methodContainer := methodContainerValue.(type) {
	case *value.Class:
		container = &methodContainer.MethodContainer
	case *value.Mixin:
		container = &methodContainer.MethodContainer
	}

	err := DefineGetter(container, name, false)
	if err != nil {
		return err
	}

	return nil
}

// Define a method alias
func (vm *VM) defineAlias() value.Value {
	newName := vm.pop().(value.Symbol)
	oldName := vm.pop().(value.Symbol)

	var err *value.Error
	methodContainerValue := vm.methodContainerValue()

	switch methodContainer := methodContainerValue.(type) {
	case *value.Class:
		err = methodContainer.DefineAlias(newName, oldName)
		if err != nil {
			return err
		}
	case *value.Mixin:
		err = methodContainer.DefineAlias(newName, oldName)
		if err != nil {
			return err
		}
	}

	return nil
}

// Define a new class
func (vm *VM) defineClass() (err value.Value) {
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
	case *value.Mixin:
		parentModule = &mod.ModulelikeObject
	default:
		return value.NewIsNotModuleError(parentModuleVal.Inspect())
	}

	var class *value.Class

	if classVal, ok := parentModule.Constants.Get(constantName); ok {
		class, ok = classVal.(*value.Class)
		if !ok {
			return value.NewRedefinedConstantError(parentModuleVal.Inspect(), constantName.Inspect())
		}
		switch superclass := superclassVal.(type) {
		case *value.Class:
			if class.Parent != superclass {
				return value.Errorf(
					value.TypeErrorClass,
					"superclass mismatch in %s, expected: %s, got: %s",
					class.Name,
					class.Parent.Name,
					superclass.Name,
				)
			}
		case value.UndefinedType:
		default:
			return value.Errorf(
				value.TypeErrorClass,
				"`%s` can't be used as a superclass", superclass.Inspect(),
			)
		}
	} else {
		class = value.NewClass()
		switch superclass := superclassVal.(type) {
		case *value.Class:
			class.Parent = superclass
		case value.UndefinedType:
		default:
			return value.Errorf(
				value.TypeErrorClass,
				"`%s` can't be used as a superclass", superclass.Inspect(),
			)
		}
		parentModule.AddConstant(constantName, class)
	}

	switch body := bodyVal.(type) {
	case *BytecodeMethod:
		vm.executeClassBody(class, body)
	case value.UndefinedType:
		vm.push(class)
	default:
		panic(fmt.Sprintf("expected undefined or a bytecode function as the class body, got: %s", bodyVal.Inspect()))
	}

	return nil
}

func (vm *VM) addCallFrame(cf CallFrame) {
	if len(vm.callFrames) == CALL_STACK_SIZE {
		panic(fmt.Sprintf("Stack overflow: %d", CALL_STACK_SIZE))
	}

	vm.callFrames = append(vm.callFrames, cf)
}

// preserve the current state of the vm in a call frame
func (vm *VM) createCurrentCallFrame() {
	vm.addCallFrame(
		CallFrame{
			bytecode: vm.bytecode,
			ip:       vm.ip,
			fp:       vm.fp,
		},
	)
}

// set up the vm to execute a class body
func (vm *VM) executeClassBody(class value.Value, body *BytecodeMethod) {
	vm.createCurrentCallFrame()

	vm.bytecode = body
	vm.fp = vm.sp
	vm.ip = 0
	// set class as `self`
	vm.push(class)
	// set class as constant container
	vm.push(class)
	// set class as method container
	vm.push(class)
}

// set up the vm to execute a mixin body
func (vm *VM) executeMixinBody(mixin value.Value, body *BytecodeMethod) {
	vm.createCurrentCallFrame()

	vm.bytecode = body
	vm.fp = vm.sp
	vm.ip = 0
	// set mixin as `self`
	vm.push(mixin)
	// set mixin as constant container
	vm.push(mixin)
	// set mixin as method container
	vm.push(mixin)
}

// set up the vm to execute a module body
func (vm *VM) executeModuleBody(module value.Value, body *BytecodeMethod) {
	vm.createCurrentCallFrame()

	vm.bytecode = body
	vm.fp = vm.sp
	vm.ip = 0
	// set module as `self`
	vm.push(module)
	// set module as constant container
	vm.push(module)
	// set module's singleton class as method container
	vm.push(module.SingletonClass())
}

// Set a local variable or value.
func (vm *VM) setLocal(index int) {
	vm.setLocalValue(index, vm.peek())
}

// Set a local variable or value.
func (vm *VM) setLocalValue(index int, val value.Value) {
	vm.stack[vm.fp+index] = val
}

// Read a local variable or value.
func (vm *VM) getLocal(index int) {
	vm.push(vm.getLocalValue(index))
}

// Read a local variable or value.
func (vm *VM) getLocalValue(index int) value.Value {
	return vm.stack[vm.fp+index]
}

// Pop a module off the stack and look for a constant with the given name.
func (vm *VM) getModuleConstant(nameIndex int) (err value.Value) {
	symbol := vm.bytecode.Values[nameIndex].(value.Symbol)
	mod := vm.pop()
	var constants value.SymbolMap

	switch m := mod.(type) {
	case *value.Class:
		constants = m.Constants
	case *value.Module:
		constants = m.Constants
	default:
		return value.Errorf(value.TypeErrorClass, "`%s` is not a module", mod.Inspect())
	}

	val, ok := constants.Get(symbol)
	if !ok {
		return value.Errorf(value.NoConstantErrorClass, "%s doesn't have a constant named `%s`", mod.Inspect(), symbol.Inspect())
	}

	vm.push(val)
	return nil
}

// Pop two values off the stack and define a constant with the given name.
func (vm *VM) defModuleConstant(nameIndex int) (err value.Value) {
	symbol := vm.bytecode.Values[nameIndex].(value.Symbol)
	mod := vm.pop()
	var constants value.SymbolMap

	switch m := mod.(type) {
	case *value.Class:
		constants = m.Constants
	case *value.Module:
		constants = m.Constants
	case *value.Mixin:
		constants = m.Constants
	default:
		return value.NewIsNotModuleError(mod.Inspect())
	}

	val := vm.peek()
	if constants.Has(symbol) {
		return value.NewRedefinedConstantError(mod.Inspect(), symbol.Inspect())
	}
	constants.Set(symbol, val)
	return nil
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

// Return the nth element on top of the stack
// without popping it.
func (vm *VM) peekAt(n int) value.Value {
	if vm.sp-n <= 0 {
		panic(fmt.Sprintf("tried to peek outside of the valid range: %d", n))
	}

	return vm.stack[vm.sp-1-n]
}

// Negate the element on top of the stack
func (vm *VM) negate() (err value.Value) {
	operand := vm.peek()
	result, builtin := value.Negate(operand)
	if !builtin {
		return value.NewNoMethodError("-", operand)
	}

	vm.replace(result)
	return nil
}

type binaryOperationWithoutErrFunc func(left value.Value, right value.Value) (value.Value, bool)

func (vm *VM) binaryOperationWithoutErr(fn binaryOperationWithoutErrFunc, methodName value.Symbol) (err value.Value) {
	right := vm.peek()
	left := vm.peekAt(1)

	result, builtin := fn(left, right)
	if builtin {
		vm.pop()
		vm.replace(result)
		return nil
	}

	er := vm.callMethodOnStack(methodName, 1)
	if er != nil {
		return er
	}

	return nil
}

func (vm *VM) negatedBinaryOperationWithoutErr(fn binaryOperationWithoutErrFunc, methodName value.Symbol) (err value.Value) {
	right := vm.peek()
	left := vm.peekAt(1)

	result, builtin := fn(left, right)
	if builtin {
		vm.pop()
		vm.replace(result)
		return nil
	}

	er := vm.callMethodOnStack(methodName, 1)
	if er != nil {
		return er
	}
	vm.replace(value.ToNotBool(vm.peek()))

	return nil
}

type binaryOperationFunc func(left value.Value, right value.Value) (value.Value, *value.Error, bool)

func (vm *VM) binaryOperation(fn binaryOperationFunc, methodName value.Symbol) value.Value {
	right := vm.peek()
	left := vm.peekAt(1)

	result, err, builtin := fn(left, right)
	if err != nil {
		return err
	}
	if builtin {
		vm.pop()
		vm.replace(result)
		return nil
	}

	er := vm.callMethodOnStack(methodName, 1)
	if er != nil {
		return er
	}
	return nil
}

var (
	andSymbol                  value.Symbol = value.ToSymbol("&")
	orSymbol                   value.Symbol = value.ToSymbol("|")
	xorSymbol                  value.Symbol = value.ToSymbol("^")
	spaceshipSymbol            value.Symbol = value.ToSymbol("<=>")
	percentSymbol              value.Symbol = value.ToSymbol("%")
	equalSymbol                value.Symbol = value.ToSymbol("==")
	strictEqualSymbol          value.Symbol = value.ToSymbol("===")
	greaterThanSymbol          value.Symbol = value.ToSymbol(">")
	greaterThanEqualSymbol     value.Symbol = value.ToSymbol(">=")
	lessThanSymbol             value.Symbol = value.ToSymbol("<")
	lessThanEqualSymbol        value.Symbol = value.ToSymbol("<=")
	leftBitshiftSymbol         value.Symbol = value.ToSymbol("<<")
	logicalLeftBitshiftSymbol  value.Symbol = value.ToSymbol("<<<")
	rightBitshiftSymbol        value.Symbol = value.ToSymbol(">>")
	logicalRightBitshiftSymbol value.Symbol = value.ToSymbol(">>>")
	addSymbol                  value.Symbol = value.ToSymbol("+")
	subtractSymbol             value.Symbol = value.ToSymbol("-")
	multiplySymbol             value.Symbol = value.ToSymbol("*")
	divideSymbol               value.Symbol = value.ToSymbol("/")
	exponentiateSymbol         value.Symbol = value.ToSymbol("**")
)

// Perform a bitwise AND and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) bitwiseAnd() (err value.Value) {
	return vm.binaryOperation(value.BitwiseAnd, andSymbol)
}

// Perform a bitwise OR and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) bitwiseOr() (err value.Value) {
	return vm.binaryOperation(value.BitwiseOr, orSymbol)
}

// Perform a bitwise XOR and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) bitwiseXor() (err value.Value) {
	return vm.binaryOperation(value.BitwiseXor, xorSymbol)
}

// Perform a comparison and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) compare() (err value.Value) {
	return vm.binaryOperation(value.Compare, spaceshipSymbol)
}

// Perform modulo and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) modulo() (err value.Value) {
	return vm.binaryOperation(value.Modulo, percentSymbol)
}

// Check whether two top elements on the stack are equal and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) equal() (err value.Value) {
	return vm.binaryOperationWithoutErr(value.Equal, equalSymbol)
}

// Check whether two top elements on the stack are not and equal push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) notEqual() (err value.Value) {
	return vm.negatedBinaryOperationWithoutErr(value.NotEqual, equalSymbol)
}

// Check whether two top elements on the stack are strictly equal push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) strictEqual() (err value.Value) {
	return vm.binaryOperationWithoutErr(value.StrictEqual, strictEqualSymbol)
}

// Check whether two top elements on the stack are strictly not equal push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) strictNotEqual() (err value.Value) {
	return vm.negatedBinaryOperationWithoutErr(value.StrictNotEqual, strictEqualSymbol)
}

// Check whether the first operand is greater than the second and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) greaterThan() (err value.Value) {
	return vm.binaryOperation(value.GreaterThan, greaterThanSymbol)
}

// Check whether the first operand is greater than or equal to the second and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) greaterThanEqual() (err value.Value) {
	return vm.binaryOperation(value.GreaterThanEqual, greaterThanEqualSymbol)
}

// Check whether the first operand is less than the second and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) lessThan() (err value.Value) {
	return vm.binaryOperation(value.LessThan, lessThanSymbol)
}

// Check whether the first operand is less than or equal to the second and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) lessThanEqual() (err value.Value) {
	return vm.binaryOperation(value.LessThanEqual, lessThanEqualSymbol)
}

// Perform a left bitshift and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) leftBitshift() (err value.Value) {
	return vm.binaryOperation(value.LeftBitshift, leftBitshiftSymbol)
}

// Perform a logical left bitshift and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) logicalLeftBitshift() (err value.Value) {
	return vm.binaryOperation(value.LogicalLeftBitshift, logicalLeftBitshiftSymbol)
}

// Perform a right bitshift and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) rightBitshift() (err value.Value) {
	return vm.binaryOperation(value.RightBitshift, rightBitshiftSymbol)
}

// Perform a logical right bitshift and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) logicalRightBitshift() (err value.Value) {
	return vm.binaryOperation(value.LogicalRightBitshift, logicalRightBitshiftSymbol)
}

// Add two operands together and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) add() (err value.Value) {
	return vm.binaryOperation(value.Add, addSymbol)
}

// Subtract two operands and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) subtract() (err value.Value) {
	return vm.binaryOperation(value.Subtract, subtractSymbol)
}

// Multiply two operands together and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) multiply() (err value.Value) {
	return vm.binaryOperation(value.Multiply, multiplySymbol)
}

// Divide two operands and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) divide() (err value.Value) {
	return vm.binaryOperation(value.Divide, divideSymbol)
}

// Exponentiate two operands and push the result to the stack.
// Returns false when an error has been raised.
func (vm *VM) exponentiate() (err value.Value) {
	return vm.binaryOperation(value.Exponentiate, exponentiateSymbol)
}

// Throw an error and attempt to find code
// that catches it.
func (vm *VM) throw(err value.Value) {
	vm.err = err
}
