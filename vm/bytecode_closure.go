package vm

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

const ClosureTerminatorFlag byte = 0xff

// Wraps a bytecode function with associated local variables
// from the outer context
type BytecodeClosure struct {
	VMID     int64 // ID of the VM that created this closure, should be -1 if it was not created through a VM
	Bytecode *BytecodeFunction
	Self     value.Value
	Upvalues []*Upvalue
}

var _ Closure = &BytecodeClosure{}

// Create a new closure
func NewBytecodeClosure(vmID int64, bytecode *BytecodeFunction, self value.Value) *BytecodeClosure {
	return &BytecodeClosure{
		VMID:     vmID,
		Bytecode: bytecode,
		Self:     self,
		Upvalues: make([]*Upvalue, bytecode.UpvalueCount),
	}
}

func (c *BytecodeClosure) HasOpenUpvalues() bool {
	for _, upvalue := range c.Upvalues {
		if upvalue.IsOpen() {
			return true
		}
	}

	return false
}

func (c *BytecodeClosure) ParameterCount() int {
	return c.Bytecode.parameterCount
}

func (c *BytecodeClosure) Location() *position.Location {
	return c.Bytecode.Location
}

func (*BytecodeClosure) Class() *value.Class {
	return value.ClosureClass
}

func (*BytecodeClosure) DirectClass() *value.Class {
	return value.ClosureClass
}

func (*BytecodeClosure) SingletonClass() *value.Class {
	return nil
}

func (c *BytecodeClosure) Copy() value.Reference {
	return c
}

func (c *BytecodeClosure) ToValue() value.Value {
	return value.Ref(c)
}

func (c *BytecodeClosure) Inspect() string {
	return fmt.Sprintf("Std::Closure{location: %s, type: :bytecode}", c.Bytecode.Location.String())
}

func (c *BytecodeClosure) Error() string {
	return c.Inspect()
}

func (*BytecodeClosure) InstanceVariables() *value.InstanceVariables {
	return nil
}
