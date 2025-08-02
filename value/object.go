// Package value contains definitions
// of Elk values, classes, structs, modules etc.
package value

import (
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-cmp/cmp"
)

var ObjectClass *Class // ::Std::Object

type Object struct {
	class             *Class
	instanceVariables InstanceVariables // Slice that stores instance variables of the value
}

func NewObjectComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *Object) bool {
		if x == nil && y == nil {
			return true
		}

		if x == nil || y == nil {
			return false
		}

		if x.class == y.class {
			return true
		}
		if x.class != nil && x.class.Name == y.class.Name {
			return true
		}

		return cmp.Equal(x.class, y.class, *opts...) &&
			cmp.Equal(x.instanceVariables, y.instanceVariables, *opts...)
	})
}

// Class constructor option function
type ObjectOption = func(*Object)

func ObjectWithClass(class *Class) ObjectOption {
	return func(o *Object) {
		o.class = class
		o.instanceVariables = make([]Value, len(class.IvarIndices))
	}
}

func ObjectWithInstanceVariables(ivars []Value) ObjectOption {
	return func(o *Object) {
		o.instanceVariables = ivars
	}
}

func ObjectWithInstanceVariablesByName(ivars SymbolMap) ObjectOption {
	return func(o *Object) {
		for key, value := range ivars {
			SetInstanceVariableByName(Ref(o), key, value)
		}
	}
}

// Create a new value.
func NewObject(opts ...ObjectOption) *Object {
	o := &Object{
		class: ObjectClass,
	}

	for _, opt := range opts {
		opt(o)
	}

	if o.instanceVariables == nil {
		o.instanceVariables = make([]Value, len(ObjectClass.IvarIndices))
	}

	return o
}

// Creates a new value.
func ObjectConstructor(class *Class) Value {
	return Ref(&Object{
		class:             class,
		instanceVariables: make([]Value, len(class.IvarIndices)),
	})
}

func (o *Object) Copy() Reference {
	newObject := &Object{
		class:             o.class,
		instanceVariables: slices.Clone(o.instanceVariables),
	}

	return newObject
}

func (o *Object) InstanceVariables() *InstanceVariables {
	return &o.instanceVariables
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
	var buff strings.Builder
	fmt.Fprintf(&buff, "%s{&: %p", o.class.PrintableName(), o)

	for ivarIndex, val := range o.instanceVariables {
		if val.IsUndefined() {
			continue
		}

		ivarName, ok := o.class.GetIvarName(ivarIndex)
		if !ok {
			panic(fmt.Sprintf("could not find ivar for class `%s` with index `%d`", o.class.Inspect(), ivarIndex))
		}

		fmt.Fprintf(&buff, ", %s: %s", ivarName.InspectContent(), val.Inspect())
	}

	buff.WriteRune('}')
	return buff.String()
}
