package types

import (
	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/position"
)

type PlaceholderModule struct {
	Module
}

func NewPlaceholderModule(name string) *PlaceholderModule {
	return &PlaceholderModule{
		Module: Module{
			NamespaceBase: MakeNamespaceBase("", name),
		},
	}
}

// Used during typechecking as a placeholder for a future
// module, class, mixin, interface etc.
type PlaceholderNamespace struct {
	name string
	Namespace
	Locations *concurrent.Slice[*position.Location]
}

func NewPlaceholderNamespace(name string) *PlaceholderNamespace {
	return &PlaceholderNamespace{
		name:      name,
		Locations: concurrent.NewSlice[*position.Location](),
		Namespace: NewPlaceholderModule(name),
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
