//go:build debug

package vm

import "github.com/elk-language/elk/value"

// Add n to the instruction pointer
func (vm *VM) ipIncrementBy(n uintptr) {
	if vm.ipOffset()+int(n) > len(vm.bytecode.Instructions) {
		panic("ip overflow")
	}
	vm.ip = vm.ip + n
}

// Subtract n from the instruction pointer
func (vm *VM) ipDecrementBy(n uintptr) {
	if vm.ipOffset()-int(n) < 0 {
		panic("ip underflow")
	}
	vm.ip = vm.ip - n
}

// Subtract n from the stack pointer
func (vm *VM) spDecrementBy(n uintptr) {
	if vm.spOffset()-n < 0 {
		panic("value stack underflow")
	}
	vm.sp = vm.sp - n*value.ValueSize
}

// Add n to the stack pointer
func (vm *VM) spIncrementBy(n uintptr) {
	if vm.spOffset()+n > len(vm.stack) {
		panic("value stack overflow")
	}
	vm.sp = vm.sp + n*value.ValueSize
}

// Subtract n from the call frame pointer
func (vm *VM) cfpDecrementBy(n uintptr) {
	if vm.cfpOffset()-n < 0 {
		panic("call stack underflow")
	}
	vm.cfp = vm.cfpSubtractRaw(n)
}
