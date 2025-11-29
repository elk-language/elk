package types

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Used during typechecking as a placeholder for a future
// constant or type in using statements
type ConstantPlaceholder struct {
	AsName    value.Symbol
	FullName  string
	Container ConstantMap
	Location  *position.Location
	Sibling   *ConstantPlaceholder
	Checked   bool
	Replaced  bool
}

func (c *ConstantPlaceholder) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(c, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(c, parent)
	}
}

func IsConstantPlaceholder(typ Type) bool {
	_, ok := typ.(*ConstantPlaceholder)
	return ok
}

func NewConstantPlaceholder(asName value.Symbol, fullName string, container ConstantMap, location *position.Location) *ConstantPlaceholder {
	return &ConstantPlaceholder{
		AsName:    asName,
		FullName:  fullName,
		Container: container,
		Location:  location,
	}
}

func (p *ConstantPlaceholder) ToNonLiteral(env *GlobalEnvironment) Type {
	return p
}

func (*ConstantPlaceholder) IsLiteral() bool {
	return false
}

func (p *ConstantPlaceholder) inspect() string {
	return fmt.Sprintf("<ConstantPlaceholder: %s>", p.FullName)
}

func (p *ConstantPlaceholder) Copy() *ConstantPlaceholder {
	return &ConstantPlaceholder{
		AsName:    p.AsName,
		FullName:  p.FullName,
		Container: p.Container,
		Location:  p.Location,
		Sibling:   p.Sibling,
		Checked:   p.Checked,
		Replaced:  p.Replaced,
	}
}

func (p *ConstantPlaceholder) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *ConstantPlaceholder {
	if newType, ok := NameToTypeOk(p.FullName, newEnv); ok {
		return newType.(*ConstantPlaceholder)
	}

	newPlaceholder := &ConstantPlaceholder{
		AsName:    p.AsName,
		FullName:  p.FullName,
		Location:  p.Location,
		Checked:   p.Checked,
		Replaced:  p.Replaced,
		Container: make(ConstantMap),
	}
	moduleConstantPath := GetConstantPath(p.FullName)
	parentNamespace := DeepCopyNamespacePath(moduleConstantPath[:len(moduleConstantPath)-1], oldEnv, newEnv)
	parentNamespace.DefineSubtype(value.ToSymbol(moduleConstantPath[len(moduleConstantPath)-1]), newPlaceholder)

	newPlaceholder.Container = ConstantsDeepCopyEnv(p.Container, oldEnv, newEnv)

	return newPlaceholder
}
