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

func (c *Class) SetSingleton(singleton *SingletonClass) {
	c.singleton = singleton
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

func (c *Class) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *Class {
	if newType, ok := NameToTypeOk(c.name, newEnv); ok {
		return newType.(*Class)
	}

	newClass := c.Copy()
	newClass.singleton = nil

	classConstantPath := GetConstantPath(c.name)
	parentNamespace := DeepCopyNamespacePath(classConstantPath[:len(classConstantPath)-1], oldEnv, newEnv)
	classConstantName := classConstantPath[len(classConstantPath)-1]
	parentNamespace.DefineSubtype(value.ToSymbol(classConstantName), newClass)

	newClass.methods = MethodsDeepCopyEnv(c.methods, oldEnv, newEnv)
	newClass.constants = ConstantsDeepCopyEnv(c.constants, oldEnv, newEnv)
	newClass.subtypes = ConstantsDeepCopyEnv(c.subtypes, oldEnv, newEnv)
	newClass.typeParameters = TypeParametersDeepCopyEnv(c.typeParameters, oldEnv, newEnv)

	if c.parent != nil {
		newClass.parent = DeepCopyEnv(c.parent, oldEnv, newEnv).(Namespace)
	}
	newClass.singleton = NewSingletonClass(
		newClass,
		DeepCopyEnv(newClass.singleton.parent, oldEnv, newEnv).(Namespace),
	)
	return newClass
}
