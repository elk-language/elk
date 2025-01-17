package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// Wraps a bytecode function with state that is necessary
// for pausing and resuming execution
type Generator struct {
	Bytecode *BytecodeFunction
	Upvalues []*Upvalue
	stack    []value.Value
	ip       uintptr
}

// Create a new generator
func NewGenerator(bytecode *BytecodeFunction) *Generator {
	return &Generator{
		Bytecode: bytecode,
		Upvalues: make([]*Upvalue, bytecode.UpvalueCount),
	}
}

func (*Generator) Class() *value.Class {
	return value.FunctionClass
}

func (*Generator) DirectClass() *value.Class {
	return value.FunctionClass
}

func (*Generator) SingletonClass() *value.Class {
	return nil
}

func (c *Generator) Copy() value.Reference {
	return c
}

func (c *Generator) Inspect() string {
	return fmt.Sprintf("Std::Generator{location: %s}", c.Bytecode.Location.String())
}

func (c *Generator) Error() string {
	return c.Inspect()
}

func (*Generator) InstanceVariables() value.SymbolMap {
	return nil
}
