package vm

import (
	"fmt"
	"unsafe"

	"github.com/elk-language/elk/value"
)

// Wraps a bytecode function with state that is necessary
// for pausing and resuming execution
type Generator struct {
	Bytecode *BytecodeFunction
	ip       uintptr
	upvalues []*Upvalue
	stack    []value.Value
}

// Create a new generator
func newGenerator(
	bytecode *BytecodeFunction,
	upvalues []*Upvalue,
	stack []value.Value,
	ip uintptr,
) *Generator {
	return &Generator{
		Bytecode: bytecode,
		upvalues: upvalues,
		stack:    stack,
		ip:       ip,
	}
}

// Create a new generator that executes the given piece of bytecode.
func NewGeneratorForBytecode(bytecode *BytecodeFunction, args ...value.Value) *Generator {
	return newGenerator(
		bytecode,
		nil,
		args,
		uintptr(unsafe.Pointer(&bytecode.Instructions[0])),
	)
}

func (*Generator) Class() *value.Class {
	return value.GeneratorClass
}

func (*Generator) DirectClass() *value.Class {
	return value.GeneratorClass
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

func (*Generator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func initGenerator() {
	// Instance methods
	c := &value.GeneratorClass.MethodContainer
	Def(
		c,
		"next",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := (*Generator)(args[0].Pointer())
			return vm.CallGeneratorNext(self)
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"reset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*Generator)(args[0].Pointer())
			catch := self.Bytecode.CatchEntries[0]
			self.ip = self.Bytecode.ipAddRaw(uintptr(catch.JumpAddress))
			return args[0], value.Undefined
		},
	)
}
