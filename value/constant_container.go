package value

import "fmt"

// Struct for embedding, contains fields
// shared by Module, Mixin, Class, Struct
type ConstantContainer struct {
	name      string
	constants SymbolMap
}

func (m *ConstantContainer) Name() string {
	return m.name
}

func (m *ConstantContainer) Constants() SymbolMap {
	return m.constants
}

// Return a human readable name.
func (m *ConstantContainer) PrintableName() string {
	if m.name == "" {
		return "<anonymous>"
	}

	return m.name
}

// Set the constant with the specified name
// to the given value.
func (m *ConstantContainer) AddConstantString(name string, val Value) {
	fullName := m.fullConstantName(name)
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case *Module:
			if v.name == "" {
				v.name = fullName
			}
		case *Class:
			if v.name == "" {
				v.name = fullName
			}
		case *Interface:
			if v.name == "" {
				v.name = fullName
			}
		}
	}

	m.constants.SetString(name, val)
	RootModule.constants.Set(ToSymbol(fullName), val)
}

// Set the constant with the specified name
// to the given value.
func (m *ConstantContainer) AddConstant(name Symbol, val Value) {
	fullName := m.fullConstantName(name.String())
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case *Module:
			if v.name == "" {
				v.name = fullName
			}
		case *Class:
			if v.name == "" {
				v.name = fullName
			}
		case *Interface:
			if v.name == "" {
				v.name = fullName
			}
		}
	}

	m.constants.Set(name, val)
	RootModule.constants.Set(ToSymbol(fullName), val)
}

func (m *ConstantContainer) fullConstantName(name string) string {
	if m.name == "Root" {
		return name
	}

	return fmt.Sprintf("%s::%s", m.name, name)
}
