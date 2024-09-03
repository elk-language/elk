package value

// Represents an Elk Mixin.
type Mixin = Class

// Mixin constructor option function
type MixinOption = func(*Class)

func MixinWithName(name string) MixinOption {
	return func(m *Class) {
		m.ConstantContainer.Name = name
	}
}

func MixinWithClass(class *Class) MixinOption {
	return func(m *Mixin) {
		m.metaClass = class
	}
}

func MixinWithConstants(constants SymbolMap) MixinOption {
	return func(m *Mixin) {
		m.Constants = constants
	}
}

func MixinWithMethods(methods MethodMap) MixinOption {
	return func(m *Mixin) {
		m.Methods = methods
	}
}

func MixinWithParent(parent *Class) MixinOption {
	return func(m *Mixin) {
		m.Parent = parent
	}
}

// Create a new mixin.
func NewMixin() *Mixin {
	m := &Mixin{
		ConstantContainer: ConstantContainer{
			Constants: make(SymbolMap),
		},
		MethodContainer: MethodContainer{
			Methods: make(MethodMap),
		},
		metaClass:         MixinClass,
		instanceVariables: make(SymbolMap),
	}
	m.SetMixin()
	return m
}

// Create a new mixin.
func NewMixinWithOptions(opts ...MixinOption) *Class {
	m := NewMixin()

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Used by the VM, create a new mixin.
func MixinConstructor(class *Class) Value {
	m := &Mixin{
		ConstantContainer: ConstantContainer{
			Constants: make(SymbolMap),
		},
		MethodContainer: MethodContainer{
			Methods: make(MethodMap),
		},
		metaClass:         MixinClass,
		instanceVariables: make(SymbolMap),
	}
	m.SetMixin()

	return m
}

// Create a proxy class that has a pointer to the
// method map of this mixin.
//
// Returns two values, the head and tail proxy classes.
// This is because of the fact that it's possible to include
// one mixin in another, so there is an entire inheritance chain.
func (m *Mixin) CreateProxyClass() *Class {
	proxy := NewClass()
	proxy.SetMixinProxy()
	proxy.Methods = m.Methods
	proxy.Name = m.Name
	proxy.metaClass = m

	return proxy
}

var MixinClass *Class // ::Std::Mixin

func initMixin() {
	MixinClass = NewClass()
	StdModule.AddConstantString("Mixin", MixinClass)
}
