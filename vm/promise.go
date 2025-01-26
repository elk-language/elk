package vm

import (
	"fmt"
	"sync"

	"github.com/elk-language/elk/value"
)

type Promise struct {
	Generator    *Generator
	Continuation *Promise
	result       value.Value
	err          value.Value
	ch           chan struct{}
	m            sync.Mutex
}

// Create a new promise
func newPromise(generator *Generator, continuation *Promise) *Promise {
	return &Promise{
		Generator:    generator,
		Continuation: continuation,
		ch:           make(chan struct{}),
	}
}

func (*Promise) Class() *value.Class {
	return value.PromiseClass
}

func (*Promise) DirectClass() *value.Class {
	return value.PromiseClass
}

func (*Promise) SingletonClass() *value.Class {
	return nil
}

func (c *Promise) Copy() value.Reference {
	return c
}

func (c *Promise) Inspect() string {
	return fmt.Sprintf("Std::Promise{&: %p, resolved: %t}", c, c.IsResolved())
}

func (c *Promise) Error() string {
	return c.Inspect()
}

func (*Promise) InstanceVariables() value.SymbolMap {
	return nil
}

func (c *Promise) IsResolved() bool {
	return c.Generator == nil
}

func initPromise() {
	// Instance methods
	// c := &value.PromiseClass.MethodContainer
	// Def(
	// 	c,
	// 	"next",
	// 	func(vm *VM, args []value.Value) (value.Value, value.Value) {
	// 		self := (*Generator)(args[0].Pointer())
	// 		return vm.CallGeneratorNext(self)
	// 	},
	// )
	// Def(
	// 	c,
	// 	"iter",
	// 	func(_ *VM, args []value.Value) (value.Value, value.Value) {
	// 		return args[0], value.Undefined
	// 	},
	// )
	// Def(
	// 	c,
	// 	"reset",
	// 	func(_ *VM, args []value.Value) (value.Value, value.Value) {
	// 		self := (*Generator)(args[0].Pointer())
	// 		catch := self.Bytecode.CatchEntries[0]
	// 		self.ip = self.Bytecode.ipAddRaw(uintptr(catch.JumpAddress))
	// 		return args[0], value.Undefined
	// 	},
	// )
}
