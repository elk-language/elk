package types

import (
	"slices"
	"strings"

	"github.com/cespare/xxhash/v2"
)

type GenericNamedType struct {
	Name           string
	Type           Ref[Type]
	TypeParameters []Ref[*TypeParameter]
	id             ID
}

func NewGenericNamedType(name string, typ Type, typeVars []Ref[*TypeParameter]) *GenericNamedType {
	t := &GenericNamedType{
		Name:           name,
		Type:           ToRef(typ),
		TypeParameters: typeVars,
	}
	return t
}

func (c *GenericNamedType) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("gnamed:")
	d.WriteString(c.Name)

	return d.Sum64()
}

func (c *GenericNamedType) EqualAny(other any) bool {
	o, ok := other.(*GenericNamedType)
	if !ok {
		return false
	}

	if c.id > 0 {
		return c.id == o.ID()
	}

	return c.Name == o.Name
}

func (c *GenericNamedType) ID() ID {
	return c.id
}

func (c *GenericNamedType) SetID(id ID) {
	c.id = id
}

func (g *GenericNamedType) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(g, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseContinue:
		return leave(g, parent)
	}

	if g.Type.Get().traverse(g, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	for _, typeParam := range g.TypeParameters {
		if typeParam.Get().traverse(g, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(g, parent)
}

func (g *GenericNamedType) ToNonLiteral() Type {
	return g
}

func (*GenericNamedType) IsLiteral() bool {
	return false
}

func (g *GenericNamedType) inspect() string {
	if len(g.TypeParameters) < 1 {
		return g.Name
	}

	buffer := new(strings.Builder)

	buffer.WriteString(g.Name)
	buffer.WriteRune('[')
	for i, typeVarRef := range g.TypeParameters {
		if i > 0 {
			buffer.WriteString(", ")
		}
		typeVar := typeVarRef.Get()
		switch typeVar.Variance {
		case COVARIANT:
			buffer.WriteRune('+')
		case CONTRAVARIANT:
			buffer.WriteRune('-')
		case BIVARIANT:
			buffer.WriteString("+-")
		}

		buffer.WriteString(typeVar.Name.String())
		if !IsNeverRef(typeVar.LowerBound) {
			buffer.WriteString(" > ")
			buffer.WriteString(Inspect(typeVar.LowerBound.Get()))
		}
		if !IsAnyRef(typeVar.UpperBound) {
			buffer.WriteString(" < ")
			buffer.WriteString(Inspect(typeVar.UpperBound.Get()))
		}
	}
	buffer.WriteRune(']')
	return buffer.String()
}

func (g *GenericNamedType) Copy() *GenericNamedType {
	return &GenericNamedType{
		Name:           g.Name,
		Type:           g.Type,
		TypeParameters: slices.Clone(g.TypeParameters),
	}
}

func (g *GenericNamedType) CopyType() Type {
	return g.Copy()
}
