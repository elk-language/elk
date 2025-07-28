package types

import (
	"fmt"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

type Class struct {
	parent         Namespace
	noinit         bool
	abstract       bool
	sealed         bool
	primitive      bool
	defined        bool
	compiled       bool
	Checked        bool
	singleton      *SingletonClass
	typeParameters []*TypeParameter
	IvarIndices    *value.IvarIndices
	NamespaceBase
}

func (c *Class) AttachedObjectName() string {
	if len(c.name) < 1 {
		return ""
	}
	if c.name[0] != '&' {
		return ""
	}

	return c.name[1:]
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

func (c *Class) SetNoInit(noinit bool) *Class {
	c.noinit = noinit
	return c
}

func (c *Class) IsNoInit() bool {
	return c.noinit
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

func getClass(namespace Namespace) Namespace {
	switch namespace := namespace.(type) {
	case *Class:
		return namespace
	case *Generic:
		if _, ok := namespace.Namespace.(*Class); ok {
			return namespace
		}
	case *TemporaryParent:
		return getClass(namespace.Namespace)
	}
	return nil
}

func (c *Class) Superclass() Namespace {
	var currentParent Namespace = c.parent
	for {
		if currentParent == nil {
			return nil
		}
		if class := getClass(currentParent); class != nil {
			return class
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

func (c *Class) RemoveTemporaryParents(env *GlobalEnvironment) {
	if _, ok := c.parent.(*TemporaryParent); !ok {
		return
	}

	c.parent = nil
	c.singleton.parent = env.StdSubtypeClass(symbol.Class)
}

func NewClass(
	docComment string,
	abstract,
	sealed,
	primitive,
	noinit bool,
	name string,
	parent Namespace,
	env *GlobalEnvironment,
) *Class {
	class := &Class{
		primitive:     primitive,
		sealed:        sealed,
		abstract:      abstract,
		noinit:        noinit,
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

func (c *Class) DefineMethod(docComment string, flags bitfield.BitFlag16, name value.Symbol, typeParams []*TypeParameter, params []*Parameter, returnType, throwType Type) *Method {
	method := NewMethod(docComment, flags, name, typeParams, params, returnType, throwType, c)
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
		noinit:         c.noinit,
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
	classConstantPath := GetConstantPath(c.name)
	parentNamespace := DeepCopyNamespacePath(classConstantPath[:len(classConstantPath)-1], oldEnv, newEnv)

	if newType, ok := NameToTypeOk(c.name, newEnv); ok {
		return newType.(*Class)
	}

	newClass := &Class{
		noinit:        c.noinit,
		abstract:      c.abstract,
		sealed:        c.sealed,
		primitive:     c.primitive,
		defined:       c.defined,
		compiled:      c.compiled,
		NamespaceBase: MakeNamespaceBase(c.docComment, c.name),
	}
	classConstantName := classConstantPath[len(classConstantPath)-1]
	parentNamespace.DefineSubtype(value.ToSymbol(classConstantName), newClass)

	newClass.singleton = nil
	newClass.singleton = DeepCopyEnv(c.singleton, oldEnv, newEnv).(*SingletonClass)

	newClass.typeParameters = TypeParametersDeepCopyEnv(c.typeParameters, oldEnv, newEnv)
	newClass.methods = MethodsDeepCopyEnv(c.methods, oldEnv, newEnv)
	newClass.instanceVariables = TypesDeepCopyEnv(c.instanceVariables, oldEnv, newEnv)
	newClass.constants = ConstantsDeepCopyEnv(c.constants, oldEnv, newEnv)
	newClass.subtypes = ConstantsDeepCopyEnv(c.subtypes, oldEnv, newEnv)

	if c.parent != nil {
		newClass.parent = DeepCopyEnv(c.parent, oldEnv, newEnv).(Namespace)
	}
	return newClass
}

// Used for debugging deep copies of types
func (c *Class) inspectInheritance() {
	fmt.Printf("Inheritance: ")
	for p := range Parents(c) {
		fmt.Printf(" -> %s(%T:%p", I(p), p, p)
		switch p := p.(type) {
		case *Generic:
			fmt.Printf(" %T:%p", p.Namespace, p.Namespace)
			fmt.Printf("[&:%p", p.ArgumentMap)
			for _, val := range p.ArgumentMap {
				switch t := val.Type.(type) {
				case *TypeParameter:
					fmt.Printf(" %s:%p", t.InspectSignatureWithColor(), t)
				}
			}
			fmt.Print("]")
			switch n := p.Namespace.(type) {
			case *InterfaceProxy:
				fmt.Printf(" %T:%p", n.Interface, n.Interface)
			case *MixinProxy:
				fmt.Printf(" %T:%p", n.Mixin, n.Mixin)
			}
		}

		fmt.Print(")")
	}
	fmt.Println()
}
