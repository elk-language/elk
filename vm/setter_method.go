package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// A simple setter method.
type SetterMethod struct {
	Doc           value.Value
	AttributeName value.Symbol
	name          value.Symbol
}

func (s *SetterMethod) Name() value.Symbol {
	return s.name
}

func (*SetterMethod) ParameterCount() int {
	return 1
}

var setterParams = []value.Symbol{value.ToSymbol("val")}

func (*SetterMethod) Parameters() []value.Symbol {
	return setterParams
}

func (*SetterMethod) OptionalParameterCount() int {
	return 0
}

func (*SetterMethod) PostRestParameterCount() int {
	return 0
}

func (*SetterMethod) NamedRestParameter() bool {
	return false
}

func (*SetterMethod) Class() *value.Class {
	return value.MethodClass
}

func (*SetterMethod) DirectClass() *value.Class {
	return value.MethodClass
}

func (*SetterMethod) SingletonClass() *value.Class {
	return nil
}

func (s *SetterMethod) Copy() value.Reference {
	return s
}

func (s *SetterMethod) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :setter}", s.AttributeName.Inspect())
}

func (s *SetterMethod) Error() string {
	return s.Inspect()
}

func (*SetterMethod) InstanceVariables() value.SymbolMap {
	return nil
}

func (s *SetterMethod) Call(self value.Value, val value.Value) (value.Value, value.Value) {
	iv := self.InstanceVariables()
	if iv == nil {
		return value.Undefined, value.Ref(value.NewCantAccessInstanceVariablesOnPrimitiveError(self.Inspect()))
	}
	iv.Set(s.AttributeName, val)
	return val, value.Undefined
}

// Create a new getter method.
func NewSetterMethod(attrName value.Symbol) *SetterMethod {
	return &SetterMethod{
		AttributeName: attrName,
		name:          value.ToSymbol(attrName.ToString() + "="),
	}
}

// Creates a setter method and attaches it to
// the given container.
func DefineSetter(
	container *value.MethodContainer,
	attrName value.Symbol,
) {
	setterMethod := NewSetterMethod(
		attrName,
	)
	container.AttachMethod(setterMethod.name, setterMethod)
}

// Utility method that creates a new setter and getter method and
// attaches them as methods to the given method map.
func DefineAccessor(
	container *value.MethodContainer,
	attrName value.Symbol,
) {
	DefineGetter(container, attrName)
	DefineSetter(container, attrName)
}

// Utility method that creates a new setter method and
// attaches it as a method to the given container.
// Panics when the method cannot be defined.
func Setter(
	container *value.MethodContainer,
	attrName string,
) {
	attrNameSymbol := value.ToSymbol(attrName)
	setterMethod := NewSetterMethod(attrNameSymbol)

	container.AttachMethod(setterMethod.name, setterMethod)
}

// Utility method that creates a new setter and getter method and
// attaches them as methods to the given container.
// Panics when the methods cannot be defined.
func Accessor(
	container *value.MethodContainer,
	attrName string,
) {
	Getter(container, attrName)
	Setter(container, attrName)
}
