package types

import "fmt"

// Type that represents the singleton class of a mixin, class etc.
type SingletonClass struct {
	AttachedObject ConstantContainer
}

func NewSingletonClass(attached ConstantContainer) *SingletonClass {
	return &SingletonClass{
		AttachedObject: attached,
	}
}

func (s *SingletonClass) ToNonLiteral(env *GlobalEnvironment) Type {
	return s
}

func (s *SingletonClass) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	otherSingleton, ok := other.(*SingletonClass)
	if !ok {
		return false
	}

	return otherSingleton.AttachedObject == s.AttachedObject
}

func (s *SingletonClass) inspect() string {
	return fmt.Sprintf("&%s", Inspect(s.AttachedObject))
}
