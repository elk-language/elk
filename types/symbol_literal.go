package types

import (
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

func (s *SymbolLiteral) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *SymbolLiteral:
		return s.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.Symbol)
	default:
		return false
	}
}

func (s *SymbolLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Symbol)
}

func (s *SymbolLiteral) Inspect() string {
	return value.InspectSymbol(s.Value)
}
