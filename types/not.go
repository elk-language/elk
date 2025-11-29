package types

import (
	"strings"
)

type Not struct {
	Type Type
}

func (n *Not) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(Void{}, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseContinue:
		return leave(n, parent)
	}

	if n.Type.traverse(n, enter, leave) == TraverseBreak {
		return TraverseBreak
	}

	return leave(n, parent)
}

func NewNot(typ Type) *Not {
	return &Not{
		Type: typ,
	}
}

func (n *Not) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func (*Not) IsLiteral() bool {
	return false
}

func (n *Not) inspect() string {
	var buf strings.Builder

	var addParens bool
	switch n.Type.(type) {
	case *Union, *Intersection, *Not:
		addParens = true
	}

	buf.WriteRune('~')
	if addParens {
		buf.WriteRune('(')
	}
	buf.WriteString(Inspect(n.Type))
	if addParens {
		buf.WriteRune(')')
	}
	return buf.String()
}

func (n *Not) Copy() *Not {
	return &Not{
		Type: n.Type,
	}
}

func (n *Not) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *Not {
	newNot := n.Copy()
	newNot.Type = DeepCopyEnv(n.Type, oldEnv, newEnv)
	return newNot
}
