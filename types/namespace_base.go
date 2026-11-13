package types

import (
	"fmt"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/position"
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

type MethodMap = map[symbol.Symbol]*Method

type TypeMap = map[symbol.Symbol]Type

type InstanceVariableMap = map[symbol.Symbol]*InstanceVariable

type ConstantMap = map[symbol.Symbol]*Constant

type Constant struct {
	FullName    string
	Type        Type
	docComment  string
	declaration *position.Location
	uses        []*position.Location
}

func (c *Constant) Declaration() *position.Location {
	return c.declaration
}

func (c *Constant) DocComment() string {
	return c.docComment
}

func (c *Constant) Uses() []*position.Location {
	return c.uses
}

func (c *Constant) AppendUse(loc ...*position.Location) {
	c.uses = append(c.uses, loc...)
	AppendUseToType(c.Type, loc...)
}

func NewConstant(docComment string, fullName string, typ Type, declaration *position.Location, uses []*position.Location) *Constant {
	return &Constant{
		docComment:  docComment,
		declaration: declaration,
		FullName:    fullName,
		Type:        typ,
		uses:        uses,
	}
}

type NamespaceBase struct {
	docComment        string
	name              string
	declarations      concurrent.Slice[*position.Location]
	uses              concurrent.Slice[*position.Location]
	constants         ConstantMap
	subtypes          ConstantMap
	instanceVariables InstanceVariableMap
	methods           MethodMap
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

func (c *NamespaceBase) Declaration() *position.Location {
	if c.declarations.Len() == 0 {
		return nil
	}

	l, _ := c.declarations.Get(0)
	return l
}

func (c *NamespaceBase) Declarations() *concurrent.Slice[*position.Location] {
	return &c.declarations
}

func (c *NamespaceBase) AppendDeclaration(decl *position.Location) {
	c.declarations.Push(decl)
}

func (c *NamespaceBase) Uses() *concurrent.Slice[*position.Location] {
	return &c.uses
}

func (c *NamespaceBase) AppendUse(decl ...*position.Location) {
	c.uses.Append(decl...)
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
func (c *NamespaceBase) Constant(name symbol.Symbol) (*Constant, bool) {
	result, ok := c.constants[name]
	return result, ok
}

// Get the constant with the given name.
func (c *NamespaceBase) ConstantString(name string) (*Constant, bool) {
	return c.Constant(symbol.ToSymbol(name))
}

// Get the subtype with the given name.
func (c *NamespaceBase) Subtype(name symbol.Symbol) (*Constant, bool) {
	result, ok := c.subtypes[name]
	return result, ok
}

// Get the subtype with the given name.
func (c *NamespaceBase) SubtypeString(name string) (*Constant, bool) {
	return c.Subtype(symbol.ToSymbol(name))
}

func (c *NamespaceBase) MustSubtype(name symbol.Symbol) Type {
	return c.subtypes[name].Type
}

func (c *NamespaceBase) MustSubtypeString(name string) Type {
	return c.subtypes[symbol.ToSymbol(name)].Type
}

// Get the method with the given name.
func (c *NamespaceBase) Method(name symbol.Symbol) *Method {
	return c.methods[name]
}

// Get the method with the given name.
func (c *NamespaceBase) MethodString(name string) *Method {
	return c.methods[symbol.ToSymbol(name)]
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

func (c *NamespaceBase) DefineConstant(docComment string, name symbol.Symbol, val Type, declaration *position.Location) {
	c.DefineConstantWithFullName(docComment, name, MakeFullConstantName(c.Name(), name.String()), val, declaration)
}

func (c *NamespaceBase) DefineConstantWithFullName(docComment string, name symbol.Symbol, fullName string, val Type, declaration *position.Location) {
	c.constants[name] = &Constant{
		FullName:    fullName,
		Type:        val,
		docComment:  docComment,
		declaration: declaration,
	}
}

func (c *NamespaceBase) DefineSubtype(name symbol.Symbol, val Type) {
	c.DefineSubtypeWithFullName(name, MakeFullConstantName(c.Name(), name.String()), val)
}

func (c *NamespaceBase) DefineSubtypeWithFullName(name symbol.Symbol, fullName string, val Type) {
	c.subtypes[name] = &Constant{
		FullName: fullName,
		Type:     val,
	}
}

func (c *NamespaceBase) SetMethod(name symbol.Symbol, method *Method) {
	c.methods[name] = method
}

// Define a new class if it does not exist
func (c *NamespaceBase) TryDefineClass(docComment string, abstract, sealed, primitive, noinit, immutable bool, name symbol.Symbol, parent Namespace, env *GlobalEnvironment, declaration *position.Location) *Class {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineClass(docComment, abstract, sealed, primitive, noinit, immutable, name, parent, env, declaration)
	}

	class := subtype.Type.(*Class)
	class.AppendDocComment(docComment)
	class.AppendDeclaration(declaration)

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
func (c *NamespaceBase) DefineClass(docComment string, abstract, sealed, primitive, noinit, immutable bool, name symbol.Symbol, parent Namespace, env *GlobalEnvironment, declaration *position.Location) *Class {
	fullName := MakeFullConstantName(c.Name(), name.String())
	class := NewClass(docComment, abstract, sealed, primitive, noinit, immutable, fullName, parent, env)
	c.DefineSubtypeWithFullName(name, fullName, class)
	c.DefineConstantWithFullName(docComment, name, fullName, class.singleton, declaration)
	return class
}

// Define a new module if it does not exist.
func (c *NamespaceBase) TryDefineModule(docComment string, name symbol.Symbol, env *GlobalEnvironment, declaration *position.Location) *Module {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineModule(docComment, name, env, declaration)
	}

	module := subtype.Type.(*Module)
	module.AppendDocComment(docComment)
	module.AppendDeclaration(declaration)
	return module
}

// Define a new module.
func (c *NamespaceBase) DefineModule(docComment string, name symbol.Symbol, env *GlobalEnvironment, declaration *position.Location) *Module {
	fullName := MakeFullConstantName(c.Name(), name.String())
	m := NewModule(docComment, fullName, env)
	c.DefineSubtypeWithFullName(name, fullName, m)
	c.DefineConstantWithFullName(docComment, name, fullName, m, declaration)
	return m
}

// Define a new mixin if it does not exist.
func (c *NamespaceBase) TryDefineMixin(docComment string, abstract bool, name symbol.Symbol, env *GlobalEnvironment, declaration *position.Location) *Mixin {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineMixin(docComment, abstract, name, env, declaration)
	}

	mixin := subtype.Type.(*Mixin)
	mixin.AppendDocComment(docComment)
	mixin.AppendDeclaration(declaration)
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
func (c *NamespaceBase) DefineMixin(docComment string, abstract bool, name symbol.Symbol, env *GlobalEnvironment, declaration *position.Location) *Mixin {
	fullName := MakeFullConstantName(c.Name(), name.String())
	m := NewMixin(docComment, abstract, fullName, env)
	c.DefineSubtypeWithFullName(name, fullName, m)
	c.DefineConstantWithFullName(docComment, name, fullName, m.singleton, declaration)
	return m
}

// Define a new module if it does not exist.
func (c *NamespaceBase) TryDefineInterface(docComment string, name symbol.Symbol, env *GlobalEnvironment, declaration *position.Location) *Interface {
	subtype, ok := c.Subtype(name)
	if !ok {
		return c.DefineInterface(docComment, name, env, declaration)
	}

	iface := subtype.Type.(*Interface)
	iface.AppendDocComment(docComment)
	iface.AppendDeclaration(declaration)
	return iface
}

// Define a new interface.
func (c *NamespaceBase) DefineInterface(docComment string, name symbol.Symbol, env *GlobalEnvironment, declaration *position.Location) *Interface {
	fullName := MakeFullConstantName(c.Name(), name.String())
	m := NewInterface(docComment, fullName, env)
	c.DefineSubtypeWithFullName(name, fullName, m)
	c.DefineConstantWithFullName(docComment, name, fullName, m.singleton, declaration)
	return m
}
