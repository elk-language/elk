package value

import (
	"fmt"
)

// Contains details like the number of arguments
// or the method name of a particular call site.
type CallSiteInfo struct {
	Name          Symbol
	ArgumentCount int
	Cache         [3]CallCache
}

type CallCache struct {
	Class  *Class
	Method Method
}

// Create a new CallSiteInfo.
func NewCallSiteInfo(methodName Symbol, argCount int) *CallSiteInfo {
	return &CallSiteInfo{
		Name:          methodName,
		ArgumentCount: argCount,
	}
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

func (c *CallSiteInfo) Copy() Reference {
	return c
}

func (c *CallSiteInfo) Inspect() string {
	return fmt.Sprintf(
		"CallSiteInfo{&: %p, name: %s, argument_count: %d}",
		c,
		c.Name.Inspect(),
		c.ArgumentCount,
	)
}

func (c *CallSiteInfo) Error() string {
	return c.Inspect()
}
