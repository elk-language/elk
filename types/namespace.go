package types

import "github.com/elk-language/elk/value"

type Namespace interface {
	Type
	Name() string
	Parent() Namespace
	SetParent(Namespace)
	IsAbstract() bool
	IsSealed() bool

	Constants() *TypeMap
	Constant(name value.Symbol) Type
	ConstantString(name string) Type
	DefineConstant(name string, val Type)

	Subtypes() *TypeMap
	Subtype(name value.Symbol) Type
	SubtypeString(name string) Type
	DefineSubtype(name string, val Type)

	Methods() *MethodMap
	Method(name value.Symbol) *Method
	MethodString(name string) *Method
	DefineMethod(name string, params []*Parameter, returnType, throwType Type) *Method
	SetMethod(name string, method *Method)

	InstanceVariables() *TypeMap
	InstanceVariable(name value.Symbol) Type
	InstanceVariableString(name string) Type
	DefineInstanceVariable(name string, val Type)

	DefineClass(name string, parent Namespace) *Class
	DefineModule(name string) *Module
	DefineMixin(name string) *Mixin
	DefineInterface(name string) *Interface
}