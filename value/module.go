package value

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

// Represents an Elk Module.
type Module struct {
	class             *Class // The class that this module is an instance of
	frozen            bool   // Frozen module can't define new constants
	instanceVariables SimpleSymbolMap
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

func ModuleWithConstants(constants SimpleSymbolMap) ModuleOption {
	return func(m *Module) {
		m.Constants = constants
	}
}

// Create a new module.
func NewModule() *Module {
	return &Module{
		class: ModuleClass,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SimpleSymbolMap),
			Methods:   make(MethodMap),
		},
		instanceVariables: make(SimpleSymbolMap),
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
func ModuleConstructor(class *Class, frozen bool) Value {
	return &Module{
		class: class,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SimpleSymbolMap),
		},
		instanceVariables: make(SimpleSymbolMap),
		frozen:            frozen,
	}
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

func (m *Module) SingletonClass() *Class {
	if m.class.IsSingleton() {
		return m.class
	}

	singletonClass := NewClass()
	singletonClass.SetSingleton()
	singletonClass.Parent = m.class
	m.class = singletonClass
	return singletonClass
}

func (m *Module) IsFrozen() bool {
	return m.frozen
}

func (m *Module) SetFrozen() {
	m.frozen = true
}

func (m *Module) Inspect() string {
	return fmt.Sprintf("module %s", m.PrintableName())
}

func (m *Module) InstanceVariables() SimpleSymbolMap {
	return m.instanceVariables
}

var ModuleComparer cmp.Option

func initModuleComparer() {
	ModuleComparer = cmp.Comparer(func(x, y *Module) bool {
		if x == y {
			return true
		}

		return x.Name == y.Name &&
			cmp.Equal(x.instanceVariables, y.instanceVariables, ValueComparerOptions...) &&
			cmp.Equal(x.Constants, y.Constants, ValueComparerOptions...) &&
			cmp.Equal(x.Methods, y.Methods, ValueComparerOptions...) &&
			cmp.Equal(x.class, y.class, ValueComparerOptions...)
	})
}

var ModuleClass *Class // ::Std::Module
var RootModule *Module // ::Root
var StdModule *Module  // ::Std
