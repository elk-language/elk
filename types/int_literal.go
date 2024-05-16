package types

import "github.com/elk-language/elk/value/symbol"

type IntLiteral struct {
	Value string
}

func NewIntLiteral(value string) *IntLiteral {
	return &IntLiteral{
		Value: value,
	}
}

func (i *IntLiteral) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	switch o := other.(type) {
	case *IntLiteral:
		return i.Value == o.Value
	case *Class:
		return o == env.StdSubtype(symbol.Int)
	default:
		return false
	}
}

func (i *IntLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype(symbol.Int)
}

func (i *IntLiteral) Inspect() string {
	return i.Value
}
