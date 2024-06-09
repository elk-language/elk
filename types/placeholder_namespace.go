package types

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/threadsafe"
)

// Used during typechecking as a placeholder for a future
// module, class, mixin, interface etc.
type PlaceholderNamespace struct {
	ConstantMap
	Replacement ConstantContainer
	Locations   *threadsafe.Slice[*position.Location]
}

func (*PlaceholderNamespace) Parent() ConstantContainer {
	return nil
}

func NewPlaceholderNamespace(name string) *PlaceholderNamespace {
	return &PlaceholderNamespace{
		ConstantMap: MakeConstantMap(name),
		Locations:   threadsafe.NewSlice[*position.Location](),
	}
}

func NewPlaceholderNamespaceWithDetails(
	name string,
	consts *TypeMap,
	subtypes *TypeMap,
	methods *MethodMap,
) *PlaceholderNamespace {
	return &PlaceholderNamespace{
		ConstantMap: ConstantMap{
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
