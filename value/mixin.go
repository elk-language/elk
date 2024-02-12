package value

import (
	"fmt"
	"maps"

	"github.com/google/go-cmp/cmp"
)

// Represents an Elk Mixin.
type Mixin struct {
	class *Class // Class that this mixin value is an instance of
	ModulelikeObject
	MethodContainer
	instanceVariables SymbolMap
}

// Mixin constructor option function
type MixinOption = func(*Mixin)

func MixinWithName(name string) MixinOption {
	return func(m *Mixin) {
		m.ModulelikeObject.Name = name
	}
}

func MixinWithClass(class *Class) MixinOption {
	return func(m *Mixin) {
		m.class = class
	}
}

func MixinWithConstants(constants SymbolMap) MixinOption {
	return func(m *Mixin) {
		m.Constants = constants
	}
}

func MixinWithMethods(methods MethodMap) MixinOption {
	return func(m *Mixin) {
		m.Methods = methods
	}
}

func MixinWithParent(parent *Class) MixinOption {
	return func(m *Mixin) {
		m.Parent = parent
	}
}

// Create a new mixin.
func NewMixin() *Mixin {
	return &Mixin{
		ModulelikeObject: ModulelikeObject{
			Constants: make(SymbolMap),
		},
		MethodContainer: MethodContainer{
			Methods: make(MethodMap),
		},
		class:             MixinClass,
		instanceVariables: make(SymbolMap),
	}
}

// Create a new mixin.
func NewMixinWithOptions(opts ...MixinOption) *Mixin {
	m := NewMixin()

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Used by the VM, create a new mixin.
func MixinConstructor(class *Class) Value {
	m := &Mixin{
		ModulelikeObject: ModulelikeObject{
			Constants: make(SymbolMap),
		},
		MethodContainer: MethodContainer{
			Methods: make(MethodMap),
		},
		class:             class,
		instanceVariables: make(SymbolMap),
	}

	return m
}

// Include the passed in mixin in this mixin.
func (m *Mixin) IncludeMixin(otherMixin *Mixin) {
	headProxy, tailProxy := otherMixin.CreateProxyClass()
	tailProxy.Parent = m.Parent
	m.Parent = headProxy
}

// Create a proxy class that has a pointer to the
// method map of this mixin.
//
// Returns two values, the head and tail proxy classes.
// This is because of the fact that it's possible to include
// one mixin in another, so there is an entire inheritance chain.
func (m *Mixin) CreateProxyClass() (head, tail *Class) {
	headProxy := NewClass()
	headProxy.SetMixinProxy()
	headProxy.Methods = m.Methods
	headProxy.Name = m.Name

	tailProxy := headProxy
	baseProxy := m.Parent
	for baseProxy != nil {
		proxyCopy := NewClass()
		proxyCopy.SetMixinProxy()
		proxyCopy.Methods = baseProxy.Methods
		proxyCopy.Name = baseProxy.Name
		tailProxy.Parent = proxyCopy

		baseProxy = baseProxy.Parent
		tailProxy = proxyCopy
	}

	return headProxy, tailProxy
}

func (m *Mixin) Class() *Class {
	if !m.class.IsSingleton() {
		return m.class
	}

	return m.class.Class()
}

func (m *Mixin) DirectClass() *Class {
	return m.class
}

func (m *Mixin) SingletonClass() *Class {
	if m.class.IsSingleton() {
		return m.class
	}

	singletonClass := NewSingletonClass(m.class, m.Name)
	m.class = singletonClass
	return singletonClass
}

func (m *Mixin) Copy() Value {
	newConstants := make(SymbolMap, len(m.Constants))
	maps.Copy(newConstants, m.Constants)

	newMethods := make(MethodMap, len(m.Methods))
	maps.Copy(newMethods, m.Methods)

	newInstanceVariables := make(SymbolMap, len(m.instanceVariables))
	maps.Copy(newInstanceVariables, m.instanceVariables)

	newMixin := &Mixin{
		ModulelikeObject: ModulelikeObject{
			Constants: newConstants,
			Name:      m.Name,
		},
		MethodContainer: MethodContainer{
			Methods: newMethods,
			Parent:  m.Parent,
		},
		class:             m.class,
		instanceVariables: newInstanceVariables,
	}

	return newMixin
}

func (m *Mixin) Inspect() string {
	return fmt.Sprintf("mixin %s", m.PrintableName())
}

func (m *Mixin) InstanceVariables() SymbolMap {
	return m.instanceVariables
}

func NewMixinComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *Mixin) bool {
		if x == y {
			return true
		}

		return x.Name == y.Name &&
			cmp.Equal(x.instanceVariables, y.instanceVariables, *opts...) &&
			cmp.Equal(x.Constants, y.Constants, *opts...) &&
			cmp.Equal(x.Methods, y.Methods, *opts...) &&
			cmp.Equal(x.class, y.class, *opts...) &&
			cmp.Equal(x.Parent, y.Parent, *opts...)
	})
}

var MixinClass *Class // ::Std::Mixin

func initMixin() {
	MixinClass = NewClass()
	StdModule.AddConstantString("Mixin", MixinClass)
}
