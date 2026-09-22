package types

import (
	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value/symbol"
)

type ModulePlaceholder struct {
	Module
}

func (m *ModulePlaceholder) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(m, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(m, parent)
	}
}

func (m *ModulePlaceholder) ToRef() Ref[*ModulePlaceholder] {
	return Ref[*ModulePlaceholder](m.id)
}

func (m *ModulePlaceholder) EqualAny(other any) bool {
	o, ok := other.(*ModulePlaceholder)
	if !ok {
		return false
	}

	if m.id > 0 {
		return m.id == o.ID()
	}

	return m.name == m.name
}

func NewModulePlaceholder(name string) *ModulePlaceholder {
	return &ModulePlaceholder{
		Module: Module{
			NamespaceBase: MakeNamespaceBase("", name),
		},
	}
}

func (m *ModulePlaceholder) Copy() *ModulePlaceholder {
	result := &ModulePlaceholder{
		Module: m.Module,
	}
	result.id = ZERO_ID
	return result
}

func (m *ModulePlaceholder) CopyType() Type {
	return m.Copy()
}

// Used during typechecking as a placeholder for a future
// module, class, mixin, interface etc.
type NamespacePlaceholder struct {
	name      string
	Namespace Ref[Namespace]
	Locations *concurrent.Slice[*position.Location]
	id        ID
}

func NewNamespacePlaceholder(name string) *NamespacePlaceholder {
	return &NamespacePlaceholder{
		name:      name,
		Locations: concurrent.NewSlice[*position.Location](),
		Namespace: CastRef[Namespace](NewModulePlaceholder(name)),
	}
}

func (m *NamespacePlaceholder) EqualAny(other any) bool {
	o, ok := other.(Namespace)
	if !ok {
		return false
	}

	if m.id > 0 {
		return m.id == o.ID()
	}

	return m.name == o.Name()
}

func (p *NamespacePlaceholder) ToNonLiteral() Type {
	return p
}

func (*NamespacePlaceholder) IsLiteral() bool {
	return false
}

func (p *NamespacePlaceholder) inspect() string {
	return p.Name()
}

func (n *NamespacePlaceholder) Copy() *NamespacePlaceholder {
	return &NamespacePlaceholder{
		name:      n.name,
		Namespace: n.Namespace,
		Locations: n.Locations,
	}
}

func (n *NamespacePlaceholder) CopyType() Type {
	return n.Copy()
}

func (n *NamespacePlaceholder) ToRef() Ref[*NamespacePlaceholder] {
	return ToRef(n)
}

func (n *NamespacePlaceholder) ID() ID {
	return n.id
}

func (n *NamespacePlaceholder) SetID(id ID) {
	n.id = id
}

func (n *NamespacePlaceholder) HashUint64() uint64 {
	d := xxhash.New()
	d.WriteString("constant:")
	d.WriteString(n.name)
	return d.Sum64()
}

func (n *NamespacePlaceholder) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(n, parent)
	}
}

func (n *NamespacePlaceholder) Name() string {
	return n.Namespace.Get().Name()
}

func (n *NamespacePlaceholder) DocComment() string {
	return n.Namespace.Get().DocComment()
}

func (n *NamespacePlaceholder) SetDocComment(comment string) {
	n.Namespace.Get().SetDocComment(comment)
}

func (n *NamespacePlaceholder) AppendDocComment(comment string) {
	n.Namespace.Get().AppendDocComment(comment)
}

func (n *NamespacePlaceholder) Parent() Namespace {
	return n.Namespace.Get().Parent()
}

func (n *NamespacePlaceholder) SetParent(parent Namespace) {
	n.Namespace.Get().SetParent(parent)
}

func (n *NamespacePlaceholder) Singleton() *SingletonClass {
	return n.Namespace.Get().Singleton()
}

func (n *NamespacePlaceholder) SetSingleton(singleton *SingletonClass) {
	n.Namespace.Get().SetSingleton(singleton)
}

func (n *NamespacePlaceholder) IsAbstract() bool {
	return n.Namespace.Get().IsAbstract()
}

func (n *NamespacePlaceholder) IsSealed() bool {
	return n.Namespace.Get().IsSealed()
}

func (n *NamespacePlaceholder) IsPrimitive() bool {
	return n.Namespace.Get().IsPrimitive()
}

func (n *NamespacePlaceholder) IsImmutable() bool {
	return n.Namespace.Get().IsImmutable()
}

func (n *NamespacePlaceholder) IsGeneric() bool {
	return n.Namespace.Get().IsGeneric()
}

func (n *NamespacePlaceholder) IsDefined() bool {
	return n.Namespace.Get().IsDefined()
}

func (n *NamespacePlaceholder) SetDefined(defined bool) {
	n.Namespace.Get().SetDefined(defined)
}

func (n *NamespacePlaceholder) IsNative() bool {
	return n.Namespace.Get().IsNative()
}

func (n *NamespacePlaceholder) TypeParameters() []Ref[*TypeParameter] {
	return n.Namespace.Get().TypeParameters()
}

func (n *NamespacePlaceholder) SetTypeParameters(t []Ref[*TypeParameter]) {
	n.Namespace.Get().SetTypeParameters(t)
}

func (n *NamespacePlaceholder) Constants() ConstantMap {
	return n.Namespace.Get().Constants()
}

func (n *NamespacePlaceholder) Constant(name symbol.Symbol) (Constant, bool) {
	return n.Namespace.Get().Constant(name)
}

func (n *NamespacePlaceholder) ConstantString(name string) (Constant, bool) {
	return n.Namespace.Get().ConstantString(name)
}

func (n *NamespacePlaceholder) DefineConstant(name symbol.Symbol, val Type) {
	n.Namespace.Get().DefineConstant(name, val)
}

func (n *NamespacePlaceholder) DefineConstantWithFullName(name symbol.Symbol, fullName string, val Type) {
	n.Namespace.Get().DefineConstantWithFullName(name, fullName, val)
}

func (n *NamespacePlaceholder) Subtypes() ConstantMap {
	return n.Namespace.Get().Subtypes()
}

func (n *NamespacePlaceholder) Subtype(name symbol.Symbol) (Constant, bool) {
	return n.Namespace.Get().Subtype(name)
}

func (n *NamespacePlaceholder) SubtypeString(name string) (Constant, bool) {
	return n.Namespace.Get().SubtypeString(name)
}

func (n *NamespacePlaceholder) MustSubtype(name symbol.Symbol) Type {
	return n.Namespace.Get().MustSubtype(name)
}

func (n *NamespacePlaceholder) MustSubtypeString(name string) Type {
	return n.Namespace.Get().MustSubtypeString(name)
}

func (n *NamespacePlaceholder) DefineSubtype(name symbol.Symbol, val Type) {
	n.Namespace.Get().DefineSubtype(name, val)
}

func (n *NamespacePlaceholder) DefineSubtypeWithFullName(name symbol.Symbol, fullName string, val Type) {
	n.Namespace.Get().DefineSubtypeWithFullName(name, fullName, val)
}

func (n *NamespacePlaceholder) Methods() MethodMap {
	return n.Namespace.Get().Methods()
}

func (n *NamespacePlaceholder) Method(name symbol.Symbol) *Method {
	return n.Namespace.Get().Method(name)
}

func (n *NamespacePlaceholder) MethodString(name string) *Method {
	return n.Namespace.Get().MethodString(name)
}

func (n *NamespacePlaceholder) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType, throwType Type) *Method {
	return n.Namespace.Get().DefineMethod(docComment, flags, name, typeParams, params, returnType, throwType)
}

func (n *NamespacePlaceholder) SetMethod(name symbol.Symbol, method *Method) {
	n.Namespace.Get().SetMethod(name, method)
}

func (n *NamespacePlaceholder) InstanceVariables() InstanceVariableMap {
	return n.Namespace.Get().InstanceVariables()
}

func (n *NamespacePlaceholder) InstanceVariable(name symbol.Symbol) *InstanceVariable {
	return n.Namespace.Get().InstanceVariable(name)
}

func (n *NamespacePlaceholder) InstanceVariableString(name string) *InstanceVariable {
	return n.Namespace.Get().InstanceVariableString(name)
}

func (n *NamespacePlaceholder) DefineInstanceVariable(name symbol.Symbol, ivar *InstanceVariable) {
	n.Namespace.Get().DefineInstanceVariable(name, ivar)
}

func (n *NamespacePlaceholder) DefineClass(docComment string, primitive, abstract, sealed, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	return n.Namespace.Get().DefineClass(docComment, primitive, abstract, sealed, noinit, immutable, name, parent)
}

func (n *NamespacePlaceholder) DefineModule(docComment string, name symbol.Symbol) *Module {
	return n.Namespace.Get().DefineModule(docComment, name)
}

func (n *NamespacePlaceholder) DefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	return n.Namespace.Get().DefineMixin(docComment, abstract, name)
}

func (n *NamespacePlaceholder) DefineInterface(docComment string, name symbol.Symbol) *Interface {
	return n.Namespace.Get().DefineInterface(docComment, name)
}
