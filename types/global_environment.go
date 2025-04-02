package types

import (
	"github.com/elk-language/elk/token"
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
		defined:       true,
	}
	env := &GlobalEnvironment{
		Root: rootModule,
		Init: true,
	}

	stdModule := &Module{
		NamespaceBase: MakeNamespaceBase("", "Std"),
		defined:       true,
	}
	rootModule.DefineConstant(symbol.Root, rootModule)
	rootModule.DefineSubtype(symbol.Root, rootModule)

	rootModule.DefineConstant(symbol.Std, stdModule)
	rootModule.DefineSubtype(symbol.Std, stdModule)

	valueClass := &Class{
		NamespaceBase: MakeNamespaceBase("", "Std::Value"),
		defined:       true,
	}
	valueClass.primitive = true
	stdModule.DefineSubtype(symbol.Value, valueClass)

	objectClass := &Class{
		parent:        valueClass,
		NamespaceBase: MakeNamespaceBase("", "Std::Object"),
		defined:       true,
	}
	stdModule.DefineSubtype(symbol.Object, objectClass)

	classClass := &Class{
		parent:        objectClass,
		NamespaceBase: MakeNamespaceBase("", "Std::Class"),
		defined:       true,
		noinit:        true,
	}
	stdModule.DefineSubtype(symbol.Class, classClass)

	valueClass.singleton = NewSingletonClass(valueClass, classClass)
	stdModule.DefineConstant(symbol.Value, valueClass.singleton)

	objectClass.singleton = NewSingletonClass(objectClass, classClass)
	stdModule.DefineConstant(symbol.Object, objectClass.singleton)

	classClass.singleton = NewSingletonClass(classClass, classClass)
	stdModule.DefineConstant(symbol.Class, classClass.singleton)

	// -- End of Bootstrapping --

	moduleClass := stdModule.DefineClass("", false, false, false, true, symbol.Module, objectClass, env)
	rootModule.parent = moduleClass
	stdModule.parent = moduleClass

	stdModule.DefineClass("", false, false, false, true, symbol.Mixin, objectClass, env)
	stdModule.DefineClass("", false, false, false, true, symbol.Interface, objectClass, env)

	stdModule.DefineModule("", symbol.Kernel, env)

	boolClass := stdModule.DefineClass("", false, true, true, true, symbol.Bool, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.True, boolClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.False, boolClass, env)

	stdModule.DefineClass("", false, true, true, true, symbol.Nil, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.String, objectClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Symbol, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Char, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Float, valueClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.BigFloat, objectClass, env)
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
	stdModule.DefineClass("", false, true, true, false, symbol.ArrayList, objectClass, env)
	stdModule.DefineClass("", false, true, true, false, symbol.ArrayTuple, objectClass, env)
	stdModule.DefineClass("", false, true, true, false, symbol.HashMap, objectClass, env)
	stdModule.DefineClass("", false, true, true, false, symbol.HashRecord, objectClass, env)
	stdModule.DefineClass("", false, true, true, false, symbol.HashSet, objectClass, env)
	stdModule.DefineClass("", false, true, true, false, symbol.Regex, objectClass, env)
	stdModule.DefineClass("", false, true, true, true, symbol.Method, objectClass, env)
	stdModule.DefineClass("", false, true, true, false, symbol.Pair, objectClass, env)

	env.Init = false
	return env
}

// Create a new global environment for type checking.
func NewGlobalEnvironment() *GlobalEnvironment {
	env := NewGlobalEnvironmentWithoutHeaders()

	env.Init = true
	setupGlobalEnvironmentFromHeaders(env)
	env.Init = false

	setupHelperTypes(env)

	return env
}

func setupHelperTypes(env *GlobalEnvironment) {
	ArrayList := env.StdSubtypeClass(symbol.ArrayList)
	ArrayTuple := env.StdSubtypeClass(symbol.ArrayTuple)
	HashSet := env.StdSubtypeClass(symbol.HashSet)
	HashMap := env.StdSubtypeClass(symbol.HashMap)
	HashRecord := env.StdSubtypeClass(symbol.HashRecord)

	Int := env.StdSubtypeClass(symbol.Int)
	Int64 := env.StdSubtypeClass(symbol.Int64)
	Int32 := env.StdSubtypeClass(symbol.Int32)
	Int16 := env.StdSubtypeClass(symbol.Int16)
	Int8 := env.StdSubtypeClass(symbol.Int8)
	UInt64 := env.StdSubtypeClass(symbol.UInt64)
	UInt32 := env.StdSubtypeClass(symbol.UInt32)
	UInt16 := env.StdSubtypeClass(symbol.UInt16)
	UInt8 := env.StdSubtypeClass(symbol.UInt8)
	Float := env.StdSubtypeClass(symbol.Float)
	Float64 := env.StdSubtypeClass(symbol.Float64)
	Float32 := env.StdSubtypeClass(symbol.Float32)
	BigFloat := env.StdSubtypeClass(symbol.BigFloat)
	String := env.StdSubtypeClass(symbol.String)
	Char := env.StdSubtypeClass(symbol.Char)
	Regex := env.StdSubtypeClass(symbol.Regex)
	ClosedRange := env.StdSubtypeClass(symbol.ClosedRange)
	OpenRange := env.StdSubtypeClass(symbol.OpenRange)
	LeftOpenRange := env.StdSubtypeClass(symbol.LeftOpenRange)
	RightOpenRange := env.StdSubtypeClass(symbol.RightOpenRange)
	Channel := env.StdSubtypeClass(symbol.Channel)

	BuiltinAddable := NewUnion(
		Int,
		Int64,
		Int32,
		Int16,
		Int8,
		UInt64,
		UInt32,
		UInt16,
		UInt8,
		Float,
		Float64,
		Float32,
		BigFloat,
		String,
		Char,
		Regex,
		ArrayList,
		ArrayTuple,
	)
	stdModule := env.Std()
	stdModule.DefineSubtype(symbol.S_BuiltinAddable, BuiltinAddable)

	BuiltinSubtractable := NewUnion(
		Int,
		Int64,
		Int32,
		Int16,
		Int8,
		UInt64,
		UInt32,
		UInt16,
		UInt8,
		Float,
		Float64,
		Float32,
		BigFloat,
		String,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinSubtractable, BuiltinSubtractable)

	BuiltinMultipliable := NewUnion(
		Int,
		Int64,
		Int32,
		Int16,
		Int8,
		UInt64,
		UInt32,
		UInt16,
		UInt8,
		Float,
		Float64,
		Float32,
		BigFloat,
		String,
		Char,
		Regex,
		ArrayList,
		ArrayTuple,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinMultipliable, BuiltinMultipliable)

	BuiltinDividable := NewUnion(
		Int,
		Int64,
		Int32,
		Int16,
		Int8,
		UInt64,
		UInt32,
		UInt16,
		UInt8,
		Float,
		Float64,
		Float32,
		BigFloat,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinDividable, BuiltinDividable)

	BuiltinNumeric := NewUnion(
		Int,
		Int64,
		Int32,
		Int16,
		Int8,
		UInt64,
		UInt32,
		UInt16,
		UInt8,
		Float,
		Float64,
		Float32,
		BigFloat,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinNumeric, BuiltinNumeric)

	BuiltinIncrementable := NewUnion(
		Int,
		Int64,
		Int32,
		Int16,
		Int8,
		UInt64,
		UInt32,
		UInt16,
		UInt8,
		Char,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinIncrementable, BuiltinIncrementable)

	BuiltinInt := NewUnion(
		Int,
		Int64,
		Int32,
		Int16,
		Int8,
		UInt64,
		UInt32,
		UInt16,
		UInt8,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinInt, BuiltinInt)

	BuiltinLogicBitshiftable := NewUnion(
		Int64,
		Int32,
		Int16,
		Int8,
		UInt64,
		UInt32,
		UInt16,
		UInt8,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinLogicBitshiftable, BuiltinLogicBitshiftable)

	BuiltinEquatable := NewUnion(
		Int,
		Int64,
		Int32,
		Int16,
		Int8,
		UInt64,
		UInt32,
		UInt16,
		UInt8,
		Float,
		Float64,
		Float32,
		BigFloat,
		String,
		Char,
		Regex,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinEquatable, BuiltinEquatable)

	BuiltinIterable := NewUnion(
		env.StdSubtype(symbol.String),
		String.MustSubtype("ByteIterator"),
		String.MustSubtype("CharIterator"),
		String.MustSubtype("GraphemeIterator"),
		ArrayList,
		ArrayList.MustSubtype("Iterator"),
		ArrayTuple,
		ArrayTuple.MustSubtype("Iterator"),
		HashMap,
		HashMap.MustSubtype("Iterator"),
		HashRecord,
		HashRecord.MustSubtype("Iterator"),
		HashSet,
		HashSet.MustSubtype("Iterator"),
		Channel,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinIterable, BuiltinIterable)

	BuiltinIterator := NewUnion(
		String.MustSubtype("ByteIterator"),
		String.MustSubtype("CharIterator"),
		String.MustSubtype("GraphemeIterator"),
		ArrayList.MustSubtype("Iterator"),
		ArrayTuple.MustSubtype("Iterator"),
		HashMap.MustSubtype("Iterator"),
		HashRecord.MustSubtype("Iterator"),
		HashSet.MustSubtype("Iterator"),
		ClosedRange.MustSubtype("Iterator"),
		OpenRange.MustSubtype("Iterator"),
		LeftOpenRange.MustSubtype("Iterator"),
		RightOpenRange.MustSubtype("Iterator"),
		Channel,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinIterator, BuiltinIterator)

	BuiltinSubscriptable := NewUnion(
		ArrayList,
		ArrayTuple,
		HashMap,
		HashRecord,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinSubscriptable, BuiltinSubscriptable)

	ElkTokenConstant, _ := env.StdSubtypeModule(symbol.Elk).Subtype(symbol.Token)
	ElkTokenClass := ElkTokenConstant.Type.(*Class)
	for _, tokenName := range token.Types() {
		ElkTokenClass.DefineConstant(value.ToSymbol(tokenName), UInt16)
	}
}

func (g *GlobalEnvironment) DeepCopyEnv() *GlobalEnvironment {
	newRoot := &Module{
		NamespaceBase: MakeNamespaceBase("", "Root"),
		defined:       true,
	}
	newEnv := &GlobalEnvironment{
		Init: g.Init,
		Root: newRoot,
	}
	newRoot.deepCopyInPlace(g.Root, g, newEnv)
	return newEnv
}
