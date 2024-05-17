package types

import (
	"strings"
)

// Intersection type represents a list of types.
// A value has to satisfy all of the types.
type Intersection struct {
	Elements []Type
}

func NewIntersection(elements ...Type) *Intersection {
	return &Intersection{
		Elements: elements,
	}
}

func (u *Intersection) ToNonLiteral(env *GlobalEnvironment) Type {
	return u
}

func (u *Intersection) inspect() string {
	var buf strings.Builder
	for i, element := range u.Elements {
		if i != 0 {
			buf.WriteString(" & ")
		}
		buf.WriteRune('(')
		buf.WriteString(Inspect(element))
		buf.WriteRune(')')
	}
	return buf.String()
}
