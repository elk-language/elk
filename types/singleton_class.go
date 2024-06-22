package types

import (
	"fmt"

	"github.com/elk-language/elk/value/symbol"
)

// Type that represents the singleton class of a mixin, class etc.
type SingletonClass struct {
	AttachedObject Namespace
	Class
}

func NewSingletonClass(attached Namespace, env *GlobalEnvironment) *SingletonClass {
	var parent Namespace
	switch attached.(type) {
	case *Mixin:
		parent = env.StdSubtypeClass(symbol.Mixin)
	case *Class:
		parent = env.StdSubtypeClass(symbol.Class)
	case *Interface:
		parent = env.StdSubtypeClass(symbol.Interface)
	default:
		panic(fmt.Sprintf("invalid object for singleton class: %T", attached))
	}

	singleton := &SingletonClass{
		AttachedObject: attached,
		Class: Class{
			parent:        parent,
			NamespaceBase: MakeNamespaceBase("&" + attached.Name()),
		},
	}
	return singleton
}

func (s *SingletonClass) ToNonLiteral(env *GlobalEnvironment) Type {
	return s
}
