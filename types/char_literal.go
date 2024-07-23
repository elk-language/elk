package types

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type CharLiteral struct {
	Value rune
}

func (c *CharLiteral) StringValue() string {
	return string(c.Value)
}

func NewCharLiteral(value rune) *CharLiteral {
	return &CharLiteral{
		Value: value,
	}
}

func (*CharLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Char)
}

func (*CharLiteral) IsLiteral() bool {
	return true
}

func (c *CharLiteral) inspect() string {
	return value.Char(c.Value).Inspect()
}
