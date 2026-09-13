package types

import (
	"encoding/binary"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/value/strings"
	"github.com/elk-language/elk/value/symbol"
)

type CharLiteral struct {
	Value rune
	id    ID
}

func (c *CharLiteral) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(c, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(c, parent)
	}
}

func (c *CharLiteral) ToRef() Ref[*CharLiteral] {
	return Ref[*CharLiteral](c.id)
}

func (c *CharLiteral) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("char:")
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(c.Value))
	d.Write(b)

	return d.Sum64()
}

func (c *CharLiteral) EqualAny(other any) bool {
	o, ok := other.(*CharLiteral)
	if !ok {
		return false
	}

	if c.id > 0 {
		return c.id == o.ID()
	}

	return c.Value == o.Value
}

func (c *CharLiteral) ID() ID {
	return c.id
}

func (c *CharLiteral) SetID(id ID) {
	c.id = id
}

func (c *CharLiteral) StringValue() string {
	return string(c.Value)
}

func NewCharLiteral(value rune) *CharLiteral {
	t := &CharLiteral{
		Value: value,
	}
	return t
}

func (*CharLiteral) ToNonLiteral() Type {
	return Env.StdSubtype(symbol.C_Char)
}

func (*CharLiteral) IsLiteral() bool {
	return true
}

func (c *CharLiteral) inspect() string {
	return strings.InspectChar(c.Value)
}
