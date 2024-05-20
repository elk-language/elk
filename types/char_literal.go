package types

import (
	"fmt"

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

func (*CharLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Char)
}

func (c *CharLiteral) inspect() string {
	return fmt.Sprintf("Std::Char(%s)", value.Char(c.Value).Inspect())
}
