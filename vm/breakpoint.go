package vm

import "github.com/elk-language/elk/value"

type BreakpointHandler interface {
	RunBreakpoint(thread *Thread, breakpointContext value.Value)
}

var BREAKPOINT_HANDLER BreakpointHandler
