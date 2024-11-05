package types

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type GlobalEnvironment struct {
	Root *Module
	Init bool // Whether the global environment is in its initialisation stage
}

func (g *GlobalEnvironment) Std() *Module {
	s, _ := g.Root.Subtype(symbol.Std)
	return s.Type.(*Module)
}

func (g *GlobalEnvironment) StdSubtype(name value.Symbol) Type {
	s, _ := g.Std().Subtype(name)
	return s.Type
}

func (g *GlobalEnvironment) StdSubtypeClass(name value.Symbol) *Class {
	s, _ := g.Std().Subtype(name)
	return s.Type.(*Class)
}

func (g *GlobalEnvironment) StdSubtypeModule(name value.Symbol) *Module {
	s, _ := g.Std().Subtype(name)
	return s.Type.(*Module)
}

func (g *GlobalEnvironment) StdSubtypeString(name string) Type {
	s, _ := g.Std().SubtypeString(name)
	return s.Type
}

func (g *GlobalEnvironment) StdConstString(name string) Type {
	s, _ := g.Std().ConstantString(name)
	return s.Type
}

func (g *GlobalEnvironment) StdConst(name value.Symbol) Type {
	s, _ := g.Std().Constant(name)
	return s.Type
}

func NewGlobalEnvironmentWithoutHeaders() *GlobalEnvironment {
	// -- Bootstrapping --

	rootModule := &Module{
		NamespaceBase: MakeNamespaceBase("", "Root"),
	}
	env := &GlobalEnvironment{
		Root: rootModule,
		Init: true,
	}

	stdModule := &Module{
		NamespaceBase: MakeNamespaceBase("", "Std"),
	}
	rootModule.DefineConstant(symbol.Std, stdModule)
	rootModule.DefineSubtype(symbol.Std, stdModule)

	valueClass := &Class{
		NamespaceBase: MakeNamespaceBase("", "Std::Value"),
	}
	valueClass.primitive = true
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

	moduleClass := stdModule.DefineClass("", false, false, false, symbol.Module, objectClass, env)
	rootModule.parent = moduleClass
	stdModule.parent = moduleClass

	stdModule.DefineClass("", false, false, false, symbol.Mixin, objectClass, env)
	stdModule.DefineClass("", false, false, false, symbol.Interface, objectClass, env)

	stdModule.DefineModule("", symbol.Kernel, env)

	boolClass := stdModule.DefineClass("", false, true, true, symbol.Bool, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.True, boolClass, env)
	stdModule.DefineClass("", false, true, true, symbol.False, boolClass, env)

	stdModule.DefineClass("", false, true, true, symbol.Nil, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.String, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Symbol, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Char, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Float, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.BigFloat, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Float64, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Float32, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Int, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Int64, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Int32, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Int16, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Int8, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.UInt64, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.UInt32, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.UInt16, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.UInt8, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.ArrayList, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.ArrayTuple, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.HashMap, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.HashRecord, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.HashSet, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Regex, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Method, valueClass, env)
	stdModule.DefineClass("", false, true, true, symbol.Pair, valueClass, env)

	env.Init = false
	return env
}

// Create a new global environment for type checking.
func NewGlobalEnvironment() *GlobalEnvironment {
	env := NewGlobalEnvironmentWithoutHeaders()

	env.Init = true
	setupGlobalEnvironmentFromHeaders(env)
	env.Init = false

	return env
}
