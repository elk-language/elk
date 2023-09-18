package value

import (
	"fmt"

	"github.com/elk-language/elk/bitfield"
)

const (
	CLASS_SINGLETON_FLAG bitfield.BitFlag8 = 1 << iota // Singleton classes are hidden classes often associated with a single value
	CLASS_ABSTRACT_FLAG                                // Abstract classes can't be instantiated
	CLASS_SEALED_FLAG                                  // Sealed classes can't be inherited from
	CLASS_IMMUTABLE_FLAG                               // Immutable classes create frozen instances
	CLASS_FROZEN_FLAG                                  // Frozen classes can't define new methods nor constants
	CLASS_NO_IVARS_FLAG                                // Instances of classes with this flag can't hold instance variables
)

// Function that creates a new instance.
type ConstructorFunc func(class *Class, frozen bool) Value

// Represents an Elk Class.
type Class struct {
	metaClass *Class // Class that this class value is an instance of
	Parent    *Class // Parent/Super class of this class
	ModulelikeObject
	ConstructorFunc   ConstructorFunc
	instanceVariables SimpleSymbolMap
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

func ClassWithImmutable() ClassOption {
	return func(c *Class) {
		c.SetImmutable()
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

func ClassWithConstructor(constructor ConstructorFunc) ClassOption {
	return func(c *Class) {
		c.ConstructorFunc = constructor
	}
}

// Create a new class.
func NewClass(opts ...ClassOption) *Class {
	c := &Class{
		Parent: ObjectClass,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SimpleSymbolMap),
		},
		ConstructorFunc:   ObjectConstructor,
		metaClass:         ClassClass,
		instanceVariables: make(SimpleSymbolMap),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Used by the VM, create a new class.
func ClassConstructor(metaClass *Class, frozen bool) Value {
	c := &Class{
		Parent: ObjectClass,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SimpleSymbolMap),
		},
		ConstructorFunc:   ObjectConstructor,
		metaClass:         metaClass,
		instanceVariables: make(SimpleSymbolMap),
	}
	if frozen {
		c.SetFrozen()
	}

	return c
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

func (c *Class) IsImmutable() bool {
	return c.bitfield.HasFlag(CLASS_IMMUTABLE_FLAG)
}

func (c *Class) SetImmutable() {
	c.bitfield.SetFlag(CLASS_IMMUTABLE_FLAG)
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
	return c.metaClass
}

func (c *Class) IsFrozen() bool {
	return c.bitfield.HasFlag(CLASS_FROZEN_FLAG)
}

func (c *Class) SetFrozen() {
	c.bitfield.SetFlag(CLASS_FROZEN_FLAG)
}

func (c *Class) Inspect() string {
	return fmt.Sprintf("class %s < %s", c.PrintableName(), c.Parent.PrintableName())
}

func (c *Class) InstanceVariables() SimpleSymbolMap {
	return c.instanceVariables
}

var ClassClass *Class // ::Std::Class
