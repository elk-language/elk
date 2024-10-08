package value

type MethodContainer struct {
	Methods MethodMap
	Parent  *Class
}

// Get the superclass (skipping any mixin proxies)
func (m *MethodContainer) Superclass() *Class {
	currentClass := m.Parent
	for {
		if currentClass == nil || !currentClass.IsMixinProxy() {
			return currentClass
		}

		currentClass = currentClass.Parent
	}
}

func (m *MethodContainer) CanOverride(name Symbol) bool {
	oldMethod := m.LookupMethod(name)
	return oldMethod == nil || !oldMethod.IsSealed()
}

// Search for a method with the given name in
// this container and its ancestors.
func (m *MethodContainer) LookupMethod(name Symbol) Method {
	if method, ok := m.Methods[name]; ok {
		return method
	}

	for currentClass := range m.Parent.Parents() {
		if method, ok := currentClass.Methods[name]; ok {
			return method
		}
	}

	return nil
}

// Attaches the given method under the given name.
func (m *MethodContainer) AttachMethod(name Symbol, method Method) *Error {
	if !m.CanOverride(name) {
		return NewCantOverrideASealedMethod(string(name.ToString()))
	}

	m.Methods[name] = method
	return nil
}

// Define an alternative name for an existing method.
func (m *MethodContainer) DefineAlias(newMethodName, oldMethodName Symbol) *Error {
	method := m.LookupMethod(oldMethodName)
	if method == nil {
		return NewCantCreateAnAliasForNonexistentMethod(string(oldMethodName.ToString()))
	}

	return m.AttachMethod(newMethodName, method)
}

// Define an alternative name for an existing method.
func (m *MethodContainer) DefineAliasString(newMethodName, oldMethodName string) *Error {
	return m.DefineAlias(ToSymbol(newMethodName), ToSymbol(oldMethodName))
}
