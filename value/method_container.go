package value

type MethodContainer struct {
	Methods MethodMap
	Parent  *Class
}

func (m *MethodContainer) CanOverride(name Symbol) bool {
	oldMethod := m.LookupMethod(name)
	return oldMethod == nil || !oldMethod.IsFrozen()
}

// Search for a method with the given name in
// this container and its ancestors.
func (m *MethodContainer) LookupMethod(name Symbol) Method {
	if method := m.Methods[name]; method != nil {
		return method
	}

	currentClass := m.Parent
	for currentClass != nil {
		if method, ok := currentClass.Methods[name]; ok {
			return method
		}
		currentClass = currentClass.Parent
	}

	return nil
}

func (m *MethodContainer) AttachMethod(name Symbol, method Method) *Error {
	if !m.CanOverride(name) {
		return NewCantOverrideAFrozenMethod(name.ToString())
	}

	m.Methods[name] = method
	return nil
}
