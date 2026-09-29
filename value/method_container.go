package value

type MethodContainer struct {
	methods MethodMap
	parent  *Class
}

func MakeMethodContainer(methods MethodMap, parent *Class) MethodContainer {
	return MethodContainer{
		methods: methods,
		parent:  parent,
	}
}

func NewMethodContainer(methods MethodMap, parent *Class) *MethodContainer {
	return &MethodContainer{
		methods: methods,
		parent:  parent,
	}
}

// Get the superclass (skipping any mixin proxies)
func (m *MethodContainer) Superclass() *Class {
	currentClass := m.parent
	for {
		if currentClass == nil || !currentClass.IsMixinProxy() {
			return currentClass
		}

		currentClass = currentClass.parent
	}
}

// Get the superclass (skipping any mixin proxies)
func (m *MethodContainer) Methods() MethodMap {
	return m.methods
}

// Search for a method with the given name in
// this container and its ancestors.
func (m *MethodContainer) LookupMethod(name Symbol) Method {
	if method, ok := m.methods[name]; ok {
		return method
	}

	for currentClass := range m.parent.Parents() {
		if method, ok := currentClass.methods[name]; ok {
			return method
		}
	}

	return nil
}

// Attaches the given method under the given name.
func (m *MethodContainer) AttachMethod(name Symbol, method Method) {
	m.methods[name] = method
}

// Define an alternative name for an existing method.
func (m *MethodContainer) DefineAlias(newMethodName, oldMethodName Symbol) {
	method := m.LookupMethod(oldMethodName)
	m.AttachMethod(newMethodName, method)
}

// Define an alternative name for an existing method.
func (m *MethodContainer) DefineAliasString(newMethodName, oldMethodName string) {
	m.DefineAlias(ToSymbol(newMethodName), ToSymbol(oldMethodName))
}
