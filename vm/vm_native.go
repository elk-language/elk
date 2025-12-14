//go:build native

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
type Thread struct {
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
func New(opts ...Option) *Thread {
	callFrames := make([]value.CallFrame, CALL_STACK_SIZE)
	// mark the end of the call stack with a sentinel value
	callFrames[len(callFrames)-1] = value.MakeSentinelCallFrame()

	id := currentID.Add(1)

	vm := &Thread{
		ID:         id,
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

func (vm *Thread) InterpretTopLevel(fn *BytecodeFunction) (value.Value, value.Value) {
	panic("cannot interpret bytecode in native mode")
}

// Execute the given bytecode chunk.
func (vm *Thread) InterpretREPL(fn *BytecodeFunction) (value.Value, value.Value) {
	panic("cannot interpret bytecode in native mode")
}

// Get the stored error.
func (vm *Thread) Err() value.Value {
	panic("cannot get stored error in native mode")
}

// Get the value on top of the value stack.
func (vm *Thread) StackTop() value.Value {
	panic("cannot get stack top in native mode")
}

func (vm *Thread) ValueStack() []value.Value {
	panic("cannot get value stack in native mode")
}

func (vm *Thread) PrintError() {
	panic("cannot print stored error")
}

func (vm *Thread) callStack() []value.CallFrame {
	cfIndex := vm.cfpOffset()
	return vm.callFrames[:cfIndex]
}

func (vm *Thread) InspectValueStack() {
	fmt.Println("value stack: none, native mode")
}

func (vm *Thread) InspectCallStack() {
	fmt.Println("call stack:")
	for i, frame := range vm.callStack() {
		fmt.Printf("%d => %#v\n", i, frame)
	}
}

func (vm *Thread) CallGeneratorNext(generator *Generator) (value.Value, value.Value) {
	// TODO: implement native generators
	// vm.createCurrentCallFrame(true)
	panic("native generators are not yet implemented")
}

// Call a callable value from Go code, preserving the state of the VM.
func (vm *Thread) CallCallable(args ...value.Value) (value.Value, value.Value) {
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
func (vm *Thread) CallClosure(closure *Closure, args ...value.Value) (value.Value, value.Value) {
	// TODO: implement native closures
	panic("native closures are not yet implemented")
}

// Call an Elk method from Go code, preserving the state of the VM.
func (vm *Thread) CallMethodByName(name value.Symbol, args ...value.Value) (value.Value, value.Value) {
	self := args[0]
	class := self.DirectClass()
	method := class.LookupMethod(name)
	return vm.CallMethod(method, args...)
}

// Call an Elk method from Go code, preserving the state of the VM.
func (vm *Thread) CallMethodByNameWithCache(name value.Symbol, cc **value.CallCache, args ...value.Value) (value.Value, value.Value) {
	self := args[0]
	class := self.DirectClass()
	method := value.LookupMethodInCache(class, name, cc)
	if method == nil {
		panic(fmt.Sprintf("no such method, class: %s, name: %s", class.Name, name))
	}
	return vm.CallMethod(method, args...)
}

func (vm *Thread) CallMethod(method value.Method, args ...value.Value) (value.Value, value.Value) {
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

func (vm *Thread) ClearStackFrames() {
	vm.cfpSet(&vm.callFrames[0])
}

func (vm *Thread) ResetError() {
	vm.state = idleState
	vm.errStackTrace = nil
}

func (vm *Thread) GetStackTrace() *value.StackTrace {
	if vm.errStackTrace != nil {
		return vm.errStackTrace
	}

	return vm.BuildStackTrace()
}

func (vm *Thread) AddCallFrame(cf value.CallFrame) {
	*vm.cfpGet() = cf
	vm.cfpIncrement()
}

func (vm *Thread) PopCallFrame() {
	vm.cfpDecrementBy(1)
}

func (vm *Thread) BuildStackTrace() *value.StackTrace {
	callStack := vm.callStack()
	newCallStack := slices.Clone(callStack)
	return (*value.StackTrace)(&newCallStack)
}

func (vm *Thread) BuildStackTracePrepend(base *value.StackTrace) *value.StackTrace {
	callStack := vm.callStack()

	stackTraceSlice := make([]value.CallFrame, 0, len(*base)+len(callStack)+1)
	stackTraceSlice = append(stackTraceSlice, callStack...)
	stackTraceSlice = append(stackTraceSlice, (*base)...)

	return (*value.StackTrace)(&stackTraceSlice)
}

// Increment the call frame pointer
func (vm *Thread) cfpIncrement() {
	vm.cfpIncrementBy(1)
}

// Add n to the call frame pointer
func (vm *Thread) cfpIncrementBy(n int) {
	ptr := vm.cfpAdd(n)
	if ptr.IsSentinel() {
		panic("call stack overflow")
	}
	vm.cfpSet(ptr)
}

func (vm *Thread) lastCallFrame() *value.CallFrame {
	return vm.cfpAdd(-1)
}

func (vm *Thread) cfpAdd(n int) *value.CallFrame {
	return vm.callFrameAdd(vm.cfpGet(), n)
}

func (vm *Thread) cfpAddRaw(n uintptr) uintptr {
	return vm.cfp + n*CallFrameSize
}

func (vm *Thread) cfpSubtractRaw(n uintptr) uintptr {
	return vm.cfp - n*CallFrameSize
}

func (vm *Thread) cfpGet() *value.CallFrame {
	return (*value.CallFrame)(unsafe.Pointer(vm.cfp))
}

// Set the typesafe call frame pointer
func (vm *Thread) cfpSet(ptr *value.CallFrame) {
	vm.cfp = uintptr(unsafe.Pointer(ptr))
}

func (vm *Thread) callFrameAdd(ptr *value.CallFrame, n int) *value.CallFrame {
	return (*value.CallFrame)(unsafe.Add(unsafe.Pointer(ptr), n*int(CallFrameSize)))
}

func (vm *Thread) callFrameAddRaw(ptr uintptr, n uintptr) uintptr {
	return ptr + n*CallFrameSize
}

func (vm *Thread) cfpOffset() int {
	return int(vm.cfp-uintptr(unsafe.Pointer(&vm.callFrames[0]))) / int(CallFrameSize)
}

// Subtract n from the call frame pointer
func (vm *Thread) cfpDecrementBy(n uintptr) {
	vm.cfp = vm.cfpSubtractRaw(n)
}
