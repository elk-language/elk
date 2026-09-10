package types

import (
	"strings"
)

// Intersection type represents a list of types.
// A value has to satisfy all of the types.
type Intersection struct {
	Elements []Ref[Type]
	id       ID
}

func NewIntersection(elements ...Ref[Type]) *Intersection {
	return &Intersection{
		Elements: elements,
	}
}

func (i *Intersection) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseContinue:
		return leave(i, parent)
	}

	for _, element := range i.Elements {
		if element.traverse(i, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(i, parent)
}

func (u *Intersection) ToNonLiteral() Type {
	return u
}

func (*Intersection) IsLiteral() bool {
	return false
}

func (u *Intersection) inspect() string {
	var buf strings.Builder
	for i, elementRef := range u.Elements {
		element := elementRef.Get()
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
		id:       i.id,
	}
}
