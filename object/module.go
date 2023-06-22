package object

import "fmt"

// Represents an Elk Module.
type Module struct {
	class  *Class // The class that this module is an instance of
	frozen bool   // Frozen module can't define new constants
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

// Create a new module.
func NewModule(opts ...ModuleOption) *Module {
	m := &Module{
		class: ModuleClass,
		ModulelikeObject: ModulelikeObject{
			Constants: make(SimpleSymbolMap),
		},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Module) Class() *Class {
	return m.class
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

var ModuleClass *Class // ::Std::Module
var RootModule *Module // ::Root
var StdModule *Module  // ::Std
