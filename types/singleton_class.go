package types

// Type that represents the singleton class of a mixin, class etc.
type SingletonClass struct {
	AttachedObject Namespace
	Class
}

func (c *SingletonClass) SetParent(parent Namespace) {
	c.parent = parent
}

func NewSingletonClass(attached Namespace, parent Namespace) *SingletonClass {
	singleton := &SingletonClass{
		AttachedObject: attached,
		Class: Class{
			parent:        parent,
			NamespaceBase: MakeNamespaceBase("", "&"+attached.Name()),
		},
	}
	return singleton
}

func (s *SingletonClass) ToNonLiteral(env *GlobalEnvironment) Type {
	return s
}
