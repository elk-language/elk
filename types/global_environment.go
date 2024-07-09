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

func NewGlobalEnvironmentWithoutHeaders() *GlobalEnvironment {
	// -- Bootstrapping --

	rootModule := NewModule("", "Root")
	env := &GlobalEnvironment{
		Root: rootModule,
	}

	stdModule := NewModule("", "Std")
	rootModule.DefineConstant("Std", stdModule)
	rootModule.DefineSubtype("Std", stdModule)

	valueClass := &Class{
		NamespaceBase: MakeNamespaceBase("", "Std::Value"),
	}
	stdModule.DefineSubtype("Value", valueClass)

	objectClass := &Class{
		parent:        valueClass,
		NamespaceBase: MakeNamespaceBase("", "Std::Object"),
	}
	stdModule.DefineSubtype("Object", objectClass)

	classClass := &Class{
		parent:        objectClass,
		NamespaceBase: MakeNamespaceBase("", "Std::Class"),
	}
	stdModule.DefineSubtype("Class", classClass)

	valueClass.singleton = NewSingletonClass(valueClass, classClass)
	stdModule.DefineConstant("Value", valueClass.singleton)

	objectClass.singleton = NewSingletonClass(objectClass, classClass)
	stdModule.DefineConstant("Object", objectClass.singleton)

	classClass.singleton = NewSingletonClass(classClass, classClass)
	stdModule.DefineConstant("Class", classClass.singleton)

	// -- End of Bootstrapping --

	stdModule.DefineClass("", false, false, false, "Module", objectClass, env)
	stdModule.DefineClass("", false, false, false, "Mixin", objectClass, env)
	stdModule.DefineClass("", false, false, false, "Interface", objectClass, env)

	boolClass := stdModule.DefineClass("", false, true, true, "Bool", objectClass, env)
	stdModule.DefineClass("", false, true, true, "True", boolClass, env)
	stdModule.DefineClass("", false, true, true, "False", boolClass, env)

	stdModule.DefineClass("", false, true, true, "Nil", objectClass, env)
	stdModule.DefineClass("", false, true, true, "String", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Symbol", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Char", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Float", objectClass, env)
	stdModule.DefineClass("", false, true, true, "BigFloat", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Float64", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Float32", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Int", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Int64", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Int32", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Int16", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Int8", objectClass, env)
	stdModule.DefineClass("", false, true, true, "UInt64", objectClass, env)
	stdModule.DefineClass("", false, true, true, "UInt32", objectClass, env)
	stdModule.DefineClass("", false, true, true, "UInt16", objectClass, env)
	stdModule.DefineClass("", false, true, true, "UInt8", objectClass, env)
	stdModule.DefineClass("", false, true, true, "ArrayList", objectClass, env)
	stdModule.DefineClass("", false, true, true, "ArrayTuple", objectClass, env)
	stdModule.DefineClass("", false, true, true, "HashMap", objectClass, env)
	stdModule.DefineClass("", false, true, true, "HashRecord", objectClass, env)
	stdModule.DefineClass("", false, true, true, "HashSet", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Regex", objectClass, env)
	stdModule.DefineClass("", false, true, true, "Method", objectClass, env)

	return env
}

// Create a new global environment for type checking.
func NewGlobalEnvironment() *GlobalEnvironment {
	env := NewGlobalEnvironmentWithoutHeaders()
	setupGlobalEnvironmentFromHeaders(env)
	return env
}
