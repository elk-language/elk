package value

import (
	"fmt"
	"maps"

	"github.com/google/go-cmp/cmp"
)

// Represents an Elk Module.
type Module struct {
	class             *Class // The class that this module is an instance of
	instanceVariables SymbolMap
	ModulelikeObject
}

// Module constructor option function.
type ModuleOption func(*Module)

func ModuleWithName(name string) ModuleOption {
	return func(m *Module) {
		m.ModulelikeObject.Name = name
	}
}

func ModuleWithClass(class *Class) ModuleOption {
	return func(m *Module) {
		m.class = class
	}
}

func ModuleWithConstants(constants SymbolMap) ModuleOption {
	return func(m *Module) {
		m.Constants = constants
	}
}

// Create a new module.
func NewModule() *Module {
	return &Module{
		class: ModuleClass,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SymbolMap),
		},
		instanceVariables: make(SymbolMap),
	}
}

// Create a new module.
func NewModuleWithOptions(opts ...ModuleOption) *Module {
	m := NewModule()

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Used by the VM, create a new module value.
func ModuleConstructor(class *Class) Value {
	return &Module{
		class: class,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SymbolMap),
		},
		instanceVariables: make(SymbolMap),
	}
}

func (m *Module) Copy() Value {
	newConstants := make(SymbolMap, len(m.Constants))
	maps.Copy(newConstants, m.Constants)

	newInstanceVariables := make(SymbolMap, len(m.instanceVariables))
	maps.Copy(newInstanceVariables, m.instanceVariables)

	newModule := &Module{
		ModulelikeObject: ModulelikeObject{
			Constants: newConstants,
			Name:      m.Name,
		},
		class:             m.class,
		instanceVariables: newInstanceVariables,
	}

	return newModule
}

func (m *Module) Class() *Class {
	if !m.class.IsSingleton() {
		return m.class
	}

	return m.class.Class()
}

func (m *Module) DirectClass() *Class {
	return m.class
}

func (m *Module) SetDirectClass(class *Class) {
	m.class = class
}

func (m *Module) SingletonClass() *Class {
	if m.class.IsSingleton() {
		return m.class
	}

	singletonClass := NewSingletonClass(m.class, m.Name)
	m.class = singletonClass
	return singletonClass
}

func (m *Module) Inspect() string {
	return fmt.Sprintf("module %s", m.PrintableName())
}

func (m *Module) InstanceVariables() SymbolMap {
	return m.instanceVariables
}

func NewModuleComparer(opts cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *Module) bool {
		if x == y {
			return true
		}

		return x.Name == y.Name &&
			cmp.Equal(x.instanceVariables, y.instanceVariables, opts...) &&
			cmp.Equal(x.Constants, y.Constants, opts...) &&
			cmp.Equal(x.class, y.class, opts...)
	})
}

var ModuleClass *Class // ::Std::Module
var RootModule *Module // ::Root
var StdModule *Module  // ::Std
