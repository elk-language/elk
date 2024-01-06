// Package value contains definitions
// of Elk values, classes, structs, modules etc.
package value

import (
	"fmt"
	"maps"

	"github.com/google/go-cmp/cmp"
)

var ObjectClass *Class // ::Std::Object

type Object struct {
	class             *Class
	instanceVariables SymbolMap // Map that stores instance variables of the value
}

func NewObjectComparer(opts cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *Object) bool {
		if x == nil && y == nil {
			return true
		}

		if x == nil || y == nil {
			return false
		}

		return cmp.Equal(x.class, y.class, opts...) &&
			cmp.Equal(x.instanceVariables, y.instanceVariables, opts...)
	})
}

// Class constructor option function
type ObjectOption = func(*Object)

func ObjectWithClass(class *Class) ObjectOption {
	return func(o *Object) {
		o.class = class
	}
}

func ObjectWithInstanceVariables(ivars SymbolMap) ObjectOption {
	return func(o *Object) {
		o.instanceVariables = ivars
	}
}

// Create a new value.
func NewObject(opts ...ObjectOption) *Object {
	o := &Object{
		class:             ObjectClass,
		instanceVariables: make(SymbolMap),
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// Creates a new value.
func ObjectConstructor(class *Class) Value {
	return &Object{
		class:             class,
		instanceVariables: make(SymbolMap),
	}
}

func (o *Object) Copy() Value {
	newInstanceVariables := make(SymbolMap, len(o.instanceVariables))
	maps.Copy(newInstanceVariables, o.instanceVariables)

	newObject := &Object{
		class:             o.class,
		instanceVariables: newInstanceVariables,
	}

	return newObject
}

func (o *Object) InstanceVariables() SymbolMap {
	return o.instanceVariables
}

func (o *Object) Class() *Class {
	if !o.class.IsSingleton() {
		return o.class
	}

	return o.class.Class()
}

func (o *Object) DirectClass() *Class {
	return o.class
}

func (o *Object) SingletonClass() *Class {
	if o.class.IsSingleton() {
		return o.class
	}

	singletonClass := NewClass()
	singletonClass.SetSingleton()
	singletonClass.Parent = o.class
	o.class = singletonClass
	return singletonClass
}

func (o *Object) Inspect() string {
	return fmt.Sprintf(
		"%s%s",
		o.class.PrintableName(),
		o.instanceVariables.Inspect(),
	)
}
