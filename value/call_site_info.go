package value

import (
	"fmt"
)

// Contains details like the number of arguments
// or the method name of a particular call site.
type CallSiteInfo struct {
	Name           Symbol
	ArgumentCount  int
	NamedArguments []Symbol
}

// Create a new CallSiteInfo.
func NewCallSiteInfo(methodName Symbol, argCount int, namedArgs []Symbol) *CallSiteInfo {
	return &CallSiteInfo{
		Name:           methodName,
		ArgumentCount:  argCount,
		NamedArguments: namedArgs,
	}
}

func (c *CallSiteInfo) PositionalArgumentCount() int {
	return c.ArgumentCount - c.NamedArgumentCount()
}

func (c *CallSiteInfo) NamedArgumentCount() int {
	return len(c.NamedArguments)
}

func (*CallSiteInfo) Class() *Class {
	return nil
}

func (*CallSiteInfo) DirectClass() *Class {
	return nil
}

func (*CallSiteInfo) SingletonClass() *Class {
	return nil
}

func (*CallSiteInfo) InstanceVariables() SymbolMap {
	return nil
}

func (c *CallSiteInfo) Inspect() string {
	if c.NamedArguments == nil {
		return fmt.Sprintf(
			"CallSiteInfo{name: %s, argument_count: %d}",
			c.Name.Inspect(),
			c.ArgumentCount,
		)
	}

	return fmt.Sprintf(
		"CallSiteInfo{name: %s, argument_count: %d, named_arguments: %s}",
		c.Name.Inspect(),
		c.ArgumentCount,
		InspectSlice(c.NamedArguments),
	)
}
