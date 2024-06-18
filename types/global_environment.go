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
	rootModule := NewModule("Root")

	stdModule := NewModule("Std")
	rootModule.DefineConstant("Std", stdModule)
	rootModule.DefineSubtype("Std", stdModule)

	valueClass := stdModule.DefineClass("Value", nil)
	objectClass := stdModule.DefineClass("Object", valueClass)

	stdModule.DefineClass("Class", objectClass)
	stdModule.DefineClass("Mixin", objectClass)
	stdModule.DefineClass("Module", objectClass)

	boolClass := stdModule.DefineClass("Bool", objectClass)
	stdModule.DefineClass("True", boolClass)
	stdModule.DefineClass("False", boolClass)

	stdModule.DefineClass("Nil", objectClass)
	stdModule.DefineClass("String", objectClass)
	stdModule.DefineClass("Symbol", objectClass)
	stdModule.DefineClass("Char", objectClass)
	stdModule.DefineClass("Float", objectClass)
	stdModule.DefineClass("BigFloat", objectClass)
	stdModule.DefineClass("Float64", objectClass)
	stdModule.DefineClass("Float32", objectClass)
	stdModule.DefineClass("Int", objectClass)
	stdModule.DefineClass("Int64", objectClass)
	stdModule.DefineClass("Int32", objectClass)
	stdModule.DefineClass("Int16", objectClass)
	stdModule.DefineClass("Int8", objectClass)
	stdModule.DefineClass("UInt64", objectClass)
	stdModule.DefineClass("UInt32", objectClass)
	stdModule.DefineClass("UInt16", objectClass)
	stdModule.DefineClass("UInt8", objectClass)
	stdModule.DefineClass("ArrayList", objectClass)
	stdModule.DefineClass("ArrayTuple", objectClass)
	stdModule.DefineClass("HashMap", objectClass)
	stdModule.DefineClass("HashRecord", objectClass)
	stdModule.DefineClass("HashSet", objectClass)
	stdModule.DefineClass("Regex", objectClass)
	stdModule.DefineClass("Method", objectClass)

	return &GlobalEnvironment{
		Root: rootModule,
	}
}
