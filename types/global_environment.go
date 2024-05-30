package types

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type GlobalEnvironment struct {
	Root *Module
}

func (g *GlobalEnvironment) Std() *Module {
	return g.Root.Subtype(symbol.Std).(*Module)
}

func (g *GlobalEnvironment) StdSubtype(name value.Symbol) Type {
	return g.Std().Subtype(name)
}

func (g *GlobalEnvironment) StdSubtypeClass(name value.Symbol) *Class {
	return g.Std().Subtype(name).(*Class)
}

func (g *GlobalEnvironment) StdSubtypeString(name string) Type {
	return g.Std().SubtypeString(name)
}

func (g *GlobalEnvironment) StdConstString(name string) Type {
	return g.Std().ConstantString(name)
}

func (g *GlobalEnvironment) StdConst(name value.Symbol) Type {
	return g.Std().Constant(name)
}

// Create a new global environment for type checking.
func NewGlobalEnvironment() *GlobalEnvironment {
	rootModule := NewModule("Root", nil, nil, nil)

	stdModule := NewModule("Std", nil, nil, nil)
	rootModule.DefineConstant("Std", stdModule)
	rootModule.DefineSubtype("Std", stdModule)

	valueClass := stdModule.DefineClass("Value", nil, nil, nil)
	objectClass := stdModule.DefineClass("Object", valueClass, nil, nil)

	stdModule.DefineClass("Class", objectClass, nil, nil)
	stdModule.DefineClass("Mixin", objectClass, nil, nil)
	stdModule.DefineClass("Module", objectClass, nil, nil)

	boolClass := stdModule.DefineClass("Bool", objectClass, nil, nil)
	stdModule.DefineClass("True", boolClass, nil, nil)
	stdModule.DefineClass("False", boolClass, nil, nil)

	stdModule.DefineClass("Nil", objectClass, nil, nil)
	stdModule.DefineClass("String", objectClass, nil, nil)
	stdModule.DefineClass("Symbol", objectClass, nil, nil)
	stdModule.DefineClass("Char", objectClass, nil, nil)
	stdModule.DefineClass("Float", objectClass, nil, nil)
	stdModule.DefineClass("BigFloat", objectClass, nil, nil)
	stdModule.DefineClass("Float64", objectClass, nil, nil)
	stdModule.DefineClass("Float32", objectClass, nil, nil)
	stdModule.DefineClass("Int", objectClass, nil, nil)
	stdModule.DefineClass("Int64", objectClass, nil, nil)
	stdModule.DefineClass("Int32", objectClass, nil, nil)
	stdModule.DefineClass("Int16", objectClass, nil, nil)
	stdModule.DefineClass("Int8", objectClass, nil, nil)
	stdModule.DefineClass("UInt64", objectClass, nil, nil)
	stdModule.DefineClass("UInt32", objectClass, nil, nil)
	stdModule.DefineClass("UInt16", objectClass, nil, nil)
	stdModule.DefineClass("UInt8", objectClass, nil, nil)
	stdModule.DefineClass("ArrayList", objectClass, nil, nil)
	stdModule.DefineClass("ArrayTuple", objectClass, nil, nil)
	stdModule.DefineClass("HashMap", objectClass, nil, nil)
	stdModule.DefineClass("HashRecord", objectClass, nil, nil)
	stdModule.DefineClass("HashSet", objectClass, nil, nil)
	stdModule.DefineClass("Regex", objectClass, nil, nil)
	stdModule.DefineClass("Method", objectClass, nil, nil)

	return &GlobalEnvironment{
		Root: rootModule,
	}
}
