package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

// An implementation of a native Elk method.
type NativeFunction func(vm *VM, args []value.Value) (returnVal, err value.Value)

// A native Elk method
type NativeMethod struct {
	Function               NativeFunction
	Doc                    value.Value
	name                   value.Symbol
	parameterCount         int
	optionalParameterCount int
}

func NewNativeMethodComparer() cmp.Option {
	return cmp.Comparer(func(x, y *NativeMethod) bool {
		return x.name == y.name &&
			x.optionalParameterCount == y.optionalParameterCount &&
			x.parameterCount == y.parameterCount
	})
}

func (n *NativeMethod) Name() value.Symbol {
	return n.name
}

func (n *NativeMethod) ParameterCount() int {
	return n.parameterCount
}

func (n *NativeMethod) OptionalParameterCount() int {
	return n.optionalParameterCount
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

func (n *NativeMethod) Copy() value.Value {
	return n
}

func (n *NativeMethod) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :native}", n.name.Inspect())
}

func (n *NativeMethod) Error() string {
	return n.Inspect()
}

func (*NativeMethod) InstanceVariables() value.SymbolMap {
	return nil
}

// Create a new native method.
func NewNativeMethod(
	name value.Symbol,
	params int,
	optParams int,
	function NativeFunction,
) *NativeMethod {
	return &NativeMethod{
		name:                   name,
		parameterCount:         params,
		optionalParameterCount: optParams,
		Function:               function,
	}
}

// Define a native method in the given container.
// Returns an error when the method couldn't be defined.
func DefineNativeMethod(
	container *value.MethodContainer,
	name value.Symbol,
	params int,
	optParams int,
	function NativeFunction,
) (err value.Value) {
	nativeMethod := NewNativeMethod(
		name,
		params,
		optParams,
		function,
	)
	container.Methods[name] = nativeMethod
	return nil
}

type DefOption func(*NativeMethod)

// Define parameters used by the method
func DefWithParameters(params int) DefOption {
	return func(n *NativeMethod) {
		n.parameterCount = params
	}
}

// Define how many parameters are optional (have default values).
// Optional arguments will be populated with `undefined` when no value was given in the call.
func DefWithOptionalParameters(optParams int) DefOption {
	return func(n *NativeMethod) {
		n.optionalParameterCount = optParams
	}
}

// Utility method that creates a new native
// method and attaches it to the given container.
//
// Panics when the method cannot be defined.
func Def(
	container *value.MethodContainer,
	name string,
	function NativeFunction,
	opts ...DefOption,
) {
	symbolName := value.ToSymbol(name)

	nativeMethod := &NativeMethod{
		name:     symbolName,
		Function: function,
	}

	for _, opt := range opts {
		opt(nativeMethod)
	}

	container.Methods[symbolName] = nativeMethod
}
