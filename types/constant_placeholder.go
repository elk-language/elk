package types

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value/symbol"
)

// Used during typechecking as a placeholder for a future
// constant or type in using statements
type ConstantPlaceholder struct {
	AsName    symbol.Symbol
	FullName  string
	Container ConstantMap
	Location  *position.Location
	Sibling   Ref[*ConstantPlaceholder]
	Checked   bool
	Replaced  bool
	id        ID
}

func (c *ConstantPlaceholder) ToRef() Ref[*ConstantPlaceholder] {
	return ToRef(c)
}

func (c *ConstantPlaceholder) HashUint64() uint64 {
	d := xxhash.New()
	d.WriteString("constant:")
	d.WriteString(c.FullName)

	return d.Sum64()
}

func (c *ConstantPlaceholder) EqualAny(other any) bool {
	o, ok := other.(*ConstantPlaceholder)
	if !ok {
		return false
	}

	if c.id > 0 {
		return c.id == o.ID()
	}

	return c.FullName == o.FullName
}

func (c *ConstantPlaceholder) ID() ID {
	return c.id
}

func (c *ConstantPlaceholder) SetID(id ID) {
	c.id = id
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

func NewConstantPlaceholder(asName symbol.Symbol, fullName string, container ConstantMap, location *position.Location) *ConstantPlaceholder {
	return &ConstantPlaceholder{
		AsName:    asName,
		FullName:  fullName,
		Container: container,
		Location:  location,
	}
}

func (p *ConstantPlaceholder) ToNonLiteral() Type {
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

func (p *ConstantPlaceholder) CopyType() Type {
	return p.Copy()
}
