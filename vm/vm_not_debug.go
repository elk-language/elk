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

// Subtract n from the stack pointer
func (vm *VM) spDecrementBy(n uintptr) {
	vm.sp = vm.sp - n*value.ValueSize
}

// Add n to the stack pointer
func (vm *VM) spIncrementBy(n uintptr) {
	vm.sp = vm.sp + n*value.ValueSize
}

// Subtract n from the call frame pointer
func (vm *VM) cfpDecrementBy(n uintptr) {
	vm.cfp = vm.cfpSubtractRaw(n)
}
