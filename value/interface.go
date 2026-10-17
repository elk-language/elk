package value

import (
	"fmt"
	"maps"

	"github.com/google/go-cmp/cmp"
)

// Represents an Elk interface.
type Interface struct {
	class *Class // The class that this interface is an instance of
	ConstantContainer
}

var _ Namespace = &Interface{}

// Interface constructor option function
type InterfaceOption = func(*Interface)

func InterfaceWithName(name string) InterfaceOption {
	return func(i *Interface) {
		i.ConstantContainer.name = name
	}
}

// Create a new module.
func NewInterface() *Interface {
	return &Interface{
		class: InterfaceClass,
		ConstantContainer: ConstantContainer{
			constants: make(SymbolMap),
		},
	}
}

// Create a new class.
func NewInterfaceWithOptions(opts ...InterfaceOption) *Interface {
	i := NewInterface()

	for _, opt := range opts {
		opt(i)
	}

	return i
}

// Used by the VM, create a new interface value.
func InterfaceConstructor(class *Class) Value {
	return Ref(NewInterface())
}

func (i *Interface) Copy() Reference {
	newConstants := make(SymbolMap, len(i.constants))
	maps.Copy(newConstants, i.constants)

	newInterface := &Interface{
		ConstantContainer: ConstantContainer{
			constants: newConstants,
			name:      i.name,
		},
	}

	return newInterface
}

func (i *Interface) ToValue() Value {
	return Ref(i)
}

func (i *Interface) Class() *Class {
	return InterfaceClass
}

func (i *Interface) DirectClass() *Class {
	return i.class
}

func (i *Interface) SingletonClass() *Class {
	if i.class.IsSingleton() {
		return i.class
	}

	singletonClass := NewSingletonClass(i.class, i.name)
	i.class = singletonClass
	return singletonClass
}

func (i *Interface) Inspect() string {
	return fmt.Sprintf("interface %s", i.PrintableName())
}

func (i *Interface) Error() string {
	return i.Inspect()
}

func (i *Interface) InstanceVariables() *InstanceVariables {
	return nil
}

func NewInterfaceComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *Interface) bool {
		if x == y {
			return true
		}

		return x.name == y.name &&
			cmp.Equal(x.Constants, y.Constants, *opts...)
	})
}

var InterfaceClass *Class // ::Std::Interface

func initInterface() {
	InterfaceClass = NewClassWithOptions()
	StdModule.AddConstantString("Interface", Ref(InterfaceClass))
	RegisterNativeClass("Std::Interface", "value.InterfaceClass")
}
