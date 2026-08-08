package types

import (
	"strings"
)

type Exact struct {
	Type Namespace
}

func NewExact(typ Namespace) *Exact {
	return &Exact{
		Type: typ,
	}
}

func (n *Exact) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
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

func (n *Exact) ToNonLiteral(env *GlobalEnvironment) Type {
	return n.Type
}

func (*Exact) IsLiteral() bool {
	return false
}

func (n *Exact) inspect() string {
	var buf strings.Builder

	buf.WriteString("exact ")
	buf.WriteString(Inspect(n.Type))

	return buf.String()
}

func (n *Exact) Copy() *Exact {
	return &Exact{
		Type: n.Type,
	}
}

func (n *Exact) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *Exact {
	newExact := n.Copy()
	newExact.Type = DeepCopyEnv(n.Type, oldEnv, newEnv).(Namespace)
	return newExact
}
