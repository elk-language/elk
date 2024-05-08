package types

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

type Module struct {
	ConstantMap
}

func NewModule(name string, consts map[value.Symbol]Type) *Module {
	return &Module{
		ConstantMap: ConstantMap{
			Name: name,
			Map:  consts,
		},
	}
}

func (m *Module) Inspect() string {
	return fmt.Sprintf("module %s", m.Name)
}

func (m *Module) IsSupertypeOf(other Type) bool {
	otherMod, ok := other.(*Module)
	if !ok {
		return false
	}

	return m == otherMod
}
