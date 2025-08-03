package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

const ClosureTerminatorFlag byte = 0xff

// Wraps a bytecode function with associated local variables
// from the outer context
type Closure struct {
	Bytecode *BytecodeFunction
	Self     value.Value
	Upvalues []*Upvalue
}

// Create a new closure
func NewClosure(bytecode *BytecodeFunction, self value.Value) *Closure {
	return &Closure{
		Bytecode: bytecode,
		Self:     self,
		Upvalues: make([]*Upvalue, bytecode.UpvalueCount),
	}
}

func (*Closure) Class() *value.Class {
	return value.ClosureClass
}

func (*Closure) DirectClass() *value.Class {
	return value.ClosureClass
}

func (*Closure) SingletonClass() *value.Class {
	return nil
}

func (c *Closure) Copy() value.Reference {
	return c
}

func (c *Closure) Inspect() string {
	return fmt.Sprintf("Std::Closure{location: %s}", c.Bytecode.Location.String())
}

func (c *Closure) Error() string {
	return c.Inspect()
}

func (*Closure) InstanceVariables() *value.InstanceVariables {
	return nil
}
