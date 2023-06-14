package object

import "fmt"

// Struct for embedding, contains fields
// shared by Module, Mixin, Class, Struct
type ModulelikeObject struct {
	Name      string
	Constants SymbolMap
}

// Set the constant with the specified name
// to the given value.
func (m *ModulelikeObject) AddConstant(name string, val Object) {
	switch v := val.(type) {
	case *Module:
		m.setObjectName(&v.ModulelikeObject, name)
	case *Class:
		m.setObjectName(&v.ModulelikeObject, name)
	}
	m.Constants.SetString(name, val)
}

// Set the name of the object when it's assigned to a constant.
func (m *ModulelikeObject) setObjectName(obj *ModulelikeObject, name string) {
	if obj.Name != "" || m.Name == "" {
		return
	}

	if m == &RootModule.ModulelikeObject {
		obj.Name = name
	}

	obj.Name = fmt.Sprintf("%s::%s", m.Name, name)
}
