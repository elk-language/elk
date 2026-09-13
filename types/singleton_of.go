package types

import (
	"encoding/binary"
	"strings"

	"github.com/cespare/xxhash/v2"
)

type SingletonOf struct {
	Type Ref[Type]
	id   ID
}

func NewSingletonOf(typ Type) *SingletonOf {
	return &SingletonOf{
		Type: ToRef(typ),
	}
}

func (s *SingletonOf) ToRef() Ref[*SingletonOf] {
	return Ref[*SingletonOf](s.id)
}

func (s *SingletonOf) HashUint64() uint64 {
	d := xxhash.New()
	d.WriteString("singletonof:")

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(s.Type))
	d.Write(b)

	return d.Sum64()
}

func (s *SingletonOf) EqualAny(other any) bool {
	o, ok := other.(*SingletonOf)
	if !ok {
		return false
	}

	if s.id > 0 {
		return s.id == o.ID()
	}

	return s.Type == o.Type
}

func (s *SingletonOf) ID() ID {
	return s.id
}

func (s *SingletonOf) SetID(id ID) {
	s.id = id
}

func (s *SingletonOf) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(s, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseContinue:
		return leave(s, parent)
	}

	if s.Type.Get().traverse(s, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(s, parent)
}

func (s *SingletonOf) ToNonLiteral() Type {
	return s
}

func (*SingletonOf) IsLiteral() bool {
	return false
}

func (s *SingletonOf) inspect() string {
	var buf strings.Builder

	var addParens bool
	typ := s.Type.Get()
	switch typ.(type) {
	case *Union, *Intersection, *Not:
		addParens = true
	}

	buf.WriteRune('&')
	if addParens {
		buf.WriteRune('(')
	}
	buf.WriteString(Inspect(typ))
	if addParens {
		buf.WriteRune(')')
	}
	return buf.String()
}

func (s *SingletonOf) Copy() *SingletonOf {
	return &SingletonOf{
		Type: s.Type,
		id:   s.id,
	}
}
