package types

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type SymbolLiteral struct {
	Value string
}

func (s *SymbolLiteral) StringValue() string {
	return s.Value
}

func NewSymbolLiteral(value string) *SymbolLiteral {
	return &SymbolLiteral{
		Value: value,
	}
}

func (s *SymbolLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Symbol)
}

func (*SymbolLiteral) IsLiteral() bool {
	return true
}

func (s *SymbolLiteral) inspect() string {
	return fmt.Sprintf(":%s", value.InspectSymbolContent(s.Value))
}
