package types

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type CharLiteral struct {
	Value rune
}

func NewCharLiteral(value rune) *CharLiteral {
	return &CharLiteral{
		Value: value,
	}
}

func (c *CharLiteral) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *CharLiteral:
		return c.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.Char)
	default:
		return false
	}
}

func (*CharLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Char)
}

func (c *CharLiteral) Inspect() string {
	return value.Char(c.Value).Inspect()
}
