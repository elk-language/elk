//go:build native

// Package vm contains the Elk Virtual Machine.
// It interprets Elk Bytecode produced by
// the Elk compiler.
package vm

import (
	"fmt"
	"io"
	"os"
	"slices"
	"unsafe"

	"github.com/elk-language/elk/value"
)

// A single instance of the Elk Virtual Machine.
type VM struct {
	ID     int64
	Stdin  io.Reader // standard output used by the VM
	Stdout io.Writer // standard input used by the VM
	Stderr io.Writer // standard error used by the VM

	cfp           uintptr           // Call frame pointer
	callFrames    []value.CallFrame // Call stack
	errStackTrace *value.StackTrace // The most recent error stack trace
	threadPool    *ThreadPool
	state         state
}

// Create a new VM instance.
func New(opts ...Option) *VM {
	callFrames := make([]value.CallFrame, CALL_STACK_SIZE)
	// mark the end of the call stack with a sentinel value
	callFrames[len(callFrames)-1] = makeSentinelCallFrame()

	id := currentID.Add(1)

	vm := &VM{
		ID:         id,
		stack:      stack,
		callFrames: callFrames,
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		threadPool: DefaultThreadPool,
	}
	vm.cfpSet(&callFrames[0])

	for _, opt := range opts {
		opt(vm)
	}

	return vm
}

func (vm *VM) InterpretTopLevel(fn *BytecodeFunction) (value.Value, value.Value) {
	panic("cannot interpret bytecode in native mode")
}

// Execute the given bytecode chunk.
func (vm *VM) InterpretREPL(fn *BytecodeFunction) (value.Value, value.Value) {
	panic("cannot interpret bytecode in native mode")
}

// Get the stored error.
func (vm *VM) Err() value.Value {
	panic("cannot get stored error in native mode")
}

// Get the value on top of the value stack.
func (vm *VM) StackTop() value.Value {
	panic("cannot get stack top in native mode")
}

func (vm *VM) ValueStack() []value.Value {
	panic("cannot get value stack in native mode")
}

func (vm *VM) PrintErrorValue(err value.Value) {
	PrintError(vm.Stderr, vm.ErrStackTrace(), err)
}

func (vm *VM) PrintError() {
	panic("cannot print stored error")
}

func (vm *VM) callStack() []CallFrame {
	cfIndex := vm.cfpOffset()
	return vm.callFrames[:cfIndex]
}

func (vm *VM) InspectValueStack() {
	fmt.Println("value stack: none, native mode")
}

func (vm *VM) InspectCallStack() {
	fmt.Println("call stack:")
	for i, frame := range vm.callStack() {
		fmt.Printf("%d => %#v\n", i, frame)
	}
}

func (vm *VM) CallGeneratorNext(generator *Generator) (value.Value, value.Value) {
	// TODO: implement native generators
	vm.createCurrentCallFrame(true)
}

// Call a callable value from Go code, preserving the state of the VM.
func (vm *VM) CallCallable(args ...value.Value) (value.Value, value.Value) {
	// TODO: implement native closures
	function := args[0]
	switch f := function.SafeAsReference().(type) {
	case *Closure:
		return vm.CallClosure(f, args[1:]...)
	default:
		return vm.CallMethodByName(callSymbol, args...)
	}
}

// Call an Elk closure from Go code, preserving the state of the VM.
func (vm *VM) CallClosure(closure *Closure, args ...value.Value) (value.Value, value.Value) {
	// TODO: implement native closures
}

// Call an Elk method from Go code, preserving the state of the VM.
func (vm *VM) CallMethodByName(name value.Symbol, args ...value.Value) (value.Value, value.Value) {
	self := args[0]
	class := self.DirectClass()
	method := class.LookupMethod(name)
	return vm.CallMethod(method, args...)
}

// Call an Elk method from Go code, preserving the state of the VM.
func (vm *VM) CallMethodByNameWithCache(name value.Symbol, cc **value.CallCache, args ...value.Value) (value.Value, value.Value) {
	self := args[0]
	class := self.DirectClass()
	method := value.LookupMethodInCache(class, name, cc)
	return vm.CallMethod(method, args...)
}

func (vm *VM) populateMissingParameters(args []value.Value, paramCount, argumentCount int) []value.Value {
	// populate missing optional arguments with undefined
	missingParams := uintptr(paramCount - argumentCount)
	if missingParams > 0 {
		newArgs := make([]value.Value, paramCount)
		clone(newArgs, args)
		return newArgs
	}

	return args
}

func (vm *VM) CallMethod(method value.Method, args ...value.Value) (value.Value, value.Value) {
	self := args[0]
	paramCount := method.ParameterCount()
	argCount := len(args) - 1

	switch m := method.(type) {
	case *NativeMethod:
		return m.Function(
			vm,
			vm.populateMissingParameters(args, paramCount, argCount),
		)
	case *GetterMethod:
		return m.Call(self)
	case *SetterMethod:
		return m.Call(self, args[1])
	default:
		panic(fmt.Sprintf("tried to call an invalid method: %#v", method))
	}
}

func (vm *VM) ClearStackFrames() {
	vm.cfpSet(&vm.callFrames[0])
}

func (vm *VM) ResetError() {
	vm.state = idleState
	vm.errStackTrace = nil
}

func (vm *VM) GetStackTrace() *value.StackTrace {
	if vm.errStackTrace != nil {
		return vm.errStackTrace
	}

	return vm.BuildStackTrace()
}

func (vm *VM) AddCallFrame(cf value.CallFrame) {
	*vm.cfpGet() = cf
	vm.cfpIncrement()
}

func (vm *VM) PopCallFrame() {
	vm.cfpDecrementBy(1)
}

func (vm *VM) BuildStackTrace() *value.StackTrace {
	callStack := vm.callStack()
	return (*value.StackTrace)(slices.Clone(callStack))
}

func (vm *VM) BuildStackTracePrepend(base *value.StackTrace) *value.StackTrace {
	callStack := vm.callStack()

	stackTraceSlice := make([]value.CallFrame, 0, len(*base)+len(callStack)+1)
	stackTraceSlice = append(stackTraceSlice, callStack...)
	stackTraceSlice = append(stackTraceSlice, (*base)...)

	return (*value.StackTrace)(&stackTraceSlice)
}

// Increment the call frame pointer
func (vm *VM) cfpIncrement() {
	vm.cfpIncrementBy(1)
}

// Add n to the call frame pointer
func (vm *VM) cfpIncrementBy(n int) {
	ptr := vm.cfpAdd(n)
	if ptr.sentinel {
		panic("call stack overflow")
	}
	vm.cfpSet(ptr)
}

func (vm *VM) lastCallFrame() *CallFrame {
	return vm.cfpAdd(-1)
}

func (vm *VM) cfpAdd(n int) *CallFrame {
	return vm.callFrameAdd(vm.cfpGet(), n)
}

func (vm *VM) cfpAddRaw(n uintptr) uintptr {
	return vm.cfp + n*CallFrameSize
}

func (vm *VM) cfpSubtractRaw(n uintptr) uintptr {
	return vm.cfp - n*CallFrameSize
}

func (vm *VM) cfpGet() *CallFrame {
	return (*CallFrame)(unsafe.Pointer(vm.cfp))
}

// Set the typesafe call frame pointer
func (vm *VM) cfpSet(ptr *CallFrame) {
	vm.cfp = uintptr(unsafe.Pointer(ptr))
}

func (vm *VM) callFrameAdd(ptr *CallFrame, n int) *CallFrame {
	return (*CallFrame)(unsafe.Add(unsafe.Pointer(ptr), n*int(CallFrameSize)))
}

func (vm *VM) callFrameAddRaw(ptr uintptr, n uintptr) uintptr {
	return ptr + n*CallFrameSize
}

func (vm *VM) cfpOffset() int {
	return int(vm.cfp-uintptr(unsafe.Pointer(&vm.callFrames[0]))) / int(CallFrameSize)
}
