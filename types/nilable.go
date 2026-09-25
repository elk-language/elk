package types

import (
	"encoding/binary"
	"strings"

	"github.com/cespare/xxhash/v2"
)

type Nilable struct {
	Type Ref[Type]
	id   ID
}

func NewNilable(typ Type) *Nilable {
	return &Nilable{
		Type: ToRef(typ),
	}
}

func (n *Nilable) ToRef() Ref[*Nilable] {
	return ToRef(n)
}

func (n *Nilable) HashUint64() uint64 {
	d := xxhash.New()
	d.WriteString("nilable:")

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(n.Type))
	d.Write(b)

	return d.Sum64()
}

func (n *Nilable) EqualAny(other any) bool {
	o, ok := other.(*Nilable)
	if !ok {
		return false
	}

	if n.id > 0 {
		return n.id == o.ID()
	}

	return n.Type == o.Type
}

func (n *Nilable) ID() ID {
	return n.id
}

func (n *Nilable) SetID(id ID) {
	n.id = id
}

func (n *Nilable) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
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

func (n *Nilable) ToNonLiteral() Type {
	return n
}

func (*Nilable) IsLiteral() bool {
	return false
}

func (n *Nilable) inspect() string {
	var buf strings.Builder

	typ := n.Type.Get()

	var addParens bool
	switch typ.(type) {
	case *Union, *Intersection:
		addParens = true
	}

	if addParens {
		buf.WriteRune('(')
	}
	buf.WriteString(Inspect(typ))
	if addParens {
		buf.WriteRune(')')
	}
	buf.WriteRune('?')
	return buf.String()
}

func (n *Nilable) Copy() *Nilable {
	return &Nilable{
		Type: n.Type,
	}
}

func (n *Nilable) CopyType() Type {
	return n.Copy()
}
