package types

type Class struct {
	parent   ConstantContainer
	Abstract bool
	Sealed   bool
	ConstantMap
}

func (c *Class) SetAbstract(abstract bool) *Class {
	c.Abstract = abstract
	return c
}

func (c *Class) SetSealed(sealed bool) *Class {
	c.Sealed = sealed
	return c
}

func (c *Class) Parent() ConstantContainer {
	return c.parent
}

func (c *Class) SetParent(parent ConstantContainer) *Class {
	c.parent = parent
	return c
}

func NewClass(name string, parent ConstantContainer) *Class {
	return &Class{
		parent:      parent,
		ConstantMap: MakeConstantMap(name),
	}
}

func NewClassWithDetails(name string, parent ConstantContainer, consts *TypeMap, subtypes *TypeMap, methods *MethodMap) *Class {
	return &Class{
		parent: parent,
		ConstantMap: ConstantMap{
			name:      name,
			constants: consts,
			subtypes:  subtypes,
			methods:   methods,
		},
	}
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
