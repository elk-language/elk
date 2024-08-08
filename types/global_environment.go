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
		FullyChecked:  true,
	}
	valueClass.primitive = true
	stdModule.DefineSubtype(symbol.Value, valueClass)

	objectClass := &Class{
		parent:        valueClass,
		NamespaceBase: MakeNamespaceBase("", "Std::Object"),
		FullyChecked:  true,
	}
	stdModule.DefineSubtype(symbol.Object, objectClass)

	classClass := &Class{
		parent:        objectClass,
		NamespaceBase: MakeNamespaceBase("", "Std::Class"),
		FullyChecked:  true,
	}
	stdModule.DefineSubtype(symbol.Class, classClass)

	valueClass.singleton = NewSingletonClass(valueClass, classClass)
	stdModule.DefineConstant(symbol.Value, valueClass.singleton)

	objectClass.singleton = NewSingletonClass(objectClass, classClass)
	stdModule.DefineConstant(symbol.Object, objectClass.singleton)

	classClass.singleton = NewSingletonClass(classClass, classClass)
	stdModule.DefineConstant(symbol.Class, classClass.singleton)

	// -- End of Bootstrapping --

	stdModule.DefineClass("", false, false, false, true, symbol.Module, objectClass, env)
	stdModule.DefineClass("", false, false, false, true, symbol.Mixin, objectClass, env)
	stdModule.DefineClass("", false, false, false, true, symbol.Interface, objectClass, env)

	boolClass := stdModule.DefineClass("", false, true, true, true, symbol.Bool, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.True, boolClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.False, boolClass, env)

	stdModule.DefineClass("", false, true, true, true, symbol.Nil, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.String, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Symbol, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Char, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Float, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.BigFloat, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Float64, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Float32, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Int, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Int64, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Int32, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Int16, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Int8, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.UInt64, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.UInt32, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.UInt16, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.UInt8, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.ArrayList, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.ArrayTuple, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.HashMap, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.HashRecord, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.HashSet, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Regex, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Method, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Pair, valueClass, env)

	return env
}

// Create a new global environment for type checking.
func NewGlobalEnvironment() *GlobalEnvironment {
	env := NewGlobalEnvironmentWithoutHeaders()
	setupGlobalEnvironmentFromHeaders(env)
	return env
}
