package types

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Module struct {
	defined bool
	parent  Namespace
	NamespaceBase
}

func (m *Module) IsGeneric() bool {
	return false
}

func (m *Module) TypeParameters() []*TypeParameter {
	return nil
}

func (m *Module) SetTypeParameters(t []*TypeParameter) {
	panic("cannot set type parameters on a module")
}

func (*Module) Singleton() *SingletonClass {
	return nil
}

func (m *Module) SetSingleton(*SingletonClass) {
	panic("cannot set singleton class of a module")
}

func (m *Module) Parent() Namespace {
	return m.parent
}

func (m *Module) SetParent(parent Namespace) {
	m.parent = parent
}

func (m *Module) IsDefined() bool {
	return m.defined
}

func (m *Module) SetDefined(defined bool) {
	m.defined = defined
}

func (m *Module) IsAbstract() bool {
	return false
}

func (m *Module) IsSealed() bool {
	return false
}

func (m *Module) IsPrimitive() bool {
	return false
}

func NewModule(docComment, name string, env *GlobalEnvironment) *Module {
	return &Module{
		defined:       env.Init,
		parent:        env.StdSubtypeClass(symbol.Module),
		NamespaceBase: MakeNamespaceBase(docComment, name),
	}
}

func NewModuleWithDetails(
	docComment string,
	name string,
	consts ConstantMap,
	subtypes ConstantMap,
	methods MethodMap,
	env *GlobalEnvironment,
) *Module {
	return &Module{
		parent: env.StdSubtypeClass(symbol.Module),
		NamespaceBase: NamespaceBase{
			docComment: docComment,
			name:       name,
			constants:  consts,
			subtypes:   subtypes,
			methods:    methods,
		},
	}
}

func (m *Module) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}

func (*Module) IsLiteral() bool {
	return false
}

func (m *Module) inspect() string {
	return m.Name()
}

func (m *Module) DefineMethod(docComment string, flags bitfield.BitFlag16, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, flags, name, typeParams, params, returnType, throwType, m)
	m.SetMethod(name, method)
	return method
}

func (m *Module) Copy() *Module {
	return &Module{
		parent:  m.parent,
		defined: m.defined,
		NamespaceBase: NamespaceBase{
			docComment: m.docComment,
			name:       m.name,
			constants:  m.constants,
			subtypes:   m.subtypes,
			methods:    m.methods,
		},
	}
}

func (m *Module) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *Module {
	moduleConstantPath := GetConstantPath(m.name)
	parentNamespace := DeepCopyNamespacePath(moduleConstantPath[:len(moduleConstantPath)-1], oldEnv, newEnv)

	if newType, ok := NameToTypeOk(m.name, newEnv); ok {
		return newType.(*Module)
	}

	newModule := &Module{
		NamespaceBase: MakeNamespaceBase(m.docComment, m.name),
		defined:       m.defined,
	}
	if parentNamespace != nil {
		parentNamespace.DefineSubtype(value.ToSymbol(moduleConstantPath[len(moduleConstantPath)-1]), newModule)
	}

	newModule.methods = MethodsDeepCopyEnv(m.methods, oldEnv, newEnv)
	newModule.instanceVariables = TypesDeepCopyEnv(m.instanceVariables, oldEnv, newEnv)
	newModule.subtypes = ConstantsDeepCopyEnv(m.subtypes, oldEnv, newEnv)
	newModule.constants = ConstantsDeepCopyEnv(m.constants, oldEnv, newEnv)

	if m.parent != nil {
		newModule.parent = DeepCopyEnv(m.parent, oldEnv, newEnv).(Namespace)
	}
	return newModule
}

func (m *Module) deepCopyInPlace(oldModule *Module, oldEnv, newEnv *GlobalEnvironment) {
	m.methods = MethodsDeepCopyEnv(oldModule.methods, oldEnv, newEnv)
	m.subtypes = ConstantsDeepCopyEnv(oldModule.subtypes, oldEnv, newEnv)
	m.constants = ConstantsDeepCopyEnv(oldModule.constants, oldEnv, newEnv)
	if m.parent != nil {
		m.parent = DeepCopyEnv(oldModule.parent, oldEnv, newEnv).(Namespace)
	}
}
