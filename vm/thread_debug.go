//go:build debug

package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

func (vm *Thread) populateMissingParametersOnStack(paramCount, argumentCount int) {
	// populate missing optional arguments with undefined
	missingParams := paramCount - argumentCount
	if missingParams < 0 {
		panic(fmt.Sprintf(
			"invalid missingParams: %d, ip: %d, bytecode:\n%s\n",
			missingParams,
			vm.ipOffset(),
			vm.bytecode.MustDisassembleString(),
		))
	}
	for range missingParams {
		vm.push(value.Undefined)
	}
}

// Add n to the instruction pointer
func (vm *Thread) ipIncrementBy(n uintptr) {
	if vm.ipOffset()+int(n) > len(vm.bytecode.Instructions) {
		panic(fmt.Sprintf(
			"ip overflow, increment: %d, ip: %d, bytecode:\n%s\n",
			n,
			vm.ipOffset(),
			vm.bytecode.MustDisassembleString(),
		))
	}
	vm.ip = vm.ip + n
}

// Subtract n from the instruction pointer
func (vm *Thread) ipDecrementBy(n uintptr) {
	if vm.ipOffset()-int(n) < 0 {
		panic(fmt.Sprintf(
			"ip underflow, decrement: %d, ip: %d, bytecode:\n%s\n",
			n,
			vm.ipOffset(),
			vm.bytecode.MustDisassembleString(),
		))
	}
	vm.ip = vm.ip - n
}

// Subtract n from the stack pointer
func (vm *Thread) spDecrementBy(n uintptr) {
	if vm.spOffset()-int(n) < 0 {
		panic(fmt.Sprintf(
			"value stack underflow, decrement: %d, ip: %d, bytecode:\n%s\n",
			n,
			vm.ipOffset(),
			vm.bytecode.MustDisassembleString(),
		))
	}
	vm.sp = vm.sp - n*value.ValueSize
}

// Add n to the stack pointer
func (vm *Thread) spIncrementBy(n uintptr) {
	if vm.spOffset()+int(n) > len(vm.stack) {
		panic(fmt.Sprintf(
			"value stack overflow, increment: %d, ip: %d, bytecode:\n%s\n",
			n,
			vm.ipOffset(),
			vm.bytecode.MustDisassembleString(),
		))
	}
	vm.sp = vm.sp + n*value.ValueSize
}

// Subtract n from the call frame pointer
func (vm *Thread) cfpDecrementBy(n uintptr) {
	if vm.cfpOffset()-int(n) < 0 {
		panic(fmt.Sprintf(
			"class stack underflow, decrement: %d, ip: %d, bytecode:\n%s\n",
			n,
			vm.ipOffset(),
			vm.bytecode.MustDisassembleString(),
		))
	}
	vm.cfp = vm.cfpSubtractRaw(n)
}
