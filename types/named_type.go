package types

import (
	"encoding/binary"

	"github.com/cespare/xxhash/v2"
)

type NamedType struct {
	Name string
	Type Ref[Type]
	id   ID
}

func NewNamedType(name string, typ Type) *NamedType {
	return &NamedType{
		Name: name,
		Type: ToRef(typ),
	}
}

func (s *NamedType) ToRef() Ref[*NamedType] {
	return Ref[*NamedType](s.id)
}

func (s *NamedType) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("named:")
	d.WriteString(s.Name)

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(s.Type))
	d.Write(b)

	return d.Sum64()
}

func (s *NamedType) EqualAny(other any) bool {
	o, ok := other.(*NamedType)
	if !ok {
		return false
	}

	if s.id > 0 {
		return s.id == o.id
	}

	return s.Name == o.Name && s.Type == o.Type
}

func (s *NamedType) ID() ID {
	return s.id
}

func (s *NamedType) SetID(id ID) {
	s.id = id
}

func (n *NamedType) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(n, parent)
	}
}

func (n *NamedType) ToNonLiteral() Type {
	return n
}

func (*NamedType) IsLiteral() bool {
	return false
}

func (n *NamedType) inspect() string {
	return n.Name
}

func (n *NamedType) Copy() *NamedType {
	return &NamedType{
		Name: n.Name,
		Type: n.Type,
	}
}
