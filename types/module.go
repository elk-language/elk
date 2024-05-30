package types

import (
	"github.com/elk-language/elk/value"
)

type Module struct {
	ConstantMap
}

func (m *Module) Parent() ConstantContainer {
	return nil
}

func NewModule(
	name string,
	consts map[value.Symbol]Type,
	subtypes map[value.Symbol]Type,
	methods MethodMap,
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
