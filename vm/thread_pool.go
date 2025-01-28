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
			return
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

func (t *ThreadPool) Close() {
	close(t.TaskQueue)
}

func initThreadPool() {
	value.ThreadPoolClass.AddConstantString("DEFAULT", value.Ref(DefaultThreadPool))

	// Instance methods
	c := &value.ThreadPoolClass.MethodContainer
	Def(
		c,
		"thread_count",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := (*ThreadPool)(args[0].Pointer())
			return value.SmallInt(self.ThreadCount()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"task_queue_size",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := (*ThreadPool)(args[0].Pointer())
			return value.SmallInt(self.TaskQueueSize()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"close",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*ThreadPool)(args[0].Pointer())
			self.Close()
			return value.Nil, value.Undefined
		},
	)
}
