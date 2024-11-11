package types

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Module struct {
	compiled bool
	parent   Namespace
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

func (m *Module) Parent() Namespace {
	return m.parent
}

func (m *Module) SetParent(parent Namespace) {
	m.parent = parent
}

func (m *Module) IsDefined() bool {
	return m.compiled
}

func (m *Module) SetDefined(compiled bool) {
	m.compiled = compiled
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
		compiled:      env.Init,
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

func (m *Module) DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, abstract, sealed, native, name, typeParams, params, returnType, throwType, m)
	m.SetMethod(name, method)
	return method
}

func (m *Module) Copy() *Module {
	return &Module{
		parent: m.parent,
		NamespaceBase: NamespaceBase{
			docComment: m.docComment,
			name:       m.name,
			constants:  m.constants,
			subtypes:   m.subtypes,
			methods:    m.methods,
		},
	}
}

func (m *Module) DeepCopy(oldEnv, newEnv *GlobalEnvironment) *Module {
	if newType, ok := NameToTypeOk(m.name, newEnv); ok {
		return newType.(*Module)
	}

	newModule := m.Copy()
	moduleConstantPath := GetConstantPath(m.name)
	parentNamespace := DeepCopyNamespacePath(moduleConstantPath[:len(moduleConstantPath)-1], oldEnv, newEnv)
	parentNamespace.DefineSubtype(value.ToSymbol(moduleConstantPath[len(moduleConstantPath)-1]), newModule)

	newMethods := make(MethodMap, len(m.methods))
	for methodName, method := range m.methods {
		newMethods[methodName] = method.Copy()
	}
	newModule.methods = newMethods

	newConstants := make(ConstantMap, len(m.constants))
	for constName, constant := range m.constants {
		newConstants[constName] = Constant{
			FullName: constant.FullName,
			Type:     DeepCopy(constant.Type, oldEnv, newEnv),
		}
	}
	newModule.constants = newConstants

	newSubtypes := make(ConstantMap, len(m.subtypes))
	for subtypeName, subtype := range m.subtypes {
		newSubtypes[subtypeName] = Constant{
			FullName: subtype.FullName,
			Type:     DeepCopy(subtype.Type, oldEnv, newEnv),
		}
	}
	newModule.subtypes = newSubtypes

	newModule.parent = DeepCopy(m.parent, oldEnv, newEnv).(Namespace)
	return newModule
}
