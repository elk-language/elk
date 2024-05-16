package types

import (
	"github.com/elk-language/elk/value"
)

type Class struct {
	Parent *Class
	ConstantMap
}

func NewClass(name string, parent *Class, consts map[value.Symbol]Type) *Class {
	return &Class{
		ConstantMap: ConstantMap{
			name:      name,
			constants: consts,
		},
	}
}

func (c *Class) inspect() string {
	return c.name
}

func (c *Class) ToNonLiteral(env *GlobalEnvironment) Type {
	return c
}
