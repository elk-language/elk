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

func (c *Class) IsCompiled() bool {
	return c.compiled
}

func (c *Class) SetCompiled(compiled bool) {
	c.compiled = compiled
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
		compiled:      env.Init,
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
		compiled:  env.Init,
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
