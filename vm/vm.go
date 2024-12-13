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
	"github.com/elk-language/elk/config"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/fatih/color"
)

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
	normalMode             mode = iota
	singleFunctionCallMode      // the VM should halt after executing a single method
	errorMode                   // the VM stopped after encountering an uncaught error
)

// A single instance of the Elk Virtual Machine.
type VM struct {
	bytecode        *BytecodeFunction
	upvalues        []*Upvalue
	openUpvalueHead *Upvalue      // linked list of open upvalues, living on the stack
	ip              uintptr       // Instruction pointer -- points to the next bytecode instruction
	sp              *value.Value  // Stack pointer -- points to the offset where the next element will be pushed to
	fp              *value.Value  // Frame pointer -- points to the offset where the section of the stack for current frame starts
	localCount      int           // the amount of registered locals
	cfp             uintptr       // Call frame pointer
	stack           []value.Value // Value stack
	callFrames      []CallFrame   // Call stack
	errStackTrace   string        // The most recent error stack trace
	Stdin           io.Reader     // standard output used by the VM
	Stdout          io.Writer     // standard input used by the VM
	Stderr          io.Writer     // standard error used by the VM
	mode            mode
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
	stack := make([]value.Value, VALUE_STACK_SIZE)
	callFrames := make([]CallFrame, CALL_STACK_SIZE)
	vm := &VM{
		stack:      stack,
		sp:         &stack[0],
		fp:         &stack[0],
		callFrames: callFrames,
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
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
	vm.run()
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
	if vm.sp == &vm.stack[0] {
		// populate the predeclared local variables
		vm.push(value.Ref(value.GlobalObject)) // populate self
		vm.localCount = 1
	} else {
		// pop the return value of the last run
		vm.pop()
	}
	vm.run()

	err := vm.Err()
	if !err.IsUndefined() {
		return value.Undefined, err
	}
	return vm.peek(), value.Undefined
}

func (vm *VM) PrintError() {
	fmt.Print(vm.ErrStackTrace())
	c := color.New(color.FgRed, color.Bold)
	c.Print("Error! Uncaught thrown value:")
	fmt.Print(" ")
	fmt.Println(lexer.Colorize(vm.Err().Inspect()))
	fmt.Println()
}

// Get the stored error.
func (vm *VM) Err() value.Value {
	if vm.mode == errorMode {
		return vm.peek()
	}

	return value.Undefined
}

// Get the stored error stack trace.
func (vm *VM) ErrStackTrace() string {
	if vm.mode == errorMode {
		return vm.errStackTrace
	}

	return ""
}

// Get the value on top of the value stack.
func (vm *VM) StackTop() value.Value {
	return vm.peek()
}

func (vm *VM) Stack() []value.Value {
	spIndex := vm.spOffset()
	return vm.stack[:spIndex]
}

func (vm *VM) InspectStack() {
	fmt.Println("stack:")
	for i, value := range vm.Stack() {
		fmt.Printf("%d => %s\n", i, value.Inspect())
	}
}

func (vm *VM) throwIfErr(err value.Value) {
	if !err.IsUndefined() {
		vm.throw(err)
	}
}

var callSymbol = value.ToSymbol("call")

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

// Call an Elk closure from Go code, preserving the state of the VM.
func (vm *VM) CallClosure(closure *Closure, args ...value.Value) (value.Value, value.Value) {
	if closure.Bytecode.ParameterCount() != len(args) {
		return value.Undefined, value.Ref(value.NewWrongArgumentCountError(
			closure.Bytecode.Name().String(),
			len(args),
			closure.Bytecode.ParameterCount(),
		))
	}

	vm.createCurrentCallFrame()
	vm.bytecode = closure.Bytecode
	vm.fp = vm.sp
	vm.ipSet(&closure.Bytecode.Instructions[0])
	vm.localCount = len(args)
	vm.upvalues = closure.Upvalues
	vm.mode = singleFunctionCallMode
	// push `self`
	vm.push(closure.Self)
	for _, arg := range args {
		vm.push(arg)
	}
	vm.run()
	err := vm.Err()
	if !err.IsUndefined() {
		vm.mode = normalMode
		vm.restoreLastFrame()

		return value.Undefined, err
	}
	vm.mode = normalMode
	return vm.pop(), value.Undefined
}

// Call an Elk method from Go code, preserving the state of the VM.
func (vm *VM) CallMethodByName(name value.Symbol, args ...value.Value) (value.Value, value.Value) {
	self := args[0]
	class := self.DirectClass()
	method := class.LookupMethod(name)
	if method == nil {
		return value.Undefined, value.Ref(value.NewNoMethodError(string(name.ToString()), self))
	}
	return vm.CallMethod(method, args...)
}

func (vm *VM) CallMethod(method value.Method, args ...value.Value) (value.Value, value.Value) {
	self := args[0]
	if method.ParameterCount() != len(args)-1 {
		return value.Undefined, value.Ref(value.NewWrongArgumentCountError(
			method.Name().String(),
			len(args)-1,
			method.ParameterCount(),
		))
	}

	switch m := method.(type) {
	case *BytecodeFunction:
		vm.createCurrentCallFrame()
		vm.bytecode = m
		vm.fp = vm.sp
		vm.ipSet(&m.Instructions[0])
		vm.localCount = len(args)
		vm.mode = singleFunctionCallMode
		for _, arg := range args {
			vm.push(arg)
		}
		vm.run()
		err := vm.Err()
		if !err.IsUndefined() {
			vm.mode = normalMode
			vm.restoreLastFrame()

			return value.Undefined, err
		}
		vm.mode = normalMode
		return vm.pop(), value.Undefined
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
		vm.createCurrentCallFrame()
		vm.bytecode = m
		vm.fp = vm.spAdd(-args - 1)
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
	if method == nil {
		return value.Ref(value.NewNoMethodError(string(name.ToString()), self))
	}

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
		// BENCHMARK: replace with a jump table
		switch instruction {
		case bytecode.RETURN_FINALLY:
			if vm.jumpToFinallyForReturn() {
				continue
			}

			// return normally
			if vm.cfp == uintptr(unsafe.Pointer(&vm.callFrames[0])) {
				return
			}
			vm.returnFromFunction()
			if vm.mode == singleFunctionCallMode {
				return
			}
		case bytecode.RETURN:
			if vm.cfp == uintptr(unsafe.Pointer(&vm.callFrames[0])) {
				return
			}
			vm.returnFromFunction()
			if vm.mode == singleFunctionCallMode {
				return
			}
		case bytecode.RETURN_FIRST_ARG:
			vm.opGetLocal(1)
			if vm.cfp == uintptr(unsafe.Pointer(&vm.callFrames[0])) {
				return
			}
			vm.returnFromFunction()
			if vm.mode == singleFunctionCallMode {
				return
			}
		case bytecode.RETURN_SELF:
			vm.self()
			if vm.cfp == uintptr(unsafe.Pointer(&vm.callFrames[0])) {
				return
			}
			vm.returnFromFunction()
			if vm.mode == singleFunctionCallMode {
				return
			}
		case bytecode.CLOSURE:
			vm.closure()
		case bytecode.JUMP_TO_FINALLY:
			leftFinallyCount := vm.peek().MustSmallInt()
			jumpOffset := vm.peekAt(1).MustSmallInt()

			if leftFinallyCount > 0 {
				vm.replace((leftFinallyCount - 1).ToValue())
				if !vm.jumpToFinallyForBreakOrContinue() {
					panic("could not find a finally block to jump to in JUMP_TO_FINALLY")
				}
				continue
			}

			vm.popN(2)
			vm.ipSetOffset(int(jumpOffset))
		case bytecode.NOOP:
		case bytecode.DUP:
			vm.push(vm.peek())
		case bytecode.SWAP:
			vm.swap()
		case bytecode.DUP_N:
			n := int(vm.readByte())
			for _, element := range unsafe.Slice(vm.spAdd(-n), n) {
				vm.push(element)
			}
		case bytecode.SELF:
			vm.self()
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
		case bytecode.DEF_METHOD_ALIAS:
			vm.opDefMethodAlias()
		case bytecode.INCLUDE:
			vm.throwIfErr(vm.opInclude())
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
		case bytecode.GET_IVAR8:
			vm.throwIfErr(
				vm.opGetIvar(int(vm.readByte())),
			)
		case bytecode.GET_IVAR16:
			vm.throwIfErr(
				vm.opGetIvar(int(vm.readUint16())),
			)
		case bytecode.GET_IVAR32:
			vm.throwIfErr(
				vm.opGetIvar(int(vm.readUint32())),
			)
		case bytecode.SET_IVAR8:
			vm.throwIfErr(
				vm.opSetIvar(int(vm.readByte())),
			)
		case bytecode.SET_IVAR16:
			vm.throwIfErr(
				vm.opSetIvar(int(vm.readUint16())),
			)
		case bytecode.SET_IVAR32:
			vm.throwIfErr(
				vm.opSetIvar(int(vm.readUint32())),
			)
		case bytecode.CALL_METHOD8:
			vm.throwIfErr(
				vm.opCallMethod(int(vm.readByte())),
			)
		case bytecode.CALL_METHOD16:
			vm.throwIfErr(
				vm.opCallMethod(int(vm.readUint16())),
			)
		case bytecode.CALL_METHOD32:
			vm.throwIfErr(
				vm.opCallMethod(int(vm.readUint32())),
			)
		case bytecode.CALL8:
			vm.throwIfErr(
				vm.opCall(int(vm.readByte())),
			)
		case bytecode.CALL16:
			vm.throwIfErr(
				vm.opCall(int(vm.readUint16())),
			)
		case bytecode.CALL32:
			vm.throwIfErr(
				vm.opCall(int(vm.readUint32())),
			)
		case bytecode.CALL_SELF8:
			vm.throwIfErr(
				vm.opCallSelf(int(vm.readByte())),
			)
		case bytecode.CALL_SELF16:
			vm.throwIfErr(
				vm.opCallSelf(int(vm.readUint16())),
			)
		case bytecode.CALL_SELF32:
			vm.throwIfErr(
				vm.opCallSelf(int(vm.readUint32())),
			)
		case bytecode.INSTANCE_OF:
			vm.throwIfErr(vm.opInstanceOf())
		case bytecode.IS_A:
			vm.throwIfErr(vm.opIsA())
		case bytecode.ROOT:
			vm.push(value.Ref(value.RootModule))
		case bytecode.UNDEFINED:
			vm.push(value.Undefined)
		case bytecode.LOAD_VALUE8:
			vm.push(vm.readValue8())
		case bytecode.LOAD_VALUE16:
			vm.push(vm.readValue16())
		case bytecode.LOAD_VALUE32:
			vm.push(vm.readValue32())
		case bytecode.ADD:
			vm.throwIfErr(vm.opAdd())
		case bytecode.SUBTRACT:
			vm.throwIfErr(vm.opSubtract())
		case bytecode.MULTIPLY:
			vm.throwIfErr(vm.opMultiply())
		case bytecode.DIVIDE:
			vm.throwIfErr(vm.opDivide())
		case bytecode.EXPONENTIATE:
			vm.throwIfErr(vm.opExponentiate())
		case bytecode.NEGATE:
			vm.throwIfErr(vm.opNegate())
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
		case bytecode.POP_ALL:
			vm.popAll()
		case bytecode.POP_N:
			vm.popN(int(vm.readByte()))
		case bytecode.POP_N_SKIP_ONE:
			vm.popNSkipOne(int(vm.readByte()))
		case bytecode.POP_SKIP_ONE:
			vm.popSkipOne()
		case bytecode.INCREMENT:
			vm.throwIfErr(vm.opIncrement())
		case bytecode.DECREMENT:
			vm.throwIfErr(vm.opDecrement())
		case bytecode.GET_LOCAL8:
			vm.opGetLocal(int(vm.readByte()))
		case bytecode.GET_LOCAL16:
			vm.opGetLocal(int(vm.readUint16()))
		case bytecode.SET_LOCAL8:
			vm.opSetLocal(int(vm.readByte()))
		case bytecode.SET_LOCAL16:
			vm.opSetLocal(int(vm.readUint16()))
		case bytecode.GET_UPVALUE8:
			vm.opGetUpvalue(int(vm.readByte()))
		case bytecode.GET_UPVALUE16:
			vm.opGetUpvalue(int(vm.readUint16()))
		case bytecode.SET_UPVALUE8:
			vm.opSetUpvalue(int(vm.readByte()))
		case bytecode.SET_UPVALUE16:
			vm.opSetUpvalue(int(vm.readUint16()))
		case bytecode.CLOSE_UPVALUE8:
			last := vm.fpAdd(int(vm.readByte()))
			vm.opCloseUpvalues(last)
		case bytecode.CLOSE_UPVALUE16:
			last := vm.fpAdd(int(vm.readUint16()))
			vm.opCloseUpvalues(last)
		case bytecode.LEAVE_SCOPE16:
			vm.opLeaveScope(int(vm.readByte()), int(vm.readByte()))
		case bytecode.LEAVE_SCOPE32:
			vm.opLeaveScope(int(vm.readUint16()), int(vm.readUint16()))
		case bytecode.PREP_LOCALS8:
			vm.opPrepLocals(int(vm.readByte()))
		case bytecode.PREP_LOCALS16:
			vm.opPrepLocals(int(vm.readUint16()))
		case bytecode.SET_SUPERCLASS:
			vm.opSetSuperclass()
		case bytecode.GET_CONST8:
			vm.throwIfErr(vm.opGetConst(int(vm.readByte())))
		case bytecode.GET_CONST16:
			vm.throwIfErr(
				vm.opGetConst(int(vm.readUint16())),
			)
		case bytecode.GET_CONST32:
			vm.throwIfErr(
				vm.opGetConst(int(vm.readUint32())),
			)
		case bytecode.DEF_CONST:
			vm.opDefConst()
		case bytecode.NEW_RANGE:
			vm.opNewRange()
		case bytecode.NEW_ARRAY_TUPLE8:
			vm.opNewArrayTuple(int(vm.readByte()))
		case bytecode.NEW_ARRAY_TUPLE32:
			vm.opNewArrayTuple(int(vm.readUint32()))
		case bytecode.NEW_ARRAY_LIST8:
			vm.throwIfErr(
				vm.opNewArrayList(int(vm.readByte())),
			)
		case bytecode.NEW_ARRAY_LIST32:
			vm.throwIfErr(
				vm.opNewArrayList(int(vm.readUint32())),
			)
		case bytecode.NEW_HASH_SET8:
			vm.throwIfErr(
				vm.opNewHashSet(int(vm.readByte())),
			)
		case bytecode.NEW_HASH_SET32:
			vm.throwIfErr(
				vm.opNewHashSet(int(vm.readUint32())),
			)
		case bytecode.NEW_HASH_MAP8:
			vm.throwIfErr(
				vm.opNewHashMap(int(vm.readByte())),
			)
		case bytecode.NEW_HASH_MAP32:
			vm.throwIfErr(
				vm.opNewHashMap(int(vm.readUint32())),
			)
		case bytecode.NEW_HASH_RECORD8:
			vm.throwIfErr(
				vm.opNewHashRecord(int(vm.readByte())),
			)
		case bytecode.NEW_HASH_RECORD32:
			vm.throwIfErr(
				vm.opNewHashRecord(int(vm.readUint32())),
			)
		case bytecode.NEW_STRING8:
			vm.throwIfErr(vm.opNewString(int(vm.readByte())))
		case bytecode.NEW_STRING32:
			vm.throwIfErr(vm.opNewString(int(vm.readUint32())))
		case bytecode.NEW_SYMBOL8:
			vm.throwIfErr(vm.opNewSymbol(int(vm.readByte())))
		case bytecode.NEW_SYMBOL32:
			vm.throwIfErr(vm.opNewSymbol(int(vm.readUint32())))
		case bytecode.NEW_REGEX8:
			vm.throwIfErr(vm.opNewRegex(vm.readByte(), int(vm.readByte())))
		case bytecode.NEW_REGEX32:
			vm.throwIfErr(vm.opNewRegex(vm.readByte(), int(vm.readUint32())))
		case bytecode.FOR_IN:
			vm.throwIfErr(vm.opForIn())
		case bytecode.GET_ITERATOR:
			vm.throwIfErr(vm.opGetIterator())
		case bytecode.JUMP_UNLESS:
			if value.Falsy(vm.peek()) {
				jump := vm.readUint16()
				vm.ipIncrementBy(int(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_IF_NIL:
			if vm.peek() == value.Nil {
				jump := vm.readUint16()
				vm.ipIncrementBy(int(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP_IF:
			if value.Truthy(vm.peek()) {
				jump := vm.readUint16()
				vm.ipIncrementBy(int(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.JUMP:
			jump := vm.readUint16()
			vm.ipIncrementBy(int(jump))
		case bytecode.JUMP_UNLESS_UNDEF:
			if !vm.peek().IsUndefined() {
				jump := vm.readUint16()
				vm.ipIncrementBy(int(jump))
				break
			}
			vm.ipIncrementBy(2)
		case bytecode.LOOP:
			jump := vm.readUint16()
			vm.ipIncrementBy(-int(jump))
		case bytecode.THROW:
			vm.throw(vm.pop())
		case bytecode.MUST:
			vm.opMust()
		case bytecode.AS:
			vm.mustAs()
		case bytecode.RETHROW:
			err := vm.pop()
			stackTrace := vm.pop().MustReference().(value.String)
			vm.rethrow(err, stackTrace)
		case bytecode.LBITSHIFT:
			vm.throwIfErr(vm.opLeftBitshift())
		case bytecode.LOGIC_LBITSHIFT:
			vm.throwIfErr(vm.opLogicalLeftBitshift())
		case bytecode.RBITSHIFT:
			vm.throwIfErr(vm.opRightBitshift())
		case bytecode.LOGIC_RBITSHIFT:
			vm.throwIfErr(vm.opLogicalRightBitshift())
		case bytecode.BITWISE_AND:
			vm.throwIfErr(vm.opBitwiseAnd())
		case bytecode.BITWISE_AND_NOT:
			vm.throwIfErr(vm.opBitwiseAndNot())
		case bytecode.BITWISE_OR:
			vm.throwIfErr(vm.opBitwiseOr())
		case bytecode.BITWISE_XOR:
			vm.throwIfErr(vm.opBitwiseXor())
		case bytecode.MODULO:
			vm.throwIfErr(vm.opModulo())
		case bytecode.COMPARE:
			vm.throwIfErr(vm.opCompare())
		case bytecode.EQUAL:
			vm.throwIfErr(vm.opEqual())
		case bytecode.NOT_EQUAL:
			vm.throwIfErr(vm.opNotEqual())
		case bytecode.LAX_EQUAL:
			vm.throwIfErr(vm.opLaxEqual())
		case bytecode.LAX_NOT_EQUAL:
			vm.throwIfErr(vm.opLaxNotEqual())
		case bytecode.STRICT_EQUAL:
			vm.throwIfErr(vm.opStrictEqual())
		case bytecode.STRICT_NOT_EQUAL:
			vm.throwIfErr(vm.opStrictNotEqual())
		case bytecode.GREATER:
			vm.throwIfErr(vm.opGreaterThan())
		case bytecode.GREATER_EQUAL:
			vm.throwIfErr(vm.opGreaterThanEqual())
		case bytecode.LESS:
			vm.throwIfErr(vm.opLessThan())
		case bytecode.LESS_EQUAL:
			vm.throwIfErr(vm.opLessThanEqual())
		case bytecode.INSPECT_STACK:
			vm.InspectStack()
		default:
			panic(fmt.Sprintf("Unknown bytecode instruction: %#v", instruction))
		}
	}

}

func (vm *VM) closure() {
	function := vm.peek().MustReference().(*BytecodeFunction)
	closure := NewClosure(function, vm.selfValue())
	vm.replace(value.Ref(closure))

	for i := range len(closure.Upvalues) {
		flags := bitfield.BitField8FromInt(vm.readByte())
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
	vm.ipIncrement() // skip past the terminator
}

func (vm *VM) captureUpvalue(location *value.Value) *Upvalue {
	var prevUpvalue *Upvalue
	currentUpvalue := vm.openUpvalueHead
	for {
		if currentUpvalue == nil ||
			(uintptr)(unsafe.Pointer(currentUpvalue.location)) <=
				(uintptr)(unsafe.Pointer(location)) {
			break
		}
		prevUpvalue = currentUpvalue
		currentUpvalue = currentUpvalue.next
	}

	if currentUpvalue != nil && currentUpvalue.location == location {
		return currentUpvalue
	}

	newUpvalue := NewUpvalue(location)
	newUpvalue.next = currentUpvalue
	if prevUpvalue != nil {
		prevUpvalue.next = newUpvalue
	} else {
		vm.openUpvalueHead = newUpvalue
	}
	return newUpvalue
}

func (vm *VM) returnFromFunction() {
	returnValue := vm.pop()
	vm.restoreLastFrame()
	vm.push(returnValue)
}

// Restore the state of the VM to the last call frame.
func (vm *VM) restoreLastFrame() {
	vm.cfpIncrementBy(-1)
	cf := vm.cfpGet()

	vm.ip = cf.ip
	vm.opCloseUpvalues(vm.fp)
	vm.popN(vm.spOffsetFrom(vm.fp))
	vm.fp = cf.fp
	vm.localCount = cf.localCount
	vm.bytecode = cf.bytecode
	vm.upvalues = cf.upvalues
}

func (vm *VM) ResetError() {
	vm.mode = normalMode
	vm.errStackTrace = ""
}

func addStackTraceEntry(output io.Writer, id int, fileName string, lineNumber int, name string) {
	// "  %d: %s:%d, in `%s`\n"
	fmt.Fprint(output, " ")
	color.New(color.FgHiBlue).Fprintf(output, "%d", id)
	fmt.Fprintf(output, ": %s:%d, in ", fileName, lineNumber)
	color.New(color.FgHiYellow).Fprintf(output, "`%s`", name)
	fmt.Fprintln(output)
}

func (vm *VM) BuildStackTrace() string {
	var buffer strings.Builder
	buffer.WriteString("Stack trace (the most recent call is last)\n")

	var i int
	if vm.cfp != uintptr(unsafe.Pointer(&vm.callFrames[0])) {
		callFrame := &vm.callFrames[0]
		for {
			addStackTraceEntry(
				&buffer,
				i,
				callFrame.FileName(),
				callFrame.LineNumber(),
				callFrame.Name().String(),
			)

			i++
			callFrame = vm.callFrameAdd(callFrame, 1)
			if uintptr(unsafe.Pointer(callFrame)) == vm.cfp {
				break
			}
		}
	}
	addStackTraceEntry(
		&buffer,
		i+1,
		vm.bytecode.FileName(),
		vm.bytecode.GetLineNumber(vm.ipOffset()-1),
		vm.bytecode.Name().String(),
	)
	// Stack trace (the most recent call is last):
	//   0: /tmp/test.elk:18, in `foo`
	//   1: /tmp/test.elk:11, in `bar`

	return buffer.String()
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

// Add n to the stack pointer
func (vm *VM) spIncrementBy(n int) {
	vm.sp = vm.spAdd(n)
}

func (vm *VM) spOffsetFrom(ptr *value.Value) int {
	return int(uintptr(unsafe.Pointer(vm.sp))-uintptr(unsafe.Pointer(ptr))) / int(value.ValueSize)
}

func (vm *VM) spOffset() int {
	return vm.spOffsetFrom(&vm.stack[0])
}

func (vm *VM) spAdd(n int) *value.Value {
	return vm.stackAdd(vm.sp, n)
}

func (vm *VM) stackAdd(ptr *value.Value, n int) *value.Value {
	return (*value.Value)(unsafe.Add(unsafe.Pointer(ptr), n*int(value.ValueSize)))
}

func (vm *VM) fpOffset() int {
	return int(uintptr(unsafe.Pointer(vm.fp))-uintptr(unsafe.Pointer(&vm.stack[0]))) / int(value.ValueSize)
}

func (vm *VM) fpAdd(n int) *value.Value {
	return (*value.Value)(unsafe.Add(unsafe.Pointer(vm.fp), n*int(value.ValueSize)))
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

// Add n to the instruction pointer
func (vm *VM) ipIncrementBy(n int) {
	vm.ip = (uintptr)(unsafe.Add(unsafe.Pointer(vm.ip), n))
}

// Increment the call frame pointer
func (vm *VM) cfpIncrement() {
	vm.cfpIncrementBy(1)
}

func (vm *VM) cfpAdd(n int) *CallFrame {
	return vm.callFrameAdd(vm.cfpGet(), n)
}

func (vm *VM) cfpGet() *CallFrame {
	return (*CallFrame)(unsafe.Pointer(vm.cfp))
}

// Add n to the call frame pointer
func (vm *VM) cfpIncrementBy(n int) {
	vm.cfpSet(vm.cfpAdd(n))
}

// Set the typesafe call frame pointer
func (vm *VM) cfpSet(ptr *CallFrame) {
	vm.cfp = uintptr(unsafe.Pointer(ptr))
}

func (vm *VM) callFrameAdd(ptr *CallFrame, n int) *CallFrame {
	return (*CallFrame)(unsafe.Add(unsafe.Pointer(ptr), n*int(CallFrameSize)))
}

// Read the next byte of code
func (vm *VM) readByte() byte {
	// BENCHMARK: compare pointer arithmetic to offsets
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
	name := vm.pop().MustSymbol()
	parentNamespace := vm.pop()

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
		newNamespace = value.Ref(value.NewClassWithOptions(value.ClassWithParent(nil)))
	case bytecode.DEF_MIXIN_FLAG:
		newNamespace = value.Ref(value.NewMixin())
	case bytecode.DEF_INTERFACE_FLAG:
		newNamespace = value.Ref(value.NewInterface())
	}

	parentConstantContainer.AddConstant(name, newNamespace)
}

func (vm *VM) opGetSingleton() (err value.Value) {
	val := vm.pop()
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
	val := vm.pop()
	class := val.Class()
	vm.push(value.Ref(class))
}

func (vm *VM) selfValue() value.Value {
	return vm.getLocalValue(0)
}

// Call a method with an implicit receiver
func (vm *VM) opCallSelf(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].MustReference().(*value.CallSiteInfo)

	self := vm.selfValue()
	class := self.DirectClass()

	method := class.LookupMethod(callInfo.Name)
	if method == nil {
		return value.Ref(value.NewNoMethodError(string(callInfo.Name.ToString()), self))
	}

	// shift all arguments one slot forward to make room for self
	for i := 0; i < callInfo.ArgumentCount; i++ {
		*vm.spAdd(-i) = *vm.spAdd(-i - 1)
	}
	*vm.spAdd(-callInfo.ArgumentCount) = self
	vm.spIncrement()

	switch m := method.(type) {
	case *BytecodeFunction:
		return vm.callBytecodeFunction(m, callInfo)
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	}

	panic(fmt.Sprintf("tried to call a method that is neither bytecode nor native: %#v", method))
}

// Set the value of an instance variable
func (vm *VM) opSetIvar(nameIndex int) (err value.Value) {
	name := vm.bytecode.Values[nameIndex].MustSymbol()
	val := vm.peek()

	self := vm.selfValue()
	ivars := self.InstanceVariables()
	if ivars == nil {
		return value.Ref(value.NewCantSetInstanceVariablesOnPrimitiveError(self.Inspect()))
	}

	ivars.Set(name, val)
	return value.Undefined
}

// Get the value of an instance variable
func (vm *VM) opGetIvar(nameIndex int) (err value.Value) {
	name := vm.bytecode.Values[nameIndex].MustSymbol()

	self := vm.selfValue()
	ivars := self.InstanceVariables()
	if ivars == nil {
		return value.Ref(value.NewCantAccessInstanceVariablesOnPrimitiveError(self.Inspect()))
	}

	val := ivars.Get(name)
	if val.IsUndefined() {
		vm.push(value.Nil)
	} else {
		vm.push(val)
	}

	return value.Undefined
}

// Pop the value on top of the stack and push its opCopy.
func (vm *VM) opCopy() {
	element := vm.peek()
	vm.replace(element.Copy())
}

// Set the value under the given key in a hash-map or hash-record
func (vm *VM) opMapSet() {
	val := vm.pop()
	key := vm.pop()
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
	element := vm.pop()
	collection := vm.peek()

	if collection.IsUndefined() {
		vm.replace(value.Ref(&value.ArrayTuple{element}))
		return
	}
	switch c := collection.MustReference().(type) {
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
		return vm.callBytecodeFunction(m, callInfo)
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	case nil:
		if callInfo.ArgumentCount == 0 {
			// no initialiser defined
			// no arguments given
			// just replace the class with the instance
			return value.Undefined
		}

		return value.Ref(value.NewWrongArgumentCountError(
			callInfo.Name.String(),
			callInfo.ArgumentCount,
			0,
		))
	default:
		panic(fmt.Sprintf("tried to call an invalid initialiser method: %#v", method))
	}
}

// Call a method in a pattern.
// Return false if the receiver does not have the method
// or it throws TypeError.
func (vm *VM) callPattern(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].MustReference().(*value.CallSiteInfo)

	self := *vm.spAdd(-callInfo.ArgumentCount - 1)
	class := self.DirectClass()

	method := class.LookupMethod(callInfo.Name)
	if method == nil {
		vm.popN(callInfo.ArgumentCount + 1)
		vm.push(value.False)
		return value.Undefined
	}
	switch m := method.(type) {
	case *BytecodeFunction:
		err = vm.callBytecodeFunction(m, callInfo)
	case *NativeMethod:
		err = vm.callNativeMethod(m, callInfo)
	case *GetterMethod:
		if callInfo.ArgumentCount != 0 {
			return value.Ref(value.NewWrongArgumentCountError(
				method.Name().String(),
				callInfo.ArgumentCount,
				0,
			))
		}
		vm.pop() // pop self
		var result value.Value
		result, err = m.Call(self)
		if err.IsUndefined() {
			vm.push(result)
		}
	case *SetterMethod:
		if callInfo.ArgumentCount != 1 {
			return value.Ref(value.NewWrongArgumentCountError(
				method.Name().String(),
				callInfo.ArgumentCount,
				1,
			))
		}
		other := vm.pop()
		vm.pop() // pop self
		var result value.Value
		result, err = m.Call(self, other)
		if err.IsUndefined() {
			vm.push(result)
		}
	default:
		panic(fmt.Sprintf("tried to call an invalid method: %#v", method))
	}

	if !err.IsUndefined() {
		if err.Class() == value.TypeErrorClass {
			vm.push(value.False)
			return value.Undefined
		}
		return err
	}

	return value.Undefined
}

// Call the `opCall` method with an explicit receiver
func (vm *VM) opCall(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].MustReference().(*value.CallSiteInfo)

	self, isClosure := vm.spAdd(-callInfo.ArgumentCount - 1).MustReference().(*Closure)
	if !isClosure {
		return vm.opCallMethod(callInfoIndex)
	}

	return vm.callClosure(self, callInfo)
}

// set up the vm to execute a closure
func (vm *VM) callClosure(closure *Closure, callInfo *value.CallSiteInfo) (err value.Value) {
	function := closure.Bytecode
	if err := vm.prepareArguments(function, callInfo); !err.IsUndefined() {
		return err
	}

	vm.createCurrentCallFrame()

	vm.localCount = function.parameterCount + 1
	vm.bytecode = function
	vm.fp = vm.spAdd(-function.parameterCount - 1)
	vm.ipSet(&function.Instructions[0])
	vm.upvalues = closure.Upvalues

	return value.Undefined
}

// Call a method with an explicit receiver
func (vm *VM) opCallMethod(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].MustReference().(*value.CallSiteInfo)

	selfPtr := vm.spAdd(-callInfo.ArgumentCount - 1)
	self := *selfPtr
	class := self.DirectClass()

	method := class.LookupMethod(callInfo.Name)
	if method == nil {
		vm.popN(callInfo.ArgumentCount + 1)
		return value.Ref(value.NewNoMethodError(string(callInfo.Name.ToString()), self))
	}
	switch m := method.(type) {
	case *BytecodeFunction:
		return vm.callBytecodeFunction(m, callInfo)
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	case *GetterMethod:
		if callInfo.ArgumentCount != 0 {
			return value.Ref(value.NewWrongArgumentCountError(
				method.Name().String(),
				callInfo.ArgumentCount,
				0,
			))
		}
		vm.pop() // pop self
		result, err := m.Call(self)
		if !err.IsUndefined() {
			return err
		}
		vm.push(result)
		return value.Undefined
	case *SetterMethod:
		if callInfo.ArgumentCount != 1 {
			return value.Ref(value.NewWrongArgumentCountError(
				method.Name().String(),
				callInfo.ArgumentCount,
				1,
			))
		}
		other := vm.pop()
		vm.pop() // pop self
		result, err := m.Call(self, other)
		if !err.IsUndefined() {
			return err
		}
		vm.push(result)
		return value.Undefined
	default:
		panic(fmt.Sprintf("tried to call an invalid method: %T", method))
	}
}

// set up the vm to execute a native method
func (vm *VM) callNativeMethod(method *NativeMethod, callInfo *value.CallSiteInfo) (err value.Value) {
	if prepErr := vm.prepareArguments(method, callInfo); !prepErr.IsUndefined() {
		return prepErr
	}

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

// set up the vm to execute a bytecode method
func (vm *VM) callBytecodeFunction(method *BytecodeFunction, callInfo *value.CallSiteInfo) (err value.Value) {
	if err := vm.prepareArguments(method, callInfo); !err.IsUndefined() {
		return err
	}

	vm.createCurrentCallFrame()

	vm.localCount = method.parameterCount + 1
	vm.bytecode = method
	vm.fp = vm.spAdd(-method.ParameterCount() - 1)
	vm.ipSet(&method.Instructions[0])

	return value.Undefined
}

func (vm *VM) prepareArguments(method value.Method, callInfo *value.CallSiteInfo) (err value.Value) {
	optParamCount := method.OptionalParameterCount()
	paramCount := method.ParameterCount()
	reqParamCount := paramCount - optParamCount

	if callInfo.ArgumentCount < reqParamCount {
		return value.Ref(value.NewWrongArgumentCountRangeError(
			method.Name().String(),
			callInfo.ArgumentCount,
			reqParamCount,
			paramCount,
		))
	}

	if optParamCount > 0 {
		// populate missing optional arguments with undefined
		missingArgCount := paramCount - callInfo.ArgumentCount
		for i := 0; i < missingArgCount; i++ {
			vm.push(value.Undefined)
		}
	} else if paramCount != callInfo.ArgumentCount {
		return value.Ref(value.NewWrongArgumentCountError(
			method.Name().String(),
			callInfo.ArgumentCount,
			paramCount,
		))
	}

	return value.Undefined
}

// Include a mixin in a class/mixin.
func (vm *VM) opInclude() (err value.Value) {
	mixinVal := vm.pop()
	targetValue := vm.pop()

	mixin, ok := mixinVal.MustReference().(*value.Mixin)
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

// Define a new method alias
func (vm *VM) opDefMethodAlias() {
	newName := vm.pop().MustSymbol()
	oldName := vm.pop().MustSymbol()
	methodContainer := vm.peek()

	switch m := methodContainer.SafeAsReference().(type) {
	case *value.Class:
		m.Methods[newName] = m.Methods[oldName]
	default:
		panic(fmt.Sprintf("invalid method container: %s", methodContainer.Inspect()))
	}
}

// Define a new method
func (vm *VM) opDefMethod() {
	name := vm.pop().MustSymbol()
	body := vm.pop().MustReference().(*BytecodeFunction)
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
	body := vm.pop().MustReference().(*BytecodeFunction)
	namespace := vm.pop()
	vm.executeNamespaceBody(namespace, body)
}

// Execute a chunk of bytecode
func (vm *VM) opExec() {
	bytecodeFunc := vm.pop().MustReference().(*BytecodeFunction)
	vm.executeFunc(bytecodeFunc)
}

// Define a getter method
func (vm *VM) opDefGetter() {
	name := vm.pop().MustSymbol()
	methodContainer := vm.peek()

	switch m := methodContainer.SafeAsReference().(type) {
	case *value.Class:
		DefineGetter(&m.MethodContainer, name)
	default:
		panic(fmt.Sprintf("cannot define a getter in an invalid method container: %s", methodContainer.Inspect()))
	}
}

// Define a setter method
func (vm *VM) opDefSetter() {
	name := vm.pop().MustSymbol()
	methodContainer := vm.peek()

	switch m := methodContainer.SafeAsReference().(type) {
	case *value.Class:
		DefineSetter(&m.MethodContainer, name)
	default:
		panic(fmt.Sprintf("cannot define a setter in an invalid method container: %s", methodContainer.Inspect()))
	}
}

func (vm *VM) addCallFrame(cf CallFrame) {
	*vm.cfpGet() = cf
	vm.cfpIncrement()
}

// preserve the current state of the vm in a call frame
func (vm *VM) createCurrentCallFrame() {
	vm.addCallFrame(
		CallFrame{
			bytecode:   vm.bytecode,
			ip:         vm.ip,
			fp:         vm.fp,
			localCount: vm.localCount,
		},
	)
}

// set up the vm to execute a namespace body
func (vm *VM) executeNamespaceBody(namespace value.Value, body *BytecodeFunction) {
	vm.createCurrentCallFrame()

	vm.bytecode = body
	vm.fp = vm.sp
	vm.ipSet(&body.Instructions[0])
	vm.localCount = 1
	// set namespace as `self`
	vm.push(namespace)
}

// set up the vm to execute a bytecode function
func (vm *VM) executeFunc(fn *BytecodeFunction) {
	vm.createCurrentCallFrame()

	vm.bytecode = fn
	vm.fp = vm.sp
	vm.ipSet(&fn.Instructions[0])
	vm.localCount = 1
	vm.push(value.Ref(value.GlobalObject))
}

// Set a local variable or value.
func (vm *VM) opSetLocal(index int) {
	vm.setLocalValue(index, vm.peek())
}

// Set a local variable or value.
func (vm *VM) setLocalValue(index int, val value.Value) {
	*vm.fpAdd(index) = val
}

// Read a local variable or value.
func (vm *VM) opGetLocal(index int) {
	vm.push(vm.getLocalValue(index))
}

// Read a local variable or value.
func (vm *VM) getLocalValue(index int) value.Value {
	return *vm.fpAdd(index)
}

// Set an upvalue.
func (vm *VM) opSetUpvalue(index int) {
	vm.setUpvalueValue(index, vm.peek())
}

// Set an upvalue.
func (vm *VM) setUpvalueValue(index int, val value.Value) {
	*vm.upvalues[index].location = val
}

// Read an upvalue.
func (vm *VM) opGetUpvalue(index int) {
	vm.push(vm.getUpvalueValue(index))
}

// Read an upvalue.
func (vm *VM) getUpvalueValue(index int) value.Value {
	return *vm.upvalues[index].location
}

// Closes all upvalues up to the given local slot.
func (vm *VM) opCloseUpvalues(lastToClose *value.Value) {
	for {
		if vm.openUpvalueHead == nil ||
			uintptr(unsafe.Pointer(vm.openUpvalueHead.location)) <
				uintptr(unsafe.Pointer(lastToClose)) {
			break
		}

		currentUpvalue := vm.openUpvalueHead
		// move the variable from the stack to the heap
		// inside of the upvalue
		currentUpvalue.closed = *currentUpvalue.location
		// the location pointer now points to the `closed` field
		// within the upvalue
		currentUpvalue.location = &currentUpvalue.closed
		vm.openUpvalueHead = currentUpvalue.next
	}
}

// Set the superclass/parent of a class
func (vm *VM) opSetSuperclass() {
	newSuperclass := vm.pop().MustReference().(*value.Class)
	class := vm.pop().MustReference().(*value.Class)

	if class.Parent != nil {
		return
	}

	class.Parent = newSuperclass
}

// Look for a constant with the given name.
func (vm *VM) opGetConst(nameIndex int) (err value.Value) {
	symbol := vm.bytecode.Values[nameIndex].MustSymbol()

	val := value.RootModule.Constants.Get(symbol)
	if val.IsUndefined() {
		return value.Ref(value.Errorf(value.NoConstantErrorClass, "undefined constant `%s`", symbol.String()))
	}

	vm.push(val)
	return value.Undefined
}

// Get the iterator of the value on top of the stack.
func (vm *VM) opGetIterator() value.Value {
	val := vm.peek()
	result, err := vm.CallMethodByName(iteratorSymbol, val)
	if !err.IsUndefined() {
		return err
	}

	vm.replace(result)
	return value.Undefined
}

var nextSymbol = value.ToSymbol("next")
var stopIterationSymbol = value.ToSymbol("stop_iteration")
var iteratorSymbol = value.ToSymbol("iter")

// Drive the for..in loop.
func (vm *VM) opForIn() value.Value {
	iterator := vm.pop()
	result, err := vm.CallMethodByName(nextSymbol, iterator)
	if err.IsSymbol() && err.AsSymbol() == stopIterationSymbol {
		vm.ipIncrementBy(int(vm.readUint16()))
		return value.Undefined
	}
	if !err.IsUndefined() {
		return err
	}

	vm.push(result)
	vm.ipIncrementBy(2)
	return value.Undefined
}

var toStringSymbol = value.ToSymbol("to_string")

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
				str, ok := strVal.MustReference().(value.String)
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
			element := elementVal.AsFloat64()
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
			element := elementVal.AsInt64()
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
			element := elementVal.AsUInt64()
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
			element := elementVal.AsSymbol()
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

	vm.spIncrementBy(-dynamicElements)
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
			element := elementVal.AsFloat64()
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
			element := elementVal.AsInt64()
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
			element := elementVal.AsUInt64()
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
			element := elementVal.AsSymbol()
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

	vm.spIncrementBy(-dynamicElements)
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
			element := elementVal.AsFloat64()
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
			element := elementVal.AsInt64()
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
			element := elementVal.AsUInt64()
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
			element := elementVal.AsSymbol()
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
	vm.spIncrementBy(-dynamicElements)
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
		to := vm.pop()
		from := vm.pop()
		newRange = value.Ref(value.NewClosedRange(from, to))
	case bytecode.OPEN_RANGE_FLAG:
		to := vm.pop()
		from := vm.pop()
		newRange = value.Ref(value.NewOpenRange(from, to))
	case bytecode.LEFT_OPEN_RANGE_FLAG:
		to := vm.pop()
		from := vm.pop()
		newRange = value.Ref(value.NewLeftOpenRange(from, to))
	case bytecode.RIGHT_OPEN_RANGE_FLAG:
		to := vm.pop()
		from := vm.pop()
		newRange = value.Ref(value.NewRightOpenRange(from, to))
	case bytecode.BEGINLESS_CLOSED_RANGE_FLAG:
		newRange = value.Ref(value.NewBeginlessClosedRange(vm.pop()))
	case bytecode.BEGINLESS_OPEN_RANGE_FLAG:
		newRange = value.Ref(value.NewBeginlessOpenRange(vm.pop()))
	case bytecode.ENDLESS_CLOSED_RANGE_FLAG:
		newRange = value.Ref(value.NewEndlessClosedRange(vm.pop()))
	case bytecode.ENDLESS_OPEN_RANGE_FLAG:
		newRange = value.Ref(value.NewEndlessOpenRange(vm.pop()))
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
	constVal := vm.pop()
	constName := vm.pop().MustSymbol()
	namespace := vm.pop()
	var constants value.ConstantContainer

	switch n := namespace.MustReference().(type) {
	case *value.Class:
		constants = n.ConstantContainer
	case *value.Module:
		constants = n.ConstantContainer
	case *value.Interface:
		constants = n.ConstantContainer
	default:
		panic(
			fmt.Sprintf(
				"tried to define a constant under an invalid namespace: %T",
				namespace,
			),
		)
	}

	constants.AddConstant(constName, constVal)
}

// Leave a local scope and pop all local variables associated with it.
func (vm *VM) opLeaveScope(lastLocalIndex, varsToPop int) {
	firstLocalIndex := lastLocalIndex - varsToPop
	for i := lastLocalIndex; i > firstLocalIndex; i-- {
		vm.stack[i] = value.Undefined
	}
}

// Register slots for local variables and values.
func (vm *VM) opPrepLocals(count int) {
	vm.spIncrementBy(count)
	vm.localCount += count
}

// Push an element on top of the value stack.
func (vm *VM) push(val value.Value) {
	*vm.sp = val
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
func (vm *VM) pop() value.Value {
	vm.spIncrementBy(-1)
	val := *vm.sp
	*vm.sp = value.Undefined
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
	vm.spIncrementBy(-n)
}

// Pop one element off the value stack skipping the first one.
func (vm *VM) popSkipOne() {
	vm.spIncrementBy(-1)
	*vm.spAdd(-1) = *vm.sp
}

// Pop n elements off the value stack skipping the first one.
func (vm *VM) popNSkipOne(n int) {
	*vm.spAdd(-n - 1) = *vm.spAdd(-1)
	for i := vm.spOffset() - 1; i >= vm.spOffset()-n; i-- {
		*vm.spAdd(i) = value.Undefined
	}
	vm.spIncrementBy(-n)
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

// Increment the element on top of the stack
func (vm *VM) opIncrement() (err value.Value) {
	return vm.unaryOperation(value.Increment, symbol.OpIncrement)
}

// Decrement the element on top of the stack
func (vm *VM) opDecrement() (err value.Value) {
	return vm.unaryOperation(value.Decrement, symbol.OpDecrement)
}

// Negate the element on top of the stack
func (vm *VM) opNegate() (err value.Value) {
	return vm.unaryOperation(value.Negate, symbol.OpNegate)
}

// Perform unary plus on the element on top of the stack
func (vm *VM) opUnaryPlus() (err value.Value) {
	return vm.unaryOperation(value.UnaryPlus, symbol.OpUnaryPlus)
}

// Preform bitwise not on the element on top of the stack
func (vm *VM) opBitwiseNot() (err value.Value) {
	return vm.unaryOperation(value.BitwiseNot, symbol.OpBitwiseNot)
}

func (vm *VM) opAppendAt() value.Value {
	val := vm.pop()
	key := vm.pop()
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
	val := vm.peek()
	key := vm.peekAt(1)
	collection := vm.peekAt(2)

	result, err := value.SubscriptSet(collection, key, val)
	if !err.IsUndefined() {
		return err
	}
	if !result.IsUndefined() {
		vm.popN(2)
		vm.replace(result)
		return value.Undefined
	}

	er := vm.callMethodOnStackByName(symbol.OpSubscriptSet, 2)
	if !er.IsUndefined() {
		return er
	}
	return value.Undefined
}

func (vm *VM) opIsA() (err value.Value) {
	classVal := vm.pop()
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
	classVal := vm.pop()
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
	return vm.binaryOperation(value.BitwiseAnd, symbol.OpAnd)
}

// Perform a bitwise AND NOT and push the result to the stack.
func (vm *VM) opBitwiseAndNot() (err value.Value) {
	return vm.binaryOperation(value.BitwiseAndNot, symbol.OpAndNot)
}

// Get the value under the given key and push the result to the stack.
func (vm *VM) opSubscript() (err value.Value) {
	return vm.binaryOperation(value.Subscript, symbol.OpSubscript)
}

// Perform a bitwise OR and push the result to the stack.
func (vm *VM) opBitwiseOr() (err value.Value) {
	return vm.binaryOperation(value.BitwiseOr, symbol.OpOr)
}

// Perform a bitwise XOR and push the result to the stack.
func (vm *VM) opBitwiseXor() (err value.Value) {
	return vm.binaryOperation(value.BitwiseXor, symbol.OpXor)
}

// Perform a comparison and push the result to the stack.
func (vm *VM) opCompare() (err value.Value) {
	return vm.binaryOperation(value.Compare, symbol.OpSpaceship)
}

// Perform opModulo and push the result to the stack.
func (vm *VM) opModulo() (err value.Value) {
	return vm.binaryOperation(value.Modulo, symbol.OpModulo)
}

// Check whether two top elements on the stack are opEqual and push the result to the stack.
func (vm *VM) opEqual() (err value.Value) {
	return vm.callEqualityOperator(value.Equal, symbol.OpEqual)
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
	return vm.callNegatedEqualityOperator(value.NotEqual, symbol.OpEqual)
}

// Check whether two top elements on the stack are equal and push the result to the stack.
func (vm *VM) opLaxEqual() (err value.Value) {
	return vm.callEqualityOperator(value.LaxEqual, symbol.OpLaxEqual)
}

// Check whether two top elements on the stack are not and equal push the result to the stack.
func (vm *VM) opLaxNotEqual() (err value.Value) {
	return vm.callNegatedEqualityOperator(value.LaxNotEqual, symbol.OpLaxEqual)
}

// Check whether two top elements on the stack are strictly equal push the result to the stack.
func (vm *VM) opStrictEqual() (err value.Value) {
	return vm.binaryOperationWithoutErr(value.StrictEqual, symbol.OpStrictEqual)
}

// Check whether two top elements on the stack are strictly not equal push the result to the stack.
func (vm *VM) opStrictNotEqual() (err value.Value) {
	return vm.negatedBinaryOperationWithoutErr(value.StrictNotEqual, symbol.OpStrictEqual)
}

// Check whether the first operand is greater than the second and push the result to the stack.
func (vm *VM) opGreaterThan() (err value.Value) {
	return vm.binaryOperation(value.GreaterThan, symbol.OpGreaterThan)
}

// Check whether the first operand is greater than or equal to the second and push the result to the stack.
func (vm *VM) opGreaterThanEqual() (err value.Value) {
	return vm.binaryOperation(value.GreaterThanEqual, symbol.OpGreaterThanEqual)
}

// Check whether the first operand is less than the second and push the result to the stack.
func (vm *VM) opLessThan() (err value.Value) {
	return vm.binaryOperation(value.LessThan, symbol.OpLessThan)
}

// Check whether the first operand is less than or equal to the second and push the result to the stack.
func (vm *VM) opLessThanEqual() (err value.Value) {
	return vm.binaryOperation(value.LessThanEqual, symbol.OpLessThanEqual)
}

// Perform a left bitshift and push the result to the stack.
func (vm *VM) opLeftBitshift() (err value.Value) {
	return vm.binaryOperation(value.LeftBitshift, symbol.OpLeftBitshift)
}

// Perform a logical left bitshift and push the result to the stack.
func (vm *VM) opLogicalLeftBitshift() (err value.Value) {
	return vm.binaryOperation(value.LogicalLeftBitshift, symbol.OpLogicalLeftBitshift)
}

// Perform a right bitshift and push the result to the stack.
func (vm *VM) opRightBitshift() (err value.Value) {
	return vm.binaryOperation(value.RightBitshift, symbol.OpRightBitshift)
}

// Perform a logical right bitshift and push the result to the stack.
func (vm *VM) opLogicalRightBitshift() (err value.Value) {
	return vm.binaryOperation(value.LogicalRightBitshift, symbol.OpLogicalRightBitshift)
}

// Add two operands together and push the result to the stack.
func (vm *VM) opAdd() (err value.Value) {
	return vm.binaryOperation(value.Add, symbol.OpAdd)
}

// Subtract two operands and push the result to the stack.
func (vm *VM) opSubtract() (err value.Value) {
	return vm.binaryOperation(value.Subtract, symbol.OpSubtract)
}

// Multiply two operands together and push the result to the stack.
func (vm *VM) opMultiply() (err value.Value) {
	return vm.binaryOperation(value.Multiply, symbol.OpMultiply)
}

// Divide two operands and push the result to the stack.
func (vm *VM) opDivide() (err value.Value) {
	return vm.binaryOperation(value.Divide, symbol.OpDivide)
}

// Exponentiate two operands and push the result to the stack.
func (vm *VM) opExponentiate() (err value.Value) {
	return vm.binaryOperation(value.Exponentiate, symbol.OpExponentiate)
}

// Throw an error when the value on top of the stack is `nil`
func (vm *VM) opMust() {
	val := vm.peek()
	if value.IsNil(val) {
		vm.throw(value.Ref(value.NewUnexpectedNilError()))
	}
}

// Throw an error when the second value on the stack is not an instance of the class/mixin on top of the stack
func (vm *VM) mustAs() {
	class := vm.pop().MustReference().(*value.Class)
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
	vm.rethrow(err, value.String(vm.BuildStackTrace()))
}

func (vm *VM) rethrow(err value.Value, stackTrace value.String) {
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

		if vm.mode == singleFunctionCallMode || vm.cfp == uintptr(unsafe.Pointer(&vm.callFrames[0])) {
			vm.mode = errorMode
			vm.errStackTrace = string(stackTrace)
			vm.push(err)
			panic(stopVM{})
		}

		vm.restoreLastFrame()
	}
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
