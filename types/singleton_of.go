package types

import (
	"strings"
)

type SingletonOf struct {
	Type Type
}

func NewSingletonOf(typ Type) *SingletonOf {
	return &SingletonOf{
		Type: typ,
	}
}

func (s *SingletonOf) ToNonLiteral(env *GlobalEnvironment) Type {
	return s
}

func (*SingletonOf) IsLiteral() bool {
	return false
}

func (s *SingletonOf) inspect() string {
	var buf strings.Builder

	var addParens bool
	switch s.Type.(type) {
	case *Union, *Intersection, *Not:
		addParens = true
	}

	buf.WriteRune('&')
	if addParens {
		buf.WriteRune('(')
	}
	buf.WriteString(Inspect(s.Type))
	if addParens {
		buf.WriteRune(')')
	}
	return buf.String()
}

func (s *SingletonOf) Copy() *SingletonOf {
	return NewSingletonOf(s.Type)
}

func (s *SingletonOf) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *SingletonOf {
	newSingleton := s.Copy()
	newSingleton.Type = DeepCopyEnv(newSingleton.Type, oldEnv, newEnv)

	return newSingleton
}
