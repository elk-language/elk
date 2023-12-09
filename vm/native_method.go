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
	name                   value.Symbol
	parameters             []value.Symbol
	optionalParameterCount int
	postRestParameterCount int
	namedRestParameter     bool
	frozen                 bool
}

func NewNativeMethodComparer() cmp.Option {
	return cmp.Comparer(func(x, y *NativeMethod) bool {
		return x.name == y.name &&
			x.optionalParameterCount == y.optionalParameterCount &&
			x.postRestParameterCount == y.postRestParameterCount &&
			x.namedRestParameter == y.namedRestParameter &&
			x.frozen == y.frozen &&
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

func (n *NativeMethod) IsFrozen() bool {
	return n.frozen
}

func (n *NativeMethod) SetFrozen() {
	n.frozen = true
}

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
	frozen bool,
	function NativeFunction,
) *NativeMethod {
	return &NativeMethod{
		name:                   name,
		parameters:             params,
		optionalParameterCount: optParams,
		postRestParameterCount: postParams,
		namedRestParameter:     namedRestParam,
		frozen:                 frozen,
		Function:               function,
	}
}

// Native method constructor option
type NativeMethodOption func(*NativeMethod)

// Define parameters used by the method
func NativeMethodWithParameters(params []value.Symbol) NativeMethodOption {
	return func(n *NativeMethod) {
		n.parameters = params
	}
}

// Define parameters used by the method
func NativeMethodWithStringParameters(params ...string) NativeMethodOption {
	return func(n *NativeMethod) {
		symbolParams := make([]value.Symbol, len(params))
		for i, param := range params {
			symbolParams[i] = value.ToSymbol(param)
		}
		n.parameters = symbolParams
	}
}

// Define how many parameters are optional (have default values).
// Optional arguments will be populated with `undefined` when no value was given in the call.
func NativeMethodWithOptionalParameters(optParams int) NativeMethodOption {
	return func(n *NativeMethod) {
		n.optionalParameterCount = optParams
	}
}

func NativeMethodWithFunction(fn NativeFunction) NativeMethodOption {
	return func(n *NativeMethod) {
		n.Function = fn
	}
}

func NativeMethodWithName(name value.Symbol) NativeMethodOption {
	return func(n *NativeMethod) {
		n.name = name
	}
}

func NativeMethodWithStringName(name string) NativeMethodOption {
	return func(n *NativeMethod) {
		n.name = value.ToSymbol(name)
	}
}

// Set the last parameter as a positional rest parameter eg. `*rest`
func NativeMethodWithPositionalRestParameter() NativeMethodOption {
	return func(n *NativeMethod) {
		n.postRestParameterCount = 0
	}
}

// Define the number of parameters that appear after
// the positional rest parameter eg. 2 for `a, *b, c, d`
func NativeMethodWithPostParameters(postParams int) NativeMethodOption {
	return func(n *NativeMethod) {
		n.postRestParameterCount = postParams
	}
}

// Set the last parameter as a named rest parameter eg. `**rest`
func NativeMethodWithNamedRestParameter() NativeMethodOption {
	return func(n *NativeMethod) {
		n.namedRestParameter = true
	}
}

// Make the method non-overridable
func NativeMethodWithFrozen() NativeMethodOption {
	return func(n *NativeMethod) {
		n.frozen = true
	}
}

// Create a new native method with options
func NewNativeMethodWithOptions(opts ...NativeMethodOption) *NativeMethod {
	method := &NativeMethod{
		postRestParameterCount: -1,
	}

	for _, opt := range opts {
		opt(method)
	}

	return method
}

// Utility method that creates a new Function and
// attaches it as a method to the given method map.
func DefineMethod(
	methodMap value.MethodMap,
	name string,
	params []string,
	optParams int,
	postParams int,
	namedRestParam bool,
	frozen bool,
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
		frozen,
		function,
	)
	methodMap[symbolName] = nativeMethod
}

// Utility method that creates a new Function and
// attaches it as a method to the given method map.
func DefineMethodWithOptions(
	methodMap value.MethodMap,
	name string,
	function NativeFunction,
	opts ...NativeMethodOption,
) {
	symbolName := value.ToSymbol(name)
	nativeMethod := &NativeMethod{
		name:                   symbolName,
		Function:               function,
		postRestParameterCount: -1,
	}

	for _, opt := range opts {
		opt(nativeMethod)
	}

	methodMap[symbolName] = nativeMethod
}
