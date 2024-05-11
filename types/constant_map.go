package types

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

type ConstantMap struct {
	name      string
	constants map[value.Symbol]Type
	subtypes  map[value.Symbol]Type
}

func (c *ConstantMap) Name() string {
	return c.name
}

func (c *ConstantMap) Constants() map[value.Symbol]Type {
	return c.constants
}

func (c *ConstantMap) Subtypes() map[value.Symbol]Type {
	return c.subtypes
}

// Get the constant with the given name.
func (c *ConstantMap) Constant(name string) Type {
	return c.constants[value.ToSymbol(name)]
}

// Get the subtype with the given name.
func (c *ConstantMap) Subtype(name string) Type {
	return c.subtypes[value.ToSymbol(name)]
}

// Get the constant with the given name.
func (c *ConstantMap) DefineConstant(name string, val Type) {
	if c.constants == nil {
		c.constants = make(map[value.Symbol]Type)
	}
	c.constants[value.ToSymbol(name)] = val
}

// Get the constant with the given name.
func (c *ConstantMap) DefineSubtype(name string, val Type) {
	if c.subtypes == nil {
		c.subtypes = make(map[value.Symbol]Type)
	}
	c.subtypes[value.ToSymbol(name)] = val
}

// Define a new class.
func (c *ConstantMap) DefineClass(name string, parent *Class, consts map[value.Symbol]Type) *Class {
	class := NewClass(fmt.Sprintf("%s::%s", c.Name(), name), parent, consts)
	c.DefineSubtype(name, class)
	c.DefineConstant(name, NewSingletonClass(class))
	return class
}

// Define a new module.
func (c *ConstantMap) DefineModule(name string, consts map[value.Symbol]Type) *Module {
	m := NewModule(fmt.Sprintf("%s::%s", c.Name(), name), consts)
	c.DefineSubtype(name, m)
	c.DefineConstant(name, m)
	return m
}
