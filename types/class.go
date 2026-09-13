package types

import (
	"fmt"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/value/ivar"
	"github.com/elk-language/elk/value/symbol"
)

type Class struct {
	parent         Ref[Namespace]
	immutable      bool
	noinit         bool
	abstract       bool
	sealed         bool
	primitive      bool
	native         bool
	defined        bool
	compiled       bool
	Checked        bool
	singleton      Ref[*SingletonClass]
	typeParameters []Ref[*TypeParameter]
	ivarIndices    *ivar.IvarIndices
	Children       ds.Set[Ref[*Class]]
	id             ID
	NamespaceBase
}

var _ Namespace = &Class{}

func (c *Class) ToRef() Ref[*Class] {
	return Ref[*Class](c.id)
}

func (c *Class) EqualAny(other any) bool {
	o, ok := other.(*Class)
	if !ok {
		return false
	}

	if c.id > 0 {
		return c.id == o.ID()
	}

	return c.name == o.name
}

func (c *Class) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(c, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(c, parent)
	}
}

func (c *Class) IvarIndices() *ivar.IvarIndices {
	return c.ivarIndices
}

func (c *Class) SetIvarIndices(in *ivar.IvarIndices) {
	c.ivarIndices = in
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

func (c *Class) TypeParameters() []Ref[*TypeParameter] {
	return c.typeParameters
}

func (c *Class) SetTypeParameters(t []Ref[*TypeParameter]) {
	c.typeParameters = t
}

func (c *Class) SetNoInit(noinit bool) *Class {
	c.noinit = noinit
	return c
}

func (c *Class) IsNoInit() bool {
	return c.noinit
}

func (c *Class) SetImmutable(immutable bool) *Class {
	c.immutable = immutable
	return c
}

func (c *Class) IsImmutable() bool {
	return c.immutable
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

func (c *Class) SetNative(native bool) *Class {
	c.native = native
	return c
}

func (c *Class) IsNative() bool {
	return c.native
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
	return c.parent.Get()
}

func (c *Class) Singleton() *SingletonClass {
	return c.singleton.Get()
}

func (c *Class) SetSingleton(singleton *SingletonClass) {
	c.singleton = singleton.ToRef()
}

func getClass(namespace Namespace) Namespace {
	switch namespace := namespace.(type) {
	case *Class:
		return namespace
	case *Generic:
		if _, ok := namespace.Namespace.Get().(*Class); ok {
			return namespace
		}
	case *TemporaryParent:
		return getClass(namespace.Namespace.Get())
	}
	return nil
}

func (c *Class) Superclass() Namespace {
	var currentParent Namespace = c.parent.Get()
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
	c.parent = ToRef(parent)
	c.registerAsChild(parent)

	superclass := c.Superclass()
	currentSingleton := c.singleton.Get()
	if superclass != nil && currentSingleton != nil {
		singleton := superclass.Singleton()
		currentSingleton.parent = CastRef[Namespace](singleton)
		singleton.registerAsChild(currentSingleton)
	}
}

func (c *Class) registerAsChild(parent Namespace) {
	switch parent := parent.(type) {
	case *Class:
		if parent.Children == nil {
			parent.Children = ds.Set[Ref[*Class]]{}
		}
		parent.Children.Add(c.ToRef())
	case *Generic:
		c.registerAsChild(parent.Namespace.Get())
	}
}

func (c *Class) RemoveTemporaryParents() {
	if _, ok := c.parent.Get().(*TemporaryParent); !ok {
		return
	}

	c.parent = 0
	c.singleton.Get().parent = Ref[Namespace](Env.StdSubtypeClass(symbol.C_Class).ID())
}

func NewClass(
	docComment string,
	abstract,
	sealed,
	primitive,
	noinit,
	immutable bool,
	name string,
	parent Namespace,
) *Class {
	class := &Class{
		primitive:     primitive,
		sealed:        sealed,
		abstract:      abstract,
		noinit:        noinit,
		immutable:     immutable,
		native:        Env.Init,
		NamespaceBase: MakeNamespaceBase(docComment, name),
	}
	class.singleton = NewSingletonClass(class, Env.StdSubtypeClass(symbol.C_Class)).ToRef()
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
) *Class {
	class := &Class{
		primitive: primitive,
		abstract:  abstract,
		sealed:    sealed,
		native:    Env.Init,
		NamespaceBase: NamespaceBase{
			docComment: docComment,
			name:       name,
			constants:  consts,
			subtypes:   subtypes,
			methods:    methods,
		},
	}
	class.singleton = NewSingletonClass(class, Env.StdSubtypeClass(symbol.C_Class)).ToRef()
	class.SetParent(parent)
	return class
}

func (c *Class) DefineMethod(docComment string, flags bitfield.BitFlag16, name symbol.Symbol, typeParams []Ref[*TypeParameter], params []Ref[*Parameter], returnType, throwType Type) *Method {
	method := NewMethod(docComment, flags, name, typeParams, params, returnType, throwType, c)
	c.SetMethod(name, method)
	return method
}

func (c *Class) inspect() string {
	return c.name
}

func (c *Class) ToNonLiteral() Type {
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
		native:         c.native,
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
				switch t := val.Type.Get().(type) {
				case *TypeParameter:
					fmt.Printf(" %s:%p", t.InspectSignatureWithColor(), t)
				}
			}
			fmt.Print("]")
			switch n := p.Namespace.Get().(type) {
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
