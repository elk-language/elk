package types

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/value/symbol"
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

type MethodMap = map[symbol.Symbol]Ref[*Method]

type TypeMap = map[symbol.Symbol]Ref[Type]

type InstanceVariableMap = map[symbol.Symbol]*InstanceVariable

type ConstantMap = map[symbol.Symbol]Constant

type Constant struct {
	FullName string
	Type     Ref[Type]
}

type NamespaceBase struct {
	docComment        string
	name              string
	constants         ConstantMap
	subtypes          ConstantMap
	instanceVariables InstanceVariableMap
	methods           MethodMap
	id                ID
}

func MakeNamespaceBase(docComment, name string) NamespaceBase {
	return NamespaceBase{
		docComment:        docComment,
		name:              name,
		constants:         make(ConstantMap),
		subtypes:          make(ConstantMap),
		instanceVariables: make(InstanceVariableMap),
		methods:           make(MethodMap),
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

func (c *NamespaceBase) Constants() ConstantMap {
	return c.constants
}

func (c *NamespaceBase) InstanceVariables() InstanceVariableMap {
	return c.instanceVariables
}

func (c *NamespaceBase) HasInstanceVariables() bool {
	return len(c.instanceVariables) > 0
}

func (c *NamespaceBase) SetInstanceVariables(iv InstanceVariableMap) {
	c.instanceVariables = iv
}

func (c *NamespaceBase) Subtypes() ConstantMap {
	return c.subtypes
}

// Get the constant with the given name.
func (c *NamespaceBase) Constant(name symbol.Symbol) (Constant, bool) {
	result, ok := c.constants[name]
	return result, ok
}

// Get the constant with the given name.
func (c *NamespaceBase) ConstantString(name string) (Constant, bool) {
	return c.Constant(symbol.ToSymbol(name))
}

// Get the subtype with the given name.
func (c *NamespaceBase) Subtype(name symbol.Symbol) (Constant, bool) {
	result, ok := c.subtypes[name]
	return result, ok
}

// Get the subtype with the given name.
func (c *NamespaceBase) SubtypeString(name string) (Constant, bool) {
	return c.Subtype(symbol.ToSymbol(name))
}

func (c *NamespaceBase) MustSubtype(name symbol.Symbol) Type {
	return c.subtypes[name].Type.Get()
}

func (c *NamespaceBase) MustSubtypeString(name string) Type {
	return c.subtypes[symbol.ToSymbol(name)].Type.Get()
}

// Get the method with the given name.
func (c *NamespaceBase) Method(name symbol.Symbol) *Method {
	return c.methods[name].Get()
}

// Get the method with the given name.
func (c *NamespaceBase) MethodString(name string) *Method {
	return c.methods[symbol.ToSymbol(name)].Get()
}

func (c *NamespaceBase) DefineInstanceVariable(name symbol.Symbol, ivar *InstanceVariable) {
	c.instanceVariables[name] = ivar
}

// Get the instance variable with the given name.
func (c *NamespaceBase) InstanceVariable(name symbol.Symbol) *InstanceVariable {
	return c.instanceVariables[name]
}

// Get the instance variable with the given name.
func (c *NamespaceBase) InstanceVariableString(name string) *InstanceVariable {
	return c.instanceVariables[symbol.ToSymbol(name)]
}

func (c *NamespaceBase) DefineConstant(name symbol.Symbol, val Type) {
	c.DefineConstantWithFullName(name, MakeFullConstantName(c.Name(), name.String()), val)
}

func (c *NamespaceBase) DefineConstantRef(name symbol.Symbol, val Ref[Type]) {
	c.DefineConstantWithFullNameRef(name, MakeFullConstantName(c.Name(), name.String()), val)
}

func (c *NamespaceBase) DefineConstantWithFullName(name symbol.Symbol, fullName string, val Type) {
	c.constants[name] = Constant{
		FullName: fullName,
		Type:     ToRef(val),
	}
}

func (c *NamespaceBase) DefineConstantWithFullNameRef(name symbol.Symbol, fullName string, val Ref[Type]) {
	c.constants[name] = Constant{
		FullName: fullName,
		Type:     val,
	}
}

func (c *NamespaceBase) DefineSubtype(name symbol.Symbol, val Type) {
	c.DefineSubtypeWithFullName(name, MakeFullConstantName(c.Name(), name.String()), val)
}

func (c *NamespaceBase) DefineSubtypeWithFullName(name symbol.Symbol, fullName string, val Type) {
	c.subtypes[name] = Constant{
		FullName: fullName,
		Type:     ToRef(val),
	}
}

func (c *NamespaceBase) SetMethod(name symbol.Symbol, method *Method) {
	c.methods[name] = method.ToRef()
}

// Define a new class if it does not exist
func (c *NamespaceBase) TryDefineClass(docComment string, abstract, sealed, primitive, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineClass(docComment, abstract, sealed, primitive, noinit, immutable, name, parent)
	}

	class := subtype.Type.Get().(*Class)
	class.AppendDocComment(docComment)

	if class.IsPrimitive() != primitive || class.IsAbstract() != abstract || class.IsSealed() != sealed {
		panic(
			fmt.Sprintf(
				"%s modifier mismatch, previous: %s, now: %s",
				InspectWithColor(class),
				InspectModifier(ModifierSet{Abstract: class.IsAbstract(), Sealed: class.IsSealed(), Primitive: class.IsPrimitive(), NoInit: class.IsNoInit(), Immutable: class.IsImmutable()}),
				InspectModifier(ModifierSet{Abstract: abstract, Sealed: sealed, Primitive: primitive, NoInit: noinit, Immutable: immutable}),
			),
		)
	}
	return class
}

// Define a new class.
func (c *NamespaceBase) DefineClass(docComment string, abstract, sealed, primitive, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	fullName := MakeFullConstantName(c.Name(), name.String())
	class := NewClass(docComment, abstract, sealed, primitive, noinit, immutable, fullName, parent)
	c.DefineSubtypeWithFullName(name, fullName, class)
	c.DefineConstantWithFullName(name, fullName, class.singleton.Get())
	return class
}

// Define a new module if it does not exist.
func (c *NamespaceBase) TryDefineModule(docComment string, name symbol.Symbol) *Module {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineModule(docComment, name)
	}

	module := subtype.Type.Get().(*Module)
	module.AppendDocComment(docComment)
	return module
}

// Define a new module.
func (c *NamespaceBase) DefineModule(docComment string, name symbol.Symbol) *Module {
	fullName := MakeFullConstantName(c.Name(), name.String())
	m := NewModule(docComment, fullName)
	c.DefineSubtypeWithFullName(name, fullName, m)
	c.DefineConstantWithFullName(name, fullName, m)
	return m
}

// Define a new mixin if it does not exist.
func (c *NamespaceBase) TryDefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineMixin(docComment, abstract, name)
	}

	mixin := subtype.Type.Get().(*Mixin)
	mixin.AppendDocComment(docComment)
	if mixin.IsAbstract() != abstract {
		panic(
			fmt.Sprintf(
				"%s modifier mismatch, previous: %s, now: %s",
				InspectWithColor(mixin),
				InspectModifier(ModifierSet{Abstract: mixin.IsAbstract()}),
				InspectModifier(ModifierSet{Abstract: abstract}),
			),
		)
	}
	return mixin
}

// Define a new mixin.
func (c *NamespaceBase) DefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	fullName := MakeFullConstantName(c.Name(), name.String())
	m := NewMixin(docComment, abstract, fullName)
	c.DefineSubtypeWithFullName(name, fullName, m)
	c.DefineConstantWithFullNameRef(name, fullName, Ref[Type](m.singleton))
	return m
}

// Define a new module if it does not exist.
func (c *NamespaceBase) TryDefineInterface(docComment string, name symbol.Symbol) *Interface {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineInterface(docComment, name)
	}

	iface := subtype.Type.Get().(*Interface)
	iface.AppendDocComment(docComment)
	return iface
}

// Define a new interface.
func (c *NamespaceBase) DefineInterface(docComment string, name symbol.Symbol) *Interface {
	fullName := MakeFullConstantName(c.Name(), name.String())
	m := NewInterface(docComment, fullName)
	c.DefineSubtypeWithFullName(name, fullName, m)
	c.DefineConstantWithFullName(name, fullName, m.singleton.Get())
	return m
}

func (c *NamespaceBase) HashUint64() uint64 {
	d := xxhash.New()
	d.WriteString("constant:")
	d.WriteString(c.name)
	return d.Sum64()
}

func (c *NamespaceBase) ID() ID {
	return c.id
}

func (c *NamespaceBase) SetID(id ID) {
	c.id = id
}
