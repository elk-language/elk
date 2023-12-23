package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// A simple getter method.
type GetterMethod struct {
	AttributeName value.Symbol
	Doc           value.Value
	sealed        bool
}

func (g *GetterMethod) Name() value.Symbol {
	return g.AttributeName
}

func (*GetterMethod) ParameterCount() int {
	return 0
}

func (*GetterMethod) Parameters() []value.Symbol {
	return nil
}

func (*GetterMethod) OptionalParameterCount() int {
	return 0
}

func (*GetterMethod) PostRestParameterCount() int {
	return 0
}

func (*GetterMethod) NamedRestParameter() bool {
	return false
}

func (*GetterMethod) Class() *value.Class {
	return value.MethodClass
}

func (*GetterMethod) DirectClass() *value.Class {
	return value.MethodClass
}

func (*GetterMethod) SingletonClass() *value.Class {
	return nil
}

func (g *GetterMethod) IsSealed() bool {
	return g.sealed
}

func (g *GetterMethod) SetSealed() {
	g.sealed = true
}

func (g *GetterMethod) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :getter}", g.AttributeName.Inspect())
}

func (*GetterMethod) InstanceVariables() value.SymbolMap {
	return nil
}

func (g *GetterMethod) Call(self value.Value) (value.Value, value.Value) {
	iv := self.InstanceVariables()
	if iv == nil {
		return nil, value.NewCantAccessInstanceVariablesOnPrimitiveError(self.Inspect())
	}
	result := iv.Get(g.AttributeName)
	if result == nil {
		return value.Nil, nil
	}
	return result, nil
}

// Create a new getter method.
func NewGetterMethod(attrName value.Symbol, sealed bool) *GetterMethod {
	return &GetterMethod{
		AttributeName: attrName,
		sealed:        sealed,
	}
}

// Creates a getter method and attaches it to
// the given container.
func DefineGetter(
	container *value.MethodContainer,
	name value.Symbol,
	sealed bool,
) *value.Error {
	getterMethod := NewGetterMethod(
		name,
		sealed,
	)
	return container.AttachMethod(name, getterMethod)
}

type GetterOption func(*GetterMethod)

func GetterWithSealed(sealed bool) GetterOption {
	return func(gm *GetterMethod) {
		gm.sealed = sealed
	}
}

// Utility method that creates a new getter method and
// attaches it as a method to the given container.
// It panics when the method can't be defined.
func Getter(
	container *value.MethodContainer,
	name string,
	opts ...GetterOption,
) {
	nameSymbol := value.ToSymbol(name)
	getterMethod := NewGetterMethod(
		nameSymbol,
		false,
	)

	for _, opt := range opts {
		opt(getterMethod)
	}

	if err := container.AttachMethod(nameSymbol, getterMethod); err != nil {
		panic(err)
	}
}
