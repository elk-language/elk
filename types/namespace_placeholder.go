package types

import (
	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
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

func (m *ModulePlaceholder) Copy() *ModulePlaceholder {
	return &ModulePlaceholder{
		Module: m.Module,
	}
}

func (m *ModulePlaceholder) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *ModulePlaceholder {
	newMod := m.Copy()

	newMod.methods = MethodsDeepCopyEnv(m.methods, oldEnv, newEnv)
	newMod.subtypes = ConstantsDeepCopyEnv(m.subtypes, oldEnv, newEnv)
	newMod.constants = ConstantsDeepCopyEnv(m.constants, oldEnv, newEnv)
	newMod.parent = DeepCopyEnv(m.parent, oldEnv, newEnv).(Namespace)

	return newMod
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

func (n *NamespacePlaceholder) Copy() *NamespacePlaceholder {
	return &NamespacePlaceholder{
		name:      n.name,
		Namespace: n.Namespace,
		Locations: n.Locations,
	}
}

func (n *NamespacePlaceholder) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) Namespace {
	if newType, ok := NameToTypeOk(n.name, newEnv); ok {
		return newType.(Namespace)
	}

	newNamespace := &NamespacePlaceholder{
		name:      n.name,
		Locations: n.Locations,
	}
	moduleConstantPath := GetConstantPath(n.name)
	parentNamespace := DeepCopyNamespacePath(moduleConstantPath[:len(moduleConstantPath)-1], oldEnv, newEnv)
	parentNamespace.DefineSubtype(value.ToSymbol(moduleConstantPath[len(moduleConstantPath)-1]), newNamespace)

	newNamespace.Namespace = DeepCopyEnv(n.Namespace, oldEnv, newEnv).(Namespace)

	return newNamespace
}
