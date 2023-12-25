package value

import (
	"strings"

	"github.com/elk-language/elk/bitfield"
	"github.com/google/go-cmp/cmp"
)

const (
	CLASS_SINGLETON_FLAG   bitfield.BitFlag8 = 1 << iota // Singleton classes are hidden classes often associated with a single value
	CLASS_ABSTRACT_FLAG                                  // Abstract classes can't be instantiated
	CLASS_SEALED_FLAG                                    // Sealed classes can't be inherited from
	CLASS_NO_IVARS_FLAG                                  // Instances of classes with this flag can't hold instance variables
	CLASS_MIXIN_PROXY_FLAG                               // This class serves as a proxy to an included mixin
)

// Function that creates a new instance.
type ConstructorFunc func(class *Class) Value

// Represents an Elk Class.
type Class struct {
	metaClass *Class // Class that this class value is an instance of
	ModulelikeObject
	MethodContainer
	ConstructorFunc   ConstructorFunc
	Flags             bitfield.Bitfield8
	instanceVariables SymbolMap
}

// Class constructor option function
type ClassOption = func(*Class)

func ClassWithName(name string) ClassOption {
	return func(c *Class) {
		c.ModulelikeObject.Name = name
	}
}

func ClassWithParent(parent *Class) ClassOption {
	return func(c *Class) {
		c.Parent = parent
	}
}

func ClassWithMetaClass(metaClass *Class) ClassOption {
	return func(c *Class) {
		c.metaClass = metaClass
	}
}

func ClassWithConstants(constants SymbolMap) ClassOption {
	return func(c *Class) {
		c.Constants = constants
	}
}

func ClassWithMethods(methods MethodMap) ClassOption {
	return func(c *Class) {
		c.Methods = methods
	}
}

func ClassWithAbstract() ClassOption {
	return func(c *Class) {
		c.SetAbstract()
	}
}

func ClassWithSingleton() ClassOption {
	return func(c *Class) {
		c.SetSingleton()
	}
}

func ClassWithSealed() ClassOption {
	return func(c *Class) {
		c.SetSealed()
	}
}

func ClassWithNoInstanceVariables() ClassOption {
	return func(c *Class) {
		c.SetNoInstanceVariables()
	}
}

func ClassWithMixinProxy() ClassOption {
	return func(c *Class) {
		c.SetMixinProxy()
	}
}

func ClassWithConstructor(constructor ConstructorFunc) ClassOption {
	return func(c *Class) {
		c.ConstructorFunc = constructor
	}
}

// Create a new class.
func NewClass() *Class {
	return &Class{
		ModulelikeObject: ModulelikeObject{
			Constants: make(SymbolMap),
		},
		MethodContainer: MethodContainer{
			Parent:  ObjectClass,
			Methods: make(MethodMap),
		},
		ConstructorFunc:   ObjectConstructor,
		metaClass:         ClassClass,
		instanceVariables: make(SymbolMap),
	}
}

func NewSingletonClass(originalClass *Class, originalName string) *Class {
	singletonClass := NewClass()
	singletonClass.SetSingleton()
	singletonClass.Parent = originalClass
	singletonClass.SetSingletonName(originalName)
	return singletonClass
}

// Create a new class.
func NewClassWithOptions(opts ...ClassOption) *Class {
	c := NewClass()

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Used by the VM, create a new class.
func ClassConstructor(metaClass *Class) Value {
	c := &Class{
		ModulelikeObject: ModulelikeObject{
			Constants: make(SymbolMap),
		},
		MethodContainer: MethodContainer{
			Parent: ObjectClass,
		},
		ConstructorFunc:   ObjectConstructor,
		metaClass:         metaClass,
		instanceVariables: make(SymbolMap),
	}

	return c
}

var docSymbol = SymbolTable.Add("doc")

func (c *Class) SetDoc(doc String) {
	c.Constants.Set(docSymbol, doc)
}

func (c *Class) Doc() Value {
	return c.Constants.Get(docSymbol)
}

// Create a new instance of the class without initialising it.
func (c *Class) CreateInstance() Value {
	return c.ConstructorFunc(c)
}

// Include the passed in mixin in this class.
func (c *Class) IncludeMixin(mixin *Mixin) {
	headProxy, tailProxy := mixin.CreateProxyClass()
	tailProxy.Parent = c.Parent
	c.Parent = headProxy
}

func (c *Class) IsSingleton() bool {
	return c.Flags.HasFlag(CLASS_SINGLETON_FLAG)
}

func (c *Class) SetSingleton() {
	c.Flags.SetFlag(CLASS_SINGLETON_FLAG)
}

func (c *Class) IsAbstract() bool {
	return c.Flags.HasFlag(CLASS_ABSTRACT_FLAG)
}

func (c *Class) SetAbstract() {
	c.Flags.SetFlag(CLASS_ABSTRACT_FLAG)
}

func (c *Class) IsSealed() bool {
	return c.Flags.HasFlag(CLASS_SEALED_FLAG)
}

func (c *Class) SetSealed() {
	c.Flags.SetFlag(CLASS_SEALED_FLAG)
}

func (c *Class) IsMixinProxy() bool {
	return c.Flags.HasFlag(CLASS_MIXIN_PROXY_FLAG)
}

func (c *Class) SetMixinProxy() {
	c.Flags.SetFlag(CLASS_MIXIN_PROXY_FLAG)
}

// Whether instances of this class can hold
// instance variables.
func (c *Class) HasNoInstanceVariables() bool {
	return c.Flags.HasFlag(CLASS_NO_IVARS_FLAG)
}

func (c *Class) SetNoInstanceVariables() {
	c.Flags.SetFlag(CLASS_NO_IVARS_FLAG)
}

func (c *Class) Class() *Class {
	if !c.metaClass.IsSingleton() {
		return c.metaClass
	}

	return c.metaClass.Class()
}

func (c *Class) DirectClass() *Class {
	return c.metaClass
}

func (c *Class) SetDirectClass(metaClass *Class) {
	c.metaClass = metaClass
}

func (c *Class) SetSingletonName(name string) {
	if name != "" {
		c.Name = "&" + name
	}
}

func (c *Class) SingletonClass() *Class {
	if c.metaClass.IsSingleton() {
		return c.metaClass
	}

	singletonClass := NewSingletonClass(c.metaClass, c.Name)
	c.metaClass = singletonClass
	return singletonClass
}

func (c *Class) Inspect() string {
	var result strings.Builder
	if c.IsAbstract() {
		result.WriteString("abstract ")
	}
	if c.IsSealed() {
		result.WriteString("sealed ")
	}
	result.WriteString("class ")
	result.WriteString(c.PrintableName())

	if c.Parent != nil {
		result.WriteString(" < ")
		result.WriteString(c.Parent.PrintableName())
	}

	return result.String()
}

func (c *Class) InstanceVariables() SymbolMap {
	return c.instanceVariables
}

func NewClassComparer(opts cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *Class) bool {
		if x == y {
			return true
		}

		if x == nil || y == nil {
			return false
		}

		if x == ClassClass || y == ClassClass {
			return false
		}

		return x.Flags == y.Flags &&
			x.Name == y.Name &&
			cmp.Equal(x.instanceVariables, y.instanceVariables, opts...) &&
			cmp.Equal(x.Constants, y.Constants, opts...) &&
			cmp.Equal(x.Methods, y.Methods, opts...) &&
			cmp.Equal(x.Parent, y.Parent, opts...) &&
			cmp.Equal(x.metaClass, y.metaClass, opts...)
	})
}

var ClassClass *Class // ::Std::Class
