package types

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

func MakeFullConstantName(containerName, constName string) string {
	if containerName == "Root" || containerName == "" {
		return constName
	}
	return fmt.Sprintf("%s::%s", containerName, constName)
}

type MethodMap map[value.Symbol]*Method

type ConstantMap struct {
	name      string
	constants map[value.Symbol]Type
	subtypes  map[value.Symbol]Type
	methods   MethodMap
}

func (c *ConstantMap) Name() string {
	return c.name
}

func (c *ConstantMap) Methods() MethodMap {
	return c.methods
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

// Get the method with the given name.
func (c *ConstantMap) Method(name string) *Method {
	return c.methods[value.ToSymbol(name)]
}

func (c *ConstantMap) DefineConstant(name string, val Type) {
	if c.constants == nil {
		c.constants = make(map[value.Symbol]Type)
	}
	c.constants[value.ToSymbol(name)] = val
}

func (c *ConstantMap) DefineSubtype(name string, val Type) {
	if c.subtypes == nil {
		c.subtypes = make(map[value.Symbol]Type)
	}
	c.subtypes[value.ToSymbol(name)] = val
}

func (c *ConstantMap) DefineMethod(name string, params ParameterMap, returnType, throwType Type) *Method {
	method := NewMethod(name, params, returnType, throwType)
	c.SetMethod(name, method)
	return method
}

func (c *ConstantMap) SetMethod(name string, method *Method) {
	if c.methods == nil {
		c.methods = make(MethodMap)
	}
	c.methods[value.ToSymbol(name)] = method
}

// Define a new class.
func (c *ConstantMap) DefineClass(name string, parent *Class, consts map[value.Symbol]Type) *Class {
	class := NewClass(MakeFullConstantName(c.Name(), name), parent, consts)
	c.DefineSubtype(name, class)
	c.DefineConstant(name, NewSingletonClass(class))
	return class
}

// Define a new module.
func (c *ConstantMap) DefineModule(name string, consts map[value.Symbol]Type, subtypes map[value.Symbol]Type) *Module {
	m := NewModule(MakeFullConstantName(c.Name(), name), consts, subtypes)
	c.DefineSubtype(name, m)
	c.DefineConstant(name, m)
	return m
}
