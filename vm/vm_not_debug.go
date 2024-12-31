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
func (vm *VM) spIncrementBy(n int) {
	vm.sp = vm.sp + uintptr(n)*value.ValueSize
}

// Add n to the call frame pointer
func (vm *VM) cfpIncrementBy(n int) {
	vm.cfpSet(vm.cfpAdd(n))
}
