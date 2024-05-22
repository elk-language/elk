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

	valueClass := stdModule.DefineClass("Value", nil, nil)
	objectClass := stdModule.DefineClass("Object", valueClass, nil)

	stdModule.DefineClass("Class", objectClass, nil)
	stdModule.DefineClass("Mixin", objectClass, nil)
	stdModule.DefineClass("Module", objectClass, nil)

	boolClass := stdModule.DefineClass("Bool", objectClass, nil)
	stdModule.DefineClass("True", boolClass, nil)
	stdModule.DefineClass("False", boolClass, nil)

	stdModule.DefineClass("Nil", objectClass, nil)
	stdModule.DefineClass("String", objectClass, nil)
	stdModule.DefineClass("Symbol", objectClass, nil)
	stdModule.DefineClass("Char", objectClass, nil)
	stdModule.DefineClass("Float", objectClass, nil)
	stdModule.DefineClass("BigFloat", objectClass, nil)
	stdModule.DefineClass("Float64", objectClass, nil)
	stdModule.DefineClass("Float32", objectClass, nil)
	stdModule.DefineClass("Int", objectClass, nil)
	stdModule.DefineClass("Int64", objectClass, nil)
	stdModule.DefineClass("Int32", objectClass, nil)
	stdModule.DefineClass("Int16", objectClass, nil)
	stdModule.DefineClass("Int8", objectClass, nil)
	stdModule.DefineClass("UInt64", objectClass, nil)
	stdModule.DefineClass("UInt32", objectClass, nil)
	stdModule.DefineClass("UInt16", objectClass, nil)
	stdModule.DefineClass("UInt8", objectClass, nil)
	stdModule.DefineClass("ArrayList", objectClass, nil)
	stdModule.DefineClass("ArrayTuple", objectClass, nil)
	stdModule.DefineClass("HashMap", objectClass, nil)
	stdModule.DefineClass("HashRecord", objectClass, nil)
	stdModule.DefineClass("HashSet", objectClass, nil)
	stdModule.DefineClass("Regex", objectClass, nil)
	stdModule.DefineClass("Method", objectClass, nil)

	return &GlobalEnvironment{
		Root: rootModule,
	}
}
