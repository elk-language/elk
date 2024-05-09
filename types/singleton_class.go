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

func (s *SingletonClass) IsSupertypeOf(other Type) bool {
	otherSingleton, ok := other.(*SingletonClass)
	if !ok {
		return false
	}

	return otherSingleton.AttachedObject == s.AttachedObject
}

func (s *SingletonClass) Inspect() string {
	return fmt.Sprintf("&%s", s.AttachedObject.Inspect())
}
