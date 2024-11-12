package types

import (
	"strings"
)

type Not struct {
	Type Type
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
