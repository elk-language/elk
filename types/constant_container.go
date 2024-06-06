package types

import "github.com/elk-language/elk/value"

type ConstantContainer interface {
	Type
	Name() string
	Parent() ConstantContainer

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

	DefineClass(name string, parent ConstantContainer) *Class
	DefineModule(name string) *Module
	DefineMixin(name string) *Mixin
}
