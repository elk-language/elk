package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// A simple setter method.
type SetterMethod struct {
	AttributeName value.Symbol
	name          value.Symbol
	frozen        bool
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
	return nil
}

func (*SetterMethod) DirectClass() *value.Class {
	return nil
}

func (*SetterMethod) SingletonClass() *value.Class {
	return nil
}

func (s *SetterMethod) IsFrozen() bool {
	return s.frozen
}

func (s *SetterMethod) SetFrozen() {
	s.frozen = true
}

func (s *SetterMethod) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :setter}", s.AttributeName.Inspect())
}

func (*SetterMethod) InstanceVariables() value.SymbolMap {
	return nil
}

func (s *SetterMethod) Call(self value.Value, val value.Value) (value.Value, value.Value) {
	iv := self.InstanceVariables()
	if iv == nil {
		return nil, value.NewCantAccessInstanceVariablesOnPrimitiveError(self.Inspect())
	}
	iv.Set(s.AttributeName, val)
	return val, nil
}

// Create a new getter method.
func NewSetterMethod(attrName value.Symbol, frozen bool) *SetterMethod {
	return &SetterMethod{
		AttributeName: attrName,
		name:          value.ToSymbol(attrName.ToString() + "="),
		frozen:        frozen,
	}
}

// Creates a setter method and attaches it to
// the given container.
func DefineSetter(
	container *value.MethodContainer,
	attrName value.Symbol,
	frozen bool,
) *value.Error {
	setterMethod := NewSetterMethod(
		attrName,
		frozen,
	)
	return container.AttachMethod(setterMethod.name, setterMethod)
}

// Utility method that creates a new setter and getter method and
// attaches them as methods to the given method map.
func DefineAccessor(
	container *value.MethodContainer,
	attrName value.Symbol,
	frozen bool,
) *value.Error {
	err := DefineGetter(container, attrName, frozen)
	if err != nil {
		return err
	}
	return DefineSetter(container, attrName, frozen)
}

type SetterOption func(*SetterMethod)

func SetterWithFrozen(frozen bool) SetterOption {
	return func(sm *SetterMethod) {
		sm.frozen = frozen
	}
}

// Utility method that creates a new setter method and
// attaches it as a method to the given container.
// Panics when the method can't be defined.
func Setter(
	container *value.MethodContainer,
	attrName string,
	opts ...SetterOption,
) {
	attrNameSymbol := value.ToSymbol(attrName)
	setterMethod := NewSetterMethod(attrNameSymbol, false)

	for _, opt := range opts {
		opt(setterMethod)
	}

	if err := container.AttachMethod(setterMethod.name, setterMethod); err != nil {
		panic(err)
	}
}

// Utility method that creates a new setter and getter method and
// attaches them as methods to the given container.
// Panics when the methods can't be defined.
func Accessor(
	container *value.MethodContainer,
	attrName string,
	frozen bool,
) {
	Getter(container, attrName, GetterWithFrozen(frozen))
	Setter(container, attrName, SetterWithFrozen(frozen))
}
