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
	rootModule.DefineConstant(symbol.Std, stdModule)
	rootModule.DefineSubtype(symbol.Std, stdModule)

	valueClass := &Class{
		NamespaceBase: MakeNamespaceBase("", "Std::Value"),
	}
	stdModule.DefineSubtype(symbol.Value, valueClass)

	objectClass := &Class{
		parent:        valueClass,
		NamespaceBase: MakeNamespaceBase("", "Std::Object"),
	}
	stdModule.DefineSubtype(symbol.Object, objectClass)

	classClass := &Class{
		parent:        objectClass,
		NamespaceBase: MakeNamespaceBase("", "Std::Class"),
	}
	stdModule.DefineSubtype(symbol.Class, classClass)

	valueClass.singleton = NewSingletonClass(valueClass, classClass)
	stdModule.DefineConstant(symbol.Value, valueClass.singleton)

	objectClass.singleton = NewSingletonClass(objectClass, classClass)
	stdModule.DefineConstant(symbol.Object, objectClass.singleton)

	classClass.singleton = NewSingletonClass(classClass, classClass)
	stdModule.DefineConstant(symbol.Class, classClass.singleton)

	// -- End of Bootstrapping --

	stdModule.DefineClass("", false, false, false, symbol.Module, objectClass, env)
	stdModule.DefineClass("", false, false, false, symbol.Mixin, objectClass, env)
	stdModule.DefineClass("", false, false, false, symbol.Interface, objectClass, env)

	boolClass := stdModule.DefineClass("", false, true, true, symbol.Bool, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.True, boolClass, env)
	stdModule.DefineClass("", false, true, true, symbol.False, boolClass, env)

	stdModule.DefineClass("", false, true, true, symbol.Nil, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.String, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Symbol, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Char, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Float, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.BigFloat, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Float64, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Float32, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Int, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Int64, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Int32, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Int16, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Int8, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.UInt64, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.UInt32, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.UInt16, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.UInt8, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.ArrayList, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.ArrayTuple, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.HashMap, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.HashRecord, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.HashSet, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Regex, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Method, objectClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Pair, objectClass, env)

	return env
}

// Create a new global environment for type checking.
func NewGlobalEnvironment() *GlobalEnvironment {
	env := NewGlobalEnvironmentWithoutHeaders()
	setupGlobalEnvironmentFromHeaders(env)
	return env
}
