package types

import (
	"strings"
)

type Nilable struct {
	Type Type
}

func NewNilable(typ Type) *Nilable {
	return &Nilable{
		Type: typ,
	}
}

func (n *Nilable) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(n, parent) {
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

func (n *Nilable) Copy() *Nilable {
	return &Nilable{
		Type: n.Type,
	}
}

func (n *Nilable) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *Nilable {
	newNilable := n.Copy()
	newNilable.Type = DeepCopyEnv(n.Type, oldEnv, newEnv)
	return newNilable
}
