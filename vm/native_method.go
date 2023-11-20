// Package bytecode implements types
// and constants that make up Elk
// bytecode.
package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// An implementation of a native Elk method.
type NativeFunction func(vm *VM, args []value.Value) (returnVal, err value.Value)

// A single unit of Elk bytecode.
type NativeMethod struct {
	Name                   value.Symbol
	Parameters             []value.Symbol
	OptionalParameterCount int
	NamedRestArgument      bool
	Function               NativeFunction
}

func (*NativeMethod) Method() {}

// Get the number of parameters
func (b *NativeMethod) ParameterCount() int {
	return len(b.Parameters)
}

func (*NativeMethod) Class() *value.Class {
	return nil
}

func (*NativeMethod) DirectClass() *value.Class {
	return nil
}

func (*NativeMethod) SingletonClass() *value.Class {
	return nil
}

func (*NativeMethod) IsFrozen() bool {
	return true
}

func (*NativeMethod) SetFrozen() {}

func (b *NativeMethod) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :native}", b.Name.Name())
}

func (*NativeMethod) InstanceVariables() value.SymbolMap {
	return nil
}

// Create a new native method.
func NewNativeMethod(
	name value.Symbol,
	params []value.Symbol,
	optParams int,
	namedRestArg bool,
	function NativeFunction,
) *NativeMethod {
	return &NativeMethod{
		Name:                   name,
		Parameters:             params,
		OptionalParameterCount: optParams,
		NamedRestArgument:      namedRestArg,
		Function:               function,
	}
}

// Utility method that creates a new Function and
// attaches it as a method to the given class.
func DefineMethod(
	class *value.Class,
	name string,
	params []string,
	optParams int,
	namedRestArg bool,
	function NativeFunction,
) {
	symbolParams := make([]value.Symbol, len(params))
	for i, param := range params {
		symbolParams[i] = value.ToSymbol(param)
	}

	symbolName := value.ToSymbol(name)
	nativeFunc := NewNativeMethod(
		symbolName,
		symbolParams,
		optParams,
		namedRestArg,
		function,
	)
	class.Methods[symbolName] = nativeFunc
}

// Define a method that takes no arguments.
func DefineMethodNoParams(
	class *value.Class,
	name string,
	function NativeFunction,
) {
	DefineMethod(class, name, nil, 0, false, function)
}

// Define a method that has required parameters.
func DefineMethodReqParams(
	class *value.Class,
	name string,
	params []string,
	function NativeFunction,
) {
	DefineMethod(class, name, params, 0, false, function)
}

// Define a method that has optional parameters.
func DefineMethodOptParams(
	class *value.Class,
	name string,
	params []string,
	function NativeFunction,
) {
	DefineMethod(class, name, params, 0, false, function)
}

func init() {
	DefineMethodReqParams(
		value.ObjectClass,
		"print",
		[]string{"val"},
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			fmt.Println(args[1])

			return value.Nil, nil
		},
	)

}
