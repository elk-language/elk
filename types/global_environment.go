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
	// -- Bootstrapping --

	rootModule := NewModule("Root")
	env := &GlobalEnvironment{
		Root: rootModule,
	}

	stdModule := NewModule("Std")
	rootModule.DefineConstant("Std", stdModule)
	rootModule.DefineSubtype("Std", stdModule)

	valueClass := &Class{
		NamespaceBase: MakeNamespaceBase("Std::Value"),
	}
	stdModule.DefineSubtype("Value", valueClass)

	objectClass := &Class{
		parent:        valueClass,
		NamespaceBase: MakeNamespaceBase("Std::Object"),
	}
	stdModule.DefineSubtype("Object", objectClass)

	classClass := &Class{
		parent:        objectClass,
		NamespaceBase: MakeNamespaceBase("Std::Class"),
	}
	stdModule.DefineSubtype("Class", classClass)

	valueClass.singleton = NewSingletonClass(valueClass, classClass)
	stdModule.DefineConstant("Value", valueClass.singleton)

	objectClass.singleton = NewSingletonClass(objectClass, classClass)
	stdModule.DefineConstant("Object", objectClass.singleton)

	classClass.singleton = NewSingletonClass(classClass, classClass)
	stdModule.DefineConstant("Class", classClass.singleton)

	// -- End of Bootstrapping --

	stdModule.DefineClass("Module", objectClass, env)
	stdModule.DefineClass("Mixin", objectClass, env)
	stdModule.DefineClass("Interface", objectClass, env)

	boolClass := stdModule.DefineClass("Bool", objectClass, env)
	stdModule.DefineClass("True", boolClass, env)
	stdModule.DefineClass("False", boolClass, env)

	stdModule.DefineClass("Nil", objectClass, env)
	stdModule.DefineClass("String", objectClass, env)
	stdModule.DefineClass("Symbol", objectClass, env)
	stdModule.DefineClass("Char", objectClass, env)
	stdModule.DefineClass("Float", objectClass, env)
	stdModule.DefineClass("BigFloat", objectClass, env)
	stdModule.DefineClass("Float64", objectClass, env)
	stdModule.DefineClass("Float32", objectClass, env)
	stdModule.DefineClass("Int", objectClass, env)
	stdModule.DefineClass("Int64", objectClass, env)
	stdModule.DefineClass("Int32", objectClass, env)
	stdModule.DefineClass("Int16", objectClass, env)
	stdModule.DefineClass("Int8", objectClass, env)
	stdModule.DefineClass("UInt64", objectClass, env)
	stdModule.DefineClass("UInt32", objectClass, env)
	stdModule.DefineClass("UInt16", objectClass, env)
	stdModule.DefineClass("UInt8", objectClass, env)
	stdModule.DefineClass("ArrayList", objectClass, env)
	stdModule.DefineClass("ArrayTuple", objectClass, env)
	stdModule.DefineClass("HashMap", objectClass, env)
	stdModule.DefineClass("HashRecord", objectClass, env)
	stdModule.DefineClass("HashSet", objectClass, env)
	stdModule.DefineClass("Regex", objectClass, env)
	stdModule.DefineClass("Method", objectClass, env)

	return env
}
