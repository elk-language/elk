package types

import (
	"strings"

	"github.com/elk-language/elk/value/symbol"
)

func IsNilable(typ Type, env *GlobalEnvironment) bool {
	switch t := typ.(type) {
	case *Nilable:
		return true
	case *Class:
		if t == env.StdSubtype(symbol.Nil) {
			return true
		}
		return false
	case *Union:
		for _, element := range t.Elements {
			if IsNilable(element, env) {
				return true
			}
		}
		return false
	case *Intersection:
		for _, element := range t.Elements {
			if IsNilable(element, env) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

type Nilable struct {
	Type Type
}

func NewNilable(typ Type) *Nilable {
	return &Nilable{
		Type: typ,
	}
}

func (n *Nilable) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func (*Nilable) IsLiteral() bool {
	return false
}

func (n *Nilable) inspect() string {
	var buf strings.Builder

	var addParens bool
	switch n.Type.(type) {
	case *Union, *Intersection:
		addParens = true
	}

	if addParens {
		buf.WriteRune('(')
	}
	buf.WriteString(Inspect(n.Type))
	if addParens {
		buf.WriteRune(')')
	}
	buf.WriteRune('?')
	return buf.String()
}
