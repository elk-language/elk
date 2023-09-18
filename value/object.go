// Package value contains definitions
// of Elk values, classes, structs, modules etc.
package value

import "fmt"

var PrimitiveObjectClass *Class // ::Std::PrimitiveObject
var ObjectClass *Class          // ::Std::Object

type Object struct {
	class             *Class
	instanceVariables SimpleSymbolMap // Map that stores instance variables of the value
	frozen            bool
}

// Class constructor option function
type ObjectOption = func(*Object)

func ObjectWithClass(class *Class) ObjectOption {
	return func(o *Object) {
		o.class = class
	}
}

func ObjectWithInstanceVariables(ivars SimpleSymbolMap) ObjectOption {
	return func(o *Object) {
		o.instanceVariables = ivars
	}
}

// Create a new value.
func NewObject(opts ...ObjectOption) *Object {
	o := &Object{
		class:             ObjectClass,
		instanceVariables: make(SimpleSymbolMap),
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// Creates a new value.
func ObjectConstructor(class *Class, frozen bool) Value {
	return &Object{
		class:             class,
		frozen:            frozen,
		instanceVariables: make(SimpleSymbolMap),
	}
}

func (o *Object) InstanceVariables() SimpleSymbolMap {
	return o.instanceVariables
}

func (o *Object) Class() *Class {
	return o.class
}

func (o *Object) Inspect() string {
	return fmt.Sprintf(
		"%s%s",
		o.class.PrintableName(),
		o.instanceVariables.Inspect(),
	)
}

func (o *Object) IsFrozen() bool {
	return o.frozen
}

func (o *Object) SetFrozen() {
	o.frozen = true
}
