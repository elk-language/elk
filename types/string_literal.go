package types

import (
	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/value/strings"
	"github.com/elk-language/elk/value/symbol"
)

type StringLiteral struct {
	Value string
	id    ID
}

func (s *StringLiteral) ToRef() Ref[*StringLiteral] {
	return Ref[*StringLiteral](s.id)
}

func (s *StringLiteral) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("str:")
	d.WriteString(s.Value)

	return d.Sum64()
}

func (s *StringLiteral) EqualAny(other any) bool {
	o, ok := other.(*StringLiteral)
	if !ok {
		return false
	}

	if s.id > 0 {
		return s.id == o.id
	}

	return s.Value == o.Value
}

func (s *StringLiteral) ID() ID {
	return s.id
}

func (s *StringLiteral) SetID(id ID) {
	s.id = id
}

func (s *StringLiteral) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(s, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(s, parent)
	}
}

func (s *StringLiteral) StringValue() string {
	return s.Value
}

func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{
		Value: value,
	}
}

func (s *StringLiteral) IsSubtypeOf(other Type) bool {
	switch o := other.(type) {
	case *StringLiteral:
		return s.Value == o.Value
	case *Class:
		return o == Env.StdSubtype(symbol.C_String)
	default:
		return false
	}
}

func (s *StringLiteral) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_String)
}

func (*StringLiteral) IsLiteral() bool {
	return true
}

func (s *StringLiteral) inspect() string {
	return strings.InspectString(s.Value)
}
