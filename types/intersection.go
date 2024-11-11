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

func (*Intersection) IsLiteral() bool {
	return false
}

func (u *Intersection) inspect() string {
	var buf strings.Builder
	for i, element := range u.Elements {
		if i != 0 {
			buf.WriteString(" & ")
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

func (i *Intersection) Copy() *Intersection {
	return &Intersection{
		Elements: i.Elements,
	}
}

func (i *Intersection) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *Intersection {
	newUnion := i.Copy()
	newElements := make([]Type, len(i.Elements))
	for i, element := range i.Elements {
		newElements[i] = DeepCopyEnv(element, oldEnv, newEnv)
	}
	newUnion.Elements = newElements

	return newUnion
}
