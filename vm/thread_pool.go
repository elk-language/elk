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
	tp := &ThreadPool{}
	tp.initThreadPool(threadCount, queueSize, opts...)
	return tp
}

func (tp *ThreadPool) initThreadPool(threadCount, queueSize int, opts ...Option) {
	tp.TaskQueue = make(chan *Promise, queueSize)

	threads := make([]*VM, threadCount)
	for i := range threads {
		thread := New(opts...)
		threads[i] = thread
		go threadWorker(thread, tp.TaskQueue)
	}
	tp.Threads = threads
}

func threadWorker(thread *VM, queue chan *Promise) {
	for task := range queue {
		thread.callPromise(task)

		switch thread.state {
		case awaitState:
			awaitedPromise := (*Promise)(thread.peek().Pointer())
			awaitedPromise.RegisterContinuationUnsafe(task)

			// promise has been locked in the VM
			awaitedPromise.m.Unlock()
		case errorState:
			err := thread.popGet()
			task.m.Lock()

			task.Generator = nil
			task.ThreadPool = nil
			task.err = err
			task.wg.Done()
			task.enqueueContinuations(queue)

			task.m.Unlock()
		default:
			result := thread.popGet()
			task.m.Lock()

			task.Generator = nil
			task.ThreadPool = nil
			task.result = result
			task.wg.Done()
			task.enqueueContinuations(queue)

			task.m.Unlock()
		}

		thread.state = idleState
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
