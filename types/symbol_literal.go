package types

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type SymbolLiteral struct {
	Value string
}

func NewSymbolLiteral(value string) *SymbolLiteral {
	return &SymbolLiteral{
		Value: value,
	}
}

func (s *SymbolLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Symbol)
}

func (s *SymbolLiteral) inspect() string {
	return fmt.Sprintf("Std::Symbol(:%s)", value.InspectSymbolContent(s.Value))
}
