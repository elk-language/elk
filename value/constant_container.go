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
	switch v := val.(type) {
	case *Module:
		m.setObjectName(&v.ConstantContainer, name)
	case *Class:
		m.setObjectName(&v.ConstantContainer, name)
	case *Mixin:
		m.setObjectName(&v.ConstantContainer, name)
	}
	m.Constants.SetString(name, val)
}

// Set the constant with the specified name
// to the given value.
func (m *ConstantContainer) AddConstant(name Symbol, val Value) {
	switch v := val.(type) {
	case *Module:
		m.setObjectName(&v.ConstantContainer, string(name.ToString()))
	case *Class:
		m.setObjectName(&v.ConstantContainer, string(name.ToString()))
	case *Mixin:
		m.setObjectName(&v.ConstantContainer, string(name.ToString()))
	}
	m.Constants.Set(name, val)
}

// Set the name of the value when it's assigned to a constant.
func (m *ConstantContainer) setObjectName(obj *ConstantContainer, name string) {
	if obj.Name != "" || m.Name == "" {
		return
	}

	if m == &RootModule.ConstantContainer {
		obj.Name = name
		return
	}

	obj.Name = fmt.Sprintf("%s::%s", m.Name, name)
}
