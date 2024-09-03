package value

import (
	"iter"
	"maps"
	"strings"

	"github.com/elk-language/elk/bitfield"
	"github.com/google/go-cmp/cmp"
)

const (
	CLASS_SINGLETON_FLAG   bitfield.BitFlag8 = 1 << iota // Singleton classes are hidden classes often associated with a single value
	CLASS_ABSTRACT_FLAG                                  // Abstract classes cannot be instantiated
	CLASS_SEALED_FLAG                                    // Sealed classes cannot be inherited from
	CLASS_NO_IVARS_FLAG                                  // Instances of classes with this flag cannot hold instance variables
	CLASS_MIXIN_PROXY_FLAG                               // This class serves as a proxy to an included mixin
	CLASS_MIXIN_FLAG                                     // This class serves as a proxy to an included mixin
)

// Function that creates a new instance.
type ConstructorFunc func(class *Class) Value

// Represents an Elk Class.
type Class struct {
	metaClass *Class // Class that this class value is an instance of
	ConstantContainer
	MethodContainer
	ConstructorFunc   ConstructorFunc
	Flags             bitfield.BitField8
	instanceVariables SymbolMap
}

// Class constructor option function
type ClassOption = func(*Class)

func ClassWithName(name string) ClassOption {
	return func(c *Class) {
		c.ConstantContainer.Name = name
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
func ClassWithMixin() ClassOption {
	return func(c *Class) {
		c.SetMixin()
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
		ConstantContainer: ConstantContainer{
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
		ConstantContainer: ConstantContainer{
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

// Iterates over every parent of the class/mixin (including itself)
func (c *Class) Parents() iter.Seq[*Class] {
	return func(yield func(parent *Class) bool) {
		classes := []*Class{c}
		for i := 0; i < len(classes); i++ {
			parent := classes[i]
			for parent != nil {
				if parent.IsMixinProxy() {
					classes = append(classes, parent.Parent)
					parent = parent.metaClass
					continue
				}
				if !yield(parent) {
					return
				}

				parent = parent.Parent
			}
		}
	}
}

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
	proxy := mixin.CreateProxyClass()
	proxy.Parent = c.Parent
	c.Parent = proxy
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

func (c *Class) IsMixin() bool {
	return c.Flags.HasFlag(CLASS_MIXIN_FLAG)
}

func (c *Class) SetMixin() {
	c.Flags.SetFlag(CLASS_MIXIN_FLAG)
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

func (c *Class) Copy() Value {
	newConstants := make(SymbolMap, len(c.Constants))
	maps.Copy(newConstants, c.Constants)

	newMethods := make(MethodMap, len(c.Methods))
	maps.Copy(newMethods, c.Methods)

	newInstanceVariables := make(SymbolMap, len(c.instanceVariables))
	maps.Copy(newInstanceVariables, c.instanceVariables)

	newClass := &Class{
		ConstantContainer: ConstantContainer{
			Constants: newConstants,
			Name:      c.Name,
		},
		MethodContainer: MethodContainer{
			Methods: newMethods,
			Parent:  c.Parent,
		},
		metaClass:         c.metaClass,
		instanceVariables: newInstanceVariables,
		Flags:             c.Flags,
		ConstructorFunc:   c.ConstructorFunc,
	}

	return newClass
}

func (c *Class) InspectInheritance() string {
	result := new(strings.Builder)
	c.inspectInheritance(result)

	return result.String()
}

func (c *Class) inspectInheritance(buff *strings.Builder) {
	first := true
	parent := c
	for parent != nil {
		if first {
			first = false
		} else {
			buff.WriteString(" < ")
		}
		buff.WriteString(parent.PrintableName())
		if parent.IsMixinProxy() {
			buff.WriteByte('[')
			parent.metaClass.inspectInheritance(buff)
			buff.WriteByte(']')
			parent = parent.Parent
			continue
		}

		parent = parent.Parent
	}
}

func (c *Class) InspectParents() string {
	result := new(strings.Builder)

	first := true
	for parent := range c.Parents() {
		if first {
			first = false
		} else {
			result.WriteString(" < ")
		}
		result.WriteString(parent.PrintableName())
	}

	return result.String()
}

func (c *Class) Inspect() string {
	var result strings.Builder
	if c.IsAbstract() {
		result.WriteString("abstract ")
	}
	if c.IsSealed() {
		result.WriteString("sealed ")
	}
	if c.IsMixin() {
		result.WriteString("mixin ")
	} else {
		result.WriteString("class ")
	}
	result.WriteString(c.PrintableName())

	s := c.Superclass()
	if s != nil {
		result.WriteString(" < ")
		result.WriteString(s.PrintableName())
	}

	return result.String()
}

func (c *Class) InstanceVariables() SymbolMap {
	return c.instanceVariables
}

func NewClassComparer(opts *cmp.Options) cmp.Option {
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
			cmp.Equal(x.instanceVariables, y.instanceVariables, *opts...) &&
			cmp.Equal(x.Constants, y.Constants, *opts...) &&
			cmp.Equal(x.Methods, y.Methods, *opts...) &&
			cmp.Equal(x.Parent, y.Parent, *opts...) &&
			cmp.Equal(x.metaClass, y.metaClass, *opts...)
	})
}

var ClassClass *Class // ::Std::Class
