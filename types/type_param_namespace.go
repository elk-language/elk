package types

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value"
)

type TypeParamNamespace struct {
	docComment string
	ForMethod  bool
	constants  ConstantMap
	subtypes   ConstantMap
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
		constants:  t.constants,
		subtypes:   t.subtypes,
	}
}

func (t *TypeParamNamespace) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *TypeParamNamespace {
	newNamespace := t.Copy()

	newNamespace.constants = ConstantsDeepCopyEnv(t.constants, oldEnv, newEnv)
	newNamespace.subtypes = ConstantsDeepCopyEnv(t.subtypes, oldEnv, newEnv)

	return newNamespace
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

func (t *TypeParamNamespace) SetParent(Namespace) {
	panic("cannot set parent of type param namespaces")
}

func (t *TypeParamNamespace) Singleton() *SingletonClass {
	return nil
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

func (t *TypeParamNamespace) IsGeneric() bool {
	return false
}

func (t *TypeParamNamespace) TypeParameters() []*TypeParameter {
	return nil
}

func (t *TypeParamNamespace) SetTypeParameters([]*TypeParameter) {
	panic("cannot set type parameters on a type parameter namespace")
}

func (t *TypeParamNamespace) Constants() ConstantMap {
	return t.constants
}

func (t *TypeParamNamespace) MethodAliases() MethodAliasMap {
	return nil
}

func (t *TypeParamNamespace) SetMethodAlias(name value.Symbol, method *Method) {
	panic("cannot define method aliases on type param namespaces")
}

func (t *TypeParamNamespace) Constant(name value.Symbol) (Constant, bool) {
	result, ok := t.constants[name]
	return result, ok
}

func (t *TypeParamNamespace) ConstantString(name string) (Constant, bool) {
	return t.Constant(value.ToSymbol(name))
}

func (t *TypeParamNamespace) DefineConstant(name value.Symbol, val Type) {
	t.constants[name] = Constant{
		Type: val,
	}
}

func (t *TypeParamNamespace) DefineConstantWithFullName(name value.Symbol, fullName string, val Type) {
	t.constants[name] = Constant{
		Type:     val,
		FullName: fullName,
	}
}

func (t *TypeParamNamespace) Subtypes() ConstantMap {
	return nil
}

func (t *TypeParamNamespace) Subtype(name value.Symbol) (Constant, bool) {
	result, ok := t.subtypes[name]
	return result, ok
}

func (t *TypeParamNamespace) MustSubtypeString(name string) Type {
	return t.subtypes[value.ToSymbol(name)].Type
}

func (t *TypeParamNamespace) MustSubtype(name value.Symbol) Type {
	return t.subtypes[name].Type
}

func (t *TypeParamNamespace) SubtypeString(name string) (Constant, bool) {
	return t.Subtype(value.ToSymbol(name))
}

func (t *TypeParamNamespace) DefineSubtype(name value.Symbol, val Type) {
	t.subtypes[name] = Constant{
		Type: val,
	}
}

func (t *TypeParamNamespace) DefineSubtypeWithFullName(name value.Symbol, fullName string, val Type) {
	t.subtypes[name] = Constant{
		Type: val,
	}
}

func (t *TypeParamNamespace) Methods() MethodMap {
	return nil
}

func (t *TypeParamNamespace) Method(name value.Symbol) *Method {
	return nil
}

func (t *TypeParamNamespace) MethodString(name string) *Method {
	return nil
}

func (t *TypeParamNamespace) DefineMethod(docComment string, flags bitfield.BitFlag16, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	panic("cannot define methods on type param namespaces")
}

func (t *TypeParamNamespace) SetMethod(name value.Symbol, method *Method) {
}

func (t *TypeParamNamespace) InstanceVariables() TypeMap {
	return nil
}

func (t *TypeParamNamespace) InstanceVariable(name value.Symbol) Type {
	return nil
}

func (t *TypeParamNamespace) InstanceVariableString(name string) Type {
	return nil
}

func (t *TypeParamNamespace) DefineInstanceVariable(name value.Symbol, val Type) {
	panic("cannot define instance variables on type param namespaces")
}

func (t *TypeParamNamespace) DefineClass(docComment string, primitive, abstract, sealed, noinit bool, name value.Symbol, parent Namespace, env *GlobalEnvironment) *Class {
	panic("cannot define classes on type param namespaces")
}

func (t *TypeParamNamespace) DefineModule(docComment string, name value.Symbol, env *GlobalEnvironment) *Module {
	panic("cannot define module on type param namespaces")
}

func (t *TypeParamNamespace) DefineMixin(docComment string, abstract bool, name value.Symbol, env *GlobalEnvironment) *Mixin {
	panic("cannot define mixins on type param namespaces")
}

func (t *TypeParamNamespace) DefineInterface(docComment string, name value.Symbol, env *GlobalEnvironment) *Interface {
	panic("cannot define interfaces on type param namespaces")
}

func (t *TypeParamNamespace) ToNonLiteral(env *GlobalEnvironment) Type {
	return t
}

func (*TypeParamNamespace) IsLiteral() bool {
	return false
}
