//go:build !native

package vm

func executeBytecodePromise(thread *Thread, queue chan *Promise, task *Promise) {
	thread.callBytecodePromise(task)

	switch thread.state {
	case awaitState:
		awaitedPromise := (*Promise)(thread.peek().Pointer())
		awaitedPromise.RegisterContinuationUnsafe(task)

		// promise has been locked in the VM
		awaitedPromise.m.Unlock()
	case errorState:
		err := thread.popGet()
		stackTrace := thread.GetStackTrace()
		task.Reject(err, stackTrace)
	default:
		result := thread.popGet()
		task.Resolve(result)
	}

}
