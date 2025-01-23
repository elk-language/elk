package value

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var CallFrameClass *Class // ::Std::CallFrame

type CallFrame struct {
	FuncName        string
	FileName        string
	LineNumber      int
	TailCallCounter int
}

func (c *CallFrame) Copy() Reference {
	return c
}

func (*CallFrame) Class() *Class {
	return CallFrameClass
}

func (*CallFrame) DirectClass() *Class {
	return CallFrameClass
}

func (*CallFrame) SingletonClass() *Class {
	return nil
}

func (c *CallFrame) String() string {
	var buffer strings.Builder

	if c.TailCallCounter > 0 {
		fmt.Fprintf(&buffer, " ... %d optimised tail call(s)\n", c.TailCallCounter)
	}
	// "%s:%d, in `%s`"
	fmt.Fprintf(&buffer, "%s:%d, in ", c.FileName, c.LineNumber)
	color.New(color.FgHiYellow).Fprintf(&buffer, "`%s`", c.FuncName)

	// /tmp/test.elk:18, in `foo`
	return buffer.String()
}

func (c *CallFrame) Inspect() string {
	return fmt.Sprintf(
		"Std::CallFrame{&: %p, func_name: %s, file_name: %s, line_number: %d, tail_calls: %d}",
		c,
		String(c.FuncName).Inspect(),
		String(c.FileName).Inspect(),
		c.LineNumber,
		c.TailCallCounter,
	)
}

func (c *CallFrame) Error() string {
	return c.Inspect()
}

func (*CallFrame) InstanceVariables() SymbolMap {
	return nil
}

func initCallFrame() {
	CallFrameClass = NewClassWithOptions()
	StdModule.AddConstantString("CallFrame", Ref(CallFrameClass))
}
