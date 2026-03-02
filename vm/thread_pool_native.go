//go:build native

package vm

func executeBytecodePromise(thread *Thread, queue chan *Promise, task *Promise) {
	panic("cannot execute bytecode promises in native mode")
}
