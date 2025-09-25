package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

const ClosureTerminatorFlag byte = 0xff

// Wraps a bytecode function with associated local variables
// from the outer context
type Closure struct {
	VMID     int64 // ID of the VM that created this closure, should be -1 if it was not created through a VM
	Bytecode *BytecodeFunction
	Self     value.Value
	Upvalues []*Upvalue
}

// Create a new native closure (created without a VM)
func NewNativeClosure(bytecode *BytecodeFunction, self value.Value) *Closure {
	return NewClosure(-1, bytecode, self)
}

// Create a new closure
func NewClosure(vmID int64, bytecode *BytecodeFunction, self value.Value) *Closure {
	return &Closure{
		VMID:     vmID,
		Bytecode: bytecode,
		Self:     self,
		Upvalues: make([]*Upvalue, bytecode.UpvalueCount),
	}
}

func (c *Closure) HasOpenUpvalues() bool {
	for _, upvalue := range c.Upvalues {
		if upvalue.IsOpen() {
			return true
		}
	}

	return false
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
