// Package bytecode implements types
// and constants that make up Elk
// bytecode.
package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

// An implementation of a native Elk method.
type NativeFunction func(vm *VM, args []value.Value) (returnVal, err value.Value)

// A single unit of Elk bytecode.
type NativeMethod struct {
	Function               NativeFunction
	name                   value.Symbol
	parameters             []value.Symbol
	optionalParameterCount int
	postRestParameterCount int
	namedRestParameter     bool
}

func NewNativeMethodComparer() cmp.Option {
	return cmp.Comparer(func(x, y *NativeMethod) bool {
		if x.Function != nil && y.Function == nil || x.Function == nil && y.Function != nil {
			return false
		}
		return x.name == y.name &&
			x.optionalParameterCount == y.optionalParameterCount &&
			x.postRestParameterCount == y.postRestParameterCount &&
			x.namedRestParameter == y.namedRestParameter &&
			cmp.Equal(x.parameters, y.parameters)
	})
}

func (n *NativeMethod) Name() value.Symbol {
	return n.name
}

func (n *NativeMethod) ParameterCount() int {
	return len(n.parameters)
}

func (n *NativeMethod) Parameters() []value.Symbol {
	return n.parameters
}

func (n *NativeMethod) OptionalParameterCount() int {
	return n.optionalParameterCount
}

func (n *NativeMethod) PostRestParameterCount() int {
	return n.postRestParameterCount
}

func (n *NativeMethod) NamedRestParameter() bool {
	return n.namedRestParameter
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

func (n *NativeMethod) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :native}", n.name.Inspect())
}

func (*NativeMethod) InstanceVariables() value.SymbolMap {
	return nil
}

// Create a new native method.
func NewNativeMethod(
	name value.Symbol,
	params []value.Symbol,
	optParams int,
	postParams int,
	namedRestParam bool,
	function NativeFunction,
) *NativeMethod {
	return &NativeMethod{
		name:                   name,
		parameters:             params,
		optionalParameterCount: optParams,
		postRestParameterCount: postParams,
		namedRestParameter:     namedRestParam,
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
	postParams int,
	namedRestParam bool,
	function NativeFunction,
) {
	symbolParams := make([]value.Symbol, len(params))
	for i, param := range params {
		symbolParams[i] = value.ToSymbol(param)
	}

	symbolName := value.ToSymbol(name)
	nativeMethod := NewNativeMethod(
		symbolName,
		symbolParams,
		optParams,
		postParams,
		namedRestParam,
		function,
	)
	class.Methods[symbolName] = nativeMethod
}

// Define a method that takes no arguments.
func DefineMethodNoParams(
	class *value.Class,
	name string,
	function NativeFunction,
) {
	DefineMethod(class, name, nil, 0, -1, false, function)
}

// Define a method that has required parameters.
func DefineMethodReqParams(
	class *value.Class,
	name string,
	params []string,
	function NativeFunction,
) {
	DefineMethod(class, name, params, 0, -1, false, function)
}

// Define a method that has optional parameters.
func DefineMethodOptParams(
	class *value.Class,
	name string,
	params []string,
	function NativeFunction,
) {
	DefineMethod(class, name, params, 0, -1, false, function)
}

// Define a method that has a splat parameter.
func DefineMethodSplatParams(
	class *value.Class,
	name string,
	params []string,
	function NativeFunction,
) {
	DefineMethod(class, name, params, 0, 0, false, function)
}

// Define a method that has a splat parameter.
func DefineMethodPostParams(
	class *value.Class,
	name string,
	params []string,
	optParams int,
	postParams int,
	function NativeFunction,
) {
	DefineMethod(class, name, params, optParams, postParams, false, function)
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
