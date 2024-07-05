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

type NamespaceBase struct {
	docComment        string
	name              string
	constants         *TypeMap
	subtypes          *TypeMap
	instanceVariables *TypeMap
	methods           *MethodMap
}

func MakeNamespaceBase(docComment, name string) NamespaceBase {
	return NamespaceBase{
		docComment:        docComment,
		name:              name,
		constants:         NewTypeMap(),
		subtypes:          NewTypeMap(),
		instanceVariables: NewTypeMap(),
		methods:           NewMethodMap(),
	}
}

func (c *NamespaceBase) DocComment() string {
	return c.docComment
}

func (c *NamespaceBase) SetDocComment(comment string) {
	c.docComment = comment
}

func (c *NamespaceBase) AppendDocComment(newComment string) {
	if newComment == "" {
		return
	}
	if c.docComment == "" {
		c.docComment = newComment
		return
	}

	c.docComment = fmt.Sprintf("%s\n\n%s", c.docComment, newComment)
}

func (c *NamespaceBase) Name() string {
	return c.name
}

func (c *NamespaceBase) Methods() *MethodMap {
	return c.methods
}

func (c *NamespaceBase) Constants() *TypeMap {
	return c.constants
}

func (c *NamespaceBase) InstanceVariables() *TypeMap {
	return c.instanceVariables
}

func (c *NamespaceBase) Subtypes() *TypeMap {
	return c.subtypes
}

// Get the constant with the given name.
func (c *NamespaceBase) Constant(name value.Symbol) Type {
	result, _ := c.constants.Get(name)
	return result
}

// Get the constant with the given name.
func (c *NamespaceBase) ConstantString(name string) Type {
	result, _ := c.constants.Get(value.ToSymbol(name))
	return result
}

// Get the subtype with the given name.
func (c *NamespaceBase) Subtype(name value.Symbol) Type {
	result, _ := c.subtypes.Get(name)
	return result
}

// Get the subtype with the given name.
func (c *NamespaceBase) SubtypeString(name string) Type {
	result, _ := c.subtypes.Get(value.ToSymbol(name))
	return result
}

// Get the method with the given name.
func (c *NamespaceBase) Method(name value.Symbol) *Method {
	result, _ := c.methods.Get(name)
	return result
}

// Get the method with the given name.
func (c *NamespaceBase) MethodString(name string) *Method {
	result, _ := c.methods.Get(value.ToSymbol(name))
	return result
}

func (c *NamespaceBase) DefineInstanceVariable(name string, val Type) {
	c.instanceVariables.Set(value.ToSymbol(name), val)
}

// Get the instance variable with the given name.
func (c *NamespaceBase) InstanceVariable(name value.Symbol) Type {
	result, _ := c.instanceVariables.Get(name)
	return result
}

// Get the instance variable with the given name.
func (c *NamespaceBase) InstanceVariableString(name string) Type {
	result, _ := c.instanceVariables.Get(value.ToSymbol(name))
	return result
}

func (c *NamespaceBase) DefineConstant(name string, val Type) {
	c.constants.Set(value.ToSymbol(name), val)
}

func (c *NamespaceBase) DefineSubtype(name string, val Type) {
	c.subtypes.Set(value.ToSymbol(name), val)
}

func (c *NamespaceBase) SetMethod(name string, method *Method) {
	c.methods.Set(value.ToSymbol(name), method)
}

// Define a new class.
func (c *NamespaceBase) DefineClass(docComment, name string, parent Namespace, env *GlobalEnvironment) *Class {
	class := NewClass(docComment, MakeFullConstantName(c.Name(), name), parent, env)
	c.DefineSubtype(name, class)
	c.DefineConstant(name, class.singleton)
	return class
}

// Define a new module.
func (c *NamespaceBase) DefineModule(docComment, name string) *Module {
	m := NewModule(docComment, MakeFullConstantName(c.Name(), name))
	c.DefineSubtype(name, m)
	c.DefineConstant(name, m)
	return m
}

// Define a new mixin.
func (c *NamespaceBase) DefineMixin(docComment string, name string, env *GlobalEnvironment) *Mixin {
	m := NewMixin(docComment, MakeFullConstantName(c.Name(), name), env)
	c.DefineSubtype(name, m)
	c.DefineConstant(name, m.singleton)
	return m
}

// Define a new interface.
func (c *NamespaceBase) DefineInterface(docComment string, name string, env *GlobalEnvironment) *Interface {
	m := NewInterface(docComment, MakeFullConstantName(c.Name(), name), env)
	c.DefineSubtype(name, m)
	c.DefineConstant(name, m.singleton)
	return m
}
