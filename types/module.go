package types

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Module struct {
	parent Namespace
	NamespaceBase
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
		parent:        env.StdSubtypeClass(symbol.Module),
		NamespaceBase: MakeNamespaceBase(docComment, name),
	}
}

func NewModuleWithDetails(
	docComment string,
	name string,
	consts *TypeMap,
	subtypes *TypeMap,
	methods *MethodMap,
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

func (m *Module) DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, abstract, sealed, native, name, params, returnType, throwType, m)
	m.SetMethod(name, method)
	return method
}
