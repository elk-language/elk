package types

import (
	"github.com/elk-language/elk/value"
)

type TypeParamNamespace struct {
	name      string
	constants *TypeMap
	subtypes  *TypeMap
}

func NewTypeParamNamespace(name string) *TypeParamNamespace {
	return &TypeParamNamespace{
		name:      name,
		constants: NewTypeMap(),
		subtypes:  NewTypeMap(),
	}
}

func (t *TypeParamNamespace) Name() string {
	return t.name
}

func (t *TypeParamNamespace) inspect() string {
	return t.name
}

func (t *TypeParamNamespace) DocComment() string {
	return ""
}

func (t *TypeParamNamespace) SetDocComment(string) {
	panic("cannot set doc comment on closures")
}

func (t *TypeParamNamespace) AppendDocComment(string) {
	panic("cannot append doc comment on closures")
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

func (t *TypeParamNamespace) Constants() *TypeMap {
	return t.constants
}

func (t *TypeParamNamespace) Constant(name value.Symbol) Type {
	result, _ := t.constants.Get(name)
	return result
}

func (t *TypeParamNamespace) ConstantString(name string) Type {
	return t.Constant(value.ToSymbol(name))
}

func (t *TypeParamNamespace) DefineConstant(name value.Symbol, val Type) {
	t.constants.Set(name, val)
}

func (t *TypeParamNamespace) Subtypes() *TypeMap {
	return nil
}

func (t *TypeParamNamespace) Subtype(name value.Symbol) Type {
	result, _ := t.subtypes.Get(name)
	return result
}

func (t *TypeParamNamespace) SubtypeString(name string) Type {
	return t.Subtype(value.ToSymbol(name))
}

func (t *TypeParamNamespace) DefineSubtype(name value.Symbol, val Type) {
	t.subtypes.Set(name, val)
}

func (t *TypeParamNamespace) Methods() *MethodMap {
	return nil
}

func (t *TypeParamNamespace) Method(name value.Symbol) *Method {
	return nil
}

func (t *TypeParamNamespace) MethodString(name string) *Method {
	return nil
}

func (t *TypeParamNamespace) DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	panic("cannot define methods on type param namespaces")
}

func (t *TypeParamNamespace) SetMethod(name value.Symbol, method *Method) {
}

func (t *TypeParamNamespace) InstanceVariables() *TypeMap {
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

func (t *TypeParamNamespace) DefineClass(docComment string, primitive, abstract, sealed bool, name value.Symbol, parent Namespace, env *GlobalEnvironment) *Class {
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
