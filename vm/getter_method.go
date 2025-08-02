package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// A simple getter method.
type GetterMethod struct {
	AttributeName value.Symbol
	Doc           value.Value
	IvarIndex     int
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

func (g *GetterMethod) Copy() value.Reference {
	return g
}

func (g *GetterMethod) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :getter}", g.AttributeName.Inspect())
}

func (g *GetterMethod) Error() string {
	return g.Inspect()
}

func (*GetterMethod) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (g *GetterMethod) Call(self value.Value) (value.Value, value.Value) {
	iv := self.InstanceVariables()
	if iv == nil {
		return value.Undefined, value.Ref(value.NewCantAccessInstanceVariablesOnPrimitiveError(self.Inspect()))
	}

	var ivarIndex int
	if g.IvarIndex != -1 {
		ivarIndex = g.IvarIndex
	} else {
		ivarIndex = self.DirectClass().IvarIndices[g.AttributeName]
	}

	result := iv.Get(ivarIndex)
	if result.IsUndefined() {
		return value.Nil, value.Undefined
	}
	return result, value.Undefined
}

// Create a new getter method.
func NewGetterMethod(attrName value.Symbol, index int) *GetterMethod {
	return &GetterMethod{
		AttributeName: attrName,
		IvarIndex:     index,
	}
}

// Creates a getter method and attaches it to
// the given container.
func DefineGetter(
	container *value.MethodContainer,
	name value.Symbol,
	index int,
) {
	getterMethod := NewGetterMethod(
		name,
		index,
	)
	container.AttachMethod(name, getterMethod)
}

// Utility method that creates a new getter method and
// attaches it as a method to the given container.
// It panics when the method cannot be defined.
func Getter(
	container *value.MethodContainer,
	name string,
) {
	nameSymbol := value.ToSymbol(name)
	getterMethod := NewGetterMethod(
		nameSymbol,
		-1,
	)
	container.AttachMethod(nameSymbol, getterMethod)
}
