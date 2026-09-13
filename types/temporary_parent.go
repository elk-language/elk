package types

import (
	"encoding/binary"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value/symbol"
)

// A temporary wrapper around a namespace parent
// used in the macro expansion phase before all namespaces
// are known.
type TemporaryParent struct {
	Namespace Ref[Namespace]
	id        ID
}

func NewTemporaryParent(namespace Namespace) *TemporaryParent {
	return &TemporaryParent{
		Namespace: ToRef(namespace),
	}
}

func (n *TemporaryParent) ToNonLiteral() Type {
	return n
}

func (*TemporaryParent) IsLiteral() bool {
	return false
}

func (n *TemporaryParent) inspect() string {
	return n.Name()
}

func (n *TemporaryParent) ToRef() Ref[*TemporaryParent] {
	return ToRef(n)
}

func (n *TemporaryParent) ID() ID {
	return n.id
}

func (n *TemporaryParent) SetID(id ID) {
	n.id = id
}

func (n *TemporaryParent) HashUint64() uint64 {
	d := xxhash.New()
	d.WriteString("tempparent:")
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(n.Namespace))
	d.Write(b)
	return d.Sum64()
}

func (n *TemporaryParent) EqualAny(other any) bool {
	o, ok := other.(*TemporaryParent)
	if !ok {
		return false
	}

	if n.id > 0 {
		return n.id == o.ID()
	}

	return n.Namespace == o.Namespace
}

func (n *TemporaryParent) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(n, parent)
	}
}

func (n *TemporaryParent) Name() string {
	return n.Namespace.Get().Name()
}

func (n *TemporaryParent) DocComment() string {
	return n.Namespace.Get().DocComment()
}

func (n *TemporaryParent) SetDocComment(comment string) {
	n.Namespace.Get().SetDocComment(comment)
}

func (n *TemporaryParent) AppendDocComment(comment string) {
	n.Namespace.Get().AppendDocComment(comment)
}

func (n *TemporaryParent) Parent() Namespace {
	return n.Namespace.Get().Parent()
}

func (n *TemporaryParent) SetParent(parent Namespace) {
	n.Namespace.Get().SetParent(parent)
}

func (n *TemporaryParent) Singleton() *SingletonClass {
	return n.Namespace.Get().Singleton()
}

func (n *TemporaryParent) SetSingleton(singleton *SingletonClass) {
	n.Namespace.Get().SetSingleton(singleton)
}

func (n *TemporaryParent) IsAbstract() bool {
	return n.Namespace.Get().IsAbstract()
}

func (n *TemporaryParent) IsSealed() bool {
	return n.Namespace.Get().IsSealed()
}

func (n *TemporaryParent) IsPrimitive() bool {
	return n.Namespace.Get().IsPrimitive()
}

func (n *TemporaryParent) IsImmutable() bool {
	return n.Namespace.Get().IsImmutable()
}

func (n *TemporaryParent) IsGeneric() bool {
	return n.Namespace.Get().IsGeneric()
}

func (n *TemporaryParent) IsDefined() bool {
	return n.Namespace.Get().IsDefined()
}

func (n *TemporaryParent) SetDefined(defined bool) {
	n.Namespace.Get().SetDefined(defined)
}

func (n *TemporaryParent) IsNative() bool {
	return n.Namespace.Get().IsNative()
}

func (n *TemporaryParent) TypeParameters() []Ref[*TypeParameter] {
	return n.Namespace.Get().TypeParameters()
}

func (n *TemporaryParent) SetTypeParameters(t []Ref[*TypeParameter]) {
	n.Namespace.Get().SetTypeParameters(t)
}

func (n *TemporaryParent) Constants() ConstantMap {
	return n.Namespace.Get().Constants()
}

func (n *TemporaryParent) Constant(name symbol.Symbol) (Constant, bool) {
	return n.Namespace.Get().Constant(name)
}

func (n *TemporaryParent) ConstantString(name string) (Constant, bool) {
	return n.Namespace.Get().ConstantString(name)
}

func (n *TemporaryParent) DefineConstant(name symbol.Symbol, val Type) {
	n.Namespace.Get().DefineConstant(name, val)
}

func (n *TemporaryParent) DefineConstantWithFullName(name symbol.Symbol, fullName string, val Type) {
	n.Namespace.Get().DefineConstantWithFullName(name, fullName, val)
}

func (n *TemporaryParent) Subtypes() ConstantMap {
	return n.Namespace.Get().Subtypes()
}

func (n *TemporaryParent) Subtype(name symbol.Symbol) (Constant, bool) {
	return n.Namespace.Get().Subtype(name)
}

func (n *TemporaryParent) SubtypeString(name string) (Constant, bool) {
	return n.Namespace.Get().SubtypeString(name)
}

func (n *TemporaryParent) MustSubtype(name symbol.Symbol) Type {
	return n.Namespace.Get().MustSubtype(name)
}

func (n *TemporaryParent) MustSubtypeString(name string) Type {
	return n.Namespace.Get().MustSubtypeString(name)
}

func (n *TemporaryParent) DefineSubtype(name symbol.Symbol, val Type) {
	n.Namespace.Get().DefineSubtype(name, val)
}

func (n *TemporaryParent) DefineSubtypeWithFullName(name symbol.Symbol, fullName string, val Type) {
	n.Namespace.Get().DefineSubtypeWithFullName(name, fullName, val)
}

func (n *TemporaryParent) Methods() MethodMap {
	return n.Namespace.Get().Methods()
}

func (n *TemporaryParent) Method(name symbol.Symbol) *Method {
	return n.Namespace.Get().Method(name)
}

func (n *TemporaryParent) MethodString(name string) *Method {
	return n.Namespace.Get().MethodString(name)
}

func (n *TemporaryParent) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType, throwType Type) *Method {
	return n.Namespace.Get().DefineMethod(docComment, flags, name, typeParams, params, returnType, throwType)
}

func (n *TemporaryParent) SetMethod(name symbol.Symbol, method *Method) {
	n.Namespace.Get().SetMethod(name, method)
}

func (n *TemporaryParent) InstanceVariables() InstanceVariableMap {
	return n.Namespace.Get().InstanceVariables()
}

func (n *TemporaryParent) InstanceVariable(name symbol.Symbol) *InstanceVariable {
	return n.Namespace.Get().InstanceVariable(name)
}

func (n *TemporaryParent) InstanceVariableString(name string) *InstanceVariable {
	return n.Namespace.Get().InstanceVariableString(name)
}

func (n *TemporaryParent) DefineInstanceVariable(name symbol.Symbol, ivar *InstanceVariable) {
	n.Namespace.Get().DefineInstanceVariable(name, ivar)
}

func (n *TemporaryParent) DefineClass(docComment string, primitive, abstract, sealed, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	return n.Namespace.Get().DefineClass(docComment, primitive, abstract, sealed, noinit, immutable, name, parent)
}

func (n *TemporaryParent) DefineModule(docComment string, name symbol.Symbol) *Module {
	return n.Namespace.Get().DefineModule(docComment, name)
}

func (n *TemporaryParent) DefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	return n.Namespace.Get().DefineMixin(docComment, abstract, name)
}

func (n *TemporaryParent) DefineInterface(docComment string, name symbol.Symbol) *Interface {
	return n.Namespace.Get().DefineInterface(docComment, name)
}
