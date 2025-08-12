package types

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value"
)

// Type that represents the singleton class of a mixin, class etc.
type SingletonClass struct {
	AttachedObject Namespace
	Class
}

func (c *SingletonClass) SetParent(parent Namespace) {
	c.parent = parent
}

func (c *SingletonClass) RemoveTemporaryParents(env *GlobalEnvironment) {
	if _, ok := c.parent.(*TemporaryParent); !ok {
		return
	}

	c.parent = c.Superclass()
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

func (c *SingletonClass) DefineMethod(docComment string, flags bitfield.BitFlag16, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, flags, name, typeParams, params, returnType, throwType, c)
	c.SetMethod(name, method)
	return method
}

func (s *SingletonClass) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *SingletonClass {
	fullConstantName := s.name[1:]
	if newType, ok := NameToConstantOk(fullConstantName, newEnv); ok {
		return newType.(*SingletonClass)
	}

	newAttachedObject := DeepCopyEnv(s.AttachedObject, oldEnv, newEnv).(Namespace)
	if newS := newAttachedObject.Singleton(); newS != nil {
		return newS
	}
	newSingleton := &SingletonClass{
		Class: Class{
			primitive:     s.primitive,
			sealed:        s.sealed,
			abstract:      s.abstract,
			defined:       s.defined,
			native:        s.native,
			compiled:      s.compiled,
			NamespaceBase: MakeNamespaceBase(s.docComment, s.name),
		},
	}
	newSingleton.AttachedObject = newAttachedObject
	singletonConstantPath := GetConstantPath(fullConstantName)
	parentNamespace := DeepCopyNamespacePath(singletonConstantPath[:len(singletonConstantPath)-1], oldEnv, newEnv)
	singletonConstantName := singletonConstantPath[len(singletonConstantPath)-1]
	parentNamespace.DefineConstant(value.ToSymbol(singletonConstantName), newSingleton)

	newSingleton.methods = MethodsDeepCopyEnv(s.methods, oldEnv, newEnv)
	newSingleton.instanceVariables = TypesDeepCopyEnv(s.instanceVariables, oldEnv, newEnv)
	newSingleton.subtypes = ConstantsDeepCopyEnv(s.subtypes, oldEnv, newEnv)
	newSingleton.constants = ConstantsDeepCopyEnv(s.constants, oldEnv, newEnv)

	if s.parent != nil {
		newSingleton.parent = DeepCopyEnv(s.parent, oldEnv, newEnv).(Namespace)
	}
	newAttachedObject.SetSingleton(newSingleton)
	return newSingleton
}
