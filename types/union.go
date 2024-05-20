package types

import (
	"strings"
)

// Union type represents a list of types.
// A value has to satisfy at least one of the types.
type Union struct {
	Elements []Type
}

func NewUnion(elements ...Type) *Union {
	return &Union{
		Elements: elements,
	}
}

func (u *Union) ToNonLiteral(env *GlobalEnvironment) Type {
	return u
}

func (u *Union) inspect() string {
	var buf strings.Builder
	for i, element := range u.Elements {
		if i != 0 {
			buf.WriteString(" | ")
		}
		var addParens bool
		switch element.(type) {
		case *Union, *Intersection:
			addParens = true
		}

		if addParens {
			buf.WriteRune('(')
		}
		buf.WriteString(Inspect(element))
		if addParens {
			buf.WriteRune(')')
		}
	}
	return buf.String()
}
