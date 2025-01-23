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
	ip              uintptr // Instruction pointer - points to the next bytecode instruction for this frame
	fp              uintptr // Frame pointer -- points to the offset on the value stack where the current frame start
	localCount      int
	tailCallCounter int
	stopVM          bool
	sentinel        bool
}

func makeSentinelCallFrame() CallFrame {
	return CallFrame{
		sentinel: true,
	}
}

func (c *CallFrame) Name() value.Symbol {
	return c.bytecode.Name()
}

func (cf *CallFrame) ToCallFrameObject() value.CallFrame {
	return value.CallFrame{
		LineNumber:      cf.LineNumber(),
		FileName:        cf.FileName(),
		FuncName:        cf.Name().String(),
		TailCallCounter: cf.tailCallCounter,
	}
}

func (c *CallFrame) ipIndex() int {
	return int(
		uintptr(unsafe.Pointer(c.ip)) -
			uintptr(unsafe.Pointer(&c.bytecode.Instructions[0])),
	)
}

func (c *CallFrame) LineNumber() int {
	return c.bytecode.GetLineNumber(c.ipIndex() - 1)
}

func (c *CallFrame) FileName() string {
	return c.bytecode.FileName()
}

// Std::CallFrame
func initCallFrame() {
	// Instance methods
	c := &value.CallFrameClass.MethodContainer
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)
	Def(
		c,
		"func_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.Ref(value.String(self.FuncName)), value.Undefined
		},
	)
	Def(
		c,
		"file_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.Ref(value.String(self.FileName)), value.Undefined
		},
	)
	Def(
		c,
		"line_number",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.SmallInt(self.LineNumber).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"tail_calls",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.SmallInt(self.TailCallCounter).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.CallFrame)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

}
