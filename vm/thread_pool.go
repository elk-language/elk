package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

type ThreadPool struct {
	Threads []*VM
	Queue   chan *Promise
}

func NewThreadPool(threadCount, queueSize int, opts ...Option) *ThreadPool {
	tp := &ThreadPool{
		Queue: make(chan *Promise, queueSize),
	}

	threads := make([]*VM, threadCount)
	for i := range threads {
		thread := New(opts...)
		threads[i] = thread
		go threadWorker(thread, tp.Queue)
	}
	tp.Threads = threads

	return tp
}

func threadWorker(thread *VM, queue <-chan *Promise) {
	// for {
	// 	task, ok := <-queue
	// 	if !ok {
	// 		break
	// 	}

	// 	// TODO
	// }
}

func (*ThreadPool) Class() *value.Class {
	return value.ThreadPoolClass
}

func (*ThreadPool) DirectClass() *value.Class {
	return value.ThreadPoolClass
}

func (*ThreadPool) SingletonClass() *value.Class {
	return nil
}

func (t *ThreadPool) Copy() value.Reference {
	return t
}

func (t *ThreadPool) Inspect() string {
	return fmt.Sprintf("Std::ThreadPool{thread_count: %d, queue_size: %d}", t.ThreadSize(), t.QueueSize())
}

func (t *ThreadPool) Error() string {
	return t.Inspect()
}

func (*ThreadPool) InstanceVariables() value.SymbolMap {
	return nil
}

func (t *ThreadPool) QueueSize() int {
	return cap(t.Queue)
}

func (t *ThreadPool) ThreadSize() int {
	return len(t.Queue)
}

func initThreadPool() {
	// Instance methods
	// c := &value.ThreadPoolClass.MethodContainer
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
