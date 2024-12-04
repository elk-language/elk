package value

import "fmt"

// Struct for embedding, contains fields
// shared by Module, Mixin, Class, Struct
type ConstantContainer struct {
	Name      string
	Constants SymbolMap
}

// Return a human readable name.
func (m *ConstantContainer) PrintableName() string {
	if m.Name == "" {
		return "<anonymous>"
	}

	return m.Name
}

// Set the constant with the specified name
// to the given value.
func (m *ConstantContainer) AddConstantString(name string, val Value) {
	fullName := m.fullConstantName(name)
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case *Module:
			if v.Name == "" {
				v.Name = fullName
			}
		case *Class:
			if v.Name == "" {
				v.Name = fullName
			}
		case *Interface:
			if v.Name == "" {
				v.Name = fullName
			}
		}
	}

	m.Constants.SetString(name, val)
	RootModule.Constants.Set(ToSymbol(fullName), val)
}

// Set the constant with the specified name
// to the given value.
func (m *ConstantContainer) AddConstant(name Symbol, val Value) {
	fullName := m.fullConstantName(name.String())
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case *Module:
			if v.Name == "" {
				v.Name = fullName
			}
		case *Class:
			if v.Name == "" {
				v.Name = fullName
			}
		case *Interface:
			if v.Name == "" {
				v.Name = fullName
			}
		}
	}

	m.Constants.Set(name, val)
	RootModule.Constants.Set(ToSymbol(fullName), val)
}

func (m *ConstantContainer) fullConstantName(name string) string {
	if m.Name == "Root" {
		return name
	}

	return fmt.Sprintf("%s::%s", m.Name, name)
}
