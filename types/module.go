package types

type Module struct {
	ConstantMap
}

func (m *Module) Parent() ConstantContainer {
	return nil
}

func (m *Module) IsAbstract() bool {
	return false
}

func (m *Module) IsSealed() bool {
	return false
}

func NewModule(name string) *Module {
	return &Module{
		ConstantMap: MakeConstantMap(name),
	}
}

func NewModuleWithDetails(
	name string,
	consts *TypeMap,
	subtypes *TypeMap,
	methods *MethodMap,
) *Module {
	return &Module{
		ConstantMap: ConstantMap{
			name:      name,
			constants: consts,
			subtypes:  subtypes,
			methods:   methods,
		},
	}
}

func (m *Module) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}

func (m *Module) inspect() string {
	return m.Name()
}

func (m *Module) DefineMethod(name string, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(name, params, returnType, throwType, m)
	m.SetMethod(name, method)
	return method
}
