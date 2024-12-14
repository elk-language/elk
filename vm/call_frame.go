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
		CALL_STACK_SIZE = 1024 // 1024 frames by default
		return
	}

	CALL_STACK_SIZE = val
}

const CallFrameSize = unsafe.Sizeof(CallFrame{})

// Contains the data of a single function call.
type CallFrame struct {
	bytecode   *BytecodeFunction
	ip         uintptr // Instruction pointer - points to the next bytecode instruction for this frame
	fp         uintptr // Frame pointer -- points to the offset on the value stack where the current frame start
	localCount int
	upvalues   []*Upvalue
}

func (c CallFrame) Name() value.Symbol {
	return c.bytecode.Name()
}

func (c *CallFrame) ipIndex() int {
	return int(
		uintptr(unsafe.Pointer(c.ip)) -
			uintptr(unsafe.Pointer(&c.bytecode.Instructions[0])),
	)
}

func (c CallFrame) LineNumber() int {
	return c.bytecode.GetLineNumber(c.ipIndex() - 1)
}

func (c CallFrame) FileName() string {
	return c.bytecode.FileName()
}
