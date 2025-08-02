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
	stackTrace    *value.StackTrace
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

// Create a new promise for a piece of bytecode executed by the VM
func NewPromiseForBytecode(threadPool *ThreadPool, bytecode *BytecodeFunction, args ...value.Value) *Promise {
	generator := NewGeneratorForBytecode(bytecode, args...)

	p := &Promise{
		ThreadPool: threadPool,
		Generator:  generator,
	}
	p.wg.Add(1)

	threadPool.AddTask(p)
	return p
}

// Returns a new native promise handled by Go code instead of the VM
func NewNativePromise(threadPool *ThreadPool) *Promise {
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

func (*Promise) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (p *Promise) IsResolved() bool {
	return p.ThreadPool == nil
}

// Wait for the result of the promise
func (p *Promise) AwaitSync() (value.Value, *value.StackTrace, value.Value) {
	p.wg.Wait()
	return p.result, p.stackTrace, p.err
}

func (p *Promise) RegisterContinuation(continuation *Promise) {
	p.m.Lock()
	p.continuations = append(p.continuations, continuation)
	p.m.Unlock()
}

func (p *Promise) RegisterContinuationUnsafe(continuation *Promise) {
	p.continuations = append(p.continuations, continuation)
}

func (p *Promise) ResolveReject(result, err value.Value) {
	p.m.Lock()

	queue := p.ThreadPool.TaskQueue
	p.Generator = nil
	p.ThreadPool = nil
	p.result = result
	p.err = err
	p.wg.Done()
	p.enqueueContinuations(queue)

	p.m.Unlock()
}

func (p *Promise) Resolve(result value.Value) {
	p.m.Lock()

	queue := p.ThreadPool.TaskQueue
	p.Generator = nil
	p.ThreadPool = nil
	p.result = result
	p.wg.Done()
	p.enqueueContinuations(queue)

	p.m.Unlock()
}

func (p *Promise) Reject(err value.Value) {
	p.m.Lock()

	queue := p.ThreadPool.TaskQueue
	p.Generator = nil
	p.ThreadPool = nil
	p.err = err
	p.wg.Done()
	p.enqueueContinuations(queue)

	p.m.Unlock()
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
	Def(
		c,
		"wait",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			collectionValue := args[1]

			p := NewNativePromise(vm.threadPool)
			go func(vm *VM, p *Promise, collection value.Value) {
				for val, err := range Iterate(vm, collection) {
					promise := (*Promise)(val.Pointer())
					if !err.IsUndefined() {
						p.Reject(err)
						return
					}
					_, _, err = promise.AwaitSync()
					if !err.IsUndefined() {
						p.Reject(err)
						return
					}
				}

				p.Resolve(value.Nil)
			}(vm, p, collectionValue)

			return value.Ref(p), value.Undefined
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.PromiseClass.MethodContainer
	Def(
		c,
		"is_resolved",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := (*Promise)(args[0].Pointer())
			return value.ToElkBool(self.IsResolved()), value.Undefined
		},
	)
}
