package vm

import (
	"fmt"

	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/google/go-cmp/cmp"
)

// An implementation of a native Elk method.
type NativeFunction func(vm *Thread, args []value.Value) (returnVal, err value.Value)

// A native Elk method
type NativeMethod struct {
	Function               NativeFunction
	namespace              *value.Class
	name                   value.Symbol
	parameterCount         int
	optionalParameterCount int
}

var _ value.Method = &NativeMethod{}

func MethodToFunc(method value.Method) NativeFunction {
	return method.(*NativeMethod).Function
}

func NewNativeMethodComparer() cmp.Option {
	return cmp.Comparer(func(x, y *NativeMethod) bool {
		return x.name == y.name &&
			x.optionalParameterCount == y.optionalParameterCount &&
			x.parameterCount == y.parameterCount &&
			x.namespace == y.namespace
	})
}

func (n *NativeMethod) MethodBody() {}

func (n *NativeMethod) Name() value.Symbol {
	return n.name
}

func (n *NativeMethod) Namespace() *value.Class {
	return n.namespace
}

func (n *NativeMethod) SetNamespace(namespace *value.Class) {
	n.namespace = namespace
}

func (n *NativeMethod) ParameterCount() int {
	return n.parameterCount
}

func (n *NativeMethod) OptionalParameterCount() int {
	return n.optionalParameterCount
}

func (*NativeMethod) Class() *value.Class {
	return value.MethodClass
}

func (*NativeMethod) DirectClass() *value.Class {
	return value.MethodClass
}

func (*NativeMethod) SingletonClass() *value.Class {
	return nil
}

func (n *NativeMethod) Copy() value.Reference {
	return n
}

func (n *NativeMethod) ToValue() value.Value {
	return value.Ref(n)
}

func (n *NativeMethod) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :native, namespace: %s}", n.name.Inspect(), n.namespace.Inspect())
}

func (n *NativeMethod) Error() string {
	return n.Inspect()
}

func (*NativeMethod) InstanceVariables() *value.InstanceVariables {
	return nil
}

// Create a new native method.
func NewNativeMethod(
	namespace *value.Class,
	name value.Symbol,
	params int,
	optParams int,
	function NativeFunction,
) *NativeMethod {
	return &NativeMethod{
		namespace:              namespace,
		name:                   name,
		parameterCount:         params,
		optionalParameterCount: optParams,
		Function:               function,
	}
}

// Define a native method in the given container.
// Returns an error when the method couldn't be defined.
func DefineNativeMethod(
	namespace *value.Class,
	name value.Symbol,
	params int,
	optParams int,
	function NativeFunction,
) (err value.Value) {
	nativeMethod := NewNativeMethod(
		namespace,
		name,
		params,
		optParams,
		function,
	)
	namespace.AttachMethod(name, nativeMethod)
	return value.Undefined
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

// Define a native macro
func DefMacro(namespace types.Namespace, docComment string, name string, params []*types.Parameter, returnType types.Type, fn NativeFunction) *types.Method {
	symbolName := symbol.ToSymbol(name)
	macro := namespace.DefineMethod(
		docComment,
		types.METHOD_MACRO_FLAG,
		symbolName,
		nil,
		params,
		returnType,
		types.Never{},
	)

	macro.Body = NewNativeMethod(
		nil,
		value.S(symbolName),
		len(params),
		0,
		fn,
	)
	return macro
}

// Utility method that creates a new native
// method and attaches it to the given container.
//
// Panics when the method cannot be defined.
func Def(
	container *value.Class,
	name string,
	function NativeFunction,
	opts ...DefOption,
) {
	symbolName := value.ToSymbol(name)

	nativeMethod := &NativeMethod{
		namespace: container,
		name:      symbolName,
		Function:  function,
	}

	for _, opt := range opts {
		opt(nativeMethod)
	}

	container.AttachMethod(symbolName, nativeMethod)
}

// Utility method that defines a new bytecode
// method in the given container.
//
// Panics when the method cannot be defined.
func DefBytecode(
	namespace *value.Class,
	name string,
	body *BytecodeFunction,
) {
	symbolName := value.ToSymbol(name)
	body.namespace = namespace
	namespace.AttachMethod(symbolName, body)
}
