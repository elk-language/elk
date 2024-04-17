package vm

import (
	"github.com/elk-language/elk/config"
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
