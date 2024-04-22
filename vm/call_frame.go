package vm

import (
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

// Contains the data of a single function call.
type CallFrame struct {
	bytecode   *BytecodeMethod
	ip         int // Instruction pointer - points to the next bytecode instruction for this frame
	fp         int // Frame pointer -- points to the offset on the value stack where the current frame start
	localCount int
}

func (c CallFrame) Name() value.Symbol {
	return c.bytecode.Name()
}

func (c CallFrame) LineNumber() int {
	return c.bytecode.GetLineNumber(c.ip - 1)
}

func (c CallFrame) FileName() string {
	return c.bytecode.FileName()
}
