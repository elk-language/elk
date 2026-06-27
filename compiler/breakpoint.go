package compiler

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

type BytecodeBreakpointContext struct {
	lastLocalIndex     int
	maxLocalIndex      int
	patternNesting     int
	upvalues           bytecodeUpvalues
	scopes             bytecodeScopes
	Location           *position.Location
	TypecheckerContext value.Reference
	value.ValueBase
}

var _ value.Reference = &BytecodeBreakpointContext{}

func (c *BytecodeBreakpointContext) Copy() value.Reference {
	return &BytecodeBreakpointContext{
		lastLocalIndex:     c.lastLocalIndex,
		maxLocalIndex:      c.maxLocalIndex,
		patternNesting:     c.patternNesting,
		upvalues:           c.upvalues,
		scopes:             c.scopes,
		Location:           c.Location,
		TypecheckerContext: c.TypecheckerContext,
	}
}

func (c *BytecodeBreakpointContext) ToValue() value.Value {
	return value.Ref(c)
}

func (c *BytecodeBreakpointContext) Inspect() string {
	return "<bytecode breakpoint context>"
}
