package types

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/value/symbol"
)

type SymbolLiteral struct {
	Value string
	id    ID
}

func (s *SymbolLiteral) Copy() *SymbolLiteral {
	return &SymbolLiteral{
		Value: s.Value,
	}
}

func (s *SymbolLiteral) CopyType() Type {
	return s.Copy()
}

func (s *SymbolLiteral) ToRef() Ref[*SymbolLiteral] {
	return Ref[*SymbolLiteral](s.id)
}

func (s *SymbolLiteral) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("str:")
	d.WriteString(s.Value)

	return d.Sum64()
}

func (s *SymbolLiteral) EqualAny(other any) bool {
	o, ok := other.(*SymbolLiteral)
	if !ok {
		return false
	}

	if s.id > 0 {
		return s.id == o.id
	}

	return s.Value == o.Value
}

func (s *SymbolLiteral) ID() ID {
	return s.id
}

func (s *SymbolLiteral) SetID(id ID) {
	s.id = id
}

func (s *SymbolLiteral) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(s, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(s, parent)
	}
}

func (s *SymbolLiteral) StringValue() string {
	return s.Value
}

func NewSymbolLiteral(value string) *SymbolLiteral {
	return &SymbolLiteral{
		Value: value,
	}
}

func (s *SymbolLiteral) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Symbol)
}

func (*SymbolLiteral) IsLiteral() bool {
	return true
}

func (s *SymbolLiteral) inspect() string {
	return fmt.Sprintf(":%s", symbol.InspectSymbolContent(s.Value))
}
