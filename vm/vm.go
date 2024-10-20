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
	normalMode             mode = iota
	singleFunctionCallMode      // the VM should halt after executing a single method
	errorMode                   // the VM stopped after encountering an uncaught error
)

// A single instance of the Elk Virtual Machine.
type VM struct {
	bytecode        *BytecodeFunction
	upvalues        []*Upvalue
	openUpvalueHead *Upvalue      // linked list of open upvalues, living on the stack
	ip              int           // Instruction pointer -- points to the next bytecode instruction
	sp              int           // Stack pointer -- points to the offset where the next element will be pushed to
	fp              int           // Frame pointer -- points to the offset where the current frame starts
	localCount      int           // the amount of registered locals
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
func (vm *VM) InterpretTopLevel(fn *BytecodeFunction) (value.Value, value.Value) {
	vm.bytecode = fn
	vm.ip = 0
	vm.push(value.GlobalObject)
	vm.push(value.RootModule)
	vm.push(value.GlobalObjectSingletonClass)
	vm.localCount = 3
	vm.run()
	err := vm.Err()
	if err != nil {
		return nil, err
	}
	return vm.peek(), nil
}

// Execute the given bytecode chunk.
func (vm *VM) InterpretREPL(fn *BytecodeFunction) (value.Value, value.Value) {
	vm.bytecode = fn
	vm.ip = 0
	if vm.sp == 0 {
		// populate the predeclared local variables
		vm.push(value.GlobalObject)               // populate self
		vm.push(value.RootModule)                 // populate constant container
		vm.push(value.GlobalObjectSingletonClass) // populate method container
		vm.localCount = 3
	} else {
		// pop the return value of the last run
		vm.pop()
	}
	vm.run()

	err := vm.Err()
	if err != nil {
		return nil, err
	}
	return vm.peek(), nil
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

	return nil
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
	return vm.stack[:vm.sp]
}

func (vm *VM) InspectStack() {
	fmt.Println("stack:")
	for i, value := range vm.Stack() {
		if value == nil {
			fmt.Printf("%d => <Go nil!>\n", i)
			continue
		}
		fmt.Printf("%d => %s\n", i, value.Inspect())
	}
}

func (vm *VM) throwIfErr(err value.Value) {
	if err != nil {
		vm.throw(err)
	}
}

var callSymbol = value.ToSymbol("call")

// Call a callable value from Go code, preserving the state of the VM.
func (vm *VM) CallCallable(args ...value.Value) (value.Value, value.Value) {
	function := args[0]
	switch f := function.(type) {
	case *Closure:
		return vm.CallClosure(f, args[1:]...)
	default:
		return vm.CallMethodByName(callSymbol, args...)
	}
}

// Call an Elk closure from Go code, preserving the state of the VM.
func (vm *VM) CallClosure(closure *Closure, args ...value.Value) (value.Value, value.Value) {
	if closure.Bytecode.ParameterCount() != len(args) {
		return nil, value.NewWrongArgumentCountError(
			closure.Bytecode.Name().String(),
			len(args),
			closure.Bytecode.ParameterCount(),
		)
	}

	vm.createCurrentCallFrame()
	vm.bytecode = closure.Bytecode
	vm.fp = vm.sp
	vm.ip = 0
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
	if err != nil {
		vm.mode = normalMode
		vm.restoreLastFrame()

		return nil, err
	}
	vm.mode = normalMode
	return vm.pop(), nil
}

// Call an Elk method from Go code, preserving the state of the VM.
func (vm *VM) CallMethodByName(name value.Symbol, args ...value.Value) (value.Value, value.Value) {
	self := args[0]
	class := self.DirectClass()
	method := class.LookupMethod(name)
	if method == nil {
		return nil, value.NewNoMethodError(string(name.ToString()), self)
	}
	return vm.CallMethod(method, args...)
}

func (vm *VM) CallMethod(method value.Method, args ...value.Value) (value.Value, value.Value) {
	self := args[0]
	if method.ParameterCount() != len(args)-1 {
		return nil, value.NewWrongArgumentCountError(
			method.Name().String(),
			len(args)-1,
			method.ParameterCount(),
		)
	}

	switch m := method.(type) {
	case *BytecodeFunction:
		vm.createCurrentCallFrame()
		vm.bytecode = m
		vm.fp = vm.sp
		vm.ip = 0
		vm.localCount = len(args)
		vm.mode = singleFunctionCallMode
		for _, arg := range args {
			vm.push(arg)
		}
		vm.run()
		err := vm.Err()
		if err != nil {
			vm.mode = normalMode
			vm.restoreLastFrame()

			return nil, err
		}
		vm.mode = normalMode
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
func (vm *VM) callMethodOnStack(method value.Method, args int) value.Value {
	switch m := method.(type) {
	case *BytecodeFunction:
		vm.createCurrentCallFrame()
		vm.bytecode = m
		vm.fp = vm.sp - args - 1
		vm.localCount = args + 1
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

func (vm *VM) callMethodOnStackByName(name value.Symbol, args int) value.Value {
	self := vm.stack[vm.sp-args-1]
	class := self.DirectClass()
	method := class.LookupMethod(name)
	if method == nil {
		return value.NewNoMethodError(string(name.ToString()), self)
	}

	return vm.callMethodOnStack(method, args)
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
		case bytecode.RETURN_FINALLY:
			if vm.jumpToFinallyForReturn() {
				continue
			}

			// return normally
			if len(vm.callFrames) == 0 {
				return
			}
			vm.returnFromFunction()
			if vm.mode == singleFunctionCallMode {
				return
			}
		case bytecode.RETURN:
			if len(vm.callFrames) == 0 {
				return
			}
			vm.returnFromFunction()
			if vm.mode == singleFunctionCallMode {
				return
			}
		case bytecode.RETURN_FIRST_ARG:
			vm.getLocal(1)
			if len(vm.callFrames) == 0 {
				return
			}
			vm.returnFromFunction()
			if vm.mode == singleFunctionCallMode {
				return
			}
		case bytecode.RETURN_SELF:
			vm.self()
			if len(vm.callFrames) == 0 {
				return
			}
			vm.returnFromFunction()
			if vm.mode == singleFunctionCallMode {
				return
			}
		case bytecode.CLOSURE:
			vm.closure()
		case bytecode.JUMP_TO_FINALLY:
			leftFinallyCount := vm.peek().(value.SmallInt)
			jumpOffset := vm.peekAt(1).(value.SmallInt)

			if leftFinallyCount > 0 {
				vm.replace(leftFinallyCount - 1)
				if !vm.jumpToFinallyForBreakOrContinue() {
					panic("could not find a finally block to jump to in JUMP_TO_FINALLY")
				}
				continue
			}

			vm.popN(2)
			vm.ip = int(jumpOffset)
		case bytecode.DUP:
			vm.push(vm.peek())
		case bytecode.SWAP:
			vm.swap()
		case bytecode.DUP_N:
			n := int(vm.readByte())
			for _, element := range vm.stack[vm.sp-n : vm.sp] {
				vm.push(element)
			}
		case bytecode.CONSTANT_CONTAINER:
			vm.constantContainer()
		case bytecode.METHOD_CONTAINER:
			vm.methodContainer()
		case bytecode.SELF:
			vm.self()
		case bytecode.DEF_SINGLETON:
			vm.throwIfErr(vm.defineSingleton())
		case bytecode.GET_SINGLETON:
			vm.throwIfErr(vm.getSingletonClass())
		case bytecode.GET_CLASS:
			vm.getClass()
		case bytecode.DEF_ALIAS:
			vm.throwIfErr(vm.defineAlias())
		case bytecode.DEF_GETTER:
			vm.throwIfErr(vm.defineGetter())
		case bytecode.DEF_SETTER:
			vm.throwIfErr(vm.defineSetter())
		case bytecode.DEF_CLASS:
			vm.throwIfErr(vm.defineClass())
		case bytecode.DEF_MODULE:
			vm.throwIfErr(vm.defineModule())
		case bytecode.DEF_MIXIN:
			vm.throwIfErr(vm.defineMixin())
		case bytecode.DEF_METHOD:
			vm.throwIfErr(vm.defineMethod())
		case bytecode.INCLUDE:
			vm.throwIfErr(vm.includeMixin())
		case bytecode.DOC_COMMENT:
			vm.throwIfErr(vm.docComment())
		case bytecode.APPEND:
			vm.appendCollection()
		case bytecode.MAP_SET:
			vm.mapSet()
		case bytecode.COPY:
			vm.copy()
		case bytecode.APPEND_AT:
			vm.throwIfErr(vm.appendAt())
		case bytecode.SUBSCRIPT:
			vm.throwIfErr(vm.subscript())
		case bytecode.SUBSCRIPT_SET:
			vm.throwIfErr(vm.subscriptSet())
		case bytecode.INSTANTIATE8:
			vm.throwIfErr(
				vm.instantiate(int(vm.readByte())),
			)
		case bytecode.INSTANTIATE16:
			vm.throwIfErr(
				vm.instantiate(int(vm.readUint16())),
			)
		case bytecode.INSTANTIATE32:
			vm.throwIfErr(
				vm.instantiate(int(vm.readUint32())),
			)
		case bytecode.GET_IVAR8:
			vm.throwIfErr(
				vm.getInstanceVariable(int(vm.readByte())),
			)
		case bytecode.GET_IVAR16:
			vm.throwIfErr(
				vm.getInstanceVariable(int(vm.readUint16())),
			)
		case bytecode.GET_IVAR32:
			vm.throwIfErr(
				vm.getInstanceVariable(int(vm.readUint32())),
			)
		case bytecode.SET_IVAR8:
			vm.throwIfErr(
				vm.setInstanceVariable(int(vm.readByte())),
			)
		case bytecode.SET_IVAR16:
			vm.throwIfErr(
				vm.setInstanceVariable(int(vm.readUint16())),
			)
		case bytecode.SET_IVAR32:
			vm.throwIfErr(
				vm.setInstanceVariable(int(vm.readUint32())),
			)
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
		case bytecode.CALL8:
			vm.throwIfErr(
				vm.call(int(vm.readByte())),
			)
		case bytecode.CALL16:
			vm.throwIfErr(
				vm.call(int(vm.readUint16())),
			)
		case bytecode.CALL32:
			vm.throwIfErr(
				vm.call(int(vm.readUint32())),
			)
		case bytecode.CALL_SELF8:
			vm.throwIfErr(
				vm.callFunction(int(vm.readByte())),
			)
		case bytecode.CALL_SELF16:
			vm.throwIfErr(
				vm.callFunction(int(vm.readUint16())),
			)
		case bytecode.CALL_SELF32:
			vm.throwIfErr(
				vm.callFunction(int(vm.readUint32())),
			)
		case bytecode.INSTANCE_OF:
			vm.throwIfErr(vm.instanceOf())
		case bytecode.IS_A:
			vm.throwIfErr(vm.isA())
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
		case bytecode.UNARY_PLUS:
			vm.throwIfErr(vm.unaryPlus())
		case bytecode.BITWISE_NOT:
			vm.throwIfErr(vm.bitwiseNot())
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
			vm.throwIfErr(vm.increment())
		case bytecode.DECREMENT:
			vm.throwIfErr(vm.decrement())
		case bytecode.GET_LOCAL8:
			vm.getLocal(int(vm.readByte()))
		case bytecode.GET_LOCAL16:
			vm.getLocal(int(vm.readUint16()))
		case bytecode.SET_LOCAL8:
			vm.setLocal(int(vm.readByte()))
		case bytecode.SET_LOCAL16:
			vm.setLocal(int(vm.readUint16()))
		case bytecode.GET_UPVALUE8:
			vm.getUpvalue(int(vm.readByte()))
		case bytecode.GET_UPVALUE16:
			vm.getUpvalue(int(vm.readUint16()))
		case bytecode.SET_UPVALUE8:
			vm.setUpvalue(int(vm.readByte()))
		case bytecode.SET_UPVALUE16:
			vm.setUpvalue(int(vm.readUint16()))
		case bytecode.CLOSE_UPVALUE8:
			last := &vm.stack[vm.fp+int(vm.readByte())]
			vm.closeUpvalues(last)
		case bytecode.CLOSE_UPVALUE16:
			last := &vm.stack[vm.fp+int(vm.readUint16())]
			vm.closeUpvalues(last)
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
		case bytecode.NEW_RANGE:
			vm.newRange()
		case bytecode.NEW_ARRAY_TUPLE8:
			vm.newArrayTuple(int(vm.readByte()))
		case bytecode.NEW_ARRAY_TUPLE32:
			vm.newArrayTuple(int(vm.readUint32()))
		case bytecode.NEW_ARRAY_LIST8:
			vm.throwIfErr(
				vm.newArrayList(int(vm.readByte())),
			)
		case bytecode.NEW_ARRAY_LIST32:
			vm.throwIfErr(
				vm.newArrayList(int(vm.readUint32())),
			)
		case bytecode.NEW_HASH_SET8:
			vm.throwIfErr(
				vm.newHashSet(int(vm.readByte())),
			)
		case bytecode.NEW_HASH_SET32:
			vm.throwIfErr(
				vm.newHashSet(int(vm.readUint32())),
			)
		case bytecode.NEW_HASH_MAP8:
			vm.throwIfErr(
				vm.newHashMap(int(vm.readByte())),
			)
		case bytecode.NEW_HASH_MAP32:
			vm.throwIfErr(
				vm.newHashMap(int(vm.readUint32())),
			)
		case bytecode.NEW_HASH_RECORD8:
			vm.throwIfErr(
				vm.newHashRecord(int(vm.readByte())),
			)
		case bytecode.NEW_HASH_RECORD32:
			vm.throwIfErr(
				vm.newHashRecord(int(vm.readUint32())),
			)
		case bytecode.NEW_STRING8:
			vm.throwIfErr(vm.newString(int(vm.readByte())))
		case bytecode.NEW_STRING32:
			vm.throwIfErr(vm.newString(int(vm.readUint32())))
		case bytecode.NEW_SYMBOL8:
			vm.throwIfErr(vm.newSymbol(int(vm.readByte())))
		case bytecode.NEW_SYMBOL32:
			vm.throwIfErr(vm.newSymbol(int(vm.readUint32())))
		case bytecode.NEW_REGEX8:
			vm.throwIfErr(vm.newRegex(vm.readByte(), int(vm.readByte())))
		case bytecode.NEW_REGEX32:
			vm.throwIfErr(vm.newRegex(vm.readByte(), int(vm.readUint32())))
		case bytecode.FOR_IN:
			vm.throwIfErr(vm.forIn())
		case bytecode.GET_ITERATOR:
			vm.throwIfErr(vm.getIterator())
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
		case bytecode.THROW:
			vm.throw(vm.pop())
		case bytecode.RETHROW:
			err := vm.pop()
			stackTrace := vm.pop().(value.String)
			vm.rethrow(err, stackTrace)
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
		case bytecode.BITWISE_AND_NOT:
			vm.throwIfErr(vm.bitwiseAndNot())
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
		case bytecode.LAX_EQUAL:
			vm.throwIfErr(vm.laxEqual())
		case bytecode.LAX_NOT_EQUAL:
			vm.throwIfErr(vm.laxNotEqual())
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
		case bytecode.INSPECT_STACK:
			vm.InspectStack()
		default:
			panic(fmt.Sprintf("Unknown bytecode instruction: %#v", instruction))
		}

		if vm.mode == errorMode {
			return
		}
	}

}

func (vm *VM) closure() {
	function := vm.peek().(*BytecodeFunction)
	closure := NewClosure(function, vm.selfValue())
	vm.replace(closure)

	for i := range len(closure.Upvalues) {
		flags := bitfield.BitField8FromInt(vm.readByte())
		var upIndex int
		if flags.HasFlag(UpvalueLongIndexFlag) {
			upIndex = int(vm.readUint16())
		} else {
			upIndex = int(vm.readByte())
		}

		if flags.HasFlag(UpvalueLocalFlag) {
			closure.Upvalues[i] = vm.captureUpvalue(&vm.stack[vm.fp+upIndex])
		} else {
			closure.Upvalues[i] = vm.upvalues[upIndex]
		}
	}
	vm.ip++ // skip past the terminator
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

func (vm *VM) lastCallFrame() *CallFrame {
	lastIndex := len(vm.callFrames) - 1
	return &vm.callFrames[lastIndex]
}

// Restore the state of the VM to the last call frame.
func (vm *VM) restoreLastFrame() {
	lastIndex := len(vm.callFrames) - 1
	cf := vm.callFrames[lastIndex]
	// reset the popped call frame
	vm.callFrames[lastIndex] = CallFrame{}
	vm.callFrames = vm.callFrames[:lastIndex]

	vm.ip = cf.ip
	vm.closeUpvalues(&vm.stack[vm.fp])
	vm.popN(vm.sp - vm.fp)
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
	for j := range len(vm.callFrames) {
		callFrame := &vm.callFrames[j]
		addStackTraceEntry(
			&buffer,
			j,
			callFrame.FileName(),
			callFrame.LineNumber(),
			callFrame.Name().String(),
		)
		i = j
	}
	addStackTraceEntry(
		&buffer,
		i+1,
		vm.bytecode.FileName(),
		vm.bytecode.GetLineNumber(vm.ip-1),
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
			"value `%s` cannot have a singleton class",
			val.Inspect(),
		)
	}

	vm.push(singleton)
	return nil
}

func (vm *VM) getClass() {
	val := vm.pop()
	class := val.Class()
	vm.push(class)
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
		return value.NewNoMethodError(string(callInfo.Name.ToString()), self)
	}

	// shift all arguments one slot forward to make room for self
	for i := 0; i < callInfo.ArgumentCount; i++ {
		vm.stack[vm.sp-i] = vm.stack[vm.sp-i-1]
	}
	vm.stack[vm.sp-callInfo.ArgumentCount] = self
	vm.sp++

	switch m := method.(type) {
	case *BytecodeFunction:
		return vm.callBytecodeFunction(m, callInfo)
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	}

	panic(fmt.Sprintf("tried to call a method that is neither bytecode nor native: %#v", method))
}

// Set the value of an instance variable
func (vm *VM) setInstanceVariable(nameIndex int) (err value.Value) {
	name := vm.bytecode.Values[nameIndex].(value.Symbol)
	val := vm.peek()

	self := vm.selfValue()
	ivars := self.InstanceVariables()
	if ivars == nil {
		return value.NewCantSetInstanceVariablesOnPrimitiveError(self.Inspect())
	}

	ivars.Set(name, val)
	return nil
}

// Get the value of an instance variable
func (vm *VM) getInstanceVariable(nameIndex int) (err value.Value) {
	name := vm.bytecode.Values[nameIndex].(value.Symbol)

	self := vm.selfValue()
	ivars := self.InstanceVariables()
	if ivars == nil {
		return value.NewCantAccessInstanceVariablesOnPrimitiveError(self.Inspect())
	}

	val := ivars.Get(name)
	if val == nil {
		vm.push(value.Nil)
	} else {
		vm.push(val)
	}

	return nil
}

// Pop the value on top of the stack and push its copy.
func (vm *VM) copy() {
	element := vm.peek()
	vm.replace(element.Copy())
}

// Set the value under the given key in a hash-map or hash-record
func (vm *VM) mapSet() {
	val := vm.pop()
	key := vm.pop()
	collection := vm.peek()

	switch c := collection.(type) {
	case *value.HashMap:
		HashMapSet(vm, c, key, val)
	case *value.HashRecord:
		HashRecordSet(vm, c, key, val)
	case value.UndefinedType:
		panic("undefined hash map base")
	default:
		panic(fmt.Sprintf("invalid map to set a value in: %#v", collection))
	}
}

// Append an element to a list, arrayTuple or hashSet.
func (vm *VM) appendCollection() {
	element := vm.pop()
	collection := vm.peek()

	switch c := collection.(type) {
	case *value.ArrayTuple:
		c.Append(element)
	case *value.ArrayList:
		c.Append(element)
	case *value.HashSet:
		HashSetAppend(vm, c, element)
	case value.UndefinedType:
		vm.replace(&value.ArrayTuple{element})
	default:
		panic(fmt.Sprintf("invalid collection to append to: %#v", collection))
	}
}

// Call a method with an explicit receiver
func (vm *VM) instantiate(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].(*value.CallSiteInfo)

	classIndex := vm.sp - callInfo.ArgumentCount - 1
	classVal := vm.stack[classIndex]
	class := classVal.(*value.Class)

	instance := class.CreateInstance()
	// replace the class with the instance
	vm.stack[classIndex] = instance
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
			return nil
		}

		return value.NewWrongArgumentCountError(
			"#init",
			callInfo.ArgumentCount,
			0,
		)
	default:
		panic(fmt.Sprintf("tried to call an invalid initialiser method: %#v", method))
	}
}

// Call a method in a pattern.
// Return false if the receiver does not have the method
// or it throws TypeError.
func (vm *VM) callPattern(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].(*value.CallSiteInfo)

	self := vm.stack[vm.sp-callInfo.ArgumentCount-1]
	class := self.DirectClass()

	method := class.LookupMethod(callInfo.Name)
	if method == nil {
		vm.popN(callInfo.ArgumentCount + 1)
		vm.push(value.False)
		return nil
	}
	switch m := method.(type) {
	case *BytecodeFunction:
		err = vm.callBytecodeFunction(m, callInfo)
	case *NativeMethod:
		err = vm.callNativeMethod(m, callInfo)
	case *GetterMethod:
		if callInfo.ArgumentCount != 0 {
			return value.NewWrongArgumentCountError(
				method.Name().String(),
				callInfo.ArgumentCount,
				0,
			)
		}
		vm.pop() // pop self
		var result value.Value
		result, err = m.Call(self)
		if err == nil {
			vm.push(result)
		}
	case *SetterMethod:
		if callInfo.ArgumentCount != 1 {
			return value.NewWrongArgumentCountError(
				method.Name().String(),
				callInfo.ArgumentCount,
				1,
			)
		}
		other := vm.pop()
		vm.pop() // pop self
		var result value.Value
		result, err = m.Call(self, other)
		if err == nil {
			vm.push(result)
		}
	default:
		panic(fmt.Sprintf("tried to call an invalid method: %#v", method))
	}

	if err != nil {
		if err.Class() == value.TypeErrorClass {
			vm.push(value.False)
			return nil
		}
		return err
	}

	return nil
}

// Call the `call` method with an explicit receiver
func (vm *VM) call(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].(*value.CallSiteInfo)

	self, isClosure := vm.stack[vm.sp-callInfo.ArgumentCount-1].(*Closure)
	if !isClosure {
		return vm.callMethod(callInfoIndex)
	}

	return vm.callClosure(self, callInfo)
}

// set up the vm to execute a closure
func (vm *VM) callClosure(closure *Closure, callInfo *value.CallSiteInfo) (err value.Value) {
	function := closure.Bytecode
	if err := vm.prepareArguments(function, callInfo); err != nil {
		return err
	}

	vm.createCurrentCallFrame()

	vm.localCount = len(function.parameters) + 1
	vm.bytecode = function
	vm.fp = vm.sp - function.ParameterCount() - 1
	vm.ip = 0
	vm.upvalues = closure.Upvalues

	return nil
}

// Call a method with an explicit receiver
func (vm *VM) callMethod(callInfoIndex int) (err value.Value) {
	callInfo := vm.bytecode.Values[callInfoIndex].(*value.CallSiteInfo)

	self := vm.stack[vm.sp-callInfo.ArgumentCount-1]
	class := self.DirectClass()

	method := class.LookupMethod(callInfo.Name)
	if method == nil {
		vm.popN(callInfo.ArgumentCount + 1)
		return value.NewNoMethodError(string(callInfo.Name.ToString()), self)
	}
	switch m := method.(type) {
	case *BytecodeFunction:
		return vm.callBytecodeFunction(m, callInfo)
	case *NativeMethod:
		return vm.callNativeMethod(m, callInfo)
	case *GetterMethod:
		if callInfo.ArgumentCount != 0 {
			return value.NewWrongArgumentCountError(
				method.Name().String(),
				callInfo.ArgumentCount,
				0,
			)
		}
		vm.pop() // pop self
		result, err := m.Call(self)
		if err != nil {
			return err
		}
		vm.push(result)
		return nil
	case *SetterMethod:
		if callInfo.ArgumentCount != 1 {
			return value.NewWrongArgumentCountError(
				method.Name().String(),
				callInfo.ArgumentCount,
				1,
			)
		}
		other := vm.pop()
		vm.pop() // pop self
		result, err := m.Call(self, other)
		if err != nil {
			return err
		}
		vm.push(result)
		return nil
	default:
		panic(fmt.Sprintf("tried to call an invalid method: %T", method))
	}
}

// set up the vm to execute a native method
func (vm *VM) callNativeMethod(method *NativeMethod, callInfo *value.CallSiteInfo) (err value.Value) {
	if err := vm.prepareArguments(method, callInfo); err != nil {
		return err
	}

	returnVal, err := method.Function(vm, vm.stack[vm.sp-method.ParameterCount()-1:vm.sp])
	vm.popN(method.ParameterCount() + 1)
	if err != nil {
		return err
	}
	vm.push(returnVal)
	return nil
}

// set up the vm to execute a bytecode method
func (vm *VM) callBytecodeFunction(method *BytecodeFunction, callInfo *value.CallSiteInfo) (err value.Value) {
	if err := vm.prepareArguments(method, callInfo); err != nil {
		return err
	}

	vm.createCurrentCallFrame()

	vm.localCount = len(method.parameters) + 1
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
	namedRestParam := method.NamedRestParameter()
	if namedRestParam {
		paramCount -= 1
	}
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
				method.Name().String(),
				posArgCount,
				requiredPosParamCount,
			)
		}

		firstPosRestArg := paramCount - method.PostRestParameterCount() - 1
		lastPosRestArg := callInfo.ArgumentCount - 1 - method.PostRestParameterCount()
		if namedRestParam {
			lastPosRestArg -= callInfo.NamedArgumentCount()
		}
		posRestArgCount := lastPosRestArg - firstPosRestArg + 1
		postArgCount := callInfo.ArgumentCount - lastPosRestArg - 1
		var postArgs []value.Value
		var restList value.ArrayList
		if postArgCount > 0 {
			postArgs = make([]value.Value, postArgCount)
			copy(postArgs, vm.stack[vm.sp-postArgCount:vm.sp])
			vm.popN(postArgCount)
		}

		if posRestArgCount > 0 {
			restList = make(value.ArrayList, posRestArgCount)
		}
		for i := 1; i <= posRestArgCount; i++ {
			restList[posRestArgCount-i] = vm.pop()
		}
		vm.push(&restList)
		for _, postArg := range postArgs {
			vm.push(postArg)
		}

		posParamNames = paramNames[:firstPosRestArg+1]

		if !namedRestParam {
			namedParamNames = paramNames[paramCount-(callInfo.ArgumentCount-posArgCount):]
		}
		spIncrease = paramCount - (callInfo.ArgumentCount - posRestArgCount + 1)
	} else {
		posParamNames = paramNames[:posArgCount]
		namedParamNames = paramNames[posArgCount:]
		spIncrease = paramCount - callInfo.ArgumentCount
	}
	if namedRestParam && len(namedParamNames) > 0 {
		namedParamNames = namedParamNames[:len(namedParamNames)-1]
	}

	var foundNamedArgCount int

	for _, paramName := range posParamNames {
		for j := 0; j < namedArgCount; j++ {
			if paramName == callInfo.NamedArguments[j] {
				return value.NewDuplicatedArgumentError(
					string(method.Name().ToString()),
					string(paramName.ToString()),
				)
			}
		}
	}

	var missingOptionalArgCount int
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
				string(method.Name().ToString()),
				string(paramName.ToString()),
			)
		}

		// the parameter is optional and is not present
		// populate it with undefined
		vm.stack[targetIndex] = value.Undefined
		missingOptionalArgCount++
	}

	unknownNamedArgCount := namedArgCount - foundNamedArgCount
	if namedRestParam {
		hmap := value.NewHashMap(unknownNamedArgCount)
		if unknownNamedArgCount != 0 {
			// construct a hashmap of named arguments
			// that are not defined in the method
			for i, namedArg := range namedArgs {
				if namedArg == value.Undefined {
					continue
				}

				HashMapSet(vm, hmap, callInfo.NamedArguments[i], namedArg)
			}
			additionalNamedArgCount := hmap.Length() - missingOptionalArgCount
			vm.popN(additionalNamedArgCount)
			spIncrease += additionalNamedArgCount
		}
		vm.push(hmap)
	} else if unknownNamedArgCount != 0 {
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

		return value.NewUnknownArgumentsError(
			method.Name().String(),
			unknownNamedArgNames,
		)
	}

	vm.sp += spIncrease
	return nil
}

func (vm *VM) preparePositionalArguments(method value.Method, callInfo *value.CallSiteInfo) (err value.Value) {
	optParamCount := method.OptionalParameterCount()
	postParamCount := method.PostRestParameterCount()
	paramCount := method.ParameterCount()
	namedRestParam := method.NamedRestParameter()
	if namedRestParam {
		paramCount--
	}
	preRestParamCount := paramCount - postParamCount - 1
	reqParamCount := paramCount - optParamCount
	if postParamCount >= 0 {
		reqParamCount -= 1
	}

	if callInfo.ArgumentCount < reqParamCount {
		if postParamCount == -1 {
			return value.NewWrongArgumentCountRangeError(
				method.Name().String(),
				callInfo.ArgumentCount,
				reqParamCount,
				paramCount,
			)
		} else {
			return value.NewWrongArgumentCountRestError(
				method.Name().String(),
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
			method.Name().String(),
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
		var restList value.ArrayList
		restArgCount := callInfo.ArgumentCount - preRestParamCount - postParamCount
		if restArgCount > 0 {
			// rest arguments
			restList = make(value.ArrayList, restArgCount)
			for i := 1; i <= restArgCount; i++ {
				restList[restArgCount-i] = vm.pop()
			}
		}
		vm.push(&restList)
		for _, postArg := range postArgs {
			vm.push(postArg)
		}
	}

	if namedRestParam {
		vm.push(&value.HashMap{})
	}

	return nil
}

// Include a mixin in a class/mixin.
func (vm *VM) includeMixin() (err value.Value) {
	targetValue := vm.pop()
	mixinVal := vm.pop()

	mixin, ok := mixinVal.(*value.Mixin)
	if !ok || !mixin.IsMixin() {
		return value.NewIsNotMixinError(mixinVal.Inspect())
	}

	switch target := targetValue.(type) {
	case *value.Class:
		target.IncludeMixin(mixin)
	default:
		return value.Errorf(
			value.TypeErrorClass,
			"cannot include into an instance of %s: `%s`",
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
	// 	target.D
	// case *value.Mixin:
	// 	target.IncludeMixin(mixin)
	// default:
	// 	return value.Errorf(
	// 		value.TypeErrorClass,
	// 		"cannot include into an instance of %s: `%s`",
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

	body := bodyVal.(*BytecodeFunction)
	name := nameVal.(value.Symbol)

	methodContainer := vm.methodContainerValue()

	switch m := methodContainer.(type) {
	case *value.Class:
		if !m.CanOverride(name) {
			return value.NewCantOverrideASealedMethod(string(name.ToString()))
		}
		m.Methods[name] = body
	default:
		panic(fmt.Sprintf("invalid method container: %s", methodContainer.Inspect()))
	}

	vm.push(body)
	return nil
}

// Define a new mixin
func (vm *VM) defineMixin() (err value.Value) {
	constantNameVal := vm.pop()
	parentModuleVal := vm.pop()
	bodyVal := vm.pop()

	constantName := constantNameVal.(value.Symbol)
	var parentModule *value.ConstantContainer

	switch m := parentModuleVal.(type) {
	case *value.Class:
		parentModule = &m.ConstantContainer
	case *value.Module:
		parentModule = &m.ConstantContainer
	default:
		return value.NewIsNotModuleError(parentModuleVal.Inspect())
	}

	var mixin *value.Mixin
	var ok bool

	if mixinVal := parentModule.Constants.Get(constantName); mixinVal != nil {
		mixin, ok = mixinVal.(*value.Mixin)
		if !ok || !mixin.IsMixin() {
			return value.NewRedefinedConstantError(parentModuleVal.Inspect(), constantName.Inspect())
		}
	} else {
		mixin = value.NewMixin()
		parentModule.AddConstant(constantName, mixin)
	}

	switch body := bodyVal.(type) {
	case *BytecodeFunction:
		vm.executeMixinBody(mixin, body)
	case value.UndefinedType:
		vm.push(mixin)
	default:
		panic(fmt.Sprintf("expected undefined or a bytecode function as the mixin body, got: %s", bodyVal.Inspect()))
	}

	return nil
}

// Define a new module
func (vm *VM) defineModule() (err value.Value) {
	constantNameVal := vm.pop()
	parentModuleVal := vm.pop()
	bodyVal := vm.pop()

	constantName := constantNameVal.(value.Symbol)
	var parentModule *value.ConstantContainer

	switch m := parentModuleVal.(type) {
	case *value.Class:
		parentModule = &m.ConstantContainer
	case *value.Module:
		parentModule = &m.ConstantContainer
	default:
		return value.NewIsNotModuleError(parentModuleVal.Inspect())
	}

	var module *value.Module
	var ok bool

	if moduleVal := parentModule.Constants.Get(constantName); moduleVal != nil {
		module, ok = moduleVal.(*value.Module)
		if !ok {
			return value.NewRedefinedConstantError(parentModuleVal.Inspect(), constantName.Inspect())
		}
	} else {
		module = value.NewModule()
		parentModule.AddConstant(constantName, module)
	}

	switch body := bodyVal.(type) {
	case *BytecodeFunction:
		vm.executeModuleBody(module, body)
	case value.UndefinedType:
		vm.push(module)
	default:
		panic(fmt.Sprintf("expected undefined or a bytecode function as the module body, got: %s", bodyVal.Inspect()))
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
	}

	err := DefineGetter(container, name, false)
	if err != nil {
		return err
	}

	return nil
}

// Define a setter method
func (vm *VM) defineSetter() value.Value {
	name := vm.pop().(value.Symbol)

	var container *value.MethodContainer
	methodContainerValue := vm.methodContainerValue()

	switch methodContainer := methodContainerValue.(type) {
	case *value.Class:
		container = &methodContainer.MethodContainer
	}

	err := DefineSetter(container, name, false)
	if err != nil {
		return err
	}

	return nil
}

// Define a method alias
func (vm *VM) defineSingleton() value.Value {
	object := vm.pop()
	bodyVal := vm.pop()
	singletonClass := object.SingletonClass()

	if singletonClass == nil {
		return value.NewSingletonError(object.Inspect())
	}

	switch body := bodyVal.(type) {
	case *BytecodeFunction:
		vm.executeClassBody(singletonClass, body)
	case value.UndefinedType:
		vm.push(singletonClass)
	default:
		panic(fmt.Sprintf("expected undefined or a bytecode function as the class body, got: %s", bodyVal.Inspect()))
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
	}

	return nil
}

// Define a new class
func (vm *VM) defineClass() (err value.Value) {
	superclassVal := vm.pop()
	constantNameVal := vm.pop()
	parentModuleVal := vm.pop()
	bodyVal := vm.pop()
	flags := bitfield.BitField8FromInt(vm.readByte())

	constantName := constantNameVal.(value.Symbol)
	var parentModule *value.ConstantContainer

	switch mod := parentModuleVal.(type) {
	case *value.Class:
		parentModule = &mod.ConstantContainer
	case *value.Module:
		parentModule = &mod.ConstantContainer
	default:
		return value.NewIsNotModuleError(parentModuleVal.Inspect())
	}

	var class *value.Class
	var ok bool

	if classVal := parentModule.Constants.Get(constantName); classVal != nil {
		class, ok = classVal.(*value.Class)
		if !ok {
			return value.NewRedefinedConstantError(parentModuleVal.Inspect(), constantName.Inspect())
		}
		switch superclass := superclassVal.(type) {
		case *value.Class:
			if class.Parent != superclass {
				return value.NewSuperclassMismatchError(
					class.Name,
					class.Parent.Name,
					superclass.Name,
				)
			}
		case value.UndefinedType:
		default:
			return value.NewInvalidSuperclassError(superclass.Inspect())
		}

		if class.IsAbstract() {
			if flags.HasFlag(value.CLASS_SEALED_FLAG) {
				return value.NewModifierMismatchError(
					class.Inspect(),
					"sealed",
					false,
				)
			}
		} else if class.IsSealed() {
			if flags.HasFlag(value.CLASS_ABSTRACT_FLAG) {
				return value.NewModifierMismatchError(
					class.Inspect(),
					"abstract",
					false,
				)
			}
		} else {
			if flags.HasFlag(value.CLASS_ABSTRACT_FLAG) {
				return value.NewModifierMismatchError(
					class.Inspect(),
					"abstract",
					false,
				)
			}
			if flags.HasFlag(value.CLASS_SEALED_FLAG) {
				return value.NewModifierMismatchError(
					class.Inspect(),
					"sealed",
					false,
				)
			}
		}
	} else {
		class = value.NewClass()
		class.Flags = flags
		switch superclass := superclassVal.(type) {
		case *value.Class:
			if superclass.IsSealed() {
				return value.NewSealedClassError(string(constantName.ToString()), superclass.Inspect())
			}
			class.Parent = superclass
		case value.UndefinedType:
		default:
			return value.NewInvalidSuperclassError(superclass.Inspect())
		}
		parentModule.AddConstant(constantName, class)
	}

	switch body := bodyVal.(type) {
	case *BytecodeFunction:
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
			bytecode:   vm.bytecode,
			ip:         vm.ip,
			fp:         vm.fp,
			localCount: vm.localCount,
		},
	)
}

// set up the vm to execute a class body
func (vm *VM) executeClassBody(class value.Value, body *BytecodeFunction) {
	vm.createCurrentCallFrame()

	vm.bytecode = body
	vm.fp = vm.sp
	vm.ip = 0
	vm.localCount = 3
	// set class as `self`
	vm.push(class)
	// set class as constant container
	vm.push(class)
	// set class as method container
	vm.push(class)
}

// set up the vm to execute a mixin body
func (vm *VM) executeMixinBody(mixin value.Value, body *BytecodeFunction) {
	vm.createCurrentCallFrame()

	vm.bytecode = body
	vm.fp = vm.sp
	vm.ip = 0
	vm.localCount = 3
	// set mixin as `self`
	vm.push(mixin)
	// set mixin as constant container
	vm.push(mixin)
	// set mixin as method container
	vm.push(mixin)
}

// set up the vm to execute a module body
func (vm *VM) executeModuleBody(module value.Value, body *BytecodeFunction) {
	vm.createCurrentCallFrame()

	vm.bytecode = body
	vm.fp = vm.sp
	vm.ip = 0
	vm.localCount = 3
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
	val := vm.stack[vm.fp+index]
	if val == nil {
		return value.Nil
	}
	return val
}

// Set an upvalue.
func (vm *VM) setUpvalue(index int) {
	vm.setUpvalueValue(index, vm.peek())
}

// Set an upvalue.
func (vm *VM) setUpvalueValue(index int, val value.Value) {
	*vm.upvalues[index].location = val
}

// Read an upvalue.
func (vm *VM) getUpvalue(index int) {
	vm.push(vm.getUpvalueValue(index))
}

// Read an upvalue.
func (vm *VM) getUpvalueValue(index int) value.Value {
	return *vm.upvalues[index].location
}

// Closes all upvalues up to the given local slot.
func (vm *VM) closeUpvalues(lastToClose *value.Value) {
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

	val := constants.Get(symbol)
	if val == nil {
		return value.Errorf(value.NoConstantErrorClass, "%s doesn't have a constant named `%s`", mod.Inspect(), symbol.Inspect())
	}

	vm.push(val)
	return nil
}

// Get the iterator of the value on top of the stack.
func (vm *VM) getIterator() value.Value {
	val := vm.peek()
	result, err := vm.CallMethodByName(iteratorSymbol, val)
	if err != nil {
		return err
	}

	vm.replace(result)
	return nil
}

var nextSymbol = value.ToSymbol("next")
var stopIterationSymbol = value.ToSymbol("stop_iteration")
var iteratorSymbol = value.ToSymbol("iter")

// Drive the for..in loop.
func (vm *VM) forIn() value.Value {
	iterator := vm.pop()
	result, err := vm.CallMethodByName(nextSymbol, iterator)
	switch e := err.(type) {
	case value.Symbol:
		if e == stopIterationSymbol {
			vm.ip += int(vm.readUint16())
			return nil
		}
		return e
	case nil:
	default:
		return e
	}

	vm.push(result)
	vm.ip += 2
	return nil
}

var toStringSymbol = value.ToSymbol("to_string")

// Create a new string.
func (vm *VM) newString(dynamicElements int) value.Value {
	firstElementIndex := vm.sp - dynamicElements

	var buffer strings.Builder
	for i, elementVal := range vm.stack[firstElementIndex:vm.sp] {
		vm.stack[firstElementIndex+i] = nil

		switch element := elementVal.(type) {
		case value.String:
			buffer.WriteString(string(element))
		case value.Char:
			buffer.WriteRune(rune(element))
		case value.Float64:
			buffer.WriteString(string(element.ToString()))
		case value.Float32:
			buffer.WriteString(string(element.ToString()))
		case value.Float:
			buffer.WriteString(string(element.ToString()))
		case value.SmallInt:
			buffer.WriteString(string(element.ToString()))
		case *value.BigInt:
			buffer.WriteString(string(element.ToString()))
		case value.Int64:
			buffer.WriteString(string(element.ToString()))
		case value.Int32:
			buffer.WriteString(string(element.ToString()))
		case value.Int16:
			buffer.WriteString(string(element.ToString()))
		case value.Int8:
			buffer.WriteString(string(element.ToString()))
		case value.UInt64:
			buffer.WriteString(string(element.ToString()))
		case value.UInt32:
			buffer.WriteString(string(element.ToString()))
		case value.UInt16:
			buffer.WriteString(string(element.ToString()))
		case value.UInt8:
			buffer.WriteString(string(element.ToString()))
		case value.NilType:
		case value.Symbol:
			buffer.WriteString(string(element.ToString()))
		case *value.Regex:
			buffer.WriteString(string(element.ToString()))
		default:
			strVal, err := vm.CallMethodByName(toStringSymbol, elementVal)
			if err != nil {
				return err
			}
			str, ok := strVal.(value.String)
			if !ok {
				return value.NewCoerceError(value.StringClass, strVal.Class())
			}
			buffer.WriteString(string(str))
		}
	}
	vm.sp -= dynamicElements
	vm.push(value.String(buffer.String()))

	return nil
}

// Create a new symbol.
func (vm *VM) newSymbol(dynamicElements int) value.Value {
	firstElementIndex := vm.sp - dynamicElements

	var buffer strings.Builder
	for i, elementVal := range vm.stack[firstElementIndex:vm.sp] {
		vm.stack[firstElementIndex+i] = nil

		switch element := elementVal.(type) {
		case value.String:
			buffer.WriteString(string(element))
		case value.Char:
			buffer.WriteRune(rune(element))
		case value.Float64:
			buffer.WriteString(string(element.ToString()))
		case value.Float32:
			buffer.WriteString(string(element.ToString()))
		case value.Float:
			buffer.WriteString(string(element.ToString()))
		case value.SmallInt:
			buffer.WriteString(string(element.ToString()))
		case *value.BigInt:
			buffer.WriteString(string(element.ToString()))
		case value.Int64:
			buffer.WriteString(string(element.ToString()))
		case value.Int32:
			buffer.WriteString(string(element.ToString()))
		case value.Int16:
			buffer.WriteString(string(element.ToString()))
		case value.Int8:
			buffer.WriteString(string(element.ToString()))
		case value.UInt64:
			buffer.WriteString(string(element.ToString()))
		case value.UInt32:
			buffer.WriteString(string(element.ToString()))
		case value.UInt16:
			buffer.WriteString(string(element.ToString()))
		case value.UInt8:
			buffer.WriteString(string(element.ToString()))
		case value.NilType:
		case value.Symbol:
			buffer.WriteString(string(element.ToString()))
		case *value.Regex:
			buffer.WriteString(string(element.ToString()))
		default:
			strVal, err := vm.CallMethodByName(toStringSymbol, elementVal)
			if err != nil {
				return err
			}
			str, ok := strVal.(value.String)
			if !ok {
				return value.NewCoerceError(value.StringClass, strVal.Class())
			}
			buffer.WriteString(string(str))
		}
	}
	vm.sp -= dynamicElements
	vm.push(value.ToSymbol(buffer.String()))

	return nil
}

// Create a new regex.
func (vm *VM) newRegex(flagByte byte, dynamicElements int) value.Value {
	flags := bitfield.BitField8FromInt(flagByte)
	firstElementIndex := vm.sp - dynamicElements

	var buffer strings.Builder
	for i, elementVal := range vm.stack[firstElementIndex:vm.sp] {
		vm.stack[firstElementIndex+i] = nil

		switch element := elementVal.(type) {
		case value.String:
			buffer.WriteString(string(element))
		case value.Char:
			buffer.WriteRune(rune(element))
		case value.Float64:
			buffer.WriteString(string(element.ToString()))
		case value.Float32:
			buffer.WriteString(string(element.ToString()))
		case value.Float:
			buffer.WriteString(string(element.ToString()))
		case value.SmallInt:
			buffer.WriteString(string(element.ToString()))
		case *value.BigInt:
			buffer.WriteString(string(element.ToString()))
		case value.Int64:
			buffer.WriteString(string(element.ToString()))
		case value.Int32:
			buffer.WriteString(string(element.ToString()))
		case value.Int16:
			buffer.WriteString(string(element.ToString()))
		case value.Int8:
			buffer.WriteString(string(element.ToString()))
		case value.UInt64:
			buffer.WriteString(string(element.ToString()))
		case value.UInt32:
			buffer.WriteString(string(element.ToString()))
		case value.UInt16:
			buffer.WriteString(string(element.ToString()))
		case value.UInt8:
			buffer.WriteString(string(element.ToString()))
		case value.NilType:
		case value.Symbol:
			buffer.WriteString(string(element.ToString()))
		case *value.Regex:
			buffer.WriteString(string(element.ToStringWithFlags()))
		default:
			strVal, err := vm.CallMethodByName(toStringSymbol, elementVal)
			if err != nil {
				return err
			}
			str, ok := strVal.(value.String)
			if !ok {
				return value.NewCoerceError(value.StringClass, strVal.Class())
			}
			buffer.WriteString(string(str))
		}
	}
	vm.sp -= dynamicElements
	re, err := value.CompileRegex(buffer.String(), flags)
	if err != nil {
		return value.NewError(value.RegexCompileErrorClass, err.Error())
	}

	vm.push(re)
	return nil
}

// Create a new hashset.
func (vm *VM) newHashSet(dynamicElements int) value.Value {
	firstElementIndex := vm.sp - dynamicElements
	capacity := vm.stack[firstElementIndex-2]
	baseSet := vm.stack[firstElementIndex-1]
	var newSet *value.HashSet

	var additionalCapacity int

	switch capacity.(type) {
	case value.UndefinedType:
	default:
		c, ok := value.ToGoInt(capacity)
		if c == -1 && !ok {
			return value.NewTooLargeCapacityError(capacity.Inspect())
		}
		if c < 0 {
			return value.NewNegativeCapacityError(capacity.Inspect())
		}
		if !ok {
			return value.NewCapacityTypeError(capacity.Inspect())
		}
		additionalCapacity = c
	}

	switch m := baseSet.(type) {
	case value.UndefinedType:
		newSet = value.NewHashSet(dynamicElements + additionalCapacity)
	case *value.HashSet:
		newSet = value.NewHashSet(m.Capacity() + additionalCapacity)
		err := HashSetCopy(vm, newSet, m)
		if err != nil {
			return err
		}
	}

	for i := firstElementIndex; i < vm.sp; i++ {
		val := vm.stack[i]
		err := HashSetAppendWithMaxLoad(vm, newSet, val, 1)
		if err != nil {
			return err
		}
	}
	vm.popN(dynamicElements + 2)

	vm.push(newSet)
	return nil
}

// Create a new hashmap.
func (vm *VM) newHashMap(dynamicElements int) value.Value {
	firstElementIndex := vm.sp - (dynamicElements * 2)
	capacity := vm.stack[firstElementIndex-2]
	baseMap := vm.stack[firstElementIndex-1]
	var newMap *value.HashMap

	var additionalCapacity int

	switch capacity.(type) {
	case value.UndefinedType:
	default:
		c, ok := value.ToGoInt(capacity)
		if c == -1 && !ok {
			return value.NewTooLargeCapacityError(capacity.Inspect())
		}
		if c < 0 {
			return value.NewNegativeCapacityError(capacity.Inspect())
		}
		if !ok {
			return value.NewCapacityTypeError(capacity.Inspect())
		}
		additionalCapacity = c
	}

	switch m := baseMap.(type) {
	case value.UndefinedType:
		newMap = value.NewHashMap(dynamicElements + additionalCapacity)
	case *value.HashMap:
		newMap = value.NewHashMap(m.Capacity() + additionalCapacity)
		err := HashMapCopy(vm, newMap, m)
		if err != nil {
			return err
		}
	}

	for i := firstElementIndex; i < vm.sp; i += 2 {
		key := vm.stack[i]
		val := vm.stack[i+1]
		err := HashMapSetWithMaxLoad(vm, newMap, key, val, 1)
		if err != nil {
			return err
		}
	}
	vm.popN((dynamicElements * 2) + 2)

	vm.push(newMap)
	return nil
}

// Create a new hash record.
func (vm *VM) newHashRecord(dynamicElements int) value.Value {
	firstElementIndex := vm.sp - (dynamicElements * 2)
	baseMap := vm.stack[firstElementIndex-1]
	var newRecord *value.HashRecord

	switch m := baseMap.(type) {
	case value.UndefinedType:
		newRecord = value.NewHashRecord(dynamicElements)
	case *value.HashRecord:
		newRecord = value.NewHashRecord(m.Length())
		HashRecordCopy(vm, newRecord, m)
	}

	for i := firstElementIndex; i < vm.sp; i += 2 {
		key := vm.stack[i]
		val := vm.stack[i+1]
		HashRecordSetWithMaxLoad(vm, newRecord, key, val, 1)
	}
	vm.popN((dynamicElements * 2) + 1)

	vm.push(newRecord)
	return nil
}

// Create a new array list.
func (vm *VM) newArrayList(dynamicElements int) value.Value {
	firstElementIndex := vm.sp - dynamicElements
	capacity := vm.stack[firstElementIndex-2]
	baseList := vm.stack[firstElementIndex-1]
	var newArrayList value.ArrayList

	var additionalCapacity int

	switch capacity.(type) {
	case value.UndefinedType:
	default:
		c, ok := value.ToGoInt(capacity)
		if c == -1 && !ok {
			return value.NewTooLargeCapacityError(capacity.Inspect())
		}
		if c < 0 {
			return value.NewNegativeCapacityError(capacity.Inspect())
		}
		if !ok {
			return value.NewCapacityTypeError(capacity.Inspect())
		}
		additionalCapacity = c
	}

	switch l := baseList.(type) {
	case value.UndefinedType:
		newArrayList = make(value.ArrayList, 0, dynamicElements+additionalCapacity)
	case *value.ArrayList:
		newArrayList = make(value.ArrayList, 0, cap(*l)+additionalCapacity)
		newArrayList = append(newArrayList, *l...)
	}

	newArrayList = append(newArrayList, vm.stack[firstElementIndex:vm.sp]...)
	vm.popN(dynamicElements + 2)

	vm.push(&newArrayList)
	return nil
}

// Create a new range.
func (vm *VM) newRange() {
	flag := vm.readByte()
	var newRange value.Value

	switch flag {
	case bytecode.CLOSED_RANGE_FLAG:
		to := vm.pop()
		from := vm.pop()
		newRange = value.NewClosedRange(from, to)
	case bytecode.OPEN_RANGE_FLAG:
		to := vm.pop()
		from := vm.pop()
		newRange = value.NewOpenRange(from, to)
	case bytecode.LEFT_OPEN_RANGE_FLAG:
		to := vm.pop()
		from := vm.pop()
		newRange = value.NewLeftOpenRange(from, to)
	case bytecode.RIGHT_OPEN_RANGE_FLAG:
		to := vm.pop()
		from := vm.pop()
		newRange = value.NewRightOpenRange(from, to)
	case bytecode.BEGINLESS_CLOSED_RANGE_FLAG:
		newRange = value.NewBeginlessClosedRange(vm.pop())
	case bytecode.BEGINLESS_OPEN_RANGE_FLAG:
		newRange = value.NewBeginlessOpenRange(vm.pop())
	case bytecode.ENDLESS_CLOSED_RANGE_FLAG:
		newRange = value.NewEndlessClosedRange(vm.pop())
	case bytecode.ENDLESS_OPEN_RANGE_FLAG:
		newRange = value.NewEndlessOpenRange(vm.pop())
	default:
		panic(fmt.Sprintf("invalid range flag: %#v", flag))
	}

	vm.push(newRange)
}

// Create a new arrayTuple.
func (vm *VM) newArrayTuple(dynamicElements int) {
	firstElementIndex := vm.sp - dynamicElements
	baseArrayTuple := vm.stack[firstElementIndex-1]
	var newArrayTuple value.ArrayTuple

	switch t := baseArrayTuple.(type) {
	case value.UndefinedType:
		newArrayTuple = make(value.ArrayTuple, 0, dynamicElements)
	case *value.ArrayTuple:
		newArrayTuple = make(value.ArrayTuple, 0, len(*t)+dynamicElements)
		newArrayTuple = append(newArrayTuple, *t...)
	}

	newArrayTuple = append(newArrayTuple, vm.stack[firstElementIndex:vm.sp]...)
	vm.popN(dynamicElements + 1)

	vm.push(&newArrayTuple)
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
	vm.localCount += count
}

// Push an element on top of the value stack.
func (vm *VM) push(val value.Value) {
	vm.stack[vm.sp] = val
	vm.sp++
}

// Push an element on top of the value stack.
func (vm *VM) swap() {
	tmp := vm.stack[vm.sp-2]
	vm.stack[vm.sp-2] = vm.stack[vm.sp-1]
	vm.stack[vm.sp-1] = tmp
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

// Pop all values on the stack leaving only slots for locals
func (vm *VM) popAll() {
	vm.popN(vm.sp - vm.localCount - 1)
}

// Pop n elements off the value stack.
func (vm *VM) popN(n int) {
	if vm.sp-n < 0 {
		panic("tried to pop more elements than are available on the value stack!")
	}

	for i := vm.sp - 1; i >= vm.sp-n; i-- {
		vm.stack[i] = nil
	}
	vm.sp -= n
}

// Pop one element off the value stack skipping the first one.
func (vm *VM) popSkipOne() {
	if vm.sp-2 < 0 {
		panic("tried to pop more elements than are available on the value stack!")
	}

	vm.sp--
	vm.stack[vm.sp-1] = vm.stack[vm.sp]
}

// Pop n elements off the value stack skipping the first one.
func (vm *VM) popNSkipOne(n int) {
	if vm.sp-n-1 < 0 {
		panic("tried to pop more elements than are available on the value stack!")
	}

	vm.stack[vm.sp-n-1] = vm.stack[vm.sp-1]
	for i := vm.sp - 1; i >= vm.sp-n; i-- {
		vm.stack[i] = nil
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

type unaryOperationFunc func(val value.Value) value.Value

func (vm *VM) unaryOperation(fn unaryOperationFunc, methodName value.Symbol) value.Value {
	operand := vm.peek()
	result := fn(operand)
	if result != nil {
		vm.replace(result)
		return nil
	}

	er := vm.callMethodOnStackByName(methodName, 0)
	if er != nil {
		return er
	}
	return nil
}

// Increment the element on top of the stack
func (vm *VM) increment() (err value.Value) {
	return vm.unaryOperation(value.Increment, symbol.OpIncrement)
}

// Decrement the element on top of the stack
func (vm *VM) decrement() (err value.Value) {
	return vm.unaryOperation(value.Decrement, symbol.OpDecrement)
}

// Negate the element on top of the stack
func (vm *VM) negate() (err value.Value) {
	return vm.unaryOperation(value.Negate, symbol.OpNegate)
}

// Perform unary plus on the element on top of the stack
func (vm *VM) unaryPlus() (err value.Value) {
	return vm.unaryOperation(value.UnaryPlus, symbol.OpUnaryPlus)
}

// Preform bitwise not on the element on top of the stack
func (vm *VM) bitwiseNot() (err value.Value) {
	return vm.unaryOperation(value.BitwiseNot, symbol.OpBitwiseNot)
}

func (vm *VM) appendAt() value.Value {
	val := vm.pop()
	key := vm.pop()
	collection := vm.peek()

	i, ok := value.ToGoInt(key)

	switch c := collection.(type) {
	case *value.ArrayTuple:
		l := len(*c)
		if !ok {
			if i == -1 {
				return value.NewIndexOutOfRangeError(key.Inspect(), l)
			}
			return value.NewCoerceError(value.IntClass, key.Class())
		}

		if i < 0 {
			return value.NewNegativeIndicesInCollectionLiteralsError(fmt.Sprint(i))
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
				return value.NewIndexOutOfRangeError(key.Inspect(), l)
			}
			return value.NewCoerceError(value.IntClass, key.Class())
		}

		if i < 0 {
			return value.NewNegativeIndicesInCollectionLiteralsError(fmt.Sprint(i))
		}

		if i >= l {
			newElementsCount := (i + 1) - l
			c.Expand(newElementsCount)
		}

		(*c)[i] = val
	default:
		panic(fmt.Sprintf("cannot APPEND_AT to: %#v", collection))
	}

	return nil
}

func (vm *VM) subscriptSet() value.Value {
	val := vm.peek()
	key := vm.peekAt(1)
	collection := vm.peekAt(2)

	result, err := value.SubscriptSet(collection, key, val)
	if err != nil {
		return err
	}
	if result != nil {
		vm.popN(2)
		vm.replace(result)
		return nil
	}

	er := vm.callMethodOnStackByName(symbol.OpSubscriptSet, 2)
	if er != nil {
		return er
	}
	return nil
}

func (vm *VM) isA() (err value.Value) {
	classVal := vm.pop()
	val := vm.peek()

	switch class := classVal.(type) {
	case *value.Class:
		vm.replace(value.ToElkBool(value.IsA(val, class)))
	default:
		vm.pop()
		return value.NewIsNotClassOrMixinError(class.Inspect())
	}

	return nil
}

func (vm *VM) instanceOf() (err value.Value) {
	classVal := vm.pop()
	val := vm.peek()

	class, ok := classVal.(*value.Class)
	if !ok || class.IsMixin() || class.IsMixinProxy() {
		vm.pop()
		return value.NewIsNotClassError(classVal.Inspect())
	}

	vm.replace(value.ToElkBool(value.InstanceOf(val, class)))
	return nil
}

type binaryOperationWithoutErrFunc func(left value.Value, right value.Value) value.Value

func (vm *VM) binaryOperationWithoutErr(fn binaryOperationWithoutErrFunc, methodName value.Symbol) (err value.Value) {
	right := vm.peek()
	left := vm.peekAt(1)

	result := fn(left, right)
	if result != nil {
		vm.pop()
		vm.replace(result)
		return nil
	}

	er := vm.callMethodOnStackByName(methodName, 1)
	if er != nil {
		return er
	}

	return nil
}

func (vm *VM) negatedBinaryOperationWithoutErr(fn binaryOperationWithoutErrFunc, methodName value.Symbol) (err value.Value) {
	right := vm.peek()
	left := vm.peekAt(1)

	result := fn(left, right)
	if result != nil {
		vm.pop()
		vm.replace(result)
		return nil
	}

	er := vm.callMethodOnStackByName(methodName, 1)
	if er != nil {
		return er
	}
	vm.replace(value.ToNotBool(vm.peek()))

	return nil
}

type binaryOperationFunc func(left value.Value, right value.Value) (value.Value, *value.Error)

func (vm *VM) binaryOperation(fn binaryOperationFunc, methodName value.Symbol) value.Value {
	right := vm.peek()
	left := vm.peekAt(1)

	result, err := fn(left, right)
	if err != nil {
		return err
	}
	if result != nil {
		vm.pop()
		vm.replace(result)
		return nil
	}

	er := vm.callMethodOnStackByName(methodName, 1)
	if er != nil {
		return er
	}
	return nil
}

// Perform a bitwise AND and push the result to the stack.
func (vm *VM) bitwiseAnd() (err value.Value) {
	return vm.binaryOperation(value.BitwiseAnd, symbol.OpAnd)
}

// Perform a bitwise AND NOT and push the result to the stack.
func (vm *VM) bitwiseAndNot() (err value.Value) {
	return vm.binaryOperation(value.BitwiseAndNot, symbol.OpAndNot)
}

// Get the value under the given key and push the result to the stack.
func (vm *VM) subscript() (err value.Value) {
	return vm.binaryOperation(value.Subscript, symbol.OpSubscript)
}

// Perform a bitwise OR and push the result to the stack.
func (vm *VM) bitwiseOr() (err value.Value) {
	return vm.binaryOperation(value.BitwiseOr, symbol.OpOr)
}

// Perform a bitwise XOR and push the result to the stack.
func (vm *VM) bitwiseXor() (err value.Value) {
	return vm.binaryOperation(value.BitwiseXor, symbol.OpXor)
}

// Perform a comparison and push the result to the stack.
func (vm *VM) compare() (err value.Value) {
	return vm.binaryOperation(value.Compare, symbol.OpSpaceship)
}

// Perform modulo and push the result to the stack.
func (vm *VM) modulo() (err value.Value) {
	return vm.binaryOperation(value.Modulo, symbol.OpModulo)
}

// Check whether two top elements on the stack are equal and push the result to the stack.
func (vm *VM) equal() (err value.Value) {
	return vm.callEqualityOperator(value.Equal, symbol.OpEqual)
}

func (vm *VM) callEqualityOperator(fn binaryOperationWithoutErrFunc, methodName value.Symbol) (err value.Value) {
	right := vm.peek()
	left := vm.peekAt(1)

	result := fn(left, right)
	if result != nil {
		vm.pop()
		vm.replace(result)
		return nil
	}

	self := vm.stack[vm.sp-2]
	class := self.DirectClass()
	method := class.LookupMethod(methodName)
	if method == nil {
		vm.push(value.ToElkBool(left == right))
		return nil
	}

	return vm.callMethodOnStack(method, 1)
}

func (vm *VM) callNegatedEqualityOperator(fn binaryOperationWithoutErrFunc, methodName value.Symbol) (err value.Value) {
	right := vm.peek()
	left := vm.peekAt(1)

	result := fn(left, right)
	if result != nil {
		vm.pop()
		vm.replace(result)
		return nil
	}

	self := vm.stack[vm.sp-2]
	class := self.DirectClass()
	method := class.LookupMethod(methodName)
	if method == nil {
		vm.push(value.ToElkBool(left != right))
		return nil
	}

	err = vm.callMethodOnStack(method, 1)
	if err != nil {
		return err
	}

	vm.replace(value.ToNotBool(vm.peek()))
	return nil
}

// Check whether two top elements on the stack are not and equal push the result to the stack.
func (vm *VM) notEqual() (err value.Value) {
	return vm.callNegatedEqualityOperator(value.NotEqual, symbol.OpEqual)
}

// Check whether two top elements on the stack are equal and push the result to the stack.
func (vm *VM) laxEqual() (err value.Value) {
	return vm.callEqualityOperator(value.LaxEqual, symbol.OpLaxEqual)
}

// Check whether two top elements on the stack are not and equal push the result to the stack.
func (vm *VM) laxNotEqual() (err value.Value) {
	return vm.callNegatedEqualityOperator(value.LaxNotEqual, symbol.OpLaxEqual)
}

// Check whether two top elements on the stack are strictly equal push the result to the stack.
func (vm *VM) strictEqual() (err value.Value) {
	return vm.binaryOperationWithoutErr(value.StrictEqual, symbol.OpStrictEqual)
}

// Check whether two top elements on the stack are strictly not equal push the result to the stack.
func (vm *VM) strictNotEqual() (err value.Value) {
	return vm.negatedBinaryOperationWithoutErr(value.StrictNotEqual, symbol.OpStrictEqual)
}

// Check whether the first operand is greater than the second and push the result to the stack.
func (vm *VM) greaterThan() (err value.Value) {
	return vm.binaryOperation(value.GreaterThan, symbol.OpGreaterThan)
}

// Check whether the first operand is greater than or equal to the second and push the result to the stack.
func (vm *VM) greaterThanEqual() (err value.Value) {
	return vm.binaryOperation(value.GreaterThanEqual, symbol.OpGreaterThanEqual)
}

// Check whether the first operand is less than the second and push the result to the stack.
func (vm *VM) lessThan() (err value.Value) {
	return vm.binaryOperation(value.LessThan, symbol.OpLessThan)
}

// Check whether the first operand is less than or equal to the second and push the result to the stack.
func (vm *VM) lessThanEqual() (err value.Value) {
	return vm.binaryOperation(value.LessThanEqual, symbol.OpLessThanEqual)
}

// Perform a left bitshift and push the result to the stack.
func (vm *VM) leftBitshift() (err value.Value) {
	return vm.binaryOperation(value.LeftBitshift, symbol.OpLeftBitshift)
}

// Perform a logical left bitshift and push the result to the stack.
func (vm *VM) logicalLeftBitshift() (err value.Value) {
	return vm.binaryOperation(value.LogicalLeftBitshift, symbol.OpLogicalLeftBitshift)
}

// Perform a right bitshift and push the result to the stack.
func (vm *VM) rightBitshift() (err value.Value) {
	return vm.binaryOperation(value.RightBitshift, symbol.OpRightBitshift)
}

// Perform a logical right bitshift and push the result to the stack.
func (vm *VM) logicalRightBitshift() (err value.Value) {
	return vm.binaryOperation(value.LogicalRightBitshift, symbol.OpLogicalRightBitshift)
}

// Add two operands together and push the result to the stack.
func (vm *VM) add() (err value.Value) {
	return vm.binaryOperation(value.Add, symbol.OpAdd)
}

// Subtract two operands and push the result to the stack.
func (vm *VM) subtract() (err value.Value) {
	return vm.binaryOperation(value.Subtract, symbol.OpSubtract)
}

// Multiply two operands together and push the result to the stack.
func (vm *VM) multiply() (err value.Value) {
	return vm.binaryOperation(value.Multiply, symbol.OpMultiply)
}

// Divide two operands and push the result to the stack.
func (vm *VM) divide() (err value.Value) {
	return vm.binaryOperation(value.Divide, symbol.OpDivide)
}

// Exponentiate two operands and push the result to the stack.
func (vm *VM) exponentiate() (err value.Value) {
	return vm.binaryOperation(value.Exponentiate, symbol.OpExponentiate)
}

// Throw an error and attempt to find code
// that catches it.
func (vm *VM) throw(err value.Value) {
	vm.rethrow(err, value.String(vm.BuildStackTrace()))
}

func (vm *VM) rethrow(err value.Value, stackTrace value.String) {
	for {
		var foundCatch *CatchEntry

		for _, catchEntry := range vm.bytecode.CatchEntries {
			if !catchEntry.Finally && vm.ip > catchEntry.From && vm.ip <= catchEntry.To {
				foundCatch = catchEntry
				break
			}
		}

		if foundCatch != nil {
			vm.ip = foundCatch.JumpAddress
			vm.push(stackTrace)
			vm.push(err)
			return
		}

		if vm.mode == singleFunctionCallMode || len(vm.callFrames) < 1 {
			vm.mode = errorMode
			vm.errStackTrace = string(stackTrace)
			vm.push(err)
			return
		}

		vm.restoreLastFrame()
	}
}

func (vm *VM) jumpToFinallyForReturn() bool {
	catchEntry := vm.findFinallyCatchEntry()
	if catchEntry == nil {
		return false
	}

	// execute finally
	vm.ip = catchEntry.JumpAddress
	return true
}

func (vm *VM) jumpToFinallyForBreakOrContinue() bool {
	catchEntry := vm.findFinallyCatchEntry()
	if catchEntry == nil {
		return false
	}

	// execute finally
	vm.ip = catchEntry.JumpAddress + 4 // skip NIL, JUMP, offsetByte1, offsetByte2
	return true
}

func (vm *VM) findFinallyCatchEntry() *CatchEntry {
	for _, catchEntry := range vm.bytecode.CatchEntries {
		if catchEntry.Finally && vm.ip > catchEntry.From && vm.ip <= catchEntry.To {
			return catchEntry
		}
	}

	return nil
}
