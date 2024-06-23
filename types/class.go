package types

import (
	"github.com/elk-language/elk/value/symbol"
)

type Class struct {
	parent    Namespace
	abstract  bool
	sealed    bool
	singleton *SingletonClass
	NamespaceBase
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

func (c *Class) Parent() Namespace {
	return c.parent
}

func (c *Class) Singleton() *SingletonClass {
	return c.singleton
}

func (c *Class) Superclass() *Class {
	var currentParent Namespace = c.parent
	for {
		if currentParent == nil {
			return nil
		}
		if currentClass, ok := currentParent.(*Class); ok {
			return currentClass
		}

		currentParent = currentParent.Parent()
	}
}

func (c *Class) SetParent(parent Namespace) {
	c.parent = parent
	superclass := c.Superclass()
	if superclass != nil && c.singleton != nil {
		c.singleton.parent = superclass.singleton
	}
}

func NewClass(name string, parent Namespace, env *GlobalEnvironment) *Class {
	class := &Class{
		NamespaceBase: MakeNamespaceBase(name),
	}
	class.singleton = NewSingletonClass(class, env.StdSubtypeClass(symbol.Class))
	class.SetParent(parent)

	return class
}

func NewClassWithDetails(name string, parent Namespace, consts *TypeMap, subtypes *TypeMap, methods *MethodMap, env *GlobalEnvironment) *Class {
	class := &Class{
		NamespaceBase: NamespaceBase{
			name:      name,
			constants: consts,
			subtypes:  subtypes,
			methods:   methods,
		},
	}
	class.singleton = NewSingletonClass(class, env.StdSubtypeClass(symbol.Class))
	class.SetParent(parent)

	return class
}

func (c *Class) DefineMethod(name string, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(name, params, returnType, throwType, c)
	c.SetMethod(name, method)
	return method
}

func (c *Class) inspect() string {
	return c.name
}

func (c *Class) ToNonLiteral(env *GlobalEnvironment) Type {
	return c
}
