package types

import (
	"fmt"
	"sync"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value/symbol"
)

// Current global type environment
var Env *GlobalEnvironment

type GlobalEnvironment struct {
	Root      *Module
	TypeSet   *ds.HashSet[Type]
	TypeIndex []Type
	Init      bool // Whether the global environment is in its initialisation stage
	mu        sync.RWMutex
}

func (g *GlobalEnvironment) GetType(id ID) Type {
	if !g.Init {
		g.mu.RLock()
		defer g.mu.RUnlock()
	}

	result := g.TypeIndex[id-1]
	return result
}

func (g *GlobalEnvironment) RegisterType(t Type) Type {
	if !g.Init {
		g.mu.Lock()
		defer g.mu.Unlock()
	}

	if existingType, ok := g.TypeSet.Get(t); ok {
		t.SetID(existingType.ID())
		return existingType
	}

	id := ID(len(g.TypeIndex)) + 1
	t.SetID(id)
	g.TypeIndex = append(g.TypeIndex, t)
	g.TypeSet.Add(t)
	return t
}

func (g *GlobalEnvironment) RegisterTypeWithID(t Type, id ID) Type {
	g.mu.Lock()
	defer g.mu.Unlock()

	if existingType, ok := g.TypeSet.Get(t); ok {
		if existingType.ID() != id {
			panic(fmt.Sprintf("tried to register a type with ID=%d but the type already has ID=%d", id, existingType.ID))
		}
		return existingType
	}

	currentId := ID(len(g.TypeIndex) + 1)
	if currentId != id {
		panic(fmt.Sprintf("tried to register a type with ID=%d but the type already has ID=%d", id, currentId))
	}

	t.SetID(id)
	g.TypeIndex = append(g.TypeIndex, t)
	g.TypeSet.Add(t)
	return t
}

func (g *GlobalEnvironment) ReplaceTypeWithID(t Type, id ID) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.TypeIndex) >= int(id) {
		panic(fmt.Sprintf("cannot replace type with ID=%d when index has length=%d", id, len(g.TypeIndex)))
	}

	t.SetID(id)
	g.TypeIndex[id-1] = t
	g.TypeSet.Add(t)
}

func (g *GlobalEnvironment) NamesToType(path ...symbol.Symbol) Type {
	return PathToNestedSubtype(path, g.Root)
}

func (g *GlobalEnvironment) NamesToNamespace(path ...symbol.Symbol) Namespace {
	return PathToNestedNamespace(path, g.Root)
}

func (g *GlobalEnvironment) Std() *Module {
	s, _ := g.Root.Subtype(symbol.C_Std)
	return s.Type.Get().(*Module)
}

func (g *GlobalEnvironment) StdSubtype(name symbol.Symbol) Type {
	s, _ := g.Std().Subtype(name)
	return s.Type.Get()
}

func (g *GlobalEnvironment) StdSubtypeClass(name symbol.Symbol) *Class {
	s, _ := g.Std().Subtype(name)
	return s.Type.Get().(*Class)
}

func (g *GlobalEnvironment) StdSubtypeModule(name symbol.Symbol) *Module {
	s, _ := g.Std().Subtype(name)
	return s.Type.Get().(*Module)
}

func (g *GlobalEnvironment) StdSubtypeString(name string) Type {
	s, _ := g.Std().SubtypeString(name)
	return s.Type.Get()
}

func (g *GlobalEnvironment) StdConstString(name string) Type {
	s, _ := g.Std().ConstantString(name)
	return s.Type.Get()
}

func (g *GlobalEnvironment) StdConst(name symbol.Symbol) Type {
	s, _ := g.Std().Constant(name)
	return s.Type.Get()
}

func (g *GlobalEnvironment) ASTSubtype(name symbol.Symbol) Type {
	elkMod := g.StdSubtypeModule(symbol.C_Elk)
	astMod := elkMod.MustSubtype(symbol.C_AST).(*Module)
	return astMod.MustSubtype(name)
}

func (g *GlobalEnvironment) ExpressionNode() *Mixin {
	return g.ASTSubtype(symbol.C_ExpressionNode).(*Mixin)
}

func (g *GlobalEnvironment) PatternNode() *Mixin {
	return g.ASTSubtype(symbol.C_PatternNode).(*Mixin)
}

func (g *GlobalEnvironment) TypeNode() *Mixin {
	return g.ASTSubtype(symbol.C_TypeNode).(*Mixin)
}

func (g *GlobalEnvironment) Copy() *GlobalEnvironment {
	g.mu.Lock()
	defer g.mu.Unlock()

	newTypeIndex := make([]Type, len(g.TypeIndex))
	newTypeSet := ds.NewHashSet[Type](g.TypeSet.Capacity())
	newEnv := &GlobalEnvironment{
		TypeIndex: newTypeIndex,
		TypeSet:   newTypeSet,
	}

	for i, typ := range g.TypeIndex {
		typeCopy := typ.CopyType()
		newEnv.ReplaceTypeWithID(typeCopy, typ.ID())
		newTypeIndex[i] = typeCopy
	}
	newEnv.Root = g.TypeIndex[g.Root.id].(*Module)

	return newEnv
}

func NewGlobalEnvironmentWithoutHeaders() *GlobalEnvironment {
	origEnv := Env
	// -- Bootstrapping --

	rootModule := &Module{
		NamespaceBase: MakeNamespaceBase("", "Root"),
		native:        true,
	}
	env := &GlobalEnvironment{
		Root:    rootModule,
		TypeSet: ds.NewHashSet[Type](30),
		Init:    true,
	}
	Env = env

	env.RegisterTypeWithID(Any{}, AnyID)
	env.RegisterTypeWithID(Bool{}, BoolID)
	env.RegisterTypeWithID(False{}, FalseID)
	env.RegisterTypeWithID(True{}, TrueID)
	env.RegisterTypeWithID(Nil{}, NilID)
	env.RegisterTypeWithID(Never{}, NeverID)
	env.RegisterTypeWithID(NoValue{}, NoValueID)
	env.RegisterTypeWithID(Self{}, SelfID)
	env.RegisterTypeWithID(Untyped{}, UntypedID)
	env.RegisterTypeWithID(Void{}, VoidID)

	env.RegisterType(rootModule)

	stdModule := &Module{
		NamespaceBase: MakeNamespaceBase("", "Std"),
		native:        true,
	}
	env.RegisterType(stdModule)
	rootModule.DefineConstant(symbol.C_Root, rootModule)
	rootModule.DefineSubtype(symbol.C_Root, rootModule)

	rootModule.DefineConstant(symbol.C_Std, stdModule)
	rootModule.DefineSubtype(symbol.C_Std, stdModule)

	valueClass := &Class{
		NamespaceBase: MakeNamespaceBase("", "Std::Value"),
		native:        true,
		immutable:     true,
	}
	valueClass.primitive = true
	env.RegisterType(valueClass)
	stdModule.DefineSubtype(symbol.C_Value, valueClass)

	objectClass := &Class{
		parent:        CastRef[Namespace](valueClass),
		NamespaceBase: MakeNamespaceBase("", "Std::Object"),
		native:        true,
		immutable:     true,
	}
	env.RegisterType(objectClass)
	stdModule.DefineSubtype(symbol.C_Object, objectClass)

	classClass := &Class{
		parent:        CastRef[Namespace](objectClass),
		NamespaceBase: MakeNamespaceBase("", "Std::Class"),
		native:        true,
		noinit:        true,
	}

	env.RegisterType(classClass)
	stdModule.DefineSubtype(symbol.C_Class, classClass)

	valueSingleton := NewSingletonClass(valueClass, classClass)
	valueClass.singleton = valueSingleton.ToRef()
	stdModule.DefineConstant(symbol.C_Value, valueSingleton)

	objectSingleton := NewSingletonClass(objectClass, classClass)
	objectClass.singleton = objectSingleton.ToRef()
	stdModule.DefineConstant(symbol.C_Object, objectSingleton)

	classSingleton := NewSingletonClass(classClass, classClass)
	classClass.singleton = classSingleton.ToRef()
	stdModule.DefineConstant(symbol.C_Class, classSingleton)

	// -- End of Bootstrapping --

	moduleClass := stdModule.DefineClass("", false, false, false, true, false, symbol.C_Module, objectClass)
	rootModule.parent = moduleClass
	stdModule.parent = moduleClass

	stdModule.DefineClass("", false, false, false, true, false, symbol.C_Mixin, objectClass)
	stdModule.DefineClass("", false, false, false, true, false, symbol.C_Interface, objectClass)

	stdModule.DefineModule("", symbol.C_Kernel)

	boxClass := stdModule.DefineClass("", false, true, true, false, false, symbol.C_Box, objectClass)
	// Set up type parameters
	var typeParam *TypeParameter
	typeParams := make([]Ref[*TypeParameter], 1)
	typeParam = NewTypeParameter(symbol.ToSymbol("Val"), boxClass, Never{}, Any{}, nil, INVARIANT)
	typeParams[0] = typeParam.ToRef()
	boxClass.DefineSubtype(symbol.ToSymbol("Val"), typeParam)
	boxClass.DefineConstant(symbol.ToSymbol("Val"), NoValue{})
	boxClass.SetTypeParameters(typeParams)

	boolClass := stdModule.DefineClass("", false, true, true, true, false, symbol.C_Bool, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_True, boolClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_False, boolClass)

	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Nil, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_String, objectClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Symbol, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Char, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Float, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_BigFloat, objectClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Float64, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Float32, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Int, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Int64, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Int32, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Int16, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Int8, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_UInt64, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_UInt32, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_UInt16, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_UInt8, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_UInt, valueClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_ArrayList, objectClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_ArrayTuple, objectClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_HashMap, objectClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_HashRecord, objectClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_HashSet, objectClass)
	stdModule.DefineClass("", false, true, true, false, false, symbol.C_Regex, objectClass)
	stdModule.DefineClass("", false, true, true, true, false, symbol.C_Method, objectClass)
	stdModule.DefineClass("", false, true, true, false, false, symbol.C_Pair, objectClass)

	tupleMixin := stdModule.DefineMixin("", true, symbol.C_Tuple)
	typeParams = make([]Ref[*TypeParameter], 1)
	typeParam = NewTypeParameter(symbol.ToSymbol("Val"), tupleMixin, Never{}, Any{}, nil, COVARIANT)
	typeParams[0] = typeParam.ToRef()
	tupleMixin.DefineSubtype(symbol.ToSymbol("Val"), typeParam)
	tupleMixin.DefineConstant(symbol.ToSymbol("Val"), NoValue{})
	tupleMixin.SetTypeParameters(typeParams)

	listMixin := stdModule.DefineMixin("", true, symbol.C_List)
	typeParams = make([]Ref[*TypeParameter], 1)
	typeParam = NewTypeParameter(symbol.ToSymbol("Val"), listMixin, Never{}, Any{}, nil, INVARIANT)
	typeParams[0] = typeParam.ToRef()
	listMixin.DefineSubtype(symbol.ToSymbol("Val"), typeParam)
	listMixin.DefineConstant(symbol.ToSymbol("Val"), NoValue{})
	listMixin.SetTypeParameters(typeParams)

	env.Init = false
	Env = origEnv
	return env
}

// Create a new global environment for type checking.
func NewGlobalEnvironment() *GlobalEnvironment {
	env := NewGlobalEnvironmentWithoutHeaders()
	origEnv := Env
	Env = env

	env.Init = true
	setupGlobalEnvironmentFromHeaders(env)
	env.Init = false

	setupHelperTypes(env)

	Env = origEnv
	return env
}

func setupHelperTypes(env *GlobalEnvironment) {
	ArrayList := env.StdSubtypeClass(symbol.C_ArrayList)
	ArrayTuple := env.StdSubtypeClass(symbol.C_ArrayTuple)
	HashSet := env.StdSubtypeClass(symbol.C_HashSet)
	HashMap := env.StdSubtypeClass(symbol.C_HashMap)
	HashRecord := env.StdSubtypeClass(symbol.C_HashRecord)

	Int := env.StdSubtypeClass(symbol.C_Int)
	Int64 := env.StdSubtypeClass(symbol.C_Int64)
	Int32 := env.StdSubtypeClass(symbol.C_Int32)
	Int16 := env.StdSubtypeClass(symbol.C_Int16)
	Int8 := env.StdSubtypeClass(symbol.C_Int8)
	UInt64 := env.StdSubtypeClass(symbol.C_UInt64)
	UInt32 := env.StdSubtypeClass(symbol.C_UInt32)
	UInt16 := env.StdSubtypeClass(symbol.C_UInt16)
	UInt8 := env.StdSubtypeClass(symbol.C_UInt8)
	Float := env.StdSubtypeClass(symbol.C_Float)
	Float64 := env.StdSubtypeClass(symbol.C_Float64)
	Float32 := env.StdSubtypeClass(symbol.C_Float32)
	BigFloat := env.StdSubtypeClass(symbol.C_BigFloat)
	String := env.StdSubtypeClass(symbol.C_String)
	Char := env.StdSubtypeClass(symbol.C_Char)
	Regex := env.StdSubtypeClass(symbol.C_Regex)
	ClosedRange := env.StdSubtypeClass(symbol.C_ClosedRange)
	OpenRange := env.StdSubtypeClass(symbol.C_OpenRange)
	LeftOpenRange := env.StdSubtypeClass(symbol.C_LeftOpenRange)
	RightOpenRange := env.StdSubtypeClass(symbol.C_RightOpenRange)
	Channel := env.StdSubtypeClass(symbol.C_Channel)

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
		env.StdSubtype(symbol.C_String),
		String.MustSubtypeString("ByteIterator"),
		String.MustSubtypeString("CharIterator"),
		String.MustSubtypeString("GraphemeIterator"),
		ArrayList,
		ArrayList.MustSubtypeString("Iterator"),
		ArrayTuple,
		ArrayTuple.MustSubtypeString("Iterator"),
		HashMap,
		HashMap.MustSubtypeString("Iterator"),
		HashRecord,
		HashRecord.MustSubtypeString("Iterator"),
		HashSet,
		HashSet.MustSubtypeString("Iterator"),
		Channel,
	)
	stdModule.DefineSubtype(symbol.S_BuiltinIterable, BuiltinIterable)

	BuiltinIterator := NewUnion(
		String.MustSubtypeString("ByteIterator"),
		String.MustSubtypeString("CharIterator"),
		String.MustSubtypeString("GraphemeIterator"),
		ArrayList.MustSubtypeString("Iterator"),
		ArrayTuple.MustSubtypeString("Iterator"),
		HashMap.MustSubtypeString("Iterator"),
		HashRecord.MustSubtypeString("Iterator"),
		HashSet.MustSubtypeString("Iterator"),
		ClosedRange.MustSubtypeString("Iterator"),
		OpenRange.MustSubtypeString("Iterator"),
		LeftOpenRange.MustSubtypeString("Iterator"),
		RightOpenRange.MustSubtypeString("Iterator"),
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

	ElkTokenConstant, _ := env.StdSubtypeModule(symbol.C_Elk).Subtype(symbol.C_Token)
	ElkTokenClass := ElkTokenConstant.Type.Get().(*Class)
	for _, tokenName := range token.Types() {
		ElkTokenClass.DefineConstant(symbol.ToSymbol(tokenName), UInt16)
	}
}
