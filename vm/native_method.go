// Package bytecode implements types
// and constants that make up Elk
// bytecode.
package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// An implementation of a native Elk method.
type NativeFunction func(vm *VM, args []value.Value) (returnVal value.Value, err *value.Error)

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

func (*NativeMethod) InstanceVariables() value.SimpleSymbolMap {
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

func init() {
	DefineMethod(
		value.ObjectClass,
		"print",
		[]string{"val"},
		0,
		false,
		func(_ *VM, args []value.Value) (value.Value, *value.Error) {
			fmt.Println(args[1])

			return value.Nil, nil
		},
	)

}
