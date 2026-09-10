package types

import (
	"encoding/binary"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value/symbol"
)

type MixinProxy struct {
	parent Ref[Namespace]
	Mixin  Ref[*Mixin]
	id     ID
}

func (m *MixinProxy) Parent() Namespace {
	return m.parent.Get()
}

func (m *MixinProxy) SetParent(parent Namespace) {
	m.parent = ToRef(parent)
}

func NewMixinProxy(mixin *Mixin, parent Namespace) *MixinProxy {
	m := &MixinProxy{
		parent: ToRef(parent),
		Mixin:  mixin.ToRef(),
	}
	Env.RegisterType(m)
	return m
}

func (m *MixinProxy) ToNonLiteral() Type {
	return m
}

func (m *MixinProxy) Copy() *MixinProxy {
	return &MixinProxy{
		parent: m.parent,
		Mixin:  m.Mixin,
		id:     m.id,
	}
}

func (m *MixinProxy) ToRef() Ref[*MixinProxy] {
	return ToRef(m)
}

func (m *MixinProxy) ID() ID {
	return m.id
}

func (m *MixinProxy) SetID(id ID) {
	m.id = id
}

func (m *MixinProxy) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("mixinproxy:")
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(m.Mixin))
	d.Write(b)
	if !m.parent.IsZero() {
		binary.LittleEndian.PutUint32(b, uint32(m.parent))
		d.Write(b)
	}

	return d.Sum64()
}

func (m *MixinProxy) EqualAny(other any) bool {
	o, ok := other.(*MixinProxy)
	if !ok {
		return false
	}

	if m.id > 0 {
		return m.id == o.ID()
	}

	return m.Mixin == o.Mixin && m.parent == o.parent
}

func (m *MixinProxy) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(m, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(m, parent)
	}
}

func (m *MixinProxy) IsLiteral() bool {
	return m.Mixin.Get().IsLiteral()
}

func (m *MixinProxy) inspect() string {
	return m.Mixin.Get().inspect()
}

func (m *MixinProxy) IsGeneric() bool {
	return m.Mixin.Get().IsGeneric()
}

func (m *MixinProxy) TypeParameters() []*TypeParameter {
	return m.Mixin.Get().TypeParameters()
}

func (m *MixinProxy) SetTypeParameters(t []*TypeParameter) {
	m.Mixin.Get().SetTypeParameters(t)
}

func (m *MixinProxy) Singleton() *SingletonClass {
	return m.Mixin.Get().Singleton()
}

func (m *MixinProxy) SetSingleton(singleton *SingletonClass) {
	m.Mixin.Get().SetSingleton(singleton)
}

func (m *MixinProxy) SetAbstract(abstract bool) *Mixin {
	return m.Mixin.Get().SetAbstract(abstract)
}

func (m *MixinProxy) IsAbstract() bool {
	return m.Mixin.Get().IsAbstract()
}

func (m *MixinProxy) IsSealed() bool {
	return m.Mixin.Get().IsSealed()
}

func (m *MixinProxy) IsNative() bool {
	return m.Mixin.Get().IsNative()
}

func (m *MixinProxy) SetNative(native bool) {
	m.Mixin.Get().SetNative(native)
}

func (m *MixinProxy) IsDefined() bool {
	return m.Mixin.Get().IsDefined()
}

func (m *MixinProxy) SetDefined(compiled bool) {
	m.Mixin.Get().SetDefined(compiled)
}

func (m *MixinProxy) IsPrimitive() bool {
	return m.Mixin.Get().IsPrimitive()
}

func (m *MixinProxy) IsImmutable() bool {
	return m.Mixin.Get().IsImmutable()
}

func (m *MixinProxy) RemoveTemporaryParents(env *GlobalEnvironment) {
	m.Mixin.Get().RemoveTemporaryParents(env)
}

func (m *MixinProxy) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	return m.Mixin.Get().DefineMethod(docComment, flags, name, typeParams, params, returnType, throwType)
}

func (m *MixinProxy) Name() string {
	return m.Mixin.Get().Name()
}

func (m *MixinProxy) SetName(name string) {
	m.Mixin.Get().SetName(name)
}

func (m *MixinProxy) DocComment() string {
	return m.Mixin.Get().DocComment()
}

func (m *MixinProxy) SetDocComment(comment string) {
	m.Mixin.Get().SetDocComment(comment)
}

func (m *MixinProxy) AppendDocComment(comment string) {
	m.Mixin.Get().AppendDocComment(comment)
}

func (m *MixinProxy) Constants() ConstantMap {
	return m.Mixin.Get().Constants()
}

func (m *MixinProxy) Constant(name symbol.Symbol) (Constant, bool) {
	return m.Mixin.Get().Constant(name)
}

func (m *MixinProxy) ConstantString(name string) (Constant, bool) {
	return m.Mixin.Get().ConstantString(name)
}

func (m *MixinProxy) DefineConstant(name symbol.Symbol, val Type) {
	m.Mixin.Get().DefineConstant(name, val)
}

func (m *MixinProxy) DefineConstantWithFullName(name symbol.Symbol, fullName string, val Type) {
	m.Mixin.Get().DefineConstantWithFullName(name, fullName, val)
}

func (m *MixinProxy) Subtypes() ConstantMap {
	return m.Mixin.Get().Subtypes()
}

func (m *MixinProxy) Subtype(name symbol.Symbol) (Constant, bool) {
	return m.Mixin.Get().Subtype(name)
}

func (m *MixinProxy) SubtypeString(name string) (Constant, bool) {
	return m.Mixin.Get().SubtypeString(name)
}

func (m *MixinProxy) MustSubtype(name symbol.Symbol) Type {
	return m.Mixin.Get().MustSubtype(name)
}

func (m *MixinProxy) MustSubtypeString(name string) Type {
	return m.Mixin.Get().MustSubtypeString(name)
}

func (m *MixinProxy) DefineSubtype(name symbol.Symbol, val Type) {
	m.Mixin.Get().DefineSubtype(name, val)
}

func (m *MixinProxy) DefineSubtypeWithFullName(name symbol.Symbol, fullName string, val Type) {
	m.Mixin.Get().DefineSubtypeWithFullName(name, fullName, val)
}

func (m *MixinProxy) Methods() MethodMap {
	return m.Mixin.Get().Methods()
}

func (m *MixinProxy) Method(name symbol.Symbol) *Method {
	return m.Mixin.Get().Method(name)
}

func (m *MixinProxy) MethodString(name string) *Method {
	return m.Mixin.Get().MethodString(name)
}

func (m *MixinProxy) SetMethod(name symbol.Symbol, method *Method) {
	m.Mixin.Get().SetMethod(name, method)
}

func (m *MixinProxy) InstanceVariables() InstanceVariableMap {
	return m.Mixin.Get().InstanceVariables()
}

func (m *MixinProxy) HasInstanceVariables() bool {
	return m.Mixin.Get().HasInstanceVariables()
}

func (m *MixinProxy) SetInstanceVariables(iv InstanceVariableMap) {
	m.Mixin.Get().SetInstanceVariables(iv)
}

func (m *MixinProxy) InstanceVariable(name symbol.Symbol) *InstanceVariable {
	return m.Mixin.Get().InstanceVariable(name)
}

func (m *MixinProxy) InstanceVariableString(name string) *InstanceVariable {
	return m.Mixin.Get().InstanceVariableString(name)
}

func (m *MixinProxy) DefineInstanceVariable(name symbol.Symbol, ivar *InstanceVariable) {
	m.Mixin.Get().DefineInstanceVariable(name, ivar)
}

func (m *MixinProxy) DefineClass(docComment string, abstract, sealed, primitive, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	return m.Mixin.Get().DefineClass(docComment, abstract, sealed, primitive, noinit, immutable, name, parent)
}

func (m *MixinProxy) TryDefineClass(docComment string, abstract, sealed, primitive, noinit, immutable bool, name symbol.Symbol, parent Namespace) *Class {
	return m.Mixin.Get().TryDefineClass(docComment, abstract, sealed, primitive, noinit, immutable, name, parent)
}

func (m *MixinProxy) DefineModule(docComment string, name symbol.Symbol) *Module {
	return m.Mixin.Get().DefineModule(docComment, name)
}

func (m *MixinProxy) TryDefineModule(docComment string, name symbol.Symbol) *Module {
	return m.Mixin.Get().TryDefineModule(docComment, name)
}

func (m *MixinProxy) DefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	return m.Mixin.Get().DefineMixin(docComment, abstract, name)
}

func (m *MixinProxy) TryDefineMixin(docComment string, abstract bool, name symbol.Symbol) *Mixin {
	return m.Mixin.Get().TryDefineMixin(docComment, abstract, name)
}

func (m *MixinProxy) DefineInterface(docComment string, name symbol.Symbol) *Interface {
	return m.Mixin.Get().DefineInterface(docComment, name)
}

func (m *MixinProxy) TryDefineInterface(docComment string, name symbol.Symbol) *Interface {
	return m.Mixin.Get().TryDefineInterface(docComment, name)
}
