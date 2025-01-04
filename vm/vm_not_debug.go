//go:build !debug

package vm

import "github.com/elk-language/elk/value"

// Add n to the instruction pointer
func (vm *VM) ipIncrementBy(n uintptr) {
	vm.ip = vm.ip + n
}

// Subtract n from the instruction pointer
func (vm *VM) ipDecrementBy(n uintptr) {
	vm.ip = vm.ip - n
}

// Add n to the stack pointer
func (vm *VM) spIncrementBy(n uintptr) {
	vm.sp = vm.sp + n*value.ValueSize
	if vm.spGet().ValueFlag() == value.SENTINEL_FLAG {
		panic("value stack overflow")
	}
}

// Subtract n from the stack pointer
func (vm *VM) spDecrementBy(n uintptr) {
	vm.sp = vm.sp - n*value.ValueSize
}

// Add n to the call frame pointer
func (vm *VM) cfpIncrementBy(n int) {
	ptr := vm.cfpAdd(n)
	if ptr.sentinel {
		panic("call stack overflow")
	}
	vm.cfpSet(ptr)
}
