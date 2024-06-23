package types

import (
	"fmt"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/position"
)

// Used during typechecking as a placeholder for a future
// module, class, mixin, interface etc.
type PlaceholderNamespace struct {
	NamespaceBase
	Replacement Namespace
	Locations   *concurrent.Slice[*position.Location]
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
		NamespaceBase: MakeNamespaceBase(name),
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

func (p *PlaceholderNamespace) inspect() string {
	return p.Name()
}

func (p *PlaceholderNamespace) DefineMethod(name string, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(name, params, returnType, throwType, p)
	p.SetMethod(name, method)
	return method
}
