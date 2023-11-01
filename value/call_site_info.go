package value

import "fmt"

// Contains details like the number of arguments
// or the method name of a particular call site.
type CallSiteInfo struct {
	MethodName    Symbol
	ArgumentCount int
}

// Create a new CallSiteInfo.
func NewCallSiteInfo(methodName Symbol, argCount int) *CallSiteInfo {
	return &CallSiteInfo{
		MethodName:    methodName,
		ArgumentCount: argCount,
	}
}

func (*CallSiteInfo) Class() *Class {
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
	return fmt.Sprintf(
		"CallSiteInfo{method_name: %s, argument_count: %d}",
		c.MethodName.Inspect(),
		c.ArgumentCount,
	)
}
