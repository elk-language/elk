package types

import (
	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/position"
)

type ModulePlaceholder struct {
	Module
}

func NewModulePlaceholder(name string) *ModulePlaceholder {
	return &ModulePlaceholder{
		Module: Module{
			NamespaceBase: MakeNamespaceBase("", name),
		},
	}
}

// Used during typechecking as a placeholder for a future
// module, class, mixin, interface etc.
type NamespacePlaceholder struct {
	name string
	Namespace
	Locations *concurrent.Slice[*position.Location]
}

func NewNamespacePlaceholder(name string) *NamespacePlaceholder {
	return &NamespacePlaceholder{
		name:      name,
		Locations: concurrent.NewSlice[*position.Location](),
		Namespace: NewModulePlaceholder(name),
	}
}

func (p *NamespacePlaceholder) ToNonLiteral(env *GlobalEnvironment) Type {
	return p
}

func (*NamespacePlaceholder) IsLiteral() bool {
	return false
}

func (p *NamespacePlaceholder) inspect() string {
	return p.Name()
}
