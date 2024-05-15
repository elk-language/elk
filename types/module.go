package types

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

type Module struct {
	ConstantMap
}

func NewModule(name string, consts map[value.Symbol]Type, subtypes map[value.Symbol]Type) *Module {
	return &Module{
		ConstantMap: ConstantMap{
			name:      name,
			constants: consts,
			subtypes:  subtypes,
		},
	}
}

func (m *Module) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}

func (m *Module) Inspect() string {
	return fmt.Sprintf("module %s", m.Name())
}

func (m *Module) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	otherMod, ok := other.(*Module)
	if !ok {
		return false
	}

	return m == otherMod
}
