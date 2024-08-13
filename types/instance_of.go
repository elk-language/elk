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

	buf.WriteRune('^')
	if addParens {
		buf.WriteRune('(')
	}
	buf.WriteString(Inspect(s.Type))
	if addParens {
		buf.WriteRune(')')
	}
	return buf.String()
}
