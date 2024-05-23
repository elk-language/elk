package types

import (
	"slices"
	"strings"

	"github.com/elk-language/elk/value/symbol"
)

func ToNonNilable(typ Type, env *GlobalEnvironment) Type {
	switch t := typ.(type) {
	case *Nilable:
		return t.Type
	case *Class:
		if t == env.StdSubtype(symbol.Nil) {
			return Never{}
		}
		return t
	case *Union:
		var newElements []Type
		for _, element := range t.Elements {
			nonNilable := ToNonNilable(element, env)
			if IsNever(nonNilable) {
				continue
			}

			newElements = append(newElements, element)
		}
		if len(newElements) == 0 {
			return Never{}
		}
		return NewUnion(newElements...)
	case *Intersection:
		for _, element := range t.Elements {
			nonNilable := ToNonNilable(element, env)
			if IsNever(nonNilable) {
				return Never{}
			}
		}
		return t
	default:
		return t
	}
}

func ToNilable(typ Type, env *GlobalEnvironment) Type {
	if IsNilable(typ, env) {
		return typ
	}

	switch t := typ.(type) {
	case *Union:
		newElements := slices.Clone(t.Elements)
		newElements = append(newElements, env.StdSubtype(symbol.Nil))
		return NewUnion(newElements...)
	case *Intersection:
		return NewNilable(typ)
	default:
		return NewNilable(typ)
	}
}

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
