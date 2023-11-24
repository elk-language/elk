package value

import (
	"fmt"

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
	Parent    *Class // Parent/Super class of this class
	ModulelikeObject
	Methods           MethodMap
	ConstructorFunc   ConstructorFunc
	instanceVariables SymbolMap
	bitfield          bitfield.Bitfield8
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
		Parent: ObjectClass,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SymbolMap),
		},
		Methods:           make(MethodMap),
		ConstructorFunc:   ObjectConstructor,
		metaClass:         ClassClass,
		instanceVariables: make(SymbolMap),
	}
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
		Parent: ObjectClass,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SymbolMap),
		},
		ConstructorFunc:   ObjectConstructor,
		metaClass:         metaClass,
		instanceVariables: make(SymbolMap),
	}

	return c
}

// Include the passed in mixin in this class.
func (c *Class) IncludeMixin(mixin *Mixin) {
	headProxy, tailProxy := mixin.CreateProxyClass()
	tailProxy.Parent = c.Parent
	c.Parent = headProxy
}

// Search for a method with the given name in this class
// and its ancestors.
func (c *Class) LookupMethod(name Symbol) Method {
	currentClass := c
	for currentClass != nil {
		if method, ok := currentClass.Methods[name]; ok {
			return method
		}
		currentClass = currentClass.Parent
	}

	return nil
}

func (c *Class) IsSingleton() bool {
	return c.bitfield.HasFlag(CLASS_SINGLETON_FLAG)
}

func (c *Class) SetSingleton() {
	c.bitfield.SetFlag(CLASS_SINGLETON_FLAG)
}

func (c *Class) IsAbstract() bool {
	return c.bitfield.HasFlag(CLASS_ABSTRACT_FLAG)
}

func (c *Class) SetAbstract() {
	c.bitfield.SetFlag(CLASS_ABSTRACT_FLAG)
}

func (c *Class) IsSealed() bool {
	return c.bitfield.HasFlag(CLASS_SEALED_FLAG)
}

func (c *Class) SetSealed() {
	c.bitfield.SetFlag(CLASS_SEALED_FLAG)
}

func (c *Class) IsMixinProxy() bool {
	return c.bitfield.HasFlag(CLASS_MIXIN_PROXY_FLAG)
}

func (c *Class) SetMixinProxy() {
	c.bitfield.SetFlag(CLASS_MIXIN_PROXY_FLAG)
}

// Whether instances of this class can hold
// instance variables.
func (c *Class) HasNoInstanceVariables() bool {
	return c.bitfield.HasFlag(CLASS_NO_IVARS_FLAG)
}

func (c *Class) SetNoInstanceVariables() {
	c.bitfield.SetFlag(CLASS_NO_IVARS_FLAG)
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

func (c *Class) SingletonClass() *Class {
	if c.metaClass.IsSingleton() {
		return c.metaClass
	}

	singletonClass := NewClass()
	singletonClass.SetSingleton()
	singletonClass.Parent = c.metaClass
	c.metaClass = singletonClass
	return singletonClass
}

func (c *Class) Inspect() string {
	if c.Parent == nil {
		return fmt.Sprintf("class %s", c.PrintableName())
	}
	return fmt.Sprintf("class %s < %s", c.PrintableName(), c.Parent.PrintableName())
}

func (c *Class) InstanceVariables() SymbolMap {
	return c.instanceVariables
}

// Define an alternative name for an existing method.
func (c *Class) DefineAlias(newMethodName, oldMethodName Symbol) (err *Error) {
	method := c.LookupMethod(oldMethodName)
	if method == nil {
		return NewCantCreateAnAliasForNonexistentMethod(oldMethodName.ToString())
	}

	newMethod := c.LookupMethod(newMethodName)
	if newMethod != nil && newMethod.IsFrozen() {
		return NewCantOverrideAFrozenMethod(newMethodName.ToString())
	}

	c.Methods[newMethodName] = method

	return nil
}

// Define an alternative name for an existing method.
func (c *Class) DefineAliasString(newMethodName, oldMethodName string) (err *Error) {
	return c.DefineAlias(ToSymbol(newMethodName), ToSymbol(oldMethodName))
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

		return x.bitfield == y.bitfield &&
			x.Name == y.Name &&
			cmp.Equal(x.instanceVariables, y.instanceVariables, opts...) &&
			cmp.Equal(x.Constants, y.Constants, opts...) &&
			cmp.Equal(x.Methods, y.Methods, opts...) &&
			cmp.Equal(x.Parent, y.Parent, opts...) &&
			cmp.Equal(x.metaClass, y.metaClass, opts...)
	})
}

var ClassClass *Class // ::Std::Class
