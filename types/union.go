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

func (u *Union) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(u, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseContinue:
		return leave(u, parent)
	}

	for _, element := range u.Elements {
		if element.traverse(u, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(u, parent)
}

func (u *Union) ToNonLiteral(env *GlobalEnvironment) Type {
	return u
}

func (*Union) IsLiteral() bool {
	return false
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

func (u *Union) Copy() *Union {
	return &Union{
		Elements: u.Elements,
	}
}

func (u *Union) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *Union {
	newUnion := u.Copy()
	newElements := make([]Type, len(u.Elements))
	for i, element := range u.Elements {
		newElements[i] = DeepCopyEnv(element, oldEnv, newEnv)
	}
	newUnion.Elements = newElements

	return newUnion
}
