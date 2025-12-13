//go:build !native

// Package vm contains the Elk Virtual Machine.
// It interprets Elk Bytecode produced by
// the Elk compiler.
package vm

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// A single instance of the Elk Virtual Machine.
type VM struct {
	ID     int64
	Stdin  io.Reader // standard output used by the VM
	Stdout io.Writer // standard input used by the VM
	Stderr io.Writer // standard error used by the VM

	bytecode        *BytecodeFunction
	upvalues        []*Upvalue
	openUpvalueHead *Upvalue // linked list of open upvalues, living on the stack
	ip              uintptr  // Instruction pointer -- points to the next bytecode instruction
	sp              uintptr  // Stack pointer -- points to the offset where the next element will be pushed to
	fp              uintptr  // Frame pointer -- points to the offset where the section of the stack for current frame starts
	localCount      int      // the amount of registered locals
	cfp             uintptr  // Call frame pointer
	tailCallCounter int
	stack           []value.Value     // Value stack
	callFrames      []CallFrame       // Call stack
	errStackTrace   *value.StackTrace // The most recent error stack trace
	threadPool      *ThreadPool
	state           state
}

// Create a new VM instance.
func New(opts ...Option) *VM {
	stack := make([]value.Value, INIT_VALUE_STACK_SIZE)
	// mark the end of the value stack with a sentinel value
	stack[len(stack)-1] = value.MakeSentinelValue()

	callFrames := make([]CallFrame, CALL_STACK_SIZE)
	// mark the end of the call stack with a sentinel value
	callFrames[len(callFrames)-1] = makeSentinelCallFrame()

	id := currentID.Add(1)

	vm := &VM{
		ID:         id,
		stack:      stack,
		sp:         uintptr(unsafe.Pointer(&stack[0])),
		fp:         uintptr(unsafe.Pointer(&stack[0])),
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

// Execute the given bytecode chunk.
func (vm *VM) InterpretTopLevel(fn *BytecodeFunction) (value.Value, value.Value) {
	vm.bytecode = fn
	vm.ipSet(&fn.Instructions[0])
	vm.push(value.Ref(value.GlobalObject))
	vm.localCount = 1
	vm.tailCallCounter = 0
	vm.runWithState()
	err := vm.Err()
	if !err.IsUndefined() {
		return value.Undefined, err
	}
	return vm.peek(), value.Undefined
}

// Execute the given bytecode chunk.
func (vm *VM) InterpretREPL(fn *BytecodeFunction) (value.Value, value.Value) {
	vm.bytecode = fn
	vm.ipSet(&fn.Instructions[0])
	vm.tailCallCounter = 0
	if vm.sp == uintptr(unsafe.Pointer(&vm.stack[0])) {
		// populate the predeclared local variables
		vm.push(value.Ref(value.GlobalObject)) // populate self
		vm.localCount = 1
	} else {
		// pop the return value of the last run
		vm.pop()
	}
	vm.runWithState()

	err := vm.Err()
	if !err.IsUndefined() {
		return value.Undefined, err
	}
	return vm.peek(), value.Undefined
}

func (vm *VM) runWithState() {
	vm.state = runningState
	vm.run()
	if vm.state != errorState {
		vm.state = terminatedState
	}
}

// Get the stored error.
func (vm *VM) Err() value.Value {
	if vm.state == errorState {
		return vm.peek()
	}

	return value.Undefined
}

// Get the value on top of the value stack.
func (vm *VM) StackTop() value.Value {
	return vm.peek()
}

func (vm *VM) ValueStack() []value.Value {
	spIndex := vm.spOffset()
	return vm.stack[:spIndex]
}

func (vm *VM) callStack() []CallFrame {
	cfIndex := vm.cfpOffset()
	return vm.callFrames[:cfIndex]
}

func (vm *VM) stackFrame() []value.Value {
	spIndex := vm.spOffset()
	fpIndex := vm.fpOffset()
	return vm.stack[fpIndex:spIndex]
}

func (vm *VM) stackFrameCopy() []value.Value {
	stack := vm.stackFrame()
	stackCopy := make([]value.Value, len(stack))
	copy(stackCopy, stack)
	return stackCopy
}

func (vm *VM) InspectValueStack() {
	fmt.Println("value stack:")
	for i, value := range vm.ValueStack() {
		fmt.Printf("%d => %s\n", i, value.Inspect())
	}
}

func (vm *VM) InspectCallStack() {
	fmt.Println("call stack:")
	for i, frame := range vm.callStack() {
		fmt.Printf("%d => %#v\n", i, frame)
	}
}

func (vm *VM) throwIfErr(err value.Value) {
	if !err.IsUndefined() {
		vm.throw(err)
	}
}

func (vm *VM) callBytecodePromise(promise *Promise) {
	vm.state = runningState
	vm.createCurrentCallFrame(true)

	generator := promise.Body.(*Generator)
	vm.bytecode = generator.Bytecode
	vm.fp = vm.sp
	vm.ip = generator.ip
	vm.localCount = generator.Bytecode.parameterCount + 1
	vm.upvalues = generator.upvalues

	baseStack := &generator.stack[0]
	stackLen := len(generator.stack)
	for i := range stackLen {
		*vm.spAdd(i) = *vm.stackAdd(baseStack, i)
	}
	vm.spIncrementBy(uintptr(stackLen))

	vm.run()

	switch vm.state {
	case awaitState:
		stack := vm.stack[vm.fpOffset():vm.spOffset()]
		stackCopy := make([]value.Value, len(stack))
		copy(stackCopy, stack)
		generator.stack = stackCopy
		generator.ip = vm.ip

		vm.restoreLastFrame()
	case errorState:
		vm.restoreLastFrame()
	}
}

func (vm *VM) CallGeneratorNext(generator *Generator) (value.Value, value.Value) {
	vm.createCurrentCallFrame(true)
	vm.bytecode = generator.Bytecode
	vm.fp = vm.sp
	vm.ip = generator.ip
	vm.localCount = generator.Bytecode.parameterCount + 1
	vm.upvalues = generator.upvalues

	baseStack := &generator.stack[0]
	stackLen := len(generator.stack)
	for i := range stackLen {
		*vm.spAdd(i) = *vm.stackAdd(baseStack, i)
	}
	vm.spIncrementBy(uintptr(stackLen))

	vm.run()

	stack := vm.stack[vm.fpOffset() : vm.spOffset()-1]
	stackCopy := make([]value.Value, len(stack))
	copy(stackCopy, stack)
	generator.stack = stackCopy
	generator.ip = vm.ip

	if vm.state == errorState {
		vm.state = runningState
		vm.restoreLastFrame()
		inst := generator.Bytecode.Instructions
		// jump to the STOP_ITERATION at the end of the generator's bytecode
		// 1 byte for STOP_ITERATION, 3 bytes for LOOP
		generator.ip = uintptr(unsafe.Pointer(&inst[len(inst)-4]))
		return value.Undefined, vm.popGet()
	}

	vm.restoreLastFrame()
	return vm.popGet(), value.Undefined
}

// Call a callable value from Go code, preserving the state of the VM.
func (vm *VM) CallCallable(args ...value.Value) (value.Value, value.Value) {
	function := args[0]
	switch f := function.SafeAsReference().(type) {
	case *Closure:
		return vm.CallClosure(f, args[1:]...)
	default:
		return vm.CallMethodByName(callSymbol, args...)
	}
}

func (vm *VM) callGo(closure *Closure) {
	vm.bytecode = closure.Bytecode
	vm.fp = vm.sp
	vm.ipSet(&closure.Bytecode.Instructions[0])
	vm.localCount = 0
	vm.upvalues = closure.Upvalues
	// push `self`
	vm.push(closure.Self)

	if closure.VMID != vm.ID && closure.HasOpenUpvalues() {
		vm.throwNoCatch(
			value.Ref(
				value.NewOpenClosureError(
					closure.VMID,
					vm.ID,
					closure.Inspect(),
				),
			),
		)
	}

	vm.run()
}

// Call an Elk closure from Go code, preserving the state of the VM.
func (vm *VM) CallClosure(closure *Closure, args ...value.Value) (value.Value, value.Value) {
	if closure.VMID != vm.ID && closure.HasOpenUpvalues() {
		return value.Undefined, value.Ref(value.NewOpenClosureError(
			closure.VMID,
			vm.ID,
			closure.Inspect(),
		))
	}
	if closure.Bytecode.ParameterCount() != len(args) {
		return value.Undefined, value.Ref(value.NewWrongArgumentCountError(
			closure.Bytecode.Name().String(),
			len(args),
			closure.Bytecode.ParameterCount(),
		))
	}

	initialState := vm.state
	vm.createCurrentCallFrame(true)
	vm.bytecode = closure.Bytecode
	vm.fp = vm.sp
	vm.ipSet(&closure.Bytecode.Instructions[0])
	vm.localCount = len(args)
	vm.upvalues = closure.Upvalues
	// push `self`
	vm.push(closure.Self)
	for _, arg := range args {
		vm.push(arg)
	}
	vm.run()
	if vm.state == errorState {
		vm.restoreLastFrame()
		vm.state = initialState
		return value.Undefined, vm.popGet()
	}
	return vm.popGet(), value.Undefined
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

func (vm *VM) CallMethod(method value.Method, args ...value.Value) (value.Value, value.Value) {
	self := args[0]
	paramCount := method.ParameterCount()
	argCount := len(args) - 1
	args = vm.populateMissingParameters(args, paramCount, argCount)

	switch m := method.(type) {
	case *BytecodeFunction:
		initialState := vm.state
		vm.createCurrentCallFrame(true)
		vm.bytecode = m
		vm.fp = vm.sp
		vm.ipSet(&m.Instructions[0])
		vm.localCount = len(args)
		for _, arg := range args {
			vm.push(arg)
		}
		vm.run()
		if vm.state == errorState {
			vm.restoreLastFrame()
			vm.state = initialState
			return value.Undefined, vm.popGet()
		}
		return vm.popGet(), value.Undefined
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
func (vm *VM) callMethodOnStack(method value.Method, args int) value.Value {
	switch m := method.(type) {
	case *BytecodeFunction:
		vm.createCurrentCallFrame(false)
		vm.bytecode = m
		vm.fp = vm.spSubtractRaw(uintptr(args) + 1)
		vm.localCount = args + 1
		vm.ipSet(&m.Instructions[0])
	case *NativeMethod:
		argsPointer := vm.spAdd(-args - 1)
		result, err := m.Function(vm, unsafe.Slice(argsPointer, args+1))
		if !err.IsUndefined() {
			return err
		}
		vm.popN(args + 1)
		vm.push(result)
	default:
		panic(fmt.Sprintf("tried to call a method that is neither bytecode nor native: %#v", method))
	}

	return value.Undefined
}

func (vm *VM) callMethodOnStackByName(name value.Symbol, args int) value.Value {
	self := *vm.spAdd(-args - 1)
	class := self.DirectClass()
	method := class.LookupMethod(name)
	return vm.callMethodOnStack(method, args)
}

// The main execution loop of the VM.
func (vm *VM) run() {
	defer func() {
		// Return normally if the panic was an elk error
		r := recover()
		if r == nil || r == (stopVM{}) {
			return
		}
		panic(r)
	}()

	for {
		instruction := bytecode.OpCode(vm.readByte())
		switch instruction {
		case bytecode.STOP_ITERATION:
			vm.state = errorState
			vm.errStackTrace = vm.BuildStackTrace()
			vm.push(symbol.L_stop_iteration.ToValue())
			return
		case bytecode.YIELD:
			return
		case bytecode.RETURN_FINALLY:
			if vm.jumpToFinallyForReturn() {
				continue
			}

			// return normally
			if vm.cfp == uintptr(unsafe.Pointer(&vm.callFrames[0])) {
				return
			}
			if vm.restoreLastFrame() {
				return
			}
		case bytecode.RETURN:
			if vm.cfp == uintptr(unsafe.Pointer(&vm.callFrames[0])) {
				return
			}
			if vm.restoreLastFrame() {
				return
			}
		case bytecode.RETURN_FIRST_ARG:
			vm.opGetLocal(1)
			if vm.cfp == uintptr(unsafe.Pointer(&vm.callFrames[0])) {
				return
			}
			if vm.restoreLastFrame() {
				return
			}
		case bytecode.RETURN_SELF:
			vm.self()
			if vm.cfp == uintptr(unsafe.Pointer(&vm.callFrames[0])) {
				return
			}
			if vm.restoreLastFrame() {
				return
			}
		case bytecode.GO:
			vm.opGo()
		case bytecode.CLOSURE:
			vm.opClosure()
		case bytecode.CLOSED_CLOSURE:
			vm.opClosedClosure()
		case bytecode.JUMP_TO_FINALLY:
			leftFinallyCount := vm.peek().AsSmallInt()
			jumpOffset := vm.peekAt(1).AsSmallInt()

			if leftFinallyCount > 0 {
				vm.replace((leftFinallyCount - 1).ToValue())
				if !vm.jumpToFinallyForBreakOrContinue() {
					panic("could not find a finally block to jump to in JUMP_TO_FINALLY")
				}
				continue
			}

			vm.popN(2)
			vm.ipSetOffset(int(jumpOffset))
		case bytecode.GENERATOR:
			vm.opGenerator()
		case bytecode.PROMISE:
			vm.opPromise()
		case bytecode.AWAIT:
			promise := (*Promise)(vm.peek().Pointer())
			promise.m.Lock()

			if !promise.IsResolved() {
				// promise is not resolved, switching contexts
				vm.state = awaitState
				return
			}

			// promise is already resolved, no need to switch contexts
			err := promise.err
			result := promise.result
			stackTrace := promise.stackTrace
			promise.m.Unlock()

			if !err.IsUndefined() {
				vm.pop()
				vm.rethrow(err, vm.BuildStackTracePrepend(stackTrace))
				return
			}

			vm.replace(result)
			vm.ipIncrement() // skip over AWAIT_RESULT
		case bytecode.AWAIT_RESULT:
			promise := (*Promise)(vm.peek().Pointer())

			if !promise.IsResolved() {
				panic("promise is still unresolved after await")
			}

			result := promise.result
			err := promise.err
			stackTrace := promise.stackTrace
			if !err.IsUndefined() {
				vm.pop()
				vm.rethrow(err, vm.BuildStackTracePrepend(stackTrace))
				return
			}

			vm.replace(result)
		case bytecode.AWAIT_SYNC:
			promise := (*Promise)(vm.peek().Pointer())

			result, stackTrace, err := promise.AwaitSync()
			if !err.IsUndefined() {
				vm.pop()
				vm.rethrow(err, vm.BuildStackTracePrepend(stackTrace))
				return
			}

			vm.replace(result)
		case bytecode.NOOP:
		case bytecode.DUP:
			vm.push(vm.peek())
		case bytecode.DUP_SECOND:
			vm.push(vm.peekAt(1))
		case bytecode.SWAP:
			vm.swap()
		case bytecode.DUP_2:
			vm.push(*vm.spAdd(-2))
			vm.push(*vm.spAdd(-2))
		case bytecode.SELF:
			vm.self()
		case bytecode.INT_M1:
			vm.push(value.SmallInt(-1).ToValue())
		case bytecode.INT_0:
			vm.push(value.SmallInt(0).ToValue())
		case bytecode.INT_1:
			vm.push(value.SmallInt(1).ToValue())
		case bytecode.INT_2:
			vm.push(value.SmallInt(2).ToValue())
		case bytecode.INT_3:
			vm.push(value.SmallInt(3).ToValue())
		case bytecode.INT_4:
			vm.push(value.SmallInt(4).ToValue())
		case bytecode.INT_5:
			vm.push(value.SmallInt(5).ToValue())
		case bytecode.LOAD_INT_8:
			vm.push(value.SmallInt(int8(vm.readByte())).ToValue())
		case bytecode.LOAD_INT_16:
			vm.push(value.SmallInt(int16(vm.readUint16())).ToValue())
		case bytecode.LOAD_CHAR_8:
			vm.push(value.Char(int8(vm.readByte())).ToValue())
		case bytecode.FLOAT_0:
			vm.push(value.Float(0).ToValue())
		case bytecode.FLOAT_1:
			vm.push(value.Float(1).ToValue())
		case bytecode.FLOAT_2:
			vm.push(value.Float(2).ToValue())
		case bytecode.LOAD_INT64_8:
			vm.push(value.Int64(int8(vm.readByte())).ToValue())
		case bytecode.LOAD_UINT64_8:
			vm.push(value.UInt64(vm.readByte()).ToValue())
		case bytecode.LOAD_INT32_8:
			vm.push(value.Int32(int8(vm.readByte())).ToValue())
		case bytecode.LOAD_UINT32_8:
			vm.push(value.UInt32(vm.readByte()).ToValue())
		case bytecode.LOAD_INT16_8:
			vm.push(value.Int16(int8(vm.readByte())).ToValue())
		case bytecode.LOAD_UINT16_8:
			vm.push(value.UInt16(vm.readByte()).ToValue())
		case bytecode.LOAD_INT8:
			vm.push(value.Int8(vm.readByte()).ToValue())
		case bytecode.LOAD_UINT8:
			vm.push(value.UInt8(vm.readByte()).ToValue())
		case bytecode.DEF_NAMESPACE:
			vm.opDefNamespace()
		case bytecode.GET_SINGLETON:
			vm.throwIfErr(vm.opGetSingleton())
		case bytecode.GET_CLASS:
			vm.getClass()
		case bytecode.DEF_GETTER:
			vm.opDefGetter()
		case bytecode.DEF_SETTER:
			vm.opDefSetter()
		case bytecode.EXEC:
			vm.opExec()
		case bytecode.INIT_NAMESPACE:
			vm.opInitNamespace()
		case bytecode.DEF_METHOD:
			vm.opDefMethod()
		case bytecode.INCLUDE:
			vm.throwIfErr(vm.opInclude())
		case bytecode.DEF_IVARS:
			vm.throwIfErr(vm.opDefIvars())
		case bytecode.APPEND:
			vm.opAppend()
		case bytecode.MAP_SET:
			vm.opMapSet()
		case bytecode.COPY:
			vm.opCopy()
		case bytecode.APPEND_AT:
			vm.throwIfErr(vm.opAppendAt())
		case bytecode.SUBSCRIPT:
			vm.throwIfErr(vm.opSubscript())
		case bytecode.SUBSCRIPT_SET:
			vm.throwIfErr(vm.opSubscriptSet())
		case bytecode.INSTANTIATE8:
			vm.throwIfErr(
				vm.opInstantiate(int(vm.readByte())),
			)
		case bytecode.INSTANTIATE16:
			vm.throwIfErr(
				vm.opInstantiate(int(vm.readUint16())),
			)
		case bytecode.GET_IVAR_0:
			vm.opGetIvar(0)
		case bytecode.GET_IVAR_1:
			vm.opGetIvar(1)
		case bytecode.GET_IVAR_2:
			vm.opGetIvar(2)
		case bytecode.GET_IVAR8:
			vm.opGetIvar(int(vm.readByte()))
		case bytecode.GET_IVAR16:
			vm.opGetIvar(int(vm.readUint16()))
		case bytecode.GET_IVAR_NAME16:
			vm.throwIfErr(
				vm.opGetIvarName(int(vm.readUint16())),
			)
		case bytecode.SET_IVAR_0:
			vm.opSetIvar(0)
		case bytecode.SET_IVAR_1:
			vm.opSetIvar(1)
		case bytecode.SET_IVAR_2:
			vm.opSetIvar(2)
		case bytecode.SET_IVAR8:
			vm.opSetIvar(int(vm.readByte()))
		case bytecode.SET_IVAR16:
			vm.opSetIvar(int(vm.readUint16()))
		case bytecode.SET_IVAR_NAME16:
			vm.throwIfErr(
				vm.opSetIvarName(int(vm.readUint16())),
			)
		case bytecode.CALL_METHOD_TCO8:
			vm.throwIfErr(
				vm.opCallMethodTCO(int(vm.readByte())),
			)
		case bytecode.CALL_METHOD_TCO16:
			vm.throwIfErr(
				vm.opCallMethodTCO(int(vm.readUint16())),
			)
		case bytecode.CALL_METHOD8:
			vm.throwIfErr(
				vm.opCallMethod(int(vm.readByte())),
			)
		case bytecode.CALL_METHOD16:
			vm.throwIfErr(
				vm.opCallMethod(int(vm.readUint16())),
			)
		case bytecode.CALL8:
			vm.throwIfErr(
				vm.opCall(int(vm.readByte())),
			)
		case bytecode.CALL16:
			vm.throwIfErr(
				vm.opCall(int(vm.readUint16())),
			)
		case bytecode.CALL_SELF_TCO8:
			vm.throwIfErr(
				vm.opCallSelfTCO(int(vm.readByte())),
			)
		case bytecode.CALL_SELF_TCO16:
			vm.throwIfErr(
				vm.opCallSelfTCO(int(vm.readUint16())),
			)
		case bytecode.CALL_SELF8:
			vm.throwIfErr(
				vm.opCallSelf(int(vm.readByte())),
			)
		case bytecode.CALL_SELF16:
			vm.throwIfErr(
				vm.opCallSelf(int(vm.readUint16())),
			)
		case bytecode.INSTANCE_OF:
			vm.throwIfErr(vm.opInstanceOf())
		case bytecode.IS_A:
			vm.throwIfErr(vm.opIsA())
		case bytecode.UNDEFINED:
			vm.push(value.Undefined)
		case bytecode.LOAD_VALUE_0:
			vm.push(vm.readValue(0))
		case bytecode.LOAD_VALUE_1:
			vm.push(vm.readValue(1))
		case bytecode.LOAD_VALUE_2:
			vm.push(vm.readValue(2))
		case bytecode.LOAD_VALUE_3:
			vm.push(vm.readValue(3))
		case bytecode.LOAD_VALUE8:
			vm.push(vm.readValue8())
		case bytecode.LOAD_VALUE16:
			vm.push(vm.readValue16())
		case bytecode.ADD:
			vm.throwIfErr(vm.opAdd())
		case bytecode.ADD_INT:
			vm.opAddInt()
		case bytecode.ADD_FLOAT:
			vm.opAddFloat()
		case bytecode.SUBTRACT:
			vm.throwIfErr(vm.opSubtract())
		case bytecode.SUBTRACT_INT:
			vm.opSubtractInt()
		case bytecode.SUBTRACT_FLOAT:
			vm.opSubtractFloat()
		case bytecode.MULTIPLY:
			vm.throwIfErr(vm.opMultiply())
		case bytecode.MULTIPLY_INT:
			vm.opMultiplyInt()
		case bytecode.MULTIPLY_FLOAT:
			vm.opMultiplyFloat()
		case bytecode.DIVIDE:
			vm.throwIfErr(vm.opDivide())
		case bytecode.DIVIDE_INT:
			vm.opDivideInt()
		case bytecode.DIVIDE_FLOAT:
			vm.opDivideFloat()
		case bytecode.EXPONENTIATE:
			vm.throwIfErr(vm.opExponentiate())
		case bytecode.EXPONENTIATE_INT:
			vm.opExponentiateInt()
		case bytecode.NEGATE:
			vm.throwIfErr(vm.opNegate())
		case bytecode.NEGATE_INT:
			vm.opNegateInt()
		case bytecode.NEGATE_FLOAT:
			vm.opNegateFloat()
		case bytecode.UNARY_PLUS:
			vm.throwIfErr(vm.opUnaryPlus())
		case bytecode.BITWISE_NOT:
			vm.throwIfErr(vm.opBitwiseNot())
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
		case bytecode.POP_2:
			vm.popN(2)
		case bytecode.POP_2_SKIP_ONE:
			vm.popNSkipOne(2)
		case bytecode.POP_SKIP_ONE:
			vm.popSkipOne()
		case bytecode.INCREMENT:
			vm.throwIfErr(vm.opIncrement())
		case bytecode.INCREMENT_INT:
			vm.opIncrementInt()
		case bytecode.DECREMENT:
			vm.throwIfErr(vm.opDecrement())
		case bytecode.DECREMENT_INT:
			vm.opDecrementInt()
		case bytecode.GET_LOCAL_1:
			vm.opGetLocal(1)
		case bytecode.GET_LOCAL_2:
			vm.opGetLocal(2)
		case bytecode.GET_LOCAL_3:
			vm.opGetLocal(3)
		case bytecode.GET_LOCAL_4:
			vm.opGetLocal(4)
		case bytecode.GET_LOCAL8:
			vm.opGetLocal(int(vm.readByte()))
		case bytecode.GET_LOCAL16:
			vm.opGetLocal(int(vm.readUint16()))
		case bytecode.BOX_LOCAL8:
			vm.opBoxLocal(int(vm.readByte()))
		case bytecode.BOX_LOCAL16:
			vm.opBoxLocal(int(vm.readUint16()))
		case bytecode.SET_LOCAL_1:
			vm.opSetLocal(1)
		case bytecode.SET_LOCAL_2:
			vm.opSetLocal(2)
		case bytecode.SET_LOCAL_3:
			vm.opSetLocal(3)
		case bytecode.SET_LOCAL_4:
			vm.opSetLocal(4)
		case bytecode.SET_LOCAL8:
			vm.opSetLocal(int(vm.readByte()))
		case bytecode.SET_LOCAL16:
			vm.opSetLocal(int(vm.readUint16()))
		case bytecode.GET_UPVALUE_0:
			vm.opGetUpvalue(0)
		case bytecode.GET_UPVALUE_1:
			vm.opGetUpvalue(1)
		case bytecode.GET_UPVALUE8:
			vm.opGetUpvalue(int(vm.readByte()))
		case bytecode.GET_UPVALUE16:
			vm.opGetUpvalue(int(vm.readUint16()))
		case bytecode.SET_UPVALUE_0:
			vm.opSetUpvalue(0)
		case bytecode.SET_UPVALUE_1:
			vm.opSetUpvalue(1)
		case bytecode.SET_UPVALUE8:
			vm.opSetUpvalue(int(vm.readByte()))
		case bytecode.SET_UPVALUE16:
			vm.opSetUpvalue(int(vm.readUint16()))
		case bytecode.CLOSE_UPVALUES_TO_1:
			last := vm.fpAddRaw(1)
			vm.opCloseUpvalues(last)
		case bytecode.CLOSE_UPVALUES_TO_2:
			last := vm.fpAddRaw(2)
			vm.opCloseUpvalues(last)
		case bytecode.CLOSE_UPVALUES_TO_3:
			last := vm.fpAddRaw(3)
			vm.opCloseUpvalues(last)
		case bytecode.CLOSE_UPVALUES_TO8:
			last := vm.fpAddRaw(uintptr(vm.readByte()))
			vm.opCloseUpvalues(last)
		case bytecode.CLOSE_UPVALUES_TO16:
			last := vm.fpAddRaw(uintptr(vm.readUint16()))
			vm.opCloseUpvalues(last)
		case bytecode.PREP_LOCALS8:
			vm.opPrepLocals(uintptr(vm.readByte()))
		case bytecode.PREP_LOCALS16:
			vm.opPrepLocals(uintptr(vm.readUint16()))
		case bytecode.SET_SUPERCLASS:
			vm.opSetSuperclass()
		case bytecode.GET_CONST8:
			vm.throwIfErr(vm.opGetConst(int(vm.readByte())))
		case bytecode.GET_CONST16:
			vm.throwIfErr(
				vm.opGetConst(int(vm.readUint16())),
			)
		case bytecode.DEF_CONST:
			vm.opDefConst()
		case bytecode.NEW_RANGE:
			vm.opNewRange()
		case bytecode.NEW_ARRAY_TUPLE8:
			vm.opNewArrayTuple(int(vm.readByte()))
		case bytecode.NEW_ARRAY_TUPLE16:
			vm.opNewArrayTuple(int(vm.readUint16()))
		case bytecode.NEW_ARRAY_LIST8:
			vm.throwIfErr(
				vm.opNewArrayList(int(vm.readByte())),
			)
		case bytecode.NEW_ARRAY_LIST16:
			vm.throwIfErr(
				vm.opNewArrayList(int(vm.readUint16())),
			)
		case bytecode.NEW_HASH_SET8:
			vm.throwIfErr(
				vm.opNewHashSet(int(vm.readByte())),
			)
		case bytecode.NEW_HASH_SET16:
			vm.throwIfErr(
				vm.opNewHashSet(int(vm.readUint16())),
			)
		case bytecode.NEW_HASH_MAP8:
			vm.throwIfErr(
				vm.opNewHashMap(int(vm.readByte())),
			)
		case bytecode.NEW_HASH_MAP16:
			vm.throwIfErr(
				vm.opNewHashMap(int(vm.readUint16())),
			)
		case bytecode.NEW_HASH_RECORD8:
			vm.throwIfErr(
				vm.opNewHashRecord(int(vm.readByte())),
			)
		case bytecode.NEW_HASH_RECORD16:
			vm.throwIfErr(
				vm.opNewHashRecord(int(vm.readUint16())),
			)
		case bytecode.NEW_STRING8:
			vm.throwIfErr(vm.opNewString(int(vm.readByte())))
		case bytecode.NEW_STRING16:
			vm.throwIfErr(vm.opNewString(int(vm.readUint16())))
		case bytecode.NEW_SYMBOL8:
			vm.throwIfErr(vm.opNewSymbol(int(vm.readByte())))
		case bytecode.NEW_SYMBOL16:
			vm.throwIfErr(vm.opNewSymbol(int(vm.readUint16())))
		case bytecode.NEW_REGEX8:
			vm.throwIfErr(vm.opNewRegex(vm.readByte(), int(vm.readByte())))
		case bytecode.NEW_REGEX16:
			vm.throwIfErr(vm.opNewRegex(vm.readByte(), int(vm.readUint16())))
		case bytecode.NEXT8:
			vm.throwIfErr(vm.opNext(int(vm.readByte())))
		case bytecode.NEXT16:
			vm.throwIfErr(vm.opNext(int(vm.readUint16())))
		case bytecode.FOR_IN_BUILTIN:
			vm.opForInBuiltin()
		case bytecode.FOR_IN:
			vm.opForIn()
		case bytecode.GET_ITERATOR:
			vm.opGetIterator()
		case bytecode.JUMP_UNLESS:
			if value.Falsy(vm.popGet()) {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_NP:
			if value.Falsy(vm.peek()) {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_UNP:
			if !vm.peek().IsUndefined() {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_UNDEF:
			if !vm.popGet().IsUndefined() {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_LE:
			right := vm.popGet()
			left := vm.popGet()

			result, err := value.LessThanEqual(left, right)
			if !err.IsUndefined() {
				vm.throw(err)
				break
			}
			if !result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_LT:
			right := vm.popGet()
			left := vm.popGet()

			result, err := value.LessThan(left, right)
			if !err.IsUndefined() {
				vm.throw(err)
				break
			}
			if !result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_GE:
			right := vm.popGet()
			left := vm.popGet()

			result, err := value.GreaterThanEqual(left, right)
			if !err.IsUndefined() {
				vm.throw(err)
				break
			}
			if !result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_GT:
			right := vm.popGet()
			left := vm.popGet()

			result, err := value.GreaterThan(left, right)
			if !err.IsUndefined() {
				vm.throw(err)
				break
			}
			if !result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_EQ:
			right := vm.popGet()
			left := vm.popGet()

			result := value.Equal(left, right)
			if !result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_IF_EQ:
			right := vm.popGet()
			left := vm.popGet()

			result := value.Equal(left, right)
			if result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_ILE:
			right := vm.popGet()
			left := vm.popGet()

			var result bool
			if left.IsSmallInt() {
				left := left.AsSmallInt()
				result, _ = left.LessThanEqual(right)
			} else {
				leftBig := left.AsReference().(*value.BigInt)
				result, _ = leftBig.LessThanEqual(right)
			}
			if !result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_ILT:
			right := vm.popGet()
			left := vm.popGet()

			var result bool
			if left.IsSmallInt() {
				left := left.AsSmallInt()
				result, _ = left.LessThan(right)
			} else {
				leftBig := left.AsReference().(*value.BigInt)
				result, _ = leftBig.LessThan(right)
			}
			if !result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_IGE:
			right := vm.popGet()
			left := vm.popGet()

			var result bool
			if left.IsSmallInt() {
				left := left.AsSmallInt()
				result, _ = left.GreaterThanEqual(right)
			} else {
				leftBig := left.AsReference().(*value.BigInt)
				result, _ = leftBig.GreaterThanEqual(right)
			}
			if !result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_IGT:
			right := vm.popGet()
			left := vm.popGet()

			var result bool
			if left.IsSmallInt() {
				left := left.AsSmallInt()
				result, _ = left.GreaterThan(right)
			} else {
				leftBig := left.AsReference().(*value.BigInt)
				result, _ = leftBig.GreaterThan(right)
			}
			if !result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_IEQ:
			right := vm.popGet()
			left := vm.popGet()

			var result bool
			if left.IsSmallInt() {
				left := left.AsSmallInt()
				result = left.StrictEqualBool(right)
			} else {
				leftBig := left.AsReference().(*value.BigInt)
				result = leftBig.Equal(right)
			}
			if !result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_IF_IEQ:
			right := vm.popGet()
			left := vm.popGet()

			var result bool
			if left.IsSmallInt() {
				left := left.AsSmallInt()
				result = left.StrictEqualBool(right)
			} else {
				leftBig := left.AsReference().(*value.BigInt)
				result = leftBig.Equal(right)
			}
			if result {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_IF_NIL:
			if vm.popGet().IsNil() {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_IF_NIL_NP:
			if vm.peek().IsNil() {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_NIL:
			if !vm.popGet().IsNil() {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_UNLESS_NNP:
			if !vm.peek().IsNil() {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_IF:
			if value.Truthy(vm.popGet()) {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_IF_NP:
			if value.Truthy(vm.peek()) {
				jump := vm.readUint16()
				vm.ipIncrementBy(uintptr(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP:
			jump := vm.readUint16()
			vm.ipIncrementBy(uintptr(jump))
		case bytecode.LOOP:
			jump := vm.readUint16()
			vm.ipDecrementBy(uintptr(jump))
		case bytecode.THROW:
			vm.throw(vm.popGet())
		case bytecode.MUST:
			vm.opMust()
		case bytecode.AS:
			vm.opAs()
		case bytecode.RETHROW:
			err := vm.popGet()
			stackTrace := vm.popGet().AsReference().(*value.StackTrace)
			vm.rethrow(err, stackTrace)
		case bytecode.LBITSHIFT:
			vm.throwIfErr(vm.opLeftBitshift())
		case bytecode.LBITSHIFT_INT:
			vm.opLeftBitshiftInt()
		case bytecode.LOGIC_LBITSHIFT:
			vm.throwIfErr(vm.opLogicalLeftBitshift())
		case bytecode.RBITSHIFT:
			vm.throwIfErr(vm.opRightBitshift())
		case bytecode.RBITSHIFT_INT:
			vm.opRightBitshiftInt()
		case bytecode.LOGIC_RBITSHIFT:
			vm.throwIfErr(vm.opLogicalRightBitshift())
		case bytecode.BITWISE_AND:
			vm.throwIfErr(vm.opBitwiseAnd())
		case bytecode.BITWISE_AND_INT:
			vm.opBitwiseAndInt()
		case bytecode.BITWISE_AND_NOT:
			vm.throwIfErr(vm.opBitwiseAndNot())
		case bytecode.BITWISE_OR:
			vm.throwIfErr(vm.opBitwiseOr())
		case bytecode.BITWISE_OR_INT:
			vm.opBitwiseOrInt()
		case bytecode.BITWISE_XOR:
			vm.throwIfErr(vm.opBitwiseXor())
		case bytecode.BITWISE_XOR_INT:
			vm.opBitwiseXorInt()
		case bytecode.MODULO:
			vm.throwIfErr(vm.opModulo())
		case bytecode.MODULO_INT:
			vm.opModuloInt()
		case bytecode.MODULO_FLOAT:
			vm.opModuloFloat()
		case bytecode.COMPARE:
			vm.throwIfErr(vm.opCompare())
		case bytecode.EQUAL:
			vm.throwIfErr(vm.opEqual())
		case bytecode.EQUAL_INT:
			vm.opEqualInt()
		case bytecode.EQUAL_FLOAT:
			vm.opEqualFloat()
		case bytecode.NOT_EQUAL:
			vm.throwIfErr(vm.opNotEqual())
		case bytecode.NOT_EQUAL_INT:
			vm.opNotEqualInt()
		case bytecode.NOT_EQUAL_FLOAT:
			vm.opNotEqualFloat()
		case bytecode.LAX_EQUAL:
			vm.throwIfErr(vm.opLaxEqual())
		case bytecode.LAX_NOT_EQUAL:
			vm.throwIfErr(vm.opLaxNotEqual())
		case bytecode.STRICT_EQUAL:
			vm.opStrictEqual()
		case bytecode.STRICT_NOT_EQUAL:
			vm.opStrictNotEqual()
		case bytecode.GREATER:
			vm.throwIfErr(vm.opGreaterThan())
		case bytecode.GREATER_INT:
			vm.opGreaterThanInt()
		case bytecode.GREATER_FLOAT:
			vm.opGreaterThanFloat()
		case bytecode.GREATER_EQUAL:
			vm.throwIfErr(vm.opGreaterThanEqual())
		case bytecode.GREATER_EQUAL_I:
			vm.opGreaterThanEqualInt()
		case bytecode.GREATER_EQUAL_F:
			vm.opGreaterThanEqualFloat()
		case bytecode.LESS:
			vm.throwIfErr(vm.opLessThan())
		case bytecode.LESS_INT:
			vm.opLessThanInt()
		case bytecode.LESS_FLOAT:
			vm.opLessThanFloat()
		case bytecode.LESS_EQUAL:
			vm.throwIfErr(vm.opLessThanEqual())
		case bytecode.LESS_EQUAL_INT:
			vm.opLessThanEqualInt()
		case bytecode.LESS_EQUAL_FLOAT:
			vm.opLessThanEqualFloat()
		case bytecode.INSPECT_STACK:
			vm.InspectValueStack()
		default:
			panic(fmt.Sprintf("Unknown bytecode instruction: %#v", instruction))
		}
	}
}

// Creates a new VM instance.
// Spins up a new goroutine and executes the closure on top of the stack in it.
func (vm *VM) opGo() {
	closure := (*Closure)(vm.peek().Pointer())
	thread := New(WithStdin(vm.Stdin), WithStdout(vm.Stdout), WithStderr(vm.Stderr))

	go func(closure *Closure, thread *VM) {
		thread.state = runningState
		thread.callGo(closure)
		if thread.state != errorState {
			thread.state = terminatedState
			return
		}

		thread.PrintError()
	}(closure, thread)

	vm.replace(value.Ref(thread))
}

func (vm *VM) opClosedClosure() {
	function := vm.peek().AsReference().(*BytecodeFunction)
	closure := NewClosure(vm.ID, function, vm.selfValue())
	vm.replace(value.Ref(closure))

	for i := 0; ; i++ {
		flagByte := vm.readByte()
		if flagByte == ClosureTerminatorFlag {
			break
		}

		flags := bitfield.BitField8FromInt(flagByte)
		var upIndex int
		if flags.HasFlag(UpvalueLongIndexFlag) {
			upIndex = int(vm.readUint16())
		} else {
			upIndex = int(vm.readByte())
		}

		if flags.HasFlag(UpvalueLocalFlag) {
			closure.Upvalues[i] = NewClosedUpvalue(*vm.fpAdd(upIndex))
		} else {
			closure.Upvalues[i] = NewClosedUpvalue(vm.upvalues[upIndex].Get())
		}
	}
}

func (vm *VM) opClosure() {
	function := vm.peek().AsReference().(*BytecodeFunction)
	closure := NewClosure(vm.ID, function, vm.selfValue())
	vm.replace(value.Ref(closure))

	for i := 0; ; i++ {
		flagByte := vm.readByte()
		if flagByte == ClosureTerminatorFlag {
			break
		}

		flags := bitfield.BitField8FromInt(flagByte)
		var upIndex int
		if flags.HasFlag(UpvalueLongIndexFlag) {
			upIndex = int(vm.readUint16())
		} else {
			upIndex = int(vm.readByte())
		}

		if flags.HasFlag(UpvalueLocalFlag) {
			closure.Upvalues[i] = vm.captureUpvalue(vm.fpAdd(upIndex))
		} else {
			closure.Upvalues[i] = vm.upvalues[upIndex]
		}
	}
}

func (vm *VM) captureUpvalue(slot *value.Value) *Upvalue {
	var prevUpvalue *Upvalue
	currentUpvalue := vm.openUpvalueHead
	for {
		if currentUpvalue == nil ||
			(uintptr)(unsafe.Pointer(currentUpvalue.slot)) <=
				(uintptr)(unsafe.Pointer(slot)) {
			break
		}
		prevUpvalue = currentUpvalue
		currentUpvalue = currentUpvalue.next
	}

	if currentUpvalue != nil && currentUpvalue.slot == slot {
		return currentUpvalue
	}

	newUpvalue := NewUpvalue(slot)
	newUpvalue.next = currentUpvalue
	if prevUpvalue != nil {
		prevUpvalue.next = newUpvalue
	} else {
		vm.openUpvalueHead = newUpvalue
	}
	return newUpvalue
}

func (vm *VM) ClearStackFrames() {
	vm.cfpSet(&vm.callFrames[0])
}

// Restore the state of the VM to the last call frame.
//
//go:inline
func (vm *VM) restoreLastFrame() bool {
	vm.cfpDecrementBy(1)
	cf := vm.cfpGet()

	returnValue := vm.peek()
	vm.ip = cf.ip
	vm.opCloseUpvalues(vm.fp)
	vm.popN(vm.spOffsetTo(vm.fpGet()) - 1)
	vm.fp = cf.fp
	vm.localCount = cf.localCount
	vm.bytecode = cf.bytecode
	vm.upvalues = cf.upvalues
	*vm.spAdd(-1) = returnValue
	return cf.stopVM
}

func (vm *VM) PopCallFrame() {
	vm.restoreLastFrame()
}

func (vm *VM) PrintErrorValue(err value.Value) {
	PrintError(vm.Stderr, vm.ErrStackTrace(), err)
}

func (vm *VM) PrintError() {
	vm.PrintErrorValue(vm.Err())
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

func (vm *VM) makeCallFrameObject() value.CallFrame {
	return value.CallFrame{
		LineNumber:      vm.bytecode.GetLineNumber(vm.ipOffset() - 1),
		FileName:        vm.bytecode.FileName(),
		FuncName:        vm.bytecode.Name().String(),
		TailCallCounter: vm.tailCallCounter,
	}
}

func (vm *VM) BuildStackTrace() *value.StackTrace {
	callStack := vm.callStack()

	stackTraceSlice := make([]value.CallFrame, 0, len(callStack)+1)
	for _, element := range callStack {
		if element.bytecode == nil {
			continue
		}
		stackTraceSlice = append(stackTraceSlice, element.ToCallFrameObject())
	}
	stackTraceSlice = append(stackTraceSlice, vm.makeCallFrameObject())

	return (*value.StackTrace)(&stackTraceSlice)
}

func (vm *VM) BuildStackTracePrepend(base *value.StackTrace) *value.StackTrace {
	callStack := vm.callStack()

	stackTraceSlice := make([]value.CallFrame, 0, len(*base)+len(callStack)+1)
	for _, element := range callStack {
		if element.bytecode == nil {
			continue
		}
		stackTraceSlice = append(stackTraceSlice, element.ToCallFrameObject())
	}
	stackTraceSlice = append(stackTraceSlice, vm.makeCallFrameObject())
	stackTraceSlice = append(stackTraceSlice, (*base)...)

	return (*value.StackTrace)(&stackTraceSlice)
}

// Treat the next 8 bits of bytecode as an index
// of a value and retrieve the value.
func (vm *VM) readValue(i int) value.Value {
	return vm.bytecode.Values[i]
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

// Increment the stack pointer
func (vm *VM) spIncrement() {
	vm.spIncrementBy(1)
}

func (vm *VM) spSet(ptr *value.Value) {
	vm.sp = uintptr(unsafe.Pointer(ptr))
}

func (vm *VM) spGet() *value.Value {
	return (*value.Value)(unsafe.Pointer(vm.sp))
}

func (vm *VM) spOffsetTo(ptr *value.Value) int {
	return int(vm.sp-uintptr(unsafe.Pointer(ptr))) / int(value.ValueSize)
}

func (vm *VM) spOffset() int {
	return vm.spOffsetTo(&vm.stack[0])
}

func (vm *VM) spAdd(n int) *value.Value {
	return vm.stackAdd(vm.spGet(), n)
}

func (vm *VM) spAddRaw(n uintptr) uintptr {
	return vm.sp + n*uintptr(value.ValueSize)
}

func (vm *VM) spSubtractRaw(n uintptr) uintptr {
	return vm.sp - n*uintptr(value.ValueSize)
}

func (vm *VM) stackAdd(ptr *value.Value, n int) *value.Value {
	return (*value.Value)(unsafe.Add(unsafe.Pointer(ptr), n*int(value.ValueSize)))
}

func (vm *VM) stackAddRaw(ptr uintptr, n uintptr) uintptr {
	return ptr + n*value.ValueSize
}

func (vm *VM) stackOffsetFromTo(from *value.Value, to *value.Value) int {
	return int(uintptr(unsafe.Pointer(from))-uintptr(unsafe.Pointer(to))) / int(value.ValueSize)
}

func (vm *VM) stackOffsetFromToRaw(from, to uintptr) int {
	return int(from-to) / int(value.ValueSize)
}

func (vm *VM) fpOffset() int {
	return int(uintptr(unsafe.Pointer(vm.fp))-uintptr(unsafe.Pointer(&vm.stack[0]))) / int(value.ValueSize)
}

func (vm *VM) fpAdd(n int) *value.Value {
	return (*value.Value)(unsafe.Add(unsafe.Pointer(vm.fp), n*int(value.ValueSize)))
}

func (vm *VM) fpAddRaw(n uintptr) uintptr {
	return vm.fp + n*value.ValueSize
}

func (vm *VM) fpSet(ptr *value.Value) {
	vm.fp = uintptr(unsafe.Pointer(ptr))
}

func (vm *VM) fpGet() *value.Value {
	return (*value.Value)(unsafe.Pointer(vm.fp))
}

func (vm *VM) ipOffset() int {
	return int(
		vm.ip -
			uintptr(unsafe.Pointer(&vm.bytecode.Instructions[0])),
	)
}

func (vm *VM) ipSetOffset(offset int) {
	vm.ipSet((*byte)(unsafe.Add(unsafe.Pointer(&vm.bytecode.Instructions[0]), offset)))
}

// Get the typesafe instruction pointer
func (vm *VM) ipGet() *byte {
	return (*byte)(unsafe.Pointer(vm.ip))
}

// Set the typesafe instruction pointer
func (vm *VM) ipSet(ptr *byte) {
	vm.ip = uintptr(unsafe.Pointer(ptr))
}

// Increment the instruction pointer
func (vm *VM) ipIncrement() {
	vm.ipIncrementBy(1)
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

// Read the next byte of code
func (vm *VM) readByte() byte {
	byt := *vm.ipGet()
	vm.ipIncrement()
	return byt
}

// Read the next 2 bytes of code
func (vm *VM) readUint16() uint16 {
	// BENCHMARK: compare manual bit shifts
	result := binary.BigEndian.Uint16(unsafe.Slice(vm.ipGet(), 2))
	vm.ipIncrementBy(2)

	return result
}

// Read the next 4 bytes of code
func (vm *VM) readUint32() uint32 {
	// BENCHMARK: compare manual bit shifts
	result := binary.BigEndian.Uint32(unsafe.Slice(vm.ipGet(), 4))
	vm.ipIncrementBy(4)

	return result
}

func (vm *VM) self() {
	vm.opGetLocal(0)
}

func (vm *VM) opDefNamespace() {
	typ := vm.readByte()
	name := vm.popGet().AsInlineSymbol()
	parentNamespace := vm.popGet()

	var parentConstantContainer value.ConstantContainer
	switch n := parentNamespace.SafeAsReference().(type) {
	case *value.Class:
		parentConstantContainer = n.ConstantContainer
	case *value.Module:
		parentConstantContainer = n.ConstantContainer
	case *value.Interface:
		parentConstantContainer = n.ConstantContainer
	default:
		panic(
			fmt.Sprintf(
				"tried to define %s under an invalid namespace: `%s`",
				name,
				parentNamespace.Inspect(),
			),
		)
	}

	if _, ok := parentConstantContainer.Constants[name]; ok {
		return
	}

	var newNamespace value.Value
	switch typ {
	case bytecode.DEF_MODULE_FLAG:
		newNamespace = value.Ref(value.NewModule())
	case bytecode.DEF_CLASS_FLAG:
		newNamespace = value.Ref(value.NewClassWithOptions(value.ClassWithSuperclass(nil)))
	case bytecode.DEF_MIXIN_FLAG:
		newNamespace = value.Ref(value.NewMixin())
	case bytecode.DEF_INTERFACE_FLAG:
		newNamespace = value.Ref(value.NewInterface())
	}

	parentConstantContainer.AddConstant(name, newNamespace)
}

func (vm *VM) opGetSingleton() (err value.Value) {
	val := vm.popGet()
	singleton := val.SingletonClass()
	if singleton == nil {
		return value.Ref(value.Errorf(
			value.TypeErrorClass,
			"value `%s` cannot have a singleton class",
			val.Inspect(),
		))
	}

	vm.push(value.Ref(singleton))
	return value.Undefined
}

func (vm *VM) getClass() {
	val := vm.popGet()
	class := val.Class()
	vm.push(value.Ref(class))
}

func (vm *VM) selfValue() value.Value {
	return vm.getLocalValue(0)
}

func (vm *VM) lookupMethod(class *value.Class, callInfo *value.CallSiteInfo, index int) value.Method {
	for i := range len(callInfo.Cache) {
		cacheEntry := callInfo.Cache[i]
		if cacheEntry.Class == class {
			return cacheEntry.Method
		}
		if cacheEntry.Class == nil {
			method := class.LookupMethod(callInfo.Name)
			newCache := callInfo.Cache
			newCache[i] = value.CallCacheEntry{
				Class:  class,
				Method: method,
			}
			vm.bytecode.Values[index] = value.Ref(&value.CallSiteInfo{
				Name:          callInfo.Name,
				ArgumentCount: callInfo.ArgumentCount,
				Cache:         newCache,
			})
			return method
		}
	}

	return class.LookupMethod(callInfo.Name)
}

// Call a method with an implicit receiver
func (vm *VM) opCallSelfTCO(callInfoIndex int) (err value.Value) {
	callInfo := (*value.CallSiteInfo)(vm.bytecode.Values[callInfoIndex].Pointer())

	self := vm.selfValue()
	class := self.DirectClass()

	// shift all arguments one slot forward to make room for self
	for i := range callInfo.ArgumentCount {
		*vm.spAdd(-i) = *vm.spAdd(-i - 1)
	}
	*vm.spAdd(-callInfo.ArgumentCount) = self
	vm.spIncrement()

	method := vm.lookupMethod(class, callInfo, callInfoIndex)
	switch m := method.(type) {
	case *BytecodeFunction:
		vm.callBytecodeFunctionTCO(m, callInfo)
		return value.Undefined
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	case *GetterMethod:
		return vm.callGetterMethod(m)
	case *SetterMethod:
		return vm.callSetterMethod(m)
	}

	panic(
		fmt.Sprintf(
			"tried to call a method that is neither bytecode nor native: %#v, %s in %s",
			method,
			callInfo.Name,
			class.Name,
		),
	)
}

// Call a method with an implicit receiver
func (vm *VM) opCallSelf(callInfoIndex int) (err value.Value) {
	callInfo := (*value.CallSiteInfo)(vm.bytecode.Values[callInfoIndex].Pointer())

	self := vm.selfValue()
	class := self.DirectClass()
	// shift all arguments one slot forward to make room for self
	for i := range callInfo.ArgumentCount {
		*vm.spAdd(-i) = *vm.spAdd(-i - 1)
	}
	*vm.spAdd(-callInfo.ArgumentCount) = self
	vm.spIncrement()

	method := vm.lookupMethod(class, callInfo, callInfoIndex)
	switch m := method.(type) {
	case *BytecodeFunction:
		vm.callBytecodeFunction(m, callInfo)
		return value.Undefined
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	case *GetterMethod:
		return vm.callGetterMethod(m)
	case *SetterMethod:
		return vm.callSetterMethod(m)
	}

	panic(
		fmt.Sprintf(
			"tried to call a method that is neither bytecode nor native: %#v, %s in %s (%s)",
			method,
			callInfo.Name,
			class.Name,
			self.Inspect(),
		),
	)
}

func (vm *VM) callGetterMethod(method *GetterMethod) value.Value {
	self := vm.popGet() // pop self
	result, err := method.Call(self)
	if !err.IsUndefined() {
		return err
	}
	vm.push(result)
	return value.Undefined
}

func (vm *VM) callSetterMethod(method *SetterMethod) value.Value {
	other := vm.popGet()
	self := vm.popGet() // pop self
	result, err := method.Call(self, other)
	if !err.IsUndefined() {
		return err
	}
	vm.push(result)
	return value.Undefined
}

// Set the value of an instance variable by name
func (vm *VM) opSetIvarName(nameIndex int) (err value.Value) {
	name := vm.bytecode.Values[nameIndex].AsInlineSymbol()
	val := vm.popGet()
	self := vm.selfValue()
	return value.SetInstanceVariableByName(self, name, val)
}

// Set the value of an instance variable by index
func (vm *VM) opSetIvar(index int) {
	val := vm.popGet()
	self := vm.selfValue()
	value.SetInstanceVariable(self, index, val)
}

// Get the value of an instance variable by name
func (vm *VM) opGetIvarName(nameIndex int) (err value.Value) {
	name := vm.bytecode.Values[nameIndex].AsInlineSymbol()
	self := vm.selfValue()
	val, err := value.GetInstanceVariableByName(self, name)
	if !err.IsUndefined() {
		return err
	}

	if val.IsUndefined() {
		vm.push(value.Nil)
	} else {
		vm.push(val)
	}

	return value.Undefined
}

// Get the value of an instance variable by name
func (vm *VM) opGetIvar(index int) {
	self := vm.selfValue()
	val := self.InstanceVariables().Get(index)

	if val.IsUndefined() {
		vm.push(value.Nil)
	} else {
		vm.push(val)
	}
}

// Create a new generator
func (vm *VM) opGenerator() {
	generator := newGenerator(
		vm.bytecode,
		vm.upvalues,
		vm.stackFrameCopy(),
		vm.ip+1,
	)
	vm.push(value.Ref(generator))
}

// Create a new promise
func (vm *VM) opPromise() {
	arg := vm.popGet()
	generator := newGenerator(
		vm.bytecode,
		vm.upvalues,
		vm.stackFrameCopy(),
		vm.ip+1,
	)

	var threadPool *ThreadPool
	if arg.IsUndefined() {
		threadPool = vm.threadPool
	} else {
		threadPool = (*ThreadPool)(arg.Pointer())
	}

	promise := NewPromise(threadPool, generator)
	vm.push(value.Ref(promise))
}

// Pop the value on top of the stack and push its opCopy.
func (vm *VM) opCopy() {
	element := vm.peek()
	vm.replace(element.Copy())
}

// Set the value under the given key in a hash-map or hash-record
func (vm *VM) opMapSet() {
	val := vm.popGet()
	key := vm.popGet()
	collection := vm.peek()

	switch c := collection.SafeAsReference().(type) {
	case *value.HashMap:
		HashMapSet(vm, c, key, val)
	case *value.HashRecord:
		HashRecordSet(vm, c, key, val)
	default:
		panic(fmt.Sprintf("invalid map to set a value in: %s", collection.Inspect()))
	}
}

// Append an element to a list, arrayTuple or hashSet.
func (vm *VM) opAppend() {
	element := vm.popGet()
	collection := vm.peek()

	if collection.IsUndefined() {
		vm.replace(value.Ref(&value.ArrayTuple{element}))
		return
	}
	switch c := collection.SafeAsReference().(type) {
	case *value.ArrayTuple:
		c.Append(element)
	case *value.ArrayList:
		c.Append(element)
	case *value.HashSet:
		HashSetAppend(vm, c, element)
	default:
		panic(fmt.Sprintf("invalid collection to append to: %s", collection.Inspect()))
	}
}

// Create a new instance of a class
func (vm *VM) opInstantiate(args int) (err value.Value) {
	callInfo := value.NewCallSiteInfo(symbol.S_init, args)
	classPtr := vm.spAdd(-callInfo.ArgumentCount - 1)
	classVal := *classPtr
	var class *value.Class
	switch c := classVal.SafeAsReference().(type) {
	case *value.Class:
		class = c
	default:
		class = classVal.Class()
	}

	instance := class.CreateInstance()
	// replace the class with the instance
	*classPtr = instance
	method := class.LookupMethod(callInfo.Name)

	switch m := method.(type) {
	case *BytecodeFunction:
		vm.callBytecodeFunction(m, callInfo)
		return value.Undefined
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	case nil:
		// no initialiser defined
		// no arguments given
		// just replace the class with the instance
		return value.Undefined
	default:
		panic(fmt.Sprintf("tried to call an invalid initialiser method: %#v", method))
	}
}

// Call the `opCall` method with an explicit receiver
func (vm *VM) opCall(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].AsReference().(*value.CallSiteInfo)

	self, isClosure := vm.spAdd(-callInfo.ArgumentCount - 1).SafeAsReference().(*Closure)
	if !isClosure {
		return vm.opCallMethod(callInfoIndex)
	}

	return vm.callClosure(self, callInfo)
}

// set up the vm to execute a closure
func (vm *VM) callClosure(closure *Closure, callInfo *value.CallSiteInfo) (err value.Value) {
	if closure.VMID != vm.ID && closure.HasOpenUpvalues() {
		return value.Ref(value.NewOpenClosureError(
			closure.VMID,
			vm.ID,
			closure.Inspect(),
		))
	}

	function := closure.Bytecode
	vm.populateMissingParametersOnStack(function.parameterCount, callInfo.ArgumentCount)
	vm.createCurrentCallFrame(false)

	vm.localCount = function.parameterCount + 1
	vm.bytecode = function
	vm.fp = vm.spSubtractRaw(uintptr(function.parameterCount) + 1)
	vm.ipSet(&function.Instructions[0])
	vm.upvalues = closure.Upvalues

	return value.Undefined
}

// Call a method with an explicit receiver with tail call optimisation
func (vm *VM) opCallMethodTCO(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].AsReference().(*value.CallSiteInfo)

	selfPtr := vm.spAdd(-callInfo.ArgumentCount - 1)
	self := *selfPtr
	class := self.DirectClass()

	method := vm.lookupMethod(class, callInfo, callInfoIndex)
	switch m := method.(type) {
	case *BytecodeFunction:
		vm.callBytecodeFunctionTCO(m, callInfo)
		return value.Undefined
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	case *GetterMethod:
		return vm.callGetterMethod(m)
	case *SetterMethod:
		return vm.callSetterMethod(m)
	default:
		panic(fmt.Sprintf("tried to call an invalid method: %T (%s) of class: %s (%s)", method, callInfo.Name, class.Name, self.Inspect()))
	}
}

// Call a method with an explicit receiver
func (vm *VM) opCallMethod(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].AsReference().(*value.CallSiteInfo)

	selfPtr := vm.spAdd(-callInfo.ArgumentCount - 1)
	self := *selfPtr
	class := self.DirectClass()

	method := vm.lookupMethod(class, callInfo, callInfoIndex)
	switch m := method.(type) {
	case *BytecodeFunction:
		vm.callBytecodeFunction(m, callInfo)
		return value.Undefined
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	case *GetterMethod:
		return vm.callGetterMethod(m)
	case *SetterMethod:
		return vm.callSetterMethod(m)
	default:
		panic(fmt.Sprintf("tried to call an invalid method: %T (%s) of class: %s (%s)", method, callInfo.Name, class.Name, self.Inspect()))
	}
}

// set up the vm to execute a native method
func (vm *VM) callNativeMethod(method *NativeMethod, callInfo *value.CallSiteInfo) (err value.Value) {
	vm.populateMissingParametersOnStack(method.parameterCount, callInfo.ArgumentCount)

	paramCount := method.ParameterCount()
	args := unsafe.Slice(vm.spAdd(-paramCount-1), paramCount+1)
	returnVal, nativeErr := method.Function(vm, args)
	vm.popN(paramCount + 1)
	if !nativeErr.IsUndefined() {
		return nativeErr
	}
	vm.push(returnVal)
	return value.Undefined
}

// set up the vm to execute a bytecode method with tail call optimisation
func (vm *VM) callBytecodeFunctionTCO(method *BytecodeFunction, callInfo *value.CallSiteInfo) {
	vm.populateMissingParametersOnStack(method.parameterCount, callInfo.ArgumentCount)

	localCount := method.parameterCount + 1
	for i := range localCount {
		*vm.fpAdd(i) = *vm.spAdd(-localCount + i)
	}
	vm.popN(vm.localCount)

	vm.localCount = localCount
	vm.bytecode = method
	vm.ipSet(&method.Instructions[0])
	vm.tailCallCounter++
}

// set up the vm to execute a bytecode method
func (vm *VM) callBytecodeFunction(method *BytecodeFunction, callInfo *value.CallSiteInfo) {
	vm.populateMissingParametersOnStack(method.parameterCount, callInfo.ArgumentCount)
	vm.createCurrentCallFrame(false)
	vm.localCount = method.parameterCount + 1
	vm.bytecode = method
	vm.fp = vm.spSubtractRaw(uintptr(method.parameterCount + 1))
	vm.ipSet(&method.Instructions[0])

	if float64(vm.spOffset()) > 0.7*float64(len(vm.stack)) {
		vm.growValueStack()
	}
}

func (vm *VM) growValueStack() {
	newSize := len(vm.stack) * 2
	if newSize >= MAX_VALUE_STACK_SIZE {
		panic("maximum value stack size exceeded")
	}
	newStack := make([]value.Value, newSize)
	copy(newStack, vm.stack)
	newStack[len(newStack)-1] = value.MakeSentinelValue()
	oldStackPtr := uintptr(unsafe.Pointer(&vm.stack[0]))
	newStackPtr := uintptr(unsafe.Pointer(&newStack[0]))
	fpOffset := uintptr(vm.fpOffset())
	spOffset := uintptr(vm.spOffset())

	for i := range vm.callFrames {
		cf := &vm.callFrames[i]
		offset := uintptr(vm.stackOffsetFromToRaw(oldStackPtr, cf.fp))
		cf.fp = vm.stackAddRaw(newStackPtr, offset)
		for _, upvalue := range cf.upvalues {
			if upvalue.IsClosed() {
				continue
			}

			offset := vm.stackOffsetFromTo(&vm.stack[0], upvalue.slot)
			upvalue.slot = vm.stackAdd(&newStack[0], offset)
		}
	}

	for _, upvalue := range vm.upvalues {
		if upvalue.IsClosed() {
			continue
		}

		offset := vm.stackOffsetFromTo(&vm.stack[0], upvalue.slot)
		upvalue.slot = vm.stackAdd(&newStack[0], offset)
	}

	vm.fp = vm.stackAddRaw(newStackPtr, fpOffset)
	vm.sp = vm.stackAddRaw(newStackPtr, spOffset)
	vm.stack = newStack
}

func (vm *VM) populateMissingParametersOnStack(paramCount, argumentCount int) {
	// populate missing optional arguments with undefined
	missingParams := uintptr(paramCount - argumentCount)
	if missingParams > 0 {
		vm.spIncrementBy(missingParams)
	}
}

// Define instance variables in a class
func (vm *VM) opDefIvars() (err value.Value) {
	ivarIndices := (*value.IvarIndices)(vm.popGet().Pointer())
	classVal := vm.popGet()

	switch class := classVal.SafeAsReference().(type) {
	case *value.Class:
		class.IvarIndices = *ivarIndices
	default:
		return value.Ref(value.Errorf(
			value.TypeErrorClass,
			"cannot define instance variables in %s: `%s`",
			classVal.Class().PrintableName(),
			classVal.Inspect(),
		))
	}

	return value.Undefined
}

// Include a mixin in a class/mixin.
func (vm *VM) opInclude() (err value.Value) {
	mixinVal := vm.popGet()
	targetValue := vm.popGet()

	mixin, ok := mixinVal.AsReference().(*value.Mixin)
	if !ok || !mixin.IsMixin() {
		return value.Ref(value.NewIsNotMixinError(mixinVal.Inspect()))
	}

	switch target := targetValue.SafeAsReference().(type) {
	case *value.Class:
		target.IncludeMixin(mixin)
	default:
		return value.Ref(value.Errorf(
			value.TypeErrorClass,
			"cannot include into an instance of %s: `%s`",
			targetValue.Class().PrintableName(),
			targetValue.Inspect(),
		))
	}

	return value.Undefined
}

// Define a new method
func (vm *VM) opDefMethod() {
	name := vm.popGet().AsInlineSymbol()
	body := vm.popGet().AsReference().(value.Method)
	methodContainer := vm.peek()

	switch m := methodContainer.SafeAsReference().(type) {
	case *value.Class:
		m.Methods[name] = body
	default:
		panic(fmt.Sprintf("invalid method container: %s", methodContainer.Inspect()))
	}
}

// Initialise a namespace
func (vm *VM) opInitNamespace() {
	body := vm.popGet().AsReference().(*BytecodeFunction)
	namespace := vm.popGet()
	vm.executeNamespaceBody(namespace, body)
}

// Execute a chunk of bytecode
func (vm *VM) opExec() {
	bytecodeFunc := vm.popGet().AsReference().(*BytecodeFunction)
	vm.executeFunc(bytecodeFunc)
}

// Define a getter method
func (vm *VM) opDefGetter() {
	index := vm.popGet().AsSmallInt()
	name := vm.popGet().AsInlineSymbol()
	methodContainer := vm.peek()

	switch m := methodContainer.SafeAsReference().(type) {
	case *value.Class:
		DefineGetter(&m.MethodContainer, name, int(index))
	default:
		panic(fmt.Sprintf("cannot define a getter in an invalid method container: %s", methodContainer.Inspect()))
	}
}

// Define a setter method
func (vm *VM) opDefSetter() {
	index := vm.popGet().AsSmallInt()
	name := vm.popGet().AsInlineSymbol()
	methodContainer := vm.peek()

	switch m := methodContainer.SafeAsReference().(type) {
	case *value.Class:
		DefineSetter(&m.MethodContainer, name, int(index))
	default:
		panic(fmt.Sprintf("cannot define a setter in an invalid method container: %s", methodContainer.Inspect()))
	}
}

func (vm *VM) AddCallFrame(cf CallFrame) {
	*vm.cfpGet() = cf
	vm.cfpIncrement()
	vm.tailCallCounter = 0
}

// preserve the current state of the vm in a call frame
func (vm *VM) createCurrentCallFrame(stopVM bool) {
	vm.AddCallFrame(
		CallFrame{
			upvalues:        vm.upvalues,
			bytecode:        vm.bytecode,
			ip:              vm.ip,
			fp:              vm.fp,
			localCount:      vm.localCount,
			tailCallCounter: vm.tailCallCounter,
			stopVM:          stopVM,
		},
	)
}

// set up the vm to execute a namespace body
func (vm *VM) executeNamespaceBody(namespace value.Value, body *BytecodeFunction) {
	vm.createCurrentCallFrame(false)

	vm.bytecode = body
	vm.fp = vm.sp
	vm.ipSet(&body.Instructions[0])
	vm.localCount = 1
	// set namespace as `self`
	vm.push(namespace)
}

// set up the vm to execute a bytecode function
func (vm *VM) executeFunc(fn *BytecodeFunction) {
	vm.createCurrentCallFrame(false)

	vm.bytecode = fn
	vm.fp = vm.sp
	vm.ipSet(&fn.Instructions[0])
	vm.localCount = 1
	vm.push(value.Ref(value.GlobalObject))
}

// Set a local variable or value.
func (vm *VM) opSetLocal(index int) {
	vm.setLocalValue(index, vm.popGet())
}

// Set a local variable or value.
func (vm *VM) setLocalValue(index int, val value.Value) {
	*vm.fpAdd(index) = val
}

// Read a local variable or value.
func (vm *VM) opGetLocal(index int) {
	vm.push(vm.getLocalValue(index))
}

// Create a box that points to a local
func (vm *VM) opBoxLocal(localIndex int) {
	upvalue := vm.captureUpvalue(vm.fpAdd(localIndex))
	box := (*LocalBox)(upvalue)
	vm.push(value.Ref(box))
}

// Read a local variable or value.
func (vm *VM) getLocalValue(index int) value.Value {
	return *vm.fpAdd(index)
}

// Set an upvalue.
func (vm *VM) opSetUpvalue(index int) {
	vm.setUpvalueValue(index, vm.popGet())
}

// Set an upvalue.
func (vm *VM) setUpvalueValue(index int, val value.Value) {
	vm.upvalues[index].Set(val)
}

// Read an upvalue.
func (vm *VM) opGetUpvalue(index int) {
	vm.push(vm.getUpvalueValue(index))
}

// Read an upvalue.
func (vm *VM) getUpvalueValue(index int) value.Value {
	return vm.upvalues[index].Get()
}

// Closes all upvalues down to the given local slot (the given slot and all above).
func (vm *VM) opCloseUpvalues(lastToClose uintptr) {
	for {
		if vm.openUpvalueHead == nil ||
			uintptr(unsafe.Pointer(vm.openUpvalueHead.slot)) <
				lastToClose {
			break
		}

		currentUpvalue := vm.openUpvalueHead
		currentUpvalue.Close()
		vm.openUpvalueHead = currentUpvalue.next
	}
}

// Set the superclass/parent of a class
func (vm *VM) opSetSuperclass() {
	newSuperclass := vm.popGet().AsReference().(*value.Class)
	class := vm.popGet().AsReference().(*value.Class)

	if class.Parent != nil {
		return
	}

	class.SetSuperclass(newSuperclass)
}

// Look for a constant with the given name.
func (vm *VM) opGetConst(nameIndex int) (err value.Value) {
	symbol := vm.bytecode.Values[nameIndex].AsInlineSymbol()

	val := value.RootModule.Constants.Get(symbol)
	if val.IsUndefined() {
		return value.Ref(value.Errorf(value.NoConstantErrorClass, "undefined constant `%s`", symbol.String()))
	}

	vm.push(val)
	return value.Undefined
}

// Get the iterator of the value on top of the stack.
func (vm *VM) opGetIterator() {
	val := vm.peek()
	result := value.Iter(val)
	vm.replace(result)
}

var stopIterationSymbol = value.ToSymbol("stop_iteration")

// Get the next element of an iterator
func (vm *VM) opNext(callInfoIndex int) value.Value {
	callInfo := vm.bytecode.Values[callInfoIndex].AsReference().(*value.CallSiteInfo)
	iterator := vm.peek()

	method := vm.lookupMethod(iterator.DirectClass(), callInfo, callInfoIndex)
	result, err := vm.CallMethod(method, iterator)
	if err.IsInlineSymbol() && err.AsInlineSymbol() == stopIterationSymbol {
		vm.replace(value.Undefined)
		return value.Undefined
	}
	if !err.IsUndefined() {
		return err
	}

	vm.replace(result)
	return value.Undefined
}

// Drives the for..in loop.
func (vm *VM) opForIn() {
	result := vm.peek()
	if result.IsUndefined() {
		vm.pop()
		vm.ipIncrementBy(uintptr(vm.readUint16()))
		return
	}

	vm.ipIncrementBy(2)
}

// Drives the for..in loop for builtin iterable types
func (vm *VM) opForInBuiltin() {
	iterator := vm.peek()
	result, err := NextBuiltin(vm, iterator)
	if !err.IsUndefined() {
		vm.pop()
		vm.ipIncrementBy(uintptr(vm.readUint16()))
		return
	}

	vm.replace(result)
	vm.ipIncrementBy(2)
}

// Create a new string.
func (vm *VM) opNewString(dynamicElements int) value.Value {
	firstElement := vm.spAdd(-dynamicElements)

	var buffer strings.Builder
	for i := range dynamicElements {
		elementPtr := vm.stackAdd(firstElement, i)
		elementVal := *elementPtr
		*elementPtr = value.Undefined

		if elementVal.IsReference() {
			switch element := elementVal.AsReference().(type) {
			case value.String:
				buffer.WriteString(string(element))
			case value.Float64:
				buffer.WriteString(string(element.ToString()))
			case *value.BigInt:
				buffer.WriteString(string(element.ToString()))
			case value.Int64:
				buffer.WriteString(string(element.ToString()))
			case value.UInt64:
				buffer.WriteString(string(element.ToString()))
			case *value.Regex:
				buffer.WriteString(string(element.ToString()))
			default:
				strVal, err := vm.CallMethodByName(toStringSymbol, elementVal)
				if !err.IsUndefined() {
					return err
				}
				str, ok := strVal.SafeAsReference().(value.String)
				if !ok {
					return value.Ref(value.NewCoerceError(value.StringClass, strVal.Class()))
				}
				buffer.WriteString(string(str))
			}
			continue
		}

		switch elementVal.ValueFlag() {
		case value.CHAR_FLAG:
			element := elementVal.AsChar()
			buffer.WriteRune(rune(element))
		case value.FLOAT64_FLAG:
			element := elementVal.AsInlineFloat64()
			buffer.WriteString(string(element.ToString()))
		case value.FLOAT32_FLAG:
			element := elementVal.AsFloat32()
			buffer.WriteString(string(element.ToString()))
		case value.FLOAT_FLAG:
			element := elementVal.AsFloat()
			buffer.WriteString(string(element.ToString()))
		case value.SMALL_INT_FLAG:
			element := elementVal.AsSmallInt()
			buffer.WriteString(string(element.ToString()))
		case value.INT64_FLAG:
			element := elementVal.AsInlineInt64()
			buffer.WriteString(string(element.ToString()))
		case value.INT32_FLAG:
			element := elementVal.AsInt32()
			buffer.WriteString(string(element.ToString()))
		case value.INT16_FLAG:
			element := elementVal.AsInt16()
			buffer.WriteString(string(element.ToString()))
		case value.INT8_FLAG:
			element := elementVal.AsInt8()
			buffer.WriteString(string(element.ToString()))
		case value.UINT64_FLAG:
			element := elementVal.AsInlineUInt64()
			buffer.WriteString(string(element.ToString()))
		case value.UINT32_FLAG:
			element := elementVal.AsInt32()
			buffer.WriteString(string(element.ToString()))
		case value.UINT16_FLAG:
			element := elementVal.AsInt16()
			buffer.WriteString(string(element.ToString()))
		case value.UINT8_FLAG:
			element := elementVal.AsUInt8()
			buffer.WriteString(string(element.ToString()))
		case value.NIL_FLAG:
		case value.SYMBOL_FLAG:
			element := elementVal.AsInlineSymbol()
			buffer.WriteString(string(element.ToString()))
		default:
			strVal, err := vm.CallMethodByName(toStringSymbol, elementVal)
			if !err.IsUndefined() {
				return err
			}
			str, ok := strVal.SafeAsReference().(value.String)
			if !ok {
				return value.Ref(value.NewCoerceError(value.StringClass, strVal.Class()))
			}
			buffer.WriteString(string(str))
		}
	}

	vm.spDecrementBy(uintptr(dynamicElements))
	vm.push(value.Ref(value.String(buffer.String())))

	return value.Undefined
}

// Create a new symbol.
func (vm *VM) opNewSymbol(dynamicElements int) value.Value {
	firstElement := vm.spAdd(-dynamicElements)

	var buffer strings.Builder
	for i := range dynamicElements {
		elementPtr := vm.stackAdd(firstElement, i)
		elementVal := *elementPtr
		*elementPtr = value.Undefined

		if elementVal.IsReference() {
			switch element := elementVal.AsReference().(type) {
			case value.String:
				buffer.WriteString(string(element))
			case value.Float64:
				buffer.WriteString(string(element.ToString()))
			case *value.BigInt:
				buffer.WriteString(string(element.ToString()))
			case value.Int64:
				buffer.WriteString(string(element.ToString()))
			case value.UInt64:
				buffer.WriteString(string(element.ToString()))
			case *value.Regex:
				buffer.WriteString(string(element.ToString()))
			default:
				strVal, err := vm.CallMethodByName(toStringSymbol, elementVal)
				if !err.IsUndefined() {
					return err
				}
				str, ok := strVal.SafeAsReference().(value.String)
				if !ok {
					return value.Ref(value.NewCoerceError(value.StringClass, strVal.Class()))
				}
				buffer.WriteString(string(str))
			}
			continue
		}

		switch elementVal.ValueFlag() {
		case value.CHAR_FLAG:
			element := elementVal.AsChar()
			buffer.WriteRune(rune(element))
		case value.FLOAT64_FLAG:
			element := elementVal.AsInlineFloat64()
			buffer.WriteString(string(element.ToString()))
		case value.FLOAT32_FLAG:
			element := elementVal.AsFloat32()
			buffer.WriteString(string(element.ToString()))
		case value.FLOAT_FLAG:
			element := elementVal.AsFloat()
			buffer.WriteString(string(element.ToString()))
		case value.SMALL_INT_FLAG:
			element := elementVal.AsSmallInt()
			buffer.WriteString(string(element.ToString()))
		case value.INT64_FLAG:
			element := elementVal.AsInlineInt64()
			buffer.WriteString(string(element.ToString()))
		case value.INT32_FLAG:
			element := elementVal.AsInt32()
			buffer.WriteString(string(element.ToString()))
		case value.INT16_FLAG:
			element := elementVal.AsInt16()
			buffer.WriteString(string(element.ToString()))
		case value.INT8_FLAG:
			element := elementVal.AsInt8()
			buffer.WriteString(string(element.ToString()))
		case value.UINT64_FLAG:
			element := elementVal.AsInlineUInt64()
			buffer.WriteString(string(element.ToString()))
		case value.UINT32_FLAG:
			element := elementVal.AsUInt32()
			buffer.WriteString(string(element.ToString()))
		case value.UINT16_FLAG:
			element := elementVal.AsUInt16()
			buffer.WriteString(string(element.ToString()))
		case value.UINT8_FLAG:
			element := elementVal.AsUInt8()
			buffer.WriteString(string(element.ToString()))
		case value.NIL_FLAG:
		case value.SYMBOL_FLAG:
			element := elementVal.AsInlineSymbol()
			buffer.WriteString(string(element.ToString()))
		default:
			strVal, err := vm.CallMethodByName(toStringSymbol, elementVal)
			if !err.IsUndefined() {
				return err
			}
			str, ok := strVal.SafeAsReference().(value.String)
			if !ok {
				return value.Ref(value.NewCoerceError(value.StringClass, strVal.Class()))
			}
			buffer.WriteString(string(str))
		}
	}

	vm.spDecrementBy(uintptr(dynamicElements))
	vm.push(value.ToSymbol(buffer.String()).ToValue())

	return value.Undefined
}

// Create a new regex.
func (vm *VM) opNewRegex(flagByte byte, dynamicElements int) value.Value {
	flags := bitfield.BitField8FromInt(flagByte)
	firstElement := vm.spAdd(-dynamicElements)

	var buffer strings.Builder
	for i := range dynamicElements {
		elementPtr := vm.stackAdd(firstElement, i)
		elementVal := *elementPtr
		*elementPtr = value.Undefined

		if elementVal.IsReference() {
			switch element := elementVal.AsReference().(type) {
			case value.String:
				buffer.WriteString(string(element))
			case value.Float64:
				buffer.WriteString(string(element.ToString()))
			case *value.BigInt:
				buffer.WriteString(string(element.ToString()))
			case value.Int64:
				buffer.WriteString(string(element.ToString()))
			case value.UInt64:
				buffer.WriteString(string(element.ToString()))
			case *value.Regex:
				buffer.WriteString(string(element.ToStringWithFlags()))
			default:
				strVal, err := vm.CallMethodByName(toStringSymbol, elementVal)
				if !err.IsUndefined() {
					return err
				}
				str, ok := strVal.SafeAsReference().(value.String)
				if !ok {
					return value.Ref(value.NewCoerceError(value.StringClass, strVal.Class()))
				}
				buffer.WriteString(string(str))
			}
			continue
		}

		switch elementVal.ValueFlag() {
		case value.CHAR_FLAG:
			element := elementVal.AsChar()
			buffer.WriteRune(rune(element))
		case value.FLOAT64_FLAG:
			element := elementVal.AsInlineFloat64()
			buffer.WriteString(string(element.ToString()))
		case value.FLOAT32_FLAG:
			element := elementVal.AsFloat32()
			buffer.WriteString(string(element.ToString()))
		case value.FLOAT_FLAG:
			element := elementVal.AsFloat()
			buffer.WriteString(string(element.ToString()))
		case value.SMALL_INT_FLAG:
			element := elementVal.AsSmallInt()
			buffer.WriteString(string(element.ToString()))
		case value.INT64_FLAG:
			element := elementVal.AsInlineInt64()
			buffer.WriteString(string(element.ToString()))
		case value.INT32_FLAG:
			element := elementVal.AsInt32()
			buffer.WriteString(string(element.ToString()))
		case value.INT16_FLAG:
			element := elementVal.AsInt16()
			buffer.WriteString(string(element.ToString()))
		case value.INT8_FLAG:
			element := elementVal.AsInt8()
			buffer.WriteString(string(element.ToString()))
		case value.UINT64_FLAG:
			element := elementVal.AsInlineUInt64()
			buffer.WriteString(string(element.ToString()))
		case value.UINT32_FLAG:
			element := elementVal.AsUInt32()
			buffer.WriteString(string(element.ToString()))
		case value.UINT16_FLAG:
			element := elementVal.AsUInt16()
			buffer.WriteString(string(element.ToString()))
		case value.UINT8_FLAG:
			element := elementVal.AsUInt8()
			buffer.WriteString(string(element.ToString()))
		case value.NIL_FLAG:
		case value.SYMBOL_FLAG:
			element := elementVal.AsInlineSymbol()
			buffer.WriteString(string(element.ToString()))
		default:
			strVal, err := vm.CallMethodByName(toStringSymbol, elementVal)
			if !err.IsUndefined() {
				return err
			}
			str, ok := strVal.SafeAsReference().(value.String)
			if !ok {
				return value.Ref(value.NewCoerceError(value.StringClass, strVal.Class()))
			}
			buffer.WriteString(string(str))
		}
	}
	vm.spDecrementBy(uintptr(dynamicElements))
	re, err := value.CompileRegex(buffer.String(), flags)
	if err != nil {
		return value.Ref(value.NewError(value.RegexCompileErrorClass, err.Error()))
	}

	vm.push(value.Ref(re))
	return value.Undefined
}

// Create a new hashset.
func (vm *VM) opNewHashSet(dynamicElements int) value.Value {
	firstElement := vm.spAdd(-dynamicElements)
	capacity := *vm.spAdd(-dynamicElements - 2)
	baseSet := *vm.spAdd(-dynamicElements - 1)
	var newSet *value.HashSet

	var additionalCapacity int

	if !capacity.IsUndefined() {
		c, ok := value.ToGoInt(capacity)
		if c == -1 && !ok {
			return value.Ref(value.NewTooLargeCapacityError(capacity.Inspect()))
		}
		if c < 0 {
			return value.Ref(value.NewNegativeCapacityError(capacity.Inspect()))
		}
		if !ok {
			return value.Ref(value.NewCapacityTypeError(capacity.Inspect()))
		}
		additionalCapacity = c
	}

	if baseSet.IsUndefined() {
		newSet = value.NewHashSet(dynamicElements + additionalCapacity)
	} else {
		switch m := baseSet.SafeAsReference().(type) {
		case *value.HashSet:
			newSet = value.NewHashSet(m.Capacity() + additionalCapacity)
			err := HashSetCopy(vm, newSet, m)
			if !err.IsUndefined() {
				return err
			}
		default:
			panic(fmt.Sprintf("invalid hash set base: %s", baseSet.Inspect()))
		}
	}

	for i := range dynamicElements {
		val := *vm.stackAdd(firstElement, i)
		err := HashSetAppendWithMaxLoad(vm, newSet, val, 1)
		if !err.IsUndefined() {
			return err
		}
	}
	vm.popN(dynamicElements + 2)

	vm.push(value.Ref(newSet))
	return value.Undefined
}

// Create a new hashmap.
func (vm *VM) opNewHashMap(dynamicElements int) value.Value {
	firstElementOffset := -(dynamicElements * 2)
	firstElement := vm.spAdd(firstElementOffset)
	capacity := *vm.spAdd(firstElementOffset - 2)
	baseMap := *vm.spAdd(firstElementOffset - 1)
	var newMap *value.HashMap

	var additionalCapacity int

	if !capacity.IsUndefined() {
		c, ok := value.ToGoInt(capacity)
		if c == -1 && !ok {
			return value.Ref(value.NewTooLargeCapacityError(capacity.Inspect()))
		}
		if c < 0 {
			return value.Ref(value.NewNegativeCapacityError(capacity.Inspect()))
		}
		if !ok {
			return value.Ref(value.NewCapacityTypeError(capacity.Inspect()))
		}
		additionalCapacity = c
	}

	if baseMap.IsUndefined() {
		newMap = value.NewHashMap(dynamicElements + additionalCapacity)
	} else {
		switch m := baseMap.SafeAsReference().(type) {
		case *value.HashMap:
			newMap = value.NewHashMap(m.Capacity() + additionalCapacity)
			err := HashMapCopy(vm, newMap, m)
			if !err.IsUndefined() {
				return err
			}
		default:
			panic(fmt.Sprintf("invalid hash map base: %s", baseMap.Inspect()))
		}
	}

	for i := 0; i < dynamicElements*2; i += 2 {
		key := *vm.stackAdd(firstElement, i)
		val := *vm.stackAdd(firstElement, i+1)
		err := HashMapSetWithMaxLoad(vm, newMap, key, val, 1)
		if !err.IsUndefined() {
			return err
		}
	}
	vm.popN((dynamicElements * 2) + 2)

	vm.push(value.Ref(newMap))
	return value.Undefined
}

// Create a new hash record.
func (vm *VM) opNewHashRecord(dynamicElements int) value.Value {
	firstElementOffset := -(dynamicElements * 2)
	firstElement := vm.spAdd(firstElementOffset)
	baseMap := *vm.spAdd(firstElementOffset - 1)
	var newRecord *value.HashRecord

	if baseMap.IsUndefined() {
		newRecord = value.NewHashRecord(dynamicElements)
	} else {
		switch m := baseMap.SafeAsReference().(type) {
		case *value.HashRecord:
			newRecord = value.NewHashRecord(m.Length())
			err := HashRecordCopy(vm, newRecord, m)
			if !err.IsUndefined() {
				return err
			}
		default:
			panic(fmt.Sprintf("invalid hash record base: %s", baseMap.Inspect()))
		}
	}

	for i := 0; i < dynamicElements*2; i += 2 {
		key := *vm.stackAdd(firstElement, i)
		val := *vm.stackAdd(firstElement, i+1)
		HashRecordSetWithMaxLoad(vm, newRecord, key, val, 1)
	}
	vm.popN((dynamicElements * 2) + 1)

	vm.push(value.Ref(newRecord))
	return value.Undefined
}

// Create a new range.
func (vm *VM) opNewRange() {
	flag := vm.readByte()
	var newRange value.Value

	switch flag {
	case bytecode.CLOSED_RANGE_FLAG:
		to := vm.popGet()
		from := vm.popGet()
		newRange = value.Ref(value.NewClosedRange(from, to))
	case bytecode.OPEN_RANGE_FLAG:
		to := vm.popGet()
		from := vm.popGet()
		newRange = value.Ref(value.NewOpenRange(from, to))
	case bytecode.LEFT_OPEN_RANGE_FLAG:
		to := vm.popGet()
		from := vm.popGet()
		newRange = value.Ref(value.NewLeftOpenRange(from, to))
	case bytecode.RIGHT_OPEN_RANGE_FLAG:
		to := vm.popGet()
		from := vm.popGet()
		newRange = value.Ref(value.NewRightOpenRange(from, to))
	case bytecode.BEGINLESS_CLOSED_RANGE_FLAG:
		newRange = value.Ref(value.NewBeginlessClosedRange(vm.popGet()))
	case bytecode.BEGINLESS_OPEN_RANGE_FLAG:
		newRange = value.Ref(value.NewBeginlessOpenRange(vm.popGet()))
	case bytecode.ENDLESS_CLOSED_RANGE_FLAG:
		newRange = value.Ref(value.NewEndlessClosedRange(vm.popGet()))
	case bytecode.ENDLESS_OPEN_RANGE_FLAG:
		newRange = value.Ref(value.NewEndlessOpenRange(vm.popGet()))
	default:
		panic(fmt.Sprintf("invalid range flag: %#v", flag))
	}

	vm.push(newRange)
}

// Create a new array list.
func (vm *VM) opNewArrayList(dynamicElements int) value.Value {
	firstElement := vm.spAdd(-dynamicElements)
	capacity := *vm.spAdd(-dynamicElements - 2)
	baseList := *vm.spAdd(-dynamicElements - 1)
	var newArrayList value.ArrayList

	var additionalCapacity int

	if !capacity.IsUndefined() {
		c, ok := value.ToGoInt(capacity)
		if c == -1 && !ok {
			return value.Ref(value.NewTooLargeCapacityError(capacity.Inspect()))
		}
		if c < 0 {
			return value.Ref(value.NewNegativeCapacityError(capacity.Inspect()))
		}
		if !ok {
			return value.Ref(value.NewCapacityTypeError(capacity.Inspect()))
		}
		additionalCapacity = c
	}

	if baseList.IsUndefined() {
		newArrayList = make(value.ArrayList, 0, dynamicElements+additionalCapacity)
	} else {
		switch l := baseList.SafeAsReference().(type) {
		case *value.ArrayList:
			newArrayList = make(value.ArrayList, 0, cap(*l)+additionalCapacity)
			newArrayList = append(newArrayList, *l...)
		default:
			panic(fmt.Sprintf("invalid array list base: %s", baseList.Inspect()))
		}
	}

	newArrayList = append(newArrayList, unsafe.Slice(firstElement, dynamicElements)...)
	vm.popN(dynamicElements + 2)

	vm.push(value.Ref(&newArrayList))
	return value.Undefined
}

// Create a new arrayTuple.
func (vm *VM) opNewArrayTuple(dynamicElements int) {
	firstElement := vm.spAdd(-dynamicElements)
	baseArrayTuple := *vm.spAdd(-dynamicElements - 1)
	var newArrayTuple value.ArrayTuple

	if baseArrayTuple.IsUndefined() {
		newArrayTuple = make(value.ArrayTuple, 0, dynamicElements)
	} else {
		switch t := baseArrayTuple.SafeAsReference().(type) {
		case *value.ArrayTuple:
			newArrayTuple = make(value.ArrayTuple, 0, len(*t)+dynamicElements)
			newArrayTuple = append(newArrayTuple, *t...)
		default:
			panic(fmt.Sprintf("invalid array tuple base: %s", baseArrayTuple.Inspect()))
		}
	}

	newArrayTuple = append(newArrayTuple, unsafe.Slice(firstElement, dynamicElements)...)
	vm.popN(dynamicElements + 1)

	vm.push(value.Ref(&newArrayTuple))
}

// Define a new constant
func (vm *VM) opDefConst() {
	constVal := vm.popGet()
	constName := vm.popGet().AsInlineSymbol()
	namespace := vm.popGet()

	value.AddConstant(namespace, constName, constVal)
}

// Register slots for local variables and values.
func (vm *VM) opPrepLocals(count uintptr) {
	vm.spIncrementBy(count)
	vm.localCount += int(count)
}

// Push an element on top of the value stack.
func (vm *VM) push(val value.Value) {
	*vm.spGet() = val
	vm.spIncrement()
}

// Push an element on top of the value stack.
func (vm *VM) swap() {
	firstPtr := vm.spAdd(-2)
	secondPtr := vm.spAdd(-1)
	tmp := *firstPtr
	*firstPtr = *secondPtr
	*secondPtr = tmp
}

// Pop an element off the value stack.
func (vm *VM) pop() {
	vm.spDecrementBy(1)
	*vm.spGet() = value.Undefined
}

// Pop an element off the value stack.
func (vm *VM) popGet() value.Value {
	vm.spDecrementBy(1)
	val := *vm.spGet()
	*vm.spGet() = value.Undefined
	return val
}

// Pop all values on the stack leaving only slots for locals
func (vm *VM) popAll() {
	vm.popN(vm.spOffset() - vm.localCount - 1)
}

// Pop n elements off the value stack.
func (vm *VM) popN(n int) {
	spOffset := vm.spOffset()
	for i := spOffset - 1; i >= spOffset; i-- {
		vm.stack[i] = value.Undefined
	}
	vm.spDecrementBy(uintptr(n))
}

// Pop one element off the value stack skipping the first one.
func (vm *VM) popSkipOne() {
	vm.spDecrementBy(1)
	*vm.spAdd(-1) = *vm.spGet()
}

// Pop n elements off the value stack skipping the first one.
func (vm *VM) popNSkipOne(n int) {
	*vm.spAdd(-n - 1) = *vm.spAdd(-1)
	for i := vm.spOffset() - 1; i >= vm.spOffset()-n; i-- {
		*vm.spAdd(i) = value.Undefined
	}
	vm.spDecrementBy(uintptr(n))
}

// Replaces the value on top of the stack without popping it.
func (vm *VM) replace(val value.Value) {
	*vm.spAdd(-1) = val
}

// Return the element on top of the stack
// without popping it.
func (vm *VM) peek() value.Value {
	return *vm.spAdd(-1)
}

// Return the nth element on top of the stack
// without popping it.
func (vm *VM) peekAt(n int) value.Value {
	return *vm.spAdd(-1 - n)
}

type unaryOperationFunc func(val value.Value) value.Value

func (vm *VM) unaryOperation(fn unaryOperationFunc, methodName value.Symbol) value.Value {
	operand := vm.peek()
	result := fn(operand)
	if !result.IsUndefined() {
		vm.replace(result)
		return value.Undefined
	}

	er := vm.callMethodOnStackByName(methodName, 0)
	if !er.IsUndefined() {
		return er
	}
	return value.Undefined
}

// Negate the Int on top of the stack
func (vm *VM) opNegateInt() {
	operand := vm.peek()
	var result value.Value
	if operand.IsSmallInt() {
		operand := operand.AsSmallInt()
		result = operand.NegateVal()
	} else {
		operand := operand.AsReference().(*value.BigInt)
		result = value.Ref(operand.Negate())
	}
	vm.replace(result)
}

// Negate the Float on top of the stack
func (vm *VM) opNegateFloat() {
	operand := vm.peek()
	o := operand.AsFloat()
	result := (-o).ToValue()
	vm.replace(result)
}

func (vm *VM) opIncrementInt() {
	operand := vm.peek()
	var result value.Value
	if operand.IsSmallInt() {
		operand := operand.AsSmallInt()
		result = operand.Increment()
	} else {
		operand := operand.AsReference().(*value.BigInt)
		result = value.Ref(operand.Increment())
	}
	vm.replace(result)
}

func (vm *VM) opDecrementInt() {
	operand := vm.peek()
	var result value.Value
	if operand.IsSmallInt() {
		operand := operand.AsSmallInt()
		result = operand.Decrement()
	} else {
		operand := operand.AsReference().(*value.BigInt)
		result = operand.DecrementVal()
	}
	vm.replace(result)
}

// Increment the element on top of the stack
func (vm *VM) opIncrement() (err value.Value) {
	return vm.unaryOperation(value.IncrementVal, symbol.OpIncrement)
}

// Decrement the element on top of the stack
func (vm *VM) opDecrement() (err value.Value) {
	return vm.unaryOperation(value.DecrementVal, symbol.OpDecrement)
}

// Negate the element on top of the stack
func (vm *VM) opNegate() (err value.Value) {
	return vm.unaryOperation(value.NegateVal, symbol.OpNegate)
}

// Perform unary plus on the element on top of the stack
func (vm *VM) opUnaryPlus() (err value.Value) {
	return vm.unaryOperation(value.UnaryPlusVal, symbol.OpUnaryPlus)
}

// Preform bitwise not on the element on top of the stack
func (vm *VM) opBitwiseNot() (err value.Value) {
	return vm.unaryOperation(value.BitwiseNotVal, symbol.OpBitwiseNot)
}

func (vm *VM) opAppendAt() value.Value {
	val := vm.popGet()
	key := vm.popGet()
	collection := vm.peek()

	i, ok := value.ToGoInt(key)

	switch c := collection.SafeAsReference().(type) {
	case *value.ArrayTuple:
		l := len(*c)
		if !ok {
			if i == -1 {
				return value.Ref(value.NewIndexOutOfRangeError(key.Inspect(), l))
			}
			return value.Ref(value.NewCoerceError(value.IntClass, key.Class()))
		}

		if i < 0 {
			return value.Ref(value.NewNegativeIndicesInCollectionLiteralsError(fmt.Sprint(i)))
		}

		if i >= l {
			newElementsCount := (i + 1) - l
			c.Expand(newElementsCount)
		}

		(*c)[i] = val
	case *value.ArrayList:
		l := len(*c)
		if !ok {
			if i == -1 {
				return value.Ref(value.NewIndexOutOfRangeError(key.Inspect(), l))
			}
			return value.Ref(value.NewCoerceError(value.IntClass, key.Class()))
		}

		if i < 0 {
			return value.Ref(value.NewNegativeIndicesInCollectionLiteralsError(fmt.Sprint(i)))
		}

		if i >= l {
			newElementsCount := (i + 1) - l
			c.Expand(newElementsCount)
		}

		(*c)[i] = val
	default:
		panic(fmt.Sprintf("cannot APPEND_AT to: %s", collection.Inspect()))
	}

	return value.Undefined
}

func (vm *VM) opSubscriptSet() value.Value {
	val := vm.popGet()
	key := vm.popGet()
	collection := vm.peek()

	err := SubscriptSetBuiltin(vm, collection, key, val)
	if !err.IsUndefined() {
		return err
	}
	vm.replace(val)
	return value.Undefined
}

func (vm *VM) opIsA() (err value.Value) {
	classVal := vm.popGet()
	val := vm.peek()

	switch class := classVal.SafeAsReference().(type) {
	case *value.Class:
		vm.replace(value.ToElkBool(value.IsA(val, class)))
	default:
		vm.pop()
		return value.Ref(value.NewIsNotClassOrMixinError(class.Inspect()))
	}

	return value.Undefined
}

func (vm *VM) opInstanceOf() (err value.Value) {
	classVal := vm.popGet()
	val := vm.peek()

	class, ok := classVal.SafeAsReference().(*value.Class)
	if !ok || class.IsMixin() || class.IsMixinProxy() {
		vm.pop()
		return value.Ref(value.NewIsNotClassError(classVal.Inspect()))
	}

	vm.replace(value.ToElkBool(value.InstanceOf(val, class)))
	return value.Undefined
}

type binaryOperationWithoutErrFunc func(left value.Value, right value.Value) value.Value

func (vm *VM) binaryOperationWithoutErr(fn binaryOperationWithoutErrFunc, methodName value.Symbol) (err value.Value) {
	right := vm.peek()
	left := vm.peekAt(1)

	result := fn(left, right)
	if !result.IsUndefined() {
		vm.pop()
		vm.replace(result)
		return value.Undefined
	}

	er := vm.callMethodOnStackByName(methodName, 1)
	if !er.IsUndefined() {
		return er
	}

	return value.Undefined
}

func (vm *VM) negatedBinaryOperationWithoutErr(fn binaryOperationWithoutErrFunc, methodName value.Symbol) (err value.Value) {
	right := vm.peek()
	left := vm.peekAt(1)

	result := fn(left, right)
	if !result.IsUndefined() {
		vm.pop()
		vm.replace(result)
		return value.Undefined
	}

	er := vm.callMethodOnStackByName(methodName, 1)
	if !er.IsUndefined() {
		return er
	}
	vm.replace(value.ToNotBool(vm.peek()))

	return value.Undefined
}

type binaryOperationFunc func(left value.Value, right value.Value) (result value.Value, err value.Value)

func (vm *VM) binaryOperation(fn binaryOperationFunc, methodName value.Symbol) value.Value {
	right := vm.peek()
	left := vm.peekAt(1)

	result, err := fn(left, right)
	if !err.IsUndefined() {
		return err
	}
	if !result.IsUndefined() {
		vm.pop()
		vm.replace(result)
		return value.Undefined
	}

	er := vm.callMethodOnStackByName(methodName, 1)
	if !er.IsUndefined() {
		return er
	}
	return value.Undefined
}

// Perform a bitwise AND and push the result to the stack.
func (vm *VM) opBitwiseAnd() (err value.Value) {
	return vm.binaryOperation(value.BitwiseAndVal, symbol.OpAnd)
}

// Perform a bitwise AND NOT and push the result to the stack.
func (vm *VM) opBitwiseAndNot() (err value.Value) {
	return vm.binaryOperation(value.BitwiseAndNotVal, symbol.OpAndNot)
}

// Get the value under the given key and push the result to the stack.
func (vm *VM) opSubscript() (err value.Value) {
	right := vm.peek()
	left := vm.peekAt(1)

	result, err := SubscriptBuiltin(vm, left, right)
	if !err.IsUndefined() {
		return err
	}
	vm.pop()
	vm.replace(result)
	return value.Undefined
}

// Perform a bitwise OR and push the result to the stack.
func (vm *VM) opBitwiseOr() (err value.Value) {
	return vm.binaryOperation(value.BitwiseOrVal, symbol.OpOr)
}

// Perform a bitwise XOR and push the result to the stack.
func (vm *VM) opBitwiseXor() (err value.Value) {
	return vm.binaryOperation(value.BitwiseXorVal, symbol.OpXor)
}

// Perform a comparison and push the result to the stack.
func (vm *VM) opCompare() (err value.Value) {
	return vm.binaryOperation(value.CompareVal, symbol.OpSpaceship)
}

// Perform opModulo and push the result to the stack.
func (vm *VM) opModulo() (err value.Value) {
	return vm.binaryOperation(value.ModuloVal, symbol.OpModulo)
}

// Check whether two top elements on the stack are opEqual and push the result to the stack.
func (vm *VM) opEqual() (err value.Value) {
	return vm.callEqualityOperator(value.EqualVal, symbol.OpEqual)
}

func (vm *VM) callEqualityOperator(fn binaryOperationWithoutErrFunc, methodName value.Symbol) (err value.Value) {
	right := vm.peek()
	left := vm.peekAt(1)

	result := fn(left, right)
	if !result.IsUndefined() {
		vm.pop()
		vm.replace(result)
		return value.Undefined
	}

	self := *vm.spAdd(-2)
	class := self.DirectClass()
	method := class.LookupMethod(methodName)
	if method == nil {
		vm.push(value.ToElkBool(left == right))
		return value.Undefined
	}

	return vm.callMethodOnStack(method, 1)
}

func (vm *VM) callNegatedEqualityOperator(fn binaryOperationWithoutErrFunc, methodName value.Symbol) (err value.Value) {
	right := vm.peek()
	left := vm.peekAt(1)

	result := fn(left, right)
	if !result.IsUndefined() {
		vm.pop()
		vm.replace(result)
		return value.Undefined
	}

	self := *vm.spAdd(-2)
	class := self.DirectClass()
	method := class.LookupMethod(methodName)
	if method == nil {
		vm.push(value.ToElkBool(left != right))
		return value.Undefined
	}

	err = vm.callMethodOnStack(method, 1)
	if !err.IsUndefined() {
		return err
	}

	vm.replace(value.ToNotBool(vm.peek()))
	return value.Undefined
}

// Check whether two top elements on the stack are not and equal push the result to the stack.
func (vm *VM) opNotEqual() (err value.Value) {
	return vm.callNegatedEqualityOperator(value.NotEqualVal, symbol.OpEqual)
}

// Check whether two top elements on the stack are equal and push the result to the stack.
func (vm *VM) opLaxEqual() (err value.Value) {
	return vm.callEqualityOperator(value.LaxEqualVal, symbol.OpLaxEqual)
}

// Check whether two top elements on the stack are not and equal push the result to the stack.
func (vm *VM) opLaxNotEqual() (err value.Value) {
	return vm.callNegatedEqualityOperator(value.LaxNotEqualVal, symbol.OpLaxEqual)
}

// Check whether two top elements on the stack are strictly equal push the result to the stack.
func (vm *VM) opStrictEqual() {
	right := vm.popGet()
	left := vm.peek()

	result := value.StrictEqualVal(left, right)
	vm.replace(result)
}

// Check whether two top elements on the stack are strictly not equal push the result to the stack.
func (vm *VM) opStrictNotEqual() {
	right := vm.popGet()
	left := vm.peek()

	result := value.StrictNotEqualVal(left, right)
	vm.replace(result)
}

// Check whether the first operand is greater than the second and push the result to the stack.
func (vm *VM) opGreaterThan() (err value.Value) {
	return vm.binaryOperation(value.GreaterThanVal, symbol.OpGreaterThan)
}

// Check whether the first operand is greater than or equal to the second and push the result to the stack.
func (vm *VM) opGreaterThanEqual() (err value.Value) {
	return vm.binaryOperation(value.GreaterThanEqualVal, symbol.OpGreaterThanEqual)
}

// Check whether the first operand is less than the second and push the result to the stack.
func (vm *VM) opLessThan() (err value.Value) {
	return vm.binaryOperation(value.LessThanVal, symbol.OpLessThan)
}

// Check whether the first operand is less than or equal to the second and push the result to the stack.
func (vm *VM) opLessThanEqual() (err value.Value) {
	return vm.binaryOperation(value.LessThanEqualVal, symbol.OpLessThanEqual)
}

// Check whether the first operand is less than or equal to the second and push the result to the stack.
func (vm *VM) opLessThanEqualInt() {
	right := vm.popGet()
	left := vm.peek()

	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.LessThanEqualVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.LessThanEqualVal(right)
	}
	vm.replace(result)
}

func (vm *VM) opLessThanEqualFloat() {
	right := vm.popGet()
	left := vm.peek()

	l := left.AsSmallInt()
	result, _ := l.LessThanEqualVal(right)
	vm.replace(result)
}

func (vm *VM) opLessThanInt() {
	right := vm.popGet()
	left := vm.peek()

	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.LessThanVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.LessThanVal(right)
	}
	vm.replace(result)
}

func (vm *VM) opLessThanFloat() {
	right := vm.popGet()
	left := vm.peek()

	l := left.AsSmallInt()
	result, _ := l.LessThanVal(right)
	vm.replace(result)
}

func (vm *VM) opGreaterThanInt() {
	right := vm.popGet()
	left := vm.peek()

	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.GreaterThanVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.GreaterThanVal(right)
	}
	vm.replace(result)
}

func (vm *VM) opGreaterThanFloat() {
	right := vm.popGet()
	left := vm.peek()

	l := left.AsFloat()
	result, _ := l.GreaterThanVal(right)
	vm.replace(result)
}

func (vm *VM) opGreaterThanEqualInt() {
	right := vm.popGet()
	left := vm.peek()

	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.GreaterThanEqualVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.GreaterThanEqualVal(right)
	}
	vm.replace(result)
}

func (vm *VM) opGreaterThanEqualFloat() {
	right := vm.popGet()
	left := vm.peek()

	l := left.AsFloat()
	result, _ := l.GreaterThanEqualVal(right)
	vm.replace(result)
}

// Perform a left bitshift and push the result to the stack.
func (vm *VM) opLeftBitshift() (err value.Value) {
	return vm.binaryOperation(value.LeftBitshiftVal, symbol.OpLeftBitshift)
}

// Perform a logical left bitshift and push the result to the stack.
func (vm *VM) opLogicalLeftBitshift() (err value.Value) {
	return vm.binaryOperation(value.LogicalLeftBitshiftVal, symbol.OpLogicalLeftBitshift)
}

// Perform a right bitshift and push the result to the stack.
func (vm *VM) opRightBitshift() (err value.Value) {
	return vm.binaryOperation(value.RightBitshiftVal, symbol.OpRightBitshift)
}

// Perform a logical right bitshift and push the result to the stack.
func (vm *VM) opLogicalRightBitshift() (err value.Value) {
	return vm.binaryOperation(value.LogicalRightBitshiftVal, symbol.OpLogicalRightBitshift)
}

// Add two operands together and push the result to the stack.
func (vm *VM) opAdd() (err value.Value) {
	return vm.binaryOperation(value.AddVal, symbol.OpAdd)
}

func (vm *VM) opBitwiseOrInt() {
	right := vm.popGet()
	left := vm.peek()
	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.BitwiseOrVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.BitwiseOrVal(right)
	}

	vm.replace(result)
}

func (vm *VM) opBitwiseXorInt() {
	right := vm.popGet()
	left := vm.peek()
	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.BitwiseXorVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.BitwiseXorVal(right)
	}

	vm.replace(result)
}

func (vm *VM) opBitwiseAndInt() {
	right := vm.popGet()
	left := vm.peek()
	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.BitwiseAndVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.BitwiseAndVal(right)
	}

	vm.replace(result)
}

func (vm *VM) opLeftBitshiftInt() {
	right := vm.popGet()
	left := vm.peek()
	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.LeftBitshiftVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.LeftBitshiftVal(right)
	}

	vm.replace(result)
}

func (vm *VM) opRightBitshiftInt() {
	right := vm.popGet()
	left := vm.peek()
	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.RightBitshiftVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.RightBitshiftVal(right)
	}

	vm.replace(result)
}

// Add an Int to another value and push the result to the stack.
func (vm *VM) opAddInt() {
	right := vm.popGet()
	left := vm.peek()
	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.AddVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.AddVal(right)
	}

	vm.replace(result)
}

// Add a Float to another value and push the result to the stack.
func (vm *VM) opAddFloat() {
	right := vm.popGet()
	left := vm.peek()
	l := left.AsFloat()
	result, _ := l.AddVal(right)
	vm.replace(result)
}

// Subtract two operands and push the result to the stack.
func (vm *VM) opSubtract() (err value.Value) {
	return vm.binaryOperation(value.SubtractVal, symbol.OpSubtract)
}

// Subtract a value from an Int another value and push the result to the stack.
func (vm *VM) opSubtractInt() {
	right := vm.popGet()
	left := vm.peek()

	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.SubtractVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.SubtractVal(right)
	}
	vm.replace(result)
}

// Subtract a value from a Float another value and push the result to the stack.
func (vm *VM) opSubtractFloat() {
	right := vm.popGet()
	left := vm.peek()
	l := left.AsSmallInt()
	result, _ := l.SubtractVal(right)
	vm.replace(result)
}

// Multiply an Int by another value and push the result to the stack.
func (vm *VM) opMultiplyInt() {
	right := vm.popGet()
	left := vm.peek()

	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.MultiplyVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.MultiplyVal(right)
	}
	vm.replace(result)
}

// Multiply a Float by another value and push the result to the stack.
func (vm *VM) opMultiplyFloat() {
	right := vm.popGet()
	left := vm.peek()
	l := left.AsFloat()
	result, _ := l.MultiplyVal(right)
	vm.replace(result)
}

// Divide an Int by another value and push the result to the stack.
func (vm *VM) opDivideInt() {
	right := vm.popGet()
	left := vm.peek()

	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.DivideVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.DivideVal(right)
	}
	vm.replace(result)
}

// Divide a Float by another value and push the result to the stack.
func (vm *VM) opDivideFloat() {
	right := vm.popGet()
	left := vm.peek()
	l := left.AsFloat()
	result, _ := l.DivideVal(right)
	vm.replace(result)
}

// Exponentiate an Int by another value and push the result to the stack.
func (vm *VM) opExponentiateInt() {
	right := vm.popGet()
	left := vm.peek()

	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.ExponentiateVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.ExponentiateVal(right)
	}
	vm.replace(result)
}

func (vm *VM) opModuloInt() {
	right := vm.popGet()
	left := vm.peek()

	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result, _ = left.ModuloVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result, _ = leftBig.ModuloVal(right)
	}
	vm.replace(result)
}

func (vm *VM) opModuloFloat() {
	right := vm.popGet()
	left := vm.peek()
	l := left.AsFloat()
	result, _ := l.ModuloVal(right)
	vm.replace(result)
}

func (vm *VM) opEqualInt() {
	right := vm.popGet()
	left := vm.peek()

	var result value.Value
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result = left.EqualVal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result = leftBig.EqualVal(right)
	}
	vm.replace(result)
}

func (vm *VM) opEqualFloat() {
	right := vm.popGet()
	left := vm.peek()
	l := left.AsFloat()
	result := l.EqualVal(right)
	vm.replace(result)
}

func (vm *VM) opNotEqualInt() {
	right := vm.popGet()
	left := vm.peek()

	var result bool
	if left.IsSmallInt() {
		left := left.AsSmallInt()
		result = left.Equal(right)
	} else {
		leftBig := left.AsReference().(*value.BigInt)
		result = leftBig.Equal(right)
	}
	vm.replace(value.ToElkBool(!result))
}

func (vm *VM) opNotEqualFloat() {
	right := vm.popGet()
	left := vm.peek()

	var result bool
	l := left.AsFloat()
	result = l.Equal(right)
	vm.replace(value.ToElkBool(!result))
}

// Multiply two operands together and push the result to the stack.
func (vm *VM) opMultiply() (err value.Value) {
	return vm.binaryOperation(value.MultiplyVal, symbol.OpMultiply)
}

// Divide two operands and push the result to the stack.
func (vm *VM) opDivide() (err value.Value) {
	return vm.binaryOperation(value.DivideVal, symbol.OpDivide)
}

// Exponentiate two operands and push the result to the stack.
func (vm *VM) opExponentiate() (err value.Value) {
	return vm.binaryOperation(value.ExponentiateVal, symbol.OpExponentiate)
}

// Throw an error when the value on top of the stack is `nil`
func (vm *VM) opMust() {
	val := vm.peek()
	if value.IsNil(val) {
		vm.throw(value.Ref(value.NewUnexpectedNilError()))
	}
}

// Throw an error when the second value on the stack is not an instance of the class/mixin on top of the stack
func (vm *VM) opAs() {
	class := vm.popGet().AsReference().(*value.Class)
	val := vm.peek()
	if !value.IsA(val, class) {
		vm.throw(
			value.Ref(value.Errorf(
				value.TypeErrorClass,
				"failed type cast, `%s` is not an instance of `%s`",
				val.Inspect(),
				class.Name,
			)),
		)
	}
}

// Throw an error and attempt to find code
// that catches it.
func (vm *VM) throw(err value.Value) {
	vm.rethrow(err, vm.BuildStackTrace())
}

func (vm *VM) rethrow(err value.Value, stackTrace *value.StackTrace) {
	for {
		var foundCatch *CatchEntry

		ipIndex := vm.ipOffset()
		for _, catchEntry := range vm.bytecode.CatchEntries {
			if !catchEntry.Finally && ipIndex > catchEntry.From && ipIndex <= catchEntry.To {
				foundCatch = catchEntry
				break
			}
		}

		if foundCatch != nil {
			vm.ipSetOffset(foundCatch.JumpAddress)
			vm.push(value.Ref(stackTrace))
			vm.push(err)
			return
		}

		if vm.cfp == uintptr(unsafe.Pointer(&vm.callFrames[0])) || vm.lastCallFrame().stopVM {
			vm.state = errorState
			vm.errStackTrace = stackTrace
			vm.push(err)
			panic(stopVM{})
		}

		vm.restoreLastFrame()
	}
}

func (vm *VM) throwNoCatch(err value.Value) {
	vm.rethrowNoCatch(err, vm.BuildStackTrace())
}

func (vm *VM) rethrowNoCatch(err value.Value, stackTrace *value.StackTrace) {
	vm.state = errorState
	vm.errStackTrace = stackTrace
	vm.push(err)
	vm.restoreLastFrame()
}

// Used in a panic to stop the VM
type stopVM struct{}

func (vm *VM) jumpToFinallyForReturn() bool {
	catchEntry := vm.findFinallyCatchEntry()
	if catchEntry == nil {
		return false
	}

	// execute finally
	vm.ipSetOffset(catchEntry.JumpAddress)
	return true
}

func (vm *VM) jumpToFinallyForBreakOrContinue() bool {
	catchEntry := vm.findFinallyCatchEntry()
	if catchEntry == nil {
		return false
	}

	// execute finally
	vm.ipSetOffset(catchEntry.JumpAddress + 4) // skip NIL, JUMP, offsetByte1, offsetByte2
	return true
}

func (vm *VM) findFinallyCatchEntry() *CatchEntry {
	ipIndex := vm.ipOffset()
	for _, catchEntry := range vm.bytecode.CatchEntries {
		if catchEntry.Finally && ipIndex > catchEntry.From && ipIndex <= catchEntry.To {
			return catchEntry
		}
	}

	return nil
}
