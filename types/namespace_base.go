package types

import (
	"fmt"

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

type MethodMap = map[value.Symbol]*Method

type TypeMap = map[value.Symbol]Type

type ConstantMap = map[value.Symbol]Constant

type Constant struct {
	FullName string
	Type     Type
}

type NamespaceBase struct {
	docComment        string
	name              string
	constants         ConstantMap
	subtypes          ConstantMap
	instanceVariables TypeMap
	methods           MethodMap
	methodAliases     MethodAliasMap
}

func MakeNamespaceBase(docComment, name string) NamespaceBase {
	return NamespaceBase{
		docComment:        docComment,
		name:              name,
		constants:         make(ConstantMap),
		subtypes:          make(ConstantMap),
		instanceVariables: make(TypeMap),
		methods:           make(MethodMap),
		methodAliases:     make(MethodAliasMap),
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

func (c *NamespaceBase) Methods() MethodMap {
	return c.methods
}

func (c *NamespaceBase) MethodAliases() MethodAliasMap {
	return c.methodAliases
}

func (c *NamespaceBase) SetMethodAlias(name value.Symbol, method *Method) {
	c.methodAliases[name] = NewMethodAlias(method)
	c.SetMethod(name, method)
}

func (c *NamespaceBase) Constants() ConstantMap {
	return c.constants
}

func (c *NamespaceBase) InstanceVariables() TypeMap {
	return c.instanceVariables
}

func (c *NamespaceBase) HasInstanceVariables() bool {
	return len(c.instanceVariables) > 0
}

func (c *NamespaceBase) SetInstanceVariables(iv TypeMap) {
	c.instanceVariables = iv
}

func (c *NamespaceBase) Subtypes() ConstantMap {
	return c.subtypes
}

// Get the constant with the given name.
func (c *NamespaceBase) Constant(name value.Symbol) (Constant, bool) {
	result, ok := c.constants[name]
	return result, ok
}

// Get the constant with the given name.
func (c *NamespaceBase) ConstantString(name string) (Constant, bool) {
	return c.Constant(value.ToSymbol(name))
}

// Get the subtype with the given name.
func (c *NamespaceBase) Subtype(name value.Symbol) (Constant, bool) {
	result, ok := c.subtypes[name]
	return result, ok
}

// Get the subtype with the given name.
func (c *NamespaceBase) SubtypeString(name string) (Constant, bool) {
	return c.Subtype(value.ToSymbol(name))
}

func (c *NamespaceBase) MustSubtype(name value.Symbol) Type {
	return c.subtypes[name].Type
}

func (c *NamespaceBase) MustSubtypeString(name string) Type {
	return c.subtypes[value.ToSymbol(name)].Type
}

// Get the method with the given name.
func (c *NamespaceBase) Method(name value.Symbol) *Method {
	return c.methods[name]
}

// Get the method with the given name.
func (c *NamespaceBase) MethodString(name string) *Method {
	return c.methods[value.ToSymbol(name)]
}

func (c *NamespaceBase) DefineInstanceVariable(name value.Symbol, val Type) {
	c.instanceVariables[name] = val
}

// Get the instance variable with the given name.
func (c *NamespaceBase) InstanceVariable(name value.Symbol) Type {
	return c.instanceVariables[name]
}

// Get the instance variable with the given name.
func (c *NamespaceBase) InstanceVariableString(name string) Type {
	return c.instanceVariables[value.ToSymbol(name)]
}

func (c *NamespaceBase) DefineConstant(name value.Symbol, val Type) {
	c.DefineConstantWithFullName(name, MakeFullConstantName(c.Name(), name.String()), val)
}

func (c *NamespaceBase) DefineConstantWithFullName(name value.Symbol, fullName string, val Type) {
	c.constants[name] = Constant{
		FullName: fullName,
		Type:     val,
	}
}

func (c *NamespaceBase) DefineSubtype(name value.Symbol, val Type) {
	c.DefineSubtypeWithFullName(name, MakeFullConstantName(c.Name(), name.String()), val)
}

func (c *NamespaceBase) DefineSubtypeWithFullName(name value.Symbol, fullName string, val Type) {
	c.subtypes[name] = Constant{
		FullName: fullName,
		Type:     val,
	}
}

func (c *NamespaceBase) SetMethod(name value.Symbol, method *Method) {
	c.methods[name] = method
}

// Define a new class if it does not exist
func (c *NamespaceBase) TryDefineClass(docComment string, abstract, sealed, primitive, noinit bool, name value.Symbol, parent Namespace, env *GlobalEnvironment) *Class {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineClass(docComment, abstract, sealed, primitive, noinit, name, parent, env)
	}

	class := subtype.Type.(*Class)
	class.AppendDocComment(docComment)

	if class.IsPrimitive() != primitive || class.IsAbstract() != abstract || class.IsSealed() != sealed {
		panic(
			fmt.Sprintf(
				"%s modifier mismatch, previous: %s, now: %s",
				InspectWithColor(class),
				InspectModifier(class.IsAbstract(), class.IsSealed(), class.IsPrimitive(), class.IsNoInit()),
				InspectModifier(abstract, sealed, primitive, noinit),
			),
		)
	}
	return class
}

// Define a new class.
func (c *NamespaceBase) DefineClass(docComment string, abstract, sealed, primitive, noinit bool, name value.Symbol, parent Namespace, env *GlobalEnvironment) *Class {
	fullName := MakeFullConstantName(c.Name(), name.String())
	class := NewClass(docComment, abstract, sealed, primitive, noinit, fullName, parent, env)
	c.DefineSubtypeWithFullName(name, fullName, class)
	c.DefineConstantWithFullName(name, fullName, class.singleton)
	return class
}

// Define a new module if it does not exist.
func (c *NamespaceBase) TryDefineModule(docComment string, name value.Symbol, env *GlobalEnvironment) *Module {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineModule(docComment, name, env)
	}

	module := subtype.Type.(*Module)
	module.AppendDocComment(docComment)
	return module
}

// Define a new module.
func (c *NamespaceBase) DefineModule(docComment string, name value.Symbol, env *GlobalEnvironment) *Module {
	fullName := MakeFullConstantName(c.Name(), name.String())
	m := NewModule(docComment, fullName, env)
	c.DefineSubtypeWithFullName(name, fullName, m)
	c.DefineConstantWithFullName(name, fullName, m)
	return m
}

// Define a new mixin if it does not exist.
func (c *NamespaceBase) TryDefineMixin(docComment string, abstract bool, name value.Symbol, env *GlobalEnvironment) *Mixin {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineMixin(docComment, abstract, name, env)
	}

	mixin := subtype.Type.(*Mixin)
	mixin.AppendDocComment(docComment)
	if mixin.IsAbstract() != abstract {
		panic(
			fmt.Sprintf(
				"%s modifier mismatch, previous: %s, now: %s",
				InspectWithColor(mixin),
				InspectModifier(mixin.IsAbstract(), false, false, false),
				InspectModifier(abstract, false, false, false),
			),
		)
	}
	return mixin
}

// Define a new mixin.
func (c *NamespaceBase) DefineMixin(docComment string, abstract bool, name value.Symbol, env *GlobalEnvironment) *Mixin {
	fullName := MakeFullConstantName(c.Name(), name.String())
	m := NewMixin(docComment, abstract, fullName, env)
	c.DefineSubtypeWithFullName(name, fullName, m)
	c.DefineConstantWithFullName(name, fullName, m.singleton)
	return m
}

// Define a new module if it does not exist.
func (c *NamespaceBase) TryDefineInterface(docComment string, name value.Symbol, env *GlobalEnvironment) *Interface {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineInterface(docComment, name, env)
	}

	iface := subtype.Type.(*Interface)
	iface.AppendDocComment(docComment)
	return iface
}

// Define a new interface.
func (c *NamespaceBase) DefineInterface(docComment string, name value.Symbol, env *GlobalEnvironment) *Interface {
	fullName := MakeFullConstantName(c.Name(), name.String())
	m := NewInterface(docComment, fullName, env)
	c.DefineSubtypeWithFullName(name, fullName, m)
	c.DefineConstantWithFullName(name, fullName, m.singleton)
	return m
}
