// Package object contains definitions
// of Elk objects, classes, structs, modules etc.
package object

import "fmt"

var PrimitiveObjectClass *Class // ::Std::PrimitiveObject
var ObjectClass *Class          // ::Std::Object

type Object struct {
	class             *Class
	InstanceVariables SimpleSymbolMap // Map that stores instance variables of the object
}

func (o *Object) Class() *Class {
	return o.class
}

func (o *Object) Inspect() string {
	return fmt.Sprintf(
		"%s%s",
		o.class.Name,
		o.InstanceVariables.Inspect(),
	)
}
