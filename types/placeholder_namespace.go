package types

import (
	"fmt"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Used during typechecking as a placeholder for a future
// module, class, mixin, interface etc.
type PlaceholderNamespace struct {
	NamespaceBase
	Replacement Namespace
	Locations   *concurrent.Slice[*position.Location]
}

func (p *PlaceholderNamespace) IsGeneric() bool {
	return false
}

func (p *PlaceholderNamespace) TypeParameters() []*TypeParameter {
	return nil
}

func (p *PlaceholderNamespace) SetTypeParameters(t []*TypeParameter) {
	panic("cannot set type parameters on a placeholder namespace")
}

func (*PlaceholderNamespace) IsAbstract() bool {
	return false
}

func (*PlaceholderNamespace) IsSealed() bool {
	return false
}

func (*PlaceholderNamespace) IsPrimitive() bool {
	return false
}

func (*PlaceholderNamespace) Singleton() *SingletonClass {
	return nil
}

func (*PlaceholderNamespace) Parent() Namespace {
	return nil
}

func (p *PlaceholderNamespace) SetParent(Namespace) {
	panic(fmt.Sprintf("cannot set the parent of placeholder namespace `%s`", p.Name()))
}

func NewPlaceholderNamespace(name string) *PlaceholderNamespace {
	return &PlaceholderNamespace{
		NamespaceBase: MakeNamespaceBase("", name),
		Locations:     concurrent.NewSlice[*position.Location](),
	}
}

func NewPlaceholderNamespaceWithDetails(
	name string,
	consts *TypeMap,
	subtypes *TypeMap,
	methods *MethodMap,
) *PlaceholderNamespace {
	return &PlaceholderNamespace{
		NamespaceBase: NamespaceBase{
			name:      name,
			constants: consts,
			subtypes:  subtypes,
			methods:   methods,
		},
	}
}

func (p *PlaceholderNamespace) ToNonLiteral(env *GlobalEnvironment) Type {
	return p
}

func (*PlaceholderNamespace) IsLiteral() bool {
	return false
}

func (p *PlaceholderNamespace) inspect() string {
	return p.Name()
}

func (p *PlaceholderNamespace) DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, abstract, sealed, native, name, typeParams, params, returnType, throwType, p)
	p.SetMethod(name, method)
	return method
}
