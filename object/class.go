package object

import "github.com/elk-language/elk/bitmap"

const (
	CLASS_SINGLETON_FLAG bitmap.BitFlag8 = 1 << iota // Singleton classes are hidden classes often associated with a single object
	CLASS_ABSTRACT_FLAG                              // Abstract classes can't be instantiated
	CLASS_SEALED_FLAG                                // Sealed classes can't be inherited from
	CLASS_IMMUTABLE_FLAG                             // Immutable classes create frozen instances
	CLASS_FROZEN_FLAG                                // Frozen classes can't define new methods nor constants
)

// Represents an Elk Class.
type Class struct {
	metaClass *Class // Class that this class object is an instance of
	Parent    *Class // Parent/Super class of this class
	ModulelikeObject
	bitmap bitmap.Bitmap8
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

// Create a new class.
func NewClass(opts ...ClassOption) *Class {
	c := &Class{
		metaClass: ClassClass,
		Parent:    ObjectClass,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SimpleSymbolMap),
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Class) IsSingleton() bool {
	return c.bitmap.HasFlag(CLASS_SINGLETON_FLAG)
}

func (c *Class) SetSingleton() {
	c.bitmap.SetFlag(CLASS_SINGLETON_FLAG)
}

func (c *Class) IsAbstract() bool {
	return c.bitmap.HasFlag(CLASS_ABSTRACT_FLAG)
}

func (c *Class) SetAbstract() {
	c.bitmap.SetFlag(CLASS_ABSTRACT_FLAG)
}

func (c *Class) IsSealed() bool {
	return c.bitmap.HasFlag(CLASS_SEALED_FLAG)
}

func (c *Class) SetSealed() {
	c.bitmap.SetFlag(CLASS_SEALED_FLAG)
}

func (c *Class) IsImmutable() bool {
	return c.bitmap.HasFlag(CLASS_IMMUTABLE_FLAG)
}

func (c *Class) SetImmutable() {
	c.bitmap.SetFlag(CLASS_IMMUTABLE_FLAG)
}

func (c *Class) Class() *Class {
	return c.metaClass
}

func (c *Class) IsFrozen() bool {
	return c.bitmap.HasFlag(CLASS_FROZEN_FLAG)
}

func (c *Class) SetFrozen() {
	c.bitmap.SetFlag(CLASS_FROZEN_FLAG)
}

var ClassClass *Class // ::Std::Class
