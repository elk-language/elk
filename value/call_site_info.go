package value

import (
	"fmt"
	"strings"
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

func (*CallSiteInfo) IsFrozen() bool {
	return true
}

func (*CallSiteInfo) SetFrozen() {}

func (*CallSiteInfo) InstanceVariables() SimpleSymbolMap {
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

	var builder strings.Builder
	builder.WriteString(
		fmt.Sprintf(
			"CallSiteInfo{name: %s, argument_count: %d, named_arguments: [",
			c.Name.Inspect(),
			c.ArgumentCount,
		),
	)

	for i, name := range c.NamedArguments {
		if i != 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(name.Inspect())
	}

	builder.WriteString("]}")

	return builder.String()
}
