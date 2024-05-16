package types

import "github.com/elk-language/elk/value"

type ConstantContainer interface {
	Type
	Name() string

	Constants() map[value.Symbol]Type
	Constant(name value.Symbol) Type
	ConstantString(name string) Type
	DefineConstant(name string, val Type)

	Subtypes() map[value.Symbol]Type
	Subtype(name value.Symbol) Type
	SubtypeString(name string) Type
	DefineSubtype(name string, val Type)

	Methods() MethodMap
	Method(name value.Symbol) *Method
	MethodString(name string) *Method
	DefineMethod(name string, params []*Parameter, returnType, throwType Type) *Method
	SetMethod(name string, method *Method)

	DefineClass(name string, parent *Class, consts map[value.Symbol]Type) *Class
	DefineModule(name string, consts map[value.Symbol]Type, subtypes map[value.Symbol]Type) *Module
}
