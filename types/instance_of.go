package types

import (
	"encoding/binary"
	"strings"

	"github.com/cespare/xxhash/v2"
)

type InstanceOf struct {
	Type Ref[Type]
	id   ID
}

func NewInstanceOf(typ Type) *InstanceOf {
	return &InstanceOf{
		Type: ToRef(typ),
	}
}

func (c *InstanceOf) ToRef() Ref[*InstanceOf] {
	return ToRef(c)
}

func (c *InstanceOf) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("instanceof:")
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(c.Type))
	d.Write(b)

	return d.Sum64()
}

func (c *InstanceOf) EqualAny(other any) bool {
	o, ok := other.(*InstanceOf)
	if !ok {
		return false
	}

	if c.id > 0 {
		return c.id == o.ID()
	}

	return c.Type == o.Type
}

func (c *InstanceOf) ID() ID {
	return c.id
}

func (c *InstanceOf) SetID(id ID) {
	c.id = id
}

func (i *InstanceOf) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseContinue:
		return leave(i, parent)
	}

	if i.Type.Get().traverse(i, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(i, parent)
}

func (s *InstanceOf) ToNonLiteral() Type {
	return s
}

func (*InstanceOf) IsLiteral() bool {
	return false
}

func (s *InstanceOf) inspect() string {
	var buf strings.Builder

	typ := s.Type.Get()
	var addParens bool
	switch typ.(type) {
	case *Union, *Intersection, *Not, *SingletonOf:
		addParens = true
	}

	buf.WriteRune('%')
	if addParens {
		buf.WriteRune('(')
	}
	buf.WriteString(Inspect(typ))
	if addParens {
		buf.WriteRune(')')
	}
	return buf.String()
}

func (i *InstanceOf) Copy() *InstanceOf {
	return &InstanceOf{
		Type: i.Type,
	}
}

func (i *InstanceOf) CopyType() Type {
	return i.Copy()
}
