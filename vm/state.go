package vm

import "github.com/elk-language/elk/value"

// VM State
type State uint8

const (
	IdleState State = iota
	RunningState
	ErrorState // the VM stopped after encountering an uncaught error
	AwaitState
	TerminatedState
	PanicState
)

var stateSymbols = [...]value.Symbol{
	IdleState:       value.ToSymbol("idle"),
	RunningState:    value.ToSymbol("running"),
	ErrorState:      value.ToSymbol("error"),
	AwaitState:      value.ToSymbol("await"),
	TerminatedState: value.ToSymbol("terminated"),
	PanicState:      value.ToSymbol("panic"),
}
