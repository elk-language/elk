package vm

import "github.com/elk-language/elk/value"

// VM state
type state uint8

const (
	idleState state = iota
	runningState
	errorState // the VM stopped after encountering an uncaught error
	awaitState
	terminatedState
)

var stateSymbols = [...]value.Symbol{
	idleState:       value.ToSymbol("idle"),
	runningState:    value.ToSymbol("running"),
	errorState:      value.ToSymbol("error"),
	terminatedState: value.ToSymbol("terminated"),
}
