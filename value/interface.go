package value

import (
	"fmt"
	"maps"

	"github.com/google/go-cmp/cmp"
)

// Represents an Elk interface.
type Interface struct {
	ConstantContainer
}

// Interface constructor option function
type InterfaceOption = func(*Interface)

func InterfaceWithName(name string) InterfaceOption {
	return func(i *Interface) {
		i.ConstantContainer.Name = name
	}
}

// Create a new module.
func NewInterface() *Interface {
	return &Interface{
		ConstantContainer: ConstantContainer{
			Constants: make(SymbolMap),
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
	return NewInterface()
}

func (i *Interface) Copy() Value {
	newConstants := make(SymbolMap, len(i.Constants))
	maps.Copy(newConstants, i.Constants)

	newInterface := &Interface{
		ConstantContainer: ConstantContainer{
			Constants: newConstants,
			Name:      i.Name,
		},
	}

	return newInterface
}

func (i *Interface) Class() *Class {
	return InterfaceClass
}

func (i *Interface) DirectClass() *Class {
	return InterfaceClass
}

func (i *Interface) SingletonClass() *Class {
	return nil
}

func (i *Interface) Inspect() string {
	return fmt.Sprintf("interface %s", i.PrintableName())
}

func (i *Interface) InstanceVariables() SymbolMap {
	return nil
}

func NewInterfaceComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *Interface) bool {
		if x == y {
			return true
		}

		return x.Name == y.Name &&
			cmp.Equal(x.Constants, y.Constants, *opts...)
	})
}

var InterfaceClass *Class // ::Std::Interface
