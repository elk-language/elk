package value

// Represents an Elk namespace for constants.
type Namespace interface {
	Name() string
	Constants() SymbolMap
	AddConstantString(name string, val Value)
	AddConstant(name Symbol, val Value)
}

// Represents an Elk namespace for constants and methods.
type MethodNamespace interface {
	Namespace
	Superclass() *Class
	Methods() MethodMap
	LookupMethod(name Symbol) Method
	AttachMethod(name Symbol, method Method)
	DefineAlias(newMethodName, oldMethodName Symbol)
	DefineAliasString(newMethodName, oldMethodName string)
}

var NamespaceClass *Class // ::Std::Namespace
