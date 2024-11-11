package types

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Class struct {
	parent         Namespace
	abstract       bool
	sealed         bool
	primitive      bool
	defined        bool
	compiled       bool
	Checked        bool
	singleton      *SingletonClass
	typeParameters []*TypeParameter
	NamespaceBase
}

func (c *Class) IsGeneric() bool {
	return len(c.typeParameters) > 0
}

func (c *Class) TypeParameters() []*TypeParameter {
	return c.typeParameters
}

func (c *Class) SetTypeParameters(t []*TypeParameter) {
	c.typeParameters = t
}

func (c *Class) SetAbstract(abstract bool) *Class {
	c.abstract = abstract
	return c
}

func (c *Class) IsAbstract() bool {
	return c.abstract
}

func (c *Class) SetSealed(sealed bool) *Class {
	c.sealed = sealed
	return c
}

func (c *Class) IsSealed() bool {
	return c.sealed
}

func (c *Class) SetPrimitive(primitive bool) *Class {
	c.primitive = primitive
	return c
}

func (c *Class) IsPrimitive() bool {
	return c.primitive
}

func (c *Class) IsDefined() bool {
	return c.defined
}

func (c *Class) SetDefined(val bool) {
	c.defined = val
}

func (c *Class) IsCompiled() bool {
	return c.compiled
}

func (c *Class) SetCompiled(val bool) {
	c.compiled = val
}

func (c *Class) Parent() Namespace {
	return c.parent
}

func (c *Class) Singleton() *SingletonClass {
	return c.singleton
}

func (c *Class) Superclass() Namespace {
	var currentParent Namespace = c.parent
	for {
		if currentParent == nil {
			return nil
		}
		switch narrowedParent := currentParent.(type) {
		case *Class:
			return narrowedParent
		case *Generic:
			if _, ok := narrowedParent.Namespace.(*Class); ok {
				return narrowedParent
			}
		}

		currentParent = currentParent.Parent()
	}
}

func (c *Class) SetParent(parent Namespace) {
	c.parent = parent
	superclass := c.Superclass()
	if superclass != nil && c.singleton != nil {
		c.singleton.parent = superclass.Singleton()
	}
}

func NewClass(
	docComment string,
	abstract,
	sealed,
	primitive bool,
	name string,
	parent Namespace,
	env *GlobalEnvironment,
) *Class {
	class := &Class{
		primitive:     primitive,
		sealed:        sealed,
		abstract:      abstract,
		defined:       env.Init,
		NamespaceBase: MakeNamespaceBase(docComment, name),
	}
	class.singleton = NewSingletonClass(class, env.StdSubtypeClass(symbol.Class))
	class.SetParent(parent)

	return class
}

func NewClassWithDetails(
	docComment string,
	abstract,
	sealed,
	primitive bool,
	name string,
	parent Namespace,
	consts ConstantMap,
	subtypes ConstantMap,
	methods MethodMap,
	env *GlobalEnvironment,
) *Class {
	class := &Class{
		primitive: primitive,
		abstract:  abstract,
		sealed:    sealed,
		defined:   env.Init,
		NamespaceBase: NamespaceBase{
			docComment: docComment,
			name:       name,
			constants:  consts,
			subtypes:   subtypes,
			methods:    methods,
		},
	}
	class.singleton = NewSingletonClass(class, env.StdSubtypeClass(symbol.Class))
	class.SetParent(parent)

	return class
}

func (c *Class) DefineMethod(docComment string, abstract, sealed, native bool, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, abstract, sealed, native, name, typeParams, params, returnType, throwType, c)
	c.SetMethod(name, method)
	return method
}

func (c *Class) inspect() string {
	return c.name
}

func (c *Class) ToNonLiteral(env *GlobalEnvironment) Type {
	return c
}

func (*Class) IsLiteral() bool {
	return false
}

func (c *Class) Copy() *Class {
	return &Class{
		parent:         c.parent,
		primitive:      c.primitive,
		sealed:         c.sealed,
		abstract:       c.abstract,
		defined:        c.defined,
		compiled:       c.compiled,
		singleton:      c.singleton,
		typeParameters: c.typeParameters,
		NamespaceBase: NamespaceBase{
			docComment: c.docComment,
			name:       c.name,
			constants:  c.constants,
			subtypes:   c.subtypes,
			methods:    c.methods,
		},
	}
}

func (c *Class) DeepCopy(oldEnv, newEnv *GlobalEnvironment) *Class {
	if newType, ok := NameToTypeOk(c.name, newEnv); ok {
		return newType.(*Class)
	}

	newClass := c.Copy()
	classConstantPath := GetConstantPath(c.name)
	parentNamespace := DeepCopyNamespacePath(classConstantPath[:len(classConstantPath)-1], oldEnv, newEnv)
	parentNamespace.DefineSubtype(value.ToSymbol(classConstantPath[len(classConstantPath)-1]), newClass)

	newMethods := make(MethodMap, len(c.methods))
	for methodName, method := range c.methods {
		newMethods[methodName] = method.Copy()
	}
	newClass.methods = newMethods

	newConstants := make(ConstantMap, len(c.constants))
	for constName, constant := range c.constants {
		newConstants[constName] = Constant{
			FullName: constant.FullName,
			Type:     DeepCopy(constant.Type, oldEnv, newEnv),
		}
	}
	newClass.constants = newConstants

	newSubtypes := make(ConstantMap, len(c.subtypes))
	for subtypeName, subtype := range c.subtypes {
		newSubtypes[subtypeName] = Constant{
			FullName: subtype.FullName,
			Type:     DeepCopy(subtype.Type, oldEnv, newEnv),
		}
	}
	newClass.subtypes = newSubtypes

	newTypeParameters := make([]*TypeParameter, len(c.typeParameters))
	for name, typeParam := range c.typeParameters {
		newTypeParameters[name] = typeParam.DeepCopy(oldEnv, newEnv)
	}
	newClass.typeParameters = newTypeParameters

	newClass.parent = DeepCopy(c.parent, oldEnv, newEnv).(Namespace)
	newClass.singleton = NewSingletonClass(
		newClass,
		DeepCopy(newClass.singleton.parent, oldEnv, newEnv).(Namespace),
	)
	return newClass
}
