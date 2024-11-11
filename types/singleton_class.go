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

func (s *SingletonClass) DeepCopy(oldEnv, newEnv *GlobalEnvironment) *SingletonClass {
	if newType, ok := NameToTypeOk(s.name, newEnv); ok {
		return newType.(*SingletonClass)
	}

	newSingleton := s.Copy()
	newSingleton.AttachedObject = DeepCopy(newSingleton.AttachedObject, oldEnv, newEnv).(Namespace)

	newMethods := make(MethodMap, len(s.methods))
	for methodName, method := range s.methods {
		newMethods[methodName] = method.Copy()
	}
	newSingleton.methods = newMethods

	newConstants := make(ConstantMap, len(s.constants))
	for constName, constant := range s.constants {
		newConstants[constName] = Constant{
			FullName: constant.FullName,
			Type:     DeepCopy(constant.Type, oldEnv, newEnv),
		}
	}
	newSingleton.constants = newConstants

	newSubtypes := make(ConstantMap, len(s.subtypes))
	for subtypeName, subtype := range s.subtypes {
		newSubtypes[subtypeName] = Constant{
			FullName: subtype.FullName,
			Type:     DeepCopy(subtype.Type, oldEnv, newEnv),
		}
	}
	newSingleton.subtypes = newSubtypes

	newSingleton.parent = DeepCopy(s.parent, oldEnv, newEnv).(Namespace)
	return newSingleton
}
