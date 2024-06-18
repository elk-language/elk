package types

import (
	"fmt"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/value"
)

func MakeFullConstantName(containerName, constName string) string {
	if containerName == "Root" || containerName == "" {
		return constName
	}
	return fmt.Sprintf("%s::%s", containerName, constName)
}

type MethodMap = concurrent.Map[value.Symbol, *Method]

func NewMethodMap() *MethodMap {
	return concurrent.NewMap[value.Symbol, *Method]()
}

type TypeMap = concurrent.Map[value.Symbol, Type]

func NewTypeMap() *TypeMap {
	return concurrent.NewMap[value.Symbol, Type]()
}

type ConstantMap struct {
	name              string
	constants         *TypeMap
	subtypes          *TypeMap
	instanceVariables *TypeMap
	methods           *MethodMap
}

func MakeConstantMap(name string) ConstantMap {
	return ConstantMap{
		name:              name,
		constants:         NewTypeMap(),
		subtypes:          NewTypeMap(),
		instanceVariables: NewTypeMap(),
		methods:           NewMethodMap(),
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

func (c *ConstantMap) InstanceVariables() *TypeMap {
	return c.instanceVariables
}

func (c *ConstantMap) Subtypes() *TypeMap {
	return c.subtypes
}

// Get the constant with the given name.
func (c *ConstantMap) Constant(name value.Symbol) Type {
	result, _ := c.constants.Get(name)
	return result
}

// Get the constant with the given name.
func (c *ConstantMap) ConstantString(name string) Type {
	result, _ := c.constants.Get(value.ToSymbol(name))
	return result
}

// Get the subtype with the given name.
func (c *ConstantMap) Subtype(name value.Symbol) Type {
	result, _ := c.subtypes.Get(name)
	return result
}

// Get the subtype with the given name.
func (c *ConstantMap) SubtypeString(name string) Type {
	result, _ := c.subtypes.Get(value.ToSymbol(name))
	return result
}

// Get the method with the given name.
func (c *ConstantMap) Method(name value.Symbol) *Method {
	result, _ := c.methods.Get(name)
	return result
}

// Get the method with the given name.
func (c *ConstantMap) MethodString(name string) *Method {
	result, _ := c.methods.Get(value.ToSymbol(name))
	return result
}

func (c *ConstantMap) DefineInstanceVariable(name string, val Type) {
	c.instanceVariables.Set(value.ToSymbol(name), val)
}

// Get the instance variable with the given name.
func (c *ConstantMap) InstanceVariable(name value.Symbol) Type {
	result, _ := c.instanceVariables.Get(name)
	return result
}

// Get the instance variable with the given name.
func (c *ConstantMap) InstanceVariableString(name string) Type {
	result, _ := c.instanceVariables.Get(value.ToSymbol(name))
	return result
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
	c.DefineConstant(name, NewSingletonClass(m))
	return m
}
