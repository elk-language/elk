package types

import "github.com/elk-language/elk/value"

// Represents a method alias, a second
// name of a particular method.
type MethodAlias struct {
	Compiled bool
	Method   *Method
}

func NewMethodAlias(method *Method) *MethodAlias {
	return &MethodAlias{
		Method: method,
	}
}

func (m *MethodAlias) IsDefinable() bool {
	return !m.Compiled
}

type MethodAliasMap = map[value.Symbol]*MethodAlias
