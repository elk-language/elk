package types

import (
	"strings"

	"github.com/elk-language/elk/value/symbol"
)

func ToNilable(typ Type, env *GlobalEnvironment) Type {
	switch t := typ.(type) {
	case *Nilable:
		return typ
	case *Union:
		var newElements []Type
		for _, element := range t.Elements {
			if elementClass, ok := element.(*Class); ok && elementClass == env.StdSubtype(symbol.Nil) {
				return typ
			}
			e := ToNilable(element, env)
			newElements = append(newElements, e)
		}
		return NewUnion(newElements...)
	case *Intersection:
		for _, element := range t.Elements {
			if elementClass, ok := element.(*Class); ok && elementClass == env.StdSubtype(symbol.Nil) {
				return typ
			}
		}
		return NewNilable(typ)
	default:
		return NewNilable(typ)
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
