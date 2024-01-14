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
	parameters             []value.Symbol
	optionalParameterCount int
	postRestParameterCount int
	namedRestParameter     bool
	sealed                 bool
}

func NewNativeMethodComparer() cmp.Option {
	return cmp.Comparer(func(x, y *NativeMethod) bool {
		return x.name == y.name &&
			x.optionalParameterCount == y.optionalParameterCount &&
			x.postRestParameterCount == y.postRestParameterCount &&
			x.namedRestParameter == y.namedRestParameter &&
			x.sealed == y.sealed &&
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

func (n *NativeMethod) IsSealed() bool {
	return n.sealed
}

func (n *NativeMethod) SetSealed() {
	n.sealed = true
}

func (n *NativeMethod) Copy() value.Value {
	return n
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
	sealed bool,
	function NativeFunction,
) *NativeMethod {
	return &NativeMethod{
		name:                   name,
		parameters:             params,
		optionalParameterCount: optParams,
		postRestParameterCount: postParams,
		namedRestParameter:     namedRestParam,
		sealed:                 sealed,
		Function:               function,
	}
}

// Define a native method in the given container.
// Returns an error when the method couldn't be defined.
func DefineNativeMethod(
	container *value.MethodContainer,
	name value.Symbol,
	params []value.Symbol,
	optParams int,
	postParams int,
	namedRestParam bool,
	sealed bool,
	function NativeFunction,
) *value.Error {
	if !container.CanOverride(name) {
		return value.NewCantOverrideASealedMethod(string(name.ToString()))
	}

	nativeMethod := NewNativeMethod(
		name,
		params,
		optParams,
		postParams,
		namedRestParam,
		sealed,
		function,
	)
	container.Methods[name] = nativeMethod
	return nil
}

type DefOption func(*NativeMethod)

// Define parameters used by the method
func DefWithParameters(params ...string) DefOption {
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
func DefWithOptionalParameters(optParams int) DefOption {
	return func(n *NativeMethod) {
		n.optionalParameterCount = optParams
	}
}

// Set the last parameter as a positional rest parameter eg. `*rest`
func DefWithPositionalRestParameter() DefOption {
	return func(n *NativeMethod) {
		n.postRestParameterCount = 0
	}
}

// Define the number of parameters that appear after
// the positional rest parameter eg. 2 for `a, *b, c, d`
func DefWithPostParameters(postParams int) DefOption {
	return func(n *NativeMethod) {
		n.postRestParameterCount = postParams
	}
}

// Set the last parameter as a named rest parameter eg. `**rest`
func DefWithNamedRestParameter() DefOption {
	return func(n *NativeMethod) {
		n.namedRestParameter = true
	}
}

// Make the method non-overridable
func DefWithSealed() DefOption {
	return func(n *NativeMethod) {
		n.sealed = true
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
	if !container.CanOverride(symbolName) {
		panic(value.NewCantOverrideASealedMethod(name))
	}

	nativeMethod := &NativeMethod{
		name:                   symbolName,
		Function:               function,
		postRestParameterCount: -1,
	}

	for _, opt := range opts {
		opt(nativeMethod)
	}

	container.Methods[symbolName] = nativeMethod
}
