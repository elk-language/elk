package vm

import (
	"fmt"
	"sync"

	"github.com/elk-language/elk/value"
)

type Promise struct {
	*Generator
	ThreadPool   *ThreadPool
	Continuation *Promise
	result       value.Value
	err          value.Value
	ch           chan struct{} // the channel gets closed when the promise is resolved, used for waiting for a promise
	m            sync.Mutex
}

// Create a new promise
func newPromise(threadPool *ThreadPool, generator *Generator, continuation *Promise) *Promise {
	p := &Promise{
		ThreadPool:   threadPool,
		Generator:    generator,
		Continuation: continuation,
		ch:           make(chan struct{}),
	}

	threadPool.AddTask(p)
	return p
}

func NewResolvedPromise(result value.Value) *Promise {
	ch := make(chan struct{})
	close(ch)

	return &Promise{
		result: result,
		ch:     ch,
	}
}

func NewRejectedPromise(err value.Value) *Promise {
	ch := make(chan struct{})
	close(ch)

	return &Promise{
		err: err,
		ch:  ch,
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

func (p *Promise) Copy() value.Reference {
	return p
}

func (p *Promise) Inspect() string {
	return fmt.Sprintf("Std::Promise{&: %p, resolved: %t}", p, p.IsResolved())
}

func (p *Promise) Error() string {
	return p.Inspect()
}

func (*Promise) InstanceVariables() value.SymbolMap {
	return nil
}

func (p *Promise) IsResolved() bool {
	return p.Generator == nil
}

func (p *Promise) AwaitSync() (value.Value, value.Value) {
	<-p.ch
	return p.result, p.err
}

func initPromise() {
	// Singleton methods
	c := &value.PromiseClass.SingletonClass().MethodContainer
	Def(
		c,
		"resolved",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			result := args[1]
			return value.Ref(NewResolvedPromise(result)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"rejected",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			err := args[1]
			return value.Ref(NewRejectedPromise(err)), value.Undefined
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.PromiseClass.MethodContainer
	Def(
		c,
		"await_sync",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := (*Promise)(args[0].Pointer())
			return self.AwaitSync()
		},
	)
	Def(
		c,
		"is_resolved",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := (*Promise)(args[0].Pointer())
			return value.ToElkBool(self.IsResolved()), value.Undefined
		},
	)
}
