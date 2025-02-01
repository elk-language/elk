package vm

import (
	"fmt"
	"sync"

	"github.com/elk-language/elk/value"
)

type Promise struct {
	*Generator
	ThreadPool    *ThreadPool
	continuations []*Promise
	result        value.Value
	err           value.Value
	wg            sync.WaitGroup // the wait group hits 0 when the promise is resolved, used for waiting for a promise
	m             sync.Mutex
}

// Create a new promise executed by the VM
func newPromise(threadPool *ThreadPool, generator *Generator) *Promise {
	p := &Promise{
		ThreadPool: threadPool,
		Generator:  generator,
	}
	p.wg.Add(1)

	threadPool.AddTask(p)
	return p
}

// Returns a new native promise handled by Go code instead of the VM
func newNativePromise(threadPool *ThreadPool) *Promise {
	p := &Promise{
		ThreadPool: threadPool,
	}
	p.wg.Add(1)
	return p
}

func NewResolvedPromise(result value.Value) *Promise {
	return &Promise{
		result: result,
	}
}

func NewRejectedPromise(err value.Value) *Promise {
	return &Promise{
		err: err,
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
	if p.IsResolved() {
		if !p.err.IsUndefined() {
			return fmt.Sprintf("Std::Promise{&: %p, resolved: true, err: %s}", p, p.err.Inspect())
		}

		return fmt.Sprintf("Std::Promise{&: %p, resolved: true, result: %s}", p, p.result.Inspect())
	}

	return fmt.Sprintf("Std::Promise{&: %p, resolved: false}", p)
}

func (p *Promise) Error() string {
	return p.Inspect()
}

func (*Promise) InstanceVariables() value.SymbolMap {
	return nil
}

func (p *Promise) IsResolved() bool {
	return p.ThreadPool == nil
}

func (p *Promise) AwaitSync() (value.Value, value.Value) {
	p.wg.Wait()
	return p.result, p.err
}

func (p *Promise) RegisterContinuation(continuation *Promise) {
	p.m.Lock()
	p.continuations = append(p.continuations, continuation)
	p.m.Unlock()
}

func (p *Promise) RegisterContinuationUnsafe(continuation *Promise) {
	p.continuations = append(p.continuations, continuation)
}

func (p *Promise) enqueueContinuations(queue chan *Promise) {
	for _, cont := range p.continuations {
		queue <- cont
	}
	p.continuations = nil
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
