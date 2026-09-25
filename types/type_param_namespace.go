package types

import (
	"maps"
	"unsafe"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value/symbol"
)

type TypeParamNamespace struct {
	docComment string
	ForMethod  bool
	constants  ConstantMap
	subtypes   ConstantMap
	id         ID
}

var _ Namespace = &TypeParamNamespace{}

func (n *TypeParamNamespace) ToRef() Ref[*TypeParamNamespace] {
	return ToRef(n)
}

func (n *TypeParamNamespace) HashUint64() uint64 {
	return uint64(uintptr(unsafe.Pointer(n)))
}

func (n *TypeParamNamespace) EqualAny(other any) bool {
	return n == other
}

func (n *TypeParamNamespace) ID() ID {
	return n.id
}

func (n *TypeParamNamespace) SetID(id ID) {
	n.id = id
}

func (t *TypeParamNamespace) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(t, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(t, parent)
	}
}

func NewTypeParamNamespace(docComment string, forMethod bool) *TypeParamNamespace {
	return &TypeParamNamespace{
		docComment: docComment,
		ForMethod:  forMethod,
		constants:  make(ConstantMap),
		subtypes:   make(ConstantMap),
	}
}

func (t *TypeParamNamespace) Copy() *TypeParamNamespace {
	return &TypeParamNamespace{
		docComment: t.docComment,
		ForMethod:  t.ForMethod,
		constants:  maps.Clone(t.constants),
		subtypes:   maps.Clone(t.subtypes),
	}
}

func (t *TypeParamNamespace) CopyType() Type {
	return t.Copy()
}

func (t *TypeParamNamespace) Name() string {
	return ""
}

func (t *TypeParamNamespace) inspect() string {
	return t.docComment
}

func (t *TypeParamNamespace) DocComment() string {
	return t.docComment
}

func (t *TypeParamNamespace) SetDocComment(string) {
	panic("cannot set doc comment on type param namespace")
}

func (t *TypeParamNamespace) AppendDocComment(string) {
	panic("cannot append doc comment on type param namespace")
}

func (t *TypeParamNamespace) Parent() Namespace {
	return nil
}

func (t *TypeParamNamespace) ParentRef() Ref[Namespace] {
	return 0
}

func (t *TypeParamNamespace) SetParent(Namespace) {
	panic("cannot set parent of type param namespaces")
}

func (t *TypeParamNamespace) Singleton() *SingletonClass {
	return nil
}

func (t *TypeParamNamespace) SingletonRef() Ref[*SingletonClass] {
	return 0
}

func (t *TypeParamNamespace) SetSingleton(*SingletonClass) {
	panic("cannot set singleton class of closure")
}

func (t *TypeParamNamespace) IsDefined() bool {
	return false
}

func (t *TypeParamNamespace) SetDefined(bool) {
	panic("cannot set `defined` in type param namespace")
}

func (t *TypeParamNamespace) IsAbstract() bool {
	return true
}

func (t *TypeParamNamespace) IsSealed() bool {
	return true
}

func (t *TypeParamNamespace) IsPrimitive() bool {
	return true
}

func (t *TypeParamNamespace) IsImmutable() bool {
	return true
}

func (t *TypeParamNamespace) IsNative() bool {
	return false
}

func (t *TypeParamNamespace) IsGeneric() bool {
	return false
}

func (t *TypeParamNamespace) TypeParameters() []Ref[*TypeParameter] {
	return nil
}

func (t *TypeParamNamespace) SetTypeParameters([]Ref[*TypeParameter]) {
	panic("cannot set type parameters on a type parameter namespace")
}

func (t *TypeParamNamespace) Constants() ConstantMap {
	return t.constants
}

func (t *TypeParamNamespace) Constant(name symbol.Symbol) (Constant, bool) {
	result, ok := t.constants[name]
	return result, ok
}

func (t *TypeParamNamespace) ConstantString(name string) (Constant, bool) {
	return t.Constant(symbol.ToSymbol(name))
}

func (t *TypeParamNamespace) DefineConstant(name symbol.Symbol, val Type) {
	t.constants[name] = Constant{
		Type: ToRef(val),
	}
}

func (t *TypeParamNamespace) DefineConstantWithFullName(name symbol.Symbol, fullName string, val Type) {
	t.constants[name] = Constant{
		Type:     ToRef(val),
		FullName: fullName,
	}
}

func (t *TypeParamNamespace) Subtypes() ConstantMap {
	return nil
}

func (t *TypeParamNamespace) Subtype(name symbol.Symbol) (Constant, bool) {
	result, ok := t.subtypes[name]
	return result, ok
}

func (t *TypeParamNamespace) MustSubtypeString(name string) Type {
	return t.subtypes[symbol.ToSymbol(name)].Type.Get()
}

func (t *TypeParamNamespace) MustSubtype(name symbol.Symbol) Type {
	return t.subtypes[name].Type.Get()
}

func (t *TypeParamNamespace) SubtypeString(name string) (Constant, bool) {
	return t.Subtype(symbol.ToSymbol(name))
}

func (t *TypeParamNamespace) DefineSubtype(name symbol.Symbol, val Type) {
	t.subtypes[name] = Constant{
		Type: ToRef(val),
	}
}

func (t *TypeParamNamespace) DefineSubtypeWithFullName(name symbol.Symbol, fullName string, val Type) {
	t.subtypes[name] = Constant{
		Type: ToRef(val),
	}
}

func (t *TypeParamNamespace) Methods() MethodMap {
	return nil
}

func (t *TypeParamNamespace) Method(name symbol.Symbol) *Method {
	return nil
}

func (t *TypeParamNamespace) MethodString(name string) *Method {
	return nil
}

func (t *TypeParamNamespace) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType, throwType Type) *Method {
	panic("cannot define methods on type param namespaces")
}

func (t *TypeParamNamespace) SetMethod(name symbol.Symbol, method *Method) {
}

func (t *TypeParamNamespace) InstanceVariables() InstanceVariableMap {
	return nil
}

func (t *TypeParamNamespace) InstanceVariable(name symbol.Symbol) *InstanceVariable {
	return nil
}

func (t *TypeParamNamespace) InstanceVariableString(name string) *InstanceVariable {
	return nil
}

func (t *TypeParamNamespace) DefineInstanceVariable(name symbol.Symbol, ivar *InstanceVariable) {
	panic("cannot define instance variables on type param namespaces")
}

func (t *TypeParamNamespace) DefineClass(docComment string, primitive, abstract, sealed, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	panic("cannot define classes on type param namespaces")
}

func (t *TypeParamNamespace) DefineModule(docComment string, name symbol.Symbol) *Module {
	panic("cannot define module on type param namespaces")
}

func (t *TypeParamNamespace) DefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	panic("cannot define mixins on type param namespaces")
}

func (t *TypeParamNamespace) DefineInterface(docComment string, name symbol.Symbol) *Interface {
	panic("cannot define interfaces on type param namespaces")
}

func (t *TypeParamNamespace) ToNonLiteral() Type {
	return t
}

func (*TypeParamNamespace) IsLiteral() bool {
	return false
}
