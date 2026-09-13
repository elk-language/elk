package types

import (
	"encoding/binary"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value/symbol"
)

type InterfaceProxy struct {
	parent    Ref[Namespace]
	Interface Ref[*Interface]
	id        ID
}

var _ Namespace = &InterfaceProxy{}

func (i *InterfaceProxy) Parent() Namespace {
	return i.parent.Get()
}

func (i *InterfaceProxy) SetParent(parent Namespace) {
	i.parent = ToRef(parent)
}

func NewInterfaceProxy(iface *Interface, parent Namespace) *InterfaceProxy {
	t := &InterfaceProxy{
		parent:    ToRef(parent),
		Interface: iface.ToRef(),
	}
	return t
}

func (i *InterfaceProxy) ToNonLiteral() Type {
	return i
}

func (i *InterfaceProxy) Copy() *InterfaceProxy {
	return &InterfaceProxy{
		parent:    i.parent,
		Interface: i.Interface,
		id:        i.id,
	}
}

func (i *InterfaceProxy) ToRef() Ref[*InterfaceProxy] {
	return ToRef(i)
}

func (i *InterfaceProxy) ID() ID {
	return i.id
}

func (i *InterfaceProxy) SetID(id ID) {
	i.id = id
}

func (i *InterfaceProxy) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("ifaceproxy:")
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i.Interface))
	d.Write(b)
	if !i.parent.IsZero() {
		binary.LittleEndian.PutUint32(b, uint32(i.parent))
		d.Write(b)
	}

	return d.Sum64()
}

func (i *InterfaceProxy) EqualAny(other any) bool {
	o, ok := other.(*InterfaceProxy)
	if !ok {
		return false
	}

	if i.id > 0 {
		return i.id == o.ID()
	}

	return i.Interface == o.Interface && i.parent == o.parent
}

func (i *InterfaceProxy) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(i, parent)
	}
}

func (i *InterfaceProxy) IsLiteral() bool {
	return i.Interface.Get().IsLiteral()
}

func (i *InterfaceProxy) inspect() string {
	return i.Interface.Get().inspect()
}

func (i *InterfaceProxy) IsGeneric() bool {
	return i.Interface.Get().IsGeneric()
}

func (i *InterfaceProxy) TypeParameters() []Ref[*TypeParameter] {
	return i.Interface.Get().TypeParameters()
}

func (i *InterfaceProxy) SetTypeParameters(t []Ref[*TypeParameter]) {
	i.Interface.Get().SetTypeParameters(t)
}

func (i *InterfaceProxy) Singleton() *SingletonClass {
	return i.Interface.Get().Singleton()
}

func (i *InterfaceProxy) SetSingleton(singleton *SingletonClass) {
	i.Interface.Get().SetSingleton(singleton)
}

func (i *InterfaceProxy) IsDefined() bool {
	return i.Interface.Get().IsDefined()
}

func (i *InterfaceProxy) SetDefined(compiled bool) {
	i.Interface.Get().SetDefined(compiled)
}

func (i *InterfaceProxy) IsAbstract() bool {
	return i.Interface.Get().IsAbstract()
}

func (i *InterfaceProxy) IsNative() bool {
	return i.Interface.Get().IsNative()
}

func (i *InterfaceProxy) IsSealed() bool {
	return i.Interface.Get().IsSealed()
}

func (i *InterfaceProxy) IsPrimitive() bool {
	return i.Interface.Get().IsPrimitive()
}

func (i *InterfaceProxy) IsImmutable() bool {
	return i.Interface.Get().IsImmutable()
}

func (i *InterfaceProxy) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType, throwType Type) *Method {
	return i.Interface.Get().DefineMethod(docComment, flags, name, typeParams, params, returnType, throwType)
}

func (i *InterfaceProxy) Name() string {
	return i.Interface.Get().Name()
}

func (i *InterfaceProxy) SetName(name string) {
	i.Interface.Get().SetName(name)
}

func (i *InterfaceProxy) DocComment() string {
	return i.Interface.Get().DocComment()
}

func (i *InterfaceProxy) SetDocComment(comment string) {
	i.Interface.Get().SetDocComment(comment)
}

func (i *InterfaceProxy) AppendDocComment(comment string) {
	i.Interface.Get().AppendDocComment(comment)
}

func (i *InterfaceProxy) Constants() ConstantMap {
	return i.Interface.Get().Constants()
}

func (i *InterfaceProxy) Constant(name symbol.Symbol) (Constant, bool) {
	return i.Interface.Get().Constant(name)
}

func (i *InterfaceProxy) ConstantString(name string) (Constant, bool) {
	return i.Interface.Get().ConstantString(name)
}

func (i *InterfaceProxy) DefineConstant(name symbol.Symbol, val Type) {
	i.Interface.Get().DefineConstant(name, val)
}

func (i *InterfaceProxy) DefineConstantWithFullName(name symbol.Symbol, fullName string, val Type) {
	i.Interface.Get().DefineConstantWithFullName(name, fullName, val)
}

func (i *InterfaceProxy) Subtypes() ConstantMap {
	return i.Interface.Get().Subtypes()
}

func (i *InterfaceProxy) Subtype(name symbol.Symbol) (Constant, bool) {
	return i.Interface.Get().Subtype(name)
}

func (i *InterfaceProxy) SubtypeString(name string) (Constant, bool) {
	return i.Interface.Get().SubtypeString(name)
}

func (i *InterfaceProxy) MustSubtype(name symbol.Symbol) Type {
	return i.Interface.Get().MustSubtype(name)
}

func (i *InterfaceProxy) MustSubtypeString(name string) Type {
	return i.Interface.Get().MustSubtypeString(name)
}

func (i *InterfaceProxy) DefineSubtype(name symbol.Symbol, val Type) {
	i.Interface.Get().DefineSubtype(name, val)
}

func (i *InterfaceProxy) DefineSubtypeWithFullName(name symbol.Symbol, fullName string, val Type) {
	i.Interface.Get().DefineSubtypeWithFullName(name, fullName, val)
}

func (i *InterfaceProxy) Methods() MethodMap {
	return i.Interface.Get().Methods()
}

func (i *InterfaceProxy) Method(name symbol.Symbol) *Method {
	return i.Interface.Get().Method(name)
}

func (i *InterfaceProxy) MethodString(name string) *Method {
	return i.Interface.Get().MethodString(name)
}

func (i *InterfaceProxy) SetMethod(name symbol.Symbol, method *Method) {
	i.Interface.Get().SetMethod(name, method)
}

func (i *InterfaceProxy) InstanceVariables() InstanceVariableMap {
	return i.Interface.Get().InstanceVariables()
}

func (i *InterfaceProxy) HasInstanceVariables() bool {
	return i.Interface.Get().HasInstanceVariables()
}

func (i *InterfaceProxy) SetInstanceVariables(iv InstanceVariableMap) {
	i.Interface.Get().SetInstanceVariables(iv)
}

func (i *InterfaceProxy) InstanceVariable(name symbol.Symbol) *InstanceVariable {
	return i.Interface.Get().InstanceVariable(name)
}

func (i *InterfaceProxy) InstanceVariableString(name string) *InstanceVariable {
	return i.Interface.Get().InstanceVariableString(name)
}

func (i *InterfaceProxy) DefineInstanceVariable(name symbol.Symbol, ivar *InstanceVariable) {
	i.Interface.Get().DefineInstanceVariable(name, ivar)
}

func (i *InterfaceProxy) DefineClass(docComment string, abstract, sealed, primitive, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	return i.Interface.Get().DefineClass(docComment, abstract, sealed, primitive, noinit, immutable, name, parent)
}

func (i *InterfaceProxy) TryDefineClass(docComment string, abstract, sealed, primitive, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	return i.Interface.Get().TryDefineClass(docComment, abstract, sealed, primitive, noinit, immutable, name, parent)
}

func (i *InterfaceProxy) DefineModule(docComment string, name symbol.Symbol) *Module {
	return i.Interface.Get().DefineModule(docComment, name)
}

func (i *InterfaceProxy) TryDefineModule(docComment string, name symbol.Symbol) *Module {
	return i.Interface.Get().TryDefineModule(docComment, name)
}

func (i *InterfaceProxy) DefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	return i.Interface.Get().DefineMixin(docComment, abstract, name)
}

func (i *InterfaceProxy) TryDefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	return i.Interface.Get().TryDefineMixin(docComment, abstract, name)
}

func (i *InterfaceProxy) DefineInterface(docComment string, name symbol.Symbol) *Interface {
	return i.Interface.Get().DefineInterface(docComment, name)
}

func (i *InterfaceProxy) TryDefineInterface(docComment string, name symbol.Symbol) *Interface {
	return i.Interface.Get().TryDefineInterface(docComment, name)
}
