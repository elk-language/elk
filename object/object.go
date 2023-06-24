// Package object contains definitions
// of Elk objects, classes, structs, modules etc.
package object

import "fmt"

var PrimitiveObjectClass *Class // ::Std::PrimitiveObject
var ObjectClass *Class          // ::Std::Object

type Object struct {
	class             *Class
	instanceVariables SimpleSymbolMap // Map that stores instance variables of the object
	frozen            bool
}

// Class constructor option function
type ObjectOption = func(*Object)

func ObjectWithClass(class *Class) ObjectOption {
	return func(o *Object) {
		o.class = class
	}
}

// Create a new object.
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
