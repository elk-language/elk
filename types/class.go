package types

import (
	"github.com/elk-language/elk/value"
)

type Class struct {
	ConstantMap
}

func NewClass(name string, parent ConstantContainer, consts map[value.Symbol]Type) *Class {
	return &Class{
		ConstantMap: ConstantMap{
			name:      name,
			constants: consts,
			parent:    parent,
		},
	}
}

func (c *Class) DefineMethod(name string, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(name, params, returnType, throwType, c)
	c.SetMethod(name, method)
	return method
}

func (c *Class) inspect() string {
	return c.name
}

func (c *Class) ToNonLiteral(env *GlobalEnvironment) Type {
	return c
}
