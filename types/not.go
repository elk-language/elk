package types

import (
	"encoding/binary"
	"strings"

	"github.com/cespare/xxhash/v2"
)

type Not struct {
	Type Ref[Type]
	id   ID
}

func (n *Not) ToRef() Ref[*Not] {
	return ToRef(n)
}

func (n *Not) HashUint64() uint64 {
	d := xxhash.New()
	d.WriteString("not:")

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(n.Type))
	d.Write(b)

	return d.Sum64()
}

func (n *Not) EqualAny(other any) bool {
	o, ok := other.(*Not)
	if !ok {
		return false
	}

	if n.id > 0 {
		return n.id == n.ID()
	}

	return n.Type == o.Type
}

func (n *Not) ID() ID {
	return n.id
}

func (n *Not) SetID(id ID) {
	n.id = id
}

func (n *Not) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(Void{}, parent) {
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

func NewNot(typ Type) *Not {
	return &Not{
		Type: ToRef(typ),
	}
}

func NewNotRef(typ Ref[Type]) *Not {
	return &Not{
		Type: typ,
	}
}

func (n *Not) ToNonLiteral() Type {
	return n
}

func (*Not) IsLiteral() bool {
	return false
}

func (n *Not) inspect() string {
	var buf strings.Builder

	typ := n.Type.Get()
	var addParens bool
	switch typ.(type) {
	case *Union, *Intersection, *Not:
		addParens = true
	}

	buf.WriteRune('~')
	if addParens {
		buf.WriteRune('(')
	}
	buf.WriteString(Inspect(typ))
	if addParens {
		buf.WriteRune(')')
	}
	return buf.String()
}

func (n *Not) Copy() *Not {
	return &Not{
		Type: n.Type,
	}
}

func (n *Not) CopyType() Type {
	return n.Copy()
}
