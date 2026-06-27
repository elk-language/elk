package vm

import (
	"unsafe"

	"github.com/elk-language/elk/config"
	"github.com/elk-language/elk/value"
)

var CALL_STACK_SIZE int

func init() {
	val, ok := config.IntFromEnvVar("ELK_CALL_STACK_SIZE")
	if !ok {
		CALL_STACK_SIZE = 74_000 / int(CallFrameSize) // 74KB by default
		return
	}

	CALL_STACK_SIZE = val / int(CallFrameSize)
}

const CallFrameSize = unsafe.Sizeof(CallFrame{})

// Contains the data of a single function call.
type CallFrame struct {
	upvalues        []*Upvalue
	bytecode        *BytecodeFunction
	ip              uintptr // Instruction pointer - points to the next bytecode instruction for this frame. In native frames stores the file name as a Symbol.
	fp              uintptr // Frame pointer -- points to the offset on the value stack where the current frame start. In native frames stores the function name as a Symbol.
	localCount      int     // Number of local variables. In native frames stores the line number.
	tailCallCounter int
	stopVM          bool
	sentinel        bool
	isNative        bool
}

func makeSentinelCallFrame() CallFrame {
	return CallFrame{
		sentinel: true,
	}
}

func makeNativeCallFrame(fileName, funcName value.Symbol, lineNumber, tailCallCounter int) CallFrame {
	return CallFrame{
		ip:              uintptr(fileName),
		fp:              uintptr(funcName),
		localCount:      lineNumber,
		tailCallCounter: tailCallCounter,
		isNative:        true,
	}
}

func (c *CallFrame) SetNativeLineNumber(lineNumber int) {
	c.localCount = lineNumber
}

func (c *CallFrame) IsNative() bool {
	return c.isNative
}

func (c *CallFrame) FuncName() value.Symbol {
	if c.isNative {
		return value.Symbol(c.fp)
	}

	return c.bytecode.Name()
}

func (c *CallFrame) LineNumber() int {
	if c.isNative {
		return c.localCount
	}

	return c.bytecode.GetLineNumber(c.ipIndex() - 1)
}

func (c *CallFrame) FileName() string {
	if c.isNative {
		return value.Symbol(c.ip).String()
	}

	return c.bytecode.FileName()
}

func (cf *CallFrame) ToCallFrameObject() value.CallFrame {
	return value.CallFrame{
		LineNumber:      cf.LineNumber(),
		FileName:        cf.FileName(),
		FuncName:        cf.FuncName().String(),
		TailCallCounter: cf.tailCallCounter,
	}
}

func (c *CallFrame) ipIndex() int {
	return int(
		c.ip -
			uintptr(unsafe.Pointer(&c.bytecode.Instructions[0])),
	)
}

// Std::CallFrame
func initCallFrame() {
	// Instance methods
	c := &value.CallFrameClass.MethodContainer
	Def(
		c,
		"to_string",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)
	Def(
		c,
		"func_name",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.Ref(value.String(self.FuncName)), value.Undefined
		},
	)
	Def(
		c,
		"file_name",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.Ref(value.String(self.FileName)), value.Undefined
		},
	)
	Def(
		c,
		"line_number",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.SmallInt(self.LineNumber).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"tail_calls",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.SmallInt(self.TailCallCounter).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

}
