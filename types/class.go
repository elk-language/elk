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
			Name:      name,
			constants: consts,
		},
	}
}

func (c *Class) Inspect() string {
	return c.Name
}

func (c *Class) IsSupertypeOf(other Type) bool {
	otherClass, ok := other.(*Class)
	if !ok {
		return false
	}

	currentOther := otherClass
	for {
		if currentOther == nil {
			return false
		}
		if currentOther == c {
			return true
		}

		currentOther = currentOther.Parent
	}
}
