package types

import "fmt"

type Module struct {
	NamespaceBase
}

func (*Module) Singleton() *SingletonClass {
	return nil
}

func (*Module) Parent() Namespace {
	return nil
}

func (m *Module) SetParent(Namespace) {
	panic(fmt.Sprintf("cannot set the parent of module `%s`", m.Name()))
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

func NewModule(docComment, name string) *Module {
	return &Module{
		NamespaceBase: MakeNamespaceBase(docComment, name),
	}
}

func NewModuleWithDetails(
	docComment string,
	name string,
	consts *TypeMap,
	subtypes *TypeMap,
	methods *MethodMap,
) *Module {
	return &Module{
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

func (m *Module) inspect() string {
	return m.Name()
}

func (m *Module) DefineMethod(docComment, name string, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, name, params, returnType, throwType, m)
	m.SetMethod(name, method)
	return method
}
