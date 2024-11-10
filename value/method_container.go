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
	m.Methods[name] = method
	return nil
}

// Define an alternative name for an existing method.
func (m *MethodContainer) DefineAlias(newMethodName, oldMethodName Symbol) *Error {
	method := m.LookupMethod(oldMethodName)
	return m.AttachMethod(newMethodName, method)
}

// Define an alternative name for an existing method.
func (m *MethodContainer) DefineAliasString(newMethodName, oldMethodName string) *Error {
	return m.DefineAlias(ToSymbol(newMethodName), ToSymbol(oldMethodName))
}
