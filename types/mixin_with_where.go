package types

import (
	"encoding/binary"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value/symbol"
)

// Represents a mixin parent with a where clause
// `extend where Bar < Baz`
type MixinWithWhere struct {
	MixinProxy Ref[*MixinProxy]
	Namespace  Ref[Namespace]
	Where      []Ref[*TypeParameter]
	id         ID
}

func IsMixinWithWhere(typ Type) bool {
	_, ok := typ.(*MixinWithWhere)
	return ok
}

func NewMixinWithWhere(mixin *MixinProxy, namespace Namespace, where []Ref[*TypeParameter]) *MixinWithWhere {
	return &MixinWithWhere{
		MixinProxy: mixin.ToRef(),
		Namespace:  ToRef(namespace),
		Where:      where,
	}
}

func (m *MixinWithWhere) ToNonLiteral() Type {
	return m
}

func (m *MixinWithWhere) inspect() string {
	return m.Namespace.Get().inspect()
}

func (m *MixinWithWhere) InspectExtend() string {
	buffer := new(strings.Builder)
	buffer.WriteString("extend where ")
	firstIteration := true
	for _, whereElement := range m.Where {
		if !firstIteration {
			buffer.WriteString(", ")
		} else {
			firstIteration = false
		}

		buffer.WriteString(whereElement.Get().InspectSignature())
	}

	return buffer.String()
}

func (m *MixinWithWhere) InspectExtendWithColor() string {
	return lexer.Colorize(m.InspectExtend())
}

func (m *MixinWithWhere) Copy() *MixinWithWhere {
	return &MixinWithWhere{
		MixinProxy: m.MixinProxy,
		Namespace:  m.Namespace,
		Where:      m.Where,
		id:         m.id,
	}
}

func (m *MixinWithWhere) ToRef() Ref[*MixinWithWhere] {
	return ToRef(m)
}

func (m *MixinWithWhere) ID() ID {
	return m.id
}

func (m *MixinWithWhere) SetID(id ID) {
	m.id = id
}

func (m *MixinWithWhere) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("mixinwhere:")
	b := make([]byte, 4)

	binary.LittleEndian.PutUint32(b, uint32(m.MixinProxy))
	d.Write(b)

	binary.LittleEndian.PutUint32(b, uint32(m.Namespace))
	d.Write(b)

	for _, whereElement := range m.Where {
		binary.LittleEndian.PutUint32(b, uint32(whereElement))
		d.Write(b)
	}

	return d.Sum64()
}

func (m *MixinWithWhere) EqualAny(other any) bool {
	o, ok := other.(*MixinWithWhere)
	if !ok {
		return false
	}

	if m.id > 0 {
		return m.id == o.ID()
	}

	if m.MixinProxy != o.MixinProxy || m.Namespace != o.Namespace || len(m.Where) != len(o.Where) {
		return false
	}

	for i := range m.Where {
		if m.Where[i] != o.Where[i] {
			return false
		}
	}

	return true
}

func (m *MixinWithWhere) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(m, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(m, parent)
	}
}

func (m *MixinWithWhere) Parent() Namespace {
	return m.MixinProxy.Get().Parent()
}

func (m *MixinWithWhere) SetParent(parent Namespace) {
	m.MixinProxy.Get().SetParent(parent)
}

func (m *MixinWithWhere) IsLiteral() bool {
	return m.MixinProxy.Get().IsLiteral()
}

func (m *MixinWithWhere) IsGeneric() bool {
	return m.MixinProxy.Get().IsGeneric()
}

func (m *MixinWithWhere) TypeParameters() []Ref[*TypeParameter] {
	return m.MixinProxy.Get().TypeParameters()
}

func (m *MixinWithWhere) SetTypeParameters(t []Ref[*TypeParameter]) {
	m.MixinProxy.Get().SetTypeParameters(t)
}

func (m *MixinWithWhere) Singleton() *SingletonClass {
	return m.MixinProxy.Get().Singleton()
}

func (m *MixinWithWhere) SetSingleton(singleton *SingletonClass) {
	m.MixinProxy.Get().SetSingleton(singleton)
}

func (m *MixinWithWhere) SetAbstract(abstract bool) *Mixin {
	return m.MixinProxy.Get().SetAbstract(abstract)
}

func (m *MixinWithWhere) IsAbstract() bool {
	return m.MixinProxy.Get().IsAbstract()
}

func (m *MixinWithWhere) IsSealed() bool {
	return m.MixinProxy.Get().IsSealed()
}

func (m *MixinWithWhere) IsNative() bool {
	return m.MixinProxy.Get().IsNative()
}

func (m *MixinWithWhere) SetNative(native bool) {
	m.MixinProxy.Get().SetNative(native)
}

func (m *MixinWithWhere) IsDefined() bool {
	return m.MixinProxy.Get().IsDefined()
}

func (m *MixinWithWhere) SetDefined(compiled bool) {
	m.MixinProxy.Get().SetDefined(compiled)
}

func (m *MixinWithWhere) IsPrimitive() bool {
	return m.MixinProxy.Get().IsPrimitive()
}

func (m *MixinWithWhere) IsImmutable() bool {
	return m.MixinProxy.Get().IsImmutable()
}

func (m *MixinWithWhere) RemoveTemporaryParents() {
	m.MixinProxy.Get().RemoveTemporaryParents()
}

func (m *MixinWithWhere) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType, throwType Type) *Method {
	return m.MixinProxy.Get().DefineMethod(docComment, flags, name, typeParams, params, returnType, throwType)
}

func (m *MixinWithWhere) Name() string {
	return m.MixinProxy.Get().Name()
}

func (m *MixinWithWhere) SetName(name string) {
	m.MixinProxy.Get().SetName(name)
}

func (m *MixinWithWhere) DocComment() string {
	return m.MixinProxy.Get().DocComment()
}

func (m *MixinWithWhere) SetDocComment(comment string) {
	m.MixinProxy.Get().SetDocComment(comment)
}

func (m *MixinWithWhere) AppendDocComment(comment string) {
	m.MixinProxy.Get().AppendDocComment(comment)
}

func (m *MixinWithWhere) Constants() ConstantMap {
	return m.MixinProxy.Get().Constants()
}

func (m *MixinWithWhere) Constant(name symbol.Symbol) (Constant, bool) {
	return m.MixinProxy.Get().Constant(name)
}

func (m *MixinWithWhere) ConstantString(name string) (Constant, bool) {
	return m.MixinProxy.Get().ConstantString(name)
}

func (m *MixinWithWhere) DefineConstant(name symbol.Symbol, val Type) {
	m.MixinProxy.Get().DefineConstant(name, val)
}

func (m *MixinWithWhere) DefineConstantWithFullName(name symbol.Symbol, fullName string, val Type) {
	m.MixinProxy.Get().DefineConstantWithFullName(name, fullName, val)
}

func (m *MixinWithWhere) Subtypes() ConstantMap {
	return m.MixinProxy.Get().Subtypes()
}

func (m *MixinWithWhere) Subtype(name symbol.Symbol) (Constant, bool) {
	return m.MixinProxy.Get().Subtype(name)
}

func (m *MixinWithWhere) SubtypeString(name string) (Constant, bool) {
	return m.MixinProxy.Get().SubtypeString(name)
}

func (m *MixinWithWhere) MustSubtype(name symbol.Symbol) Type {
	return m.MixinProxy.Get().MustSubtype(name)
}

func (m *MixinWithWhere) MustSubtypeString(name string) Type {
	return m.MixinProxy.Get().MustSubtypeString(name)
}

func (m *MixinWithWhere) DefineSubtype(name symbol.Symbol, val Type) {
	m.MixinProxy.Get().DefineSubtype(name, val)
}

func (m *MixinWithWhere) DefineSubtypeWithFullName(name symbol.Symbol, fullName string, val Type) {
	m.MixinProxy.Get().DefineSubtypeWithFullName(name, fullName, val)
}

func (m *MixinWithWhere) Methods() MethodMap {
	return m.MixinProxy.Get().Methods()
}

func (m *MixinWithWhere) Method(name symbol.Symbol) *Method {
	return m.MixinProxy.Get().Method(name)
}

func (m *MixinWithWhere) MethodString(name string) *Method {
	return m.MixinProxy.Get().MethodString(name)
}

func (m *MixinWithWhere) SetMethod(name symbol.Symbol, method *Method) {
	m.MixinProxy.Get().SetMethod(name, method)
}

func (m *MixinWithWhere) InstanceVariables() InstanceVariableMap {
	return m.MixinProxy.Get().InstanceVariables()
}

func (m *MixinWithWhere) HasInstanceVariables() bool {
	return m.MixinProxy.Get().HasInstanceVariables()
}

func (m *MixinWithWhere) SetInstanceVariables(iv InstanceVariableMap) {
	m.MixinProxy.Get().SetInstanceVariables(iv)
}

func (m *MixinWithWhere) InstanceVariable(name symbol.Symbol) *InstanceVariable {
	return m.MixinProxy.Get().InstanceVariable(name)
}

func (m *MixinWithWhere) InstanceVariableString(name string) *InstanceVariable {
	return m.MixinProxy.Get().InstanceVariableString(name)
}

func (m *MixinWithWhere) DefineInstanceVariable(name symbol.Symbol, ivar *InstanceVariable) {
	m.MixinProxy.Get().DefineInstanceVariable(name, ivar)
}

func (m *MixinWithWhere) DefineClass(docComment string, abstract, sealed, primitive, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	return m.MixinProxy.Get().DefineClass(docComment, abstract, sealed, primitive, noinit, immutable, name, parent)
}

func (m *MixinWithWhere) TryDefineClass(docComment string, abstract, sealed, primitive, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	return m.MixinProxy.Get().TryDefineClass(docComment, abstract, sealed, primitive, noinit, immutable, name, parent)
}

func (m *MixinWithWhere) DefineModule(docComment string, name symbol.Symbol) *Module {
	return m.MixinProxy.Get().DefineModule(docComment, name)
}

func (m *MixinWithWhere) TryDefineModule(docComment string, name symbol.Symbol) *Module {
	return m.MixinProxy.Get().TryDefineModule(docComment, name)
}

func (m *MixinWithWhere) DefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	return m.MixinProxy.Get().DefineMixin(docComment, abstract, name)
}

func (m *MixinWithWhere) TryDefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	return m.MixinProxy.Get().TryDefineMixin(docComment, abstract, name)
}

func (m *MixinWithWhere) DefineInterface(docComment string, name symbol.Symbol) *Interface {
	return m.MixinProxy.Get().DefineInterface(docComment, name)
}

func (m *MixinWithWhere) TryDefineInterface(docComment string, name symbol.Symbol) *Interface {
	return m.MixinProxy.Get().TryDefineInterface(docComment, name)
}
