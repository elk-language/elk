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

func (c *Class) Inspect() string {
	return c.name
}

func (c *Class) ToNonLiteral(env *GlobalEnvironment) Type {
	return c
}

func (c *Class) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	otherClass, ok := other.(*Class)
	if !ok {
		return false
	}

	currentOther := c
	for {
		if currentOther == nil {
			return false
		}
		if currentOther == otherClass {
			return true
		}

		currentOther = currentOther.Parent
	}
}
