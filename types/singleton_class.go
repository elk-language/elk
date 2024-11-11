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

func (s *SingletonClass) Copy() *SingletonClass {
	return &SingletonClass{
		AttachedObject: s.AttachedObject,
		Class:          s.Class,
	}
}

func (s *SingletonClass) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *SingletonClass {
	if newType, ok := NameToTypeOk(s.name, newEnv); ok {
		return newType.(*SingletonClass)
	}

	newAttachedObject := DeepCopyEnv(s.AttachedObject, oldEnv, newEnv).(Namespace)
	if newS := newAttachedObject.Singleton(); newS != nil {
		return newS
	}
	newSingleton := s.Copy()
	newSingleton.AttachedObject = newAttachedObject

	newSingleton.methods = MethodsDeepCopyEnv(s.methods, oldEnv, newEnv)
	newSingleton.subtypes = ConstantsDeepCopyEnv(s.subtypes, oldEnv, newEnv)
	newSingleton.constants = ConstantsDeepCopyEnv(s.constants, oldEnv, newEnv)

	if s.parent != nil {
		newSingleton.parent = DeepCopyEnv(s.parent, oldEnv, newEnv).(Namespace)
	}
	newAttachedObject.SetSingleton(newSingleton)
	return newSingleton
}
