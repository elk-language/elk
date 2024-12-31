//go:build debug

package vm

import "github.com/elk-language/elk/value"

// Add n to the instruction pointer
func (vm *VM) ipIncrementBy(n uintptr) {
	if vm.ipOffset()+int(n) >= len(vm.bytecode.Instructions) {
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

// Add n to the stack pointer
func (vm *VM) spIncrementBy(n int) {
	if vm.spOffset()+n >= VALUE_STACK_SIZE {
		panic("value stack overflow")
	}
	vm.sp = vm.sp + uintptr(n)*value.ValueSize
}

// Add n to the call frame pointer
func (vm *VM) cfpIncrementBy(n int) {
	if vm.cfpOffset()+n >= VALUE_STACK_SIZE {
		panic("call stack overflow")
	}
	vm.cfpSet(vm.cfpAdd(n))
}
