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
	if containerName[0] == '&' {
		return fmt.Sprintf("(%s)::%s", containerName, constName)
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

func (c *NamespaceBase) SetName(name string) {
	c.name = name
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

func (c *NamespaceBase) DefineInstanceVariable(name value.Symbol, val Type) {
	c.instanceVariables.Set(name, val)
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

func (c *NamespaceBase) DefineConstant(name value.Symbol, val Type) {
	c.constants.Set(name, val)
}

func (c *NamespaceBase) DefineSubtype(name value.Symbol, val Type) {
	c.subtypes.Set(name, val)
}

func (c *NamespaceBase) SetMethod(name value.Symbol, method *Method) {
	c.methods.Set(name, method)
}

// Define a new class if it does not exist
func (c *NamespaceBase) TryDefineClass(docComment string, abstract, sealed, primitive bool, name value.Symbol, parent Namespace, env *GlobalEnvironment) *Class {
	subtype := c.Subtype(name)
	if subtype == nil {
		return c.DefineClass(docComment, primitive, abstract, sealed, true, name, parent, env)
	}

	class := subtype.(*Class)
	class.AppendDocComment(docComment)
	if class.IsPrimitive() != primitive || class.IsAbstract() != abstract || class.IsSealed() != sealed {
		panic(
			fmt.Sprintf(
				"%s modifier mismatch, previous: %s, now: %s",
				InspectWithColor(class),
				InspectModifier(class.IsAbstract(), class.IsSealed(), class.IsPrimitive()),
				InspectModifier(abstract, sealed, primitive),
			),
		)
	}
	return class
}

// Define a new class.
func (c *NamespaceBase) DefineClass(docComment string, abstract, sealed, primitive, fullyChecked bool, name value.Symbol, parent Namespace, env *GlobalEnvironment) *Class {
	class := NewClass(docComment, abstract, sealed, primitive, MakeFullConstantName(c.Name(), name.String()), parent, env)
	class.FullyChecked = fullyChecked
	c.DefineSubtype(name, class)
	c.DefineConstant(name, class.singleton)
	return class
}

// Define a new module if it does not exist.
func (c *NamespaceBase) TryDefineModule(docComment string, name value.Symbol) *Module {
	subtype := c.Subtype(name)
	if subtype == nil {
		return c.DefineModule(docComment, name)
	}

	module := subtype.(*Module)
	module.AppendDocComment(docComment)
	return module
}

// Define a new module.
func (c *NamespaceBase) DefineModule(docComment string, name value.Symbol) *Module {
	m := NewModule(docComment, MakeFullConstantName(c.Name(), name.String()))
	c.DefineSubtype(name, m)
	c.DefineConstant(name, m)
	return m
}

// Define a new mixin if it does not exist.
func (c *NamespaceBase) TryDefineMixin(docComment string, abstract bool, name value.Symbol, env *GlobalEnvironment) *Mixin {
	subtype := c.Subtype(name)
	if subtype == nil {
		return c.DefineMixin(docComment, abstract, name, env)
	}

	mixin := subtype.(*Mixin)
	mixin.AppendDocComment(docComment)
	if mixin.IsAbstract() != abstract {
		panic(
			fmt.Sprintf(
				"%s modifier mismatch, previous: %s, now: %s",
				InspectWithColor(mixin),
				InspectModifier(mixin.IsAbstract(), false, false),
				InspectModifier(abstract, false, false),
			),
		)
	}
	return mixin
}

// Define a new mixin.
func (c *NamespaceBase) DefineMixin(docComment string, abstract bool, name value.Symbol, env *GlobalEnvironment) *Mixin {
	m := NewMixin(docComment, abstract, MakeFullConstantName(c.Name(), name.String()), env)
	c.DefineSubtype(name, m)
	c.DefineConstant(name, m.singleton)
	return m
}

// Define a new module if it does not exist.
func (c *NamespaceBase) TryDefineInterface(docComment string, name value.Symbol, env *GlobalEnvironment) *Interface {
	subtype := c.Subtype(name)
	if subtype == nil {
		return c.DefineInterface(docComment, name, env)
	}

	iface := subtype.(*Interface)
	iface.AppendDocComment(docComment)
	return iface
}

// Define a new interface.
func (c *NamespaceBase) DefineInterface(docComment string, name value.Symbol, env *GlobalEnvironment) *Interface {
	m := NewInterface(docComment, MakeFullConstantName(c.Name(), name.String()), false, env)
	c.DefineSubtype(name, m)
	c.DefineConstant(name, m.singleton)
	return m
}
