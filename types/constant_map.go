package types

import (
	"fmt"

	"github.com/elk-language/elk/threadsafe"
	"github.com/elk-language/elk/value"
)

func MakeFullConstantName(containerName, constName string) string {
	if containerName == "Root" || containerName == "" {
		return constName
	}
	return fmt.Sprintf("%s::%s", containerName, constName)
}

type MethodMap = threadsafe.Map[value.Symbol, *Method]

func NewMethodMap() *MethodMap {
	return threadsafe.NewMap[value.Symbol, *Method]()
}

type TypeMap = threadsafe.Map[value.Symbol, Type]

func NewTypeMap() *TypeMap {
	return threadsafe.NewMap[value.Symbol, Type]()
}

type ConstantMap struct {
	name      string
	constants *TypeMap
	subtypes  *TypeMap
	methods   *MethodMap
}

func MakeConstantMap(name string) ConstantMap {
	return ConstantMap{
		name:      name,
		constants: NewTypeMap(),
		subtypes:  NewTypeMap(),
		methods:   NewMethodMap(),
	}
}

func (c *ConstantMap) Name() string {
	return c.name
}

func (c *ConstantMap) Methods() *MethodMap {
	return c.methods
}

func (c *ConstantMap) Constants() *TypeMap {
	return c.constants
}

func (c *ConstantMap) Subtypes() *TypeMap {
	return c.subtypes
}

// Get the constant with the given name.
func (c *ConstantMap) Constant(name value.Symbol) Type {
	return c.constants.Get(name)
}

// Get the constant with the given name.
func (c *ConstantMap) ConstantString(name string) Type {
	return c.constants.Get(value.ToSymbol(name))
}

// Get the subtype with the given name.
func (c *ConstantMap) Subtype(name value.Symbol) Type {
	return c.subtypes.Get(name)
}

// Get the subtype with the given name.
func (c *ConstantMap) SubtypeString(name string) Type {
	return c.subtypes.Get(value.ToSymbol(name))
}

// Get the method with the given name.
func (c *ConstantMap) Method(name value.Symbol) *Method {
	return c.methods.Get(name)
}

// Get the method with the given name.
func (c *ConstantMap) MethodString(name string) *Method {
	return c.methods.Get(value.ToSymbol(name))
}

func (c *ConstantMap) DefineConstant(name string, val Type) {
	c.constants.Set(value.ToSymbol(name), val)
}

func (c *ConstantMap) DefineSubtype(name string, val Type) {
	c.subtypes.Set(value.ToSymbol(name), val)
}

func (c *ConstantMap) SetMethod(name string, method *Method) {
	c.methods.Set(value.ToSymbol(name), method)
}

// Define a new class.
func (c *ConstantMap) DefineClass(name string, parent ConstantContainer) *Class {
	class := NewClass(MakeFullConstantName(c.Name(), name), parent)
	c.DefineSubtype(name, class)
	c.DefineConstant(name, NewSingletonClass(class))
	return class
}

// Define a new module.
func (c *ConstantMap) DefineModule(name string) *Module {
	m := NewModule(MakeFullConstantName(c.Name(), name))
	c.DefineSubtype(name, m)
	c.DefineConstant(name, m)
	return m
}

// Define a new mixin.
func (c *ConstantMap) DefineMixin(name string) *Mixin {
	m := NewMixin(MakeFullConstantName(c.Name(), name))
	c.DefineSubtype(name, m)
	c.DefineConstant(name, m)
	return m
}
