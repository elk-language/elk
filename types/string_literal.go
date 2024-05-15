package types

import "github.com/elk-language/elk/value"

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
		return o == env.StdSubtype("String")
	default:
		return false
	}
}

func (s *StringLiteral) ToNonLiteral(env *GlobalEnvironment) Type {
	return env.StdSubtype("String")
}

func (s *StringLiteral) Inspect() string {
	return value.String(s.Value).Inspect()
}
