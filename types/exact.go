package types

import (
	"encoding/binary"
	"strings"

	"github.com/cespare/xxhash/v2"
)

type Exact struct {
	Type Ref[Namespace]
	id   ID
}

func NewExact(typ Namespace) *Exact {
	return &Exact{
		Type: ToRef(typ),
	}
}

func (c *Exact) ToRef() Ref[*Exact] {
	return ToRef(c)
}

func (c *Exact) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("exact:")
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(c.Type))
	d.Write(b)

	return d.Sum64()
}

func (c *Exact) EqualAny(other any) bool {
	o, ok := other.(*Exact)
	if !ok {
		return false
	}

	if c.id > 0 {
		return c.id == o.ID()
	}

	return c.Type == o.Type
}

func (c *Exact) ID() ID {
	return c.id
}

func (c *Exact) SetID(id ID) {
	c.id = id
}

func (n *Exact) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseContinue:
		return leave(n, parent)
	}

	if n.Type.Get().traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func (n *Exact) ToNonLiteral() Type {
	return n.Type.Get()
}

func (*Exact) IsLiteral() bool {
	return false
}

func (n *Exact) inspect() string {
	var buf strings.Builder

	buf.WriteString("exact ")
	buf.WriteString(Inspect(n.Type.Get()))

	return buf.String()
}

func (n *Exact) Copy() *Exact {
	return &Exact{
		Type: n.Type,
		id:   n.id,
	}
}
