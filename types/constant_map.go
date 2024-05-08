package types

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

type ConstantMap struct {
	Name string
	Map  map[value.Symbol]Type
}

func (c *ConstantMap) Constants() map[value.Symbol]Type {
	return c.Map
}

// Get the constant with the given name.
func (c *ConstantMap) Constant(name string) Type {
	return c.Map[value.ToSymbol(name)]
}

// Get the constant with the given name.
func (c *ConstantMap) DefineConstant(name string, val Type) {
	if c.Map == nil {
		c.Map = make(map[value.Symbol]Type)
	}
	c.Map[value.ToSymbol(name)] = val
}

// Define a new class.
func (c *ConstantMap) DefineClass(name string, parent *Class, consts map[value.Symbol]Type) *Class {
	class := NewClass(fmt.Sprintf("%s::%s", c.Name, name), parent, consts)
	c.DefineConstant(name, class)
	return class
}

// Define a new module.
func (c *ConstantMap) DefineModule(name string, consts map[value.Symbol]Type) *Module {
	m := NewModule(fmt.Sprintf("%s::%s", c.Name, name), consts)
	c.DefineConstant(name, m)
	return m
}
