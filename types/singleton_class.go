package types

// Type that represents the singleton class of a mixin, class etc.
type SingletonClass struct {
	AttachedObject Namespace
	Class
}

func NewSingletonClass(attached Namespace) *SingletonClass {
	return &SingletonClass{
		AttachedObject: attached,
		Class: Class{
			parent:        nil,
			NamespaceBase: MakeNamespaceBase("&" + attached.Name()),
		},
	}
}

func (s *SingletonClass) ToNonLiteral(env *GlobalEnvironment) Type {
	return s
}
