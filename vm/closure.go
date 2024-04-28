package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// Wraps a bytecode function with associated local variables
// from the outer context
type Closure struct {
	Bytecode *BytecodeFunction
}

// Create a new closure
func NewClosure(bytecode *BytecodeFunction) *Closure {
	return &Closure{
		Bytecode: bytecode,
	}
}

func (*Closure) Class() *value.Class {
	return value.FunctionClass
}

func (*Closure) DirectClass() *value.Class {
	return value.FunctionClass
}

func (*Closure) SingletonClass() *value.Class {
	return nil
}

func (c *Closure) Copy() value.Value {
	return c
}

func (c *Closure) Inspect() string {
	return fmt.Sprintf("Function{location: %s}", c.Bytecode.Location.String())
}

func (*Closure) InstanceVariables() value.SymbolMap {
	return nil
}
