package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

type ThreadPool struct {
	Threads   []*VM
	TaskQueue chan *Promise
}

func NewThreadPool(threadCount, queueSize int, opts ...Option) *ThreadPool {
	tp := &ThreadPool{
		TaskQueue: make(chan *Promise, queueSize),
	}

	threads := make([]*VM, threadCount)
	for i := range threads {
		thread := New(opts...)
		threads[i] = thread
		go threadWorker(thread, tp.TaskQueue)
	}
	tp.Threads = threads

	return tp
}

func threadWorker(thread *VM, queue <-chan *Promise) {
	for {
		task, ok := <-queue
		if !ok {
			break
		}

		thread.callPromise(task)
		switch thread.state {
		case awaitState:
			// TODO
		case errorState:
			err := thread.popGet()
			task.reject(err)
		default:
			result := thread.popGet()
			task.resolve(result)
		}
	}
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
	return fmt.Sprintf("Std::ThreadPool{thread_count: %d, task_queue_size: %d}", t.ThreadCount(), t.TaskQueueSize())
}

func (t *ThreadPool) Error() string {
	return t.Inspect()
}

func (*ThreadPool) InstanceVariables() value.SymbolMap {
	return nil
}

func (t *ThreadPool) TaskQueueSize() int {
	return cap(t.TaskQueue)
}

func (t *ThreadPool) ThreadCount() int {
	return len(t.Threads)
}

func (t *ThreadPool) AddTask(promise *Promise) {
	t.TaskQueue <- promise
}

func initThreadPool() {
	value.ThreadPoolClass.AddConstantString("DEFAULT", value.Ref(DefaultThreadPool))

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
