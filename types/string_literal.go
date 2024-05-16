package types

import (
	"fmt"

	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type StringLiteral struct {
	Value string
}

func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{
		Value: value,
	}
}

func (s *StringLiteral) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *StringLiteral:
		return s.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.String)
	default:
		return false
	}
}

func (s *StringLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.String)
}

func (s *StringLiteral) Inspect() string {
	return fmt.Sprintf("Std::String(%s)", value.String(s.Value).Inspect())
}
