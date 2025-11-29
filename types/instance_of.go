package types

import (
	"strings"
)

type InstanceOf struct {
	Type Type
}

func NewInstanceOf(typ Type) *InstanceOf {
	return &InstanceOf{
		Type: typ,
	}
}

func (i *InstanceOf) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseContinue:
		return leave(i, parent)
	}

	if i.Type.traverse(i, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(i, parent)
}

func (s *InstanceOf) ToNonLiteral(env *GlobalEnvironment) Type {
	return s
}

func (*InstanceOf) IsLiteral() bool {
	return false
}

func (s *InstanceOf) inspect() string {
	var buf strings.Builder

	var addParens bool
	switch s.Type.(type) {
	case *Union, *Intersection, *Not, *SingletonOf:
		addParens = true
	}

	buf.WriteRune('%')
	if addParens {
		buf.WriteRune('(')
	}
	buf.WriteString(Inspect(s.Type))
	if addParens {
		buf.WriteRune(')')
	}
	return buf.String()
}

func (i *InstanceOf) Copy() *InstanceOf {
	return &InstanceOf{
		Type: i.Type,
	}
}

func (i *InstanceOf) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *InstanceOf {
	newInstanceOf := i.Copy()
	newInstanceOf.Type = DeepCopyEnv(i.Type, oldEnv, newEnv)
	return newInstanceOf
}
